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
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/shing1211/webullapi4go/client"
)

// The census must be testable without a credential, because a test that needs
// the sandbox cannot prove the loop records a 404 correctly. Every case here
// therefore runs against an httptest server, and the only strings that look
// like credentials are the two placeholders below, which the SDK signs with and
// the server ignores.
const (
	probeTestCredentialA = "live-probe-test-credential-a"
	probeTestCredentialB = "live-probe-test-credential-b"
)

// probeDiscoveredAccount stands in for the account id phase 0 would discover.
// It is obviously synthetic so a value that escaped into a test failure could
// not be read as a real account, and it is a constant so a case can assert that
// the discovered value reached the wire in preference to the page's example.
const probeDiscoveredAccount = "live-probe-discovered-account"

// probePageExampleAccount stands in for the account_id a published page puts in
// its schema, which is an account nobody holds. It is obviously synthetic so
// the case reads as "the page's value and the discovered value differ" without
// this file carrying an identifier Webull issued.
const probePageExampleAccount = "live-probe-page-example-account"

// syntheticAccountID is a number shaped like a Webull account id, used by the
// cases that hold the reduction to placeholders: a body that carries a reading
// has to be refused for the same reason a body carrying a price is, and the
// reading has to be one a reader would recognise as a reading.
//
// It is a constant, and this is its comment, because it was a bare literal in
// four files until a whole-branch review asked which it was. It is NOT an
// identifier Webull issued: it is not in the docgen cache, it is not a documented
// example, it appears nowhere under docs/ or conformance/testdata, and it
// predates the live run (it arrived with d13085e, in the reduction itself).
// Commit fc75d6c removed the account id the trading schema publishes from
// census_test.go on the principle that no hand-written file here should carry
// one, and that principle is what this constant records: a reader who sees the
// digits know they are a fixture, because something now says so.
const syntheticAccountID = "9110101000000000001"

// newProbeTestClient returns a client aimed at baseURL. Auto-token is off and
// retry is off so that a case observes exactly one request and one status.
func newProbeTestClient(t *testing.T, baseURL string) *client.Client {
	t.Helper()
	cl, err := client.New(
		client.WithCredentials(probeTestCredentialA, probeTestCredentialB),
		client.WithRegion(client.HK),
		client.WithEnvironment(client.Sandbox),
		client.WithBaseURL(baseURL),
		client.WithAutoToken(false),
		client.WithoutRetry(),
	)
	if err != nil {
		t.Fatalf("client.New: %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })
	return cl
}

// probeSandboxBaseURL is the base URL of the one sandbox host AGENTS.md permits in
// committed material, written as a constant so the tests that need to present a
// sandbox cannot drift from the host the live run will use and cannot be
// read as an arbitrary string.
const probeSandboxBaseURL = "https://api.sandbox.webull.hk"

// newProbeTestClientPresenting returns a client whose base URL is presentedURL
// while every request it sends is dialled at target.
//
// The mutating gate judges the host in the base URL, so a test that exercises the
// gate's OPEN path has to present a host the SDK derives for a sandbox. A client
// aimed straight at 127.0.0.1 presents "127.0.0.1", which the gate correctly
// refuses as not a sandbox - so with only that helper the open path could be
// reached only through the override, and the case a live run actually takes
// (opt-in plus the sandbox host) would go unproven.
//
// The DIAL is redirected, not the name resolved: no DNS lookup happens and no
// request leaves the process, and the SDK still signs over the presented host,
// which is what a live run signs and what these tests are about. Presenting the
// host is the point; where the bytes land is the test server's business.
func newProbeTestClientPresenting(t *testing.T, presentedURL string, target *httptest.Server) *client.Client {
	t.Helper()
	targetURL, err := url.Parse(target.URL)
	if err != nil {
		t.Fatalf("parsing the test server URL: %v", err)
	}
	cl, err := client.New(
		client.WithCredentials(probeTestCredentialA, probeTestCredentialB),
		client.WithRegion(client.HK),
		client.WithEnvironment(client.Sandbox),
		client.WithBaseURL(presentedURL),
		client.WithHTTPClient(&http.Client{Transport: &dialAtTarget{
			target: targetURL,
			next:   http.DefaultTransport,
		}}),
		client.WithAutoToken(false),
		client.WithoutRetry(),
	)
	if err != nil {
		t.Fatalf("client.New: %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })
	return cl
}

// dialAtTarget sends every request to one address whatever host its URL carries,
// so a client can present a Webull host while the bytes go to a test server.
type dialAtTarget struct {
	target *url.URL
	next   http.RoundTripper
}

func (d *dialAtTarget) RoundTrip(req *http.Request) (*http.Response, error) {
	redirected := req.Clone(req.Context())
	redirected.URL.Scheme = d.target.Scheme
	redirected.URL.Host = d.target.Host
	return d.next.RoundTrip(redirected)
}

// The loop's job is to record outcomes, so the case that matters most is the one
// that proves it distinguishes "the endpoint answered 404" from "the probe could
// not ask". A blocked endpoint must never reach the network: a fabricated value
// would be rejected by the server and the rejection recorded as the endpoint's
// own answer.
func TestCensusRecordsNonSuccessStatuses(t *testing.T) {
	var mu sync.Mutex
	var seen []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		seen = append(seen, r.URL.Path)
		mu.Unlock()
		switch r.URL.Path {
		case "/ok":
			_, _ = w.Write([]byte(`{"a":1}`))
		case "/gone":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusUnauthorized)
		}
	}))
	defer srv.Close()

	cl := newProbeTestClient(t, srv.URL)
	got, err := Census(context.Background(), cl, []Endpoint{
		{Symbol: "ok", Method: "GET", Path: "/ok"},
		{Symbol: "gone", Method: "GET", Path: "/gone"},
		{Symbol: "blocked", Method: "GET", Path: "/nope",
			RequiredParams: []NamedParam{{Name: "account_id"}}},
	})
	if err != nil {
		t.Fatalf("Census: %v", err)
	}
	by := map[string]Outcome{}
	for _, o := range got {
		by[o.Symbol] = o
	}
	if by["ok"].Status != http.StatusOK {
		t.Errorf("ok status = %d, want 200", by["ok"].Status)
	}
	if by["gone"].Status != http.StatusNotFound {
		t.Errorf("gone status = %d, want 404", by["gone"].Status)
	}
	// The unresolvable-parameter case must never reach the network: it is not an
	// endpoint that failed to answer, it is an endpoint the probe could not ask.
	if by["blocked"].Blocked == "" {
		t.Error("blocked endpoint has no Blocked reason")
	}
	if !strings.Contains(by["blocked"].Blocked, "account_id") {
		t.Errorf("blocked reason = %q, want it to name account_id", by["blocked"].Blocked)
	}
	if by["blocked"].Status != 0 {
		t.Errorf("blocked endpoint has status %d; it should not have been called",
			by["blocked"].Status)
	}
	mu.Lock()
	defer mu.Unlock()
	for _, p := range seen {
		if p == "/nope" {
			t.Error("the probe called an endpoint whose parameters it could not resolve")
		}
	}
}

