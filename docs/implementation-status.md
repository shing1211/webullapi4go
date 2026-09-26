# Implementation Status

Last updated: 2026-09-26

- Latest repository tag: **`v2.1.2`** (2026-09-26)
- Current hardening: **tagged in repository `v2.1.2`**; introduced in `v2.1.1`
- Module path: **`github.com/shing1211/webullapi4go`**, kept on the v1 import path
  **by decision**; no `/v2` migration is planned
- Installable version: **`v1.1.1`** — the module proxy serves only the `v1.x`
  line, so `v2.x` tags are not installable with `go get`. Consume newer work by
  pinning a commit
- Offline evidence recorded 2026-09-25: module-aware race/vet checks passed;
  `make cover` measured 71.6% aggregate root coverage and 80.8% in nested
  `broker/` (measurements, not guarantees)
- Live verification: partial historical evidence only; the current hardening
  was not newly live-verified. The repository tag is not a published
  Go-semver v2 module.

## How to read this page

| State | Meaning |
|---|---|
| Implemented | The surface exists in the current tree |
| Offline-tested | Credential-free unit tests or local fakes pass |
| Live-verified | Successfully exercised against an available Webull environment |
| Blocked | Credentials, scope, entitlement, host, or sandbox data prevent verification |
| Repository-tagged | Included in the authorized repository Git tag `v2.1.1`; this is not a published Go-semver v2 module |
| Released | Included in a published module release; the unchanged root path means `v2.x` Git tags are not published Go-semver v2 modules |

Implementation and offline tests do not imply live verification.

## Architecture status

| Layer | Canonical location | Status |
|---|---|---|
| Core client | `client/` | Preserved; current pipeline hardening is offline-tested |
| Public services | `data/`, `stream/`, `trade/`, `events/`, `connect/`, `display/`, `brokerfd/`, `brokerfd/events/` | Preserved at repository root |
| Broker HK | `broker/` | Preserved as a separate Go module |
| Shared foundations | `pkg/errors`, `pkg/observability`, `pkg/resilience`, `pkg/transport`, `pkg/types` | Implemented and offline-tested |
| Domain | `pkg/domain/money`, `pkg/domain/order` | `money.Money` and OMS reconciliation implemented and offline-tested |
| Convenience | `webull/` | Thin core-client aliases, not an aggregate facade |

The proposed full relocation of service packages under `pkg/`, root-service deprecated shims, and raw `decimal.Decimal` DTO migration are superseded and closed. The stream path has no `sync.Pool`.

## Service matrix

| Area | Implemented | Offline-tested | Live-verified | Blocked or unverified |
|---|---:|---:|---|---|
| Core `client/` | Yes | Yes | Core auth and selected HK HTTP paths | Current pipeline hardening is not newly live-verified |
| Authentication | Yes | Yes | Token lifecycle in HK sandbox | Production activation requires the Webull App 2FA flow |
| Market Data HTTP | Yes | Yes | Selected AAPL/fundamentals/watchlist calls | US-only and entitlement-gated calls remain partial |
| Market Data MQTT | Yes | Yes | Basic MQTT-over-WebSocket paths | Current channel/health/shutdown hardening is offline-tested only |
| Trading HTTP | Yes | Yes | Accounts, assets, previews, guarded order paths | Multi-leg/futures/event live behavior is unavailable or rejected in HK |
| Trading Events | Yes | Yes | Basic stream path was previously exercised | Current telemetry is offline-tested only |
| Connect OAuth | Yes | Yes | — | US live access unavailable |
| Display Solution | Yes | Yes | — | HK host returns `403` |
| Broker API HK | Yes | Yes | — | HK returns `401 ROUTE_NOT_PERMITTED` because scope is missing |
| Broker FD HTTP | Yes | Yes | — | US-only; HK returns `404`, no US credentials |
| Broker FD Events | Yes | Yes | — | US-only; no live verification available |
| Shared foundations | Yes | Yes | Not applicable | — |

## Endpoint reconciliation

The generated [SDK ↔ API Reconciliation](reconciliation.md) is authoritative. Its 2026-09-22 snapshot reports:

| Measure | Count |
|---|---:|
| Implemented endpoints | 209 |
| Documented-only gaps | 0 |
| Exact OpenAPI JSON path matches | 180 |
| Matches docs summary only | 4 |
| Differs from both sources | 0 |
| Unresolved SDK path | 25 |

There are no documented-only endpoint gaps, but the snapshot is not a zero-discrepancy report. Generated pages under `webull-api/` are not hand-edited.

## v2.1.1 repository-tagged hardening

| Work | State |
|---|---|
| Category matching and identity-specific semantic sentinels | Implemented, offline-tested |
| HTTP 417 compatibility mapping and MQTT/event terminal mappings | Implemented, offline-tested |
| Shared `Do` / `DoBroker` / `DoStream` request pipeline | Implemented, offline-tested |
| One-based attempt telemetry, stable correlation IDs, W3C propagation, and `SafeErrorText` | Implemented with real OTel/local-server assertions |
| Provider-order-independent resilience metrics and caller-owned custom breakers | Implemented, offline-tested |
| OMS snapshot reconciliation and failure-scene mapping | Implemented, offline-tested |
| Account-scoped tracking and stable auto-generated client IDs | Implemented, offline-tested |
| Compare-and-swap stream state and deterministic data-only health recovery | Implemented with deterministic watchdog tests |
| Idempotent MQTT/channel shutdown, blocked-dispatch cancellation, and registration-order head-of-line behavior | Implemented, offline-tested for all channel types |
| MQTT Connect cancellation, Connect/Close races, and late-callback suppression | Implemented, offline-tested |
| Serialized resubscription and active-subscription replay | Implemented, offline-tested |
| Trading and Broker FD event telemetry, including failure/cancellation outcomes | Implemented, local-gRPC tested |
| Trading `Close` cancellation of all active `Run` calls | Implemented, offline-tested |
| Cancellation, strict data/event leak checks, and public resilience/transport/type/`webull` coverage | Implemented, offline evidence recorded |
| Account monitor and updated order/stream examples | Implemented, included in offline module build/test |
| Documentation reconciliation | Implemented; final F1 report records strict docs/link/diff gates |

