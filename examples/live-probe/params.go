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
	"strings"
	"sync"
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

// Omitted is the value Param returns for a documented parameter that the
// request must not carry. It is a value, not an error: the parameter is real and
// the probe knows what it is, and the correct request is the one that leaves it
// out.
//
// A caller must therefore check IsOmitted and drop the key from the query. It
// must not send a placeholder cursor and must not send an empty string: both are
// a different request from an absent key, and a server is free to answer them
// differently, which would be recorded as the endpoint's own answer.
//
// Omitted is distinguishable from every value Param can otherwise return, and in
// particular from a resolved empty string, because it is the only value of its
// own unexported type. A caller may compare against it directly (`v == Omitted`);
// IsOmitted is the form that needs no type knowledge at the call site.
var Omitted = omitted{}

// omitted is the type of Omitted. It is unexported and carries nothing, so the
// only way to hold a value of it is to compare against Omitted.
type omitted struct{}

// IsOmitted reports whether v is the Omitted marker, which is how a caller of
// Param tells "leave this parameter out of the request" apart from every other
// resolved value, including the empty string.
func IsOmitted(v any) bool {
	_, ok := v.(omitted)
	return ok
}

// session holds the access token the client already holds from the token
// exchange, for the one documented parameter that is a session value rather than
// a probe input. It is guarded because the census may resolve parameters from
// several goroutines while main records the token, and an unguarded read would
// be a data race the race detector is entitled to fail.
var session struct {
	mu    sync.RWMutex
	token string
}

// SetSessionToken records the access token the client obtained from the token
// exchange, so that a page documenting `access_token` as a request parameter can
// be answered with the real one. main calls it once, after EnsureToken succeeds
// and before any endpoint is walked.
//
// An empty string clears the recorded token, which is how a test returns the
// process to its start state. While no token is recorded, `access_token` is
// unresolvable: the probe reports the endpoint as blocked rather than sending a
// request carrying a token it made up.
func SetSessionToken(token string) {
	session.mu.Lock()
	defer session.mu.Unlock()
	session.token = token
}

// SessionToken reports the recorded access token, or the empty string where none
// is recorded. It is exported so a caller can assert what it set without
// reaching into the resolver.
func SessionToken() string {
	session.mu.RLock()
	defer session.mu.RUnlock()
	return session.token
}

// sdkSignedHeaders are the request headers the SDK signs and supplies itself.
// They appear in the published OpenAPI document of an endpoint, which is why a
// caller can hand one to Param, and they must never be resolved: x-app-key and
// x-app-secret are credentials, and the rest of the set is a signature over
// them. Task 3 filters them out of the documented parameter list before
// resolving anything, so this refusal is a backstop for a future caller rather
// than a path the census takes.
var sdkSignedHeaders = map[string]struct{}{
	"x-app-key":             {},
	"x-app-secret":          {},
	"x-timestamp":           {},
	"x-access-token":        {},
	"x-signature":           {},
	"x-signature-algorithm": {},
	"x-signature-nonce":     {},
	"x-signature-version":   {},
	"x-version":             {},
}

// isSDKSignedHeader reports whether name is one of the headers the SDK owns. The
// comparison is case-insensitive and treats `_` as `-`, because an HTTP header
// name is case-insensitive and the corpus spells these names both ways; a refusal
// that could be walked past with a capital letter would not be a refusal.
func isSDKSignedHeader(name string) bool {
	normalized := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), "_", "-"))
	_, ok := sdkSignedHeaders[normalized]
	return ok
}

// callerGeneratedIDs are the documented request parameters the caller mints
// itself, which the server accepts any unique value in. They are correlation
// keys, not server-issued identifiers, so a synthetic value is the correct
// answer rather than a guess.
var callerGeneratedIDs = map[string]struct{}{
	"reqid":             {},
	"client_request_id": {},
	"client_order_id":   {},
	"event_id":          {},
	"milestone_id":      {},
}

// paginationKey is the optional cursor a caller echoes back from a previous
// page. On a first request there is no previous page, so the correct request is
// the one that omits it.
const paginationKey = "pagination_key"

// sessionTokenParam is the documented request parameter that carries the access
// token the client already holds. The SDK sends the same value in a header, so
// this is the alternate spelling an endpoint documents rather than a credential
// the probe has to produce.
const sessionTokenParam = "access_token"

