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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/shing1211/webullapi4go/client"
)

// The mutating gate is a safety property, so every case here asserts on the
// request side: an httptest server records what reached the wire, and the
// assertion is on that list rather than on a return value. A test that only
// checked the summary's numbers would pass against an implementation that sent
// the request and then wrote "skipped" over the result, which is the failure mode
// this gate exists to make impossible.

// probeMutatingEndpoints is the endpoint set both gate cases walk: one POST in
// each of the three mutating areas, plus one market-data GET that must be called
// either way.
//
// The bodies matter. trade.PlaceOrder carries a required array and the other
// three carry only scalars, so the set spans both halves of the argument: the
// array rule cannot help the three scalar ones, and the gate is the only barrier
// in front of them. That is the case the empty-array rule structurally cannot
// cover, and it is why the gate exists at all.
func probeMutatingEndpoints() []Endpoint {
	return []Endpoint{
		{
			Symbol: "trade.PlaceOrder", Fixture: "trading/POST-trading-orders-place.json",
			Method: "POST", Path: "/trading/orders/place", HasBody: true,
			BodyParams: []NamedParam{
				{Name: "new_orders", Spec: ParamSpec{Type: "array"}},
			},
		},
		{
			Symbol: "trade.CancelOrder", Fixture: "trading/POST-trading-orders-cancel.json",
			Method: "POST", Path: "/trading/orders/cancel", HasBody: true,
			BodyParams: []NamedParam{
				{Name: "order_id", Spec: ParamSpec{Type: "string", Example: "1", HasExample: true}},
			},
		},
		{
			Symbol:  "broker.CreateVirtualAccount",
			Fixture: "broker-hk/POST-broker-accounts-virtual-accounts-create.json",
			Method:  "POST", Path: "/broker/accounts/virtual-accounts/create", HasBody: true,
			BodyParams: []NamedParam{
				{Name: "label", Spec: ParamSpec{Type: "string", Example: "probe", HasExample: true}},
			},
		},
		{
			Symbol:  "brokerFD.CreateACHRelationship",
			Fixture: "broker-fd-us/POST-broker-funding-ach-relationships-create.json",
			Method:  "POST", Path: "/broker/funding/ach/relationships/create", HasBody: true,
			BodyParams: []NamedParam{
				{Name: "processor_token", Spec: ParamSpec{Type: "string", Example: "probe", HasExample: true}},
			},
		},
		{
			Symbol: "data.GetQuote", Fixture: "market-data-stock/GET-market-data-stocks-quote.json",
			Method: "GET", Path: "/market-data/stocks/quotes",
		},
	}
}

// recordingServer starts a server that records "METHOD path" for every request
// and answers 200 to all of them, and returns the recorded list.
func recordingServer(t *testing.T) (*httptest.Server, func() []string) {
	t.Helper()
	var mu sync.Mutex
	var seen []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		seen = append(seen, r.Method+" "+r.URL.Path)
		mu.Unlock()
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(srv.Close)
	return srv, func() []string {
		mu.Lock()
		defer mu.Unlock()
		out := make([]string, len(seen))
		copy(out, seen)
		sort.Strings(out)
		return out
	}
}

func outcomesBySymbol(outcomes []Outcome) map[string]Outcome {
	by := map[string]Outcome{}
	for _, o := range outcomes {
		by[o.Symbol] = o
	}
	return by
}

// The load-bearing case. With the opt-in absent nothing that can mutate may reach
// the network, and the market-data GET must still be called - otherwise a gate
// that simply muted the whole census would pass a test that only checks that
// nothing dangerous happened.
//
// The client presents a sandbox host, so this case isolates ONE of the gate's two
// conditions. That is deliberate: it is what makes the absence of the opt-in the
// only possible cause of the refusal, and the reason is asserted to name it.
func TestCensusDoesNotCallMutatingEndpointsWithoutTheOptIn(t *testing.T) {
	// Set to empty rather than left alone, so the case holds on a machine whose
	// shell exports the variable for the mutating sandbox tests.
	t.Setenv(MutateOptInEnv, "")
	t.Setenv(NonSandboxOverrideEnv, "")

	srv, recorded := recordingServer(t)
	cl := newProbeTestClientPresenting(t, probeSandboxBaseURL, srv)
	outcomes, err := Census(context.Background(), cl, probeMutatingEndpoints())
	if err != nil {
		t.Fatalf("Census: %v", err)
	}

	// The assertion is on the wire, not on a return value.
	if got, want := recorded(), []string{"GET /market-data/stocks/quotes"}; !equalStrings(got, want) {
		t.Errorf("requests that reached the server = %v, want exactly %v.\n"+
			"Anything else means the census called an endpoint that can change state "+
			"with the gate closed.", got, want)
	}

	by := outcomesBySymbol(outcomes)
	for _, symbol := range []string{
		"trade.PlaceOrder", "trade.CancelOrder",
		"broker.CreateVirtualAccount", "brokerFD.CreateACHRelationship",
	} {
		o, ok := by[symbol]
		if !ok {
			t.Fatalf("no outcome for %s", symbol)
		}
		if o.Skipped == "" {
			t.Errorf("%s: Skipped = %q, want the gate's reason", symbol, o.Skipped)
		}
		if o.Skipped != skippedNotOptedIn() {
			t.Errorf("%s: Skipped = %q, want the not-opted-in reason %q. A run against a "+
				"sandbox host failed only the opt-in, so any other reason would be this "+
				"gate misreporting which condition it refused on.", symbol, o.Skipped,
				skippedNotOptedIn())
		}
		if o.Status != 0 {
			t.Errorf("%s: Status = %d, want 0; a skipped endpoint has no status",
				symbol, o.Status)
		}
		if o.Err != "" {
			t.Errorf("%s: Err = %q, want empty; no request was sent, so no error "+
				"code describes the server", symbol, o.Err)
		}
		if o.Blocked != "" {
			t.Errorf("%s: Blocked = %q, want empty; the request was buildable and was "+
				"refused by policy, which is a different claim", symbol, o.Blocked)
		}
		if !o.Mutating {
			t.Errorf("%s: Mutating = false, want true", symbol)
		}
	}

	// The non-mutating endpoint is called, and it is called as a call.
	if o := by["data.GetQuote"]; o.Status != http.StatusOK || o.Skipped != "" || o.Mutating {
		t.Errorf("data.GetQuote = %+v, want a 200 and no gate", o)
	}

	summary := Summarise(outcomes)
	if summary.Skipped != 4 {
		t.Errorf("Skipped = %d, want 4", summary.Skipped)
	}
	if summary.Reachable != 1 {
		t.Errorf("Reachable = %d, want 1; a refused endpoint must never be counted as "+
			"reachable", summary.Reachable)
	}
	if summary.Blocked != 0 || summary.Unanswered != 0 {
		t.Errorf("Blocked = %d, Unanswered = %d, want 0 and 0; 'we did not call it' is "+
			"neither a build failure nor a transport failure", summary.Blocked, summary.Unanswered)
	}
	if summary.Mutating != 4 {
		t.Errorf("Mutating = %d, want 4", summary.Mutating)
	}
	if sum := summary.Reachable + summary.Blocked + summary.Skipped + summary.Unanswered; sum != summary.Total {
		t.Errorf("Reachable+Blocked+Skipped+Unanswered = %d, want Total %d; the four classes "+
			"must partition the corpus", sum, summary.Total)
	}
	// The status map is the other way a refused endpoint could be reported as a
	// server behaviour. It must hold the one status that was really observed.
	if len(summary.Statuses) != 1 || summary.Statuses["200"] != 1 {
		t.Errorf("Statuses = %v, want only 200=1; a skipped row contributes no status",
			summary.Statuses)
	}
	if area := summary.Area["trading"]; area.Skipped != 2 || area.Reachable != 0 || area.Mutating != 2 {
		t.Errorf("trading = %+v, want 2 skipped, 0 reachable, 2 mutating", area)
	}
	if area := summary.Area["market-data-stock"]; area.Skipped != 0 || area.Reachable != 1 {
		t.Errorf("market-data-stock = %+v, want 1 reachable and 0 skipped", area)
	}
}

