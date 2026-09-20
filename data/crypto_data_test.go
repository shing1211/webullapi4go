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

func TestGetCryptoBars_pathAndQuery(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/market-data/bars" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/market-data/bars")
		}
		q := r.URL.Query()
		if q.Get("category") != "CRYPTO" {
			t.Errorf("category = %q, want %q", q.Get("category"), "CRYPTO")
		}
		if q.Get("symbol") != "BTCUSD" {
			t.Errorf("symbol = %q, want %q", q.Get("symbol"), "BTCUSD")
		}
		if q.Get("timespan") != "D" {
			t.Errorf("timespan = %q, want %q", q.Get("timespan"), "D")
		}
		if q.Get("count") != "5" {
			t.Errorf("count = %q, want %q", q.Get("count"), "5")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetCryptoBars(context.Background(), data.CryptoBarsQuery{
		Symbol:   "BTCUSD",
		Interval: data.BarTimespanDay,
		Count:    5,
	})
	if err != nil {
		t.Fatalf("GetCryptoBars() error = %v", err)
	}
}

func TestGetCryptoBars_responseFields(t *testing.T) {
	t.Parallel()
	const body = `[{"symbol":"BTCUSD","exchange":"CRYPTO","currency":"USD","open":"95000","high":"97000","low":"94000","close":"96000","volume":"12345","turnover":"1180000000","timestamp":"1739100939000","timespan":"d1"}]`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	bars, err := c.GetCryptoBars(context.Background(), data.CryptoBarsQuery{Symbol: "BTCUSD"})
	if err != nil {
		t.Fatalf("GetCryptoBars() error = %v", err)
	}
	if len(bars) != 1 {
		t.Fatalf("len(bars) = %d, want 1", len(bars))
	}
	b := bars[0]
	if b.Symbol != "BTCUSD" {
		t.Errorf("Symbol = %q, want %q", b.Symbol, "BTCUSD")
	}
	if b.Open != "95000" || b.Close != "96000" {
		t.Errorf("open/close = %q/%q, want 95000/96000", b.Open, b.Close)
	}
	if b.Volume != "12345" {
		t.Errorf("Volume = %q, want %q", b.Volume, "12345")
	}
}

func TestGetCryptoBars_omitsOptionalParams(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		for _, key := range []string{"timespan", "count"} {
			if _, ok := q[key]; ok {
				t.Errorf("%s present, want omitted", key)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetCryptoBars(context.Background(), data.CryptoBarsQuery{Symbol: "ETHUSD"})
	if err != nil {
		t.Fatalf("GetCryptoBars() error = %v", err)
	}
}

func TestGetCryptoTick_pathAndQuery(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/market-data/tick" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/market-data/tick")
		}
		q := r.URL.Query()
		if q.Get("category") != "CRYPTO" {
			t.Errorf("category = %q, want %q", q.Get("category"), "CRYPTO")
		}
		if q.Get("symbol") != "ETHUSD" {
			t.Errorf("symbol = %q, want %q", q.Get("symbol"), "ETHUSD")
		}
		if q.Get("count") != "10" {
			t.Errorf("count = %q, want %q", q.Get("count"), "10")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetCryptoTick(context.Background(), data.CryptoTickQuery{Symbol: "ETHUSD", Count: 10})
	if err != nil {
		t.Fatalf("GetCryptoTick() error = %v", err)
	}
}

func TestGetCryptoTick_responseFields(t *testing.T) {
	t.Parallel()
	const body = `[{"symbol":"ETHUSD","exchange":"CRYPTO","currency":"USD","price":"3500","volume":"100","turnover":"350000","timestamp":"1739100939000","direction":"BUY","trade_id":"12345"}]`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	ticks, err := c.GetCryptoTick(context.Background(), data.CryptoTickQuery{Symbol: "ETHUSD"})
	if err != nil {
		t.Fatalf("GetCryptoTick() error = %v", err)
	}
	if len(ticks) != 1 {
		t.Fatalf("len(ticks) = %d, want 1", len(ticks))
	}
	if ticks[0].Price != "3500" || ticks[0].Direction != "BUY" {
		t.Errorf("tick = %+v", ticks[0])
	}
}

func TestGetCryptoDepth_pathAndQuery(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/market-data/depth" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/market-data/depth")
		}
		q := r.URL.Query()
		if q.Get("category") != "CRYPTO" {
			t.Errorf("category = %q, want %q", q.Get("category"), "CRYPTO")
		}
		if q.Get("symbol") != "BTCUSD" {
			t.Errorf("symbol = %q, want %q", q.Get("symbol"), "BTCUSD")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"symbol":"BTCUSD","exchange":"CRYPTO","currency":"USD","timestamp":"1739100939000","levels":[]}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetCryptoDepth(context.Background(), data.CryptoDepthQuery{Symbol: "BTCUSD"})
	if err != nil {
		t.Fatalf("GetCryptoDepth() error = %v", err)
	}
}

