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
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/pkg/domain/money"
	"github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/trade"
)

// newOrderTestClient returns a trade.Client with no guardrails that targets
// baseURL.
func newOrderTestClient(t *testing.T, baseURL string) *trade.Client {
	t.Helper()
	return newTradeClient(t, baseURL)
}

// newTradeClient returns a trade.Client configured with opts that targets
// baseURL, using the shared test credentials.
func newTradeClient(t *testing.T, baseURL string, opts ...trade.Option) *trade.Client {
	t.Helper()
	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithBaseURL(baseURL),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })
	return trade.New(cl, opts...)
}

// validOrder returns a well-formed simple US equity limit order. Tests mutate a
// copy to isolate a single validation rule.
func validOrder() trade.OrderRequest {
	return trade.OrderRequest{
		ClientOrderID:         "test-order-1",
		ComboType:             trade.ComboTypeNormal,
		InstrumentType:        trade.InstrumentTypeEquity,
		Market:                trade.MarketUS,
		Symbol:                "AAPL",
		OrderType:             trade.OrderTypeLimit,
		Side:                  trade.OrderSideBuy,
		Quantity:              mp("1"),
		EntrustType:           trade.EntrustTypeQty,
		TimeInForce:           trade.TimeInForceDay,
		SupportTradingSession: trade.TradingSessionCore,
		LimitPrice:            mp("180.00"),
	}
}

// validPlaceRequest returns a well-formed single-order place/preview request.
func validPlaceRequest() trade.PlaceOrderRequest {
	return trade.PlaceOrderRequest{
		AccountID: "ACC1",
		NewOrders: []trade.OrderRequest{validOrder()},
	}
}

func TestPreviewOrder(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/trading/orders/preview"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("Content-Type"), "application/json"; got != want {
			t.Errorf("Content-Type = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("x-version"), client.APIVersionV3; got != want {
			t.Errorf("x-version = %q, want %q", got, want)
		}
		raw := readBody(t, r)
		var body trade.PlaceOrderRequest
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if body.AccountID != "ACC1" || len(body.NewOrders) != 1 {
			t.Errorf("body = %+v, want account ACC1 with 1 order", body)
		}
		o := body.NewOrders[0]
		if o.ClientOrderID != "test-order-1" || o.Symbol != "AAPL" || o.Market != trade.MarketUS {
			t.Errorf("order identity = %+v", o)
		}
		if o.OrderType != trade.OrderTypeLimit || o.LimitPrice.Cmp(money.Must(money.NewFromString("180.00"))) != 0 || o.Quantity.Cmp(money.Must(money.NewFromString("1"))) != 0 {
			t.Errorf("order terms = %+v", o)
		}
		if o.ComboType != trade.ComboTypeNormal || o.EntrustType != trade.EntrustTypeQty || o.Side != trade.OrderSideBuy {
			t.Errorf("order enums = %+v", o)
		}
		if !strings.Contains(string(raw), `"support_trading_session":"CORE"`) {
			t.Errorf("body missing support_trading_session: %s", raw)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"estimated_cost":"180.00","estimated_transaction_fee":"1.00"}`))
	}))
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)
	got, err := c.PreviewOrder(context.Background(), validPlaceRequest())
	if err != nil {
		t.Fatalf("PreviewOrder() error = %v", err)
	}
	if got.EstimatedCost.Cmp(money.Must(money.NewFromString("180.00"))) != 0 || got.EstimatedTransactionFee.Cmp(money.Must(money.NewFromString("1.00"))) != 0 {
		t.Errorf("PreviewOrder() = %+v", got)
	}
}

func TestPlaceOrder(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/trading/orders/place"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("x-version"), client.APIVersionV3; got != want {
			t.Errorf("x-version = %q, want %q", got, want)
		}
		var body trade.PlaceOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if body.AccountID != "ACC1" || len(body.NewOrders) != 1 {
			t.Errorf("body = %+v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"client_order_id":"test-order-1","order_id":"OID-1"}`))
	}))
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)
	got, err := c.PlaceOrder(context.Background(), validPlaceRequest())
	if err != nil {
		t.Fatalf("PlaceOrder() error = %v", err)
	}
	if got.ClientOrderID != "test-order-1" || got.OrderID != "OID-1" {
		t.Errorf("PlaceOrder() = %+v", got)
	}
}

