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

// Command streaming connects to the Webull streaming broker over MQTT-over-
// WebSocket, subscribes to AAPL quote/snapshot/tick pushes, and prints messages
// until interrupted with Ctrl+C.
//
// MQTT-over-WebSocket is used because plain MQTT on port 1883 is blocked on some
// networks. Credentials are read from the environment:
//
//	WEBULL_APP_KEY=...
//	WEBULL_APP_SECRET=...
//	WEBULL_ENVIRONMENT=sandbox
//
// Run it with:
//
//	go run ./examples/streaming
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/shing1211/webullapi4go/client"
	marketdatav1 "github.com/shing1211/webullapi4go/gen/webull/marketdata/v1"
	"github.com/shing1211/webullapi4go/stream"
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

	// Subscribe and unsubscribe are HTTP calls, so a token must be available
	// before streaming starts.
	if _, err := cl.EnsureToken(context.Background()); err != nil {
		log.Fatal(err)
	}

	// WithWebSocket switches from the plain TCP broker to wss://...:8883/mqtt.
	s, err := stream.New(cl, stream.WithWebSocket(true))
	if err != nil {
		log.Fatal(err)
	}

	// Ctrl+C or SIGTERM cancels the context and ends the stream cleanly.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	s.OnConnect(func() { log.Println("connected") })
	s.OnDisconnect(func(err error) { log.Printf("disconnected: %v", err) })
	s.OnError(func(err error) { log.Printf("stream error: %v", err) })

	s.OnQuote(func(q *marketdatav1.Quote) {
		log.Printf("quote %s asks=%d bids=%d",
			q.GetBasic().GetSymbol(), len(q.GetAsks()), len(q.GetBids()))
	})
	s.OnSnapshot(func(snap *marketdatav1.Snapshot) {
		log.Printf("snapshot %s price=%s", snap.GetBasic().GetSymbol(), snap.GetPrice())
	})
	s.OnTick(func(t *marketdatav1.Tick) {
		log.Printf("tick %s price=%s volume=%s side=%s",
			t.GetBasic().GetSymbol(), t.GetPrice(), t.GetVolume(), t.GetSide())
	})

	if err := s.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer func() { _ = s.Close() }()

	if err := s.Subscribe(ctx, stream.SubscribeRequest{
		Symbols:  []string{"AAPL"},
		Category: stream.CategoryUSStock,
		SubTypes: []stream.SubType{stream.SubTypeQuote, stream.SubTypeSnapshot, stream.SubTypeTick},
		Grab:     true,
	}); err != nil {
		log.Fatal(err)
	}

	log.Println("streaming AAPL; press Ctrl+C to stop")
	<-ctx.Done()
	log.Println("shutting down")
}
