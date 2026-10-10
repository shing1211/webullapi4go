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
	"errors"
	"strings"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// fakeMessage is a paho.Message whose accessors return fixed values, so the
// adapter in handleMessage can be exercised without a broker.
type fakeMessage struct {
	topic     string
	payload   []byte
	qos       byte
	retained  bool
	duplicate bool
	acked     bool
}

func (m *fakeMessage) Duplicate() bool   { return m.duplicate }
func (m *fakeMessage) Qos() byte         { return m.qos }
func (m *fakeMessage) Retained() bool    { return m.retained }
func (m *fakeMessage) Topic() string     { return m.topic }
func (m *fakeMessage) MessageID() uint16 { return 1 }
func (m *fakeMessage) Payload() []byte   { return m.payload }
func (m *fakeMessage) Ack()              { m.acked = true }

// fakeToken is a paho.Token that is not a *paho.ConnectToken, so
// connectReturnCode must report zero for it.
type fakeToken struct {
	done chan struct{}
	err  error
}

func newFakeToken(err error) *fakeToken {
	t := &fakeToken{done: make(chan struct{}), err: err}
	close(t.done)
	return t
}

func (t *fakeToken) Wait() bool                     { return true }
func (t *fakeToken) WaitTimeout(time.Duration) bool { return true }
func (t *fakeToken) Done() <-chan struct{}          { return t.done }
func (t *fakeToken) Error() error                   { return t.err }

