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

package retry_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/pkg/resilience/clock"
	"github.com/shing1211/webullapi4go/pkg/resilience/retry"
)

type backoffBarrierClock struct {
	now          time.Time
	afterStarted chan struct{}
	afterOnce    sync.Once
	after        chan time.Time
}

func newBackoffBarrierClock(now time.Time) *backoffBarrierClock {
	return &backoffBarrierClock{
		now:          now,
		afterStarted: make(chan struct{}),
		after:        make(chan time.Time),
	}
}

func (c *backoffBarrierClock) Now() time.Time {
	return c.now
}

func (c *backoffBarrierClock) After(time.Duration) <-chan time.Time {
	c.afterOnce.Do(func() { close(c.afterStarted) })
	return c.after
}

func TestFullJitterNoRetryOnPermanent(t *testing.T) {
	t.Parallel()

	r := retry.New(
		retry.WithMaxAttempts(3),
		retry.WithBaseDelay(5*time.Millisecond),
		retry.WithMaxDelay(20*time.Millisecond),
		retry.WithFullJitter(true),
	)

	calls := 0
	err := r.Do(context.Background(), func(ctx context.Context) error {
		calls++
		return retry.Permanent(errors.New("permanent"))
	})
	if !retry.IsPermanent(err) {
		t.Errorf("expected permanent error, got %v", err)
	}
	if calls != 1 {
		t.Errorf("calls = %d, want 1", calls)
	}
}

func TestFullJitterRespectsContextCancel(t *testing.T) {
	t.Parallel()

	fake := clock.NewFake(time.Now())
	r := retry.New(
		retry.WithMaxAttempts(5),
		retry.WithBaseDelay(time.Second),
		retry.WithMaxDelay(4*time.Second),
		retry.WithFullJitter(true),
		retry.WithClock(fake),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	err := r.Do(ctx, func(ctx context.Context) error {
		calls++
		return errors.New("transient")
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if calls != 0 {
		t.Errorf("calls = %d, want 0 (cancelled before first attempt)", calls)
	}
}

func TestCancellationInterruptsBackoffWait(t *testing.T) {
	t.Parallel()

	fake := newBackoffBarrierClock(time.Now())
	r := retry.New(
		retry.WithMaxAttempts(3),
		retry.WithBaseDelay(time.Hour),
		retry.WithMaxDelay(time.Hour),
		retry.WithJitter(false),
		retry.WithIsRetryable(func(error) bool { return true }),
		retry.WithClock(fake),
	)

	ctx, cancel := context.WithCancel(context.Background())
	attempted := make(chan struct{})
	result := make(chan error, 1)
	go func() {
		result <- r.Do(ctx, func(context.Context) error {
			close(attempted)
			return errors.New("transient")
		})
	}()

	waitForRetrySignal(t, attempted, "first attempt")
	waitForRetrySignal(t, fake.afterStarted, "backoff wait")
	cancel()

	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Do() error = %v, want context.Canceled", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Do did not return after cancellation")
	}
}

func waitForRetrySignal(t *testing.T, ch <-chan struct{}, what string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for %s", what)
	}
}
