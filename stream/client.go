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
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shing1211/webullapi4go/client"
	marketdatav1 "github.com/shing1211/webullapi4go/gen/webull/marketdata/v1"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/pkg/observability"
	mqtt "github.com/shing1211/webullapi4go/pkg/transport/mqtt"

	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// password is sent as the MQTT password. Webull documents that any value is
// accepted; the App Key is carried as the user name.
const password = "webullapi4go"

// streamMetrics holds the OTel instruments for stream telemetry.
// All fields are nil when no Meter is configured.
type streamMetrics struct {
	reconnectCounter    metric.Int64Counter
	quoteDropCounter    metric.Int64Counter
	snapshotDropCounter metric.Int64Counter
	tickDropCounter     metric.Int64Counter
}

// Client is a Webull market-data streaming client. It owns an MQTT connection
// and routes decoded pushes to the handlers registered with the On* methods.
//
// A Client is safe for concurrent use. Handlers may be registered before or
// after [Client.Connect]. The underlying [client.Client] is owned by the
// caller and is not closed by [Client.Close].
type Client struct {
	core *client.Client
	cfg  config
	mqtt *mqtt.Client

	closeOnce sync.Once

	// subs is the active-subscription registry used to re-issue HTTP
	// subscriptions after a reconnect.
	subs subscriptionRegistry
	// resubMu serialises re-subscription so concurrent reconnects cannot
	// interleave their subscribe calls.
	resubMu sync.Mutex
	// everConnected distinguishes the first connection from a reconnect.
	everConnected atomic.Bool
	// reconnecting mirrors the MQTT transport reconnect state.
	reconnecting atomic.Bool
	// resubCtx bounds background re-subscription; it is cancelled by Close.
	resubCtx    context.Context
	resubCancel context.CancelFunc

	// connCtx bounds the health watchdog; cancelled by Close.
	connCtx    context.Context
	connCancel context.CancelFunc
	healthOnce sync.Once
	// state is the current connection lifecycle state.
	state atomic.Value
	// lastMessageAt records the arrival time of the most recent data message
	// (quote/snapshot/tick). It is used by the health watchdog.
	lastMessageAt atomic.Int64

	// chanReg is the per-subscription channel registry.
	chanReg *chanRegistry
	// metrics holds OTel instruments. nil when no Meter is configured.
	metrics *streamMetrics
	// obs holds the telemetry configuration inherited from the core client.
	obs *observability.Config
	// tracer creates dispatch spans when tracing is enabled.
	tracer trace.Tracer

	mu             sync.RWMutex
	onQuote        []func(*marketdatav1.Quote)
	onSnapshot     []func(*marketdatav1.Snapshot)
	onTick         []func(*marketdatav1.Tick)
	onNotice       []func([]byte)
	onError        []func(error)
	onConnect      []func()
	onDisconnect   []func(error)
	onReconnecting []func()
	onStateChange  []func(State, State)
}

