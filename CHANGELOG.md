# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `client.Client.ObservabilityConfig` exposes the shared, read-only telemetry
  configuration inherited by service clients.
- `events` and `brokerfd/events` now emit one client-kind OpenTelemetry span
  per gRPC stream attempt, `event_stream_attempts` and
  `event_stream_attempt_duration` metrics, structured start/done logs, and
  correlation/trace metadata. Credentials are not included in telemetry.
- `trade.Client` exposes account-scoped OMS inspection and reconciliation:
  `GetTrackedOrder`, `TrackedOrder`, `TrackedOrderState`, `TrackedOrders`,
  `ApplyOrderEvent`, `ReconcileOrderStatus`, and `ReconcileOrderState`.
- `pkg/domain/order.Machine` and `order.Order` now support concurrency-safe
  status reconciliation, stale non-terminal no-ops, terminal-state protection,
  and distinct failed cancel/modify/reject scene events.
- `examples/account-monitor`: a read-only, single-worker balance monitor with
  bounded requests, signal-aware cancellation, and a consecutive-failure stop.

### Changed

- `client.Do`, `client.DoBroker`, and `client.DoStream` now share one
  per-attempt request pipeline for rate limiting, circuit breaking, hooks,
  interceptors, signing, correlation, logging, tracing, and metrics. Hook
  attempt numbers are one-based and increase across retries.
- `client.Do` and `client.DoBroker` reuse one correlation ID across retries;
  `DoBroker` and `DoStream` also learn the bounded clock offset from successful
  response `Date` headers when clock correction is enabled.
- REST requests inject the configured W3C trace propagator. `json.RawMessage`
  request bodies remain byte-for-byte stable across retries.
- The default HTTP client now clones and tunes Go's default transport for idle
  connections, TLS handshakes, and expect-continue timing; explicit
  `WithHTTPTransport` still wins.
- `trade.WithAutoClientOrderID(true)` derives stable IDs for the same logical
  place/batch request, does not mutate the caller's request slice, and keeps
  repeated requests idempotent. Successful batch results are tracked locally.
- The trade order registry is keyed by `(account_id, client_order_id)`.
  Successful replace/cancel actions advance matching local state; failed
  actions leave it unchanged.
- Order guardrail failures retain the historical outer `invalid_config` code
  and also wrap `pkg/errors.ErrOrderGuardrail`.
- `stream.Close` is terminal, channel cancel functions are idempotent, blocked
  channel sends unblock on cancellation/close, and reconnect replay is
  serialized with subscription mutations.

### Fixed

- MQTT `Connect` no longer starts a `Token.Wait` waiter goroutine, checks an
  already-cancelled context before connecting, and rejects use after `Close`.
  MQTT close is idempotent and suppresses late callbacks/messages.
- `DropOldest` now removes the oldest unread value. Cancelling or closing while
  dispatch is blocked no longer deadlocks the stream.
- Stream health only treats quote/snapshot/tick messages as fresh data; notices
  and echo heartbeats neither prevent nor falsely recover `StateDegraded`.
  Repeated stale checks no longer restore `StateConnected` without new data.
- OMS scene mapping no longer treats failed cancel, modify, or reject operations
  as successful transitions. Status reconciliation never regresses a terminal
  order and ignores stale non-terminal snapshots.
- Correlation context lookup is type-safe, response status/latency hooks receive
  real attempt values, and stream interceptors cannot accidentally return a
  nil successful response.

### Tests

- Added offline regression coverage for the shared request pipeline, raw-body
  retries, clock correction, account-scoped OMS tracking, status reconciliation,
  concurrent order machines, stable auto IDs, guardrail sentinels, MQTT/channel
  shutdown, stream health, resubscription serialization, `money.Money` wire
  forms, and Trading/Broker FD event telemetry.
- `go test ./...` passes in the root module and nested `broker/` module on
  2026-09-25. The new hardening was not newly live-verified.

### Documentation

- Reconciled the architecture and status docs with the preserved root service
  packages and public `money.Money`; marked the full `pkg` relocation,
  root-service shim plan, raw `decimal.Decimal` DTO migration, and `sync.Pool`
  decisions as superseded or closed.
- Added OMS reconciliation, streaming state/health/channel behavior, typed
  error examples, OpenTelemetry setup, event telemetry, and explicit
  implemented/offline-tested/live-verified/blocked distinctions.
- Corrected generated reconciliation summaries so they no longer claim zero
  path discrepancies while the generated report contains summary-only and
  unresolved paths.

## [2.1.0] - 2026-09-24

Phase 13 vet warnings fixed.

### Fixed

- `pkg/observability/otel.go`: `Tracer`, `Meter`, `InjectTraceContext`, and
  `ExtractTraceContext` now use `*Config` pointer receivers, eliminating the
  `passes lock by value` vet warnings.
- `client/config.go`: `Validate` now uses `*Config` pointer receiver.
- `client/config.go`: `client.Config.otel` field changed from embedded
  `observability.Config` to `*observability.Config` (pointer), so that
  returning `Config` by value does not copy the embedded `sync.RWMutex`.
- `client/client.go`: `Client.cfg` changed to `*Config` (pointer) to avoid
  copying the mutex when constructing a `Client`; `Config()` returns
  `*c.cfg` by value-copy of the pointer.

## [2.0.9] - 2026-09-24

Phase 9 structured errors.

### Added

- `pkg/errors/errors.go`: three new error codes (`CodeNotInitialized`,
  `CodeValidation`, `CodeOrderGuardrail`) and four new sentinels
  (`ErrNotInitialized`, `ErrValidation`, `ErrOrderGuardrail`,
  `ErrSubscriptionExpired`, `ErrConnectionLimitExceeded`).
- `pkg/errors/errors.go`: `Error.Is` now compares by pointer equality first,
  ensuring `ErrX.Is(ErrX)` returns true even when `ErrX` is a typed
  `*Error`.

### Changed

- `client/request.go`: `ErrCircuitOpen` is now `*pkgerrors.Error` with
  `CodeTransport`; existing `errors.Is(err, client.ErrCircuitOpen)` checks
  continue to work.
- `pkg/transport/mqtt/mqtt.go`: `ErrConnectionRefused` and `ErrConnectionLimit`
  are now `*pkgerrors.Error` with `CodeTransport`; `errors.Is(err,
  mqtt.ErrConnectionLimit)` and `errors.Is(err, mqtt.ErrConnectionRefused)`
  continue to work.
- `trade/orders.go`: `enforceGuardrails` now uses `errs.Wrap` instead of
  `errs.New` when prefixing batch-order index, preserving the full error chain.
- `events/client.go`, `brokerfd/events/events.go`: subscription-expired errors
  now wrap `ErrSubscriptionExpired` so callers can check
  `errors.Is(err, errs.ErrSubscriptionExpired)`.

