# API Reference

The pkg.go.dev pages are authoritative for the complete exported API. This page
maps the current public architecture and highlights behavior that is easy to
miss.

## Package map

| Package | Purpose |
|---|---|
| [`client`](https://pkg.go.dev/github.com/shing1211/webullapi4go/client) | Core configuration, signing, tokens, HTTP transport, request pipeline, and resilience |
| [`data`](https://pkg.go.dev/github.com/shing1211/webullapi4go/data) | Market Data HTTP endpoints and DTOs |
| [`stream`](https://pkg.go.dev/github.com/shing1211/webullapi4go/stream) | Market Data MQTT streaming, reconnect, health, and channels |
| [`trade`](https://pkg.go.dev/github.com/shing1211/webullapi4go/trade) | Trading HTTP endpoints and local OMS integration |
| [`events`](https://pkg.go.dev/github.com/shing1211/webullapi4go/events) | Trading events over gRPC |
| [`connect`](https://pkg.go.dev/github.com/shing1211/webullapi4go/connect) | OAuth 2.0 authorization-code flow |
| [`display`](https://pkg.go.dev/github.com/shing1211/webullapi4go/display) | Display Solution client-token service |
| [`broker`](https://pkg.go.dev/github.com/shing1211/webullapi4go/broker) | Broker API HK; separate Go module |
| [`brokerfd`](https://pkg.go.dev/github.com/shing1211/webullapi4go/brokerfd) | Broker FD US HTTP endpoints |
| [`brokerfd/events`](https://pkg.go.dev/github.com/shing1211/webullapi4go/brokerfd/events) | Broker FD events over gRPC |
| [`webull`](https://pkg.go.dev/github.com/shing1211/webullapi4go/webull) | Optional thin aliases for core client types and options |
| [`pkg/errors`](https://pkg.go.dev/github.com/shing1211/webullapi4go/pkg/errors) | Typed errors, stable codes, and sentinels |
| [`pkg/observability`](https://pkg.go.dev/github.com/shing1211/webullapi4go/pkg/observability) | OTel handles, helpers, and shared instruments |
| [`pkg/resilience`](https://pkg.go.dev/github.com/shing1211/webullapi4go/pkg/resilience) | Retry, rate limit, circuit breaker, and clock primitives |
| [`pkg/transport`](https://pkg.go.dev/github.com/shing1211/webullapi4go/pkg/transport) | Public HTTP and MQTT transport |
| [`pkg/domain/money`](https://pkg.go.dev/github.com/shing1211/webullapi4go/pkg/domain/money) | Public `Money` decimal wrapper |
| [`pkg/domain/order`](https://pkg.go.dev/github.com/shing1211/webullapi4go/pkg/domain/order) | Concurrency-safe order state machine and status reconciliation |
| [`pkg/types`](https://pkg.go.dev/github.com/shing1211/webullapi4go/pkg/types) | Shared public market and instrument types |
| [`gen/webull`](https://pkg.go.dev/github.com/shing1211/webullapi4go/gen/webull) | Committed generated protobuf packages |

The root service packages are canonical. `pkg/` contains shared foundations,
not relocated services. The `webull` package is a convenience alias layer; it
does not replace `data.New`, `trade.New`, `stream.New`, or `events.New`.

## `client`

Construct one core client and share it with service packages.

Core entry points:

- `client.New(opts ...Option) (*Client, error)`
- `Client.Do(ctx, method, path, body, out) error`
- `Client.DoBroker(ctx, method, path, body, out) error`
- `Client.DoStream(ctx, method, path, body) (*http.Response, error)`
- `Client.Close() error`

`Do` and `DoBroker` use the configured retry policy. Only GET and HEAD are
retried unless `RetryConfig.RetryNonIdempotent` is enabled. `DoStream` never
retries because a successful streaming body is owned by the caller and cannot
be replayed safely.

All three entry points share rate limiting, circuit breaking, interceptors,
hooks, logging, tracing, metrics, signing, token injection, correlation IDs,
and W3C trace-context injection. Hook attempt numbers are one-based per logical
request and increment across retries.

Configuration options include:

- Credentials and environment: `WithAppKey`, `WithAppSecret`,
  `WithCredentials`, `WithRegion`, `WithEnvironment`, `WithSandbox`, `WithEnv`
- Endpoints and HTTP: `WithEndpoints`, `WithBaseURL`, `WithHTTPClient`,
  `WithHTTPTransport`, `WithTimeout`, `WithUserAgent`
- Resilience: `WithRetry`, `WithoutRetry`, `WithRateLimiter`,
  `NewRateLimiter`, `WithBreaker`, `NewBreaker`,
  `WithResiliencePreset`, `WithClockDriftCorrection`
- Pipeline: `WithInterceptor`, `Hooks`, `WithHooks`
- Versions and tokens: `WithAPIVersion`, `WithAPIVersionFor`, `WithAutoToken`
- Observability: `WithLogger`, `WithTracerProvider`, `WithMeterProvider`,
  `WithPropagator`

Accessors include `Config`, `ObservabilityConfig`, `Region`, `Environment`,
`Endpoints`, `HTTPClient`, `AppKey`, and `AppSecret`. `ObservabilityConfig`
returns a shared pointer that service clients inherit and callers must treat as
read-only.

Token lifecycle:

- `CreateToken`, `CheckToken`, `EnsureToken`
- `CurrentToken`, `AccessToken`, `SetToken`
- `SetTokenPollInterval`, `SetTokenPollTimeout`, `EnableTokenInjection`

Correlation helpers are `WithCorrelationID` and `CorrelationIDFromContext`.
`ErrAccessTokenRequired` and `ErrCircuitOpen` remain public sentinels. The full
typed error model lives in [`pkg/errors`](errors.md).

## `data`

Construct with `data.New(*client.Client)`. The package covers the generated
reconciliation's Market Data, fundamentals, fund, event-contract, and Display
routes. Query and response types are exported alongside endpoint methods.

Financial DTOs use:

- `money.Money` for required/response values
- `*money.Money` for optional/request values

The wire JSON remains a decimal string. Application code should use
`money.MustNew` for trusted constants or `money.NewFromString` for untrusted
input; public DTOs do not expose raw `decimal.Decimal` migration semantics.

## `stream`

Construct with `stream.New(*client.Client, ...Option)`.

Lifecycle and state:

- `Connect`, `Close`, `IsConnected`, `Reconnecting`, `SessionID`
- `State`, `OnStateChange`, `OnReconnecting`
- States: `StateDisconnected`, `StateConnecting`, `StateConnected`,
  `StateReconnecting`, `StateDegraded`, `StateClosed`

`StateClosed` is terminal. Handlers and channel sends that race with `Close`
are ignored or unblocked.

Subscriptions and handlers:

- `Subscribe`, `Unsubscribe`
- `OnQuote`, `OnSnapshot`, `OnTick`, `OnNotice`, `OnError`, `OnConnect`,
  `OnDisconnect`
- `SubscribeRequest`, `UnsubscribeRequest`, `Category`, `SubType`
- Topics: `TopicQuote`, `TopicSnapshot`, `TopicTick`, `TopicNotice`, `TopicEcho`

Channel subscriptions:

- `SubscribeQuoteChan`, `SubscribeSnapshotChan`, `SubscribeTickChan`
- `ChannelConfig` with `DropBlock`, `DropOldest`, or `DropSample`
- Default buffer: 100 messages; returned cancel functions are idempotent and
  close their channel

`WithHealthWatchdog(interval)` enables message-age health monitoring. Only
quote, snapshot, and tick messages count as fresh data; notice and echo
messages do not prevent `StateDegraded`. A fresh data message restores
`StateConnected`. Webull MQTT payloads do not carry sequence numbers, so this is
a data-age signal, not sequence-gap detection.

Reconnect/resubscribe is serialized with subscription mutations. After a
successful reconnect, active HTTP subscriptions are replayed before
`OnConnect` handlers run. See [Streaming](streaming.md) for lifecycle and
channel behavior.

## `trade`

Construct with `trade.New(*client.Client, ...Option)`. Requests under
`/trading/` default to API version v3.

Main methods:

- Accounts/assets: `ListAccounts`, `GetBalance`, `GetPositions`,
  `GetCashActivities`, `GetCashActivitiesPage`, `GetAllCashActivities`
- Orders: `PreviewOrder`, `PlaceOrder`, `ReplaceOrder`, `CancelOrder`,
  `BatchPlaceOrder`
- Queries: `GetOpenOrders`, `GetOpenOrdersPage`, `GetAllOpenOrders`,
  `GetOrderHistory`, `GetOrderHistoryPage`, `GetAllOrderHistory`,
  `GetOrderDetail`

`PlaceOrder` returns `*pkg/domain/order.Order`, not the raw place response. It
embeds `PlaceOrderResult`, includes the account ID, and starts a local state
machine. `BatchPlaceOrder` registers each returned client order ID locally.

OMS methods:

- `GetTrackedOrder` / `TrackedOrder`
- `TrackedOrderState`, `TrackedOrders`
- `ApplyOrderEvent`
- `ReconcileOrderStatus`, `ReconcileOrderState`

The registry key is `(account ID, client order ID)`, so identical client IDs in
different accounts do not collide. Reconciliation accepts a Webull status
string, `trade.OrderStatus`, or `order.State`; stale non-terminal snapshots do
not regress state, and terminal states never become non-terminal.

`ReplaceOrder` and `CancelOrder` skip the network for a locally tracked terminal
order. A successful API action advances local state; a failed action leaves it
unchanged. The registry is in-memory and is not a durable source of truth.

Guardrails and builders:

- `WithMaxOrderNotional`, `WithMaxOrderQuantity`
- `WithAutoClientOrderID`
- `NewPlaceOrderRequest`, `NewEquityOrder`, `EquityOrderBuilder`
- `NewCancelOrderRequest`, `NewModifyOrderRequest`
- `NewClientOrderID`, `ClientOrderIDFrom`, `ValidClientOrderID`

Guardrail failures retain the historical `errs.CodeInvalidConfig` code and also
wrap `errs.ErrOrderGuardrail`.

## `events`

Construct with `events.New(*client.Client, ...Option)`.

- Lifecycle: `Run`, `Close`
- Handlers: `OnConnect`, `OnPing`, `OnEvent`, `OnOrder`, `OnPosition`,
  `OnOption`, `OnError`
- Payloads: `OrderEvent`, `PositionEvent`, `OptionEvent`
- Subscribe flags: `SubscribeOrder`, `SubscribePosition`, `SubscribeOption`,
  `SubscribeAll`
- Connection/reconnect options: `WithGRPCEndpoint`, `WithGRPCPort`, `WithTLS`,
  `WithDialTimeout`, `WithGRPCDialOption`, `WithSubscribeTypes`,
  `WithAccounts`, `WithAutoReconnect`, `WithReconnectBaseDelay`,
  `WithReconnectMaxDelay`, `WithMaxReconnectAttempts`

Every stream attempt emits an optional client-kind span named
`/grpc.trade.event.EventService/Subscribe`, records optional attempt and
duration metrics, and writes optional structured start/done logs. Correlation
IDs and configured trace propagators are sent as gRPC metadata. Credentials
are not telemetry attributes or log fields.

## `brokerfd` and `brokerfd/events`

`brokerfd.New` exposes the Broker FD US HTTP surface. The root module owns this
package; it is not a separate Go module.

`brokerfd/events.New` exposes the Broker FD gRPC stream and uses the
`grpc.event.EventService` RPC. It follows the same per-attempt telemetry model
as Trading Events, with instrumentation scope
`webullapi4go/brokerfd/events`, span name `/grpc.event.EventService/Subscribe`,
and shared metric names scoped separately by the OTel meter.

## Shared foundations

### `pkg/domain/money`

`Money` wraps an exact decimal and marshals as a JSON string. Key APIs are
`New`, `NewFromString`, `MustNew`, `NewFromInt64`, `Zero`, arithmetic,
comparison, `String`, `Rat`, and JSON methods. `Money` is the public DTO type;
raw `decimal.Decimal` is not a DTO migration target.

### `pkg/domain/order`

The package exports `State`, `Event`, `Machine`, `Order`, terminal-state
helpers, Webull status mapping, and scene-type mapping. `Machine` is safe for
concurrent use. `ReconcileStatus` advances by the smallest valid event path,
ignores stale non-terminal snapshots, and never regresses a terminal state.
`OnStateChanged` runs after the state is committed and without holding the
machine lock.

### `pkg/errors`

`Error` carries `Code`, `Message`, optional HTTP `Status`, and an optional
wrapped cause. Use `errors.Is`, `errors.As`, `errs.Is`, and exported sentinels
instead of parsing strings. See [Errors](errors.md).

### `pkg/observability`

The package provides OTel aliases, lazy instruments, span helpers, and trace
propagation helpers. SDK service clients use the core configuration inherited
from `client.Client`. See [Observability](observability.md).

## Internal packages

`internal/auth` contains the protocol-specific REST and gRPC signing logic.
Other internal packages are implementation details or narrow deprecated
compatibility aliases for older internal imports. They are not public service
packages, are not referenced by public signatures, and should not be imported
by applications.
