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

import "strings"

// MQTT topics published by the Webull streaming broker.
const (
	// TopicQuote carries real-time order-book depth as protobuf.
	TopicQuote = "quote"
	// TopicSnapshot carries market snapshots as protobuf.
	TopicSnapshot = "snapshot"
	// TopicTick carries tick-by-tick trades as protobuf.
	TopicTick = "tick"
	// TopicNotice carries server notifications as JSON.
	TopicNotice = "notice"
	// TopicEcho is a heartbeat with a null payload and is ignored.
	TopicEcho = "echo"
)

// normalizeTopic reduces a broker topic to its routing key. Topics are
// case-insensitive and may carry a hierarchical prefix, so the comparison uses
// the lower-cased final path segment.
func normalizeTopic(topic string) string {
	t := strings.ToLower(strings.TrimSpace(topic))
	if i := strings.LastIndexByte(t, '/'); i >= 0 {
		t = t[i+1:]
	}
	return t
}
