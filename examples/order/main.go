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
//	WEBULL_ORDER_IDEMPOTENCY_KEY=...
//	WEBULL_ORDER_PLACE=1
//
// The account is taken from WEBULL_ACCOUNT_ID; when it is unset, the first
// account returned by the API is used. Mutating placement also requires a
// caller-generated WEBULL_ORDER_IDEMPOTENCY_KEY that is persisted before the
// first request and reused for every retry. Without WEBULL_ORDER_PLACE=1 the
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
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/pkg/domain/money"
	"github.com/shing1211/webullapi4go/trade"
)

const (
	orderLimitPrice          = "1.00"
	orderIDEnvironment       = "WEBULL_ORDER_IDEMPOTENCY_KEY"
	maxClientOrderIDLength   = 32
	orderTokenTimeout        = 5 * time.Minute
	orderRequestTimeout      = 20 * time.Second
	orderCancellationTimeout = 30 * time.Second
	previewOrderIDSeed       = "preview:AAPL:BUY:1:1.00"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if os.Getenv("WEBULL_APP_KEY") == "" {
		fmt.Println("example: WEBULL_APP_KEY not set, skipping")
		return nil
	}

	placing := os.Getenv("WEBULL_ORDER_PLACE") == "1"
	clientOrderID, err := resolveClientOrderID(placing)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cl, err := client.New(client.WithEnv())
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}
	defer func() { _ = cl.Close() }()

	tokenCtx, cancelToken := context.WithTimeout(ctx, orderTokenTimeout)
	_, err = cl.EnsureToken(tokenCtx)
	cancelToken()
	if err != nil {
		return fmt.Errorf("ensure access token: %w", err)
	}

	trading := trade.New(cl,
		trade.WithMaxOrderQuantity("10"),
		trade.WithMaxOrderNotional("2500.00"),
	)
	defer func() { _ = trading.Close() }()

	accountCtx, cancelAccount := context.WithTimeout(ctx, orderRequestTimeout)
	accountID, err := resolveAccountID(accountCtx, trading)
	cancelAccount()
	if err != nil {
		return fmt.Errorf("resolve account: %w", err)
	}

	req := buildOrder(accountID, clientOrderID)
	previewCtx, cancelPreview := context.WithTimeout(ctx, orderRequestTimeout)
	preview, err := trading.PreviewOrder(previewCtx, req)
	cancelPreview()
	if err != nil {
		return fmt.Errorf("preview order: %w", err)
	}

	order := req.NewOrders[0]
	log.Printf("preview ok client_order_id=%s symbol=%s side=%s type=%s qty=%s limit=%s estimated_cost=%s estimated_fee=%s",
		order.ClientOrderID, order.Symbol, order.Side, order.OrderType,
		order.Quantity.String(), order.LimitPrice.String(),
		preview.EstimatedCost.String(), preview.EstimatedTransactionFee.String())

	if !placing {
		log.Println("preview only: set WEBULL_ORDER_PLACE=1 and WEBULL_ORDER_IDEMPOTENCY_KEY to place and immediately cancel this limit order")
		return nil
	}

	return placeAndCancel(ctx, trading, req)
}

func resolveClientOrderID(placing bool) (string, error) {
	key := strings.TrimSpace(os.Getenv(orderIDEnvironment))
	if key == "" {
		if placing {
			return "", fmt.Errorf("%s is required for placement; persist a new key before the first request and reuse it for retries", orderIDEnvironment)
		}
		return trade.ClientOrderIDFrom([]byte(previewOrderIDSeed)), nil
	}
	if len(key) > maxClientOrderIDLength || !trade.ValidClientOrderID(key) {
		return "", fmt.Errorf("%s must be at most %d characters using only letters, digits, '-' and '_'", orderIDEnvironment, maxClientOrderIDLength)
	}
	return key, nil
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

func buildOrder(accountID, clientOrderID string) trade.PlaceOrderRequest {
	quantity := money.MustNew("1")
	limitPrice := money.MustNew(orderLimitPrice)
	order := trade.NewEquityOrderBuilder("AAPL", trade.OrderSideBuy, &quantity).
		LimitPrice(&limitPrice).
		Build()
	order.ClientOrderID = clientOrderID
	return trade.NewPlaceOrderRequest(accountID, order)
}

func placeAndCancel(ctx context.Context, trading *trade.Client, req trade.PlaceOrderRequest) error {
	clientOrderID := req.NewOrders[0].ClientOrderID
	log.Println("WARNING: WEBULL_ORDER_PLACE=1 is set; placing a real order mutates the account")

	placeCtx, cancelPlace := context.WithTimeout(ctx, orderRequestTimeout)
	res, err := trading.PlaceOrder(placeCtx, req)
	cancelPlace()
	if err != nil {
		return fmt.Errorf("place order with client_order_id %q (retry with the same key if the outcome is unknown): %w", clientOrderID, err)
	}
	log.Printf("placed client_order_id=%s order_id=%s", res.ClientOrderID, res.OrderID)

	cancelCtx, cancelCancel := context.WithTimeout(context.Background(), orderCancellationTimeout)
	cancelResult, err := trading.CancelOrder(cancelCtx, trade.NewCancelOrderRequest(req.AccountID, clientOrderID))
	cancelCancel()
	if err != nil {
		return fmt.Errorf("cancel client_order_id %q; the order may still be working: %w", clientOrderID, err)
	}
	log.Printf("cancelled client_order_id=%s order_id=%s", cancelResult.ClientOrderID, cancelResult.OrderID)
	return nil
}
