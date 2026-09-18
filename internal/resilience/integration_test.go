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

package resilience_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
)

// newClient builds a client pointed at baseURL. It lives here rather than in
// client's own tests so that the resilience integration can be exercised
// without modifying the client package's test files.
func newClient(t *testing.T, baseURL string, opts ...client.Option) *client.Client {
	t.Helper()
	base := []client.Option{
		client.WithAppKey("test-app-key"),
		client.WithAppSecret("test-app-secret"),
		client.WithBaseURL(baseURL),
	}
	cl, err := client.New(append(base, opts...)...)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })
	return cl
}

func TestDefaultRetryRetriesTransientGET(t *testing.T) {
	t.Parallel()

	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"boom"}`))
			return
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	cl := newClient(t, srv.URL, client.WithRetry(client.RetryConfig{
		BaseDelay: time.Millisecond,
		Jitter:    false,
	}))

	var out map[string]any
	if err := cl.Do(context.Background(), http.MethodGet, "/openapi/x", nil, &out); err != nil {
		t.Fatalf("Do() error = %v, want nil after retry", err)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("server calls = %d, want 2", got)
	}
	if ok, _ := out["ok"].(bool); !ok {
		t.Fatalf("decoded response = %v, want ok=true", out)
	}
}

func TestNonIdempotentNotRetriedByDefault(t *testing.T) {
	t.Parallel()

	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"boom"}`))
	}))
	defer srv.Close()

	cl := newClient(t, srv.URL, client.WithRetry(client.RetryConfig{
		BaseDelay: time.Millisecond,
		Jitter:    false,
	}))

	if err := cl.Do(context.Background(), http.MethodPost, "/openapi/x", map[string]string{"a": "b"}, nil); err == nil {
		t.Fatal("Do() error = nil, want server error")
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("server calls = %d, want 1 (POST must not be retried)", got)
	}
}

func TestWithoutRetryAttemptsOnce(t *testing.T) {
	t.Parallel()

	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"boom"}`))
	}))
	defer srv.Close()

	cl := newClient(t, srv.URL, client.WithoutRetry())
	if err := cl.Do(context.Background(), http.MethodGet, "/openapi/x", nil, nil); err == nil {
		t.Fatal("Do() error = nil, want server error")
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("server calls = %d, want 1", got)
	}
}

func TestBreakerRejectsWhenOpen(t *testing.T) {
	t.Parallel()

	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"boom"}`))
	}))
	defer srv.Close()

	cl := newClient(t, srv.URL,
		client.WithoutRetry(),
		client.WithBreaker(client.NewBreaker(1, time.Hour)),
	)

	if err := cl.Do(context.Background(), http.MethodGet, "/openapi/x", nil, nil); err == nil {
		t.Fatal("first Do() error = nil, want server error")
	}
	err := cl.Do(context.Background(), http.MethodGet, "/openapi/x", nil, nil)
	if !errors.Is(err, client.ErrCircuitOpen) {
		t.Fatalf("second Do() error = %v, want ErrCircuitOpen", err)
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("server calls = %d, want 1 (second call rejected by breaker)", got)
	}
}

// countingLimiter records the keys it is asked to wait on.
type countingLimiter struct {
	mu   sync.Mutex
	keys []string
	err  error
}

func (l *countingLimiter) Wait(_ context.Context, key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.keys = append(l.keys, key)
	return l.err
}

func TestRateLimiterSeesRequestPath(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	lim := &countingLimiter{}
	cl := newClient(t, srv.URL, client.WithRateLimiter(lim))

	if err := cl.Do(context.Background(), http.MethodGet, "/openapi/account/list", nil, nil); err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	lim.mu.Lock()
	defer lim.mu.Unlock()
	if len(lim.keys) != 1 || lim.keys[0] != "/openapi/account/list" {
		t.Fatalf("limiter keys = %v, want [/openapi/account/list]", lim.keys)
	}
}

func TestRateLimiterErrorStopsRequest(t *testing.T) {
	t.Parallel()

	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	lim := &countingLimiter{err: context.DeadlineExceeded}
	cl := newClient(t, srv.URL, client.WithoutRetry(), client.WithRateLimiter(lim))

	err := cl.Do(context.Background(), http.MethodGet, "/openapi/x", nil, nil)
	if err == nil {
		t.Fatal("Do() error = nil, want limiter error")
	}
	if got := atomic.LoadInt32(&calls); got != 0 {
		t.Fatalf("server calls = %d, want 0", got)
	}
}
