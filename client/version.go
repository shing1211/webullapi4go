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
	"strings"

	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

// API interface versions accepted by the Webull OpenAPI x-version header.
const (
	// APIVersionV2 is the legacy API version and the SDK default.
	APIVersionV2 = "v2"
	// APIVersionV3 is the current API version used by the newer reference
	// endpoints, notably the /trading/ paths and the market-data
	// fundamentals and futures static data paths.
	APIVersionV3 = "v3"
	// DefaultAPIVersion is the x-version value used for request paths that
	// have no more specific built-in default, such as the market-data paths.
	// It remains v2 so that existing behaviour is preserved.
	DefaultAPIVersion = APIVersionV2
)

// IsValidAPIVersion reports whether v is one of the API versions supported by
// the SDK.
func IsValidAPIVersion(v string) bool {
	return v == APIVersionV2 || v == APIVersionV3
}

// ErrAccessTokenRequired is returned by [Client.Do] when [WithAutoToken] is
// enabled, no usable access token is cached, and the client targets the
// production environment. Production tokens require the caller to complete the
// Webull App 2FA flow, which the SDK refuses to start implicitly; call
// [Client.EnsureToken] once to create and activate a token, then retry.
var ErrAccessTokenRequired = errs.New(
	errs.CodeAuth,
	"access token required: call client.EnsureToken to complete the 2FA flow before running in production",
)

// versionOverride pins a non-default API version for requests whose path
// matches a prefix.
type versionOverride struct {
	prefix  string
	version string
}

// WithAPIVersion sets the x-version header sent with every request, taking
// precedence over the SDK's built-in per-path defaults. Supported values are
// [APIVersionV2] and [APIVersionV3]. When it is not called, paths under
// /trading/ default to [APIVersionV3] and every other path to
// [DefaultAPIVersion] ("v2"). [New] rejects any unsupported version.
func WithAPIVersion(v string) Option {
	return func(c *Config) {
		c.APIVersion = v
		c.apiVersionSet = true
	}
}

// WithAPIVersionFor overrides the x-version header for requests whose path
// starts with pathPrefix. The longest matching prefix wins, so a single
// endpoint can opt into [APIVersionV3] while the rest of the client stays on
// the configured default. pathPrefix must be non-empty and v must be a
// supported version; [New] rejects anything else.
func WithAPIVersionFor(pathPrefix, v string) Option {
	return func(c *Config) {
		c.versionOverrides = append(c.versionOverrides, versionOverride{prefix: pathPrefix, version: v})
	}
}

// WithAutoToken enables or disables automatic access-token handling for the
// client. It is off by default, so [Client.Do], [Client.DoBroker], and
// [Client.DoStream] never obtain a token on their own and only attach one cached
// by [Client.EnsureToken] or [Client.SetToken].
//
// When enabled and no usable token is cached, the first request that is not a
// token-lifecycle endpoint behaves as follows:
//
//   - in the sandbox environment [Client.EnsureToken] is called automatically
//     (sandbox tokens are activated without 2FA) and the token is attached to
//     the request and every later one;
//   - in production [Client.Do] fails with [ErrAccessTokenRequired], leaving
//     the caller to drive the 2FA flow explicitly with [Client.EnsureToken].
//
// Production 2FA is never started implicitly.
func WithAutoToken(enabled bool) Option {
	return func(c *Config) { c.autoToken = enabled }
}

// tradingPathPrefix identifies the Webull trading API paths. Requests under it
// default to [APIVersionV3], the version documented for the /trading/ endpoints.
const tradingPathPrefix = "/trading/"

// defaultAPIVersionFor returns the built-in x-version for path: [APIVersionV3]
// for the trading API and [DefaultAPIVersion] otherwise.
func defaultAPIVersionFor(path string) string {
	if strings.HasPrefix(path, tradingPathPrefix) {
		return APIVersionV3
	}
	return DefaultAPIVersion
}

// apiVersionFor returns the x-version value for path. The configured default
// applies only when [WithAPIVersion] was passed explicitly; otherwise the
// built-in per-path default does. In both cases the longest matching
// [WithAPIVersionFor] override wins.
func (c *Client) apiVersionFor(path string) string {
	version := defaultAPIVersionFor(path)
	if c.cfg.apiVersionSet {
		version = c.cfg.APIVersion
	}
	best := -1
	for _, o := range c.cfg.versionOverrides {
		if len(o.prefix) > best && strings.HasPrefix(path, o.prefix) {
			best = len(o.prefix)
			version = o.version
		}
	}
	return version
}

// ensureAutoToken makes sure a usable access token is cached before a
// token-consuming request when [WithAutoToken] is enabled. Token-lifecycle
// endpoints are exempt so that [Client.EnsureToken] can create the first token
// without recursing.
func (c *Client) ensureAutoToken(ctx context.Context, path string) error {
	if !c.cfg.autoToken || isTokenEndpoint(path) {
		return nil
	}
	state := tokenStateFor(c)
	if c.cachedValidToken(state) {
		return nil
	}

	c.autoTokenMu.Lock()
	defer c.autoTokenMu.Unlock()

	// Another request may have populated the cache while we waited.
	if c.cachedValidToken(state) {
		return nil
	}
	// Only the sandbox activates tokens without 2FA. In production the caller
	// must run [Client.EnsureToken] themselves.
	if c.cfg.Environment != Sandbox {
		return ErrAccessTokenRequired
	}
	_, err := c.EnsureToken(ctx)
	return err
}

// cachedValidToken reports whether state holds a usable token and, when it
// does, makes sure the injection transport is installed.
func (c *Client) cachedValidToken(state *tokenState) bool {
	t := state.get()
	if t == nil || !t.Valid() {
		return false
	}
	c.enableTokenInjection(state)
	return true
}
