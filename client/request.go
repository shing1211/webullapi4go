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

package client

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/shing1211/webullapi4go/internal/auth"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/pkg/observability"
	"github.com/shing1211/webullapi4go/pkg/resilience/retry"
	"github.com/shing1211/webullapi4go/pkg/transport"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Headers that the SDK sets on every request but that do not participate in the
// canonical signature.
const (
	headerSignature     = "x-signature"
	headerVersion       = "x-version"
	headerCorrelationID = "x-correlation-id"
)

// contextKey is a value of unique type used as a context key.
type contextKey struct{}

// CorrelationIDKey is the context key for the per-request correlation
// identifier. It is used by [WithCorrelationID] and [CorrelationIDFromContext].
var CorrelationIDKey = contextKey{}

// WithCorrelationID stores a correlation identifier in the context. The SDK
// generates one automatically if not already present; callers may set it before
// [Client.Do] to propagate an externally-generated ID.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, id)
}

// CorrelationIDFromContext returns the correlation identifier stored in ctx, or
// an empty string if none is set.
func CorrelationIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(CorrelationIDKey).(string); ok {
		return v
	}
	return ""
}

// newCorrelationID generates a fresh 16-character hex correlation identifier
// using crypto/rand.
func newCorrelationID() string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "0000000000000000"
	}
	return hex.EncodeToString(b[:])
}

func ensureCorrelationID(ctx context.Context) (context.Context, string) {
	if id := CorrelationIDFromContext(ctx); id != "" {
		return ctx, id
	}
	id := newCorrelationID()
	return WithCorrelationID(ctx, id), id
}

// ErrCircuitOpen is wrapped into the error returned by [Client.Do],
// [Client.DoBroker], and [Client.DoStream] when a configured circuit breaker
// rejects a call.
var ErrCircuitOpen = errs.New(errs.CodeTransport, "circuit breaker is open")

// Do performs a signed request against the Webull OpenAPI and decodes the JSON
// response into out. It is the single transport entry point for the SDK.
//
// path is resolved against the configured base URL and may carry a query
// string, for example "/openapi/account/list?page=1". body is optional: when
// non-nil it is serialized as compact JSON without HTML escaping, and the exact
// bytes transmitted are the bytes that are signed. A nil body sends no request
// body. out is optional: when non-nil the response body is decoded into it with
// [encoding/json].
//
// Do applies the configured resilience primitives without changing this
// signature: the rate limiter waits before every attempt, the circuit breaker
// may reject an attempt, and transient errors are retried according to the
// configured [RetryConfig]. By default only idempotent requests (GET and HEAD)
// are retried.
//
// On a 2xx response Do returns nil. On any other status it returns a typed
// *[errs.Error] classified by [errs.FromHTTPStatus]; network and encoding
// failures return [errs.CodeTransport] or [errs.CodeAPI] errors respectively.
func (c *Client) Do(ctx context.Context, method, path string, body, out any) error {
	reqPath, query, err := splitPathQuery(path)
	if err != nil {
		return err
	}
	query = normalizeQuery(query)
	ctx, _ = ensureCorrelationID(ctx)

	if err := c.ensureAutoToken(ctx, reqPath); err != nil {
		return err
	}

	bodyBytes, err := marshalRequestBody(body)
	if err != nil {
		return err
	}

	retryable := c.cfg.retry != nil && (c.cfg.retryNonIdempotent || isIdempotent(method))
	if !retryable {
		return c.attempt(ctx, method, reqPath, query, bodyBytes, out)
	}
	return c.runWithRetry(ctx, func(attemptCtx context.Context, attempt int) error {
		_, err := c.attemptWithNumber(attemptCtx, method, reqPath, query, bodyBytes, out, attempt, c.transport)
		return err
	})
}

// DoBroker performs a signed request against the Broker API and decodes the JSON
// response into out. It is identical to [Client.Do] except that it uses the
// Broker HTTP endpoint (cfg.Endpoints.BrokerHTTP) instead of the default HTTP
// endpoint. If BrokerHTTP is not configured, DoBroker returns an error.
func (c *Client) DoBroker(ctx context.Context, method, path string, body, out any) error {
	if c.brokerTr == nil {
		return errs.New(errs.CodeInvalidConfig, "broker HTTP endpoint is not configured")
	}
	reqPath, query, err := splitPathQuery(path)
	if err != nil {
		return err
	}
	query = normalizeQuery(query)
	ctx, _ = ensureCorrelationID(ctx)

	if err := c.ensureAutoToken(ctx, reqPath); err != nil {
		return err
	}

	bodyBytes, err := marshalRequestBody(body)
	if err != nil {
		return err
	}

	retryable := c.cfg.retry != nil && (c.cfg.retryNonIdempotent || isIdempotent(method))
	if !retryable {
		return c.attemptBroker(ctx, method, reqPath, query, bodyBytes, out)
	}
	return c.runWithRetry(ctx, func(attemptCtx context.Context, attempt int) error {
		_, err := c.attemptBrokerWithNumber(attemptCtx, method, reqPath, query, bodyBytes, out, attempt)
		return err
	})
}

