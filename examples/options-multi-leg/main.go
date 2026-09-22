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

// Command options-multi-leg probes every multi-leg option strategy value via
// PreviewOrder. It is read-only; no live orders are placed.
//
// Run it with:
//
//	WEBULL_OPTIONS_TEST=1
//	WEBULL_APP_KEY=...
//	WEBULL_APP_SECRET=...
//	WEBULL_ENVIRONMENT=sandbox
//	WEBULL_ACCOUNT_ID=...  # optional; uses first account if unset
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/trade"
)

func main() {
	if os.Getenv("WEBULL_OPTIONS_TEST") != "1" {
		fmt.Println("options multi-leg probe disabled (set WEBULL_OPTIONS_TEST=1 to run)")
		os.Exit(0)
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

	trading := trade.New(cl)
	defer func() { _ = trading.Close() }()

	accountID, err := resolveAccountID(ctx, trading)
	if err != nil {
		log.Fatal(err)
	}

	strategies := []struct {
		name       string
		strategy   trade.OptionStrategy
		legs       []trade.OrderLeg
		quantity   string
		limitPrice string
	}{
		{
			name:       "NORMAL",
			strategy:   trade.OptionStrategySingle,
			legs:       []trade.OrderLeg{aaplCallLeg("220.00", "2026-12-18")},
			quantity:   "1",
			limitPrice: "11.25",
		},
		{
			name:     "VERTICAL",
			strategy: trade.OptionStrategyVertical,
			legs: []trade.OrderLeg{
				{Symbol: "AAPL", Side: trade.OrderSideBuy, OptionType: trade.OptionTypeCall, StrikePrice: "220.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
				{Symbol: "AAPL", Side: trade.OrderSideSell, OptionType: trade.OptionTypeCall, StrikePrice: "230.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
			},
			quantity:   "1",
			limitPrice: "1.50",
		},
		{
			name:     "STRADDLE",
			strategy: trade.OptionStrategyStraddle,
			legs: []trade.OrderLeg{
				{Symbol: "AAPL", Side: trade.OrderSideBuy, OptionType: trade.OptionTypeCall, StrikePrice: "220.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
				{Symbol: "AAPL", Side: trade.OrderSideBuy, OptionType: trade.OptionTypePut, StrikePrice: "220.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
			},
			quantity:   "1",
			limitPrice: "3.00",
		},
		{
			name:     "STRANGLE",
			strategy: trade.OptionStrategyStrangle,
			legs: []trade.OrderLeg{
				{Symbol: "AAPL", Side: trade.OrderSideBuy, OptionType: trade.OptionTypeCall, StrikePrice: "220.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
				{Symbol: "AAPL", Side: trade.OrderSideBuy, OptionType: trade.OptionTypePut, StrikePrice: "210.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
			},
			quantity:   "1",
			limitPrice: "2.50",
		},
		{
			name:     "IRON_CONDOR",
			strategy: trade.OptionStrategyIronCondor,
			legs: []trade.OrderLeg{
				{Symbol: "AAPL", Side: trade.OrderSideBuy, OptionType: trade.OptionTypePut, StrikePrice: "200.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
				{Symbol: "AAPL", Side: trade.OrderSideSell, OptionType: trade.OptionTypePut, StrikePrice: "210.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
				{Symbol: "AAPL", Side: trade.OrderSideSell, OptionType: trade.OptionTypeCall, StrikePrice: "230.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
				{Symbol: "AAPL", Side: trade.OrderSideBuy, OptionType: trade.OptionTypeCall, StrikePrice: "240.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
			},
			quantity:   "1",
			limitPrice: "1.00",
		},
		{
			name:     "IRON_BUTTERFLY",
			strategy: trade.OptionStrategyIronButterfly,
			legs: []trade.OrderLeg{
				{Symbol: "AAPL", Side: trade.OrderSideBuy, OptionType: trade.OptionTypePut, StrikePrice: "200.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
				{Symbol: "AAPL", Side: trade.OrderSideSell, OptionType: trade.OptionTypePut, StrikePrice: "220.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
				{Symbol: "AAPL", Side: trade.OrderSideSell, OptionType: trade.OptionTypeCall, StrikePrice: "220.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
				{Symbol: "AAPL", Side: trade.OrderSideBuy, OptionType: trade.OptionTypeCall, StrikePrice: "240.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
			},
			quantity:   "1",
			limitPrice: "1.50",
		},
		{
			name:     "BUTTERFLY",
			strategy: trade.OptionStrategyButterfly,
			legs: []trade.OrderLeg{
				{Symbol: "AAPL", Side: trade.OrderSideBuy, OptionType: trade.OptionTypeCall, StrikePrice: "210.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
				{Symbol: "AAPL", Side: trade.OrderSideSell, OptionType: trade.OptionTypeCall, StrikePrice: "220.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "2"},
				{Symbol: "AAPL", Side: trade.OrderSideBuy, OptionType: trade.OptionTypeCall, StrikePrice: "230.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
			},
			quantity:   "1",
			limitPrice: "2.00",
		},
		{
			name:     "COLLAR",
			strategy: trade.OptionStrategyCollar,
			legs: []trade.OrderLeg{
				{Symbol: "AAPL", Side: trade.OrderSideBuy, OptionType: trade.OptionTypePut, StrikePrice: "210.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
				{Symbol: "AAPL", Side: trade.OrderSideSell, OptionType: trade.OptionTypeCall, StrikePrice: "230.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
			},
			quantity:   "1",
			limitPrice: "1.00",
		},
		{
			name:     "CALENDAR",
			strategy: trade.OptionStrategyCalendar,
			legs: []trade.OrderLeg{
				{Symbol: "AAPL", Side: trade.OrderSideBuy, OptionType: trade.OptionTypeCall, StrikePrice: "220.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
				{Symbol: "AAPL", Side: trade.OrderSideSell, OptionType: trade.OptionTypeCall, StrikePrice: "220.00", OptionExpireDate: "2027-01-15", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
			},
			quantity:   "1",
			limitPrice: "1.50",
		},
		{
			name:     "DIAGONAL",
			strategy: trade.OptionStrategyDiagonal,
			legs: []trade.OrderLeg{
				{Symbol: "AAPL", Side: trade.OrderSideBuy, OptionType: trade.OptionTypeCall, StrikePrice: "210.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
				{Symbol: "AAPL", Side: trade.OrderSideSell, OptionType: trade.OptionTypeCall, StrikePrice: "230.00", OptionExpireDate: "2027-01-15", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
			},
			quantity:   "1",
			limitPrice: "1.00",
		},
		{
			name:     "RATIO",
			strategy: trade.OptionStrategyRatio,
			legs: []trade.OrderLeg{
				{Symbol: "AAPL", Side: trade.OrderSideBuy, OptionType: trade.OptionTypeCall, StrikePrice: "210.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "2"},
				{Symbol: "AAPL", Side: trade.OrderSideSell, OptionType: trade.OptionTypeCall, StrikePrice: "230.00", OptionExpireDate: "2026-12-18", InstrumentType: trade.InstrumentTypeOption, Market: trade.MarketUS, Quantity: "1"},
			},
			quantity:   "1",
			limitPrice: "2.00",
		},
	}

	fmt.Printf("Probing %d option strategies against account %s\n\n", len(strategies), accountID)

	for _, tc := range strategies {
		req := trade.PlaceOrderRequest{
			AccountID: accountID,
			NewOrders: []trade.OrderRequest{
				{
					ClientOrderID:  fmt.Sprintf("probe-%d-%s", time.Now().Unix(), tc.strategy),
					ComboType:      trade.ComboTypeNormal,
					InstrumentType: trade.InstrumentTypeOption,
					Market:         trade.MarketUS,
					Symbol:         "AAPL",
					OrderType:      trade.OrderTypeLimit,
					Side:           trade.OrderSideBuy,
					Quantity:       tc.quantity,
					EntrustType:    trade.EntrustTypeQty,
					TimeInForce:    trade.TimeInForceDay,
					LimitPrice:     tc.limitPrice,
					OptionStrategy: tc.strategy,
					Legs:           tc.legs,
				},
			},
		}

		preview, err := trading.PreviewOrder(ctx, req)
		if err != nil {
			fmt.Printf("[%s] FAIL: %v\n", tc.name, err)
			continue
		}
		fmt.Printf("[%s] PASS: estimated_cost=%s estimated_fee=%s\n", tc.name, preview.EstimatedCost, preview.EstimatedTransactionFee)
	}
}

func resolveAccountID(ctx context.Context, trading *trade.Client) (string, error) {
	if id := os.Getenv("WEBULL_TRADE_ACCOUNT_ID"); id != "" {
		return id, nil
	}
	accounts, err := trading.ListAccounts(ctx)
	if err != nil {
		return "", err
	}
	if len(accounts) == 0 {
		return "", errors.New("no trading accounts available")
	}
	return accounts[0].AccountID, nil
}

func aaplCallLeg(strike, expiry string) trade.OrderLeg {
	return trade.OrderLeg{
		InstrumentType:   trade.InstrumentTypeOption,
		Market:           trade.MarketUS,
		Symbol:           "AAPL",
		Side:             trade.OrderSideBuy,
		StrikePrice:      strike,
		OptionExpireDate: expiry,
		OptionType:       trade.OptionTypeCall,
		Quantity:         "1",
	}
}
