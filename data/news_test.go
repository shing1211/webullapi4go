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
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shing1211/webullapi4go/data"
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
}
