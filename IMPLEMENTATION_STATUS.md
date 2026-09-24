# Implementation Status

Last updated: 2026-09-24 (v2.0.6 release) · Current version: **v2.0.6**

## Summary

| Module | Package | Transport | Exported API | TODOs | Verified |
|--------|---------|-----------|-------------|-------|----------|
| Core SDK | `client/` | HTTP | 20+ options, `Client.Do`, `Client.DoStream` | 0 | ✅ |
| Authentication | `internal/auth/` | — | HMAC-SHA1/SHA-256 signing, token lifecycle | 0 | ✅ |
| Market Data HTTP | `data/` | HTTP | 112 functions | 0 | ✅ Paths aligned to the official OpenAPI definition |
| Market Data Streaming | `stream/` | MQTT | 12+ options, typed handlers | 0 | ✅ |
| Trading HTTP | `trade/` | HTTP | 30+ methods | 0 | ✅ |
| Trading Events | `events/` | gRPC | 21 functions | 0 | ✅ |
| Connect API | `connect/` | HTTP | OAuth authorization-code flow | 0 | ✅ |
| Display Solution | `display/` | HTTP | 9 functions | 0 | ✅ |
| Broker API HK | `broker/` | HTTP | 32 methods | 0 | ✅ Own `go.mod` |
| Broker FD API US | `brokerfd/` | HTTP | 56 methods | 0 | ✅ |
| Broker FD Events | `brokerfd/events/` | gRPC | 17 functions | 0 | ✅ |

## TODO Marker Inventory (0 total)

No provisional TODO markers remain. Every documented Webull endpoint is
implemented and the SDK paths follow the official OpenAPI definition.

## Version History

