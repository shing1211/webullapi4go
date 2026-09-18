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

	"github.com/shing1211/webullapi4go/internal/errs"
)

// Combo orders are a US-equity-only feature. These are the documented per-role
// order-type matrices:
//
//   - a take-profit/stop-loss MASTER is a MARKET or LIMIT order;
//   - an OTO or OTOCO MASTER is a MARKET, LIMIT, STOP_LOSS, or
//     STOP_LOSS_LIMIT order;
//   - an OTO follow-up order is a MARKET, LIMIT, STOP_LOSS, or
//     STOP_LOSS_LIMIT order;
//   - an OCO or OTOCO leg is a LIMIT, STOP_LOSS, or STOP_LOSS_LIMIT order.
var (
	// comboTPSLMasterOrderTypes are the order types allowed for the MASTER of
	// a take-profit/stop-loss group.
	comboTPSLMasterOrderTypes = []OrderType{OrderTypeMarket, OrderTypeLimit}
	// comboOTOMasterOrderTypes are the order types allowed for the MASTER of
	// an OTO or OTOCO group.
	comboOTOMasterOrderTypes = []OrderType{
		OrderTypeMarket, OrderTypeLimit, OrderTypeStopLoss, OrderTypeStopLossLimit,
	}
	// comboOTOlegOrderTypes are the order types allowed for an OTO follow-up
	// order.
	comboOTOlegOrderTypes = []OrderType{
		OrderTypeMarket, OrderTypeLimit, OrderTypeStopLoss, OrderTypeStopLossLimit,
	}
	// comboOCOlegOrderTypes are the order types allowed for an OCO order.
	comboOCOlegOrderTypes = []OrderType{
		OrderTypeLimit, OrderTypeStopLoss, OrderTypeStopLossLimit,
	}
	// comboOTOCOlegOrderTypes are the order types allowed for an OTOCO leg.
	comboOTOCOlegOrderTypes = []OrderType{
		OrderTypeLimit, OrderTypeStopLoss, OrderTypeStopLossLimit,
	}
)

// Combo group leg-count bounds from the API combo-order table. A take-profit or
// stop-loss sub-order is optional and is bounded separately by
// [PlaceOrderRequest.validateTakeProfitStopLoss].
const (
	// comboOTOminLegs is the minimum number of OTO follow-up orders.
	comboOTOminLegs = 1
	// comboOTOMaxLegs is the maximum number of OTO follow-up orders.
	comboOTOMaxLegs = 6
	// comboOCOminLegs is the minimum number of OCO orders.
	comboOCOminLegs = 2
	// comboOCOmaxLegs is the maximum number of OCO orders.
	comboOCOmaxLegs = 6
	// comboOTOCOminLegs is the minimum number of OTOCO legs.
	comboOTOCOminLegs = 1
	// comboOTOCOmaxLegs is the maximum number of OTOCO legs.
	comboOTOCOmaxLegs = 6
)

// validateComboRules enforces the combo-order group composition documented for
// US equity orders before any network call. It complements the per-order
// validation in [OrderRequest.Validate], which still checks each order's
// required fields and market rules, by ensuring the order set forms a valid
// MASTER/OTO/OCO/OTOCO or take-profit/stop-loss group.
//
// A request whose orders are all NORMAL is left unchanged. Otherwise the whole
// request must be a single combo group and must satisfy:
//
//   - every order uses a non-NORMAL combo_type;
//   - every order is a US equity order, because combo orders are equity-only
//     and US-only (option and futures orders support only NORMAL);
//   - client_combo_order_id is non-empty;
//   - the MASTER and leg counts, the per-role order types, and the
//     sell-to-close side rule match the documented table for the group's kind.
//
// It returns a typed [errs.Error] with [errs.CodeInvalidConfig] on the first
// problem found.
func (r PlaceOrderRequest) validateComboRules() error {
	counts := make(map[ComboType]int, len(r.NewOrders))
	hasCombo := false
	for i := range r.NewOrders {
		ct := r.NewOrders[i].ComboType
		counts[ct]++
		if ct != ComboTypeNormal {
			hasCombo = true
		}
	}
	if !hasCombo {
		return nil
	}

	if strings.TrimSpace(r.ClientComboOrderID) == "" {
		return comboFail("client_combo_order_id is required when any order has a combo_type other than NORMAL")
	}

	for i := range r.NewOrders {
		o := &r.NewOrders[i]
		switch {
		case o.ComboType == ComboTypeNormal:
			return comboFail("new_orders[%d]: NORMAL orders cannot be mixed with combo orders", i)
		case o.InstrumentType != InstrumentTypeEquity:
			return comboFail("new_orders[%d]: combo_type %s is only supported for EQUITY orders, got instrument_type %s",
				i, o.ComboType, o.InstrumentType)
		case o.Market != MarketUS:
			return comboFail("new_orders[%d]: combo_type %s is only supported for US orders, got market %s",
				i, o.ComboType, o.Market)
		}
	}

	switch {
	case counts[ComboTypeOTO] > 0:
		return r.validateOTO(counts)
	case counts[ComboTypeOTOCO] > 0:
		return r.validateOTOCO(counts)
	case counts[ComboTypeOCO] > 0:
		return r.validateOCO(counts)
	default:
		return r.validateTakeProfitStopLoss(counts)
	}
}

