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
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/internal/errs"
)

// versionRecorder is an httptest handler that records the x-version header keyed
// by request path.
type versionRecorder struct {
	mu       sync.Mutex
	versions map[string]string
}

func newVersionRecorder() *versionRecorder {
	return &versionRecorder{versions: make(map[string]string)}
}

func (r *versionRecorder) handler(w http.ResponseWriter, req *http.Request) {
	r.mu.Lock()
	r.versions[req.URL.Path] = req.Header.Get("x-version")
	r.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{}`))
}

func (r *versionRecorder) version(path string) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.versions[path]
}

func TestDoUsesDefaultAPIVersionV2(t *testing.T) {
	t.Parallel()

	rec := newVersionRecorder()
	srv := httptest.NewServer(http.HandlerFunc(rec.handler))
	defer srv.Close()

	cl := newTestClient(t, srv.URL)
	if err := cl.Do(context.Background(), http.MethodGet, "/openapi/account/list", nil, nil); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if got := rec.version("/openapi/account/list"); got != client.APIVersionV2 {
		t.Fatalf("x-version = %q, want default %q", got, client.APIVersionV2)
	}
}

func TestDoUsesConfiguredAPIVersion(t *testing.T) {
	t.Parallel()

	rec := newVersionRecorder()
	srv := httptest.NewServer(http.HandlerFunc(rec.handler))
	defer srv.Close()

	cl := newTokenClient(t, srv.URL, client.WithAPIVersion(client.APIVersionV3))
	if err := cl.Do(context.Background(), http.MethodGet, "/market-data/fundamentals/company-profiles/get", nil, nil); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if got := rec.version("/market-data/fundamentals/company-profiles/get"); got != client.APIVersionV3 {
		t.Fatalf("x-version = %q, want %q", got, client.APIVersionV3)
	}
}

func TestWithAPIVersionForOverridesByLongestPrefix(t *testing.T) {
	t.Parallel()

	rec := newVersionRecorder()
	srv := httptest.NewServer(http.HandlerFunc(rec.handler))
	defer srv.Close()

	cl := newTokenClient(t, srv.URL,
		client.WithAPIVersion(client.APIVersionV2),
		client.WithAPIVersionFor("/market-data", client.APIVersionV3),
		client.WithAPIVersionFor("/market-data/fundamentals", client.APIVersionV2),
	)

	cases := []struct {
		path string
		want string
	}{
		{"/market-data/fundamentals/company-profiles/get", client.APIVersionV2},
		{"/market-data/quotes/list", client.APIVersionV3},
		{"/openapi/account/list", client.APIVersionV2},
	}
	for _, tc := range cases {
		if err := cl.Do(context.Background(), http.MethodGet, tc.path, nil, nil); err != nil {
			t.Fatalf("Do(%s) error = %v", tc.path, err)
		}
		if got := rec.version(tc.path); got != tc.want {
			t.Errorf("Do(%s) x-version = %q, want %q", tc.path, got, tc.want)
		}
	}
}

func TestNewRejectsInvalidAPIVersion(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		opt  client.Option
	}{
		{"unknown version", client.WithAPIVersion("v9")},
		{"empty version", client.WithAPIVersion("")},
		{"unknown override", client.WithAPIVersionFor("/market-data", "v9")},
		{"empty override prefix", client.WithAPIVersionFor("", client.APIVersionV3)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := client.New(
				client.WithAppKey(testAppKey),
				client.WithAppSecret(testAppSecret),
				tc.opt,
			)
			if !errs.Is(err, errs.CodeInvalidConfig) {
				t.Fatalf("New(%s) error = %v, want invalid_config", tc.name, err)
			}
		})
	}
}

func TestAutoTokenSandboxEnsuresAndAttachesToken(t *testing.T) {
	t.Parallel()

	srv := &tokenServer{
		create: tokenResponse{Token: "tok-auto", ExpiresAt: futureMillis(), Status: "NORMAL"},
	}
	ts := httptest.NewServer(srv.handler(t))
	defer ts.Close()

	cl := newTokenClient(t, ts.URL, client.WithSandbox(), client.WithAutoToken(true))

	if err := cl.Do(context.Background(), http.MethodGet, "/openapi/account/list", nil, nil); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if got := srv.header(); got != "tok-auto" {
		t.Fatalf("%s header = %q, want %q", client.AccessTokenHeader, got, "tok-auto")
	}
	if got := cl.AccessToken(); got != "tok-auto" {
		t.Fatalf("AccessToken() = %q, want %q", got, "tok-auto")
	}

	// A second request must reuse the cached token, not create another one.
	if err := cl.Do(context.Background(), http.MethodGet, "/openapi/account/list", nil, nil); err != nil {
		t.Fatalf("second Do() error = %v", err)
	}
	if create, check := srv.counts(); create != 1 || check != 0 {
		t.Fatalf("server calls create=%d check=%d, want 1/0 (token must be cached)", create, check)
	}
}

// TestAutoTokenConcurrentCallsSingleFlight verifies that concurrent first
// requests on one client serialise automatic token acquisition through the
// client's own mutex: exactly one token is created and every request carries it.
func TestAutoTokenConcurrentCallsSingleFlight(t *testing.T) {
	t.Parallel()

	srv := &tokenServer{
		create: tokenResponse{Token: "tok-concurrent", ExpiresAt: futureMillis(), Status: "NORMAL"},
	}
	ts := httptest.NewServer(srv.handler(t))
	defer ts.Close()

	cl := newTokenClient(t, ts.URL, client.WithSandbox(), client.WithAutoToken(true))

	const goroutines = 16
	var wg sync.WaitGroup
	errCh := make(chan error, goroutines)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := cl.Do(context.Background(), http.MethodGet, "/openapi/account/list", nil, nil); err != nil {
				errCh <- err
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatalf("concurrent Do() error = %v", err)
	}

	if create, check := srv.counts(); create != 1 || check != 0 {
		t.Fatalf("server calls create=%d check=%d, want 1/0 (first-token acquisition must be single-flight)", create, check)
	}
	if got := srv.header(); got != "tok-concurrent" {
		t.Fatalf("%s header = %q, want %q", client.AccessTokenHeader, got, "tok-concurrent")
	}
}

func TestAutoTokenProductionReturnsTypedErrorWithoutNetwork(t *testing.T) {
	t.Parallel()

	var hits int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		_, _ = w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	cl := newTokenClient(t, ts.URL,
		client.WithEnvironment(client.Production),
		client.WithAutoToken(true),
	)

	err := cl.Do(context.Background(), http.MethodGet, "/market-data/fundamentals/company-profiles/get", nil, nil)
	if !errs.Is(err, errs.CodeAuth) {
		t.Fatalf("Do() error = %v, want auth error", err)
	}
	if !errors.Is(err, client.ErrAccessTokenRequired) {
		t.Fatalf("errors.Is(err, ErrAccessTokenRequired) = false, want true")
	}
	if !strings.Contains(err.Error(), "EnsureToken") {
		t.Fatalf("error %q does not mention EnsureToken", err)
	}
	if hits != 0 {
		t.Fatalf("server received %d requests, want 0 (production must not auto-start 2FA)", hits)
	}
}

func TestAutoTokenDisabledByDefault(t *testing.T) {
	t.Parallel()

	srv := &tokenServer{}
	ts := httptest.NewServer(srv.handler(t))
	defer ts.Close()

	cl := newTokenClient(t, ts.URL, client.WithSandbox())

	if err := cl.Do(context.Background(), http.MethodGet, "/openapi/account/list", nil, nil); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if got := srv.header(); got != "" {
		t.Fatalf("%s header = %q, want empty when auto-token is off", client.AccessTokenHeader, got)
	}
	if create, check := srv.counts(); create != 0 || check != 0 {
		t.Fatalf("server calls create=%d check=%d, want 0/0", create, check)
	}
}

func TestAutoTokenUsesCachedTokenWithoutNetwork(t *testing.T) {
	t.Parallel()

	srv := &tokenServer{}
	ts := httptest.NewServer(srv.handler(t))
	defer ts.Close()

	cl := newTokenClient(t, ts.URL, client.WithSandbox(), client.WithAutoToken(true))
	cl.SetToken(&client.Token{
		Value:     "tok-cached",
		Status:    client.TokenStatusNormal,
		ExpiresAt: time.Now().Add(time.Hour),
	})

	if err := cl.Do(context.Background(), http.MethodGet, "/openapi/account/list", nil, nil); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if got := srv.header(); got != "tok-cached" {
		t.Fatalf("%s header = %q, want %q", client.AccessTokenHeader, got, "tok-cached")
	}
	if create, check := srv.counts(); create != 0 || check != 0 {
		t.Fatalf("server calls create=%d check=%d, want 0/0 (cached token must skip network)", create, check)
	}
}
