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
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	marketdatav1 "github.com/shing1211/webullapi4go/gen/webull/marketdata/v1"
	mqtt "github.com/shing1211/webullapi4go/pkg/transport/mqtt"
)

func TestStreamHardeningDropOldestRemovesTheOldestBufferedMessage(t *testing.T) {
	c := &Client{chanReg: newChanRegistry()}
	ch, cancel := c.SubscribeQuoteChan(ChannelConfig{Policy: DropOldest, BufferSize: 2})
	defer cancel()

	first := &marketdatav1.Quote{Basic: &marketdatav1.Basic{Symbol: "first"}}
	second := &marketdatav1.Quote{Basic: &marketdatav1.Basic{Symbol: "second"}}
	third := &marketdatav1.Quote{Basic: &marketdatav1.Basic{Symbol: "third"}}
	c.emitQuote(first)
	c.emitQuote(second)
	c.emitQuote(third)

	if got := <-ch; got != second {
		t.Fatalf("first buffered message = %q, want %q", got.GetBasic().GetSymbol(), "second")
	}
	if got := <-ch; got != third {
		t.Fatalf("second buffered message = %q, want %q", got.GetBasic().GetSymbol(), "third")
	}

	c.chanReg.mu.RLock()
	var drops int64
	for _, cfg := range c.chanReg.quote {
		drops = cfg.DropCount()
	}
	c.chanReg.mu.RUnlock()
	if drops != 1 {
		t.Fatalf("DropCount() = %d, want 1", drops)
	}
}

func TestStreamHardeningChannelCancellationUnblocksDispatchAndIsIdempotent(t *testing.T) {
	c := &Client{chanReg: newChanRegistry()}
	ch, cancel := c.SubscribeQuoteChan(ChannelConfig{Policy: DropBlock, BufferSize: 1})
	first := &marketdatav1.Quote{Basic: &marketdatav1.Basic{Symbol: "first"}}
	second := &marketdatav1.Quote{Basic: &marketdatav1.Basic{Symbol: "second"}}
	c.emitQuote(first)

	c.chanReg.mu.RLock()
	var entry *chanQuote
	for candidate := range c.chanReg.quote {
		entry = candidate
	}
	c.chanReg.mu.RUnlock()
	if entry == nil {
		t.Fatal("quote channel entry was not registered")
	}

	dispatchStarted := make(chan struct{})
	dispatchDone := make(chan struct{})
	go func() {
		close(dispatchStarted)
		c.emitQuote(second)
		close(dispatchDone)
	}()
	<-dispatchStarted

	deadline := time.Now().Add(time.Second)
	for entry.lifecycle.sendMu.TryLock() {
		entry.lifecycle.sendMu.Unlock()
		runtime.Gosched()
		if time.Now().After(deadline) {
			t.Fatal("dispatch did not reach the blocked send")
		}
	}

	cancelDone := make(chan struct{})
	go func() {
		cancel()
		cancel()
		close(cancelDone)
	}()
	select {
	case <-cancelDone:
	case <-time.After(time.Second):
		t.Fatal("channel cancellation remained blocked")
	}
	select {
	case <-dispatchDone:
	case <-time.After(time.Second):
		t.Fatal("dispatch remained blocked after cancellation")
	}

	c.chanReg.stopAll()
	assertStreamHardeningChannelClosed(t, ch)
}

func TestStreamHardeningClientCloseClosesChannelsWhileDispatchIsBlocked(t *testing.T) {
	c := &Client{chanReg: newChanRegistry()}
	ch, cancel := c.SubscribeQuoteChan(ChannelConfig{Policy: DropBlock, BufferSize: 1})
	c.emitQuote(&marketdatav1.Quote{Basic: &marketdatav1.Basic{Symbol: "first"}})

	c.chanReg.mu.RLock()
	var entry *chanQuote
	for candidate := range c.chanReg.quote {
		entry = candidate
	}
	c.chanReg.mu.RUnlock()

	dispatchDone := make(chan struct{})
	go func() {
		c.emitQuote(&marketdatav1.Quote{Basic: &marketdatav1.Basic{Symbol: "second"}})
		close(dispatchDone)
	}()
	deadline := time.Now().Add(time.Second)
	for entry.lifecycle.sendMu.TryLock() {
		entry.lifecycle.sendMu.Unlock()
		runtime.Gosched()
		if time.Now().After(deadline) {
			t.Fatal("dispatch did not reach the blocked send")
		}
	}

	closeDone := make(chan struct{})
	go func() {
		_ = c.Close()
		close(closeDone)
	}()
	select {
	case <-closeDone:
	case <-time.After(time.Second):
		t.Fatal("Client.Close() remained blocked")
	}
	select {
	case <-dispatchDone:
	case <-time.After(time.Second):
		t.Fatal("dispatch remained blocked after Client.Close()")
	}
	if got := c.State(); got != StateClosed {
		t.Fatalf("State() after Close() = %s, want closed", got)
	}
	cancel()
	cancel()
	assertStreamHardeningChannelClosed(t, ch)
}

func assertStreamHardeningChannelClosed[T any](t *testing.T, ch <-chan T) {
	t.Helper()
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case <-timer.C:
			t.Fatal("channel did not close")
		}
	}
}

func TestStreamHardeningHealthDegradationWaitsForFreshDataToRecover(t *testing.T) {
	interval := time.Minute
	c := &Client{cfg: config{healthWatchdogInterval: interval}}
	c.state.Store(StateConnected)
	now := time.Unix(1_700_000_000, 0)
	c.lastMessageAt.Store(now.Add(-2 * interval).UnixMilli())

	c.checkHealth(now)
	if got := c.State(); got != StateDegraded {
		t.Fatalf("State() after stale data = %s, want degraded", got)
	}
	c.checkHealth(now.Add(2 * interval))
	if got := c.State(); got != StateDegraded {
		t.Fatalf("State() without fresh data = %s, want degraded", got)
	}
	c.markDataMessage()
	if got := c.State(); got != StateConnected {
		t.Fatalf("State() after fresh data = %s, want connected", got)
	}
}

