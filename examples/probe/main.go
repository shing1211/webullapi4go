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

// probe tests data package endpoints and arbitrary paths against the live HK and US sandboxes.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/data"
	"github.com/shing1211/webullapi4go/pkg/errors"
)

func main() {
	appKey := os.Getenv("WEBULL_APP_KEY")
	appSecret := os.Getenv("WEBULL_APP_SECRET")
	baseURL := os.Getenv("WEBULL_BASE_URL")
	regionStr := os.Getenv("WEBULL_REGION")

	if appKey == "" || appSecret == "" {
		fmt.Fprintln(os.Stderr, "WEBULL_APP_KEY and WEBULL_APP_SECRET must be set")
		os.Exit(1)
	}

	region := client.HK
	if regionStr == "US" {
		region = client.US
	}

	opts := []client.Option{
		client.WithAppKey(appKey),
		client.WithAppSecret(appSecret),
		client.WithRegion(region),
		client.WithEnvironment(client.Sandbox),
	}
	if baseURL != "" {
		opts = append(opts, client.WithBaseURL(baseURL))
	}

	cl, err := client.New(opts...)
	if err != nil {
		log.Fatalf("client.New: %v", err)
	}
	defer func() { _ = cl.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if _, err := cl.EnsureToken(ctx); err != nil {
		log.Fatalf("EnsureToken: %v", err)
	}
	fmt.Printf("[OK] Token acquired (region=%s base=%s)\n", region, baseURL)

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

	// === Known working sanity check ===
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

	// === Track A: Crypto bars — corrected path ===
	// The correct path is /market-data/bars with category=CRYPTO&symbol=BTCUSD (query params, not path params).
	fmt.Println("\n--- Crypto Bars (Corrected Path) ---")

	cryptoSymbols := []string{"BTCUSD", "ETHUSD", "BTC-USD", "ETH-USD"}
	cryptoCategories := []string{"CRYPTO", "US_CRYPTO"}
	cryptoIntervals := []string{"d1", "m1", "h1"}

	for _, sym := range cryptoSymbols {
		for _, cat := range cryptoCategories {
			for _, interval := range cryptoIntervals {
				q := fmt.Sprintf("category=%s&symbol=%s&timespan=%s&count=5", cat, sym, interval)
				path := "/market-data/bars?" + q
				name := fmt.Sprintf("crypto_bars_%s_%s_%s", sym, interval, cat)
				probe(name, func(ctx context.Context) any {
					var out json.RawMessage
					err := cl.Do(ctx, http.MethodGet, path, nil, &out)
					if err != nil {
						return err
					}
					return out
				})
			}
		}
	}

	// Try symbols= (plural) as the webullpay.com doc shows
	for _, sym := range []string{"BTCUSD", "ETHUSD"} {
		for _, cat := range []string{"CRYPTO", "US_CRYPTO"} {
			path := fmt.Sprintf("/market-data/bars?category=%s&symbols=%s&timespan=d1&count=5", cat, sym)
			name := fmt.Sprintf("crypto_bars_plural_%s_%s", sym, cat)
			probe(name, func(ctx context.Context) any {
				var out json.RawMessage
				err := cl.Do(ctx, http.MethodGet, path, nil, &out)
				if err != nil {
					return err
				}
				return out
			})
		}
	}

	// Also try the old path patterns (they return 400 but we keep them to confirm)
	fmt.Println("\n--- Crypto Bars (Legacy Paths — expecting 400) ---")
	for _, ticker := range []string{"BTCUSD", "ETHUSD"} {
		for _, interval := range []string{"d1", "m1"} {
			path := fmt.Sprintf("/market-data/crypto/%s/bars?interval=%s&count=5", ticker, interval)
			name := fmt.Sprintf("crypto_legacy_%s_bars_%s", ticker, interval)
			probe(name, func(ctx context.Context) any {
				var out json.RawMessage
				err := cl.Do(ctx, http.MethodGet, path, nil, &out)
				if err != nil {
					return err
				}
				return out
			})
		}
	}

	// === Track B: Display Solution endpoints ===
	fmt.Println("\n--- Display Solution Endpoints ---")

	probe("corp_actions_list", func(ctx context.Context) any {
		v, _, err := dataCl.GetCorporateActions(ctx, data.CorporateActionQuery{
			Symbols:  []string{"AAPL"},
			PageSize: 5,
		})
		if err != nil {
			return err
		}
		return v
	})

	probe("corp_actions_market", func(ctx context.Context) any {
		v, _, err := dataCl.GetCorporateActionsByMarket(ctx, data.CorporateActionQuery{
			Market:   "US",
			PageSize: 5,
		})
		if err != nil {
			return err
		}
		return v
	})

	// Screener v2 POST
	body := map[string]any{
		"filter":    map[string]any{"field": "price", "operator": "gt", "value": "100"},
		"sort":      map[string]any{"field": "price", "direction": "DESC"},
		"page_size": 5,
	}
	for _, path := range []string{
		"/wlas/screener/ng/query",
		"/market-data/screener/ng/query",
		"/market-data/screeners/ng/query",
	} {
		name := "screener_v2_" + strings.ReplaceAll(strings.ReplaceAll(path, "/", "_"), ".", "_")
		probe(name, func(ctx context.Context) any {
			var out json.RawMessage
			err := cl.Do(ctx, http.MethodPost, path, body, &out)
			if err != nil {
				return err
			}
			return out
		})
	}

	fmt.Println("\n=== Summary ===")
	for _, r := range results {
		var icon string
		switch r.Status {
		case 200:
			icon = "✅"
		case 400:
			icon = "❗"
		case 404:
			icon = "❌"
		default:
			icon = "---"
		}
		fmt.Printf("  %s [%d] %s\n", icon, r.Status, r.Name)
	}

	fmt.Println("\n=== Full Results JSON ===")
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
	}
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
