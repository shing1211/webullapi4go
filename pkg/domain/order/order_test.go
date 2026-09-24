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
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestMachineApplyEvent(t *testing.T) {
	tests := []struct {
		name    string
		initial State
		event   Event
		want    State
		wantErr bool
	}{
		{name: "unknown place", initial: StateUnknown, event: EventPlace, want: StateSubmitted},
		{name: "pending place", initial: StatePending, event: EventPlace, want: StateSubmitted},
		{name: "unknown acknowledge", initial: StateUnknown, event: EventAcknowledge, want: StatePreSubmitted},
		{name: "pending acknowledge", initial: StatePending, event: EventAcknowledge, want: StatePreSubmitted},
		{name: "submitted acknowledge", initial: StateSubmitted, event: EventAcknowledge, want: StateConfirmed},
		{name: "pre-submitted acknowledge", initial: StatePreSubmitted, event: EventAcknowledge, want: StatePreSubmitted},
		{name: "pre-submitted confirm", initial: StatePreSubmitted, event: EventConfirm, want: StateConfirmed},
		{name: "submitted fill", initial: StateSubmitted, event: EventFill, want: StateFilled},
		{name: "pre-submitted fill", initial: StatePreSubmitted, event: EventFill, want: StateFilled},
		{name: "confirmed fill", initial: StateConfirmed, event: EventFill, want: StateFilled},
		{name: "partial fill", initial: StatePartialFilled, event: EventFill, want: StateFilled},
		{name: "submitted partial fill", initial: StateSubmitted, event: EventPartialFill, want: StatePartialFilled},
		{name: "pre-submitted partial fill", initial: StatePreSubmitted, event: EventPartialFill, want: StatePartialFilled},
		{name: "confirmed partial fill", initial: StateConfirmed, event: EventPartialFill, want: StatePartialFilled},
		{name: "partial fill idempotent", initial: StatePartialFilled, event: EventPartialFill, want: StatePartialFilled},
		{name: "confirmed replace", initial: StateConfirmed, event: EventReplace, want: StateSubmitted},
		{name: "pending cancel", initial: StatePending, event: EventCancel, want: StateCancelled},
		{name: "submitted cancel", initial: StateSubmitted, event: EventCancel, want: StateCancelled},
		{name: "pre-submitted cancel", initial: StatePreSubmitted, event: EventCancel, want: StateCancelled},
		{name: "confirmed cancel", initial: StateConfirmed, event: EventCancel, want: StateCancelled},
		{name: "partial-filled cancel", initial: StatePartialFilled, event: EventCancel, want: StateCancelled},
		{name: "pending reject", initial: StatePending, event: EventReject, want: StateFailed},
		{name: "pre-submitted reject", initial: StatePreSubmitted, event: EventReject, want: StateFailed},
		{name: "submitted reject", initial: StateSubmitted, event: EventReject, want: StateFailed},
		{name: "pending expire", initial: StatePending, event: EventExpire, want: StateExpired},
		{name: "pre-submitted expire", initial: StatePreSubmitted, event: EventExpire, want: StateExpired},
		{name: "submitted expire", initial: StateSubmitted, event: EventExpire, want: StateExpired},
		{name: "confirmed expire", initial: StateConfirmed, event: EventExpire, want: StateExpired},
		{name: "partial-filled expire", initial: StatePartialFilled, event: EventExpire, want: StateExpired},
		{name: "pending cannot fill", initial: StatePending, event: EventFill, want: StatePending, wantErr: true},
		{name: "unknown cannot fill", initial: StateUnknown, event: EventFill, want: StateUnknown, wantErr: true},
		{name: "submitted cannot confirm", initial: StateSubmitted, event: EventConfirm, want: StateSubmitted, wantErr: true},
		{name: "confirmed cannot reject", initial: StateConfirmed, event: EventReject, want: StateConfirmed, wantErr: true},
		{name: "partial cannot reject", initial: StatePartialFilled, event: EventReject, want: StatePartialFilled, wantErr: true},
		{name: "partial cannot replace", initial: StatePartialFilled, event: EventReplace, want: StatePartialFilled, wantErr: true},
		{name: "unknown event", initial: StatePending, event: Event("UNKNOWN"), want: StatePending, wantErr: true},
		{name: "invalid initial state", initial: State("INVALID"), event: EventPlace, want: State("INVALID"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			machine := New(tt.initial)
			got, err := machine.ApplyEvent(tt.event)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ApplyEvent(%q) error = %v, wantErr %v", tt.event, err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, ErrInvalidTransition) {
				t.Fatalf("ApplyEvent(%q) error = %v, want ErrInvalidTransition", tt.event, err)
			}
			if got != tt.want {
				t.Fatalf("ApplyEvent(%q) state = %q, want %q", tt.event, got, tt.want)
			}
			if state := machine.State(); state != tt.want {
				t.Fatalf("State() = %q, want %q", state, tt.want)
			}
		})
	}
}

