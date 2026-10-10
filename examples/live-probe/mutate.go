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
	"net/url"
	"os"
	"strings"

	"github.com/shing1211/webullapi4go/client"
)

// MutateOptInEnv names the environment variable that is the FIRST of the two
// conditions that open the mutating gate.
//
// The name is the one this repository already uses for the same class of action:
// the mutating trade sandbox test (trade/orders_sandbox_test.go:119) and the
// mutating order-event test (events/sandbox_test.go:117) are both gated behind
// it, and AGENTS.md records it as a separate opt-in from the ordinary sandbox
// gate. Reusing it is the point: one convention, one variable, one meaning -
// "this run may change state" - rather than a second spelling of the same idea
// that a reader has to learn separately.
//
// IT IS NOT SUFFICIENT. A shell that exports WEBULL_TRADE_MUTATE=1 to run the
// mutating tests has, by that one variable alone, opened nothing: the second
// condition - a base URL on a host the SDK derives for a sandbox - is also
// required, with the one deliberate exception of [NonSandboxOverrideEnv]. See
// [ResolveMutationGate]. The reason is the hazard the shared name creates, which
// a second condition removes rather than merely documents: the base URL comes
// from the environment too, so exporting the mutate variable and pointing the SDK
// at production would otherwise issue the 34 endpoints below - three of them
// scalar-bodied, and therefore made inert by no other rule - against a real
// account.
//
// The gate state is still written into census.json rather than left implicit, so
// a census that called them can never be read as one that did not.
const MutateOptInEnv = "WEBULL_TRADE_MUTATE"

// MutateOptInValue is the only value of MutateOptInEnv that satisfies the first
// condition.
//
// "1" and nothing else, exactly as both existing readers compare it
// (`os.Getenv("WEBULL_TRADE_MUTATE") != "1"`). A looser reading - any non-empty
// value, "true", "yes" - would give one variable two meanings, which is the
// collision this gate has to avoid rather than the convenience it would buy.
// An unset variable and an empty one are both closed.
const MutateOptInValue = "1"

// NonSandboxOverrideEnv names the environment variable that waives the SECOND of
// the two conditions, the sandbox-host one, and only that one.
//
// It is a deliberately separate name. Reinterpreting WEBULL_TRADE_MUTATE=1 to
// mean "yes, really, not a sandbox" would leave a run pointed at production with
// one variable set and nothing in the shell to show that a second and larger
// decision had been made. Two names, both deliberately set, is what makes the
// intent legible at the moment of the mistake rather than afterwards.
//
// The name is scoped to this program (WEBULL_LIVE_PROBE_...) because the
// concession is not one the trade sandbox tests are making: they always run
// against a sandbox, so a variable of theirs would name a situation that cannot
// arise for them. It is also the census-specific half of the pair, so a reader
// who has learned WEBULL_TRADE_MUTATE from the mutating trade test is not misled
// about what this one does.
const NonSandboxOverrideEnv = "WEBULL_LIVE_PROBE_MUTATE_NON_SANDBOX"

// NonSandboxOverrideValue is the only value of NonSandboxOverrideEnv that
// waives the sandbox-host condition.
//
// The same single token as [MutateOptInValue], and deliberately so: both
// variables are read with the same exact comparison, so a reader who has learned
// one has learned the other, and a value that opens one of them cannot quietly
// do something else. The two variables remain separately named and separately
// required; sharing a value shares no meaning.
const NonSandboxOverrideValue = MutateOptInValue

// safeMethods are the HTTP methods that cannot change server state by definition
// (RFC 9110 §9.2.1). A documented method outside this set is treated as unsafe,
// which includes a method the census could not read at all: the gate fails
// closed on an endpoint whose method is unknown.
var safeMethods = map[string]struct{}{
	"GET":     {},
	"HEAD":    {},
	"OPTIONS": {},
	"TRACE":   {},
}

// mutatingAreas are the documentation areas that act on account state: the
// trading API and the two broker APIs. The area comes from the manifest's own
// area - the first segment of the fixture name, read by areaOf - so the gate and
// the per-area summary count the same partition of the corpus by the same field
// and cannot disagree about which area a row is in.
//
// An area set rather than a path list is deliberate. A list of paths is
// hand-maintained per endpoint, so it is out of date the moment Webull
// documents a new one, and it reads as exhaustive while covering 34 of the 51
// unsafe methods in the corpus. An area is a section of the documentation that
// a reader already navigates by, and a new area fails closed: it is not on this
// list, so a POST in it is called, and the summary prints the per-area mutating
// count so the exposure is visible in the artifact rather than assumed absent.
var mutatingAreas = map[string]struct{}{
	"trading":      {},
	"broker-hk":    {},
	"broker-fd-us": {},
}

