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

func TestGetOptionTick(t *testing.T) {
	t.Parallel()

	const body = `{"symbol":"AAPL260522C00300000","instrument_id":"470059643",` +
		`"result":[{"time":"1761182953043","price":"48.07","volume":"1","side":"S"}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/market-data/options/ticks/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbol"), "AAPL260522C00300000"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_OPTION"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("count"), "30"; got != want {
			t.Errorf("count = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetOptionTick(context.Background(), data.OptionTickQuery{
		Symbol: "AAPL260522C00300000",
		Count:  30,
	})
	if err != nil {
		t.Fatalf("GetOptionTick() error = %v", err)
	}
	if got.Symbol != "AAPL260522C00300000" || got.InstrumentID != "470059643" {
		t.Errorf("identity = %+v", got)
	}
	if len(got.Result) != 1 || got.Result[0].Price != "48.07" || got.Result[0].Side != "S" {
		t.Errorf("ticks = %+v", got.Result)
	}
}

func TestGetOptionSnapshot(t *testing.T) {
	t.Parallel()

	const body = `[{"symbol":"AAPL260522C00300000","instrument_id":"470059643",` +
		`"price":"47.35","open":"48.17","high":"48.655","low":"46.855",` +
		`"pre_close":"47.704","volume":"48906","change":"-0.354",` +
		`"change_ratio":"-0.0074","last_trade_time":1761131406558,"close":"0.05",` +
		`"strike_price":"300.0","gamma":"1.0E-4","delta":"-0.0023","rho":"-0.001",` +
		`"theta":"-0.0037","vega":"0.0077","imp_vol":"0.609","open_interest":"14331",` +
		`"quote_time":1761131409276,"bid":"47.345","ask":"47.355","ask_size":"2",` +
		`"bid_size":"1","deal_amount":"70267.5"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/options/snapshots/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbols"), "AAPL260522C00300000"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_OPTION"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetOptionSnapshot(context.Background(), data.OptionSnapshotQuery{
		Symbols: []string{"AAPL260522C00300000"},
	})
	if err != nil {
		t.Fatalf("GetOptionSnapshot() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d snapshots, want 1", len(got))
	}
	snap := got[0]
	if snap.Symbol != "AAPL260522C00300000" || snap.Price != "47.35" || snap.Bid != "47.345" {
		t.Errorf("snapshot = %+v", snap)
	}
	if snap.StrikePrice != "300.0" || snap.ImpVol != "0.609" || snap.OpenInterest != "14331" {
		t.Errorf("greeks/interest = %+v", snap)
	}
	if snap.LastTradeTime != 1761131406558 || snap.QuoteTime != 1761131409276 {
		t.Errorf("times = %d/%d", snap.LastTradeTime, snap.QuoteTime)
	}
}

func TestGetOptionBars(t *testing.T) {
	t.Parallel()

	const body = `{"result":[{"symbol":"AAPL260522C00300000","instrument_id":"470059643",` +
		`"result":[{"time":"2021-12-28T09:00:09.945+0000","open":"1.3362","close":"1.3362",` +
		`"high":"1.3362","low":"1.3362","volume":"10"}]}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/options/bars/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbols"), "AAPL260522C00300000"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("timespan"), "D"; got != want {
			t.Errorf("timespan = %q, want %q", got, want)
		}
		if got, want := q.Get("count"), "5"; got != want {
			t.Errorf("count = %q, want %q", got, want)
		}
		if got, want := q.Get("real_time_required"), "true"; got != want {
			t.Errorf("real_time_required = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetOptionBars(context.Background(), data.OptionBarsQuery{
		Symbols:          []string{"AAPL260522C00300000"},
		Timespan:         data.OptionBarTimespanD,
		Count:            5,
		RealTimeRequired: true,
	})
	if err != nil {
		t.Fatalf("GetOptionBars() error = %v", err)
	}
	if len(got) != 1 || got[0].Symbol != "AAPL260522C00300000" {
		t.Fatalf("bars = %+v", got)
	}
	if len(got[0].Result) != 1 || got[0].Result[0].Close != "1.3362" {
		t.Errorf("bar = %+v", got[0].Result)
	}
}

func TestGetOptionContracts(t *testing.T) {
	t.Parallel()

	const body = `{"data":[{"instrument_id":"470059643",` +
		`"symbol":"AAPL260116C00300000","underlying_symbol":"AAPL",` +
		`"option_type":"CALL","strike_price":"300.0","expiration_date":"2026-01-16",` +
		`"exchange_code":"OPRA","category":"US_OPTION","currency":"USD","lot_size":"100"}],` +
		`"pagination_key":"page-2"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/trading/instruments/options/contracts/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbol"), "AAPL"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_OPTION"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("expiration"), "2026-01-16"; got != want {
			t.Errorf("expiration = %q, want %q", got, want)
		}
		if got, want := q.Get("option_type"), "CALL"; got != want {
			t.Errorf("option_type = %q, want %q", got, want)
		}
		if got, want := q.Get("strike_min"), "290"; got != want {
			t.Errorf("strike_min = %q, want %q", got, want)
		}
		if got, want := q.Get("strike_max"), "310"; got != want {
			t.Errorf("strike_max = %q, want %q", got, want)
		}
		if got, want := q.Get("pagination_key"), "page-1"; got != want {
			t.Errorf("pagination_key = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetOptionContracts(context.Background(), data.OptionContractsQuery{
		Symbol:        "AAPL",
		Expiration:    "2026-01-16",
		OptionType:    data.OptionTypeCall,
		StrikeMin:     "290",
		StrikeMax:     "310",
		PaginationKey: "page-1",
	})
	if err != nil {
		t.Fatalf("GetOptionContracts() error = %v", err)
	}
	if len(got.Contracts) != 1 {
		t.Fatalf("got %d contracts, want 1", len(got.Contracts))
	}
	con := got.Contracts[0]
	if con.InstrumentID != "470059643" || con.Symbol != "AAPL260116C00300000" {
		t.Errorf("identity = %+v", con)
	}
	if con.UnderlyingSymbol != "AAPL" || con.OptionType != data.OptionTypeCall {
		t.Errorf("underlying/type = %+v", con)
	}
	if con.StrikePrice != "300.0" || con.ExpirationDate != "2026-01-16" {
		t.Errorf("strike/expiration = %+v", con)
	}
	if con.ExchangeCode != "OPRA" || con.Category != data.OptionCategoryUS {
		t.Errorf("exchange/category = %+v", con)
	}
	if con.Currency != "USD" || con.LotSize != "100" {
		t.Errorf("currency/lot = %+v", con)
	}
	if got.PaginationKey != "page-2" {
		t.Errorf("paginationKey = %q, want %q", got.PaginationKey, "page-2")
	}
}

func TestGetOptionContractsOmitsOptionalParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got, want := q.Get("symbol"), "AAPL"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		for _, key := range []string{"expiration", "option_type", "strike_min", "strike_max", "pagination_key"} {
			if got := q.Get(key); got != "" {
				t.Errorf("%s = %q, want omitted", key, got)
			}
		}
		_, _ = w.Write([]byte(`{"data":[],"pagination_key":""}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetOptionContracts(context.Background(), data.OptionContractsQuery{Symbol: "AAPL"})
	if err != nil {
		t.Fatalf("GetOptionContracts() error = %v", err)
	}
	if len(got.Contracts) != 0 {
		t.Fatalf("got %d contracts, want 0", len(got.Contracts))
	}
}

// TestGetOptionContractsPartialOptionalParams confirms each optional filter is
// omitted independently and that a call-put selector is forwarded verbatim.
func TestGetOptionContractsPartialOptionalParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got, want := q.Get("symbol"), "AAPL"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_OPTION"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("option_type"), "PUT"; got != want {
			t.Errorf("option_type = %q, want %q", got, want)
		}
		if got, want := q.Get("strike_min"), "100"; got != want {
			t.Errorf("strike_min = %q, want %q", got, want)
		}
		for _, key := range []string{"expiration", "strike_max", "pagination_key"} {
			if got := q.Get(key); got != "" {
				t.Errorf("%s = %q, want omitted", key, got)
			}
		}
		_, _ = w.Write([]byte(`{"data":[],"pagination_key":""}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetOptionContracts(context.Background(), data.OptionContractsQuery{
		Symbol:     "AAPL",
		OptionType: data.OptionTypePut,
		StrikeMin:  "100",
	})
	if err != nil {
		t.Fatalf("GetOptionContracts() error = %v", err)
	}
	if got == nil || len(got.Contracts) != 0 {
		t.Fatalf("GetOptionContracts() = %+v, want empty result", got)
	}
}

// TestOptionEmptyAndMissingData drives every option endpoint with an empty or
// absent data envelope and asserts the call returns without panic and yields an
// empty result.
func TestOptionEmptyAndMissingData(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		path string
		body string
		call func(*data.Client) (int, error)
	}{
		{
			name: "tick empty object",
			path: "/market-data/options/ticks/list",
			body: `{}`,
			call: func(c *data.Client) (int, error) {
				res, err := c.GetOptionTick(context.Background(), data.OptionTickQuery{Symbol: "AAPL260522C00300000"})
				if res == nil {
					return 0, err
				}
				return len(res.Result), err
			},
		},
		{
			name: "snapshot null",
			path: "/market-data/options/snapshots/list",
			body: `null`,
			call: func(c *data.Client) (int, error) {
				res, err := c.GetOptionSnapshot(context.Background(), data.OptionSnapshotQuery{Symbols: []string{"AAPL260522C00300000"}})
				return len(res), err
			},
		},
		{
			name: "bars empty object",
			path: "/market-data/options/bars/list",
			body: `{}`,
			call: func(c *data.Client) (int, error) {
				res, err := c.GetOptionBars(context.Background(), data.OptionBarsQuery{
					Symbols:  []string{"AAPL260522C00300000"},
					Timespan: data.OptionBarTimespanD,
				})
				return len(res), err
			},
		},
		{
			name: "contracts missing data",
			path: "/trading/instruments/options/contracts/list",
			body: `{}`,
			call: func(c *data.Client) (int, error) {
				res, err := c.GetOptionContracts(context.Background(), data.OptionContractsQuery{Symbol: "AAPL"})
				if res == nil {
					return 0, err
				}
				return len(res.Contracts), err
			},
		},
		{
			name: "contracts null data",
			path: "/trading/instruments/options/contracts/list",
			body: `{"data":null,"pagination_key":""}`,
			call: func(c *data.Client) (int, error) {
				res, err := c.GetOptionContracts(context.Background(), data.OptionContractsQuery{Symbol: "AAPL"})
				if res == nil {
					return 0, err
				}
				return len(res.Contracts), err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.URL.Path; got != tc.path {
					t.Errorf("path = %q, want %q", got, tc.path)
				}
				_, _ = w.Write([]byte(tc.body))
			}))
			defer srv.Close()

			c := newTestClient(t, srv.URL)
			count, err := tc.call(c)
			if err != nil {
				t.Fatalf("call error = %v, want nil", err)
			}
			if count != 0 {
				t.Fatalf("result count = %d, want 0", count)
			}
		})
	}
}

func TestGetOptionBarsOmitsOptionalParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got, want := q.Get("category"), "US_OPTION"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got := q.Get("count"); got != "" {
			t.Errorf("count = %q, want omitted", got)
		}
		if got := q.Get("real_time_required"); got != "" {
			t.Errorf("real_time_required = %q, want omitted", got)
		}
		_, _ = w.Write([]byte(`{"result":[]}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetOptionBars(context.Background(), data.OptionBarsQuery{
		Symbols:  []string{"AAPL260522C00300000"},
		Timespan: data.OptionBarTimespanM1,
	})
	if err != nil {
		t.Fatalf("GetOptionBars() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("got %d symbol bars, want 0", len(got))
	}
}
