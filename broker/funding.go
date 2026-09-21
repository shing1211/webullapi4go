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
	pathFXRate               = "/openapi/v1/broker/funding/fx/rate"
	pathFXExchange           = "/openapi/v1/broker/funding/fx/exchange"
	pathFXExchangeDetail     = "/openapi/v1/broker/funding/fx/exchange/detail"
	pathInstantExchange      = "/openapi/v1/broker/funding/fx/instant-exchange"
	pathInstantFunding       = "/openapi/v1/broker/funding/instant"
	pathInstantFundingDetail = "/openapi/v1/broker/funding/instant/detail"
)

type FXRate struct {
	FromCurrency string `json:"from_currency"`
	ToCurrency   string `json:"to_currency"`
	Rate         string `json:"rate"`
	Timestamp    string `json:"timestamp"`
}

func (c *Client) GetFXRate(ctx context.Context, from, to string) (*FXRate, error) {
	q := url.Values{}
	q.Set("from", from)
	q.Set("to", to)
	var out FXRate
	if err := c.get(ctx, pathFXRate, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type CreateFXExchangeRequest struct {
	AccountID    string `json:"account_id"`
	FromCurrency string `json:"from_currency"`
	ToCurrency   string `json:"to_currency"`
	Amount       string `json:"amount"`
}

type FXExchange struct {
	ExchangeID   string `json:"exchange_id"`
	AccountID    string `json:"account_id"`
	FromCurrency string `json:"from_currency"`
	ToCurrency   string `json:"to_currency"`
	FromAmount   string `json:"from_amount"`
	ToAmount     string `json:"to_amount"`
	Rate         string `json:"rate"`
	Status       string `json:"status"`
	CreateTime   string `json:"create_time"`
}

func (c *Client) CreateFXExchange(ctx context.Context, req CreateFXExchangeRequest) (*FXExchange, error) {
	var out FXExchange
	if err := c.post(ctx, pathFXExchange, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetFXExchangeDetail(ctx context.Context, exchangeID string) (*FXExchange, error) {
	q := url.Values{}
	q.Set("exchange_id", exchangeID)
	var out FXExchange
	if err := c.get(ctx, pathFXExchangeDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type CreateInstantExchangeRequest struct {
	AccountID    string `json:"account_id"`
	FromCurrency string `json:"from_currency"`
	ToCurrency   string `json:"to_currency"`
	Amount       string `json:"amount"`
}

type InstantExchange struct {
	ExchangeID   string `json:"exchange_id"`
	AccountID    string `json:"account_id"`
	FromCurrency string `json:"from_currency"`
	ToCurrency   string `json:"to_currency"`
	Amount       string `json:"amount"`
	Rate         string `json:"rate"`
	Status       string `json:"status"`
	CreateTime   string `json:"create_time"`
}

func (c *Client) CreateInstantExchange(ctx context.Context, req CreateInstantExchangeRequest) (*InstantExchange, error) {
	var out InstantExchange
	if err := c.post(ctx, pathInstantExchange, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetInstantExchangeDetail(ctx context.Context, exchangeID string) (*InstantExchange, error) {
	q := url.Values{}
	q.Set("exchange_id", exchangeID)
	var out InstantExchange
	if err := c.get(ctx, pathFXExchangeDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type CreateInstantFundingRequest struct {
	AccountID string `json:"account_id"`
	Amount    string `json:"amount"`
	Currency  string `json:"currency"`
}

type InstantFunding struct {
	FundingID  string `json:"funding_id"`
	AccountID  string `json:"account_id"`
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
	Status     string `json:"status"`
	CreateTime string `json:"create_time"`
}

func (c *Client) CreateInstantFunding(ctx context.Context, req CreateInstantFundingRequest) (*InstantFunding, error) {
	var out InstantFunding
	if err := c.post(ctx, pathInstantFunding, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetInstantFundingDetail(ctx context.Context, fundingID string) (*InstantFunding, error) {
	q := url.Values{}
	q.Set("funding_id", fundingID)
	var out InstantFunding
	if err := c.get(ctx, pathInstantFundingDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
