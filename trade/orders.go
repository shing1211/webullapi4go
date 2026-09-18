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
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/shing1211/webullapi4go/internal/errs"
)

// Stock and option order endpoint paths.
const (
	// pathOrdersPreview estimates the cost of an order without placing it.
	pathOrdersPreview = "/trading/orders/preview"
	// pathOrdersPlace submits one or more orders.
	pathOrdersPlace = "/trading/orders/place"
)

// maxClientOrderIDLength is the documented maximum length of a client order
// identifier.
const maxClientOrderIDLength = 32

// OrderSide is the intended trading direction of an order.
type OrderSide string

// Order sides accepted by the trading API.
const (
	// OrderSideBuy buys or opens a long position.
	OrderSideBuy OrderSide = "BUY"
	// OrderSideSell sells or closes a long position.
	OrderSideSell OrderSide = "SELL"
	// OrderSideShort opens or closes a short position.
	OrderSideShort OrderSide = "SHORT"
)

// OrderType is the execution instruction for an order. The set of types that is
// valid for an order depends on the market and instrument; that matrix is
// enforced in later releases.
type OrderType string

// Order types accepted by the trading API.
const (
	// OrderTypeLimit executes at limit_price or better.
	OrderTypeLimit OrderType = "LIMIT"
	// OrderTypeMarket executes at the best available price.
	OrderTypeMarket OrderType = "MARKET"
	// OrderTypeStopLoss becomes a market order at stop_price.
	OrderTypeStopLoss OrderType = "STOP_LOSS"
	// OrderTypeStopLossLimit becomes a limit order at stop_price.
	OrderTypeStopLossLimit OrderType = "STOP_LOSS_LIMIT"
	// OrderTypeEnhancedLimit is a Hong Kong enhanced limit order.
	OrderTypeEnhancedLimit OrderType = "ENHANCED_LIMIT"
	// OrderTypeAtAuction is a Hong Kong at-auction order.
	OrderTypeAtAuction OrderType = "AT_AUCTION"
	// OrderTypeAtAuctionLimit is a Hong Kong at-auction limit order.
	OrderTypeAtAuctionLimit OrderType = "AT_AUCTION_LIMIT"
	// OrderTypeMarketOnOpen executes at the opening price.
	OrderTypeMarketOnOpen OrderType = "MARKET_ON_OPEN"
	// OrderTypeMarketOnClose executes at the closing price.
	OrderTypeMarketOnClose OrderType = "MARKET_ON_CLOSE"
	// OrderTypeTrailingStopLoss trails a stop price by trailing_stop_step.
	OrderTypeTrailingStopLoss OrderType = "TRAILING_STOP_LOSS"
	// OrderTypeTrailingStopLossLimit trails a stop that submits a limit order.
	OrderTypeTrailingStopLossLimit OrderType = "TRAILING_STOP_LOSS_LIMIT"
	// OrderTypeTouchMkt becomes a market order when the trigger price is touched.
	OrderTypeTouchMkt OrderType = "TOUCH_MKT"
	// OrderTypeTouchLmt becomes a limit order when the trigger price is touched.
	OrderTypeTouchLmt OrderType = "TOUCH_LMT"
)

// TimeInForce is how long an order remains active.
type TimeInForce string

// Time-in-force values accepted by the trading API.
const (
	// TimeInForceDay expires at the end of the trading day.
	TimeInForceDay TimeInForce = "DAY"
	// TimeInForceGTD expires on expire_date; currently US only.
	TimeInForceGTD TimeInForce = "GTD"
	// TimeInForceGTC remains active until filled or cancelled.
	TimeInForceGTC TimeInForce = "GTC"
)

// ComboType identifies the role an order plays within a combo order.
// OTO, OCO, and OTOCO are equity-only; their leg-count rules are enforced in a
// later release.
type ComboType string

