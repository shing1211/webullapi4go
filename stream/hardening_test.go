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
	"reflect"
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

func TestStreamHardeningLifecycleCallbacksAndTerminalState(t *testing.T) {
	c := &Client{cfg: config{autoReconnect: true}, chanReg: newChanRegistry()}
	c.state.Store(StateDisconnected)
	var events []string
	record := func(event string) {
		events = append(events, event)
	}
	c.OnStateChange(func(prev, next State) {
		record("state:" + prev.String() + ">" + next.String())
	})
	c.OnConnect(func() {
		record("connect:" + c.State().String())
	})
	c.OnDisconnect(func(err error) {
		if err != io.EOF {
			t.Errorf("OnDisconnect error = %v, want EOF", err)
		}
		record("disconnect:" + c.State().String())
	})
	c.OnReconnecting(func() {
		record("reconnecting:" + c.State().String())
	})

	c.setState(StateConnecting)
	c.handleConnect()
	interval := time.Minute
	now := time.Unix(1_700_000_000, 0)
	c.lastMessageAt.Store(now.Add(-2 * interval).UnixMilli())
	c.cfg.healthWatchdogInterval = interval
	c.checkHealth(now)
	c.markDataMessage()
	c.handleConnectionLost(io.EOF)
	c.handleReconnecting()
	c.handleConnect()
	if err := c.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	c.handleConnectionLost(io.EOF)
	c.handleReconnecting()
	c.handleConnect()
	c.markDataMessage()
	c.checkHealth(now.Add(2 * interval))
	c.setState(StateConnected)
	c.setState(StateClosed)

	want := []string{
		"state:disconnected>connecting",
		"state:connecting>connected",
		"connect:connected",
		"state:connected>degraded",
		"state:degraded>connected",
		"state:connected>disconnected",
		"disconnect:disconnected",
		"state:disconnected>reconnecting",
		"reconnecting:reconnecting",
		"state:reconnecting>connected",
		"connect:connected",
		"state:connected>closed",
	}
	if !reflect.DeepEqual(events, want) {
		t.Fatalf("lifecycle events = %v, want %v", events, want)
	}
	if got := c.State(); got != StateClosed {
		t.Fatalf("terminal State() = %s, want closed", got)
	}
}

type streamHardeningWatchdogTicker struct {
	ch       chan time.Time
	stopped  chan struct{}
	stopOnce sync.Once
}

func (t *streamHardeningWatchdogTicker) C() <-chan time.Time {
	return t.ch
}

func (t *streamHardeningWatchdogTicker) Stop() {
	t.stopOnce.Do(func() { close(t.stopped) })
}

