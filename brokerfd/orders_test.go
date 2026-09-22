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

func TestPreviewFDOrder(t *testing.T) {
	var gotReq FDOrderPreviewRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != pathFDOrderPreview {
			t.Errorf("expected %s, got %s", pathFDOrderPreview, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotReq); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(FDOrderPreview{
			OrderID:        "prev123",
			EstimatedFee:   "0.99",
			EstimatedTotal: "100.99",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	c := New(cl)

	req := FDOrderPreviewRequest{
		AccountID:   "acc1",
		Symbol:      "AAPL",
		OrderType:   "LIMIT",
		Side:        "BUY",
		Quantity:    "10",
		LimitPrice:  "150.00",
		TimeInForce: "DAY",
	}
	got, err := c.PreviewFDOrder(context.Background(), req)
	if err != nil {
		t.Fatalf("PreviewFDOrder failed: %v", err)
	}
	if got.OrderID != "prev123" {
		t.Errorf("expected order_id prev123, got %s", got.OrderID)
	}
	if gotReq.AccountID != "acc1" {
		t.Errorf("expected account_id acc1, got %s", gotReq.AccountID)
	}
}

func TestPlaceFDOrder(t *testing.T) {
	var gotReq FDOrderPreviewRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != pathFDOrderPlace {
			t.Errorf("expected %s, got %s", pathFDOrderPlace, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotReq); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(FDOrder{
			OrderID:        "ord123",
			AccountID:      "acc1",
			Symbol:         "AAPL",
			OrderType:      "LIMIT",
			Side:           "BUY",
			Quantity:       "10",
			LimitPrice:     "150.00",
			TimeInForce:    "DAY",
			Status:         "SUBMITTED",
			FilledQuantity: "0",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	c := New(cl)

	req := FDOrderPreviewRequest{
		AccountID:   "acc1",
		Symbol:      "AAPL",
		OrderType:   "LIMIT",
		Side:        "BUY",
		Quantity:    "10",
		LimitPrice:  "150.00",
		TimeInForce: "DAY",
	}
	got, err := c.PlaceFDOrder(context.Background(), req)
	if err != nil {
		t.Fatalf("PlaceFDOrder failed: %v", err)
	}
	if got.OrderID != "ord123" {
		t.Errorf("expected order_id ord123, got %s", got.OrderID)
	}
	if gotReq.Symbol != "AAPL" {
		t.Errorf("expected symbol AAPL, got %s", gotReq.Symbol)
	}
}

func TestReplaceFDOrder(t *testing.T) {
	var gotReq ReplaceFDOrderRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != pathFDOrderReplace {
			t.Errorf("expected %s, got %s", pathFDOrderReplace, r.URL.Path)
		}
		if r.URL.Query().Get("order_id") != "ord123" {
			t.Errorf("expected order_id ord123 in query, got %s", r.URL.Query().Get("order_id"))
		}
		if err := json.NewDecoder(r.Body).Decode(&gotReq); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(FDOrder{
			OrderID:     "ord123",
			AccountID:   "acc1",
			Symbol:      "AAPL",
			OrderType:   "LIMIT",
			Side:        "BUY",
			Quantity:    "20",
			LimitPrice:  "155.00",
			TimeInForce: "DAY",
			Status:      "SUBMITTED",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	c := New(cl)

	req := ReplaceFDOrderRequest{
		OrderID:    "ord123",
		LimitPrice: "155.00",
		Quantity:   "20",
	}
	got, err := c.ReplaceFDOrder(context.Background(), "ord123", req)
	if err != nil {
		t.Fatalf("ReplaceFDOrder failed: %v", err)
	}
	if got.OrderID != "ord123" {
		t.Errorf("expected order_id ord123, got %s", got.OrderID)
	}
	if gotReq.LimitPrice != "155.00" {
		t.Errorf("expected limit_price 155.00, got %s", gotReq.LimitPrice)
	}
}

func TestCancelFDOrder(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != pathFDOrderCancel {
			t.Errorf("expected %s, got %s", pathFDOrderCancel, r.URL.Path)
		}
		if r.URL.Query().Get("order_id") != "ord123" {
			t.Errorf("expected order_id ord123 in query, got %s", r.URL.Query().Get("order_id"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	c := New(cl)

	if err := c.CancelFDOrder(context.Background(), "ord123"); err != nil {
		t.Fatalf("CancelFDOrder failed: %v", err)
	}
}

func TestGetFDOrderDetail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != pathFDOrderDetail {
			t.Errorf("expected %s, got %s", pathFDOrderDetail, r.URL.Path)
		}
		if r.URL.Query().Get("order_id") != "ord123" {
			t.Errorf("expected order_id ord123 in query, got %s", r.URL.Query().Get("order_id"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(FDOrder{
			OrderID:        "ord123",
			AccountID:      "acc1",
			Symbol:         "AAPL",
			OrderType:      "LIMIT",
			Side:           "BUY",
			Quantity:       "10",
			LimitPrice:     "150.00",
			TimeInForce:    "DAY",
			Status:         "FILLED",
			FilledQuantity: "10",
			AvgFillPrice:   "150.25",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	c := New(cl)

	got, err := c.GetFDOrderDetail(context.Background(), "ord123")
	if err != nil {
		t.Fatalf("GetFDOrderDetail failed: %v", err)
	}
	if got.OrderID != "ord123" {
		t.Errorf("expected order_id ord123, got %s", got.OrderID)
	}
	if got.Status != "FILLED" {
		t.Errorf("expected status FILLED, got %s", got.Status)
	}
}

func TestGetFDOrderHistory(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != pathFDOrderHistory {
			t.Errorf("expected %s, got %s", pathFDOrderHistory, r.URL.Path)
		}
		if r.URL.Query().Get("account_id") != "acc1" {
			t.Errorf("expected account_id acc1 in query, got %s", r.URL.Query().Get("account_id"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]FDOrder{
			{OrderID: "ord1", AccountID: "acc1", Symbol: "AAPL", Status: "FILLED"},
			{OrderID: "ord2", AccountID: "acc1", Symbol: "TSLA", Status: "CANCELLED"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	c := New(cl)

	got, err := c.GetFDOrderHistory(context.Background(), "acc1")
	if err != nil {
		t.Fatalf("GetFDOrderHistory failed: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 orders, got %d", len(got))
	}
	if got[0].Symbol != "AAPL" {
		t.Errorf("expected first symbol AAPL, got %s", got[0].Symbol)
	}
}

func TestGetFDOpenOrders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != pathFDOrderOpen {
			t.Errorf("expected %s, got %s", pathFDOrderOpen, r.URL.Path)
		}
		if r.URL.Query().Get("account_id") != "acc1" {
			t.Errorf("expected account_id acc1 in query, got %s", r.URL.Query().Get("account_id"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]FDOrder{
			{OrderID: "ord1", AccountID: "acc1", Symbol: "AAPL", Status: "SUBMITTED"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	c := New(cl)

	got, err := c.GetFDOpenOrders(context.Background(), "acc1")
	if err != nil {
		t.Fatalf("GetFDOpenOrders failed: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 order, got %d", len(got))
	}
	if got[0].Status != "SUBMITTED" {
		t.Errorf("expected status SUBMITTED, got %s", got[0].Status)
	}
}
