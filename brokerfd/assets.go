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
type FDPosition struct {
	PositionID     string      `json:"position_id"`
	AccountID      string      `json:"account_id"`
	Symbol         string      `json:"symbol"`
	Quantity       string      `json:"quantity"`
	AverageCost    money.Money `json:"average_cost"`
	MarketValue    money.Money `json:"market_value"`
	UnrealizedPL   money.Money `json:"unrealized_pl"`
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
