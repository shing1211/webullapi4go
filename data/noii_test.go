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

func TestGetNOIIBars(t *testing.T) {
	t.Parallel()

	const body = `[{"instrument_id":"913256135","symbol":"AAPL",` +
		`"imbalance_time":1711262998500,"imbalance_ref_price":"172.35",` +
		`"imbalance_near_price":"173.1","imbalance_far_price":"175.5",` +
		`"imbalance_action_type":"PRE_OPEN"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/stocks/noii-bars/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbol"), "AAPL"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("imbalance_action_type"), "PRE_OPEN"; got != want {
			t.Errorf("imbalance_action_type = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetNOIIBars(context.Background(), data.NOIIQuery{
		Symbol:              "AAPL",
		Category:            data.StockCategoryUS,
		ImbalanceActionType: data.NOIIActionPreOpen,
	})
	if err != nil {
		t.Fatalf("GetNOIIBars() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d bars, want 1", len(got))
	}
	bar := got[0]
	if bar.Symbol != "AAPL" || bar.InstrumentID != "913256135" {
		t.Errorf("identity = %q/%q, want AAPL/913256135", bar.Symbol, bar.InstrumentID)
	}
	if bar.ImbalanceTime != 1711262998500 {
		t.Errorf("imbalance_time = %d, want 1711262998500", bar.ImbalanceTime)
	}
	if bar.ImbalanceNearPrice != "173.1" || bar.ImbalanceActionType != data.NOIIActionPreOpen {
		t.Errorf("decoded bar = %+v", bar)
	}
}

func TestGetNOIISnapshot(t *testing.T) {
	t.Parallel()

	const body = `{"instrument_id":"913256135","symbol":"AAPL",` +
		`"paired_shares":"701859","imbalance_shares":"5715","imbalance_side":"2",` +
		`"imbalance_ref_price":"253.83","imbalance_near_price":"253.93",` +
		`"imbalance_far_price":"253.98","imbalance_action_type":"PRE_CLOSE",` +
		`"imbalance_time":1774272599000,"imbalance_var_indicator":"10"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/stocks/noii-snapshots/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("imbalance_action_type"), "PRE_CLOSE"; got != want {
			t.Errorf("imbalance_action_type = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetNOIISnapshot(context.Background(), data.NOIIQuery{
		Symbol:              "AAPL",
		Category:            data.StockCategoryUS,
		ImbalanceActionType: data.NOIIActionPreClose,
	})
	if err != nil {
		t.Fatalf("GetNOIISnapshot() error = %v", err)
	}
	if got.Symbol != "AAPL" || got.ImbalanceShares != "5715" {
		t.Errorf("decoded snapshot = %+v", got)
	}
	if got.ImbalanceTime != 1774272599000 {
		t.Errorf("imbalance_time = %d, want 1774272599000", got.ImbalanceTime)
	}
	if got.ImbalanceActionType != data.NOIIActionPreClose || got.ImbalanceVarIndicator != "10" {
		t.Errorf("decoded snapshot = %+v", got)
	}
}
