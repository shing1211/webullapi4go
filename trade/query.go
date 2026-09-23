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
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/shing1211/webullapi4go/pkg/errors"
)

// Order query endpoint paths.
const (
	// pathOrdersOpen lists the currently working orders of an account.
	pathOrdersOpen = "/trading/orders/open-orders/list"
	// pathOrdersHistory lists an account's orders over a time range.
	pathOrdersHistory = "/trading/orders/historical-orders/list"
	// pathOrdersDetail returns a single order by client order identifier.
	pathOrdersDetail = "/trading/orders/get"
)

// MaxOrderQueryPages bounds how many cursor pages [Client.GetAllOpenOrders] and
// [Client.GetAllOrderHistory] fetch in a single call. The API returns a
// pagination_key on every page but the last, so a misbehaving server could
// otherwise page forever; after this many requests the call fails with an
// [errs.CodeAPI] error instead. It is exported so callers can reason about the
// bound.
const MaxOrderQueryPages = 100

// OrderPage is a single page of grouped orders returned by the open-order and
// order-history endpoints. Orders holds the page's groups; PaginationKey, when
// non-empty, is the cursor to pass to the next request. An empty PaginationKey
// means the page is the last one.
type OrderPage struct {
	// Orders is the page of order groups.
	Orders []OrderGroup `json:"data"`
	// PaginationKey is the cursor for the next page, empty on the last page.
	PaginationKey string `json:"pagination_key"`
}

// OrderGroup is a client order together with its child orders. A NORMAL order
// has a single entry in Orders; combo orders (OTO, OCO, OTOCO, take-profit and
// stop-loss sets) carry their legs here.
type OrderGroup struct {
	// ClientOrderID is the caller-supplied identifier grouping the orders.
	ClientOrderID string `json:"client_order_id"`
	// ComboType is the role the group plays within a combo order.
	ComboType ComboType `json:"combo_type"`
	// Orders are the individual orders in the group.
	Orders []Order `json:"orders"`
}

// Order is the state of a single order as returned by the order query
// endpoints. Numeric values are decimal strings to preserve precision; a field
// the API omits for a given order is left empty.
type Order struct {
	// ClientOrderID is the caller-supplied order identifier.
	ClientOrderID string `json:"client_order_id"`
	// OrderID is the system-generated order identifier.
	OrderID string `json:"order_id"`
	// Symbol is the trading symbol of the instrument.
	Symbol string `json:"symbol"`
	// Side is the intended trading direction.
	Side OrderSide `json:"side"`
	// Status is the current lifecycle state of the order.
	Status OrderStatus `json:"status"`
	// OrderType is the execution instruction.
	OrderType OrderType `json:"order_type"`
	// InstrumentType is the kind of instrument traded.
	InstrumentType InstrumentType `json:"instrument_type"`
	// SupportTradingSession is the US trading session the order is restricted
	// to; it is empty when the API applies its default.
	SupportTradingSession TradingSession `json:"support_trading_session"`
	// TimeInForce is how long the order remains active.
	TimeInForce TimeInForce `json:"time_in_force"`
	// TotalQuantity is the total order quantity, as a decimal string.
	TotalQuantity string `json:"total_quantity"`
	// FilledQuantity is the quantity executed so far, as a decimal string.
	FilledQuantity string `json:"filled_quantity"`
	// FilledPrice is the average execution price, as a decimal string. It may
	// be zero or empty when nothing has filled.
	FilledPrice string `json:"filled_price"`
	// LimitPrice is the limit price, as a decimal string.
	LimitPrice string `json:"limit_price"`
	// StopPrice is the stop trigger price, as a decimal string.
	StopPrice string `json:"stop_price"`
	// TrailingType is how TrailingStopStep is expressed.
	TrailingType TrailingType `json:"trailing_type"`
	// TrailingStopStep is the trailing spread, as a decimal string.
	TrailingStopStep string `json:"trailing_stop_step"`
	// TrailingLimitPriceOffset is the offset between the triggered stop price
	// and the submitted limit price, as a decimal string.
	TrailingLimitPriceOffset string `json:"trailing_limit_price_offset"`
	// TriggerPriceType is the market price a touch or stop order triggers on.
	TriggerPriceType TriggerPriceType `json:"trigger_price_type"`
	// PlaceTimeAt is the order placement time in ISO8601 UTC form, for example
	// "2025-11-11T05:44:35.385Z".
	PlaceTimeAt string `json:"place_time_at"`
	// FilledTimeAt is the time of the last execution in ISO8601 UTC form.
	FilledTimeAt string `json:"filled_time_at"`
	// Legs lists the option legs. It is empty for non-option orders.
	Legs []OrderLegDetail `json:"legs"`
	// Commission is the commission breakdown. It is only present on the
	// order-detail response.
	Commission *OrderCommission `json:"commission,omitempty"`
	// Fees is the fee breakdown. It is only present on the order-detail
	// response.
	Fees []OrderFee `json:"fees,omitempty"`
}