// New returns a streaming client bound to cl. When no session id is supplied
// with [WithSessionID] or [WithClientID], New generates a unique one. New does
// not open a connection; call [Client.Connect] for that.
//
// The broker address is taken from cl's resolved endpoints unless overridden
// with [WithMQTTURL]; by default the plain TCP endpoint is used and
// [WithWebSocket] selects the WebSocket endpoint instead.
func New(cl *client.Client, opts ...Option) (*Client, error) {
	if cl == nil {
		return nil, errs.New(errs.CodeInvalidConfig, "stream: client is required")
	}
	cfg := defaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.sessionID == "" {
		cfg.sessionID = newSessionID()
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	broker, err := cfg.resolveBroker(cl)
	if err != nil {
		return nil, err
	}

	mc, err := mqtt.New(mqtt.Config{
		Broker:               broker,
		ClientID:             cfg.sessionID,
		Username:             cl.Config().AppKey,
		Password:             password,
		KeepAlive:            cfg.keepAlive,
		ConnectTimeout:       cfg.connectTimeout,
		WriteTimeout:         cfg.writeTimeout,
		MessageChannelDepth:  cfg.messageChannelDepth,
		CleanSession:         cfg.cleanSession,
		AutoReconnect:        cfg.autoReconnect,
		MaxReconnectInterval: DefaultMaxReconnectInterval,
		TLSConfig:            cfg.tlsConfig,
	})
	if err != nil {
		return nil, errs.Wrap(errs.CodeInvalidConfig, "stream: invalid mqtt configuration", err)
	}

	obs := cl.ObservabilityConfig()
	if cfg.meter == nil && obs != nil {
		cfg.meter = obs.Meter("webullapi4go/stream")
	}
	c := &Client{core: cl, cfg: cfg, mqtt: mc, obs: obs}
	if obs != nil {
		c.tracer = obs.Tracer("webullapi4go/stream")
	}
	c.resubCtx, c.resubCancel = context.WithCancel(context.Background())
	c.connCtx, c.connCancel = context.WithCancel(context.Background())
	c.state.Store(StateDisconnected)
	c.chanReg = newChanRegistry()
	c.metrics = newStreamMetrics(cfg.meter)
	mc.SetMessageHandler(c.handleMessage)
	mc.SetConnectHandler(c.handleConnect)
	mc.SetConnectionLostHandler(c.handleConnectionLost)
	mc.SetReconnectHandler(c.handleReconnecting)
	mc.SetErrorHandler(c.emitError)
	return c, nil
}

// resolveBroker returns the broker URL for cfg, preferring an explicit
// [WithMQTTURL] override over the endpoints resolved by the core client.
func (cfg config) resolveBroker(cl *client.Client) (string, error) {
	if cfg.mqttURL != "" {
		return cfg.mqttURL, nil
	}
	endpoints := cl.Endpoints()
	if cfg.websocket {
		if endpoints.MQTTWebSocket == "" {
			return "", errs.New(errs.CodeInvalidConfig, "stream: no MQTT WebSocket endpoint configured")
		}
		return endpoints.MQTTWebSocket, nil
	}
	if endpoints.MQTT == "" {
		return "", errs.New(errs.CodeInvalidConfig, "stream: no MQTT endpoint configured")
	}
	return endpoints.MQTT, nil
}

// Connect establishes the MQTT connection, blocking until it succeeds, fails,
// or ctx is done. The first successful connection fires the [Client.OnConnect]
// handlers.
//
// Webull does not restore subscriptions after a connection is lost. When a
// connection is re-established (see [WithAutoReconnect]), the client
// automatically re-issues the HTTP subscribe calls for every active
// subscription before the [Client.OnConnect] handlers run, so callers do not
// need to subscribe again.
//
// A failure caused by Webull code 105 (the App Key already has five concurrent
// connections) is reported with a message explaining the limit and the roughly
// one-minute server retention window.
func (c *Client) Connect(ctx context.Context) error {
	if c.mqtt == nil {
		return errs.New(errs.CodeInvalidConfig, "stream: client is not initialized")
	}
	if c.State() == StateClosed {
		return errs.New(errs.CodeInvalidConfig, "stream: client is closed")
	}
	if c.IsConnected() {
		return nil
	}
	c.setState(StateConnecting)
	if err := c.mqtt.Connect(ctx); err != nil {
		if c.State() != StateClosed {
			c.setState(StateDisconnected)
		}
		return streamConnectError(err)
	}
	if c.cfg.healthWatchdogInterval > 0 {
		c.healthOnce.Do(func() { go c.healthWatchdog() })
	}
	return nil
}

// streamConnectError wraps a transport connect failure with a clear message for
// the Webull-specific cases.
func streamConnectError(err error) error {
	switch {
	case errors.Is(err, mqtt.ErrConnectionLimit):
		return errs.Wrap(errs.CodeTransport,
			"stream: connect rejected: exceeds the 5 concurrent connections per App Key (Webull code 105); wait about 1 minute after a disconnect before reconnecting",
			err)
	case errors.Is(err, mqtt.ErrConnectionRefused):
		return errs.Wrap(errs.CodeTransport,
			"stream: connect refused by broker; verify the App Key and that the session id is not already in use",
			err)
	default:
		return errs.Wrap(errs.CodeTransport, "stream: connect", err)
	}
}

// Close disconnects from the broker. It is idempotent and does not close the
// underlying [client.Client].
func (c *Client) Close() error {
	c.closeOnce.Do(func() {
		c.reconnecting.Store(false)
		c.setState(StateClosed)
		if c.connCancel != nil {
			c.connCancel()
		}
		if c.resubCancel != nil {
			c.resubCancel()
		}
		if c.chanReg != nil {
			c.chanReg.stopAll()
		}
		if c.mqtt != nil {
			_ = c.mqtt.Close()
		}
	})
	return nil
}

// IsConnected reports whether the MQTT connection is currently live. It is
// false while a reconnect is in progress.
func (c *Client) IsConnected() bool {
	return c.mqtt != nil && c.mqtt.IsConnected() && !c.Reconnecting()
}

// Reconnecting reports whether the client is currently attempting to
// re-establish a connection that was lost. It is always false before the first
// connection and after [Client.Close].
func (c *Client) Reconnecting() bool {
	return c.reconnecting.Load()
}

// SessionID returns the session id used as the MQTT client id and in
// subscribe/unsubscribe calls. It is stable for the lifetime of the Client.
func (c *Client) SessionID() string { return c.cfg.sessionID }

// OnQuote registers a handler for order-book pushes on the quote topic.
func (c *Client) OnQuote(fn func(*marketdatav1.Quote)) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onQuote = append(c.onQuote, fn)
	c.mu.Unlock()
}

