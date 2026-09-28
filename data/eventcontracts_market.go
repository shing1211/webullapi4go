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

// Event contract market data endpoints.
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
//
// The page documents 15 required names and this type originally carried 7, several
// under names the page does not use: it has LastPrice where the page says price, and no
// size companion for any side, so a caller could not read the depth a snapshot exists to
// report. The documented names are carried **alongside** the existing ones rather than
// replacing them, because a required name is one a conforming server always sends and
// replacing a field would break a caller that reads it. Which set a live server sends
// is unverified, so both are present and exactly one is populated per response.
type EventSnapshot struct {
	// Symbol is the event-contract symbol.
	Symbol string `json:"symbol"`
	// LastPrice is the last traded price, as a decimal string. The page calls this
	// price; see Price.
	LastPrice money.Money `json:"last_price"`
	// YesBid is the best bid price for the yes side, as a decimal string.
	YesBid money.Money `json:"yes_bid"`
	// YesAsk is the best ask price for the yes side, as a decimal string.
	YesAsk money.Money `json:"yes_ask"`
	// Volume is the traded volume, as a decimal string.
	Volume string `json:"volume"`
	// OpenInterest is the open interest, as a decimal string.
	OpenInterest string `json:"open_interest"`
	// Timestamp is the snapshot time, as a string.
	Timestamp string `json:"timestamp"`

	// The fields below are the documented names. They are required by the page.

	// InstrumentID is the unique identifier of the instrument.
	InstrumentID string `json:"instrument_id"`
	// Name is the event's display name.
	Name string `json:"name"`
	// Price is the documented name for the snapshot price, as a decimal string.
	Price money.Money `json:"price"`
	// LastTradeTime is the last trade time, as a Unix epoch millisecond timestamp.
	LastTradeTime int64 `json:"last_trade_time"`
	// YesBidSize is the size resting at the best yes bid.
	YesBidSize string `json:"yes_bid_size"`
	// YesAskSize is the size resting at the best yes ask.
	YesAskSize string `json:"yes_ask_size"`
	// NoBid is the best bid price for the no side, as a decimal string.
	NoBid money.Money `json:"no_bid"`
	// NoBidSize is the size resting at the best no bid.
	NoBidSize string `json:"no_bid_size"`
	// NoAsk is the best ask price for the no side, as a decimal string.
	NoAsk money.Money `json:"no_ask"`
	// NoAskSize is the size resting at the best no ask.
	NoAskSize string `json:"no_ask_size"`
}

// GetEventSnapshot retrieves real-time market snapshots for one or more event
// contract symbols.
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
	Price money.Money `json:"price"`
	// Size is the aggregate quantity at the level, as a decimal string.
	Size string `json:"size"`
}

// EventDepth is the order book for one event contract instrument.
//
// The page documents the 200 body as an **array** of these, each requiring
// instrument_id, symbol, quote_time, yes_bids, yes_asks, no_bids and no_asks, so
// [Client.GetEventDepth] returns a slice of them. The no-side book was absent here
// entirely, and the depth time was named timestamp rather than the documented
// quote_time, so a response honouring the documentation left both unreadable.
// QuoteTime is carried alongside Timestamp because a required name is one a
// conforming server always sends and which of the two a live server sends is
// unverified.
type EventDepth struct {
	// Symbol is the event-contract symbol.
	Symbol string `json:"symbol"`
	// Timestamp is the SDK's own spelling of the depth time.
	Timestamp string `json:"timestamp"`
	// YesBids is the bid side of the book, best (highest) price first.
	YesBids []DepthLevel `json:"yes_bids"`
	// YesAsks is the ask side of the book, best (lowest) price first.
	YesAsks []DepthLevel `json:"yes_asks"`

	// The fields below are the documented names, all required by the page.

	// InstrumentID is the unique identifier of the instrument.
	InstrumentID string `json:"instrument_id"`
	// QuoteTime is the documented depth time, as a Unix epoch millisecond timestamp.
	QuoteTime int64 `json:"quote_time"`
	// NoBids is the no-side bid book, best (highest) price first. The SDK carried no
	// no-side book at all, so a caller could not read the no side of an event
	// contract's depth, which is the side such an instrument is named for.
	NoBids []DepthLevel `json:"no_bids"`
	// NoAsks is the no-side ask book, best (lowest) price first.
	NoAsks []DepthLevel `json:"no_asks"`
}

