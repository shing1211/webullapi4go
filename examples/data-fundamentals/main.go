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

// Command data-fundamentals fetches fundamental data for AAPL over the
// Market Data HTTP API.
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
//	go run ./examples/data-fundamentals
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
	cat := data.StockCategoryUS

	profile, err := market.GetCompanyProfile(ctx, "AAPL", cat)
	if err != nil {
		log.Fatal("GetCompanyProfile: ", err)
	}
	printJSON("company profile", profile)

	target, err := market.GetAnalystTargetPrice(ctx, "AAPL", cat)
	if err != nil {
		log.Fatal("GetAnalystTargetPrice: ", err)
	}
	printJSON("analyst target price", target)

	rating, err := market.GetAnalystRating(ctx, "AAPL", cat)
	if err != nil {
		log.Fatal("GetAnalystRating: ", err)
	}
	printJSON("analyst rating", rating)

	capitalFlow, err := market.GetCapitalFlow(ctx, "AAPL", cat, 5)
	if err != nil {
		log.Fatal("GetCapitalFlow: ", err)
	}
	printJSON("capital flow", capitalFlow)

	industry, err := market.GetIndustryComparison(ctx, "AAPL", cat, "EPS_TTM")
	if err != nil {
		log.Fatal("GetIndustryComparison: ", err)
	}
	printJSON("industry comparison", industry)

	earnings, err := market.GetEarningsCalendar(ctx, "AAPL", cat)
	if err != nil {
		log.Fatal("GetEarningsCalendar: ", err)
	}
	printJSON("earnings calendar", earnings)

	dividends, err := market.GetDividendCalendar(ctx, "AAPL", cat)
	if err != nil {
		log.Fatal("GetDividendCalendar: ", err)
	}
	printJSON("dividend calendar", dividends)

	filings, err := market.GetFilings(ctx, "AAPL", cat)
	if err != nil {
		log.Fatal("GetFilings: ", err)
	}
	printJSON("SEC filings", filings)

	income, err := market.GetIncomeStatement(ctx, "AAPL", cat)
	if err != nil {
		log.Fatal("GetIncomeStatement: ", err)
	}
	printJSON("income statement", income)

	balance, err := market.GetBalanceSheet(ctx, "AAPL", cat)
	if err != nil {
		log.Fatal("GetBalanceSheet: ", err)
	}
	printJSON("balance sheet", balance)

	cashFlow, err := market.GetCashFlow(ctx, "AAPL", cat)
	if err != nil {
		log.Fatal("GetCashFlow: ", err)
	}
	printJSON("cash flow", cashFlow)

	indicators, err := market.GetFinancialIndicators(ctx, "AAPL", cat)
	if err != nil {
		log.Fatal("GetFinancialIndicators: ", err)
	}
	printJSON("financial indicators", indicators)

	alert, err := market.GetFinancialAlert(ctx, "AAPL", cat)
	if err != nil {
		log.Fatal("GetFinancialAlert: ", err)
	}
	printJSON("financial alert", alert)

	forecast, err := market.GetForecastEPS(ctx, "AAPL", cat)
	if err != nil {
		log.Fatal("GetForecastEPS: ", err)
	}
	printJSON("forecast EPS", forecast)
}

func printJSON(label string, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("%s: error marshaling: %v", label, err)
		return
	}
	log.Printf("%s:\n%s", label, b)
}
