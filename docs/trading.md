# Trading

The `trade` package exposes the Webull Trading HTTP API. It is a thin, typed
layer over the core [client](api.md#client), so request signing, access-token
handling, retries, rate limiting, and error classification are shared with the
rest of the SDK.

The v0.2.1 foundation covers read-only account and asset access: listing
accounts, reading an account's balance, and listing open positions. Order
placement, replacement, cancellation, and order queries arrive in later v0.2
patches.

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

## API version

Requests under `/trading/` default to the v3 API version, sent as the
`x-version: v3` header. Market Data paths still default to `v2`. The default can
be overridden globally with `client.WithAPIVersion` or per prefix with
`client.WithAPIVersionFor`; the longest matching prefix wins.

## Order guardrails

The client accepts advisory order guardrails. They are configuration in the
foundation release and are enforced by the order methods before an order is
built:

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

## Sandbox

Use the sandbox while developing. Set the environment and credentials, and
select an account with `WEBULL_ACCOUNT_ID`:

```sh
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
export WEBULL_ENVIRONMENT="sandbox"
export WEBULL_ACCOUNT_ID="your-sandbox-account-id"
```

The runnable [`examples/account`](https://github.com/shing1211/webullapi4go/tree/main/examples/account)
program lists accounts and prints the balance and positions for
`WEBULL_ACCOUNT_ID`, or for the first account when it is unset. It is read-only.

Credentials, account IDs, and tokens are per-account secrets and must never be
committed. Only the sandbox host `api.sandbox.webull.hk` belongs in committed
material. The sandbox integration test is gated by
`WEBULL_TRADE_SANDBOX=1` together with `WEBULL_TRADE_APP_KEY`,
`WEBULL_TRADE_APP_SECRET`, and `WEBULL_TRADE_ACCOUNT_ID`; see
[Sandbox](sandbox.md) for the environment table and known limitations.

## Related

- [Getting Started](getting-started.md) — install and credentials.
- [Authentication](authentication.md) — signing and tokens.
- [Sandbox](sandbox.md) — environments, test credentials, limitations.
- [`trade` reference](https://pkg.go.dev/github.com/shing1211/webullapi4go/trade) — full API.
