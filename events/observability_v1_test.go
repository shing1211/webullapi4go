// Copyright 2026 shing1211
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package events_test

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	grpcstatus "google.golang.org/grpc/status"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/events"
	eventsevents "github.com/shing1211/webullapi4go/gen/webull/trade/events/v1"
)

type eventV1SlogRecord struct {
	level   slog.Level
	message string
	attrs   map[string]slog.Value
}

type eventV1SlogHandler struct {
	mu      sync.Mutex
	records []eventV1SlogRecord
}

func (h *eventV1SlogHandler) Enabled(context.Context, slog.Level) bool { return true }

func (h *eventV1SlogHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := make(map[string]slog.Value, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value
		return true
	})
	h.mu.Lock()
	h.records = append(h.records, eventV1SlogRecord{level: record.Level, message: record.Message, attrs: attrs})
	h.mu.Unlock()
	return nil
}

func (h *eventV1SlogHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *eventV1SlogHandler) WithGroup(string) slog.Handler      { return h }

func (h *eventV1SlogHandler) snapshot() []eventV1SlogRecord {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]eventV1SlogRecord, len(h.records))
	copy(out, h.records)
	return out
}

func eventV1NewTracer(t *testing.T) (*sdktrace.TracerProvider, *tracetest.SpanRecorder) {
	t.Helper()
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })
	return provider, recorder
}

func eventV1RemoteParent(t *testing.T) (context.Context, trace.SpanContext) {
	t.Helper()
	traceID, err := trace.TraceIDFromHex("0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatal(err)
	}
	spanID, err := trace.SpanIDFromHex("0123456789abcdef")
	if err != nil {
		t.Fatal(err)
	}
	parent := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
		Remote:     true,
	})
	return trace.ContextWithRemoteSpanContext(context.Background(), parent), parent
}

type eventV1Server struct {
	eventsevents.UnimplementedEventServiceServer
	attempts atomic.Int32
	mu       sync.Mutex
	metadata []metadata.MD
	behavior func(int32, grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error
}

func (s *eventV1Server) Subscribe(_ *eventsevents.SubscribeRequest, stream grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error {
	md, _ := metadata.FromIncomingContext(stream.Context())
	captured := metadata.MD{}
	for key, values := range md {
		captured[key] = append([]string(nil), values...)
	}
	attempt := s.attempts.Add(1)
	s.mu.Lock()
	s.metadata = append(s.metadata, captured)
	s.mu.Unlock()
	if s.behavior == nil {
		return nil
	}
	return s.behavior(attempt, stream)
}

func (s *eventV1Server) captured() []metadata.MD {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]metadata.MD, len(s.metadata))
	for i, md := range s.metadata {
		out[i] = metadata.MD{}
		for key, values := range md {
			out[i][key] = append([]string(nil), values...)
		}
	}
	return out
}

func eventV1SpanAttrs(span sdktrace.ReadOnlySpan) map[string]attribute.Value {
	attrs := make(map[string]attribute.Value, len(span.Attributes()))
	for _, attr := range span.Attributes() {
		attrs[string(attr.Key)] = attr.Value
	}
	return attrs
}

func eventV1AssertStringAttr(t *testing.T, span sdktrace.ReadOnlySpan, key, want string) {
	t.Helper()
	got, ok := eventV1SpanAttrs(span)[key]
	if !ok || got.AsString() != want {
		t.Fatalf("span %q attribute %s = %v/%t, want %q/true", span.Name(), key, got, ok, want)
	}
}

func eventV1AssertIntAttr(t *testing.T, span sdktrace.ReadOnlySpan, key string, want int64) {
	t.Helper()
	got, ok := eventV1SpanAttrs(span)[key]
	if !ok || got.AsInt64() != want {
		t.Fatalf("span %q attribute %s = %v/%t, want %d/true", span.Name(), key, got, ok, want)
	}
}

