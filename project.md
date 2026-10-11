---
type: project
---

# webullapi4go — Project Overview

## Identity

- **Module**: `github.com/shing1211/webullapi4go`
- **Type**: Go SDK (idiomatic, typed client)
- **Latest tag**: `v2.1.39` (docs); installable `v1.1.1` line
- **Go version**: 1.26+ (no cgo)
- **License**: Apache-2.0

## What this project is

An idiomatic Go SDK for the [Webull OpenAPI](https://developer.webull.hk/apis/docs/).
Provides typed HTTP, MQTT, and gRPC clients for Webull's trading, market data, and broker services.

## What this project is NOT

- Not a published v2 module — import path stays on v1 by decision
- Not live-verified end-to-end — HK sandbox tested, US production access unavailable

## Feature areas

| Area | Status | Notes |
|---|---|---|
| Auth (HMAC-SHA1 REST, token flow) | ✅ Implemented | Live-verified HK |
| Market Data HTTP | ✅ Implemented | Offline + HK live-verified |
| Market Data MQTT (QUOTE/SNAPSHOT/TICK) | ✅ Implemented | Offline-tested |
| Trading HTTP + OMS | ✅ Implemented | Offline + HK live-verified |
| Trading gRPC events | ✅ Implemented | Offline-tested |
| Broker HK | ✅ Implemented | Live blocked (no HK creds) |
| Broker FD US | ✅ Implemented | Live blocked (no US creds) |
| Display Solution | ✅ Implemented | Live blocked (host/entitlement) |
| Connect OAuth | ✅ Implemented | US live unavailable |

## Architecture

- `broker/` — virtual account broker API (HK)
- `brokerfd/` — broker FD API (US)
- `pkg/observability/` — OpenTelemetry tracing/metrics (otel, no-op by default)
- `pkg/resilience/` — rate limiting, circuit breaker, retry helpers
- `pkg/transport/mqtt/` — MQTT client wrapper
- `stream/` — streaming client abstraction
- `trade/` — trading HTTP client
- `display/` — Display Solution client
- `data/` — market data HTTP client
- `gen/` — generated protobuf/grpc types

## Key constraints

- No cgo — pure Go
- No published v2 import path
- Live verification blocked by credential availability

## OpenSpec

- Config: `openspec/config.yaml`
- Changes: `openspec/changes/`
- Specs: `openspec/specs/`
- ADRs: `docs/adr/`
