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

package eventsevents

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

var _ EventServiceClient = (*eventServiceClient)(nil)

func TestSubscribeRequestRoundTrip(t *testing.T) {
	want := &SubscribeRequest{
		SubscribeType: 7,
		Timestamp:     1700000000000,
		ContentType:   "application/json",
		Payload:       `{"userId":"u1"}`,
		Accounts:      []string{"acct-1", "acct-2"},
	}

	data, err := proto.Marshal(want)
	if err != nil {
		t.Fatalf("proto.Marshal: %v", err)
	}

	got := new(SubscribeRequest)
	if err := proto.Unmarshal(data, got); err != nil {
		t.Fatalf("proto.Unmarshal: %v", err)
	}

	if !proto.Equal(want, got) {
		t.Fatalf("round-trip mismatch:\n want %v\n  got %v", want, got)
	}
	if got.GetSubscribeType() != 7 {
		t.Errorf("SubscribeType = %d, want 7", got.GetSubscribeType())
	}
	if got.GetContentType() != "application/json" {
		t.Errorf("ContentType = %q, want %q", got.GetContentType(), "application/json")
	}
	if len(got.GetAccounts()) != 2 || got.GetAccounts()[1] != "acct-2" {
		t.Errorf("Accounts = %v, want [acct-1 acct-2]", got.GetAccounts())
	}
}

func TestSubscribeResponseDecode(t *testing.T) {
	want := &SubscribeResponse{
		EventType:     EventType_SubscribeSuccess,
		SubscribeType: 1024,
		ContentType:   "application/json",
		Payload:       `{"orderStatus":"FILLED"}`,
		RequestId:     "req-1",
		Timestamp:     1700000001000,
	}

	data, err := proto.Marshal(want)
	if err != nil {
		t.Fatalf("proto.Marshal: %v", err)
	}

	got := new(SubscribeResponse)
	if err := proto.Unmarshal(data, got); err != nil {
		t.Fatalf("proto.Unmarshal: %v", err)
	}

	if got.GetEventType() != EventType_SubscribeSuccess {
		t.Errorf("EventType = %v, want %v", got.GetEventType(), EventType_SubscribeSuccess)
	}
	if got.GetEventType() != 0 {
		t.Errorf("EventType = %d, want 0", int32(got.GetEventType()))
	}
	if got.GetRequestId() != "req-1" {
		t.Errorf("RequestId = %q, want %q", got.GetRequestId(), "req-1")
	}
}

func TestEventTypeValues(t *testing.T) {
	cases := map[EventType]int32{
		EventType_SubscribeSuccess: 0,
		EventType_Ping:             1,
		EventType_AuthError:        2,
		EventType_NumOfConnExceed:  3,
		EventType_SubscribeExpired: 4,
	}
	for et, want := range cases {
		if int32(et) != want {
			t.Errorf("%v = %d, want %d", et, int32(et), want)
		}
		if got := et.String(); got == "" {
			t.Errorf("EventType(%d).String() is empty", want)
		}
	}
}
