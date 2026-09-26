# Observability

The SDK supports structured logging, OpenTelemetry tracing and metrics, and
correlation IDs across REST, MQTT, and gRPC. All integrations are opt-in and
default to no-op providers or no logging.

!!! note "Status"
    REST and core stream telemetry has a historical baseline. Trading/Broker FD
    event attempt telemetry and the current request/OMS/stream hardening are
    tagged in repository `v2.1.2`, implemented, and offline-tested; they were
    not newly live-verified. `v2.1.2` is a repository patch release, not a
    published Go-semver v2 module; the module stays on the v1 import path by
    decision, so `v2.x` tags are repository-only.

## Minimal configuration

The SDK depends only on the OpenTelemetry API. The application chooses the SDK,
exporter, sampler, and backend. For a console example, add the SDK and stdout
exporters to the application module:

```sh
go get go.opentelemetry.io/otel/sdk go.opentelemetry.io/otel/exporters/stdout/stdouttrace go.opentelemetry.io/otel/exporters/stdout/stdoutmetric
```

```go
package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/shing1211/webullapi4go/client"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/propagation"
)

func main() {
	traceExporter, err := stdouttrace.New(stdouttrace.WithWriter(os.Stdout))
	if err != nil {
		panic(err)
	}
	metricExporter, err := stdoutmetric.New(stdoutmetric.WithWriter(os.Stdout))
	if err != nil {
		panic(err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
	)
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
	)
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)
	otel.SetTextMapPropagator(propagator)

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tracerProvider.Shutdown(ctx)
	}()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = meterProvider.Shutdown(ctx)
	}()

	cl, err := client.New(
		client.WithEnv(),
		client.WithTracerProvider(tracerProvider),
		client.WithMeterProvider(meterProvider),
		client.WithPropagator(propagator),
		client.WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil))),
		client.WithResiliencePreset(client.ProductionPreset),
	)
	if err != nil {
		panic(err)
	}
	defer func() { _ = cl.Close() }()

	_ = cl
}
```

The stdout exporters are illustrative. Use an OTLP or other production exporter
with the sampling, batching, retry, and resource policy appropriate for the
application.

`WithMeterProvider` and `WithResiliencePreset` are order-independent. The
production preset creates its circuit breaker after all client options have been
applied, using the final meter provider. A breaker supplied explicitly with
`WithBreaker` remains caller-owned in either order; `WithBreaker(nil)` explicitly
disables the preset breaker. Construct service clients only after the core
client is configured so they inherit the final providers.

## Core options

| Option | Default | Behavior |
|---|---|---|
| `WithLogger(*slog.Logger)` | No logging | Structured request and stream lifecycle logs |
| `WithTracerProvider` | Global no-op tracer | REST, MQTT, and gRPC event spans |
| `WithMeterProvider` | Global no-op meter | SDK metrics |
| `WithPropagator` | No injection | Injects configured trace context into HTTP headers and gRPC metadata |
| `WithHooks` | No callbacks | Per-attempt request lifecycle integration |
| `WithInterceptor` | None | Custom request-pipeline behavior |

`Client.ObservabilityConfig()` exposes the shared observability configuration
used by service clients. Treat the returned pointer as read-only after
construction.

## REST spans, logs, and metrics

`Client.Do`, `Client.DoBroker`, and `Client.DoStream` share one per-attempt
pipeline. Each attempt starts a client-kind span named `METHOD /path` in scope
`webullapi4go/client` and ends it when that attempt finishes.

The span carries:

- `http.method` and `http.route` (the path without its query string)
- `webull.attempt` (one-based and increasing across retries)
- `webull.correlation_id`
- `http.duration_ms`
- `http.status_code` when a response was received
- `error` and a recorded error event containing only sanitized error text when
  the attempt failed

Successful attempts have OpenTelemetry status `OK`; failed attempts have status
`ERROR`. The configured propagator injects the active span into the outgoing
HTTP headers, and the correlation ID is sent as `x-correlation-id`.

