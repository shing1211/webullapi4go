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
	"reflect"
	"strings"
	"testing"
)

func TestParam(t *testing.T) {
	// Class: caller-supplied correlation key. Any unique value is accepted.
	got, err := Param("client_request_id", ParamSpec{Type: "string"})
	if err != nil || got == "" {
		t.Errorf("client_request_id = %v, %v; want a synthesised non-empty value", got, err)
	}
	// Two calls must differ, or a second endpoint's correlation key would collide.
	a, _ := Param("client_request_id", ParamSpec{Type: "string"})
	b, _ := Param("client_request_id", ParamSpec{Type: "string"})
	if a == b {
		t.Error("two correlation keys are identical")
	}
	// Class: enum. First value wins, and the page's own example is preferred.
	got, _ = Param("category", ParamSpec{Type: "string", Enum: []any{"US_STOCK", "HK_STOCK"}})
	if got != "US_STOCK" {
		t.Errorf("enum first = %v, want US_STOCK", got)
	}
	got, _ = Param("category", ParamSpec{Type: "string", Enum: []any{"US_STOCK"}, Example: "HK_STOCK", HasExample: true})
	if got != "HK_STOCK" {
		t.Errorf("documented example = %v, want HK_STOCK: the page is authoritative", got)
	}
	// Class: declared scalar, respecting the published range.
	got, _ = Param("count", ParamSpec{Type: "integer", Minimum: ptr(1.0), Maximum: ptr(20.0)})
	if n, ok := got.(int); !ok || n < 1 || n > 20 {
		t.Errorf("count = %v, want an int within [1,20]", got)
	}
	// Class: domain constant.
	got, _ = Param("symbol", ParamSpec{Type: "string"})
	if got != DefaultSymbol {
		t.Errorf("symbol = %v, want %q", got, DefaultSymbol)
	}
	// Class: caller identity. Not synthesisable — the server issues it.
	if _, err := Param("account_id", ParamSpec{Type: "string"}); err == nil {
		t.Error("account_id resolved without discovery; it must be threaded in, not guessed")
	}
}

func TestParamRejectsUnresolvable(t *testing.T) {
	if _, err := Param("mystery", ParamSpec{}); err == nil {
		t.Error("a spec with no type, enum or example resolved; it must not")
	}
}

// ptr is the published-bounds helper the cases above use. The bounds are
// pointers in ParamSpec because an absent bound and a bound of zero are
// different documents, and a plain float64 could not say which was declared.
func ptr(f float64) *float64 { return &f }

