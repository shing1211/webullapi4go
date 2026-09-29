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
	"github.com/shing1211/webullapi4go/pkg/types"
)

const (
	pathFDBankAccounts         = "/broker/funding/bank-relationships/list"
	pathFDBankAccountDetail    = "/broker-fd/funding/bank-account/detail"
	pathFDBankAccountAdd       = "/broker/funding/bank-relationships/create"
	pathFDBankAccountRemove    = "/broker/funding/bank-relationships/delete"
	pathFDAchAccounts          = "/broker/funding/ach-relationships/list"
	pathFDAchAccountDetail     = "/broker-fd/funding/ach-account/detail"
	pathFDAchAccountAdd        = "/broker/funding/ach-relationships/create"
	pathFDAchAccountRemove     = "/broker/funding/ach-relationships/delete"
	pathFDTransfers            = "/broker/funding/transfers/list"
	pathFDTransferDetail       = "/broker/funding/transfers/get"
	pathFDTransferInitiate     = "/broker/funding/transfers/create"
	pathFDInstantFunding       = "/broker/funding/instant-funding/create"
	pathFDInstantFundingDetail = "/broker/funding/instant-funding/get"
	pathFDTransferFees         = "/broker/fees/get"
	pathFDCreditInfo           = "/broker/credits/get"
)

