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

package data_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/data"
)

// TestSandboxStockInstruments calls the shared Webull sandbox and asserts that
// the stock instrument list returns a non-empty result for AAPL.
//
// It runs only when WEBULL_SANDBOX=1 and WEBULL_APP_KEY / WEBULL_APP_SECRET are
// set; otherwise it skips so the default test run stays hermetic. Never commit
// the sandbox credentials to the repository. Sandbox market-data access is
// limited to AAPL.
func TestSandboxStockInstruments(t *testing.T) {
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
	// WEBULL_BASE_URL lets a run target a different region's sandbox (for
	// example https://api.sandbox.webull.com) without changing the test.
	if baseURL := os.Getenv("WEBULL_BASE_URL"); baseURL != "" {
		opts = append(opts, client.WithBaseURL(baseURL))
	}

	cl, err := client.New(opts...)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	market := data.New(cl)
	instruments, err := market.GetStockInstruments(ctx, data.StockInstrumentQuery{
		Category: data.StockCategoryUS,
		Symbols:  []string{"AAPL"},
	})
	if err != nil {
		t.Fatalf("GetStockInstruments(AAPL, US) error = %v", err)
	}
	if len(instruments) == 0 {
		t.Fatal("GetStockInstruments(AAPL, US) returned no instruments")
	}

	var found bool
	for _, inst := range instruments {
		if inst.Symbol == "AAPL" {
			found = true
			t.Logf("AAPL instrument: id=%s exchange=%s currency=%s status=%s",
				inst.InstrumentID, inst.ExchangeCode, inst.Currency, inst.Status)
		}
	}
	if !found {
		t.Fatalf("AAPL not present in %d returned instruments", len(instruments))
	}
}