func TestStreamHardeningWatchdogStartTickRecoverAndClose(t *testing.T) {
	interval := time.Minute
	now := time.Unix(1_700_000_000, 0)
	ticker := &streamHardeningWatchdogTicker{
		ch:      make(chan time.Time, 2),
		stopped: make(chan struct{}),
	}
	connCtx, connCancel := context.WithCancel(context.Background())
	c := &Client{
		cfg:        config{healthWatchdogInterval: interval},
		connCtx:    connCtx,
		connCancel: connCancel,
		chanReg:    newChanRegistry(),
	}
	c.state.Store(StateConnected)
	c.lastMessageAt.Store(now.Add(-2 * interval).UnixMilli())
	created := make(chan time.Duration, 1)
	c.newWatchdogTicker = func(got time.Duration) watchdogTicker {
		created <- got
		return ticker
	}
	c.healthClock = func() time.Time { return now }
	transitions := make(chan [2]State, 2)
	c.OnStateChange(func(prev, next State) {
		transitions <- [2]State{prev, next}
	})

	c.startHealthWatchdog()
	select {
	case got := <-created:
		if got != interval {
			t.Fatalf("watchdog interval = %v, want %v", got, interval)
		}
	case <-time.After(time.Second):
		t.Fatal("watchdog did not start")
	}
	c.startHealthWatchdog()
	select {
	case <-created:
		t.Fatal("watchdog started more than once")
	default:
	}

	ticker.ch <- now
	select {
	case got := <-transitions:
		if got != [2]State{StateConnected, StateDegraded} {
			t.Fatalf("stale transition = %v, want connected>degraded", got)
		}
	case <-time.After(time.Second):
		t.Fatal("watchdog did not process stale tick")
	}

	c.lastMessageAt.Store(now.UnixMilli())
	ticker.ch <- now
	select {
	case got := <-transitions:
		if got != [2]State{StateDegraded, StateConnected} {
			t.Fatalf("fresh transition = %v, want degraded>connected", got)
		}
	case <-time.After(time.Second):
		t.Fatal("watchdog did not process fresh tick")
	}

	if err := c.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	select {
	case <-ticker.stopped:
	case <-time.After(time.Second):
		t.Fatal("watchdog ticker was not stopped")
	}
	select {
	case <-c.healthDone:
	case <-time.After(time.Second):
		t.Fatal("watchdog goroutine did not finish")
	}
	c.startHealthWatchdog()
	select {
	case <-created:
		t.Fatal("closed client started another watchdog")
	default:
	}
	if got := c.State(); got != StateClosed {
		t.Fatalf("State() after watchdog Close() = %s, want closed", got)
	}
}

func TestStreamHardeningStaleHealthCheckCannotOverwriteReconnect(t *testing.T) {
	interval := time.Minute
	now := time.Unix(1_700_000_000, 0)
	c := &Client{cfg: config{healthWatchdogInterval: interval}, chanReg: newChanRegistry()}
	c.state.Store(StateConnected)
	c.lastMessageAt.Store(now.Add(-2 * interval).UnixMilli())
	observed := make(chan struct{})
	release := make(chan struct{})
	done := make(chan struct{})
	go func() {
		c.checkHealthWithTransition(now, func(expected, next State) bool {
			close(observed)
			<-release
			return c.setStateIfCurrent(expected, next)
		})
		close(done)
	}()
	<-observed
	c.setState(StateReconnecting)
	close(release)
	<-done
	if got := c.State(); got != StateReconnecting {
		t.Fatalf("State() after stale health check = %s, want reconnecting", got)
	}
}

func TestStreamHardeningHealthTransitionsRequireCurrentState(t *testing.T) {
	interval := time.Minute
	now := time.Unix(1_700_000_000, 0)
	tests := []struct {
		name  string
		state State
		last  time.Time
		now   time.Time
		want  State
	}{
		{name: "stale connected", state: StateConnected, last: now.Add(-2 * interval), now: now, want: StateDegraded},
		{name: "fresh degraded", state: StateDegraded, last: now, now: now, want: StateConnected},
		{name: "stale reconnecting", state: StateReconnecting, last: now.Add(-2 * interval), now: now, want: StateReconnecting},
		{name: "stale disconnected", state: StateDisconnected, last: now.Add(-2 * interval), now: now, want: StateDisconnected},
		{name: "stale closed", state: StateClosed, last: now.Add(-2 * interval), now: now, want: StateClosed},
		{name: "fresh reconnecting", state: StateReconnecting, last: now, now: now, want: StateReconnecting},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{cfg: config{healthWatchdogInterval: interval}, chanReg: newChanRegistry()}
			c.state.Store(tt.state)
			c.lastMessageAt.Store(tt.last.UnixMilli())
			c.checkHealth(tt.now)
			if got := c.State(); got != tt.want {
				t.Fatalf("State() = %s, want %s", got, tt.want)
			}
		})
	}

	c := &Client{chanReg: newChanRegistry()}
	c.state.Store(StateReconnecting)
	c.markDataMessage()
	if got := c.State(); got != StateReconnecting {
		t.Fatalf("State() after stale data mark = %s, want reconnecting", got)
	}
}

