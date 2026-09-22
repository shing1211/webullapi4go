# Webull API Master Reference

This page maps every official [Webull OpenAPI](https://developer.webull.hk/apis/docs/) endpoint to its corresponding `webullapi4go` function, grouped by API category.

## Official API Documentation

Base URL: <https://developer.webull.hk/apis/docs/>

---

## Market Data API

### Non-Display Solution

#### Stock

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /market-data/quotes/snapshot` | `data.GetSnapshot` | Batch snapshot |
| `GET /market-data/quotes/tick` | `data.GetTick` | Single symbol tick |
| `GET /market-data/bars/list` | `data.GetBars` | Single symbol historical bars |
| `POST /market-data/bars/batch` | `data.GetBatchBars` | Multi-symbol historical bars |
| `GET /market-data/depths/list` | `data.GetQuotes` | Order book depth |
| `GET /market-data/footprints/list` | `data.GetFootprint` | Footprint bars |
| `GET /market-data/nioi/bars` | `data.GetNOIIBars` | NOI bars |
| `GET /market-data/nioi/snapshot` | `data.GetNOIISnapshot` | NOI snapshot |

#### Options

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /market-data/options/ticks/list` | `data.GetOptionTick` | Option tick |
| `GET /market-data/options/snapshots/list` | `data.GetOptionSnapshot` | Option snapshot (batch) |
| `GET /market-data/options/bars/list` | `data.GetOptionBars` | Option historical bars |

#### Futures

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /market-data/futures/ticks/list` | `data.GetFuturesTick` | Futures tick |
| `GET /market-data/futures/snapshots/list` | `data.GetFuturesSnapshot` | Futures snapshot |
| `GET /market-data/futures/bars/list` | `data.GetFuturesBars` | Futures historical bars |
| `GET /market-data/futures/depths/list` | `data.GetFuturesDepth` | Futures order book depth |
| `GET /market-data/futures/footprints/list` | `data.GetFuturesFootprint` | Futures footprint bars |

#### News

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /market-data/news` | `data.GetNewsSummary` | SSE stream news summary |

#### Screeners

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /market-data/screener/tickers` | `data.GetTopGainersLosers` | Top gainers/losers |
| `GET /market-data/screener/top-active` | `data.GetMostActive` | Top actives |

#### Watchlist

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /market-data/vip/wetal/simplelist` | `data.GetWatchlists` | List watchlists |
| `POST /market-data/vip/wetal/simplelist` | `data.CreateWatchlist` | Create watchlist |
| `PUT /market-data/vip/wetal/simplelist` | `data.UpdateWatchlist` | Update watchlist |
| `DELETE /market-data/vip/wetal/simplelist` | `data.DeleteWatchlist` | Delete watchlist |
| `GET /market-data/vip/wetal/stock/list` | `data.GetWatchlistInstruments` | Get watchlist instruments |
| `POST /market-data/vip/wetal/stock/list` | `data.AddWatchlistInstruments` | Add instruments |
| `DELETE /market-data/vip/wetal/stock/list` | `data.RemoveWatchlistInstruments` | Remove instruments |
| `PUT /market-data/vip/wetal/stock/list` | `data.UpdateWatchlistInstruments` | Update instruments |

#### Fundamentals

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /market-data/fundamental/capital-flow` | `data.GetCapitalFlow` | Capital flow |
| `GET /market-data/fundamental/industry/compare` | `data.GetIndustryComparison` | Industry comparison |
| `GET /market-data/fundamental/earnings-calendar` | `data.GetEarningsCalendar` | Earnings calendar |
| `GET /market-data/fundamental/dividend-calendar` | `data.GetDividendCalendar` | Dividend calendar |
| `GET /market-data/fundamental/filings` | `data.GetFilings` | SEC filings |
| `GET /market-data/fundamental/financial/income` | `data.GetIncomeStatement` | Income statement |
| `GET /market-data/fundamental/financial/balance-sheet` | `data.GetBalanceSheet` | Balance sheet |
| `GET /market-data/fundamental/financial/cash-flow` | `data.GetCashFlow` | Cash flow |
| `GET /market-data/fundamental/financial/indicator` | `data.GetFinancialIndicators` | Financial indicators |
| `GET /market-data/fundamental/financial/alert` | `data.GetFinancialAlert` | Financial alert |
| `GET /market-data/fundamental/forecast/eps` | `data.GetForecastEPS` | Forecast EPS |

#### Fund Data

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /market-data/fund/funds/{ticker}/nav` | `data.GetFundNav` | Fund NAV |
| `GET /market-data/fund/funds/{ticker}/info` | `data.GetFundInfo` | Fund info |
| `GET /market-data/fund/funds/{ticker}/dividend` | `data.GetFundDividends` | Fund dividends |
| `GET /market-data/fund/funds/list` | `data.GetFundList` | Fund list |

!!! warning
    Fund data endpoints are **unconfirmed** against live API. The HK sandbox returns `404`.

---

### Display Solution

#### Quotes

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /market-data/v1/snapshot` | `data.GetStockSnapshot` | Stock snapshot (Display) |
| `POST /market-data/v1/bars` | `data.GetStockHistoricalBars` | Stock historical bars (Display) |
| `GET /market-data/v1/tick` | `data.GetStockTick` | Stock tick (Display) |
| `GET /market-data/v1/quotes` | `data.GetStockDepth` | Stock depth (Display) |

#### Corporate Actions

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /market-data/v1/corp-action` | `data.GetCorporateActionsList` | Corporate actions |
| `GET /market-data/v1/corp-action/market` | `data.GetCorporateActionsMarket` | Market corporate actions |

#### Instruments

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /market-data/v1/instruments` | `data.GetStockInstruments` | Stock instruments |
| `POST /market-data/v1/logo` | `data.GetLogosBatch` | Batch logos |
| `GET /market-data/v1/company-profile` | `data.GetCompanyProfile` | Company profile |
| `GET /market-data/v1/analyst-target-price` | `data.GetAnalystTargetPrice` | Analyst target price |
| `GET /market-data/v1/analyst-rating` | `data.GetAnalystRating` | Analyst rating |

#### Screeners (Display Solution)

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /market-data/screener/tickers` | `data.GetTopGainersLosers` | Top gainers/losers |
| `GET /market-data/screener/top-active` | `data.GetMostActive` | Top actives |
| `GET /market-data/screener/sectors` | `data.GetMarketSectors` | Market sectors |
| `GET /market-data/screener/sector-detail` | `data.GetMarketSectorDetail` | Sector detail |
| `GET /market-data/screener/high-dividend` | `data.GetHighDividendRank` | High dividend rank |
| `GET /market-data/screener/week-52` | `data.GetWeek52HighLow` | 52-week high/low |

#### News (Display Solution)

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /market-data/news` | `data.GetNewsSummary` | News summary (SSE) |

---

## Trading API

### Instruments

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /trading/instruments/options/contracts/list` | `data.GetOptionContracts` | Option contracts (chain) |

!!! warning
    Field mappings and wire values are **unconfirmed**. HK sandbox returns `404`. US sandbox credentials required.

### Futures Instruments

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /trading/future/options/information` | `data.GetFuturesInstruments` | Futures instruments |
| `GET /trading/future/options/product-codes` | `data.GetFuturesProductCodes` | Product codes |
| `GET /trading/future/options/product-classes` | `data.GetFuturesProductClasses` | Product classes |

### Order Lifecycle

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `POST /trading/v5/order/oderpreview` | `trade.PreviewOrder` | Preview order |
| `POST /trading/v5/order/place` | `trade.PlaceOrder` | Place order |
| `POST /trading/v5/order/cancel` | `trade.CancelOrder` | Cancel order |
| `POST /trading/v5/order/modify` | `trade.ReplaceOrder` | Replace order |

### Order Queries

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /trading/v5/order/query/current` | `trade.GetOpenOrders` | Current open orders |
| `GET /trading/v5/order/query/history` | `trade.GetOrderHistory` | Order history |

### Account & Asset

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /trading/v5/account/陶0` | `trade.ListAccounts` | List accounts |
| `GET /trading/v5/account/balance` | `trade.GetBalance` | Account balance |
| `GET /trading/v5/position/list` | `trade.GetPositions` | Positions |

---

## Broker API

### Broker API HK

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /broker/v1/account/陶0` | `broker.GetBrokerAccount` | Broker account |
| `GET /broker/v1/position/list` | `broker.GetBrokerPositions` | Positions |
| `GET /broker/v1/order/query/current` | `broker.GetBrokerOpenOrders` | Open orders |
| `GET /broker/v1/order/query/history` | `broker.GetBrokerOrderHistory` | Order history |
| `POST /broker/v1/order/place` | `broker.PlaceBrokerOrder` | Place order |
| `POST /broker/v1/order/cancel` | `broker.CancelBrokerOrder` | Cancel order |

!!! warning
    Broker API HK returns `401 ROUTE_NOT_PERMITTED` in HK sandbox — the app lacks the required scope.

### Broker FD API US

| Official Endpoint | SDK Function | Notes |
|-----------------|--------------|-------|
| `GET /brokerfd/v1/account/summary` | `brokerfd.GetAccountsSummary` | Account summary |
| `GET /brokerfd/v1/position/list` | `brokerfd.GetPositions` | Positions |
| Event service (gRPC) | `brokerfd/events.New` | Push events |

!!! warning
    Broker FD API US returns `404` in HK sandbox — US-only endpoint.

---

## Provisional (Unconfirmed Endpoints)

The following are present in the SDK but have **not been confirmed** against the live Webull API. They are marked `TODO` in code and may change.

| Function | Issue |
|----------|-------|
| `GetOptionContracts` field mappings | HK sandbox returns `404`; requires US sandbox |
| `OptionStrategy` values (`VERTICAL`, `CALENDAR`, `STRANGLE`, `BUTTERFLY`, `RATIO`, `CONDOR`, `DIAGONAL`) | Wire values unconfirmed; only `SINGLE` confirmed in official docs |
| `InstrumentTypeFutures` order support | Futures order validation provisional |
| Fund data endpoints (`GetFundNav`, `GetFundInfo`, `GetFundDividends`, `GetFundList`) | HK sandbox returns `404` |

---

## Removed Functions

The following were removed in v1.0.2 because they are **not documented** in the official Webull OpenAPI:

| Removed Function | Reason |
|-----------------|---------|
| `GetOptionExpirations` | No such endpoint in official API; option chain response includes expirations |
| `GetOptionChain` | Was using wrong path (`/market-data/options/contracts/list`); correct path is `/trading/instruments/options/contracts/list` |
| `GetHKOptionExpirations` | Same as above — no standalone expiration endpoint |
| `GetHKOptionChain` | Same as above |
| `GetHKFuturesTick`, `GetHKFuturesSnapshot`, `GetHKFuturesBars`, `GetHKFuturesDepth`, `GetHKFuturesFootprint` | Redundant with non-HK variants; same paths with `category=HK` |
| `GetCryptoBars`, `GetCryptoTick`, `GetCryptoDepth`, `GetCryptoSnapshot`, `GetCryptoProfile`, `GetCryptoList` | Not in official API; HK sandbox returns `417` |
| `GetScreenerV2` | Not in official API; HK sandbox returns `404` |

---

## Reference

- Official Webull OpenAPI Docs: <https://developer.webull.hk/apis/docs/>
- SDK pkg.go.dev: <https://pkg.go.dev/github.com/shing1211/webullapi4go>
- This document: `docs/webull-api.md`