// OrderLegDetail is one leg of an option order as returned by the order query
// endpoints. It is the response counterpart of the request-side [OrderLeg].
type OrderLegDetail struct {
	// Symbol is the trading symbol of the leg.
	Symbol string `json:"symbol"`
	// Side is the intended direction of the leg.
	Side OrderSide `json:"side"`
	// Quantity is the leg quantity, as a decimal string.
	Quantity string `json:"quantity"`
	// OptionType is whether the leg is a call or a put.
	OptionType OptionType `json:"option_type"`
	// OptionCategory is the option's exercise style, for example AMERICAN.
	OptionCategory string `json:"option_category"`
	// OptionStrategy identifies the option strategy, for example SINGLE.
	OptionStrategy OptionStrategy `json:"option_strategy"`
	// StrikePrice is the option strike, as a decimal string.
	StrikePrice string `json:"strike_price"`
	// OptionContractMultiplier is the number of shares one contract
	// represents, as a decimal string.
	OptionContractMultiplier string `json:"option_contract_multiplier"`
	// OptionContractDeliverable is the number of shares delivered on exercise
	// of one contract, as a decimal string.
	OptionContractDeliverable string `json:"option_contract_deliverable"`
	// OptionExpireDate is the option expiry in yyyy-MM-dd form.
	OptionExpireDate string `json:"option_expire_date"`
}

// OrderCommission is the commission breakdown of a filled order.
type OrderCommission struct {
	// ActualCommission is the commission collected, as a decimal string.
	ActualCommission string `json:"actual_commission"`
	// ReceivableCommission is the commission still receivable, as a decimal
	// string.
	ReceivableCommission string `json:"receivable_commission"`
}

// OrderFee is a single fee line of a filled order.
type OrderFee struct {
	// Type is the fee type, for example "FINRA_CAT_REGULATORY_FEE".
	Type string `json:"type"`
	// ActualValue is the fee collected, as a decimal string.
	ActualValue string `json:"actual_value"`
	// ReceivableValue is the fee still receivable, as a decimal string.
	ReceivableValue string `json:"receivable_value"`
}

// OrderHistoryQuery parameterizes [Client.GetOrderHistory] and
// [Client.GetAllOrderHistory]. AccountID is required; the remaining fields are
// optional and omitted from the request when empty.
type OrderHistoryQuery struct {
	// AccountID is the account whose history is queried. Required.
	AccountID string
	// StartTime is the inclusive start of the range in RFC3339 UTC form, for
	// example "2025-01-05T22:59:59.012Z". When empty the API defaults to the
	// last seven days; the API allows at most six months of look-back.
	StartTime string
	// EndTime is the exclusive end of the range in RFC3339 UTC form. When empty
	// the API defaults to now.
	EndTime string
	// PaginationKey resumes from a previous page. It is normally left empty and
	// set from an [OrderPage.PaginationKey] for manual paging.
	PaginationKey string
}