func (c *Client) runWithRetry(ctx context.Context, fn func(context.Context, int) error) error {
	attempt := 0
	return c.cfg.retry.Do(ctx, func(attemptCtx context.Context) error {
		attempt++
		return fn(attemptCtx, attempt)
	})
}

type attemptResult struct {
	status   int
	response *http.Response
}

// runAttempt executes one common request attempt. Rate limiting and circuit
// breaking run before user interceptors, and every HTTP response is recorded by
// the same hook, span, logger, and metric path.
func (c *Client) runAttempt(ctx context.Context, method, reqPath string, attempt int, requireResponse bool, operation func(context.Context) (attemptResult, error)) (attemptResult, error) {
	ctx, correlationID := ensureCorrelationID(ctx)
	if attempt < 1 {
		attempt = 1
	}
	start := time.Now()

	tracer := c.cfg.otel.Tracer("webullapi4go/client")
	ctx, span := tracer.Start(ctx, observability.SpanName(method, reqPath),
		trace.WithAttributes(observability.SpanAttributes(method, reqPath, attempt)...),
		trace.WithAttributes(attribute.String("webull.correlation_id", correlationID)),
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	if log := c.cfg.otel.Logger; log != nil {
		log.LogAttrs(ctx, slog.LevelInfo, "webull request start",
			slog.String("method", method),
			slog.String("path", reqPath),
			slog.Int("attempt", attempt),
			slog.String("correlation_id", correlationID),
		)
	}

	var result attemptResult
	var err error
	if c.cfg.rateLimiter != nil {
		err = c.cfg.rateLimiter.Wait(ctx, reqPath)
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				err = ctxErr
			} else {
				err = errs.Wrap(errs.CodeRateLimited, "rate limiter", err)
			}
		}
	}
	if err == nil && c.cfg.breaker != nil && !c.cfg.breaker.Allow() {
		err = retry.Permanent(errs.Wrap(errs.CodeTransport, "circuit breaker open", ErrCircuitOpen))
	}
	if err == nil {
		if h := c.cfg.hooks; h.OnRequest != nil {
			h.OnRequest(method, reqPath)
		}

		next := func(operationCtx context.Context) error {
			operationCtx = WithCorrelationID(operationCtx, correlationID)
			if !trace.SpanFromContext(operationCtx).SpanContext().IsValid() {
				operationCtx = trace.ContextWithSpan(operationCtx, span)
			}
			var operationErr error
			result, operationErr = operation(operationCtx)
			c.recordBreaker(operationErr)
			return operationErr
		}
		for i := len(c.cfg.interceptors) - 1; i >= 0; i-- {
			interceptor := c.cfg.interceptors[i]
			nextFn := next
			next = func(operationCtx context.Context) error {
				return interceptor(operationCtx, nextFn)
			}
		}
		err = next(ctx)
		if err == nil && requireResponse && result.response == nil {
			err = errs.New(errs.CodeTransport, "stream interceptor returned no response")
		}
	}

	latency := time.Since(start)
	if h := c.cfg.hooks; h.OnLatency != nil {
		h.OnLatency(attempt, latency)
	}
	if err == nil {
		if h := c.cfg.hooks; h.OnResponse != nil {
			h.OnResponse(result.status, latency)
		}
	} else if h := c.cfg.hooks; h.OnError != nil {
		h.OnError(err, latency)
	}

	span.SetAttributes(attribute.Int64("http.duration_ms", latency.Milliseconds()))
	if result.status != 0 {
		span.SetAttributes(attribute.Int("http.status_code", result.status))
	}
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
	}

	if hist := c.cfg.otel.ClientLatencyHistogram(); hist != nil {
		hist.Record(ctx, float64(latency.Milliseconds()), metric.WithAttributes(attribute.String("http.route", reqPath)))
	}

	if log := c.cfg.otel.Logger; log != nil {
		attrs := []slog.Attr{
			slog.String("method", method),
			slog.String("path", reqPath),
			slog.Int("attempt", attempt),
			slog.Int("status", result.status),
			slog.String("correlation_id", correlationID),
			slog.Duration("latency", latency),
		}
		if err != nil {
			attrs = append(attrs, slog.String("error", err.Error()))
			log.LogAttrs(ctx, slog.LevelError, "webull request done", attrs...)
		} else {
			log.LogAttrs(ctx, slog.LevelInfo, "webull request done", attrs...)
		}
	}

	return result, err
}

