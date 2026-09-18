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
	"sort"
	"strings"
	"sync"
)

// subKey identifies one active subscription: a data type for one symbol in one
// category. It is the unit of idempotency for [subscriptionRegistry]:
// subscribing the same key twice does not create a second entry.
type subKey struct {
	category Category
	symbol   string
	subType  SubType
}

// subOptions are the request-level options a subscription was registered with.
// When the same key is subscribed again with different options, the most recent
// options win.
type subOptions struct {
	grab              bool
	depth             string
	overnightRequired bool
}

// group buckets subscriptions that can be re-issued in the same subscribe
// request: they share a category and request-level options.
type group struct {
	category Category
	options  subOptions
}

// subscriptionRegistry tracks the subscriptions that are active for a session
// so they can be re-issued after a reconnect. Webull does not restore
// subscriptions across reconnects.
//
// A subscriptionRegistry is safe for concurrent use. The zero value is an empty
// registry. It only tracks the client's own session; a [SubscribeRequest] that
// overrides the session id is still recorded but re-issued with the client's
// session id.
type subscriptionRegistry struct {
	mu   sync.Mutex
	subs map[subKey]subOptions
}

// add records every (symbol, sub type) pair in body as active.
func (r *subscriptionRegistry) add(body subscribeBody) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.subs == nil {
		r.subs = make(map[subKey]subOptions)
	}
	opts := subOptions{
		grab:              body.Grab == "true",
		depth:             body.Depth,
		overnightRequired: body.OvernightRequired,
	}
	for _, symbol := range body.Symbols {
		for _, subType := range body.SubTypes {
			r.subs[subKey{category: body.Category, symbol: symbol, subType: subType}] = opts
		}
	}
}

// remove deletes the subscriptions named by body. A body with UnsubscribeAll
// set clears the registry.
func (r *subscriptionRegistry) remove(body unsubscribeBody) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if body.UnsubscribeAll {
		r.subs = nil
		return
	}
	for _, symbol := range body.Symbols {
		for _, subType := range body.SubTypes {
			delete(r.subs, subKey{category: body.Category, symbol: symbol, subType: subType})
		}
	}
}

// empty reports whether no subscriptions are active.
func (r *subscriptionRegistry) empty() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.subs) == 0
}

// requests reconstructs the minimum set of [SubscribeRequest]s that exactly
// restores the active subscriptions. Subscriptions sharing a category, request
// options, and sub-type set are merged into a single request, and symbols with
// different sub-type sets are never combined (which would over-subscribe).
// Every active (symbol, sub type) pair appears in exactly one returned request.
func (r *subscriptionRegistry) requests() []SubscribeRequest {
	r.mu.Lock()
	groups := make(map[group]map[string]map[SubType]struct{})
	for key, opts := range r.subs {
		g := group{category: key.category, options: opts}
		symbols := groups[g]
		if symbols == nil {
			symbols = make(map[string]map[SubType]struct{})
			groups[g] = symbols
		}
		subTypes := symbols[key.symbol]
		if subTypes == nil {
			subTypes = make(map[SubType]struct{})
			symbols[key.symbol] = subTypes
		}
		subTypes[key.subType] = struct{}{}
	}
	r.mu.Unlock()

	reqs := make([]SubscribeRequest, 0, len(groups))
	for g, symbols := range groups {
		bySet := make(map[string]*requestCombo)
		for symbol, set := range symbols {
			subTypes := sortedSubTypes(set)
			key := subTypeSetKey(subTypes)
			combo := bySet[key]
			if combo == nil {
				combo = &requestCombo{subTypes: subTypes}
				bySet[key] = combo
			}
			combo.symbols = append(combo.symbols, symbol)
		}
		for _, combo := range bySet {
			sort.Strings(combo.symbols)
			reqs = append(reqs, SubscribeRequest{
				Symbols:           combo.symbols,
				Category:          g.category,
				SubTypes:          combo.subTypes,
				Grab:              g.options.grab,
				Depth:             g.options.depth,
				OvernightRequired: g.options.overnightRequired,
			})
		}
	}
	sortSubscribeRequests(reqs)
	return reqs
}

// requestCombo accumulates the symbols that share one sub-type set.
type requestCombo struct {
	subTypes []SubType
	symbols  []string
}

// sortedSubTypes returns the sub types of set in ascending order so the
// resulting request is deterministic.
func sortedSubTypes(set map[SubType]struct{}) []SubType {
	out := make([]SubType, 0, len(set))
	for subType := range set {
		out = append(out, subType)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// subTypeSetKey returns a canonical key for an already sorted sub-type slice.
func subTypeSetKey(subTypes []SubType) string {
	parts := make([]string, len(subTypes))
	for i, subType := range subTypes {
		parts[i] = string(subType)
	}
	return strings.Join(parts, ",")
}

// sortSubscribeRequests orders requests deterministically so that a reconnect
// re-issues them in a stable order.
func sortSubscribeRequests(reqs []SubscribeRequest) {
	sort.Slice(reqs, func(i, j int) bool {
		a, b := reqs[i], reqs[j]
		if a.Category != b.Category {
			return a.Category < b.Category
		}
		if a.Depth != b.Depth {
			return a.Depth < b.Depth
		}
		if a.Grab != b.Grab {
			return !a.Grab
		}
		if a.OvernightRequired != b.OvernightRequired {
			return !a.OvernightRequired
		}
		if ka, kb := subTypeSetKey(a.SubTypes), subTypeSetKey(b.SubTypes); ka != kb {
			return ka < kb
		}
		return strings.Join(a.Symbols, ",") < strings.Join(b.Symbols, ",")
	})
}
