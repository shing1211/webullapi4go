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

package data_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/data"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

func TestGetNewsSummary(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/news/summaries/get"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got := r.Header.Get(client.AccessTokenHeader); got != "tok-news" {
			t.Errorf("%s = %q, want %q", client.AccessTokenHeader, got, "tok-news")
		}
		if got := r.Header.Get("x-version"); got != client.APIVersionV2 {
			t.Errorf("x-version = %q, want %q (news must go through the client pipeline)", got, client.APIVersionV2)
		}
		if r.Header.Get("x-signature") == "" {
			t.Error("missing x-signature header: request was not signed")
		}
		var req struct {
			CategorySymbols []struct {
				Category string   `json:"category"`
				Symbols  []string `json:"symbols"`
			} `json:"category_symbols"`
			Lang string `json:"lang"`
		}
		decodeRequestBody(t, r, &req)
		if len(req.CategorySymbols) != 1 || req.CategorySymbols[0].Category != "US_STOCK" {
			t.Errorf("category_symbols = %+v", req.CategorySymbols)
		}
		if len(req.CategorySymbols[0].Symbols) != 1 || req.CategorySymbols[0].Symbols[0] != "AAPL" {
			t.Errorf("symbols = %+v", req.CategorySymbols[0].Symbols)
		}
		if req.Lang != "en" {
			t.Errorf("lang = %q, want en", req.Lang)
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		flusher, _ := w.(http.Flusher)
		writeFrame := func(s string) {
			_, _ = io.WriteString(w, s)
			if flusher != nil {
				flusher.Flush()
			}
		}
		writeFrame(": keep-alive\n\n")
		writeFrame("event: message\n")
		writeFrame("data: {\"type\":\"meta\",\"args\":{\"sessionId\":\"1\",\"convId\":451107711450219}}\n\n")
		writeFrame("data: {\"type\":\"text\",\"message\":\"Hello\"}\n\n")
		writeFrame("data: {\"type\":\"table\",\"headers\":[{\"text\":\"Symbol\"}]," +
			"\"rows\":[[{\"text\":\"AAPL\"}]]}\n\n")
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	c.Core().SetToken(&client.Token{
		Value:     "tok-news",
		Status:    client.TokenStatusNormal,
		ExpiresAt: time.Now().Add(time.Hour),
	})
	stream, err := c.GetNewsSummary(context.Background(), data.NewsSummaryParam{
		CategorySymbols: []data.NewsCategorySymbols{
			{Category: data.StockCategoryUS, Symbols: []string{"AAPL"}},
		},
		Lang: "en",
	})
	if err != nil {
		t.Fatalf("GetNewsSummary() error = %v", err)
	}
	defer func() { _ = stream.Close() }()

	meta, err := stream.Next()
	if err != nil {
		t.Fatalf("Next() meta error = %v", err)
	}
	if meta.Type != "meta" || !strings.Contains(string(meta.Args), "sessionId") {
		t.Errorf("meta = %+v", meta)
	}

	text, err := stream.Next()
	if err != nil {
		t.Fatalf("Next() text error = %v", err)
	}
	if text.Type != "text" || text.Message != "Hello" {
		t.Errorf("text = %+v", text)
	}

	table, err := stream.Next()
	if err != nil {
		t.Fatalf("Next() table error = %v", err)
	}
	if table.Type != "table" || len(table.Headers) != 1 || table.Headers[0].Text != "Symbol" {
		t.Errorf("table headers = %+v", table.Headers)
	}
	if len(table.Rows) != 1 || len(table.Rows[0]) != 1 || table.Rows[0][0].Text != "AAPL" {
		t.Errorf("table rows = %+v", table.Rows)
	}

	if _, err := stream.Next(); err != io.EOF {
		t.Errorf("Next() at end = %v, want io.EOF", err)
	}
}

