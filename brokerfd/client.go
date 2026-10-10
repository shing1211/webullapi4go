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

package brokerfd

import (
	"context"
	"net/http"
	"net/url"

	"github.com/shing1211/webullapi4go/client"
)

// Client provides access to the Broker FD (US) HTTP API via the shared transport.
// It wraps a [client.Client] and delegates all HTTP calls to it, using the Broker
// host rather than the core host. See the do method for why.
type Client struct {
	core *client.Client
}

// New returns a new Broker FD client backed by the supplied shared [client.Client].
func New(c *client.Client) *Client { return &Client{core: c} }

// Core returns the underlying shared [client.Client] used for HTTP transport.
func (c *Client) Core() *client.Client { return c.core }

// do executes an HTTP request with the given method, path, query parameters, body,
// and populates the response into out. It appends query parameters to the path.
//
// The request goes to the Broker host, not the core host. Broker FD is a Broker API
// surface, so it is reached with [client.Client.DoBroker] exactly as the sibling
// broker/ module is; routing this package through the core transport sent every
// Broker FD request to a host that does not serve the surface, which surfaced as a
// 404 indistinguishable from a genuinely absent endpoint. That is a code-structure
// fact rather than an empirical one, so it needed no credential to establish: the
// transport, the endpoint field, and the sibling module's use of it all already
// existed.
func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	if len(query) > 0 {
		path += "?" + query.Encode()
	}
	return c.core.DoBroker(ctx, method, path, body, out)
}

// get issues a GET request to the Broker FD API.
func (c *Client) get(ctx context.Context, path string, query url.Values, out any) error {
	return c.do(ctx, http.MethodGet, path, query, nil, out)
}

// post issues a POST request to the Broker FD API.
func (c *Client) post(ctx context.Context, path string, query url.Values, body, out any) error {
	return c.do(ctx, http.MethodPost, path, query, body, out)
}

// Close is a placeholder for future resource cleanup. Currently returns nil.
func (c *Client) Close() error { return nil }
