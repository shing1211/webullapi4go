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
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	"github.com/shing1211/webullapi4go/client"
)

type v1SlogRecord struct {
	level   slog.Level
	message string
	attrs   map[string]slog.Value
}

type v1SlogHandler struct {
	mu      sync.Mutex
	records []v1SlogRecord
}

func (h *v1SlogHandler) Enabled(context.Context, slog.Level) bool { return true }

func (h *v1SlogHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := make(map[string]slog.Value, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value
		return true
	})
	h.mu.Lock()
	h.records = append(h.records, v1SlogRecord{level: record.Level, message: record.Message, attrs: attrs})
	h.mu.Unlock()
	return nil
}

func (h *v1SlogHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *v1SlogHandler) WithGroup(string) slog.Handler      { return h }

func (h *v1SlogHandler) snapshot() []v1SlogRecord {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]v1SlogRecord, len(h.records))
	copy(out, h.records)
	return out
}

func v1NewTracer(t *testing.T) (*sdktrace.TracerProvider, *tracetest.SpanRecorder) {
	t.Helper()
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })
	return provider, recorder
}

func v1RemoteParent(t *testing.T) (context.Context, trace.SpanContext) {
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

type v1HeaderCapture struct {
	mu     sync.Mutex
	values []http.Header
}

func (c *v1HeaderCapture) add(header http.Header) {
	copyHeader := header.Clone()
	c.mu.Lock()
	c.values = append(c.values, copyHeader)
	c.mu.Unlock()
}

func (c *v1HeaderCapture) snapshot() []http.Header {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]http.Header, len(c.values))
	for i, header := range c.values {
		out[i] = header.Clone()
	}
	return out
}

func v1SpanAttrs(span sdktrace.ReadOnlySpan) map[string]attribute.Value {
	attrs := make(map[string]attribute.Value, len(span.Attributes()))
	for _, attr := range span.Attributes() {
		attrs[string(attr.Key)] = attr.Value
	}
	return attrs
}

func v1AssertStringAttr(t *testing.T, span sdktrace.ReadOnlySpan, key, want string) {
	t.Helper()
	attrs := v1SpanAttrs(span)
	got, ok := attrs[key]
	if !ok || got.AsString() != want {
		t.Fatalf("span %q attribute %s = %v/%t, want %q/true", span.Name(), key, got, ok, want)
	}
}

func v1AssertIntAttr(t *testing.T, span sdktrace.ReadOnlySpan, key string, want int64) {
	t.Helper()
	attrs := v1SpanAttrs(span)
	got, ok := attrs[key]
	if !ok || got.AsInt64() != want {
		t.Fatalf("span %q attribute %s = %v/%t, want %d/true", span.Name(), key, got, ok, want)
	}
}

func v1AssertSlogValue(t *testing.T, record v1SlogRecord, key string, want any) {
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
		t.Fatalf("unsupported slog expectation %T", want)
	}
}

func v1AssertSlogAttrs(t *testing.T, record v1SlogRecord, want map[string]any) {
	t.Helper()
	if len(record.attrs) != len(want) {
		t.Fatalf("log %q attributes = %v, want exactly %v", record.message, record.attrs, want)
	}
	for key, value := range want {
		v1AssertSlogValue(t, record, key, value)
	}
}

func v1SpanText(span sdktrace.ReadOnlySpan) string {
	var builder strings.Builder
	builder.WriteString(span.Name())
	builder.WriteString(span.Status().Description)
	for _, attr := range span.Attributes() {
		builder.WriteString(string(attr.Key))
		builder.WriteString(fmt.Sprint(attr.Value.AsInterface()))
	}
	for _, event := range span.Events() {
		builder.WriteString(event.Name)
		for _, attr := range event.Attributes {
			builder.WriteString(string(attr.Key))
			builder.WriteString(fmt.Sprint(attr.Value.AsInterface()))
		}
	}
	return builder.String()
}