func (c *Client) attempt(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte, out any) error {
	_, err := c.attemptWithNumber(ctx, method, reqPath, query, bodyBytes, out, 1, c.transport)
	return err
}

func (c *Client) attemptWithNumber(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte, out any, attempt int, tr *transport.Client) (attemptResult, error) {
	return c.runAttempt(ctx, method, reqPath, attempt, false, func(operationCtx context.Context) (attemptResult, error) {
		return c.executeBuffered(operationCtx, method, reqPath, query, bodyBytes, tr, out)
	})
}

func (c *Client) attemptBroker(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte, out any) error {
	_, err := c.attemptBrokerWithNumber(ctx, method, reqPath, query, bodyBytes, out, 1)
	return err
}

func (c *Client) attemptBrokerWithNumber(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte, out any, attempt int) (attemptResult, error) {
	return c.attemptWithNumber(ctx, method, reqPath, query, bodyBytes, out, attempt, c.brokerTr)
}

func (c *Client) executeBuffered(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte, tr *transport.Client, out any) (attemptResult, error) {
	var result attemptResult
	req, err := c.buildSignedRequestForTransport(ctx, method, reqPath, query, bodyBytes, tr)
	if err != nil {
		return result, err
	}

	resp, err := tr.Do(req)
	if resp != nil {
		result.status = resp.StatusCode
	}
	if err != nil {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		return result, errs.Wrap(errs.CodeTransport, method+" "+reqPath, err)
	}
	if resp == nil {
		return result, errs.New(errs.CodeTransport, method+" "+reqPath+" returned nil response")
	}

	if resp.Body == nil {
		return result, errs.New(errs.CodeTransport, method+" "+reqPath+" returned nil response body")
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, errs.Wrap(errs.CodeTransport, "reading response body", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return result, errs.FromHTTPStatus(resp.StatusCode, data)
	}
	c.maybeUpdateClockOffset(resp.Header)
	if out == nil || len(data) == 0 {
		return result, nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return result, errs.Wrap(errs.CodeAPI, "decoding response body", err)
	}
	return result, nil
}

// recordBreaker reports the outcome of an attempt to the circuit breaker.
// Server and transport failures trip the breaker; every other error means the
// service responded, so it counts as a success for breaker purposes.
func (c *Client) recordBreaker(err error) {
	if c.cfg.breaker == nil {
		return
	}
	if err != nil && (errs.Is(err, errs.CodeServer) || errs.Is(err, errs.CodeTransport)) {
		c.cfg.breaker.RecordFailure()
		return
	}
	c.cfg.breaker.RecordSuccess()
}

// isIdempotent reports whether method may be retried by default. Only the safe,
// idempotent methods are retried; callers can opt in to others with
// [RetryConfig.RetryNonIdempotent].
func isIdempotent(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead:
		return true
	default:
		return false
	}
}

// buildSignedRequest constructs the outgoing request for one API call and
// attaches its signing headers. It is shared by [Client.Do],
// [Client.DoBroker], and [Client.DoStream] so that a streamed request is built
// and signed exactly like a buffered one: the x-version header follows the same
// per-path defaults, the host that is signed is the host that is sent, and the
// exact bytes that are signed are the bytes that are transmitted.
func (c *Client) buildSignedRequest(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte) (*http.Request, error) {
	return c.buildSignedRequestForTransport(ctx, method, reqPath, query, bodyBytes, c.transport)
}

