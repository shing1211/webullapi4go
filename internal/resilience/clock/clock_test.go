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
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/internal/resilience/clock"
)

func TestFakeAdvanceFiresWaiters(t *testing.T) {
	t.Parallel()

	start := time.Unix(1000, 0)
	f := clock.NewFake(start)

	soon := f.After(50 * time.Millisecond)
	later := f.After(150 * time.Millisecond)

	f.Advance(60 * time.Millisecond)
	select {
	case <-soon:
	default:
		t.Fatal("waiter for 50ms did not fire after advancing 60ms")
	}
	select {
	case <-later:
		t.Fatal("waiter for 150ms fired after advancing only 60ms")
	default:
	}

	if got, want := f.Now(), start.Add(60*time.Millisecond); !got.Equal(want) {
		t.Fatalf("Now() = %v, want %v", got, want)
	}

	f.Advance(100 * time.Millisecond)
	select {
	case <-later:
	default:
		t.Fatal("waiter for 150ms did not fire after advancing 160ms total")
	}
}

func TestFakeZeroDurationFiresOnAdvance(t *testing.T) {
	t.Parallel()

	f := clock.NewFake(time.Unix(0, 0))
	ch := f.After(0)
	f.Advance(time.Nanosecond)
	select {
	case <-ch:
	default:
		t.Fatal("waiter for zero duration did not fire")
	}
}
