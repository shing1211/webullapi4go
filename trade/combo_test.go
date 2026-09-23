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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/trade"
)

// comboOrder returns a well-formed US equity order for a combo role. The prices
// required by ot are populated so a failure isolates the composition rule under
// test rather than a missing-field rule.
func comboOrder(id string, ct trade.ComboType, ot trade.OrderType, side trade.OrderSide) trade.OrderRequest {
	o := trade.OrderRequest{
		ClientOrderID:         id,
		ComboType:             ct,
		InstrumentType:        trade.InstrumentTypeEquity,
		Market:                trade.MarketUS,
		Symbol:                "AAPL",
		OrderType:             ot,
		Side:                  side,
		Quantity:              mp("1"),
		EntrustType:           trade.EntrustTypeQty,
		TimeInForce:           trade.TimeInForceDay,
		SupportTradingSession: trade.TradingSessionCore,
	}
	switch ot {
	case trade.OrderTypeLimit:
		o.LimitPrice = mp("180.00")
	case trade.OrderTypeStopLoss:
		o.StopPrice = mp("170.00")
	case trade.OrderTypeStopLossLimit:
		o.StopPrice = mp("170.00")
		o.LimitPrice = mp("169.00")
	}
	return o
}

// comboRequest groups orders under a client combo order identifier.
func comboRequest(comboID string, orders ...trade.OrderRequest) trade.PlaceOrderRequest {
	return trade.PlaceOrderRequest{
		AccountID:          "ACC1",
		ClientComboOrderID: comboID,
		NewOrders:          orders,
	}
}

// otoLegs returns n OTO follow-up orders.
func otoLegs(n int) []trade.OrderRequest {
	legs := make([]trade.OrderRequest, n)
	for i := range legs {
		legs[i] = comboOrder(fmt.Sprintf("oto-leg-%d", i), trade.ComboTypeOTO, trade.OrderTypeLimit, trade.OrderSideSell)
	}
	return legs
}

// ocoLegs returns n OCO orders cycling through the allowed order types.
func ocoLegs(n int) []trade.OrderRequest {
	types := []trade.OrderType{
		trade.OrderTypeLimit,
		trade.OrderTypeStopLoss,
		trade.OrderTypeStopLossLimit,
	}
	legs := make([]trade.OrderRequest, n)
	for i := range legs {
		legs[i] = comboOrder(fmt.Sprintf("oco-leg-%d", i), trade.ComboTypeOCO, types[i%len(types)], trade.OrderSideSell)
	}
	return legs
}

// otocoLegs returns n OTOCO trigger legs.
func otocoLegs(n int) []trade.OrderRequest {
	legs := make([]trade.OrderRequest, n)
	for i := range legs {
		legs[i] = comboOrder(fmt.Sprintf("otoco-leg-%d", i), trade.ComboTypeOTOCO, trade.OrderTypeLimit, trade.OrderSideSell)
	}
	return legs
}

// validOptionForCombo is a well-formed single-leg option order carrying a combo
// role, used to prove combo orders reject non-equity instruments.
func validOptionForCombo(ct trade.ComboType) trade.OrderRequest {
	return trade.OrderRequest{
		ClientOrderID:  "opt-1",
		ComboType:      ct,
		InstrumentType: trade.InstrumentTypeOption,
		Market:         trade.MarketUS,
		Symbol:         "AAPL",
		OrderType:      trade.OrderTypeLimit,
		Side:           trade.OrderSideBuy,
		Quantity:       mp("1"),
		EntrustType:    trade.EntrustTypeQty,
		TimeInForce:    trade.TimeInForceDay,
		LimitPrice:     mp("11.00"),
		OptionStrategy: trade.OptionStrategySingle,
		Legs: []trade.OrderLeg{{
			InstrumentType:   trade.InstrumentTypeOption,
			Market:           trade.MarketUS,
			Symbol:           "AAPL",
			Side:             trade.OrderSideBuy,
			StrikePrice:      mp("220.00"),
			OptionExpireDate: "2026-12-18",
			OptionType:       trade.OptionTypeCall,
			Quantity:         mp("1"),
		}},
	}
}

