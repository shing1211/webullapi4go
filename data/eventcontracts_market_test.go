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

func TestGetEventSnapshot(t *testing.T) {
	t.Parallel()

	const body = `[{"symbol":"AAPL_250919C00240000","last_price":"3.50",` +
		`"yes_bid":"3.45","yes_ask":"3.55","volume":"1200",` +
		`"open_interest":"5000","timestamp":"2025-09-19T10:00:00Z"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/event-contracts/snapshots/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbols"), "AAPL_250919C00240000,TSLA_250919C00250000"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetEventSnapshot(context.Background(), data.EventSnapshotQuery{
		Symbols:  []string{"AAPL_250919C00240000", "TSLA_250919C00250000"},
		Category: "US_STOCK",
	})
	if err != nil {
		t.Fatalf("GetEventSnapshot() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d snapshots, want 1", len(got))
	}
	s := got[0]
	if s.Symbol != "AAPL_250919C00240000" {
		t.Errorf("symbol = %q, want AAPL_250919C00240000", s.Symbol)
	}
	if s.LastPrice.Cmp(money.Must(money.NewFromString("3.50"))) != 0 || s.YesBid.Cmp(money.Must(money.NewFromString("3.45"))) != 0 || s.YesAsk.Cmp(money.Must(money.NewFromString("3.55"))) != 0 {
		t.Errorf("prices = last=%s bid=%s ask=%s", s.LastPrice, s.YesBid, s.YesAsk)
	}
	if s.Volume != "1200" || s.OpenInterest != "5000" {
		t.Errorf("volume/open_interest = %s/%s", s.Volume, s.OpenInterest)
	}
	if s.Timestamp != "2025-09-19T10:00:00Z" {
		t.Errorf("timestamp = %q", s.Timestamp)
	}
}

func TestGetEventDepth(t *testing.T) {
	t.Parallel()

	// The documented 200 body is an array of objects, each requiring
	// instrument_id, symbol, quote_time, yes_bids, yes_asks, no_bids and
	// no_asks. The no-side book is what this change makes readable.
	const body = `[{"instrument_id":"504279491","symbol":"AAPL_250919C00240000",` +
		`"quote_time":1768872168870,` +
		`"yes_bids":[{"price":"3.45","size":"100"},{"price":"3.40","size":"200"}],` +
		`"yes_asks":[{"price":"3.55","size":"150"},{"price":"3.60","size":"250"}],` +
		`"no_bids":[{"price":"0.10","size":"300"}],` +
		`"no_asks":[{"price":"0.12","size":"400"}]}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/event-contracts/depths/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbol"), "AAPL_250919C00240000"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_STOCK"; got != want {
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
	got, err := c.GetEventDepth(context.Background(), data.EventDepthQuery{
		Symbol:   "AAPL_250919C00240000",
		Category: "US_STOCK",
		Depth:    10,
	})
	if err != nil {
		t.Fatalf("GetEventDepth() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1 element per instrument", len(got))
	}
	book := got[0]
	if book.Symbol != "AAPL_250919C00240000" {
		t.Errorf("symbol = %q", book.Symbol)
	}
	if book.InstrumentID != "504279491" {
		t.Errorf("InstrumentID = %q, want 504279491", book.InstrumentID)
	}
	if book.QuoteTime != 1768872168870 {
		t.Errorf("quote_time = %d, want 1768872168870", book.QuoteTime)
	}
	if len(book.YesBids) != 2 || len(book.YesAsks) != 2 {
		t.Fatalf("bids/asks lengths = %d/%d, want 2/2", len(book.YesBids), len(book.YesAsks))
	}
	// The no-side book is new: the SDK could not read the no side of an event
	// contract's depth at all before this change.
	if len(book.NoBids) != 1 || len(book.NoAsks) != 1 {
		t.Fatalf("no-bids/no-asks lengths = %d/%d, want 1/1", len(book.NoBids), len(book.NoAsks))
	}
	if book.NoBids[0].Price.Cmp(money.Must(money.NewFromString("0.10"))) != 0 ||
		book.NoBids[0].Size != "300" {
		t.Errorf("first no bid = %+v", book.NoBids[0])
	}
	if book.YesBids[0].Price.Cmp(money.Must(money.NewFromString("3.45"))) != 0 || book.YesBids[0].Size != "100" {
		t.Errorf("first bid = %+v", book.YesBids[0])
	}
	if book.YesAsks[0].Price.Cmp(money.Must(money.NewFromString("3.55"))) != 0 || book.YesAsks[0].Size != "150" {
		t.Errorf("first ask = %+v", book.YesAsks[0])
	}
}

func TestGetEventBars(t *testing.T) {
	t.Parallel()

	// The documented 200 body is an array of objects carrying instrument_id,
	// symbol and result, where result is the bar array. The SDK used to
	// decode the inner array and so could not read the grouping key.
	const body = `[{"instrument_id":"504279491","symbol":"AAPL_250919C00240000",` +
		`"result":[{"open":"3.40","high":"3.60","low":"3.35","close":"3.50",` +
		`"volume":"5000","time":"2025-09-19T10:00:00Z",` +
		`"timespan":"M1"}]}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/event-contracts/bars/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbols"), "AAPL_250919C00240000"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("timespan"), "M1"; got != want {
			t.Errorf("timespan = %q, want %q", got, want)
		}
		if got, want := q.Get("count"), "100"; got != want {
			t.Errorf("count = %q, want %q", got, want)
		}
		if got, want := q.Get("real_time_required"), "false"; got != want {
			t.Errorf("real_time_required = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	realTime := false
	c := newTestClient(t, srv.URL)
	got, err := c.GetEventBars(context.Background(), data.EventBarsQuery{
		Symbols:          []string{"AAPL_250919C00240000"},
		Category:         "US_STOCK",
		Timespan:         data.BarTimespanM1,
		Count:            100,
		RealTimeRequired: &realTime,
	})
	if err != nil {
		t.Fatalf("GetEventBars() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1 element per instrument", len(got))
	}
	if got[0].InstrumentID != "504279491" {
		t.Errorf("InstrumentID = %q, want 504279491: the grouping key is the whole "+
			"point of the documented shape", got[0].InstrumentID)
	}
	if got[0].Symbol != "AAPL_250919C00240000" {
		t.Errorf("Symbol = %q", got[0].Symbol)
	}
	bars := got[0].Result
	if len(bars) != 1 {
		t.Fatalf("len(Result) = %d, want 1", len(bars))
	}
	bar := bars[0]
	// Symbol is the SDK's own spelling and the page does not require it on the inner
	// bar, so the documented body leaves it empty rather than borrowing the wrapper's.
	if bar.Symbol != "" {
		t.Errorf("bar.Symbol = %q, want empty: the inner bar carries no documented "+
			"symbol, so a populated value could only have come from the wrapper", bar.Symbol)
	}
	if bar.Open.Cmp(money.Must(money.NewFromString("3.40"))) != 0 || bar.High.Cmp(money.Must(money.NewFromString("3.60"))) != 0 || bar.Low.Cmp(money.Must(money.NewFromString("3.35"))) != 0 || bar.Close.Cmp(money.Must(money.NewFromString("3.50"))) != 0 {
		t.Errorf("OHLC = %s/%s/%s/%s", bar.Open, bar.High, bar.Low, bar.Close)
	}
	if bar.Volume != "5000" || bar.Time != "2025-09-19T10:00:00Z" {
		t.Errorf("volume/time = %s/%s", bar.Volume, bar.Time)
	}
	if bar.Timespan != data.BarTimespanM1 {
		t.Errorf("timespan = %q, want M1", bar.Timespan)
	}
}

func TestGetEventTick(t *testing.T) {
	t.Parallel()

	// The documented 200 body wraps the tick series in an object carrying
	// instrument_id, symbol and result, so the two ticks belong in one wrapper.
	const body = `[{"instrument_id":"504279491","symbol":"AAPL_250919C00240000",` +
		`"result":[{"yes_price":"3.50","no_price":"0.50","side":"B","volume":"10",` +
		`"trade_id":"T123","time":"2025-09-19T10:00:00Z"},` +
		`{"yes_price":"3.48","no_price":"0.52","side":"S","volume":"5",` +
		`"trade_id":"T124","time":"2025-09-19T10:00:01Z"}]}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/event-contracts/ticks/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbol"), "AAPL_250919C00240000"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_STOCK"; got != want {
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
	got, err := c.GetEventTick(context.Background(), data.EventTickQuery{
		Symbol:   "AAPL_250919C00240000",
		Category: "US_STOCK",
		Count:    30,
	})
	if err != nil {
		t.Fatalf("GetEventTick() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d elements, want 1 element per instrument", len(got))
	}
	if got[0].InstrumentID != "504279491" {
		t.Errorf("InstrumentID = %q, want 504279491: the grouping key is the whole "+
			"point of the documented shape", got[0].InstrumentID)
	}
	if got[0].Symbol != "AAPL_250919C00240000" {
		t.Errorf("Symbol = %q", got[0].Symbol)
	}
	ticks := got[0].Result
	if len(ticks) != 2 {
		t.Fatalf("got %d ticks, want 2", len(ticks))
	}
	tick := ticks[0]
	if tick.YesPrice.Cmp(money.Must(money.NewFromString("3.50"))) != 0 || tick.NoPrice.Cmp(money.Must(money.NewFromString("0.50"))) != 0 {
		t.Errorf("yes/no price = %s/%s", tick.YesPrice, tick.NoPrice)
	}
	if tick.Side != "B" || tick.Volume != "10" || tick.TradeID != "T123" {
		t.Errorf("side/volume/trade_id = %s/%s/%s", tick.Side, tick.Volume, tick.TradeID)
	}
	if tick.Time != "2025-09-19T10:00:00Z" {
		t.Errorf("time = %q", tick.Time)
	}
	if ticks[1].Side != "S" {
		t.Errorf("second tick side = %q, want S", ticks[1].Side)
	}
}

func TestGetEventSnapshotOmitsOptionalParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got, want := q.Get("symbols"), "AAPL_250919C00240000"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(`[{"symbol":"AAPL_250919C00240000","last_price":"1.00",` +
			`"yes_bid":"0.99","yes_ask":"1.01","volume":"0","open_interest":"0","timestamp":"t"}]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetEventSnapshot(context.Background(), data.EventSnapshotQuery{
		Symbols:  []string{"AAPL_250919C00240000"},
		Category: "US_STOCK",
	})
	if err != nil {
		t.Fatalf("GetEventSnapshot() error = %v", err)
	}
}

func TestGetEventDepthOmitsOptionalParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got, want := q.Get("symbol"), "AAPL_250919C00240000"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if q.Get("depth") != "" {
			t.Errorf("depth should be omitted, got %q", q.Get("depth"))
		}
		_, _ = w.Write([]byte(`[{"instrument_id":"504279491",` +
			`"symbol":"AAPL_250919C00240000","quote_time":1768872168870,` +
			`"yes_bids":[],"yes_asks":[],"no_bids":[],"no_asks":[]}]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetEventDepth(context.Background(), data.EventDepthQuery{
		Symbol:   "AAPL_250919C00240000",
		Category: "US_STOCK",
	})
	if err != nil {
		t.Fatalf("GetEventDepth() error = %v", err)
	}
}

func TestGetEventBarsOmitsOptionalParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got, want := q.Get("symbols"), "AAPL_250919C00240000"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("timespan"), "D"; got != want {
			t.Errorf("timespan = %q, want %q", got, want)
		}
		if q.Get("count") != "" {
			t.Errorf("count should be omitted, got %q", q.Get("count"))
		}
		if q.Get("real_time_required") != "" {
			t.Errorf("real_time_required should be omitted, got %q", q.Get("real_time_required"))
		}
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetEventBars(context.Background(), data.EventBarsQuery{
		Symbols:  []string{"AAPL_250919C00240000"},
		Category: "US_STOCK",
		Timespan: data.BarTimespanDay,
	})
	if err != nil {
		t.Fatalf("GetEventBars() error = %v", err)
	}
}

func TestGetEventTickOmitsOptionalParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got, want := q.Get("symbol"), "AAPL_250919C00240000"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if q.Get("count") != "" {
			t.Errorf("count should be omitted, got %q", q.Get("count"))
		}
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetEventTick(context.Background(), data.EventTickQuery{
		Symbol:   "AAPL_250919C00240000",
		Category: "US_STOCK",
	})
	if err != nil {
		t.Fatalf("GetEventTick() error = %v", err)
	}
}
