# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.3] - 2026-09-18

Market-specific order validation and Hong Kong BCAN support for the `trade`
package.

### Added

- Market-specific equity order-type validation: US accepts limit, market, stop,
  stop-limit, market-on-open, market-on-close, touch, and trailing stop orders;
  HK accepts enhanced limit, at-auction, at-auction limit, stop, stop-limit,
  touch, and trailing stop orders; CN accepts only `LIMIT`.
- Hong Kong BCAN: `HK` equity orders now require at least one `no_party_ids`
  entry with a non-blank `party_id`, `party_id_source` `"D"`, and `party_role`
  `"3"`. A `no_party_ids` list on a non-HK-equity order is rejected.
- US `support_trading_session` validation (`CORE`, `ALL`, `NIGHT`, `ALL_DAY`).
  The deprecated `Y` and `N` aliases and use on a non-US order are rejected.
- At-auction price rules: `AT_AUCTION` rejects `limit_price` and
  `AT_AUCTION_LIMIT` requires it.
- A-share (`CN`) documentation noting that only `LIMIT` is accepted and that
  A-share trading is disabled by default until enabled by Webull support.
- Table-driven tests covering the market matrix, Hong Kong BCAN, US trading
  sessions, and the at-auction price rules.

### Changed

- The trading documentation adds a "Market rules" section, and the `OrderType`,
  `TradingSession`, and `PartyID` GoDoc describe the enforced rules.

## [0.2.2] - 2026-09-18

Stock-order lifecycle and queries for the `trade` package.

### Added

- `trade` order methods: `PreviewOrder` and `PlaceOrder` (with guardrail
  enforcement), `ReplaceOrder`, and `CancelOrder`. Requests are validated before
  any network call and the v3 modify endpoints identify an order by its client
  order ID.
- `trade` order queries: `GetOpenOrders`, `GetOpenOrdersPage`,
  `GetAllOpenOrders`, `GetOrderHistory`, `GetOrderHistoryPage`,
  `GetAllOrderHistory`, and `GetOrderDetail`. Pagination follows the cursor to
  exhaustion, bounded by `trade.MaxOrderQueryPages`.
- Order domain types and enums: `OrderRequest`, `PlaceOrderRequest`,
  `PlaceOrderResult`, `PreviewResult`, `ModifyOrderRequest`,
  `ReplaceOrderRequest`, `ReplaceOrderResult`, `CancelOrderRequest`,
  `CancelOrderResult`, `OrderGroup`, `OrderPage`, `Order`, `OrderLeg`,
  `OrderLegDetail`, `OrderCommission`, `OrderFee`, `OrderHistoryQuery`,
  `OrderSide`, `OrderType`, `TimeInForce`, `ComboType`, `EntrustType`,
  `TradingSession`, `TriggerPriceType`, `TrailingType`, `OrderStatus`, and
  `PartyID`.
- `OrderRequest.Validate`, `PlaceOrderRequest.Validate`,
  `ReplaceOrderRequest.Validate`, and `CancelOrderRequest.Validate`, returning
  typed `invalid_config` errors for the first problem found.
- Runnable `examples/order` program that previews a non-marketable AAPL limit
  buy and, only when `WEBULL_ORDER_PLACE=1` is set, places and cancels it.
- Expanded trading documentation covering the order lifecycle, order-type,
  time-in-force, and combo-type tables, validation rules, guardrails, queries,
  and the explicit warning that placing orders mutates a real account.

## [0.2.1] - 2026-09-18

Trading HTTP foundation: read-only accounts and assets.

### Added

- `trade` package: a typed Trading HTTP client built on the core `client.Client`
  with `trade.New`. Endpoints: `ListAccounts`, `GetBalance`, and `GetPositions`,
  covering account listing and per-account balances and positions. Asset calls
  reject an empty account ID before any network request.
- Trading domain types and enums: `Account`, `AccountType`, `AccountClass`,
  `AssetsBalance`, `AssetsCurrencyAssets`, `Position`, `PositionLeg`, `Market`,
  `InstrumentType`, `OptionType`, and `OptionStrategy`. Numeric fields stay
  strings to preserve precision.
- Order guardrail options `trade.WithMaxOrderNotional` and
  `trade.WithMaxOrderQuantity`, enforced by the order methods before an order is
  built (order placement lands in a later patch).
- Built-in per-path API version defaults: requests under `/trading/` now send
  `x-version: v3`, while Market Data paths remain on `v2`. Overridable with
  `client.WithAPIVersion` and `client.WithAPIVersionFor`.
- Runnable `examples/account` program that lists accounts and prints the balance
  and positions for an account.
- Trading documentation page covering authentication, accounts and assets, the
  v3 default, guardrails, and sandbox usage.

## [0.1.0] - 2026-09-18

Initial public release.

### Added

- Core `client` package: `client.New`, the single signed-request entry point
  `Client.Do`, and functional options for credentials, region, environment,
  endpoints, transport, retries, rate limiting, circuit breaking, API version,
  and environment configuration (`WithEnv`).
- Authentication: HMAC-SHA1 request signing over a percent-encoded canonical
  string, the token lifecycle (`CreateToken`, `CheckToken`, `EnsureToken`,
  `CurrentToken`, `AccessToken`, `SetToken`), automatic token injection, and
  automatic sandbox token acquisition with `WithAutoToken`.
- `data` package: Market Data HTTP endpoints for instruments, company profile,
  analyst data, futures static data, snapshot, tick, quotes/depth, single and
  batch bars, footprint, NOII, screener, watchlist CRUD, options, and news.
- `stream` package: Market Data streaming over MQTT and MQTT-over-WebSocket with
  typed `Quote`, `Snapshot`, and `Tick` handlers, automatic reconnect, and
  idempotent automatic re-subscription.
- Generated protobuf types in `gen/webull/marketdata/v1` for streamed messages.
- `pkg/types` shared domain types.
- Resilient transport primitives: retry with exponential backoff, a keyed token
  bucket rate limiter, and a circuit breaker.
- Documentation site (MkDocs Material) with getting-started, authentication,
  market-data, streaming, sandbox, errors, API reference, and troubleshooting
  pages, plus Architecture Decision Records.
- Runnable examples under `examples/` for auth, market data, streaming, and
  watchlists.

[Unreleased]: https://github.com/shing1211/webullapi4go/compare/v0.2.3...HEAD
[0.2.3]: https://github.com/shing1211/webullapi4go/releases/tag/v0.2.3
[0.2.2]: https://github.com/shing1211/webullapi4go/releases/tag/v0.2.2
[0.2.1]: https://github.com/shing1211/webullapi4go/releases/tag/v0.2.1
[0.1.0]: https://github.com/shing1211/webullapi4go/releases/tag/v0.1.0
