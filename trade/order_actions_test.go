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
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/trade"
)

// validReplaceRequest returns a well-formed single-order replace request.
func validReplaceRequest() trade.ReplaceOrderRequest {
	return trade.ReplaceOrderRequest{
		AccountID: "ACC1",
		ModifyOrders: []trade.ModifyOrderRequest{
			{
				ClientOrderID: "test-order-1",
				Quantity:      "2",
				LimitPrice:    "175.50",
			},
		},
	}
}

// validCancelRequest returns a well-formed cancel request.
func validCancelRequest() trade.CancelOrderRequest {
	return trade.CancelOrderRequest{
		AccountID:     "ACC1",
		ClientOrderID: "test-order-1",
	}
}

func TestReplaceOrder(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/trading/orders/replace"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("Content-Type"), "application/json"; got != want {
			t.Errorf("Content-Type = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("x-version"), client.APIVersionV3; got != want {
			t.Errorf("x-version = %q, want %q", got, want)
		}
		var body trade.ReplaceOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if body.AccountID != "ACC1" || len(body.ModifyOrders) != 1 {
			t.Errorf("body = %+v, want account ACC1 with 1 modify order", body)
		}
		o := body.ModifyOrders[0]
		if o.ClientOrderID != "test-order-1" || o.Quantity != "2" || o.LimitPrice != "175.50" {
			t.Errorf("modify order = %+v", o)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"client_order_id":"test-order-1","order_id":"OID-1"}`))
	}))
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)
	got, err := c.ReplaceOrder(context.Background(), validReplaceRequest())
	if err != nil {
		t.Fatalf("ReplaceOrder() error = %v", err)
	}
	if got.ClientOrderID != "test-order-1" || got.OrderID != "OID-1" {
		t.Errorf("ReplaceOrder() = %+v", got)
	}
}

func TestCancelOrder(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/trading/orders/cancel"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("Content-Type"), "application/json"; got != want {
			t.Errorf("Content-Type = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("x-version"), client.APIVersionV3; got != want {
			t.Errorf("x-version = %q, want %q", got, want)
		}
		var body trade.CancelOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if body.AccountID != "ACC1" || body.ClientOrderID != "test-order-1" {
			t.Errorf("body = %+v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"client_order_id":"test-order-1","order_id":"OID-1"}`))
	}))
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)
	got, err := c.CancelOrder(context.Background(), validCancelRequest())
	if err != nil {
		t.Fatalf("CancelOrder() error = %v", err)
	}
	if got.ClientOrderID != "test-order-1" || got.OrderID != "OID-1" {
		t.Errorf("CancelOrder() = %+v", got)
	}
}

