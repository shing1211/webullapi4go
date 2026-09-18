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
)

// pathStockBarsList is the historical-bars endpoint.
//
// Reference: https://developer.webull.hk/apis/docs/reference/historical-bars.md
// The documented current path is /market-data/stocks/bars/list; the legacy
// /openapi/... alias is not used.
//
// The historical single-symbol endpoint GET /market-data/stocks/bars/get is
// retired: both the official Webull Python SDK (MarketData.get_history_bar) and
// the live sandbox report it as no longer available. [Client.GetBars] therefore
// issues a one-symbol batch request against this same path.
const pathStockBarsList = "/market-data/stocks/bars/list"

// BarTimespan is the time granularity of a bar. It is the "timespan" body field
// of the batch bars endpoint.
type BarTimespan string

// Bar time granularities accepted by the bars endpoints.
const (
	// BarTimespanS5 is a 5-second bar.
	BarTimespanS5 BarTimespan = "S5"
	// BarTimespanS15 is a 15-second bar.
	BarTimespanS15 BarTimespan = "S15"
	// BarTimespanM1 is a 1-minute bar.
	BarTimespanM1 BarTimespan = "M1"
	// BarTimespanM5 is a 5-minute bar.
	BarTimespanM5 BarTimespan = "M5"
	// BarTimespanM15 is a 15-minute bar.
	BarTimespanM15 BarTimespan = "M15"
	// BarTimespanM30 is a 30-minute bar.
	BarTimespanM30 BarTimespan = "M30"
	// BarTimespanM60 is a 60-minute bar.
	BarTimespanM60 BarTimespan = "M60"
	// BarTimespanM120 is a 120-minute bar.
	BarTimespanM120 BarTimespan = "M120"
	// BarTimespanM240 is a 240-minute bar.
	BarTimespanM240 BarTimespan = "M240"
	// BarTimespanDay is a daily bar.
	BarTimespanDay BarTimespan = "D"
	// BarTimespanWeek is a weekly bar.
	BarTimespanWeek BarTimespan = "W"
	// BarTimespanMonth is a monthly bar.
	BarTimespanMonth BarTimespan = "M"
	// BarTimespanYear is a yearly bar.
	BarTimespanYear BarTimespan = "Y"
)

// Bar is a single OHLCV candlestick. Trading sessions use the shared
// [TradingSession] type.
type Bar struct {
	// Time is the bar time, as returned by the server. The batch endpoint
	// documents an ISO-8601 string.
	Time string `json:"time"`
	// Open is the open price, as a decimal string.
	Open string `json:"open"`
	// Close is the close price, as a decimal string.
	Close string `json:"close"`
	// High is the high price, as a decimal string.
	High string `json:"high"`
	// Low is the low price, as a decimal string.
	Low string `json:"low"`
	// Volume is the volume, as a decimal string.
	Volume string `json:"volume"`
	// TradingSession is the session the bar belongs to, when supplied.
	TradingSession TradingSession `json:"trading_sessions,omitempty"`
}

// StockBars is the historical-bars response for a single symbol.
type StockBars struct {
	// Symbol is the security symbol.
	Symbol string `json:"symbol"`
	// InstrumentID is the unique identifier of the security.
	InstrumentID string `json:"instrument_id"`
	// Result is the list of bars.
	Result []Bar `json:"result"`
}

// BarQuery parameterizes [Client.GetBars]. Symbol, Category and Interval are
// required.
type BarQuery struct {
	// Symbol is the security symbol, for example "AAPL".
	Symbol string
	// Category is the market to query. Required.
	Category StockCategory
	// Interval is the bar granularity. Required.
	Interval BarTimespan
	// Count is the number of bars to return. Zero means the server default
	// (200); the documented maximum is 1200 (1650 for M1).
	Count int
	// RealTimeRequired requests the latest market data when true. When nil the
	// server default (true) applies.
	RealTimeRequired *bool
	// TradingSessions restricts the result to the given trading sessions.
	TradingSessions []TradingSession
}

