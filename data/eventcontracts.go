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
type EventContractCategory struct {
	Category string `json:"category"`
	Name     string `json:"name"`
}

// EventContractSeries represents a series of related events.
type EventContractSeries struct {
	SeriesSymbol string `json:"series_symbol"`
	Category     string `json:"category"`
	Name         string `json:"name"`
	Status       string `json:"status"`
}

// EventContractEvent represents a single event within a series.
type EventContractEvent struct {
	EventSymbol  string `json:"event_symbol"`
	SeriesSymbol string `json:"series_symbol"`
	Name         string `json:"name"`
	Status       string `json:"status"`
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
// Reference: https://developer.webull.com/apis/docs/reference/event-categories-list.md
func (c *Client) GetEventContractSeries(ctx context.Context, q EventContractSeriesQuery) ([]EventContractSeries, error) {
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
	var out []EventContractSeries
	if err := c.get(ctx, pathEventContractSeries, query, &out); err != nil {
		return nil, err
	}
	return out, nil
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
// Reference: https://developer.webull.com/apis/docs/reference/event-categories-list.md
func (c *Client) GetEventContractMarkets(ctx context.Context, q EventContractMarketsQuery) ([]EventContractMarket, error) {
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
	var out []EventContractMarket
	if err := c.get(ctx, pathEventContractMarkets, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}
