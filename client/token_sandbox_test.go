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

package client_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
)

// TestSandboxToken creates an access token against the shared Webull sandbox and
// asserts that it is NORMAL. Sandbox tokens are activated automatically (no 2FA
// step), so this exercises the create/poll/cache path end to end.
//
// It runs only when WEBULL_SANDBOX=1 and WEBULL_APP_KEY / WEBULL_APP_SECRET are
// set; otherwise it skips so the default test run stays hermetic. Never commit
// the sandbox credentials to the repository.
func TestSandboxToken(t *testing.T) {
	if os.Getenv("WEBULL_SANDBOX") != "1" {
		t.Skip("set WEBULL_SANDBOX=1, WEBULL_APP_KEY and WEBULL_APP_SECRET to run the sandbox test")
	}
	appKey := os.Getenv("WEBULL_APP_KEY")
	appSecret := os.Getenv("WEBULL_APP_SECRET")
	if appKey == "" || appSecret == "" {
		t.Skip("set WEBULL_APP_KEY and WEBULL_APP_SECRET to run the sandbox test")
	}

	opts := []client.Option{
		client.WithAppKey(appKey),
		client.WithAppSecret(appSecret),
		client.WithRegion(client.HK),
		client.WithEnvironment(client.Sandbox),
	}
	// WEBULL_BASE_URL lets a run target a different region's sandbox (for
	// example https://api.sandbox.webull.com) without changing the test.
	if baseURL := os.Getenv("WEBULL_BASE_URL"); baseURL != "" {
		opts = append(opts, client.WithBaseURL(baseURL))
	}

	cl, err := client.New(opts...)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tok, err := cl.EnsureToken(ctx)
	if err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}
	t.Logf("sandbox token status=%s expires_at=%s value=%.8s...",
		tok.Status, tok.ExpiresAt.Format(time.RFC3339), tok.Value)
	if tok.Status != client.TokenStatusNormal {
		t.Fatalf("sandbox token status = %q, want %q", tok.Status, client.TokenStatusNormal)
	}
	if cl.AccessToken() != tok.Value {
		t.Errorf("AccessToken() = %q, want cached %q", cl.AccessToken(), tok.Value)
	}

	// A second call must reuse the cached token rather than create a new one.
	again, err := cl.EnsureToken(ctx)
	if err != nil {
		t.Fatalf("second EnsureToken() error = %v", err)
	}
	if again.Value != tok.Value {
		t.Fatalf("second EnsureToken() value = %q, want cached %q", again.Value, tok.Value)
	}
}
