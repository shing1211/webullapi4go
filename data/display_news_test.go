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
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/data"
	"github.com/shing1211/webullapi4go/display"
)

func newDSNewsTestServer(t *testing.T) (*httptest.Server, *display.Service) {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/client-tokens/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":       "test-access-token",
			"expires_at":         0,
			"refresh_token":      "refresh",
			"refresh_expires_at": 0,
		})
	})
	mux.HandleFunc("/market-data/news/summaries/get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"title": "Watchlist News", "content": "Content 1", "source": "Reuters", "symbol": "AAPL"},
		})
	})
	mux.HandleFunc("/market-data/news/market-news/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"title": "Market News", "content": "Market content", "source": "Bloomberg", "symbol": "SPY"},
		})
	})
	mux.HandleFunc("/market-data/news/symbol-news/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"title": "Ticker News", "content": "AAPL update", "source": "CNBC", "symbol": "AAPL"},
		})
	})
	mux.HandleFunc("/market-data/news/latest-news/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"title": "Latest News", "content": "Breaking", "source": "AP", "symbol": ""},
		})
	})

	srv := httptest.NewServer(mux)
	dsvc := display.NewService(testAppKey, testAppSecret,
		display.WithBaseURL(srv.URL),
		display.WithHTTPClient(srv.Client()),
	)
	return srv, dsvc
}

func newDSNewsDataClient(t *testing.T, dsvc *display.Service) *data.Client {
	t.Helper()
	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })
	dc := data.New(cl)
	dc.SetDisplaySvcForTesting(dsvc)
	return dc
}

func TestGetDSNewsSummary(t *testing.T) {
	t.Parallel()
	srv, dsvc := newDSNewsTestServer(t)
	defer srv.Close()
	dc := newDSNewsDataClient(t, dsvc)

	items, err := dc.GetDSNewsSummary(context.Background(), []string{"AAPL"})
	if err != nil {
		t.Fatalf("GetDSNewsSummary() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].Title != "Watchlist News" {
		t.Errorf("Title = %q, want %q", items[0].Title, "Watchlist News")
	}
}

func TestGetDSMarketNews(t *testing.T) {
	t.Parallel()
	srv, dsvc := newDSNewsTestServer(t)
	defer srv.Close()
	dc := newDSNewsDataClient(t, dsvc)

	items, err := dc.GetDSMarketNews(context.Background(), "US")
	if err != nil {
		t.Fatalf("GetDSMarketNews() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].Title != "Market News" {
		t.Errorf("Title = %q, want %q", items[0].Title, "Market News")
	}
}

func TestGetDSSymbolNews(t *testing.T) {
	t.Parallel()
	srv, dsvc := newDSNewsTestServer(t)
	defer srv.Close()
	dc := newDSNewsDataClient(t, dsvc)

	items, err := dc.GetDSSymbolNews(context.Background(), "AAPL")
	if err != nil {
		t.Fatalf("GetDSSymbolNews() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].Title != "Ticker News" {
		t.Errorf("Title = %q, want %q", items[0].Title, "Ticker News")
	}
}

func TestGetDSLatestNews(t *testing.T) {
	t.Parallel()
	srv, dsvc := newDSNewsTestServer(t)
	defer srv.Close()
	dc := newDSNewsDataClient(t, dsvc)

	items, err := dc.GetDSLatestNews(context.Background())
	if err != nil {
		t.Fatalf("GetDSLatestNews() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].Title != "Latest News" {
		t.Errorf("Title = %q, want %q", items[0].Title, "Latest News")
	}
}

func TestGetDSNewsSummary_bodyContainsSymbols(t *testing.T) {
	t.Parallel()
	var capturedBody []byte
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/client-tokens/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":       "test-access-token",
			"expires_at":         0,
			"refresh_token":      "refresh",
			"refresh_expires_at": 0,
		})
	})
	mux.HandleFunc("/market-data/news/summaries/get", func(w http.ResponseWriter, r *http.Request) {
		capturedBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	dsvc := display.NewService(testAppKey, testAppSecret,
		display.WithBaseURL(srv.URL),
		display.WithHTTPClient(srv.Client()),
	)
	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	defer cl.Close()
	dc := data.New(cl)
	dc.SetDisplaySvcForTesting(dsvc)

	symbols := []string{"AAPL", "GOOG"}
	_, err = dc.GetDSNewsSummary(context.Background(), symbols)
	if err != nil {
		t.Fatalf("GetDSNewsSummary() error = %v", err)
	}

	var got []string
	if err := json.Unmarshal(capturedBody, &got); err != nil {
		t.Fatalf("decoding captured body: %v", err)
	}
	if len(got) != 2 || got[0] != "AAPL" || got[1] != "GOOG" {
		t.Errorf("body symbols = %v, want [AAPL GOOG]", got)
	}
}
