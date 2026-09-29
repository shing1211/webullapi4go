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

func TestGetFDActivities(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []FDActivity{
				{ActivityID: "ACT1", AccountID: "A1", Type: "DEPOSIT", Amount: money.Must(money.NewFromString("1000")), Currency: "USD", Status: "completed", CreateTime: "2026-01-01", Description: "Deposit"},
				{ActivityID: "ACT2", AccountID: "A1", Type: "WITHDRAWAL", Amount: money.Must(money.NewFromString("500")), Currency: "USD", Status: "completed", CreateTime: "2026-01-02", Description: "Withdrawal"},
			},
			"pagination_key": "eyJ2IjoxLCJsYXN0SWQiOiIwIiwicGFnZU9mZnNldCI6MX0=",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithEndpoints(client.Endpoints{HTTP: srv.URL, BrokerHTTP: srv.URL}), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDActivities(context.Background(), "A1", "")
	if err != nil {
		t.Fatalf("GetFDActivities error = %v", err)
	}
	if len(got.Data) != 2 {
		t.Fatalf("len(got.Data) = %d, want 2", len(got.Data))
	}
	if capturedReq.URL.Path != pathFDActivities {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDActivities)
	}
	if capturedReq.URL.Query().Get("account_id") != "A1" {
		t.Fatalf("account_id = %s, want A1", capturedReq.URL.Query().Get("account_id"))
	}
	if got.PaginationKey != "eyJ2IjoxLCJsYXN0SWQiOiIwIiwicGFnZU9mZnNldCI6MX0=" {
		t.Errorf("PaginationKey = %q, want the documented cursor: without it the endpoint cannot be paged",
			got.PaginationKey)
	}
}

func TestGetFDActivitiesEmpty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":           []FDActivity{},
			"pagination_key": "eyJ2IjoxLCJsYXN0SWQiOiIwIiwicGFnZU9mZnNldCI6MX0=",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithEndpoints(client.Endpoints{HTTP: srv.URL, BrokerHTTP: srv.URL}), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDActivities(context.Background(), "A2", "")
	if err != nil {
		t.Fatalf("GetFDActivities error = %v", err)
	}
	if len(got.Data) != 0 {
		t.Fatalf("len(got.Data) = %d, want 0", len(got.Data))
	}
	if got.PaginationKey != "eyJ2IjoxLCJsYXN0SWQiOiIwIiwicGFnZU9mZnNldCI6MX0=" {
		t.Errorf("PaginationKey = %q, want the documented cursor: without it the endpoint cannot be paged",
			got.PaginationKey)
	}
}

func TestGetFDActivitiesSingle(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []FDActivity{
				{ActivityID: "ACT3", AccountID: "A1", Type: "DIVIDEND", Amount: money.Must(money.NewFromString("50")), Currency: "USD", Status: "completed", CreateTime: "2026-01-03", Description: "Dividend"},
			},
			"pagination_key": "eyJ2IjoxLCJsYXN0SWQiOiIwIiwicGFnZU9mZnNldCI6MX0=",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithEndpoints(client.Endpoints{HTTP: srv.URL, BrokerHTTP: srv.URL}), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDActivities(context.Background(), "A1", "")
	if err != nil {
		t.Fatalf("GetFDActivities error = %v", err)
	}
	if len(got.Data) != 1 {
		t.Fatalf("len(got.Data) = %d, want 1", len(got.Data))
	}
	if got.Data[0].ActivityID != "ACT3" {
		t.Fatalf("ActivityID = %s, want ACT3", got.Data[0].ActivityID)
	}
	if capturedReq.URL.Path != pathFDActivities {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDActivities)
	}
	if got.PaginationKey != "eyJ2IjoxLCJsYXN0SWQiOiIwIiwicGFnZU9mZnNldCI6MX0=" {
		t.Errorf("PaginationKey = %q, want the documented cursor: without it the endpoint cannot be paged",
			got.PaginationKey)
	}
}
