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

	eventsevents "github.com/shing1211/webullapi4go/gen/webull/trade/events/v1"
)

// SubscribeType is the bitmask of event categories requested in a
// SubscribeRequest. The server accepts the value as the bitwise OR of the
// individual categories.
type SubscribeType uint32

const (
	// SubscribeOrder selects order status-change events.
	SubscribeOrder SubscribeType = 1
	// SubscribePosition selects event-contract position settlement events.
	SubscribePosition SubscribeType = 2
	// SubscribeOption selects option status-change events.
	SubscribeOption SubscribeType = 4
	// SubscribeAll selects every event category and is the default.
	SubscribeAll SubscribeType = SubscribeOrder | SubscribePosition | SubscribeOption
)

// Data event kinds carried by a SubscribeResponse's event type. They are
// distinct from the control event types (success, ping, and errors) that the
// server defines as an enum, and are delivered to [Client.OnEvent].
const (
	// EventOrder identifies an order status-change event.
	EventOrder uint32 = 1024
	// EventPosition identifies an event-contract position event.
	EventPosition uint32 = 1028
	// EventOption identifies an option status-change event.
	EventOption uint32 = 1032
)

// metadataSignature is the gRPC metadata key carrying the computed signature.
// It is not part of the canonical string itself.
const metadataSignature = "x-signature"

// newSubscribeRequest builds the request in cfg. The server derives the
// timestamp from the signed body, so it is set to the current time in
// milliseconds.
func (cfg config) newSubscribeRequest(now time.Time) *eventsevents.SubscribeRequest {
	return &eventsevents.SubscribeRequest{
		SubscribeType: uint32(cfg.subscribeTypes),
		Timestamp:     now.UnixMilli(),
		Accounts:      append([]string(nil), cfg.accounts...),
	}
}
