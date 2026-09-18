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

// decodeRequestBody reads r's JSON body into out, failing the test on error.
func decodeRequestBody(t *testing.T, r *http.Request, out any) {
	t.Helper()
	defer func() { _ = r.Body.Close() }()
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		t.Fatalf("decoding request body: %v", err)
	}
}

func TestGetWatchlists(t *testing.T) {
	t.Parallel()

	const body = `[{"watchlist_id":"12345678","name":"My HK Stocks","sort":1,` +
		`"create_time":"2026-04-01T10:30:00Z","update_time":"2026-04-30T15:45:00Z"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/watchlists/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetWatchlists(context.Background())
	if err != nil {
		t.Fatalf("GetWatchlists() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d watchlists, want 1", len(got))
	}
	if got[0].WatchlistID != "12345678" || got[0].Name != "My HK Stocks" || got[0].Sort != 1 {
		t.Errorf("watchlist = %+v", got[0])
	}
	if got[0].UpdateTime != "2026-04-30T15:45:00Z" {
		t.Errorf("update_time = %q", got[0].UpdateTime)
	}
}

func TestCreateWatchlist(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/watchlists/create"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		var req struct {
			Name string `json:"name"`
			Sort int32  `json:"sort"`
		}
		decodeRequestBody(t, r, &req)
		if req.Name != "My HK Stocks" || req.Sort != 3 {
			t.Errorf("body = %+v", req)
		}
		_, _ = w.Write([]byte(`{"watchlist_id":"98765432"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.CreateWatchlist(context.Background(), data.CreateWatchlistParams{
		Name: "My HK Stocks",
		Sort: 3,
	})
	if err != nil {
		t.Fatalf("CreateWatchlist() error = %v", err)
	}
	if got.WatchlistID != "98765432" {
		t.Errorf("watchlist_id = %q, want 98765432", got.WatchlistID)
	}
}

func TestUpdateWatchlist(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/watchlists/update"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		var req struct {
			WatchlistID string `json:"watchlist_id"`
			Name        string `json:"name"`
			Sort        int32  `json:"sort"`
		}
		decodeRequestBody(t, r, &req)
		if req.WatchlistID != "12345678" || req.Name != "Renamed" || req.Sort != 2 {
			t.Errorf("body = %+v", req)
		}
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.UpdateWatchlist(context.Background(), data.UpdateWatchlistParams{
		WatchlistID: "12345678",
		Name:        "Renamed",
		Sort:        2,
	})
	if err != nil {
		t.Fatalf("UpdateWatchlist() error = %v", err)
	}
	if !got.Success {
		t.Errorf("success = false, want true")
	}
}

func TestDeleteWatchlist(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/watchlists/delete"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		var req struct {
			WatchlistID string `json:"watchlist_id"`
		}
		decodeRequestBody(t, r, &req)
		if req.WatchlistID != "12345678" {
			t.Errorf("watchlist_id = %q, want 12345678", req.WatchlistID)
		}
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.DeleteWatchlist(context.Background(), "12345678")
	if err != nil {
		t.Fatalf("DeleteWatchlist() error = %v", err)
	}
	if !got.Success {
		t.Errorf("success = false, want true")
	}
}

func TestGetWatchlistInstruments(t *testing.T) {
	t.Parallel()

	const body = `{"watchlist_id":"12345678","instruments":[` +
		`{"instrument_id":"913256135","symbol":"AAPL","name":"Apple Inc.",` +
		`"exchange_code":"NSQ","sort":1,"added_time":"2024-05-11T01:38:36.839Z"}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/watchlists/instruments/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("watchlist_id"), "12345678"; got != want {
			t.Errorf("watchlist_id = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetWatchlistInstruments(context.Background(), "12345678")
	if err != nil {
		t.Fatalf("GetWatchlistInstruments() error = %v", err)
	}
	if got.WatchlistID != "12345678" || len(got.Instruments) != 1 {
		t.Fatalf("result = %+v", got)
	}
	inst := got.Instruments[0]
	if inst.Symbol != "AAPL" || inst.InstrumentID != "913256135" || inst.ExchangeCode != "NSQ" {
		t.Errorf("instrument = %+v", inst)
	}
	if inst.Name != "Apple Inc." || inst.Sort != 1 || inst.AddedTime != "2024-05-11T01:38:36.839Z" {
		t.Errorf("instrument = %+v", inst)
	}
}

// watchlistInstrumentPaths maps the three instrument mutation methods to their
// documented endpoints.
var watchlistInstrumentPaths = map[string]string{
	"add":    "/market-data/watchlists/instruments/add",
	"remove": "/market-data/watchlists/instruments/remove",
	"update": "/market-data/watchlists/instruments/update",
}

func TestWatchlistInstrumentMutations(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		call func(*data.Client, context.Context, data.WatchlistInstrumentsParam) (*data.SuccessResponse, error)
	}{
		{"add", (*data.Client).AddWatchlistInstruments},
		{"remove", (*data.Client).RemoveWatchlistInstruments},
		{"update", (*data.Client).UpdateWatchlistInstruments},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("method = %q, want POST", r.Method)
				}
				if got, want := r.URL.Path, watchlistInstrumentPaths[tc.name]; got != want {
					t.Errorf("path = %q, want %q", got, want)
				}
				var req struct {
					WatchlistID string `json:"watchlist_id"`
					Instruments []struct {
						Symbol   string `json:"symbol"`
						Category string `json:"category"`
						Sort     int32  `json:"sort"`
					} `json:"instruments"`
				}
				decodeRequestBody(t, r, &req)
				if req.WatchlistID != "12345678" || len(req.Instruments) != 1 {
					t.Fatalf("body = %+v", req)
				}
				inst := req.Instruments[0]
				if inst.Symbol != "AAPL" || inst.Category != "US_STOCK" || inst.Sort != 4 {
					t.Errorf("instrument = %+v", inst)
				}
				_, _ = w.Write([]byte(`{"success":true}`))
			}))
			defer srv.Close()

			c := newTestClient(t, srv.URL)
			got, err := tc.call(c, context.Background(), data.WatchlistInstrumentsParam{
				WatchlistID: "12345678",
				Instruments: []data.WatchlistInstrumentParam{
					{Symbol: "AAPL", Category: data.StockCategoryUS, Sort: 4},
				},
			})
			if err != nil {
				t.Fatalf("%s() error = %v", tc.name, err)
			}
			if !got.Success {
				t.Errorf("success = false, want true")
			}
		})
	}
}
