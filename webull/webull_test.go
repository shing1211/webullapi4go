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

package webull_test

import (
	"net/http"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/webull"
)

var (
	_ *client.Client     = (*webull.Client)(nil)
	_ client.Option      = webull.WithCredentials("", "")
	_ client.Region      = webull.RegionHK
	_ client.Environment = webull.EnvSandbox
	_ client.Endpoints   = webull.Endpoints{}
)

type webullRoundTripper struct {
	closeCalls atomic.Int32
}

func (*webullRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       http.NoBody,
	}, nil
}

func (t *webullRoundTripper) CloseIdleConnections() { t.closeCalls.Add(1) }

func TestExportedFunctionAliases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		got  any
		want any
	}{
		{name: "WithCredentials", got: webull.WithCredentials, want: client.WithCredentials},
		{name: "WithAppKey", got: webull.WithAppKey, want: client.WithAppKey},
		{name: "WithAppSecret", got: webull.WithAppSecret, want: client.WithAppSecret},
		{name: "WithRegion", got: webull.WithRegion, want: client.WithRegion},
		{name: "WithEnvironment", got: webull.WithEnvironment, want: client.WithEnvironment},
		{name: "WithSandbox", got: webull.WithSandbox, want: client.WithSandbox},
		{name: "WithEndpoints", got: webull.WithEndpoints, want: client.WithEndpoints},
		{name: "WithBaseURL", got: webull.WithBaseURL, want: client.WithBaseURL},
		{name: "WithHTTPClient", got: webull.WithHTTPClient, want: client.WithHTTPClient},
		{name: "WithTimeout", got: webull.WithTimeout, want: client.WithTimeout},
		{name: "WithUserAgent", got: webull.WithUserAgent, want: client.WithUserAgent},
		{name: "WithRetry", got: webull.WithRetry, want: client.WithRetry},
		{name: "WithoutRetry", got: webull.WithoutRetry, want: client.WithoutRetry},
		{name: "WithRateLimiter", got: webull.WithRateLimiter, want: client.WithRateLimiter},
		{name: "WithBreaker", got: webull.WithBreaker, want: client.WithBreaker},
		{name: "NewRateLimiter", got: webull.NewRateLimiter, want: client.NewRateLimiter},
		{name: "NewBreaker", got: webull.NewBreaker, want: client.NewBreaker},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotPointer := reflect.ValueOf(tt.got).Pointer()
			wantPointer := reflect.ValueOf(tt.want).Pointer()
			if gotPointer != wantPointer {
				t.Fatalf("%s points to %x, want client function %x", tt.name, gotPointer, wantPointer)
			}
		})
	}
}

func TestConstantAliases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		got  string
		want string
	}{
		{name: "RegionHK", got: webull.RegionHK.String(), want: client.HK.String()},
		{name: "RegionUS", got: webull.RegionUS.String(), want: client.US.String()},
		{name: "EnvProduction", got: webull.EnvProduction.String(), want: client.Production.String()},
		{name: "EnvSandbox", got: webull.EnvSandbox.String(), want: client.Sandbox.String()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.got != tt.want {
				t.Fatalf("%s = %q, want %q", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestNewDelegatesOptionsAndLifecycle(t *testing.T) {
	t.Parallel()

	tracker := &webullRoundTripper{}
	hc := &http.Client{Transport: tracker}
	delegatedOption := webull.Option(func(cfg *client.Config) {
		cfg.UserAgent = "delegated-user-agent"
	})
	cl, err := webull.New(
		webull.WithCredentials("app-key", "app-secret"),
		webull.WithRegion(webull.RegionUS),
		webull.WithEnvironment(webull.EnvSandbox),
		webull.WithBaseURL("https://example.test/base"),
		webull.WithHTTPClient(hc),
		webull.WithoutRetry(),
		delegatedOption,
		nil,
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	cfg := cl.Config()
	if cfg.AppKey != "app-key" || cfg.AppSecret != "app-secret" {
		t.Fatalf("credentials = %q, %q", cfg.AppKey, cfg.AppSecret)
	}
	if cfg.Region != webull.RegionUS || cfg.Environment != webull.EnvSandbox {
		t.Fatalf("region/environment = %q/%q", cfg.Region, cfg.Environment)
	}
	if cfg.Endpoints.HTTP != "https://example.test/base" {
		t.Fatalf("HTTP endpoint = %q", cfg.Endpoints.HTTP)
	}
	if cfg.HTTPClient != hc {
		t.Fatal("New() did not retain the supplied HTTP client")
	}
	if cfg.UserAgent != "delegated-user-agent" {
		t.Fatalf("delegated option UserAgent = %q", cfg.UserAgent)
	}
	if err := cl.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if got := tracker.closeCalls.Load(); got != 1 {
		t.Fatalf("underlying CloseIdleConnections() calls = %d, want 1", got)
	}
}

func TestNewDelegatesValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts []webull.Option
	}{
		{name: "missing app key", opts: []webull.Option{webull.WithAppSecret("secret")}},
		{name: "missing app secret", opts: []webull.Option{webull.WithAppKey("key")}},
		{name: "unknown region", opts: []webull.Option{webull.WithCredentials("key", "secret"), webull.WithRegion(client.Region("unknown"))}},
		{name: "invalid base URL", opts: []webull.Option{webull.WithCredentials("key", "secret"), webull.WithBaseURL("relative")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cl, err := webull.New(tt.opts...)
			if cl != nil {
				t.Fatal("New() returned a client for invalid configuration")
			}
			if !errs.Is(err, errs.CodeInvalidConfig) {
				t.Fatalf("New() error = %v, want invalid-config code", err)
			}
		})
	}
}

func TestEndpointTypeAliasAcceptsClientValue(t *testing.T) {
	t.Parallel()

	value := client.Endpoints{HTTP: "https://example.test"}
	acceptEndpoints := func(webull.Endpoints) {}
	acceptEndpoints(value)
	cl, err := webull.New(
		webull.WithCredentials("key", "secret"),
		webull.WithEndpoints(value),
		webull.WithTimeout(time.Second),
	)
	if err != nil {
		t.Fatalf("New() with aliased Endpoints error = %v", err)
	}
	if err := cl.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
