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

package main

import (
	"fmt"
	"math"
	"sync/atomic"
	"time"
)

// DefaultSymbol is the instrument every synthesised symbol parameter resolves
// to. The sandbox carries AAPL and nothing else (AGENTS.md, "Known
// constraints"), so a symbol resolved to any other ticker names an instrument
// the sandbox does not have, and the endpoint answers 404 for a reason that
// says nothing about the endpoint. It is a request value the probe chooses, not
// a value the server sent, so it carries nothing into a committed file.
const DefaultSymbol = "AAPL"

// defaultCount is the value a declared numeric parameter resolves to before its
// published range is applied. One is the smallest count a paged endpoint can
// honour and it is in the page's range wherever the page publishes no bound, so
// it is a request rather than an error for the common parameter (`count`,
// `limit`, `pageSize`) without costing the endpoint a page of work.
const defaultCount = 1

// correlationSequence numbers the correlation keys this process has issued. It
// is atomic because the census resolves parameters from whatever goroutine is
// walking an endpoint, and two goroutines sharing a sequence number would
// collide exactly where the race detector is not looking.
var correlationSequence atomic.Uint64

// ParamSpec is what a reference page documents about one required query
// parameter: its declared JSON type, the closed set of values it permits, the
// example the page itself publishes, and the bounds it publishes.
//
// Type is the page's own type name, which is not the same thing as the kind
// SkeletonKind reports for a reduced body: this is a string the documentation
// uses for a request, and it is read here as documentation. A page may also use
// `integer`, which is not a JSON type name, so both names reach the numeric
// class.
//
// The remaining fields are separate rather than merged because an absent field
// and a zero one are different documents. HasExample distinguishes a page that
// published an example from a page whose example was the empty string, and the
// Minimum and Maximum pointers distinguish a page that published a bound of
// zero from a page that published none.
type ParamSpec struct {
	// Type is the declared type name: "string", "integer", "number" or
	// "boolean" in the corpus. An unrecognised name resolves to an error rather
	// than to a guess.
	Type string
	// Enum is the closed set the page declares, in the order it declares it. The
	// first member is used, so the order is part of the contract: a page that
	// lists a default first is read as naming one.
	Enum []any
	// Example is the value the page publishes, returned as written.
	Example any
	// HasExample reports whether Example was published.
	HasExample bool
	// Minimum and Maximum are the published bounds, or nil where the page
	// publishes none.
	Minimum *float64
	Maximum *float64
}

// Param decides the value to send for one required query parameter, from what
// the reference page documents about it. It returns the value to place in the
// query, and a non-nil error when no class of parameter can fill it.
//
// The classes are tried in a fixed order, and the order is the contract:
//
//  1. the page's own example, which is the only value the page says is valid;
//  2. the first member of a declared closed set, which the page says is
//     acceptable;
//  3. account_id, which is refused;
//  4. a symbol name, resolved to DefaultSymbol;
//  5. a declared number, resolved to 1 inside the published range;
//  6. a declared boolean, resolved to true;
//  7. a caller-supplied correlation key, resolved to a value unique to this
//     call;
//  8. otherwise an error.
//
// Steps 1 and 2 come first because the page is authoritative about a parameter
// it documents, and a synthesised value is a guess about one it does not.
//
// The Go type of a resolved value is part of the contract, because the caller
// puts it on the wire and the page declared a JSON type: a numeric parameter
// resolves to an int, a boolean to the bool true, a string to a string, a
// plural symbol to a []string, and a page value to the page's own value with its
// type unchanged. A []string is returned as a list and the caller encodes it in
// the spelling the page documents, which this function does not know: the
// spelling is a property of the request being built, not of the parameter.
//
// An error is the signal the census records as `blocked: unresolvable
// parameter`, and it is why the last step is a refusal rather than a default.
// The probe cannot distinguish "the endpoint did not answer" from "the probe
// never sent a request the endpoint could answer", so a class that cannot be
// filled must not invent a value to send: a fabricated value reaches the
// endpoint and is rejected for a reason that is recorded as the endpoint's own
// answer. The returned error names the parameter, so the census row says which
// one blocked the call.
func Param(name string, spec ParamSpec) (any, error) {
	if spec.HasExample {
		return spec.Example, nil
	}
	if len(spec.Enum) > 0 {
		return spec.Enum[0], nil
	}
	// The identity of the caller is issued by the server, and this branch is
	// ahead of every type-based default for that reason: an account_id that
	// resolved to a fabricated string would send 34 account-scoped endpoints at
	// an account that does not exist, and each would be recorded as an endpoint
	// that failed to answer rather than as a call the probe could not make. A
	// discovered account id belongs here as a parameter class, not in the page's
	// schema.
	if name == "account_id" {
		return nil, fmt.Errorf("live-probe: %q is issued by the server: discover an account and thread its id in", name)
	}
	// The instrument. The plural is resolved to a list of one because the name
	// is the only place the page's arity is visible: the spec carries no field
	// for it, and a plural name over a singular value is a request the endpoint
	// answers with a type error.
	switch name {
	case "symbol", "series_symbol":
		return DefaultSymbol, nil
	case "symbols":
		return []string{DefaultSymbol}, nil
	}
	switch spec.Type {
	case "integer", "number":
		return clampCount(spec.Minimum, spec.Maximum), nil
	case "boolean":
		return true, nil
	case "string":
		// A correlation key is the one class the server accepts any value in,
		// so it is the one class where a synthesised value is always right. It
		// has to be unique: a repeated key is rejected as a duplicate, which is
		// an error attributed to the endpoint rather than to the key.
		if name == "client_request_id" || name == "client_order_id" {
			return newCorrelationID(), nil
		}
	}
	return nil, fmt.Errorf("live-probe: no value is synthesable for parameter %q of declared type %q", name, spec.Type)
}

// newCorrelationID returns a value no other call in this process returns, for a
// caller-supplied correlation key.
//
// The counter is what makes two calls in one run differ, and it is what the
// probe would rely on if two calls were issued in the same nanosecond. The
// timestamp is what separates two runs, so a key from a previous run is not
// resubmitted and rejected as a duplicate by a server that remembers them.
func newCorrelationID() string {
	return fmt.Sprintf("live-probe-%d-%d", correlationSequence.Add(1), time.Now().UnixNano())
}

// clampCount resolves a declared numeric parameter to an integer inside the
// range the page publishes, preferring the smallest value an endpoint can serve.
//
// Either published bound is applied, not only a range with both bounds: a page
// that documents `count` as at most 100 and a page that documents it as at least
// 5 are each a range the value must respect, and a bound left unapplied sends a
// request the endpoint rejects for a reason the probe would record as the
// endpoint's answer. With no bound published the value is defaultCount.
//
// A published bound need not be a whole number, and the value sent is an int, so
// the integer is clamped a second time against the bounds rounded outward. When
// the published range contains no integer at all, no buildable value satisfies
// it and the nearest one is sent; the probe prefers a request the endpoint can
// answer to a refusal that leaves the endpoint untested.
func clampCount(minimum, maximum *float64) int {
	v := float64(defaultCount)
	if minimum != nil && v < *minimum {
		v = *minimum
	}
	if maximum != nil && v > *maximum {
		v = *maximum
	}
	n := int(v)
	if minimum != nil && float64(n) < *minimum {
		n = int(math.Ceil(*minimum))
	}
	if maximum != nil && float64(n) > *maximum {
		n = int(math.Floor(*maximum))
	}
	return n
}