type streamHardeningChannelHarness[T any] struct {
	subscribe      func(*Client, ChannelConfig) (<-chan T, func())
	emit           func(*Client, T)
	config         func(*Client) *channelConfig
	lifecycle      func(*Client) *channelLifecycle
	lifecycleOrder func(*Client, uint64) *channelLifecycle
	first          T
	second         T
	third          T
}

func streamHardeningQuoteConfig(c *Client) *channelConfig {
	c.chanReg.mu.RLock()
	defer c.chanReg.mu.RUnlock()
	for _, cfg := range c.chanReg.quote {
		return cfg
	}
	return nil
}

func streamHardeningQuoteLifecycle(c *Client) *channelLifecycle {
	c.chanReg.mu.RLock()
	defer c.chanReg.mu.RUnlock()
	for entry := range c.chanReg.quote {
		return &entry.lifecycle
	}
	return nil
}

func streamHardeningQuoteLifecycleByOrder(c *Client, order uint64) *channelLifecycle {
	c.chanReg.mu.RLock()
	defer c.chanReg.mu.RUnlock()
	for entry := range c.chanReg.quote {
		if entry.order == order {
			return &entry.lifecycle
		}
	}
	return nil
}

func streamHardeningSnapshotLifecycleByOrder(c *Client, order uint64) *channelLifecycle {
	c.chanReg.mu.RLock()
	defer c.chanReg.mu.RUnlock()
	for entry := range c.chanReg.snapshot {
		if entry.order == order {
			return &entry.lifecycle
		}
	}
	return nil
}

func streamHardeningSnapshotConfig(c *Client) *channelConfig {
	c.chanReg.mu.RLock()
	defer c.chanReg.mu.RUnlock()
	for _, cfg := range c.chanReg.snapshot {
		return cfg
	}
	return nil
}

func streamHardeningSnapshotLifecycle(c *Client) *channelLifecycle {
	c.chanReg.mu.RLock()
	defer c.chanReg.mu.RUnlock()
	for entry := range c.chanReg.snapshot {
		return &entry.lifecycle
	}
	return nil
}

func streamHardeningTickLifecycleByOrder(c *Client, order uint64) *channelLifecycle {
	c.chanReg.mu.RLock()
	defer c.chanReg.mu.RUnlock()
	for entry := range c.chanReg.tick {
		if entry.order == order {
			return &entry.lifecycle
		}
	}
	return nil
}

func streamHardeningTickConfig(c *Client) *channelConfig {
	c.chanReg.mu.RLock()
	defer c.chanReg.mu.RUnlock()
	for _, cfg := range c.chanReg.tick {
		return cfg
	}
	return nil
}

func streamHardeningTickLifecycle(c *Client) *channelLifecycle {
	c.chanReg.mu.RLock()
	defer c.chanReg.mu.RUnlock()
	for entry := range c.chanReg.tick {
		return &entry.lifecycle
	}
	return nil
}

func receiveStreamHardeningValue[T any](t *testing.T, ch <-chan T) T {
	t.Helper()
	select {
	case value := <-ch:
		return value
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for channel value")
		var zero T
		return zero
	}
}

func waitForStreamHardeningSend(t *testing.T, lifecycle *channelLifecycle) {
	t.Helper()
	if lifecycle == nil {
		t.Fatal("channel entry was not registered")
	}
	deadline := time.Now().Add(time.Second)
	for lifecycle.sendMu.TryLock() {
		lifecycle.sendMu.Unlock()
		runtime.Gosched()
		if time.Now().After(deadline) {
			t.Fatal("dispatch did not reach the channel send")
		}
	}
}

