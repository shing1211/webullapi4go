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
	"strings"
)

// pathCorpActionsList is the corporate actions by-symbol endpoint.
//
// Reference: https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/corp-action-using-get
const pathCorpActionsList = "/market-data/instruments/stocks/corporate-actions/list"

// pathCorpActionsMarket is the corporate actions bulk-by-market endpoint.
//
// Reference: https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/corp-market-using-get
const pathCorpActionsMarket = "/market-data/instruments/stocks/corporate-actions/market"

// CorporateActionQuery parameterizes [Client.GetCorporateActions].
type CorporateActionQuery struct {
	Symbols       []string
	Market        string
	StartDate     string
	EndDate       string
	PageSize      int
	PaginationKey string
}

// CorporateAction represents a corporate action event such as a dividend or stock split.
// Field names and types are best-effort; live probe needed to confirm all fields.
type CorporateAction struct {
	// Symbol is the security symbol, for example "AAPL".
	Symbol string `json:"symbol"`
	// EventType is the type of corporate action, for example "DIVIDEND", "SPLIT".
	EventType string `json:"event_type"`
	// DeclarationDate is the announcement date, as a date string.
	DeclarationDate string `json:"declaration_date"`
	// ExDate is the ex-dividend or ex-split date.
	ExDate string `json:"ex_date"`
	// PayDate is the payment date.
	PayDate string `json:"pay_date"`
	// RecordDate is the record date for determining shareholders eligible for the action.
	RecordDate string `json:"record_date"`
	// Amount is the per-share amount for dividends, or the ratio for splits.
	Amount string `json:"amount"`
	// Currency is the currency of the amount, for example "USD".
	Currency string `json:"currency"`
	// Frequency is the dividend frequency, for example "QUARTERLY", "ANNUAL".
	Frequency string `json:"frequency"`
	// Fields not yet confirmed by live probe.
	Extra map[string]string `json:"-"`
}

// GetCorporateActions retrieves corporate action events for one or more symbols.
//
// Reference: https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/corp-action-using-get
func (c *Client) GetCorporateActions(ctx context.Context, q CorporateActionQuery) ([]CorporateAction, error) {
	query := url.Values{}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	if q.Market != "" {
		query.Set("market", q.Market)
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
	if q.PaginationKey != "" {
		query.Set("pagination_key", q.PaginationKey)
	}

	var out []CorporateAction
	if err := c.get(ctx, pathCorpActionsList, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCorporateActionsByMarket retrieves corporate action events for all securities in a market.
//
// Reference: https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/corp-market-using-get
func (c *Client) GetCorporateActionsByMarket(ctx context.Context, q CorporateActionQuery) ([]CorporateAction, error) {
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
	if q.PageSize > 0 {
		query.Set("page_size", strconv.Itoa(q.PageSize))
	}
	if q.PaginationKey != "" {
		query.Set("pagination_key", q.PaginationKey)
	}

	var out []CorporateAction
	if err := c.get(ctx, pathCorpActionsMarket, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}
