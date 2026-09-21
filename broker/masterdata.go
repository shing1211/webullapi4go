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

package broker

import (
	"context"
	"net/url"
)

// Master data endpoints.
const (
	pathTradeCalendar = "/openapi/v1/broker/master-data/trade-calendar"
)

// TradeCalendar represents a trading calendar entry.
type TradeCalendar struct {
	Date      string `json:"date"`
	Market    string `json:"market"`
	Status    string `json:"status"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
}

// GetTradeCalendar retrieves the trade calendar for a date range.
func (c *Client) GetTradeCalendar(ctx context.Context, market, startDate, endDate string) ([]TradeCalendar, error) {
	q := url.Values{}
	q.Set("market", market)
	q.Set("start_date", startDate)
	q.Set("end_date", endDate)
	var out []TradeCalendar
	if err := c.get(ctx, pathTradeCalendar, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}
