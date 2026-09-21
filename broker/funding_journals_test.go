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

package broker_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shing1211/webullapi4go/broker"
)

func TestGetFXRate(t *testing.T) {
	t.Parallel()

	const body = `{"from_currency":"USD","to_currency":"HKD","rate":"7.85","timestamp":"2026-01-01T00:00:00Z"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/openapi/v1/broker/funding/fx/rate"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("from"), "USD"; got != want {
			t.Errorf("from = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("to"), "HKD"; got != want {
			t.Errorf("to = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	got, err := c.GetFXRate(context.Background(), "USD", "HKD")
	if err != nil {
		t.Fatalf("GetFXRate() error = %v", err)
	}
	if got.FromCurrency != "USD" || got.ToCurrency != "HKD" {
		t.Errorf("currencies = %+v, want USD/HKD", got)
	}
	if got.Rate != "7.85" {
		t.Errorf("rate = %q, want 7.85", got.Rate)
	}
}

func TestCreateFXExchange(t *testing.T) {
	t.Parallel()

	const respBody = `{"exchange_id":"FX1","account_id":"ACC1","from_currency":"USD","to_currency":"HKD","from_amount":"1000.00","to_amount":"7850.00","rate":"7.85","status":"COMPLETED","create_time":"2026-01-01T00:00:00Z"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/openapi/v1/broker/funding/fx/exchange"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var req broker.CreateFXExchangeRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if req.AccountID != "ACC1" {
			t.Errorf("account_id = %q, want ACC1", req.AccountID)
		}
		if req.FromCurrency != "USD" {
			t.Errorf("from_currency = %q, want USD", req.FromCurrency)
		}
		if req.ToCurrency != "HKD" {
			t.Errorf("to_currency = %q, want HKD", req.ToCurrency)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respBody))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	got, err := c.CreateFXExchange(context.Background(), broker.CreateFXExchangeRequest{
		AccountID:    "ACC1",
		FromCurrency: "USD",
		ToCurrency:   "HKD",
		Amount:       "1000.00",
	})
	if err != nil {
		t.Fatalf("CreateFXExchange() error = %v", err)
	}
	if got.ExchangeID != "FX1" {
		t.Errorf("exchange_id = %q, want FX1", got.ExchangeID)
	}
	if got.Status != "COMPLETED" {
		t.Errorf("status = %q, want COMPLETED", got.Status)
	}
}

func TestGetFXExchangeDetail(t *testing.T) {
	t.Parallel()

	const body = `{"exchange_id":"FX1","account_id":"ACC1","from_currency":"USD","to_currency":"HKD","from_amount":"1000.00","to_amount":"7850.00","rate":"7.85","status":"COMPLETED","create_time":"2026-01-01T00:00:00Z"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/openapi/v1/broker/funding/fx/exchange/detail"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("exchange_id"), "FX1"; got != want {
			t.Errorf("exchange_id = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	got, err := c.GetFXExchangeDetail(context.Background(), "FX1")
	if err != nil {
		t.Fatalf("GetFXExchangeDetail() error = %v", err)
	}
	if got.ExchangeID != "FX1" {
		t.Errorf("exchange_id = %q, want FX1", got.ExchangeID)
	}
}

func TestCreateInstantExchange(t *testing.T) {
	t.Parallel()

	const respBody = `{"exchange_id":"IE1","account_id":"ACC1","from_currency":"USD","to_currency":"HKD","amount":"1000.00","rate":"7.85","status":"COMPLETED","create_time":"2026-01-01T00:00:00Z"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/openapi/v1/broker/funding/fx/instant-exchange"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var req broker.CreateInstantExchangeRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if req.AccountID != "ACC1" {
			t.Errorf("account_id = %q, want ACC1", req.AccountID)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respBody))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	got, err := c.CreateInstantExchange(context.Background(), broker.CreateInstantExchangeRequest{
		AccountID:    "ACC1",
		FromCurrency: "USD",
		ToCurrency:   "HKD",
		Amount:       "1000.00",
	})
	if err != nil {
		t.Fatalf("CreateInstantExchange() error = %v", err)
	}
	if got.ExchangeID != "IE1" {
		t.Errorf("exchange_id = %q, want IE1", got.ExchangeID)
	}
}

func TestGetInstantExchangeDetail(t *testing.T) {
	t.Parallel()

	const body = `{"exchange_id":"IE1","account_id":"ACC1","from_currency":"USD","to_currency":"HKD","amount":"1000.00","rate":"7.85","status":"COMPLETED","create_time":"2026-01-01T00:00:00Z"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/openapi/v1/broker/funding/fx/exchange/detail"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("exchange_id"), "IE1"; got != want {
			t.Errorf("exchange_id = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	got, err := c.GetInstantExchangeDetail(context.Background(), "IE1")
	if err != nil {
		t.Fatalf("GetInstantExchangeDetail() error = %v", err)
	}
	if got.ExchangeID != "IE1" {
		t.Errorf("exchange_id = %q, want IE1", got.ExchangeID)
	}
}

func TestCreateInstantFunding(t *testing.T) {
	t.Parallel()

	const respBody = `{"funding_id":"IF1","account_id":"ACC1","amount":"5000.00","currency":"USD","status":"COMPLETED","create_time":"2026-01-01T00:00:00Z"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/openapi/v1/broker/funding/instant"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var req broker.CreateInstantFundingRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if req.AccountID != "ACC1" {
			t.Errorf("account_id = %q, want ACC1", req.AccountID)
		}
		if req.Amount != "5000.00" {
			t.Errorf("amount = %q, want 5000.00", req.Amount)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respBody))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	got, err := c.CreateInstantFunding(context.Background(), broker.CreateInstantFundingRequest{
		AccountID: "ACC1",
		Amount:    "5000.00",
		Currency:  "USD",
	})
	if err != nil {
		t.Fatalf("CreateInstantFunding() error = %v", err)
	}
	if got.FundingID != "IF1" {
		t.Errorf("funding_id = %q, want IF1", got.FundingID)
	}
}

