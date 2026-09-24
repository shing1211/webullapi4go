# Observability

The SDK supports structured logging, OpenTelemetry tracing and metrics, and
correlation IDs across REST, MQTT, and gRPC. All integrations are opt-in and
default to no-op providers or no logging.

!!! note "Status"
    REST and core stream telemetry has a released baseline. Trading/Broker FD
    event attempt telemetry and the current request/OMS/stream hardening are
    **Unreleased**, implemented, and offline-tested; they are not newly
    live-verified.

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

Set `WithMeterProvider` **before** `WithResiliencePreset`. The preset constructs
its circuit breaker with the meter available at that point. Construct service
clients only after the core client is configured so they inherit the final
providers.

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

## REST spans and metrics

`Client.Do`, `Client.DoBroker`, and `Client.DoStream` share one per-attempt
pipeline. Each attempt can emit a client-kind span named `METHOD /path` with:

- `http.method`
- `http.route`
- `webull.attempt` (one-based and stable across retries)
- `webull.correlation_id`
- `http.duration_ms`
- `http.status_code` when a response was received
- `error` when the attempt failed

The same attempt records the `request_latency` float histogram in milliseconds
from instrumentation scope `webullapi4go/client`, with `http.route` as an
attribute.

The request pipeline is:

```text
rate limit → circuit breaker → hooks → interceptors → sign/send → record
```

`Do` and `DoBroker` may retry eligible idempotent requests. The correlation ID
and logical request remain stable across those retries. `DoStream` performs one
attempt because replaying an open response body is not safe.

## Streaming telemetry

The stream client inherits the core tracer, meter, and logger unless
`stream.WithMeter` overrides the meter.

- A trace can be emitted for each MQTT message dispatch as the consumer span
  `mqtt.dispatch`.
- Span attributes include `messaging.system=mqtt`, the normalized destination
  (`quote`, `snapshot`, `tick`, `notice`, or `echo`), and payload size.
- `reconnects` in scope `webullapi4go/stream` counts reconnect starts.
- `channel_drops` in the same scope counts channel discards with a `topic`
  attribute.
- Structured logs record connect, disconnect, and reconnect transitions.

MQTT callback messages do not carry an application context, so dispatch spans
start from a background context. They are not automatically children of the
HTTP call that created the subscription. Event gRPC spans do inherit the
context passed to `Run`.

Channel drop metrics are recorded when `DropOldest` or `DropSample` discards a
value. A `DropBlock` subscription that is cancelled while waiting does not count
as a dropped market-data message.

## Event telemetry

`events` and `brokerfd/events` create telemetry for each gRPC stream attempt,
including reconnect attempts.

| Client | Span name | Instrumentation scope |
|---|---|---|
| Trading Events | `/grpc.trade.event.EventService/Subscribe` | `webullapi4go/events` |
| Broker FD Events | `/grpc.event.EventService/Subscribe` | `webullapi4go/brokerfd/events` |

Each client-kind span includes:

- `rpc.system=grpc`
- RPC service and method
- `webull.attempt`
- `webull.correlation_id`
- gRPC status code on completion
- `webull.stream.outcome` (`ok` or `error`)
- `webull.stream.duration_ms`

The duration histogram covers the lifetime of one stream attempt, not an
individual event message. It can therefore be large for a healthy long-lived
subscription.

Each scope exports:

| Instrument | Type | Attributes |
|---|---|---|
| `event_stream_attempts` | Int64 counter | RPC system, service, method |
| `event_stream_attempt_duration` | Float64 histogram in milliseconds | RPC fields, gRPC status, outcome |

The event client injects configured trace headers and `x-correlation-id` into
outgoing gRPC metadata alongside the required Webull signing metadata. The
same integration exists for Broker FD Events.

## Resilience metrics

`WithResiliencePreset(client.ProductionPreset)` wires a circuit breaker to the
core meter. Its `breaker.state_transitions` Int64 counter uses `from` and `to`
attributes.

A custom breaker can be constructed with a meter from
`pkg/resilience/breaker` and supplied through `client.WithBreaker`; a breaker
created with `client.NewBreaker` does not take a meter option.

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
its structured log fields or telemetry attributes. Event telemetry tests also
assert that credentials are absent from logs.

Application handlers can still introduce sensitive data. Configure `slog`
redaction, avoid logging raw order/account payloads at sensitive verbosity, and
use your backend's access controls. Logs include request paths without query
strings in the core request logger; avoid copying secrets into context-derived
attributes.

## No-op and cost behavior

With no logger, propagator, or real providers, the SDK uses OTel no-op paths and
does not create application-visible spans or metrics. The event and stream
metrics are optional instruments; recording callbacks still execute when a real
meter is configured.

The SDK does not use `otelhttp` and does not automatically export data. Trace
volume for high-rate MQTT dispatch spans can be significant; sample or filter at
the application-owned `TracerProvider` according to the workload.
