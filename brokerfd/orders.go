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

package brokerfd

import (
	"context"
	"net/url"

	"github.com/shing1211/webullapi4go/pkg/domain/money"
	"github.com/shing1211/webullapi4go/pkg/types"
)

const (
	pathFDOrderPreview = "/broker/orders/preview"
	pathFDOrderPlace   = "/broker/orders/place"
	pathFDOrderReplace = "/broker/orders/replace"
	pathFDOrderCancel  = "/broker/orders/cancel"
	pathFDOrderDetail  = "/broker/orders/get"
	pathFDOrderHistory = "/broker/orders/historical-orders/list"
	pathFDOrderOpen    = "/broker/orders/open-orders/list"
)

// FDOrder represents a fractional share order with execution details.
type FDOrder struct {
	OrderID        string      `json:"order_id"`
	AccountID      string      `json:"account_id"`
	Symbol         string      `json:"symbol"`
	OrderType      string      `json:"order_type"`
	Side           string      `json:"side"`
	Quantity       money.Money `json:"quantity"`
	LimitPrice     money.Money `json:"limit_price,omitempty"`
	StopPrice      money.Money `json:"stop_price,omitempty"`
	TimeInForce    string      `json:"time_in_force"`
	Status         string      `json:"status"`
	FilledQuantity money.Money `json:"filled_quantity"`
	AvgFillPrice   money.Money `json:"avg_fill_price,omitempty"`
	CreateTime     string      `json:"create_time"`

	// ClientRequestID is declared by the place and replace pages but not marked
	// required, so it was absent here and a response carrying it decoded the value to
	// the empty string with no error reported. That is weaker evidence than a
	// missing-required-name row: those pages publish no required list, so the name
	// may be optional or conditionally sent.
	//
	// It is the caller's own correlation key, which is what makes an order placed
	// through this SDK traceable to the request that placed it. The sibling broker/
	// module carries it in the same role. Not live-verified, for the reason given
	// on [TransferFee.FeeID].
	ClientRequestID string `json:"client_order_id"`
}

// FDOrderPreviewRequest contains the order parameters for a preview request.
type FDOrderPreviewRequest struct {
	AccountID   string       `json:"account_id"`
	Symbol      string       `json:"symbol"`
	OrderType   string       `json:"order_type"`
	Side        string       `json:"side"`
	Quantity    *money.Money `json:"quantity"`
	LimitPrice  *money.Money `json:"limit_price,omitempty"`
	StopPrice   *money.Money `json:"stop_price,omitempty"`
	TimeInForce string       `json:"time_in_force"`
}

// FDOrderPreview contains the estimated cost breakdown for a fractional order.
type FDOrderPreview struct {
	OrderID        string      `json:"order_id"`
	EstimatedFee   money.Money `json:"estimated_fee"`
	EstimatedTotal money.Money `json:"estimated_total"`

	// The two fields below are required by the page and were absent here, so they
	// decoded to zero with no error reported. The page splits the preview into a base
	// cost and a separate transaction fee; this type merges both into EstimatedTotal
	// and names only one of them EstimatedFee, so a caller cannot tell how much of the
	// quoted total is brokerage and how much is the underlying purchase. Not
	// live-verified, for the reason given on [TransferFee.FeeID].
	EstimatedCost           money.Money `json:"estimated_cost"`
	EstimatedTransactionFee money.Money `json:"estimated_transaction_fee"`
}

// PreviewFDOrder submits a fractional order for fee and cost estimation without execution.
func (c *Client) PreviewFDOrder(ctx context.Context, req FDOrderPreviewRequest) (*FDOrderPreview, error) {
	var out FDOrderPreview
	if err := c.post(ctx, pathFDOrderPreview, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PlaceFDOrder submits a fractional order for immediate execution.
func (c *Client) PlaceFDOrder(ctx context.Context, req FDOrderPreviewRequest) (*FDOrder, error) {
	var out FDOrder
	if err := c.post(ctx, pathFDOrderPlace, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ReplaceFDOrderRequest contains the fields that may be updated on an existing fractional order.
type ReplaceFDOrderRequest struct {
	OrderID    string       `json:"order_id"`
	LimitPrice *money.Money `json:"limit_price,omitempty"`
	StopPrice  *money.Money `json:"stop_price,omitempty"`
	Quantity   *money.Money `json:"quantity,omitempty"`
}

// ReplaceFDOrder modifies an existing fractional order with new parameters.
func (c *Client) ReplaceFDOrder(ctx context.Context, orderID string, req ReplaceFDOrderRequest) (*FDOrder, error) {
	q := url.Values{}
	q.Set("order_id", orderID)
	var out FDOrder
	if err := c.post(ctx, pathFDOrderReplace, q, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelFDOrder cancels a pending fractional order by ID.
func (c *Client) CancelFDOrder(ctx context.Context, orderID string) error {
	q := url.Values{}
	q.Set("order_id", orderID)
	return c.post(ctx, pathFDOrderCancel, q, nil, nil)
}

// GetFDOrderDetail retrieves the current state of a single fractional order.
func (c *Client) GetFDOrderDetail(ctx context.Context, orderID string) (*FDOrder, error) {
	q := url.Values{}
	q.Set("order_id", orderID)
	var out FDOrder
	if err := c.get(ctx, pathFDOrderDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetFDOrderHistory returns all filled and cancelled fractional orders for an account.
func (c *Client) GetFDOrderHistory(ctx context.Context, accountID, paginationKey string) (*types.Page[FDOrder], error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	if paginationKey != "" {
		q.Set("pagination_key", paginationKey)
	}
	var out types.Page[FDOrder]
	if err := c.get(ctx, pathFDOrderHistory, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetFDOpenOrders returns all open (unfilled) fractional orders for an account.
func (c *Client) GetFDOpenOrders(ctx context.Context, accountID, paginationKey string) (*types.Page[FDOrder], error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	if paginationKey != "" {
		q.Set("pagination_key", paginationKey)
	}
	var out types.Page[FDOrder]
	if err := c.get(ctx, pathFDOrderOpen, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
