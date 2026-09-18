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
	"strings"
)

// pathStockInstruments is the stock instrument-list endpoint.
//
// Reference: https://developer.webull.hk/apis/docs/reference/instrument-list.md
// (the linked reference documents the v3 path
// /trading/instruments/stocks/profiles/list, which returns the same instrument
// fields wrapped in a {data, pagination_key} envelope; the v2 path below
// returns the instrument array directly).
const pathStockInstruments = "/openapi/instrument/stock/list"

// StockCategory identifies the market of a stock instrument.
type StockCategory string

// Stock market categories accepted by the instrument-list endpoint.
const (
	// StockCategoryUS identifies United States stocks.
	StockCategoryUS StockCategory = "US_STOCK"
	// StockCategoryHK identifies Hong Kong stocks.
	StockCategoryHK StockCategory = "HK_STOCK"
	// StockCategoryCN identifies mainland China stocks.
	StockCategoryCN StockCategory = "CN_STOCK"
)

// InstrumentStatus is the tradable status of a stock or futures instrument.
type InstrumentStatus string

// Instrument tradable status values.
const (
	// InstrumentStatusTradable means the instrument is available for trading
	// (OC).
	InstrumentStatusTradable InstrumentStatus = "OC"
	// InstrumentStatusLiquidateOnly means the instrument can only be sold, with
	// no purchases allowed (CO).
	InstrumentStatusLiquidateOnly InstrumentStatus = "CO"
	// InstrumentStatusNonTradable means the instrument cannot be traded (NT).
	InstrumentStatusNonTradable InstrumentStatus = "NT"
)

// StockSubCategory is a finer classification of a stock instrument. It is only
// effective when a symbol list is not supplied.
type StockSubCategory string

// Stock sub-categories. All values are accepted for US stocks; Hong Kong and
// mainland China stocks support only [StockSubCategoryCommonStock] and
// [StockSubCategoryETF].
const (
	// StockSubCategoryCommonStock identifies common shares.
	StockSubCategoryCommonStock StockSubCategory = "COMMON_STOCK"
	// StockSubCategoryETF identifies exchange-traded funds.
	StockSubCategoryETF StockSubCategory = "ETF"
	// StockSubCategoryPreferredStock identifies preferred shares.
	StockSubCategoryPreferredStock StockSubCategory = "PREFERRED_STOCK"
	// StockSubCategoryWarrant identifies warrants.
	StockSubCategoryWarrant StockSubCategory = "WARRANT"
	// StockSubCategoryUnits identifies units.
	StockSubCategoryUnits StockSubCategory = "UNITS"
	// StockSubCategoryRight identifies rights.
	StockSubCategoryRight StockSubCategory = "RIGHT"
)

// StockInstrumentQuery parameterizes [Client.GetStockInstruments]. Category is
// required; every other field is optional.
type StockInstrumentQuery struct {
	// Category is the market to query. Required.
	Category StockCategory
	// Symbols restricts the result to the named symbols, at most 100 per query.
	// When empty the endpoint pages through all instruments in Category.
	Symbols []string
	// Status filters by tradable status.
	Status InstrumentStatus
	// SubCategory filters by sub-category. It is only effective when Symbols is
	// empty.
	SubCategory StockSubCategory
	// PaginationKey continues from a previous page. It is only used by
	// paginated endpoints.
	PaginationKey string
}

// StockInstrument is the profile of a single stock instrument.
type StockInstrument struct {
	// Name is the display name, for example "APPLE INC".
	Name string `json:"name"`
	// InstrumentID is the unique identifier of the security.
	InstrumentID string `json:"instrument_id"`
	// ExchangeCode is the exchange code, for example "NSQ" (Nasdaq).
	ExchangeCode string `json:"exchange_code"`
	// Category is the instrument's market.
	Category StockCategory `json:"category"`
	// Symbol is the trading symbol, for example "AAPL".
	Symbol string `json:"symbol"`
	// Status is the tradable status.
	Status InstrumentStatus `json:"status"`
	// Shortable reports whether the instrument can be sold short.
	Shortable bool `json:"shortable"`
	// Fractionable reports whether fractional trading is supported.
	Fractionable bool `json:"fractionable"`
	//nolint:misspell // "marginable" is the Webull API field name, not a typo.
	// Marginable reports whether the instrument can be bought on margin.
	Marginable bool `json:"marginable"` //nolint:misspell
	// OvernightTradingSupported reports whether overnight trading is supported.
	OvernightTradingSupported bool `json:"overnight_trading_supported"`
	// MarginRequirementLong is the margin requirement ratio for a long
	// position, as a decimal string.
	MarginRequirementLong string `json:"margin_requirement_long"`
	// MarginRequirementShort is the margin requirement ratio for a short
	// position, as a decimal string.
	MarginRequirementShort string `json:"margin_requirement_short"`
	// IntradayMarginLong is the intraday margin requirement ratio for a long
	// position, as a decimal string.
	IntradayMarginLong string `json:"intraday_margin_long"`
	// IntradayMarginShort is the intraday margin requirement ratio for a short
	// position, as a decimal string.
	IntradayMarginShort string `json:"intraday_margin_short"`
	// MaintenanceMarginLong is the maintenance margin ratio for a long
	// position, as a decimal string.
	MaintenanceMarginLong string `json:"maintenance_margin_long"`
	// MaintenanceMarginShort is the maintenance margin ratio for a short
	// position, as a decimal string.
	MaintenanceMarginShort string `json:"maintenance_margin_short"`
	// EasyToBorrow reports whether the instrument is easy to borrow.
	EasyToBorrow bool `json:"easy_to_borrow"`
	// LotSize is the minimum tradable quantity, as a decimal string.
	LotSize string `json:"lot_size"`
	// Currency is the trading currency, for example "USD".
	Currency string `json:"currency"`
	// SubCategory is the finer instrument classification.
	SubCategory StockSubCategory `json:"sub_category"`
}

// GetStockInstruments retrieves profile information for one or more stock
// instruments. When q.Symbols is empty it returns the instruments of q.Category
// page by page; when set, it returns the requested symbols directly.
//
// Reference: https://developer.webull.hk/apis/docs/reference/instrument-list.md
func (c *Client) GetStockInstruments(ctx context.Context, q StockInstrumentQuery) ([]StockInstrument, error) {
	query := url.Values{}
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	if q.Status != "" {
		query.Set("status", string(q.Status))
	}
	if q.SubCategory != "" {
		query.Set("sub_category", string(q.SubCategory))
	}
	if q.PaginationKey != "" {
		query.Set("pagination_key", q.PaginationKey)
	}

	var out []StockInstrument
	if err := c.get(ctx, pathStockInstruments, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}
