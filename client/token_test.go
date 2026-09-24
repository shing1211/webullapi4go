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
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/internal/auth"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

const (
	tokenTestAppKey    = "token-test-app-key"
	tokenTestAppSecret = "token-test-app-secret"
)

// tokenResponse is a canned TokenRespVo body for the fake server.
type tokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	Status    string `json:"status"`
}

// tokenServer is a fake Webull token API. It records how often each endpoint is
// called and the body/token header of the corresponding requests.
type tokenServer struct {
	create tokenResponse
	check  tokenResponse

	mu          sync.Mutex
	createCalls int
	checkCalls  int
	checkBodies []string
	accessToken string
}

func (s *tokenServer) handler(t *testing.T) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		assertTokenSigningHeaders(t, r)
		verifyTokenSignature(t, r, body)

		switch r.URL.Path {
		case "/openapi/auth/token/create":
			s.mu.Lock()
			s.createCalls++
			s.mu.Unlock()
			writeTokenJSON(t, w, s.create)
		case "/openapi/auth/token/check":
			s.mu.Lock()
			s.checkCalls++
			s.checkBodies = append(s.checkBodies, string(body))
			s.mu.Unlock()
			writeTokenJSON(t, w, s.check)
		case "/openapi/account/list":
			s.mu.Lock()
			s.accessToken = r.Header.Get(client.AccessTokenHeader)
			s.mu.Unlock()
			_, _ = w.Write([]byte(`[]`))
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
			http.NotFound(w, r)
		}
	}
}

func (s *tokenServer) counts() (create, check int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.createCalls, s.checkCalls
}

func (s *tokenServer) header() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.accessToken
}

func writeTokenJSON(t *testing.T, w http.ResponseWriter, r tokenResponse) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(r); err != nil {
		t.Errorf("encoding token response: %v", err)
	}
}

// assertTokenSigningHeaders checks that every request carries the standard
// signing headers.
func assertTokenSigningHeaders(t *testing.T, r *http.Request) {
	t.Helper()
	for _, name := range []string{
		"x-app-key", "x-timestamp", "x-signature",
		"x-signature-algorithm", "x-signature-version", "x-signature-nonce", "x-version",
	} {
		if r.Header.Get(name) == "" {
			t.Errorf("missing signing header %q on %s", name, r.URL.Path)
		}
	}
	if got := r.Header.Get("x-app-key"); got != tokenTestAppKey {
		t.Errorf("x-app-key = %q, want %q", got, tokenTestAppKey)
	}
}

// verifyTokenSignature recomputes the HMAC-SHA1 signature and compares it to the
// x-signature header, proving the request was signed correctly.
func verifyTokenSignature(t *testing.T, r *http.Request, body []byte) {
	t.Helper()
	headers := make(http.Header)
	headers.Set(auth.HeaderAppKey, r.Header.Get(auth.HeaderAppKey))
	headers.Set(auth.HeaderSignatureAlgorithm, r.Header.Get(auth.HeaderSignatureAlgorithm))
	headers.Set(auth.HeaderSignatureVersion, r.Header.Get(auth.HeaderSignatureVersion))
	headers.Set(auth.HeaderSignatureNonce, r.Header.Get(auth.HeaderSignatureNonce))
	headers.Set(auth.HeaderTimestamp, r.Header.Get(auth.HeaderTimestamp))
	headers.Set(auth.HeaderHost, r.Host)

	got, err := auth.Sign(auth.SignParams{
		Method:    r.Method,
		Path:      r.URL.Path,
		Query:     r.URL.Query(),
		Headers:   headers,
		Body:      body,
		AppSecret: tokenTestAppSecret,
	})
	if err != nil {
		t.Errorf("recomputing signature: %v", err)
		return
	}
	if want := r.Header.Get("x-signature"); got != want {
		t.Errorf("signature mismatch on %s: recomputed %q, header %q", r.URL.Path, got, want)
	}
}

func newTokenClient(t *testing.T, baseURL string, opts ...client.Option) *client.Client {
	t.Helper()
	all := append([]client.Option{
		client.WithAppKey(tokenTestAppKey),
		client.WithAppSecret(tokenTestAppSecret),
		client.WithBaseURL(baseURL),
	}, opts...)
	cl, err := client.New(all...)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })
	return cl
}

func futureMillis() int64 { return time.Now().Add(24 * time.Hour).UnixMilli() }

