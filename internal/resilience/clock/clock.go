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

// Package clock provides the time source used by the SDK's resilience
// primitives so that retry, rate limiting, and circuit-breaking behaviour can
// be tested with an injected, deterministic clock.
package clock

import (
	"sync"
	"time"
)

// Clock is the minimal time source the resilience primitives depend on.
type Clock interface {
	// Now returns the current time.
	Now() time.Time
	// After returns a channel that receives the current time once at least d
	// has elapsed.
	After(d time.Duration) <-chan time.Time
}

// system is the real-time [Clock] backed by the time package.
type system struct{}

// Now returns the current wall-clock time.
func (system) Now() time.Time { return time.Now() }

// After delegates to [time.After].
func (system) After(d time.Duration) <-chan time.Time { return time.After(d) }

// System returns the real-time [Clock].
func System() Clock { return system{} }

// Fake is a manually advanced [Clock] intended for tests. It is safe for
// concurrent use.
type Fake struct {
	mu      sync.Mutex
	now     time.Time
	waiters []fakeWaiter
}

// fakeWaiter records a pending After deadline and its delivery channel.
type fakeWaiter struct {
	at time.Time
	ch chan time.Time
}

// NewFake returns a [Fake] clock positioned at start.
func NewFake(start time.Time) *Fake { return &Fake{now: start} }

// Now returns the fake clock's current time.
func (f *Fake) Now() time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.now
}

// After registers a waiter that fires when the clock is advanced past
// Now()+d. The returned channel is buffered so delivery never blocks.
func (f *Fake) After(d time.Duration) <-chan time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	ch := make(chan time.Time, 1)
	f.waiters = append(f.waiters, fakeWaiter{at: f.now.Add(d), ch: ch})
	return ch
}

// Advance moves the clock forward by d and fires every waiter whose deadline
// has been reached, delivering the new current time.
func (f *Fake) Advance(d time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.now = f.now.Add(d)
	now := f.now
	remaining := f.waiters[:0]
	for _, w := range f.waiters {
		if w.at.After(now) {
			remaining = append(remaining, w)
			continue
		}
		w.ch <- now
	}
	f.waiters = remaining
}