// OnSnapshot registers a handler for market-snapshot pushes.
func (c *Client) OnSnapshot(fn func(*marketdatav1.Snapshot)) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onSnapshot = append(c.onSnapshot, fn)
	c.mu.Unlock()
}

// OnTick registers a handler for tick-by-tick pushes.
func (c *Client) OnTick(fn func(*marketdatav1.Tick)) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onTick = append(c.onTick, fn)
	c.mu.Unlock()
}

// OnNotice registers a handler for server notifications. The payload is
// delivered as raw JSON because the notice topic is not protobuf encoded.
func (c *Client) OnNotice(fn func([]byte)) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onNotice = append(c.onNotice, fn)
	c.mu.Unlock()
}

// OnError registers a handler for asynchronous errors, such as a lost
// connection, a failed decode, or an unknown topic.
func (c *Client) OnError(fn func(error)) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onError = append(c.onError, fn)
	c.mu.Unlock()
}

// OnConnect registers a handler invoked after a successful connection,
// including automatic reconnections.
func (c *Client) OnConnect(fn func()) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onConnect = append(c.onConnect, fn)
	c.mu.Unlock()
}

// OnDisconnect registers a handler invoked when an established connection is
// lost, with the reason.
func (c *Client) OnDisconnect(fn func(error)) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onDisconnect = append(c.onDisconnect, fn)
	c.mu.Unlock()
}

// OnReconnecting registers a handler invoked when the client begins attempting
// to reconnect after a connection loss.
func (c *Client) OnReconnecting(fn func()) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onReconnecting = append(c.onReconnecting, fn)
	c.mu.Unlock()
}

// OnStateChange registers a handler invoked whenever the connection state
// changes, receiving the previous and next state.
func (c *Client) OnStateChange(fn func(State, State)) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onStateChange = append(c.onStateChange, fn)
	c.mu.Unlock()
}

// State returns the current connection state.
func (c *Client) State() State {
	if c.state.Load() == nil {
		c.state.Store(StateDisconnected)
	}
	value, _ := c.state.Load().(State)
	return value
}