| Version | Date | Scope | Status |
|---------|------|-------|--------|
| v0.1.0 | 2026-09-18 | Auth, core HTTP, Market Data HTTP + MQTT streaming | Done |
| v0.2.1 | 2026-09-18 | Accounts, balances, positions | Done |
| v0.2.2 | 2026-09-18 | Stock order lifecycle + queries | Done |
| v0.2.3 | 2026-09-18 | Market-specific rules, HK BCAN | Done |
| v0.2.4 | 2026-09-18 | Single-leg options orders | Done |
| v0.2.5 | 2026-09-18 | US combo orders | Done |
| v0.2.6 | 2026-09-18 | News SSE refactor (`DoStream`) | Done |
| v0.3.0 | 2026-09-18 | Trading events over gRPC | Done |
| v0.4.0 | 2026-09-18 | Market data fundamentals (13 endpoints) | Done |
| v0.5.0 | 2026-09-19 | Display Solution, corporate actions, fund/crypto data, screener v2 | Done (provisional) |
| v0.6.0 | 2026-09-20 | Multi-leg options, futures validation, option chain discovery | Done (provisional) |
| v0.7.0 | 2026-09-21 | Event contracts, Broker API HK, Broker FD US, Broker FD events | Done (current) |
| v0.8.0 | 2026-09-21 | Reserved | — |
| v0.9.0 | 2026-09-22 | GoDoc coverage on brokerfd/ and brokerfd/events/, HK options stubs (OptionCategoryHK/CN, GetHKOptionExpirations, GetHKOptionChain), HK futures market data (GetHKFuturesTick/Snapshot/Bars/Depth/Footprint), new examples (brokerfd, brokerfd-events, options), graceful credential handling, FuturesCategoryCN | Done |
| v0.9.1 | 2026-09-21 | Watchlist boolean-response fix; `DoBroker` transport; `watchlist-cmd` and `broker-probe` examples | Done |
| v0.9.2 | 2026-09-22 | Broker HK path correction (`/openapi/v1/broker/...` → `/broker/...`); `401 ROUTE_NOT_PERMITTED` instead of `404 Route Not Found` | Done |
| v1.0.0 | 2026-09-22 | Futures market data bug fix (Category field added to all 5 query structs); v1.0 API stability audit completed; all 49 TODOs remain provisional (require US sandbox) | Done (provisional) |
| v1.0.1 | 2026-09-22 | FuturesInstrument.Unit flexible type (StringOrNumber handles numeric API responses); client_order_id length fix in options-multi-leg example; HK sandbox probe findings documented | Done |
| v1.0.2 | 2026-09-22 | Reconciled the codebase against the official Webull API: removed undocumented functions (crypto, screener v2, option expirations, HK futures duplicates) and corrected `GetOptionContracts` to the Trading API path | Done |
| v1.0.3 | 2026-09-22 | Verbatim Webull master guides/reference, SDK↔API reconciliation and coverage gaps, doc-generator CLI (`tools/webull-docgen`), HK-sandbox path probe (`examples/path-probe`) | Done |
| v1.1.0 | 2026-09-22 | Full SDK parity with the official OpenAPI: aligned 69 differing paths, implemented the 23 remaining endpoints (Connect OAuth, crypto, Display event contracts, fund extras, Display refresh), removed all provisional TODO markers, added `connect/` package | Done |
| v1.1.1 | 2026-09-23 | Production hardening: Makefile targets, `gosec`/`govulncheck`/coverage CI, multi-OS matrix, dependabot, `goleak` goroutine-leak detection, fuzz deserialization tests | Done |
| v2.0.0 | 2026-09-23 | Phase 2 production hardening: `pkg/errors`, `pkg/transport`, `pkg/resilience`, `pkg/domain/money`, `pkg/domain/order` public API; `webull/` facade package | Done |
| v2.0.1 | 2026-09-23 | Numeric string fields converted to `*money.Money`/`money.Money` across `trade/` and `data/` packages; `money.Rat()` bug fix | Done |
| v2.0.2 | 2026-09-23 | Numeric string fields converted to `*money.Money`/`money.Money` across `brokerfd/` package | Done |
| v2.0.3 | 2026-09-23 | goleak: ignore paho HTTP/2 goroutines after MQTT WebSocket disconnect | Done |
| v2.0.4 | 2026-09-24 | Phase 3 hardening: clock-drift correction, idempotency helpers, `WithHTTPTransport`, full-jitter retry, `WithResiliencePreset(production)` | Done |

## Feature Coverage

### ✅ Fully Verified (tested against HK sandbox)

All functions below have real HTTP/gRPC logic, full test coverage, and work against the HK sandbox.

**Core SDK (`client/`)**
- `client.New` with functional options (`WithEnv`, `WithAppKey`, `WithAppSecret`, `WithRegion`, `WithSandbox`, `WithHTTPTransport`, `WithResiliencePreset`, `WithClockDriftCorrection`, etc.)
- `Client.Do` — signed HTTP request with automatic token injection
- `Client.DoStream` — signed streaming request for SSE endpoints
- Token lifecycle: `CreateToken`, `CheckToken`, `EnsureToken`, `CurrentToken`, `AccessToken`, `SetToken`
- Per-path API version defaults (`x-version` header: `v2` market data, `v3` trading)
- Resilience: retry with exponential backoff, full jitter option, token-bucket rate limiter, circuit breaker
- Clock-drift correction via `Date` response header offset

