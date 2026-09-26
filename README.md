# webullapi4go

An idiomatic Go SDK for the [Webull OpenAPI](https://developer.webull.hk/apis/docs/).
It provides typed clients for Webull's HTTP, MQTT, and gRPC services. The latest
repository tag is **`v2.1.3`** (2026-09-26); the current request, OMS,
streaming, telemetry, and documentation hardening is tagged in repository
`v2.1.3` and was introduced in `v2.1.1`. This is a repository patch release, not a published Go-semver v2
module: the root module path remains `github.com/shing1211/webullapi4go` and
stays on the v1 import path **by decision**, with no `/v2` migration planned.
The tagged tree is therefore not an installable published v2 module.

- Module: `github.com/shing1211/webullapi4go`
- Documentation: https://shing1211.github.io/webullapi4go/
- License: Apache-2.0
- Requires Go 1.26 or newer; no cgo.

## Feature matrix

Status distinguishes implementation and live availability. “Implemented” does
not mean every endpoint is live-verified; see the
[implementation status](IMPLEMENTATION_STATUS.md) for verification boundaries.

| Area | Implementation | Verification | Details |
|---|---|---|---|
| Authentication | Implemented | Offline-tested; core token flow live-verified in HK | HMAC-SHA1 REST signing, token create/check/ensure, automatic token injection |
| Market Data HTTP | Implemented | Offline-tested; selected HK calls live-verified | Instruments, fundamentals, futures, snapshot, tick, quotes/depth, bars, watchlists, options, news, event contracts, crypto/funds, and Display routes |
| Market Data MQTT | Implemented | Offline-tested; basic HK stream path previously live-verified | QUOTE, SNAPSHOT, and TICK over MQTT or MQTT-over-WebSocket, reconnect/resubscribe, health state, and bounded channels |
| Trading HTTP + OMS | Implemented | Offline-tested; selected HK account/preview paths live-verified | Accounts, assets, stock/single-leg/multi-leg/futures/event order validation, order queries, guardrails, and local order-state reconciliation |
| Trading events | Implemented | Offline-tested; basic stream path previously exercised | Order, position, and option streams over server-streaming gRPC, reconnect, correlation, spans, and metrics |
| Broker API HK | Implemented and offline-tested | Live blocked | Scope-protected virtual accounts, instruments, assets, orders, funding, journals, and events; separate `broker/` module |
| Broker FD API US | Implemented and offline-tested | Live blocked without US credentials | Agreements, accounts, documents, assets, activity, funding, instruments, orders, journals, and master data |
| Broker FD events | Implemented and offline-tested | Live blocked without US credentials | Broker FD event stream over gRPC using `grpc.event.EventService` |
| Display Solution | Implemented and offline-tested | Live blocked by host/entitlement | Company profile, analyst data, news, streaming, screeners, and quotes; HK host returns `403` |
| Connect OAuth | Implemented and offline-tested | US live access unavailable | Authorization-code URL builder and token exchange/refresh |

## Install

The module path is `github.com/shing1211/webullapi4go` and does not include the
`/v2` suffix required by Go's semantic import versioning; it stays that way by
decision, and no `/v2` migration is planned. The module proxy therefore serves
only the `v1.x` line. The `v2.1.3` repository tag does not make the tagged tree
installable, and an unqualified `go get` installs `v1.1.1`, which predates it.

To install the current tree, pin a commit:

```sh
go get github.com/shing1211/webullapi4go@78c164c
```

Go resolves that to a pseudo-version (`v1.1.2-0.20260926035012-78c164c22fc5` at
the time of writing). A commit pin moves only when you change it, whereas
`@main` tracks the newest code.

To pin the released v1.x line instead:

```sh
go get github.com/shing1211/webullapi4go@v1.1.1
```

The `broker/` package is a separate Broker API HK module. Its `go.mod` replaces
the root module with `../`; it is not a nested package that root
`go build ./...` traverses.

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
| `WithHTTPClient`, `WithHTTPTransport`, `WithTimeout`, `WithUserAgent` | Tune HTTP transport and connection pooling |
| `WithRetry(RetryConfig)`, `WithoutRetry()`, `WithResiliencePreset` | Configure transient-failure handling |
| `WithRateLimiter`, `NewRateLimiter(rate, burst)` | Throttle requests per path |
| `WithBreaker`, `NewBreaker(threshold, cooldown)` | Add a circuit breaker |
| `WithClockDriftCorrection` | Learn a bounded signing-clock offset from response `Date` headers |
| `WithInterceptor`, `WithHooks` | Add request-pipeline behavior and lifecycle callbacks |
| `WithAPIVersion("v2"\|"v3")`, `WithAPIVersionFor(prefix, version)` | Select the `x-version` header |
| `WithAutoToken(true)` | Obtain a sandbox token automatically before the first token-consuming request |
| `WithLogger`, `WithTracerProvider`, `WithMeterProvider`, `WithPropagator` | Configure structured logs and OpenTelemetry |
| `WithEnv()` | Fill configuration from the environment |

Streaming is configured with `stream.WithSessionID`, `WithMQTTURL`,
`WithWebSocket`, `WithAutoReconnect`, `WithAutoResubscribe`,
`WithResubscribeTimeout`, `WithKeepAlive`, `WithConnectTimeout`,
`WithWriteTimeout`, `WithMessageChannelDepth`, `WithCleanSession`,
`WithTLSConfig`, `WithHealthWatchdog`, and `WithMeter`.

Trading is configured with `trade.WithMaxOrderNotional`,
`trade.WithMaxOrderQuantity`, and `trade.WithAutoClientOrderID`. The first two
are advisory configuration guardrails enforced before placement. The notional
cap does not cover multi-leg option orders, which have no single top-level
notional; the quantity cap still applies. Auto client-order IDs are stable for
the same logical place request and do not mutate the caller's request.

Trading events are configured with `events.WithSubscribeTypes`,
`events.WithAccounts`, `events.WithGRPCEndpoint`, `events.WithGRPCPort`,
`events.WithTLS`, `events.WithDialTimeout`, `events.WithGRPCDialOption`,
`events.WithAutoReconnect`, `events.WithReconnectBaseDelay`,
`events.WithReconnectMaxDelay`, and `events.WithMaxReconnectAttempts`.

See [Observability](docs/observability.md) for OTel setup, metric names, event
spans, and correlation propagation.

## Sandbox testing

Use the sandbox while developing. Point the client at the environment and it
resolves the Hong Kong sandbox host `https://api.sandbox.webull.hk`:

```sh
export WEBULL_ENVIRONMENT="sandbox"
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
```

Core, data, stream, and events sandbox tests are skipped unless explicitly
enabled. Use the actual `Sandbox` test selector, and never commit credentials:

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

Trading-package sandbox tests use separate credentials and the
`WEBULL_TRADE_SANDBOX=1`, `WEBULL_TRADE_APP_KEY`,
`WEBULL_TRADE_APP_SECRET`, and `WEBULL_TRADE_ACCOUNT_ID` variables:

```sh
WEBULL_TRADE_SANDBOX=1 \
WEBULL_TRADE_APP_KEY=your-trading-sandbox-app-key \
WEBULL_TRADE_APP_SECRET=your-trading-sandbox-app-secret \
WEBULL_TRADE_ACCOUNT_ID=your-trading-account-id \
go test ./trade -run Sandbox
```

The mutating trade test additionally requires `WEBULL_TRADE_MUTATE=1`. The
mutating order-event test uses the generic `WEBULL_SANDBOX` credentials plus
`WEBULL_TRADE_ACCOUNT_ID` and `WEBULL_TRADE_MUTATE=1`. MQTT over WebSocket
tests additionally read `WEBULL_MQTT_WEBSOCKET=1`. See the local
[sandbox guide](docs/sandbox.md) and [troubleshooting guide](docs/troubleshooting.md)
for the complete gates and known limitations (a single supported symbol,
entitlement-gated footprint data, empty depth outside market hours, and network
blocking of plain MQTT on port 1883). HTTP 417 is mapped to the historical
`INVALID_TOKEN` category for compatibility, but Webull also uses it for invalid
symbols and other business validation; inspect the API message before treating
it as a token failure.

## Package map

| Package | Purpose |
|---------|---------|
| `client` | Canonical core SDK: configuration, signing, tokens, HTTP transport, request pipeline, and resilience |
| `data` | Canonical Market Data HTTP client and DTOs |
| `stream` | Canonical Market Data MQTT client with reconnect, health, and channel policies |
| `trade` | Canonical Trading HTTP client and OMS tracking integration |
| `events` | Canonical Trading Events gRPC client |
| `connect` | OAuth 2.0 authorization-code flow for third-party apps (US only) |
| `display` | Display Solution client-to-server authentication and token management |
| `broker` | Broker API HK (separate Go module; `broker/go.mod`: `replace github.com/shing1211/webullapi4go => ../`) |
| `brokerfd` | Broker FD US HTTP endpoints (accounts, orders, funding, instruments, etc.) |
| `brokerfd/events` | Broker FD US events over gRPC |
| `webull` | Optional thin aliases for the core client; service clients remain in their root packages |
| `pkg/errors` | Public typed errors, codes, and sentinels |
| `pkg/observability` | OpenTelemetry handles, span helpers, shared instruments, and sanitized telemetry error text |
| `pkg/resilience` | Public retry, rate-limit, circuit-breaker, and clock primitives |
| `pkg/transport` | Public HTTP transport and MQTT transport |
| `pkg/domain/money` | `money.Money`, the public DTO decimal type |
| `pkg/domain/order` | Public order state machine and reconciliation model |
| `pkg/types` | Shared public market/instrument types |
| `gen/webull/...` | Committed generated protobuf types; do not hand-edit |
| `internal/*` | Non-public authentication and compatibility implementation details |

## Roadmap and release status

| Version | Scope | Status |
|---|---|---|
| v0.x–v1.0 | Core API, Trading, Events, full endpoint coverage, examples, and API stabilization | Released |
| v1.1.1 | Production test/security tooling, multi-OS CI, leak checks, and fuzzing | Released |
| v2.0.0–v2.0.2 | Public error/resilience/transport foundations, OMS domain, `money.Money`, and thin `webull` aliases; root services retained | Historical tag; not a published v2 module |
| v2.0.3–v2.0.4 | Context hygiene, structured resilience, clock correction, idempotency, and transport tuning | Historical tag; not a published v2 module |
| v2.0.5–v2.0.7 | Request interceptors/hooks, initial OMS, stream state/channels, slog, OTel tracing, and metrics | Historical tag; not a published v2 module |
| v2.0.8–v2.0.9 | Context and typed-error hardening | Historical tag; not a published v2 module |
| v2.1.0 | `go vet` mutex-copy fixes | Historical tag; not a published v2 module |
| v2.1.1 | Error matching specificity, request-pipeline parity, OMS reconciliation, stream/MQTT lifecycle hardening, gRPC event telemetry, cancellation/leak coverage, and documentation reconciliation | Repository tag; not a published v2 module; offline-tested and not newly live-verified |

The next version number is intentionally unassigned for future work. The
current hardening is represented by repository tag `v2.1.3`. The module-path
decision is settled (stay on the v1 import path, no `/v2` migration); live
verification remains a separate follow-up decision. See
[CHANGELOG.md](CHANGELOG.md) and [PLAN.md](PLAN.md) for the decision record.

## Links

- Documentation: https://shing1211.github.io/webullapi4go/
- API reference: https://pkg.go.dev/github.com/shing1211/webullapi4go
- `data` reference: https://pkg.go.dev/github.com/shing1211/webullapi4go/data
- `stream` reference: https://pkg.go.dev/github.com/shing1211/webullapi4go/stream
- `trade` reference: https://pkg.go.dev/github.com/shing1211/webullapi4go/trade
- `events` reference: https://pkg.go.dev/github.com/shing1211/webullapi4go/events
- Questions and ideas: [GitHub Discussions](https://github.com/shing1211/webullapi4go/discussions)
- Architecture decisions: [ADR index](docs/adr/index.md)
- Implementation status: [IMPLEMENTATION_STATUS.md](IMPLEMENTATION_STATUS.md)
- Observability: [docs/observability.md](docs/observability.md)
- Changelog: [CHANGELOG.md](CHANGELOG.md)
- Contributing: [CONTRIBUTING.md](CONTRIBUTING.md)

## Development

| Target | Description |
|--------|-------------|
| `make build` | compile the root module and every nested module listed by the Makefile |
| `make vet` | vet the root module and every nested module listed by the Makefile |
| `make test` | unit tests across those modules (sandbox tests remain environment-gated) |
| `make test-race` | tests with the race detector across those modules |
| `make cover` | coverage profiling across those modules |
| `make lint` | `golangci-lint` (includes `gosec`) |
| `make fuzz` | fuzz the data deserializers |
| `make vuln` | `govulncheck` in every module |
| `make docs` | build the MkDocs site (`mkdocs build --strict`) |

Root `go build ./...` and `go test ./...` cover only the root module; use the
Makefile targets for module-aware verification. The 2026-09-26 offline run
measured 73.6% aggregate coverage in the root module and 80.8% in the nested
`broker/` module. These are reproducible measurements, not behavior or release
guarantees; CI's 60% gate is root-only.

The current hardening was not newly live-verified. Historical live results are
listed separately and must not be read as verification of the `v2.1.1`
repository tag.

## License

Apache-2.0. See [LICENSE](LICENSE).

## Disclaimer

This SDK is an independent, unofficial wrapper around the Webull OpenAPI and is
provided for reference and educational use only. It is not affiliated with or
endorsed by Webull. Nothing here is investment advice. Use it at your own risk,
and always follow Webull's terms of service and your own risk controls.
