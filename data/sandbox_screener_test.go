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

// newSandboxClient returns a sandbox market-data client, skipping the test when
// the sandbox environment variables are absent. The caller owns the underlying
// client via t.Cleanup.
//
// It runs only when WEBULL_SANDBOX=1 and WEBULL_APP_KEY / WEBULL_APP_SECRET are
// set; never commit the sandbox credentials to the repository.
func newSandboxClient(t *testing.T) (*client.Client, *data.Client) {
	t.Helper()
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
	return cl, data.New(cl)
}

// TestSandboxMarketDataScreener exercises the footprint, NOII, and screener
// endpoints against the shared Webull sandbox.
//
// The screener endpoints are market-wide and assert a non-empty result.
// Footprint and NOII are US-stock-only order-flow feeds and the sandbox may not
// permit them, so their outcome is reported without failing the test.
func TestSandboxMarketDataScreener(t *testing.T) {
	cl, market := newSandboxClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if _, err := cl.EnsureToken(ctx); err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}

	t.Run("gainers-losers", func(t *testing.T) {
		got, err := market.GetTopGainersLosers(ctx, data.GainersLosersQuery{
			RankType: data.GainersLosersRankDay1,
			Category: data.StockCategoryUS,
			SortBy:   data.ScreenerSortChangeRatio,
		})
		if err != nil {
			t.Fatalf("GetTopGainersLosers() error = %v", err)
		}
		if len(got) == 0 {
			t.Fatal("GetTopGainersLosers() returned no stocks")
		}
		first := got[0]
		t.Logf("gainers/losers: %d stocks, top %s change_ratio=%s", len(got), first.Symbol, first.ChangeRatio)
	})

	t.Run("most-active", func(t *testing.T) {
		got, err := market.GetMostActive(ctx, data.MostActiveQuery{
			Category: data.StockCategoryUS,
			RankType: data.MostActiveRankVolume,
			SortBy:   data.ScreenerSortVolume,
		})
		if err != nil {
			t.Fatalf("GetMostActive() error = %v", err)
		}
		if len(got) == 0 {
			t.Fatal("GetMostActive() returned no stocks")
		}
		t.Logf("most-active: %d stocks, top %s volume=%s", len(got), got[0].Symbol, got[0].Volume)
	})

	t.Run("footprint", func(t *testing.T) {
		got, err := market.GetFootprint(ctx, data.FootprintQuery{
			Symbols:  []string{"AAPL"},
			Category: data.StockCategoryUS,
			Timespan: data.FootprintTimespanM1,
		})
		if err != nil {
			t.Logf("GetFootprint(AAPL) unavailable in sandbox: %v", err)
			return
		}
		t.Logf("footprint: AAPL returned %d symbols", len(got))
	})

	t.Run("noii-bars", func(t *testing.T) {
		got, err := market.GetNOIIBars(ctx, data.NOIIQuery{
			Symbol:              "AAPL",
			Category:            data.StockCategoryUS,
			ImbalanceActionType: data.NOIIActionPreOpen,
		})
		if err != nil {
			t.Logf("GetNOIIBars(AAPL) unavailable in sandbox: %v", err)
			return
		}
		t.Logf("noii-bars: AAPL returned %d bars", len(got))
	})

	t.Run("noii-snapshot", func(t *testing.T) {
		got, err := market.GetNOIISnapshot(ctx, data.NOIIQuery{
			Symbol:              "AAPL",
			Category:            data.StockCategoryUS,
			ImbalanceActionType: data.NOIIActionPreOpen,
		})
		if err != nil {
			t.Logf("GetNOIISnapshot(AAPL) unavailable in sandbox: %v", err)
			return
		}
		t.Logf("noii-snapshot: AAPL shares=%s side=%s", got.ImbalanceShares, got.ImbalanceSide)
	})
}