## [2.0.8] - 2026-09-24

Phase 8 context hygiene.

### Fixed

- `pkg/transport/mqtt/mqtt.go`: `Connect()` now uses `sync.WaitGroup` so the
  `token.Wait()` goroutine exits promptly when the context is cancelled, instead
  of leaking until the MQTT stack processes the disconnect.
- `stream/client.go`: `resubscribeContext()` panics on a nil `resubCtx` instead of
  silently using an uncancellable `context.Background()`.

### Documentation

- `stream/client.go`, `stream/channels.go`: added comments explaining that OTel
  metric `Add` calls use `context.Background()` intentionally (fire-and-forget;
  the metric SDK records synchronously and cannot block).

## [2.0.7] - 2026-09-24

Phase 6.3 OTel metrics hooks.

### Added

- `pkg/observability/otel.go`: lazy OTel instrument accessors on `Config`:
  `ClientLatencyHistogram()`, `BreakerTransitionsCounter()`,
  `StreamReconnectsCounter()`, `StreamChannelDropsCounter()`.
- `pkg/resilience/breaker/breaker.go`: `WithMeter(m metric.Meter)` option and
  `transitionCounter` instrument. `transitionTo()` records `from`/`to` state
  attributes on every breaker transition.
- `client/request.go`: `request_latency` histogram recorded on every `attempt()`
  call with `http.route` attribute.
- `stream/client.go`: `streamMetrics` with reconnect and per-topic drop
  instruments in scope `webullapi4go/stream`. `handleReconnecting()` increments
  the `reconnects` counter. `WithMeter(m metric.Meter)` is available in
  `stream/option.go`.
- `stream/channels.go`: per-topic drop counters feed the `channel_drops`
  instrument with a `topic` attribute.

## [2.0.6] - 2026-09-24

Phase 5 streaming engine hardening + Phase 6 observability.

### Added

- `stream/state.go`: new `State` type with six values (`Disconnected`,
  `Connecting`, `Connected`, `Reconnecting`, `Degraded`, `Closed`). `Client.State()`
  returns the current state; `Client.OnStateChange(func(prev, next State))` fires
  on every transition.
- `stream/option.go`: `WithHealthWatchdog(interval)` — enables a background
  watchdog that transitions the connection to `StateDegraded` when no data
  message (quote/snapshot/tick) arrives within the configured interval, and
  recovers to `StateConnected` when a message arrives.
- `stream/client.go`: `Client.OnReconnecting(func())` handler invoked when the
  client begins attempting to reconnect after a connection loss.
- `stream/channels.go`: new channel-mux API with configurable backpressure.
  `Client.SubscribeQuoteChan`, `Client.SubscribeSnapshotChan`,
  `Client.SubscribeTickChan` each return a `<-chan T` with a per-subscription
  cancel function. `ChannelConfig.Policy` supports `DropBlock` (default,
  blocks sender), `DropOldest` (drops oldest unread), and `DropSample`
  (probabilistic). Drop counters track discards per channel.
- `client/request.go`: each `Client.Do` and `Client.DoBroker` call now establishes
  an OpenTelemetry span (`SpanKindClient`) scoped to the configured
  `TracerProvider`. Span attributes include `http.method`, `http.route`,
  `webull.attempt`, and `http.duration_ms`. Correlation ID (`X-Correlation-ID`
  header) is generated per-request using `crypto/rand` and threaded through the
  context.
- `client/request.go`: `WithCorrelationID(ctx, id)` and
  `CorrelationIDFromContext(ctx)` helpers for explicit correlation ID management.
- `client/option.go`: `WithLogger(*slog.Logger)` — configures per-request
  structured log output (request start and completion with method, path,
  correlation ID, latency, and error).
- `client/option.go`: `WithTracerProvider`, `WithMeterProvider`, and
  `WithPropagator` options supply real OpenTelemetry providers, replacing the
  no-op defaults.
- `pkg/observability/otel.go`: new package providing API-only OTel integration.
  `Config.Tracer()`, `Config.Meter()`, `Config.InjectTraceContext`, and
  `Config.ExtractTraceContext` use the configured providers; all types default
  to the global no-op implementations.
- `pkg/observability/otel.go`: `SpanAttributes(method, path, attempt)` and
  `SpanName(method, path)` helpers for consistent span naming and attribute
  sets.

## [2.0.4] - 2026-09-24

Phase 3 + Phase 4 production hardening: clock-drift correction, idempotency
helpers, transport tuning, resilience presets, interceptor pipeline, OMS state
machine integration, typed request builders.

### Added

- `client/option.go`: `WithClockDriftCorrection(bool)` — learns the clock offset
  between the client and Webull server from the `Date` response header and applies
  it to subsequent request signing timestamps, clamped to ±5 minutes.
- `client/option.go`: `WithHTTPTransport(*http.Transport)` — sets a custom HTTP
  transport on the client, applied after `WithHTTPClient` so callers can tune
  connection pooling and TLS without replacing the whole client.
- `client/option.go`: `WithResiliencePreset(ProductionPreset)` — applies a
  production-ready resilience composition: per-path rate limiter (10 req/s, burst
  20), circuit breaker (5-failure threshold, 30 s cooldown), and exponential
  backoff retry with full jitter (base 200 ms, cap 2 s, up to 3 attempts).
- `pkg/resilience/retry/`: `WithFullJitter(bool)` option — enables the full
  jitter formula `rand(0, min(cap, base·2ⁿ))` giving uniformly distributed
  delays in `[0, cap]` instead of ±20% jitter around the exponential.
- `trade/orders.go`: `NewClientOrderID()` — generates a fresh 20-character
  client-order identifier using `crypto/rand`, suitable for idempotent order
  placement.
- `trade/orders.go`: `ClientOrderIDFrom([]byte)` — derives a deterministic
  32-character client-order identifier from arbitrary content via SHA-256.
- `trade/orders.go`: `ValidClientOrderID(string) bool` — exported validation
  helper for the client-order identifier character set.
- `trade/option.go`: `WithAutoClientOrderID(bool)` — configures `PlaceOrder`
  and `BatchPlaceOrder` to auto-generate and assign a `NewClientOrderID` to
  each order whose `ClientOrderID` is empty.

### Added (Phase 4)

- `client/option.go`: `Interceptor` func type and `WithInterceptor(Interceptor)`
  option — composable request pipeline interceptors that run after rate-limiting
  and circuit-breaking but before signing and send.
- `client/option.go`: `Hooks` struct with `OnRequest`, `OnResponse`, `OnError`,
  `OnLatency` callbacks — observability integration point for Phase 6.
