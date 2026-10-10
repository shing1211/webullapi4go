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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/shing1211/webullapi4go/client"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/trade"
)

// NamedParam is one documented request parameter: the name the page uses for it
// and what the page documents about it.
//
// The split is this file's, not ParamSpec's. Param takes the name as its first
// argument, so a container that has to carry a name alongside a spec needs a
// type of its own, and adding a Name field to ParamSpec would have changed a
// shipped, reviewed interface to serve a caller that did not exist yet.
type NamedParam struct {
	// Name is the parameter name exactly as the page spells it.
	Name string
	// Spec is what the page documents about the parameter.
	Spec ParamSpec
}

// Endpoint is one documented endpoint to ask: a row per conformance fixture.
//
// RequiredParams and BodyParams carry only what the request cannot be built
// without, so a reader can tell an endpoint the probe can construct from one it
// cannot. An optional parameter is absent by design: the minimal conforming
// request omits it, and an omitted key is a different request from an empty one.
type Endpoint struct {
	// Symbol is the SDK method the manifest maps this endpoint to. It is
	// documentation for the reader; the census sends the documented path, not
	// the SDK's own, because the question is whether the endpoint answers.
	Symbol string
	// Fixture is the manifest fixture file name, "<area>/<METHOD>-<slug>.json".
	// The area is its first segment, which is where the areas diverge.
	Fixture string
	// Method is the documented HTTP method, upper case. The cached documents
	// spell it lower case and gen_fixtures.py upper-cases it before writing the
	// manifest, so the census does the same.
	Method string
	// Path is the documented path, with any {template} segment unsubstituted.
	// The substituted path is never recorded on an Outcome: it carries the
	// account id, and a live value must not reach a committed file.
	Path string
	// RequiredParams are the required query and header parameters, in page
	// order, with the headers the SDK signs already removed.
	RequiredParams []NamedParam
	// PathParams are the {template} segments Path carries, in the order they
	// appear. No page in the current cache has one; the field exists because a
	// literal brace in a URL is a different URL from the documented one.
	PathParams []NamedParam
	// BodyParams are the required request-body properties, in page order.
	BodyParams []NamedParam
	// HasBody reports that the page documents a request body, so the request
	// must carry one. 50 of the 193 documented endpoints do, and calling them
	// with an empty body records the probe's failure as the server's.
	HasBody bool
	// DroppedSDKSignedHeaders counts the documented header parameters removed
	// when this endpoint was read because the SDK signs them.
	DroppedSDKSignedHeaders int
	// Unreadable is why no schema could be read for this endpoint, and is
	// non-empty only when the cached page is absent or carries no JSON block.
	// It becomes the Outcome's Blocked reason, so a page the probe could not
	// read is never confused with an endpoint the probe could not call.
	Unreadable string
}

// Outcome is what happened to one endpoint.
//
// It records whether the sandbox answered and how, never what it said. Err
// carries the SDK's error Code and not the server's message, because
// [errs.FromHTTPStatus] folds a response body into that message and a body can
// carry an account number. The three numbers a reader must never add together
// are Reachable, Blocked and Skipped, and [Outcome] keeps them apart: a status
// means the endpoint answered, Blocked means the probe could not build a
// request, and Skipped means the probe built one and declined to send it.
type Outcome struct {
	// Symbol is the SDK method the manifest maps the endpoint to.
	Symbol string `json:"symbol"`
	// Fixture is the manifest fixture file name, which carries the area.
	Fixture string `json:"fixture"`
	// Area is the first segment of Fixture, e.g. "trading" or "broker-fd-us".
	Area string `json:"area"`
	// Method is the HTTP method the request was sent with.
	Method string `json:"method"`
	// Path is the documented path, unsubstituted. The path actually sent
	// carries a discovered account id and is deliberately not recorded.
	Path string `json:"path"`
	// Status is the HTTP status the endpoint answered with, or 0 when no
	// response was received. A 404 is a status; a Blocked is not.
	Status int `json:"status"`
	// Err is the SDK error Code, or "" on a 2xx and on a blocked endpoint. It
	// is a classification and never the server's message text.
	Err string `json:"err,omitempty"`
	// Blocked is why the probe could not ask, or "" when it asked. A blocked
	// endpoint has Status 0 by construction: the loop never issues a request it
	// already knows it cannot build, because a fabricated value would be
	// rejected by the server and the rejection recorded as the endpoint's own
	// answer.
	Blocked string `json:"blocked,omitempty"`
	// Skipped is why the probe declined to call an endpoint it could have
	// called, or "" when no gate stopped it. A skipped endpoint has Status 0 and
	// Blocked "" by construction, and it is the only Outcome carrying a server's
	// absence: reporting "we did not call it" as anything else - a 404, a
	// transport failure, a block - would state a server behaviour the probe never
	// observed. The three classes are counted separately and never summed into
	// one.
	Skipped string `json:"skipped,omitempty"`
	// Mutating reports that the endpoint is in the mutating class (see
	// [IsMutating]), whether or not the gate stopped the call. It is recorded on
	// every run, including a dry run that sends nothing, so the size of the class
	// is knowable without a credential and a reader can see which rows a live
	// run would not have called.
	Mutating bool `json:"mutating,omitempty"`
	// QueryKeys names the query parameters the request carried. The values are
	// not recorded, for the same reason Path is not substituted.
	QueryKeys []string `json:"queryKeys,omitempty"`
	// BodyBytes is the length of the body sent, and BodyKeys its property
	// names. Neither the body nor any value in it is recorded.
	BodyBytes int      `json:"bodyBytes,omitempty"`
	BodyKeys  []string `json:"bodyKeys,omitempty"`
	// SynthesisedEmptyArrays names the required body arrays the probe filled
	// with an empty array. It is reported because 16 endpoints rest on that
	// choice and a reader should be able to see which.
	SynthesisedEmptyArrays []string `json:"synthesisedEmptyArrays,omitempty"`
	// DroppedSDKSignedHeaders counts the documented header parameters removed
	// because the SDK signs them. A non-zero value is the normal case on 164 of
	// the 193 pages, and a row that still carries one of those names is a
	// filter bug rather than a finding.
	DroppedSDKSignedHeaders int `json:"droppedSdkSignedHeaders,omitempty"`
	// DryRun reports that the endpoint was built but not sent.
	DryRun bool `json:"dryRun,omitempty"`
}

