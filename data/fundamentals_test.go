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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shing1211/webullapi4go/data"
)

func TestGetCompanyProfile(t *testing.T) {
	t.Parallel()

	const body = `{"symbol":"AAPL","category":"US_STOCK","company_name":"Apple Inc",` +
		`"establish_date":"1977-01-03","exhibition_code":"NASDAQ",` +
		`"profile":"Apple Inc. designs and markets consumer electronics.",` +
		`"employees":"164000","address":"One Apple Park Way, Cupertino, CA",` +
		`"ceo":"Tim Cook","industries":["Consumer Electronics","Technology"]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/fundamentals/company-profiles/get"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("symbol"), "AAPL"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		if got, want := q.Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetCompanyProfile(context.Background(), "AAPL", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetCompanyProfile() error = %v", err)
	}
	if got.Symbol != "AAPL" || got.CompanyName != "Apple Inc" || got.CEO != "Tim Cook" {
		t.Errorf("profile = %+v", got)
	}
	if got.ExhibitionCode != "NASDAQ" || got.Employees != "164000" {
		t.Errorf("listing/employees = %q/%q", got.ExhibitionCode, got.Employees)
	}
	if len(got.Industries) != 2 || got.Industries[0] != "Consumer Electronics" {
		t.Errorf("industries = %v", got.Industries)
	}
}

func TestGetAnalystTargetPrice(t *testing.T) {
	t.Parallel()

	const body = `{"symbol":"AAPL","category":"US_STOCK","mean":"327.83744","low":"215",` +
		`"high":"405","median":"340","currency":"USD",` +
		`"effective_start_date":"2026-09-16T07:38:47.000+0000"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/fundamentals/analysis/target-prices/get"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("symbol"), "AAPL"; got != want {
			t.Errorf("symbol = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetAnalystTargetPrice(context.Background(), "AAPL", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetAnalystTargetPrice() error = %v", err)
	}
	if got.Mean != "327.83744" || got.Low != "215" || got.High != "405" || got.Median != "340" {
		t.Errorf("target price = %+v", got)
	}
	if got.Currency != "USD" {
		t.Errorf("currency = %q, want USD", got.Currency)
	}
}

func TestGetAnalystRating(t *testing.T) {
	t.Parallel()

	const body = `{"symbol":"AAPL","category":"US_STOCK","number":"44","under_perform":"3",` +
		`"buy":"6","sell":"3","strong_buy":"19","hold":"13",` +
		`"effective_start_date":"2026-09-16T22:12:31.000+0000"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/fundamentals/analysis/ratings/get"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetAnalystRating(context.Background(), "AAPL", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetAnalystRating() error = %v", err)
	}
	if got.Number != "44" || got.StrongBuy != "19" || got.Buy != "6" || got.Hold != "13" {
		t.Errorf("rating = %+v", got)
	}
	if got.UnderPerform != "3" || got.Sell != "3" {
		t.Errorf("negative rating counts = %+v", got)
	}
}
