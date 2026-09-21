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
)

func TestListFDAccounts(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]FDAccount{
			{AccountID: "A1", AccountNumber: "123", AccountType: "CASH", AccountClass: "INDIVIDUAL_CASH", Status: "active", Currency: "USD", CreateTime: "2026-01-01"},
			{AccountID: "A2", AccountNumber: "456", AccountType: "MARGIN", AccountClass: "INDIVIDUAL_MRGN", Status: "active", Currency: "USD", CreateTime: "2026-01-02"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.ListFDAccounts(context.Background())
	if err != nil {
		t.Fatalf("ListFDAccounts error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if capturedReq.URL.Path != pathFDAccountList {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDAccountList)
	}
}

func TestGetFDAccountDetail(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(FDAccount{
			AccountID: "A1", AccountNumber: "123", AccountType: "CASH", AccountClass: "INDIVIDUAL_CASH", Status: "active", Currency: "USD", CreateTime: "2026-01-01",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDAccountDetail(context.Background(), "A1")
	if err != nil {
		t.Fatalf("GetFDAccountDetail error = %v", err)
	}
	if got.AccountID != "A1" {
		t.Fatalf("AccountID = %s, want A1", got.AccountID)
	}
	if capturedReq.URL.Path != pathFDAccountDetail {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDAccountDetail)
	}
	if capturedReq.URL.Query().Get("account_id") != "A1" {
		t.Fatalf("account_id = %s, want A1", capturedReq.URL.Query().Get("account_id"))
	}
}

func TestCreateFDAccount(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(FDAccount{
			AccountID: "A3", AccountNumber: "789", AccountType: "CASH", AccountClass: "INDIVIDUAL_CASH", Status: "active", Currency: "USD", CreateTime: "2026-01-03",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.CreateFDAccount(context.Background(), CreateFDAccountRequest{
		AccountType:  "CASH",
		AccountClass: "INDIVIDUAL_CASH",
		Currency:     "USD",
	})
	if err != nil {
		t.Fatalf("CreateFDAccount error = %v", err)
	}
	if got.AccountID != "A3" {
		t.Fatalf("AccountID = %s, want A3", got.AccountID)
	}
	if capturedReq.URL.Path != pathFDAccountCreate {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDAccountCreate)
	}
	if capturedReq.Method != http.MethodPost {
		t.Fatalf("method = %s, want %s", capturedReq.Method, http.MethodPost)
	}
}

func TestUpdateFDAccount(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(FDAccount{
			AccountID: "A1", AccountNumber: "123", AccountType: "MARGIN", AccountClass: "INDIVIDUAL_MRGN", Status: "active", Currency: "USD", CreateTime: "2026-01-01",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.UpdateFDAccount(context.Background(), UpdateFDAccountRequest{
		AccountID:    "A1",
		AccountType:  "MARGIN",
		AccountClass: "INDIVIDUAL_MRGN",
	})
	if err != nil {
		t.Fatalf("UpdateFDAccount error = %v", err)
	}
	if got.AccountType != "MARGIN" {
		t.Fatalf("AccountType = %s, want MARGIN", got.AccountType)
	}
	if capturedReq.URL.Path != pathFDAccountUpdate {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDAccountUpdate)
	}
	if capturedReq.Method != http.MethodPut {
		t.Fatalf("method = %s, want %s", capturedReq.Method, http.MethodPut)
	}
}

func TestCloseFDAccount(t *testing.T) {
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

	err = c.CloseFDAccount(context.Background(), "A1")
	if err != nil {
		t.Fatalf("CloseFDAccount error = %v", err)
	}
	if capturedReq.URL.Path != pathFDAccountClose {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDAccountClose)
	}
	if capturedReq.URL.Query().Get("account_id") != "A1" {
		t.Fatalf("account_id = %s, want A1", capturedReq.URL.Query().Get("account_id"))
	}
}

func TestListAccountForms(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]AccountForm{
			{FormID: "F1", FormType: "W-8BEN", Status: "pending", CreateTime: "2026-01-01"},
			{FormID: "F2", FormType: "W-9", Status: "completed", CreateTime: "2026-01-02"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.ListAccountForms(context.Background(), "A1")
	if err != nil {
		t.Fatalf("ListAccountForms error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if capturedReq.URL.Path != pathFDAccountForms {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDAccountForms)
	}
	if capturedReq.URL.Query().Get("account_id") != "A1" {
		t.Fatalf("account_id = %s, want A1", capturedReq.URL.Query().Get("account_id"))
	}
}

func TestGetAccountFormDetail(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(AccountForm{
			FormID: "F1", FormType: "W-8BEN", Status: "pending", CreateTime: "2026-01-01", Content: "full form content",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetAccountFormDetail(context.Background(), "F1")
	if err != nil {
		t.Fatalf("GetAccountFormDetail error = %v", err)
	}
	if got.FormID != "F1" {
		t.Fatalf("FormID = %s, want F1", got.FormID)
	}
	if capturedReq.URL.Path != pathFDAccountFormDetail {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDAccountFormDetail)
	}
	if capturedReq.URL.Query().Get("form_id") != "F1" {
		t.Fatalf("form_id = %s, want F1", capturedReq.URL.Query().Get("form_id"))
	}
}

func TestSubmitAccountForm(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(AccountForm{
			FormID: "F1", FormType: "W-8BEN", Status: "submitted", CreateTime: "2026-01-01",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.SubmitAccountForm(context.Background(), SubmitAccountFormRequest{
		FormID:  "F1",
		Content: "completed form data",
	})
	if err != nil {
		t.Fatalf("SubmitAccountForm error = %v", err)
	}
	if got.Status != "submitted" {
		t.Fatalf("Status = %s, want submitted", got.Status)
	}
	if capturedReq.URL.Path != pathFDAccountFormSubmit {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDAccountFormSubmit)
	}
	if capturedReq.Method != http.MethodPost {
		t.Fatalf("method = %s, want %s", capturedReq.Method, http.MethodPost)
	}
}

func TestGetAccountFormStatus(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(AccountForm{
			FormID: "F1", FormType: "W-8BEN", Status: "approved", CreateTime: "2026-01-01",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetAccountFormStatus(context.Background(), "F1")
	if err != nil {
		t.Fatalf("GetAccountFormStatus error = %v", err)
	}
	if got.Status != "approved" {
		t.Fatalf("Status = %s, want approved", got.Status)
	}
	if capturedReq.URL.Path != pathFDAccountFormStatus {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDAccountFormStatus)
	}
	if capturedReq.URL.Query().Get("form_id") != "F1" {
		t.Fatalf("form_id = %s, want F1", capturedReq.URL.Query().Get("form_id"))
	}
}
