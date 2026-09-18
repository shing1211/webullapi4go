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

// batchBarsRequestBody mirrors the JSON body sent to the batch bars endpoint so
// tests can assert on it.
type batchBarsRequestBody struct {
	Symbols          []string `json:"symbols"`
	Category         string   `json:"category"`
	Timespan         string   `json:"timespan"`
	Count            int      `json:"count"`
	RealTimeRequired *bool    `json:"real_time_required"`
	TradingSessions  string   `json:"trading_sessions"`
	StartTime        int64    `json:"start_time"`
	EndTime          int64    `json:"end_time"`
}

func TestGetBars(t *testing.T) {
	t.Parallel()

	const body = `{"result":[{"symbol":"AAPL","instrument_id":"913256135","result":[` +
		`{"time":"2021-12-28T09:00:09.945+0000","open":"150.25","close":"152.3",` +
		`"high":"153.15","low":"149.8","volume":"1250000"}]}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/stocks/bars/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("Content-Type"), "application/json"; got != want {
			t.Errorf("Content-Type = %q, want %q", got, want)
		}
		var req batchBarsRequestBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		if len(req.Symbols) != 1 || req.Symbols[0] != "AAPL" {
			t.Errorf("symbols = %v, want [AAPL]", req.Symbols)
		}
		if req.Category != "US_STOCK" || req.Timespan != "M1" || req.Count != 5 {
			t.Errorf("category/timespan/count = %q/%q/%d", req.Category, req.Timespan, req.Count)
		}
		if req.RealTimeRequired == nil || *req.RealTimeRequired {
			t.Errorf("real_time_required = %v, want false", req.RealTimeRequired)
		}
		if req.TradingSessions != "RTH" {
			t.Errorf("trading_sessions = %q, want RTH", req.TradingSessions)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	realTime := false
	c := newTestClient(t, srv.URL)
	got, err := c.GetBars(context.Background(), data.BarQuery{
		Symbol:           "AAPL",
		Category:         data.StockCategoryUS,
		Interval:         data.BarTimespanM1,
		Count:            5,
		RealTimeRequired: &realTime,
		TradingSessions:  []data.TradingSession{data.TradingSessionRTH},
	})
	if err != nil {
		t.Fatalf("GetBars() error = %v", err)
	}
	if got.Symbol != "AAPL" || got.InstrumentID != "913256135" {
		t.Errorf("identity = %+v", got)
	}
	if len(got.Result) != 1 {
		t.Fatalf("got %d bars, want 1", len(got.Result))
	}
	bar := got.Result[0]
	if bar.Open != "150.25" || bar.Close != "152.3" || bar.High != "153.15" ||
		bar.Low != "149.8" || bar.Volume != "1250000" {
		t.Errorf("bar = %+v", bar)
	}
}

func TestGetBarsOmitsOptionalParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req batchBarsRequestBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		if len(req.Symbols) != 1 || req.Symbols[0] != "AAPL" || req.Category != "US_STOCK" {
			t.Errorf("symbols/category = %v/%q", req.Symbols, req.Category)
		}
		if req.Timespan != "D" {
			t.Errorf("timespan = %q, want D", req.Timespan)
		}
		if req.Count != 0 || req.RealTimeRequired != nil || req.TradingSessions != "" ||
			req.StartTime != 0 || req.EndTime != 0 {
			t.Errorf("optional fields should be omitted: %+v", req)
		}
		_, _ = w.Write([]byte(`{"result":[{"symbol":"AAPL","instrument_id":"1","result":[]}]}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetBars(context.Background(), data.BarQuery{
		Symbol:   "AAPL",
		Category: data.StockCategoryUS,
		Interval: data.BarTimespanDay,
	})
	if err != nil {
		t.Fatalf("GetBars() error = %v", err)
	}
	if got.Symbol != "AAPL" || len(got.Result) != 0 {
		t.Errorf("result = %+v", got)
	}
}

func TestGetBatchBars(t *testing.T) {
	t.Parallel()

	const body = `{"result":[{"symbol":"AAPL","instrument_id":"913256135","result":[` +
		`{"time":"2021-12-28T09:00:09.945+0000","open":"150.25","close":"152.3",` +
		`"high":"153.15","low":"149.8","volume":"1250000"}]},` +
		`{"symbol":"TSLA","instrument_id":"913256136","result":[]}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/stocks/bars/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("Content-Type"), "application/json"; got != want {
			t.Errorf("Content-Type = %q, want %q", got, want)
		}
		var req batchBarsRequestBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		if len(req.Symbols) != 2 || req.Symbols[0] != "AAPL" || req.Symbols[1] != "TSLA" {
			t.Errorf("symbols = %v, want [AAPL TSLA]", req.Symbols)
		}
		if req.Category != "US_STOCK" || req.Timespan != "M1" || req.Count != 500 {
			t.Errorf("category/timespan/count = %q/%q/%d", req.Category, req.Timespan, req.Count)
		}
		if req.RealTimeRequired == nil || *req.RealTimeRequired {
			t.Errorf("real_time_required = %v, want false", req.RealTimeRequired)
		}
		if req.TradingSessions != "PRE,RTH" {
			t.Errorf("trading_sessions = %q, want PRE,RTH", req.TradingSessions)
		}
		if req.StartTime != 1711262998500 || req.EndTime != 1711349398500 {
			t.Errorf("start/end = %d/%d", req.StartTime, req.EndTime)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	realTime := false
	c := newTestClient(t, srv.URL)
	got, err := c.GetBatchBars(context.Background(), data.BatchBarQuery{
		Symbols:          []string{"AAPL", "TSLA"},
		Category:         data.StockCategoryUS,
		Timespan:         data.BarTimespanM1,
		Count:            500,
		RealTimeRequired: &realTime,
		TradingSessions:  []data.TradingSession{data.TradingSessionPre, data.TradingSessionRTH},
		StartTime:        1711262998500,
		EndTime:          1711349398500,
	})
	if err != nil {
		t.Fatalf("GetBatchBars() error = %v", err)
	}
	if len(got.Result) != 2 {
		t.Fatalf("got %d symbol results, want 2", len(got.Result))
	}
	if got.Result[0].Symbol != "AAPL" || got.Result[0].InstrumentID != "913256135" {
		t.Errorf("first result identity = %+v", got.Result[0])
	}
	if len(got.Result[0].Result) != 1 || got.Result[0].Result[0].Close != "152.3" {
		t.Errorf("first result bars = %+v", got.Result[0].Result)
	}
	if got.Result[1].Symbol != "TSLA" || len(got.Result[1].Result) != 0 {
		t.Errorf("second result = %+v", got.Result[1])
	}
}