// setState transitions to next, firing OnStateChange handlers outside the lock.
func (c *Client) setState(next State) {
	c.State()
	var prev State
	for {
		prev, _ = c.state.Load().(State)
		if prev == StateClosed && next != StateClosed {
			return
		}
		if c.state.CompareAndSwap(prev, next) {
			break
		}
	}
	if prev != next {
		c.mu.RLock()
		handlers := make([]func(State, State), len(c.onStateChange))
		copy(handlers, c.onStateChange)
		c.mu.RUnlock()
		for _, h := range handlers {
			h(prev, next)
		}
	}
}

// handleMessage is the low-level callback. It reports dispatch failures to the
// error handlers and records the arrival time for health tracking.
func (c *Client) handleMessage(m mqtt.Message) {
	if c.State() == StateClosed {
		return
	}
	var span trace.Span
	if c.tracer != nil {
		_, span = c.tracer.Start(context.Background(), "mqtt.dispatch",
			trace.WithSpanKind(trace.SpanKindConsumer),
			trace.WithAttributes(
				attribute.String("messaging.system", "mqtt"),
				attribute.String("messaging.destination", streamTopicKind(m.Topic)),
				attribute.Int("messaging.message.body.size", len(m.Payload)),
			),
		)
		defer span.End()
	}
	if isDataTopic(m.Topic) {
		c.markDataMessage()
	}
	if err := c.dispatch(m.Topic, m.Payload); err != nil {
		if span != nil {
			span.RecordError(err)
			span.SetStatus(otelcodes.Error, "stream message dispatch failed")
		}
		c.emitError(err)
	}
}

func isDataTopic(topic string) bool {
	switch normalizeTopic(topic) {
	case TopicQuote, TopicSnapshot, TopicTick:
		return true
	default:
		return false
	}
}

func streamTopicKind(topic string) string {
	switch normalizeTopic(topic) {
	case TopicQuote:
		return "quote"
	case TopicSnapshot:
		return "snapshot"
	case TopicTick:
		return "tick"
	case TopicNotice:
		return "notice"
	case TopicEcho:
		return "echo"
	default:
		return "unknown"
	}
}

func (c *Client) markDataMessage() {
	if c.State() == StateClosed {
		return
	}
	c.lastMessageAt.Store(time.Now().UnixMilli())
	if c.State() == StateDegraded {
		c.setState(StateConnected)
	}
}

// ChannelConfig configures a channel-based subscription. The default is a
// blocking policy with a buffer of 100 messages.
type ChannelConfig struct {
	// Policy controls what happens when the buffer is full. The default is
	// [DropBlock]. Use [DropOldest] to discard the oldest unread message or
	// [DropSample] to probabilistically discard before the buffer is full.
	Policy DropPolicy
	// BufferSize sets the channel buffer depth. Values ≤ 0 use the default of 100.
	BufferSize int
}

func (c ChannelConfig) resolve() (policy DropPolicy, bufSize int) {
	policy = c.Policy
	bufSize = c.BufferSize
	if bufSize <= 0 {
		bufSize = defaultChannelBuffer
	}
	return
}

// SubscribeQuoteChan returns a channel that receives decoded quote pushes.
// The channel is closed when the subscription is cancelled. The returned
// cancel function must be called to release the subscription.
func (c *Client) SubscribeQuoteChan(cfg ChannelConfig) (<-chan *marketdatav1.Quote, func()) {
	policy, bufSize := cfg.resolve()
	ch := make(chan *marketdatav1.Quote, bufSize)
	entry := newChanQuote(ch)
	otelCounter := metric.Int64Counter(nil)
	if c.metrics != nil {
		otelCounter = c.metrics.quoteDropCounter
	}
	config := &channelConfig{policy: policy, bufSize: bufSize, otelCounter: otelCounter, topic: "quote"}
	if !c.chanReg.addQuote(entry, config) {
		entry.shutdown()
	}
	return ch, func() { c.chanReg.removeQuote(entry) }
}