func assertStreamHardeningDropOldest[T any](t *testing.T, h streamHardeningChannelHarness[T]) {
	c := &Client{chanReg: newChanRegistry()}
	ch, cancel := h.subscribe(c, ChannelConfig{Policy: DropOldest, BufferSize: 2})
	h.emit(c, h.first)
	h.emit(c, h.second)
	h.emit(c, h.third)
	if got := receiveStreamHardeningValue(t, ch); !reflect.DeepEqual(got, h.second) {
		t.Fatalf("first retained value = %v, want second", got)
	}
	if got := receiveStreamHardeningValue(t, ch); !reflect.DeepEqual(got, h.third) {
		t.Fatalf("second retained value = %v, want third", got)
	}
	cfg := h.config(c)
	if cfg == nil || cfg.DropCount() != 1 {
		t.Fatalf("DropOldest drop count = %v, want 1", cfg)
	}
	cancel()
}

func assertStreamHardeningDropSample[T any](t *testing.T, h streamHardeningChannelHarness[T]) {
	c := &Client{chanReg: newChanRegistry()}
	ch, cancel := h.subscribe(c, ChannelConfig{Policy: DropSample, BufferSize: 1})
	cfg := h.config(c)
	if cfg == nil {
		t.Fatal("DropSample channel was not registered")
	}
	var keep atomic.Bool
	cfg.sampleKeep = keep.Load
	h.emit(c, h.first)
	h.emit(c, h.second)
	if got := len(ch); got != 1 {
		t.Fatalf("buffer length after deterministic drop = %d, want 1", got)
	}
	if got := receiveStreamHardeningValue(t, ch); !reflect.DeepEqual(got, h.first) {
		t.Fatalf("first DropSample value = %v, want first", got)
	}
	keep.Store(true)
	h.emit(c, h.third)
	if got := receiveStreamHardeningValue(t, ch); !reflect.DeepEqual(got, h.third) {
		t.Fatalf("second DropSample value = %v, want third", got)
	}
	if got := cfg.DropCount(); got != 1 {
		t.Fatalf("DropSample drop count = %d, want 1", got)
	}
	cancel()
}

func assertStreamHardeningBlockingPolicy[T any](t *testing.T, h streamHardeningChannelHarness[T], policy DropPolicy) {
	c := &Client{chanReg: newChanRegistry()}
	ch, cancel := h.subscribe(c, ChannelConfig{Policy: policy, BufferSize: 1})
	cfg := h.config(c)
	if cfg == nil {
		t.Fatal("blocking channel was not registered")
	}
	if policy == DropPolicy(99) && policy.String() != "block" {
		t.Fatalf("invalid DropPolicy.String() = %q, want block", policy.String())
	}
	h.emit(c, h.first)
	started := make(chan struct{})
	done := make(chan struct{})
	go func() {
		close(started)
		h.emit(c, h.second)
		close(done)
	}()
	<-started
	waitForStreamHardeningSend(t, h.lifecycle(c))
	select {
	case <-done:
		t.Fatal("DropBlock dispatch completed while the channel was full")
	default:
	}
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("DropBlock dispatch did not unblock after cancellation")
	}
	assertStreamHardeningChannelClosed(t, ch)
}

func assertStreamHardeningAllPolicies[T any](t *testing.T, h streamHardeningChannelHarness[T]) {
	t.Run("drop-oldest", func(t *testing.T) {
		assertStreamHardeningDropOldest(t, h)
	})
	t.Run("drop-sample", func(t *testing.T) {
		assertStreamHardeningDropSample(t, h)
	})
	t.Run("drop-block", func(t *testing.T) {
		assertStreamHardeningBlockingPolicy(t, h, DropBlock)
	})
}

func assertStreamHardeningDefaultsAndInvalidPolicy[T any](t *testing.T, h streamHardeningChannelHarness[T]) {
	for _, tt := range []struct {
		name string
		size int
	}{
		{name: "zero", size: 0},
		{name: "negative", size: -1},
	} {
		t.Run("default-buffer-"+tt.name, func(t *testing.T) {
			c := &Client{chanReg: newChanRegistry()}
			ch, cancel := h.subscribe(c, ChannelConfig{BufferSize: tt.size})
			if got := cap(ch); got != defaultChannelBuffer {
				t.Fatalf("default channel capacity = %d, want %d", got, defaultChannelBuffer)
			}
			cfg := h.config(c)
			if cfg == nil || cfg.bufSize != defaultChannelBuffer || cfg.policy != DropBlock {
				t.Fatalf("default channel config = %+v, want DropBlock/%d", cfg, defaultChannelBuffer)
			}
			cancel()
		})
	}
	t.Run("invalid-policy", func(t *testing.T) {
		assertStreamHardeningBlockingPolicy(t, h, DropPolicy(99))
	})
}