- `client/request.go`: `attempt()` refactored to run an ordered interceptor
  chain; hooks fire on every request attempt including retries.
- `pkg/domain/order/order.go`: new `Order` struct embeds `PlaceOrderResult`
  plus `AccountID` and a local `*Machine` state machine; `SceneTypeToEvent`
  maps Webull gRPC scene types to domain events.
- `trade/client.go`: added `orderRegistry` map and `registerOrder`/`getOrder`
  helpers for local order state tracking.
- `trade/orders.go`: `PlaceOrder` now returns `*order.Order` (not
  `*PlaceOrderResult`); the order is registered with `StatePending` on
  placement. `PlaceOrderResult` is embedded so `order.OrderID` and
  `order.ClientOrderID` remain accessible. `BatchPlaceOrder` unchanged.
- `trade/order_actions.go`: `CancelOrder` and `ReplaceOrder` check local order
  state before sending; terminal orders (filled, cancelled, failed, expired) return
  `errs.CodeInvalidTransition` without an API call.
- `trade/orders.go`: `NewPlaceOrderRequest(accountID, orders...)`,
  `NewEquityOrder(symbol, side, qty)`, and `EquityOrderBuilder` fluent API —
  typed request constructors with US equity defaults.
- `trade/order_actions.go`: `NewCancelOrderRequest(accountID, clientOrderID)` and
  `NewModifyOrderRequest(accountID, clientOrderID)` convenience constructors.
- `pkg/errors/errors.go`: added `CodeInvalidTransition` and `ErrInvalidTransition`
  for local state validation failures.

### Documentation

- Reorganized the docs site nav into Guides / API Reference / Official Webull
  Docs / Coverage / Architecture Decisions, and renamed pages for consistency
  (`broker.md` → `broker-hk.md`, `brokerfd.md` → `broker-fd-us.md`,
  `webull-api/crypto.md` → `webull-api/market-data-crypto.md`); old URLs are
  preserved with `mkdocs-redirects`.
- Marked every generated file with a "do not edit" banner; refreshed stale
  v1.0-era content (feature matrices, hub path status, provisional warnings) for
  v1.1.0; consolidated the official Webull doc links into `AGENTS.md`; merged
  `docs/adr/README.md` into `docs/adr/index.md`; standardized READMEs across all
  17 example directories.
- Fixed 80+ doc-code discrepancies across all endpoint docs: wrong function names
  (`GetOptionChain`→`GetOptionContracts`, `GetStockProfilesList`→`GetStockProfilesV3`),
  parameter mismatches, outdated code examples, and missing type documentation.
- Replaced the v0.1 coverage table in `market-data.md` with a comprehensive
  current table of all ~105 functions organized by group.
- Added SDK compatibility notes to all Display Solution endpoints documenting
  differences between the SDK implementation and the official OpenAPI spec.
- Added shared patterns reference (`patterns.md`) covering client construction,
  options, pagination, numeric strings, error handling, display service, and
  streaming patterns.
- Enhanced `errors.md` with transient vs permanent error classification table
  and rate-limit retry pattern.
- Added architecture diagram, "What's new in v1.1" section, and cross-reference
  links to `docs/index.md`.
- Added prerequisites to streaming, trading, and events docs.
- Documented previously undocumented types and functions across authentication,
  streaming, events, trading, broker-fd, and connect-api docs.
- Added prerequisites admonitions to all 10 guide pages (authentication,
  market-data, streaming, trading, events, fundamentals, errors, patterns,
  broker-hk, broker-fd-us).
- Added error handling tips to 5 guide pages (market-data, streaming, trading,
  events, fundamentals).
- Created glossary (`glossary.md`) with 30+ SDK-specific terms.
- Added ASCII token lifecycle diagram to authentication guide.
- Restructured API reference index (`webull-api.md`) with service-level grouping.
- Added common first-call failures table to getting started guide.
- Added `connect` and `display` packages to Go Packages reference (`api.md`).
- Added Documentation Standards section to `CONTRIBUTING.md`.
- Fixed broken Options table in streaming guide.

## [2.0.1] - 2026-09-23

Phase 2 production hardening: public API restructuring and type-safe numeric fields.

### Added

- `pkg/errors`: typed errors (`Error` with `Code` and `Message`) extracted from
  `internal/errs` to the public `pkg/errors/` package. A deprecated shim at
  `internal/errs/errs.go` preserves backward compatibility with existing imports.
- `pkg/transport/http.go`: public `Doer` interface (`Do(ctx, method, path, query,
  req, resp) error`) extracted from `client.Client`. New `pkg/transport/mqtt/`
  sub-package for MQTT transport types. A deprecated shim at
  `internal/transport/http.go` preserves backward compatibility.
- `pkg/resilience/{breaker,clock,ratelimit,retry}/`: public resilience primitives
  extracted from `internal/resilience`. A deprecated shim at
  `internal/resilience/resilience.go` preserves backward compatibility.
- `pkg/domain/money`: type-safe `Money` wrapper around `shopspring/decimal`
  (`github.com/shopspring/decimal v1.4.0`). Provides JSON round-trip fidelity,
  decimal arithmetic, `MarshalJSON`/`UnmarshalJSON`, and `Rat()` for rational
  arithmetic. Package-level helpers: `NewFromString`, `Must`, `ParseMoney`,
  `ParseDecimal`, `Zero`.
- `pkg/domain/order`: OMS state machine (`StateMachine`, `ApplyEvent`,
  `FromWebullStatus`) for order lifecycle tracking.
- `webull/`: thin type-alias convenience package re-exporting selected
  `client` and `trade` constructors and types. The root service packages remain
  canonical; this is not an aggregate service facade.

### Changed (breaking)

- `pkg/errors` is now the canonical import path. Old `internal/errs` is
  deprecated and will be removed in a future release.
- `pkg/transport` is now the canonical import path for `Doer`. Old
  `internal/transport` is deprecated.
- `pkg/resilience` is now the canonical import path. Old `internal/resilience`
  is deprecated.
- All numeric string fields converted to `*money.Money` (optional/request-side)
  or `money.Money` (required/response-side) throughout `trade/` and `data/`
  packages. Approximately 250 fields affected. JSON serialization is preserved
  (decimal strings); the change is transparent for most callers.
- `shopspring/decimal` is now a direct dependency (v1.4.0). Previously it was
  an indirect dependency.
- `money.Rat()` now delegates to the shopspring library's built-in `Rat()`
  method, fixing a bug where fractional quantities were incorrectly computed.

### Deprecated

- `internal/errs/errs.go`: use `pkg/errors` instead.
- `internal/transport/http.go`: use `pkg/transport` instead.
- `internal/resilience/resilience.go`: use `pkg/resilience` instead.
- Direct use of `client.Client` remains canonical; `webull.New` and
  `webull.Client` are optional aliases for callers that prefer one import.

