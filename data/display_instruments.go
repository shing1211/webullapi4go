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

// Display Solution instrument endpoints.
//
// TODO(ds): Confirm exact paths against US sandbox.
const (
	pathDSCompanyProfile = "/openapi/market-data/stock/company-profile"
	pathDSAnalystTarget  = "/openapi/market-data/stock/analyst-target-price"
	pathDSAnalystRating  = "/openapi/market-data/stock/analyst-rating"
)

// GetDSCompanyProfile retrieves the company profile for symbol via the Display
// Solution endpoint.
//
// TODO(ds): Confirm exact paths against US sandbox.
func (c *Client) GetDSCompanyProfile(ctx context.Context, symbol string) (*CompanyProfile, error) {
	query := url.Values{"symbol": {symbol}}
	var out CompanyProfile
	if err := c.DisplayService().Get(ctx, pathDSCompanyProfile, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetDSAnalystTargetPrice retrieves aggregate analyst target-price data for
// symbol via the Display Solution endpoint.
//
// TODO(ds): Confirm exact paths against US sandbox.
func (c *Client) GetDSAnalystTargetPrice(ctx context.Context, symbol string) (*AnalystTargetPrice, error) {
	query := url.Values{"symbol": {symbol}}
	var out AnalystTargetPrice
	if err := c.DisplayService().Get(ctx, pathDSAnalystTarget, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetDSAnalystRating retrieves aggregate analyst rating counts for symbol via
// the Display Solution endpoint.
//
// TODO(ds): Confirm exact paths against US sandbox.
func (c *Client) GetDSAnalystRating(ctx context.Context, symbol string) (*AnalystRating, error) {
	query := url.Values{"symbol": {symbol}}
	var out AnalystRating
	if err := c.DisplayService().Get(ctx, pathDSAnalystRating, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
