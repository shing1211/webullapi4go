# Market Data

Market Data is available over two transports: HTTP for on-demand queries and
MQTT for real-time streaming. Both share the same authentication described in
[Authentication](authentication.md).

!!! note "Prerequisites"
    - A [Webull account](https://developer.webull.hk/apis/docs/sdk#test-accounts) (sandbox or production)
    - Go 1.26+
    - An authenticated client — see [Authentication](authentication.md)

!!! tip "Error handling"
    All SDK functions return `error`. See [Errors](errors.md) for the typed error model, transient vs permanent classification, and retry patterns.

## HTTP

The `data` package exposes the HTTP Market Data API. Build it from the public
core client:

```go
cl, err := client.New(client.WithEnv())
if err != nil {
	return err
}
defer func() { _ = cl.Close() }()

if _, err := cl.EnsureToken(ctx); err != nil {
	return err
}

market := data.New(cl)
```

| Group | Methods |
|-------|---------|
| Instruments | `GetStockInstruments` |
| Instrument v3 (Display) | `GetStockProfilesV3` |
| Logos (Display) | `GetLogos` |
| Corporate actions (Display) | `GetCorporateActions`, `GetCorporateActionsByMarket` |
| Profile and analyst | `GetCompanyProfile`, `GetAnalystTargetPrice`, `GetAnalystRating` |
| Futures static data | `GetFuturesInstruments`, `GetFuturesProductCodes`, `GetFuturesProductClasses` |
| Snapshot and quotes | `GetSnapshot`, `GetQuotes` |
| Ticks and bars | `GetTick`, `GetBars`, `GetBatchBars` |
| Depth analytics | `GetFootprint`, `GetNOIIBars`, `GetNOIISnapshot` |
| Discovery | `GetTopGainersLosers`, `GetMostActive` |
| Watchlists | `GetWatchlists`, `CreateWatchlist`, `UpdateWatchlist`, `DeleteWatchlist`, `GetWatchlistInstruments`, `AddWatchlistInstruments`, `RemoveWatchlistInstruments`, `UpdateWatchlistInstruments` |
| Derivatives and news | `GetOptionTick`, `GetOptionSnapshot`, `GetOptionBars`, `GetOptionContracts`, `GetNewsSummary` |
| Event contracts | `GetEventContractCategories`, `GetEventContractSeries`, `GetEventContractEvents`, `GetEventContractMarkets` |
| Event contract market data | `GetEventSnapshot`, `GetEventDepth`, `GetEventBars`, `GetEventTick` |
| Event contracts display | `GetEventContractTags`, `GetEventContractEventsList`, `GetEventContractMilestones`, `GetEventContractSeriesList`, `GetEventContractSportsFilters`, `GetEventGameStats`, `GetEventLiveData`, `GetEventMarketBars`, `GetEventMarketBarsByEvent`, `GetEventMarketDepth`, `GetEventMarketSnapshot` |
| Futures market data | `GetFuturesTick`, `GetFuturesSnapshot`, `GetFuturesBars`, `GetFuturesDepth`, `GetFuturesFootprint` |
| Fund data | `GetFundNav`, `GetFundInfo`, `GetFundDividends`, `GetFundList` |
| Fund extras | `GetFundPerformance`, `GetFundHoldings`, `GetFundRating`, `GetFundSplits`, `GetFundFiles`, `GetFundAllocation` |
| Non-display screener | `GetMarketSectors`, `GetMarketSectorDetail`, `GetHighDividendRank`, `GetWeek52HighLow` |
| Crypto US | `GetCryptoSnapshot`, `GetCryptoBars`, `GetCryptoInstruments` |
| Fundamentals | `GetCapitalFlow`, `GetIndustryComparison`, `GetEarningsCalendar`, `GetDividendCalendar`, `GetFilings`, `GetIncomeStatement`, `GetBalanceSheet`, `GetCashFlow`, `GetFinancialIndicators`, `GetFinancialAlert`, `GetForecastEPS` |

`GetOptionContracts` follows the official option-contract list endpoint. The HK
sandbox may return `404` for this US-only surface.

Example: snapshot and bars for `AAPL` on the `US` market.

```go
snaps, err := market.GetSnapshot(ctx, data.SnapshotQuery{
	Symbols:  []string{"AAPL"},
	Category: data.StockCategoryUS,
})
if err != nil {
	return err
}

bars, err := market.GetBars(ctx, data.BarQuery{
	Symbol:   "AAPL",
	Category: data.StockCategoryUS,
	Interval: data.BarTimespanDay,
	Count:    5,
})
if err != nil {
	return err
}
```

Hosts:

| Environment | Host |
|-------------|------|
| Production (Hong Kong) | `https://api.webull.hk` |
| Sandbox (Hong Kong) | `https://api.sandbox.webull.hk` |

Requests use the same signed headers as the rest of the API. Responses are
compact JSON, and numeric fields such as prices are returned as strings to
preserve precision. Timestamps are Unix milliseconds where documented; bar times
come back as strings. `GetBars` issues a one-symbol batch request because the
historical single-symbol endpoint has been retired.

`GetNewsSummary` is the one HTTP Market Data method that replies with a
Server-Sent Events stream rather than a JSON body. It flows through the same
core client pipeline as every other request — the access token, the `x-version`
header (news stays on `v2`), and the configured rate limiter and circuit breaker
are applied, and the request is signed identically. Streamed requests are never
retried, because the long-lived connection cannot be safely replayed. Read
events with `NewsSummaryStream.Next` and close the stream when finished.

## MQTT streaming

Real-time quotes, snapshots, and ticks are delivered over MQTT. The `stream`
package:

- Connects to the region's broker (plain MQTT on port 1883, or
  MQTT-over-WebSocket on port 8883).
