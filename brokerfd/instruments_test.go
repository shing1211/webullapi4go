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

func TestGetFDStockInstruments(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]FDStockInstrument{
			{Symbol: "AAPL", Name: "Apple Inc.", Exchange: "NASDAQ", Currency: "USD", LotSize: "100", Status: "active"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDStockInstruments(context.Background(), []string{"AAPL"})
	if err != nil {
		t.Fatalf("GetFDStockInstruments error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
	if capturedReq.URL.Path != pathFDStockInstruments {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDStockInstruments)
	}
	if got[0].Symbol != "AAPL" {
		t.Fatalf("Symbol = %s, want AAPL", got[0].Symbol)
	}
}

func TestGetFDStockLocate(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]FDStockLocate{
			{Symbol: "AAPL", LocateQuantity: "1000", Available: "800", Rate: money.Must(money.NewFromString("0.05"))},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDStockLocate(context.Background(), "AAPL")
	if err != nil {
		t.Fatalf("GetFDStockLocate error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
	if capturedReq.URL.Path != pathFDStockLocate {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDStockLocate)
	}
	if capturedReq.URL.Query().Get("symbol") != "AAPL" {
		t.Fatalf("symbol = %s, want AAPL", capturedReq.URL.Query().Get("symbol"))
	}
}

func TestGetFDCorporateActions(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]FDCorporateAction{
			{ActionID: "CA1", Symbol: "AAPL", ActionType: "DIVIDEND", ExDate: "2026-01-10", RecordDate: "2026-01-11", PayDate: "2026-01-15", Ratio: "0.25"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDCorporateActions(context.Background(), "AAPL")
	if err != nil {
		t.Fatalf("GetFDCorporateActions error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
	if capturedReq.URL.Path != pathFDCorporateActions {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDCorporateActions)
	}
	if got[0].ActionType != "DIVIDEND" {
		t.Fatalf("ActionType = %s, want DIVIDEND", got[0].ActionType)
	}
}

func TestGetFDCorporateActionDetail(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(FDCorporateAction{ActionID: "CA1", Symbol: "AAPL", ActionType: "SPLIT", ExDate: "2026-02-01", RecordDate: "2026-02-02", PayDate: "2026-02-03", Ratio: "4:1"})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDCorporateActionDetail(context.Background(), "CA1")
	if err != nil {
		t.Fatalf("GetFDCorporateActionDetail error = %v", err)
	}
	if capturedReq.URL.Path != pathFDCorporateActionDetail {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDCorporateActionDetail)
	}
	if capturedReq.URL.Query().Get("action_id") != "CA1" {
		t.Fatalf("action_id = %s, want CA1", capturedReq.URL.Query().Get("action_id"))
	}
	if got.ActionType != "SPLIT" {
		t.Fatalf("ActionType = %s, want SPLIT", got.ActionType)
	}
}

func TestGetFDECInstruments(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]FDECInstrument{
			{Symbol: "EC2026AAPL", EventID: "EVT1", SeriesID: "SER1", StrikePrice: money.Must(money.NewFromString("150.00")), ExpirationDate: "2026-06-20", Status: "active"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDECInstruments(context.Background(), "EVT1")
	if err != nil {
		t.Fatalf("GetFDECInstruments error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
	if capturedReq.URL.Path != pathFDECInstruments {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDECInstruments)
	}
	if got[0].EventID != "EVT1" {
		t.Fatalf("EventID = %s, want EVT1", got[0].EventID)
	}
}

func TestGetFDECInstrumentDetail(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(FDECInstrument{Symbol: "EC2026AAPL", EventID: "EVT1", SeriesID: "SER1", StrikePrice: money.Must(money.NewFromString("150.00")), ExpirationDate: "2026-06-20", Status: "active"})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDECInstrumentDetail(context.Background(), "EC2026AAPL")
	if err != nil {
		t.Fatalf("GetFDECInstrumentDetail error = %v", err)
	}
	if capturedReq.URL.Path != pathFDECInstrumentDetail {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDECInstrumentDetail)
	}
	if capturedReq.URL.Query().Get("symbol") != "EC2026AAPL" {
		t.Fatalf("symbol = %s, want EC2026AAPL", capturedReq.URL.Query().Get("symbol"))
	}
	if got.StrikePrice.Cmp(money.Must(money.NewFromString("150.00"))) != 0 {
		t.Fatalf("StrikePrice = %v, want 150.00", got.StrikePrice)
	}
}
