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
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/pkg/resilience/clock"
	"github.com/shing1211/webullapi4go/pkg/resilience/retry"
)

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

func TestFullJitterRetriesWithRealClock(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	t.Parallel()

	r := retry.New(
		retry.WithMaxAttempts(3),
		retry.WithBaseDelay(10*time.Millisecond),
		retry.WithMaxDelay(40*time.Millisecond),
		retry.WithFullJitter(true),
		retry.WithIsRetryable(func(err error) bool { return err != nil }),
	)

	calls := 0
	err := r.Do(context.Background(), func(ctx context.Context) error {
		calls++
		return errors.New("transient")
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if calls != 3 {
		t.Errorf("calls = %d, want 3", calls)
	}
}
