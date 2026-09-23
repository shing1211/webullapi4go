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
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/pkg/domain/money"
	"github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/trade"
)

func TestListAccounts(t *testing.T) {
	t.Parallel()

	const body = `[{"account_id":"ACC1","account_number":"10010048",` +
		`"account_type":"CASH","account_class":"INDIVIDUAL_CASH"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/trading/accounts/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.RawQuery, ""; got != want {
			t.Errorf("RawQuery = %q, want empty", got)
		}
		if got, want := r.Header.Get("x-version"), client.APIVersionV3; got != want {
			t.Errorf("x-version = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.ListAccounts(context.Background())
	if err != nil {
		t.Fatalf("ListAccounts() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d accounts, want 1", len(got))
	}
	acct := got[0]
	if acct.AccountID != "ACC1" || acct.AccountNumber != "10010048" {
		t.Errorf("identity = %+v, want ACC1/10010048", acct)
	}
	if acct.AccountType != trade.AccountTypeCash {
		t.Errorf("account_type = %q, want %q", acct.AccountType, trade.AccountTypeCash)
	}
	if acct.AccountClass != trade.AccountClassIndividualCash {
		t.Errorf("account_class = %q, want %q", acct.AccountClass, trade.AccountClassIndividualCash)
	}
}

func TestGetBalance(t *testing.T) {
	t.Parallel()

	const body = `{"total_asset_currency":"USD","total_cash_balance":"485705.0",` +
		`"total_market_value":"995705.0","total_unrealized_profit_loss":"227689.0",` +
		`"init_margin":"18000.0","account_currency_assets":[` +
		`{"currency":"USD","cash_balance":"485705.95","settled_cash":"485705.95",` +
		`"unsettled_cash":"0.0","market_value":"0.0","held_amount":"0.0",` +
		`"frozen_amount":"485705","buying_power":"484551",` +
		`"unrealized_profit_loss":"227689","available_withdrawal":"3.0558743194E8",` +
		`"interests_unpaid":"0.0","init_margin":"18000.0"}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/trading/assets/balances/get"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("account_id"), "ACC1"; got != want {
			t.Errorf("account_id = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("x-version"), client.APIVersionV3; got != want {
			t.Errorf("x-version = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetBalance(context.Background(), "ACC1")
	if err != nil {
		t.Fatalf("GetBalance() error = %v", err)
	}
	if got.TotalAssetCurrency != "USD" || got.TotalCashBalance.Cmp(money.Must(money.NewFromString("485705.0"))) != 0 {
		t.Errorf("totals = %+v", got)
	}
	if got.TotalUnrealizedProfitLoss.Cmp(money.Must(money.NewFromString("227689.0"))) != 0 || got.InitMargin.Cmp(money.Must(money.NewFromString("18000.0"))) != 0 {
		t.Errorf("profit/init margin = %+v", got)
	}
	if len(got.AccountCurrencyAssets) != 1 {
		t.Fatalf("got %d currency assets, want 1", len(got.AccountCurrencyAssets))
	}
	cur := got.AccountCurrencyAssets[0]
	if cur.Currency != "USD" || cur.CashBalance.Cmp(money.Must(money.NewFromString("485705.95"))) != 0 || cur.BuyingPower.Cmp(money.Must(money.NewFromString("484551"))) != 0 {
		t.Errorf("currency assets = %+v", cur)
	}
	if cur.AvailableWithdrawal.Cmp(money.Must(money.NewFromString("3.0558743194E8"))) != 0 || cur.InterestsUnpaid.Cmp(money.Must(money.NewFromString("0.0"))) != 0 {
		t.Errorf("withdrawal/interest = %+v", cur)
	}
}

func TestGetPositions(t *testing.T) {
	t.Parallel()

	const body = `[{"position_id":"POS1","currency":"USD","quantity":"1",` +
		`"symbol":"AAPL","option_strategy":"SINGLE","instrument_type":"OPTION",` +
		`"last_price":"10.0","cost_price":"11.12","unrealized_profit_loss":"0.08",` +
		`"legs":[{"symbol":"AAPL","quantity":"4","option_type":"CALL",` +
		`"option_expire_date":"2019-09-20","option_exercise_price":"11.0",` +
		`"option_contract_multiplier":"100","option_contract_deliverable":"100",` +
		`"expiration_type":"AM"}]}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/trading/assets/positions/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("account_id"), "ACC1"; got != want {
			t.Errorf("account_id = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("x-version"), client.APIVersionV3; got != want {
			t.Errorf("x-version = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetPositions(context.Background(), "ACC1")
	if err != nil {
		t.Fatalf("GetPositions() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d positions, want 1", len(got))
	}
	pos := got[0]
	if pos.PositionID != "POS1" || pos.Symbol != "AAPL" || pos.Quantity.Cmp(money.Must(money.NewFromString("1"))) != 0 {
		t.Errorf("identity = %+v", pos)
	}
	if pos.InstrumentType != trade.InstrumentTypeOption {
		t.Errorf("instrument_type = %q, want %q", pos.InstrumentType, trade.InstrumentTypeOption)
	}
	if pos.OptionStrategy != trade.OptionStrategySingle {
		t.Errorf("option_strategy = %q, want %q", pos.OptionStrategy, trade.OptionStrategySingle)
	}
	if len(pos.Legs) != 1 {
		t.Fatalf("got %d legs, want 1", len(pos.Legs))
	}
	leg := pos.Legs[0]
	if leg.OptionType != trade.OptionTypeCall || leg.OptionExercisePrice.Cmp(money.Must(money.NewFromString("11.0"))) != 0 {
		t.Errorf("leg = %+v", leg)
	}
	if leg.OptionContractMultiplier.Cmp(money.Must(money.NewFromString("100"))) != 0 || leg.ExpirationType != "AM" {
		t.Errorf("leg contract/expiration = %+v", leg)
	}
}

func TestAccountScopedEndpointsRejectEmptyAccountID(t *testing.T) {
	t.Parallel()

	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	cases := []struct {
		name string
		call func(context.Context, string) error
	}{
		{"GetBalance", func(ctx context.Context, id string) error {
			_, err := c.GetBalance(ctx, id)
			return err
		}},
		{"GetPositions", func(ctx context.Context, id string) error {
			_, err := c.GetPositions(ctx, id)
			return err
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, id := range []string{"", "   "} {
				err := tc.call(context.Background(), id)
				if !errs.Is(err, errs.CodeInvalidConfig) {
					t.Errorf("%s(%q) error = %v, want invalid_config", tc.name, id, err)
				}
			}
		})
	}
	if got := hits.Load(); got != 0 {
		t.Fatalf("server received %d requests, want 0 (validation must precede the network)", got)
	}
}

func TestGetBalancePropagatesTypedAPIError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error_code":"FORBIDDEN","message":"Insufficient permission"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetBalance(context.Background(), "ACC1")
	if !errs.Is(err, errs.CodeForbidden) {
		t.Fatalf("GetBalance() error = %v, want forbidden", err)
	}
	var typed *errs.Error
	if !errors.As(err, &typed) {
		t.Fatalf("error %v is not *errs.Error", err)
	}
	if typed.Status != http.StatusForbidden {
		t.Errorf("status = %d, want %d", typed.Status, http.StatusForbidden)
	}
}

func TestTradingRequestsCarryAccessToken(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Header.Get(client.AccessTokenHeader), "tok-trade"; got != want {
			t.Errorf("%s = %q, want %q", client.AccessTokenHeader, got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	c.Core().SetToken(&client.Token{
		Value:     "tok-trade",
		Status:    client.TokenStatusNormal,
		ExpiresAt: time.Now().Add(time.Hour),
	})
	if _, err := c.ListAccounts(context.Background()); err != nil {
		t.Fatalf("ListAccounts() error = %v", err)
	}
}