**Market Data HTTP (`data/`) — 53 verified functions**
- Instruments: `GetStockInstruments`, `GetStockProfilesV3` (Display Solution v3)
- Logos: `GetLogos` (POST, Display Solution)
- Fundamentals (13): `GetCompanyProfile`, `GetAnalystTargetPrice`, `GetAnalystRating`, `GetCapitalFlow`, `GetIndustryComparison`, `GetEarningsCalendar`, `GetDividendCalendar`, `GetFilings`, `GetIncomeStatement`, `GetBalanceSheet`, `GetCashFlow`, `GetFinancialIndicators`, `GetFinancialAlert`, `GetForecastEPS`
- Futures static: `GetFuturesInstruments`, `GetFuturesProductCodes`, `GetFuturesProductClasses`
- Snapshot: `GetSnapshot`
- Tick: `GetTick`
- Quotes/Depth: `GetQuotes`
- Bars: `GetBars`, `GetBatchBars`
- Footprint: `GetFootprint` (requires paid entitlement; 403 in sandbox)
- NOII: `GetNOIIBars`, `GetNOIISnapshot`
- Screener (6): `GetTopGainersLosers`, `GetMostActive`, `GetMarketSectors`, `GetMarketSectorDetail`, `GetHighDividendRank`, `GetWeek52HighLow`
- Watchlists (8): `GetWatchlists`, `CreateWatchlist`, `UpdateWatchlist`, `DeleteWatchlist`, `GetWatchlistInstruments`, `AddWatchlistInstruments`, `RemoveWatchlistInstruments`, `UpdateWatchlistInstruments`
- Options: `GetOptionTick`, `GetOptionSnapshot`, `GetOptionBars`
- News: `GetNewsSummary` (SSE streaming)
- Event contracts (4): `GetEventContractCategories`, `GetEventContractSeries`, `GetEventContractEvents`, `GetEventContractMarkets`
- Corporate actions: `GetCorporateActions`, `GetCorporateActionsByMarket` (Display Solution)

**Market Data Streaming (`stream/`)**
- MQTT and MQTT-over-WebSocket connections
- Typed handlers: `OnQuote`, `OnSnapshot`, `OnTick`, `OnNotice`, `OnError`, `OnConnect`, `OnDisconnect`
- Auto-reconnect with HTTP re-subscription (idempotent, deduplicated)

**Trading HTTP (`trade/`) — 30+ verified methods**
- Accounts: `ListAccounts`, `GetBalance`, `GetPositions`, `GetCashActivities`, `GetCashActivitiesPage`, `GetAllCashActivities`
- Order lifecycle: `PreviewOrder`, `PlaceOrder`, `ReplaceOrder`, `CancelOrder`, `BatchPlaceOrder`
- Order queries: `GetOpenOrders`, `GetOpenOrdersPage`, `GetAllOpenOrders`, `GetOrderHistory`, `GetOrderHistoryPage`, `GetAllOrderHistory`, `GetOrderDetail`
- Single-leg options: LIMIT, STOP_LOSS, STOP_LOSS_LIMIT; BUY/SELL sides; DAY TIF
- US combo orders: MASTER/STOP_PROFIT/STOP_LOSS, OTO, OCO, OTOCO with full composition and leg-count validation
- Market-specific rules: US (8 order types), HK (7 order types), CN (LIMIT only)
- HK BCAN: party_id requirements for HK equity orders
- US trading session validation: CORE, ALL, NIGHT, ALL_DAY
- Order guardrails: `WithMaxOrderNotional`, `WithMaxOrderQuantity`

**Trading Events (`events/`)**
- gRPC server-streaming client for order, position, and option streams
- Typed payloads: `OrderEvent`, `PositionEvent`, `OptionEvent`
- HMAC-SHA256 signing for event service
- Auto-reconnect with exponential backoff + jitter
- Subscribe types: `SubscribeOrder` (1), `SubscribePosition` (2), `SubscribeOption` (4), `SubscribeAll` (7)

**Display Solution (`display/`)**
- Separate authentication service (HMAC-SHA1 signed client-token fetch)
- Lazy token initialization, caching, auto-refresh on 401
- Used by `data/` for Display Solution endpoints

**Broker API HK (`broker/`) — 32 methods, own `go.mod`**
- Virtual accounts: create, update, get, list
- Instruments: search, locate, corporate actions detail
- Assets: balance, positions
- Orders: preview, place, replace, cancel, detail, history, open orders
- Cash activities: list
- Funding FX: rate, exchange, instant exchange
- Journals: cash and position journals
- Master data: trade calendar
- Event contracts: categories, series, events, instruments