// Combo types accepted by the trading API.
const (
	// ComboTypeNormal is a standard single order.
	ComboTypeNormal ComboType = "NORMAL"
	// ComboTypeMaster is the primary order that triggers its siblings.
	ComboTypeMaster ComboType = "MASTER"
	// ComboTypeStopProfit is a take-profit sub-order.
	ComboTypeStopProfit ComboType = "STOP_PROFIT"
	// ComboTypeStopLoss is a stop-loss sub-order.
	ComboTypeStopLoss ComboType = "STOP_LOSS"
	// ComboTypeOTO is the follow-up order of a one-triggers-the-other pair.
	ComboTypeOTO ComboType = "OTO"
	// ComboTypeOCO is one of a pair where filling either cancels the other.
	ComboTypeOCO ComboType = "OCO"
	// ComboTypeOTOCO is one of the order set triggered by a OTOCO master.
	ComboTypeOTOCO ComboType = "OTOCO"
)

// EntrustType is how an order size is expressed.
type EntrustType string

// Entrust types accepted by the trading API.
const (
	// EntrustTypeQty sizes an order by quantity of shares or contracts.
	EntrustTypeQty EntrustType = "QTY"
	// EntrustTypeAmount sizes an order by total cash amount, for US
	// fractional share trading.
	EntrustTypeAmount EntrustType = "AMOUNT"
)

// TradingSession is the US trading session an order may execute in.
type TradingSession string

// Trading-session values accepted by the trading API.
const (
	// TradingSessionY includes extended hours. Deprecated by the API; prefer
	// TradingSessionAll.
	TradingSessionY TradingSession = "Y"
	// TradingSessionN restricts to regular hours. Deprecated by the API;
	// prefer TradingSessionCore.
	TradingSessionN TradingSession = "N"
	// TradingSessionNight restricts to night trading.
	TradingSessionNight TradingSession = "NIGHT"
	// TradingSessionAll includes extended hours.
	TradingSessionAll TradingSession = "ALL"
	// TradingSessionCore restricts to regular trading hours.
	TradingSessionCore TradingSession = "CORE"
	// TradingSessionAllDay includes overnight hours from 8:00 p.m. ET to
	// 8:00 p.m. ET the next day.
	TradingSessionAllDay TradingSession = "ALL_DAY"
)

// TriggerPriceType is the market price a touch or stop order triggers on.
type TriggerPriceType string

// Trigger price types accepted by the trading API.
const (
	// TriggerPriceTypePrice uses the latest transaction price.
	TriggerPriceTypePrice TriggerPriceType = "PRICE"
	// TriggerPriceTypePriceBid uses the best bid.
	TriggerPriceTypePriceBid TriggerPriceType = "PRICE_BID"
	// TriggerPriceTypePriceAsk uses the best ask.
	TriggerPriceTypePriceAsk TriggerPriceType = "PRICE_ASK"
)

// TrailingType is how a trailing stop step is expressed.
type TrailingType string

// Trailing types accepted by the trading API.
const (
	// TrailingTypeAmount trails by a fixed price amount.
	TrailingTypeAmount TrailingType = "AMOUNT"
	// TrailingTypePercentage trails by a percentage, where "0.01" is 1%.
	TrailingTypePercentage TrailingType = "PERCENTAGE"
)

// OrderStatus is the lifecycle state of an order.
type OrderStatus string

// Order lifecycle states returned by the trading API.
const (
	// OrderStatusPending is awaiting submission.
	OrderStatusPending OrderStatus = "PENDING"
	// OrderStatusSubmitted is accepted and working.
	OrderStatusSubmitted OrderStatus = "SUBMITTED"
	// OrderStatusCancelled was cancelled before completing.
	OrderStatusCancelled OrderStatus = "CANCELLED"
	// OrderStatusFilled completed in full.
	OrderStatusFilled OrderStatus = "FILLED"
	// OrderStatusFailed was rejected or failed.
	OrderStatusFailed OrderStatus = "FAILED"
	// OrderStatusPartialFilled completed in part.
	OrderStatusPartialFilled OrderStatus = "PARTIAL_FILLED"
)

