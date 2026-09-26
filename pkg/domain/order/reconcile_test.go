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

package order

import (
	"errors"
	"sync"
	"testing"
)

func TestMachineStateNilReceiver(t *testing.T) {
	t.Parallel()

	var m *Machine
	if got := m.State(); got != StateUnknown {
		t.Fatalf("nil Machine.State() = %q, want %q", got, StateUnknown)
	}
}

func TestSetOnStateChangedNilAndNilReceiver(t *testing.T) {
	t.Parallel()

	var m *Machine
	m.SetOnStateChanged(func(State, State) {}) // must not panic

	mach := New(StatePending)
	mach.SetOnStateChanged(func(State, State) { t.Error("callback fired after being disabled") })
	mach.SetOnStateChanged(nil)
	if _, err := mach.ApplyEvent(EventPlace); err != nil {
		t.Fatalf("ApplyEvent() error = %v", err)
	}
}

func TestApplyEventNilReceiver(t *testing.T) {
	t.Parallel()

	var m *Machine
	if got, err := m.ApplyEvent(EventPlace); got != StateUnknown || !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("nil ApplyEvent() = (%q, %v), want (%q, ErrInvalidOrder)", got, err, StateUnknown)
	}
}

func TestReconcileStatusNilReceiver(t *testing.T) {
	t.Parallel()

	var m *Machine
	if got, err := m.ReconcileStatus(StateFilled); got != StateUnknown || !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("nil ReconcileStatus() = (%q, %v), want (%q, ErrInvalidOrder)", got, err, StateUnknown)
	}
}

// TestReconcileStatusUnknownStatus covers the unmappable-status contract: the
// current state is preserved and ErrInvalidStatus is returned.
func TestReconcileStatusUnknownStatus(t *testing.T) {
	t.Parallel()

	m := New(StatePending)
	got, err := m.ReconcileStatus(StateUnknown)
	if !errors.Is(err, ErrInvalidStatus) {
		t.Fatalf("ReconcileStatus(StateUnknown) error = %v, want ErrInvalidStatus", err)
	}
	if got != StatePending {
		t.Fatalf("ReconcileStatus(StateUnknown) = %q, want unchanged %q", got, StatePending)
	}
}

// TestReconcileStatusStaleAndTerminalNoOps covers the three ways reconciliation
// deliberately does not move the machine.
func TestReconcileStatusStaleAndTerminalNoOps(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		initial State
		target  State
	}{
		{"same state", StateConfirmed, StateConfirmed},
		{"stale non-terminal", StateConfirmed, StateSubmitted},
		{"stale behind partial fill", StatePartialFilled, StateConfirmed},
		{"terminal never regressed to pending", StateFilled, StatePending},
		{"terminal never regressed to confirmed", StateCancelled, StateConfirmed},
		{"expired is terminal", StateExpired, StateSubmitted},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := New(tc.initial)
			m.SetOnStateChanged(func(from, to State) {
				t.Errorf("OnStateChanged(%q, %q) fired for a no-op reconcile", from, to)
			})

			got, err := m.ReconcileStatus(tc.target)
			if err != nil {
				t.Fatalf("ReconcileStatus(%q) error = %v, want nil", tc.target, err)
			}
			if got != tc.initial {
				t.Fatalf("ReconcileStatus(%q) = %q, want unchanged %q", tc.target, got, tc.initial)
			}
			if m.State() != tc.initial {
				t.Fatalf("machine state = %q, want unchanged %q", m.State(), tc.initial)
			}
		})
	}
}

