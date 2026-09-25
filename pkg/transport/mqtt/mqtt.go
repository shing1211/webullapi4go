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

	pkgerrs "github.com/shing1211/webullapi4go/pkg/errors"
)

const (
	DefaultKeepAlive                 = 30 * time.Second
	DefaultConnectTimeout            = 15 * time.Second
	DefaultWriteTimeout              = 5 * time.Second
	DefaultPingTimeout               = 10 * time.Second
	DefaultMessageChannelDepth  uint = 100
	DefaultMaxReconnectInterval      = 30 * time.Second
)

const ConnackConnectionLimit byte = 105

var (
	ErrConnectionRefused = pkgerrs.NewSentinel(pkgerrs.CodeTransport, "mqtt: connection refused by broker")
	ErrConnectionLimit   = pkgerrs.NewSentinel(pkgerrs.CodeTransport, "mqtt: exceeds connection limit (code 105): at most 5 concurrent connections per App Key; wait about 1 minute after a disconnect before reconnecting")
	errClientClosed      = errors.New("mqtt: client is closed")
)

type ConnackError struct {
	Code byte
	Err  error
}

func (e *ConnackError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("mqtt: connection refused (code %d): %v", e.Code, e.Err)
	}
	return fmt.Sprintf("mqtt: connection refused (code %d)", e.Code)
}

func (e *ConnackError) Unwrap() error { return e.Err }

func (e *ConnackError) Is(target error) bool {
	if target == ErrConnectionLimit {
		return e.Code == ConnackConnectionLimit
	}
	if target == ErrConnectionRefused {
		return true
	}
	return errors.Is(e.Err, target)
}

type Message struct {
	Topic     string
	Payload   []byte
	QoS       byte
	Retained  bool
	Duplicate bool
}

type Handler func(Message)

type Config struct {
	Broker               string
	ClientID             string
	Username             string
	Password             string
	KeepAlive            time.Duration
	ConnectTimeout       time.Duration
	WriteTimeout         time.Duration
	PingTimeout          time.Duration
	MessageChannelDepth  uint
	CleanSession         bool
	AutoReconnect        bool
	MaxReconnectInterval time.Duration
	TLSConfig            *tls.Config
}

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

type Client struct {
	mu          sync.RWMutex
	lifecycleMu sync.Mutex
	cfg         Config
	pc          paho.Client
	closedCh    chan struct{}

	closeOnce    sync.Once
	closed       atomic.Bool
	reconnecting atomic.Bool

	onMessage        Handler
	onConnect        func()
	onConnectionLost func(error)
	onReconnect      func()
	onError          func(error)
}

func New(cfg Config) (*Client, error) {
	cfg = cfg.withDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	c := &Client{cfg: cfg, closedCh: make(chan struct{})}

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
	opts.SetConnectRetry(false)
	opts.SetOrderMatters(false)
	if cfg.TLSConfig != nil {
		opts.SetTLSConfig(cfg.TLSConfig)
	}

	opts.SetDefaultPublishHandler(func(_ paho.Client, m paho.Message) {
		c.handleMessage(m)
	})
	opts.SetOnConnectHandler(func(_ paho.Client) {
		c.handlePahoConnect()
	})
	opts.SetConnectionLostHandler(func(_ paho.Client, err error) {
		c.handlePahoConnectionLost(err)
	})
	opts.SetReconnectingHandler(func(_ paho.Client, _ *paho.ClientOptions) {
		c.handlePahoReconnecting()
	})

	c.pc = paho.NewClient(opts)
	return c, nil
}

func (c *Client) handlePahoConnect() {
	if c.closed.Load() {
		return
	}
	c.reconnecting.Store(false)
	c.mu.RLock()
	h := c.onConnect
	c.mu.RUnlock()
	if h != nil {
		h()
	}
}

