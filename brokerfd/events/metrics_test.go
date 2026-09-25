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
	"testing"

	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/shing1211/webullapi4go/client"
	eventsevents "github.com/shing1211/webullapi4go/gen/webull/brokerfd/events/v1"
)

func newBrokerEventMetricReader(t *testing.T) (*sdkmetric.ManualReader, *sdkmetric.MeterProvider) {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })
	return reader, provider
}

func collectBrokerEventMetrics(t *testing.T, reader *sdkmetric.ManualReader) metricdata.ResourceMetrics {
	t.Helper()
	var rm metricdata.ResourceMetrics
	if err := reader.Collect(context.Background(), &rm); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	return rm
}

func findBrokerEventMetric(t *testing.T, rm metricdata.ResourceMetrics, scope, name string) metricdata.Metrics {
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

func assertBrokerEventAttemptMetric(t *testing.T, metric metricdata.Metrics) {
	t.Helper()
	data, ok := metric.Data.(metricdata.Sum[int64])
	if !ok || len(data.DataPoints) != 1 || data.DataPoints[0].Value != 1 {
		t.Fatalf("broker event attempt data = %#v, want one value of 1", metric.Data)
	}
	attrs := data.DataPoints[0].Attributes
	for key, want := range map[string]string{
		"rpc.system":  "grpc",
		"rpc.service": "grpc.event.EventService",
		"rpc.method":  "Subscribe",
	} {
		if got, ok := attrs.Value(attribute.Key(key)); !ok || got.AsString() != want {
			t.Errorf("attempt attribute %s = %v/%t, want %s/true", key, got, ok, want)
		}
	}
}

func assertBrokerEventDurationMetric(t *testing.T, metric metricdata.Metrics, status int64, outcome string) {
	t.Helper()
	data, ok := metric.Data.(metricdata.Histogram[float64])
	if !ok || len(data.DataPoints) != 1 || data.DataPoints[0].Count != 1 {
		t.Fatalf("broker event duration data = %#v, want one sample", metric.Data)
	}
	attrs := data.DataPoints[0].Attributes
	for key, want := range map[string]string{
		"rpc.system":  "grpc",
		"rpc.service": "grpc.event.EventService",
		"rpc.method":  "Subscribe",
	} {
		if got, ok := attrs.Value(attribute.Key(key)); !ok || got.AsString() != want {
			t.Errorf("duration attribute %s = %v/%t, want %s/true", key, got, ok, want)
		}
	}
	if got, ok := attrs.Value(attribute.Key("rpc.grpc.status_code")); !ok || got.AsInt64() != status {
		t.Errorf("duration status = %v/%t, want %d/true", got, ok, status)
	}
	if got, ok := attrs.Value(attribute.Key("webull.stream.outcome")); !ok || got.AsString() != outcome {
		t.Errorf("duration outcome = %v/%t, want %s/true", got, ok, outcome)
	}
}

func TestBrokerEventMetricsRecordAttemptAndDurationForSuccess(t *testing.T) {
	reader, provider := newBrokerEventMetricReader(t)
	core := newCoreWithOptions(t, client.WithMeterProvider(provider))
	cl := newClientWithCore(t, core, &metadataServer{}, WithAutoReconnect(false))
	if err := cl.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	rm := collectBrokerEventMetrics(t, reader)
	assertBrokerEventAttemptMetric(t, findBrokerEventMetric(t, rm, "webullapi4go/brokerfd/events", "event_stream_attempts"))
	assertBrokerEventDurationMetric(t, findBrokerEventMetric(t, rm, "webullapi4go/brokerfd/events", "event_stream_attempt_duration"), 0, "ok")
}

func TestBrokerEventMetricsRecordAttemptAndDurationForCancellation(t *testing.T) {
	reader, provider := newBrokerEventMetricReader(t)
	core := newCoreWithOptions(t, client.WithMeterProvider(provider))
	server := &fakeServer{messages: []*eventsevents.SubscribeResponse{{EventType: eventsevents.EventType_SubscribeSuccess}}}
	cl := newClientWithCore(t, core, server, WithAutoReconnect(false))
	connected := make(chan struct{}, 1)
	cl.OnConnect(func() { trySend(connected, struct{}{}) })

	run := startTestRun(t, cl, context.Background())
	waitSignal(t, connected, "OnConnect")
	if err := run.stop(t); !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want context.Canceled", err)
	}

	rm := collectBrokerEventMetrics(t, reader)
	assertBrokerEventAttemptMetric(t, findBrokerEventMetric(t, rm, "webullapi4go/brokerfd/events", "event_stream_attempts"))
	assertBrokerEventDurationMetric(t, findBrokerEventMetric(t, rm, "webullapi4go/brokerfd/events", "event_stream_attempt_duration"), 1, "error")
}

func TestBrokerEventWithNilMeterProviderDoesNotPanic(t *testing.T) {
	core := newCoreWithOptions(t, client.WithMeterProvider(nil))
	cl := newClientWithCore(t, core, &metadataServer{}, WithAutoReconnect(false))
	if err := cl.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
}
