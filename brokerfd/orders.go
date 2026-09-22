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
	OrderID        string `json:"order_id"`
	AccountID      string `json:"account_id"`
	Symbol         string `json:"symbol"`
	OrderType      string `json:"order_type"`
	Side           string `json:"side"`
	Quantity       string `json:"quantity"`
	LimitPrice     string `json:"limit_price,omitempty"`
	StopPrice      string `json:"stop_price,omitempty"`
	TimeInForce    string `json:"time_in_force"`
	Status         string `json:"status"`
	FilledQuantity string `json:"filled_quantity"`
	AvgFillPrice   string `json:"avg_fill_price,omitempty"`
	CreateTime     string `json:"create_time"`
}

// FDOrderPreviewRequest contains the order parameters for a preview request.
type FDOrderPreviewRequest struct {
	AccountID   string `json:"account_id"`
	Symbol      string `json:"symbol"`
	OrderType   string `json:"order_type"`
	Side        string `json:"side"`
	Quantity    string `json:"quantity"`
	LimitPrice  string `json:"limit_price,omitempty"`
	StopPrice   string `json:"stop_price,omitempty"`
	TimeInForce string `json:"time_in_force"`
}

// FDOrderPreview contains the estimated cost breakdown for a fractional order.
type FDOrderPreview struct {
	OrderID        string `json:"order_id"`
	EstimatedFee   string `json:"estimated_fee"`
	EstimatedTotal string `json:"estimated_total"`
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
	OrderID    string `json:"order_id"`
	LimitPrice string `json:"limit_price,omitempty"`
	StopPrice  string `json:"stop_price,omitempty"`
	Quantity   string `json:"quantity,omitempty"`
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
func (c *Client) GetFDOrderHistory(ctx context.Context, accountID string) ([]FDOrder, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out []FDOrder
	if err := c.get(ctx, pathFDOrderHistory, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFDOpenOrders returns all open (unfilled) fractional orders for an account.
func (c *Client) GetFDOpenOrders(ctx context.Context, accountID string) ([]FDOrder, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out []FDOrder
	if err := c.get(ctx, pathFDOrderOpen, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}
