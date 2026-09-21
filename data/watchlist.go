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
	"encoding/json"
	"net/http"
	"net/url"
)

// Watchlist endpoints.
//
// These are the current documented paths. Older Webull documentation exposed
// them under the /openapi/... alias; prefer the /market-data/... paths below.
//
// References:
//   - https://developer.webull.hk/apis/docs/reference/get-watchlist.md
//   - https://developer.webull.hk/apis/docs/reference/create-watchlist.md
//   - https://developer.webull.hk/apis/docs/reference/update-watchlist.md
//   - https://developer.webull.hk/apis/docs/reference/delete-watchlist.md
const (
	pathWatchlists                 = "/market-data/watchlists/list"
	pathWatchlistCreate            = "/market-data/watchlists/create"
	pathWatchlistUpdate            = "/market-data/watchlists/update"
	pathWatchlistDelete            = "/market-data/watchlists/delete"
	pathWatchlistInstruments       = "/market-data/watchlists/instruments/list"
	pathWatchlistInstrumentsAdd    = "/market-data/watchlists/instruments/add"
	pathWatchlistInstrumentsRemove = "/market-data/watchlists/instruments/remove"
	pathWatchlistInstrumentsUpdate = "/market-data/watchlists/instruments/update"
)

// Watchlist is the metadata of one user watchlist.
type Watchlist struct {
	// WatchlistID is the unique identifier of the watchlist, for example
	// "12345678".
	WatchlistID string `json:"watchlist_id"`
	// Name is the display name of the watchlist.
	Name string `json:"name"`
	// Sort is the display ordering number.
	Sort int32 `json:"sort"`
	// CreateTime is the watchlist creation time in ISO 8601 form.
	CreateTime string `json:"create_time"`
	// UpdateTime is the time of the last watchlist update in ISO 8601 form.
	UpdateTime string `json:"update_time"`
}

// WatchlistInstrument is one instrument held in a watchlist.
type WatchlistInstrument struct {
	// InstrumentID is the unique identifier of the instrument.
	InstrumentID string `json:"instrument_id"`
	// Symbol is the trading symbol, for example "AAPL" or "00700".
	Symbol string `json:"symbol"`
	// Name is the display name of the instrument.
	Name string `json:"name"`
	// ExchangeCode is the standardized exchange code, for example "NSQ".
	ExchangeCode string `json:"exchange_code"`
	// Sort is the sort order within the watchlist.
	Sort int32 `json:"sort"`
	// AddedTime is when the instrument was added, in ISO 8601 form.
	AddedTime string `json:"added_time"`
}

// WatchlistInstruments is a watchlist together with the instruments it holds.
type WatchlistInstruments struct {
	// WatchlistID is the unique identifier of the watchlist.
	WatchlistID string `json:"watchlist_id"`
	// Instruments lists the instruments in the watchlist.
	Instruments []WatchlistInstrument `json:"instruments"`
}

// CreateWatchlistParams parameterizes [Client.CreateWatchlist]. Name is
// required; a zero Sort lets the server choose the next display order.
type CreateWatchlistParams struct {
	// Name is the watchlist name. Required.
	Name string `json:"name"`
	// Sort is the display ordering number. Zero means unset.
	Sort int32 `json:"sort,omitempty"`
}

// CreateWatchlistResult is returned by [Client.CreateWatchlist].
type CreateWatchlistResult struct {
	// WatchlistID is the identifier assigned to the new watchlist.
	WatchlistID string `json:"watchlist_id"`
}

// UpdateWatchlistParams parameterizes [Client.UpdateWatchlist]. WatchlistID is
// required; a zero Sort or empty Name leaves the corresponding field unchanged.
type UpdateWatchlistParams struct {
	// WatchlistID identifies the watchlist to update. Required.
	WatchlistID string `json:"watchlist_id"`
	// Name is the new watchlist name. Empty means unchanged.
	Name string `json:"name,omitempty"`
	// Sort is the new display ordering number. Zero means unchanged.
	Sort int32 `json:"sort,omitempty"`
}

// WatchlistInstrumentParam identifies one instrument in a watchlist mutation
// request. Symbol and Category are required; Sort is only meaningful for
// [Client.UpdateWatchlistInstruments].
type WatchlistInstrumentParam struct {
	// Symbol is the security symbol, for example "AAPL".
	Symbol string `json:"symbol"`
	// Category is the security category, for example
	// [StockCategoryUS]. Required.
	Category StockCategory `json:"category"`
	// Sort is the display ordering number. Zero means unset.
	Sort int32 `json:"sort,omitempty"`
}

// WatchlistInstrumentsParam parameterizes the watchlist instrument mutation
// endpoints ([Client.AddWatchlistInstruments],
// [Client.RemoveWatchlistInstruments], and
// [Client.UpdateWatchlistInstruments]).
type WatchlistInstrumentsParam struct {
	// WatchlistID identifies the watchlist to mutate. Required.
	WatchlistID string `json:"watchlist_id"`
	// Instruments is the non-empty list of instruments to add, remove, or
	// reorder. Required.
	Instruments []WatchlistInstrumentParam `json:"instruments"`
}