func TestPlaceOrderSerializesComboAndAmountFields(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := string(readBody(t, r))
		for _, want := range []string{
			`"client_combo_order_id":"combo-1"`,
			`"total_cash_amount":"100.4"`,
			`"support_trading_session":"CORE"`,
		} {
			if !strings.Contains(raw, want) {
				t.Errorf("body missing %s: %s", want, raw)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"client_order_id":"test-order-1","order_id":"OID-1"}`))
	}))
	defer srv.Close()

	req := validPlaceRequest()
	req.ClientComboOrderID = "combo-1"
	order := validOrder()
	order.ComboType = trade.ComboTypeMaster
	order.EntrustType = trade.EntrustTypeAmount
	order.TotalCashAmount = mp("100.40")
	order.Quantity = nil
	req.NewOrders = []trade.OrderRequest{order}

	c := newOrderTestClient(t, srv.URL)
	if _, err := c.PlaceOrder(context.Background(), req); err != nil {
		t.Fatalf("PlaceOrder() error = %v", err)
	}
}

func TestOrderRequestValidate(t *testing.T) {
	t.Parallel()

	overlong := strings.Repeat("a", 33)
	longest := strings.Repeat("a", 32)

	cases := []struct {
		name   string
		mutate func(*trade.OrderRequest)
	}{
		{"valid", func(o *trade.OrderRequest) {}},
		{"client order id at max length", func(o *trade.OrderRequest) { o.ClientOrderID = longest }},
		{"missing client_order_id", func(o *trade.OrderRequest) { o.ClientOrderID = "" }},
		{"client_order_id too long", func(o *trade.OrderRequest) { o.ClientOrderID = overlong }},
		{"client_order_id bad charset", func(o *trade.OrderRequest) { o.ClientOrderID = "bad id!" }},
		{"missing combo_type", func(o *trade.OrderRequest) { o.ComboType = "" }},
		{"invalid combo_type", func(o *trade.OrderRequest) { o.ComboType = "BOGUS" }},
		{"invalid instrument_type", func(o *trade.OrderRequest) { o.InstrumentType = "BOND" }},
		{"invalid market", func(o *trade.OrderRequest) { o.Market = "JP" }},
		{"missing symbol", func(o *trade.OrderRequest) { o.Symbol = "  " }},
		{"invalid order_type", func(o *trade.OrderRequest) { o.OrderType = "BOGUS" }},
		{"invalid side", func(o *trade.OrderRequest) { o.Side = "HOLD" }},
		{"invalid time_in_force", func(o *trade.OrderRequest) { o.TimeInForce = "FOK" }},
		{"invalid support_trading_session", func(o *trade.OrderRequest) { o.SupportTradingSession = "WEEKEND" }},
		{"missing quantity for QTY", func(o *trade.OrderRequest) { o.Quantity = nil }},
		{"zero quantity", func(o *trade.OrderRequest) { o.Quantity = mp("0") }},
		{"negative quantity", func(o *trade.OrderRequest) { o.Quantity = mp("-1") }},
		{"invalid entrust_type", func(o *trade.OrderRequest) { o.EntrustType = "SHARES" }},
		{"AMOUNT missing total_cash_amount", func(o *trade.OrderRequest) {
			o.EntrustType = trade.EntrustTypeAmount
			o.TotalCashAmount = nil
		}},
		{"AMOUNT non-positive total_cash_amount", func(o *trade.OrderRequest) {
			o.EntrustType = trade.EntrustTypeAmount
			o.TotalCashAmount = mp("0")
		}},
		{"LIMIT missing limit_price", func(o *trade.OrderRequest) { o.LimitPrice = nil }},
		{"STOP_LOSS missing stop_price", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeStopLoss
			o.LimitPrice = nil
		}},
		{"STOP_LOSS_LIMIT missing limit_price", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeStopLossLimit
			o.StopPrice = mp("170.00")
			o.LimitPrice = nil
		}},
		{"STOP_LOSS_LIMIT missing stop_price", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeStopLossLimit
			o.LimitPrice = mp("170.00")
		}},
		{"TOUCH_MKT missing stop_price", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeTouchMkt
			o.LimitPrice = nil
		}},
		{"TOUCH_LMT missing stop_price", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeTouchLmt
			o.LimitPrice = mp("170.00")
		}},
		{"trailing missing trailing_type", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeTrailingStopLoss
			o.LimitPrice = nil
		}},
		{"trailing missing trailing_stop_step", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeTrailingStopLoss
			o.LimitPrice = nil
			o.TrailingType = trade.TrailingTypeAmount
		}},
		{"GTD missing expire_date", func(o *trade.OrderRequest) { o.TimeInForce = trade.TimeInForceGTD }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			o := validOrder()
			tc.mutate(&o)
			err := o.Validate()
			if tc.name == "valid" || tc.name == "client order id at max length" {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}
			if !errs.Is(err, errs.CodeInvalidConfig) {
				t.Fatalf("Validate() error = %v, want invalid_config", err)
			}
		})
	}
}

func TestPlaceOrderRequestValidate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		mutate  func(*trade.PlaceOrderRequest)
		wantErr bool
	}{
		{"valid", func(r *trade.PlaceOrderRequest) {}, false},
		{"missing account_id", func(r *trade.PlaceOrderRequest) { r.AccountID = "" }, true},
		{"blank account_id", func(r *trade.PlaceOrderRequest) { r.AccountID = "   " }, true},
		{"no orders", func(r *trade.PlaceOrderRequest) { r.NewOrders = nil }, true},
		{"order invalid", func(r *trade.PlaceOrderRequest) { r.NewOrders[0].Side = "" }, true},
		{"duplicate client_order_id", func(r *trade.PlaceOrderRequest) {
			second := validOrder()
			r.NewOrders = append(r.NewOrders, second)
		}, true},
		{"distinct client_order_id", func(r *trade.PlaceOrderRequest) {
			second := validOrder()
			second.ClientOrderID = "test-order-2"
			r.NewOrders = append(r.NewOrders, second)
		}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := validPlaceRequest()
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

func TestOrderValidationPrecedesNetwork(t *testing.T) {
	t.Parallel()

	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newOrderTestClient(t, srv.URL)
	bad := validPlaceRequest()
	bad.NewOrders[0].ClientOrderID = ""

	if _, err := c.PreviewOrder(context.Background(), bad); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("PreviewOrder() error = %v, want invalid_config", err)
	}
	if _, err := c.PlaceOrder(context.Background(), bad); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("PlaceOrder() error = %v, want invalid_config", err)
	}
	if got := hits.Load(); got != 0 {
		t.Fatalf("server received %d requests, want 0 (validation must precede the network)", got)
	}
}

func TestOrderGuardrails(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		opts    []trade.Option
		mutate  func(*trade.OrderRequest)
		wantErr bool
	}{
		{"quantity exceeded", []trade.Option{trade.WithMaxOrderQuantity("1")},
			func(o *trade.OrderRequest) { o.Quantity = mp("2") }, true},
		{"quantity at cap", []trade.Option{trade.WithMaxOrderQuantity("2")},
			func(o *trade.OrderRequest) { o.Quantity = mp("2") }, false},
		{"quantity fractional exceeded", []trade.Option{trade.WithMaxOrderQuantity("1.5")},
			func(o *trade.OrderRequest) { o.Quantity = mp("1.6") }, true},
		{"notional from limit exceeded", []trade.Option{trade.WithMaxOrderNotional("100")},
			func(o *trade.OrderRequest) { o.Quantity = mp("2"); o.LimitPrice = mp("60") }, true},
		{"notional from limit at cap", []trade.Option{trade.WithMaxOrderNotional("120")},
			func(o *trade.OrderRequest) { o.Quantity = mp("2"); o.LimitPrice = mp("60") }, false},
		{"notional from amount exceeded", []trade.Option{trade.WithMaxOrderNotional("100")},
			func(o *trade.OrderRequest) {
				o.EntrustType = trade.EntrustTypeAmount
				o.TotalCashAmount = mp("500")
				o.Quantity = nil
			}, true},
		{"market skips notional", []trade.Option{trade.WithMaxOrderNotional("1")},
			func(o *trade.OrderRequest) {
				o.OrderType = trade.OrderTypeMarket
				o.LimitPrice = nil
				o.Quantity = mp("1")
			}, false},
		{"market still enforces quantity", []trade.Option{
			trade.WithMaxOrderNotional("1"),
			trade.WithMaxOrderQuantity("1"),
		}, func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeMarket
			o.LimitPrice = nil
			o.Quantity = mp("2")
		}, true},
		{"guardrails disabled", nil,
			func(o *trade.OrderRequest) { o.Quantity = mp("1000000"); o.LimitPrice = mp("1000000") }, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var hits atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				hits.Add(1)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"estimated_cost":"1.00","estimated_transaction_fee":"0.10"}`))
			}))
			defer srv.Close()

			c := newTradeClient(t, srv.URL, tc.opts...)
			req := validPlaceRequest()
			tc.mutate(&req.NewOrders[0])

			_, err := c.PreviewOrder(context.Background(), req)
			if tc.wantErr {
				if !errs.Is(err, errs.CodeInvalidConfig) {
					t.Fatalf("PreviewOrder() error = %v, want invalid_config", err)
				}
				if got := hits.Load(); got != 0 {
					t.Fatalf("server received %d requests, want 0 (guardrail must precede the network)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("PreviewOrder() error = %v, want nil", err)
			}
		})
	}
}

func TestPlaceOrderEnforcesGuardrails(t *testing.T) {
	t.Parallel()

	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"client_order_id":"test-order-1","order_id":"OID-1"}`))
	}))
	defer srv.Close()

	c := newTradeClient(t, srv.URL, trade.WithMaxOrderQuantity("1"))
	req := validPlaceRequest()
	req.NewOrders[0].Quantity = mp("2")

	if _, err := c.PlaceOrder(context.Background(), req); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("PlaceOrder() error = %v, want invalid_config", err)
	}
	if got := hits.Load(); got != 0 {
		t.Fatalf("server received %d requests, want 0", got)
	}
}

// readBody reads and returns r's request body.
func readBody(t *testing.T, r *http.Request) []byte {
	t.Helper()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return data
}
