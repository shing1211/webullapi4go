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

package breaker

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func newBreakerMetricReader(t *testing.T) (*sdkmetric.ManualReader, *sdkmetric.MeterProvider) {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })
	return reader, provider
}

func collectBreakerMetrics(t *testing.T, reader *sdkmetric.ManualReader) metricdata.ResourceMetrics {
	t.Helper()
	var rm metricdata.ResourceMetrics
	if err := reader.Collect(context.Background(), &rm); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	return rm
}

func findBreakerMetric(t *testing.T, rm metricdata.ResourceMetrics, scope, name string) metricdata.Metrics {
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

func TestBreakerTransitionCounterRecordsValuesAndAttributes(t *testing.T) {
	reader, provider := newBreakerMetricReader(t)
	counter, err := provider.Meter("webullapi4go/resilience").Int64Counter("breaker.state_transitions")
	if err != nil {
		t.Fatalf("Int64Counter() error = %v", err)
	}
	b := NewWithCounter(counter, WithThreshold(2), WithCooldown(time.Hour))

	b.RecordFailure()
	b.RecordFailure()
	if b.Allow() {
		t.Fatal("Allow() = true after breaker opened")
	}
	b.Reset()

	metric := findBreakerMetric(t, collectBreakerMetrics(t, reader), "webullapi4go/resilience", "breaker.state_transitions")
	sum, ok := metric.Data.(metricdata.Sum[int64])
	if !ok {
		t.Fatalf("breaker metric data type = %T, want metricdata.Sum[int64]", metric.Data)
	}
	if len(sum.DataPoints) != 2 {
		t.Fatalf("transition data points = %d, want 2", len(sum.DataPoints))
	}
	got := make(map[string]int64)
	for _, point := range sum.DataPoints {
		from, fromOK := point.Attributes.Value(attribute.Key("from"))
		to, toOK := point.Attributes.Value(attribute.Key("to"))
		if !fromOK || !toOK {
			t.Fatalf("transition attributes = %v, want from and to", point.Attributes.ToSlice())
		}
		got[from.AsString()+"->"+to.AsString()] = point.Value
	}
	if got["closed->open"] != 1 || got["open->closed"] != 1 {
		t.Fatalf("transition values = %v, want closed->open=1 and open->closed=1", got)
	}
}
