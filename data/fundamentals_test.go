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
	"github.com/shing1211/webullapi4go/pkg/domain/money"
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
	if got.Mean.Cmp(money.Must(money.NewFromString("327.83744"))) != 0 || got.Low.Cmp(money.Must(money.NewFromString("215"))) != 0 || got.High.Cmp(money.Must(money.NewFromString("405"))) != 0 || got.Median.Cmp(money.Must(money.NewFromString("340"))) != 0 {
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

func TestGetCapitalFlow(t *testing.T) {
	t.Parallel()

	const body = `[{"date":"20260915","large_in":"1.78E8","large_out":"2.32E8",` +
		`"medium_in":"5.38E7","medium_out":"6.09E7","small_in":"5.25E7","small_out":"5.44E7"},` +
		`{"date":"20260916","large_in":"2.10E8","large_out":"1.95E8",` +
		`"medium_in":"4.90E7","medium_out":"5.20E7","small_in":"4.80E7","small_out":"4.60E7"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/fundamentals/capital-flows/get"; got != want {
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
	got, err := c.GetCapitalFlow(context.Background(), "AAPL", data.StockCategoryUS, 5)
	if err != nil {
		t.Fatalf("GetCapitalFlow() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(capitalFlow) = %d, want 2", len(got))
	}
	if got[0].Date != "20260915" || got[0].LargeIn.Cmp(money.Must(money.NewFromString("1.78E8"))) != 0 {
		t.Errorf("capitalFlow[0] = %+v", got[0])
	}
}

func TestGetIndustryComparison(t *testing.T) {
	t.Parallel()

	const body = `{"fiscal_year":2026,"fiscal_period":1,` +
		`"industry_name":"Phones & Handheld Devices","type":"EPS_TTM",` +
		`"data":[{"symbol":"AAPL","name":"Apple","rank":1,"value":"8.266"},` +
		`{"symbol":"SAMSUNG","name":"Samsung","rank":2,"value":"5.100"}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/fundamentals/industry-comparisons/get"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetIndustryComparison(context.Background(), "AAPL", data.StockCategoryUS, "EPS_TTM")
	if err != nil {
		t.Fatalf("GetIndustryComparison() error = %v", err)
	}
	if got.IndustryName != "Phones & Handheld Devices" {
		t.Errorf("industry = %q", got.IndustryName)
	}
	if len(got.Data) != 2 {
		t.Errorf("len(data) = %d, want 2", len(got.Data))
	}
	if got.Data[0].Symbol != "AAPL" || got.Data[0].Rank != 1 {
		t.Errorf("data[0] = %+v", got.Data[0])
	}
}

func TestGetEarningsCalendar(t *testing.T) {
	t.Parallel()

	const body = `[{"fiscal_year":2026,"fiscal_period":1,"currency":"USD",` +
		`"expected_publish_date":"2026-02-05","eps_actual":"2.18","eps_est":"2.10",` +
		`"rev_actual":"119580000000","rev_est":"117900000000"},` +
		`{"fiscal_year":2025,"fiscal_period":4,"currency":"USD",` +
		`"expected_publish_date":"2025-11-01","eps_actual":"","eps_est":"2.05",` +
		`"rev_actual":"","rev_est":"115000000000"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/fundamentals/earnings-calendars/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetEarningsCalendar(context.Background(), "AAPL", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetEarningsCalendar() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(earnings) = %d, want 2", len(got))
	}
	if got[0].EPSActual != "2.18" || got[0].RevActual != "119580000000" {
		t.Errorf("earnings[0] = %+v", got[0])
	}
}

func TestGetDividendCalendar(t *testing.T) {
	t.Parallel()

	const body = `[{"symbol":"AAPL","market":"US","currency":"USD","amount":"0.25",` +
		`"div_type":"CASH_DIVIDEND","declare_date":"2026-05-01",` +
		`"ex_div_date":"2026-05-08","record_date":"2026-05-11","pay_date":"2026-05-15"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/fundamentals/dividend-calendars/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetDividendCalendar(context.Background(), "AAPL", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetDividendCalendar() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(dividends) = %d, want 1", len(got))
	}
	if got[0].DivType != "CASH_DIVIDEND" || got[0].Amount.Cmp(money.Must(money.NewFromString("0.25"))) != 0 {
		t.Errorf("dividend[0] = %+v", got[0])
	}
}

func TestGetFilings(t *testing.T) {
	t.Parallel()

	const body = `{"symbol":"AAPL","category":"US_STOCK",` +
		`"filings":[{"title":"8-K | Apple Inc.","url":"https://www.sec.gov/...","publish_date":"2026-07-31"}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/fundamentals/filings/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFilings(context.Background(), "AAPL", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetFilings() error = %v", err)
	}
	if len(got.Filings) != 1 {
		t.Fatalf("len(filings) = %d, want 1", len(got.Filings))
	}
	if got.Filings[0].Title == "" {
		t.Errorf("filing title is empty")
	}
}

func TestGetIncomeStatement(t *testing.T) {
	t.Parallel()

	const body = `[{"fiscal_year":"2025","fiscal_period":"Q1","revenue":"119580000000","net_income":"23600"},` +
		`{"fiscal_year":"2024","fiscal_period":"Q4","revenue":"115600000000","net_income":"20800"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/fundamentals/income-statements/get"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetIncomeStatement(context.Background(), "AAPL", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetIncomeStatement() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(income) = %d, want 2", len(got))
	}
}

func TestGetBalanceSheet(t *testing.T) {
	t.Parallel()

	const body = `[{"fiscal_year":"2025","fiscal_period":"Q1","total_assets":"352500000000",` +
		`"total_liabilities":"267300000000"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/fundamentals/balance-sheets/get"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetBalanceSheet(context.Background(), "AAPL", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetBalanceSheet() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(balance) = %d, want 1", len(got))
	}
}

func TestGetCashFlow(t *testing.T) {
	t.Parallel()

	const body = `[{"fiscal_year":"2025","fiscal_period":"Q1",` +
		`"operating_cash_flow":"28450000000","free_cash_flow":"23600000000"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/fundamentals/cash-flows/get"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetCashFlow(context.Background(), "AAPL", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetCashFlow() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(cashFlow) = %d, want 1", len(got))
	}
}

func TestGetFinancialIndicators(t *testing.T) {
	t.Parallel()

	const body = `{"roa":"8.5","roe":"25.0","eps":"6.56","net_margin":"19.7","debt_ratio":"0.75"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/fundamentals/indicators/get"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFinancialIndicators(context.Background(), "AAPL", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetFinancialIndicators() error = %v", err)
	}
	if got.ROA != "8.5" || got.ROE != "25.0" {
		t.Errorf("indicators = %+v", got)
	}
}

func TestGetFinancialAlert(t *testing.T) {
	t.Parallel()

	const body = `{"symbol":"AAPL","category":"US_STOCK",` +
		`"expected_report_date":"2026-02-05","estimated_eps":"2.10","last_year_eps":"2.18"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/fundamentals/financial-alerts/get"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFinancialAlert(context.Background(), "AAPL", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetFinancialAlert() error = %v", err)
	}
	if got.EstimatedEPS != "2.10" || got.LastYearEPS != "2.18" {
		t.Errorf("alert = %+v", got)
	}
}

func TestGetForecastEPS(t *testing.T) {
	t.Parallel()

	const body = `[{"fiscal_year":2026,"fiscal_period":1,"actual":"2.18","est":"2.10","reported":true},` +
		`{"fiscal_year":2026,"fiscal_period":2,"actual":"","est":"2.25","reported":false}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/market-data/fundamentals/forecast-eps/get"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetForecastEPS(context.Background(), "AAPL", data.StockCategoryUS)
	if err != nil {
		t.Fatalf("GetForecastEPS() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(forecastEPS) = %d, want 2", len(got))
	}
	if got[0].FiscalYear != 2026 || got[0].Reported != true {
		t.Errorf("forecastEPS[0] = %+v", got[0])
	}
}
