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
	"errors"
	"log/slog"
	"strings"
	"sync"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"

	"github.com/shing1211/webullapi4go/client"
	marketdatav1 "github.com/shing1211/webullapi4go/gen/webull/marketdata/v1"
	mqtt "github.com/shing1211/webullapi4go/pkg/transport/mqtt"
)

type streamV1SlogRecord struct {
	level   slog.Level
	message string
	attrs   map[string]slog.Value
}

type streamV1SlogHandler struct {
	mu      sync.Mutex
	records []streamV1SlogRecord
}

func (h *streamV1SlogHandler) Enabled(context.Context, slog.Level) bool { return true }

func (h *streamV1SlogHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := make(map[string]slog.Value, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value
		return true
	})
	h.mu.Lock()
	h.records = append(h.records, streamV1SlogRecord{level: record.Level, message: record.Message, attrs: attrs})
	h.mu.Unlock()
	return nil
}

func (h *streamV1SlogHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *streamV1SlogHandler) WithGroup(string) slog.Handler      { return h }

func (h *streamV1SlogHandler) snapshot() []streamV1SlogRecord {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]streamV1SlogRecord, len(h.records))
	copy(out, h.records)
	return out
}

func streamV1NewTracer(t *testing.T) (*sdktrace.TracerProvider, *tracetest.SpanRecorder) {
	t.Helper()
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })
	return provider, recorder
}

func streamV1NewClient(t *testing.T, provider *sdktrace.TracerProvider, logger *slog.Logger) *Client {
	t.Helper()
	core, err := client.New(
		client.WithAppKey("stream-v1-app-key"),
		client.WithAppSecret("stream-v1-app-secret"),
		client.WithBaseURL("http://127.0.0.1:1"),
		client.WithTracerProvider(provider),
		client.WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = core.Close() })
	streamClient, err := New(core, WithSessionID("stream-v1-session"), WithMQTTURL("tcp://127.0.0.1:1"), WithAutoReconnect(true))
	if err != nil {
		t.Fatalf("stream.New() error = %v", err)
	}
	t.Cleanup(func() { _ = streamClient.Close() })
	return streamClient
}

func streamV1SpanAttrs(span sdktrace.ReadOnlySpan) map[string]attribute.Value {
	attrs := make(map[string]attribute.Value, len(span.Attributes()))
	for _, attr := range span.Attributes() {
		attrs[string(attr.Key)] = attr.Value
	}
	return attrs
}

func streamV1AssertStringAttr(t *testing.T, span sdktrace.ReadOnlySpan, key, want string) {
	t.Helper()
	got, ok := streamV1SpanAttrs(span)[key]
	if !ok || got.AsString() != want {
		t.Fatalf("span %q attribute %s = %v/%t, want %q/true", span.Name(), key, got, ok, want)
	}
}

func streamV1AssertIntAttr(t *testing.T, span sdktrace.ReadOnlySpan, key string, want int64) {
	t.Helper()
	got, ok := streamV1SpanAttrs(span)[key]
	if !ok || got.AsInt64() != want {
		t.Fatalf("span %q attribute %s = %v/%t, want %d/true", span.Name(), key, got, ok, want)
	}
}

