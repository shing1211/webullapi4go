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
	"net/http"
	"net/url"
	"strings"
)

// pathStockProfilesV3 is the Display Solution stock instruments endpoint (v3).
//
// Reference: https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-using-get
// Unlike the v2 path (/openapi/instrument/stock/list), this path uses POST
// and returns a {data, pagination_key} envelope instead of a direct array.
const pathStockProfilesV3 = "/market-data/instruments/stocks/profiles/list"

// StockProfilesV3Query parameterizes [Client.GetStockProfilesV3].
type StockProfilesV3Query struct {
	Symbols       []string
	Category      StockCategory
	SubCategory   StockSubCategory
	ExchangeCodes []string
	Status        InstrumentStatus
	PaginationKey string
}

// instrumentProfilesV3Resp is the response envelope for the v3 instruments endpoint.
type instrumentProfilesV3Resp struct {
	Data          []StockInstrument `json:"data"`
	PaginationKey string            `json:"pagination_key"`
}

// StockProfilesV3Result is the result of [Client.GetStockProfilesV3].
type StockProfilesV3Result struct {
	Instruments   []StockInstrument
	PaginationKey string
}

// GetStockProfilesV3 retrieves profile information for one or more stock instruments
// using the Display Solution (v3) endpoint. Unlike the v2 endpoint
// ([Client.GetStockInstruments]), this uses POST and returns a paginated
// {data, pagination_key} envelope.
//
// Reference: https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/list-using-get
func (c *Client) GetStockProfilesV3(ctx context.Context, q StockProfilesV3Query) (*StockProfilesV3Result, error) {
	query := url.Values{}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	if q.SubCategory != "" {
		query.Set("sub_category", string(q.SubCategory))
	}
	if len(q.ExchangeCodes) > 0 {
		query.Set("exchange_codes", strings.Join(q.ExchangeCodes, ","))
	}
	if q.Status != "" {
		query.Set("status", string(q.Status))
	}
	if q.PaginationKey != "" {
		query.Set("pagination_key", q.PaginationKey)
	}

	var resp instrumentProfilesV3Resp
	if err := c.do(ctx, http.MethodPost, pathStockProfilesV3, query, nil, &resp); err != nil {
		return nil, err
	}
	return &StockProfilesV3Result{
		Instruments:   resp.Data,
		PaginationKey: resp.PaginationKey,
	}, nil
}
