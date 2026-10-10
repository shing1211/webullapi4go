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

	const body = `{"data":[{"account_id":"VA1","account_number":"Paper",` +
		`"account_type":"PAPER","account_status":"ACTIVE",` +
		`"belong_account_id":"ACC1","belong_account_number":"Main",` +
		`"trading_permissions":["STOCK","OPTION"],"option_level":"Level 1",` +
		`"commission_code":"C001","w8ben_info":"","china_connect_investor_info":""}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/broker/accounts/virtual-accounts/list"; got != want {
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
	if got[0].AccountID != "VA1" || got[0].AccountNumber != "Paper" {
		t.Errorf("account = %+v, want VA1/Paper", got[0])
	}
}

func TestGetVirtualAccount(t *testing.T) {
	t.Parallel()

	const body = `{"account_id":"VA1","account_number":"Paper",` +
		`"account_type":"PAPER","account_status":"ACTIVE",` +
		`"belong_account_id":"ACC1","belong_account_number":"Main",` +
		`"trading_permissions":["STOCK"],"option_level":"Level 1",` +
		`"commission_code":"C001","w8ben_info":"","china_connect_investor_info":""}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/broker/accounts/virtual-accounts/get"; got != want {
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
	if got.AccountID != "VA1" || got.AccountNumber != "Paper" {
		t.Errorf("account = %+v", got)
	}
}

func TestCreateVirtualAccount(t *testing.T) {
	t.Parallel()

	const respBody = `{"account_id":"VA2","account_number":"Trading",` +
		`"account_type":"PAPER","account_status":"ACTIVE",` +
		`"belong_account_id":"ACC1","belong_account_number":"Main",` +
		`"trading_permissions":["STOCK"],"option_level":"Level 1",` +
		`"commission_code":"C002","w8ben_info":"","china_connect_investor_info":""}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/broker/accounts/virtual-accounts/create"; got != want {
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
	if got.AccountID != "VA2" || got.AccountNumber != "Trading" {
		t.Errorf("account = %+v", got)
	}
}

func TestUpdateVirtualAccount(t *testing.T) {
	t.Parallel()

	const respBody = `{"account_id":"VA1","account_number":"Renamed",` +
		`"account_type":"PAPER","account_status":"ACTIVE",` +
		`"belong_account_id":"ACC1","belong_account_number":"Main",` +
		`"trading_permissions":["STOCK"],"option_level":"Level 1",` +
		`"commission_code":"C001","w8ben_info":"","china_connect_investor_info":""}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// POST, not PUT: the page documents POST for this path, and the previous
		// implementation's PUT was the defect this test now pins.
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/broker/accounts/virtual-accounts/update"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		// account_id belongs in the body. It was a query parameter before, and the
		// page lists no query parameter other than the auth headers.
		if got := r.URL.Query().Get("account_id"); got != "" {
			t.Errorf("account_id in query = %q, want it in the body instead", got)
		}
		if got := r.URL.RawQuery; got != "" {
			t.Errorf("RawQuery = %q, want empty: the documented request takes no query parameters", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var req broker.UpdateVirtualAccountRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if req.ClientRequestID != "req-1" {
			t.Errorf("client_request_id = %q, want req-1", req.ClientRequestID)
		}
		// The body is decoded into the unexported wire shape for the required pair,
		// because account_id arrives from the argument rather than the request.
		var wire map[string]any
		if err := json.Unmarshal(body, &wire); err != nil {
			t.Fatalf("unmarshal body as object: %v", err)
		}
		if got, _ := wire["account_id"].(string); got != "VA1" {
			t.Errorf("account_id in body = %q, want VA1", got)
		}
		if got, _ := wire["option_level"].(string); got != "LV1" {
			t.Errorf("option_level = %q, want LV1", got)
		}
		if perms, _ := wire["trading_permissions"].([]any); len(perms) != 1 {
			t.Errorf("trading_permissions = %v, want one entry", wire["trading_permissions"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respBody))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	got, err := c.UpdateVirtualAccount(context.Background(), "VA1", broker.UpdateVirtualAccountRequest{
		ClientRequestID:    "req-1",
		TradingPermissions: []string{"US_STOCK_NORMAL"},
		OptionLevel:        "LV1",
	})
	if err != nil {
		t.Fatalf("UpdateVirtualAccount() error = %v", err)
	}
	if got.AccountID != "VA1" || got.AccountNumber != "Renamed" {
		t.Errorf("account = %+v", got)
	}
}

// TestUpdateVirtualAccountOmitsUnsetOptionals guards the encoding choice: an unset
// optional field is omitted, not sent as an empty string, because a field the page
// does not require is meaningful when absent and misleading when blank.
func TestUpdateVirtualAccountOmitsUnsetOptionals(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var wire map[string]any
		if err := json.Unmarshal(body, &wire); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		for _, absent := range []string{
			"trading_permissions", "option_level", "commission_code",
			"w8ben_info", "china_connect_investor_info",
		} {
			if v, ok := wire[absent]; ok {
				t.Errorf("%s = %v present, want omitted", absent, v)
			}
		}
		// The two the page requires must still be there.
		for _, required := range []string{"account_id", "client_request_id"} {
			if _, ok := wire[required]; !ok {
				t.Errorf("%s missing, want present", required)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"account_id":"VA1","account_number":"N1"}`))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	if _, err := c.UpdateVirtualAccount(context.Background(), "VA1", broker.UpdateVirtualAccountRequest{
		ClientRequestID: "req-1",
	}); err != nil {
		t.Fatalf("UpdateVirtualAccount() error = %v", err)
	}
}
