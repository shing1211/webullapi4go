# Production Hardening Plan and Decision Record

Last updated: 2026-09-25

Status: **Latest release `v2.1.0`; current hardening is implemented and offline-tested under Unreleased, not released or newly live-verified.**

## Architecture baseline

The current architecture preserves the original public service packages. Moving every service under `pkg/` was considered and is closed.

| Layer | Canonical packages | Responsibility |
|---|---|---|
| Core | `client/` | Configuration, signing, tokens, HTTP transport, resilience, and request pipeline |
| Services | `data/`, `stream/`, `trade/`, `events/`, `connect/`, `display/`, `brokerfd/`, `brokerfd/events/` | Market data, streaming, trading, events, OAuth, and client surfaces |
| Broker HK | `broker/` | Separate Go module for Broker API HK |
| Shared foundations | `pkg/errors/`, `pkg/observability/`, `pkg/resilience/`, `pkg/transport/`, `pkg/types/`, `pkg/domain/money/`, `pkg/domain/order/` | Reusable public primitives and domain types; not relocated services |
| Convenience | `webull/` | Thin type and constructor aliases for the core client; it is not an aggregate service facade |
| Implementation details | `internal/` | Authentication, compatibility aliases, legacy tests, and other non-public code |

Public DTOs that model decimal financial values use `money.Money` (required/response fields) or `*money.Money` (optional/request fields). Raw `decimal.Decimal` is not a public DTO migration target. JSON remains a decimal string on the wire.

## Verification vocabulary

| State | Meaning |
|---|---|
| Implemented | Public code exists in the current tree |
| Offline-tested | Credential-free unit or local-server tests pass |
| Live-verified | Exercised successfully against an available Webull sandbox or production surface |
| Blocked | Cannot be live-verified with the currently available credentials, host, or entitlement |
| Released | Present in a tagged release; Unreleased work is not included |

An implemented endpoint is not automatically live-verified. Offline tests do not prove that an upstream host, entitlement, symbol set, or response schema is available.

## Current resolved decisions

| # | Decision | Resolution |
|---|---|---|
| T1 | Tracing and metrics | Direct `go.opentelemetry.io/otel` API integration with no-op defaults, `WithTracerProvider`, `WithMeterProvider`, `WithPropagator`, and `WithLogger`; no `otelhttp` dependency |
| T2 | Resilience defaults | Keep the released v1 behavior; rate limiting, circuit breaking, and full-jitter presets remain opt-in |
| T3 | Package layout | Preserve root service packages; use `pkg/` only for shared foundations |
| T4 | Numeric DTOs | Use `money.Money` and `*money.Money`; do not migrate public DTOs to raw `decimal.Decimal` |
| T5 | OMS | Keep the SDK-side, concurrency-safe order state machine and add status reconciliation from events and snapshots |
| T6 | HMAC signing | REST remains HMAC-SHA1 with an MD5 body digest; event gRPC remains HMAC-SHA256; no protocol change |
| T7 | Go version | `go.mod` requires Go 1.26; CI validates supported 1.26 toolchains |
| T8 | Streaming allocations | Do not use `sync.Pool` in the streaming path; ownership and channel lifecycle remain explicit |
| T9 | Live verification | Record credential and entitlement blockers instead of treating implementation as proof of live behavior |

## Superseded and closed decisions

| Earlier decision | Status | Current resolution |
|---|---|---|
| Relocate `client`, `data`, `trade`, `stream`, `events`, and the other services under a full `pkg/` hierarchy | **Superseded; closed** | Root service packages remain canonical; only shared foundations live under `pkg/` |
| Turn root service packages into deprecated forwarding shims | **Superseded; closed** | Root packages contain the implementations and are not deprecated. Limited deprecated aliases under `internal/` remain only for old internal imports |
| Migrate all DTO decimals to raw `decimal.Decimal` | **Superseded; closed** | Public financial values use `money.Money`; JSON decimal strings are preserved |
| Add `sync.Pool` unconditionally, or as an assumed streaming optimization | **Closed without implementation** | The stream path has no `sync.Pool`; a future pool requires a separate benchmark-backed proposal |
| Build a full aggregate `webull` service facade | **Superseded; closed** | `webull/` remains a thin alias/convenience package; service clients remain in their root packages |
| Wrap REST with `otelhttp` | **Superseded; closed** | The SDK emits direct OpenTelemetry client spans and metrics through the OTel API |

