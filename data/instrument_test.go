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

func TestGetStockInstruments(t *testing.T) {
	t.Parallel()

	const body = `[{"name":"APPLE INC","instrument_id":"913256135","exchange_code":"NSQ",` +
		`"category":"US_STOCK","symbol":"AAPL","status":"OC","shortable":true,` +
		`"fractionable":true,"marginable":true,"overnight_trading_supported":true,` + //nolint:misspell
		`"margin_requirement_long":"0.5","margin_requirement_short":"0.5",` +
		`"intraday_margin_long":"0.5","intraday_margin_short":"0.5",` +
		`"maintenance_margin_long":"0.45","maintenance_margin_short":"0.45",` +
		`"easy_to_borrow":true,"lot_size":"1.0","currency":"USD","sub_category":"COMMON_STOCK"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/openapi/instrument/stock/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("category"), "US_STOCK"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("symbols"), "AAPL,TSLA"; got != want {
			t.Errorf("symbols = %q, want %q", got, want)
		}
		if got, want := q.Get("status"), "OC"; got != want {
			t.Errorf("status = %q, want %q", got, want)
		}
		if got, want := q.Get("sub_category"), "COMMON_STOCK"; got != want {
			t.Errorf("sub_category = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetStockInstruments(context.Background(), data.StockInstrumentQuery{
		Category:    data.StockCategoryUS,
		Symbols:     []string{"AAPL", "TSLA"},
		Status:      data.InstrumentStatusTradable,
		SubCategory: data.StockSubCategoryCommonStock,
	})
	if err != nil {
		t.Fatalf("GetStockInstruments() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d instruments, want 1", len(got))
	}
	inst := got[0]
	if inst.Symbol != "AAPL" || inst.InstrumentID != "913256135" || inst.ExchangeCode != "NSQ" {
		t.Errorf("identity = %+v, want AAPL/913256135/NSQ", inst)
	}
	if inst.Category != data.StockCategoryUS || inst.Status != data.InstrumentStatusTradable {
		t.Errorf("category/status = %q/%q, want US_STOCK/OC", inst.Category, inst.Status)
	}
	if !inst.Shortable || !inst.Fractionable || !inst.Marginable || !inst.EasyToBorrow { //nolint:misspell
		t.Errorf("boolean flags not decoded: %+v", inst)
	}
	if inst.MaintenanceMarginLong != "0.45" || inst.LotSize != "1.0" || inst.Currency != "USD" {
		t.Errorf("decoded values = %+v", inst)
	}
	if inst.SubCategory != data.StockSubCategoryCommonStock {
		t.Errorf("sub_category = %q, want COMMON_STOCK", inst.SubCategory)
	}
}

func TestGetStockInstrumentsOmitsOptionalParams(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.RawQuery, "category=HK_STOCK"; got != want {
			t.Errorf("RawQuery = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetStockInstruments(context.Background(), data.StockInstrumentQuery{
		Category: data.StockCategoryHK,
	})
	if err != nil {
		t.Fatalf("GetStockInstruments() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("got %d instruments, want 0", len(got))
	}
}
