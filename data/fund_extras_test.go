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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shing1211/webullapi4go/data"
	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

func fundTestServer(t *testing.T, wantPath, wantSymbol, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Path; got != wantPath {
			t.Errorf("path = %q, want %q", got, wantPath)
		}
		q := r.URL.Query()
		if q.Get("symbol") != wantSymbol {
			t.Errorf("symbol = %q, want %q", q.Get("symbol"), wantSymbol)
		}
		if q.Get("category") != "US_STOCK" {
			t.Errorf("category = %q, want US_STOCK", q.Get("category"))
		}
		_, _ = w.Write([]byte(body))
	}))
}

func TestGetFundPerformance(t *testing.T) {
	t.Parallel()
	srv := fundTestServer(t, "/market-data/fundamentals/fund-performances/get", "SPY",
		`{"currency":"USD","end_date":"2026-09-19","return_1y":"0.12"}`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetFundPerformance(context.Background(), "SPY", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetFundPerformance() error = %v", err)
	}
	if got.Return1Y != "0.12" {
		t.Errorf("return_1y = %q", got.Return1Y)
	}
}

func TestGetFundHoldings(t *testing.T) {
	t.Parallel()
	srv := fundTestServer(t, "/market-data/fundamentals/fund-holdings/get", "SPY",
		`[{"target_symbol":"AAPL","stock_name":"Apple","share_held_pct":"7.0"}]`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetFundHoldings(context.Background(), "SPY", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetFundHoldings() error = %v", err)
	}
	if len(got) != 1 || got[0].TargetSymbol != "AAPL" {
		t.Fatalf("holdings = %+v", got)
	}
}

func TestGetFundRating(t *testing.T) {
	t.Parallel()
	srv := fundTestServer(t, "/market-data/fundamentals/fund-ratings/get", "SPY",
		`[{"rating_date":"2026-01-01","rating_agency":"MORNINGSTAR","rating_cycle":"Q","rating_results":5}]`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetFundRating(context.Background(), "SPY", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetFundRating() error = %v", err)
	}
	if len(got) != 1 || got[0].RatingResults != 5 {
		t.Fatalf("rating = %+v", got)
	}
}

func TestGetFundSplits(t *testing.T) {
	t.Parallel()
	srv := fundTestServer(t, "/market-data/fundamentals/fund-splits/get", "SPY",
		`[{"split_date":"2026-01-01","split_type":"FORWARD","split_ratio":"2:1","from":1,"to":2}]`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetFundSplits(context.Background(), "SPY", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetFundSplits() error = %v", err)
	}
	if len(got) != 1 || got[0].To != 2 {
		t.Fatalf("splits = %+v", got)
	}
}

func TestGetFundFiles(t *testing.T) {
	t.Parallel()
	srv := fundTestServer(t, "/market-data/fundamentals/fund-files/get", "SPY",
		`[{"publish_date":"2026-01-01","url":"https://x","type":1,"file_name":"a.pdf"}]`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetFundFiles(context.Background(), "SPY", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetFundFiles() error = %v", err)
	}
	if len(got) != 1 || got[0].FileName != "a.pdf" {
		t.Fatalf("files = %+v", got)
	}
}

func TestGetFundAllocation(t *testing.T) {
	t.Parallel()
	srv := fundTestServer(t, "/market-data/fundamentals/fund-allocations/get", "SPY",
		`[{"date":"2026-01-01","aum":"1000","stock":{"pct":"99"},"cash":{"pct":"1"}}]`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetFundAllocation(context.Background(), "SPY", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetFundAllocation() error = %v", err)
	}
	if len(got) != 1 || got[0].Aum.Cmp(money.Must(money.NewFromString("1000"))) != 0 || got[0].Stock["pct"] != "99" {
		t.Fatalf("allocation = %+v", got)
	}
}