// OrderLeg is one leg of an option order. It is only populated for option
// orders; validation of the leg set is added in a later release.
type OrderLeg struct {
	// InstrumentType is the kind of instrument the leg references.
	InstrumentType InstrumentType `json:"instrument_type"`
	// Market is the market the leg trades in.
	Market Market `json:"market"`
	// Symbol is the trading symbol of the leg.
	Symbol string `json:"symbol"`
	// Side is the intended direction of the leg.
	Side OrderSide `json:"side"`
	// StrikePrice is the option strike, as a decimal string.
	StrikePrice string `json:"strike_price,omitempty"`
	// OptionExpireDate is the option expiry in yyyy-MM-dd form.
	OptionExpireDate string `json:"option_expire_date,omitempty"`
	// OptionType is whether the leg is a call or a put.
	OptionType OptionType `json:"option_type,omitempty"`
	// Quantity is the leg quantity, as a decimal string.
	Quantity string `json:"quantity,omitempty"`
}

// PartyID identifies a party to a Hong Kong order for regulatory reporting.
// It is relevant only for Hong Kong stock orders.
type PartyID struct {
	// PartyID is the broker client identifier, for example "ABC123.2568".
	PartyID string `json:"party_id"`
	// PartyIDSource is the identifier source; the API requires "D".
	PartyIDSource string `json:"party_id_source"`
	// PartyRole is the party role; the API requires "3".
	PartyRole string `json:"party_role"`
}

// OrderRequest is a single order within a [PlaceOrderRequest]. Numeric values
// are strings to preserve precision.
//
// The fields required for every order are ClientOrderID, ComboType,
// InstrumentType, Market, Symbol, OrderType, Side, EntrustType, and
// TimeInForce. Conditional requirements are enforced by [OrderRequest.Validate].
type OrderRequest struct {
	// ClientOrderID is a caller-supplied identifier, at most 32 characters of
	// [A-Za-z0-9_-], unique per account.
	ClientOrderID string `json:"client_order_id"`
	// ComboType is the order's role within a combo order.
	ComboType ComboType `json:"combo_type"`
	// InstrumentType is the kind of instrument to trade.
	InstrumentType InstrumentType `json:"instrument_type"`
	// Market is the market to route the order to.
	Market Market `json:"market"`
	// Symbol is the trading symbol of the instrument.
	Symbol string `json:"symbol"`
	// OrderType is the execution instruction.
	OrderType OrderType `json:"order_type"`
	// Side is the intended trading direction.
	Side OrderSide `json:"side"`
	// Quantity is the order quantity, as a decimal string. Required when
	// EntrustType is QTY; ignored when EntrustType is AMOUNT.
	Quantity string `json:"quantity,omitempty"`
	// EntrustType is whether the order is sized by quantity or cash amount.
	EntrustType EntrustType `json:"entrust_type"`
	// TimeInForce is how long the order remains active.
	TimeInForce TimeInForce `json:"time_in_force"`
	// SupportTradingSession restricts a US order to a trading session. It is
	// optional; the API applies its own default when it is empty.
	SupportTradingSession TradingSession `json:"support_trading_session,omitempty"`
	// LimitPrice is the limit price, as a decimal string. Required for
	// LIMIT, STOP_LOSS_LIMIT, and TOUCH_LMT orders.
	LimitPrice string `json:"limit_price,omitempty"`
	// StopPrice is the trigger price, as a decimal string. Required for
	// STOP_LOSS, STOP_LOSS_LIMIT, TOUCH_MKT, and TOUCH_LMT orders.
	StopPrice string `json:"stop_price,omitempty"`
	// TotalCashAmount is the cash amount, as a decimal string. Required when
	// EntrustType is AMOUNT.
	TotalCashAmount string `json:"total_cash_amount,omitempty"`
	// TriggerPriceType is the market price a touch or stop order triggers on.
	TriggerPriceType TriggerPriceType `json:"trigger_price_type,omitempty"`
	// TrailingType is how TrailingStopStep is expressed.
	TrailingType TrailingType `json:"trailing_type,omitempty"`
	// TrailingStopStep is the trailing spread, as a decimal string.
	TrailingStopStep string `json:"trailing_stop_step,omitempty"`
	// TrailingLimitPriceOffset is the offset between the triggered stop price
	// and the submitted limit price, as a decimal string.
	TrailingLimitPriceOffset string `json:"trailing_limit_price_offset,omitempty"`
	// ExpireDate is the GTD expiry in yyyy-MM-dd form. Required when
	// TimeInForce is GTD.
	ExpireDate string `json:"expire_date,omitempty"`
	// OptionStrategy identifies the option strategy.
	OptionStrategy OptionStrategy `json:"option_strategy,omitempty"`
	// Legs lists the option legs. It is empty for non-option orders.
	Legs []OrderLeg `json:"legs,omitempty"`
	// NoPartyIDs lists the regulatory parties for a Hong Kong order.
	NoPartyIDs []PartyID `json:"no_party_ids,omitempty"`
	// SenderSubID identifies the firm or sub-account on whose behalf a
	// third-party order is placed.
	SenderSubID string `json:"sender_sub_id,omitempty"`
}