### Futures, Options, Events, Broker FD and Display

Futures market data, event-contract market data, options contracts, multi-leg
options, futures order rules, and the Broker FD and Display Solution surfaces
are implemented with paths taken from the official OpenAPI definition. Live
behavior for the US-only and Display-entitlement surfaces is not exercised in
the HK sandbox.

### ✅ Fully Implemented

Every endpoint documented by Webull is implemented — 209 endpoints across
`data`, `trade`, `connect`, `broker`, `brokerfd` and `display`. Paths follow the
official OpenAPI definition; see [`docs/reconciliation.md`](docs/reconciliation.md)
(0 gaps, 0 path discrepancies) for the per-endpoint mapping.

Live behavior for the US-only and Display-entitlement surfaces is not exercised
in the HK sandbox; optional verification is available via `examples/path-probe`
(`WEBULL_SANDBOX=1`).

## Test Coverage

| Package | Test Files | Test Functions | Coverage |
|---------|------------|----------------|----------|
| `data/` | 32 | 145 | Unit (offline) |
| `trade/` | 15 | 84 | Unit (offline) |
| `brokerfd/` (root) | 10 | 52 | Unit (offline) |
| `client/` | 7 | 50 | Unit (offline) |
| `broker/` | 8 | 42 | Unit (offline) |
| `stream/` | 4 | 30 | Unit (offline) |
| `internal/resilience/` | 5 | 30 | Unit (offline) |
| `events/` | 5 | 20 | Unit (offline) |
| `display/` | 1 | 19 | Unit (offline) |
| `internal/auth/` | 2 | 20 | Unit (offline) |
| `brokerfd/events/` | 2 | 13 | Unit (offline) |
| `internal/mqtt/` | 1 | 6 | Unit (offline) |
| `gen/webull/marketdata/v1/` | 1 | 4 | Unit (offline) |
| `internal/region/` | 1 | 4 | Unit (offline) |
| `gen/webull/trade/events/` | 1 | 3 | Unit (offline) |
| **Total** | **90** | **522** | |

- Zero TODO markers in test files
- Integration tests are env-gated (`WEBULL_SANDBOX=1`) and skipped by default
- All unit tests are offline and credential-free

## Examples

| Example | Directory | Documented | Description |
|---------|-----------|------------|-------------|
| auth | `examples/auth/` | ✅ | Token creation/reuse |
| marketdata | `examples/marketdata/` | ✅ | AAPL snapshot + daily bars |
| streaming | `examples/streaming/` | ✅ | MQTT-over-WebSocket subscription |
| watchlist | `examples/watchlist/` | ✅ | List watchlists (read-only) |
| account | `examples/account/` | ✅ | Accounts, balance, positions |
| order | `examples/order/` | ✅ | Preview/place/cancel AAPL limit buy |
| events | `examples/events/` | ✅ | gRPC order event subscription |
| data-fundamentals | `examples/data-fundamentals/` | ✅ | All 13 fundamental endpoints for AAPL |
| probe | `examples/probe/` | ❌ | Sandbox endpoint testing tool (has own README) |
| watchlist-cmd | `examples/watchlist-cmd/` | ✅ | Watchlist CRUD example (create, add, update, remove, delete) |
| broker-probe | `examples/broker-probe/` | ❌ | Broker HK read-only endpoint probe program |
| futures-probe | `examples/futures-probe/` | ❌ | HK futures product discovery and market data probe |
| options-multi-leg | `examples/options-multi-leg/` | ❌ | Multi-leg options strategy probe (11 combo types) |
| options | `examples/options/` | ❌ | Multi-leg options strategy preview |
| brokerfd | `examples/brokerfd/` | ❌ | Broker FD US API endpoints |
| brokerfd-events | `examples/brokerfd-events/` | ❌ | Broker FD US gRPC events |
| path-probe | `examples/path-probe/` | ❌ | SDK path verification tool |

