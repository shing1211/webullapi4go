# Market Data — Futures

Futures market data and instruments.

[<- Webull API Reference](../webull-api.md)

## Futures Tick

`GET /market-data/futures/ticks/list`

> Retrieves futures tick-by-tick trade data.

| | |
|---|---|
| **SDK** | `data.GetFuturesTick` |
| **Reference** | [futures-tick.md](https://developer.webull.hk/apis/docs/reference/futures-tick.md) |
| **Note** | Paths unconfirmed against live API. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Futures symbol. |
| `category` | query | string | yes | Security type. Currently only `US_FUTURES` is supported for this interface. — one of: `US_FUTURES`, `HK_FUTURES` |
| `count` | query | string | yes | Number of ticks, maximum limit 1200 . |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Futures contract symbol used in trading and market data, e.g. front-month code provided by the exchange. |
| `instrument_id` | string | yes | Unique instrument identifier for this futures contract in the Webull system or exchange. |
| `result` | array<object> | yes | List of tick details (trade prints) for this futures contract. |

*Nested — `result`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `time` | string | yes | Trade time of this tick, expressed as Unix epoch timestamp in milliseconds |
| `price` | string | yes | Executed trade price for this futures contract at this tick |
| `volume` | string | yes | Executed trade volume at this tick, expressed in number of futures contracts |
| `side` | string | yes | Such as: B S G L N |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Futures Snapshot

`GET /market-data/futures/snapshots/list`

> Retrieves futures real-time snapshot data.

| | |
|---|---|
| **SDK** | `data.GetFuturesSnapshot` |
| **Reference** | [futures-snapshot.md](https://developer.webull.hk/apis/docs/reference/futures-snapshot.md) |
| **Note** | Paths unconfirmed. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbols` | query | string | yes | List of futures symbols, separated by commas; maximum 20 symbols per query. Example: SILZ5,6BM6. |
| `category` | query | string | yes | Security type. Currently only `US_FUTURES` is supported for this interface. — one of: `US_FUTURES`, `HK_FUTURES` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Futures contract symbol used in trading and market data, e.g. front-month code provided by the exchange. |
| `instrument_id` | string |  | Unique instrument identifier for this futures contract in the Webull system or exchange. |
| `price` | string |  | Last traded price (last done) of the futures contract, quoted in the contract's trading currency (e.g. USD). |
| `open` | string |  | Session open price for this futures contract. Represents the first traded price of the current regular trading session. If no trade has occurred in the session, this field may be empty. |
| `high` | string |  | Session high price for this futures contract during the current regular trading session. If no trade has occurred in the session, this field may be empty. |
| `low` | string |  | Session low price for this futures contract during the current regular trading session. If no trade has occurred in the session, this field may be empty. |
| `pre_close` | string |  | Previous settlement/close price of the futures contract (typically the official settlement price of the previous trading day), quoted in the contract's trading currency. |
| `volume` | string |  | Accumulated traded volume for the current session, expressed in number of futures contracts. If no trading occurred in the session, this field may be empty. |
| `change` | string |  | Absolute price change of the last traded price relative to the previous settlement/close price. If no valid reference price or last trade exists, this field may be empty. |
| `change_ratio` | string |  | Price change ratio of the last traded price relative to the previous settlement/close price, Expressed as a decimal (e.g., -0.0074 represents -0.74%) |
| `last_trade_time` | integer |  | Timestamp of the last executed trade for this futures contract, expressed as Unix epoch time in milliseconds |
| `open_interest` | string |  | Open interest, representing the total number of outstanding and unsettled futures contracts for this instrument, expressed in number of contracts. |
| `quote_time` | integer |  | Quote timestamp of this snapshot, expressed as Unix epoch time in milliseconds (UTC). Represents the time when this snapshot data was generated. |
| `bid` | string |  | Best bid price (top of book), i.e. the highest price currently offered by buyers, quoted in the contract's trading currency. |
| `ask` | string |  | Best ask price (top of book), i.e. the lowest price currently offered by sellers, quoted in the contract's trading currency. |
| `bid_size` | string |  | Best bid size, i.e. the total quantity available at the best bid price, expressed in number of futures contracts (whole contract units). |
| `ask_size` | string |  | Best ask size, i.e. the total quantity available at the best ask price, expressed in number of futures contracts (whole contract units). |
| `settle_date` | string |  | Settlement date of the latest official daily settlement, typically in ISO-8601 datetime format with timezone information. |
| `settle_price` | string |  | Settlement price of the latest official daily settlement, used for marking the contract to market and determining margin requirements, quoted in the contract's trading currency. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Futures Footprint

`GET /market-data/futures/footprints/list`

> Retrieves futures footprint data.

| | |
|---|---|
| **SDK** | `data.GetFuturesFootprint` |
| **Reference** | [futures-footprint.md](https://developer.webull.hk/apis/docs/reference/futures-footprint.md) |
| **Note** | Paths unconfirmed. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbols` | query | string | yes | List of security symbols, supports JSON array format, multiple symbols separated by commas; maximum 20 symbols per query. |
| `category` | query | string | yes | Security type. Category values are as shown in the enum; Only US_FUTURES type queries are supported. — one of: `US_FUTURES`, `HK_FUTURES` |
| `timespan` | query | string | yes | Supports granularities such as S5, S15, M1, M5, and M30. — one of: `S5`, `S15`, `M1`, `M5`, `M30` |
| `count` | query | string |  | Number of bars, default 200, maximum limit 1200. |
| `real_time_required` | query | string | yes | Whether to include the latest unfinalized bar. Default: false. Applies to minute level timespans only |
| `trading_sessions` | query | string |  | Specify trading hours. Only supported RTH — one of: `PRE`, `RTH`, `ATH`, `OVN` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Futures contract symbol used in trading and market data, e.g. front-month code provided by the exchange. |
| `instrument_id` | string |  | Unique instrument identifier for this futures contract in the Webull system or exchange. |
| `result` | array<object> |  | Footprint result |

*Nested — `result`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `time` | string |  | Bar timestamp in ISO-8601 format |
| `trading_session` | string |  | Trading session identifier only supported RTH (e.g.,RTH) |
| `total` | string |  | The sum of the main buy and sell volumes |
| `delta` | string |  | The difference in trading volume (primary buyers - primary sellers) |
| `buy_total` | string |  | Buy-initiated volume |
| `sell_total` | string |  | Sell-initiated volume |
| `buy_detail` | object |  | The main purchase footprint details (quantity combined for items with the same price). |
| `sell_detail` | object |  | The main seller's footprint shows details (quantities combined for the same price). |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Futures Quotes (Depth)

`GET /market-data/futures/depths/list`

> Retrieves futures depth of book data.

| | |
|---|---|
| **SDK** | `data.GetFuturesDepth` |
| **Reference** | [futures-depth-of-book.md](https://developer.webull.hk/apis/docs/reference/futures-depth-of-book.md) |
| **Note** | Paths unconfirmed. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Futures symbol. |
| `category` | query | string | yes | Security type. Currently only `US_FUTURES` is supported for this interface. — one of: `US_FUTURES`, `HK_FUTURES` |
| `depth` | query | string | yes | User-defined number of bid/ask levels (per side) to return in the Level-2 order book. Valid range: 1 – 10. (Note: Level-1 data must be requested separately through the Snapshot API.) |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Futures contract symbol used in trading and market data, e.g. front-month code provided by the exchange. |
| `instrument_id` | string | yes | Unique instrument identifier for this futures contract in the Webull system or exchange. |
| `quote_time` | integer | yes | Quote timestamp of this snapshot, expressed as Unix epoch time in milliseconds. |
| `asks` | array<object> | yes | Array of ask orders |
| `bids` | array<object> | yes | Array of bid orders |

*Nested — `asks`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `price` | string | yes | Ask price for this level in the futures order book. |
| `size` | string | yes | Ask size (quantity) at this price level, expressed in number of contracts. |

*Nested — `bids`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `price` | string | yes | Bid price for this level in the futures order book. |
| `size` | string | yes | Bid size (quantity) at this price level, expressed in number of contracts. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Futures Historical Bars

`GET /market-data/futures/bars/list`

> Retrieves futures historical bars data.

| | |
|---|---|
| **SDK** | `data.GetFuturesBars` |
| **Reference** | [futures-historical-bars.md](https://developer.webull.hk/apis/docs/reference/futures-historical-bars.md) |
| **Note** | Paths unconfirmed. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbols` | query | string | yes | List of futures symbols, separated by commas; maximum 20 symbols per query. Example: SILZ5,6BM6. |
| `category` | query | string | yes | Security type. Currently only `US_FUTURES` is supported for this interface. — one of: `US_FUTURES`, `HK_FUTURES` |
| `timespan` | query | string | yes | Bar time granularity. eg: M1, M5, M15, M30, M60, M120, M240, D, W, M, Y |
| `count` | query | string |  | Number of bars, maximum limit 1200 . |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `result` | array<object> | yes | List of batch bar data results, each element contains historical bar data for one futures. |

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

## Futures Instrument List

`GET /trading/instruments/futures/contracts/list`

> Retrieves detail information of futures instruments.

| | |
|---|---|
| **SDK** | `data.GetFuturesInstruments` |
| **Reference** | [futures-instrument-list.md](https://developer.webull.hk/apis/docs/reference/futures-instrument-list.md) |
| **Note** | Product-codes path confirmed against HK sandbox. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `category` | query | string | yes | Security type. Supported values: US_FUTURES, HK_FUTURES — one of: `US_FUTURES`, `HK_FUTURES` |
| `symbols` | query | string |  | List of futures trading symbols. Accepts JSON array format or comma-separated strings. Maximum of 100 symbols per request. Note: Either symbols or code must be provided. |
| `code` | query | string |  | List of futures trading code, remark:Either 'symbols' or 'code' must be present. |
| `status` | query | string |  | Tradable Status OC - Tradable: Security is available for trading CO - Liquidate only: Security can only be sold, no purchases allowed NT - Non-Tradable: Security cannot be traded Default: OC |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Futures contract symbol used in trading and market data, e.g. front-month or continuous contract code such as ESZ5, ESmain, etc. |
| `instrument_id` | string |  | Unique identifier for this futures instrument in the Webull system or exchange. If the symbol represents a main/continuous contract, this ID is for the main contract itself. For order placement, it needs to be mapped to the actual month contract ID (see `contractId`). |
| `exchange_code` | string |  | Exchange code, for example: CBOE, GLOBEX, XNYM, XCEC, XCME, XCBT, CDE. |
| `code` | string |  | Code for this futures contract, for example: ES. |
| `name` | string |  | Display name of the futures contract. |
| `product_class_id` | integer |  | Futures product class id, For example: 2 |
| `product_class_name` | string |  | Futures product class name, For example: 2 |
| `status` | string |  | Tradable status: OC (Tradable), CO (Liquidate only), NT (Non-Tradable) OC, CO, NT |
| `currency` | string |  | Trading currency of this futures contract, for example: USD. |
| `contract_month` | string |  | Contract delivery month in the format yyyyMM, for example: 202512 means Dec 2025 (year + month). |
| `settlement_date` | string |  | Final settlement (delivery) date of the contract in the format yyyy-MM-dd, for example: 2025-12-29. |
| `size` | string |  | Contract size (multiplier). The notional value of one contract equals futures price multiplied by this size. |
| `unit` | string |  | Contract unit, describing the pricing unit and quantity (for example: index points x USD). — one of: `1 - Index points`, `2 - Hong Kong dollars`, `3 - US dollars`, `4 - Bushels`, `5 - Bushels 2`, `6 - Futures contract`, `7 - Short tons, 2000 pounds`, `8 - Pounds`, `9 - Gallons`, `10 - Metric tons, 2204.6 pounds`, `11 - Brazilian real`, `12 - Troy ounces`, `13 - British pounds`, `14 - Euros`, `15 - Mexican peso`, `16 - Czech koruna`, `17 - Polish zloty`, `18 - Israeli shekel`, `19 - Barrels`, `20 - Metric ton`, `21 - Australian dollar`, `22 - New Zealand dollar`, `23 - Canadian dollar`, `24 - Swiss franc`, `25 - Japanese yen`, `26 - South African rand`, `27 - Hungarian forint`, `28 - Korean won`, `29 - Million British thermal units`, `30 - Chinese renminbi`, `31 - Megawatt hours`, `41 - Megawatt`, `42 - Therms`, `51 - Environmental offset`, `52 - Basis points`, `53 - Metric tons (thousands)`, `54 - Gross tons`, `55 - Tons (thousands)`, `56 - Ton`, `57 - Bitcoin`, `58 - Russian ruble`, `59 - Indian rupee`, `60 - 1 day of time charter`, `61 - Cubic meter`, `62 - Kiloliters`, `63 - Kilos`, `64 - Chilean peso`, `65 - Regional Greenhouse Gas Initiative allowances (RGGI)`, `66 - Hundredweight, 100 pounds`, `67 - Norwegian krone`, `68 - Allowance (emission)`, `69 - Board feet`, `70 - Grams`, `71 - Swedish krona`, `72 - Environmental credit`, `73 - Dry metric tons`, `74 - Shares`, `75 - Metric ton`, `76 - Malaysian ringgit`, `77 - Ether`, `78 - Pounds net weight`, `79 - Renewable Identification Number (RIN)`, `80 - Barrels (thousands)`, `81 - Troy ounce (millions)` |
| `min_tick` | string |  | Minimum price increment (tick size) for the futures price. |
| `first_notice_date` | string |  | First notice date. For physically delivered contracts, this is the first date on which physical delivery can be assigned. After this date, new long positions cannot be opened, and existing long positions are typically forced to close a few trading days before this date. For cash-settled or index futures, this field is usually empty. |
| `last_notice_date` | string |  | Last notice date, i.e. the last date on which the buyer can be notified to take physical delivery. |
| `first_trading_date` | string |  | First trading date on which this futures contract becomes tradable. |
| `last_trading_date` | string |  | Last trading date, i.e. the final trading day in the delivery month. After this date, any outstanding futures positions must be closed out through physical delivery or cash settlement. For cash-settled contracts, trading is allowed normally before the last trading deadline. For non-cash-settled contracts, opening new positions is usually restricted from three trading days before the earlier of the last trading date or first notice date. |
| `contract_type` | string |  | Contract type. MONTHLY means regular month contract; MAIN means main/continuous contract. — one of: `MONTHLY`, `MAIN` |
| `settlement` | string |  | Settlement method of the contract. Cash means cash settlement; Physical means physical delivery. — one of: `Cash`, `Physical` |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Futures Product Codes

`GET /trading/instruments/futures/product-codes/list`

> Retrieves futures product codes list.

| | |
|---|---|
| **SDK** | `data.GetFuturesProductCodes` |
| **Reference** | [futures-products.md](https://developer.webull.hk/apis/docs/reference/futures-products.md) |
| **Note** | Confirmed against HK sandbox. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `category` | query | string | yes | Security type. Supported values: US_FUTURES, HK_FUTURES — one of: `US_FUTURES`, `HK_FUTURES` |
| `product_class_id` | query | integer |  | Product class id. |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string |  | Display name of the futures product, e.g., E-Mini S&P 500 |
| `code` | string |  | Futures product code: often one to three letter codes identifying the asset that is attached to a specific contract. For example: ES. |
| `product_class_id` | integer |  | Futures product class id, For example: 2 |
| `product_class_name` | string |  | Futures product class name, For example: 2 |
| `exchange_code` | string |  | Futures product for exchange code, For example: XCME |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Futures Product Classes

`GET /trading/instruments/futures/product-classes/list`

> Retrieves futures product classes.

| | |
|---|---|
| **SDK** | `data.GetFuturesProductClasses` |
| **Reference** | [futures-products-class.md](https://developer.webull.hk/apis/docs/reference/futures-products-class.md) |
| **Note** | Unconfirmed. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `category` | query | string | yes | Security type. Supported values: US_FUTURES, HK_FUTURES — one of: `US_FUTURES`, `HK_FUTURES` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `product_class_id` | integer |  | Futures product class id |
| `product_class_name` | string |  | Futures product class name |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