func TestCreateTokenReturnsPending(t *testing.T) {
	t.Parallel()

	srv := &tokenServer{
		create: tokenResponse{Token: "tok-pending", ExpiresAt: futureMillis(), Status: "PENDING"},
	}
	ts := httptest.NewServer(srv.handler(t))
	defer ts.Close()

	cl := newTokenClient(t, ts.URL)
	tok, err := cl.CreateToken(context.Background())
	if err != nil {
		t.Fatalf("CreateToken() error = %v", err)
	}
	if tok.Value != "tok-pending" {
		t.Errorf("Value = %q, want %q", tok.Value, "tok-pending")
	}
	if tok.Status != client.TokenStatusPending {
		t.Errorf("Status = %q, want %q", tok.Status, client.TokenStatusPending)
	}
	if tok.ExpiresAt.IsZero() || !tok.ExpiresAt.After(time.Now()) {
		t.Errorf("ExpiresAt = %v, want a future time", tok.ExpiresAt)
	}
	if tok.Valid() {
		t.Error("Valid() = true for a PENDING token, want false")
	}
}

func TestCheckTokenSendsTokenAndParsesStatus(t *testing.T) {
	t.Parallel()

	srv := &tokenServer{
		check: tokenResponse{Token: "tok-normal", ExpiresAt: futureMillis(), Status: "NORMAL"},
	}
	ts := httptest.NewServer(srv.handler(t))
	defer ts.Close()

	cl := newTokenClient(t, ts.URL)
	tok, err := cl.CheckToken(context.Background(), "tok-normal")
	if err != nil {
		t.Fatalf("CheckToken() error = %v", err)
	}
	if tok.Status != client.TokenStatusNormal {
		t.Errorf("Status = %q, want %q", tok.Status, client.TokenStatusNormal)
	}
	if !tok.Valid() {
		t.Error("Valid() = false for a NORMAL future token, want true")
	}
	srv.mu.Lock()
	bodies := append([]string(nil), srv.checkBodies...)
	srv.mu.Unlock()
	if len(bodies) != 1 {
		t.Fatalf("check calls = %d, want 1", len(bodies))
	}
	var got map[string]string
	if err := json.Unmarshal([]byte(bodies[0]), &got); err != nil {
		t.Fatalf("check body %q is not JSON: %v", bodies[0], err)
	}
	if got["token"] != "tok-normal" {
		t.Errorf("check body token = %q, want %q", got["token"], "tok-normal")
	}
}

func TestCheckTokenRequiresToken(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request to %s", r.URL.Path)
	}))
	defer ts.Close()

	cl := newTokenClient(t, ts.URL)
	if _, err := cl.CheckToken(context.Background(), "  "); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("CheckToken(empty) error = %v, want invalid_config", err)
	}
}

func TestEnsureTokenPollsAndCaches(t *testing.T) {
	t.Parallel()

	srv := &tokenServer{
		create: tokenResponse{Token: "tok-1", ExpiresAt: futureMillis(), Status: "PENDING"},
		check:  tokenResponse{Token: "tok-1", ExpiresAt: futureMillis(), Status: "NORMAL"},
	}
	ts := httptest.NewServer(srv.handler(t))
	defer ts.Close()

	cl := newTokenClient(t, ts.URL)
	cl.SetTokenPollInterval(time.Millisecond)
	cl.SetTokenPollTimeout(2 * time.Second)

	first, err := cl.EnsureToken(context.Background())
	if err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}
	if first.Status != client.TokenStatusNormal || first.Value != "tok-1" {
		t.Fatalf("EnsureToken() = %+v, want NORMAL tok-1", first)
	}
	if got := cl.AccessToken(); got != "tok-1" {
		t.Errorf("AccessToken() = %q, want %q", got, "tok-1")
	}

	// A second call must reuse the cached, valid token with no new requests.
	second, err := cl.EnsureToken(context.Background())
	if err != nil {
		t.Fatalf("second EnsureToken() error = %v", err)
	}
	if second.Value != first.Value {
		t.Errorf("second EnsureToken() value = %q, want cached %q", second.Value, first.Value)
	}

	create, check := srv.counts()
	if create != 1 || check != 1 {
		t.Fatalf("server calls create=%d check=%d, want 1/1 (cached token must be reused)", create, check)
	}
}

