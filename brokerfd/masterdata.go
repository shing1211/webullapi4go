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
)

const (
	pathFDEnums         = "/broker-fd/master-data/enums"
	pathFDTradeCalendar = "/broker-fd/master-data/trade-calendar"
)

// FDEnum represents a key-value pair with a human-readable label, returned by the
// broker-fd master-data enums endpoint.
type FDEnum struct {
	EnumType string `json:"enum_type"`
	Value    string `json:"value"`
	Label    string `json:"label"`
}

// GetFDEnums returns the list of broker-fd master-data enums.
// The set of enum types and values is fixed for the lifetime of the API.
func (c *Client) GetFDEnums(ctx context.Context) ([]FDEnum, error) {
	var out []FDEnum
	if err := c.get(ctx, pathFDEnums, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FDTradeCalendarEntry represents a single trading day for a market.
type FDTradeCalendarEntry struct {
	Date      string `json:"date"`
	Market    string `json:"market"`
	Status    string `json:"status"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
}

// GetFDTradeCalendar returns the trading calendar for a given market between startDate and endDate.
// Dates are in YYYY-MM-DD format. The market parameter identifies the exchange or region.
func (c *Client) GetFDTradeCalendar(ctx context.Context, market, startDate, endDate string) ([]FDTradeCalendarEntry, error) {
	q := url.Values{}
	q.Set("market", market)
	q.Set("start_date", startDate)
	q.Set("end_date", endDate)
	var out []FDTradeCalendarEntry
	if err := c.get(ctx, pathFDTradeCalendar, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}
