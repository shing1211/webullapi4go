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
	pathBalance   = "/broker/assets/balances/get"
	pathPositions = "/broker/assets/positions/list"
)

type CurrencyAsset struct {
	Currency string `json:"currency"`
	Balance  string `json:"balance"`
}

type Balance struct {
	TotalAssetCurrency    string          `json:"total_asset_currency"`
	TotalCashBalance      string          `json:"total_cash_balance"`
	TotalMarketValue      string          `json:"total_market_value"`
	TotalUnrealizedPL     string          `json:"total_unrealized_profit_loss"`
	InitMargin            string          `json:"init_margin"`
	AccountCurrencyAssets []CurrencyAsset `json:"account_currency_assets"`
}

func (c *Client) GetBalance(ctx context.Context, accountID string) (*Balance, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out Balance
	if err := c.get(ctx, pathBalance, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type Position struct {
	Symbol         string `json:"symbol"`
	Quantity       string `json:"quantity"`
	AverageCost    string `json:"average_cost"`
	MarketValue    string `json:"market_value"`
	UnrealizedPL   string `json:"unrealized_pl"`
	InstrumentType string `json:"instrument_type"`
	Currency       string `json:"currency"`
}

type positionsResponse struct {
	Data []Position `json:"data"`
}

func (c *Client) GetPositions(ctx context.Context, accountID string) ([]Position, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out positionsResponse
	if err := c.get(ctx, pathPositions, q, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
