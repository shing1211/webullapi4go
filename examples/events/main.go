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

// Command events connects to the Webull gRPC trade-event stream, subscribes to
// order events for one trading account, and prints each decoded event until
// interrupted with Ctrl+C.
//
// Unlike the REST and MQTT paths, the event service signs every Subscribe call
// with HMAC-SHA256 over the serialized request, so no access token is required.
// Credentials are read from the environment:
//
//	WEBULL_APP_KEY=...
//	WEBULL_APP_SECRET=...
//	WEBULL_ENVIRONMENT=sandbox
//	WEBULL_TRADE_ACCOUNT_ID=...
//
// Run it with:
//
//	go run ./examples/events
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/events"
)

func main() {
	cl, err := client.New(client.WithEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = cl.Close() }()

	opts := []events.Option{events.WithSubscribeTypes(events.SubscribeOrder)}
	if account := os.Getenv("WEBULL_TRADE_ACCOUNT_ID"); account != "" {
		opts = append(opts, events.WithAccounts([]string{account}))
	}

	ev, err := events.New(cl, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = ev.Close() }()

	ev.OnConnect(func() { log.Println("subscribed") })
	ev.OnError(func(err error) { log.Printf("event stream error: %v", err) })
	ev.OnOrder(func(o *events.OrderEvent) {
		log.Printf("order %s %s %s status=%s scene=%s symbol=%s qty=%s filled=%s@%s",
			o.OrderID, o.Side, o.OrderType, o.OrderStatus, o.SceneType,
			o.Symbol, o.Quantity, o.FilledQuantity, o.FilledPrice)
	})

	// Ctrl+C or SIGTERM cancels the context and ends the stream cleanly.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	runErr := make(chan error, 1)
	go func() { runErr <- ev.Run(ctx) }()

	select {
	case <-ctx.Done():
		log.Println("shutting down")
	case err := <-runErr:
		if err != nil {
			log.Fatal(err)
		}
	}
}
