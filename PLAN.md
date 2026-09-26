# Production Hardening Plan and Decision Record

Last updated: 2026-09-25

Status: **Current request, OMS, streaming, event-telemetry, and documentation
hardening is tagged in repository `v2.1.3` (2026-09-26) and offline-tested. It
is a repository patch release, not a published Go-semver v2 module, and was
not newly live-verified.** The root module path is unchanged and stays on the v1
import path by decision, with no `/v2` migration planned.

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
| Repository-tagged | Included in the authorized repository Git tag `v2.1.1`; this is not a published Go-semver v2 module |
| Released | Present in a published module release; the unchanged root path means `v2.x` Git tags are not published Go-semver v2 modules |

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

| Workstream | Published/historical baseline | v2.1.1 repository-tagged follow-up | Verification |
|---|---|---|---|
| DevOps and tests | Build/test/lint targets, security checks, multi-OS CI, dependabot, leak and fuzz tests | Direct public-primitive tests, deterministic cancellation/leak coverage, and reproducible per-module coverage measurement | Module-aware race/vet evidence recorded 2026-09-26; root 73.6% and broker 80.8% measurements are not guarantees |
| Error contracts | Public typed categories and compatibility sentinels | Category matching remains stable; `NewSentinel` provides identity-specific semantics; HTTP 417, MQTT, and event terminal mappings are tested | Offline-tested |
| Shared foundations | Public errors, transport, resilience, observability, `money.Money`, and order domain packages | Provider-order-independent metrics, safe error text, cancellation-safe rate limiting, and direct public-package tests | Offline-tested |
| REST pipeline | Interceptors, hooks, resilience, clock correction, tracing, metrics, and logging | `Do`, `DoBroker`, and `DoStream` share one attempt pipeline; one-based attempt telemetry; stable correlation IDs; W3C propagation; response status and clock-offset parity | Offline-tested; not newly live-verified |
| Trading OMS | Placement tracking and terminal preflight | Account-scoped tracking, status reconciliation, stable auto-generated client IDs, and success-only action transitions | Offline-tested; not newly live-verified |
| MQTT streaming | State machine, reconnect/resubscribe, callbacks, and bounded channel policies | Compare-and-swap state, deterministic data-only recovery, terminal close, idempotent channel cancellation, blocked-dispatch cancellation, serialized replay, and explicit synchronous head-of-line dispatch | Offline-tested; not newly live-verified |
| Event telemetry | Typed event streams and reconnect | Per-attempt gRPC spans/metrics/logs, sanitized failure text, cancellation telemetry, correlation propagation, and Trading `Close` cancellation of all active runs | Offline-tested with local gRPC servers; not newly live-verified |
| Documentation | Guides, generated API pages, and reconciliation | Architecture/status synchronization, errors/417, streaming lifecycle, Broker FD raw-event limits, testing boundaries, and exact telemetry contracts | `mkdocs build --strict` is the release gate |

## v2.1.1 repository-tagged scope

Current work is included in the authorized repository tag `v2.1.1`. It is not
a published Go-semver v2 module, and the current hardening was not newly
live-verified.

1. Error-category/semantic-sentinel contracts, HTTP 417 caveats, and terminal
   MQTT/gRPC mappings.
2. Request-pipeline parity and observability hardening, including provider-order
   independence and sanitized error text.
3. OMS status reconciliation and account-scoped order tracking.
4. MQTT/channel shutdown, reconnect serialization, compare-and-swap state, and
   deterministic health recovery.
5. Trading and Broker FD event telemetry, cancellation metrics, and Trading
   all-runs shutdown.
6. Cancellation, leak, flakiness, and public-primitive test coverage.
7. Updated account-monitor/order/stream examples and synchronized documentation.

No item in this section is represented as a published module release merely
because it is included in the repository tag.

## Live-verification boundaries

| Surface | Current state |
|---|---|
| Core auth, selected HK market data, accounts, read-only assets, MQTT streaming, and Trading events | Previously exercised against available HK sandbox paths; not all endpoints or v2.1.1 hardening changes are live-verified |
| Footprint | Implemented and offline-tested; live request blocked by entitlement (`403`) |
| Options contracts and multi-leg orders | Implemented and offline-tested; HK sandbox may return `417`, and US live behavior is unavailable |
| US-only crypto, fund, screener, Broker FD, and related data | Implemented and offline-tested where covered; live verification blocked without US sandbox credentials |
| Broker API HK | Implemented and offline-tested in its own module; HK sandbox returns `401 ROUTE_NOT_PERMITTED` because scope is missing |
| Display Solution | Implemented and offline-tested; HK host returns `403` before endpoint-specific behavior can be verified |
| SSE news | Implemented; HK upstream currently returns `504` |

The generated [SDK ↔ API reconciliation](docs/reconciliation.md) remains authoritative for endpoint counts and path states. Its 2026-09-22 snapshot reports 209 implemented endpoints and 0 documented-only gaps, while also reporting summary-only and unresolved paths; do not summarize that report as zero discrepancies.

## Verification gates

Run from the repository root. These Makefile targets traverse the root module
and every nested module in `MODULES`:

```sh
make build
make vet
make test
make test-race
make cover
make lint
gofmt -l .                 # must print nothing
make docs
```

Root-only `go build ./...` and `go test ./...` are useful for package iteration
but do not traverse `broker/` or the nested example modules. `make cover`
records per-module measurements; it does not make coverage a behavioral
guarantee.

CI runs the root race matrix and nested-module build/vet/race checks, but its
60% coverage gate is root-only and it does not run the strict documentation
build. The local Makefile gates remain the release verification source.

## Remaining risks

- The v2.1.1 repository-tagged request, OMS, stream, and event-telemetry
  behavior has not been newly exercised against a live Webull service.
- US-only, Display Solution, Broker HK, entitlement-gated, SSE, and unsupported
  sandbox product paths remain blocked or unverified as listed in
  `IMPLEMENTATION_STATUS.md`.
- Stream dispatch is intentionally synchronous. A slow callback or full
  `DropBlock` channel creates head-of-line latency; cancellation and terminal
  close guarantee release, not asynchronous delivery.
- Broker FD events remain raw-only, do not expose response request ID/timestamp
  through `OnData`, support one active `Run`, and have no public option that
  injects a non-zero raw subscribe bitmask.
- The 2026-09-22 generated reconciliation still has four summary-only and 25
  unresolved SDK paths despite zero documented-only endpoint gaps.
- CI does not enforce nested coverage or a strict docs build; root aggregate
  coverage percentages are measurements, not correctness guarantees.
- The repository patch release is recorded as `v2.1.3`; the root module remains
  on `github.com/shing1211/webullapi4go` and stays on the v1 import path by
  decision, so `v2.x` tags are not installable Go-semver v2 modules and
  `v1.1.1` is the newest installable version. A `/v2` migration is declined
  rather than deferred; it would require a breaking import-path change and
  explicit maintainer approval.

## Out of scope

- Reintroducing a full `pkg/` service relocation.
- Deprecating the root service packages.
- Raw `decimal.Decimal` DTO migration.
- Unconditional or speculative `sync.Pool` use.
- Editing generated files under `docs/webull-api/`.
- Editing run artifacts under `docs/runs/` other than the explicitly authorized release/status records.
- Claiming blocked or newly hardened surfaces are live-verified without the required access.

## Reference

- Repository guide: `AGENTS.md`
- Build/test/lint: `Makefile`
- Current status: `IMPLEMENTATION_STATUS.md`
- Changelog: `CHANGELOG.md`
- Generated endpoint reconciliation: `docs/reconciliation.md`