// AreaSummary is one area's census. Reachable counts endpoints that answered,
// whatever they answered: a 404 from broker-fd-us is a live answer and a 200 is
// a live answer, and collapsing the two would lose the information the census
// exists to collect. Blocked counts endpoints the probe could not ask, Skipped
// counts the ones it chose not to ask, and Unanswered counts requests that
// produced no response at all. Those four partition Total.
type AreaSummary struct {
	Total      int `json:"total"`
	Reachable  int `json:"reachable"`
	Blocked    int `json:"blocked"`
	Skipped    int `json:"skipped"`
	Unanswered int `json:"unanswered"`
	// Mutating is how many of the area's endpoints are in the mutating class,
	// so Skipped can be read against the size of the class rather than against
	// the area total.
	Mutating int            `json:"mutating"`
	Statuses map[string]int `json:"statuses,omitempty"`
}

// add records one outcome in this area.
//
// The order of the first two cases is the safety property. A skipped endpoint is
// counted before a blocked one and neither reaches the status map, so "the probe
// chose not to send this" can never be tallied as an answer or as a build
// failure.
func (a *AreaSummary) add(o Outcome) {
	a.Total++
	if o.Mutating {
		a.Mutating++
	}
	switch {
	case gateClosed(o):
		a.Skipped++
	case o.Blocked != "":
		a.Blocked++
	case o.Status == 0:
		a.Unanswered++
	default:
		a.Reachable++
	}
	if o.Status != 0 {
		if a.Statuses == nil {
			a.Statuses = map[string]int{}
		}
		a.Statuses[strconv.Itoa(o.Status)]++
	}
}

// CensusSummary is the whole census, and the distinction it exists to keep is
// between a number of endpoints the sandbox answered and a number of endpoints
// the probe was able to ask. Those are different numbers and they diverge by
// area sharply: request-constructibility is 100% of market data and 0% of
// display-solution, while sandbox-reachability is dominated by 404, 403 and 401
// answers that have nothing to do with whether the probe can build a request.
type CensusSummary struct {
	Total      int `json:"total"`
	Reachable  int `json:"reachable"`
	Blocked    int `json:"blocked"`
	Skipped    int `json:"skipped"`
	Unanswered int `json:"unanswered"`
	// Mutating is the size of the mutating class over the whole corpus. It is
	// reported whether or not the gate stopped anything, because the number a
	// reader needs before authorising a live run is the number a live run would
	// not cover.
	Mutating int `json:"mutating"`
	// Area is keyed by the area name and holds one AreaSummary per area.
	Area map[string]AreaSummary `json:"area"`
	// Statuses counts every status the corpus answered with, over all areas.
	Statuses map[string]int `json:"statuses,omitempty"`
	// BlockedReasons counts the distinct Blocked strings, so a reader can see
	// which class of probe limitation produced the blocked rows.
	BlockedReasons map[string]int `json:"blockedReasons,omitempty"`
	// Errs counts the SDK error codes, which classify a request that produced
	// no status.
	Errs map[string]int `json:"errs,omitempty"`
}

// Summarise aggregates outcomes into the counts a reader looks at first. The
// four totals partition Total: Reachable, Blocked, Skipped and Unanswered, and
// none of them is the same measure as another.
func Summarise(outcomes []Outcome) CensusSummary {
	summary := CensusSummary{
		Total:          len(outcomes),
		Area:           map[string]AreaSummary{},
		Statuses:       map[string]int{},
		BlockedReasons: map[string]int{},
		Errs:           map[string]int{},
	}
	for _, o := range outcomes {
		area := summary.Area[o.Area]
		area.add(o)
		summary.Area[o.Area] = area

		if o.Status != 0 {
			summary.Statuses[strconv.Itoa(o.Status)]++
		}
		if o.Err != "" {
			summary.Errs[o.Err]++
		}
		if o.Mutating {
			summary.Mutating++
		}
		switch {
		case gateClosed(o):
			summary.Skipped++
		case o.Blocked != "":
			summary.Blocked++
			summary.BlockedReasons[o.Blocked]++
		case o.Status == 0:
			summary.Unanswered++
		default:
			summary.Reachable++
		}
	}
	return summary
}

// errCacheAbsent reports that the gitignored docgen cache holds no pages, so
// the census has no request schemas to read. It is a sentinel rather than a
// plain error because the right response is to skip, not to fail: the cache is
// not committed, so a fresh checkout cannot judge the corpus either way. The
// conformance-fixtures target skips for the same reason.
var errCacheAbsent = errors.New("no docgen cache")

// account holds the account id phase 0 discovered, for the documented account_id
// parameters and body properties that need one. It is guarded for the same
// reason the session token is: a test sets it and the loop reads it, and an
// unguarded read is a race the detector is entitled to fail.
var account struct {
	mu sync.RWMutex
	id string
}

// SetAccountID records the account id every account-scoped request is addressed
// to. main calls it once, after phase 0 and before any endpoint is walked, and
// calls it with the empty string when phase 0 did not run or did not succeed.
//
// While no id is recorded, account_id is unresolvable and every endpoint that
// needs one is reported as blocked rather than called against an account that
// does not exist.
func SetAccountID(id string) {
	account.mu.Lock()
	defer account.mu.Unlock()
	account.id = id
}

