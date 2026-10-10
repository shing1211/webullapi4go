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
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	"github.com/shing1211/webullapi4go/client"
	eventsevents "github.com/shing1211/webullapi4go/gen/webull/brokerfd/events/v1"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

type capture struct {
	request  *eventsevents.SubscribeRequest
	metadata metadata.MD
}

type fakeServer struct {
	eventsevents.UnimplementedEventServiceServer

	messages []*eventsevents.SubscribeResponse

	mu  sync.Mutex
	got capture
}

func (s *fakeServer) Subscribe(req *eventsevents.SubscribeRequest, stream grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error {
	md, _ := metadata.FromIncomingContext(stream.Context())
	s.mu.Lock()
	s.got = capture{request: req, metadata: md}
	s.mu.Unlock()
	for _, msg := range s.messages {
		if err := stream.Send(msg); err != nil {
			return err
		}
	}
	<-stream.Context().Done()
	return stream.Context().Err()
}

func (s *fakeServer) captured() capture {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.got
}

type reconnectServer struct {
	eventsevents.UnimplementedEventServiceServer
	subscribes      atomic.Int32
	secondSubscribe chan struct{}
}

func (s *reconnectServer) Subscribe(_ *eventsevents.SubscribeRequest, stream grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error {
	if s.subscribes.Add(1) == 1 {
		return nil
	}
	close(s.secondSubscribe)
	<-stream.Context().Done()
	return stream.Context().Err()
}

func newCore(t *testing.T) *client.Client {
	t.Helper()
	return newCoreWithOptions(t)
}

func newCoreWithOptions(t *testing.T, opts ...client.Option) *client.Client {
	t.Helper()
	base := []client.Option{
		client.WithAppKey("test-app-key"),
		client.WithAppSecret("test-app-secret"),
	}
	cl, err := client.New(append(base, opts...)...)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })
	return cl
}

func newClient(t *testing.T, impl eventsevents.EventServiceServer, opts ...Option) *Client {
	t.Helper()
	return newClientWithCore(t, newCore(t), impl, opts...)
}

func newClientWithCore(t *testing.T, core *client.Client, impl eventsevents.EventServiceServer, opts ...Option) *Client {
	t.Helper()
	lis := bufconn.Listen(1 << 20)
	srv := grpc.NewServer()
	eventsevents.RegisterEventServiceServer(srv, impl)
	serveDone := make(chan struct{})
	go func() {
		defer close(serveDone)
		_ = srv.Serve(lis)
	}()
	t.Cleanup(func() {
		srv.Stop()
		<-serveDone
	})

	base := []Option{
		WithGRPCEndpoint("passthrough:///bufnet"),
		WithTLS(false),
		WithDialTimeout(5 * time.Second),
		WithGRPCDialOption(grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		})),
	}
	cl, err := New(core, append(base, opts...)...)
	if err != nil {
		t.Fatalf("events.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })
	return cl
}

type receivedData struct {
	subscribeType uint32
	contentType   string
	payload       []byte
}

type testRun struct {
	cancel context.CancelFunc
	done   chan struct{}
	err    error
}

func startTestRun(t *testing.T, cl *Client, ctx context.Context) *testRun {
	t.Helper()
	ctx, cancel := context.WithCancel(ctx)
	run := &testRun{cancel: cancel, done: make(chan struct{})}
	go func() {
		run.err = cl.Run(ctx)
		close(run.done)
	}()
	t.Cleanup(func() { _ = run.stop(t) })
	return run
}

func (r *testRun) wait(t *testing.T) error {
	t.Helper()
	select {
	case <-r.done:
		return r.err
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return")
		return nil
	}
}

func (r *testRun) stop(t *testing.T) error {
	t.Helper()
	r.cancel()
	return r.wait(t)
}

func TestRunDispatchesSignedEvents(t *testing.T) {
	fake := &fakeServer{messages: []*eventsevents.SubscribeResponse{
		{EventType: eventsevents.EventType_SubscribeSuccess, Payload: `{"status":"ok"}`},
		{EventType: eventsevents.EventType_Ping},
		{
			EventType:     1024,
			SubscribeType: 1024,
			ContentType:   "application/json",
			Payload:       `{"account_id":"acct-1"}`,
		},
	}}
	cl := newClient(t, fake, WithAccounts([]string{"acct-1"}))

	connected := make(chan struct{}, 1)
	pinged := make(chan struct{}, 1)
	data := make(chan receivedData, 4)
	errCh := make(chan error, 4)
	cl.OnConnect(func() { trySend(connected, struct{}{}) })
	cl.OnPing(func() { trySend(pinged, struct{}{}) })
	cl.OnData(func(st uint32, ct string, payload []byte) {
		trySend(data, receivedData{st, ct, payload})
	})
	cl.OnError(func(err error) { trySend(errCh, err) })

	run := startTestRun(t, cl, context.Background())

	waitSignal(t, connected, "OnConnect")
	waitSignal(t, pinged, "OnPing")

	var got receivedData
	select {
	case got = <-data:
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for OnData (errors: %v)", drainErr(errCh))
	}
	if got.subscribeType != 1024 {
		t.Errorf("subscribeType = %d, want 1024", got.subscribeType)
	}
	if got.contentType != "application/json" {
		t.Errorf("contentType = %q, want application/json", got.contentType)
	}
	if string(got.payload) != `{"account_id":"acct-1"}` {
		t.Errorf("payload = %q", string(got.payload))
	}

	cap := fake.captured()
	for _, key := range []string{
		"x-app-key",
		"x-signature-algorithm",
		"x-signature-version",
		"x-signature-nonce",
		"x-timestamp",
		"x-signature",
	} {
		if len(cap.metadata.Get(key)) == 0 {
			t.Errorf("metadata %q missing from %v", key, cap.metadata)
		}
	}
	if got := cap.metadata.Get("x-signature-algorithm"); len(got) != 1 || got[0] != "HMAC-SHA256" {
		t.Errorf("x-signature-algorithm = %v, want [HMAC-SHA256]", got)
	}
	if got := cap.metadata.Get("x-app-key"); len(got) != 1 || got[0] != "test-app-key" {
		t.Errorf("x-app-key = %v, want [test-app-key]", got)
	}
	for _, key := range []string{"x-signature-nonce", "x-signature"} {
		if v := cap.metadata.Get(key); len(v) == 1 && v[0] == "" {
			t.Errorf("metadata %q is empty", key)
		}
	}
	if len(cap.metadata.Get("host")) != 0 {
		t.Errorf("host metadata must not be sent, got %v", cap.metadata.Get("host"))
	}
	if values := cap.metadata.Get("x-correlation-id"); len(values) != 1 || values[0] == "" {
		t.Errorf("x-correlation-id = %v, want one non-empty value", values)
	}

	if cap.request == nil {
		t.Fatal("server did not receive a SubscribeRequest")
	}
	if cap.request.GetTimestamp() <= 0 {
		t.Errorf("timestamp = %d, want a positive epoch-millis value", cap.request.GetTimestamp())
	}
	if got := cap.request.GetAccounts(); len(got) != 1 || got[0] != "acct-1" {
		t.Errorf("accounts = %v, want [acct-1]", got)
	}

	if err := run.stop(t); !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want context.Canceled", err)
	}
}

