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
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/internal/auth"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

const (
	testAppKey    = "test-app-key"
	testAppSecret = "test-app-secret"
)

// requiredHeaders are the headers every signed request must carry.
var requiredHeaders = []string{
	"x-app-key",
	"x-timestamp",
	"x-signature",
	"x-signature-algorithm",
	"x-signature-version",
	"x-signature-nonce",
	"x-version",
}

// assertRequiredHeaders verifies the presence and shape of the signing headers.
// It uses t.Errorf rather than t.Fatalf because it runs on the server goroutine.
func assertRequiredHeaders(t *testing.T, r *http.Request) {
	t.Helper()
	for _, name := range requiredHeaders {
		if r.Header.Get(name) == "" {
			t.Errorf("missing required header %q", name)
		}
	}
	if got := r.Header.Get("x-app-key"); got != testAppKey {
		t.Errorf("x-app-key = %q, want %q", got, testAppKey)
	}
	if got := r.Header.Get("x-signature-algorithm"); got != auth.SignatureAlgorithm {
		t.Errorf("x-signature-algorithm = %q, want %q", got, auth.SignatureAlgorithm)
	}
	if got := r.Header.Get("x-signature-version"); got != auth.SignatureVersion {
		t.Errorf("x-signature-version = %q, want %q", got, auth.SignatureVersion)
	}
	if got := r.Header.Get("x-version"); got != "v2" {
		t.Errorf("x-version = %q, want %q", got, "v2")
	}
	if nonce := r.Header.Get("x-signature-nonce"); len(nonce) != 32 {
		t.Errorf("x-signature-nonce = %q, want 32 hex characters", nonce)
	}
	if _, err := time.Parse(auth.TimestampFormat, r.Header.Get("x-timestamp")); err != nil {
		t.Errorf("x-timestamp %q is not in %q format: %v", r.Header.Get("x-timestamp"), auth.TimestampFormat, err)
	}
}

// verifySignature recomputes the request signature from the received request
// and compares it to the x-signature header.
func verifySignature(t *testing.T, r *http.Request, body []byte) {
	t.Helper()
	headers := make(http.Header)
	headers.Set(auth.HeaderAppKey, r.Header.Get(auth.HeaderAppKey))
	headers.Set(auth.HeaderSignatureAlgorithm, r.Header.Get(auth.HeaderSignatureAlgorithm))
	headers.Set(auth.HeaderSignatureVersion, r.Header.Get(auth.HeaderSignatureVersion))
	headers.Set(auth.HeaderSignatureNonce, r.Header.Get(auth.HeaderSignatureNonce))
	headers.Set(auth.HeaderTimestamp, r.Header.Get(auth.HeaderTimestamp))
	// net/http moves the Host header into Request.Host on the server side.
	headers.Set(auth.HeaderHost, r.Host)

	got, err := auth.Sign(auth.SignParams{
		Method:    r.Method,
		Path:      r.URL.Path,
		Query:     r.URL.Query(),
		Headers:   headers,
		Body:      body,
		AppSecret: testAppSecret,
	})
	if err != nil {
		t.Errorf("recomputing signature: %v", err)
		return
	}
	if want := r.Header.Get("x-signature"); got != want {
		t.Errorf("signature mismatch: recomputed %q, header %q", got, want)
	}
}

func newTestClient(t *testing.T, baseURL string) *client.Client {
	t.Helper()
	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithBaseURL(baseURL),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() {
		if err := cl.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	})
	return cl
}

