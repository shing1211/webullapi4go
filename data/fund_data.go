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

	"github.com/shing1211/webullapi4go/pkg/domain/money"
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
	Symbol         string      `json:"symbol"`
	Name           string      `json:"name"`
	Currency       string      `json:"currency"`
	Exchange       string      `json:"exchange"`
	Nav            money.Money `json:"nav"`
	NavDate        string      `json:"nav_date"`
	PrevNav        money.Money `json:"prev_nav"`
	NavChange      money.Money `json:"nav_change"`
	NavChangeRatio string      `json:"nav_change_ratio"`

	// The two fields below are declared by the page but not marked required, so they
	// were absent here and a response carrying them decoded the value to the zero
	// value with no error reported. That is weaker evidence than a
	// missing-required-name row: the page publishes no required list for this call,
	// so a name may be optional or conditionally sent.
	//
	// The page's NetValue is this type's Nav and its Date is this type's NavDate, so
	// both sets are carried and a response populates whichever it sends. This type
	// already carries Currency, which the page's element also declares.
	NetValue money.Money `json:"net_value"`
	Date     string      `json:"date"`
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
	Symbol        string      `json:"symbol"`
	Name          string      `json:"name"`
	Currency      string      `json:"currency"`
	Exchange      string      `json:"exchange"`
	Aum           money.Money `json:"aum"`
	ExpenseRatio  string      `json:"expense_ratio"`
	DividendYield string      `json:"dividend_yield"`
	InceptionDate string      `json:"inception_date"`
	FundType      string      `json:"fund_type"`
	Category      string      `json:"category"`

	// The six fields below are declared by the page but not marked required, so they
	// were absent here and a response carrying them decoded the value to the zero
	// value with no error reported. That is weaker evidence than a
	// missing-required-name row: the page publishes no required list for this call,
	// so a name may be optional or conditionally sent.
	//
	// The page's LaunchDate is the same fact this type spells InceptionDate, so both
	// are carried. The remaining five have no counterpart above.
	Issuer              string            `json:"issuer"`
	Custodian           string            `json:"custodian"`
	Benchmark           string            `json:"benchmark"`
	InvestmentObjective string            `json:"investment_objective"`
	LaunchDate          string            `json:"launch_date"`
	Managers            []FundInfoManager `json:"managers"`
}

// FundInfoManager is one manager of a fund.
//
// The page declares every field without marking any required, so each may be absent.
type FundInfoManager struct {
	// Name is the manager's name.
	Name string `json:"name"`
	// Title is the manager's role, for example "Manager".
	Title string `json:"title"`
	// StartDate is when the manager took the role, as YYYY-MM-DD.
	StartDate string `json:"start_date"`
	// EndDate is when the manager left the role, as YYYY-MM-DD, or empty while
	// incumbent.
	EndDate string `json:"end_date"`
	// IsIncumbent reports whether the manager currently holds the role. The page
	// documents it as an integer, so it is carried as int64 and a caller reads
	// Incumbent rather than assuming a bool.
	IsIncumbent int64 `json:"is_incumbent"`
	// TenureDays is the manager's tenure in days.
	TenureDays int64 `json:"tenure_days"`
	// TenureYears is the manager's tenure in years, as a decimal string.
	TenureYears string `json:"tenure_years"`
	// TenureReturn is the return over the manager's tenure, as a decimal string.
	TenureReturn string `json:"tenure_return"`
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
	Symbol     string      `json:"symbol"`
	Name       string      `json:"name"`
	Currency   string      `json:"currency"`
	Exchange   string      `json:"exchange"`
	Amount     money.Money `json:"amount"`
	ExDate     string      `json:"ex_date"`
	PayDate    string      `json:"pay_date"`
	RecordDate string      `json:"record_date"`
	Frequency  string      `json:"frequency"`
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
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	Currency      string `json:"currency"`
	Exchange      string `json:"exchange"`
	FundType      string `json:"fund_type"`
	Category      string `json:"category"`
	DividendYield string `json:"dividend_yield"`
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
