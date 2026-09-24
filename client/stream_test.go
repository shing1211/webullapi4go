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
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

// stubRateLimiter records the keys passed to Wait.
type stubRateLimiter struct {
	keys []string
}

func (s *stubRateLimiter) Wait(_ context.Context, key string) error {
	s.keys = append(s.keys, key)
	return nil
}

// stubBreaker is a controllable CircuitBreaker implementation.
type stubBreaker struct {
	allow    bool
	allowed  int
	success  int
	failures int
}

func (b *stubBreaker) Allow() bool {
	b.allowed++
	return b.allow
}

func (b *stubBreaker) RecordSuccess() { b.success++ }

func (b *stubBreaker) RecordFailure() { b.failures++ }

func TestDoStreamSignsRequest(t *testing.T) {
	t.Parallel()

	var received []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		received = body
		assertRequiredHeaders(t, r)
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/news/summaries/get"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}
		verifySignature(t, r, body)
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "data: {\"type\":\"text\"}\n\n")
	}))
	defer srv.Close()

	cl := newTestClient(t, srv.URL)

	resp, err := cl.DoStream(context.Background(), http.MethodPost, "/market-data/news/summaries/get", map[string]string{
		"lang": "<en>&more",
	})
	if err != nil {
		t.Fatalf("DoStream() error = %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	streamed, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading stream body: %v", err)
	}
	if !strings.Contains(string(streamed), "text") {
		t.Errorf("streamed body = %q, want it to contain the event", streamed)
	}

	const want = `{"lang":"<en>&more"}`
	if string(received) != want {
		t.Fatalf("sent body = %q, want %q", received, want)
	}
	if strings.Contains(string(received), `\u003c`) || strings.Contains(string(received), ": ") {
		t.Fatalf("sent body is not compact, unescaped JSON: %q", received)
	}
}

func TestDoStreamHonoursPathAPIVersion(t *testing.T) {
	t.Parallel()

	rec := newVersionRecorder()
	srv := httptest.NewServer(http.HandlerFunc(rec.handler))
	defer srv.Close()

	cl := newTestClient(t, srv.URL)

	cases := []struct {
		path string
		want string
	}{
		{"/market-data/news/summaries/get", client.APIVersionV2},
		{"/trading/orders/open-orders/list", client.APIVersionV3},
	}
	for _, tc := range cases {
		resp, err := cl.DoStream(context.Background(), http.MethodGet, tc.path, nil)
		if err != nil {
			t.Fatalf("DoStream(%s) error = %v", tc.path, err)
		}
		_ = resp.Body.Close()
		if got := rec.version(tc.path); got != tc.want {
			t.Errorf("DoStream(%s) x-version = %q, want %q", tc.path, got, tc.want)
		}
	}
}

func TestDoStreamAttachesCachedAccessToken(t *testing.T) {
	t.Parallel()

	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get(client.AccessTokenHeader)
		_, _ = io.WriteString(w, "data: {}\n\n")
	}))
	defer srv.Close()

	cl := newTestClient(t, srv.URL)
	cl.SetToken(&client.Token{
		Value:     "tok-stream",
		Status:    client.TokenStatusNormal,
		ExpiresAt: time.Now().Add(time.Hour),
	})

	resp, err := cl.DoStream(context.Background(), http.MethodGet, "/market-data/news/summaries/get", nil)
	if err != nil {
		t.Fatalf("DoStream() error = %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if got != "tok-stream" {
		t.Fatalf("%s = %q, want %q", client.AccessTokenHeader, got, "tok-stream")
	}
}

func TestDoStreamMapsHTTPStatusToTypedError(t *testing.T) {
	t.Parallel()

	cases := []struct {
		status   int
		code     errs.Code
		sentinel error
	}{
		{http.StatusUnauthorized, errs.CodeUnauthorized, errs.ErrUnauthorized},
		{http.StatusForbidden, errs.CodeForbidden, errs.ErrForbidden},
		{http.StatusTooManyRequests, errs.CodeRateLimited, errs.ErrRateLimited},
		{http.StatusInternalServerError, errs.CodeServer, errs.ErrServer},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(strconv.Itoa(tc.status), func(t *testing.T) {
			t.Parallel()

			status := tc.status
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(status)
				_, _ = w.Write([]byte(`{"message":"boom"}`))
			}))
			defer srv.Close()

			cl := newTestClient(t, srv.URL)

			resp, err := cl.DoStream(context.Background(), http.MethodGet, "/market-data/news/summaries/get", nil)
			if resp != nil {
				t.Errorf("resp = %v, want nil on error", resp)
			}
			if err == nil {
				t.Fatalf("DoStream() error = nil, want status %d error", status)
			}
			var e *errs.Error
			if !errors.As(err, &e) {
				t.Fatalf("DoStream() error type = %T, want *errs.Error", err)
			}
			if e.Code != tc.code || e.Status != status {
				t.Errorf("error = %+v, want code %q status %d", e, tc.code, status)
			}
			if !strings.Contains(e.Message, "boom") {
				t.Errorf("message = %q, want it to contain the API message", e.Message)
			}
			if !errors.Is(err, tc.sentinel) {
				t.Errorf("errors.Is(err, %v) = false, want true", tc.sentinel)
			}
		})
	}
}

func TestDoStreamDoesNotRetry(t *testing.T) {
	t.Parallel()

	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"boom"}`))
	}))
	defer srv.Close()

	cl := newTestClient(t, srv.URL)

	_, err := cl.DoStream(context.Background(), http.MethodGet, "/market-data/news/summaries/get", nil)
	if err == nil {
		t.Fatal("DoStream() error = nil, want a server error")
	}
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Fatalf("server received %d requests, want 1 (a stream must not be retried)", got)
	}
}

func TestDoStreamAppliesRateLimiterAndBreaker(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "data: {}\n\n")
	}))
	defer srv.Close()

	limiter := &stubRateLimiter{}
	br := &stubBreaker{allow: true}
	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithBaseURL(srv.URL),
		client.WithRateLimiter(limiter),
		client.WithBreaker(br),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	resp, err := cl.DoStream(context.Background(), http.MethodPost, "/market-data/news/summaries/get", nil)
	if err != nil {
		t.Fatalf("DoStream() error = %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if len(limiter.keys) != 1 || limiter.keys[0] != "/market-data/news/summaries/get" {
		t.Errorf("rate limiter keys = %v, want the request path", limiter.keys)
	}
	if br.allowed != 1 || br.success != 1 || br.failures != 0 {
		t.Errorf("breaker allowed=%d success=%d failures=%d, want 1/1/0", br.allowed, br.success, br.failures)
	}
}

func TestDoStreamBreakerOpenSkipsRequest(t *testing.T) {
	t.Parallel()

	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&hits, 1)
		_, _ = io.WriteString(w, "data: {}\n\n")
	}))
	defer srv.Close()

	br := &stubBreaker{allow: false}
	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithBaseURL(srv.URL),
		client.WithBreaker(br),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	_, err = cl.DoStream(context.Background(), http.MethodGet, "/market-data/news/summaries/get", nil)
	if !errors.Is(err, client.ErrCircuitOpen) {
		t.Fatalf("DoStream() error = %v, want ErrCircuitOpen", err)
	}
	if got := atomic.LoadInt32(&hits); got != 0 {
		t.Fatalf("server received %d requests, want 0 when the breaker is open", got)
	}
	if br.allowed != 1 || br.success != 0 || br.failures != 0 {
		t.Errorf("breaker allowed=%d success=%d failures=%d, want 1/0/0", br.allowed, br.success, br.failures)
	}
}
