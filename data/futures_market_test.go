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

func TestGetFuturesTick(t *testing.T) {
	t.Parallel()

	const body = `{"symbol":"ESZ5","instrument_id":"123456789","result":[` +
		`{"time":"1761182953043","price":"4807.00","volume":"1","side":"S"},` +
		`{"time":"1761182953042","price":"4806.50","volume":"3","side":"B"}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/futures/ticks/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbol"), "ESZ5"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_FUTURES"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("count"), "30"; got != want {
			t.Errorf("count = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFuturesTick(context.Background(), data.FuturesTickQuery{
		Symbol: "ESZ5",
		Count:  30,
	})
	if err != nil {
		t.Fatalf("GetFuturesTick() error = %v", err)
	}
	if got.Symbol != "ESZ5" || got.InstrumentID != "123456789" {
		t.Errorf("identity = %+v", got)
	}
	if len(got.Result) != 2 {
		t.Fatalf("got %d ticks, want 2", len(got.Result))
	}
	if got.Result[0].Price.Cmp(money.Must(money.NewFromString("4807.00"))) != 0 || got.Result[0].Side != "S" {
		t.Errorf("first tick = %+v", got.Result[0])
	}
	if got.Result[1].Side != "B" {
		t.Errorf("second tick side = %q, want B", got.Result[1].Side)
	}
}

func TestGetFuturesTickOmitsOptionalParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got, want := q.Get("symbol"), "ESZ5"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_FUTURES"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if _, ok := q["count"]; ok {
			t.Errorf("count present (%q), want omitted", q.Get("count"))
		}
		_, _ = w.Write([]byte(`{"symbol":"ESZ5","instrument_id":"1","result":[]}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	if _, err := c.GetFuturesTick(context.Background(), data.FuturesTickQuery{
		Symbol: "ESZ5",
	}); err != nil {
		t.Fatalf("GetFuturesTick() error = %v", err)
	}
}

