# API Reference

The generated reference documentation lives on pkg.go.dev:

| Package | Reference |
|---------|-----------|
| Core client | https://pkg.go.dev/github.com/shing1211/webullapi4go/client |
| Market Data (HTTP) | https://pkg.go.dev/github.com/shing1211/webullapi4go/data |
| Market Data (streaming) | https://pkg.go.dev/github.com/shing1211/webullapi4go/stream |
| Trading (HTTP) | https://pkg.go.dev/github.com/shing1211/webullapi4go/trade |
| Streamed protobuf types | https://pkg.go.dev/github.com/shing1211/webullapi4go/gen/webull/marketdata/v1 |
| Shared domain types | https://pkg.go.dev/github.com/shing1211/webullapi4go/pkg/types |
| Module root | https://pkg.go.dev/github.com/shing1211/webullapi4go |

Every exported identifier carries GoDoc comments. This page is a map of the
public surface; the pkg.go.dev pages are authoritative.

## `client`

The core SDK. Construct one value and share it across API packages.

Constructor and transport:

- `client.New(opts ...Option) (*Client, error)`
- `Client.Do(ctx, method, path, body, out) error` — the single signed-request entry point
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
- Types: `Account`, `AccountType`, `AccountClass`, `AssetsBalance`,
  `AssetsCurrencyAssets`, `Position`, `PositionLeg`, `Market`, `InstrumentType`,
  `OptionType`, `OptionStrategy`
- Options: `WithMaxOrderNotional`, `WithMaxOrderQuantity` (order guardrails)

Order placement, replacement, cancellation, and order queries are added in later
v0.2 patches.

## `gen/webull/marketdata/v1`

Generated protobuf types for streamed messages (`Quote`, `Snapshot`, `Tick`,
`Basic`, `AskBid`, and so on) plus the `DecodeQuote`, `DecodeSnapshot`, and
`DecodeTick` helpers. Generated code is committed so builds do not require a
protobuf toolchain.

## `pkg/types`

Shared public domain types that are safe to import from outside the module, such
as `Market` and `InstrumentType`.

## `internal/*`

Implementation details: request signing and token lifecycle, region endpoints,
the HTTP transport, resilience primitives (retry, rate limit, circuit breaker),
and the low-level MQTT client. Internal packages are not part of the public API
and are never referenced by exported signatures; do not import them.
