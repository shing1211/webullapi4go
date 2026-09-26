# Implementation Status

Last updated: 2026-09-26

- Latest repository tag: **`v2.1.3`** (2026-09-26)
- Current hardening: **tagged in repository `v2.1.3`**; introduced in `v2.1.1`
- Module path: **`github.com/shing1211/webullapi4go`**, kept on the v1 import path
  **by decision**; no `/v2` migration is planned
- Installable version: **`v1.1.1`** — the module proxy serves only the `v1.x`
  line, so `v2.x` tags are not installable with `go get`. Consume newer work by
  pinning a commit
- Offline evidence recorded 2026-09-26: module-aware race/vet checks passed;
  `make cover` measured 73.6% aggregate root coverage and 80.8% in nested
  `broker/` (measurements, not behavior guarantees)
- Live verification: partial historical evidence only; the current hardening
  was not newly live-verified. The repository tag is not a published
  Go-semver v2 module.

## Status definitions

| State | Meaning |
|---|---|
| Implemented | The public SDK surface exists in the current tree |
| Offline-tested | Credential-free unit tests or local HTTP/gRPC/MQTT fakes pass |
| Live-verified | Successfully exercised against an available Webull environment |
| Blocked | Live verification is prevented by credentials, scope, entitlement, host, or sandbox data |
| Repository-tagged | Included in the authorized repository Git tag `v2.1.1`; this is not a published Go-semver v2 module |
| Released | Included in a published module release; the unchanged root path means `v2.x` Git tags are not published Go-semver v2 modules |

Implementation and offline tests do not imply live verification.

## Current architecture

| Layer | Package | Purpose |
|---|---|---|
| Core | `client/` | Config, signing, tokens, HTTP transport, request pipeline, resilience |
| Services | `data/`, `stream/`, `trade/`, `events/`, `connect/`, `display/`, `brokerfd/`, `brokerfd/events/` | Public service clients remain at the repository root |
| Broker HK | `broker/` | Separate Go module |
| Shared foundations | `pkg/errors`, `pkg/observability`, `pkg/resilience`, `pkg/transport`, `pkg/types` | Reusable public primitives |
| Domain | `pkg/domain/money`, `pkg/domain/order` | `money.Money` and the OMS state machine |
| Convenience | `webull/` | Thin aliases for the core client; not a replacement service facade |

The full-service relocation to `pkg/`, root-service deprecation shims, and raw `decimal.Decimal` DTO migration are superseded and closed. The streaming path does not use `sync.Pool`.

## Service status

| Area | Implemented | Offline-tested | Live-verified | Blocked or unverified |
|---|---:|---:|---|---|
| Core `client/` | Yes | Yes | Core auth and selected HK HTTP paths | Current request-pipeline hardening is not newly live-verified |
| Authentication | Yes | Yes | Token create/check/ensure in HK sandbox | Production token activation still requires the Webull App 2FA flow |
| Market Data HTTP `data/` | Yes | Yes | Selected AAPL, fundamentals, watchlist, and related HK calls | US-only data and entitlement-gated calls remain partial |
| Market Data MQTT `stream/` | Yes | Yes | Basic MQTT-over-WebSocket receive/reconnect paths | Current channel, shutdown, and health hardening is offline-tested only |
| Trading HTTP `trade/` | Yes | Yes | Accounts, balances, positions, previews, and guarded order paths | Multi-leg/futures/event-contract live behavior is unavailable or rejected in HK |
| Trading Events `events/` | Yes | Yes | Basic stream path was previously exercised | Current telemetry/reconciliation behavior is offline-tested only |
| Connect OAuth `connect/` | Yes | Yes | — | US live access unavailable |
| Display Solution `display/` | Yes | Yes | — | HK host returns `403` before endpoint verification |
| Broker API HK `broker/` | Yes | Yes, separate module | — | HK sandbox returns `401 ROUTE_NOT_PERMITTED` because scope is missing |
| Broker FD HTTP `brokerfd/` | Yes | Yes | — | US-only; no US sandbox credentials, HK returns `404` |
| Broker FD Events `brokerfd/events/` | Yes | Yes | — | US-only; no live verification available |
| Shared foundations | Yes | Yes | Not applicable | — |

