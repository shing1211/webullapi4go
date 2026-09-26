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

package events

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"math/rand/v2"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/shing1211/webullapi4go/client"
	eventsevents "github.com/shing1211/webullapi4go/gen/webull/brokerfd/events/v1"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

// Client is a gRPC-based Broker FD event subscription client. It maintains a
// persistent connection to the event stream, optionally reconnecting on
// transport errors, and dispatches connection lifecycle events and data
// payloads to registered handlers.
type Client struct {
	core    *client.Client
	cfg     config
	conn    *grpc.ClientConn
	svc     eventsevents.EventServiceClient
	metrics *eventMetrics

	closeOnce sync.Once

	mu          sync.RWMutex
	onConnect   []func()
	onPing      []func()
	onError     []func(error)
	onData      []func(subscribeType uint32, contentType string, payload []byte)
	onDataEvent []func(*DataEvent)

	runMu     sync.Mutex
	runCancel context.CancelFunc
}

// New returns a Broker FD event client configured with cl and the supplied
// options. It validates the configuration and establishes the gRPC connection
// but does not start the event loop; call [Client.Run] to begin receiving
// events. New returns an error if cl is nil, no endpoint is available, or
// the dial options are invalid.
func New(cl *client.Client, opts ...Option) (*Client, error) {
	if cl == nil {
		return nil, errs.New(errs.CodeInvalidConfig, "brokerfd/events: client is required")
	}
	cfg := defaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.endpoint == "" {
		cfg.endpoint = cl.Endpoints().GRPC
	}
	if cfg.endpoint == "" {
		return nil, errs.New(errs.CodeInvalidConfig, "brokerfd/events: no gRPC endpoint configured")
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	dialOpts := make([]grpc.DialOption, 0, len(cfg.dialOptions)+1)
	if cfg.tls {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})))
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	dialOpts = append(dialOpts, cfg.dialOptions...)

	conn, err := grpc.NewClient(dialTarget(cfg.endpoint, cfg.port), dialOpts...)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInvalidConfig, "brokerfd/events: invalid gRPC target", err)
	}
	return &Client{
		core:    cl,
		cfg:     cfg,
		conn:    conn,
		svc:     eventsevents.NewEventServiceClient(conn),
		metrics: newEventMetrics(cl.ObservabilityConfig()),
	}, nil
}

func dialTarget(host string, port int) string {
	if strings.Contains(host, "://") {
		return host
	}
	return net.JoinHostPort(host, strconv.Itoa(port))
}

// OnConnect registers fn to be called when the client establishes a
// successful subscription with the server. The callbacks are invoked
// synchronously inside the stream receive loop. Multiple callbacks are
// supported; they are called in registration order.
func (c *Client) OnConnect(fn func()) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onConnect = append(c.onConnect, fn)
	c.mu.Unlock()
}

// OnPing registers fn to be called when the server sends a keep-alive ping.
// The callbacks are invoked synchronously inside the stream receive loop.
func (c *Client) OnPing(fn func()) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onPing = append(c.onPing, fn)
	c.mu.Unlock()
}

// OnError registers fn to be called when a transport or protocol error
// occurs during streaming. It is not called for context cancellation or
// normal stream termination. The error is the reason the stream failed; the
// client may still be reconnecting depending on the configured options.
func (c *Client) OnError(fn func(error)) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onError = append(c.onError, fn)
	c.mu.Unlock()
}

// OnData registers fn to be called for each data event received from the
// server. subscribeType identifies the event category, contentType is the MIME
// type of the payload (for example "application/json"), and payload holds the
// raw event data. The callbacks are invoked synchronously inside the stream
// receive loop.
func (c *Client) OnData(fn func(subscribeType uint32, contentType string, payload []byte)) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onData = append(c.onData, fn)
	c.mu.Unlock()
}

// OnDataEvent registers fn to be called for each data event received from the
// server, delivering the whole [DataEvent] including the server-assigned
// RequestId and Timestamp. Use this instead of [Client.OnData] when the event
// metadata matters; OnData discards those two fields because its callback
// signature predates them.
//
// The two registrations are independent and may be combined: each event is
// delivered to every OnData handler and then to every OnDataEvent handler, in
// registration order within each group. Callbacks are invoked synchronously
// inside the stream receive loop, so a slow handler delays later delivery, the
// other handler group, and the receive path until it returns.
func (c *Client) OnDataEvent(fn func(*DataEvent)) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onDataEvent = append(c.onDataEvent, fn)
	c.mu.Unlock()
}

