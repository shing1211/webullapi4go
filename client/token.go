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
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/shing1211/webullapi4go/internal/auth"
	"github.com/shing1211/webullapi4go/pkg/errors"
)

// Token endpoint paths, relative to the configured HTTP base URL.
const (
	createTokenPath = "/openapi/auth/token/create"
	checkTokenPath  = "/openapi/auth/token/check"
)

// AccessTokenHeader is the request header that carries an active access token,
// as documented in Webull's authentication guide. It is not part of the request
// signature, so it can be attached after signing.
const AccessTokenHeader = "x-access-token"

// Default token polling parameters. Webull issues a token as PENDING and gives
// the user five minutes to complete 2FA, so the default timeout matches that
// window. The default interval respects the documented limit of 10 requests per
// 30 seconds for the token endpoints.
const (
	// DefaultTokenPollInterval is the delay between check-token requests while
	// waiting for a PENDING token to become NORMAL.
	DefaultTokenPollInterval = 5 * time.Second
	// DefaultTokenPollTimeout is how long [Client.EnsureToken] waits for a
	// PENDING token to become NORMAL before giving up.
	DefaultTokenPollTimeout = 5 * time.Minute
)

// TokenStatus is the lifecycle state of an access token, mirrored from the
// Webull TokenRespVo schema so that callers never need an internal package.
type TokenStatus string

// Token lifecycle states. A token is created PENDING and becomes NORMAL once
// 2FA is completed; sandbox tokens are NORMAL immediately. INVALID means the
// token was unused for 15 consecutive days (or does not exist) and EXPIRED
// means 2FA was not completed within five minutes.
const (
	// TokenStatusPending means the token awaits verification.
	TokenStatusPending TokenStatus = "PENDING"
	// TokenStatusNormal means the token is valid and usable.
	TokenStatusNormal TokenStatus = "NORMAL"
	// TokenStatusInvalid means the token is invalid or was never used.
	TokenStatusInvalid TokenStatus = "INVALID"
	// TokenStatusExpired means the verification window elapsed.
	TokenStatusExpired TokenStatus = "EXPIRED"
)

// String returns the wire representation of the status.
func (s TokenStatus) String() string { return string(s) }

// Valid reports whether s is one of the four statuses defined by Webull.
func (s TokenStatus) Valid() bool {
	switch s {
	case TokenStatusPending, TokenStatusNormal, TokenStatusInvalid, TokenStatusExpired:
		return true
	default:
		return false
	}
}

// Token is an access token and its lifecycle metadata. It is the public mirror
// of the internal token DTO; every exported signature in this package uses this
// type.
type Token struct {
	// Value is the access token string sent as the [AccessTokenHeader] header.
	Value string
	// ExpiresAt is the token expiry. The zero time means unspecified.
	ExpiresAt time.Time
	// Status is the current lifecycle state.
	Status TokenStatus
}

// Valid reports whether the token is usable right now: it has a value, its
// status is [TokenStatusNormal], and it has not expired.
func (t Token) Valid() bool {
	if t.Value == "" || t.Status != TokenStatusNormal {
		return false
	}
	if t.ExpiresAt.IsZero() {
		return true
	}
	return time.Now().Before(t.ExpiresAt)
}

// Expired reports whether the token has an expiry that is at or before now. A
// token without an expiry is never considered expired.
func (t Token) Expired() bool {
	return !t.ExpiresAt.IsZero() && !time.Now().Before(t.ExpiresAt)
}

// newToken converts an internal token DTO into the public [Token].
func newToken(t auth.Token) *Token {
	out := &Token{Value: t.Value, Status: TokenStatus(t.Status)}
	if t.ExpiresAt != 0 {
		out.ExpiresAt = time.UnixMilli(t.ExpiresAt)
	}
	return out
}

// CreateToken creates a new access token with
// POST /openapi/auth/token/create.
//
// In production the returned token is [TokenStatusPending] until the user
// completes the Webull App 2FA flow (see [Client.EnsureToken] for an
// alternative that waits for that to happen); in the sandbox environment the
// token is [TokenStatusNormal] immediately. CreateToken does not touch the
// client's token cache; most callers want [Client.EnsureToken].
func (c *Client) CreateToken(ctx context.Context) (*Token, error) {
	var resp auth.Token
	if err := c.Do(ctx, http.MethodPost, createTokenPath, auth.CreateTokenRequest{}, &resp); err != nil {
		return nil, err
	}
	return newToken(resp), nil
}

