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

package trade

import "strings"

// marketEquityOrderTypes is the authoritative matrix of equity order types
// accepted by each supported market. It reflects the Webull Stock Trading API
// reference:
//
//   - US accepts limit, market, stop, stop-limit, market-on-open,
//     market-on-close, touch, and trailing stop orders.
//   - HK accepts enhanced limit, at-auction, at-auction limit, stop,
//     stop-limit, touch, and trailing stop orders.
//   - CN (A-share Stock Connect) accepts only LIMIT. A-share trading is
//     disabled by default server-side and must be enabled by Webull support.
//
// The matrix applies to equity orders. Option and futures order-type rules are
// enforced by their own validation.
var marketEquityOrderTypes = map[Market][]OrderType{
	MarketUS: {
		OrderTypeLimit,
		OrderTypeMarket,
		OrderTypeStopLoss,
		OrderTypeStopLossLimit,
		OrderTypeMarketOnOpen,
		OrderTypeMarketOnClose,
		OrderTypeTouchMkt,
		OrderTypeTouchLmt,
		OrderTypeTrailingStopLoss,
		OrderTypeTrailingStopLossLimit,
	},
	MarketHK: {
		OrderTypeEnhancedLimit,
		OrderTypeAtAuction,
		OrderTypeAtAuctionLimit,
		OrderTypeStopLoss,
		OrderTypeStopLossLimit,
		OrderTypeTouchMkt,
		OrderTypeTouchLmt,
		OrderTypeTrailingStopLoss,
		OrderTypeTrailingStopLossLimit,
	},
	MarketCN: {
		OrderTypeLimit,
	},
}

// allowsEquityOrderType reports whether m accepts t for an equity order.
func (m Market) allowsEquityOrderType(t OrderType) bool {
	for _, allowed := range marketEquityOrderTypes[m] {
		if allowed == t {
			return true
		}
	}
	return false
}

// orderTypeList renders ts as a comma-separated list for an error message.
func orderTypeList(ts []OrderType) string {
	parts := make([]string, len(ts))
	for i, t := range ts {
		parts[i] = string(t)
	}
	return strings.Join(parts, ", ")
}

// validateMarketRules enforces the market-specific order constraints that can
// be checked before any network call: the equity order-type matrix, the US-only
// trading session, the Hong Kong BCAN party identifiers, and the at-auction
// price rules.
//
// fail formats and returns the caller's typed error, and prefix locates the
// order within a batch. The rules are ordered so the most fundamental problem
// is reported first, and it returns the first problem found:
//
//   - an order type the market does not accept,
//   - a trading session on a non-US order, or a deprecated session,
//   - a missing or malformed HK BCAN party identifier, or party identifiers on
//     a non-HK-equity order,
//   - a limit price on an AT_AUCTION order.
//
// The requirement that AT_AUCTION_LIMIT carries a limit price is enforced by
// [OrderType.needsLimitPrice].
func (r OrderRequest) validateMarketRules(fail func(string, ...any) error) error {
	if r.InstrumentType == InstrumentTypeEquity && !r.Market.allowsEquityOrderType(r.OrderType) {
		return fail("order_type %s is not supported for %s equity orders; supported types: %s",
			r.OrderType, r.Market, orderTypeList(marketEquityOrderTypes[r.Market]))
	}

	if r.SupportTradingSession != "" {
		if r.Market != MarketUS {
			return fail("support_trading_session is only valid for US orders")
		}
		switch r.SupportTradingSession {
		case TradingSessionY, TradingSessionN:
			return fail("support_trading_session %q is deprecated; use CORE, ALL, NIGHT or ALL_DAY",
				r.SupportTradingSession)
		}
	}

	isHKEQUITY := r.Market == MarketHK && r.InstrumentType == InstrumentTypeEquity
	if len(r.NoPartyIDs) > 0 && !isHKEQUITY {
		return fail("no_party_ids is only valid for HK equity orders")
	}
	if isHKEQUITY {
		if len(r.NoPartyIDs) == 0 {
			return fail("no_party_ids is required for HK equity orders")
		}
		for i, p := range r.NoPartyIDs {
			switch {
			case strings.TrimSpace(p.PartyID) == "":
				return fail("no_party_ids[%d]: party_id is required", i)
			case p.PartyIDSource != "D":
				return fail("no_party_ids[%d]: party_id_source must be %q, got %q", i, "D", p.PartyIDSource)
			case p.PartyRole != "3":
				return fail("no_party_ids[%d]: party_role must be %q, got %q", i, "3", p.PartyRole)
			}
		}
	}

	if r.OrderType == OrderTypeAtAuction && strings.TrimSpace(r.LimitPrice) != "" {
		return fail("limit_price must not be set for AT_AUCTION orders")
	}

	return nil
}
