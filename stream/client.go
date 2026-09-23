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
	"sync"
	"sync/atomic"

	"github.com/shing1211/webullapi4go/client"
	marketdatav1 "github.com/shing1211/webullapi4go/gen/webull/marketdata/v1"
	"github.com/shing1211/webullapi4go/pkg/errors"
	mqtt "github.com/shing1211/webullapi4go/pkg/transport/mqtt"
)

// password is sent as the MQTT password. Webull documents that any value is
// accepted; the App Key is carried as the user name.
const password = "webullapi4go"

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

	mu           sync.RWMutex
	onQuote      []func(*marketdatav1.Quote)
	onSnapshot   []func(*marketdatav1.Snapshot)
	onTick       []func(*marketdatav1.Tick)
	onNotice     []func([]byte)
	onError      []func(error)
	onConnect    []func()
	onDisconnect []func(error)
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

	c := &Client{core: cl, cfg: cfg, mqtt: mc}
	c.resubCtx, c.resubCancel = context.WithCancel(context.Background())
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
	if err := c.mqtt.Connect(ctx); err != nil {
		return streamConnectError(err)
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
		if c.resubCancel != nil {
			c.resubCancel()
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
	return c.mqtt != nil && c.mqtt.IsConnected()
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

// handleMessage is the low-level callback. It reports dispatch failures to the
// error handlers.
func (c *Client) handleMessage(m mqtt.Message) {
	if err := c.dispatch(m.Topic, m.Payload); err != nil {
		c.emitError(err)
	}
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

// emitQuote invokes every registered quote handler outside the lock.
func (c *Client) emitQuote(msg *marketdatav1.Quote) {
	c.mu.RLock()
	handlers := make([]func(*marketdatav1.Quote), len(c.onQuote))
	copy(handlers, c.onQuote)
	c.mu.RUnlock()
	for _, h := range handlers {
		h(msg)
	}
}

// emitSnapshot invokes every registered snapshot handler outside the lock.
func (c *Client) emitSnapshot(msg *marketdatav1.Snapshot) {
	c.mu.RLock()
	handlers := make([]func(*marketdatav1.Snapshot), len(c.onSnapshot))
	copy(handlers, c.onSnapshot)
	c.mu.RUnlock()
	for _, h := range handlers {
		h(msg)
	}
}

// emitTick invokes every registered tick handler outside the lock.
func (c *Client) emitTick(msg *marketdatav1.Tick) {
	c.mu.RLock()
	handlers := make([]func(*marketdatav1.Tick), len(c.onTick))
	copy(handlers, c.onTick)
	c.mu.RUnlock()
	for _, h := range handlers {
		h(msg)
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
	if err == nil {
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

// handleConnect is the MQTT connect callback. It runs on the initial
// connection and on every successful reconnect. On a reconnect it re-issues the
// active subscriptions before notifying connect handlers, so handlers observe a
// restored session.
func (c *Client) handleConnect() {
	c.reconnecting.Store(false)
	if c.everConnected.Swap(true) {
		c.resubscribe()
	}
	c.emitConnect()
}

// handleConnectionLost is the MQTT connection-lost callback.
func (c *Client) handleConnectionLost(err error) {
	c.reconnecting.Store(false)
	c.emitDisconnect(err)
}

// handleReconnecting is the MQTT reconnect-start callback.
func (c *Client) handleReconnecting() {
	c.reconnecting.Store(true)
}

// resubscribe re-issues every active subscription through the core client. It
// is serialized and works from a deduplicated registry snapshot, so each active
// subscription is re-issued exactly once per reconnect and never duplicated.
// Failures are reported to the error handlers without aborting the remaining
// subscriptions.
func (c *Client) resubscribe() {
	if c.core == nil || !c.cfg.autoResubscribe || c.subs.empty() {
		return
	}
	reqs := c.subs.requests()
	if len(reqs) == 0 {
		return
	}
	c.resubMu.Lock()
	defer c.resubMu.Unlock()
	ctx, cancel := c.resubscribeContext()
	defer cancel()
	for _, req := range reqs {
		if err := c.Subscribe(ctx, req); err != nil {
			c.emitError(err)
		}
	}
}

// resubscribeContext returns a context bounded by the re-subscribe timeout and
// tied to the client lifetime.
func (c *Client) resubscribeContext() (context.Context, context.CancelFunc) {
	base := c.resubCtx
	if base == nil {
		base = context.Background()
	}
	return context.WithTimeout(base, c.cfg.resubscribeTimeout)
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
