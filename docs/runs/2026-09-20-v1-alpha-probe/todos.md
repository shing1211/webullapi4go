# Todos — v1.0 SDK Completion

Single source of truth for run `2026-09-20-v1-alpha-probe`.

## Discovery (2026-09-20)

Prior probe artifacts in `examples/probe/results/` show **every unconfirmed
endpoint returns 404 on the HK sandbox**: Screener v2 (`/wlas/...`,
`/market-data/screener/ng/query`), Fund Data, Crypto Data, Corporate Actions
(market-data host), Broker FD. Instrument v3 and Logos also route through the
market-data host.

Consequence: T3, T4, T5, T7, T11 cannot be verified without US sandbox
credentials. They are marked **blocked**, not implemented with guessed schemas.

## Active Work

| ID | Task | Role | Status | Depends On | Acceptance |
|----|------|------|--------|------------|------------|
| T8 | Multi-leg options orders | backend | done | — | `go test ./trade/...` passes; new `OptionStrategy` values; cross-leg validation |
| T9 | Futures order placement | backend | done | T8 | `validateFuturesRules` wired into `OrderRequest.validate`; `go test ./trade/...` passes |
| T10 | Options chain discovery | backend | done | T9 | `GetOptionChain` added with `// TODO: live probe`; `go build` passes |
| T14 | Tests for new work | tester | done | T8,T9,T10 | Offline table-driven tests pass; race clean |
| T15 | Quality review | reviewer | done | T14 | 8 findings; 5 fixed |
| T16 | Docs sync | docs | done | T14 | `mkdocs build --strict` passes; no stale refs |

## Blocked — need US sandbox credentials

| ID | Task | Blocker |
|----|------|---------|
| T3 | Screener v2 refine | HK 404 |
| T4 | Display Solution refine | HK 404 (market-data host) |
| T5 | Instrument v3 + Logos refine | HK 404 |
| T7 | API-lock pass | depends on T3–T5 |
| T11 | Options MQTT stream | live topic unconfirmed |

## Note on T10

The official Webull OpenAPI docs (`developer.webull.hk/apis/llms.txt`) expose
only three option endpoints (ticks, snapshots, bars) — there is **no published
option chain/expiration endpoint**. `GetOptionChain`/`GetOptionExpirations`
were implemented against placeholder paths to satisfy the requested surface and
are marked `TODO(t10): unconfirmed path`. They may return empty results until
verified; treat as provisional.

## Deferred — out of scope this run

| ID | Task | Blocker |
|----|------|---------|
| T1 | Fund Data confirm | US sandbox |
| T2 | Crypto Data confirm | US sandbox |
| T6 | Broker FD HTTP paths | US sandbox |
| T12 | Broker FD gRPC client | US sandbox |
| T13 | HK/CN derivatives | live probe |

Statuses: todo · doing · blocked · review · done · deferred
