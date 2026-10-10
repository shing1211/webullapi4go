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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/data"
	"github.com/shing1211/webullapi4go/display"
)

type corpTestVars struct {
	path          string
	symbol        string
	category      string
	market        string
	eventTypes    string
	pageSize      string
	paginationKey string
}

func newCorpTestServer(t *testing.T) (*httptest.Server, *display.Service, *corpTestVars) {
	t.Helper()
	cv := &corpTestVars{}
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
	mux.HandleFunc("/market-data/instruments/stocks/corporate-actions/list", func(w http.ResponseWriter, r *http.Request) {
		cv.path = r.URL.Path
		cv.symbol = r.URL.Query().Get("symbol")
		cv.category = r.URL.Query().Get("category")
		cv.eventTypes = r.URL.Query().Get("event_types")
		cv.pageSize = r.URL.Query().Get("page_size")
		cv.paginationKey = r.URL.Query().Get("pagination_key")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"instrument_id": 913256135,
					"symbol":        "AAPL",
					"exchange_code": "NSQ",
					"event_type":    "DIVIDEND",
					"event_action":  "DISTRIBUTE",
					"event_id":      12345,
					"source":        "WEBULL",
					"ratio_old":     "1.0",
					"ratio_new":     "1.0",
					"event_date":    "2026-09-19",
					"update_time":   "2026-09-19T10:00:00Z",
					"create_time":   "2026-09-19T10:00:00Z",
				},
			},
			"pagination_key": "next-page",
		})
	})
	mux.HandleFunc("/market-data/instruments/stocks/corporate-actions/list-by-market", func(w http.ResponseWriter, r *http.Request) {
		cv.path = r.URL.Path
		cv.market = r.URL.Query().Get("market")
		cv.eventTypes = r.URL.Query().Get("event_types")
		cv.pageSize = r.URL.Query().Get("page_size")
		w.Header().Set("Content-Type", "application/json")
		q := r.URL.Query()
		if q.Get("market") == "" {
			http.Error(w, "missing market param", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"instrument_id": 913256135,
					"symbol":        "AAPL",
					"exchange_code": "NSQ",
					"event_type":    "SPLIT",
					"event_action":  "SPLIT",
					"event_id":      67890,
					"source":        "WEBULL",
					"ratio_old":     "4.0",
					"ratio_new":     "1.0",
					"event_date":    "2026-09-20",
					"update_time":   "2026-09-20T10:00:00Z",
					"create_time":   "2026-09-20T10:00:00Z",
				},
			},
			"pagination_key": "",
		})
	})

	srv := httptest.NewServer(mux)
	dsvc := display.NewService(testAppKey, testAppSecret,
		display.WithBaseURL(srv.URL),
		display.WithHTTPClient(srv.Client()),
	)
	return srv, dsvc, cv
}

