# webullapi4go

`webullapi4go` is an idiomatic Go SDK for the
[Webull OpenAPI](https://developer.webull.hk/apis/docs/). It wraps Webull's HTTP,
MQTT, and gRPC services in typed Go clients.

The latest tagged release is `v2.1.0`. Current request-pipeline, OMS,
streaming, and event-telemetry hardening is **Unreleased**, implemented, and
offline-tested; it is not presented as live-verified. See the
[implementation status](implementation-status.md) and
[changelog](https://github.com/shing1211/webullapi4go/blob/main/CHANGELOG.md).

## Coverage

The generated [SDK ↔ API Reconciliation](reconciliation.md) reports 209
implemented endpoints and 0 documented-only gaps in its 2026-09-22 snapshot.
It also reports four paths that match only the docs summary and 25 unresolved
paths, so coverage does not mean every endpoint has been live-verified.

## Feature status

| Area | Implementation | Verification boundary |
|---|---|---|
| Authentication | Implemented | Core token flow live-verified in HK; production 2FA is user-driven |
| Market Data HTTP | Implemented | Offline-tested; selected HK calls live-verified, US/entitlement surfaces partial |
| Market Data MQTT | Implemented | Offline-tested; basic HK streaming path previously live-verified |
| Trading HTTP + OMS | Implemented | Offline-tested; accounts/assets/preview and guarded paths exercised in HK |
| Trading Events | Implemented | Offline-tested; basic stream path previously exercised |
| Broker API HK | Implemented | Offline-tested; HK blocked by missing route scope (`401`) |
| Broker FD US | Implemented | Offline-tested; live verification blocked without US credentials |
| Display Solution | Implemented | Offline-tested; HK host returns `403` |
| Connect OAuth | Implemented | Offline-tested; US live access unavailable |

## Architecture

```text
┌──────────────────────────────────────────────────────────────────────┐
│                          Your Application                            │
├───────────┬───────────┬───────────┬───────────┬──────────────────────┤
│  client   │   data    │   trade   │  stream   │ events / connect /   │
│  (core)   │ (market)  │ (trading) │  (MQTT)   │ display / brokerfd   │
├───────────┴───────────┴───────────┴───────────┴──────────────────────┤
│ Shared foundations: pkg/errors · observability · resilience ·        │
│ transport · types · domain/money · domain/order                      │
├──────────────────────────────────────────────────────────────────────┤
│ internal: signing and compatibility implementation details           │
└──────────────────────────────────────────────────────────────────────┘
         broker/ = separate Broker API HK Go module
         webull/ = optional thin aliases for the core client
```

The root service packages are canonical. A full relocation under `pkg/` is
superseded and closed. Public financial DTOs use `money.Money` or
`*money.Money`, not raw `decimal.Decimal`.

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
	for _, snap := range snaps {
		log.Printf("%s %s", snap.Symbol, snap.Price)
	}
}
```

Set credentials through the environment or `client.WithCredentials`:

```sh
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
export WEBULL_ENVIRONMENT="sandbox"
```

## Core patterns

- `client.New` constructs one shared core; `data.New`, `trade.New`,
  `stream.New`, and `events.New` bind service clients to it.
- Financial DTOs preserve decimal precision through `money.Money` while JSON
  remains a string.
- `pkg/errors` exposes stable `Code` values and sentinels for `errors.Is` and
  `errors.As`; do not match error strings.
- Trading tracks successful place/batch results and can reconcile order
  snapshots or event-driven states locally.
- Streaming exposes lifecycle state, a data-age watchdog, and bounded channels
  with blocking, drop-oldest, or sampling policies.
- REST, MQTT, and both gRPC event clients inherit core logging, correlation,
  and OTel configuration.

## Documentation

- [Getting Started](getting-started.md) — install, credentials, sandbox, first call.
- [Authentication](authentication.md) — request signing and token lifecycle.
- [Market Data](market-data.md) — HTTP market queries.
- [Streaming](streaming.md) — MQTT state, reconnect, health, and channels.
- [Trading](trading.md) — accounts, orders, guardrails, and OMS reconciliation.
- [Trading Events](events.md) — order, position, and option events over gRPC.
- [Broker API HK](broker-hk.md) — HK broker endpoints.
- [Broker FD US](broker-fd-us.md) — US Broker FD endpoints.
- [Connect API](webull-api/connect-api.md) — OAuth 2.0 authorization-code flow.
- [Sandbox](sandbox.md) — environments, credentials, and limitations.
- [Errors](errors.md) — typed errors and stable classification.
- [Observability](observability.md) — slog, OpenTelemetry, metrics, and correlation.
- [Patterns](patterns.md) — shared client, money, pagination, and error patterns.
- [Go Packages](api.md) — public package map.
- [Webull API Reference](webull-api.md) — documented endpoints mapped to SDK methods.
- [SDK ↔ API Reconciliation](reconciliation.md) — generated coverage and path states.
- [Implementation Status](implementation-status.md) — released, Unreleased, offline-tested, live-verified, and blocked work.
- [Troubleshooting](troubleshooting.md) — symptoms and fixes.

## Links

- Source and issues: [github.com/shing1211/webullapi4go](https://github.com/shing1211/webullapi4go)
- Questions and ideas: [GitHub Discussions](https://github.com/shing1211/webullapi4go/discussions)
- Architecture decisions: [ADR index](adr/index.md)
- Official Webull docs: [developer.webull.hk](https://developer.webull.hk/apis/docs/)

## Disclaimer

This SDK is an independent, unofficial wrapper and is not affiliated with or
endorsed by Webull. It is provided for reference and educational use and is not
investment advice.
