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

// Package brokerfd provides the Webull Broker FD (US) API client.
//
// A Client is constructed from the public client, so signing, token handling,
// retries, and rate limiting are shared with the rest of the SDK:
//
//	cl, err := client.New(
//		client.WithAppKey(key),
//		client.WithAppSecret(secret),
//		client.WithSandbox(),
//	)
//	if err != nil {
//		return err
//	}
//	bfd := brokerfd.New(cl)
package brokerfd

import (
	"context"
)

// AccountSummary holds cash and buying power data for a Broker FD account.
// Numeric fields are returned as strings to preserve precision.
type AccountSummary struct {
	AccountID    string `json:"account_id"`
	Currency     string `json:"currency"`
	NetLiquidity string `json:"net_liquidity"`
	CashBalance  string `json:"cash_balance"`
	MarketValue  string `json:"market_value"`
	BuyingPower  string `json:"buying_power"`
}

// PositionSummary holds position-level data for a Broker FD account.
// Numeric fields are returned as strings to preserve precision.
type PositionSummary struct {
	Symbol       string `json:"symbol"`
	Quantity     string `json:"quantity"`
	MarketValue  string `json:"market_value"`
	CostBasis    string `json:"cost_basis"`
	UnrealizedPL string `json:"unrealized_pl"`
}

// GetAccountsSummary returns cash and buying power data for all Broker FD accounts.
func (c *Client) GetAccountsSummary(ctx context.Context) ([]AccountSummary, error) {
	path := "/broker-fd/accounts"
	var out []AccountSummary
	if err := c.get(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetPositions returns position-level data for all positions held across Broker FD accounts.
func (c *Client) GetPositions(ctx context.Context) ([]PositionSummary, error) {
	path := "/broker-fd/positions"
	var out []PositionSummary
	if err := c.get(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