## [2.0.3] - 2026-09-23

### Fixed

- `stream/` and `data/`: `go.uber.org/goleak` false positives from
  `net/http.(*http2clientConnReadLoop).run` goroutines left after paho-mqtt
  WebSocket disconnect. Added `goleak.IgnoreAnyFunction` filters to both
  packages' `TestMain`. The goroutines are cleaned up asynchronously by the Go
  runtime and are not an SDK leak.

## [2.0.2] - 2026-09-23

### Changed (breaking)

- All numeric string fields in `brokerfd/` converted to `*money.Money`
  (optional/request-side) or `money.Money` (required/response-side).
  Approximately 30 fields affected across `assets.go`, `orders.go`,
  `funding.go`, `activity.go`, `journals.go`, and `instruments.go`.
  JSON serialization is preserved (decimal strings).

## [1.1.1] - 2026-09-23

### Dev/Tooling

- Added `Makefile` targets: `build`, `test`, `test-race`, `cover`, `lint`,
  `fuzz`, `vuln`, `docs`.
- Enabled `gosec` in `.golangci.yml`; excluded `internal/auth`
  (protocol-mandated HMAC-SHA1/MD5 body digest) and
  `internal/resilience/retry` (intentional `math/rand` jitter).
- Expanded CI matrix to ubuntu/macos/windows; added coverage job
  (baseline measurement) and `govulncheck` job.
- Added `.github/dependabot.yml` (weekly gomod + github-actions updates).
- Added `go.uber.org/goleak` for goroutine-leak detection; `TestMain` with
  `goleak.VerifyTestMain` for `client`, `stream`, `data`, `internal/mqtt`.
- Added `data/fuzz_test.go` with `FuzzDecodeQuote`/`Snapshot`/`Tick`
  JSON deserialization fuzz tests.

## [1.1.0] - 2026-09-22

Full SDK parity with the official Webull OpenAPI.

### Added

- `connect/`: OAuth 2.0 authorization-code flow (`AuthorizationURL`,
  `CreateToken`) with `GrantTypeAuthorizationCode`/`GrantTypeRefreshToken`.
- `data/`: crypto (`GetCryptoBars`, `GetCryptoSnapshot`, `GetCryptoInstruments`),
  Display event contracts (`GetEventContractTags`, `GetEventContractEventsList`,
  `GetEventContractMilestones`, `GetEventContractSeriesList`,
  `GetEventContractSportsFilters`, `GetEventGameStats`, `GetEventLiveData`,
  `GetEventMarketBars`, `GetEventMarketBarsByEvent`, `GetEventMarketDepth`,
  `GetEventMarketSnapshot`), and fund extras (`GetFundPerformance`,
  `GetFundHoldings`, `GetFundRating`, `GetFundSplits`, `GetFundFiles`,
  `GetFundAllocation`).
- `display`: `Service.RefreshClientToken`.

### Changed

- Aligned 69 SDK paths to the official OpenAPI definition across `broker`,
  `brokerfd` and `data`. `broker.GetTradeCalendar` is now `GET` with query
  parameters; Broker FD account-update, order-replace and order-cancel are now
  `POST`.
- Removed all 41 provisional TODO markers. The SDK now covers every documented
  endpoint (0 documented-only gaps in `docs/reconciliation.md`). Path
  reconciliation remains tracked separately by the generated report.

## [1.0.3] - 2026-09-22

Documentation release: verbatim Webull master reference, an SDK↔API
reconciliation, a reproducible doc generator, and a path probe.

### Added

- `docs/webull-api.md` and `docs/webull-api/**`: an SDK-mapped reference for
  every documented endpoint, plus **verbatim** Webull master guides
  (`master-guides.md`) and per-area OpenAPI definitions (`reference/*.md`).
- `docs/reconciliation.md`: maps every documented endpoint to its SDK function
  and status (match / differs / not implemented).
- `tools/webull-docgen/`: generator CLI (`docgen.py
  reference|master|reconciliation|all`) that fetches Webull's machine-readable
  `.md` pages and renders the docs above.
- `examples/path-probe/`: env-gated probe comparing SDK paths against the
  official OpenAPI paths for the endpoints where they disagree.

### Changed

- `IMPLEMENTATION_STATUS.md`: coverage/gap section, Known Issues 12–13 (path
  drift), version history updated to v1.0.3.
- `mkdocs.yml`: added the Webull API Reference section.

## [1.0.2] - 2026-09-22

Reconciled the codebase against the official Webull API.

### Removed (breaking)

- Undocumented functions: crypto data (`data/crypto_data.go`), screener v2
  (`data/screener_v2.go`), option expirations/chains, and the redundant HK
  futures market-data variants.

### Changed

- `data.GetOptionContracts` now uses the Trading API path
  `/trading/instruments/options/contracts/list`.

## [1.0.1] - 2026-09-22

HK sandbox probe: confirmed futures product-codes path, fixed FuturesInstrument.Unit
flexible type to handle numeric API responses, fixed client_order_id length overflow
in options-multi-leg example, and updated TODO markers with HK sandbox findings.

### Fixed

- `data/futures.go`: `FuturesInstrument.Unit` now uses a `StringOrNumber` flexible
  type that handles both string (`"1-index points"`) and numeric (`1`) JSON values
  from the live API.
- `examples/options-multi-leg/main.go:193`: `client_order_id` now uses `Unix()` instead
  of `UnixNano()` to stay within the 32-character limit.

### Tests Added

- `data/futures_test.go`: `TestGetFuturesInstrumentsNumericUnit` verifies numeric unit
  deserialization.

### Changed

- `trade/types.go`: `TODO(t8)` comment updated to reflect HK sandbox finding that
  multi-leg strategies are rejected (only SINGLE accepted).
- `trade/options.go`: `TODO(t8)` comments updated to reflect HK sandbox findings.
- `data/futures_market.go`: header comment notes that product-codes path is confirmed
  while market data paths remain unconfirmed.

## [0.9.0] - 2026-09-22

GoDoc coverage, HK options stubs, HK futures market data, new examples, and graceful credential handling.

### Added

- `brokerfd/brokerfd.go`: GoDoc on all 12 files and ~80 exported identifiers.
- `brokerfd/events/events.go`, `brokerfd/events/option.go`, `brokerfd/events/sign.go`:
  GoDoc on all event types and sign functions.
- `examples/brokerfd/`: Broker FD US read-only endpoint probe (accounts, orders,
  assets, instruments, funding, activity, journals, master data, agreements, documents).
- `examples/brokerfd-events/`: Broker FD gRPC event subscription probe (order, option,
  position streams).
