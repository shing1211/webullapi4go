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

// Position represents a single position in a Broker API HK account.
//
// CostPrice, LastPrice and UnrealizedProfitLoss exist because the documented
// response for `GET /broker/assets/positions/list` requires `cost_price`,
// `last_price` and `unrealized_profit_loss`, while this type carries
// `average_cost`, `market_value` and `unrealized_pl`. The two sets are held side
// by side rather than one replacing the other, because which a live server sends
// is unverified -- Broker API HK returns 401 ROUTE_NOT_PERMITTED in the HK
// sandbox -- and a retag would break a call that currently works for anyone whose
// server sends the SDK's own names. A response populates whichever of the two
// names it carries, so a caller reading either name reads a value when the server
// sends it; nothing here asserts which of the two a live server sends, because
// that is exactly what is unverified.
//
// This is additive and therefore not a breaking change. Before it, a caller
// reading the documented names got a successful call and three silently zeroed
// values with nothing in the return path to distinguish them from a genuinely
// zero position. brokerfd.FDPosition carries the same three names for the same
// reason.
type Position struct {
	Symbol      string `json:"symbol"`
	Quantity    string `json:"quantity"`
	AverageCost string `json:"average_cost"`
	// CostPrice is the documented name for the cost basis, carried alongside
	// AverageCost. See the type comment on which of the two names is populated.
	CostPrice   string `json:"cost_price"`
	MarketValue string `json:"market_value"`
	// LastPrice is the documented name for the last traded price, carried
	// alongside MarketValue. See the type comment.
	LastPrice    string `json:"last_price"`
	UnrealizedPL string `json:"unrealized_pl"`
	// UnrealizedProfitLoss is the documented name for open P&L, carried alongside
	// UnrealizedPL. See the type comment.
	UnrealizedProfitLoss string `json:"unrealized_profit_loss"`
	InstrumentType       string `json:"instrument_type"`
	Currency             string `json:"currency"`
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
