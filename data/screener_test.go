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

const screenerBody = `[{"instrument_id":"913256135","symbol":"AAPL",` +
	`"name":"Apple Inc.","exchange_code":"NSQ","currency_code":"USD",` +
	`"pre_close":"380.2","open":"382.0","high":"388.6","low":"380.0",` +
	`"close":"385.6","price":"385.6","change":"5.4","change_ratio":"0.0142",` +
	`"volume":"12345678","turnover":"4756789012","turnover_rate":"0.0013",` +
	`"market_value":"3650000000000","amplitude":"0.0226",` +
	`"relative_volume_10d":"11.91"}]`

func TestGetTopGainersLosers(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/screeners/gainers-losers/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("rank_type"), "DAY_1"; got != want {
			t.Errorf("rank_type = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("sort_by"), "CHANGE_RATIO"; got != want {
			t.Errorf("sort_by = %q, want %q", got, want)
		}
		if got, want := q.Get("direction"), "ASC"; got != want {
			t.Errorf("direction = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(screenerBody))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetTopGainersLosers(context.Background(), data.GainersLosersQuery{
		RankType:  data.GainersLosersRankDay1,
		Category:  data.StockCategoryUS,
		SortBy:    data.ScreenerSortChangeRatio,
		Direction: data.SortDirectionAsc,
	})
	if err != nil {
		t.Fatalf("GetTopGainersLosers() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d stocks, want 1", len(got))
	}
	stock := got[0]
	if stock.Symbol != "AAPL" || stock.InstrumentID != "913256135" {
		t.Errorf("identity = %q/%q, want AAPL/913256135", stock.Symbol, stock.InstrumentID)
	}
	if stock.ExchangeCode != "NSQ" || stock.CurrencyCode != "USD" {
		t.Errorf("exchange/currency = %q/%q, want NSQ/USD", stock.ExchangeCode, stock.CurrencyCode)
	}
	if stock.ChangeRatio != "0.0142" || stock.RelativeVolume10D != "11.91" {
		t.Errorf("decoded ratios = %+v", stock)
	}
}

func TestGetMostActive(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/screeners/top-actives/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("rank_type"), "VOLUME"; got != want {
			t.Errorf("rank_type = %q, want %q", got, want)
		}
		if got, want := q.Get("sort_by"), "VOLUME"; got != want {
			t.Errorf("sort_by = %q, want %q", got, want)
		}
		if got, want := q.Get("direction"), "DESC"; got != want {
			t.Errorf("direction = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(screenerBody))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetMostActive(context.Background(), data.MostActiveQuery{
		Category:  data.StockCategoryUS,
		RankType:  data.MostActiveRankVolume,
		SortBy:    data.ScreenerSortVolume,
		Direction: data.SortDirectionDesc,
	})
	if err != nil {
		t.Fatalf("GetMostActive() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d stocks, want 1", len(got))
	}
	if got[0].Symbol != "AAPL" || got[0].RelativeVolume10D != "11.91" {
		t.Errorf("decoded stock = %+v", got[0])
	}
}

func TestGetMostActiveOmitsOptionalParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got, want := q.Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		for _, key := range []string{"rank_type", "sort_by", "direction"} {
			if _, ok := q[key]; ok {
				t.Errorf("%s present (%q), want omitted", key, q.Get(key))
			}
		}
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetMostActive(context.Background(), data.MostActiveQuery{
		Category: data.StockCategoryUS,
	})
	if err != nil {
		t.Fatalf("GetMostActive() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("got %d stocks, want 0", len(got))
	}
}