// The other half of the gate: with BOTH conditions satisfied - the opt-in set
// and the client presenting the sandbox host - the same set is called.
//
// This is the path a live run takes, so it is proven rather than assumed. Without
// it the gate could be an unconditional refusal that satisfies every "nothing
// dangerous happened" assertion in the file, and the open path could exist only
// behind the override, which is not what a sandbox run uses.
func TestCensusCallsMutatingEndpointsWithTheOptInAndASandboxHost(t *testing.T) {
	t.Setenv(MutateOptInEnv, MutateOptInValue)
	t.Setenv(NonSandboxOverrideEnv, "")

	srv, recorded := recordingServer(t)
	cl := newProbeTestClientPresenting(t, probeSandboxBaseURL, srv)
	outcomes, err := Census(context.Background(), cl, probeMutatingEndpoints())
	if err != nil {
		t.Fatalf("Census: %v", err)
	}

	want := []string{
		"GET /market-data/stocks/quotes",
		"POST /broker/accounts/virtual-accounts/create",
		"POST /broker/funding/ach/relationships/create",
		"POST /trading/orders/cancel",
		"POST /trading/orders/place",
	}
	if got := recorded(); !equalStrings(got, want) {
		t.Errorf("requests that reached the server = %v, want %v; with %s=%s and a sandbox "+
			"host the gate must not block", got, want, MutateOptInEnv, MutateOptInValue)
	}

	by := outcomesBySymbol(outcomes)
	for _, symbol := range []string{
		"trade.PlaceOrder", "trade.CancelOrder",
		"broker.CreateVirtualAccount", "brokerFD.CreateACHRelationship",
	} {
		o := by[symbol]
		if o.Skipped != "" {
			t.Errorf("%s: Skipped = %q, want empty when the gate is open", symbol, o.Skipped)
		}
		if o.Status != http.StatusOK {
			t.Errorf("%s: Status = %d, want 200", symbol, o.Status)
		}
		if !o.Mutating {
			t.Errorf("%s: Mutating = false; the row must still say which class it is in", symbol)
		}
	}

	// With the gate open, the scalar body really is sent. That is what makes this
	// case a live-mutation test rather than a second refusal test, and it is why
	// it is pointed at httptest and nowhere else.
	if o := by["trade.CancelOrder"]; o.BodyKeys == nil || o.BodyKeys[0] != "order_id" {
		t.Errorf("trade.CancelOrder BodyKeys = %v, want [order_id]", o.BodyKeys)
	}

	summary := Summarise(outcomes)
	if summary.Reachable != 5 || summary.Skipped != 0 || summary.Mutating != 4 {
		t.Errorf("Reachable = %d, Skipped = %d, Mutating = %d, want 5, 0, 4",
			summary.Reachable, summary.Skipped, summary.Mutating)
	}
}

