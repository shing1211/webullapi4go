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

package events_test

import (
	"context"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/events"
	eventsevents "github.com/shing1211/webullapi4go/gen/webull/trade/events/v1"
	"github.com/shing1211/webullapi4go/internal/errs"
)

// orderEventJSON is a representative order status-change payload, modelled on a
// live sandbox CANCEL_SUCCESS event and Webull's documented order-event fields.
// The trailing key is not modelled and must be ignored rather than rejected.
const orderEventJSON = `{"secAccountId":66600004338,"requestId":"EVENT-REQ-1","account_id":"acct-1","request_id":"ORDER-REQ-1","order_id":"ORDER-1","client_order_id":"cli-1","instrument_id":"913256135","order_status":"CANCELLED","symbol":"AAPL","short_name":"Apple Inc","qty":"1.0000000000","filled_price":"0E-10","filled_qty":"0E-10","side":"BUY","scene_type":"CANCEL_SUCCESS","category":"US_STOCK","order_type":"LIMIT","actual_commission":"0E-10","receivable_commission":"0E-10","future_field":"ignored"}`

// positionEventJSON is a representative event-contract position payload.
const positionEventJSON = `{"position_id":"POS-1","account_id":"acct-1","symbol":"KXRATECUTCOUNT-26DEC31-T14","event_name":"Number of rate cuts in 2025","yes_condition":"Exactly 8 cuts","settle_result":"Yes","settle_side":"Yes","quantity":"40","cost":"20.00","settle_amount":"40.00","biz_type":"EVENT_POSITION_SETTLED"}`

// optionEventJSON is a representative option status-change payload. Option
// events reuse the order schema and report category US_OPTION.
const optionEventJSON = `{"secAccountId":66600004338,"requestId":"EVENT-REQ-2","account_id":"acct-1","request_id":"ORDER-REQ-2","order_id":"ORDER-2","client_order_id":"cli-2","instrument_id":"470059643","order_status":"SUBMITTED","symbol":"AAPL260522C00300000","qty":"1","filled_qty":"0","side":"BUY","scene_type":"MODIFY_SUCCESS","category":"US_OPTION","order_type":"LIMIT"}`

// orderResponse builds a data response carrying payload with the given event
// kind.
func dataResponse(kind uint32, contentType, payload string) *eventsevents.SubscribeResponse {
	return &eventsevents.SubscribeResponse{
		EventType:   eventsevents.EventType(kind),
		ContentType: contentType,
		Payload:     payload,
	}
}

// TestDecodesOrderEvent verifies that an order data event decodes into an
// *OrderEvent with the documented fields, preserves the raw bytes, and ignores
// unknown keys.
func TestDecodesOrderEvent(t *testing.T) {
	fake := &fakeServer{messages: []*eventsevents.SubscribeResponse{
		{EventType: eventsevents.EventType_SubscribeSuccess},
		dataResponse(events.EventOrder, "application/json", orderEventJSON),
	}}
	cl := newClient(t, fake)

	orders := make(chan *events.OrderEvent, 1)
	cl.OnOrder(func(ev *events.OrderEvent) { trySend(orders, ev) })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = cl.Run(ctx) }()

	var ev *events.OrderEvent
	select {
	case ev = <-orders:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnOrder")
	}

	if ev.AccountID != "acct-1" {
		t.Errorf("AccountID = %q, want acct-1", ev.AccountID)
	}
	if ev.OrderID != "ORDER-1" || ev.ClientOrderID != "cli-1" {
		t.Errorf("order/client ids = %q/%q", ev.OrderID, ev.ClientOrderID)
	}
	if ev.OrderStatus != "CANCELLED" || ev.SceneType != "CANCEL_SUCCESS" {
		t.Errorf("status/scene = %q/%q", ev.OrderStatus, ev.SceneType)
	}
	if ev.Quantity != "1.0000000000" || ev.FilledQuantity != "0E-10" {
		t.Errorf("qty/filled = %q/%q", ev.Quantity, ev.FilledQuantity)
	}
	if ev.Symbol != "AAPL" || ev.Category != "US_STOCK" || ev.OrderType != "LIMIT" {
		t.Errorf("symbol/category/type = %q/%q/%q", ev.Symbol, ev.Category, ev.OrderType)
	}
	if got := ev.SecAccountID.String(); got != "66600004338" {
		t.Errorf("SecAccountID = %q, want 66600004338", got)
	}
	if string(ev.Raw) != orderEventJSON {
		t.Errorf("Raw = %q, want the original payload", ev.Raw)
	}
}

