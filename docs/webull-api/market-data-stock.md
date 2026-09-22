# Market Data — Stock

HTTP on-demand stock/ETF market data (Non-Display Solution).

[<- Webull API Reference](../webull-api.md)

## Stock Snapshot

`GET /market-data/stocks/snapshots/list`

> Retrieves stock real-time snapshot data.

| | |
|---|---|
| **SDK** | `data.GetSnapshot` |
| **Reference** | [snapshot.md](https://developer.webull.hk/apis/docs/reference/snapshot.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbols` | query | string | yes | List of security symbols, supports JSON array format, multiple symbols separated by commas; maximum 100 symbols per query. |
| `category` | query | string | yes | Security type. Category values are as shown in the enum; US_OPTION type query is currently not supported. — one of: `US_STOCK`, `US_ETF`, `HK_STOCK`, `CN_STOCK` |
| `extend_hour_required` | query | string |  | Whether to include pre-market and after-hours trading data. |
| `overnight_required` | query | string |  | Whether to include overnight trading data. |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string |  | Instrument ID |
| `pre_close` | string |  | Previous close price |
| `change_ratio` | string |  | Change ratio |
| `symbol` | string | yes | Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives). |
| `last_trade_time` | integer |  | Last trade time |
| `price` | string |  | Current price |
| `open` | string |  | Open price, for US stocks it's intraday open price, excluding pre/post market data. No return value if no trading occurred on the day |
| `close` | string |  | Intraday close price |
| `high` | string |  | Today's high price, for US stocks it's intraday high, excluding pre/post market data. No return value if no trading occurred on the day |
| `low` | string |  | Today's low price, for US stocks it's intraday low, excluding pre/post market data. No return value if no trading occurred on the day |
| `volume` | string |  | Volume. No return value if no trading occurred on the day |
| `change` | string |  | Change amount. No return value if no trading occurred on the day |
| `ask` | string |  | Ask |
| `ask_size` | string |  | Ask size (Quantity) |
| `bid` | string |  | Bid |
| `bid_size` | string |  | Bid size (Quantity) |
| `turnover` | string |  | Turnover rate. |
| `eps` | string |  | Earnings Per Share. |
| `eps_ttm` | string |  | Earnings Per Share (TTM). |
| `lot_size` | string |  | Shares per lot. |
| `bps` | string |  | Book Value Per Share. |
| `extend_hour_last_price` | string |  | Pre/post market latest price |
| `extend_hour_high` | string |  | Pre/post market high price |
| `extend_hour_low` | string |  | Pre/post market low price |
| `extend_hour_change` | string |  | Pre/post market change amount |
| `extend_hour_change_ratio` | string |  | Pre/post market change ratio |
| `extend_hour_volume` | string |  | Pre/post market volume |
| `extend_hour_last_trade_time` | integer |  | Current pre/post market trade time |
| `ovn_price` | string |  | Overnight price |
| `ovn_high` | string |  | Overnight high price |
| `ovn_low` | string |  | Overnight low price |
| `ovn_volume` | string |  | Overnight volume |
| `ovn_change` | string |  | Overnight change amount |
| `ovn_change_ratio` | string |  | Overnight change ratio |
| `ovn_last_trade_time` | integer |  | Overnight trade time |
| `ovn_ask` | string |  | Overnight ask |
| `ovn_ask_size` | string |  | Overnight ask size (Quantity) |
| `ovn_bid` | string |  | Overnight bid |
| `ovn_bid_size` | string |  | Overnight bid size (Quantity) |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Stock Tick

`GET /market-data/stocks/ticks/list`

> Retrieves stock tick-by-tick trade data.

