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

// Package transport is a deprecated alias for [github.com/shing1211/webullapi4go/pkg/transport].
// New code should import that package directly. This package will be removed in v3.
package transport

import (
	"net/http"

	pkgtrans "github.com/shing1211/webullapi4go/pkg/transport"
)

// Doer is the minimal interface for an HTTP executor.
//
// Deprecated: use [pkgtrans.Doer].
type Doer = pkgtrans.Doer

// ErrNilHTTPClient is returned by [New] when no *http.Client is supplied.
//
// Deprecated: use [pkgtrans.ErrNilHTTPClient].
var ErrNilHTTPClient = pkgtrans.ErrNilHTTPClient

// Client is a thin HTTP executor.
//
// Deprecated: use [pkgtrans.Client].
type Client = pkgtrans.Client

// New returns a Client that sends requests rooted at baseURL using hc.
//
// Deprecated: use [pkgtrans.New].
func New(baseURL string, hc *http.Client, userAgent string) (*Client, error) {
	return pkgtrans.New(baseURL, hc, userAgent)
}
