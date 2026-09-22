# Webull API Reference

`webullapi4go` implements the [Webull OpenAPI](https://developer.webull.hk/apis/docs).
This section is generated from the **official OpenAPI definitions** and maps each
documented endpoint to its Go method, request fields and response fields.

- **Source of truth:** the machine-readable [`llms.txt`](https://developer.webull.hk/apis/llms.txt)
  index and the `.md` variant of each reference page (which embeds the full
  OpenAPI definition JSON). Broker FD US is documented on the
  [US site](https://developer.webull.com/apis/llms.txt).
- Every page below links back to the official reference for that endpoint.
- The official `path` is authoritative; the SDK follows it (v1.1.0 aligned all
  paths — see [Path status](#path-status) below).

Go type references live on pkg.go.dev:
[`client`](https://pkg.go.dev/github.com/shing1211/webullapi4go/client) ·
[`data`](https://pkg.go.dev/github.com/shing1211/webullapi4go/data) ·
[`trade`](https://pkg.go.dev/github.com/shing1211/webullapi4go/trade) ·
[`broker`](https://pkg.go.dev/github.com/shing1211/webullapi4go/broker) ·
[`brokerfd`](https://pkg.go.dev/github.com/shing1211/webullapi4go/brokerfd) ·
[`stream`](https://pkg.go.dev/github.com/shing1211/webullapi4go/stream) ·
[`events`](https://pkg.go.dev/github.com/shing1211/webullapi4go/events).

## How to read an entry

Each endpoint entry is laid out as:

| Row | Meaning |
|-----|---------|
| `METHOD /path` | The official endpoint path from the OpenAPI definition. |
| **SDK** | The `webullapi4go` method that calls it. |
| **Reference** | The official page (`.md` variant). |
| **Note** | SDK-specific caveats (sandbox behaviour, entitlements, region limits). |
| **Request — parameters** | Query/path parameters with type, required flag, and enum values. |
| **Request body** | JSON body schema; nested objects are shown as follow-up tables. |
| **Response 200** | Success schema fields. |
| **Errors** | `401` / `417` / `500`; see [Errors](errors.md) for the typed model. |

Prices, sizes, and quantities are **strings on the wire** and are kept as strings
in the Go DTOs to preserve precision.

## Environments and base URLs

| API | Service | Production | Sandbox |
|-----|---------|------------|---------|
| Trading API | HTTP | `api.webull.hk` | `api.sandbox.webull.hk` |
| Trading API | Events (gRPC) | `events-api.webull.hk` | `events-api.sandbox.webull.hk` |
| Market Data API | HTTP | `api.webull.hk` | `api.sandbox.webull.hk` |
| Market Data API | Streaming (MQTT) | `data-api.webull.hk` | `data-api.sandbox.webull.hk` |
| Broker API HK | HTTP | `broker-api.webull.hk` | `broker-api.sandbox.webull.hk` |
| Display Solution | HTTP | `co-branding-openapi.webull.hk` | `hk-co-branding-openapi.uat.webullbroker.com` |

Regions are selected with `client.WithRegion`; the endpoint set is derived by
`client.EndpointsFor`. Only the HK hosts are officially published; other regions
are inferred by analogy.

## Authentication

Webull uses a **dual layer**: a signed request plus an access token.

- **Server-to-server (Trading, Broker, Non-Display Market Data)** — every request
  carries the six `x-signature*`/`host` headers computed with **HMAC-SHA1** over a
  canonical string (`x-app-key`, `x-signature-algorithm`, `x-signature-version`,
  `x-signature-nonce`, `x-timestamp`, `host`, the upper-cased MD5 body digest,
  and the request path), plus `x-access-token` and `x-version`.
- **Client-to-server (Display Solution)** — a separate host using an OAuth-style
  client token sent as `Authorization: Bearer`; the SDK fetches it via
  `display.Service`.
- **Events (gRPC)** — HMAC-SHA256 over the serialized request, sent as gRPC
  metadata; no `host` participates.

See [Authentication](webull-api/authentication.md) and the
[Authentication](authentication.md) and [Errors](errors.md) guides.

## API versioning

`x-version` selects the interface version (`v2` or `v3`). The SDK defaults to
`v2`, except paths under `/trading/`, which default to `v3`. Override with
`client.WithAPIVersion` or `client.WithAPIVersionFor`.

## Rate limits

| Scope | Limit |
|-------|-------|
| Create / Check token | 10 req / 30s |
| Client token create / refresh (Display) | 600 req / min |
| Market data — general | 600 req / min |
| Stock tick / snapshot / quotes / bars | 60 req / 60s |
| Footprint | 600 req / min |
| Streaming subscribe / unsubscribe | 60 req / 60s |
| Order preview | 40 req / 10s |
| Order place / replace / cancel | 15 req/s (US), 1 req/s (HK / A-share) |
| Order open / history / detail | 40 req / 2s |
| MQTT connections | max 5 concurrent per App Key |
| MQTT throughput | ~3 messages/sec/connection |

The SDK additionally applies client-side retry, rate limiting and a circuit
breaker (`client.WithRetry`, `client.NewRateLimiter`, `client.NewBreaker`).

## Errors

Business failures return HTTP `417` with `{ "error_code", "message" }`;
`401` is unauthorized and `500` is a server error. `401` on Display Solution
means the client token expired and is refreshed automatically. The SDK maps all
of these onto typed `errs` codes — see [Errors](errors.md).

## Path status

Since v1.1.0 every SDK path follows the official OpenAPI definition —
[SDK ↔ API Reconciliation](reconciliation.md) reports **209 implemented
endpoints, 0 gaps, and 0 path discrepancies**. What remains unverified is live
behaviour in environments the HK sandbox cannot exercise:

| Area | Status |
|------|--------|
| Display Solution | Implemented; HK sandbox returns `403` (paid entitlement required). |
| Broker API HK | Implemented; HK sandbox returns `401 ROUTE_NOT_PERMITTED` (app scope missing). |
| Broker FD US, crypto, option-chain, some futures and event-contract data | Implemented; US-only — HK sandbox returns `404` / `417`. Verify with US sandbox credentials. |
| Auth and instrument paths | Webull's `llms.txt` summary and its own OpenAPI JSON disagree for a few endpoints; the SDK follows the summary path. `examples/path-probe` can confirm both against a live sandbox. |

Multi-leg option strategies, futures order validation and option-chain
discovery are fully implemented; the HK sandbox accepts only `SINGLE` orders
(`417` for every other strategy).

## Not implemented

None — every endpoint documented by Webull is implemented. See
[SDK ↔ API Reconciliation](reconciliation.md) for the per-endpoint mapping.

## Pages

- [Authentication](webull-api/authentication.md) — create/check token, client token.
- [Market Data — Stock](webull-api/market-data-stock.md)
- [Market Data — Option](webull-api/market-data-option.md)
- [Market Data — Crypto](webull-api/market-data-crypto.md)
- [Market Data — Futures](webull-api/market-data-futures.md)
- [Market Data — News](webull-api/market-data-news.md)
- [Market Data — Screener](webull-api/market-data-screener.md)
- [Market Data — Watchlist](webull-api/market-data-watchlist.md)
- [Fundamentals and Fund Data](webull-api/fundamentals.md)
- [Event Contracts](webull-api/event-contracts.md)
- [Trading API](webull-api/trading.md)
- [Broker API — HK](webull-api/broker-hk.md)
- [Broker API — FD (US)](webull-api/broker-fd-us.md)
- [Display Solution](webull-api/display-solution.md)
- [Streaming (MQTT)](webull-api/streaming.md)
- [Events (gRPC)](webull-api/events.md)
- [Connect API (OAuth)](webull-api/connect-api.md)

## Raw Webull data (verbatim)

Complete snapshots of Webull's own published documentation, with no
SDK-specific content:

- **[Master Guides (verbatim)](webull-api/master-guides.md)** — authentication,
  signature algorithm, token lifecycle, market-data and streaming guides,
  trading rules, broker/connect guides, error codes, FAQ and changelog.
- **[Master Reference (verbatim)](webull-api/master-reference.md)** — the
  OpenAPI definition of every documented endpoint.
- **[SDK ↔ API Reconciliation](reconciliation.md)** — for every documented
  endpoint, whether the SDK implements it and whether the SDK path matches
  the official definition.
