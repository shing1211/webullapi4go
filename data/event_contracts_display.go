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

	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

// Event-contract instrument and market-data endpoints (US only).
const (
	pathEventTags           = "/market-data/instruments/event-contracts/categories/tags/list"
	pathEventEventsList     = "/market-data/instruments/event-contracts/events/list"
	pathEventMilestones     = "/market-data/instruments/event-contracts/milestones/list"
	pathEventSeriesList     = "/market-data/instruments/event-contracts/series/list"
	pathEventSportsFilters  = "/market-data/instruments/event-contracts/sports-filters/list"
	pathEventGameStats      = "/market-data/event-contracts/game-stats/get"
	pathEventLiveData       = "/market-data/event-contracts/live-data/get"
	pathEventMarketBars     = "/market-data/event-contracts/markets/bars/list"
	pathEventMarketBarsEvt  = "/market-data/event-contracts/markets/bars/list-by-event"
	pathEventMarketDepths   = "/market-data/event-contracts/markets/depths/list"
	pathEventMarketSnapshot = "/market-data/event-contracts/markets/snapshots/list"
)

// EventContractTag groups the tags available for a category.
type EventContractTag struct {
	Tags         []string `json:"tags"`
	CategoryID   int      `json:"category_id"`
	CategoryName string   `json:"category_name"`
	CategoryCode string   `json:"category_code"`
}

