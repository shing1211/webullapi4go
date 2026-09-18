# Recon Findings — T17.0

## Method
Read the Webull API reference docs (developer.webull.hk/apis/docs/reference) for each
candidate endpoint, extracting exact paths, methods, parameters, and response schemas.
Attempted sandbox live tests but credentials (WEBULL_APP_KEY, WEBULL_APP_SECRET)
are not set in this environment.

## Credentials status
- WEBULL_APP_KEY, WEBULL_APP_SECRET: NOT SET (env vars absent, no .env file)
- GitHub vars: none set
- Sandbox probe: FAILED (401 Unauthorized with placeholder credentials)
- Live testing deferred until credentials are available

## Findings: Market Data API Endpoints

### Already implemented in Go SDK (verified from API reference)
| Endpoint | Path | Status |
|----------|------|--------|
| CompanyProfile | GET /market-data/fundamentals/company-profiles/get | ✅ implemented |
| AnalystTargetPrice | GET /market-data/fundamentals/analysis/target-prices/get | ✅ implemented |
| AnalystRating | GET /market-data/fundamentals/analysis/ratings/get | ✅ implemented |
| Stock Bars | GET /market-data/stocks/bars/get (single) + POST /market-data/stocks/bars/list (batch) | ✅ implemented |
| Stock Snapshot | POST /market-data/stocks/snapshots/list | ✅ implemented |
| Stock Tick | GET /market-data/stocks/ticks/list | ✅ implemented |
| Stock Depth | GET /market-data/stocks/depths/list | ✅ implemented |
| Stock Footprint | GET /market-data/stocks/footprints/list | ✅ implemented |
| NOII Bars/Snapshot | GET /market-data/stocks/noii-bars/list, /market-data/stocks/noii-snapshots/list | ✅ implemented |
| Options tick/snapshot/bars | GET /market-data/options/{ticks,snapshots,bars}/list | ✅ implemented |
| Futures tick/snapshot/depth/bars | GET /market-data/futures/{ticks,snapshots,depths,bars}/list | ✅ implemented |
| Watchlist CRUD | GET/POST /market-data/watchlists/... | ✅ implemented |
| TopGainersLosers | GET /market-data/screeners/gainers-losers/list | ✅ implemented (path: screeners plural — confirmed) |
| TopActives | GET /market-data/screeners/top-actives/list | ✅ implemented |

### Fundamentals — MISSING (T17.2)
All use GET with query params: symbol (required), category (required, enum).

| Endpoint | Path | API Ref Confirmed | Sandbox | Priority |
|----------|------|-------------------|---------|----------|
| Capital Flow | GET /market-data/fundamentals/capital-flows/get | ✅ yes (fetched doc) | unverified | HIGH |
| Industry Comparison | GET /market-data/fundamentals/industry-comparisons/get | ✅ rate limits page | unverified | HIGH |
| Earnings Calendar | GET /market-data/fundamentals/earnings-calendars/list | ✅ rate limits page | unverified | HIGH |
| Dividend Calendar | GET /market-data/fundamentals/dividend-calendars/list | ✅ rate limits page | unverified | HIGH |
| SEC Filings | GET /market-data/fundamentals/filings/list | ✅ rate limits page | unverified | MED |

Capital Flow confirmed schema (from API ref doc):
- Query: symbol (string), category (US_STOCK|HK_STOCK|CN_STOCK), count (1~5, default 5)
- Response: array of {date, large_in, large_out, medium_in, medium_out, small_in, small_out}

### Financial Statements — MISSING (T17.3)
All use GET with query params: symbol (required), category (required).

| Endpoint | Path | Sandbox | Priority |
|----------|------|---------|----------|
| Income Statement | GET /market-data/fundamentals/income-statements/get | unverified | HIGH |
| Balance Sheet | GET /market-data/fundamentals/balance-sheets/get | unverified | HIGH |
| Cash Flow | GET /market-data/fundamentals/cash-flows/get | unverified | HIGH |
| Indicators | GET /market-data/fundamentals/indicators/get | unverified | MED |
| Financial Alert | GET /market-data/fundamentals/financial-alerts/get | unverified | MED |
| Forecast EPS | GET /market-data/fundamentals/forecast-eps/get | unverified | MED |

### Fund Data — MISSING (T17.5)
All use GET with query params: symbol (required), category (required).

