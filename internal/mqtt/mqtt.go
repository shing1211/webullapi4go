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

package mqtt

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// Default connection parameters applied by [New] when the corresponding
// [Config] field is zero.
const (
	// DefaultKeepAlive is the MQTT keep-alive interval.
	DefaultKeepAlive = 30 * time.Second
	// DefaultConnectTimeout bounds a single connection attempt.
	DefaultConnectTimeout = 15 * time.Second
	// DefaultWriteTimeout bounds writing a control packet to the broker.
	DefaultWriteTimeout = 5 * time.Second
	// DefaultPingTimeout is how long to wait for a PINGRESP before treating
	// the connection as lost.
	DefaultPingTimeout = 10 * time.Second
	// DefaultMessageChannelDepth is the number of inbound messages buffered
	// before incoming traffic is throttled.
	DefaultMessageChannelDepth uint = 100
	// DefaultMaxReconnectInterval caps the exponential reconnect backoff when
	// auto-reconnect is enabled.
	DefaultMaxReconnectInterval = 30 * time.Second
)

// ConnackConnectionLimit is Webull's CONNACK return code reporting that the
// App Key already holds the maximum number of concurrent connections. MQTT
// 3.1.1 only standardises return codes 0-5, so paho does not translate this
// code and may complete a refused connection without an error; [Client.Connect]
// detects that case separately.
const ConnackConnectionLimit byte = 105

// Connection-level errors. They are matched with [errors.Is] through
// [ConnackError.Is], so callers can branch on the cause without parsing text.
var (
	// ErrConnectionRefused reports that the broker rejected the CONNECT
	// packet or the connection did not become usable.
	ErrConnectionRefused = errors.New("mqtt: connection refused by broker")
	// ErrConnectionLimit reports Webull error code 105: the App Key already
	// has five concurrent connections. The server retains a disconnected
	// session for about one minute, so a caller that has hit the limit must
	// wait before reconnecting.
	ErrConnectionLimit = errors.New("mqtt: exceeds connection limit (code 105): at most 5 concurrent connections per App Key; wait about 1 minute after a disconnect before reconnecting")
)

// ConnackError reports a broker rejection of the MQTT CONNECT packet. Code is
// the raw CONNACK return code, including Webull's non-standard 100-105 codes.
type ConnackError struct {
	// Code is the raw CONNACK return code.
	Code byte
	// Err is the optional underlying cause.
	Err error
}

// Error implements the error interface.
func (e *ConnackError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("mqtt: connection refused (code %d): %v", e.Code, e.Err)
	}
	return fmt.Sprintf("mqtt: connection refused (code %d)", e.Code)
}

// Unwrap returns the wrapped cause, if any.
func (e *ConnackError) Unwrap() error { return e.Err }

// Is reports whether e matches target. A [ConnackError] matches
// [ErrConnectionLimit] when its code is [ConnackConnectionLimit] and
// [ErrConnectionRefused] for any rejection.
func (e *ConnackError) Is(target error) bool {
	switch target {
	case ErrConnectionLimit:
		return e.Code == ConnackConnectionLimit
	case ErrConnectionRefused:
		return true
	default:
		return errors.Is(e.Err, target)
	}
}

// Message is a message received from the broker. It is a transport-level copy
// of the payload; topic routing and decoding happen in the caller.
type Message struct {
	// Topic is the MQTT topic the message was published on.
	Topic string
	// Payload is the raw message body.
	Payload []byte
	// QoS is the quality-of-service level the message was delivered with.
	QoS byte
	// Retained reports whether the message was a retained delivery.
	Retained bool
	// Duplicate reports whether the broker flagged the delivery as a
	// duplicate.
	Duplicate bool
}

// Handler receives a decoded transport [Message].
type Handler func(Message)

