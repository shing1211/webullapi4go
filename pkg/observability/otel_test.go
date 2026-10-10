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

package observability

import (
	"context"
	stderrors "errors"
	"strings"
	"sync"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/trace"

	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

func TestDefaultConfigCreatesUsableNoopHandles(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.TracerProvider == nil || cfg.MeterProvider == nil || cfg.Propagator == nil {
		t.Fatal("DefaultConfig() contains a nil provider")
	}
	if cfg.Tracer("test") == nil || cfg.Meter("test") == nil {
		t.Fatal("DefaultConfig() did not create tracer and meter handles")
	}
	if histogram := cfg.ClientLatencyHistogram(); histogram == nil {
		t.Fatal("ClientLatencyHistogram() = nil")
	}
	if counter := cfg.BreakerTransitionsCounter(); counter == nil {
		t.Fatal("BreakerTransitionsCounter() = nil")
	}
	if counter := cfg.StreamReconnectsCounter(); counter == nil {
		t.Fatal("StreamReconnectsCounter() = nil")
	}
	if counter := cfg.StreamChannelDropsCounter(); counter == nil {
		t.Fatal("StreamChannelDropsCounter() = nil")
	}
}

func TestTraceContextRoundTrip(t *testing.T) {
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
	})
	ctx := trace.ContextWithSpanContext(context.Background(), parent)
	cfg := DefaultConfig()
	cfg.Propagator = propagation.TraceContext{}
	carrier := propagation.MapCarrier{}
	cfg.InjectTraceContext(ctx, carrier)
	if !strings.Contains(carrier.Get("traceparent"), traceID.String()) {
		t.Fatalf("traceparent = %q, want trace ID %q", carrier.Get("traceparent"), traceID)
	}
	extracted := cfg.ExtractTraceContext(context.Background(), carrier)
	if got := trace.SpanContextFromContext(extracted); got.TraceID() != traceID {
		t.Fatalf("extracted trace ID = %q, want %q", got.TraceID(), traceID)
	}
}

func TestTraceContextIsNoopWithoutPropagator(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Propagator = nil
	carrier := propagation.MapCarrier{}
	cfg.InjectTraceContext(context.Background(), carrier)
	if len(carrier) != 0 {
		t.Fatalf("carrier = %v, want empty", carrier)
	}
	if cfg.ExtractTraceContext(context.Background(), carrier) == nil {
		t.Fatal("ExtractTraceContext() = nil")
	}
}

func newObservabilityMetricReader(t *testing.T) (*sdkmetric.ManualReader, *sdkmetric.MeterProvider) {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })
	return reader, provider
}

func findObservabilityMetric(t *testing.T, rm metricdata.ResourceMetrics, scope, name string) metricdata.Metrics {
	t.Helper()
	for _, scopeMetrics := range rm.ScopeMetrics {
		if scopeMetrics.Scope.Name != scope {
			continue
		}
		for _, metric := range scopeMetrics.Metrics {
			if metric.Name == name {
				return metric
			}
		}
	}
	t.Fatalf("metric %q in scope %q not found", name, scope)
	return metricdata.Metrics{}
}

func collectObservabilityMetrics(t *testing.T, reader *sdkmetric.ManualReader) metricdata.ResourceMetrics {
	t.Helper()
	var rm metricdata.ResourceMetrics
	if err := reader.Collect(context.Background(), &rm); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	return rm
}

