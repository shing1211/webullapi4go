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

package events

import (
	"encoding/json"

	eventsevents "github.com/shing1211/webullapi4go/gen/webull/trade/events/v1"
	"github.com/shing1211/webullapi4go/internal/errs"
)

// contentTypeJSON is the MIME type the server uses for event payloads that are
// JSON documents. A data event whose content type is anything else is delivered
// only to the raw [Client.OnEvent] handlers.
const contentTypeJSON = "application/json"

// OrderEvent is a decoded order status-change payload, the JSON carried by an
// [EventOrder] data event. Webull does not publish a formal schema, so the
// fields mirror the documented and observed payload keys and every numeric
// value is kept as a string to preserve precision. Unknown JSON keys are
// ignored and [OrderEvent.Raw] holds the exact bytes the server sent, so a
// consumer can always fall back to the wire form.
//
// An order event is emitted for placement, modification, cancellation, and
// fills. [OrderEvent.SceneType] distinguishes the scenario (for example
// "CANCEL_SUCCESS" or "FINAL_FILLED") and [OrderEvent.OrderStatus] carries the
// resulting lifecycle state.
type OrderEvent struct {
	// Raw is the exact JSON payload received from the server. It is always
	// populated for a decoded event and lets callers read fields that this
	// version does not model.
	Raw json.RawMessage `json:"-"`
	// SecAccountID is the numeric securities account identifier. The server
	// encodes it as a bare JSON number; json.Number preserves it either as a
	// number or, if the encoding ever changes, as a string.
	SecAccountID json.Number `json:"secAccountId"`
	// EventRequestID is the event-level request identifier (the JSON key
	// "requestId"), distinct from the order request identifier.
	EventRequestID string `json:"requestId"`
	// AccountID is the account the order belongs to.
	AccountID string `json:"account_id"`
	// RequestID is the order request identifier (the JSON key "request_id").
	RequestID string `json:"request_id"`
	// OrderID is the system-generated order identifier.
	OrderID string `json:"order_id"`
	// ClientOrderID is the caller-supplied order identifier.
	ClientOrderID string `json:"client_order_id"`
	// InstrumentID is the Webull instrument identifier.
	InstrumentID string `json:"instrument_id"`
	// OrderStatus is the resulting lifecycle state, for example SUBMITTED,
	// FILLED, or CANCELLED.
	OrderStatus string `json:"order_status"`
	// Symbol is the trading symbol of the instrument.
	Symbol string `json:"symbol"`
	// ShortName is the human-readable instrument name, when supplied.
	ShortName string `json:"short_name"`
	// Quantity is the total order quantity, as a decimal string.
	Quantity string `json:"qty"`
	// FilledQuantity is the quantity executed so far, as a decimal string.
	FilledQuantity string `json:"filled_qty"`
	// FilledPrice is the average execution price, as a decimal string.
	FilledPrice string `json:"filled_price"`
	// FilledTime is the time of the last execution, when supplied.
	FilledTime string `json:"filled_time"`
	// Side is the order direction, for example BUY or SELL.
	Side string `json:"side"`
	// SceneType is the event scenario, for example CANCEL_SUCCESS, FILLED, or
	// FINAL_FILLED.
	SceneType string `json:"scene_type"`
	// Category is the instrument category, for example US_STOCK or US_EVENT.
	Category string `json:"category"`
	// OrderType is the execution instruction, for example LIMIT.
	OrderType string `json:"order_type"`
	// BizType is the business type, fixed as TRADE for order events.
	BizType string `json:"biz_type"`
	// ActualCommission is the commission collected, as a decimal string.
	ActualCommission string `json:"actual_commission"`
	// ReceivableCommission is the commission still receivable, as a decimal
	// string.
	ReceivableCommission string `json:"receivable_commission"`
}