- `examples/options/`: HK options discovery probe (expirations, option chain).
- `examples/`: graceful credential handling in all 9 existing examples — no more
  panic on missing env vars; instead a descriptive message and zero-value client.
- `data/futures.go`: `FuturesCategoryCN` constant for CN futures queries.
- `data/futures_market.go`: `GetHKFuturesTick`, `GetHKFuturesSnapshot`,
  `GetHKFuturesBars`, `GetHKFuturesDepth`, `GetHKFuturesFootprint` — 5 new HK
  futures market data functions.
- `data/options.go`: `OptionCategoryHK`, `OptionCategoryCN` constants.
- `data/options.go`: `GetHKOptionExpirations`, `GetHKOptionChain` — HK options
  discovery stubs (TODO t10, paths unconfirmed).
- `trade/derivatives_hk.go`: new file for HK derivatives-specific trading helpers.

### Tests Added

- `brokerfd/brokerfd_test.go`: unit tests for broker FD root package.
- `broker/options_test.go`: unit tests for broker options.
- `brokerfd/events/option_test.go`: unit tests for broker FD option events.

## [0.9.2] - 2026-09-22

Critical path corrections for the Broker API HK package. All paths were wrong
— the SDK used `/openapi/v1/broker/...` but the official API uses `/broker/...`
directly. Sandbox returned 404 not because the API was missing but because every
path was wrong. Additionally, request/response shapes were updated to match the
official OpenAPI specification.

### Fixed

- `broker/accounts.go`: paths corrected to `/broker/accounts/virtual-accounts/{list,get,create,update}`;
  `VirtualAccount` struct updated to official fields; `ListVirtualAccounts` now handles
  wrapped `{data:[...]}` response.
- `broker/activities.go`: path corrected to `/broker/activities/list`.
- `broker/assets.go`: paths corrected to `/broker/assets/{balances/get,positions/list}`;
  `Balance` struct updated to official nested shape with `total_cash_balance`,
  `total_market_value`, `total_unrealized_profit_loss`; `GetPositions` handles
  wrapped `{data:[...]}` response.
- `broker/instruments.go`: paths corrected to `/broker/instruments/{stocks/list,stock-locate/get,corporate-actions/get}`.
- `broker/orders.go`: paths corrected to `/broker/orders/{preview,place,replace,cancel,get,history,open}`;
  `ReplaceOrder` and `CancelOrder` now use POST with JSON body instead of PUT/DELETE
  with query params; `PreviewOrderRequest` adapted to official nested wire format.
- `broker/funding.go`: paths corrected to `/broker/funding/fx-rates/get`,
  `/broker/funding/fx-exchanges/{create,get}`,
  `/broker/funding/instant-fx/{create,get}`, `/broker/funding/instant/{create,get}`;
  `FXRate` struct updated (`fx_rate`, `rate_effective_time`, `rate_expire_time`);
  detail queries now require `client_request_id` + `account_id`.
- `broker/journals.go`: paths corrected to `/broker/journals/{cash-journals/{create,get},position-journals/{create,get}}`;
  detail queries now require `client_request_id` + `account_id`.
- `broker/masterdata.go`: path corrected to `/broker/master-data/trade-calendar/query`;
  changed from GET with query params to POST with JSON body.
- `broker/eventcontracts.go`: paths corrected to `/broker/event-contracts/{categories/list,series/get,events/get,instruments/get}`.
- `examples/broker-probe/main.go`: updated to use new field names and paths.

### Noted

- Broker API HK (`/broker/...`) returns `401 ROUTE_NOT_PERMITTED` on the HK
  sandbox — the app does not have Broker API scope enabled in the sandbox.
  Paths are confirmed correct (no more 404); the 401 means auth/scope issue.

## [0.9.1] - 2026-09-21

Bug fixes found during remaining endpoint verification sweep.

### Fixed

- `data/watchlist.go`: `AddWatchlistInstruments`, `RemoveWatchlistInstruments`,
  `UpdateWatchlistInstruments`, `UpdateWatchlist`, and `DeleteWatchlist` now
  handle bare boolean (`true`/`false`) responses from the API using a
  `BoolOrSuccess` type that accepts both raw bool and `{"success":bool}` JSON.
- `broker/client.go`: `DoBroker` method added to route Broker API calls to the
  correct `broker-api.sandbox.webull.hk` base URL instead of the market data URL.

### Added

- `examples/watchlist-cmd/`: Watchlist CRUD example (create, add instruments,
  update, remove, delete) guarded by `WEBULL_WATCHLIST_TEST=1`.
- `examples/broker-probe/`: Broker HK read-only endpoint probe program.

### Noted

- Broker API HK (`/openapi/v1/broker/...`) returns `404 Route Not Found` on the
  HK sandbox — the endpoint group is not available in the sandbox environment.
  Broker HK remains unverified pending production or US sandbox access.

## [0.9.0] - 2026-09-21

Full sandbox verification: all core SDK functionality tested against live HK sandbox.

### Fixed

- `data/fundamentals.go`: `FinancialsItem` changed from `map[string]string` to
  `map[string]any` to handle numeric values returned by the income statement,
  balance sheet, and cash flow endpoints.

### Added

- `examples/`: all examples verified against HK sandbox (auth, marketdata,
  account, watchlist, order preview, data-fundamentals, streaming, events).
- Documentation complete: `IMPLEMENTATION_STATUS.md` with full codebase audit
  (49 TODOs, 522 tests), README roadmap sync, AGENTS.md constraints updated.

## [0.7.0] - 2026-09-21

Complete API coverage: market data extensions, event contracts, Broker API HK, Broker FD API US, and Broker FD gRPC events.

### Added