// Config configures an MQTT [Client]. Only Broker and ClientID are required;
// zero-valued timing and buffer fields are replaced by the package defaults.
type Config struct {
	// Broker is the broker address. It accepts a paho URL
	// (tcp://host:port, ssl://, tls://, ws://, wss://) or a bare host:port,
	// which is treated as tcp://host:port.
	Broker string
	// ClientID is the MQTT client identifier. Webull requires this to be the
	// session id also used for subscribe/unsubscribe.
	ClientID string
	// Username is the MQTT user name; Webull expects the App Key.
	Username string
	// Password is the MQTT password; Webull accepts any value.
	Password string
	// KeepAlive is the keep-alive interval.
	KeepAlive time.Duration
	// ConnectTimeout bounds a single connection attempt.
	ConnectTimeout time.Duration
	// WriteTimeout bounds writing a control packet.
	WriteTimeout time.Duration
	// PingTimeout is how long to wait for a PINGRESP.
	PingTimeout time.Duration
	// MessageChannelDepth is the inbound message buffer size.
	MessageChannelDepth uint
	// CleanSession requests a clean MQTT session. Webull connections are
	// stateless, so this is normally true.
	CleanSession bool
	// AutoReconnect enables the client's background reconnect loop.
	AutoReconnect bool
	// MaxReconnectInterval caps the reconnect backoff.
	MaxReconnectInterval time.Duration
	// TLSConfig overrides the TLS configuration used for tls/ssl/wss
	// connections. A nil value lets the paho client use a default
	// configuration.
	TLSConfig *tls.Config
}

// withDefaults returns a copy of c with zero-valued fields replaced by the
// package defaults.
func (c Config) withDefaults() Config {
	if c.KeepAlive <= 0 {
		c.KeepAlive = DefaultKeepAlive
	}
	if c.ConnectTimeout <= 0 {
		c.ConnectTimeout = DefaultConnectTimeout
	}
	if c.WriteTimeout <= 0 {
		c.WriteTimeout = DefaultWriteTimeout
	}
	if c.PingTimeout <= 0 {
		c.PingTimeout = DefaultPingTimeout
	}
	if c.MessageChannelDepth == 0 {
		c.MessageChannelDepth = DefaultMessageChannelDepth
	}
	if c.MaxReconnectInterval <= 0 {
		c.MaxReconnectInterval = DefaultMaxReconnectInterval
	}
	return c
}

// Validate reports whether c is usable. It does not make any network calls.
func (c Config) Validate() error {
	if c.Broker == "" {
		return errors.New("mqtt: broker is required")
	}
	if c.ClientID == "" {
		return errors.New("mqtt: client id is required")
	}
	if _, err := normalizeBroker(c.Broker); err != nil {
		return err
	}
	return nil
}

// normalizeBroker validates a broker address and returns it as a paho URL
// string. A bare host:port is prefixed with tcp://.
func normalizeBroker(broker string) (string, error) {
	raw := broker
	if !hasScheme(broker) {
		raw = "tcp://" + broker
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", errors.New("mqtt: invalid broker address: " + err.Error())
	}
	switch u.Scheme {
	case "tcp", "ssl", "tls", "ws", "wss":
	default:
		return "", errors.New("mqtt: unsupported broker scheme: " + u.Scheme)
	}
	if u.Host == "" {
		return "", errors.New("mqtt: broker address has no host")
	}
	return raw, nil
}

// hasScheme reports whether s begins with a URL scheme such as "tcp://".
func hasScheme(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9', c == '+', c == '-', c == '.':
			continue
		case c == ':':
			return i+2 < len(s) && s[i+1] == '/' && s[i+2] == '/'
		default:
			return false
		}
	}
	return false
}

// Client is a thin wrapper over a paho MQTT client. It exposes only the
// operations Webull streaming needs and keeps paho out of exported
// signatures. A Client is safe for concurrent use.
type Client struct {
	mu  sync.RWMutex
	cfg Config
	pc  paho.Client

	reconnecting atomic.Bool

	onMessage        Handler
	onConnect        func()
	onConnectionLost func(error)
	onReconnect      func()
	onError          func(error)
}

