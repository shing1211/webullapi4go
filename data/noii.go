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

	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

// NOII (Net Order Imbalance Indicator) endpoints.
const (
	// pathNOIIBars is the NOII bars endpoint.
	//
	// Reference: https://developer.webull.hk/apis/docs/reference/get-noii-bars.md
	pathNOIIBars = "/market-data/stocks/noii-bars/list"
	// pathNOIISnapshot is the NOII snapshot endpoint.
	//
	// Reference: https://developer.webull.hk/apis/docs/reference/get-noii-snapshot.md
	pathNOIISnapshot = "/market-data/stocks/noii-snapshots/list"
)

// NOIIActionType selects which auction imbalance a NOII request describes.
type NOIIActionType string

// NOII auction types.
const (
	// NOIIActionPreOpen is the opening auction imbalance.
	NOIIActionPreOpen NOIIActionType = "PRE_OPEN"
	// NOIIActionPreClose is the closing auction imbalance.
	NOIIActionPreClose NOIIActionType = "PRE_CLOSE"
)

// NOIIQuery parameterizes [Client.GetNOIIBars] and [Client.GetNOIISnapshot].
// Every field is required.
type NOIIQuery struct {
	// Symbol is the security symbol. Only a single symbol is supported.
	Symbol string
	// Category is the security type. Only [StockCategoryUS] is supported.
	Category StockCategory
	// ImbalanceActionType selects the opening or closing auction.
	ImbalanceActionType NOIIActionType
}

// NOIIBar is a single Net Order Imbalance Indicator bar.
type NOIIBar struct {
	// InstrumentID is the unique identifier of the security.
	InstrumentID string `json:"instrument_id"`
	// Symbol is the security symbol, for example "AAPL".
	Symbol string `json:"symbol"`
	// ImbalanceTime is the data publish time, in milliseconds since the Unix
	// epoch.
	ImbalanceTime int64 `json:"imbalance_time"`
	// ImbalanceRefPrice is the reference price, as a decimal string.
	ImbalanceRefPrice money.Money `json:"imbalance_ref_price"`
	// ImbalanceNearPrice is the indicative match price (the most likely
	// execution price), as a decimal string.
	ImbalanceNearPrice money.Money `json:"imbalance_near_price"`
	// ImbalanceFarPrice is the far price (the price at which orders could
	// execute in extreme scenarios), as a decimal string.
	ImbalanceFarPrice money.Money `json:"imbalance_far_price"`
	// ImbalanceActionType is the auction the bar belongs to.
	ImbalanceActionType NOIIActionType `json:"imbalance_action_type"`
}

// NOIISnapshot is the latest Net Order Imbalance Indicator snapshot.
type NOIISnapshot struct {
	// InstrumentID is the unique identifier of the security.
	InstrumentID string `json:"instrument_id"`
	// Symbol is the security symbol, for example "AAPL".
	Symbol string `json:"symbol"`
	// PairedShares is the number of shares that can be matched under current
	// conditions, as a decimal string.
	PairedShares string `json:"paired_shares"`
	// ImbalanceShares is the number of unmatched buy or sell shares, as a
	// decimal string.
	ImbalanceShares string `json:"imbalance_shares"`
	// ImbalanceSide is the direction of the imbalance.
	ImbalanceSide string `json:"imbalance_side"`
	// ImbalanceRefPrice is the reference price, as a decimal string.
	ImbalanceRefPrice money.Money `json:"imbalance_ref_price"`
	// ImbalanceNearPrice is the indicative match price (the most likely
	// execution price), as a decimal string.
	ImbalanceNearPrice money.Money `json:"imbalance_near_price"`
	// ImbalanceFarPrice is the far price (the price at which orders could
	// execute in extreme scenarios), as a decimal string.
	ImbalanceFarPrice money.Money `json:"imbalance_far_price"`
	// ImbalanceActionType is the auction the snapshot describes.
	ImbalanceActionType NOIIActionType `json:"imbalance_action_type"`
	// ImbalanceTime is the snapshot timestamp, in milliseconds since the Unix
	// epoch.
	ImbalanceTime int64 `json:"imbalance_time"`
	// ImbalanceVarIndicator is the volatility/imbalance status indicator.
	ImbalanceVarIndicator string `json:"imbalance_var_indicator"`
}

// GetNOIIBars retrieves Net Order Imbalance Indicator (NOII) bars for the
// opening or closing auction of a US stock. Only a single symbol is supported.
//
// Reference: https://developer.webull.hk/apis/docs/reference/get-noii-bars.md
func (c *Client) GetNOIIBars(ctx context.Context, q NOIIQuery) ([]NOIIBar, error) {
	var out []NOIIBar
	if err := c.get(ctx, pathNOIIBars, q.values(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetNOIISnapshot retrieves the latest Net Order Imbalance Indicator (NOII)
// snapshot for the opening or closing auction of a US stock. Only a single
// symbol is supported.
//
// Reference: https://developer.webull.hk/apis/docs/reference/get-noii-snapshot.md
func (c *Client) GetNOIISnapshot(ctx context.Context, q NOIIQuery) (*NOIISnapshot, error) {
	var out NOIISnapshot
	if err := c.get(ctx, pathNOIISnapshot, q.values(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// values encodes q as the query parameters shared by the NOII endpoints.
func (q NOIIQuery) values() url.Values {
	query := url.Values{}
	if q.Symbol != "" {
		query.Set("symbol", q.Symbol)
	}
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	if q.ImbalanceActionType != "" {
		query.Set("imbalance_action_type", string(q.ImbalanceActionType))
	}
	return query
}
