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
	"errors"
	"testing"
)

func TestConfigValidate(t *testing.T) {
	base := Config{Broker: "data-api.webull.hk:1883", ClientID: "abc"}

	if err := base.Validate(); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
	if err := (Config{ClientID: "abc"}).Validate(); err == nil {
		t.Error("Validate() with empty broker = nil, want error")
	}
	if err := (Config{Broker: base.Broker}).Validate(); err == nil {
		t.Error("Validate() with empty client id = nil, want error")
	}
	if err := (Config{Broker: "ftp://host:1883", ClientID: "abc"}).Validate(); err == nil {
		t.Error("Validate() with unsupported scheme = nil, want error")
	}
	if err := (Config{Broker: "tcp://", ClientID: "abc"}).Validate(); err == nil {
		t.Error("Validate() with hostless broker = nil, want error")
	}
}

func TestWithDefaults(t *testing.T) {
	got := (Config{Broker: "host:1883", ClientID: "abc"}).withDefaults()
	if got.KeepAlive != DefaultKeepAlive {
		t.Errorf("KeepAlive = %v, want %v", got.KeepAlive, DefaultKeepAlive)
	}
	if got.ConnectTimeout != DefaultConnectTimeout {
		t.Errorf("ConnectTimeout = %v, want %v", got.ConnectTimeout, DefaultConnectTimeout)
	}
	if got.MessageChannelDepth != DefaultMessageChannelDepth {
		t.Errorf("MessageChannelDepth = %d, want %d", got.MessageChannelDepth, DefaultMessageChannelDepth)
	}

	set := Config{
		Broker:              "host:1883",
		ClientID:            "abc",
		KeepAlive:           5,
		ConnectTimeout:      6,
		MessageChannelDepth: 7,
	}.withDefaults()
	if set.KeepAlive != 5 || set.ConnectTimeout != 6 || set.MessageChannelDepth != 7 {
		t.Errorf("withDefaults() overwrote explicit values: %+v", set)
	}
}

func TestNormalizeBroker(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{in: "data-api.webull.hk:1883", want: "tcp://data-api.webull.hk:1883"},
		{in: "tcp://data-api.webull.hk:1883", want: "tcp://data-api.webull.hk:1883"},
		{in: "wss://data-api.webull.hk:8883/mqtt", want: "wss://data-api.webull.hk:8883/mqtt"},
		{in: "ws://localhost:8080", want: "ws://localhost:8080"},
		{in: "tls://host:8883", want: "tls://host:8883"},
		{in: "ftp://host:21", wantErr: true},
		{in: "tcp://", wantErr: true},
		{in: "http://host", wantErr: true},
	}
	for _, tt := range tests {
		got, err := normalizeBroker(tt.in)
		if tt.wantErr {
			if err == nil {
				t.Errorf("normalizeBroker(%q) = %q, nil; want error", tt.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("normalizeBroker(%q) error = %v", tt.in, err)
			continue
		}
		if got != tt.want {
			t.Errorf("normalizeBroker(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestNewBuildsClientWithoutConnecting(t *testing.T) {
	c, err := New(Config{Broker: "data-api.webull.hk:1883", ClientID: "session-1"})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if c.IsConnected() {
		t.Error("IsConnected() = true, want false before Connect")
	}
	if c.IsReconnecting() {
		t.Error("IsReconnecting() = true, want false before Connect")
	}
	if got := c.Broker(); got != "tcp://data-api.webull.hk:1883" {
		t.Errorf("Broker() = %q, want %q", got, "tcp://data-api.webull.hk:1883")
	}
	if err := c.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestConnackErrorClassification(t *testing.T) {
	limit := &ConnackError{Code: ConnackConnectionLimit, Err: ErrConnectionLimit}
	if !errors.Is(limit, ErrConnectionLimit) {
		t.Error("code 105 error does not match ErrConnectionLimit")
	}
	if !errors.Is(limit, ErrConnectionRefused) {
		t.Error("code 105 error does not match ErrConnectionRefused")
	}

	var asConnack *ConnackError
	if !errors.As(limit, &asConnack) || asConnack.Code != ConnackConnectionLimit {
		t.Errorf("errors.As() did not recover the ConnackError: %v", asConnack)
	}

	other := &ConnackError{Code: 4}
	if errors.Is(other, ErrConnectionLimit) {
		t.Error("non-105 code matches ErrConnectionLimit")
	}
	if !errors.Is(other, ErrConnectionRefused) {
		t.Error("rejection code does not match ErrConnectionRefused")
	}
}

func TestIsConnectionLimitError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "code text", err: errors.New("connection refused: code 105"), want: true},
		{name: "json code", err: errors.New(`{"code":"105","message":"limit"}`), want: true},
		{name: "limit text", err: errors.New("exceeds connection limit"), want: true},
		{name: "sentinel", err: ErrConnectionLimit, want: true},
		{name: "address contains 105", err: errors.New("dial tcp 10.0.1.105:1883: timeout"), want: false},
		{name: "unrelated", err: errors.New("dial tcp: timeout"), want: false},
	}
	for _, tt := range tests {
		if got := isConnectionLimitError(tt.err); got != tt.want {
			t.Errorf("isConnectionLimitError(%v) = %v, want %v", tt.err, got, tt.want)
		}
	}
}
