# Market Data — Screener

Ranked lists and sector data. Some rankers require the Display Solution entitlement (marked below).

[<- Webull API Reference](../webull-api.md)

## Top Gainers/Losers

`GET /market-data/screeners/gainers-losers/list`

> • Function description: Top Gainers/Losers. The only difference between Top Gainers and Top Losers is that when order=CHANGE_RATIO, direction is passed as ASC (Losers) and DESC (Gainers). Returns top 200 results without pagination. • Frequency limit: Market-data interfaces rate limit is 600 requests per minute.

| | |
|---|---|
| **SDK** | `data.GetTopGainersLosers` |
| **Reference** | [get-gainers-losers.md](https://developer.webull.hk/apis/docs/reference/get-gainers-losers.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `rank_type` | query | string | yes | Ranking time dimension that determines the calculation period for price change. Default: DAY_1. — one of: `PRE_MARKET`, `AFTER_MARKET`, `MIN_3`, `MIN_5`, `DAY_1`, `DAY_5`, `MONTH_1`, `MONTH_3`, `WEEK_52` |
| `category` | query | string | yes | Security market category. Default: US_STOCK. — one of: `US_STOCK` |
| `sort_by` | query | string | yes | Secondary sort field for further ordering within the ranking. Default: CHANGE_RATIO. — one of: `CHANGE_RATIO`, `RELATIVE_VOLUME_10D`, `MARKET_VALUE`, `CLOSE`, `PRICE`, `PE_TTM`, `HIGH`, `LOW`, `AMPLITUDE`, `TURNOVER`, `VOLUME` |
| `direction` | query | string |  | Sort direction. Default: DESC for Gainers, use ASC for Losers. — one of: `ASC`, `DESC` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string |  | Unique identifier of the tradable instrument |
| `symbol` | string |  | Trading symbol of the financial instrument |
| `name` | string |  | Full name of the instrument |
| `exchange_code` | string |  | Standardized exchange code |
| `currency_code` | string |  | Denomination currency of the instrument (ISO 4217) |
| `pre_close` | string |  | Previous trading day's closing price |
| `open` | string |  | Opening price for the current trading day |
| `high` | string |  | Intraday high price for the current trading day |
| `low` | string |  | Intraday low price for the current trading day |
| `close` | string |  | Latest traded price for the current trading day |
| `price` | string |  | Most recent quoted price within the selected time interval (pre market / post market / intraday) |
| `change` | string |  | Absolute price change within the selected time interval. Returns pre/post market change when in extended hours |
| `change_ratio` | string |  | Price change percentage within the selected time interval (decimal ratio). Returns pre/post market change percentage when in extended hours |
| `volume` | string |  | Cumulative traded volume for the current day (in shares) |
| `turnover` | string |  | Cumulative turnover amount in denomination currency |
| `turnover_rate` | string |  | Turnover rate as a decimal ratio (e.g. 0.05 represents 5%) |
| `market_value` | string |  | Total market capitalization in denomination currency |
| `amplitude` | string |  | Price amplitude ((high - low) / pre_close) as a decimal ratio |
| `relative_volume_10d` | string |  | Relative volume (current day volume / 10 day average volume) |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Top Actives

`GET /market-data/screeners/top-actives/list`

> • Function description: Stock Top Active Rank. Most actively traded stocks ranked by volume, relative volume, turnover, turnover rate, or amplitude. The relative_volume_10d field is unique to the Top Active response compared to the Gainers/Losers endpoint. Returns top 200 results without pagination. • Frequency limit: Market-data interfaces rate limit is 600 requests per minute.

| | |
|---|---|
| **SDK** | `data.GetMostActive` |
| **Reference** | [get-top-active.md](https://developer.webull.hk/apis/docs/reference/get-top-active.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `category` | query | string | yes | Security market category. Default: US_STOCK. — one of: `US_STOCK` |
| `rank_type` | query | string |  | Ranking dimension that determines which activity metric is used for filtering. Default: VOLUME. — one of: `VOLUME`, `RELATIVE_VOLUME_10D`, `TURNOVER`, `TURNOVER_RATE`, `AMPLITUDE` |
| `sort_by` | query | string |  | Secondary sort field for further ordering within the ranking. Default: VOLUME. — one of: `CHANGE_RATIO`, `RELATIVE_VOLUME_10D`, `MARKET_VALUE`, `CLOSE`, `PRICE`, `PE_TTM`, `HIGH`, `LOW`, `AMPLITUDE`, `TURNOVER`, `VOLUME` |
| `direction` | query | string |  | Sort direction. Default: DESC. — one of: `ASC`, `DESC` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string |  | Unique identifier of the tradable instrument |
| `symbol` | string |  | Trading symbol of the financial instrument |
| `name` | string |  | Full name of the instrument |
| `exchange_code` | string |  | Standardized exchange code |
| `currency_code` | string |  | Denomination currency of the instrument (ISO 4217) |
| `pre_close` | string |  | Previous trading day's closing price |
| `open` | string |  | Opening price for the current trading day |
| `high` | string |  | Intraday high price for the current trading day |
| `low` | string |  | Intraday low price for the current trading day |
| `close` | string |  | Latest traded price for the current trading day |
| `price` | string |  | Most recent quoted price within the selected time interval (pre market / post market / intraday) |
| `change` | string |  | Absolute price change within the selected time interval. Returns pre/post market change when in extended hours |
| `change_ratio` | string |  | Price change percentage within the selected time interval (decimal ratio). Returns pre/post market change percentage when in extended hours |
| `volume` | string |  | Cumulative traded volume for the current day (in shares) |
| `turnover` | string |  | Cumulative turnover amount in denomination currency |
| `turnover_rate` | string |  | Turnover rate as a decimal ratio (e.g. 0.05 represents 5%) |
| `market_value` | string |  | Total market capitalization in denomination currency |
| `amplitude` | string |  | Price amplitude ((high - low) / pre_close) as a decimal ratio |
| `relative_volume_10d` | string |  | Relative volume (current day volume / 10 day average volume) |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Market Sectors

`GET /market-data/screeners/market-sectors/list`

> • Function description: Get all sector overview data including sector name, change ratio, volume, market value, and leading stocks. • Frequency limit: Market-data interfaces rate limit is 600 requests per minute.

| | |
|---|---|
| **SDK** | `data.GetMarketSectors` |
| **Reference** | [get-market-sectors.md](https://developer.webull.hk/apis/docs/reference/get-market-sectors.md) |
| **Note** | Display Solution; HK sandbox host blocked (403). |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `category` | query | string | yes | Security market category — one of: `US_STOCK` |
| `agg_type` | query | string |  | Statistics type, default is MARKET_VALUE. — one of: `MARKET_VALUE`, `VOLUME` |
| `period` | query | string |  | Statistics period, default is D1. — one of: `D1`, `D5`, `MO1`, `MO3` |
| `direction` | query | string |  | Sort direction. Default: ASC. — one of: `ASC`, `DESC` |
| `pagination_key` | query | string |  | Pagination key from previous response for next page |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | List of market sectors |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `id` | string |  | Sector ID |
| `name` | string |  | Sector Name |
| `change_ratio` | string |  | Price change ratio relative to previous close. Expressed as a decimal (e.g., 0.0111 = 1.11%) |
| `volume` | string |  | Trading volume within the current statistical period |
| `market_value` | string |  | Market value within the current statistical period |
| `declined` | string |  | Number of stocks that have fallen |
| `advanced` | string |  | Number of stocks that have risen |
| `flat` | string |  | Number of stocks with a price change of 0 |
| `data` | array<object> |  | Leading stocks in the sector |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string |  | Security ID |
| `name` | string |  | Security name |
| `symbol` | string |  | Security symbol |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Market Sector Detail

`GET /market-data/screeners/market-sectors/get`

> • Function description: Get stock list and statistics for a specific sector. • Frequency limit: Market-data interfaces rate limit is 600 requests per minute.

| | |
|---|---|
| **SDK** | `data.GetMarketSectorDetail` |
| **Reference** | [get-market-sectors-detail.md](https://developer.webull.hk/apis/docs/reference/get-market-sectors-detail.md) |
| **Note** | Display Solution. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `sector_id` | query | string | yes | Sector ID |
| `category` | query | string | yes | Security market category — one of: `US_STOCK` |
| `period` | query | string |  | Statistics period, default is D1. — one of: `D1`, `D5`, `MO1`, `MO3` |
| `sort_by` | query | string |  | Sort field, default is CHANGE_RATIO. — one of: `CHANGE_RATIO`, `RELATIVE_VOLUME_10D`, `MARKET_VALUE`, `CLOSE`, `PRICE`, `PE_TTM`, `HIGH`, `LOW`, `AMPLITUDE`, `TURNOVER`, `VOLUME`, `YIELD`, `DIVIDEND` |
| `direction` | query | string |  | Sort direction. Default: ASC. — one of: `ASC`, `DESC` |
| `pagination_key` | query | string |  | Pagination key from previous response for next page |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `id` | string |  | Sector ID |
| `name` | string |  | Sector Name |
| `change_ratio` | string |  | Price change ratio |
| `declined` | string |  | Number of stocks that have fallen |
| `advanced` | string |  | Number of stocks that have risen |
| `flat` | string |  | Number of stocks with a price change of 0 |
| `data` | array<object> |  | List of stocks in the sector |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string |  | Security ID |
| `category` | string |  | Security category |
| `currency` | string |  | Currency code |
| `name` | string |  | Security name |
| `symbol` | string |  | Security symbol |
| `exchange_code` | string |  | Exchange code |
| `close` | string |  | Latest intraday price |
| `change_ratio` | string |  | Price change ratio |
| `price` | string |  | Latest price |
| `volume` | string |  | Trade volume |
| `market_value` | string |  | Market Value |
| `turnover_rate` | string |  | Turnover Rate |
| `amplitude` | string |  | Amplitude Ratio |
| `high` | string |  | Today's high |
| `low` | string |  | Today's low |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## High Dividend Rank

`GET /market-data/screeners/high-dividend-ranks/list`

> • Function description: Get high dividend rank list. Returns top 200 results without pagination. • Frequency limit: Market-data interfaces rate limit is 600 requests per minute.

| | |
|---|---|
| **SDK** | `data.GetHighDividendRank` |
| **Reference** | [get-high-dividend.md](https://developer.webull.hk/apis/docs/reference/get-high-dividend.md) |
| **Note** | Display Solution. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `category` | query | string | yes | Security market category — one of: `US_STOCK` |
| `sort_by` | query | string |  | Sort field, default is YIELD. — one of: `CHANGE_RATIO`, `RELATIVE_VOLUME_10D`, `MARKET_VALUE`, `CLOSE`, `PRICE`, `PE_TTM`, `HIGH`, `LOW`, `AMPLITUDE`, `TURNOVER`, `VOLUME`, `YIELD`, `DIVIDEND` |
| `direction` | query | string |  | Sort direction. Default: DESC. — one of: `ASC`, `DESC` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string |  | Security ID |
| `category` | string |  | Security category |
| `currency` | string |  | Currency code |
| `name` | string |  | Security name |
| `symbol` | string |  | Security symbol |
| `exchange_code` | string |  | Exchange code |
| `close` | string |  | Latest intraday price |
| `change` | string |  | Trade change for the current trading day |
| `change_ratio` | string |  | Price change ratio |
| `price` | string |  | Latest price |
| `volume` | string |  | Trade volume |
| `market_value` | string |  | Market Value |
| `turnover_rate` | string |  | Turnover Rate |
| `amplitude` | string |  | Amplitude Ratio |
| `high` | string |  | Today's high |
| `low` | string |  | Today's low |
| `turnover` | string |  | Turnover amount |
| `yield` | string |  | Dividend Yield |
| `dividend` | string |  | Dividend |
| `ex_date` | string |  | Ex-Date |
| `pe_ttm` | string |  | PE TTM |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## 52-Week High/Low

`GET /market-data/screeners/week52-high-low/list`

> • Function description: Get 52 week high/low rank list. Returns top 200 results without pagination. • Frequency limit: Market-data interfaces rate limit is 600 requests per minute.

| | |
|---|---|
| **SDK** | `data.GetWeek52HighLow` |
| **Reference** | [get-week-52-high-low.md](https://developer.webull.hk/apis/docs/reference/get-week-52-high-low.md) |
| **Note** | Display Solution. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `rank_type` | query | string |  | 52 week rank type. — one of: `NEW_HIGH`, `NEAR_HIGH`, `NEW_LOW`, `NEAR_LOW` |
| `category` | query | string | yes | Security market category — one of: `US_STOCK` |
| `sort_by` | query | string |  | Sort field, default is CHANGE_RATIO_52W. — one of: `CHANGE_RATIO`, `RELATIVE_VOLUME_10D`, `MARKET_VALUE`, `CLOSE`, `PRICE`, `PE_TTM`, `HIGH`, `LOW`, `AMPLITUDE`, `TURNOVER`, `VOLUME`, `YIELD`, `DIVIDEND` |
| `direction` | query | string |  | Sort direction. Default: ASC. — one of: `ASC`, `DESC` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string |  | Security ID |
| `category` | string |  | Security category |
| `currency` | string |  | Currency code |
| `name` | string |  | Security name |
| `symbol` | string |  | Security symbol |
| `exchange_code` | string |  | Exchange code |
| `close` | string |  | Latest intraday price |
| `change` | string |  | Trade change for the current trading day |
| `change_ratio` | string |  | Price change ratio |
| `price` | string |  | Latest price |
| `volume` | string |  | Trade Volume |
| `market_value` | string |  | Market Value |
| `turnover_rate` | string |  | Turnover Rate |
| `high` | string |  | Today's high |
| `low` | string |  | Today's low |
| `turnover` | string |  | Turnover amount |
| `amplitude` | string |  | Amplitude Ratio |
| `price_1w` | string |  | This Week High/Low |
| `price_52w` | string |  | 52W Last High/Low |
| `change_ratio_52w` | string |  | Change ratio since previous highest/lowest |
| `pe_ttm` | string |  | PE TTM |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