func TestEnsureTokenAttachesAccessTokenOnSignedRequest(t *testing.T) {
	t.Parallel()

	srv := &tokenServer{
		create: tokenResponse{Token: "tok-header", ExpiresAt: futureMillis(), Status: "PENDING"},
		check:  tokenResponse{Token: "tok-header", ExpiresAt: futureMillis(), Status: "NORMAL"},
	}
	ts := httptest.NewServer(srv.handler(t))
	defer ts.Close()

	cl := newTokenClient(t, ts.URL)
	cl.SetTokenPollInterval(time.Millisecond)
	if _, err := cl.EnsureToken(context.Background()); err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}

	var out []json.RawMessage
	if err := cl.Do(context.Background(), http.MethodGet, "/openapi/account/list", nil, &out); err != nil {
		t.Fatalf("Do(account/list) error = %v", err)
	}
	if got := srv.header(); got != "tok-header" {
		t.Fatalf("%s header = %q, want %q", client.AccessTokenHeader, got, "tok-header")
	}
}

func TestSetTokenInjectsHeaderWithoutNetwork(t *testing.T) {
	t.Parallel()

	srv := &tokenServer{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		verifyTokenSignature(t, r, body)
		srv.mu.Lock()
		srv.accessToken = r.Header.Get(client.AccessTokenHeader)
		srv.mu.Unlock()
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	cl := newTokenClient(t, ts.URL)
	cl.SetToken(&client.Token{
		Value:     "tok-set",
		Status:    client.TokenStatusNormal,
		ExpiresAt: time.Now().Add(time.Hour),
	})

	var out []json.RawMessage
	if err := cl.Do(context.Background(), http.MethodGet, "/openapi/account/list", nil, &out); err != nil {
		t.Fatalf("Do(account/list) error = %v", err)
	}
	if got := srv.header(); got != "tok-set" {
		t.Fatalf("%s header = %q, want %q", client.AccessTokenHeader, got, "tok-set")
	}
}

func TestEnsureTokenTerminalStatus(t *testing.T) {
	t.Parallel()

	for _, status := range []string{"INVALID", "EXPIRED"} {
		status := status
		t.Run(status, func(t *testing.T) {
			t.Parallel()

			srv := &tokenServer{
				create: tokenResponse{Token: "tok-bad", ExpiresAt: futureMillis(), Status: "PENDING"},
				check:  tokenResponse{Token: "tok-bad", ExpiresAt: futureMillis(), Status: status},
			}
			ts := httptest.NewServer(srv.handler(t))
			defer ts.Close()

			cl := newTokenClient(t, ts.URL)
			cl.SetTokenPollInterval(time.Millisecond)
			_, err := cl.EnsureToken(context.Background())
			if !errs.Is(err, errs.CodeAuth) {
				t.Fatalf("EnsureToken() error = %v, want auth error for %s", err, status)
			}
			if got := cl.AccessToken(); got != "" {
				t.Errorf("AccessToken() = %q, want empty after terminal status", got)
			}
		})
	}
}

func TestEnsureTokenSandboxShortCircuits(t *testing.T) {
	t.Parallel()

	srv := &tokenServer{
		create: tokenResponse{Token: "tok-sandbox", ExpiresAt: futureMillis(), Status: "PENDING"},
		check:  tokenResponse{Token: "tok-sandbox", ExpiresAt: futureMillis(), Status: "NORMAL"},
	}
	ts := httptest.NewServer(srv.handler(t))
	defer ts.Close()

	cl := newTokenClient(t, ts.URL, client.WithSandbox())
	_, err := cl.EnsureToken(context.Background())
	if !errs.Is(err, errs.CodeAuth) {
		t.Fatalf("EnsureToken() error = %v, want auth error (sandbox must not poll)", err)
	}
	_, check := srv.counts()
	if check != 0 {
		t.Fatalf("check calls = %d, want 0 in sandbox", check)
	}
}

func TestEnsureTokenTimeout(t *testing.T) {
	t.Parallel()

	srv := &tokenServer{
		create: tokenResponse{Token: "tok-slow", ExpiresAt: futureMillis(), Status: "PENDING"},
		check:  tokenResponse{Token: "tok-slow", ExpiresAt: futureMillis(), Status: "PENDING"},
	}
	ts := httptest.NewServer(srv.handler(t))
	defer ts.Close()

	cl := newTokenClient(t, ts.URL)
	cl.SetTokenPollInterval(time.Millisecond)
	cl.SetTokenPollTimeout(10 * time.Millisecond)

	_, err := cl.EnsureToken(context.Background())
	if !errs.Is(err, errs.CodeAuth) {
		t.Fatalf("EnsureToken() error = %v, want timeout auth error", err)
	}
}

func TestEnsureTokenRespectsContext(t *testing.T) {
	t.Parallel()

	srv := &tokenServer{
		create: tokenResponse{Token: "tok-ctx", ExpiresAt: futureMillis(), Status: "PENDING"},
		check:  tokenResponse{Token: "tok-ctx", ExpiresAt: futureMillis(), Status: "PENDING"},
	}
	ts := httptest.NewServer(srv.handler(t))
	defer ts.Close()

	cl := newTokenClient(t, ts.URL)
	cl.SetTokenPollInterval(time.Hour)
	cl.SetTokenPollTimeout(time.Hour)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := cl.EnsureToken(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("EnsureToken(cancelled) error = %v, want context.Canceled", err)
	}
}

func TestTokenStatusValid(t *testing.T) {
	t.Parallel()

	for _, s := range []client.TokenStatus{
		client.TokenStatusPending,
		client.TokenStatusNormal,
		client.TokenStatusInvalid,
		client.TokenStatusExpired,
	} {
		if !s.Valid() {
			t.Errorf("TokenStatus(%q).Valid() = false, want true", s)
		}
	}
	if client.TokenStatus("WEIRD").Valid() {
		t.Error(`TokenStatus("WEIRD").Valid() = true, want false`)
	}
}

// TestTokenStateIsolatedPerClient is a regression test for a correctness bug in
// which token state lived in package-level sync.Maps keyed by *Client. Those
// entries were never removed, so a client released by the garbage collector and
// a later client allocated at the same address could observe the same state: the
// new client's injection hook was never installed and it could send another
// client's (possibly stale) token. State now lives inside Client, so no
// address-keyed registry exists and clients are independent by construction.
func TestTokenStateIsolatedPerClient(t *testing.T) {
	t.Parallel()

	const clientCount = 8
	var mu sync.Mutex
	observed := make(map[string]string) // request path -> x-access-token seen
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		observed[r.URL.Path] = r.Header.Get(client.AccessTokenHeader)
		mu.Unlock()
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	clients := make([]*client.Client, 0, clientCount)
	for i := 0; i < clientCount; i++ {
		cl := newTokenClient(t, srv.URL)
		tok := fmt.Sprintf("tok-%d", i)
		cl.SetToken(&client.Token{
			Value:     tok,
			Status:    client.TokenStatusNormal,
			ExpiresAt: time.Now().Add(time.Hour),
		})
		if got := cl.AccessToken(); got != tok {
			t.Fatalf("client %d AccessToken() = %q, want %q", i, got, tok)
		}
		clients = append(clients, cl)
	}

	// Every client must inject only its own token, even though all clients share
	// the same server and were created in a tight loop (the address-reuse
	// scenario that triggered the original bug).
	for i, cl := range clients {
		path := fmt.Sprintf("/openapi/account/list/%d", i)
		if err := cl.Do(context.Background(), http.MethodGet, path, nil, nil); err != nil {
			t.Fatalf("client %d Do() error = %v", i, err)
		}
		mu.Lock()
		got := observed[path]
		mu.Unlock()
		want := fmt.Sprintf("tok-%d", i)
		if got != want {
			t.Fatalf("client %d injected %s = %q, want %q", i, client.AccessTokenHeader, got, want)
		}
	}
}

func TestClearTokenStopsInjection(t *testing.T) {
	t.Parallel()

	srv := &tokenServer{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.mu.Lock()
		srv.accessToken = r.Header.Get(client.AccessTokenHeader)
		srv.mu.Unlock()
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	cl := newTokenClient(t, ts.URL)
	cl.SetToken(&client.Token{Value: "tok-clear", Status: client.TokenStatusNormal, ExpiresAt: time.Now().Add(time.Hour)})
	cl.SetToken(nil)

	var out []json.RawMessage
	if err := cl.Do(context.Background(), http.MethodGet, "/openapi/account/list", nil, &out); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if got := srv.header(); got != "" {
		t.Fatalf("%s header = %q, want empty after clearing the token", client.AccessTokenHeader, got)
	}
}