With `WithLogger`, every attempt emits `webull request start` and
`webull request done` records. The done record contains method, path, attempt,
status (zero before a response), correlation ID, and latency. Failure records
also contain the sanitized `error` value. The SDK does not log the query string
or response body.

The same attempt records one `request_latency` Float64 histogram sample in
milliseconds from scope `webullapi4go/client`, with `http.route` as its only
metric attribute. A request with two attempts therefore records two samples.

The request pipeline is:

```text
rate limit → circuit breaker → hooks → interceptors → sign/send → record
```

`Do` and `DoBroker` may retry eligible idempotent requests. The correlation ID
and logical request remain stable across those retries. `DoStream` performs one
attempt because replaying an open response body is not safe.

## Streaming telemetry

The stream client inherits the core tracer, meter, and logger unless
`stream.WithMeter` overrides the meter. With the inherited core meter, the
instrumentation scope is `webullapi4go/stream`; an explicitly supplied meter
uses whatever scope the application used to create that meter.

- Each MQTT message dispatch starts a consumer span named `mqtt.dispatch` in
  scope `webullapi4go/stream`.
- Span attributes are `messaging.system=mqtt`, normalized
  `messaging.destination` (`quote`, `snapshot`, `tick`, `notice`, `echo`, or
  `unknown`), and `messaging.message.body.size`.
