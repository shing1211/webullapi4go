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

package clock_test

import (
	"sync"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/pkg/resilience/clock"
)

var (
	_ clock.Clock = clock.System()
	_ clock.Clock = (*clock.Fake)(nil)
)

func TestSystemClockNowIsWithinCallWindow(t *testing.T) {
	t.Parallel()

	before := time.Now()
	got := clock.System().Now()
	after := time.Now()
	if got.Before(before) || got.After(after) {
		t.Fatalf("System().Now() = %v, want between %v and %v", got, before, after)
	}
}

func TestSystemClockAfter(t *testing.T) {
	t.Parallel()

	select {
	case <-clock.System().After(0):
	case <-time.After(5 * time.Second):
		t.Fatal("System().After(0) did not fire")
	}
}

func TestFakeAdvanceTable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		wait          time.Duration
		beforeAdvance time.Duration
		finalAdvance  time.Duration
	}{
		{name: "negative duration", wait: -time.Second, beforeAdvance: 0, finalAdvance: 0},
		{name: "zero duration", wait: 0, beforeAdvance: 0, finalAdvance: 0},
		{name: "future duration", wait: time.Second, beforeAdvance: 999 * time.Millisecond, finalAdvance: time.Millisecond},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			start := time.Unix(1_700_000_000, 0)
			fake := clock.NewFake(start)
			waiter := fake.After(tt.wait)
			fake.Advance(tt.beforeAdvance)
			if tt.finalAdvance != 0 {
				select {
				case <-waiter:
					t.Fatal("waiter fired before its deadline")
				default:
				}
			}

			fake.Advance(tt.finalAdvance)
			select {
			case got := <-waiter:
				if want := fake.Now(); !got.Equal(want) {
					t.Fatalf("waiter time = %v, want current fake time %v", got, want)
				}
			default:
				t.Fatal("waiter did not fire at its deadline")
			}

			fake.Advance(time.Hour)
			select {
			case got := <-waiter:
				t.Fatalf("waiter fired more than once with %v", got)
			default:
			}
		})
	}
}

func TestFakeConcurrentWaiters(t *testing.T) {
	t.Parallel()

	const waiters = 32
	start := time.Unix(1_700_000_000, 0)
	fake := clock.NewFake(start)
	ready := make(chan struct{})
	results := make(chan time.Time, waiters)
	var wg sync.WaitGroup
	for range waiters {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = fake.Now()
			wait := fake.After(time.Second)
			ready <- struct{}{}
			results <- <-wait
		}()
	}
	for range waiters {
		<-ready
	}
	fake.Advance(time.Second)
	wg.Wait()
	close(results)

	want := start.Add(time.Second)
	for got := range results {
		if !got.Equal(want) {
			t.Fatalf("concurrent waiter time = %v, want %v", got, want)
		}
	}
	if got := fake.Now(); !got.Equal(want) {
		t.Fatalf("Now() = %v, want %v", got, want)
	}
}