func TestGetInstantFundingDetail(t *testing.T) {
	t.Parallel()

	const body = `{"funding_id":"IF1","account_id":"ACC1","amount":"5000.00","currency":"USD","status":"COMPLETED","create_time":"2026-01-01T00:00:00Z"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/openapi/v1/broker/funding/instant/detail"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("funding_id"), "IF1"; got != want {
			t.Errorf("funding_id = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	got, err := c.GetInstantFundingDetail(context.Background(), "IF1")
	if err != nil {
		t.Fatalf("GetInstantFundingDetail() error = %v", err)
	}
	if got.FundingID != "IF1" {
		t.Errorf("funding_id = %q, want IF1", got.FundingID)
	}
}

func TestCreateCashJournal(t *testing.T) {
	t.Parallel()

	const respBody = `{"journal_id":"CJ1","account_id":"ACC1","type":"DEPOSIT","amount":"10000.00","currency":"HKD","status":"COMPLETED","create_time":"2026-01-01T00:00:00Z"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/openapi/v1/broker/journals/cash"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var req broker.CreateCashJournalRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if req.AccountID != "ACC1" {
			t.Errorf("account_id = %q, want ACC1", req.AccountID)
		}
		if req.Type != "DEPOSIT" {
			t.Errorf("type = %q, want DEPOSIT", req.Type)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respBody))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	got, err := c.CreateCashJournal(context.Background(), broker.CreateCashJournalRequest{
		AccountID: "ACC1",
		Type:      "DEPOSIT",
		Amount:    "10000.00",
		Currency:  "HKD",
	})
	if err != nil {
		t.Fatalf("CreateCashJournal() error = %v", err)
	}
	if got.JournalID != "CJ1" {
		t.Errorf("journal_id = %q, want CJ1", got.JournalID)
	}
}

func TestGetCashJournalDetail(t *testing.T) {
	t.Parallel()

	const body = `{"journal_id":"CJ1","account_id":"ACC1","type":"DEPOSIT","amount":"10000.00","currency":"HKD","status":"COMPLETED","create_time":"2026-01-01T00:00:00Z"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/openapi/v1/broker/journals/cash/detail"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("journal_id"), "CJ1"; got != want {
			t.Errorf("journal_id = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	got, err := c.GetCashJournalDetail(context.Background(), "CJ1")
	if err != nil {
		t.Fatalf("GetCashJournalDetail() error = %v", err)
	}
	if got.JournalID != "CJ1" {
		t.Errorf("journal_id = %q, want CJ1", got.JournalID)
	}
}

func TestCreatePositionJournal(t *testing.T) {
	t.Parallel()

	const respBody = `{"journal_id":"PJ1","account_id":"ACC1","symbol":"00700.HK","quantity":"100","action":"TRANSFER_IN","status":"COMPLETED","create_time":"2026-01-01T00:00:00Z"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/openapi/v1/broker/journals/position"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var req broker.CreatePositionJournalRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if req.AccountID != "ACC1" {
			t.Errorf("account_id = %q, want ACC1", req.AccountID)
		}
		if req.Symbol != "00700.HK" {
			t.Errorf("symbol = %q, want 00700.HK", req.Symbol)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respBody))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	got, err := c.CreatePositionJournal(context.Background(), broker.CreatePositionJournalRequest{
		AccountID: "ACC1",
		Symbol:    "00700.HK",
		Quantity:  "100",
		Action:    "TRANSFER_IN",
	})
	if err != nil {
		t.Fatalf("CreatePositionJournal() error = %v", err)
	}
	if got.JournalID != "PJ1" {
		t.Errorf("journal_id = %q, want PJ1", got.JournalID)
	}
}

func TestGetPositionJournalDetail(t *testing.T) {
	t.Parallel()

	const body = `{"journal_id":"PJ1","account_id":"ACC1","symbol":"00700.HK","quantity":"100","action":"TRANSFER_IN","status":"COMPLETED","create_time":"2026-01-01T00:00:00Z"}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/openapi/v1/broker/journals/position/detail"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("journal_id"), "PJ1"; got != want {
			t.Errorf("journal_id = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newBrokerClient(t, srv.URL)
	got, err := c.GetPositionJournalDetail(context.Background(), "PJ1")
	if err != nil {
		t.Fatalf("GetPositionJournalDetail() error = %v", err)
	}
	if got.JournalID != "PJ1" {
		t.Errorf("journal_id = %q, want PJ1", got.JournalID)
	}
}
