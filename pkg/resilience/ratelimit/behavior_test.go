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

type limiterWaitClock struct {
	*clock.Fake
	afterStarted chan struct{}
	afterOnce    sync.Once
}

func (c *limiterWaitClock) After(d time.Duration) <-chan time.Time {
	ch := c.Fake.After(d)
	c.afterOnce.Do(func() { close(c.afterStarted) })
	return ch
}

func TestWaitAcquiresRefilledToken(t *testing.T) {
	t.Parallel()

	fake := &limiterWaitClock{
		Fake:         clock.NewFake(time.Unix(1_700_000_000, 0)),
		afterStarted: make(chan struct{}),
	}
	limiter := ratelimit.New(10, 1, ratelimit.WithClock(fake))
	if !limiter.Allow() {
		t.Fatal("Allow() = false for the initial token")
	}

	result := make(chan error, 1)
	go func() { result <- limiter.Wait(context.Background()) }()
	select {
	case <-fake.afterStarted:
	case <-time.After(5 * time.Second):
		t.Fatal("Wait did not start a refill wait")
	}
	fake.Advance(100 * time.Millisecond)

	select {
	case err := <-result:
		if err != nil {
			t.Fatalf("Wait() error = %v, want nil", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Wait did not return after the fake clock advanced")
	}
	if got := limiter.Tokens(); got != 0 {
		t.Fatalf("Tokens() = %v after Wait, want 0", got)
	}
}

func TestNewRejectsInvalidConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		call func()
	}{
		{name: "zero rate", call: func() { ratelimit.New(0, 1) }},
		{name: "negative rate", call: func() { ratelimit.New(-1, 1) }},
		{name: "zero burst", call: func() { ratelimit.New(1, 0) }},
		{name: "negative burst", call: func() { ratelimit.New(1, -1) }},
		{name: "nil keyed factory", call: func() { ratelimit.NewKeyed(nil) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				if recover() == nil {
					t.Fatal("constructor did not panic")
				}
			}()
			tt.call()
		})
	}
}

func TestNewAcceptsNilOptionsAndClock(t *testing.T) {
	t.Parallel()

	limiter := ratelimit.New(1, 1, nil, ratelimit.WithClock(nil))
	if !limiter.Allow() {
		t.Fatal("Allow() = false for the initial token")
	}
	if limiter.Allow() {
		t.Fatal("Allow() = true before the system clock advances")
	}
}

func TestKeyedReusesFactoryLimiterAndSetReplacesIt(t *testing.T) {
	t.Parallel()

	fake := clock.NewFake(time.Unix(0, 0))
	factoryCalls := 0
	keyed := ratelimit.NewKeyed(func() *ratelimit.Limiter {
		factoryCalls++
		return ratelimit.New(1, 1, ratelimit.WithClock(fake))
	})

	first := keyed.For("alpha")
	if second := keyed.For("alpha"); second != first {
		t.Fatal("For() created more than one limiter for the same key")
	}
	if factoryCalls != 1 {
		t.Fatalf("factory calls = %d, want 1", factoryCalls)
	}
	if !first.Allow() || first.Allow() {
		t.Fatal("factory limiter did not enforce its one-token burst")
	}

	replacement := ratelimit.New(10, 2, ratelimit.WithClock(fake))
	keyed.Set("alpha", replacement)
	if got := keyed.For("alpha"); got != replacement {
		t.Fatal("For() did not return the limiter installed by Set")
	}
	firstReplacement := keyed.Allow("alpha")
	secondReplacement := keyed.Allow("alpha")
	thirdReplacement := keyed.Allow("alpha")
	if !firstReplacement || !secondReplacement || thirdReplacement {
		t.Fatal("Allow() did not delegate to the replacement limiter")
	}
}

func TestKeyedWaitDelegatesCancellation(t *testing.T) {
	t.Parallel()

	keyed := ratelimit.NewKeyed(func() *ratelimit.Limiter { return ratelimit.New(1, 1) })
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := keyed.Wait(ctx, "cancelled"); !errors.Is(err, context.Canceled) {
		t.Fatalf("Wait() error = %v, want context.Canceled", err)
	}
	if got := keyed.For("cancelled").Tokens(); got != 1 {
		t.Fatalf("Tokens() = %v after cancelled Wait, want 1", got)
	}
}
