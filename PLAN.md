# Enhancement / Production Hardening Plan

Status: **Approved · pending execution** · Target release: **`v2.0.0`** (breaking Go-API changes)

> Approving this plan explicitly authorizes edits to `.github/workflows/ci.yml` and `.golangci.yml` that AGENTS.md otherwise restricts.

## Compatibility strategy

v2.0.0 introduces breaking Go-API changes (package relocation, `decimal.Decimal` on DTOs, new facade import path). v1.x is preserved on a branch for non-migrating consumers; deprecated shims at the old import paths delegate to the new locations for one release cycle and are removed in v3. JSON wire format is unchanged — `decimal.Decimal` marshal/unmarshal round-trips as strings.

## Structure

```
Resolved Decisions        ← what + why
Phases 1–7              ← the work (sequenced by execution order)
Execution Order         ← checkpoint plan
Risks & Mitigations
Out of Scope
Verification Gates
Reference
```

## Resolved Decisions

### Scope & strategy

| Decision | Resolution |
|---|---|
| Scope | All items in scope — SDK hardening, OMS, decimal, `/pkg` relayout, observability, DevOps |
| Compatibility | v2.0.0 breaking changes; v1.x preserved; deprecated shims for one cycle |
| goreleaser | **Skip** — a library ships no binaries; revisit only if CLI tools are added |

### Technical decisions

| # | Decision | Resolution |
|---|---|---|
| T1 | Tracing/metrics | **Direct `go.opentelemetry.io/otel` API-only dep**, no-op default, `WithTracerProvider`/`WithMeterProvider`; `otelhttp` for REST instrumentation |
| T2 | Resilience defaults | **Keep opt-in + add `WithResiliencePreset()`**; defaults unchanged from v1.1.0 (non-breaking) |
| T3 | Package layout | **Relocate to `/pkg/...`** (`pkg/domain`, `pkg/services`, `pkg/transport`, `pkg/errors`); old root packages become deprecated shims |
| T4 | Decimal | **`decimal.Decimal` on DTOs** (strict wrapper; JSON round-trips as strings) |
| T5 | OMS | **SDK-side domain model** (`OrderState` + `ApplyEvent(Event) (State, error)`) |
| T6 | HMAC signing | **REST = HMAC-SHA1 (MD5 body digest); gRPC = HMAC-SHA256** — documented, not changed |
| T7 | Go version matrix | `[1.26.x, stable]` — go.mod is `1.26`, so `1.22+` is impossible |

## Phases

### Phase 1 — DevOps & testing infrastructure

Pure config; unblocks later gates.

| # | Item | Notes |
|---|---|---|
| 1.1 | Makefile targets | `build`, `test`, `test-race`, `cover`, `lint`, `fuzz`, `vuln`, `docs` (today only proto codegen) |
| 1.2 | `gosec` | Enable in `.golangci.yml`; exclude known false positives (e.g. `math/rand` in jitter) |
| 1.3 | `govulncheck` | Makefile target + CI job |
| 1.4 | Coverage gate | **Measure baseline first** (`go test ./... -cover`), then a job + realistic gate (85% on touched/new code first; raise repo-wide only if baseline supports it) |
| 1.5 | CI matrix | `ubuntu-latest` / `macos-latest` / `windows-latest` × Go `[1.26.x, stable]`; Windows needs a `GOTMPDIR` workaround for the known `a.out.exe` file-lock issue; `gofmt` and `go vet` jobs retained |
| 1.6 | Dependabot | `gomod` + `github-actions`, weekly |
| 1.7 | Hardening tests | `goleak` in test mains; `FuzzDecodeQuote/Snapshot/Tick` + notice-JSON fuzz; mock MQTT broker; burst/reconnect/drop simulations; table-driven tests throughout |

### Phase 2 — Architecture & domain foundation (blocking)

Largest phase; everything else depends on the new layout. Perform relocation in one pass with shims.