// GetOpenOrders returns the working orders of accountID from the first page of
// the open-order feed. The result is empty when the account has no working
// orders, and groups a combo order's legs into a single [OrderGroup].
//
// Only the first page is returned; use [Client.GetAllOpenOrders] to follow the
// pagination cursor to exhaustion, or [Client.GetOpenOrdersPage] to page
// manually. accountID is required and is validated with an
// [errs.CodeInvalidConfig] error before any network call.
//
// Reference: https://developer.webull.hk/apis/docs/reference/order-open.md
func (c *Client) GetOpenOrders(ctx context.Context, accountID string) ([]OrderGroup, error) {
	page, err := c.GetOpenOrdersPage(ctx, accountID, "")
	if err != nil {
		return nil, err
	}
	return page.Orders, nil
}

// GetOpenOrdersPage returns one page of the open orders of accountID.
// paginationKey, when non-empty, resumes from a previous page and is available
// as [OrderPage.PaginationKey]. The returned page preserves the cursor so
// callers can page explicitly.
//
// accountID is required and is validated with an [errs.CodeInvalidConfig] error
// before any network call.
//
// Reference: https://developer.webull.hk/apis/docs/reference/order-open.md
func (c *Client) GetOpenOrdersPage(ctx context.Context, accountID, paginationKey string) (*OrderPage, error) {
	query, err := accountQuery(accountID)
	if err != nil {
		return nil, err
	}
	if paginationKey != "" {
		query.Set("pagination_key", paginationKey)
	}
	return c.queryOrderPage(ctx, pathOrdersOpen, query)
}

// GetAllOpenOrders returns every working order of accountID by following the
// pagination cursor to exhaustion. The result is the concatenation of every
// page's groups.
//
// The walk is bounded by [MaxOrderQueryPages]; a server that keeps returning a
// non-empty pagination_key fails the call with an [errs.CodeAPI] error rather
// than looping forever. accountID is required and is validated with an
// [errs.CodeInvalidConfig] error before any network call.
//
// Reference: https://developer.webull.hk/apis/docs/reference/order-open.md
func (c *Client) GetAllOpenOrders(ctx context.Context, accountID string) ([]OrderGroup, error) {
	query, err := accountQuery(accountID)
	if err != nil {
		return nil, err
	}
	return c.collectOrderGroups(ctx, pathOrdersOpen, query)
}

// GetOrderHistory returns the orders of q.AccountID from the first page of the
// history feed. When q.StartTime and q.EndTime are empty the API defaults to
// the last seven days. Only the first page is returned; use
// [Client.GetAllOrderHistory] to follow the cursor, or
// [Client.GetOrderHistoryPage] to page manually.
//
// q.AccountID is required, and a non-empty time must be an RFC3339 timestamp;
// either violation is reported as an [errs.CodeInvalidConfig] error before any
// network call.
//
// Reference: https://developer.webull.hk/apis/docs/reference/order-history.md
func (c *Client) GetOrderHistory(ctx context.Context, q OrderHistoryQuery) ([]OrderGroup, error) {
	page, err := c.GetOrderHistoryPage(ctx, q)
	if err != nil {
		return nil, err
	}
	return page.Orders, nil
}

// GetOrderHistoryPage returns one page of the order history described by q.
// q.PaginationKey, when non-empty, resumes from a previous page. The returned
// page preserves the cursor so callers can page explicitly.
//
// q.AccountID is required, and a non-empty time must be an RFC3339 timestamp;
// either violation is reported as an [errs.CodeInvalidConfig] error before any
// network call.
//
// Reference: https://developer.webull.hk/apis/docs/reference/order-history.md
func (c *Client) GetOrderHistoryPage(ctx context.Context, q OrderHistoryQuery) (*OrderPage, error) {
	query, err := q.values()
	if err != nil {
		return nil, err
	}
	return c.queryOrderPage(ctx, pathOrdersHistory, query)
}

