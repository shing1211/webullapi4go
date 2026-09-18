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
	"net/http"
	"time"

	"github.com/shing1211/webullapi4go/internal/errs"
	"github.com/shing1211/webullapi4go/internal/resilience/retry"
)

// Defaults applied by [DefaultConfig].
const (
	// DefaultTimeout is the per-request HTTP timeout.
	DefaultTimeout = 30 * time.Second
	// DefaultUserAgent is the User-Agent sent with every request.
	DefaultUserAgent = "webullapi4go"
	// DefaultRetryAttempts is the total number of attempts, including the
	// first, for retryable idempotent requests.
	DefaultRetryAttempts = 2
	// DefaultRetryBaseDelay is the delay before the first retry.
	DefaultRetryBaseDelay = 100 * time.Millisecond
	// DefaultRetryMaxDelay caps the exponential retry backoff.
	DefaultRetryMaxDelay = 2 * time.Second
)

// Config holds the resolved configuration for a [Client]. It is populated by
// [Option] functions and treated as immutable once the client is built.
//
// Every field uses a public type: callers outside the module never need to
// import an internal package.
type Config struct {
	// AppKey is the Webull OpenAPI app key.
	AppKey string
	// AppSecret is the Webull OpenAPI app secret.
	AppSecret string
	// Region is the deployment region.
	Region Region
	// Environment is the deployment environment.
	Environment Environment
	// Endpoints are the resolved service addresses.
	Endpoints Endpoints
	// HTTPClient is the client used for REST calls. When nil, [New] creates
	// one with the configured timeout.
	HTTPClient *http.Client
	// Timeout is the per-request timeout. It is applied only to the HTTP
	// client that [New] creates; a caller-supplied [Config.HTTPClient] keeps
	// its own timeout.
	Timeout time.Duration
	// UserAgent is the User-Agent header sent with every request.
	UserAgent string
	// APIVersion is the x-version header sent with every request when it has
	// been set explicitly with [WithAPIVersion]. Supported values are
	// [APIVersionV2] and [APIVersionV3]. When it is left at its default,
	// [Client] applies built-in per-path defaults instead: [APIVersionV3] for
	// paths under /trading/ and [APIVersionV2] everywhere else. Per-path
	// overrides registered with [WithAPIVersionFor] take precedence over both.
	APIVersion string

	// autoToken, when true, makes [Client.Do] obtain an access token
	// automatically before the first token-consuming request. It is enabled
	// with [WithAutoToken].
	autoToken bool
	// versionOverrides apply a per-path x-version override; the longest
	// matching prefix wins. [Config.APIVersion] is the fallback.
	versionOverrides []versionOverride
	// apiVersionSet records whether [WithAPIVersion] set
	// [Config.APIVersion] explicitly. When false, [Client.apiVersionFor]
	// applies the built-in per-path defaults instead of Config.APIVersion.
	apiVersionSet bool

	// endpointsOverride records whether an option set Endpoints or the base
	// URL explicitly. When false, [New] recomputes Endpoints from Region and
	// Environment so that those options take effect.
	endpointsOverride bool

	// The flags below record whether an option set the corresponding
	// environment-overridable field explicitly. [WithEnv] is a
	// lower-precedence default source: it fills a field only when the flag is
	// false, so an explicit option always wins regardless of the order in
	// which the options are passed to [New].
	appKeySet      bool
	appSecretSet   bool
	regionSet      bool
	environmentSet bool

	// retry executes transient idempotent requests. It is never nil unless
	// disabled with [WithoutRetry].
	retry *retry.Retrier
	// retryNonIdempotent allows retrying requests whose method is not
	// idempotent.
	retryNonIdempotent bool
	// rateLimiter, when non-nil, throttles outgoing requests per path.
	rateLimiter RateLimiter
	// breaker, when non-nil, gates outgoing requests.
	breaker CircuitBreaker
}

// DefaultConfig returns a [Config] pre-filled with production Hong Kong
// endpoints and the SDK's transport defaults.
func DefaultConfig() Config {
	return Config{
		Region:      HK,
		Environment: Production,
		Endpoints:   DefaultEndpoints(),
		Timeout:     DefaultTimeout,
		UserAgent:   DefaultUserAgent,
		APIVersion:  DefaultAPIVersion,
		retry:       defaultRetrier(),
	}
}

// Validate reports whether the configuration is usable. It does not make any
// network calls.
func (c Config) Validate() error {
	switch {
	case c.AppKey == "":
		return errs.New(errs.CodeInvalidConfig, "app key is required")
	case c.AppSecret == "":
		return errs.New(errs.CodeInvalidConfig, "app secret is required")
	case !c.Region.IsValid():
		return errs.New(errs.CodeInvalidConfig, "unknown region: "+string(c.Region))
	case !c.Environment.IsValid():
		return errs.New(errs.CodeInvalidConfig, "unknown environment: "+string(c.Environment))
	case c.Endpoints.HTTP == "":
		return errs.New(errs.CodeInvalidConfig, "http endpoint is required")
	case !IsValidAPIVersion(c.APIVersion):
		return errs.New(errs.CodeInvalidConfig, "unknown api version: "+c.APIVersion)
	}
	for _, o := range c.versionOverrides {
		if o.prefix == "" {
			return errs.New(errs.CodeInvalidConfig, "api version override prefix is required")
		}
		if !IsValidAPIVersion(o.version) {
			return errs.New(errs.CodeInvalidConfig, "unknown api version override: "+o.version)
		}
	}
	return nil
}
