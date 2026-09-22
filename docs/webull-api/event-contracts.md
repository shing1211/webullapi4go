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

## Event Contract Tags

`GET /market-data/instruments/event-contracts/categories/tags/list`

> Retrieves all available category tags for event contract series. Use the returned tags to filter series via the tags parameter in the Series List endpoint.

| | |
|---|---|
| **SDK** | `data.GetEventContractTags` |
| **Reference** | [all-tags-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/all-tags-using-get.md) |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `tags` | array<string> |  | Category Tag list |
| `category_id` | integer |  | Category ID |
| `category_name` | string |  | Category Name |
| `category_code` | string |  | Category Code |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Contract Events List

`GET /market-data/instruments/event-contracts/events/list`

> Retrieves a list of tradable events under a series. Each event represents a specific question or market (e.g., '2026-27 College Football National Championship Winner'). Filter by series_symbol to get events for a specific series.

| | |
|---|---|
| **SDK** | `data.GetEventContractEventsList` |
| **Reference** | [event-list-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-list-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `series_symbol` | query | string |  | series symbol |
| `status` | query | string |  | status, ACTIVE=event is open for trading; INACTIVE=event is closed or settled, default:ACTIVE |
| `pagination_key` | query | string |  | Pagination key returned from previous page response. Pass null or omit for first page. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Data list |
| `pagination_key` | string |  | Pagination key for next page. null means no more data. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Event symbol |
| `name` | string |  | Event name |
| `status` | string |  | status, eg:ACTIVE=event is open for trading；INACTIVE=event is closed or settled, default:ACTIVE |
| `series_id` | integer |  | Series ID |
| `event_id` | integer |  | Event ID |
| `short_name` | string |  | Event short name |
| `strike_date` | string |  | strike date,Only present for sports category events. Represents the settlement date/period of the event. |
| `strike_period` | string |  | strike period, Only present for sports category events. Represents the settlement date/period of the event. |
| `mutually_exclusive` | boolean |  | mutually exclusive, true or false |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Contract Milestones

`GET /market-data/instruments/event-contracts/milestones/list`

> Retrieves a paginated list of milestones (individual game or economic release events). Each milestone contains match details, team info, and related event contract symbols.

| | |
|---|---|
| **SDK** | `data.GetEventContractMilestones` |
| **Reference** | [milestones-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/milestones-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `minimum_start_date` | query | integer |  | Query data with match start time (millisecond timestamp) that is later than that time. |
| `category` | query | string |  | category code |
| `competition` | query | string |  | competition name |
| `related_event_symbol` | query | string |  | related event symbol |
| `pagination_key` | query | string |  | Pagination key returned from previous page response. Pass null or omit for first page. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Data list |
| `pagination_key` | string |  | Pagination key for next page. null means no more data. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `category` | string |  | Category Code |
| `type` | string |  | Milestone type |
| `title` | string |  | Milestone title |
| `details` | object |  | Detail fields vary by type. Refer to SportsGameDetails when type=$SPORTS_GAME, EconomicReleaseDetails when type=ECONOMIC_RELEASE. |
| `status` | string |  | status, eg:NOT_STARTED,INPROGRESS,CLOSED,CANCELLED,POSTPONED,DELAYED,SUSPENDED,UNKNOWN |
| `milestone_id` | string |  | Milestone Unique Identifier |
| `start_date` | string |  | Start date |
| `end_date` | string |  | End date |
| `related_event_symbols` | array<string> |  | related event symbol |
| `primary_event_symbols` | array<string> |  | primary event symbol |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Contract Series List

`GET /market-data/instruments/event-contracts/series/list`

> Retrieves a paginated list of event contract series. A series represents a recurring competition or event category (e.g., CFP National Champion). Use category and tags to filter by sport type.

| | |
|---|---|
| **SDK** | `data.GetEventContractSeriesList` |
| **Reference** | [series-list-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/series-list-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `category` | query | string |  | category code |
| `tags` | query | string |  | tag name |
| `symbols` | query | string |  | series symbols |
| `pagination_key` | query | string |  | Pagination key returned from previous page response. Pass null or omit for first page. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Data list |
| `pagination_key` | string |  | Pagination key for next page. null means no more data. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Series Symbol |
| `name` | string |  | Series Name |
| `category` | string |  | Category Code |
| `frequency` | string | yes | frequency, e.g., HOURLY,DAILY,WEEKLY,MONTHLY,ANNUAL,ONE_OFF,CUSTOM. |
| `tags` | array<string> |  | tags name |
| `series_id` | integer |  | Series ID |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Contract Sports Filters

`GET /market-data/instruments/event-contracts/sports-filters/list`

> Retrieves available filter options for sports event contracts, including sport tags, competitions, and scopes. Use this to populate filter UI or discover available sports categories before querying series or events.

| | |
|---|---|
| **SDK** | `data.GetEventContractSportsFilters` |
| **Reference** | [sports-filter-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/sports-filter-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `tag` | query | string |  | tag name |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `tag` | string |  | Tag name |
| `competitions` | array<object> |  | competition list |
| `scopes` | array<string> |  | List of available scope categories under this tag (e.g., Games, Futures, Awards). |

*Nested — `competitions`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `competition` | string |  | Competition name (e.g., College Football, Pro Football). |
| `scopes` | array<string> |  | scope list |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Game Stats

`GET /market-data/event-contracts/game-stats/get`

> Retrieves detailed play-by-play or drive-by-drive game statistics for a live sporting event. The response is a flat structure (wide table) containing fields for all sport types. Only fields relevant to the queried milestone's sport will be populated; others will be absent.

| | |
|---|---|
| **SDK** | `data.GetEventGameStats` |
| **Reference** | [event-game-stats-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-game-stats-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `milestone_id` | query | string | yes | Milestone ID. A milestone represents a specific game or real-world occurrence tied to events. |
| `category` | query | string |  | Category, default is US_EVENT, currently only US_EVENT is supported. — one of: `US_EVENT` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `milestone_id` | string | yes | Milestone ID |
| `periods` | array<object> |  | Game periods/innings list |

*Nested — `periods`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `period_number` | integer |  | Period/inning number |
| `period_type` | string |  | Period type, e.g. quarter, top, bottom, period, half. Varies by sport. |
| `events` | array<object> |  | Event list for this period |
| `period_name` | string |  | Period name. May appear in baseball, soccer. |
| `attribution` | string |  | Attribution team ID. May appear in baseball. |
| `half` | string |  | Half inning indicator. May appear in baseball. |
| `away_score` | integer |  | Away team current score. May appear in baseball, hockey, soccer. |
| `away_team_id` | string |  | Away team ID. May appear in baseball, hockey, soccer. |
| `home_score` | integer |  | Home team current score. May appear in baseball, hockey, soccer. |
| `home_team_id` | string |  | Home team ID. May appear in baseball, hockey, soccer. |

*Nested — `events`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `type` | string |  | Event type identifier, e.g. football_drive, basketball_play, baseball_play, hockey_play, soccer_event. |
| `description` | string |  | Event description. |
| `attribution` | string |  | Attribution team/player ID. |
| `away_points` | integer |  | Away team points at this event. |
| `home_points` | integer |  | Home team points at this event. |
| `clock` | string |  | Game clock. |
| `possession` | string |  | Possession team ID. |
| `event_type` | string |  | Event sub-type (e.g. three_point_made, field_goal_made, goal, yellow_card). May appear in basketball, soccer. |
| `wall_clock` | integer |  | Event wall clock UTC timestamp in seconds. May appear in basketball. |
| `half` | string |  | Half inning indicator. May appear in baseball. |
| `strength` | string |  | Strength status, e.g. even/powerplay/shorthanded. May appear in hockey. |
| `competitor` | string |  | Competitor side: home/away. May appear in soccer. |
| `match_time` | integer |  | Match time in minutes. May appear in soccer. |
| `player_name` | string |  | Player name. May appear in soccer. |
| `stoppage_time` | integer |  | Stoppage time in minutes. May appear in soccer. |
| `def_points` | integer |  | Defensive points. May appear in football. |
| `end_reason` | string |  | Drive end reason. May appear in football. |
| `gain` | integer |  | Drive gain in yards. May appear in football. |
| `off_points` | integer |  | Offensive points. May appear in football. |
| `play_count` | integer |  | Number of plays in drive. May appear in football. |
| `plays` | array<object> |  | Play details list. May appear in football. |

*Nested — `plays`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `type` | string |  | Play type identifier. |
| `clock` | string |  | Game clock. |
| `description` | string |  | Play description. |
| `down` | integer |  | Current down number. |
| `yfd` | integer |  | Yards to first down. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Live Data

`GET /market-data/event-contracts/live-data/get`

> Retrieves real-time live game/event data for a specified milestone, including scores, game clock, period, and winner information. Use this to display live match status alongside event contract prices.

| | |
|---|---|
| **SDK** | `data.GetEventLiveData` |
| **Reference** | [event-live-data-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-live-data-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `milestone_id` | query | string | yes | Milestone ID. A milestone represents a specific game or real-world occurrence tied to events. |
| `category` | query | string |  | Category, default is US_EVENT, currently only US_EVENT is supported. — one of: `US_EVENT` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `type` | string | yes | Sport type identifier, e.g. basketball_game, baseball_game, baseball_tournament, football_game, hockey_match, hockey_tournament, soccer_tournament_multi_leg, tennis_tournament_singles, golf_tournament, cricket_match. |
| `milestone_id` | string | yes | Milestone ID. |
| `status` | string | yes | Game status. Values: NOT_STARTED (not begun), INPROGRESS (underway), CLOSED (concluded with result), CANCELLED (cancelled, contracts may be voided), POSTPONED (delayed to future date), DELAYED (temporarily paused, expected to resume), SUSPENDED (indefinitely paused, outcome pending). |
| `winner` | string |  | Winner identifier. Empty string or absent if not yet determined. |
| `last_play` | object |  | EventContractLiveLastPlayVo |
| `last_updated_ts` | integer |  | Last updated timestamp in seconds. |
| `details` | object |  | Sport-specific details. Structure varies by type field. Type mapping: • basketball_game → EventContractLiveBasketballDetailsVo • baseball_game → EventContractLiveBaseballDetailsVo • baseball_tournament → EventContractLiveBaseballDetailsVo • football_game → EventContractLiveFootballDetailsVo • hockey_match → EventContractLiveHockeyDetailsVo • hockey_tournament → EventContractLiveHockeyDetailsVo • soccer_tournament_multi_leg → EventContractLiveSoccerDetailsVo • tennis_tournament_singles → EventContractLiveTennisDetailsVo • golf_tournament → EventContractLiveGolfDetailsVo • cricket_match → EventContractLiveCricketDetailsVo |

*Nested — `last_play`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `description` | string |  | Last play description. |
| `occurence_ts` | integer |  | Occurrence timestamp in seconds. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Market Bars

`GET /market-data/event-contracts/markets/bars/list`

> Retrieves historical OHLCV bar data for one or more event contract markets, keyed by market symbol. Use this endpoint to build price charts for individual contract markets. Maximum 100 symbols per request.

| | |
|---|---|
| **SDK** | `data.GetEventMarketBars` |
| **Reference** | [event-market-bars-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-market-bars-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbols` | query | string | yes | Comma-separated market symbols. A market is a single binary contract within an event. Maximum 100 symbols. |
| `category` | query | string |  | Category, default is US_EVENT, currently only US_EVENT is supported. — one of: `US_EVENT` |
| `start_time` | query | integer |  | Start time (unix timestamp in milliseconds). Empty means no lower bound. |
| `end_time` | query | integer |  | End time (unix timestamp in milliseconds). Empty means no upper bound. |
| `count` | query | integer |  | Number of bars. Range: 1-1200, default 200. |
| `timespan` | query | string | yes | Bar time granularity. M1=1min, M5=5min, M15=15min, M30=30min, M60=1hour, M120=2hour, M240=4hour, D=Daily, W=Weekly, M=Monthly, Y=Yearly. |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `end_period_time` | string | yes | Bar end period time (UTC datetime) |
| `volume` | string | yes | Volume during period |
| `open` | string |  | Open price (nullable) |
| `high` | string |  | High price (nullable) |
| `low` | string |  | Low price (nullable) |
| `close` | string |  | Close price (nullable) |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Market Bars By Event

`GET /market-data/event-contracts/markets/bars/list-by-event`

> Retrieves historical OHLCV bar data for all markets under a given event, keyed by market symbol. Use this endpoint to compare price movements across all contract outcomes within a single event.

| | |
|---|---|
| **SDK** | `data.GetEventMarketBarsByEvent` |
| **Reference** | [event-market-bars-by-event-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-market-bars-by-event-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `event_symbol` | query | string | yes | Event unique identifier. An event contains multiple markets (binary contracts). |
| `category` | query | string |  | Category, default is US_EVENT, currently only US_EVENT is supported. — one of: `US_EVENT` |
| `start_time` | query | integer |  | Start time (unix timestamp in milliseconds). Empty means no lower bound. |
| `end_time` | query | integer |  | End time (unix timestamp in milliseconds). Empty means no upper bound. |
| `count` | query | integer |  | Number of bars. Range: 1-1200, default 200. |
| `timespan` | query | string | yes | Bar time granularity. M1=1min, M5=5min, M15=15min, M30=30min, M60=1hour, M120=2hour, M240=4hour, D=Daily, W=Weekly, M=Monthly, Y=Yearly. |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `end_period_time` | string | yes | Bar end period time (UTC datetime) |
| `volume` | string | yes | Volume during period |
| `open` | string |  | Open price (nullable) |
| `high` | string |  | High price (nullable) |
| `low` | string |  | Low price (nullable) |
| `close` | string |  | Close price (nullable) |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Market Depth

`GET /market-data/event-contracts/markets/depths/list`

> Retrieves the order book (bid/ask depth) for a single event contract market. Each level shows the price and aggregate size (quantity of open orders) at that price point.

| | |
|---|---|
| **SDK** | `data.GetEventMarketDepth` |
| **Reference** | [event-market-depth-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-market-depth-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Market symbol. A market is a single binary contract within an event. |
| `category` | query | string |  | Category, default is US_EVENT, currently only US_EVENT is supported. — one of: `US_EVENT` |
| `depth` | query | integer |  | Market depth levels. Range: 0-100. Default 0 returns all available levels. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Contract symbol |
| `instrument_id` | string | yes | Contract instrument ID |
| `yes_asks` | array<object> | yes | Yes side ask orders, sorted by price ascending |
| `yes_bids` | array<object> | yes | Yes side bid orders, sorted by price descending |
| `no_asks` | array<object> | yes | No side ask orders, sorted by price ascending |
| `no_bids` | array<object> | yes | No side bid orders, sorted by price descending |

*Nested — `yes_asks`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `price` | string | yes | Price |
| `size` | string | yes | Size (Quantity) |

*Nested — `yes_bids`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `price` | string | yes | Price |
| `size` | string | yes | Size (Quantity) |

*Nested — `no_asks`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `price` | string | yes | Price |
| `size` | string | yes | Size (Quantity) |

*Nested — `no_bids`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `price` | string | yes | Price |
| `size` | string | yes | Size (Quantity) |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Event Market Snapshot

`GET /market-data/event-contracts/markets/snapshots/list`

> Retrieves the latest market snapshot for a single event contract market, including yes/no bid-ask prices, last trade price, volume, open interest, and market status.

| | |
|---|---|
| **SDK** | `data.GetEventMarketSnapshot` |
| **Reference** | [event-market-snapshot-using-get.md](https://developer.webull.com/apis/docs/reference/broker-market-data-api/event-market-snapshot-using-get.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Market symbol. A market is a single binary contract within an event. |
| `category` | query | string |  | Category, default is US_EVENT, currently only US_EVENT is supported. — one of: `US_EVENT` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Contract symbol |
| `instrument_id` | string | yes | Contract instrument ID |
| `event_symbol` | string |  | Parent event symbol |
| `yes_sub_title` | string |  | Yes side subtitle (e.g. 'Lakers Win') |
| `no_sub_title` | string |  | No side subtitle (e.g. 'Celtics Win') |
| `status` | string |  | Market status. ACTIVE=Market is open for trading; INACTIVE=Market is closed, no new orders accepted. |
| `yes_bid` | string |  | Yes side best bid price |
| `yes_ask` | string |  | Yes side best ask price |
| `no_bid` | string |  | No side best bid price |
| `no_ask` | string |  | No side best ask price |
| `price` | string |  | Price for the last traded YES contract on this market in dollars. |
| `volume` | string |  | String representation of the market volume in contracts. |
| `open_interest` | string |  | String representation of the number of contracts bought on this market disregarding netting. |
| `last_trade_time` | string |  | Timestamp of the most recent trade. Format: ISO 8601 with timezone offset. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

