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
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

// Futures static-data endpoints.
const (
	pathFuturesInstruments    = "/trading/instruments/futures/contracts/list"
	pathFuturesProductCodes   = "/trading/instruments/futures/product-codes/list"
	pathFuturesProductClasses = "/trading/instruments/futures/product-classes/list"
)

// FuturesCategory identifies the futures market to query.
type FuturesCategory string

// Futures market categories.
const (
	// FuturesCategoryUS identifies United States futures.
	FuturesCategoryUS FuturesCategory = "US_FUTURES"
	// FuturesCategoryHK identifies Hong Kong futures.
	FuturesCategoryHK FuturesCategory = "HK_FUTURES"
	// FuturesCategoryCN identifies China futures.
	FuturesCategoryCN FuturesCategory = "CN_FUTURES"
)

// FuturesContractType distinguishes regular month contracts from continuous
// main contracts.
type FuturesContractType string

// Futures contract types.
const (
	// FuturesContractTypeMonthly is a regular delivery-month contract.
	FuturesContractTypeMonthly FuturesContractType = "MONTHLY"
	// FuturesContractTypeMain is a main/continuous contract.
	FuturesContractTypeMain FuturesContractType = "MAIN"
)

// FuturesSettlement is the settlement method of a futures contract.
type FuturesSettlement string

// Futures settlement methods.
const (
	// FuturesSettlementCash settles in cash.
	FuturesSettlementCash FuturesSettlement = "Cash"
	// FuturesSettlementPhysical settles by physical delivery.
	FuturesSettlementPhysical FuturesSettlement = "Physical"
)

// StringOrNumber handles JSON values that may be either a string or a numeric
// type. The Webull API occasionally sends numeric values where a string is
// expected (for example, the futures instrument unit field).
type StringOrNumber struct {
	Str string
}

func (s *StringOrNumber) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if data[0] == '"' {
		return json.Unmarshal(data, &s.Str)
	}
	n, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return err
	}
	s.Str = strconv.FormatFloat(n, 'f', -1, 64)
	return nil
}

func (s StringOrNumber) String() string {
	return s.Str
}

// FuturesInstrumentQuery parameterizes [Client.GetFuturesInstruments]. Category
// is required, and at least one of Symbols or Code must be provided.
type FuturesInstrumentQuery struct {
	// Category is the futures market to query. Required.
	Category FuturesCategory
	// Symbols restricts the result to the named trading symbols (for example
	// "ESZ5" or "NQZ5"), at most 100 per query.
	Symbols []string
	// Code restricts the result to a product code (for example "ES"). Either
	// Symbols or Code must be set.
	Code string
	// Status filters by tradable status.
	Status InstrumentStatus
}

// FuturesInstrument is the static profile of a single futures contract.
type FuturesInstrument struct {
	// Symbol is the contract symbol used in trading and market data, for
	// example "ESZ5" or the continuous "ESmain".
	Symbol string `json:"symbol"`
	// InstrumentID is the unique identifier of the contract. For a main or
	// continuous contract this identifies the main contract itself.
	InstrumentID string `json:"instrument_id"`
	// ExchangeCode is the exchange code, for example "XCME".
	ExchangeCode string `json:"exchange_code"`
	// Code is the product code, for example "ES".
	Code string `json:"code"`
	// Name is the display name of the contract.
	Name string `json:"name"`
	// ProductClassID identifies the product class, for example 2.
	ProductClassID int32 `json:"product_class_id"`
	// ProductClassName is the product class name, for example "Equities".
	ProductClassName string `json:"product_class_name"`
	// Status is the tradable status.
	Status InstrumentStatus `json:"status"`
	// Currency is the trading currency, for example "USD".
	Currency string `json:"currency"`
	// ContractMonth is the delivery month in yyyyMM form, for example "202512".
	ContractMonth string `json:"contract_month"`
	// SettlementDate is the final settlement date (YYYY-MM-DD).
	SettlementDate string `json:"settlement_date"`
	// Size is the contract multiplier, as a decimal string.
	Size string `json:"size"`
	// Unit describes the pricing unit and quantity. The API returns this as
	// either a string (e.g., "1-index points") or a bare number.
	Unit StringOrNumber `json:"unit"`
	// MinTick is the minimum price increment, as a decimal string.
	MinTick string `json:"min_tick"`
	// FirstNoticeDate is the first notice date (YYYY-MM-DD), when applicable.
	FirstNoticeDate string `json:"first_notice_date"`
	// LastNoticeDate is the last notice date (YYYY-MM-DD), when applicable.
	LastNoticeDate string `json:"last_notice_date"`
	// FirstTradingDate is the first tradable date (YYYY-MM-DD).
	FirstTradingDate string `json:"first_trading_date"`
	// LastTradingDate is the final trading date (YYYY-MM-DD).
	LastTradingDate string `json:"last_trading_date"`
	// ContractType distinguishes monthly from main/continuous contracts.
	ContractType FuturesContractType `json:"contract_type"`
	// Settlement is the settlement method.
	Settlement FuturesSettlement `json:"settlement"`
}

