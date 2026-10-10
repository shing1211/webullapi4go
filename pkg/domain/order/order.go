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

// Package order implements the SDK's order management domain model: a
// validated state machine that enforces the legal transitions of a Webull
// order through its lifecycle.
package order

import (
	"errors"
	"strings"
	"sync"
)

type State string

const (
	StateUnknown       State = ""
	StatePending       State = "PENDING"
	StateSubmitted     State = "SUBMITTED"
	StatePreSubmitted  State = "PRE_SUBMITTED"
	StateConfirmed     State = "CONFIRMED"
	StatePartialFilled State = "PARTIAL_FILLED"
	StateFilled        State = "FILLED"
	StateCancelled     State = "CANCELLED"
	StateFailed        State = "FAILED"
	StateExpired       State = "EXPIRED"
)

type Event string

const (
	EventPlace         Event = "PLACED"
	EventAcknowledge   Event = "ACKNOWLEDGED"
	EventConfirm       Event = "CONFIRMED"
	EventFill          Event = "FILL"
	EventPartialFill   Event = "PARTIAL_FILL"
	EventCancelRequest Event = "CANCEL_REQUESTED"
	EventCancel        Event = "CANCELLED"
	// EventCancelFailed reports that a cancellation request did not succeed.
	// It does not change the order state.
	EventCancelFailed Event = "CANCEL_FAILED"
	EventReject       Event = "REJECTED"
	// EventRejectFailed reports that a rejection operation did not succeed.
	// It does not change the order state.
	EventRejectFailed Event = "REJECT_FAILED"
	EventExpire       Event = "EXPIRED"
	EventReplace      Event = "REPLACED"
	// EventModifyFailed reports that a replacement request did not succeed.
	// It does not change the order state.
	EventModifyFailed Event = "MODIFY_FAILED"
)

var (
	// ErrInvalidTransition indicates that an event is not legal for the
	// machine's current state.
	ErrInvalidTransition = errors.New("order: invalid state transition")
	// ErrInvalidStatus indicates that a status value cannot be mapped to an
	// order state.
	ErrInvalidStatus = errors.New("order: invalid status")
	// ErrInvalidOrder indicates that an order or its state machine is nil.
	ErrInvalidOrder = errors.New("order: order is not initialized")
)

// OnStateChanged is called after a machine changes state. The callback runs
// without the machine lock and concurrent transitions may invoke it concurrently.
type OnStateChanged func(from, to State)

// Machine tracks an order lifecycle and is safe for concurrent use.
type Machine struct {
	mu        sync.RWMutex
	state     State
	onChanged OnStateChanged
}

// New returns a machine initialized to initial.
func New(initial State) *Machine {
	return &Machine{state: initial}
}