## Workstream status

| Workstream | Released baseline | Unreleased follow-up | Verification |
|---|---|---|---|
| DevOps and tests | Build/test/lint targets, security checks, multi-OS CI, dependabot, leak and fuzz tests | Current changes | Offline suite passes 2026-09-25 |
| Shared foundations | Public errors, transport, resilience, observability, `money.Money`, and order domain packages | Additional money and order-state regression coverage | Offline-tested |
| REST pipeline | Interceptors, hooks, resilience, clock correction, tracing, metrics, and logging | `Do`, `DoBroker`, and `DoStream` share one attempt pipeline; one-based attempt telemetry; stable correlation IDs; W3C propagation; response status and clock-offset parity | Offline-tested; not newly live-verified |
| Trading OMS | Placement tracking and terminal preflight | Account-scoped tracking, status reconciliation, stable auto-generated client IDs, and success-only action transitions | Offline-tested; not newly live-verified |
| MQTT streaming | State machine, reconnect/resubscribe, callbacks, and bounded channel policies | Terminal close, idempotent channel cancellation, blocked-dispatch cancellation, serialized resubscription, data-only health recovery, and MQTT dispatch tracing | Offline-tested; not newly live-verified |
| Event telemetry | Typed event streams and reconnect | Per-attempt gRPC spans, attempt/duration metrics, structured logs, correlation metadata, and trace propagation for Trading and Broker FD events | Offline-tested with local gRPC servers; not newly live-verified |
| Documentation | Guides, generated API pages, and reconciliation | Architecture/status reconciliation, OMS, streaming health/channels, typed errors, money, and OTel setup | `mkdocs build --strict` is the release gate |

## Unreleased scope

Current work is intentionally recorded as Unreleased until a tag is created.

1. Request-pipeline parity and observability hardening.
2. OMS status reconciliation and account-scoped order tracking.
3. MQTT/channel shutdown, reconnect serialization, and health correctness.
4. Trading and Broker FD event telemetry.
5. Updated account-monitor/order/stream examples and documentation.

No item in this section is represented as released merely because it exists in the working tree.

## Live-verification boundaries

| Surface | Current state |
|---|---|
| Core auth, selected HK market data, accounts, read-only assets, MQTT streaming, and Trading events | Previously exercised against available HK sandbox paths; not all endpoints or current hardening changes are live-verified |
| Footprint | Implemented and offline-tested; live request blocked by entitlement (`403`) |
| Options contracts and multi-leg orders | Implemented and offline-tested; HK sandbox may return `417`, and US live behavior is unavailable |
| US-only crypto, fund, screener, Broker FD, and related data | Implemented and offline-tested where covered; live verification blocked without US sandbox credentials |
| Broker API HK | Implemented and offline-tested in its own module; HK sandbox returns `401 ROUTE_NOT_PERMITTED` because scope is missing |
| Display Solution | Implemented and offline-tested; HK host returns `403` before endpoint-specific behavior can be verified |
| SSE news | Implemented; HK upstream currently returns `504` |

The generated [SDK ↔ API reconciliation](docs/reconciliation.md) remains authoritative for endpoint counts and path states. Its 2026-09-22 snapshot reports 209 implemented endpoints and 0 documented-only gaps, while also reporting summary-only and unresolved paths; do not summarize that report as zero discrepancies.

## Verification gates

```sh
go build ./...
go vet ./...
gofmt -l .                 # must print nothing
go test ./...
go test -race -count=1 ./...
golangci-lint run ./...
mkdocs build --strict
```

The nested `broker/` module is tested separately with `go test ./...` from that directory.

## Out of scope

- Reintroducing a full `pkg/` service relocation.
- Deprecating the root service packages.
- Raw `decimal.Decimal` DTO migration.
- Unconditional or speculative `sync.Pool` use.
- Editing generated files under `docs/webull-api/`.
- Editing internal run artifacts under `docs/runs/`.
- Claiming blocked or newly hardened surfaces are live-verified without the required access.

## Reference

- Repository guide: `AGENTS.md`
- Build/test/lint: `Makefile`
- Current status: `IMPLEMENTATION_STATUS.md`
- Changelog: `CHANGELOG.md`
- Generated endpoint reconciliation: `docs/reconciliation.md`
