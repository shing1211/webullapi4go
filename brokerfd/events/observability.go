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
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	grpcstatus "google.golang.org/grpc/status"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/pkg/observability"
)

const (
	brokerEventGRPCService            = "grpc.event.EventService"
	brokerEventGRPCMethod             = "Subscribe"
	brokerEventGRPCPath               = "/grpc.event.EventService/Subscribe"
	brokerEventCorrelationMetadataKey = "x-correlation-id"
)

type eventMetrics struct {
	attempts metric.Int64Counter
	duration metric.Float64Histogram
}

func newEventMetrics(cfg *observability.Config) *eventMetrics {
	if cfg == nil {
		return nil
	}
	meter := cfg.Meter("webullapi4go/brokerfd/events")
	if meter == nil {
		return nil
	}
	attempts, _ := meter.Int64Counter(
		"event_stream_attempts",
		metric.WithDescription("gRPC event stream attempts"),
		metric.WithUnit("{attempt}"),
	)
	duration, _ := meter.Float64Histogram(
		"event_stream_attempt_duration",
		metric.WithDescription("Duration of gRPC event stream attempts"),
		metric.WithUnit("ms"),
	)
	if attempts == nil && duration == nil {
		return nil
	}
	return &eventMetrics{attempts: attempts, duration: duration}
}

type eventAttempt struct {
	ctx           context.Context
	span          trace.Span
	cfg           *observability.Config
	metrics       *eventMetrics
	started       time.Time
	attempt       int
	correlationID string
}

func (c *Client) observabilityConfig() *observability.Config {
	if c == nil || c.core == nil {
		return nil
	}
	return c.core.ObservabilityConfig()
}

func (c *Client) startEventAttempt(ctx context.Context, attempt int) (context.Context, *eventAttempt) {
	if attempt < 1 {
		attempt = 1
	}
	ctx, correlationID := ensureCorrelationID(ctx)
	state := &eventAttempt{
		ctx:           ctx,
		cfg:           c.observabilityConfig(),
		metrics:       c.metrics,
		started:       time.Now(),
		attempt:       attempt,
		correlationID: correlationID,
	}

	if state.cfg != nil {
		if tracer := state.cfg.Tracer("webullapi4go/brokerfd/events"); tracer != nil {
			ctx, state.span = tracer.Start(ctx, brokerEventGRPCPath,
				trace.WithSpanKind(trace.SpanKindClient),
				trace.WithAttributes(
					attribute.String("rpc.system", "grpc"),
					attribute.String("rpc.service", brokerEventGRPCService),
					attribute.String("rpc.method", brokerEventGRPCMethod),
					attribute.Int("webull.attempt", attempt),
					attribute.String("webull.correlation_id", correlationID),
				),
			)
			state.ctx = ctx
		}
		if log := state.cfg.Logger; log != nil {
			log.LogAttrs(state.ctx, slog.LevelInfo, "webull brokerfd events stream start",
				slog.String("transport", "grpc"),
				slog.String("service", brokerEventGRPCService),
				slog.String("method", brokerEventGRPCMethod),
				slog.Int("attempt", attempt),
				slog.String("correlation_id", correlationID),
			)
		}
	}

	if state.metrics != nil && state.metrics.attempts != nil {
		state.metrics.attempts.Add(state.ctx, 1, metric.WithAttributes(
			attribute.String("rpc.system", "grpc"),
			attribute.String("rpc.service", brokerEventGRPCService),
			attribute.String("rpc.method", brokerEventGRPCMethod),
		))
	}
	return ctx, state
}

func (a *eventAttempt) finish(err error) {
	if a == nil {
		return
	}
	statusCode := eventStatusCode(a.ctx, err)
	outcome := "ok"
	level := slog.LevelInfo
	if err != nil {
		outcome = "error"
		level = slog.LevelError
	}

	if a.span != nil {
		a.span.SetAttributes(
			attribute.Int("rpc.grpc.status_code", int(statusCode)),
			attribute.String("webull.stream.outcome", outcome),
		)
		if err != nil {
			a.span.RecordError(err)
			a.span.SetStatus(otelcodes.Error, "gRPC event stream attempt failed")
		} else {
			a.span.SetStatus(otelcodes.Ok, "")
		}
	}

	latency := time.Since(a.started)
	if a.span != nil {
		a.span.SetAttributes(attribute.Int64("webull.stream.duration_ms", latency.Milliseconds()))
	}
	if a.metrics != nil && a.metrics.duration != nil {
		a.metrics.duration.Record(a.ctx, float64(latency.Microseconds())/1000, metric.WithAttributes(
			attribute.String("rpc.system", "grpc"),
			attribute.String("rpc.service", brokerEventGRPCService),
			attribute.String("rpc.method", brokerEventGRPCMethod),
			attribute.Int("rpc.grpc.status_code", int(statusCode)),
			attribute.String("webull.stream.outcome", outcome),
		))
	}

	if a.cfg != nil && a.cfg.Logger != nil {
		attrs := []slog.Attr{
			slog.String("transport", "grpc"),
			slog.String("service", brokerEventGRPCService),
			slog.String("method", brokerEventGRPCMethod),
			slog.Int("attempt", a.attempt),
			slog.Int("status_code", int(statusCode)),
			slog.String("status", statusCode.String()),
			slog.String("correlation_id", a.correlationID),
			slog.Duration("latency", latency),
		}
		if err != nil {
			attrs = append(attrs, slog.String("error", "gRPC event stream attempt failed"))
		}
		a.cfg.Logger.LogAttrs(a.ctx, level, "webull brokerfd events stream done", attrs...)
	}
	if a.span != nil {
		a.span.End()
	}
}

func eventStatusCode(ctx context.Context, err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if errors.Is(err, context.Canceled) || (ctx != nil && errors.Is(ctx.Err(), context.Canceled)) {
		return codes.Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) || (ctx != nil && errors.Is(ctx.Err(), context.DeadlineExceeded)) {
		return codes.DeadlineExceeded
	}
	return grpcstatus.Code(err)
}

func (c *Client) addEventMetadata(ctx context.Context, md metadata.MD) metadata.MD {
	if md == nil {
		md = metadata.MD{}
	}
	if cfg := c.observabilityConfig(); cfg != nil {
		carrier := propagation.MapCarrier{}
		cfg.InjectTraceContext(ctx, carrier)
		for key, value := range carrier {
			md.Set(key, value)
		}
	}
	if id := client.CorrelationIDFromContext(ctx); id != "" {
		md.Set(brokerEventCorrelationMetadataKey, id)
	}
	return md
}

func ensureCorrelationID(ctx context.Context) (context.Context, string) {
	if id := client.CorrelationIDFromContext(ctx); id != "" {
		return ctx, id
	}
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		id := "0000000000000000"
		return client.WithCorrelationID(ctx, id), id
	}
	id := hex.EncodeToString(b[:])
	return client.WithCorrelationID(ctx, id), id
}
