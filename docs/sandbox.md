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

Webull publishes shared sandbox test accounts in its
[getting-started guide](https://developer.webull.com/apis/docs/getting-started).
Use them to try the API without applying for access.

Credentials are not reproduced in this repository and must never be committed.
Supply the sandbox app key and secret at run time:

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

See [Troubleshooting](troubleshooting.md) for the symptoms of each limitation and
how to respond.

## Related

- [Getting Started](getting-started.md) — credentials and the first call.
- [Authentication](authentication.md) — token lifecycle.
- [Streaming](streaming.md) — MQTT and WebSocket options.
