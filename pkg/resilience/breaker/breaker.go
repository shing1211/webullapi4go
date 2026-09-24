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
package breaker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/shing1211/webullapi4go/pkg/resilience/clock"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var ErrOpen = errors.New("circuit breaker is open")

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

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

const (
	DefaultThreshold   = 5
	DefaultCooldown    = 30 * time.Second
	DefaultHalfOpenMax = 1
)

type Config struct {
	Threshold     int
	Cooldown      time.Duration
	HalfOpenMax   int
	OnStateChange func(from, to State)
	Clock         clock.Clock
	Meter         metric.Meter
}

type Option func(*Config)

func WithThreshold(n int) Option {
	return func(c *Config) { c.Threshold = n }
}

func WithCooldown(d time.Duration) Option {
	return func(c *Config) { c.Cooldown = d }
}

func WithHalfOpenMax(n int) Option {
	return func(c *Config) { c.HalfOpenMax = n }
}

func WithMeter(m metric.Meter) Option {
	return func(c *Config) { c.Meter = m }
}

func WithOnStateChange(fn func(from, to State)) Option {
	return func(c *Config) { c.OnStateChange = fn }
}

func WithClock(clk clock.Clock) Option {
	return func(c *Config) { c.Clock = clk }
}

type Breaker struct {
	mu                sync.Mutex
	state             State
	failures          int
	threshold         int
	cooldown          time.Duration
	halfOpenMax       int
	halfOpenInFlight  int
	openedAt          time.Time
	onStateChange     func(from, to State)
	clk               clock.Clock
	transitionCounter metric.Int64Counter
}

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
	b := &Breaker{
		state:         StateClosed,
		threshold:     cfg.Threshold,
		cooldown:      cfg.Cooldown,
		halfOpenMax:   cfg.HalfOpenMax,
		onStateChange: cfg.OnStateChange,
		clk:           cfg.Clock,
	}
	if cfg.Meter != nil {
		b.transitionCounter, _ = cfg.Meter.Int64Counter(
			"breaker.state_transitions",
			metric.WithDescription("Circuit breaker state transitions"),
		)
	}
	return b
}

func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

func (b *Breaker) Failures() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.failures
}

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

func (b *Breaker) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.transitionTo(StateClosed)
}

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
	if b.transitionCounter != nil {
		b.transitionCounter.Add(context.Background(), 1,
			metric.WithAttributes(
				attribute.String("from", prev.String()),
				attribute.String("to", next.String()),
			))
	}
}
