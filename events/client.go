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
	eventsevents "github.com/shing1211/webullapi4go/gen/webull/trade/events/v1"
	"github.com/shing1211/webullapi4go/pkg/errors"
)

// Client is a Webull trade-event streaming client. It owns a gRPC connection to
// the event service and routes pushed messages to the handlers registered with
// the On* methods.
//
// A Client is safe for concurrent use. Handlers may be registered before or
// after [Client.Run]. The underlying [client.Client] is owned by the caller and
// is not closed by [Client.Close].
type Client struct {
	core   *client.Client
	cfg    config
	conn   *grpc.ClientConn
	events eventsevents.EventServiceClient

	closeOnce sync.Once

	mu         sync.RWMutex
	onConnect  []func()
	onPing     []func()
	onError    []func(error)
	onEvent    []func(subscribeType uint32, contentType string, payload []byte)
	onOrder    []func(*OrderEvent)
	onPosition []func(*PositionEvent)
	onOption   []func(*OptionEvent)

	runMu     sync.Mutex
	runCancel context.CancelFunc
}

// New returns an event client bound to cl. The gRPC endpoint is taken from cl's
// resolved endpoints unless overridden with [WithGRPCEndpoint]. New validates
// the configuration and creates the (lazily connected) gRPC channel; it does
// not open the stream. Call [Client.Run] for that.
//
// TLS is used by default. The client is owned by cl: New never closes cl, and
// [Client.Close] does not either.
func New(cl *client.Client, opts ...Option) (*Client, error) {
	if cl == nil {
		return nil, errs.New(errs.CodeInvalidConfig, "events: client is required")
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
		return nil, errs.New(errs.CodeInvalidConfig, "events: no gRPC endpoint configured")
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	dialOpts := make([]grpc.DialOption, 0, len(cfg.dialOptions)+1)
	if cfg.tls {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	dialOpts = append(dialOpts, cfg.dialOptions...)

	conn, err := grpc.NewClient(dialTarget(cfg.endpoint, cfg.port), dialOpts...)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInvalidConfig, "events: invalid gRPC target", err)
	}
	return &Client{
		core:   cl,
		cfg:    cfg,
		conn:   conn,
		events: eventsevents.NewEventServiceClient(conn),
	}, nil
}

// dialTarget returns the gRPC target for host. A host that already carries a
// resolver scheme (for example "passthrough:///bufnet") is used verbatim so its
// scheme is preserved; a bare host is combined with port.
func dialTarget(host string, port int) string {
	if strings.Contains(host, "://") {
		return host
	}
	return net.JoinHostPort(host, strconv.Itoa(port))
}

// OnConnect registers a handler invoked after the server acknowledges the
// subscription (a SubscribeSuccess event).
func (c *Client) OnConnect(fn func()) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onConnect = append(c.onConnect, fn)
	c.mu.Unlock()
}

// OnPing registers a handler invoked for each server heartbeat (a Ping event).
func (c *Client) OnPing(fn func()) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onPing = append(c.onPing, fn)
	c.mu.Unlock()
}

// OnError registers a handler for asynchronous errors: a failed connection or
// subscription, a receive failure, a payload that could not be decoded, or one
// of the terminal server events (AuthError, NumOfConnExceed, SubscribeExpired).
func (c *Client) OnError(fn func(error)) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onError = append(c.onError, fn)
	c.mu.Unlock()
}

// OnEvent registers a handler for every data event. The first argument is the
// event kind: [EventOrder], [EventPosition], or [EventOption]. contentType is
// the MIME type of payload ("application/json" for JSON payloads), which is
// delivered verbatim as raw bytes. For decoded payloads use [Client.OnOrder],
// [Client.OnPosition], or [Client.OnOption]; OnEvent fires first so both may be
// registered.
func (c *Client) OnEvent(fn func(subscribeType uint32, contentType string, payload []byte)) {
	if fn == nil {
		return
	}
	c.mu.Lock()
	c.onEvent = append(c.onEvent, fn)
	c.mu.Unlock()
}

// Run connects, subscribes, and blocks pumping events until ctx is cancelled.
// It reconnects and re-subscribes when the stream ends or fails, unless
// automatic reconnect is disabled with [WithAutoReconnect].
//
// With reconnect enabled (the default), a clean server end and transient
// transport failures are retried with exponential backoff and jitter; a
// terminal failure ends Run with a typed error and invokes the [Client.OnError]
// handlers. Terminal failures are authentication, permission, account, and
// configuration errors, plus the terminal server events (AuthError,
// NumOfConnExceed, SubscribeExpired). Consecutive attempts are bounded by
// [WithMaxReconnectAttempts]; the default is unlimited. [Client.OnConnect]
// fires after every SubscribeSuccess acknowledgement.
//
// When reconnect is disabled, Run returns after the first stream ends. Run
// returns nil on a clean server end and ctx.Err() when ctx is cancelled
// (including via [Client.Close]).
//
// Run uses the credentials, subscribe types, and accounts captured at [New];
// each (re)connect signs and sends a fresh Subscribe request.
func (c *Client) Run(ctx context.Context) error {
	if c == nil || c.conn == nil {
		return errs.New(errs.CodeInvalidConfig, "events: client is not initialized")
	}
	ctx, cancel := context.WithCancel(ctx)
	c.setRunCancel(cancel)
	defer func() {
		c.setRunCancel(nil)
		cancel()
	}()

	if !c.cfg.autoReconnect {
		return c.fail(ctx, c.runOnce(ctx))
	}
	return c.runReconnecting(ctx)
}

