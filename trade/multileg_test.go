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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/trade"
)

// optionLeg builds an option leg from the single-leg fixture with the given
// side, type, strike, and expiration. Tests use it to assemble multi-leg
// strategies leg by leg.
func optionLeg(side trade.OrderSide, optType trade.OptionType, strike, expire string) trade.OrderLeg {
	l := validOptionLeg()
	l.Side = side
	l.OptionType = optType
	l.StrikePrice = strike
	l.OptionExpireDate = expire
	return l
}

// multiLegOrder returns a well-formed multi-leg option order carrying the given
// strategy and legs. The top-level fields satisfy the shared option rules; the
// legality of the strategy's structure is what the tests vary.
func multiLegOrder(strategy trade.OptionStrategy, legs ...trade.OrderLeg) trade.OrderRequest {
	return trade.OrderRequest{
		ClientOrderID:  "multi-leg-1",
		ComboType:      trade.ComboTypeNormal,
		InstrumentType: trade.InstrumentTypeOption,
		Market:         trade.MarketUS,
		Symbol:         "AAPL",
		OrderType:      trade.OrderTypeLimit,
		Side:           trade.OrderSideBuy,
		Quantity:       "1",
		EntrustType:    trade.EntrustTypeQty,
		TimeInForce:    trade.TimeInForceDay,
		LimitPrice:     "1.50",
		OptionStrategy: strategy,
		Legs:           legs,
	}
}

// verticalLegs returns a bull call spread: long the 220 call and short the 230
// call, same expiration.
func verticalLegs() []trade.OrderLeg {
	return []trade.OrderLeg{
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "220.00", "2026-12-18"),
		optionLeg(trade.OrderSideSell, trade.OptionTypeCall, "230.00", "2026-12-18"),
	}
}

// ironCondorLegs returns a four-leg iron condor: a short put spread below the
// underlying and a short call spread above it.
func ironCondorLegs() []trade.OrderLeg {
	return []trade.OrderLeg{
		optionLeg(trade.OrderSideBuy, trade.OptionTypePut, "210.00", "2026-12-18"),
		optionLeg(trade.OrderSideSell, trade.OptionTypePut, "220.00", "2026-12-18"),
		optionLeg(trade.OrderSideSell, trade.OptionTypeCall, "230.00", "2026-12-18"),
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "240.00", "2026-12-18"),
	}
}

// butterflyLegs returns a three-leg butterfly: long the 210 and 230 calls and
// short the 220 call, all the same expiration. It covers an odd leg count above
// the two-leg minimum.
func butterflyLegs() []trade.OrderLeg {
	return []trade.OrderLeg{
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "210.00", "2026-12-18"),
		optionLeg(trade.OrderSideSell, trade.OptionTypeCall, "220.00", "2026-12-18"),
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "230.00", "2026-12-18"),
	}
}

// manyLegs returns a six-leg set with a distinct side, type, and strike per
// leg. It covers leg-count growth well past the minimum.
func manyLegs() []trade.OrderLeg {
	return []trade.OrderLeg{
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "200.00", "2026-12-18"),
		optionLeg(trade.OrderSideSell, trade.OptionTypeCall, "210.00", "2026-12-18"),
		optionLeg(trade.OrderSideBuy, trade.OptionTypePut, "190.00", "2026-12-18"),
		optionLeg(trade.OrderSideSell, trade.OptionTypePut, "180.00", "2026-12-18"),
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "250.00", "2027-01-15"),
		optionLeg(trade.OrderSideSell, trade.OptionTypeCall, "260.00", "2027-01-15"),
	}
}

