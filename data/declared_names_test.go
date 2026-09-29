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

package data

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

// 76 declared names across 17 symbols reached no field on their response types. Each
// one decoded to the zero value with no error reported, so a caller reading a name the
// page publishes got an empty string, a zero or an empty slice.
//
// The evidence is weaker than the required-name work in v2.1.33, and these tests are
// written to reflect that. Every page here declares its properties but publishes no
// required list, so a name may be optional or conditionally sent: the field is here
// because the page publishes it, not because a server must send it. What each test
// asserts is that the published name reaches a field and holds the published value,
// not that a live response carries it.
//
// The bodies are the documented shapes rather than the SDK's own types, because a
// round-trip is green whether the type is right or wrong: that blindness is how these
// names went missing with a green suite.

// TestFundInfoReadsTheDeclaredNames covers the six descriptive names plus the
// managers array, and pins the one type decision on this type: the page documents
// is_incumbent as an integer, so it is carried as int64 and a caller reads a number
// rather than assuming a bool.
func TestFundInfoReadsTheDeclaredNames(t *testing.T) {
	var got FundInfo
	if err := json.Unmarshal([]byte(`{
		"symbol": "TQQQ",
		"name": "ProShares UltraPro QQQ",
		"currency": "USD",
		"aum": "2.9858225645E10",
		"issuer": "ProShares",
		"custodian": "JPMorgan Chase Bank, N.A.",
		"benchmark": "NASDAQ 100 TR USD",
		"investment_objective": "The investment seeks daily investment results equal to 200% of the daily performance of the Nasdaq-100 Index.",
		"launch_date": "2010-02-09",
		"managers": [
			{"name": "Michael Neches", "title": "Manager",
			 "start_date": "2013-10-01", "end_date": "2023-05-01",
			 "is_incumbent": 1, "tenure_days": 4587,
			 "tenure_years": "12.5", "tenure_return": "61.332016"}
		]
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for name, pair := range map[string][2]string{
		"Issuer":     {got.Issuer, "ProShares"},
		"Custodian":  {got.Custodian, "JPMorgan Chase Bank, N.A."},
		"Benchmark":  {got.Benchmark, "NASDAQ 100 TR USD"},
		"LaunchDate": {got.LaunchDate, "2010-02-09"},
	} {
		if pair[0] != pair[1] {
			t.Errorf("%s = %q, want %q", name, pair[0], pair[1])
		}
	}
	// The objective is a long sentence, so it is checked for its prefix rather than
	// repeated in full: the point is that the published text reaches the field at all.
	if !strings.HasPrefix(got.InvestmentObjective, "The investment seeks daily investment results") {
		t.Errorf("InvestmentObjective = %q, want the published sentence: it is the only "+
			"statement of what the fund is for", got.InvestmentObjective)
	}
	if len(got.Managers) != 1 {
		t.Fatalf("len(Managers) = %d, want 1: without it a caller cannot see who runs "+
			"the fund, which is a fund's most consequential descriptive fact", len(got.Managers))
	}
	m := got.Managers[0]
	if m.Name != "Michael Neches" || m.TenureDays != 4587 {
		t.Errorf("manager = %+v", m)
	}
	// The page says integer, so this must not decode into a bool-backed field. A bool
	// here would also have accepted "1", which the page never documents.
	if m.IsIncumbent != 1 {
		t.Errorf("IsIncumbent = %d, want 1", m.IsIncumbent)
	}
	if m.IsIncumbent == 0 {
		t.Error("IsIncumbent read as false for a documented 1")
	}
}

// TestFundNavReadsTheDocumentedNames covers the pair where the page spells a fact this
// type has already named differently. Both sets are carried, so a response populates
// whichever it sends.
func TestFundNavReadsTheDocumentedNames(t *testing.T) {
	var got FundNav
	if err := json.Unmarshal([]byte(`{
		"symbol": "TQQQ", "name": "ProShares UltraPro QQQ",
		"currency": "USD", "exchange": "NASDAQ",
		"net_value": "60.2908", "date": "2026-04-22"
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if want := money.Must(money.NewFromString("60.2908")); got.NetValue.Cmp(want) != 0 {
		t.Errorf("NetValue = %s, want %s", got.NetValue, want)
	}
	if got.Date != "2026-04-22" {
		t.Errorf("Date = %q, want 2026-04-22: the page's date is this type's nav date, "+
			"and a caller reading either name needs a value", got.Date)
	}
	// The SDK's own spellings must be untouched by a body that carries only the
	// documented ones, which is the point of carrying both. Money's zero value is
	// compared by string rather than IsZero: a never-assigned Money wraps a
	// decimal.Decimal that has not been initialised, and its IsZero does not
	// report that state.
	if got.Nav.String() == "60.2908" || got.NavDate != "" {
		t.Errorf("Nav/NavDate = %s/%q, want them still empty: a body sending only "+
			"net_value and date must not populate the SDK's own pair",
			got.Nav, got.NavDate)
	}
}

// TestFinancialAlertReadsTheDeclaredNames covers the largest remaining batch and pins
// two type decisions: fiscal_year and fiscal_period are integers on the page, and the
// four money figures are decimal strings.
func TestFinancialAlertReadsTheDeclaredNames(t *testing.T) {
	var got FinancialAlert
	if err := json.Unmarshal([]byte(`{
		"symbol": "AAPL", "category": "EARNINGS",
		"expected_report_date": "2026-04-30",
		"currency": "USD",
		"start_date": "2026-04-30", "end_date": "2026-04-30",
		"fiscal_year": 2026, "fiscal_period": 1,
		"eps_est": "1.9439", "eps_ly": "1.9439",
		"rev_est": "109614867330", "rev_ly": "109614867330"
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.FiscalYear != 2026 || got.FiscalPeriod != 1 {
		t.Errorf("FiscalYear/FiscalPeriod = %d/%d, want 2026/1: the page documents both "+
			"as integers, so a string-backed field would have dropped them",
			got.FiscalYear, got.FiscalPeriod)
	}
	if got.Currency != "USD" {
		t.Errorf("Currency = %q, want USD", got.Currency)
	}
	for name, pair := range map[string][2]string{
		"StartDate":   {got.StartDate, "2026-04-30"},
		"EndDate":     {got.EndDate, "2026-04-30"},
		"EPSEstimate": {got.EPSEstimate, "1.9439"},
		"EPSLastYear": {got.EPSLastYear, "1.9439"},
		"RevEstimate": {got.RevEstimate, "109614867330"},
		"RevLastYear": {got.RevLastYear, "109614867330"},
	} {
		if pair[0] != pair[1] {
			t.Errorf("%s = %q, want %q", name, pair[0], pair[1])
		}
	}
}

// TestScreenerStockReadsTheNamesThreePagesDeclare covers the type three separate pages
// resolve to, which is why it accumulated the most names. A caller of the 52-week page
// could not read price_52w, and a caller of the dividend page could not read ex_date.
func TestScreenerStockReadsTheNamesThreePagesDeclare(t *testing.T) {
	var got ScreenerStock
	if err := json.Unmarshal([]byte(`{
		"symbol": "TSLA", "name": "Tesla Inc", "category": "US_STOCK",
		"currency": "USD", "pe_ttm": "0.1111",
		"change_ratio_52w": "0.0111", "price_1w": "1.01", "price_52w": "1.01",
		"dividend": "0.1111", "ex_date": "2026-01-08", "yield": "0.1111",
		"id": "6391", "advanced": "110", "declined": "39", "flat": "107"
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if want := money.Must(money.NewFromString("1.01")); got.Price52W.Cmp(want) != 0 {
		t.Errorf("Price52W = %s, want %s", got.Price52W, want)
	}
	if want := money.Must(money.NewFromString("1.01")); got.Price1W.Cmp(want) != 0 {
		t.Errorf("Price1W = %s, want %s", got.Price1W, want)
	}
	if want := money.Must(money.NewFromString("0.1111")); got.Dividend.Cmp(want) != 0 {
		t.Errorf("Dividend = %s, want %s", got.Dividend, want)
	}
	if got.ExDate != "2026-01-08" {
		t.Errorf("ExDate = %q, want 2026-01-08: without it a caller cannot tell when the "+
			"dividend detached, which is the whole point of the dividend page", got.ExDate)
	}
	if got.Yield != "0.1111" || got.ChangeRatio52W != "0.0111" || got.PETTM != "0.1111" {
		t.Errorf("Yield/ChangeRatio52W/PETTM = %q/%q/%q", got.Yield, got.ChangeRatio52W, got.PETTM)
	}
	// The sector-detail page's counts are string-typed on the page even though they
	// are counts, so they are carried as string to match the wire.
	if got.Advanced != "110" || got.Declined != "39" || got.Flat != "107" || got.ID != "6391" {
		t.Errorf("Advanced/Declined/Flat/ID = %q/%q/%q/%q", got.Advanced, got.Declined, got.Flat, got.ID)
	}
}

// TestDSNewsSummaryItemReadsTheNewsNames covers the four news pages that resolve to one
// type. The page's id is an integer, and a string-backed field would have dropped a
// 17-digit value.
func TestDSNewsSummaryItemReadsTheNewsNames(t *testing.T) {
	var got DSNewsSummaryItem
	if err := json.Unmarshal([]byte(`{
		"id": 10739820929238016, "news_time": "2026-01-08 10:00:00",
		"news_url": "https://news.example/a", "source_name": "Reuters",
		"thumbnail": "https://img.example/a.jpg", "title": "Headline"
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.ID != 10739820929238016 {
		t.Errorf("ID = %d, want 10739820929238016: the page documents an integer, and a "+
			"string-backed field would have dropped a 17-digit value", got.ID)
	}
	if got.NewsTime == "" || got.NewsURL == "" || got.SourceName == "" || got.Thumbnail == "" {
		t.Errorf("news names = %q/%q/%q/%q, all four reach a field",
			got.NewsTime, got.NewsURL, got.SourceName, got.Thumbnail)
	}
	if got.Title != "Headline" {
		t.Errorf("Title = %q, want the SDK's own spelling still populated", got.Title)
	}
}

// TestDSNewsSummaryItemReadsTheSummaryNames covers the summary page, whose type
// discriminator picks which of the free-form fields is populated. Without Type a caller
// cannot know which one to read.
func TestDSNewsSummaryItemReadsTheSummaryNames(t *testing.T) {
	var got DSNewsSummaryItem
	if err := json.Unmarshal([]byte(`{
		"type": "table", "message": "",
		"headers": {"symbol": "Symbol", "last": "Last"},
		"rows": {"symbol": "AAPL", "last": "190.1"},
		"args": {"depth": 3}
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Type != "table" {
		t.Fatalf("Type = %q, want table: it selects which of Message, Headers and Rows "+
			"carries the payload, so a caller cannot read any of them without it", got.Type)
	}
	// The page declares these three as free-form objects with no property inside, so
	// they are carried untyped. That is a weaker result than a typed struct and the
	// reason is recorded rather than hidden.
	if got.Headers["symbol"] != "Symbol" {
		t.Errorf("Headers = %v, want the published table headers", got.Headers)
	}
	if got.Rows["last"] != "190.1" {
		t.Errorf("Rows = %v, want the published table row", got.Rows)
	}
	if got.Args["depth"] != float64(3) {
		t.Errorf("Args = %v, want the published arguments", got.Args)
	}
}

// TestFinancialIndicatorReadsTheDeclaredNames covers the free-form values object, whose
// keys the page does not name, so the indicator's own name is what keys it.
func TestFinancialIndicatorReadsTheDeclaredNames(t *testing.T) {
	var got FinancialIndicator
	if err := json.Unmarshal([]byte(`{
		"roa": "0.05", "roe": "0.12",
		"currency": "USD",
		"values": {"pe_ttm": "31.2", "pb": "8.1"}
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Currency != "USD" {
		t.Errorf("Currency = %q, want USD", got.Currency)
	}
	if got.Values["pe_ttm"] != "31.2" || got.Values["pb"] != "8.1" {
		t.Errorf("Values = %v, want both published keys: the page declares the object "+
			"without naming any property inside, so nothing here can name its fields", got.Values)
	}
}

// TestLogoReadsTheDeclaredNames covers the pair where the page's logo_url is this
// type's own logo.
func TestLogoReadsTheDeclaredNames(t *testing.T) {
	var got Logo
	if err := json.Unmarshal([]byte(`{
		"symbol": "AAPL", "category": "US_STOCK",
		"logo_url": "https://quotes-static.webull.com/AAPL.png"
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.LogoURL != "https://quotes-static.webull.com/AAPL.png" {
		t.Errorf("LogoURL = %q", got.LogoURL)
	}
	if got.Category != "US_STOCK" {
		t.Errorf("Category = %q, want US_STOCK", got.Category)
	}
	if got.Logo != "" {
		t.Errorf("Logo = %q, want it still empty: a body sending only logo_url must not "+
			"populate the SDK's own spelling", got.Logo)
	}
}

// TestStockInstrumentReadsTheDeclaredNames covers the v3 profile names. The page
// documents is_adr as a string spelling a boolean, so it is carried as string and a
// caller reads "true" rather than true.
func TestStockInstrumentReadsTheDeclaredNames(t *testing.T) {
	var got StockInstrument
	if err := json.Unmarshal([]byte(`{
		"symbol": "AAPL", "name": "Apple",
		"is_adr": "false", "subtype": "COMMON_STOCK"
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Subtype != "COMMON_STOCK" {
		t.Errorf("Subtype = %q, want COMMON_STOCK", got.Subtype)
	}
	if got.IsADR != "false" {
		t.Errorf("IsADR = %q, want the quoted string \"false\": the page sends a string, "+
			"so a bool-backed field would have rejected it", got.IsADR)
	}
}
