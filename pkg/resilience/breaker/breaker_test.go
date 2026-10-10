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

package breaker_test

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go.opentelemetry.io/otel/metric/noop"

	"github.com/shing1211/webullapi4go/pkg/resilience/breaker"
	"github.com/shing1211/webullapi4go/pkg/resilience/clock"
)

type stateTransition struct {
	from breaker.State
	to   breaker.State
}

type breakerContextKey struct{}

func TestStateString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		state breaker.State
		want  string
	}{
		{name: "closed", state: breaker.StateClosed, want: "closed"},
		{name: "open", state: breaker.StateOpen, want: "open"},
		{name: "half open", state: breaker.StateHalfOpen, want: "half-open"},
		{name: "unknown", state: breaker.State(99), want: "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.state.String(); got != tt.want {
				t.Fatalf("State(%d).String() = %q, want %q", tt.state, got, tt.want)
			}
		})
	}
}

func TestBreakerLifecycle(t *testing.T) {
	t.Parallel()

	probeErr := errors.New("probe failed")
	tests := []struct {
		name        string
		outcome     error
		wantState   breaker.State
		wantFailure int
	}{
		{name: "success closes", wantState: breaker.StateClosed},
		{name: "failure reopens", outcome: probeErr, wantState: breaker.StateOpen, wantFailure: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fake := clock.NewFake(time.Unix(1_700_000_000, 0))
			var transitions []stateTransition
			b := breaker.New(
				nil,
				breaker.WithClock(fake),
				breaker.WithThreshold(1),
				breaker.WithCooldown(time.Minute),
				breaker.WithHalfOpenMax(1),
				breaker.WithMeter(noop.NewMeterProvider().Meter("breaker-test")),
				breaker.WithOnStateChange(func(from, to breaker.State) {
					transitions = append(transitions, stateTransition{from: from, to: to})
				}),
			)

			if !b.Allow() || b.State() != breaker.StateClosed {
				t.Fatalf("initial state = %v, allow = %v; want closed and allowed", b.State(), b.Allow())
			}
			b.RecordFailure()
			if got := b.State(); got != breaker.StateOpen {
				t.Fatalf("state after threshold = %v, want open", got)
			}
			if b.Allow() {
				t.Fatal("Allow() = true before cooldown")
			}

			called := false
			if err := b.Do(context.Background(), func(context.Context) error {
				called = true
				return nil
			}); !errors.Is(err, breaker.ErrOpen) {
				t.Fatalf("Do() while open error = %v, want ErrOpen", err)
			}
			if called {
				t.Fatal("Do() invoked the function while open")
			}

			fake.Advance(time.Minute)
			ctx := context.WithValue(context.Background(), breakerContextKey{}, "probe")
			err := b.Do(ctx, func(got context.Context) error {
				if got != ctx {
					t.Fatal("Do() did not forward the caller's context")
				}
				return tt.outcome
			})
			if tt.outcome == nil && err != nil {
				t.Fatalf("successful probe returned %v", err)
			}
			if tt.outcome != nil && !errors.Is(err, tt.outcome) {
				t.Fatalf("probe error = %v, want %v", err, tt.outcome)
			}
			if got := b.State(); got != tt.wantState {
				t.Fatalf("final state = %v, want %v", got, tt.wantState)
			}
			if got := b.Failures(); got != tt.wantFailure {
				t.Fatalf("Failures() = %d, want %d", got, tt.wantFailure)
			}

			want := []stateTransition{
				{from: breaker.StateClosed, to: breaker.StateOpen},
				{from: breaker.StateOpen, to: breaker.StateHalfOpen},
			}
			if tt.wantState == breaker.StateClosed {
				want = append(want, stateTransition{from: breaker.StateHalfOpen, to: breaker.StateClosed})
			} else {
				want = append(want, stateTransition{from: breaker.StateHalfOpen, to: breaker.StateOpen})
			}
			if len(transitions) != len(want) {
				t.Fatalf("transitions = %v, want %v", transitions, want)
			}
			for i := range want {
				if transitions[i] != want[i] {
					t.Fatalf("transitions[%d] = %+v, want %+v", i, transitions[i], want[i])
				}
			}
		})
	}
}

