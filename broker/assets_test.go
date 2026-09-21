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

package broker_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shing1211/webullapi4go/broker"
)

func TestGetBalance(t *testing.T) {
	t.Parallel()

	const body = `{"total_asset_currency":"USD","total_cash_balance":"200000.00",` +
		`"total_market_value":"300000.00","total_unrealized_profit_loss":"15000.00",` +
		`"init_margin":"50000.00","account_currency_assets":[{"currency":"USD","balance":"200000.00"}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/broker/assets/balances/get"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("account_id"), "ACC1"; got != want {
			t.Errorf("account_id = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := broker.New(newTestClient(t, srv.URL))
	got, err := c.GetBalance(context.Background(), "ACC1")
	if err != nil {
		t.Fatalf("GetBalance() error = %v", err)
	}
	if got.TotalCashBalance != "200000.00" {
		t.Errorf("total_cash_balance = %q, want 200000.00", got.TotalCashBalance)
	}
	if got.TotalMarketValue != "300000.00" {
		t.Errorf("total_market_value = %q, want 300000.00", got.TotalMarketValue)
	}
	if got.TotalUnrealizedPL != "15000.00" {
		t.Errorf("total_unrealized_profit_loss = %q, want 15000.00", got.TotalUnrealizedPL)
	}
	if got.InitMargin != "50000.00" {
		t.Errorf("init_margin = %q, want 50000.00", got.InitMargin)
	}
	if len(got.AccountCurrencyAssets) != 1 {
		t.Errorf("account_currency_assets len = %d, want 1", len(got.AccountCurrencyAssets))
	}
	if got.AccountCurrencyAssets[0].Currency != "USD" {
		t.Errorf("account_currency_assets[0].currency = %q, want USD", got.AccountCurrencyAssets[0].Currency)
	}
}

func TestGetPositions(t *testing.T) {
	t.Parallel()

	const body = `{"data":[{"symbol":"00700.HK","quantity":"100",` +
		`"average_cost":"350.00","market_value":"36000.00",` +
		`"unrealized_pl":"1000.00","instrument_type":"EQUITY",` +
		`"currency":"HKD"}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/broker/assets/positions/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("account_id"), "ACC1"; got != want {
			t.Errorf("account_id = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := broker.New(newTestClient(t, srv.URL))
	got, err := c.GetPositions(context.Background(), "ACC1")
	if err != nil {
		t.Fatalf("GetPositions() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d positions, want 1", len(got))
	}
	pos := got[0]
	if pos.Symbol != "00700.HK" {
		t.Errorf("symbol = %q, want 00700.HK", pos.Symbol)
	}
	if pos.Quantity != "100" {
		t.Errorf("quantity = %q, want 100", pos.Quantity)
	}
	if pos.AverageCost != "350.00" || pos.MarketValue != "36000.00" {
		t.Errorf("cost/market = %+v", pos)
	}
	if pos.UnrealizedPL != "1000.00" {
		t.Errorf("unrealized_pl = %q, want 1000.00", pos.UnrealizedPL)
	}
	if pos.InstrumentType != "EQUITY" {
		t.Errorf("instrument_type = %q, want EQUITY", pos.InstrumentType)
	}
	if pos.Currency != "HKD" {
		t.Errorf("currency = %q, want HKD", pos.Currency)
	}
}

func TestGetPositionsEmpty(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer srv.Close()

	c := broker.New(newTestClient(t, srv.URL))
	got, err := c.GetPositions(context.Background(), "ACC1")
	if err != nil {
		t.Fatalf("GetPositions() error = %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %d positions, want 0", len(got))
	}
}
