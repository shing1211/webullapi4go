# Glossary

SDK-specific terms and Webull OpenAPI concepts used throughout the documentation.

## Authentication and access

| Term | Definition |
|------|------------|
| **App Key** | Your Webull OpenAPI application identifier. Sent as the `x-app-key` header on every request. Obtained from the Webull developer portal. |
| **App Secret** | The HMAC signing key paired with your App Key. Used only to compute `x-signature`; never transmitted over the wire. |
| **Access Token** | A time-limited bearer token (default 15 days) obtained from the token endpoint. Sent as `x-access-token` on authenticated requests. Created with `POST /openapi/auth/token/create`. |
| **HMAC Signing** | Every REST request is signed with HMAC-SHA1 over a canonical string built from query parameters and participating headers. The gRPC events API uses HMAC-SHA256 with a different canonical string. |
| **Token Lifecycle** | The state machine governing an access token: `PENDING` → `NORMAL` → `EXPIRED`. Sandbox tokens skip `PENDING` and go directly to `NORMAL`. |
| **Display Token** | A client-token used by Display Solution endpoints. Managed automatically by `data.Client.DisplayService()`. Separate from the regular access token. |

## Security and regions

| Term | Definition |
|------|------------|
| **Sandbox** | A simulated trading environment at `api.sandbox.webull.hk`. Real orders are placed but with test funds. Selected with `client.WithSandbox()` or `WEBULL_ENVIRONMENT=sandbox`. |
| **Region** | The Webull geographic region (`HK` or `US`). Determines the base URL and available endpoints. Set with `client.WithRegion()`. |
| **BCAN** | Broker Client Account Number. A unique identifier for Hong Kong equity trading accounts, used in `no_party_ids` for HK orders. |

## Market Data

| Term | Definition |
|------|------------|
| **Category** | The security type for market data queries: `US_STOCK`, `US_ETF`, `HK_STOCK`, `CN_STOCK`, `US_OPTION`. Required on snapshot, bars, and quote requests. |
| **Instrument ID** | Webull's internal numeric identifier for a financial instrument. Returned in responses; not the same as a ticker symbol. |
| **Snapshot** | A point-in-time market data snapshot containing price, volume, bid/ask, and other fields. |
| **Bars** | Historical OHLCV candle data. The `BarQuery.Interval` selects the timespan (1m, 5m, 15m, 30m, 1h, 1d, etc.). |
| **Depth** | Level-2 order-book depth showing bid/ask levels at multiple price points. Used with MQTT streaming and Display Solution. |
| **Footprint** | An advanced depth analytics endpoint showing order-flow data. Requires a paid entitlement. |
| **Pagination Key** | A cursor token returned by list endpoints. Pass it as `PaginationKey` to fetch the next page. Absence indicates the last page. |

## Trading

| Term | Definition |
|------|------------|
| **Client Order ID** | A caller-supplied unique identifier (1–32 chars, `[A-Za-z0-9_-]`) for an order. Used to replace, cancel, and query orders. Must be unique per account. |
| **Instrument Type** | The asset class of an order: `EQUITY`, `OPTION`, `FUTURES`, or `EVENT`. |
| **Order Type** | The execution instruction: `LIMIT`, `MARKET`, `STOP_LOSS`, `STOP_LOSS_LIMIT`, `TRAILING_STOP_LOSS`, etc. |
| **Time in Force** | How long an order stays active: `DAY` (expires end of trading day), `GTC` (good till cancelled), `GTD` (good till date). |
| **Entrust Type** | Whether the order specifies quantity (`QTY`) or total cash amount (`AMOUNT`). `AMOUNT` supports US fractional share trading. |
| **Combo Type** | The role an order plays in a combo (multi-order) structure: `NORMAL`, `MASTER`, `STOP_PROFIT`, `STOP_LOSS`, `OTO`, `OCO`, `OTOCO`. |
| **Option Strategy** | The structure of a multi-leg options order: `SINGLE`, `VERTICAL`, `STRADDLE`, `STRANGLE`, `IRON_CONDOR`, etc. |
| **Preview** | A dry-run validation of an order. Returns estimated cost and fees without mutating the account. Always preview before placing. |
| **Party ID** | Hong Kong equity order metadata specifying the BCAN, source (`D`), and role (`3`). Required for all HK equity orders. |

## Streaming and events

| Term | Definition |
|------|------------|
| **MQTT** | The lightweight pub/sub protocol used for real-time Market Data streaming. The SDK supports plain TCP (port 1883) and WebSocket (port 8883). |
| **gRPC** | The RPC framework used for Trading Events. Server-streaming over TLS on port 443. |
| **Subscribe / Unsubscribe** | HTTP calls that tell the server which symbols to push data for. The MQTT connection only carries pushes. |
| **Session ID** | A unique identifier for an MQTT connection. A new connection reusing an existing session ID disconnects the previous one. |
| **Reconnect** | Automatic reconnection after a dropped MQTT or gRPC connection. The SDK re-subscribes to active subscriptions transparently. |
| **Event Order** | A JSON payload delivered over gRPC when an order's status changes (submitted, filled, cancelled, etc.). |
| **Event Position** | A JSON payload delivered over gRPC when an event-contract position is settled. |
| **Event Option** | A JSON payload delivered over gRPC when an option order's status changes. |

## Error handling

| Term | Definition |
|------|------------|
| **Transient Error** | An error caused by temporary conditions (rate limit, network, server error). Safe to retry with backoff. |
| **Permanent Error** | An error caused by invalid input or missing permissions. Retrying will not help. |
| **Sentinel Error** | A predefined error value (e.g., `client.ErrAccessTokenRequired`) matched with `errors.Is`. |
| **Circuit Breaker** | A resilience pattern that stops sending requests after consecutive failures. Configured with `client.WithBreaker`. |

## Related

- [Getting Started](getting-started.md) — install and credentials.
- [Authentication](authentication.md) — signing and token lifecycle.
- [Patterns](patterns.md) — shared SDK conventions.
