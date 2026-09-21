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

package broker_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/shing1211/webullapi4go/broker"
)

func TestGetStockInstruments(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/broker/instruments/stocks/list" {
			t.Errorf("path = %q, want /broker/instruments/stocks/list", r.URL.Path)
		}
		got := r.URL.Query().Get("symbols")
		if got != "AAPL,TSLA" {
			t.Errorf("symbols = %q, want AAPL,TSLA", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"symbol":"AAPL","name":"APPLE INC","exchange":"NSQ","currency":"USD","lot_size":"1","status":"OC"}]`))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)

	out, err := c.GetStockInstruments(context.Background(), []string{"AAPL", "TSLA"})
	if err != nil {
		t.Fatalf("GetStockInstruments() error = %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("got %d results, want 1", len(out))
	}
	if out[0].Symbol != "AAPL" {
		t.Errorf("Symbol = %q, want AAPL", out[0].Symbol)
	}
}

func TestGetStockLocate(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/broker/instruments/stock-locate/get" {
			t.Errorf("path = %q, want /broker/instruments/stock-locate/get", r.URL.Path)
		}
		if got := r.URL.Query().Get("symbol"); got != "TSLA" {
			t.Errorf("symbol = %q, want TSLA", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"symbol":"TSLA","locate_quantity":"1000","available":"500","rate":"0.015"}]`))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)

	out, err := c.GetStockLocate(context.Background(), "TSLA")
	if err != nil {
		t.Fatalf("GetStockLocate() error = %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("got %d results, want 1", len(out))
	}
	if out[0].Symbol != "TSLA" {
		t.Errorf("Symbol = %q, want TSLA", out[0].Symbol)
	}
	if out[0].LocateQuantity != "1000" {
		t.Errorf("LocateQuantity = %q, want 1000", out[0].LocateQuantity)
	}
}

func TestGetCorporateActionsDetail(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/broker/instruments/corporate-actions/get" {
			t.Errorf("path = %q, want /broker/instruments/corporate-actions/get", r.URL.Path)
		}
		if got := r.URL.Query().Get("symbol"); got != "MSFT" {
			t.Errorf("symbol = %q, want MSFT", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"symbol":"MSFT","action_type":"DIVIDEND","ex_date":"2026-06-10","record_date":"2026-06-11","payable_date":"2026-07-11","ratio":"0.75","description":"Quarterly dividend"}]`))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)

	out, err := c.GetCorporateActionsDetail(context.Background(), "MSFT")
	if err != nil {
		t.Fatalf("GetCorporateActionsDetail() error = %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("got %d results, want 1", len(out))
	}
	if out[0].ActionType != "DIVIDEND" {
		t.Errorf("ActionType = %q, want DIVIDEND", out[0].ActionType)
	}
}

func TestGetStockInstrumentsEmptySymbols(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.URL.Query().Get("symbols")
		if got != "" {
			t.Errorf("symbols = %q, want empty string", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)

	out, err := c.GetStockInstruments(context.Background(), []string{})
	if err != nil {
		t.Fatalf("GetStockInstruments() error = %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("got %d results, want 0", len(out))
	}
}

func TestGetStockInstrumentsError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"server error"}`))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)

	_, err := c.GetStockInstruments(context.Background(), []string{"AAPL"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetStockLocateMultipleResults(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"symbol":"AAPL","locate_quantity":"500","available":"250","rate":"0.01"},{"symbol":"AAPL","locate_quantity":"300","available":"100","rate":"0.02"}]`))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)

	out, err := c.GetStockLocate(context.Background(), "AAPL")
	if err != nil {
		t.Fatalf("GetStockLocate() error = %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("got %d results, want 2", len(out))
	}

	var rates []string
	for _, loc := range out {
		rates = append(rates, loc.Rate)
	}
	sort.Strings(rates)
	want := []string{"0.01", "0.02"}
	if len(rates) != len(want) {
		t.Fatalf("got %d rates, want %d", len(rates), len(want))
	}
	for i := range rates {
		if rates[i] != want[i] {
			t.Errorf("rate[%d] = %q, want %q", i, rates[i], want[i])
		}
	}
}

func TestGetCorporateActionsDetailMultipleActions(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"symbol":"GOOG","action_type":"SPLIT","ex_date":"2026-01-20","ratio":"20:1","description":"Stock split"},{"symbol":"GOOG","action_type":"DIVIDEND","ex_date":"2026-03-15","ratio":"0.20","description":"Quarterly dividend"}]`))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)

	out, err := c.GetCorporateActionsDetail(context.Background(), "GOOG")
	if err != nil {
		t.Fatalf("GetCorporateActionsDetail() error = %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("got %d results, want 2", len(out))
	}

	var types []string
	for _, ca := range out {
		types = append(types, ca.ActionType)
	}
	sort.Strings(types)
	want := []string{"DIVIDEND", "SPLIT"}
	if !equalStringSlice(types, want) {
		t.Errorf("action types = %v, want %v", types, want)
	}
}

func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestStockInstrumentJSONDeserialization(t *testing.T) {
	t.Parallel()

	payload := `{"symbol":"NVDA","name":"NVIDIA CORP","exchange":"NSQ","currency":"USD","lot_size":"1","status":"OC"}`
	var inst broker.StockInstrument
	if err := json.Unmarshal([]byte(payload), &inst); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if inst.Symbol != "NVDA" {
		t.Errorf("Symbol = %q, want NVDA", inst.Symbol)
	}
	if inst.Name != "NVIDIA CORP" {
		t.Errorf("Name = %q, want NVIDIA CORP", inst.Name)
	}
}

func TestStockLocateJSONDeserialization(t *testing.T) {
	t.Parallel()

	payload := `{"symbol":"AMZN","locate_quantity":"2000","available":"800","rate":"0.025"}`
	var loc broker.StockLocate
	if err := json.Unmarshal([]byte(payload), &loc); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if loc.Symbol != "AMZN" {
		t.Errorf("Symbol = %q, want AMZN", loc.Symbol)
	}
	if loc.Rate != "0.025" {
		t.Errorf("Rate = %q, want 0.025", loc.Rate)
	}
}

func TestCorporateActionDetailJSONDeserialization(t *testing.T) {
	t.Parallel()

	payload := `{"symbol":"META","action_type":"DIVIDEND","ex_date":"2026-09-15","record_date":"2026-09-16","payable_date":"2026-10-16","ratio":"0.50","description":"Quarterly dividend"}`
	var ca broker.CorporateActionDetail
	if err := json.Unmarshal([]byte(payload), &ca); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if ca.Symbol != "META" {
		t.Errorf("Symbol = %q, want META", ca.Symbol)
	}
	if ca.ExDate != "2026-09-15" {
		t.Errorf("ExDate = %q, want 2026-09-15", ca.ExDate)
	}
}

func TestGetStockInstrumentsQueryParamsEncoding(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.URL.Query().Get("symbols")
		if !strings.Contains(got, "BRK.B") {
			t.Errorf("symbols = %q, want it to contain BRK.B", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"symbol":"BRK.B","name":"BERKSHIRE HATHAWAY","exchange":"NYS","currency":"USD","lot_size":"1","status":"OC"}]`))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)

	out, err := c.GetStockInstruments(context.Background(), []string{"BRK.B"})
	if err != nil {
		t.Fatalf("GetStockInstruments() error = %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("got %d results, want 1", len(out))
	}
}