// State returns the current order state.
func (m *Machine) State() State {
	if m == nil {
		return StateUnknown
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

// SetOnStateChanged installs fn for future state changes. Passing nil disables
// the callback.
func (m *Machine) SetOnStateChanged(fn OnStateChanged) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onChanged = fn
}

// ApplyEvent applies e and returns the resulting state. An invalid transition
// leaves the state unchanged and returns [ErrInvalidTransition].
func (m *Machine) ApplyEvent(e Event) (State, error) {
	if m == nil {
		return StateUnknown, ErrInvalidOrder
	}
	m.mu.Lock()
	prev := m.state
	next, err := transition(prev, e)
	if err != nil {
		m.mu.Unlock()
		return prev, err
	}
	m.state = next
	onChanged := m.onChanged
	m.mu.Unlock()

	if onChanged != nil && prev != next {
		onChanged(prev, next)
	}
	return next, nil
}

// ReconcileStatus advances m to the state represented by status using the
// smallest valid event path, or records an authoritative terminal snapshot
// when the normal event path is unavailable. A stale non-terminal status is
// treated as a no-op, and a terminal state is never regressed. It returns the
// resulting state and [ErrInvalidStatus] when status is not recognized.
func (m *Machine) ReconcileStatus(status State) (State, error) {
	if m == nil {
		return StateUnknown, ErrInvalidOrder
	}
	if status == StateUnknown {
		return m.State(), ErrInvalidStatus
	}
	m.mu.Lock()
	prev := m.state
	if prev == status || IsTerminal(prev) {
		m.mu.Unlock()
		return prev, nil
	}
	events, ok := statusEvents(prev, status)
	if !ok {
		if IsTerminal(status) {
			m.state = status
			onChanged := m.onChanged
			m.mu.Unlock()
			if onChanged != nil && prev != status {
				onChanged(prev, status)
			}
			return status, nil
		}
		m.mu.Unlock()
		return prev, ErrInvalidTransition
	}
	next := prev
	for _, event := range events {
		var err error
		next, err = transition(next, event)
		if err != nil {
			m.mu.Unlock()
			return prev, err
		}
	}
	m.state = next
	onChanged := m.onChanged
	m.mu.Unlock()

	if onChanged != nil && prev != next {
		onChanged(prev, next)
	}
	return next, nil
}

func statusEvents(current, target State) ([]Event, bool) {
	if current == StateUnknown && target == StatePending {
		return nil, true
	}
	if current != StateUnknown && !IsTerminal(target) && stateProgress(target) < stateProgress(current) {
		return nil, true
	}
	switch target {
	case StateSubmitted:
		if current == StateUnknown || current == StatePending {
			return []Event{EventPlace}, true
		}
	case StatePreSubmitted:
		if current == StateUnknown || current == StatePending {
			return []Event{EventAcknowledge}, true
		}
	case StateConfirmed:
		switch current {
		case StateUnknown, StatePending:
			return []Event{EventAcknowledge, EventConfirm}, true
		case StateSubmitted:
			return []Event{EventAcknowledge}, true
		case StatePreSubmitted:
			return []Event{EventConfirm}, true
		}
	case StatePartialFilled:
		switch current {
		case StateUnknown, StatePending:
			return []Event{EventPlace, EventPartialFill}, true
		case StateSubmitted, StatePreSubmitted, StateConfirmed:
			return []Event{EventPartialFill}, true
		}
	case StateFilled:
		switch current {
		case StateUnknown, StatePending:
			return []Event{EventPlace, EventFill}, true
		case StateSubmitted, StatePreSubmitted, StateConfirmed, StatePartialFilled:
			return []Event{EventFill}, true
		}
	case StateCancelled:
		if current == StateUnknown || current == StatePending || current == StateSubmitted || current == StatePreSubmitted || current == StateConfirmed || current == StatePartialFilled {
			return []Event{EventCancel}, true
		}
	case StateFailed:
		if current == StateUnknown || current == StatePending || current == StateSubmitted || current == StatePreSubmitted {
			return []Event{EventReject}, true
		}
	case StateExpired:
		if current == StateUnknown || current == StatePending || current == StateSubmitted || current == StatePreSubmitted || current == StateConfirmed || current == StatePartialFilled {
			return []Event{EventExpire}, true
		}
	}
	return nil, false
}

func stateProgress(s State) int {
	switch s {
	case StatePending:
		return 1
	case StateSubmitted:
		return 2
	case StatePreSubmitted:
		return 3
	case StateConfirmed:
		return 4
	case StatePartialFilled:
		return 5
	default:
		return 0
	}
}

func transition(s State, e Event) (State, error) {
	if IsTerminal(s) {
		return s, nil
	}
	switch e {
	case EventCancelRequest, EventCancelFailed, EventRejectFailed, EventModifyFailed:
		switch s {
		case StateUnknown, StatePending, StateSubmitted, StatePreSubmitted, StateConfirmed, StatePartialFilled:
			return s, nil
		}
	}

	switch s {
	case StateUnknown, StatePending:
		switch e {
		case EventPlace:
			return StateSubmitted, nil
		case EventAcknowledge:
			return StatePreSubmitted, nil
		case EventReplace:
			return StateSubmitted, nil
		case EventCancel:
			return StateCancelled, nil
		case EventReject:
			return StateFailed, nil
		case EventExpire:
			return StateExpired, nil
		}
	case StatePreSubmitted:
		switch e {
		case EventConfirm:
			return StateConfirmed, nil
		case EventAcknowledge:
			return StatePreSubmitted, nil
		case EventReplace:
			return StateSubmitted, nil
		case EventFill:
			return StateFilled, nil
		case EventPartialFill:
			return StatePartialFilled, nil
		case EventCancel:
			return StateCancelled, nil
		case EventReject:
			return StateFailed, nil
		case EventExpire:
			return StateExpired, nil
		}
	case StateSubmitted:
		switch e {
		case EventAcknowledge:
			return StateConfirmed, nil
		case EventReplace:
			return StateSubmitted, nil
		case EventFill:
			return StateFilled, nil
		case EventPartialFill:
			return StatePartialFilled, nil
		case EventCancel:
			return StateCancelled, nil
		case EventReject:
			return StateFailed, nil
		case EventExpire:
			return StateExpired, nil
		}
	case StateConfirmed:
		switch e {
		case EventFill:
			return StateFilled, nil
		case EventPartialFill:
			return StatePartialFilled, nil
		case EventCancel:
			return StateCancelled, nil
		case EventReplace:
			return StateSubmitted, nil
		case EventExpire:
			return StateExpired, nil
		}
	case StatePartialFilled:
		switch e {
		case EventFill:
			return StateFilled, nil
		case EventPartialFill:
			return StatePartialFilled, nil
		case EventCancel:
			return StateCancelled, nil
		case EventExpire:
			return StateExpired, nil
		}
	}
	return s, ErrInvalidTransition
}

func IsTerminal(s State) bool {
	switch s {
	case StateFilled, StateCancelled, StateFailed, StateExpired:
		return true
	}
	return false
}

func FromWebullStatus(s string) State {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "PENDING":
		return StatePending
	case "SUBMITTED":
		return StateSubmitted
	case "PRE_SUBMITTED":
		return StatePreSubmitted
	case "CONFIRMED":
		return StateConfirmed
	case "PARTIAL_FILLED":
		return StatePartialFilled
	case "FILLED":
		return StateFilled
	case "CANCELLED":
		return StateCancelled
	case "FAILED":
		return StateFailed
	case "EXPIRED":
		return StateExpired
	default:
		return StateUnknown
	}
}

