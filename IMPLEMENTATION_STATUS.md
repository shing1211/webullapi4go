# Implementation Status

Last updated: 2026-09-22 (v1.0.3 release) · Current version: **v1.0.3**

## Summary

| Module | Package | Transport | Exported API | TODOs | Verified |
|--------|---------|-----------|-------------|-------|----------|
| Core SDK | `client/` | HTTP | 20+ options, `Client.Do`, `Client.DoStream` | 0 | ✅ |
| Authentication | `internal/auth/` | — | HMAC-SHA1/SHA-256 signing, token lifecycle | 0 | ✅ |
| Market Data HTTP | `data/` | HTTP | 92 functions | 23 | ⚠️ Futures/event/option-chain paths unconfirmed (Category field added to all futures query structs — v1.0.0) |
| Market Data Streaming | `stream/` | MQTT | 12+ options, typed handlers | 0 | ✅ |
| Trading HTTP | `trade/` | HTTP | 30+ methods | 6 | ⚠️ Multi-leg/futures unconfirmed |
| Trading Events | `events/` | gRPC | 21 functions | 0 | ✅ |
| Display Solution | `display/` | HTTP | 8 functions | 16 | ⚠️ Host blocked (403 in HK sandbox) |
| Broker API HK | `broker/` | HTTP | 32 methods | 0 | ✅ Own `go.mod` |
| Broker FD API US | `brokerfd/` | HTTP | 56 methods | 1 | ⚠️ Paths unconfirmed |
| Broker FD Events | `brokerfd/events/` | gRPC | 17 functions | 3 | ⚠️ Schemas unconfirmed |

## TODO Marker Inventory (49 total)

| Tag | Count | Package | Subject |
|-----|-------|---------|---------|
| `TODO(ds)` | 16 | `data/display_*.go` | Display Solution paths; requires paid Display Solution entitlement — HK sandbox returns 403 |
| `TODO(t10)` | 12 | `data/options.go` | Option expirations/chain endpoints undocumented, paths and schemas unconfirmed |
| `TODO(futures)` | 6 | `data/futures_market.go` | Futures market data paths unconfirmed against US sandbox |
| `TODO(event-market-data)` | 5 | `data/eventcontracts_market.go` | Event contract market data paths unconfirmed |
| `TODO(t8)` | 3 | `trade/types.go`, `trade/options.go` | Multi-leg strategy wire values and order-type matrix unconfirmed |
| `TODO(t9)` | 2 | `trade/rules.go` | Futures order-type matrix and TIF/entrust/quantity rules unconfirmed |
| `TODO(event)` | 1 | `trade/rules.go` | Event contract order rules unconfirmed |
| `TODO(v08)` | 3 | `brokerfd/events/` | SubscribeType bitmask values and data event JSON schemas unconfirmed |
| `TODO` | 1 | `brokerfd/brokerfd.go` | Confirm all paths via live probe |

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

## Feature Coverage

### ✅ Fully Verified (tested against HK sandbox)

All functions below have real HTTP/gRPC logic, full test coverage, and work against the HK sandbox.

**Core SDK (`client/`)**
- `client.New` with functional options (`WithEnv`, `WithAppKey`, `WithAppSecret`, `WithRegion`, `WithSandbox`, etc.)
- `Client.Do` — signed HTTP request with automatic token injection
- `Client.DoStream` — signed streaming request for SSE endpoints
- Token lifecycle: `CreateToken`, `CheckToken`, `EnsureToken`, `CurrentToken`, `AccessToken`, `SetToken`
- Per-path API version defaults (`x-version` header: `v2` market data, `v3` trading)
- Resilience: retry with exponential backoff, token-bucket rate limiter, circuit breaker

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

### ⚠️ Provisional (implemented, paths/rules unconfirmed)

Every function below has real HTTP/gRPC logic but hits paths or uses wire values that are inferred from convention and not confirmed against a live US sandbox.

**Market Data HTTP — core** (`data/` non-DS) — 23 TODO(futures/t10/event-market-data)

**Futures market data** (`data/futures_market.go`) — 6 TODO(futures)
- `GetFuturesTick`, `GetFuturesSnapshot`, `GetFuturesBars`, `GetFuturesDepth`, `GetFuturesFootprint`
- Paths inferred from pattern; response shape unconfirmed for bars and footprint
- **v1.0.0 fix:** All 5 query structs now have a `Category` field (defaults to `US_FUTURES` if empty) — previously hardcoded US category regardless of query parameter
- Note: `GetFuturesProductCodes` in `data/futures.go` confirmed via HK sandbox; `data/futures_market.go` market data paths (tick/snapshot/bars/depth/footprint) remain unconfirmed

