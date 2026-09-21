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

package data_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/data"
	"github.com/shing1211/webullapi4go/display"
)

func newDSInstrumentTestServer(t *testing.T) (*httptest.Server, *display.Service) {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/client-tokens/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":       "test-access-token",
			"expires_at":         0,
			"refresh_token":      "refresh",
			"refresh_expires_at": 0,
		})
	})
	mux.HandleFunc("/openapi/market-data/stock/company-profile", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"symbol":       "AAPL",
			"category":     "US",
			"company_name": "Apple Inc.",
			"profile":      "Technology company",
			"employees":    "150000",
			"address":      "Cupertino, CA",
			"ceo":          "Tim Cook",
		})
	})
	mux.HandleFunc("/openapi/market-data/stock/analyst-target-price", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"symbol":   "AAPL",
			"category": "US",
			"mean":     "200.00",
			"low":      "180.00",
			"high":     "220.00",
			"median":   "195.00",
			"currency": "USD",
		})
	})
	mux.HandleFunc("/openapi/market-data/stock/analyst-rating", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"symbol":   "AAPL",
			"category": "US",
			"number":   "30",
			"buy":      "20",
			"sell":     "2",
			"hold":     "8",
		})
	})

	srv := httptest.NewServer(mux)
	dsvc := display.NewService(testAppKey, testAppSecret,
		display.WithBaseURL(srv.URL),
		display.WithHTTPClient(srv.Client()),
	)
	return srv, dsvc
}

func newDSInstrumentDataClient(t *testing.T, dsvc *display.Service) *data.Client {
	t.Helper()
	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })
	dc := data.New(cl)
	dc.SetDisplaySvcForTesting(dsvc)
	return dc
}

func TestGetDSCompanyProfile(t *testing.T) {
	t.Parallel()
	srv, dsvc := newDSInstrumentTestServer(t)
	defer srv.Close()
	dc := newDSInstrumentDataClient(t, dsvc)

	profile, err := dc.GetDSCompanyProfile(context.Background(), "AAPL")
	if err != nil {
		t.Fatalf("GetDSCompanyProfile() error = %v", err)
	}
	if profile.Symbol != "AAPL" {
		t.Errorf("Symbol = %q, want %q", profile.Symbol, "AAPL")
	}
	if profile.CompanyName != "Apple Inc." {
		t.Errorf("CompanyName = %q, want %q", profile.CompanyName, "Apple Inc.")
	}
	if profile.CEO != "Tim Cook" {
		t.Errorf("CEO = %q, want %q", profile.CEO, "Tim Cook")
	}
}

func TestGetDSAnalystTargetPrice(t *testing.T) {
	t.Parallel()
	srv, dsvc := newDSInstrumentTestServer(t)
	defer srv.Close()
	dc := newDSInstrumentDataClient(t, dsvc)

	tp, err := dc.GetDSAnalystTargetPrice(context.Background(), "AAPL")
	if err != nil {
		t.Fatalf("GetDSAnalystTargetPrice() error = %v", err)
	}
	if tp.Symbol != "AAPL" {
		t.Errorf("Symbol = %q, want %q", tp.Symbol, "AAPL")
	}
	if tp.Mean != "200.00" {
		t.Errorf("Mean = %q, want %q", tp.Mean, "200.00")
	}
	if tp.Currency != "USD" {
		t.Errorf("Currency = %q, want %q", tp.Currency, "USD")
	}
}

func TestGetDSAnalystRating(t *testing.T) {
	t.Parallel()
	srv, dsvc := newDSInstrumentTestServer(t)
	defer srv.Close()
	dc := newDSInstrumentDataClient(t, dsvc)

	rating, err := dc.GetDSAnalystRating(context.Background(), "AAPL")
	if err != nil {
		t.Fatalf("GetDSAnalystRating() error = %v", err)
	}
	if rating.Symbol != "AAPL" {
		t.Errorf("Symbol = %q, want %q", rating.Symbol, "AAPL")
	}
	if rating.Buy != "20" {
		t.Errorf("Buy = %q, want %q", rating.Buy, "20")
	}
	if rating.Sell != "2" {
		t.Errorf("Sell = %q, want %q", rating.Sell, "2")
	}
}