// Run starts the event subscription loop. It blocks until the context is
// cancelled, the stream terminates, or an unrecoverable error occurs. If
// auto-reconnect is enabled (the default), Run retries transient transport
// errors indefinitely up to [WithMaxReconnectAttempts]; when disabled it
// returns after a single stream attempt. Run must not be called more than once
// concurrently. [Client.Close] is terminal for this client; callers must not
// reuse it by calling Run again.
func (c *Client) Run(ctx context.Context) error {
	if c == nil || c.conn == nil {
		return errs.New(errs.CodeInvalidConfig, "brokerfd/events: client is not initialized")
	}
	ctx, _ = ensureCorrelationID(ctx)
	ctx, cancel := context.WithCancel(ctx)
	c.setRunCancel(cancel)
	defer func() {
		c.setRunCancel(nil)
		cancel()
	}()

	if !c.cfg.autoReconnect {
		return c.fail(ctx, c.runOnce(ctx, 1))
	}
	return c.runReconnecting(ctx)
}

func (c *Client) runOnce(ctx context.Context, attempt int) (err error) {
	ctx, telemetry := c.startEventAttempt(ctx, attempt)
	defer func() {
		telemetry.finish(err)
	}()
	if err = c.waitForReady(ctx); err != nil {
		return err
	}
	stream, err := c.open(ctx)
	if err != nil {
		return err
	}
	for {
		resp, recvErr := stream.Recv()
		if recvErr != nil {
			return c.recvError(ctx, recvErr)
		}
		if dispatchErr := c.dispatch(resp); dispatchErr != nil {
			return dispatchErr
		}
	}
}

func (c *Client) runReconnecting(ctx context.Context) error {
	attempt := 0
	for {
		err := c.runOnce(ctx, attempt+1)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if !isRetryableStreamError(err) {
			return c.fail(ctx, err)
		}
		if err != nil {
			c.emitError(err)
		}
		if c.cfg.maxReconnectAttempts > 0 && attempt >= c.cfg.maxReconnectAttempts {
			exhausted := errs.New(errs.CodeTransport, "brokerfd/events: reconnect attempts exhausted")
			c.emitError(exhausted)
			return exhausted
		}
		delay := c.reconnectDelay(attempt)
		attempt++
		if !sleepContext(ctx, delay) {
			return ctx.Err()
		}
	}
}

func isRetryableStreamError(err error) bool {
	if err == nil {
		return true
	}
	if errors.Is(err, errs.ErrConnectionLimitExceeded) {
		return false
	}
	return errs.Is(err, errs.CodeTransport)
}

func (c *Client) reconnectDelay(attempt int) time.Duration {
	base := c.cfg.reconnectBaseDelay
	if base <= 0 {
		base = DefaultReconnectBaseDelay
	}
	maxDelay := c.cfg.reconnectMaxDelay
	if maxDelay < base {
		maxDelay = base
	}
	delay := base
	for i := 0; i < attempt && delay < maxDelay; i++ {
		delay *= 2
	}
	if delay > maxDelay {
		delay = maxDelay
	}
	half := delay / 2
	if half <= 0 {
		return delay
	}
	return half + time.Duration(rand.Int64N(int64(half)+1))
}

