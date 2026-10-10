# Market Data — Crypto

> ⚠️ **Generated file — do not edit.** Regenerate with `python tools/webull-docgen/docgen.py <target>` (`reference`, `master`, `reconciliation` or `all`).

Crypto market data and instruments (US only).

[<- Webull API Reference](../webull-api.md)

## Crypto Snapshot

`GET /market-data/crypto/snapshots/list`

> Retrieve real-time market snapshot data for one or more crypto symbols. The response includes key market indicators such as latest price, price change, price change percentage, bid/ask quotes, and other real-time metrics. Supports querying up to <b>20 symbols</b> per request. <b>Rate Limits:</b> • 1 request per second per App Key • Market Data Global Limit: 600 requests per minute

| | |
|---|---|
| **SDK** | `data.GetCryptoSnapshot` |
| **Reference** | [crypto-snapshot.md](https://developer.webull.com/apis/docs/reference/crypto-snapshot.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbols` | query | string | yes | List of crypto trading symbols. Supports JSON array format or comma-separated values. Maximum 20 symbols. |
| `category` | query | string | yes | Asset category (e.g., US_CRYPTO). Only crypto categories are supported for this API. — one of: `US_CRYPTO` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string |  | Instrument ID associated with the trading symbol |
| `symbol` | string |  | Trading symbol of the crypto |
| `pre_close` | string |  | Previous closing price |
| `last_trade_time` | integer |  | Timestamp of the most recent trade (Unix timestamp in milliseconds). |
| `price` | string |  | Latest traded price |
| `open` | string |  | Opening price of the current trading day. May be empty if no trades occurred. |
| `high` | string |  | Highest traded price of the current trading day. Empty if no trading occurred. |
| `low` | string |  | Lowest traded price of the current trading day. Empty if no trading occurred. |
| `change` | string |  | Absolute price change compared with the previous close |
| `change_ratio` | string |  | Price change ratio compared with the previous close |
| `quote_time` | string |  | Timestamp of the latest quote update (Unix timestamp in milliseconds). |
| `bid` | string |  | Best bid price (Bid 1) |
| `bid_size` | string |  | The total volume (in lots/shares) of buy orders at the best bid price (buy-1 price) |
| `ask` | string |  | Best ask price (Ask 1) |
| `ask_size` | string |  | The total volume (in lots/shares) of sell orders at the best ask price (sell-1 price). |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Crypto Bars

`GET /market-data/crypto/bars/list`

> Retrieve historical candlestick (K-line) data for a specified crypto symbol. Supports multiple time intervals such as M1, M5, H1, D, etc. • Daily and higher intervals return forward-adjusted bars • Minute intervals return non-adjusted bars Supports retrieving the most recent N bars: • Range: 1–1200 bars (all intervals) <b>Rate Limits:</b> • 1 request per second per App Key • Market Data Global Limit: 600 requests per minute

| | |
|---|---|
| **SDK** | `data.GetCryptoBars` |
| **Reference** | [crypto-bars.md](https://developer.webull.com/apis/docs/reference/crypto-bars.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbols` | query | string | yes | List of crypto trading symbols. Supports JSON array format or comma-separated values. Maximum 20 symbols. |
| `category` | query | string | yes | Asset category (e.g., US_CRYPTO). Only US_CRYPTO are supported. — one of: `US_CRYPTO` |
| `timespan` | query | string | yes | Time interval of the candlesticks (e.g., M1, M5, D). — one of: `M1`, `M5`, `M15`, `M30`, `M60`, `M120`, `M240`, `D`, `W`, `M`, `Y` |
| `count` | query | string |  | Number of bars to return. Default is 200. Range 1–1200; |
| `real_time_required` | query | string | yes | Whether to include the most recent in-progress bar. • true: Only completed historical bars are returned • false: Includes the latest in-progress bar |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Crypto trading symbol |
| `instrument_id` | string |  | Internal instrument identifier |
| `result` | array<object> |  | List of historical K-line (candlestick) records |

*Nested — `result`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `time` | string |  | Bar timestamp in ISO-8601 or Unix time format |
| `open` | string |  | Open price |
| `close` | string |  | Close price |
| `high` | string |  | Highest price within the bar interval |
| `low` | string |  | Lowest price within the bar interval |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Crypto Instruments

`GET /trading/instruments/crypto/profiles/list`

> Retrieves profile information for one or more crypto instruments.

| | |
|---|---|
| **SDK** | `data.GetCryptoInstruments` |
| **Reference** | [crypto-instrument-list.md](https://developer.webull.com/apis/docs/reference/crypto-instrument-list.md) |
| **Note** | Trading API surface. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `category` | query | string | yes | Instrument type — one of: `US_CRYPTO` |
| `symbols` | query | string |  | List of crypto trading symbols, maximum 100 symbols per query. |
| `status` | query | string |  | Tradable status: OC (Tradable), CO (Liquidate only), NT (Non-Tradable) — one of: `OC`, `CO`, `NT` |
| `pagination_key` | query | string |  | Pagination key from previous response for next page |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Result data list |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string |  | Symbol name, e.g. BTCUSD |
| `instrument_id` | string |  | Unique identifier of the security |
| `exchange_code` | string |  | Exchange code, e.g. CCC |
| `category` | string |  | Instrument type, e.g. US_CRYPTO — one of: `US_STOCK` |
| `symbol` | string |  | Symbol of the instrument |
| `status` | string |  | Tradable status: OC (Tradable), CO (Liquidate only), NT (Non-Tradable) — one of: `OC`, `CO`, `NT` |
| `min_trade_amt` | string |  | Minimum trade amount |
| `max_trade_amt` | string |  | Maximum trade amount |
| `min_trade_qty` | string |  | Minimum trade quantity |
| `max_trade_qty` | string |  | Maximum trade quantity |
| `price_step` | string |  | Price step increment |
| `lot_size` | string |  | Lot size |
| `currency` | string |  | currency |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