func TestGetFuturesSnapshot(t *testing.T) {
	t.Parallel()

	const body = `[{"instrument_id":"123456789","pre_close":"4800","change_ratio":"0.01",` +
		`"symbol":"ESZ5","last_trade_time":1640688000000,"price":"4848","open":"4800",` +
		`"close":"4800","high":"4850","low":"4790","volume":"50000","change":"48",` +
		`"ask":"4848.50","ask_size":"10","bid":"4848.00","bid_size":"15"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/futures/snapshots/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbols"), "ESZ5,NQZ5"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_FUTURES"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFuturesSnapshot(context.Background(), data.FuturesSnapshotQuery{
		Symbols: []string{"ESZ5", "NQZ5"},
	})
	if err != nil {
		t.Fatalf("GetFuturesSnapshot() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d snapshots, want 1", len(got))
	}
	if got[0].Symbol != "ESZ5" || got[0].Price.Cmp(money.Must(money.NewFromString("4848"))) != 0 {
		t.Errorf("identity/price = %+v", got[0])
	}
	if got[0].LastTradeTime != 1640688000000 {
		t.Errorf("last_trade_time = %d, want 1640688000000", got[0].LastTradeTime)
	}
}

func TestGetFuturesSnapshotOmitsEmptySymbols(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if _, ok := q["symbols"]; ok {
			t.Errorf("symbols present (%q), want omitted", q.Get("symbols"))
		}
		if got, want := q.Get("category"), "US_FUTURES"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFuturesSnapshot(context.Background(), data.FuturesSnapshotQuery{})
	if err != nil {
		t.Fatalf("GetFuturesSnapshot() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("got %d snapshots, want 0", len(got))
	}
}

func TestGetFuturesBars(t *testing.T) {
	t.Parallel()

	const body = `{"result":[{"symbol":"ESZ5","instrument_id":"123456789","result":[` +
		`{"time":"2025-12-28T09:00:09.945+0000","open":"4800.00","close":"4820.00",` +
		`"high":"4825.00","low":"4795.00","volume":"125000"}]}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/futures/bars/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbols"), "ESZ5"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_FUTURES"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("timespan"), "M5"; got != want {
			t.Errorf("timespan = %q, want %q", got, want)
		}
		if got, want := q.Get("count"), "100"; got != want {
			t.Errorf("count = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFuturesBars(context.Background(), data.FuturesBarsQuery{
		Symbols:  []string{"ESZ5"},
		Interval: data.BarTimespanM5,
		Count:    100,
	})
	if err != nil {
		t.Fatalf("GetFuturesBars() error = %v", err)
	}
	if len(got.Result) != 1 {
		t.Fatalf("got %d symbol results, want 1", len(got.Result))
	}
	if got.Result[0].Symbol != "ESZ5" || got.Result[0].InstrumentID != "123456789" {
		t.Errorf("identity = %+v", got.Result[0])
	}
	if len(got.Result[0].Result) != 1 {
		t.Fatalf("got %d bars, want 1", len(got.Result[0].Result))
	}
	bar := got.Result[0].Result[0]
	if bar.Open.Cmp(money.Must(money.NewFromString("4800.00"))) != 0 || bar.Close.Cmp(money.Must(money.NewFromString("4820.00"))) != 0 || bar.Volume != "125000" {
		t.Errorf("bar = %+v", bar)
	}
}

func TestGetFuturesBarsOmitsOptionalParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got, want := q.Get("symbols"), "ESZ5"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_FUTURES"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if _, ok := q["timespan"]; ok {
			t.Errorf("timespan present (%q), want omitted", q.Get("timespan"))
		}
		if _, ok := q["count"]; ok {
			t.Errorf("count present (%q), want omitted", q.Get("count"))
		}
		_, _ = w.Write([]byte(`{"result":[{"symbol":"ESZ5","instrument_id":"1","result":[]}]}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFuturesBars(context.Background(), data.FuturesBarsQuery{
		Symbols: []string{"ESZ5"},
	})
	if err != nil {
		t.Fatalf("GetFuturesBars() error = %v", err)
	}
	if len(got.Result) != 1 || got.Result[0].Symbol != "ESZ5" || len(got.Result[0].Result) != 0 {
		t.Errorf("result = %+v", got)
	}
}

func TestGetFuturesDepth(t *testing.T) {
	t.Parallel()

	const body = `{"symbol":"ESZ5","instrument_id":"123456789","quote_time":1640688000000,` +
		`"asks":[{"price":"4848.50","size":"10","order":[{"mpid":"CME","size":"10"}]}],` +
		`"bids":[{"price":"4848.00","size":"15","order":[{"mpid":"CME","size":"15"}]}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/futures/depths/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbol"), "ESZ5"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_FUTURES"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("depth"), "10"; got != want {
			t.Errorf("depth = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFuturesDepth(context.Background(), data.FuturesDepthQuery{
		Symbol: "ESZ5",
		Depth:  10,
	})
	if err != nil {
		t.Fatalf("GetFuturesDepth() error = %v", err)
	}
	if got.Symbol != "ESZ5" || got.InstrumentID != "123456789" {
		t.Errorf("identity = %+v", got)
	}
	if len(got.Asks) != 1 || len(got.Bids) != 1 {
		t.Fatalf("asks/bids lengths = %d/%d, want 1/1", len(got.Asks), len(got.Bids))
	}
	if got.Asks[0].Price.Cmp(money.Must(money.NewFromString("4848.50"))) != 0 || got.Bids[0].Price.Cmp(money.Must(money.NewFromString("4848.00"))) != 0 {
		t.Errorf("ask/bid = %+v/%+v", got.Asks[0], got.Bids[0])
	}
	if got.QuoteTime != 1640688000000 {
		t.Errorf("quote_time = %d, want 1640688000000", got.QuoteTime)
	}
}

func TestGetFuturesDepthOmitsOptionalParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got, want := q.Get("symbol"), "ESZ5"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_FUTURES"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if _, ok := q["depth"]; ok {
			t.Errorf("depth present (%q), want omitted", q.Get("depth"))
		}
		_, _ = w.Write([]byte(`{"symbol":"ESZ5","instrument_id":"1","quote_time":0,"asks":[],"bids":[]}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	if _, err := c.GetFuturesDepth(context.Background(), data.FuturesDepthQuery{
		Symbol: "ESZ5",
	}); err != nil {
		t.Fatalf("GetFuturesDepth() error = %v", err)
	}
}

func TestGetFuturesFootprint(t *testing.T) {
	t.Parallel()

	const body = `[{"symbol":"ESZ5","instrument_id":"123456789","result":[` +
		`{"time":"2025-09-30T05:47:00.000+0000","trading_session":"RTH",` +
		`"total":"1000","delta":"200","buy_total":"600","sell_total":"400",` +
		`"buy_detail":{"4820.00":"100","4821.00":"60"},` +
		`"sell_detail":{"4820.00":"50","4821.00":"50"}}]}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/futures/footprints/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbols"), "ESZ5"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_FUTURES"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("timespan"), "M1"; got != want {
			t.Errorf("timespan = %q, want %q", got, want)
		}
		if got, want := q.Get("count"), "500"; got != want {
			t.Errorf("count = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFuturesFootprint(context.Background(), data.FuturesFootprintQuery{
		Symbols:  []string{"ESZ5"},
		Timespan: data.FootprintTimespanM1,
		Count:    500,
	})
	if err != nil {
		t.Fatalf("GetFuturesFootprint() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d footprints, want 1", len(got))
	}
	fp := got[0]
	if fp.Symbol != "ESZ5" || fp.InstrumentID != "123456789" {
		t.Errorf("identity = %q/%q, want ESZ5/123456789", fp.Symbol, fp.InstrumentID)
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
	if bar.BuyDetail["4821.00"] != "60" || bar.SellDetail["4820.00"] != "50" {
		t.Errorf("detail maps = buy:%v sell:%v", bar.BuyDetail, bar.SellDetail)
	}
}

func TestGetFuturesFootprintOmitsOptionalParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got, want := q.Get("symbols"), "ESZ5"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("timespan"), "M5"; got != want {
			t.Errorf("timespan = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_FUTURES"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if _, ok := q["count"]; ok {
			t.Errorf("count present (%q), want omitted", q.Get("count"))
		}
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFuturesFootprint(context.Background(), data.FuturesFootprintQuery{
		Symbols:  []string{"ESZ5"},
		Timespan: data.FootprintTimespanM5,
	})
	if err != nil {
		t.Fatalf("GetFuturesFootprint() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("got %d footprints, want 0", len(got))
	}
}
