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

// Package transport provides the HTTP transport layer used by the public client.
// It joins request paths onto a base URL, applies the default User-Agent, and
// delegates the round trip to an *http.Client, which owns timeouts and connection
// pooling.
//
// Request signing, JSON encoding, and typed error mapping live in the client
// package and in pkg/auth and pkg/errors. The resilience middleware (retry,
// rate limiting, and circuit breaking) lives in pkg/resilience.
package transport

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ErrNilHTTPClient is returned by [New] when no *http.Client is supplied.
var ErrNilHTTPClient = errors.New("transport: http client is required")

// Doer is the minimal interface for an HTTP executor. It is satisfied by
// [*Client]. Consumers may substitute a mock for testing.
type Doer interface {
	// Do sends req and returns the raw response. The caller owns the response
	// body.
	Do(req *http.Request) (*http.Response, error)
	// CloseIdleConnections releases idle connections held by the underlying
	// transport. It is safe to call more than once.
	CloseIdleConnections()
}

// Client is a thin HTTP executor. It joins request paths onto a fixed base URL,
// applies a default User-Agent, and delegates the actual round trip to the
// configured *http.Client, which owns timeouts and connection pooling.
//
// A Client is safe for concurrent use.
type Client struct {
	baseURL   *url.URL
	http      *http.Client
	userAgent string
}

// New returns a Client that sends requests rooted at baseURL using hc. The
// baseURL must be an absolute http or https URL (for example
// "https://api.webull.hk"); its path, if any, is preserved as a prefix.
// Requests carry userAgent when it is non-empty. hc must not be nil.
func New(baseURL string, hc *http.Client, userAgent string) (*Client, error) {
	if hc == nil {
		return nil, ErrNilHTTPClient
	}
	u, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, errors.New("transport: base URL must be absolute and include a host")
	}
	return &Client{baseURL: u, http: hc, userAgent: userAgent}, nil
}

// BaseURL returns the resolved base URL.
func (c *Client) BaseURL() *url.URL {
	u := *c.baseURL
	return &u
}

// HTTPClient returns the underlying *http.Client.
func (c *Client) HTTPClient() *http.Client { return c.http }

// NewRequest builds a request for method against the client's base URL. path is
// appended to the base path and query, when non-nil, becomes the encoded query
// string. body, when non-empty, is sent verbatim. The returned request is bound
// to ctx.
func (c *Client) NewRequest(ctx context.Context, method, path string, query url.Values, body []byte) (*http.Request, error) {
	u := *c.baseURL
	full, err := url.Parse(joinPath(c.baseURL.Path, path))
	if err != nil {
		return nil, err
	}
	u.Path = full.Path
	u.RawQuery = full.RawQuery
	if query != nil {
		u.RawQuery = query.Encode()
	}
	var reader io.Reader
	if len(body) > 0 {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), reader)
	if err != nil {
		return nil, err
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	return req, nil
}

// Do sends req and returns the raw response. The caller owns the response body.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.http.Do(req)
}

// CloseIdleConnections releases idle connections held by the underlying
// transport. It is safe to call more than once. It does not close connections
// that are still in use.
func (c *Client) CloseIdleConnections() {
	c.http.CloseIdleConnections()
}

// joinPath appends path to base, ensuring exactly one "/" between them.
func joinPath(base, path string) string {
	if path == "" {
		return base
	}
	if strings.HasPrefix(path, "/") {
		return strings.TrimRight(base, "/") + path
	}
	return strings.TrimRight(base, "/") + "/" + path
}
