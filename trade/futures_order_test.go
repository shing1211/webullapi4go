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
	"strings"
	"testing"

	"github.com/shing1211/webullapi4go/internal/errs"
	"github.com/shing1211/webullapi4go/trade"
)

// validFuturesOrder returns a well-formed US futures limit order. Tests mutate a
// copy to isolate a single futures rule.
func validFuturesOrder() trade.OrderRequest {
	return trade.OrderRequest{
		ClientOrderID:  "futures-order-1",
		ComboType:      trade.ComboTypeNormal,
		InstrumentType: trade.InstrumentTypeFutures,
		Market:         trade.MarketUS,
		Symbol:         "ESZ5",
		OrderType:      trade.OrderTypeLimit,
		Side:           trade.OrderSideBuy,
		Quantity:       "1",
		EntrustType:    trade.EntrustTypeQty,
		TimeInForce:    trade.TimeInForceDay,
		LimitPrice:     "4500.00",
	}
}

// TestFuturesRequestValidate exercises every futures-specific rule. Cases with
// wantNil are expected to pass; every other case must fail with
// [errs.CodeInvalidConfig] before any network call.
func TestFuturesRequestValidate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		mutate  func(*trade.OrderRequest)
		wantNil bool
	}{
		{"valid US futures limit", func(o *trade.OrderRequest) {}, true},
		{"valid US futures market", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeMarket
			o.LimitPrice = ""
		}, true},
		{"valid US futures GTC", func(o *trade.OrderRequest) {
			o.TimeInForce = trade.TimeInForceGTC
		}, true},
		{"valid HK futures limit", func(o *trade.OrderRequest) {
			o.Market = trade.MarketHK
			o.Symbol = "HSIQ6"
			o.LimitPrice = "25000.00"
		}, true},
		{"valid US futures stop loss", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeStopLoss
			o.LimitPrice = ""
			o.StopPrice = "4490.00"
		}, true},
		{"valid US futures stop loss limit", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeStopLossLimit
			o.StopPrice = "4490.00"
			o.LimitPrice = "4480.00"
		}, true},
		{"valid HK futures market", func(o *trade.OrderRequest) {
			o.Market = trade.MarketHK
			o.Symbol = "HSIQ6"
			o.OrderType = trade.OrderTypeMarket
			o.LimitPrice = ""
		}, true},
		{"amount entrust rejected", func(o *trade.OrderRequest) {
			o.EntrustType = trade.EntrustTypeAmount
			o.Quantity = ""
			o.TotalCashAmount = "1000.00"
		}, false},
		{"amount entrust with quantity rejected", func(o *trade.OrderRequest) {
			o.EntrustType = trade.EntrustTypeAmount
			o.TotalCashAmount = "1000.00"
		}, false},
		{"fractional quantity rejected", func(o *trade.OrderRequest) {
			o.Quantity = "1.5"
		}, false},
		{"integral decimal quantity rejected", func(o *trade.OrderRequest) {
			o.Quantity = "1.0"
		}, false},
		{"signed quantity rejected", func(o *trade.OrderRequest) {
			o.Quantity = "+1"
		}, false},
		{"negative quantity rejected", func(o *trade.OrderRequest) {
			o.Quantity = "-1"
		}, false},
		{"zero quantity rejected", func(o *trade.OrderRequest) {
			o.Quantity = "0"
		}, false},
		{"missing quantity rejected", func(o *trade.OrderRequest) {
			o.Quantity = ""
		}, false},
		{"blank quantity rejected", func(o *trade.OrderRequest) {
			o.Quantity = "   "
		}, false},
		{"invalid time_in_force rejected", func(o *trade.OrderRequest) {
			o.TimeInForce = trade.TimeInForce("FOK")
		}, false},
		{"option strategy rejected", func(o *trade.OrderRequest) {
			o.OptionStrategy = trade.OptionStrategySingle
		}, false},
		{"legs rejected", func(o *trade.OrderRequest) {
			o.Legs = []trade.OrderLeg{validOptionLeg()}
		}, false},
		{"no_party_ids rejected", func(o *trade.OrderRequest) {
			o.NoPartyIDs = []trade.PartyID{{PartyID: "ABC123.2568", PartyIDSource: "D", PartyRole: "3"}}
		}, false},
		{"HK futures no_party_ids rejected", func(o *trade.OrderRequest) {
			o.Market = trade.MarketHK
			o.Symbol = "HSIQ6"
			o.LimitPrice = "25000.00"
			o.NoPartyIDs = []trade.PartyID{{PartyID: "ABC123.2568", PartyIDSource: "D", PartyRole: "3"}}
		}, false},
		{"support_trading_session rejected", func(o *trade.OrderRequest) {
			o.SupportTradingSession = trade.TradingSessionCore
		}, false},
		{"GTD rejected", func(o *trade.OrderRequest) {
			o.TimeInForce = trade.TimeInForceGTD
			o.ExpireDate = "2026-12-18"
		}, false},
		{"unsupported US order type rejected", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeMarketOnOpen
			o.LimitPrice = ""
		}, false},
		{"unsupported HK order type rejected", func(o *trade.OrderRequest) {
			o.Market = trade.MarketHK
			o.Symbol = "HSIQ6"
			o.OrderType = trade.OrderTypeAtAuction
		}, false},
		{"CN futures market rejected", func(o *trade.OrderRequest) {
			o.Market = trade.MarketCN
			o.Symbol = "IF2601"
		}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			o := validFuturesOrder()
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

// TestFuturesOrderTypeMessages checks that the futures order-type matrix reports
// the unsupported type and lists the market's supported set.
func TestFuturesOrderTypeMessages(t *testing.T) {
	t.Parallel()

	o := validFuturesOrder()
	o.OrderType = trade.OrderTypeMarketOnOpen
	err := o.Validate()
	if !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("Validate() error = %v, want invalid_config", err)
	}
	if !strings.Contains(err.Error(), "is not supported for US futures orders") {
		t.Fatalf("Validate() error = %q, want a futures matrix message", err)
	}
	if !strings.Contains(err.Error(), string(trade.OrderTypeLimit)) {
		t.Fatalf("Validate() error = %q, want the supported types listed", err)
	}
}

