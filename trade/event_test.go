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
	"errors"
	"strings"
	"testing"

	"github.com/shing1211/webullapi4go/internal/errs"
	"github.com/shing1211/webullapi4go/trade"
)

// validEventOrder returns a well-formed US event contract limit order. Tests
// mutate a copy to isolate a single event rule.
func validEventOrder() trade.OrderRequest {
	return trade.OrderRequest{
		ClientOrderID:  "event-order-1",
		ComboType:      trade.ComboTypeNormal,
		InstrumentType: trade.InstrumentTypeEvent,
		Market:         trade.MarketUS,
		Symbol:         "AAPL-EVENT-20261218",
		OrderType:      trade.OrderTypeLimit,
		Side:           trade.OrderSideBuy,
		Quantity:       "10",
		EntrustType:    trade.EntrustTypeQty,
		TimeInForce:    trade.TimeInForceDay,
		LimitPrice:     "0.55",
	}
}

// TestEventOrderValid asserts that a well-formed event contract order passes
// validation.
func TestEventOrderValid(t *testing.T) {
	t.Parallel()
	o := validEventOrder()
	if err := o.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestEventOrderRejectsMarket asserts that a MARKET order type is rejected for
// event contracts.
func TestEventOrderRejectsMarket(t *testing.T) {
	t.Parallel()
	o := validEventOrder()
	o.OrderType = trade.OrderTypeMarket
	o.LimitPrice = ""
	err := o.Validate()
	if err == nil {
		t.Fatal("expected error for MARKET order type, got nil")
	}
	var e *errs.Error
	if !errors.As(err, &e) || e.Code != errs.CodeInvalidConfig {
		t.Fatalf("expected CodeInvalidConfig, got %v", err)
	}
	if !strings.Contains(e.Message, "order_type MARKET is not supported for event contract orders") {
		t.Fatalf("unexpected message: %s", e.Message)
	}
}

// TestEventOrderRejectsGTC asserts that GTC time-in-force is rejected for event
// contracts.
func TestEventOrderRejectsGTC(t *testing.T) {
	t.Parallel()
	o := validEventOrder()
	o.TimeInForce = trade.TimeInForceGTC
	err := o.Validate()
	if err == nil {
		t.Fatal("expected error for GTC time-in-force, got nil")
	}
	var e *errs.Error
	if !errors.As(err, &e) || e.Code != errs.CodeInvalidConfig {
		t.Fatalf("expected CodeInvalidConfig, got %v", err)
	}
	if !strings.Contains(e.Message, "time_in_force GTC is not supported for event contract orders") {
		t.Fatalf("unexpected message: %s", e.Message)
	}
}

// TestEventOrderRejectsOver50k asserts that a quantity exceeding 50,000 is
// rejected for event contracts.
func TestEventOrderRejectsOver50k(t *testing.T) {
	t.Parallel()
	o := validEventOrder()
	o.Quantity = "50001"
	err := o.Validate()
	if err == nil {
		t.Fatal("expected error for quantity >50000, got nil")
	}
	var e *errs.Error
	if !errors.As(err, &e) || e.Code != errs.CodeInvalidConfig {
		t.Fatalf("expected CodeInvalidConfig, got %v", err)
	}
	if !strings.Contains(e.Message, "exceeds the maximum 50000") {
		t.Fatalf("unexpected message: %s", e.Message)
	}
}

// TestEventOrderRejectsAMOUNT asserts that AMOUNT entrust type is rejected for
// event contracts.
func TestEventOrderRejectsAMOUNT(t *testing.T) {
	t.Parallel()
	o := validEventOrder()
	o.EntrustType = trade.EntrustTypeAmount
	o.Quantity = ""
	o.TotalCashAmount = "55.00"
	err := o.Validate()
	if err == nil {
		t.Fatal("expected error for AMOUNT entrust type, got nil")
	}
	var e *errs.Error
	if !errors.As(err, &e) || e.Code != errs.CodeInvalidConfig {
		t.Fatalf("expected CodeInvalidConfig, got %v", err)
	}
	if !strings.Contains(e.Message, "entrust_type \"AMOUNT\" must be QTY for event contract orders") {
		t.Fatalf("unexpected message: %s", e.Message)
	}
}

// TestEventOrderRejectsOptionFields asserts that option_strategy is rejected for
// event contract orders.
func TestEventOrderRejectsOptionFields(t *testing.T) {
	t.Parallel()
	o := validEventOrder()
	o.OptionStrategy = trade.OptionStrategySingle
	err := o.Validate()
	if err == nil {
		t.Fatal("expected error for option_strategy, got nil")
	}
	var e *errs.Error
	if !errors.As(err, &e) || e.Code != errs.CodeInvalidConfig {
		t.Fatalf("expected CodeInvalidConfig, got %v", err)
	}
	if !strings.Contains(e.Message, "option_strategy is only valid for OPTION orders") {
		t.Fatalf("unexpected message: %s", e.Message)
	}
}

// TestEventOrderRejectsLegs asserts that legs are rejected for event contract
// orders.
func TestEventOrderRejectsLegs(t *testing.T) {
	t.Parallel()
	o := validEventOrder()
	o.Legs = []trade.OrderLeg{{
		InstrumentType: trade.InstrumentTypeOption,
		Market:         trade.MarketUS,
		Symbol:         "AAPL",
		Side:           trade.OrderSideBuy,
		StrikePrice:    "220.00",
		OptionType:     trade.OptionTypeCall,
		Quantity:       "1",
	}}
	err := o.Validate()
	if err == nil {
		t.Fatal("expected error for legs, got nil")
	}
	var e *errs.Error
	if !errors.As(err, &e) || e.Code != errs.CodeInvalidConfig {
		t.Fatalf("expected CodeInvalidConfig, got %v", err)
	}
	if !strings.Contains(e.Message, "legs is only valid for OPTION orders") {
		t.Fatalf("unexpected message: %s", e.Message)
	}
}

// TestEventOrderRejectsTradingSession asserts that support_trading_session is
// rejected for event contract orders.
func TestEventOrderRejectsTradingSession(t *testing.T) {
	t.Parallel()
	o := validEventOrder()
	o.SupportTradingSession = trade.TradingSessionAll
	err := o.Validate()
	if err == nil {
		t.Fatal("expected error for support_trading_session, got nil")
	}
	var e *errs.Error
	if !errors.As(err, &e) || e.Code != errs.CodeInvalidConfig {
		t.Fatalf("expected CodeInvalidConfig, got %v", err)
	}
	if !strings.Contains(e.Message, "support_trading_session is not valid for event contract orders") {
		t.Fatalf("unexpected message: %s", e.Message)
	}
}

// TestEventOrderRejectsPartyIDs asserts that no_party_ids is rejected for event
// contract orders.
func TestEventOrderRejectsPartyIDs(t *testing.T) {
	t.Parallel()
	o := validEventOrder()
	o.NoPartyIDs = []trade.PartyID{{PartyID: "A", PartyIDSource: "D", PartyRole: "3"}}
	err := o.Validate()
	if err == nil {
		t.Fatal("expected error for no_party_ids, got nil")
	}
	var e *errs.Error
	if !errors.As(err, &e) || e.Code != errs.CodeInvalidConfig {
		t.Fatalf("expected CodeInvalidConfig, got %v", err)
	}
	if !strings.Contains(e.Message, "no_party_ids is not valid for event contract orders") {
		t.Fatalf("unexpected message: %s", e.Message)
	}
}
