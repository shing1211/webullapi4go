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
	"net"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/internal/resilience/clock"
	"github.com/shing1211/webullapi4go/internal/resilience/retry"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

// advanceUntil drives the fake clock forward until done is closed.
func advanceUntil(clk *clock.Fake, done <-chan struct{}) {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			clk.Advance(time.Second)
		}
	}
}

// runRetrier runs r.Do on a goroutine while driving the fake clock, and
// returns the result.
func runRetrier(clk *clock.Fake, r *retry.Retrier, fn func(context.Context) error) error {
	result := make(chan error, 1)
	done := make(chan struct{})
	go func() {
		result <- r.Do(context.Background(), fn)
		close(done)
	}()
	advanceUntil(clk, done)
	return <-result
}

func TestDoRetriesUntilSuccess(t *testing.T) {
	t.Parallel()

	clk := clock.NewFake(time.Unix(0, 0))
	r := retry.New(
		retry.WithClock(clk),
		retry.WithJitter(false),
		retry.WithMaxAttempts(3),
		retry.WithBaseDelay(time.Second),
	)

	calls := 0
	err := runRetrier(clk, r, func(context.Context) error {
		calls++
		if calls < 3 {
			return errs.New(errs.CodeServer, "boom")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Do() error = %v, want nil", err)
	}
	if calls != 3 {
		t.Fatalf("calls = %d, want 3", calls)
	}
}

func TestDoExhaustsAttempts(t *testing.T) {
	t.Parallel()

	clk := clock.NewFake(time.Unix(0, 0))
	r := retry.New(
		retry.WithClock(clk),
		retry.WithJitter(false),
		retry.WithMaxAttempts(2),
		retry.WithBaseDelay(time.Second),
	)

	calls := 0
	err := runRetrier(clk, r, func(context.Context) error {
		calls++
		return errs.New(errs.CodeServer, "boom")
	})
	if err == nil {
		t.Fatal("Do() error = nil, want exhausted error")
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
}

func TestDoStopsOnNonRetryable(t *testing.T) {
	t.Parallel()

	r := retry.New(
		retry.WithMaxAttempts(5),
		retry.WithBaseDelay(0),
		retry.WithJitter(false),
	)
	calls := 0
	err := r.Do(context.Background(), func(context.Context) error {
		calls++
		return errs.New(errs.CodeAPI, "bad request")
	})
	if err == nil {
		t.Fatal("Do() error = nil, want error")
	}
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
}

func TestDoTreatsPermanentAsNonRetryable(t *testing.T) {
	t.Parallel()

	r := retry.New(retry.WithMaxAttempts(5), retry.WithBaseDelay(0))
	calls := 0
	err := r.Do(context.Background(), func(context.Context) error {
		calls++
		return retry.Permanent(errs.New(errs.CodeServer, "do not retry"))
	})
	if err == nil {
		t.Fatal("Do() error = nil, want error")
	}
	if !retry.IsPermanent(err) {
		t.Fatalf("IsPermanent(%v) = false, want true", err)
	}
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
}

func TestDoHonoursContextCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := retry.New(retry.WithMaxAttempts(5))
	calls := 0
	err := r.Do(ctx, func(context.Context) error {
		calls++
		cancel()
		return errs.New(errs.CodeServer, "boom")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Do() error = %v, want context.Canceled", err)
	}
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
}

func TestDoObservesContextBeforeFirstAttempt(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	r := retry.New(retry.WithMaxAttempts(3))
	calls := 0
	err := r.Do(ctx, func(context.Context) error {
		calls++
		return nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Do() error = %v, want context.Canceled", err)
	}
	if calls != 0 {
		t.Fatalf("calls = %d, want 0", calls)
	}
}

// timeoutError is a net.Error whose Timeout reports true.
type timeoutError struct{}

func (timeoutError) Error() string   { return "timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

var _ net.Error = timeoutError{}

func TestDefaultIsRetryable(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"server", errs.New(errs.CodeServer, "500"), true},
		{"rate limited", errs.New(errs.CodeRateLimited, "429"), true},
		{"transport", errs.New(errs.CodeTransport, "dial"), true},
		{"client error", errs.New(errs.CodeAPI, "400"), false},
		{"auth", errs.New(errs.CodeUnauthorized, "401"), false},
		{"network timeout", timeoutError{}, true},
		{"permanent", retry.Permanent(errs.New(errs.CodeServer, "x")), false},
		{"context canceled", context.Canceled, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := retry.DefaultIsRetryable(tc.err); got != tc.want {
				t.Fatalf("DefaultIsRetryable(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

func TestNewWithConfigDefaults(t *testing.T) {
	t.Parallel()

	r := retry.NewWithConfig(retry.Config{})
	if got := r.MaxAttempts(); got != 1 {
		t.Fatalf("MaxAttempts() = %d, want 1 for zero config", got)
	}
}
