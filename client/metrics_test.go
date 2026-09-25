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

package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/shing1211/webullapi4go/client"
)

func newClientMetricReader(t *testing.T) (*sdkmetric.ManualReader, *sdkmetric.MeterProvider) {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })
	return reader, provider
}

func collectClientMetrics(t *testing.T, reader *sdkmetric.ManualReader) metricdata.ResourceMetrics {
	t.Helper()
	var rm metricdata.ResourceMetrics
	if err := reader.Collect(context.Background(), &rm); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	return rm
}

func findClientMetric(t *testing.T, rm metricdata.ResourceMetrics, scope, name string) metricdata.Metrics {
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

func hasClientMetric(rm metricdata.ResourceMetrics, scope, name string) bool {
	for _, scopeMetrics := range rm.ScopeMetrics {
		if scopeMetrics.Scope.Name != scope {
			continue
		}
		for _, metric := range scopeMetrics.Metrics {
			if metric.Name == name {
				return true
			}
		}
	}
	return false
}

func TestClientLatencyRecordsOneSamplePerAttempt(t *testing.T) {
	reader, provider := newClientMetricReader(t)
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if attempts.Add(1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{}`))
			return
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	cl, err := client.New(
		client.WithCredentials(testAppKey, testAppSecret),
		client.WithBaseURL(srv.URL),
		client.WithMeterProvider(provider),
		client.WithRetry(client.RetryConfig{
			MaxAttempts: 2,
			BaseDelay:   time.Nanosecond,
			MaxDelay:    time.Nanosecond,
		}),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	defer cl.Close()

	if err := cl.Do(context.Background(), http.MethodGet, "/metrics/test?symbol=AAPL", nil, nil); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if got := attempts.Load(); got != 2 {
		t.Fatalf("HTTP attempts = %d, want 2", got)
	}

	metric := findClientMetric(t, collectClientMetrics(t, reader), "webullapi4go/client", "request_latency")
	data, ok := metric.Data.(metricdata.Histogram[float64])
	if !ok || len(data.DataPoints) != 1 {
		t.Fatalf("latency data = %T with %d points, want one histogram point", metric.Data, len(data.DataPoints))
	}
	point := data.DataPoints[0]
	if point.Count != 2 {
		t.Fatalf("latency sample count = %d, want 2", point.Count)
	}
	if route, ok := point.Attributes.Value(attribute.Key("http.route")); !ok || route.AsString() != "/metrics/test" {
		t.Fatalf("latency route = %v/%t, want /metrics/test/true", route, ok)
	}
}

func TestProductionPresetMeterProviderOrderIsIndependent(t *testing.T) {
	for _, tc := range []struct {
		name string
		opts func(*sdkmetric.MeterProvider) []client.Option
	}{
		{
			name: "provider-before-preset",
			opts: func(provider *sdkmetric.MeterProvider) []client.Option {
				return []client.Option{
					client.WithMeterProvider(provider),
					client.WithResiliencePreset(client.ProductionPreset),
				}
			},
		},
		{
			name: "provider-after-preset",
			opts: func(provider *sdkmetric.MeterProvider) []client.Option {
				return []client.Option{
					client.WithResiliencePreset(client.ProductionPreset),
					client.WithMeterProvider(provider),
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reader, provider := newClientMetricReader(t)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{}`))
			}))
			defer srv.Close()

			opts := []client.Option{
				client.WithCredentials(testAppKey, testAppSecret),
				client.WithBaseURL(srv.URL),
			}
			opts = append(opts, tc.opts(provider)...)
			opts = append(opts, client.WithoutRetry())
			cl, err := client.New(opts...)
			if err != nil {
				t.Fatalf("client.New() error = %v", err)
			}
			defer cl.Close()

			for i := 0; i < 5; i++ {
				if err := cl.Do(context.Background(), http.MethodGet, "/breaker", nil, nil); err == nil {
					t.Fatal("Do() error = nil for a failing response")
				}
			}

			metric := findClientMetric(t, collectClientMetrics(t, reader), "webullapi4go/resilience", "breaker.state_transitions")
			data, ok := metric.Data.(metricdata.Sum[int64])
			if !ok || len(data.DataPoints) != 1 {
				t.Fatalf("breaker metric data = %T with %d points, want one sum point", metric.Data, len(data.DataPoints))
			}
			point := data.DataPoints[0]
			from, fromOK := point.Attributes.Value(attribute.Key("from"))
			to, toOK := point.Attributes.Value(attribute.Key("to"))
			if point.Value != 1 || !fromOK || !toOK || from.AsString() != "closed" || to.AsString() != "open" {
				t.Fatalf("breaker transition = %d %s->%s, want 1 closed->open", point.Value, from.AsString(), to.AsString())
			}
		})
	}
}

type clientMetricsBreaker struct {
	allows   atomic.Int32
	failures atomic.Int32
}

func (b *clientMetricsBreaker) Allow() bool {
	b.allows.Add(1)
	return true
}

func (b *clientMetricsBreaker) RecordSuccess() {}

func (b *clientMetricsBreaker) RecordFailure() {
	b.failures.Add(1)
}

func TestProductionPresetPreservesCustomBreakerOwnership(t *testing.T) {
	for _, customFirst := range []bool{true, false} {
		name := "custom-after-preset"
		if customFirst {
			name = "custom-before-preset"
		}
		t.Run(name, func(t *testing.T) {
			reader, provider := newClientMetricReader(t)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{}`))
			}))
			defer srv.Close()

			custom := &clientMetricsBreaker{}
			preset := client.WithResiliencePreset(client.ProductionPreset)
			breaker := client.WithBreaker(custom)
			ordered := []client.Option{preset, breaker}
			if customFirst {
				ordered = []client.Option{breaker, preset}
			}
			opts := []client.Option{
				client.WithCredentials(testAppKey, testAppSecret),
				client.WithBaseURL(srv.URL),
				client.WithMeterProvider(provider),
			}
			opts = append(opts, ordered...)
			opts = append(opts, client.WithoutRetry())
			cl, err := client.New(opts...)
			if err != nil {
				t.Fatalf("client.New() error = %v", err)
			}
			defer cl.Close()

			for i := 0; i < 5; i++ {
				_ = cl.Do(context.Background(), http.MethodGet, "/breaker", nil, nil)
			}
			if got := custom.allows.Load(); got != 5 {
				t.Fatalf("custom breaker allows = %d, want 5", got)
			}
			if got := custom.failures.Load(); got != 5 {
				t.Fatalf("custom breaker failures = %d, want 5", got)
			}
			if hasClientMetric(collectClientMetrics(t, reader), "webullapi4go/resilience", "breaker.state_transitions") {
				t.Fatal("production breaker metric was emitted for a custom breaker")
			}
		})
	}
}

func TestClientWithNilMeterProviderDoesNotPanic(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()
	cl, err := client.New(
		client.WithCredentials(testAppKey, testAppSecret),
		client.WithBaseURL(srv.URL),
		client.WithMeterProvider(nil),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	defer cl.Close()
	if err := cl.Do(context.Background(), http.MethodGet, "/nil-meter", nil, nil); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
}