| # | Item | Location | Notes |
|---|---|---|---|
| 2.1 | New canonical layout | `pkg/domain`, `pkg/services`, `pkg/transport`, `pkg/errors` | `pkg/domain` holds models, enums, decimal helpers, public errors; `pkg/services` holds MarketData, Trading, Account, Streaming; `pkg/transport` holds REST/MQTT/gRPC factories and the `Doer` interface |
| 2.2 | Deprecated shims | `client/`, `data/`, `trade/`, `stream/`, `events/`, `broker/`, `brokerfd/`, `display/`, `connect/` | Re-export/delegate to `pkg/` locations with `// Deprecated:`; removed in v3 |
| 2.3 | Facade package | new `webull/` | `New()` + `Client`, `MarketData`, `Trading`, `Account`, `Streaming` interfaces and concrete implementations; compile-time `var _ MarketData = (*svc.MarketData)(nil)` asserts |
| 2.4 | Decimal on DTOs | ~250 numeric `string` fields across `data/`, `trade/`, `pkg/domain` | Strict `decimal.Decimal` wrapper (errors on non-numeric); JSON round-trips as strings; Go-side break requires a v2 migration guide |
| 2.5 | Computed helpers | `pkg/domain/money` | `Position.UnrealizedPnL()`, `AssetsBalance.MarginUtilization()`, `OrderRequest.Notional()` |
| 2.6 | Interfaces & `Doer` | `pkg/transport` and per-service packages | `RESTClient`, `MQTTClient`, `GRPCClient` enable pure unit mocks |
| 2.7 | OMS state machine | `pkg/domain/order` | `OrderState` enum with validated transitions; `ApplyEvent(Event) (State, error)`; events: Placed, Acknowledged, Fill, PartialFill, CancelRequest, Cancel, Reject, Expire; `OnStateChanged` hook |

### Phase 3 — Security & resilience

| # | Item | Notes |
|---|---|---|
| 3.1 | Clock-drift correction | Offset learned from response `Date` header → `x-timestamp`; bounded; `WithClockDriftCorrection(false)` opt-out |
| 3.2 | Idempotency helpers | `NewClientOrderID()` (UUIDv4), `ClientOrderIDFrom()` (deterministic hash), `WithAutoClientOrderID(true)`; reuse/persistence policy stays app-side |
| 3.3 | Transport tuning | Custom `http.Transport` defaults (pool/keep-alive/timeouts); `WithHTTPTransport` override wins |
| 3.4 | Retry full jitter | `rand(0, min(cap, base·2ⁿ))` behind `WithResiliencePreset()`; default stays ±20% (non-breaking) |
| 3.5 | Resilience preset | `WithResiliencePreset(production)` composes rate-limit + breaker + full-jitter defaults |
| 3.6 | HMAC correctness | REST = HMAC-SHA1, gRPC = HMAC-SHA256 — document in `docs/authentication.md`; do **not** change REST to SHA256 |

### Phase 4 — REST client & trading hardening

| # | Item | Notes |
|---|---|---|
| 4.1 | Interceptor pipeline | Composably ordered: rate-limit → breaker → sign → send → record; logging/metrics/tracing hooks plug in |
| 4.2 | OMS integration | `PlaceOrder` returns an `Order` with initial state; cancel/amend validate state transitions; websocket events feed `ApplyEvent` |
| 4.3 | Validation helpers | Enum/state checks before send; typed request builders |

### Phase 5 — Streaming engine (MQTT)

| # | Item | Notes |
|---|---|---|
| 5.1 | Connection state machine + health | Disconnected → Connecting → Connected → Reconnecting → Degraded → Closed; message-age watchdog; `OnHealth`; re-snapshot via `Grab:true`. Note: Webull MQTT protobuf carries **no sequence numbers**, so true seq-gap detection is impossible; health watchdog + re-snapshot is the substitute |
| 5.2 | Channel mux + backpressure | Per-subscription bounded `<-chan`; configurable drop policy (`block` default / drop-oldest / sample); drop counters |
| 5.3 | `sync.Pool` hot path | Include, gated by benchmarks — only land if allocation pressure is proven |
| 5.4 | Resubscribe & state preservation | Formalize existing behavior; expose lifecycle events |

### Phase 6 — Observability