// The second condition, and the hazard it removes. The opt-in alone is not
// enough: a shell set up to run the mutating trade tests already exports
// WEBULL_TRADE_MUTATE=1, and the base URL comes from the environment too, so
// "export the variable, point the SDK at production, run the census" would
// otherwise issue the endpoints below - three of them scalar-bodied, and
// therefore made inert by no other rule - against a real account.
//
// Every case is asserted on the wire rather than on a return value, because a gate
// that sent the request and then wrote "skipped" over it would satisfy every
// count-based assertion in this file.
func TestCensusRequiresBothTheOptInAndASandboxHost(t *testing.T) {
	const productionBaseURL = "https://api.webull.hk"
	mutatingPaths := []string{
		"POST /broker/accounts/virtual-accounts/create",
		"POST /broker/funding/ach/relationships/create",
		"POST /trading/orders/cancel",
		"POST /trading/orders/place",
	}
	// The market-data GET is expected in every case, so a gate that silenced the
	// whole census fails here rather than passing.
	const readRequest = "GET /market-data/stocks/quotes"

	cases := []struct {
		name string
		// optedIn and overridden are the two environment variables as an operator
		// would set them; presented is the base URL the client's host - and
		// therefore the gate's - is read from.
		optedIn    string
		overridden string
		presented  string
		wantCalled bool
		// wantReason is compared exactly, so a refusal can never be recorded under
		// the wrong condition's reason.
		wantReason string
	}{
		{
			// This is the live run's path, and the reason it is a case at all.
			name: "opt-in and a sandbox host: called", optedIn: MutateOptInValue,
			presented: probeSandboxBaseURL, wantCalled: true,
		},
		{
			name:       "opt-in and a production host: refused for the host",
			optedIn:    MutateOptInValue,
			presented:  productionBaseURL,
			wantReason: skippedNotSandbox("api.webull.hk"),
		},
		{
			// The escape hatch has to be a working door, not a branch nothing
			// reaches. An override that never opened anything would be worse than
			// no override: it would look like the way through.
			name:       "opt-in, override and a production host: called",
			optedIn:    MutateOptInValue,
			overridden: NonSandboxOverrideValue,
			presented:  productionBaseURL, wantCalled: true,
		},
		{
			// Independence, the other direction: the override waives the host
			// condition and nothing else.
			name:       "override and a sandbox host, no opt-in: refused for the opt-in",
			overridden: NonSandboxOverrideValue, presented: probeSandboxBaseURL,
			wantReason: skippedNotOptedIn(),
		},
		{
			// Both unmet, and the opt-in is the one reported: a run that was never
			// authorised is the more basic fact, and reporting the host first
			// would point an operator at the wrong variable.
			name:       "override and a production host, no opt-in: refused for the opt-in",
			overridden: NonSandboxOverrideValue, presented: productionBaseURL,
			wantReason: skippedNotOptedIn(),
		},
		{
			name:      "neither variable set, a sandbox host: refused for the opt-in",
			presented: probeSandboxBaseURL, wantReason: skippedNotOptedIn(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(MutateOptInEnv, tc.optedIn)
			t.Setenv(NonSandboxOverrideEnv, tc.overridden)

			srv, recorded := recordingServer(t)
			cl := newProbeTestClientPresenting(t, tc.presented, srv)
			outcomes, err := Census(context.Background(), cl, probeMutatingEndpoints())
			if err != nil {
				t.Fatalf("Census: %v", err)
			}

			want := []string{readRequest}
			if tc.wantCalled {
				want = append(want, mutatingPaths...)
			}
			if got := recorded(); !equalStrings(got, want) {
				t.Errorf("requests that reached the server = %v, want %v", got, want)
			}

			by := outcomesBySymbol(outcomes)
			for _, symbol := range []string{
				"trade.PlaceOrder", "trade.CancelOrder",
				"broker.CreateVirtualAccount", "brokerFD.CreateACHRelationship",
			} {
				o := by[symbol]
				if tc.wantCalled {
					if o.Skipped != "" || o.Status != http.StatusOK {
						t.Errorf("%s: Skipped = %q, Status = %d; want an open gate and a 200",
							symbol, o.Skipped, o.Status)
					}
					continue
				}
				if o.Skipped != tc.wantReason {
					t.Errorf("%s: Skipped = %q, want %q", symbol, o.Skipped, tc.wantReason)
				}
				if o.Status != 0 {
					t.Errorf("%s: Status = %d, want 0; a refused endpoint has no status",
						symbol, o.Status)
				}
			}

			summary := Summarise(outcomes)
			if tc.wantCalled {
				if summary.Skipped != 0 || summary.Reachable != 5 {
					t.Errorf("Skipped = %d, Reachable = %d, want 0 and 5",
						summary.Skipped, summary.Reachable)
				}
				return
			}
			// A refusal must never be counted as an answer, which is the
			// arithmetic the summary exists to keep honest.
			if summary.Skipped != 4 || summary.Reachable != 1 {
				t.Errorf("Skipped = %d, Reachable = %d, want 4 and 1",
					summary.Skipped, summary.Reachable)
			}
			if summary.Mutating != 4 {
				t.Errorf("Mutating = %d, want 4; the class is recorded whether or not "+
					"the gate stopped it", summary.Mutating)
			}
		})
	}
}

// The two reasons must be tellable apart, because they call for different actions
// and reporting one as the other sends an operator to change the wrong thing - or,
// worse, lets a run pointed at production be read as merely un-authorised.
//
// So this asserts on the CONTENT of both strings and not only on the two states a
// walk can reach: the opt-in reason must not mention the host or the override,
// and the host reason must not mention the opt-in variable.
func TestSkipReasonsAreDistinctPerFailingCondition(t *testing.T) {
	const host = "api.webull.hk"
	notOptedIn := skippedNotOptedIn()
	notSandbox := skippedNotSandbox(host)

	if notOptedIn == notSandbox {
		t.Fatalf("both conditions produce the same reason %q; a reader could not tell "+
			"an unauthorised run from a run pointed at production", notOptedIn)
	}
	for _, want := range []string{MutateOptInEnv, MutateOptInValue, "not opted in"} {
		if !strings.Contains(notOptedIn, want) {
			t.Errorf("the not-opted-in reason %q does not contain %q", notOptedIn, want)
		}
	}
	for _, unwanted := range []string{NonSandboxOverrideEnv, host, "sandbox host"} {
		if strings.Contains(notOptedIn, unwanted) {
			t.Errorf("the not-opted-in reason %q mentions %q; that run failed only the "+
				"opt-in and must not name the condition that was not the problem",
				notOptedIn, unwanted)
		}
	}
	for _, want := range []string{host, NonSandboxOverrideEnv, "not a sandbox host"} {
		if !strings.Contains(notSandbox, want) {
			t.Errorf("the not-a-sandbox reason %q does not contain %q", notSandbox, want)
		}
	}
	for _, unwanted := range []string{MutateOptInEnv, "not opted in"} {
		if strings.Contains(notSandbox, unwanted) {
			t.Errorf("the not-a-sandbox reason %q mentions %q; the opt-in WAS satisfied, and "+
				"saying otherwise would send an operator to set a variable already set",
				notSandbox, unwanted)
		}
	}
	// An empty host is the only case where the reason cannot name a deployment, and
	// it must still read as a claim about the host rather than as a blank.
	if reason := skippedNotSandbox("  "); !strings.Contains(reason, "the configured base URL") {
		t.Errorf("the not-a-sandbox reason for an unknown host is %q, want it to say the host "+
			"could not be identified", reason)
	}
}

