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

package trade_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	errs "github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/trade"
)

func TestBatchPlaceOrder(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/trading/orders/batch-place"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("Content-Type"), "application/json"; got != want {
			t.Errorf("Content-Type = %q, want %q", got, want)
		}
		var body trade.PlaceOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if body.AccountID != "ACC1" {
			t.Errorf("body.AccountID = %q, want ACC1", body.AccountID)
		}
		if len(body.NewOrders) != 2 {
			t.Errorf("len(body.NewOrders) = %d, want 2", len(body.NewOrders))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[{"client_order_id":"bo-1","order_id":"OID-B1"},{"client_order_id":"bo-2","order_id":"OID-B2"}]}`))
	}))
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)

	second := validOrder()
	second.ClientOrderID = "bo-2"
	req := trade.PlaceOrderRequest{
		AccountID: "ACC1",
		NewOrders: []trade.OrderRequest{
			func() trade.OrderRequest { o := validOrder(); o.ClientOrderID = "bo-1"; return o }(),
			second,
		},
	}

	got, err := c.BatchPlaceOrder(context.Background(), req)
	if err != nil {
		t.Fatalf("BatchPlaceOrder() error = %v", err)
	}
	if len(got.Results) != 2 {
		t.Fatalf("len(Results) = %d, want 2", len(got.Results))
	}
	if got.Results[0].ClientOrderID != "bo-1" || got.Results[0].OrderID != "OID-B1" {
		t.Errorf("Results[0] = %+v", got.Results[0])
	}
	if got.Results[1].ClientOrderID != "bo-2" || got.Results[1].OrderID != "OID-B2" {
		t.Errorf("Results[1] = %+v", got.Results[1])
	}
}

func TestBatchPlaceOrderOverLimit(t *testing.T) {
	t.Parallel()

	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)

	orders := make([]trade.OrderRequest, 51)
	for i := range orders {
		o := validOrder()
		o.ClientOrderID = fmt.Sprintf("bo-%d", i+1)
		orders[i] = o
	}
	req := trade.PlaceOrderRequest{
		AccountID: "ACC1",
		NewOrders: orders,
	}

	_, err := c.BatchPlaceOrder(context.Background(), req)
	if !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("BatchPlaceOrder() error = %v, want invalid_config", err)
	}
	if hits != 0 {
		t.Fatalf("server received %d requests, want 0", hits)
	}
}

func TestBatchPlaceOrderNonEquity(t *testing.T) {
	t.Parallel()

	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)

	o := validOrder()
	o.InstrumentType = trade.InstrumentTypeOption
	o.OptionStrategy = trade.OptionStrategySingle
	o.Side = trade.OrderSideBuy
	o.TimeInForce = trade.TimeInForceGTC
	o.Legs = []trade.OrderLeg{{
		InstrumentType:   trade.InstrumentTypeOption,
		Market:           trade.MarketUS,
		Symbol:           "AAPL",
		Side:             trade.OrderSideBuy,
		StrikePrice:      mp("180.00"),
		OptionExpireDate: "2026-12-19",
		OptionType:       trade.OptionTypeCall,
		Quantity:         mp("1"),
	}}
	req := trade.PlaceOrderRequest{
		AccountID: "ACC1",
		NewOrders: []trade.OrderRequest{o},
	}

	_, err := c.BatchPlaceOrder(context.Background(), req)
	if !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("BatchPlaceOrder() error = %v, want invalid_config", err)
	}
	if !strings.Contains(err.Error(), "EQUITY") {
		t.Errorf("error message should mention EQUITY: %v", err)
	}
	if hits != 0 {
		t.Fatalf("server received %d requests, want 0", hits)
	}
}

func TestBatchPlaceOrderComboType(t *testing.T) {
	t.Parallel()

	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)

	o := validOrder()
	o.ComboType = trade.ComboTypeMaster
	req := trade.PlaceOrderRequest{
		AccountID: "ACC1",
		NewOrders: []trade.OrderRequest{o},
	}

	_, err := c.BatchPlaceOrder(context.Background(), req)
	if !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("BatchPlaceOrder() error = %v, want invalid_config", err)
	}
	if !strings.Contains(err.Error(), "NORMAL") {
		t.Errorf("error message should mention NORMAL: %v", err)
	}
	if hits != 0 {
		t.Fatalf("server received %d requests, want 0", hits)
	}
}
