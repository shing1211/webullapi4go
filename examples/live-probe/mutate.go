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

import "os"

// MutateOptInEnv names the environment variable that opens the mutating gate.
//
// The name is the one this repository already uses for the same class of action:
// the mutating trade sandbox test (trade/orders_sandbox_test.go:119) and the
// mutating order-event test (events/sandbox_test.go:117) are both gated behind
// it, and AGENTS.md records it as a separate opt-in from the ordinary sandbox
// gate. Reusing it is the point: one convention, one variable, one meaning -
// "this run may change state" - rather than a second spelling of the same idea
// that a reader has to learn separately.
//
// It is also a real hazard, and the reason the gate state is written into
// census.json rather than left implicit. A shell that exports
// WEBULL_TRADE_MUTATE=1 to run the mutating tests has, by the same variable,
// opted a census into calling the 34 endpoints below. The report therefore always
// records whether the gate was open for the run that produced it, so a census
// that called them can never be read as one that did not.
const MutateOptInEnv = "WEBULL_TRADE_MUTATE"

// MutateOptInValue is the only value of MutateOptInEnv that opens the gate.
//
// "1" and nothing else, exactly as both existing readers compare it
// (`os.Getenv("WEBULL_TRADE_MUTATE") != "1"`). A looser reading - any non-empty
// value, "true", "yes" - would give one variable two meanings, which is the
// collision this gate has to avoid rather than the convenience it would buy.
// An unset variable and an empty one are both closed.
const MutateOptInValue = "1"

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
// probe must not call unless a run has explicitly opted in through
// MutateOptInEnv.
//
// THE RULE, in two conditions, both read from committed or cached
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
//   - The host. The gate does not check that the run is pointed at a sandbox, so
//     WEBULL_TRADE_MUTATE=1 with a production credential calls these 34 for
//     real. That is why the default is closed: the operator who opens the gate is
//     the operator who chose the base URL, and the closed default means a
//     mistaken credential or a mistaken WEBULL_BASE_URL does nothing.
func IsMutating(ep Endpoint) bool {
	if _, safe := safeMethods[ep.Method]; safe {
		return false
	}
	_, mutating := mutatingAreas[areaOf(ep.Fixture)]
	return mutating
}

// MutationOptedIn reports whether this process has opted into calling the
// mutating class. It reads MutateOptInEnv once per call and compares the value
// exactly, so a shell variable set for another purpose, a stray "true" or a
// padded " 1" all leave the gate closed.
func MutationOptedIn() bool {
	return os.Getenv(MutateOptInEnv) == MutateOptInValue
}

// skippedMutationReason is the Skipped text of every endpoint the gate declined
// to call. It names the variable that opens the gate, because a row that says
// only "skipped" leaves a reader with a hole in the census and no way to fill
// it, and this is the one class of skipped row there is.
const skippedMutationReason = "mutating endpoint: the probe did not call it because " +
	MutateOptInEnv + "=" + MutateOptInValue + " is not set"

// gateClosed reports whether the gate refused this outcome. It is a named
// predicate rather than a field comparison so the loop, the summary and the
// report all key off one definition of "the probe chose not to send this".
func gateClosed(o Outcome) bool {
	return o.Skipped != ""
}
