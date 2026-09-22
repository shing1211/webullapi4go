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

// Command watchlist-cmd tests watchlist CRUD operations against the sandbox.
// It is guarded by WEBULL_WATCHLIST_TEST=1; no mutations run without it.
//
// Credentials are read from the environment:
//
//	WEBULL_APP_KEY=...
//	WEBULL_APP_SECRET=...
//	WEBULL_ENVIRONMENT=sandbox
//
// Run it with:
//
//	WEBULL_WATCHLIST_TEST=1 go run ./examples/watchlist-cmd
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/data"
)

func main() {
	if os.Getenv("WEBULL_WATCHLIST_TEST") != "1" {
		fmt.Println("watchlist test disabled (set WEBULL_WATCHLIST_TEST=1 to run)")
		os.Exit(0)
	}

	if os.Getenv("WEBULL_APP_KEY") == "" {
		fmt.Println("example: WEBULL_APP_KEY not set, skipping")
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

	market := data.New(cl)

	fmt.Println("[GetWatchlists]")
	lists, err := market.GetWatchlists(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var watchlistID string
	if len(lists) > 0 {
		watchlistID = lists[0].WatchlistID
		fmt.Printf("using existing watchlist id=%s name=%q\n", watchlistID, lists[0].Name)
	} else {
		fmt.Println("no watchlists found, creating one")
		fmt.Println("[CreateWatchlist]")
		res, err := market.CreateWatchlist(ctx, data.CreateWatchlistParams{Name: "test-wl"})
		if err != nil {
			log.Fatal(err)
		}
		watchlistID = res.WatchlistID
		fmt.Printf("created watchlist id=%s\n", watchlistID)
	}

	fmt.Println("[AddWatchlistInstruments]")
	_, err = market.AddWatchlistInstruments(ctx, data.WatchlistInstrumentsParam{
		WatchlistID: watchlistID,
		Instruments: []data.WatchlistInstrumentParam{
			{Symbol: "AAPL", Category: data.StockCategoryUS},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("added AAPL")

	fmt.Println("[UpdateWatchlist]")
	_, err = market.UpdateWatchlist(ctx, data.UpdateWatchlistParams{
		WatchlistID: watchlistID,
		Name:        "test-wl-renamed",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("renamed to test-wl-renamed")

	fmt.Println("[RemoveWatchlistInstruments]")
	_, err = market.RemoveWatchlistInstruments(ctx, data.WatchlistInstrumentsParam{
		WatchlistID: watchlistID,
		Instruments: []data.WatchlistInstrumentParam{
			{Symbol: "AAPL", Category: data.StockCategoryUS},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("removed AAPL")

	fmt.Println("[DeleteWatchlist]")
	_, err = market.DeleteWatchlist(ctx, watchlistID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("deleted watchlist")

	fmt.Println("all operations completed")
}