// PositionEvent is a decoded event-contract position settlement payload, the
// JSON carried by an [EventPosition] data event. Webull does not publish a
// formal schema, so the fields mirror the documented payload keys and every
// numeric value is kept as a string. Unknown JSON keys are ignored and
// [PositionEvent.Raw] holds the exact bytes the server sent.
type PositionEvent struct {
	// Raw is the exact JSON payload received from the server.
	Raw json.RawMessage `json:"-"`
	// PositionID is the position identifier.
	PositionID string `json:"position_id"`
	// AccountID is the account the position belongs to.
	AccountID string `json:"account_id"`
	// Symbol is the trading symbol of the event contract.
	Symbol string `json:"symbol"`
	// EventName is the human-readable event contract name.
	EventName string `json:"event_name"`
	// YesCondition describes the condition for a YES outcome.
	YesCondition string `json:"yes_condition"`
	// SettleResult is the settlement result of the event contract.
	SettleResult string `json:"settle_result"`
	// SettleSide is the event outcome decision.
	SettleSide string `json:"settle_side"`
	// Quantity is the held quantity, as a decimal string.
	Quantity string `json:"quantity"`
	// Cost is the total cost basis, as a decimal string.
	Cost string `json:"cost"`
	// SettleAmount is the settlement amount, as a decimal string.
	SettleAmount string `json:"settle_amount"`
	// BizType is the business type, fixed as EVENT_POSITION_SETTLED.
	BizType string `json:"biz_type"`
}

// OptionEvent is a decoded option status-change payload, the JSON carried by an
// [EventOption] data event. Webull does not publish a schema for option events;
// they are order status changes for an option contract and carry the same keys
// as [OrderEvent], distinguished by [OptionEvent.Category] (usually
// "US_OPTION"). Unknown JSON keys are ignored and [OptionEvent.Raw] holds the
// exact bytes the server sent, so an unmodelled option-specific key remains
// accessible.
type OptionEvent struct {
	// Raw is the exact JSON payload received from the server.
	Raw json.RawMessage `json:"-"`
	// SecAccountID is the numeric securities account identifier, preserved as
	// a JSON number or string.
	SecAccountID json.Number `json:"secAccountId"`
	// EventRequestID is the event-level request identifier (the JSON key
	// "requestId").
	EventRequestID string `json:"requestId"`
	// AccountID is the account the option order belongs to.
	AccountID string `json:"account_id"`
	// RequestID is the order request identifier (the JSON key "request_id").
	RequestID string `json:"request_id"`
	// OrderID is the system-generated order identifier.
	OrderID string `json:"order_id"`
	// ClientOrderID is the caller-supplied order identifier.
	ClientOrderID string `json:"client_order_id"`
	// InstrumentID is the Webull instrument identifier of the option.
	InstrumentID string `json:"instrument_id"`
	// OrderStatus is the resulting lifecycle state.
	OrderStatus string `json:"order_status"`
	// Symbol is the option contract symbol.
	Symbol string `json:"symbol"`
	// ShortName is the human-readable instrument name, when supplied.
	ShortName string `json:"short_name"`
	// Quantity is the total order quantity, as a decimal string.
	Quantity string `json:"qty"`
	// FilledQuantity is the quantity executed so far, as a decimal string.
	FilledQuantity string `json:"filled_qty"`
	// FilledPrice is the average execution price, as a decimal string.
	FilledPrice string `json:"filled_price"`
	// FilledTime is the time of the last execution, when supplied.
	FilledTime string `json:"filled_time"`
	// Side is the order direction, for example BUY or SELL.
	Side string `json:"side"`
	// SceneType is the event scenario, for example CANCEL_SUCCESS.
	SceneType string `json:"scene_type"`
	// Category is the instrument category; option events report US_OPTION.
	Category string `json:"category"`
	// OrderType is the execution instruction, for example LIMIT.
	OrderType string `json:"order_type"`
	// BizType is the business type, fixed as TRADE for order events.
	BizType string `json:"biz_type"`
	// ActualCommission is the commission collected, as a decimal string.
	ActualCommission string `json:"actual_commission"`
	// ReceivableCommission is the commission still receivable, as a decimal
	// string.
	ReceivableCommission string `json:"receivable_commission"`
}

