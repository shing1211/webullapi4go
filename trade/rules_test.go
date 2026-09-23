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
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/trade"
)

// allEquityOrderTypes is every order type the API defines. The market matrix
// tests iterate all of them against every market.
var allEquityOrderTypes = []trade.OrderType{
	trade.OrderTypeLimit,
	trade.OrderTypeMarket,
	trade.OrderTypeStopLoss,
	trade.OrderTypeStopLossLimit,
	trade.OrderTypeEnhancedLimit,
	trade.OrderTypeAtAuction,
	trade.OrderTypeAtAuctionLimit,
	trade.OrderTypeMarketOnOpen,
	trade.OrderTypeMarketOnClose,
	trade.OrderTypeTrailingStopLoss,
	trade.OrderTypeTrailingStopLossLimit,
	trade.OrderTypeTouchMkt,
	trade.OrderTypeTouchLmt,
}

// marketAllowedOrderTypes mirrors the documented market matrix for equity
// orders. It is intentionally spelled out rather than derived from production
// code so the test fails if the matrix changes unexpectedly.
var marketAllowedOrderTypes = map[trade.Market]map[trade.OrderType]bool{
	trade.MarketUS: {
		trade.OrderTypeLimit:                 true,
		trade.OrderTypeMarket:                true,
		trade.OrderTypeStopLoss:              true,
		trade.OrderTypeStopLossLimit:         true,
		trade.OrderTypeMarketOnOpen:          true,
		trade.OrderTypeMarketOnClose:         true,
		trade.OrderTypeTouchMkt:              true,
		trade.OrderTypeTouchLmt:              true,
		trade.OrderTypeTrailingStopLoss:      true,
		trade.OrderTypeTrailingStopLossLimit: true,
	},
	trade.MarketHK: {
		trade.OrderTypeEnhancedLimit:         true,
		trade.OrderTypeAtAuction:             true,
		trade.OrderTypeAtAuctionLimit:        true,
		trade.OrderTypeStopLoss:              true,
		trade.OrderTypeStopLossLimit:         true,
		trade.OrderTypeTouchMkt:              true,
		trade.OrderTypeTouchLmt:              true,
		trade.OrderTypeTrailingStopLoss:      true,
		trade.OrderTypeTrailingStopLossLimit: true,
	},
	trade.MarketCN: {
		trade.OrderTypeLimit: true,
	},
}

// symbolFor returns a plausible symbol for the market.
func symbolFor(market trade.Market) string {
	switch market {
	case trade.MarketHK:
		return "00700"
	case trade.MarketCN:
		return "600519"
	default:
		return "AAPL"
	}
}

// orderFor returns a valid equity order for the given market and order type,
// populating the market-required fields (HK BCAN, US session) and the prices
// the order type requires. It is used to isolate the market rule under test:
// for every market/type pair the market accepts, Validate reports no error.
func orderFor(market trade.Market, ot trade.OrderType) trade.OrderRequest {
	o := trade.OrderRequest{
		ClientOrderID:  "order-1",
		ComboType:      trade.ComboTypeNormal,
		InstrumentType: trade.InstrumentTypeEquity,
		Market:         market,
		Symbol:         symbolFor(market),
		OrderType:      ot,
		Side:           trade.OrderSideBuy,
		Quantity:       "1",
		EntrustType:    trade.EntrustTypeQty,
		TimeInForce:    trade.TimeInForceDay,
	}
	switch market {
	case trade.MarketUS:
		o.SupportTradingSession = trade.TradingSessionCore
	case trade.MarketHK:
		o.NoPartyIDs = []trade.PartyID{{PartyID: "ABC123.2568", PartyIDSource: "D", PartyRole: "3"}}
	}
	switch ot {
	case trade.OrderTypeLimit, trade.OrderTypeEnhancedLimit, trade.OrderTypeAtAuctionLimit:
		o.LimitPrice = "10.00"
	case trade.OrderTypeStopLoss, trade.OrderTypeTouchMkt:
		o.StopPrice = "9.00"
	case trade.OrderTypeStopLossLimit, trade.OrderTypeTouchLmt:
		o.StopPrice = "9.00"
		o.LimitPrice = "10.00"
	case trade.OrderTypeTrailingStopLoss:
		o.TrailingType = trade.TrailingTypeAmount
		o.TrailingStopStep = "1.00"
	case trade.OrderTypeTrailingStopLossLimit:
		o.TrailingType = trade.TrailingTypeAmount
		o.TrailingStopStep = "1.00"
		o.TrailingLimitPriceOffset = "1.00"
	}
	return o
}

