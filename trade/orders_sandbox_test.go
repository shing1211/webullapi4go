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
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/internal/errs"
	"github.com/shing1211/webullapi4go/trade"
)

// newSandboxTradeClient builds a trade client against the dedicated trading
// sandbox account, or skips the test when the environment is not configured. It
// reuses the token across the test and is gated by WEBULL_TRADE_SANDBOX plus
// WEBULL_TRADE_APP_KEY, WEBULL_TRADE_APP_SECRET, and WEBULL_TRADE_ACCOUNT_ID.
// Credentials are read from the environment only and are never committed.
func newSandboxTradeClient(t *testing.T) (*trade.Client, string) {
	t.Helper()
	if os.Getenv("WEBULL_TRADE_SANDBOX") != "1" {
		t.Skip("set WEBULL_TRADE_SANDBOX=1, WEBULL_TRADE_APP_KEY, WEBULL_TRADE_APP_SECRET and WEBULL_TRADE_ACCOUNT_ID to run the sandbox test")
	}
	appKey := os.Getenv("WEBULL_TRADE_APP_KEY")
	appSecret := os.Getenv("WEBULL_TRADE_APP_SECRET")
	accountID := os.Getenv("WEBULL_TRADE_ACCOUNT_ID")
	if appKey == "" || appSecret == "" || accountID == "" {
		t.Skip("set WEBULL_TRADE_APP_KEY, WEBULL_TRADE_APP_SECRET and WEBULL_TRADE_ACCOUNT_ID to run the sandbox test")
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
	cl, err := client.New(opts...)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	t.Cleanup(cancel)
	if _, err := cl.EnsureToken(ctx); err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}
	return trade.New(cl), accountID
}

// sandboxPreviewOrder is a small, non-marketable AAPL limit buy used to exercise
// the preview endpoint without any risk of execution.
func sandboxPreviewOrder(accountID string) trade.PlaceOrderRequest {
	return trade.PlaceOrderRequest{
		AccountID: accountID,
		NewOrders: []trade.OrderRequest{
			{
				ClientOrderID:         "sdk-preview-aapl-1",
				ComboType:             trade.ComboTypeNormal,
				InstrumentType:        trade.InstrumentTypeEquity,
				Market:                trade.MarketUS,
				Symbol:                "AAPL",
				OrderType:             trade.OrderTypeLimit,
				Side:                  trade.OrderSideBuy,
				Quantity:              "1",
				EntrustType:           trade.EntrustTypeQty,
				TimeInForce:           trade.TimeInForceDay,
				SupportTradingSession: trade.TradingSessionCore,
				LimitPrice:            "180.00",
			},
		},
	}
}

// TestSandboxPreviewOrder estimates the cost of a small AAPL limit buy on the
// sandbox account. It is read-only: preview never places or mutates an order.
func TestSandboxPreviewOrder(t *testing.T) {
	trading, accountID := newSandboxTradeClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := trading.PreviewOrder(ctx, sandboxPreviewOrder(accountID))
	if err != nil {
		t.Fatalf("PreviewOrder() error = %v", err)
	}
	t.Logf("preview estimated_cost=%s estimated_transaction_fee=%s",
		res.EstimatedCost, res.EstimatedTransactionFee)
}

// TestSandboxPlaceOrder places a tiny, non-marketable AAPL limit buy far below
// the market, captures the identifiers, and cancels it in t.Cleanup. It is
// gated by WEBULL_TRADE_SANDBOX plus WEBULL_TRADE_MUTATE so that an ordinary
// offline run never trades. A market order is never used, and the cancel is
// best-effort so a cleanup failure is logged rather than hidden.
func TestSandboxPlaceOrder(t *testing.T) {
	if os.Getenv("WEBULL_TRADE_SANDBOX") != "1" {
		t.Skip("set WEBULL_TRADE_SANDBOX=1 and credentials to run the sandbox test")
	}
	if os.Getenv("WEBULL_TRADE_MUTATE") != "1" {
		t.Skip("set WEBULL_TRADE_MUTATE=1 to run the mutating sandbox test")
	}
	trading, accountID := newSandboxTradeClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := sandboxPreviewOrder(accountID)
	req.NewOrders[0].ClientOrderID = fmt.Sprintf("sdk-place-%d", time.Now().UnixNano())
	req.NewOrders[0].LimitPrice = "1.00"

	res, err := trading.PlaceOrder(ctx, req)
	if err != nil {
		t.Fatalf("PlaceOrder() error = %v", err)
	}
	t.Logf("placed client_order_id=%s order_id=%s", res.ClientOrderID, res.OrderID)

	t.Cleanup(func() {
		cctx, ccancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer ccancel()
		cres, err := trading.CancelOrder(cctx, trade.CancelOrderRequest{
			AccountID:     accountID,
			ClientOrderID: res.ClientOrderID,
		})
		if err != nil {
			t.Logf("cleanup CancelOrder(%s) failed (order_id=%s): %v", res.ClientOrderID, res.OrderID, err)
			return
		}
		t.Logf("cleanup cancelled client_order_id=%s order_id=%s", cres.ClientOrderID, cres.OrderID)
	})
}

// TestSandboxReplaceCancelNonexistent exercises replace and cancel against a
// client order identifier that does not exist. It is read-only with respect to
// real orders: the API is expected to reject the request as a business error,
// which must surface as a typed [errs.Error] rather than a panic. A nil result
// is logged and tolerated because it indicates the sandbox accepted the no-op.
func TestSandboxReplaceCancelNonexistent(t *testing.T) {
	trading, accountID := newSandboxTradeClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	id := fmt.Sprintf("sdk-noop-%d", time.Now().UnixNano())

	res, err := trading.CancelOrder(ctx, trade.CancelOrderRequest{
		AccountID:     accountID,
		ClientOrderID: id,
	})
	if err == nil {
		t.Logf("CancelOrder(nonexistent) unexpectedly succeeded: %+v", res)
	} else {
		var e *errs.Error
		if !errors.As(err, &e) {
			t.Fatalf("CancelOrder(nonexistent) error = %v (%T), want *errs.Error", err, err)
		}
		t.Logf("CancelOrder(nonexistent) typed error: code=%s message=%s", e.Code, e.Message)
	}

	rres, rerr := trading.ReplaceOrder(ctx, trade.ReplaceOrderRequest{
		AccountID: accountID,
		ModifyOrders: []trade.ModifyOrderRequest{
			{ClientOrderID: id, LimitPrice: "1.00"},
		},
	})
	if rerr == nil {
		t.Logf("ReplaceOrder(nonexistent) unexpectedly succeeded: %+v", rres)
		return
	}
	var re *errs.Error
	if !errors.As(rerr, &re) {
		t.Fatalf("ReplaceOrder(nonexistent) error = %v (%T), want *errs.Error", rerr, rerr)
	}
	t.Logf("ReplaceOrder(nonexistent) typed error: code=%s message=%s", re.Code, re.Message)
}
