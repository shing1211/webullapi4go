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

package marketdatav1

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// DecodeSnapshot unmarshals a protobuf payload received on the `snapshot` MQTT
// topic into a Snapshot.
func DecodeSnapshot(data []byte) (*Snapshot, error) {
	msg := new(Snapshot)
	if err := proto.Unmarshal(data, msg); err != nil {
		return nil, fmt.Errorf("webull/marketdata: decode snapshot: %w", err)
	}
	return msg, nil
}

// DecodeQuote unmarshals a protobuf payload received on the `quote` MQTT topic
// into a Quote.
func DecodeQuote(data []byte) (*Quote, error) {
	msg := new(Quote)
	if err := proto.Unmarshal(data, msg); err != nil {
		return nil, fmt.Errorf("webull/marketdata: decode quote: %w", err)
	}
	return msg, nil
}

// DecodeTick unmarshals a protobuf payload received on the `tick` MQTT topic
// into a Tick.
func DecodeTick(data []byte) (*Tick, error) {
	msg := new(Tick)
	if err := proto.Unmarshal(data, msg); err != nil {
		return nil, fmt.Errorf("webull/marketdata: decode tick: %w", err)
	}
	return msg, nil
}