// validateTakeProfitStopLoss validates a take-profit/stop-loss group. It
// accepts either the buy-to-open form (exactly one MASTER with 0-1 STOP_PROFIT
// and 0-1 STOP_LOSS) or the sell-to-close form (only STOP_PROFIT/STOP_LOSS
// orders, side SELL, and no MASTER).
func (r PlaceOrderRequest) validateTakeProfitStopLoss(counts map[ComboType]int) error {
	if counts[ComboTypeMaster] > 1 {
		return comboFail("new_orders: take-profit/stop-loss combo allows at most one MASTER order, got %d",
			counts[ComboTypeMaster])
	}
	if counts[ComboTypeStopProfit] > 1 {
		return comboFail("new_orders: take-profit/stop-loss combo allows at most one STOP_PROFIT order, got %d",
			counts[ComboTypeStopProfit])
	}
	if counts[ComboTypeStopLoss] > 1 {
		return comboFail("new_orders: take-profit/stop-loss combo allows at most one STOP_LOSS order, got %d",
			counts[ComboTypeStopLoss])
	}

	sellToClose := counts[ComboTypeMaster] == 0
	for i := range r.NewOrders {
		o := &r.NewOrders[i]
		switch o.ComboType {
		case ComboTypeMaster:
			if !containsOrderType(comboTPSLMasterOrderTypes, o.OrderType) {
				return comboOrderTypeFail(i, o.ComboType, o.OrderType, comboTPSLMasterOrderTypes)
			}
		case ComboTypeStopProfit:
			if o.OrderType != OrderTypeLimit {
				return comboOrderTypeFail(i, o.ComboType, o.OrderType, []OrderType{OrderTypeLimit})
			}
		case ComboTypeStopLoss:
			if o.OrderType != OrderTypeStopLoss {
				return comboOrderTypeFail(i, o.ComboType, o.OrderType, []OrderType{OrderTypeStopLoss})
			}
		}
		if sellToClose && o.Side != OrderSideSell {
			return comboFail("new_orders[%d]: %s orders in a sell-to-close take-profit/stop-loss combo must use side SELL, got %s",
				i, o.ComboType, o.Side)
		}
	}
	return nil
}

// validateOTO validates a one-triggers-the-other group: exactly one MASTER plus
// between one and six OTO follow-up orders.
func (r PlaceOrderRequest) validateOTO(counts map[ComboType]int) error {
	if err := r.checkComboGroup("OTO", ComboTypeMaster, ComboTypeOTO); err != nil {
		return err
	}
	if counts[ComboTypeMaster] != 1 {
		return comboFail("new_orders: OTO combo requires exactly one MASTER order, got %d",
			counts[ComboTypeMaster])
	}
	if n := counts[ComboTypeOTO]; n < comboOTOminLegs || n > comboOTOMaxLegs {
		return comboFail("new_orders: OTO combo requires between %d and %d OTO orders, got %d",
			comboOTOminLegs, comboOTOMaxLegs, n)
	}
	for i := range r.NewOrders {
		o := &r.NewOrders[i]
		allowed := comboOTOlegOrderTypes
		if o.ComboType == ComboTypeMaster {
			allowed = comboOTOMasterOrderTypes
		}
		if !containsOrderType(allowed, o.OrderType) {
			return comboOrderTypeFail(i, o.ComboType, o.OrderType, allowed)
		}
	}
	return nil
}