func eventV1SpanText(span sdktrace.ReadOnlySpan) string {
	var builder strings.Builder
	builder.WriteString(span.Name())
	builder.WriteString(span.Status().Description)
	for _, attr := range span.Attributes() {
		builder.WriteString(string(attr.Key))
		builder.WriteString(attr.Value.String())
	}
	for _, event := range span.Events() {
		builder.WriteString(event.Name)
		for _, attr := range event.Attributes {
			builder.WriteString(string(attr.Key))
			builder.WriteString(attr.Value.String())
		}
	}
	return builder.String()
}

func eventV1SlogText(records []eventV1SlogRecord) string {
	var builder strings.Builder
	for _, record := range records {
		builder.WriteString(record.message)
		for key, value := range record.attrs {
			builder.WriteString(key)
			builder.WriteString(value.String())
		}
	}
	return builder.String()
}

func eventV1AssertLogValue(t *testing.T, record eventV1SlogRecord, key string, want any) {
	t.Helper()
	got, ok := record.attrs[key]
	if !ok {
		t.Fatalf("log %q missing attribute %q: %v", record.message, key, record.attrs)
	}
	switch expected := want.(type) {
	case string:
		if got.Kind() != slog.KindString || got.String() != expected {
			t.Fatalf("log %q attribute %s = %v, want %q", record.message, key, got, expected)
		}
	case int:
		if got.Kind() != slog.KindInt64 || got.Int64() != int64(expected) {
			t.Fatalf("log %q attribute %s = %v, want %d", record.message, key, got, expected)
		}
	case time.Duration:
		if got.Kind() != slog.KindDuration || got.Duration() < 0 {
			t.Fatalf("log %q attribute %s = %v, want a non-negative duration", record.message, key, got)
		}
	default:
		t.Fatalf("unsupported log expectation %T", want)
	}
}

func eventV1AssertLogAttrs(t *testing.T, record eventV1SlogRecord, want map[string]any) {
	t.Helper()
	if len(record.attrs) != len(want) {
		t.Fatalf("log %q attributes = %v, want exactly %v", record.message, record.attrs, want)
	}
	for key, value := range want {
		eventV1AssertLogValue(t, record, key, value)
	}
}

func eventV1AssertAttemptLog(t *testing.T, record eventV1SlogRecord, attempt int, statusCode codes.Code, failed bool) {
	t.Helper()
	if failed {
		if record.level != slog.LevelError || record.message != "webull events stream done" {
			t.Fatalf("failed log = %+v, want ERROR done", record)
		}
		want := map[string]any{
			"transport":      "grpc",
			"service":        "grpc.trade.event.EventService",
			"method":         "Subscribe",
			"attempt":        attempt,
			"status_code":    int(statusCode),
			"status":         statusCode.String(),
			"correlation_id": "v1-event-correlation",
			"latency":        time.Duration(0),
			"error":          "gRPC event stream attempt failed",
		}
		eventV1AssertLogAttrs(t, record, want)
		return
	}
	if record.level != slog.LevelInfo || record.message != "webull events stream done" {
		t.Fatalf("successful log = %+v, want INFO done", record)
	}
	eventV1AssertLogAttrs(t, record, map[string]any{
		"transport":      "grpc",
		"service":        "grpc.trade.event.EventService",
		"method":         "Subscribe",
		"attempt":        attempt,
		"status_code":    int(statusCode),
		"status":         statusCode.String(),
		"correlation_id": "v1-event-correlation",
		"latency":        time.Duration(0),
	})
}

func eventV1AssertStartLog(t *testing.T, record eventV1SlogRecord, attempt int) {
	t.Helper()
	if record.level != slog.LevelInfo || record.message != "webull events stream start" {
		t.Fatalf("start log = %+v, want INFO start", record)
	}
	eventV1AssertLogAttrs(t, record, map[string]any{
		"transport":      "grpc",
		"service":        "grpc.trade.event.EventService",
		"method":         "Subscribe",
		"attempt":        attempt,
		"correlation_id": "v1-event-correlation",
	})
}