func TestMachineNonMutatingEvents(t *testing.T) {
	states := []State{
		StateUnknown,
		StatePending,
		StateSubmitted,
		StatePreSubmitted,
		StateConfirmed,
		StatePartialFilled,
	}
	events := []Event{
		EventCancelRequest,
		EventCancelFailed,
		EventRejectFailed,
		EventModifyFailed,
	}

	for _, initial := range states {
		for _, event := range events {
			t.Run(fmt.Sprintf("%s/%s", initial, event), func(t *testing.T) {
				machine := New(initial)
				got, err := machine.ApplyEvent(event)
				if err != nil {
					t.Fatalf("ApplyEvent(%q) error = %v", event, err)
				}
				if got != initial {
					t.Fatalf("ApplyEvent(%q) state = %q, want %q", event, got, initial)
				}
			})
		}
	}
}

func TestMachineTerminalStatesAreIdempotent(t *testing.T) {
	terminals := []State{StateFilled, StateCancelled, StateFailed, StateExpired}
	events := []Event{
		EventPlace,
		EventAcknowledge,
		EventConfirm,
		EventFill,
		EventPartialFill,
		EventCancelRequest,
		EventCancel,
		EventCancelFailed,
		EventReject,
		EventRejectFailed,
		EventExpire,
		EventReplace,
		EventModifyFailed,
		Event("UNKNOWN"),
	}

	for _, initial := range terminals {
		for _, event := range events {
			t.Run(fmt.Sprintf("%s/%s", initial, event), func(t *testing.T) {
				var callbackCalls atomic.Int64
				machine := New(initial)
				machine.SetOnStateChanged(func(State, State) {
					callbackCalls.Add(1)
				})

				for range 2 {
					got, err := machine.ApplyEvent(event)
					if err != nil {
						t.Fatalf("ApplyEvent(%q) error = %v", event, err)
					}
					if got != initial {
						t.Fatalf("ApplyEvent(%q) state = %q, want terminal state %q", event, got, initial)
					}
				}
				if got := callbackCalls.Load(); got != 0 {
					t.Fatalf("callback called %d times, want 0", got)
				}
			})
		}
	}
}

func TestSceneTypeToEvent(t *testing.T) {
	tests := []struct {
		sceneType string
		want      Event
	}{
		{sceneType: "PLACE", want: EventPlace},
		{sceneType: "PLACE_FAILED", want: EventReject},
		{sceneType: "REJECTED", want: EventReject},
		{sceneType: "REJECT_FAIL", want: EventRejectFailed},
		{sceneType: "REJECT_FAILED", want: EventRejectFailed},
		{sceneType: "ACKNOWLEDGED", want: EventAcknowledge},
		{sceneType: "CONFIRM_SUCCESS", want: EventConfirm},
		{sceneType: "CANCEL_REQUESTED", want: EventCancelRequest},
		{sceneType: "CANCEL_SUCCESS", want: EventCancel},
		{sceneType: "CANCEL_FAIL", want: EventCancelFailed},
		{sceneType: "CANCEL_FAILED", want: EventCancelFailed},
		{sceneType: "MODIFY_SUCCESS", want: EventReplace},
		{sceneType: "MODIFY_FAIL", want: EventModifyFailed},
		{sceneType: "MODIFY_FAILED", want: EventModifyFailed},
		{sceneType: "FILLED", want: EventPartialFill},
		{sceneType: "PARTIAL_FILLED", want: EventPartialFill},
		{sceneType: "FINAL_FILLED", want: EventFill},
		{sceneType: "EXPIRED", want: EventExpire},
		{sceneType: "final_filled", want: EventFill},
		{sceneType: "other", want: Event("OTHER")},
	}

	for _, tt := range tests {
		t.Run(tt.sceneType, func(t *testing.T) {
			if got := SceneTypeToEvent(tt.sceneType); got != tt.want {
				t.Fatalf("SceneTypeToEvent(%q) = %q, want %q", tt.sceneType, got, tt.want)
			}
		})
	}
}