// accountParam is the caller's identity, which the server issues. It is named
// here because the refusal below is the one class whose reason is discovery
// rather than synthesis.
const accountParam = "account_id"

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
// query, the Omitted marker for a parameter the request must not carry, and a
// non-nil error when no class of parameter can fill it.
//
// The classes are tried in a fixed order, and the order is the contract:
//
//  1. a header the SDK signs itself, which is refused;
//  2. the page's own example, which is the only value the page says is valid;
//  3. the first member of a declared closed set, which the page accepts;
//  4. account_id, which is refused;
//  5. a symbol name, resolved to DefaultSymbol or a one-element list;
//  6. a caller-generated id, resolved to a value unique to this call;
//  7. an optional pagination cursor, resolved to Omitted;
//  8. access_token, resolved from the recorded session token;
//  9. a declared number, inside the published range;
//
// 10. a declared boolean, resolved to true;
// 11. otherwise an error.
//
// Steps 2 and 3 come first because the page is authoritative about a parameter
// it documents, and a synthesised value is a guess about one it does not.
//
// Step 1 sits above all of them, which is the one place this order departs from
// "the page is authoritative": a page that publishes an example for `x-app-key`
// is still a page whose header the SDK signs, and the value that would leave
// this process is a credential. Refusing is the only safe reading, and it costs
// nothing in the census, because Task 3 filters these nine names out before
// resolving anything.
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
	// The SDK signs these nine headers itself. Refusing them is deliberate: a
	// resolved x-app-key or x-app-secret is a credential this process invented,
	// and a resolved x-signature is a signature over a credential this process
	// invented. Either would leave the process and be rejected, and the census
	// would record the rejection as the endpoint's answer. Unresolvable is the
	// honest answer, and it is the only one that cannot be mistaken for a
	// working credential. Task 3 filters these names out of the documented
	// parameter list, so this is a backstop, not the census path.
	if isSDKSignedHeader(name) {
		return nil, fmt.Errorf("live-probe: %q is a header the SDK signs and supplies: the request must not carry one", name)
	}
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
	if name == accountParam {
		return nil, fmt.Errorf("live-probe: %q is issued by the server: discover an account and thread its id in", name)
	}
	// The instrument. The plural is resolved to a list of one because the name
	// is the only place the page's arity is visible: the spec carries no field
	// for it, and a plural name over a singular value is a request the endpoint
	// answers with a type error. A page that declares a plural name as an array
	// is the same case, so the declared type is not consulted here.
	switch name {
	case "symbol", "series_symbol", "event_symbol", "root_symbol", "underlying_symbol":
		return DefaultSymbol, nil
	case "symbols", "option_symbols", "category_symbols", "instruments":
		return []string{DefaultSymbol}, nil
	}
	// A correlation key is the one class the server accepts any value in, so it
	// is the one class where a synthesised value is always right. It has to be
	// unique: a repeated key is rejected as a duplicate, which is an error
	// attributed to the endpoint rather than to the key.
	if _, ok := callerGeneratedIDs[name]; ok && isStringType(spec) {
		return newCorrelationID(), nil
	}
	// The cursor for the next page, on a request for the first page. Omitting it
	// is a different request from sending a placeholder, and only the omission
	// asks the endpoint for what the probe meant.
	if name == paginationKey {
		return Omitted, nil
	}
	// The access token is a value the client already holds, so the probe reads it
	// from the session rather than inventing one. With no token recorded the
	// parameter is unresolvable: a fabricated bearer token would be sent to an
	// endpoint that rejects it, and the rejection would be recorded as the
	// endpoint's answer. This is the one class where the correct answer is
	// available to the caller and absent here, and an error says so.
	if name == sessionTokenParam && isStringType(spec) {
		token := SessionToken()
		if token == "" {
			return nil, fmt.Errorf("live-probe: %q is the session token: no token is recorded, so the probe must not invent one", name)
		}
		return token, nil
	}
	switch spec.Type {
	case "integer", "number":
		return clampCount(spec.Minimum, spec.Maximum), nil
	case "boolean":
		return true, nil
	}
	return nil, fmt.Errorf("live-probe: no value is synthesable for parameter %q of declared type %q", name, spec.Type)
}

// isStringType reports whether the page declared the parameter as a string, or
// declared no type at all. An absent type is not a contradiction of a name the
// corpus spells as a string, and refusing a correlation key because a page
// omitted its type would block a call the probe can make; any other declared
// type is a contradiction, and a string sent for it is a request the endpoint
// answers with a type error.
func isStringType(spec ParamSpec) bool {
	return spec.Type == "" || spec.Type == "string"
}

// newCorrelationID returns a value no other call in this process returns, for a
// caller-generated correlation key.
//
// The counter is what makes two calls in one run differ, and it is what the
// probe would rely on if two calls were issued in the same nanosecond. The
// timestamp is what separates two runs, so a key from a previous run is not
// resubmitted and rejected as a duplicate by a server that remembers them.
//
// The `live-probe-` prefix is what makes the value obviously synthetic. These
// keys are echoed back in logs and in some responses, and a value that could be
// mistaken for an identifier the server issued would be evidence about this
// process that it is not. A counter and a clock are used rather than randomness
// because a value a test can predict is a value a test can assert about: a random
// generator would make "two calls differ" the only property available, and would
// leave the format untested.
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
