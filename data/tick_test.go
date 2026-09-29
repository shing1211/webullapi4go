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
	"strings"
	"testing"

	"github.com/shing1211/webullapi4go/data"
	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

func TestGetTick(t *testing.T) {
	t.Parallel()

	const body = `{"symbol":"AAPL","instrument_id":"913256135","result":[` +
		`{"time":"1761182953043","price":"48.07","volume":"1","side":"S"},` +
		`{"time":"1761182953042","price":"48.06","volume":"3","side":"B"}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/stocks/ticks/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbol"), "AAPL"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("count"), "30"; got != want {
			t.Errorf("count = %q, want %q", got, want)
		}
		if got, want := q.Get("trading_sessions"), "PRE,RTH"; got != want {
			t.Errorf("trading_sessions = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetTick(context.Background(), data.TickQuery{
		Symbol:          "AAPL",
		Category:        data.StockCategoryUS,
		Count:           30,
		TradingSessions: []data.TradingSession{data.TradingSessionPre, data.TradingSessionRTH},
	})
	if err != nil {
		t.Fatalf("GetTick() error = %v", err)
	}
	if got.Symbol != "AAPL" || got.InstrumentID != "913256135" {
		t.Errorf("identity = %+v", got)
	}
	if len(got.Result) != 2 {
		t.Fatalf("got %d ticks, want 2", len(got.Result))
	}
	if got.Result[0].Time != "1761182953043" || got.Result[0].Price.Cmp(money.Must(money.NewFromString("48.07"))) != 0 ||
		got.Result[0].Volume != "1" || got.Result[0].Side != "S" {
		t.Errorf("first tick = %+v", got.Result[0])
	}
	if got.Result[1].Side != "B" {
		t.Errorf("second tick side = %q, want B", got.Result[1].Side)
	}
}

func TestGetTickOmitsOptionalParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.RawQuery, "category=US_STOCK&symbol=AAPL"; got != want {
			t.Errorf("RawQuery = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(`{"symbol":"AAPL","instrument_id":"1","result":[]}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	if _, err := c.GetTick(context.Background(), data.TickQuery{
		Symbol:   "AAPL",
		Category: data.StockCategoryUS,
	}); err != nil {
		t.Fatalf("GetTick() error = %v", err)
	}
}

// TestStockTicksInstrumentIDSpellings pins the two spellings Webull sends and the
// precedence between them. The sandbox futures response carries instrumentId while
// the page documents instrument_id, and encoding/json reaches neither from the
// other's tag, so a regression here is a caller silently reading "".
func TestStockTicksInstrumentIDSpellings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		want string
	}{
		{"documented snake_case", `{"symbol":"A","instrument_id":"1","result":[]}`, "1"},
		{"observed camelCase", `{"symbol":"A","instrumentId":"2","result":[]}`, "2"},
		{
			"both prefer documented",
			`{"symbol":"A","instrument_id":"1","instrumentId":"2","result":[]}`,
			"1",
		},
		{"neither", `{"symbol":"A","result":[]}`, ""},
		{"null camelCase", `{"symbol":"A","instrumentId":null,"result":[]}`, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var got data.StockTicks
			if err := json.Unmarshal([]byte(tc.body), &got); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			if got.InstrumentID != tc.want {
				t.Errorf("InstrumentID = %q, want %q", got.InstrumentID, tc.want)
			}
		})
	}
}

// TestStockTicksMarshalKeepsDocumentedName holds the wire form stable. A caller
// that round-trips a decoded value must still see instrument_id, or the fix for
// reading a response would have broken writing one.
func TestStockTicksMarshalKeepsDocumentedName(t *testing.T) {
	t.Parallel()

	var got data.StockTicks
	if err := json.Unmarshal([]byte(`{"instrumentId":"2","symbol":"A","result":[]}`), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	enc, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if !strings.Contains(string(enc), `"instrument_id":"2"`) {
		t.Errorf("encoded = %s, want the instrument_id member", enc)
	}
	if strings.Contains(string(enc), "instrumentId") {
		t.Errorf("encoded = %s, must not carry the camelCase member", enc)
	}
}

// TestGetFuturesTickAcceptsCamelCaseInstrumentID is the endpoint the defect was
// observed on, driven end to end rather than through json.Unmarshal alone.
func TestGetFuturesTickAcceptsCamelCaseInstrumentID(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/futures/ticks/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("category"), "US_FUTURES"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(`{"symbol":"ESZ5","instrumentId":"913256135","result":[]}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFuturesTick(context.Background(), data.FuturesTickQuery{Symbol: "ESZ5"})
	if err != nil {
		t.Fatalf("GetFuturesTick() error = %v", err)
	}
	if got.InstrumentID != "913256135" {
		t.Errorf("InstrumentID = %q, want %q", got.InstrumentID, "913256135")
	}
}
