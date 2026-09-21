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

package broker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetEventContractCategories(t *testing.T) {
	t.Parallel()

	const body = `[{"category_id":"CAT1","category_name":"Earnings"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/openapi/v1/broker/event-contracts/categories"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestBrokerClient(t, srv.URL)
	got, err := c.GetEventContractCategories(context.Background())
	if err != nil {
		t.Fatalf("GetEventContractCategories() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d categories, want 1", len(got))
	}
	if got[0].CategoryID != "CAT1" || got[0].CategoryName != "Earnings" {
		t.Errorf("category = %+v", got[0])
	}
}

func TestGetEventContractSeries(t *testing.T) {
	t.Parallel()

	const body = `[{"series_id":"SER1","category_id":"CAT1","name":"AAPL Earnings","status":"ACTIVE"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/openapi/v1/broker/event-contracts/series"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("category_id"), "CAT1"; got != want {
			t.Errorf("category_id = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestBrokerClient(t, srv.URL)
	got, err := c.GetEventContractSeries(context.Background(), "CAT1")
	if err != nil {
		t.Fatalf("GetEventContractSeries() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d series, want 1", len(got))
	}
	if got[0].SeriesID != "SER1" || got[0].Name != "AAPL Earnings" {
		t.Errorf("series = %+v", got[0])
	}
}

func TestGetEventContractEvents(t *testing.T) {
	t.Parallel()

	const body = `[{"event_id":"EVT1","series_id":"SER1","name":"Q1 2026","status":"ACTIVE","event_date":"2026-04-15"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/openapi/v1/broker/event-contracts/events"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("series_id"), "SER1"; got != want {
			t.Errorf("series_id = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestBrokerClient(t, srv.URL)
	got, err := c.GetEventContractEvents(context.Background(), "SER1")
	if err != nil {
		t.Fatalf("GetEventContractEvents() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d events, want 1", len(got))
	}
	if got[0].EventID != "EVT1" || got[0].EventDate != "2026-04-15" {
		t.Errorf("event = %+v", got[0])
	}
}

func TestGetEventContractInstruments(t *testing.T) {
	t.Parallel()

	const body = `[{"symbol":"AAPL260415C150","event_id":"EVT1","series_id":"SER1","status":"ACTIVE","strike_price":"150","expiration_date":"2026-04-15"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/openapi/v1/broker/event-contracts/instruments"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("event_id"), "EVT1"; got != want {
			t.Errorf("event_id = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestBrokerClient(t, srv.URL)
	got, err := c.GetEventContractInstruments(context.Background(), "EVT1")
	if err != nil {
		t.Fatalf("GetEventContractInstruments() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d instruments, want 1", len(got))
	}
	if got[0].Symbol != "AAPL260415C150" || got[0].StrikePrice != "150" {
		t.Errorf("instrument = %+v", got[0])
	}
}
