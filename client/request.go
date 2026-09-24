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
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/shing1211/webullapi4go/internal/auth"
	"github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/pkg/observability"
	"github.com/shing1211/webullapi4go/pkg/resilience/retry"
	"github.com/shing1211/webullapi4go/pkg/transport"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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
	if v := ctx.Value(CorrelationIDKey); v != nil {
		return v.(string)
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

// ErrCircuitOpen is wrapped into the error returned by [Client.Do] when a
// configured circuit breaker rejects a call.
var ErrCircuitOpen = errors.New("webull: circuit breaker is open")

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

	// Obtain an access token before the first token-consuming request when
	// [WithAutoToken] is enabled. Token-lifecycle endpoints are exempt, so this
	// is a no-op for the EnsureToken flow itself.
	if err := c.ensureAutoToken(ctx, reqPath); err != nil {
		return err
	}

	bodyBytes, err := marshalRequestBody(body)
	if err != nil {
		return err
	}

	attempt := func(ctx context.Context) error {
		return c.attempt(ctx, method, reqPath, query, bodyBytes, out)
	}
	if c.cfg.retry != nil && (c.cfg.retryNonIdempotent || isIdempotent(method)) {
		return c.cfg.retry.Do(ctx, attempt)
	}
	return attempt(ctx)
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

	if err := c.ensureAutoToken(ctx, reqPath); err != nil {
		return err
	}

	bodyBytes, err := marshalRequestBody(body)
	if err != nil {
		return err
	}

	attempt := func(ctx context.Context) error {
		return c.attemptBroker(ctx, method, reqPath, query, bodyBytes, out)
	}
	if c.cfg.retry != nil && (c.cfg.retryNonIdempotent || isIdempotent(method)) {
		return c.cfg.retry.Do(ctx, attempt)
	}
	return attempt(ctx)
}

func (c *Client) attemptBroker(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte, out any) error {
	correlationID := CorrelationIDFromContext(ctx)
	if correlationID == "" {
		correlationID = newCorrelationID()
		ctx = WithCorrelationID(ctx, correlationID)
	}

	tracer := c.cfg.otel.Tracer("webullapi4go/client")
	ctx, span := tracer.Start(ctx, observability.SpanName(method, reqPath),
		trace.WithAttributes(observability.SpanAttributes(method, reqPath, 0)...),
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	if c.cfg.rateLimiter != nil {
		if err := c.cfg.rateLimiter.Wait(ctx, reqPath); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return ctxErr
			}
			return errs.Wrap(errs.CodeRateLimited, "rate limiter", err)
		}
	}
	if c.cfg.breaker != nil && !c.cfg.breaker.Allow() {
		return retry.Permanent(errs.Wrap(errs.CodeTransport, "circuit breaker open", ErrCircuitOpen))
	}

	err := c.executeBroker(ctx, method, reqPath, query, bodyBytes, out)
	if c.cfg.breaker != nil {
		c.recordBreaker(err)
	}
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
	}
	return err
}

