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

const (
	pathTradeCalendar = "/broker/master-data/trading-calendars/list"
)

type TradeCalendar struct {
	Date      string `json:"date"`
	Market    string `json:"market"`
	Status    string `json:"status"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
}

func (c *Client) GetTradeCalendar(ctx context.Context, market, startDate, endDate string) ([]TradeCalendar, error) {
	query := url.Values{}
	if market != "" {
		query.Set("market", market)
	}
	if startDate != "" {
		query.Set("start_date", startDate)
	}
	if endDate != "" {
		query.Set("end_date", endDate)
	}
	var out []TradeCalendar
	if err := c.get(ctx, pathTradeCalendar, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}