- A successful dispatch has status `OK`. A decode/topic failure records only
  [`SafeErrorText`](#safeerrortext-redaction) and has status `ERROR`.
- `reconnects` is an Int64 counter with no attributes and counts reconnect
  starts.
- `channel_drops` is an Int64 counter with `topic=quote|snapshot|tick` and
  counts values discarded by `DropOldest` or `DropSample`.
- Structured logs record `webull stream connected`, `webull stream
  disconnected`, and `webull stream reconnecting` with `transport=mqtt` and the
  current `state`.

MQTT callback messages do not carry an application context, so dispatch spans
start from a background context. They are not automatically children of the
HTTP call that created the subscription. Event gRPC spans do inherit the
context passed to `Run`.

Channel drop metrics are recorded when `DropOldest` or `DropSample` discards a
value. A `DropBlock` subscription that is cancelled or closed while waiting
does not count as a dropped market-data message.

## Event telemetry

`events` and `brokerfd/events` create one telemetry record set for every gRPC
stream attempt, including reconnects.

| Client | Span name | Instrumentation scope |
|---|---|---|
| Trading Events | `/grpc.trade.event.EventService/Subscribe` | `webullapi4go/events` |
| Broker FD Events | `/grpc.event.EventService/Subscribe` | `webullapi4go/brokerfd/events` |

Each client-kind span includes:

- `rpc.system=grpc`, `rpc.service`, and `rpc.method=Subscribe`
- one-based `webull.attempt` and stable `webull.correlation_id`
- terminal `rpc.grpc.status_code`
- `webull.stream.outcome` (`ok` or `error`)
- `webull.stream.duration_ms`

A clean stream end records gRPC `OK` and outcome `ok`. A transport/server
failure records its gRPC status and outcome `error`, records only
`SafeErrorText`, and sets span status to `ERROR`. Context cancellation records
gRPC `Canceled`, outcome `error`, and duration even though cancellation is
normal shutdown and is not sent to `OnError`.

The duration histogram covers the lifetime of one stream attempt, not an
individual event message. It can therefore be large for a healthy long-lived
subscription.

Each scope exports:

| Instrument | Type | Attributes |
|---|---|---|
| `event_stream_attempts` | Int64 counter (`{attempt}`) | `rpc.system`, `rpc.service`, `rpc.method` |
| `event_stream_attempt_duration` | Float64 histogram in milliseconds | RPC fields, `rpc.grpc.status_code`, `webull.stream.outcome` |

With `WithLogger`, each attempt emits package-specific `... stream start` and
`... stream done` records with transport, service, method, attempt, status,
correlation ID, and latency. Failed done records use the fixed text `gRPC event
stream attempt failed`, not the server error message.

The event client injects configured trace headers and `x-correlation-id` into
outgoing gRPC metadata alongside the required Webull signing metadata. The
correlation ID and application trace parent survive reconnects; the attempt
number changes. The same integration exists for Broker FD Events.

## Resilience metrics


`WithResiliencePreset(client.ProductionPreset)` wires its circuit breaker to
the final core meter. The Int64 counter `breaker.state_transitions` in scope
`webullapi4go/resilience` records every actual transition with `from` and `to`
attributes (`closed`, `open`, and `half-open`). A no-op transition is not
counted.

A custom breaker can be constructed with a meter from
`pkg/resilience/breaker` and supplied through `client.WithBreaker`; the
production preset never replaces or instruments an explicitly supplied
breaker. A breaker created with `client.NewBreaker` has no meter option.

## Instrument reference

The inherited SDK meter creates the following instruments:

| Scope | Instrument | Type | Attributes |
|---|---|---|---|
| `webullapi4go/client` | `request_latency` | Float64 histogram (`ms`) | `http.route` |
| `webullapi4go/resilience` | `breaker.state_transitions` | Int64 counter | `from`, `to` |
| `webullapi4go/stream` | `reconnects` | Int64 counter | none |
| `webullapi4go/stream` | `channel_drops` | Int64 counter | `topic` |
| `webullapi4go/events` | `event_stream_attempts` | Int64 counter (`{attempt}`) | RPC system/service/method |
| `webullapi4go/events` | `event_stream_attempt_duration` | Float64 histogram (`ms`) | RPC fields, gRPC status, outcome |
| `webullapi4go/brokerfd/events` | `event_stream_attempts` | Int64 counter (`{attempt}`) | RPC system/service/method |
| `webullapi4go/brokerfd/events` | `event_stream_attempt_duration` | Float64 histogram (`ms`) | RPC fields, gRPC status, outcome |

Metric names are SDK contracts; exporters may apply their own unit or resource
normalization. `stream.WithMeter` preserves the stream instrument names but
uses the application-defined scope of the supplied meter.

## Correlation IDs

Every logical REST request gets a correlation ID if the context does not already
contain one. Retries reuse it. Callers can supply their own value:

```go
ctx = client.WithCorrelationID(ctx, "checkout-7f9c")
```

The SDK sends the value in the HTTP `X-Correlation-ID` header and event gRPC
metadata, and includes it in relevant logs and spans. An empty stored value is
treated as absent.

A correlation ID is not an authentication token and does not replace trace
sampling. Do not place secrets in it.

## Logs and sensitive data

The SDK does not put App Keys, App Secrets, access tokens, or signing values in
its structured log fields or telemetry attributes. It also does not copy an
HTTP response body, gRPC status message, or wrapped error cause into telemetry.

### `SafeErrorText` redaction

`observability.SafeErrorText(err)` is the public sanitizer used by REST, MQTT,
and both gRPC event clients:

| Input | Returned text |
|---|---|
| `nil` | `""` |
| `context.Canceled` | `"context canceled"` |
| `context.DeadlineExceeded` | `"context deadline exceeded"` |
| Wrapped `*errs.Error` with a code | `"webull: <CODE>"` |
| Any other error | `"operation failed"` |

It deliberately omits `Error.Message`, HTTP bodies, gRPC status text, and
wrapped causes. This prevents a server response from echoing credentials into
spans or logs. It is a stable error classifier, not a general-purpose secret
scanner and not a replacement for application redaction.

Application handlers can still introduce sensitive data. Configure `slog`
redaction, avoid logging raw order/account payloads at sensitive verbosity, and
use your backend's access controls. Core request logs contain the path without
its query string, but paths, correlation IDs, trace baggage, metric labels, and
application callback data can still be operationally sensitive. Do not copy
secrets into context-derived attributes.

## No-op and cost behavior

With no logger, propagator, or real providers, the SDK uses OTel no-op paths and
does not create application-visible spans or metrics. Passing a nil provider is
safe and falls back to the global no-op behavior. Lazy shared instruments are
created once and are safe under concurrent first use.

The SDK does not use `otelhttp` and does not automatically export data. Trace
volume for high-rate MQTT dispatch spans can be significant; sample or filter at
the application-owned `TracerProvider` according to the workload.
