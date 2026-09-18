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
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/data"
)

// sandboxOptionSymbol is the AAPL call contract used by the Webull
// documentation examples. Sandbox market data is limited to AAPL.
const sandboxOptionSymbol = "AAPL260522C00300000"

// TestSandboxOptionMarketData exercises the option tick, snapshot, and
// historical-bars endpoints against the shared Webull sandbox.
//
// The sandbox does not guarantee option-chain coverage: an option endpoint may
// reject the documented AAPL contract. Each sub-test therefore reports an
// unavailable feed instead of failing, matching the treatment of the other
// limited sandbox market-data feeds.
func TestSandboxOptionMarketData(t *testing.T) {
	cl, market := newT44SandboxClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if _, err := cl.EnsureToken(ctx); err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}

	t.Run("snapshot", func(t *testing.T) {
		got, err := market.GetOptionSnapshot(ctx, data.OptionSnapshotQuery{
			Symbols: []string{sandboxOptionSymbol},
		})
		if err != nil {
			t.Logf("GetOptionSnapshot(%s) unavailable in sandbox: %v", sandboxOptionSymbol, err)
			return
		}
		t.Logf("option snapshots: %d returned", len(got))
		for _, snap := range got {
			t.Logf("snapshot %s price=%s bid=%s ask=%s", snap.Symbol, snap.Price, snap.Bid, snap.Ask)
		}
	})

	t.Run("tick", func(t *testing.T) {
		got, err := market.GetOptionTick(ctx, data.OptionTickQuery{
			Symbol: sandboxOptionSymbol,
			Count:  30,
		})
		if err != nil {
			t.Logf("GetOptionTick(%s) unavailable in sandbox: %v", sandboxOptionSymbol, err)
			return
		}
		t.Logf("option ticks: %d returned", len(got.Result))
	})

	t.Run("bars", func(t *testing.T) {
		got, err := market.GetOptionBars(ctx, data.OptionBarsQuery{
			Symbols:  []string{sandboxOptionSymbol},
			Timespan: data.OptionBarTimespanD,
			Count:    5,
		})
		if err != nil {
			t.Logf("GetOptionBars(%s) unavailable in sandbox: %v", sandboxOptionSymbol, err)
			return
		}
		t.Logf("option bars: %d symbols returned", len(got))
	})
}
