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

package client

import "github.com/shing1211/webullapi4go/internal/region"

// Region identifies a Webull OpenAPI deployment region. It is a public mirror
// of the internal region identifier so that external callers never need to
// import an internal package.
type Region string

const (
	// HK is the Hong Kong region (the SDK's primary target).
	HK Region = "hk"
	// US is the United States region.
	US Region = "us"
	// JP is the Japan region.
	JP Region = "jp"
	// SG is the Singapore region.
	SG Region = "sg"
	// TH is the Thailand region.
	TH Region = "th"
	// AU is the Australia region.
	AU Region = "au"
	// MY is the Malaysia region.
	MY Region = "my"
	// UK is the United Kingdom region.
	UK Region = "uk"
	// BR is the Brazil region.
	BR Region = "br"
	// MX is the Mexico region.
	MX Region = "mx"
	// ZA is the South Africa region.
	ZA Region = "za"
	// EU is the European Union region.
	EU Region = "eu"
)

// String returns the lowercase region identifier.
func (r Region) String() string { return string(r) }

// IsValid reports whether r is a region known to the SDK.
func (r Region) IsValid() bool {
	_, ok := region.Domain(region.Region(r))
	return ok
}

// Environment selects between the production and sandbox deployments.
type Environment string

const (
	// Production is the live Webull OpenAPI environment.
	Production Environment = "production"
	// Sandbox is the Webull test environment.
	Sandbox Environment = "sandbox"
)

// String returns the environment name.
func (e Environment) String() string { return string(e) }

// IsValid reports whether e is a known environment.
func (e Environment) IsValid() bool {
	return e == Production || e == Sandbox
}

// Endpoints holds the resolved service addresses for one region/environment
// pair. It is a public mirror of the internal endpoint set.
type Endpoints struct {
	// HTTP is the REST API base URL used for trading and market data.
	HTTP string
	// MQTT is the MQTT broker address (host:port) used for market-data
	// streaming.
	MQTT string
	// MQTTWebSocket is the MQTT-over-WebSocket URL used for market-data
	// streaming.
	MQTTWebSocket string
	// GRPC is the host of the gRPC trading-event stream.
	GRPC string
	// BrokerHTTP is the Broker API REST endpoint.
	BrokerHTTP string
}

// EndpointsFor returns the service endpoints for the given region and
// environment. An unknown region falls back to Hong Kong and an unknown
// environment falls back to production; call [Region.IsValid] and
// [Environment.IsValid] first to detect this.
func EndpointsFor(r Region, env Environment) Endpoints {
	return endpointsFromInternal(region.EndpointsFor(region.Region(r), region.Environment(env)))
}

// DefaultEndpoints returns the production Hong Kong endpoints.
func DefaultEndpoints() Endpoints {
	return endpointsFromInternal(region.DefaultEndpoints())
}

// endpointsFromInternal converts the internal endpoint set into the public
// mirror. It keeps internal types out of the public API.
func endpointsFromInternal(e region.Endpoints) Endpoints {
	return Endpoints{
		HTTP:          e.HTTP,
		MQTT:          e.MQTT,
		MQTTWebSocket: e.MQTTWebSocket,
		GRPC:          e.GRPC,
		BrokerHTTP:    e.BrokerHTTP,
	}
}
