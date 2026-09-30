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

package brokerfd

import (
	"context"
	"net/url"
	"strings"

	"github.com/shing1211/webullapi4go/pkg/domain/money"
	"github.com/shing1211/webullapi4go/pkg/types"
)

const (
	pathFDStockInstruments      = "/broker/instruments/stocks/profiles/list"
	pathFDStockLocate           = "/broker-fd/instruments/stock-locate"
	pathFDCorporateActions      = "/broker/instruments/stocks/corporate-actions/get"
	pathFDCorporateActionDetail = "/broker-fd/instruments/corporate-action/detail"
	pathFDECInstruments         = "/broker/instruments/event-contracts/markets/list"
	pathFDECInstrumentDetail    = "/broker-fd/instruments/event-contract/detail"
)

// FDStockInstrument represents a fractional share stock instrument available for trading.
type FDStockInstrument struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Currency string `json:"currency"`
	LotSize  string `json:"lot_size"`
	Status   string `json:"status"`
}

// GetFDStockInstruments retrieves fractional share stock instruments by symbol list.
func (c *Client) GetFDStockInstruments(ctx context.Context, symbols []string, paginationKey string) (*types.Page[FDStockInstrument], error) {
	q := url.Values{}
	q.Set("symbols", strings.Join(symbols, ","))
	if paginationKey != "" {
		q.Set("pagination_key", paginationKey)
	}
	var out types.Page[FDStockInstrument]
	if err := c.get(ctx, pathFDStockInstruments, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FDStockLocate represents locate (borrowed shares) information for a fractional share.
type FDStockLocate struct {
	Symbol         string      `json:"symbol"`
	LocateQuantity string      `json:"locate_quantity"`
	Available      string      `json:"available"`
	Rate           money.Money `json:"rate"`
}

// GetFDStockLocate retrieves locate information for a fractional share symbol.
//
// No Webull documentation page exists for this endpoint. The path uses the
// /broker-fd/ namespace which does not appear in any cached page. This method
// was added from inference and has not been verified against a live response.
//
// See IMPLEMENTATION_STATUS.md item 17 "no documented endpoint" class.
func (c *Client) GetFDStockLocate(ctx context.Context, symbol string) ([]FDStockLocate, error) {
	q := url.Values{}
	q.Set("symbol", symbol)
	var out []FDStockLocate
	if err := c.get(ctx, pathFDStockLocate, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FDCorporateAction represents a corporate action (dividend, split, etc.) for fractional shares.
type FDCorporateAction struct {
	ActionID   string `json:"action_id"`
	Symbol     string `json:"symbol"`
	ActionType string `json:"action_type"`
	ExDate     string `json:"ex_date"`
	RecordDate string `json:"record_date"`
	PayDate    string `json:"pay_date"`
	Ratio      string `json:"ratio"`

	// The fields below are declared by the page but not marked required, so they were
	// absent here and a response carrying them decoded the value to the zero value
	// with no error reported. That is weaker evidence than a missing-required-name
	// row: the page publishes no required list for this call, so a name may be
	// optional or conditionally sent, and nothing here asserts that a live server
	// sends it. Not live-verified: the Broker FD host returns 404 in the HK sandbox
	// and no US credential was available here.
	Category             string                    `json:"category"`
	CountryCode          string                    `json:"country_code"`
	IssuerCountryCode    string                    `json:"issuer_country_code"`
	ListingCountryOfCode string                    `json:"listing_country_of_code"`
	EventID              string                    `json:"event_id"`
	EventType            string                    `json:"event_type"`
	EventVersion         string                    `json:"event_version"`
	PaymentDate          string                    `json:"payment_date"`
	FinalPayDate         string                    `json:"final_pay_date"`
	InstrumentID         string                    `json:"instrument_id"`
	From                 FDCorporateActionPosition `json:"from"`
	To                   []FDCorporateActionTarget `json:"to"`
}

// FDCorporateActionPosition is the position's instrument on a corporate action, named
// "from" on the page.
type FDCorporateActionPosition struct {
	// Symbol is the trading symbol of the position's instrument.
	Symbol string `json:"symbol"`
	// Name is the display name of the position's instrument.
	Name string `json:"name"`
	// Exchange is the exchange the position's instrument trades on.
	Exchange string `json:"exchange"`
}

// FDCorporateActionTarget is one target instrument of a corporate action.
//
// The page documents the element with default_option_flag, description,
// option_number and payouts, and publishes no required list for any of them, so each
// is optional in practice. Payouts is a free-form object on the page: it declares no
// property, so nothing here can name its fields and it is carried as a map.
type FDCorporateActionTarget struct {
	// DefaultOptionFlag reports whether this target is the default option.
	DefaultOptionFlag string `json:"default_option_flag"`
	// Description is the human-readable description of the target.
	Description string `json:"description"`
	// OptionNumber is the option's number within the contract.
	OptionNumber string `json:"option_number"`
	// Payouts is the free-form payout detail for this target. The page declares the
	// property but no property inside it, so the value is carried untyped.
	Payouts map[string]any `json:"payouts"`
}

// GetFDCorporateActions retrieves corporate actions for a fractional share symbol.
func (c *Client) GetFDCorporateActions(ctx context.Context, symbol string) ([]FDCorporateAction, error) {
	q := url.Values{}
	q.Set("symbol", symbol)
	var out []FDCorporateAction
	if err := c.get(ctx, pathFDCorporateActions, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFDCorporateActionDetail retrieves details for a specific corporate action by its ID.
func (c *Client) GetFDCorporateActionDetail(ctx context.Context, actionID string) (*FDCorporateAction, error) {
	q := url.Values{}
	q.Set("action_id", actionID)
	var out FDCorporateAction
	if err := c.get(ctx, pathFDCorporateActionDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FDECInstrument represents an event contract (EC) instrument for fractional shares.
type FDECInstrument struct {
	Symbol         string      `json:"symbol"`
	EventID        string      `json:"event_id"`
	SeriesID       string      `json:"series_id"`
	StrikePrice    money.Money `json:"strike_price"`
	ExpirationDate string      `json:"expiration_date"`
	Status         string      `json:"status"`
}

// GetFDECInstruments retrieves event contract instruments by event ID.
func (c *Client) GetFDECInstruments(ctx context.Context, eventID, paginationKey string) (*types.Page[FDECInstrument], error) {
	q := url.Values{}
	q.Set("event_id", eventID)
	if paginationKey != "" {
		q.Set("pagination_key", paginationKey)
	}
	var out types.Page[FDECInstrument]
	if err := c.get(ctx, pathFDECInstruments, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetFDECInstrumentDetail retrieves details for a specific event contract by its symbol.
//
// No Webull documentation page exists for this endpoint. The path uses the
// /broker-fd/ namespace which does not appear in any cached page. The cached
// /broker/instruments/event-contracts/markets/list page covers GetFDECInstruments
// only; it does not name a /broker-fd/instruments/event-contract/detail endpoint.
//
// See IMPLEMENTATION_STATUS.md item 17 "no documented endpoint" class.
func (c *Client) GetFDECInstrumentDetail(ctx context.Context, symbol string) (*FDECInstrument, error) {
	q := url.Values{}
	q.Set("symbol", symbol)
	var out FDECInstrument
	if err := c.get(ctx, pathFDECInstrumentDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