// TestMultiLegOptionValidate covers the structural rules for multi-leg option
// orders. Cases with wantErr set must fail with [errs.CodeInvalidConfig] before
// any network call; wantSubstr, when set, pins the diagnostic to the offending
// leg.
func TestMultiLegOptionValidate(t *testing.T) {
	t.Parallel()

	degenerate := []trade.OrderLeg{
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "220.00", "2026-12-18"),
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "220.00", "2027-01-15"),
	}
	duplicate := []trade.OrderLeg{
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "220.00", "2026-12-18"),
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "220.00", "2026-12-18"),
	}
	nonOptionLeg := verticalLegs()
	nonOptionLeg[1].InstrumentType = trade.InstrumentTypeEquity

	marketOrder := multiLegOrder(trade.OptionStrategyVertical, verticalLegs()...)
	marketOrder.OrderType = trade.OrderTypeMarket
	marketOrder.LimitPrice = ""

	stopLossOrder := multiLegOrder(trade.OptionStrategyVertical, verticalLegs()...)
	stopLossOrder.OrderType = trade.OrderTypeStopLoss
	stopLossOrder.LimitPrice = ""
	stopLossOrder.StopPrice = "1.00"

	duplicateNonAdjacent := []trade.OrderLeg{
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "210.00", "2026-12-18"),
		optionLeg(trade.OrderSideSell, trade.OptionTypeCall, "220.00", "2026-12-18"),
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "210.00", "2026-12-18"),
		optionLeg(trade.OrderSideSell, trade.OptionTypeCall, "230.00", "2026-12-18"),
	}

	sameStrikeDifferentSides := []trade.OrderLeg{
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "220.00", "2026-12-18"),
		optionLeg(trade.OrderSideSell, trade.OptionTypeCall, "220.00", "2026-12-18"),
	}

	degenerateSymbols := []trade.OrderLeg{
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "220.00", "2026-12-18"),
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "220.00", "2026-12-18"),
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "220.00", "2026-12-18"),
	}
	degenerateSymbols[0].Symbol = "AAPL"
	degenerateSymbols[1].Symbol = "MSFT"
	degenerateSymbols[2].Symbol = "TSLA"

	duplicatePrecision := []trade.OrderLeg{
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "220.0", "2026-12-18"),
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "220.00", "2026-12-18"),
	}

	degeneratePrecision := []trade.OrderLeg{
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "220.0", "2026-12-18"),
		optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "220.00", "2027-01-15"),
	}

	cases := []struct {
		name       string
		order      trade.OrderRequest
		wantErr    bool
		wantSubstr string
	}{
		{"valid 2-leg vertical", multiLegOrder(trade.OptionStrategyVertical, verticalLegs()...), false, ""},
		{"valid 2-leg straddle", multiLegOrder(trade.OptionStrategyStraddle,
			optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "220.00", "2026-12-18"),
			optionLeg(trade.OrderSideBuy, trade.OptionTypePut, "220.00", "2026-12-18"),
		), false, ""},
		{"valid 4-leg iron condor", multiLegOrder(trade.OptionStrategyIronCondor, ironCondorLegs()...), false, ""},
		{"valid 3-leg butterfly", multiLegOrder(trade.OptionStrategyButterfly, butterflyLegs()...), false, ""},
		{"valid 6-leg set", multiLegOrder(trade.OptionStrategyRatio, manyLegs()...), false, ""},
		{"valid same strike different sides", multiLegOrder(trade.OptionStrategyRatio,
			sameStrikeDifferentSides...), false, ""},
		{"valid single leg unchanged", validOptionOrder(), false, ""},
		{"unknown multi-leg strategy rejected", multiLegOrder(trade.OptionStrategy("NOT_A_STRATEGY"),
			verticalLegs()...), true, "not a supported option strategy"},
		{"empty multi-leg strategy rejected", multiLegOrder(trade.OptionStrategy(""),
			verticalLegs()...), true, "not a supported option strategy"},
		{"one-leg multi-leg strategy rejected", multiLegOrder(trade.OptionStrategyVertical,
			optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "220.00", "2026-12-18"),
		), true, "at least two legs"},
		{"no-leg multi-leg strategy rejected", multiLegOrder(trade.OptionStrategyVertical), true, "at least two legs"},
		{"duplicate legs rejected", multiLegOrder(trade.OptionStrategyVertical, duplicate...), true, "duplicate leg"},
		{"duplicate legs differing only in strike precision rejected", multiLegOrder(trade.OptionStrategyButterfly,
			duplicatePrecision...), true, "duplicate leg"},
		{"duplicate non-adjacent legs rejected", multiLegOrder(trade.OptionStrategyButterfly,
			duplicateNonAdjacent...), true, "legs[2]: duplicate leg"},
		{"degenerate same side/type/strike across symbols rejected", multiLegOrder(trade.OptionStrategyRatio,
			degenerateSymbols...), true, "same side"},
		{"degenerate same side/type/strike across strike precision rejected", multiLegOrder(trade.OptionStrategyCalendar,
			degeneratePrecision...), true, "same side"},
		{"non-option leg rejected", multiLegOrder(trade.OptionStrategyVertical, nonOptionLeg...), true, "legs[1]"},
		{"short-side leg rejected", multiLegOrder(trade.OptionStrategyVertical,
			optionLeg(trade.OrderSideBuy, trade.OptionTypeCall, "220.00", "2026-12-18"),
			optionLeg(trade.OrderSideShort, trade.OptionTypeCall, "230.00", "2026-12-18"),
		), true, "legs[1]"},
		{"degenerate same side/type/strike rejected", multiLegOrder(trade.OptionStrategyCalendar, degenerate...), true, "same side"},
		{"multi-leg market order rejected", marketOrder, true, "order_type"},
		{"multi-leg stop loss order rejected", stopLossOrder, true, "order_type"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.order.Validate()
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
				t.Fatalf("Validate() error = %v, want substring %q", err, tc.wantSubstr)
			}
		})
	}
}

