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
)

// Fundamental-data endpoints.
const (
	pathCompanyProfile = "/market-data/fundamentals/company-profiles/get"
	pathAnalystTarget  = "/market-data/fundamentals/analysis/target-prices/get"
	pathAnalystRating  = "/market-data/fundamentals/analysis/ratings/get"
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
