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

	"github.com/shing1211/webullapi4go/data"
	"github.com/shing1211/webullapi4go/display"
)

const dsSnapshotBody = `[{"instrument_id":"913256135","pre_close":"101",` +
	`"change_ratio":"0.05","symbol":"AAPL","last_trade_time":1640688000000,` +
	`"price":"100","open":"100","close":"101","high":"105","low":"99",` +
	`"volume":"1000","change":"1.0","ask":"13.9","ask_size":"5","bid":"13.8",` +
	`"bid_size":"7","turnover":"0.01","eps":"7.465","eps_ttm":"7.465",` +
	`"lot_size":"1","bps":"4.991"}]`

const dsBarsBody = `{"result":[{"symbol":"AAPL","instrument_id":"913256135",` +
	`"result":[{"time":"2026-09-20T10:00:00Z","open":"380.0","close":"385.0",` +
	`"high":"388.0","low":"379.0","volume":"1000000"}]}]}`

const dsSingleBarBody = `{"symbol":"AAPL","instrument_id":"913256135",` +
	`"result":[{"time":"2026-09-20T10:00:00Z","open":"380.0","close":"385.0",` +
	`"high":"388.0","low":"379.0","volume":"1000000"}]}`

const dsTickBody = `{"symbol":"AAPL","instrument_id":"913256135",` +
	`"result":[{"time":"1640688000000","price":"385.0","volume":"100","side":"B"}]}`

const dsDepthBody = `{"symbol":"AAPL","instrument_id":"913256135",` +
	`"quote_time":1640688000000,` +
	`"asks":[{"price":"386.0","size":"100","order":[],"broker":[]}],` +
	`"bids":[{"price":"385.0","size":"200","order":[],"broker":[]}]}`

type dsQuotesVars struct {
	snapshotPath string
	snapshotSyms string
	snapshotCat  string

	barsPath string
	barsBody struct {
		Symbols  []string `json:"symbols"`
		Category string   `json:"category"`
		Timespan string   `json:"timespan"`
		Count    int      `json:"count"`
	}

	singleBarsPath string
	singleBarsSym  string
	singleBarsCat  string
	singleBarsTS   string

	tickPath string
	tickSym  string
	tickCat  string

	depthPath string
	depthSym  string
	depthCat  string
}

func newDisplayQuotesTestServer(t *testing.T) (*httptest.Server, *display.Service, *dsQuotesVars) {
	t.Helper()
	sv := &dsQuotesVars{}
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/client-tokens/create", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":       "test-access-token",
			"expires_at":         0,
			"refresh_token":      "refresh",
			"refresh_expires_at": 0,
		})
	})

	mux.HandleFunc("/openapi/market-data/stock/snapshot", func(w http.ResponseWriter, r *http.Request) {
		sv.snapshotPath = r.URL.Path
		sv.snapshotSyms = r.URL.Query().Get("symbols")
		sv.snapshotCat = r.URL.Query().Get("category")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(dsSnapshotBody))
	})

	mux.HandleFunc("/market-data/stocks/bars/list", func(w http.ResponseWriter, r *http.Request) {
		sv.barsPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &sv.barsBody)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(dsBarsBody))
	})

	mux.HandleFunc("/market-data/stocks/bars/get", func(w http.ResponseWriter, r *http.Request) {
		sv.singleBarsPath = r.URL.Path
		sv.singleBarsSym = r.URL.Query().Get("symbol")
		sv.singleBarsCat = r.URL.Query().Get("category")
		sv.singleBarsTS = r.URL.Query().Get("timespan")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(dsSingleBarBody))
	})

	mux.HandleFunc("/market-data/stocks/ticks/list", func(w http.ResponseWriter, r *http.Request) {
		sv.tickPath = r.URL.Path
		sv.tickSym = r.URL.Query().Get("symbol")
		sv.tickCat = r.URL.Query().Get("category")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(dsTickBody))
	})

	mux.HandleFunc("/market-data/stocks/depths/list", func(w http.ResponseWriter, r *http.Request) {
		sv.depthPath = r.URL.Path
		sv.depthSym = r.URL.Query().Get("symbol")
		sv.depthCat = r.URL.Query().Get("category")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(dsDepthBody))
	})

	srv := httptest.NewServer(mux)
	dsvc := display.NewService(testAppKey, testAppSecret,
		display.WithBaseURL(srv.URL),
		display.WithHTTPClient(srv.Client()),
	)
	return srv, dsvc, sv
}

