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

// Command brokerfd demonstrates read-only Broker FD (US) API operations.
// It is read-only: it never places, replaces, or cancels an order.
//
// Credentials are read from the environment:
//
//	WEBULL_APP_KEY=...
//	WEBULL_APP_SECRET=...
//	WEBULL_TRADE_ACCOUNT_ID=...
//	WEBULL_ENVIRONMENT=sandbox
//
// Run it with:
//
//	go run ./examples/brokerfd
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/webullapi4go/brokerfd"
	"github.com/shing1211/webullapi4go/client"
)

func main() {
	if os.Getenv("WEBULL_APP_KEY") == "" {
		fmt.Println("credentials not set, skipping brokerfd example")
		return
	}

	cl, err := client.New(client.WithEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = cl.Close() }()

	bfd := brokerfd.New(cl)
	ctx := context.Background()

	summary, err := bfd.GetAccountsSummary(ctx)
	if err != nil {
		log.Printf("GetAccountsSummary error: %v", err)
	} else {
		fmt.Printf("Accounts Summary:\n")
		for _, s := range summary {
			fmt.Printf("  account_id=%s currency=%s net_liquidity=%s cash=%s market_value=%s buying_power=%s\n",
				s.AccountID, s.Currency, s.NetLiquidity, s.CashBalance, s.MarketValue, s.BuyingPower)
		}
	}

	positions, err := bfd.GetPositions(ctx)
	if err != nil {
		log.Printf("GetPositions error: %v", err)
	} else {
		fmt.Printf("Positions:\n")
		for _, p := range positions {
			fmt.Printf("  symbol=%s quantity=%s market_value=%s cost_basis=%s unrealized_pl=%s\n",
				p.Symbol, p.Quantity, p.MarketValue, p.CostBasis, p.UnrealizedPL)
		}
	}

	accountID := os.Getenv("WEBULL_TRADE_ACCOUNT_ID")
	if accountID != "" {
		orders, err := bfd.GetFDOpenOrders(ctx, accountID)
		if err != nil {
			log.Printf("GetFDOpenOrders error: %v", err)
		} else {
			fmt.Printf("Open Orders for account %s:\n", accountID)
			for _, o := range orders {
				fmt.Printf("  order_id=%s symbol=%s side=%s qty=%s filled=%s status=%s\n",
					o.OrderID, o.Symbol, o.Side, o.Quantity, o.FilledQuantity, o.Status)
			}
		}
	}
}
