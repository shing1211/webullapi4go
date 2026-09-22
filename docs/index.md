# webullapi4go

`webullapi4go` is an idiomatic Go SDK for the [Webull OpenAPI](https://developer.webull.hk/apis/docs/).
It wraps Webull's HTTP and MQTT services in typed Go, starting with the Hong Kong
region.

The current release (v1.1.0) covers authentication, a core signed REST client,
Market Data HTTP and MQTT streaming, the Trading HTTP API, Trading events over
gRPC, fundamentals and fund data, crypto data, Display Solution, event
contracts, options and futures, Broker API HK, Broker FD US, and the Connect
API. Every endpoint documented by Webull is implemented, with paths taken from
the official OpenAPI definition — see the
[reconciliation report](reconciliation.md). The module is licensed under
Apache-2.0.

## Feature matrix

| Area | Status |
|------|--------|
| Authentication | Supported |
| Market Data (HTTP) | Supported |
| Market Data Fundamentals | Supported — capital flows, industry comparisons, earnings/dividend calendars, SEC filings, financial statements |
| Market Data (MQTT streaming) | Supported |
| Trading (HTTP) | Supported — including multi-leg options and futures order validation |
| Trading events (gRPC) | Supported — order, position, and option streams |
| Event contracts | Supported |
| Crypto market data | Supported — US-only; HK sandbox returns `404` |
| Fund data | Supported — info, NAV, dividends, plus performance, holdings, rating, splits, files, allocation |
| Broker API HK | Supported — HK sandbox returns `401 ROUTE_NOT_PERMITTED` (app scope missing) |
| Broker FD US | Supported — US-only; HK sandbox returns `404` |
| Display Solution | Supported, entitlement-gated — HK sandbox returns `403` |
| Connect API (OAuth) | Supported |

## What's new in v1.1

- **209 endpoints** — full parity with the official Webull OpenAPI (was ~140 in v1.0)
- **6 new packages** — `connect`, `display`, `broker`, `brokerfd`, `brokerfd/events`, plus new protobuf types
- **Zero TODO markers** — all provisional scaffolding removed
- **Sandbox-tested** — 20/20 integration tests pass against the HK sandbox

## Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                       Your Application                           │
├──────────┬──────────┬──────────┬──────────┬──────────────────────┤
│  client  │   data   │  trade   │  stream  │       events         │
│  (core)  │ (market) │ (orders) │  (MQTT)  │       (gRPC)         │
├──────────┴──────────┴──────────┴──────────┴──────────────────────┤
│                 internal/auth · errs · region                    │
└──────────────────────────────────────────────────────────────────┘
  Display: data.DisplayService()   │  Broker: broker/, brokerfd/
  Connect: connect/ (OAuth)        │  Types:  pkg/types/
```

## Install

```sh
go get github.com/shing1211/webullapi4go
```

Requires Go 1.26 or newer.

## Quickstart

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

Set the credentials in the environment first (see
[Getting Started](getting-started.md)):

```sh
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
export WEBULL_ENVIRONMENT="sandbox"
```

## Documentation

- [Getting Started](getting-started.md) — install, credentials, sandbox, first call.
- [Authentication](authentication.md) — request signing and token lifecycle.
- [Market Data](market-data.md) — HTTP queries (requires [Authentication](authentication.md)).
- [Fundamentals](fundamentals.md) — capital flows, industry comparisons, earnings/dividend calendars, SEC filings, financial statements.
- [Streaming](streaming.md) — real-time MQTT pushes (requires [Authentication](authentication.md)).
- [Trading](trading.md) — accounts, balances, positions, and orders (requires [Authentication](authentication.md)).
- [Trading Events](events.md) — order, position, and option events over gRPC (requires account ID from [Trading](trading.md)).
- [Display Solution](webull-api/display-solution.md) — entitlement-gated market data (compare with [Market Data](market-data.md)).
- [Broker API HK](broker-hk.md) — HK broker endpoints.
- [Broker FD US](broker-fd-us.md) — US Broker FD endpoints (US equivalent of [Broker HK](broker-hk.md)).
- [Connect API](webull-api/connect-api.md) — OAuth 2.0 authorization-code flow.
- [Sandbox](sandbox.md) — environments, test credentials, limitations.
- [Errors](errors.md) — typed errors and classification.
- [Patterns](patterns.md) — shared SDK patterns (pagination, options, error handling).
- [API Reference](api.md) — package overview and pkg.go.dev links.
- [Webull API Reference](webull-api.md) — every official endpoint mapped to its SDK method.
- [SDK ↔ API Reconciliation](reconciliation.md) — coverage report: implemented endpoints, gaps, path parity.
- [Troubleshooting](troubleshooting.md) — symptoms and fixes.

## Links

- Source and issues: [github.com/shing1211/webullapi4go](https://github.com/shing1211/webullapi4go)
- Questions and ideas: [GitHub Discussions](https://github.com/shing1211/webullapi4go/discussions)
- Architecture decisions: [ADR index](adr/index.md)
- Webull API reference: [developer.webull.hk](https://developer.webull.hk/apis/docs/)

## Disclaimer

This SDK is an independent, unofficial wrapper and is not affiliated with or
endorsed by Webull. It is provided for reference and educational use only and is
not investment advice.