## Endpoint coverage

The generated [SDK ↔ API reconciliation](docs/reconciliation.md) is authoritative. Its 2026-09-22 snapshot reports:

| Measure | Count |
|---|---:|
| Implemented endpoints | 209 |
| Documented-only gaps | 0 |
| Exact OpenAPI JSON path matches | 180 |
| Matches the docs summary only | 4 |
| Differs from both sources | 0 |
| Unresolved SDK path | 25 |

The SDK therefore has no documented-only endpoint gaps, but it is inaccurate to call the snapshot a zero-discrepancy report. Do not hand-edit generated pages under `docs/webull-api/`; regenerate derived documentation with the doc generator when its inputs or the SDK change.

## v2.1.1 repository-tagged work

| Work | State | Verification |
|---|---|---|
| Category matching plus identity-specific semantic sentinels | Implemented | Offline-tested across public client/MQTT/event sentinels |
| HTTP 417 compatibility mapping and MQTT/event terminal mappings | Implemented | Offline-tested; 417 API-message caveat preserved |
| Shared `Do`, `DoBroker`, and `DoStream` attempt pipeline | Implemented | Offline-tested; not newly live-verified |
| One-based hook attempts, stable correlation IDs, status/latency telemetry, W3C propagation, and `SafeErrorText` | Implemented | Offline-tested with real OTel readers and local servers |
| Provider-order-independent resilience metrics and caller-owned custom breakers | Implemented | Offline-tested |
| OMS status snapshots, stale/terminal reconciliation, and correct failure-scene mapping | Implemented | Offline-tested |
| Account-scoped order registry and successful replace/cancel transitions | Implemented | Offline-tested |
| Stable auto-generated client order IDs without mutating caller requests | Implemented | Offline-tested |
| Compare-and-swap stream state and deterministic data-only health recovery | Implemented | Offline-tested with deterministic watchdog ticks/barriers |
| Terminal MQTT/channel shutdown, cancellation unblocking, and synchronous registration-order dispatch | Implemented | Offline-tested for quote/snapshot/tick channels |
| MQTT Connect cancellation, Connect/Close races, and late-callback suppression | Implemented | Offline-tested |
| Serialized resubscription and active-subscription replay | Implemented | Offline-tested |
| Trading and Broker FD gRPC event telemetry, including cancellation/failure outcomes | Implemented | Offline-tested with local gRPC servers |
| Trading `Close` cancellation of all active `Run` calls | Implemented | Offline-tested |
| Cancellation, strict data/event leak checks, and public resilience/transport/type/`webull` coverage | Implemented | Offline race/vet evidence recorded 2026-09-25 |
| Account monitor and updated order/stream examples | Implemented | Included in offline module build/test |
| Documentation reconciliation | Implemented | Final strict MkDocs/link/diff gates recorded in the F1 report |

These items are included in repository tag `v2.1.1` (2026-09-25). They are
not a published Go-semver v2 module, and the current hardening was not newly
live-verified.

## Version history

The `v2.x` rows below are repository-tag records, not published Go-semver v2
modules; `v2.1.3` is the current authorized repository tag and earlier rows are
historical. The root module path remains `github.com/shing1211/webullapi4go` and
stays on the v1 import path by decision.

