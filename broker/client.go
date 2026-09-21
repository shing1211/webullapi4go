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

package broker

import (
	"context"
	"net/http"
	"net/url"

	"github.com/shing1211/webullapi4go/client"
)

// Client exposes the Webull Broker HTTP API. It is a thin, typed layer on top
// of [client.Client]: every request goes through the public client so that
// signing, token injection, transport, and error handling are shared with the
// rest of the SDK.
//
// A Client is safe for concurrent use and does not own the underlying
// [client.Client]; callers close that client themselves.
type Client struct {
	core *client.Client
	cfg  config
}

// New returns a broker client bound to c, configured by opts. Options are
// applied in order on top of the defaults, and nil options are ignored. The
// underlying client is owned by the caller and is not closed by [Client.Close].
func New(c *client.Client, opts ...Option) *Client {
	cfg := defaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return &Client{core: c, cfg: cfg}
}

// Core returns the underlying public client.
func (c *Client) Core() *client.Client { return c.core }

// do performs a signed request and decodes the JSON response into out. path is
// the API path without a query string; query, when non-empty, is appended as an
// encoded query string. body is optional and out is optional, exactly as in
// [client.Client.DoBroker].
//
// It is the single shared helper used by every endpoint method in this package.
func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	if len(query) > 0 {
		path += "?" + query.Encode()
	}
	return c.core.DoBroker(ctx, method, path, body, out)
}

// get is a convenience wrapper around [Client.do] for GET requests with no body.
func (c *Client) get(ctx context.Context, path string, query url.Values, out any) error {
	return c.do(ctx, http.MethodGet, path, query, nil, out)
}

// post is a convenience wrapper around [Client.do] for POST requests.
func (c *Client) post(ctx context.Context, path string, query url.Values, body, out any) error {
	return c.do(ctx, http.MethodPost, path, query, body, out)
}

// put is a convenience wrapper around [Client.do] for PUT requests.
func (c *Client) put(ctx context.Context, path string, query url.Values, body, out any) error {
	return c.do(ctx, http.MethodPut, path, query, body, out)
}

// delete is a convenience wrapper around [Client.do] for DELETE requests.
func (c *Client) delete(ctx context.Context, path string, query url.Values, body, out any) error {
	return c.do(ctx, http.MethodDelete, path, query, body, out)
}

// Close releases resources held by the client. The underlying [client.Client]
// is owned by the caller and is not closed here.
func (c *Client) Close() error { return nil }
