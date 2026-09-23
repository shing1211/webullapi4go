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
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/shing1211/webullapi4go/events"
	eventsevents "github.com/shing1211/webullapi4go/gen/webull/trade/events/v1"
	"github.com/shing1211/webullapi4go/pkg/errors"
)

// countingServer records how many Subscribe calls it received and delegates
// each to behavior, which may be swapped per test.
type countingServer struct {
	eventsevents.UnimplementedEventServiceServer

	subscribes atomic.Int32
	behavior   func(attempt int32, stream grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error
}

func (s *countingServer) Subscribe(_ *eventsevents.SubscribeRequest, stream grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error {
	n := s.subscribes.Add(1)
	if s.behavior == nil {
		return nil
	}
	return s.behavior(n, stream)
}

// TestReconnectsAndResubscribes verifies that a clean end of stream triggers a
// reconnect with a fresh Subscribe request and that OnConnect fires again.
func TestReconnectsAndResubscribes(t *testing.T) {
	srv := &countingServer{}
	srv.behavior = func(attempt int32, stream grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error {
		if err := stream.Send(&eventsevents.SubscribeResponse{EventType: eventsevents.EventType_SubscribeSuccess}); err != nil {
			return err
		}
		if attempt == 1 {
			if err := stream.Send(dataResponse(events.EventOrder, "application/json", orderEventJSON)); err != nil {
				return err
			}
			return nil
		}
		<-stream.Context().Done()
		return stream.Context().Err()
	}

	cl := newClient(t, srv,
		events.WithReconnectBaseDelay(time.Millisecond),
		events.WithReconnectMaxDelay(5*time.Millisecond),
		events.WithMaxReconnectAttempts(10),
	)

	var connects atomic.Int32
	orders := make(chan *events.OrderEvent, 4)
	cl.OnConnect(func() { connects.Add(1) })
	cl.OnOrder(func(ev *events.OrderEvent) { trySend(orders, ev) })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	runErr := make(chan error, 1)
	go func() { runErr <- cl.Run(ctx) }()

	deadline := time.After(5 * time.Second)
	for connects.Load() < 2 {
		select {
		case <-deadline:
			t.Fatalf("client did not reconnect: connects=%d subscribes=%d", connects.Load(), srv.subscribes.Load())
		case <-time.After(5 * time.Millisecond):
		}
	}

	select {
	case ev := <-orders:
		if ev.OrderID != "ORDER-1" {
			t.Errorf("OrderID = %q, want ORDER-1", ev.OrderID)
		}
	case <-time.After(time.Second):
		t.Fatal("order event from the first stream was not delivered")
	}

	if got := srv.subscribes.Load(); got != 2 {
		t.Errorf("Subscribe calls = %d, want 2", got)
	}

	cancel()
	select {
	case err := <-runErr:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Run() error = %v, want context.Canceled", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after cancellation")
	}
}

// TestTerminalErrorStopsReconnect verifies that an AuthError ends Run without
// retrying, so the client cannot loop forever on bad credentials.
func TestTerminalErrorStopsReconnect(t *testing.T) {
	srv := &countingServer{}
	srv.behavior = func(_ int32, stream grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error {
		return stream.Send(&eventsevents.SubscribeResponse{EventType: eventsevents.EventType_AuthError})
	}

	cl := newClient(t, srv,
		events.WithReconnectBaseDelay(time.Millisecond),
		events.WithMaxReconnectAttempts(3),
	)
	errCh := make(chan error, 4)
	cl.OnError(func(err error) { trySend(errCh, err) })

	runErr := cl.Run(context.Background())
	if !errs.Is(runErr, errs.CodeAuth) {
		t.Fatalf("Run() error = %v, want an auth error", runErr)
	}
	if got := srv.subscribes.Load(); got != 1 {
		t.Errorf("Subscribe calls = %d, want 1 (no retry on a terminal error)", got)
	}
	select {
	case err := <-errCh:
		if !errs.Is(err, errs.CodeAuth) {
			t.Errorf("OnError = %v, want an auth error", err)
		}
	case <-time.After(time.Second):
		t.Error("OnError was not invoked")
	}
}

// TestConnectionLimitIsTerminal verifies that NumOfConnExceed is not retried
// even though it is classified as a transport error.
func TestConnectionLimitIsTerminal(t *testing.T) {
	srv := &countingServer{}
	srv.behavior = func(_ int32, stream grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error {
		return stream.Send(&eventsevents.SubscribeResponse{EventType: eventsevents.EventType_NumOfConnExceed})
	}

	cl := newClient(t, srv,
		events.WithReconnectBaseDelay(time.Millisecond),
		events.WithMaxReconnectAttempts(3),
	)
	runErr := cl.Run(context.Background())
	if !errs.Is(runErr, errs.CodeTransport) {
		t.Fatalf("Run() error = %v, want a transport error", runErr)
	}
	if got := srv.subscribes.Load(); got != 1 {
		t.Errorf("Subscribe calls = %d, want 1 (no retry on a connection-limit event)", got)
	}
}

// TestMaxReconnectAttemptsBoundsRetries verifies that a permanently failing
// stream gives up after the configured number of attempts.
func TestMaxReconnectAttemptsBoundsRetries(t *testing.T) {
	srv := &countingServer{} // behavior nil: every Subscribe ends immediately
	cl := newClient(t, srv,
		events.WithReconnectBaseDelay(time.Millisecond),
		events.WithReconnectMaxDelay(2*time.Millisecond),
		events.WithMaxReconnectAttempts(2),
	)

	runErr := cl.Run(context.Background())
	if runErr == nil {
		t.Fatal("Run() returned nil, want an exhausted-retries error")
	}
	if !errs.Is(runErr, errs.CodeTransport) {
		t.Errorf("Run() error = %v, want a transport error", runErr)
	}
	if got := srv.subscribes.Load(); got != 3 {
		t.Errorf("Subscribe calls = %d, want 3 (initial plus 2 retries)", got)
	}
}

// TestAutoReconnectDisabled verifies that WithAutoReconnect(false) returns after
// a single clean stream end.
func TestAutoReconnectDisabled(t *testing.T) {
	srv := &countingServer{}
	srv.behavior = func(_ int32, stream grpc.ServerStreamingServer[eventsevents.SubscribeResponse]) error {
		return stream.Send(&eventsevents.SubscribeResponse{EventType: eventsevents.EventType_SubscribeSuccess})
	}

	cl := newClient(t, srv,
		events.WithAutoReconnect(false),
		events.WithReconnectBaseDelay(time.Millisecond),
	)
	if err := cl.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v, want nil on a clean end", err)
	}
	if got := srv.subscribes.Load(); got != 1 {
		t.Errorf("Subscribe calls = %d, want 1", got)
	}
}

// TestDialTimeoutIsRetryable verifies that a connect timeout is classified as a
// transient transport failure and retried rather than treated as terminal.
func TestDialTimeoutIsRetryable(t *testing.T) {
	cl, err := events.New(newCore(t),
		events.WithGRPCEndpoint("passthrough:///unreachable.invalid:1"),
		events.WithTLS(false),
		events.WithDialTimeout(50*time.Millisecond),
		events.WithReconnectBaseDelay(time.Millisecond),
		events.WithReconnectMaxDelay(2*time.Millisecond),
		events.WithMaxReconnectAttempts(2),
	)
	if err != nil {
		t.Fatalf("events.New() error = %v", err)
	}
	defer func() { _ = cl.Close() }()

	runErr := cl.Run(context.Background())
	if !errs.Is(runErr, errs.CodeTransport) {
		t.Fatalf("Run() error = %v, want a transport error", runErr)
	}
}
