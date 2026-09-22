# Plan: HK Sandbox Probe 2026-09-22

## Goals
1. Probe HK sandbox for v0.10 viability (futures, options multi-leg strategies, Broker API HK)
2. Identify all TODO items confirmed/denied by live sandbox
3. Fix critical bugs surfaced during probe
4. Write findings to docs/runs/

## Probe Targets
- [x] `GetFuturesProductCodes` (HK) — confirm path and response shape
- [x] `GetFuturesInstruments` (HK) — check Unit field type (suspected numeric in live data)
- [x] Multi-leg option strategies via `PreviewOrder` (VERTICAL, STRADDLE, STRANGLE, etc.)
- [x] `GetOptionsExpirations`, `GetOptionChain` (HK) — path confirmation
- [x] Watchlist CRUD on HK sandbox
- [x] US-only endpoints (confirm 404)
- [x] Broker API HK — confirm scope requirements

## Bugs to Fix
- [ ] `FuturesInstrument.Unit` flexible type (string vs numeric)
- [ ] `client_order_id` length overflow in `examples/options-multi-leg/main.go:193`
- [ ] Update TODO markers with HK findings

## Output
- `docs/runs/2026-09-22-hk-sandbox-probe/plan.md` (this file)
- `docs/runs/2026-09-22-hk-sandbox-probe/todos.md`
- `docs/runs/2026-09-22-hk-sandbox-probe/report.md`
