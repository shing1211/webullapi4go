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

func TestGetSnapshot(t *testing.T) {
	t.Parallel()

	const body = `[{"instrument_id":"913256135","pre_close":"101","change_ratio":"0.05",` +
		`"symbol":"AAPL","last_trade_time":1640688000000,"price":"100","open":"100",` +
		`"close":"101","high":"105","low":"99","volume":"1000","change":"1.0",` +
		`"ask":"13.9","ask_size":"5","bid":"13.8","bid_size":"7","turnover":"0.01",` +
		`"eps":"7.465","eps_ttm":"7.465","lot_size":"1","bps":"4.991",` +
		`"extend_hour_last_price":"100.5","extend_hour_high":"101.0",` +
		`"extend_hour_low":"99.5","extend_hour_change":"0.5",` +
		`"extend_hour_change_ratio":"0.005","extend_hour_volume":"200",` +
		`"extend_hour_last_trade_time":1640688000000,"ovn_price":"100.25",` +
		`"ovn_high":"101.0","ovn_low":"99.5","ovn_volume":"500","ovn_change":"0.25",` +
		`"ovn_change_ratio":"0.0025","ovn_last_trade_time":1640688000000,` +
		`"ovn_ask":"13.9","ovn_ask_size":"5","ovn_bid":"13.8","ovn_bid_size":"7"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/stocks/snapshots/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbols"), "AAPL,TSLA"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("extend_hour_required"), "true"; got != want {
			t.Errorf("extend_hour_required = %q, want %q", got, want)
		}
		if got, want := q.Get("overnight_required"), "true"; got != want {
			t.Errorf("overnight_required = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetSnapshot(context.Background(), data.SnapshotQuery{
		Symbols:            []string{"AAPL", "TSLA"},
		Category:           data.StockCategoryUS,
		ExtendHourRequired: true,
		OvernightRequired:  true,
	})
	if err != nil {
		t.Fatalf("GetSnapshot() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d snapshots, want 1", len(got))
	}
	s := got[0]
	if s.Symbol != "AAPL" || s.InstrumentID != "913256135" || s.Price != "100" {
		t.Errorf("identity/price = %+v", s)
	}
	if s.LastTradeTime != 1640688000000 || s.ExtendHourLastTradeTime != 1640688000000 || s.OvnLastTradeTime != 1640688000000 {
		t.Errorf("timestamps not decoded: %+v", s)
	}
	if s.EPS != "7.465" || s.EPSTTM != "7.465" || s.BPS != "4.991" || s.LotSize != "1" {
		t.Errorf("fundamental fields = %+v", s)
	}
	if s.OvnPrice != "100.25" || s.OvnAskSize != "5" || s.OvnBidSize != "7" {
		t.Errorf("overnight fields = %+v", s)
	}
}