func newDataClientWithDisplay(t *testing.T, dsvc *display.Service) *data.Client {
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

func TestGetCorporateActions_pathAndQuery(t *testing.T) {
	t.Parallel()
	srv, dsvc, cv := newCorpTestServer(t)
	defer srv.Close()
	dc := newDataClientWithDisplay(t, dsvc)

	_, _, err := dc.GetCorporateActions(context.Background(), data.CorporateActionQuery{
		Symbols: []string{"AAPL"},
		Market:  "US",
	})
	if err != nil {
		t.Fatalf("GetCorporateActions() error = %v", err)
	}
	if cv.path != "/market-data/instruments/stocks/corporate-actions/list" {
		t.Errorf("path = %q, want %q", cv.path, "/market-data/instruments/stocks/corporate-actions/list")
	}
	if cv.symbol != "AAPL" {
		t.Errorf("symbol = %q, want %q", cv.symbol, "AAPL")
	}
	if cv.category != "US" {
		t.Errorf("category = %q, want %q", cv.category, "US")
	}
}

func TestGetCorporateActions_responseFields(t *testing.T) {
	t.Parallel()
	srv, dsvc, _ := newCorpTestServer(t)
	defer srv.Close()
	dc := newDataClientWithDisplay(t, dsvc)

	actions, paginationKey, err := dc.GetCorporateActions(context.Background(), data.CorporateActionQuery{
		Symbols: []string{"AAPL"},
		Market:  "US",
	})
	if err != nil {
		t.Fatalf("GetCorporateActions() error = %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("len(actions) = %d, want 1", len(actions))
	}
	a := actions[0]
	if a.InstrumentID != 913256135 {
		t.Errorf("InstrumentID = %d, want 913256135", a.InstrumentID)
	}
	if a.Symbol != "AAPL" {
		t.Errorf("Symbol = %q, want %q", a.Symbol, "AAPL")
	}
	if a.EventType != "DIVIDEND" {
		t.Errorf("EventType = %q, want %q", a.EventType, "DIVIDEND")
	}
	if a.EventAction != "DISTRIBUTE" {
		t.Errorf("EventAction = %q, want %q", a.EventAction, "DISTRIBUTE")
	}
	if a.RatioOld != "1.0" || a.RatioNew != "1.0" {
		t.Errorf("ratio = %q/%q, want 1.0/1.0", a.RatioOld, a.RatioNew)
	}
	if paginationKey != "next-page" {
		t.Errorf("paginationKey = %q, want %q", paginationKey, "next-page")
	}
}

func TestGetCorporateActions_paginationKey(t *testing.T) {
	t.Parallel()
	srv, dsvc, _ := newCorpTestServer(t)
	defer srv.Close()
	dc := newDataClientWithDisplay(t, dsvc)

	_, pk, err := dc.GetCorporateActions(context.Background(), data.CorporateActionQuery{
		Symbols:       []string{"AAPL"},
		Market:        "US",
		PaginationKey: "page-2",
	})
	if err != nil {
		t.Fatalf("GetCorporateActions() error = %v", err)
	}
	if pk != "next-page" {
		t.Errorf("paginationKey = %q, want %q", pk, "next-page")
	}
}

func TestGetCorporateActionsByMarket_pathAndQuery(t *testing.T) {
	t.Parallel()
	srv, dsvc, cv := newCorpTestServer(t)
	defer srv.Close()
	dc := newDataClientWithDisplay(t, dsvc)

	_, _, err := dc.GetCorporateActionsByMarket(context.Background(), data.CorporateActionQuery{
		Market: "US",
	})
	if err != nil {
		t.Fatalf("GetCorporateActionsByMarket() error = %v", err)
	}
	if cv.path != "/market-data/instruments/stocks/corporate-actions/list-by-market" {
		t.Errorf("path = %q, want %q", cv.path, "/market-data/instruments/stocks/corporate-actions/list-by-market")
	}
	if cv.market != "US" {
		t.Errorf("market = %q, want %q", cv.market, "US")
	}
}

func TestGetCorporateActionsByMarket_responseFields(t *testing.T) {
	t.Parallel()
	srv, dsvc, _ := newCorpTestServer(t)
	defer srv.Close()
	dc := newDataClientWithDisplay(t, dsvc)

	actions, pk, err := dc.GetCorporateActionsByMarket(context.Background(), data.CorporateActionQuery{
		Market: "US",
	})
	if err != nil {
		t.Fatalf("GetCorporateActionsByMarket() error = %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("len(actions) = %d, want 1", len(actions))
	}
	a := actions[0]
	if a.Symbol != "AAPL" {
		t.Errorf("Symbol = %q, want %q", a.Symbol, "AAPL")
	}
	if a.EventType != "SPLIT" {
		t.Errorf("EventType = %q, want %q", a.EventType, "SPLIT")
	}
	if a.RatioOld != "4.0" || a.RatioNew != "1.0" {
		t.Errorf("ratio = %q/%q, want 4.0/1.0", a.RatioOld, a.RatioNew)
	}
	if pk != "" {
		t.Errorf("paginationKey = %q, want empty", pk)
	}
}
