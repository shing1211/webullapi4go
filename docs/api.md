# API Reference

The generated reference documentation lives on pkg.go.dev:

| Package | Reference |
|---------|-----------|
| Core client | https://pkg.go.dev/github.com/shing1211/webullapi4go/client |
| Market Data (HTTP) | https://pkg.go.dev/github.com/shing1211/webullapi4go/data |
| Market Data (streaming) | https://pkg.go.dev/github.com/shing1211/webullapi4go/stream |
| Trading (HTTP) | https://pkg.go.dev/github.com/shing1211/webullapi4go/trade |
| Trading events (gRPC) | https://pkg.go.dev/github.com/shing1211/webullapi4go/events |
| Broker FD (HTTP) | https://pkg.go.dev/github.com/shing1211/webullapi4go/brokerfd |
| Broker FD event protobuf types | https://pkg.go.dev/github.com/shing1211/webullapi4go/gen/webull/brokerfd/v1 |
| Streamed protobuf types | https://pkg.go.dev/github.com/shing1211/webullapi4go/gen/webull/marketdata/v1 |
| Event protobuf types | https://pkg.go.dev/github.com/shing1211/webullapi4go/gen/webull/trade/events/v1 |
| Broker FD event protobuf types | https://pkg.go.dev/github.com/shing1211/webullapi4go/gen/webull/brokerfd/events/v1 |
| Broker API HK | https://pkg.go.dev/github.com/shing1211/webullapi4go/broker |
| Broker FD events | https://pkg.go.dev/github.com/shing1211/webullapi4go/brokerfd/events |
| Connect API (OAuth) | https://pkg.go.dev/github.com/shing1211/webullapi4go/connect |
| Display Solution | https://pkg.go.dev/github.com/shing1211/webullapi4go/display |
| Shared domain types | https://pkg.go.dev/github.com/shing1211/webullapi4go/pkg/types |
| Module root | https://pkg.go.dev/github.com/shing1211/webullapi4go |

Every exported identifier carries GoDoc comments. This page is a map of the
public surface; the pkg.go.dev pages are authoritative.

## `client`

The core SDK. Construct one value and share it across API packages.

Constructor and transport:

- `client.New(opts ...Option) (*Client, error)`
- `Client.Do(ctx, method, path, body, out) error` — the single signed-request entry point
- `Client.DoStream(ctx, method, path, body) (*http.Response, error)` — signed
  streaming request that returns the open response for endpoints such as
  Server-Sent Events; applies the rate limiter and circuit breaker, is never
  retried, and the caller closes the response body
- `Client.Close() error`
- Accessors: `Config`, `Region`, `Environment`, `Endpoints`, `HTTPClient`

Token lifecycle:

- `Client.CreateToken`, `Client.CheckToken`, `Client.EnsureToken`
- `Client.CurrentToken`, `Client.AccessToken`, `Client.SetToken`
- `Client.SetTokenPollInterval`, `Client.SetTokenPollTimeout`
- `Client.EnableTokenInjection`
- Sentinel: `client.ErrAccessTokenRequired`

Configuration options:

- Credentials and region: `WithAppKey`, `WithAppSecret`, `WithCredentials`,
  `WithRegion`, `WithEnvironment`, `WithSandbox`, `WithEnv`
- Endpoints and transport: `WithEndpoints`, `WithBaseURL`, `WithHTTPClient`,
  `WithTimeout`, `WithUserAgent`
- Resilience: `WithRetry`, `WithoutRetry`, `NewRateLimiter`, `WithRateLimiter`,
  `NewBreaker`, `WithBreaker`
- Versioning and tokens: `WithAPIVersion`, `WithAPIVersionFor`, `WithAutoToken`

Sentinels and accessors: `client.ErrCircuitOpen`, `client.AccessTokenHeader`.

## `data`

The Market Data HTTP client. Construct it with `data.New(*client.Client)`.

