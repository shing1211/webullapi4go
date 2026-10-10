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

func TestGetQuotes(t *testing.T) {
	t.Parallel()

	const body = `{"symbol":"AAPL","instrument_id":"913256135","quote_time":1640688000000,` +
		`"asks":[{"price":"13.9","size":"5","order":[{"mpid":"NSDQ","size":"5"}],` +
		`"broker":[{"bid":"1","name":"BRK"}]}],` +
		`"bids":[{"price":"13.8","size":"7","order":[{"mpid":"ARCA","size":"7"}]}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/stocks/depths/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbol"), "AAPL"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("depth"), "10"; got != want {
			t.Errorf("depth = %q, want %q", got, want)
		}
		if got, want := q.Get("overnight_required"), "false"; got != want {
			t.Errorf("overnight_required = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetQuotes(context.Background(), data.DepthQuery{
		Symbol:   "AAPL",
		Category: data.StockCategoryUS,
		Depth:    10,
	})
	if err != nil {
		t.Fatalf("GetQuotes() error = %v", err)
	}
	if got.Symbol != "AAPL" || got.InstrumentID != "913256135" || got.QuoteTime != 1640688000000 {
		t.Errorf("identity = %+v", got)
	}
	if len(got.Asks) != 1 || len(got.Bids) != 1 {
		t.Fatalf("asks/bids lengths = %d/%d, want 1/1", len(got.Asks), len(got.Bids))
	}
	ask := got.Asks[0]
	if ask.Price.Cmp(money.Must(money.NewFromString("13.9"))) != 0 || ask.Size != "5" {
		t.Errorf("ask = %+v", ask)
	}
	if len(ask.Order) != 1 || ask.Order[0].MPID != "NSDQ" || ask.Order[0].Size != "5" {
		t.Errorf("ask orders = %+v", ask.Order)
	}
	if len(ask.Broker) != 1 || ask.Broker[0].Bid != "1" || ask.Broker[0].Name != "BRK" {
		t.Errorf("ask brokers = %+v", ask.Broker)
	}
	if got.Bids[0].Order[0].MPID != "ARCA" {
		t.Errorf("bid mpid = %q, want ARCA", got.Bids[0].Order[0].MPID)
	}
}

func TestGetQuotesOmitsDepth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.RawQuery, "category=US_STOCK&overnight_required=false&symbol=AAPL"; got != want {
			t.Errorf("RawQuery = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(`{"symbol":"AAPL","instrument_id":"1","quote_time":0,"asks":[],"bids":[]}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	if _, err := c.GetQuotes(context.Background(), data.DepthQuery{
		Symbol:   "AAPL",
		Category: data.StockCategoryUS,
	}); err != nil {
		t.Fatalf("GetQuotes() error = %v", err)
	}
}
