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

package stream_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	marketdatav1 "github.com/shing1211/webullapi4go/gen/webull/marketdata/v1"
	"github.com/shing1211/webullapi4go/stream"
)

// TestSandboxStreaming connects to the shared Webull sandbox MQTT broker,
// subscribes to AAPL quote/snapshot/tick over HTTP, waits briefly for a push,
// then unsubscribes and closes.
//
// It runs only when WEBULL_SANDBOX=1 and WEBULL_APP_KEY / WEBULL_APP_SECRET are
// set; otherwise it skips so the default test run stays hermetic. Never commit
// the sandbox credentials. A market that is closed may push nothing, which is
// tolerated; connect and subscribe failures are not.
func TestSandboxStreaming(t *testing.T) {
	if os.Getenv("WEBULL_SANDBOX") != "1" {
		t.Skip("set WEBULL_SANDBOX=1, WEBULL_APP_KEY and WEBULL_APP_SECRET to run the sandbox test")
	}
	appKey := os.Getenv("WEBULL_APP_KEY")
	appSecret := os.Getenv("WEBULL_APP_SECRET")
	if appKey == "" || appSecret == "" {
		t.Skip("set WEBULL_APP_KEY and WEBULL_APP_SECRET to run the sandbox test")
	}

	opts := []client.Option{
		client.WithAppKey(appKey),
		client.WithAppSecret(appSecret),
		client.WithRegion(client.HK),
		client.WithEnvironment(client.Sandbox),
	}
	if baseURL := os.Getenv("WEBULL_BASE_URL"); baseURL != "" {
		opts = append(opts, client.WithBaseURL(baseURL))
	}

	core, err := client.New(opts...)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = core.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	if _, err := core.EnsureToken(ctx); err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}

	s, err := stream.New(core, sandboxStreamOptions()...)
	if err != nil {
		t.Fatalf("stream.New() error = %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })

	received := make(chan string, 32)
	streamErrors := make(chan error, 8)
	s.OnQuote(func(*marketdatav1.Quote) { enqueue(received, "quote") })
	s.OnSnapshot(func(*marketdatav1.Snapshot) { enqueue(received, "snapshot") })
	s.OnTick(func(*marketdatav1.Tick) { enqueue(received, "tick") })
	s.OnNotice(func([]byte) { enqueue(received, "notice") })
	s.OnError(func(err error) { enqueueErr(streamErrors, err) })

	if err := s.Connect(ctx); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	t.Logf("connected: session_id=%s", s.SessionID())

	sub := stream.SubscribeRequest{
		Symbols:  []string{"AAPL"},
		Category: stream.CategoryUSStock,
		SubTypes: []stream.SubType{stream.SubTypeQuote, stream.SubTypeSnapshot, stream.SubTypeTick},
		Grab:     true,
	}
	if err := s.Subscribe(ctx, sub); err != nil {
		t.Fatalf("Subscribe(AAPL, US_STOCK) error = %v", err)
	}
	t.Log("subscribe acknowledged by the HTTP API")

	select {
	case topic := <-received:
		t.Logf("received %s push for AAPL", topic)
	case err := <-streamErrors:
		t.Fatalf("stream error after subscribe: %v", err)
	case <-time.After(20 * time.Second):
		t.Log("no push within 20s (market may be closed); subscribe ack was clean")
	}

	unsub := stream.UnsubscribeRequest{
		Symbols:        []string{"AAPL"},
		Category:       stream.CategoryUSStock,
		SubTypes:       []stream.SubType{stream.SubTypeQuote, stream.SubTypeSnapshot, stream.SubTypeTick},
		UnsubscribeAll: true,
	}
	if err := s.Unsubscribe(ctx, unsub); err != nil {
		t.Fatalf("Unsubscribe() error = %v", err)
	}
	t.Log("unsubscribe acknowledged; closing")
}

// sandboxStreamOptions lets a run target the WebSocket transport or an
// alternate broker without changing the test. WEBULL_MQTT_WEBSOCKET=1 selects
// the wss endpoint; WEBULL_MQTT_URL overrides the broker address entirely.
func sandboxStreamOptions() []stream.Option {
	var opts []stream.Option
	if os.Getenv("WEBULL_MQTT_WEBSOCKET") == "1" {
		opts = append(opts, stream.WithWebSocket(true))
	}
	if rawURL := os.Getenv("WEBULL_MQTT_URL"); rawURL != "" {
		opts = append(opts, stream.WithMQTTURL(rawURL))
	}
	return opts
}

// enqueue performs a non-blocking send, dropping the value when the buffer is
// full so a slow test can never block the MQTT callback goroutine.
func enqueue[T any](ch chan T, v T) {
	select {
	case ch <- v:
	default:
	}
}

// enqueueErr performs a non-blocking error send.
func enqueueErr(ch chan error, err error) {
	select {
	case ch <- err:
	default:
	}
}