// AccountID reports the recorded account id, or the empty string where none is
// recorded. It is exported so a caller can assert what it set.
func AccountID() string {
	account.mu.RLock()
	defer account.mu.RUnlock()
	return account.id
}

// blockedAccountReason is the Blocked text of every endpoint that needs an
// account id and has none. It names the cause rather than a generic resolution
// failure, because "blocked" alone leaves a reader with nothing to act on: the
// fix is a successful phase 0, not a wider resolver.
const blockedAccountReason = "account_id: account discovery did not yield an id"

// accountIDNames are the spellings of the caller's identity the corpus uses.
// They are compared after normalisation, so accountId and account_id are the
// same parameter.
var accountIDNames = map[string]struct{}{
	"account_id": {},
	"accountid":  {},
}

// isAccountParam reports whether name is the caller's identity.
func isAccountParam(name string) bool {
	_, ok := accountIDNames[normaliseParamName(name)]
	return ok
}

// normaliseParamName reduces a documented parameter name to the form the
// resolver's name classes are written in: lower case, non-alphanumerics
// collapsed to a single underscore. A page that spells the same parameter
// accountId and account_id documents one parameter, and a refusal that could be
// walked past with a capital letter would not be a refusal.
//
// This is NOT conformance.fold, which lives in the conformance package and does
// the opposite thing with the same intent: it DELETES every underscore and
// upper-cases the rest, so "instrumentId" and "instrument_id" both become
// "INSTRUMENTID" and compare equal. The two normalisers answer different
// questions - "is this one parameter name?" against a page, and "is this one wire
// member?" against a body encoding/json already declined to match - so they have
// different rules and neither is a substitute for the other. See conformance.fold
// for the side that says which.
func normaliseParamName(name string) string {
	var b strings.Builder
	b.Grow(len(name))
	underscore := false
	for _, r := range strings.ToLower(strings.TrimSpace(name)) {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
			underscore = false
		case !underscore && b.Len() > 0:
			b.WriteByte('_')
			underscore = true
		}
	}
	return strings.TrimRight(b.String(), "_")
}

// resolveValue decides the value to send for one documented parameter.
//
// It is [Param] with one class in front of it. account_id is the caller's
// identity and the pages publish an example for it that is an account nobody
// holds, so the page's own example is exactly the fabrication this probe exists
// to rule out: sending it would point 25 account-scoped endpoints at an account
// that does not exist and record the rejection as each endpoint's own answer.
// Param deliberately lets the example win for this name, and its own comment
// names the fix as "a discovered account threaded in ahead of this call"
// (params_test.go:270). Threading it in ahead of the call is this function.
//
// Every other class is Param's, unchanged and in Param's order.
func resolveValue(name string, spec ParamSpec) (any, error) {
	if isAccountParam(name) {
		id := AccountID()
		if id == "" {
			return nil, errors.New(blockedAccountReason)
		}
		return id, nil
	}
	return Param(name, spec)
}

// errNoAccountReported is phase 0 finding no account to address. It is a sentinel
// rather than a formatted error because the report renders this case by identity:
// a caller has to be able to say "the list came back empty" without quoting the
// error, and the error's own text is prose this program wrote, so quoting it would
// be safe but a sentinel keeps the one line in the artefact from depending on a
// message that can be reworded.
var errNoAccountReported = errors.New("live-probe: account discovery: the account list is empty, " +
	"so no account id can be threaded into an account-scoped request")

// DiscoverAccountID returns the first account id the trading account list
// reports.
//
// This is phase 0. It runs once, before any endpoint is walked, because an
// account-scoped request is only meaningful against a real account: a fabricated
// id produces a 404 or a 400 that says nothing about whether the endpoint is
// reachable, which is the confusion a reachability census exists to avoid.
//
// An error is returned rather than swallowed so main can decide what to do with
// it, and an empty list is an error rather than an empty string, because an
// empty account id threaded into a request is the same fabrication as a
// fabricated one.
func DiscoverAccountID(ctx context.Context, cl *client.Client) (string, error) {
	accounts, err := trade.New(cl).ListAccounts(ctx)
	if err != nil {
		return "", fmt.Errorf("live-probe: account discovery: %w", err)
	}
	if len(accounts) == 0 {
		return "", errNoAccountReported
	}
	return accounts[0].AccountID, nil
}

// preparedRequest is the request [Census] would send for one endpoint. It is
// separated from the loop so a test can assert the query string and the body
// without an HTTP server, and so a dry run can build every request without a
// credential.
type preparedRequest struct {
	// Path is the request path with every {template} substituted. It carries
	// the account id and is never recorded on an Outcome.
	Path string
	// Query is the query string as key/value pairs, in no particular order;
	// url.Values.Encode is what fixes the order on the wire.
	Query url.Values
	// Body is the JSON request body, or nil where the request carries none.
	Body []byte
	// Blocked is why the request cannot be built, and is non-empty exactly when
	// the endpoint must not be called.
	Blocked string
	// SynthesisedEmptyArrays names the required body arrays filled with [].
	SynthesisedEmptyArrays []string
	// DroppedSDKSignedHeaders counts the documented headers removed because the
	// SDK signs them.
	DroppedSDKSignedHeaders int
	// BodyKeys names the body properties the request carried.
	BodyKeys []string
}

