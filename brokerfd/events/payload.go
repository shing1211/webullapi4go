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
	"time"

	eventsevents "github.com/shing1211/webullapi4go/gen/webull/brokerfd/events/v1"
)

// SubscribeType is the bitmask of Broker FD event categories requested in a
// SubscribeRequest. The server accepts the value as the bitwise OR of the
// individual categories.
type SubscribeType uint32

// DataEvent is the raw form of a data event delivered by the Broker FD event
// stream. It carries the subscribe type, content type, raw payload bytes, and
// the server-assigned request identifier and timestamp.
//
// Receive one with [Client.OnDataEvent], which delivers this whole struct.
// [Client.OnData] is also supported and is simpler, but its callback signature
// carries only the subscribe type, content type, and payload, so RequestId and
// Timestamp are dropped there.
//
// Once the JSON schema for each event category is confirmed via live probe,
// typed structs (analogous to [OrderEvent] in the trade events package) will be
// added.
type DataEvent struct {
	// SubscribeType is the event category that generated this data event.
	SubscribeType uint32
	// ContentType is the MIME type of the payload, for example "application/json".
	ContentType string
	// Payload is the raw event data. Decode it according to ContentType.
	Payload []byte
	// RequestId is the server-assigned request identifier for this event. It is
	// empty when the server does not supply one.
	RequestId string
	// Timestamp is the server timestamp in milliseconds. It is zero when the
	// server does not supply one.
	Timestamp int64
}

// SubscribeResponse wraps the generated protobuf type and provides accessor
// methods for the fields this package uses.
type SubscribeResponse struct {
	*eventsevents.SubscribeResponse
}

// EventType returns the control event type.
func (r *SubscribeResponse) EventType() eventsevents.EventType { return r.GetEventType() }

// ContentType returns the payload content type.
func (r *SubscribeResponse) ContentType() string { return r.GetContentType() }

// Payload returns the raw payload bytes.
func (r *SubscribeResponse) Payload() []byte { return []byte(r.GetPayload()) }

// RequestId returns the server-assigned request identifier.
func (r *SubscribeResponse) RequestId() string { return r.GetRequestId() }

// Timestamp returns the server timestamp in milliseconds.
func (r *SubscribeResponse) Timestamp() int64 { return r.GetTimestamp() }

// ToDataEvent converts a SubscribeResponse to a DataEvent for delivery to
// registered handlers.
func (r *SubscribeResponse) ToDataEvent() *DataEvent {
	return &DataEvent{
		SubscribeType: r.GetSubscribeType(),
		ContentType:   r.GetContentType(),
		Payload:       []byte(r.GetPayload()),
		RequestId:     r.GetRequestId(),
		Timestamp:     r.GetTimestamp(),
	}
}

// NewSubscribeRequest builds a SubscribeRequest for the Broker FD event service.
// subscribeType is the raw bitmask of categories; accounts filters events to
// those account IDs (empty subscribes to all accessible accounts).
func NewSubscribeRequest(subscribeType uint32, accounts []string) *eventsevents.SubscribeRequest {
	return &eventsevents.SubscribeRequest{
		SubscribeType: subscribeType,
		Timestamp:     time.Now().UnixMilli(),
		ContentType:   "application/json",
		Accounts:      accounts,
	}
}
