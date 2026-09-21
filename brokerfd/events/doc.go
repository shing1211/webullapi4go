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

// Package events provides a gRPC event-streaming client for the Webull Broker
// FD (US) API. It connects to the Broker FD event service, signs each
// Subscribe request with HMAC-SHA256 over gRPC metadata, and delivers pushed
// events to registered handlers.
//
// The package mirrors the architecture of the trade-events client
// (github.com/shing1211/webullapi4go/events) and is safe for concurrent use.
//
// # Provisional status
//
// This package is marked provisional because two key aspects of the Broker FD
// event service require live US sandbox confirmation before they can be
// finalized:
//
//  1. SubscribeType bitmask values: The plan references 13 event categories
//     but the exact integer values (for example which bit maps to which
//     category) are unconfirmed. Callers pass the raw uint32 until constants
//     are added after a live probe.
//
//  2. Data event payload shapes: The JSON schema for each data event category
//     is unknown. The package delivers raw []byte payloads via Client.OnData
//     rather than attempting to decode into typed structs. Once the schema is
//     confirmed, typed structs (analogous to OrderEvent in the trade events
//     package) will be added.
//
// TODO(v08): confirm SubscribeType bitmask values against live US sandbox.
// TODO(v08): confirm data event JSON schemas via live US sandbox capture.
//
// # Construct
//
// Build a Client from the public client and start the run loop:
//
//	func Example() {
//		cl, err := client.New(client.WithEnv())
//		if err != nil { /* ... */ }
//
//		fd := brokerfd.New(cl)
//		ec := events.New(cl, events.WithAccounts([]string{accountID}))
//
//		ec.OnConnect(func() { log.Println("connected") })
//		ec.OnData(func(st uint32, ct string, payload []byte) {
//			log.Printf("event st=%d ct=%s size=%d", st, ct, len(payload))
//		})
//		ec.OnError(func(err error) { log.Printf("error: %v", err) })
//
//		if err := ec.Run(context.Background()); err != nil {
//			log.Fatalf("events run ended: %v", err)
//		}
//	}
//
// # gRPC endpoint
//
// The gRPC endpoint is resolved from the public client's Endpoints().GRPC field
// unless overridden with events.WithGRPCEndpoint. For the US Broker FD service
// the host is events-api.webull.com (production) or
// events-api.sandbox.webull.hk (sandbox).
package events