// newTestClient builds a client that is never connected to a broker.
func newTestClient(t *testing.T) *Client {
	t.Helper()
	c, err := New(Config{Broker: "data-api.webull.hk:1883", ClientID: "test-session"})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

func TestConnackErrorMessage(t *testing.T) {
	t.Parallel()

	withCause := (&ConnackError{Code: ConnackConnectionLimit, Err: ErrConnectionLimit}).Error()
	if !strings.Contains(withCause, "105") {
		t.Errorf("Error() = %q, want it to mention the return code", withCause)
	}
	if !strings.Contains(withCause, "connection refused") {
		t.Errorf("Error() = %q, want it to mention the refusal", withCause)
	}
	if !strings.Contains(withCause, ErrConnectionLimit.Error()) {
		t.Errorf("Error() = %q, want it to include the wrapped cause", withCause)
	}

	noCause := (&ConnackError{Code: 4}).Error()
	if !strings.Contains(noCause, "4") {
		t.Errorf("Error() = %q, want it to mention the return code", noCause)
	}
	if strings.Contains(noCause, "<nil>") {
		t.Errorf("Error() = %q, want no nil cause text", noCause)
	}
}

func TestHandlePahoConnectInvokesHandler(t *testing.T) {
	t.Parallel()

	c := newTestClient(t)
	var calls int
	c.SetConnectHandler(func() { calls++ })

	c.reconnecting.Store(true)
	c.handlePahoConnect()

	if calls != 1 {
		t.Fatalf("connect handler called %d times, want 1", calls)
	}
	if c.IsReconnecting() {
		t.Fatal("IsReconnecting() = true after a connect callback, want false")
	}
}

func TestHandlePahoConnectionLostPassesCause(t *testing.T) {
	t.Parallel()

	c := newTestClient(t)
	want := errors.New("connection lost")
	var got error
	var calls int
	c.SetConnectionLostHandler(func(err error) { got = err; calls++ })

	c.handlePahoConnectionLost(want)

	if calls != 1 {
		t.Fatalf("connection-lost handler called %d times, want 1", calls)
	}
	if !errors.Is(got, want) {
		t.Fatalf("handler received %v, want %v", got, want)
	}
}

func TestHandlePahoReconnectingSetsState(t *testing.T) {
	t.Parallel()

	c := newTestClient(t)
	var calls int
	c.SetReconnectHandler(func() { calls++ })

	c.handlePahoReconnecting()

	if calls != 1 {
		t.Fatalf("reconnect handler called %d times, want 1", calls)
	}
	if !c.IsReconnecting() {
		t.Fatal("IsReconnecting() = false after a reconnecting callback, want true")
	}
}

// TestPahoCallbacksAreSuppressedAfterClose pins that a late broker callback
// cannot reach user handlers once the client is closed.
func TestPahoCallbacksAreSuppressedAfterClose(t *testing.T) {
	t.Parallel()

	c := newTestClient(t)
	var calls int
	c.SetConnectHandler(func() { calls++ })
	c.SetConnectionLostHandler(func(error) { calls++ })
	c.SetReconnectHandler(func() { calls++ })
	c.SetMessageHandler(func(Message) { calls++ })
	c.SetErrorHandler(func(error) { calls++ })

	if err := c.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	c.handlePahoConnect()
	c.handlePahoConnectionLost(errors.New("late"))
	c.handlePahoReconnecting()
	c.handleMessage(&fakeMessage{topic: "t", payload: []byte("p")})
	c.emitError(errors.New("late"))

	if calls != 0 {
		t.Fatalf("handlers called %d times after Close(), want 0", calls)
	}
}

func TestHandleMessageAdaptsEveryField(t *testing.T) {
	t.Parallel()

	c := newTestClient(t)
	var got Message
	c.SetMessageHandler(func(m Message) { got = m })

	c.handleMessage(&fakeMessage{
		topic:     "quotesAAPL",
		payload:   []byte(`{"a":1}`),
		qos:       1,
		retained:  true,
		duplicate: true,
	})

	if got.Topic != "quotesAAPL" {
		t.Errorf("Topic = %q, want %q", got.Topic, "quotesAAPL")
	}
	if string(got.Payload) != `{"a":1}` {
		t.Errorf("Payload = %q, want %q", got.Payload, `{"a":1}`)
	}
	if got.QoS != 1 {
		t.Errorf("QoS = %d, want 1", got.QoS)
	}
	if !got.Retained {
		t.Error("Retained = false, want true")
	}
	if !got.Duplicate {
		t.Error("Duplicate = false, want true")
	}
}

func TestHandleMessageWithoutHandlerIsNoOp(t *testing.T) {
	t.Parallel()

	c := newTestClient(t)
	c.handleMessage(&fakeMessage{topic: "t", payload: []byte("p")})
}

func TestEmitErrorDeliversAndIgnoresNil(t *testing.T) {
	t.Parallel()

	c := newTestClient(t)
	var got error
	var calls int
	c.SetErrorHandler(func(err error) { got = err; calls++ })

	c.emitError(nil)
	if calls != 0 {
		t.Fatalf("emitError(nil) called the handler %d times, want 0", calls)
	}

	want := errors.New("publish failed")
	c.emitError(want)
	if calls != 1 {
		t.Fatalf("emitError() called the handler %d times, want 1", calls)
	}
	if !errors.Is(got, want) {
		t.Fatalf("handler received %v, want %v", got, want)
	}
}

func TestEmitErrorWithoutHandlerIsNoOp(t *testing.T) {
	t.Parallel()

	c := newTestClient(t)
	c.emitError(errors.New("nobody is listening"))
}

func TestConnectReturnCodeForNonConnectToken(t *testing.T) {
	t.Parallel()

	if got := connectReturnCode(newFakeToken(nil)); got != 0 {
		t.Fatalf("connectReturnCode(non-ConnectToken) = %d, want 0", got)
	}
	if got := connectReturnCode(&paho.ConnectToken{}); got != 0 {
		t.Fatalf("connectReturnCode(zero ConnectToken) = %d, want 0", got)
	}
}

// TestClassifyConnectError covers the branches reachable without a broker. A
// non-zero paho return code cannot be constructed from outside the paho
// package, so the code-bearing branches are covered through the
// isConnectionLimitError text path instead.
func TestClassifyConnectError(t *testing.T) {
	t.Parallel()

	t.Run("nil error is passed through", func(t *testing.T) {
		t.Parallel()
		if got := classifyConnectError(newFakeToken(nil), nil); got != nil {
			t.Fatalf("classifyConnectError(nil) = %v, want nil", got)
		}
	})

	t.Run("unrelated error is passed through", func(t *testing.T) {
		t.Parallel()
		want := errors.New("dial tcp: timeout")
		if got := classifyConnectError(newFakeToken(want), want); !errors.Is(got, want) {
			t.Fatalf("classifyConnectError() = %v, want %v", got, want)
		}
	})

	t.Run("connection limit text becomes code 105", func(t *testing.T) {
		t.Parallel()
		limit := errors.New("connection refused: code 105")
		got := classifyConnectError(newFakeToken(limit), limit)

		var connack *ConnackError
		if !errors.As(got, &connack) {
			t.Fatalf("classifyConnectError() = %v, want a *ConnackError", got)
		}
		if connack.Code != ConnackConnectionLimit {
			t.Fatalf("code = %d, want %d", connack.Code, ConnackConnectionLimit)
		}
		if !errors.Is(got, ErrConnectionLimit) {
			t.Fatalf("classifyConnectError() = %v, want it to match ErrConnectionLimit", got)
		}
	})
}

func TestConnectRejectsInvalidStates(t *testing.T) {
	t.Parallel()

	t.Run("uninitialized client", func(t *testing.T) {
		t.Parallel()
		c := &Client{}
		err := c.Connect(context.Background())
		if err == nil || !strings.Contains(err.Error(), "not initialized") {
			t.Fatalf("Connect() on an uninitialized client = %v, want a not-initialized error", err)
		}
	})

	t.Run("closed client", func(t *testing.T) {
		t.Parallel()
		c := newTestClient(t)
		if err := c.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
		if err := c.Connect(context.Background()); !errors.Is(err, errClientClosed) {
			t.Fatalf("Connect() after Close() = %v, want errClientClosed", err)
		}
	})

	t.Run("already cancelled context", func(t *testing.T) {
		t.Parallel()
		c := newTestClient(t)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := c.Connect(ctx)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Connect() with a cancelled context = %v, want context.Canceled", err)
		}
	})
}

