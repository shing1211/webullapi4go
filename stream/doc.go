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

// Package stream provides Webull market-data streaming over MQTT.
//
// A streaming client is built from the public [client.Client], which supplies
// the credentials, region, and the resolved MQTT endpoints:
//
//	cl, err := client.New(
//		client.WithAppKey(key),
//		client.WithAppSecret(secret),
//		client.WithSandbox(),
//	)
//	if err != nil {
//		return err
//	}
//	if _, err := cl.EnsureToken(ctx); err != nil {
//		return err
//	}
//
//	s, err := stream.New(cl)
//	if err != nil {
//		return err
//	}
//	defer s.Close()
//
//	s.OnQuote(func(q *marketdatav1.Quote) { ... })
//	if err := s.Connect(ctx); err != nil {
//		return err
//	}
//	if err := s.Subscribe(ctx, stream.SubscribeRequest{
//		Symbols:  []string{"AAPL"},
//		Category: stream.CategoryUSStock,
//		SubTypes: []stream.SubType{stream.SubTypeQuote, stream.SubTypeSnapshot},
//	}); err != nil {
//		return err
//	}
//
// # Protocol
//
// The MQTT connection carries only pushes: subscribe and unsubscribe are HTTP
// calls made with [Client.Subscribe] and [Client.Unsubscribe]. The MQTT
// CONNECT packet uses a freshly generated session id as the client id, the App
// Key as the user name, and an arbitrary password.
//
// Incoming messages are routed by topic and decoded before the registered
// handler is invoked:
//
//   - [TopicQuote] is decoded into a marketdata Quote and delivered to
//     [Client.OnQuote]
//   - [TopicSnapshot] is decoded into a Snapshot and delivered to
//     [Client.OnSnapshot]
//   - [TopicTick] is decoded into a Tick and delivered to [Client.OnTick]
//   - [TopicNotice] is delivered to [Client.OnNotice] as raw JSON
//   - [TopicEcho] is a heartbeat and is ignored
//
// # Reconnection
//
// Webull does not restore subscriptions after a connection is lost, and the
// MQTT broker never restores them either. With [WithAutoReconnect], the client
// detects the lost connection, reconnects, and then re-issues the HTTP
// subscribe calls for every active subscription before the [Client.OnConnect]
// handlers run. Re-subscription is idempotent: each active subscription is
// re-issued exactly once and a symbol subscribed twice is only restored once.
// [Client.OnDisconnect] fires when a live connection is lost; [Client.OnConnect]
// fires after the initial connection and after every successful reconnect;
// [Client.Reconnecting] reports an in-progress reconnect attempt.
//
// # Limits
//
// Webull applies the following connection rules, which the SDK reflects where
// it can:
//
//   - An App Key supports at most five concurrent MQTT connections. Exceeding
//     the limit fails the connection with Webull error code 105. When that
//     happens [Client.Connect] returns an error explaining the limit.
//   - A new connection that reuses an existing session id disconnects the
//     previous connection. [New] therefore generates a unique session id by
//     default; only set [WithSessionID] when exclusivity is guaranteed.
//   - After a disconnect the server retains the connection state for about one
//     minute. When the limit is reached, wait roughly a minute before
//     reconnecting rather than retrying immediately.
//   - The server pushes at most three messages per second per connection. The
//     client does not need to throttle; this only bounds the inbound rate.
package stream