func TestGetDisplaySnapshot(t *testing.T) {
	t.Parallel()
	srv, dsvc, sv := newDisplayQuotesTestServer(t)
	defer srv.Close()
	dc := newDataClientWithDisplay(t, dsvc)

	got, err := dc.GetDisplaySnapshot(context.Background(), data.SnapshotQuery{
		Symbols:  []string{"AAPL", "TSLA"},
		Category: data.StockCategoryUS,
	})
	if err != nil {
		t.Fatalf("GetDisplaySnapshot() error = %v", err)
	}
	if sv.snapshotPath != "/openapi/market-data/stock/snapshot" {
		t.Errorf("path = %q, want %q", sv.snapshotPath, "/openapi/market-data/stock/snapshot")
	}
	if sv.snapshotSyms != "AAPL,TSLA" {
		t.Errorf("symbols = %q, want %q", sv.snapshotSyms, "AAPL,TSLA")
	}
	if sv.snapshotCat != "US_STOCK" {
		t.Errorf("category = %q, want %q", sv.snapshotCat, "US_STOCK")
	}
	if len(got) != 1 {
		t.Fatalf("got %d snapshots, want 1", len(got))
	}
	s := got[0]
	if s.Symbol != "AAPL" || s.InstrumentID != "913256135" || s.Price != "100" {
		t.Errorf("identity/price = %+v", s)
	}
	if s.LastTradeTime != 1640688000000 {
		t.Errorf("LastTradeTime = %d, want 1640688000000", s.LastTradeTime)
	}
	if s.EPS != "7.465" || s.BPS != "4.991" {
		t.Errorf("fundamental fields = %+v", s)
	}
}

func TestGetDisplayBars(t *testing.T) {
	t.Parallel()
	srv, dsvc, sv := newDisplayQuotesTestServer(t)
	defer srv.Close()
	dc := newDataClientWithDisplay(t, dsvc)

	got, err := dc.GetDisplayBars(context.Background(), data.BatchBarQuery{
		Symbols:  []string{"AAPL"},
		Category: data.StockCategoryUS,
		Timespan: data.BarTimespanDay,
		Count:    100,
	})
	if err != nil {
		t.Fatalf("GetDisplayBars() error = %v", err)
	}
	if sv.barsPath != "/market-data/stocks/bars/list" {
		t.Errorf("path = %q, want %q", sv.barsPath, "/market-data/stocks/bars/list")
	}
	if len(sv.barsBody.Symbols) != 1 || sv.barsBody.Symbols[0] != "AAPL" {
		t.Errorf("body symbols = %v, want [AAPL]", sv.barsBody.Symbols)
	}
	if sv.barsBody.Category != "US_STOCK" {
		t.Errorf("body category = %q, want US_STOCK", sv.barsBody.Category)
	}
	if sv.barsBody.Timespan != "D" {
		t.Errorf("body timespan = %q, want D", sv.barsBody.Timespan)
	}
	if sv.barsBody.Count != 100 {
		t.Errorf("body count = %d, want 100", sv.barsBody.Count)
	}
	if len(got.Result) != 1 {
		t.Fatalf("got %d bar groups, want 1", len(got.Result))
	}
	if got.Result[0].Symbol != "AAPL" {
		t.Errorf("Symbol = %q, want AAPL", got.Result[0].Symbol)
	}
	if len(got.Result[0].Result) != 1 {
		t.Fatalf("got %d bars, want 1", len(got.Result[0].Result))
	}
	if got.Result[0].Result[0].Open != "380.0" || got.Result[0].Result[0].Close != "385.0" {
		t.Errorf("bar OHLCV = %+v", got.Result[0].Result[0])
	}
}

