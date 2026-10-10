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

//go:build tools

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/shing1211/webullapi4go/client"
)

func main() {
	appKey := os.Getenv("WEBULL_APP_KEY")
	appSecret := os.Getenv("WEBULL_APP_SECRET")
	if appKey == "" || appSecret == "" {
		fmt.Fprintf(os.Stderr, "WEBULL_APP_KEY and WEBULL_APP_SECRET must be set\n")
		os.Exit(1)
	}

	cl, err := client.New(
		client.WithAppKey(appKey),
		client.WithAppSecret(appSecret),
		client.WithRegion(client.HK),
		client.WithEnvironment(client.Sandbox),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "client.New: %v\n", err)
		os.Exit(1)
	}
	defer cl.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tok, err := cl.EnsureToken(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "EnsureToken: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(tok.Value)
}
