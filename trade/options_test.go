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
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/internal/errs"
	"github.com/shing1211/webullapi4go/trade"
)

// validOptionLeg returns a well-formed single option leg. Tests mutate a copy to
// isolate a single leg rule.
func validOptionLeg() trade.OrderLeg {
	return trade.OrderLeg{
		InstrumentType:   trade.InstrumentTypeOption,
		Market:           trade.MarketUS,
		Symbol:           "AAPL",
		Side:             trade.OrderSideBuy,
		StrikePrice:      "220.00",
		OptionExpireDate: "2026-12-18",
		OptionType:       trade.OptionTypeCall,
		Quantity:         "1",
	}
}

// validOptionOrder returns a well-formed single-leg buy-to-open AAPL call limit
// order as documented for options.
func validOptionOrder() trade.OrderRequest {
	return trade.OrderRequest{
		ClientOrderID:  "option-order-1",
		ComboType:      trade.ComboTypeNormal,
		InstrumentType: trade.InstrumentTypeOption,
		Market:         trade.MarketUS,
		Symbol:         "AAPL",
		OrderType:      trade.OrderTypeLimit,
		Side:           trade.OrderSideBuy,
		Quantity:       "1",
		EntrustType:    trade.EntrustTypeQty,
		TimeInForce:    trade.TimeInForceDay,
		LimitPrice:     "11.25",
		OptionStrategy: trade.OptionStrategySingle,
		Legs:           []trade.OrderLeg{validOptionLeg()},
	}
}

// validOptionPlaceRequest returns a well-formed single-order option place
// request.
func validOptionPlaceRequest() trade.PlaceOrderRequest {
	return trade.PlaceOrderRequest{
		AccountID: "ACC1",
		NewOrders: []trade.OrderRequest{validOptionOrder()},
	}
}