func TestDoSignsGETRequest(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		assertRequiredHeaders(t, r)
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/openapi/account/list" {
			t.Errorf("path = %q, want /openapi/account/list", r.URL.Path)
		}
		if len(body) != 0 {
			t.Errorf("GET body = %q, want empty", body)
		}
		verifySignature(t, r, body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"acct-1"}]`))
	}))
	defer srv.Close()

	cl := newTestClient(t, srv.URL)

	var out []json.RawMessage
	if err := cl.Do(context.Background(), http.MethodGet, "/openapi/account/list", nil, &out); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("decoded %d records, want 1", len(out))
	}
}

func TestDoSignsQueryString(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		assertRequiredHeaders(t, r)
		if got, want := r.URL.RawQuery, "a=1&b=2"; got != want {
			t.Errorf("RawQuery = %q, want %q", got, want)
		}
		verifySignature(t, r, body)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	cl := newTestClient(t, srv.URL)

	var out json.RawMessage
	if err := cl.Do(context.Background(), http.MethodGet, "/openapi/account/list?b=2&a=1", nil, &out); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
}

func TestDoSignsPOSTCompactBody(t *testing.T) {
	t.Parallel()

	var received []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		received = body
		assertRequiredHeaders(t, r)
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}
		verifySignature(t, r, body)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	cl := newTestClient(t, srv.URL)

	var out json.RawMessage
	err := cl.Do(context.Background(), http.MethodPost, "/openapi/token/create", map[string]string{
		"name": "<a>&b",
	}, &out)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	const want = `{"name":"<a>&b"}`
	if string(received) != want {
		t.Fatalf("sent body = %q, want %q", received, want)
	}
	if strings.Contains(string(received), `\u003c`) || strings.Contains(string(received), `\u0026`) {
		t.Fatalf("sent body is HTML-escaped: %q", received)
	}
	if strings.Contains(string(received), ": ") || strings.Contains(string(received), ", ") {
		t.Fatalf("sent body is not compact: %q", received)
	}
}

func TestDoMapsHTTPStatusToTypedError(t *testing.T) {
	t.Parallel()

	cases := []struct {
		status   int
		code     errs.Code
		sentinel error
	}{
		{http.StatusUnauthorized, errs.CodeUnauthorized, errs.ErrUnauthorized},
		{http.StatusForbidden, errs.CodeForbidden, errs.ErrForbidden},
		{http.StatusExpectationFailed, errs.CodeInvalidToken, errs.ErrInvalidToken},
		{http.StatusTooManyRequests, errs.CodeRateLimited, errs.ErrRateLimited},
		{http.StatusInternalServerError, errs.CodeServer, errs.ErrServer},
		{http.StatusBadGateway, errs.CodeServer, errs.ErrServer},
		{http.StatusNotFound, errs.CodeAPI, nil},
	}
	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.status), func(t *testing.T) {
			t.Parallel()

			status := tc.status
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				_, _ = w.Write([]byte(`{"message":"boom"}`))
			}))
			defer srv.Close()

			cl := newTestClient(t, srv.URL)

			err := cl.Do(context.Background(), http.MethodGet, "/openapi/account/list", nil, nil)
			if err == nil {
				t.Fatalf("Do() error = nil, want status %d error", status)
			}

			var e *errs.Error
			if !errors.As(err, &e) {
				t.Fatalf("Do() error type = %T, want *errs.Error", err)
			}
			if e.Code != tc.code {
				t.Errorf("error code = %q, want %q", e.Code, tc.code)
			}
			if e.Status != status {
				t.Errorf("error status = %d, want %d", e.Status, status)
			}
			if !strings.Contains(e.Message, "boom") {
				t.Errorf("error message = %q, want it to contain the API message %q", e.Message, "boom")
			}
			if tc.sentinel != nil && !errors.Is(err, tc.sentinel) {
				t.Errorf("errors.Is(err, %v) = false, want true", tc.sentinel)
			}
		})
	}
}

func TestNewResolvesEndpointsFromRegionAndEnvironment(t *testing.T) {
	t.Parallel()

	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithRegion(client.HK),
		client.WithEnvironment(client.Sandbox),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	if got := cl.Endpoints().HTTP; got != "https://api.sandbox.webull.hk" {
		t.Fatalf("sandbox HK endpoint = %q, want %q", got, "https://api.sandbox.webull.hk")
	}
	if cl.Region() != client.HK || cl.Environment() != client.Sandbox {
		t.Fatalf("Region/Environment = %q/%q, want hk/sandbox", cl.Region(), cl.Environment())
	}

	us, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithRegion(client.US),
		client.WithEnvironment(client.Production),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	if got := us.Endpoints().HTTP; got != "https://api.webull.com" {
		t.Fatalf("US production endpoint = %q, want %q", got, "https://api.webull.com")
	}
}

func TestNewEndpointOverrides(t *testing.T) {
	t.Parallel()

	base, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithRegion(client.HK),
		client.WithEnvironment(client.Sandbox),
		client.WithBaseURL("http://127.0.0.1:9090"),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	if got := base.Endpoints().HTTP; got != "http://127.0.0.1:9090" {
		t.Fatalf("WithBaseURL endpoint = %q, want %q", got, "http://127.0.0.1:9090")
	}

	full := client.Endpoints{HTTP: "http://example.test", MQTT: "mqtt.example.test:1883"}
	ep, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithRegion(client.US),
		client.WithEnvironment(client.Production),
		client.WithEndpoints(full),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	if got := ep.Endpoints(); got != full {
		t.Fatalf("WithEndpoints = %+v, want %+v", got, full)
	}
}

func TestNewValidatesConfig(t *testing.T) {
	t.Parallel()

	if _, err := client.New(client.WithAppSecret(testAppSecret)); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("New() missing app key error = %v, want invalid_config", err)
	}
	if _, err := client.New(client.WithAppKey(testAppKey)); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("New() missing app secret error = %v, want invalid_config", err)
	}
	if _, err := client.New(client.WithAppKey(testAppKey), client.WithAppSecret(testAppSecret), client.WithBaseURL("not-a-url")); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("New() invalid endpoint error = %v, want invalid_config", err)
	}
}

func TestClockDriftCorrection(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	var capturedTimestamp string
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		capturedTimestamp = r.Header.Get(auth.HeaderTimestamp)
		w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"tok","token_type":"Bearer","expires_in":86400}`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithBaseURL(srv.URL+"/"),
		client.WithClockDriftCorrection(true),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	ctx := context.Background()

	err = cl.Do(ctx, http.MethodGet, "/test", nil, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if capturedTimestamp == "" {
		t.Fatal("timestamp header was not captured")
	}
	ts, err := time.Parse(auth.TimestampFormat, capturedTimestamp)
	if err != nil {
		t.Fatalf("time.Parse(%q) error = %v", capturedTimestamp, err)
	}

	now := time.Now()
	drift := ts.Sub(now)
	if drift < -time.Minute || drift > time.Minute {
		t.Fatalf("captured timestamp drift = %v, want roughly 0 (within 1 minute)", drift)
	}
}

