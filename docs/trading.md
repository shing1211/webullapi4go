# Trading

The `trade` package exposes the Webull Trading HTTP API. It is a thin, typed
layer over the core [client](api.md#client), so request signing, access-token
handling, retries, rate limiting, and error classification are shared with the
rest of the SDK.

The v0.2.1 foundation covers read-only account and asset access. The v0.2.2
release adds the stock-order lifecycle — preview, place, replace, cancel — and
order queries. The v0.2.3 release adds the per-market order rules. The v0.2.4
release adds single-leg options orders, and the v0.2.5 release adds US combo
orders.

## Authentication

Trading requests require an access token, sent as the `x-access-token` header.
Call `EnsureToken` once before the first request; it creates or reuses a token
and the client then injects it into every request.

```go
cl, err := client.New(client.WithEnv())
if err != nil {
	return err
}
defer func() { _ = cl.Close() }()

if _, err := cl.EnsureToken(ctx); err != nil {
	return err
}

trading := trade.New(cl)
defer func() { _ = trading.Close() }()
```

See [Authentication](authentication.md) for the token lifecycle.

## Accounts and assets

| Method | Endpoint | Returns |
|--------|----------|---------|
| `ListAccounts(ctx)` | `GET /trading/accounts/list` | `[]Account` |
| `GetBalance(ctx, accountID)` | `GET /trading/assets/balances/get` | `*AssetsBalance` |
| `GetPositions(ctx, accountID)` | `GET /trading/assets/positions/list` | `[]Position` |

`ListAccounts` returns the accounts available to the authenticated user. Every
other asset call is scoped to one account: `GetBalance` and `GetPositions`
require the `accountID` from `Account.AccountID`. An empty account ID is
rejected with an `invalid_config` error before any network call, so a request is
never sent without an account.

```go
accounts, err := trading.ListAccounts(ctx)
if err != nil {
	return err
}
if len(accounts) == 0 {
	return errors.New("no trading accounts")
}

accountID := accounts[0].AccountID

balance, err := trading.GetBalance(ctx, accountID)
if err != nil {
	return err
}
log.Printf("cash=%s market_value=%s pnl=%s",
	balance.TotalCashBalance, balance.TotalMarketValue,
	balance.TotalUnrealizedProfitLoss)

positions, err := trading.GetPositions(ctx, accountID)
if err != nil {
	return err
}
for _, pos := range positions {
	log.Printf("%s qty=%s last=%s", pos.Symbol, pos.Quantity, pos.LastPrice)
}
```

`AssetsBalance` includes a per-currency breakdown in
`AccountCurrencyAssets`. `Position` carries the held quantity, average cost,
last price, unrealized P/L, and, for multi-leg option positions, the `Legs`
slice. Numeric fields such as quantities and prices are strings, preserving the
precision of the wire values.

## Order lifecycle

!!! warning "Placing an order mutates the account"

    `PlaceOrder`, `ReplaceOrder`, and `CancelOrder` change the state of a real
    brokerage account. The sandbox is still a real trading account: an order
    placed there is an order. Always preview first, use the guardrails, keep
    orders small and non-marketable, and gate any placement behind an explicit
    opt-in such as `WEBULL_ORDER_PLACE=1`. Never send a market order from an
    automated test or example.

A stock order moves through five operations:

1. **Preview** (`PreviewOrder`) — validate the request and estimate its cost.
   Read-only.
2. **Place** (`PlaceOrder`) — submit the order and receive its identifiers.
3. **Replace** (`ReplaceOrder`) — change the working terms of an order, matched
   by client order ID.
4. **Cancel** (`CancelOrder`) — cancel a working order, matched by client order
   ID.
5. **Query** (`GetOpenOrders`, `GetOrderHistory`, `GetOrderDetail`) — inspect
   working orders, historical orders, and a single order's detail.

| Method | Endpoint | Returns |
|--------|----------|---------|
| `PreviewOrder(ctx, req)` | `POST /trading/orders/preview` | `*PreviewResult` |
| `PlaceOrder(ctx, req)` | `POST /trading/orders/place` | `*PlaceOrderResult` |
| `ReplaceOrder(ctx, req)` | `POST /trading/orders/replace` | `*ReplaceOrderResult` |
| `CancelOrder(ctx, req)` | `POST /trading/orders/cancel` | `*CancelOrderResult` |
| `GetOpenOrders(ctx, accountID)` | `GET /trading/orders/open-orders/list` | `[]OrderGroup` |
| `GetOpenOrdersPage(ctx, accountID, key)` | `GET /trading/orders/open-orders/list` | `*OrderPage` |
| `GetAllOpenOrders(ctx, accountID)` | `GET /trading/orders/open-orders/list` | `[]OrderGroup` |
| `GetOrderHistory(ctx, q)` | `GET /trading/orders/historical-orders/list` | `[]OrderGroup` |
| `GetOrderHistoryPage(ctx, q)` | `GET /trading/orders/historical-orders/list` | `*OrderPage` |
| `GetAllOrderHistory(ctx, q)` | `GET /trading/orders/historical-orders/list` | `[]OrderGroup` |
| `GetOrderDetail(ctx, accountID, clientOrderID)` | `GET /trading/orders/get` | `*OrderGroup` |

The order endpoints take a `PlaceOrderRequest`, which carries the `AccountID`
and one or more `OrderRequest` values. The modify endpoints take a
`ReplaceOrderRequest` or `CancelOrderRequest`. Requests are validated before any
network call, so a malformed request never reaches the API.

The v3 modify and query endpoints identify an order by its **client order ID**,
not the system order ID: reuse the `ClientOrderID` supplied at placement time.
`PlaceOrderResult.OrderID` is informational.

```go
req := trade.PlaceOrderRequest{
	AccountID: accountID,
	NewOrders: []trade.OrderRequest{
		{
			ClientOrderID:         "demo-aapl-buy-1",
			ComboType:             trade.ComboTypeNormal,
			InstrumentType:        trade.InstrumentTypeEquity,
			Market:                trade.MarketUS,
			Symbol:                "AAPL",
			OrderType:             trade.OrderTypeLimit,
			Side:                  trade.OrderSideBuy,
			Quantity:              "1",
			EntrustType:           trade.EntrustTypeQty,
			TimeInForce:           trade.TimeInForceDay,
			SupportTradingSession: trade.TradingSessionCore,
			LimitPrice:            "1.00", // far below market, so it will not fill
		},
	},
}

preview, err := trading.PreviewOrder(ctx, req)
if err != nil {
	return err
}
log.Printf("estimated cost=%s fee=%s",
	preview.EstimatedCost, preview.EstimatedTransactionFee)

placed, err := trading.PlaceOrder(ctx, req)
if err != nil {
	return err
}
log.Printf("placed order_id=%s", placed.OrderID)

cancelled, err := trading.CancelOrder(ctx, trade.CancelOrderRequest{
	AccountID:     accountID,
	ClientOrderID: placed.ClientOrderID,
})
if err != nil {
	return err
}
log.Printf("cancelled order_id=%s", cancelled.OrderID)
```

To change a working order instead of cancelling it:

```go
_, err := trading.ReplaceOrder(ctx, trade.ReplaceOrderRequest{
	AccountID: accountID,
	ModifyOrders: []trade.ModifyOrderRequest{
		{
			ClientOrderID: "demo-aapl-buy-1",
			LimitPrice:    "1.50",
		},
	},
})
```

Only the fields set on each `ModifyOrderRequest` are changed. `ClientOrderID` is
required and selects the order; `TimeInForce`, `Quantity`, `LimitPrice`,
`StopPrice`, `TriggerPriceType`, `TrailingType`, `TrailingStopStep`,
`TrailingLimitPriceOffset`, and `ExpireDate` are optional.

## Order types

`OrderRequest.OrderType` is the execution instruction. The SDK accepts the stock
order types below; the required price fields are the ones `Validate` checks
before a request is sent. The per-market validity matrix (which type is allowed
for a given market and instrument) is also enforced locally; see
[Market rules](#market-rules).

| Order type | Extra fields the SDK requires | Notes |
|------------|-------------------------------|-------|
| `LIMIT` | `limit_price` | Executes at the limit price or better |
| `MARKET` | none | Executes at the best available price. Never place one from a test or example |
| `STOP_LOSS` | `stop_price` | Becomes a market order once the stop triggers |
| `STOP_LOSS_LIMIT` | `stop_price`, `limit_price` | Becomes a limit order once the stop triggers |
| `TOUCH_MKT` | `stop_price` | Becomes a market order when the trigger price is touched |
| `TOUCH_LMT` | `stop_price`, `limit_price` | Becomes a limit order when the trigger price is touched |
| `TRAILING_STOP_LOSS` | `trailing_type`, `trailing_stop_step` | Trails the market by a fixed amount or percentage |
| `TRAILING_STOP_LOSS_LIMIT` | `trailing_type`, `trailing_stop_step` | Trailing stop that submits a limit order |
| `ENHANCED_LIMIT` | `limit_price` | Hong Kong enhanced limit order |
| `AT_AUCTION` | none, and `limit_price` must be empty | Hong Kong at-auction order |
| `AT_AUCTION_LIMIT` | `limit_price` | Hong Kong at-auction limit order |
| `MARKET_ON_OPEN` | none | Executes at the opening price |
| `MARKET_ON_CLOSE` | none | Executes at the closing price |

`TriggerPriceType` selects the market price a touch or stop order triggers on:
`PRICE` (last trade), `PRICE_BID` (best bid), or `PRICE_ASK` (best ask).
`TrailingType` is `AMOUNT` for a fixed price spread or `PERCENTAGE` for a
percentage where `"0.01"` is 1%.

## Market rules

`OrderRequest.Validate` enforces the market-specific rules below before any
network call. A violation is an `invalid_config` error that names the offending
field.

### Equity order types by market

The accepted equity order types depend on `OrderRequest.Market`:

| Market | Accepted order types |
|--------|----------------------|
| `US` | `LIMIT`, `MARKET`, `STOP_LOSS`, `STOP_LOSS_LIMIT`, `MARKET_ON_OPEN`, `MARKET_ON_CLOSE`, `TOUCH_MKT`, `TOUCH_LMT`, `TRAILING_STOP_LOSS`, `TRAILING_STOP_LOSS_LIMIT` |
| `HK` | `ENHANCED_LIMIT`, `AT_AUCTION`, `AT_AUCTION_LIMIT`, `STOP_LOSS`, `STOP_LOSS_LIMIT`, `TOUCH_MKT`, `TOUCH_LMT`, `TRAILING_STOP_LOSS`, `TRAILING_STOP_LOSS_LIMIT` |
| `CN` | `LIMIT` |

An order type outside its market's row is rejected. The matrix applies to equity
orders; option and futures order types are validated by their own rules.

### Hong Kong BCAN party IDs

Every Hong Kong equity order must carry at least one `no_party_ids` entry. Each
entry is a `trade.PartyID`, and all three fields are required:

| Field | Required value |
|-------|----------------|
| `party_id` | The broker client identifier (BCAN), for example `ABC123.2568`; non-blank |
| `party_id_source` | `"D"` |
| `party_role` | `"3"` |

`no_party_ids` is rejected on a non-HK-equity order: a US or CN order, or a
non-equity HK order, must leave it empty.

```go
NoPartyIDs: []trade.PartyID{{
	PartyID:       "ABC123.2568",
	PartyIDSource: "D",
	PartyRole:     "3",
}},
```

The BCAN is per-account and must never be committed. Read it from the
environment, as the sandbox test does with `WEBULL_TRADE_PARTY_ID`.

### US trading sessions

`support_trading_session` is optional and valid only on US orders; setting it on
a non-US order is rejected. The accepted values are:

| Value | Meaning |
|-------|---------|
| `CORE` | Regular trading hours |
| `ALL` | Extended hours |
| `NIGHT` | Night trading |
| `ALL_DAY` | Overnight, 8:00 p.m. ET to 8:00 p.m. ET the next day |

The deprecated aliases `Y` (use `ALL`) and `N` (use `CORE`) are rejected.

### At-auction price rules

Hong Kong at-auction orders differ in whether they carry a price:

- `AT_AUCTION` must **not** set `limit_price`; a limit price is rejected.
- `AT_AUCTION_LIMIT` **must** set `limit_price`.

### A-share orders are disabled by default

`CN` (A-share Stock Connect) accepts only `LIMIT`, and A-share trading is
disabled on the account by default: it must be enabled by Webull support before
any A-share order is accepted. The SDK validates the request shape but cannot
enable the entitlement.

## Options orders

An options order is a single-leg `SINGLE` order: `instrument_type` is `OPTION`,
`option_strategy` is `SINGLE`, and `legs` holds exactly one leg. The
`OrderRequest` still carries the top-level `symbol`, `side`, `order_type`, and
`quantity`, and the leg repeats the contract fields the API attributes the fill
to. `OrderRequest.Validate` rejects an option order that carries any other
strategy, no leg, or more than one leg.

Only the following order types are accepted for options:

| Order type | Extra fields the SDK requires |
|------------|-------------------------------|
| `LIMIT` | `limit_price` |
| `STOP_LOSS` | `stop_price` |
| `STOP_LOSS_LIMIT` | `stop_price`, `limit_price` |

`MARKET` and every other stock order type are rejected with an `invalid_config`
error before any network call.

`side` must be `BUY` or `SELL`; `SHORT` is rejected. A sell-side option order
must use `time_in_force` `DAY`: the API does not accept `GTC` or `GTD` when
closing an options position. A buy-side order may use `DAY` or `GTC`, but `GTD`
is rejected for options in both directions.

Each `OrderLeg` in `legs` carries:

| Field | Required value |
|-------|----------------|
| `instrument_type` | `OPTION` |
| `market` | `US` |
| `symbol` | The option contract symbol; non-blank |
| `side` | `BUY` or `SELL` |
| `strike_price` | The strike, a positive decimal string |
| `option_expire_date` | Expiration date in `YYYY-MM-DD` form |
| `option_type` | `CALL` or `PUT` |
| `quantity` | The leg quantity, a positive decimal string |

The example below previews a non-marketable single-leg AAPL call; placing it
follows the same pattern as a stock order and mutates the account.

```go
option := trade.OrderRequest{
	ClientOrderID:  "demo-aapl-call-1",
	ComboType:      trade.ComboTypeNormal,
	InstrumentType: trade.InstrumentTypeOption,
	Market:         trade.MarketUS,
	Symbol:         "AAPL",
	OrderType:      trade.OrderTypeLimit,
	Side:           trade.OrderSideBuy,
	Quantity:       "1",
	EntrustType:    trade.EntrustTypeQty,
	TimeInForce:    trade.TimeInForceDay,
	LimitPrice:     "0.05", // far below market, so it will not fill
	OptionStrategy: trade.OptionStrategySingle,
	Legs: []trade.OrderLeg{{
		InstrumentType:   trade.InstrumentTypeOption,
		Market:           trade.MarketUS,
		Symbol:           "AAPL",
		Side:             trade.OrderSideBuy,
		StrikePrice:      "100.00",
		OptionExpireDate: "2026-01-16",
		OptionType:       trade.OptionTypeCall,
		Quantity:         "1",
	}},
}

preview, err := trading.PreviewOrder(ctx, trade.PlaceOrderRequest{
	AccountID: accountID,
	NewOrders: []trade.OrderRequest{option},
})
if err != nil {
	return err
}
log.Printf("estimated cost=%s", preview.EstimatedCost)
```

!!! note "Sandbox option-contract availability"

    Sandbox market data is limited to `AAPL`, and the sandbox may not list the
    option contract you have in mind: a preview can fail with `417 Invalid
    Symbol` when the contract does not exist. The expiration date above is
    illustrative; pick a real listed contract for the account and environment.
    Footprint and other entitlement-gated data may also return `403
    Insufficient permission` in the sandbox.

## Time in force

| Value | Meaning | Notes |
|-------|---------|-------|
| `DAY` | Expires at the end of the trading day | Default choice for a non-marketable test order |
| `GTC` | Remains active until filled or cancelled | |
| `GTD` | Expires on `expire_date` | `expire_date` (yyyy-MM-dd) is required; currently US only |

## Combo types

`OrderRequest.ComboType` identifies the role an order plays within a combo
order. A plain stock order is `NORMAL`.

| Value | Role |
|-------|------|
| `NORMAL` | A standard single order |
| `MASTER` | The primary order that triggers its siblings |
| `STOP_PROFIT` | A take-profit sub-order |
| `STOP_LOSS` | A stop-loss sub-order |
| `OTO` | The follow-up order of a one-triggers-the-other pair |
| `OCO` | One of a pair where filling either cancels the other |
| `OTOCO` | One of the order set triggered by an OTOCO master |

### Combo orders (US only)

Combo orders are supported for US equity orders only. When any order in a
`PlaceOrderRequest` uses a non-NORMAL `combo_type`, the whole request must form
one valid combo group, and `client_combo_order_id` must be set to group it. A
request whose orders are all `NORMAL` is unaffected.

`PlaceOrderRequest.Validate` enforces the following composition rules before any
network call, and returns an `invalid_config` error naming the first problem:

- Every order in a combo request uses a non-NORMAL `combo_type`; a `NORMAL`
  order cannot be mixed with combo orders.
- Every order is a US equity order (`instrument_type` `EQUITY`, `market` `US`).
- `client_combo_order_id` is non-empty.
- The orders form exactly one group kind; roles from different groups cannot be
  mixed.

| Group | Composition | Leg order types |
|-------|-------------|-----------------|
| Take-profit/stop-loss (buy-to-open) | Exactly one `MASTER`, with an optional `STOP_PROFIT` and an optional `STOP_LOSS` | `MASTER`: `LIMIT`, `MARKET`; `STOP_PROFIT`: `LIMIT`; `STOP_LOSS`: `STOP_LOSS` |
| Take-profit/stop-loss (sell-to-close) | One `STOP_PROFIT` and/or one `STOP_LOSS`, side `SELL`, no `MASTER` | `STOP_PROFIT`: `LIMIT`; `STOP_LOSS`: `STOP_LOSS` |
| `OTO` | Exactly one `MASTER` plus 1–6 `OTO` orders | `MASTER` and `OTO`: `LIMIT`, `MARKET`, `STOP_LOSS`, `STOP_LOSS_LIMIT` |
| `OCO` | 2–6 `OCO` orders, no `MASTER` | `OCO`: `LIMIT`, `STOP_LOSS`, `STOP_LOSS_LIMIT` |
| `OTOCO` | Exactly one `MASTER` plus 1–6 `OTOCO` orders | `MASTER`: `LIMIT`, `MARKET`, `STOP_LOSS`, `STOP_LOSS_LIMIT`; `OTOCO`: `LIMIT`, `STOP_LOSS`, `STOP_LOSS_LIMIT` |

A take-profit/stop-loss group with no `MASTER` is the sell-to-close form, used
to attach a profit target and a stop to an existing long position: every
sub-order uses side `SELL`, and the two sub-orders are alternatives that cancel
each other.

The example previews a buy-to-open take-profit/stop-loss group on `AAPL`: a
limit `MASTER` far below the market so it will not fill, a `STOP_PROFIT` above
it, and a `STOP_LOSS` below it. Placing it follows the same pattern as a stock
order and mutates the account.

```go
combo := trade.PlaceOrderRequest{
	AccountID:          accountID,
	ClientComboOrderID: "demo-aapl-tpsl",
	NewOrders: []trade.OrderRequest{
		{
			ClientOrderID:  "demo-aapl-tpsl-master",
			ComboType:      trade.ComboTypeMaster,
			InstrumentType: trade.InstrumentTypeEquity,
			Market:         trade.MarketUS,
			Symbol:         "AAPL",
			OrderType:      trade.OrderTypeLimit,
			Side:           trade.OrderSideBuy,
			Quantity:       "1",
			EntrustType:    trade.EntrustTypeQty,
			TimeInForce:    trade.TimeInForceDay,
			LimitPrice:     "1.00", // far below market, so it will not fill
		},
		{
			ClientOrderID:  "demo-aapl-tpsl-profit",
			ComboType:      trade.ComboTypeStopProfit,
			InstrumentType: trade.InstrumentTypeEquity,
			Market:         trade.MarketUS,
			Symbol:         "AAPL",
			OrderType:      trade.OrderTypeLimit,
			Side:           trade.OrderSideSell,
			Quantity:       "1",
			EntrustType:    trade.EntrustTypeQty,
			TimeInForce:    trade.TimeInForceDay,
			LimitPrice:     "999.00",
		},
		{
			ClientOrderID:  "demo-aapl-tpsl-loss",
			ComboType:      trade.ComboTypeStopLoss,
			InstrumentType: trade.InstrumentTypeEquity,
			Market:         trade.MarketUS,
			Symbol:         "AAPL",
			OrderType:      trade.OrderTypeStopLoss,
			Side:           trade.OrderSideSell,
			Quantity:       "1",
			EntrustType:    trade.EntrustTypeQty,
			TimeInForce:    trade.TimeInForceDay,
			StopPrice:      "1.00",
		},
	},
}

preview, err := trading.PreviewOrder(ctx, combo)
if err != nil {
	return err
}
log.Printf("estimated cost=%s", preview.EstimatedCost)
```

## Validation rules

Every order method validates its request and returns a typed error with code
`invalid_config` on the first problem, before any network call. The rules are:

- **`PlaceOrderRequest`** — `account_id` is required; `new_orders` must contain
  at least one order; no two orders may reuse a `client_order_id`.
- **`OrderRequest`** — `client_order_id`, `combo_type`, `instrument_type`,
  `market`, `symbol`, `order_type`, `side`, `entrust_type`, and `time_in_force`
  are required. `instrument_type` is `EQUITY`, `OPTION`, or `FUTURES` and
  `market` is `US`, `HK`, or `CN`.
- **Market rules** — the equity order-type matrix, the Hong Kong BCAN
  `no_party_ids` requirement, the US-only `support_trading_session`, and the
  at-auction price rules are enforced per market; see
  [Market rules](#market-rules).
- **Option rules** — an `OPTION` order must use `option_strategy` `SINGLE` with
  exactly one leg, an allowed order type, and a `BUY`/`SELL` side (sell-side
  only `DAY`); `option_strategy` and `legs` are rejected on a non-option order.
  See [Options orders](#options-orders).
- **Combo rules** — when any order uses a non-NORMAL `combo_type`, the request
  must be one valid US-equity combo group with a `client_combo_order_id`; see
  [Combo orders (US only)](#combo-orders-us-only).
- **`client_order_id`** — 1 to 32 characters from `[A-Za-z0-9_-]`, and unique
  per account. Generate a fresh identifier for each new order.
- **Size** — when `entrust_type` is `QTY`, `quantity` is required and must be a
  positive decimal; when it is `AMOUNT`, `total_cash_amount` is required and
  must be a positive decimal. Sizes are strings so precision is preserved;
  `AMOUNT` supports US fractional share trading.
- **Prices and expiry** — the conditional fields in the order-type table above
  are required; `expire_date` is required when `time_in_force` is `GTD`.
- **`ReplaceOrderRequest`** — `account_id` is required; `modify_orders` must
  contain at least one order; each `client_order_id` must be valid and unique
  within the request.
- **`CancelOrderRequest`** — `account_id` and a valid `client_order_id` are
  required.
- **`OrderHistoryQuery`** — `account_id` is required; `start_time` and
  `end_time` must be RFC3339 timestamps; the API allows at most six months of
  look-back and defaults to the last seven days. Following the pagination
  cursor is bounded by `trade.MaxOrderQueryPages` (100).

## Order guardrails

The client accepts advisory order guardrails. They are enforced by
`PreviewOrder` and `PlaceOrder` before an order is built, so an over-limit order
never leaves the process:

```go
trading := trade.New(cl,
	trade.WithMaxOrderNotional("2500.00"),
	trade.WithMaxOrderQuantity("100"),
)
```

`WithMaxOrderNotional` caps a single order's notional value and
`WithMaxOrderQuantity` caps its quantity. Both take a non-negative decimal
string; an empty string disables the cap (the default). Passing a value that is
not a non-negative, finite number panics at configuration time.

The notional cap compares `total_cash_amount` for `AMOUNT` orders and
`quantity` times `limit_price` for orders that carry both. When the notional
cannot be computed — for example a `MARKET` order with no limit price — the
notional cap is skipped, but the quantity cap still applies. The guardrails
apply to preview and place; `ReplaceOrder` is not bounded by them, so re-check
modified sizes yourself.

## Order queries

The query methods return `OrderGroup` values: a client order together with its
child orders. A `NORMAL` order has a single entry in `OrderGroup.Orders`; combo
orders carry their legs there. Each `Order` reports its `Status` (`PENDING`,
`SUBMITTED`, `CANCELLED`, `FILLED`, `FAILED`, or `PARTIAL_FILLED`), quantities
and prices as decimal strings, and timestamps in ISO8601 UTC form.

The list endpoints are cursor paginated. The plain getters return the first page;
the `Page` variants return one page with its `PaginationKey`; the `All` variants
follow the cursor to exhaustion.

```go
open, err := trading.GetOpenOrders(ctx, accountID)
if err != nil {
	return err
}
for _, group := range open {
	for _, order := range group.Orders {
		log.Printf("open %s %s status=%s filled=%s/%s",
			order.Symbol, order.Side, order.Status,
			order.FilledQuantity, order.TotalQuantity)
	}
}

history, err := trading.GetOrderHistory(ctx, trade.OrderHistoryQuery{
	AccountID: accountID,
	StartTime: "2025-01-05T22:59:59.012Z", // optional; defaults to last 7 days
	EndTime:   "2025-01-06T22:59:59.012Z", // optional; defaults to now
})
if err != nil {
	return err
}
_ = history

detail, err := trading.GetOrderDetail(ctx, accountID, "demo-aapl-buy-1")
if err != nil {
	return err
}
if detail.Orders[0].Commission != nil {
	log.Printf("commission=%s", detail.Orders[0].Commission.ActualCommission)
}
```

`GetOrderDetail` is the only query that returns the `Commission` and `Fees`
breakdowns.

## API version

Requests under `/trading/` default to the v3 API version, sent as the
`x-version: v3` header. Market Data paths still default to `v2`. The default can
be overridden globally with `client.WithAPIVersion` or per prefix with
`client.WithAPIVersionFor`; the longest matching prefix wins.

## Sandbox

Use the sandbox while developing. Set the environment and credentials, and
select an account with `WEBULL_ACCOUNT_ID`:

```sh
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
export WEBULL_ENVIRONMENT="sandbox"
export WEBULL_ACCOUNT_ID="your-sandbox-account-id"
```

The runnable
[`examples/account`](https://github.com/shing1211/webullapi4go/tree/main/examples/account)
program lists accounts and prints the balance and positions for
`WEBULL_ACCOUNT_ID`, or for the first account when it is unset. It is read-only.

The
[`examples/order`](https://github.com/shing1211/webullapi4go/tree/main/examples/order)
program previews a small AAPL limit buy and, only when `WEBULL_ORDER_PLACE=1` is
set, places it far below the market and immediately cancels it. Without that
variable it is preview-only and mutates nothing. The example never submits a
market order:

```sh
WEBULL_ORDER_PLACE=1 go run ./examples/order
```

Credentials, account IDs, and tokens are per-account secrets and must never be
committed. Only the sandbox host `api.sandbox.webull.hk` belongs in committed
material. The sandbox integration test is gated by
`WEBULL_TRADE_SANDBOX=1` together with `WEBULL_TRADE_APP_KEY`,
`WEBULL_TRADE_APP_SECRET`, and `WEBULL_TRADE_ACCOUNT_ID`; the mutating
place/cancel test additionally requires `WEBULL_TRADE_MUTATE=1`. See
[Sandbox](sandbox.md) for the environment table and known limitations.

## Related

- [Getting Started](getting-started.md) — install and credentials.
- [Authentication](authentication.md) — signing and tokens.
- [Sandbox](sandbox.md) — environments, test credentials, limitations.
- [Trading Events](events.md) — real-time order, position, and option status over gRPC; the counterpart to the trading HTTP API.
- [`trade` reference](https://pkg.go.dev/github.com/shing1211/webullapi4go/trade) — full API.
