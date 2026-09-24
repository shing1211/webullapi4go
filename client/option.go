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
	"log/slog"
	"net/http"
	"time"

	"github.com/shing1211/webullapi4go/pkg/observability"
	"github.com/shing1211/webullapi4go/pkg/resilience/breaker"
	"github.com/shing1211/webullapi4go/pkg/resilience/ratelimit"
	"github.com/shing1211/webullapi4go/pkg/resilience/retry"
)

// Option mutates a [Config] during [New]. Options are applied in order on top
// of [DefaultConfig], so a later option overrides an earlier one.
type Option func(*Config)

// WithAppKey sets the Webull OpenAPI app key.
func WithAppKey(appKey string) Option {
	return func(c *Config) {
		c.AppKey = appKey
		c.appKeySet = true
	}
}

// WithAppSecret sets the Webull OpenAPI app secret.
func WithAppSecret(appSecret string) Option {
	return func(c *Config) {
		c.AppSecret = appSecret
		c.appSecretSet = true
	}
}

// WithCredentials sets both the app key and the app secret.
func WithCredentials(appKey, appSecret string) Option {
	return func(c *Config) {
		c.AppKey, c.AppSecret = appKey, appSecret
		c.appKeySet, c.appSecretSet = true, true
	}
}

// WithRegion sets the deployment region.
func WithRegion(r Region) Option {
	return func(c *Config) {
		c.Region = r
		c.regionSet = true
	}
}

// WithEnvironment sets the deployment environment.
func WithEnvironment(env Environment) Option {
	return func(c *Config) {
		c.Environment = env
		c.environmentSet = true
	}
}

// WithSandbox switches the client to the sandbox environment.
func WithSandbox() Option {
	return WithEnvironment(Sandbox)
}

// WithEndpoints overrides the service endpoints resolved from the region and
// environment.
func WithEndpoints(e Endpoints) Option {
	return func(c *Config) {
		c.Endpoints = e
		c.endpointsOverride = true
	}
}

// WithBaseURL overrides only the REST base URL, for example
// "https://api.sandbox.webull.hk" or a local test server. It takes precedence
// over the address resolved from the region and environment.
func WithBaseURL(baseURL string) Option {
	return func(c *Config) {
		c.Endpoints.HTTP = baseURL
		c.endpointsOverride = true
	}
}

// WithHTTPClient sets the HTTP client used for REST calls. When omitted, [New]
// creates one with the configured timeout.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Config) { c.HTTPClient = hc }
}

// WithHTTPTransport sets the [http.Transport] used for REST calls. It is
// applied after [WithHTTPClient]; if no HTTP client has been set, [New] creates
// one with the configured timeout before applying the transport. This allows
// callers to tune connection pooling, TLS, timeouts, and other transport-level
// parameters without replacing the entire HTTP client.
func WithHTTPTransport(tr *http.Transport) Option {
	return func(c *Config) { c.httpTransport = tr }
}

// WithTimeout sets the per-request timeout used when [New] creates the HTTP
// client.
func WithTimeout(d time.Duration) Option {
	return func(c *Config) { c.Timeout = d }
}

// WithUserAgent sets the User-Agent header sent with every request.
func WithUserAgent(ua string) Option {
	return func(c *Config) { c.UserAgent = ua }
}

// RetryConfig configures the retries performed by [Client.Do] for transient
// failures. Zero-valued duration and attempt fields fall back to the SDK
// defaults ([DefaultRetryAttempts], [DefaultRetryBaseDelay],
// [DefaultRetryMaxDelay]); a nil IsRetryable falls back to the SDK classifier.
// To disable retries entirely use [WithoutRetry].
type RetryConfig struct {
	// MaxAttempts is the total number of attempts, including the first. One
	// disables retrying while keeping the retry machinery active.
	MaxAttempts int
	// BaseDelay is the delay before the first retry and the base of the
	// exponential backoff.
	BaseDelay time.Duration
	// MaxDelay caps the computed delay.
	MaxDelay time.Duration
	// Jitter randomizes each delay by plus or minus 20 percent.
	Jitter bool
	// RetryNonIdempotent allows retrying methods that are not idempotent
	// (for example POST). It is disabled by default so that only GET and HEAD
	// requests are retried.
	RetryNonIdempotent bool
	// IsRetryable reports whether an error is transient. When nil the SDK
	// classifier retries transport failures, HTTP 429, and HTTP 5xx.
	IsRetryable func(error) bool
}