// TestReconcileStatusAdvancesAlongEventPath covers every target that has a
// normal event path, including the multi-step paths.
func TestReconcileStatusAdvancesAlongEventPath(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		initial State
		target  State
		want    State
	}{
		{"unknown to pending is a no-op", StateUnknown, StatePending, StateUnknown},
		{"pending to submitted", StatePending, StateSubmitted, StateSubmitted},
		{"unknown to submitted", StateUnknown, StateSubmitted, StateSubmitted},
		{"pending to pre-submitted", StatePending, StatePreSubmitted, StatePreSubmitted},
		{"submitted to confirmed", StateSubmitted, StateConfirmed, StateConfirmed},
		{"pre-submitted to confirmed", StatePreSubmitted, StateConfirmed, StateConfirmed},
		{"unknown to confirmed", StateUnknown, StateConfirmed, StateConfirmed},
		{"pending to partial filled", StatePending, StatePartialFilled, StatePartialFilled},
		{"confirmed to partial filled", StateConfirmed, StatePartialFilled, StatePartialFilled},
		{"partial filled to filled", StatePartialFilled, StateFilled, StateFilled},
		{"unknown to filled", StateUnknown, StateFilled, StateFilled},
		{"submitted to cancelled", StateSubmitted, StateCancelled, StateCancelled},
		{"confirmed to cancelled", StateConfirmed, StateCancelled, StateCancelled},
		{"submitted to failed", StateSubmitted, StateFailed, StateFailed},
		{"unknown to failed", StateUnknown, StateFailed, StateFailed},
		{"confirmed to expired", StateConfirmed, StateExpired, StateExpired},
		{"submitted to expired", StateSubmitted, StateExpired, StateExpired},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := New(tc.initial)
			var changes int
			var mu sync.Mutex
			m.SetOnStateChanged(func(from, to State) {
				mu.Lock()
				changes++
				mu.Unlock()
				if from == to {
					t.Errorf("OnStateChanged(%q, %q) reported a self-transition", from, to)
				}
			})

			got, err := m.ReconcileStatus(tc.target)
			if err != nil {
				t.Fatalf("ReconcileStatus(%q) error = %v", tc.target, err)
			}
			if got != tc.want {
				t.Fatalf("ReconcileStatus(%q) = %q, want %q", tc.target, got, tc.want)
			}
			if m.State() != tc.want {
				t.Fatalf("machine state = %q, want %q", m.State(), tc.want)
			}
			mu.Lock()
			defer mu.Unlock()
			wantChanges := 0
			if tc.want != tc.initial {
				wantChanges = 1
			}
			if changes != wantChanges {
				t.Fatalf("OnStateChanged fired %d times, want %d", changes, wantChanges)
			}
		})
	}
}

// TestReconcileStatusAuthoritativeTerminalSnapshot covers a terminal target that
// has no legal event path. The machine must record the snapshot rather than
// fail, because the exchange is authoritative about a terminal state.
func TestReconcileStatusAuthoritativeTerminalSnapshot(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		initial State
		target  State
	}{
		{"confirmed to failed has no event path", StateConfirmed, StateFailed},
		{"partial filled to failed has no event path", StatePartialFilled, StateFailed},
		{"partial filled to rejected snapshot", StatePartialFilled, StateFailed},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := New(tc.initial)
			var gotFrom, gotTo State
			var fired bool
			m.SetOnStateChanged(func(from, to State) {
				gotFrom, gotTo, fired = from, to, true
			})

			got, err := m.ReconcileStatus(tc.target)
			if err != nil {
				t.Fatalf("ReconcileStatus(%q) error = %v, want nil for an authoritative terminal snapshot", tc.target, err)
			}
			if got != tc.target {
				t.Fatalf("ReconcileStatus(%q) = %q, want %q", tc.target, got, tc.target)
			}
			if !fired {
				t.Fatal("OnStateChanged did not fire for an authoritative terminal snapshot")
			}
			if gotFrom != tc.initial || gotTo != tc.target {
				t.Fatalf("OnStateChanged(%q, %q), want (%q, %q)", gotFrom, gotTo, tc.initial, tc.target)
			}
		})
	}
}

// TestReconcileStatusInvalidRegression covers a non-terminal target that has no
// event path and is not authoritative. That is a caller error.
func TestReconcileStatusInvalidRegression(t *testing.T) {
	t.Parallel()

	m := New(StateSubmitted)
	got, err := m.ReconcileStatus(StatePreSubmitted)
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("ReconcileStatus(PreSubmitted) error = %v, want ErrInvalidTransition", err)
	}
	if got != StateSubmitted {
		t.Fatalf("ReconcileStatus(PreSubmitted) = %q, want unchanged %q", got, StateSubmitted)
	}
	if m.State() != StateSubmitted {
		t.Fatalf("machine state = %q, want unchanged %q", m.State(), StateSubmitted)
	}
}

