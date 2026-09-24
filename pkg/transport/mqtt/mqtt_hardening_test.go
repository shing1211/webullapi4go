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
	"sync"
	"sync/atomic"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type mqttHardeningFakeToken struct {
	done chan struct{}
	err  error
	once sync.Once
}

func (t *mqttHardeningFakeToken) Wait() bool {
	<-t.done
	return true
}

func (t *mqttHardeningFakeToken) WaitTimeout(time.Duration) bool {
	select {
	case <-t.done:
		return true
	case <-time.After(time.Second):
		return false
	}
}

func (t *mqttHardeningFakeToken) Done() <-chan struct{} { return t.done }
func (t *mqttHardeningFakeToken) Error() error          { return t.err }

func (t *mqttHardeningFakeToken) complete() {
	t.once.Do(func() { close(t.done) })
}

type mqttHardeningFakePahoClient struct {
	token           *mqttHardeningFakeToken
	open            atomic.Bool
	connectCalls    atomic.Int32
	disconnectCalls atomic.Int32
}

func (f *mqttHardeningFakePahoClient) IsConnected() bool      { return f.open.Load() }
func (f *mqttHardeningFakePahoClient) IsConnectionOpen() bool { return f.open.Load() }
func (f *mqttHardeningFakePahoClient) Connect() paho.Token {
	f.connectCalls.Add(1)
	f.open.Store(true)
	f.token.complete()
	return f.token
}
func (f *mqttHardeningFakePahoClient) Disconnect(uint) {
	f.disconnectCalls.Add(1)
	f.open.Store(false)
}
func (f *mqttHardeningFakePahoClient) Publish(string, byte, bool, interface{}) paho.Token { return nil }
func (f *mqttHardeningFakePahoClient) Subscribe(string, byte, paho.MessageHandler) paho.Token {
	return nil
}
func (f *mqttHardeningFakePahoClient) SubscribeMultiple(map[string]byte, paho.MessageHandler) paho.Token {
	return nil
}
func (f *mqttHardeningFakePahoClient) Unsubscribe(...string) paho.Token     { return nil }
func (f *mqttHardeningFakePahoClient) AddRoute(string, paho.MessageHandler) {}
func (f *mqttHardeningFakePahoClient) OptionsReader() paho.ClientOptionsReader {
	return paho.ClientOptionsReader{}
}

func TestMQTTHardeningConnectCancellationDoesNotStartAConnection(t *testing.T) {
	fake := &mqttHardeningFakePahoClient{token: &mqttHardeningFakeToken{done: make(chan struct{})}}
	c := &Client{pc: fake}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := c.Connect(ctx); err != context.Canceled {
		t.Fatalf("Connect() error = %v, want context.Canceled", err)
	}
	if got := fake.connectCalls.Load(); got != 0 {
		t.Fatalf("Connect() calls = %d, want 0", got)
	}
}

func TestMQTTHardeningCloseIsTerminalAndIdempotent(t *testing.T) {
	fake := &mqttHardeningFakePahoClient{token: &mqttHardeningFakeToken{done: make(chan struct{})}}
	c := &Client{pc: fake}
	if err := c.Connect(context.Background()); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	var messages atomic.Int32
	c.SetMessageHandler(func(Message) { messages.Add(1) })
	if err := c.Close(); err != nil {
		t.Fatalf("first Close() error = %v", err)
	}
	if err := c.Close(); err != nil {
		t.Fatalf("second Close() error = %v", err)
	}
	if got := fake.disconnectCalls.Load(); got != 1 {
		t.Fatalf("Disconnect() calls = %d, want 1", got)
	}
	if c.IsConnected() {
		t.Fatal("IsConnected() = true after Close()")
	}

	c.handleMessage(nil)
	if got := messages.Load(); got != 0 {
		t.Fatalf("messages after Close() = %d, want 0", got)
	}
	if err := c.Connect(context.Background()); err == nil {
		t.Fatal("Connect() after Close() = nil, want error")
	}
	if got := fake.connectCalls.Load(); got != 1 {
		t.Fatalf("Connect() calls after Close() = %d, want 1", got)
	}
}
