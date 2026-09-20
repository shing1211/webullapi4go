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
)

func TestGetFundNav_pathAndQuery(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/market-data/fund/SPY/nav" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/market-data/fund/SPY/nav")
		}
		q := r.URL.Query()
		if q.Get("start_date") != "2026-01-01" {
			t.Errorf("start_date = %q, want %q", q.Get("start_date"), "2026-01-01")
		}
		if q.Get("end_date") != "2026-09-19" {
			t.Errorf("end_date = %q, want %q", q.Get("end_date"), "2026-09-19")
		}
		if q.Get("page_size") != "10" {
			t.Errorf("page_size = %q, want %q", q.Get("page_size"), "10")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetFundNav(context.Background(), data.FundNavQuery{
		Symbol:    "SPY",
		StartDate: "2026-01-01",
		EndDate:   "2026-09-19",
		PageSize:  10,
	})
	if err != nil {
		t.Fatalf("GetFundNav() error = %v", err)
	}
}

func TestGetFundNav_responseFields(t *testing.T) {
	t.Parallel()
	const body = `[{"symbol":"SPY","name":"SPDR S&P 500 ETF","currency":"USD","exchange":"NYSE","nav":"520.50","nav_date":"2026-09-19","prev_nav":"519.25","nav_change":"1.25","nav_change_ratio":"0.0024"}]`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	navs, err := c.GetFundNav(context.Background(), data.FundNavQuery{Symbol: "SPY"})
	if err != nil {
		t.Fatalf("GetFundNav() error = %v", err)
	}
	if len(navs) != 1 {
		t.Fatalf("len(navs) = %d, want 1", len(navs))
	}
	if navs[0].Nav != "520.50" || navs[0].NavChange != "1.25" {
		t.Errorf("nav = %+v", navs[0])
	}
}

func TestGetFundInfo_path(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/market-data/fund/VOO/info" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/market-data/fund/VOO/info")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetFundInfo(context.Background(), data.FundInfoQuery{Symbol: "VOO"})
	if err != nil {
		t.Fatalf("GetFundInfo() error = %v", err)
	}
}

func TestGetFundInfo_responseFields(t *testing.T) {
	t.Parallel()
	const body = `{"symbol":"VOO","name":"Vanguard S&P 500 ETF","currency":"USD","exchange":"NYSE","aum":"1000000","expense_ratio":"0.0003","dividend_yield":"0.012","inception_date":"2010-09-07","fund_type":"ETF","category":"US_STOCK"}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	info, err := c.GetFundInfo(context.Background(), data.FundInfoQuery{Symbol: "VOO"})
	if err != nil {
		t.Fatalf("GetFundInfo() error = %v", err)
	}
	if info.Aum != "1000000" || info.ExpenseRatio != "0.0003" {
		t.Errorf("info = %+v", info)
	}
	if info.FundType != "ETF" {
		t.Errorf("FundType = %q, want ETF", info.FundType)
	}
}

func TestGetFundDividends_pathAndQuery(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/market-data/fund/SPY/dividends" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/market-data/fund/SPY/dividends")
		}
		q := r.URL.Query()
		if q.Get("start_date") != "2026-01-01" {
			t.Errorf("start_date = %q, want %q", q.Get("start_date"), "2026-01-01")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetFundDividends(context.Background(), data.FundDividendsQuery{
		Symbol:    "SPY",
		StartDate: "2026-01-01",
	})
	if err != nil {
		t.Fatalf("GetFundDividends() error = %v", err)
	}
}

func TestGetFundDividends_responseFields(t *testing.T) {
	t.Parallel()
	const body = `[{"symbol":"SPY","name":"SPDR S&P 500 ETF","currency":"USD","exchange":"NYSE","amount":"1.89","ex_date":"2026-09-19","pay_date":"2026-09-22","record_date":"2026-09-19","frequency":"QUARTERLY"}]`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	divs, err := c.GetFundDividends(context.Background(), data.FundDividendsQuery{Symbol: "SPY"})
	if err != nil {
		t.Fatalf("GetFundDividends() error = %v", err)
	}
	if len(divs) != 1 {
		t.Fatalf("len(divs) = %d, want 1", len(divs))
	}
	if divs[0].Amount != "1.89" || divs[0].Frequency != "QUARTERLY" {
		t.Errorf("div = %+v", divs[0])
	}
}

func TestGetFundList_pathAndQuery(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/market-data/fund/list" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/market-data/fund/list")
		}
		q := r.URL.Query()
		if q.Get("market") != "US" {
			t.Errorf("market = %q, want %q", q.Get("market"), "US")
		}
		if q.Get("category") != "ETF" {
			t.Errorf("category = %q, want %q", q.Get("category"), "ETF")
		}
		if q.Get("page_size") != "5" {
			t.Errorf("page_size = %q, want %q", q.Get("page_size"), "5")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetFundList(context.Background(), data.FundListQuery{
		Market:   "US",
		Category: "ETF",
		PageSize: 5,
	})
	if err != nil {
		t.Fatalf("GetFundList() error = %v", err)
	}
}

func TestGetFundList_responseFields(t *testing.T) {
	t.Parallel()
	const body = `[{"symbol":"SPY","name":"SPDR S&P 500 ETF","currency":"USD","exchange":"NYSE","fund_type":"ETF","category":"US_STOCK","dividend_yield":"0.012"}]`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	funds, err := c.GetFundList(context.Background(), data.FundListQuery{Market: "US"})
	if err != nil {
		t.Fatalf("GetFundList() error = %v", err)
	}
	if len(funds) != 1 {
		t.Fatalf("len(funds) = %d, want 1", len(funds))
	}
	if funds[0].Symbol != "SPY" || funds[0].FundType != "ETF" {
		t.Errorf("fund = %+v", funds[0])
	}
}

func TestGetFundNav_omitsOptionalParams(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		for _, key := range []string{"start_date", "end_date", "page_size"} {
			if _, ok := q[key]; ok {
				t.Errorf("%s present, want omitted", key)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetFundNav(context.Background(), data.FundNavQuery{Symbol: "SPY"})
	if err != nil {
		t.Fatalf("GetFundNav() error = %v", err)
	}
}
