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

// pathStockSnapshots is the stock snapshot endpoint.
//
// Reference: https://developer.webull.hk/apis/docs/reference/snapshot.md
// The documented current path is /market-data/stocks/snapshots/list; the legacy
// /openapi/... alias is not used.
const pathStockSnapshots = "/market-data/stocks/snapshots/list"

// SnapshotQuery parameterizes [Client.GetSnapshot]. Symbols and Category are
// required; the remaining fields are optional.
type SnapshotQuery struct {
	// Symbols is the list of security symbols to query, at most 100.
	Symbols []string
	// Category is the market to query. Required. Besides the
	// [StockCategoryUS], [StockCategoryHK] and [StockCategoryCN] values it
	// accepts the snapshot-specific StockCategory("US_ETF").
	Category StockCategory
	// ExtendHourRequired includes pre-market and after-hours trading data when
	// true.
	ExtendHourRequired bool
	// OvernightRequired includes overnight trading data when true.
	OvernightRequired bool
}

// Snapshot is a real-time market snapshot for a single security.
type Snapshot struct {
	// InstrumentID is the unique identifier of the security.
	InstrumentID string `json:"instrument_id"`
	// PreClose is the previous close price, as a decimal string.
	PreClose string `json:"pre_close"`
	// ChangeRatio is the price change ratio, as a decimal string.
	ChangeRatio string `json:"change_ratio"`
	// Symbol is the trading symbol, for example "AAPL".
	Symbol string `json:"symbol"`
	// LastTradeTime is the last trade time, as a Unix timestamp in
	// milliseconds.
	LastTradeTime int64 `json:"last_trade_time"`
	// Price is the current price, as a decimal string.
	Price string `json:"price"`
	// Open is the intraday open price, as a decimal string. For US stocks this
	// excludes pre/post-market data and is unset when no trading occurred.
	Open string `json:"open"`
	// Close is the intraday close price, as a decimal string.
	Close string `json:"close"`
	// High is the intraday high price, as a decimal string.
	High string `json:"high"`
	// Low is the intraday low price, as a decimal string.
	Low string `json:"low"`
	// Volume is the traded volume, as a decimal string.
	Volume string `json:"volume"`
	// Change is the price change amount, as a decimal string.
	Change string `json:"change"`
	// Ask is the best ask price, as a decimal string.
	Ask string `json:"ask"`
	// AskSize is the best ask size, as a decimal string.
	AskSize string `json:"ask_size"`
	// Bid is the best bid price, as a decimal string.
	Bid string `json:"bid"`
	// BidSize is the best bid size, as a decimal string.
	BidSize string `json:"bid_size"`
	// Turnover is the turnover rate, as a decimal string.
	Turnover string `json:"turnover"`
	// EPS is the earnings per share, as a decimal string.
	EPS string `json:"eps"`
	// EPSTTM is the trailing-twelve-month earnings per share, as a decimal
	// string.
	EPSTTM string `json:"eps_ttm"`
	// LotSize is the number of shares per lot, as a decimal string.
	LotSize string `json:"lot_size"`
	// BPS is the book value per share, as a decimal string.
	BPS string `json:"bps"`
	// ExtendHourLastPrice is the pre/post-market latest price, as a decimal
	// string.
	ExtendHourLastPrice string `json:"extend_hour_last_price"`
	// ExtendHourHigh is the pre/post-market high price, as a decimal string.
	ExtendHourHigh string `json:"extend_hour_high"`
	// ExtendHourLow is the pre/post-market low price, as a decimal string.
	ExtendHourLow string `json:"extend_hour_low"`
	// ExtendHourChange is the pre/post-market change amount, as a decimal
	// string.
	ExtendHourChange string `json:"extend_hour_change"`
	// ExtendHourChangeRatio is the pre/post-market change ratio, as a decimal
	// string.
	ExtendHourChangeRatio string `json:"extend_hour_change_ratio"`
	// ExtendHourVolume is the pre/post-market volume, as a decimal string.
	ExtendHourVolume string `json:"extend_hour_volume"`
	// ExtendHourLastTradeTime is the pre/post-market last trade time, as a Unix
	// timestamp in milliseconds.
	ExtendHourLastTradeTime int64 `json:"extend_hour_last_trade_time"`
	// OvnPrice is the overnight price, as a decimal string.
	OvnPrice string `json:"ovn_price"`
	// OvnHigh is the overnight high price, as a decimal string.
	OvnHigh string `json:"ovn_high"`
	// OvnLow is the overnight low price, as a decimal string.
	OvnLow string `json:"ovn_low"`
	// OvnVolume is the overnight volume, as a decimal string.
	OvnVolume string `json:"ovn_volume"`
	// OvnChange is the overnight change amount, as a decimal string.
	OvnChange string `json:"ovn_change"`
	// OvnChangeRatio is the overnight change ratio, as a decimal string.
	OvnChangeRatio string `json:"ovn_change_ratio"`
	// OvnLastTradeTime is the overnight last trade time, as a Unix timestamp in
	// milliseconds.
	OvnLastTradeTime int64 `json:"ovn_last_trade_time"`
	// OvnAsk is the overnight best ask price, as a decimal string.
	OvnAsk string `json:"ovn_ask"`
	// OvnAskSize is the overnight best ask size, as a decimal string.
	OvnAskSize string `json:"ovn_ask_size"`
	// OvnBid is the overnight best bid price, as a decimal string.
	OvnBid string `json:"ovn_bid"`
	// OvnBidSize is the overnight best bid size, as a decimal string.
	OvnBidSize string `json:"ovn_bid_size"`
}

// GetSnapshot retrieves real-time market snapshots for one or more symbols.
//
// Reference: https://developer.webull.hk/apis/docs/reference/snapshot.md
func (c *Client) GetSnapshot(ctx context.Context, q SnapshotQuery) ([]Snapshot, error) {
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
	if err := c.get(ctx, pathStockSnapshots, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}
