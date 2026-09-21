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
	pathFDOrderPreview = "/broker-fd/orders/preview"
	pathFDOrderPlace   = "/broker-fd/orders/place"
	pathFDOrderReplace = "/broker-fd/orders/replace"
	pathFDOrderCancel  = "/broker-fd/orders/cancel"
	pathFDOrderDetail  = "/broker-fd/orders/detail"
	pathFDOrderHistory = "/broker-fd/orders/history"
	pathFDOrderOpen    = "/broker-fd/orders/open"
)

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

type FDOrderPreview struct {
	OrderID        string `json:"order_id"`
	EstimatedFee   string `json:"estimated_fee"`
	EstimatedTotal string `json:"estimated_total"`
}

func (c *Client) PreviewFDOrder(ctx context.Context, req FDOrderPreviewRequest) (*FDOrderPreview, error) {
	var out FDOrderPreview
	if err := c.post(ctx, pathFDOrderPreview, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) PlaceFDOrder(ctx context.Context, req FDOrderPreviewRequest) (*FDOrder, error) {
	var out FDOrder
	if err := c.post(ctx, pathFDOrderPlace, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type ReplaceFDOrderRequest struct {
	OrderID    string `json:"order_id"`
	LimitPrice string `json:"limit_price,omitempty"`
	StopPrice  string `json:"stop_price,omitempty"`
	Quantity   string `json:"quantity,omitempty"`
}

func (c *Client) ReplaceFDOrder(ctx context.Context, orderID string, req ReplaceFDOrderRequest) (*FDOrder, error) {
	q := url.Values{}
	q.Set("order_id", orderID)
	var out FDOrder
	if err := c.put(ctx, pathFDOrderReplace, q, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CancelFDOrder(ctx context.Context, orderID string) error {
	q := url.Values{}
	q.Set("order_id", orderID)
	return c.delete(ctx, pathFDOrderCancel, q, nil, nil)
}

func (c *Client) GetFDOrderDetail(ctx context.Context, orderID string) (*FDOrder, error) {
	q := url.Values{}
	q.Set("order_id", orderID)
	var out FDOrder
	if err := c.get(ctx, pathFDOrderDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetFDOrderHistory(ctx context.Context, accountID string) ([]FDOrder, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out []FDOrder
	if err := c.get(ctx, pathFDOrderHistory, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetFDOpenOrders(ctx context.Context, accountID string) ([]FDOrder, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out []FDOrder
	if err := c.get(ctx, pathFDOrderOpen, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}