// IsMutating reports whether an endpoint is in the mutating class: one this
// probe must not call unless a run has satisfied both conditions of
// [ResolveMutationGate].
//
// THE CLASSIFICATION, in two conditions, both read from committed or cached
// documentation rather than from a list written here:
//
//  1. the documented HTTP method is not a safe one; and
//  2. the endpoint's area is trading, broker-hk or broker-fd-us.
//
// Against the committed corpus that is 34 of 193 endpoints, and it contains every
// endpoint that can place an order, cancel one, move money, or create, modify or
// close an account - including /trading/orders/cancel,
// /broker/accounts/virtual-accounts/create, /broker/accounts/virtual-accounts/update
// and /broker/funding/ach/relationships/create, whose bodies are all scalars and
// therefore carry no array for any value rule to empty.
//
// WHAT IT CANNOT CATCH. Stated here because a classifier that reads as
// exhaustive and is not is worse than a narrow one:
//
//   - An unsafe method outside those three areas. The corpus has 17: 6 in
//     market-data-watchlist, 5 in display-solution, 4 in authentication, 1 in
//     connect-api and 1 in market-data-stock. The gate calls all 17 by default.
//     The watchlist ones add and remove symbols on a user's list and the
//     display-solution ones write to a hosted display; the line this rule draws
//     is money and orders, not every write a service offers. If that line is
//     wrong the fix is to widen mutatingAreas, not to add paths.
//   - A documented safe method that mutates. Nothing in this corpus does, but a
//     GET that changed state would pass both conditions and be reported as a
//     read. The gate trusts the published method and cannot audit the server.
//   - The preview endpoints. /trading/orders/preview and /broker/orders/preview
//     are in the class and are called like any other. They are almost certainly
//     inert, and exempting them would mean a per-path exception list, which is
//     the maintenance this rule exists to avoid.
//   - Anything about the request's content, and therefore anything the probe
//     would send. The gate asks what the endpoint is; it does not ask what the
//     body looks like. It is not a substitute for the synthesised empty array in
//     prepare, and the two are independent: /trading/orders/place is both gated
//     and inert, /broker/accounts/create is gated with nothing to empty, and a
//     market-data-watchlist POST is called with an empty array and protected by
//     nothing else. Both barriers are required, and neither is sufficient alone.
func IsMutating(ep Endpoint) bool {
	if _, safe := safeMethods[ep.Method]; safe {
		return false
	}
	_, mutating := mutatingAreas[areaOf(ep.Fixture)]
	return mutating
}

// knownRegions are the regions the SDK publishes, listed so that a sandbox host
// in ANY of them is recognised rather than only the one AGENTS.md names as the
// host that may appear in committed material.
//
// The list is a list of the SDK's own region identifiers, never of hosts: every
// host it yields is derived through [client.EndpointsFor], so there is no host
// string in this file that can drift from the SDK's own knowledge of where
// Webull serves each region.
//
// The SDK exposes no enumeration of its regions, so the identifiers have to be
// named, and a region added to the SDK after this was written is not on the list.
// That direction is safe: such a census is REFUSED for a non-sandbox host, and
// the operator sees the not-a-sandbox-host reason, rather than being let through.
// TestKnownRegionsAreAllValidAndYieldNoProductionHost pins every entry against
// the SDK, so a typo here cannot silently shrink or widen the set.
var knownRegions = []client.Region{
	client.HK, client.US, client.JP, client.SG, client.TH, client.AU,
	client.MY, client.UK, client.BR, client.MX, client.ZA, client.EU,
}

