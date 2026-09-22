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

// Futures market data endpoints. Path convention follows the stock market
// data pattern (/market-data/{asset}/{resource}/list). These paths are
// inferred from the documented convention and are unconfirmed against the
// live API.
//
// TODO(futures): Confirm exact paths against US sandbox.
const (
	pathFuturesTick      = "/market-data/futures/ticks/list"
	pathFuturesSnapshot  = "/market-data/futures/snapshots/list"
	pathFuturesBars      = "/market-data/futures/bars/list"
	pathFuturesDepth     = "/market-data/futures/depths/list"
	pathFuturesFootprint = "/market-data/futures/footprints/list"
)

// FuturesTickQuery parameterizes [Client.GetFuturesTick]. Symbol is required.
type FuturesTickQuery struct {
	// Symbol is the futures contract symbol, for example "ESZ5".
	Symbol string
	// Category is the futures market category. Defaults to US_FUTURES if empty.
	Category FuturesCategory
	// Count is the number of ticks to return. Zero means the server default
	// (30); the documented maximum is 1000.
	Count int
}

// GetFuturesTick retrieves tick-by-tick trade data for a single futures
// contract.
//
// TODO(futures): Confirm path and query parameters against US sandbox.
func (c *Client) GetFuturesTick(ctx context.Context, q FuturesTickQuery) (*StockTicks, error) {
	query := url.Values{}
	if q.Symbol != "" {
		query.Set("symbol", q.Symbol)
	}
	category := q.Category
	if category == "" {
		category = FuturesCategoryUS
	}
	query.Set("category", string(category))
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}

	var out StockTicks
	if err := c.get(ctx, pathFuturesTick, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FuturesSnapshotQuery parameterizes [Client.GetFuturesSnapshot]. Symbols is
// required.
type FuturesSnapshotQuery struct {
	// Symbols is the list of futures contract symbols to query, at most 100.
	Symbols []string
	// Category is the futures market category. Defaults to US_FUTURES if empty.
	Category FuturesCategory
}

// GetFuturesSnapshot retrieves real-time market snapshots for one or more
// futures contracts.
//
// TODO(futures): Confirm path and query parameters against US sandbox.
func (c *Client) GetFuturesSnapshot(ctx context.Context, q FuturesSnapshotQuery) ([]Snapshot, error) {
	query := url.Values{}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	category := q.Category
	if category == "" {
		category = FuturesCategoryUS
	}
	query.Set("category", string(category))

	var out []Snapshot
	if err := c.get(ctx, pathFuturesSnapshot, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FuturesBarsQuery parameterizes [Client.GetFuturesBars]. Symbols and Interval
// are required.
type FuturesBarsQuery struct {
	// Symbols is the list of futures contract symbols to query.
	Symbols []string
	// Category is the futures market category. Defaults to US_FUTURES if empty.
	Category FuturesCategory
	// Interval is the bar granularity. Required.
	Interval BarTimespan
	// Count is the number of bars to return per symbol. Zero means the server
	// default (200); the documented maximum is 1200 (1650 for M1).
	Count int
}

// GetFuturesBars retrieves historical bars for one or more futures contracts.
//
// TODO(futures): Confirm path, query parameters, and response shape against
// US sandbox.
func (c *Client) GetFuturesBars(ctx context.Context, q FuturesBarsQuery) (*BatchBars, error) {
	query := url.Values{}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	category := q.Category
	if category == "" {
		category = FuturesCategoryUS
	}
	query.Set("category", string(category))
	if q.Interval != "" {
		query.Set("timespan", string(q.Interval))
	}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}

	var out BatchBars
	if err := c.get(ctx, pathFuturesBars, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FuturesDepthQuery parameterizes [Client.GetFuturesDepth]. Symbol is required.
type FuturesDepthQuery struct {
	// Symbol is the futures contract symbol, for example "ESZ5".
	Symbol string
	// Category is the futures market category. Defaults to US_FUTURES if empty.
	Category FuturesCategory
	// Depth is the number of order-book levels to return: 1 for L1, 10 for the
	// default L2 depth. Zero means the server default.
	Depth int
}

// GetFuturesDepth retrieves the latest bid/ask order-book depth for a single
// futures contract.
//
// TODO(futures): Confirm path and query parameters against US sandbox.
func (c *Client) GetFuturesDepth(ctx context.Context, q FuturesDepthQuery) (*Quote, error) {
	query := url.Values{}
	if q.Symbol != "" {
		query.Set("symbol", q.Symbol)
	}
	category := q.Category
	if category == "" {
		category = FuturesCategoryUS
	}
	query.Set("category", string(category))
	if q.Depth > 0 {
		query.Set("depth", strconv.Itoa(q.Depth))
	}

	var out Quote
	if err := c.get(ctx, pathFuturesDepth, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FuturesFootprintQuery parameterizes [Client.GetFuturesFootprint]. Symbols and
// Timespan are required.
type FuturesFootprintQuery struct {
	// Symbols are the futures contract symbols to query, at most 20 per request.
	Symbols []string
	// Category is the futures market category. Defaults to US_FUTURES if empty.
	Category FuturesCategory
	// Timespan is the bar granularity. Required.
	Timespan FootprintTimespan
	// Count is the number of bars to return, between 1 and 1200. Zero means
	// the server default of 200.
	Count int
}

// GetFuturesFootprint retrieves footprint (order-flow) bars for one or more
// futures contracts.
//
// TODO(futures): Confirm path, query parameters, and response shape against
// US sandbox.
func (c *Client) GetFuturesFootprint(ctx context.Context, q FuturesFootprintQuery) ([]StockFootprint, error) {
	query := url.Values{}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	category := q.Category
	if category == "" {
		category = FuturesCategoryUS
	}
	query.Set("category", string(category))
	if q.Timespan != "" {
		query.Set("timespan", string(q.Timespan))
	}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}

	var out []StockFootprint
	if err := c.get(ctx, pathFuturesFootprint, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// HKFuturesTickQuery parameterizes [Client.GetHKFuturesTick]. Symbol is required;
// Category is pre-set to [FuturesCategoryHK].
type HKFuturesTickQuery struct {
	Symbol string
	Count  int
}

// GetHKFuturesTick retrieves tick-by-tick trade data for a single Hong Kong
// futures contract.
//
// TODO(futures-hk): unconfirmed path — requires live probe
func (c *Client) GetHKFuturesTick(ctx context.Context, q HKFuturesTickQuery) (*StockTicks, error) {
	query := url.Values{
		"symbol":   {q.Symbol},
		"category": {string(FuturesCategoryHK)},
	}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}
	var out StockTicks
	if err := c.get(ctx, pathFuturesTick, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// HKFuturesSnapshotQuery parameterizes [Client.GetHKFuturesSnapshot]. Symbols is
// required; Category is pre-set to [FuturesCategoryHK].
type HKFuturesSnapshotQuery struct {
	Symbols []string
}

// GetHKFuturesSnapshot retrieves real-time market snapshots for one or more
// Hong Kong futures contracts.
//
// TODO(futures-hk): unconfirmed path — requires live probe
func (c *Client) GetHKFuturesSnapshot(ctx context.Context, q HKFuturesSnapshotQuery) ([]Snapshot, error) {
	query := url.Values{
		"symbols":  {strings.Join(q.Symbols, ",")},
		"category": {string(FuturesCategoryHK)},
	}
	var out []Snapshot
	if err := c.get(ctx, pathFuturesSnapshot, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// HKFuturesBarsQuery parameterizes [Client.GetHKFuturesBars]. Symbols and Interval
// are required; Category is pre-set to [FuturesCategoryHK].
type HKFuturesBarsQuery struct {
	Symbols  []string
	Interval BarTimespan
	Count    int
}

// GetHKFuturesBars retrieves historical bars for one or more Hong Kong futures
// contracts.
//
// TODO(futures-hk): unconfirmed path — requires live probe
func (c *Client) GetHKFuturesBars(ctx context.Context, q HKFuturesBarsQuery) (*BatchBars, error) {
	query := url.Values{
		"symbols":  {strings.Join(q.Symbols, ",")},
		"category": {string(FuturesCategoryHK)},
	}
	if q.Interval != "" {
		query.Set("timespan", string(q.Interval))
	}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}
	var out BatchBars
	if err := c.get(ctx, pathFuturesBars, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// HKFuturesDepthQuery parameterizes [Client.GetHKFuturesDepth]. Symbol is required;
// Category is pre-set to [FuturesCategoryHK].
type HKFuturesDepthQuery struct {
	Symbol string
	Depth  int
}

// GetHKFuturesDepth retrieves the latest bid/ask order-book depth for a single
// Hong Kong futures contract.
//
// TODO(futures-hk): unconfirmed path — requires live probe
func (c *Client) GetHKFuturesDepth(ctx context.Context, q HKFuturesDepthQuery) (*Quote, error) {
	query := url.Values{
		"symbol":   {q.Symbol},
		"category": {string(FuturesCategoryHK)},
	}
	if q.Depth > 0 {
		query.Set("depth", strconv.Itoa(q.Depth))
	}
	var out Quote
	if err := c.get(ctx, pathFuturesDepth, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// HKFuturesFootprintQuery parameterizes [Client.GetHKFuturesFootprint]. Symbols and
// Timespan are required; Category is pre-set to [FuturesCategoryHK].
type HKFuturesFootprintQuery struct {
	Symbols  []string
	Timespan FootprintTimespan
	Count    int
}

// GetHKFuturesFootprint retrieves footprint (order-flow) bars for one or more
// Hong Kong futures contracts.
//
// TODO(futures-hk): unconfirmed path — requires live probe
func (c *Client) GetHKFuturesFootprint(ctx context.Context, q HKFuturesFootprintQuery) ([]StockFootprint, error) {
	query := url.Values{
		"symbols":  {strings.Join(q.Symbols, ",")},
		"category": {string(FuturesCategoryHK)},
	}
	if q.Timespan != "" {
		query.Set("timespan", string(q.Timespan))
	}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}
	var out []StockFootprint
	if err := c.get(ctx, pathFuturesFootprint, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}