// runOnce opens one stream and pumps it until it ends, fails, or ctx is
// cancelled. A clean end of stream is reported as a nil error.
func (c *Client) runOnce(ctx context.Context) error {
	if err := c.waitForReady(ctx); err != nil {
		return err
	}
	stream, err := c.open(ctx)
	if err != nil {
		return err
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			return c.recvError(ctx, err)
		}
		if err := c.dispatch(resp); err != nil {
			return err
		}
	}
}

// runReconnecting drives runOnce, retrying retryable failures and clean stream
// ends with exponential backoff and jitter until ctx is cancelled, a terminal
// error occurs, or the attempt budget is exhausted.
func (c *Client) runReconnecting(ctx context.Context) error {
	attempt := 0
	for {
		err := c.runOnce(ctx)
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
			if err != nil {
				return err
			}
			exhausted := errs.New(errs.CodeTransport, "events: reconnect attempts exhausted")
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

// errTerminalStream marks a stream failure caused by a terminal server event
// that must not be retried even though its transport classification would
// otherwise look transient.
var errTerminalStream = errors.New("events: terminal stream event")

// isRetryableStreamError reports whether Run should reconnect after err. A nil
// error is a clean end of stream and is retryable; a typed transport error is
// retryable unless it wraps the terminal marker; every other typed error
// (authentication, permission, account, configuration, unsupported) is
// terminal.
func isRetryableStreamError(err error) bool {
	if err == nil {
		return true
	}
	if errors.Is(err, errTerminalStream) {
		return false
	}
	return errs.Is(err, errs.CodeTransport)
}

// reconnectDelay returns the backoff before the next attempt. The base delay
// doubles per attempt up to the configured maximum, then a jittered value in
// the upper half of that window is returned so concurrent clients do not
// reconnect in lockstep.
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

// sleepContext waits for d or until ctx is done, reporting false when ctx won.
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

// Close cancels an in-progress [Client.Run] and closes the gRPC connection. It
// is idempotent and does not close the underlying [client.Client].
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

// setRunCancel stores or clears the cancel function of the active Run.
func (c *Client) setRunCancel(cancel context.CancelFunc) {
	c.runMu.Lock()
	c.runCancel = cancel
	c.runMu.Unlock()
}

// fail reports err to the error handlers unless it is a context cancellation,
// which is a normal shutdown rather than a stream error.
func (c *Client) fail(ctx context.Context, err error) error {
	if err == nil || ctx.Err() != nil {
		return err
	}
	c.emitError(err)
	return err
}

// waitForReady triggers the connection and blocks until it is ready, the dial
// timeout elapses, or ctx is done. A dial timeout is reported as a typed
// transport error so [Client.Run] can retry it; only cancellation of the
// caller's context is reported as a context error.
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
			return errs.New(errs.CodeTransport, "events: gRPC connection not ready before the dial timeout")
		}
	}
}

// open signs and issues the Subscribe RPC, returning the response stream.
func (c *Client) open(ctx context.Context) (grpc.ServerStreamingClient[eventsevents.SubscribeResponse], error) {
	now := time.Now()
	req := c.cfg.newSubscribeRequest(now)
	body, err := proto.Marshal(req)
	if err != nil {
		return nil, errs.Wrap(errs.CodeTransport, "events: marshal subscribe request", err)
	}
	md, err := subscribeMetadata(c.core.Config().AppKey, c.core.Config().AppSecret, now, body)
	if err != nil {
		return nil, err
	}
	stream, err := c.events.Subscribe(metadata.NewOutgoingContext(ctx, md), req, grpc.WaitForReady(true))
	if err != nil {
		return nil, rpcError("events: subscribe", err)
	}
	return stream, nil
}

// recvError classifies a stream receive failure. A clean end of stream is not
// an error and a cancellation is reported as the context's error.
func (c *Client) recvError(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if errors.Is(err, io.EOF) {
		return nil
	}
	return rpcError("events: stream receive", err)
}

// rpcError maps a gRPC status error onto the SDK's typed error model so callers
// can branch on a category rather than a status code.
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

// dispatch routes one response by event type.
func (c *Client) dispatch(resp *eventsevents.SubscribeResponse) error {
	switch resp.GetEventType() {
	case eventsevents.EventType_SubscribeSuccess:
		c.emitConnect()
		return nil
	case eventsevents.EventType_Ping:
		c.emitPing()
		return nil
	case eventsevents.EventType_AuthError:
		return errs.New(errs.CodeAuth, "events: authentication failed; verify the App Key, App Secret, and signing parameters")
	case eventsevents.EventType_NumOfConnExceed:
		return errs.Wrap(errs.CodeTransport, "events: connection limit exceeded; Webull allows at most 5 concurrent event connections per App Key", errTerminalStream)
	case eventsevents.EventType_SubscribeExpired:
		return errs.Wrap(errs.CodeAuth, "events: subscription expired; reconnect to resume", errs.ErrSubscriptionExpired)
	default:
		c.routeDataEvent(resp)
		return nil
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

// emitPing invokes every registered ping handler outside the lock.
func (c *Client) emitPing() {
	c.mu.RLock()
	handlers := make([]func(), len(c.onPing))
	copy(handlers, c.onPing)
	c.mu.RUnlock()
	for _, h := range handlers {
		h()
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

// emitEvent invokes every registered data-event handler outside the lock.
func (c *Client) emitEvent(subscribeType uint32, contentType string, payload []byte) {
	c.mu.RLock()
	handlers := make([]func(uint32, string, []byte), len(c.onEvent))
	copy(handlers, c.onEvent)
	c.mu.RUnlock()
	for _, h := range handlers {
		h(subscribeType, contentType, payload)
	}
}