func (c *Client) buildSignedRequestForTransport(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte, tr *transport.Client) (*http.Request, error) {
	if tr == nil {
		return nil, errs.New(errs.CodeInvalidConfig, "HTTP transport is not configured")
	}
	req, err := tr.NewRequest(ctx, method, reqPath, query, bodyBytes)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInvalidConfig, "building request", err)
	}

	host := signingHost(req.URL)
	req.Host = host

	signingHeaders, err := auth.NewSigningHeaders(c.cfg.AppKey, host, c.signingTime())
	if err != nil {
		return nil, errs.Wrap(errs.CodeAuth, "building signing headers", err)
	}
	applySigningHeaders(req.Header, signingHeaders)
	req.Header.Set(headerVersion, c.apiVersionFor(reqPath))
	if len(bodyBytes) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	signature, err := auth.Sign(auth.SignParams{
		Method:    method,
		Path:      reqPath,
		Query:     query,
		Headers:   signingHeaders,
		Body:      bodyBytes,
		AppSecret: c.cfg.AppSecret,
	})
	if err != nil {
		return nil, errs.Wrap(errs.CodeAuth, "signing request", err)
	}
	req.Header.Set(headerSignature, signature)
	c.cfg.otel.InjectTraceContext(ctx, propagation.HeaderCarrier(req.Header))

	if cid := CorrelationIDFromContext(ctx); cid != "" {
		req.Header.Set(headerCorrelationID, cid)
	}

	return req, nil
}

// splitPathQuery separates an optional query string from path and parses it.
func splitPathQuery(path string) (string, url.Values, error) {
	i := strings.IndexByte(path, '?')
	if i < 0 {
		return path, nil, nil
	}
	query, err := url.ParseQuery(path[i+1:])
	if err != nil {
		return "", nil, errs.Wrap(errs.CodeInvalidConfig, "invalid query string in path", err)
	}
	return path[:i], query, nil
}

// normalizeQuery sorts the values of each key so that the query string that is
// sent is byte-for-byte the query string that is signed. A nil or empty result
// is returned as nil.
func normalizeQuery(q url.Values) url.Values {
	if len(q) == 0 {
		return nil
	}
	for key := range q {
		sort.Strings(q[key])
	}
	return q
}

// marshalRequestBody encodes body as compact JSON without HTML escaping. A
// []byte (including json.RawMessage) is sent verbatim so callers can supply
// pre-encoded payloads.
func marshalRequestBody(body any) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	switch raw := body.(type) {
	case []byte:
		return raw, nil
	case json.RawMessage:
		return raw, nil
	}
	data, err := auth.MarshalBody(body)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInvalidConfig, "encoding request body", err)
	}
	return data, nil
}

// signingHost returns the value that must be signed as the "host" header: the
// request host as "hostname[:port]", without the scheme. The port is included
// only when it is not the default for the scheme (443 for https, 80 for http),
// matching Webull's documented example of "api.webull.com".
func signingHost(u *url.URL) string {
	host := u.Hostname()
	port := u.Port()
	if port == "" {
		return host
	}
	if (u.Scheme == "https" && port == "443") || (u.Scheme == "http" && port == "80") {
		return host
	}
	return net.JoinHostPort(host, port)
}

// applySigningHeaders copies the signing headers into dst, skipping "host": the
// Host header is sent from [http.Request.Host] by the HTTP stack, and copying
// it here would be redundant.
func applySigningHeaders(dst, src http.Header) {
	for name, values := range src {
		if strings.EqualFold(name, auth.HeaderHost) {
			continue
		}
		for _, v := range values {
			dst.Add(name, v)
		}
	}
}

// signingTime returns the time used in the x-timestamp signing header. When
// clock-drift correction is enabled, it includes the learned offset between the
// local clock and the server's clock.
func (c *Client) signingTime() time.Time {
	if !c.cfg.clockDriftCorrection {
		return time.Now()
	}
	c.clockOffsetMu.Lock()
	offset := c.clockOffset
	c.clockOffsetMu.Unlock()
	return time.Now().Add(offset)
}

// maybeUpdateClockOffset parses the Date header from a successful response and
// updates the learned clock offset, clamped to the range [-5 minutes, +5
// minutes]. It is a no-op when clock-drift correction is disabled.
func (c *Client) maybeUpdateClockOffset(respHeaders http.Header) {
	if !c.cfg.clockDriftCorrection {
		return
	}
	dateStr := respHeaders.Get("Date")
	if dateStr == "" {
		return
	}
	serverTime, err := parseDateHeader(dateStr)
	if err != nil {
		return
	}
	offset := time.Until(serverTime)
	clamped := offset
	const maxOffset = 5 * time.Minute
	if offset > maxOffset {
		clamped = maxOffset
	} else if offset < -maxOffset {
		clamped = -maxOffset
	}
	c.clockOffsetMu.Lock()
	c.clockOffset = clamped
	c.clockOffsetMu.Unlock()
}

// parseDateHeader parses an RFC 7231 Date header value (e.g.
// "Tue, 23 Sep 2025 12:34:56 GMT") into a time.Time in UTC.
func parseDateHeader(s string) (time.Time, error) {
	return http.ParseTime(s)
}
