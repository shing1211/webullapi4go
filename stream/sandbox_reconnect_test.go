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
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/stream"
)

// sandboxSubscribePath mirrors the package's streaming subscribe endpoint. The
// sandbox test cannot reference the unexported constant, so it repeats it.
const sandboxSubscribePath = "/market-data/streaming/subscribe"

// countingTransport counts POSTs to the streaming subscribe endpoint so the
// sandbox test can observe re-subscription over the real API without a proxy.
type countingTransport struct {
	base http.RoundTripper
	mu   sync.Mutex
	subs int
}

func (c *countingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method == http.MethodPost && req.URL.Path == sandboxSubscribePath {
		c.mu.Lock()
		c.subs++
		c.mu.Unlock()
	}
	return c.base.RoundTrip(req)
}

func (c *countingTransport) subscribeCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.subs
}

// TestSandboxReconnectResubscribe connects over WebSocket, subscribes to AAPL,
// forces a client-side disconnect, and asserts the client reconnects and
// re-issues the subscription automatically.
//
// It runs only when WEBULL_SANDBOX=1, WEBULL_MQTT_WEBSOCKET=1, and
// WEBULL_APP_KEY / WEBULL_APP_SECRET are set; otherwise it skips. Never commit
// the sandbox credentials. Market pushes are not required: the assertion is on
// the reconnect and the HTTP subscribe call count, not on market data.
func TestSandboxReconnectResubscribe(t *testing.T) {
	if os.Getenv("WEBULL_SANDBOX") != "1" {
		t.Skip("set WEBULL_SANDBOX=1, WEBULL_APP_KEY and WEBULL_APP_SECRET to run the sandbox test")
	}
	if os.Getenv("WEBULL_MQTT_WEBSOCKET") != "1" {
		t.Skip("set WEBULL_MQTT_WEBSOCKET=1 to run the sandbox reconnect test over WebSocket")
	}
	appKey := os.Getenv("WEBULL_APP_KEY")
	appSecret := os.Getenv("WEBULL_APP_SECRET")
	if appKey == "" || appSecret == "" {
		t.Skip("set WEBULL_APP_KEY and WEBULL_APP_SECRET to run the sandbox test")
	}

	ct := &countingTransport{base: http.DefaultTransport}
	opts := []client.Option{
		client.WithAppKey(appKey),
		client.WithAppSecret(appSecret),
		client.WithRegion(client.HK),
		client.WithEnvironment(client.Sandbox),
		client.WithHTTPClient(&http.Client{Transport: ct}),
	}
	if baseURL := os.Getenv("WEBULL_BASE_URL"); baseURL != "" {
		opts = append(opts, client.WithBaseURL(baseURL))
	}
	core, err := client.New(opts...)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = core.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// A single token acquisition keeps the sandbox token rate limit (10/30s)
	// satisfied and avoids rapid repeats.
	if _, err := core.EnsureToken(ctx); err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}

	s, err := stream.New(core, stream.WithWebSocket(true), stream.WithAutoReconnect(true))
	if err != nil {
		t.Fatalf("stream.New() error = %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })

	var connects, disconnects int32
	reconnected := make(chan struct{}, 4)
	s.OnConnect(func() {
		if atomic.AddInt32(&connects, 1) > 1 {
			select {
			case reconnected <- struct{}{}:
			default:
			}
		}
	})
	s.OnDisconnect(func(err error) {
		atomic.AddInt32(&disconnects, 1)
		if err != nil {
			t.Logf("connection lost: %v", err)
		}
	})

	if err := s.Connect(ctx); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	t.Logf("connected: session_id=%s", s.SessionID())

	if err := s.Subscribe(ctx, stream.SubscribeRequest{
		Symbols:  []string{"AAPL"},
		Category: stream.CategoryUSStock,
		SubTypes: []stream.SubType{stream.SubTypeQuote, stream.SubTypeSnapshot, stream.SubTypeTick},
		Grab:     true,
	}); err != nil {
		t.Fatalf("Subscribe(AAPL) error = %v", err)
	}
	initialSubs := ct.subscribeCount()
	t.Logf("subscribed to AAPL (subscribe calls=%d)", initialSubs)

	// Force a client-side disconnect: opening another connection with the same
	// session id makes the server drop the first one, per Webull's connection
	// rules. The first client then auto-reconnects and re-subscribes.
	kick, err := stream.New(core, stream.WithWebSocket(true), stream.WithSessionID(s.SessionID()))
	if err != nil {
		t.Fatalf("stream.New(kick) error = %v", err)
	}
	t.Cleanup(func() { _ = kick.Close() })
	if err := kick.Connect(ctx); err != nil {
		t.Fatalf("kick Connect() error = %v", err)
	}
	t.Log("opened a second connection with the same session id to force a drop")

	select {
	case <-reconnected:
		t.Log("client reconnected")
	case <-time.After(60 * time.Second):
		t.Fatalf("client did not reconnect within 60s (OnConnect=%d, OnDisconnect=%d)",
			atomic.LoadInt32(&connects), atomic.LoadInt32(&disconnects))
	}

	// Re-subscription precedes the OnConnect callback, but allow a short window
	// in case the callback was observed slightly ahead of the HTTP completion.
	deadline := time.Now().Add(15 * time.Second)
	for ct.subscribeCount() <= initialSubs && time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
	}
	if got := ct.subscribeCount(); got <= initialSubs {
		t.Fatalf("no re-subscription after reconnect (subscribe calls=%d, initial=%d)", got, initialSubs)
	}
	t.Logf("re-subscribed after reconnect (subscribe calls=%d, initial=%d)", ct.subscribeCount(), initialSubs)
}
