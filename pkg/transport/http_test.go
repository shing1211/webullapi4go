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

package transport_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/pkg/transport"
)

var _ transport.Doer = (*transport.Client)(nil)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

type requestContextKey struct{}

type closeTrackingRoundTripper struct {
	roundTrip  roundTripFunc
	closeCalls atomic.Int32
}

func (t *closeTrackingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.roundTrip(req)
}

func (t *closeTrackingRoundTripper) CloseIdleConnections() { t.closeCalls.Add(1) }

func TestNewValidation(t *testing.T) {
	t.Parallel()

	if _, err := transport.New("https://example.test", nil, ""); !errors.Is(err, transport.ErrNilHTTPClient) {
		t.Fatalf("New() with nil HTTP client error = %v, want ErrNilHTTPClient", err)
	}

	tests := []struct {
		name     string
		baseURL  string
		wantBase string
		wantErr  bool
	}{
		{name: "trimmed HTTPS URL", baseURL: "  https://example.test/root/  ", wantBase: "https://example.test/root/"},
		{name: "HTTP URL", baseURL: "http://example.test", wantBase: "http://example.test"},
		{name: "empty", baseURL: "", wantErr: true},
		{name: "relative", baseURL: "/relative", wantErr: true},
		{name: "missing host", baseURL: "https://", wantErr: true},
		{name: "invalid escape", baseURL: "https://example.test/%zz", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hc := &http.Client{}
			got, err := transport.New(tt.baseURL, hc, "test-agent")
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if got != nil {
					t.Fatal("New() returned a client with an invalid base URL")
				}
				return
			}
			if got.HTTPClient() != hc {
				t.Fatal("HTTPClient() did not return the supplied client")
			}
			if base := got.BaseURL().String(); base != tt.wantBase {
				t.Fatalf("BaseURL() = %q, want %q", base, tt.wantBase)
			}
		})
	}
}

func TestNewRequestTable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     string
		path       string
		query      url.Values
		body       []byte
		userAgent  string
		wantURL    string
		wantBody   string
		wantLength int64
	}{
		{
			name:      "relative path and encoded query",
			method:    http.MethodGet,
			path:      "account/list",
			query:     url.Values{"symbol": {"AAPL", "MSFT"}, "account": {"id 1"}},
			userAgent: "webull-test/1.0",
			wantURL:   "https://example.test/root/account/list?account=id+1&symbol=AAPL&symbol=MSFT",
		},
		{
			name:       "rooted path",
			method:     http.MethodPost,
			path:       "/account",
			body:       []byte(`{"accountId":"1"}`),
			wantURL:    "https://example.test/root/account",
			wantBody:   `{"accountId":"1"}`,
			wantLength: 17,
		},
		{
			name:    "path query",
			method:  http.MethodHead,
			path:    "account?active=true",
			wantURL: "https://example.test/root/account?active=true",
		},
		{
			name:    "explicit query overrides path query",
			method:  http.MethodGet,
			path:    "account?ignored=true",
			query:   url.Values{},
			wantURL: "https://example.test/root/account",
		},
		{
			name:    "empty path preserves base",
			method:  http.MethodGet,
			path:    "",
			query:   url.Values{"page": {"1"}},
			wantURL: "https://example.test/root/?page=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, err := transport.New(" https://example.test/root/ ", &http.Client{}, tt.userAgent)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			ctx := context.WithValue(context.Background(), requestContextKey{}, "request")
			req, err := c.NewRequest(ctx, tt.method, tt.path, tt.query, tt.body)
			if err != nil {
				t.Fatalf("NewRequest() error = %v", err)
			}
			if got := req.URL.String(); got != tt.wantURL {
				t.Fatalf("request URL = %q, want %q", got, tt.wantURL)
			}
			if req.Method != tt.method {
				t.Fatalf("request method = %q, want %q", req.Method, tt.method)
			}
			if req.Context() != ctx {
				t.Fatal("request did not retain the supplied context")
			}
			if got := req.Header.Get("User-Agent"); got != tt.userAgent {
				t.Fatalf("User-Agent = %q, want %q", got, tt.userAgent)
			}
			if req.ContentLength != tt.wantLength {
				t.Fatalf("ContentLength = %d, want %d", req.ContentLength, tt.wantLength)
			}
			if tt.body == nil {
				if req.Body != nil {
					t.Fatal("request has a body when none was supplied")
				}
				if req.GetBody != nil {
					t.Fatal("GetBody is set for an empty request")
				}
				return
			}
			gotBody, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll(request body) error = %v", err)
			}
			if string(gotBody) != tt.wantBody {
				t.Fatalf("request body = %q, want %q", gotBody, tt.wantBody)
			}
			if req.GetBody == nil {
				t.Fatal("GetBody is nil for a byte-backed request body")
			}
		})
	}
}

