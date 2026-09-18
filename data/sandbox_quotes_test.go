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

// TestSandboxQuotes exercises the market-data endpoints added in T4.2 against
// the shared Webull sandbox: AAPL snapshot, ticks, order-book depth, 1-minute
// bars and the batch-bars request.
//
// It runs only when WEBULL_SANDBOX=1 and WEBULL_APP_KEY / WEBULL_APP_SECRET are
// set; otherwise it skips so the default test run stays hermetic. Never commit
// the sandbox credentials to the repository. Sandbox market-data access is
// limited to AAPL.
func TestSandboxQuotes(t *testing.T) {
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

	cl, err := client.New(opts...)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	token, err := cl.EnsureToken(ctx)
	if err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}
	t.Logf("token status = %s", token.Status)

	market := data.New(cl)

	snapshots, err := market.GetSnapshot(ctx, data.SnapshotQuery{
		Symbols:  []string{"AAPL"},
		Category: data.StockCategoryUS,
	})
	if err != nil {
		t.Errorf("GetSnapshot(AAPL) error = %v", err)
	} else if len(snapshots) == 0 {
		t.Errorf("GetSnapshot(AAPL) returned no snapshots")
	} else {
		s := snapshots[0]
		t.Logf("snapshot: symbol=%s price=%s change=%s pre_close=%s last_trade_time=%d",
			s.Symbol, s.Price, s.Change, s.PreClose, s.LastTradeTime)
	}

	ticks, err := market.GetTick(ctx, data.TickQuery{
		Symbol:          "AAPL",
		Category:        data.StockCategoryUS,
		Count:           30,
		TradingSessions: []data.TradingSession{data.TradingSessionRTH},
	})
	if err != nil {
		t.Errorf("GetTick(AAPL) error = %v", err)
	} else {
		t.Logf("ticks: symbol=%s instrument_id=%s count=%d", ticks.Symbol, ticks.InstrumentID, len(ticks.Result))
		if len(ticks.Result) > 0 {
			tick := ticks.Result[0]
			t.Logf("first tick: time=%s price=%s volume=%s side=%s", tick.Time, tick.Price, tick.Volume, tick.Side)
		}
	}

	quote, err := market.GetQuotes(ctx, data.DepthQuery{
		Symbol:   "AAPL",
		Category: data.StockCategoryUS,
		Depth:    10,
	})
	if err != nil {
		t.Errorf("GetQuotes(AAPL) error = %v", err)
	} else {
		t.Logf("quotes: symbol=%s quote_time=%d asks=%d bids=%d",
			quote.Symbol, quote.QuoteTime, len(quote.Asks), len(quote.Bids))
		if len(quote.Bids) > 0 {
			t.Logf("best bid: price=%s size=%s", quote.Bids[0].Price, quote.Bids[0].Size)
		}
		if len(quote.Asks) > 0 {
			t.Logf("best ask: price=%s size=%s", quote.Asks[0].Price, quote.Asks[0].Size)
		}
	}

	bars, err := market.GetBars(ctx, data.BarQuery{
		Symbol:          "AAPL",
		Category:        data.StockCategoryUS,
		Interval:        data.BarTimespanM1,
		Count:           5,
		TradingSessions: []data.TradingSession{data.TradingSessionRTH},
	})
	if err != nil {
		t.Errorf("GetBars(AAPL, M1) error = %v", err)
	} else {
		t.Logf("bars: symbol=%s instrument_id=%s count=%d", bars.Symbol, bars.InstrumentID, len(bars.Result))
		if len(bars.Result) > 0 {
			bar := bars.Result[len(bars.Result)-1]
			t.Logf("last bar: time=%s open=%s high=%s low=%s close=%s volume=%s",
				bar.Time, bar.Open, bar.High, bar.Low, bar.Close, bar.Volume)
		}
	}

	batch, err := market.GetBatchBars(ctx, data.BatchBarQuery{
		Symbols:  []string{"AAPL"},
		Category: data.StockCategoryUS,
		Timespan: data.BarTimespanM1,
		Count:    5,
	})
	if err != nil {
		t.Errorf("GetBatchBars(AAPL, M1) error = %v", err)
	} else {
		t.Logf("batch bars: symbols=%d", len(batch.Result))
		for _, r := range batch.Result {
			t.Logf("batch bars: symbol=%s instrument_id=%s count=%d", r.Symbol, r.InstrumentID, len(r.Result))
		}
	}
}