- Instruments: `GetStockInstruments`
- Instrument v3 (Display Solution): `GetStockProfilesV3` (batch, POST)
- Logos: `GetLogos` (batch, POST)
- Corporate actions: `GetCorporateActions`, `GetCorporateActionsByMarket`
- Profile and analyst: `GetCompanyProfile`, `GetAnalystTargetPrice`, `GetAnalystRating`
- Futures static data: `GetFuturesInstruments`, `GetFuturesProductCodes`, `GetFuturesProductClasses`
- Snapshot and quotes: `GetSnapshot`, `GetQuotes`
- Ticks and bars: `GetTick`, `GetBars`, `GetBatchBars`
- Depth analytics: `GetFootprint`, `GetNOIIBars`, `GetNOIISnapshot`
- Discovery: `GetTopGainersLosers`, `GetMostActive`
- Watchlists: `GetWatchlists`, `CreateWatchlist`, `UpdateWatchlist`,
  `DeleteWatchlist`, `GetWatchlistInstruments`, `AddWatchlistInstruments`,
  `RemoveWatchlistInstruments`, `UpdateWatchlistInstruments`
- Derivatives and news: `GetOptionTick`, `GetOptionSnapshot`, `GetOptionBars`,
  `GetNewsSummary`
- Option contracts: `GetOptionContracts` with `OptionContractsQuery`,
  `OptionContract`, and `OptionType` (`Call`/`Put`); path follows the official
  definition (HK sandbox returns `404` — US-only surface)
- Fund data: `GetFundNav`, `GetFundInfo`, `GetFundDividends`, `GetFundList`,
  plus `GetFundPerformance`, `GetFundHoldings`, `GetFundRating`, `GetFundSplits`,
  `GetFundFiles`, `GetFundAllocation` (HK sandbox returns `404`)
- Fundamentals: `GetCapitalFlow`, `GetIndustryComparison`,
  `GetEarningsCalendar`, `GetDividendCalendar`, `GetFilings`,
  `GetIncomeStatement`, `GetBalanceSheet`, `GetCashFlow`,
  `GetFinancialIndicators`, `GetFinancialAlert`, `GetForecastEPS`

The package also exposes its query and response types (for example
`SnapshotQuery`, `BarQuery`, `Watchlist`) and the `StockCategory`,
`BarTimespan`, and `TradingSession` enumerations.

## `stream`

The Market Data streaming client. Construct it with `stream.New(*client.Client, ...Option)`.

- Lifecycle: `Connect`, `Close`, `IsConnected`, `Reconnecting`, `SessionID`
- Subscriptions: `Subscribe`, `Unsubscribe`
- Handlers: `OnQuote`, `OnSnapshot`, `OnTick`, `OnNotice`, `OnError`,
  `OnConnect`, `OnDisconnect`
- Types: `SubscribeRequest`, `UnsubscribeRequest`, `Category`, `SubType`
- Topic constants: `TopicQuote`, `TopicSnapshot`, `TopicTick`, `TopicNotice`, `TopicEcho`
- Options: `WithSessionID`, `WithClientID`, `WithMQTTURL`, `WithWebSocket`,
  `WithAutoReconnect`, `WithAutoResubscribe`, `WithResubscribeTimeout`,
  `WithKeepAlive`, `WithConnectTimeout`, `WithWriteTimeout`,
  `WithMessageChannelDepth`, `WithCleanSession`, `WithTLSConfig`

## `trade`

The Trading HTTP client. Construct it with `trade.New(*client.Client, ...Option)`.
Requests require an access token and default to the v3 API version under
`/trading/`.

- Accounts and assets: `ListAccounts`, `GetBalance`, `GetPositions`
- Order lifecycle: `PreviewOrder`, `PlaceOrder`, `ReplaceOrder`, `CancelOrder`
- Order queries: `GetOpenOrders`, `GetOpenOrdersPage`, `GetAllOpenOrders`,
  `GetOrderHistory`, `GetOrderHistoryPage`, `GetAllOrderHistory`,
  `GetOrderDetail`
- Order request and response types: `OrderRequest`, `PlaceOrderRequest`,
  `PlaceOrderResult`, `PreviewResult`, `ModifyOrderRequest`,
  `ReplaceOrderRequest`, `ReplaceOrderResult`, `CancelOrderRequest`,
  `CancelOrderResult`, `OrderGroup`, `OrderPage`, `Order`, `OrderLeg`,
  `OrderLegDetail`, `OrderCommission`, `OrderFee`, `OrderHistoryQuery`
- Enumerations: `OrderSide`, `OrderType`, `TimeInForce`, `ComboType`,
  `EntrustType`, `TradingSession`, `TriggerPriceType`, `TrailingType`,
  `OrderStatus`