func TestEquityOrderTypeMarketMatrix(t *testing.T) {
	t.Parallel()

	markets := []trade.Market{trade.MarketUS, trade.MarketHK, trade.MarketCN}
	for _, market := range markets {
		for _, ot := range allEquityOrderTypes {
			t.Run(string(market)+"/"+string(ot), func(t *testing.T) {
				t.Parallel()
				o := orderFor(market, ot)
				err := o.Validate()
				if marketAllowedOrderTypes[market][ot] {
					if err != nil {
						t.Fatalf("Validate() error = %v, want nil for allowed %s/%s", err, market, ot)
					}
					return
				}
				if !errs.Is(err, errs.CodeInvalidConfig) {
					t.Fatalf("Validate() error = %v, want invalid_config for disallowed %s/%s", err, market, ot)
				}
				if !strings.Contains(err.Error(), "is not supported for") {
					t.Fatalf("Validate() error = %q, want a market-matrix message", err)
				}
			})
		}
	}
}

func TestCNOnlyAcceptsLimit(t *testing.T) {
	t.Parallel()

	if err := orderFor(trade.MarketCN, trade.OrderTypeLimit).Validate(); err != nil {
		t.Fatalf("CN LIMIT Validate() error = %v, want nil", err)
	}
	for _, ot := range allEquityOrderTypes {
		if ot == trade.OrderTypeLimit {
			continue
		}
		err := orderFor(trade.MarketCN, ot).Validate()
		if !errs.Is(err, errs.CodeInvalidConfig) {
			t.Fatalf("CN %s Validate() error = %v, want invalid_config", ot, err)
		}
	}
}

