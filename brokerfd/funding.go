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
	pathFDBankAccounts         = "/broker-fd/funding/bank-accounts"
	pathFDBankAccountDetail    = "/broker-fd/funding/bank-account/detail"
	pathFDBankAccountAdd       = "/broker-fd/funding/bank-account/add"
	pathFDBankAccountRemove    = "/broker-fd/funding/bank-account/remove"
	pathFDAchAccounts          = "/broker-fd/funding/ach-accounts"
	pathFDAchAccountDetail     = "/broker-fd/funding/ach-account/detail"
	pathFDAchAccountAdd        = "/broker-fd/funding/ach-account/add"
	pathFDAchAccountRemove     = "/broker-fd/funding/ach-account/remove"
	pathFDTransfers            = "/broker-fd/funding/transfers"
	pathFDTransferDetail       = "/broker-fd/funding/transfer/detail"
	pathFDTransferInitiate     = "/broker-fd/funding/transfer/initiate"
	pathFDInstantFunding       = "/broker-fd/funding/instant"
	pathFDInstantFundingDetail = "/broker-fd/funding/instant/detail"
	pathFDTransferFees         = "/broker-fd/funding/transfer/fees"
	pathFDCreditInfo           = "/broker-fd/funding/credit"
)

// BankAccount represents a linked bank account for funding operations.
type BankAccount struct {
	BankID        string `json:"bank_id"`
	AccountID     string `json:"account_id"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	RoutingNumber string `json:"routing_number"`
	Status        string `json:"status"`
}

// ListFDBankAccounts retrieves all linked bank accounts for a broker FD account.
func (c *Client) ListFDBankAccounts(ctx context.Context, accountID string) ([]BankAccount, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out []BankAccount
	if err := c.get(ctx, pathFDBankAccounts, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFDBankAccountDetail retrieves details for a specific bank account by its ID.
func (c *Client) GetFDBankAccountDetail(ctx context.Context, bankID string) (*BankAccount, error) {
	q := url.Values{}
	q.Set("bank_id", bankID)
	var out BankAccount
	if err := c.get(ctx, pathFDBankAccountDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddFDBankAccountRequest contains the parameters to link a new bank account.
type AddFDBankAccountRequest struct {
	AccountID     string `json:"account_id"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	RoutingNumber string `json:"routing_number"`
}

// AddFDBankAccount links a new bank account to the broker FD account.
func (c *Client) AddFDBankAccount(ctx context.Context, req AddFDBankAccountRequest) (*BankAccount, error) {
	var out BankAccount
	if err := c.post(ctx, pathFDBankAccountAdd, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RemoveFDBankAccount unlinks a bank account from the broker FD account.
func (c *Client) RemoveFDBankAccount(ctx context.Context, bankID string) error {
	q := url.Values{}
	q.Set("bank_id", bankID)
	return c.post(ctx, pathFDBankAccountRemove, q, nil, nil)
}

// ACHAccount represents an ACH account linked for electronic fund transfers.
type ACHAccount struct {
	ACHID         string `json:"ach_id"`
	AccountID     string `json:"account_id"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	Status        string `json:"status"`
}

// ListFDAchAccounts retrieves all linked ACH accounts for a broker FD account.
func (c *Client) ListFDAchAccounts(ctx context.Context, accountID string) ([]ACHAccount, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out []ACHAccount
	if err := c.get(ctx, pathFDAchAccounts, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFDAchAccountDetail retrieves details for a specific ACH account by its ID.
func (c *Client) GetFDAchAccountDetail(ctx context.Context, achID string) (*ACHAccount, error) {
	q := url.Values{}
	q.Set("ach_id", achID)
	var out ACHAccount
	if err := c.get(ctx, pathFDAchAccountDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddFDAchAccountRequest contains the parameters to link a new ACH account.
type AddFDAchAccountRequest struct {
	AccountID     string `json:"account_id"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
}

// AddFDAchAccount links a new ACH account to the broker FD account.
func (c *Client) AddFDAchAccount(ctx context.Context, req AddFDAchAccountRequest) (*ACHAccount, error) {
	var out ACHAccount
	if err := c.post(ctx, pathFDAchAccountAdd, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RemoveFDAchAccount unlinks an ACH account from the broker FD account.
func (c *Client) RemoveFDAchAccount(ctx context.Context, achID string) error {
	q := url.Values{}
	q.Set("ach_id", achID)
	return c.post(ctx, pathFDAchAccountRemove, q, nil, nil)
}

// Transfer represents a fund transfer transaction.
type Transfer struct {
	TransferID string `json:"transfer_id"`
	AccountID  string `json:"account_id"`
	Type       string `json:"type"`
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
	Status     string `json:"status"`
	CreateTime string `json:"create_time"`
}

// ListFDTransfers retrieves all fund transfers for a broker FD account.
func (c *Client) ListFDTransfers(ctx context.Context, accountID string) ([]Transfer, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out []Transfer
	if err := c.get(ctx, pathFDTransfers, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFDTransferDetail retrieves details for a specific transfer by its ID.
func (c *Client) GetFDTransferDetail(ctx context.Context, transferID string) (*Transfer, error) {
	q := url.Values{}
	q.Set("transfer_id", transferID)
	var out Transfer
	if err := c.get(ctx, pathFDTransferDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// InitiateTransferRequest contains the parameters to initiate a fund transfer.
type InitiateTransferRequest struct {
	AccountID string `json:"account_id"`
	Type      string `json:"type"`
	Amount    string `json:"amount"`
	Currency  string `json:"currency"`
}

// InitiateFDTransfer initiates a new fund transfer for the broker FD account.
func (c *Client) InitiateFDTransfer(ctx context.Context, req InitiateTransferRequest) (*Transfer, error) {
	var out Transfer
	if err := c.post(ctx, pathFDTransferInitiate, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// InstantFunding represents an instant funding transaction.
type InstantFunding struct {
	FundingID  string `json:"funding_id"`
	AccountID  string `json:"account_id"`
	Amount     string `json:"amount"`
	Status     string `json:"status"`
	CreateTime string `json:"create_time"`
}

// CreateFDInstantFunding creates an instant funding transaction for immediate funds.
func (c *Client) CreateFDInstantFunding(ctx context.Context, accountID, amount string) (*InstantFunding, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	q.Set("amount", amount)
	var out InstantFunding
	if err := c.post(ctx, pathFDInstantFunding, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetFDInstantFundingDetail retrieves details for a specific instant funding by its ID.
func (c *Client) GetFDInstantFundingDetail(ctx context.Context, fundingID string) (*InstantFunding, error) {
	q := url.Values{}
	q.Set("funding_id", fundingID)
	var out InstantFunding
	if err := c.get(ctx, pathFDInstantFundingDetail, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// TransferFee represents the fee associated with a transfer type.
type TransferFee struct {
	Type     string `json:"type"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

// GetFDTransferFees retrieves all available transfer fees.
func (c *Client) GetFDTransferFees(ctx context.Context) ([]TransferFee, error) {
	var out []TransferFee
	if err := c.get(ctx, pathFDTransferFees, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreditInfo represents credit/margin information for a broker FD account.
type CreditInfo struct {
	AccountID       string `json:"account_id"`
	CreditLimit     string `json:"credit_limit"`
	UsedCredit      string `json:"used_credit"`
	AvailableCredit string `json:"available_credit"`
}

// GetFDCreditInfo retrieves credit information for a broker FD account.
func (c *Client) GetFDCreditInfo(ctx context.Context, accountID string) (*CreditInfo, error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	var out CreditInfo
	if err := c.get(ctx, pathFDCreditInfo, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
