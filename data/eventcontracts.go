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

	"github.com/shing1211/webullapi4go/pkg/domain/money"
	"github.com/shing1211/webullapi4go/pkg/types"
)

// Event contract instrument discovery endpoints.
//
// Reference: https://developer.webull.com/apis/docs/reference/event-categories-list.md
const (
	pathEventContractCategories = "/trading/instruments/event-contracts/categories/list"
	pathEventContractSeries     = "/trading/instruments/event-contracts/series/list"
	pathEventContractEvents     = "/trading/instruments/event-contracts/events/list"
	pathEventContractMarkets    = "/trading/instruments/event-contracts/markets/list"
)

// EventContractCategory represents an event contract category.
//
// The page documents category_id, category_code and category_name, and marks all
// three required. Category and Name are the SDK's own spellings and are carried
// alongside them rather than replaced: a response that honours the documentation
// previously left all three required values unreadable, with nothing in the return
// path to distinguish that from a genuinely empty category. Exactly one naming set is
// populated per response.
type EventContractCategory struct {
	// Category is the SDK's own spelling, kept for callers that already read it.
	Category string `json:"category"`
	// Name is the SDK's own spelling, kept for callers that already read it.
	Name string `json:"name"`
	// CategoryID is the documented numeric identifier.
	CategoryID int64 `json:"category_id"`
	// CategoryCode is the documented short code, for example ECONOMICS.
	CategoryCode string `json:"category_code"`
	// CategoryName is the documented display name.
	CategoryName string `json:"category_name"`
}

// EventContractSeries represents a series of related events.
type EventContractSeries struct {
	SeriesSymbol string `json:"series_symbol"`
	Category     string `json:"category"`
	Name         string `json:"name"`
	Status       string `json:"status"`
}

// EventContractEvent represents a single event within a series.
//
// The page documents symbol, series_id, name, status, short_name, strike_date,
// strike_period and mutually_exclusive, and requires six of them. EventSymbol and
// SeriesSymbol are the SDK's own spellings for what the page calls symbol and
// series_id; both sets are carried so a caller reading either name gets a value, and
// so no call that works today breaks. Exactly one naming set is populated per
// response.
type EventContractEvent struct {
	// EventSymbol is the SDK's own spelling of the documented symbol.
	EventSymbol string `json:"event_symbol"`
	// SeriesSymbol is the SDK's own spelling of the documented series_id.
	SeriesSymbol string `json:"series_symbol"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	// Symbol is the documented event symbol.
	Symbol string `json:"symbol"`
	// SeriesID is the documented series identifier.
	SeriesID string `json:"series_id"`
	// ShortName is the documented abbreviated name.
	ShortName string `json:"short_name"`
	// MutuallyExclusive reports whether the event is mutually exclusive within its
	// series, as the page documents it.
	MutuallyExclusive bool `json:"mutually_exclusive"`
	// StrikeDate and StrikePeriod are documented but not required; they are omitted
	// when the server does not send them.
	StrikeDate   string `json:"strike_date,omitempty"`
	StrikePeriod string `json:"strike_period,omitempty"`
}

// EventContractMarket represents a tradable event contract instrument.
type EventContractMarket struct {
	Symbol         string      `json:"symbol"`
	EventSymbol    string      `json:"event_symbol"`
	SeriesSymbol   string      `json:"series_symbol"`
	Status         string      `json:"status"`
	StrikePrice    money.Money `json:"strike_price,omitempty"`
	ExpirationDate string      `json:"expiration_date,omitempty"`
}

// EventContractSeriesQuery parameterizes [Client.GetEventContractSeries].
type EventContractSeriesQuery struct {
	Category      string
	Symbols       []string
	PaginationKey string
}

// EventContractEventsQuery parameterizes [Client.GetEventContractEvents].
type EventContractEventsQuery struct {
	SeriesSymbol string
	Symbols      []string
	Status       string
}

// EventContractMarketsQuery parameterizes [Client.GetEventContractMarkets].
type EventContractMarketsQuery struct {
	SeriesSymbol        string
	EventSymbol         string
	Symbols             []string
	ExpirationDateAfter string
	PaginationKey       string
}

// GetEventContractCategories retrieves the list of event contract categories.
//
// Reference: https://developer.webull.com/apis/docs/reference/event-categories-list.md
func (c *Client) GetEventContractCategories(ctx context.Context) ([]EventContractCategory, error) {
	var out []EventContractCategory
	if err := c.get(ctx, pathEventContractCategories, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetEventContractSeries retrieves event contract series, optionally filtered
// by category and symbols.
//
// **Breaking, in v2.1.35.** The page documents the 200 body as
// `{"data": [...], "pagination_key": "..."}`, which a bare slice cannot
// decode, so this method failed outright against a conforming server. It
// returns a [types.Page] now: a caller reads out.Data instead of out, and
// passes out.PaginationKey back to fetch the following page. The new form
// succeeds where the old one could not.
//
// Reference: https://developer.webull.com/apis/docs/reference/event-categories-list.md
func (c *Client) GetEventContractSeries(ctx context.Context, q EventContractSeriesQuery) (*types.Page[EventContractSeries], error) {
	query := make(url.Values)
	if q.Category != "" {
		query.Set("category", q.Category)
	}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	if q.PaginationKey != "" {
		query.Set("pagination_key", q.PaginationKey)
	}
	var out types.Page[EventContractSeries]
	if err := c.get(ctx, pathEventContractSeries, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetEventContractEvents retrieves events within a series, optionally filtered
// by symbols and status.
//
// Reference: https://developer.webull.com/apis/docs/reference/event-categories-list.md
func (c *Client) GetEventContractEvents(ctx context.Context, q EventContractEventsQuery) ([]EventContractEvent, error) {
	query := make(url.Values)
	if q.SeriesSymbol != "" {
		query.Set("series_symbol", q.SeriesSymbol)
	}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	if q.Status != "" {
		query.Set("status", q.Status)
	}
	var out []EventContractEvent
	if err := c.get(ctx, pathEventContractEvents, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetEventContractMarkets retrieves tradable event contract markets, optionally
// filtered by series, event, symbols, and expiration date.
//
// **Breaking, in v2.1.35.** The page documents the 200 body as
// `{"data": [...], "pagination_key": "..."}`, which a bare slice cannot
// decode, so this method failed outright against a conforming server. It
// returns a [types.Page] now: a caller reads out.Data instead of out, and
// passes out.PaginationKey back to fetch the following page. The new form
// succeeds where the old one could not.
//
// Reference: https://developer.webull.com/apis/docs/reference/event-categories-list.md
func (c *Client) GetEventContractMarkets(ctx context.Context, q EventContractMarketsQuery) (*types.Page[EventContractMarket], error) {
	query := make(url.Values)
	if q.SeriesSymbol != "" {
		query.Set("series_symbol", q.SeriesSymbol)
	}
	if q.EventSymbol != "" {
		query.Set("event_symbol", q.EventSymbol)
	}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	if q.ExpirationDateAfter != "" {
		query.Set("expiration_date_after", q.ExpirationDateAfter)
	}
	if q.PaginationKey != "" {
		query.Set("pagination_key", q.PaginationKey)
	}
	var out types.Page[EventContractMarket]
	if err := c.get(ctx, pathEventContractMarkets, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
