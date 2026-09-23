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
	"strings"

	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

const (
	pathFDStockInstruments      = "/broker/instruments/stocks/profiles/list"
	pathFDStockLocate           = "/broker-fd/instruments/stock-locate"
	pathFDCorporateActions      = "/broker/instruments/stocks/corporate-actions/get"
	pathFDCorporateActionDetail = "/broker-fd/instruments/corporate-action/detail"
	pathFDECInstruments         = "/broker/instruments/event-contracts/markets/list"
	pathFDECInstrumentDetail    = "/broker-fd/instruments/event-contract/detail"
)

// FDStockInstrument represents a fractional share stock instrument available for trading.
type FDStockInstrument struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Currency string `json:"currency"`
	LotSize  string `json:"lot_size"`
	Status   string `json:"status"`
}

// GetFDStockInstruments retrieves fractional share stock instruments by symbol list.
func (c *Client) GetFDStockInstruments(ctx context.Context, symbols []string) ([]FDStockInstrument, error) {
	q := url.Values{}
	q.Set("symbols", strings.Join(symbols, ","))
	var out []FDStockInstrument
	if err := c.get(ctx, pathFDStockInstruments, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FDStockLocate represents locate (borrowed shares) information for a fractional share.
type FDStockLocate struct {
	Symbol         string      `json:"symbol"`
	LocateQuantity string      `json:"locate_quantity"`
	Available      string      `json:"available"`
	Rate           money.Money `json:"rate"`
}

// GetFDStockLocate retrieves locate information for a fractional share symbol.
func (c *Client) GetFDStockLocate(ctx context.Context, symbol string) ([]FDStockLocate, error) {
	q := url.Values{}
	q.Set("symbol", symbol)
	var out []FDStockLocate
	if err := c.get(ctx, pathFDStockLocate, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FDCorporateAction represents a corporate action (dividend, split, etc.) for fractional shares.
type FDCorporateAction struct {
	ActionID   string `json:"action_id"`
	Symbol     string `json:"symbol"`
	ActionType string `json:"action_type"`
	ExDate     string `json:"ex_date"`
	RecordDate string `json:"record_date"`
	PayDate    string `json:"pay_date"`
	Ratio      string `json:"ratio"`
}

// GetFDCorporateActions retrieves corporate actions for a fractional share symbol.
func (c *Client) GetFDCorporateActions(ctx context.Context, symbol string) ([]FDCorporateAction, error) {
	q := url.Values{}
	q.Set("symbol", symbol)
	var out []FDCorporateAction
	if err := c.get(ctx, pathFDCorporateActions, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFDCorporateActionDetail retrieves details for a specific corporate action by its ID.
func (c *Client) GetFDCorporateActionDetail(ctx context.Context, actionID string) (*FDCorporateAction, error) {
	q := url.Values{}
	q.Set("action_id", actionID)
	var out FDCorporateAction
	if err := c.get(ctx, pathFDCorporateActionDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FDECInstrument represents an event contract (EC) instrument for fractional shares.
type FDECInstrument struct {
	Symbol         string      `json:"symbol"`
	EventID        string      `json:"event_id"`
	SeriesID       string      `json:"series_id"`
	StrikePrice    money.Money `json:"strike_price"`
	ExpirationDate string      `json:"expiration_date"`
	Status         string      `json:"status"`
}

// GetFDECInstruments retrieves event contract instruments by event ID.
func (c *Client) GetFDECInstruments(ctx context.Context, eventID string) ([]FDECInstrument, error) {
	q := url.Values{}
	q.Set("event_id", eventID)
	var out []FDECInstrument
	if err := c.get(ctx, pathFDECInstruments, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFDECInstrumentDetail retrieves details for a specific event contract by its symbol.
func (c *Client) GetFDECInstrumentDetail(ctx context.Context, symbol string) (*FDECInstrument, error) {
	q := url.Values{}
	q.Set("symbol", symbol)
	var out FDECInstrument
	if err := c.get(ctx, pathFDECInstrumentDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
