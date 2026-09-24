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
	"bytes"
	"context"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/events"
	eventsevents "github.com/shing1211/webullapi4go/gen/webull/trade/events/v1"
)

type metadataServer struct {
	eventsevents.UnimplementedEventServiceServer
	mu sync.Mutex
	md metadata.MD
}

func (s *metadataServer) Subscribe(_ *eventsevents.SubscribeRequest, stream grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error {
	md, _ := metadata.FromIncomingContext(stream.Context())
	s.mu.Lock()
	s.md = md
	s.mu.Unlock()
	return nil
}

func (s *metadataServer) captured() metadata.MD {
	s.mu.Lock()
	defer s.mu.Unlock()
	captured := metadata.MD{}
	for key, values := range s.md {
		captured[key] = append([]string(nil), values...)
	}
	return captured
}

type countingTracerProvider struct {
	trace.TracerProvider
	tracer *countingTracer
}

func (p *countingTracerProvider) Tracer(string, ...trace.TracerOption) trace.Tracer {
	return p.tracer
}

type countingTracer struct {
	trace.Tracer
	starts atomic.Int32
	ends   atomic.Int32
	mu     sync.Mutex
	spans  []*countingSpan
}

func (t *countingTracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	t.starts.Add(1)
	cfg := trace.NewSpanStartConfig(opts...)
	ctx, span := t.Tracer.Start(ctx, name, opts...)
	recording := &countingSpan{
		Span:  span,
		owner: t,
		name:  name,
		kind:  cfg.SpanKind(),
		attrs: append([]attribute.KeyValue(nil), cfg.Attributes()...),
	}
	t.mu.Lock()
	t.spans = append(t.spans, recording)
	t.mu.Unlock()
	return ctx, recording
}

func (t *countingTracer) snapshot() []*countingSpan {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]*countingSpan(nil), t.spans...)
}

type countingSpan struct {
	trace.Span
	owner      *countingTracer
	name       string
	kind       trace.SpanKind
	attrs      []attribute.KeyValue
	mu         sync.Mutex
	status     otelcodes.Code
	errorCount int
	ended      bool
}

func (s *countingSpan) SetStatus(code otelcodes.Code, description string) {
	s.mu.Lock()
	s.status = code
	s.mu.Unlock()
	s.Span.SetStatus(code, description)
}

func (s *countingSpan) RecordError(err error, opts ...trace.EventOption) {
	s.mu.Lock()
	s.errorCount++
	s.mu.Unlock()
	s.Span.RecordError(err, opts...)
}

func (s *countingSpan) SetAttributes(attrs ...attribute.KeyValue) {
	s.mu.Lock()
	s.attrs = append(s.attrs, attrs...)
	s.mu.Unlock()
	s.Span.SetAttributes(attrs...)
}

func (s *countingSpan) End(opts ...trace.SpanEndOption) {
	s.mu.Lock()
	if !s.ended {
		s.ended = true
		s.owner.ends.Add(1)
	}
	s.mu.Unlock()
	s.Span.End(opts...)
}

func newCountingTracerProvider() *countingTracerProvider {
	base := noop.NewTracerProvider()
	return &countingTracerProvider{
		TracerProvider: base,
		tracer:         &countingTracer{Tracer: base.Tracer("test")},
	}
}

func TestRunPropagatesMetadataAndRecordsOneAttempt(t *testing.T) {
	provider := newCountingTracerProvider()
	var logs bytes.Buffer
	core := newCoreWithOptions(t,
		client.WithTracerProvider(provider),
		client.WithPropagator(propagation.TraceContext{}),
		client.WithLogger(slog.New(slog.NewJSONHandler(&logs, nil))),
	)
	server := &metadataServer{}
	cl := newClientWithCore(t, core, server, events.WithAutoReconnect(false))

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
	ctx := trace.ContextWithRemoteSpanContext(context.Background(), parent)
	ctx = client.WithCorrelationID(ctx, "corr-1")
	if err := cl.Run(ctx); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	md := server.captured()
	if got := md.Get("x-correlation-id"); len(got) != 1 || got[0] != "corr-1" {
		t.Fatalf("correlation metadata = %v, want [corr-1]", got)
	}
	if got := md.Get("traceparent"); len(got) != 1 || !strings.Contains(got[0], traceID.String()) || !strings.Contains(got[0], spanID.String()) {
		t.Fatalf("traceparent = %v, want parent trace and span IDs", got)
	}
	for _, key := range []string{
		"x-app-key",
		"x-signature-algorithm",
		"x-signature-version",
		"x-signature-nonce",
		"x-timestamp",
		"x-signature",
	} {
		if len(md.Get(key)) != 1 {
			t.Errorf("signed metadata %q = %v, want one value", key, md.Get(key))
		}
	}
	if strings.Contains(logs.String(), "test-app-secret") || strings.Contains(logs.String(), "test-app-key") {
		t.Errorf("logs contain credentials: %s", logs.String())
	}

	spans := provider.tracer.snapshot()
	if len(spans) != 1 {
		t.Fatalf("spans = %d, want 1", len(spans))
	}
	if spans[0].kind != trace.SpanKindClient {
		t.Errorf("span kind = %v, want client", spans[0].kind)
	}
	if spans[0].name != "/grpc.trade.event.EventService/Subscribe" {
		t.Errorf("span name = %q", spans[0].name)
	}
	hasCorrelation := false
	for _, attr := range spans[0].attrs {
		if attr.Key == "webull.correlation_id" && attr.Value.AsString() == "corr-1" {
			hasCorrelation = true
		}
	}
	if !hasCorrelation {
		t.Error("span is missing the correlation ID attribute")
	}
	if provider.tracer.starts.Load() != 1 || provider.tracer.ends.Load() != 1 {
		t.Errorf("span starts/ends = %d/%d, want 1/1", provider.tracer.starts.Load(), provider.tracer.ends.Load())
	}
	spans[0].mu.Lock()
	status, ended, errorCount := spans[0].status, spans[0].ended, spans[0].errorCount
	spans[0].mu.Unlock()
	if status != otelcodes.Ok || !ended || errorCount != 0 {
		t.Errorf("span status/ended/errors = %v/%v/%d, want OK/true/0", status, ended, errorCount)
	}
}
