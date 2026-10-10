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

package trade

import (
	"math"
	"strconv"
)

// config is the resolved trade client configuration. It is populated by
// [Option] functions during [New] and treated as immutable afterwards.
type config struct {
	// maxOrderNotional, when non-empty, is the advisory cap on a single
	// order's notional value, as a decimal string.
	maxOrderNotional string
	// maxOrderQuantity, when non-empty, is the advisory cap on a single
	// order's quantity, as a decimal string.
	maxOrderQuantity string
	// autoClientOrderID, when true, causes [Client.PlaceOrder] and
	// [Client.BatchPlaceOrder] to derive a stable client order identifier for
	// any order whose ClientOrderID is empty.
	autoClientOrderID bool
}

// defaultConfig returns the trade client defaults: both order guardrails are
// disabled (empty).
func defaultConfig() config { return config{} }

// Option mutates the [config] resolved by [New]. Options are applied in order,
// so a later option overrides an earlier one.
type Option func(*config)

// WithMaxOrderNotional sets an advisory cap on the notional value of a single
// order, expressed as a non-negative decimal string such as "2500.00". The
// default is empty, which disables the cap.
//
// The cap does not cover multi-leg option orders: each leg is priced
// separately, so the order has no single top-level notional and the cap is
// skipped. WithMaxOrderQuantity still applies to multi-leg orders.
//
// The guardrail is configuration only in this release; the order methods
// enforce it before an order is built. WithMaxOrderNotional panics when v is
// not a non-negative, finite decimal number.
func WithMaxOrderNotional(v string) Option {
	notional := normalizeGuardrail("WithMaxOrderNotional", v)
	return func(c *config) { c.maxOrderNotional = notional }
}

// WithMaxOrderQuantity sets an advisory cap on the quantity of a single order,
// expressed as a non-negative decimal string such as "100" or "2.5". The
// default is empty, which disables the cap.
//
// The guardrail is configuration only in this release; the order methods
// enforce it before an order is built. WithMaxOrderQuantity panics when v is
// not a non-negative, finite decimal number.
func WithMaxOrderQuantity(v string) Option {
	quantity := normalizeGuardrail("WithMaxOrderQuantity", v)
	return func(c *config) { c.maxOrderQuantity = quantity }
}

// normalizeGuardrail validates a decimal guardrail value. The empty string
// disables the guardrail and is returned unchanged; any other value must parse
// as a finite float64 and be non-negative, otherwise normalizeGuardrail panics.
func normalizeGuardrail(name, v string) string {
	if v == "" {
		return ""
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) || f < 0 {
		panic("trade: " + name + ": " + strconv.Quote(v) + " is not a non-negative decimal number")
	}
	return v
}

// WithAutoClientOrderID controls whether [Client.PlaceOrder] and
// [Client.BatchPlaceOrder] derive and assign a client order identifier to any
// order whose ClientOrderID is empty. The default is false (disabled). Derived
// identifiers are stable for the same logical request across retries, and the
// caller's request slice is not modified.
func WithAutoClientOrderID(v bool) Option {
	return func(c *config) { c.autoClientOrderID = v }
}
