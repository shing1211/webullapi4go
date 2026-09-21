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
	"strconv"
	"strings"
)

// Display Solution quote endpoints.
//
// TODO(ds): Confirm exact paths against US sandbox.
const (
	// pathDSSnapshot is the Display Solution snapshot endpoint.
	pathDSSnapshot = "/openapi/market-data/stock/snapshot"
	// pathDSBars is the Display Solution batch-bars endpoint (POST).
	pathDSBars = "/openapi/market-data/stock/batch-bars"
	// pathDSBarsSingle is the Display Solution single-bar endpoint (GET).
	pathDSBarsSingle = "/openapi/market-data/stock/bars"
	// pathDSTick is the Display Solution tick endpoint.
	pathDSTick = "/openapi/market-data/stock/tick"
	// pathDSDepth is the Display Solution depth/quotes endpoint.
	pathDSDepth = "/openapi/market-data/stock/quotes"
)

// GetDisplaySnapshot retrieves real-time market snapshots via Display Solution.
func (c *Client) GetDisplaySnapshot(ctx context.Context, q SnapshotQuery) ([]Snapshot, error) {
	query := url.Values{}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	query.Set("extend_hour_required", strconv.FormatBool(q.ExtendHourRequired))
	query.Set("overnight_required", strconv.FormatBool(q.OvernightRequired))

	var out []Snapshot
	if err := c.DisplayService().Get(ctx, pathDSSnapshot, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDisplayBars retrieves historical bars for multiple symbols via Display
// Solution. The request is sent as a POST with a JSON body.
func (c *Client) GetDisplayBars(ctx context.Context, q BatchBarQuery) (*BatchBars, error) {
	body := batchBarsRequest{
		Symbols:          q.Symbols,
		Category:         q.Category,
		Timespan:         q.Timespan,
		Count:            q.Count,
		RealTimeRequired: q.RealTimeRequired,
		StartTime:        q.StartTime,
		EndTime:          q.EndTime,
	}
	if len(q.TradingSessions) > 0 {
		body.TradingSessions = joinTradingSessions(q.TradingSessions)
	}

	var out BatchBars
	if err := c.DisplayService().Do(ctx, http.MethodPost, pathDSBars, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetDisplayBarsSingle retrieves historical bars for a single symbol via
// Display Solution. This is a GET request against the single-symbol bars path.
func (c *Client) GetDisplayBarsSingle(ctx context.Context, q BarQuery) (*StockBars, error) {
	out := &StockBars{Symbol: q.Symbol, Result: []Bar{}}
	if q.Symbol == "" {
		return out, nil
	}
	query := url.Values{}
	query.Set("symbol", q.Symbol)
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	if q.Interval != "" {
		query.Set("timespan", string(q.Interval))
	}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}

	if err := c.DisplayService().Get(ctx, pathDSBarsSingle, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDisplayTick retrieves tick-by-tick trade data via Display Solution.
func (c *Client) GetDisplayTick(ctx context.Context, q TickQuery) (*StockTicks, error) {
	query := url.Values{}
	if q.Symbol != "" {
		query.Set("symbol", q.Symbol)
	}
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}
	if len(q.TradingSessions) > 0 {
		query.Set("trading_sessions", joinTradingSessions(q.TradingSessions))
	}

	var out StockTicks
	if err := c.DisplayService().Get(ctx, pathDSTick, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetDisplayDepth retrieves the latest bid/ask order-book depth via Display
// Solution.
func (c *Client) GetDisplayDepth(ctx context.Context, q DepthQuery) (*Quote, error) {
	query := url.Values{}
	if q.Symbol != "" {
		query.Set("symbol", q.Symbol)
	}
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	if q.Depth > 0 {
		query.Set("depth", strconv.Itoa(q.Depth))
	}
	query.Set("overnight_required", strconv.FormatBool(q.OvernightRequired))

	var out Quote
	if err := c.DisplayService().Get(ctx, pathDSDepth, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