// New returns a client configured by cfg. It does not open a connection;
// call [Client.Connect] for that. New validates cfg and returns an error when
// the broker address or client id is missing or malformed.
func New(cfg Config) (*Client, error) {
	cfg = cfg.withDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	c := &Client{cfg: cfg}

	opts := paho.NewClientOptions()
	opts.AddBroker(cfg.Broker)
	opts.SetClientID(cfg.ClientID)
	opts.SetUsername(cfg.Username)
	opts.SetPassword(cfg.Password)
	opts.SetKeepAlive(cfg.KeepAlive)
	opts.SetConnectTimeout(cfg.ConnectTimeout)
	opts.SetWriteTimeout(cfg.WriteTimeout)
	opts.SetPingTimeout(cfg.PingTimeout)
	opts.SetMessageChannelDepth(cfg.MessageChannelDepth)
	opts.SetCleanSession(cfg.CleanSession)
	opts.SetAutoReconnect(cfg.AutoReconnect)
	opts.SetMaxReconnectInterval(cfg.MaxReconnectInterval)
	// The connect attempt must fail fast rather than retry forever, so the
	// caller's context stays in control of the wait.
	opts.SetConnectRetry(false)
	opts.SetOrderMatters(false)
	if cfg.TLSConfig != nil {
		opts.SetTLSConfig(cfg.TLSConfig)
	}

	opts.SetDefaultPublishHandler(func(_ paho.Client, m paho.Message) {
		c.handleMessage(m)
	})
	opts.SetOnConnectHandler(func(_ paho.Client) {
		c.reconnecting.Store(false)
		c.mu.RLock()
		h := c.onConnect
		c.mu.RUnlock()
		if h != nil {
			h()
		}
	})
	opts.SetConnectionLostHandler(func(_ paho.Client, err error) {
		c.mu.RLock()
		h := c.onConnectionLost
		c.mu.RUnlock()
		if h != nil {
			h(err)
		}
	})
	opts.SetReconnectingHandler(func(_ paho.Client, _ *paho.ClientOptions) {
		c.reconnecting.Store(true)
		c.mu.RLock()
		h := c.onReconnect
		c.mu.RUnlock()
		if h != nil {
			h()
		}
	})

	c.pc = paho.NewClient(opts)
	return c, nil
}

// SetMessageHandler registers the callback invoked for every inbound message.
// A nil handler clears it. It may be called before or after [Client.Connect].
func (c *Client) SetMessageHandler(h Handler) {
	c.mu.Lock()
	c.onMessage = h
	c.mu.Unlock()
}

// SetConnectHandler registers the callback invoked after a successful
// connection (including automatic reconnections).
func (c *Client) SetConnectHandler(h func()) {
	c.mu.Lock()
	c.onConnect = h
	c.mu.Unlock()
}

// SetConnectionLostHandler registers the callback invoked when an established
// connection is lost.
func (c *Client) SetConnectionLostHandler(h func(error)) {
	c.mu.Lock()
	c.onConnectionLost = h
	c.mu.Unlock()
}

// SetReconnectHandler registers the callback invoked when the client begins a
// reconnect attempt.
func (c *Client) SetReconnectHandler(h func()) {
	c.mu.Lock()
	c.onReconnect = h
	c.mu.Unlock()
}

// SetErrorHandler registers a callback for transport-level errors, such as a
// failed connection attempt. A nil handler clears it.
func (c *Client) SetErrorHandler(h func(error)) {
	c.mu.Lock()
	c.onError = h
	c.mu.Unlock()
}

// handleMessage adapts a paho message to the transport [Message] and invokes
// the registered handler.
func (c *Client) handleMessage(m paho.Message) {
	c.mu.RLock()
	h := c.onMessage
	c.mu.RUnlock()
	if h == nil {
		return
	}
	h(Message{
		Topic:     m.Topic(),
		Payload:   m.Payload(),
		QoS:       m.Qos(),
		Retained:  m.Retained(),
		Duplicate: m.Duplicate(),
	})
}

// emitError forwards err to the error handler, if any.
func (c *Client) emitError(err error) {
	if err == nil {
		return
	}
	c.mu.RLock()
	h := c.onError
	c.mu.RUnlock()
	if h != nil {
		h(err)
	}
}