func TestHKBCANRules(t *testing.T) {
	t.Parallel()

	validParty := trade.PartyID{PartyID: "ABC123.2568", PartyIDSource: "D", PartyRole: "3"}

	cases := []struct {
		name    string
		order   func() trade.OrderRequest
		wantErr bool
	}{
		{"HK equity with BCAN", func() trade.OrderRequest {
			return orderFor(trade.MarketHK, trade.OrderTypeEnhancedLimit)
		}, false},
		{"HK equity without BCAN", func() trade.OrderRequest {
			o := orderFor(trade.MarketHK, trade.OrderTypeEnhancedLimit)
			o.NoPartyIDs = nil
			return o
		}, true},
		{"HK equity with empty BCAN list", func() trade.OrderRequest {
			o := orderFor(trade.MarketHK, trade.OrderTypeEnhancedLimit)
			o.NoPartyIDs = []trade.PartyID{}
			return o
		}, true},
		{"HK equity missing party_id", func() trade.OrderRequest {
			o := orderFor(trade.MarketHK, trade.OrderTypeEnhancedLimit)
			o.NoPartyIDs = []trade.PartyID{{PartyIDSource: "D", PartyRole: "3"}}
			return o
		}, true},
		{"HK equity blank party_id", func() trade.OrderRequest {
			o := orderFor(trade.MarketHK, trade.OrderTypeEnhancedLimit)
			o.NoPartyIDs = []trade.PartyID{{PartyID: "  ", PartyIDSource: "D", PartyRole: "3"}}
			return o
		}, true},
		{"HK equity bad party_id_source", func() trade.OrderRequest {
			o := orderFor(trade.MarketHK, trade.OrderTypeEnhancedLimit)
			o.NoPartyIDs = []trade.PartyID{{PartyID: "ABC123.2568", PartyIDSource: "X", PartyRole: "3"}}
			return o
		}, true},
		{"HK equity bad party_role", func() trade.OrderRequest {
			o := orderFor(trade.MarketHK, trade.OrderTypeEnhancedLimit)
			o.NoPartyIDs = []trade.PartyID{{PartyID: "ABC123.2568", PartyIDSource: "D", PartyRole: "1"}}
			return o
		}, true},
		{"HK equity multiple valid parties", func() trade.OrderRequest {
			o := orderFor(trade.MarketHK, trade.OrderTypeEnhancedLimit)
			o.NoPartyIDs = []trade.PartyID{validParty, {PartyID: "DEF456.9999", PartyIDSource: "D", PartyRole: "3"}}
			return o
		}, false},
		{"US with no_party_ids rejected", func() trade.OrderRequest {
			o := orderFor(trade.MarketUS, trade.OrderTypeLimit)
			o.NoPartyIDs = []trade.PartyID{validParty}
			return o
		}, true},
		{"CN with no_party_ids rejected", func() trade.OrderRequest {
			o := orderFor(trade.MarketCN, trade.OrderTypeLimit)
			o.NoPartyIDs = []trade.PartyID{validParty}
			return o
		}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.order().Validate()
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

func TestSupportTradingSessionRules(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		order   func() trade.OrderRequest
		wantErr bool
	}{
		{"US empty", func() trade.OrderRequest {
			o := orderFor(trade.MarketUS, trade.OrderTypeLimit)
			o.SupportTradingSession = ""
			return o
		}, false},
		{"US CORE", func() trade.OrderRequest {
			o := orderFor(trade.MarketUS, trade.OrderTypeLimit)
			o.SupportTradingSession = trade.TradingSessionCore
			return o
		}, false},
		{"US ALL", func() trade.OrderRequest {
			o := orderFor(trade.MarketUS, trade.OrderTypeLimit)
			o.SupportTradingSession = trade.TradingSessionAll
			return o
		}, false},
		{"US NIGHT", func() trade.OrderRequest {
			o := orderFor(trade.MarketUS, trade.OrderTypeLimit)
			o.SupportTradingSession = trade.TradingSessionNight
			return o
		}, false},
		{"US ALL_DAY", func() trade.OrderRequest {
			o := orderFor(trade.MarketUS, trade.OrderTypeLimit)
			o.SupportTradingSession = trade.TradingSessionAllDay
			return o
		}, false},
		{"US deprecated Y", func() trade.OrderRequest {
			o := orderFor(trade.MarketUS, trade.OrderTypeLimit)
			o.SupportTradingSession = trade.TradingSessionY
			return o
		}, true},
		{"US deprecated N", func() trade.OrderRequest {
			o := orderFor(trade.MarketUS, trade.OrderTypeLimit)
			o.SupportTradingSession = trade.TradingSessionN
			return o
		}, true},
		{"HK session rejected", func() trade.OrderRequest {
			o := orderFor(trade.MarketHK, trade.OrderTypeEnhancedLimit)
			o.SupportTradingSession = trade.TradingSessionCore
			return o
		}, true},
		{"CN session rejected", func() trade.OrderRequest {
			o := orderFor(trade.MarketCN, trade.OrderTypeLimit)
			o.SupportTradingSession = trade.TradingSessionAll
			return o
		}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.order().Validate()
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

// TestSandboxPreviewHKEnhancedLimit previews a small, non-marketable HK
// ENHANCED_LIMIT buy for 00700 (a holding of the dedicated sandbox account),
// exercising the BCAN requirement against the live sandbox. It is read-only:
// preview never places or mutates an order.
//
// The account's real BCAN party id is per-account and is never committed, so it
// is read from WEBULL_TRADE_PARTY_ID. The test is skipped unless
// WEBULL_TRADE_SANDBOX=1 and WEBULL_TRADE_PARTY_ID are set. A sandbox rejection
// of the BCAN payload is tolerated as a typed [errs.Error] rather than a
// failure, so the rule can be observed without a certified party id.
func TestSandboxPreviewHKEnhancedLimit(t *testing.T) {
	partyID := os.Getenv("WEBULL_TRADE_PARTY_ID")
	if partyID == "" {
		t.Skip("set WEBULL_TRADE_PARTY_ID to run the HK BCAN sandbox preview")
	}
	trading, accountID := newSandboxTradeClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := trade.PlaceOrderRequest{
		AccountID: accountID,
		NewOrders: []trade.OrderRequest{{
			ClientOrderID:  "sdk-preview-00700-1",
			ComboType:      trade.ComboTypeNormal,
			InstrumentType: trade.InstrumentTypeEquity,
			Market:         trade.MarketHK,
			Symbol:         "00700",
			OrderType:      trade.OrderTypeEnhancedLimit,
			Side:           trade.OrderSideBuy,
			Quantity:       "100",
			EntrustType:    trade.EntrustTypeQty,
			TimeInForce:    trade.TimeInForceDay,
			LimitPrice:     "1.00",
			NoPartyIDs: []trade.PartyID{{
				PartyID:       partyID,
				PartyIDSource: "D",
				PartyRole:     "3",
			}},
		}},
	}

	res, err := trading.PreviewOrder(ctx, req)
	if err != nil {
		var e *errs.Error
		if !errors.As(err, &e) {
			t.Fatalf("PreviewOrder() error = %v (%T), want *errs.Error", err, err)
		}
		t.Logf("PreviewOrder(00700) typed error (tolerated): code=%s message=%s", e.Code, e.Message)
		return
	}
	t.Logf("preview 00700 estimated_cost=%s estimated_transaction_fee=%s",
		res.EstimatedCost, res.EstimatedTransactionFee)
}

func TestAuctionPriceRules(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		order   func() trade.OrderRequest
		wantErr bool
	}{
		{"HK AT_AUCTION without limit_price", func() trade.OrderRequest {
			return orderFor(trade.MarketHK, trade.OrderTypeAtAuction)
		}, false},
		{"HK AT_AUCTION with limit_price", func() trade.OrderRequest {
			o := orderFor(trade.MarketHK, trade.OrderTypeAtAuction)
			o.LimitPrice = "10.00"
			return o
		}, true},
		{"HK AT_AUCTION_LIMIT with limit_price", func() trade.OrderRequest {
			return orderFor(trade.MarketHK, trade.OrderTypeAtAuctionLimit)
		}, false},
		{"HK AT_AUCTION_LIMIT without limit_price", func() trade.OrderRequest {
			o := orderFor(trade.MarketHK, trade.OrderTypeAtAuctionLimit)
			o.LimitPrice = ""
			return o
		}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.order().Validate()
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
