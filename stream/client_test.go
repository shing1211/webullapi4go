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

package stream

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/shing1211/webullapi4go/client"
	marketdatav1 "github.com/shing1211/webullapi4go/gen/webull/marketdata/v1"
	mqtt "github.com/shing1211/webullapi4go/pkg/transport/mqtt"
)

// newTestCore returns a core client pointed at a stub HTTP server. It needs
// credentials only for signing, so no sandbox access is required.
func newTestCore(t *testing.T, baseURL string) *client.Client {
	t.Helper()
	core, err := client.New(
		client.WithAppKey("test-key"),
		client.WithAppSecret("test-secret"),
		client.WithBaseURL(baseURL),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = core.Close() })
	return core
}

func TestDispatchQuote(t *testing.T) {
	c := &Client{}
	got := make(chan *marketdatav1.Quote, 1)
	c.OnQuote(func(q *marketdatav1.Quote) { got <- q })

	want := &marketdatav1.Quote{
		Basic: &marketdatav1.Basic{Symbol: "AAPL", InstrumentId: "913256135"},
		Asks:  []*marketdatav1.AskBid{{Price: "190.20", Size: "300"}},
	}
	data, err := proto.Marshal(want)
	if err != nil {
		t.Fatalf("proto.Marshal() error = %v", err)
	}
	if err := c.dispatch(TopicQuote, data); err != nil {
		t.Fatalf("dispatch(quote) error = %v", err)
	}
	select {
	case q := <-got:
		if q.GetBasic().GetSymbol() != "AAPL" || q.GetAsks()[0].GetPrice() != "190.20" {
			t.Fatalf("decoded quote = %v, want AAPL at 190.20", q)
		}
	default:
		t.Fatal("OnQuote handler was not invoked")
	}
}

func TestDispatchSnapshot(t *testing.T) {
	c := &Client{}
	got := make(chan *marketdatav1.Snapshot, 1)
	c.OnSnapshot(func(s *marketdatav1.Snapshot) { got <- s })

	want := &marketdatav1.Snapshot{Basic: &marketdatav1.Basic{Symbol: "AAPL"}, Price: "190.12"}
	data, err := proto.Marshal(want)
	if err != nil {
		t.Fatalf("proto.Marshal() error = %v", err)
	}
	if err := c.dispatch(TopicSnapshot, data); err != nil {
		t.Fatalf("dispatch(snapshot) error = %v", err)
	}
	select {
	case s := <-got:
		if s.GetPrice() != "190.12" {
			t.Fatalf("decoded snapshot price = %q, want 190.12", s.GetPrice())
		}
	default:
		t.Fatal("OnSnapshot handler was not invoked")
	}
}

func TestDispatchTick(t *testing.T) {
	c := &Client{}
	got := make(chan *marketdatav1.Tick, 1)
	c.OnTick(func(tk *marketdatav1.Tick) { got <- tk })

	want := &marketdatav1.Tick{Basic: &marketdatav1.Basic{Symbol: "AAPL"}, Price: "190.12", Side: "B"}
	data, err := proto.Marshal(want)
	if err != nil {
		t.Fatalf("proto.Marshal() error = %v", err)
	}
	if err := c.dispatch(TopicTick, data); err != nil {
		t.Fatalf("dispatch(tick) error = %v", err)
	}
	select {
	case tk := <-got:
		if tk.GetSide() != "B" || tk.GetPrice() != "190.12" {
			t.Fatalf("decoded tick = %v, want side B at 190.12", tk)
		}
	default:
		t.Fatal("OnTick handler was not invoked")
	}
}

func TestDispatchNotice(t *testing.T) {
	c := &Client{}
	got := make(chan []byte, 1)
	c.OnNotice(func(b []byte) { got <- b })

	payload := []byte(`{"type":"status","rtt":100}`)
	if err := c.dispatch(TopicNotice, payload); err != nil {
		t.Fatalf("dispatch(notice) error = %v", err)
	}
	select {
	case b := <-got:
		if string(b) != string(payload) {
			t.Fatalf("notice payload = %s, want %s", b, payload)
		}
	default:
		t.Fatal("OnNotice handler was not invoked")
	}
}

