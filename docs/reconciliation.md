# SDK ↔ Webull API Reconciliation

> ⚠️ **Generated file — do not edit.** Regenerate with `python tools/webull-docgen/docgen.py <target>` (`reference`, `master`, `reconciliation` or `all`).

> Reconciles every implemented `webullapi4go` function against the official Webull OpenAPI. **Official (OpenAPI JSON)** is the canonical path embedded in the docs; **Official (llms.txt summary)** is the path in Webull's machine-readable index (they disagree for some endpoints). **SDK path** is what the code actually calls.

| | |
|---|---|
| **Snapshot** | 2026-09-22 |
| **Sources** | [HK llms.txt](https://developer.webull.hk/apis/llms.txt), [US llms.txt](https://developer.webull.com/apis/llms.txt) |
| **Implemented endpoints** | 209 |
| **Documented-only endpoints (gaps)** | 0 |
| ✅ Path matches OpenAPI JSON | 180 |
| 🟡 Matches docs summary only | 4 |
| ⚠️ Path differs from both | 0 |
| ❓ Unresolved | 25 |
| ℹ️ Intentionally not implemented | 0 |

## Implemented endpoints

## Authentication

### Create Token

| | |
|---|---|
| **SDK** | `client.CreateToken` |
| **Official (OpenAPI JSON)** | `POST /auth/tokens/create` |
| **Official (llms.txt summary)** | `/openapi/auth/token/create` |
| **SDK path** | `/openapi/auth/token/create` (createTokenPath) |
| **Status** | 🟡 SDK matches docs summary, not OpenAPI JSON |
| **Note** | Returns a token in `PENDING`; completes 2FA in the Webull App. Sandbox tokens are `NORMAL` immediately. |

Reference: [create-token.md](https://developer.webull.hk/apis/docs/reference/create-token.md)

### Check Token

| | |
|---|---|
| **SDK** | `client.CheckToken` |
| **Official (OpenAPI JSON)** | `POST /auth/tokens/check` |
| **Official (llms.txt summary)** | `/openapi/auth/token/check` |
| **SDK path** | `/openapi/auth/token/check` (checkTokenPath) |
| **Status** | 🟡 SDK matches docs summary, not OpenAPI JSON |
| **Note** | Polled by `client.EnsureToken` until the token is `NORMAL`. |

Reference: [check-token.md](https://developer.webull.hk/apis/docs/reference/check-token.md)

### Create Client Token (Display)

| | |
|---|---|
| **SDK** | `display.Service.EnsureToken` |
| **Official (OpenAPI JSON)** | `POST /auth/client-tokens/create` |
| **Official (llms.txt summary)** | `/openapi/auth/client/token/create` |
| **Status** | ❓ SDK path unresolved |
| **Note** | Internal to `display.Service`; `client_user_id` is fixed to `openapi_client`. |

Reference: [create-client-token.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/create-client-token.md)

### Refresh Client Token (Display)

| | |
|---|---|
| **SDK** | `display.Service.RefreshClientToken` |
| **Official (OpenAPI JSON)** | `POST /auth/client-tokens/refresh` |
| **Official (llms.txt summary)** | `/openapi/auth/client/token/refresh` |
| **SDK path** | `/auth/client-tokens/refresh` (displayTokenRefreshPath) |
| **Status** | ✅ match |
| **Note** | The SDK re-creates the client token instead of refreshing it. |

Reference: [refresh-client-token.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/refresh-client-token.md)

## Market Data — Stock

### Stock Snapshot

| | |
|---|---|
| **SDK** | `data.GetSnapshot` |
| **Official (OpenAPI JSON)** | `GET /market-data/stocks/snapshots/list` |
| **Official (llms.txt summary)** | `/openapi/market-data/stock/snapshot` |
| **SDK path** | `/market-data/stocks/snapshots/list` (pathStockSnapshots) |
| **Status** | ✅ match |

Reference: [snapshot.md](https://developer.webull.hk/apis/docs/reference/snapshot.md)

### Stock Tick

| | |
|---|---|
| **SDK** | `data.GetTick` |
| **Official (OpenAPI JSON)** | `GET /market-data/stocks/ticks/list` |
| **Official (llms.txt summary)** | `/openapi/market-data/stock/tick` |
| **SDK path** | `/market-data/stocks/ticks/list` (pathStockTicks) |
| **Status** | ✅ match |

Reference: [tick.md](https://developer.webull.hk/apis/docs/reference/tick.md)

### Stock Quotes (Depth)

| | |
|---|---|
| **SDK** | `data.GetQuotes` |
| **Official (OpenAPI JSON)** | `GET /market-data/stocks/depths/list` |
| **Official (llms.txt summary)** | `/openapi/market-data/stock/quotes` |
| **SDK path** | `/market-data/stocks/depths/list` (pathStockDepths) |
| **Status** | ✅ match |

Reference: [quotes.md](https://developer.webull.hk/apis/docs/reference/quotes.md)

### Stock Historical Bars (Batch)

| | |
|---|---|
| **SDK** | `data.GetBatchBars` |
| **Official (OpenAPI JSON)** | `POST /market-data/stocks/bars/list` |
| **Official (llms.txt summary)** | `/openapi/market-data/stock/batch-bars` |
| **SDK path** | `/market-data/stocks/bars/list` (pathStockBarsList) |
| **Status** | ✅ match |
| **Note** | `data.GetBars` POSTs a one-symbol batch against this same endpoint. |

Reference: [historical-bars.md](https://developer.webull.hk/apis/docs/reference/historical-bars.md)

### Stock Footprint

| | |
|---|---|
| **SDK** | `data.GetFootprint` |
| **Official (OpenAPI JSON)** | `GET /market-data/stocks/footprints/list` |
| **Official (llms.txt summary)** | `/openapi/market-data/stock/footprint` |
| **SDK path** | `/market-data/stocks/footprints/list` (pathStockFootprints) |
| **Status** | ✅ match |
| **Note** | Requires a paid entitlement. |

Reference: [footprint.md](https://developer.webull.hk/apis/docs/reference/footprint.md)

### NOII Bars

| | |
|---|---|
| **SDK** | `data.GetNOIIBars` |
| **Official (OpenAPI JSON)** | `GET /market-data/stocks/noii-bars/list` |
| **SDK path** | `/market-data/stocks/noii-bars/list` (pathNOIIBars) |
| **Status** | ✅ match |

Reference: [get-noii-bars.md](https://developer.webull.hk/apis/docs/reference/get-noii-bars.md)

### NOII Snapshot

| | |
|---|---|
| **SDK** | `data.GetNOIISnapshot` |
| **Official (OpenAPI JSON)** | `GET /market-data/stocks/noii-snapshots/list` |
| **SDK path** | `/market-data/stocks/noii-snapshots/list` (pathNOIISnapshot) |
| **Status** | ✅ match |

Reference: [get-noii-snapshot.md](https://developer.webull.hk/apis/docs/reference/get-noii-snapshot.md)

## Market Data — Option

### Option Tick

| | |
|---|---|
| **SDK** | `data.GetOptionTick` |
| **Official (OpenAPI JSON)** | `GET /market-data/options/ticks/list` |
| **SDK path** | `/market-data/options/ticks/list` (pathOptionTicks) |
| **Status** | ✅ match |

Reference: [option-tick.md](https://developer.webull.hk/apis/docs/reference/option-tick.md)

### Option Snapshot

| | |
|---|---|
| **SDK** | `data.GetOptionSnapshot` |
| **Official (OpenAPI JSON)** | `GET /market-data/options/snapshots/list` |
| **SDK path** | `/market-data/options/snapshots/list` (pathOptionSnapshots) |
| **Status** | ✅ match |
| **Note** | At most 20 symbols per query. |

Reference: [option-snapshot.md](https://developer.webull.hk/apis/docs/reference/option-snapshot.md)

### Option Historical Bars

| | |
|---|---|
| **SDK** | `data.GetOptionBars` |
| **Official (OpenAPI JSON)** | `GET /market-data/options/bars/list` |
| **SDK path** | `/market-data/options/bars/list` (pathOptionBars) |
| **Status** | ✅ match |
| **Note** | At most 20 symbols per query. |

Reference: [option-historical-bars.md](https://developer.webull.hk/apis/docs/reference/option-historical-bars.md)

### Option Contract List (Chain)

| | |
|---|---|
| **SDK** | `data.GetOptionContracts` |
| **Official (OpenAPI JSON)** | `GET /trading/instruments/options/contracts/list` |
| **SDK path** | `/trading/instruments/options/contracts/list` (pathOptionContracts) |
| **Status** | ✅ match |
| **Note** | Trading API surface. |

Reference: [option-contract-list.md](https://developer.webull.hk/apis/docs/reference/option-contract-list.md)

## Market Data — Crypto

### Crypto Snapshot

| | |
|---|---|
| **SDK** | `data.GetCryptoSnapshot` |
| **Official (OpenAPI JSON)** | `GET /market-data/crypto/snapshots/list` |
| **SDK path** | `/market-data/crypto/snapshots/list` (pathCryptoSnapshots) |
| **Status** | ✅ match |

Reference: [crypto-snapshot.md](https://developer.webull.com/apis/docs/reference/crypto-snapshot.md)

### Crypto Bars

| | |
|---|---|
| **SDK** | `data.GetCryptoBars` |
| **Official (OpenAPI JSON)** | `GET /market-data/crypto/bars/list` |
| **SDK path** | `/market-data/crypto/bars/list` (pathCryptoBars) |
| **Status** | ✅ match |

Reference: [crypto-bars.md](https://developer.webull.com/apis/docs/reference/crypto-bars.md)

### Crypto Instruments

| | |
|---|---|
| **SDK** | `data.GetCryptoInstruments` |
| **Official (OpenAPI JSON)** | `GET /trading/instruments/crypto/profiles/list` |
| **SDK path** | `/trading/instruments/crypto/profiles/list` (pathCryptoInstruments) |
| **Status** | ✅ match |
| **Note** | Trading API surface. |

Reference: [crypto-instrument-list.md](https://developer.webull.com/apis/docs/reference/crypto-instrument-list.md)

## Market Data — Futures

### Futures Tick

| | |
|---|---|
| **SDK** | `data.GetFuturesTick` |
| **Official (OpenAPI JSON)** | `GET /market-data/futures/ticks/list` |
| **SDK path** | `/market-data/futures/ticks/list` (pathFuturesTick) |
| **Status** | ✅ match |

Reference: [futures-tick.md](https://developer.webull.hk/apis/docs/reference/futures-tick.md)

### Futures Snapshot

| | |
|---|---|
| **SDK** | `data.GetFuturesSnapshot` |
| **Official (OpenAPI JSON)** | `GET /market-data/futures/snapshots/list` |
| **SDK path** | `/market-data/futures/snapshots/list` (pathFuturesSnapshot) |
| **Status** | ✅ match |

Reference: [futures-snapshot.md](https://developer.webull.hk/apis/docs/reference/futures-snapshot.md)

### Futures Footprint

| | |
|---|---|
| **SDK** | `data.GetFuturesFootprint` |
| **Official (OpenAPI JSON)** | `GET /market-data/futures/footprints/list` |
| **SDK path** | `/market-data/futures/footprints/list` (pathFuturesFootprint) |
| **Status** | ✅ match |

Reference: [futures-footprint.md](https://developer.webull.hk/apis/docs/reference/futures-footprint.md)

### Futures Quotes (Depth)

| | |
|---|---|
| **SDK** | `data.GetFuturesDepth` |
| **Official (OpenAPI JSON)** | `GET /market-data/futures/depths/list` |
| **SDK path** | `/market-data/futures/depths/list` (pathFuturesDepth) |
| **Status** | ✅ match |

Reference: [futures-depth-of-book.md](https://developer.webull.hk/apis/docs/reference/futures-depth-of-book.md)

### Futures Historical Bars

| | |
|---|---|
| **SDK** | `data.GetFuturesBars` |
| **Official (OpenAPI JSON)** | `GET /market-data/futures/bars/list` |
| **SDK path** | `/market-data/futures/bars/list` (pathFuturesBars) |
| **Status** | ✅ match |

Reference: [futures-historical-bars.md](https://developer.webull.hk/apis/docs/reference/futures-historical-bars.md)

### Futures Instrument List

| | |
|---|---|
| **SDK** | `data.GetFuturesInstruments` |
| **Official (OpenAPI JSON)** | `GET /trading/instruments/futures/contracts/list` |
| **SDK path** | `/trading/instruments/futures/contracts/list` (pathFuturesInstruments) |
| **Status** | ✅ match |

Reference: [futures-instrument-list.md](https://developer.webull.hk/apis/docs/reference/futures-instrument-list.md)

### Futures Product Codes

| | |
|---|---|
| **SDK** | `data.GetFuturesProductCodes` |
| **Official (OpenAPI JSON)** | `GET /trading/instruments/futures/product-codes/list` |
| **SDK path** | `/trading/instruments/futures/product-codes/list` (pathFuturesProductCodes) |
| **Status** | ✅ match |

Reference: [futures-products.md](https://developer.webull.hk/apis/docs/reference/futures-products.md)

### Futures Product Classes

| | |
|---|---|
| **SDK** | `data.GetFuturesProductClasses` |
| **Official (OpenAPI JSON)** | `GET /trading/instruments/futures/product-classes/list` |
| **SDK path** | `/trading/instruments/futures/product-classes/list` (pathFuturesProductClasses) |
| **Status** | ✅ match |

Reference: [futures-products-class.md](https://developer.webull.hk/apis/docs/reference/futures-products-class.md)

## Market Data — News

### News Summary

| | |
|---|---|
| **SDK** | `data.GetNewsSummary` |
| **Official (OpenAPI JSON)** | `POST /market-data/news/summaries/get` |
| **SDK path** | `/market-data/news/summaries/get` (pathNewsSummaries) |
| **Status** | ✅ match |
| **Note** | SSE stream via `client.DoStream`; HK sandbox upstream returns 504. |

Reference: [news-summary.md](https://developer.webull.hk/apis/docs/reference/news-summary.md)

## Market Data — Screener

### Top Gainers/Losers

| | |
|---|---|
| **SDK** | `data.GetTopGainersLosers` |
| **Official (OpenAPI JSON)** | `GET /market-data/screeners/gainers-losers/list` |
| **SDK path** | `/market-data/screeners/gainers-losers/list` (pathGainersLosers) |
| **Status** | ✅ match |

Reference: [get-gainers-losers.md](https://developer.webull.hk/apis/docs/reference/get-gainers-losers.md)

### Top Actives

| | |
|---|---|
| **SDK** | `data.GetMostActive` |
| **Official (OpenAPI JSON)** | `GET /market-data/screeners/top-actives/list` |
| **SDK path** | `/market-data/screeners/top-actives/list` (pathTopActives) |
| **Status** | ✅ match |

Reference: [get-top-active.md](https://developer.webull.hk/apis/docs/reference/get-top-active.md)

### Market Sectors

| | |
|---|---|
| **SDK** | `data.GetMarketSectors` |
| **Official (OpenAPI JSON)** | `GET /market-data/screeners/market-sectors/list` |
| **SDK path** | `/market-data/screeners/market-sectors/list` (pathMarketSectors) |
| **Status** | ✅ match |
| **Note** | Display Solution; HK sandbox host blocked (403). |

Reference: [get-market-sectors.md](https://developer.webull.hk/apis/docs/reference/get-market-sectors.md)

### Market Sector Detail

| | |
|---|---|
| **SDK** | `data.GetMarketSectorDetail` |
| **Official (OpenAPI JSON)** | `GET /market-data/screeners/market-sectors/get` |
| **SDK path** | `/market-data/screeners/market-sectors/get` (pathMarketSectorDetail) |
| **Status** | ✅ match |
| **Note** | Display Solution. |

Reference: [get-market-sectors-detail.md](https://developer.webull.hk/apis/docs/reference/get-market-sectors-detail.md)

### High Dividend Rank

| | |
|---|---|
| **SDK** | `data.GetHighDividendRank` |
| **Official (OpenAPI JSON)** | `GET /market-data/screeners/high-dividend-ranks/list` |
| **SDK path** | `/market-data/screeners/high-dividend-ranks/list` (pathHighDividend) |
| **Status** | ✅ match |
| **Note** | Display Solution. |

Reference: [get-high-dividend.md](https://developer.webull.hk/apis/docs/reference/get-high-dividend.md)

### 52-Week High/Low

| | |
|---|---|
| **SDK** | `data.GetWeek52HighLow` |
| **Official (OpenAPI JSON)** | `GET /market-data/screeners/week52-high-low/list` |
| **SDK path** | `/market-data/screeners/week52-high-low/list` (pathWeek52HighLow) |
| **Status** | ✅ match |
| **Note** | Display Solution. |

Reference: [get-week-52-high-low.md](https://developer.webull.hk/apis/docs/reference/get-week-52-high-low.md)

## Market Data — Watchlist

### Get Watchlists

| | |
|---|---|
| **SDK** | `data.GetWatchlists` |
| **Official (OpenAPI JSON)** | `GET /market-data/watchlists/list` |
| **SDK path** | `/market-data/watchlists/list` (pathWatchlists) |
| **Status** | ✅ match |

Reference: [get-watchlist.md](https://developer.webull.hk/apis/docs/reference/get-watchlist.md)

### Create Watchlist

| | |
|---|---|
| **SDK** | `data.CreateWatchlist` |
| **Official (OpenAPI JSON)** | `POST /market-data/watchlists/create` |
| **SDK path** | `/market-data/watchlists/create` (pathWatchlistCreate) |
| **Status** | ✅ match |

Reference: [create-watchlist.md](https://developer.webull.hk/apis/docs/reference/create-watchlist.md)

### Update Watchlist

| | |
|---|---|
| **SDK** | `data.UpdateWatchlist` |
| **Official (OpenAPI JSON)** | `POST /market-data/watchlists/update` |
| **SDK path** | `/market-data/watchlists/update` (pathWatchlistUpdate) |
| **Status** | ✅ match |

Reference: [update-watchlist.md](https://developer.webull.hk/apis/docs/reference/update-watchlist.md)

### Delete Watchlist

| | |
|---|---|
| **SDK** | `data.DeleteWatchlist` |
| **Official (OpenAPI JSON)** | `POST /market-data/watchlists/delete` |
| **SDK path** | `/market-data/watchlists/delete` (pathWatchlistDelete) |
| **Status** | ✅ match |

Reference: [delete-watchlist.md](https://developer.webull.hk/apis/docs/reference/delete-watchlist.md)

### Get Watchlist Instruments

| | |
|---|---|
| **SDK** | `data.GetWatchlistInstruments` |
| **Official (OpenAPI JSON)** | `GET /market-data/watchlists/instruments/list` |
| **SDK path** | `/market-data/watchlists/instruments/list` (pathWatchlistInstruments) |
| **Status** | ✅ match |

Reference: [get-watchlist-instruments.md](https://developer.webull.hk/apis/docs/reference/get-watchlist-instruments.md)

### Add Watchlist Instruments

| | |
|---|---|
| **SDK** | `data.AddWatchlistInstruments` |
| **Official (OpenAPI JSON)** | `POST /market-data/watchlists/instruments/add` |
| **SDK path** | `/market-data/watchlists/instruments/add` (pathWatchlistInstrumentsAdd) |
| **Status** | ✅ match |

Reference: [add-watchlist-instruments.md](https://developer.webull.hk/apis/docs/reference/add-watchlist-instruments.md)

### Remove Watchlist Instruments

| | |
|---|---|
| **SDK** | `data.RemoveWatchlistInstruments` |
| **Official (OpenAPI JSON)** | `POST /market-data/watchlists/instruments/remove` |
| **SDK path** | `/market-data/watchlists/instruments/remove` (pathWatchlistInstrumentsRemove) |
| **Status** | ✅ match |

Reference: [remove-watchlist-instruments.md](https://developer.webull.hk/apis/docs/reference/remove-watchlist-instruments.md)

### Update Watchlist Instruments

| | |
|---|---|
| **SDK** | `data.UpdateWatchlistInstruments` |
| **Official (OpenAPI JSON)** | `POST /market-data/watchlists/instruments/update` |
| **SDK path** | `/market-data/watchlists/instruments/update` (pathWatchlistInstrumentsUpdate) |
| **Status** | ✅ match |

Reference: [update-watchlist-instruments.md](https://developer.webull.hk/apis/docs/reference/update-watchlist-instruments.md)

## Fundamentals and Fund Data

### Company Profile

| | |
|---|---|
| **SDK** | `data.GetCompanyProfile` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/company-profiles/get` |
| **SDK path** | `/market-data/fundamentals/company-profiles/get` (pathCompanyProfile) |
| **Status** | ✅ match |
| **Note** | US only in practice. |

Reference: [get-company-profile.md](https://developer.webull.hk/apis/docs/reference/get-company-profile.md)

### Analyst Target Price

| | |
|---|---|
| **SDK** | `data.GetAnalystTargetPrice` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/analysis/target-prices/get` |
| **SDK path** | `/market-data/fundamentals/analysis/target-prices/get` (pathAnalystTarget) |
| **Status** | ✅ match |

Reference: [get-analyst-target-price.md](https://developer.webull.hk/apis/docs/reference/get-analyst-target-price.md)

### Analyst Rating

| | |
|---|---|
| **SDK** | `data.GetAnalystRating` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/analysis/ratings/get` |
| **SDK path** | `/market-data/fundamentals/analysis/ratings/get` (pathAnalystRating) |
| **Status** | ✅ match |

Reference: [get-analyst-rating.md](https://developer.webull.hk/apis/docs/reference/get-analyst-rating.md)

### Capital Flow

| | |
|---|---|
| **SDK** | `data.GetCapitalFlow` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/capital-flows/get` |
| **SDK path** | `/market-data/fundamentals/capital-flows/get` (pathCapitalFlow) |
| **Status** | ✅ match |

Reference: [capital-flow.md](https://developer.webull.hk/apis/docs/reference/capital-flow.md)

### Industry Comparison

| | |
|---|---|
| **SDK** | `data.GetIndustryComparison` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/industry-comparisons/get` |
| **SDK path** | `/market-data/fundamentals/industry-comparisons/get` (pathIndustryComp) |
| **Status** | ✅ match |

Reference: [industry-comparison.md](https://developer.webull.hk/apis/docs/reference/industry-comparison.md)

### Earnings Calendar

| | |
|---|---|
| **SDK** | `data.GetEarningsCalendar` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/earnings-calendars/list` |
| **SDK path** | `/market-data/fundamentals/earnings-calendars/list` (pathEarningsCal) |
| **Status** | ✅ match |

Reference: [earnings-calendar.md](https://developer.webull.hk/apis/docs/reference/earnings-calendar.md)

### Dividend Calendar

| | |
|---|---|
| **SDK** | `data.GetDividendCalendar` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/dividend-calendars/list` |
| **SDK path** | `/market-data/fundamentals/dividend-calendars/list` (pathDividendCal) |
| **Status** | ✅ match |

Reference: [dividend-calendar.md](https://developer.webull.hk/apis/docs/reference/dividend-calendar.md)

### Filings

| | |
|---|---|
| **SDK** | `data.GetFilings` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/filings/list` |
| **SDK path** | `/market-data/fundamentals/filings/list` (pathFilings) |
| **Status** | ✅ match |
| **Note** | US/SEC only. |

Reference: [filings.md](https://developer.webull.hk/apis/docs/reference/filings.md)

### Income Statement

| | |
|---|---|
| **SDK** | `data.GetIncomeStatement` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/income-statements/get` |
| **SDK path** | `/market-data/fundamentals/income-statements/get` (pathIncomeStmt) |
| **Status** | ✅ match |
| **Note** | Returns `map[string]any` per period. |

Reference: [financial-income.md](https://developer.webull.hk/apis/docs/reference/financial-income.md)

### Balance Sheet

| | |
|---|---|
| **SDK** | `data.GetBalanceSheet` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/balance-sheets/get` |
| **SDK path** | `/market-data/fundamentals/balance-sheets/get` (pathBalanceSheet) |
| **Status** | ✅ match |
| **Note** | Returns `map[string]any` per period. |

Reference: [financial-balancesheet.md](https://developer.webull.hk/apis/docs/reference/financial-balancesheet.md)

### Cash Flow

| | |
|---|---|
| **SDK** | `data.GetCashFlow` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/cash-flows/get` |
| **SDK path** | `/market-data/fundamentals/cash-flows/get` (pathCashFlow) |
| **Status** | ✅ match |
| **Note** | Returns `map[string]any` per period. |

Reference: [financial-cashflow.md](https://developer.webull.hk/apis/docs/reference/financial-cashflow.md)

### Financial Indicators

| | |
|---|---|
| **SDK** | `data.GetFinancialIndicators` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/indicators/get` |
| **SDK path** | `/market-data/fundamentals/indicators/get` (pathIndicators) |
| **Status** | ✅ match |

Reference: [financial-indicators.md](https://developer.webull.hk/apis/docs/reference/financial-indicators.md)

### Financial Alert

| | |
|---|---|
| **SDK** | `data.GetFinancialAlert` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/financial-alerts/get` |
| **SDK path** | `/market-data/fundamentals/financial-alerts/get` (pathFinancialAlert) |
| **Status** | ✅ match |

Reference: [financial-alert.md](https://developer.webull.hk/apis/docs/reference/financial-alert.md)

### Forecast EPS

| | |
|---|---|
| **SDK** | `data.GetForecastEPS` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/forecast-eps/get` |
| **SDK path** | `/market-data/fundamentals/forecast-eps/get` (?) |
| **Status** | ✅ match |
| **Note** | Most recent 5 quarters. |

Reference: [forecast-eps.md](https://developer.webull.hk/apis/docs/reference/forecast-eps.md)

### Fund Brief

| | |
|---|---|
| **SDK** | `data.GetFundInfo` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/fund-brief/get` |
| **SDK path** | `/market-data/fundamentals/fund-brief/get` (pathFundInfo) |
| **Status** | ✅ match |

Reference: [fund-brief.md](https://developer.webull.hk/apis/docs/reference/fund-brief.md)

### Fund Performance

| | |
|---|---|
| **SDK** | `data.GetFundPerformance` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/fund-performances/get` |
| **SDK path** | `/market-data/fundamentals/fund-performances/get` (pathFundPerformance) |
| **Status** | ✅ match |

Reference: [fund-performance.md](https://developer.webull.hk/apis/docs/reference/fund-performance.md)

### Fund Net Value

| | |
|---|---|
| **SDK** | `data.GetFundNav` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/fund-net-values/get` |
| **SDK path** | `/market-data/fundamentals/fund-net-values/get` (pathFundNav) |
| **Status** | ✅ match |

Reference: [fund-net-value.md](https://developer.webull.hk/apis/docs/reference/fund-net-value.md)

### Fund Holdings

| | |
|---|---|
| **SDK** | `data.GetFundHoldings` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/fund-holdings/get` |
| **SDK path** | `/market-data/fundamentals/fund-holdings/get` (pathFundHoldings) |
| **Status** | ✅ match |

Reference: [fund-holdings.md](https://developer.webull.hk/apis/docs/reference/fund-holdings.md)

### Fund Dividends

| | |
|---|---|
| **SDK** | `data.GetFundDividends` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/fund-dividends/get` |
| **SDK path** | `/market-data/fundamentals/fund-dividends/get` (pathFundDividends) |
| **Status** | ✅ match |

Reference: [fund-dividends.md](https://developer.webull.hk/apis/docs/reference/fund-dividends.md)

### Fund Rating

| | |
|---|---|
| **SDK** | `data.GetFundRating` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/fund-ratings/get` |
| **SDK path** | `/market-data/fundamentals/fund-ratings/get` (pathFundRating) |
| **Status** | ✅ match |

Reference: [fund-rating.md](https://developer.webull.hk/apis/docs/reference/fund-rating.md)

### Fund Splits

| | |
|---|---|
| **SDK** | `data.GetFundSplits` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/fund-splits/get` |
| **SDK path** | `/market-data/fundamentals/fund-splits/get` (pathFundSplits) |
| **Status** | ✅ match |

Reference: [fund-splits.md](https://developer.webull.hk/apis/docs/reference/fund-splits.md)

### Fund Files

| | |
|---|---|
| **SDK** | `data.GetFundFiles` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/fund-files/get` |
| **SDK path** | `/market-data/fundamentals/fund-files/get` (pathFundFiles) |
| **Status** | ✅ match |

Reference: [fund-files.md](https://developer.webull.hk/apis/docs/reference/fund-files.md)

### Fund Allocation

| | |
|---|---|
| **SDK** | `data.GetFundAllocation` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/fund-allocations/get` |
| **SDK path** | `/market-data/fundamentals/fund-allocations/get` (pathFundAllocation) |
| **Status** | ✅ match |

Reference: [fund-allocation.md](https://developer.webull.hk/apis/docs/reference/fund-allocation.md)

## Event Contracts

### Event Contract Categories

| | |
|---|---|
| **SDK** | `data.GetEventContractCategories` |
| **Official (OpenAPI JSON)** | `GET /trading/instruments/event-contracts/categories/list` |
| **SDK path** | `/trading/instruments/event-contracts/categories/list` (pathEventContractCategories) |
| **Status** | ✅ match |
| **Note** | Documented on the US site. |

Reference: [event-categories-list.md](https://developer.webull.com/apis/docs/reference/event-categories-list.md)

### Event Contract Series

| | |
|---|---|
| **SDK** | `data.GetEventContractSeries` |
| **Official (OpenAPI JSON)** | `GET /trading/instruments/event-contracts/series/list` |
| **SDK path** | `/trading/instruments/event-contracts/series/list` (pathEventContractSeries) |
| **Status** | ✅ match |

Reference: [event-series-list.md](https://developer.webull.com/apis/docs/reference/event-series-list.md)

### Event Contract Events

| | |
|---|---|
| **SDK** | `data.GetEventContractEvents` |
| **Official (OpenAPI JSON)** | `GET /trading/instruments/event-contracts/events/list` |
| **SDK path** | `/trading/instruments/event-contracts/events/list` (pathEventContractEvents) |
| **Status** | ✅ match |

Reference: [event-events-list.md](https://developer.webull.com/apis/docs/reference/event-events-list.md)

### Event Contract Instruments

| | |
|---|---|
| **SDK** | `data.GetEventContractMarkets` |
| **Official (OpenAPI JSON)** | `GET /trading/instruments/event-contracts/markets/list` |
| **SDK path** | `/trading/instruments/event-contracts/markets/list` (pathEventContractMarkets) |
| **Status** | ✅ match |

Reference: [event-market-list.md](https://developer.webull.com/apis/docs/reference/event-market-list.md)

### Event Snapshot

| | |
|---|---|
| **SDK** | `data.GetEventSnapshot` |
| **Official (OpenAPI JSON)** | `GET /market-data/event-contracts/snapshots/list` |
| **SDK path** | `/market-data/event-contracts/snapshots/list` (pathEventSnapshot) |
| **Status** | ✅ match |

Reference: [event-snapshot.md](https://developer.webull.com/apis/docs/reference/event-snapshot.md)

### Event Depth

| | |
|---|---|
| **SDK** | `data.GetEventDepth` |
| **Official (OpenAPI JSON)** | `GET /market-data/event-contracts/depths/list` |
| **SDK path** | `/market-data/event-contracts/depths/list` (pathEventDepth) |
| **Status** | ✅ match |

Reference: [event-depth.md](https://developer.webull.com/apis/docs/reference/event-depth.md)

### Event Bars

| | |
|---|---|
| **SDK** | `data.GetEventBars` |
| **Official (OpenAPI JSON)** | `GET /market-data/event-contracts/bars/list` |
| **SDK path** | `/market-data/event-contracts/bars/list` (pathEventBars) |
| **Status** | ✅ match |

Reference: [event-bars.md](https://developer.webull.com/apis/docs/reference/event-bars.md)

### Event Tick

| | |
|---|---|
| **SDK** | `data.GetEventTick` |
| **Official (OpenAPI JSON)** | `GET /market-data/event-contracts/ticks/list` |
| **SDK path** | `/market-data/event-contracts/ticks/list` (pathEventTick) |
| **Status** | ✅ match |

Reference: [event-tick.md](https://developer.webull.com/apis/docs/reference/event-tick.md)

### Event Contract Tags

| | |
|---|---|
| **SDK** | `data.GetEventContractTags` |
| **Official (OpenAPI JSON)** | `GET /market-data/instruments/event-contracts/categories/tags/list` |
| **SDK path** | `/market-data/instruments/event-contracts/categories/tags/list` (pathEventTags) |
| **Status** | ✅ match |

Reference: [all-tags-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/all-tags-using-get.md)

### Event Contract Events List

| | |
|---|---|
| **SDK** | `data.GetEventContractEventsList` |
| **Official (OpenAPI JSON)** | `GET /market-data/instruments/event-contracts/events/list` |
| **SDK path** | `/market-data/instruments/event-contracts/events/list` (pathEventEventsList) |
| **Status** | ✅ match |

Reference: [event-list-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-list-using-get.md)

### Event Contract Milestones

| | |
|---|---|
| **SDK** | `data.GetEventContractMilestones` |
| **Official (OpenAPI JSON)** | `GET /market-data/instruments/event-contracts/milestones/list` |
| **SDK path** | `/market-data/instruments/event-contracts/milestones/list` (pathEventMilestones) |
| **Status** | ✅ match |

Reference: [milestones-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/milestones-using-get.md)

### Event Contract Series List

| | |
|---|---|
| **SDK** | `data.GetEventContractSeriesList` |
| **Official (OpenAPI JSON)** | `GET /market-data/instruments/event-contracts/series/list` |
| **SDK path** | `/market-data/instruments/event-contracts/series/list` (pathEventSeriesList) |
| **Status** | ✅ match |

Reference: [series-list-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/series-list-using-get.md)

### Event Contract Sports Filters

| | |
|---|---|
| **SDK** | `data.GetEventContractSportsFilters` |
| **Official (OpenAPI JSON)** | `GET /market-data/instruments/event-contracts/sports-filters/list` |
| **SDK path** | `/market-data/instruments/event-contracts/sports-filters/list` (pathEventSportsFilters) |
| **Status** | ✅ match |

Reference: [sports-filter-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/sports-filter-using-get.md)

### Event Game Stats

| | |
|---|---|
| **SDK** | `data.GetEventGameStats` |
| **Official (OpenAPI JSON)** | `GET /market-data/event-contracts/game-stats/get` |
| **SDK path** | `/market-data/event-contracts/game-stats/get` (pathEventGameStats) |
| **Status** | ✅ match |

Reference: [event-game-stats-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-game-stats-using-get.md)

### Event Live Data

| | |
|---|---|
| **SDK** | `data.GetEventLiveData` |
| **Official (OpenAPI JSON)** | `GET /market-data/event-contracts/live-data/get` |
| **SDK path** | `/market-data/event-contracts/live-data/get` (pathEventLiveData) |
| **Status** | ✅ match |

Reference: [event-live-data-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-live-data-using-get.md)

### Event Market Bars

| | |
|---|---|
| **SDK** | `data.GetEventMarketBars` |
| **Official (OpenAPI JSON)** | `GET /market-data/event-contracts/markets/bars/list` |
| **SDK path** | `/market-data/event-contracts/markets/bars/list` (pathEventMarketBars) |
| **Status** | ✅ match |

Reference: [event-market-bars-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-market-bars-using-get.md)

### Event Market Bars By Event

| | |
|---|---|
| **SDK** | `data.GetEventMarketBarsByEvent` |
| **Official (OpenAPI JSON)** | `GET /market-data/event-contracts/markets/bars/list-by-event` |
| **SDK path** | `/market-data/event-contracts/markets/bars/list-by-event` (pathEventMarketBarsEvt) |
| **Status** | ✅ match |

Reference: [event-market-bars-by-event-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-market-bars-by-event-using-get.md)

### Event Market Depth

| | |
|---|---|
| **SDK** | `data.GetEventMarketDepth` |
| **Official (OpenAPI JSON)** | `GET /market-data/event-contracts/markets/depths/list` |
| **SDK path** | `/market-data/event-contracts/markets/depths/list` (pathEventMarketDepths) |
| **Status** | ✅ match |

Reference: [event-market-depth-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-market-depth-using-get.md)

### Event Market Snapshot

| | |
|---|---|
| **SDK** | `data.GetEventMarketSnapshot` |
| **Official (OpenAPI JSON)** | `GET /market-data/event-contracts/markets/snapshots/list` |
| **SDK path** | `/market-data/event-contracts/markets/snapshots/list` (pathEventMarketSnapshot) |
| **Status** | ✅ match |

Reference: [event-market-snapshot-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-market-snapshot-using-get.md)

## Trading API

### Get Instruments

| | |
|---|---|
| **SDK** | `data.GetStockInstruments` |
| **Official (OpenAPI JSON)** | `GET /trading/instruments/stocks/profiles/list` |
| **Official (llms.txt summary)** | `/openapi/instrument/stock/list` |
| **SDK path** | `/openapi/instrument/stock/list` (pathStockInstruments) |
| **Status** | 🟡 SDK matches docs summary, not OpenAPI JSON |

Reference: [instrument-list.md](https://developer.webull.hk/apis/docs/reference/instrument-list.md)

### Account List

| | |
|---|---|
| **SDK** | `trade.ListAccounts` |
| **Official (OpenAPI JSON)** | `GET /trading/accounts/list` |
| **Official (llms.txt summary)** | `/openapi/account/list` |
| **SDK path** | `/trading/accounts/list` (pathAccountsList) |
| **Status** | ✅ match |

Reference: [account-list.md](https://developer.webull.hk/apis/docs/reference/account-list.md)

### Account Balance

| | |
|---|---|
| **SDK** | `trade.GetBalance` |
| **Official (OpenAPI JSON)** | `GET /trading/assets/balances/get` |
| **Official (llms.txt summary)** | `/openapi/assets/balance` |
| **SDK path** | `/trading/assets/balances/get` (pathBalanceGet) |
| **Status** | ✅ match |

Reference: [query-account-balance.md](https://developer.webull.hk/apis/docs/reference/query-account-balance.md)

### Account Positions

| | |
|---|---|
| **SDK** | `trade.GetPositions` |
| **Official (OpenAPI JSON)** | `GET /trading/assets/positions/list` |
| **Official (llms.txt summary)** | `/openapi/assets/positions` |
| **SDK path** | `/trading/assets/positions/list` (pathPositionsList) |
| **Status** | ✅ match |

Reference: [query-account-position.md](https://developer.webull.hk/apis/docs/reference/query-account-position.md)

### Order Preview

| | |
|---|---|
| **SDK** | `trade.PreviewOrder` |
| **Official (OpenAPI JSON)** | `POST /trading/orders/preview` |
| **Official (llms.txt summary)** | `/openapi/trade/order/preview` |
| **SDK path** | `/trading/orders/preview` (pathOrdersPreview) |
| **Status** | ✅ match |
| **Note** | Validates and enforces order guardrails. |

Reference: [common-order-preview.md](https://developer.webull.hk/apis/docs/reference/common-order-preview.md)

### Order Place

| | |
|---|---|
| **SDK** | `trade.PlaceOrder` |
| **Official (OpenAPI JSON)** | `POST /trading/orders/place` |
| **Official (llms.txt summary)** | `/openapi/trade/order/place` |
| **SDK path** | `/trading/orders/place` (pathOrdersPlace) |
| **Status** | ✅ match |
| **Note** | Creates live orders. |

Reference: [common-order-place.md](https://developer.webull.hk/apis/docs/reference/common-order-place.md)

### Batch Place Orders

| | |
|---|---|
| **SDK** | `trade.BatchPlaceOrder` |
| **Official (OpenAPI JSON)** | `POST /trading/orders/batch-place` |
| **SDK path** | `/trading/orders/batch-place` (pathOrdersBatchPlace) |
| **Status** | ✅ match |
| **Note** | Equity only, up to 50 orders. |

Reference: [order-batch-place.md](https://developer.webull.com/apis/docs/reference/order-batch-place.md)

### Order Replace

| | |
|---|---|
| **SDK** | `trade.ReplaceOrder` |
| **Official (OpenAPI JSON)** | `POST /trading/orders/replace` |
| **Official (llms.txt summary)** | `/openapi/trade/order/replace` |
| **SDK path** | `/trading/orders/replace` (pathOrdersReplace) |
| **Status** | ✅ match |

Reference: [common-order-replace.md](https://developer.webull.hk/apis/docs/reference/common-order-replace.md)

### Order Cancel

| | |
|---|---|
| **SDK** | `trade.CancelOrder` |
| **Official (OpenAPI JSON)** | `POST /trading/orders/cancel` |
| **Official (llms.txt summary)** | `/openapi/trade/order/cancel` |
| **SDK path** | `/trading/orders/cancel` (pathOrdersCancel) |
| **Status** | ✅ match |

Reference: [common-order-cancel.md](https://developer.webull.hk/apis/docs/reference/common-order-cancel.md)

### Open Orders

| | |
|---|---|
| **SDK** | `trade.GetOpenOrders` |
| **Official (OpenAPI JSON)** | `GET /trading/orders/open-orders/list` |
| **Official (llms.txt summary)** | `/openapi/trade/order/open` |
| **SDK path** | `/trading/orders/open-orders/list` (pathOrdersOpen) |
| **Status** | ✅ match |

Reference: [order-open.md](https://developer.webull.hk/apis/docs/reference/order-open.md)

### Order History

| | |
|---|---|
| **SDK** | `trade.GetOrderHistory` |
| **Official (OpenAPI JSON)** | `GET /trading/orders/historical-orders/list` |
| **Official (llms.txt summary)** | `/openapi/trade/order/history` |
| **SDK path** | `/trading/orders/historical-orders/list` (pathOrdersHistory) |
| **Status** | ✅ match |

Reference: [order-history.md](https://developer.webull.hk/apis/docs/reference/order-history.md)

### Order Detail

| | |
|---|---|
| **SDK** | `trade.GetOrderDetail` |
| **Official (OpenAPI JSON)** | `GET /trading/orders/get` |
| **Official (llms.txt summary)** | `/openapi/trade/order/detail` |
| **SDK path** | `/trading/orders/get` (pathOrdersDetail) |
| **Status** | ✅ match |

Reference: [order-detail.md](https://developer.webull.hk/apis/docs/reference/order-detail.md)

### Cash Activities

| | |
|---|---|
| **SDK** | `trade.GetCashActivities` |
| **Official (OpenAPI JSON)** | `GET /trading/activities/cash-activities/list` |
| **SDK path** | `/trading/activities/cash-activities/list` (pathCashActivitiesList) |
| **Status** | ✅ match |
| **Note** | Documented on the US site. |

Reference: [trade-cash-activity-by-type.md](https://developer.webull.com/apis/docs/reference/trade-cash-activity-by-type.md)

## Broker API — HK

### Create Virtual Account

| | |
|---|---|
| **SDK** | `broker.CreateVirtualAccount` |
| **Official (OpenAPI JSON)** | `POST /broker/accounts/virtual-accounts/create` |
| **SDK path** | `/broker/accounts/virtual-accounts/create` (pathVirtualAccountsCreate) |
| **Status** | ✅ match |

Reference: [broker-account-create.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-account-create.md)

### Update Virtual Account

| | |
|---|---|
| **SDK** | `broker.UpdateVirtualAccount` |
| **Official (OpenAPI JSON)** | `POST /broker/accounts/virtual-accounts/update` |
| **Status** | ❓ SDK path unresolved |

Reference: [broker-account-update.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-account-update.md)

### Get Virtual Account Detail

| | |
|---|---|
| **SDK** | `broker.GetVirtualAccount` |
| **Official (OpenAPI JSON)** | `GET /broker/accounts/virtual-accounts/get` |
| **SDK path** | `/broker/accounts/virtual-accounts/get` (pathVirtualAccountsGet) |
| **Status** | ✅ match |

Reference: [broker-account-detail.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-account-detail.md)

### List Virtual Accounts

| | |
|---|---|
| **SDK** | `broker.ListVirtualAccounts` |
| **Official (OpenAPI JSON)** | `GET /broker/accounts/virtual-accounts/list` |
| **SDK path** | `/broker/accounts/virtual-accounts/list` (pathVirtualAccountsList) |
| **Status** | ✅ match |

Reference: [broker-account-list.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-account-list.md)

### Get Stock Instrument

| | |
|---|---|
| **SDK** | `broker.GetStockInstruments` |
| **Official (OpenAPI JSON)** | `GET /broker/instruments/stocks/profiles/list` |
| **SDK path** | `/broker/instruments/stocks/profiles/list` (pathStockInstruments) |
| **Status** | ✅ match |

Reference: [broker-instrument-list.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-instrument-list.md)

### Get Stock Locate Detail

| | |
|---|---|
| **SDK** | `broker.GetStockLocate` |
| **Official (OpenAPI JSON)** | `GET /broker/instruments/stock-locates/get` |
| **SDK path** | `/broker/instruments/stock-locates/get` (pathStockLocate) |
| **Status** | ✅ match |

Reference: [broker-stock-locate-detail.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-stock-locate-detail.md)

### Get Corporate Actions Detail

| | |
|---|---|
| **SDK** | `broker.GetCorporateActionsDetail` |
| **Official (OpenAPI JSON)** | `GET /broker/instruments/stocks/corporate-actions/get` |
| **SDK path** | `/broker/instruments/stocks/corporate-actions/get` (pathCorporateActionsDetail) |
| **Status** | ✅ match |

Reference: [broker-corporate-actions-detail.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-corporate-actions-detail.md)

### Account Activities By Type

| | |
|---|---|
| **SDK** | `broker.GetCashActivities` |
| **Official (OpenAPI JSON)** | `GET /broker/activities/cash-activities/list` |
| **SDK path** | `/broker/activities/cash-activities/list` (pathActivities) |
| **Status** | ✅ match |

Reference: [broker-activity-by-type.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-activity-by-type.md)

### Account Balance

| | |
|---|---|
| **SDK** | `broker.GetBalance` |
| **Official (OpenAPI JSON)** | `GET /broker/assets/balances/get` |
| **SDK path** | `/broker/assets/balances/get` (pathBalance) |
| **Status** | ✅ match |

Reference: [broker-assets-balance.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-assets-balance.md)

### Account Positions

| | |
|---|---|
| **SDK** | `broker.GetPositions` |
| **Official (OpenAPI JSON)** | `GET /broker/assets/positions/list` |
| **SDK path** | `/broker/assets/positions/list` (pathPositions) |
| **Status** | ✅ match |

Reference: [broker-assets-positions.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-assets-positions.md)

### Order Preview

| | |
|---|---|
| **SDK** | `broker.PreviewOrder` |
| **Official (OpenAPI JSON)** | `POST /broker/orders/preview` |
| **SDK path** | `/broker/orders/preview` (pathOrderPreview) |
| **Status** | ✅ match |

Reference: [broker-order-preview.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-order-preview.md)

### Order Place

| | |
|---|---|
| **SDK** | `broker.PlaceOrder` |
| **Official (OpenAPI JSON)** | `POST /broker/orders/place` |
| **SDK path** | `/broker/orders/place` (pathOrderPlace) |
| **Status** | ✅ match |

Reference: [broker-order-place.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-order-place.md)

### Order Replace

| | |
|---|---|
| **SDK** | `broker.ReplaceOrder` |
| **Official (OpenAPI JSON)** | `POST /broker/orders/replace` |
| **SDK path** | `/broker/orders/replace` (pathOrderReplace) |
| **Status** | ✅ match |

Reference: [broker-order-replace.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-order-replace.md)

### Order Cancel

| | |
|---|---|
| **SDK** | `broker.CancelOrder` |
| **Official (OpenAPI JSON)** | `POST /broker/orders/cancel` |
| **SDK path** | `/broker/orders/cancel` (pathOrderCancel) |
| **Status** | ✅ match |

Reference: [broker-order-cancel.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-order-cancel.md)

### Order Detail

| | |
|---|---|
| **SDK** | `broker.GetOrderDetail` |
| **Official (OpenAPI JSON)** | `GET /broker/orders/get` |
| **SDK path** | `/broker/orders/get` (pathOrderDetail) |
| **Status** | ✅ match |

Reference: [broker-order-detail.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-order-detail.md)

### Order History

| | |
|---|---|
| **SDK** | `broker.GetOrderHistory` |
| **Official (OpenAPI JSON)** | `GET /broker/orders/historical-orders/list` |
| **SDK path** | `/broker/orders/historical-orders/list` (pathOrderHistory) |
| **Status** | ✅ match |

Reference: [broker-order-history.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-order-history.md)

### Open Orders

| | |
|---|---|
| **SDK** | `broker.GetOpenOrders` |
| **Official (OpenAPI JSON)** | `GET /broker/orders/open-orders/list` |
| **SDK path** | `/broker/orders/open-orders/list` (pathOpenOrders) |
| **Status** | ✅ match |

Reference: [broker-order-open.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-order-open.md)

### Get FX Rate

| | |
|---|---|
| **SDK** | `broker.GetFXRate` |
| **Official (OpenAPI JSON)** | `GET /broker/funding/fx-rates/get` |
| **SDK path** | `/broker/funding/fx-rates/get` (pathFXRate) |
| **Status** | ✅ match |

Reference: [broker-funding-query-rate.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-funding-query-rate.md)

### Create FX Request

| | |
|---|---|
| **SDK** | `broker.CreateFXExchange` |
| **Official (OpenAPI JSON)** | `POST /broker/funding/fx-exchanges/create` |
| **SDK path** | `/broker/funding/fx-exchanges/create` (pathFXExchange) |
| **Status** | ✅ match |

Reference: [broker-funding-create-fx.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-funding-create-fx.md)

### Get FX Detail

| | |
|---|---|
| **SDK** | `broker.GetFXExchangeDetail` |
| **Official (OpenAPI JSON)** | `GET /broker/funding/fx-exchanges/get` |
| **SDK path** | `/broker/funding/fx-exchanges/get` (pathFXExchangeDetail) |
| **Status** | ✅ match |

Reference: [broker-funding-query-fx.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-funding-query-fx.md)

### Create Instant Exchange

| | |
|---|---|
| **SDK** | `broker.CreateInstantExchange` |
| **Official (OpenAPI JSON)** | `POST /broker/funding/instant-exchanges/create` |
| **SDK path** | `/broker/funding/instant-exchanges/create` (pathInstantExchange) |
| **Status** | ✅ match |

Reference: [broker-funding-create-instant-fx.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-funding-create-instant-fx.md)

### Get Instant Exchange Detail

| | |
|---|---|
| **SDK** | `broker.GetInstantExchangeDetail` |
| **Official (OpenAPI JSON)** | `GET /broker/funding/instant-exchanges/get` |
| **SDK path** | `/broker/funding/instant-exchanges/get` (pathInstantExchangeDetail) |
| **Status** | ✅ match |

Reference: [broker-funding-query-instant-fx.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-funding-query-instant-fx.md)

### Create Instant Funding

| | |
|---|---|
| **SDK** | `broker.CreateInstantFunding` |
| **Official (OpenAPI JSON)** | `POST /broker/funding/instant-funding/create` |
| **SDK path** | `/broker/funding/instant-funding/create` (pathInstantFunding) |
| **Status** | ✅ match |

Reference: [broker-funding-instant-create.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-funding-instant-create.md)

### Get Instant Funding Detail

| | |
|---|---|
| **SDK** | `broker.GetInstantFundingDetail` |
| **Official (OpenAPI JSON)** | `GET /broker/funding/instant-funding/get` |
| **SDK path** | `/broker/funding/instant-funding/get` (pathInstantFundingDetail) |
| **Status** | ✅ match |

Reference: [broker-funding-instant-query.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-funding-instant-query.md)

### Create Cash Journal

| | |
|---|---|
| **SDK** | `broker.CreateCashJournal` |
| **Official (OpenAPI JSON)** | `POST /broker/journals/cash-journals/create` |
| **SDK path** | `/broker/journals/cash-journals/create` (pathCashJournal) |
| **Status** | ✅ match |

Reference: [broker-journal-cash-create.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-journal-cash-create.md)

### Get Cash Journal Detail

| | |
|---|---|
| **SDK** | `broker.GetCashJournalDetail` |
| **Official (OpenAPI JSON)** | `GET /broker/journals/cash-journals/get` |
| **SDK path** | `/broker/journals/cash-journals/get` (pathCashJournalDetail) |
| **Status** | ✅ match |

Reference: [broker-journal-cash-query.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-journal-cash-query.md)

### Create Position Journal

| | |
|---|---|
| **SDK** | `broker.CreatePositionJournal` |
| **Official (OpenAPI JSON)** | `POST /broker/journals/position-journals/create` |
| **SDK path** | `/broker/journals/position-journals/create` (pathPositionJournal) |
| **Status** | ✅ match |

Reference: [broker-journal-position-create.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-journal-position-create.md)

### Get Position Journal Detail

| | |
|---|---|
| **SDK** | `broker.GetPositionJournalDetail` |
| **Official (OpenAPI JSON)** | `GET /broker/journals/position-journals/get` |
| **SDK path** | `/broker/journals/position-journals/get` (pathPositionJournalDetail) |
| **Status** | ✅ match |

Reference: [broker-journal-position-query.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-journal-position-query.md)

### Trade Calendar

| | |
|---|---|
| **SDK** | `broker.GetTradeCalendar` |
| **Official (OpenAPI JSON)** | `GET /broker/master-data/trading-calendars/list` |
| **SDK path** | `/broker/master-data/trading-calendars/list` (pathTradeCalendar) |
| **Status** | ✅ match |

Reference: [broker-trade-calendar.md](https://developer.webull.hk/apis/docs/reference/broker-api/broker-trade-calendar.md)

### Account Events

| | |
|---|---|
| **SDK** | `—` |
| **Official** | _no OpenAPI schema_ |
| **Status** | ❓ SDK path unresolved |
| **Note** | gRPC subscription. |

Reference: [broker-account-events.md](https://developer.webull.hk/apis/docs/reference/custom/broker-account-events.md)

### Instrument Events

| | |
|---|---|
| **SDK** | `—` |
| **Official** | _no OpenAPI schema_ |
| **Status** | ❓ SDK path unresolved |
| **Note** | gRPC subscription. |

Reference: [broker-instrument-events.md](https://developer.webull.hk/apis/docs/reference/custom/broker-instrument-events.md)

### Corporate Actions Events

| | |
|---|---|
| **SDK** | `—` |
| **Official** | _no OpenAPI schema_ |
| **Status** | ❓ SDK path unresolved |
| **Note** | gRPC subscription. |

Reference: [broker-ca-events.md](https://developer.webull.hk/apis/docs/reference/custom/broker-ca-events.md)

### Trade Events

| | |
|---|---|
| **SDK** | `—` |
| **Official** | _no OpenAPI schema_ |
| **Status** | ❓ no OpenAPI schema on page |
| **Note** | gRPC subscription. |

Reference: [broker-trade-events.md](https://developer.webull.hk/apis/docs/reference/custom/broker-trade-events.md)

### Funding Events

| | |
|---|---|
| **SDK** | `—` |
| **Official** | _no OpenAPI schema_ |
| **Status** | ❓ no OpenAPI schema on page |
| **Note** | gRPC subscription. |

Reference: [broker-funding-events.md](https://developer.webull.hk/apis/docs/reference/custom/broker-funding-events.md)

### Journal Events

| | |
|---|---|
| **SDK** | `—` |
| **Official** | _no OpenAPI schema_ |
| **Status** | ❓ no OpenAPI schema on page |
| **Note** | gRPC subscription. |

Reference: [broker-journal-events.md](https://developer.webull.hk/apis/docs/reference/custom/broker-journal-events.md)

### Master Data Events

| | |
|---|---|
| **SDK** | `—` |
| **Official** | _no OpenAPI schema_ |
| **Status** | ❓ no OpenAPI schema on page |
| **Note** | gRPC subscription. |

Reference: [broker-master-data-events.md](https://developer.webull.hk/apis/docs/reference/custom/broker-master-data-events.md)

## Broker API — FD (US)

### List Accounts

| | |
|---|---|
| **SDK** | `brokerfd.ListFDAccounts` |
| **Official (OpenAPI JSON)** | `GET /broker/accounts/list` |
| **SDK path** | `/broker/accounts/list` (pathFDAccountList) |
| **Status** | ✅ match |

Reference: [list-accounts.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/list-accounts.md)

### Account Detail

| | |
|---|---|
| **SDK** | `brokerfd.GetFDAccountDetail` |
| **Official (OpenAPI JSON)** | `GET /broker/accounts/get` |
| **SDK path** | `/broker/accounts/get` (pathFDAccountDetail) |
| **Status** | ✅ match |

Reference: [get-account-detail.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/get-account-detail.md)

### Create Account

| | |
|---|---|
| **SDK** | `brokerfd.CreateFDAccount` |
| **Official (OpenAPI JSON)** | `POST /broker/accounts/create` |
| **SDK path** | `/broker/accounts/create` (pathFDAccountCreate) |
| **Status** | ✅ match |

Reference: [create-account-apply.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/create-account-apply.md)

### Update Account

| | |
|---|---|
| **SDK** | `brokerfd.UpdateFDAccount` |
| **Official (OpenAPI JSON)** | `POST /broker/accounts/update` |
| **SDK path** | `/broker/accounts/update` (pathFDAccountUpdate) |
| **Status** | ✅ match |

Reference: [update-account-apply.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/update-account-apply.md)

### Close Account

| | |
|---|---|
| **SDK** | `brokerfd.CloseFDAccount` |
| **Official (OpenAPI JSON)** | `POST /broker/accounts/close` |
| **SDK path** | `/broker/accounts/close` (pathFDAccountClose) |
| **Status** | ✅ match |

Reference: [close-account.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/close-account.md)

### Account Application Detail

| | |
|---|---|
| **SDK** | `—` |
| **Official (OpenAPI JSON)** | `GET /broker/accounts/applications/get` |
| **Status** | ❓ SDK path unresolved |

Reference: [get-account-application-detail.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/get-account-application-detail.md)

### List Forms

| | |
|---|---|
| **SDK** | `brokerfd.ListAccountForms` |
| **Official (OpenAPI JSON)** | `GET /broker/forms/list` |
| **SDK path** | `/broker/forms/list` (pathFDAccountForms) |
| **Status** | ✅ match |

Reference: [get-form-list.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/get-form-list.md)

### List Form Versions

| | |
|---|---|
| **SDK** | `—` |
| **Official (OpenAPI JSON)** | `GET /broker/forms/versions/list` |
| **Status** | ❓ SDK path unresolved |

Reference: [get-form-version-list.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/get-form-version-list.md)

### Form Content

| | |
|---|---|
| **SDK** | `—` |
| **Official (OpenAPI JSON)** | `GET /broker/forms/get` |
| **Status** | ❓ SDK path unresolved |

Reference: [get-form-content.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/get-form-content.md)

### Upload Document

| | |
|---|---|
| **SDK** | `brokerfd.UploadDocument` |
| **Official (OpenAPI JSON)** | `POST /broker/documents/upload` |
| **SDK path** | `/broker/documents/upload` (pathDocumentUpload) |
| **Status** | ✅ match |

Reference: [document-upload.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/document-upload.md)

### Download Document

| | |
|---|---|
| **SDK** | `brokerfd.DownloadDocument` |
| **Official (OpenAPI JSON)** | `GET /broker/documents/download` |
| **SDK path** | `/broker/documents/download` (pathDocumentDownload) |
| **Status** | ✅ match |

Reference: [document-download.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/document-download.md)

### Assets Summary

| | |
|---|---|
| **SDK** | `brokerfd.GetAccountsSummary / GetFDAssetsSummary` |
| **Official (OpenAPI JSON)** | `GET /broker/assets/summaries/get` |
| **Status** | ❓ SDK path unresolved |

Reference: [summary.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/summary.md)

### Assets Detail

| | |
|---|---|
| **SDK** | `brokerfd.GetFDAssetsDetail` |
| **Official (OpenAPI JSON)** | `GET /broker/assets/balances/get` |
| **SDK path** | `/broker/assets/balances/get` (pathFDAssetsDetail) |
| **Status** | ✅ match |

Reference: [account-balance.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/account-balance.md)

### Positions

| | |
|---|---|
| **SDK** | `brokerfd.GetFDPositions / GetPositions` |
| **Official (OpenAPI JSON)** | `GET /broker/assets/positions/list` |
| **SDK path** | `/broker/assets/positions/list` (pathFDAssetsPositions) |
| **Status** | ✅ match |

Reference: [account-position.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/account-position.md)

### Cash Activities By Type

| | |
|---|---|
| **SDK** | `brokerfd.GetFDActivities` |
| **Official (OpenAPI JSON)** | `GET /broker/activities/cash-activities/list` |
| **SDK path** | `/broker/activities/cash-activities/list` (pathFDActivities) |
| **Status** | ✅ match |

Reference: [broker-cash-activity-by-type.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-cash-activity-by-type.md)

### Create Bank Relationship

| | |
|---|---|
| **SDK** | `brokerfd.AddFDBankAccount` |
| **Official (OpenAPI JSON)** | `POST /broker/funding/bank-relationships/create` |
| **SDK path** | `/broker/funding/bank-relationships/create` (pathFDBankAccountAdd) |
| **Status** | ✅ match |

Reference: [create-bank-relationship.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/create-bank-relationship.md)

### Delete Bank Relationship

| | |
|---|---|
| **SDK** | `brokerfd.RemoveFDBankAccount` |
| **Official (OpenAPI JSON)** | `POST /broker/funding/bank-relationships/delete` |
| **SDK path** | `/broker/funding/bank-relationships/delete` (pathFDBankAccountRemove) |
| **Status** | ✅ match |

Reference: [delete-bank-relationship.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/delete-bank-relationship.md)

### List Bank Accounts

| | |
|---|---|
| **SDK** | `brokerfd.ListFDBankAccounts` |
| **Official (OpenAPI JSON)** | `GET /broker/funding/bank-relationships/list` |
| **SDK path** | `/broker/funding/bank-relationships/list` (pathFDBankAccounts) |
| **Status** | ✅ match |

Reference: [list-linked-bank-accounts.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/list-linked-bank-accounts.md)

### Create ACH Relationship

| | |
|---|---|
| **SDK** | `brokerfd.AddFDAchAccount` |
| **Official (OpenAPI JSON)** | `POST /broker/funding/ach-relationships/create` |
| **SDK path** | `/broker/funding/ach-relationships/create` (pathFDAchAccountAdd) |
| **Status** | ✅ match |

Reference: [create-ach-relationship.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/create-ach-relationship.md)

### Delete ACH Relationship

| | |
|---|---|
| **SDK** | `brokerfd.RemoveFDAchAccount` |
| **Official (OpenAPI JSON)** | `POST /broker/funding/ach-relationships/delete` |
| **SDK path** | `/broker/funding/ach-relationships/delete` (pathFDAchAccountRemove) |
| **Status** | ✅ match |

Reference: [delete-ach-relationship.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/delete-ach-relationship.md)

### List ACH Relationships

| | |
|---|---|
| **SDK** | `brokerfd.ListFDAchAccounts` |
| **Official (OpenAPI JSON)** | `GET /broker/funding/ach-relationships/list` |
| **SDK path** | `/broker/funding/ach-relationships/list` (pathFDAchAccounts) |
| **Status** | ✅ match |

Reference: [list-ach-relationships.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/list-ach-relationships.md)

### Create Transfer

| | |
|---|---|
| **SDK** | `brokerfd.InitiateFDTransfer` |
| **Official (OpenAPI JSON)** | `POST /broker/funding/transfers/create` |
| **SDK path** | `/broker/funding/transfers/create` (pathFDTransferInitiate) |
| **Status** | ✅ match |

Reference: [create-transfer.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/create-transfer.md)

### Transfer List

| | |
|---|---|
| **SDK** | `brokerfd.ListFDTransfers` |
| **Official (OpenAPI JSON)** | `GET /broker/funding/transfers/list` |
| **SDK path** | `/broker/funding/transfers/list` (pathFDTransfers) |
| **Status** | ✅ match |

Reference: [transfer-list.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/transfer-list.md)

### Transfer Detail

| | |
|---|---|
| **SDK** | `brokerfd.GetFDTransferDetail` |
| **Official (OpenAPI JSON)** | `GET /broker/funding/transfers/get` |
| **SDK path** | `/broker/funding/transfers/get` (pathFDTransferDetail) |
| **Status** | ✅ match |

Reference: [transfer-detail.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/transfer-detail.md)

### Cancel Transfer

| | |
|---|---|
| **SDK** | `—` |
| **Official (OpenAPI JSON)** | `POST /broker/funding/transfers/cancel` |
| **Status** | ❓ SDK path unresolved |

Reference: [cancel-transfer.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/cancel-transfer.md)

### Create Instant Funding

| | |
|---|---|
| **SDK** | `brokerfd.CreateFDInstantFunding` |
| **Official (OpenAPI JSON)** | `POST /broker/funding/instant-funding/create` |
| **SDK path** | `/broker/funding/instant-funding/create` (pathFDInstantFunding) |
| **Status** | ✅ match |

Reference: [broker-funding-instant-create.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-instant-create.md)

### Instant Funding Detail

| | |
|---|---|
| **SDK** | `brokerfd.GetFDInstantFundingDetail` |
| **Official (OpenAPI JSON)** | `GET /broker/funding/instant-funding/get` |
| **SDK path** | `/broker/funding/instant-funding/get` (pathFDInstantFundingDetail) |
| **Status** | ✅ match |

Reference: [broker-funding-instant-query.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-instant-query.md)

### Create Fee

| | |
|---|---|
| **SDK** | `—` |
| **Official (OpenAPI JSON)** | `POST /broker/fees/create` |
| **Status** | ❓ SDK path unresolved |

Reference: [broker-funding-fee-create.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-fee-create.md)

### Get Fee

| | |
|---|---|
| **SDK** | `brokerfd.GetFDTransferFees` |
| **Official (OpenAPI JSON)** | `GET /broker/fees/get` |
| **SDK path** | `/broker/fees/get` (pathFDTransferFees) |
| **Status** | ✅ match |

Reference: [broker-funding-fee-query.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-fee-query.md)

### Create Credit

| | |
|---|---|
| **SDK** | `—` |
| **Official (OpenAPI JSON)** | `POST /broker/credits/create` |
| **Status** | ❓ SDK path unresolved |

Reference: [broker-funding-credit-create.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-credit-create.md)

### Get Credit

| | |
|---|---|
| **SDK** | `brokerfd.GetFDCreditInfo` |
| **Official (OpenAPI JSON)** | `GET /broker/credits/get` |
| **SDK path** | `/broker/credits/get` (pathFDCreditInfo) |
| **Status** | ✅ match |

Reference: [broker-funding-credit-query.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-funding-credit-query.md)

### List Stock Instruments

| | |
|---|---|
| **SDK** | `brokerfd.GetFDStockInstruments` |
| **Official (OpenAPI JSON)** | `GET /broker/instruments/stocks/profiles/list` |
| **SDK path** | `/broker/instruments/stocks/profiles/list` (pathFDStockInstruments) |
| **Status** | ✅ match |

Reference: [list-stock-instruments.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/list-stock-instruments.md)

### Event Contract Categories

| | |
|---|---|
| **SDK** | `—` |
| **Official (OpenAPI JSON)** | `GET /broker/instruments/event-contracts/categories/list` |
| **Status** | ❓ SDK path unresolved |

Reference: [broker-event-categories-list.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-event-categories-list.md)

### Event Contract Series

| | |
|---|---|
| **SDK** | `—` |
| **Official (OpenAPI JSON)** | `GET /broker/instruments/event-contracts/series/list` |
| **Status** | ❓ SDK path unresolved |

Reference: [broker-event-series-list.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-event-series-list.md)

### Event Contract Events

| | |
|---|---|
| **SDK** | `—` |
| **Official (OpenAPI JSON)** | `GET /broker/instruments/event-contracts/events/list` |
| **Status** | ❓ SDK path unresolved |

Reference: [broker-event-events-list.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-event-events-list.md)

### Event Contract Instruments

| | |
|---|---|
| **SDK** | `brokerfd.GetFDECInstruments` |
| **Official (OpenAPI JSON)** | `GET /broker/instruments/event-contracts/markets/list` |
| **SDK path** | `/broker/instruments/event-contracts/markets/list` (pathFDECInstruments) |
| **Status** | ✅ match |

Reference: [broker-event-market-list.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-event-market-list.md)

### Corporate Actions Detail

| | |
|---|---|
| **SDK** | `brokerfd.GetFDCorporateActions` |
| **Official (OpenAPI JSON)** | `GET /broker/instruments/stocks/corporate-actions/get` |
| **SDK path** | `/broker/instruments/stocks/corporate-actions/get` (pathFDCorporateActions) |
| **Status** | ✅ match |

Reference: [broker-corporate-actions-detail.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-corporate-actions-detail.md)

### Order Preview

| | |
|---|---|
| **SDK** | `brokerfd.PreviewFDOrder` |
| **Official (OpenAPI JSON)** | `POST /broker/orders/preview` |
| **SDK path** | `/broker/orders/preview` (pathFDOrderPreview) |
| **Status** | ✅ match |

Reference: [common-order-preview.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/common-order-preview.md)

### Order Place

| | |
|---|---|
| **SDK** | `brokerfd.PlaceFDOrder` |
| **Official (OpenAPI JSON)** | `POST /broker/orders/place` |
| **SDK path** | `/broker/orders/place` (pathFDOrderPlace) |
| **Status** | ✅ match |

Reference: [common-order-place.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/common-order-place.md)

### Order Replace

| | |
|---|---|
| **SDK** | `brokerfd.ReplaceFDOrder` |
| **Official (OpenAPI JSON)** | `POST /broker/orders/replace` |
| **SDK path** | `/broker/orders/replace` (pathFDOrderReplace) |
| **Status** | ✅ match |

Reference: [common-order-replace.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/common-order-replace.md)

### Order Cancel

| | |
|---|---|
| **SDK** | `brokerfd.CancelFDOrder` |
| **Official (OpenAPI JSON)** | `POST /broker/orders/cancel` |
| **SDK path** | `/broker/orders/cancel` (pathFDOrderCancel) |
| **Status** | ✅ match |

Reference: [common-order-cancel.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/common-order-cancel.md)

### Open Orders

| | |
|---|---|
| **SDK** | `brokerfd.GetFDOpenOrders` |
| **Official (OpenAPI JSON)** | `GET /broker/orders/open-orders/list` |
| **SDK path** | `/broker/orders/open-orders/list` (pathFDOrderOpen) |
| **Status** | ✅ match |

Reference: [order-open.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/order-open.md)

### Order Detail

| | |
|---|---|
| **SDK** | `brokerfd.GetFDOrderDetail` |
| **Official (OpenAPI JSON)** | `GET /broker/orders/get` |
| **SDK path** | `/broker/orders/get` (pathFDOrderDetail) |
| **Status** | ✅ match |

Reference: [order-detail.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/order-detail.md)

### Order History

| | |
|---|---|
| **SDK** | `brokerfd.GetFDOrderHistory` |
| **Official (OpenAPI JSON)** | `GET /broker/orders/historical-orders/list` |
| **SDK path** | `/broker/orders/historical-orders/list` (pathFDOrderHistory) |
| **Status** | ✅ match |

Reference: [order-history.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/order-history.md)

### Create Cash Journal

| | |
|---|---|
| **SDK** | `—` |
| **Official (OpenAPI JSON)** | `POST /broker/journals/cash-journals/create` |
| **Status** | ❓ SDK path unresolved |

Reference: [broker-journal-cash-create.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-journal-cash-create.md)

### Cash Journal Detail

| | |
|---|---|
| **SDK** | `brokerfd.GetFDCashJournalDetail` |
| **Official (OpenAPI JSON)** | `GET /broker/journals/cash-journals/get` |
| **SDK path** | `/broker/journals/cash-journals/get` (pathFDCashJournalDetail) |
| **Status** | ✅ match |

Reference: [broker-journal-cash-query.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-journal-cash-query.md)

### List Enums

| | |
|---|---|
| **SDK** | `brokerfd.GetFDEnums` |
| **Official (OpenAPI JSON)** | `GET /broker/master-data/enums/list` |
| **SDK path** | `/broker/master-data/enums/list` (pathFDEnums) |
| **Status** | ✅ match |

Reference: [list-enums.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/list-enums.md)

### Trade Calendar

| | |
|---|---|
| **SDK** | `brokerfd.GetFDTradeCalendar` |
| **Official (OpenAPI JSON)** | `GET /broker/master-data/trading-calendars/list` |
| **SDK path** | `/broker/master-data/trading-calendars/list` (pathFDTradeCalendar) |
| **Status** | ✅ match |

Reference: [list-trade-calendar.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/list-trade-calendar.md)

### List Agreements

| | |
|---|---|
| **SDK** | `brokerfd.ListAgreements` |
| **Official (OpenAPI JSON)** | `GET /broker/agreements/list` |
| **SDK path** | `/broker/agreements/list` (pathAgreements) |
| **Status** | ✅ match |

Reference: [broker-list-agreements-by-type.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-list-agreements-by-type.md)

### Agreement Details

| | |
|---|---|
| **SDK** | `brokerfd.GetAgreementDetail` |
| **Official (OpenAPI JSON)** | `GET /broker/agreements/get` |
| **SDK path** | `/broker/agreements/get` (pathAgreementDetail) |
| **Status** | ✅ match |

Reference: [broker-get-agreement-details.md](https://developer.webull.com/apis/docs/reference/broker-fd-api/broker-get-agreement-details.md)

## Display Solution

### Stock Top Gainers/Losers

| | |
|---|---|
| **SDK** | `data.GetDisplayGainersLosers` |
| **Official (OpenAPI JSON)** | `GET /market-data/screeners/gainers-losers/list` |
| **SDK path** | `/market-data/screeners/gainers-losers/list` (pathDSGainersLosers) |
| **Status** | ✅ match |

Reference: [top-gainers-using-get-new.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/top-gainers-using-get-new.md)

### Top Active

| | |
|---|---|
| **SDK** | `data.GetDisplayTopActive` |
| **Official (OpenAPI JSON)** | `GET /market-data/screeners/top-actives/list` |
| **SDK path** | `/market-data/screeners/top-actives/list` (pathDSTopActive) |
| **Status** | ✅ match |

Reference: [top-active-using-get-new.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/top-active-using-get-new.md)

### Snapshot

| | |
|---|---|
| **SDK** | `data.GetDisplaySnapshot` |
| **Official (OpenAPI JSON)** | `POST /market-data/stocks/snapshots/list` |
| **Official (llms.txt summary)** | `/openapi/market-data/stock/snapshot` |
| **SDK path** | `/openapi/market-data/stock/snapshot` (pathDSSnapshot) |
| **Status** | 🟡 SDK matches docs summary, not OpenAPI JSON |

Reference: [snapshot-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/snapshot-using-get.md)

### Historical Bars (Batch)

| | |
|---|---|
| **SDK** | `data.GetDisplayBars` |
| **Official (OpenAPI JSON)** | `POST /market-data/stocks/bars/list` |
| **SDK path** | `/market-data/stocks/bars/list` (pathDSBars) |
| **Status** | ✅ match |

Reference: [query-batch-bars-using-post.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/query-batch-bars-using-post.md)

### Historical Bars (Single)

| | |
|---|---|
| **SDK** | `data.GetDisplayBarsSingle` |
| **Official (OpenAPI JSON)** | `GET /market-data/stocks/bars/get` |
| **SDK path** | `/market-data/stocks/bars/get` (pathDSBarsSingle) |
| **Status** | ✅ match |

Reference: [bars-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/bars-using-get.md)

### Tick

| | |
|---|---|
| **SDK** | `data.GetDisplayTick` |
| **Official (OpenAPI JSON)** | `GET /market-data/stocks/ticks/list` |
| **SDK path** | `/market-data/stocks/ticks/list` (pathDSTick) |
| **Status** | ✅ match |

Reference: [tick-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/tick-using-get.md)

### Quotes Depth

| | |
|---|---|
| **SDK** | `data.GetDisplayDepth` |
| **Official (OpenAPI JSON)** | `GET /market-data/stocks/depths/list` |
| **SDK path** | `/market-data/stocks/depths/list` (pathDSDepth) |
| **Status** | ✅ match |

Reference: [quotes-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/quotes-using-get.md)

### News Summary

| | |
|---|---|
| **SDK** | `data.GetDSNewsSummary` |
| **Official (OpenAPI JSON)** | `POST /market-data/news/summaries/get` |
| **SDK path** | `/market-data/news/summaries/get` (pathDSNewsSummary) |
| **Status** | ✅ match |

Reference: [watchlist-summary-using-post.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/watchlist-summary-using-post.md)

### Market News

| | |
|---|---|
| **SDK** | `data.GetDSMarketNews` |
| **Official (OpenAPI JSON)** | `GET /market-data/news/market-news/list` |
| **SDK path** | `/market-data/news/market-news/list` (pathDSMarketNews) |
| **Status** | ✅ match |

Reference: [list-news-by-market-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-news-by-market-using-get.md)

### Symbol News

| | |
|---|---|
| **SDK** | `data.GetDSSymbolNews` |
| **Official (OpenAPI JSON)** | `GET /market-data/news/symbol-news/list` |
| **SDK path** | `/market-data/news/symbol-news/list` (pathDSSymbolNews) |
| **Status** | ✅ match |

Reference: [list-news-by-ticker-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-news-by-ticker-using-get.md)

### Latest News

| | |
|---|---|
| **SDK** | `data.GetDSLatestNews` |
| **Official (OpenAPI JSON)** | `GET /market-data/news/latest-news/list` |
| **SDK path** | `/market-data/news/latest-news/list` (pathDSLatestNews) |
| **Status** | ✅ match |

Reference: [list-latest-news-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-latest-news-using-get.md)

### Corporate Actions

| | |
|---|---|
| **SDK** | `data.GetCorporateActions` |
| **Official (OpenAPI JSON)** | `GET /market-data/instruments/stocks/corporate-actions/list` |
| **SDK path** | `/market-data/instruments/stocks/corporate-actions/list` (pathCorpActionsList) |
| **Status** | ✅ match |

Reference: [corp-action-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/corp-action-using-get.md)

### Corporate Actions By Market

| | |
|---|---|
| **SDK** | `data.GetCorporateActionsByMarket` |
| **Official (OpenAPI JSON)** | `GET /market-data/instruments/stocks/corporate-actions/list-by-market` |
| **SDK path** | `/market-data/instruments/stocks/corporate-actions/list-by-market` (pathCorpActionsMarket) |
| **Status** | ✅ match |

Reference: [corp-market-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/corp-market-using-get.md)

### Get Instruments

| | |
|---|---|
| **SDK** | `data.GetStockProfilesV3` |
| **Official (OpenAPI JSON)** | `POST /market-data/instruments/stocks/profiles/list` |
| **SDK path** | `/market-data/instruments/stocks/profiles/list` (pathStockProfilesV3) |
| **Status** | ✅ match |

Reference: [list-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-using-get.md)

### Batch Logos

| | |
|---|---|
| **SDK** | `data.GetLogos` |
| **Official (OpenAPI JSON)** | `POST /market-data/fundamentals/logos/list` |
| **SDK path** | `/market-data/fundamentals/logos/list` (pathLogosBatch) |
| **Status** | ✅ match |

Reference: [batch-logo-using-post.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/batch-logo-using-post.md)

### Company Profile

| | |
|---|---|
| **SDK** | `data.GetDSCompanyProfile` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/company-profiles/get` |
| **SDK path** | `/market-data/fundamentals/company-profiles/get` (pathDSCompanyProfile) |
| **Status** | ✅ match |

Reference: [list-company-profile-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-company-profile-using-get.md)

### Analyst Target Price

| | |
|---|---|
| **SDK** | `data.GetDSAnalystTargetPrice` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/analysis/target-prices/get` |
| **SDK path** | `/market-data/fundamentals/analysis/target-prices/get` (pathDSAnalystTarget) |
| **Status** | ✅ match |

Reference: [list-analyst-target-price-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-analyst-target-price-using-get.md)

### Analyst Rating

| | |
|---|---|
| **SDK** | `data.GetDSAnalystRating` |
| **Official (OpenAPI JSON)** | `GET /market-data/fundamentals/analysis/ratings/get` |
| **SDK path** | `/market-data/fundamentals/analysis/ratings/get` (pathDSAnalystRating) |
| **Status** | ✅ match |

Reference: [list-analyst-rating-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-analyst-rating-using-get.md)

### Streaming Subscribe

| | |
|---|---|
| **SDK** | `data.DSSubscribe` |
| **Official (OpenAPI JSON)** | `POST /market-data/streaming/subscribe` |
| **SDK path** | `/market-data/streaming/subscribe` (pathDSSubscribe) |
| **Status** | ✅ match |
| **Note** | US-site reference. |

Reference: [subscribe-using-post.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/subscribe-using-post.md)

### Streaming Unsubscribe

| | |
|---|---|
| **SDK** | `data.DSUnsubscribe` |
| **Official (OpenAPI JSON)** | `POST /market-data/streaming/unsubscribe` |
| **SDK path** | `/market-data/streaming/unsubscribe` (pathDSUnsubscribe) |
| **Status** | ✅ match |
| **Note** | US-site reference. |

Reference: [unsubscribe-using-post.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/unsubscribe-using-post.md)

## Streaming (MQTT)

### Subscribe

| | |
|---|---|
| **SDK** | `stream.Subscribe` |
| **Official (OpenAPI JSON)** | `POST /market-data/streaming/subscribe` |
| **Official (llms.txt summary)** | `/openapi/market-data/streaming/subscribe` |
| **SDK path** | `/market-data/streaming/subscribe` (subscribePath) |
| **Status** | ✅ match |
| **Note** | HTTP; registers the MQTT session. |

Reference: [subscribe.md](https://developer.webull.hk/apis/docs/reference/subscribe.md)

### Unsubscribe

| | |
|---|---|
| **SDK** | `stream.Unsubscribe` |
| **Official (OpenAPI JSON)** | `POST /market-data/streaming/unsubscribe` |
| **Official (llms.txt summary)** | `/openapi/market-data/streaming/unsubscribe` |
| **SDK path** | `/market-data/streaming/unsubscribe` (unsubscribePath) |
| **Status** | ✅ match |
| **Note** | HTTP; `UnsubscribeAll` clears everything. |

Reference: [unsubscribe.md](https://developer.webull.hk/apis/docs/reference/unsubscribe.md)

## Events (gRPC)

### Subscribe Trade Events

| | |
|---|---|
| **SDK** | `events.New / events.Run` |
| **Official** | _no OpenAPI schema_ |
| **Status** | ❓ no OpenAPI schema on page |
| **Note** | Order status changes. |

Reference: [subscribe-trade-events.md](https://developer.webull.hk/apis/docs/reference/custom/subscribe-trade-events.md)

### Subscribe Position Events

| | |
|---|---|
| **SDK** | `events + SubscribePosition` |
| **Official** | _no OpenAPI schema_ |
| **Status** | ❓ no OpenAPI schema on page |
| **Note** | Position changes. |

Reference: [subscribe-position-events.md](https://developer.webull.com/apis/docs/reference/custom/subscribe-position-events.md)

### Subscribe Events (FD)

| | |
|---|---|
| **SDK** | `brokerfd/events` |
| **Official** | _no OpenAPI schema_ |
| **Status** | ❓ no OpenAPI schema on page |

Reference: [subscribe-events.md](https://developer.webull.com/apis/docs/reference/fd-events/subscribe-events.md)

## Connect API (OAuth)

### Authorization Code

| | |
|---|---|
| **SDK** | `connect.AuthorizationURL` |
| **Official (OpenAPI JSON)** | `GET /oauth2/auth-codes/get` |
| **Status** | ❓ SDK path unresolved |
| **Note** | Browser redirect URL builder. |

Reference: [get-authorization-code.md](https://developer.webull.com/apis/docs/reference/connect-api/get-authorization-code.md)

### Create Token

| | |
|---|---|
| **SDK** | `connect.CreateToken` |
| **Official (OpenAPI JSON)** | `POST /oauth2/tokens/create` |
| **Status** | ❓ SDK path unresolved |

Reference: [create-and-refresh-token.md](https://developer.webull.com/apis/docs/reference/connect-api/create-and-refresh-token.md)

## Documented but not implemented

Unique official endpoints (deduplicated by method and path) that `webullapi4go` does not implement. Duplicate HK/US references to the same endpoint are collapsed into one row.

