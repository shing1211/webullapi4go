# Trading API

> ⚠️ **Generated file — do not edit.** Regenerate with `python tools/webull-docgen/docgen.py <target>` (`reference`, `master`, `reconciliation` or `all`).

Accounts, assets, the order lifecycle, order queries and instruments. Requests require an access token and default to API version `v3`.

[<- Webull API Reference](../webull-api.md)

## Get Instruments

`GET /trading/instruments/stocks/profiles/list`

> Retrieves profile information for one or more stock instruments.

| | |
|---|---|
| **SDK** | `data.GetStockInstruments` |
| **Reference** | [instrument-list.md](https://developer.webull.hk/apis/docs/reference/instrument-list.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `category` | query | string | yes | Security type. — one of: `US_STOCK`, `HK_STOCK`, `CN_STOCK` |
| `symbols` | query | string |  | List of security symbols, maximum 100 symbols per query. |
| `status` | query | string |  | Tradable status: OC (Tradable), CO (Liquidate only), NT (Non-Tradable) — one of: `OC`, `CO`, `NT` |
| `sub_category` | query | string |  | Sub-category of the instrument. Only effective when symbols is not specified. When category = US_STOCK, supported values: COMMON_STOCK, ETF, PREFERRED_STOCK, WARRANT, UNITS, RIGHT. When category = HK_STOCK, CN_STOCK, supported values: COMMON_STOCK, ETF. If not specified, returns all sub-categories. — one of: `COMMON_STOCK`, `ETF`, `PREFERRED_STOCK`, `WARRANT`, `UNITS`, `RIGHT` |
| `pagination_key` | query | string |  | Pagination key from previous response for next page |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Result data list |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string |  | Symbol name, e.g. Apple |
| `instrument_id` | string |  | Unique identifier of the security |
| `exchange_code` | string |  | Exchange code, e.g. CCC |
| `category` | string |  | Instrument type, e.g. US_STOCK — one of: `US_STOCK`, `HK_STOCK`, `CN_STOCK` |
| `symbol` | string |  | Symbol of the instrument |
| `status` | string |  | Tradable status: OC (Tradable), CO (Liquidate only), NT (Non-Tradable) — one of: `OC`, `CO`, `NT` |
| `shortable` | boolean |  | Instrument is shortable or not |
| `fractionable` | boolean |  | Instrument is fractionable or not |
| `marginable` | boolean |  | Instrument is marginable or not |
| `overnight_trading_supported` | boolean |  | Instrument support overnight trading or not |
| `margin_requirement_long` | string |  | Margin requirement ratio for long position |
| `margin_requirement_short` | string |  | Margin requirement ratio for short position |
| `intraday_margin_long` | string |  | Intraday margin requirement ratio for long position |
| `intraday_margin_short` | string |  | Intraday margin requirement ratio for short position |
| `maintenance_margin_long` | string |  | Maintenance margin requirement ratio for long position |
| `maintenance_margin_short` | string |  | Maintenance margin requirement ratio for short position |
| `easy_to_borrow` | boolean |  | Instrument is easy to borrow or not |
| `lot_size` | string |  | Lot size |
| `currency` | string |  | currency |
| `sub_category` | string |  | Sub-category of the instrument. — one of: `COMMON_STOCK`, `ETF`, `PREFERRED_STOCK`, `WARRANT`, `UNITS`, `RIGHT` |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Account List

`GET /trading/accounts/list`

> Retrieves the account list and returns account information.

| | |
|---|---|
| **SDK** | `trade.ListAccounts` |
| **Reference** | [account-list.md](https://developer.webull.hk/apis/docs/reference/account-list.md) |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `account_id` | string |  | Account identifier |
| `account_number` | string |  | Brokerage account |
| `account_type` | string |  | Account type — one of: `MARGIN`, `CASH` |
| `account_class` | string |  | Account Class — one of: `INDIVIDUAL_CASH`, `INDIVIDUAL_MRGN`, `FUTURES_MRGN`, `INSTITUTIONAL_CASH`, `INSTITUTIONAL_MRGN`, `INSTITUTIONAL_FUTURES_MRGN` |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Account Balance

`GET /trading/assets/balances/get`

> Retrieves account details by account ID.

| | |
|---|---|
| **SDK** | `trade.GetBalance` |
| **Reference** | [query-account-balance.md](https://developer.webull.hk/apis/docs/reference/query-account-balance.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `account_id` | query | String | yes | Account identifier |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `total_asset_currency` | string | yes | Currency — one of: `CNH`, `HKD`, `USD` |
| `total_cash_balance` | string | yes | Cash Balance |
| `total_market_value` | string | yes | Total holding market value |
| `total_unrealized_profit_loss` | string | yes | Open P&L |
| `init_margin` | string |  | Initial margin |
| `account_currency_assets` | array<object> | yes | Currency assets Details |

*Nested — `account_currency_assets`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `currency` | string | yes | Currency — one of: `CNH`, `HKD`, `USD` |
| `cash_balance` | string | yes | Cash Balance |
| `settled_cash` | string | yes | Settled Cash |
| `unsettled_cash` | string | yes | Unsettled Cash |
| `market_value` | string | yes | holding market value |
| `held_amount` | string |  | In-transit funds |
| `frozen_amount` | string |  | Frozen funds |
| `buying_power` | string | yes | Buying Power |
| `unrealized_profit_loss` | string | yes | Open P&L |
| `available_withdrawal` | string | yes | The withdrawable amount |
| `interests_unpaid` | string | yes | Interest to be paid |
| `init_margin` | string |  | Init margin |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Account Positions

`GET /trading/assets/positions/list`

> Retrieves positions according to the account ID

| | |
|---|---|
| **SDK** | `trade.GetPositions` |
| **Reference** | [query-account-position.md](https://developer.webull.hk/apis/docs/reference/query-account-position.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `account_id` | query | String | yes | Account identifier |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `position_id` | string | yes | Position ID |
| `currency` | string | yes | Currency — one of: `CNH`, `HKD`, `USD` |
| `quantity` | string | yes | Quantity of the order. Specifies the number of shares or units to transact. For US stocks, fractional quantities are allowed and can include decimals. |
| `symbol` | string | yes | Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives). |
| `option_strategy` | string | yes | Type of options strategy - SINGLE: Indicates a single-leg options order — one of: `SINGLE` |
| `instrument_type` | string | yes | Type of financial instrument associated with the request. — one of: `EQUITY`, `OPTION`, `FUTURES` |
| `last_price` | string | yes | Last Price |
| `cost_price` | string | yes | Cost Basis |
| `unrealized_profit_loss` | string | yes | Open P&L |
| `legs` | array<object> |  | legs |

*Nested — `legs`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives). |
| `quantity` | string |  | Quantity of the order. Specifies the number of shares or units to transact. |
| `option_type` | string |  | Type of the option. - CALL: Right to buy the underlying asset. - PUT: Right to sell the underlying asset. — one of: `CALL`, `PUT` |
| `option_expire_date` | string |  | Option expiration date. Format: yyyy-MM-dd |
| `option_exercise_price` | string |  | Exercise Price |
| `option_contract_multiplier` | string |  | The number of shares corresponding to each option contract |
| `option_contract_deliverable` | string |  | The number of shares required to exercise each contract |
| `expiration_type` | string |  | Option expiration types |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Order Preview

`POST /trading/orders/preview`

> Calculates the estimated amount and cost based on the provided information. Supports simple orders.

| | |
|---|---|
| **SDK** | `trade.PreviewOrder` |
| **Reference** | [common-order-preview.md](https://developer.webull.hk/apis/docs/reference/common-order-preview.md) |
| **Note** | Validates and enforces order guardrails. |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `account_id` | string | yes | Account identifier |
| `client_combo_order_id` | string |  | Unique client-defined identifier for the combined order If combo_type = NORMAL and client_combo_order_id not need to set If combo_type != NORMAL and client_combo_order_id not provided the server will automatically generate one. To sell and close an existing position with take-profit/stop-loss, submit only STOP_PROFIT/STOP_LOSS sub-orders (side = SELL) grouped under the same client_combo_order_id; no MASTER order is required in this scenario. |
| `new_orders` | array<object> | yes | Order Details |

*Nested — `new_orders`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `combo_type` | string | yes | Specifies the type of order combination. For details, please refer to Combo Order. Futures Trading currently support only the NORMAL type. - NORMAL: A standard single order. - MASTER: A primary order that triggers a take-profit or stop-loss order upon execution - STOP_PROFIT: A take-profit order - STOP_LOSS: A stop-loss order - OTO: An order that triggers another order upon execution (One-Triggers-the-Other) - OCO: A pair of orders where the execution of one cancels the other (One-Cancels-the-Other) - OTOCO: An order that triggers an OCO order set upon execution (One-Triggers-One-Cancels-the-Other) Note: When placing take-profit/stop-loss orders to sell and close an existing position, submit only STOP_PROFIT/STOP_LOSS sub-orders (side = SELL) under the same client_combo_order_id; no MASTER order is required or supported in this scenario, since no new position is being opened. Note: OTO, OCO and OTOCO combo types are only supported for stock (EQUITY) orders; option orders (including option_strategy = SINGLE) do not support OTO, OCO or OTOCO. |
| `client_order_id` | string | yes | Unique client-defined identifier for the order. Maximum length is 32 characters and must be unique per account. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). Used to track or reference the order when interacting with the system. |
| `instrument_type` | string | yes | Type of financial instrument associated with the request. — one of: `EQUITY`, `OPTION`, `FUTURES` |
| `market` | string | yes | Market code indicating the trading venue or regulatory region of the financial instrument.Used together with symbol and instrument_type to uniquely identify a tradable instrument. — one of: `US`, `HK`, `CN` |
| `symbol` | string | yes | Trading symbol of the financial instrument. Represents the unique identifier of the security in the specified market. |
| `order_type` | string | yes | Specifies the type of order to be placed. Determines how the order will be executed in the market. Available order types depend on the market and instrument type. Options trading Only LIMIT,STOP_LOSS,STOP_LOSS_LIMIT are supported. U.S. Stock - <b>LIMIT:</b> Limit Order - <b>MARKET:</b> Market Order - <b>STOP_LOSS:</b> Stop Order - <b>STOP_LOSS_LIMIT:</b> Stop Limit Order - <b>MARKET_ON_OPEN:</b> Opening market order - <b>MARKET_ON_CLOSE:</b> Closing market order - <b>TOUCH_MKT:</b> Touch Market Order - <b>TOUCH_LMT:</b> Touch Limit Order - <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order - <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order Hong Kong Stock - <b>ENHANCED_LIMIT:</b> Enhanced Limit Order - <b>AT_AUCTION:</b> At-auction order - <b>AT_AUCTION_LIMIT:</b> At-auction limit order - <b>STOP_LOSS:</b> Stop Order - <b>STOP_LOSS_LIMIT:</b> Stop Limit Order - <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order - <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order - <b>TOUCH_MKT:</b> Touch Market Order - <b>TOUCH_LMT:</b> Touch Limit Order - <b>ODD_LOT_LIMIT:</b> Odd Lot Limit Order China Connect - <b>LIMIT:</b> Limit Order — one of: `LIMIT`, `MARKET`, `STOP_LOSS`, `STOP_LOSS_LIMIT`, `ENHANCED_LIMIT`, `AT_AUCTION`, `AT_AUCTION_LIMIT`, `MARKET_ON_OPEN`, `TRAILING_STOP_LOSS`, `TRAILING_STOP_LOSS_LIMIT`, `TOUCH_MKT`, `TOUCH_LMT`, `ODD_LOT_LIMIT` |
| `entrust_type` | string | yes | Specifies the method for placing the order. - QTY: Order specified by quantity of shares or units. - AMOUNT: Order specified by total cash amount, applicable for fractional share trading of US stocks. — one of: `QTY`, `AMOUNT` |
| `support_trading_session` | string | yes | Specifies the trading session for the order. Applicable to U.S. stock market orders only. Deprecated values: - Y: [Deprecated]Include extended trading hours. - N: [Deprecated]Only support regular trading hours. Active values: - NIGHT: Only supports night trading. - ALL: Include extended trading hours. - CORE: Only support regular trading hours. - ALL_DAY: Included Overnight Hours, 8:00 p.m.ET - 8:00 p.m.ET(the next day) — one of: `Y`, `N`, `NIGHT`, `ALL`, `CORE`, `ALL_DAY` |
| `time_in_force` | string | yes | Specifies the duration for which the order remains active in the market (Time-In-Force). - DAY: The order is valid only for the current trading day and expires at the end of the day. - GTD: order that will automatically expire and be cancelled at a specific future date and time,Currently only supports the US market. - GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days). — one of: `DAY`, `GTD`, `GTC` |
| `side` | string | yes | The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash). — one of: `BUY`, `SELL`, `SHORT` |
| `quantity` | string |  | Transaction quantity. You can specify decimals when placing fractional lot orders for US stocks. |
| `limit_price` | string |  | Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT. Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed. |
| `stop_price` | string |  | Stop price of the order. Required when order_type is STOP_LOSS or STOP_LOSS_LIMIT. Specifies the trigger price at which the stop order becomes active. |
| `option_strategy` | string |  | Type of options strategy - SINGLE: Indicates a single-leg options order — one of: `SINGLE` |
| `legs` | array<object> |  | Option leg detail. Only required when previewing option orders. |