// WithRetry replaces the default retry policy. Only idempotent requests are
// retried unless [RetryConfig.RetryNonIdempotent] is set.
func WithRetry(cfg RetryConfig) Option {
	if cfg.MaxAttempts == 0 {
		cfg.MaxAttempts = DefaultRetryAttempts
	}
	if cfg.BaseDelay == 0 {
		cfg.BaseDelay = DefaultRetryBaseDelay
	}
	if cfg.MaxDelay == 0 {
		cfg.MaxDelay = DefaultRetryMaxDelay
	}
	if cfg.IsRetryable == nil {
		cfg.IsRetryable = retry.DefaultIsRetryable
	}
	return func(c *Config) {
		c.retry = retry.NewWithConfig(retry.Config{
			MaxAttempts: cfg.MaxAttempts,
			BaseDelay:   cfg.BaseDelay,
			MaxDelay:    cfg.MaxDelay,
			Jitter:      cfg.Jitter,
			IsRetryable: cfg.IsRetryable,
		})
		c.retryNonIdempotent = cfg.RetryNonIdempotent
	}
}

// WithoutRetry disables retries. Requests are attempted exactly once.
func WithoutRetry() Option {
	return func(c *Config) {
		c.retry = nil
		c.retryNonIdempotent = false
	}
}

// defaultRetrier returns the SDK's default retry policy: up to
// [DefaultRetryAttempts] attempts for transient errors.
func defaultRetrier() *retry.Retrier {
	return retry.NewWithConfig(retry.Config{
		MaxAttempts: DefaultRetryAttempts,
		BaseDelay:   DefaultRetryBaseDelay,
		MaxDelay:    DefaultRetryMaxDelay,
		Jitter:      true,
		IsRetryable: retry.DefaultIsRetryable,
	})
}

// RateLimiter throttles outgoing requests. [Client.Do] calls Wait before every
// attempt, passing the request path as the key so implementations can apply
// per-endpoint limits. Wait must return ctx.Err() when ctx is done.
type RateLimiter interface {
	Wait(ctx context.Context, key string) error
}

// WithRateLimiter sets the rate limiter consulted before every request
// attempt. A nil limiter disables rate limiting (the default).
func WithRateLimiter(l RateLimiter) Option {
	return func(c *Config) { c.rateLimiter = l }
}

// NewRateLimiter returns a keyed token-bucket limiter that allows rate tokens
// per second with the given burst, applied independently per request path.
// NewRateLimiter panics if rate or burst is not positive.
func NewRateLimiter(rate float64, burst int) RateLimiter {
	return ratelimit.NewKeyed(func() *ratelimit.Limiter {
		return ratelimit.New(rate, burst)
	})
}

// CircuitBreaker gates outgoing requests. Allow reports whether a call may
// proceed; RecordSuccess and RecordFailure report its outcome.
type CircuitBreaker interface {
	Allow() bool
	RecordSuccess()
	RecordFailure()
}

// WithBreaker sets the circuit breaker consulted before every request attempt.
// A nil breaker disables circuit breaking (the default).
func WithBreaker(b CircuitBreaker) Option {
	return func(c *Config) { c.breaker = b }
}

// NewBreaker returns a circuit breaker that opens after threshold consecutive
// server or transport failures and probes again after cooldown.
func NewBreaker(threshold int, cooldown time.Duration) CircuitBreaker {
	return breaker.New(
		breaker.WithThreshold(threshold),
		breaker.WithCooldown(cooldown),
	)
}

// WithClockDriftCorrection enables clock-drift correction. When enabled, the
// client learns the offset between the local clock and the Webull server's clock
// by observing the Date response header on each successful request, and applies
// that offset to subsequent x-timestamp values used in request signing. The offset
// is clamped to the range [-5 minutes, +5 minutes] to prevent extreme offsets from
// being injected. Disable with false to opt out.
func WithClockDriftCorrection(enabled bool) Option {
	return func(c *Config) {
		c.clockDriftCorrection = enabled
	}
}

// Interceptor is a function that wraps the request pipeline. It receives the next
// function in the chain and may inspect, decorate, or short-circuit the request.
// Interceptors are invoked in the order they are supplied to [WithInterceptor].
// The first interceptor in the chain receives a function that performs the
// underlying HTTP call (rate-limit, breaker, signing, and send).
//
// For example, a logging interceptor:
//
//	func(ctx context.Context, next func(context.Context) error) error {
//	    start := time.Now()
//	    err := next(ctx)
//	    log.Printf("request took %s: %v", time.Since(start), err)
//	    return err
//	}
type Interceptor func(ctx context.Context, next func(context.Context) error) error