// OnOrder registers a handler invoked for every decoded order event. Handlers
// are called from the stream pump and must not block; a slow handler delays
// later events. A payload that fails to decode is reported to [Client.OnError]
// instead and the stream continues.
func (c *Client) OnOrder(fn func(*OrderEvent)) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onOrder = append(c.onOrder, fn)
	c.mu.Unlock()
}

// OnPosition registers a handler invoked for every decoded event-contract
// position event. Handlers are called from the stream pump and must not block.
// A payload that fails to decode is reported to [Client.OnError] instead and
// the stream continues.
func (c *Client) OnPosition(fn func(*PositionEvent)) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onPosition = append(c.onPosition, fn)
	c.mu.Unlock()
}

// OnOption registers a handler invoked for every decoded option event.
// Handlers are called from the stream pump and must not block. A payload that
// fails to decode is reported to [Client.OnError] instead and the stream
// continues.
func (c *Client) OnOption(fn func(*OptionEvent)) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onOption = append(c.onOption, fn)
	c.mu.Unlock()
}

// routeDataEvent delivers one data event to the raw [Client.OnEvent] handlers
// and, when the payload is JSON, to the matching typed handler. A decoding
// failure is reported to [Client.OnError] and never ends the stream.
func (c *Client) routeDataEvent(resp *eventsevents.SubscribeResponse) {
	kind := uint32(resp.GetEventType())
	contentType := resp.GetContentType()
	payload := []byte(resp.GetPayload())

	c.emitEvent(kind, contentType, payload)

	if contentType != contentTypeJSON || len(payload) == 0 {
		return
	}
	switch kind {
	case EventOrder:
		var ev OrderEvent
		if err := decodeEvent(payload, &ev); err != nil {
			c.emitError(err)
			return
		}
		ev.Raw = cloneRaw(payload)
		c.emitOrder(&ev)
	case EventPosition:
		var ev PositionEvent
		if err := decodeEvent(payload, &ev); err != nil {
			c.emitError(err)
			return
		}
		ev.Raw = cloneRaw(payload)
		c.emitPosition(&ev)
	case EventOption:
		var ev OptionEvent
		if err := decodeEvent(payload, &ev); err != nil {
			c.emitError(err)
			return
		}
		ev.Raw = cloneRaw(payload)
		c.emitOption(&ev)
	}
}

// decodeEvent unmarshals a JSON payload into out, returning a typed error that
// carries the payload's event kind so callers can locate the failure.
func decodeEvent(payload []byte, out any) error {
	if err := json.Unmarshal(payload, out); err != nil {
		return errs.Wrap(errs.CodeAPI, "events: decode JSON payload", err)
	}
	return nil
}

// cloneRaw copies a payload so the decoded event does not retain the backing
// array of the protobuf response.
func cloneRaw(payload []byte) json.RawMessage {
	return append(json.RawMessage(nil), payload...)
}

// emitOrder invokes every registered order handler outside the lock.
func (c *Client) emitOrder(ev *OrderEvent) {
	c.mu.RLock()
	handlers := make([]func(*OrderEvent), len(c.onOrder))
	copy(handlers, c.onOrder)
	c.mu.RUnlock()
	for _, h := range handlers {
		h(ev)
	}
}

// emitPosition invokes every registered position handler outside the lock.
func (c *Client) emitPosition(ev *PositionEvent) {
	c.mu.RLock()
	handlers := make([]func(*PositionEvent), len(c.onPosition))
	copy(handlers, c.onPosition)
	c.mu.RUnlock()
	for _, h := range handlers {
		h(ev)
	}
}

// emitOption invokes every registered option handler outside the lock.
func (c *Client) emitOption(ev *OptionEvent) {
	c.mu.RLock()
	handlers := make([]func(*OptionEvent), len(c.onOption))
	copy(handlers, c.onOption)
	c.mu.RUnlock()
	for _, h := range handlers {
		h(ev)
	}
}