| Version | Date | Scope | Status |
|---|---|---|---|
| v0.1.0 | 2026-09-18 | Auth, core HTTP, Market Data HTTP, and MQTT streaming | Released |
| v0.2.1 | 2026-09-18 | Accounts, balances, and positions | Released |
| v0.2.2 | 2026-09-18 | Stock-order lifecycle and queries | Released |
| v0.2.3 | 2026-09-18 | Market-specific rules and HK BCAN | Released |
| v0.2.4 | 2026-09-18 | Single-leg options orders | Released |
| v0.2.5 | 2026-09-18 | US combo orders | Released |
| v0.2.6 | 2026-09-18 | News SSE through `DoStream` | Released |
| v0.3.0 | 2026-09-18 | Trading events over gRPC | Released |
| v0.4.0 | 2026-09-18 | Market Data fundamentals | Released |
| v0.5.0 | 2026-09-19 | Display, corporate actions, fund, crypto, and screener surfaces | Released; historical provisional paths later reconciled |
| v0.6.0 | 2026-09-20 | Multi-leg options, futures validation, and discovery helpers | Released; historical provisional behavior later reconciled |
| v0.7.0 | 2026-09-21 | Event contracts, Broker HK, Broker FD, and Broker FD events | Released |
| v0.8.0 | 2026-09-21 | Reserved | No code release |
| v0.9.0 | 2026-09-22 | GoDoc coverage, HK derivatives helpers, and examples | Released |
| v0.9.1 | 2026-09-21 | Watchlist boolean response, `DoBroker`, and examples | Released |
| v0.9.2 | 2026-09-21 | Broker HK path correction and live scope finding | Released |
| v1.0.0 | 2026-09-22 | Futures fixes and API stability audit | Released |
| v1.0.1 | 2026-09-22 | Flexible futures unit and order-ID example fix | Released |
| v1.0.2 | 2026-09-22 | Removed undocumented endpoints and corrected the option-contract path | Released |
| v1.0.3 | 2026-09-22 | Generated API reference, reconciliation, and path probe | Released |
| v1.1.0 | 2026-09-22 | Full documented endpoint coverage and zero documented-only gaps | Released |
| v1.1.1 | 2026-09-23 | DevOps, security checks, multi-OS CI, leak tests, and fuzzing | Released |
| v2.0.0 | 2026-09-23 | Public error/transport/resilience/domain foundations and thin `webull` aliases; root services retained | Historical tag |
| v2.0.1 | 2026-09-23 | `money.Money` conversion in `data/` and `trade/` | Historical tag |
| v2.0.2 | 2026-09-23 | `money.Money` conversion in `brokerfd/` | Historical tag |
| v2.0.3 | 2026-09-23 | MQTT/paho goroutine cleanup compatibility in leak tests | Historical tag |
| v2.0.4 | 2026-09-24 | Clock correction, idempotency helpers, transport tuning, resilience preset | Historical tag |
| v2.0.5 | 2026-09-24 | Request interceptors/hooks and initial OMS integration | Historical tag |
| v2.0.6 | 2026-09-24 | Stream state/channels, structured logging, and OTel tracing | Historical tag |
| v2.0.7 | 2026-09-24 | HTTP latency, breaker, reconnect, and drop metrics | Historical tag |
| v2.0.8 | 2026-09-24 | MQTT and resubscription context cleanup | Historical tag |
| v2.0.9 | 2026-09-24 | Structured public error codes and wrapped sentinels | Historical tag |
| v2.1.0 | 2026-09-24 | `go vet` mutex-copy fixes | Historical tag; not a published v2 module |
| v2.1.1 | 2026-09-25 | Error specificity, request/OMS, stream/MQTT, event telemetry, cancellation/leak coverage, and documentation hardening | Historical tag; not a published v2 module |
| v2.1.2 | 2026-09-26 | OpenTelemetry v1.46.0 across all modules, GitHub Actions bumps, pinned `govulncheck`, Dependabot grouping, and documentation policy updates | Historical tag; not a published v2 module |
| v2.1.3 | 2026-09-26 | Honest nightly live-verification signal, removal of the dead `internal/` shims and duplicate resilience tree, and direct tests for the order reconciliation, MQTT, money, error, and region surfaces | Repository tag; not a published v2 module; offline-tested and not newly live-verified |

## Previously exercised HK surface

Earlier sandbox runs exercised the core token flow, selected AAPL market-data and fundamentals calls, watchlists, accounts, balances, positions, order preview paths, MQTT streaming, and Trading Events. This is not a claim that all 209 endpoints or all current hardening behavior passed a live test.

## Test status

- The module-aware offline race/vet checks recorded for this run passed on
  2026-09-25. `make test` and `make test-race` traverse the root, `broker/`, and
  nested example modules; root-only `go test ./...` does not.
- `make cover` measured 73.6% aggregate coverage for the root module and 80.8%
  for `broker/` on 2026-09-26. These dated measurements are not a guarantee of
  behavior or correctness and must not be combined into one aggregate
  percentage.
- Unit tests are offline and credential-free. Live tests use the actual
  `Sandbox` selector and remain environment-gated and skipped by default.
