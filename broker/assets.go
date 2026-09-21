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

// Broker asset endpoints.
const (
	pathBalance   = "/openapi/v1/broker/assets/balance"
	pathPositions = "/openapi/v1/broker/assets/positions"
)

// Balance represents the account balance.
type Balance struct {
	AccountID    string `json:"account_id"`
	TotalEquity  string `json:"total_equity"`
	CashBalance  string `json:"cash_balance"`
	MarketValue  string `json:"market_value"`
	BuyingPower  string `json:"buying_power"`
	UnrealizedPL string `json:"unrealized_pl"`
	Margin       string `json:"margin"`
	Currency     string `json:"currency"`
}

// GetBalance retrieves the account balance.
func (c *Client) GetBalance(ctx context.Context, accountID string) (*Balance, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out Balance
	if err := c.get(ctx, pathBalance, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Position represents a holding.
type Position struct {
	Symbol         string `json:"symbol"`
	Quantity       string `json:"quantity"`
	AverageCost    string `json:"average_cost"`
	MarketValue    string `json:"market_value"`
	UnrealizedPL   string `json:"unrealized_pl"`
	InstrumentType string `json:"instrument_type"`
	Currency       string `json:"currency"`
}

// GetPositions retrieves positions for an account.
func (c *Client) GetPositions(ctx context.Context, accountID string) ([]Position, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out []Position
	if err := c.get(ctx, pathPositions, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}