func v1SlogText(records []v1SlogRecord) string {
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

func v1AssertTraceparent(t *testing.T, header http.Header, parent trace.SpanContext, child trace.SpanContext) {
	t.Helper()
	raw := header.Get("traceparent")
	if raw == "" {
		t.Fatal("traceparent header is empty")
	}
	extracted := propagation.TraceContext{}.Extract(context.Background(), propagation.HeaderCarrier(header))
	got := trace.SpanContextFromContext(extracted)
	if got.TraceID() != parent.TraceID() || got.SpanID() != child.SpanID() {
		t.Fatalf("traceparent = %q, extracted trace/span = %s/%s, want %s/%s", raw, got.TraceID(), got.SpanID(), parent.TraceID(), child.SpanID())
	}
}

func TestV1RESTEntrypointsRecordTraceAndCorrelation(t *testing.T) {
	cases := []struct {
		name string
		call func(context.Context, *client.Client) error
	}{
		{
			name: "do",
			call: func(ctx context.Context, cl *client.Client) error {
				return cl.Do(ctx, http.MethodGet, "/v1/accounts?page=1", nil, nil)
			},
		},
		{
			name: "broker",
			call: func(ctx context.Context, cl *client.Client) error {
				return cl.DoBroker(ctx, http.MethodGet, "/v1/accounts?page=1", nil, nil)
			},
		},
		{
			name: "stream",
			call: func(ctx context.Context, cl *client.Client) error {
				response, err := cl.DoStream(ctx, http.MethodGet, "/v1/accounts?page=1", nil)
				if response != nil {
					_ = response.Body.Close()
				}
				return err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			provider, recorder := v1NewTracer(t)
			handler := &v1SlogHandler{}
			capture := &v1HeaderCapture{}
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capture.add(r.Header)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{}`))
			}))
			t.Cleanup(srv.Close)

			cl, err := client.New(
				client.WithCredentials("v1-app-key", "v1-app-secret"),
				client.WithEndpoints(client.Endpoints{HTTP: srv.URL, BrokerHTTP: srv.URL}),
				client.WithTracerProvider(provider),
				client.WithPropagator(propagation.TraceContext{}),
				client.WithLogger(slog.New(handler)),
				client.WithoutRetry(),
			)
			if err != nil {
				t.Fatalf("client.New() error = %v", err)
			}
			t.Cleanup(func() { _ = cl.Close() })
			cl.SetToken(&client.Token{Value: "v1-access-token", Status: client.TokenStatusNormal})

			parentCtx, parent := v1RemoteParent(t)
			ctx := client.WithCorrelationID(parentCtx, "v1-correlation")
			if err := tc.call(ctx, cl); err != nil {
				t.Fatalf("%s() error = %v", tc.name, err)
			}

			headers := capture.snapshot()
			if len(headers) != 1 {
				t.Fatalf("request count = %d, want 1", len(headers))
			}
			if got := headers[0].Get("x-correlation-id"); got != "v1-correlation" {
				t.Fatalf("x-correlation-id = %q, want v1-correlation", got)
			}
			spans := recorder.Ended()
			if len(spans) != 1 {
				t.Fatalf("ended spans = %d, want 1", len(spans))
			}
			span := spans[0]
			if span.Name() != "GET /v1/accounts" {
				t.Fatalf("span name = %q, want GET /v1/accounts", span.Name())
			}
			if span.SpanKind() != trace.SpanKindClient {
				t.Fatalf("span kind = %v, want client", span.SpanKind())
			}
			if span.InstrumentationScope().Name != "webullapi4go/client" {
				t.Fatalf("span scope = %q, want webullapi4go/client", span.InstrumentationScope().Name)
			}
			if span.Parent().TraceID() != parent.TraceID() || span.Parent().SpanID() != parent.SpanID() {
				t.Fatalf("span parent = %s/%s, want %s/%s", span.Parent().TraceID(), span.Parent().SpanID(), parent.TraceID(), parent.SpanID())
			}
			if !span.SpanContext().IsValid() || span.EndTime().IsZero() {
				t.Fatal("REST span did not complete with a valid span context")
			}
			v1AssertStringAttr(t, span, "http.method", http.MethodGet)
			v1AssertStringAttr(t, span, "http.route", "/v1/accounts")
			v1AssertStringAttr(t, span, "webull.correlation_id", "v1-correlation")
			v1AssertIntAttr(t, span, "webull.attempt", 1)
			v1AssertIntAttr(t, span, "http.status_code", http.StatusOK)
			if got := v1SpanAttrs(span)["http.duration_ms"]; got.Type() != attribute.INT64 || got.AsInt64() < 0 {
				t.Fatalf("http.duration_ms = %v, want a non-negative duration", got)
			}
			if _, ok := v1SpanAttrs(span)["error"]; ok {
				t.Fatal("successful REST span contains an error attribute")
			}
			if span.Status().Code != otelcodes.Ok {
				t.Fatalf("span status = %v, want OK", span.Status())
			}
			v1AssertTraceparent(t, headers[0], parent, span.SpanContext())

			records := handler.snapshot()
			if len(records) != 2 {
				t.Fatalf("log record count = %d, want start and done", len(records))
			}
			if records[0].level != slog.LevelInfo || records[0].message != "webull request start" {
				t.Fatalf("first log = %+v, want INFO start", records[0])
			}
			v1AssertSlogAttrs(t, records[0], map[string]any{
				"method":         http.MethodGet,
				"path":           "/v1/accounts",
				"attempt":        1,
				"correlation_id": "v1-correlation",
			})
			if records[1].level != slog.LevelInfo || records[1].message != "webull request done" {
				t.Fatalf("second log = %+v, want INFO done", records[1])
			}
			v1AssertSlogAttrs(t, records[1], map[string]any{
				"method":         http.MethodGet,
				"path":           "/v1/accounts",
				"attempt":        1,
				"status":         http.StatusOK,
				"correlation_id": "v1-correlation",
				"latency":        time.Duration(0),
			})
			for _, secret := range []string{"v1-app-key", "v1-app-secret", "v1-access-token"} {
				if strings.Contains(v1SlogText(records), secret) {
					t.Fatalf("logs contain credential %q: %s", secret, v1SlogText(records))
				}
			}
		})
	}
}

func TestV1RESTEntrypointsRecordFailureTelemetryAndRedaction(t *testing.T) {
	const (
		failureKey      = "v1-failure-app-key"
		failureValue    = "v1-failure-app-secret"
		failureMaterial = "v1-failure-access-token"
		correlation     = "v1-failure-correlation"
		route           = "/v1/failure"
	)
	cases := []struct {
		name string
		call func(context.Context, *client.Client) error
	}{
		{
			name: "do",
			call: func(ctx context.Context, cl *client.Client) error {
				return cl.Do(ctx, http.MethodGet, route, nil, nil)
			},
		},
		{
			name: "broker",
			call: func(ctx context.Context, cl *client.Client) error {
				return cl.DoBroker(ctx, http.MethodGet, route, nil, nil)
			},
		},
		{
			name: "stream",
			call: func(ctx context.Context, cl *client.Client) error {
				response, err := cl.DoStream(ctx, http.MethodGet, route, nil)
				if response != nil {
					_ = response.Body.Close()
				}
				return err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			provider, recorder := v1NewTracer(t)
			handler := &v1SlogHandler{}
			capture := &v1HeaderCapture{}
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capture.add(r.Header)
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"message": "app-key=" + failureKey + " app-secret=" + failureValue + " access-token=" + failureMaterial + " x-signature=" + r.Header.Get("x-signature"),
				})
			}))
			t.Cleanup(srv.Close)
			cl, err := client.New(
				client.WithCredentials(failureKey, failureValue),
				client.WithEndpoints(client.Endpoints{HTTP: srv.URL, BrokerHTTP: srv.URL}),
				client.WithTracerProvider(provider),
				client.WithPropagator(propagation.TraceContext{}),
				client.WithLogger(slog.New(handler)),
				client.WithoutRetry(),
			)
			if err != nil {
				t.Fatalf("client.New() error = %v", err)
			}
			t.Cleanup(func() { _ = cl.Close() })
			cl.SetToken(&client.Token{Value: failureMaterial, Status: client.TokenStatusNormal})
			parentCtx, parent := v1RemoteParent(t)
			ctx := client.WithCorrelationID(parentCtx, correlation)
			if err := tc.call(ctx, cl); err == nil {
				t.Fatalf("%s() error = nil, want failure", tc.name)
			}

			headers := capture.snapshot()
			if len(headers) != 1 {
				t.Fatalf("request count = %d, want 1", len(headers))
			}
			if got := headers[0].Get("x-correlation-id"); got != correlation {
				t.Fatalf("x-correlation-id = %q, want %q", got, correlation)
			}
			spans := recorder.Ended()
			if len(spans) != 1 {
				t.Fatalf("ended spans = %d, want 1", len(spans))
			}
			span := spans[0]
			if span.Name() != "GET "+route || span.SpanKind() != trace.SpanKindClient {
				t.Fatalf("span = %q/%v, want GET route/client", span.Name(), span.SpanKind())
			}
			if span.Parent().TraceID() != parent.TraceID() || span.Parent().SpanID() != parent.SpanID() || span.EndTime().IsZero() {
				t.Fatalf("span parent/completion = %s/%s/%v", span.Parent().TraceID(), span.Parent().SpanID(), span.EndTime())
			}
			v1AssertStringAttr(t, span, "http.method", http.MethodGet)
			v1AssertStringAttr(t, span, "http.route", route)
			v1AssertStringAttr(t, span, "webull.correlation_id", correlation)
			v1AssertIntAttr(t, span, "webull.attempt", 1)
			v1AssertIntAttr(t, span, "http.status_code", http.StatusInternalServerError)
			v1AssertStringAttr(t, span, "error", "webull: SERVER_ERROR")
			if got := v1SpanAttrs(span)["http.duration_ms"]; got.Type() != attribute.INT64 || got.AsInt64() < 0 {
				t.Fatalf("http.duration_ms = %v, want non-negative duration", got)
			}
			if span.Status().Code != otelcodes.Error || len(span.Events()) == 0 {
				t.Fatalf("span status/events = %v/%d, want ERROR/non-empty", span.Status(), len(span.Events()))
			}
			v1AssertTraceparent(t, headers[0], parent, span.SpanContext())

			logs := handler.snapshot()
			if len(logs) != 2 {
				t.Fatalf("log count = %d, want start and done", len(logs))
			}
			if logs[0].level != slog.LevelInfo || logs[0].message != "webull request start" {
				t.Fatalf("start log = %+v", logs[0])
			}
			v1AssertSlogAttrs(t, logs[0], map[string]any{
				"method":         http.MethodGet,
				"path":           route,
				"attempt":        1,
				"correlation_id": correlation,
			})
			if logs[1].level != slog.LevelError || logs[1].message != "webull request done" {
				t.Fatalf("failure log = %+v", logs[1])
			}
			v1AssertSlogAttrs(t, logs[1], map[string]any{
				"method":         http.MethodGet,
				"path":           route,
				"attempt":        1,
				"status":         http.StatusInternalServerError,
				"correlation_id": correlation,
				"latency":        time.Duration(0),
				"error":          "webull: SERVER_ERROR",
			})
			text := v1SlogText(logs) + v1SpanText(span)
			for _, secret := range []string{failureKey, failureValue, failureMaterial, headers[0].Get("x-signature")} {
				if strings.Contains(text, secret) {
					t.Fatalf("failure telemetry contains sensitive value %q: %s", secret, text)
				}
			}
		})
	}
}

func TestV1RESTRetryLogsAndSpansRedactSensitiveErrors(t *testing.T) {
	const (
		appKey         = "v1-retry-app-key"
		firstMaterial  = "v1-retry-app-secret"
		secondMaterial = "v1-retry-access-token"
		thirdMaterial  = "v1-retry-query-token"
		correlation    = "v1-retry-correlation"
		route          = "/v1/retry"
	)
	cases := []struct {
		name string
		call func(context.Context, *client.Client) error
	}{
		{
			name: "do",
			call: func(ctx context.Context, cl *client.Client) error {
				return cl.Do(ctx, http.MethodGet, route+"?access_token="+thirdMaterial, nil, nil)
			},
		},
		{
			name: "broker",
			call: func(ctx context.Context, cl *client.Client) error {
				return cl.DoBroker(ctx, http.MethodGet, route+"?access_token="+thirdMaterial, nil, nil)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			provider, recorder := v1NewTracer(t)
			handler := &v1SlogHandler{}
			var attempts atomic.Int32
			var signaturesMu sync.Mutex
			var signatures []string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempt := attempts.Add(1)
				signaturesMu.Lock()
				signatures = append(signatures, r.Header.Get("x-signature"))
				signaturesMu.Unlock()
				if attempt == 1 {
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(map[string]string{
						"message": "app-key=" + appKey + " app-secret=" + firstMaterial + " access-token=" + secondMaterial + " x-signature=" + r.Header.Get("x-signature"),
					})
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{}`))
			}))
			t.Cleanup(srv.Close)

			cl, err := client.New(
				client.WithCredentials(appKey, firstMaterial),
				client.WithEndpoints(client.Endpoints{HTTP: srv.URL, BrokerHTTP: srv.URL}),
				client.WithTracerProvider(provider),
				client.WithPropagator(propagation.TraceContext{}),
				client.WithLogger(slog.New(handler)),
				client.WithRetry(client.RetryConfig{MaxAttempts: 2, BaseDelay: time.Nanosecond, MaxDelay: time.Nanosecond}),
			)
			if err != nil {
				t.Fatalf("client.New() error = %v", err)
			}
			t.Cleanup(func() { _ = cl.Close() })
			cl.SetToken(&client.Token{Value: secondMaterial, Status: client.TokenStatusNormal})

			parentCtx, parent := v1RemoteParent(t)
			ctx := client.WithCorrelationID(parentCtx, correlation)
			if err := tc.call(ctx, cl); err != nil {
				t.Fatalf("%s() error = %v", tc.name, err)
			}

			spans := recorder.Ended()
			if len(spans) != 2 {
				t.Fatalf("ended spans = %d, want two attempts", len(spans))
			}
			for i, span := range spans {
				wantAttempt := int64(i + 1)
				if span.Name() != "GET "+route || span.SpanKind() != trace.SpanKindClient {
					t.Fatalf("span %d = %q/%v, want GET route/client", i, span.Name(), span.SpanKind())
				}
				if span.Parent().TraceID() != parent.TraceID() || span.Parent().SpanID() != parent.SpanID() {
					t.Fatalf("span %d parent = %s/%s, want application parent", i, span.Parent().TraceID(), span.Parent().SpanID())
				}
				if span.EndTime().IsZero() {
					t.Fatalf("span %d did not end", i)
				}
				v1AssertIntAttr(t, span, "webull.attempt", wantAttempt)
				v1AssertStringAttr(t, span, "http.route", route)
				v1AssertStringAttr(t, span, "webull.correlation_id", correlation)
				if i == 0 {
					v1AssertIntAttr(t, span, "http.status_code", http.StatusInternalServerError)
					v1AssertStringAttr(t, span, "error", "webull: SERVER_ERROR")
					if span.Status().Code != otelcodes.Error {
						t.Fatalf("failed span status = %v, want ERROR", span.Status())
					}
				} else {
					v1AssertIntAttr(t, span, "http.status_code", http.StatusOK)
					if _, ok := v1SpanAttrs(span)["error"]; ok {
						t.Fatal("successful retry span contains an error attribute")
					}
					if span.Status().Code != otelcodes.Ok {
						t.Fatalf("successful retry span status = %v, want OK", span.Status())
					}
				}
			}

			records := handler.snapshot()
			if len(records) != 4 {
				t.Fatalf("log record count = %d, want two start/done pairs", len(records))
			}
			wantAttempts := []int{1, 1, 2, 2}
			wantLevels := []slog.Level{slog.LevelInfo, slog.LevelError, slog.LevelInfo, slog.LevelInfo}
			wantMessages := []string{"webull request start", "webull request done", "webull request start", "webull request done"}
			wantStatuses := []int{0, http.StatusInternalServerError, 0, http.StatusOK}
			for i, record := range records {
				if record.level != wantLevels[i] || record.message != wantMessages[i] {
					t.Fatalf("log %d = %+v, want level/message %v/%q", i, record, wantLevels[i], wantMessages[i])
				}
				if i%2 == 0 {
					v1AssertSlogAttrs(t, record, map[string]any{
						"method":         http.MethodGet,
						"path":           route,
						"attempt":        wantAttempts[i],
						"correlation_id": correlation,
					})
				} else {
					want := map[string]any{
						"method":         http.MethodGet,
						"path":           route,
						"attempt":        wantAttempts[i],
						"status":         wantStatuses[i],
						"correlation_id": correlation,
						"latency":        time.Duration(0),
					}
					if i == 1 {
						want["error"] = "webull: SERVER_ERROR"
					}
					v1AssertSlogAttrs(t, record, want)
				}
			}

			signaturesMu.Lock()
			gotSignatures := append([]string(nil), signatures...)
			signaturesMu.Unlock()
			if len(gotSignatures) != 2 || gotSignatures[0] == "" {
				t.Fatalf("captured signatures = %v, want two non-empty values", gotSignatures)
			}
			for _, span := range spans {
				text := v1SpanText(span)
				for _, secret := range []string{appKey, firstMaterial, secondMaterial, thirdMaterial, gotSignatures[0], gotSignatures[1]} {
					if strings.Contains(text, secret) {
						t.Fatalf("span telemetry contains sensitive value %q: %s", secret, text)
					}
				}
			}
			logText := v1SlogText(records)
			for _, secret := range []string{appKey, firstMaterial, secondMaterial, thirdMaterial, gotSignatures[0], gotSignatures[1]} {
				if strings.Contains(logText, secret) {
					t.Fatalf("log telemetry contains sensitive value %q: %s", secret, logText)
				}
			}
		})
	}
}

func TestV1RESTNilLoggerIsSafe(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(srv.Close)
	cl, err := client.New(
		client.WithCredentials("nil-logger-key", "nil-logger-secret"),
		client.WithBaseURL(srv.URL),
		client.WithLogger(nil),
		client.WithoutRetry(),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })
	if err := cl.Do(context.Background(), http.MethodGet, "/nil-logger", nil, nil); err != nil {
		t.Fatalf("Do() with nil logger error = %v", err)
	}
}
