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

package events

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
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	grpcstatus "google.golang.org/grpc/status"

	"github.com/shing1211/webullapi4go/client"
	eventsevents "github.com/shing1211/webullapi4go/gen/webull/brokerfd/events/v1"
)

type brokerV1SlogRecord struct {
	level   slog.Level
	message string
	attrs   map[string]slog.Value
}

type brokerV1SlogHandler struct {
	mu      sync.Mutex
	records []brokerV1SlogRecord
}

func (h *brokerV1SlogHandler) Enabled(context.Context, slog.Level) bool { return true }

func (h *brokerV1SlogHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := make(map[string]slog.Value, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value
		return true
	})
	h.mu.Lock()
	h.records = append(h.records, brokerV1SlogRecord{level: record.Level, message: record.Message, attrs: attrs})
	h.mu.Unlock()
	return nil
}

func (h *brokerV1SlogHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *brokerV1SlogHandler) WithGroup(string) slog.Handler      { return h }

func (h *brokerV1SlogHandler) snapshot() []brokerV1SlogRecord {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]brokerV1SlogRecord, len(h.records))
	copy(out, h.records)
	return out
}

func brokerV1NewTracer(t *testing.T) (*sdktrace.TracerProvider, *tracetest.SpanRecorder) {
	t.Helper()
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })
	return provider, recorder
}