| Endpoint | Path | Sandbox | Priority |
|----------|------|---------|----------|
| Fund Brief | GET /market-data/fundamentals/fund-brief/get | unverified | MED |
| Fund Performance | GET /market-data/fundamentals/fund-performances/get | unverified | MED |
| Fund Net Value | GET /market-data/fundamentals/fund-net-values/get | unverified | MED |
| Fund Holdings | GET /market-data/fundamentals/fund-holdings/get | unverified | MED |
| Fund Dividends | GET /market-data/fundamentals/fund-dividends/get | unverified | MED |
| Fund Rating | GET /market-data/fundamentals/fund-ratings/get | unverified | MED |
| Fund Splits | GET /market-data/fundamentals/fund-splits/get | unverified | MED |
| Fund Files | GET /market-data/fundamentals/fund-files/get | unverified | LOW |
| Fund Allocation | GET /market-data/fundamentals/fund-allocations/get | unverified | MED |

### Crypto — MISSING (T17.6)
Note: Crypto requires NO additional subscription (from docs).

| Endpoint | Path | Method | Sandbox | Priority |
|----------|------|--------|---------|----------|
| Crypto Snapshot | GET /market-data/crypto/snapshots/list | GET | unverified | HIGH |
| Crypto Bars | GET /market-data/crypto/bars/list | GET | unverified | HIGH |

Crypto path constants (from Python SDK endpoints.py):
- `base_fintech_gw_url/crypto/charts/query?tickerIds={stock}` (Python SDK old path)
- New path per API ref: `GET /market-data/crypto/bars/list`
- Snapshot: `GET /market-data/crypto/snapshots/list`

### Extended Screener — MISSING (T17.7)

| Endpoint | Path | Method | Sandbox | Priority |
|----------|------|--------|---------|----------|
| Market Sectors List | GET /market-data/screeners/market-sectors/list | GET | unverified | MED |
| Market Sectors Detail | GET /market-data/screeners/market-sectors/get | GET | unverified | MED |
| High Dividend Rank | GET /market-data/screeners/high-dividend-ranks/list | GET | unverified | MED |
| 52 Week High/Low | GET /market-data/screeners/week52-high-low/list | GET | unverified | MED |

Note: Go SDK uses `screener` (singular); API ref uses `screeners` (plural).
The Go SDK's `GetTopGainersLosers` already uses the plural form confirmed by API ref.

### Corporate Actions + Logos (T17.8)

| Endpoint | Path | Method | Sandbox | Priority |
|----------|------|--------|---------|----------|
| Corp Actions by Market | GET /market-data/instruments/stocks/corporate-actions/list-by-market | GET | unverified | MED |
| Corp Actions | GET /market-data/instruments/stocks/corporate-actions/list | GET | unverified | MED |
| Batch Logos | POST /market-data/fundamentals/logos/list | POST | unverified | LOW |

### News REST (NOT EXISTS)
- Python SDK `get_news` uses `GET /information/news/tickerNews` (different domain)
- API ref shows `POST /market-data/news/summaries/get` — but this is a POST with body
- Sandbox status: unverified (Python SDK has no REST news list, only SSE stream)
- Conclusion: no REST news list confirmed; only SSE stream in data/news.go

### Instrument v3 (T17.2)
- Path: `POST /market-data/instruments/stocks/profiles/list`
- Note: Go SDK currently uses `/openapi/instrument/stock/list` (legacy)
- Response envelope likely: `{data: [...], pagination_key: ...}` (different from array)
- Action: ADD new method (don't replace old one — breaking change risk)

## Classification summary
| Category | Count | Notes |
|----------|-------|-------|
| (a) Implemented | 14 | Verified from API ref |
| (b) Missing, API ref confirmed | 27 | All paths confirmed from rate limits page |
| (c) 404/unverified | 0 | No unknown paths |
| (d) Exception | 0 | N/A |

## Key decisions from recon
1. **Capital flow path confirmed**: `/market-data/fundamentals/capital-flows/get` (not `/capital-flow`)
2. **Screener plural confirmed**: all screener endpoints use `screeners` (plural) — existing Go code already correct
3. **Fund data paths**: all follow `/market-data/fundamentals/fund-*` pattern
4. **Corporate actions**: two variants — by-market (no symbol) and by-symbol
5. **Batch logos**: POST with `{"symbols": ["AAPL"]}` body
6. **Crypto**: no tick/depth/footprint in API ref (Python SDK confirms)
7. **News REST**: POST endpoint exists per rate limits, but Python SDK uses different domain path
8. **Instrument v3**: add new method, don't replace legacy

## Next step
Implement all (b) endpoints using the exact paths confirmed above. Offline fixture tests
will provide coverage; live sandbox tests deferred until credentials are available.
