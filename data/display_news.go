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
	"net/http"
	"net/url"
)

// Display Solution news endpoints.
//
// TODO(ds): Confirm exact paths against US sandbox.
const (
	pathDSNewsSummary = "/openapi/market-data/news/watchlist-summary"
	pathDSMarketNews  = "/openapi/market-data/news/market"
	pathDSSymbolNews  = "/openapi/market-data/news/ticker"
	pathDSLatestNews  = "/openapi/market-data/news/latest"
)

// DSNewsSummaryItem is a single news item returned by the Display Solution
// news endpoints.
type DSNewsSummaryItem struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Source   string `json:"source"`
	URL      string `json:"url"`
	PubTime  string `json:"pub_time"`
	Symbol   string `json:"symbol"`
	Category string `json:"category"`
}

// GetDSNewsSummary retrieves news summaries for the given symbols via the
// Display Solution endpoint.
//
// TODO(ds): Confirm exact paths against US sandbox.
func (c *Client) GetDSNewsSummary(ctx context.Context, symbols []string) ([]DSNewsSummaryItem, error) {
	var out []DSNewsSummaryItem
	if err := c.DisplayService().Do(ctx, http.MethodPost, pathDSNewsSummary, nil, symbols, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDSMarketNews retrieves market-wide news for the given category via the
// Display Solution endpoint.
//
// TODO(ds): Confirm exact paths against US sandbox.
func (c *Client) GetDSMarketNews(ctx context.Context, category string) ([]DSNewsSummaryItem, error) {
	query := url.Values{"category": {category}}
	var out []DSNewsSummaryItem
	if err := c.DisplayService().Get(ctx, pathDSMarketNews, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDSSymbolNews retrieves news for a specific symbol via the Display
// Solution endpoint.
//
// TODO(ds): Confirm exact paths against US sandbox.
func (c *Client) GetDSSymbolNews(ctx context.Context, symbol string) ([]DSNewsSummaryItem, error) {
	query := url.Values{"symbol": {symbol}}
	var out []DSNewsSummaryItem
	if err := c.DisplayService().Get(ctx, pathDSSymbolNews, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetDSLatestNews retrieves the latest news headlines via the Display
// Solution endpoint.
//
// TODO(ds): Confirm exact paths against US sandbox.
func (c *Client) GetDSLatestNews(ctx context.Context) ([]DSNewsSummaryItem, error) {
	var out []DSNewsSummaryItem
	if err := c.DisplayService().Get(ctx, pathDSLatestNews, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