// The gate's arithmetic, stated once as a table so it cannot drift from the walk
// cases above. The override widens the host condition; it does not touch the
// opt-in, which is what makes the two conditions independent.
func TestMutationGateOpenNeedsTheOptInAndOneOfTwoHostPaths(t *testing.T) {
	cases := []struct {
		optedIn    bool
		sandbox    bool
		overridden bool
		wantOpen   bool
	}{
		{false, false, false, false},
		{false, true, false, false},
		{false, true, true, false},
		{false, false, true, false},
		{true, true, false, true},
		{true, false, false, false},
		{true, false, true, true},
		{true, true, true, true},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("optedIn=%t sandbox=%t overridden=%t",
			tc.optedIn, tc.sandbox, tc.overridden), func(t *testing.T) {
			gate := MutationGate{
				OptedIn:     tc.optedIn,
				SandboxHost: tc.sandbox,
				Overridden:  tc.overridden,
				Host:        "api.webull.hk",
			}
			if got := gate.Open(); got != tc.wantOpen {
				t.Errorf("Open() = %t, want %t", got, tc.wantOpen)
			}
			// An open gate has no reason to record and a closed one always does: a
			// refused row that says nothing is the hole this gate exists to avoid.
			if reason := gate.SkipReason(); (reason != "") == tc.wantOpen {
				t.Errorf("SkipReason() = %q with Open() = %t; the two must always disagree",
					reason, tc.wantOpen)
			}
		})
	}
}

// ResolveMutationGate is the single reading of the environment per run, and the
// override is held to the same exact-value discipline as the opt-in: a "true" that
// opened one and not the other would give the two variables different weights.
func TestResolveMutationGateReadsBothVariables(t *testing.T) {
	const sandboxHost = "api.sandbox.webull.hk"
	cases := []struct {
		optIn      string
		overridden string
		host       string
		want       MutationGate
	}{
		{"", "", sandboxHost, MutationGate{SandboxHost: true}},
		{MutateOptInValue, "", sandboxHost, MutationGate{OptedIn: true, SandboxHost: true}},
		{MutateOptInValue, "", "api.webull.hk", MutationGate{OptedIn: true}},
		{MutateOptInValue, NonSandboxOverrideValue, "api.webull.hk",
			MutationGate{OptedIn: true, Overridden: true}},
		{"true", "yes", sandboxHost, MutationGate{SandboxHost: true}},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("optIn=%q override=%q host=%q", tc.optIn, tc.overridden, tc.host),
			func(t *testing.T) {
				t.Setenv(MutateOptInEnv, tc.optIn)
				t.Setenv(NonSandboxOverrideEnv, tc.overridden)
				got := ResolveMutationGate(tc.host)
				want := tc.want
				want.Host = tc.host
				if got != want {
					t.Errorf("ResolveMutationGate(%q) = %+v, want %+v", tc.host, got, want)
				}
			})
	}
}

// The opt-in reads one value, exactly. Every other reading - any non-empty value,
// "true", a padded "1" - would give one variable two meanings, which is the
// collision the shared name creates and the reason these are pinned here. The
// override is pinned the same way, in the same table shape, because the two are
// read by the same comparison and must not drift apart in what they accept.
func TestMutateOptInAcceptsOnlyTheTradeTestsValue(t *testing.T) {
	cases := []struct {
		value string
		want  bool
	}{
		{"1", true},
		{"", false},
		{"0", false},
		{"true", false},
		{"TRUE", false},
		{"yes", false},
		{" 1", false},
		{"1 ", false},
		{"11", false},
	}
	for _, tc := range cases {
		t.Run("value "+fmt.Sprintf("%q", tc.value), func(t *testing.T) {
			t.Setenv(MutateOptInEnv, tc.value)
			if got := MutationOptedIn(); got != tc.want {
				t.Errorf("MutationOptedIn() with %s=%q = %t, want %t",
					MutateOptInEnv, tc.value, got, tc.want)
			}
		})
	}
}

func TestNonSandboxOverrideAcceptsOnlyTheSameSingleValue(t *testing.T) {
	cases := []struct {
		value string
		want  bool
	}{
		{"1", true},
		{"", false},
		{"0", false},
		{"true", false},
		{"yes", false},
		{" 1", false},
		{"11", false},
	}
	for _, tc := range cases {
		t.Run("value "+fmt.Sprintf("%q", tc.value), func(t *testing.T) {
			t.Setenv(NonSandboxOverrideEnv, tc.value)
			if got := NonSandboxOverridden(); got != tc.want {
				t.Errorf("NonSandboxOverridden() with %s=%q = %t, want %t",
					NonSandboxOverrideEnv, tc.value, got, tc.want)
			}
		})
	}
}

