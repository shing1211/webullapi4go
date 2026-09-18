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

// Command marketdata fetches a market snapshot and recent daily bars for AAPL
// over the Market Data HTTP API.
//
// Credentials are read from the environment. The sandbox only serves a small
// symbol set, so run it against the sandbox first:
//
//	WEBULL_APP_KEY=...
//	WEBULL_APP_SECRET=...
//	WEBULL_ENVIRONMENT=sandbox
//
// Run it with:
//
//	go run ./examples/marketdata
package main

import (
	"context"
	"log"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/data"
)

func main() {
	cl, err := client.New(client.WithEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = cl.Close() }()

	ctx := context.Background()

	// Market Data endpoints require an access token. EnsureToken creates or
	// reuses one; the sandbox activates it without 2FA.
	if _, err := cl.EnsureToken(ctx); err != nil {
		log.Fatal(err)
	}

	market := data.New(cl)

	snaps, err := market.GetSnapshot(ctx, data.SnapshotQuery{
		Symbols:  []string{"AAPL"},
		Category: data.StockCategoryUS,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range snaps {
		log.Printf("snapshot %s price=%s pre_close=%s change=%s", s.Symbol, s.Price, s.PreClose, s.Change)
	}

	bars, err := market.GetBars(ctx, data.BarQuery{
		Symbol:   "AAPL",
		Category: data.StockCategoryUS,
		Interval: data.BarTimespanDay,
		Count:    5,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range bars.Result {
		log.Printf("bar %s open=%s high=%s low=%s close=%s volume=%s", b.Time, b.Open, b.High, b.Low, b.Close, b.Volume)
	}
}