- Subscribes and unsubscribes through paired HTTP calls, then receives
  protobuf-encoded messages on the broker.
- Decodes `Quote`, `Snapshot`, and `Tick` payloads into Go types.
- Reconnects and re-subscribes automatically after a dropped session, respecting
  Webull's connection limits.

Hosts:

| Environment | Broker | WebSocket |
|-------------|--------|-----------|
| Production (Hong Kong) | `data-api.webull.hk:1883` | `wss://data-api.webull.hk:8883/mqtt` |
| Sandbox (Hong Kong) | `data-api.sandbox.webull.hk:1883` | `wss://data-api.sandbox.webull.hk:8883/mqtt` |

Streaming is scoped to the symbols you subscribe to. In the sandbox the available
symbol set is limited (currently `AAPL`). See [Streaming](streaming.md) for the
full API and [Troubleshooting](troubleshooting.md) for network and entitlement
limitations.

## Display Solution

Display Solution endpoints use a separate signing and client-token mechanism
and are accessed via `data.Client.DisplayService()`. Their paths follow the
official documentation; the HK sandbox host returns `403 Forbidden` (paid
entitlement required), so live behaviour could not be verified there.

Coverage: screener (`GetDisplayGainersLosers`, `GetDisplayTopActive`), quotes
(`GetDisplaySnapshot`, `GetDisplayBars`, `GetDisplayBarsSingle`, `GetDisplayTick`,
`GetDisplayDepth`), instruments (`GetDSCompanyProfile`, `GetDSAnalystTargetPrice`,
`GetDSAnalystRating`), news (`GetDSNewsSummary`, `GetDSMarketNews`, `GetDSSymbolNews`,
`GetDSLatestNews`), streaming (`DSSubscribe`, `DSUnsubscribe`), instruments v3
(`GetStockProfilesV3`), logos (`GetLogos`), corporate actions (`GetCorporateActions`,
`GetCorporateActionsByMarket`).

## Related

- [Getting Started](getting-started.md) — install and credentials.
- [Authentication](authentication.md) — signing and tokens.
- [Fundamentals](fundamentals.md) — capital flows, industry comparisons, earnings/dividend calendars, SEC filings, and financial statements.
- [Streaming](streaming.md) — real-time pushes and reconnection.
