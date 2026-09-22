# Display Solution

> ⚠️ **Generated file — do not edit.** Regenerate with `python tools/webull-docgen/docgen.py <target>` (`reference`, `master`, `reconciliation` or `all`).

Hosted Display Solution: a separate entitlement and host with Client-to-Server (Bearer) authentication. The SDK routes these through `display.Service`.

[<- Webull API Reference](../webull-api.md)

## Stock Top Gainers/Losers

`GET /market-data/screeners/gainers-losers/list`

> Retrieves a ranked list of top gaining or losing stocks for a specified time period. To get top gainers, pass direction=DESC; for top losers, pass direction=ASC. The rank_type parameter controls the time window (e.g., D1=today, W52=52-week)

| | |
|---|---|
| **SDK** | `data.GetDisplayGainersLosers` |
| **Reference** | [top-gainers-using-get-new.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/top-gainers-using-get-new.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `rank_type` | query | string |  | Ranking type — one of: `PRE_MARKET`, `AFTER_MARKET`, `M3`, `M5`, `D1`, `D5`, `MO1`, `MO3`, `W52` |
| `category` | query | string | yes | Security category — one of: `US_STOCK` |
| `sort_by` | query | string |  | Sort Field — one of: `CHANGE_RATIO`, `RELATIVE_VOLUME_10D`, `MARKET_VALUE`, `CLOSE`, `PRICE`, `PE_TTM`, `HIGH`, `LOW`, `AMPLITUDE`, `TURNOVER`, `VOLUME` |
| `direction` | query | string |  | Sorting direction. Ascending order: ASC, descending order: DESC — one of: `ASC`, `DESC` |
| `pagination_key` | query | string |  | Pagination key returned from previous page response. Pass null or omit for first page. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Data list |
| `pagination_key` | string |  | Pagination key for next page. null means no more data. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Security symbol |
| `name` | string |  | Security name |
| `exchange_code` | string |  | Exchange code |
| `currency` | string |  | Currency code |
| `pre_close` | string |  | The closing price of the previous trading day |
| `open` | string |  | Open price for the current trading day |
| `high` | string |  | Today’s high |
| `low` | string |  | Today’s low |
| `close` | string |  | Latest intraday prices for the current trading day |
| `price` | string |  | The latest price for the current trading day |
| `change` | string |  | Trade change for the current trading day |
| `change_ratio` | string |  | Price change ratio relative to previous close. Expressed as a decimal (e.g., 0.0111 = 1.11%) |
| `volume` | string |  | Trade volume |
| `turnover` | string |  | Transaction amount, US market stocks and ETFs do not return this data under Nb authorization |
| `turnover_rate` | string |  | Turnover Rate |
| `market_value` | string |  | Market Value |
| `amplitude` | string |  | Amplitude Ratio |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Top Active

`GET /market-data/screeners/top-actives/list`

> Retrieves stock top active rank list

| | |
|---|---|
| **SDK** | `data.GetDisplayTopActive` |
| **Reference** | [top-active-using-get-new.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/top-active-using-get-new.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `rank_type` | query | string |  | Rank list type — one of: `VOLUME`, `RELATIVE_VOLUME_10D`, `TURNOVER`, `TURNOVER_RATE`, `AMPLITUDE` |
| `category` | query | string | yes | Security category — one of: `US_STOCK` |
| `sort_by` | query | string |  | Sort Field — one of: `CHANGE_RATIO`, `RELATIVE_VOLUME_10D`, `MARKET_VALUE`, `CLOSE`, `PRICE`, `PE_TTM`, `HIGH`, `LOW`, `AMPLITUDE`, `TURNOVER`, `VOLUME` |
| `direction` | query | string |  | Sorting direction. Ascending order: ASC, descending order: DESC — one of: `ASC`, `DESC` |
| `pagination_key` | query | string |  | Pagination key returned from previous page response. Pass null or omit for first page. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Data list |
| `pagination_key` | string |  | Pagination key for next page. null means no more data. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Security symbol |
| `name` | string |  | Security name |
| `exchange_code` | string |  | Exchange code |
| `currency` | string |  | Currency code |
| `pre_close` | string |  | The closing price of the previous trading day |
| `open` | string |  | Open price for the current trading day |
| `high` | string |  | Today’s high |
| `low` | string |  | Today’s low |
| `close` | string |  | Latest intraday prices for the current trading day |
| `price` | string |  | The latest price for the current trading day |
| `change` | string |  | Trade change for the current trading day |
| `change_ratio` | string |  | Price change ratio relative to previous close. Expressed as a decimal (e.g., 0.0111 = 1.11%) |
| `volume` | string |  | Trade volume |
| `turnover` | string |  | Transaction amount, US market stocks and ETFs do not return this data under Nb authorization |
| `turnover_rate` | string |  | Turnover Rate |
| `market_value` | string |  | Market Value |
| `amplitude` | string |  | Amplitude Ratio |
| `relative_volume_10d` | string |  | 10 day average trading volume ratio |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Snapshot

`POST /market-data/stocks/snapshots/list`

> Retrieves real-time market snapshot data for a security. Returns key market indicators such as latest price, price change, volume, turnover rate, etc. Supports querying various security types including US stocks, etc., with optional inclusion of pre-market, after-hours, and overnight trading data.

| | |
|---|---|
| **SDK** | `data.GetDisplaySnapshot` |
| **Reference** | [snapshot-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/snapshot-using-get.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `category_symbols` | array<object> | yes | List of security symbols by category; maximum 100 symbols per query. |
| `extend_hour_required` | string |  | Whether to include extend hour trading data. |
| `overnight_required` | string |  | Whether to include overnight trading data. |

*Nested — `category_symbols`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `category` | string | yes | Security type. Category values are as shown in the enum. — one of: `US_STOCK` |
| `symbols` | array<string> | yes | List of security symbols, supports JSON array format. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Security symbol |
| `pre_close` | string |  | Previous close price |
| `change_ratio` | string |  | Price change ratio relative to previous close. Expressed as a decimal (e.g., 0.0111 = 1.11%) |
| `last_trade_time` | integer |  | Last trade time |
| `price` | string |  | Current price |
| `open` | string |  | Open price, for US stocks it's intraday open price, excluding pre/post market data. No return value if no trading occurred on the day |
| `high` | string |  | Today's high price, for US stocks it's intraday high, excluding pre/post market data. No return value if no trading occurred on the day |
| `low` | string |  | Today's low price, for US stocks it's intraday low, excluding pre/post market data. No return value if no trading occurred on the day |
| `volume` | string |  | Volume. No return value if no trading occurred on the day |
| `change` | string |  | Change amount. No return value if no trading occurred on the day |
| `close` | string |  | Intraday close price |
| `ask` | string |  | Ask Price |
| `ask_size` | string |  | Ask Size (Quantity) |
| `bid` | string |  | Bid Price |
| `bid_size` | string |  | Bid Size (Quantity) |
| `extend_hour_last_price` | string |  | Pre/post market latest price |
| `extend_hour_change` | string |  | Pre/post market change amount |
| `extend_hour_change_ratio` | string |  | Pre/post market change ratio |
| `extend_hour_volume` | string |  | Pre/post market volume |
| `extend_hour_last_trade_time` | integer |  | Current pre/post market trade time |
| `extend_hour_high` | string |  | Pre/post market high price |
| `extend_hour_low` | string |  | Pre/post market low price |
| `ovn_price` | string |  | Overnight price |
| `ovn_high` | string |  | Overnight high price |
| `ovn_low` | string |  | Overnight low price |
| `ovn_volume` | string |  | Overnight volume |
| `ovn_change` | string |  | Overnight change amount |
| `ovn_change_ratio` | string |  | Overnight change ratio |
| `ovn_last_trade_time` | integer |  | Overnight trade time |
| `ovn_ask` | string |  | Overnight Ask Price |
| `ovn_ask_size` | string |  | Overnight Ask Size (Quantity) |
| `ovn_bid` | string |  | Overnight Bid Price |
| `ovn_bid_size` | string |  | Overnight Bid Size (Quantity) |
| `pb_ratio` | string |  | Price to Book Ratio |
| `ps_ratio` | string |  | Price to Sales Ratio |
| `pe_ratio` | string |  | Price to Earnings Ratio |
| `market_value` | string |  | Total market value |
| `neg_market_value` | string |  | Free-Float market value |
| `yield` | string |  | Dividend yield |
| `total_shares` | string |  | Total shares outstanding |
| `out_standing_shares` | string |  | Free-Float shares outstanding, |
| `fifty_two_wk_high` | string |  | 52 weeks high |
| `fifty_two_wk_low` | string |  | 52 weeks low |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Historical Bars (Batch)

`POST /market-data/stocks/bars/list`

> Retrieves the recent N bars of data based on stock symbols, time granularity, and type. Supports historical bars of various granularities like M1, M5, etc. Currently, daily bars (D) and above only provide forward-adjusted bars; minute bars provide unadjusted bars.

| | |
|---|---|
| **SDK** | `data.GetDisplayBars` |
| **Reference** | [query-batch-bars-using-post.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/query-batch-bars-using-post.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `category_symbols` | array<object> | yes | List of security symbols by category; maximum 100 symbols per query. |
| `count` | string |  | 1-1200（M1：1-1650） |
| `trading_sessions` | string |  | Specify trading session(s). Multiple sessions separated by ",". |
| `interval` | string |  | Bar time granularity:M1, M5, M15, M30, M60, M120, M240, D, W, M, Y |
| `last_time` | string |  | Last time of pre page. Example: 1763555670 |
| `real_time_required` | boolean |  | Return the latest trading data, default is true;\n false: Pulls only the completed bars from the previous period at the nearest whole hour at the time of request.\n true: The returned data includes the latest market data. |

*Nested — `category_symbols`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `category` | string | yes | Security type. Category values are as shown in the enum. — one of: `US_STOCK` |
| `symbols` | array<string> | yes | List of security symbols, supports JSON array format. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | symbol |
| `result` | array<object> | yes | k-line data |
| `times` | array<object> | yes | exchange trading hours |
| `instrument_id` | string | yes | instrument_id |
| `special_times` | array<object> | yes | Special Trading Hours at the Exchange |

*Nested — `result`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `time` | string | yes | Bar timestamp, Unix timestamp format. |
| `open` | string | yes | Open price |
| `close` | string | yes | Close price |
| `high` | string | yes | High price |
| `low` | string | yes | Low price |
| `volume` | string | yes | Volume |
| `trading_sessions` | string | yes | Trading session, e.g., RTH (Regular Trading Hours), PRE (Pre-market). |

*Nested — `times`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `start` | string | yes | exchange start time |
| `end` | string | yes | exchange end time |
| `trading_session` | string | yes | Trading session, e.g., RTH (Regular Trading Hours), PRE (Pre-market). |

*Nested — `special_times`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `start` | integer | yes | exchange start time |
| `end` | integer | yes | exchange end time |
| `trading_session` | string | yes | Trading session, e.g., RTH (Regular Trading Hours), PRE (Pre-market). |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Historical Bars (Single)

`GET /market-data/stocks/bars/get`

> Retrieves the recent N bars of data based on stock symbol, time granularity, and type. Supports historical bars of various granularities like M1, M5, etc. Currently, daily bars (D) and above only provide forward-adjusted bars; minute bars provide unadjusted bars.

| | |
|---|---|
| **SDK** | `data.GetDisplayBarsSingle` |
| **Reference** | [bars-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/bars-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Category values are as shown in the enum. — one of: `US_STOCK` |
| `interval` | query | string | yes | Bar time granularity:M1, M5, M15, M30, M60, M120, M240, D, W, M, Y — one of: `M1`, `M5`, `M15`, `M30`, `M60`, `M120`, `M240`, `D`, `W`, `M`, `Y` |
| `last_time` | query | string |  | Last time of pre page. |
| `count` | query | string |  | Number of bars, default 200, maximum limit 1200 (M1 supports up to 1650). |
| `real_time_required` | query | string |  | Return the latest trading data, default is true;\n false: Pulls only the completed bars from the previous period at the nearest whole hour at the time of request.\n true: The returned data includes the latest market data. |
| `trading_sessions` | query | string |  | Specify trading hours. Multiple selections are allowed. Separate multiple items with ",". — one of: `PRE`, `RTH`, `ATH`, `OVN` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | symbol |
| `result` | array<object> | yes | k-line data |
| `times` | array<object> | yes | exchange trading hours |
| `instrument_id` | string | yes | instrument_id |
| `special_times` | array<object> | yes | Special Trading Hours at the Exchange |

*Nested — `result`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `time` | string | yes | Bar timestamp, Unix timestamp format. |
| `open` | string | yes | Open price |
| `close` | string | yes | Close price |
| `high` | string | yes | High price |
| `low` | string | yes | Low price |
| `volume` | string | yes | Volume |
| `trading_sessions` | string | yes | Trading session, e.g., RTH (Regular Trading Hours), PRE (Pre-market). |

*Nested — `times`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `start` | string | yes | exchange start time |
| `end` | string | yes | exchange end time |
| `trading_session` | string | yes | Trading session, e.g., RTH (Regular Trading Hours), PRE (Pre-market). |

*Nested — `special_times`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `start` | integer | yes | exchange start time |
| `end` | integer | yes | exchange end time |
| `trading_session` | string | yes | Trading session, e.g., RTH (Regular Trading Hours), PRE (Pre-market). |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Tick

`GET /market-data/stocks/ticks/list`

> Retrieves tick-by-tick trade data for a security. Returns detailed tick trade records within a specified time range for a given security, including trade time, price, volume, direction, and other details. Data is sorted in reverse chronological order (latest first).

| | |
|---|---|
| **SDK** | `data.GetDisplayTick` |
| **Reference** | [tick-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/tick-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Category values are as shown in the enum. — one of: `US_STOCK` |
| `last_time` | query | string |  | Last time of pre page. |
| `count` | query | string | yes | Number of ticks, default 100, maximum limit 1000. |
| `trading_sessions` | query | string | yes | Specify trading hours. Multiple selections are allowed. Separate multiple items with ",". — one of: `PRE`, `RTH`, `ATH`, `OVN` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Security symbol |
| `result` | array<object> | yes | Tick details |

*Nested — `result`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `time` | string | yes | Trade time |
| `price` | string | yes | Price |
| `volume` | string | yes | Volume |
| `side` | string | yes | Trade direction, see trade direction for details |
| `trading_session` | string | yes | Trading session |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Quotes Depth

`GET /market-data/stocks/depths/list`

> Retrieves the latest bid/ask data for a security. Returns bid/ask information for a specified depth, including price, quantity, order details, etc.

| | |
|---|---|
| **SDK** | `data.GetDisplayDepth` |
| **Reference** | [quotes-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/quotes-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Category values are as shown in the enum. — one of: `US_STOCK` |
| `depth` | query | string | yes | Market depth, L2-default 10 levels, etc. |
| `overnight_required` | query | string | yes | Whether to include overnight trading data. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Security symbol |
| `quote_time` | string | yes | Quote time |
| `asks` | array<object> | yes | Array of ask orders |
| `bids` | array<object> | yes | Array of bid orders |

*Nested — `asks`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `price` | string | yes | Price |
| `size` | string | yes | Size (Quantity) |
| `order` | array<object> | yes | Array of order details |

*Nested — `order`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `mpid` | string | yes | Market participant ID |
| `size` | string | yes | Size (Quantity) |

*Nested — `bids`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `price` | string | yes | Price |
| `size` | string | yes | Size (Quantity) |
| `order` | array<object> | yes | Array of order details |

*Nested — `order`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `mpid` | string | yes | Market participant ID |
| `size` | string | yes | Size (Quantity) |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## News Summary

`POST /market-data/news/summaries/get`

> Invokes LLM to generate news summaries for watchlist.

| | |
|---|---|
| **SDK** | `data.GetDSNewsSummary` |
| **Reference** | [watchlist-summary-using-post.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/watchlist-summary-using-post.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `category_symbols` | array<object> |  | List of security symbols by category. |
| `lang` | string |  | Support language, enum: [en]. |

*Nested — `category_symbols`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `category` | string | yes | Security type. Category values are as shown in the enum. — one of: `US_STOCK` |
| `symbols` | array<string> | yes | List of security symbols, supports JSON array format. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `type` | string |  | Message type — one of: `meta`, `text`, `table` |
| `message` | string |  | Message content (for text type) |
| `args` | object |  | Additional arguments (for meta type) |
| `headers` | object |  | Table headers (for table type) |
| `rows` | object |  | Table rows (for table type) |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Market News

`GET /market-data/news/market-news/list`

> Retrieves news from the market within the past 3 days.

| | |
|---|---|
| **SDK** | `data.GetDSMarketNews` |
| **Reference** | [list-news-by-market-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-news-by-market-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `market` | query | string | yes | Region，eg: United States：US，Thailand：TH，default US. |
| `language` | query | string | yes | News language. |
| `last_news_id` | query | integer |  | The ID of the last data item on the previous page,default 0. |
| `page_size` | query | integer |  | Number of data per page, default 15. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `id` | integer |  | News id |
| `title` | string |  | News title |
| `source_name` | string |  | News publish source |
| `news_time` | string |  | News publish time |
| `news_url` | string |  | News link |
| `thumbnail` | string |  | News thumbnail |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Symbol News

`GET /market-data/news/symbol-news/list`

> Get news on stocks within the past 3 days.

| | |
|---|---|
| **SDK** | `data.GetDSSymbolNews` |
| **Reference** | [list-news-by-ticker-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-news-by-ticker-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security category. |
| `language` | query | string | yes | News language. |
| `last_news_id` | query | integer |  | The ID of the last data item on the previous page,default 0. |
| `page_size` | query | integer |  | Number of data per page, default 15. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `id` | integer |  | News id |
| `title` | string |  | News title |
| `source_name` | string |  | News publish source |
| `news_time` | string |  | News publish time |
| `news_url` | string |  | News link |
| `thumbnail` | string |  | News thumbnail |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Latest News

`GET /market-data/news/latest-news/list`

> Retrieves latest news within the past 3 days.

| | |
|---|---|
| **SDK** | `data.GetDSLatestNews` |
| **Reference** | [list-latest-news-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-latest-news-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `language` | query | string |  | News language. |
| `last_news_id` | query | integer |  | The ID of the last data item on the previous page,default 0. |
| `page_size` | query | integer |  | Number of data per page, default 15. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `id` | integer |  | News id |
| `title` | string |  | News title |
| `source_name` | string |  | News publish source |
| `news_time` | string |  | News publish time |
| `news_url` | string |  | News link |
| `thumbnail` | string |  | News thumbnail |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Corporate Actions

`GET /market-data/instruments/stocks/corporate-actions/list`

> Supports the query of the corporate events for stock splits and reverse stock split, including past and upcoming events.

| | |
|---|---|
| **SDK** | `data.GetCorporateActions` |
| **Reference** | [corp-action-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/corp-action-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. |
| `start_date` | query | string |  | Event start date, UTC time. Format: yyyy-MM-dd |
| `end_date` | query | string |  | Event end date, UTC time. Format: yyyy-MM-dd |
| `event_types` | query | string |  | Event type collection. Multiple event_types should be separated by , |
| `pagination_key` | query | string |  | Pagination key returned from previous page response. Pass null or omit for first page. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Data list |
| `pagination_key` | string |  | Pagination key for next page. null means no more data. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | integer |  | Security ID |
| `symbol` | string |  | Security symbol, e.g., AAPL, GOOG. |
| `exchange_code` | string |  | Exchange code, e.g., NAS, OTC. |
| `event_type` | string |  | Event type — one of: `NAME_CHANGE`, `CASH_DIVIDEND`, `STOCK_DIVIDEND`, `REVERSE_SPLIT`, `FORWARD_SPLIT`, `SPIN_OFF`, `UNIT_SPLIT`, `MERGER`, `REDEMPTION` |
| `event_action` | string |  | Event status, e.g., I(Insert, Valid)/U(Update, Valid)/C(Cancellation, invalid)/D(Deletion, invalid). |
| `event_id` | integer |  | Event id |
| `source` | string |  | Event source, e.g., WEBULL_ARTIFICIAL(Webull artificial) |
| `ratio_old` | string |  | Old ratio, Before the change |
| `ratio_new` | string |  | New ratio, After the change, During the REVERSE_SPLIT process, when ratio_old is 10 and ratio_new is 5, it means that 10 shares are combined into 5. |
| `event_date` | string |  | Event date, UTC time, e.g: 2021-12-28 |
| `update_time` | string |  | Update time, UTC time, e.g: 2021-12-28T09:00:09.945+0000 |
| `create_time` | string |  | Create time, UTC time, e.g: 2021-12-28T09:00:09.945+0000 |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Corporate Actions By Market

`GET /market-data/instruments/stocks/corporate-actions/list-by-market`

> Retrieves corporate action events for all securities in a specified market within a date range. Use this endpoint to bulk-fetch events across the entire US market, rather than querying by individual symbol.

| | |
|---|---|
| **SDK** | `data.GetCorporateActionsByMarket` |
| **Reference** | [corp-market-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/corp-market-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `market` | query | string | yes | Currently only `US` (US Stock Market) is supported. |
| `start_date` | query | string |  | Event start date, UTC time. Format: yyyy-MM-dd |
| `end_date` | query | string |  | Event end date, UTC time. Format: yyyy-MM-dd |
| `event_types` | query | string |  | Event type collection. Multiple event_types should be separated by , |
| `pagination_key` | query | string |  | Pagination key returned from previous page response. Pass null or omit for first page. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Data list |
| `pagination_key` | string |  | Pagination key for next page. null means no more data. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | integer |  | Security ID |
| `symbol` | string |  | Security symbol, e.g., AAPL, GOOG. |
| `exchange_code` | string |  | Exchange code, e.g., NAS, OTC. |
| `event_type` | string |  | Event type — one of: `NAME_CHANGE`, `CASH_DIVIDEND`, `STOCK_DIVIDEND`, `REVERSE_SPLIT`, `FORWARD_SPLIT`, `SPIN_OFF`, `UNIT_SPLIT`, `MERGER`, `REDEMPTION` |
| `event_action` | string |  | Event status, e.g., I(Insert, Valid)/U(Update, Valid)/C(Cancellation, invalid)/D(Deletion, invalid). |
| `event_id` | integer |  | Event id |
| `source` | string |  | Event source, e.g., WEBULL_ARTIFICIAL(Webull artificial) |
| `ratio_old` | string |  | Old ratio, Before the change |
| `ratio_new` | string |  | New ratio, After the change, During the REVERSE_SPLIT process, when ratio_old is 10 and ratio_new is 5, it means that 10 shares are combined into 5. |
| `event_date` | string |  | Event date, UTC time, e.g: 2021-12-28 |
| `update_time` | string |  | Update time, UTC time, e.g: 2021-12-28T09:00:09.945+0000 |
| `create_time` | string |  | Create time, UTC time, e.g: 2021-12-28T09:00:09.945+0000 |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Get Instruments

`POST /market-data/instruments/stocks/profiles/list`

> Retrieves security information for one or more instruments.

| | |
|---|---|
| **SDK** | `data.GetStockProfilesV3` |
| **Reference** | [list-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-using-get.md) |

**Request body**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `category` | string | yes | Security type. Category values are as shown in the enum. — one of: `US_STOCK` |
| `symbols` | array<string> | yes | List of security symbols, supports JSON array format. |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string |  | Security name |
| `symbol` | string |  | Security symbol. |
| `category` | string |  | Security category. |
| `exchange_code` | string |  | Exchange code |
| `currency` | string |  | Currency |
| `subtype` | string |  | Security Subtypes — one of: `COMMON_STOCK`, `ETF`, `INDEX`, `PREFERRED_STOCK`, `WARRANT`, `UNITS`, `RIGHT` |
| `is_adr` | string |  | Is ADR, true or false |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Batch Logos

`POST /market-data/fundamentals/logos/list`

> Retrieves logo image URLs for the specified securities. URLs are hosted on Webull's CDN.

| | |
|---|---|
| **SDK** | `data.GetLogos` |
| **Reference** | [batch-logo-using-post.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/batch-logo-using-post.md) |

**Request body**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `category` | string | yes | Security type. Category values are as shown in the enum. — one of: `US_STOCK` |
| `symbols` | array<string> | yes | List of security symbols, supports JSON array format, multiple symbols separated by commas; maximum 100 symbols per query. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Security symbol |
| `category` | string |  | Security type. Category values are as shown in the enum — one of: `US_STOCK` |
| `logo_url` | string |  | Ticker's icon url. Returns null if no logo is available for the symbol |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Company Profile

`GET /market-data/fundamentals/company-profiles/get`

> Retrieves company profile for one instrument.

| | |
|---|---|
| **SDK** | `data.GetDSCompanyProfile` |
| **Reference** | [list-company-profile-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-company-profile-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Category values are as shown in the enum. default is US_STOCK — one of: `US_STOCK` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Security symbol |
| `category` | string |  | Security type |
| `company_name` | string |  | Company name |
| `establish_date` | string |  | Date of incorporation |
| `exhibition_code` | string |  | The exchange or market where the security is listed (e.g., NASDAQ, NYSE) |
| `profile` | string |  | Company profile |
| `employees` | string |  | Number of employees |
| `address` | string |  | Headquarters address |
| `ceo` | string |  | Company CEO |
| `industries` | array<string> |  | Company industries |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Analyst Target Price

`GET /market-data/fundamentals/analysis/target-prices/get`

> Retrieves analyst target price for one instrument.

| | |
|---|---|
| **SDK** | `data.GetDSAnalystTargetPrice` |
| **Reference** | [list-analyst-target-price-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-analyst-target-price-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Category values are as shown in the enum. default is US_STOCK — one of: `US_STOCK` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Security symbol |
| `category` | string |  | Security type |
| `mean` | string |  | Average target price |
| `low` | string |  | Lowest target price |
| `high` | string |  | Highest target price |
| `median` | string |  | Median target price |
| `currency` | string |  | Currency |
| `effective_start_date` | string |  | The date from which the current consensus rating is effective, in ISO 8601 format (UTC). |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Analyst Rating

`GET /market-data/fundamentals/analysis/ratings/get`

> Retrieves analyst rating for one instrument.

| | |
|---|---|
| **SDK** | `data.GetDSAnalystRating` |
| **Reference** | [list-analyst-rating-using-get.md](https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-analyst-rating-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Category values are as shown in the enum. default is US_STOCK — one of: `US_STOCK` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Security symbol |
| `category` | string |  | Security type |
| `number` | string |  | Total number of analysts |
| `under_perform` | string |  | Under perform count |
| `buy` | string |  | Buy count |
| `sell` | string |  | Sell count |
| `strong_buy` | string |  | Strong buy count |
| `hold` | string |  | Hold (neutral) count |
| `effective_start_date` | string |  | The date from which the current consensus rating is effective, in ISO 8601 format (UTC). |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Streaming Subscribe

`POST /market-data/streaming/subscribe`

> Subscribe to real-time market data streaming. This interface allows you to subscribe to various types of market data including quotes, snapshots, and tick data for specified securities.

| | |
|---|---|
| **SDK** | `data.DSSubscribe` |
| **Reference** | [subscribe-using-post.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/subscribe-using-post.md) |
| **Note** | US-site reference. |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `session_id` | string | yes | The session_id used to create the connection, and the connection must be successfully established. |
| `category_symbols` | array<object> | yes | List of security symbols by category; maximum 100 symbols per query. |
| `sub_types` | string | yes | Subscription data type(s), multiple types separated by commas ",", enum, refer to: SubType, e.g.: [SNAPSHOT] — one of: `QUOTE`, `SNAPSHOT`, `TICK` |
| `depth` | string |  | LV2 subscription depth, default 10 levels, US stocks max 50 levels. |
| `overnight_required` | boolean |  | Whether to include overnight session, true/false. For US stock subscriptions, includes overnight session, only effective for US stocks, default is not included. |

*Nested — `category_symbols`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `category` | string | yes | Security type. Category values are as shown in the enum. — one of: `US_STOCK` |
| `symbols` | array<string> | yes | List of security symbols, supports JSON array format. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Streaming Unsubscribe

`POST /market-data/streaming/unsubscribe`

> After successfully establishing the market data streaming MQTT connection, call this interface to unsubscribe from real-time market data push. Successful call returns no value; failures return an Error. Unsubscribing will release the topic quota. Frequency limit: 1 call per second per App Key.

| | |
|---|---|
| **SDK** | `data.DSUnsubscribe` |
| **Reference** | [unsubscribe-using-post.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/unsubscribe-using-post.md) |
| **Note** | US-site reference. |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `session_id` | string |  | The session_id used to create the connection, and the connection must be successfully established. |
| `category_symbols` | array<object> | yes | List of security symbols by category; maximum 100 symbols per query. |
| `sub_types` | string | yes | Subscription data type(s), multiple types separated by commas ",", enum, refer to: SubType, e.g.: [SNAPSHOT] — one of: `QUOTE`, `SNAPSHOT`, `TICK` |
| `unsubscribe_all` | boolean |  | Whether to unsubscribe all, true/false. When set to true, all subscriptions will be cancelled. |

*Nested — `category_symbols`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `category` | string | yes | Security type. Category values are as shown in the enum. — one of: `US_STOCK` |
| `symbols` | array<string> | yes | List of security symbols, supports JSON array format. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

