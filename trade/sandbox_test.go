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
	"os"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/trade"
)

// TestSandboxAccountsAssets reads the dedicated trading sandbox account: it
// lists the accounts, then fetches the balance and positions for
// WEBULL_TRADE_ACCOUNT_ID. It is read-only and never places or mutates an
// order.
//
// It runs only when WEBULL_TRADE_SANDBOX=1, WEBULL_TRADE_APP_KEY,
// WEBULL_TRADE_APP_SECRET, and WEBULL_TRADE_ACCOUNT_ID are set; otherwise it
// skips so the default test run stays hermetic. Never commit the sandbox
// credentials to the repository.
func TestSandboxAccountsAssets(t *testing.T) {
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
	// WEBULL_BASE_URL lets a run target a different region's sandbox without
	// changing the test.
	if baseURL := os.Getenv("WEBULL_BASE_URL"); baseURL != "" {
		opts = append(opts, client.WithBaseURL(baseURL))
	}

	cl, err := client.New(opts...)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	if _, err := cl.EnsureToken(ctx); err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}

	trading := trade.New(cl)

	accounts, err := trading.ListAccounts(ctx)
	if err != nil {
		t.Fatalf("ListAccounts() error = %v", err)
	}
	if len(accounts) == 0 {
		t.Fatal("ListAccounts() returned no accounts")
	}
	var found bool
	for _, acct := range accounts {
		t.Logf("account id=%s number=%s type=%s class=%s",
			acct.AccountID, acct.AccountNumber, acct.AccountType, acct.AccountClass)
		if acct.AccountID == accountID {
			found = true
		}
	}
	if !found {
		t.Fatalf("WEBULL_TRADE_ACCOUNT_ID %s not present in %d returned accounts", accountID, len(accounts))
	}

	balance, err := trading.GetBalance(ctx, accountID)
	if err != nil {
		t.Fatalf("GetBalance(%s) error = %v", accountID, err)
	}
	t.Logf("balance currency=%s cash=%s market_value=%s pnl=%s currencies=%d",
		balance.TotalAssetCurrency, balance.TotalCashBalance, balance.TotalMarketValue,
		balance.TotalUnrealizedProfitLoss, len(balance.AccountCurrencyAssets))

	positions, err := trading.GetPositions(ctx, accountID)
	if err != nil {
		t.Fatalf("GetPositions(%s) error = %v", accountID, err)
	}
	t.Logf("positions=%d", len(positions))
	for _, pos := range positions {
		t.Logf("position symbol=%s type=%s qty=%s currency=%s",
			pos.Symbol, pos.InstrumentType, pos.Quantity, pos.Currency)
	}
}