// SubscribeSnapshotChan returns a channel that receives decoded snapshot pushes.
// The channel is closed when the subscription is cancelled. The returned
// cancel function must be called to release the subscription.
func (c *Client) SubscribeSnapshotChan(cfg ChannelConfig) (<-chan *marketdatav1.Snapshot, func()) {
	policy, bufSize := cfg.resolve()
	ch := make(chan *marketdatav1.Snapshot, bufSize)
	entry := newChanSnapshot(ch)
	otelCounter := metric.Int64Counter(nil)
	if c.metrics != nil {
		otelCounter = c.metrics.snapshotDropCounter
	}
	config := &channelConfig{policy: policy, bufSize: bufSize, otelCounter: otelCounter, topic: "snapshot"}
	if !c.chanReg.addSnapshot(entry, config) {
		entry.shutdown()
	}
	return ch, func() { c.chanReg.removeSnapshot(entry) }
}

// SubscribeTickChan returns a channel that receives decoded tick pushes.
// The channel is closed when the subscription is cancelled. The returned
// cancel function must be called to release the subscription.
func (c *Client) SubscribeTickChan(cfg ChannelConfig) (<-chan *marketdatav1.Tick, func()) {
	policy, bufSize := cfg.resolve()
	ch := make(chan *marketdatav1.Tick, bufSize)
	entry := newChanTick(ch)
	otelCounter := metric.Int64Counter(nil)
	if c.metrics != nil {
		otelCounter = c.metrics.tickDropCounter
	}
	config := &channelConfig{policy: policy, bufSize: bufSize, otelCounter: otelCounter, topic: "tick"}
	if !c.chanReg.addTick(entry, config) {
		entry.shutdown()
	}
	return ch, func() { c.chanReg.removeTick(entry) }
}

// dispatch routes a raw message by topic and invokes the matching typed
// handlers. It returns an error for an unknown topic or a malformed payload;
// the echo heartbeat is silently ignored.
func (c *Client) dispatch(topic string, payload []byte) error {
	switch normalizeTopic(topic) {
	case TopicQuote:
		msg, err := marketdatav1.DecodeQuote(payload)
		if err != nil {
			return err
		}
		c.emitQuote(msg)
	case TopicSnapshot:
		msg, err := marketdatav1.DecodeSnapshot(payload)
		if err != nil {
			return err
		}
		c.emitSnapshot(msg)
	case TopicTick:
		msg, err := marketdatav1.DecodeTick(payload)
		if err != nil {
			return err
		}
		c.emitTick(msg)
	case TopicNotice:
		c.emitNotice(payload)
	case TopicEcho:
		// Heartbeat; no payload and no handler.
	case "":
		return errors.New("stream: message with empty topic")
	default:
		return errors.New("stream: unknown topic " + topic)
	}
	return nil
}

// emitQuote invokes every registered quote handler and channel subscriber
// outside the lock.
func (c *Client) emitQuote(msg *marketdatav1.Quote) {
	c.mu.RLock()
	handlers := make([]func(*marketdatav1.Quote), len(c.onQuote))
	copy(handlers, c.onQuote)
	c.mu.RUnlock()
	for _, h := range handlers {
		h(msg)
	}
	if c.chanReg != nil {
		c.chanReg.dispatchQuote(msg)
	}
}

// emitSnapshot invokes every registered snapshot handler and channel subscriber
// outside the lock.
func (c *Client) emitSnapshot(msg *marketdatav1.Snapshot) {
	c.mu.RLock()
	handlers := make([]func(*marketdatav1.Snapshot), len(c.onSnapshot))
	copy(handlers, c.onSnapshot)
	c.mu.RUnlock()
	for _, h := range handlers {
		h(msg)
	}
	if c.chanReg != nil {
		c.chanReg.dispatchSnapshot(msg)
	}
}

// emitTick invokes every registered tick handler and channel subscriber
// outside the lock.
func (c *Client) emitTick(msg *marketdatav1.Tick) {
	c.mu.RLock()
	handlers := make([]func(*marketdatav1.Tick), len(c.onTick))
	copy(handlers, c.onTick)
	c.mu.RUnlock()
	for _, h := range handlers {
		h(msg)
	}
	if c.chanReg != nil {
		c.chanReg.dispatchTick(msg)
	}
}