// PlaceOrderRequest is the body shared by the preview and place endpoints. It
// carries one or more orders and, for combo orders, an optional grouping
// identifier.
type PlaceOrderRequest struct {
	// AccountID is the account the orders belong to.
	AccountID string `json:"account_id"`
	// ClientComboOrderID groups the orders of a combo. It is optional for
	// NORMAL orders and auto-generated by the server for combo orders that
	// omit it.
	ClientComboOrderID string `json:"client_combo_order_id,omitempty"`
	// NewOrders are the orders to preview or place. At least one is required.
	NewOrders []OrderRequest `json:"new_orders"`
}

// PlaceOrderResult is the response of [Client.PlaceOrder] for a simple order.
type PlaceOrderResult struct {
	// ClientOrderID echoes the caller-supplied order identifier.
	ClientOrderID string `json:"client_order_id"`
	// OrderID is the system-generated order identifier.
	OrderID string `json:"order_id"`
}

// PreviewResult is the estimated cost of an order, as returned by
// [Client.PreviewOrder]. Amounts are decimal strings and the actual values may
// differ based on execution.
type PreviewResult struct {
	// EstimatedCost is the estimated capital required for the order.
	EstimatedCost string `json:"estimated_cost"`
	// EstimatedTransactionFee is the estimated transaction fee, including
	// exchange, clearing, and commission fees.
	EstimatedTransactionFee string `json:"estimated_transaction_fee"`
}

// Validate reports whether r is well formed. It checks the required fields, the
// client order identifier's length and character set, the enum fields, and the
// conditional price and size fields. It returns a typed [errs.Error] with
// [errs.CodeInvalidConfig] and an actionable message on the first problem
// found, so callers can reject a bad order before any network call.
func (r OrderRequest) Validate() error { return r.validate("") }

