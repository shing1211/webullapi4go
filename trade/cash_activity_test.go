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

package trade_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/pkg/domain/money"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/trade"
)

const cashActivityJSON = `{"id":"ACT1","account_id":"ACC1","account_number":"10010048",` +
	`"activity_type":"TRADE","activity_sub_type":"BUY","currency":"USD",` +
	`"market":"US","symbol":"AAPL","trade_date":"2025-11-11","net_amount":"-17950.00",` +
	`"biz_time":"2025-11-11T05:44:35.385Z"}`

func TestGetCashActivities(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/trading/activities/cash-activities/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("account_id"), "ACC1"; got != want {
			t.Errorf("account_id = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("activity_type"), "TRADE"; got != want {
			t.Errorf("activity_type = %q, want %q", got, want)
		}
		if _, ok := r.URL.Query()["pagination_key"]; ok {
			t.Errorf("pagination_key = %q, want absent", r.URL.Query().Get("pagination_key"))
		}
		if got, want := r.Header.Get("x-version"), client.APIVersionV3; got != want {
			t.Errorf("x-version = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[` + cashActivityJSON + `],"pagination_key":"next-key"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetCashActivities(context.Background(), trade.CashActivityQuery{
		AccountID:    "ACC1",
		ActivityType: trade.CashActivityTypeTrade,
	})
	if err != nil {
		t.Fatalf("GetCashActivities() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d activities, want 1", len(got))
	}
	a := got[0]
	if a.ID != "ACT1" || a.AccountID != "ACC1" || a.AccountNumber != "10010048" {
		t.Errorf("identity = %+v", a)
	}
	if a.ActivityType != trade.CashActivityTypeTrade || a.ActivitySubType != "BUY" {
		t.Errorf("activity type = %+v", a)
	}
	if a.Currency != "USD" || a.Market != "US" || a.Symbol != "AAPL" {
		t.Errorf("market data = %+v", a)
	}
	if a.TradeDate != "2025-11-11" || a.NetAmount.Cmp(money.Must(money.NewFromString("-17950"))) != 0 {
		t.Errorf("trade data = %+v", a)
	}
}

func TestGetCashActivitiesPage(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Query().Get("account_id"), "ACC1"; got != want {
			t.Errorf("account_id = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("pagination_key"), "resume-key"; got != want {
			t.Errorf("pagination_key = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[` + cashActivityJSON + `],"pagination_key":"next-key"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	page, err := c.GetCashActivitiesPage(context.Background(), trade.CashActivityQuery{
		AccountID:     "ACC1",
		ActivityType:  trade.CashActivityTypeDividend,
		PaginationKey: "resume-key",
	})
	if err != nil {
		t.Fatalf("GetCashActivitiesPage() error = %v", err)
	}
	if page.PaginationKey != "next-key" {
		t.Errorf("PaginationKey = %q, want %q", page.PaginationKey, "next-key")
	}
	if len(page.Activities) != 1 {
		t.Fatalf("got %d activities, want 1", len(page.Activities))
	}
	if page.Activities[0].ID != "ACT1" {
		t.Errorf("activity ID = %q, want %q", page.Activities[0].ID, "ACT1")
	}
}

func TestGetCashActivitiesEmptyAccount(t *testing.T) {
	t.Parallel()

	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	for _, id := range []string{"", "   "} {
		_, err := c.GetCashActivities(context.Background(), trade.CashActivityQuery{AccountID: id})
		if !errs.Is(err, errs.CodeInvalidConfig) {
			t.Errorf("GetCashActivities(%q) error = %v, want invalid_config", id, err)
		}
		_, err = c.GetCashActivitiesPage(context.Background(), trade.CashActivityQuery{AccountID: id})
		if !errs.Is(err, errs.CodeInvalidConfig) {
			t.Errorf("GetCashActivitiesPage(%q) error = %v, want invalid_config", id, err)
		}
		_, err = c.GetAllCashActivities(context.Background(), trade.CashActivityQuery{AccountID: id})
		if !errs.Is(err, errs.CodeInvalidConfig) {
			t.Errorf("GetAllCashActivities(%q) error = %v, want invalid_config", id, err)
		}
	}
	if n := hits.Load(); n != 0 {
		t.Fatalf("server received %d requests, want 0 (validation must precede the network)", n)
	}
}

func TestGetAllCashActivities(t *testing.T) {
	t.Parallel()

	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := hits.Add(1)
		key := r.URL.Query().Get("pagination_key")
		w.Header().Set("Content-Type", "application/json")
		switch n {
		case 1:
			if key != "" {
				t.Errorf("first request pagination_key = %q, want empty", key)
			}
			_, _ = w.Write([]byte(`{"data":[` + cashActivityJSON + `],"pagination_key":"k1"}`))
		case 2:
			if key != "k1" {
				t.Errorf("second request pagination_key = %q, want %q", key, "k1")
			}
			_, _ = w.Write([]byte(`{"data":[` + cashActivityJSON + `]}`))
		default:
			t.Errorf("unexpected request %d", n)
		}
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetAllCashActivities(context.Background(), trade.CashActivityQuery{
		AccountID:    "ACC1",
		ActivityType: trade.CashActivityTypeInterest,
	})
	if err != nil {
		t.Fatalf("GetAllCashActivities() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d activities, want 2", len(got))
	}
	if n := hits.Load(); n != 2 {
		t.Errorf("server hits = %d, want 2", n)
	}
}

func TestGetCashActivitiesPropagatesTypedAPIError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error_code":"FORBIDDEN","message":"Insufficient permission"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetCashActivities(context.Background(), trade.CashActivityQuery{AccountID: "ACC1"})
	if !errs.Is(err, errs.CodeForbidden) {
		t.Fatalf("GetCashActivities() error = %v, want forbidden", err)
	}
	var typed *errs.Error
	if !errors.As(err, &typed) {
		t.Fatalf("error %v is not *errs.Error", err)
	}
	if typed.Status != http.StatusForbidden {
		t.Errorf("status = %d, want %d", typed.Status, http.StatusForbidden)
	}
}