// CheckToken returns the current status of token with
// POST /openapi/auth/token/check. It does not touch the client's token cache.
func (c *Client) CheckToken(ctx context.Context, token string) (*Token, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errs.New(errs.CodeInvalidConfig, "token is required")
	}
	var resp auth.Token
	if err := c.Do(ctx, http.MethodPost, checkTokenPath, auth.CheckTokenRequest{Token: token}, &resp); err != nil {
		return nil, err
	}
	return newToken(resp), nil
}

// EnsureToken returns a usable access token, creating and, when necessary,
// waiting for one to become [TokenStatusNormal].
//
// A cached, still-valid token is reused without any network call. Otherwise a
// token is created: in the sandbox it is expected to be NORMAL immediately and
// is returned at once; in production a PENDING token is polled with
// [Client.CheckToken] until it becomes NORMAL, becomes INVALID or EXPIRED (an
// error is returned), or the poll timeout elapses. The resolved token is cached
// and installed as the [AccessTokenHeader] value for subsequent [Client.Do]
// calls.
//
// Polling is bounded by [DefaultTokenPollInterval] and
// [DefaultTokenPollTimeout], which can be changed with
// [Client.SetTokenPollInterval] and [Client.SetTokenPollTimeout].
func (c *Client) EnsureToken(ctx context.Context) (*Token, error) {
	state := tokenStateFor(c)
	if cached := state.get(); cached != nil && cached.Valid() {
		c.enableTokenInjection(state)
		return cached, nil
	}
	return c.fetchToken(ctx, state)
}

// fetchToken creates a token and waits for it to become usable.
func (c *Client) fetchToken(ctx context.Context, state *tokenState) (*Token, error) {
	created, err := c.CreateToken(ctx)
	if err != nil {
		return nil, err
	}
	if created.Status == TokenStatusNormal {
		return c.storeToken(state, created), nil
	}
	if c.Environment() == Sandbox {
		// Sandbox tokens are NORMAL by default, so a non-NORMAL token means
		// something is wrong with the account. Fail fast instead of polling.
		return created, errs.New(errs.CodeAuth, "sandbox token is not NORMAL: "+created.Status.String())
	}
	return c.pollToken(ctx, state, created.Value)
}

// pollToken repeatedly checks token until it becomes NORMAL or a terminal
// status or timeout is reached.
func (c *Client) pollToken(ctx context.Context, state *tokenState, value string) (*Token, error) {
	interval, timeout := state.pollConfig()
	deadline := time.Now().Add(timeout)
	for {
		checked, err := c.CheckToken(ctx, value)
		if err != nil {
			return nil, err
		}
		switch checked.Status {
		case TokenStatusNormal:
			return c.storeToken(state, checked), nil
		case TokenStatusInvalid, TokenStatusExpired:
			return nil, errs.New(errs.CodeAuth, "token is "+checked.Status.String())
		default:
			// PENDING (or an unrecognized status): keep polling.
		}
		if time.Now().After(deadline) {
			return nil, errs.New(errs.CodeAuth, "timed out waiting for token verification")
		}
		if err := sleepContext(ctx, interval); err != nil {
			return nil, err
		}
	}
}

// storeToken caches t and enables automatic header injection.
func (c *Client) storeToken(state *tokenState, t *Token) *Token {
	cp := *t
	state.set(&cp)
	c.enableTokenInjection(state)
	return &cp
}

// CurrentToken returns a copy of the cached access token, or nil when none is
// cached. It never makes a network call.
func (c *Client) CurrentToken() *Token {
	t := tokenStateFor(c).get()
	if t == nil {
		return nil
	}
	cp := *t
	return &cp
}

// AccessToken returns the cached access token value, or an empty string when no
// token is cached. It never makes a network call.
func (c *Client) AccessToken() string {
	if t := tokenStateFor(c).get(); t != nil {
		return t.Value
	}
	return ""
}

