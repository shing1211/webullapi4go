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

// Package breaker implements the circuit-breaker pattern.
//
// A breaker starts closed and admits every call. Once the number of
// consecutive failures reaches the configured threshold it opens and rejects
// calls immediately with [ErrOpen]. After the cooldown elapses it moves to
// half-open and admits a limited number of probe calls: a success closes the
// breaker, a failure re-opens it. A [Breaker] is safe for concurrent use.
package breaker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/shing1211/webullapi4go/internal/resilience/clock"
)

// ErrOpen is returned by [Breaker.Allow] consumers (and [Breaker.Do]) when the
// circuit is open.
var ErrOpen = errors.New("circuit breaker is open")

// State is a circuit-breaker state.
type State int

// The breaker states.
const (
	// StateClosed admits every call.
	StateClosed State = iota
	// StateOpen rejects every call until the cooldown elapses.
	StateOpen
	// StateHalfOpen admits probe calls to test recovery.
	StateHalfOpen
)

// String returns the human-readable state name.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Defaults applied by [New].
const (
	// DefaultThreshold is the number of consecutive failures that opens the
	// breaker.
	DefaultThreshold = 5
	// DefaultCooldown is how long the breaker stays open before probing.
	DefaultCooldown = 30 * time.Second
	// DefaultHalfOpenMax is the number of concurrent probe calls allowed in
	// the half-open state.
	DefaultHalfOpenMax = 1
)

// Config configures a [Breaker]. The zero value is completed with defaults by
// [New].
type Config struct {
	// Threshold is the number of consecutive failures that opens the breaker.
	Threshold int
	// Cooldown is how long the breaker stays open before allowing probes.
	Cooldown time.Duration
	// HalfOpenMax is the number of concurrent probe calls admitted while
	// half-open.
	HalfOpenMax int
	// OnStateChange, when set, is invoked after a state transition with the
	// previous and new states. It is called while the breaker's lock is held,
	// so it must not call back into the breaker.
	OnStateChange func(from, to State)
	// Clock is the time source. When nil, the real-time clock is used.
	Clock clock.Clock
}

// Option mutates a [Config]. nil options are ignored.
type Option func(*Config)

// WithThreshold sets the failure threshold that opens the breaker.
func WithThreshold(n int) Option {
	return func(c *Config) { c.Threshold = n }
}

// WithCooldown sets how long the breaker stays open before probing.
func WithCooldown(d time.Duration) Option {
	return func(c *Config) { c.Cooldown = d }
}

// WithHalfOpenMax sets the number of concurrent half-open probe calls.
func WithHalfOpenMax(n int) Option {
	return func(c *Config) { c.HalfOpenMax = n }
}

// WithOnStateChange registers a state-transition callback.
func WithOnStateChange(fn func(from, to State)) Option {
	return func(c *Config) { c.OnStateChange = fn }
}

// WithClock overrides the time source. A nil clock is replaced by the
// real-time clock.
func WithClock(clk clock.Clock) Option {
	return func(c *Config) { c.Clock = clk }
}

// Breaker is a thread-safe circuit breaker.
type Breaker struct {
	mu               sync.Mutex
	state            State
	failures         int
	threshold        int
	cooldown         time.Duration
	halfOpenMax      int
	halfOpenInFlight int
	openedAt         time.Time
	onStateChange    func(from, to State)
	clk              clock.Clock
}

// New returns a [Breaker] configured by opts. Zero-valued fields fall back to
// the package defaults.
func New(opts ...Option) *Breaker {
	cfg := Config{
		Threshold:   DefaultThreshold,
		Cooldown:    DefaultCooldown,
		HalfOpenMax: DefaultHalfOpenMax,
		Clock:       clock.System(),
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return NewWithConfig(cfg)
}

// NewWithConfig returns a [Breaker] for cfg, filling zero-valued fields from
// the package defaults.
func NewWithConfig(cfg Config) *Breaker {
	if cfg.Threshold < 1 {
		cfg.Threshold = DefaultThreshold
	}
	if cfg.Cooldown <= 0 {
		cfg.Cooldown = DefaultCooldown
	}
	if cfg.HalfOpenMax < 1 {
		cfg.HalfOpenMax = DefaultHalfOpenMax
	}
	if cfg.Clock == nil {
		cfg.Clock = clock.System()
	}
	return &Breaker{
		state:         StateClosed,
		threshold:     cfg.Threshold,
		cooldown:      cfg.Cooldown,
		halfOpenMax:   cfg.HalfOpenMax,
		onStateChange: cfg.OnStateChange,
		clk:           cfg.Clock,
	}
}

// State returns the current state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// Failures returns the number of consecutive recorded failures.
func (b *Breaker) Failures() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.failures
}

// Allow reports whether a call may proceed. In the closed and half-open states
// it reserves a probe slot as needed.
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateClosed:
		return true
	case StateOpen:
		if b.clk.Now().Sub(b.openedAt) >= b.cooldown {
			b.transitionTo(StateHalfOpen)
			b.halfOpenInFlight++
			return true
		}
		return false
	case StateHalfOpen:
		if b.halfOpenInFlight < b.halfOpenMax {
			b.halfOpenInFlight++
			return true
		}
		return false
	default:
		return false
	}
}

// RecordSuccess records a successful call. A success while half-open closes
// the breaker; while closed it decrements the failure count.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateClosed:
		if b.failures > 0 {
			b.failures--
		}
	case StateHalfOpen:
		if b.halfOpenInFlight > 0 {
			b.halfOpenInFlight--
		}
		b.transitionTo(StateClosed)
	default:
	}
}

// RecordFailure records a failed call. When the threshold is reached while
// closed, or when a half-open probe fails, the breaker opens.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateClosed:
		b.failures++
		if b.failures >= b.threshold {
			b.transitionTo(StateOpen)
		}
	case StateHalfOpen:
		if b.halfOpenInFlight > 0 {
			b.halfOpenInFlight--
		}
		b.transitionTo(StateOpen)
	default:
	}
}

// Reset closes the breaker and clears the failure count.
func (b *Breaker) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.transitionTo(StateClosed)
}

// Do runs fn when the circuit admits the call, recording the outcome. It
// returns [ErrOpen] without invoking fn when the circuit rejects the call.
func (b *Breaker) Do(ctx context.Context, fn func(context.Context) error) error {
	if !b.Allow() {
		return ErrOpen
	}
	err := fn(ctx)
	if err != nil {
		b.RecordFailure()
	} else {
		b.RecordSuccess()
	}
	return err
}

// transitionTo changes state and fires the change callback. The caller must
// hold b.mu.
func (b *Breaker) transitionTo(next State) {
	if b.state == next {
		return
	}
	prev := b.state
	b.state = next
	switch next {
	case StateOpen:
		b.openedAt = b.clk.Now()
	case StateHalfOpen:
		b.halfOpenInFlight = 0
	case StateClosed:
		b.failures = 0
		b.halfOpenInFlight = 0
	}
	if b.onStateChange != nil {
		b.onStateChange(prev, next)
	}
}
