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

// pathStockFootprints is the stock footprint endpoint.
//
// Reference: https://developer.webull.hk/apis/docs/reference/footprint.md
// (the documented v3 path is /market-data/stocks/footprints/list; the payload
// is returned as a top-level array rather than a {data, ...} envelope).
const pathStockFootprints = "/market-data/stocks/footprints/list"

// FootprintTimespan is the bar granularity of a footprint request.
type FootprintTimespan string

// Footprint granularities accepted by the footprint endpoint. Only these values
// are supported; the endpoint does not accept arbitrary intervals.
const (
	// FootprintTimespanS5 is a five-second bar.
	FootprintTimespanS5 FootprintTimespan = "S5"
	// FootprintTimespanS15 is a fifteen-second bar.
	FootprintTimespanS15 FootprintTimespan = "S15"
	// FootprintTimespanM1 is a one-minute bar.
	FootprintTimespanM1 FootprintTimespan = "M1"
	// FootprintTimespanM5 is a five-minute bar.
	FootprintTimespanM5 FootprintTimespan = "M5"
	// FootprintTimespanM30 is a thirty-minute bar.
	FootprintTimespanM30 FootprintTimespan = "M30"
)

// TradingSession identifies a portion of the trading day. It is shared by the
// market-data endpoints that accept a trading_sessions parameter.
type TradingSession string

// Trading sessions accepted by the market-data endpoints. The footprint
// endpoint does not accept [TradingSessionOvernight].
const (
	// TradingSessionPre is the pre-market session.
	TradingSessionPre TradingSession = "PRE"
	// TradingSessionRTH is the regular trading hours session.
	TradingSessionRTH TradingSession = "RTH"
	// TradingSessionAfter is the after-hours session.
	TradingSessionAfter TradingSession = "ATH"
	// TradingSessionOvernight is the overnight session.
	TradingSessionOvernight TradingSession = "OVN"
)

// FootprintQuery parameterizes [Client.GetFootprint]. Symbols, Category and
// Timespan are required; the remaining fields are optional.
type FootprintQuery struct {
	// Symbols are the security symbols to query, at most 20 per request.
	// Required.
	Symbols []string
	// Category is the security type. Required, and only [StockCategoryUS] is
	// supported.
	Category StockCategory
	// Timespan is the bar granularity. Required.
	Timespan FootprintTimespan
	// Count is the number of bars to return, between 1 and 1200. Zero means
	// the server default of 200.
	Count int
	// RealTimeRequired reports whether bars that are not yet finalized should
	// be included. It only applies to minute timespans.
	RealTimeRequired bool
	// TradingSessions restricts the result to one trading session. Empty means
	// all sessions. The endpoint does not accept [TradingSessionOvernight].
	TradingSessions TradingSession
}

// StockFootprint is the footprint (order-flow) chart for a single symbol.
type StockFootprint struct {
	// Symbol is the security symbol, for example "AAPL".
	Symbol string `json:"symbol"`
	// InstrumentID is the unique identifier of the security.
	InstrumentID string `json:"instrument_id"`
	// Result is the ordered footprint bars.
	Result []FootprintBar `json:"result"`
}

// FootprintBar is a single footprint bar: aggregated buy and sell volume with a
// per-price breakdown.
type FootprintBar struct {
	// Time is the bar timestamp, for example
	// "2025-09-30T05:47:00.000+0000".
	Time string `json:"time"`
	// TradingSession is the session the bar belongs to.
	TradingSession TradingSession `json:"trading_session"`
	// Total is the sum of buy and sell volume, as a decimal string.
	Total string `json:"total"`
	// Delta is buy volume minus sell volume, as a decimal string.
	Delta string `json:"delta"`
	// BuyTotal is the buy-initiated volume, as a decimal string.
	BuyTotal string `json:"buy_total"`
	// SellTotal is the sell-initiated volume, as a decimal string.
	SellTotal string `json:"sell_total"`
	// BuyDetail maps price levels to buy volume at that price; both the price
	// and the volume are decimal strings.
	BuyDetail map[string]string `json:"buy_detail"`
	// SellDetail maps price levels to sell volume at that price; both the price
	// and the volume are decimal strings.
	SellDetail map[string]string `json:"sell_detail"`
}

// GetFootprint retrieves footprint (order-flow) bars for one or more US stocks.
// Symbols, Category and Timespan are required.
//
// Reference: https://developer.webull.hk/apis/docs/reference/footprint.md
func (c *Client) GetFootprint(ctx context.Context, q FootprintQuery) ([]StockFootprint, error) {
	query := url.Values{}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	if q.Timespan != "" {
		query.Set("timespan", string(q.Timespan))
	}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}
	query.Set("real_time_required", strconv.FormatBool(q.RealTimeRequired))
	if q.TradingSessions != "" {
		query.Set("trading_sessions", string(q.TradingSessions))
	}

	var out []StockFootprint
	if err := c.get(ctx, pathStockFootprints, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}
