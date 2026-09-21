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

func TestListAgreements(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]Agreement{
			{AgreementID: "A1", AgreementType: "margin", Version: "1.0", Status: "active", EffectiveDate: "2026-01-01"},
			{AgreementID: "A2", AgreementType: "cash", Version: "2.0", Status: "active", EffectiveDate: "2026-01-02"},
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.ListAgreements(context.Background())
	if err != nil {
		t.Fatalf("ListAgreements error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if capturedReq.URL.Path != pathAgreements {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathAgreements)
	}
}

func TestGetAgreementDetail(t *testing.T) {
	var capturedReq *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Agreement{
			AgreementID: "A1", AgreementType: "margin", Version: "1.0", Status: "active", EffectiveDate: "2026-01-01", Content: "full text",
		})
	}))
	defer srv.Close()

	cl, err := client.New(client.WithBaseURL(srv.URL), client.WithAppKey("test-key"), client.WithAppSecret("test-secret"), client.WithRegion(client.HK))
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	c := New(cl)

	got, err := c.GetAgreementDetail(context.Background(), "A1")
	if err != nil {
		t.Fatalf("GetAgreementDetail error = %v", err)
	}
	if got.AgreementID != "A1" {
		t.Fatalf("AgreementID = %s, want A1", got.AgreementID)
	}
	if capturedReq.URL.Path != pathAgreementDetail {
		t.Fatalf("path = %s, want %s", capturedReq.URL.Path, pathAgreementDetail)
	}
	if capturedReq.URL.Query().Get("agreement_id") != "A1" {
		t.Fatalf("agreement_id = %s, want A1", capturedReq.URL.Query().Get("agreement_id"))
	}
}
