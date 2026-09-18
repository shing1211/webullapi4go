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
	"github.com/shing1211/webullapi4go/internal/errs"
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

	// === Track A: Crypto param sweep ===
	// The path /market-data/crypto/bars returns 400 (path exists, params wrong).
	// Try every combination of ticker formats, categories, and intervals.
	fmt.Println("\n--- Crypto Param Sweep ---")

	cryptoTickers := []string{
		"BTCUSD",
		"BTC-USD",
		"BTC_USD",
		"BTC",
		"ETHUSD",
		"ETH-USD",
		"ETH",
	}
	cryptoCategories := []string{"CRYPTO", "US", "", "ETF", "FUND"}
	intervals := []string{"1d", "1m", "1h"}

	for _, ticker := range cryptoTickers {
		for _, cat := range cryptoCategories {
			for _, interval := range intervals {
				q := fmt.Sprintf("ticker=%s&interval=%s&count=5", ticker, interval)
				if cat != "" {
					q += "&category=" + cat
				}
				path := "/market-data/crypto/bars?" + q
				name := fmt.Sprintf("crypto_bars_%s_%s_%s", ticker, interval, cat)
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

	// Also try with type=ETF in query for crypto (some endpoints use this)
	for _, ticker := range []string{"BTCUSD", "ETHUSD"} {
		for _, interval := range []string{"1d", "1m"} {
			path := fmt.Sprintf("/market-data/crypto/bars?ticker=%s&type=ETF&interval=%s&count=5", ticker, interval)
			name := fmt.Sprintf("crypto_bars_%s_%s_typeETF", ticker, interval)
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

	// Also try /market-data/crypto/{ticker}/bars style path
	for _, ticker := range []string{"BTCUSD", "ETHUSD", "BTC-USD"} {
		for _, interval := range []string{"1d", "1m"} {
			path := fmt.Sprintf("/market-data/crypto/%s/bars?interval=%s&count=5", ticker, interval)
			name := fmt.Sprintf("crypto_%s_bars_%s", ticker, interval)
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

	// Fund / ETF
	for _, path := range []string{
		"/market-data/fund/list",
		"/market-data/etf/list",
		"/market-data/funds/list",
	} {
		name := "fund_" + strings.ReplaceAll(path, "/", "_")
		probe(name, func(ctx context.Context) any {
			var out json.RawMessage
			err := cl.Do(ctx, http.MethodGet, path+"?category=US&page_size=5", nil, &out)
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
