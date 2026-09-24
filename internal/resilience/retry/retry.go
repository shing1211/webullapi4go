// Copyright 2026 shing1211
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package retry provides context-aware retries with exponential backoff and
// optional jitter.
//
// A [Retrier] retries a function until it succeeds, returns a non-retryable
// error, or exhausts its attempt budget. Retryability is delegated to an
// [Config.IsRetryable] predicate; [DefaultIsRetryable] treats timeouts, HTTP
// 429 responses, and HTTP 5xx responses as transient. Errors wrapped with
// [Permanent] are never retried.
package retry

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/shing1211/webullapi4go/internal/resilience/clock"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

// Defaults applied by [DefaultConfig].
const (
	// DefaultMaxAttempts is the total number of attempts, including the first.
	DefaultMaxAttempts = 2
	// DefaultBaseDelay is the delay before the first retry.
	DefaultBaseDelay = 100 * time.Millisecond
	// DefaultMaxDelay caps the exponentially growing delay.
	DefaultMaxDelay = 2 * time.Second
)

// Config configures a [Retrier]. The zero value is completed with defaults by
// [New].
type Config struct {
	// MaxAttempts is the total number of attempts, including the first. It is
	// clamped to at least 1.
	MaxAttempts int
	// BaseDelay is the delay before the first retry and the base of the
	// exponential backoff.
	BaseDelay time.Duration
	// MaxDelay caps the computed delay.
	MaxDelay time.Duration
	// Jitter randomizes each delay by plus or minus 20 percent.
	Jitter bool
	// IsRetryable reports whether err should be retried. When nil,
	// [DefaultIsRetryable] is used.
	IsRetryable func(error) bool
	// Clock is the time source. When nil, the real-time clock is used.
	Clock clock.Clock
}

// Option mutates a [Config]. nil options are ignored.
type Option func(*Config)

// DefaultConfig returns a [Config] pre-filled with the package defaults.
func DefaultConfig() Config {
	return Config{
		MaxAttempts: DefaultMaxAttempts,
		BaseDelay:   DefaultBaseDelay,
		MaxDelay:    DefaultMaxDelay,
		Jitter:      true,
		IsRetryable: DefaultIsRetryable,
		Clock:       clock.System(),
	}
}

// WithMaxAttempts sets the total number of attempts, including the first.
func WithMaxAttempts(n int) Option {
	return func(c *Config) { c.MaxAttempts = n }
}

// WithBaseDelay sets the delay before the first retry.
func WithBaseDelay(d time.Duration) Option {
	return func(c *Config) { c.BaseDelay = d }
}

// WithMaxDelay caps the computed delay.
func WithMaxDelay(d time.Duration) Option {
	return func(c *Config) { c.MaxDelay = d }
}

// WithJitter enables or disables jitter.
func WithJitter(enabled bool) Option {
	return func(c *Config) { c.Jitter = enabled }
}

// WithIsRetryable overrides the retryability predicate. A nil predicate is
// replaced by [DefaultIsRetryable].
func WithIsRetryable(fn func(error) bool) Option {
	return func(c *Config) { c.IsRetryable = fn }
}

// WithClock overrides the time source. A nil clock is replaced by the
// real-time clock.
func WithClock(clk clock.Clock) Option {
	return func(c *Config) { c.Clock = clk }
}

// Retrier executes a function with retries. It is safe for concurrent use.
type Retrier struct {
	attempts    int
	baseDelay   time.Duration
	maxDelay    time.Duration
	jitter      bool
	isRetryable func(error) bool
	clk         clock.Clock
}

// New returns a [Retrier] configured by opts. Zero-valued fields fall back to
// [DefaultConfig].
func New(opts ...Option) *Retrier {
	cfg := DefaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return NewWithConfig(cfg)
}

// NewWithConfig returns a [Retrier] for cfg, filling zero-valued fields from
// [DefaultConfig].
func NewWithConfig(cfg Config) *Retrier {
	if cfg.MaxAttempts < 1 {
		cfg.MaxAttempts = 1
	}
	if cfg.BaseDelay < 0 {
		cfg.BaseDelay = 0
	}
	if cfg.MaxDelay < cfg.BaseDelay {
		cfg.MaxDelay = cfg.BaseDelay
	}
	if cfg.IsRetryable == nil {
		cfg.IsRetryable = DefaultIsRetryable
	}
	if cfg.Clock == nil {
		cfg.Clock = clock.System()
	}
	return &Retrier{
		attempts:    cfg.MaxAttempts,
		baseDelay:   cfg.BaseDelay,
		maxDelay:    cfg.MaxDelay,
		jitter:      cfg.Jitter,
		isRetryable: cfg.IsRetryable,
		clk:         cfg.Clock,
	}
}

// MaxAttempts returns the configured attempt budget.
func (r *Retrier) MaxAttempts() int { return r.attempts }

// Do invokes fn until it succeeds, returns a non-retryable error, or the
// attempt budget is exhausted. It observes ctx cancellation between attempts
// and while waiting for a backoff delay.
//
// On exhaustion the last error is returned. When ctx is canceled or its
// deadline passes, ctx.Err() is returned.
func (r *Retrier) Do(ctx context.Context, fn func(context.Context) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	var last error
	for attempt := 0; attempt < r.attempts; attempt++ {
		if attempt > 0 {
			if err := r.wait(ctx, r.delay(attempt-1)); err != nil {
				return err
			}
		}
		err := fn(ctx)
		if err == nil {
			return nil
		}
		last = err
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		if !r.isRetryable(err) {
			return err
		}
	}
	return last
}

// wait blocks for d or until ctx is done, whichever happens first.
func (r *Retrier) wait(ctx context.Context, d time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if d <= 0 {
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.clk.After(d):
		return nil
	}
}

// delay returns the backoff before the retry that follows attempt index
// attempt (zero-based among retries).
func (r *Retrier) delay(attempt int) time.Duration {
	d := float64(r.baseDelay) * math.Pow(2, float64(attempt))
	if max := float64(r.maxDelay); d > max {
		d = max
	}
	if r.jitter {
		d *= 0.8 + 0.4*rand.Float64()
	}
	if d < 0 {
		d = 0
	}
	return time.Duration(d)
}

// permanentError marks err as non-retryable.
type permanentError struct{ err error }

// Error implements the error interface.
func (e *permanentError) Error() string { return e.err.Error() }

// Unwrap returns the wrapped cause.
func (e *permanentError) Unwrap() error { return e.err }

// Permanent marks err as non-retryable. A nil error is returned unchanged.
func Permanent(err error) error {
	if err == nil {
		return nil
	}
	return &permanentError{err: err}
}

// IsPermanent reports whether err was wrapped with [Permanent].
func IsPermanent(err error) bool {
	var p *permanentError
	return errors.As(err, &p)
}

// DefaultIsRetryable reports whether err is transient. Transport failures
// (including timeouts), HTTP 429 responses, and HTTP 5xx responses are
// retryable; everything else, including errors marked with [Permanent], is
// not.
func DefaultIsRetryable(err error) bool {
	if err == nil || IsPermanent(err) {
		return false
	}
	var e *errs.Error
	if errors.As(err, &e) {
		switch e.Code {
		case errs.CodeRateLimited, errs.CodeServer, errs.CodeTransport:
			return true
		}
		switch {
		case e.Status == http.StatusTooManyRequests, e.Status >= 500:
			return true
		}
		return false
	}
	if errors.Is(err, context.Canceled) {
		return false
	}
	var ne net.Error
	if errors.As(err, &ne) {
		return ne.Timeout()
	}
	return false
}
