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
// It wraps a [client.Client] and delegates all HTTP calls to it.
type Client struct {
	core *client.Client
}

// New returns a new Broker FD client backed by the supplied shared [client.Client].
func New(c *client.Client) *Client { return &Client{core: c} }

// Core returns the underlying shared [client.Client] used for HTTP transport.
func (c *Client) Core() *client.Client { return c.core }

// do executes an HTTP request with the given method, path, query parameters, body,
// and populates the response into out. It appends query parameters to the path.
func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	if len(query) > 0 {
		path += "?" + query.Encode()
	}
	return c.core.Do(ctx, method, path, body, out)
}

// get issues a GET request to the Broker FD API.
func (c *Client) get(ctx context.Context, path string, query url.Values, out any) error {
	return c.do(ctx, http.MethodGet, path, query, nil, out)
}

// post issues a POST request to the Broker FD API.
func (c *Client) post(ctx context.Context, path string, query url.Values, body, out any) error {
	return c.do(ctx, http.MethodPost, path, query, body, out)
}

// put issues a PUT request to the Broker FD API.
func (c *Client) put(ctx context.Context, path string, query url.Values, body, out any) error {
	return c.do(ctx, http.MethodPut, path, query, body, out)
}

// delete issues a DELETE request to the Broker FD API.
func (c *Client) delete(ctx context.Context, path string, query url.Values, body, out any) error {
	return c.do(ctx, http.MethodDelete, path, query, body, out)
}

// Close is a placeholder for future resource cleanup. Currently returns nil.
func (c *Client) Close() error { return nil }