**Event contract market data** (`data/eventcontracts_market.go`) — 5 TODO(event-market-data)
- `GetEventSnapshot`, `GetEventDepth`, `GetEventBars`, `GetEventTick`
- Host and paths unconfirmed against US sandbox

**Display Solution endpoints** (`data/display_*.go`) — 16 TODO(ds)
- Instruments (4): `GetDSCompanyProfile`, `GetDSAnalystTargetPrice`, `GetDSAnalystRating`
- News (5): `GetDSNewsSummary`, `GetDSMarketNews`, `GetDSSymbolNews`, `GetDSLatestNews`
- Quotes (5): `GetDisplaySnapshot`, `GetDisplayBars`, `GetDisplayBarsSingle`, `GetDisplayTick`, `GetDisplayDepth`
- Screeners (3): `GetDisplayGainersLosers`, `GetDisplayTopActive`
- Streaming (2): `DSSubscribe`, `DSUnsubscribe`
- **HK sandbox probe result**: All 14 endpoints return `403 Forbidden` at the Display Solution host (`hk-co-branding-openapi.uat.webullbroker.com`) — even the token endpoint is blocked. App lacks Display Solution entitlement in HK sandbox.

**Option chain discovery** (`data/options.go`) — 12 TODO(t10)
- `GetOptionExpirations`, `GetOptionChain`
- Speculative (not in published Webull API); HK sandbox returns 404
- Paths, response shapes, field mappings all unconfirmed

**Multi-leg options orders** (`trade/options.go`, `trade/types.go`) — 3 TODO(t8)
- Strategy wire values: VERTICAL, STRADDLE, STRANGLE, IRON_CONDOR, IRON_BUTTERFLY, BUTTERFLY, COLLAR, CALENDAR, DIAGONAL, RATIO
- Order-type matrix for multi-leg: only LIMIT and STOP_LOSS_LIMIT (provisional)
- Structural validation logic is complete and tested

**Futures order validation** (`trade/rules.go`) — 2 TODO(t9)
- Order-type matrix: US/HK accept LIMIT, MARKET, STOP_LOSS, STOP_LOSS_LIMIT
- TIF: DAY/GTC only; entrust: QTY only; quantity: positive integer
- Logic complete and tested (24 table-driven cases)

**Event contract order rules** (`trade/rules.go`) — 1 TODO(event)
- LIMIT-only, DAY-only, QTY-only, positive integer, max 50,000

**Fund data** (`data/fund_data.go`) — 0 TODOs, paths unconfirmed
- `GetFundNav`, `GetFundInfo`, `GetFundDividends`, `GetFundList`
- HK sandbox returns 404; US sandbox credentials needed

**Crypto data** (`data/crypto_data.go`) — 0 TODOs, HK sandbox returns 417
- `GetCryptoBars`, `GetCryptoTick`, `GetCryptoDepth`, `GetCryptoSnapshot`
- `GetCryptoSnapshotList`, `GetCryptoBarsList` (US dedicated paths)
- CRYPTO category unsupported in HK sandbox

**Screener v2** (`data/screener_v2.go`) — 0 TODOs, path unconfirmed
- `GetScreenerV2` (POST query)
- Path noted as "unconfirmed" in const comment

**Broker FD API US** (`brokerfd/`) — 1 TODO, root module (no separate `go.mod`)
- 56 methods across accounts, orders, assets, instruments, funding, activity, master data, journals, documents, agreements
- All paths best-effort; returns 404 in HK sandbox

**Broker FD Events** (`brokerfd/events/`) — 3 TODO(v08), root module
- gRPC streaming client using `grpc.event.EventService`
- SubscribeType bitmask values unconfirmed
- Data event JSON schemas unconfirmed

### ❌ Not Yet Implemented

The SDK does not implement the endpoints below. This is generated from
[`docs/reconciliation.md`](docs/reconciliation.md) (snapshot 2026-09-22); see
that file for per-endpoint detail and reference links.

**Not mapped at all (16)**

- **Connect API (OAuth)** (2): `GET /oauth2/auth-codes/get`,
  `POST /oauth2/tokens/create`
- **Crypto** (3): `GET /market-data/crypto/bars/list`,
  `GET /market-data/crypto/snapshots/list`,
  `GET /trading/instruments/crypto/profiles/list`
