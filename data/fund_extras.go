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

package data

import (
	"context"
	"net/url"

	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

// Additional fund-data endpoints (US only).
const (
	pathFundPerformance = "/market-data/fundamentals/fund-performances/get"
	pathFundHoldings    = "/market-data/fundamentals/fund-holdings/get"
	pathFundRating      = "/market-data/fundamentals/fund-ratings/get"
	pathFundSplits      = "/market-data/fundamentals/fund-splits/get"
	pathFundFiles       = "/market-data/fundamentals/fund-files/get"
	pathFundAllocation  = "/market-data/fundamentals/fund-allocations/get"
)

func fundSymbolQuery(symbol string, category StockCategory) url.Values {
	query := url.Values{}
	if symbol != "" {
		query.Set("symbol", symbol)
	}
	if category != "" {
		query.Set("category", string(category))
	}
	return query
}

// FundPerformance is the performance-return history of a fund.
type FundPerformance struct {
	Currency  string `json:"currency"`
	EndDate   string `json:"end_date"`
	Return1M  string `json:"return_1m"`
	Return3M  string `json:"return_3m"`
	Return6M  string `json:"return_6m"`
	Return1Y  string `json:"return_1y"`
	Return3Y  string `json:"return_3y"`
	Return5Y  string `json:"return_5y"`
	Return10Y string `json:"return_10y"`
	ReturnSI  string `json:"return_si"`
}

// GetFundPerformance retrieves the performance returns of a fund.
func (c *Client) GetFundPerformance(ctx context.Context, symbol string, category StockCategory) (*FundPerformance, error) {
	var out FundPerformance
	if err := c.get(ctx, pathFundPerformance, fundSymbolQuery(symbol, category), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FundHolding is one fund holding.
type FundHolding struct {
	TargetSymbol    string `json:"target_symbol"`
	StockName       string `json:"stock_name"`
	ShareHeldPct    string `json:"share_held_pct"`
	ShareHeldChgPct string `json:"share_held_chg_pct"`
	MaturityDate    string `json:"maturity_date"`
	UpdateTime      string `json:"update_time"`
}

// GetFundHoldings retrieves the top holdings of a fund.
func (c *Client) GetFundHoldings(ctx context.Context, symbol string, category StockCategory) ([]FundHolding, error) {
	var out []FundHolding
	if err := c.get(ctx, pathFundHoldings, fundSymbolQuery(symbol, category), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FundRating is one fund rating record.
type FundRating struct {
	RatingDate    string `json:"rating_date"`
	RatingAgency  string `json:"rating_agency"`
	RatingCycle   string `json:"rating_cycle"`
	RatingResults int    `json:"rating_results"`
}

// GetFundRating retrieves the rating history of a fund.
func (c *Client) GetFundRating(ctx context.Context, symbol string, category StockCategory) ([]FundRating, error) {
	var out []FundRating
	if err := c.get(ctx, pathFundRating, fundSymbolQuery(symbol, category), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FundSplit is one fund split record.
type FundSplit struct {
	SplitDate  string  `json:"split_date"`
	SplitType  string  `json:"split_type"`
	SplitRatio string  `json:"split_ratio"`
	From       float64 `json:"from"`
	To         float64 `json:"to"`
}

// GetFundSplits retrieves the split history of a fund.
func (c *Client) GetFundSplits(ctx context.Context, symbol string, category StockCategory) ([]FundSplit, error) {
	var out []FundSplit
	if err := c.get(ctx, pathFundSplits, fundSymbolQuery(symbol, category), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FundFile is one fund document.
type FundFile struct {
	PublishDate string `json:"publish_date"`
	URL         string `json:"url"`
	Type        int    `json:"type"`
	FileName    string `json:"file_name"`
}

// GetFundFiles retrieves the documents of a fund.
func (c *Client) GetFundFiles(ctx context.Context, symbol string, category StockCategory) ([]FundFile, error) {
	var out []FundFile
	if err := c.get(ctx, pathFundFiles, fundSymbolQuery(symbol, category), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FundAllocation is one fund asset-allocation record. The asset-class fields
// are returned as decoded objects because their schema is not fixed.
type FundAllocation struct {
	Date        string         `json:"date"`
	Aum         money.Money    `json:"aum"`
	Cash        map[string]any `json:"cash"`
	Bond        map[string]any `json:"bond"`
	Stock       map[string]any `json:"stock"`
	Preferred   map[string]any `json:"preferred"`
	Convertible map[string]any `json:"convertible"`
	Other       map[string]any `json:"other"`
}

// GetFundAllocation retrieves the asset allocation history of a fund.
func (c *Client) GetFundAllocation(ctx context.Context, symbol string, category StockCategory) ([]FundAllocation, error) {
	var out []FundAllocation
	if err := c.get(ctx, pathFundAllocation, fundSymbolQuery(symbol, category), &out); err != nil {
		return nil, err
	}
	return out, nil
}
