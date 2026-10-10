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
	"encoding/json"
	"testing"

	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

// 58 required names across 18 Broker FD symbols reached no field on their response
// types, so every documented value decoded to the zero value and no error was
// reported. The bodies below are the documented shapes. They are decoded here rather
// than round-tripped from the SDK's own types, because a round-trip is green whether
// the type is right or wrong: that blindness is how these names went missing with a
// green suite, and it is why each test asserts the value the page sends reaches a
// field rather than that a field survives a marshal.

// documentedBankAccount is the response of POST /broker/funding/bank-relationships/create
// and of GET /broker/funding/bank-relationships/list.
const documentedBankAccount = `{
	"bank_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
	"account_id": "93IUJ28O9VO2KBGHDHR4H9",
	"bank_name": "Bank of America",
	"account_number": "123456789",
	"routing_number": "021000021",
	"status": "ACTIVE",
	"client_request_id": "9d8f1a2b",
	"create_time": "2025-01-05T22:59:59.012Z",
	"update_time": "2025-01-05T22:59:59.012Z",
	"bank_relationship_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
	"bank_account_name": "ABA",
	"bank_account_number": "123456789",
	"bank_routing_number": "021000021",
	"bank_code": "ABA",
	"bank_code_type": "ABA"
}`

// documentedACHAccount is the response of POST /broker/funding/ach-relationships/create
// and of GET /broker/funding/ach-relationships/list.
const documentedACHAccount = `{
	"ach_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
	"account_id": "93IUJ28O9VO2KBGHDHR4H9",
	"bank_name": "Bank of America",
	"account_number": "123456789",
	"status": "ACTIVE",
	"client_request_id": "9d8f1a2b",
	"create_time": "2025-01-05T22:59:59.012Z",
	"update_time": "2025-01-05T22:59:59.012Z",
	"ach_relationship_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
	"account_owner_name": "Jane Doe",
	"bank_account_number": "123456789",
	"bank_account_type": "CHECKING",
	"bank_routing_number": "021000021"
}`