// TestOptionRequestValidate exercises every option-specific rule. Cases that set
// wantNil are expected to pass; every other case must fail with
// [errs.CodeInvalidConfig] before any network call.
func TestOptionRequestValidate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		mutate  func(*trade.OrderRequest)
		wantNil bool
	}{
		{"valid buy call limit DAY", func(o *trade.OrderRequest) {}, true},
		{"valid buy put GTC", func(o *trade.OrderRequest) {
			o.TimeInForce = trade.TimeInForceGTC
			o.Legs[0].OptionType = trade.OptionTypePut
		}, true},
		{"valid sell call limit DAY", func(o *trade.OrderRequest) {
			o.Side = trade.OrderSideSell
			o.Legs[0].Side = trade.OrderSideSell
		}, true},
		{"valid stop loss sell DAY", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeStopLoss
			o.LimitPrice = ""
			o.StopPrice = "3.00"
			o.Side = trade.OrderSideSell
			o.Legs[0].Side = trade.OrderSideSell
		}, true},
		{"valid stop loss limit sell DAY", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeStopLossLimit
			o.StopPrice = "4.00"
			o.LimitPrice = "3.80"
			o.Side = trade.OrderSideSell
			o.Legs[0].Side = trade.OrderSideSell
		}, true},
		{"market rejected", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeMarket
			o.LimitPrice = ""
		}, false},
		{"unsupported order type rejected", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeEnhancedLimit
		}, false},
		{"missing option_strategy", func(o *trade.OrderRequest) {
			o.OptionStrategy = ""
		}, false},
		{"unknown option_strategy", func(o *trade.OrderRequest) {
			o.OptionStrategy = "VERTICAL"
		}, false},
		{"short side rejected", func(o *trade.OrderRequest) {
			o.Side = trade.OrderSideShort
		}, false},
		{"sell gtc rejected", func(o *trade.OrderRequest) {
			o.Side = trade.OrderSideSell
			o.Legs[0].Side = trade.OrderSideSell
			o.TimeInForce = trade.TimeInForceGTC
		}, false},
		{"sell gtd rejected", func(o *trade.OrderRequest) {
			o.Side = trade.OrderSideSell
			o.Legs[0].Side = trade.OrderSideSell
			o.TimeInForce = trade.TimeInForceGTD
			o.ExpireDate = "2026-12-18"
		}, false},
		{"buy gtd rejected", func(o *trade.OrderRequest) {
			o.TimeInForce = trade.TimeInForceGTD
			o.ExpireDate = "2026-12-18"
		}, false},
		{"no legs", func(o *trade.OrderRequest) {
			o.Legs = nil
		}, false},
		{"two legs", func(o *trade.OrderRequest) {
			second := validOptionLeg()
			second.StrikePrice = "230.00"
			o.Legs = append(o.Legs, second)
		}, false},
		{"limit missing limit_price", func(o *trade.OrderRequest) {
			o.LimitPrice = ""
		}, false},
		{"stop loss missing stop_price", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeStopLoss
			o.LimitPrice = ""
		}, false},
		{"stop loss limit missing stop_price", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeStopLossLimit
			o.LimitPrice = "3.80"
		}, false},
		{"leg wrong instrument_type", func(o *trade.OrderRequest) {
			o.Legs[0].InstrumentType = trade.InstrumentTypeEquity
		}, false},
		{"leg wrong market", func(o *trade.OrderRequest) {
			o.Legs[0].Market = trade.MarketHK
		}, false},
		{"leg missing symbol", func(o *trade.OrderRequest) {
			o.Legs[0].Symbol = "  "
		}, false},
		{"leg short side rejected", func(o *trade.OrderRequest) {
			o.Legs[0].Side = trade.OrderSideShort
		}, false},
		{"leg missing strike_price", func(o *trade.OrderRequest) {
			o.Legs[0].StrikePrice = ""
		}, false},
		{"leg non-positive strike_price", func(o *trade.OrderRequest) {
			o.Legs[0].StrikePrice = "0"
		}, false},
		{"leg missing option_expire_date", func(o *trade.OrderRequest) {
			o.Legs[0].OptionExpireDate = ""
		}, false},
		{"leg malformed option_expire_date", func(o *trade.OrderRequest) {
			o.Legs[0].OptionExpireDate = "2026/12/18"
		}, false},
		{"leg unpadded option_expire_date", func(o *trade.OrderRequest) {
			o.Legs[0].OptionExpireDate = "2026-1-8"
		}, false},
		{"leg invalid calendar option_expire_date", func(o *trade.OrderRequest) {
			o.Legs[0].OptionExpireDate = "2026-13-40"
		}, false},
		{"leg missing option_type", func(o *trade.OrderRequest) {
			o.Legs[0].OptionType = ""
		}, false},
		{"leg unknown option_type", func(o *trade.OrderRequest) {
			o.Legs[0].OptionType = "STRADDLE"
		}, false},
		{"leg missing quantity", func(o *trade.OrderRequest) {
			o.Legs[0].Quantity = ""
		}, false},
		{"leg zero quantity", func(o *trade.OrderRequest) {
			o.Legs[0].Quantity = "0"
		}, false},
		{"option_strategy on equity rejected", func(o *trade.OrderRequest) {
			e := validOrder()
			e.OptionStrategy = trade.OptionStrategySingle
			*o = e
		}, false},
		{"legs on equity rejected", func(o *trade.OrderRequest) {
			e := validOrder()
			e.Legs = []trade.OrderLeg{validOptionLeg()}
			*o = e
		}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			o := validOptionOrder()
			tc.mutate(&o)
			err := o.Validate()
			if tc.wantNil {
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

// TestOrderLegValidate covers the exported single-leg validator directly.
func TestOrderLegValidate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		mutate  func(*trade.OrderLeg)
		wantNil bool
	}{
		{"valid call", func(l *trade.OrderLeg) {}, true},
		{"valid put", func(l *trade.OrderLeg) { l.OptionType = trade.OptionTypePut }, true},
		{"wrong instrument_type", func(l *trade.OrderLeg) { l.InstrumentType = trade.InstrumentTypeFutures }, false},
		{"wrong market", func(l *trade.OrderLeg) { l.Market = trade.MarketCN }, false},
		{"bad date", func(l *trade.OrderLeg) { l.OptionExpireDate = "18-12-2026" }, false},
		{"short side", func(l *trade.OrderLeg) { l.Side = trade.OrderSideShort }, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := validOptionLeg()
			tc.mutate(&l)
			err := l.Validate()
			if tc.wantNil {
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

// TestOptionOrderSerializesLegs posts an option order to the preview and place
// endpoints through httptest and asserts the serialized legs array and the
// account/identity fields, then checks the decoded response.
func TestOptionOrderSerializesLegs(t *testing.T) {
	t.Parallel()

	wantLeg := `"legs":[{"instrument_type":"OPTION","market":"US","symbol":"AAPL","side":"BUY","strike_price":"220.00","option_expire_date":"2026-12-18","option_type":"CALL","quantity":"1"}]`

	endpoints := []struct {
		name string
		path string
		call func(*trade.Client, context.Context, trade.PlaceOrderRequest) error
	}{
		{
			name: "preview",
			path: "/trading/orders/preview",
			call: func(c *trade.Client, ctx context.Context, req trade.PlaceOrderRequest) error {
				_, err := c.PreviewOrder(ctx, req)
				return err
			},
		},
		{
			name: "place",
			path: "/trading/orders/place",
			call: func(c *trade.Client, ctx context.Context, req trade.PlaceOrderRequest) error {
				_, err := c.PlaceOrder(ctx, req)
				return err
			},
		},
	}

	for _, ep := range endpoints {
		t.Run(ep.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("method = %q, want POST", r.Method)
				}
				if got := r.URL.Path; got != ep.path {
					t.Errorf("path = %q, want %q", got, ep.path)
				}
				if got := r.Header.Get("x-version"); got != client.APIVersionV3 {
					t.Errorf("x-version = %q, want %q", got, client.APIVersionV3)
				}
				raw := readBody(t, r)
				if !strings.Contains(string(raw), wantLeg) {
					t.Errorf("body missing serialized legs:\n got %s\nwant %s", raw, wantLeg)
				}
				for _, want := range []string{
					`"option_strategy":"SINGLE"`,
					`"instrument_type":"OPTION"`,
					`"order_type":"LIMIT"`,
					`"limit_price":"11.25"`,
				} {
					if !strings.Contains(string(raw), want) {
						t.Errorf("body missing %s: %s", want, raw)
					}
				}

				var body trade.PlaceOrderRequest
				if err := json.Unmarshal(raw, &body); err != nil {
					t.Errorf("decode body: %v", err)
				}
				if body.AccountID != "ACC1" || len(body.NewOrders) != 1 || len(body.NewOrders[0].Legs) != 1 {
					t.Errorf("body = %+v, want ACC1 with one single-leg order", body)
				}
				leg := body.NewOrders[0].Legs[0]
				if leg.Symbol != "AAPL" || leg.StrikePrice != "220.00" || leg.OptionExpireDate != "2026-12-18" ||
					leg.OptionType != trade.OptionTypeCall || leg.Side != trade.OrderSideBuy || leg.Market != trade.MarketUS {
					t.Errorf("leg = %+v", leg)
				}

				w.Header().Set("Content-Type", "application/json")
				if ep.name == "preview" {
					_, _ = w.Write([]byte(`{"estimated_cost":"1125.00","estimated_transaction_fee":"1.00"}`))
					return
				}
				_, _ = w.Write([]byte(`{"client_order_id":"option-order-1","order_id":"OID-OPT-1"}`))
			}))
			defer srv.Close()

			c := newOrderTestClient(t, srv.URL)
			if err := ep.call(c, context.Background(), validOptionPlaceRequest()); err != nil {
				t.Fatalf("%s option order error = %v", ep.name, err)
			}
		})
	}
}