func TestWithHTTPTransport(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"tok","token_type":"Bearer","expires_in":86400}`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	customTransport := &http.Transport{
		MaxIdleConns:        99,
		MaxIdleConnsPerHost: 50,
	}

	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithBaseURL(srv.URL+"/"),
		client.WithHTTPTransport(customTransport),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	cfg := cl.Config()
	if got := cfg.HTTPClient.Transport.(*http.Transport).MaxIdleConnsPerHost; got != 50 {
		t.Errorf("MaxIdleConnsPerHost = %d, want 50", got)
	}

	ctx := context.Background()
	err = cl.Do(ctx, http.MethodGet, "/test", nil, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
}

func TestWithResiliencePreset(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"tok","token_type":"Bearer","expires_in":86400}`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithBaseURL(srv.URL+"/"),
		client.WithResiliencePreset(client.ProductionPreset),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	ctx := context.Background()
	err = cl.Do(ctx, http.MethodGet, "/test", nil, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
}

func TestInterceptorChain(t *testing.T) {
	t.Parallel()

	var callOrder []string

	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		callOrder = append(callOrder, "server")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"tok","token_type":"Bearer","expires_in":86400}`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithBaseURL(srv.URL+"/"),
		client.WithInterceptor(func(ctx context.Context, next func(context.Context) error) error {
			callOrder = append(callOrder, "interceptor-a")
			return next(ctx)
		}),
		client.WithInterceptor(func(ctx context.Context, next func(context.Context) error) error {
			callOrder = append(callOrder, "interceptor-b")
			return next(ctx)
		}),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	err = cl.Do(context.Background(), http.MethodGet, "/test", nil, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if len(callOrder) != 3 {
		t.Fatalf("callOrder = %v, want [interceptor-a, interceptor-b, server]", callOrder)
	}
	if callOrder[0] != "interceptor-a" || callOrder[1] != "interceptor-b" || callOrder[2] != "server" {
		t.Errorf("callOrder = %v, want [interceptor-a, interceptor-b, server]", callOrder)
	}
}

func TestHooks(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"tok","token_type":"Bearer","expires_in":86400}`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	var hookOrder []string
	hooks := client.Hooks{
		OnRequest: func(method, path string) {
			hookOrder = append(hookOrder, "OnRequest:"+path)
		},
		OnLatency: func(attempt int, d time.Duration) {
			hookOrder = append(hookOrder, "OnLatency")
		},
		OnResponse: func(status int, latency time.Duration) {
			hookOrder = append(hookOrder, "OnResponse")
		},
	}

	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithBaseURL(srv.URL+"/"),
		client.WithHooks(hooks),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	err = cl.Do(context.Background(), http.MethodGet, "/test", nil, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if len(hookOrder) == 0 {
		t.Fatal("no hooks were called")
	}
	if hookOrder[0] != "OnRequest:/test" {
		t.Errorf("first hook = %q, want OnRequest:/test", hookOrder[0])
	}
}
