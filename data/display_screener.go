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

	"github.com/shing1211/webullapi4go/pkg/types"
)

// Display Solution screener endpoints. These use the Display Solution host
// (co-branding-openapi.webull.hk) with C2S Bearer token auth.
const (
	// pathDSGainersLosers is the Display Solution gainers/losers endpoint.
	pathDSGainersLosers = "/market-data/screeners/gainers-losers/list"
	// pathDSTopActive is the Display Solution top-active endpoint.
	pathDSTopActive = "/market-data/screeners/top-actives/list"
)

//
// **Breaking, in v2.1.35.** The page documents the 200 body as
// `{"data": [...], "pagination_key": "..."}`, which a bare slice cannot
// decode, so this method failed outright against a conforming server. It
// returns a [types.Page] now: a caller reads out.Data instead of out, and
// passes out.PaginationKey back to fetch the following page. The new form
// succeeds where the old one could not.

// GetDisplayGainersLosers retrieves gainers/losers via Display Solution.
func (c *Client) GetDisplayGainersLosers(ctx context.Context, q GainersLosersQuery) (*types.Page[ScreenerStock], error) {
	query := url.Values{}
	if q.RankType != "" {
		query.Set("rank_type", string(q.RankType))
	}
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	if q.SortBy != "" {
		query.Set("sort_by", string(q.SortBy))
	}
	if q.Direction != "" {
		query.Set("direction", string(q.Direction))
	}

	if q.PaginationKey != "" {
		query.Set("pagination_key", q.PaginationKey)
	}

	var out types.Page[ScreenerStock]
	if err := c.DisplayService().Get(ctx, pathDSGainersLosers, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

//
// **Breaking, in v2.1.35.** The page documents the 200 body as
// `{"data": [...], "pagination_key": "..."}`, which a bare slice cannot
// decode, so this method failed outright against a conforming server. It
// returns a [types.Page] now: a caller reads out.Data instead of out, and
// passes out.PaginationKey back to fetch the following page. The new form
// succeeds where the old one could not.

// GetDisplayTopActive retrieves top active stocks via Display Solution.
func (c *Client) GetDisplayTopActive(ctx context.Context, q MostActiveQuery) (*types.Page[ScreenerStock], error) {
	query := url.Values{}
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	if q.RankType != "" {
		query.Set("rank_type", string(q.RankType))
	}
	if q.SortBy != "" {
		query.Set("sort_by", string(q.SortBy))
	}
	if q.Direction != "" {
		query.Set("direction", string(q.Direction))
	}

	if q.PaginationKey != "" {
		query.Set("pagination_key", q.PaginationKey)
	}

	var out types.Page[ScreenerStock]
	if err := c.DisplayService().Get(ctx, pathDSTopActive, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
