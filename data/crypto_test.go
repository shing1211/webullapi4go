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

func TestGetCryptoBars(t *testing.T) {
	t.Parallel()
	const body = `[{"symbol":"BTCUSD","instrument_id":"1","result":[` +
		`{"time":"2026-01-01T00:00:00Z","open":"1","close":"2","high":"3","low":"0.5","volume":"10"}]}]`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/crypto/bars/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if q.Get("symbols") != "BTCUSD" || q.Get("category") != "US_CRYPTO" {
			t.Errorf("query = %v", q)
		}
		if q.Get("timespan") != "D" || q.Get("real_time_required") != "true" {
			t.Errorf("query = %v", q)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetCryptoBars(context.Background(), data.CryptoBarsQuery{
		Symbols:          []string{"BTCUSD"},
		Category:         "US_CRYPTO",
		Timespan:         data.BarTimespanDay,
		Count:            5,
		RealTimeRequired: true,
	})
	if err != nil {
		t.Fatalf("GetCryptoBars() error = %v", err)
	}
	if len(got) != 1 || got[0].Symbol != "BTCUSD" || len(got[0].Result) != 1 {
		t.Fatalf("bars = %+v", got)
	}
	if got[0].Result[0].Close != "2" {
		t.Errorf("close = %q", got[0].Result[0].Close)
	}
}

func TestGetCryptoSnapshot(t *testing.T) {
	t.Parallel()
	const body = `[{"instrument_id":"1","symbol":"BTCUSD","pre_close":"100","last_trade_time":1,` +
		`"price":"101","open":"100","high":"102","low":"99","change":"1","change_ratio":"0.01",` +
		`"quote_time":"2026-01-01","bid":"100.5","bid_size":"1","ask":"101.5","ask_size":"2"}]`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/crypto/snapshots/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if r.URL.Query().Get("symbols") != "BTCUSD" {
			t.Errorf("symbols = %q", r.URL.Query().Get("symbols"))
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetCryptoSnapshot(context.Background(), data.CryptoSnapshotQuery{
		Symbols:  []string{"BTCUSD"},
		Category: "US_CRYPTO",
	})
	if err != nil {
		t.Fatalf("GetCryptoSnapshot() error = %v", err)
	}
	if len(got) != 1 || got[0].Price != "101" || got[0].Bid != "100.5" {
		t.Fatalf("snapshot = %+v", got)
	}
}

func TestGetCryptoInstruments(t *testing.T) {
	t.Parallel()
	const body = `{"data":[{"symbol":"BTCUSD","name":"Bitcoin","category":"US_CRYPTO",` +
		`"currency":"USD","status":"OC","instrument_id":"1"}],"pagination_key":"p2"}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/trading/instruments/crypto/profiles/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if r.URL.Query().Get("category") != "US_CRYPTO" {
			t.Errorf("category = %q", r.URL.Query().Get("category"))
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetCryptoInstruments(context.Background(), data.CryptoInstrumentQuery{Category: "US_CRYPTO"})
	if err != nil {
		t.Fatalf("GetCryptoInstruments() error = %v", err)
	}
	if len(got.Instruments) != 1 || got.Instruments[0].Symbol != "BTCUSD" || got.PaginationKey != "p2" {
		t.Fatalf("instruments = %+v", got)
	}
}
