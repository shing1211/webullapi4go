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

package events_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/events"
	"github.com/shing1211/webullapi4go/pkg/domain/money"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/trade"
)

func mp(s string) *money.Money {
	m := money.Must(money.NewFromString(s))
	return &m
}

// TestSandboxEvents dials Webull's sandbox event service over TLS, subscribes,
// and waits for either a SubscribeSuccess acknowledgement or a clean AuthError.
//
// It runs only when WEBULL_SANDBOX=1 and WEBULL_APP_KEY / WEBULL_APP_SECRET are
// set; otherwise it skips so the default test run stays hermetic and
// credential-free. Set WEBULL_TRADE_ACCOUNT_ID to subscribe for a specific
// account; the account list is optional. Never commit the credentials.
func TestSandboxEvents(t *testing.T) {
	if os.Getenv("WEBULL_SANDBOX") != "1" {
		t.Skip("set WEBULL_SANDBOX=1, WEBULL_APP_KEY and WEBULL_APP_SECRET to run the sandbox test")
	}
	appKey := os.Getenv("WEBULL_APP_KEY")
	appSecret := os.Getenv("WEBULL_APP_SECRET")
	if appKey == "" || appSecret == "" {
		t.Skip("set WEBULL_APP_KEY and WEBULL_APP_SECRET to run the sandbox test")
	}

	opts := []client.Option{
		client.WithAppKey(appKey),
		client.WithAppSecret(appSecret),
		client.WithRegion(client.HK),
		client.WithEnvironment(client.Sandbox),
	}
	if baseURL := os.Getenv("WEBULL_BASE_URL"); baseURL != "" {
		opts = append(opts, client.WithBaseURL(baseURL))
	}
	core, err := client.New(opts...)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = core.Close() })

	var eopts []events.Option
	if account := os.Getenv("WEBULL_TRADE_ACCOUNT_ID"); account != "" {
		eopts = append(eopts, events.WithAccounts([]string{account}))
	}
	cl, err := events.New(core, eopts...)
	if err != nil {
		t.Fatalf("events.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	connected := make(chan struct{}, 1)
	errCh := make(chan error, 4)
	cl.OnConnect(func() { trySend(connected, struct{}{}) })
	cl.OnError(func(err error) { trySend(errCh, err) })

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	runErr := make(chan error, 1)
	go func() { runErr <- cl.Run(ctx) }()

	select {
	case <-connected:
		t.Log("subscribe acknowledged with SubscribeSuccess")
	case err := <-errCh:
		if errs.Is(err, errs.CodeAuth) || errs.Is(err, errs.CodeForbidden) {
			t.Logf("sandbox rejected the subscription with a clean typed error: %v", err)
			return
		}
		t.Fatalf("event stream error = %v", err)
	case <-ctx.Done():
		t.Skip("no SubscribeSuccess or AuthError within the timeout")
	}
}

// TestSandboxOrderEvent subscribes to the sandbox order stream, places a tiny
// non-marketable AAPL limit buy far below the market through the trade package,
// and asserts that at least one typed OrderEvent is decoded. The place is
// cancelled in cleanup.
//
// It is gated by WEBULL_SANDBOX=1, WEBULL_APP_KEY, WEBULL_APP_SECRET,
// WEBULL_TRADE_ACCOUNT_ID, and WEBULL_TRADE_MUTATE=1. In practice the sandbox
// does not push a placement event for a resting order, so the test cancels the
// order in the body after a short wait to provoke the CANCEL_SUCCESS event it
// decodes; cleanup then tolerates the already-cancelled order. Credentials are
// read from the environment only and are never committed.
func TestSandboxOrderEvent(t *testing.T) {
	if os.Getenv("WEBULL_SANDBOX") != "1" {
		t.Skip("set WEBULL_SANDBOX=1 and credentials to run the sandbox test")
	}
	if os.Getenv("WEBULL_TRADE_MUTATE") != "1" {
		t.Skip("set WEBULL_TRADE_MUTATE=1 to run the mutating sandbox test")
	}
	appKey := os.Getenv("WEBULL_APP_KEY")
	appSecret := os.Getenv("WEBULL_APP_SECRET")
	accountID := os.Getenv("WEBULL_TRADE_ACCOUNT_ID")
	if appKey == "" || appSecret == "" || accountID == "" {
		t.Skip("set WEBULL_APP_KEY, WEBULL_APP_SECRET and WEBULL_TRADE_ACCOUNT_ID to run the sandbox test")
	}

	opts := []client.Option{
		client.WithAppKey(appKey),
		client.WithAppSecret(appSecret),
		client.WithRegion(client.HK),
		client.WithEnvironment(client.Sandbox),
	}
	if baseURL := os.Getenv("WEBULL_BASE_URL"); baseURL != "" {
		opts = append(opts, client.WithBaseURL(baseURL))
	}
	core, err := client.New(opts...)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = core.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	if _, err := core.EnsureToken(ctx); err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}
	trading := trade.New(core)

	cl, err := events.New(core,
		events.WithAccounts([]string{accountID}),
		events.WithSubscribeTypes(events.SubscribeOrder),
		events.WithMaxReconnectAttempts(5),
	)
	if err != nil {
		t.Fatalf("events.New() error = %v", err)
	}
	defer func() { _ = cl.Close() }()

	var (
		mu          sync.Mutex
		orderEvents []*events.OrderEvent
	)
	countEvents := func() int {
		mu.Lock()
		defer mu.Unlock()
		return len(orderEvents)
	}

	connected := make(chan struct{}, 4)
	errCh := make(chan error, 8)
	cl.OnConnect(func() { trySend(connected, struct{}{}) })
	cl.OnError(func(err error) { trySend(errCh, err) })
	cl.OnOrder(func(ev *events.OrderEvent) {
		mu.Lock()
		orderEvents = append(orderEvents, ev)
		mu.Unlock()
	})

	runCtx, runCancel := context.WithCancel(ctx)
	defer runCancel()
	go func() { _ = cl.Run(runCtx) }()

	select {
	case <-connected:
	case err := <-errCh:
		if errs.Is(err, errs.CodeAuth) || errs.Is(err, errs.CodeForbidden) {
			t.Skipf("sandbox rejected the events subscription: %v", err)
		}
		t.Fatalf("event stream error = %v", err)
	case <-time.After(15 * time.Second):
		t.Skip("no SubscribeSuccess within the timeout")
	}

	clientOrderID := fmt.Sprintf("sdk-evt-%d", time.Now().UnixNano())
	placed, err := trading.PlaceOrder(ctx, trade.PlaceOrderRequest{
		AccountID: accountID,
		NewOrders: []trade.OrderRequest{{
			ClientOrderID:         clientOrderID,
			ComboType:             trade.ComboTypeNormal,
			InstrumentType:        trade.InstrumentTypeEquity,
			Market:                trade.MarketUS,
			Symbol:                "AAPL",
			OrderType:             trade.OrderTypeLimit,
			Side:                  trade.OrderSideBuy,
			Quantity:              mp("1"),
			EntrustType:           trade.EntrustTypeQty,
			TimeInForce:           trade.TimeInForceDay,
			SupportTradingSession: trade.TradingSessionCore,
			LimitPrice:            mp("1.00"),
		}},
	})
	if err != nil {
		t.Fatalf("PlaceOrder() error = %v", err)
	}
	t.Logf("placed order_id=%s client_order_id=%s", placed.OrderID, placed.ClientOrderID)

	var cancelOnce sync.Once
	cancelOrder := func() {
		cancelOnce.Do(func() {
			cctx, ccancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer ccancel()
			if _, err := trading.CancelOrder(cctx, trade.CancelOrderRequest{
				AccountID:     accountID,
				ClientOrderID: clientOrderID,
			}); err != nil {
				t.Logf("cancel order %s failed: %v", clientOrderID, err)
				return
			}
			t.Log("cancelled order")
		})
	}
	t.Cleanup(cancelOrder)

	time.Sleep(3 * time.Second)
	if countEvents() == 0 {
		cancelOrder()
	}

	deadline := time.Now().Add(20 * time.Second)
	for countEvents() == 0 && time.Now().Before(deadline) {
		select {
		case err := <-errCh:
			t.Logf("event stream error while waiting: %v", err)
		case <-time.After(200 * time.Millisecond):
		}
	}

	mu.Lock()
	got := append([]*events.OrderEvent(nil), orderEvents...)
	mu.Unlock()
	if len(got) == 0 {
		t.Errorf("no OrderEvent received within the timeout")
		return
	}
	for _, ev := range got {
		t.Logf("order event order_id=%s status=%s scene=%s symbol=%s raw=%s",
			ev.OrderID, ev.OrderStatus, ev.SceneType, ev.Symbol, ev.Raw)
	}
}
