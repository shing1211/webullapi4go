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
	"fmt"
	"strings"
	"time"

	"github.com/shing1211/webullapi4go/pkg/domain/money"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

// optionOrderTypes are the only order types the API accepts for single-leg
// option orders. MARKET is deliberately excluded; single-leg options support
// LIMIT, STOP_LOSS and STOP_LOSS_LIMIT only.
var optionOrderTypes = []OrderType{
	OrderTypeLimit,
	OrderTypeStopLoss,
	OrderTypeStopLossLimit,
}

// optionMultiLegOrderTypes is the set of order types accepted for multi-leg
// option orders: LIMIT and STOP_LOSS_LIMIT.
var optionMultiLegOrderTypes = []OrderType{
	OrderTypeLimit,
	OrderTypeStopLossLimit,
}

// optionExpireDateLayout is the required wire format of an option expiration
// date.
const optionExpireDateLayout = "2006-01-02"

// Validate reports whether l is a well-formed single-leg option order. It checks
// the leg instrument, market, symbol, side, strike price, expiration date,
// option type, and quantity, and returns a typed [errs.Error] with
// [errs.CodeInvalidConfig] on the first problem found.
//
// A leg must use instrument_type OPTION, market US, and side BUY or SELL
// (SHORT is rejected); option_expire_date must be in YYYY-MM-DD form. This is
// the same validation [OrderRequest.Validate] applies to each option order's
// single leg.
func (l OrderLeg) Validate() error { return checkOptionLeg(l, optionLegFail) }

// checkOptionLeg is the shared implementation of [OrderLeg.Validate] and the
// option-order leg check. It reports the first problem through fail so callers
// can locate the leg within a batch.
func checkOptionLeg(l OrderLeg, fail func(string, ...any) error) error {
	if l.InstrumentType != InstrumentTypeOption {
		return fail("instrument_type must be OPTION, got %q", l.InstrumentType)
	}
	if l.Market != MarketUS {
		return fail("market must be US, got %q", l.Market)
	}
	if strings.TrimSpace(l.Symbol) == "" {
		return fail("symbol is required")
	}
	switch l.Side {
	case OrderSideBuy, OrderSideSell:
	default:
		return fail("side %q must be BUY or SELL", l.Side)
	}
	if l.StrikePrice == nil {
		return fail("strike_price is required")
	}
	if !l.StrikePrice.IsPositive() {
		return fail("strike_price must be a positive decimal number")
	}
	if strings.TrimSpace(l.OptionExpireDate) == "" {
		return fail("option_expire_date is required")
	}
	if !validOptionExpireDate(l.OptionExpireDate) {
		return fail("option_expire_date %q must be in YYYY-MM-DD format", l.OptionExpireDate)
	}
	switch l.OptionType {
	case OptionTypeCall, OptionTypePut:
	default:
		return fail("option_type %q must be CALL or PUT", l.OptionType)
	}
	if l.Quantity == nil {
		return fail("quantity is required")
	}
	if !l.Quantity.IsPositive() {
		return fail("quantity must be a positive decimal number")
	}
	return nil
}

// optionLegFail builds a typed [errs.Error] with [errs.CodeInvalidConfig] for a
// leg validation problem.
func optionLegFail(format string, args ...any) error {
	return errs.New(errs.CodeInvalidConfig, fmt.Sprintf(format, args...))
}

// validOptionExpireDate reports whether s is a calendar date in the required
// YYYY-MM-DD form. The round-trip format comparison rejects zero-padding and
// any format the parser would otherwise tolerate.
func validOptionExpireDate(s string) bool {
	t, err := time.Parse(optionExpireDateLayout, s)
	if err != nil {
		return false
	}
	return t.Format(optionExpireDateLayout) == s
}

// valid reports whether s is a recognized option strategy.
func (s OptionStrategy) valid() bool {
	switch s {
	case OptionStrategySingle, OptionStrategyVertical, OptionStrategyStraddle,
		OptionStrategyStrangle, OptionStrategyIronCondor, OptionStrategyIronButterfly,
		OptionStrategyButterfly, OptionStrategyCollar, OptionStrategyCalendar,
		OptionStrategyDiagonal, OptionStrategyRatio:
		return true
	default:
		return false
	}
}