// IsSandboxHost reports whether host is a service host the SDK derives for the
// sandbox environment of a region it knows.
//
// Both the REST host and the broker REST host are considered for every region.
// The census addresses every documented path, including the /broker/... ones,
// through the single REST base URL - it calls [client.Client.DoStream], never
// [client.Client.DoBroker] - so checking only the REST host would be the narrower
// of the two readings of "where is this run pointed", and either host is a
// sandbox Webull published.
//
// The comparison is EXACT, after normalisation, against hosts derived from the
// SDK - not a substring test for the word "sandbox". A heuristic here would be a
// way to switch the gate off by choosing a similar-looking host, which is the
// precise failure this function exists to prevent. It also means a host that
// merely contains the label, such as sandbox.attacker.example, is refused.
//
// WHY THE CONFIGURED ENVIRONMENT IS NOT CONSULTED. [client.Client.Environment] can
// disagree with the base URL the run will actually use: [client.WithEnv] applies
// WEBULL_BASE_URL after WEBULL_ENVIRONMENT and overwrites only
// Config.Endpoints.HTTP, so WEBULL_ENVIRONMENT=sandbox together with
// WEBULL_BASE_URL=https://api.webull.hk resolves to a client that reports
// "sandbox" and talks to production. A gate that read the environment field would
// open on exactly the configuration that must not open it. Only the host answers
// the question being asked.
func IsSandboxHost(host string) bool {
	normalised := normaliseHost(host)
	if normalised == "" {
		return false
	}
	for _, region := range knownRegions {
		sandbox := client.EndpointsFor(region, client.Sandbox)
		if normalised == hostnameOf(sandbox.HTTP) || normalised == hostnameOf(sandbox.BrokerHTTP) {
			return true
		}
	}
	return false
}

// hostnameOf returns the host of a URL with any port, userinfo and path removed,
// and "" when raw is not a URL with a host. It is a second reading of the same
// field [hostOf] reports, and the difference is deliberate: the report's `host`
// shows the operator exactly what was configured, while a comparison needs the
// name alone, so a base URL written with an explicit port still matches.
func hostnameOf(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return parsed.Hostname()
}

// normaliseHost reduces a host to the form the comparison in [IsSandboxHost] is
// written in: lower case, no port, no trailing dot, no surrounding whitespace.
//
// host is a host, optionally with a port, as [hostOf] returns it; a full base URL
// is accepted too and reduced to its host, because refusing to recognise a known
// sandbox because it arrived as a URL rather than as a bare name would be a false
// refusal for no gain - the accepted SET is the same either way, and this function
// cannot add to it. A DNS name is case-insensitive and a fully qualified one may
// carry a trailing dot, and an operator writing a base URL by hand may use either.
// Normalising them is not leniency about WHICH host it is - the comparison
// afterwards is still exact equality against a set derived from the SDK - it is
// only refusing to fail in the safe direction over a capital letter. The gate's
// cost when it is wrong in the unsafe direction is an irreversible request; its
// cost when it is wrong in the safe direction is one visible reason string and an
// override.
func normaliseHost(host string) string {
	host = strings.ToLower(strings.TrimSpace(host))
	if host == "" {
		return ""
	}
	// A URL is read as a URL; anything else is read as a host:port authority,
	// which is what drops the port and unwraps an IPv6 literal's brackets.
	if parsed, err := url.Parse(host); err == nil && parsed.Host != "" {
		host = parsed.Host
	} else if parsed, err := url.Parse("//" + host); err == nil {
		if name := parsed.Hostname(); name != "" {
			host = name
		}
	}
	return strings.TrimSuffix(host, ".")
}

// MutationOptedIn reports whether this process has satisfied the FIRST of the two
// gate conditions: MutateOptInEnv carries MutateOptInValue.
//
// It reads the variable once per call and compares the value exactly, so a shell
// variable set for another purpose, a stray "true" or a padded " 1" all leave the
// condition unmet. It is deliberately a predicate on ONE condition and not a
// predicate on the gate: the gate is [ResolveMutationGate], and a reader who needs
// "may this run mutate anything" wants that.
func MutationOptedIn() bool {
	return os.Getenv(MutateOptInEnv) == MutateOptInValue
}

// NonSandboxOverridden reports whether this process has satisfied the override
// that waives the sandbox-host condition. It reads NonSandboxOverrideEnv once per
// call and compares the value exactly, for the same reason [MutationOptedIn] does:
// one accepted spelling per variable, so neither can acquire a second meaning.
func NonSandboxOverridden() bool {
	return os.Getenv(NonSandboxOverrideEnv) == NonSandboxOverrideValue
}

