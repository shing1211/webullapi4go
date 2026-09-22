# Market Data — Option

> ⚠️ **Generated file — do not edit.** Regenerate with `python tools/webull-docgen/docgen.py <target>` (`reference`, `master`, `reconciliation` or `all`).

Option tick, snapshot and historical bars (Non-Display Solution), plus the option contract list used to build an option chain.

[<- Webull API Reference](../webull-api.md)

## Option Tick

`GET /market-data/options/ticks/list`

> Retrieves option tick-by-tick trade data.

| | |
|---|---|
| **SDK** | `data.GetOptionTick` |
| **Reference** | [option-tick.md](https://developer.webull.hk/apis/docs/reference/option-tick.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Option symbol. |
| `category` | query | string | yes | Security type. Currently only `US_OPTION` is supported for this interface. — one of: `US_OPTION` |
| `count` | query | string | yes | Number of ticks, maximum limit 1200 . |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Option contract symbol used in trading and market data, e.g. front-month code provided by the exchange. |
| `instrument_id` | string | yes | Unique instrument identifier for this option contract in the Webull system or exchange. |
| `result` | array<object> | yes | List of tick details (trade prints) for this option contract. |

*Nested — `result`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `time` | string | yes | Trade time of this tick, expressed as Unix epoch timestamp in milliseconds |
| `price` | string | yes | Executed trade price for this futures contract at this tick |
| `volume` | string | yes | Executed trade volume at this tick, expressed in number of futures contracts |
| `side` | string | yes | Such as: B S G L N |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Option Snapshot

`GET /market-data/options/snapshots/list`

> Retrieves option real-time snapshot data.

| | |
|---|---|
| **SDK** | `data.GetOptionSnapshot` |
| **Reference** | [option-snapshot.md](https://developer.webull.hk/apis/docs/reference/option-snapshot.md) |
| **Note** | At most 20 symbols per query. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbols` | query | string | yes | List of option symbols, separated by commas; maximum 20 symbols per query. Example: AAPL260522C00300000. |
| `category` | query | string | yes | Security type. Currently only `US_OPTION` is supported for this interface. — one of: `US_OPTION` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string |  | Unique instrument identifier for this option contract in the Webull system or exchange. |
| `symbol` | string | yes | Option contract symbol used in trading and market data, e.g. front-month code provided by the exchange. |
| `price` | string |  | Last traded price (last done) of the option contract, quoted in the contract's trading currency (e.g. USD). |
| `open` | string |  | Session open price for this option contract. Represents the first traded price of the current regular trading session. If no trade has occurred in the session, this field may be empty. |
| `high` | string |  | Session high price for this option contract during the current regular trading session. If no trade has occurred in the session, this field may be empty. |
| `low` | string |  | Session low price for this option contract during the current regular trading session. If no trade has occurred in the session, this field may be empty. |
| `pre_close` | string |  | Previous settlement/close price of the option contract (typically the official settlement price of the previous trading day), quoted in the contract's trading currency. |
| `volume` | string |  | Accumulated traded volume for the current session, expressed in number of option contracts. If no trading occurred in the session, this field may be empty. |
| `change` | string |  | Absolute price change of the last traded price relative to the previous settlement/close price. If no valid reference price or last trade exists, this field may be empty. |
| `change_ratio` | string |  | Price change ratio of the last traded price relative to the previous settlement/close price, Expressed as a decimal (e.g., -0.0074 represents -0.74%) |
| `last_trade_time` | integer |  | Timestamp of the last executed trade for this option contract, expressed as Unix epoch time in milliseconds |
| `close` | string |  | Close price of the option contract, quoted in the contract's trading currency. |
| `strike_price` | string |  | Strike price of the option contract, the price at which the option can be exercised. |
| `gamma` | string |  | Gamma, the rate of change of delta with respect to the underlying asset's price. Measures the convexity of the option's value. |
| `delta` | string |  | Delta, the rate of change of the option price with respect to the underlying asset's price. Ranges from -1 to 1 for options. |
| `rho` | string |  | Rho, the rate of change of the option price with respect to the risk-free interest rate. |
| `theta` | string |  | Theta, the rate of change of the option price with respect to time (time decay). Usually expressed as the change in option price per day. |
| `vega` | string |  | Vega, the rate of change of the option price with respect to volatility. Measures sensitivity to implied volatility changes. |
| `imp_vol` | string |  | Implied volatility, the market's forecast of the underlying asset's likely movement. Expressed as a decimal (e.g., 0.609 represents 60.9%). |
| `open_interest` | string |  | Open interest, representing the total number of outstanding and unsettled option contracts for this instrument, expressed in number of contracts. |
| `quote_time` | integer |  | Quote timestamp of this snapshot, expressed as Unix epoch time in milliseconds (UTC). Represents the time when this snapshot data was generated. |
| `bid` | string |  | Best bid price (top of book), i.e. the highest price currently offered by buyers, quoted in the contract's trading currency. |
| `ask` | string |  | Best ask price (top of book), i.e. the lowest price currently offered by sellers, quoted in the contract's trading currency. |
| `ask_size` | string |  | Best ask size, i.e. the total quantity available at the best ask price, expressed in number of option contracts (whole contract units). |
| `bid_size` | string |  | Best bid size, i.e. the total quantity available at the best bid price, expressed in number of option contracts (whole contract units). |
| `deal_amount` | string |  | Total deal amount (traded value) for the option contract, quoted in the contract's trading currency. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Option Historical Bars

`GET /market-data/options/bars/list`

> Retrieves option historical bars data.

| | |
|---|---|
| **SDK** | `data.GetOptionBars` |
| **Reference** | [option-historical-bars.md](https://developer.webull.hk/apis/docs/reference/option-historical-bars.md) |
| **Note** | At most 20 symbols per query. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbols` | query | string | yes | List of option symbols, separated by commas; maximum 20 symbols per query. Example: AAPL260522C00300000 |
| `category` | query | string | yes | Security type. Currently only `US_OPTION` is supported for this interface. — one of: `US_OPTION` |
| `timespan` | query | string | yes | Bar time granularity. eg: M1, M5, M15, M30, M60, M120, M240, D, W, M, Y — one of: `M1`, `M5`, `M15`, `M30`, `M60`, `M120`, `M240`, `D`, `W`, `M`, `Y` |
| `count` | query | string |  | Number of bars, maximum limit 1200 . |
| `real_time_required` | query | string |  | Include the latest data |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `result` | array<object> | yes | List of batch bar data results, each element contains historical bar data for one option. |

*Nested — `result`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Option contract symbol used in trading and market data, e.g. front-month code provided by the exchange. |
| `instrument_id` | string | yes | Unique instrument identifier for this option contract in the Webull system or exchange. |
| `result` | array<object> | yes | List of historical bar data for this option. |

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

## Option Contract List (Chain)

`GET /trading/instruments/options/contracts/list`

> Retrieves option contracts filtered by underlying symbol, status and other attributes. Contains static contract information.

| | |
|---|---|
| **SDK** | `data.GetOptionContracts` |
| **Reference** | [option-contract-list.md](https://developer.webull.hk/apis/docs/reference/option-contract-list.md) |
| **Note** | Trading API surface. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `category` | query | string | yes | Option category. Currently only US_OPTION is supported. — one of: `US_OPTION` |
| `option_symbols` | query | string |  | Option symbols, multiple separated by commas. |
| `underlying_symbols` | query | string |  | Underlying symbols, multiple separated by commas. |
| `status` | query | string |  | Contract status, default LISTING. Enum: LISTING, DELISTING. |
| `start_date` | query | string |  | Exact expiration date, format: YYYY-MM-DD. |
| `end_date` | query | string |  | Expiration date lower bound (inclusive), format: YYYY-MM-DD. |
| `root_symbol` | query | string |  | Root symbol filter (series symbol, e.g. SPXW). Mainly used for index options and post-CA non-standard contracts. |
| `option_type` | query | string |  | Contract type: CALL / PUT. |
| `style` | query | string |  | Exercise style: AMERICAN / EUROPEAN. |
| `strike_price_gte` | query | string |  | Strike price lower bound (inclusive). |
| `strike_price_lte` | query | string |  | Strike price upper bound (inclusive). |
| `ppind` | query | boolean |  | Penny Program Indicator: true = Penny Pilot contract, false = non-Penny Pilot. |
| `show_deliverables` | query | boolean |  | Whether to return deliverables array in response: TRUE / FALSE, default FALSE. |
| `pagination_key` | query | string |  | Pagination key from previous response for next page. Not required for first request. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Result data list |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_id` | string |  | Contract unique identifier |
| `symbol` | string |  | OCC contract symbol |
| `status` | string |  | Contract status: LISTING, DELISTING |
| `tradable_status` | string |  | Trading restriction: OC (Tradable), CO (Liquidate only), NT (Non-Tradable) |
| `expiration_date` | string |  | Expiration date (effective expiration date after CA events), format: YYYY-MM-DD |
| `root_symbol` | string |  | Root symbol (series symbol). For most equity options, same as underlying_symbol; may differ for index and post-CA contracts (e.g., SPXW) |
| `underlying_symbol` | string |  | Underlying symbol |
| `underlying_instrument_id` | string |  | Underlying instrument id |
| `underlying_type` | string |  | Underlying type: EQUITY_PUT_OPTION / EQUITY_CALL_OPTION / INDEX_CALL_OPTION / ETF_CALL_OPTION / ETF_PUT_OPTION |
| `option_type` | string |  | Contract type: CALL / PUT |
| `style` | string |  | Exercise style: AMERICAN / EUROPEAN |
| `strike_price` | string |  | Strike price |
| `multiplier` | string |  | Contract multiplier, typically 100 for US equity options |
| `settlement_method` | string |  | Settlement method: PHYSICAL / CASH |
| `expired_cycle` | string |  | Expiration cycle: DAILY / WEEKLY / MONTHLY / QUARTERLY / EOM |
| `ppind` | boolean |  | Penny Program Indicator: true = Penny Pilot, false = non-Penny Pilot |
| `currency` | string |  | Pricing currency |
| `def_type` | string |  | Definition type: STANDARD / BINARY / FLEX |
| `listed_exchanges` | array<string> |  | Listed exchanges |
| `deliverables` | array<object> |  | Deliverables configuration; after CA a single contract may correspond to multiple underlyings + cash |

*Nested — `deliverables`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `asset_type` | string |  | Asset type. settlement_method=CAFX/CFR/CADF returns CASH; settlement_method=BTOB/CCC/POST/PHYS returns EQUITY |
| `symbol` | string |  | Deliverable underlying symbol |
| `instrument_id` | string |  | Deliverable underlying instrument id |
| `amount` | string |  | Delivery amount: number of shares when type=EQUITY, cash amount when type=CASH |
| `allocation_percentage` | string |  | Allocation percentage (0-100) |
| `settlement_type` | string |  | Settlement cycle: T_0 / T_1 / T_2 / T_3 / T_4 |
| `settlement_method` | string |  | Clearing method: PHYSICAL (BTOB) / CASH_DIFF (CADF) / CASH_FIXED (CAFX) / CCC / CASH_FIXED_RETURN (CFR) / POSITIONAL (POST) |
| `settlement_status` | string |  | Settlement status: DELAYED / REGULAR |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