func TestNewWithConfigNormalizesInvalidValues(t *testing.T) {
	t.Parallel()

	fake := clock.NewFake(time.Unix(0, 0))
	b := breaker.NewWithConfig(breaker.Config{Clock: fake})
	for i := 1; i < breaker.DefaultThreshold; i++ {
		b.RecordFailure()
		if got := b.State(); got != breaker.StateClosed {
			t.Fatalf("state after failure %d = %v, want closed", i, got)
		}
	}
	b.RecordFailure()
	if got := b.State(); got != breaker.StateOpen {
		t.Fatalf("state after default threshold = %v, want open", got)
	}

	fake.Advance(breaker.DefaultCooldown - time.Nanosecond)
	if b.Allow() {
		t.Fatal("Allow() = true before the default cooldown")
	}
	fake.Advance(time.Nanosecond)
	if !b.Allow() || b.State() != breaker.StateHalfOpen {
		t.Fatalf("state at default cooldown = %v, allow = %v; want half-open and allowed", b.State(), b.Allow())
	}
	if b.Allow() {
		t.Fatal("default half-open limit admitted a second probe")
	}
	b.RecordFailure()
	if got := b.State(); got != breaker.StateOpen {
		t.Fatalf("failed default probe state = %v, want open", got)
	}
}

func TestClosedSuccessReducesFailuresWithoutUnderflow(t *testing.T) {
	t.Parallel()

	b := breaker.New(breaker.WithThreshold(3))
	b.RecordFailure()
	b.RecordFailure()
	b.RecordSuccess()
	if got := b.Failures(); got != 1 {
		t.Fatalf("Failures() after success = %d, want 1", got)
	}
	b.RecordSuccess()
	b.RecordSuccess()
	if got := b.Failures(); got != 0 {
		t.Fatalf("Failures() after additional successes = %d, want 0", got)
	}
}

func TestResetClosesAndClearsBreaker(t *testing.T) {
	t.Parallel()

	b := breaker.New(breaker.WithThreshold(1))
	b.RecordFailure()
	if b.State() != breaker.StateOpen {
		t.Fatalf("state before Reset() = %v, want open", b.State())
	}
	b.Reset()
	if b.State() != breaker.StateClosed || b.Failures() != 0 || !b.Allow() {
		t.Fatalf("after Reset(): state=%v failures=%d allow=%v; want closed, 0, true", b.State(), b.Failures(), b.Allow())
	}
}

func TestHalfOpenAdmissionLimit(t *testing.T) {
	t.Parallel()

	for _, halfOpenMax := range []int{1, 3} {
		t.Run(strconv.Itoa(halfOpenMax), func(t *testing.T) {
			t.Parallel()

			fake := clock.NewFake(time.Unix(0, 0))
			b := breaker.New(
				breaker.WithClock(fake),
				breaker.WithThreshold(1),
				breaker.WithCooldown(time.Second),
				breaker.WithHalfOpenMax(halfOpenMax),
			)
			b.RecordFailure()
			fake.Advance(time.Second)

			var allowed atomic.Int32
			var wg sync.WaitGroup
			for range 32 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if b.Allow() {
						allowed.Add(1)
					}
				}()
			}
			wg.Wait()
			if got := int(allowed.Load()); got != halfOpenMax {
				t.Fatalf("admitted probes = %d, want %d", got, halfOpenMax)
			}
			for range halfOpenMax {
				b.RecordFailure()
			}
			if got := b.State(); got != breaker.StateOpen {
				t.Fatalf("state after failed probes = %v, want open", got)
			}
		})
	}
}
