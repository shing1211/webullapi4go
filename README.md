# webullapi4go

An idiomatic Go SDK for the [Webull OpenAPI](https://developer.webull.com/apis/docs/).
It wraps Webull's HTTP and MQTT services in typed Go, starting with the Hong Kong
region. The v0.1 release covers authentication, a core signed REST client, the
Market Data HTTP API, and real-time Market Data streaming over MQTT.

- Module: `github.com/shing1211/webullapi4go`
- Documentation: https://shing1211.github.io/webullapi4go/
- License: Apache-2.0
- Requires Go 1.26 or newer; no cgo.

## Feature matrix (v0.1)

| Area | Status | Details |
|------|--------|---------|
| Authentication | Supported | HMAC-SHA1 request signing, token create/check/ensure, automatic token injection |
| Market Data (HTTP) | Supported | Instruments, company profile, analyst data, futures static, snapshot, tick, quotes/depth, bars (single and batch), footprint, NOII, screener, watchlists, options, news |
| Market Data (MQTT streaming) | Supported | QUOTE, SNAPSHOT, and TICK pushes over MQTT or MQTT-over-WebSocket, with auto-reconnect and auto-resubscribe |
| Trading (HTTP) | Not yet | Planned for v0.2 |
| Trading events (gRPC) | Not yet | Planned for v0.3 |
| Display Solution | Not yet | Planned for v0.4 |
| Broker API | Not yet | Planned for v0.5 |

## Install

```sh
go get github.com/shing1211/webullapi4go
```

## Quickstart

### HTTP market data

```go
package main

import (
	"context"
	"log"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/data"
)

func main() {
	// WithEnv reads WEBULL_APP_KEY, WEBULL_APP_SECRET, WEBULL_REGION, and
	// WEBULL_ENVIRONMENT.
	cl, err := client.New(client.WithEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = cl.Close() }()

	ctx := context.Background()
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

### Streaming over MQTT

```go
package main

import (
	"context"
	"log"

	"github.com/shing1211/webullapi4go/client"
	marketdatav1 "github.com/shing1211/webullapi4go/gen/webull/marketdata/v1"
	"github.com/shing1211/webullapi4go/stream"
)

func main() {
	cl, err := client.New(client.WithEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = cl.Close() }()

	ctx := context.Background()
	if _, err := cl.EnsureToken(ctx); err != nil {
		log.Fatal(err)
	}

	s, err := stream.New(cl, stream.WithWebSocket(true))
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = s.Close() }()

	s.OnSnapshot(func(snap *marketdatav1.Snapshot) {
		log.Printf("%s %s", snap.GetBasic().GetSymbol(), snap.GetPrice())
	})

	if err := s.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	if err := s.Subscribe(ctx, stream.SubscribeRequest{
		Symbols:  []string{"AAPL"},
		Category: stream.CategoryUSStock,
		SubTypes: []stream.SubType{stream.SubTypeQuote, stream.SubTypeSnapshot, stream.SubTypeTick},
	}); err != nil {
		log.Fatal(err)
	}

	select {} // block until the process is interrupted
}
```

More runnable programs live in [`examples/`](examples/README.md).

## Configuration

Credentials and environment are never hard-coded. `client.WithEnv()` reads the
following process environment variables:

| Variable | Purpose |
|----------|---------|
| `WEBULL_APP_KEY` | Webull OpenAPI app key (required) |
| `WEBULL_APP_SECRET` | Webull OpenAPI app secret, used to sign requests (required) |
| `WEBULL_REGION` | Deployment region: `hk` (default), `us`, `jp`, `sg`, `th`, `au`, `my`, `uk`, `br`, `mx`, `za`, `eu` |
| `WEBULL_ENVIRONMENT` | `prod` / `production` (default) or `uat` / `sandbox` |
| `WEBULL_BASE_URL` | Optional override of the REST base URL only |
| `WEBULL_MQTT_URL` | Optional override of the MQTT broker address only |

Explicit functional options always take precedence over `WithEnv()` regardless of
order. Common options:

| Option | Purpose |
|--------|---------|
| `WithAppKey`, `WithAppSecret`, `WithCredentials` | Set credentials directly |
| `WithRegion(client.HK)` | Select the deployment region |
| `WithEnvironment(client.Sandbox)`, `WithSandbox()` | Select the environment |
| `WithBaseURL`, `WithEndpoints` | Override resolved service endpoints |
| `WithHTTPClient`, `WithTimeout`, `WithUserAgent` | Tune the HTTP transport |
| `WithRetry(RetryConfig)`, `WithoutRetry()` | Configure transient-failure retries |
| `WithRateLimiter`, `NewRateLimiter(rate, burst)` | Throttle requests per path |
| `WithBreaker`, `NewBreaker(threshold, cooldown)` | Add a circuit breaker |
| `WithAPIVersion("v2"\|"v3")`, `WithAPIVersionFor(prefix, version)` | Select the `x-version` header |
| `WithAutoToken(true)` | Obtain a token automatically (sandbox) before the first request |
| `WithEnv()` | Fill configuration from the environment |

Streaming is configured with `stream.WithSessionID`, `WithMQTTURL`,
`WithWebSocket`, `WithAutoReconnect`, `WithAutoResubscribe`,
`WithResubscribeTimeout`, `WithKeepAlive`, `WithConnectTimeout`,
`WithWriteTimeout`, `WithMessageChannelDepth`, `WithCleanSession`, and
`WithTLSConfig`.

## Sandbox testing

Use the sandbox while developing. Point the client at the environment and it
resolves the Hong Kong sandbox host `https://api.sandbox.webull.hk`:

```sh
export WEBULL_ENVIRONMENT="sandbox"
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
```

The integration tests are skipped unless explicitly enabled and require sandbox
credentials. Never commit the values.

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

MQTT over WebSocket tests additionally read `WEBULL_MQTT_WEBSOCKET=1`. See the
[sandbox documentation](https://shing1211.github.io/webullapi4go/sandbox/) and
[troubleshooting guide](https://shing1211.github.io/webullapi4go/troubleshooting/)
for known sandbox limitations (a single supported symbol, entitlement-gated
footprint data, empty depth outside market hours, and network blocking of plain
MQTT on port 1883).

## Package map

| Package | Purpose |
|---------|---------|
| `client` | Core SDK: configuration, options, signing, tokens, transport, and `Client.Do` |
| `data` | Market Data HTTP endpoints (typed requests and responses) |
| `stream` | Market Data streaming over MQTT, with reconnect and resubscribe |
| `gen/webull/marketdata/v1` | Generated protobuf types for streamed messages |
| `pkg/types` | Shared public domain types (markets, instrument types) |
| `internal/*` | Implementation details: signing, token lifecycle, region endpoints, transport, resilience, MQTT |

## Roadmap

| Version | Scope | Status |
|---------|-------|--------|
| v0.1 | Authentication, core HTTP client, Market Data HTTP + MQTT streaming | Done |
| v0.2 | Trading (HTTP): orders, positions, accounts | Planned |
| v0.3 | Trading events over gRPC | Planned |
| v0.4 | Display Solution | Planned |
| v0.5 | Broker API | Planned |
| v1.0 | Stable public API, full documentation, semver guarantees | Planned |

## Links

- Documentation: https://shing1211.github.io/webullapi4go/
- API reference: https://pkg.go.dev/github.com/shing1211/webullapi4go
- `data` reference: https://pkg.go.dev/github.com/shing1211/webullapi4go/data
- `stream` reference: https://pkg.go.dev/github.com/shing1211/webullapi4go/stream
- Questions and ideas: [GitHub Discussions](https://github.com/shing1211/webullapi4go/discussions)
- Architecture decisions: [ADR index](docs/adr/index.md)
- Changelog: [CHANGELOG.md](CHANGELOG.md)
- Contributing: [CONTRIBUTING.md](CONTRIBUTING.md)

## License

Apache-2.0. See [LICENSE](LICENSE).

## Disclaimer

This SDK is an independent, unofficial wrapper around the Webull OpenAPI and is
provided for reference and educational use only. It is not affiliated with or
endorsed by Webull. Nothing here is investment advice. Use it at your own risk,
and always follow Webull's terms of service and your own risk controls.
