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
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

// subscribeRecorder is a stub HTTP server that records the decoded bodies of
// streaming subscribe requests. Unsubscribe requests are acknowledged but not
// recorded.
type subscribeRecorder struct {
	mu   sync.Mutex
	subs []map[string]any
}

func (rec *subscribeRecorder) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == subscribePath {
		body, _ := io.ReadAll(r.Body)
		var decoded map[string]any
		if err := json.Unmarshal(body, &decoded); err == nil {
			rec.mu.Lock()
			rec.subs = append(rec.subs, decoded)
			rec.mu.Unlock()
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (rec *subscribeRecorder) snapshot() []map[string]any {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	out := make([]map[string]any, len(rec.subs))
	copy(out, rec.subs)
	return out
}

func (rec *subscribeRecorder) len() int {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	return len(rec.subs)
}

// newReconnectTestClient builds a streaming client whose core points at a stub
// subscribe server, so disconnect/reconnect can be simulated without a broker.
func newReconnectTestClient(t *testing.T, rec *subscribeRecorder, opts ...Option) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(rec.serveHTTP))
	t.Cleanup(srv.Close)

	core := newTestCore(t, srv.URL)
	core.SetToken(&client.Token{Value: "tok", Status: client.TokenStatusNormal})

	all := append([]Option{WithSessionID("sess-reconnect")}, opts...)
	s, err := New(core, all...)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

// mustSubscribe subscribes symbols with the given sub types and fails the test
// on error.
func mustSubscribe(t *testing.T, s *Client, ctx context.Context, symbols []string, subTypes ...SubType) {
	t.Helper()
	err := s.Subscribe(ctx, SubscribeRequest{
		Symbols:  symbols,
		Category: CategoryUSStock,
		SubTypes: subTypes,
	})
	if err != nil {
		t.Fatalf("Subscribe(%v) error = %v", symbols, err)
	}
}

// flattenRecorded counts each recorded (symbol, sub type) pair so duplicates are
// observable.
func flattenRecorded(reqs []map[string]any) map[string]int {
	counts := make(map[string]int)
	for _, req := range reqs {
		symbols, _ := req["symbols"].([]any)
		subTypes, _ := req["sub_types"].([]any)
		for _, symbol := range symbols {
			for _, subType := range subTypes {
				counts[symbol.(string)+"/"+subType.(string)]++
			}
		}
	}
	return counts
}

// flattenRequested counts each (symbol, sub type) pair in a reconstructed
// request slice.
func flattenRequested(reqs []SubscribeRequest) map[string]int {
	counts := make(map[string]int)
	for _, req := range reqs {
		for _, symbol := range req.Symbols {
			for _, subType := range req.SubTypes {
				counts[symbol+"/"+string(subType)]++
			}
		}
	}
	return counts
}

// assertCounts compares got with want and reports both missing/extra keys and
// wrong multiplicities.
func assertCounts(t *testing.T, got, want map[string]int) {
	t.Helper()
	for key, n := range want {
		if got[key] != n {
			t.Errorf("subscription %s issued %d times, want %d (all=%v)", key, got[key], n, got)
		}
	}
	for key := range got {
		if _, ok := want[key]; !ok {
			t.Errorf("unexpected subscription %s re-issued (all=%v)", key, got)
		}
	}
}

func addBody(t *testing.T, reg *subscriptionRegistry, req SubscribeRequest) {
	t.Helper()
	body, err := req.toBody("sess")
	if err != nil {
		t.Fatalf("toBody() error = %v", err)
	}
	reg.add(body)
}

func TestSubscriptionRegistryGrouping(t *testing.T) {
	var reg subscriptionRegistry
	addBody(t, &reg, SubscribeRequest{Symbols: []string{"AAPL"}, Category: CategoryUSStock, SubTypes: []SubType{SubTypeQuote}})
	addBody(t, &reg, SubscribeRequest{Symbols: []string{"TSLA"}, Category: CategoryUSStock, SubTypes: []SubType{SubTypeQuote, SubTypeSnapshot}})

	reqs := reg.requests()
	if len(reqs) != 2 {
		t.Fatalf("requests() returned %d requests, want 2: %+v", len(reqs), reqs)
	}
	for _, req := range reqs {
		switch req.Symbols[0] {
		case "AAPL":
			if len(req.SubTypes) != 1 || req.SubTypes[0] != SubTypeQuote {
				t.Errorf("AAPL request = %+v, want only QUOTE", req)
			}
		case "TSLA":
			if len(req.SubTypes) != 2 {
				t.Errorf("TSLA request = %+v, want QUOTE and SNAPSHOT", req)
			}
		default:
			t.Errorf("unexpected request = %+v", req)
		}
	}
}

func TestSubscriptionRegistryIdempotent(t *testing.T) {
	var reg subscriptionRegistry
	req := SubscribeRequest{Symbols: []string{"AAPL"}, Category: CategoryUSStock, SubTypes: []SubType{SubTypeQuote, SubTypeSnapshot}}
	addBody(t, &reg, req)
	addBody(t, &reg, req)

	reqs := reg.requests()
	if len(reqs) != 1 {
		t.Fatalf("requests() returned %d requests, want 1", len(reqs))
	}
	assertCounts(t, flattenRequested(reqs), map[string]int{"AAPL/QUOTE": 1, "AAPL/SNAPSHOT": 1})
}

func TestSubscriptionRegistryRemove(t *testing.T) {
	var reg subscriptionRegistry
	addBody(t, &reg, SubscribeRequest{Symbols: []string{"AAPL", "TSLA"}, Category: CategoryUSStock, SubTypes: []SubType{SubTypeQuote}})

	unsub, err := UnsubscribeRequest{Symbols: []string{"TSLA"}, Category: CategoryUSStock, SubTypes: []SubType{SubTypeQuote}}.toBody("sess")
	if err != nil {
		t.Fatalf("toBody(unsubscribe) error = %v", err)
	}
	reg.remove(unsub)
	assertCounts(t, flattenRequested(reg.requests()), map[string]int{"AAPL/QUOTE": 1})

	all, err := UnsubscribeRequest{UnsubscribeAll: true}.toBody("sess")
	if err != nil {
		t.Fatalf("toBody(unsubscribe all) error = %v", err)
	}
	reg.remove(all)
	if !reg.empty() {
		t.Fatal("registry is not empty after UnsubscribeAll")
	}
	if reqs := reg.requests(); len(reqs) != 0 {
		t.Fatalf("requests() after UnsubscribeAll = %+v, want none", reqs)
	}
}

func TestAutoResubscribeOnReconnect(t *testing.T) {
	rec := &subscribeRecorder{}
	s := newReconnectTestClient(t, rec, WithAutoReconnect(true))

	var connects, disconnects int32
	s.OnConnect(func() { atomic.AddInt32(&connects, 1) })
	s.OnDisconnect(func(error) { atomic.AddInt32(&disconnects, 1) })

	ctx := context.Background()
	mustSubscribe(t, s, ctx, []string{"AAPL"}, SubTypeQuote, SubTypeSnapshot)
	mustSubscribe(t, s, ctx, []string{"TSLA"}, SubTypeQuote, SubTypeSnapshot)
	if got := rec.len(); got != 2 {
		t.Fatalf("initial subscribe calls = %d, want 2", got)
	}

	// The first connection must not replay anything.
	s.handleConnect()
	if got := rec.len(); got != 2 {
		t.Fatalf("first connect replayed subscriptions: calls = %d, want 2", got)
	}
	if got := atomic.LoadInt32(&connects); got != 1 {
		t.Fatalf("OnConnect count = %d, want 1", got)
	}

	// Lost connection, reconnect: the active set is re-issued exactly once.
	s.handleReconnecting()
	if !s.Reconnecting() {
		t.Fatal("Reconnecting() = false after a reconnect started")
	}
	s.handleConnectionLost(errors.New("connection lost"))
	s.handleConnect()
	if got := atomic.LoadInt32(&connects); got != 2 {
		t.Fatalf("OnConnect count = %d, want 2", got)
	}
	if got := atomic.LoadInt32(&disconnects); got != 1 {
		t.Fatalf("OnDisconnect count = %d, want 1", got)
	}
	if s.Reconnecting() {
		t.Fatal("Reconnecting() = true after a successful reconnect")
	}
	assertCounts(t, flattenRecorded(rec.snapshot()[2:]), map[string]int{
		"AAPL/QUOTE": 1, "AAPL/SNAPSHOT": 1, "TSLA/QUOTE": 1, "TSLA/SNAPSHOT": 1,
	})

	// Unsubscribing removes the subscription so a later reconnect omits it.
	if err := s.Unsubscribe(ctx, UnsubscribeRequest{
		Symbols:  []string{"TSLA"},
		Category: CategoryUSStock,
		SubTypes: []SubType{SubTypeQuote, SubTypeSnapshot},
	}); err != nil {
		t.Fatalf("Unsubscribe(TSLA) error = %v", err)
	}
	before := rec.len()
	s.handleConnectionLost(errors.New("again"))
	s.handleConnect()
	assertCounts(t, flattenRecorded(rec.snapshot()[before:]), map[string]int{
		"AAPL/QUOTE": 1, "AAPL/SNAPSHOT": 1,
	})

	// UnsubscribeAll clears the registry: the next reconnect issues nothing.
	if err := s.Unsubscribe(ctx, UnsubscribeRequest{UnsubscribeAll: true}); err != nil {
		t.Fatalf("Unsubscribe(all) error = %v", err)
	}
	before = rec.len()
	s.handleConnectionLost(errors.New("third"))
	s.handleConnect()
	if got := rec.len(); got != before {
		t.Fatalf("reconnect after UnsubscribeAll re-issued %d subscriptions", got-before)
	}
}

func TestAutoResubscribeDisabled(t *testing.T) {
	rec := &subscribeRecorder{}
	s := newReconnectTestClient(t, rec, WithAutoReconnect(true), WithAutoResubscribe(false))

	ctx := context.Background()
	mustSubscribe(t, s, ctx, []string{"AAPL"}, SubTypeQuote)

	s.handleConnect()
	before := rec.len()
	s.handleConnectionLost(errors.New("connection lost"))
	s.handleConnect()
	if got := rec.len(); got != before {
		t.Fatalf("auto-resubscribe disabled but %d subscriptions were re-issued", got-before)
	}
}

func TestResubscribeUsesClientSession(t *testing.T) {
	rec := &subscribeRecorder{}
	s := newReconnectTestClient(t, rec, WithAutoReconnect(true))
	ctx := context.Background()
	mustSubscribe(t, s, ctx, []string{"AAPL"}, SubTypeTick)

	s.handleConnect()
	s.handleConnectionLost(errors.New("connection lost"))
	s.handleConnect()

	replayed := rec.snapshot()[1:]
	if len(replayed) != 1 {
		t.Fatalf("replayed %d requests, want 1", len(replayed))
	}
	if got := replayed[0]["session_id"]; got != "sess-reconnect" {
		t.Fatalf("replayed session_id = %v, want sess-reconnect", got)
	}
}

func TestSubscribeAndUnsubscribeAfterClose(t *testing.T) {
	rec := &subscribeRecorder{}
	s := newReconnectTestClient(t, rec)
	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if err := s.Subscribe(context.Background(), SubscribeRequest{
		Symbols:  []string{"AAPL"},
		Category: CategoryUSStock,
		SubTypes: []SubType{SubTypeQuote},
	}); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("Subscribe() after Close error = %v, want invalid_config", err)
	}
	if err := s.Unsubscribe(context.Background(), UnsubscribeRequest{UnsubscribeAll: true}); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("Unsubscribe() after Close error = %v, want invalid_config", err)
	}
	if got := rec.len(); got != 0 {
		t.Fatalf("HTTP requests after Close = %d, want 0", got)
	}
}

func TestSubscribeConcurrentWithCloseDoesNotRecordSubscription(t *testing.T) {
	entered := make(chan struct{})
	release := make(chan struct{})
	var enteredOnce sync.Once
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == subscribePath {
			enteredOnce.Do(func() { close(entered) })
			<-release
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	released := false
	defer func() {
		if !released {
			close(release)
		}
	}()

	core := newTestCore(t, srv.URL)
	core.SetToken(&client.Token{Value: "tok", Status: client.TokenStatusNormal})
	s, err := New(core, WithSessionID("close-race-session"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })

	done := make(chan error, 1)
	go func() {
		done <- s.Subscribe(context.Background(), SubscribeRequest{
			Symbols:  []string{"AAPL"},
			Category: CategoryUSStock,
			SubTypes: []SubType{SubTypeQuote},
		})
	}()
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("Subscribe request did not reach the server")
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	close(release)
	released = true
	select {
	case err := <-done:
		if !errs.Is(err, errs.CodeInvalidConfig) {
			t.Fatalf("Subscribe() racing Close error = %v, want invalid_config", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Subscribe() racing Close did not finish")
	}
	if !s.subs.empty() {
		t.Fatal("subscription was recorded after Close")
	}
}
