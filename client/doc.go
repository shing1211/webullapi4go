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

// Package client is the public entry point to the Webull OpenAPI SDK. It owns
// the HTTP client, request signing, compact-JSON encoding, and typed error
// mapping, and it is the type passed to API-specific packages such as data.
//
// Create a client with [New] and the functional [Option]s, then issue signed
// requests with [Client.Do]:
//
//	cl, err := client.New(
//		client.WithAppKey(key),
//		client.WithAppSecret(secret),
//		client.WithRegion(client.HK),
//		client.WithEnvironment(client.Sandbox),
//	)
//	if err != nil {
//		return err
//	}
//	defer cl.Close()
//
//	var out json.RawMessage
//	if err := cl.Do(ctx, http.MethodGet, "/openapi/account/list", nil, &out); err != nil {
//		return err
//	}
//
// Every exported type in this package is defined here, so callers outside the
// module never need to import an internal package.
package client
