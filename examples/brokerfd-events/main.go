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

// Command brokerfd-events connects to the Webull Broker FD gRPC event stream,
// subscribes to events for one or more trading accounts, and prints each
// received data event until interrupted with Ctrl+C.
//
// Credentials are read from the environment:
//
//	WEBULL_APP_KEY=...
//	WEBULL_APP_SECRET=...
//	WEBULL_ENVIRONMENT=sandbox
//	WEBULL_TRADE_ACCOUNT_ID=...   (optional; omit to subscribe without account filter)
//
// Run it with:
//
//	go run ./examples/brokerfd-events
package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/shing1211/webullapi4go/brokerfd/events"
	"github.com/shing1211/webullapi4go/client"
)

func main() {
	if os.Getenv("WEBULL_APP_KEY") == "" {
		log.Println("credentials not set, skipping brokerfd-events example")
		return
	}

	cl, err := client.New(client.WithEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = cl.Close() }()

	opts := []events.Option{}
	if account := os.Getenv("WEBULL_TRADE_ACCOUNT_ID"); account != "" {
		opts = append(opts, events.WithAccounts([]string{account}))
	}

	ev, err := events.New(cl, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = ev.Close() }()

	ev.OnConnect(func() { log.Println("connected") })
	ev.OnPing(func() { log.Println("ping") })
	ev.OnError(func(err error) { log.Printf("event stream error: %v", err) })
	ev.OnData(func(subscribeType uint32, contentType string, payload []byte) {
		log.Printf("data: type=%d content=%s payload=%s", subscribeType, contentType, hex.EncodeToString(payload))
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	runErr := make(chan error, 1)
	go func() { runErr <- ev.Run(ctx) }()

	select {
	case <-ctx.Done():
		fmt.Println("shutting down")
	case err := <-runErr:
		if err != nil {
			log.Fatal(err)
		}
	}
}
