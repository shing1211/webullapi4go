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

// Package mqtt is the low-level MQTT transport used by the public stream
// package to receive Webull market-data pushes.
//
// It wraps github.com/eclipse/paho.mqtt.golang behind a small, dependency-free
// surface: [Client.Connect] establishes a TCP or WebSocket connection,
// [Client.Disconnect] tears it down, and message and connection-state callbacks
// are registered with the Set*Handler methods. No paho type appears in an
// exported signature, so the MQTT library can be swapped without changing
// public API.
//
// The wrapper deliberately performs no topic subscription of its own. Webull
// manages subscriptions through the HTTP streaming API (see the public stream
// package), and the broker pushes messages to the connected session; a
// per-message callback is installed with [Client.SetMessageHandler].
package mqtt
