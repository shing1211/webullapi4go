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

// Package ratelimit implements a token-bucket rate limiter and a keyed
// variant for applying independent limits per endpoint.
package ratelimit

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/shing1211/webullapi4go/pkg/resilience/clock"
)

type Limiter struct {
	mu    sync.Mutex
	rate  float64
	burst float64
	toks  float64
	last  time.Time
	clk   clock.Clock
}

type Option func(*config)

type config struct {
	clock clock.Clock
}

func WithClock(clk clock.Clock) Option {
	return func(c *config) { c.clock = clk }
}

func New(rate float64, burst int, opts ...Option) *Limiter {
	if rate <= 0 {
		panic("ratelimit: rate must be positive")
	}
	if burst <= 0 {
		panic("ratelimit: burst must be positive")
	}
	cfg := config{clock: clock.System()}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.clock == nil {
		cfg.clock = clock.System()
	}
	return &Limiter{
		rate:  rate,
		burst: float64(burst),
		toks:  float64(burst),
		last:  cfg.clock.Now(),
		clk:   cfg.clock,
	}
}

func (l *Limiter) refill(now time.Time) float64 {
	elapsed := now.Sub(l.last).Seconds()
	if elapsed > 0 {
		l.toks = math.Min(l.burst, l.toks+elapsed*l.rate)
		l.last = now
	}
	return l.toks
}

func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.refill(l.clk.Now()) >= 1 {
		l.toks--
		return true
	}
	return false
}

func (l *Limiter) Tokens() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.refill(l.clk.Now())
}

func (l *Limiter) Wait(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		l.mu.Lock()
		if err := ctx.Err(); err != nil {
			l.mu.Unlock()
			return err
		}
		now := l.clk.Now()
		if l.refill(now) >= 1 {
			l.toks--
			l.mu.Unlock()
			return nil
		}
		need := (1 - l.toks) / l.rate
		l.mu.Unlock()

		delay := time.Duration(need * float64(time.Second))
		if delay < time.Millisecond {
			delay = time.Millisecond
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-l.clk.After(delay):
		}
	}
}

type Keyed struct {
	mu       sync.Mutex
	factory  func() *Limiter
	limiters map[string]*Limiter
}

func NewKeyed(factory func() *Limiter) *Keyed {
	if factory == nil {
		panic("ratelimit: NewKeyed factory must not be nil")
	}
	return &Keyed{factory: factory, limiters: make(map[string]*Limiter)}
}

func (k *Keyed) For(key string) *Limiter {
	k.mu.Lock()
	defer k.mu.Unlock()
	l, ok := k.limiters[key]
	if !ok {
		l = k.factory()
		k.limiters[key] = l
	}
	return l
}

func (k *Keyed) Set(key string, l *Limiter) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.limiters[key] = l
}

func (k *Keyed) Allow(key string) bool { return k.For(key).Allow() }

func (k *Keyed) Wait(ctx context.Context, key string) error { return k.For(key).Wait(ctx) }
