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

package broker

import (
	"context"
	"net/url"
)

const (
	pathOrderPreview = "/openapi/v1/broker/orders/preview"
	pathOrderPlace   = "/openapi/v1/broker/orders/place"
	pathOrderReplace = "/openapi/v1/broker/orders/replace"
	pathOrderCancel  = "/openapi/v1/broker/orders/cancel"
	pathOrderDetail  = "/openapi/v1/broker/orders/detail"
	pathOrderHistory = "/openapi/v1/broker/orders/history"
	pathOpenOrders   = "/openapi/v1/broker/orders/open"
)

type BrokerOrder struct {
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
	CreateTime     string `json:"create_time"`
	UpdateTime     string `json:"update_time"`
	FilledQuantity string `json:"filled_quantity"`
	AvgFillPrice   string `json:"avg_fill_price,omitempty"`
}

type PreviewOrderRequest struct {
	AccountID   string `json:"account_id"`
	Symbol      string `json:"symbol"`
	OrderType   string `json:"order_type"`
	Side        string `json:"side"`
	Quantity    string `json:"quantity"`
	LimitPrice  string `json:"limit_price,omitempty"`
	StopPrice   string `json:"stop_price,omitempty"`
	TimeInForce string `json:"time_in_force"`
}

type OrderPreview struct {
	OrderID        string `json:"order_id"`
	EstimatedFee   string `json:"estimated_fee"`
	EstimatedTotal string `json:"estimated_total"`
	WarningMessage string `json:"warning_message,omitempty"`
}

func (c *Client) PreviewOrder(ctx context.Context, req PreviewOrderRequest) (*OrderPreview, error) {
	var out OrderPreview
	if err := c.post(ctx, pathOrderPreview, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) PlaceOrder(ctx context.Context, req PreviewOrderRequest) (*BrokerOrder, error) {
	var out BrokerOrder
	if err := c.post(ctx, pathOrderPlace, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type ReplaceOrderRequest struct {
	OrderID    string `json:"order_id"`
	LimitPrice string `json:"limit_price,omitempty"`
	StopPrice  string `json:"stop_price,omitempty"`
	Quantity   string `json:"quantity,omitempty"`
}

func (c *Client) ReplaceOrder(ctx context.Context, orderID string, req ReplaceOrderRequest) (*BrokerOrder, error) {
	path := pathOrderReplace + "?order_id=" + orderID
	var out BrokerOrder
	if err := c.put(ctx, path, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CancelOrder(ctx context.Context, orderID string) error {
	path := pathOrderCancel + "?order_id=" + orderID
	return c.delete(ctx, path, nil, nil, nil)
}

func (c *Client) GetOrderDetail(ctx context.Context, orderID string) (*BrokerOrder, error) {
	q := url.Values{}
	q.Set("order_id", orderID)
	var out BrokerOrder
	if err := c.get(ctx, pathOrderDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetOrderHistory(ctx context.Context, accountID string) ([]BrokerOrder, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out []BrokerOrder
	if err := c.get(ctx, pathOrderHistory, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetOpenOrders(ctx context.Context, accountID string) ([]BrokerOrder, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out []BrokerOrder
	if err := c.get(ctx, pathOpenOrders, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}