// Which hosts open the gate. The refusals are the interesting half of this table:
// it is the set that stands between an exported WEBULL_TRADE_MUTATE=1 and a
// production account, and a lookalike that slipped in would be a silent reopening
// of exactly the hazard the second condition exists to close.
func TestIsSandboxHostAcceptsOnlyHostsTheSDKDerivesForASandbox(t *testing.T) {
	cases := []struct {
		name string
		host string
		want bool
	}{
		// The host AGENTS.md names as the one sandbox host that may appear in
		// committed material, which is also the host a live HK run uses.
		{"the HK sandbox", "api.sandbox.webull.hk", true},
		{"the HK broker sandbox", "broker-api.sandbox.webull.hk", true},
		{"another region's sandbox", "api.sandbox.webull.com", true},
		{"a base URL rather than a bare host", "https://api.sandbox.webull.hk", true},
		{"normalisation is not a second identity", "  API.SANDBOX.WEBULL.HK  ", true},
		{"an explicit port is the same deployment", "api.sandbox.webull.hk:8443", true},
		{"a fully qualified name is the same deployment", "api.sandbox.webull.hk.", true},
		// Every production host the SDK derives must be refused, in every
		// environment, because production is what the gate exists to keep out.
		{"HK production", "api.webull.hk", false},
		{"HK broker production", "broker-api.webull.hk", false},
		{"US production", "api.webull.com", false},
		{"a lookalike that merely contains the label", "sandbox.api.webull.hk", false},
		{"a lookalike that appends the label", "api.webull.hk.sandbox", false},
		{"a lookalike domain", "api.sandbox.webull.hk.example.com", false},
		{"a lookalike host under a lookalike label", "sandbox.attacker.example", false},
		{"a regional events host, which the census never dials", "events-api.sandbox.webull.hk", false},
		{"a local test server", "127.0.0.1:52341", false},
		{"the loopback name", "localhost", false},
		{"nothing at all", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsSandboxHost(tc.host); got != tc.want {
				t.Errorf("IsSandboxHost(%q) = %t, want %t", tc.host, got, tc.want)
			}
		})
	}
}

// The recognised set is derived, not written, and this asserts the two things a
// hand-written list could get wrong: that every entry is a region the SDK really
// has (so a typo cannot silently shrink the set), and that nothing derived for a
// PRODUCTION environment is in it (so the derivation cannot be pointed at the
// wrong environment by a later edit).
func TestKnownRegionsAreAllValidAndYieldNoProductionHost(t *testing.T) {
	if len(knownRegions) == 0 {
		t.Fatal("knownRegions is empty, so no host is a sandbox host and the gate is " +
			"permanently closed for every base URL")
	}
	seen := map[client.Region]bool{}
	for _, region := range knownRegions {
		if !region.IsValid() {
			t.Errorf("%q is not a region the SDK knows, so it silently contributes no "+
				"sandbox host; EndpointsFor falls back to Hong Kong for it", region)
		}
		if seen[region] {
			t.Errorf("%q is listed twice", region)
		}
		seen[region] = true

		sandbox := client.EndpointsFor(region, client.Sandbox)
		production := client.EndpointsFor(region, client.Production)
		for _, host := range []string{sandbox.HTTP, sandbox.BrokerHTTP} {
			if !strings.Contains(host, "sandbox.") {
				t.Errorf("%s sandbox host %q carries no sandbox label; a sandbox host "+
					"that does not look like one is a derivation bug, not a feature",
					region, host)
			}
			if host == production.HTTP || host == production.BrokerHTTP {
				t.Errorf("%s sandbox host %q is also its production host; the gate would "+
					"open on production", region, host)
			}
		}
	}
	// The host the report will actually name for an HK sandbox run, asserted
	// through the same derivation the gate uses rather than as a literal, so a
	// change in the SDK's endpoint derivation fails here and not silently in a
	// live run.
	if got := hostnameOf(client.EndpointsFor(client.HK, client.Sandbox).HTTP); got != "api.sandbox.webull.hk" {
		t.Errorf("the HK sandbox REST host resolves to %q, want api.sandbox.webull.hk; "+
			"AGENTS.md and the tests both name that host", got)
	}
}

// The classifier is the safety property, so it is pinned as a table including the
// cases it deliberately does not catch. The three "false on a POST" rows are
// documented limits, not oversights: they are what IsMutating's own comment says
// it cannot see, and a reader who finds one of them called should find it
// explained here rather than merely unhandled.
func TestIsMutatingClassifiesFromMethodAndArea(t *testing.T) {
	cases := []struct {
		name     string
		fixture  string
		method   string
		mutating bool
	}{
		{"trading POST places an order", "trading/POST-x.json", "POST", true},
		{"broker-hk POST creates an account", "broker-hk/POST-x.json", "POST", true},
		{"broker-fd-us POST creates a bank relationship", "broker-fd-us/POST-x.json", "POST", true},
		{"trading GET only reads", "trading/GET-x.json", "GET", false},
		{"broker-fd-us GET only reads", "broker-fd-us/GET-x.json", "GET", false},
		{"an unknown method fails closed", "trading/x.json", "", true},
		{"a lower-case method is still unsafe", "trading/POST-x.json", "post", true},
		// The three documented limits. See IsMutating.
		{"a watchlist POST is outside the class", "market-data-watchlist/POST-x.json", "POST", false},
		{"a display-solution POST is outside the class", "display-solution/POST-x.json", "POST", false},
		{"a token exchange POST is outside the class", "authentication/POST-x.json", "POST", false},
		{"an area this corpus has never had is outside it", "unreleased-area/POST-x.json", "POST", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsMutating(Endpoint{Fixture: tc.fixture, Method: tc.method})
			if got != tc.mutating {
				t.Errorf("IsMutating(%s %s) = %t, want %t", tc.method, tc.fixture, got, tc.mutating)
			}
		})
	}
}

