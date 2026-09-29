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

// Reading the cursor is half of pagination; sending it is the other half, and the
// second half is the one that fails quietly. A method can decode the envelope
// correctly, populate PaginationKey, and still drop the argument on the way out, and
// then every call returns page one forever with no error anywhere.
//
// That is not hypothetical. The first pass of the v2.1.35 change built the query in
// GetMarketSectors, wrote the cursor into it, and then passed nil to c.get -- so the
// cursor was constructed and discarded. Nothing in the tree caught it except a compile
// error about an undefined variable, which does not distinguish a dropped query from
// any other mistake. These tests exist so the next one is caught by an assertion.

const testCursor = "eyJ2IjoxLCJsYXN0SWQiOiIwIiwicGFnZU9mZnNldCI6MX0="

// pageBody is a documented one-row page carrying testCursor.
const pageBody = `{"data":[{"sector_name":"Technology"}],` +
	`"pagination_key":"` + testCursor + `"}`

// TestEnvelopeCursorIsSent covers every converted data method in one table. Each case
// calls the method with a non-empty cursor and asserts the server saw it.
func TestEnvelopeCursorIsSent(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		call func(c *data.Client) error
	}{
		{
			name: "GetMarketSectors",
			call: func(c *data.Client) error {
				_, err := c.GetMarketSectors(context.Background(), testCursor)
				return err
			},
		},
		{
			name: "GetMarketSectorDetail",
			call: func(c *data.Client) error {
				_, err := c.GetMarketSectorDetail(context.Background(), data.MarketSectorDetailQuery{
					SectorName:    "Technology",
					Category:      data.StockCategoryUS,
					PaginationKey: testCursor,
				})
				return err
			},
		},
		{
			name: "GetFundDividends",
			call: func(c *data.Client) error {
				_, err := c.GetFundDividends(context.Background(), data.FundDividendsQuery{
					Symbol:        "SPY",
					PaginationKey: testCursor,
				})
				return err
			},
		},
		{
			name: "GetEventContractSeries",
			call: func(c *data.Client) error {
				_, err := c.GetEventContractSeries(context.Background(), data.EventContractSeriesQuery{
					Category:      "PREDICTION",
					PaginationKey: testCursor,
				})
				return err
			},
		},
		{
			name: "GetEventContractMarkets",
			call: func(c *data.Client) error {
				_, err := c.GetEventContractMarkets(context.Background(), data.EventContractMarketsQuery{
					SeriesSymbol:  "AAPL_PREDICTION",
					PaginationKey: testCursor,
				})
				return err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var seen string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				seen = r.URL.Query().Get("pagination_key")
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(pageBody))
			}))
			defer srv.Close()

			if err := tc.call(newTestClient(t, srv.URL)); err != nil {
				t.Fatalf("call error = %v", err)
			}
			if seen != testCursor {
				t.Errorf("server saw pagination_key = %q, want %q: the cursor was decoded "+
					"but never sent, so this endpoint would return page one forever",
					seen, testCursor)
			}
		})
	}
}

// TestEnvelopeOmitCursorOnFirstPage pins the other direction: a caller asking for the
// first page passes an empty cursor, and an empty cursor must not be sent as an empty
// query parameter. A server that treats "" as a real cursor could answer with an error
// or with a wrong page, and neither is visible from the SDK.
func TestEnvelopeOmitCursorOnFirstPage(t *testing.T) {
	t.Parallel()

	var seen string
	var present bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, present = r.URL.Query()["pagination_key"]
		seen = r.URL.Query().Get("pagination_key")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[],"pagination_key":""}`))
	}))
	defer srv.Close()

	if _, err := newTestClient(t, srv.URL).GetMarketSectors(context.Background(), ""); err != nil {
		t.Fatalf("GetMarketSectors() error = %v", err)
	}
	if present {
		t.Errorf("pagination_key was sent as %q on a first-page request, want the "+
			"parameter omitted entirely", seen)
	}
}

// TestEnvelopeRoundTripsAcrossTwoPages is the end-to-end shape of the capability: read
// a page, take its cursor, and send it to get the next one. This is the loop a caller
// has to write, so it is asserted here once rather than left to each of the fourteen
// methods.
func TestEnvelopeRoundTripsAcrossTwoPages(t *testing.T) {
	t.Parallel()

	page := 0
	var cursors []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cursors = append(cursors, r.URL.Query().Get("pagination_key"))
		w.Header().Set("Content-Type", "application/json")
		page++
		if page == 1 {
			_, _ = w.Write([]byte(pageBody))
			return
		}
		// Last page: no cursor, which is how the caller knows to stop.
		_, _ = w.Write([]byte(`{"data":[{"sector_name":"Energy"}],` +
			`"pagination_key":""}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	var symbols []string
	cursor := ""
	for range 2 {
		pg, err := c.GetMarketSectors(context.Background(), cursor)
		if err != nil {
			t.Fatalf("GetMarketSectors() error = %v", err)
		}
		for _, s := range pg.Data {
			symbols = append(symbols, s.SectorName)
		}
		cursor = pg.PaginationKey
		if cursor == "" {
			break
		}
	}

	if len(symbols) != 2 || symbols[0] != "Technology" || symbols[1] != "Energy" {
		t.Errorf("collected %v, want [Technology Energy] across the two pages", symbols)
	}
	if page != 2 {
		t.Errorf("server served %d pages, want 2: an empty cursor must end the loop", page)
	}
	if len(cursors) != 2 || cursors[0] != "" || cursors[1] != testCursor {
		t.Errorf("cursors sent = %v, want [\"\" %q]: the second request must carry the "+
			"first page's cursor", cursors, testCursor)
	}
}
