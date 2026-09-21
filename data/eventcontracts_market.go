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

// Event contract market data endpoints.
//
// TODO(event-market-data): Confirm host and paths against US sandbox.
const (
	pathEventSnapshot = "/market-data/event-contracts/snapshots/list"
	pathEventDepth    = "/market-data/event-contracts/depths/list"
	pathEventBars     = "/market-data/event-contracts/bars/list"
	pathEventTick     = "/market-data/event-contracts/ticks/list"
)

// EventSnapshotQuery parameterizes [Client.GetEventSnapshot]. Symbols and
// Category are required.
type EventSnapshotQuery struct {
	// Symbols is the list of event-contract symbols to query, at most 100.
	Symbols []string
	// Category is the market to query. Required.
	Category string
}

// EventSnapshot is a real-time snapshot for an event contract.
type EventSnapshot struct {
	// Symbol is the event-contract symbol.
	Symbol string `json:"symbol"`
	// LastPrice is the last traded price, as a decimal string.
	LastPrice string `json:"last_price"`
	// YesBid is the best bid price for the yes side, as a decimal string.
	YesBid string `json:"yes_bid"`
	// YesAsk is the best ask price for the yes side, as a decimal string.
	YesAsk string `json:"yes_ask"`
	// Volume is the traded volume, as a decimal string.
	Volume string `json:"volume"`
	// OpenInterest is the open interest, as a decimal string.
	OpenInterest string `json:"open_interest"`
	// Timestamp is the snapshot time, as a string.
	Timestamp string `json:"timestamp"`
	// Extra holds additional fields not mapped to the struct.
	Extra map[string]string `json:"-"`
}

// GetEventSnapshot retrieves real-time market snapshots for one or more event
// contract symbols.
//
// TODO(event-market-data): Confirm host and paths against US sandbox.
func (c *Client) GetEventSnapshot(ctx context.Context, q EventSnapshotQuery) ([]EventSnapshot, error) {
	query := url.Values{}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	if q.Category != "" {
		query.Set("category", q.Category)
	}

	var out []EventSnapshot
	if err := c.get(ctx, pathEventSnapshot, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// EventDepthQuery parameterizes [Client.GetEventDepth]. Symbol, Category and
// Depth are required.
type EventDepthQuery struct {
	// Symbol is the event-contract symbol.
	Symbol string
	// Category is the market to query. Required.
	Category string
	// Depth is the number of order-book levels to return.
	Depth int
}

// DepthLevel is a single price level in the event-contract order book.
type DepthLevel struct {
	// Price is the level price, as a decimal string.
	Price string `json:"price"`
	// Size is the aggregate quantity at the level, as a decimal string.
	Size string `json:"size"`
}

// EventDepth is the order book for an event contract.
type EventDepth struct {
	// Symbol is the event-contract symbol.
	Symbol string `json:"symbol"`
	// Timestamp is the depth time, as a string.
	Timestamp string `json:"timestamp"`
	// YesBids is the bid side of the book, best (highest) price first.
	YesBids []DepthLevel `json:"yes_bids"`
	// YesAsks is the ask side of the book, best (lowest) price first.
	YesAsks []DepthLevel `json:"yes_asks"`
}

// GetEventDepth retrieves the bid/ask order-book depth for a single event
// contract symbol.
//
// TODO(event-market-data): Confirm host and paths against US sandbox.
func (c *Client) GetEventDepth(ctx context.Context, q EventDepthQuery) (*EventDepth, error) {
	query := url.Values{}
	if q.Symbol != "" {
		query.Set("symbol", q.Symbol)
	}
	if q.Category != "" {
		query.Set("category", q.Category)
	}
	if q.Depth > 0 {
		query.Set("depth", strconv.Itoa(q.Depth))
	}

	var out EventDepth
	if err := c.get(ctx, pathEventDepth, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EventBarsQuery parameterizes [Client.GetEventBars]. Symbols, Category and
// Timespan are required.
type EventBarsQuery struct {
	// Symbols is the list of event-contract symbols to query.
	Symbols []string
	// Category is the market to query. Required.
	Category string
	// Timespan is the bar granularity. Required.
	Timespan BarTimespan
	// Count is the number of bars to return per symbol. Zero means the server
	// default.
	Count int
	// RealTimeRequired requests the latest market data when true. When nil the
	// server default applies.
	RealTimeRequired *bool
}

// EventBar is a candlestick bar for an event contract.
type EventBar struct {
	// Symbol is the event-contract symbol.
	Symbol string `json:"symbol"`
	// Open is the open price, as a decimal string.
	Open string `json:"open"`
	// High is the high price, as a decimal string.
	High string `json:"high"`
	// Low is the low price, as a decimal string.
	Low string `json:"low"`
	// Close is the close price, as a decimal string.
	Close string `json:"close"`
	// Volume is the volume, as a decimal string.
	Volume string `json:"volume"`
	// Timestamp is the bar time, as a string.
	Timestamp string `json:"timestamp"`
	// Timespan is the bar granularity.
	Timespan BarTimespan `json:"timespan"`
}

// GetEventBars retrieves historical bars for one or more event contract
// symbols.
//
// TODO(event-market-data): Confirm host and paths against US sandbox.
func (c *Client) GetEventBars(ctx context.Context, q EventBarsQuery) ([]EventBar, error) {
	query := url.Values{}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	if q.Category != "" {
		query.Set("category", q.Category)
	}
	if q.Timespan != "" {
		query.Set("timespan", string(q.Timespan))
	}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}
	if q.RealTimeRequired != nil {
		query.Set("real_time_required", strconv.FormatBool(*q.RealTimeRequired))
	}

	var out []EventBar
	if err := c.get(ctx, pathEventBars, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// EventTickQuery parameterizes [Client.GetEventTick]. Symbol, Category and
// Count are required.
type EventTickQuery struct {
	// Symbol is the event-contract symbol.
	Symbol string
	// Category is the market to query. Required.
	Category string
	// Count is the number of ticks to return.
	Count int
}

// EventTick is a tick trade for an event contract.
type EventTick struct {
	// Symbol is the event-contract symbol.
	Symbol string `json:"symbol"`
	// YesPrice is the yes-side trade price, as a decimal string.
	YesPrice string `json:"yes_price"`
	// NoPrice is the no-side trade price, as a decimal string.
	NoPrice string `json:"no_price"`
	// Side is the aggressor side.
	Side string `json:"side"`
	// Volume is the executed trade volume, as a decimal string.
	Volume string `json:"volume"`
	// TradeID is the unique trade identifier.
	TradeID string `json:"trade_id"`
	// Timestamp is the trade time, as a string.
	Timestamp string `json:"timestamp"`
}

// GetEventTick retrieves tick-by-tick trade data for a single event contract
// symbol.
//
// TODO(event-market-data): Confirm host and paths against US sandbox.
func (c *Client) GetEventTick(ctx context.Context, q EventTickQuery) ([]EventTick, error) {
	query := url.Values{}
	if q.Symbol != "" {
		query.Set("symbol", q.Symbol)
	}
	if q.Category != "" {
		query.Set("category", q.Category)
	}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}

	var out []EventTick
	if err := c.get(ctx, pathEventTick, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}
