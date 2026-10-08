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

	errs "github.com/shing1211/webullapi4go/pkg/errors"
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
// a cached access token is attached, the host that is signed is the host that
// is sent, and the configured rate limiter, circuit breaker, interceptors,
// hooks, logger, spans, and metrics apply.
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
	ctx, _ = ensureCorrelationID(ctx)

	if err := c.ensureAutoToken(ctx, reqPath); err != nil {
		return nil, err
	}

	bodyBytes, err := marshalRequestBody(body)
	if err != nil {
		return nil, err
	}

	result, err := c.attemptStreamWithNumber(ctx, method, reqPath, query, bodyBytes, 1)
	if err != nil {
		if result.response != nil && result.response.Body != nil {
			_ = result.response.Body.Close()
		}
		return nil, err
	}
	if result.response == nil {
		return nil, errs.New(errs.CodeTransport, "stream interceptor returned no response")
	}
	return result.response, nil
}

func (c *Client) attemptStreamWithNumber(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte, attempt int) (attemptResult, error) {
	return c.runAttempt(ctx, method, reqPath, attempt, true, func(operationCtx context.Context) (attemptResult, error) {
		return c.executeStream(operationCtx, method, reqPath, query, bodyBytes)
	})
}

func (c *Client) executeStream(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte) (attemptResult, error) {
	var result attemptResult
	req, err := c.buildSignedRequest(ctx, method, reqPath, query, bodyBytes)
	if err != nil {
		return result, err
	}

	resp, err := c.transport.Do(req)
	if resp != nil {
		result.status = resp.StatusCode
		result.response = resp
	}
	if err != nil {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		result.response = nil
		return result, errs.Wrap(errs.CodeTransport, method+" "+reqPath, err)
	}
	if resp == nil {
		return result, errs.New(errs.CodeTransport, method+" "+reqPath+" returned nil response")
	}

	if resp.Body == nil {
		result.response = nil
		return result, errs.New(errs.CodeTransport, method+" "+reqPath+" returned nil response body")
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		data, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		result.response = nil
		if readErr != nil {
			return result, errs.Wrap(errs.CodeTransport, "reading response body", readErr)
		}
		return result, errs.FromHTTPStatus(resp.StatusCode, data)
	}
	c.maybeUpdateClockOffset(resp.Header)
	return result, nil
}