// emitNotice invokes every registered notice handler outside the lock.
func (c *Client) emitNotice(payload []byte) {
	c.mu.RLock()
	handlers := make([]func([]byte), len(c.onNotice))
	copy(handlers, c.onNotice)
	c.mu.RUnlock()
	for _, h := range handlers {
		h(payload)
	}
}

// emitError invokes every registered error handler outside the lock.
func (c *Client) emitError(err error) {
	if err == nil || c.State() == StateClosed {
		return
	}
	c.mu.RLock()
	handlers := make([]func(error), len(c.onError))
	copy(handlers, c.onError)
	c.mu.RUnlock()
	for _, h := range handlers {
		h(err)
	}
}

func (c *Client) logStreamEvent(level slog.Level, message string) {
	if c.obs == nil || c.obs.Logger == nil {
		return
	}
	c.obs.Logger.LogAttrs(context.Background(), level, message,
		slog.String("transport", "mqtt"),
		slog.String("state", c.State().String()),
	)
}

// handleConnect is the MQTT connect callback. It runs on the initial
// connection and on every successful reconnect. On a reconnect it re-issues the
// active subscriptions before notifying connect handlers, so handlers observe a
// restored session.
func (c *Client) handleConnect() {
	if c.State() == StateClosed {
		return
	}
	c.reconnecting.Store(false)
	c.lastMessageAt.Store(time.Now().UnixMilli())
	c.setState(StateConnected)
	c.logStreamEvent(slog.LevelInfo, "webull stream connected")
	if c.everConnected.Swap(true) {
		c.resubscribe()
	}
	if c.State() == StateClosed {
		return
	}
	c.emitConnect()
}

// handleConnectionLost is the MQTT connection-lost callback.
func (c *Client) handleConnectionLost(err error) {
	if c.State() == StateClosed {
		return
	}
	if !c.cfg.autoReconnect {
		c.reconnecting.Store(false)
	}
	if !c.cfg.autoReconnect || !c.Reconnecting() {
		c.setState(StateDisconnected)
	}
	if c.State() == StateClosed {
		return
	}
	c.logStreamEvent(slog.LevelWarn, "webull stream disconnected")
	c.emitDisconnect(err)
}

// handleReconnecting is the MQTT reconnect-start callback.
func (c *Client) handleReconnecting() {
	if c.State() == StateClosed {
		return
	}
	c.reconnecting.Store(true)
	c.setState(StateReconnecting)
	// OTel metric recording is intentionally fire-and-forget: the metric SDK
	// records synchronously and cannot block, so context.Background() is safe.
	if c.metrics != nil && c.metrics.reconnectCounter != nil {
		c.metrics.reconnectCounter.Add(context.Background(), 1)
	}
	if c.State() == StateClosed {
		return
	}
	c.logStreamEvent(slog.LevelInfo, "webull stream reconnecting")
	c.mu.RLock()
	handlers := make([]func(), len(c.onReconnecting))
	copy(handlers, c.onReconnecting)
	c.mu.RUnlock()
	for _, h := range handlers {
		h()
	}
}

// resubscribe re-issues every active subscription through the core client. It
// is serialized and works from a deduplicated registry snapshot, so each active
// subscription is re-issued exactly once per reconnect and never duplicated.
// Failures are reported to the error handlers without aborting the remaining
// subscriptions.
func (c *Client) resubscribe() {
	if c.core == nil || !c.cfg.autoResubscribe {
		return
	}
	var resubErrs []error
	c.resubMu.Lock()
	func() {
		defer c.resubMu.Unlock()
		if c.State() == StateClosed {
			return
		}
		if c.subs.empty() {
			return
		}
		reqs := c.subs.requests()
		if len(reqs) == 0 {
			return
		}
		ctx, cancel := c.resubscribeContext()
		defer cancel()
		for _, req := range reqs {
			if err := c.subscribe(ctx, req); err != nil {
				resubErrs = append(resubErrs, err)
			}
		}
	}()
	if c.State() == StateClosed {
		return
	}
	for _, err := range resubErrs {
		c.emitError(err)
	}
}

