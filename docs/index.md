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
- [Market Data](market-data.md) — HTTP queries.
- [Fundamentals](fundamentals.md) — capital flows, industry comparisons, earnings/dividend calendars, SEC filings, financial statements.
- [Streaming](streaming.md) — real-time MQTT pushes.
- [Trading](trading.md) — accounts, balances, positions, and orders.
- [Trading Events](events.md) — order, position, and option events over gRPC.
- [Broker API HK](broker-hk.md) — HK broker endpoints.
- [Broker FD US](broker-fd-us.md) — US Broker FD endpoints.
- [Sandbox](sandbox.md) — environments, test credentials, limitations.
- [Errors](errors.md) — typed errors and classification.
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