// GetBars retrieves historical bars for a single symbol.
//
// The historical single-symbol endpoint GET /market-data/stocks/bars/get is
// retired: both the official Webull Python SDK (MarketData.get_history_bar) and
// the live sandbox report it as no longer available. GetBars therefore issues a
// one-symbol batch request against POST /market-data/stocks/bars/list and
// unwraps the single result.
//
// Reference: https://developer.webull.hk/apis/docs/reference/historical-bars.md
func (c *Client) GetBars(ctx context.Context, q BarQuery) (*StockBars, error) {
	out := &StockBars{Symbol: q.Symbol, Result: []Bar{}}
	if q.Symbol == "" {
		return out, nil
	}
	batch, err := c.GetBatchBars(ctx, BatchBarQuery{
		Symbols:          []string{q.Symbol},
		Category:         q.Category,
		Timespan:         q.Interval,
		Count:            q.Count,
		RealTimeRequired: q.RealTimeRequired,
		TradingSessions:  q.TradingSessions,
	})
	if err != nil {
		return nil, err
	}
	for _, r := range batch.Result {
		if r.Symbol == q.Symbol {
			out.Symbol = r.Symbol
			out.InstrumentID = r.InstrumentID
			out.Result = r.Result
			break
		}
	}
	return out, nil
}

// BatchBarQuery parameterizes [Client.GetBatchBars]. Symbols, Category and
// Timespan are required.
type BatchBarQuery struct {
	// Symbols is the list of security symbols to query, at most 100.
	Symbols []string
	// Category is the market to query. Required.
	Category StockCategory
	// Timespan is the bar granularity. Required. The batch endpoint documents
	// this parameter as "timespan".
	Timespan BarTimespan
	// Count is the number of bars to return per symbol. Zero means the server
	// default (200); the documented maximum is 1200 (1650 for M1).
	Count int
	// RealTimeRequired requests the latest market data when true. When nil the
	// server default (true) applies.
	RealTimeRequired *bool
	// TradingSessions restricts the result to the given trading sessions.
	TradingSessions []TradingSession
	// StartTime restricts the result to bars at or after this Unix timestamp in
	// milliseconds. Zero means unbounded.
	StartTime int64
	// EndTime restricts the result to bars at or before this Unix timestamp in
	// milliseconds. Zero means unbounded.
	EndTime int64
}

// batchBarsRequest is the wire body of [Client.GetBatchBars].
type batchBarsRequest struct {
	Symbols          []string      `json:"symbols"`
	Category         StockCategory `json:"category"`
	Timespan         BarTimespan   `json:"timespan"`
	Count            int           `json:"count,omitempty"`
	RealTimeRequired *bool         `json:"real_time_required,omitempty"`
	TradingSessions  string        `json:"trading_sessions,omitempty"`
	StartTime        int64         `json:"start_time,omitempty"`
	EndTime          int64         `json:"end_time,omitempty"`
}

// BatchBarSymbol groups the bars returned for one symbol.
type BatchBarSymbol struct {
	// Symbol is the security symbol.
	Symbol string `json:"symbol"`
	// InstrumentID is the unique identifier of the security.
	InstrumentID string `json:"instrument_id"`
	// Result is the list of bars for Symbol.
	Result []Bar `json:"result"`
}

// BatchBars is the batch historical-bars response.
type BatchBars struct {
	// Result groups the bars by symbol.
	Result []BatchBarSymbol `json:"result"`
}

// GetBatchBars retrieves historical bars for multiple symbols in one request.
//
// Reference: https://developer.webull.hk/apis/docs/reference/historical-bars.md
func (c *Client) GetBatchBars(ctx context.Context, q BatchBarQuery) (*BatchBars, error) {
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
	if err := c.do(ctx, http.MethodPost, pathStockBarsList, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
