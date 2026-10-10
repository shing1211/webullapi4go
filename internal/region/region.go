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

// Package region describes Webull OpenAPI deployment regions and the service
// endpoints (HTTP, MQTT, gRPC) for each environment.
//
// Hosts are region specific: every region has its own root domain and the API
// hosts are derived from it, for example Hong Kong uses api.webull.hk and the
// United States uses api.webull.com. [EndpointsFor] resolves the full set of
// hosts for a region and environment; [DefaultEndpoints] returns the production
// Hong Kong endpoints, which are the SDK's primary target.
//
// Only the Hong Kong hosts are published explicitly in Webull's HK
// documentation. The subdomain pattern for the other regions
// (data-api.<domain>, events-api.<domain>, and their sandbox variants) is
// inferred by analogy from that documentation and should be verified against
// each region before use in production.
package region

import "strings"

// Region identifies a Webull OpenAPI region.
type Region string

const (
	HK Region = "hk"
	US Region = "us"
	JP Region = "jp"
	SG Region = "sg"
	TH Region = "th"
	AU Region = "au"
	MY Region = "my"
	UK Region = "uk"
	BR Region = "br"
	MX Region = "mx"
	ZA Region = "za"
	EU Region = "eu"
)

// String returns the lowercase region identifier.
func (r Region) String() string { return string(r) }

// IsValid reports whether r is a known region.
func (r Region) IsValid() bool {
	_, ok := domainByRegion[r]
	return ok
}

// ParseRegion parses a case-insensitive region identifier. The boolean result
// reports whether the identifier names a known region.
func ParseRegion(s string) (Region, bool) {
	r := Region(strings.ToLower(strings.TrimSpace(s)))
	return r, r.IsValid()
}

// domainByRegion maps each region to its Webull root domain. It is derived from
// the developer documentation domains by stripping the "developer." prefix; the
// API host for a region is "api." followed by this domain.
var domainByRegion = map[Region]string{
	HK: "webull.hk",
	US: "webull.com",
	JP: "webull.co.jp",
	SG: "webull.com.sg",
	TH: "webull.co.th",
	AU: "webull.com.au",
	MY: "webull.com.my",
	UK: "webull-uk.com",
	BR: "webull.com.br",
	MX: "webull.com.mx",
	ZA: "webull.co.za",
	EU: "webull.eu",
}

// Domain returns the Webull root domain for r and reports whether r is known.
// For example, HK returns "webull.hk".
func Domain(r Region) (string, bool) {
	d, ok := domainByRegion[r]
	return d, ok
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

// Endpoints holds the service addresses for one region/environment pair.
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

// Documented Hong Kong service hosts. These are the values returned by
// EndpointsFor(HK, Production) and EndpointsFor(HK, Sandbox) respectively and
// are exposed for readability and tests.
const (
	// Hong Kong production.
	HKProductionHTTP          = "https://api.webull.hk"
	HKProductionGRPC          = "events-api.webull.hk"
	HKProductionBrokerHTTP    = "https://broker-api.webull.hk"
	HKProductionMQTT          = "data-api.webull.hk:1883"
	HKProductionMQTTWebSocket = "wss://data-api.webull.hk:8883/mqtt"

	// Hong Kong sandbox.
	HKSandboxHTTP          = "https://api.sandbox.webull.hk"
	HKSandboxGRPC          = "events-api.sandbox.webull.hk"
	HKSandboxBrokerHTTP    = "https://broker-api.sandbox.webull.hk"
	HKSandboxMQTT          = "data-api.sandbox.webull.hk:1883"
	HKSandboxMQTTWebSocket = "wss://data-api.sandbox.webull.hk:8883/mqtt"
)

// Default MQTT ports.
const (
	// DefaultMQTTPort is the plain MQTT broker port.
	DefaultMQTTPort = 1883
	// DefaultMQTTWebSocketPort is the MQTT-over-WebSocket port.
	DefaultMQTTWebSocketPort = 8883
)

// EndpointsFor returns the service endpoints for the given region and
// environment.
//
// Hosts are region specific. An unknown region cannot be mapped to a host
// safely, so it falls back to the primary region, Hong Kong, rather than
// producing a malformed address; an unknown environment falls back to
// production. Callers that need to detect this should validate the inputs with
// [Region.IsValid] and [Environment.IsValid] first.
func EndpointsFor(r Region, env Environment) Endpoints {
	domain, ok := domainByRegion[r]
	if !ok {
		domain = domainByRegion[HK]
	}
	if env != Production && env != Sandbox {
		env = Production
	}
	return endpointsFromDomain(domain, env)
}

// DefaultEndpoints returns the production Hong Kong endpoints, the SDK's
// primary target.
func DefaultEndpoints() Endpoints {
	return EndpointsFor(HK, Production)
}

// endpointsFromDomain derives the service hosts for a root domain and
// environment. The sandbox environment adds a "sandbox." label below each
// service subdomain.
func endpointsFromDomain(domain string, env Environment) Endpoints {
	prefix := ""
	if env == Sandbox {
		prefix = "sandbox."
	}
	return Endpoints{
		HTTP:          "https://api." + prefix + domain,
		MQTT:          "data-api." + prefix + domain + ":1883",
		MQTTWebSocket: "wss://data-api." + prefix + domain + ":8883/mqtt",
		GRPC:          "events-api." + prefix + domain,
		BrokerHTTP:    "https://broker-api." + prefix + domain,
	}
}
