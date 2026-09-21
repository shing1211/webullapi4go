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
	pathCashJournal           = "/openapi/v1/broker/journals/cash"
	pathCashJournalDetail     = "/openapi/v1/broker/journals/cash/detail"
	pathPositionJournal       = "/openapi/v1/broker/journals/position"
	pathPositionJournalDetail = "/openapi/v1/broker/journals/position/detail"
)

type CashJournal struct {
	JournalID  string `json:"journal_id"`
	AccountID  string `json:"account_id"`
	Type       string `json:"type"`
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
	Status     string `json:"status"`
	CreateTime string `json:"create_time"`
}

type CreateCashJournalRequest struct {
	AccountID string `json:"account_id"`
	Type      string `json:"type"`
	Amount    string `json:"amount"`
	Currency  string `json:"currency"`
}

func (c *Client) CreateCashJournal(ctx context.Context, req CreateCashJournalRequest) (*CashJournal, error) {
	var out CashJournal
	if err := c.post(ctx, pathCashJournal, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetCashJournalDetail(ctx context.Context, journalID string) (*CashJournal, error) {
	q := url.Values{}
	q.Set("journal_id", journalID)
	var out CashJournal
	if err := c.get(ctx, pathCashJournalDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type PositionJournal struct {
	JournalID  string `json:"journal_id"`
	AccountID  string `json:"account_id"`
	Symbol     string `json:"symbol"`
	Quantity   string `json:"quantity"`
	Action     string `json:"action"`
	Status     string `json:"status"`
	CreateTime string `json:"create_time"`
}

type CreatePositionJournalRequest struct {
	AccountID string `json:"account_id"`
	Symbol    string `json:"symbol"`
	Quantity  string `json:"quantity"`
	Action    string `json:"action"`
}

func (c *Client) CreatePositionJournal(ctx context.Context, req CreatePositionJournalRequest) (*PositionJournal, error) {
	var out PositionJournal
	if err := c.post(ctx, pathPositionJournal, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetPositionJournalDetail(ctx context.Context, journalID string) (*PositionJournal, error) {
	q := url.Values{}
	q.Set("journal_id", journalID)
	var out PositionJournal
	if err := c.get(ctx, pathPositionJournalDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
