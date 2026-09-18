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
	"net/http"
	"net/url"
)

// AccountSummary is a best-effort account summary type.
// Live probe needed to confirm all fields.
type AccountSummary struct {
	AccountID    string `json:"account_id"`
	Currency     string `json:"currency"`
	NetLiquidity string `json:"net_liquidity"`
	CashBalance  string `json:"cash_balance"`
	MarketValue  string `json:"market_value"`
	BuyingPower  string `json:"buying_power"`
}

// PositionSummary is a best-effort position type.
type PositionSummary struct {
	Symbol       string `json:"symbol"`
	Quantity     string `json:"quantity"`
	MarketValue  string `json:"market_value"`
	CostBasis    string `json:"cost_basis"`
	UnrealizedPL string `json:"unrealized_pl"`
}

// GetAccountsSummary retrieves account summaries for the linked account.
//
// Path and schema unconfirmed; live probe required.
func (c *Client) GetAccountsSummary(ctx context.Context) ([]AccountSummary, error) {
	path := "/broker-fd/accounts"
	var out []AccountSummary
	if err := c.get(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetPositions retrieves positions for the linked account.
//
// Path and schema unconfirmed; live probe required.
func (c *Client) GetPositions(ctx context.Context) ([]PositionSummary, error) {
	path := "/broker-fd/positions"
	var out []PositionSummary
	if err := c.get(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) get(ctx context.Context, path string, query url.Values, out any) error {
	if len(query) > 0 {
		path += "?" + query.Encode()
	}
	return c.core.Do(ctx, http.MethodGet, path, nil, out)
}