func eventV1AssertMetricAttempt(t *testing.T, metric metricdata.Metrics, want int64) {
	t.Helper()
	data, ok := metric.Data.(metricdata.Sum[int64])
	if !ok {
		t.Fatalf("event attempt metric type = %T, want metricdata.Sum[int64]", metric.Data)
	}
	var total int64
	for _, point := range data.DataPoints {
		total += point.Value
		for key, value := range map[string]string{
			"rpc.system":  "grpc",
			"rpc.service": "grpc.trade.event.EventService",
			"rpc.method":  "Subscribe",
		} {
			got, ok := point.Attributes.Value(attribute.Key(key))
			if !ok || got.AsString() != value {
				t.Fatalf("event attempt attribute %s = %v/%t, want %q/true", key, got, ok, value)
			}
		}
	}
	if total != want {
		t.Fatalf("event attempt total = %d, want %d", total, want)
	}
}

func eventV1AssertDuration(t *testing.T, metric metricdata.Metrics, wantStatus int64, wantOutcome string) {
	t.Helper()
	data, ok := metric.Data.(metricdata.Histogram[float64])
	if !ok {
		t.Fatalf("event duration metric type = %T, want metricdata.Histogram[float64]", metric.Data)
	}
	if len(data.DataPoints) != 1 || data.DataPoints[0].Count != 1 {
		t.Fatalf("event duration points = %d, want one", len(data.DataPoints))
	}
	point := data.DataPoints[0]
	if got, ok := point.Attributes.Value(attribute.Key("rpc.grpc.status_code")); !ok || got.AsInt64() != wantStatus {
		t.Fatalf("duration status = %v/%t, want %d/true", got, ok, wantStatus)
	}
	if got, ok := point.Attributes.Value(attribute.Key("webull.stream.outcome")); !ok || got.AsString() != wantOutcome {
		t.Fatalf("duration outcome = %v/%t, want %s/true", got, ok, wantOutcome)
	}
	for key, value := range map[string]string{
		"rpc.system":  "grpc",
		"rpc.service": "grpc.trade.event.EventService",
		"rpc.method":  "Subscribe",
	} {
		got, ok := point.Attributes.Value(attribute.Key(key))
		if !ok || got.AsString() != value {
			t.Fatalf("duration attribute %s = %v/%t, want %q/true", key, got, ok, value)
		}
	}
}

func eventV1AssertTraceparent(t *testing.T, md metadata.MD, parent trace.SpanContext, child trace.SpanContext) {
	t.Helper()
	carrier := propagation.MapCarrier{}
	for key, values := range md {
		if len(values) > 0 {
			carrier[key] = values[0]
		}
	}
	raw := carrier.Get("traceparent")
	if raw == "" {
		t.Fatal("gRPC metadata has no traceparent")
	}
	extracted := propagation.TraceContext{}.Extract(context.Background(), carrier)
	got := trace.SpanContextFromContext(extracted)
	if got.TraceID() != parent.TraceID() || got.SpanID() != child.SpanID() {
		t.Fatalf("traceparent = %q, extracted trace/span = %s/%s, want %s/%s", raw, got.TraceID(), got.SpanID(), parent.TraceID(), child.SpanID())
	}
}

func eventV1NewMetricProvider(t *testing.T) (*sdkmetric.ManualReader, *sdkmetric.MeterProvider) {
	return newEventMetricReader(t)
}

