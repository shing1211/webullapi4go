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

// Package data provides the Webull Market Data API client.
//
// A Client is constructed from the public client:
//
//	cl, err := client.New(
//		client.WithAppKey(key),
//		client.WithAppSecret(secret),
//		client.WithSandbox(),
//	)
//	if err != nil {
//		return err
//	}
//	market := data.New(cl)
//
// The individual endpoints are grouped by area: instruments in instrument.go,
// fundamentals in fundamentals.go, futures static data in futures.go, and
// snapshot, tick, quotes, bars, footprint, NOII, screener, watchlist, options,
// and news in their own files.
//
// Numeric values are returned as strings to preserve precision, and timestamps
// are strings in the format documented by the Webull API.
package data

import (
	"context"
	"net/http"
	"net/url"
	"sync"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/display"
)

// Client exposes the Webull Market Data HTTP API. It is a thin, typed layer on
// top of [client.Client]: every request goes through the public client so that
// signing, transport, and error handling are shared with the rest of the SDK.
//
// A Client is safe for concurrent use and does not own the underlying
// [client.Client]; callers close that client themselves.
type Client struct {
	core *client.Client

	displayOnce sync.Once
	displaySvc  *display.Service
}

// New returns a market-data client bound to c. The underlying client is owned
// by the caller and is not closed by [Client.Close].
func New(c *client.Client) *Client { return &Client{core: c} }

// Core returns the underlying public client.
func (c *Client) Core() *client.Client { return c.core }

// DisplayService returns the [display.Service] for Display Solution API calls.
// It is created lazily on first use using the same app credentials as the
// underlying [client.Client].
func (c *Client) DisplayService() *display.Service {
	c.displayOnce.Do(func() {
		sandbox := c.core.Environment() == client.Sandbox
		c.displaySvc = display.NewService(
			c.core.AppKey(),
			c.core.AppSecret(),
			display.WithSandbox(sandbox),
		)
	})
	return c.displaySvc
}

// do performs a signed request and decodes the JSON response into out. path is
// the API path without a query string; query, when non-empty, is appended as an
// encoded query string. body is optional and out is optional, exactly as in
// [client.Client.Do].
//
// It is the single shared helper used by every endpoint method in this package.
func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	if len(query) > 0 {
		path += "?" + query.Encode()
	}
	return c.core.Do(ctx, method, path, body, out)
}

// get is a convenience wrapper around [Client.do] for GET requests with no body.
func (c *Client) get(ctx context.Context, path string, query url.Values, out any) error {
	return c.do(ctx, http.MethodGet, path, query, nil, out)
}

// Close releases resources held by the client. The underlying [client.Client]
// is owned by the caller and is not closed here.
func (c *Client) Close() error { return nil }
