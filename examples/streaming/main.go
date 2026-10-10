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
// WebSocket, subscribes to AAPL quote/snapshot/tick pushes, and consumes bounded
// channels until interrupted with Ctrl+C or the two-minute run limit expires.
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
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shing1211/webullapi4go/client"
	marketdatav1 "github.com/shing1211/webullapi4go/gen/webull/marketdata/v1"
	"github.com/shing1211/webullapi4go/stream"
)

const (
	setupTimeout     = 5 * time.Minute
	requestTimeout   = 20 * time.Second
	streamRunTimeout = 2 * time.Minute
	channelBuffer    = 32
)

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cl, err := client.New(client.WithEnv())
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}
	defer func() { _ = cl.Close() }()

	tokenCtx, cancelToken := context.WithTimeout(ctx, setupTimeout)
	_, err = cl.EnsureToken(tokenCtx)
	cancelToken()
	if err != nil {
		return fmt.Errorf("ensure access token: %w", err)
	}

	s, err := stream.New(cl, stream.WithWebSocket(true))
	if err != nil {
		return fmt.Errorf("create stream client: %w", err)
	}
	defer func() { _ = s.Close() }()

	channelConfig := stream.ChannelConfig{
		Policy:     stream.DropOldest,
		BufferSize: channelBuffer,
	}
	quotes, cancelQuotes := s.SubscribeQuoteChan(channelConfig)
	defer cancelQuotes()
	snapshots, cancelSnapshots := s.SubscribeSnapshotChan(channelConfig)
	defer cancelSnapshots()
	ticks, cancelTicks := s.SubscribeTickChan(channelConfig)
	defer cancelTicks()

	s.OnConnect(func() { log.Println("connected") })
	s.OnDisconnect(func(err error) { log.Printf("disconnected: %v", err) })
	s.OnReconnecting(func() { log.Println("reconnecting") })
	s.OnError(func(err error) { log.Printf("stream error: %v", err) })

	connectCtx, cancelConnect := context.WithTimeout(ctx, requestTimeout)
	err = s.Connect(connectCtx)
	cancelConnect()
	if err != nil {
		return fmt.Errorf("connect stream: %w", err)
	}

	subscribeCtx, cancelSubscribe := context.WithTimeout(ctx, requestTimeout)
	err = s.Subscribe(subscribeCtx, stream.SubscribeRequest{
		Symbols:  []string{"AAPL"},
		Category: stream.CategoryUSStock,
		SubTypes: []stream.SubType{stream.SubTypeQuote, stream.SubTypeSnapshot, stream.SubTypeTick},
		Grab:     true,
	})
	cancelSubscribe()
	if err != nil {
		return fmt.Errorf("subscribe stream: %w", err)
	}

	runCtx, cancelRun := context.WithTimeout(ctx, streamRunTimeout)
	defer cancelRun()
	log.Printf("streaming AAPL for at most %s; press Ctrl+C to stop", streamRunTimeout)
	consume(runCtx, quotes, snapshots, ticks)
	return nil
}

func consume(
	ctx context.Context,
	quotes <-chan *marketdatav1.Quote,
	snapshots <-chan *marketdatav1.Snapshot,
	ticks <-chan *marketdatav1.Tick,
) {
	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				log.Println("stream run limit reached")
			} else {
				log.Println("shutting down")
			}
			return
		case quote, ok := <-quotes:
			if !ok {
				log.Println("quote channel closed")
				return
			}
			log.Printf("quote %s asks=%d bids=%d",
				quote.GetBasic().GetSymbol(), len(quote.GetAsks()), len(quote.GetBids()))
		case snapshot, ok := <-snapshots:
			if !ok {
				log.Println("snapshot channel closed")
				return
			}
			log.Printf("snapshot %s price=%s", snapshot.GetBasic().GetSymbol(), snapshot.GetPrice())
		case tick, ok := <-ticks:
			if !ok {
				log.Println("tick channel closed")
				return
			}
			log.Printf("tick %s price=%s volume=%s side=%s",
				tick.GetBasic().GetSymbol(), tick.GetPrice(), tick.GetVolume(), tick.GetSide())
		}
	}
}