// GetEventContractTags lists the tags per event-contract category.
func (c *Client) GetEventContractTags(ctx context.Context) ([]EventContractTag, error) {
	var out []EventContractTag
	if err := c.get(ctx, pathEventTags, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// EventListQuery parameterizes [Client.GetEventContractEventsList].
type EventListQuery struct {
	SeriesSymbol  string
	Status        string
	PaginationKey string
}

// EventListPage is a paginated list of event items.
type EventListPage struct {
	Data          []map[string]any `json:"data"`
	PaginationKey string           `json:"pagination_key"`
}

// GetEventContractEventsList lists event-contract events.
func (c *Client) GetEventContractEventsList(ctx context.Context, q EventListQuery) (*EventListPage, error) {
	query := url.Values{}
	if q.SeriesSymbol != "" {
		query.Set("series_symbol", q.SeriesSymbol)
	}
	if q.Status != "" {
		query.Set("status", q.Status)
	}
	if q.PaginationKey != "" {
		query.Set("pagination_key", q.PaginationKey)
	}
	var out EventListPage
	if err := c.get(ctx, pathEventEventsList, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MilestoneQuery parameterizes [Client.GetEventContractMilestones].
type MilestoneQuery struct {
	MinimumStartDate   string
	Category           string
	Competition        string
	RelatedEventSymbol string
	PaginationKey      string
}

// MilestonePage is a paginated list of milestones.
type MilestonePage struct {
	Data          []map[string]any `json:"data"`
	PaginationKey string           `json:"pagination_key"`
}

// GetEventContractMilestones lists event-contract milestones.
func (c *Client) GetEventContractMilestones(ctx context.Context, q MilestoneQuery) (*MilestonePage, error) {
	query := url.Values{}
	if q.MinimumStartDate != "" {
		query.Set("minimum_start_date", q.MinimumStartDate)
	}
	if q.Category != "" {
		query.Set("category", q.Category)
	}
	if q.Competition != "" {
		query.Set("competition", q.Competition)
	}
	if q.RelatedEventSymbol != "" {
		query.Set("related_event_symbol", q.RelatedEventSymbol)
	}
	if q.PaginationKey != "" {
		query.Set("pagination_key", q.PaginationKey)
	}
	var out MilestonePage
	if err := c.get(ctx, pathEventMilestones, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EventSeriesListQuery parameterizes [Client.GetEventContractSeriesList].
type EventSeriesListQuery struct {
	Category      string
	Tags          string
	Symbols       []string
	PaginationKey string
}

// EventSeriesListPage is a paginated list of event series.
type EventSeriesListPage struct {
	Data          []map[string]any `json:"data"`
	PaginationKey string           `json:"pagination_key"`
}

// GetEventContractSeriesList lists event-contract series.
func (c *Client) GetEventContractSeriesList(ctx context.Context, q EventSeriesListQuery) (*EventSeriesListPage, error) {
	query := url.Values{}
	if q.Category != "" {
		query.Set("category", q.Category)
	}
	if q.Tags != "" {
		query.Set("tags", q.Tags)
	}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	if q.PaginationKey != "" {
		query.Set("pagination_key", q.PaginationKey)
	}
	var out EventSeriesListPage
	if err := c.get(ctx, pathEventSeriesList, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EventSportsFilter is a sports filter for event contracts.
type EventSportsFilter struct {
	Tag          string   `json:"tag"`
	Competitions []string `json:"competitions"`
	Scopes       []string `json:"scopes"`
}

// GetEventContractSportsFilters lists the sports filters.
func (c *Client) GetEventContractSportsFilters(ctx context.Context, tag string) ([]EventSportsFilter, error) {
	query := url.Values{}
	if tag != "" {
		query.Set("tag", tag)
	}
	var out []EventSportsFilter
	if err := c.get(ctx, pathEventSportsFilters, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// EventGameStats is the game statistics for an event-contract milestone.
type EventGameStats struct {
	MilestoneID string           `json:"milestone_id"`
	Periods     []map[string]any `json:"periods"`
}

// GetEventGameStats retrieves game statistics for a milestone.
func (c *Client) GetEventGameStats(ctx context.Context, milestoneID, category string) (*EventGameStats, error) {
	query := url.Values{}
	if milestoneID != "" {
		query.Set("milestone_id", milestoneID)
	}
	if category != "" {
		query.Set("category", category)
	}
	var out EventGameStats
	if err := c.get(ctx, pathEventGameStats, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EventLiveData is the live state of an event-contract milestone.
type EventLiveData struct {
	Type          string         `json:"type"`
	MilestoneID   string         `json:"milestone_id"`
	Status        string         `json:"status"`
	Winner        string         `json:"winner"`
	LastPlay      map[string]any `json:"last_play"`
	LastUpdatedTS int64          `json:"last_updated_ts"`
	Details       map[string]any `json:"details"`
}

// GetEventLiveData retrieves live data for a milestone.
func (c *Client) GetEventLiveData(ctx context.Context, milestoneID, category string) (*EventLiveData, error) {
	query := url.Values{}
	if milestoneID != "" {
		query.Set("milestone_id", milestoneID)
	}
	if category != "" {
		query.Set("category", category)
	}
	var out EventLiveData
	if err := c.get(ctx, pathEventLiveData, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EventMarketBar is a single event-contract market bar.
type EventMarketBar struct {
	EndPeriodTime string      `json:"end_period_time"`
	Volume        string      `json:"volume"`
	Open          money.Money `json:"open"`
	High          money.Money `json:"high"`
	Low           money.Money `json:"low"`
	Close         money.Money `json:"close"`
}

// EventMarketBarsQuery parameterizes [Client.GetEventMarketBars].
type EventMarketBarsQuery struct {
	Symbols   []string
	Category  string
	StartTime int64
	EndTime   int64
	Count     int
	Timespan  string
}

// GetEventMarketBars retrieves bars for event-contract market symbols.
func (c *Client) GetEventMarketBars(ctx context.Context, q EventMarketBarsQuery) ([]EventMarketBar, error) {
	query := url.Values{}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	if q.Category != "" {
		query.Set("category", q.Category)
	}
	if q.StartTime > 0 {
		query.Set("start_time", strconv.FormatInt(q.StartTime, 10))
	}
	if q.EndTime > 0 {
		query.Set("end_time", strconv.FormatInt(q.EndTime, 10))
	}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}
	if q.Timespan != "" {
		query.Set("timespan", q.Timespan)
	}
	var out []EventMarketBar
	if err := c.get(ctx, pathEventMarketBars, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// EventMarketBarsByEventQuery parameterizes [Client.GetEventMarketBarsByEvent].
type EventMarketBarsByEventQuery struct {
	EventSymbol string
	Category    string
	StartTime   int64
	EndTime     int64
	Count       int
	Timespan    string
}

// GetEventMarketBarsByEvent retrieves bars for an event symbol.
func (c *Client) GetEventMarketBarsByEvent(ctx context.Context, q EventMarketBarsByEventQuery) ([]EventMarketBar, error) {
	query := url.Values{}
	if q.EventSymbol != "" {
		query.Set("event_symbol", q.EventSymbol)
	}
	if q.Category != "" {
		query.Set("category", q.Category)
	}
	if q.StartTime > 0 {
		query.Set("start_time", strconv.FormatInt(q.StartTime, 10))
	}
	if q.EndTime > 0 {
		query.Set("end_time", strconv.FormatInt(q.EndTime, 10))
	}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}
	if q.Timespan != "" {
		query.Set("timespan", q.Timespan)
	}
	var out []EventMarketBar
	if err := c.get(ctx, pathEventMarketBarsEvt, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// EventMarketDepth is the order-book depth of an event-contract market.
type EventMarketDepth struct {
	Symbol       string           `json:"symbol"`
	InstrumentID string           `json:"instrument_id"`
	YesAsks      []map[string]any `json:"yes_asks"`
	YesBids      []map[string]any `json:"yes_bids"`
	NoAsks       []map[string]any `json:"no_asks"`
	NoBids       []map[string]any `json:"no_bids"`
}

// GetEventMarketDepth retrieves the depth of an event-contract market.
func (c *Client) GetEventMarketDepth(ctx context.Context, symbol, category string, depth int) (*EventMarketDepth, error) {
	query := url.Values{}
	if symbol != "" {
		query.Set("symbol", symbol)
	}
	if category != "" {
		query.Set("category", category)
	}
	if depth > 0 {
		query.Set("depth", strconv.Itoa(depth))
	}
	var out EventMarketDepth
	if err := c.get(ctx, pathEventMarketDepths, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EventMarketSnapshot is the market snapshot of an event-contract market.
type EventMarketSnapshot struct {
	Symbol        string      `json:"symbol"`
	InstrumentID  string      `json:"instrument_id"`
	EventSymbol   string      `json:"event_symbol"`
	YesSubTitle   string      `json:"yes_sub_title"`
	NoSubTitle    string      `json:"no_sub_title"`
	Status        string      `json:"status"`
	YesBid        money.Money `json:"yes_bid"`
	YesAsk        money.Money `json:"yes_ask"`
	NoBid         money.Money `json:"no_bid"`
	NoAsk         money.Money `json:"no_ask"`
	Price         money.Money `json:"price"`
	Volume        string      `json:"volume"`
	OpenInterest  string      `json:"open_interest"`
	LastTradeTime string      `json:"last_trade_time"`
}

// GetEventMarketSnapshot retrieves the snapshot of an event-contract market.
func (c *Client) GetEventMarketSnapshot(ctx context.Context, symbol, category string) (*EventMarketSnapshot, error) {
	query := url.Values{}
	if symbol != "" {
		query.Set("symbol", symbol)
	}
	if category != "" {
		query.Set("category", category)
	}
	var out EventMarketSnapshot
	if err := c.get(ctx, pathEventMarketSnapshot, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
