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

type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

type system struct{}

func (system) Now() time.Time { return time.Now() }

func (system) After(d time.Duration) <-chan time.Time { return time.After(d) }

func System() Clock { return system{} }

type Fake struct {
	mu      sync.Mutex
	now     time.Time
	waiters []fakeWaiter
}

type fakeWaiter struct {
	at time.Time
	ch chan time.Time
}

func NewFake(start time.Time) *Fake { return &Fake{now: start} }

func (f *Fake) Now() time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.now
}

func (f *Fake) After(d time.Duration) <-chan time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	ch := make(chan time.Time, 1)
	f.waiters = append(f.waiters, fakeWaiter{at: f.now.Add(d), ch: ch})
	return ch
}

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