*Nested — `legs`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_type` | string | yes | Type of financial instrument associated with the request. — one of: `EQUITY`, `OPTION`, `FUTURES` |
| `market` | string | yes | Market code indicating the trading venue or regulatory region of the financial instrument.Used together with symbol and instrument_type to uniquely identify a tradable instrument. — one of: `US`, `HK`, `CN` |
| `symbol` | string | yes | Trading symbol of the financial instrument. Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives). |
| `side` | string | yes | The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash). — one of: `BUY`, `SELL`, `SHORT` |
| `strike_price` | string |  | Exercise price (strike price) of the option. Specifies the price at which the underlying asset can be bought (CALL) or sold (PUT) upon exercise. |
| `option_expire_date` | string |  | Expiration date. Format: yyyy-MM-dd |
| `option_type` | string |  | Type of the option. - CALL: Right to buy the underlying asset. - PUT: Right to sell the underlying asset. — one of: `CALL`, `PUT` |
| `quantity` | string |  | Quantity of the order or strategy leg. For stock legs, specifies the number of shares to transact For option legs, specifies the number of option contracts to transact for this leg and is expressed in whole contracts. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `estimated_cost` | string | yes | Estimated capital required for the order. The meaning varies by product type: The actual fee may differ based on final execution. |
| `estimated_transaction_fee` | string | yes | Estimated transaction fee for placing the order, including exchange, clearing, and commission fees. The actual fee may differ based on final execution. |
| `estimated_transaction_fee_detail` | object |  | Breakdown of the estimated transaction fee, including commission and itemized fees. |

*Nested — `estimated_transaction_fee_detail`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `commission` | object |  | Commission breakdown |
| `fees` | array<object> |  | Fee breakdown |

*Nested — `commission`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `actual_commission` | string |  | Actual commission collected |
| `receivable_commission` | string |  | Receivable commission |

*Nested — `fees`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `type` | string |  | Fee type |
| `actual_value` | string |  | Actual fee collected |
| `receivable_value` | string |  | Receivable fee |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Order Place

`POST /trading/orders/place`

> Places equity and options orders. The A-Share trading function is disabled by default.

| | |
|---|---|
| **SDK** | `trade.PlaceOrder` |
| **Reference** | [common-order-place.md](https://developer.webull.hk/apis/docs/reference/common-order-place.md) |
| **Note** | Creates live orders. |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `account_id` | string | yes | Account identifier |
| `client_combo_order_id` | string |  | Unique client-defined identifier for the combined order If combo_type = NORMAL and client_combo_order_id not need to set If combo_type != NORMAL and client_combo_order_id not provided the server will automatically generate one. To sell and close an existing position with take-profit/stop-loss, submit only STOP_PROFIT/STOP_LOSS sub-orders (side = SELL) grouped under the same client_combo_order_id; no MASTER order is required in this scenario. |
| `new_orders` | array<object> | yes | Order Details |

*Nested — `new_orders`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `combo_type` | string | yes | Specifies the type of order combination. For details, please refer to Combo Order. Futures Trading currently support only the NORMAL type. - NORMAL: A standard single order. - MASTER: A primary order that triggers a take-profit or stop-loss order upon execution - STOP_PROFIT: A take-profit order - STOP_LOSS: A stop-loss order - OTO: An order that triggers another order upon execution (One-Triggers-the-Other) - OCO: A pair of orders where the execution of one cancels the other (One-Cancels-the-Other) - OTOCO: An order that triggers an OCO order set upon execution (One-Triggers-One-Cancels-the-Other) Note: When placing take-profit/stop-loss orders to sell and close an existing position, submit only STOP_PROFIT/STOP_LOSS sub-orders (side = SELL) under the same client_combo_order_id; no MASTER order is required or supported in this scenario, since no new position is being opened. Note: OTO, OCO and OTOCO combo types are only supported for stock (EQUITY) orders; option orders (including option_strategy = SINGLE) do not support OTO, OCO or OTOCO. Sub-order quantity limits by combo_type: \| Scenario \| combo_type \| Supported order_type \| Quantity \| Description \| \|---\|---\|---\|---\|---\| \| Take-Profit/Stop-Loss \| MASTER \| MARKET, LIMIT \| 1 \| Master Order \| \| Take-Profit/Stop-Loss \| STOP_PROFIT \| LIMIT \| 0-1 \| Take Profit Order \| \| Take-Profit/Stop-Loss \| STOP_LOSS \| STOP_LOSS \| 0-1 \| Stop Loss Order \| \| OTO \| MASTER \| MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 1 \| Master Order \| \| OTO \| OTO \| MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 1-6 \| Triggered Order(s) \| \| OCO \| OCO \| LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 2-6 \| Mutually Cancelling Orders \| \| OTOCO \| MASTER \| MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 1 \| Master Order \| \| OTOCO \| OTOCO \| LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 1-6 \| OCO Order Set Triggered by MASTER \| |
| `client_order_id` | string | yes | Unique client-defined identifier for the order. Maximum length is 32 characters and must be unique per account. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). Used to track or reference the order when interacting with the system. |
| `instrument_type` | string | yes | Type of financial instrument associated with the request. — one of: `EQUITY`, `OPTION`, `FUTURES` |
| `market` | string | yes | Market code indicating the trading venue or regulatory region of the financial instrument.Used together with symbol and instrument_type to uniquely identify a tradable instrument. — one of: `US`, `HK`, `CN` |
| `symbol` | string | yes | Trading symbol of the financial instrument. Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives). |
| `order_type` | string | yes | Specifies the type of order to be placed. Determines how the order will be executed in the market. Available order types depend on the market and instrument type. Options trading Only LIMIT,STOP_LOSS,STOP_LOSS_LIMIT are supported. U.S. Stock - <b>LIMIT:</b> Limit Order - <b>MARKET:</b> Market Order - <b>STOP_LOSS:</b> Stop Order - <b>STOP_LOSS_LIMIT:</b> Stop Limit Order - <b>MARKET_ON_OPEN:</b> Opening market order - <b>MARKET_ON_CLOSE:</b> Closing market order - <b>TOUCH_MKT:</b> Touch Market Order - <b>TOUCH_LMT:</b> Touch Limit Order - <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order - <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order Hong Kong Stock - <b>ENHANCED_LIMIT:</b> Enhanced Limit Order - <b>AT_AUCTION:</b> At-auction order - <b>AT_AUCTION_LIMIT:</b> At-auction limit order - <b>STOP_LOSS:</b> Stop Order - <b>STOP_LOSS_LIMIT:</b> Stop Limit Order - <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order - <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order - <b>TOUCH_MKT:</b> Touch Market Order - <b>TOUCH_LMT:</b> Touch Limit Order - <b>ODD_LOT_LIMIT:</b> Odd Lot Limit Order China Connect - <b>LIMIT:</b> Limit Order — one of: `LIMIT`, `MARKET`, `STOP_LOSS`, `STOP_LOSS_LIMIT`, `ENHANCED_LIMIT`, `AT_AUCTION`, `AT_AUCTION_LIMIT`, `MARKET_ON_OPEN`, `TRAILING_STOP_LOSS`, `TRAILING_STOP_LOSS_LIMIT`, `TOUCH_MKT`, `TOUCH_LMT`, `ODD_LOT_LIMIT` |
| `entrust_type` | string | yes | Specifies the method for placing the order. - QTY: Order specified by quantity of shares or units. - AMOUNT: Order specified by total cash amount, applicable for fractional share trading of US stocks. — one of: `QTY`, `AMOUNT` |
| `support_trading_session` | string |  | Specifies the trading session for the order. Applicable to U.S. stock market orders only. Deprecated values: - Y: [Deprecated]Include extended trading hours. - N: [Deprecated]Only support regular trading hours. Active values: - NIGHT: Only supports night trading. - ALL: Include extended trading hours. - CORE: Only support regular trading hours. - ALL_DAY: Included Overnight Hours, 8:00 p.m.ET - 8:00 p.m.ET(the next day) — one of: `Y`, `N`, `NIGHT`, `ALL`, `CORE`, `ALL_DAY` |
| `time_in_force` | string | yes | Specifies the duration for which the order remains active in the market (Time-In-Force). - DAY: The order is valid only for the current trading day and expires at the end of the day. - GTD: order that will automatically expire and be cancelled at a specific future date and time,Currently only supports the US market. - GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days). — one of: `DAY`, `GTD`, `GTC` |
| `side` | string | yes | The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash). — one of: `BUY`, `SELL`, `SHORT` |
| `quantity` | string |  | Transaction quantity. You can specify decimals when placing fractional lot orders for US stocks. |
| `total_cash_amount` | string |  | The total order amount is currently only applicable to US stock fractional share transactions and when the order is placed by amount. |
| `limit_price` | string |  | Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT or TOUCH_LMT. Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed. |
| `stop_price` | string |  | Stop price of the order. Required when order_type is STOP_LOSS, STOP_LOSS_LIMIT, TOUCH_MKT or TOUCH_LMT. Specifies the trigger price at which the stop order becomes active. |
| `trailing_type` | string |  | When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. Required when order_type is TRAILING_STOP_LOSS or TRAILING_STOP_LOSS_LIMIT. - AMOUNT: By amount. - PERCENTAGE: By percentage. — one of: `AMOUNT`, `PERCENTAGE` |
| `trailing_stop_step` | string |  | Trailing Stop Spread. If the tracking type is percentage, the tracking spread can not exceed 1,0.01 means 1%. Required when order_type is TRAILING_STOP_LOSS or TRAILING_STOP_LOSS_LIMIT. |
| `trailing_limit_price_offset` | string |  | The offset amount between the triggered stop price and the submitted limit price for a trailing stop-limit order. Required when order_type is TRAILING_STOP_LOSS_LIMIT. When triggered, the limit order price is calculated as: - Buy: limit price = stop price + trailing_limit_price_offset - Sell: limit price = stop price - trailing_limit_price_offset If the calculated limit price does not align with the instrument's tick size, it will be rounded to the nearest valid tick: rounded up for buy orders, rounded down for sell orders. |
| `trigger_price_type` | string |  | Trigger price type of the order - PRICE: Latest transaction price. - PRICE_BID: Buy at one price. - PRICE_ASK: Sell at one price. — one of: `PRICE`, `PRICE_BID`, `PRICE_ASK` |
| `sender_sub_id` | string |  | Identifier for the firm or sub-account in third-party transactions. For brokers, this field should contain the UUID of the broker user. Used to distinguish different entities or users within the same firm. |
| `no_party_ids` | array<object> |  | List of party identifiers. Applicable only for Hong Kong stock orders. Required for Relevant Regulated Intermediaries; should be omitted otherwise. |
| `expire_date` | string |  | GTD order expire date. format (UTC). The value must be in yyyy-MM-dd format |
| `option_strategy` | string |  | Type of options strategy - SINGLE: Indicates a single-leg options order — one of: `SINGLE` |
| `legs` | array<object> |  | Option leg detail. Only required when placing option orders. |

*Nested — `no_party_ids`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `party_id` | string | yes | ID of the broker client submitting the order. Format examples: - CE Number Format: ABC123 - BCAN Format: 2568 Combined format example: ABC123.2568 Must be certified by the BCAN system of the Hong Kong Stock Exchange. |
| `party_id_source` | string | yes | Source of the party ID. Value must be "D" (Proprietary/Custom Code). |
| `party_role` | string | yes | Role of the party. Value must be "3" (Client ID, BCAN Field). |

*Nested — `legs`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `instrument_type` | string | yes | Type of financial instrument associated with the request. — one of: `EQUITY`, `OPTION`, `FUTURES` |
| `market` | string | yes | Market code indicating the trading venue or regulatory region of the financial instrument.Used together with symbol and instrument_type to uniquely identify a tradable instrument. — one of: `US`, `HK`, `CN` |
| `symbol` | string | yes | Trading symbol of the financial instrument. Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives). |
| `side` | string | yes | The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash). — one of: `BUY`, `SELL`, `SHORT` |
| `strike_price` | string |  | Exercise price (strike price) of the option. Specifies the price at which the underlying asset can be bought (CALL) or sold (PUT) upon exercise. |
| `option_expire_date` | string |  | Expiration date. Format: yyyy-MM-dd |
| `option_type` | string |  | Type of the option. - CALL: Right to buy the underlying asset. - PUT: Right to sell the underlying asset. — one of: `CALL`, `PUT` |
| `quantity` | string |  | Quantity of the order or strategy leg. For stock legs, specifies the number of shares to transact For option legs, specifies the number of option contracts to transact for this leg and is expressed in whole contracts. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string |  | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order for NORMAL order. |
| `order_id` | string |  | System-generated order identifier. Returned in the response for simple orders. Represents the unique Webull order ID assigned by the system when placing the order for NORMAL order. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Batch Place Orders

`POST /trading/orders/batch-place`

> Places multiple orders in a single request. A maximum of 50 orders can be submitted once, Currently only stocks are supported. This service is not currently available to all clients. Please contact Webull if you require assistance.

| | |
|---|---|
| **SDK** | `trade.BatchPlaceOrder` |
| **Reference** | [order-batch-place.md](https://developer.webull.com/apis/docs/reference/order-batch-place.md) |
| **Note** | Equity only, up to 50 orders. |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `account_id` | string | yes | Account identifier |
| `batch_orders` | array<object> | yes | Batch Orders |

*Nested — `batch_orders`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string | yes | Unique client-defined identifier for the order. Maximum length is 32 characters and must be unique per account. Used to track or reference the order when interacting with the system. |
| `combo_type` | string | yes | Type of order combination. Currently only NORMAL is supported. It may be expanded to support other types in the future - NORMAL: Indicates a standard single order |
| `instrument_type` | string | yes | Type of financial instrument associated with the request. — one of: `EQUITY` |
| `entrust_type` | string | yes | Specifies the method for placing the order. - QTY: Order specified by quantity of shares or units. — one of: `QTY` |
| `support_trading_session` | string | yes | Specifies the trading session for the order. Applicable to U.S. stock market orders only. Algorithmic trading order currently supports only regular trading hours. - NIGHT: Only supports night trading. - ALL: Include extended trading hours. - CORE: Only support regular trading hours. — one of: `ALL`, `CORE`, `NIGHT` |
| `symbol` | string | yes | Trading symbol of the financial instrument. Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities). |
| `market` | string | yes | Market code indicating the trading venue or regulatory region of the financial instrument. Used together with symbol and instrument_type to uniquely identify a tradable instrument. — one of: `US` |
| `side` | string | yes | The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash). — one of: `BUY`, `SELL` |
| `order_type` | string | yes | Specifies the type of order to be placed. Determines how the order will be executed in the market. Available order types depend on the market and instrument type. - <b>LIMIT:</b> Limit Order - <b>MARKET:</b> Market Order — one of: `MARKET`, `LIMIT` |
| `time_in_force` | string | yes | Specifies the duration for which the order remains active in the market (Time-In-Force). - DAY: The order is valid only for the current trading day and expires at the end of the day. — one of: `DAY` |
| `quantity` | string | yes | Transaction quantity. You can specify decimals when placing fractional lot orders for US stocks. |
| `limit_price` | string |  | Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT. Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed. When event_trade_mode is set for an event contract trade, limit_price is not required. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `total` | integer | yes | The total number of orders submitted each time |
| `success` | integer | yes | The number of orders successfully submitted to the webull system |
| `failed` | integer | yes | The number of failed order submitted to the webull system |
| `batch_orders` | array<object> | yes | Batch Order place result |

*Nested — `batch_orders`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string |  | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order for NORMAL order. |
| `order_id` | string |  | System-generated order identifier. Returned in the response for simple orders. Represents the unique Webull order ID assigned by the system when placing the order for NORMAL order. |
| `error_code` | string |  | Order place failed code |
| `message` | string |  | Order place failed and detail failed reason. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Order Replace

`POST /trading/orders/replace`

> Modifies equity and options orders.

| | |
|---|---|
| **SDK** | `trade.ReplaceOrder` |
| **Reference** | [common-order-replace.md](https://developer.webull.hk/apis/docs/reference/common-order-replace.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `account_id` | string | yes | Account identifier |
| `modify_orders` | array<object> | yes | Order Details |

*Nested — `modify_orders`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string | yes | Unique client-defined identifier for the order. Maximum length is 32 characters and must be unique per account. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). Used to track or reference the order when interacting with the system. |
| `time_in_force` | string |  | Specifies the duration for which the order remains active in the market (Time-In-Force). - DAY: The order is valid only for the current trading day and expires at the end of the day. - GTD: order that will automatically expire and be cancelled at a specific future date and time,Currently only supports the US market. - GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days). — one of: `DAY`, `GTD`, `GTC` |
| `quantity` | string |  | Transaction quantity. You can specify decimals when placing fractional lot orders for US stocks. |
| `expire_date` | string |  | GTD order expire date. format (UTC). The value must be in yyyy-MM-dd format |
| `limit_price` | string |  | Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT or TOUCH_LMT. Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed. |
| `stop_price` | string |  | Stop price of the order. Required when order_type is STOP_LOSS, STOP_LOSS_LIMIT, TOUCH_MKT or TOUCH_LMT. Specifies the trigger price at which the stop order becomes active. |
| `trailing_type` | string |  | When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. Required when order_type is TRAILING_STOP_LOSS or TRAILING_STOP_LOSS_LIMIT. - AMOUNT: By amount. - PERCENTAGE: By percentage. — one of: `AMOUNT`, `PERCENTAGE` |
| `trailing_stop_step` | string |  | Trailing Stop Spread. If the tracking type is percentage, the tracking spread can not exceed 1,0.01 means 1%. Required when order_type is TRAILING_STOP_LOSS or TRAILING_STOP_LOSS_LIMIT. |
| `trailing_limit_price_offset` | string |  | The offset amount between the triggered stop price and the submitted limit price for a trailing stop-limit order. Required when order_type is TRAILING_STOP_LOSS_LIMIT. When triggered, the limit order price is calculated as: - Buy: limit price = stop price + trailing_limit_price_offset - Sell: limit price = stop price - trailing_limit_price_offset If the calculated limit price does not align with the instrument's tick size, it will be rounded to the nearest valid tick: rounded up for buy orders, rounded down for sell orders. |
| `trigger_price_type` | string |  | Trigger price type of the order - PRICE: Latest transaction price. - PRICE_BID: Buy at one price. - PRICE_ASK: Sell at one price. — one of: `PRICE`, `PRICE_BID`, `PRICE_ASK` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string |  | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order for NORMAL order. |
| `order_id` | string |  | System-generated order identifier. Returned in the response for simple orders. Represents the unique Webull order ID assigned by the system when placing the order for NORMAL order. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Order Cancel

`POST /trading/orders/cancel`

> Cancels orders for equities and options.

| | |
|---|---|
| **SDK** | `trade.CancelOrder` |
| **Reference** | [common-order-cancel.md](https://developer.webull.hk/apis/docs/reference/common-order-cancel.md) |

**Request body**

| Field | Type | Required | Description |
|---|---|---|---|
| `account_id` | string | yes | Account identifier |
| `client_order_id` | string | yes | Unique client-defined identifier for the order. Maximum length is 32 characters and must be unique per account. Allowed characters: letters (A–Z, a–z), digits (0–9), hyphen (-), underscore (_). Used to track or reference the order when interacting with the system. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string |  | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order for NORMAL order. |
| `order_id` | string |  | System-generated order identifier. Returned in the response for simple orders. Represents the unique Webull order ID assigned by the system when placing the order for NORMAL order. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Open Orders

`GET /trading/orders/open-orders/list`

> Retrieves pending orders by page. Orders can be modified or cancelled based on client_order_id.

| | |
|---|---|
| **SDK** | `trade.GetOpenOrders` |
| **Reference** | [order-open.md](https://developer.webull.hk/apis/docs/reference/order-open.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `account_id` | query | String | yes | Account identifier. |
| `pagination_key` | query | String |  | Pagination key from previous response for next page. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Result data list |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string | yes | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order. |
| `combo_type` | string | yes | Specifies the type of order combination. For details, please refer to Combo Order. Futures Trading currently support only the NORMAL type. - NORMAL: A standard single order. - MASTER: A primary order that triggers a take-profit or stop-loss order upon execution - STOP_PROFIT: A take-profit order - STOP_LOSS: A stop-loss order - OTO: An order that triggers another order upon execution (One-Triggers-the-Other) - OCO: A pair of orders where the execution of one cancels the other (One-Cancels-the-Other) - OTOCO: An order that triggers an OCO order set upon execution (One-Triggers-One-Cancels-the-Other) Note: When placing take-profit/stop-loss orders to sell and close an existing position, submit only STOP_PROFIT/STOP_LOSS sub-orders (side = SELL) under the same client_combo_order_id; no MASTER order is required or supported in this scenario, since no new position is being opened. Note: OTO, OCO and OTOCO combo types are only supported for stock (EQUITY) orders; option orders (including option_strategy = SINGLE) do not support OTO, OCO or OTOCO. Sub-order quantity limits by combo_type: \| Scenario \| combo_type \| Supported order_type \| Quantity \| Description \| \|---\|---\|---\|---\|---\| \| Take-Profit/Stop-Loss \| MASTER \| MARKET, LIMIT \| 1 \| Master Order \| \| Take-Profit/Stop-Loss \| STOP_PROFIT \| LIMIT \| 0-1 \| Take Profit Order \| \| Take-Profit/Stop-Loss \| STOP_LOSS \| STOP_LOSS \| 0-1 \| Stop Loss Order \| \| OTO \| MASTER \| MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 1 \| Master Order \| \| OTO \| OTO \| MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 1-6 \| Triggered Order(s) \| \| OCO \| OCO \| LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 2-6 \| Mutually Cancelling Orders \| \| OTOCO \| MASTER \| MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 1 \| Master Order \| \| OTOCO \| OTOCO \| LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 1-6 \| OCO Order Set Triggered by MASTER \| — one of: `NORMAL`, `MASTER`, `STOP_PROFIT`, `STOP_LOSS`, `OTO`, `OCO`, `OTOCO` |
| `orders` | array<object> | yes | Order Details |

*Nested — `orders`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string | yes | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order. |
| `order_id` | string | yes | System-generated order identifier. Returned in the response for simple orders. Represents the unique Webull order ID assigned by the system. |
| `symbol` | string | yes | Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives). |
| `side` | string | yes | The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash). — one of: `BUY`, `SELL`, `SHORT` |
| `status` | string | yes | - PENDING: Indicates that the order has been submitted to the exchange and is awaiting completion - SUBMITTED: Indicates that the order has been submitted to the exchange and is awaiting completion - CANCELLED: Indicates that the order has been successfully cancelled - FILLED: Indicates that the order has been fully executed - FAILED: Indicates a failed order, such as REJECTED - PARTIAL_FILLED: Refers to the portion of the order that has been completed, but not all of it has been completed — one of: `PENDING`, `SUBMITTED`, `CANCELLED`, `FILLED`, `FAILED`, `PARTIAL_FILLED` |
| `order_type` | string | yes | Specifies the type of order to be placed. Determines how the order will be executed in the market. Available order types depend on the market and instrument type. Options trading Only LIMIT,STOP_LOSS,STOP_LOSS_LIMIT are supported. U.S. Stock - <b>LIMIT:</b> Limit Order - <b>MARKET:</b> Market Order - <b>STOP_LOSS:</b> Stop Order - <b>STOP_LOSS_LIMIT:</b> Stop Limit Order - <b>MARKET_ON_OPEN:</b> Opening market order - <b>MARKET_ON_CLOSE:</b> Closing market order - <b>TOUCH_MKT:</b> Touch Market Order - <b>TOUCH_LMT:</b> Touch Limit Order - <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order - <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order Hong Kong Stock - <b>ENHANCED_LIMIT:</b> Enhanced Limit Order - <b>AT_AUCTION:</b> At-auction order - <b>AT_AUCTION_LIMIT:</b> At-auction limit order - <b>STOP_LOSS:</b> Stop Order - <b>STOP_LOSS_LIMIT:</b> Stop Limit Order - <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order - <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order - <b>TOUCH_MKT:</b> Touch Market Order - <b>TOUCH_LMT:</b> Touch Limit Order - <b>ODD_LOT_LIMIT:</b> Odd Lot Limit Order China Connect - <b>LIMIT:</b> Limit Order — one of: `LIMIT`, `MARKET`, `STOP_LOSS`, `STOP_LOSS_LIMIT`, `ENHANCED_LIMIT`, `AT_AUCTION`, `AT_AUCTION_LIMIT`, `MARKET_ON_OPEN`, `TRAILING_STOP_LOSS`, `TRAILING_STOP_LOSS_LIMIT`, `TOUCH_MKT`, `TOUCH_LMT`, `ODD_LOT_LIMIT` |
| `instrument_type` | string |  | Type of financial instrument associated with the request. — one of: `EQUITY`, `OPTION`, `FUTURES` |
| `support_trading_session` | string |  | Specifies the trading session for the order. Applicable to U.S. stock market orders only. Deprecated values: - Y: [Deprecated]Include extended trading hours. - N: [Deprecated]Only support regular trading hours. Active values: - NIGHT: Only supports night trading. - ALL: Include extended trading hours. - CORE: Only support regular trading hours. - ALL_DAY: Included Overnight Hours, 8:00 p.m.ET - 8:00 p.m.ET(the next day) — one of: `Y`, `N`, `NIGHT`, `ALL`, `CORE`, `ALL_DAY` |
| `time_in_force` | string | yes | Specifies the duration for which the order remains active in the market (Time-In-Force). - DAY: The order is valid only for the current trading day and expires at the end of the day. - GTD: order that will automatically expire and be cancelled at a specific future date and time,Currently only supports the US market. - GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days). — one of: `DAY`, `GTD`, `GTC` |
| `total_quantity` | string | yes | Total order quantity. Represents the total number of units submitted for this order. |
| `filled_quantity` | string |  | Quantity that has been executed. Represents the number of units that have been filled so far. |
| `filled_price` | string |  | Average transaction price of the filled quantity. If the order has not been executed yet, this may be zero or null. |
| `limit_price` | string |  | Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT or TOUCH_LMT. Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed. |
| `stop_price` | string |  | Stop price of the order. Required when order_type is STOP_LOSS, STOP_LOSS_LIMIT, TOUCH_MKT or TOUCH_LMT. Specifies the trigger price at which the stop order becomes active. |
| `trailing_type` | string |  | When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. Required when order_type is TRAILING_STOP_LOSS or TRAILING_STOP_LOSS_LIMIT. - AMOUNT: By amount. - PERCENTAGE: By percentage. — one of: `AMOUNT`, `PERCENTAGE` |
| `trailing_stop_step` | string |  | Trailing Stop Spread. If the tracking type is percentage, the tracking spread can not exceed 1,0.01 means 1%. Required when order_type is TRAILING_STOP_LOSS or TRAILING_STOP_LOSS_LIMIT. |
| `trailing_limit_price_offset` | string |  | The offset amount between the triggered stop price and the submitted limit price for a trailing stop-limit order. Required when order_type is TRAILING_STOP_LOSS_LIMIT. When triggered, the limit order price is calculated as: - Buy: limit price = stop price + trailing_limit_price_offset - Sell: limit price = stop price - trailing_limit_price_offset If the calculated limit price does not align with the instrument's tick size, it will be rounded to the nearest valid tick: rounded up for buy orders, rounded down for sell orders. |
| `trigger_price_type` | string |  | Trigger price type of the order - PRICE: Latest transaction price. - PRICE_BID: Buy at one price. - PRICE_ASK: Sell at one price. — one of: `PRICE`, `PRICE_BID`, `PRICE_ASK` |
| `place_time` | string |  | Order placement time in milliseconds since Unix epoch. |
| `place_time_at` | string | yes | Order placement time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ |
| `filled_time` | string |  | Time of the last executed trade in milliseconds since Unix epoch. |
| `filled_time_at` | string |  | Time of the last executed trade in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ |
| `legs` | array<object> |  | Leg detail |

*Nested — `legs`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives). |
| `side` | string | yes | The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash). — one of: `BUY`, `SELL`, `SHORT` |
| `quantity` | string | yes | Quantity of the order. Specifies the number of shares or units to transact. For US stocks, fractional quantities are allowed and can include decimals. |
| `option_type` | string | yes | Type of the option. - CALL: Right to buy the underlying asset. - PUT: Right to sell the underlying asset. — one of: `CALL`, `PUT` |
| `option_category` | string | yes | Category of the option, indicating its exercise style. Possible values: - AMERICAN: Can be exercised any time before expiration. - EUROPEAN: Can only be exercised at expiration. — one of: `AMERICAN`, `EUROPEAN` |
| `option_strategy` | string |  | Type of options strategy - SINGLE: Indicates a single-leg options order — one of: `SINGLE` |
| `strike_price` | string | yes | Exercise price (strike price) of the option. Specifies the price at which the underlying asset can be bought (CALL) or sold (PUT) upon exercise. |
| `option_contract_multiplier` | string |  | The number of shares corresponding to each option contract |
| `option_contract_deliverable` | string |  | The number of shares required to exercise each contract |
| `option_expire_date` | string | yes | Expiration date of the option. Format: yyyy-MM-dd. After this date, the option will no longer be valid. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Order History

`GET /trading/orders/historical-orders/list`

> Retrieves historical orders for the past 7 days. If orders are group orders, they will be returned together.

| | |
|---|---|
| **SDK** | `trade.GetOrderHistory` |
| **Reference** | [order-history.md](https://developer.webull.hk/apis/docs/reference/order-history.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `account_id` | query | String | yes | Account identifier. |
| `start_time` | query | String |  | The start date of the query period. If not provided, the default query period is the last 7 days. Users can specify an earlier date, but the maximum allowed look-back period is 6 months. Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'. |
| `end_time` | query | String |  | The end time of the query period. If not provided, the default query period is the last 7 days. Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'. |
| `pagination_key` | query | String |  | Pagination key from previous response for next page. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Result data list |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string | yes | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order. |
| `combo_type` | string | yes | Specifies the type of order combination. For details, please refer to Combo Order. Futures Trading currently support only the NORMAL type. - NORMAL: A standard single order. - MASTER: A primary order that triggers a take-profit or stop-loss order upon execution - STOP_PROFIT: A take-profit order - STOP_LOSS: A stop-loss order - OTO: An order that triggers another order upon execution (One-Triggers-the-Other) - OCO: A pair of orders where the execution of one cancels the other (One-Cancels-the-Other) - OTOCO: An order that triggers an OCO order set upon execution (One-Triggers-One-Cancels-the-Other) Note: When placing take-profit/stop-loss orders to sell and close an existing position, submit only STOP_PROFIT/STOP_LOSS sub-orders (side = SELL) under the same client_combo_order_id; no MASTER order is required or supported in this scenario, since no new position is being opened. Note: OTO, OCO and OTOCO combo types are only supported for stock (EQUITY) orders; option orders (including option_strategy = SINGLE) do not support OTO, OCO or OTOCO. Sub-order quantity limits by combo_type: \| Scenario \| combo_type \| Supported order_type \| Quantity \| Description \| \|---\|---\|---\|---\|---\| \| Take-Profit/Stop-Loss \| MASTER \| MARKET, LIMIT \| 1 \| Master Order \| \| Take-Profit/Stop-Loss \| STOP_PROFIT \| LIMIT \| 0-1 \| Take Profit Order \| \| Take-Profit/Stop-Loss \| STOP_LOSS \| STOP_LOSS \| 0-1 \| Stop Loss Order \| \| OTO \| MASTER \| MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 1 \| Master Order \| \| OTO \| OTO \| MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 1-6 \| Triggered Order(s) \| \| OCO \| OCO \| LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 2-6 \| Mutually Cancelling Orders \| \| OTOCO \| MASTER \| MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 1 \| Master Order \| \| OTOCO \| OTOCO \| LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 1-6 \| OCO Order Set Triggered by MASTER \| — one of: `NORMAL`, `MASTER`, `STOP_PROFIT`, `STOP_LOSS`, `OTO`, `OCO`, `OTOCO` |
| `orders` | array<object> | yes | Order Details |

*Nested — `orders`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string | yes | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order. |
| `order_id` | string | yes | System-generated order identifier. Returned in the response for simple orders. Represents the unique Webull order ID assigned by the system. |
| `symbol` | string | yes | Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives). |
| `side` | string | yes | The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash). — one of: `BUY`, `SELL`, `SHORT` |
| `status` | string | yes | - PENDING: Indicates that the order has been submitted to the exchange and is awaiting completion - SUBMITTED: Indicates that the order has been submitted to the exchange and is awaiting completion - CANCELLED: Indicates that the order has been successfully cancelled - FILLED: Indicates that the order has been fully executed - FAILED: Indicates a failed order, such as REJECTED - PARTIAL_FILLED: Refers to the portion of the order that has been completed, but not all of it has been completed — one of: `PENDING`, `SUBMITTED`, `CANCELLED`, `FILLED`, `FAILED`, `PARTIAL_FILLED` |
| `order_type` | string | yes | Specifies the type of order to be placed. Determines how the order will be executed in the market. Available order types depend on the market and instrument type. Options trading Only LIMIT,STOP_LOSS,STOP_LOSS_LIMIT are supported. U.S. Stock - <b>LIMIT:</b> Limit Order - <b>MARKET:</b> Market Order - <b>STOP_LOSS:</b> Stop Order - <b>STOP_LOSS_LIMIT:</b> Stop Limit Order - <b>MARKET_ON_OPEN:</b> Opening market order - <b>MARKET_ON_CLOSE:</b> Closing market order - <b>TOUCH_MKT:</b> Touch Market Order - <b>TOUCH_LMT:</b> Touch Limit Order - <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order - <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order Hong Kong Stock - <b>ENHANCED_LIMIT:</b> Enhanced Limit Order - <b>AT_AUCTION:</b> At-auction order - <b>AT_AUCTION_LIMIT:</b> At-auction limit order - <b>STOP_LOSS:</b> Stop Order - <b>STOP_LOSS_LIMIT:</b> Stop Limit Order - <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order - <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order - <b>TOUCH_MKT:</b> Touch Market Order - <b>TOUCH_LMT:</b> Touch Limit Order - <b>ODD_LOT_LIMIT:</b> Odd Lot Limit Order China Connect - <b>LIMIT:</b> Limit Order — one of: `LIMIT`, `MARKET`, `STOP_LOSS`, `STOP_LOSS_LIMIT`, `ENHANCED_LIMIT`, `AT_AUCTION`, `AT_AUCTION_LIMIT`, `MARKET_ON_OPEN`, `TRAILING_STOP_LOSS`, `TRAILING_STOP_LOSS_LIMIT`, `TOUCH_MKT`, `TOUCH_LMT`, `ODD_LOT_LIMIT` |
| `instrument_type` | string |  | Type of financial instrument associated with the request. — one of: `EQUITY`, `OPTION`, `FUTURES` |
| `support_trading_session` | string |  | Specifies the trading session for the order. Applicable to U.S. stock market orders only. Deprecated values: - Y: [Deprecated]Include extended trading hours. - N: [Deprecated]Only support regular trading hours. Active values: - NIGHT: Only supports night trading. - ALL: Include extended trading hours. - CORE: Only support regular trading hours. - ALL_DAY: Included Overnight Hours, 8:00 p.m.ET - 8:00 p.m.ET(the next day) — one of: `Y`, `N`, `NIGHT`, `ALL`, `CORE`, `ALL_DAY` |
| `time_in_force` | string | yes | Specifies the duration for which the order remains active in the market (Time-In-Force). - DAY: The order is valid only for the current trading day and expires at the end of the day. - GTD: order that will automatically expire and be cancelled at a specific future date and time,Currently only supports the US market. - GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days). — one of: `DAY`, `GTD`, `GTC` |
| `total_quantity` | string | yes | Total order quantity. Represents the total number of units submitted for this order. |
| `filled_quantity` | string |  | Quantity that has been executed. Represents the number of units that have been filled so far. |
| `filled_price` | string |  | Average transaction price of the filled quantity. If the order has not been executed yet, this may be zero or null. |
| `limit_price` | string |  | Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT or TOUCH_LMT. Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed. |
| `stop_price` | string |  | Stop price of the order. Required when order_type is STOP_LOSS, STOP_LOSS_LIMIT, TOUCH_MKT or TOUCH_LMT. Specifies the trigger price at which the stop order becomes active. |
| `trailing_type` | string |  | When market continues to fall, the stop price to buy follows, or trails, the lowest price of a stock by a trail that you set. Required when order_type is TRAILING_STOP_LOSS or TRAILING_STOP_LOSS_LIMIT. - AMOUNT: By amount. - PERCENTAGE: By percentage. — one of: `AMOUNT`, `PERCENTAGE` |
| `trailing_stop_step` | string |  | Trailing Stop Spread. If the tracking type is percentage, the tracking spread can not exceed 1,0.01 means 1%. Required when order_type is TRAILING_STOP_LOSS or TRAILING_STOP_LOSS_LIMIT. |
| `trailing_limit_price_offset` | string |  | The offset amount between the triggered stop price and the submitted limit price for a trailing stop-limit order. Required when order_type is TRAILING_STOP_LOSS_LIMIT. When triggered, the limit order price is calculated as: - Buy: limit price = stop price + trailing_limit_price_offset - Sell: limit price = stop price - trailing_limit_price_offset If the calculated limit price does not align with the instrument's tick size, it will be rounded to the nearest valid tick: rounded up for buy orders, rounded down for sell orders. |
| `trigger_price_type` | string |  | Trigger price type of the order - PRICE: Latest transaction price. - PRICE_BID: Buy at one price. - PRICE_ASK: Sell at one price. — one of: `PRICE`, `PRICE_BID`, `PRICE_ASK` |
| `place_time` | string |  | Order placement time in milliseconds since Unix epoch. |
| `place_time_at` | string | yes | Order placement time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ |
| `filled_time` | string |  | Time of the last executed trade in milliseconds since Unix epoch. |
| `filled_time_at` | string |  | Time of the last executed trade in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ |
| `legs` | array<object> |  | Leg detail |

*Nested — `legs`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives). |
| `side` | string | yes | The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash). — one of: `BUY`, `SELL`, `SHORT` |
| `quantity` | string | yes | Quantity of the order. Specifies the number of shares or units to transact. For US stocks, fractional quantities are allowed and can include decimals. |
| `option_type` | string | yes | Type of the option. - CALL: Right to buy the underlying asset. - PUT: Right to sell the underlying asset. — one of: `CALL`, `PUT` |
| `option_category` | string | yes | Category of the option, indicating its exercise style. Possible values: - AMERICAN: Can be exercised any time before expiration. - EUROPEAN: Can only be exercised at expiration. — one of: `AMERICAN`, `EUROPEAN` |
| `option_strategy` | string |  | Type of options strategy - SINGLE: Indicates a single-leg options order — one of: `SINGLE` |
| `strike_price` | string | yes | Exercise price (strike price) of the option. Specifies the price at which the underlying asset can be bought (CALL) or sold (PUT) upon exercise. |
| `option_contract_multiplier` | string |  | The number of shares corresponding to each option contract |
| `option_contract_deliverable` | string |  | The number of shares required to exercise each contract |
| `option_expire_date` | string | yes | Expiration date of the option. Format: yyyy-MM-dd. After this date, the option will no longer be valid. |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Order Detail

`GET /trading/orders/get`

> Retrieves the specified order details through the order ID.

| | |
|---|---|
| **SDK** | `trade.GetOrderDetail` |
| **Reference** | [order-detail.md](https://developer.webull.hk/apis/docs/reference/order-detail.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `account_id` | query | String | yes | Account identifier. |
| `client_order_id` | query | String | yes | Unique client-defined identifier for the order. Maximum length is 32 characters and must be unique per account. Used to track or reference the order when interacting with the system. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string | yes | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order. |
| `combo_type` | string | yes | Specifies the type of order combination. For details, please refer to Combo Order. Futures Trading currently support only the NORMAL type. - NORMAL: A standard single order. - MASTER: A primary order that triggers a take-profit or stop-loss order upon execution - STOP_PROFIT: A take-profit order - STOP_LOSS: A stop-loss order - OTO: An order that triggers another order upon execution (One-Triggers-the-Other) - OCO: A pair of orders where the execution of one cancels the other (One-Cancels-the-Other) - OTOCO: An order that triggers an OCO order set upon execution (One-Triggers-One-Cancels-the-Other) Note: When placing take-profit/stop-loss orders to sell and close an existing position, submit only STOP_PROFIT/STOP_LOSS sub-orders (side = SELL) under the same client_combo_order_id; no MASTER order is required or supported in this scenario, since no new position is being opened. Note: OTO, OCO and OTOCO combo types are only supported for stock (EQUITY) orders; option orders (including option_strategy = SINGLE) do not support OTO, OCO or OTOCO. Sub-order quantity limits by combo_type: \| Scenario \| combo_type \| Supported order_type \| Quantity \| Description \| \|---\|---\|---\|---\|---\| \| Take-Profit/Stop-Loss \| MASTER \| MARKET, LIMIT \| 1 \| Master Order \| \| Take-Profit/Stop-Loss \| STOP_PROFIT \| LIMIT \| 0-1 \| Take Profit Order \| \| Take-Profit/Stop-Loss \| STOP_LOSS \| STOP_LOSS \| 0-1 \| Stop Loss Order \| \| OTO \| MASTER \| MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 1 \| Master Order \| \| OTO \| OTO \| MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 1-6 \| Triggered Order(s) \| \| OCO \| OCO \| LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 2-6 \| Mutually Cancelling Orders \| \| OTOCO \| MASTER \| MARKET, LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 1 \| Master Order \| \| OTOCO \| OTOCO \| LIMIT, STOP_LOSS, STOP_LOSS_LIMIT \| 1-6 \| OCO Order Set Triggered by MASTER \| — one of: `NORMAL`, `MASTER`, `STOP_PROFIT`, `STOP_LOSS`, `OTO`, `OCO`, `OTOCO` |
| `orders` | array<object> | yes | Order Details |

*Nested — `orders`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `client_order_id` | string | yes | Client-defined order identifier. Returned in the response for simple orders. Represents the unique order ID assigned by the user when placing the order. |
| `order_id` | string | yes | System-generated order identifier. Returned in the response for simple orders. Represents the unique Webull order ID assigned by the system. |
| `symbol` | string | yes | Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives). |
| `side` | string | yes | The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash). — one of: `BUY`, `SELL`, `SHORT` |
| `status` | string | yes | - PENDING: Indicates that the order has been submitted to the exchange and is awaiting completion - SUBMITTED: Indicates that the order has been submitted to the exchange and is awaiting completion - CANCELLED: Indicates that the order has been successfully cancelled - FILLED: Indicates that the order has been fully executed - FAILED: Indicates a failed order, such as REJECTED - PARTIAL_FILLED: Refers to the portion of the order that has been completed, but not all of it has been completed — one of: `PENDING`, `SUBMITTED`, `CANCELLED`, `FILLED`, `FAILED`, `PARTIAL_FILLED` |
| `order_type` | string | yes | Specifies the type of order to be placed. Determines how the order will be executed in the market. Available order types depend on the market and instrument type. Options trading Only LIMIT,STOP_LOSS,STOP_LOSS_LIMIT are supported. U.S. Stock - <b>LIMIT:</b> Limit Order - <b>MARKET:</b> Market Order - <b>STOP_LOSS:</b> Stop Order - <b>STOP_LOSS_LIMIT:</b> Stop Limit Order - <b>MARKET_ON_OPEN:</b> Opening market order - <b>MARKET_ON_CLOSE:</b> Closing market order - <b>TOUCH_MKT:</b> Touch Market Order - <b>TOUCH_LMT:</b> Touch Limit Order - <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order - <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order Hong Kong Stock - <b>ENHANCED_LIMIT:</b> Enhanced Limit Order - <b>AT_AUCTION:</b> At-auction order - <b>AT_AUCTION_LIMIT:</b> At-auction limit order - <b>STOP_LOSS:</b> Stop Order - <b>STOP_LOSS_LIMIT:</b> Stop Limit Order - <b>TRAILING_STOP_LOSS:</b> Trailing Stop Order - <b>TRAILING_STOP_LOSS_LIMIT:</b> Trailing Stop Limit Order - <b>TOUCH_MKT:</b> Touch Market Order - <b>TOUCH_LMT:</b> Touch Limit Order - <b>ODD_LOT_LIMIT:</b> Odd Lot Limit Order China Connect - <b>LIMIT:</b> Limit Order — one of: `LIMIT`, `MARKET`, `STOP_LOSS`, `STOP_LOSS_LIMIT`, `ENHANCED_LIMIT`, `AT_AUCTION`, `AT_AUCTION_LIMIT`, `MARKET_ON_OPEN`, `TRAILING_STOP_LOSS`, `TRAILING_STOP_LOSS_LIMIT`, `TOUCH_MKT`, `TOUCH_LMT`, `ODD_LOT_LIMIT` |
| `instrument_type` | string |  | Type of financial instrument associated with the request. — one of: `EQUITY`, `OPTION`, `FUTURES` |
| `support_trading_session` | string |  | Specifies the trading session for the order. Applicable to U.S. stock market orders only. Deprecated values: - Y: [Deprecated]Include extended trading hours. - N: [Deprecated]Only support regular trading hours. Active values: - NIGHT: Only supports night trading. - ALL: Include extended trading hours. - CORE: Only support regular trading hours. - ALL_DAY: Included Overnight Hours, 8:00 p.m.ET - 8:00 p.m.ET(the next day) — one of: `Y`, `N`, `NIGHT`, `ALL`, `CORE`, `ALL_DAY` |
| `time_in_force` | string | yes | Specifies the duration for which the order remains active in the market (Time-In-Force). - DAY: The order is valid only for the current trading day and expires at the end of the day. - GTD: order that will automatically expire and be cancelled at a specific future date and time,Currently only supports the US market. - GTC: Good-Till-Canceled, the order remains active until it is executed, explicitly canceled, or reaches the maximum allowed duration (typically 90 days). — one of: `DAY`, `GTD`, `GTC` |
| `total_quantity` | string | yes | Total order quantity. Represents the total number of units submitted for this order. |
| `filled_quantity` | string |  | Quantity that has been executed. Represents the number of units that have been filled so far. |
| `filled_price` | string |  | Average transaction price of the filled quantity. If the order has not been executed yet, this may be zero or null. |
| `limit_price` | string |  | Limit price of the order. Required when order_type is LIMIT, STOP_LOSS_LIMIT. Specifies the maximum (for buy) or minimum (for sell) price at which the order can be executed. |
| `stop_price` | string |  | Stop price of the order. Required when order_type is STOP_LOSS, STOP_LOSS_LIMIT. Specifies the trigger price at which the stop order becomes active. |
| `place_time` | string |  | Order placement time in milliseconds since Unix epoch. |
| `place_time_at` | string | yes | Order placement time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ |
| `filled_time` | string |  | Time of the last executed trade in milliseconds since Unix epoch. |
| `filled_time_at` | string |  | Time of the last executed trade in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSSZ |
| `legs` | array<object> |  | Leg detail |
| `commission` | object |  | Commission breakdown |
| `fees` | array<object> |  | Fee breakdown |

*Nested — `legs`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string | yes | Trading symbol of the financial instrument.Represents the unique identifier of the security in the specified market (e.g., ticker symbol for equities or option symbol code for derivatives). |
| `side` | string | yes | The order side indicating the intended trading direction of the transaction. The meaning of side may vary depending on the instrument_type and account type (e.g., margin vs. cash). — one of: `BUY`, `SELL`, `SHORT` |
| `quantity` | string | yes | Quantity of the order. Specifies the number of shares or units to transact. For US stocks, fractional quantities are allowed and can include decimals. |
| `option_type` | string | yes | Type of the option. - CALL: Right to buy the underlying asset. - PUT: Right to sell the underlying asset. — one of: `CALL`, `PUT` |
| `option_category` | string | yes | Category of the option, indicating its exercise style. Possible values: - AMERICAN: Can be exercised any time before expiration. - EUROPEAN: Can only be exercised at expiration. — one of: `AMERICAN`, `EUROPEAN` |
| `option_strategy` | string |  | Type of options strategy - SINGLE: Indicates a single-leg options order — one of: `SINGLE` |
| `strike_price` | string | yes | Exercise price (strike price) of the option. Specifies the price at which the underlying asset can be bought (CALL) or sold (PUT) upon exercise. |
| `option_contract_multiplier` | string |  | The number of shares corresponding to each option contract |
| `option_contract_deliverable` | string |  | The number of shares required to exercise each contract |
| `option_expire_date` | string | yes | Expiration date of the option. Format: yyyy-MM-dd. After this date, the option will no longer be valid. |

*Nested — `commission`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `actual_commission` | string |  | Actual commission collected |
| `receivable_commission` | string |  | Receivable commission |

*Nested — `fees`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `type` | string |  | Fee type |
| `actual_value` | string |  | Actual fee collected |
| `receivable_value` | string |  | Receivable fee |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Cash Activities

`GET /trading/activities/cash-activities/list`

> Lists an account's cash activities, filterable by type and time range. Defaults to the last 7 days if no date is provided.

| | |
|---|---|
| **SDK** | `trade.GetCashActivities` |
| **Reference** | [trade-cash-activity-by-type.md](https://developer.webull.com/apis/docs/reference/trade-cash-activity-by-type.md) |
| **Note** | Documented on the US site. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `account_id` | query | String | yes | Provide the target account id |
| `activity_types` | query | string |  | Account activity types. Note: EC_STATEMENT is deprecated; use EC_SETTLEMENT instead (EC_STATEMENT behaves the same as EC_SETTLEMENT). Note: Crypto accounts only support: TRADE, DEPOSIT, WITHDRAW, and FEES. — one of: `TRADE`, `DEPOSIT`, `WITHDRAW`, `FEES`, `TRANSFER`, `DIVIDENDS`, `TAX`, `INTERESTS`, `CORPORATE_ACTION`, `OPTION_EA`, `EC_SETTLEMENT`, `JOURNAL` |
| `start_time` | query | String |  | Activity query start time, time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'. Cross-year queries are not supported; start_time and end_time must be within the same year. |
| `end_time` | query | String |  | If not provided, the default query is the last 7 days. Activity query end time, time in ISO8601 format (UTC). Format: yyyy-MM-dd'T'HH:mm:ss.SSS'Z'. Cross-year queries are not supported; start_time and end_time must be within the same year. |
| `pagination_key` | query | String |  | Pagination key from previous response for next page. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | Result data list |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `id` | string | yes | Unique ID |
| `account_id` | string | yes | Account ID |
| `activity_type` | string | yes | Activity Type \| Code \| Description \| \|------------------\|--------------------------------------\| \| TRADE \| Trade Activity Type \| \| DEPOSIT \| Deposit Activity Type \| \| WITHDRAW \| Withdraw Activity Type \| \| FEES \| Fees Activity Type \| \| TRANSFER \| Transfer Activity Type \| \| DIVIDENDS \| Dividends Activity Type \| \| TAX \| Tax Activity Type \| \| INTERESTS \| Interests Activity Type \| \| CORPORATE_ACTION \| Corporate Action Activity Type \| \| OPTION_EA \| Option Exercise and Assignment Type \| \| JOURNAL \| Journal Activity Type \| \| EC_SETTLEMENT \| EC Settlement Activity Type \| \| OTHER \| Other Activity Type \| — one of: `TRADE`, `DEPOSIT`, `WITHDRAW`, `FEES`, `TRANSFER`, `DIVIDENDS`, `TAX`, `INTERESTS`, `CORPORATE_ACTION`, `OPTION_EA`, `JOURNAL`, `EC_SETTLEMENT`, `OTHER` |
| `activity_sub_type` | string | yes | Activity Type and Sub Type Mapping. \| ActivityType \| ActivitySubType \| \|------------------\|---------------------------\| \| TRADE \| BUY \| \| TRADE \| SELL \| \| TRADE \| BUY_CANCELLED \| \| TRADE \| SELL_CANCELLED \| \| TRADE \| PENNY_FOR_LOT \| \| TRADE \| FX_EXCHANGE \| \| DEPOSIT \| WIRE \| \| DEPOSIT \| ACH \| \| DEPOSIT \| REVERSAL \| \| DEPOSIT \| CHECK \| \| DEPOSIT \| ACH_REVERSE \| \| DEPOSIT \| INTERNAL_TRANSFER \| \| WITHDRAW \| WIRE \| \| WITHDRAW \| ACH \| \| WITHDRAW \| REVERSAL \| \| WITHDRAW \| CHECK \| \| WITHDRAW \| INTERNAL_TRANSFER \| \| FEES \| WIRE_FEE \| \| FEES \| REVERSAL_FEE \| \| FEES \| TRANSFER_ACATS \| \| FEES \| CA_HANDLING_FEE \| \| FEES \| ADR \| \| FEES \| PAPER_STATEMENT_FEE \| \| FEES \| PAPER_CONFIRM_FEE \| \| FEES \| WRITE_OFF \| \| FEES \| CHECK_FEE \| \| FEES \| ADVISORY_FEE \| \| FEES \| SUBSCRIPTION_FEE \| \| FEES \| SERVICE_FEE \| \| FEES \| ACH_REVERSE_FEE \| \| FEES \| OTHER \| \| TRANSFER \| ACATS_IN \| \| TRANSFER \| ACATS_OUT \| \| TRANSFER \| INTERNAL_TRANSFER \| \| TRANSFER \| BANK_SWEEP \| \| DIVIDENDS \| INCOME \| \| DIVIDENDS \| PAYMENT_IN_LIEU \| \| DIVIDENDS \| CASH_IN_LIEU \| \| DIVIDENDS \| LP_DISTRIBUTION \| \| TAX \| FOREIGN_TAX_WITHHELD \| \| TAX \| US_TAX_WITHHOLDING \| \| TAX \| IRA_FED_WITHHOLDING \| \| TAX \| STATE_WITHHOLDING \| \| TAX \| WITHHOLDING_TAX \| \| TAX \| GP_TAX_WITHHELD \| \| TAX \| OTHER \| \| INTERESTS \| CREDIT \| \| INTERESTS \| DEBIT \| \| INTERESTS \| STOCK_BORROW_INTEREST \| \| INTERESTS \| SECURITIES_LENDING_INCOME \| \| INTERESTS \| ADJUSTMENT \| \| INTERESTS \| PAYMENT \| \| INTERESTS \| INTEREST_REBATE \| \| CORPORATE_ACTION \| CASH_IN_LIEU \| \| CORPORATE_ACTION \| REDEMPTION \| \| CORPORATE_ACTION \| MISC_ADJUSTMENT \| \| CORPORATE_ACTION \| MERGER \| \| CORPORATE_ACTION \| RIGHTS_OFFERING \| \| CORPORATE_ACTION \| IDENTIFIER_CHANGE \| \| CORPORATE_ACTION \| REVERSE_SPLIT \| \| CORPORATE_ACTION \| FORWARD_SPLIT \| \| CORPORATE_ACTION \| SPIN_OFF \| \| CORPORATE_ACTION \| CONVERSION \| \| CORPORATE_ACTION \| LIQUIDATION \| \| OPTION_EA \| CONTRACT_CLOSE \| \| OPTION_EA \| OPTION_ASSIGNMENT \| \| OPTION_EA \| OPTION_EXPIRATION \| \| OPTION_EA \| OPTION_EXERCISE \| \| JOURNAL \| CASH_JOURNAL \| \| EC_SETTLEMENT \| EC_EXPIRATION \| \| EC_SETTLEMENT \| EC_PAYOUT \| \| OTHER \| INCOME \| \| OTHER \| LENDING_REBATE \| \| OTHER \| DVP \| \| OTHER \| GRID_TRANSFER \| \| OTHER \| OTHER \| — one of: `BUY`, `SELL`, `BUY_CANCELLED`, `SELL_CANCELLED`, `PENNY_FOR_LOT`, `TRADE`, `FX_EXCHANGE`, `OPTION_EXPIRATION`, `WIRE`, `ACH`, `REVERSAL`, `CHECK`, `ACH_REVERSE`, `WIRE_FEE`, `REVERSAL_FEE`, `TRANSFER_ACATS`, `CA_HANDLING_FEE`, `ADR`, `PAPER_STATEMENT_FEE`, `PAPER_CONFIRM_FEE`, `WRITE_OFF`, `CHECK_FEE`, `ADVISORY_FEE`, `SUBSCRIPTION_FEE`, `SERVICE_FEE`, `ACH_REVERSE_FEE`, `ACATS_IN`, `ACATS_OUT`, `INTERNAL_TRANSFER`, `BANK_SWEEP`, `INCOME`, `PAYMENT_IN_LIEU`, `CASH_IN_LIEU`, `LP_DISTRIBUTION`, `FOREIGN_TAX_WITHHELD`, `US_TAX_WITHHOLDING`, `IRA_FED_WITHHOLDING`, `STATE_WITHHOLDING`, `WITHHOLDING_TAX`, `GP_TAX_WITHHELD`, `CREDIT`, `DEBIT`, `STOCK_BORROW_INTEREST`, `SECURITIES_LENDING_INCOME`, `ADJUSTMENT`, `PAYMENT`, `INTEREST_REBATE`, `REDEMPTION`, `MISC_ADJUSTMENT`, `MERGER`, `RIGHTS_OFFERING`, `IDENTIFIER_CHANGE`, `REVERSE_SPLIT`, `FORWARD_SPLIT`, `SPIN_OFF`, `CONVERSION`, `LIQUIDATION`, `CONTRACT_CLOSE`, `OPTION_ASSIGNMENT`, `OPTION_EXERCISE`, `CASH_JOURNAL`, `EC_EXPIRATION`, `EC_PAYOUT`, `LENDING_REBATE`, `DVP`, `GRID_TRANSFER`, `OTHER` |
| `currency` | string | yes | Currency — one of: `USD` |
| `market` | string |  | Market Code US - US Market — one of: `US` |
| `symbol` | string |  | Activity Symbol |
| `trade_date` | string | yes | Accounting date of the transaction (trade date), format: yyyy-MM-dd |
| `net_amount` | string | yes | Net change amount of the transaction (positive for credit, negative for debit) |
| `biz_time` | string | yes | Business event time when the transaction occurred |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

