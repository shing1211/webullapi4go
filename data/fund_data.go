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
	"strconv"
)

// pathFundNav is the fund NAV history endpoint.
const pathFundNav = "/market-data/fundamentals/fund-net-values/get"

// pathFundInfo is the fund basic info endpoint.
const pathFundInfo = "/market-data/fundamentals/fund-brief/get"

// pathFundDividends is the fund dividend history endpoint.
const pathFundDividends = "/market-data/fundamentals/fund-dividends/get"

// pathFundList is the fund list by market/category endpoint.
const pathFundList = "/market-data/fund/list"

// FundNavQuery parameterizes [Client.GetFundNav].
type FundNavQuery struct {
	Symbol    string
	StartDate string
	EndDate   string
	PageSize  int
}

// FundNav represents fund NAV (Net Asset Value) history data.
type FundNav struct {
	Symbol         string            `json:"symbol"`
	Name           string            `json:"name"`
	Currency       string            `json:"currency"`
	Exchange       string            `json:"exchange"`
	Nav            string            `json:"nav"`
	NavDate        string            `json:"nav_date"`
	PrevNav        string            `json:"prev_nav"`
	NavChange      string            `json:"nav_change"`
	NavChangeRatio string            `json:"nav_change_ratio"`
	Extra          map[string]string `json:"-"`
}

// GetFundNav retrieves NAV history for a fund or ETF.
func (c *Client) GetFundNav(ctx context.Context, q FundNavQuery) ([]FundNav, error) {
	query := url.Values{}
	if q.Symbol != "" {
		query.Set("symbol", q.Symbol)
	}
	if q.StartDate != "" {
		query.Set("start_date", q.StartDate)
	}
	if q.EndDate != "" {
		query.Set("end_date", q.EndDate)
	}
	if q.PageSize > 0 {
		query.Set("page_size", strconv.Itoa(q.PageSize))
	}

	var out []FundNav
	if err := c.get(ctx, pathFundNav, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FundInfoQuery parameterizes [Client.GetFundInfo].
type FundInfoQuery struct {
	Symbol string
}

// FundInfo represents basic information for a fund or ETF.
type FundInfo struct {
	Symbol        string            `json:"symbol"`
	Name          string            `json:"name"`
	Currency      string            `json:"currency"`
	Exchange      string            `json:"exchange"`
	Aum           string            `json:"aum"`
	ExpenseRatio  string            `json:"expense_ratio"`
	DividendYield string            `json:"dividend_yield"`
	InceptionDate string            `json:"inception_date"`
	FundType      string            `json:"fund_type"`
	Category      string            `json:"category"`
	Extra         map[string]string `json:"-"`
}

// GetFundInfo retrieves basic information for a fund or ETF.
func (c *Client) GetFundInfo(ctx context.Context, q FundInfoQuery) (*FundInfo, error) {
	query := url.Values{}
	if q.Symbol != "" {
		query.Set("symbol", q.Symbol)
	}
	var out FundInfo
	if err := c.get(ctx, pathFundInfo, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FundDividendsQuery parameterizes [Client.GetFundDividends].
type FundDividendsQuery struct {
	Symbol    string
	StartDate string
	EndDate   string
	PageSize  int
}

// FundDividend represents a fund or ETF dividend event.
type FundDividend struct {
	Symbol     string            `json:"symbol"`
	Name       string            `json:"name"`
	Currency   string            `json:"currency"`
	Exchange   string            `json:"exchange"`
	Amount     string            `json:"amount"`
	ExDate     string            `json:"ex_date"`
	PayDate    string            `json:"pay_date"`
	RecordDate string            `json:"record_date"`
	Frequency  string            `json:"frequency"`
	Extra      map[string]string `json:"-"`
}

// GetFundDividends retrieves dividend history for a fund or ETF.
func (c *Client) GetFundDividends(ctx context.Context, q FundDividendsQuery) ([]FundDividend, error) {
	query := url.Values{}
	if q.Symbol != "" {
		query.Set("symbol", q.Symbol)
	}
	if q.StartDate != "" {
		query.Set("start_date", q.StartDate)
	}
	if q.EndDate != "" {
		query.Set("end_date", q.EndDate)
	}
	if q.PageSize > 0 {
		query.Set("page_size", strconv.Itoa(q.PageSize))
	}

	var out []FundDividend
	if err := c.get(ctx, pathFundDividends, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FundListQuery parameterizes [Client.GetFundList].
type FundListQuery struct {
	Market   string
	Category string
	Exchange string
	PageSize int
}

// FundListItem represents a fund or ETF in a list response.
type FundListItem struct {
	Symbol        string            `json:"symbol"`
	Name          string            `json:"name"`
	Currency      string            `json:"currency"`
	Exchange      string            `json:"exchange"`
	FundType      string            `json:"fund_type"`
	Category      string            `json:"category"`
	DividendYield string            `json:"dividend_yield"`
	Extra         map[string]string `json:"-"`
}

// GetFundList retrieves a list of funds or ETFs by market and category.
func (c *Client) GetFundList(ctx context.Context, q FundListQuery) ([]FundListItem, error) {
	query := url.Values{}
	if q.Market != "" {
		query.Set("market", q.Market)
	}
	if q.Category != "" {
		query.Set("category", q.Category)
	}
	if q.Exchange != "" {
		query.Set("exchange", q.Exchange)
	}
	if q.PageSize > 0 {
		query.Set("page_size", strconv.Itoa(q.PageSize))
	}

	var out []FundListItem
	if err := c.get(ctx, pathFundList, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}
