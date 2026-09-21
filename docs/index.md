# webullapi4go

`webullapi4go` is an idiomatic Go SDK for the [Webull OpenAPI](https://developer.webull.com/apis/docs/).
It wraps Webull's HTTP and MQTT services in typed Go, starting with the Hong Kong
region.

The v0.1 surface covers authentication, a core signed REST client, the Market
Data HTTP API, and real-time Market Data streaming over MQTT. The v0.2 releases
add the Trading HTTP API for accounts, balances, positions, and the stock,
options, and combo order lifecycle. The v0.3.0 release adds Trading events over
gRPC. The v0.4.0 release adds Market Data fundamentals: capital flows, industry
comparisons, earnings and dividend calendars, SEC filings, and financial
statements. The v0.5 release adds Display Solution (corporate actions, instrument
profiles, logos), fund and crypto data, and screener v2. The v0.6 release adds
multi-leg options orders, futures order validation (provisional), and speculative
option chain and expiration discovery. The v0.7 release adds event contracts,
Broker API HK, Broker FD US, and Broker FD gRPC events. The module is licensed
under Apache-2.0.

## Feature matrix

| Area | Status |
|------|--------|
| Authentication | Supported |
| Market Data (HTTP) | Supported |
| Market Data Fundamentals | Supported (v0.4.0: capital flows, industry comparisons, earnings/dividend calendars, SEC filings, financial statements) |
| Market Data (MQTT streaming) | Supported |
| Trading (HTTP) | Supported |
| Trading events (gRPC) | Supported (v0.3.0: order, position, and option streams) |
| Broker API HK | Supported (v0.7) |
| Broker FD US | Provisional (v0.7) |
| Display Solution | Provisional (v0.7) | Company profile, analyst data, news, streaming, screeners, quotes — requires paid Webull subscription |

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
- [Sandbox](sandbox.md) — environments, test credentials, limitations.
- [Errors](errors.md) — typed errors and classification.
- [API Reference](api.md) — package overview and pkg.go.dev links.
- [Troubleshooting](troubleshooting.md) — symptoms and fixes.

## Links

- Source and issues: [github.com/shing1211/webullapi4go](https://github.com/shing1211/webullapi4go)
- Questions and ideas: [GitHub Discussions](https://github.com/shing1211/webullapi4go/discussions)
- Architecture decisions: [ADR index](adr/index.md)
- Webull API reference: [developer.webull.com](https://developer.webull.com/apis/docs/)

## Disclaimer

This SDK is an independent, unofficial wrapper and is not affiliated with or
endorsed by Webull. It is provided for reference and educational use only and is
not investment advice.
