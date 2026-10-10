# Examples

Runnable examples for the current API surface of `webullapi4go`. Examples are
split between the root module and standalone modules. From the repository root,
use the module-aware Makefile targets:

```sh
make build
make vet
```

Those targets cover the root module plus `broker/`,
`examples/watchlist-cmd/`, `examples/broker-probe/`,
`examples/futures-probe/`, and `examples/options-multi-leg/`. A root
`go build ./...` or `go vet ./...` does not traverse nested modules.

## Credentials

Every example constructs its client with
[`client.WithEnv()`](https://pkg.go.dev/github.com/shing1211/webullapi4go/client#WithEnv),
which reads:

| Variable | Purpose |
|----------|---------|
| `WEBULL_APP_KEY` | Webull OpenAPI app key (required) |
| `WEBULL_APP_SECRET` | Webull OpenAPI app secret (required) |
| `WEBULL_REGION` | Region, for example `hk` or `us` (optional, defaults to `hk`) |
| `WEBULL_ENVIRONMENT` | `sandbox` / `uat` or `prod` / `production` (optional, defaults to production) |
| `WEBULL_ACCOUNT_ID` | Trading account to inspect in `account` or monitor in `account-monitor` (optional) |
| `WEBULL_TRADE_ACCOUNT_ID` | Trading account to subscribe to in the `events` example (optional) |
| `WEBULL_MONITOR_INTERVAL` | `account-monitor` poll interval, such as `10s` (optional, defaults to `30s`) |
| `WEBULL_ORDER_IDEMPOTENCY_KEY` | Persisted `client_order_id` required for a mutating `order` run (optional, preview-only by default) |
| `WEBULL_ORDER_PLACE` | Example only: set to `1` in `order` to opt in to placing a real order (optional, off by default) |

Credentials are never hard-coded and must never be committed. Webull publishes
shared sandbox test accounts in its
[getting-started guide](https://developer.webull.hk/apis/docs/getting-started);
supply the values at run time through your shell or a secret manager. Only the
sandbox host `api.sandbox.webull.hk` belongs in committed material.

Most development happens against the sandbox:

```sh
# macOS / Linux
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
export WEBULL_ENVIRONMENT="sandbox"
```

```powershell
# Windows PowerShell
$env:WEBULL_APP_KEY = "your-sandbox-app-key"
$env:WEBULL_APP_SECRET = "your-sandbox-app-secret"
$env:WEBULL_ENVIRONMENT = "sandbox"
```

See [Sandbox](../docs/sandbox.md) in the documentation site for the full
environment table and known limitations.

## auth

Creates or reuses an access token and prints its status and expiry.

```sh
go run ./examples/auth
```

Use this first to confirm credentials work. In the sandbox the token is `NORMAL`
immediately; in production the call waits for the Webull App 2FA verification
window (up to five minutes by default).

## marketdata

Fetches a snapshot and recent daily bars for `AAPL` on the `US` market, then
prints the fields.

```sh
go run ./examples/marketdata
```

The sandbox only serves a limited symbol set (currently `AAPL`). Bars are
requested through the batch endpoint because the historical single-symbol
endpoint has been retired.

## streaming

Connects to the streaming broker over MQTT-over-WebSocket, subscribes to
`AAPL` `QUOTE`, `SNAPSHOT`, and `TICK` pushes, and consumes each typed channel.
Channels have bounded buffers with a drop-oldest policy, and the run stops after
two minutes or on Ctrl+C (or `SIGTERM`).

```sh
go run ./examples/streaming
```

WebSocket is used because plain MQTT on port `1883` is blocked on some networks;
the example targets `wss://data-api.sandbox.webull.hk:8883/mqtt` in the sandbox.
The stream also needs an access token, which the example obtains with
`EnsureToken`.

## watchlist

Lists the authenticated user's watchlists. This example is read-only.

```sh
go run ./examples/watchlist
```

## account

Lists the authenticated user's trading accounts and, for one account, prints its
balance and open positions. This example is read-only: it never places, replaces,
or cancels an order.

```sh
go run ./examples/account
```

The account is taken from `WEBULL_ACCOUNT_ID`; when it is unset, the first
account returned by the API is used. Trading requests require an access token,
which the example obtains with `EnsureToken`.

## account-monitor

Periodically prints one account's balance, buying power, and initial margin. A
single worker performs bounded, read-only `GetBalance` requests on a ticker;
ticks that arrive during a request are skipped so calls never overlap. Each
request has a 15-second timeout, and the monitor stops after three consecutive
failures or when Ctrl+C (or `SIGTERM`) cancels it.

```sh
go run ./examples/account-monitor

# Optional: poll every 10 seconds instead of the 30-second default.
WEBULL_MONITOR_INTERVAL=10s go run ./examples/account-monitor
```

`WEBULL_ACCOUNT_ID` selects the account; when unset, the first account returned
by the API is monitored.

## order

Previews a small AAPL limit buy and, only when `WEBULL_ORDER_PLACE=1` and a
persisted `WEBULL_ORDER_IDEMPOTENCY_KEY` are set, places it far below the market
and immediately cancels it. Previewing is read-only; placing mutates the
account, so the example is preview-only by default and never sends a market
order.

```sh
# Preview only (mutates nothing)
go run ./examples/order

# Choose and persist a new key before the first mutating request.
# Reuse this exact key for retries of the same logical order.
WEBULL_ORDER_IDEMPOTENCY_KEY=phase7-aapl-001 \
WEBULL_ORDER_PLACE=1 \
go run ./examples/order
```

```powershell
# Windows PowerShell
$env:WEBULL_ORDER_IDEMPOTENCY_KEY = "phase7-aapl-001"
$env:WEBULL_ORDER_PLACE = "1"
go run ./examples/order
```

Use it against a sandbox account only. The account is taken from
`WEBULL_ACCOUNT_ID`; when it is unset, the first account returned by the API is
used. Mutating placement refuses to start unless
`WEBULL_ORDER_IDEMPOTENCY_KEY` is a valid, caller-persisted key of at most 32
characters using letters, digits, `-`, or `_`; reuse it after an unknown result
and choose a new one for each intentionally new order. Preview-only runs use
`trade.ClientOrderIDFrom` for a deterministic identifier.

The example configures the `WithMaxOrderQuantity("10")` and
`WithMaxOrderNotional("2500.00")` guardrails, which `PreviewOrder` and
`PlaceOrder` enforce before any network call. After placement it uses an
independent 30-second cancellation timeout so Ctrl+C cannot skip cleanup.

## events

Connects to the Webull gRPC trade-event stream, subscribes to order events, and
prints each decoded `OrderEvent` until Ctrl+C. The event service signs each
`Subscribe` call with HMAC-SHA256, so no access token is required.

```sh
go run ./examples/events
```

The subscription is unfiltered when `WEBULL_TRADE_ACCOUNT_ID` is unset. When
it is set, the selected account must belong to the App Key. The stream
reconnects and re-subscribes automatically after a transient drop. In the
sandbox, placement events may not be pushed for a resting order; a cancellation
produces the observed `CANCEL_SUCCESS` event.

## data-fundamentals

Fetches all available fundamental data for `AAPL` on the `US` market:
company profile, analyst target and rating, capital flow, industry comparison,
earnings and dividend calendars, SEC filings, income statement, balance sheet,
cash flow, financial indicators, financial alert, and forecast EPS.

```sh
go run ./examples/data-fundamentals
```

## probe

Sandbox endpoint testing tool. Hits live sandbox endpoints to verify response
schemas for v0.5+ features (crypto, fund data, screener v2, display solution,
broker FD). Includes a Go program and a Python equivalent.

```sh
export WEBULL_ENVIRONMENT="sandbox"
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
go run ./examples/probe
```

Results are saved to `examples/probe/results/`. Not a user-facing example —
used during SDK development to verify endpoint paths and response shapes.

## watchlist-cmd

Watchlist CRUD example demonstrating create, add instruments, update, remove
instruments, and delete. Guarded by `WEBULL_WATCHLIST_TEST=1` since it mutates
the watchlist.

```sh
cd examples/watchlist-cmd
WEBULL_WATCHLIST_TEST=1 go run .
```

## broker-probe

Broker HK read-only endpoint probe program. Tests virtual account, instrument,
asset, order, cash activity, FX, and journal endpoints against the sandbox.

```sh
cd examples/broker-probe
export WEBULL_ENVIRONMENT="sandbox"
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
go run .
```

Note: Broker API HK (`/broker/...`) returns `401 ROUTE_NOT_PERMITTED` in
the HK sandbox — the app lacks the required scope, not a path issue.

## brokerfd

Broker FD US read-only endpoint probe. Tests accounts, orders, assets,
instruments, funding, activity, journals, master data, agreements, and
documents endpoints against the sandbox.

```sh
export WEBULL_ENVIRONMENT="sandbox"
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
go run ./examples/brokerfd
```

## brokerfd-events

Broker FD gRPC event subscription probe. Subscribes to order, option, and
position event streams over gRPC and prints decoded payloads until Ctrl+C.

```sh
export WEBULL_ENVIRONMENT="sandbox"
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
go run ./examples/brokerfd-events
```

## options

Preview-only multi-leg US options order: builds an AAPL vertical call spread
(long 220 call / short 230 call) and calls `PreviewOrder` to estimate its cost.
No order is ever placed.

```sh
go run ./examples/options
```

The HK sandbox accepts only `SINGLE` strategies (`417` for multi-leg) and may
not have option contracts for `AAPL` (`417 Invalid Symbol`), so a US sandbox is
needed to exercise this example end to end.

## options-multi-leg

Probes every multi-leg strategy value (VERTICAL, STRADDLE, STRANGLE,
IRON_CONDOR, IRON_BUTTERFLY, BUTTERFLY, CALENDAR, DIAGONAL, RATIO, COLLAR)
through the read-only `PreviewOrder`; no live orders are sent. Its own Go
module. Off by default:

```sh
cd examples/options-multi-leg
WEBULL_OPTIONS_TEST=1 go run .
```

The HK sandbox accepts only `SINGLE` (`417` otherwise).

## futures-probe

Probes HK futures discovery and market data (instrument list, product codes,
product classes, snapshot, bars). Its own Go module. Off by default:

```sh
cd examples/futures-probe
WEBULL_FUTURES_TEST=1 go run .
```

## path-probe

Compares the paths the SDK calls against the official Webull OpenAPI definition
for the endpoints where Webull's `llms.txt` summary and its own OpenAPI JSON
disagree (Known Issue 12 in `IMPLEMENTATION_STATUS.md`). Env-gated:

```sh
WEBULL_SANDBOX=1 WEBULL_APP_KEY=... WEBULL_APP_SECRET=... go run ./examples/path-probe
```

## Sandbox limitations

While trying these examples against the sandbox, expect a few restrictions:

- Market data is limited to `AAPL`.
- Footprint returns `403 Insufficient permission` without a paid entitlement.
- Option contracts for `AAPL` may not exist (`417 Invalid Symbol`).
- Order-book depth can be empty outside regular trading hours.
- Plain MQTT on `:1883` may be blocked; use MQTT over WebSocket on `:8883/mqtt`.
- The token endpoint allows 10 requests per 30 seconds, and MQTT allows at most
  5 concurrent connections per App Key.
- The gRPC event stream may not push a placement event for a resting order; only
  `CANCEL_SUCCESS` has been observed. If an account filter is configured, the
  selected account must belong to the App Key.

See [Troubleshooting](../docs/troubleshooting.md) for symptoms and fixes.
