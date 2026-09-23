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
	"net/http"
	"strings"

	"github.com/shing1211/webullapi4go/pkg/errors"
)

// Order modification endpoint paths.
const (
	// pathOrdersReplace modifies the terms of working orders.
	pathOrdersReplace = "/trading/orders/replace"
	// pathOrdersCancel cancels working orders.
	pathOrdersCancel = "/trading/orders/cancel"
)

// ModifyOrderRequest is a single order to modify within a
// [ReplaceOrderRequest]. Only the fields that are set are changed; every field
// is optional except ClientOrderID, which selects the order.
//
// The endpoint identifies an order by its client order identifier, so the
// caller must reuse the ClientOrderID supplied when the order was placed.
// Numeric values are strings to preserve precision.
type ModifyOrderRequest struct {
	// ClientOrderID selects the order to modify. It is required and must match
	// the identifier used when the order was placed: at most 32 characters of
	// [A-Za-z0-9_-].
	ClientOrderID string `json:"client_order_id"`
	// TimeInForce changes how long the order remains active.
	TimeInForce TimeInForce `json:"time_in_force,omitempty"`
	// Quantity changes the order quantity, as a decimal string.
	Quantity string `json:"quantity,omitempty"`
	// ExpireDate changes the GTD expiry in yyyy-MM-dd form. It is relevant only
	// when TimeInForce is GTD.
	ExpireDate string `json:"expire_date,omitempty"`
	// LimitPrice changes the limit price, as a decimal string.
	LimitPrice string `json:"limit_price,omitempty"`
	// StopPrice changes the trigger price, as a decimal string.
	StopPrice string `json:"stop_price,omitempty"`
	// TrailingType changes how TrailingStopStep is expressed.
	TrailingType TrailingType `json:"trailing_type,omitempty"`
	// TrailingStopStep changes the trailing spread, as a decimal string.
	TrailingStopStep string `json:"trailing_stop_step,omitempty"`
	// TrailingLimitPriceOffset changes the offset between the triggered stop
	// price and the submitted limit price, as a decimal string.
	TrailingLimitPriceOffset string `json:"trailing_limit_price_offset,omitempty"`
	// TriggerPriceType changes the market price a touch or stop order triggers
	// on.
	TriggerPriceType TriggerPriceType `json:"trigger_price_type,omitempty"`
}

// ReplaceOrderRequest is the body of [Client.ReplaceOrder]. It carries the
// account and one or more orders to modify.
type ReplaceOrderRequest struct {
	// AccountID is the account the orders belong to.
	AccountID string `json:"account_id"`
	// ModifyOrders are the orders to modify. At least one is required.
	ModifyOrders []ModifyOrderRequest `json:"modify_orders"`
}

// CancelOrderRequest is the body of [Client.CancelOrder]. The v3 cancel
// endpoint identifies an order by its client order identifier, so the caller
// must supply the ClientOrderID used when the order was placed.
type CancelOrderRequest struct {
	// AccountID is the account the order belongs to.
	AccountID string `json:"account_id"`
	// ClientOrderID selects the order to cancel: at most 32 characters of
	// [A-Za-z0-9_-].
	ClientOrderID string `json:"client_order_id"`
}

// ReplaceOrderResult is the response of [Client.ReplaceOrder].
type ReplaceOrderResult struct {
	// ClientOrderID echoes the caller-supplied order identifier.
	ClientOrderID string `json:"client_order_id"`
	// OrderID is the system-generated order identifier.
	OrderID string `json:"order_id"`
}

// CancelOrderResult is the response of [Client.CancelOrder].
type CancelOrderResult struct {
	// ClientOrderID echoes the caller-supplied order identifier.
	ClientOrderID string `json:"client_order_id"`
	// OrderID is the system-generated order identifier.
	OrderID string `json:"order_id"`
}