// validHKForCombo is a well-formed HK equity order carrying a combo role, used
// to prove combo orders reject non-US markets.
func validHKForCombo(ct trade.ComboType) trade.OrderRequest {
	return trade.OrderRequest{
		ClientOrderID:  "hk-1",
		ComboType:      ct,
		InstrumentType: trade.InstrumentTypeEquity,
		Market:         trade.MarketHK,
		Symbol:         "00700",
		OrderType:      trade.OrderTypeEnhancedLimit,
		Side:           trade.OrderSideBuy,
		Quantity:       mp("100"),
		EntrustType:    trade.EntrustTypeQty,
		TimeInForce:    trade.TimeInForceDay,
		LimitPrice:     mp("10.00"),
		NoPartyIDs:     []trade.PartyID{{PartyID: "ABC123.2568", PartyIDSource: "D", PartyRole: "3"}},
	}
}

func TestComboCompositionValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		build      func() trade.PlaceOrderRequest
		wantErr    bool
		wantSubstr string
	}{
		{
			name: "valid TP/SL buy-to-open",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-tpsl",
					comboOrder("m", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy),
					comboOrder("p", trade.ComboTypeStopProfit, trade.OrderTypeLimit, trade.OrderSideSell),
					comboOrder("l", trade.ComboTypeStopLoss, trade.OrderTypeStopLoss, trade.OrderSideSell),
				)
			},
		},
		{
			name: "valid TP/SL master only",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-tpsl",
					comboOrder("m", trade.ComboTypeMaster, trade.OrderTypeMarket, trade.OrderSideBuy),
				)
			},
		},
		{
			name: "valid TP/SL sell-to-close",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-tpsl",
					comboOrder("p", trade.ComboTypeStopProfit, trade.OrderTypeLimit, trade.OrderSideSell),
					comboOrder("l", trade.ComboTypeStopLoss, trade.OrderTypeStopLoss, trade.OrderSideSell),
				)
			},
		},
		{
			name: "valid TP/SL sell-to-close only stop profit",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-tpsl",
					comboOrder("p", trade.ComboTypeStopProfit, trade.OrderTypeLimit, trade.OrderSideSell),
				)
			},
		},
		{
			name: "valid OTO one leg",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-oto",
					comboOrder("m", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy),
					comboOrder("l", trade.ComboTypeOTO, trade.OrderTypeLimit, trade.OrderSideSell),
				)
			},
		},
		{
			name: "valid OTO stop master with six legs",
			build: func() trade.PlaceOrderRequest {
				m := comboOrder("m", trade.ComboTypeMaster, trade.OrderTypeStopLoss, trade.OrderSideBuy)
				return comboRequest("combo-oto", append([]trade.OrderRequest{m}, otoLegs(6)...)...)
			},
		},
		{
			name: "valid OTO leg stop loss limit",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-oto",
					comboOrder("m", trade.ComboTypeMaster, trade.OrderTypeStopLossLimit, trade.OrderSideBuy),
					comboOrder("l", trade.ComboTypeOTO, trade.OrderTypeStopLossLimit, trade.OrderSideSell),
				)
			},
		},
		{
			name: "valid OCO two legs",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-oco", ocoLegs(2)...)
			},
		},
		{
			name: "valid OCO six legs",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-oco", ocoLegs(6)...)
			},
		},
		{
			name: "valid OTOCO",
			build: func() trade.PlaceOrderRequest {
				m := comboOrder("m", trade.ComboTypeMaster, trade.OrderTypeMarket, trade.OrderSideBuy)
				return comboRequest("combo-otoco", append([]trade.OrderRequest{m}, otocoLegs(1)...)...)
			},
		},
		{
			name: "valid NORMAL unaffected",
			build: func() trade.PlaceOrderRequest {
				return validPlaceRequest()
			},
		},
		{
			name: "missing client_combo_order_id",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("",
					comboOrder("m", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy),
				)
			},
			wantErr:    true,
			wantSubstr: "client_combo_order_id is required",
		},
		{
			name: "NORMAL mixed with combo",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-tpsl",
					validOrder(),
					comboOrder("m", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy),
				)
			},
			wantErr:    true,
			wantSubstr: "NORMAL orders cannot be mixed",
		},
		{
			name: "option order with combo type",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-opt", validOptionForCombo(trade.ComboTypeMaster))
			},
			wantErr:    true,
			wantSubstr: "only supported for EQUITY orders",
		},
		{
			name: "futures order with combo type",
			build: func() trade.PlaceOrderRequest {
				o := comboOrder("f", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy)
				o.InstrumentType = trade.InstrumentTypeFutures
				o.SupportTradingSession = ""
				return comboRequest("combo-fut", o)
			},
			wantErr:    true,
			wantSubstr: "only supported for EQUITY orders",
		},
		{
			name: "HK order with combo type",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-hk", validHKForCombo(trade.ComboTypeMaster))
			},
			wantErr:    true,
			wantSubstr: "only supported for US orders",
		},
		{
			name: "TP/SL two masters",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-tpsl",
					comboOrder("m1", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy),
					comboOrder("m2", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy),
				)
			},
			wantErr:    true,
			wantSubstr: "at most one MASTER",
		},
		{
			name: "TP/SL two stop profits",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-tpsl",
					comboOrder("p1", trade.ComboTypeStopProfit, trade.OrderTypeLimit, trade.OrderSideSell),
					comboOrder("p2", trade.ComboTypeStopProfit, trade.OrderTypeLimit, trade.OrderSideSell),
				)
			},
			wantErr:    true,
			wantSubstr: "at most one STOP_PROFIT",
		},
		{
			name: "TP/SL two stop losses",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-tpsl",
					comboOrder("l1", trade.ComboTypeStopLoss, trade.OrderTypeStopLoss, trade.OrderSideSell),
					comboOrder("l2", trade.ComboTypeStopLoss, trade.OrderTypeStopLoss, trade.OrderSideSell),
				)
			},
			wantErr:    true,
			wantSubstr: "at most one STOP_LOSS",
		},
		{
			name: "TP/SL master stop loss type rejected",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-tpsl",
					comboOrder("m", trade.ComboTypeMaster, trade.OrderTypeStopLoss, trade.OrderSideBuy),
				)
			},
			wantErr:    true,
			wantSubstr: "not supported for MASTER orders",
		},
		{
			name: "TP/SL stop profit market type rejected",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-tpsl",
					comboOrder("p", trade.ComboTypeStopProfit, trade.OrderTypeMarket, trade.OrderSideSell),
				)
			},
			wantErr:    true,
			wantSubstr: "not supported for STOP_PROFIT orders",
		},
		{
			name: "TP/SL stop loss limit type rejected",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-tpsl",
					comboOrder("l", trade.ComboTypeStopLoss, trade.OrderTypeStopLossLimit, trade.OrderSideSell),
				)
			},
			wantErr:    true,
			wantSubstr: "not supported for STOP_LOSS orders",
		},
		{
			name: "TP/SL sell-to-close buy side rejected",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-tpsl",
					comboOrder("p", trade.ComboTypeStopProfit, trade.OrderTypeLimit, trade.OrderSideBuy),
					comboOrder("l", trade.ComboTypeStopLoss, trade.OrderTypeStopLoss, trade.OrderSideBuy),
				)
			},
			wantErr:    true,
			wantSubstr: "must use side SELL",
		},
		{
			name: "OTO missing master",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-oto", otoLegs(1)...)
			},
			wantErr:    true,
			wantSubstr: "exactly one MASTER",
		},
		{
			name: "OTO two masters",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-oto",
					comboOrder("m1", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy),
					comboOrder("m2", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy),
					comboOrder("l", trade.ComboTypeOTO, trade.OrderTypeLimit, trade.OrderSideSell),
				)
			},
			wantErr:    true,
			wantSubstr: "exactly one MASTER",
		},
		{
			name: "OTO more than six legs",
			build: func() trade.PlaceOrderRequest {
				m := comboOrder("m", trade.ComboTypeMaster, trade.OrderTypeMarket, trade.OrderSideBuy)
				return comboRequest("combo-oto", append([]trade.OrderRequest{m}, otoLegs(7)...)...)
			},
			wantErr:    true,
			wantSubstr: "between 1 and 6 OTO orders",
		},
		{
			name: "OTO master trailing type rejected",
			build: func() trade.PlaceOrderRequest {
				m := comboOrder("m", trade.ComboTypeMaster, trade.OrderTypeTrailingStopLoss, trade.OrderSideBuy)
				m.TrailingType = trade.TrailingTypeAmount
				m.TrailingStopStep = mp("1.00")
				return comboRequest("combo-oto", m,
					comboOrder("l", trade.ComboTypeOTO, trade.OrderTypeLimit, trade.OrderSideSell))
			},
			wantErr:    true,
			wantSubstr: "not supported for MASTER orders",
		},
		{
			name: "OTO leg market type allowed",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-oto",
					comboOrder("m", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy),
					comboOrder("l", trade.ComboTypeOTO, trade.OrderTypeMarket, trade.OrderSideSell))
			},
		},
		{
			name: "OTO with stop profit mixed",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-oto",
					comboOrder("m", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy),
					comboOrder("l", trade.ComboTypeOTO, trade.OrderTypeLimit, trade.OrderSideSell),
					comboOrder("p", trade.ComboTypeStopProfit, trade.OrderTypeLimit, trade.OrderSideSell))
			},
			wantErr:    true,
			wantSubstr: "cannot be combined in a OTO group",
		},
		{
			name: "OCO one leg",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-oco", ocoLegs(1)...)
			},
			wantErr:    true,
			wantSubstr: "between 2 and 6 OCO orders",
		},
		{
			name: "OCO seven legs",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-oco", ocoLegs(7)...)
			},
			wantErr:    true,
			wantSubstr: "between 2 and 6 OCO orders",
		},
		{
			name: "OCO with master mixed",
			build: func() trade.PlaceOrderRequest {
				m := comboOrder("m", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy)
				return comboRequest("combo-oco", append([]trade.OrderRequest{m}, ocoLegs(2)...)...)
			},
			wantErr:    true,
			wantSubstr: "cannot be combined in a OCO group",
		},
		{
			name: "OCO market leg rejected",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-oco",
					comboOrder("l1", trade.ComboTypeOCO, trade.OrderTypeMarket, trade.OrderSideSell),
					comboOrder("l2", trade.ComboTypeOCO, trade.OrderTypeLimit, trade.OrderSideSell))
			},
			wantErr:    true,
			wantSubstr: "not supported for OCO orders",
		},
		{
			name: "OTOCO missing master",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-otoco", otocoLegs(1)...)
			},
			wantErr:    true,
			wantSubstr: "exactly one MASTER",
		},
		{
			name: "OTOCO two masters",
			build: func() trade.PlaceOrderRequest {
				return comboRequest("combo-otoco",
					comboOrder("m1", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy),
					comboOrder("m2", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy),
					comboOrder("l", trade.ComboTypeOTOCO, trade.OrderTypeLimit, trade.OrderSideSell))
			},
			wantErr:    true,
			wantSubstr: "exactly one MASTER",
		},
		{
			name: "OTOCO seven legs",
			build: func() trade.PlaceOrderRequest {
				m := comboOrder("m", trade.ComboTypeMaster, trade.OrderTypeMarket, trade.OrderSideBuy)
				return comboRequest("combo-otoco", append([]trade.OrderRequest{m}, otocoLegs(7)...)...)
			},
			wantErr:    true,
			wantSubstr: "between 1 and 6 OTOCO orders",
		},
		{
			name: "OTOCO market leg rejected",
			build: func() trade.PlaceOrderRequest {
				m := comboOrder("m", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy)
				return comboRequest("combo-otoco", append([]trade.OrderRequest{m},
					comboOrder("l", trade.ComboTypeOTOCO, trade.OrderTypeMarket, trade.OrderSideSell))...)
			},
			wantErr:    true,
			wantSubstr: "not supported for OTOCO orders",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.build().Validate()
			if !tc.wantErr {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}
			if !errs.Is(err, errs.CodeInvalidConfig) {
				t.Fatalf("Validate() error = %v, want invalid_config", err)
			}
			if tc.wantSubstr != "" && !strings.Contains(err.Error(), tc.wantSubstr) {
				t.Fatalf("Validate() error = %q, want it to contain %q", err, tc.wantSubstr)
			}
		})
	}
}

