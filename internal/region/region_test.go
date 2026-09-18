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

package region_test

import (
	"testing"

	"github.com/shing1211/webullapi4go/internal/region"
)

func TestDefaultEndpointsIsHKProduction(t *testing.T) {
	got := region.DefaultEndpoints()
	if got.HTTP != "https://api.webull.hk" {
		t.Fatalf("DefaultEndpoints().HTTP = %q, want %q", got.HTTP, "https://api.webull.hk")
	}
	if got != (region.Endpoints{
		HTTP:          region.HKProductionHTTP,
		MQTT:          region.HKProductionMQTT,
		MQTTWebSocket: region.HKProductionMQTTWebSocket,
		GRPC:          region.HKProductionGRPC,
	}) {
		t.Fatalf("DefaultEndpoints() = %+v, want the documented HK production endpoints", got)
	}
}

func TestEndpointsFor(t *testing.T) {
	tests := []struct {
		name string
		r    region.Region
		env  region.Environment
		want region.Endpoints
	}{
		{
			name: "hk production",
			r:    region.HK,
			env:  region.Production,
			want: region.Endpoints{
				HTTP:          "https://api.webull.hk",
				MQTT:          "data-api.webull.hk:1883",
				MQTTWebSocket: "wss://data-api.webull.hk:8883/mqtt",
				GRPC:          "events-api.webull.hk",
			},
		},
		{
			name: "hk sandbox",
			r:    region.HK,
			env:  region.Sandbox,
			want: region.Endpoints{
				HTTP:          "https://api.sandbox.webull.hk",
				MQTT:          "data-api.sandbox.webull.hk:1883",
				MQTTWebSocket: "wss://data-api.sandbox.webull.hk:8883/mqtt",
				GRPC:          "events-api.sandbox.webull.hk",
			},
		},
		{
			name: "us production",
			r:    region.US,
			env:  region.Production,
			want: region.Endpoints{
				HTTP:          "https://api.webull.com",
				MQTT:          "data-api.webull.com:1883",
				MQTTWebSocket: "wss://data-api.webull.com:8883/mqtt",
				GRPC:          "events-api.webull.com",
			},
		},
		{
			name: "jp sandbox",
			r:    region.JP,
			env:  region.Sandbox,
			want: region.Endpoints{
				HTTP:          "https://api.sandbox.webull.co.jp",
				MQTT:          "data-api.sandbox.webull.co.jp:1883",
				MQTTWebSocket: "wss://data-api.sandbox.webull.co.jp:8883/mqtt",
				GRPC:          "events-api.sandbox.webull.co.jp",
			},
		},
		{
			name: "unknown region falls back to hk",
			r:    region.Region("xx"),
			env:  region.Production,
			want: region.EndpointsFor(region.HK, region.Production),
		},
		{
			name: "unknown environment falls back to production",
			r:    region.HK,
			env:  region.Environment("staging"),
			want: region.EndpointsFor(region.HK, region.Production),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := region.EndpointsFor(tt.r, tt.env); got != tt.want {
				t.Fatalf("EndpointsFor(%q, %q) = %+v, want %+v", tt.r, tt.env, got, tt.want)
			}
		})
	}
}

func TestDomain(t *testing.T) {
	if d, ok := region.Domain(region.HK); !ok || d != "webull.hk" {
		t.Fatalf("Domain(HK) = %q, %v; want %q, true", d, ok, "webull.hk")
	}
	if _, ok := region.Domain(region.Region("xx")); ok {
		t.Fatalf("Domain(unknown) reports ok")
	}
}

func TestRegionIsValid(t *testing.T) {
	if !region.HK.IsValid() || !region.EU.IsValid() {
		t.Fatalf("known region reported invalid")
	}
	if region.Region("xx").IsValid() {
		t.Fatalf("unknown region reported valid")
	}
}
