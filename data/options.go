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

// Option market-data endpoints.
//
// References:
//   - https://developer.webull.hk/apis/docs/reference/option-tick.md
//   - https://developer.webull.hk/apis/docs/reference/option-snapshot.md
//   - https://developer.webull.hk/apis/docs/reference/option-historical-bars.md
const (
	pathOptionTicks     = "/market-data/options/ticks/list"
	pathOptionSnapshots = "/market-data/options/snapshots/list"
	pathOptionBars      = "/market-data/options/bars/list"
)

// OptionCategory identifies the option market. The option endpoints currently
// support US options only.
type OptionCategory string

// Option market categories.
const (
	// OptionCategoryUS identifies United States options, the only category
	// currently supported by the option endpoints.
	OptionCategoryUS OptionCategory = "US_OPTION"
)

// OptionBarTimespan is the time granularity of option historical bars.
type OptionBarTimespan string

// Option bar timespans.
const (
	// OptionBarTimespanM1 is a one-minute bar.
	OptionBarTimespanM1 OptionBarTimespan = "M1"
	// OptionBarTimespanM5 is a five-minute bar.
	OptionBarTimespanM5 OptionBarTimespan = "M5"
	// OptionBarTimespanM15 is a fifteen-minute bar.
	OptionBarTimespanM15 OptionBarTimespan = "M15"
	// OptionBarTimespanM30 is a thirty-minute bar.
	OptionBarTimespanM30 OptionBarTimespan = "M30"
	// OptionBarTimespanM60 is a sixty-minute bar.
	OptionBarTimespanM60 OptionBarTimespan = "M60"
	// OptionBarTimespanM120 is a two-hour bar.
	OptionBarTimespanM120 OptionBarTimespan = "M120"
	// OptionBarTimespanM240 is a four-hour bar.
	OptionBarTimespanM240 OptionBarTimespan = "M240"
	// OptionBarTimespanD is a daily bar.
	OptionBarTimespanD OptionBarTimespan = "D"
	// OptionBarTimespanW is a weekly bar.
	OptionBarTimespanW OptionBarTimespan = "W"
	// OptionBarTimespanM is a monthly bar.
	OptionBarTimespanM OptionBarTimespan = "M"
	// OptionBarTimespanY is a yearly bar.
	OptionBarTimespanY OptionBarTimespan = "Y"
)

// OptionTickQuery parameterizes [Client.GetOptionTick]. Symbol is required;
// Category defaults to [OptionCategoryUS] and Count defaults to the server's
// value when zero.
type OptionTickQuery struct {
	// Symbol is the option contract symbol, for example
	// "AAPL260522C00300000". Required.
	Symbol string
	// Category is the option market. Empty means [OptionCategoryUS].
	Category OptionCategory
	// Count is the number of ticks to return, at most 1200. Zero means unset.
	Count int
}

// OptionTick is a single executed option trade.
type OptionTick struct {
	// Time is the trade time as a Unix epoch millisecond timestamp string.
	Time string `json:"time"`
	// Price is the executed trade price, as a decimal string.
	Price string `json:"price"`
	// Volume is the executed trade volume, as a string.
	Volume string `json:"volume"`
	// Side is the aggressor side, for example "B" or "S".
	Side string `json:"side"`
}

// OptionTickResult is the tick history of one option contract.
type OptionTickResult struct {
	// Symbol is the option contract symbol.
	Symbol string `json:"symbol"`
	// InstrumentID is the unique identifier of the option contract.
	InstrumentID string `json:"instrument_id"`
	// Result lists the trade ticks.
	Result []OptionTick `json:"result"`
}

// OptionSnapshotQuery parameterizes [Client.GetOptionSnapshot]. Symbols is
// required and is limited to 20 symbols per query; Category defaults to
// [OptionCategoryUS].
type OptionSnapshotQuery struct {
	// Symbols is the option contract symbols to query. Required.
	Symbols []string
	// Category is the option market. Empty means [OptionCategoryUS].
	Category OptionCategory
}

