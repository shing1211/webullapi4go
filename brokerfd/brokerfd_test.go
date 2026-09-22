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

func TestGetAccountsSummary(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]AccountSummary{
			{AccountID: "A1", Currency: "USD", NetLiquidity: "100000.00", CashBalance: "50000.00", MarketValue: "45000.00", BuyingPower: "200000.00"},
			{AccountID: "A2", Currency: "USD", NetLiquidity: "200000.00", CashBalance: "80000.00", MarketValue: "115000.00", BuyingPower: "400000.00"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetAccountsSummary(context.Background())
	if err != nil {
		t.Fatalf("GetAccountsSummary error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if got[0].AccountID != "A1" {
		t.Fatalf("AccountID = %s, want A1", got[0].AccountID)
	}
	if got[1].AccountID != "A2" {
		t.Fatalf("AccountID = %s, want A2", got[1].AccountID)
	}
	if capturedReq.URL.Path != "/broker-fd/accounts" {
		t.Fatalf("path = %s, want /broker-fd/accounts", capturedReq.URL.Path)
	}
	if capturedReq.Method != http.MethodGet {
		t.Fatalf("method = %s, want %s", capturedReq.Method, http.MethodGet)
	}
}

func TestGetPositions(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]PositionSummary{
			{Symbol: "AAPL", Quantity: "100", MarketValue: "17500.00", CostBasis: "15000.00", UnrealizedPL: "2500.00"},
			{Symbol: "GOOGL", Quantity: "50", MarketValue: "9000.00", CostBasis: "8000.00", UnrealizedPL: "1000.00"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetPositions(context.Background())
	if err != nil {
		t.Fatalf("GetPositions error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if got[0].Symbol != "AAPL" {
		t.Fatalf("Symbol = %s, want AAPL", got[0].Symbol)
	}
	if got[1].Symbol != "GOOGL" {
		t.Fatalf("Symbol = %s, want GOOGL", got[1].Symbol)
	}
	if capturedReq.URL.Path != "/broker-fd/positions" {
		t.Fatalf("path = %s, want /broker-fd/positions", capturedReq.URL.Path)
	}
	if capturedReq.Method != http.MethodGet {
		t.Fatalf("method = %s, want %s", capturedReq.Method, http.MethodGet)
	}
}