// BankAccount represents a linked bank account for funding operations.
type BankAccount struct {
	BankID        string `json:"bank_id"`
	AccountID     string `json:"account_id"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	RoutingNumber string `json:"routing_number"`
	Status        string `json:"status"`

	// The three fields below are required by the page and are entity metadata on this
	// object rather than part of an envelope. The sibling broker/ module carries
	// ClientRequestID in the same role, so the spelling follows it.
	ClientRequestID string `json:"client_request_id"`
	CreateTime      string `json:"create_time"`
	UpdateTime      string `json:"update_time"`

	// The six fields below are also required by the page and were absent here, so they
	// decoded to the zero value with no error reported. The type's own names for three
	// of the same facts are BankID, AccountNumber and RoutingNumber, and the page's are
	// BankRelationshipID, BankAccountNumber and BankRoutingNumber; which of the two a
	// live server sends is unverified, because the Broker FD host returns 404 in the HK
	// sandbox and no US credential was available here. Both sets are therefore carried
	// side by side, and a response populates whichever it carries, so a caller reading
	// either name reads a value when the server sends it. BankAccountName, BankCode and
	// BankCodeType have no SDK-spelled counterpart and are new.
	BankRelationshipID string `json:"bank_relationship_id"`
	BankAccountName    string `json:"bank_account_name"`
	BankAccountNumber  string `json:"bank_account_number"`
	BankRoutingNumber  string `json:"bank_routing_number"`
	BankCode           string `json:"bank_code"`
	BankCodeType       string `json:"bank_code_type"`
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

	// The three fields below are required by the page and are entity metadata on this
	// object rather than part of an envelope. The sibling broker/ module carries
	// ClientRequestID in the same role, so the spelling follows it.
	ClientRequestID string `json:"client_request_id"`
	CreateTime      string `json:"create_time"`
	UpdateTime      string `json:"update_time"`

	// The five fields below are also required by the page and were absent here, so they
	// decoded to the zero value with no error reported. The type's own names for two of
	// the same facts are ACHID and AccountNumber, and the page's are
	// ACHRelationshipID and BankAccountNumber; which of the two a live server sends is
	// unverified, for the reason given on [BankAccount.BankRelationshipID]. Both sets
	// are carried side by side. AccountOwnerName, BankAccountType and
	// BankRoutingNumber have no SDK-spelled counterpart and are new.
	ACHRelationshipID string `json:"ach_relationship_id"`
	AccountOwnerName  string `json:"account_owner_name"`
	BankAccountNumber string `json:"bank_account_number"`
	BankAccountType   string `json:"bank_account_type"`
	BankRoutingNumber string `json:"bank_routing_number"`
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
	TransferID string      `json:"transfer_id"`
	AccountID  string      `json:"account_id"`
	Type       string      `json:"type"`
	Amount     money.Money `json:"amount"`
	Currency   string      `json:"currency"`
	Status     string      `json:"status"`
	CreateTime string      `json:"create_time"`

	// ClientRequestID and UpdateTime are required by the page and are entity
	// metadata on this object. The sibling broker/ module carries
	// ClientRequestID in the same role.
	ClientRequestID string `json:"client_request_id"`
	UpdateTime      string `json:"update_time"`

	// The two fields below are also required by the page and were absent here, so they
	// decoded to the zero value with no error reported. The page's TransferType is the
	// same fact this type spells Type, and which of the two names a live server sends
	// is unverified for the reason given on [BankAccount.BankRelationshipID], so both
	// are carried. Direction has no SDK-spelled counterpart and is new.
	TransferType string `json:"transfer_type"`
	Direction    string `json:"direction"`
}

// ListFDTransfers retrieves all fund transfers for a broker FD account.
func (c *Client) ListFDTransfers(ctx context.Context, accountID, paginationKey string) (*types.Page[Transfer], error) {
	q := url.Values{}
	q.Set("account_id", accountID)
	if paginationKey != "" {
		q.Set("pagination_key", paginationKey)
	}
	var out types.Page[Transfer]
	if err := c.get(ctx, pathFDTransfers, q, &out); err != nil {
		return nil, err
	}
	return &out, nil
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
	AccountID string       `json:"account_id"`
	Type      string       `json:"type"`
	Amount    *money.Money `json:"amount"`
	Currency  string       `json:"currency"`
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
	FundingID  string      `json:"funding_id"`
	AccountID  string      `json:"account_id"`
	Amount     money.Money `json:"amount"`
	Status     string      `json:"status"`
	CreateTime string      `json:"create_time"`

	// ClientRequestID and UpdateTime are required by the page and are entity
	// metadata on this object. The sibling broker/ module carries
	// ClientRequestID in the same role.
	ClientRequestID string `json:"client_request_id"`
	UpdateTime      string `json:"update_time"`

	// The three fields below are also required by the page and were absent here, so
	// they decoded to the zero value with no error reported. The page's
	// InstantFundingID is the same fact this type spells FundingID, and which of the
	// two names a live server sends is unverified for the reason given on
	// [BankAccount.BankRelationshipID], so both are carried. Currency and Type have no
	// SDK-spelled counterpart and are new.
	InstantFundingID string `json:"instant_funding_id"`
	Currency         string `json:"currency"`
	Type             string `json:"type"`
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
	Type     string      `json:"type"`
	Amount   money.Money `json:"amount"`
	Currency string      `json:"currency"`

	// The three fields below are required by the page and are entity metadata on this
	// object rather than part of an envelope. The sibling broker/ module carries
	// ClientRequestID in the same role, so the spelling follows it.
	ClientRequestID string `json:"client_request_id"`
	CreateTime      string `json:"create_time"`
	UpdateTime      string `json:"update_time"`

	// The four fields below are also required by the page and were absent here, so they
	// decoded to the zero value with no error reported. Nothing here is
	// live-verified: the Broker FD host returns 404 in the HK sandbox and no US
	// credential was available here.
	//
	// GetFDTransferFees returns a slice of this type while the page documents a single
	// object, so the array-against-object question is left recorded rather than decided
	// here. These four names describe the fee entity itself, so they are readable
	// whichever shape the server sends.
	FeeID         string `json:"fee_id"`
	Account       string `json:"account"`
	ContraAccount string `json:"contra_account"`
	Status        string `json:"status"`
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
	AccountID       string      `json:"account_id"`
	CreditLimit     money.Money `json:"credit_limit"`
	UsedCredit      money.Money `json:"used_credit"`
	AvailableCredit money.Money `json:"available_credit"`

	// The three fields below are required by the page and are entity metadata on this
	// object rather than part of an envelope. The sibling broker/ module carries
	// ClientRequestID in the same role, so the spelling follows it.
	ClientRequestID string `json:"client_request_id"`
	CreateTime      string `json:"create_time"`
	UpdateTime      string `json:"update_time"`

	// The seven fields below are also required by the page and were absent here, so
	// they decoded to the zero value with no error reported. This type's AccountID is
	// the same fact the page spells Account, and its Amount and Currency have no
	// SDK-spelled counterpart on a credit record. Nothing here is live-verified, for
	// the reason given on [TransferFee.FeeID].
	CreditID      string      `json:"credit_id"`
	Account       string      `json:"account"`
	ContraAccount string      `json:"contra_account"`
	Type          string      `json:"type"`
	Status        string      `json:"status"`
	Amount        money.Money `json:"amount"`
	Currency      string      `json:"currency"`
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