// FuturesProduct is a futures underlying product and its code.
type FuturesProduct struct {
	// Name is the display name, for example "E-Mini S&P 500".
	Name string `json:"name"`
	// Code is the product code, for example "ES".
	Code string `json:"code"`
	// ProductClassID identifies the product class.
	ProductClassID int32 `json:"product_class_id"`
	// ProductClassName is the product class name, for example "Equities".
	ProductClassName string `json:"product_class_name"`
	// ExchangeCode is the exchange code, for example "XCME".
	ExchangeCode string `json:"exchange_code"`
}

// FuturesProductClass is a futures product classification group.
type FuturesProductClass struct {
	// ProductClassID is the product class identifier.
	ProductClassID int32 `json:"product_class_id"`
	// ProductClassName is the product class name, for example "Equities".
	ProductClassName string `json:"product_class_name"`
}

// FuturesProductCodeQuery parameterizes [Client.GetFuturesProductCodes].
// Category is required.
type FuturesProductCodeQuery struct {
	// Category is the futures market to query. Required.
	Category FuturesCategory
	// ProductClassID optionally restricts the result to one product class.
	// Zero means no filter.
	ProductClassID int32
}

// GetFuturesInstruments retrieves static detail for one or more futures
// contracts. Category is required, and at least one of Symbols or Code must be
// provided.
//
// Reference: https://developer.webull.hk/apis/docs/reference/futures-instrument-list.md
func (c *Client) GetFuturesInstruments(ctx context.Context, q FuturesInstrumentQuery) ([]FuturesInstrument, error) {
	query := url.Values{}
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	if q.Code != "" {
		query.Set("code", q.Code)
	}
	if q.Status != "" {
		query.Set("status", string(q.Status))
	}

	var out []FuturesInstrument
	if err := c.get(ctx, pathFuturesInstruments, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFuturesProductCodes retrieves the futures products and their product codes
// for q.Category.
//
// Reference: https://developer.webull.hk/apis/docs/reference/futures-products.md
func (c *Client) GetFuturesProductCodes(ctx context.Context, q FuturesProductCodeQuery) ([]FuturesProduct, error) {
	query := url.Values{}
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	if q.ProductClassID > 0 {
		query.Set("product_class_id", strconv.FormatInt(int64(q.ProductClassID), 10))
	}

	var out []FuturesProduct
	if err := c.get(ctx, pathFuturesProductCodes, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFuturesProductClasses retrieves the futures product classification groups
// for category.
//
// Reference: https://developer.webull.hk/apis/docs/reference/futures-products-class.md
func (c *Client) GetFuturesProductClasses(ctx context.Context, category FuturesCategory) ([]FuturesProductClass, error) {
	query := url.Values{}
	if category != "" {
		query.Set("category", string(category))
	}

	var out []FuturesProductClass
	if err := c.get(ctx, pathFuturesProductClasses, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}