// GetEventDepth retrieves the bid/ask order-book depth for a single event
// contract symbol.
//
// **Breaking, in v2.1.32.** It now returns []EventDepth, because the page documents
// the 200 body as an array of one entry per instrument rather than a single object.
// That inverts which shape fails to decode: previously a conforming response could not
// be read, and now a single-object response cannot be. Neither shape has been verified
// against a live host.
func (c *Client) GetEventDepth(ctx context.Context, q EventDepthQuery) ([]EventDepth, error) {
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

	var out []EventDepth
	if err := c.get(ctx, pathEventDepth, query, &out); err != nil {
		return nil, err
	}
	return out, nil
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
	Open money.Money `json:"open"`
	// High is the high price, as a decimal string.
	High money.Money `json:"high"`
	// Low is the low price, as a decimal string.
	Low money.Money `json:"low"`
	// Close is the close price, as a decimal string.
	Close money.Money `json:"close"`
	// Volume is the volume, as a decimal string.
	Volume string `json:"volume"`
	// Timestamp is the bar time, as a string. The page calls this time; see Time.
	Timestamp string `json:"timestamp"`
	// Timespan is the bar granularity.
	Timespan BarTimespan `json:"timespan"`
	// Time is the documented name for the bar time, as a string. The page requires
	// time on every bar and does not mention timestamp, so a response honouring the
	// documentation left this type's own field empty.
	Time string `json:"time"`
}

// EventBarsResult is one element of the documented response to
// [Client.GetEventBars].
//
// The page documents the 200 body as an array of objects, each requiring
// instrument_id, result and symbol, where result is itself an array of bars. This
// type is that outer element. The SDK previously returned the inner bars directly,
// which discarded the grouping key: bars from several instruments arrived in one flat
// list with nothing to say which instrument each belonged to, and instrument_id and
// symbol were unreachable.
type EventBarsResult struct {
	// InstrumentID is the unique identifier of the instrument these bars belong to.
	InstrumentID string `json:"instrument_id"`
	// Symbol is the event-contract symbol these bars belong to.
	Symbol string `json:"symbol"`
	// Result is the bar series, in the order the page documents them.
	Result []EventBar `json:"result"`
}

// GetEventBars retrieves historical bars for one or more event contract
// symbols.
//
// **Breaking, in v2.1.32.** It now returns []EventBarsResult, one entry per
// instrument, because that is the shape the page documents. It previously returned
// []EventBar - the elements of the inner result array - which flattened several
// instruments' bars into one list and dropped the grouping key. A caller reading
// out.Result[i] instead of out[i] adapts in one line; a caller that used out[i] as a
// bar must read out[i].Result.
func (c *Client) GetEventBars(ctx context.Context, q EventBarsQuery) ([]EventBarsResult, error) {
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

	var out []EventBarsResult
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
	YesPrice money.Money `json:"yes_price"`
	// NoPrice is the no-side trade price, as a decimal string.
	NoPrice money.Money `json:"no_price"`
	// Side is the aggressor side.
	Side string `json:"side"`
	// Volume is the executed trade volume, as a decimal string.
	Volume string `json:"volume"`
	// TradeID is the unique trade identifier.
	TradeID string `json:"trade_id"`
	// Timestamp is the trade time, as a string. The page calls this time; see Time.
	Timestamp string `json:"timestamp"`
	// Time is the documented name for the trade time, as a string. The page requires
	// time on every tick and does not mention timestamp, so a response honouring the
	// documentation left this type's own field empty.
	Time string `json:"time"`
}

// EventTickResult is one element of the documented response to
// [Client.GetEventTick].
//
// The page documents the 200 body as an array of objects, each requiring
// instrument_id, result and symbol, where result is itself an array of ticks. This
// type is that outer element. The SDK previously returned the inner ticks directly,
// which discarded the grouping key, so ticks from several instruments arrived in one
// flat list with nothing to say which instrument each belonged to.
type EventTickResult struct {
	// InstrumentID is the unique identifier of the instrument these ticks belong to.
	InstrumentID string `json:"instrument_id"`
	// Symbol is the event-contract symbol these ticks belong to.
	Symbol string `json:"symbol"`
	// Result is the tick series, in the order the page documents them.
	Result []EventTick `json:"result"`
}

// GetEventTick retrieves tick-by-tick trade data for a single event contract
// symbol.
//
// **Breaking, in v2.1.32.** It now returns []EventTickResult, one entry per
// instrument, because that is the shape the page documents. It previously returned
// []EventTick - the elements of the inner result array. A caller reading
// out.Result[i] instead of out[i] adapts in one line.
func (c *Client) GetEventTick(ctx context.Context, q EventTickQuery) ([]EventTickResult, error) {
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

	var out []EventTickResult
	if err := c.get(ctx, pathEventTick, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}
