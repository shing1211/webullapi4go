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

func TestGetFDActivities(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]FDActivity{
			{ActivityID: "ACT1", AccountID: "A1", Type: "DEPOSIT", Amount: "1000", Currency: "USD", Status: "completed", CreateTime: "2026-01-01", Description: "Deposit"},
			{ActivityID: "ACT2", AccountID: "A1", Type: "WITHDRAWAL", Amount: "500", Currency: "USD", Status: "completed", CreateTime: "2026-01-02", Description: "Withdrawal"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDActivities(context.Background(), "A1")
	if err != nil {
		t.Fatalf("GetFDActivities error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if capturedReq.URL.Path != pathFDActivities {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDActivities)
	}
	if capturedReq.URL.Query().Get("account_id") != "A1" {
		t.Fatalf("account_id = %s, want A1", capturedReq.URL.Query().Get("account_id"))
	}
}

func TestGetFDActivitiesEmpty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]FDActivity{})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDActivities(context.Background(), "A2")
	if err != nil {
		t.Fatalf("GetFDActivities error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("len(got) = %d, want 0", len(got))
	}
}

func TestGetFDActivitiesSingle(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]FDActivity{
			{ActivityID: "ACT3", AccountID: "A1", Type: "DIVIDEND", Amount: "50", Currency: "USD", Status: "completed", CreateTime: "2026-01-03", Description: "Dividend"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDActivities(context.Background(), "A1")
	if err != nil {
		t.Fatalf("GetFDActivities error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
	if got[0].ActivityID != "ACT3" {
		t.Fatalf("ActivityID = %s, want ACT3", got[0].ActivityID)
	}
	if capturedReq.URL.Path != pathFDActivities {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDActivities)
	}
}
