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
	"fmt"
	"net/url"
)

// Fundamental-data endpoints.
const (
	pathCompanyProfile = "/market-data/fundamentals/company-profiles/get"
	pathAnalystTarget  = "/market-data/fundamentals/analysis/target-prices/get"
	pathAnalystRating  = "/market-data/fundamentals/analysis/ratings/get"
	pathCapitalFlow    = "/market-data/fundamentals/capital-flows/get"
	pathIndustryComp   = "/market-data/fundamentals/industry-comparisons/get"
	pathEarningsCal    = "/market-data/fundamentals/earnings-calendars/list"
	pathDividendCal    = "/market-data/fundamentals/dividend-calendars/list"
	pathFilings        = "/market-data/fundamentals/filings/list"
	pathIncomeStmt     = "/market-data/fundamentals/income-statements/get"
	pathBalanceSheet   = "/market-data/fundamentals/balance-sheets/get"
	pathCashFlow       = "/market-data/fundamentals/cash-flows/get"
	pathIndicators     = "/market-data/fundamentals/indicators/get"
	pathFinancialAlert = "/market-data/fundamentals/financial-alerts/get"
)

// CompanyProfile describes a company's static profile.
type CompanyProfile struct {
	// Symbol is the security symbol, for example "AAPL".
	Symbol string `json:"symbol"`
	// Category is the instrument's market.
	Category StockCategory `json:"category"`
	// CompanyName is the registered company name.
	CompanyName string `json:"company_name"`
	// EstablishDate is the date of incorporation (YYYY-MM-DD).
	EstablishDate string `json:"establish_date"`
	// ExhibitionCode is the market where the company is listed.
	ExhibitionCode string `json:"exhibition_code"`
	// Profile is the free-text business description.
	Profile string `json:"profile"`
	// Employees is the number of employees, as a string.
	Employees string `json:"employees"`
	// Address is the headquarters address.
	Address string `json:"address"`
	// CEO is the name of the chief executive officer.
	CEO string `json:"ceo"`
	// Industries lists the company's industries.
	Industries []string `json:"industries"`
}

// AnalystTargetPrice holds analyst target-price statistics for a security.
type AnalystTargetPrice struct {
	// Symbol is the security symbol.
	Symbol string `json:"symbol"`
	// Category is the instrument's market.
	Category StockCategory `json:"category"`
	// Mean is the average target price, as a decimal string.
	Mean string `json:"mean"`
	// Low is the lowest target price, as a decimal string.
	Low string `json:"low"`
	// High is the highest target price, as a decimal string.
	High string `json:"high"`
	// Median is the median target price, as a decimal string.
	Median string `json:"median"`
	// Currency is the target-price currency, for example "USD".
	Currency string `json:"currency"`
	// EffectiveStartDate is when the figures became effective.
	EffectiveStartDate string `json:"effective_start_date"`
}

// AnalystRating holds analyst rating counts for a security.
type AnalystRating struct {
	// Symbol is the security symbol.
	Symbol string `json:"symbol"`
	// Category is the instrument's market.
	Category StockCategory `json:"category"`
	// Number is the total number of analysts, as a string.
	Number string `json:"number"`
	// UnderPerform is the number of under-perform ratings, as a string.
	UnderPerform string `json:"under_perform"`
	// Buy is the number of buy ratings, as a string.
	Buy string `json:"buy"`
	// Sell is the number of sell ratings, as a string.
	Sell string `json:"sell"`
	// StrongBuy is the number of strong-buy ratings, as a string.
	StrongBuy string `json:"strong_buy"`
	// Hold is the number of hold (neutral) ratings, as a string.
	Hold string `json:"hold"`
	// EffectiveStartDate is when the figures became effective.
	EffectiveStartDate string `json:"effective_start_date"`
}