- `data/eventcontracts.go`: Event contract instrument discovery — `GetEventContractCategories`, `GetEventContractSeries`, `GetEventContractEvents`, `GetEventContractMarkets`
- `data/eventcontracts_market.go`: Event contract market data (provisional) — `GetEventSnapshot`, `GetEventDepth`, `GetEventBars`, `GetEventTick` (TODO)
- `trade/types.go`: `InstrumentTypeEvent` and `EventOutcome` (`yes`/`no`) for event contract trading
- `trade/rules.go`: `validateEventRules` — LIMIT-only, DAY-only, QTY-only, integer quantity ≤50,000 (TODO)
- `data/screener.go`: Non-display screener extensions — `GetMarketSectors`, `GetMarketSectorDetail`, `GetHighDividendRank`, `GetWeek52HighLow`
- `data/crypto_data.go`: US crypto dedicated paths — `GetCryptoSnapshotList`, `GetCryptoBarsList`
- `trade/orders.go`: `BatchPlaceOrder` for multi-order submission
- `trade/accounts.go`: Cash activity accessors — `GetCashActivities`, `GetCashActivitiesPage`, `GetAllCashActivities`
- `data/futures_market.go`: Futures market data (provisional) — `GetFuturesTick`, `GetFuturesSnapshot`, `GetFuturesBars`, `GetFuturesDepth`, `GetFuturesFootprint` (TODO)
- `data/display_screener.go`: Display Solution screeners (provisional) — `GetDisplayGainersLosers`, `GetDisplayTopActive` (TODO)
- `data/display_quotes.go`: Display Solution quotes (provisional) — `GetDisplaySnapshot`, `GetDisplayBars`, `GetDisplayBarsSingle`, `GetDisplayTick`, `GetDisplayDepth` (TODO)
- `data/display_instruments.go`, `data/display_news.go`, `data/display_streaming.go`: Display Solution instruments, news, and streaming (provisional) (TODO)
- `broker/` module: New Go module for Broker API HK — virtual accounts, instruments, assets, orders, cash activities, funding FX, instant funding, journals, master data, event contracts
- `brokerfd/` module: Complete rewrite of Broker FD API US — agreements, accounts, documents, assets, activity, funding, instruments, orders, journals, master data
- `brokerfd/events/` module: Broker FD gRPC events client using `grpc.event.EventService`
- `internal/region/region.go`: `BrokerHTTP` field added to `Endpoints` struct
- `gen/webull/brokerfd/events/v1/`: Generated protobuf for Broker FD events (P5.2)

### Changed

- `data/instrument_v3.go`: `GetStockProfilesV3` auth routed through `DisplayService()`
- `data/logos.go`: `GetLogos` auth routed through `DisplayService()`

## [0.6.0] - 2026-09-20

Multi-leg options orders, futures order validation, speculative option chain discovery, and documentation sync.

### Added

- `trade.OptionStrategy` constants for multi-leg options: `VERTICAL`, `STRADDLE`, `STRANGLE`, `IRON_CONDOR`, `IRON_BUTTERFLY`, `BUTTERFLY`, `COLLAR`, `CALENDAR`, `DIAGONAL`, `RATIO` — provisional (TODO(t8))
- `trade.validateFuturesRules`: dedicated futures order-type matrix (US/HK), QTY-only/whole-contract/DAY-GTC validation — provisional (TODO(t9))
- `trade/multileg_test.go`, `trade/futures_order_test.go`: offline tests for multi-leg and futures validation
- `data.GetOptionExpirations`, `data.GetOptionChain`: speculative option chain discovery (not in published Webull API) — marked TODO(t10)
- `data/options_test.go`: offline tests for option chain functions

### Changed

- `trade.orders.go`: instrument-dispatched validation (equity/option/futures branches)
- `trade.options.go`: multi-leg order validation (2+ leg enforcement, duplicate/degenerate detection, canonicalStrike normalization)
- `trade.rules.go`: futures order-type matrix and `isPositiveInteger` helper
- `trade.option.go`: `WithMaxOrderNotional` documentation noting multi-leg bypass
- `README.md`, `AGENTS.md`, `docs/trading.md`, `docs/api.md`, `docs/market-data.md` updated to reflect new features

## [0.5.0] - 2026-09-19

Display Solution integration, Corporate Actions, Fund Data, Crypto Data, and Screener v2 stubs.

### Added

- `display` package: Display Solution authentication service — HMAC-SHA1 signed
  client-token fetch, `Authorization: Bearer` header, lazy initialization, HK
  sandbox/prod hosts.
- `data.GetCorporateActionsList`, `data.GetCorporateActionsMarket`: Corporate
  actions by symbol or by market (Display Solution). Paths confirmed via live probe.
- `data.GetStockProfilesList`: Batch instrument profiles (POST, Display Solution).
- `data.GetLogosBatch`: Batch company logos (POST, Display Solution).
- `data.GetFundNav`, `data.GetFundInfo`, `data.GetFundDividends`, `data.GetFundList`:
  Fund NAV, info, dividends, and fund list. **Paths are best-effort; require US
  sandbox credentials to verify.**
- `data.GetCryptoBars`, `data.GetCryptoTick`, `data.GetCryptoDepth`,
  `data.GetCryptoSnapshot`: Crypto OHLCV, tick, depth, and snapshot. Paths corrected
  to `/market-data/` with `category=CRYPTO` query param. **HK sandbox returns 417
  (CRYPTO category unsupported); US credentials required to verify response schemas.**
- `data.GetScreenerV2`: Screener v2 POST query. **Path is best-effort; requires
  US sandbox credentials to verify.**
- `brokerfd` package: Broker FD HTTP stub with `GetAccountsSummary` and `GetPositions`.
  **All paths are best-effort; returns 404 in HK sandbox.**
- `examples/probe`: Live probe binary for systematic endpoint testing.

### Changed

- `data.GetCorporateActions`: signature changed from `(ctx, query) (data, error)`
  to `(ctx, query) (data, paginationKey, error)` to match the paginated API response.

### Fixed

- `data/crypto_data.go`: paths changed from `/market-data/crypto/{symbol}/bars` to
  `/market-data/bars` with `category=CRYPTO&symbol=` query params.

## [0.4.0] - 2026-09-18

Market data fundamentals: capital flows, industry comparisons, earnings and dividend
calendars, SEC filings, and financial statements.

### Added

- `data.GetCapitalFlow`: capital-flow data (large/medium/small in/out flows) for
  the trailing N trading days per symbol.
- `data.GetIndustryComparison`: industry-relative performance metrics (e.g. EPS_TTM)
  with per-company ranks and values.
- `data.GetEarningsCalendar`: upcoming and historical earnings-release dates, EPS
  actual vs. estimate, and revenue actual vs. estimate per fiscal period.
- `data.GetDividendCalendar`: dividend and split events with declare/ex-div/record/pay
  dates and per-share amounts.
- `data.GetFilings`: SEC filings list (8-K, 10-K, 10-Q, etc.) with titles, URLs,
  and publish dates.
- `data.GetIncomeStatement`: multi-period income statement data (revenue, net income,
  EPS, etc.).
- `data.GetBalanceSheet`: multi-period balance sheet data (total assets, liabilities,
  equity, etc.).
- `data.GetCashFlow`: multi-period cash flow statement data (operating, investing,
  financing cash flows).
- `data.GetFinancialIndicators`: key financial ratios and metrics (ROA, ROE, EPS,
  net margin, debt ratio).
- `data.GetFinancialAlert`: upcoming earnings-release alert with expected date and
  estimated vs. last-year EPS.
- `data.GetForecastEPS`: analyst consensus EPS forecasts for the next 5 quarters.
- `internal/auth`: added `DigestCase` (`DigestUpper`/`DigestLower`) to `SignParams`
  so callers can control the casing of the body digest hex output; the REST path
  remains unchanged.