func TestGetCryptoDepth_responseFields(t *testing.T) {
	t.Parallel()
	const body = `{"symbol":"BTCUSD","exchange":"CRYPTO","currency":"USD","timestamp":"1739100939000","levels":[{"bid_price":"96000","bid_size":"1","ask_price":"96001","ask_size":"2"}]}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	depth, err := c.GetCryptoDepth(context.Background(), data.CryptoDepthQuery{Symbol: "BTCUSD"})
	if err != nil {
		t.Fatalf("GetCryptoDepth() error = %v", err)
	}
	if depth.Symbol != "BTCUSD" {
		t.Errorf("Symbol = %q, want %q", depth.Symbol, "BTCUSD")
	}
	if len(depth.Levels) != 1 {
		t.Fatalf("len(levels) = %d, want 1", len(depth.Levels))
	}
	if depth.Levels[0].BidPrice != "96000" {
		t.Errorf("bid_price = %q, want 96000", depth.Levels[0].BidPrice)
	}
}

func TestGetCryptoSnapshot_pathAndQuery(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/market-data/snapshot" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/market-data/snapshot")
		}
		q := r.URL.Query()
		if q.Get("category") != "CRYPTO" {
			t.Errorf("category = %q, want %q", q.Get("category"), "CRYPTO")
		}
		if q.Get("symbol") != "BTCUSD" {
			t.Errorf("symbol = %q, want %q", q.Get("symbol"), "BTCUSD")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetCryptoSnapshot(context.Background(), data.CryptoSnapshotQuery{Symbol: "BTCUSD"})
	if err != nil {
		t.Fatalf("GetCryptoSnapshot() error = %v", err)
	}
}

func TestGetCryptoSnapshot_responseFields(t *testing.T) {
	t.Parallel()
	const body = `{"symbol":"BTCUSD","exchange":"CRYPTO","currency":"USD","last_price":"96000","open":"95000","high":"97000","low":"94000","close":"96000","volume":"12345","turnover":"1180000000","bid":"95999","ask":"96001","timestamp":"1739100939000"}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var snap data.CryptoSnapshot
		if err := json.Unmarshal([]byte(body), &snap); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}
		if snap.LastPrice != "96000" {
			t.Errorf("LastPrice = %q, want 96000", snap.LastPrice)
		}
		if snap.Bid != "95999" || snap.Ask != "96001" {
			t.Errorf("bid/ask = %q/%q, want 95999/96001", snap.Bid, snap.Ask)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	snap, err := c.GetCryptoSnapshot(context.Background(), data.CryptoSnapshotQuery{Symbol: "BTCUSD"})
	if err != nil {
		t.Fatalf("GetCryptoSnapshot() error = %v", err)
	}
	if snap.LastPrice != "96000" || snap.Bid != "95999" {
		t.Errorf("snapshot = %+v", snap)
	}
}