// GetCompanyProfile retrieves the profile of the company behind symbol. The
// Webull documentation currently supports US stocks only, so category should
// normally be [StockCategoryUS].
//
// Reference: https://developer.webull.hk/apis/docs/reference/get-company-profile.md
func (c *Client) GetCompanyProfile(ctx context.Context, symbol string, category StockCategory) (*CompanyProfile, error) {
	query := url.Values{
		"symbol":   {symbol},
		"category": {string(category)},
	}
	var out CompanyProfile
	if err := c.get(ctx, pathCompanyProfile, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAnalystTargetPrice retrieves aggregate analyst target-price data for
// symbol.
//
// Reference: https://developer.webull.hk/apis/docs/reference/get-analyst-target-price.md
func (c *Client) GetAnalystTargetPrice(ctx context.Context, symbol string, category StockCategory) (*AnalystTargetPrice, error) {
	query := url.Values{
		"symbol":   {symbol},
		"category": {string(category)},
	}
	var out AnalystTargetPrice
	if err := c.get(ctx, pathAnalystTarget, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAnalystRating retrieves aggregate analyst rating counts for symbol.
//
// Reference: https://developer.webull.hk/apis/docs/reference/get-analyst-rating.md
func (c *Client) GetAnalystRating(ctx context.Context, symbol string, category StockCategory) (*AnalystRating, error) {
	query := url.Values{
		"symbol":   {symbol},
		"category": {string(category)},
	}
	var out AnalystRating
	if err := c.get(ctx, pathAnalystRating, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CapitalFlowEntry describes one trading day's capital flow breakdown.
type CapitalFlowEntry struct {
	Date      string `json:"date"`
	LargeIn   string `json:"large_in"`
	LargeOut  string `json:"large_out"`
	MediumIn  string `json:"medium_in"`
	MediumOut string `json:"medium_out"`
	SmallIn   string `json:"small_in"`
	SmallOut  string `json:"small_out"`
}

// GetCapitalFlow retrieves the capital flow breakdown for symbol over the most
// recent count trading days (default 5, maximum 5). Results are sorted in
// ascending chronological order.
//
// Reference: https://developer.webull.hk/apis/docs/reference/capital-flow.md
func (c *Client) GetCapitalFlow(ctx context.Context, symbol string, category StockCategory, count int) ([]CapitalFlowEntry, error) {
	q := url.Values{
		"symbol":   {symbol},
		"category": {string(category)},
	}
	if count > 0 {
		q.Set("count", fmt.Sprintf("%d", count))
	}
	var out []CapitalFlowEntry
	if err := c.get(ctx, pathCapitalFlow, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// IndustryComparisonItem is one entry in an industry comparison.
type IndustryComparisonItem struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Rank   int    `json:"rank"`
	Value  string `json:"value"`
}

// IndustryComparison holds industry comparison data for a single fiscal period.
type IndustryComparison struct {
	FiscalYear   int                      `json:"fiscal_year"`
	FiscalPeriod int                      `json:"fiscal_period"`
	IndustryName string                   `json:"industry_name"`
	Type         string                   `json:"type"`
	Data         []IndustryComparisonItem `json:"data"`
}

// GetIndustryComparison retrieves peer comparison data for the industry that
// symbol belongs to. sortBy defaults to EPS_TTM.
//
// Reference: https://developer.webull.hk/apis/docs/reference/industry-comparison.md
func (c *Client) GetIndustryComparison(ctx context.Context, symbol string, category StockCategory, sortBy string) (*IndustryComparison, error) {
	q := url.Values{
		"symbol":   {symbol},
		"category": {string(category)},
	}
	if sortBy != "" {
		q["sort_by"] = []string{sortBy}
	}
	var out IndustryComparison
	if err := c.get(ctx, pathIndustryComp, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EarningsCalendarEntry is one earnings announcement.
type EarningsCalendarEntry struct {
	FiscalYear          int    `json:"fiscal_year"`
	FiscalPeriod        int    `json:"fiscal_period"`
	Currency            string `json:"currency"`
	ExpectedPublishDate string `json:"expected_publish_date"`
	EPSActual           string `json:"eps_actual"`
	EPSEst              string `json:"eps_est"`
	RevActual           string `json:"rev_actual"`
	RevEst              string `json:"rev_est"`
}

// GetEarningsCalendar retrieves the earnings calendar for symbol.
//
// Reference: https://developer.webull.hk/apis/docs/reference/earnings-calendar.md
func (c *Client) GetEarningsCalendar(ctx context.Context, symbol string, category StockCategory) ([]EarningsCalendarEntry, error) {
	q := url.Values{
		"symbol":   {symbol},
		"category": {string(category)},
	}
	var out []EarningsCalendarEntry
	if err := c.get(ctx, pathEarningsCal, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DividendCalendarEntry is one dividend event.
type DividendCalendarEntry struct {
	Symbol      string `json:"symbol"`
	Market      string `json:"market"`
	Currency    string `json:"currency"`
	Amount      string `json:"amount"`
	DivType     string `json:"div_type"`
	DeclareDate string `json:"declare_date"`
	ExDivDate   string `json:"ex_div_date"`
	RecordDate  string `json:"record_date"`
	PayDate     string `json:"pay_date"`
}

// GetDividendCalendar retrieves the dividend calendar for symbol.
//
// Reference: https://developer.webull.hk/apis/docs/reference/dividend-calendar.md
func (c *Client) GetDividendCalendar(ctx context.Context, symbol string, category StockCategory) ([]DividendCalendarEntry, error) {
	q := url.Values{
		"symbol":   {symbol},
		"category": {string(category)},
	}
	var out []DividendCalendarEntry
	if err := c.get(ctx, pathDividendCal, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FilingEntry is one SEC filing.
type FilingEntry struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	PublishDate string `json:"publish_date"`
}

// FilingsResponse wraps the filings list.
type FilingsResponse struct {
	Symbol   string        `json:"symbol"`
	Category string        `json:"category"`
	Filings  []FilingEntry `json:"filings"`
}

// GetFilings retrieves SEC filings for symbol.
//
// Reference: https://developer.webull.hk/apis/docs/reference/filings.md
func (c *Client) GetFilings(ctx context.Context, symbol string, category StockCategory) (*FilingsResponse, error) {
	q := url.Values{
		"symbol":   {symbol},
		"category": {string(category)},
	}
	var out FilingsResponse
	if err := c.get(ctx, pathFilings, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FinancialsItem is one period's financial data entry. Field names are preserved
// from the API response. Values are any because the API returns a mix of strings
// and numbers (e.g. integer fiscal periods).
type FinancialsItem map[string]any

// GetIncomeStatement retrieves the income statement for symbol.
//
// Reference: https://developer.webull.hk/apis/docs/reference/income-statement.md
func (c *Client) GetIncomeStatement(ctx context.Context, symbol string, category StockCategory) ([]FinancialsItem, error) {
	q := url.Values{
		"symbol":   {symbol},
		"category": {string(category)},
	}
	var out []FinancialsItem
	if err := c.get(ctx, pathIncomeStmt, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetBalanceSheet retrieves the balance sheet for symbol.
//
// Reference: https://developer.webull.hk/apis/docs/reference/balance-sheet.md
func (c *Client) GetBalanceSheet(ctx context.Context, symbol string, category StockCategory) ([]FinancialsItem, error) {
	q := url.Values{
		"symbol":   {symbol},
		"category": {string(category)},
	}
	var out []FinancialsItem
	if err := c.get(ctx, pathBalanceSheet, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCashFlow retrieves the cash flow statement for symbol.
//
// Reference: https://developer.webull.hk/apis/docs/reference/cash-flow-statement.md
func (c *Client) GetCashFlow(ctx context.Context, symbol string, category StockCategory) ([]FinancialsItem, error) {
	q := url.Values{
		"symbol":   {symbol},
		"category": {string(category)},
	}
	var out []FinancialsItem
	if err := c.get(ctx, pathCashFlow, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FinancialIndicator holds key financial ratios and metrics.
type FinancialIndicator struct {
	ROA       string `json:"roa"`
	ROE       string `json:"roe"`
	EPS       string `json:"eps"`
	NetMargin string `json:"net_margin"`
	DebtRatio string `json:"debt_ratio"`
}

// GetFinancialIndicators retrieves financial indicators for symbol.
//
// Reference: https://developer.webull.hk/apis/docs/reference/financial-indicators.md
func (c *Client) GetFinancialIndicators(ctx context.Context, symbol string, category StockCategory) (*FinancialIndicator, error) {
	q := url.Values{
		"symbol":   {symbol},
		"category": {string(category)},
	}
	var out FinancialIndicator
	if err := c.get(ctx, pathIndicators, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FinancialAlert holds upcoming earnings-release alert information.
type FinancialAlert struct {
	Symbol             string `json:"symbol"`
	Category           string `json:"category"`
	ExpectedReportDate string `json:"expected_report_date"`
	EstimatedEPS       string `json:"estimated_eps"`
	LastYearEPS        string `json:"last_year_eps"`
}

// GetFinancialAlert retrieves upcoming earnings-release alert for symbol.
//
// Reference: https://developer.webull.hk/apis/docs/reference/financial-alert.md
func (c *Client) GetFinancialAlert(ctx context.Context, symbol string, category StockCategory) (*FinancialAlert, error) {
	q := url.Values{
		"symbol":   {symbol},
		"category": {string(category)},
	}
	var out FinancialAlert
	if err := c.get(ctx, pathFinancialAlert, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ForecastEPSEntry is one quarter's EPS forecast data.
type ForecastEPSEntry struct {
	FiscalYear   int    `json:"fiscal_year"`
	FiscalPeriod int    `json:"fiscal_period"`
	Actual       string `json:"actual"`
	Est          string `json:"est"`
	Reported     bool   `json:"reported"`
}

// GetForecastEPS retrieves forecast EPS data for symbol for the most recent
// 5 quarters.
//
// Reference: https://developer.webull.hk/apis/docs/reference/forecast-eps.md
func (c *Client) GetForecastEPS(ctx context.Context, symbol string, category StockCategory) ([]ForecastEPSEntry, error) {
	q := url.Values{
		"symbol":   {symbol},
		"category": {string(category)},
	}
	var out []ForecastEPSEntry
	if err := c.get(ctx, "/market-data/fundamentals/forecast-eps/get", q, &out); err != nil {
		return nil, err
	}
	return out, nil
}
