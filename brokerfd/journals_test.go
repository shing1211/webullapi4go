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

func TestListFDCashJournals(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]FDCashJournal{
			{JournalID: "J1", AccountID: "A1", Type: "DEPOSIT", Amount: money.Must(money.NewFromString("1000")), Currency: "USD", Status: "completed", CreateTime: "2026-01-01"},
			{JournalID: "J2", AccountID: "A1", Type: "WITHDRAWAL", Amount: money.Must(money.NewFromString("500")), Currency: "USD", Status: "completed", CreateTime: "2026-01-02"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.ListFDCashJournals(context.Background(), "A1")
	if err != nil {
		t.Fatalf("ListFDCashJournals error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if capturedReq.URL.Path != pathFDCashJournal {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDCashJournal)
	}
	if capturedReq.URL.Query().Get("account_id") != "A1" {
		t.Fatalf("account_id = %s, want A1", capturedReq.URL.Query().Get("account_id"))
	}
}

func TestGetFDCashJournalDetail(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(FDCashJournal{
			JournalID: "J1", AccountID: "A1", Type: "DEPOSIT", Amount: money.Must(money.NewFromString("1000")), Currency: "USD", Status: "completed", CreateTime: "2026-01-01",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDCashJournalDetail(context.Background(), "J1")
	if err != nil {
		t.Fatalf("GetFDCashJournalDetail error = %v", err)
	}
	if got.JournalID != "J1" {
		t.Fatalf("JournalID = %s, want J1", got.JournalID)
	}
	if capturedReq.URL.Path != pathFDCashJournalDetail {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDCashJournalDetail)
	}
	if capturedReq.URL.Query().Get("journal_id") != "J1" {
		t.Fatalf("journal_id = %s, want J1", capturedReq.URL.Query().Get("journal_id"))
	}
}