func sleepContext(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		return ctx.Err() == nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

// Close gracefully shuts down the client. It cancels the run context, waits
// for any in-flight handlers to finish, and closes the underlying gRPC
// connection. Close is safe to call multiple times; subsequent calls return
// nil.
func (c *Client) Close() error {
	var err error
	c.closeOnce.Do(func() {
		c.runMu.Lock()
		cancel := c.runCancel
		c.runMu.Unlock()
		if cancel != nil {
			cancel()
		}
		if c.conn != nil {
			err = c.conn.Close()
		}
	})
	return err
}

func (c *Client) setRunCancel(cancel context.CancelFunc) {
	c.runMu.Lock()
	c.runCancel = cancel
	c.runMu.Unlock()
}

func (c *Client) fail(ctx context.Context, err error) error {
	if err == nil || ctx.Err() != nil {
		return err
	}
	c.emitError(err)
	return err
}

func (c *Client) waitForReady(ctx context.Context) error {
	parent := ctx
	ctx, cancel := context.WithTimeout(ctx, c.cfg.dialTimeout)
	defer cancel()
	c.conn.Connect()
	for {
		state := c.conn.GetState()
		if state == connectivity.Ready {
			return nil
		}
		if !c.conn.WaitForStateChange(ctx, state) {
			if parent.Err() != nil {
				return parent.Err()
			}
			return errs.New(errs.CodeTransport, "brokerfd/events: gRPC connection not ready before the dial timeout")
		}
	}
}

func (c *Client) open(ctx context.Context) (*streamClient, error) {
	now := time.Now()
	req := NewSubscribeRequest(uint32(c.cfg.subscribeTypes), c.cfg.accounts)
	body, err := proto.Marshal(req)
	if err != nil {
		return nil, errs.Wrap(errs.CodeTransport, "brokerfd/events: marshal subscribe request", err)
	}
	md, err := subscribeMetadata(c.core.Config().AppKey, c.core.Config().AppSecret, now, body)
	if err != nil {
		return nil, err
	}
	md = c.addEventMetadata(ctx, md)
	stream, err := c.svc.Subscribe(metadata.NewOutgoingContext(ctx, md), req, grpc.WaitForReady(true))
	if err != nil {
		return nil, rpcError("brokerfd/events: subscribe", err)
	}
	return &streamClient{stream: stream}, nil
}

type streamClient struct {
	stream eventsevents.EventService_SubscribeClient
}

func (s *streamClient) Recv() (*SubscribeResponse, error) {
	resp, err := s.stream.Recv()
	if err == io.EOF {
		return nil, io.EOF
	}
	if err != nil {
		return nil, err
	}
	return &SubscribeResponse{SubscribeResponse: resp}, nil
}

func (c *Client) recvError(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if errors.Is(err, io.EOF) {
		return nil
	}
	return rpcError("brokerfd/events: stream receive", err)
}

func rpcError(message string, err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return errs.Wrap(errs.CodeTransport, message, err)
	}
	code := errs.CodeTransport
	switch st.Code() {
	case codes.Unauthenticated:
		code = errs.CodeAuth
	case codes.PermissionDenied:
		code = errs.CodeForbidden
	case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.FailedPrecondition:
		code = errs.CodeAPI
	case codes.ResourceExhausted:
		code = errs.CodeRateLimited
	case codes.Unimplemented:
		code = errs.CodeUnsupported
	}
	return errs.Wrap(code, message+": "+st.Message(), err)
}

func (c *Client) dispatch(resp *SubscribeResponse) error {
	switch resp.EventType() {
	case eventsevents.EventType_SubscribeSuccess:
		c.emitConnect()
		return nil
	case eventsevents.EventType_Ping:
		c.emitPing()
		return nil
	case eventsevents.EventType_AuthError:
		return errs.New(errs.CodeAuth, "brokerfd/events: authentication failed; verify the App Key, App Secret, and signing parameters")
	case eventsevents.EventType_NumOfConnExceed:
		return errs.Wrap(errs.CodeTransport, "brokerfd/events: connection limit exceeded; Webull allows at most 5 concurrent event connections per App Key", errs.ErrConnectionLimitExceeded)
	case eventsevents.EventType_SubscribeExpired:
		return errs.Wrap(errs.CodeAuth, "brokerfd/events: subscription expired; reconnect to resume", errs.ErrSubscriptionExpired)
	default:
		c.emitData(resp)
		return nil
	}
}

func (c *Client) emitConnect() {
	c.mu.RLock()
	handlers := make([]func(), len(c.onConnect))
	copy(handlers, c.onConnect)
	c.mu.RUnlock()
	for _, h := range handlers {
		h()
	}
}

func (c *Client) emitPing() {
	c.mu.RLock()
	handlers := make([]func(), len(c.onPing))
	copy(handlers, c.onPing)
	c.mu.RUnlock()
	for _, h := range handlers {
		h()
	}
}

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

func (c *Client) emitData(resp *SubscribeResponse) {
	c.mu.RLock()
	handlers := make([]func(uint32, string, []byte), len(c.onData))
	copy(handlers, c.onData)
	eventHandlers := make([]func(*DataEvent), len(c.onDataEvent))
	copy(eventHandlers, c.onDataEvent)
	c.mu.RUnlock()

	de := resp.ToDataEvent()
	for _, h := range handlers {
		h(de.SubscribeType, de.ContentType, de.Payload)
	}
	for _, h := range eventHandlers {
		h(de)
	}
}