// A census row that cannot be traced back to a request is not evidence, so the
// loop has to report the method and path it actually used, and a transport
// failure has to be distinguishable from an HTTP answer.
func TestCensusRecordsTheRequestItAttempted(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	cl := newProbeTestClient(t, srv.URL)
	got, err := Census(context.Background(), cl, []Endpoint{{
		Symbol:  "trade.ListAccounts",
		Fixture: "trading/GET-trading-accounts-list.json",
		Method:  "GET",
		Path:    "/trading/accounts/list",
	}})
	if err != nil {
		t.Fatalf("Census: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
	o := got[0]
	if o.Symbol != "trade.ListAccounts" || o.Method != "GET" || o.Path != "/trading/accounts/list" {
		t.Errorf("outcome = %+v, want the endpoint it attempted", o)
	}
	if o.Area != "trading" {
		t.Errorf("Area = %q, want it derived from the fixture name", o.Area)
	}
	if o.Blocked != "" {
		t.Errorf("Blocked = %q, want empty for an endpoint that answered", o.Blocked)
	}
	if o.Err != "" {
		t.Errorf("Err = %q, want empty for a 2xx", o.Err)
	}
}

// The deliverable: a request built from what the page documents, asserted on
// the wire. The query and the body are the whole product of this task, so they
// are asserted byte for byte rather than through a decode.
//
// The gate is opened here, and only here, because the case is about the bytes on
// the wire and a request that is never sent produces no bytes. Both of the gate's
// conditions are satisfied for this test alone - the opt-in is set, and the
// client presents the sandbox host while its dial is redirected at a local
// httptest server, so nothing leaves the process. t.Setenv scopes the variable to
// this test; the gate's own behaviour is pinned in mutate_test.go against the
// same server.
func TestCensusSendsTheDocumentedQueryAndBody(t *testing.T) {
	t.Setenv(MutateOptInEnv, MutateOptInValue)
	t.Cleanup(func() { SetAccountID("") })
	SetAccountID(probeDiscoveredAccount)

	type request struct {
		method string
		path   string
		query  string
		body   string
	}
	var mu sync.Mutex
	var got []request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, r.ContentLength)
		if r.ContentLength > 0 {
			_, _ = r.Body.Read(buf)
		}
		mu.Lock()
		got = append(got, request{r.Method, r.URL.Path, r.URL.RawQuery, string(buf)})
		mu.Unlock()
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	cl := newProbeTestClientPresenting(t, probeSandboxBaseURL, srv)

	outcomes, err := Census(context.Background(), cl, []Endpoint{{
		Symbol:  "trade.PreviewOrder",
		Fixture: "trading/POST-trading-orders-preview.json",
		Method:  "POST",
		Path:    "/trading/orders/preview",
		RequiredParams: []NamedParam{
			{Name: "account_id"},
			{Name: "symbol", Spec: ParamSpec{Type: "string"}},
			{Name: "count", Spec: ParamSpec{Type: "integer"}},
			{Name: "dry_run", Spec: ParamSpec{Type: "boolean"}},
		},
		BodyParams: []NamedParam{
			{Name: "account_id"},
			{Name: "symbol", Spec: ParamSpec{Type: "string", Example: "AAPL", HasExample: true}},
			{Name: "quantity", Spec: ParamSpec{Type: "string", Example: "1", HasExample: true}},
		},
		HasBody: true,
	}})
	if err != nil {
		t.Fatalf("Census: %v", err)
	}
	if outcomes[0].Blocked != "" {
		t.Fatalf("Blocked = %q, want the request to be buildable", outcomes[0].Blocked)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 {
		t.Fatalf("sent %d requests, want 1", len(got))
	}
	want := request{
		method: "POST",
		path:   "/trading/orders/preview",
		query:  "account_id=" + probeDiscoveredAccount + "&count=1&dry_run=true&symbol=AAPL",
		body:   `{"account_id":"` + probeDiscoveredAccount + `","quantity":"1","symbol":"AAPL"}`,
	}
	if got[0] != want {
		t.Errorf("request =\n  %+v\nwant\n  %+v", got[0], want)
	}
}

// Landmine 1. The nine signed headers appear on 164 of the 193 documented pages
// and are supplied by the SDK, so they are dropped before resolution. A row
// naming one of them is a filter bug, and filling x-app-key would be fabricating
// a credential.
func TestCensusDropsSDKSignedHeaderParameters(t *testing.T) {
	headers := []string{
		"x-app-key", "x-app-secret", "x-timestamp", "x-access-token", "x-signature",
		"x-signature-algorithm", "x-signature-nonce", "x-signature-version", "x-version",
	}
	for _, name := range headers {
		t.Run(name, func(t *testing.T) {
			// Each header is given a published example, which is the case that
			// would otherwise resolve to a value and leave the process.
			doc := map[string]any{
				"method": "get",
				"path":   "/x",
				"parameters": []any{
					map[string]any{"name": name, "in": "header", "required": true,
						"schema":  map[string]any{"type": "string"},
						"example": "would-be-a-credential"},
				},
				"responses": map[string]any{},
			}
			ep := endpointFromDocument("t/x", "t/x.json", doc)
			if len(ep.RequiredParams) != 0 {
				t.Fatalf("RequiredParams = %+v, want the signed header dropped", ep.RequiredParams)
			}
			got := prepare(ep)
			if got.Blocked != "" {
				t.Errorf("Blocked = %q, want a buildable request", got.Blocked)
			}
			if len(got.Query) != 0 {
				t.Errorf("query = %v, want no %s on the wire", got.Query, name)
			}
		})
	}
}

// The filter is scoped to the SDK's nine names. access_token and reqid are
// documented as headers too, and they are the probe's to resolve.
func TestCensusKeepsHeaderParametersTheSDKDoesNotSign(t *testing.T) {
	doc := map[string]any{
		"method": "get",
		"path":   "/x",
		"parameters": []any{
			map[string]any{"name": "access_token", "in": "header", "required": true,
				"schema": map[string]any{"type": "string"}},
			map[string]any{"name": "reqid", "in": "header", "required": true,
				"schema": map[string]any{"type": "string"}, "example": "req-1"},
		},
		"responses": map[string]any{},
	}
	ep := endpointFromDocument("t/x", "t/x.json", doc)
	names := make([]string, 0, len(ep.RequiredParams))
	for _, p := range ep.RequiredParams {
		names = append(names, p.Name)
	}
	sort.Strings(names)
	if want := []string{"access_token", "reqid"}; !reflect.DeepEqual(names, want) {
		t.Errorf("names = %v, want %v", names, want)
	}
}

// Landmine 2. The corpus carries 433 examples on the parameter object and none
// at parameter["schema"]["example"], so a reader that looks only at the schema
// loses every query-parameter example and quietly drops 22 endpoints with no
// error at all.
func TestParameterExampleComesOffTheParameterObject(t *testing.T) {
	cases := []struct {
		name  string
		param map[string]any
		want  any
	}{
		{
			name: "example on the parameter wins",
			param: map[string]any{"name": "period", "in": "query", "required": true,
				"schema":  map[string]any{"type": "string"},
				"example": "day"},
			want: "day",
		},
		{
			name: "schema example is the fallback",
			param: map[string]any{"name": "period", "in": "query", "required": true,
				"schema": map[string]any{"type": "string", "example": "week"}},
			want: "week",
		},
		{
			name: "no example anywhere leaves HasExample false",
			param: map[string]any{"name": "period", "in": "query", "required": true,
				"schema": map[string]any{"type": "string"}},
			want: nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			spec := paramSpecOf(tc.param)
			if tc.want == nil {
				if spec.HasExample {
					t.Errorf("HasExample = true, Example = %#v; want no example", spec.Example)
				}
				return
			}
			if !spec.HasExample {
				t.Fatalf("HasExample = false, want %#v", tc.want)
			}
			if !reflect.DeepEqual(spec.Example, tc.want) {
				t.Errorf("Example = %#v, want %#v", spec.Example, tc.want)
			}
		})
	}
}

// A numeric example must reach the wire as the number the page published, not
// as a float64 that re-encodes differently.
func TestParameterNumericExampleKeepsItsLiteral(t *testing.T) {
	doc := map[string]any{
		"method": "get",
		"path":   "/x",
		"parameters": []any{
			map[string]any{"name": "count", "in": "query", "required": true,
				"schema": map[string]any{"type": "integer"}, "example": 10},
		},
		"responses": map[string]any{},
	}
	got := prepare(endpointFromDocument("t/x", "t/x.json", doc))
	if got.Query.Get("count") != "10" {
		t.Errorf("count = %q, want the page's own 10", got.Query.Get("count"))
	}
}

// account_id is the caller's identity and the pages publish an example for it
// that is an account nobody holds. Param deliberately lets the example win
// (params_test.go:274) and its own comment names the fix: "a discovered account
// threaded in ahead of this call". The census is where that threading happens,
// so it intercepts the name before Param and never sends the page's value.
//
// The page's example is written as an obviously synthetic value rather than the
// identifier the published schema carries. The case only needs the two to
// differ, and no hand-written file in this repository should carry an
// identifier Webull issued, however well published it is.
func TestCensusThreadsTheDiscoveredAccountAheadOfThePageExample(t *testing.T) {
	t.Cleanup(func() { SetAccountID("") })
	SetAccountID(probeDiscoveredAccount)

	doc := map[string]any{
		"method": "post",
		"path":   "/trading/orders/cancel",
		"parameters": []any{
			map[string]any{"name": "account_id", "in": "query", "required": true,
				"schema": map[string]any{"type": "string"}, "example": probePageExampleAccount},
		},
		"requestBody": jsonRequestBody(map[string]any{
			"type":       "object",
			"required":   []any{"account_id"},
			"properties": map[string]any{"account_id": map[string]any{"type": "string", "example": probePageExampleAccount}},
		}),
		"responses": map[string]any{},
	}
	got := prepare(endpointFromDocument("t/cancel", "trading/x.json", doc))
	if got.Blocked != "" {
		t.Fatalf("Blocked = %q", got.Blocked)
	}
	if q := got.Query.Get("account_id"); q != probeDiscoveredAccount {
		t.Errorf("query account_id = %q, want the discovered %q", q, probeDiscoveredAccount)
	}
	var body map[string]any
	if err := json.Unmarshal(got.Body, &body); err != nil {
		t.Fatalf("body: %v", err)
	}
	if body["account_id"] != probeDiscoveredAccount {
		t.Errorf("body account_id = %#v, want the discovered %q", body["account_id"], probeDiscoveredAccount)
	}
}

// With no discovered account the same endpoint is blocked, and the reason says
// so rather than naming a generic resolution failure.
func TestCensusBlocksWhenAccountDiscoveryDidNotRun(t *testing.T) {
	t.Cleanup(func() { SetAccountID("") })
	SetAccountID("")

	doc := map[string]any{
		"method": "post", "path": "/trading/orders/cancel",
		"parameters": []any{map[string]any{"name": "account_id", "in": "query",
			"required": true, "schema": map[string]any{"type": "string"}}},
		"responses": map[string]any{},
	}
	got := prepare(endpointFromDocument("t/cancel", "trading/x.json", doc))
	if !strings.Contains(got.Blocked, "account discovery") {
		t.Errorf("Blocked = %q, want it to name account discovery", got.Blocked)
	}
}

// The optional cursor is omitted rather than sent empty, because an empty value
// is a different request from an absent key.
func TestCensusOmitsTheOptionalCursor(t *testing.T) {
	t.Cleanup(func() { SetSessionToken("") })
	SetSessionToken("probe-placeholder")

	doc := map[string]any{
		"method": "get", "path": "/x",
		"parameters": []any{
			map[string]any{"name": "pagination_key", "in": "query", "required": true,
				"schema": map[string]any{"type": "string"}},
			map[string]any{"name": "access_token", "in": "header", "required": true,
				"schema": map[string]any{"type": "string"}},
		},
		"responses": map[string]any{},
	}
	got := prepare(endpointFromDocument("t/x", "t/x.json", doc))
	if got.Blocked != "" {
		t.Fatalf("Blocked = %q", got.Blocked)
	}
	if _, present := got.Query["pagination_key"]; present {
		t.Error("pagination_key is on the query; the cursor must be omitted, not emptied")
	}
	if got.Query.Get("access_token") != "probe-placeholder" {
		t.Errorf("access_token = %q, want the recorded session token", got.Query.Get("access_token"))
	}
}

// A required string with no example and no closed set is unresolvable, in a
// body exactly as in a query. Sending "" is a request the endpoint answers for a
// reason that says nothing about whether it is reachable.
func TestCensusBlocksAnUndocumentedRequiredBodyString(t *testing.T) {
	doc := map[string]any{
		"method": "post", "path": "/broker/funding/ach/relationships/create",
		"requestBody": jsonRequestBody(map[string]any{
			"type":     "object",
			"required": []any{"processor_token"},
			"properties": map[string]any{
				"processor_token": map[string]any{"type": "string", "description": "issued by the processor"},
			},
		}),
		"responses": map[string]any{},
	}
	got := prepare(endpointFromDocument("t/ach", "broker-fd-us/x.json", doc))
	if !strings.Contains(got.Blocked, "processor_token") {
		t.Errorf("Blocked = %q, want it to name processor_token", got.Blocked)
	}
	if len(got.Body) != 0 {
		t.Errorf("body = %s, want no body for a blocked endpoint", got.Body)
	}
}

// No page in the cache declares minItems on a required array, so an empty array
// is a conforming instance and the probe can ask. It is also the only value that
// cannot place an order, so 50 POST endpoints are walked without the census
// mutating anything. The flag is recorded because 16 endpoints rest on it.
func TestCensusSynthesisesAnEmptyRequiredArray(t *testing.T) {
	doc := map[string]any{
		"method": "post", "path": "/trading/orders/place",
		"requestBody": jsonRequestBody(map[string]any{
			"type":     "object",
			"required": []any{"new_orders"},
			"properties": map[string]any{
				"new_orders": map[string]any{"type": "array",
					"items": map[string]any{"type": "object"}},
			},
		}),
		"responses": map[string]any{},
	}
	got := prepare(endpointFromDocument("t/place", "trading/x.json", doc))
	if got.Blocked != "" {
		t.Fatalf("Blocked = %q", got.Blocked)
	}
	if string(got.Body) != `{"new_orders":[]}` {
		t.Errorf("body = %s, want an empty required array", got.Body)
	}
	if want := []string{"new_orders"}; !reflect.DeepEqual(got.SynthesisedEmptyArrays, want) {
		t.Errorf("SynthesisedEmptyArrays = %v, want %v", got.SynthesisedEmptyArrays, want)
	}
}

// Two of the 193 documented bodies are not JSON at all, so there is no schema
// to build a conforming body from and the probe must not invent one.
func TestCensusBlocksANonJSONRequestBody(t *testing.T) {
	cases := []struct {
		name        string
		requestBody map[string]any
	}{
		{"multipart", map[string]any{"content": map[string]any{
			"multipart/form-data": map[string]any{"schema": map[string]any{"type": "object"}}}}},
		{"form urlencoded", map[string]any{"content": map[string]any{
			"application/x-www-form-urlencoded": map[string]any{"schema": map[string]any{"type": "object"}}}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			doc := map[string]any{
				"method": "post", "path": "/x",
				"requestBody": tc.requestBody,
				"responses":   map[string]any{},
			}
			got := prepare(endpointFromDocument("t/x", "t/x.json", doc))
			if !strings.Contains(got.Blocked, "no JSON request body") {
				t.Errorf("Blocked = %q, want it to name the missing JSON body", got.Blocked)
			}
		})
	}
}

// A body whose schema declares no required property is still a body: {} is the
// minimal instance the documented schema admits.
func TestCensusSendsAnEmptyObjectWhenNothingIsRequired(t *testing.T) {
	doc := map[string]any{
		"method": "post", "path": "/x",
		"requestBody": jsonRequestBody(map[string]any{
			"type":       "object",
			"properties": map[string]any{"note": map[string]any{"type": "string"}},
		}),
		"responses": map[string]any{},
	}
	got := prepare(endpointFromDocument("t/x", "t/x.json", doc))
	if got.Blocked != "" {
		t.Fatalf("Blocked = %q", got.Blocked)
	}
	if string(got.Body) != `{}` {
		t.Errorf("body = %s, want {}", got.Body)
	}
}

// A GET carries no body even when the page documents one, and a templated path
// segment is filled from the same resolution the query uses.
func TestCensusFillsAPathTemplateAndOmitsABodyForGET(t *testing.T) {
	t.Cleanup(func() { SetAccountID("") })
	SetAccountID(probeDiscoveredAccount)

	doc := map[string]any{
		"method": "get", "path": "/trading/accounts/{account_id}/positions",
		"responses": map[string]any{},
	}
	ep := endpointFromDocument("t/positions", "trading/x.json", doc)
	if len(ep.PathParams) != 1 || ep.PathParams[0].Name != "account_id" {
		t.Fatalf("PathParams = %+v, want the templated segment to become a parameter", ep.PathParams)
	}
	got := prepare(ep)
	if got.Blocked != "" {
		t.Fatalf("Blocked = %q", got.Blocked)
	}
	if want := "/trading/accounts/" + probeDiscoveredAccount + "/positions"; got.Path != want {
		t.Errorf("path = %q, want %q", got.Path, want)
	}
	if len(got.Body) != 0 {
		t.Errorf("body = %s, want none for a GET", got.Body)
	}
}

// A path template nothing resolves is blocked rather than sent with a literal
// brace in it, which is a different URL from the documented one.
func TestCensusBlocksAnUnresolvablePathTemplate(t *testing.T) {
	doc := map[string]any{"method": "get", "path": "/x/{order_id}", "responses": map[string]any{}}
	got := prepare(endpointFromDocument("t/x", "t/x.json", doc))
	if !strings.Contains(got.Blocked, "order_id") {
		t.Errorf("Blocked = %q, want it to name order_id", got.Blocked)
	}
}

// Reachable and blocked are different numbers and must never be summed into
// one. The summary is what a reader of census.json looks at first, so it has to
// carry the per-area split the areas diverge on.
func TestSummariseSeparatesReachableFromBlocked(t *testing.T) {
	outcomes := []Outcome{
		{Symbol: "a", Area: "fundamentals", Status: 200},
		{Symbol: "b", Area: "fundamentals", Status: 404},
		{Symbol: "c", Area: "display-solution", Blocked: "body:bank_code"},
		{Symbol: "d", Area: "display-solution", Status: 403},
		{Symbol: "e", Area: "fundamentals", Err: "transport"},
	}
	got := Summarise(outcomes)
	if got.Total != 5 {
		t.Errorf("Total = %d, want 5", got.Total)
	}
	if got.Reachable != 3 {
		t.Errorf("Reachable = %d, want 3 (a 200, a 404 and a 403 are all answers)", got.Reachable)
	}
	if got.Blocked != 1 {
		t.Errorf("Blocked = %d, want 1", got.Blocked)
	}
	if got.Unanswered != 1 {
		t.Errorf("Unanswered = %d, want 1 (a transport failure is not an answer)", got.Unanswered)
	}
	if got.Area["fundamentals"].Total != 3 || got.Area["fundamentals"].Reachable != 2 {
		t.Errorf("fundamentals = %+v, want 2 of 3 reachable", got.Area["fundamentals"])
	}
	if got.Area["display-solution"].Blocked != 1 {
		t.Errorf("display-solution = %+v, want 1 blocked", got.Area["display-solution"])
	}
}

// The cache filename is a copy of three lines of _common.fetch, because calling
// fetch would create the directory and re-fetch on a miss. If the derivation
// drifts, every page is a miss and the census silently walks nothing.
func TestCacheFilenameMatchesTheDocgenDerivation(t *testing.T) {
	cases := []struct {
		url  string
		want string
	}{
		{"https://developer.webull.hk/apis/docs/reference/create-token.md",
			"https_developer_webull_hk_apis_docs_reference_create_token_md.md"},
		{"https://developer.webull.com/apis/docs/reference/a b?c=1",
			"https_developer_webull_com_apis_docs_reference_a_b_c_1.md"},
	}
	for _, tc := range cases {
		if got := cacheFilename(tc.url); got != tc.want {
			t.Errorf("cacheFilename(%q) = %q, want %q", tc.url, got, tc.want)
		}
	}
	long := "https://developer.webull.hk/apis/docs/reference/" + strings.Repeat("z", 200) + ".md"
	if got := cacheFilename(long); len(got) != 120+len(".md") {
		t.Errorf("cacheFilename of a long url is %d chars, want the last 120 plus .md", len(got))
	}
}

// A cached page's JSON block is a single-endpoint document. select_block has to
// keep the "first block lacks path" case that a hand-rolled parser gets wrong:
// the first block is often a guide example, and taking it yields an endpoint
// count of zero across the whole corpus with no error at all.
func TestSelectBlockPrefersTheEndpointDocument(t *testing.T) {
	cases := []struct {
		name      string
		md        string
		wantIndex int
		wantPath  string
	}{
		{
			name:      "first block is the endpoint document",
			md:        "```json\n{\"path\":\"/a\",\"responses\":{}}\n```\n",
			wantIndex: 0,
			wantPath:  "/a",
		},
		{
			name: "a leading guide block is skipped",
			md: "prose\n```json\n{\"title\":\"guide\",\"examples\":[]}\n```\n" +
				"more prose\n```json\n{\"path\":\"/b\",\"responses\":{}}\n```\n",
			wantIndex: 1,
			wantPath:  "/b",
		},
		{
			name:      "a leading block that does not parse is skipped",
			md:        "```json\n{not json\n```\n```json\n{\"path\":\"/c\",\"responses\":{}}\n```\n",
			wantIndex: 1,
			wantPath:  "/c",
		},
		{
			name:      "no block at all",
			md:        "nothing here",
			wantIndex: -1,
			wantPath:  "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			index, doc := selectBlock(tc.md)
			if index != tc.wantIndex {
				t.Fatalf("index = %d, want %d", index, tc.wantIndex)
			}
			if tc.wantPath == "" {
				if doc != nil {
					t.Errorf("doc = %+v, want nil", doc)
				}
				return
			}
			if doc["path"] != tc.wantPath {
				t.Errorf("path = %v, want %q", doc["path"], tc.wantPath)
			}
		})
	}
}