// validate is the shared implementation of [OrderRequest.Validate]. prefix, when
// non-empty, is prepended to the message to locate the order within a batch.
func (r OrderRequest) validate(prefix string) error {
	fail := func(format string, args ...any) error {
		return errs.New(errs.CodeInvalidConfig, prefix+fmt.Sprintf(format, args...))
	}

	switch {
	case strings.TrimSpace(r.ClientOrderID) == "":
		return fail("client_order_id is required")
	case len(r.ClientOrderID) > maxClientOrderIDLength:
		return fail("client_order_id must be at most %d characters, got %d", maxClientOrderIDLength, len(r.ClientOrderID))
	case !validClientOrderID(r.ClientOrderID):
		return fail("client_order_id %q may contain only letters, digits, '-' and '_'", r.ClientOrderID)
	}

	if r.ComboType == "" {
		return fail("combo_type is required")
	}
	if !r.ComboType.valid() {
		return fail("combo_type %q is not one of %s", r.ComboType, comboTypeList)
	}
	if !r.InstrumentType.valid() {
		return fail("instrument_type %q must be EQUITY, OPTION or FUTURES", r.InstrumentType)
	}
	if !r.Market.valid() {
		return fail("market %q must be US, HK or CN", r.Market)
	}
	if strings.TrimSpace(r.Symbol) == "" {
		return fail("symbol is required")
	}
	if !r.OrderType.valid() {
		return fail("order_type %q is not a supported order type", r.OrderType)
	}
	if !r.Side.valid() {
		return fail("side %q must be BUY, SELL or SHORT", r.Side)
	}
	if !r.TimeInForce.valid() {
		return fail("time_in_force %q must be DAY, GTD or GTC", r.TimeInForce)
	}
	if r.SupportTradingSession != "" && !r.SupportTradingSession.valid() {
		return fail("support_trading_session %q is not a supported session", r.SupportTradingSession)
	}
	if r.TriggerPriceType != "" && !r.TriggerPriceType.valid() {
		return fail("trigger_price_type %q must be PRICE, PRICE_BID or PRICE_ASK", r.TriggerPriceType)
	}
	if r.TrailingType != "" && !r.TrailingType.valid() {
		return fail("trailing_type %q must be AMOUNT or PERCENTAGE", r.TrailingType)
	}

	switch r.EntrustType {
	case EntrustTypeQty:
		if strings.TrimSpace(r.Quantity) == "" {
			return fail("quantity is required when entrust_type is QTY")
		}
		if !isPositiveDecimal(r.Quantity) {
			return fail("quantity %q must be a positive decimal number", r.Quantity)
		}
	case EntrustTypeAmount:
		if strings.TrimSpace(r.TotalCashAmount) == "" {
			return fail("total_cash_amount is required when entrust_type is AMOUNT")
		}
		if !isPositiveDecimal(r.TotalCashAmount) {
			return fail("total_cash_amount %q must be a positive decimal number", r.TotalCashAmount)
		}
	default:
		return fail("entrust_type %q must be QTY or AMOUNT", r.EntrustType)
	}

	if r.OrderType.needsLimitPrice() && strings.TrimSpace(r.LimitPrice) == "" {
		return fail("limit_price is required for %s orders", r.OrderType)
	}
	if r.OrderType.needsStopPrice() && strings.TrimSpace(r.StopPrice) == "" {
		return fail("stop_price is required for %s orders", r.OrderType)
	}
	if r.OrderType.isTrailing() {
		if r.TrailingType == "" {
			return fail("trailing_type is required for %s orders", r.OrderType)
		}
		if strings.TrimSpace(r.TrailingStopStep) == "" {
			return fail("trailing_stop_step is required for %s orders", r.OrderType)
		}
	}
	if r.TimeInForce == TimeInForceGTD && strings.TrimSpace(r.ExpireDate) == "" {
		return fail("expire_date is required when time_in_force is GTD")
	}
	return nil
}

// Validate reports whether r is well formed: the account is set, at least one
// order is present, every order passes [OrderRequest.Validate], and no two
// orders reuse a client order identifier. It returns a typed [errs.Error] with
// [errs.CodeInvalidConfig] on the first problem found.
func (r PlaceOrderRequest) Validate() error {
	if strings.TrimSpace(r.AccountID) == "" {
		return errs.New(errs.CodeInvalidConfig, "account_id is required")
	}
	if len(r.NewOrders) == 0 {
		return errs.New(errs.CodeInvalidConfig, "new_orders must contain at least one order")
	}
	seen := make(map[string]struct{}, len(r.NewOrders))
	for i := range r.NewOrders {
		order := &r.NewOrders[i]
		if err := order.validate(fmt.Sprintf("new_orders[%d]: ", i)); err != nil {
			return err
		}
		if _, dup := seen[order.ClientOrderID]; dup {
			return errs.New(errs.CodeInvalidConfig,
				fmt.Sprintf("new_orders[%d]: client_order_id %q is duplicated within the request", i, order.ClientOrderID))
		}
		seen[order.ClientOrderID] = struct{}{}
	}
	return nil
}

