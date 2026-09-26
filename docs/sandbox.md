# Sandbox

The Webull sandbox is a separate deployment you can use for development and
integration testing without touching production data. Select it with
`client.WithSandbox()` (equivalent to `client.WithEnvironment(client.Sandbox)`) or
by setting `WEBULL_ENVIRONMENT=sandbox` / `WEBULL_ENVIRONMENT=uat`.

## Environments

For the Hong Kong region:

| Environment | REST | MQTT | MQTT over WebSocket |
|-------------|------|------|---------------------|
| Production | `https://api.webull.hk` | `data-api.webull.hk:1883` | `wss://data-api.webull.hk:8883/mqtt` |
| Sandbox | `https://api.sandbox.webull.hk` | `data-api.sandbox.webull.hk:1883` | `wss://data-api.sandbox.webull.hk:8883/mqtt` |

Other regions derive their hosts from the region's root domain, with a
`sandbox.` label added below each service subdomain. Only the Hong Kong hosts are
published explicitly by Webull; verify other regions before production use.

## Test credentials

Webull publishes shared sandbox test accounts for both regions:

- **Hong Kong:** <https://developer.webull.hk/apis/docs/sdk#test-accounts>
- **US / International:** <https://developer.webull.com/apis/docs/sdk#test-accounts>

Use them to try the API without applying for access.

### Shared HK sandbox accounts

Webull publishes shared HK sandbox test accounts for immediate use — no
application required. Each entry provides an account ID, an app key, and an app
secret:

- **Hong Kong:** <https://developer.webull.hk/apis/docs/sdk#test-accounts>
- **US / International:** <https://developer.webull.com/apis/docs/sdk#test-accounts>

Read the current values from that page rather than copying them into this
repository. Webull can rotate or retire shared accounts at any time, and a copy
pasted here would silently go stale; fetching the values at run time keeps the
procedure below correct without a repository change.

The account table also appears in [Master Guides](webull-api/master-guides.md),
which is a verbatim snapshot of Webull's own published documentation and is not
maintained by this project. Values are not inlined in the hand-written pages.

!!! caution
    These are shared public accounts. Other users may place orders on them
    at any time. Use them for read-only exploration and testing only. For
    dedicated sandbox accounts, apply through the
    [Sandbox environment application](https://developer.webull.hk/apis/docs/authentication/TradingAPIApplication).

### Setting credentials

Private credentials must never be committed. Supply the sandbox app key and
secret at run time:

```sh
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
export WEBULL_ENVIRONMENT="sandbox"
```

Only the host `api.sandbox.webull.hk` belongs in committed files. App keys, app
secrets, and access tokens are per-account secrets.

## Sandbox integration tests

The SDK's integration tests are skipped unless their environment gates are
enabled. The actual test selector is `Sandbox`. Before any live run, use
`make test` for the credential-free root and nested-module suite; root-only
`go test ./...` does not traverse nested modules.

### Core, market-data, streaming, and events gates

| Variable | Purpose |
|----------|---------|
| `WEBULL_SANDBOX=1` | Enables the core, data, stream, and events sandbox tests |
| `WEBULL_APP_KEY` | Sandbox app key for those tests |
| `WEBULL_APP_SECRET` | Sandbox app secret for those tests |
| `WEBULL_MQTT_WEBSOCKET=1` | Additionally runs MQTT-over-WebSocket tests |
| `WEBULL_BASE_URL` | Optional override for another regional sandbox |
| `WEBULL_MQTT_URL` | Optional MQTT broker address override |

The trading-events read-only test may omit an account filter. Set
`WEBULL_TRADE_ACCOUNT_ID` to scope that test to one account. The mutating
order-event test uses these generic credentials and additionally requires
`WEBULL_TRADE_ACCOUNT_ID` and `WEBULL_TRADE_MUTATE=1`.

```sh
# macOS / Linux
WEBULL_SANDBOX=1 \
WEBULL_APP_KEY=your-sandbox-app-key \
WEBULL_APP_SECRET=your-sandbox-app-secret \
go test ./... -run Sandbox
```

```powershell
# Windows PowerShell
$env:WEBULL_SANDBOX = "1"
$env:WEBULL_APP_KEY = "your-sandbox-app-key"
$env:WEBULL_APP_SECRET = "your-sandbox-app-secret"
go test ./... -run Sandbox
```

### Trading-package gates

Trading-package sandbox tests intentionally use a separate credential set:

| Variable | Purpose |
|----------|---------|
| `WEBULL_TRADE_SANDBOX=1` | Enables the dedicated trading sandbox tests |
| `WEBULL_TRADE_APP_KEY` | Trading sandbox app key |
| `WEBULL_TRADE_APP_SECRET` | Trading sandbox app secret |
| `WEBULL_TRADE_ACCOUNT_ID` | Account whose assets, orders, and previews are exercised |
| `WEBULL_TRADE_MUTATE=1` | Additional opt-in for the test that places and cancels an order |
| `WEBULL_TRADE_PARTY_ID` | Optional HK BCAN party identifier for the gated preview |

```sh
WEBULL_TRADE_SANDBOX=1 \
WEBULL_TRADE_APP_KEY=your-trading-sandbox-app-key \
WEBULL_TRADE_APP_SECRET=your-trading-sandbox-app-secret \
WEBULL_TRADE_ACCOUNT_ID=your-trading-account-id \
go test ./trade -run Sandbox
```

Leave `WEBULL_TRADE_MUTATE` unset for the normal read-only and preview tests.
Set it only for an explicitly authorized sandbox mutation run; never point a
mutating test at production or a shared account.

## Known sandbox limitations

Sandbox data is intentionally constrained. Expect the following while
developing:

- **Single symbol.** Market-data access is limited to `AAPL`.
- **Footprint entitlement.** Footprint requests return `403 Insufficient
  permission` because the sandbox account lacks the paid entitlement.
- **Options.** Option contracts for `AAPL` may not exist in the sandbox;
  option requests can return `417 Invalid Symbol`. HTTP 417 is not uniquely a
  token failure—the SDK retains the historical `INVALID_TOKEN` category for
  compatibility, so inspect the API message before renewing a token.
- **Empty depth.** Order-book depth can be empty outside regular trading hours.
- **Network blocking.** Plain MQTT on port `1883` is blocked by some networks;
  use MQTT over WebSocket on `wss://...:8883/mqtt` (in `stream`, pass
  `WithWebSocket(true)`).
- **Token rate limit.** The token endpoint allows 10 requests per 30 seconds.
- **MQTT concurrency.** At most 5 concurrent MQTT connections per App Key; the
  server retains a disconnected session for about a minute.
- **Broker API HK.** The HK sandbox returns `401 ROUTE_NOT_PERMITTED` for Broker API
  HK endpoints (`/broker/...`). The app lacks the required scope, not a path issue.
- **SSE news.** The news SSE endpoint returns `504 Gateway Timeout` in the
  sandbox.

See [Troubleshooting](troubleshooting.md) for the symptoms of each limitation and
how to respond.

## Related

- [Getting Started](getting-started.md) — credentials and the first call.
- [Authentication](authentication.md) — token lifecycle.
- [Streaming](streaming.md) — MQTT and WebSocket options.