| # | Item | Notes |
|---|---|---|
| 6.1 | `slog` | `WithLogger(*slog.Logger)`, no-op default; per-request correlation ID in logs and error context |
| 6.2 | OpenTelemetry | `go.opentelemetry.io/otel` (API-only, no-op default) + `otelhttp` for REST client instrumentation; spans for MQTT dispatch and gRPC; `WithTracerProvider`/`WithMeterProvider` options |
| 6.3 | Metrics hooks | Latency histogram, drop/reconnect counters, breaker state transitions |
| 6.4 | Correlation propagation | Request ID threaded through REST/MQTT/gRPC → logs + span context |

### Phase 7 — Docs, examples & release

| # | Item | Notes |
|---|---|---|
| 7.1 | Documentation | README (architecture, thread-safety, error handling); error-handling guide (public errors); OTel integration guide; decimal migration guide; v1→v2 migration guide |
| 7.2 | Examples | Channel streaming with `context` timeout; limit order with idempotency key; async margin/balance monitor |
| 7.3 | Changelog & status | `CHANGELOG.md` `[Unreleased]`; update `docs/implementation-status.md` |
| 7.4 | Optional ADR | ADR-0003 to document layout/decimal/observability decisions (ADRs 0001/0002 are immutable) |
| 7.5 | Release | Tag `v2.0.0`; preserve v1.x on a branch for non-migrating consumers |

## Execution order & checkpoints

Each checkpoint is a focused `git commit -s` (Conventional Commits + Signed-off-by).

1. **P1** DevOps & testing infrastructure → commit
2. **P2** Architecture & domain foundation (blocking) → tests → commit
3. **P3** Security & resilience → commit
4. **P4** REST client & trading → commit
5. **P5** Streaming engine (+ `sync.Pool` only if benchmarks justify) → race + fuzz → commit
6. **P6** Observability → commit
7. **P7** Docs/examples + `mkdocs build --strict` → commit
8. Full verification gates → push `origin` + `gitee` → **tag `v2.0.0`**

## Risks & mitigations

| Risk | Mitigation |
|---|---|
| Relocation is large/mechanical (209 endpoints) | Deprecated shims preserve compatibility; use `gofmt -r` + full test coverage; do in one pass |
| Decimal type changes break Go consumers | v2 migration guide; JSON wire format unchanged; strict wrapper errors on bad input |
| OTel dependency adds module weight | API-only dep is lean; no-op default means no telemetry unless opted in |
| Full-jitter behind a preset (defaults unchanged) | Non-breaking by design |
| `gosec` false positives | Configure exclusions per-finding; review each flagged item |
| Coverage gate unrealistic on the full repo | Measure baseline first; gate on touched/new code initially |
| CI matrix Windows + race detector | `GOTMPDIR` workaround step |
| Deprecated shims go stale | Maintain one release cycle, remove in v3 |
| `gofmt -l .` fails after relocation | Run `gofmt -w .` as part of the relocation step; verify before commit |

## Out of scope

- Broker FD US / Display Solution / broker HK surfaces are **environment-blocked** (404/403/401 in the HK sandbox), not SDK issues.
- `docs/runs/**` is excluded from the published site and must never be edited.
- ADRs 0001 and 0002 are immutable; supersede with new ADRs.
- Stray root `.exe` binaries (`brokerfd.exe`, `brokerfd-events.exe`, `options.exe`) are gitignored; optional cleanup.

## Verification gates (every checkpoint)

- `go build ./...`
- `go vet ./...`
- `gofmt -l .` must print nothing
- `go test -race -count=1 ./...`
- `golangci-lint run` (including gosec once enabled)
- Coverage gate (once P1.4 lands)
- Fuzz smoke runs
- `mkdocs build --strict` (after any doc changes)

Final step: push to `origin` and `gitee`, then tag `v2.0.0`.

## Reference

- Repository guide: `AGENTS.md`
- Build/test/lint: `Makefile`, `.golangci.yml`, `.github/workflows/ci.yml`
- Live sandbox tests: `.github/workflows/nightly-live.yml`
- Docs site: `mkdocs.yml`, `docs/`
- Official Webull docs: `https://developer.webull.hk/apis/docs` (HK)