// prepare builds the request for one endpoint, or the reason it cannot be
// built.
//
// A blocked request is never returned with a partial path, query or body: the
// caller checks Blocked first and must not send what is left, because half a
// request is a request the endpoint answers for a reason that has nothing to do
// with the endpoint.
//
// Three classes of documented input are deliberately not built:
//
//   - the nine headers the SDK signs, which are removed when the Endpoint is
//     read rather than here, so that the count of them removed is recorded and
//     a row still carrying one is visible as a filter bug;
//   - a required string with no published example and no closed set, for which
//     "" is a value the page never sanctioned;
//   - a request body the page does not document as JSON.
func prepare(ep Endpoint) preparedRequest {
	if ep.Unreadable != "" {
		return preparedRequest{Blocked: ep.Unreadable}
	}
	if ep.Method == "" {
		return preparedRequest{Blocked: "the page documents no HTTP method"}
	}
	if ep.Path == "" {
		return preparedRequest{Blocked: "the page documents no path"}
	}

	out := preparedRequest{Path: ep.Path, Query: url.Values{}}

	// Path templates first: a path that still holds a brace is a different URL
	// from the documented one, so it is either filled or blocked.
	for _, p := range ep.PathParams {
		value, err := resolveValue(p.Name, p.Spec)
		if err != nil {
			return preparedRequest{Blocked: fmt.Sprintf("path parameter %s: %v", p.Name, err)}
		}
		text, ok := scalarText(value)
		if !ok {
			return preparedRequest{Blocked: fmt.Sprintf(
				"path parameter %s: resolved to a value that is not a path segment", p.Name)}
		}
		out.Path = strings.ReplaceAll(out.Path, "{"+p.Name+"}", url.PathEscape(text))
	}
	if strings.Contains(out.Path, "{") {
		return preparedRequest{Blocked: "the documented path holds an unresolved template: " + out.Path}
	}

	for _, p := range ep.RequiredParams {
		value, err := resolveValue(p.Name, p.Spec)
		if err != nil {
			return preparedRequest{Blocked: fmt.Sprintf("parameter %s: %v", p.Name, err)}
		}
		// The cursor for the next page has no page to come from. Omitting the
		// key is the request that asks for what the probe meant; an empty value
		// is a different request the server may answer differently.
		if IsOmitted(value) {
			continue
		}
		switch typed := value.(type) {
		case []string:
			for _, item := range typed {
				out.Query.Add(p.Name, item)
			}
		case string:
			out.Query.Set(p.Name, typed)
		default:
			text, ok := scalarText(value)
			if !ok {
				return preparedRequest{Blocked: fmt.Sprintf(
					"parameter %s: resolved to a value that cannot be sent as a query value", p.Name)}
			}
			out.Query.Set(p.Name, text)
		}
	}

	if !ep.HasBody {
		return out
	}
	body := map[string]any{}
	for _, p := range ep.BodyParams {
		if p.Spec.Type == "array" && !p.Spec.HasExample && len(p.Spec.Enum) == 0 {
			// An empty array is the minimal instance a required array admits: no
			// page in the cache declares minItems, so [] satisfies the schema.
			// It is also the only value that cannot place an order, cancel one
			// or move money, which is what lets a called trading endpoint be
			// inert. The choice is recorded on the Outcome so a reader can see
			// which endpoints rest on it.
			//
			// This is the SECOND of two barriers, and it is not a substitute for
			// the first. [IsMutating] gates the 34 endpoints that can change
			// state; this rule makes the ones that are called harmless. They are
			// independent and each covers a case the other cannot:
			//
			//   - /trading/orders/place holds a required array, so it is both
			//     gated and inert; the gate is what protects it, since prepare
			//     does not consult the gate at all.
			//   - /broker/accounts/create, /trading/orders/cancel and
			//     /broker/funding/ach/relationships/create have scalar bodies, so
			//     there is no array here to empty and this rule says nothing about
			//     them. The gate is the only barrier.
			//   - a market-data-watchlist POST holds a required array and is not
			//     in the mutating class, so it is called, and this rule is the
			//     only barrier.
			//
			// Either mechanism alone leaves a class of endpoint that can be
			// changed, so both are load-bearing and both are tested.
			body[p.Name] = []any{}
			out.SynthesisedEmptyArrays = append(out.SynthesisedEmptyArrays, p.Name)
			continue
		}
		value, err := resolveValue(p.Name, p.Spec)
		if err != nil {
			return preparedRequest{Blocked: fmt.Sprintf("body property %s: %v", p.Name, err)}
		}
		body[p.Name] = value
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		return preparedRequest{Blocked: fmt.Sprintf("body: %v", err)}
	}
	out.Body = encoded
	out.BodyKeys = make([]string, 0, len(body))
	for name := range body {
		out.BodyKeys = append(out.BodyKeys, name)
	}
	sort.Strings(out.BodyKeys)
	return out
}

// scalarText renders a resolved value as the single string a query value or a
// path segment carries. A resolved value is a scalar for every class Param
// documents, so a non-scalar here means a value the resolver produced that this
// request cannot express, and saying so beats coercing it.
func scalarText(value any) (string, bool) {
	switch typed := value.(type) {
	case string:
		return typed, true
	case bool:
		return strconv.FormatBool(typed), true
	case int:
		return strconv.Itoa(typed), true
	case int64:
		return strconv.FormatInt(typed, 10), true
	case float64:
		return strconv.FormatFloat(typed, 'g', -1, 64), true
	case json.Number:
		return typed.String(), true
	}
	return "", false
}

// sortedKeys returns the keys of v in ascending order, so a recorded key list is
// stable between runs and diffs cleanly.
func sortedKeys(v url.Values) []string {
	if len(v) == 0 {
		return nil
	}
	keys := make([]string, 0, len(v))
	for name := range v {
		keys = append(keys, name)
	}
	sort.Strings(keys)
	return keys
}