// WithInterceptor adds an interceptor to the request pipeline. Interceptors are
// invoked after rate-limiting and circuit-breaking but before the request is
// signed and sent, and after the response is received. Multiple interceptors
// can be added; they fire in the order they are supplied.
func WithInterceptor(i Interceptor) Option {
	return func(c *Config) {
		c.interceptors = append(c.interceptors, i)
	}
}

// Hooks holds optional lifecycle callbacks invoked by [Client.Do] around each
// request attempt. All fields are optional; nil fields are no-ops.
//
// These hooks are the integration point for observability tools (structured
// logging, OpenTelemetry tracing, Prometheus metrics). Phase 6 of the
// production hardening plan adds actual instrumentation via these hooks.
type Hooks struct {
	// OnRequest is called before each attempt, after rate-limiting and
	// circuit-breaking, with the HTTP method and path.
	OnRequest func(method, path string)
	// OnResponse is called on a successful response (2xx), passing the
	// HTTP status code and the round-trip latency.
	OnResponse func(status int, latency time.Duration)
	// OnError is called when the request returns an error or a non-retryable
	// HTTP status, passing the error and the round-trip latency.
	OnError func(err error, latency time.Duration)
	// OnLatency is called after every attempt (success, error, or retry) with
	// the attempt number and the round-trip latency.
	OnLatency func(attempt int, latency time.Duration)
}

// WithHooks installs lifecycle hooks for observability integration.
func WithHooks(h Hooks) Option {
	return func(c *Config) {
		c.hooks = h
	}
}

// ResiliencePreset applies a named set of resilience defaults. Presets are
// composable: each call adds to or overrides the existing configuration.
type ResiliencePreset string

const (
	// ProductionPreset applies sensible production defaults: a per-path rate
	// limiter (10 requests per second, burst 20), a circuit breaker (opens after
	// 5 consecutive failures, 30-second cooldown), and exponential backoff retry
	// with full jitter (base 200 ms, cap 2 s, up to 3 attempts).
	ProductionPreset ResiliencePreset = "production"
)

// WithResiliencePreset applies a named resilience preset. Currently only
// [ProductionPreset] is defined. The preset is applied by setting a per-path
// rate limiter, circuit breaker, and retry policy; earlier options take precedence
// so callers can override individual components.
func WithResiliencePreset(preset ResiliencePreset) Option {
	return func(c *Config) {
		switch preset {
		case ProductionPreset:
			c.rateLimiter = NewRateLimiter(10, 20)
			c.breaker = breaker.NewWithConfig(breaker.Config{
				Threshold: 5,
				Cooldown:  30 * time.Second,
				Meter:     c.otel.Meter("webullapi4go/resilience"),
			})
			c.retry = retry.New(
				retry.WithMaxAttempts(3),
				retry.WithBaseDelay(200*time.Millisecond),
				retry.WithMaxDelay(2*time.Second),
				retry.WithFullJitter(true),
				retry.WithIsRetryable(retry.DefaultIsRetryable),
			)
		}
	}
}

// WithLogger sets the structured logger used by [Client.Do] for per-request
// log output. When nil (the default), no structured logging is produced.
// A no-op logger is used when a real logger is not supplied.
func WithLogger(log *slog.Logger) Option {
	return func(c *Config) {
		c.otel.Logger = log
	}
}

// WithTracerProvider sets the OpenTelemetry [TracerProvider] used to create
// spans for [Client.Do] and [Client.DoBroker] requests. When nil (the default),
// the global no-op tracer is used.
func WithTracerProvider(tp observability.TracerProvider) Option {
	return func(c *Config) {
		c.otel.TracerProvider = tp
	}
}

// WithMeterProvider sets the OpenTelemetry [MeterProvider] used to create
// meters for SDK-level metrics. When nil (the default), the global no-op
// meter is used.
func WithMeterProvider(mp observability.MeterProvider) Option {
	return func(c *Config) {
		c.otel.MeterProvider = mp
	}
}

// WithPropagator sets the OpenTelemetry propagator used to extract and inject
// trace context on requests. When nil (the default), no trace context
// propagation is performed.
func WithPropagator(p observability.TextMapPropagator) Option {
	return func(c *Config) {
		c.otel.Propagator = p
	}
}
