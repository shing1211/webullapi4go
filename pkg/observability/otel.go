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

// Package observability provides OpenTelemetry tracing and metrics integration for
// the Webull API SDK. All tracer and meter types are API-only: the no-op
// implementation is used by default, so the package compiles and works without
// any telemetry backend. Supply a real [TracerProvider] or [MeterProvider] via
// [WithTracerProvider] or [WithMeterProvider] to enable telemetry.
package observability

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// TracerProvider is the interface for supplying a tracer. It is satisfied by
// otel.TracerProvider.
type TracerProvider = trace.TracerProvider

// MeterProvider is the interface for supplying a meter. It is satisfied by
// otel.MeterProvider.
type MeterProvider = metric.MeterProvider

// TextMapPropagator is the interface for trace context propagation.
// It is satisfied by otel's propagation.TextMapPropagator.
type TextMapPropagator = propagation.TextMapPropagator

// Tracer is the interface for OpenTelemetry tracing. It is satisfied by
// otel.Tracer.
type Tracer = trace.Tracer

// Span is the interface for an in-flight trace. It is satisfied by
// otel.Span.
type Span = trace.Span

// Config holds the observability handles used by the SDK. All fields are
// optional and default to no-op implementations.
type Config struct {
	TracerProvider TracerProvider
	MeterProvider  MeterProvider
	Logger         *slog.Logger
	Propagator     TextMapPropagator
}

// DefaultConfig returns a config with no-op tracer and meter providers and a
// no-op logger.
func DefaultConfig() Config {
	return Config{
		TracerProvider: otel.GetTracerProvider(),
		MeterProvider:  otel.GetMeterProvider(),
		Propagator:     otel.GetTextMapPropagator(),
	}
}

// Tracer returns a tracer for the given name, scoped to the configured
// TracerProvider. If no TracerProvider is configured, the global no-op is used.
func (c Config) Tracer(name string, opts ...trace.TracerOption) Tracer {
	if c.TracerProvider == nil {
		return otel.Tracer(name, opts...)
	}
	return c.TracerProvider.Tracer(name, opts...)
}

// Meter returns a meter for the given name, scoped to the configured
// MeterProvider. If no MeterProvider is configured, the global no-op is used.
func (c Config) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	if c.MeterProvider == nil {
		return otel.Meter(name, opts...)
	}
	return c.MeterProvider.Meter(name, opts...)
}

// SpanAttributes returns the standard set of span attributes used by the SDK.
func SpanAttributes(method, path string, attempt int) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("http.method", method),
		attribute.String("http.route", path),
		attribute.Int("webull.attempt", attempt),
	}
}

// SpanName returns a span name from an HTTP method and path.
func SpanName(method, path string) string {
	return method + " " + path
}

// InjectTraceContext propagates trace context from ctx into carrier using
// cfg.Propagator. It is a no-op when no propagator is configured.
func (c Config) InjectTraceContext(ctx context.Context, carrier propagation.TextMapCarrier) {
	if c.Propagator != nil {
		c.Propagator.Inject(ctx, carrier)
	}
}

// ExtractTraceContext extracts trace context from carrier into ctx using
// cfg.Propagator. It returns ctx unchanged when no propagator is configured.
func (c Config) ExtractTraceContext(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	if c.Propagator != nil {
		return c.Propagator.Extract(ctx, carrier)
	}
	return ctx
}
