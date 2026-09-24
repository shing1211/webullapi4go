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
	"strings"
	"testing"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
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
