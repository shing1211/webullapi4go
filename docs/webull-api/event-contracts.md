# Event Contracts

Event-contract instruments and market data (US documentation).

[<- Webull API Reference](../webull-api.md)

## Event Contract Categories

`GET /trading/instruments/event-contracts/categories/list`

> Retrieves all categories under the Event Contract.

| | |
|---|---|
| **SDK** | `data.GetEventContractCategories` |
| **Reference** | [event-categories-list.md](https://developer.webull.com/apis/docs/reference/event-categories-list.md) |
| **Note** | Documented on the US site. |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `category_id` | integer | yes | Category ID. |
| `category_code` | string | yes | Unique identifier code for category. |
| `category_name` | string | yes | Category corresponding name. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Contract Series

`GET /trading/instruments/event-contracts/series/list`

> Retrieves multiple series with specified filters. A series represents a template for recurring events that follow the same format and rules (e.g., "Monthly Jobs Report").

| | |
|---|---|
| **SDK** | `data.GetEventContractSeries` |
| **Reference** | [event-series-list.md](https://developer.webull.com/apis/docs/reference/event-series-list.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `category` | query | string |  | The category which this series belongs to. — one of: `ECONOMICS`, `FINANCIALS`, `POLITICS`, `ENTERTAINMENT`, `SCIENCE_TECHNOLOGY`, `CLIMATE_WEATHER`, `TRANSPORTATION`, `CRYPTO`, `SPORTS` |
| `symbols` | query | string |  | List of series symbols, maximum 100 symbols per query. |
| `pagination_key` | query | string |  | Pagination key from previous response for next page |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Result data list |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `category` | string | yes | The category which this series belongs to. |
| `series_id` | string | yes | ID that identifies this series. |
| `symbol` | string | yes | Symbol that identifies this series. |
| `name` | string | yes | Name that describes the series. |
| `frequency` | string | yes | Description of the frequency of the series. — one of: `HOURLY`, `DAILY`, `WEEKLY`, `MONTHLY`, `ANNUAL`, `ONE_OFF`, `CUSTOM` |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Contract Events

`GET /trading/instruments/event-contracts/events/list`

> Retrieves events under the Event Contract matching the query.

| | |
|---|---|
| **SDK** | `data.GetEventContractEvents` |
| **Reference** | [event-events-list.md](https://developer.webull.com/apis/docs/reference/event-events-list.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `series_symbol` | query | string | yes | Symbol that identifies this series. |
| `symbols` | query | string |  | List of events symbols, maximum 100 symbols per query. |
| `status` | query | string |  | The status of the event. — one of: `ACTIVE`, `INACTIVE` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `series_id` | string | yes | ID identification of series. |
| `symbol` | string | yes | The symbol of the event |
| `name` | string | yes | The name of the event. |
| `status` | string | yes | The status of the event. — one of: `ACTIVE`, `INACTIVE` |
| `short_name` | string | yes | The abbreviation for event. |
| `strike_date` | string |  | Exercise Date. |
| `strike_period` | string |  | Exercise period. |
| `mutually_exclusive` | boolean | yes | Whether mutually exclusive. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Contract Instruments

`GET /trading/instruments/event-contracts/markets/list`

> Retrieves event contract market instruments for the given series symbol.

| | |
|---|---|
| **SDK** | `data.GetEventContractMarkets` |
| **Reference** | [event-market-list.md](https://developer.webull.com/apis/docs/reference/event-market-list.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `series_symbol` | query | string |  | Symbol that identifies this series. |
| `event_symbol` | query | string |  | Symbol of the event events. |
| `symbols` | query | string |  | List of security symbols, maximum 100 symbols per query. |
| `expiration_date_after` | query | string |  | Used to filter items whose expiration date is later than a specified date; the default selection is the current day (inclusive). |
| `pagination_key` | query | string |  | Pagination key from previous response for next page |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Result data list |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `series_id` | string | yes | ID that identifies this series. |
| `series_symbol` | string | yes | Symbol that identifies this series. |
| `series_name` | string | yes | Name that describes the series. |
| `event_symbol` | string | yes | Symbol that identifies this events. |
| `event_name` | string | yes | Name that describes the events. |
| `instrument_id` | string | yes | Unique id of the event market. |
| `symbol` | string | yes | Symbol of the event market. |
| `name` | string | yes | Name of the event market. |
| `yes_condition` | string | yes | Conditions for a 'Yes' outcome. |
| `last_trading_date` | string | yes | Last Notice Day. |
| `status` | string | yes | Listing status. — one of: `NOT_SET`, `LISTING`, `DELISTING`, `OTHER`, `UNRECOGNIZED` |
| `tradable_status` | string | yes | Tradable status: OC (Tradable), CO (Liquidate only), NT (Non-Tradable) — one of: `OC`, `CO`, `NT` |
| `can_close_early` | boolean | yes | Can the contract close early? |
| `expected_exp_date` | string | yes | Expected expiration date of contract. |
| `latest_exp_date` | string | yes | Latest expiration date. |
| `payout_date` | string | yes | Settlement/Payment Date. |
| `fractionable` | boolean | yes | Support fragmented event contracts. |
| `price_ranges` | array<object> |  | Price range. |

*Nested — `price_ranges`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `start` | string | yes | Start price. |
| `end` | string | yes | End price. |
| `step` | string | yes | Step length. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Snapshot

`GET /market-data/event-contracts/snapshots/list`

> Retrieves a real-time snapshot for an event instrument.

| | |
|---|---|
| **SDK** | `data.GetEventSnapshot` |
| **Reference** | [event-snapshot.md](https://developer.webull.com/apis/docs/reference/event-snapshot.md) |
| **Note** | Host/paths unconfirmed. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbols` | query | string | yes | Symbol of the event market, supports JSON array format, multiple symbols separated by commas; maximum 100 symbols per query. |
| `category` | query | string |  | Category, default is US_EVENT, currently only US_EVENT is supported. — one of: `US_EVENT` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string | yes | Unique id of the event market. |
| `symbol` | string | yes | Symbol of the event market. |
| `name` | string | yes | Name of the event market. |
| `price` | string | yes | The current market price for buying or selling an event contract |
| `volume` | string | yes | Number of contracts bought on this event market. |
| `last_trade_time` | integer | yes | Price for the last traded YES contract on this market. |
| `open_interest` | string | yes | Number of contracts bought on this event market disconsidering netting. |
| `yes_bid` | string | yes | Price for the highest YES buy offer on this event market. |
| `yes_bid_size` | string | yes | Size for the highest YES buy offer on this event market. |
| `yes_ask` | string | yes | Price for the lowest YES sell offer on this event market. |
| `yes_ask_size` | string | yes | Size for the lowest YES sell offer on this event market. |
| `no_bid` | string | yes | Price for the highest NO buy offer on this event market. |
| `no_bid_size` | string | yes | Size for the highest NO buy offer on this event market. |
| `no_ask` | string | yes | Price for the lowest NO sell offer on this event market. |
| `no_ask_size` | string | yes | Size for the lowest NO sell offer on this event market. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Depth

`GET /market-data/event-contracts/depths/list`

> Retrieves the order book for an event instrument. Only yes/no bids are returned (in binary markets a yes bid at X equals a no ask at 100-X).

| | |
|---|---|
| **SDK** | `data.GetEventDepth` |
| **Reference** | [event-depth.md](https://developer.webull.com/apis/docs/reference/event-depth.md) |
| **Note** | Host/paths unconfirmed. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Symbol of the event market. |
| `category` | query | string |  | Category, default is US_EVENT, currently only US_EVENT is supported. — one of: `US_EVENT` |
| `depth` | query | string |  | Depth of buying and selling orders, default 10 levels, etc. |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string | yes | Unique id of the event market. |
| `symbol` | string | yes | Symbol of the event market. |
| `quote_time` | integer | yes | Quotation Time. |
| `yes_bids` | array<object> | yes | Yes, buy order array. |
| `yes_asks` | array<object> | yes | Yes, sell order array. |
| `no_bids` | array<object> | yes | No, buy order array. |
| `no_asks` | array<object> | yes | No, sell order array. |

*Nested — `yes_bids`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `price` | string | yes | Price |
| `size` | string | yes | Trading volume. |

*Nested — `yes_asks`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `price` | string | yes | Price |
| `size` | string | yes | Trading volume. |

*Nested — `no_bids`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `price` | string | yes | Price |
| `size` | string | yes | Trading volume. |

*Nested — `no_asks`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `price` | string | yes | Price |
| `size` | string | yes | Trading volume. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Bars

`GET /market-data/event-contracts/bars/list`

> Retrieves the most recent N bars for an event symbol.

| | |
|---|---|
| **SDK** | `data.GetEventBars` |
| **Reference** | [event-bars.md](https://developer.webull.com/apis/docs/reference/event-bars.md) |
| **Note** | Host/paths unconfirmed. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbols` | query | string | yes | Symbol of the event market, supports JSON array format, multiple symbols separated by commas; maximum 100 symbols per query. |
| `category` | query | string |  | Category, default is US_EVENT, currently only US_EVENT is supported. — one of: `US_EVENT` |
| `timespan` | query | string | yes | Bar time granularity. — one of: `M1`, `M5`, `M15`, `M30`, `M60`, `M120`, `M240`, `D` |
| `count` | query | string |  | Number of bars, default 200, maximum limit 1200. |
| `real_time_required` | query | string | yes | Does it include the latest data? For candlesticks that are not yet finalized, the default is false (does not include). Only minute timespan is used. |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string | yes | Unique id of the event market. |
| `symbol` | string | yes | Symbol of the event market. |
| `result` | array<object> | yes | K-line data. |

*Nested — `result`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `open` | string | yes | Opening price. |
| `close` | string | yes | Close price. |
| `high` | string | yes | High price. |
| `low` | string | yes | Low price. |
| `volume` | string | yes | Turnover. |
| `time` | string | yes | UTC time |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Tick

`GET /market-data/event-contracts/ticks/list`

> Retrieves tick-by-tick trades for an event contract, sorted latest first.

| | |
|---|---|
| **SDK** | `data.GetEventTick` |
| **Reference** | [event-tick.md](https://developer.webull.com/apis/docs/reference/event-tick.md) |
| **Note** | Host/paths unconfirmed. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Symbol of the event market. |
| `category` | query | string |  | Category, default is US_EVENT, currently only US_EVENT is supported. — one of: `US_EVENT` |
| `count` | query | string |  | Number of tick, default 30, maximum limit 1200. |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string | yes | Unique id of the event market. |
| `symbol` | string | yes | Symbol of the event market. |
| `result` | array<object> | yes | Tick data. |

*Nested — `result`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `time` | string | yes | Timestamp when this trade was executed. |
| `yes_price` | string | yes | Yes price for this trade in dollars. |
| `no_price` | string | yes | No price for this trade in dollars. |
| `volume` | string | yes | Turnover. |
| `side` | string | yes | Side for the taker of this trade(yes/no). |
| `trade_id` | string | yes | Unique identifier for this trade. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

