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

// pathCorpActionsList is the Display Solution endpoint for corporate actions by symbol.
//
// Reference: https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/corp-action-using-get
const pathCorpActionsList = "/market-data/instruments/stocks/corporate-actions/list"

// pathCorpActionsMarket is the Display Solution endpoint for corporate actions by market.
//
// Reference: https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/corp-market-using-get
const pathCorpActionsMarket = "/market-data/instruments/stocks/corporate-actions/list-by-market"

// CorporateActionQuery parameterizes [Client.GetCorporateActions].
type CorporateActionQuery struct {
	Symbols       []string
	Market        string
	StartDate     string
	EndDate       string
	EventTypes    []string
	PageSize      int
	PaginationKey string
}

// corpActionResponse is the raw API response wrapper.
type corpActionResponse struct {
	Data          []CorporateAction `json:"data"`
	PaginationKey string            `json:"pagination_key"`
}

// CorporateAction represents a corporate action event such as a dividend or stock split.
//
// Reference: https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/corp-action-using-get
type CorporateAction struct {
	InstrumentID int64  `json:"instrument_id"`
	Symbol       string `json:"symbol"`
	ExchangeCode string `json:"exchange_code"`
	EventType    string `json:"event_type"`
	EventAction  string `json:"event_action"`
	EventID      int64  `json:"event_id"`
	Source       string `json:"source"`
	RatioOld     string `json:"ratio_old"`
	RatioNew     string `json:"ratio_new"`
	EventDate    string `json:"event_date"`
	UpdateTime   string `json:"update_time"`
	CreateTime   string `json:"create_time"`
}

// GetCorporateActions retrieves corporate action events for one or more symbols
// using the Display Solution API.
//
// category is required and must be one of: "US_STOCK", "HK_STOCK", "CN_STOCK".
//
// Reference: https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/corp-action-using-get
func (c *Client) GetCorporateActions(ctx context.Context, q CorporateActionQuery) ([]CorporateAction, string, error) {
	query := url.Values{}
	if len(q.Symbols) > 0 {
		query.Set("symbol", q.Symbols[0])
	}
	if q.Market != "" {
		query.Set("category", q.Market)
	}
	if q.StartDate != "" {
		query.Set("start_date", q.StartDate)
	}
	if q.EndDate != "" {
		query.Set("end_date", q.EndDate)
	}
	if len(q.EventTypes) > 0 {
		query.Set("event_types", joinStrings(q.EventTypes, ","))
	}
	if q.PageSize > 0 {
		query.Set("page_size", itoa(q.PageSize))
	}
	if q.PaginationKey != "" {
		query.Set("pagination_key", q.PaginationKey)
	}

	var out corpActionResponse
	if err := c.DisplayService().Get(ctx, pathCorpActionsList, query, &out); err != nil {
		return nil, "", err
	}
	return out.Data, out.PaginationKey, nil
}

// GetCorporateActionsByMarket retrieves corporate action events for all securities
// in a market using the Display Solution API.
//
// Currently only "US" market is supported.
//
// Reference: https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/corp-market-using-get
func (c *Client) GetCorporateActionsByMarket(ctx context.Context, q CorporateActionQuery) ([]CorporateAction, string, error) {
	query := url.Values{}
	if q.Market != "" {
		query.Set("market", q.Market)
	}
	if q.StartDate != "" {
		query.Set("start_date", q.StartDate)
	}
	if q.EndDate != "" {
		query.Set("end_date", q.EndDate)
	}
	if len(q.EventTypes) > 0 {
		query.Set("event_types", joinStrings(q.EventTypes, ","))
	}
	if q.PageSize > 0 {
		query.Set("page_size", itoa(q.PageSize))
	}
	if q.PaginationKey != "" {
		query.Set("pagination_key", q.PaginationKey)
	}

	var out corpActionResponse
	if err := c.DisplayService().Get(ctx, pathCorpActionsMarket, query, &out); err != nil {
		return nil, "", err
	}
	return out.Data, out.PaginationKey, nil
}

func joinStrings(vals []string, sep string) string {
	if len(vals) == 0 {
		return ""
	}
	if len(vals) == 1 {
		return vals[0]
	}
	result := vals[0]
	for _, v := range vals[1:] {
		result += sep + v
	}
	return result
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	if i < 0 {
		return "-" + uitoa(uint(-i))
	}
	return uitoa(uint(i))
}

func uitoa(val uint) string {
	if val == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for val > 0 {
		i--
		buf[i] = byte('0' + val%10)
		val /= 10
	}
	return string(buf[i:])
}
