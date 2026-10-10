# Plan — v1.0 SDK Completion

## Goal

Cover all known Webull API functions. Build everything verifiable with the HK
sandbox now; mark US-only items for later verification once US sandbox
credentials are available.

## Config

| Field | Value |
|-------|-------|
| Project | webullapi4go |
| Repo root | `D:\github\webullapi4go` |
| Stack | Go, Protobuf, MQTT, gRPC |
| VCS | GitHub + Gitee, `main` |
| Verification | `go build ./... && go vet ./... && gofmt -l . && go test ./... && golangci-lint run ./... && mkdocs build --strict` |
| Constraints | Do not touch `.github/workflows/ci.yml`, `.golangci.yml`, `LICENSE` |
| Credentials | HK sandbox available; US sandbox deferred |

## Scope — Task List

### Group A — Build now (HK-verifiable or pure logic)

| ID | Task | Role | Size | Acceptance |
|----|------|------|------|------------|
| T3 | Screener v2 refine | backend | M | `GetScreenerV2` request/response fields concrete; `go test ./data/... -run Screener` passes |
| T4 | Display Solution refine | backend | M | Corporate actions DTOs have concrete fields; `go test ./data/... -run Corporate` passes |
| T5 | Instrument v3 + Logos refine | backend | S | DTO field types confirmed; `go test ./data/... -run "Profile\|Logo"` passes |
| T7 | API-lock pass | backend | M | No guessed types remain in HK-verifiable files |
| T8 | Multi-leg options orders | backend | L | New `OptionStrategy` values; multi-leg validation; `go test ./trade/...` passes |
| T9 | Futures order placement | backend | M | `validateFuturesRules` added and wired; `go test ./trade/...` passes |
| T14 | Tests for new work | tester | M | Table-driven offline tests for T3–T9; no network |

### Group B — Build now, mark `// TODO: live probe`

| ID | Task | Role | Size | Acceptance |
|----|------|------|------|------------|
| T10 | Options chain discovery | backend | M | `GetOptionChain` added with documented TODO; `go build` passes |
| T11 | Options MQTT stream | backend | M | Options push subscription added; `go test ./stream/...` passes |

### Group C — Deferred until US sandbox

| ID | Task | Size | Blocker |
|----|------|------|---------|
| T1 | Fund Data confirm | M | HK 404 |
| T2 | Crypto Data confirm | M | HK 417 |
| T6 | Broker FD HTTP paths | M | HK 404 |
| T12 | Broker FD gRPC streaming client | L | Proto unconfirmed |
| T13 | HK/CN derivatives expansion | L | Scope unconfirmed |

### Group D — Close-out

| ID | Task | Role | Size |
|----|------|------|------|
| T15 | v1.0 API freeze pass | reviewer | S |
| T16 | Docs sync | docs | M |

## Approach

Two-phase: (1) refine and verify existing stubs against the HK sandbox and
offline tests; (2) build the two largest missing features — multi-leg option
orders and futures order validation — as pure logic with offline tests, since
both are verifiable without credentials. Deferred items are documented in
`next-phase.md` rather than half-implemented.

Rejected alternatives:

- **Implement all 16 tasks now** — US-only paths would be guesswork with no way
  to verify; high risk of shipping wrong schemas.
- **Wait for US credentials before any work** — unnecessarily blocks 9 tasks
  that are fully verifiable today.

## Key Design Decisions

| Decision | Rationale |
|----------|-----------|
| Reuse `OrderLeg` + `OrderRequest.Legs` for multi-leg | Wire shape already matches; only strategy enum and validation need extending |
| Add `validateFuturesRules` separate from equity/option | `InstrumentTypeFutures` is defined but unvalidated today |
| Reconcile `"FUTURES"` (trade) vs `"FUTURE"` (pkg/types) | Existing enum mismatch before v1.0 lock |
| Mark options-chain path as TODO | No existing endpoint; cannot confirm without live probe |

## Risks

| Risk | Mitigation |
|------|-----------|
| Multi-leg order wire format differs from single-leg | Ship validation + structure; sandbox integration test gated by env |
| US-only items stay unverified | `// TODO: US sandbox` markers + next-phase.md |
| Options chain path wrong | Isolated behind one function |
| Scope creep in T13 | Keep deferred |

## Verification

```
go build ./...
go vet ./...
gofmt -l .
go test ./...
golangci-lint run ./...
mkdocs build --strict
```

All must pass. Sandbox integration tests gated by `WEBULL_SANDBOX=1` must skip
by default.

## Actuals vs Plan

| ID | Planned | Actual |
|----|---------|--------|
| T8 | Multi-leg options | Done — 10 strategies, validation, notional skip, tests |
| T9 | Futures validation | Done — per-market matrix, QTY/whole-contract/DAY-GTC, tests |
| T10 | Options chain | Done but **speculative** — endpoint not in published OpenAPI |
| T14 | Tests | Done — edge cases + race clean |
| T15 | v1.0 API freeze | **Reframed** as quality review (freeze impossible while unverified); 5 findings fixed |
| T16 | Docs sync | Done — 6 files updated; mkdocs strict passes |
| T3,T4,T5,T7,T11 | Build now | **Blocked** — HK sandbox 404; need US credentials |
| T1,T2,T6,T12,T13 | Deferred | Confirmed deferred |

Net change: 5 tasks shipped, 1 shipped speculative, 5 blocked, 5 deferred.
The plan's assumption that T3–T5 were HK-verifiable was **wrong** — all probe
artifacts show 404. This was discovered at execution start and drove the
re-scoping.