// TestMultiLegOptionSkipsNotionalGuardrail confirms a multi-leg order is not
// rejected by the notional cap even when the top-level quantity and limit price
// would exceed it, because each leg is priced separately.
func TestMultiLegOptionSkipsNotionalGuardrail(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"estimated_cost":"1.00","estimated_transaction_fee":"0.10"}`))
	}))
	defer srv.Close()

	c := newTradeClient(t, srv.URL, trade.WithMaxOrderNotional("1"))

	order := multiLegOrder(trade.OptionStrategyVertical, verticalLegs()...)
	order.Quantity = "1000"
	order.LimitPrice = "1000.00"

	req := trade.PlaceOrderRequest{AccountID: "ACC1", NewOrders: []trade.OrderRequest{order}}
	if _, err := c.PreviewOrder(context.Background(), req); err != nil {
		t.Fatalf("PreviewOrder() error = %v, want nil (multi-leg notional must be skipped)", err)
	}
}

// TestMultiLegOptionComboRejected confirms a multi-leg option order cannot take
// part in any equity combo group. The client combo identifier is set so the
// failure is the equity-only instrument rule rather than a missing group id.
func TestMultiLegOptionComboRejected(t *testing.T) {
	t.Parallel()

	comboTypes := []trade.ComboType{
		trade.ComboTypeMaster,
		trade.ComboTypeStopProfit,
		trade.ComboTypeStopLoss,
		trade.ComboTypeOTO,
		trade.ComboTypeOCO,
		trade.ComboTypeOTOCO,
	}

	for _, ct := range comboTypes {
		t.Run(string(ct), func(t *testing.T) {
			t.Parallel()

			order := multiLegOrder(trade.OptionStrategyVertical, verticalLegs()...)
			order.ComboType = ct

			req := trade.PlaceOrderRequest{
				AccountID:          "ACC1",
				ClientComboOrderID: "combo-multi-leg",
				NewOrders:          []trade.OrderRequest{order},
			}
			err := req.Validate()
			if !errs.Is(err, errs.CodeInvalidConfig) {
				t.Fatalf("Validate() error = %v, want invalid_config", err)
			}
			if !strings.Contains(err.Error(), "only supported for EQUITY orders") {
				t.Fatalf("Validate() error = %q, want the equity-only combo message", err)
			}
		})
	}
}

// TestMultiLegOptionNormalPlaceRequestValid confirms a NORMAL multi-leg order is
// accepted by the batch validator, so only combo participation is rejected.
func TestMultiLegOptionNormalPlaceRequestValid(t *testing.T) {
	t.Parallel()

	order := multiLegOrder(trade.OptionStrategyIronCondor, ironCondorLegs()...)
	req := trade.PlaceOrderRequest{AccountID: "ACC1", NewOrders: []trade.OrderRequest{order}}
	if err := req.Validate(); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
}

// TestMultiLegOptionSkipsNotionalGuardrailFourLegs confirms the notional skip
// covers any leg count above one, not just the two-leg boundary.
func TestMultiLegOptionSkipsNotionalGuardrailFourLegs(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"estimated_cost":"1.00","estimated_transaction_fee":"0.10"}`))
	}))
	defer srv.Close()

	c := newTradeClient(t, srv.URL, trade.WithMaxOrderNotional("1"))

	order := multiLegOrder(trade.OptionStrategyIronCondor, ironCondorLegs()...)
	order.Quantity = "1000"
	order.LimitPrice = "1000.00"

	req := trade.PlaceOrderRequest{AccountID: "ACC1", NewOrders: []trade.OrderRequest{order}}
	if _, err := c.PreviewOrder(context.Background(), req); err != nil {
		t.Fatalf("PreviewOrder() error = %v, want nil (four-leg notional must be skipped)", err)
	}
}

// TestSingleLegOptionNotionalGuardrail confirms the notional guardrail still
// applies to a single-leg option order, which is the boundary the multi-leg
// skip must not cross.
func TestSingleLegOptionNotionalGuardrail(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"estimated_cost":"1.00","estimated_transaction_fee":"0.10"}`))
	}))
	defer srv.Close()

	c := newTradeClient(t, srv.URL, trade.WithMaxOrderNotional("100"))

	order := validOptionOrder()
	order.Quantity = "100"
	order.LimitPrice = "10.00"

	req := trade.PlaceOrderRequest{AccountID: "ACC1", NewOrders: []trade.OrderRequest{order}}
	if _, err := c.PreviewOrder(context.Background(), req); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("PreviewOrder() error = %v, want invalid_config (single-leg notional is enforced)", err)
	}
}
