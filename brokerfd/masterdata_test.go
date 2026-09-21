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

func TestGetFDEnums(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]FDEnum{
			{EnumType: "account_type", Value: "CASH", Label: "Cash Account"},
			{EnumType: "account_type", Value: "MARGIN", Label: "Margin Account"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDEnums(context.Background())
	if err != nil {
		t.Fatalf("GetFDEnums error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if capturedReq.URL.Path != pathFDEnums {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDEnums)
	}
}

func TestGetFDTradeCalendar(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]FDTradeCalendarEntry{
			{Date: "2026-01-01", Market: "US", Status: "OPEN", OpenTime: "09:30", CloseTime: "16:00"},
			{Date: "2026-01-02", Market: "US", Status: "OPEN", OpenTime: "09:30", CloseTime: "16:00"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetFDTradeCalendar(context.Background(), "US", "2026-01-01", "2026-01-31")
	if err != nil {
		t.Fatalf("GetFDTradeCalendar error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if capturedReq.URL.Path != pathFDTradeCalendar {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathFDTradeCalendar)
	}
	if capturedReq.URL.Query().Get("market") != "US" {
		t.Fatalf("market = %s, want US", capturedReq.URL.Query().Get("market"))
	}
	if capturedReq.URL.Query().Get("start_date") != "2026-01-01" {
		t.Fatalf("start_date = %s, want 2026-01-01", capturedReq.URL.Query().Get("start_date"))
	}
	if capturedReq.URL.Query().Get("end_date") != "2026-01-31" {
		t.Fatalf("end_date = %s, want 2026-01-31", capturedReq.URL.Query().Get("end_date"))
	}
}
