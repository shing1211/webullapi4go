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
	pathOrderPreview = "/broker/orders/preview"
	pathOrderPlace   = "/broker/orders/place"
	pathOrderReplace = "/broker/orders/replace"
	pathOrderCancel  = "/broker/orders/cancel"
	pathOrderDetail  = "/broker/orders/get"
	pathOrderHistory = "/broker/orders/historical-orders/list"
	pathOpenOrders   = "/broker/orders/open-orders/list"
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
	OrderID            string `json:"order_id"`
	EstimatedCost      string `json:"estimated_cost"`
	EstimatedFee       string `json:"estimated_transaction_fee"`
	EstimatedFeeDetail any    `json:"estimated_transaction_fee_detail,omitempty"`
	WarningMessage     string `json:"warning_message,omitempty"`
}

type previewOrderInternal struct {
	AccountID          string              `json:"account_id"`
	ClientComboOrderID string              `json:"client_combo_order_id,omitempty"`
	NewOrders          []orderLineInternal `json:"new_orders"`
}

type orderLineInternal struct {
	ComboType             string `json:"combo_type"`
	ClientOrderID         string `json:"client_order_id"`
	InstrumentType        string `json:"instrument_type"`
	Market                string `json:"market"`
	Symbol                string `json:"symbol"`
	OrderType             string `json:"order_type"`
	EntrustType           string `json:"entrust_type"`
	SupportTradingSession string `json:"support_trading_session"`
	TimeInForce           string `json:"time_in_force"`
	Side                  string `json:"side"`
	Quantity              string `json:"quantity"`
	LimitPrice            string `json:"limit_price,omitempty"`
	StopPrice             string `json:"stop_price,omitempty"`
}

func (c *Client) PreviewOrder(ctx context.Context, req PreviewOrderRequest) (*OrderPreview, error) {
	body := previewOrderInternal{
		AccountID: req.AccountID,
		NewOrders: []orderLineInternal{{
			ComboType:             "SINGLE",
			ClientOrderID:         "",
			InstrumentType:        "EQUITY",
			Market:                "",
			Symbol:                req.Symbol,
			OrderType:             req.OrderType,
			EntrustType:           "QTY",
			SupportTradingSession: "ALL",
			TimeInForce:           req.TimeInForce,
			Side:                  req.Side,
			Quantity:              req.Quantity,
			LimitPrice:            req.LimitPrice,
			StopPrice:             req.StopPrice,
		}},
	}
	var out OrderPreview
	if err := c.post(ctx, pathOrderPreview, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) PlaceOrder(ctx context.Context, req PreviewOrderRequest) (*BrokerOrder, error) {
	body := previewOrderInternal{
		AccountID: req.AccountID,
		NewOrders: []orderLineInternal{{
			ComboType:             "SINGLE",
			ClientOrderID:         "",
			InstrumentType:        "EQUITY",
			Market:                "",
			Symbol:                req.Symbol,
			OrderType:             req.OrderType,
			EntrustType:           "QTY",
			SupportTradingSession: "ALL",
			TimeInForce:           req.TimeInForce,
			Side:                  req.Side,
			Quantity:              req.Quantity,
			LimitPrice:            req.LimitPrice,
			StopPrice:             req.StopPrice,
		}},
	}
	var out BrokerOrder
	if err := c.post(ctx, pathOrderPlace, nil, body, &out); err != nil {
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
	req.OrderID = orderID
	var out BrokerOrder
	if err := c.post(ctx, pathOrderReplace, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type cancelOrderRequest struct {
	OrderID string `json:"order_id"`
}

func (c *Client) CancelOrder(ctx context.Context, orderID string) error {
	return c.post(ctx, pathOrderCancel, nil, cancelOrderRequest{OrderID: orderID}, nil)
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