func TestStreamHardeningChannelPoliciesForEveryTopic(t *testing.T) {
	quote := streamHardeningChannelHarness[*marketdatav1.Quote]{
		subscribe: func(c *Client, cfg ChannelConfig) (<-chan *marketdatav1.Quote, func()) {
			return c.SubscribeQuoteChan(cfg)
		},
		emit:           func(c *Client, msg *marketdatav1.Quote) { c.emitQuote(msg) },
		config:         streamHardeningQuoteConfig,
		lifecycle:      streamHardeningQuoteLifecycle,
		lifecycleOrder: streamHardeningQuoteLifecycleByOrder,
		first:          &marketdatav1.Quote{Basic: &marketdatav1.Basic{Symbol: "first"}},
		second:         &marketdatav1.Quote{Basic: &marketdatav1.Basic{Symbol: "second"}},
		third:          &marketdatav1.Quote{Basic: &marketdatav1.Basic{Symbol: "third"}},
	}
	snapshot := streamHardeningChannelHarness[*marketdatav1.Snapshot]{
		subscribe: func(c *Client, cfg ChannelConfig) (<-chan *marketdatav1.Snapshot, func()) {
			return c.SubscribeSnapshotChan(cfg)
		},
		emit:           func(c *Client, msg *marketdatav1.Snapshot) { c.emitSnapshot(msg) },
		config:         streamHardeningSnapshotConfig,
		lifecycle:      streamHardeningSnapshotLifecycle,
		lifecycleOrder: streamHardeningSnapshotLifecycleByOrder,
		first:          &marketdatav1.Snapshot{Basic: &marketdatav1.Basic{Symbol: "first"}},
		second:         &marketdatav1.Snapshot{Basic: &marketdatav1.Basic{Symbol: "second"}},
		third:          &marketdatav1.Snapshot{Basic: &marketdatav1.Basic{Symbol: "third"}},
	}
	tick := streamHardeningChannelHarness[*marketdatav1.Tick]{
		subscribe: func(c *Client, cfg ChannelConfig) (<-chan *marketdatav1.Tick, func()) {
			return c.SubscribeTickChan(cfg)
		},
		emit:           func(c *Client, msg *marketdatav1.Tick) { c.emitTick(msg) },
		config:         streamHardeningTickConfig,
		lifecycle:      streamHardeningTickLifecycle,
		lifecycleOrder: streamHardeningTickLifecycleByOrder,
		first:          &marketdatav1.Tick{Basic: &marketdatav1.Basic{Symbol: "first"}},
		second:         &marketdatav1.Tick{Basic: &marketdatav1.Basic{Symbol: "second"}},
		third:          &marketdatav1.Tick{Basic: &marketdatav1.Basic{Symbol: "third"}},
	}

	t.Run("quote", func(t *testing.T) {
		assertStreamHardeningAllPolicies(t, quote)
		assertStreamHardeningDefaultsAndInvalidPolicy(t, quote)
	})
	t.Run("snapshot", func(t *testing.T) {
		assertStreamHardeningAllPolicies(t, snapshot)
		assertStreamHardeningDefaultsAndInvalidPolicy(t, snapshot)
	})
	t.Run("tick", func(t *testing.T) {
		assertStreamHardeningAllPolicies(t, tick)
		assertStreamHardeningDefaultsAndInvalidPolicy(t, tick)
	})
}

