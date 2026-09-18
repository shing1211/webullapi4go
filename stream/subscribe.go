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

package stream

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/shing1211/webullapi4go/internal/errs"
)

// Streaming subscription endpoints, relative to the core client's HTTP base
// URL. These are HTTP calls; the MQTT connection only carries pushes.
const (
	subscribePath   = "/market-data/streaming/subscribe"
	unsubscribePath = "/market-data/streaming/unsubscribe"
	// maxSymbols is the documented maximum number of symbols per request.
	maxSymbols = 100
)

// Category is the Webull security type accepted by the streaming
// subscribe/unsubscribe API.
type Category string

// Supported security categories.
const (
	// CategoryUSStock selects United States stocks.
	CategoryUSStock Category = "US_STOCK"
	// CategoryUSETF selects United States ETFs.
	CategoryUSETF Category = "US_ETF"
	// CategoryHKStock selects Hong Kong stocks.
	CategoryHKStock Category = "HK_STOCK"
	// CategoryCNStock selects China Mainland A-Shares (Stock Connect).
	CategoryCNStock Category = "CN_STOCK"
)

// Valid reports whether c is a category accepted by Webull.
func (c Category) Valid() bool {
	switch c {
	case CategoryUSStock, CategoryUSETF, CategoryHKStock, CategoryCNStock:
		return true
	default:
		return false
	}
}

// SubType is a streaming data type accepted by the subscribe/unsubscribe API.
type SubType string

// Supported subscription data types.
const (
	// SubTypeQuote selects real-time order-book data.
	SubTypeQuote SubType = "QUOTE"
	// SubTypeSnapshot selects market snapshots.
	SubTypeSnapshot SubType = "SNAPSHOT"
	// SubTypeTick selects tick-by-tick trades.
	SubTypeTick SubType = "TICK"
)

// Valid reports whether s is a data type accepted by Webull.
func (s SubType) Valid() bool {
	switch s {
	case SubTypeQuote, SubTypeSnapshot, SubTypeTick:
		return true
	default:
		return false
	}
}

// SubscribeRequest is the body of an HTTP streaming subscribe call.
type SubscribeRequest struct {
	// SessionID overrides the client's session id for this call. It must match
	// the session id of the live MQTT connection; leave it empty to use
	// [Client.SessionID].
	SessionID string
	// Symbols are the security symbols to subscribe to, for example "AAPL".
	// At most 100 symbols are allowed per request.
	Symbols []string
	// Category is the security type.
	Category Category
	// SubTypes are the data types to receive.
	SubTypes []SubType
	// Grab requests that a snapshot be pushed immediately on subscription.
	Grab bool
	// Depth is the level-2 order-book depth. It is optional and defaults to
	// the server's value (10 levels); US stocks support at most 50 levels.
	Depth string
	// OvernightRequired includes the overnight session for US stocks.
	OvernightRequired bool
}

// subscribeBody is the exact wire representation of a subscribe request.
// Grab is encoded as a string, matching the published schema.
type subscribeBody struct {
	SessionID         string    `json:"session_id"`
	Symbols           []string  `json:"symbols"`
	Category          Category  `json:"category"`
	SubTypes          []SubType `json:"sub_types"`
	Grab              string    `json:"grab"`
	Depth             string    `json:"depth,omitempty"`
	OvernightRequired bool      `json:"overnight_required,omitempty"`
}

// toBody validates r and converts it to its wire representation. defaultSession
// is used when r does not name a session.
func (r SubscribeRequest) toBody(defaultSession string) (subscribeBody, error) {
	session, err := resolveSession(r.SessionID, defaultSession)
	if err != nil {
		return subscribeBody{}, err
	}
	symbols, err := normalizeSymbols(r.Symbols)
	if err != nil {
		return subscribeBody{}, err
	}
	if !r.Category.Valid() {
		return subscribeBody{}, errs.New(errs.CodeInvalidConfig, "stream: invalid category: "+string(r.Category))
	}
	subTypes, err := normalizeSubTypes(r.SubTypes)
	if err != nil {
		return subscribeBody{}, err
	}
	if err := validateDepth(r.Depth); err != nil {
		return subscribeBody{}, err
	}
	return subscribeBody{
		SessionID:         session,
		Symbols:           symbols,
		Category:          r.Category,
		SubTypes:          subTypes,
		Grab:              strconv.FormatBool(r.Grab),
		Depth:             r.Depth,
		OvernightRequired: r.OvernightRequired,
	}, nil
}

// UnsubscribeRequest is the body of an HTTP streaming unsubscribe call.
type UnsubscribeRequest struct {
	// SessionID overrides the client's session id for this call. Leave it
	// empty to use [Client.SessionID].
	SessionID string
	// Symbols are the security symbols to unsubscribe from. They are ignored
	// when UnsubscribeAll is true.
	Symbols []string
	// Category is the security type. It is ignored when UnsubscribeAll is
	// true.
	Category Category
	// SubTypes are the data types to stop receiving. They are ignored when
	// UnsubscribeAll is true.
	SubTypes []SubType
	// UnsubscribeAll cancels every subscription for the session; when set,
	// the other selection fields may be empty.
	UnsubscribeAll bool
}