func TestStatusEventsPaths(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		current    State
		target     State
		wantEvents []Event
		wantOK     bool
	}{
		{"unknown to pending is a no-op path", StateUnknown, StatePending, nil, true},
		{"stale confirmed to submitted is a no-op path", StateConfirmed, StateSubmitted, nil, true},
		{"stale partial to pre-submitted is a no-op path", StatePartialFilled, StatePreSubmitted, nil, true},
		{"pending to submitted", StatePending, StateSubmitted, []Event{EventPlace}, true},
		{"pending to pre-submitted", StatePending, StatePreSubmitted, []Event{EventAcknowledge}, true},
		{"unknown to confirmed", StateUnknown, StateConfirmed, []Event{EventAcknowledge, EventConfirm}, true},
		{"submitted to confirmed", StateSubmitted, StateConfirmed, []Event{EventAcknowledge}, true},
		{"pre-submitted to confirmed", StatePreSubmitted, StateConfirmed, []Event{EventConfirm}, true},
		{"pending to filled", StatePending, StateFilled, []Event{EventPlace, EventFill}, true},
		{"confirmed to filled", StateConfirmed, StateFilled, []Event{EventFill}, true},
		{"partial to filled", StatePartialFilled, StateFilled, []Event{EventFill}, true},
		{"confirmed to cancelled", StateConfirmed, StateCancelled, []Event{EventCancel}, true},
		{"submitted to failed", StateSubmitted, StateFailed, []Event{EventReject}, true},
		{"confirmed to expired", StateConfirmed, StateExpired, []Event{EventExpire}, true},
		{"confirmed to failed has no path", StateConfirmed, StateFailed, nil, false},
		{"partial to failed has no path", StatePartialFilled, StateFailed, nil, false},
		{"partial to rejected has no path", StatePartialFilled, StateFailed, nil, false},
		{"submitted to pre-submitted has no path", StateSubmitted, StatePreSubmitted, nil, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, ok := statusEvents(tc.current, tc.target)
			if ok != tc.wantOK {
				t.Fatalf("statusEvents(%q, %q) ok = %v, want %v", tc.current, tc.target, ok, tc.wantOK)
			}
			if len(got) != len(tc.wantEvents) {
				t.Fatalf("statusEvents(%q, %q) = %v, want %v", tc.current, tc.target, got, tc.wantEvents)
			}
			for i := range got {
				if got[i] != tc.wantEvents[i] {
					t.Fatalf("statusEvents(%q, %q) = %v, want %v", tc.current, tc.target, got, tc.wantEvents)
				}
			}
		})
	}
}

func TestStateProgressOrdering(t *testing.T) {
	t.Parallel()

	want := map[State]int{
		StateUnknown:       0,
		StatePending:       1,
		StateSubmitted:     2,
		StatePreSubmitted:  3,
		StateConfirmed:     4,
		StatePartialFilled: 5,
		// Terminal states intentionally share the default progress value; the
		// staleness check guards on IsTerminal before comparing progress.
		StateFilled:    0,
		StateCancelled: 0,
		StateFailed:    0,
		StateExpired:   0,
	}

	for state, wantProgress := range want {
		if got := stateProgress(state); got != wantProgress {
			t.Fatalf("stateProgress(%q) = %d, want %d", state, got, wantProgress)
		}
	}
}

// TestTransitionNonMutatingEvents pins the contract that failure events for
// cancel, modify, and reject report an unsuccessful operation without moving
// the order.
func TestTransitionNonMutatingEvents(t *testing.T) {
	t.Parallel()

	active := []State{StateUnknown, StatePending, StateSubmitted, StatePreSubmitted, StateConfirmed, StatePartialFilled}
	nonMutating := []Event{EventCancelRequest, EventCancelFailed, EventRejectFailed, EventModifyFailed}

	for _, s := range active {
		for _, e := range nonMutating {
			got, err := transition(s, e)
			if err != nil {
				t.Fatalf("transition(%q, %q) error = %v, want nil", s, e, err)
			}
			if got != s {
				t.Fatalf("transition(%q, %q) = %q, want unchanged %q", s, e, got, s)
			}
		}
	}
}