func TestNewRequestRejectsInvalidInputs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     string
		path       string
		nilContext bool
	}{
		{name: "invalid URL escape", method: http.MethodGet, path: "%zz"},
		{name: "invalid method", method: "bad method", path: "account"},
		{name: "nil context", method: http.MethodGet, path: "account", nilContext: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, err := transport.New("https://example.test", &http.Client{}, "")
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			ctx := context.Background()
			if tt.nilContext {
				ctx = nil
			}
			if got, err := c.NewRequest(ctx, tt.method, tt.path, nil, nil); err == nil || got != nil {
				t.Fatalf("NewRequest() = %v, %v; want nil request and error", got, err)
			}
		})
	}
}

func TestAccessorsReturnDefensiveValues(t *testing.T) {
	t.Parallel()

	hc := &http.Client{}
	c, err := transport.New("https://example.test/root", hc, "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if c.HTTPClient() != hc {
		t.Fatal("HTTPClient() did not return the configured client")
	}
	base := c.BaseURL()
	base.Path = "/changed"
	if got := c.BaseURL().String(); got != "https://example.test/root" {
		t.Fatalf("mutating BaseURL() result changed client to %q", got)
	}
}

func TestDoDelegatesToHTTPClient(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("round trip failed")
	tests := []struct {
		name         string
		response     *http.Response
		err          error
		wantResponse bool
		wantError    bool
	}{
		{
			name: "response",
			response: &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("ok")),
			},
			wantResponse: true,
		},
		{name: "error", err: sentinel, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var received *http.Request
			hc := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				received = req
				return tt.response, tt.err
			})}
			c, err := transport.New("https://example.test", hc, "")
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			req, err := c.NewRequest(context.Background(), http.MethodGet, "account", nil, nil)
			if err != nil {
				t.Fatalf("NewRequest() error = %v", err)
			}
			got, err := c.Do(req)
			if received != req {
				t.Fatal("Do() did not pass the original request to the HTTP client")
			}
			if (err != nil) != tt.wantError {
				t.Fatalf("Do() error = %v, wantError = %v", err, tt.wantError)
			}
			if tt.wantError && !errors.Is(err, sentinel) {
				t.Fatalf("Do() error = %v, want wrapped sentinel", err)
			}
			if tt.wantResponse {
				if got != tt.response {
					t.Fatalf("Do() response = %p, want delegated response %p", got, tt.response)
				}
				if err := got.Body.Close(); err != nil {
					t.Fatalf("response body Close() error = %v", err)
				}
			} else if got != nil {
				t.Fatalf("Do() response = %v on error, want nil", got)
			}
		})
	}
}

func TestDoPropagatesCancellation(t *testing.T) {
	t.Parallel()

	started := make(chan struct{})
	hc := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		close(started)
		<-req.Context().Done()
		return nil, req.Context().Err()
	})}
	c, err := transport.New("https://example.test", hc, "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	req, err := c.NewRequest(ctx, http.MethodGet, "account", nil, nil)
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}
	result := make(chan error, 1)
	go func() {
		_, err := c.Do(req)
		result <- err
	}()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("HTTP round trip did not start")
	}
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Do() error = %v, want context.Canceled", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Do did not return after cancellation")
	}
}

func TestCloseIdleConnectionsDelegatesEveryCall(t *testing.T) {
	t.Parallel()

	tracker := &closeTrackingRoundTripper{roundTrip: func(*http.Request) (*http.Response, error) {
		return nil, errors.New("unused")
	}}
	c, err := transport.New("https://example.test", &http.Client{Transport: tracker}, "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	c.CloseIdleConnections()
	c.CloseIdleConnections()
	if got := tracker.closeCalls.Load(); got != 2 {
		t.Fatalf("underlying CloseIdleConnections() calls = %d, want 2", got)
	}
}
