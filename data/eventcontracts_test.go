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

func TestGetEventContractCategories(t *testing.T) {
	t.Parallel()

	const body = `[{"category":"PREDICTION","name":"Predictions"},{"category":"POLITICS","name":"Politics"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/trading/instruments/event-contracts/categories/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if q := r.URL.Query(); len(q) != 0 {
			t.Errorf("query = %v, want empty", q)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetEventContractCategories(context.Background())
	if err != nil {
		t.Fatalf("GetEventContractCategories() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Category != "PREDICTION" || got[0].Name != "Predictions" {
		t.Errorf("category[0] = %+v", got[0])
	}
	if got[1].Category != "POLITICS" || got[1].Name != "Politics" {
		t.Errorf("category[1] = %+v", got[1])
	}
}

func TestGetEventContractSeries(t *testing.T) {
	t.Parallel()

	const body = `[{"series_symbol":"AAPL_PREDICTION","category":"PREDICTION",` +
		`"name":"AAPL Price Prediction","status":"ACTIVE"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/trading/instruments/event-contracts/series/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("category"), "PREDICTION"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("symbols"), "AAPL,MSFT"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("pagination_key"), "next_page"; got != want {
			t.Errorf("pagination_key = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetEventContractSeries(context.Background(), data.EventContractSeriesQuery{
		Category:      "PREDICTION",
		Symbols:       []string{"AAPL", "MSFT"},
		PaginationKey: "next_page",
	})
	if err != nil {
		t.Fatalf("GetEventContractSeries() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].SeriesSymbol != "AAPL_PREDICTION" || got[0].Status != "ACTIVE" {
		t.Errorf("series[0] = %+v", got[0])
	}
}

func TestGetEventContractEvents(t *testing.T) {
	t.Parallel()

	const body = `[{"event_symbol":"AAPL_PREDICTION_2026Q4",` +
		`"series_symbol":"AAPL_PREDICTION","name":"Q4 2026 Prediction","status":"OPEN"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/trading/instruments/event-contracts/events/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("series_symbol"), "AAPL_PREDICTION"; got != want {
			t.Errorf("series_symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("symbols"), "AAPL"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("status"), "OPEN"; got != want {
			t.Errorf("status = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetEventContractEvents(context.Background(), data.EventContractEventsQuery{
		SeriesSymbol: "AAPL_PREDICTION",
		Symbols:      []string{"AAPL"},
		Status:       "OPEN",
	})
	if err != nil {
		t.Fatalf("GetEventContractEvents() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].EventSymbol != "AAPL_PREDICTION_2026Q4" || got[0].Status != "OPEN" {
		t.Errorf("event[0] = %+v", got[0])
	}
}

func TestGetEventContractMarkets(t *testing.T) {
	t.Parallel()

	const body = `[{"symbol":"AAPL_P450","event_symbol":"AAPL_PREDICTION_2026Q4",` +
		`"series_symbol":"AAPL_PREDICTION","status":"TRADEABLE",` +
		`"strike_price":"450.00","expiration_date":"2026-12-31"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/trading/instruments/event-contracts/markets/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("series_symbol"), "AAPL_PREDICTION"; got != want {
			t.Errorf("series_symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("event_symbol"), "AAPL_PREDICTION_2026Q4"; got != want {
			t.Errorf("event_symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("symbols"), "AAPL"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("expiration_date_after"), "2026-10-01"; got != want {
			t.Errorf("expiration_date_after = %q, want %q", got, want)
		}
		if got, want := q.Get("pagination_key"), "cursor123"; got != want {
			t.Errorf("pagination_key = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetEventContractMarkets(context.Background(), data.EventContractMarketsQuery{
		SeriesSymbol:        "AAPL_PREDICTION",
		EventSymbol:         "AAPL_PREDICTION_2026Q4",
		Symbols:             []string{"AAPL"},
		ExpirationDateAfter: "2026-10-01",
		PaginationKey:       "cursor123",
	})
	if err != nil {
		t.Fatalf("GetEventContractMarkets() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].Symbol != "AAPL_P450" || got[0].Status != "TRADEABLE" {
		t.Errorf("market[0] = %+v", got[0])
	}
	if got[0].StrikePrice.Cmp(money.Must(money.NewFromString("450.00"))) != 0 || got[0].ExpirationDate != "2026-12-31" {
		t.Errorf("strike/expiration = %q/%q", got[0].StrikePrice, got[0].ExpirationDate)
	}
}
