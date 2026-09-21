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
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shing1211/webullapi4go/broker"
	"github.com/shing1211/webullapi4go/client"
)

const (
	testAppKey    = "test-app-key"
	testAppSecret = "test-app-secret"
)

type requestInfo struct {
	Method string
	Path   string
	Body   []byte
}

func newTestCoreClient(t *testing.T, baseURL string) *client.Client {
	t.Helper()
	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithEndpoints(client.Endpoints{
			HTTP:       baseURL,
			BrokerHTTP: baseURL,
		}),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() {
		if err := cl.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	})
	return cl
}

func TestPreviewOrder(t *testing.T) {
	t.Parallel()
	var captured requestInfo
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		captured = requestInfo{Method: r.Method, Path: r.URL.RequestURI(), Body: body}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"order_id":"O-1","estimated_cost":"501.23","estimated_transaction_fee":"1.23"}`))
	}))
	defer srv.Close()

	cl := broker.New(newTestCoreClient(t, srv.URL))
	req := broker.PreviewOrderRequest{
		AccountID:   "ACC-1",
		Symbol:      "AAPL",
		OrderType:   "LIMIT",
		Side:        "BUY",
		Quantity:    "10",
		LimitPrice:  "150.00",
		TimeInForce: "DAY",
	}
	out, err := cl.PreviewOrder(context.Background(), req)
	if err != nil {
		t.Fatalf("PreviewOrder: %v", err)
	}
	if captured.Method != http.MethodPost {
		t.Errorf("method = %s, want POST", captured.Method)
	}
	if captured.Path != "/broker/orders/preview" {
		t.Errorf("path = %s, want /broker/orders/preview", captured.Path)
	}
	if out.OrderID != "O-1" {
		t.Errorf("order_id = %s, want O-1", out.OrderID)
	}
	if out.EstimatedFee != "1.23" {
		t.Errorf("estimated_transaction_fee = %s, want 1.23", out.EstimatedFee)
	}
}

func TestPlaceOrder(t *testing.T) {
	t.Parallel()
	var captured requestInfo
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		captured = requestInfo{Method: r.Method, Path: r.URL.RequestURI(), Body: body}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"order_id":"O-2","account_id":"ACC-1","symbol":"AAPL","order_type":"LIMIT","side":"BUY","quantity":"10","time_in_force":"DAY","status":"SUBMITTED"}`))
	}))
	defer srv.Close()

	cl := broker.New(newTestCoreClient(t, srv.URL))
	req := broker.PreviewOrderRequest{
		AccountID:   "ACC-1",
		Symbol:      "AAPL",
		OrderType:   "LIMIT",
		Side:        "BUY",
		Quantity:    "10",
		LimitPrice:  "150.00",
		TimeInForce: "DAY",
	}
	out, err := cl.PlaceOrder(context.Background(), req)
	if err != nil {
		t.Fatalf("PlaceOrder: %v", err)
	}
	if captured.Method != http.MethodPost {
		t.Errorf("method = %s, want POST", captured.Method)
	}
	if captured.Path != "/broker/orders/place" {
		t.Errorf("path = %s, want /broker/orders/place", captured.Path)
	}
	if out.OrderID != "O-2" {
		t.Errorf("order_id = %s, want O-2", out.OrderID)
	}
	if out.Status != "SUBMITTED" {
		t.Errorf("status = %s, want SUBMITTED", out.Status)
	}
}