func streamV1SpanText(span sdktrace.ReadOnlySpan) string {
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

func streamV1AssertLog(t *testing.T, record streamV1SlogRecord, level slog.Level, message, state string) {
	t.Helper()
	if record.level != level || record.message != message {
		t.Fatalf("log = %+v, want level %v/message %q", record, level, message)
	}
	want := map[string]string{"transport": "mqtt", "state": state}
	if len(record.attrs) != len(want) {
		t.Fatalf("log %q attributes = %v, want exactly %v", message, record.attrs, want)
	}
	for key, value := range want {
		got, ok := record.attrs[key]
		if !ok || got.Kind() != slog.KindString || got.String() != value {
			t.Fatalf("log %q attribute %s = %v/%t, want %q/true", message, key, got, ok, value)
		}
	}
}

func TestV1MQTTDispatchSpansDataNonDataAndError(t *testing.T) {
	provider, recorder := streamV1NewTracer(t)
	client := streamV1NewClient(t, provider, nil)
	quotePayload, err := proto.Marshal(&marketdatav1.Quote{Basic: &marketdatav1.Basic{Symbol: "AAPL"}})
	if err != nil {
		t.Fatalf("proto.Marshal() error = %v", err)
	}

	cases := []struct {
		name        string
		topic       string
		payload     []byte
		destination string
		wantStatus  otelcodes.Code
		wantError   bool
	}{
		{name: "data", topic: TopicQuote, payload: quotePayload, destination: "quote", wantStatus: otelcodes.Ok},
		{name: "non-data", topic: TopicNotice, payload: []byte(`{"type":"status"}`), destination: "notice", wantStatus: otelcodes.Ok},
		{name: "error", topic: "v1-sensitive-topic", payload: []byte("stream-v1-app-secret"), destination: "unknown", wantStatus: otelcodes.Error, wantError: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client.handleMessage(mqttMessage(tc.topic, tc.payload))
		})
	}

	spans := recorder.Ended()
	if len(spans) != len(cases) {
		t.Fatalf("dispatch spans = %d, want %d", len(spans), len(cases))
	}
	for i, tc := range cases {
		span := spans[i]
		if span.Name() != "mqtt.dispatch" || span.SpanKind() != trace.SpanKindConsumer {
			t.Fatalf("span %d = %q/%v, want mqtt.dispatch/consumer", i, span.Name(), span.SpanKind())
		}
		if span.Parent().IsValid() {
			t.Fatalf("dispatch span %d has application parent %s/%s", i, span.Parent().TraceID(), span.Parent().SpanID())
		}
		if span.EndTime().IsZero() || !span.SpanContext().IsValid() {
			t.Fatalf("dispatch span %d did not complete", i)
		}
		if span.Status().Code != tc.wantStatus {
			t.Fatalf("dispatch span %d status = %v, want %v", i, span.Status().Code, tc.wantStatus)
		}
		if span.InstrumentationScope().Name != "webullapi4go/stream" {
			t.Fatalf("dispatch span %d scope = %q, want webullapi4go/stream", i, span.InstrumentationScope().Name)
		}
		streamV1AssertStringAttr(t, span, "messaging.system", "mqtt")
		streamV1AssertStringAttr(t, span, "messaging.destination", tc.destination)
		streamV1AssertIntAttr(t, span, "messaging.message.body.size", int64(len(tc.payload)))
		text := streamV1SpanText(span)
		if tc.wantError {
			if !strings.Contains(text, "operation failed") {
				t.Fatalf("error dispatch span text = %q, want safe error text", text)
			}
			if strings.Contains(text, "stream-v1-app-secret") {
				t.Fatalf("dispatch span contains sensitive payload: %s", text)
			}
		} else if strings.Contains(text, "error") {
			t.Fatalf("successful dispatch span contains error text: %s", text)
		}
	}
}

func TestV1StreamLifecycleLogsConnectDisconnectAndReconnect(t *testing.T) {
	provider, _ := streamV1NewTracer(t)
	handler := &streamV1SlogHandler{}
	client := streamV1NewClient(t, provider, slog.New(handler))

	client.handleConnect()
	client.handleConnectionLost(errors.New("v1 disconnect"))
	client.handleReconnecting()

	records := handler.snapshot()
	if len(records) != 3 {
		t.Fatalf("lifecycle log count = %d, want 3", len(records))
	}
	streamV1AssertLog(t, records[0], slog.LevelInfo, "webull stream connected", "connected")
	streamV1AssertLog(t, records[1], slog.LevelWarn, "webull stream disconnected", "disconnected")
	streamV1AssertLog(t, records[2], slog.LevelInfo, "webull stream reconnecting", "reconnecting")
}

func TestV1StreamNilLoggerIsSafe(t *testing.T) {
	provider, _ := streamV1NewTracer(t)
	client := streamV1NewClient(t, provider, nil)
	client.handleConnect()
	client.handleConnectionLost(errors.New("v1 disconnect"))
	client.handleReconnecting()
}

func mqttMessage(topic string, payload []byte) mqtt.Message {
	return mqtt.Message{Topic: topic, Payload: payload}
}