- Multi-leg `OptionStrategy` values (`VERTICAL` through `RATIO`) and
  `InstrumentTypeFutures` order validation follow the official documentation;
  the HK sandbox accepts only `SINGLE` strategies (`417` otherwise) — see
  [Trading](trading.md#multi-leg-orders) and [Trading](trading.md#futures-orders)
- Types: `Account`, `AccountType`, `AccountClass`, `AssetsBalance`,
  `AssetsCurrencyAssets`, `Position`, `PositionLeg`, `Market`, `InstrumentType`,
  `OptionType`, `OptionStrategy`, `PartyID`
- Options: `WithMaxOrderNotional`, `WithMaxOrderQuantity` (order guardrails,
  enforced by `PreviewOrder` and `PlaceOrder`)
- Bounds: `MaxOrderQueryPages`

## `events`

The Trading events client over gRPC. Construct it with
`events.New(*client.Client, ...Option)`.

- Lifecycle: `New`, `Run`, `Close`
- Handlers: `OnConnect`, `OnPing`, `OnEvent`, `OnOrder`, `OnPosition`,
  `OnOption`, `OnError`
- Types: `OrderEvent`, `PositionEvent`, `OptionEvent`, `SubscribeType`
- Subscribe bitmask: `SubscribeOrder`, `SubscribePosition`, `SubscribeOption`,
  `SubscribeAll`
- Data event kinds: `EventOrder`, `EventPosition`, `EventOption`
- Options: `WithGRPCEndpoint`, `WithGRPCPort`, `WithTLS`, `WithDialTimeout`,
  `WithGRPCDialOption`, `WithSubscribeTypes`, `WithAccounts`,
  `WithAutoReconnect`, `WithReconnectBaseDelay`, `WithReconnectMaxDelay`,
  `WithMaxReconnectAttempts`
- Defaults: `DefaultGRPCPort`, `DefaultDialTimeout`, `DefaultReconnectBaseDelay`,
  `DefaultReconnectMaxDelay`

The event service signs each `Subscribe` call with HMAC-SHA256 over the
serialized request, with no `host` in the canonical string. See
[Trading events](events.md) for the signing rules, dispatch model, and payload
schemas.

## `brokerfd`

The Broker FD (Fund Data) HTTP client. Construct it with `brokerfd.New(*client.Client)`.

- Accounts: `GetAccountsSummary`, `GetPositions`, plus orders, assets,
  instruments, funding, activity, journals, master data, agreements, and
  documents endpoints (paths follow the official definition; US-only — the HK
  sandbox returns `404`)
- Event types (gRPC, in `gen/webull/brokerfd/v1`): `BrokerFDEventTypeAccountPush`,
  `BrokerFDEventTypeOrderPush`, `BrokerFDEventTypePositionPush`, `BrokerFDEventTypeTradePush`,
  `BrokerFDEventTypeAssetDetail`, `BrokerFDEventTypeRiskPush`, `BrokerFDEventTypeOrderFill`,
  `BrokerFDEventTypePositionSync`, `BrokerFDEventTypeAssetSync`, `BrokerFDEventTypeAssetsPush`,
  `BrokerFDEventTypeOrdersPush`, `BrokerFDEventTypeCashSecLiab`

## `gen/webull/marketdata/v1`

Generated protobuf types for streamed messages (`Quote`, `Snapshot`, `Tick`,
`Basic`, `AskBid`, and so on) plus the `DecodeQuote`, `DecodeSnapshot`, and
`DecodeTick` helpers. Generated code is committed so builds do not require a
protobuf toolchain.

## `gen/webull/brokerfd/v1`

Generated protobuf types for Broker FD gRPC events (`BrokerFDEvent`,
`BrokerFDEventType`, `AccountPush`, `OrderPush`, `PositionPush`, `TradePush`,
`AssetDetail`, `RiskPush`, `OrderFillPush`, `PositionSync`, `AssetSync`,
`AssetsPush`, `OrdersPush`, `CashSecLiability`). Schema is best-effort;
live probe required to confirm field names and types.

## `pkg/types`

Shared public domain types that are safe to import from outside the module, such
as `Market` and `InstrumentType`.

## `internal/*`

Implementation details: request signing and token lifecycle, region endpoints,
the HTTP transport, resilience primitives (retry, rate limit, circuit breaker),
and the low-level MQTT client. `internal/auth` carries an algorithm-parameterised
signer: the REST API uses HMAC-SHA1 over a canonical string that includes the
`host` and an upper-cased MD5 body digest, while the gRPC event service uses
HMAC-SHA256 with no `host` and a lower-cased SHA-256 body digest. Internal
packages are not part of the public API and are never referenced by exported
signatures; do not import them.