// Census walks endpoints in order and records one [Outcome] each.
//
// An endpoint whose request cannot be built is recorded as blocked with no
// status and no request: that is the whole point of the distinction the census
// draws. A probe that filled an unresolvable parameter and sent the result would
// record the server's rejection of the fabrication as the endpoint's own answer,
// and a census of 193 rejections would look like a finding.
//
// The mutating gate is the second distinction, and it is applied before the
// request is built, so a refused endpoint is refused on the fact of what it is
// rather than on whether this run happened to be able to construct a request for
// it. A refused endpoint gets no request, no status, no error code and no
// Blocked reason: it gets Skipped, which is a fourth class the summary counts on
// its own. Collapsing "we did not call it" into "the server said no" is the one
// outcome this walk must never produce, because it would report a sandbox
// restriction as a server behaviour.
//
// The gate is on sending, not on building. [prepare] stays pure, so a dry run
// still measures whether a mutating endpoint's request is constructible, and that
// is the only place the constructibility of the 34 is measurable at all - a live
// run that does not open the gate learns nothing about them beyond their
// existence.
//
// The gate is resolved ONCE, before the loop, and from this client's own base URL
// rather than from the environment, so a process that changed its own
// environment or reconfigured its client mid-walk could not end up with a
// half-open gate - and so the two conditions are judged against the host these
// requests will actually be sent to. [ResolveMutationGate] takes [hostOf] of
// this very client.
//
// Endpoints are walked sequentially. The rate limits are per App Key and the
// token endpoint allows ten requests per thirty seconds, so a concurrent walk
// would turn a census into a self-inflicted 429, and a 429 is an answer about the
// probe rather than about the endpoint.
//
// A non-nil error is returned only for a failure that stops the walk. An error
// from one endpoint is recorded on its Outcome and the walk continues, because
// one unreachable endpoint says nothing about the other 192.
func Census(ctx context.Context, cl *client.Client, endpoints []Endpoint) ([]Outcome, error) {
	if cl == nil {
		return nil, errors.New("live-probe: Census needs a client")
	}
	gate := ResolveMutationGate(hostOf(cl))
	outcomes := make([]Outcome, 0, len(endpoints))
	for _, ep := range endpoints {
		outcome := Outcome{
			Symbol:   ep.Symbol,
			Fixture:  ep.Fixture,
			Area:     areaOf(ep.Fixture),
			Method:   ep.Method,
			Path:     ep.Path,
			Mutating: IsMutating(ep),
		}
		if outcome.Mutating {
			if reason := gate.SkipReason(); reason != "" {
				outcome.Skipped = reason
				outcomes = append(outcomes, outcome)
				continue
			}
		}
		request := prepare(ep)
		if request.Blocked != "" {
			outcome.Blocked = request.Blocked
			outcomes = append(outcomes, outcome)
			continue
		}
		outcome.QueryKeys = sortedKeys(request.Query)
		outcome.BodyBytes = len(request.Body)
		outcome.BodyKeys = request.BodyKeys
		outcome.SynthesisedEmptyArrays = request.SynthesisedEmptyArrays
		outcome.DroppedSDKSignedHeaders = ep.DroppedSDKSignedHeaders
		outcome.Status = send(ctx, cl, ep.Method, request, &outcome)
		outcomes = append(outcomes, outcome)
	}
	return outcomes, nil
}

// send issues one census request and records its status and error code on
// outcome.
//
// It uses [client.Client.DoStream] rather than [client.Client.Do] for two
// reasons. Do returns nil on any 2xx and does not report the status, and a
// census whose success row cannot say 200 from 204 cannot answer its own
// question. And DoStream never retries, which is what a reachability census
// wants: one request, one answer, and a retried endpoint reports the last
// attempt rather than the first.
//
// The response body is drained and discarded. A census records whether an
// endpoint answered; capturing what it said is a later task, and a body that
// reached a committed file would be a value the server sent.
func send(ctx context.Context, cl *client.Client, method string, request preparedRequest, outcome *Outcome) int {
	target := request.Path
	if encoded := request.Query.Encode(); encoded != "" {
		target += "?" + encoded
	}
	var body any
	if len(request.Body) > 0 {
		body = json.RawMessage(request.Body)
	}

	resp, err := cl.DoStream(ctx, method, target, body)
	if resp != nil {
		outcome.Status = resp.StatusCode
		if resp.Body != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}
	}
	if err != nil {
		outcome.Err = classify(err)
		// A non-2xx answer arrives as a typed error rather than a response, so
		// the status the endpoint answered with is on the error. A 404 the census
		// cannot record is a 404 the reader would count as an unanswered request,
		// which is the opposite of what happened.
		if outcome.Status == 0 {
			var typed *errs.Error
			if errors.As(err, &typed) {
				outcome.Status = typed.Status
			}
		}
	}
	return outcome.Status
}

// classify reduces an SDK error to its Code.
//
// The message is dropped on purpose. [errs.FromHTTPStatus] folds the response
// body into the message of every non-2xx error, and a body can carry an account
// number, a position or a balance, so a message that reaches census.json is a
// live value in a file a reader may commit. The Code is a fixed vocabulary and
// the status is recorded beside it, which is everything a reachability row needs.
func classify(err error) string {
	var typed *errs.Error
	if errors.As(err, &typed) {
		return string(typed.Code)
	}
	return "unclassified"
}

// undeterminedArrays names the required body properties the probe would fill
// with an empty array, for ep.
//
// It is a pure reading of the document and resolves nothing, which is why it is
// separate from [prepare]. In a run with no discovered account every
// account-scoped endpoint is blocked before its body is built, so a count taken
// from built bodies would report the 16 endpoints that rest on the empty-array
// choice as far fewer, and the choice would look smaller than it is.
func undeterminedArrays(ep Endpoint) []string {
	var names []string
	for _, p := range ep.BodyParams {
		if p.Spec.Type == "array" && !p.Spec.HasExample && len(p.Spec.Enum) == 0 {
			names = append(names, p.Name)
		}
	}
	return names
}

// areaOf returns the area a fixture belongs to, which is its first path
// segment: "trading/POST-trading-orders-place.json" is in "trading". The area
// comes from the fixture name rather than a second field so that the two can
// never disagree, and a fixture with no separator is its own area rather than an
// empty one, because an empty key would merge unrelated rows in the summary.
func areaOf(fixture string) string {
	if index := strings.Index(fixture, "/"); index > 0 {
		return fixture[:index]
	}
	if fixture == "" {
		return "unknown"
	}
	return fixture
}

