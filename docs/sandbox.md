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

Webull publishes these shared accounts for immediate use — no application
required. They are public and shared across all developers.

| # | Account ID | App Key | App Secret |
|---|-----------|---------|------------|
| 1 | `V4H6R3L4VRI33UQ4TGR2NM1VI9` | `4b2b7acd2bf0d30d8aea173fceefa238` | `840b4353a6a31ce3ab91e2f99a510272` |
| 2 | `OGG4RRLC6EDE98HI920KRBVSKB` | `42bd186fb65ea76de309d69cf12f024e` | `29feb64b59d6b1b6b2d2aa8cea8a1b8d` |
| 3 | `2DHSQ9B1DMPBFPMPFU2R5SDPB8` | `64fc722617af8b5ebb746f50a910e91f` | `a268416fc681d438533f9e9316bab576` |

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

The SDK's integration tests are skipped unless explicitly enabled. They read:

| Variable | Purpose |
|----------|---------|
| `WEBULL_SANDBOX=1` | Enables the sandbox integration tests |
| `WEBULL_APP_KEY` | Sandbox app key |
| `WEBULL_APP_SECRET` | Sandbox app secret |
| `WEBULL_MQTT_WEBSOCKET=1` | Runs the MQTT-over-WebSocket streaming tests |

```sh
# macOS / Linux
WEBULL_SANDBOX=1 \
WEBULL_APP_KEY=your-sandbox-app-key \
WEBULL_APP_SECRET=your-sandbox-app-secret \
go test ./... -run Integration
```

```powershell
# Windows PowerShell
$env:WEBULL_SANDBOX = "1"
$env:WEBULL_APP_KEY = "your-sandbox-app-key"
$env:WEBULL_APP_SECRET = "your-sandbox-app-secret"
go test ./... -run Integration
```

## Known sandbox limitations

Sandbox data is intentionally constrained. Expect the following while
developing:

- **Single symbol.** Market-data access is limited to `AAPL`.
- **Footprint entitlement.** Footprint requests return `403 Insufficient
  permission` because the sandbox account lacks the paid entitlement.
- **Options.** Option contracts for `AAPL` may not exist in the sandbox;
  option requests can return `417 Invalid Symbol`.
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
