# Findings — Phase 0 + Live Probe

## HK Sandbox Live Probe Results (2026-09-19)

### Verified Working on HK Sandbox
| Endpoint | Status | Notes |
|----------|--------|-------|
| `GET /market-data/stocks/bars/list` (AAPL) | **200 OK** | Returns real OHLCV data |
| Token auth | **200 OK** | HMAC-SHA1 signing works |

### Corporate Actions
| Path | Status | Notes |
|------|--------|-------|
| `/market-data/instruments/stocks/corporate-actions/list` | **404** | Path not found in HK sandbox |
| `/market-data/instruments/stocks/corporate-actions/market` | **404** | Path not found in HK sandbox |

**Conclusion**: Corporate Actions paths from HK API reference do not exist in the HK sandbox.
May require US sandbox credentials or Display Solution authentication.

### Fund / ETF Data
| Path | Status | Notes |
|------|--------|-------|
| `/market-data/fund/AAPL/nav` | 404 | Not found |
| `/market-data/fund/list` | 404 | Not found |
| `/market-data/etf/list` | 404 | Not found |

**Conclusion**: Fund/ETF data endpoints do not exist in HK sandbox market data API.

### Crypto Data
| Path | Status | Notes |
|------|--------|-------|
| `/market-data/crypto/bars` | **400** | Path exists but params invalid |
| `/market-data/crypto/BTCUSD/bars` | 404 | Path not found |
| `/market-data/crypto/ETHUSD/bars` | 404 | Path not found |

**Conclusion**: `/market-data/crypto/bars` is a real path but requires correct params.
Category `CRYPTO` may be wrong — the HK sandbox may use a different category value.

### Screener v2
| Path | Status | Notes |
|------|--------|-------|
| `/wlas/screener/ng/query` | 404 | Not found |
| `/market-data/screener/ng/query` | 404 | Not found |

**Conclusion**: Screener v2 paths do not exist in HK sandbox.

### Fund Data, Broker FD
All tried paths return 404 — these are not available in the HK sandbox market data API.

## Key Insight

The HK sandbox (`api.sandbox.webull.hk`) supports **US stock market data** (AAPL bars work with `category=US`)
but does NOT support Display Solution endpoints (Corporate Actions, Crypto, Screener v2).

These Display Solution endpoints appear to require:
- Either US sandbox credentials (`api.sandbox.webull.com`)
- Or Display Solution authentication (separate host: `us-global-openapi.uat.webullbroker.com`)

## T0.2–T0.7 Status Update

| ID | Task | Status | Notes |
|----|------|--------|-------|
| T0.2 | Corporate Actions | **Stub only** | Path confirmed from API ref but 404 in sandbox |
| T0.3 | Fund Data | **Stub only** | All paths 404 in sandbox |
| T0.4 | Crypto Data | **Partial** | Path `/market-data/crypto/bars` exists (400), params wrong |
| T0.5 | Screener v2 | **Stub only** | All paths 404 in sandbox |
| T0.6 | Broker FD HTTP | **Stub only** | All paths 404 |
| T0.7 | Broker FD gRPC | **Stub only** | Proto not live-probed |

## Probe Binary

The Go probe is at `examples/probe/main.go`. Run with:

```sh
WEBULL_APP_KEY=xxx WEBULL_APP_SECRET=yyy go run examples/probe/main.go
```

Build first: `go build -o examples/probe/probe.exe examples/probe/main.go`