- Data and both event packages use strict goroutine-leak checks; DNS-dependent
  failure tests use local deterministic dialers.
- CI runs root race tests across three operating systems and build/vet/race for
  nested modules. Its 60% coverage gate is root-only; nested coverage and a
  strict MkDocs build are not CI gates.
- No live verification was added for the v2.1.1 repository-tagged hardening.

## Examples

| Example | Directory | Purpose |
|---|---|---|
| auth | `examples/auth/` | Token create/check/reuse |
| marketdata | `examples/marketdata/` | AAPL snapshot and bars |
| streaming | `examples/streaming/` | MQTT-over-WebSocket subscription |
| watchlist | `examples/watchlist/` | Read-only watchlists |
| account | `examples/account/` | Accounts, balances, and positions |
| account-monitor | `examples/account-monitor/` | Bounded, read-only periodic balance monitor |
| order | `examples/order/` | Guarded AAPL preview/place/cancel flow with persisted idempotency key |
| events | `examples/events/` | Trading order events over gRPC |
| data-fundamentals | `examples/data-fundamentals/` | Fundamentals endpoints |
| watchlist-cmd | `examples/watchlist-cmd/` | Explicitly gated watchlist mutations |
| broker-probe | `examples/broker-probe/` | Broker HK read-only probe |
| futures-probe | `examples/futures-probe/` | Futures discovery and market-data probe |
| options-multi-leg | `examples/options-multi-leg/` | Multi-leg strategy probe |
| options | `examples/options/` | HK options discovery probe |
| brokerfd | `examples/brokerfd/` | Broker FD US probe |
| brokerfd-events | `examples/brokerfd-events/` | Broker FD gRPC probe |
| probe | `examples/probe/` | General sandbox probe |
| path-probe | `examples/path-probe/` | Generated path discrepancy probe |

## Blocked and unverified surfaces

1. **Display Solution:** the HK host returns `403`, including token creation.
2. **US surfaces:** no US sandbox credentials are available for crypto, funds, screener, Broker FD, or other US-only behavior.
3. **Broker API HK:** returns `401 ROUTE_NOT_PERMITTED` because the app lacks the required scope.
4. **HK symbols:** sandbox data is limited to `AAPL`.
5. **Footprint:** requires an entitlement and returns `403` in the sandbox.
6. **Option contracts:** may be absent for the sandbox account/symbol and return `417 Invalid Symbol`. HTTP 417 is mapped to `INVALID_TOKEN` for compatibility even when the API message reports business validation.
7. **Multi-leg options:** HK rejects non-`SINGLE` strategies with `417`; US behavior is unverified.
8. **Futures/event-contract trading:** request validation is offline-tested; live product behavior is unverified.
9. **Order-book depth:** may be empty outside regular trading hours.
10. **Plain MQTT:** port `1883` may be blocked; prefer the configured MQTT-over-WebSocket endpoint.
11. **SSE news:** the HK upstream currently returns `504`.
12. **Generated path status:** four paths match only the docs summary and 25 are unresolved; see the generated reconciliation rather than claiming zero discrepancies.
13. **Broker FD event API:** data remains raw, `OnData` omits response request ID/timestamp, only one active `Run` is supported, and no public option injects a non-zero raw subscribe bitmask.
14. **Stream backpressure:** callbacks and channel dispatch are synchronous; a slow handler or full `DropBlock` subscriber causes head-of-line delay until cancellation or terminal close.
15. **CI gaps:** nested coverage and strict documentation builds are Makefile/release gates, not CI gates; measured percentages are not behavior guarantees.

## Next steps

1. Live-verify the v2.1.1 repository-tagged request, OMS, stream, and event-telemetry changes when suitable credentials and non-production test access are available.
2. Supply US sandbox credentials for the US-only surfaces.
3. Resolve the 25 generated unresolved paths and the four summary-only matches through the doc generator and official OpenAPI sources.
4. Decide whether Broker FD needs a public raw-subscribe option, richer `OnData` metadata, and all-runs lifecycle parity before any future release tag.
5. Benchmark a separately approved asynchronous stream-dispatch design only if synchronous head-of-line latency is unacceptable.
6. Run the full race, vet, formatting, lint, and strict documentation gates before any future release tag.
