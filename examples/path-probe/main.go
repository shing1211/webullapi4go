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

// Command path-probe compares the paths the SDK calls against the paths in the
// official Webull OpenAPI definition for endpoints where the two disagree. It
// settles Known Issue 12 in IMPLEMENTATION_STATUS.md.
//
// It is env-gated and skipped by default. To run:
//
//	WEBULL_SANDBOX=1 WEBULL_APP_KEY=... WEBULL_APP_SECRET=... go run ./examples/path-probe
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/data"
	"github.com/shing1211/webullapi4go/internal/errs"
)

func main() {
	if os.Getenv("WEBULL_SANDBOX") != "1" {
		fmt.Println("path-probe: skipped (set WEBULL_SANDBOX=1 with credentials to run)")
		return
	}
	key, secret := os.Getenv("WEBULL_APP_KEY"), os.Getenv("WEBULL_APP_SECRET")
	if key == "" || secret == "" {
		fmt.Println("path-probe: skipped (WEBULL_APP_KEY and WEBULL_APP_SECRET required)")
		return
	}

	c, err := client.New(
		client.WithAppKey(key),
		client.WithAppSecret(secret),
		client.WithSandbox(),
	)
	if err != nil {
		fmt.Println("path-probe: client init failed:", err)
		os.Exit(1)
	}
	defer func() { _ = c.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tok, err := c.CreateToken(ctx)
	report("POST /openapi/auth/token/create  (SDK path)", err)
	if err == nil {
		_, err = c.CheckToken(ctx, tok.Value)
		report("POST /openapi/auth/token/check   (SDK path)", err)
		c.SetToken(tok)
		c.EnableTokenInjection()
	}

	// Alternate auth paths from the OpenAPI JSON definition.
	probe(ctx, c, "POST", "/auth/tokens/create")
	probe(ctx, c, "POST", "/auth/tokens/check")

	// Instrument list: SDK calls /openapi/instrument/stock/list.
	d := data.New(c)
	_, err = d.GetStockInstruments(ctx, data.StockInstrumentQuery{
		Category: data.StockCategoryUS,
		Symbols:  []string{"AAPL"},
	})
	report("GET /openapi/instrument/stock/list (SDK path)", err)

	// Alternate instrument path from the OpenAPI JSON definition.
	probe(ctx, c, "GET", "/trading/instruments/stocks/profiles/list?category=US_STOCK&symbols=AAPL")

	fmt.Println("path-probe: done")
}

func probe(ctx context.Context, c *client.Client, method, path string) {
	var out any
	err := c.Do(ctx, method, path, nil, &out)
	report(fmt.Sprintf("%s %s", method, path), err)
}

func report(label string, err error) {
	fmt.Printf("%-58s %s\n", label, status(err))
}

func status(err error) string {
	if err == nil {
		return "200 OK"
	}
	var e *errs.Error
	if errors.As(err, &e) && e.Status != 0 {
		return fmt.Sprintf("HTTP %d (%s)", e.Status, e.Code)
	}
	return "error: " + err.Error()
}
