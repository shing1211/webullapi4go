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
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	runErr := make(chan error, 1)
	go func() { runErr <- cl.Run(ctx) }()

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

	cancel()
	select {
	case err := <-runErr:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Run() error = %v, want context.Canceled", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after context cancellation")
	}
}

func TestRunStopsOnContextCancel(t *testing.T) {
	fake := &fakeServer{messages: []*eventsevents.SubscribeResponse{
		{EventType: eventsevents.EventType_SubscribeSuccess},
	}}
	cl := newClient(t, fake)
	connected := make(chan struct{}, 1)
	cl.OnConnect(func() { trySend(connected, struct{}{}) })

	ctx, cancel := context.WithCancel(context.Background())
	runErr := make(chan error, 1)
	go func() { runErr <- cl.Run(ctx) }()

	waitSignal(t, connected, "OnConnect")
	cancel()

	select {
	case err := <-runErr:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Run() error = %v, want context.Canceled", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after context cancellation")
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
	select {
	case err := <-errCh:
		if !errs.Is(err, errs.CodeTransport) {
			t.Errorf("OnError received %v, want transport error", err)
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
	_, err := New(cl, WithGRPCEndpoint("buffered-host"))
	if err != nil {
		t.Fatalf("New with explicit endpoint failed: %v", err)
	}
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
