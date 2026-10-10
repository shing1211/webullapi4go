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

// UpdateVirtualAccountRequest is the documented request body of
// `POST /broker/accounts/virtual-accounts/update`.
//
// The page requires account_id and client_request_id. account_id is not a field
// here: [Client.UpdateVirtualAccount] takes it as its accountID argument and puts
// it in the body, so a caller cannot construct a request that omits it. Every
// field is optional in the encoding sense, so an unset one is omitted rather than
// sent as an empty string.
//
// This replaces a request type carrying a single AccountName field. The
// documented request has no account_name property at all and requires two fields
// the old type could not express, so the old shape was wrong in every respect:
// wrong verb, account_id in the query instead of the body, both required fields
// absent, and the one field it did send undocumented.
type UpdateVirtualAccountRequest struct {
	// ClientRequestID is required by the page: a caller-supplied id, unique per
	// request, at most 32 characters, limited to letters, digits, hyphen and
	// underscore.
	ClientRequestID string `json:"client_request_id"`
	// TradingPermissions is a list of permission codes, for example
	// US_STOCK_NORMAL, HK_STOCK_NORMAL, US_OPTION_NORMAL or CN_STOCK_NORMAL.
	TradingPermissions []string `json:"trading_permissions,omitempty"`
	// OptionLevel must be less than or equal to the account's permitted level.
	OptionLevel string `json:"option_level,omitempty"`
	// CommissionCode is for example STANDARD_COMMISSION.
	CommissionCode string `json:"commission_code,omitempty"`
	// W8BENInfo carries the W-8BEN tax certification. It is an object in the
	// request, though the matching VirtualAccount response field is a string.
	W8BENInfo *W8BENInfo `json:"w8ben_info,omitempty"`
	// ChinaConnectInvestorInfo carries the China Connect investor identity, and
	// like W8BENInfo is an object in the request.
	ChinaConnectInvestorInfo *ChinaConnectInvestorInfo `json:"china_connect_investor_info,omitempty"`
}

// PostalAddress is a postal address in a W-8BEN certification. The page documents
// it for both home_address and mail_address.
type PostalAddress struct {
	Country       string `json:"country,omitempty"`
	State         string `json:"state,omitempty"`
	City          string `json:"city,omitempty"`
	StreetAddress string `json:"street_address,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
}

// W8BENInfo is the W-8BEN tax certification in an update request.
type W8BENInfo struct {
	TreatyCountry string         `json:"treaty_country,omitempty"`
	TaxID         string         `json:"tax_id,omitempty"`
	SignDate      string         `json:"sign_date,omitempty"`
	FirstName     string         `json:"first_name,omitempty"`
	LastName      string         `json:"last_name,omitempty"`
	MiddleName    string         `json:"middle_name,omitempty"`
	HomeAddress   *PostalAddress `json:"home_address,omitempty"`
	MailAddress   *PostalAddress `json:"mail_address,omitempty"`
}

// ChinaConnectInvestorInfo is the investor identity in an update request.
type ChinaConnectInvestorInfo struct {
	FirstName         string `json:"first_name,omitempty"`
	LastName          string `json:"last_name,omitempty"`
	MiddleName        string `json:"middle_name,omitempty"`
	IDType            string `json:"id_type,omitempty"`
	IDNumber          string `json:"id_number,omitempty"`
	CountryOfIssuance string `json:"country_of_issuance,omitempty"`
}

// updateVirtualAccountBody is the wire shape: the documented fields plus the
// account_id argument, so the required pair reaches the body together.
type updateVirtualAccountBody struct {
	AccountID                string                    `json:"account_id"`
	ClientRequestID          string                    `json:"client_request_id"`
	TradingPermissions       []string                  `json:"trading_permissions,omitempty"`
	OptionLevel              string                    `json:"option_level,omitempty"`
	CommissionCode           string                    `json:"commission_code,omitempty"`
	W8BENInfo                *W8BENInfo                `json:"w8ben_info,omitempty"`
	ChinaConnectInvestorInfo *ChinaConnectInvestorInfo `json:"china_connect_investor_info,omitempty"`
}

func (c *Client) CreateVirtualAccount(ctx context.Context, req CreateVirtualAccountRequest) (*VirtualAccount, error) {
	var out VirtualAccount
	if err := c.post(ctx, pathVirtualAccountsCreate, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateVirtualAccount updates a Broker API HK virtual account.
//
// It issues POST, the documented verb for this path; the previous implementation
// issued PUT, which no documented source describes, and the reconciler reported
// the row as a clean path match because it compares paths and not verbs. accountID
// is sent in the JSON body, as the page requires, rather than as a query
// parameter.
func (c *Client) UpdateVirtualAccount(ctx context.Context, accountID string, req UpdateVirtualAccountRequest) (*VirtualAccount, error) {
	body := updateVirtualAccountBody{
		AccountID:                accountID,
		ClientRequestID:          req.ClientRequestID,
		TradingPermissions:       req.TradingPermissions,
		OptionLevel:              req.OptionLevel,
		CommissionCode:           req.CommissionCode,
		W8BENInfo:                req.W8BENInfo,
		ChinaConnectInvestorInfo: req.ChinaConnectInvestorInfo,
	}
	var out VirtualAccount
	if err := c.post(ctx, pathVirtualAccountsUpdate, nil, body, &out); err != nil {
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