// --------------------------------------------------------------------------
// Reading the manifest and the docgen cache
// --------------------------------------------------------------------------
//
// The census reads the request schemas from the gitignored docgen cache, because
// the committed fixtures are response-only: tools/conformance/gen_fixtures.py
// never reads `parameters` or `requestBody`, so nothing committed in this
// repository states what an endpoint takes as input.
//
// The cache is parsed here in Go rather than by shelling out to the Python
// tools. Three reasons, in order of weight:
//
//  1. A cached page's JSON block is a single-endpoint document with top-level
//     `path`, `method`, `parameters` and `responses`, not an OpenAPI `paths`
//     map, and gen_fixtures.py is the only reader in the repository that knows
//     how to find it. Reimplementing that reader in Go is the same code twice;
//     reusing it from Go is not possible without a subprocess or a rewrite of
//     the tool, and a rewrite would put the committed fixtures behind a change
//     the conformance gate depends on.
//  2. The census is a Go program whose output the operator inspects, and a
//     subprocess would make the request-building path untestable in-process.
//  3. The selection rule is small and is pinned by a test rather than trusted:
//     [selectBlock] reproduces gen_fixtures.py's select_block, including the
//     "first JSON block is a guide example, not the endpoint" case that a
//     hand-rolled reader gets wrong and that yields zero endpoints across the
//     whole corpus with no error at all.
//
// The drift this exposes is bounded and observable: a page whose JSON block
// stops parsing becomes an endpoint whose Blocked reason says so, and the count
// of such endpoints is in the summary rather than hidden.

// cacheUnsafe matches every run of characters a cache filename may not contain.
// It is the Go spelling of the derivation in _common.fetch, which
// gen_fixtures.py also copies rather than calls.
var cacheUnsafe = regexp.MustCompile(`[^A-Za-z0-9]+`)

// jsonFence matches one fenced ```json block and captures its body. It is the
// Go spelling of the pattern in _common.extract_json and
// gen_fixtures._json_blocks, so a block this does not see is a block neither
// does.
var jsonFence = regexp.MustCompile("(?s)```json\\s*\n(.*?)\n```")

// pathTemplate matches one {template} segment of a documented path.
var pathTemplate = regexp.MustCompile(`\{([^{}]+)\}`)

// cacheFilename returns the cache file a URL's page is stored in.
//
// It is a copy of the derivation in _common.fetch and gen_fixtures.cache_filename
// rather than a call to either: fetch creates the cache directory and re-fetches
// on a miss, and the census may do neither. A page whose name no longer resolves
// this way becomes an endpoint reported as not cached, so the duplication cannot
// drift without the census saying so.
func cacheFilename(pageURL string) string {
	safe := cacheUnsafe.ReplaceAllString(pageURL, "_")
	if len(safe) > 120 {
		safe = safe[len(safe)-120:]
	}
	return safe + ".md"
}

// readPage returns the cached markdown for a page URL, or the empty string when
// the page is not cached. The census never fetches: it is a read of a local
// snapshot, and a run that reached the network to fill a gap would be measuring
// the network rather than the endpoints.
func readPage(cacheDir, pageURL string) (string, error) {
	path := filepath.Join(cacheDir, cacheFilename(pageURL))
	info, err := os.Stat(path)
	if err != nil || info.Size() == 0 {
		return "", nil
	}
	// The pages carry emoji, and the Windows default codec cannot decode them.
	raw, err := os.ReadFile(path) //nolint:gosec // G304: the path is derived from the manifest's own URL, inside a caller-supplied cache directory.
	if err != nil {
		return "", fmt.Errorf("live-probe: reading cached page %s: %w", path, err)
	}
	return string(raw), nil
}

// decodeDocument decodes one JSON block. It is configured with UseNumber so a
// numeric example reaches the wire as the literal the page published rather than
// as a float64 that re-encodes differently, which for a price or a size is a
// different request from the documented one.
func decodeDocument(raw []byte) (map[string]any, error) {
	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.UseNumber()
	var doc map[string]any
	if err := dec.Decode(&doc); err != nil {
		return nil, err
	}
	return doc, nil
}

// selectBlock returns the index and the endpoint document in a cached page.
//
// It reproduces gen_fixtures.py's select_block, and the case it exists for is
// the first one: a page whose first ```json block is a guide example rather than
// the endpoint document. Taking that first block yields a document with no
// `path` for the whole corpus, so every endpoint reports no parameters and no
// body and the census claims a clean sweep of an empty question. The rule is
// therefore: the first block when it carries both `path` and `responses`,
// otherwise the first block that does, otherwise the first block at all, and -1
// when the page has no JSON block.
//
// The index is the position among the page's fenced blocks, counting blocks that
// do not parse, which is what gen_fixtures.py reports as jsonBlockIndex.
func selectBlock(md string) (int, map[string]any) {
	matches := jsonFence.FindAllStringSubmatch(md, -1)
	first, hadFirst := map[string]any{}, false
	for i, match := range matches {
		doc, err := decodeDocument([]byte(match[1]))
		if err != nil {
			continue
		}
		if i == 0 {
			first, hadFirst = doc, true
		}
		if _, ok := doc["path"]; ok {
			if _, ok := doc["responses"]; ok {
				return i, doc
			}
		}
	}
	if hadFirst {
		return 0, first
	}
	return -1, nil
}