// The method in a cached document is lower case and the census sends upper, the
// same normalisation gen_fixtures.py applies before it writes the manifest.
func TestEndpointFromDocumentNormalisesTheMethod(t *testing.T) {
	ep := endpointFromDocument("t/x", "t/x.json", map[string]any{
		"method": "post", "path": "/x", "responses": map[string]any{},
	})
	if ep.Method != "POST" {
		t.Errorf("Method = %q, want POST", ep.Method)
	}
}

// -assume-discovered substitutes a placeholder for a real account and a real
// session token, so it must be unusable outside a dry run. A run that reached
// the sandbox with a placeholder would record 59 endpoints answering 404 against
// an account nobody holds, which is the exact fabrication this probe exists to
// rule out.
func TestAssumeDiscoveredIsRefusedOutsideADryRun(t *testing.T) {
	err := run(context.Background(), []string{
		"-assume-discovered", "-out", filepath.Join(t.TempDir(), "census.json"),
	})
	if err == nil {
		t.Fatal("run accepted -assume-discovered without -dry-run")
	}
	if !strings.Contains(err.Error(), "-dry-run") {
		t.Errorf("err = %v, want it to say the flag needs -dry-run", err)
	}
}

// An absent cache is a skip, not a failure: it is gitignored, so a fresh
// checkout cannot judge the corpus either way. The same rule the
// conformance-fixtures target uses.
func TestLoadEndpointsReportsAnAbsentCache(t *testing.T) {
	dir := t.TempDir()
	manifest := filepath.Join(dir, "manifest.json")
	body := `{"fixtures":[{"id":"t/GET-x","fixture":"t/GET-x.json","area":"t",` +
		`"sdkSymbol":"x.Y","source":{"url":"https://example.test/x.md"},` +
		`"documented":{"method":"GET","path":"/x"}}]}`
	if err := os.WriteFile(manifest, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	empty := filepath.Join(dir, "cache")
	if err := os.MkdirAll(empty, 0o750); err != nil {
		t.Fatal(err)
	}
	_, err := LoadEndpoints(manifest, empty)
	if !errors.Is(err, errCacheAbsent) {
		t.Fatalf("err = %v, want errCacheAbsent", err)
	}
	if !strings.Contains(err.Error(), empty) {
		t.Errorf("err = %v, want it to name the cache directory it looked in", err)
	}
}

// With the cache present every fixture becomes an endpoint, and the two numbers
// have to agree or the corpus has moved under the census.
func TestLoadEndpointsWalksEveryCachedFixture(t *testing.T) {
	dir := t.TempDir()
	cache := filepath.Join(dir, "cache")
	if err := os.MkdirAll(cache, 0o750); err != nil {
		t.Fatal(err)
	}
	const url = "https://example.test/reference/get-thing.md"
	page := "```json\n" + `{"method":"get","path":"/thing",` +
		`"parameters":[{"name":"symbol","in":"query","required":true,` +
		`"schema":{"type":"string"},"example":"AAPL"}],"responses":{}}` + "\n```\n"
	if err := os.WriteFile(filepath.Join(cache, cacheFilename(url)), []byte(page), 0o600); err != nil {
		t.Fatal(err)
	}
	manifest := filepath.Join(dir, "manifest.json")
	body := `{"fixtures":[{"id":"t/GET-thing","fixture":"t/GET-thing.json","area":"t",` +
		`"sdkSymbol":"data.GetThing","source":{"url":"` + url + `"},` +
		`"documented":{"method":"GET","path":"/thing"}}]}`
	if err := os.WriteFile(manifest, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}

	endpoints, err := LoadEndpoints(manifest, cache)
	if err != nil {
		t.Fatalf("LoadEndpoints: %v", err)
	}
	if len(endpoints) != 1 {
		t.Fatalf("len(endpoints) = %d, want 1", len(endpoints))
	}
	ep := endpoints[0]
	if ep.Symbol != "data.GetThing" || ep.Method != "GET" || ep.Path != "/thing" {
		t.Errorf("endpoint = %+v", ep)
	}
	if len(ep.RequiredParams) != 1 || ep.RequiredParams[0].Name != "symbol" {
		t.Errorf("RequiredParams = %+v, want symbol", ep.RequiredParams)
	}
	if got := prepare(ep); got.Query.Get("symbol") != "AAPL" {
		t.Errorf("symbol = %q, want the page's AAPL", got.Query.Get("symbol"))
	}
}

// The two landmine signatures have to be separated from every other blocked
// class, because a row blocked on a signed header or on a parameter the page
// published an example for is a probe bug and would otherwise sit in the same
// column as a real 404. Each case plants one row of one class and asserts that
// it lands in that class and no other.
func TestDiagnoseFilesEachBlockedClassSeparately(t *testing.T) {
	cases := []struct {
		name string
		row  Outcome
		want string
	}{
		{
			name: "a signed header that survived the filter",
			row:  Outcome{Fixture: "t/a.json", Blocked: `parameter x-app-key: "x-app-key" is a header the SDK signs`},
			want: "sdk-signed header not filtered",
		},
		{
			name: "a parameter the page published an example for",
			row:  Outcome{Fixture: "t/b.json", Blocked: "parameter period: no value is synthesable for parameter \"period\" of declared type \"string\""},
			want: "example published but not read",
		},
		{
			name: "a required body string the page never documented",
			row:  Outcome{Fixture: "t/c.json", Blocked: "body property bank_code: ..."},
			want: "undocumented required body string",
		},
		{
			name: "a body that is not JSON",
			row:  Outcome{Fixture: "t/d.json", Blocked: "no JSON request body is documented, so the probe must not invent one: multipart/form-data"},
			want: "request body is not JSON",
		},
		{
			name: "an account id phase 0 did not find",
			row:  Outcome{Fixture: "t/e.json", Blocked: "parameter account_id: " + blockedAccountReason},
			want: "account id unavailable",
		},
		{
			name: "a session token that was never recorded",
			row:  Outcome{Fixture: "t/f.json", Blocked: `parameter access_token: live-probe: "access_token" is the session token: no token is recorded`},
			want: "session token unavailable",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := diagnose([]Outcome{tc.row})
			if len(got) != 1 {
				t.Fatalf("diagnose returned %d classes, want exactly 1: %+v", len(got), got)
			}
			if _, ok := got[tc.want]; !ok {
				t.Fatalf("diagnose classes = %v, want %q", classNames(got), tc.want)
			}
		})
	}
}

