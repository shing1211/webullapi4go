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
)

// Display Solution screener endpoints. These use the Display Solution host
// (co-branding-openapi.webull.hk) with C2S Bearer token auth.
//
// TODO(ds): Confirm exact paths against US sandbox.
const (
	// pathDSGainersLosers is the Display Solution gainers/losers endpoint.
	pathDSGainersLosers = "/openapi/market-data/screener/rank"
	// pathDSTopActive is the Display Solution top-active endpoint.
	pathDSTopActive = "/openapi/market-data/screener/top-active"
)

// GetDisplayGainersLosers retrieves gainers/losers via Display Solution.
//
// TODO(ds): path unconfirmed
func (c *Client) GetDisplayGainersLosers(ctx context.Context, q GainersLosersQuery) ([]ScreenerStock, error) {
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

	var out []ScreenerStock
	if err := c.DisplayService().Get(ctx, pathDSGainersLosers, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDisplayTopActive retrieves top active stocks via Display Solution.
//
// TODO(ds): path unconfirmed
func (c *Client) GetDisplayTopActive(ctx context.Context, q MostActiveQuery) ([]ScreenerStock, error) {
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

	var out []ScreenerStock
	if err := c.DisplayService().Get(ctx, pathDSTopActive, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}
