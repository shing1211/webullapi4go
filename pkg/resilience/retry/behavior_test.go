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
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	errs "github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/pkg/resilience/clock"
	"github.com/shing1211/webullapi4go/pkg/resilience/retry"
)

type retryRecordingClock struct {
	*clock.Fake
	waits chan time.Duration
}

func newRetryRecordingClock(start time.Time) *retryRecordingClock {
	return &retryRecordingClock{
		Fake:  clock.NewFake(start),
		waits: make(chan time.Duration, 8),
	}
}

func (c *retryRecordingClock) After(d time.Duration) <-chan time.Time {
	ch := c.Fake.After(d)
	c.waits <- d
	return ch
}

type immediateRetryClock struct {
	now   time.Time
	calls atomic.Int32
}

func (c *immediateRetryClock) Now() time.Time { return c.now }

func (c *immediateRetryClock) After(time.Duration) <-chan time.Time {
	c.calls.Add(1)
	ch := make(chan time.Time, 1)
	ch <- c.now
	return ch
}

type retryTimeoutError struct {
	timeout bool
}

func (e retryTimeoutError) Error() string   { return "network error" }
func (e retryTimeoutError) Timeout() bool   { return e.timeout }
func (e retryTimeoutError) Temporary() bool { return e.timeout }

func TestDoOutcomesAndDelaySchedule(t *testing.T) {
	t.Parallel()

	serverErr := errs.New(errs.CodeServer, "server error")
	apiErr := errs.New(errs.CodeAPI, "bad request")
	tests := []struct {
		name        string
		maxAttempts int
		maxDelay    time.Duration
		failures    int
		failure     error
		classifier  func(error) bool
		wantCalls   int
		wantDelays  []time.Duration
		wantError   bool
	}{
		{name: "success after retries", maxAttempts: 3, maxDelay: 20 * time.Millisecond, failures: 2, wantCalls: 3, wantDelays: []time.Duration{10 * time.Millisecond, 20 * time.Millisecond}},
		{name: "attempts exhausted", maxAttempts: 2, maxDelay: 20 * time.Millisecond, failures: 2, wantCalls: 2, wantDelays: []time.Duration{10 * time.Millisecond}, wantError: true},
		{name: "typed non-retryable", maxAttempts: 5, maxDelay: 20 * time.Millisecond, failures: 1, failure: apiErr, wantCalls: 1, wantError: true},
		{name: "custom non-retryable", maxAttempts: 5, maxDelay: 20 * time.Millisecond, failures: 1, failure: errors.New("stop"), classifier: func(error) bool { return false }, wantCalls: 1, wantError: true},
		{name: "delay capped", maxAttempts: 4, maxDelay: 15 * time.Millisecond, failures: 3, wantCalls: 4, wantDelays: []time.Duration{10 * time.Millisecond, 15 * time.Millisecond, 15 * time.Millisecond}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fake := newRetryRecordingClock(time.Unix(1_700_000_000, 0))
			opts := []retry.Option{
				retry.WithMaxAttempts(tt.maxAttempts),
				retry.WithBaseDelay(10 * time.Millisecond),
				retry.WithMaxDelay(tt.maxDelay),
				retry.WithJitter(false),
				retry.WithClock(fake),
			}
			if tt.classifier != nil {
				opts = append(opts, retry.WithIsRetryable(tt.classifier))
			}
			retrier := retry.New(opts...)
			if got := retrier.MaxAttempts(); got != tt.maxAttempts {
				t.Fatalf("MaxAttempts() = %d, want %d", got, tt.maxAttempts)
			}
			calls := 0
			result := make(chan error, 1)
			go func() {
				result <- retrier.Do(context.Background(), func(context.Context) error {
					calls++
					if calls <= tt.failures {
						if tt.failure != nil {
							return tt.failure
						}
						return serverErr
					}
					return nil
				})
			}()

			for i, want := range tt.wantDelays {
				select {
				case got := <-fake.waits:
					if got != want {
						t.Fatalf("wait %d = %v, want %v", i, got, want)
					}
					fake.Advance(got)
				case <-time.After(5 * time.Second):
					t.Fatalf("timed out waiting for delay %d", i)
				}
			}
			select {
			case err := <-result:
				if (err != nil) != tt.wantError {
					t.Fatalf("Do() error = %v, wantError = %v", err, tt.wantError)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("Do did not return")
			}
			if calls != tt.wantCalls {
				t.Fatalf("calls = %d, want %d", calls, tt.wantCalls)
			}
			if got := len(fake.waits); got != 0 {
				t.Fatalf("unconsumed waits = %d, want 0", got)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := retry.DefaultConfig()
	if cfg.MaxAttempts != retry.DefaultMaxAttempts || cfg.BaseDelay != retry.DefaultBaseDelay || cfg.MaxDelay != retry.DefaultMaxDelay {
		t.Fatalf("DefaultConfig() timing = %+v", cfg)
	}
	if !cfg.Jitter || cfg.FullJitter || cfg.IsRetryable == nil || cfg.Clock == nil {
		t.Fatalf("DefaultConfig() behavior = %+v", cfg)
	}
}

func TestNewWithConfigNormalizesInvalidTiming(t *testing.T) {
	t.Parallel()

	fake := &immediateRetryClock{now: time.Unix(0, 0)}
	retrier := retry.NewWithConfig(retry.Config{
		MaxAttempts: 2,
		BaseDelay:   -time.Second,
		MaxDelay:    -time.Second,
		Jitter:      false,
		IsRetryable: func(error) bool { return true },
		Clock:       fake,
	})
	calls := 0
	err := retrier.Do(context.Background(), func(context.Context) error {
		calls++
		return errors.New("retry")
	})
	if err == nil {
		t.Fatal("Do() error = nil, want final retry error")
	}
	if got := retrier.MaxAttempts(); got != 2 {
		t.Fatalf("MaxAttempts() = %d, want 2", got)
	}
	if got := retry.NewWithConfig(retry.Config{}).MaxAttempts(); got != 1 {
		t.Fatalf("zero Config MaxAttempts() = %d, want 1", got)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
	if got := fake.calls.Load(); got != 0 {
		t.Fatalf("After() calls = %d for normalized zero delay, want 0", got)
	}
}

func TestPermanentPreservesCauseAndWrapping(t *testing.T) {
	t.Parallel()

	if retry.Permanent(nil) != nil {
		t.Fatal("Permanent(nil) != nil")
	}
	cause := errors.New("cause")
	permanent := retry.Permanent(cause)
	if !errors.Is(permanent, cause) || !retry.IsPermanent(permanent) {
		t.Fatal("Permanent() did not preserve its cause")
	}
	wrapped := fmt.Errorf("outer: %w", permanent)
	if !retry.IsPermanent(wrapped) {
		t.Fatal("IsPermanent() did not traverse a wrapped permanent error")
	}
	if retry.IsPermanent(cause) || retry.IsPermanent(nil) {
		t.Fatal("ordinary error classified as permanent")
	}
}

func TestDefaultIsRetryable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "rate limited code", err: errs.New(errs.CodeRateLimited, "limited"), want: true},
		{name: "server code", err: errs.New(errs.CodeServer, "server"), want: true},
		{name: "transport code", err: errs.New(errs.CodeTransport, "transport"), want: true},
		{name: "HTTP 429 status", err: &errs.Error{Code: errs.CodeAPI, Status: http.StatusTooManyRequests}, want: true},
		{name: "HTTP 503 status", err: &errs.Error{Code: errs.CodeAPI, Status: http.StatusServiceUnavailable}, want: true},
		{name: "HTTP 400 status", err: &errs.Error{Code: errs.CodeAPI, Status: http.StatusBadRequest}, want: false},
		{name: "wrapped typed error", err: fmt.Errorf("wrapped: %w", errs.New(errs.CodeServer, "server")), want: true},
		{name: "context canceled", err: context.Canceled, want: false},
		{name: "network timeout", err: retryTimeoutError{timeout: true}, want: true},
		{name: "non-timeout network error", err: retryTimeoutError{}, want: false},
		{name: "ordinary error", err: errors.New("ordinary"), want: false},
		{name: "permanent typed error", err: retry.Permanent(errs.New(errs.CodeServer, "server")), want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := retry.DefaultIsRetryable(tt.err); got != tt.want {
				t.Fatalf("DefaultIsRetryable(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
