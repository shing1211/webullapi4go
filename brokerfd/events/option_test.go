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
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
)

func TestWithReconnectMaxDelay_LessThanBaseDelay(t *testing.T) {
	cl, err := client.New(client.WithAppKey("testkey"), client.WithAppSecret("testsecret"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = New(cl,
		WithGRPCEndpoint("localhost:50051"),
		WithReconnectMaxDelay(500*time.Millisecond),
	)
	if err == nil {
		t.Error("expected error when max delay is less than base delay, got nil")
	}
}

func TestWithMaxReconnectAttempts_Zero(t *testing.T) {
	cl, err := client.New(client.WithAppKey("testkey"), client.WithAppSecret("testsecret"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	c, err := New(cl,
		WithGRPCEndpoint("localhost:50051"),
		WithMaxReconnectAttempts(0),
	)
	if err != nil {
		t.Errorf("expected no error for max reconnect attempts 0 (unlimited), got: %v", err)
	}
	if c == nil {
		t.Error("expected client to be created, got nil")
		return
	}
	if c.cfg.maxReconnectAttempts != 0 {
		t.Errorf("expected maxReconnectAttempts to be 0, got %d", c.cfg.maxReconnectAttempts)
	}
}

func TestWithAccounts_EmptyList(t *testing.T) {
	cl, err := client.New(client.WithAppKey("testkey"), client.WithAppSecret("testsecret"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	c, err := New(cl,
		WithGRPCEndpoint("localhost:50051"),
		WithAccounts([]string{}),
	)
	if err != nil {
		t.Errorf("expected no error for empty accounts list, got: %v", err)
	}
	if c == nil {
		t.Error("expected client to be created, got nil")
		return
	}
	if len(c.cfg.accounts) != 0 {
		t.Errorf("expected accounts to be empty, got %d", len(c.cfg.accounts))
	}
}

func TestWithGRPCEndpoint_Empty(t *testing.T) {
	cl, err := client.New(
		client.WithAppKey("testkey"),
		client.WithAppSecret("testsecret"),
		client.WithEndpoints(client.Endpoints{HTTP: "http://localhost"}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = New(cl,
		WithGRPCEndpoint(""),
	)
	if err == nil {
		t.Error("expected error for empty endpoint, got nil")
	}
}