func TestOrderApplySceneType(t *testing.T) {
	tests := []struct {
		name      string
		initial   State
		sceneType string
		want      State
		wantErr   bool
	}{
		{name: "cancel request is not success", initial: StateConfirmed, sceneType: "CANCEL_REQUESTED", want: StateConfirmed},
		{name: "cancel success", initial: StateConfirmed, sceneType: "CANCEL_SUCCESS", want: StateCancelled},
		{name: "cancel failure", initial: StateConfirmed, sceneType: "CANCEL_FAILED", want: StateConfirmed},
		{name: "modify failure", initial: StateConfirmed, sceneType: "MODIFY_FAILED", want: StateConfirmed},
		{name: "reject failure", initial: StateSubmitted, sceneType: "REJECT_FAILED", want: StateSubmitted},
		{name: "partial fill", initial: StateConfirmed, sceneType: "FILLED", want: StatePartialFilled},
		{name: "final fill", initial: StateConfirmed, sceneType: "FINAL_FILLED", want: StateFilled},
		{name: "place failure", initial: StatePending, sceneType: "PLACE_FAILED", want: StateFailed},
		{name: "invalid known transition", initial: StatePending, sceneType: "CONFIRM_SUCCESS", want: StatePending, wantErr: true},
		{name: "unknown scene", initial: StatePending, sceneType: "UNKNOWN_SCENE", want: StatePending, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := &Order{Machine: New(tt.initial)}
			got, err := order.ApplySceneType(tt.sceneType)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ApplySceneType(%q) error = %v, wantErr %v", tt.sceneType, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("ApplySceneType(%q) state = %q, want %q", tt.sceneType, got, tt.want)
			}
			if state := order.Machine.State(); state != tt.want {
				t.Fatalf("Machine.State() = %q, want %q", state, tt.want)
			}
		})
	}
}

func TestOrderCancellationLifecycle(t *testing.T) {
	tests := []struct {
		sceneType string
		want      State
	}{
		{sceneType: "CANCEL_REQUESTED", want: StateConfirmed},
		{sceneType: "CANCEL_FAILED", want: StateConfirmed},
		{sceneType: "CANCEL_SUCCESS", want: StateCancelled},
		{sceneType: "CANCEL_REQUESTED", want: StateCancelled},
		{sceneType: "CANCEL_FAILED", want: StateCancelled},
		{sceneType: "CANCEL_SUCCESS", want: StateCancelled},
	}

	order := &Order{Machine: New(StateConfirmed)}
	for _, tt := range tests {
		got, err := order.ApplySceneType(tt.sceneType)
		if err != nil {
			t.Fatalf("ApplySceneType(%q) error = %v", tt.sceneType, err)
		}
		if got != tt.want {
			t.Fatalf("ApplySceneType(%q) state = %q, want %q", tt.sceneType, got, tt.want)
		}
	}
}

func TestOnStateChanged(t *testing.T) {
	tests := []struct {
		name      string
		initial   State
		event     Event
		want      State
		wantErr   bool
		wantCalls int
	}{
		{name: "state change", initial: StatePending, event: EventPlace, want: StateSubmitted, wantCalls: 1},
		{name: "cancel request", initial: StateConfirmed, event: EventCancelRequest, want: StateConfirmed},
		{name: "cancel failure", initial: StateConfirmed, event: EventCancelFailed, want: StateConfirmed},
		{name: "modify failure", initial: StateConfirmed, event: EventModifyFailed, want: StateConfirmed},
		{name: "self transition", initial: StatePartialFilled, event: EventPartialFill, want: StatePartialFilled},
		{name: "invalid transition", initial: StatePending, event: EventFill, want: StatePending, wantErr: true},
		{name: "terminal idempotency", initial: StateFilled, event: EventCancel, want: StateFilled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			machine := New(tt.initial)
			var calls atomic.Int64
			machine.SetOnStateChanged(func(from, to State) {
				if from != tt.initial {
					t.Errorf("callback from = %q, want %q", from, tt.initial)
				}
				if to != tt.want {
					t.Errorf("callback to = %q, want %q", to, tt.want)
				}
				calls.Add(1)
			})

			got, err := machine.ApplyEvent(tt.event)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ApplyEvent(%q) error = %v, wantErr %v", tt.event, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("ApplyEvent(%q) state = %q, want %q", tt.event, got, tt.want)
			}
			if got := calls.Load(); got != int64(tt.wantCalls) {
				t.Fatalf("callback called %d times, want %d", got, tt.wantCalls)
			}
		})
	}
}