// paramSpecOf reads a [ParamSpec] out of one documented parameter object.
//
// The example is read off the parameter object first and off parameter.schema
// only as a fallback, and the order is the contract rather than a preference:
// measured across the cache, 433 parameter objects carry `example` and zero
// carry `schema.example`. A reader that looks only at the schema loses every
// query-parameter example, which is not an error but a quietly worse census:
// endpoints fall through to a class-based guess or to blocked, and nothing says
// which happened.
func paramSpecOf(param map[string]any) ParamSpec {
	schema, _ := param["schema"].(map[string]any)
	spec := ParamSpec{}
	if schema != nil {
		if declared, ok := schema["type"].(string); ok {
			spec.Type = declared
		}
		if enum, ok := schema["enum"].([]any); ok {
			spec.Enum = enum
		}
		if bound, ok := numberValue(schema["minimum"]); ok {
			spec.Minimum = &bound
		}
		if bound, ok := numberValue(schema["maximum"]); ok {
			spec.Maximum = &bound
		}
	}
	if value, ok := param["example"]; ok {
		spec.Example, spec.HasExample = value, true
		return spec
	}
	if schema != nil {
		if value, ok := schema["example"]; ok {
			spec.Example, spec.HasExample = value, true
		}
	}
	return spec
}

// numberValue reads a JSON number as a float64 whatever its decoded type. A
// published bound is only ever compared against, so the precision UseNumber
// preserves is not needed here, while accepting both shapes is: a page may have
// been written by a generator that emitted 5 and one that emitted 5.0.
func numberValue(value any) (float64, bool) {
	switch typed := value.(type) {
	case json.Number:
		parsed, err := typed.Float64()
		return parsed, err == nil
	case float64:
		return typed, true
	case int:
		return float64(typed), true
	}
	return 0, false
}

// jsonRequestBody wraps a schema in the requestBody envelope a cached page uses.
func jsonRequestBody(schema map[string]any) map[string]any {
	return map[string]any{
		"required": true,
		"content": map[string]any{
			"application/json": map[string]any{"schema": schema},
		},
	}
}

// requestBodyParams reads the required request-body properties of a document.
//
// It returns a blocked reason rather than an error for the two cases the census
// must not paper over: a body the page does not document as JSON, and a required
// property with no declared schema. There is no conforming body to build in
// either case, and a body the probe invented would be recorded as the endpoint's
// answer.
//
// The JSON media type is chosen from the sorted media type names, so a page
// documenting both application/json and a vendor type always resolves the same
// way. Go map iteration is randomised and a census that picked a different
// schema on each run could not be diffed.
func requestBodyParams(doc map[string]any) (params []NamedParam, hasBody bool, blocked string) {
	body, ok := doc["requestBody"].(map[string]any)
	if !ok {
		return nil, false, ""
	}
	content, _ := body["content"].(map[string]any)
	if len(content) == 0 {
		return nil, true, "the page documents a request body with no content"
	}
	names := make([]string, 0, len(content))
	for name := range content {
		names = append(names, name)
	}
	sort.Strings(names)

	schema := map[string]any(nil)
	for _, name := range names {
		if !strings.Contains(strings.ToLower(name), "json") {
			continue
		}
		media, _ := content[name].(map[string]any)
		schema, _ = media["schema"].(map[string]any)
		break
	}
	if schema == nil {
		return nil, true, "no JSON request body is documented, so the probe must not invent one: " +
			strings.Join(names, ", ")
	}

	required, _ := schema["required"].([]any)
	properties, _ := schema["properties"].(map[string]any)
	// The page's own order, not Go's map order, so a body and a report line up
	// with the page's property table.
	ordered := make([]string, 0, len(required))
	for _, name := range required {
		text, ok := name.(string)
		if !ok {
			continue
		}
		ordered = append(ordered, text)
	}
	for _, name := range ordered {
		node, ok := properties[name].(map[string]any)
		if !ok {
			return nil, true, "required body property " + name + " has no declared schema"
		}
		// A body property publishes its example on the property itself, which is
		// where the 144 instances in the cache are.
		spec := ParamSpec{}
		if declared, ok := node["type"].(string); ok {
			spec.Type = declared
		}
		if enum, ok := node["enum"].([]any); ok {
			spec.Enum = enum
		}
		if bound, ok := numberValue(node["minimum"]); ok {
			spec.Minimum = &bound
		}
		if bound, ok := numberValue(node["maximum"]); ok {
			spec.Maximum = &bound
		}
		if value, ok := node["example"]; ok {
			spec.Example, spec.HasExample = value, true
		}
		params = append(params, NamedParam{Name: name, Spec: spec})
	}
	return params, true, ""
}

// requiredParams reads the required query and header parameters of a document,
// with the headers the SDK signs removed.
//
// The removal is the first thing that happens to a documented parameter, before
// any resolution, and it is a removal rather than a resolution because a
// resolved x-app-key or x-signature is a credential this process invented. Those
// nine names appear on 164 of the 193 documented pages; without this filter the
// census reports almost every endpoint blocked on a header the SDK already sent,
// which is a probe bug indistinguishable from a finding.
//
// Only `in: header` is filtered. access_token and reqid are documented as headers
// too and are the probe's to resolve: access_token is the token the client
// already holds and reqid is a correlation key the server accepts any value in.
func requiredParams(doc map[string]any) (params []NamedParam, dropped int) {
	parameters, _ := doc["parameters"].([]any)
	for _, entry := range parameters {
		param, ok := entry.(map[string]any)
		if !ok {
			continue
		}
		if required, ok := param["required"].(bool); !ok || !required {
			continue
		}
		name, _ := param["name"].(string)
		if name == "" {
			continue
		}
		if in, _ := param["in"].(string); strings.EqualFold(in, "header") && isSDKSignedHeader(name) {
			dropped++
			continue
		}
		params = append(params, NamedParam{Name: name, Spec: paramSpecOf(param)})
	}
	return params, dropped
}