func assertStreamHardeningHeadOfLine[T any](t *testing.T, h streamHardeningChannelHarness[T]) {
	c := &Client{chanReg: newChanRegistry()}
	blocked, cancelBlocked := h.subscribe(c, ChannelConfig{Policy: DropBlock, BufferSize: 1})
	healthy, cancelHealthy := h.subscribe(c, ChannelConfig{Policy: DropOldest, BufferSize: 1})
	h.emit(c, h.first)
	if got := receiveStreamHardeningValue(t, healthy); !reflect.DeepEqual(got, h.first) {
		t.Fatalf("later subscriber initial value = %v, want first", got)
	}

	dispatchStarted := make(chan struct{})
	dispatchDone := make(chan struct{})
	go func() {
		close(dispatchStarted)
		h.emit(c, h.second)
		close(dispatchDone)
	}()
	<-dispatchStarted
	waitForStreamHardeningSend(t, h.lifecycleOrder(c, 0))
	select {
	case msg := <-healthy:
		t.Fatalf("later subscriber received %v before blocked subscriber was cancelled", msg)
	default:
	}
	select {
	case <-dispatchDone:
		t.Fatal("multiple-subscriber dispatch completed while the first subscriber was blocked")
	default:
	}

	cancelBlocked()
	select {
	case <-dispatchDone:
	case <-time.After(time.Second):
		t.Fatal("dispatch did not continue after blocked subscriber cancellation")
	}
	if got := receiveStreamHardeningValue(t, healthy); !reflect.DeepEqual(got, h.second) {
		t.Fatalf("later subscriber value = %v, want second", got)
	}
	if got := h.config(c).DropCount(); got != 0 {
		t.Fatalf("later subscriber drop count = %d, want 0", got)
	}
	cancelHealthy()
	assertStreamHardeningChannelClosed(t, blocked)
}

func TestStreamHardeningSubscriberDispatchPreservesOrderAndBlocksHeadOfLine(t *testing.T) {
	quote := streamHardeningChannelHarness[*marketdatav1.Quote]{
		subscribe: func(c *Client, cfg ChannelConfig) (<-chan *marketdatav1.Quote, func()) {
			return c.SubscribeQuoteChan(cfg)
		},
		emit:           func(c *Client, msg *marketdatav1.Quote) { c.emitQuote(msg) },
		config:         streamHardeningQuoteConfig,
		lifecycleOrder: streamHardeningQuoteLifecycleByOrder,
		first:          &marketdatav1.Quote{},
		second:         &marketdatav1.Quote{},
	}
	snapshot := streamHardeningChannelHarness[*marketdatav1.Snapshot]{
		subscribe: func(c *Client, cfg ChannelConfig) (<-chan *marketdatav1.Snapshot, func()) {
			return c.SubscribeSnapshotChan(cfg)
		},
		emit:           func(c *Client, msg *marketdatav1.Snapshot) { c.emitSnapshot(msg) },
		config:         streamHardeningSnapshotConfig,
		lifecycleOrder: streamHardeningSnapshotLifecycleByOrder,
		first:          &marketdatav1.Snapshot{},
		second:         &marketdatav1.Snapshot{},
	}
	tick := streamHardeningChannelHarness[*marketdatav1.Tick]{
		subscribe: func(c *Client, cfg ChannelConfig) (<-chan *marketdatav1.Tick, func()) {
			return c.SubscribeTickChan(cfg)
		},
		emit:           func(c *Client, msg *marketdatav1.Tick) { c.emitTick(msg) },
		config:         streamHardeningTickConfig,
		lifecycleOrder: streamHardeningTickLifecycleByOrder,
		first:          &marketdatav1.Tick{},
		second:         &marketdatav1.Tick{},
	}
	t.Run("quote", func(t *testing.T) { assertStreamHardeningHeadOfLine(t, quote) })
	t.Run("snapshot", func(t *testing.T) { assertStreamHardeningHeadOfLine(t, snapshot) })
	t.Run("tick", func(t *testing.T) { assertStreamHardeningHeadOfLine(t, tick) })
}