// Order is a tracked Webull order with a local state machine. It is returned by
// [Client.PlaceOrder] and updated by applying events from the Webull order event
// stream. The embedded [PlaceOrderResult] provides field-level compatibility with
// the wire response.
type Order struct {
	PlaceOrderResult
	// AccountID is the securities account the order belongs to.
	AccountID string `json:"account_id"`
	// Machine is the local state machine tracking the order's lifecycle.
	Machine *Machine
}

// PlaceOrderResult mirrors the fields returned by the Webull place-order endpoint.
type PlaceOrderResult struct {
	// ClientOrderID echoes the caller-supplied order identifier.
	ClientOrderID string `json:"client_order_id"`
	// OrderID is the system-generated order identifier.
	OrderID string `json:"order_id"`
}

// State returns the order's current lifecycle state. It returns
// [StateUnknown] when the order is nil or has no state machine.
func (o *Order) State() State {
	if o == nil || o.Machine == nil {
		return StateUnknown
	}
	return o.Machine.State()
}

// ApplyEvent applies a domain order event to the order's state machine. It
// leaves the state unchanged when the event is not a legal transition.
func (o *Order) ApplyEvent(event Event) (State, error) {
	if o == nil || o.Machine == nil {
		return StateUnknown, ErrInvalidOrder
	}
	return o.Machine.ApplyEvent(event)
}

// ReconcileStatus reconciles the order with a Webull status string using
// [Machine.ReconcileStatus]. Unknown status values return [ErrInvalidStatus].
func (o *Order) ReconcileStatus(status string) (State, error) {
	if o == nil || o.Machine == nil {
		return StateUnknown, ErrInvalidOrder
	}
	return o.Machine.ReconcileStatus(FromWebullStatus(status))
}

// ReconcileState reconciles the order with an already normalized [State].
func (o *Order) ReconcileState(state State) (State, error) {
	if o == nil || o.Machine == nil {
		return StateUnknown, ErrInvalidOrder
	}
	return o.Machine.ReconcileStatus(state)
}

// ApplySceneType applies a state transition for the given Webull scene type string
// (for example "CANCEL_SUCCESS", "FILLED", "FINAL_FILLED"). It returns the new
// state and an error if the scene type is unknown or the transition is invalid.
func (o *Order) ApplySceneType(sceneType string) (State, error) {
	return o.ApplyEvent(SceneTypeToEvent(sceneType))
}

// SceneTypeToEvent maps a Webull gRPC scene type string to an order [Event].
// Failure scenes have distinct events so callers do not mistake failed cancel,
// modify, or reject operations for successful transitions.
func SceneTypeToEvent(sceneType string) Event {
	normalized := strings.ToUpper(sceneType)
	switch normalized {
	case "PLACE":
		return EventPlace
	case "PLACE_FAILED", "REJECTED":
		return EventReject
	case "REJECT_FAIL", "REJECT_FAILED":
		return EventRejectFailed
	case "ACKNOWLEDGED":
		return EventAcknowledge
	case "CONFIRM_SUCCESS":
		return EventConfirm
	case "CANCEL_REQUESTED":
		return EventCancelRequest
	case "CANCEL_SUCCESS":
		return EventCancel
	case "CANCEL_FAIL", "CANCEL_FAILED":
		return EventCancelFailed
	case "MODIFY_SUCCESS":
		return EventReplace
	case "MODIFY_FAIL", "MODIFY_FAILED":
		return EventModifyFailed
	case "FILLED", "PARTIAL_FILLED":
		return EventPartialFill
	case "FINAL_FILLED":
		return EventFill
	case "EXPIRED":
		return EventExpire
	default:
		return Event(normalized)
	}
}