// The classifier is checked against the committed corpus, so the number a reader
// is given before a live run - 34 of 193 - is a measured one and not a claim. It
// also pins the four endpoints named in the brief as the reason the gate exists:
// each has a body of scalars, so no value rule could have made them inert.
func TestMutatingClassAgainstTheCommittedManifest(t *testing.T) {
	manifest, err := resolveManifest("")
	if err != nil {
		t.Skipf("conformance/testdata/manifest.json is not reachable from here: %v", err)
	}
	raw, err := os.ReadFile(manifest) //nolint:gosec // G304: the path came from resolveManifest, which locates the committed manifest.
	if err != nil {
		t.Fatalf("reading the manifest: %v", err)
	}
	var parsed struct {
		Fixtures []manifestRecord `json:"fixtures"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatalf("parsing the manifest: %v", err)
	}

	byArea := map[string]int{}
	paths := map[string]bool{}
	for _, record := range parsed.Fixtures {
		fixture := record.Fixture
		if fixture == "" {
			fixture = record.ID + ".json"
		}
		ep := Endpoint{Fixture: fixture, Method: strings.ToUpper(record.Documented.Method)}
		if !IsMutating(ep) {
			continue
		}
		byArea[areaOf(fixture)]++
		paths[record.Documented.Path] = true
	}

	wantAreas := map[string]int{"trading": 5, "broker-hk": 11, "broker-fd-us": 18}
	if len(byArea) != len(wantAreas) {
		t.Errorf("mutating areas = %v, want exactly %v", byArea, wantAreas)
	}
	total := 0
	for area, want := range wantAreas {
		if byArea[area] != want {
			t.Errorf("%s contributes %d mutating endpoints, want %d", area, byArea[area], want)
		}
		total += want
	}
	if total != 34 {
		t.Errorf("mutating total = %d, want 34", total)
	}
	if len(parsed.Fixtures) != 193 {
		t.Errorf("the manifest holds %d fixtures; the percentages in the report are "+
			"quoted against 193 and both numbers move together", len(parsed.Fixtures))
	}

	// The scalar-body cases, by path. These are the endpoints no array rule can
	// protect, so the gate is the only thing standing in front of them. The paths
	// are the manifest's own; note that the ACH relationship endpoints are
	// /broker/funding/ach-relationships/... with a hyphen, not the
	// ach/relationships spelling a reader might assume from the broker area's
	// /broker/funding/bank-relationships/... neighbour.
	for _, path := range []string{
		"/trading/orders/cancel",
		"/broker/accounts/virtual-accounts/create",
		"/broker/accounts/virtual-accounts/update",
		"/broker/funding/ach-relationships/create",
	} {
		if !paths[path] {
			t.Errorf("%s is not in the mutating class; the gate would not protect it", path)
		}
	}
	// A read in a mutating area must stay out of it, or the census would refuse
	// 164 endpoints it could safely ask.
	if paths["/trading/accounts/list"] {
		t.Error("/trading/accounts/list is a GET and must not be in the mutating class")
	}
}

// The three-way model, at the level a reader of census.json looks at. A skipped
// row is a fourth thing, and the test that matters is the one that would fail if
// it were folded into any of the other three.
func TestSummariseSeparatesSkippedFromEveryOtherClass(t *testing.T) {
	outcomes := []Outcome{
		{Symbol: "a", Area: "trading", Status: 200, Mutating: true},
		{Symbol: "b", Area: "trading", Status: 404},
		{Symbol: "c", Area: "trading", Skipped: skippedNotOptedIn(), Mutating: true},
		{Symbol: "d", Area: "display-solution", Blocked: "body property bank_code"},
		{Symbol: "e", Area: "fundamentals", Err: "transport"},
	}
	got := Summarise(outcomes)
	if got.Total != 5 {
		t.Errorf("Total = %d, want 5", got.Total)
	}
	if got.Reachable != 2 {
		t.Errorf("Reachable = %d, want 2; a skipped row is not an answer", got.Reachable)
	}
	if got.Skipped != 1 {
		t.Errorf("Skipped = %d, want 1", got.Skipped)
	}
	if got.Blocked != 1 {
		t.Errorf("Blocked = %d, want 1; a skipped row is not a build failure", got.Blocked)
	}
	if got.Unanswered != 1 {
		t.Errorf("Unanswered = %d, want 1; a skipped row is not a transport failure", got.Unanswered)
	}
	if got.Mutating != 2 {
		t.Errorf("Mutating = %d, want 2", got.Mutating)
	}
	if got.BlockedReasons["body property bank_code"] != 1 {
		t.Errorf("BlockedReasons = %v, want the build failure only; a skipped row has no "+
			"blocked reason to attribute", got.BlockedReasons)
	}
	if got.Area["trading"].Skipped != 1 || got.Area["trading"].Reachable != 2 {
		t.Errorf("trading = %+v, want 1 skipped and 2 reachable", got.Area["trading"])
	}
}

// The gate is a class of its own in the diagnostics, and it must not be able to
// hide inside a class that means a probe bug.
func TestDiagnoseFilesTheMutatingClassSeparately(t *testing.T) {
	row := Outcome{
		Fixture: "trading/POST-trading-orders-cancel.json", Area: "trading",
		Method: "POST", Mutating: true, Skipped: skippedNotOptedIn(),
	}
	got := diagnose([]Outcome{row})
	if len(got) != 1 {
		t.Fatalf("diagnose returned %d classes, want exactly 1: %v", len(got), classNames(got))
	}
	if _, ok := got["endpoint can mutate, gated by "+MutateOptInEnv]; !ok {
		t.Errorf("diagnose classes = %v, want the mutating class", classNames(got))
	}
	// The same row in a dry run: nothing was sent, so nothing was skipped, and the
	// class must still be reported, because it is the number a reader needs before
	// a live run.
	dry := Outcome{Fixture: row.Fixture, Area: row.Area, Method: row.Method, Mutating: true, DryRun: true}
	if got := diagnose([]Outcome{dry}); len(got) != 1 {
		t.Errorf("diagnose on a dry run returned %v, want the mutating class alone", classNames(got))
	}
	// A row that is neither gated nor blocked must produce no class at all, so the
	// mutating class cannot become a catch-all.
	plain := Outcome{Fixture: "fundamentals/GET-x.json", Area: "fundamentals", Status: 200}
	if got := diagnose([]Outcome{plain}); len(got) != 0 {
		t.Errorf("diagnose on a plain 200 = %v, want no classes", classNames(got))
	}
}

// The counts are half the requirement; the printed table is the other half. A
// reader who sees only the table must be able to tell a refusal from an answer,
// which means the word "skipped" is on the page and the reachable column does not
// include the refused rows.
func TestPrintSummaryShowsSkippedAsItsOwnClass(t *testing.T) {
	built := report{
		Mode:     "live",
		Host:     "api.sandbox.webull.hk",
		Account:  "succeeded: one account id is threaded into every account-scoped request",
		Mutation: mutationNote("live", MutationGate{Host: "api.webull.hk"}),
		Outcomes: []Outcome{
			{Symbol: "trade.CancelOrder", Area: "trading", Method: "POST",
				Path: "/trading/orders/cancel", Mutating: true, Skipped: skippedNotOptedIn()},
			{Symbol: "data.GetQuote", Area: "market-data-stock", Method: "GET",
				Path: "/market-data/stocks/quotes", Status: http.StatusOK},
		},
	}
	built.Summary = Summarise(built.Outcomes)

	out := captureStdout(t, func() { printSummary(built) })

	for _, want := range []string{
		"mutating gate: closed",
		MutateOptInEnv + "=" + MutateOptInValue + " is not set",
		"reachable 1",
		"blocked 0",
		"skipped 1",
		"mutating class 1",
		"skipped: mutating",
		"trading  ",
		// The footer has to tell a reader how to get past the gate, and there are
		// now two things to get past. Printing only the opt-in would send an
		// operator pointed at a sandbox to set a variable they already set.
		NonSandboxOverrideEnv,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("printed table does not contain %q\n---\n%s", want, out)
		}
	}
	// The per-area row must carry the refusal in the skipped column and nothing in
	// the reachable one, so a reader scanning the table cannot mistake it for a
	// row that answered.
	var tradingLine string
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "trading ") {
			tradingLine = line
		}
	}
	fields := strings.Fields(tradingLine)
	if len(fields) < 6 {
		t.Fatalf("trading row = %q, want total, reachable, blocked, skipped, mutating and statuses", tradingLine)
	}
	if fields[1] != "1" || fields[2] != "0" || fields[3] != "0" || fields[4] != "1" || fields[5] != "1" {
		t.Errorf("trading row = %q, want total 1, reachable 0, blocked 0, skipped 1, mutating 1", tradingLine)
	}
}

// The gate state travels with the artifact. A census.json is designed to be
// committed and diffed, so a reader months later must be able to tell whether the
// 34 endpoints were called without knowing anything about the machine that wrote
// it - and which of the two conditions was the one that stopped them.
//
// The four live cases are all distinct, because collapsing any pair of them is the
// misreading this field exists to prevent: a sandbox run and a production run
// differ, and an unauthorised run differs from both.
func TestReportRecordsTheMutationGateState(t *testing.T) {
	const host = "api.sandbox.webull.hk"
	cases := []struct {
		name string
		mode string
		gate MutationGate
		want string
	}{
		{
			name: "live, not opted in", mode: "live",
			gate: MutationGate{Host: host},
			want: "closed: " + MutateOptInEnv,
		},
		{
			name: "live, opted in against a non-sandbox host", mode: "live",
			gate: MutationGate{OptedIn: true, Host: "api.webull.hk"},
			want: "closed: " + MutateOptInEnv + "=" + MutateOptInValue + " was set but " +
				"api.webull.hk is not a sandbox host",
		},
		{
			name: "live, opted in against a sandbox host", mode: "live",
			gate: MutationGate{OptedIn: true, SandboxHost: true, Host: host},
			want: "OPEN: " + MutateOptInEnv + "=" + MutateOptInValue,
		},
		{
			name: "live, opted in and overridden against a non-sandbox host", mode: "live",
			gate: MutationGate{OptedIn: true, Overridden: true, Host: "api.webull.hk"},
			want: "OPEN BY OVERRIDE: " + NonSandboxOverrideEnv,
		},
		{
			name: "dry run", mode: "dry-run: no request was sent",
			gate: MutationGate{},
			want: "not applicable",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			note := mutationNote(tc.mode, tc.gate)
			if !strings.HasPrefix(note, tc.want) {
				t.Errorf("mutationNote = %q, want it to start with %q", note, tc.want)
			}
			// The two open cases must say which one they are: "OPEN" for a sandbox
			// run and "OPEN BY OVERRIDE" for a run against a host that is not one.
			// A reader who finds orders placed against production must be able to
			// see that from the artifact alone.
			if tc.gate.Overridden && strings.HasPrefix(note, "OPEN:") {
				t.Errorf("mutationNote = %q, an override run must not read as a plain OPEN", note)
			}
			encoded, err := json.Marshal(report{Mode: tc.mode, Mutation: note})
			if err != nil {
				t.Fatalf("marshalling: %v", err)
			}
			var decoded map[string]any
			if err := json.Unmarshal(encoded, &decoded); err != nil {
				t.Fatalf("unmarshalling: %v", err)
			}
			// No omitempty: an absent field would be indistinguishable from a
			// report written before the gate existed.
			if decoded["mutationGate"] != note {
				t.Errorf("mutationGate = %v, want %q", decoded["mutationGate"], note)
			}
		})
	}
}

// Collision guard. The census deliberately shares its opt-in name with the
// mutating sandbox tests, which is what makes it one convention rather than two.
// The hazard is that the two drift apart in what the name means: one reader
// accepting "true", or a third reader appearing with an unrelated use. Both are
// cheap to detect here and expensive to detect in a run, so every other Go file
// in the repository that reads the literal must want the same value.
func TestMutateOptInNameIsNotReusedForAnotherMeaning(t *testing.T) {
	if MutateOptInEnv != "WEBULL_TRADE_MUTATE" {
		t.Fatalf("MutateOptInEnv = %q; the name is shared with the mutating sandbox "+
			"tests, so changing it silently creates a second convention", MutateOptInEnv)
	}
	root, err := repositoryRoot()
	if err != nil {
		t.Skipf("the repository root is not reachable from here, so the collision scan "+
			"did not run: %v", err)
	}

	var readers []string
	walkErr := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", "site", "gen", filepath.Base(currentDir(t)):
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		raw, readErr := os.ReadFile(path) //nolint:gosec // G304: a walk of the repository the tests live in.
		if readErr != nil {
			return readErr
		}
		for i, line := range strings.Split(string(raw), "\n") {
			if !strings.Contains(line, MutateOptInEnv) {
				continue
			}
			// Only a line that actually reads the environment is a reader. The
			// other two occurrences in each file are the comment above the gate and
			// the skip message inside it, and both are prose that happens to name
			// the variable; a reader's meaning is the one that compares it.
			if !strings.Contains(line, "Getenv(") {
				continue
			}
			readers = append(readers, fmt.Sprintf("%s:%d", path, i+1))
			if !strings.Contains(line, `"`+MutateOptInValue+`"`) {
				t.Errorf("%s:%d reads %s but does not require %q: %s",
					path, i+1, MutateOptInEnv, MutateOptInValue, strings.TrimSpace(line))
			}
		}
		return nil
	})
	if walkErr != nil {
		t.Skipf("walking the repository failed, so the collision scan did not run: %v", walkErr)
	}
	if len(readers) == 0 {
		t.Fatal("no other file in the repository reads " + MutateOptInEnv + ", so the claim " +
			"that this is a shared convention is no longer true. Either the mutating " +
			"sandbox tests moved - update the comment on MutateOptInEnv - or the name is " +
			"now free and the census may keep it.")
	}
	t.Logf("%s is read by %d other file(s): %s", MutateOptInEnv, len(readers), strings.Join(readers, ", "))
}

// The override's name is the other half of the gate and has three properties the
// second condition depends on, each cheap to break and expensive to discover in a
// run:
//
//  1. it is NOT the opt-in's name, and does not contain it. If it were, or did,
//     then one exported variable would satisfy both conditions and the gate would
//     be back to one condition;
//  2. it is scoped to this program, so a reader who has learned the trade-test
//     convention from WEBULL_TRADE_MUTATE cannot guess what it does;
//  3. it is read in exactly one place, this program, so the concession cannot be
//     inherited by something that never calls these endpoints.
//
// The third is not a second string scan: this module reads a constant, so a scan
// for the literal would be vacuous. Its value discipline is asserted where it is
// read instead, by TestNonSandboxOverrideAcceptsOnlyTheSameSingleValue.
func TestNonSandboxOverrideIsASeparateDeliberateName(t *testing.T) {
	if NonSandboxOverrideEnv == MutateOptInEnv {
		t.Fatalf("both conditions are read from %s, so one variable opens the gate again "+
			"and the second condition is not a second condition", NonSandboxOverrideEnv)
	}
	if NonSandboxOverrideEnv == "" {
		t.Fatal("NonSandboxOverrideEnv is empty; the override needs a name to be deliberate")
	}
	if !strings.HasPrefix(NonSandboxOverrideEnv, "WEBULL_LIVE_PROBE_") {
		t.Errorf("NonSandboxOverrideEnv = %q, want a WEBULL_LIVE_PROBE_ prefix so it is "+
			"scoped to this program and cannot be confused with the trade-test convention",
			NonSandboxOverrideEnv)
	}
	if strings.Contains(NonSandboxOverrideEnv, MutateOptInEnv) ||
		strings.Contains(MutateOptInEnv, NonSandboxOverrideEnv) {
		t.Errorf("%q and %q overlap; two names that contain one another are one name "+
			"spelled two ways", MutateOptInEnv, NonSandboxOverrideEnv)
	}
	if NonSandboxOverrideValue != MutateOptInValue {
		t.Errorf("NonSandboxOverrideValue = %q, want the same single token as "+
			"MutateOptInValue (%q); the two variables are read by the same comparison and "+
			"must not acquire different weights", NonSandboxOverrideValue, MutateOptInValue)
	}
	root, err := repositoryRoot()
	if err != nil {
		t.Skipf("the repository root is not reachable from here, so the scope scan did "+
			"not run: %v", err)
	}
	// The scan looks for the IDENTIFIER, not the literal. This module reads the
	// constant, so a scan for the string would find nothing at all and would pass
	// vacuously. The value discipline is not re-checked here because it is checked
	// where it is read, by TestNonSandboxOverrideAcceptsOnlyTheSameSingleValue and
	// by the equality assertion above; what is checked here is SCOPE, which is the
	// third property and the one a name loses silently.
	moduleDir := currentDir(t)
	var readers, foreign []string
	walkErr := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", "site", "gen":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		raw, readErr := os.ReadFile(path) //nolint:gosec // G304: a walk of the repository the tests live in.
		if readErr != nil {
			return readErr
		}
		for i, line := range strings.Split(string(raw), "\n") {
			if !strings.Contains(line, "Getenv(") || !strings.Contains(line, "NonSandboxOverrideEnv") {
				continue
			}
			where := fmt.Sprintf("%s:%d", path, i+1)
			readers = append(readers, where)
			if !strings.HasPrefix(path, moduleDir) {
				foreign = append(foreign, where)
			}
		}
		return nil
	})
	if walkErr != nil {
		t.Skipf("walking the repository failed, so the scope scan did not run: %v", walkErr)
	}
	if len(readers) == 0 {
		t.Fatal("no file reads " + NonSandboxOverrideEnv + ", so the override is " +
			"unreachable and the not-a-sandbox-host reason names a variable that does nothing")
	}
	if len(foreign) > 0 {
		t.Errorf("%s is also read outside this module: %s. It is the census's own "+
			"concession and means nothing to a program that never calls these "+
			"endpoints; a reader elsewhere would reasonably assume it did.",
			NonSandboxOverrideEnv, strings.Join(foreign, ", "))
	}
	t.Logf("%s is read by %d file(s), all in this module: %s", NonSandboxOverrideEnv,
		len(readers), strings.Join(readers, ", "))
}

// equalStrings compares two string slices, so a mismatch prints the whole list
// rather than the first element that differs.
func equalStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

// captureStdout runs fn with os.Stdout redirected to a pipe and returns what it
// printed. The pipe is drained by a goroutine so a table wider than the pipe
// buffer cannot deadlock, which the 193-row table in a real run would.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	original := os.Stdout
	os.Stdout = writer
	collected := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, reader)
		collected <- buf.String()
	}()
	func() {
		defer func() {
			os.Stdout = original
			_ = writer.Close()
		}()
		fn()
	}()
	out := <-collected
	_ = reader.Close()
	return out
}

// currentDir returns this process's working directory, used to keep the collision
// scan out of this module's own files.
func currentDir(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}
	return dir
}

// repositoryRoot walks up from the working directory to the directory holding the
// committed conformance manifest, which is this repository's root.
func repositoryRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, "conformance", "testdata", "manifest.json")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", errors.New("no conformance/testdata/manifest.json in any parent directory")
}
