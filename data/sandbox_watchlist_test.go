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

package data_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/data"
)

// newT44SandboxClient returns a sandbox market-data client for the T4.4
// endpoints, skipping the test when the sandbox environment variables are
// absent. The caller owns the underlying client via t.Cleanup.
//
// It runs only when WEBULL_SANDBOX=1 and WEBULL_APP_KEY / WEBULL_APP_SECRET are
// set; never commit the sandbox credentials to the repository.
func newT44SandboxClient(t *testing.T) (*client.Client, *data.Client) {
	t.Helper()
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
	if baseURL := os.Getenv("WEBULL_BASE_URL"); baseURL != "" {
		opts = append(opts, client.WithBaseURL(baseURL))
	}

	cl, err := client.New(opts...)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })
	return cl, data.New(cl)
}

// TestSandboxWatchlists exercises the read-only watchlist endpoints against the
// shared Webull sandbox. The sandbox accounts are public and shared, so the
// test never creates, renames, or deletes a watchlist.
func TestSandboxWatchlists(t *testing.T) {
	cl, market := newT44SandboxClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if _, err := cl.EnsureToken(ctx); err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}

	watchlists, err := market.GetWatchlists(ctx)
	if err != nil {
		t.Fatalf("GetWatchlists() error = %v", err)
	}
	t.Logf("sandbox watchlists: %d", len(watchlists))

	for _, wl := range watchlists {
		got, err := market.GetWatchlistInstruments(ctx, wl.WatchlistID)
		if err != nil {
			t.Logf("GetWatchlistInstruments(%s) unavailable: %v", wl.WatchlistID, err)
			continue
		}
		t.Logf("watchlist %s (%s): %d instruments", wl.WatchlistID, wl.Name, len(got.Instruments))
	}
}