// unsubscribeBody is the exact wire representation of an unsubscribe request.
type unsubscribeBody struct {
	SessionID      string    `json:"session_id"`
	Symbols        []string  `json:"symbols,omitempty"`
	Category       Category  `json:"category,omitempty"`
	SubTypes       []SubType `json:"sub_types,omitempty"`
	UnsubscribeAll bool      `json:"unsubscribe_all,omitempty"`
}

// toBody validates r and converts it to its wire representation.
func (r UnsubscribeRequest) toBody(defaultSession string) (unsubscribeBody, error) {
	session, err := resolveSession(r.SessionID, defaultSession)
	if err != nil {
		return unsubscribeBody{}, err
	}
	if r.UnsubscribeAll {
		return unsubscribeBody{SessionID: session, UnsubscribeAll: true}, nil
	}
	symbols, err := normalizeSymbols(r.Symbols)
	if err != nil {
		return unsubscribeBody{}, err
	}
	if !r.Category.Valid() {
		return unsubscribeBody{}, errs.New(errs.CodeInvalidConfig, "stream: invalid category: "+string(r.Category))
	}
	subTypes, err := normalizeSubTypes(r.SubTypes)
	if err != nil {
		return unsubscribeBody{}, err
	}
	return unsubscribeBody{
		SessionID: session,
		Symbols:   symbols,
		Category:  r.Category,
		SubTypes:  subTypes,
	}, nil
}

// resolveSession returns the explicit session id, falling back to the client's
// session id.
func resolveSession(explicit, fallback string) (string, error) {
	if s := strings.TrimSpace(explicit); s != "" {
		return s, nil
	}
	if fallback == "" {
		return "", errs.New(errs.CodeInvalidConfig, "stream: session id is required")
	}
	return fallback, nil
}

// normalizeSymbols trims, drops empty entries, and bounds the symbol list.
func normalizeSymbols(in []string) ([]string, error) {
	out := make([]string, 0, len(in))
	for _, s := range in {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		return nil, errs.New(errs.CodeInvalidConfig, "stream: at least one symbol is required")
	}
	if len(out) > maxSymbols {
		return nil, errs.New(errs.CodeInvalidConfig, "stream: too many symbols (maximum 100)")
	}
	return out, nil
}

// normalizeSubTypes trims and validates the requested data types.
func normalizeSubTypes(in []SubType) ([]SubType, error) {
	out := make([]SubType, 0, len(in))
	for _, s := range in {
		t := SubType(strings.ToUpper(strings.TrimSpace(string(s))))
		if !t.Valid() {
			return nil, errs.New(errs.CodeInvalidConfig, "stream: invalid sub type: "+string(s))
		}
		out = append(out, t)
	}
	if len(out) == 0 {
		return nil, errs.New(errs.CodeInvalidConfig, "stream: at least one sub type is required")
	}
	return out, nil
}

// validateDepth checks an optional order-book depth.
func validateDepth(depth string) error {
	if depth == "" {
		return nil
	}
	n, err := strconv.Atoi(depth)
	if err != nil || n <= 0 {
		return errs.New(errs.CodeInvalidConfig, "stream: depth must be a positive integer")
	}
	return nil
}

// Subscribe starts streaming data for the requested symbols with
// POST /market-data/streaming/subscribe. The MQTT connection must already be
// established with [Client.Connect], and an access token must be available
// (see [client.Client.EnsureToken]).
//
// A successful subscription is recorded in the client's registry. Webull does
// not restore subscriptions across reconnects, so when the MQTT connection is
// re-established the client re-issues the active subscriptions automatically
// (see [WithAutoResubscribe]); callers do not need to call Subscribe again
// after a reconnect. Subscribing the same symbol and data type twice is
// idempotent and does not create a duplicate registration.
func (c *Client) Subscribe(ctx context.Context, req SubscribeRequest) error {
	if c.core == nil {
		return errs.New(errs.CodeInvalidConfig, "stream: client is not initialized")
	}
	body, err := req.toBody(c.cfg.sessionID)
	if err != nil {
		return err
	}
	if err := c.core.Do(ctx, http.MethodPost, subscribePath, body, nil); err != nil {
		return err
	}
	c.subs.add(body)
	return nil
}

// Unsubscribe stops streaming data for the requested symbols with
// POST /market-data/streaming/unsubscribe. Set
// [UnsubscribeRequest.UnsubscribeAll] to cancel every subscription for the
// session.
//
// A successful unsubscribe removes the matching entries from the client's
// registry, so an unsubscribed symbol is not re-issued after a reconnect.
func (c *Client) Unsubscribe(ctx context.Context, req UnsubscribeRequest) error {
	if c.core == nil {
		return errs.New(errs.CodeInvalidConfig, "stream: client is not initialized")
	}
	body, err := req.toBody(c.cfg.sessionID)
	if err != nil {
		return err
	}
	if err := c.core.Do(ctx, http.MethodPost, unsubscribePath, body, nil); err != nil {
		return err
	}
	c.subs.remove(body)
	return nil
}