// resubscribeContext returns a context bounded by the re-subscribe timeout and
// tied to the client lifetime.
func (c *Client) resubscribeContext() (context.Context, context.CancelFunc) {
	if c.resubCtx == nil {
		panic("stream: resubCtx is nil — New() was not called or context was already cancelled")
	}
	return context.WithTimeout(c.resubCtx, c.cfg.resubscribeTimeout)
}

// healthWatchdog monitors message-age and transitions to Degraded when no data
// arrives within the configured interval. It recovers to Connected when a
// message arrives. The watchdog is started by [Client.Connect] and cancelled by
// [Client.Close].
func (c *Client) healthWatchdog() {
	if c.connCtx == nil || c.cfg.healthWatchdogInterval <= 0 {
		return
	}
	ticker := time.NewTicker(c.cfg.healthWatchdogInterval)
	defer ticker.Stop()
	for {
		select {
		case <-c.connCtx.Done():
			return
		case <-ticker.C:
			c.checkHealth(time.Now())
		}
	}
}

func (c *Client) checkHealth(now time.Time) {
	interval := c.cfg.healthWatchdogInterval
	if interval <= 0 {
		return
	}
	last := c.lastMessageAt.Load()
	if last == 0 {
		return
	}
	age := now.Sub(time.UnixMilli(last))
	state := c.State()
	if state == StateConnected && age > interval {
		c.setState(StateDegraded)
		return
	}
	if state == StateDegraded && age <= interval {
		c.setState(StateConnected)
	}
}

// emitConnect invokes every registered connect handler outside the lock.
func (c *Client) emitConnect() {
	c.mu.RLock()
	handlers := make([]func(), len(c.onConnect))
	copy(handlers, c.onConnect)
	c.mu.RUnlock()
	for _, h := range handlers {
		h()
	}
}

// emitDisconnect invokes every registered disconnect handler outside the lock.
func (c *Client) emitDisconnect(err error) {
	c.mu.RLock()
	handlers := make([]func(error), len(c.onDisconnect))
	copy(handlers, c.onDisconnect)
	c.mu.RUnlock()
	for _, h := range handlers {
		h(err)
	}
}

// newSessionID returns a random RFC 4122 version 4 UUID string for use as an
// MQTT session id. crypto/rand is used so ids are not reused across calls.
func newSessionID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand failure is effectively unreachable; fall back to a
		// fixed-width hex encoding of an all-zero buffer rather than panic,
		// leaving uniqueness enforcement to the caller's retry.
		return "00000000-0000-4000-8000-000000000000"
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return hex.EncodeToString(b[0:4]) + "-" +
		hex.EncodeToString(b[4:6]) + "-" +
		hex.EncodeToString(b[6:8]) + "-" +
		hex.EncodeToString(b[8:10]) + "-" +
		hex.EncodeToString(b[10:16])
}

// newStreamMetrics creates OTel instruments for stream telemetry if a meter is
// supplied; otherwise it returns nil.
func newStreamMetrics(m metric.Meter) *streamMetrics {
	if m == nil {
		return nil
	}
	return &streamMetrics{
		reconnectCounter:    newCounter(m, "reconnects", "Stream reconnection events"),
		quoteDropCounter:    newCounter(m, "channel_drops", "Channel message drops"),
		snapshotDropCounter: newCounter(m, "channel_drops", "Channel message drops"),
		tickDropCounter:     newCounter(m, "channel_drops", "Channel message drops"),
	}
}

func newCounter(m metric.Meter, name, desc string) metric.Int64Counter {
	c, _ := m.Int64Counter(name, metric.WithDescription(desc))
	return c
}
