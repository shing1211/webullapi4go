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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

func TestListFDBankAccounts(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]BankAccount{
			{BankID: "B1", AccountID: "A1", BankName: "Bank A", AccountNumber: "123", RoutingNumber: "111", Status: "active"},
			{BankID: "B2", AccountID: "A1", BankName: "Bank B", AccountNumber: "456", RoutingNumber: "222", Status: "active"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.ListFDBankAccounts(context.Background(), "A1")
	if err != nil {
		t.Fatalf("ListFDBankAccounts error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if capturedReq.URL.Path != pathFDBankAccounts {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDBankAccounts)
	}
	if capturedReq.URL.Query().Get("account_id") != "A1" {
		t.Fatalf("account_id = %s, want A1", capturedReq.URL.Query().Get("account_id"))
	}
}

func TestGetFDBankAccountDetail(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(BankAccount{
			BankID: "B1", AccountID: "A1", BankName: "Bank A", AccountNumber: "123", RoutingNumber: "111", Status: "active",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDBankAccountDetail(context.Background(), "B1")
	if err != nil {
		t.Fatalf("GetFDBankAccountDetail error = %v", err)
	}
	if got.BankID != "B1" {
		t.Fatalf("BankID = %s, want B1", got.BankID)
	}
	if capturedReq.URL.Path != pathFDBankAccountDetail {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDBankAccountDetail)
	}
	if capturedReq.URL.Query().Get("bank_id") != "B1" {
		t.Fatalf("bank_id = %s, want B1", capturedReq.URL.Query().Get("bank_id"))
	}
}

func TestAddFDBankAccount(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(BankAccount{
			BankID: "B3", AccountID: "A1", BankName: "Bank C", AccountNumber: "789", RoutingNumber: "333", Status: "active",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.AddFDBankAccount(context.Background(), AddFDBankAccountRequest{
		AccountID: "A1", BankName: "Bank C", AccountNumber: "789", RoutingNumber: "333",
	})
	if err != nil {
		t.Fatalf("AddFDBankAccount error = %v", err)
	}
	if got.BankID != "B3" {
		t.Fatalf("BankID = %s, want B3", got.BankID)
	}
	if capturedReq.URL.Path != pathFDBankAccountAdd {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDBankAccountAdd)
	}
	if capturedReq.Method != http.MethodPost {
		t.Fatalf("method = %s, want %s", capturedReq.Method, http.MethodPost)
	}
}

func TestRemoveFDBankAccount(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	err = c.RemoveFDBankAccount(context.Background(), "B1")
	if err != nil {
		t.Fatalf("RemoveFDBankAccount error = %v", err)
	}
	if capturedReq.URL.Path != pathFDBankAccountRemove {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDBankAccountRemove)
	}
	if capturedReq.URL.Query().Get("bank_id") != "B1" {
		t.Fatalf("bank_id = %s, want B1", capturedReq.URL.Query().Get("bank_id"))
	}
}

func TestListFDAchAccounts(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]ACHAccount{
			{ACHID: "ACH1", AccountID: "A1", BankName: "Bank A", AccountNumber: "123", Status: "active"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.ListFDAchAccounts(context.Background(), "A1")
	if err != nil {
		t.Fatalf("ListFDAchAccounts error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
	if capturedReq.URL.Path != pathFDAchAccounts {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDAchAccounts)
	}
}

func TestGetFDAchAccountDetail(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ACHAccount{
			ACHID: "ACH1", AccountID: "A1", BankName: "Bank A", AccountNumber: "123", Status: "active",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDAchAccountDetail(context.Background(), "ACH1")
	if err != nil {
		t.Fatalf("GetFDAchAccountDetail error = %v", err)
	}
	if got.ACHID != "ACH1" {
		t.Fatalf("ACHID = %s, want ACH1", got.ACHID)
	}
	if capturedReq.URL.Path != pathFDAchAccountDetail {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDAchAccountDetail)
	}
}

func TestListFDTransfers(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]Transfer{
			{TransferID: "T1", AccountID: "A1", Type: "DEPOSIT", Amount: money.Must(money.NewFromString("1000")), Currency: "USD", Status: "completed", CreateTime: "2026-01-01"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.ListFDTransfers(context.Background(), "A1")
	if err != nil {
		t.Fatalf("ListFDTransfers error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
	if capturedReq.URL.Path != pathFDTransfers {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDTransfers)
	}
}

func TestGetFDTransferDetail(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Transfer{
			TransferID: "T1", AccountID: "A1", Type: "DEPOSIT", Amount: money.Must(money.NewFromString("1000")), Currency: "USD", Status: "completed", CreateTime: "2026-01-01",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDTransferDetail(context.Background(), "T1")
	if err != nil {
		t.Fatalf("GetFDTransferDetail error = %v", err)
	}
	if got.TransferID != "T1" {
		t.Fatalf("TransferID = %s, want T1", got.TransferID)
	}
	if capturedReq.URL.Path != pathFDTransferDetail {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDTransferDetail)
	}
}

func TestInitiateFDTransfer(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Transfer{
			TransferID: "T2", AccountID: "A1", Type: "DEPOSIT", Amount: money.Must(money.NewFromString("500")), Currency: "USD", Status: "pending", CreateTime: "2026-01-02",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.InitiateFDTransfer(context.Background(), InitiateTransferRequest{
		AccountID: "A1", Type: "DEPOSIT", Amount: func() *money.Money { m := money.Must(money.NewFromString("500")); return &m }(), Currency: "USD",
	})
	if err != nil {
		t.Fatalf("InitiateFDTransfer error = %v", err)
	}
	if got.TransferID != "T2" {
		t.Fatalf("TransferID = %s, want T2", got.TransferID)
	}
	if capturedReq.URL.Path != pathFDTransferInitiate {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDTransferInitiate)
	}
	if capturedReq.Method != http.MethodPost {
		t.Fatalf("method = %s, want %s", capturedReq.Method, http.MethodPost)
	}
}

func TestCreateFDInstantFunding(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(InstantFunding{
			FundingID: "F1", AccountID: "A1", Amount: money.Must(money.NewFromString("1000")), Status: "pending", CreateTime: "2026-01-01",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.CreateFDInstantFunding(context.Background(), "A1", "1000")
	if err != nil {
		t.Fatalf("CreateFDInstantFunding error = %v", err)
	}
	if got.FundingID != "F1" {
		t.Fatalf("FundingID = %s, want F1", got.FundingID)
	}
	if capturedReq.URL.Path != pathFDInstantFunding {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDInstantFunding)
	}
	if capturedReq.URL.Query().Get("account_id") != "A1" {
		t.Fatalf("account_id = %s, want A1", capturedReq.URL.Query().Get("account_id"))
	}
	if capturedReq.URL.Query().Get("amount") != "1000" {
		t.Fatalf("amount = %s, want 1000", capturedReq.URL.Query().Get("amount"))
	}
}

func TestGetFDInstantFundingDetail(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(InstantFunding{
			FundingID: "F1", AccountID: "A1", Amount: money.Must(money.NewFromString("1000")), Status: "completed", CreateTime: "2026-01-01",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDInstantFundingDetail(context.Background(), "F1")
	if err != nil {
		t.Fatalf("GetFDInstantFundingDetail error = %v", err)
	}
	if got.FundingID != "F1" {
		t.Fatalf("FundingID = %s, want F1", got.FundingID)
	}
	if capturedReq.URL.Path != pathFDInstantFundingDetail {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDInstantFundingDetail)
	}
}

func TestGetFDTransferFees(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]TransferFee{
			{Type: "ACH", Amount: money.Must(money.NewFromString("0")), Currency: "USD"},
			{Type: "WIRE", Amount: money.Must(money.NewFromString("25")), Currency: "USD"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDTransferFees(context.Background())
	if err != nil {
		t.Fatalf("GetFDTransferFees error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if capturedReq.URL.Path != pathFDTransferFees {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDTransferFees)
	}
}

func TestGetFDCreditInfo(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(CreditInfo{
			AccountID: "A1", CreditLimit: money.Must(money.NewFromString("10000")), UsedCredit: money.Must(money.NewFromString("2000")), AvailableCredit: money.Must(money.NewFromString("8000")),
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDCreditInfo(context.Background(), "A1")
	if err != nil {
		t.Fatalf("GetFDCreditInfo error = %v", err)
	}
	if got.AccountID != "A1" {
		t.Fatalf("AccountID = %s, want A1", got.AccountID)
	}
	if capturedReq.URL.Path != pathFDCreditInfo {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDCreditInfo)
	}
	if capturedReq.URL.Query().Get("account_id") != "A1" {
		t.Fatalf("account_id = %s, want A1", capturedReq.URL.Query().Get("account_id"))
	}
}
