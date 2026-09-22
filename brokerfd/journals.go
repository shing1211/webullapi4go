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
	pathFDCashJournal       = "/broker-fd/journals/cash"
	pathFDCashJournalDetail = "/broker-fd/journals/cash/detail"
)

// FDCashJournal represents a cash journal entry recording a deposit, withdrawal, or other
// cash movement for a broker-fd account.
type FDCashJournal struct {
	JournalID  string `json:"journal_id"`
	AccountID  string `json:"account_id"`
	Type       string `json:"type"`
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
	Status     string `json:"status"`
	CreateTime string `json:"create_time"`
}

// ListFDCashJournals returns all cash journal entries for the specified broker-fd account.
// Results are ordered by creation time descending.
func (c *Client) ListFDCashJournals(ctx context.Context, accountID string) ([]FDCashJournal, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out []FDCashJournal
	if err := c.get(ctx, pathFDCashJournal, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFDCashJournalDetail returns a single cash journal entry by its journal ID.
func (c *Client) GetFDCashJournalDetail(ctx context.Context, journalID string) (*FDCashJournal, error) {
	q := url.Values{}
	q.Set("journal_id", journalID)
	var out FDCashJournal
	if err := c.get(ctx, pathFDCashJournalDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
