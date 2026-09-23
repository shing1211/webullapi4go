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

// Crypto market-data and instrument endpoints. Documented on the US site.
const (
	pathCryptoBars        = "/market-data/crypto/bars/list"
	pathCryptoSnapshots   = "/market-data/crypto/snapshots/list"
	pathCryptoInstruments = "/trading/instruments/crypto/profiles/list"
)

// CryptoBarsQuery parameterizes [Client.GetCryptoBars].
type CryptoBarsQuery struct {
	// Symbols is the list of crypto symbols to query. Required.
	Symbols []string
	// Category is the crypto market category. Required.
	Category string
	// Timespan is the bar granularity. Required.
	Timespan BarTimespan
	// Count is the number of bars to return. Zero means the server default.
	Count int
	// RealTimeRequired asks the server to include the latest in-progress bar.
	RealTimeRequired bool
}

// CryptoBar is a single crypto price bar.
type CryptoBar struct {
	Time   string      `json:"time"`
	Open   money.Money `json:"open"`
	Close  money.Money `json:"close"`
	High   money.Money `json:"high"`
	Low    money.Money `json:"low"`
	Volume string      `json:"volume"`
}

// CryptoSymbolBars is the historical bars of one crypto symbol.
type CryptoSymbolBars struct {
	Symbol       string      `json:"symbol"`
	InstrumentID string      `json:"instrument_id"`
	Result       []CryptoBar `json:"result"`
}

// GetCryptoBars retrieves historical bars for one or more crypto symbols.
func (c *Client) GetCryptoBars(ctx context.Context, q CryptoBarsQuery) ([]CryptoSymbolBars, error) {
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
	query.Set("real_time_required", strconv.FormatBool(q.RealTimeRequired))

	var out []CryptoSymbolBars
	if err := c.get(ctx, pathCryptoBars, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CryptoSnapshotQuery parameterizes [Client.GetCryptoSnapshot].
type CryptoSnapshotQuery struct {
	// Symbols is the list of crypto symbols to query. Required.
	Symbols []string
	// Category is the crypto market category. Required.
	Category string
}

// CryptoSnapshot is the real-time market snapshot of one crypto symbol.
type CryptoSnapshot struct {
	InstrumentID  string      `json:"instrument_id"`
	Symbol        string      `json:"symbol"`
	PreClose      money.Money `json:"pre_close"`
	LastTradeTime int64       `json:"last_trade_time"`
	Price         money.Money `json:"price"`
	Open          money.Money `json:"open"`
	High          money.Money `json:"high"`
	Low           money.Money `json:"low"`
	Change        money.Money `json:"change"`
	ChangeRatio   string      `json:"change_ratio"`
	QuoteTime     string      `json:"quote_time"`
	Bid           money.Money `json:"bid"`
	BidSize       string      `json:"bid_size"`
	Ask           money.Money `json:"ask"`
	AskSize       string      `json:"ask_size"`
}

// GetCryptoSnapshot retrieves real-time snapshots for one or more crypto symbols.
func (c *Client) GetCryptoSnapshot(ctx context.Context, q CryptoSnapshotQuery) ([]CryptoSnapshot, error) {
	query := url.Values{}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	if q.Category != "" {
		query.Set("category", q.Category)
	}
	var out []CryptoSnapshot
	if err := c.get(ctx, pathCryptoSnapshots, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CryptoInstrumentQuery parameterizes [Client.GetCryptoInstruments].
type CryptoInstrumentQuery struct {
	// Category is the crypto market category. Required.
	Category string
	// Symbols optionally restricts the result to specific symbols.
	Symbols []string
	// Status optionally filters by instrument status.
	Status string
	// PaginationKey continues from a previous page. Empty means unset.
	PaginationKey string
}

// CryptoInstrument is the profile of one crypto instrument.
type CryptoInstrument struct {
	Symbol       string `json:"symbol"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	Currency     string `json:"currency"`
	Status       string `json:"status"`
	InstrumentID string `json:"instrument_id"`
}

// cryptoInstrumentsResponse is the {data, pagination_key} envelope.
type cryptoInstrumentsResponse struct {
	Data          []CryptoInstrument `json:"data"`
	PaginationKey string             `json:"pagination_key"`
}

// CryptoInstrumentsResult is the result of [Client.GetCryptoInstruments].
type CryptoInstrumentsResult struct {
	Instruments   []CryptoInstrument
	PaginationKey string
}

// GetCryptoInstruments retrieves crypto instrument profiles.
func (c *Client) GetCryptoInstruments(ctx context.Context, q CryptoInstrumentQuery) (*CryptoInstrumentsResult, error) {
	query := url.Values{}
	if q.Category != "" {
		query.Set("category", q.Category)
	}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}
	if q.Status != "" {
		query.Set("status", q.Status)
	}
	if q.PaginationKey != "" {
		query.Set("pagination_key", q.PaginationKey)
	}
	var resp cryptoInstrumentsResponse
	if err := c.get(ctx, pathCryptoInstruments, query, &resp); err != nil {
		return nil, err
	}
	return &CryptoInstrumentsResult{Instruments: resp.Data, PaginationKey: resp.PaginationKey}, nil
}
