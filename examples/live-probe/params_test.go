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
// corpus spells it nine ways, so all nine are pinned: a name the class covers
// that resolves to an error is reported as an unbuildable request, and one that
// resolves to a bare string where the page declares a list is a request the
// endpoint answers with a 400 the probe would misread.
//
// The names below are the ones the corpus spells, so a page that starts spelling
// the class a tenth way is a case to add here rather than a surprise in a census.
func TestParamSymbolNames(t *testing.T) {
	cases := []struct {
		name  string
		param string
		spec  ParamSpec
		want  any
	}{
		{"symbol", "symbol", ParamSpec{Type: "string"}, DefaultSymbol},
		{"series_symbol", "series_symbol", ParamSpec{Type: "string"}, DefaultSymbol},
		{"event_symbol", "event_symbol", ParamSpec{Type: "string"}, DefaultSymbol},
		{"root_symbol", "root_symbol", ParamSpec{Type: "string"}, DefaultSymbol},
		{"underlying_symbol", "underlying_symbol", ParamSpec{Type: "string"}, DefaultSymbol},
		// The plural is the one place the resolution returns a container. The
		// page declares the arity in the name, since the spec carries no field
		// for it, so `symbols` is a list of one and `symbol` is a bare string.
		{"symbols", "symbols", ParamSpec{Type: "string"}, []string{DefaultSymbol}},
		{"symbols declared as an array", "symbols", ParamSpec{Type: "array"}, []string{DefaultSymbol}},
		{"option_symbols", "option_symbols", ParamSpec{Type: "string"}, []string{DefaultSymbol}},
		{"category_symbols", "category_symbols", ParamSpec{Type: "string"}, []string{DefaultSymbol}},
		// `instruments` is a plural symbol name as a query parameter and a list
		// of order bodies as a request body. Param only ever sees the former, so
		// the body case stays with the caller's body synthesis: resolving it
		// here cannot make a body buildable.
		{"instruments", "instruments", ParamSpec{Type: "string"}, []string{DefaultSymbol}},
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
// to call, and the brief pins one name of the two. All five names the corpus
// spells are pinned here because the property is the uniqueness rather than the
// name: a probe that sends the same key to two endpoints gets the second request
// rejected as a duplicate, which is a false negative about the endpoint rather
// than about the key.
func TestParamCorrelationKeysAreUnique(t *testing.T) {
	for _, name := range []string{"client_request_id", "client_order_id", "reqid", "event_id", "milestone_id"} {
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

// A correlation key is echoed into logs and sometimes into a response, so its
// value has to be recognisable as this probe's own. The property pinned here is
// the prefix: the corpus publishes a UUID for `reqid` and `client_request_id`,
// so an unprefixed 32-hex string is exactly the shape a reader could not tell
// apart from an identifier the server issued.
func TestParamCorrelationKeysAreObviouslySynthetic(t *testing.T) {
	for _, name := range []string{"client_request_id", "client_order_id", "reqid", "event_id", "milestone_id"} {
		t.Run(name, func(t *testing.T) {
			got, err := Param(name, ParamSpec{Type: "string"})
			if err != nil {
				t.Fatalf("err = %v", err)
			}
			key, ok := got.(string)
			if !ok {
				t.Fatalf("got %#v (%T), want a string", got, got)
			}
			if !strings.HasPrefix(key, "live-probe-") {
				t.Errorf("key %q does not say that it is the probe's own", key)
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

// probePlaceholderToken stands in for a session token in the two cases that pin
// the session-token class. It is not a token and is not credential-shaped, and
// nothing sends it anywhere: the class is about which string the resolver hands
// back, and that is observable without a credential. No case in this file
// supplies a real token and none can, because the value enters the process
// through SetSessionToken at run time.
const probePlaceholderToken = "probe-placeholder"

// The access token is the one value the probe must never invent: the client
// already holds it from the token exchange, so the resolver reads the recorded
// one and refuses the parameter when there is none. Both halves are pinned,
// because the failure mode being guarded against is a fabricated bearer token
// leaving this process, and a resolver that quietly filled the gap would be
// indistinguishable from a correct one in a census.
func TestParamSessionToken(t *testing.T) {
	t.Cleanup(func() { SetSessionToken("") })

	got, err := Param("access_token", ParamSpec{Type: "string"})
	if err == nil {
		t.Fatalf("resolved to %#v with no token recorded; the probe must not invent one", got)
	}
	if got != nil {
		t.Errorf("got %#v alongside an error, want nil", got)
	}
	if !strings.Contains(err.Error(), "access_token") {
		t.Errorf("error %q does not name the parameter", err)
	}

	SetSessionToken(probePlaceholderToken)
	got, err = Param("access_token", ParamSpec{Type: "string"})
	if err != nil {
		t.Fatalf("err = %v with a token recorded", err)
	}
	if got != probePlaceholderToken {
		t.Errorf("got %#v, want the recorded token back unchanged", got)
	}
	if SessionToken() != probePlaceholderToken {
		t.Errorf("SessionToken() = %q, want the recorded token", SessionToken())
	}
}

// Setting an empty token is how a caller returns the process to its start state,
// and it is a clear rather than a value to send: a recorded empty token must not
// resolve the parameter, because that would put an empty bearer on the wire.
func TestParamEmptySessionTokenIsNotAToken(t *testing.T) {
	t.Cleanup(func() { SetSessionToken("") })
	SetSessionToken("")
	if got, err := Param("access_token", ParamSpec{Type: "string"}); err == nil {
		t.Errorf("resolved to %#v; an empty token is an absent token", got)
	}
	if SessionToken() != "" {
		t.Errorf("SessionToken() = %q, want empty", SessionToken())
	}
}

// The optional cursor is the one parameter whose correct request is the one that
// leaves it out, so the marker has to be distinguishable from every value Param
// could otherwise return. The empty string is the case that matters: a caller
// that serialises a resolved "" as a present key has sent a different request
// from one that omits the key, and the server is free to answer them differently.
func TestParamOmitsTheOptionalCursor(t *testing.T) {
	for _, tc := range []struct {
		name string
		spec ParamSpec
	}{
		{"string", ParamSpec{Type: "string"}},
		{"undeclared type", ParamSpec{}},
		{"array", ParamSpec{Type: "array"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Param("pagination_key", tc.spec)
			if err != nil {
				t.Fatalf("omitting an optional parameter must not be an error: %v", err)
			}
			if !IsOmitted(got) {
				t.Errorf("got %#v (%T), want the Omitted marker", got, got)
			}
			if got == "" {
				t.Error("omitted resolved to the empty string, which is a different request")
			}
		})
	}
}

// The marker is recognised as itself and as nothing else. Without this, a caller
// testing for omission with a type assertion would also match a page's own
// empty-string example, which is a value to send.
func TestIsOmittedMatchesOnlyTheMarker(t *testing.T) {
	if !IsOmitted(Omitted) {
		t.Error("IsOmitted(Omitted) = false")
	}
	for _, v := range []any{nil, "", "cursor", 0, false, []string{DefaultSymbol}, DefaultSymbol} {
		if IsOmitted(v) {
			t.Errorf("IsOmitted(%#v) = true, want false", v)
		}
	}
}

// The SDK signs nine headers itself and they appear in the published document of
// an endpoint, so a caller can hand one to Param. Every one is refused, and the
// refusal is case-insensitive and separator-insensitive because an HTTP header
// name is: a refusal that a capital letter walks past is not a refusal. The
// documented example and enum cases in the loop are the ordering claim — the
// refusal sits above both, because a page's own example for x-app-key is still a
// credential that must not leave this process.
func TestParamRefusesSDKSignedHeaders(t *testing.T) {
	names := []string{
		"x-app-key", "x-app-secret", "x-timestamp", "x-access-token", "x-signature",
		"x-signature-algorithm", "x-signature-nonce", "x-signature-version", "x-version",
		"X-App-Key", "X-SIGNATURE", "x_app_key", "  x-signature  ",
	}
	specs := []struct {
		name string
		spec ParamSpec
	}{
		{"bare", ParamSpec{Type: "string"}},
		{"documented example", ParamSpec{Type: "string", Example: "documented", HasExample: true}},
		{"documented enum", ParamSpec{Type: "string", Enum: []any{"documented"}}},
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			for _, s := range specs {
				got, err := Param(name, s.spec)
				if err == nil {
					t.Errorf("with a %s: resolved to %#v; the SDK signs this header", s.name, got)
				}
				if got != nil {
					t.Errorf("with a %s: got %#v alongside an error, want nil", s.name, got)
				}
			}
		})
	}
}

// The refusal is narrow: a header-shaped name that the SDK does not sign is not
// refused, because refusing names it does not own would block a request the
// probe can make and the refusal would read as evidence about the endpoint.
func TestParamDoesNotRefuseUnownedHeaderNames(t *testing.T) {
	for _, name := range []string{"x-request-id", "x-trace-id", "accept", "user-agent"} {
		t.Run(name, func(t *testing.T) {
			if _, err := Param(name, ParamSpec{Type: "string"}); err == nil {
				t.Logf("%q resolves through the type default", name)
			}
			if isSDKSignedHeader(name) {
				t.Errorf("%q is not a header the SDK signs, so it must not be refused", name)
			}
		})
	}
}

// Every name the tier-B rules claim to cover has to resolve, and the names are
// written out here rather than read from the implementation, so that a rule
// dropped from params.go fails this case instead of quietly returning a blocked
// census row. The counts the corpus measures for each name are in the
// task-2 report; this case is the floor, not the census.
func TestParamCoversEveryDocumentedTierBName(t *testing.T) {
	t.Cleanup(func() { SetSessionToken("") })
	SetSessionToken(probePlaceholderToken)
	names := []string{
		// Rule 1, default symbol.
		"symbol", "symbols", "series_symbol", "event_symbol", "root_symbol",
		"underlying_symbol", "option_symbols", "instruments", "category_symbols",
		// Rule 2, caller-generated id.
		"reqid", "client_request_id", "client_order_id", "event_id", "milestone_id",
		// Rule 3, optional cursor omitted.
		"pagination_key",
		// Rule 4, session token.
		"access_token",
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			if got, err := Param(name, ParamSpec{Type: "string"}); err != nil {
				t.Errorf("Param(%q) = %v; a name the tier-B rules cover must resolve", name, err)
			} else if got == nil {
				t.Errorf("Param(%q) = nil with no error; a resolved value is never nil", name)
			}
		})
	}
}
