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

package brokerfd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shing1211/webullapi4go/client"
)

// Reading the cursor is half of pagination; sending it is the other half, and the
// second half fails quietly. A method can decode the envelope, populate
// PaginationKey, and then drop the argument on the way out, and every call returns
// page one forever with nothing reporting an error.
//
// All seven brokerfd methods take the cursor as a positional parameter rather than
// through a query struct, which is exactly the shape in which an accepted-and-ignored
// argument is easy to write and invisible to review. These tests assert the server saw
// it.

const fdTestCursor = "eyJ2IjoxLCJsYXN0SWQiOiIwIiwicGFnZU9mZnNldCI6MX0="

// TestFDCursorIsSent covers all seven converted methods in one table.
func TestFDCursorIsSent(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		call func(c *Client) error
	}{
		{"ListFDAccounts", func(c *Client) error {
			_, err := c.ListFDAccounts(context.Background(), fdTestCursor)
			return err
		}},
		{"ListFDTransfers", func(c *Client) error {
			_, err := c.ListFDTransfers(context.Background(), "A1", fdTestCursor)
			return err
		}},
		{"GetFDActivities", func(c *Client) error {
			_, err := c.GetFDActivities(context.Background(), "A1", fdTestCursor)
			return err
		}},
		{"GetFDECInstruments", func(c *Client) error {
			_, err := c.GetFDECInstruments(context.Background(), "EVT1", fdTestCursor)
			return err
		}},
		{"GetFDOpenOrders", func(c *Client) error {
			_, err := c.GetFDOpenOrders(context.Background(), "A1", fdTestCursor)
			return err
		}},
		{"GetFDOrderHistory", func(c *Client) error {
			_, err := c.GetFDOrderHistory(context.Background(), "A1", fdTestCursor)
			return err
		}},
		{"GetFDStockInstruments", func(c *Client) error {
			_, err := c.GetFDStockInstruments(context.Background(), []string{"AAPL"}, fdTestCursor)
			return err
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var seen string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				seen = r.URL.Query().Get("pagination_key")
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]any{
					"data":           []map[string]any{{"account_id": "A1"}},
					"pagination_key": fdTestCursor,
				})
			}))
			defer srv.Close()

			cl, err := client.New(
				client.WithAppKey(testAppKey), client.WithAppSecret(testAppSecret),
				client.WithEndpoints(client.Endpoints{HTTP: srv.URL, BrokerHTTP: srv.URL}))
			if err != nil {
				t.Fatalf("client.New error = %v", err)
			}
			if err := tc.call(New(cl)); err != nil {
				t.Fatalf("call error = %v", err)
			}
			if seen != fdTestCursor {
				t.Errorf("server saw pagination_key = %q, want %q: the cursor was accepted "+
					"as an argument and then dropped, so this endpoint would return page "+
					"one forever", seen, fdTestCursor)
			}
		})
	}
}

// TestFDOmitCursorOnFirstPage pins the other direction: an empty cursor must be
// omitted rather than sent as an empty query parameter.
func TestFDOmitCursorOnFirstPage(t *testing.T) {
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

	cl, err := client.New(
		client.WithAppKey(testAppKey), client.WithAppSecret(testAppSecret),
		client.WithEndpoints(client.Endpoints{HTTP: srv.URL, BrokerHTTP: srv.URL}))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	if _, err := New(cl).ListFDAccounts(context.Background(), ""); err != nil {
		t.Fatalf("ListFDAccounts error = %v", err)
	}
	if present {
		t.Errorf("pagination_key was sent as %q on a first-page request, want the "+
			"parameter omitted entirely", seen)
	}
}
