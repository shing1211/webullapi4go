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

// probe tests data package endpoints and arbitrary paths against the live HK sandbox.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/data"
	"github.com/shing1211/webullapi4go/internal/errs"
)

func main() {
	appKey := os.Getenv("WEBULL_APP_KEY")
	appSecret := os.Getenv("WEBULL_APP_SECRET")
	if appKey == "" || appSecret == "" {
		fmt.Fprintln(os.Stderr, "WEBULL_APP_KEY and WEBULL_APP_SECRET must be set")
		os.Exit(1)
	}

	cl, err := client.New(
		client.WithAppKey(appKey),
		client.WithAppSecret(appSecret),
		client.WithRegion(client.HK),
		client.WithEnvironment(client.Sandbox),
	)
	if err != nil {
		log.Fatalf("client.New: %v", err)
	}
	defer cl.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if _, err := cl.EnsureToken(ctx); err != nil {
		log.Fatalf("EnsureToken: %v", err)
	}
	fmt.Println("[OK] Token acquired")

	dataCl := data.New(cl)

	type result struct {
		Name   string `json:"name"`
		Status int    `json:"status"`
		Error  string `json:"error,omitempty"`
		Data   any    `json:"data,omitempty"`
	}
	var results []result

	probe := func(name string, fn func(context.Context) any) {
		v := fn(ctx)
		r := result{Name: name}
		switch val := v.(type) {
		case nil:
			r.Status = 200
		case error:
			r.Status = statusFromError(val)
			r.Error = val.Error()
		case *http.Response:
			r.Status = val.StatusCode
		default:
			r.Status = 200
			r.Data = val
		}
		results = append(results, r)
		fmt.Printf("[%-3d] %s\n", r.Status, name)
		if r.Error != "" && r.Status != 200 && r.Status != 404 {
			fmt.Printf("      ERR: %s\n", r.Error)
		}
	}

	// --- Known working endpoints (sanity check) ---
	probe("gainers_losers", func(ctx context.Context) any {
		v, err := dataCl.GetTopGainersLosers(ctx, data.GainersLosersQuery{
			RankType: data.GainersLosersRankDay1,
			Category: data.StockCategoryUS,
			SortBy:   data.ScreenerSortChangeRatio,
		})
		if err != nil {
			return err
		}
		return v
	})

	probe("bars_aapl", func(ctx context.Context) any {
		v, err := dataCl.GetBars(ctx, data.BarQuery{
			Symbol:   "AAPL",
			Category: data.StockCategoryUS,
			Interval: data.BarTimespanDay,
			Count:    5,
		})
		if err != nil {
			return err
		}
		return v
	})

	// --- Typed data package calls (our stubs) ---
	probe("corp_actions_list", func(ctx context.Context) any {
		v, err := dataCl.GetCorporateActions(ctx, data.CorporateActionQuery{
			Symbols:  []string{"AAPL"},
			PageSize: 5,
		})
		if err != nil {
			return err
		}
		return v
	})

	probe("corp_actions_market", func(ctx context.Context) any {
		v, err := dataCl.GetCorporateActionsByMarket(ctx, data.CorporateActionQuery{
			Market:   "US",
			PageSize: 5,
		})
		if err != nil {
			return err
		}
		return v
	})

	// --- Arbitrary path probing using cl.Do ---
	// Try with category=HK for HK sandbox
	type pathResult struct {
		path  string
		query map[string]string
	}

	pathsToProbe := []string{
		"/market-data/crypto/bars?ticker=BTCUSD&interval=1d&count=5&category=US",
		"/market-data/crypto/bars?ticker=BTCUSD&interval=1m&count=5&category=US",
		"/market-data/crypto/bars?ticker=BTCUSD&interval=1d&count=5",
		"/market-data/crypto/bars?ticker=BTCUSD&interval=1m&count=5",
		"/market-data/crypto/ETHUSD/bars?interval=1d&count=5",
		"/market-data/fund/list?category=US&page_size=5",
		"/market-data/etf/list?category=US&page_size=5",
	}

	body := map[string]any{
		"filter":   map[string]any{"field": "price", "operator": "gt", "value": "100"},
		"sort":      map[string]any{"field": "price", "direction": "DESC"},
		"page_size": 5,
	}

	for _, path := range pathsToProbe {
		name := path
		probe(name, func(ctx context.Context) any {
			var out json.RawMessage
			err := cl.Do(ctx, http.MethodGet, path, nil, &out)
			if err != nil {
				return err
			}
			return out
		})
	}

	// POST screener v2
	for _, path := range []string{"/wlas/screener/ng/query", "/market-data/screener/ng/query"} {
		probe(path, func(ctx context.Context) any {
			var out json.RawMessage
			err := cl.Do(ctx, http.MethodPost, path, body, &out)
			if err != nil {
				return err
			}
			return out
		})
	}

	// Screener v2 guessed paths
	for _, path := range []string{
		"/wlas/screener/ng/query",
		"/market-data/screener/ng/query",
	} {
		name := path
		body := map[string]any{
			"filter":  map[string]any{"field": "price", "operator": "gt", "value": "100"},
			"sort":    map[string]any{"field": "price", "direction": "DESC"},
			"page_size": 5,
		}
		probe(name, func(ctx context.Context) any {
			var out json.RawMessage
			err := cl.Do(ctx, http.MethodPost, path, body, &out)
			if err != nil {
				return err
			}
			return out
		})
	}

	fmt.Println("\n=== Full Results ===")
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(results)
}

func statusFromError(err error) int {
	if err == nil {
		return 200
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return 408
	}
	var se *errs.Error
	if errors.As(err, &se) {
		return se.Status
	}
	return 500
}