// OptionSnapshot is the real-time market snapshot of one option contract.
type OptionSnapshot struct {
	// InstrumentID is the unique identifier of the option contract.
	InstrumentID string `json:"instrument_id"`
	// Symbol is the option contract symbol.
	Symbol string `json:"symbol"`
	// Price is the last traded price, as a decimal string.
	Price string `json:"price"`
	// Open is the session open price, as a decimal string.
	Open string `json:"open"`
	// High is the session high price, as a decimal string.
	High string `json:"high"`
	// Low is the session low price, as a decimal string.
	Low string `json:"low"`
	// PreClose is the previous settlement or close price, as a decimal string.
	PreClose string `json:"pre_close"`
	// Volume is the accumulated session volume, as a string.
	Volume string `json:"volume"`
	// Change is the absolute change from PreClose, as a decimal string.
	Change string `json:"change"`
	// ChangeRatio is the change relative to PreClose, as a decimal string.
	ChangeRatio string `json:"change_ratio"`
	// LastTradeTime is the last trade time as a Unix epoch millisecond
	// timestamp.
	LastTradeTime int64 `json:"last_trade_time"`
	// Close is the close price, as a decimal string.
	Close string `json:"close"`
	// StrikePrice is the option strike price, as a decimal string.
	StrikePrice string `json:"strike_price"`
	// Gamma is the option gamma, as a decimal string.
	Gamma string `json:"gamma"`
	// Delta is the option delta, as a decimal string.
	Delta string `json:"delta"`
	// Rho is the option rho, as a decimal string.
	Rho string `json:"rho"`
	// Theta is the option theta, as a decimal string.
	Theta string `json:"theta"`
	// Vega is the option vega, as a decimal string.
	Vega string `json:"vega"`
	// ImpVol is the implied volatility, as a decimal string.
	ImpVol string `json:"imp_vol"`
	// OpenInterest is the open interest, as a string.
	OpenInterest string `json:"open_interest"`
	// QuoteTime is the quote time as a Unix epoch millisecond timestamp.
	QuoteTime int64 `json:"quote_time"`
	// Bid is the best bid price, as a decimal string.
	Bid string `json:"bid"`
	// Ask is the best ask price, as a decimal string.
	Ask string `json:"ask"`
	// AskSize is the quantity available at the best ask, as a string.
	AskSize string `json:"ask_size"`
	// BidSize is the quantity available at the best bid, as a string.
	BidSize string `json:"bid_size"`
	// DealAmount is the accumulated traded value, as a decimal string.
	DealAmount string `json:"deal_amount"`
}

// OptionBarsQuery parameterizes [Client.GetOptionBars]. Symbols, Category, and
// Timespan are required; Count defaults to the server's value when zero.
type OptionBarsQuery struct {
	// Symbols is the option contract symbols to query, at most 20 per query.
	// Required.
	Symbols []string
	// Category is the option market. Empty means [OptionCategoryUS].
	Category OptionCategory
	// Timespan is the bar granularity. Required.
	Timespan OptionBarTimespan
	// Count is the number of bars to return, at most 1200. Zero means unset.
	Count int
	// RealTimeRequired, when true, asks the server to include the latest
	// in-progress bar.
	RealTimeRequired bool
}

// OptionBar is a single historical option price bar.
type OptionBar struct {
	// Time is the bar time in ISO 8601 form.
	Time string `json:"time"`
	// Open is the open price, as a decimal string.
	Open string `json:"open"`
	// Close is the close price, as a decimal string.
	Close string `json:"close"`
	// High is the high price, as a decimal string.
	High string `json:"high"`
	// Low is the low price, as a decimal string.
	Low string `json:"low"`
	// Volume is the bar volume, as a string.
	Volume string `json:"volume"`
}

// OptionSymbolBars is the historical bars of one option contract.
type OptionSymbolBars struct {
	// Symbol is the option contract symbol.
	Symbol string `json:"symbol"`
	// InstrumentID is the unique identifier of the option contract.
	InstrumentID string `json:"instrument_id"`
	// Result lists the historical bars.
	Result []OptionBar `json:"result"`
}

// optionBarsResponse is the {result: [...]} envelope returned by the option
// bars endpoint.
type optionBarsResponse struct {
	Result []OptionSymbolBars `json:"result"`
}

// optionCategory returns the category to send, defaulting to
// [OptionCategoryUS] when the caller left it unset.
func optionCategory(category OptionCategory) OptionCategory {
	if category == "" {
		return OptionCategoryUS
	}
	return category
}

// GetOptionTick retrieves tick-by-tick trade data for one option contract.
//
// Reference: https://developer.webull.hk/apis/docs/reference/option-tick.md
func (c *Client) GetOptionTick(ctx context.Context, q OptionTickQuery) (*OptionTickResult, error) {
	query := url.Values{
		"symbol":   {q.Symbol},
		"category": {string(optionCategory(q.Category))},
	}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}
	var out OptionTickResult
	if err := c.get(ctx, pathOptionTicks, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetOptionSnapshot retrieves real-time snapshots for up to 20 option
// contracts.
//
// Reference: https://developer.webull.hk/apis/docs/reference/option-snapshot.md
func (c *Client) GetOptionSnapshot(ctx context.Context, q OptionSnapshotQuery) ([]OptionSnapshot, error) {
	query := url.Values{
		"category": {string(optionCategory(q.Category))},
	}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	var out []OptionSnapshot
	if err := c.get(ctx, pathOptionSnapshots, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetOptionBars retrieves historical bars for one or more option contracts.
//
// Reference: https://developer.webull.hk/apis/docs/reference/option-historical-bars.md
func (c *Client) GetOptionBars(ctx context.Context, q OptionBarsQuery) ([]OptionSymbolBars, error) {
	query := url.Values{
		"category": {string(optionCategory(q.Category))},
	}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	if q.Timespan != "" {
		query.Set("timespan", string(q.Timespan))
	}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}
	if q.RealTimeRequired {
		query.Set("real_time_required", strconv.FormatBool(true))
	}
	var out optionBarsResponse
	if err := c.get(ctx, pathOptionBars, query, &out); err != nil {
		return nil, err
	}
	return out.Result, nil
}