// TestConnectToleratesNilContext documents that a nil context is replaced with
// context.Background instead of panicking. The uninitialized client returns
// before any network attempt, so this needs no broker.
func TestConnectToleratesNilContext(t *testing.T) {
	t.Parallel()

	c := &Client{}
	//nolint:staticcheck // a nil context is the documented input under test
	err := c.Connect(nil)
	if err == nil || !strings.Contains(err.Error(), "not initialized") {
		t.Fatalf("Connect(nil) on an uninitialized client = %v, want a not-initialized error", err)
	}
}

func TestDisconnectAfterCloseIsNoOp(t *testing.T) {
	t.Parallel()

	c := newTestClient(t)
	c.Disconnect(0)

	if err := c.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	// Must not panic or block after the client is closed.
	c.Disconnect(0)
}

func TestDisconnectOnUninitializedClientIsNoOp(t *testing.T) {
	t.Parallel()

	c := &Client{}
	c.Disconnect(0)
}

func TestBrokerFallsBackWhenNormalizationFails(t *testing.T) {
	t.Parallel()

	// An unsupported scheme cannot be normalized, so Broker must return the
	// configured value rather than an empty string.
	c := &Client{cfg: Config{Broker: "ftp://host:21"}}
	if got := c.Broker(); got != "ftp://host:21" {
		t.Fatalf("Broker() = %q, want the raw configured value %q", got, "ftp://host:21")
	}
}

func TestHasScheme(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in   string
		want bool
	}{
		{"tcp://host:1883", true},
		{"ssl://host:8883", true},
		{"ws://host:8080", true},
		{"wss://host:8883/mqtt", true},
		{"tls://host:8883", true},
		{"mqtt://host:1883", true},
		{"host:1883", false},
		{"host", false},
		{"", false},
		// The separator is present even though the scheme name is empty, so
		// hasScheme reports true and normalizeBroker rejects it by scheme.
		{"://host", true},
	}

	for _, tc := range cases {
		if got := hasScheme(tc.in); got != tc.want {
			t.Errorf("hasScheme(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
