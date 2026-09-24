# Implementation Status

Last updated: 2026-09-25

- Latest tagged release: **`v2.1.0`**
- Current hardening: **Unreleased**
- Offline verification: `go test ./...` passes for the root module and nested `broker/` module
- Live verification: partial

## How to read this page

| State | Meaning |
|---|---|
| Implemented | The surface exists in the current tree |
| Offline-tested | Credential-free unit tests or local fakes pass |
| Live-verified | Successfully exercised against an available Webull environment |
| Blocked | Credentials, scope, entitlement, host, or sandbox data prevent verification |
| Released | Included in a tag; current working-tree hardening is not released |

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

## Unreleased hardening

| Work | State |
|---|---|
| Shared `Do` / `DoBroker` / `DoStream` request pipeline | Implemented, offline-tested |
| One-based attempt telemetry, stable correlation IDs, status/latency data, and trace propagation | Implemented, offline-tested |
| OMS snapshot reconciliation and failure-scene mapping | Implemented, offline-tested |
| Account-scoped tracking and stable auto-generated client IDs | Implemented, offline-tested |
| Idempotent MQTT/channel shutdown and blocked-dispatch cancellation | Implemented, offline-tested |
| Serialized resubscription and data-only health recovery | Implemented, offline-tested |
| Trading and Broker FD event telemetry | Implemented, local-gRPC tested |
| Account monitor and updated order/stream examples | Implemented, included in root build/test |
| Documentation reconciliation | In progress; strict MkDocs build is the gate |

These items are not released until a tag is created.

## Released milestones

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

Earlier v0.x and v1.0 milestones remain recorded in the root `IMPLEMENTATION_STATUS.md`.

## Verification performed

- `go test ./...` passes in the root module on 2026-09-25.
- `go test ./...` passes in `broker/` on 2026-09-25.
- Unit tests are offline and credential-free.
- Live tests remain environment-gated.
- No `TODO` markers remain in Go source.
- No new live validation was performed for the current Unreleased hardening.

## Live-verification boundaries

Previously exercised HK paths include core tokens, selected AAPL data and fundamentals, watchlists, accounts, balances, positions, order previews, MQTT streaming, and Trading Events. That does not cover all 209 endpoints or the current hardening.

Blocked or unverified areas:

- Display Solution: host-level `403`.
- US-only crypto, fund, screener, and Broker FD: no US sandbox credentials.
- Broker API HK: missing scope (`401 ROUTE_NOT_PERMITTED`).
- Footprint: missing entitlement (`403`).
- Options/multi-leg: limited sandbox contracts and non-`SINGLE` strategies rejected with `417`.
- Futures and event-contract trading: validation is offline-tested; live product behavior is unverified.
- SSE news: upstream `504`.
- Generated paths: four summary-only matches and 25 unresolved paths remain.

## Next steps

1. Live-verify current request, OMS, stream, and event telemetry with suitable non-production access.
2. Add US sandbox verification for US-only surfaces.
3. Resolve generated path states through the doc generator and official sources.
4. Run race, vet, formatting, lint, and `mkdocs build --strict` before tagging.
