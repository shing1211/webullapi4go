# Market Data

Market Data is available over two transports: HTTP for on-demand queries and
MQTT for real-time streaming. Both share the same authentication described in
[Authentication](authentication.md).

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

Coverage in v0.1:

| Group | Methods |
|-------|---------|
| Instruments | `GetStockInstruments` |
| Profile and analyst | `GetCompanyProfile`, `GetAnalystTargetPrice`, `GetAnalystRating` |
| Futures static data | `GetFuturesInstruments`, `GetFuturesProductCodes`, `GetFuturesProductClasses` |
| Snapshot and quotes | `GetSnapshot`, `GetQuotes` |
| Ticks and bars | `GetTick`, `GetBars`, `GetBatchBars` |
| Depth analytics | `GetFootprint`, `GetNOIIBars`, `GetNOIISnapshot` |
| Discovery | `GetTopGainersLosers`, `GetMostActive` |
| Watchlists | `GetWatchlists`, `CreateWatchlist`, `UpdateWatchlist`, `DeleteWatchlist`, `GetWatchlistInstruments`, `AddWatchlistInstruments`, `RemoveWatchlistInstruments`, `UpdateWatchlistInstruments` |
| Derivatives and news | `GetOptionTick`, `GetOptionSnapshot`, `GetOptionBars`, `GetOptionExpirations`*, `GetOptionChain`*, `GetNewsSummary` |

\* Speculative: see the note below.

!!! warning "Unpublished option-discovery endpoints"

    `GetOptionExpirations` and `GetOptionChain` are **speculative**. The Webull
    OpenAPI reference documents only option ticks, snapshots, and historical
    bars, so these two endpoints are not part of the published API. Their paths,
    parameters, and response shapes follow the naming convention of the
    documented option endpoints, are marked `TODO(t10)` in the code, and may
    return empty results. Verify them against a live US account before relying
    on them.

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

## Related

- [Getting Started](getting-started.md) — install and credentials.
- [Authentication](authentication.md) — signing and tokens.
- [Fundamentals](fundamentals.md) — capital flows, industry comparisons, earnings/dividend calendars, SEC filings, and financial statements.
- [Streaming](streaming.md) — real-time pushes and reconnection.