func brokerV1RemoteParent(t *testing.T) (context.Context, trace.SpanContext) {
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

type brokerV1Server struct {
	eventsevents.UnimplementedEventServiceServer
	attempts atomic.Int32
	mu       sync.Mutex
	metadata []metadata.MD
	behavior func(int32, grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error
}

func (s *brokerV1Server) Subscribe(_ *eventsevents.SubscribeRequest, stream grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error {
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

func (s *brokerV1Server) captured() []metadata.MD {
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

func brokerV1SpanAttrs(span sdktrace.ReadOnlySpan) map[string]attribute.Value {
	attrs := make(map[string]attribute.Value, len(span.Attributes()))
	for _, attr := range span.Attributes() {
		attrs[string(attr.Key)] = attr.Value
	}
	return attrs
}

func brokerV1AssertStringAttr(t *testing.T, span sdktrace.ReadOnlySpan, key, want string) {
	t.Helper()
	got, ok := brokerV1SpanAttrs(span)[key]
	if !ok || got.AsString() != want {
		t.Fatalf("span %q attribute %s = %v/%t, want %q/true", span.Name(), key, got, ok, want)
	}
}

func brokerV1AssertIntAttr(t *testing.T, span sdktrace.ReadOnlySpan, key string, want int64) {
	t.Helper()
	got, ok := brokerV1SpanAttrs(span)[key]
	if !ok || got.AsInt64() != want {
		t.Fatalf("span %q attribute %s = %v/%t, want %d/true", span.Name(), key, got, ok, want)
	}
}

func brokerV1SpanText(span sdktrace.ReadOnlySpan) string {
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

func brokerV1SlogText(records []brokerV1SlogRecord) string {
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

func brokerV1AssertLogValue(t *testing.T, record brokerV1SlogRecord, key string, want any) {
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

func brokerV1AssertLogAttrs(t *testing.T, record brokerV1SlogRecord, want map[string]any) {
	t.Helper()
	if len(record.attrs) != len(want) {
		t.Fatalf("log %q attributes = %v, want exactly %v", record.message, record.attrs, want)
	}
	for key, value := range want {
		brokerV1AssertLogValue(t, record, key, value)
	}
}

func brokerV1AssertStartLog(t *testing.T, record brokerV1SlogRecord, attempt int) {
	t.Helper()
	if record.level != slog.LevelInfo || record.message != "webull brokerfd events stream start" {
		t.Fatalf("start log = %+v, want INFO start", record)
	}
	brokerV1AssertLogAttrs(t, record, map[string]any{
		"transport":      "grpc",
		"service":        "grpc.event.EventService",
		"method":         "Subscribe",
		"attempt":        attempt,
		"correlation_id": "v1-broker-correlation",
	})
}

func brokerV1AssertDoneLog(t *testing.T, record brokerV1SlogRecord, attempt int, statusCode codes.Code, failed bool) {
	t.Helper()
	want := map[string]any{
		"transport":      "grpc",
		"service":        "grpc.event.EventService",
		"method":         "Subscribe",
		"attempt":        attempt,
		"status_code":    int(statusCode),
		"status":         statusCode.String(),
		"correlation_id": "v1-broker-correlation",
		"latency":        time.Duration(0),
	}
	if failed {
		if record.level != slog.LevelError || record.message != "webull brokerfd events stream done" {
			t.Fatalf("failed log = %+v, want ERROR done", record)
		}
		want["error"] = "gRPC event stream attempt failed"
	} else {
		if record.level != slog.LevelInfo || record.message != "webull brokerfd events stream done" {
			t.Fatalf("successful log = %+v, want INFO done", record)
		}
	}
	brokerV1AssertLogAttrs(t, record, want)
}

func brokerV1AssertMetricAttempt(t *testing.T, metric metricdata.Metrics, want int64) {
	t.Helper()
	data, ok := metric.Data.(metricdata.Sum[int64])
	if !ok {
		t.Fatalf("broker event attempt metric type = %T, want metricdata.Sum[int64]", metric.Data)
	}
	var total int64
	for _, point := range data.DataPoints {
		total += point.Value
		for key, value := range map[string]string{
			"rpc.system":  "grpc",
			"rpc.service": "grpc.event.EventService",
			"rpc.method":  "Subscribe",
		} {
			got, ok := point.Attributes.Value(attribute.Key(key))
			if !ok || got.AsString() != value {
				t.Fatalf("broker event attempt attribute %s = %v/%t, want %q/true", key, got, ok, value)
			}
		}
	}
	if total != want {
		t.Fatalf("broker event attempt total = %d, want %d", total, want)
	}
}

func brokerV1AssertDuration(t *testing.T, metric metricdata.Metrics, wantStatus int64, wantOutcome string) {
	t.Helper()
	data, ok := metric.Data.(metricdata.Histogram[float64])
	if !ok {
		t.Fatalf("broker event duration metric type = %T, want metricdata.Histogram[float64]", metric.Data)
	}
	if len(data.DataPoints) != 1 || data.DataPoints[0].Count != 1 {
		t.Fatalf("broker event duration points = %d, want one", len(data.DataPoints))
	}
	point := data.DataPoints[0]
	if got, ok := point.Attributes.Value(attribute.Key("rpc.grpc.status_code")); !ok || got.AsInt64() != wantStatus {
		t.Fatalf("broker duration status = %v/%t, want %d/true", got, ok, wantStatus)
	}
	if got, ok := point.Attributes.Value(attribute.Key("webull.stream.outcome")); !ok || got.AsString() != wantOutcome {
		t.Fatalf("broker duration outcome = %v/%t, want %s/true", got, ok, wantOutcome)
	}
	for key, value := range map[string]string{
		"rpc.system":  "grpc",
		"rpc.service": "grpc.event.EventService",
		"rpc.method":  "Subscribe",
	} {
		got, ok := point.Attributes.Value(attribute.Key(key))
		if !ok || got.AsString() != value {
			t.Fatalf("broker duration attribute %s = %v/%t, want %q/true", key, got, ok, value)
		}
	}
}

func brokerV1AssertTraceparent(t *testing.T, md metadata.MD, parent trace.SpanContext, child trace.SpanContext) {
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

func TestV1BrokerEventAttemptTelemetryForOutcomesAndMetrics(t *testing.T) {
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
				return grpcstatus.Error(codes.Unavailable, "app-key=broker-v1-app-key app-secret=broker-v1-app-secret access-token=broker-v1-token")
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
			tracerProvider, recorder := brokerV1NewTracer(t)
			reader, meterProvider := newBrokerEventMetricReader(t)
			logger := &brokerV1SlogHandler{}
			server := &brokerV1Server{behavior: func(_ int32, stream grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error {
				return tc.behavior(stream)
			}}
			core := newCoreWithOptions(t,
				client.WithAppKey("broker-v1-app-key"),
				client.WithAppSecret("broker-v1-app-secret"),
				client.WithTracerProvider(tracerProvider),
				client.WithMeterProvider(meterProvider),
				client.WithPropagator(propagation.TraceContext{}),
				client.WithLogger(slog.New(logger)),
			)
			cl := newClientWithCore(t, core, server, WithAutoReconnect(false))
			parentCtx, parent := brokerV1RemoteParent(t)
			ctx := client.WithCorrelationID(parentCtx, "v1-broker-correlation")

			if tc.name == "cancellation" {
				connected := make(chan struct{}, 1)
				cl.OnConnect(func() { trySend(connected, struct{}{}) })
				run := startTestRun(t, cl, ctx)
				waitSignal(t, connected, "broker event subscription acknowledgement")
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
				t.Fatalf("ended broker event spans = %d, want 1", len(spans))
			}
			span := spans[0]
			if span.Name() != "/grpc.event.EventService/Subscribe" || span.SpanKind() != trace.SpanKindClient {
				t.Fatalf("broker event span = %q/%v, want Subscribe/client", span.Name(), span.SpanKind())
			}
			if span.InstrumentationScope().Name != "webullapi4go/brokerfd/events" {
				t.Fatalf("broker event span scope = %q, want webullapi4go/brokerfd/events", span.InstrumentationScope().Name)
			}
			if span.Parent().TraceID() != parent.TraceID() || span.Parent().SpanID() != parent.SpanID() {
				t.Fatalf("broker event span parent = %s/%s, want application parent", span.Parent().TraceID(), span.Parent().SpanID())
			}
			if span.EndTime().IsZero() || !span.SpanContext().IsValid() {
				t.Fatal("broker event attempt span did not complete")
			}
			brokerV1AssertStringAttr(t, span, "rpc.system", "grpc")
			brokerV1AssertStringAttr(t, span, "rpc.service", "grpc.event.EventService")
			brokerV1AssertStringAttr(t, span, "rpc.method", "Subscribe")
			brokerV1AssertIntAttr(t, span, "webull.attempt", 1)
			brokerV1AssertStringAttr(t, span, "webull.correlation_id", "v1-broker-correlation")
			brokerV1AssertIntAttr(t, span, "rpc.grpc.status_code", int64(tc.wantStatus))
			brokerV1AssertStringAttr(t, span, "webull.stream.outcome", tc.wantOutcome)
			if got := brokerV1SpanAttrs(span)["webull.stream.duration_ms"]; got.Type() != attribute.INT64 || got.AsInt64() < 0 {
				t.Fatalf("webull.stream.duration_ms = %v, want non-negative duration", got)
			}
			if tc.wantError {
				if span.Status().Code != otelcodes.Error || len(span.Events()) == 0 {
					t.Fatalf("failed broker event span status/events = %v/%d, want ERROR/non-empty", span.Status(), len(span.Events()))
				}
			} else if span.Status().Code != otelcodes.Ok || len(span.Events()) != 0 {
				t.Fatalf("successful broker event span status/events = %v/%d, want OK/0", span.Status(), len(span.Events()))
			}
			md := server.captured()
			if len(md) != 1 || len(md[0].Get("x-correlation-id")) != 1 || md[0].Get("x-correlation-id")[0] != "v1-broker-correlation" {
				t.Fatalf("broker event correlation metadata = %v, want v1-broker-correlation", md)
			}
			signature := md[0].Get("x-signature")
			if len(signature) != 1 || signature[0] == "" {
				t.Fatalf("broker event signature metadata = %v, want one non-empty value", signature)
			}
			brokerV1AssertTraceparent(t, md[0], parent, span.SpanContext())

			rm := collectBrokerEventMetrics(t, reader)
			brokerV1AssertMetricAttempt(t, findBrokerEventMetric(t, rm, "webullapi4go/brokerfd/events", "event_stream_attempts"), 1)
			brokerV1AssertDuration(t, findBrokerEventMetric(t, rm, "webullapi4go/brokerfd/events", "event_stream_attempt_duration"), int64(tc.wantStatus), tc.wantOutcome)

			logs := logger.snapshot()
			if len(logs) != 2 {
				t.Fatalf("broker event log count = %d, want 2", len(logs))
			}
			brokerV1AssertStartLog(t, logs[0], 1)
			brokerV1AssertDoneLog(t, logs[1], 1, tc.wantStatus, tc.wantError)
			for _, secret := range []string{"broker-v1-app-key", "broker-v1-app-secret", "broker-v1-token", signature[0]} {
				if strings.Contains(brokerV1SpanText(span), secret) || strings.Contains(brokerV1SlogText(logs), secret) {
					t.Fatalf("broker event telemetry contains credential %q", secret)
				}
			}
		})
	}
}

func TestV1BrokerEventCorrelationAndTraceContextSurviveReconnect(t *testing.T) {
	tracerProvider, recorder := brokerV1NewTracer(t)
	reader, meterProvider := newBrokerEventMetricReader(t)
	logger := &brokerV1SlogHandler{}
	secondSubscribe := make(chan struct{})
	var secondOnce sync.Once
	server := &brokerV1Server{behavior: func(attempt int32, stream grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error {
		if attempt == 1 {
			return nil
		}
		secondOnce.Do(func() { close(secondSubscribe) })
		<-stream.Context().Done()
		return stream.Context().Err()
	}}
	core := newCoreWithOptions(t,
		client.WithAppKey("broker-reconnect-v1-app-key"),
		client.WithAppSecret("broker-reconnect-v1-app-secret"),
		client.WithTracerProvider(tracerProvider),
		client.WithMeterProvider(meterProvider),
		client.WithPropagator(propagation.TraceContext{}),
		client.WithLogger(slog.New(logger)),
	)
	cl := newClientWithCore(t, core, server,
		WithAutoReconnect(true),
		WithReconnectBaseDelay(time.Millisecond),
		WithReconnectMaxDelay(2*time.Millisecond),
		WithMaxReconnectAttempts(2),
	)
	parentCtx, parent := brokerV1RemoteParent(t)
	ctx := client.WithCorrelationID(parentCtx, "v1-broker-correlation")
	run := startTestRun(t, cl, ctx)
	waitSignal(t, secondSubscribe, "reconnected broker Subscribe")
	run.cancel()
	if err := run.wait(t); !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want context.Canceled", err)
	}

	spans := recorder.Ended()
	if len(spans) != 2 {
		t.Fatalf("broker reconnect spans = %d, want 2", len(spans))
	}
	for i, span := range spans {
		if span.Parent().TraceID() != parent.TraceID() || span.Parent().SpanID() != parent.SpanID() {
			t.Fatalf("broker span %d parent = %s/%s, want application parent", i, span.Parent().TraceID(), span.Parent().SpanID())
		}
		brokerV1AssertStringAttr(t, span, "webull.correlation_id", "v1-broker-correlation")
		brokerV1AssertIntAttr(t, span, "webull.attempt", int64(i+1))
		if i == 0 {
			brokerV1AssertIntAttr(t, span, "rpc.grpc.status_code", int64(codes.OK))
			brokerV1AssertStringAttr(t, span, "webull.stream.outcome", "ok")
		} else {
			brokerV1AssertIntAttr(t, span, "rpc.grpc.status_code", int64(codes.Canceled))
			brokerV1AssertStringAttr(t, span, "webull.stream.outcome", "error")
		}
	}
	md := server.captured()
	if len(md) != 2 {
		t.Fatalf("broker metadata attempts = %d, want 2", len(md))
	}
	for i := range md {
		if got := md[i].Get("x-correlation-id"); len(got) != 1 || got[0] != "v1-broker-correlation" {
			t.Fatalf("broker metadata correlation %d = %v", i, got)
		}
		brokerV1AssertTraceparent(t, md[i], parent, spans[i].SpanContext())
	}

	rm := collectBrokerEventMetrics(t, reader)
	brokerV1AssertMetricAttempt(t, findBrokerEventMetric(t, rm, "webullapi4go/brokerfd/events", "event_stream_attempts"), 2)
	duration := findBrokerEventMetric(t, rm, "webullapi4go/brokerfd/events", "event_stream_attempt_duration")
	durationData, ok := duration.Data.(metricdata.Histogram[float64])
	if !ok || len(durationData.DataPoints) != 2 {
		t.Fatalf("broker reconnect duration data = %T with %d points, want two", duration.Data, len(durationData.DataPoints))
	}
	statuses := map[int64]string{}
	for _, point := range durationData.DataPoints {
		status, statusOK := point.Attributes.Value(attribute.Key("rpc.grpc.status_code"))
		outcome, outcomeOK := point.Attributes.Value(attribute.Key("webull.stream.outcome"))
		if !statusOK || !outcomeOK {
			t.Fatalf("broker duration point lacks status/outcome: %v", point.Attributes)
		}
		statuses[status.AsInt64()] = outcome.AsString()
	}
	if statuses[int64(codes.OK)] != "ok" || statuses[int64(codes.Canceled)] != "error" {
		t.Fatalf("broker reconnect duration status/outcomes = %v", statuses)
	}

	logs := logger.snapshot()
	if len(logs) != 4 {
		t.Fatalf("broker reconnect log count = %d, want 4", len(logs))
	}
	brokerV1AssertStartLog(t, logs[0], 1)
	brokerV1AssertDoneLog(t, logs[1], 1, codes.OK, false)
	brokerV1AssertStartLog(t, logs[2], 2)
	brokerV1AssertDoneLog(t, logs[3], 2, codes.Canceled, true)
	for _, span := range spans {
		for _, secret := range []string{"broker-reconnect-v1-app-key", "broker-reconnect-v1-app-secret"} {
			if strings.Contains(brokerV1SpanText(span), secret) {
				t.Fatalf("broker reconnect span contains credential %q", secret)
			}
		}
	}
}

func TestV1BrokerEventNilLoggerIsSafe(t *testing.T) {
	provider, _ := brokerV1NewTracer(t)
	core := newCoreWithOptions(t, client.WithTracerProvider(provider))
	cl := newClientWithCore(t, core, &brokerV1Server{}, WithAutoReconnect(false))
	if err := cl.Run(context.Background()); err != nil {
		t.Fatalf("Run() with nil logger error = %v", err)
	}
}
