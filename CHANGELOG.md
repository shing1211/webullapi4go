# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-09-18

Initial public release.

### Added

- Core `client` package: `client.New`, the single signed-request entry point
  `Client.Do`, and functional options for credentials, region, environment,
  endpoints, transport, retries, rate limiting, circuit breaking, API version,
  and environment configuration (`WithEnv`).
- Authentication: HMAC-SHA1 request signing over a percent-encoded canonical
  string, the token lifecycle (`CreateToken`, `CheckToken`, `EnsureToken`,
  `CurrentToken`, `AccessToken`, `SetToken`), automatic token injection, and
  automatic sandbox token acquisition with `WithAutoToken`.
- `data` package: Market Data HTTP endpoints for instruments, company profile,
  analyst data, futures static data, snapshot, tick, quotes/depth, single and
  batch bars, footprint, NOII, screener, watchlist CRUD, options, and news.
- `stream` package: Market Data streaming over MQTT and MQTT-over-WebSocket with
  typed `Quote`, `Snapshot`, and `Tick` handlers, automatic reconnect, and
  idempotent automatic re-subscription.
- Generated protobuf types in `gen/webull/marketdata/v1` for streamed messages.
- `pkg/types` shared domain types.
- Resilient transport primitives: retry with exponential backoff, a keyed token
  bucket rate limiter, and a circuit breaker.
- Documentation site (MkDocs Material) with getting-started, authentication,
  market-data, streaming, sandbox, errors, API reference, and troubleshooting
  pages, plus Architecture Decision Records.
- Runnable examples under `examples/` for auth, market data, streaming, and
  watchlists.

[Unreleased]: https://github.com/shing1211/webullapi4go/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/shing1211/webullapi4go/releases/tag/v0.1.0