func TestReplaceOrder(t *testing.T) {
	t.Parallel()
	var captured requestInfo
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		captured = requestInfo{Method: r.Method, Path: r.URL.RequestURI(), Body: body}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"order_id":"O-3","account_id":"ACC-1","symbol":"AAPL","order_type":"LIMIT","side":"BUY","quantity":"20","time_in_force":"DAY","status":"SUBMITTED"}`))
	}))
	defer srv.Close()

	cl := broker.New(newTestCoreClient(t, srv.URL))
	req := broker.ReplaceOrderRequest{
		OrderID:    "O-3",
		LimitPrice: "155.00",
		Quantity:   "20",
	}
	out, err := cl.ReplaceOrder(context.Background(), "O-3", req)
	if err != nil {
		t.Fatalf("ReplaceOrder: %v", err)
	}
	if captured.Method != http.MethodPost {
		t.Errorf("method = %s, want POST", captured.Method)
	}
	if captured.Path != "/broker/orders/replace" {
		t.Errorf("path = %s, want /broker/orders/replace", captured.Path)
	}
	if out.OrderID != "O-3" {
		t.Errorf("order_id = %s, want O-3", out.OrderID)
	}

	var got broker.ReplaceOrderRequest
	if err := json.Unmarshal(captured.Body, &got); err != nil {
		t.Fatalf("unmarshalling request body: %v", err)
	}
	if got.LimitPrice != "155.00" {
		t.Errorf("request limit_price = %s, want 155.00", got.LimitPrice)
	}
	if got.Quantity != "20" {
		t.Errorf("request quantity = %s, want 20", got.Quantity)
	}
}

func TestCancelOrder(t *testing.T) {
	t.Parallel()
	var captured requestInfo
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		captured = requestInfo{Method: r.Method, Path: r.URL.RequestURI(), Body: body}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	cl := broker.New(newTestCoreClient(t, srv.URL))
	err := cl.CancelOrder(context.Background(), "O-4")
	if err != nil {
		t.Fatalf("CancelOrder: %v", err)
	}
	if captured.Method != http.MethodPost {
		t.Errorf("method = %s, want POST", captured.Method)
	}
	if captured.Path != "/broker/orders/cancel" {
		t.Errorf("path = %s, want /broker/orders/cancel", captured.Path)
	}
}

func TestGetOrderDetail(t *testing.T) {
	t.Parallel()
	var captured requestInfo
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		captured = requestInfo{Method: r.Method, Path: r.URL.RequestURI(), Body: body}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"order_id":"O-5","account_id":"ACC-1","symbol":"TSLA","order_type":"MARKET","side":"BUY","quantity":"5","time_in_force":"DAY","status":"FILLED","avg_fill_price":"200.00"}`))
	}))
	defer srv.Close()

	cl := broker.New(newTestCoreClient(t, srv.URL))
	out, err := cl.GetOrderDetail(context.Background(), "O-5")
	if err != nil {
		t.Fatalf("GetOrderDetail: %v", err)
	}
	if captured.Method != http.MethodGet {
		t.Errorf("method = %s, want GET", captured.Method)
	}
	if captured.Path != "/broker/orders/get?order_id=O-5" {
		t.Errorf("path = %s, want /broker/orders/get?order_id=O-5", captured.Path)
	}
	if out.OrderID != "O-5" {
		t.Errorf("order_id = %s, want O-5", out.OrderID)
	}
	if out.AvgFillPrice != "200.00" {
		t.Errorf("avg_fill_price = %s, want 200.00", out.AvgFillPrice)
	}
}

func TestGetOrderHistory(t *testing.T) {
	t.Parallel()
	var captured requestInfo
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		captured = requestInfo{Method: r.Method, Path: r.URL.RequestURI(), Body: body}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"order_id":"O-6","account_id":"ACC-1","symbol":"AAPL","status":"FILLED"},{"order_id":"O-7","account_id":"ACC-1","symbol":"MSFT","status":"CANCELLED"}]`))
	}))
	defer srv.Close()

	cl := broker.New(newTestCoreClient(t, srv.URL))
	out, err := cl.GetOrderHistory(context.Background(), "ACC-1")
	if err != nil {
		t.Fatalf("GetOrderHistory: %v", err)
	}
	if captured.Method != http.MethodGet {
		t.Errorf("method = %s, want GET", captured.Method)
	}
	if captured.Path != "/broker/orders/history?account_id=ACC-1" {
		t.Errorf("path = %s, want /broker/orders/history?account_id=ACC-1", captured.Path)
	}
	if len(out) != 2 {
		t.Fatalf("len(out) = %d, want 2", len(out))
	}
	if out[0].OrderID != "O-6" {
		t.Errorf("out[0].order_id = %s, want O-6", out[0].OrderID)
	}
	if out[1].Status != "CANCELLED" {
		t.Errorf("out[1].status = %s, want CANCELLED", out[1].Status)
	}
}

func TestGetOpenOrders(t *testing.T) {
	t.Parallel()
	var captured requestInfo
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		captured = requestInfo{Method: r.Method, Path: r.URL.RequestURI(), Body: body}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"order_id":"O-8","account_id":"ACC-1","symbol":"GOOG","status":"SUBMITTED"}]`))
	}))
	defer srv.Close()

	cl := broker.New(newTestCoreClient(t, srv.URL))
	out, err := cl.GetOpenOrders(context.Background(), "ACC-1")
	if err != nil {
		t.Fatalf("GetOpenOrders: %v", err)
	}
	if captured.Method != http.MethodGet {
		t.Errorf("method = %s, want GET", captured.Method)
	}
	if captured.Path != "/broker/orders/open?account_id=ACC-1" {
		t.Errorf("path = %s, want /broker/orders/open?account_id=ACC-1", captured.Path)
	}
	if len(out) != 1 {
		t.Fatalf("len(out) = %d, want 1", len(out))
	}
	if out[0].OrderID != "O-8" {
		t.Errorf("order_id = %s, want O-8", out[0].OrderID)
	}
}