// validateOptionRules enforces the option-specific constraints that can be
// checked before any network call. It covers both directions of the option
// distinction:
//
//   - a non-OPTION order must not carry option_strategy or legs;
//   - an OPTION order requires a supported option_strategy, a side of BUY or
//     SELL, and, for a sell-side order, time_in_force DAY. GTD is rejected for
//     options and GTC is allowed only for buy-side orders;
//   - a SINGLE order requires exactly one leg and an order type of LIMIT,
//     STOP_LOSS or STOP_LOSS_LIMIT;
//   - a multi-leg order must carry at least two structurally well-formed legs
//     and a top-level order type of LIMIT or STOP_LOSS_LIMIT (see
//     [validateOptionLegSet]).
//
// fail formats and returns the caller's typed error with the batch prefix
// already applied, and it returns the first problem found.
func (r OrderRequest) validateOptionRules(fail func(string, ...any) error) error {
	if r.InstrumentType != InstrumentTypeOption {
		if r.OptionStrategy != "" {
			return fail("option_strategy is only valid for OPTION orders")
		}
		if len(r.Legs) > 0 {
			return fail("legs is only valid for OPTION orders")
		}
		return nil
	}

	if !r.OptionStrategy.valid() {
		return fail("option_strategy %q is not a supported option strategy", r.OptionStrategy)
	}

	multiLeg := r.OptionStrategy != OptionStrategySingle
	allowed := optionOrderTypes
	if multiLeg {
		allowed = optionMultiLegOrderTypes
	}
	if !containsOrderType(allowed, r.OrderType) {
		return fail("order_type %s is not supported for options; supported types: %s",
			r.OrderType, orderTypeList(allowed))
	}

	switch r.Side {
	case OrderSideBuy:
		if r.TimeInForce == TimeInForceGTD {
			return fail("time_in_force GTD is not supported for options")
		}
	case OrderSideSell:
		if r.TimeInForce != TimeInForceDay {
			return fail("time_in_force must be DAY for option sell-side orders, got %q", r.TimeInForce)
		}
	default:
		return fail("side %q is not supported for options; only BUY and SELL are allowed", r.Side)
	}

	if multiLeg {
		return validateOptionLegSet(r.Legs, fail)
	}

	switch len(r.Legs) {
	case 1:
	case 0:
		return fail("legs must contain exactly one leg for a SINGLE option order, got none")
	default:
		return fail("legs must contain exactly one leg for a SINGLE option order, got %d", len(r.Legs))
	}

	legFail := func(format string, args ...any) error { return fail("legs[0]: "+format, args...) }
	return checkOptionLeg(r.Legs[0], legFail)
}

// optionLegKey is the identity of an option leg for duplicate detection. Two
// legs with the same key are the same order line and must be combined rather
// than submitted twice.
type optionLegKey struct {
	symbol     string
	side       OrderSide
	strike     string
	expiration string
	optionType OptionType
}

// canonicalStrike returns a canonical form of an option strike price for
// equality comparison, so values that differ only in decimal precision (for
// example "220.0" and "220.00") compare equal. It normalizes through
// [parseDecimal] and does not alter the stored or emitted strike string. If the
// canonicalStrike returns a canonical string representation of the strike price
// by rounding to the decimal's string form. Since validation already ensures the
// value is a positive decimal, the parse failure path is defensive only.
func canonicalStrike(s *money.Money) string {
	if s == nil {
		return ""
	}
	return s.Decimal().String()
}

// validateOptionLegSet checks that legs form a structurally well-formed
// multi-leg option order: at least two legs, every leg validated in slice order
// through [checkOptionLeg], no duplicate legs, and no degenerate set whose legs
// all share a side, option type, and strike. Errors locate the offending leg as
// legs[i].
//
// It validates structure only: it does not price the strategy or evaluate its
// risk profile.
func validateOptionLegSet(legs []OrderLeg, fail func(string, ...any) error) error {
	if len(legs) < 2 {
		return fail("legs must contain at least two legs for a multi-leg option order, got %d", len(legs))
	}
	seen := make(map[optionLegKey]struct{}, len(legs))
	for i := range legs {
		legFail := func(format string, args ...any) error {
			return fail(fmt.Sprintf("legs[%d]: ", i)+format, args...)
		}
		if err := checkOptionLeg(legs[i], legFail); err != nil {
			return err
		}
		key := optionLegKey{
			symbol:     legs[i].Symbol,
			side:       legs[i].Side,
			strike:     canonicalStrike(legs[i].StrikePrice),
			expiration: legs[i].OptionExpireDate,
			optionType: legs[i].OptionType,
		}
		if _, dup := seen[key]; dup {
			return fail("legs[%d]: duplicate leg; each leg must differ in symbol, side, strike, expiration or type", i)
		}
		seen[key] = struct{}{}
	}
	if sameSideTypeStrike(legs) {
		return fail("legs must not all share the same side, option type and strike price")
	}
	return nil
}

// sameSideTypeStrike reports whether every leg shares the same side, option
// type, and strike price, which makes the leg set degenerate. See
// [validateOptionLegSet]. Strikes are compared numerically so values that
// differ only in decimal precision are treated as equal.
func sameSideTypeStrike(legs []OrderLeg) bool {
	first := legs[0]
	firstStrike := canonicalStrike(first.StrikePrice)
	for _, l := range legs[1:] {
		if l.Side != first.Side || l.OptionType != first.OptionType || canonicalStrike(l.StrikePrice) != firstStrike {
			return false
		}
	}
	return true
}
