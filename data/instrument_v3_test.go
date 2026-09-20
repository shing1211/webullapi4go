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

func newProfileTestServer(t *testing.T) (*httptest.Server, *display.Service) {
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
	mux.HandleFunc("/market-data/instruments/stocks/profiles/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"symbol":        "AAPL",
					"instrument_id": "913256135",
					"exchange_code": "NSQ",
					"name":          "Apple Inc.",
					"currency":      "USD",
					"status":        "ACTIVE",
				},
				{
					"symbol":        "TSLA",
					"instrument_id": "877",
					"exchange_code": "NGQ",
					"name":          "Tesla Inc.",
					"currency":      "USD",
					"status":        "ACTIVE",
				},
			},
			"pagination_key": "page-2",
		})
	})
	srv := httptest.NewServer(mux)
	dsvc := display.NewService(testAppKey, testAppSecret,
		display.WithBaseURL(srv.URL),
		display.WithHTTPClient(srv.Client()),
	)
	return srv, dsvc
}

func TestGetStockProfilesV3_pathAndQuery(t *testing.T) {
	t.Parallel()
	srv, dsvc := newProfileTestServer(t)
	defer srv.Close()

	cl, _ := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithBaseURL(srv.URL),
	)
	dc := data.New(cl)
	dc.SetDisplaySvcForTesting(dsvc)

	_, err := dc.GetStockProfilesV3(context.Background(), data.StockProfilesV3Query{
		Symbols:  []string{"AAPL", "TSLA"},
		Category: data.StockCategoryUS,
	})
	if err != nil {
		t.Fatalf("GetStockProfilesV3() error = %v", err)
	}
}

func TestGetStockProfilesV3_responseFields(t *testing.T) {
	t.Parallel()
	srv, dsvc := newProfileTestServer(t)
	defer srv.Close()

	cl, _ := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithBaseURL(srv.URL),
	)
	dc := data.New(cl)
	dc.SetDisplaySvcForTesting(dsvc)

	result, err := dc.GetStockProfilesV3(context.Background(), data.StockProfilesV3Query{
		Symbols:  []string{"AAPL", "TSLA"},
		Category: data.StockCategoryUS,
	})
	if err != nil {
		t.Fatalf("GetStockProfilesV3() error = %v", err)
	}
	if len(result.Instruments) != 2 {
		t.Fatalf("len(instruments) = %d, want 2", len(result.Instruments))
	}
	if result.PaginationKey != "page-2" {
		t.Errorf("paginationKey = %q, want %q", result.PaginationKey, "page-2")
	}
	if result.Instruments[0].Symbol != "AAPL" {
		t.Errorf("instrument[0] = %+v", result.Instruments[0])
	}
}

func TestGetStockProfilesV3_paginationKey(t *testing.T) {
	t.Parallel()
	srv, dsvc := newProfileTestServer(t)
	defer srv.Close()

	cl, _ := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithBaseURL(srv.URL),
	)
	dc := data.New(cl)
	dc.SetDisplaySvcForTesting(dsvc)

	result, err := dc.GetStockProfilesV3(context.Background(), data.StockProfilesV3Query{
		Symbols:       []string{"AAPL"},
		Category:      data.StockCategoryUS,
		PaginationKey: "page-1",
	})
	if err != nil {
		t.Fatalf("GetStockProfilesV3() error = %v", err)
	}
	if result.PaginationKey != "page-2" {
		t.Errorf("paginationKey = %q, want %q", result.PaginationKey, "page-2")
	}
}
