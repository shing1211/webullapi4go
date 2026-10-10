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

// Command options demonstrates a multi-leg US options order preview. It builds a
// vertical call spread (buy lower strike, sell higher strike) and previews its
// estimated cost without placing it.
//
// This is preview-only: no order is submitted to the account.
//
// Credentials are read from the environment:
//
//	WEBULL_APP_KEY=...
//	WEBULL_APP_SECRET=...
//	WEBULL_ACCOUNT_ID=...
//	WEBULL_ENVIRONMENT=sandbox
//
// Run it with:
//
//	go run ./examples/options
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/pkg/domain/money"
	"github.com/shing1211/webullapi4go/trade"
)

func moneyPtr(s string) *money.Money {
	m := money.Must(money.NewFromString(s))
	return &m
}

func main() {
	if os.Getenv("WEBULL_APP_KEY") == "" {
		fmt.Println("credentials not set, skipping options example")
		return
	}

	cl, err := client.New(client.WithEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = cl.Close() }()

	ctx := context.Background()

	if _, err := cl.EnsureToken(ctx); err != nil {
		log.Fatal(err)
	}

	trading := trade.New(cl,
		trade.WithMaxOrderQuantity("10"),
		trade.WithMaxOrderNotional("2500.00"),
	)
	defer func() { _ = trading.Close() }()

	accountID, err := resolveAccountID(ctx, trading)
	if err != nil {
		log.Fatal(err)
	}

	req := buildVerticalCallSpread(accountID)

	preview, err := trading.PreviewOrder(ctx, req)
	if err != nil {
		log.Fatalf("PreviewOrder: %v", err)
	}
	order := req.NewOrders[0]
	log.Printf("preview ok symbol=%s strategy=%s side=%s type=%s qty=%s limit=%s estimated_cost=%s estimated_fee=%s",
		order.Symbol, order.OptionStrategy, order.Side, order.OrderType, order.Quantity, order.LimitPrice,
		preview.EstimatedCost, preview.EstimatedTransactionFee)
}

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

func buildVerticalCallSpread(accountID string) trade.PlaceOrderRequest {
	expiration := "2026-12-18"
	return trade.PlaceOrderRequest{
		AccountID: accountID,
		NewOrders: []trade.OrderRequest{
			{
				ClientOrderID:  fmt.Sprintf("sdk-option-%d", time.Now().UnixNano()),
				ComboType:      trade.ComboTypeNormal,
				InstrumentType: trade.InstrumentTypeOption,
				Market:         trade.MarketUS,
				Symbol:         "AAPL",
				OrderType:      trade.OrderTypeLimit,
				Side:           trade.OrderSideBuy,
				Quantity:       moneyPtr("1"),
				EntrustType:    trade.EntrustTypeQty,
				TimeInForce:    trade.TimeInForceDay,
				LimitPrice:     moneyPtr("1.50"),
				OptionStrategy: trade.OptionStrategyVertical,
				Legs: []trade.OrderLeg{
					{
						InstrumentType:   trade.InstrumentTypeOption,
						Market:           trade.MarketUS,
						Symbol:           "AAPL",
						Side:             trade.OrderSideBuy,
						StrikePrice:      moneyPtr("220.00"),
						OptionExpireDate: expiration,
						OptionType:       trade.OptionTypeCall,
						Quantity:         moneyPtr("1"),
					},
					{
						InstrumentType:   trade.InstrumentTypeOption,
						Market:           trade.MarketUS,
						Symbol:           "AAPL",
						Side:             trade.OrderSideSell,
						StrikePrice:      moneyPtr("230.00"),
						OptionExpireDate: expiration,
						OptionType:       trade.OptionTypeCall,
						Quantity:         moneyPtr("1"),
					},
				},
			},
		},
	}
}