// A body resting on a synthesised empty array is a choice, not a blockage, so it
// is reported from the document rather than from a blocked row. Read from the
// built body it would vanish for every account-scoped endpoint, because those
// are blocked before the body is built.
func TestUndeterminedArraysReadsTheDocumentNotTheBuiltBody(t *testing.T) {
	ep := Endpoint{
		Fixture: "trading/POST-trading-orders-place.json",
		Method:  "POST",
		Path:    "/trading/orders/place",
		BodyParams: []NamedParam{
			{Name: "account_id"},
			{Name: "symbol", Spec: ParamSpec{Type: "string", Example: "AAPL", HasExample: true}},
			{Name: "new_orders", Spec: ParamSpec{Type: "array"}},
		},
		HasBody:    true,
		Unreadable: blockedAccountReason,
	}
	if got := prepare(ep); got.Blocked == "" {
		t.Fatal("prepare built a request for an unreadable endpoint")
	}
	want := []string{"new_orders"}
	if got := undeterminedArrays(ep); !reflect.DeepEqual(got, want) {
		t.Errorf("undeterminedArrays = %v, want %v", got, want)
	}
}

// An array property that does publish an example is not synthesised: the page
// said what to send, and an empty array would be a different request.
func TestUndeterminedArraysLeavesAnExampledArrayAlone(t *testing.T) {
	ep := Endpoint{BodyParams: []NamedParam{
		{Name: "symbols", Spec: ParamSpec{Type: "array", Example: []any{"AAPL"}, HasExample: true}},
		{Name: "tags", Spec: ParamSpec{Type: "array", Enum: []any{"a"}}},
	}}
	if got := undeterminedArrays(ep); len(got) != 0 {
		t.Errorf("undeterminedArrays = %v, want none", got)
	}
}