// SetToken caches token and enables automatic header injection, replacing any
// previously cached token. Pass nil to clear the cache; a cleared cache also
// disables the token header for subsequent requests.
func (c *Client) SetToken(token *Token) {
	state := tokenStateFor(c)
	if token == nil {
		state.set(nil)
		return
	}
	cp := *token
	state.set(&cp)
	c.enableTokenInjection(state)
}

// SetTokenPollInterval sets the delay between check-token requests used by
// [Client.EnsureToken]. A non-positive duration is ignored. It should be set
// before polling begins.
func (c *Client) SetTokenPollInterval(d time.Duration) {
	if d <= 0 {
		return
	}
	state := tokenStateFor(c)
	state.mu.Lock()
	state.interval = d
	state.mu.Unlock()
}

// SetTokenPollTimeout sets how long [Client.EnsureToken] waits for a PENDING
// token to become NORMAL. A non-positive duration is ignored. It should be set
// before polling begins.
func (c *Client) SetTokenPollTimeout(d time.Duration) {
	if d <= 0 {
		return
	}
	state := tokenStateFor(c)
	state.mu.Lock()
	state.timeout = d
	state.mu.Unlock()
}

// EnableTokenInjection installs the transport hook that attaches the cached
// access token as the [AccessTokenHeader] header on outgoing requests. It is
// called automatically by [Client.EnsureToken] and [Client.SetToken], and is
// idempotent.
//
// INTEGRATION SEAM: this is how the SDK attaches x-access-token without
// changing the core request pipeline. The hook is an [http.RoundTripper] that
// wraps the client's existing transport, so it applies to every request sent
// through [Client.Do] (including the token endpoints, which are deliberately
// skipped). The header is added after signing and is therefore not part of the
// signature. If the request pipeline later grows an explicit token hook, it can
// consult [Client.AccessToken] instead; this wrapper remains compatible.
func (c *Client) EnableTokenInjection() {
	c.enableTokenInjection(tokenStateFor(c))
}

// enableTokenInjection installs the token transport for state exactly once.
func (c *Client) enableTokenInjection(state *tokenState) {
	state.injectOnce.Do(func() {
		hc := c.HTTPClient()
		if hc == nil {
			return
		}
		base := hc.Transport
		if base == nil {
			base = http.DefaultTransport
		}
		hc.Transport = &tokenTransport{base: base, state: state}
	})
}

// tokenTransport injects the cached access token into outgoing requests.
type tokenTransport struct {
	base  http.RoundTripper
	state *tokenState
}

// RoundTrip implements [http.RoundTripper]. It adds [AccessTokenHeader] when a
// token is cached, except for the token endpoints themselves.
func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.state != nil {
		if token := t.state.get(); token != nil && token.Value != "" && !isTokenEndpoint(req.URL.Path) {
			req.Header.Set(AccessTokenHeader, token.Value)
		}
	}
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}

// isTokenEndpoint reports whether path is one of the token lifecycle endpoints,
// which must not carry a previously cached token.
func isTokenEndpoint(path string) bool {
	return strings.HasPrefix(path, createTokenPath) || strings.HasPrefix(path, checkTokenPath)
}

// tokenState is the per-client token cache and polling configuration. It is
// embedded in [Client] so that no global, address-keyed registry is required.
type tokenState struct {
	mu         sync.RWMutex
	token      *Token
	interval   time.Duration
	timeout    time.Duration
	injectOnce sync.Once
}

// tokenStateFor returns c's own token state. It exists to keep the call sites
// uniform; the state is a field of the client, so there is no shared map and no
// possibility of observing another client's state.
func tokenStateFor(c *Client) *tokenState {
	return &c.tok
}

// get returns the cached token, or nil.
func (s *tokenState) get() *Token {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.token
}

// set replaces the cached token.
func (s *tokenState) set(t *Token) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = t
}

// pollConfig returns the current polling interval and timeout, falling back to
// the package defaults when they have not been set. This keeps the zero value of
// tokenState usable.
func (s *tokenState) pollConfig() (interval, timeout time.Duration) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	interval, timeout = s.interval, s.timeout
	if interval <= 0 {
		interval = DefaultTokenPollInterval
	}
	if timeout <= 0 {
		timeout = DefaultTokenPollTimeout
	}
	return interval, timeout
}

// sleepContext sleeps for d unless ctx is cancelled first.
func sleepContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
