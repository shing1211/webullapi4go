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

package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/internal/resilience/clock"
	"github.com/shing1211/webullapi4go/internal/resilience/ratelimit"
)

func TestAllowConsumesAndRefills(t *testing.T) {
	t.Parallel()

	clk := clock.NewFake(time.Unix(0, 0))
	l := ratelimit.New(10, 2, ratelimit.WithClock(clk))

	first := l.Allow()
	second := l.Allow()
	if !first || !second {
		t.Fatal("first burst of two tokens should be allowed")
	}
	if l.Allow() {
		t.Fatal("third call within the same instant should be rejected")
	}

	clk.Advance(100 * time.Millisecond) // one token at 10/s
	if !l.Allow() {
		t.Fatal("call after refill should be allowed")
	}
	if l.Allow() {
		t.Fatal("second call after a single refill should be rejected")
	}

	clk.Advance(time.Second)
	if got := l.Tokens(); got != 2 {
		t.Fatalf("Tokens() = %v, want capped at 2", got)
	}
}

func TestWaitBlocksUntilToken(t *testing.T) {
	t.Parallel()

	clk := clock.NewFake(time.Unix(0, 0))
	l := ratelimit.New(10, 1, ratelimit.WithClock(clk)) // token every 100ms
	if !l.Allow() {
		t.Fatal("initial token should be allowed")
	}

	result := make(chan error, 1)
	done := make(chan struct{})
	go func() {
		result <- l.Wait(context.Background())
		close(done)
	}()

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-done:
		case <-ticker.C:
			clk.Advance(time.Second)
			continue
		}
		break
	}
	if err := <-result; err != nil {
		t.Fatalf("Wait() error = %v, want nil", err)
	}
}

func TestWaitHonoursContext(t *testing.T) {
	t.Parallel()

	clk := clock.NewFake(time.Unix(0, 0))
	l := ratelimit.New(1, 1, ratelimit.WithClock(clk))
	if !l.Allow() {
		t.Fatal("initial token should be allowed")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := l.Wait(ctx); err != context.Canceled {
		t.Fatalf("Wait() error = %v, want context.Canceled", err)
	}
}

func TestKeyedLimitsAreIndependent(t *testing.T) {
	t.Parallel()

	clk := clock.NewFake(time.Unix(0, 0))
	k := ratelimit.NewKeyed(func() *ratelimit.Limiter {
		return ratelimit.New(1, 1, ratelimit.WithClock(clk))
	})

	if !k.Allow("a") {
		t.Fatal("first call for key a should be allowed")
	}
	if k.Allow("a") {
		t.Fatal("second call for key a should be rejected")
	}
	if !k.Allow("b") {
		t.Fatal("first call for key b should be allowed")
	}
}

func TestKeyedSetOverridesFactory(t *testing.T) {
	t.Parallel()

	clk := clock.NewFake(time.Unix(0, 0))
	k := ratelimit.NewKeyed(func() *ratelimit.Limiter {
		return ratelimit.New(1, 1, ratelimit.WithClock(clk))
	})
	k.Set("slow", ratelimit.New(10, 5, ratelimit.WithClock(clk)))

	for i := 0; i < 5; i++ {
		if !k.Allow("slow") {
			t.Fatalf("call %d for overridden key should be allowed", i)
		}
	}
	if k.Allow("slow") {
		t.Fatal("sixth call for overridden key should be rejected")
	}
}

func TestNewPanicsOnInvalidConfig(t *testing.T) {
	t.Parallel()

	assertPanics := func(name string, fn func()) {
		t.Helper()
		defer func() {
			if recover() == nil {
				t.Fatalf("%s: expected panic", name)
			}
		}()
		fn()
	}
	assertPanics("zero rate", func() { ratelimit.New(0, 1) })
	assertPanics("zero burst", func() { ratelimit.New(1, 0) })
	assertPanics("nil factory", func() { ratelimit.NewKeyed(nil) })
}