// PreviewOrder estimates the cost of the orders in req without placing them.
// The request is validated and the configured order guardrails are enforced
// before any network call, so a rejected request never reaches the API.
//
// Reference: https://developer.webull.hk/apis/docs/reference/common-order-preview.md
func (c *Client) PreviewOrder(ctx context.Context, req PlaceOrderRequest) (*PreviewResult, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := c.enforceGuardrails(req); err != nil {
		return nil, err
	}
	var out PreviewResult
	if err := c.do(ctx, http.MethodPost, pathOrdersPreview, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PlaceOrder submits the orders in req and returns the resulting identifiers
// for a simple order. The request is validated and the configured order
// guardrails are enforced before any network call, so a rejected request never
// reaches the API.
//
// PlaceOrder can create live orders. Prefer [Client.PreviewOrder] to validate an
// order first, and configure the guardrails with [WithMaxOrderNotional] and
// [WithMaxOrderQuantity] to bound what can be sent.
//
// Reference: https://developer.webull.hk/apis/docs/reference/common-order-place.md
func (c *Client) PlaceOrder(ctx context.Context, req PlaceOrderRequest) (*PlaceOrderResult, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := c.enforceGuardrails(req); err != nil {
		return nil, err
	}
	var out PlaceOrderResult
	if err := c.do(ctx, http.MethodPost, pathOrdersPlace, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// enforceGuardrails applies the configured order caps to every order in req. It
// runs before the network call and returns a typed [errs.Error] with
// [errs.CodeInvalidConfig] when an order exceeds a cap.
//
// The quantity cap compares each order's quantity. The notional cap compares
// total_cash_amount for AMOUNT orders, and quantity times limit_price for orders
// that carry both; when a notional cannot be computed (for example a MARKET
// order with no limit price) the notional cap is skipped, but the quantity cap
// still applies.
func (c *Client) enforceGuardrails(req PlaceOrderRequest) error {
	for i := range req.NewOrders {
		if err := c.enforceOrderGuardrails(&req.NewOrders[i]); err != nil {
			var e *errs.Error
			if errors.As(err, &e) {
				return errs.New(e.Code, fmt.Sprintf("new_orders[%d]: %s", i, e.Message))
			}
			return err
		}
	}
	return nil
}

// enforceOrderGuardrails applies the configured caps to a single order.
func (c *Client) enforceOrderGuardrails(o *OrderRequest) error {
	if c.cfg.maxOrderQuantity != "" && strings.TrimSpace(o.Quantity) != "" {
		if qty, ok := parseDecimal(o.Quantity); ok {
			if max, ok := parseDecimal(c.cfg.maxOrderQuantity); ok && qty.Cmp(max) > 0 {
				return errs.New(errs.CodeInvalidConfig,
					fmt.Sprintf("quantity %s exceeds the configured maximum %s", o.Quantity, c.cfg.maxOrderQuantity))
			}
		}
	}
	if c.cfg.maxOrderNotional == "" {
		return nil
	}
	notional, ok := orderNotional(o)
	if !ok {
		return nil
	}
	max, ok := parseDecimal(c.cfg.maxOrderNotional)
	if !ok {
		return nil
	}
	if notional.Cmp(max) > 0 {
		return errs.New(errs.CodeInvalidConfig,
			fmt.Sprintf("order notional %s exceeds the configured maximum %s", notional.FloatString(2), c.cfg.maxOrderNotional))
	}
	return nil
}

// orderNotional returns the notional value of o and whether it could be
// computed. AMOUNT orders use total_cash_amount; other orders use quantity
// times limit_price when both parse as decimals.
func orderNotional(o *OrderRequest) (*big.Rat, bool) {
	if o.EntrustType == EntrustTypeAmount {
		return parseDecimal(o.TotalCashAmount)
	}
	if strings.TrimSpace(o.Quantity) == "" || strings.TrimSpace(o.LimitPrice) == "" {
		return nil, false
	}
	qty, ok := parseDecimal(o.Quantity)
	if !ok {
		return nil, false
	}
	price, ok := parseDecimal(o.LimitPrice)
	if !ok {
		return nil, false
	}
	return new(big.Rat).Mul(qty, price), true
}

// parseDecimal parses a trimmed decimal string exactly. It reports false when s
// is not a number, including the empty string.
func parseDecimal(s string) (*big.Rat, bool) {
	r, ok := new(big.Rat).SetString(strings.TrimSpace(s))
	return r, ok
}

// isPositiveDecimal reports whether s is a decimal number greater than zero.
func isPositiveDecimal(s string) bool {
	r, ok := parseDecimal(s)
	return ok && r.Sign() > 0
}

// validClientOrderID reports whether id contains only the characters the API
// allows in a client order identifier: letters, digits, hyphen, and underscore.
func validClientOrderID(id string) bool {
	for i := 0; i < len(id); i++ {
		switch c := id[i]; {
		case c >= 'A' && c <= 'Z':
		case c >= 'a' && c <= 'z':
		case c >= '0' && c <= '9':
		case c == '-' || c == '_':
		default:
			return false
		}
	}
	return true
}

// valid reports whether t is a recognized instrument type.
func (t InstrumentType) valid() bool {
	switch t {
	case InstrumentTypeEquity, InstrumentTypeOption, InstrumentTypeFutures:
		return true
	default:
		return false
	}
}

// valid reports whether m is a recognized market.
func (m Market) valid() bool {
	switch m {
	case MarketUS, MarketHK, MarketCN:
		return true
	default:
		return false
	}
}

// valid reports whether s is a recognized order side.
func (s OrderSide) valid() bool {
	switch s {
	case OrderSideBuy, OrderSideSell, OrderSideShort:
		return true
	default:
		return false
	}
}

// valid reports whether t is a recognized order type.
func (t OrderType) valid() bool {
	switch t {
	case OrderTypeLimit, OrderTypeMarket, OrderTypeStopLoss, OrderTypeStopLossLimit,
		OrderTypeEnhancedLimit, OrderTypeAtAuction, OrderTypeAtAuctionLimit,
		OrderTypeMarketOnOpen, OrderTypeMarketOnClose, OrderTypeTrailingStopLoss,
		OrderTypeTrailingStopLossLimit, OrderTypeTouchMkt, OrderTypeTouchLmt:
		return true
	default:
		return false
	}
}

// needsLimitPrice reports whether t requires limit_price.
func (t OrderType) needsLimitPrice() bool {
	switch t {
	case OrderTypeLimit, OrderTypeStopLossLimit, OrderTypeTouchLmt:
		return true
	default:
		return false
	}
}

// needsStopPrice reports whether t requires stop_price.
func (t OrderType) needsStopPrice() bool {
	switch t {
	case OrderTypeStopLoss, OrderTypeStopLossLimit, OrderTypeTouchMkt, OrderTypeTouchLmt:
		return true
	default:
		return false
	}
}

// isTrailing reports whether t is a trailing stop order type.
func (t OrderType) isTrailing() bool {
	switch t {
	case OrderTypeTrailingStopLoss, OrderTypeTrailingStopLossLimit:
		return true
	default:
		return false
	}
}

// valid reports whether t is a recognized time-in-force.
func (t TimeInForce) valid() bool {
	switch t {
	case TimeInForceDay, TimeInForceGTD, TimeInForceGTC:
		return true
	default:
		return false
	}
}

// valid reports whether c is a recognized combo type.
func (c ComboType) valid() bool {
	switch c {
	case ComboTypeNormal, ComboTypeMaster, ComboTypeStopProfit, ComboTypeStopLoss,
		ComboTypeOTO, ComboTypeOCO, ComboTypeOTOCO:
		return true
	default:
		return false
	}
}

// comboTypeList renders the accepted combo types for error messages.
const comboTypeList = "NORMAL, MASTER, STOP_PROFIT, STOP_LOSS, OTO, OCO, OTOCO"

// valid reports whether s is a recognized trading session.
func (s TradingSession) valid() bool {
	switch s {
	case TradingSessionY, TradingSessionN, TradingSessionNight, TradingSessionAll,
		TradingSessionCore, TradingSessionAllDay:
		return true
	default:
		return false
	}
}

// valid reports whether t is a recognized trigger price type.
func (t TriggerPriceType) valid() bool {
	switch t {
	case TriggerPriceTypePrice, TriggerPriceTypePriceBid, TriggerPriceTypePriceAsk:
		return true
	default:
		return false
	}
}

// valid reports whether t is a recognized trailing type.
func (t TrailingType) valid() bool {
	switch t {
	case TrailingTypeAmount, TrailingTypePercentage:
		return true
	default:
		return false
	}
}