func TestGetNewsSummaryError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error_code":"UNAUTHORIZED","message":"Insufficient permission"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetNewsSummary(context.Background(), data.NewsSummaryParam{
		CategorySymbols: []data.NewsCategorySymbols{
			{Category: data.StockCategoryUS, Symbols: []string{"AAPL"}},
		},
	})
	if err == nil {
		t.Fatal("GetNewsSummary() error = nil, want unauthorized error")
	}
	if !strings.Contains(err.Error(), "Insufficient permission") {
		t.Errorf("error = %v, want permission message", err)
	}
	var e *errs.Error
	if !errors.As(err, &e) {
		t.Fatalf("GetNewsSummary() error type = %T, want *errs.Error", err)
	}
	if e.Code != errs.CodeUnauthorized {
		t.Errorf("error code = %q, want %q", e.Code, errs.CodeUnauthorized)
	}
	if e.Status != http.StatusUnauthorized {
		t.Errorf("error status = %d, want %d", e.Status, http.StatusUnauthorized)
	}
}

// TestSandboxNewsSummary opens the news-summary SSE stream against the shared
// Webull sandbox, proving that the stream flows through the client pipeline
// (token, signing, x-version) against a real server.
//
// It runs only when WEBULL_SANDBOX=1 and WEBULL_APP_KEY / WEBULL_APP_SECRET are
// set; otherwise it skips so the default test run stays hermetic. The stream is
// best-effort: an entitlement error, a server-side timeout of the upstream
// summarizer, or a clean, empty close is reported rather than failing the
// build. Never commit the sandbox credentials.
func TestSandboxNewsSummary(t *testing.T) {
	if os.Getenv("WEBULL_SANDBOX") != "1" {
		t.Skip("set WEBULL_SANDBOX=1, WEBULL_APP_KEY and WEBULL_APP_SECRET to run the sandbox test")
	}
	appKey := os.Getenv("WEBULL_APP_KEY")
	appSecret := os.Getenv("WEBULL_APP_SECRET")
	if appKey == "" || appSecret == "" {
		t.Skip("set WEBULL_APP_KEY and WEBULL_APP_SECRET to run the sandbox test")
	}

	opts := []client.Option{
		client.WithAppKey(appKey),
		client.WithAppSecret(appSecret),
		client.WithRegion(client.HK),
		client.WithEnvironment(client.Sandbox),
	}
	if baseURL := os.Getenv("WEBULL_BASE_URL"); baseURL != "" {
		opts = append(opts, client.WithBaseURL(baseURL))
	}

	cl, err := client.New(opts...)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if _, err := cl.EnsureToken(ctx); err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}

	market := data.New(cl)
	stream, err := market.GetNewsSummary(ctx, data.NewsSummaryParam{
		CategorySymbols: []data.NewsCategorySymbols{
			{Category: data.StockCategoryUS, Symbols: []string{"AAPL"}},
		},
		Lang: "en",
	})
	if err != nil {
		// The summary is generated by an upstream service that the sandbox
		// does not always serve (observed: HTTP 504 GATEWAY_TIMEOUT), and it
		// may require an entitlement the sandbox account does not hold.
		// Tolerate those conditions so the gated smoke run reports them
		// without failing; any other error is a real regression.
		if errs.Is(err, errs.CodeForbidden) || errs.Is(err, errs.CodeServer) {
			t.Skipf("news summary unavailable in sandbox: %v", err)
		}
		t.Fatalf("GetNewsSummary() error = %v", err)
	}
	defer func() { _ = stream.Close() }()

	// The summary may be empty; read a few events and accept a clean EOF.
	for i := 0; i < 5; i++ {
		event, err := stream.Next()
		if err == io.EOF {
			t.Logf("news summary stream closed cleanly after %d events", i)
			return
		}
		if err != nil {
			t.Fatalf("Next() error = %v", err)
		}
		t.Logf("news summary event %d: type=%q", i, event.Type)
	}
}