func TestGetDisplayBarsSingle(t *testing.T) {
	t.Parallel()
	srv, dsvc, sv := newDisplayQuotesTestServer(t)
	defer srv.Close()
	dc := newDataClientWithDisplay(t, dsvc)

	got, err := dc.GetDisplayBarsSingle(context.Background(), data.BarQuery{
		Symbol:   "AAPL",
		Category: data.StockCategoryUS,
		Interval: data.BarTimespanDay,
		Count:    50,
	})
	if err != nil {
		t.Fatalf("GetDisplayBarsSingle() error = %v", err)
	}
	if sv.singleBarsPath != "/market-data/stocks/bars/get" {
		t.Errorf("path = %q, want %q", sv.singleBarsPath, "/market-data/stocks/bars/get")
	}
	if sv.singleBarsSym != "AAPL" {
		t.Errorf("symbol = %q, want AAPL", sv.singleBarsSym)
	}
	if sv.singleBarsCat != "US_STOCK" {
		t.Errorf("category = %q, want US_STOCK", sv.singleBarsCat)
	}
	if sv.singleBarsTS != "D" {
		t.Errorf("timespan = %q, want D", sv.singleBarsTS)
	}
	if got.Symbol != "AAPL" || got.InstrumentID != "913256135" {
		t.Errorf("identity = %q/%q, want AAPL/913256135", got.Symbol, got.InstrumentID)
	}
	if len(got.Result) != 1 {
		t.Fatalf("got %d bars, want 1", len(got.Result))
	}
	if got.Result[0].Open != "380.0" {
		t.Errorf("bar open = %q, want 380.0", got.Result[0].Open)
	}
}

func TestGetDisplayTick(t *testing.T) {
	t.Parallel()
	srv, dsvc, sv := newDisplayQuotesTestServer(t)
	defer srv.Close()
	dc := newDataClientWithDisplay(t, dsvc)

	got, err := dc.GetDisplayTick(context.Background(), data.TickQuery{
		Symbol:   "AAPL",
		Category: data.StockCategoryUS,
		Count:    10,
	})
	if err != nil {
		t.Fatalf("GetDisplayTick() error = %v", err)
	}
	if sv.tickPath != "/market-data/stocks/ticks/list" {
		t.Errorf("path = %q, want %q", sv.tickPath, "/market-data/stocks/ticks/list")
	}
	if sv.tickSym != "AAPL" {
		t.Errorf("symbol = %q, want AAPL", sv.tickSym)
	}
	if sv.tickCat != "US_STOCK" {
		t.Errorf("category = %q, want US_STOCK", sv.tickCat)
	}
	if got.Symbol != "AAPL" || got.InstrumentID != "913256135" {
		t.Errorf("identity = %q/%q, want AAPL/913256135", got.Symbol, got.InstrumentID)
	}
	if len(got.Result) != 1 {
		t.Fatalf("got %d ticks, want 1", len(got.Result))
	}
	tk := got.Result[0]
	if tk.Price != "385.0" || tk.Volume != "100" || tk.Side != "B" {
		t.Errorf("tick = %+v, want price=385.0 volume=100 side=B", tk)
	}
}

func TestGetDisplayDepth(t *testing.T) {
	t.Parallel()
	srv, dsvc, sv := newDisplayQuotesTestServer(t)
	defer srv.Close()
	dc := newDataClientWithDisplay(t, dsvc)

	got, err := dc.GetDisplayDepth(context.Background(), data.DepthQuery{
		Symbol:   "AAPL",
		Category: data.StockCategoryUS,
		Depth:    10,
	})
	if err != nil {
		t.Fatalf("GetDisplayDepth() error = %v", err)
	}
	if sv.depthPath != "/market-data/stocks/depths/list" {
		t.Errorf("path = %q, want %q", sv.depthPath, "/market-data/stocks/depths/list")
	}
	if sv.depthSym != "AAPL" {
		t.Errorf("symbol = %q, want AAPL", sv.depthSym)
	}
	if sv.depthCat != "US_STOCK" {
		t.Errorf("category = %q, want US_STOCK", sv.depthCat)
	}
	if got.Symbol != "AAPL" || got.InstrumentID != "913256135" {
		t.Errorf("identity = %q/%q, want AAPL/913256135", got.Symbol, got.InstrumentID)
	}
	if len(got.Asks) != 1 || len(got.Bids) != 1 {
		t.Fatalf("ask/bid levels = %d/%d, want 1/1", len(got.Asks), len(got.Bids))
	}
	if got.Asks[0].Price != "386.0" || got.Asks[0].Size != "100" {
		t.Errorf("ask = %+v, want price=386.0 size=100", got.Asks[0])
	}
	if got.Bids[0].Price != "385.0" || got.Bids[0].Size != "200" {
		t.Errorf("bid = %+v, want price=385.0 size=200", got.Bids[0])
	}
}
