# Next Phase — v0.5 Live Probe

## What was completed

Live probe run against HK sandbox (`api.sandbox.webull.hk`) using Go probe (`examples/probe/main.go`):

- **AAPL bars**: 200 ✅ — HK sandbox supports US stocks, auth + signing works
- **Corporate Actions**: 404 — paths from HK API ref don't exist in HK sandbox
- **Fund/ETF data**: 404 — doesn't exist in HK sandbox
- **Crypto `/market-data/crypto/bars`**: **400** — path exists, params wrong (progress!)
- **Crypto BTC/ETH in path**: 404 — different path structure
- **Screener v2**: 404 — paths don't exist in HK sandbox
- **Broker FD**: 404 — not available in HK sandbox market data API

## Key insight

HK sandbox supports US stocks (AAPL works with `category=US`) but does NOT support
Display Solution endpoints. Those require either:
1. US sandbox credentials (`api.sandbox.webull.com`)
2. Display Solution host (`us-global-openapi.uat.webullbroker.com`)

## Remaining gaps

| Gap | What we need | Priority |
|-----|--------------|----------|
| Corporate Actions schema | US sandbox or Display Solution host | HIGH |
| Crypto params | US sandbox to find correct `category`/params for `/market-data/crypto/bars` | HIGH |
| Screener v2 path | US sandbox or Display Solution host | MEDIUM |
| Fund data | May not be in Webull API at all | LOW |
| Broker FD | HK Broker API host not tested (different auth) | MEDIUM |

## Recommended next phase

### Option A: Get US sandbox credentials (HIGH value)
Obtain US sandbox test account credentials. The Display Solution endpoints
(Corporate Actions, Crypto, Screener v2) likely work with US sandbox.
Effort: S. Value: confirms all v0.5 schemas.

### Option B: Ship v0.5 as-is (no new coding)
Release v0.5 with current best-effort stubs. All stubs compile, have `Extra map[string]string`,
and are forward-compatible with real data. Open issues for schema confirmation.
Effort: XS (just tag + changelog). Value: ships what we have.

### Option C: Investigate crypto 400 (MEDIUM value)
The `/market-data/crypto/bars` path exists (400 not 404). Try different `category`
values and param names to find the correct crypto data format.
Effort: S. Risk: may still need US sandbox.

## Open questions for the human

1. Do you have access to US sandbox credentials (`api.sandbox.webull.com`)?
2. Do you want to ship v0.5 as-is with best-effort stubs, or wait for US sandbox probing?
3. Should we investigate the crypto 400 path further (try different category values)?