func TestLazyInstrumentsRecordValuesAndAttributes(t *testing.T) {
	reader, provider := newObservabilityMetricReader(t)
	cfg := Config{MeterProvider: provider}

	cfg.ClientLatencyHistogram().Record(context.Background(), 12.5,
		metric.WithAttributes(attribute.String("http.route", "/quotes")))
	cfg.BreakerTransitionsCounter().Add(context.Background(), 2,
		metric.WithAttributes(
			attribute.String("from", "closed"),
			attribute.String("to", "open"),
		))
	cfg.StreamReconnectsCounter().Add(context.Background(), 3)
	cfg.StreamChannelDropsCounter().Add(context.Background(), 4,
		metric.WithAttributes(attribute.String("topic", "quote")))

	rm := collectObservabilityMetrics(t, reader)
	latency := findObservabilityMetric(t, rm, "webullapi4go/client", "request_latency")
	latencyData, ok := latency.Data.(metricdata.Histogram[float64])
	if !ok || len(latencyData.DataPoints) != 1 {
		t.Fatalf("latency data = %T with %d points, want one histogram point", latency.Data, len(latencyData.DataPoints))
	}
	point := latencyData.DataPoints[0]
	if point.Count != 1 || point.Sum != 12.5 {
		t.Fatalf("latency count/sum = %d/%v, want 1/12.5", point.Count, point.Sum)
	}
	if route, ok := point.Attributes.Value(attribute.Key("http.route")); !ok || route.AsString() != "/quotes" {
		t.Fatalf("latency route = %v/%t, want /quotes/true", route, ok)
	}

	transitions := findObservabilityMetric(t, rm, "webullapi4go/resilience", "breaker.state_transitions")
	transitionData, ok := transitions.Data.(metricdata.Sum[int64])
	if !ok || len(transitionData.DataPoints) != 1 || transitionData.DataPoints[0].Value != 2 {
		t.Fatalf("breaker transition data = %#v, want one value of 2", transitions.Data)
	}

	reconnects := findObservabilityMetric(t, rm, "webullapi4go/stream", "reconnects")
	reconnectData, ok := reconnects.Data.(metricdata.Sum[int64])
	if !ok || len(reconnectData.DataPoints) != 1 || reconnectData.DataPoints[0].Value != 3 {
		t.Fatalf("reconnect data = %#v, want one value of 3", reconnects.Data)
	}

	drops := findObservabilityMetric(t, rm, "webullapi4go/stream", "channel_drops")
	dropData, ok := drops.Data.(metricdata.Sum[int64])
	if !ok || len(dropData.DataPoints) != 1 || dropData.DataPoints[0].Value != 4 {
		t.Fatalf("drop data = %#v, want one value of 4", drops.Data)
	}
	if topic, ok := dropData.DataPoints[0].Attributes.Value(attribute.Key("topic")); !ok || topic.AsString() != "quote" {
		t.Fatalf("drop topic = %v/%t, want quote/true", topic, ok)
	}
}

func TestLazyInstrumentCreationIsConcurrentSafe(t *testing.T) {
	_, provider := newObservabilityMetricReader(t)
	cfg := Config{MeterProvider: provider}
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cfg.ClientLatencyHistogram()
			cfg.BreakerTransitionsCounter()
			cfg.StreamReconnectsCounter()
			cfg.StreamChannelDropsCounter()
		}()
	}
	wg.Wait()
}

func TestSafeErrorTextDoesNotExposeErrorBodies(t *testing.T) {
	errorBody := "app-secret=secret-value access_token=secret-token"
	if got := SafeErrorText(errs.Wrap(errs.CodeServer, errorBody, stderrors.New(errorBody))); got != "webull: SERVER_ERROR" {
		t.Fatalf("SafeErrorText(typed error) = %q, want webull: SERVER_ERROR", got)
	}
	if got := SafeErrorText(context.Canceled); got != context.Canceled.Error() {
		t.Fatalf("SafeErrorText(canceled) = %q, want %q", got, context.Canceled.Error())
	}
	if got := SafeErrorText(stderrors.New(errorBody)); got != "operation failed" {
		t.Fatalf("SafeErrorText(untyped error) = %q, want operation failed", got)
	}
	if got := SafeErrorText(nil); got != "" {
		t.Fatalf("SafeErrorText(nil) = %q, want empty", got)
	}
}

func TestNilAndReaderlessMeterProvidersAreSafe(t *testing.T) {
	previous := otel.GetMeterProvider()
	otel.SetMeterProvider(metricnoop.NewMeterProvider())
	t.Cleanup(func() { otel.SetMeterProvider(previous) })

	var cfg Config
	if cfg.Meter("webullapi4go/test") == nil {
		t.Fatal("Meter() with nil provider = nil")
	}
	ctx := context.Background()
	if cfg.ClientLatencyHistogram().Enabled(ctx) {
		t.Fatal("nil provider enabled client latency histogram")
	}
	if cfg.BreakerTransitionsCounter().Enabled(ctx) {
		t.Fatal("nil provider enabled breaker transition counter")
	}
	if cfg.StreamReconnectsCounter().Enabled(ctx) {
		t.Fatal("nil provider enabled stream reconnect counter")
	}
	if cfg.StreamChannelDropsCounter().Enabled(ctx) {
		t.Fatal("nil provider enabled stream channel drop counter")
	}
	cfg.ClientLatencyHistogram().Record(ctx, 1)
	cfg.BreakerTransitionsCounter().Add(ctx, 1)
	cfg.StreamReconnectsCounter().Add(ctx, 1)
	cfg.StreamChannelDropsCounter().Add(ctx, 1)

	readerless := sdkmetric.NewMeterProvider()
	t.Cleanup(func() { _ = readerless.Shutdown(context.Background()) })
	readerlessCfg := Config{MeterProvider: readerless}
	if readerlessCfg.ClientLatencyHistogram().Enabled(context.Background()) {
		t.Fatal("readerless provider enabled client latency histogram")
	}
	readerlessCfg.BreakerTransitionsCounter().Add(context.Background(), 1)
	readerlessCfg.StreamReconnectsCounter().Add(context.Background(), 1)
	readerlessCfg.StreamChannelDropsCounter().Add(context.Background(), 1)
}