| | |
|---|---|
| **SDK** | `data.GetTick` |
| **Reference** | [tick.md](https://developer.webull.hk/apis/docs/reference/tick.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Category values are as shown in the enum; US_OPTION type query is currently not supported. — one of: `US_STOCK`, `US_ETF`, `HK_STOCK`, `CN_STOCK` |
| `count` | query | string | yes | Number of ticks, default 30, maximum limit 1000. |
| `trading_sessions` | query | string | yes | Specify trading hours. Multiple selections are allowed. Separate multiple items with ",". — one of: `PRE`, `RTH`, `ATH`, `OVN` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Security symbol |
| `instrument_id` | string | yes | Instrument ID |
| `result` | array<object> | yes | Tick details |

*Nested — `result`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `time` | string | yes | Trade time of this tick, expressed as Unix epoch timestamp in milliseconds |
| `price` | string | yes | Executed trade price for this futures contract at this tick |
| `volume` | string | yes | Executed trade volume at this tick, expressed in number of futures contracts |
| `side` | string | yes | Such as: B S G L N |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Stock Quotes (Depth)

`GET /market-data/stocks/depths/list`

> Retrieves stock quotes data.

| | |
|---|---|
| **SDK** | `data.GetQuotes` |
| **Reference** | [quotes.md](https://developer.webull.hk/apis/docs/reference/quotes.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Category values are as shown in the enum; US_OPTION type query is currently not supported. — one of: `US_STOCK`, `US_ETF`, `HK_STOCK`, `CN_STOCK` |
| `depth` | query | string | yes | Market depth, L1-1 level, L2-default 10 levels, etc. |
| `overnight_required` | query | string | yes | Whether to include overnight trading data. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Security symbol |
| `instrument_id` | string | yes | Instrument ID |
| `quote_time` | string | yes | Quote time |
| `asks` | array<object> | yes | Array of ask orders |
| `bids` | array<object> | yes | Array of bid orders |

*Nested — `asks`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `price` | string | yes | Price |
| `size` | string | yes | Size (Quantity) |
| `order` | array<object> | yes | Array of order details |
| `broker` | array<object> |  |  |

*Nested — `order`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `mpid` | string | yes | Market participant ID |
| `size` | string | yes | Size (Quantity) |

*Nested — `broker`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `bid` | string | yes | Broker ID |
| `name` | string | yes | Broker Name |

*Nested — `bids`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `price` | string | yes | Price |
| `size` | string | yes | Size (Quantity) |
| `order` | array<object> | yes | Array of order details |
| `broker` | array<object> |  |  |

*Nested — `order`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `mpid` | string | yes | Market participant ID |
| `size` | string | yes | Size (Quantity) |

*Nested — `broker`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `bid` | string | yes | Broker ID |
| `name` | string | yes | Broker Name |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Stock Historical Bars (Batch)

`POST /market-data/stocks/bars/list`

> Retrieves historical bars data for multiple stock symbols in batch.

| | |
|---|---|
| **SDK** | `data.GetBatchBars` |
| **Reference** | [historical-bars.md](https://developer.webull.hk/apis/docs/reference/historical-bars.md) |
| **Note** | `data.GetBars` POSTs a one-symbol batch against this same endpoint. |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbols` | array<string> | yes | List of security symbols, supports JSON array format, multiple symbols separated by commas; maximum 100 symbols per query. |
| `category` | string | yes | Security type. Category values are as shown in the enum; US_OPTION type query is currently not supported. — one of: `US_STOCK`, `US_ETF`, `HK_STOCK`, `CN_STOCK` |
| `timespan` | string | yes | Bar time granularity. — one of: `S5`, `S15`, `M1`, `M5`, `M15`, `M30`, `M60`, `M120`, `M240`, `D`, `W`, `M`, `Y` |
| `count` | integer |  | Number of bars, default 200, maximum limit 1200 (M1 supports up to 1650). |
| `real_time_required` | boolean |  | Return the latest trading data, default is true; true: Pulls only the completed bars from the previous period at the nearest whole hour at the time of request. false: The returned data includes the latest market data. |
| `trading_sessions` | string |  | Specify trading session(s). Multiple sessions separated by ",". — one of: `PRE`, `RTH`, `ATH`, `OVN` |
| `start_time` | integer |  | Start time as timestamp in milliseconds. Used to specify the beginning of the time range for bar data. |
| `end_time` | integer |  | End time as timestamp in milliseconds. Used to specify the end of the time range for bar data. Delayed permission will automatically offset the time. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `result` | array<object> | yes | List of batch bar data results, each element contains historical bar data for one stock. |

*Nested — `result`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Futures contract symbol used in trading and market data, e.g. front-month code provided by the exchange. |
| `instrument_id` | string | yes | Unique instrument identifier for this futures contract in the Webull system or exchange. |
| `result` | array<object> | yes | List of historical bar data for this futures. |

*Nested — `result`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `time` | string | yes | Bar UTC time |
| `open` | string | yes | Open price |
| `close` | string | yes | Close price |
| `high` | string | yes | High price |
| `low` | string | yes | Low price |
| `volume` | string | yes | Volume |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Stock Footprint

`GET /market-data/stocks/footprints/list`

> Retrieves stock footprint data.

| | |
|---|---|
| **SDK** | `data.GetFootprint` |
| **Reference** | [footprint.md](https://developer.webull.hk/apis/docs/reference/footprint.md) |
| **Note** | Requires a paid entitlement. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbols` | query | string | yes | List of security symbols, supports JSON array format, multiple symbols separated by commas; maximum 20 symbols per query. |
| `category` | query | string | yes | Security type. Category values are as shown in the enum; Only US_STOCK type queries are supported. — one of: `US_STOCK` |
| `timespan` | query | string | yes | Supports granularities such as S5, S15, M1, M5, and M30. — one of: `S5`, `S15`, `M1`, `M5`, `M30` |
| `count` | query | string |  | Number of bars, default 200, maximum limit 1200. |
| `real_time_required` | query | string | yes | Does it include the latest data? For candlesticks that are not yet finalized, the default is false (does not include). Only minute timespan is used. |
| `trading_sessions` | query | string |  | Specify trading hours. OVN type not supported. — one of: `PRE`, `RTH`, `ATH`, `OVN` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Security symbol |
| `instrument_id` | string |  | Unique Identifier for Securities |
| `result` | array<object> |  | Footprint chart candlestick chart |

*Nested — `result`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `time` | string |  | Transaction date |
| `trading_session` | string |  | Trading Hours |
| `total` | string |  | The sum of the main buy and sell volumes |
| `delta` | string |  | The difference in trading volume (primary buyers - primary sellers) |
| `buy_total` | string |  | Buy-initiated volume |
| `sell_total` | string |  | Sell-initiated volume |
| `buy_detail` | object |  | The main purchase footprint details (quantity combined for items with the same price). |
| `sell_detail` | object |  | The main seller's footprint shows details (quantities combined for the same price). |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## NOII Bars

`GET /market-data/stocks/noii-bars/list`

> Retrieves NOII bars data.

| | |
|---|---|
| **SDK** | `data.GetNOIIBars` |
| **Reference** | [get-noii-bars.md](https://developer.webull.hk/apis/docs/reference/get-noii-bars.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. Currently only supports single symbol query. |
| `category` | query | string | yes | Security type. Currently only supports US_STOCK. — one of: `US_STOCK` |
| `imbalance_action_type` | query | string | yes | Imbalance action type: PRE_OPEN (opening imbalance), PRE_CLOSE (closing imbalance). — one of: `PRE_OPEN`, `PRE_CLOSE` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string |  | Instrument unique identifier |
| `symbol` | string |  | Security symbol |
| `imbalance_time` | integer |  | Timestamp of the imbalance data in milliseconds (data publish time) |
| `imbalance_ref_price` | string |  | Reference price |
| `imbalance_near_price` | string |  | Indicative Match Price - the most likely execution price |
| `imbalance_far_price` | string |  | Far Price - the price at which orders could execute in extreme scenarios |
| `imbalance_action_type` | string |  | Imbalance action type: PRE_OPEN (opening imbalance), PRE_CLOSE (closing imbalance) |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## NOII Snapshot

`GET /market-data/stocks/noii-snapshots/list`

> Retrieves NOII snapshot data.

| | |
|---|---|
| **SDK** | `data.GetNOIISnapshot` |
| **Reference** | [get-noii-snapshot.md](https://developer.webull.hk/apis/docs/reference/get-noii-snapshot.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. Currently only supports single symbol query. |
| `category` | query | string | yes | Security type. Currently only supports US_STOCK. — one of: `US_STOCK` |
| `imbalance_action_type` | query | string | yes | Imbalance action type: PRE_OPEN (opening imbalance), PRE_CLOSE (closing imbalance). — one of: `PRE_OPEN`, `PRE_CLOSE` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string |  | Instrument unique identifier |
| `symbol` | string |  | Security symbol |
| `paired_shares` | string |  | Paired shares - the number of shares that can be matched under current conditions |
| `imbalance_shares` | string |  | Imbalance shares - the number of unmatched buy/sell shares |
| `imbalance_side` | string |  | Imbalance side (direction of imbalance) |
| `imbalance_ref_price` | string |  | Reference price |
| `imbalance_near_price` | string |  | Indicative Match Price - the most likely execution price |
| `imbalance_far_price` | string |  | Far Price - the price at which orders could execute in extreme scenarios |
| `imbalance_action_type` | string |  | Imbalance action type: PRE_OPEN (opening imbalance), PRE_CLOSE (closing imbalance) |
| `imbalance_time` | integer |  | Timestamp in milliseconds |
| `imbalance_var_indicator` | string |  | Volatility/imbalance status indicator |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

