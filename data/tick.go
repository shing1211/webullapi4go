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

// pathStockTicks is the stock tick-by-tick endpoint.
//
// Reference: https://developer.webull.hk/apis/docs/reference/tick.md
// The documented current path is /market-data/stocks/ticks/list; the legacy
// /openapi/... alias is not used.
const pathStockTicks = "/market-data/stocks/ticks/list"

// TickQuery parameterizes [Client.GetTick]. Symbol and Category are required.
type TickQuery struct {
	// Symbol is the security symbol, for example "AAPL".
	Symbol string
	// Category is the market to query. Required.
	Category StockCategory
	// Count is the number of ticks to return. Zero means the server default
	// (30); the documented maximum is 1000.
	Count int
	// TradingSessions restricts the result to the given trading sessions.
	// When empty the server default applies.
	TradingSessions []TradingSession
}

// Tick is a single executed trade.
type Tick struct {
	// Time is the trade time, as a Unix timestamp in milliseconds encoded as a
	// string.
	Time string `json:"time"`
	// Price is the executed trade price, as a decimal string.
	Price money.Money `json:"price"`
	// Volume is the executed trade volume, as a decimal string.
	Volume string `json:"volume"`
	// Side is the aggressor side. Documented values include "B", "S", "G",
	// "L" and "N".
	Side string `json:"side"`
}

// StockTicks is the tick-by-tick response for a single security.
type StockTicks struct {
	// Symbol is the security symbol.
	Symbol string `json:"symbol"`
	// InstrumentID is the unique identifier of the security.
	InstrumentID string `json:"instrument_id"`
	// Result is the list of executed trades, newest first.
	Result []Tick `json:"result"`
}

// GetTick retrieves tick-by-tick trade data for a single symbol.
//
// Reference: https://developer.webull.hk/apis/docs/reference/tick.md
func (c *Client) GetTick(ctx context.Context, q TickQuery) (*StockTicks, error) {
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
	if err := c.get(ctx, pathStockTicks, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// joinTradingSessions encodes sessions as the comma-separated wire form used by
// the market-data endpoints.
func joinTradingSessions(sessions []TradingSession) string {
	parts := make([]string, len(sessions))
	for i, s := range sessions {
		parts[i] = string(s)
	}
	return strings.Join(parts, ",")
}
