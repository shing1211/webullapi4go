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
	"strings"
)

// pathLogosBatch is the batch logos endpoint.
//
// Reference: https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/batch-logo-using-post
const pathLogosBatch = "/market-data/instruments/logos/batch"

// LogoQuery parameterizes [Client.GetLogos].
type LogoQuery struct {
	Symbols []string
}

// Logo URLs returned by the batch logo endpoint. URLs are hosted on Webull's CDN.
type Logo struct {
	Symbol string `json:"symbol"`
	Logo   string `json:"logo"`
}

// GetLogos retrieves logo image URLs for the specified securities.
//
// Reference: https://developer.webull.hk/apis/docs/reference/market-display-solution-data-api/batch-logo-using-post
func (c *Client) GetLogos(ctx context.Context, q LogoQuery) ([]Logo, error) {
	query := url.Values{}
	if len(q.Symbols) > 0 {
		query.Set("symbols", strings.Join(q.Symbols, ","))
	}

	var out []Logo
	if err := c.DisplayService().Do(ctx, http.MethodPost, pathLogosBatch, query, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