func TestReplaceOrderSerializesOnlySetFields(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var raw map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			t.Errorf("decode body: %v", err)
		}
		items, ok := raw["modify_orders"]
		if !ok {
			t.Fatalf("body missing modify_orders: %v", raw)
		}
		var orders []map[string]json.RawMessage
		if err := json.Unmarshal(items, &orders); err != nil {
			t.Fatalf("decode modify_orders: %v", err)
		}
		if len(orders) != 1 {
			t.Fatalf("modify_orders len = %d, want 1", len(orders))
		}
		if _, ok := orders[0]["stop_price"]; ok {
			t.Errorf("stop_price should be omitted when unset: %v", orders[0])
		}
		if _, ok := orders[0]["expire_date"]; ok {
			t.Errorf("expire_date should be omitted when unset: %v", orders[0])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"client_order_id":"test-order-1","order_id":"OID-1"}`))
	}))
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)
	if _, err := c.ReplaceOrder(context.Background(), validReplaceRequest()); err != nil {
		t.Fatalf("ReplaceOrder() error = %v", err)
	}
}

func TestReplaceOrderRequestValidate(t *testing.T) {
	t.Parallel()

	overlong := strings.Repeat("a", 33)

	cases := []struct {
		name    string
		mutate  func(*trade.ReplaceOrderRequest)
		wantErr bool
	}{
		{"valid", func(r *trade.ReplaceOrderRequest) {}, false},
		{"client_order_id at max length", func(r *trade.ReplaceOrderRequest) {
			r.ModifyOrders[0].ClientOrderID = strings.Repeat("a", 32)
		}, false},
		{"missing account_id", func(r *trade.ReplaceOrderRequest) { r.AccountID = "" }, true},
		{"blank account_id", func(r *trade.ReplaceOrderRequest) { r.AccountID = "   " }, true},
		{"no modify orders", func(r *trade.ReplaceOrderRequest) { r.ModifyOrders = nil }, true},
		{"missing client_order_id", func(r *trade.ReplaceOrderRequest) { r.ModifyOrders[0].ClientOrderID = "" }, true},
		{"client_order_id too long", func(r *trade.ReplaceOrderRequest) { r.ModifyOrders[0].ClientOrderID = overlong }, true},
		{"client_order_id bad charset", func(r *trade.ReplaceOrderRequest) { r.ModifyOrders[0].ClientOrderID = "bad id!" }, true},
		{"invalid time_in_force", func(r *trade.ReplaceOrderRequest) { r.ModifyOrders[0].TimeInForce = "FOK" }, true},
		{"invalid trailing_type", func(r *trade.ReplaceOrderRequest) { r.ModifyOrders[0].TrailingType = "BOGUS" }, true},
		{"invalid trigger_price_type", func(r *trade.ReplaceOrderRequest) { r.ModifyOrders[0].TriggerPriceType = "LAST" }, true},
		{"duplicate client_order_id", func(r *trade.ReplaceOrderRequest) {
			second := r.ModifyOrders[0]
			r.ModifyOrders = append(r.ModifyOrders, second)
		}, true},
		{"distinct client_order_id", func(r *trade.ReplaceOrderRequest) {
			second := r.ModifyOrders[0]
			second.ClientOrderID = "test-order-2"
			r.ModifyOrders = append(r.ModifyOrders, second)
		}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := validReplaceRequest()
			tc.mutate(&r)
			err := r.Validate()
			if tc.wantErr {
				if !errs.Is(err, errs.CodeInvalidConfig) {
					t.Fatalf("Validate() error = %v, want invalid_config", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Validate() error = %v, want nil", err)
			}
		})
	}
}

func TestCancelOrderRequestValidate(t *testing.T) {
	t.Parallel()

	overlong := strings.Repeat("a", 33)

	cases := []struct {
		name    string
		mutate  func(*trade.CancelOrderRequest)
		wantErr bool
	}{
		{"valid", func(r *trade.CancelOrderRequest) {}, false},
		{"client_order_id at max length", func(r *trade.CancelOrderRequest) {
			r.ClientOrderID = strings.Repeat("a", 32)
		}, false},
		{"missing account_id", func(r *trade.CancelOrderRequest) { r.AccountID = "" }, true},
		{"blank account_id", func(r *trade.CancelOrderRequest) { r.AccountID = "  " }, true},
		{"missing client_order_id", func(r *trade.CancelOrderRequest) { r.ClientOrderID = "" }, true},
		{"client_order_id too long", func(r *trade.CancelOrderRequest) { r.ClientOrderID = overlong }, true},
		{"client_order_id bad charset", func(r *trade.CancelOrderRequest) { r.ClientOrderID = "bad id!" }, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := validCancelRequest()
			tc.mutate(&r)
			err := r.Validate()
			if tc.wantErr {
				if !errs.Is(err, errs.CodeInvalidConfig) {
					t.Fatalf("Validate() error = %v, want invalid_config", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Validate() error = %v, want nil", err)
			}
		})
	}
}

func TestOrderActionValidationPrecedesNetwork(t *testing.T) {
	t.Parallel()

	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)

	badReplace := validReplaceRequest()
	badReplace.ModifyOrders[0].ClientOrderID = ""
	if _, err := c.ReplaceOrder(context.Background(), badReplace); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("ReplaceOrder() error = %v, want invalid_config", err)
	}

	badCancel := validCancelRequest()
	badCancel.ClientOrderID = ""
	if _, err := c.CancelOrder(context.Background(), badCancel); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("CancelOrder() error = %v, want invalid_config", err)
	}

	if got := hits.Load(); got != 0 {
		t.Fatalf("server received %d requests, want 0 (validation must precede the network)", got)
	}
}
