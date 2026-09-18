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

func TestGetFootprint(t *testing.T) {
	t.Parallel()

	const body = `[{"symbol":"AAPL","instrument_id":"913256135","result":[` +
		`{"time":"2025-09-30T05:47:00.000+0000","trading_session":"RTH",` +
		`"total":"1000","delta":"200","buy_total":"600","sell_total":"400",` +
		`"buy_detail":{"24.20":"100","24.21":"60"},` +
		`"sell_detail":{"24.20":"50","24.21":"50"}}]}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/stocks/footprints/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbols"), "AAPL,TSLA"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("timespan"), "S5"; got != want {
			t.Errorf("timespan = %q, want %q", got, want)
		}
		if got, want := q.Get("count"), "500"; got != want {
			t.Errorf("count = %q, want %q", got, want)
		}
		if got, want := q.Get("real_time_required"), "true"; got != want {
			t.Errorf("real_time_required = %q, want %q", got, want)
		}
		if got, want := q.Get("trading_sessions"), "RTH"; got != want {
			t.Errorf("trading_sessions = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFootprint(context.Background(), data.FootprintQuery{
		Symbols:          []string{"AAPL", "TSLA"},
		Category:         data.StockCategoryUS,
		Timespan:         data.FootprintTimespanS5,
		Count:            500,
		RealTimeRequired: true,
		TradingSessions:  data.TradingSessionRTH,
	})
	if err != nil {
		t.Fatalf("GetFootprint() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d footprints, want 1", len(got))
	}
	fp := got[0]
	if fp.Symbol != "AAPL" || fp.InstrumentID != "913256135" {
		t.Errorf("identity = %q/%q, want AAPL/913256135", fp.Symbol, fp.InstrumentID)
	}
	if len(fp.Result) != 1 {
		t.Fatalf("got %d bars, want 1", len(fp.Result))
	}
	bar := fp.Result[0]
	if bar.TradingSession != data.TradingSessionRTH {
		t.Errorf("trading_session = %q, want RTH", bar.TradingSession)
	}
	if bar.Total != "1000" || bar.Delta != "200" || bar.BuyTotal != "600" || bar.SellTotal != "400" {
		t.Errorf("aggregates = %+v", bar)
	}
	if bar.BuyDetail["24.21"] != "60" || bar.SellDetail["24.20"] != "50" {
		t.Errorf("detail maps = buy:%v sell:%v", bar.BuyDetail, bar.SellDetail)
	}
}

func TestGetFootprintOmitsOptionalParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got, want := q.Get("symbols"), "AAPL"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("timespan"), "M1"; got != want {
			t.Errorf("timespan = %q, want %q", got, want)
		}
		if _, ok := q["count"]; ok {
			t.Errorf("count present (%q), want omitted", q.Get("count"))
		}
		if _, ok := q["trading_sessions"]; ok {
			t.Errorf("trading_sessions present (%q), want omitted", q.Get("trading_sessions"))
		}
		if got, want := q.Get("real_time_required"), "false"; got != want {
			t.Errorf("real_time_required = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFootprint(context.Background(), data.FootprintQuery{
		Symbols:  []string{"AAPL"},
		Category: data.StockCategoryUS,
		Timespan: data.FootprintTimespanM1,
	})
	if err != nil {
		t.Fatalf("GetFootprint() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("got %d footprints, want 0", len(got))
	}
}
