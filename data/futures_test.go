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

func TestGetFuturesInstruments(t *testing.T) {
	t.Parallel()

	const body = `[{"symbol":"ESZ5","instrument_id":"470059643","exchange_code":"XCME",` +
		`"code":"ES","name":"E-mini S&P 500 Futures Dec 2025","product_class_id":2,` +
		`"product_class_name":"Equities","status":"OC","currency":"USD",` +
		`"contract_month":"202512","settlement_date":"2025-12-29","size":"50.0",` +
		`"unit":"1-index points","min_tick":"0.25","first_notice_date":"2025-11-25",` +
		`"last_notice_date":"2025-11-28","first_trading_date":"2024-12-01",` +
		`"last_trading_date":"2025-12-19","contract_type":"MONTHLY","settlement":"Cash"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/trading/instruments/futures/contracts/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("category"), "US_FUTURES"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("code"), "ES"; got != want {
			t.Errorf("code = %q, want %q", got, want)
		}
		if got, want := q.Get("status"), "OC"; got != want {
			t.Errorf("status = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFuturesInstruments(context.Background(), data.FuturesInstrumentQuery{
		Category: data.FuturesCategoryUS,
		Code:     "ES",
		Status:   data.InstrumentStatusTradable,
	})
	if err != nil {
		t.Fatalf("GetFuturesInstruments() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d instruments, want 1", len(got))
	}
	inst := got[0]
	if inst.Symbol != "ESZ5" || inst.Code != "ES" || inst.InstrumentID != "470059643" {
		t.Errorf("identity = %+v", inst)
	}
	if inst.ProductClassID != 2 || inst.ProductClassName != "Equities" {
		t.Errorf("product class = %d/%q", inst.ProductClassID, inst.ProductClassName)
	}
	if inst.ContractType != data.FuturesContractTypeMonthly || inst.Settlement != data.FuturesSettlementCash {
		t.Errorf("contract type/settlement = %q/%q", inst.ContractType, inst.Settlement)
	}
	if inst.MinTick != "0.25" || inst.LastTradingDate != "2025-12-19" {
		t.Errorf("values = %+v", inst)
	}
}

func TestGetFuturesInstrumentsNumericUnit(t *testing.T) {
	t.Parallel()

	const body = `[{"symbol":"HSIZ5","instrument_id":"119540372","exchange_code":"HKEX","code":"HSIZ","name":"Hang Seng Index Futures Dec 2025","product_class_id":2,"product_class_name":"Index","status":"OC","currency":"HKD","contract_month":"202512","settlement_date":"2025-12-29","size":"50.0","unit":1,"min_tick":"1.0","first_notice_date":"","last_notice_date":"","first_trading_date":"2024-12-02","last_trading_date":"2025-12-23","contract_type":"MONTHLY","settlement":"Cash"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/trading/instruments/futures/contracts/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFuturesInstruments(context.Background(), data.FuturesInstrumentQuery{
		Category: data.FuturesCategoryHK,
		Code:     "HSIZ",
	})
	if err != nil {
		t.Fatalf("GetFuturesInstruments() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d instruments, want 1", len(got))
	}
	inst := got[0]
	if inst.Unit.String() != "1" {
		t.Errorf("Unit = %q, want %q", inst.Unit.String(), "1")
	}
	if inst.MinTick != "1.0" {
		t.Errorf("MinTick = %q, want %q", inst.MinTick, "1.0")
	}
}

func TestGetFuturesProductCodes(t *testing.T) {
	t.Parallel()

	const body = `[{"name":"E-Mini S&P 500","code":"ES","product_class_id":2,` +
		`"product_class_name":"Equities","exchange_code":"XCME"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/trading/instruments/futures/product-codes/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("category"), "US_FUTURES"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		if got, want := q.Get("product_class_id"), "2"; got != want {
			t.Errorf("product_class_id = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFuturesProductCodes(context.Background(), data.FuturesProductCodeQuery{
		Category:       data.FuturesCategoryUS,
		ProductClassID: 2,
	})
	if err != nil {
		t.Fatalf("GetFuturesProductCodes() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d products, want 1", len(got))
	}
	if got[0].Code != "ES" || got[0].Name != "E-Mini S&P 500" || got[0].ProductClassID != 2 {
		t.Errorf("product = %+v", got[0])
	}
}

func TestGetFuturesProductClasses(t *testing.T) {
	t.Parallel()

	const body = `[{"product_class_id":1,"product_class_name":"Energy"},` +
		`{"product_class_id":2,"product_class_name":"Equities"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/trading/instruments/futures/product-classes/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("category"), "HK_FUTURES"; got != want {
			t.Errorf("category = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetFuturesProductClasses(context.Background(), data.FuturesCategoryHK)
	if err != nil {
		t.Fatalf("GetFuturesProductClasses() error = %v", err)
	}
	if len(got) != 2 || got[0].ProductClassID != 1 || got[1].ProductClassName != "Equities" {
		t.Errorf("classes = %+v", got)
	}
}
