# webullapi4go

An idiomatic Go SDK for the [Webull OpenAPI](https://developer.webull.com/apis/docs/).
It wraps Webull's HTTP and MQTT services in typed Go, starting with the Hong Kong
region. The v0.1 release covers authentication, a core signed REST client, the
Market Data HTTP API, and real-time Market Data streaming over MQTT. The v0.2
release adds the Trading HTTP API: v0.2.1 introduced account listing, balances,
and positions, v0.2.2 adds the stock-order lifecycle (preview, place, replace,
cancel) and order queries, v0.2.3 adds market-specific order rules for US, HK,
and CN, including Hong Kong BCAN party IDs, v0.2.4 adds single-leg options
orders, and v0.2.5 adds US combo orders. The v0.3 release adds real-time Trading
events over gRPC. The v0.4 release adds Market Data fundamentals: capital flows,
industry comparisons, earnings and dividend calendars, SEC filings, and financial
statements.

- Module: `github.com/shing1211/webullapi4go`
- Documentation: https://shing1211.github.io/webullapi4go/
- License: Apache-2.0
- Requires Go 1.26 or newer; no cgo.

## Feature matrix (v0.9.2)

| Area | Status | Details |
|------|--------|---------|
| Authentication | Supported | HMAC-SHA1 request signing, token create/check/ensure, automatic token injection |
| Market Data (HTTP) | Supported | Instruments, fundamentals, futures static, snapshot, tick, quotes/depth, bars, footprint, NOII, screener, watchlists, options, news, event contracts. **Provisional**: futures market data, Display Solution endpoints (TODO markers) |
| Market Data (MQTT streaming) | Supported | QUOTE, SNAPSHOT, and TICK pushes over MQTT or MQTT-over-WebSocket, with auto-reconnect and auto-resubscribe |
| Trading (HTTP) | Supported | Accounts, balances, and positions; stock order preview, place, replace, cancel, and order queries; US/HK/CN order-type rules, Hong Kong BCAN; single-leg options orders; US combo orders; multi-leg options orders (provisional); futures order validation (provisional); event contract orders (provisional); batch place orders |
| Trading events (gRPC) | Supported | Order, event-contract position, and option status-change streams over server-streaming gRPC |
| Broker API HK | Supported | Virtual accounts, instruments, assets, orders, cash activities, funding FX, instant funding, journals, master data, event contracts (`broker/` module) |
| Broker FD API US | Supported | Agreements, accounts, documents, assets, activity, funding, instruments, orders, journals, master data (`brokerfd/` module) |
| Broker FD events (gRPC) | Supported | Broker FD event stream over gRPC using `grpc.event.EventService` (`brokerfd/events/` module) |
| Display Solution | Supported | Company profile, analyst data, news, streaming, screeners, quotes (provisional paths, TODO markers) |

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

### Trading: accounts and positions

```go
package main

import (
	"context"
	"log"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/trade"
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

	trading := trade.New(cl)
	accounts, err := trading.ListAccounts(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, acct := range accounts {
		balance, err := trading.GetBalance(ctx, acct.AccountID)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s cash=%s market_value=%s", acct.AccountID,
			balance.TotalCashBalance, balance.TotalMarketValue)
	}
}
```

### Trading events over gRPC

```go
package main

import (
	"context"
	"log"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/events"
)

func main() {
	cl, err := client.New(client.WithEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = cl.Close() }()

	ev, err := events.New(cl, events.WithSubscribeTypes(events.SubscribeOrder))
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = ev.Close() }()

	ev.OnConnect(func() { log.Println("subscribed") })
	ev.OnOrder(func(o *events.OrderEvent) {
		log.Printf("%s %s %s", o.OrderID, o.OrderStatus, o.SceneType)
	})

	if err := ev.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
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

Trading is configured with `trade.WithMaxOrderNotional` and
`trade.WithMaxOrderQuantity`, advisory order guardrails that the order methods
enforce before an order is built. The notional cap does not cover multi-leg
option orders, which have no single top-level notional; the quantity cap still
applies to them.

Trading events are configured with `events.WithSubscribeTypes`,
`events.WithAccounts`, `events.WithGRPCEndpoint`, `events.WithGRPCPort`,
`events.WithTLS`, `events.WithDialTimeout`, `events.WithGRPCDialOption`,
`events.WithAutoReconnect`, `events.WithReconnectBaseDelay`,
`events.WithReconnectMaxDelay`, and `events.WithMaxReconnectAttempts`.

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
| `trade` | Trading HTTP endpoints (accounts, balances, positions, stock, single-leg and multi-leg options, futures validation, US combo orders, and order queries) |
| `events` | Trading events over gRPC: order, position, and option streams with typed payloads and reconnect |
| `gen/webull/marketdata/v1` | Generated protobuf types for streamed messages |
| `gen/webull/trade/events/v1` | Generated protobuf types for the gRPC event service |
| `pkg/types` | Shared public domain types (markets, instrument types) |
| `internal/*` | Implementation details: signing, token lifecycle, region endpoints, transport, resilience, MQTT |

## Roadmap

| Version | Scope | Status |
|---------|-------|--------|
| v0.1 | Authentication, core HTTP client, Market Data HTTP + MQTT streaming | Done |
| v0.2 | Trading (HTTP): accounts, balances, positions (v0.2.1), stock orders (v0.2.2), market-specific rules and HK BCAN (v0.2.3), single-leg options orders (v0.2.4), US combo orders (v0.2.5), then a Market Data news SSE refactor (v0.2.6) | Done |
| v0.3 | Trading events over gRPC | Done |
| v0.4 | Market Data fundamentals: capital flows, industry comparisons, earnings/dividend calendars, SEC filings, financial statements | Done |
| v0.5 | Fund data, crypto data, screener v2, corporate actions, instrument v3 migration | Done (provisional) |
| v0.6 | Broker API (HK + FD US), multi-leg options, futures validation, option chain discovery | Done (provisional) |
| v0.7.0 | Event contracts, Broker API HK, Broker FD US, Broker FD gRPC events | Done |
| v0.8.0 | — | — |
| v0.9.0 | GoDoc coverage on brokerfd/ and brokerfd/events/, HK options stubs, HK futures market data, new examples (brokerfd, brokerfd-events, options), graceful credential handling | Done |
| v0.9.1 | Watchlist boolean-response fix; `DoBroker` transport; `watchlist-cmd` and `broker-probe` examples | Done |
| v0.9.2 | Broker HK path correction (`/openapi/v1/broker/...` → `/broker/...`); `401 ROUTE_NOT_PERMITTED` instead of `404 Route Not Found` | Done |
| v1.0 | Stable public API, full documentation, semver guarantees | Planned |

## Links

- Documentation: https://shing1211.github.io/webullapi4go/
- API reference: https://pkg.go.dev/github.com/shing1211/webullapi4go
- `data` reference: https://pkg.go.dev/github.com/shing1211/webullapi4go/data
- `stream` reference: https://pkg.go.dev/github.com/shing1211/webullapi4go/stream
- `trade` reference: https://pkg.go.dev/github.com/shing1211/webullapi4go/trade
- `events` reference: https://pkg.go.dev/github.com/shing1211/webullapi4go/events
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
