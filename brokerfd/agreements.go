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
	pathAgreements      = "/broker-fd/agreements"
	pathAgreementDetail = "/broker-fd/agreements/detail"
)

type Agreement struct {
	AgreementID   string `json:"agreement_id"`
	AgreementType string `json:"agreement_type"`
	Version       string `json:"version"`
	Status        string `json:"status"`
	EffectiveDate string `json:"effective_date"`
	Content       string `json:"content,omitempty"`
}

func (c *Client) ListAgreements(ctx context.Context) ([]Agreement, error) {
	var out []Agreement
	if err := c.get(ctx, pathAgreements, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetAgreementDetail(ctx context.Context, agreementID string) (*Agreement, error) {
	q := url.Values{}
	q.Set("agreement_id", agreementID)
	var out Agreement
	if err := c.get(ctx, pathAgreementDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
