# Report — v1.0 SDK Completion (v1-alpha-probe)

## Summary

Built the two largest missing Trading features — **multi-leg options orders**
(T8) and **futures order validation** (T9) — plus speculative **option
chain/expiration discovery** (T10), with tests and docs. Discovered that the
HK sandbox cannot verify any of the remaining unconfirmed endpoints, which
re-scoped several planned tasks to blocked/deferred.

## Shipped vs. Deferred vs. Blocked

| ID | Task | Status | Notes |
|----|------|--------|-------|
| T8 | Multi-leg options orders | **Shipped** | 10 new `OptionStrategy` values; `OrderRequest.Legs` now supports 2+ legs; duplicate/degenerate/cross-leg validation; `orderNotional` skips multi-leg |
| T9 | Futures order placement | **Shipped** | `validateFuturesRules` wired into the instrument dispatch; US/HK order-type matrix; QTY-only, whole-contract, DAY/GTC |
| T10 | Options chain discovery | **Shipped (speculative)** | `GetOptionExpirations` + `GetOptionChain`; endpoint not in published OpenAPI; `TODO(t10)` |
| T14 | Tests for new work | **Shipped** | Edge cases for T8/T9/T10; offline; race clean |
| T15 | Quality review | **Shipped** | 8 findings; 5 fixed (honesty markers, numeric strike compare, notional doc, CN message) |
| T16 | Docs sync | **Shipped** | README, CHANGELOG, AGENTS.md, docs/trading.md, docs/api.md, docs/market-data.md |
| T3 | Screener v2 refine | **Blocked** | HK sandbox 404 |
| T4 | Display Solution refine | **Blocked** | HK sandbox 404 |
| T5 | Instrument v3 + Logos refine | **Blocked** | HK sandbox 404 |
| T7 | API-lock pass | **Blocked** | depends on T3–T5 |
| T11 | Options MQTT stream | **Blocked** | live topic unconfirmed |
| T1 | Fund Data confirm | **Deferred** | US sandbox |
| T2 | Crypto Data confirm | **Deferred** | US sandbox |
| T6 | Broker FD HTTP paths | **Deferred** | US sandbox |
| T12 | Broker FD gRPC client | **Deferred** | US sandbox |
| T13 | HK/CN derivatives | **Deferred** | live probe |

## Key Discovery

`examples/probe/results/` shows every unconfirmed endpoint returns **404** on the
HK sandbox: Screener v2 (`/wlas/...`, `/market-data/screener/ng/query`), Fund
Data, Crypto Data, Corporate Actions (market-data host), Broker FD, Instrument
v3, and Logos. The official OpenAPI docs
(`developer.webull.hk/apis/llms.txt`) expose **only three** option endpoints
(ticks, snapshots, bars) — no option chain/expiration endpoint exists publicly.

Consequence: schema-refinement tasks cannot be completed or verified without a
**US sandbox account**. Rather than guess, they are marked blocked/deferred.

## Files Touched (uncommitted)

**Code**
- `trade/types.go` — multi-leg `OptionStrategy` values; futures enum note
- `trade/options.go` — multi-leg validation, `canonicalStrike`, provisional markers
- `trade/orders.go` — instrument-dispatched validation; `orderNotional` multi-leg skip
- `trade/rules.go` — `marketFuturesOrderTypes`, `validateFuturesRules`, `isPositiveInteger`
- `trade/option.go` — notional-cap caveat
- `data/options.go` — `GetOptionExpirations`, `GetOptionChain`, `OptionContract`, etc.

**Tests**
- `trade/multileg_test.go` (new)
- `trade/futures_order_test.go` (new)
- `trade/options_test.go`, `trade/combo_test.go` (updated)

**Docs**
- `README.md`, `CHANGELOG.md`, `AGENTS.md`, `docs/api.md`, `docs/market-data.md`, `docs/trading.md`

**Run artifacts**
- `docs/runs/2026-09-20-v1-alpha-probe/{plan,todos,report,next-phase}.md`

## Risks

| Risk | Severity | Mitigation |
|------|----------|-----------|
| Multi-leg strategy wire values provisional | High | `TODO(t8)` markers; isolated to constants + `OptionStrategy.valid()` |
| Futures order-type matrix provisional | Medium | `TODO(t9)`; single matrix to update |
| Option chain endpoints speculative (may return empty) | Medium | `TODO(t10)`; strongly worded docs |
| `WithMaxOrderNotional` bypassable by adding a leg | Low | Documented on the option; quantity cap still applies |
| Top-level `Market` not validated for OPTION orders | Low | Pre-existing; flagged for follow-up |

## Verification

```
go build ./...                          → OK
go vet ./...                            → OK
gofmt -l .                              → clean
go test ./... -count=1                  → all packages ok
go test -race -count=1 ./trade/... ./data/...  → ok
golangci-lint run ./...                 → 0 issues
mkdocs build --strict                   → exit 0; site/ excludes docs/runs/
```

## Follow-ups

1. **Obtain a US sandbox account** — unblocks T1, T2, T3, T4, T5, T6, T7, T11, T12.
2. Confirm `TODO(t8)`/`TODO(t9)`/`TODO(t10)` against the live API.
3. Decide whether to keep the speculative option-chain functions or remove them
   until the endpoint is publicly documented.
4. Reconcile release-status lag in README/SECURITY.md with the released v0.5.0.
5. Add top-level market validation for OPTION orders (review finding 5).