// TestPlaceOrderSerializesComboGroup asserts that a validated TP/SL combo is
// serialized as one grouped body: a MASTER, a STOP_PROFIT, and a STOP_LOSS
// sharing a single client_combo_order_id.
func TestPlaceOrderSerializesComboGroup(t *testing.T) {
	t.Parallel()

	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if got, want := r.URL.Path, "/trading/orders/place"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		raw := string(readBody(t, r))
		for _, want := range []string{
			`"client_combo_order_id":"combo-tpsl-1"`,
			`"combo_type":"MASTER"`,
			`"combo_type":"STOP_PROFIT"`,
			`"combo_type":"STOP_LOSS"`,
			`"order_type":"STOP_LOSS"`,
			`"stop_price":"170"`,
		} {
			if !strings.Contains(raw, want) {
				t.Errorf("body missing %s: %s", want, raw)
			}
		}

		var body trade.PlaceOrderRequest
		if err := json.Unmarshal([]byte(raw), &body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if body.ClientComboOrderID != "combo-tpsl-1" {
			t.Errorf("client_combo_order_id = %q, want combo-tpsl-1", body.ClientComboOrderID)
		}
		got := []trade.ComboType{}
		for _, o := range body.NewOrders {
			got = append(got, o.ComboType)
		}
		want := []trade.ComboType{trade.ComboTypeMaster, trade.ComboTypeStopProfit, trade.ComboTypeStopLoss}
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("combo types = %v, want %v", got, want)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"client_order_id":"tpsl-master","order_id":"OID-1"}`))
	}))
	defer srv.Close()

	req := comboRequest("combo-tpsl-1",
		comboOrder("tpsl-master", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy),
		comboOrder("tpsl-profit", trade.ComboTypeStopProfit, trade.OrderTypeLimit, trade.OrderSideSell),
		comboOrder("tpsl-loss", trade.ComboTypeStopLoss, trade.OrderTypeStopLoss, trade.OrderSideSell),
	)

	c := newOrderTestClient(t, srv.URL)
	if _, err := c.PlaceOrder(context.Background(), req); err != nil {
		t.Fatalf("PlaceOrder() error = %v", err)
	}
	if hits != 1 {
		t.Fatalf("server received %d requests, want 1", hits)
	}
}

// TestComboValidationPrecedesNetwork asserts that an invalid combo group is
// rejected before the request reaches the API.
func TestComboValidationPrecedesNetwork(t *testing.T) {
	t.Parallel()

	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	bad := comboRequest("combo-oco", ocoLegs(1)...)

	c := newOrderTestClient(t, srv.URL)
	if _, err := c.PreviewOrder(context.Background(), bad); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("PreviewOrder() error = %v, want invalid_config", err)
	}
	if _, err := c.PlaceOrder(context.Background(), bad); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("PlaceOrder() error = %v, want invalid_config", err)
	}
	if hits != 0 {
		t.Fatalf("server received %d requests, want 0 (validation must precede the network)", hits)
	}
}

// TestSandboxPreviewComboOrder previews a small AAPL buy-to-open take-profit/
// stop-loss combo on the sandbox account. It is read-only: preview never places
// or mutates an order. A sandbox rejection, for example when the account lacks
// combo entitlement, is tolerated as a typed [errs.Error] rather than a
// failure.
//
// It runs only when WEBULL_TRADE_SANDBOX=1, WEBULL_TRADE_APP_KEY,
// WEBULL_TRADE_APP_SECRET, and WEBULL_TRADE_ACCOUNT_ID are set; otherwise it
// skips so the default test run stays hermetic.
func TestSandboxPreviewComboOrder(t *testing.T) {
	trading, accountID := newSandboxTradeClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := comboRequest(fmt.Sprintf("sdk-combo-%d", time.Now().UnixNano()),
		comboOrder("sdk-combo-master", trade.ComboTypeMaster, trade.OrderTypeLimit, trade.OrderSideBuy),
		comboOrder("sdk-combo-profit", trade.ComboTypeStopProfit, trade.OrderTypeLimit, trade.OrderSideSell),
		comboOrder("sdk-combo-loss", trade.ComboTypeStopLoss, trade.OrderTypeStopLoss, trade.OrderSideSell),
	)
	req.AccountID = accountID

	res, err := trading.PreviewOrder(ctx, req)
	if err != nil {
		var e *errs.Error
		if !errors.As(err, &e) {
			t.Fatalf("PreviewOrder(combo) error = %v (%T), want *errs.Error", err, err)
		}
		t.Logf("PreviewOrder(combo) typed error (tolerated): code=%s message=%s", e.Code, e.Message)
		return
	}
	t.Logf("preview combo estimated_cost=%s estimated_transaction_fee=%s",
		res.EstimatedCost, res.EstimatedTransactionFee)
}
