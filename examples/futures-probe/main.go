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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/data"
)

func main() {
	if os.Getenv("WEBULL_FUTURES_TEST") != "1" {
		fmt.Println("futures probe disabled (set WEBULL_FUTURES_TEST=1 to run)")
		os.Exit(0)
	}

	ctx := context.Background()
	cl, err := client.New(client.WithEnv())
	if err != nil {
		log.Fatalf("client.New: %v", err)
	}
	if _, err := cl.EnsureToken(ctx); err != nil {
		log.Fatalf("EnsureToken: %v", err)
	}
	dataCl := data.New(cl)

	fmt.Println("=== HK Futures Probe ===")
	fmt.Println()

	products, err := dataCl.GetFuturesProductCodes(ctx, data.FuturesProductCodeQuery{Category: data.FuturesCategoryHK})
	if err != nil {
		log.Fatalf("GetFuturesProductCodes: %v", err)
	}
	if len(products) == 0 {
		fmt.Println("no HK futures products found")
		os.Exit(0)
	}
	fmt.Printf("GetFuturesProductCodes: PASS count=%d\n", len(products))
	firstProduct := products[0]
	fmt.Printf("  -> using product code=%s name=%s\n", firstProduct.Code, firstProduct.Name)

	fmt.Println()

	instruments, err := dataCl.GetFuturesInstruments(ctx, data.FuturesInstrumentQuery{Category: data.FuturesCategoryHK, Code: firstProduct.Code})
	if err != nil {
		log.Fatalf("GetFuturesInstruments: %v", err)
	}
	fmt.Printf("GetFuturesInstruments: PASS count=%d\n", len(instruments))

	var firstSymbol string
	if len(instruments) > 0 {
		firstSymbol = instruments[0].Symbol
		fmt.Printf("  -> using symbol=%s\n", firstSymbol)
	} else {
		firstSymbol = firstProduct.Code
		fmt.Printf("  -> no instruments returned, trying symbol=%s\n", firstSymbol)
	}

	fmt.Println()

	bars, err := dataCl.GetFuturesBars(ctx, data.FuturesBarsQuery{
		Symbols:  []string{firstSymbol},
		Category: data.FuturesCategoryHK,
		Interval: data.BarTimespanDay,
		Count:    5,
	})
	if err != nil {
		fmt.Printf("GetFuturesBars: FAIL %v\n", err)
	} else {
		fmt.Println("GetFuturesBars: PASS")
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(bars); err != nil {
			fmt.Printf("  JSON encode error: %v\n", err)
		}
	}

	fmt.Println()

	snapshots, err := dataCl.GetFuturesSnapshot(ctx, data.FuturesSnapshotQuery{
		Symbols:  []string{firstSymbol},
		Category: data.FuturesCategoryHK,
	})
	if err != nil {
		fmt.Printf("GetFuturesSnapshot: FAIL %v\n", err)
	} else {
		fmt.Printf("GetFuturesSnapshot: PASS count=%d\n", len(snapshots))
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(snapshots); err != nil {
			fmt.Printf("  JSON encode error: %v\n", err)
		}
	}

	_ = cl.Close()
}
