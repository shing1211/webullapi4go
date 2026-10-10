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

// Command account-monitor periodically prints an account's balance, buying
// power, and initial margin. It performs one read-only balance request at a
// time, bounds every request with a timeout, and stops after three consecutive
// failures. Ctrl+C or SIGTERM cancels the monitor and its active request.
//
// Credentials are read from the environment:
//
//	WEBULL_APP_KEY=...
//	WEBULL_APP_SECRET=...
//	WEBULL_ENVIRONMENT=sandbox
//	WEBULL_ACCOUNT_ID=...
//	WEBULL_MONITOR_INTERVAL=30s
//
// Run it with:
//
//	go run ./examples/account-monitor
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/trade"
)

const (
	defaultMonitorInterval = 30 * time.Second
	minimumMonitorInterval = time.Second
	monitorTokenTimeout    = 5 * time.Minute
	monitorRequestTimeout  = 15 * time.Second
	maximumPollErrors      = 3
)

type pollResult struct {
	balance *trade.AssetsBalance
	err     error
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if os.Getenv("WEBULL_APP_KEY") == "" {
		fmt.Println("example: WEBULL_APP_KEY not set, skipping")
		return nil
	}

	interval, err := monitorInterval()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cl, err := client.New(client.WithEnv())
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}
	defer func() { _ = cl.Close() }()

	tokenCtx, cancelToken := context.WithTimeout(ctx, monitorTokenTimeout)
	_, err = cl.EnsureToken(tokenCtx)
	cancelToken()
	if err != nil {
		return fmt.Errorf("ensure access token: %w", err)
	}

	trading := trade.New(cl)
	defer func() { _ = trading.Close() }()

	accountCtx, cancelAccount := context.WithTimeout(ctx, monitorRequestTimeout)
	accountID, err := resolveAccountID(accountCtx, trading)
	cancelAccount()
	if err != nil {
		return fmt.Errorf("resolve account: %w", err)
	}

	log.Printf("monitoring account %s every %s; press Ctrl+C to stop", accountID, interval)
	return monitorBalance(ctx, trading, accountID, interval)
}

func monitorInterval() (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv("WEBULL_MONITOR_INTERVAL"))
	if value == "" {
		return defaultMonitorInterval, nil
	}
	interval, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("parse WEBULL_MONITOR_INTERVAL: %w", err)
	}
	if interval < minimumMonitorInterval {
		return 0, fmt.Errorf("WEBULL_MONITOR_INTERVAL must be at least %s", minimumMonitorInterval)
	}
	return interval, nil
}

func resolveAccountID(ctx context.Context, trading *trade.Client) (string, error) {
	if id := os.Getenv("WEBULL_ACCOUNT_ID"); id != "" {
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

func monitorBalance(ctx context.Context, trading *trade.Client, accountID string, interval time.Duration) error {
	workerCtx, stopWorker := context.WithCancel(ctx)
	requests := make(chan struct{})
	results := make(chan pollResult)
	var worker sync.WaitGroup
	worker.Add(1)
	go func() {
		defer worker.Done()
		for {
			select {
			case <-workerCtx.Done():
				return
			case <-requests:
				result := fetchBalance(workerCtx, trading, accountID)
				select {
				case results <- result:
				case <-workerCtx.Done():
					return
				}
			}
		}
	}()
	defer func() {
		stopWorker()
		worker.Wait()
	}()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	inFlight := false
	startPoll := func() bool {
		if inFlight {
			return false
		}
		select {
		case requests <- struct{}{}:
			inFlight = true
			return true
		case <-ctx.Done():
			return false
		}
	}
	if !startPoll() {
		return nil
	}

	consecutiveErrors := 0
	for {
		select {
		case <-ctx.Done():
			log.Println("account monitor stopped")
			return nil
		case <-ticker.C:
			startPoll()
		case result := <-results:
			inFlight = false
			if result.err != nil {
				if ctx.Err() != nil {
					return nil
				}
				consecutiveErrors++
				log.Printf("balance request failed (%d/%d): %v", consecutiveErrors, maximumPollErrors, result.err)
				if consecutiveErrors >= maximumPollErrors {
					return fmt.Errorf("account monitor stopped after %d consecutive failures: %w", maximumPollErrors, result.err)
				}
				continue
			}
			consecutiveErrors = 0
			logBalance(accountID, result.balance)
		}
	}
}

func fetchBalance(ctx context.Context, trading *trade.Client, accountID string) pollResult {
	requestCtx, cancel := context.WithTimeout(ctx, monitorRequestTimeout)
	defer cancel()
	balance, err := trading.GetBalance(requestCtx, accountID)
	return pollResult{balance: balance, err: err}
}

func logBalance(accountID string, balance *trade.AssetsBalance) {
	log.Printf("balance account=%s currency=%s cash=%s market_value=%s pnl=%s initial_margin=%s",
		accountID,
		balance.TotalAssetCurrency,
		balance.TotalCashBalance.String(),
		balance.TotalMarketValue.String(),
		balance.TotalUnrealizedProfitLoss.String(),
		balance.InitMargin.String(),
	)
	for _, asset := range balance.AccountCurrencyAssets {
		log.Printf("margin account=%s currency=%s buying_power=%s initial_margin=%s held=%s frozen=%s unpaid_interest=%s",
			accountID,
			asset.Currency,
			asset.BuyingPower.String(),
			asset.InitMargin.String(),
			asset.HeldAmount.String(),
			asset.FrozenAmount.String(),
			asset.InterestsUnpaid.String(),
		)
	}
}
