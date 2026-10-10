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
// The root module provides the shared client in client and service clients in
// data, stream, trade, events, connect, display, brokerfd, and brokerfd/events.
// Broker API HK is provided by the separate broker module. The optional webull
// package aliases selected core-client types and options; it is not an
// aggregate service facade.
//
// Shared public foundations live in pkg/errors, pkg/observability,
// pkg/resilience, pkg/transport, pkg/types, pkg/domain/money, and
// pkg/domain/order. Public financial DTOs use money.Money for required and
// response values, or *money.Money for optional and request values; decimal
// JSON strings remain the wire representation.
//
// The current request, OMS, streaming, event-telemetry, and documentation
// hardening is tagged in repository v2.1.1. This is a repository patch release,
// not a published Go-semver v2 module; v2 module publication remains deferred,
// and this module remains github.com/shing1211/webullapi4go.
package webullapi4go