// Validate reports whether r is well formed: the account is set, at least one
// order is present, every order carries a valid client order identifier, and no
// two orders reuse one. It returns a typed [errs.Error] with
// [errs.CodeInvalidConfig] on the first problem found, so callers can reject a
// bad request before any network call.
func (r ReplaceOrderRequest) Validate() error {
	if strings.TrimSpace(r.AccountID) == "" {
		return errs.New(errs.CodeInvalidConfig, "account_id is required")
	}
	if len(r.ModifyOrders) == 0 {
		return errs.New(errs.CodeInvalidConfig, "modify_orders must contain at least one order")
	}
	seen := make(map[string]struct{}, len(r.ModifyOrders))
	for i := range r.ModifyOrders {
		order := &r.ModifyOrders[i]
		if err := order.validate(fmt.Sprintf("modify_orders[%d]: ", i)); err != nil {
			return err
		}
		if _, dup := seen[order.ClientOrderID]; dup {
			return errs.New(errs.CodeInvalidConfig,
				fmt.Sprintf("modify_orders[%d]: client_order_id %q is duplicated within the request", i, order.ClientOrderID))
		}
		seen[order.ClientOrderID] = struct{}{}
	}
	return nil
}

// validate checks a single [ModifyOrderRequest]. prefix, when non-empty, is
// prepended to the message to locate the order within a batch.
func (r ModifyOrderRequest) validate(prefix string) error {
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
	if r.TimeInForce != "" && !r.TimeInForce.valid() {
		return fail("time_in_force %q must be DAY, GTD or GTC", r.TimeInForce)
	}
	if r.TrailingType != "" && !r.TrailingType.valid() {
		return fail("trailing_type %q must be AMOUNT or PERCENTAGE", r.TrailingType)
	}
	if r.TriggerPriceType != "" && !r.TriggerPriceType.valid() {
		return fail("trigger_price_type %q must be PRICE, PRICE_BID or PRICE_ASK", r.TriggerPriceType)
	}
	return nil
}

// Validate reports whether r is well formed: the account and client order
// identifier are set and the identifier is within the documented length and
// character set. It returns a typed [errs.Error] with [errs.CodeInvalidConfig]
// on the first problem found, so callers can reject a bad request before any
// network call.
func (r CancelOrderRequest) Validate() error {
	fail := func(format string, args ...any) error {
		return errs.New(errs.CodeInvalidConfig, fmt.Sprintf(format, args...))
	}
	if strings.TrimSpace(r.AccountID) == "" {
		return fail("account_id is required")
	}
	switch {
	case strings.TrimSpace(r.ClientOrderID) == "":
		return fail("client_order_id is required")
	case len(r.ClientOrderID) > maxClientOrderIDLength:
		return fail("client_order_id must be at most %d characters, got %d", maxClientOrderIDLength, len(r.ClientOrderID))
	case !validClientOrderID(r.ClientOrderID):
		return fail("client_order_id %q may contain only letters, digits, '-' and '_'", r.ClientOrderID)
	}
	return nil
}

// ReplaceOrder modifies the working orders in req, matched by client order
// identifier, and returns the resulting identifiers. The request is validated
// before any network call, so a malformed request never reaches the API.
//
// Only the fields set on each [ModifyOrderRequest] are changed. ReplaceOrder
// can mutate live orders; callers should confirm the order identifiers first.
//
// Reference: https://developer.webull.hk/apis/docs/reference/common-order-replace.md
func (c *Client) ReplaceOrder(ctx context.Context, req ReplaceOrderRequest) (*ReplaceOrderResult, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var out ReplaceOrderResult
	if err := c.do(ctx, http.MethodPost, pathOrdersReplace, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelOrder cancels the order in req, matched by client order identifier, and
// returns the resulting identifiers. The request is validated before any
// network call, so a malformed request never reaches the API.
//
// CancelOrder can mutate live orders. The v3 endpoint identifies an order by
// [CancelOrderRequest.ClientOrderID]; the [CancelOrderResult.OrderID] is only
// echoed in the response.
//
// Reference: https://developer.webull.hk/apis/docs/reference/common-order-cancel.md
func (c *Client) CancelOrder(ctx context.Context, req CancelOrderRequest) (*CancelOrderResult, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var out CancelOrderResult
	if err := c.do(ctx, http.MethodPost, pathOrdersCancel, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