// endpointFromDocument builds an [Endpoint] from one cached page document.
//
// A page that carries no endpoint document is not dropped: it becomes an
// endpoint whose Unreadable reason says so, so a missing page is reported in the
// summary instead of quietly reducing the endpoint count.
func endpointFromDocument(symbol, fixture string, doc map[string]any) Endpoint {
	ep := Endpoint{Symbol: symbol, Fixture: fixture}
	if method, ok := doc["method"].(string); ok {
		ep.Method = strings.ToUpper(method)
	}
	if path, ok := doc["path"].(string); ok {
		ep.Path = path
	}
	for _, match := range pathTemplate.FindAllStringSubmatch(ep.Path, -1) {
		ep.PathParams = append(ep.PathParams, NamedParam{Name: match[1], Spec: ParamSpec{Type: "string"}})
	}
	ep.RequiredParams, ep.DroppedSDKSignedHeaders = requiredParams(doc)
	ep.BodyParams, ep.HasBody, ep.Unreadable = requestBodyParams(doc)
	return ep
}

// manifestRecord is the part of a conformance/testdata/manifest.json entry the
// census reads. The manifest is committed and response-only, which is why it
// supplies the endpoint inventory and the page URL and nothing about inputs.
type manifestRecord struct {
	ID        string `json:"id"`
	Fixture   string `json:"fixture"`
	SDKSymbol string `json:"sdkSymbol"`
	Source    struct {
		URL string `json:"url"`
	} `json:"source"`
	Documented struct {
		Method string `json:"method"`
		Path   string `json:"path"`
	} `json:"documented"`
}

// LoadEndpoints turns the committed manifest plus the cached pages into one
// [Endpoint] per fixture.
//
// The manifest is the inventory and the cache is the input schema, and neither
// alone is enough: the manifest names 193 endpoints and states nothing about what
// any of them takes, and the cache holds 348 pages and does not say which are
// endpoints.
//
// A missing cache returns an error wrapping [errCacheAbsent], which main treats
// as a skip and reports with exit 0. The cache is gitignored, so a fresh checkout
// cannot judge the corpus either way and failing would be a false alarm about
// the SDK.
func LoadEndpoints(manifestPath, cacheDir string) ([]Endpoint, error) {
	populated, err := cacheIsPopulated(cacheDir)
	if err != nil {
		return nil, err
	}
	if !populated {
		return nil, fmt.Errorf("live-probe: no docgen cache at %s: %w", cacheDir, errCacheAbsent)
	}

	raw, err := os.ReadFile(manifestPath) //nolint:gosec // G304: the path is the operator's -manifest flag, a file the census was asked to read.
	if err != nil {
		return nil, fmt.Errorf("live-probe: reading manifest %s: %w", manifestPath, err)
	}
	var manifest struct {
		Fixtures []manifestRecord `json:"fixtures"`
	}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return nil, fmt.Errorf("live-probe: parsing manifest %s: %w", manifestPath, err)
	}

	endpoints := make([]Endpoint, 0, len(manifest.Fixtures))
	for _, record := range manifest.Fixtures {
		symbol := record.SDKSymbol
		if symbol == "" {
			symbol = record.ID
		}
		fixture := record.Fixture
		if fixture == "" {
			fixture = record.ID + ".json"
		}
		md, err := readPage(cacheDir, record.Source.URL)
		if err != nil {
			return nil, err
		}
		if md == "" {
			endpoints = append(endpoints, Endpoint{
				Symbol: symbol, Fixture: fixture,
				Method: record.Documented.Method,
				Path:   record.Documented.Path,
				Unreadable: "the cached page for " + record.Source.URL +
					" is absent, so the probe knows nothing about what this endpoint takes",
			})
			continue
		}
		index, doc := selectBlock(md)
		if index < 0 {
			endpoints = append(endpoints, Endpoint{
				Symbol: symbol, Fixture: fixture,
				Method: record.Documented.Method,
				Path:   record.Documented.Path,
				Unreadable: "no JSON block in the cached page parses, so the probe knows " +
					"nothing about what this endpoint takes",
			})
			continue
		}
		ep := endpointFromDocument(symbol, fixture, doc)
		// The manifest is the authority on the method and path, because it was
		// written from the same page through the same selection rule and a
		// disagreement between the two would be a finding about the corpus rather
		// than about the endpoint. It is not the authority here, so a page that
		// documents neither keeps whatever the manifest recorded.
		if ep.Method == "" {
			ep.Method = strings.ToUpper(record.Documented.Method)
		}
		if ep.Path == "" {
			ep.Path = record.Documented.Path
		}
		endpoints = append(endpoints, ep)
	}
	return endpoints, nil
}

// cacheIsPopulated reports whether a cache directory holds at least one page.
// An unreadable directory is an error rather than an absence, because a
// permissions failure and a fresh checkout are different situations and only one
// of them should skip.
func cacheIsPopulated(cacheDir string) (bool, error) {
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("live-probe: reading cache directory %s: %w", cacheDir, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			return true, nil
		}
	}
	return false, nil
}

// statusText renders an outcome's status for a report line, so a blocked row and
// an unanswered row are never both written as a bare zero. A skipped row reads
// as its own word rather than as a status, because the one thing it must never
// look like is an answer from the server.
func statusText(o Outcome) string {
	switch {
	case gateClosed(o):
		return "skipped: mutating"
	case o.Blocked != "":
		return "blocked"
	case o.Status == 0:
		if o.Err != "" {
			return "no response (" + o.Err + ")"
		}
		return "no response"
	}
	return http.StatusText(o.Status)
}

// rowReason renders why a row is not a plain success: why the probe declined to
// call it, why it could not ask, or the classified error code when a request
// produced no status. It is empty for a 2xx, so a column of it is mostly empty on
// purpose.
func rowReason(o Outcome) string {
	switch {
	case gateClosed(o):
		return o.Skipped
	case o.Blocked != "":
		return o.Blocked
	case o.Err != "":
		return o.Err
	}
	return ""
}