// MutationGate is the resolved state of the mutating gate for one run: the two
// conditions, the host they were judged against, and the override.
//
// It is a value rather than a package-level bool so that the loop, the report
// line and the tests cannot disagree about why a run did or did not mutate
// anything. [ResolveMutationGate] reads the environment once and every consumer
// reads this struct, so there is exactly one reading of the environment per run.
type MutationGate struct {
	// OptedIn reports condition one: MutateOptInEnv is set to MutateOptInValue.
	OptedIn bool
	// SandboxHost reports condition two: Host is a host the SDK derives for a
	// sandbox environment. See [IsSandboxHost] for what "sandbox" means here and
	// why the configured environment is not consulted instead.
	SandboxHost bool
	// Overridden reports that NonSandboxOverrideEnv was set to
	// NonSandboxOverrideValue, which waives condition two. It cannot waive
	// condition one.
	Overridden bool
	// Host is the base URL host the gate judged, exactly as [hostOf] reported it
	// and exactly as the report records it, so the reason string and the
	// artifact's `host` field can never name different deployments. It is empty
	// only where there is no client to ask, which is a closed gate.
	Host string
}

// Open reports whether this run may call the mutating class.
//
// Both conditions must hold, except that the second is satisfied either by a
// sandbox host or by the override:
//
//	opted in AND (sandbox host OR override)
//
// The override cannot open a run that has not opted in. That is what makes the
// two conditions independent rather than two ways in, and it is why a run that
// forgot the opt-in is skipped with the not-opted-in reason even when it is
// pointed at a sandbox and even when the override is set.
func (g MutationGate) Open() bool {
	if !g.OptedIn {
		return false
	}
	return g.SandboxHost || g.Overridden
}

// SkipReason is the text to record as an [Outcome.Skipped] for a mutating
// endpoint this gate refuses, and "" when the gate is open.
//
// The two reasons are different FACTS and a reader has to be able to tell them
// apart: "nobody opted in" is fixed by setting one variable, while "this was not
// a sandbox host" says the run may be pointed at real account state and is fixed
// by pointing it somewhere else, or by setting the override on purpose. Reporting
// one as the other would send an operator to make the wrong change - and would let
// a production run be read as merely un-authorised, which is the confusion this
// whole gate exists to remove.
//
// The opt-in is reported first when both are unmet, because it is the condition
// that holds regardless of where the run is pointed and is therefore the more
// basic of the two.
func (g MutationGate) SkipReason() string {
	if g.Open() {
		return ""
	}
	if !g.OptedIn {
		return skippedNotOptedIn()
	}
	return skippedNotSandbox(g.Host)
}

// skippedNotOptedIn is the Skipped text of every mutating endpoint this gate
// declined because MutateOptInEnv is not set to MutateOptInValue.
//
// It names the variable that opens the gate, because a row that says only
// "skipped" leaves a reader with a hole in the census and no way to fill it. It
// deliberately does not name the host, the override variable, or anything else:
// this reason says the run was not authorised, and adding the conditions that
// were not the problem would blur exactly the distinction the two reasons exist
// to keep.
func skippedNotOptedIn() string {
	return "mutating endpoint: not opted in - " + MutateOptInEnv + "=" + MutateOptInValue +
		" is not set, so the probe did not call it"
}

// skippedNotSandbox is the Skipped text of every mutating endpoint this gate
// declined because the run is not pointed at a sandbox host and the override was
// not set.
//
// It names the host that was judged - so a reader can see which deployment the
// gate thought it was protecting, and so an operator who expected a sandbox is
// told which host is actually configured - and it names the override variable as
// the deliberate way past this condition, because a run that genuinely needs a
// non-sandbox base URL should have to ask for it by name rather than discover it
// by relaxing the opt-in.
func skippedNotSandbox(host string) string {
	judged := host
	if strings.TrimSpace(judged) == "" {
		judged = "the configured base URL"
	}
	return "mutating endpoint: not a sandbox host - the probe is pointed at " + judged +
		", which is not a sandbox host the SDK derives, so it did not call it; set " +
		NonSandboxOverrideEnv + "=" + NonSandboxOverrideValue + " to override this condition"
}

// ResolveMutationGate reads both conditions once, for a run pointed at host, and
// returns the resolved gate. [Census] and main both call it with the same
// [hostOf] of the same client, so the rows in the artifact and the mutationGate
// line above them are decided by one reading of the environment and cannot
// disagree.
func ResolveMutationGate(host string) MutationGate {
	return MutationGate{
		OptedIn:     MutationOptedIn(),
		SandboxHost: IsSandboxHost(host),
		Overridden:  NonSandboxOverridden(),
		Host:        host,
	}
}

// gateClosed reports whether the gate refused this outcome. It is a named
// predicate rather than a field comparison so the loop, the summary and the
// report all key off one definition of "the probe chose not to send this".
func gateClosed(o Outcome) bool {
	return o.Skipped != ""
}