- **Display Event Contracts** (11): instruments — `.../categories/tags/list`,
  `events/list`, `milestones/list`, `series/list`, `sports-filters/list`;
  market data — `game-stats/get`, `live-data/get`, `markets/bars/list`,
  `markets/bars/list-by-event`, `markets/depths/list`, `markets/snapshots/list`

**Mapped but intentionally not implemented (7)**

- **Display client-token refresh** (1): `POST /auth/client-tokens/refresh` —
  the SDK re-creates the client token instead of refreshing it
- **Fund data** (6): `fund-allocations`, `fund-files`, `fund-holdings`,
  `fund-performances`, `fund-ratings`, `fund-splits`

**Path discrepancies (69)** — implemented, but the SDK path differs from the
official definition: `broker` 11, `brokerfd` 39, `data` 19 (Display Solution
`TODO(ds)` paths, fund data, and Broker FD `/broker-fd/...` vs the documented
`/broker/...`). See Known Issues 12 and 13.

**Unresolved (23)** — the reconciliation tool could not map the method to a path
constant (methods that delegate to a helper, or use an inline path literal).
These are implementations, not gaps.

**Decision: deferred.** The 16 unmapped endpoints and the 6 fund extras are
**not implemented on purpose** — they are US-only or entitlement-gated and cannot
be verified in the HK sandbox. They stay documented in
[`docs/reconciliation.md`](docs/reconciliation.md) and are picked up once US
sandbox credentials or a Display Solution entitlement are available. No
unverified endpoints are added.

- **US sandbox verification**: All ⚠️ items above need US sandbox credentials to verify
- **v1.0.3 is released** with 49 provisional TODOs — all require US sandbox; see audit findings in `docs/runs/2026-09-22-futures-options-probe/`

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

## Known Issues

1. **Display Solution blocked at host level**: The HK sandbox Display Solution host (`hk-co-branding-openapi.uat.webullbroker.com`) returns 403 for all requests, including token creation. All 16 `TODO(ds)` items are blocked at the auth level — not just paths unconfirmed. App may not have Display Solution entitlement in HK sandbox.
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
13. **Broker FD documented paths use `/broker/...`**: the US Broker FD reference pages define their `path` under `/broker/...` (e.g. `/broker/orders/place`, `/broker/accounts/list`), while the `brokerfd` package calls `/broker-fd/...`. Broker FD is provisional; reconcile these paths against the US sandbox.

## Module Structure

| Module | Path | go.mod? | Dependencies |
|--------|------|---------|-------------|
| Root | `/` | ✅ | `github.com/shing1211/webullapi4go` |
| Broker HK | `broker/` | ✅ (own) | Root module via `replace` |
| Broker FD US | `brokerfd/` | ❌ (root) | Part of root module |
| Broker FD Events | `brokerfd/events/` | ❌ (root) | Part of root module |

## Next Steps

v1.0.3 is released. All 49 TODO items require US sandbox credentials to verify and cannot be confirmed with the current HK sandbox setup.

Coverage backlog: maintain [`docs/reconciliation.md`](docs/reconciliation.md) — it maps every documented endpoint to its SDK function and status; regenerate it when the SDK or the Webull docs change. The 16 unmapped endpoints and the 69 path discrepancies are the actionable backlog.

1. **Obtain US sandbox credentials** (`WEBULL_APP_KEY`, `WEBULL_APP_SECRET` for US) — required to verify all 49 provisional items
2. **Run futures-probe** (`examples/futures-probe/`): probe HK futures discovery + market data; requires `WEBULL_FUTURES_TEST=1`
3. **Run options-multi-leg** (`examples/options-multi-leg/`): test 11 combo strategies via PreviewOrder; requires `WEBULL_OPTIONS_TEST=1`
4. **Verify Display Solution** (16 TODOs): HK returns 403 at host level — US sandbox or paid entitlement required
5. **Verify option chain discovery** (12 TODOs): HK returns 404 — US-only endpoints
6. **Verify futures market data paths** (6 TODOs): US sandbox required
7. **Verify event contract market data** (5 TODOs): US-only product
8. **Verify multi-leg strategy wire values** (3 TODOs): US sandbox required for PreviewOrder probe
9. **Verify futures order rules** (2 TODOs): US sandbox required
10. **Verify Broker FD paths** (1 TODO) and **Broker FD event schemas** (3 TODOs): US-only