// Connect establishes the broker connection, blocking until it succeeds, fails,
// or ctx is done. A failed attempt is reported to the registered error handler
// and returned to the caller. Calling Connect on an already-connected client is
// a no-op.
//
// When the broker rejects the CONNECT packet the returned error is a
// [ConnackError]; it matches [ErrConnectionLimit] for Webull code 105 and
// [ErrConnectionRefused] for any other rejection. paho v1.5.1 silently
// completes a token for non-standard return codes (Webull uses 100-105), so a
// refusal is also detected by the absence of a live connection.
func (c *Client) Connect(ctx context.Context) error {
	if c.pc == nil {
		return errors.New("mqtt: client is not initialized")
	}
	if c.IsConnected() {
		return nil
	}

	token := c.pc.Connect()
	done := make(chan struct{})
	go func() {
		token.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		c.pc.Disconnect(0)
		return ctx.Err()
	case <-done:
		if err := classifyConnectError(token, token.Error()); err != nil {
			c.emitError(err)
			return err
		}
		if !c.pc.IsConnectionOpen() {
			// A refused connection with a non-standard return code completes
			// the paho token without an error. Report it as a refusal and,
			// when the code is known, as the connection limit.
			refused := error(ErrConnectionRefused)
			if rc := connectReturnCode(token); rc == ConnackConnectionLimit {
				refused = &ConnackError{Code: rc, Err: ErrConnectionLimit}
			}
			c.emitError(refused)
			return refused
		}
		return nil
	}
}

// classifyConnectError converts a completed paho connect token into a typed
// connection error, or nil on success.
func classifyConnectError(token paho.Token, err error) error {
	code := connectReturnCode(token)
	switch {
	case code == ConnackConnectionLimit:
		return &ConnackError{Code: code, Err: ErrConnectionLimit}
	case code != 0:
		return &ConnackError{Code: code, Err: err}
	case err != nil && isConnectionLimitError(err):
		return &ConnackError{Code: ConnackConnectionLimit, Err: ErrConnectionLimit}
	default:
		return err
	}
}

// connectReturnCode extracts the CONNACK return code from a paho connect token.
// It returns 0 when the concrete token does not expose a return code.
func connectReturnCode(token paho.Token) byte {
	if ct, ok := token.(*paho.ConnectToken); ok {
		return ct.ReturnCode()
	}
	return 0
}

// isConnectionLimitError reports whether err text identifies Webull code 105.
// It is a best-effort fallback for transports that report the code only as
// text; the match is deliberately narrow to avoid treating an address or port
// that happens to contain "105" as a limit error.
func isConnectionLimitError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrConnectionLimit) {
		return true
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "connection limit") {
		return true
	}
	for _, marker := range []string{"code 105", "code:105", "code: 105", "code=105", `"105"`} {
		if strings.Contains(msg, marker) {
			return true
		}
	}
	return false
}

// Disconnect closes the connection. quiesce is the number of milliseconds to
// wait for in-flight messages to complete; zero disconnects immediately.
// Disconnect is safe to call when already disconnected.
func (c *Client) Disconnect(quiesce uint) {
	if c.pc == nil {
		return
	}
	c.pc.Disconnect(quiesce)
}

// Close releases the connection. It is equivalent to [Client.Disconnect] with
// a short quiesce period and is safe to call more than once.
func (c *Client) Close() error {
	c.Disconnect(250)
	return nil
}

// IsConnected reports whether the client currently has a live connection. It is
// false while a reconnect is in progress; see [Client.IsReconnecting].
func (c *Client) IsConnected() bool {
	if c.pc == nil {
		return false
	}
	return c.pc.IsConnectionOpen()
}

// IsReconnecting reports whether the client is currently attempting to
// re-establish a connection it lost.
func (c *Client) IsReconnecting() bool {
	return c.reconnecting.Load()
}

// Broker returns the normalized broker URL the client connects to.
func (c *Client) Broker() string {
	normalized, err := normalizeBroker(c.cfg.Broker)
	if err != nil {
		return c.cfg.Broker
	}
	return normalized
}
