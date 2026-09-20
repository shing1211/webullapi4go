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

	"github.com/shing1211/webullapi4go/data"
)

func TestGetScreenerV2_pathIsPOST(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/wlas/screener/ng/query" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/wlas/screener/ng/query")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[],"pagination_key":""}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetScreenerV2(context.Background(), data.ScreenerV2Query{})
	if err != nil {
		t.Fatalf("GetScreenerV2() error = %v", err)
	}
}

func TestGetScreenerV2_requestBody(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body data.ScreenerV2Query
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode error: %v", err)
		}
		if body.PageSize != 20 {
			t.Errorf("PageSize = %d, want 20", body.PageSize)
		}
		if body.Filter == nil {
			t.Fatal("Filter is nil")
		}
		if len(body.Filter.Conditions) != 1 {
			t.Fatalf("len(conditions) = %d, want 1", len(body.Filter.Conditions))
		}
		if body.Filter.Conditions[0].Field != "price" {
			t.Errorf("condition[0].Field = %q, want price", body.Filter.Conditions[0].Field)
		}
		if body.Filter.Conditions[0].Operator != data.ScreenerV2OpGT {
			t.Errorf("condition[0].Operator = %q, want gt", body.Filter.Conditions[0].Operator)
		}
		if len(body.Sort) != 1 {
			t.Fatalf("len(sort) = %d, want 1", len(body.Sort))
		}
		if body.Sort[0].Field != "volume" {
			t.Errorf("sort[0].Field = %q, want volume", body.Sort[0].Field)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[],"pagination_key":""}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetScreenerV2(context.Background(), data.ScreenerV2Query{
		Filter: &data.ScreenerV2Filter{
			Conditions: []data.ScreenerV2Condition{
				{Field: "price", Operator: data.ScreenerV2OpGT, Value: "100"},
			},
		},
		Sort:     []data.ScreenerV2Sort{{Field: "volume", Direction: data.SortDirectionDesc}},
		PageSize: 20,
	})
	if err != nil {
		t.Fatalf("GetScreenerV2() error = %v", err)
	}
}

func TestGetScreenerV2_nestedFilter(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body data.ScreenerV2Query
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode error: %v", err)
		}
		if body.Filter == nil || len(body.Filter.Filters) != 1 {
			t.Fatalf("expected nested filter, got %+v", body.Filter)
		}
		if body.Filter.Logic != data.ScreenerV2LogicAnd {
			t.Errorf("filter.Logic = %q, want AND", body.Filter.Logic)
		}
		if body.Filter.Filters[0].Logic != data.ScreenerV2LogicOr {
			t.Errorf("subfilter.Logic = %q, want OR", body.Filter.Filters[0].Logic)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[],"pagination_key":""}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetScreenerV2(context.Background(), data.ScreenerV2Query{
		Filter: &data.ScreenerV2Filter{
			Logic: data.ScreenerV2LogicAnd,
			Filters: []data.ScreenerV2Filter{
				{
					Logic: data.ScreenerV2LogicOr,
					Conditions: []data.ScreenerV2Condition{
						{Field: "price", Operator: data.ScreenerV2OpGT, Value: "100"},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("GetScreenerV2() error = %v", err)
	}
}

func TestGetScreenerV2_responseFields(t *testing.T) {
	t.Parallel()
	const body = `{"data":[{"instrument_id":"913256135","symbol":"AAPL","name":"Apple Inc.","exchange_code":"NSQ","currency_code":"USD","pre_close":"380.2","open":"382.0","high":"388.6","low":"380.0","close":"385.6","price":"385.6","change":"5.4","change_ratio":"0.0142","volume":"12345678","turnover":"4756789012","market_value":"3650000000000"}],"pagination_key":"next-key"}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	result, err := c.GetScreenerV2(context.Background(), data.ScreenerV2Query{PageSize: 10})
	if err != nil {
		t.Fatalf("GetScreenerV2() error = %v", err)
	}
	if len(result.Stocks) != 1 {
		t.Fatalf("len(stocks) = %d, want 1", len(result.Stocks))
	}
	if result.PaginationKey != "next-key" {
		t.Errorf("paginationKey = %q, want %q", result.PaginationKey, "next-key")
	}
	s := result.Stocks[0]
	if s.Symbol != "AAPL" {
		t.Errorf("symbol = %q, want AAPL", s.Symbol)
	}
}

func TestGetScreenerV2_paginationKey(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body data.ScreenerV2Query
		json.NewDecoder(r.Body).Decode(&body)
		if body.PaginationKey != "current-key" {
			t.Errorf("paginationKey = %q, want %q", body.PaginationKey, "current-key")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[],"pagination_key":"next-key"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	result, err := c.GetScreenerV2(context.Background(), data.ScreenerV2Query{
		PaginationKey: "current-key",
	})
	if err != nil {
		t.Fatalf("GetScreenerV2() error = %v", err)
	}
	if result.PaginationKey != "next-key" {
		t.Errorf("paginationKey = %q, want %q", result.PaginationKey, "next-key")
	}
}
