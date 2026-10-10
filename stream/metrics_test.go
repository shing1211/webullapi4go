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

package stream

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/shing1211/webullapi4go/client"
	marketdatav1 "github.com/shing1211/webullapi4go/gen/webull/marketdata/v1"
)

func newStreamMetricReader(t *testing.T) (*sdkmetric.ManualReader, *sdkmetric.MeterProvider) {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })
	return reader, provider
}

func collectStreamMetrics(t *testing.T, reader *sdkmetric.ManualReader) metricdata.ResourceMetrics {
	t.Helper()
	var rm metricdata.ResourceMetrics
	if err := reader.Collect(context.Background(), &rm); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	return rm
}

func findStreamMetric(t *testing.T, rm metricdata.ResourceMetrics, scope, name string) metricdata.Metrics {
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

func TestStreamRecordsReconnectsAndPerTopicChannelDrops(t *testing.T) {
	reader, provider := newStreamMetricReader(t)
	core, err := client.New(
		client.WithAppKey("test-key"),
		client.WithAppSecret("test-secret"),
		client.WithBaseURL("http://127.0.0.1:1"),
		client.WithMeterProvider(provider),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	defer core.Close()

	s, err := New(core, WithSessionID("metrics-session"), WithMQTTURL("tcp://127.0.0.1:1"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer s.Close()

	_, cancelQuote := s.SubscribeQuoteChan(ChannelConfig{Policy: DropOldest, BufferSize: 1})
	defer cancelQuote()
	_, cancelSnapshot := s.SubscribeSnapshotChan(ChannelConfig{Policy: DropOldest, BufferSize: 1})
	defer cancelSnapshot()
	_, cancelTick := s.SubscribeTickChan(ChannelConfig{Policy: DropOldest, BufferSize: 1})
	defer cancelTick()

	s.emitQuote(&marketdatav1.Quote{Basic: &marketdatav1.Basic{Symbol: "AAPL"}})
	s.emitQuote(&marketdatav1.Quote{Basic: &marketdatav1.Basic{Symbol: "AAPL"}})
	s.emitSnapshot(&marketdatav1.Snapshot{Basic: &marketdatav1.Basic{Symbol: "AAPL"}})
	s.emitSnapshot(&marketdatav1.Snapshot{Basic: &marketdatav1.Basic{Symbol: "AAPL"}})
	s.emitTick(&marketdatav1.Tick{Basic: &marketdatav1.Basic{Symbol: "AAPL"}})
	s.emitTick(&marketdatav1.Tick{Basic: &marketdatav1.Basic{Symbol: "AAPL"}})
	s.handleReconnecting()

	rm := collectStreamMetrics(t, reader)
	reconnects := findStreamMetric(t, rm, "webullapi4go/stream", "reconnects")
	reconnectData, ok := reconnects.Data.(metricdata.Sum[int64])
	if !ok || len(reconnectData.DataPoints) != 1 || reconnectData.DataPoints[0].Value != 1 {
		t.Fatalf("reconnect data = %#v, want one value of 1", reconnects.Data)
	}

	drops := findStreamMetric(t, rm, "webullapi4go/stream", "channel_drops")
	dropData, ok := drops.Data.(metricdata.Sum[int64])
	if !ok {
		t.Fatalf("drop data type = %T, want metricdata.Sum[int64]", drops.Data)
	}
	got := make(map[string]int64)
	for _, point := range dropData.DataPoints {
		topic, ok := point.Attributes.Value(attribute.Key("topic"))
		if !ok {
			t.Fatalf("drop attributes = %v, want topic", point.Attributes.ToSlice())
		}
		got[topic.AsString()] = point.Value
	}
	for _, topic := range []string{"quote", "snapshot", "tick"} {
		if got[topic] != 1 {
			t.Errorf("channel_drops[%s] = %d, want 1", topic, got[topic])
		}
	}
}

func TestStreamWithNilMeterProviderDoesNotPanic(t *testing.T) {
	core, err := client.New(
		client.WithAppKey("test-key"),
		client.WithAppSecret("test-secret"),
		client.WithBaseURL("http://127.0.0.1:1"),
		client.WithMeterProvider(nil),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	defer core.Close()
	s, err := New(core, WithSessionID("nil-meter-session"), WithMQTTURL("tcp://127.0.0.1:1"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer s.Close()
	s.handleReconnecting()
	_, cancel := s.SubscribeQuoteChan(ChannelConfig{Policy: DropOldest, BufferSize: 1})
	cancel()
}
