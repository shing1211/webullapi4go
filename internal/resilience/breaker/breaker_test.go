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
	"sync"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/internal/resilience/breaker"
	"github.com/shing1211/webullapi4go/internal/resilience/clock"
)

func TestOpensAfterThreshold(t *testing.T) {
	t.Parallel()

	clk := clock.NewFake(time.Unix(0, 0))
	b := breaker.New(
		breaker.WithClock(clk),
		breaker.WithThreshold(2),
		breaker.WithCooldown(time.Second),
	)

	if !b.Allow() {
		t.Fatal("initial call should be allowed")
	}
	b.RecordFailure()
	if b.State() != breaker.StateClosed {
		t.Fatalf("State() = %v after one failure, want closed", b.State())
	}
	if !b.Allow() {
		t.Fatal("second call should be allowed")
	}
	b.RecordFailure()
	if b.State() != breaker.StateOpen {
		t.Fatalf("State() = %v after threshold, want open", b.State())
	}
	if b.Allow() {
		t.Fatal("call while open should be rejected")
	}
}

func TestHalfOpenClosesOnSuccess(t *testing.T) {
	t.Parallel()

	clk := clock.NewFake(time.Unix(0, 0))
	b := breaker.New(
		breaker.WithClock(clk),
		breaker.WithThreshold(1),
		breaker.WithCooldown(time.Second),
	)
	b.Allow()
	b.RecordFailure() // opens
	if b.Allow() {
		t.Fatal("call before cooldown should be rejected")
	}

	clk.Advance(time.Second)
	if !b.Allow() {
		t.Fatal("probe after cooldown should be allowed")
	}
	if b.State() != breaker.StateHalfOpen {
		t.Fatalf("State() = %v, want half-open", b.State())
	}
	b.RecordSuccess()
	if b.State() != breaker.StateClosed {
		t.Fatalf("State() = %v after probe success, want closed", b.State())
	}
}

func TestHalfOpenReopensOnFailure(t *testing.T) {
	t.Parallel()

	clk := clock.NewFake(time.Unix(0, 0))
	b := breaker.New(
		breaker.WithClock(clk),
		breaker.WithThreshold(1),
		breaker.WithCooldown(time.Second),
	)
	b.Allow()
	b.RecordFailure()

	clk.Advance(time.Second)
	if !b.Allow() {
		t.Fatal("probe after cooldown should be allowed")
	}
	b.RecordFailure()
	if b.State() != breaker.StateOpen {
		t.Fatalf("State() = %v after probe failure, want open", b.State())
	}
	if b.Allow() {
		t.Fatal("call after failed probe should be rejected")
	}
}

func TestHalfOpenLimitsConcurrentProbes(t *testing.T) {
	t.Parallel()

	clk := clock.NewFake(time.Unix(0, 0))
	b := breaker.New(
		breaker.WithClock(clk),
		breaker.WithThreshold(1),
		breaker.WithCooldown(time.Second),
		breaker.WithHalfOpenMax(1),
	)
	b.Allow()
	b.RecordFailure()
	clk.Advance(time.Second)

	if !b.Allow() {
		t.Fatal("first probe should be allowed")
	}
	if b.Allow() {
		t.Fatal("second concurrent probe should be rejected")
	}
	b.RecordSuccess()
	if !b.Allow() {
		t.Fatal("call after closing should be allowed")
	}
}

func TestStateChangeCallback(t *testing.T) {
	t.Parallel()

	clk := clock.NewFake(time.Unix(0, 0))
	var mu sync.Mutex
	var transitions [][2]breaker.State
	b := breaker.New(
		breaker.WithClock(clk),
		breaker.WithThreshold(1),
		breaker.WithCooldown(time.Second),
		breaker.WithOnStateChange(func(from, to breaker.State) {
			mu.Lock()
			defer mu.Unlock()
			transitions = append(transitions, [2]breaker.State{from, to})
		}),
	)
	b.Allow()
	b.RecordFailure()
	clk.Advance(time.Second)
	b.Allow() // half-open
	b.RecordSuccess()

	mu.Lock()
	defer mu.Unlock()
	want := [][2]breaker.State{
		{breaker.StateClosed, breaker.StateOpen},
		{breaker.StateOpen, breaker.StateHalfOpen},
		{breaker.StateHalfOpen, breaker.StateClosed},
	}
	if len(transitions) != len(want) {
		t.Fatalf("transitions = %v, want %v", transitions, want)
	}
	for i := range want {
		if transitions[i] != want[i] {
			t.Fatalf("transitions[%d] = %v, want %v", i, transitions[i], want[i])
		}
	}
}

func TestDoReturnsErrOpen(t *testing.T) {
	t.Parallel()

	clk := clock.NewFake(time.Unix(0, 0))
	b := breaker.New(
		breaker.WithClock(clk),
		breaker.WithThreshold(1),
		breaker.WithCooldown(time.Second),
	)
	if err := b.Do(context.Background(), func(context.Context) error { return nil }); err != nil {
		t.Fatalf("Do() error = %v, want nil", err)
	}
	_ = b.Do(context.Background(), func(context.Context) error { return errors.New("boom") })

	called := false
	err := b.Do(context.Background(), func(context.Context) error {
		called = true
		return nil
	})
	if !errors.Is(err, breaker.ErrOpen) {
		t.Fatalf("Do() error = %v, want ErrOpen", err)
	}
	if called {
		t.Fatal("fn should not run while the circuit is open")
	}
}

func TestReset(t *testing.T) {
	t.Parallel()

	b := breaker.New(breaker.WithThreshold(1))
	b.Allow()
	b.RecordFailure()
	b.Reset()
	if b.State() != breaker.StateClosed || b.Failures() != 0 {
		t.Fatalf("after Reset: state = %v, failures = %d; want closed, 0", b.State(), b.Failures())
	}
}

func TestConcurrentAccess(t *testing.T) {
	t.Parallel()

	b := breaker.New(breaker.WithThreshold(3), breaker.WithCooldown(time.Millisecond))
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if b.Allow() {
				if i%2 == 0 {
					b.RecordSuccess()
				} else {
					b.RecordFailure()
				}
			}
			_ = b.State()
		}(i)
	}
	wg.Wait()
}
