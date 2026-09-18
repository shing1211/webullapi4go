# Report — v0.5 Phase 2: Fund Data, Crypto Data, Screener v2

## Shipped vs. Deferred

| Item | Status | Notes |
|------|--------|-------|
| `data/fund_data.go` | Shipped | 4 endpoints, best-effort, `map[string]string` rows |
| `data/crypto_data.go` | Shipped | 4 endpoints, best-effort, `map[string]string` rows |
| `data/screener_v2.go` | Shipped | 1 POST endpoint, best-effort, reuses `ScreenerStock` |
| `docs/api.md` update | Shipped | New endpoints listed under `data` section |
| Live probe (Fund/Crypto/Screener) | Deferred | Blocked — needs `WEBULL_APP_KEY` + `WEBULL_APP_SECRET` |
| Unit tests for new stubs | Deferred | Can add before or after live probe |

## Risks

- All endpoint paths are guessed — may differ from actual API
- Schema fields unknown — `map[string]string` on all structs
- Screener v2 uses POST body — signing behavior unconfirmed (same HMAC-SHA1 as other market data?)

## Follow-ups

1. **Live probe needed** — confirm paths and refine schemas
2. **Screener v2 signing** — verify POST body signing matches GET endpoint behavior
3. **Add tests** — `fund_data_test.go`, `crypto_data_test.go`, `screener_v2_test.go`

## Commits

| Commit | Message |
|--------|---------|
| `f41b803` | feat: v0.5 phase 2 — Fund Data, Crypto Data, Screener v2 best-effort stubs |
| `e1c20bf` | docs: update runs/index.md for v05-fund-crypto release |

Both pushed to GitHub (`origin/main`) and Gitee (`gitee/main`).