// SuccessResponse reports the outcome of a mutating watchlist operation.
type SuccessResponse struct {
	// Success reports whether the server applied the operation.
	Success bool `json:"success"`
}

type BoolOrSuccess struct {
	Success bool
}

func (b *BoolOrSuccess) UnmarshalJSON(data []byte) error {
	if string(data) == "true" {
		b.Success = true
		return nil
	}
	if string(data) == "false" {
		b.Success = false
		return nil
	}
	var obj struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	b.Success = obj.Success
	return nil
}

// watchlistIDRequest is the JSON body shared by the body-only watchlist
// operations that take a single identifier.
type watchlistIDRequest struct {
	WatchlistID string `json:"watchlist_id"`
}

// GetWatchlists retrieves the authenticated user's watchlists.
//
// Reference: https://developer.webull.hk/apis/docs/reference/get-watchlist.md
func (c *Client) GetWatchlists(ctx context.Context) ([]Watchlist, error) {
	var out []Watchlist
	if err := c.get(ctx, pathWatchlists, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateWatchlist creates a new watchlist and returns its identifier.
//
// Reference: https://developer.webull.hk/apis/docs/reference/create-watchlist.md
func (c *Client) CreateWatchlist(ctx context.Context, params CreateWatchlistParams) (*CreateWatchlistResult, error) {
	var out CreateWatchlistResult
	if err := c.do(ctx, http.MethodPost, pathWatchlistCreate, nil, params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateWatchlist renames a watchlist and/or changes its sort order.
//
// Reference: https://developer.webull.hk/apis/docs/reference/update-watchlist.md
func (c *Client) UpdateWatchlist(ctx context.Context, params UpdateWatchlistParams) (*SuccessResponse, error) {
	var out BoolOrSuccess
	if err := c.do(ctx, http.MethodPost, pathWatchlistUpdate, nil, params, &out); err != nil {
		return nil, err
	}
	return &SuccessResponse{Success: out.Success}, nil
}

// DeleteWatchlist deletes the watchlist identified by watchlistID.
//
// Reference: https://developer.webull.hk/apis/docs/reference/delete-watchlist.md
func (c *Client) DeleteWatchlist(ctx context.Context, watchlistID string) (*SuccessResponse, error) {
	var out BoolOrSuccess
	body := watchlistIDRequest{WatchlistID: watchlistID}
	if err := c.do(ctx, http.MethodPost, pathWatchlistDelete, nil, body, &out); err != nil {
		return nil, err
	}
	return &SuccessResponse{Success: out.Success}, nil
}

// GetWatchlistInstruments retrieves the instruments held in watchlistID.
//
// Reference: https://developer.webull.hk/apis/docs/reference/get-watchlist-instruments.md
func (c *Client) GetWatchlistInstruments(ctx context.Context, watchlistID string) (*WatchlistInstruments, error) {
	query := url.Values{"watchlist_id": {watchlistID}}
	var out WatchlistInstruments
	if err := c.get(ctx, pathWatchlistInstruments, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddWatchlistInstruments adds instruments to a watchlist.
//
// Reference: https://developer.webull.hk/apis/docs/reference/add-watchlist-instruments.md
func (c *Client) AddWatchlistInstruments(ctx context.Context, params WatchlistInstrumentsParam) (*SuccessResponse, error) {
	var out BoolOrSuccess
	if err := c.do(ctx, http.MethodPost, pathWatchlistInstrumentsAdd, nil, params, &out); err != nil {
		return nil, err
	}
	return &SuccessResponse{Success: out.Success}, nil
}

// RemoveWatchlistInstruments removes instruments from a watchlist.
//
// Reference: https://developer.webull.hk/apis/docs/reference/remove-watchlist-instruments.md
func (c *Client) RemoveWatchlistInstruments(ctx context.Context, params WatchlistInstrumentsParam) (*SuccessResponse, error) {
	var out BoolOrSuccess
	if err := c.do(ctx, http.MethodPost, pathWatchlistInstrumentsRemove, nil, params, &out); err != nil {
		return nil, err
	}
	return &SuccessResponse{Success: out.Success}, nil
}

// UpdateWatchlistInstruments changes the sort order of instruments in a
// watchlist.
//
// Reference: https://developer.webull.hk/apis/docs/reference/update-watchlist-instruments.md
func (c *Client) UpdateWatchlistInstruments(ctx context.Context, params WatchlistInstrumentsParam) (*SuccessResponse, error) {
	var out BoolOrSuccess
	if err := c.do(ctx, http.MethodPost, pathWatchlistInstrumentsUpdate, nil, params, &out); err != nil {
		return nil, err
	}
	return &SuccessResponse{Success: out.Success}, nil
}
