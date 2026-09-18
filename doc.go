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

// Package webullapi4go is an idiomatic Go SDK for the Webull OpenAPI.
//
// It wraps the Webull HTTP and MQTT services in typed Go, starting with the
// Hong Kong region. The v0.1 surface covers authentication, the core REST
// client, Market Data over HTTP, and Market Data streaming over MQTT.
//
// Package layout:
//
//   - client: the public client, configuration, options, signing, and transport
//     entry point.
//   - data: the Market Data API client (built on client).
//   - stream: Market Data streaming over MQTT (built on client).
//   - gen/webull/marketdata/v1: generated protobuf types for streaming pushes.
//   - pkg/types: shared domain types that are safe for external use.
//   - internal/auth: request signing and token lifecycle.
//   - internal/region: deployment regions and their service endpoints.
//   - internal/errs: the SDK's typed error model.
//   - internal/transport: the thin HTTP executor.
//   - internal/mqtt: the low-level MQTT transport used by the stream package.
//   - internal/resilience: retry, rate limiting, and circuit breaking.
//
// Trading over HTTP, gRPC trade events, the Display Solution API, and the
// Broker API are deferred to later releases and are not part of v0.1.
package webullapi4go
