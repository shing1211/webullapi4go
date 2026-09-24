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

package events_test

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
	"github.com/shing1211/webullapi4go/events"
	eventsevents "github.com/shing1211/webullapi4go/gen/webull/trade/events/v1"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

// capture holds what the fake server observed for one Subscribe call.
type capture struct {
	request  *eventsevents.SubscribeRequest
	metadata metadata.MD
}

// fakeServer is an in-process EventService that records the incoming request
// and metadata, replays a fixed sequence of responses, then holds the stream
// open until the client cancels it.
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

// newCore returns a credential-only client; no network call is made here.
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

// newClient returns an events client wired to impl over an in-memory bufconn
// listener, so the test never touches the network.
func newClient(t *testing.T, impl eventsevents.EventServiceServer, opts ...events.Option) *events.Client {
	t.Helper()
	return newClientWithCore(t, newCore(t), impl, opts...)
}

func newClientWithCore(t *testing.T, core *client.Client, impl eventsevents.EventServiceServer, opts ...events.Option) *events.Client {
	t.Helper()
	lis := bufconn.Listen(1 << 20)
	srv := grpc.NewServer()
	eventsevents.RegisterEventServiceServer(srv, impl)
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)

	base := []events.Option{
		events.WithGRPCEndpoint("passthrough:///bufnet"),
		events.WithTLS(false),
		events.WithDialTimeout(5 * time.Second),
		events.WithGRPCDialOption(grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		})),
	}
	cl, err := events.New(core, append(base, opts...)...)
	if err != nil {
		t.Fatalf("events.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })
	return cl
}

// dataEvent is a decoded OnEvent callback.
type dataEvent struct {
	kind        uint32
	contentType string
	payload     string
}

func TestRunDispatchesSignedEvents(t *testing.T) {
	fake := &fakeServer{messages: []*eventsevents.SubscribeResponse{
		{EventType: eventsevents.EventType_SubscribeSuccess, Payload: `{"status":"ok"}`},
		{EventType: eventsevents.EventType_Ping},
		{
			EventType:     eventsevents.EventType(events.EventOrder),
			SubscribeType: 1,
			ContentType:   "application/json",
			Payload:       `{"order_status":"FILLED"}`,
		},
	}}
	cl := newClient(t, fake, events.WithAccounts([]string{"acct-1"}))

	connected := make(chan struct{}, 1)
	pinged := make(chan struct{}, 1)
	data := make(chan dataEvent, 4)
	errCh := make(chan error, 4)
	cl.OnConnect(func() { trySend(connected, struct{}{}) })
	cl.OnPing(func() { trySend(pinged, struct{}{}) })
	cl.OnEvent(func(kind uint32, contentType string, payload []byte) {
		trySend(data, dataEvent{kind: kind, contentType: contentType, payload: string(payload)})
	})
	cl.OnError(func(err error) { trySend(errCh, err) })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	runErr := make(chan error, 1)
	go func() { runErr <- cl.Run(ctx) }()

	waitSignal(t, connected, "OnConnect")
	waitSignal(t, pinged, "OnPing")

	var got dataEvent
	select {
	case got = <-data:
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for OnEvent (errors: %v)", drainErr(errCh))
	}
	if got.kind != events.EventOrder {
		t.Errorf("event kind = %d, want %d", got.kind, events.EventOrder)
	}
	if got.contentType != "application/json" {
		t.Errorf("content type = %q, want application/json", got.contentType)
	}
	if got.payload != `{"order_status":"FILLED"}` {
		t.Errorf("payload = %q", got.payload)
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
	if got := cap.request.GetSubscribeType(); got != uint32(events.SubscribeAll) {
		t.Errorf("subscribeType = %d, want %d", got, events.SubscribeAll)
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

func TestWithSubscribeTypes(t *testing.T) {
	fake := &fakeServer{messages: []*eventsevents.SubscribeResponse{
		{EventType: eventsevents.EventType_AuthError},
	}}
	cl := newClient(t, fake, events.WithSubscribeTypes(events.SubscribeOrder, events.SubscribeOption))
	_ = cl.Run(context.Background())

	want := uint32(events.SubscribeOrder | events.SubscribeOption)
	if got := fake.captured().request.GetSubscribeType(); got != want {
		t.Errorf("subscribeType = %d, want %d", got, want)
	}
}

func TestSubscribeTypeValues(t *testing.T) {
	cases := map[events.SubscribeType]uint32{
		events.SubscribeOrder:    1,
		events.SubscribePosition: 2,
		events.SubscribeOption:   4,
		events.SubscribeAll:      7,
	}
	for st, want := range cases {
		if uint32(st) != want {
			t.Errorf("SubscribeType(%d) = %d, want %d", st, uint32(st), want)
		}
	}
}

func TestEventKindValues(t *testing.T) {
	cases := map[uint32]uint32{
		events.EventOrder:    1024,
		events.EventPosition: 1028,
		events.EventOption:   1032,
	}
	for kind, want := range cases {
		if kind != want {
			t.Errorf("event kind = %d, want %d", kind, want)
		}
	}
}

func TestNewRequiresClient(t *testing.T) {
	if _, err := events.New(nil); !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("events.New(nil) error = %v, want invalid config", err)
	}
}

// trySend performs a non-blocking send so a handler can never block the stream
// pump.
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