// TestDecodesPositionAndOptionEvents verifies the position and option routes.
func TestDecodesPositionAndOptionEvents(t *testing.T) {
	fake := &fakeServer{messages: []*eventsevents.SubscribeResponse{
		{EventType: eventsevents.EventType_SubscribeSuccess},
		dataResponse(events.EventPosition, "application/json", positionEventJSON),
		dataResponse(events.EventOption, "application/json", optionEventJSON),
	}}
	cl := newClient(t, fake)

	positions := make(chan *events.PositionEvent, 1)
	options := make(chan *events.OptionEvent, 1)
	cl.OnPosition(func(ev *events.PositionEvent) { trySend(positions, ev) })
	cl.OnOption(func(ev *events.OptionEvent) { trySend(options, ev) })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = cl.Run(ctx) }()

	select {
	case ev := <-positions:
		if ev.PositionID != "POS-1" || ev.SettleResult != "Yes" || ev.SettleAmount != "40.00" {
			t.Errorf("position = %+v", ev)
		}
		if ev.BizType != "EVENT_POSITION_SETTLED" {
			t.Errorf("position biz_type = %q", ev.BizType)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnPosition")
	}

	select {
	case ev := <-options:
		if ev.OrderID != "ORDER-2" || ev.Category != "US_OPTION" || ev.InstrumentID != "470059643" {
			t.Errorf("option = %+v", ev)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnOption")
	}
}

// TestMalformedPayloadReportsErrorAndContinues verifies that a malformed JSON
// data event is reported to OnError without ending the stream, and that a later
// valid event is still delivered.
func TestMalformedPayloadReportsErrorAndContinues(t *testing.T) {
	fake := &fakeServer{messages: []*eventsevents.SubscribeResponse{
		{EventType: eventsevents.EventType_SubscribeSuccess},
		dataResponse(events.EventOrder, "application/json", `{"order_id":`),
		dataResponse(events.EventOrder, "application/json", orderEventJSON),
	}}
	cl := newClient(t, fake)

	errCh := make(chan error, 4)
	orders := make(chan *events.OrderEvent, 2)
	cl.OnError(func(err error) { trySend(errCh, err) })
	cl.OnOrder(func(ev *events.OrderEvent) { trySend(orders, ev) })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = cl.Run(ctx) }()

	select {
	case err := <-errCh:
		if !errs.Is(err, errs.CodeAPI) {
			t.Errorf("OnError = %v, want an api code", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for the decode error")
	}

	select {
	case ev := <-orders:
		if ev.OrderID != "ORDER-1" {
			t.Errorf("OrderID = %q, want ORDER-1", ev.OrderID)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("stream did not continue after the malformed payload")
	}
}

// TestNonJSONPayloadSkipsTypedDecode verifies that a data event whose content
// type is not application/json reaches OnEvent but no typed handler.
func TestNonJSONPayloadSkipsTypedDecode(t *testing.T) {
	fake := &fakeServer{messages: []*eventsevents.SubscribeResponse{
		{EventType: eventsevents.EventType_SubscribeSuccess},
		dataResponse(events.EventOrder, "text/plain", orderEventJSON),
	}}
	cl := newClient(t, fake)

	raw := make(chan dataEvent, 1)
	orders := make(chan *events.OrderEvent, 1)
	cl.OnEvent(func(kind uint32, contentType string, payload []byte) {
		trySend(raw, dataEvent{kind: kind, contentType: contentType, payload: string(payload)})
	})
	cl.OnOrder(func(ev *events.OrderEvent) { trySend(orders, ev) })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = cl.Run(ctx) }()

	select {
	case got := <-raw:
		if got.contentType != "text/plain" {
			t.Errorf("contentType = %q", got.contentType)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for OnEvent")
	}
	select {
	case ev := <-orders:
		t.Errorf("OnOrder fired for a non-JSON payload: %+v", ev)
	case <-time.After(200 * time.Millisecond):
	}
}
