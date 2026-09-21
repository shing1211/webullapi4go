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

// Virtual account management endpoint paths.
const (
	pathVirtualAccounts       = "/openapi/v1/broker/virtual-accounts"
	pathVirtualAccountsDetail = "/openapi/v1/broker/virtual-accounts/detail"
)

// VirtualAccount represents a virtual trading account.
type VirtualAccount struct {
	AccountID   string `json:"account_id"`
	AccountName string `json:"account_name"`
	AccountType string `json:"account_type"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
}

// CreateVirtualAccountRequest is the body for [Client.CreateVirtualAccount].
type CreateVirtualAccountRequest struct {
	AccountName string `json:"account_name"`
	AccountType string `json:"account_type"`
	Currency    string `json:"currency"`
}

// UpdateVirtualAccountRequest is the body for [Client.UpdateVirtualAccount].
type UpdateVirtualAccountRequest struct {
	AccountName string `json:"account_name,omitempty"`
}

// CreateVirtualAccount creates a new virtual account.
func (c *Client) CreateVirtualAccount(ctx context.Context, req CreateVirtualAccountRequest) (*VirtualAccount, error) {
	var out VirtualAccount
	if err := c.post(ctx, pathVirtualAccounts, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateVirtualAccount updates an existing virtual account identified by accountID.
func (c *Client) UpdateVirtualAccount(ctx context.Context, accountID string, req UpdateVirtualAccountRequest) (*VirtualAccount, error) {
	path := pathVirtualAccountsDetail + "?account_id=" + accountID
	var out VirtualAccount
	if err := c.put(ctx, path, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetVirtualAccount retrieves a single virtual account by accountID.
func (c *Client) GetVirtualAccount(ctx context.Context, accountID string) (*VirtualAccount, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out VirtualAccount
	if err := c.get(ctx, pathVirtualAccountsDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListVirtualAccounts lists all virtual accounts.
func (c *Client) ListVirtualAccounts(ctx context.Context) ([]VirtualAccount, error) {
	var out []VirtualAccount
	if err := c.get(ctx, pathVirtualAccounts, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