func TestDispatchTopicVariants(t *testing.T) {
	for _, topic := range []string{"quote", "QUOTE", "webull/quote", "  quote  ", "a/b/Quote"} {
		t.Run(topic, func(t *testing.T) {
			c := &Client{}
			called := false
			c.OnQuote(func(*marketdatav1.Quote) { called = true })
			data, err := proto.Marshal(&marketdatav1.Quote{Basic: &marketdatav1.Basic{Symbol: "AAPL"}})
			if err != nil {
				t.Fatalf("proto.Marshal() error = %v", err)
			}
			if err := c.dispatch(topic, data); err != nil {
				t.Fatalf("dispatch(%q) error = %v", topic, err)
			}
			if !called {
				t.Fatalf("dispatch(%q) did not invoke the quote handler", topic)
			}
		})
	}
}

func TestDispatchUnknownTopic(t *testing.T) {
	c := &Client{}
	err := c.dispatch("bogus", []byte("x"))
	if err == nil {
		t.Fatal("dispatch(bogus) = nil, want error")
	}
	if !strings.Contains(err.Error(), "unknown topic") {
		t.Fatalf("error = %v, want unknown topic", err)
	}
}

func TestDispatchInvalidPayload(t *testing.T) {
	c := &Client{}
	if err := c.dispatch(TopicQuote, []byte{0xff, 0xff, 0xff, 0xff}); err == nil {
		t.Fatal("dispatch(quote, invalid) = nil, want decode error")
	}
}

func TestDispatchEmptyTopic(t *testing.T) {
	c := &Client{}
	if err := c.dispatch("", []byte("x")); err == nil {
		t.Fatal("dispatch(empty) = nil, want error")
	}
}

func TestDispatchEchoIgnored(t *testing.T) {
	c := &Client{}
	called := false
	c.OnQuote(func(*marketdatav1.Quote) { called = true })
	if err := c.dispatch(TopicEcho, nil); err != nil {
		t.Fatalf("dispatch(echo) error = %v", err)
	}
	if called {
		t.Fatal("echo dispatch invoked a data handler")
	}
}

func TestHandleMessageReportsError(t *testing.T) {
	c := &Client{}
	got := make(chan error, 1)
	c.OnError(func(err error) { got <- err })
	c.handleMessage(mqtt.Message{Topic: "bogus", Payload: []byte("x")})
	select {
	case err := <-got:
		if !strings.Contains(err.Error(), "unknown topic") {
			t.Fatalf("error = %v, want unknown topic", err)
		}
	default:
		t.Fatal("handleMessage did not report the dispatch error")
	}
}

func TestNewRequiresClient(t *testing.T) {
	if _, err := New(nil); err == nil {
		t.Fatal("New(nil) = nil error, want error")
	}
}

func TestNewSessionIDUnique(t *testing.T) {
	core := newTestCore(t, "http://127.0.0.1:1")
	seen := make(map[string]struct{}, 64)
	for i := 0; i < 64; i++ {
		s, err := New(core)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		id := s.SessionID()
		if id == "" {
			t.Fatal("SessionID() = empty")
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate session id %q", id)
		}
		seen[id] = struct{}{}
		_ = s.Close()
	}
}

func TestNewSessionIDFormat(t *testing.T) {
	id := newSessionID()
	if len(id) != 36 {
		t.Fatalf("newSessionID() = %q, want a 36-character UUID", id)
	}
	if id[14] != '4' {
		t.Fatalf("newSessionID() = %q, want version 4 UUID", id)
	}
}