func TestRunStopsOnContextCancel(t *testing.T) {
	fake := &fakeServer{messages: []*eventsevents.SubscribeResponse{
		{EventType: eventsevents.EventType_SubscribeSuccess},
	}}
	cl := newClient(t, fake)
	connected := make(chan struct{}, 1)
	errCh := make(chan error, 1)
	cl.OnConnect(func() { trySend(connected, struct{}{}) })
	cl.OnError(func(err error) { trySend(errCh, err) })

	run := startTestRun(t, cl, context.Background())
	waitSignal(t, connected, "OnConnect")
	if err := run.stop(t); !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want context.Canceled", err)
	}
	select {
	case err := <-errCh:
		t.Fatalf("OnError received cancellation error %v", err)
	default:
	}
}

func TestCloseCancelsActiveRunWithoutReportingError(t *testing.T) {
	fake := &fakeServer{messages: []*eventsevents.SubscribeResponse{
		{EventType: eventsevents.EventType_SubscribeSuccess},
	}}
	cl := newClient(t, fake)
	connected := make(chan struct{}, 1)
	errCh := make(chan error, 1)
	cl.OnConnect(func() { trySend(connected, struct{}{}) })
	cl.OnError(func(err error) { trySend(errCh, err) })

	run := startTestRun(t, cl, context.Background())
	waitSignal(t, connected, "OnConnect")
	if err := cl.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if err := run.wait(t); !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want context.Canceled", err)
	}
	select {
	case err := <-errCh:
		t.Fatalf("OnError received shutdown error %v", err)
	default:
	}
}

func TestReconnectCancellationStopsWithoutReportingError(t *testing.T) {
	srv := &reconnectServer{secondSubscribe: make(chan struct{})}
	cl := newClient(t, srv,
		WithReconnectBaseDelay(time.Millisecond),
		WithReconnectMaxDelay(2*time.Millisecond),
		WithMaxReconnectAttempts(2),
	)
	errCh := make(chan error, 1)
	cl.OnError(func(err error) { trySend(errCh, err) })
	run := startTestRun(t, cl, context.Background())
	waitSignal(t, srv.secondSubscribe, "reconnected Subscribe")
	run.cancel()
	if err := run.wait(t); !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want context.Canceled", err)
	}
	select {
	case err := <-errCh:
		t.Fatalf("OnError received reconnect cancellation %v", err)
	default:
	}
}

