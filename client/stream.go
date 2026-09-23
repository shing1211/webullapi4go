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

package client

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/shing1211/webullapi4go/pkg/errors"
)

// DoStream sends a signed request and returns the raw, open response for
// endpoints that stream their result, such as Server-Sent Events, instead of
// buffering and JSON-decoding it.
//
// DoStream applies the same request pipeline as [Client.Do]: path is resolved
// against the configured base URL and may carry a query string, body is
// serialized as compact JSON without HTML escaping (a []byte or
// [encoding/json.RawMessage] is sent verbatim, nil sends no body), the
// x-version header follows the same built-in per-path defaults and overrides,
// a cached access token is attached, and the host that is signed is the host
// that is sent. The configured rate limiter waits before the request and the
// circuit breaker gates it, exactly as for a buffered call.
//
// Unlike [Client.Do], DoStream never retries. A streaming response is a
// long-lived connection that is consumed incrementally and is not safe to
// replay, so the configured retry policy is deliberately bypassed.
//
// On a non-2xx response the body is read, closed, and returned as a typed
// *[errs.Error] classified by [errs.FromHTTPStatus]. On success the returned
// response has an open body that the caller owns: the caller must read and
// close resp.Body. DoStream never closes a successful response body.
func (c *Client) DoStream(ctx context.Context, method, path string, body any) (*http.Response, error) {
	reqPath, query, err := splitPathQuery(path)
	if err != nil {
		return nil, err
	}
	query = normalizeQuery(query)

	if err := c.ensureAutoToken(ctx, reqPath); err != nil {
		return nil, err
	}

	bodyBytes, err := marshalRequestBody(body)
	if err != nil {
		return nil, err
	}

	if c.cfg.rateLimiter != nil {
		if err := c.cfg.rateLimiter.Wait(ctx, reqPath); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return nil, ctxErr
			}
			return nil, errs.Wrap(errs.CodeRateLimited, "rate limiter", err)
		}
	}
	if c.cfg.breaker != nil && !c.cfg.breaker.Allow() {
		return nil, errs.Wrap(errs.CodeTransport, "circuit breaker open", ErrCircuitOpen)
	}

	resp, err := c.sendStream(ctx, method, reqPath, query, bodyBytes)
	if c.cfg.breaker != nil {
		c.recordBreaker(err)
	}
	return resp, err
}

// sendStream performs a single streaming request attempt. It builds and signs
// the request with the shared pipeline and returns the open response, or a
// typed error. The caller owns a successful response body.
func (c *Client) sendStream(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte) (*http.Response, error) {
	req, err := c.buildSignedRequest(ctx, method, reqPath, query, bodyBytes)
	if err != nil {
		return nil, err
	}

	resp, err := c.transport.Do(req)
	if err != nil {
		return nil, errs.Wrap(errs.CodeTransport, method+" "+reqPath, err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		data, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			return nil, errs.Wrap(errs.CodeTransport, "reading response body", readErr)
		}
		return nil, errs.FromHTTPStatus(resp.StatusCode, data)
	}
	return resp, nil
}