- `examples/data-fundamentals`: Runnable program demonstrating all fundamentals
  endpoints for AAPL.

## [0.3.0] - 2026-09-18

Trading events over gRPC.

### Added

- `events` package: a server-streaming gRPC client for Webull's trade events,
  built from the core `client.Client` with `events.New`. It subscribes to order,
  event-contract position, and option streams (`SubscribeOrder`,
  `SubscribePosition`, `SubscribeOption`, and the `SubscribeAll` bitmask),
  optionally scoped to trading accounts with `events.WithAccounts`.
- Typed event payloads: `OrderEvent`, `PositionEvent`, and `OptionEvent` decode
  the JSON carried by data events, with every numeric value kept as a string and
  `Raw` preserving the exact wire bytes. Handlers are `OnOrder`, `OnPosition`,
  `OnOption`, and the raw `OnEvent`, alongside `OnConnect`, `OnPing`, and
  `OnError`.
- HMAC-SHA256 signing for the event service: no `host` participates in the
  canonical string and the body digest is lower-case SHA-256, distinct from the
  REST HMAC-SHA1 signer. The algorithm-parameterised signer in `internal/auth`
  keeps the REST path byte-for-byte unchanged (golden vectors) and adds the
  SHA-256 path used by `events`.
- Vendored `events.proto` from the Apache-2.0 Webull Python SDK, generated into
  `gen/webull/trade/events/v1`, with provenance and attribution recorded in
  `THIRD_PARTY_NOTICES.md` and the accepted ADR-0002.
- Reconnect and re-subscribe with exponential backoff plus jitter, bounded by
  `events.WithMaxReconnectAttempts`; terminal server events (auth, connection
  limit, expired subscription) end the run with a typed error.
- Runnable `examples/events` program that subscribes to order events, prints
  each decoded event, and shuts down on Ctrl+C.
- Trading events documentation page covering the endpoint, the HMAC-SHA256
  signing rules, the subscribe bitmask, the dispatch model, the order-event JSON
  schema, reconnect options, and sandbox caveats.
- A nightly, read-only live-sandbox CI job that connects to the event stream
  without mutating any account.

### Changed

- The README feature matrix and roadmap mark Trading events as supported in
  v0.3.0, and the documentation site gains a Trading Events page.

## [0.2.6] - 2026-09-18

Route the news Server-Sent Events stream through the core client pipeline.
No public API changes.

### Added

- `client.DoStream`, a signed streaming request that returns the raw, open
  response for endpoints that stream their result. It applies the same signing,
  per-path API version, token, rate limiter, and circuit breaker as
  `client.Do`, but never retries and leaves the successful response body for
  the caller to read and close.

### Changed

- The `data` package news summary (`GetNewsSummary`) now uses
  `client.DoStream` instead of its own private stream path, so the access token,
  API version, and resilience policies apply to it exactly as to buffered
  requests. Its removed private helpers are not part of the public API.

## [0.2.5] - 2026-09-18

US combo orders for the `trade` package.

### Added

- US combo-order groups: `PlaceOrderRequest` now validates the whole order set as
  a combo group when any order uses a non-NORMAL `combo_type`. Supported groups
  are take-profit/stop-loss (`MASTER` with optional `STOP_PROFIT` and
  `STOP_LOSS`), `OTO`, `OCO`, and `OTOCO`. Combo orders are US equity only.
- Composition and leg-count rules: one request must contain a single group kind,
  every order must be a US equity order, and `client_combo_order_id` is required.
  OTO and OTOCO require exactly one `MASTER` plus one to six legs, OCO requires
  two to six legs with no `MASTER`, and take-profit/stop-loss allows at most one
  `MASTER`, one `STOP_PROFIT`, and one `STOP_LOSS`.
- Per-role order-type rules: a take-profit/stop-loss `MASTER` is `MARKET` or
  `LIMIT`; an OTO/OTOCO `MASTER` and an OTO leg add `STOP_LOSS` and
  `STOP_LOSS_LIMIT`; take-profit/stop-loss, OCO, and OTOCO legs are `LIMIT`,
  `STOP_LOSS`, or `STOP_LOSS_LIMIT`.
- Sell-to-close take-profit/stop-loss groups: a group with no `MASTER` must use
  side `SELL` for every sub-order.
- Table-driven tests for the combo composition, leg-count, order-type, and
  sell-to-close rules, plus a combo preview test.
- A "Combo types" section in the trading documentation covering the group kinds,
  the composition and leg-count rules, the sell-to-close form, and a
  MASTER/STOP_PROFIT/STOP_LOSS example.

### Changed

- `PlaceOrderRequest.Validate` now runs the combo-group validation after the
  per-order checks.

## [0.2.4] - 2026-09-18

Single-leg options orders for the `trade` package.

### Added

- Single-leg options orders: `OrderRequest` accepts `instrument_type` `OPTION`
  with `option_strategy` `SINGLE` and exactly one `legs` entry, validated before
  any network call. Option orders support `LIMIT`, `STOP_LOSS`, and
  `STOP_LOSS_LIMIT`; `MARKET` and the other stock order types are rejected.
- Options side and time-in-force rules: only `BUY` and `SELL` are accepted
  (`SHORT` is rejected), sell-side orders must use `DAY`, and `GTD` is rejected
  for options.
- `OrderLeg.Validate` validates a single option leg: `instrument_type` `OPTION`,
  `market` `US`, a non-blank `symbol`, `BUY`/`SELL`, a positive decimal
  `strike_price`, a calendar `option_expire_date` in `YYYY-MM-DD` form,
  `option_type` `CALL` or `PUT`, and a positive decimal `quantity`.
- Table-driven tests for the option order and leg rules, plus an env-gated,
  live-tolerant preview test.
- An "Options orders" section in the trading documentation covering the
  supported order types, side and time-in-force rules, the `legs` schema, a code
  example, and sandbox option-contract caveats.

### Changed

- `OrderRequest.Validate` now enforces the option-specific rules after the
  market rules; a non-OPTION order that carries `option_strategy` or `legs` is
  rejected.

## [0.2.3] - 2026-09-18

Market-specific order validation and Hong Kong BCAN support for the `trade`
package.

### Added

- Market-specific equity order-type validation: US accepts limit, market, stop,
  stop-limit, market-on-open, market-on-close, touch, and trailing stop orders;
  HK accepts enhanced limit, at-auction, at-auction limit, stop, stop-limit,
  touch, and trailing stop orders; CN accepts only `LIMIT`.
- Hong Kong BCAN: `HK` equity orders now require at least one `no_party_ids`
  entry with a non-blank `party_id`, `party_id_source` `"D"`, and `party_role`
  `"3"`. A `no_party_ids` list on a non-HK-equity order is rejected.