func (c *Client) handlePahoConnectionLost(err error) {
	if c.closed.Load() {
		return
	}
	c.mu.RLock()
	h := c.onConnectionLost
	c.mu.RUnlock()
	if h != nil {
		h(err)
	}
}

func (c *Client) handlePahoReconnecting() {
	if c.closed.Load() {
		return
	}
	c.reconnecting.Store(true)
	c.mu.RLock()
	h := c.onReconnect
	c.mu.RUnlock()
	if h != nil {
		h()
	}
}

func (c *Client) SetMessageHandler(h Handler) {
	c.mu.Lock()
	c.onMessage = h
	c.mu.Unlock()
}

func (c *Client) SetConnectHandler(h func()) {
	c.mu.Lock()
	c.onConnect = h
	c.mu.Unlock()
}

func (c *Client) SetConnectionLostHandler(h func(error)) {
	c.mu.Lock()
	c.onConnectionLost = h
	c.mu.Unlock()
}

func (c *Client) SetReconnectHandler(h func()) {
	c.mu.Lock()
	c.onReconnect = h
	c.mu.Unlock()
}

func (c *Client) SetErrorHandler(h func(error)) {
	c.mu.Lock()
	c.onError = h
	c.mu.Unlock()
}

func (c *Client) handleMessage(m paho.Message) {
	if c.closed.Load() {
		return
	}
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

func (c *Client) emitError(err error) {
	if err == nil || c.closed.Load() {
		return
	}
	c.mu.RLock()
	h := c.onError
	c.mu.RUnlock()
	if h != nil {
		h(err)
	}
}

func (c *Client) Connect(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	c.lifecycleMu.Lock()
	if c.pc == nil {
		c.lifecycleMu.Unlock()
		return errors.New("mqtt: client is not initialized")
	}
	if c.closed.Load() {
		c.lifecycleMu.Unlock()
		return errClientClosed
	}
	if c.IsConnected() {
		c.lifecycleMu.Unlock()
		return nil
	}
	if err := ctx.Err(); err != nil {
		c.lifecycleMu.Unlock()
		return err
	}
	if c.closedCh == nil {
		c.closedCh = make(chan struct{})
	}
	closedCh := c.closedCh
	token := c.pc.Connect()
	c.lifecycleMu.Unlock()

	select {
	case <-ctx.Done():
		if c.closed.Load() {
			return errClientClosed
		}
		c.Disconnect(0)
		return ctx.Err()
	case <-closedCh:
		return errClientClosed
	case <-token.Done():
		if c.closed.Load() {
			return errClientClosed
		}
		if err := classifyConnectError(token, token.Error()); err != nil {
			c.emitError(err)
			return err
		}
		if !c.pc.IsConnectionOpen() {
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

func connectReturnCode(token paho.Token) byte {
	if ct, ok := token.(*paho.ConnectToken); ok {
		return ct.ReturnCode()
	}
	return 0
}

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

func (c *Client) Disconnect(quiesce uint) {
	c.reconnecting.Store(false)
	if c.closed.Load() || c.pc == nil {
		return
	}
	c.pc.Disconnect(quiesce)
}

func (c *Client) Close() error {
	c.closeOnce.Do(func() {
		c.lifecycleMu.Lock()
		c.closed.Store(true)
		c.reconnecting.Store(false)
		if c.closedCh == nil {
			c.closedCh = make(chan struct{})
		}
		close(c.closedCh)
		pc := c.pc
		c.lifecycleMu.Unlock()
		if pc != nil {
			pc.Disconnect(250)
		}
	})
	return nil
}

func (c *Client) IsConnected() bool {
	if c.closed.Load() || c.pc == nil {
		return false
	}
	return c.pc.IsConnectionOpen()
}

func (c *Client) IsReconnecting() bool {
	return c.reconnecting.Load()
}

func (c *Client) Broker() string {
	normalized, err := normalizeBroker(c.cfg.Broker)
	if err != nil {
		return c.cfg.Broker
	}
	return normalized
}