func TestRunReportsTerminalEvent(t *testing.T) {
	fake := &fakeServer{messages: []*eventsevents.SubscribeResponse{
		{EventType: eventsevents.EventType_AuthError},
	}}
	cl := newClient(t, fake)
	errCh := make(chan error, 4)
	cl.OnError(func(err error) { trySend(errCh, err) })

	runErr := cl.Run(context.Background())
	if !errs.Is(runErr, errs.CodeAuth) {
		t.Fatalf("Run() error = %v, want an auth error", runErr)
	}
	select {
	case err := <-errCh:
		if !errs.Is(err, errs.CodeAuth) {
			t.Errorf("OnError received %v, want an auth error", err)
		}
	case <-time.After(time.Second):
		t.Error("OnError was not invoked for AuthError")
	}
}

func TestRunReportsTerminalSubscribeExpired(t *testing.T) {
	fake := &fakeServer{messages: []*eventsevents.SubscribeResponse{
		{EventType: eventsevents.EventType_SubscribeExpired},
	}}
	cl := newClient(t, fake)
	errCh := make(chan error, 4)
	cl.OnError(func(err error) { trySend(errCh, err) })

	runErr := cl.Run(context.Background())
	if !errs.Is(runErr, errs.CodeAuth) {
		t.Fatalf("Run() error = %v, want auth error", runErr)
	}
	if !errors.Is(runErr, errs.ErrSubscriptionExpired) {
		t.Fatalf("Run() error = %v, want ErrSubscriptionExpired", runErr)
	}
	select {
	case err := <-errCh:
		if !errors.Is(err, errs.ErrSubscriptionExpired) {
			t.Errorf("OnError received %v, want ErrSubscriptionExpired", err)
		}
	case <-time.After(time.Second):
		t.Error("OnError was not invoked for SubscribeExpired")
	}
}

func TestRunReportsTerminalNumOfConnExceed(t *testing.T) {
	fake := &fakeServer{messages: []*eventsevents.SubscribeResponse{
		{EventType: eventsevents.EventType_NumOfConnExceed},
	}}
	cl := newClient(t, fake)
	errCh := make(chan error, 4)
	cl.OnError(func(err error) { trySend(errCh, err) })

	runErr := cl.Run(context.Background())
	if !errs.Is(runErr, errs.CodeTransport) {
		t.Fatalf("Run() error = %v, want transport error", runErr)
	}
	if !errors.Is(runErr, errs.ErrConnectionLimitExceeded) {
		t.Fatalf("Run() error = %v, want ErrConnectionLimitExceeded", runErr)
	}
	select {
	case err := <-errCh:
		if !errs.Is(err, errs.CodeTransport) {
			t.Errorf("OnError received %v, want transport error", err)
		}
		if !errors.Is(err, errs.ErrConnectionLimitExceeded) {
			t.Errorf("OnError received %v, want ErrConnectionLimitExceeded", err)
		}
	case <-time.After(time.Second):
		t.Error("OnError was not invoked for NumOfConnExceed")
	}
}

func TestNewRequiresClient(t *testing.T) {
	if _, err := New(nil); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("New(nil) error = %v, want invalid config", err)
	}
}

func TestNewWithGRPCEndpointOverride(t *testing.T) {
	cl := newCore(t)
	ev, err := New(cl, WithGRPCEndpoint("buffered-host"))
	if err != nil {
		t.Fatalf("New with explicit endpoint failed: %v", err)
	}
	t.Cleanup(func() { _ = ev.Close() })
}

func TestPayloadToDataEvent(t *testing.T) {
	resp := &SubscribeResponse{
		&eventsevents.SubscribeResponse{
			SubscribeType: 1024,
			ContentType:   "application/json",
			Payload:       `{"key":"value"}`,
			RequestId:     "req-1",
			Timestamp:     1700000000000,
		},
	}
	de := resp.ToDataEvent()
	if de.SubscribeType != 1024 {
		t.Errorf("SubscribeType = %d, want 1024", de.SubscribeType)
	}
	if de.ContentType != "application/json" {
		t.Errorf("ContentType = %q", de.ContentType)
	}
	if string(de.Payload) != `{"key":"value"}` {
		t.Errorf("Payload = %q", string(de.Payload))
	}
	if de.RequestId != "req-1" {
		t.Errorf("RequestId = %q", de.RequestId)
	}
	if de.Timestamp != 1700000000000 {
		t.Errorf("Timestamp = %d", de.Timestamp)
	}
}

func trySend[T any](ch chan T, v T) {
	select {
	case ch <- v:
	default:
	}
}

func waitSignal(t *testing.T, ch <-chan struct{}, what string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for %s", what)
	}
}

func drainErr(ch <-chan error) error {
	select {
	case err := <-ch:
		return err
	default:
		return nil
	}
}