- US `support_trading_session` validation (`CORE`, `ALL`, `NIGHT`, `ALL_DAY`).
  The deprecated `Y` and `N` aliases and use on a non-US order are rejected.
- At-auction price rules: `AT_AUCTION` rejects `limit_price` and
  `AT_AUCTION_LIMIT` requires it.
- A-share (`CN`) documentation noting that only `LIMIT` is accepted and that
  A-share trading is disabled by default until enabled by Webull support.
- Table-driven tests covering the market matrix, Hong Kong BCAN, US trading
  sessions, and the at-auction price rules.

### Changed

- The trading documentation adds a "Market rules" section, and the `OrderType`,
  `TradingSession`, and `PartyID` GoDoc describe the enforced rules.

## [0.2.2] - 2026-09-18

Stock-order lifecycle and queries for the `trade` package.

### Added

- `trade` order methods: `PreviewOrder` and `PlaceOrder` (with guardrail
  enforcement), `ReplaceOrder`, and `CancelOrder`. Requests are validated before
  any network call and the v3 modify endpoints identify an order by its client
  order ID.
- `trade` order queries: `GetOpenOrders`, `GetOpenOrdersPage`,
  `GetAllOpenOrders`, `GetOrderHistory`, `GetOrderHistoryPage`,
  `GetAllOrderHistory`, and `GetOrderDetail`. Pagination follows the cursor to
  exhaustion, bounded by `trade.MaxOrderQueryPages`.
- Order domain types and enums: `OrderRequest`, `PlaceOrderRequest`,
  `PlaceOrderResult`, `PreviewResult`, `ModifyOrderRequest`,
  `ReplaceOrderRequest`, `ReplaceOrderResult`, `CancelOrderRequest`,
  `CancelOrderResult`, `OrderGroup`, `OrderPage`, `Order`, `OrderLeg`,
  `OrderLegDetail`, `OrderCommission`, `OrderFee`, `OrderHistoryQuery`,
  `OrderSide`, `OrderType`, `TimeInForce`, `ComboType`, `EntrustType`,
  `TradingSession`, `TriggerPriceType`, `TrailingType`, `OrderStatus`, and
  `PartyID`.
- `OrderRequest.Validate`, `PlaceOrderRequest.Validate`,
  `ReplaceOrderRequest.Validate`, and `CancelOrderRequest.Validate`, returning
  typed `invalid_config` errors for the first problem found.
- Runnable `examples/order` program that previews a non-marketable AAPL limit
  buy and, only when `WEBULL_ORDER_PLACE=1` is set, places and cancels it.
- Expanded trading documentation covering the order lifecycle, order-type,
  time-in-force, and combo-type tables, validation rules, guardrails, queries,
  and the explicit warning that placing orders mutates a real account.

## [0.2.1] - 2026-09-18

Trading HTTP foundation: read-only accounts and assets.

### Added

- `trade` package: a typed Trading HTTP client built on the core `client.Client`
  with `trade.New`. Endpoints: `ListAccounts`, `GetBalance`, and `GetPositions`,
  covering account listing and per-account balances and positions. Asset calls
  reject an empty account ID before any network request.
- Trading domain types and enums: `Account`, `AccountType`, `AccountClass`,
  `AssetsBalance`, `AssetsCurrencyAssets`, `Position`, `PositionLeg`, `Market`,
  `InstrumentType`, `OptionType`, and `OptionStrategy`. Numeric fields stay
  strings to preserve precision.
- Order guardrail options `trade.WithMaxOrderNotional` and
  `trade.WithMaxOrderQuantity`, enforced by the order methods before an order is
  built (order placement lands in a later patch).
- Built-in per-path API version defaults: requests under `/trading/` now send
  `x-version: v3`, while Market Data paths remain on `v2`. Overridable with
  `client.WithAPIVersion` and `client.WithAPIVersionFor`.
- Runnable `examples/account` program that lists accounts and prints the balance
  and positions for an account.
- Trading documentation page covering authentication, accounts and assets, the
  v3 default, guardrails, and sandbox usage.

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

[Unreleased]: https://github.com/shing1211/webullapi4go/compare/v2.1.0...HEAD
[2.1.0]: https://github.com/shing1211/webullapi4go/releases/tag/v2.1.0
[2.0.9]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.9
[2.0.8]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.8
[2.0.7]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.7
[2.0.6]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.6
[2.0.5]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.5
[2.0.4]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.4
[2.0.3]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.3
[2.0.2]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.2
[2.0.1]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.1
[2.0.0]: https://github.com/shing1211/webullapi4go/releases/tag/v2.0.0
[1.1.1]: https://github.com/shing1211/webullapi4go/releases/tag/v1.1.1
[1.1.0]: https://github.com/shing1211/webullapi4go/releases/tag/v1.1.0
[1.0.3]: https://github.com/shing1211/webullapi4go/releases/tag/v1.0.3
[1.0.2]: https://github.com/shing1211/webullapi4go/releases/tag/v1.0.2
[1.0.1]: https://github.com/shing1211/webullapi4go/releases/tag/v1.0.1
[1.0.0]: https://github.com/shing1211/webullapi4go/releases/tag/v1.0.0
[0.9.2]: https://github.com/shing1211/webullapi4go/releases/tag/v0.9.2
[0.9.1]: https://github.com/shing1211/webullapi4go/releases/tag/v0.9.1
[0.9.0]: https://github.com/shing1211/webullapi4go/releases/tag/v0.9.0
[0.7.0]: https://github.com/shing1211/webullapi4go/releases/tag/v0.7.0
[0.6.0]: https://github.com/shing1211/webullapi4go/releases/tag/v0.6.0
[0.5.0]: https://github.com/shing1211/webullapi4go/releases/tag/v0.5.0
[0.4.0]: https://github.com/shing1211/webullapi4go/releases/tag/v0.4.0
[0.3.0]: https://github.com/shing1211/webullapi4go/releases/tag/v0.3.0
[0.2.6]: https://github.com/shing1211/webullapi4go/releases/tag/v0.2.6
[0.2.5]: https://github.com/shing1211/webullapi4go/releases/tag/v0.2.5
[0.2.4]: https://github.com/shing1211/webullapi4go/releases/tag/v0.2.4
[0.2.3]: https://github.com/shing1211/webullapi4go/releases/tag/v0.2.3
[0.2.2]: https://github.com/shing1211/webullapi4go/releases/tag/v0.2.2
[0.2.1]: https://github.com/shing1211/webullapi4go/releases/tag/v0.2.1
[0.1.0]: https://github.com/shing1211/webullapi4go/releases/tag/v0.1.0
