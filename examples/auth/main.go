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

// Command auth creates or reuses a Webull access token and prints its status.
//
// Credentials are read from the environment (see [client.WithEnv]):
//
//	WEBULL_APP_KEY=...
//	WEBULL_APP_SECRET=...
//	WEBULL_ENVIRONMENT=sandbox
//
// Run it with:
//
//	go run ./examples/auth
package main

import (
	"context"
	"log"
	"time"

	"github.com/shing1211/webullapi4go/client"
)

func main() {
	// WithEnv reads WEBULL_APP_KEY, WEBULL_APP_SECRET, WEBULL_REGION, and
	// WEBULL_ENVIRONMENT. New fails when the credentials are missing.
	cl, err := client.New(client.WithEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = cl.Close() }()

	// EnsureToken reuses a cached, still-valid token or creates a new one. In
	// the sandbox the token is NORMAL immediately; in production the call waits
	// for the Webull App 2FA verification window.
	tok, err := cl.EnsureToken(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	expires := "unset"
	if !tok.ExpiresAt.IsZero() {
		expires = tok.ExpiresAt.Format(time.RFC3339)
	}
	log.Printf("token status=%s expires_at=%s", tok.Status, expires)
}
