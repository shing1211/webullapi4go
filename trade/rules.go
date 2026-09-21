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

import (
	"math/big"
	"strings"
)

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
// The matrix applies to equity orders. Option order types are enforced by
// [OrderRequest.validateOptionRules] and futures order types by
// [OrderRequest.validateFuturesRules].
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

// marketFuturesOrderTypes is the matrix of futures order types accepted by each
// futures-capable market. The set is deliberately conservative: only the
// widely-supported types are enabled until the exact futures support is
// confirmed.
//
// TODO(t9): confirm futures order-type matrix against live API.
var marketFuturesOrderTypes = map[Market][]OrderType{
	MarketUS: {
		OrderTypeLimit,
		OrderTypeMarket,
		OrderTypeStopLoss,
		OrderTypeStopLossLimit,
	},
	MarketHK: {
		OrderTypeLimit,
		OrderTypeMarket,
		OrderTypeStopLoss,
		OrderTypeStopLossLimit,
	},
}

// allowsFuturesOrderType reports whether m accepts t for a futures order.
func (m Market) allowsFuturesOrderType(t OrderType) bool {
	for _, allowed := range marketFuturesOrderTypes[m] {
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

// validateFuturesRules enforces the futures-specific constraints that can be
// checked before any network call. Futures orders are single-instrument,
// whole-contract orders, so it rejects the option and equity/HK fields that do
// not apply and applies the per-market futures order-type matrix:
//
//   - support_trading_session is a US-equity session selector and no_party_ids
//     is Hong Kong equity regulatory reporting; neither applies to futures and
//     both are rejected;
//   - option_strategy and legs are option-only and are rejected;
//   - the order type must be one the market accepts for futures;
//   - time_in_force must be DAY or GTC; GTD is rejected because futures
//     expire_date semantics are unconfirmed;
//   - entrust_type must be QTY, because futures are not sized by total cash
//     amount;
//   - quantity must be a positive integer, because futures trade whole
//     contracts and do not support fractional quantities.
//
// The time_in_force, entrust_type and whole-contract quantity rules below are
// provisional assumptions about the API, not confirmed guarantees.
//
// TODO(t9): confirm futures time_in_force, entrust_type and whole-contract quantity rules against live API
//
// fail formats and returns the caller's typed error with the batch prefix
// already applied, and it returns the first problem found.
func (r OrderRequest) validateFuturesRules(fail func(string, ...any) error) error {
	if r.SupportTradingSession != "" {
		return fail("support_trading_session is not valid for futures orders")
	}
	if len(r.NoPartyIDs) > 0 {
		return fail("no_party_ids is not valid for futures orders")
	}
	if r.OptionStrategy != "" {
		return fail("option_strategy is only valid for OPTION orders")
	}
	if len(r.Legs) > 0 {
		return fail("legs is only valid for OPTION orders")
	}
	if !r.Market.allowsFuturesOrderType(r.OrderType) {
		allowed, known := marketFuturesOrderTypes[r.Market]
		if !known {
			return fail("futures are not supported for this market")
		}
		return fail("order_type %s is not supported for %s futures orders; supported types: %s",
			r.OrderType, r.Market, orderTypeList(allowed))
	}
	switch r.TimeInForce {
	case TimeInForceDay, TimeInForceGTC:
	case TimeInForceGTD:
		return fail("time_in_force GTD is not supported for futures orders")
	default:
		return fail("time_in_force %q must be DAY or GTC for futures orders", r.TimeInForce)
	}
	if r.EntrustType != EntrustTypeQty {
		return fail("entrust_type %q must be QTY for futures orders", r.EntrustType)
	}
	if !isPositiveInteger(r.Quantity) {
		return fail("quantity %q must be a positive integer for futures orders", r.Quantity)
	}
	return nil
}

// validateEventRules enforces event-contract-specific constraints:
//   - LIMIT orders only
//   - DAY time-in-force only
//   - EntrustType must be QTY
//   - Quantity must be a positive integer (max 50,000)
//   - No option_strategy or legs
//   - No support_trading_session or no_party_ids
//
// TODO(event): confirm event contract order rules against live API.
func (r OrderRequest) validateEventRules(fail func(string, ...any) error) error {
	if r.SupportTradingSession != "" {
		return fail("support_trading_session is not valid for event contract orders")
	}
	if len(r.NoPartyIDs) > 0 {
		return fail("no_party_ids is not valid for event contract orders")
	}
	if r.OptionStrategy != "" {
		return fail("option_strategy is only valid for OPTION orders")
	}
	if len(r.Legs) > 0 {
		return fail("legs is only valid for OPTION orders")
	}
	if r.OrderType != OrderTypeLimit {
		return fail("order_type %s is not supported for event contract orders; only LIMIT is supported", r.OrderType)
	}
	if r.TimeInForce != TimeInForceDay {
		return fail("time_in_force %s is not supported for event contract orders; only DAY is supported", r.TimeInForce)
	}
	if r.EntrustType != EntrustTypeQty {
		return fail("entrust_type %q must be QTY for event contract orders", r.EntrustType)
	}
	if !isPositiveInteger(r.Quantity) {
		return fail("quantity %q must be a positive integer for event contract orders", r.Quantity)
	}
	if qty, ok := new(big.Int).SetString(strings.TrimSpace(r.Quantity), 10); ok {
		max := new(big.Int).SetInt64(50000)
		if qty.Cmp(max) > 0 {
			return fail("quantity %s exceeds the maximum 50000 contracts for event contract orders", r.Quantity)
		}
	}
	return nil
}

// isPositiveInteger reports whether s is a base-10 integer string greater than
// zero. It rejects signs, decimal points, and any non-digit character, so a
// fractional futures quantity such as "1.5" fails.
func isPositiveInteger(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	n, ok := new(big.Int).SetString(s, 10)
	return ok && n.Sign() > 0
}
