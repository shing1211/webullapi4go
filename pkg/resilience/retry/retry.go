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
package retry

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"net"
	"net/http"
	"time"

	pkgerrs "github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/pkg/resilience/clock"
)

const (
	DefaultMaxAttempts = 2
	DefaultBaseDelay   = 100 * time.Millisecond
	DefaultMaxDelay    = 2 * time.Second
)

type Config struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Jitter      bool
	FullJitter  bool
	IsRetryable func(error) bool
	Clock       clock.Clock
}

type Option func(*Config)

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

func WithMaxAttempts(n int) Option {
	return func(c *Config) { c.MaxAttempts = n }
}

func WithBaseDelay(d time.Duration) Option {
	return func(c *Config) { c.BaseDelay = d }
}

func WithMaxDelay(d time.Duration) Option {
	return func(c *Config) { c.MaxDelay = d }
}

func WithJitter(enabled bool) Option {
	return func(c *Config) { c.Jitter = enabled }
}

func WithFullJitter(enabled bool) Option {
	return func(c *Config) {
		c.FullJitter = enabled
		if enabled {
			c.Jitter = true
		}
	}
}

func WithIsRetryable(fn func(error) bool) Option {
	return func(c *Config) { c.IsRetryable = fn }
}

func WithClock(clk clock.Clock) Option {
	return func(c *Config) { c.Clock = clk }
}

type Retrier struct {
	attempts    int
	baseDelay   time.Duration
	maxDelay    time.Duration
	jitter      bool
	fullJitter  bool
	isRetryable func(error) bool
	clk         clock.Clock
}

func New(opts ...Option) *Retrier {
	cfg := DefaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return NewWithConfig(cfg)
}

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
		fullJitter:  cfg.FullJitter,
		isRetryable: cfg.IsRetryable,
		clk:         cfg.Clock,
	}
}

func (r *Retrier) MaxAttempts() int { return r.attempts }

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

func (r *Retrier) delay(attempt int) time.Duration {
	cap := float64(r.baseDelay) * math.Pow(2, float64(attempt))
	if max := float64(r.maxDelay); cap > max {
		cap = max
	}
	var d float64
	if r.fullJitter {
		d = rand.Float64() * cap
	} else if r.jitter {
		d = cap * (0.8 + 0.4*rand.Float64())
	} else {
		d = cap
	}
	if d < 0 {
		d = 0
	}
	return time.Duration(d)
}

type permanentError struct{ err error }

func (e *permanentError) Error() string { return e.err.Error() }

func (e *permanentError) Unwrap() error { return e.err }

func Permanent(err error) error {
	if err == nil {
		return nil
	}
	return &permanentError{err: err}
}

func IsPermanent(err error) bool {
	var p *permanentError
	return errors.As(err, &p)
}

func DefaultIsRetryable(err error) bool {
	if err == nil || IsPermanent(err) {
		return false
	}
	var e *pkgerrs.Error
	if errors.As(err, &e) {
		switch e.Code {
		case pkgerrs.CodeRateLimited, pkgerrs.CodeServer, pkgerrs.CodeTransport:
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