// classNames returns the diagnostic class names, for a failure message.
func classNames(set map[string]DiagnosticSet) []string {
	names := make([]string, 0, len(set))
	for name := range set {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// A page missing from the cache is reported, not skipped silently: the manifest
// names 193 endpoints and 16 of its pages are outside the corpus, and a reader has
// to be able to tell those from a page that resolved.
func TestLoadEndpointsReportsAPageMissingFromTheCache(t *testing.T) {
	dir := t.TempDir()
	cache := filepath.Join(dir, "cache")
	if err := os.MkdirAll(cache, 0o750); err != nil {
		t.Fatal(err)
	}
	// One unrelated page, so the cache is populated and the manifest's own page
	// is the missing one. An empty cache is the skip case, pinned separately.
	present := "```json\n{\"path\":\"/other\",\"responses\":{}}\n```\n"
	if err := os.WriteFile(filepath.Join(cache, cacheFilename("https://example.test/other.md")),
		[]byte(present), 0o600); err != nil {
		t.Fatal(err)
	}
	manifest := filepath.Join(dir, "manifest.json")
	body := `{"fixtures":[{"id":"t/GET-x","fixture":"t/GET-x.json","area":"t",` +
		`"sdkSymbol":"x.Y","source":{"url":"https://example.test/absent.md"},` +
		`"documented":{"method":"GET","path":"/x"}}]}`
	if err := os.WriteFile(manifest, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	endpoints, err := LoadEndpoints(manifest, cache)
	if err != nil {
		t.Fatalf("LoadEndpoints: %v", err)
	}
	if len(endpoints) != 1 {
		t.Fatalf("len(endpoints) = %d, want the page kept as an endpoint to report on", len(endpoints))
	}
	if got := prepare(endpoints[0]); !strings.Contains(got.Blocked, "is absent") {
		t.Errorf("Blocked = %q, want it to say the page is absent from the cache", got.Blocked)
	}
}