func TestV1EventAttemptTelemetryForOutcomesAndMetrics(t *testing.T) {
	cases := []struct {
		name        string
		behavior    func(grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error
		wantStatus  codes.Code
		wantOutcome string
		wantError   bool
	}{
		{
			name:        "success",
			behavior:    func(grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error { return nil },
			wantStatus:  codes.OK,
			wantOutcome: "ok",
		},
		{
			name: "failure",
			behavior: func(grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error {
				return grpcstatus.Error(codes.Unavailable, "app-key=event-v1-app-key app-secret=event-v1-app-secret access-token=event-v1-token")
			},
			wantStatus:  codes.Unavailable,
			wantOutcome: "error",
			wantError:   true,
		},
		{
			name: "cancellation",
			behavior: func(stream grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error {
				if err := stream.Send(&eventsevents.SubscribeResponse{EventType: eventsevents.EventType_SubscribeSuccess}); err != nil {
					return err
				}
				<-stream.Context().Done()
				return stream.Context().Err()
			},
			wantStatus:  codes.Canceled,
			wantOutcome: "error",
			wantError:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tracerProvider, recorder := eventV1NewTracer(t)
			reader, meterProvider := eventV1NewMetricProvider(t)
			logger := &eventV1SlogHandler{}
			server := &eventV1Server{behavior: func(_ int32, stream grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error {
				return tc.behavior(stream)
			}}
			core := newCoreWithOptions(t,
				client.WithAppKey("event-v1-app-key"),
				client.WithAppSecret("event-v1-app-secret"),
				client.WithTracerProvider(tracerProvider),
				client.WithMeterProvider(meterProvider),
				client.WithPropagator(propagation.TraceContext{}),
				client.WithLogger(slog.New(logger)),
			)
			cl := newClientWithCore(t, core, server, events.WithAutoReconnect(false))
			parentCtx, parent := eventV1RemoteParent(t)
			ctx := client.WithCorrelationID(parentCtx, "v1-event-correlation")

			if tc.name == "cancellation" {
				connected := make(chan struct{}, 1)
				cl.OnConnect(func() { trySend(connected, struct{}{}) })
				run := startTestRun(t, cl, ctx)
				waitSignal(t, connected, "event subscription acknowledgement")
				run.cancel()
				if err := run.wait(t); !errors.Is(err, context.Canceled) {
					t.Fatalf("Run() error = %v, want context.Canceled", err)
				}
			} else if err := cl.Run(ctx); err != nil {
				if tc.name == "success" {
					t.Fatalf("Run() error = %v", err)
				}
				if !strings.Contains(err.Error(), "Unavailable") {
					t.Fatalf("Run() error = %v, want unavailable failure", err)
				}
			}

			spans := recorder.Ended()
			if len(spans) != 1 {
				t.Fatalf("ended event spans = %d, want 1", len(spans))
			}
			span := spans[0]
			if span.Name() != "/grpc.trade.event.EventService/Subscribe" || span.SpanKind() != trace.SpanKindClient {
				t.Fatalf("event span = %q/%v, want Subscribe/client", span.Name(), span.SpanKind())
			}
			if span.InstrumentationScope().Name != "webullapi4go/events" {
				t.Fatalf("event span scope = %q, want webullapi4go/events", span.InstrumentationScope().Name)
			}
			if span.Parent().TraceID() != parent.TraceID() || span.Parent().SpanID() != parent.SpanID() {
				t.Fatalf("event span parent = %s/%s, want application parent", span.Parent().TraceID(), span.Parent().SpanID())
			}
			if span.EndTime().IsZero() || !span.SpanContext().IsValid() {
				t.Fatal("event attempt span did not complete")
			}
			eventV1AssertStringAttr(t, span, "rpc.system", "grpc")
			eventV1AssertStringAttr(t, span, "rpc.service", "grpc.trade.event.EventService")
			eventV1AssertStringAttr(t, span, "rpc.method", "Subscribe")
			eventV1AssertIntAttr(t, span, "webull.attempt", 1)
			eventV1AssertStringAttr(t, span, "webull.correlation_id", "v1-event-correlation")
			eventV1AssertIntAttr(t, span, "rpc.grpc.status_code", int64(tc.wantStatus))
			eventV1AssertStringAttr(t, span, "webull.stream.outcome", tc.wantOutcome)
			if got := eventV1SpanAttrs(span)["webull.stream.duration_ms"]; got.Type() != attribute.INT64 || got.AsInt64() < 0 {
				t.Fatalf("webull.stream.duration_ms = %v, want non-negative duration", got)
			}
			if tc.wantError {
				if span.Status().Code != otelcodes.Error {
					t.Fatalf("event span status = %v, want ERROR", span.Status())
				}
				if len(span.Events()) == 0 {
					t.Fatal("failed event span has no recorded error event")
				}
			} else if span.Status().Code != otelcodes.Ok || len(span.Events()) != 0 {
				t.Fatalf("successful event span status/events = %v/%d, want OK/0", span.Status(), len(span.Events()))
			}
			md := server.captured()
			if len(md) != 1 || len(md[0].Get("x-correlation-id")) != 1 || md[0].Get("x-correlation-id")[0] != "v1-event-correlation" {
				t.Fatalf("event correlation metadata = %v, want v1-event-correlation", md)
			}
			signature := md[0].Get("x-signature")
			if len(signature) != 1 || signature[0] == "" {
				t.Fatalf("event signature metadata = %v, want one non-empty value", signature)
			}
			eventV1AssertTraceparent(t, md[0], parent, span.SpanContext())

			rm := collectEventMetrics(t, reader)
			eventV1AssertMetricAttempt(t, findEventMetric(t, rm, "webullapi4go/events", "event_stream_attempts"), 1)
			eventV1AssertDuration(t, findEventMetric(t, rm, "webullapi4go/events", "event_stream_attempt_duration"), int64(tc.wantStatus), tc.wantOutcome)

			logs := logger.snapshot()
			if len(logs) != 2 {
				t.Fatalf("event log count = %d, want start and done", len(logs))
			}
			eventV1AssertStartLog(t, logs[0], 1)
			eventV1AssertAttemptLog(t, logs[1], 1, tc.wantStatus, tc.wantError)
			for _, secret := range []string{"event-v1-app-key", "event-v1-app-secret", "event-v1-token", signature[0]} {
				if strings.Contains(eventV1SpanText(span), secret) || strings.Contains(eventV1SlogText(logs), secret) {
					t.Fatalf("event telemetry contains credential %q", secret)
				}
			}
		})
	}
}

func TestV1EventCorrelationAndTraceContextSurviveReconnect(t *testing.T) {
	tracerProvider, recorder := eventV1NewTracer(t)
	reader, meterProvider := newEventMetricReader(t)
	logger := &eventV1SlogHandler{}
	secondSubscribe := make(chan struct{})
	var secondOnce sync.Once
	server := &eventV1Server{behavior: func(attempt int32, stream grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error {
		if attempt == 1 {
			return nil
		}
		secondOnce.Do(func() { close(secondSubscribe) })
		<-stream.Context().Done()
		return stream.Context().Err()
	}}
	core := newCoreWithOptions(t,
		client.WithAppKey("reconnect-v1-app-key"),
		client.WithAppSecret("reconnect-v1-app-secret"),
		client.WithTracerProvider(tracerProvider),
		client.WithMeterProvider(meterProvider),
		client.WithPropagator(propagation.TraceContext{}),
		client.WithLogger(slog.New(logger)),
	)
	cl := newClientWithCore(t, core, server,
		events.WithAutoReconnect(true),
		events.WithReconnectBaseDelay(time.Millisecond),
		events.WithReconnectMaxDelay(2*time.Millisecond),
		events.WithMaxReconnectAttempts(2),
	)
	parentCtx, parent := eventV1RemoteParent(t)
	ctx := client.WithCorrelationID(parentCtx, "v1-event-correlation")
	run := startTestRun(t, cl, ctx)
	waitSignal(t, secondSubscribe, "reconnected Subscribe")
	run.cancel()
	if err := run.wait(t); !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want context.Canceled", err)
	}

	spans := recorder.Ended()
	if len(spans) != 2 {
		t.Fatalf("reconnect spans = %d, want 2", len(spans))
	}
	for i, span := range spans {
		attempt := int64(i + 1)
		if span.Parent().TraceID() != parent.TraceID() || span.Parent().SpanID() != parent.SpanID() {
			t.Fatalf("span %d parent = %s/%s, want application parent", i, span.Parent().TraceID(), span.Parent().SpanID())
		}
		eventV1AssertStringAttr(t, span, "webull.correlation_id", "v1-event-correlation")
		eventV1AssertIntAttr(t, span, "webull.attempt", attempt)
		if i == 0 {
			eventV1AssertIntAttr(t, span, "rpc.grpc.status_code", int64(codes.OK))
			eventV1AssertStringAttr(t, span, "webull.stream.outcome", "ok")
			if span.Status().Code != otelcodes.Ok {
				t.Fatalf("first reconnect span status = %v, want OK", span.Status())
			}
		} else {
			eventV1AssertIntAttr(t, span, "rpc.grpc.status_code", int64(codes.Canceled))
			eventV1AssertStringAttr(t, span, "webull.stream.outcome", "error")
			if span.Status().Code != otelcodes.Error {
				t.Fatalf("second reconnect span status = %v, want ERROR", span.Status())
			}
		}
	}

	md := server.captured()
	if len(md) != 2 {
		t.Fatalf("metadata attempts = %d, want 2", len(md))
	}
	for i := range md {
		if got := md[i].Get("x-correlation-id"); len(got) != 1 || got[0] != "v1-event-correlation" {
			t.Fatalf("metadata correlation %d = %v, want stable ID", i, got)
		}
		eventV1AssertTraceparent(t, md[i], parent, spans[i].SpanContext())
	}

	rm := collectEventMetrics(t, reader)
	eventV1AssertMetricAttempt(t, findEventMetric(t, rm, "webullapi4go/events", "event_stream_attempts"), 2)
	duration := findEventMetric(t, rm, "webullapi4go/events", "event_stream_attempt_duration")
	durationData, ok := duration.Data.(metricdata.Histogram[float64])
	if !ok || len(durationData.DataPoints) != 2 {
		t.Fatalf("reconnect duration data = %T with %d points, want two", duration.Data, len(durationData.DataPoints))
	}
	statuses := map[int64]string{}
	for _, point := range durationData.DataPoints {
		status, ok := point.Attributes.Value(attribute.Key("rpc.grpc.status_code"))
		outcome, outcomeOK := point.Attributes.Value(attribute.Key("webull.stream.outcome"))
		if !ok || !outcomeOK {
			t.Fatalf("reconnect duration point lacks status/outcome: %v", point.Attributes)
		}
		statuses[status.AsInt64()] = outcome.AsString()
	}
	if statuses[int64(codes.OK)] != "ok" || statuses[int64(codes.Canceled)] != "error" {
		t.Fatalf("reconnect duration status/outcomes = %v", statuses)
	}

	logs := logger.snapshot()
	if len(logs) != 4 {
		t.Fatalf("reconnect log count = %d, want 4", len(logs))
	}
	eventV1AssertStartLog(t, logs[0], 1)
	eventV1AssertAttemptLog(t, logs[1], 1, codes.OK, false)
	eventV1AssertStartLog(t, logs[2], 2)
	eventV1AssertAttemptLog(t, logs[3], 2, codes.Canceled, true)
	for _, span := range spans {
		text := eventV1SpanText(span)
		for _, secret := range []string{"reconnect-v1-app-key", "reconnect-v1-app-secret"} {
			if strings.Contains(text, secret) {
				t.Fatalf("reconnect span contains credential %q", secret)
			}
		}
	}
}

func TestV1EventNilLoggerIsSafe(t *testing.T) {
	provider, _ := eventV1NewTracer(t)
	core := newCoreWithOptions(t, client.WithTracerProvider(provider))
	cl := newClientWithCore(t, core, &eventV1Server{}, events.WithAutoReconnect(false))
	if err := cl.Run(context.Background()); err != nil {
		t.Fatalf("Run() with nil logger error = %v", err)
	}
}