func TestWithSessionID(t *testing.T) {
	core := newTestCore(t, "http://127.0.0.1:1")
	s, err := New(core, WithSessionID("session-123"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got := s.SessionID(); got != "session-123" {
		t.Fatalf("SessionID() = %q, want session-123", got)
	}
}

func TestWithClientIDAliasesSessionID(t *testing.T) {
	core := newTestCore(t, "http://127.0.0.1:1")
	s, err := New(core, WithClientID("client-9"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got := s.SessionID(); got != "client-9" {
		t.Fatalf("SessionID() = %q, want client-9", got)
	}
}

func TestConfigValidate(t *testing.T) {
	valid := defaultConfig()
	valid.sessionID = "s"
	if err := valid.validate(); err != nil {
		t.Fatalf("valid config error = %v", err)
	}

	cases := map[string]func(*config){
		"empty session":      func(c *config) { c.sessionID = "" },
		"zero keepalive":     func(c *config) { c.keepAlive = 0 },
		"zero connect":       func(c *config) { c.connectTimeout = 0 },
		"zero write":         func(c *config) { c.writeTimeout = 0 },
		"zero channel depth": func(c *config) { c.messageChannelDepth = 0 },
	}
	for name, mutate := range cases {
		c := defaultConfig()
		c.sessionID = "s"
		mutate(&c)
		if err := c.validate(); err == nil {
			t.Errorf("validate() with %s = nil, want error", name)
		}
	}
}

func TestSubscribeRequestBody(t *testing.T) {
	req := SubscribeRequest{
		Symbols:  []string{" AAPL ", "", "TSLA"},
		Category: CategoryUSStock,
		SubTypes: []SubType{SubTypeQuote, "snapshot"},
		Grab:     true,
		Depth:    "10",
	}
	body, err := req.toBody("sess-1")
	if err != nil {
		t.Fatalf("toBody() error = %v", err)
	}
	if body.SessionID != "sess-1" {
		t.Errorf("SessionID = %q, want sess-1", body.SessionID)
	}
	if len(body.Symbols) != 2 || body.Symbols[0] != "AAPL" || body.Symbols[1] != "TSLA" {
		t.Errorf("Symbols = %v, want [AAPL TSLA]", body.Symbols)
	}
	if body.Grab != "true" {
		t.Errorf("Grab = %q, want true", body.Grab)
	}
	if body.SubTypes[1] != SubTypeSnapshot {
		t.Errorf("SubTypes = %v, want normalized SNAPSHOT", body.SubTypes)
	}
}

func TestSubscribeRequestValidation(t *testing.T) {
	cases := map[string]SubscribeRequest{
		"no symbols":   {Category: CategoryUSStock, SubTypes: []SubType{SubTypeQuote}},
		"bad category": {Symbols: []string{"AAPL"}, Category: "BOGUS", SubTypes: []SubType{SubTypeQuote}},
		"no subtypes":  {Symbols: []string{"AAPL"}, Category: CategoryUSStock},
		"bad subtype":  {Symbols: []string{"AAPL"}, Category: CategoryUSStock, SubTypes: []SubType{"NOPE"}},
		"bad depth":    {Symbols: []string{"AAPL"}, Category: CategoryUSStock, SubTypes: []SubType{SubTypeQuote}, Depth: "abc"},
	}
	for name, req := range cases {
		if _, err := req.toBody("s"); err == nil {
			t.Errorf("toBody() with %s = nil, want error", name)
		}
	}
}

func TestUnsubscribeRequestBody(t *testing.T) {
	body, err := UnsubscribeRequest{
		Symbols:  []string{"AAPL"},
		Category: CategoryUSStock,
		SubTypes: []SubType{SubTypeTick},
	}.toBody("s")
	if err != nil {
		t.Fatalf("toBody() error = %v", err)
	}
	if body.UnsubscribeAll || body.Symbols[0] != "AAPL" {
		t.Fatalf("unsubscribe body = %+v", body)
	}

	all, err := UnsubscribeRequest{UnsubscribeAll: true}.toBody("s")
	if err != nil {
		t.Fatalf("toBody(all) error = %v", err)
	}
	if !all.UnsubscribeAll || len(all.Symbols) != 0 {
		t.Fatalf("unsubscribe-all body = %+v", all)
	}
}

// capturedRequest records one request received by the stub server.
type capturedRequest struct {
	path  string
	body  []byte
	token string
}

// TestSubscribeHTTPWiring verifies that Subscribe and Unsubscribe post the
// documented payloads to the documented HTTP paths through the core client.
func TestSubscribeHTTPWiring(t *testing.T) {
	captured := make(chan capturedRequest, 2)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		captured <- capturedRequest{
			path:  r.URL.Path,
			body:  body,
			token: r.Header.Get(client.AccessTokenHeader),
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	core := newTestCore(t, srv.URL)
	core.SetToken(&client.Token{Value: "tok", Status: client.TokenStatusNormal})

	s, err := New(core, WithSessionID("sess-42"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	ctx := context.Background()

	if err := s.Subscribe(ctx, SubscribeRequest{
		Symbols:  []string{"AAPL"},
		Category: CategoryUSStock,
		SubTypes: []SubType{SubTypeQuote, SubTypeSnapshot, SubTypeTick},
		Grab:     true,
	}); err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	req := <-captured
	if req.path != subscribePath {
		t.Errorf("subscribe path = %q, want %q", req.path, subscribePath)
	}
	var sub map[string]any
	if err := json.Unmarshal(req.body, &sub); err != nil {
		t.Fatalf("decode subscribe body: %v (body=%s)", err, req.body)
	}
	if sub["session_id"] != "sess-42" {
		t.Errorf("session_id = %v, want sess-42", sub["session_id"])
	}
	if sub["category"] != "US_STOCK" {
		t.Errorf("category = %v, want US_STOCK", sub["category"])
	}
	if sub["grab"] != "true" {
		t.Errorf("grab = %v, want string \"true\"", sub["grab"])
	}
	if req.token != "tok" {
		t.Errorf("x-access-token = %q, want tok", req.token)
	}

	if err := s.Unsubscribe(ctx, UnsubscribeRequest{UnsubscribeAll: true}); err != nil {
		t.Fatalf("Unsubscribe() error = %v", err)
	}
	req = <-captured
	if req.path != unsubscribePath {
		t.Errorf("unsubscribe path = %q, want %q", req.path, unsubscribePath)
	}
	var unsub map[string]any
	if err := json.Unmarshal(req.body, &unsub); err != nil {
		t.Fatalf("decode unsubscribe body: %v (body=%s)", err, req.body)
	}
	if unsub["session_id"] != "sess-42" || unsub["unsubscribe_all"] != true {
		t.Errorf("unsubscribe body = %v", unsub)
	}
}

// TestHandlerConcurrentRegistration exercises the handler registry under the
// race detector: registering while messages are dispatched must be safe.
func TestHandlerConcurrentRegistration(t *testing.T) {
	c := &Client{}
	data, err := proto.Marshal(&marketdatav1.Quote{Basic: &marketdatav1.Basic{Symbol: "AAPL"}})
	if err != nil {
		t.Fatalf("proto.Marshal() error = %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				c.OnQuote(func(*marketdatav1.Quote) {})
				if err := c.dispatch(TopicQuote, data); err != nil {
					t.Errorf("dispatch() error = %v", err)
				}
			}
		}()
	}
	wg.Wait()
}

// TestResolveBroker verifies endpoint selection between TCP, WebSocket, and an
// explicit override.
func TestResolveBroker(t *testing.T) {
	core := newTestCore(t, "http://127.0.0.1:1")
	endpoints := core.Endpoints()

	cfg := defaultConfig()
	cfg.sessionID = "s"
	got, err := cfg.resolveBroker(core)
	if err != nil {
		t.Fatalf("resolveBroker(tcp) error = %v", err)
	}
	if got != endpoints.MQTT {
		t.Errorf("resolveBroker(tcp) = %q, want %q", got, endpoints.MQTT)
	}

	cfg.websocket = true
	got, err = cfg.resolveBroker(core)
	if err != nil {
		t.Fatalf("resolveBroker(ws) error = %v", err)
	}
	if got != endpoints.MQTTWebSocket {
		t.Errorf("resolveBroker(ws) = %q, want %q", got, endpoints.MQTTWebSocket)
	}

	cfg.mqttURL = "wss://example.test/mqtt"
	got, err = cfg.resolveBroker(core)
	if err != nil {
		t.Fatalf("resolveBroker(override) error = %v", err)
	}
	if got != "wss://example.test/mqtt" {
		t.Errorf("resolveBroker(override) = %q", got)
	}
}