// The brief's case pins `symbol` alone. The class is the instrument, and the
// corpus spells it three ways, so all three are pinned: a name the class
// covers that resolves to an error is reported as an unbuildable request, and
// one that resolves to a bare string where the page declares a list is a
// request the endpoint answers with a 400 the probe would misread.
func TestParamSymbolNames(t *testing.T) {
	cases := []struct {
		name  string
		param string
		spec  ParamSpec
		want  any
	}{
		{"symbol", "symbol", ParamSpec{Type: "string"}, DefaultSymbol},
		{"series_symbol", "series_symbol", ParamSpec{Type: "string"}, DefaultSymbol},
		// The plural is the one place the resolution returns a container. The
		// page declares the arity in the name, since the spec carries no field
		// for it, so `symbols` is a list of one and `symbol` is a bare string.
		{"symbols", "symbols", ParamSpec{Type: "string"}, []string{DefaultSymbol}},
		{"symbols declared as an array", "symbols", ParamSpec{Type: "array"}, []string{DefaultSymbol}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Param(tc.param, tc.spec)
			if err != nil {
				t.Fatalf("err = %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}

// The correlation-key class is the only one whose value has to differ from call
// to call, and the brief pins one name of the two. Both are pinned here
// because the property is the uniqueness rather than the name: a probe that
// sends the same key to two endpoints gets the second request rejected as a
// duplicate, which is a false negative about the endpoint rather than about the
// key.
func TestParamCorrelationKeysAreUnique(t *testing.T) {
	for _, name := range []string{"client_request_id", "client_order_id"} {
		t.Run(name, func(t *testing.T) {
			first, err := Param(name, ParamSpec{Type: "string"})
			if err != nil {
				t.Fatalf("err = %v", err)
			}
			key, ok := first.(string)
			if !ok || key == "" {
				t.Fatalf("got %#v, want a non-empty string", first)
			}
			second, err := Param(name, ParamSpec{Type: "string"})
			if err != nil {
				t.Fatalf("err = %v", err)
			}
			if key == second {
				t.Errorf("two calls both produced %q", key)
			}
		})
	}
}

// The class is a declared scalar, and the brief pins the range that every
// declared bound sits inside. These are the ranges a page actually publishes:
// one bound, and bounds that do not admit the default. The value sent has to be
// inside the published range in each, or the endpoint answers 400 and the probe
// records a parameter error as an endpoint that did not answer.
func TestParamNumericRange(t *testing.T) {
	cases := []struct {
		name string
		spec ParamSpec
		want int
	}{
		{"no published bound", ParamSpec{Type: "integer"}, 1},
		{"default is inside the range", ParamSpec{Type: "integer", Minimum: ptr(1.0), Maximum: ptr(20.0)}, 1},
		{"maximum below the default", ParamSpec{Type: "integer", Maximum: ptr(0.0)}, 0},
		{"maximum above the default is left alone", ParamSpec{Type: "integer", Maximum: ptr(100.0)}, 1},
		// A minimum above the default is the case a "clamp only when both are
		// published" reading gets wrong: it sends 1 to a page that documents
		// `count` as 5 or more, and the endpoint rejects the request for a
		// reason that says nothing about the endpoint.
		{"minimum above the default", ParamSpec{Type: "integer", Minimum: ptr(5.0)}, 5},
		{"minimum and maximum both above the default", ParamSpec{Type: "integer", Minimum: ptr(5.0), Maximum: ptr(9.0)}, 5},
		// A bound that is not a whole number cannot be sent as one, so the
		// nearest integer the range admits is sent. A fractional maximum the
		// default already sits inside is left alone: the bound is respected by
		// sending 1, and rounding it outward would send a larger request for
		// nothing. A range that admits no integer at all is the documented
		// limit of the class: no value the probe can build satisfies it, and the
		// closest one is sent rather than an out-of-range 1.
		{"fractional minimum is rounded up", ParamSpec{Type: "integer", Minimum: ptr(1.5)}, 2},
		{"fractional maximum above the default is left alone", ParamSpec{Type: "integer", Maximum: ptr(2.5)}, 1},
		{"a range with no integer sends the closest one", ParamSpec{Type: "integer", Minimum: ptr(0.2), Maximum: ptr(0.5)}, 0},
		{"a narrow range with no integer sends the closest one", ParamSpec{Type: "integer", Minimum: ptr(1.5), Maximum: ptr(1.8)}, 1},
		// A declared number resolves through the same class as a declared
		// integer: both are JSON numbers, and the range is the page's.
		{"number is the same class", ParamSpec{Type: "number", Minimum: ptr(1.0), Maximum: ptr(20.0)}, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Param("count", tc.spec)
			if err != nil {
				t.Fatalf("err = %v", err)
			}
			n, ok := got.(int)
			if !ok {
				t.Fatalf("got %#v (%T), want an int", got, got)
			}
			if n != tc.want {
				t.Errorf("got %d, want %d", n, tc.want)
			}
		})
	}
}

// The boolean class is one value, and which one is a judgement the brief fixes:
// true, because every documented flag defaults to off and a probe asking for
// the richer response is asking for the flag the page advertises. The property
// pinned here is the Go type, because a string "true" on the wire is a
// different request from a JSON true.
func TestParamBoolean(t *testing.T) {
	got, err := Param("extendedHours", ParamSpec{Type: "boolean"})
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if b, ok := got.(bool); !ok || !b {
		t.Errorf("got %#v (%T), want the bool true", got, got)
	}
}

// A page's example and a closed set are returned as the page wrote them, with
// no coercion to the declared type. Coercing them would be the probe deciding
// the page meant something else, and a value of the wrong JSON type is a
// request the endpoint rejects.
func TestParamReturnsPageValuesVerbatim(t *testing.T) {
	cases := []struct {
		name string
		spec ParamSpec
		want any
	}{
		{"numeric enum member", ParamSpec{Type: "integer", Enum: []any{5.5}}, 5.5},
		{"boolean enum member", ParamSpec{Type: "boolean", Enum: []any{false}}, false},
		{"numeric example", ParamSpec{Type: "number", Example: 0.25, HasExample: true}, 0.25},
		{"false example is not an absent example", ParamSpec{Type: "boolean", Example: false, HasExample: true}, false},
		{"empty-string example is an example", ParamSpec{Type: "string", Example: "", HasExample: true}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Param("category", tc.spec)
			if err != nil {
				t.Fatalf("err = %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %#v (%T), want %#v (%T)", got, got, tc.want, tc.want)
			}
		})
	}
}

// The ordering is the requirement, so the case the ordering produces is pinned
// rather than left to be discovered in a run: a page that publishes an example
// for account_id is answered with that example, because the page is
// authoritative about a value it documents, and the account-identity refusal
// below it exists for a page that publishes nothing. An account-scoped
// endpoint called this way is then reported against an account that does not
// exist, which the census records as the endpoint's own answer rather than as a
// blocked call. The fix is a discovered account threaded in ahead of this call,
// not a change to the class.
func TestParamExampleOutranksAccountIdentity(t *testing.T) {
	got, err := Param("account_id", ParamSpec{Type: "string", Example: "1234567890", HasExample: true})
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if got != "1234567890" {
		t.Errorf("got %#v, want the page's own example", got)
	}
}

// Every class the resolver cannot fill must be an error rather than a guess,
// because a guess is a request that reaches the endpoint and is answered for
// the wrong reason. The corpus is the undeclared cases a documented page
// produces: a string with no closed set and no example, a numeric parameter
// whose type the page omits, a structured parameter, and a name the classes do
// not cover.
func TestParamUnresolvableClasses(t *testing.T) {
	cases := []struct {
		name string
		spec ParamSpec
	}{
		{"no type, no enum, no example", ParamSpec{}},
		{"string with nothing to choose from", ParamSpec{Type: "string"}},
		{"numeric with no declared type", ParamSpec{Minimum: ptr(1.0), Maximum: ptr(9.0)}},
		{"structured parameter", ParamSpec{Type: "object"}},
		{"array parameter that is not a symbol list", ParamSpec{Type: "array"}},
		{"a type the classes do not name", ParamSpec{Type: "null"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := Param("period", tc.spec); err == nil {
				t.Errorf("resolved to %#v, want an error", got)
			}
		})
	}
}

// An error is the census's only signal that a call must not be attempted, and
// the census records which parameter blocked it. The message therefore has to
// name the parameter, and the account-identity refusal has to name account
// discovery, because "unresolvable" alone leaves a reader of the census with
// nothing to act on.
func TestParamErrorNamesTheParameter(t *testing.T) {
	_, err := Param("period", ParamSpec{Type: "string"})
	if err == nil {
		t.Fatal("err = nil, want an error")
	}
	if !strings.Contains(err.Error(), "period") {
		t.Errorf("error %q does not name the parameter", err)
	}

	_, err = Param("account_id", ParamSpec{Type: "string"})
	if err == nil {
		t.Fatal("err = nil, want an error")
	}
	if !strings.Contains(err.Error(), "discover") {
		t.Errorf("error %q does not name account discovery", err)
	}
}

// A refused parameter returns no value. A caller that ignores the error and
// sends the result anyway would be sending a nil, which is a request the
// endpoint cannot distinguish from an absent parameter.
func TestParamErrorReturnsNoValue(t *testing.T) {
	got, err := Param("account_id", ParamSpec{Type: "string"})
	if err == nil {
		t.Fatal("err = nil, want an error")
	}
	if got != nil {
		t.Errorf("got %#v alongside an error, want nil", got)
	}
}