// GetAllOrderHistory returns every order described by q by following the
// pagination cursor to exhaustion. The walk is bounded by
// [MaxOrderQueryPages]; a server that keeps returning a non-empty
// pagination_key fails the call with an [errs.CodeAPI] error rather than
// looping forever.
//
// q.AccountID is required, and a non-empty time must be an RFC3339 timestamp;
// either violation is reported as an [errs.CodeInvalidConfig] error before any
// network call.
//
// Reference: https://developer.webull.hk/apis/docs/reference/order-history.md
func (c *Client) GetAllOrderHistory(ctx context.Context, q OrderHistoryQuery) ([]OrderGroup, error) {
	query, err := q.values()
	if err != nil {
		return nil, err
	}
	return c.collectOrderGroups(ctx, pathOrdersHistory, query)
}

// GetOrderDetail returns a single order group by its client order identifier.
// The response carries the same order fields as the list endpoints plus the
// commission and fee breakdowns.
//
// The reference endpoint identifies the order by client_order_id, so
// clientOrderID is the caller-supplied identifier returned by
// [Client.PlaceOrder]. Both accountID and clientOrderID are required and are
// validated with an [errs.CodeInvalidConfig] error before any network call.
//
// Reference: https://developer.webull.hk/apis/docs/reference/order-detail.md
func (c *Client) GetOrderDetail(ctx context.Context, accountID, clientOrderID string) (*OrderGroup, error) {
	query, err := accountQuery(accountID)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(clientOrderID) == "" {
		return nil, errs.New(errs.CodeInvalidConfig, "client_order_id is required")
	}
	query.Set("client_order_id", clientOrderID)

	var out OrderGroup
	if err := c.get(ctx, pathOrdersDetail, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// queryOrderPage performs a single paginated order request and decodes the page.
func (c *Client) queryOrderPage(ctx context.Context, path string, query url.Values) (*OrderPage, error) {
	var out OrderPage
	if err := c.get(ctx, path, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// collectOrderGroups follows the pagination_key of path until the server stops
// returning one, concatenating every page's groups. query is mutated with the
// cursor as the walk advances. It returns an [errs.CodeAPI] error when the walk
// exceeds [MaxOrderQueryPages].
func (c *Client) collectOrderGroups(ctx context.Context, path string, query url.Values) ([]OrderGroup, error) {
	var (
		all []OrderGroup
		key string
	)
	for page := 1; ; page++ {
		if page > MaxOrderQueryPages {
			return nil, errs.New(errs.CodeAPI,
				fmt.Sprintf("order query exceeded the maximum of %d pages", MaxOrderQueryPages))
		}
		if key != "" {
			query.Set("pagination_key", key)
		}
		res, err := c.queryOrderPage(ctx, path, query)
		if err != nil {
			return nil, err
		}
		all = append(all, res.Orders...)
		if res.PaginationKey == "" {
			return all, nil
		}
		key = res.PaginationKey
	}
}

// values validates q and renders it as the query string of the history
// endpoint. It returns an [errs.CodeInvalidConfig] error when AccountID is
// empty or a time is not RFC3339, so callers never send a malformed request.
func (q OrderHistoryQuery) values() (url.Values, error) {
	query, err := accountQuery(q.AccountID)
	if err != nil {
		return nil, err
	}
	if err := validateOrderTime("start_time", q.StartTime); err != nil {
		return nil, err
	}
	if err := validateOrderTime("end_time", q.EndTime); err != nil {
		return nil, err
	}
	if q.StartTime != "" {
		query.Set("start_time", q.StartTime)
	}
	if q.EndTime != "" {
		query.Set("end_time", q.EndTime)
	}
	if q.PaginationKey != "" {
		query.Set("pagination_key", q.PaginationKey)
	}
	return query, nil
}

// validateOrderTime reports whether value, when non-empty, is an RFC3339
// timestamp such as "2025-01-05T22:59:59.012Z". name locates the field in the
// returned error.
func validateOrderTime(name, value string) error {
	if value == "" {
		return nil
	}
	if _, err := time.Parse(time.RFC3339, value); err != nil {
		return errs.Wrap(errs.CodeInvalidConfig,
			fmt.Sprintf("%s %q must be an RFC3339 timestamp such as 2025-01-05T22:59:59.012Z", name, value), err)
	}
	return nil
}