// TestFuturesUnsupportedMarketMessage pins the diagnostic for a market with no
// futures entries in the matrix. It must not render the empty "supported types:"
// list that an unknown market would otherwise produce.
func TestFuturesUnsupportedMarketMessage(t *testing.T) {
	t.Parallel()

	o := validFuturesOrder()
	o.Market = trade.MarketCN
	o.Symbol = "IF2601"

	err := o.Validate()
	if !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("Validate() error = %v, want invalid_config", err)
	}
	if !strings.Contains(err.Error(), "futures are not supported for this market") {
		t.Fatalf("Validate() error = %q, want a distinct unsupported-market message", err)
	}
	if strings.Contains(err.Error(), "supported types:") {
		t.Fatalf("Validate() error = %q, must not list an empty supported types set", err)
	}
}

// TestFuturesRequestRejectionMessages pins the diagnostic for every futures
// rejection path so a rule cannot silently change which constraint it reports.
func TestFuturesRequestRejectionMessages(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		mutate     func(*trade.OrderRequest)
		wantSubstr string
	}{
		{"amount entrust", func(o *trade.OrderRequest) {
			o.EntrustType = trade.EntrustTypeAmount
			o.TotalCashAmount = "1000.00"
		}, "must be QTY for futures orders"},
		{"fractional quantity", func(o *trade.OrderRequest) {
			o.Quantity = "1.5"
		}, "must be a positive integer for futures orders"},
		{"option strategy", func(o *trade.OrderRequest) {
			o.OptionStrategy = trade.OptionStrategySingle
		}, "option_strategy is only valid for OPTION orders"},
		{"legs", func(o *trade.OrderRequest) {
			o.Legs = []trade.OrderLeg{validOptionLeg()}
		}, "legs is only valid for OPTION orders"},
		{"no_party_ids", func(o *trade.OrderRequest) {
			o.NoPartyIDs = []trade.PartyID{{PartyID: "ABC123.2568", PartyIDSource: "D", PartyRole: "3"}}
		}, "no_party_ids is not valid for futures orders"},
		{"support_trading_session", func(o *trade.OrderRequest) {
			o.SupportTradingSession = trade.TradingSessionCore
		}, "support_trading_session is not valid for futures orders"},
		{"GTD", func(o *trade.OrderRequest) {
			o.TimeInForce = trade.TimeInForceGTD
			o.ExpireDate = "2026-12-18"
		}, "time_in_force GTD is not supported for futures orders"},
		{"unsupported US order type", func(o *trade.OrderRequest) {
			o.OrderType = trade.OrderTypeMarketOnOpen
			o.LimitPrice = ""
		}, "is not supported for US futures orders"},
		{"unsupported CN market", func(o *trade.OrderRequest) {
			o.Market = trade.MarketCN
			o.Symbol = "IF2601"
		}, "futures are not supported for this market"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			o := validFuturesOrder()
			tc.mutate(&o)
			err := o.Validate()
			if !errs.Is(err, errs.CodeInvalidConfig) {
				t.Fatalf("Validate() error = %v, want invalid_config", err)
			}
			if !strings.Contains(err.Error(), tc.wantSubstr) {
				t.Fatalf("Validate() error = %q, want substring %q", err, tc.wantSubstr)
			}
		})
	}
}

// TestFuturesOrderRejectedFromCombo is a regression test that a futures order
// remains rejected when it carries a combo role. Combo orders are US-equity
// only.
func TestFuturesOrderRejectedFromCombo(t *testing.T) {
	t.Parallel()

	o := validFuturesOrder()
	o.ComboType = trade.ComboTypeMaster
	req := trade.PlaceOrderRequest{
		AccountID:          "ACC1",
		ClientComboOrderID: "combo-fut",
		NewOrders:          []trade.OrderRequest{o},
	}

	err := req.Validate()
	if !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("Validate() error = %v, want invalid_config", err)
	}
	if !strings.Contains(err.Error(), "only supported for EQUITY orders") {
		t.Fatalf("Validate() error = %q, want a combo equity-only message", err)
	}
}