func TestStreamHardeningConnectionLossPreservesReconnectState(t *testing.T) {
	c := &Client{cfg: config{autoReconnect: true}}
	c.state.Store(StateReconnecting)
	c.reconnecting.Store(true)
	c.handleConnectionLost(io.EOF)
	if got := c.State(); got != StateReconnecting {
		t.Fatalf("State() after connection loss = %s, want reconnecting", got)
	}
	if !c.Reconnecting() {
		t.Fatal("Reconnecting() = false after connection loss")
	}
}

func TestStreamHardeningNonDataMessagesDoNotRefreshHealth(t *testing.T) {
	c := &Client{}
	c.handleMessage(mqtt.Message{Topic: TopicNotice, Payload: []byte(`{"type":"status"}`)})
	if got := c.lastMessageAt.Load(); got != 0 {
		t.Fatalf("lastMessageAt after notice = %d, want 0", got)
	}
}

type streamHardeningResubscribeBarrier struct {
	subscribeCount             atomic.Int32
	replayStarted              chan struct{}
	releaseReplay              chan struct{}
	unsubscribeMutationArrived chan struct{}
	subscribeMutationArrived   chan struct{}
	replayOnce                 sync.Once
	unsubscribeMutationOnce    sync.Once
	subscribeMutationOnce      sync.Once
}

func (b *streamHardeningResubscribeBarrier) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = io.ReadAll(r.Body)
	if r.URL.Path == subscribePath {
		count := b.subscribeCount.Add(1)
		if count == 2 {
			b.replayOnce.Do(func() { close(b.replayStarted) })
			<-b.releaseReplay
		}
		if count > 2 {
			b.subscribeMutationOnce.Do(func() { close(b.subscribeMutationArrived) })
		}
	}
	if r.URL.Path == unsubscribePath {
		b.unsubscribeMutationOnce.Do(func() { close(b.unsubscribeMutationArrived) })
	}
	w.WriteHeader(http.StatusOK)
}

func newStreamHardeningBarrierClient(t *testing.T, barrier *streamHardeningResubscribeBarrier) *Client {
	t.Helper()
	srv := httptest.NewServer(barrier)
	t.Cleanup(srv.Close)
	core := newTestCore(t, srv.URL)
	core.SetToken(&client.Token{Value: "tok", Status: client.TokenStatusNormal})
	s, err := New(core, WithSessionID("barrier-session"), WithAutoReconnect(true))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestStreamHardeningResubscribeSerializesWithSubscriptionMutations(t *testing.T) {
	barrier := &streamHardeningResubscribeBarrier{
		replayStarted:              make(chan struct{}),
		releaseReplay:              make(chan struct{}),
		unsubscribeMutationArrived: make(chan struct{}),
		subscribeMutationArrived:   make(chan struct{}),
	}
	s := newStreamHardeningBarrierClient(t, barrier)
	ctx := context.Background()
	if err := s.Subscribe(ctx, SubscribeRequest{
		Symbols:  []string{"AAPL"},
		Category: CategoryUSStock,
		SubTypes: []SubType{SubTypeQuote},
	}); err != nil {
		t.Fatalf("initial Subscribe() error = %v", err)
	}
	s.handleConnect()

	reconnectDone := make(chan struct{})
	go func() {
		s.handleConnect()
		close(reconnectDone)
	}()
	select {
	case <-barrier.replayStarted:
	case <-time.After(time.Second):
		t.Fatal("re-subscribe request did not start")
	}

	unsubscribeStarted := make(chan struct{})
	unsubscribeDone := make(chan error, 1)
	go func() {
		close(unsubscribeStarted)
		unsubscribeDone <- s.Unsubscribe(ctx, UnsubscribeRequest{
			Symbols:  []string{"AAPL"},
			Category: CategoryUSStock,
			SubTypes: []SubType{SubTypeQuote},
		})
	}()
	subscribeStarted := make(chan struct{})
	subscribeDone := make(chan error, 1)
	go func() {
		close(subscribeStarted)
		subscribeDone <- s.Subscribe(ctx, SubscribeRequest{
			Symbols:  []string{"TSLA"},
			Category: CategoryUSStock,
			SubTypes: []SubType{SubTypeQuote},
		})
	}()
	<-unsubscribeStarted
	<-subscribeStarted
	premature := false
	select {
	case <-barrier.unsubscribeMutationArrived:
		premature = true
	case <-barrier.subscribeMutationArrived:
		premature = true
	case <-time.After(100 * time.Millisecond):
	}
	close(barrier.releaseReplay)

	select {
	case <-reconnectDone:
	case <-time.After(time.Second):
		t.Fatal("re-subscribe did not finish")
	}
	select {
	case err := <-unsubscribeDone:
		if err != nil {
			t.Fatalf("Unsubscribe() error = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Unsubscribe() did not finish")
	}
	select {
	case err := <-subscribeDone:
		if err != nil {
			t.Fatalf("Subscribe() error = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Subscribe() did not finish")
	}
	if premature {
		t.Fatal("a subscription mutation reached the server while re-subscribe was in progress")
	}
	assertCounts(t, flattenRequested(s.subs.requests()), map[string]int{"TSLA/QUOTE": 1})

	s.handleConnect()
	if got := barrier.subscribeCount.Load(); got != 4 {
		t.Fatalf("subscribe request count = %d, want 4", got)
	}
}