func TestTransitionTerminalIsAbsorbing(t *testing.T) {
	t.Parallel()

	terminals := []State{StateFilled, StateCancelled, StateFailed, StateExpired}
	all := []Event{
		EventPlace, EventAcknowledge, EventConfirm, EventFill, EventPartialFill,
		EventCancelRequest, EventCancel, EventCancelFailed, EventReject,
		EventRejectFailed, EventExpire, EventReplace, EventModifyFailed,
	}

	for _, s := range terminals {
		for _, e := range all {
			got, err := transition(s, e)
			if err != nil {
				t.Fatalf("transition(%q, %q) error = %v, want nil for a terminal state", s, e, err)
			}
			if got != s {
				t.Fatalf("transition(%q, %q) = %q, want terminal state %q preserved", s, e, got, s)
			}
		}
	}
}

func TestTransitionFromActiveStates(t *testing.T) {
	t.Parallel()

	cases := []struct {
		from  State
		event Event
		want  State
	}{
		{StateUnknown, EventPlace, StateSubmitted},
		{StateUnknown, EventAcknowledge, StatePreSubmitted},
		{StateUnknown, EventReplace, StateSubmitted},
		{StateUnknown, EventCancel, StateCancelled},
		{StateUnknown, EventReject, StateFailed},
		{StateUnknown, EventExpire, StateExpired},
		{StatePreSubmitted, EventConfirm, StateConfirmed},
		{StatePreSubmitted, EventAcknowledge, StatePreSubmitted},
		{StatePreSubmitted, EventReplace, StateSubmitted},
		{StatePreSubmitted, EventReject, StateFailed},
		{StateSubmitted, EventAcknowledge, StateConfirmed},
		{StateSubmitted, EventFill, StateFilled},
		{StateSubmitted, EventReplace, StateSubmitted},
		{StateSubmitted, EventReject, StateFailed},
		{StateConfirmed, EventFill, StateFilled},
		{StateConfirmed, EventPartialFill, StatePartialFilled},
		{StateConfirmed, EventReplace, StateSubmitted},
		{StateConfirmed, EventCancel, StateCancelled},
		{StateConfirmed, EventExpire, StateExpired},
		{StatePartialFilled, EventFill, StateFilled},
		{StatePartialFilled, EventPartialFill, StatePartialFilled},
		{StatePartialFilled, EventCancel, StateCancelled},
		{StatePartialFilled, EventExpire, StateExpired},
	}

	for _, tc := range cases {
		got, err := transition(tc.from, tc.event)
		if err != nil {
			t.Fatalf("transition(%q, %q) error = %v", tc.from, tc.event, err)
		}
		if got != tc.want {
			t.Fatalf("transition(%q, %q) = %q, want %q", tc.from, tc.event, got, tc.want)
		}
	}
}

func TestTransitionRejectsIllegalEvents(t *testing.T) {
	t.Parallel()

	cases := []struct {
		from  State
		event Event
	}{
		{StateSubmitted, EventConfirm},
		{StateConfirmed, EventAcknowledge},
		{StateConfirmed, EventReject},
		{StatePartialFilled, EventAcknowledge},
		{StatePartialFilled, EventReplace},
		{StatePartialFilled, EventReject},
		{StateUnknown, Event("NOT_A_REAL_EVENT")},
		{StateSubmitted, Event("NOT_A_REAL_EVENT")},
	}

	for _, tc := range cases {
		got, err := transition(tc.from, tc.event)
		if !errors.Is(err, ErrInvalidTransition) {
			t.Fatalf("transition(%q, %q) error = %v, want ErrInvalidTransition", tc.from, tc.event, err)
		}
		if got != tc.from {
			t.Fatalf("transition(%q, %q) = %q, want unchanged %q", tc.from, tc.event, got, tc.from)
		}
	}
}

func TestOrderNilAndMissingMachine(t *testing.T) {
	t.Parallel()

	var nilOrder *Order
	if got := nilOrder.State(); got != StateUnknown {
		t.Fatalf("nil Order.State() = %q, want %q", got, StateUnknown)
	}
	if got, err := nilOrder.ApplyEvent(EventPlace); got != StateUnknown || !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("nil Order.ApplyEvent() = (%q, %v), want ErrInvalidOrder", got, err)
	}
	if got, err := nilOrder.ReconcileStatus("FILLED"); got != StateUnknown || !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("nil Order.ReconcileStatus() = (%q, %v), want ErrInvalidOrder", got, err)
	}
	if got, err := nilOrder.ReconcileState(StateFilled); got != StateUnknown || !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("nil Order.ReconcileState() = (%q, %v), want ErrInvalidOrder", got, err)
	}

	// A non-nil Order with no machine must behave the same way.
	empty := &Order{}
	if got := empty.State(); got != StateUnknown {
		t.Fatalf("machine-less Order.State() = %q, want %q", got, StateUnknown)
	}
	if _, err := empty.ApplyEvent(EventPlace); !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("machine-less Order.ApplyEvent() error = %v, want ErrInvalidOrder", err)
	}
	if _, err := empty.ReconcileStatus("FILLED"); !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("machine-less Order.ReconcileStatus() error = %v, want ErrInvalidOrder", err)
	}
	if _, err := empty.ReconcileState(StateFilled); !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("machine-less Order.ReconcileState() error = %v, want ErrInvalidOrder", err)
	}
}