func (c *Client) executeBroker(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte, out any) error {
	req, err := c.buildSignedRequestForTransport(ctx, method, reqPath, query, bodyBytes, c.brokerTr)
	if err != nil {
		return err
	}

	resp, err := c.brokerTr.Do(req)
	if err != nil {
		return errs.Wrap(errs.CodeTransport, method+" "+reqPath, err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errs.Wrap(errs.CodeTransport, "reading response body", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return errs.FromHTTPStatus(resp.StatusCode, data)
	}
	c.maybeUpdateClockOffset(resp.Header)
	if out == nil || len(data) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return errs.Wrap(errs.CodeAPI, "decoding response body", err)
	}
	return nil
}

func (c *Client) buildSignedRequestForTransport(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte, tr *transport.Client) (*http.Request, error) {
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

	if cid := CorrelationIDFromContext(ctx); cid != "" {
		req.Header.Set(headerCorrelationID, cid)
	}

	return req, nil
}

// attempt executes a single request attempt through the ordered interceptor
// chain. The fixed order is:
//
//	rate-limiter.Wait → circuit-breaker.Allow → user-interceptors → sign → send
//
// Each user-supplied interceptor receives the next function in the chain and
// may inspect, decorate, or short-circuit the request. After the response,
// the circuit-breaker outcome is recorded so that retries participate in
// breaker accounting.
//
// attempt establishes an OpenTelemetry span for the request and records the
// round-trip latency and correlation ID to the configured logger.
func (c *Client) attempt(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte, out any) error {
	start := time.Now()

	// Extract or generate a correlation ID for this request.
	correlationID := CorrelationIDFromContext(ctx)
	if correlationID == "" {
		correlationID = newCorrelationID()
		ctx = WithCorrelationID(ctx, correlationID)
	}

	// Create an OTel span for this attempt.
	tracer := c.cfg.otel.Tracer("webullapi4go/client")
	ctx, span := tracer.Start(ctx, observability.SpanName(method, reqPath),
		trace.WithAttributes(observability.SpanAttributes(method, reqPath, 0)...),
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	// Log request start if a logger is configured.
	if log := c.cfg.otel.Logger; log != nil {
		log.LogAttrs(ctx, slog.LevelInfo, "webull request start",
			slog.String("method", method),
			slog.String("path", reqPath),
			slog.String("correlation_id", correlationID),
		)
	}

	// The core function that performs rate-limiting, circuit-breaking, and the
	// HTTP call (sign + send).
	core := func(ctx context.Context) error {
		if c.cfg.rateLimiter != nil {
			if err := c.cfg.rateLimiter.Wait(ctx, reqPath); err != nil {
				if ctxErr := ctx.Err(); ctxErr != nil {
					return ctxErr
				}
				return errs.Wrap(errs.CodeRateLimited, "rate limiter", err)
			}
		}
		if c.cfg.breaker != nil && !c.cfg.breaker.Allow() {
			return retry.Permanent(errs.Wrap(errs.CodeTransport, "circuit breaker open", ErrCircuitOpen))
		}
		err := c.execute(ctx, method, reqPath, query, bodyBytes, out)
		c.recordBreaker(err)
		return err
	}

	// Build the interceptor chain from the inside out: the last supplied
	// interceptor is closest to the core HTTP call.
	// Chain: interceptors[0] → interceptors[1] → ... → interceptors[n] → core
	next := core
	for i := len(c.cfg.interceptors) - 1; i >= 0; i-- {
		fn := c.cfg.interceptors[i]
		n := next
		next = func(ctx context.Context) error { return fn(ctx, n) }
	}

	// Fire the OnRequest hook before the chain runs.
	if h := c.cfg.hooks; h.OnRequest != nil {
		h.OnRequest(method, reqPath)
	}

	// Wrap the whole chain for total latency tracking and hook dispatch.
	var err error
	var latency time.Duration
	func() {
		t0 := time.Now()
		err = next(ctx)
		latency = time.Since(t0)
		if h := c.cfg.hooks; h.OnLatency != nil {
			h.OnLatency(0, latency)
		}
	}()

	if err == nil {
		if h := c.cfg.hooks; h.OnResponse != nil {
			h.OnResponse(0, latency)
		}
	} else {
		if h := c.cfg.hooks; h.OnError != nil {
			h.OnError(err, latency)
		}
	}

	span.SetAttributes(attribute.Int64("http.duration_ms", latency.Milliseconds()))
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
	}

	if hist := c.cfg.otel.ClientLatencyHistogram(); hist != nil {
		hist.Record(ctx, float64(latency.Milliseconds()),
			metric.WithAttributes(attribute.String("http.route", reqPath)))
	}

	if log := c.cfg.otel.Logger; log != nil {
		if err != nil {
			log.LogAttrs(ctx, slog.LevelError, "webull request done",
				slog.String("method", method),
				slog.String("path", reqPath),
				slog.String("correlation_id", correlationID),
				slog.Duration("latency", latency),
				slog.String("error", err.Error()),
			)
		} else {
			log.LogAttrs(ctx, slog.LevelInfo, "webull request done",
				slog.String("method", method),
				slog.String("path", reqPath),
				slog.String("correlation_id", correlationID),
				slog.Duration("latency", latency),
			)
		}
	}

	_ = start
	return err
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

// execute builds, signs, sends, and decodes a single request. It is the part of
// the request path that is replayed on retry.
func (c *Client) execute(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte, out any) error {
	req, err := c.buildSignedRequest(ctx, method, reqPath, query, bodyBytes)
	if err != nil {
		return err
	}

	resp, err := c.transport.Do(req)
	if err != nil {
		return errs.Wrap(errs.CodeTransport, method+" "+reqPath, err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errs.Wrap(errs.CodeTransport, "reading response body", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return errs.FromHTTPStatus(resp.StatusCode, data)
	}
	c.maybeUpdateClockOffset(resp.Header)
	if out == nil || len(data) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return errs.Wrap(errs.CodeAPI, "decoding response body", err)
	}
	return nil
}

// buildSignedRequest constructs the outgoing request for one API call and
// attaches its signing headers. It is shared by [Client.Do] and
// [Client.DoStream] so that a streamed request is built and signed exactly like
// a buffered one: the x-version header follows the same per-path defaults, the
// host that is signed is the host that is sent, and the exact bytes that are
// signed are the bytes that are transmitted.
func (c *Client) buildSignedRequest(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte) (*http.Request, error) {
	req, err := c.transport.NewRequest(ctx, method, reqPath, query, bodyBytes)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInvalidConfig, "building request", err)
	}

	// The host that is signed must be the host that is actually sent. Setting
	// req.Host makes the value explicit even when a custom (non-default) port
	// is in play.
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
	if raw, ok := body.([]byte); ok {
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
// updates the learned clock offset, clamped to the range [-5 minutes, +5 minutes].
// It is a no-op when clock-drift correction is disabled.
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
	offset := serverTime.Sub(time.Now())
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
	return time.Parse(http.TimeFormat, s)
}
