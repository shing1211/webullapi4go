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
	"context"
	"testing"
	"time"

	eventsevents "github.com/shing1211/webullapi4go/gen/webull/brokerfd/events/v1"
)

// TestOnDataEventDeliversRequestMetadata covers the additive metadata surface:
// OnDataEvent receives the whole DataEvent, so the server-assigned RequestId and
// Timestamp that OnData discards are available to callers.
func TestOnDataEventDeliversRequestMetadata(t *testing.T) {
	fake := &fakeServer{messages: []*eventsevents.SubscribeResponse{
		{EventType: eventsevents.EventType_SubscribeSuccess, Payload: `{"status":"ok"}`},
		{
			EventType:     1024,
			SubscribeType: 1024,
			ContentType:   "application/json",
			Payload:       `{"account_id":"acct-1"}`,
			RequestId:     "req-abc-123",
			Timestamp:     1758800000123,
		},
	}}
	cl := newClient(t, fake, WithAccounts([]string{"acct-1"}))

	got := make(chan *DataEvent, 4)
	errCh := make(chan error, 4)
	cl.OnDataEvent(func(de *DataEvent) { trySend(got, de) })
	cl.OnError(func(err error) { trySend(errCh, err) })

	startTestRun(t, cl, context.Background())

	var de *DataEvent
	select {
	case de = <-got:
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for OnDataEvent (errors: %v)", drainErr(errCh))
	}

	if de.SubscribeType != 1024 {
		t.Errorf("SubscribeType = %d, want 1024", de.SubscribeType)
	}
	if de.ContentType != "application/json" {
		t.Errorf("ContentType = %q, want application/json", de.ContentType)
	}
	if string(de.Payload) != `{"account_id":"acct-1"}` {
		t.Errorf("Payload = %q", de.Payload)
	}
	if de.RequestId != "req-abc-123" {
		t.Errorf("RequestId = %q, want req-abc-123", de.RequestId)
	}
	if de.Timestamp != 1758800000123 {
		t.Errorf("Timestamp = %d, want 1758800000123", de.Timestamp)
	}
}

// TestOnDataEventAndOnDataBothFire pins that the two registrations are
// independent: registering one does not suppress the other, and each event
// reaches both groups.
func TestOnDataEventAndOnDataBothFire(t *testing.T) {
	fake := &fakeServer{messages: []*eventsevents.SubscribeResponse{
		{EventType: eventsevents.EventType_SubscribeSuccess, Payload: `{"status":"ok"}`},
		{
			EventType:     1024,
			SubscribeType: 7,
			ContentType:   "application/json",
			Payload:       `{"n":1}`,
			RequestId:     "rid-1",
			Timestamp:     42,
		},
	}}
	cl := newClient(t, fake)

	legacy := make(chan receivedData, 4)
	rich := make(chan *DataEvent, 4)
	errCh := make(chan error, 4)
	cl.OnData(func(st uint32, ct string, payload []byte) {
		trySend(legacy, receivedData{st, ct, payload})
	})
	cl.OnDataEvent(func(de *DataEvent) { trySend(rich, de) })
	cl.OnError(func(err error) { trySend(errCh, err) })

	startTestRun(t, cl, context.Background())

	var lg receivedData
	select {
	case lg = <-legacy:
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for OnData (errors: %v)", drainErr(errCh))
	}
	var de *DataEvent
	select {
	case de = <-rich:
	case <-time.After(5 * time.Second):
		t.Fatalf("OnData did not fire alongside OnDataEvent (errors: %v)", drainErr(errCh))
	}

	if lg.subscribeType != 7 || lg.contentType != "application/json" || string(lg.payload) != `{"n":1}` {
		t.Errorf("OnData received %+v, want the same subscribe type, content type, and payload", lg)
	}
	if de.RequestId != "rid-1" || de.Timestamp != 42 {
		t.Errorf("OnDataEvent received %+v, want the request metadata", de)
	}
}

// TestOnDataEventPreservesRegistrationOrder pins the documented ordering
// guarantee within the OnDataEvent group.
func TestOnDataEventPreservesRegistrationOrder(t *testing.T) {
	fake := &fakeServer{messages: []*eventsevents.SubscribeResponse{
		{EventType: eventsevents.EventType_SubscribeSuccess, Payload: `{"status":"ok"}`},
		{EventType: 1024, SubscribeType: 1, ContentType: "application/json", Payload: `{"i":1}`, RequestId: "r1"},
	}}
	cl := newClient(t, fake)

	order := make(chan int, 8)
	done := make(chan struct{}, 1)
	cl.OnDataEvent(func(*DataEvent) { order <- 1 })
	cl.OnDataEvent(func(*DataEvent) { order <- 2 })
	cl.OnDataEvent(func(*DataEvent) { order <- 3; trySend(done, struct{}{}) })

	startTestRun(t, cl, context.Background())

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for the third OnDataEvent handler")
	}

	for want := 1; want <= 3; want++ {
		select {
		case got := <-order:
			if got != want {
				t.Fatalf("handler %d ran in position %d", got, want)
			}
		default:
			t.Fatalf("only %d of 3 handlers ran", want-1)
		}
	}
}

func TestOnDataEventIgnoresNilHandler(t *testing.T) {
	t.Parallel()

	cl := newClient(t, &fakeServer{})
	cl.OnDataEvent(nil)
	if len(cl.onDataEvent) != 0 {
		t.Fatalf("OnDataEvent(nil) registered %d handlers, want 0", len(cl.onDataEvent))
	}
}

// TestEmitDataIsSafeWithNoHandlers covers the empty-registration path for both
// groups.
func TestEmitDataIsSafeWithNoHandlers(t *testing.T) {
	t.Parallel()

	cl := newClient(t, &fakeServer{})
	resp := &SubscribeResponse{SubscribeResponse: &eventsevents.SubscribeResponse{
		EventType:     1024,
		SubscribeType: 3,
		ContentType:   "application/json",
		Payload:       `{"x":1}`,
		RequestId:     "r",
		Timestamp:     9,
	}}
	cl.emitData(resp)
}

// TestToDataEventCarriesEveryField pins the conversion that both handler groups
// depend on.
func TestToDataEventCarriesEveryField(t *testing.T) {
	t.Parallel()

	resp := &SubscribeResponse{SubscribeResponse: &eventsevents.SubscribeResponse{
		EventType:     1024,
		SubscribeType: 2048,
		ContentType:   "application/json",
		Payload:       `{"k":"v"}`,
		RequestId:     "req-9",
		Timestamp:     1758800000999,
	}}

	de := resp.ToDataEvent()
	if de.SubscribeType != 2048 {
		t.Errorf("SubscribeType = %d, want 2048", de.SubscribeType)
	}
	if de.ContentType != "application/json" {
		t.Errorf("ContentType = %q, want application/json", de.ContentType)
	}
	if string(de.Payload) != `{"k":"v"}` {
		t.Errorf("Payload = %q", de.Payload)
	}
	if de.RequestId != "req-9" {
		t.Errorf("RequestId = %q, want req-9", de.RequestId)
	}
	if de.Timestamp != 1758800000999 {
		t.Errorf("Timestamp = %d, want 1758800000999", de.Timestamp)
	}

	// The accessors and the struct conversion must agree.
	if resp.RequestId() != de.RequestId {
		t.Errorf("RequestId() = %q, ToDataEvent().RequestId = %q", resp.RequestId(), de.RequestId)
	}
	if resp.Timestamp() != de.Timestamp {
		t.Errorf("Timestamp() = %d, ToDataEvent().Timestamp = %d", resp.Timestamp(), de.Timestamp)
	}
}
