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

package data

import (
	"context"
	"net/http"
)

// Display Solution streaming endpoints.
//
// TODO(ds): Confirm exact paths against US sandbox.
const (
	pathDSSubscribe   = "/openapi/market-data/streaming/subscribe"
	pathDSUnsubscribe = "/openapi/market-data/streaming/unsubscribe"
)

// DSSubscribeRequest is the request body for [Client.DSSubscribe].
type DSSubscribeRequest struct {
	Symbols []string `json:"symbols"`
}

// DSUnsubscribeRequest is the request body for [Client.DSUnsubscribe].
type DSUnsubscribeRequest struct {
	Symbols []string `json:"symbols"`
}

// DSSubscribe subscribes to real-time streaming data for the given symbols
// via the Display Solution endpoint.
//
// TODO(ds): Confirm exact paths against US sandbox.
func (c *Client) DSSubscribe(ctx context.Context, req DSSubscribeRequest) error {
	return c.DisplayService().Do(ctx, http.MethodPost, pathDSSubscribe, nil, req, nil)
}

// DSUnsubscribe unsubscribes from real-time streaming data for the given
// symbols via the Display Solution endpoint.
//
// TODO(ds): Confirm exact paths against US sandbox.
func (c *Client) DSUnsubscribe(ctx context.Context, req DSUnsubscribeRequest) error {
	return c.DisplayService().Do(ctx, http.MethodPost, pathDSUnsubscribe, nil, req, nil)
}
