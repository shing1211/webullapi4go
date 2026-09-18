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

// Command account lists the authenticated user's trading accounts and, for one
// account, prints its balance and open positions. It is read-only: it never
// places, replaces, or cancels an order.
//
// Credentials are read from the environment:
//
//	WEBULL_APP_KEY=...
//	WEBULL_APP_SECRET=...
//	WEBULL_ENVIRONMENT=sandbox
//
// The account to inspect is taken from WEBULL_ACCOUNT_ID. When that is unset,
// the first account returned by the API is used. Run it with:
//
//	go run ./examples/account
package main

import (
	"context"
	"log"
	"os"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/trade"
)

func main() {
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

	trading := trade.New(cl)
	defer func() { _ = trading.Close() }()

	accounts, err := trading.ListAccounts(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if len(accounts) == 0 {
		log.Println("no trading accounts")
		return
	}
	for _, acct := range accounts {
		log.Printf("account id=%s number=%s type=%s class=%s",
			acct.AccountID, acct.AccountNumber, acct.AccountType, acct.AccountClass)
	}

	accountID := os.Getenv("WEBULL_ACCOUNT_ID")
	if accountID == "" {
		accountID = accounts[0].AccountID
	}

	balance, err := trading.GetBalance(ctx, accountID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("balance %s currency=%s cash=%s market_value=%s pnl=%s",
		accountID, balance.TotalAssetCurrency, balance.TotalCashBalance,
		balance.TotalMarketValue, balance.TotalUnrealizedProfitLoss)

	positions, err := trading.GetPositions(ctx, accountID)
	if err != nil {
		log.Fatal(err)
	}
	if len(positions) == 0 {
		log.Printf("account %s holds no positions", accountID)
		return
	}
	for _, pos := range positions {
		log.Printf("position symbol=%s type=%s qty=%s currency=%s last=%s cost=%s pnl=%s",
			pos.Symbol, pos.InstrumentType, pos.Quantity, pos.Currency,
			pos.LastPrice, pos.CostPrice, pos.UnrealizedProfitLoss)
	}
}
