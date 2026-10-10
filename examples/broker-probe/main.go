package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shing1211/webullapi4go/broker"
	"github.com/shing1211/webullapi4go/client"
)

func main() {
	ctx := context.Background()
	cl, err := client.New(client.WithEnv())
	if err != nil {
		log.Fatalf("client.New: %v", err)
	}
	if _, err := cl.EnsureToken(ctx); err != nil {
		log.Fatalf("EnsureToken: %v", err)
	}
	bk := broker.New(cl)

	fmt.Println("=== Broker HK Probe ===")
	fmt.Println()

	accounts, err := bk.ListVirtualAccounts(ctx)
	if err != nil {
		fmt.Printf("ListVirtualAccounts: FAIL %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("ListVirtualAccounts: PASS count=%d\n", len(accounts))
	if len(accounts) == 0 {
		fmt.Println("No accounts found; cannot proceed with account-dependent tests")
		os.Exit(1)
	}
	accountID := accounts[0].AccountID
	fmt.Printf("  -> using account_id=%s\n", accountID)

	fmt.Println()

	{
		va, err := bk.GetVirtualAccount(ctx, accountID)
		if err != nil {
			fmt.Printf("GetVirtualAccount: FAIL %v\n", err)
		} else {
			fmt.Printf("GetVirtualAccount: PASS account_id=%s number=%s status=%s\n", va.AccountID, va.AccountNumber, va.AccountStatus)
		}
	}

	{
		bal, err := bk.GetBalance(ctx, accountID)
		if err != nil {
			fmt.Printf("GetBalance: FAIL %v\n", err)
		} else {
			fmt.Printf("GetBalance: PASS total_cash=%s total_market=%s currency=%s\n", bal.TotalCashBalance, bal.TotalMarketValue, bal.TotalAssetCurrency)
		}
	}

	{
		pos, err := bk.GetPositions(ctx, accountID)
		if err != nil {
			fmt.Printf("GetPositions: FAIL %v\n", err)
		} else {
			fmt.Printf("GetPositions: PASS count=%d\n", len(pos))
		}
	}

	{
		acts, err := bk.GetCashActivities(ctx, accountID)
		if err != nil {
			fmt.Printf("GetCashActivities: FAIL %v\n", err)
		} else {
			fmt.Printf("GetCashActivities: PASS count=%d\n", len(acts))
		}
	}

	{
		instrs, err := bk.GetStockInstruments(ctx, []string{"AAPL"})
		if err != nil {
			fmt.Printf("GetStockInstruments: FAIL %v\n", err)
		} else {
			fmt.Printf("GetStockInstruments: PASS count=%d\n", len(instrs))
		}
	}

	{
		locates, err := bk.GetStockLocate(ctx, "AAPL")
		if err != nil {
			fmt.Printf("GetStockLocate: FAIL %v\n", err)
		} else {
			fmt.Printf("GetStockLocate: PASS count=%d\n", len(locates))
		}
	}

	{
		ca, err := bk.GetCorporateActionsDetail(ctx, "AAPL")
		if err != nil {
			fmt.Printf("GetCorporateActionsDetail: FAIL %v\n", err)
		} else {
			fmt.Printf("GetCorporateActionsDetail: PASS count=%d\n", len(ca))
		}
	}

	{
		cats, err := bk.GetEventContractCategories(ctx)
		if err != nil {
			fmt.Printf("GetEventContractCategories: FAIL %v\n", err)
		} else {
			fmt.Printf("GetEventContractCategories: PASS count=%d\n", len(cats))
			if len(cats) > 0 {
				catID := cats[0].CategoryID
				fmt.Printf("  -> using category_id=%s\n", catID)

				series, err := bk.GetEventContractSeries(ctx, catID)
				if err != nil {
					fmt.Printf("GetEventContractSeries: FAIL %v\n", err)
				} else {
					fmt.Printf("GetEventContractSeries: PASS count=%d\n", len(series))
					if len(series) > 0 {
						seriesID := series[0].SeriesID
						fmt.Printf("  -> using series_id=%s\n", seriesID)

						events, err := bk.GetEventContractEvents(ctx, seriesID)
						if err != nil {
							fmt.Printf("GetEventContractEvents: FAIL %v\n", err)
						} else {
							fmt.Printf("GetEventContractEvents: PASS count=%d\n", len(events))
							if len(events) > 0 {
								eventID := events[0].EventID
								fmt.Printf("  -> using event_id=%s\n", eventID)

								instrs, err := bk.GetEventContractInstruments(ctx, eventID)
								if err != nil {
									fmt.Printf("GetEventContractInstruments: FAIL %v\n", err)
								} else {
									fmt.Printf("GetEventContractInstruments: PASS count=%d\n", len(instrs))
								}
							} else {
								fmt.Println("GetEventContractInstruments: SKIP (no events returned)")
							}
						}
					} else {
						fmt.Println("GetEventContractEvents: SKIP (no series returned)")
					}
				}
			} else {
				fmt.Println("GetEventContractSeries: SKIP (no categories returned)")
			}
		}
	}

	{
		fx, err := bk.GetFXRate(ctx, "USD", "HKD")
		if err != nil {
			fmt.Printf("GetFXRate: FAIL %v\n", err)
		} else {
			fmt.Printf("GetFXRate: PASS fx_rate=%s from=%s to=%s\n", fx.FXRate, fx.FromCurrency, fx.ToCurrency)
		}
	}

	{
		fmt.Println("GetFXExchangeDetail: SKIP (no exchange_id available)")
	}

	{
		fmt.Println("GetInstantExchangeDetail: SKIP (no exchange_id available)")
	}

	{
		fmt.Println("GetInstantFundingDetail: SKIP (no funding_id available)")
	}

	{
		preview, err := bk.PreviewOrder(ctx, broker.PreviewOrderRequest{
			AccountID:   accountID,
			Symbol:      "AAPL",
			OrderType:   "LO",
			Side:        "BUY",
			Quantity:    "1",
			LimitPrice:  "150.00",
			TimeInForce: "DAY",
		})
		if err != nil {
			fmt.Printf("PreviewOrder: FAIL %v\n", err)
		} else {
			fmt.Printf("PreviewOrder: PASS estimated_cost=%s estimated_fee=%s\n", preview.EstimatedCost, preview.EstimatedFee)
		}
	}

	{
		fmt.Println("GetOrderDetail: SKIP (no order_id available)")
	}

	{
		hist, err := bk.GetOrderHistory(ctx, accountID)
		if err != nil {
			fmt.Printf("GetOrderHistory: FAIL %v\n", err)
		} else {
			fmt.Printf("GetOrderHistory: PASS count=%d\n", len(hist))
		}
	}

	{
		open, err := bk.GetOpenOrders(ctx, accountID)
		if err != nil {
			fmt.Printf("GetOpenOrders: FAIL %v\n", err)
		} else {
			fmt.Printf("GetOpenOrders: PASS count=%d\n", len(open))
		}
	}

	_ = bk.Close()
}
