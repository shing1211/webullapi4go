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

package broker_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shing1211/webullapi4go/broker"
)

func TestListVirtualAccounts(t *testing.T) {
	t.Parallel()

	const body = `[{"account_id":"VA1","account_name":"Paper",` +
		`"account_type":"PAPER","currency":"USD","status":"ACTIVE",` +
		`"create_time":"2026-01-01","update_time":"2026-01-02"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/openapi/v1/broker/virtual-accounts"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	got, err := c.ListVirtualAccounts(context.Background())
	if err != nil {
		t.Fatalf("ListVirtualAccounts() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d accounts, want 1", len(got))
	}
	if got[0].AccountID != "VA1" || got[0].AccountName != "Paper" {
		t.Errorf("account = %+v, want VA1/Paper", got[0])
	}
}

func TestGetVirtualAccount(t *testing.T) {
	t.Parallel()

	const body = `{"account_id":"VA1","account_name":"Paper",` +
		`"account_type":"PAPER","currency":"USD","status":"ACTIVE",` +
		`"create_time":"2026-01-01","update_time":"2026-01-02"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/openapi/v1/broker/virtual-accounts/detail"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("account_id"), "VA1"; got != want {
			t.Errorf("account_id = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	got, err := c.GetVirtualAccount(context.Background(), "VA1")
	if err != nil {
		t.Fatalf("GetVirtualAccount() error = %v", err)
	}
	if got.AccountID != "VA1" || got.Currency != "USD" {
		t.Errorf("account = %+v", got)
	}
}

func TestCreateVirtualAccount(t *testing.T) {
	t.Parallel()

	const respBody = `{"account_id":"VA2","account_name":"Trading",` +
		`"account_type":"PAPER","currency":"HKD","status":"ACTIVE",` +
		`"create_time":"2026-02-01","update_time":"2026-02-01"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/openapi/v1/broker/virtual-accounts"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var req broker.CreateVirtualAccountRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if req.AccountName != "Trading" {
			t.Errorf("account_name = %q, want Trading", req.AccountName)
		}
		if req.Currency != "HKD" {
			t.Errorf("currency = %q, want HKD", req.Currency)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respBody))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	got, err := c.CreateVirtualAccount(context.Background(), broker.CreateVirtualAccountRequest{
		AccountName: "Trading",
		AccountType: "PAPER",
		Currency:    "HKD",
	})
	if err != nil {
		t.Fatalf("CreateVirtualAccount() error = %v", err)
	}
	if got.AccountID != "VA2" || got.AccountName != "Trading" {
		t.Errorf("account = %+v", got)
	}
}

func TestUpdateVirtualAccount(t *testing.T) {
	t.Parallel()

	const respBody = `{"account_id":"VA1","account_name":"Renamed",` +
		`"account_type":"PAPER","currency":"USD","status":"ACTIVE",` +
		`"create_time":"2026-01-01","update_time":"2026-03-01"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %q, want PUT", r.Method)
		}
		if got, want := r.URL.Path, "/openapi/v1/broker/virtual-accounts/detail"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("account_id"), "VA1"; got != want {
			t.Errorf("account_id = %q, want %q", got, want)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var req broker.UpdateVirtualAccountRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if req.AccountName != "Renamed" {
			t.Errorf("account_name = %q, want Renamed", req.AccountName)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respBody))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	got, err := c.UpdateVirtualAccount(context.Background(), "VA1", broker.UpdateVirtualAccountRequest{
		AccountName: "Renamed",
	})
	if err != nil {
		t.Fatalf("UpdateVirtualAccount() error = %v", err)
	}
	if got.AccountID != "VA1" || got.AccountName != "Renamed" {
		t.Errorf("account = %+v", got)
	}
}
