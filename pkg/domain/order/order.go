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
	EventReject        Event = "REJECTED"
	EventExpire        Event = "EXPIRED"
	EventReplace       Event = "REPLACED"
)

var (
	ErrInvalidTransition = errors.New("order: invalid state transition")
)

type OnStateChanged func(from, to State)

type Machine struct {
	state     State
	onChanged OnStateChanged
}

func New(initial State) *Machine {
	return &Machine{state: initial}
}

func (m *Machine) State() State {
	return m.state
}

func (m *Machine) SetOnStateChanged(fn OnStateChanged) {
	m.onChanged = fn
}

func (m *Machine) ApplyEvent(e Event) (State, error) {
	next, err := transition(m.state, e)
	if err != nil {
		return m.state, err
	}
	prev := m.state
	m.state = next
	if m.onChanged != nil {
		m.onChanged(prev, next)
	}
	return next, nil
}

func transition(s State, e Event) (State, error) {
	switch s {
	case StateUnknown, StatePending:
		switch e {
		case EventPlace:
			return StateSubmitted, nil
		case EventAcknowledge:
			return StatePreSubmitted, nil
		case EventReject:
			return StateFailed, nil
		case EventCancelRequest:
			return StateCancelled, nil
		case EventExpire:
			return StateExpired, nil
		}
	case StatePreSubmitted:
		switch e {
		case EventConfirm:
			return StateConfirmed, nil
		case EventAcknowledge:
			return StatePreSubmitted, nil
		case EventFill:
			return StateFilled, nil
		case EventPartialFill:
			return StatePartialFilled, nil
		case EventCancelRequest:
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
		case EventFill:
			return StateFilled, nil
		case EventPartialFill:
			return StatePartialFilled, nil
		case EventCancelRequest:
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
		case EventCancelRequest:
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
		case EventCancelRequest:
			return StateCancelled, nil
		case EventExpire:
			return StateExpired, nil
		}
	case StateFilled, StateCancelled, StateFailed, StateExpired:
		return s, nil
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
	switch s {
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

// ApplySceneType applies a state transition for the given Webull scene type string
// (for example "CANCEL_SUCCESS", "FILLED", "FINAL_FILLED"). It returns the new
// state and an error if the scene type is unknown or the transition is invalid.
func (o *Order) ApplySceneType(sceneType string) (State, error) {
	return o.Machine.ApplyEvent(SceneTypeToEvent(sceneType))
}

// SceneTypeToEvent maps a Webull gRPC scene type string to an order [Event].
// Scene types are documented in the Webull OpenAPI event payload schema.
func SceneTypeToEvent(sceneType string) Event {
	switch sceneType {
	case "PLACE":
		return EventPlace
	case "ACKNOWLEDGED":
		return EventAcknowledge
	case "CONFIRM_SUCCESS":
		return EventConfirm
	case "CANCEL_SUCCESS":
		return EventCancel
	case "CANCEL_FAIL":
		return EventCancel
	case "MODIFY_SUCCESS":
		return EventReplace
	case "MODIFY_FAIL":
		return EventReplace
	case "FILLED", "FINAL_FILLED":
		return EventFill
	case "PARTIAL_FILLED":
		return EventPartialFill
	case "REJECTED", "REJECT_FAIL":
		return EventReject
	case "EXPIRED":
		return EventExpire
	case "CANCEL_REQUESTED":
		return EventCancelRequest
	default:
		return Event(strings.ToUpper(sceneType))
	}
}