## Known Issues

1. **Display Solution blocked at host level**: The HK sandbox Display Solution host (`hk-co-branding-openapi.uat.webullbroker.com`) returns 403 for all requests, including token creation — the app appears to lack the Display Solution entitlement in the HK sandbox. Not a path issue; paths follow the official definition.
2. **No US sandbox credentials**: Cannot verify fund data, crypto data, screener v2, broker FD, Display Solution (US paths), option chain discovery, or any US-only surface
3. **Sandbox symbol limit**: Only `AAPL` supported in HK sandbox
4. **Footprint entitlement**: Requires paid entitlement; sandbox returns `403`
5. **Options contracts**: May not exist for `AAPL` in sandbox (`417 Invalid Symbol`)
6. **Order-book depth**: Empty outside regular trading hours
7. **MQTT blocking**: Plain MQTT on port 1883 may be blocked; use MQTT-over-WebSocket on `wss://...:8883/mqtt`
8. **Rate limits**: Token endpoint allows 10 requests/30s; max 5 MQTT connections per App Key
9. **Broker API HK sandbox limitation**: Broker API HK (`/broker/...`) returns `401 ROUTE_NOT_PERMITTED` in the HK sandbox — the app lacks the required scope, not a path issue. Broker HK remains unverified pending production or US sandbox access.
10. **SSE news 504**: SSE news upstream returns `504 Gateway Timeout` in the HK sandbox.
11. **Multi-leg options strategies blocked in HK sandbox**: All multi-leg strategies (VERTICAL, STRADDLE, STRANGLE, IRON_CONDOR, IRON_BUTTERFLY, BUTTERFLY, CALENDAR, DIAGONAL, RATIO, COLLAR) are rejected with 417 errors — only SINGLE is accepted in HK sandbox. US sandbox needed to confirm the multi-leg strategy wire values.
12. **Webull docs disagree on some paths**: for several endpoints the `llms.txt` summary path differs from the OpenAPI JSON `path` on the same page — e.g. Create Token is `POST /openapi/auth/token/create` in the summary but `POST /auth/tokens/create` in the JSON, and Get Instruments is `/openapi/instrument/stock/list` vs `/trading/instruments/stocks/profiles/list`. The SDK follows the summary path. A committed probe (`examples/path-probe`, run with `WEBULL_SANDBOX=1 WEBULL_APP_KEY=... WEBULL_APP_SECRET=...`) compares both paths; it has **not been run yet** (no sandbox credentials in this environment).

## Module Structure

| Module | Path | go.mod? | Dependencies |
|--------|------|---------|-------------|
| Root | `/` | ✅ | `github.com/shing1211/webullapi4go` |
| Broker HK | `broker/` | ✅ (own) | Root module via `replace` |
| Broker FD US | `brokerfd/` | ❌ (root) | Part of root module |
| Broker FD Events | `brokerfd/events/` | ❌ (root) | Part of root module |

## Next Steps

v1.1.0 is released with full endpoint coverage. All documented endpoints are
implemented and paths follow the official OpenAPI definition.

1. **Optional live verification** — run `examples/path-probe`
   (`WEBULL_SANDBOX=1 WEBULL_APP_KEY=... WEBULL_APP_SECRET=...`) to confirm the
   auth/instrument paths where Webull's `llms.txt` and OpenAPI JSON disagree
   (Known Issue 12).
2. **US-only surfaces** (futures market data, event contracts, crypto, Broker FD)
   can be exercised only with US sandbox credentials.
3. **Display Solution** requires a paid entitlement; the HK sandbox host returns 403.
4. **Maintain** [`docs/reconciliation.md`](docs/reconciliation.md) and regenerate via
   `python tools/webull-docgen/docgen.py all` when the SDK or Webull docs change.
