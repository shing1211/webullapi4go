# Troubleshooting

This page lists common symptoms, their causes, and fixes.

## Diagnose an error

Run the auth example first to confirm credentials and connectivity
independently of any endpoint:

```sh
go run ./examples/auth
```

The SDK error message carries a stable code (for example `webull: FORBIDDEN:
http 403: Insufficient permission`). See [Errors](errors.md) for the full code
table.

## HTTP errors

| Symptom | Likely cause | Fix |
|---------|--------------|-----|
| `401 UNAUTHORIZED` | Missing or badly signed request | Confirm `WEBULL_APP_KEY` / `WEBULL_APP_SECRET`; let `Client.Do` sign the request |
| `403 FORBIDDEN` / `Insufficient permission` | Endpoint requires a paid entitlement (for example footprint) | Use an account with the entitlement; on the sandbox this is expected |
| `417 INVALID_TOKEN` | Access token missing, expired, or invalid | Call `Client.EnsureToken` and retry |
| `417 Invalid Symbol` | Symbol not carried by the environment | In the sandbox, use `AAPL` |
| `429 RATE_LIMITED` | Too many requests | Back off; the token endpoint allows 10 requests per 30 seconds |
| `5xx SERVER_ERROR` | Webull service failure | Transient; retry, or configure `WithRetry` / `WithBreaker` |

Empty order-book depth is not an error: markets push no depth outside regular
trading hours, so `GetQuotes` or a `QUOTE` subscription may simply have nothing to
return.

## Sandbox-specific limitations

- Market-data access is limited to the symbol `AAPL`.
- Footprint returns `403 Insufficient permission` without a paid entitlement.
- Option contracts for `AAPL` may not exist and can return `417 Invalid Symbol`.
- Order-book depth can be empty outside regular trading hours.
- Plain MQTT on `:1883` can be blocked by some networks.
- The token endpoint allows 10 requests per 30 seconds.
- MQTT allows at most 5 concurrent connections per App Key.

See [Sandbox](sandbox.md) for the environment table and test-credential setup.

## Streaming connection problems

### Connection limit (Webull code 105)

An App Key supports at most **5 concurrent MQTT connections**. Exceeding the
limit fails the connection with Webull error code 105. `stream.Connect` returns
an error explaining the limit. The server retains a disconnected session for
about one minute, so wait roughly a minute before reconnecting rather than
retrying immediately.

### Plain MQTT is blocked

If `stream.New(cl)` without `WithWebSocket` cannot connect, port `1883` may be
blocked on your network. Use MQTT-over-WebSocket instead:

```go
s, err := stream.New(cl, stream.WithWebSocket(true))
```

That targets `wss://data-api.sandbox.webull.hk:8883/mqtt` in the sandbox.

### Session id reuse

Reusing a session id across connections disconnects the previous connection.
`stream.New` generates a unique id by default; only pass `WithSessionID` when
exclusivity is guaranteed.

### Subscriptions stop after a reconnect

Subscriptions are restored automatically only when `WithAutoReconnect` is
enabled. It is off by default. Webull does not restore subscriptions itself.

## Production token stays PENDING

Production tokens start `PENDING` and become `NORMAL` only after the Webull App
2FA flow completes within five minutes. `Client.EnsureToken` polls until then and
returns an `auth` error if the token becomes `INVALID` or `EXPIRED` or the poll
timeout elapses. Increase the window with `Client.SetTokenPollTimeout` if needed.

## Docs build

To build the documentation site locally:

```sh
python -m venv .venv
.venv/Scripts/pip install -r requirements-docs.txt   # Windows
pip install -r requirements-docs.txt                 # macOS / Linux
mkdocs build --strict
```

Internal run artifacts under `docs/runs/` are workspace-only and excluded from
the built site.

## Getting help

- [GitHub Discussions](https://github.com/shing1211/webullapi4go/discussions) for
  questions and ideas.
- [Issues](https://github.com/shing1211/webullapi4go/issues) for bugs and feature
  requests.
- For vulnerabilities, follow [SECURITY.md](https://github.com/shing1211/webullapi4go/blob/main/SECURITY.md).