func TestOrderStateAndApplyEvent(t *testing.T) {
	t.Parallel()

	o := &Order{
		PlaceOrderResult: PlaceOrderResult{ClientOrderID: "c-1", OrderID: "o-1"},
		AccountID:        "acc-1",
		Machine:          New(StatePending),
	}

	if got := o.State(); got != StatePending {
		t.Fatalf("Order.State() = %q, want %q", got, StatePending)
	}
	got, err := o.ApplyEvent(EventPlace)
	if err != nil {
		t.Fatalf("Order.ApplyEvent() error = %v", err)
	}
	if got != StateSubmitted {
		t.Fatalf("Order.ApplyEvent() = %q, want %q", got, StateSubmitted)
	}
	if o.State() != StateSubmitted {
		t.Fatalf("Order.State() = %q, want %q", o.State(), StateSubmitted)
	}
	if o.OrderID != "o-1" || o.ClientOrderID != "c-1" {
		t.Fatalf("embedded PlaceOrderResult = %+v, want echoed identifiers", o.PlaceOrderResult)
	}
}

func TestOrderReconcileStatusFromWireString(t *testing.T) {
	t.Parallel()

	o := &Order{Machine: New(StatePending)}

	got, err := o.ReconcileStatus("CONFIRMED")
	if err != nil {
		t.Fatalf("Order.ReconcileStatus(CONFIRMED) error = %v", err)
	}
	if got != StateConfirmed {
		t.Fatalf("Order.ReconcileStatus(CONFIRMED) = %q, want %q", got, StateConfirmed)
	}

	// A wire value that cannot be mapped must not move the order.
	got, err = o.ReconcileStatus("NOT_A_STATUS")
	if !errors.Is(err, ErrInvalidStatus) {
		t.Fatalf("Order.ReconcileStatus(NOT_A_STATUS) error = %v, want ErrInvalidStatus", err)
	}
	if got != StateConfirmed {
		t.Fatalf("Order.ReconcileStatus(NOT_A_STATUS) = %q, want unchanged %q", got, StateConfirmed)
	}
}

func TestOrderReconcileStateAcceptsNormalizedState(t *testing.T) {
	t.Parallel()

	o := &Order{Machine: New(StatePending)}

	got, err := o.ReconcileState(StateFilled)
	if err != nil {
		t.Fatalf("Order.ReconcileState(Filled) error = %v", err)
	}
	if got != StateFilled {
		t.Fatalf("Order.ReconcileState(Filled) = %q, want %q", got, StateFilled)
	}

	if _, err := o.ReconcileState(StateUnknown); !errors.Is(err, ErrInvalidStatus) {
		t.Fatalf("Order.ReconcileState(Unknown) error = %v, want ErrInvalidStatus", err)
	}
}

// TestReconcileStatusConcurrent exercises reconciliation under the race detector
// and asserts the machine always lands on a terminal or stable state.
func TestReconcileStatusConcurrent(t *testing.T) {
	t.Parallel()

	m := New(StatePending)
	targets := []State{StateSubmitted, StateConfirmed, StatePartialFilled, StateFilled}

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 32; j++ {
				if _, err := m.ReconcileStatus(targets[(i+j)%len(targets)]); err != nil {
					t.Errorf("ReconcileStatus() error = %v", err)
					return
				}
			}
		}(i)
	}
	wg.Wait()

	final := m.State()
	if !IsTerminal(final) && final != StateSubmitted && final != StateConfirmed && final != StatePartialFilled {
		t.Fatalf("final state = %q, want a progressed or terminal state", final)
	}
}
