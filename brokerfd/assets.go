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

	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

const (
	pathFDAssetsSummary   = "/broker-fd/assets/summary"
	pathFDAssetsDetail    = "/broker/assets/balances/get"
	pathFDAssetsPositions = "/broker/assets/positions/list"
)

// FDAssetsSummary contains aggregate asset data for a fractional shares account.
// Values are presented as strings to preserve numeric precision.
type FDAssetsSummary struct {
	AccountID    string      `json:"account_id"`
	TotalEquity  money.Money `json:"total_equity"`
	CashBalance  money.Money `json:"cash_balance"`
	MarketValue  money.Money `json:"market_value"`
	BuyingPower  money.Money `json:"buying_power"`
	UnrealizedPL money.Money `json:"unrealized_pl"`
	RealizedPL   money.Money `json:"realized_pl"`
	Currency     string      `json:"currency"`
}

// GetFDAssetsSummary retrieves the aggregate asset summary for a fractional shares account.
// It requires a valid accountID and returns summary fields including equity, cash, and P/L.
func (c *Client) GetFDAssetsSummary(ctx context.Context, accountID string) (*FDAssetsSummary, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out FDAssetsSummary
	if err := c.get(ctx, pathFDAssetsSummary, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FDAssetDetail provides a per-currency breakdown of cash and buying power
// for a fractional shares account. Values are presented as strings to preserve numeric precision.
type FDAssetDetail struct {
	Currency      string      `json:"currency"`
	CashBalance   money.Money `json:"cash_balance"`
	MarketValue   money.Money `json:"market_value"`
	BuyingPower   money.Money `json:"buying_power"`
	AvailableCash money.Money `json:"available_cash"`
}

// GetFDAssetsDetail retrieves per-currency asset detail for a fractional shares account.
// Returns a slice of FDAssetDetail, one per supported currency.
func (c *Client) GetFDAssetsDetail(ctx context.Context, accountID string) ([]FDAssetDetail, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out []FDAssetDetail
	if err := c.get(ctx, pathFDAssetsDetail, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FDPosition represents a single position held in a fractional shares account.
// Values such as quantity, cost, and P/L are strings to preserve numeric precision.
//
// The three CostPrice, LastPrice and UnrealizedProfitLoss fields exist because the
// documented response for `GET /broker/assets/positions/list` requires
// `cost_price`, `last_price` and `unrealized_profit_loss`, while this type has
// always carried `average_cost`, `market_value` and `unrealized_pl` instead. The
// documented schema and the SDK therefore disagreed in both directions: three
// required names reached no field, and three SDK names appeared in no documented
// property. The two sets are carried side by side rather than one replacing the
// other, because which of them a live server sends is unverified — the Broker FD
// host returns 404 in the HK sandbox and no US credential was available here — and
// a retag would break a call that currently works for anyone whose server sends the
// SDK's own names. A response populates whichever of the two names it carries,
// so a caller reading either name reads a value when the server sends it; nothing
// here asserts which of the two a live server sends, because that is exactly what
// is unverified.
//
// This is additive and therefore not a breaking change: a retag would have been.
// Before this, a caller reading the documented names got a successful call, a
// non-nil slice, and three silently zeroed money.Money values where cost basis,
// last price and open P&L belong, with nothing in the return path to distinguish
// that from a genuinely zero position.
type FDPosition struct {
	PositionID  string      `json:"position_id"`
	AccountID   string      `json:"account_id"`
	Symbol      string      `json:"symbol"`
	Quantity    string      `json:"quantity"`
	AverageCost money.Money `json:"average_cost"`
	// CostPrice is the documented name for the cost basis, carried alongside
	// AverageCost. See the type comment on which of the two names is populated.
	CostPrice   money.Money `json:"cost_price"`
	MarketValue money.Money `json:"market_value"`
	// LastPrice is the documented name for the last traded price, carried
	// alongside MarketValue. See the type comment.
	LastPrice    money.Money `json:"last_price"`
	UnrealizedPL money.Money `json:"unrealized_pl"`
	// UnrealizedProfitLoss is the documented name for open P&L, carried alongside
	// UnrealizedPL. See the type comment.
	UnrealizedProfitLoss money.Money `json:"unrealized_profit_loss"`
	// RealizedPL is closed P&L. Its tag was `unrealized_pl`, duplicating the field
	// above, which is a pre-existing defect this change corrects: when two fields
	// of one struct encode under the same name, encoding/json drops BOTH, so
	// UnrealizedPL was silently zero for every caller and realized_pl was never
	// decodable at all. Verified against the decoder before and after.
	RealizedPL     money.Money `json:"realized_pl"`
	InstrumentType string      `json:"instrument_type"`
	Currency       string      `json:"currency"`
}

// GetFDPositions retrieves all open positions for a fractional shares account.
func (c *Client) GetFDPositions(ctx context.Context, accountID string) ([]FDPosition, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out []FDPosition
	if err := c.get(ctx, pathFDAssetsPositions, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}