// TestBankAccountReadsBothSpellings is the dual-spelling case. This type has always
// carried bank_id, account_number and routing_number while the page sends
// bank_relationship_id, bank_account_number and bank_routing_number for the same
// facts, so a body carrying only the documented names used to decode to empty strings
// on those three. Both sets are carried, so each reads the value the server sends.
func TestBankAccountReadsBothSpellings(t *testing.T) {
	var got BankAccount
	if err := json.Unmarshal([]byte(documentedBankAccount), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for name, pair := range map[string][2]string{
		"BankRelationshipID": {got.BankRelationshipID, "89JGKD4LKIVI5UU6L3KHNU10IA"},
		"BankAccountName":    {got.BankAccountName, "ABA"},
		"BankAccountNumber":  {got.BankAccountNumber, "123456789"},
		"BankRoutingNumber":  {got.BankRoutingNumber, "021000021"},
		"BankCode":           {got.BankCode, "ABA"},
		"BankCodeType":       {got.BankCodeType, "ABA"},
	} {
		if pair[0] != pair[1] {
			t.Errorf("%s = %q, want %q", name, pair[0], pair[1])
		}
	}
	// The SDK's own spellings are untouched by a body that carries only the
	// documented ones, which is the point of carrying both: whichever a server
	// sends, a caller reading either name gets a value.
	if got.BankName != "Bank of America" || got.Status != "ACTIVE" {
		t.Errorf("pre-existing fields = %q/%q, want them still populated",
			got.BankName, got.Status)
	}
}

func TestACHAccountReadsTheDocumentedNames(t *testing.T) {
	var got ACHAccount
	if err := json.Unmarshal([]byte(documentedACHAccount), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for name, pair := range map[string][2]string{
		"ACHRelationshipID": {got.ACHRelationshipID, "89JGKD4LKIVI5UU6L3KHNU10IA"},
		"AccountOwnerName":  {got.AccountOwnerName, "Jane Doe"},
		"BankAccountNumber": {got.BankAccountNumber, "123456789"},
		"BankAccountType":   {got.BankAccountType, "CHECKING"},
		"BankRoutingNumber": {got.BankRoutingNumber, "021000021"},
	} {
		if pair[0] != pair[1] {
			t.Errorf("%s = %q, want %q", name, pair[0], pair[1])
		}
	}
}

// TestCreditInfoReadsEveryDocumentedName covers the largest single batch. The page
// requires ten names and the type carried seven, so a caller could not read which
// account a credit applied to or how much was involved.
func TestCreditInfoReadsEveryDocumentedName(t *testing.T) {
	var got CreditInfo
	if err := json.Unmarshal([]byte(`{
		"account_id": "93IUJ28O9VO2KBGHDHR4H9",
		"credit_limit": "10000",
		"used_credit": "100",
		"available_credit": "9900",
		"client_request_id": "9d8f1a2b",
		"create_time": "2025-01-05T22:59:59.012Z",
		"update_time": "2025-01-05T22:59:59.012Z",
		"credit_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
		"account": "93IUJ28O9VO2KBGHDHR4H9",
		"contra_account": "41IO9QG4M5O65B0EA4LSJ4UJ99",
		"type": "GIFTING",
		"status": "SUBMITTED",
		"amount": "100",
		"currency": "USD"
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.CreditID != "89JGKD4LKIVI5UU6L3KHNU10IA" {
		t.Errorf("CreditID = %q", got.CreditID)
	}
	if got.Account != got.AccountID {
		t.Errorf("Account = %q, want the same account as AccountID (%q): the page sends "+
			"account and the type spelled account_id, so both must read",
			got.Account, got.AccountID)
	}
	if got.ContraAccount != "41IO9QG4M5O65B0EA4LSJ4UJ99" {
		t.Errorf("ContraAccount = %q", got.ContraAccount)
	}
	if got.Type != "GIFTING" || got.Status != "SUBMITTED" || got.Currency != "USD" {
		t.Errorf("Type/Status/Currency = %q/%q/%q", got.Type, got.Status, got.Currency)
	}
	if want := money.Must(money.NewFromString("100")); got.Amount.Cmp(want) != 0 {
		t.Errorf("Amount = %s, want %s", got.Amount, want)
	}
}

func TestInstantFundingReadsTheDocumentedNames(t *testing.T) {
	var got InstantFunding
	if err := json.Unmarshal([]byte(`{
		"funding_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
		"account_id": "93IUJ28O9VO2KBGHDHR4H9",
		"amount": "100",
		"status": "SUBMITTED",
		"create_time": "2025-01-05T22:59:59.012Z",
		"client_request_id": "9d8f1a2b",
		"update_time": "2025-01-05T22:59:59.012Z",
		"instant_funding_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
		"currency": "USD",
		"type": "DEPOSIT"
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.InstantFundingID != got.FundingID {
		t.Errorf("InstantFundingID = %q, want the same value as FundingID (%q)",
			got.InstantFundingID, got.FundingID)
	}
	if got.Currency != "USD" || got.Type != "DEPOSIT" {
		t.Errorf("Currency/Type = %q/%q, want USD/DEPOSIT", got.Currency, got.Type)
	}
}

func TestTransferReadsTheDocumentedNames(t *testing.T) {
	var got Transfer
	if err := json.Unmarshal([]byte(`{
		"transfer_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
		"account_id": "93IUJ28O9VO2KBGHDHR4H9",
		"type": "DEPOSIT",
		"amount": "100",
		"currency": "USD",
		"status": "SUBMITTED",
		"create_time": "2025-01-05T22:59:59.012Z",
		"client_request_id": "9d8f1a2b",
		"update_time": "2025-01-05T22:59:59.012Z",
		"transfer_type": "ACH",
		"direction": "IN"
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.TransferType != "ACH" {
		t.Errorf("TransferType = %q, want ACH: the type spells the same fact Type, so "+
			"a caller could not tell an ACH transfer from any other kind", got.TransferType)
	}
	if got.Direction != "IN" {
		t.Errorf("Direction = %q, want IN: without it a transfer cannot be read as "+
			"inbound or outbound", got.Direction)
	}
}

func TestTransferFeeReadsTheDocumentedNames(t *testing.T) {
	var got TransferFee
	if err := json.Unmarshal([]byte(`{
		"type": "SERVICE_FEE",
		"amount": "1",
		"currency": "USD",
		"client_request_id": "9d8f1a2b",
		"create_time": "2025-01-05T22:59:59.012Z",
		"update_time": "2025-01-05T22:59:59.012Z",
		"fee_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
		"account": "93IUJ28O9VO2KBGHDHR4H9",
		"contra_account": "41IO9QG4M5O65B0EA4LSJ4UJ99",
		"status": "SUBMITTED"
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.FeeID == "" || got.Account == "" || got.ContraAccount == "" || got.Status == "" {
		t.Errorf("fee identity fields = %q/%q/%q/%q, all four are required by the page",
			got.FeeID, got.Account, got.ContraAccount, got.Status)
	}
}

func TestAgreementReadsTheDocumentedNames(t *testing.T) {
	var got Agreement
	if err := json.Unmarshal([]byte(`{
		"agreement_id": "JPGSRHN6QCAI68SEO24UL316TB_1",
		"agreement_type": "CUSTOMER",
		"version": "1",
		"status": "ACTIVE",
		"effective_date": "2025-01-05",
		"content": "text",
		"id": "JPGSRHN6QCAI68SEO24UL316TB_1",
		"name": "Customer Agreement",
		"content_type": "TEXT"
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.ID != got.AgreementID {
		t.Errorf("ID = %q, want the same value as AgreementID (%q)", got.ID, got.AgreementID)
	}
	if got.Name != "Customer Agreement" {
		t.Errorf("Name = %q: without it a caller could not tell which agreement this is", got.Name)
	}
	if got.ContentType != "TEXT" {
		t.Errorf("ContentType = %q, want TEXT", got.ContentType)
	}
}

func TestFDAccountReadsApplicationID(t *testing.T) {
	var got FDAccount
	if err := json.Unmarshal([]byte(`{
		"account_id": "93IUJ28O9VO2KBGHDHR4H9",
		"account_number": "123456789",
		"account_type": "FRACTIONAL",
		"account_class": "CASH",
		"status": "ACTIVE",
		"currency": "USD",
		"create_time": "2025-01-05T22:59:59.012Z",
		"client_request_id": "9d8f1a2b",
		"application_id": "4F5A13B5FDFD43D29C7A3A856B2458A2"
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.ApplicationID != "4F5A13B5FDFD43D29C7A3A856B2458A2" {
		t.Errorf("ApplicationID = %q, want the documented value: it is how an "+
			"application tells an account it opened from one it only acts for",
			got.ApplicationID)
	}
}

func TestFDCashJournalReadsBothCounterparties(t *testing.T) {
	var got FDCashJournal
	if err := json.Unmarshal([]byte(`{
		"journal_id": "89JGKD4LKIVI5UU6L3KHNU10IA",
		"account_id": "93IUJ28O9VO2KBGHDHR4H9",
		"type": "DEPOSIT",
		"amount": "100",
		"currency": "USD",
		"status": "POSTED",
		"create_time": "2025-01-05T22:59:59.012Z",
		"client_request_id": "9d8f1a2b",
		"update_time": "2025-01-05T22:59:59.012Z",
		"from_account": "93IUJ28O9VO2KBGHDHR4H9",
		"to_account": "41IO9QG4M5O65B0EA4LSJ4UJ99"
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.FromAccount != "93IUJ28O9VO2KBGHDHR4H9" || got.ToAccount != "41IO9QG4M5O65B0EA4LSJ4UJ99" {
		t.Errorf("FromAccount/ToAccount = %q/%q, want both: they are what distinguishes "+
			"a deposit from a withdrawal, and AccountID is neither of them",
			got.FromAccount, got.ToAccount)
	}
}

func TestFDTradeCalendarEntryReadsTheBooleans(t *testing.T) {
	var got []FDTradeCalendarEntry
	if err := json.Unmarshal([]byte(`[{"date":"2025-11-28","market":"US","status":"CLOSED",
		"open_time":"09:30:00","close_time":"16:00:00",
		"is_trading_day":false,"is_settlement_day":false},
		{"date":"2025-11-26","market":"US","status":"OPEN",
		"open_time":"09:30:00","close_time":"13:00:00",
		"is_trading_day":true,"is_settlement_day":true}]`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].IsTradingDay || got[0].IsSettlementDay {
		t.Errorf("day 1 read as trading=%v settlement=%v, want both false: a half-day "+
			"or a holiday cannot be told from a status word", got[0].IsTradingDay, got[0].IsSettlementDay)
	}
	if !got[1].IsTradingDay || !got[1].IsSettlementDay {
		t.Errorf("day 2 read as trading=%v settlement=%v, want both true",
			got[1].IsTradingDay, got[1].IsSettlementDay)
	}
}

func TestFDOrderPreviewReadsTheCostSplit(t *testing.T) {
	var got FDOrderPreview
	if err := json.Unmarshal([]byte(`{
		"order_id": "0352U72LQI6DT0KF41GK000000",
		"estimated_fee": "1",
		"estimated_total": "101",
		"estimated_cost": "100",
		"estimated_transaction_fee": "1"
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if want := money.Must(money.NewFromString("100")); got.EstimatedCost.Cmp(want) != 0 {
		t.Errorf("EstimatedCost = %s, want %s", got.EstimatedCost, want)
	}
	if want := money.Must(money.NewFromString("1")); got.EstimatedTransactionFee.Cmp(want) != 0 {
		t.Errorf("EstimatedTransactionFee = %s, want %s: the page splits the quote "+
			"into base cost and fee, and a caller could not see the split before",
			got.EstimatedTransactionFee, want)
	}
}