func TestStreamHardeningSubscribeChannelsAfterClose(t *testing.T) {
	c := &Client{chanReg: newChanRegistry()}
	if err := c.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	quote, cancelQuote := c.SubscribeQuoteChan(ChannelConfig{})
	snapshot, cancelSnapshot := c.SubscribeSnapshotChan(ChannelConfig{})
	tick, cancelTick := c.SubscribeTickChan(ChannelConfig{})
	assertStreamHardeningChannelClosed(t, quote)
	assertStreamHardeningChannelClosed(t, snapshot)
	assertStreamHardeningChannelClosed(t, tick)
	cancelQuote()
	cancelQuote()
	cancelSnapshot()
	cancelSnapshot()
	cancelTick()
	cancelTick()
}

func assertStreamHardeningConcurrentDispatchCancelClose[T any](t *testing.T, h streamHardeningChannelHarness[T]) {
	for iteration := 0; iteration < 32; iteration++ {
		c := &Client{chanReg: newChanRegistry()}
		ch, cancel := h.subscribe(c, ChannelConfig{Policy: DropBlock, BufferSize: 1})
		h.emit(c, h.first)
		start := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(5)
		go func() {
			defer wg.Done()
			<-start
			h.emit(c, h.second)
		}()
		go func() {
			defer wg.Done()
			<-start
			h.emit(c, h.third)
		}()
		go func() {
			defer wg.Done()
			<-start
			cancel()
			cancel()
		}()
		go func() {
			defer wg.Done()
			<-start
			_ = c.Close()
			_ = c.Close()
		}()
		go func() {
			defer wg.Done()
			close(start)
		}()
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatalf("iteration %d leaked a blocked lifecycle operation", iteration)
		}
		assertStreamHardeningChannelClosed(t, ch)
	}
}

func TestStreamHardeningConcurrentDispatchCancelAndTerminalClose(t *testing.T) {
	quote := streamHardeningChannelHarness[*marketdatav1.Quote]{
		subscribe: func(c *Client, cfg ChannelConfig) (<-chan *marketdatav1.Quote, func()) {
			return c.SubscribeQuoteChan(cfg)
		},
		emit:   func(c *Client, msg *marketdatav1.Quote) { c.emitQuote(msg) },
		first:  &marketdatav1.Quote{},
		second: &marketdatav1.Quote{},
		third:  &marketdatav1.Quote{},
	}
	snapshot := streamHardeningChannelHarness[*marketdatav1.Snapshot]{
		subscribe: func(c *Client, cfg ChannelConfig) (<-chan *marketdatav1.Snapshot, func()) {
			return c.SubscribeSnapshotChan(cfg)
		},
		emit:   func(c *Client, msg *marketdatav1.Snapshot) { c.emitSnapshot(msg) },
		first:  &marketdatav1.Snapshot{},
		second: &marketdatav1.Snapshot{},
		third:  &marketdatav1.Snapshot{},
	}
	tick := streamHardeningChannelHarness[*marketdatav1.Tick]{
		subscribe: func(c *Client, cfg ChannelConfig) (<-chan *marketdatav1.Tick, func()) {
			return c.SubscribeTickChan(cfg)
		},
		emit:   func(c *Client, msg *marketdatav1.Tick) { c.emitTick(msg) },
		first:  &marketdatav1.Tick{},
		second: &marketdatav1.Tick{},
		third:  &marketdatav1.Tick{},
	}
	t.Run("quote", func(t *testing.T) { assertStreamHardeningConcurrentDispatchCancelClose(t, quote) })
	t.Run("snapshot", func(t *testing.T) { assertStreamHardeningConcurrentDispatchCancelClose(t, snapshot) })
	t.Run("tick", func(t *testing.T) { assertStreamHardeningConcurrentDispatchCancelClose(t, tick) })
}
