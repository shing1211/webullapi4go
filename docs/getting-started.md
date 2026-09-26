# Getting Started

This page covers installation, credentials, the sandbox environment, and your
first call.

## Install

The current hardening is recorded in repository tag `v2.1.4`. The module path
remains `github.com/shing1211/webullapi4go` and stays on the v1 import path **by
decision**; no `/v2` migration is planned. Because the import path carries no
major-version suffix, the module proxy serves only the `v1.x` line, so an
unqualified `go get` installs `v1.1.1` — which predates the tagged hardening —
and the `v2.1.4` tag is not a published module version. To install the current
tree, pin a commit:

```sh
go get github.com/shing1211/webullapi4go@78c164c
```

Go resolves that to a pseudo-version (`v1.1.2-0.20260926035012-78c164c22fc5` at
the time of writing). To pin the released v1.x line instead:

```sh
go get github.com/shing1211/webullapi4go@v1.1.1
```

A commit pin moves only when you change it, whereas a branch reference such as
`@main` tracks the newest code. The module requires Go 1.26 or newer and has no
cgo dependencies.

## Credentials

The SDK reads credentials from the environment so secrets are never hard-coded.
`client.WithEnv()` reads:

| Variable | Purpose |
|----------|---------|
| `WEBULL_APP_KEY` | Webull OpenAPI app key (required) |
| `WEBULL_APP_SECRET` | Webull OpenAPI app secret, used to compute request signatures (required) |
| `WEBULL_REGION` | Optional region: `hk` (default), `us`, `jp`, `sg`, `th`, `au`, `my`, `uk`, `br`, `mx`, `za`, `eu` |
| `WEBULL_ENVIRONMENT` | Optional environment: `prod` / `production` (default) or `uat` / `sandbox` |
| `WEBULL_BASE_URL` | Optional override of the REST base URL only |
| `WEBULL_MQTT_URL` | Optional override of the MQTT broker address only |

```sh
export WEBULL_APP_KEY="..."
export WEBULL_APP_SECRET="..."
export WEBULL_ENVIRONMENT="sandbox"
```

`WEBULL_APP_SECRET` is used only to sign requests on the client side. It is never
sent as a request header. Explicit options such as `client.WithCredentials(key,
secret)` always take precedence over `WithEnv()` regardless of the order in which
they are passed to `client.New`.

## Sandbox

Use the sandbox while developing. Selecting the environment is enough: the SDK
resolves the Hong Kong sandbox hosts from the region and environment.

```go
cl, err := client.New(
	client.WithEnv(),          // credentials from the environment
	client.WithSandbox(),      // force the sandbox environment
)
if err != nil {
	return err
}
```

`client.WithSandbox()` is equivalent to `client.WithEnvironment(client.Sandbox)`.
It resolves the Hong Kong sandbox REST host:

```
https://api.sandbox.webull.hk
```

Webull publishes shared test accounts in its
[getting-started guide](https://developer.webull.hk/apis/docs/getting-started)
so you can try the API without applying for access. Only the host above belongs
in committed material. App keys, app secrets, and access tokens are per-account
secrets and must be supplied through the environment or your own secret store,
never checked into the repository or docs.

Sandbox data is intentionally limited. Expect a small symbol set (currently
`AAPL`) and token status `NORMAL` without 2FA. See [Sandbox](sandbox.md) for the
environment table and [Troubleshooting](troubleshooting.md) for limitations.

## First call

Construct a core client, obtain an access token, then bind a Market Data client
to it.

```go
package main

import (
	"context"
	"log"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/data"
)

func main() {
	cl, err := client.New(client.WithEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = cl.Close() }()

	ctx := context.Background()

	// Market Data requests require an access token. In the sandbox the token
	// is activated immediately; in production this waits for the 2FA window.
	if _, err := cl.EnsureToken(ctx); err != nil {
		log.Fatal(err)
	}

	market := data.New(cl)
	snaps, err := market.GetSnapshot(ctx, data.SnapshotQuery{
		Symbols:  []string{"AAPL"},
		Category: data.StockCategoryUS,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range snaps {
		log.Printf("%s %s", s.Symbol, s.Price)
	}
}
```

`data.New` takes the public `*client.Client`, so signing, retries, rate limiting,
and error handling are shared across every request. For streaming, see
[Streaming](streaming.md).

## Common first-call failures

| Symptom | Cause | Fix |
|---------|-------|-----|
| `UNAUTHORIZED: http 401` | Bad app key or secret | Check `WEBULL_APP_KEY` and `WEBULL_APP_SECRET` — no trailing whitespace |
| `INVALID_TOKEN: http 417` | Token failure, or another business failure using status 417 | Check the API message: verify token/region for a token error; do not assume every 417 is token-related |
| Token stays `PENDING` | Production 2FA not completed | In production, complete the Webull App verification within 5 minutes; in sandbox, tokens are `NORMAL` immediately |
| `FORBIDDEN: http 403` | Missing entitlement | The endpoint requires a paid subscription (e.g., Footprint, Display Solution) |
| `417 Invalid Symbol` | Symbol not in sandbox | Sandbox data is limited to `AAPL`; try that symbol first |
| Connection refused | Wrong host or network | Verify `WEBULL_ENVIRONMENT` is `sandbox`; check firewall and DNS |

The SDK maps every HTTP 417 to the historical `INVALID_TOKEN` code so existing
callers remain compatible, but Webull also uses 417 for unsupported categories,
invalid symbols, and rejected trading strategies. Use `Error.Status` and the
human-readable `Error.Message` for diagnosis; do not match the message text in
program logic.

## Next steps

- [Authentication](authentication.md) — how requests are signed and tokens are managed.
- [Market Data](market-data.md) — HTTP queries and the available endpoint groups.
- [Fundamentals](fundamentals.md) — capital flows, industry comparisons, earnings/dividend calendars, SEC filings, and financial statements.
- [Streaming](streaming.md) — real-time MQTT and MQTT-over-WebSocket pushes.
- [Trading](trading.md) — accounts, balances, positions, and the order lifecycle.
- [Trading Events](events.md) — order, position, and option events over gRPC.
- [Errors](errors.md) — typed errors and how to branch on them.