// TestOptionValidationPrecedesNetwork ensures an invalid option order is
// rejected before any request reaches the server.
func TestOptionValidationPrecedesNetwork(t *testing.T) {
	t.Parallel()

	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	req := validOptionPlaceRequest()
	req.NewOrders[0].OrderType = trade.OrderTypeMarket

	c := newOrderTestClient(t, srv.URL)
	if _, err := c.PreviewOrder(context.Background(), req); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("PreviewOrder() error = %v, want invalid_config", err)
	}
	if _, err := c.PlaceOrder(context.Background(), req); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("PlaceOrder() error = %v, want invalid_config", err)
	}
	if hits != 0 {
		t.Fatalf("server received %d requests, want 0 (validation must precede the network)", hits)
	}
}

// TestSandboxPreviewOption previews a small AAPL call option on the dedicated
// trading sandbox account. It is read-only: preview never places or mutates an
// order. The sandbox may not hold the option contract, in which case the API
// returns a typed business error (for example 417 Invalid Symbol) that this test
// tolerates and logs rather than failing, so the option payload can be observed
// live without a guaranteed contract.
//
// It is gated by WEBULL_TRADE_SANDBOX=1 plus WEBULL_TRADE_APP_KEY,
// WEBULL_TRADE_APP_SECRET, and WEBULL_TRADE_ACCOUNT_ID; credentials are read
// from the environment only and are never committed.
func TestSandboxPreviewOption(t *testing.T) {
	trading, accountID := newSandboxTradeClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := validOptionPlaceRequest()
	req.AccountID = accountID
	req.NewOrders[0].ClientOrderID = "sdk-preview-option-aapl-1"
	req.NewOrders[0].Symbol = "AAPL"
	req.NewOrders[0].LimitPrice = "1.00"
	req.NewOrders[0].Legs[0].StrikePrice = "220.00"
	req.NewOrders[0].Legs[0].OptionExpireDate = "2026-12-18"

	res, err := trading.PreviewOrder(ctx, req)
	if err != nil {
		var e *errs.Error
		if !errors.As(err, &e) {
			t.Fatalf("PreviewOrder(option) error = %v (%T), want *errs.Error", err, err)
		}
		t.Logf("PreviewOrder(option) typed error (tolerated): code=%s message=%s", e.Code, e.Message)
		return
	}
	t.Logf("preview option estimated_cost=%s estimated_transaction_fee=%s",
		res.EstimatedCost, res.EstimatedTransactionFee)
}
