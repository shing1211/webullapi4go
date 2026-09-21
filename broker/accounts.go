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
	pathVirtualAccountsList   = "/broker/accounts/virtual-accounts/list"
	pathVirtualAccountsGet    = "/broker/accounts/virtual-accounts/get"
	pathVirtualAccountsCreate = "/broker/accounts/virtual-accounts/create"
	pathVirtualAccountsUpdate = "/broker/accounts/virtual-accounts/update"
)

type VirtualAccount struct {
	AccountID                string   `json:"account_id"`
	AccountNumber            string   `json:"account_number"`
	AccountStatus            string   `json:"account_status"`
	AccountType              string   `json:"account_type"`
	BelongAccountID          string   `json:"belong_account_id"`
	BelongAccountNumber      string   `json:"belong_account_number"`
	TradingPermissions       []string `json:"trading_permissions"`
	OptionLevel              string   `json:"option_level"`
	CommissionCode           string   `json:"commission_code"`
	W8benInfo                string   `json:"w8ben_info"`
	ChinaConnectInvestorInfo string   `json:"china_connect_investor_info"`
}

type virtualAccountsResponse struct {
	Data []VirtualAccount `json:"data"`
}

type CreateVirtualAccountRequest struct {
	AccountName string `json:"account_name"`
	AccountType string `json:"account_type"`
	Currency    string `json:"currency"`
}

type UpdateVirtualAccountRequest struct {
	AccountName string `json:"account_name,omitempty"`
}

func (c *Client) CreateVirtualAccount(ctx context.Context, req CreateVirtualAccountRequest) (*VirtualAccount, error) {
	var out VirtualAccount
	if err := c.post(ctx, pathVirtualAccountsCreate, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateVirtualAccount(ctx context.Context, accountID string, req UpdateVirtualAccountRequest) (*VirtualAccount, error) {
	path := pathVirtualAccountsUpdate + "?account_id=" + accountID
	var out VirtualAccount
	if err := c.put(ctx, path, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetVirtualAccount(ctx context.Context, accountID string) (*VirtualAccount, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out VirtualAccount
	if err := c.get(ctx, pathVirtualAccountsGet, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListVirtualAccounts(ctx context.Context) ([]VirtualAccount, error) {
	var out virtualAccountsResponse
	if err := c.get(ctx, pathVirtualAccountsList, nil, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
