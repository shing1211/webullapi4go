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
)

const (
	pathFDStockInstruments      = "/broker-fd/instruments/stocks"
	pathFDStockLocate           = "/broker-fd/instruments/stock-locate"
	pathFDCorporateActions      = "/broker-fd/instruments/corporate-actions"
	pathFDCorporateActionDetail = "/broker-fd/instruments/corporate-action/detail"
	pathFDECInstruments         = "/broker-fd/instruments/event-contracts"
	pathFDECInstrumentDetail    = "/broker-fd/instruments/event-contract/detail"
)

type FDStockInstrument struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Currency string `json:"currency"`
	LotSize  string `json:"lot_size"`
	Status   string `json:"status"`
}

func (c *Client) GetFDStockInstruments(ctx context.Context, symbols []string) ([]FDStockInstrument, error) {
	q := url.Values{}
	q.Set("symbols", strings.Join(symbols, ","))
	var out []FDStockInstrument
	if err := c.get(ctx, pathFDStockInstruments, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

type FDStockLocate struct {
	Symbol         string `json:"symbol"`
	LocateQuantity string `json:"locate_quantity"`
	Available      string `json:"available"`
	Rate           string `json:"rate"`
}

func (c *Client) GetFDStockLocate(ctx context.Context, symbol string) ([]FDStockLocate, error) {
	q := url.Values{}
	q.Set("symbol", symbol)
	var out []FDStockLocate
	if err := c.get(ctx, pathFDStockLocate, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

type FDCorporateAction struct {
	ActionID   string `json:"action_id"`
	Symbol     string `json:"symbol"`
	ActionType string `json:"action_type"`
	ExDate     string `json:"ex_date"`
	RecordDate string `json:"record_date"`
	PayDate    string `json:"pay_date"`
	Ratio      string `json:"ratio"`
}

func (c *Client) GetFDCorporateActions(ctx context.Context, symbol string) ([]FDCorporateAction, error) {
	q := url.Values{}
	q.Set("symbol", symbol)
	var out []FDCorporateAction
	if err := c.get(ctx, pathFDCorporateActions, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetFDCorporateActionDetail(ctx context.Context, actionID string) (*FDCorporateAction, error) {
	q := url.Values{}
	q.Set("action_id", actionID)
	var out FDCorporateAction
	if err := c.get(ctx, pathFDCorporateActionDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type FDECInstrument struct {
	Symbol         string `json:"symbol"`
	EventID        string `json:"event_id"`
	SeriesID       string `json:"series_id"`
	StrikePrice    string `json:"strike_price"`
	ExpirationDate string `json:"expiration_date"`
	Status         string `json:"status"`
}

func (c *Client) GetFDECInstruments(ctx context.Context, eventID string) ([]FDECInstrument, error) {
	q := url.Values{}
	q.Set("event_id", eventID)
	var out []FDECInstrument
	if err := c.get(ctx, pathFDECInstruments, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetFDECInstrumentDetail(ctx context.Context, symbol string) (*FDECInstrument, error) {
	q := url.Values{}
	q.Set("symbol", symbol)
	var out FDECInstrument
	if err := c.get(ctx, pathFDECInstrumentDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
