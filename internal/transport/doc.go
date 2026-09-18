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

// Package transport implements the thin HTTP executor used by the public client
// package. It joins request paths onto a base URL, applies the default
// User-Agent, and delegates the round trip to an *http.Client, which owns
// timeouts and connection pooling.
//
// Request signing, JSON encoding, and typed error mapping live in the public
// client package and in internal/auth and internal/errs. The resilience
// middleware (retry, rate limiting, and circuit breaking) lives in
// internal/resilience.
package transport
