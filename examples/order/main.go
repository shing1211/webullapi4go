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

// Command order demonstrates the stock-order lifecycle: it previews a small
// AAPL limit buy and, only when explicitly opted in, places it far below the
// market and immediately cancels it. A market order is never submitted.
//
// Previewing is read-only. Placing an order mutates the account, so the
// placement step is gated behind WEBULL_ORDER_PLACE=1 and should only ever be
// pointed at a sandbox account:
//
//	WEBULL_APP_KEY=...
//	WEBULL_APP_SECRET=...
//	WEBULL_ENVIRONMENT=sandbox
//	WEBULL_ACCOUNT_ID=...
//	WEBULL_ORDER_PLACE=1
//
// The account is taken from WEBULL_ACCOUNT_ID; when it is unset, the first
// account returned by the API is used. Without WEBULL_ORDER_PLACE=1 the
// program previews the order and prints a note, then exits.
//
// Run it with:
//
//	go run ./examples/order
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/trade"
)

// orderLimitPrice is far below the AAPL market price, so the order is not
// marketable and will not execute before it is cancelled.
const orderLimitPrice = "1.00"

func main() {
	if os.Getenv("WEBULL_APP_KEY") == "" {
		fmt.Println("example: WEBULL_APP_KEY not set, skipping")
		return
	}

	cl, err := client.New(client.WithEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = cl.Close() }()

	ctx := context.Background()

	// Trading endpoints require an access token, sent as the x-access-token
	// header. EnsureToken creates or reuses one; in the sandbox the token is
	// activated without 2FA.
	if _, err := cl.EnsureToken(ctx); err != nil {
		log.Fatal(err)
	}

	// The guardrails bound what this example can send, even if the order is
	// later edited: at most 10 shares and a notional of at most 2500.00.
	trading := trade.New(cl,
		trade.WithMaxOrderQuantity("10"),
		trade.WithMaxOrderNotional("2500.00"),
	)
	defer func() { _ = trading.Close() }()

	accountID, err := resolveAccountID(ctx, trading)
	if err != nil {
		log.Fatal(err)
	}

	req := buildOrder(accountID)

	// Preview validates the request and estimates its cost without placing it.
	preview, err := trading.PreviewOrder(ctx, req)
	if err != nil {
		log.Fatalf("PreviewOrder: %v", err)
	}
	order := req.NewOrders[0]
	log.Printf("preview ok symbol=%s side=%s type=%s qty=%s limit=%s estimated_cost=%s estimated_fee=%s",
		order.Symbol, order.Side, order.OrderType, order.Quantity, order.LimitPrice,
		preview.EstimatedCost, preview.EstimatedTransactionFee)

	if os.Getenv("WEBULL_ORDER_PLACE") != "1" {
		log.Println("preview only: set WEBULL_ORDER_PLACE=1 to place this non-marketable limit order and cancel it immediately")
		return
	}

	placeAndCancel(ctx, trading, req)
}

// resolveAccountID returns WEBULL_ACCOUNT_ID when set, otherwise the first
// account available to the authenticated user.
func resolveAccountID(ctx context.Context, trading *trade.Client) (string, error) {
	if id := os.Getenv("WEBULL_ACCOUNT_ID"); id != "" {
		return id, nil
	}
	accounts, err := trading.ListAccounts(ctx)
	if err != nil {
		return "", err
	}
	if len(accounts) == 0 {
		return "", errors.New("no trading accounts available")
	}
	return accounts[0].AccountID, nil
}

// buildOrder returns a small, non-marketable AAPL limit buy with a unique
// client order identifier. Stock orders are NORMAL, EQUITY, and sized by QTY.
func buildOrder(accountID string) trade.PlaceOrderRequest {
	return trade.PlaceOrderRequest{
		AccountID: accountID,
		NewOrders: []trade.OrderRequest{
			{
				ClientOrderID:         fmt.Sprintf("sdk-order-%d", time.Now().UnixNano()),
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
				LimitPrice:            orderLimitPrice,
			},
		},
	}
}

// placeAndCancel submits the order and cancels it with a fresh timeout so the
// cleanup still runs if the caller's context is already done.
func placeAndCancel(ctx context.Context, trading *trade.Client, req trade.PlaceOrderRequest) {
	log.Println("WARNING: WEBULL_ORDER_PLACE=1 is set; placing a real order mutates the account")

	res, err := trading.PlaceOrder(ctx, req)
	if err != nil {
		log.Fatalf("PlaceOrder: %v", err)
	}
	log.Printf("placed client_order_id=%s order_id=%s", res.ClientOrderID, res.OrderID)

	cctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cres, err := trading.CancelOrder(cctx, trade.CancelOrderRequest{
		AccountID:     req.AccountID,
		ClientOrderID: res.ClientOrderID,
	})
	if err != nil {
		log.Fatalf("CancelOrder(%s): %v", res.ClientOrderID, err)
	}
	log.Printf("cancelled client_order_id=%s order_id=%s", cres.ClientOrderID, cres.OrderID)
}
