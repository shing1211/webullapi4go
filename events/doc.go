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

// Package events provides a client for Webull's gRPC trade-event stream.
//
// The event service is a server-streaming gRPC RPC (grpc.trade.event.
// EventService/Subscribe) that pushes order, position, and option events over a
// persistent TLS connection. Unlike the REST API, it authenticates each
// Subscribe call with an HMAC-SHA256 signature over the serialized request:
// the five x-signature headers plus x-signature are sent as gRPC metadata and
// no host header participates in the canonical string.
//
// A client is built from the public [client.Client], which supplies the
// credentials, region, and the resolved gRPC endpoint:
//
//	cl, err := client.New(
//		client.WithAppKey(key),
//		client.WithAppSecret(secret),
//		client.WithEnvironment(client.Sandbox),
//	)
//	if err != nil {
//		return err
//	}
//
//	ev, err := events.New(cl, events.WithAccounts([]string{accountID}))
//	if err != nil {
//		return err
//	}
//	defer ev.Close()
//
//	ev.OnConnect(func() { log.Print("subscribed") })
//	ev.OnError(func(err error) { log.Print(err) })
//	ev.OnOrder(func(o *events.OrderEvent) {
//		log.Printf("order %s is %s", o.OrderID, o.OrderStatus)
//	})
//
//	if err := ev.Run(ctx); err != nil {
//		return err
//	}
//
// [Client.Run] blocks, pumping the stream until the context is cancelled or the
// stream ends. By default it reconnects and re-subscribes with exponential
// backoff and jitter after a clean end or a transient failure, and stops on
// terminal errors (authentication, permission, account, and configuration
// failures) or when the attempt budget set with [WithMaxReconnectAttempts] is
// exhausted; disable it with [WithAutoReconnect](false). Each (re)connect signs
// and sends a fresh Subscribe request.
//
// # Event dispatch
//
// Incoming messages are routed by their event type:
//
//   - SubscribeSuccess fires the [Client.OnConnect] handlers.
//   - Ping fires the [Client.OnPing] handlers (heartbeat).
//   - AuthError, NumOfConnExceed, and SubscribeExpired fire [Client.OnError]
//     with a typed error and end Run.
//   - Every other event type is a data event: it is delivered raw to
//     [Client.OnEvent] with its kind ([EventOrder], [EventPosition], or
//     [EventOption]) and, when the content type is application/json, decoded
//     into [OrderEvent], [PositionEvent], or [OptionEvent] and delivered to the
//     matching typed handler ([Client.OnOrder], [Client.OnPosition], or
//     [Client.OnOption]). A decoding failure goes to [Client.OnError] without
//     ending the stream.
package events