These items are included in repository tag `v2.1.1` (2026-09-25). They are
not a published Go-semver v2 module, and the current hardening was not newly
live-verified.

## Milestones and historical tags

The `v2.x` rows below are repository-tag records, not published Go-semver v2
modules; `v2.1.2` is the current authorized repository tag and earlier rows are
historical. The root module path remains `github.com/shing1211/webullapi4go` and
stays on the v1 import path by decision.

| Version | Scope |
|---|---|
| v1.1.0 | Full documented endpoint coverage; 209 implemented, 0 documented-only gaps |
| v1.1.1 | DevOps, security checks, multi-OS CI, leak tests, and fuzzing |
| v2.0.0 | Public error/transport/resilience/domain foundations and thin `webull` aliases; root services retained |
| v2.0.1–v2.0.2 | `money.Money` DTO conversion in `data/`, `trade/`, and `brokerfd/` |
| v2.0.3–v2.0.4 | MQTT cleanup, clock correction, idempotency helpers, transport tuning, and resilience preset |
| v2.0.5–v2.0.7 | Interceptors/hooks, initial OMS, stream state/channels, slog, OTel tracing, and metrics |
| v2.0.8–v2.0.9 | Context hygiene and structured public errors |
| v2.1.0 | `go vet` mutex-copy fixes |
| v2.1.1 | Current hardening, including error contracts, request/OMS reconciliation, stream/MQTT lifecycle, event telemetry, cancellation/leak coverage, and documentation reconciliation; repository tag, not a published v2 module |
| v2.1.2 | OpenTelemetry v1.46.0 across all modules, GitHub Actions bumps, pinned `govulncheck`, Dependabot grouping, and documentation policy updates; repository tag, not a published v2 module |

Earlier v0.x and v1.0 milestones remain recorded in the root `IMPLEMENTATION_STATUS.md`.

## Verification performed

- Module-aware offline race/vet checks were recorded as passing on 2026-09-25.
  `make test` and `make test-race` traverse the root, `broker/`, and nested
  example modules; root-only `go test ./...` does not.
- `make cover` measured 71.6% aggregate root coverage and 80.8% in `broker/`.
  These dated measurements are not correctness guarantees and are not combined
  into one repository-wide percentage.
- Unit tests are offline and credential-free. Live tests use the `Sandbox`
  selector and remain environment-gated.
- Data and both event packages use strict goroutine-leak checks; event tests
  cancel and wait for completed runs.
- CI's 60% coverage gate is root-only. Nested coverage and
  `mkdocs build --strict` remain Makefile/release gates rather than CI gates.
- No new live validation was performed for the v2.1.1 repository-tagged hardening.

## Live-verification boundaries

Previously exercised HK paths include core tokens, selected AAPL data and fundamentals, watchlists, accounts, balances, positions, order previews, MQTT streaming, and Trading Events. That does not cover all 209 endpoints or the current hardening.

Blocked or unverified areas:

- Display Solution: host-level `403`.
- US-only crypto, fund, screener, and Broker FD: no US sandbox credentials.
- Broker API HK: missing scope (`401 ROUTE_NOT_PERMITTED`).
- Footprint: missing entitlement (`403`).
- Options/multi-leg: limited sandbox contracts and non-`SINGLE` strategies rejected with `417`. The SDK's `INVALID_TOKEN` category is a compatibility mapping for every HTTP 417 and does not prove a token failure.
- Futures and event-contract trading: validation is offline-tested; live product behavior is unverified.
- SSE news: upstream `504`.
- Generated paths: four summary-only matches and 25 unresolved paths remain.

## Remaining risks

- Current v2.1.1 repository-tagged behavior has not been newly live-verified;
  US-only, Display, Broker HK, entitlement-gated, SSE, and unsupported sandbox
  paths remain blocked or unverified.
- Broker FD events are raw-only, `OnData` omits request ID/timestamp, one active
  `Run` is supported, and no public option injects a non-zero raw subscribe
  bitmask.
- Stream callbacks/channels are synchronous; a slow handler or full
  `DropBlock` subscriber causes head-of-line delay until cancellation or close.
- Four summary-only and 25 unresolved generated paths remain.
- Nested coverage and strict docs are not CI gates; percentages are measurements,
  not behavior guarantees.

## Next steps

1. Live-verify the v2.1.1 repository-tagged request, OMS, stream, and event telemetry with suitable non-production access.
2. Add US sandbox verification for US-only surfaces.
3. Resolve generated path states through the doc generator and official sources.
4. Decide whether Broker FD needs public subscribe-bitmask, richer raw metadata, and all-runs lifecycle APIs before release.
5. Run module-aware race, vet, formatting, lint, and `mkdocs build --strict` before any separately approved release tag.
