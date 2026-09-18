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

// Command watchlist lists the authenticated user's watchlists. It is read-only:
// it never creates, updates, or deletes anything.
//
// Credentials are read from the environment:
//
//	WEBULL_APP_KEY=...
//	WEBULL_APP_SECRET=...
//	WEBULL_ENVIRONMENT=sandbox
//
// Run it with:
//
//	go run ./examples/watchlist
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

	// Watchlist endpoints require an access token.
	if _, err := cl.EnsureToken(ctx); err != nil {
		log.Fatal(err)
	}

	market := data.New(cl)

	lists, err := market.GetWatchlists(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if len(lists) == 0 {
		log.Println("no watchlists")
		return
	}
	for _, w := range lists {
		log.Printf("watchlist id=%s name=%q sort=%d", w.WatchlistID, w.Name, w.Sort)
	}
}
