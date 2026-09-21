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

func TestGetFDAssetsSummary(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(FDAssetsSummary{
			AccountID:    "A1",
			TotalEquity:  "100000",
			CashBalance:  "50000",
			MarketValue:  "50000",
			BuyingPower:  "100000",
			UnrealizedPL: "1000",
			RealizedPL:   "500",
			Currency:     "USD",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithAppKey(testAppKey), client.WithAppSecret(testAppSecret), client.WithBaseURL(srv.URL))
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
	if got.TotalEquity != "100000" {
		t.Fatalf("TotalEquity = %s, want 100000", got.TotalEquity)
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
		_ = json.NewEncoder(w).Encode([]FDAssetDetail{
			{Currency: "USD", CashBalance: "50000", MarketValue: "50000", BuyingPower: "100000", AvailableCash: "40000"},
			{Currency: "HKD", CashBalance: "100000", MarketValue: "0", BuyingPower: "200000", AvailableCash: "100000"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithAppKey(testAppKey), client.WithAppSecret(testAppSecret), client.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDAssetsDetail(context.Background(), "A1")
	if err != nil {
		t.Fatalf("GetFDAssetsDetail error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if got[0].Currency != "USD" {
		t.Fatalf("Currency[0] = %s, want USD", got[0].Currency)
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
			{PositionID: "P1", AccountID: "A1", Symbol: "AAPL", Quantity: "100", AverageCost: "150.00", MarketValue: "17500", UnrealizedPL: "2500", RealizedPL: "0", InstrumentType: "STOCK", Currency: "USD"},
			{PositionID: "P2", AccountID: "A1", Symbol: "TSLA", Quantity: "50", AverageCost: "200.00", MarketValue: "10000", UnrealizedPL: "0", RealizedPL: "500", InstrumentType: "STOCK", Currency: "USD"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithAppKey(testAppKey), client.WithAppSecret(testAppSecret), client.WithBaseURL(srv.URL))
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
