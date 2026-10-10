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
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/pkg/resilience/clock"
	"github.com/shing1211/webullapi4go/pkg/resilience/ratelimit"
)

type observedFakeClock struct {
	*clock.Fake
	afterStarted chan struct{}
	afterOnce    sync.Once
}

func (c *observedFakeClock) After(d time.Duration) <-chan time.Time {
	ch := c.Fake.After(d)
	c.afterOnce.Do(func() { close(c.afterStarted) })
	return ch
}

func TestLimiterRefillsWithFakeClock(t *testing.T) {
	t.Parallel()

	fake := clock.NewFake(time.Unix(1_700_000_000, 0))
	limiter := ratelimit.New(1, 1, ratelimit.WithClock(fake))
	if !limiter.Allow() {
		t.Fatal("Allow() = false for the initial token")
	}
	if limiter.Allow() {
		t.Fatal("Allow() = true before the bucket refills")
	}

	fake.Advance(999 * time.Millisecond)
	if limiter.Allow() {
		t.Fatal("Allow() = true before one full token is available")
	}
	fake.Advance(time.Millisecond)
	if !limiter.Allow() {
		t.Fatal("Allow() = false after one token refilled")
	}
}

func TestKeyedLimiterKeepsKeysIsolated(t *testing.T) {
	t.Parallel()

	fake := clock.NewFake(time.Unix(1_700_000_000, 0))
	keyed := ratelimit.NewKeyed(func() *ratelimit.Limiter {
		return ratelimit.New(1, 1, ratelimit.WithClock(fake))
	})

	if !keyed.Allow("alpha") {
		t.Fatal("alpha consumed no initial token")
	}
	if !keyed.Allow("beta") {
		t.Fatal("beta was limited by alpha")
	}
	if keyed.Allow("alpha") {
		t.Fatal("alpha consumed more than its burst")
	}
	if !keyed.Allow("gamma") {
		t.Fatal("gamma was limited by another key")
	}
}

func TestLimiterWaitRejectsAlreadyCancelledContext(t *testing.T) {
	t.Parallel()

	limiter := ratelimit.New(1, 1)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := limiter.Wait(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Wait() error = %v, want context.Canceled", err)
	}
	if got := limiter.Tokens(); got != 1 {
		t.Fatalf("Tokens() = %v after cancelled Wait, want 1", got)
	}
}

func TestLimiterWaitCancellationInterruptsRefill(t *testing.T) {
	t.Parallel()

	fake := &observedFakeClock{
		Fake:         clock.NewFake(time.Unix(1_700_000_000, 0)),
		afterStarted: make(chan struct{}),
	}
	limiter := ratelimit.New(1, 1, ratelimit.WithClock(fake))
	if !limiter.Allow() {
		t.Fatal("Allow() = false for the initial token")
	}

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() { result <- limiter.Wait(ctx) }()
	waitForRateLimitSignal(t, fake.afterStarted, "refill wait")
	cancel()

	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Wait() error = %v, want context.Canceled", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Wait did not return after cancellation")
	}
	if got := limiter.Tokens(); got != 0 {
		t.Fatalf("Tokens() = %v during refill wait, want 0", got)
	}
	fake.Advance(time.Second)
	if got := limiter.Tokens(); got != 1 {
		t.Fatalf("Tokens() = %v after refill, want 1", got)
	}
}

func waitForRateLimitSignal(t *testing.T, ch <-chan struct{}, what string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for %s", what)
	}
}
