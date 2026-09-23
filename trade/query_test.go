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

// orderGroupJSON is a single NORMAL order group with the fields the query
// endpoints return.
const orderGroupJSON = `{"client_order_id":"coid-1","combo_type":"NORMAL","orders":[` +
	`{"client_order_id":"coid-1","order_id":"OID-1","symbol":"AAPL","side":"BUY",` +
	`"status":"SUBMITTED","order_type":"LIMIT","instrument_type":"EQUITY",` +
	`"support_trading_session":"CORE","time_in_force":"DAY","total_quantity":"1",` +
	`"filled_quantity":"0","filled_price":"0","limit_price":"180.00","stop_price":"",` +
	`"place_time_at":"2025-11-11T05:44:35.385Z","filled_time_at":"","legs":[]}]}`

func TestGetOpenOrders(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/trading/orders/open-orders/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("account_id"), "ACC1"; got != want {
			t.Errorf("account_id = %q, want %q", got, want)
		}
		if _, ok := r.URL.Query()["pagination_key"]; ok {
			t.Errorf("pagination_key = %q, want absent", r.URL.Query().Get("pagination_key"))
		}
		if got, want := r.Header.Get("x-version"), client.APIVersionV3; got != want {
			t.Errorf("x-version = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[` + orderGroupJSON + `],"pagination_key":"next-key"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetOpenOrders(context.Background(), "ACC1")
	if err != nil {
		t.Fatalf("GetOpenOrders() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d groups, want 1", len(got))
	}
	g := got[0]
	if g.ClientOrderID != "coid-1" || g.ComboType != trade.ComboTypeNormal {
		t.Errorf("group = %+v", g)
	}
	if len(g.Orders) != 1 {
		t.Fatalf("got %d orders, want 1", len(g.Orders))
	}
	o := g.Orders[0]
	if o.OrderID != "OID-1" || o.Symbol != "AAPL" || o.Status != trade.OrderStatusSubmitted {
		t.Errorf("order identity = %+v", o)
	}
	if o.OrderType != trade.OrderTypeLimit || o.InstrumentType != trade.InstrumentTypeEquity {
		t.Errorf("order enums = %+v", o)
	}
	if o.TotalQuantity.Cmp(money.Must(money.NewFromString("1"))) != 0 || o.LimitPrice.Cmp(money.Must(money.NewFromString("180.00"))) != 0 || o.TimeInForce != trade.TimeInForceDay {
		t.Errorf("order terms = %+v", o)
	}
	if o.PlaceTimeAt != "2025-11-11T05:44:35.385Z" || o.SupportTradingSession != trade.TradingSessionCore {
		t.Errorf("order times/session = %+v", o)
	}
}

func TestGetOpenOrdersPageCarriesPaginationKey(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Query().Get("pagination_key"), "resume-key"; got != want {
			t.Errorf("pagination_key = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[` + orderGroupJSON + `],"pagination_key":"next-key"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	page, err := c.GetOpenOrdersPage(context.Background(), "ACC1", "resume-key")
	if err != nil {
		t.Fatalf("GetOpenOrdersPage() error = %v", err)
	}
	if page.PaginationKey != "next-key" {
		t.Errorf("PaginationKey = %q, want %q", page.PaginationKey, "next-key")
	}
	if len(page.Orders) != 1 {
		t.Fatalf("got %d groups, want 1", len(page.Orders))
	}
}

func TestGetAllOpenOrdersPaginates(t *testing.T) {
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
			_, _ = w.Write([]byte(`{"data":[` + orderGroupJSON + `],"pagination_key":"k1"}`))
		case 2:
			if key != "k1" {
				t.Errorf("second request pagination_key = %q, want %q", key, "k1")
			}
			_, _ = w.Write([]byte(`{"data":[` + orderGroupJSON + `]}`))
		default:
			t.Errorf("unexpected request %d", n)
		}
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetAllOpenOrders(context.Background(), "ACC1")
	if err != nil {
		t.Fatalf("GetAllOpenOrders() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d groups, want 2", len(got))
	}
	if n := hits.Load(); n != 2 {
		t.Errorf("server hits = %d, want 2", n)
	}
}

func TestGetAllOpenOrdersMaxPageGuard(t *testing.T) {
	t.Parallel()

	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[],"pagination_key":"loop"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetAllOpenOrders(context.Background(), "ACC1")
	if !errs.Is(err, errs.CodeAPI) {
		t.Fatalf("GetAllOpenOrders() error = %v, want api", err)
	}
	if n := hits.Load(); n != int32(trade.MaxOrderQueryPages) {
		t.Errorf("server hits = %d, want %d", n, trade.MaxOrderQueryPages)
	}
}

func TestGetOrderHistory(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/trading/orders/historical-orders/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		q := r.URL.Query()
		if got, want := q.Get("account_id"), "ACC1"; got != want {
			t.Errorf("account_id = %q, want %q", got, want)
		}
		if got, want := q.Get("start_time"), "2025-01-05T22:59:59.012Z"; got != want {
			t.Errorf("start_time = %q, want %q", got, want)
		}
		if got, want := q.Get("end_time"), "2025-01-06T22:59:59.012Z"; got != want {
			t.Errorf("end_time = %q, want %q", got, want)
		}
		if _, ok := q["pagination_key"]; ok {
			t.Errorf("pagination_key = %q, want absent", q.Get("pagination_key"))
		}
		if got, want := r.Header.Get("x-version"), client.APIVersionV3; got != want {
			t.Errorf("x-version = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[` + orderGroupJSON + `]}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetOrderHistory(context.Background(), trade.OrderHistoryQuery{
		AccountID: "ACC1",
		StartTime: "2025-01-05T22:59:59.012Z",
		EndTime:   "2025-01-06T22:59:59.012Z",
	})
	if err != nil {
		t.Fatalf("GetOrderHistory() error = %v", err)
	}
	if len(got) != 1 || got[0].Orders[0].Symbol != "AAPL" {
		t.Fatalf("GetOrderHistory() = %+v", got)
	}
}

func TestGetAllOrderHistoryPaginates(t *testing.T) {
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
			_, _ = w.Write([]byte(`{"data":[` + orderGroupJSON + `],"pagination_key":"h1"}`))
		case 2:
			if key != "h1" {
				t.Errorf("second request pagination_key = %q, want %q", key, "h1")
			}
			_, _ = w.Write([]byte(`{"data":[` + orderGroupJSON + `]}`))
		default:
			t.Errorf("unexpected request %d", n)
		}
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetAllOrderHistory(context.Background(), trade.OrderHistoryQuery{AccountID: "ACC1"})
	if err != nil {
		t.Fatalf("GetAllOrderHistory() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d groups, want 2", len(got))
	}
	if n := hits.Load(); n != 2 {
		t.Errorf("server hits = %d, want 2", n)
	}
}

func TestOrderHistoryRejectsBadInput(t *testing.T) {
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
		q    trade.OrderHistoryQuery
	}{
		{"empty account", trade.OrderHistoryQuery{AccountID: ""}},
		{"blank account", trade.OrderHistoryQuery{AccountID: "   "}},
		{"bad start_time", trade.OrderHistoryQuery{AccountID: "ACC1", StartTime: "01/05/2025"}},
		{"bad end_time", trade.OrderHistoryQuery{AccountID: "ACC1", EndTime: "yesterday"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := c.GetOrderHistory(context.Background(), tc.q)
			if !errs.Is(err, errs.CodeInvalidConfig) {
				t.Errorf("GetOrderHistory(%+v) error = %v, want invalid_config", tc.q, err)
			}
		})
	}
	if n := hits.Load(); n != 0 {
		t.Fatalf("server received %d requests, want 0 (validation must precede the network)", n)
	}
}

func TestGetOrderDetail(t *testing.T) {
	t.Parallel()

	const body = `{"client_order_id":"coid-9","combo_type":"NORMAL","orders":[` +
		`{"client_order_id":"coid-9","order_id":"OID-9","symbol":"AAPL","side":"BUY",` +
		`"status":"FILLED","order_type":"LIMIT","instrument_type":"EQUITY",` +
		`"time_in_force":"DAY","total_quantity":"1","filled_quantity":"1",` +
		`"filled_price":"179.50","limit_price":"180.00",` +
		`"place_time_at":"2025-11-11T05:44:35.385Z",` +
		`"filled_time_at":"2025-11-11T05:44:36.000Z",` +
		`"legs":[{"symbol":"AAPL","side":"BUY","quantity":"1","option_type":"CALL",` +
		`"option_category":"AMERICAN","option_strategy":"SINGLE","strike_price":"190.0",` +
		`"option_contract_multiplier":"100","option_contract_deliverable":"100",` +
		`"option_expire_date":"2025-11-21"}],` +
		`"commission":{"actual_commission":"1.0","receivable_commission":"0.0"},` +
		`"fees":[{"type":"FINRA_CAT_REGULATORY_FEE","actual_value":"0.01",` +
		`"receivable_value":"0.0"}]}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/trading/orders/get"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("account_id"), "ACC1"; got != want {
			t.Errorf("account_id = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("client_order_id"), "coid-9"; got != want {
			t.Errorf("client_order_id = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("x-version"), client.APIVersionV3; got != want {
			t.Errorf("x-version = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	got, err := c.GetOrderDetail(context.Background(), "ACC1", "coid-9")
	if err != nil {
		t.Fatalf("GetOrderDetail() error = %v", err)
	}
	if got.ClientOrderID != "coid-9" || len(got.Orders) != 1 {
		t.Fatalf("GetOrderDetail() = %+v", got)
	}
	o := got.Orders[0]
	if o.Status != trade.OrderStatusFilled || o.FilledQuantity.Cmp(money.Must(money.NewFromString("1"))) != 0 || o.FilledPrice.Cmp(money.Must(money.NewFromString("179.50"))) != 0 {
		t.Errorf("fill state = %+v", o)
	}
	if len(o.Legs) != 1 {
		t.Fatalf("got %d legs, want 1", len(o.Legs))
	}
	leg := o.Legs[0]
	if leg.OptionType != trade.OptionTypeCall || leg.OptionStrategy != trade.OptionStrategySingle {
		t.Errorf("leg enums = %+v", leg)
	}
	if leg.StrikePrice.Cmp(money.Must(money.NewFromString("190.0"))) != 0 || leg.OptionContractMultiplier.Cmp(money.Must(money.NewFromString("100"))) != 0 || leg.OptionExpireDate != "2025-11-21" {
		t.Errorf("leg terms = %+v", leg)
	}
	if o.Commission == nil || o.Commission.ActualCommission.Cmp(money.Must(money.NewFromString("1.0"))) != 0 {
		t.Errorf("commission = %+v", o.Commission)
	}
	if len(o.Fees) != 1 || o.Fees[0].Type != "FINRA_CAT_REGULATORY_FEE" || o.Fees[0].ActualValue.Cmp(money.Must(money.NewFromString("0.01"))) != 0 {
		t.Errorf("fees = %+v", o.Fees)
	}
}

func TestOrderQueryRejectsEmptyInput(t *testing.T) {
	t.Parallel()

	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	ctx := context.Background()
	for _, id := range []string{"", "   "} {
		if _, err := c.GetOpenOrders(ctx, id); !errs.Is(err, errs.CodeInvalidConfig) {
			t.Errorf("GetOpenOrders(%q) error = %v, want invalid_config", id, err)
		}
		if _, err := c.GetOpenOrdersPage(ctx, id, ""); !errs.Is(err, errs.CodeInvalidConfig) {
			t.Errorf("GetOpenOrdersPage(%q) error = %v, want invalid_config", id, err)
		}
		if _, err := c.GetAllOpenOrders(ctx, id); !errs.Is(err, errs.CodeInvalidConfig) {
			t.Errorf("GetAllOpenOrders(%q) error = %v, want invalid_config", id, err)
		}
		if _, err := c.GetOrderHistory(ctx, trade.OrderHistoryQuery{AccountID: id}); !errs.Is(err, errs.CodeInvalidConfig) {
			t.Errorf("GetOrderHistory(%q) error = %v, want invalid_config", id, err)
		}
		if _, err := c.GetOrderDetail(ctx, id, "coid-1"); !errs.Is(err, errs.CodeInvalidConfig) {
			t.Errorf("GetOrderDetail(%q) error = %v, want invalid_config", id, err)
		}
	}
	for _, coid := range []string{"", "   "} {
		if _, err := c.GetOrderDetail(ctx, "ACC1", coid); !errs.Is(err, errs.CodeInvalidConfig) {
			t.Errorf("GetOrderDetail(ACC1, %q) error = %v, want invalid_config", coid, err)
		}
	}
	if n := hits.Load(); n != 0 {
		t.Fatalf("server received %d requests, want 0 (validation must precede the network)", n)
	}
}

func TestGetOrderHistoryPropagatesTypedAPIError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error_code":"FORBIDDEN","message":"Insufficient permission"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.GetOrderHistory(context.Background(), trade.OrderHistoryQuery{AccountID: "ACC1"})
	if !errs.Is(err, errs.CodeForbidden) {
		t.Fatalf("GetOrderHistory() error = %v, want forbidden", err)
	}
	var typed *errs.Error
	if !errors.As(err, &typed) {
		t.Fatalf("error %v is not *errs.Error", err)
	}
	if typed.Status != http.StatusForbidden {
		t.Errorf("status = %d, want %d", typed.Status, http.StatusForbidden)
	}
}

// TestSandboxOrderQueriesReadOnly reads the working orders and the last seven
// days of history for the dedicated trading sandbox account. It is read-only:
// open orders and history never place or mutate an order, so an empty list is a
// valid result. It shares newSandboxTradeClient's gating on
// WEBULL_TRADE_SANDBOX=1 plus WEBULL_TRADE_APP_KEY, WEBULL_TRADE_APP_SECRET,
// and WEBULL_TRADE_ACCOUNT_ID, and skips otherwise so the default run stays
// hermetic.
func TestSandboxOrderQueriesReadOnly(t *testing.T) {
	trading, accountID := newSandboxTradeClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	open, err := trading.GetOpenOrders(ctx, accountID)
	if err != nil {
		t.Fatalf("GetOpenOrders(%s) error = %v", accountID, err)
	}
	t.Logf("open orders groups=%d", len(open))
	for _, g := range open {
		t.Logf("open group client_order_id=%s combo_type=%s orders=%d",
			g.ClientOrderID, g.ComboType, len(g.Orders))
	}

	history, err := trading.GetOrderHistory(ctx, trade.OrderHistoryQuery{AccountID: accountID})
	if err != nil {
		t.Fatalf("GetOrderHistory(%s) error = %v", accountID, err)
	}
	t.Logf("history groups=%d", len(history))
	for _, g := range history {
		for _, o := range g.Orders {
			t.Logf("history order_id=%s symbol=%s status=%s", o.OrderID, o.Symbol, o.Status)
		}
	}
}
