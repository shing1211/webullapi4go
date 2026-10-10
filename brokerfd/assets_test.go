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
	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

func TestGetFDAssetsSummary(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(FDAssetsSummary{
			AccountID:    "A1",
			TotalEquity:  money.Must(money.NewFromString("100000")),
			CashBalance:  money.Must(money.NewFromString("50000")),
			MarketValue:  money.Must(money.NewFromString("50000")),
			BuyingPower:  money.Must(money.NewFromString("100000")),
			UnrealizedPL: money.Must(money.NewFromString("1000")),
			RealizedPL:   money.Must(money.NewFromString("500")),
			Currency:     "USD",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithAppKey(testAppKey), client.WithAppSecret(testAppSecret), client.WithEndpoints(client.Endpoints{HTTP: srv.URL, BrokerHTTP: srv.URL}))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDAssetsSummary(context.Background(), "A1")
	if err != nil {
		t.Fatalf("GetFDAssetsSummary error = %v", err)
	}
	if got.AccountID != "A1" {
		t.Fatalf("AccountID = %s, want A1", got.AccountID)
	}
	if got.TotalEquity.Cmp(money.Must(money.NewFromString("100000"))) != 0 {
		t.Fatalf("TotalEquity = %v, want 100000", got.TotalEquity)
	}
	if capturedReq.URL.Path != pathFDAssetsSummary {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDAssetsSummary)
	}
	if capturedReq.URL.Query().Get("account_id") != "A1" {
		t.Fatalf("account_id = %s, want A1", capturedReq.URL.Query().Get("account_id"))
	}
}

func TestGetFDAssetsDetail(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		// The documented 200 body is an object wrapping the currency array together
		// with the account-level totals. Serving the SDK's own shape here instead
		// would make this test pass whichever shape the method decoded, which is how
		// the totals went missing without any test failing.
		_, _ = w.Write([]byte(`{
			"account_currency_assets": [
				{"currency":"USD","cash_balance":"50000"},
				{"currency":"HKD","cash_balance":"100000"}
			],
			"total_asset_currency": "USD",
			"total_cash_balance": "485705.0"
		}`))
	}))
	defer srv.Close()

	cl, err := client.New(client.WithAppKey(testAppKey), client.WithAppSecret(testAppSecret), client.WithEndpoints(client.Endpoints{HTTP: srv.URL, BrokerHTTP: srv.URL}))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDAssetsDetail(context.Background(), "A1")
	if err != nil {
		t.Fatalf("GetFDAssetsDetail error = %v", err)
	}
	if len(got.AccountCurrencyAssets) != 2 {
		t.Fatalf("len(AccountCurrencyAssets) = %d, want 2", len(got.AccountCurrencyAssets))
	}
	if got.AccountCurrencyAssets[0].Currency != "USD" {
		t.Fatalf("Currency[0] = %s, want USD", got.AccountCurrencyAssets[0].Currency)
	}
	// The two totals are the reason the envelope exists. Before v2.1.33 they reached
	// no field and were dropped with no error reported.
	if got.TotalAssetCurrency != "USD" {
		t.Errorf("TotalAssetCurrency = %q, want USD", got.TotalAssetCurrency)
	}
	if want := money.Must(money.NewFromString("485705.0")); got.TotalCashBalance.Cmp(want) != 0 {
		t.Errorf("TotalCashBalance = %s, want %s", got.TotalCashBalance, want)
	}
	if capturedReq.URL.Path != pathFDAssetsDetail {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDAssetsDetail)
	}
	if capturedReq.URL.Query().Get("account_id") != "A1" {
		t.Fatalf("account_id = %s, want A1", capturedReq.URL.Query().Get("account_id"))
	}
}

func TestGetFDPositions(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]FDPosition{
			{PositionID: "P1", AccountID: "A1", Symbol: "AAPL", Quantity: "100", AverageCost: money.Must(money.NewFromString("150.00")), MarketValue: money.Must(money.NewFromString("17500")), UnrealizedPL: money.Must(money.NewFromString("2500")), RealizedPL: money.Must(money.NewFromString("0")), InstrumentType: "STOCK", Currency: "USD"},
			{PositionID: "P2", AccountID: "A1", Symbol: "TSLA", Quantity: "50", AverageCost: money.Must(money.NewFromString("200.00")), MarketValue: money.Must(money.NewFromString("10000")), UnrealizedPL: money.Must(money.NewFromString("0")), RealizedPL: money.Must(money.NewFromString("500")), InstrumentType: "STOCK", Currency: "USD"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithAppKey(testAppKey), client.WithAppSecret(testAppSecret), client.WithEndpoints(client.Endpoints{HTTP: srv.URL, BrokerHTTP: srv.URL}))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDPositions(context.Background(), "A1")
	if err != nil {
		t.Fatalf("GetFDPositions error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if got[0].Symbol != "AAPL" {
		t.Fatalf("Symbol[0] = %s, want AAPL", got[0].Symbol)
	}
	if capturedReq.URL.Path != pathFDAssetsPositions {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDAssetsPositions)
	}
	if capturedReq.URL.Query().Get("account_id") != "A1" {
		t.Fatalf("account_id = %s, want A1", capturedReq.URL.Query().Get("account_id"))
	}
}
