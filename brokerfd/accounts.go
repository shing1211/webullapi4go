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
	pathFDAccountList       = "/broker-fd/account/list"
	pathFDAccountDetail     = "/broker-fd/account/detail"
	pathFDAccountCreate     = "/broker-fd/account/create"
	pathFDAccountUpdate     = "/broker-fd/account/update"
	pathFDAccountClose      = "/broker-fd/account/close"
	pathFDAccountForms      = "/broker-fd/account/forms"
	pathFDAccountFormDetail = "/broker-fd/account/form/detail"
	pathFDAccountFormSubmit = "/broker-fd/account/form/submit"
	pathFDAccountFormStatus = "/broker-fd/account/form/status"
)

type FDAccount struct {
	AccountID     string `json:"account_id"`
	AccountNumber string `json:"account_number"`
	AccountType   string `json:"account_type"`
	AccountClass  string `json:"account_class"`
	Status        string `json:"status"`
	Currency      string `json:"currency"`
	CreateTime    string `json:"create_time"`
}

func (c *Client) ListFDAccounts(ctx context.Context) ([]FDAccount, error) {
	var out []FDAccount
	if err := c.get(ctx, pathFDAccountList, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetFDAccountDetail(ctx context.Context, accountID string) (*FDAccount, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out FDAccount
	if err := c.get(ctx, pathFDAccountDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type CreateFDAccountRequest struct {
	AccountType  string `json:"account_type"`
	AccountClass string `json:"account_class"`
	Currency     string `json:"currency"`
}

func (c *Client) CreateFDAccount(ctx context.Context, req CreateFDAccountRequest) (*FDAccount, error) {
	var out FDAccount
	if err := c.post(ctx, pathFDAccountCreate, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type UpdateFDAccountRequest struct {
	AccountID    string `json:"account_id"`
	AccountType  string `json:"account_type,omitempty"`
	AccountClass string `json:"account_class,omitempty"`
}

func (c *Client) UpdateFDAccount(ctx context.Context, req UpdateFDAccountRequest) (*FDAccount, error) {
	var out FDAccount
	if err := c.put(ctx, pathFDAccountUpdate, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CloseFDAccount(ctx context.Context, accountID string) error {
	q := url.Values{}
	q.Set("account_id", accountID)
	return c.post(ctx, pathFDAccountClose, q, nil, nil)
}

type AccountForm struct {
	FormID     string `json:"form_id"`
	FormType   string `json:"form_type"`
	Status     string `json:"status"`
	CreateTime string `json:"create_time"`
	Content    string `json:"content,omitempty"`
}

func (c *Client) ListAccountForms(ctx context.Context, accountID string) ([]AccountForm, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out []AccountForm
	if err := c.get(ctx, pathFDAccountForms, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetAccountFormDetail(ctx context.Context, formID string) (*AccountForm, error) {
	q := url.Values{}
	q.Set("form_id", formID)
	var out AccountForm
	if err := c.get(ctx, pathFDAccountFormDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type SubmitAccountFormRequest struct {
	FormID  string `json:"form_id"`
	Content string `json:"content"`
}

func (c *Client) SubmitAccountForm(ctx context.Context, req SubmitAccountFormRequest) (*AccountForm, error) {
	var out AccountForm
	if err := c.post(ctx, pathFDAccountFormSubmit, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetAccountFormStatus(ctx context.Context, formID string) (*AccountForm, error) {
	q := url.Values{}
	q.Set("form_id", formID)
	var out AccountForm
	if err := c.get(ctx, pathFDAccountFormStatus, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