// validateOCO validates a one-cancels-the-other group: between two and six OCO
// orders with no MASTER.
func (r PlaceOrderRequest) validateOCO(counts map[ComboType]int) error {
	if err := r.checkComboGroup("OCO", ComboTypeOCO); err != nil {
		return err
	}
	if n := counts[ComboTypeOCO]; n < comboOCOminLegs || n > comboOCOmaxLegs {
		return comboFail("new_orders: OCO combo requires between %d and %d OCO orders, got %d",
			comboOCOminLegs, comboOCOmaxLegs, n)
	}
	for i := range r.NewOrders {
		o := &r.NewOrders[i]
		if !containsOrderType(comboOCOlegOrderTypes, o.OrderType) {
			return comboOrderTypeFail(i, o.ComboType, o.OrderType, comboOCOlegOrderTypes)
		}
	}
	return nil
}

// validateOTOCO validates a one-triggers-one-cancels-the-other group: exactly
// one MASTER plus between one and six OTOCO legs.
func (r PlaceOrderRequest) validateOTOCO(counts map[ComboType]int) error {
	if err := r.checkComboGroup("OTOCO", ComboTypeMaster, ComboTypeOTOCO); err != nil {
		return err
	}
	if counts[ComboTypeMaster] != 1 {
		return comboFail("new_orders: OTOCO combo requires exactly one MASTER order, got %d",
			counts[ComboTypeMaster])
	}
	if n := counts[ComboTypeOTOCO]; n < comboOTOCOminLegs || n > comboOTOCOmaxLegs {
		return comboFail("new_orders: OTOCO combo requires between %d and %d OTOCO orders, got %d",
			comboOTOCOminLegs, comboOTOCOmaxLegs, n)
	}
	for i := range r.NewOrders {
		o := &r.NewOrders[i]
		allowed := comboOTOCOlegOrderTypes
		if o.ComboType == ComboTypeMaster {
			allowed = comboOTOMasterOrderTypes
		}
		if !containsOrderType(allowed, o.OrderType) {
			return comboOrderTypeFail(i, o.ComboType, o.OrderType, allowed)
		}
	}
	return nil
}

// checkComboGroup reports the first order whose combo_type is not part of the
// named group's allowed roles. It rejects mixture such as a MASTER inside an OCO
// group or a take-profit leg inside an OTO group.
func (r PlaceOrderRequest) checkComboGroup(kind string, allowed ...ComboType) error {
	for i := range r.NewOrders {
		if !containsComboType(allowed, r.NewOrders[i].ComboType) {
			return comboFail("new_orders[%d]: combo_type %s cannot be combined in a %s group",
				i, r.NewOrders[i].ComboType, kind)
		}
	}
	return nil
}

// comboFail builds a typed [errs.Error] with [errs.CodeInvalidConfig] for a
// combo composition problem.
func comboFail(format string, args ...any) error {
	return errs.New(errs.CodeInvalidConfig, fmt.Sprintf(format, args...))
}

// comboOrderTypeFail builds a typed [errs.Error] for an order type that is not
// allowed for a combo role, listing the supported types.
func comboOrderTypeFail(index int, role ComboType, ot OrderType, allowed []OrderType) error {
	return comboFail("new_orders[%d]: order_type %s is not supported for %s orders; supported types: %s",
		index, ot, role, orderTypeList(allowed))
}

// containsOrderType reports whether allowed contains t.
func containsOrderType(allowed []OrderType, t OrderType) bool {
	for _, a := range allowed {
		if a == t {
			return true
		}
	}
	return false
}

// containsComboType reports whether allowed contains c.
func containsComboType(allowed []ComboType, c ComboType) bool {
	for _, a := range allowed {
		if a == c {
			return true
		}
	}
	return false
}
