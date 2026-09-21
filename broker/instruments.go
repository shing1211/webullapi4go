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
	"strings"
)

const (
	pathStockInstruments       = "/openapi/v1/broker/instruments/stocks"
	pathStockLocate            = "/openapi/v1/broker/instruments/stock-locate"
	pathCorporateActionsDetail = "/openapi/v1/broker/instruments/corporate-actions/detail"
)

// StockInstrument represents a stock instrument in the broker system.
type StockInstrument struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Currency string `json:"currency"`
	LotSize  string `json:"lot_size"`
	Status   string `json:"status"`
}

// GetStockInstruments retrieves stock instruments for a list of symbols.
func (c *Client) GetStockInstruments(ctx context.Context, symbols []string) ([]StockInstrument, error) {
	q := url.Values{}
	q.Set("symbols", strings.Join(symbols, ","))
	var out []StockInstrument
	if err := c.get(ctx, pathStockInstruments, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLocate represents a stock locate entry.
type StockLocate struct {
	Symbol         string `json:"symbol"`
	LocateQuantity string `json:"locate_quantity"`
	Available      string `json:"available"`
	Rate           string `json:"rate"`
}

// GetStockLocate retrieves locate information for a symbol.
func (c *Client) GetStockLocate(ctx context.Context, symbol string) ([]StockLocate, error) {
	q := url.Values{}
	q.Set("symbol", symbol)
	var out []StockLocate
	if err := c.get(ctx, pathStockLocate, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CorporateActionDetail represents corporate action details.
type CorporateActionDetail struct {
	Symbol      string `json:"symbol"`
	ActionType  string `json:"action_type"`
	ExDate      string `json:"ex_date"`
	RecordDate  string `json:"record_date"`
	PayableDate string `json:"payable_date"`
	Ratio       string `json:"ratio"`
	Description string `json:"description"`
}

// GetCorporateActionsDetail retrieves corporate action details for a symbol.
func (c *Client) GetCorporateActionsDetail(ctx context.Context, symbol string) ([]CorporateActionDetail, error) {
	q := url.Values{}
	q.Set("symbol", symbol)
	var out []CorporateActionDetail
	if err := c.get(ctx, pathCorporateActionsDetail, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}