func TestOnStateChangedCanReadMachine(t *testing.T) {
	machine := New(StatePending)
	observed := make(chan State, 1)
	machine.SetOnStateChanged(func(_, _ State) {
		observed <- machine.State()
	})

	if got, err := machine.ApplyEvent(EventPlace); err != nil || got != StateSubmitted {
		t.Fatalf("ApplyEvent(EventPlace) = (%q, %v), want (%q, nil)", got, err, StateSubmitted)
	}
	if got := <-observed; got != StateSubmitted {
		t.Fatalf("State() in callback = %q, want %q", got, StateSubmitted)
	}
}

func TestMachineConcurrentAccess(t *testing.T) {
	machine := New(StatePending)
	var callbackCalls atomic.Int64
	callback := func(State, State) {
		callbackCalls.Add(1)
	}
	machine.SetOnStateChanged(callback)

	const (
		workers    = 32
		iterations = 200
	)
	start := make(chan struct{})
	applyEvent := func(event Event) {
		_, _ = machine.ApplyEvent(event)
	}
	var wait sync.WaitGroup
	wait.Add(workers)
	for worker := range workers {
		go func(worker int) {
			defer wait.Done()
			<-start
			for iteration := range iterations {
				switch (worker + iteration) % 7 {
				case 0:
					applyEvent(EventPlace)
				case 1:
					applyEvent(EventAcknowledge)
				case 2:
					applyEvent(EventPartialFill)
				case 3:
					applyEvent(EventFill)
				case 4:
					machine.State()
				case 5:
					machine.SetOnStateChanged(callback)
				case 6:
					applyEvent(EventCancelRequest)
				}
			}
		}(worker)
	}
	close(start)
	wait.Wait()

	switch state := machine.State(); state {
	case StateUnknown, StatePending, StateSubmitted, StatePreSubmitted, StateConfirmed, StatePartialFilled, StateFilled, StateCancelled, StateFailed, StateExpired:
	default:
		t.Fatalf("State() = %q, want a known state", state)
	}
	if got := callbackCalls.Load(); got == 0 {
		t.Fatal("state-change callback was never called")
	}
}

func TestIsTerminal(t *testing.T) {
	tests := []struct {
		state State
		want  bool
	}{
		{state: StateUnknown},
		{state: StatePending},
		{state: StateSubmitted},
		{state: StatePreSubmitted},
		{state: StateConfirmed},
		{state: StatePartialFilled},
		{state: StateFilled, want: true},
		{state: StateCancelled, want: true},
		{state: StateFailed, want: true},
		{state: StateExpired, want: true},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			if got := IsTerminal(tt.state); got != tt.want {
				t.Fatalf("IsTerminal(%q) = %v, want %v", tt.state, got, tt.want)
			}
		})
	}
}

func TestFromWebullStatus(t *testing.T) {
	tests := []struct {
		status string
		want   State
	}{
		{status: "PENDING", want: StatePending},
		{status: "SUBMITTED", want: StateSubmitted},
		{status: "PRE_SUBMITTED", want: StatePreSubmitted},
		{status: "CONFIRMED", want: StateConfirmed},
		{status: "PARTIAL_FILLED", want: StatePartialFilled},
		{status: "FILLED", want: StateFilled},
		{status: "CANCELLED", want: StateCancelled},
		{status: "FAILED", want: StateFailed},
		{status: "EXPIRED", want: StateExpired},
		{status: "UNRECOGNIZED", want: StateUnknown},
		{status: "", want: StateUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			if got := FromWebullStatus(tt.status); got != tt.want {
				t.Fatalf("FromWebullStatus(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}
