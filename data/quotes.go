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

	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

// pathStockDepths is the stock order-book depth endpoint.
//
// Reference: https://developer.webull.hk/apis/docs/reference/quotes.md
// The documented current path is /market-data/stocks/depths/list; the legacy
// /openapi/... alias is not used.
const pathStockDepths = "/market-data/stocks/depths/list"

// DepthQuery parameterizes [Client.GetQuotes]. Symbol and Category are
// required.
type DepthQuery struct {
	// Symbol is the security symbol, for example "AAPL".
	Symbol string
	// Category is the market to query. Required.
	Category StockCategory
	// Depth is the number of order-book levels to return: 1 for L1, 10 for the
	// default L2 depth. Zero means the server default.
	Depth int
	// OvernightRequired includes overnight trading data when true.
	OvernightRequired bool
}

// QuoteOrder is a single market-participant order contributing to a price
// level.
type QuoteOrder struct {
	// MPID is the market participant identifier, for example "NSDQ".
	MPID string `json:"mpid"`
	// Size is the quantity contributed by the participant, as a decimal string.
	Size string `json:"size"`
}

// QuoteBroker is a broker attribution for a price level.
type QuoteBroker struct {
	// Bid is the broker identifier.
	Bid string `json:"bid"`
	// Name is the broker display name.
	Name string `json:"name"`
}

// QuoteLevel is one side of the order book at a single price.
type QuoteLevel struct {
	// Price is the level price, as a decimal string.
	Price money.Money `json:"price"`
	// Size is the aggregate quantity at the level, as a decimal string.
	Size string `json:"size"`
	// Order lists the contributing market-participant orders.
	Order []QuoteOrder `json:"order"`
	// Broker lists the contributing brokers, when supplied.
	Broker []QuoteBroker `json:"broker"`
}

// Quote is the order-book depth for a single security.
type Quote struct {
	// Symbol is the security symbol.
	Symbol string `json:"symbol"`
	// InstrumentID is the unique identifier of the security.
	InstrumentID string `json:"instrument_id"`
	// QuoteTime is the quote time, as a Unix timestamp in milliseconds. The
	// sandbox returns it as a JSON number even though the reference documents
	// it as a string.
	QuoteTime int64 `json:"quote_time"`
	// Asks is the ask side of the book, best (lowest) price first.
	Asks []QuoteLevel `json:"asks"`
	// Bids is the bid side of the book, best (highest) price first.
	Bids []QuoteLevel `json:"bids"`
}

// GetQuotes retrieves the latest bid/ask order-book depth for a single symbol.
//
// Reference: https://developer.webull.hk/apis/docs/reference/quotes.md
func (c *Client) GetQuotes(ctx context.Context, q DepthQuery) (*Quote, error) {
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
	if err := c.get(ctx, pathStockDepths, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
