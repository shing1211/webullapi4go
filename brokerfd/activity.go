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

	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

const (
	pathFDActivities = "/broker/activities/cash-activities/list"
)

// FDActivity represents a single activity or transaction in a fractional shares account.
// Amount is a string to preserve numeric precision; type and status are free-form strings.
type FDActivity struct {
	ActivityID  string      `json:"activity_id"`
	AccountID   string      `json:"account_id"`
	Type        string      `json:"type"`
	Amount      money.Money `json:"amount"`
	Currency    string      `json:"currency"`
	Status      string      `json:"status"`
	CreateTime  string      `json:"create_time"`
	Description string      `json:"description"`
}

// GetFDActivities retrieves all account activities for a fractional shares account,
// including trades, deposits, withdrawals, and other transactions.
func (c *Client) GetFDActivities(ctx context.Context, accountID string) ([]FDActivity, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out []FDActivity
	if err := c.get(ctx, pathFDActivities, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}
