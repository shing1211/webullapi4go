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

package stream

import "fmt"

// State describes the lifecycle state of the MQTT connection.
type State int

const (
	// StateDisconnected means no connection has been attempted or the previous
	// connection was closed cleanly.
	StateDisconnected State = iota
	// StateConnecting means a connection attempt is in progress.
	StateConnecting
	// StateConnected means the MQTT connection is live.
	StateConnected
	// StateReconnecting means the connection was lost and the client is
	// attempting to re-establish it automatically.
	StateReconnecting
	// StateDegraded means the connection is live but message throughput has
	// dropped below the health threshold (message-age watchdog triggered).
	StateDegraded
	// StateClosed means [Client.Close] was called; the client is permanent.
	StateClosed
)

func (s State) String() string {
	switch s {
	case StateDisconnected:
		return "disconnected"
	case StateConnecting:
		return "connecting"
	case StateConnected:
		return "connected"
	case StateReconnecting:
		return "reconnecting"
	case StateDegraded:
		return "degraded"
	case StateClosed:
		return "closed"
	default:
		return fmt.Sprintf("state(%d)", int(s))
	}
}
