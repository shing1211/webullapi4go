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

package conformance

import "strconv"

// The reasons and unblock requirements for the committed live baseline.
//
// They are here, and not in the generator that writes live-divergences.json,
// because a required reason is not a required *true* reason and the only defence
// against a false one is that a person wrote it after reading the finding. The
// discipline is the one known-divergences.json uses, and the failure it has
// already suffered is the reason this file exists in this shape: a false reason,
// repeated across many entries, satisfies a required-field check perfectly.
//
// Each entry below is keyed by the same tuple the finding is, so a finding that
// moves is a missing key rather than a stale sentence, and a key with no entry
// below stops the regeneration.

// liveBaselinePolicy is the gate rule, written into the file so the rule travels
// with the data it governs.
const liveBaselinePolicy = "The gate passes when the classified set exactly equals the recorded " +
	"set, and it fails in both directions. A finding in a bucket the file does not record fails " +
	"immediately, because it is new. A recorded finding that no longer reproduces also fails, " +
	"because a finding that was fixed and a check that changed meaning look identical from here " +
	"and both need a person to look at the tree. Every entry carries a reason saying why it is " +
	"recorded rather than fixed, and an unblock saying what would close it; a required reason is " +
	"not a required true reason, and the reasons here are held to being true by review rather than " +
	"by the schema. This set is separate from known-divergences.json on purpose: a live finding is " +
	"a claim about Webull's server, and recording one in the SDK's backlog would make the " +
	"documented rows and the live rows read as one list of defects when the live rows are not " +
	"claims about this SDK."

// liveBaselineNotes are the statements a reader cannot get from a count: what
// each bucket means, and what the coverage numbers are a count of.
func liveBaselineNotes(run LiveRun) map[string]string {
	counts := BucketCounts(run.Classifications)
	return map[string]string{
		"liveOnly": "live-only is the new signal: " +
			"the server sent something the documentation does not describe, which is the class the " +
			"documented comparison is structurally blind to. " +
			"It holds " + countNote(counts[LiveOnly]) + " finding(s).",
		"both": "both means the same finding was observed on both sides with the same detail, so " +
			"it is a pre-existing documentation-versus-SDK disagreement the live server reproduces " +
			"rather than new evidence. It holds " + countNote(counts[LiveBoth]) + " finding(s).",
		"bothDetailChanged": "both-detail-changed means the same key was observed on both sides " +
			"with *different* detail, so the two bodies disagree about the same disagreement. It is " +
			"neither agreement nor two findings, and a naive set difference loses it. It holds " +
			countNote(counts[LiveBothDetailChanged]) + " finding(s), and it is empty here, which is a " +
			"result rather than a gap: no endpoint in this tree has the SDK disagreeing with the " +
			"server and the page in two different ways at once.",
		"docsOnly": "docs-only means the documentation and the SDK disagree where the server and " +
			"the SDK do not, so it is evidence about the page rather than about the SDK. It holds " +
			countNote(counts[LiveDocsOnly]) + " finding(s).",
		"notComparable": "notComparable names the " + countNote(len(run.Coverage.NotComparable)) +
			" row(s) whose SDK method sends a path its page does not document. The documented half " +
			"of those rows claims nothing, and the live half still runs: a skeleton was captured by " +
			"calling the SDK method itself, so the path question does not bear on it.",
		"decodeRejected": "decodeRejected names the " + countNote(len(run.Coverage.DecodeRejected)) +
			" captured body the SDK's own type could not unmarshal. All of them are top-level-kind " +
			"mismatches, and that is the strongest signal this harness can produce: the documented " +
			"comparison is blind to it by construction, because a page and a fixture that agree " +
			"with each other decode cleanly whether or not either matches the server.",
		"topLevelIsNotLive": "checkShape reads the top level from the manifest's recorded checks " +
			"rather than from the body, so a top-level-shape-mismatch is the same observation in " +
			"both runs and is a claim about the page. The live evidence for a container-kind " +
			"disagreement arrives as a decode rejection instead, and the two are different kinds " +
			"because they are different checks. See the package comment in live.go.",
		"oneRun": "Every figure here comes from one authorised run against api.sandbox.webull.hk on " +
			"2026-09-29 and there is no second one to compare it against. A finding here is what " +
			"that host answered once; nothing here says what production answers.",
	}
}

// countNote renders a count for the notes above, so they read as sentences
// rather than as a format string.
func countNote(n int) string { return strconv.Itoa(n) }

// liveReasonFor returns the reason and unblock requirement for a classification.
//
// It is keyed on the same tuple the finding carries, so a finding that changes
// kind, name or side is a missing key -- which the regeneration fails on -- and
// not a sentence that has quietly stopped applying.
func liveReasonFor(c LiveClassification) (reason, unblock string, ok bool) {
	key := c.Symbol + "|" + string(DivergenceKind(c.Kind)) + "|" + c.Name
	switch key {
	case "data.GetCapitalFlow|top-level-shape-mismatch|":
		return "The page documents an object and data.GetCapitalFlow decodes a slice, so the " +
				"documented comparison reports the inversion in both directions -- the shape check " +
				"reads the manifest rather than the body, which is why the two runs produce one " +
				"identical finding. It is recorded here rather than in known-divergences.json because " +
				"the live evidence changes what it means: the sandbox answered an array, so the server " +
				"and the SDK agree and the page is the outlier. A documented comparison cannot see " +
				"that, because the fixture and the type were made to match each other.",
			"None needed for the live evidence, which is already recorded. Closing the " +
				"documentation half needs a correction on the reference page, which is Webull's to " +
				"make; the fixture regenerates from it with no credential.", true

	case "data.GetCapitalFlow|decode-failure|":
		return "The documented fixture does not unmarshal into []data.CapitalFlowEntry, and the " +
				"decode check reports the same top-level inversion the shape check named. It is " +
				"docs-only because the live body -- an array -- decodes cleanly, so the live server " +
				"does not reproduce the failure. The decode check is the harness's weakest and it is " +
				"recorded beside the sharper finding rather than instead of it, which is why both rows " +
				"exist for one defect.",
			"Defers to the top-level-shape-mismatch row for data.GetCapitalFlow, which " +
				"names the same defect: a correction on the reference page. No credential is " +
				"required to close the documented half.", true

	case "data.GetDisplayGainersLosers|decode-failure|":
		return "The sandbox answered a bare JSON array and the SDK decodes " +
				"types.Page[data.ScreenerStock], so the response does not decode. This is a " +
				"top-level-kind mismatch and it is invisible to the documented comparison by " +
				"construction: the page and the documented fixture agree with the SDK type, so that " +
				"run is green. A caller on this endpoint gets an error or an empty page where the " +
				"screener's rows are. It is recorded rather than fixed because the fix depends on " +
				"which is right -- the SDK decoding the bare payload, or the endpoint wrapping it -- " +
				"and the live response is evidence about the server, not a specification.",
			"A decision from Webull on whether this endpoint is paginated, plus a US " +
				"production or sandbox credential to confirm the shape is not region-specific. The HK " +
				"sandbox cannot settle it: Display Solution needs a paid subscription and its host " +
				"returns 403 even with valid credentials.", true

	case "data.GetDisplayTopActive|decode-failure|":
		return "Identical in kind to data.GetDisplayGainersLosers and for the same reason: the " +
				"sandbox answered a bare array and the SDK decodes a paginated page type. Recorded " +
				"separately because it is a separate endpoint with a separate response and a separate " +
				"fix, and merging the two would make one endpoint's repair close the other's finding.",
			"The same decision and the same credential as data.GetDisplayGainersLosers: the " +
				"endpoint's pagination contract from Webull, and a Display Solution entitlement the HK " +
				"sandbox does not grant.", true

	case "data.GetFuturesBars|decode-failure|":
		return "The sandbox answered a bare JSON array and the SDK decodes data.BatchBars, a " +
				"struct, so the response does not decode. A top-level-kind mismatch, and again one the " +
				"documented comparison cannot see: the page documents an object and the documented " +
				"fixture carries the {symbol, result} envelope the SDK type expects, so that run is " +
				"green. Note the live body does contain a result key, but at the top level of each " +
				"array element rather than under an envelope, which is what the decode rejection is " +
				"reporting.",
			"A US sandbox or production credential. HK futures are blocked here, so the capture " +
				"was made against an endpoint whose full response shape this harness has not seen.", true

	case "data.GetFuturesTick|missing-required-name|instrument_id":
		return "The strongest and most actionable finding in the set. The page requires " +
				"instrument_id and the SDK tags the field instrument_id, so the documented comparison " +
				"is green. The live response sends instrumentId -- camelCase -- and encoding/json " +
				"matches a member exactly and then case-insensitively, and the underscore is not a " +
				"case, so the field is left at its zero value with no error. A caller reading " +
				"StockTicks.InstrumentID from this endpoint gets \"\". It is live-only precisely " +
				"because the documentation does not describe the wire: this is the class of defect " +
				"the whole live harness exists to find.",
			"Nothing further is needed to fix it -- the tag or the field name is wrong on one " +
				"side. It is recorded rather than fixed because it was found by a single HK sandbox " +
				"run and changing a public DTO's wire name is a breaking change for anyone who has " +
				"already written against it, so it needs the maintainer's decision rather than a " +
				"silent edit. Confirming against a second host would establish whether the camelCase " +
				"spelling is endpoint-specific or SDK-wide.", true

	case "data.GetStockInstruments|top-level-shape-mismatch|":
		return "The SDK decodes []data.StockInstrument and the live response is an object carrying " +
				"data and pagination_key, so the two disagree at the top level. Recorded live-only " +
				"because the row is not comparable on the documented side: the SDK method sends " +
				"/openapi/instrument/stock/list where the page documents " +
				"/trading/instruments/stocks/profiles/list, so the page is not this call's contract " +
				"and claims nothing. The live run still executes, and it is one of the four bodies " +
				"the SDK's own type could not decode.",
			"A decision on which path this endpoint is meant to serve, and then a credential for " +
				"whichever environment serves it. US-only surfaces return 404 from the HK sandbox, " +
				"so no HK credential can close this.", true

	case "data.GetStockInstruments|decode-failure|":
		return "The same container disagreement as the row above, reported by the decode check " +
				"because that is the check that reads the body. Recorded beside the sharper finding " +
				"rather than instead of it, which is why one defect produces two rows: the shape check " +
				"names the inversion and the decode check reports what a caller experiences.",
			"Defers to the top-level-shape-mismatch row for data.GetStockInstruments, which " +
				"names the same container disagreement and the same fix.", true

	case "data.GetStockInstruments|missing-declared-name|data":
		return "The page declares a top-level name data and the SDK's []data.StockInstrument " +
				"carries no field tags for it, because a slice of elements has nowhere to put a " +
				"top-level name. The live response does send data, so a server-side envelope is being " +
				"answered and the SDK cannot read it -- which is the same defect as the two rows " +
				"above, observed from a third angle. It is recorded rather than suppressed because " +
				"naming a name the SDK cannot reach is what the check is for, and suppressing it here " +
				"would be a special case inside a safety check.",
			"Defers to the top-level-shape-mismatch row for data.GetStockInstruments. No " +
				"credential is required to know the field cannot be carried; a credential is " +
				"required to know whether the endpoint should be returning a slice at all.", true

	case "data.GetStockInstruments|missing-declared-name|pagination_key":
		return "As for data: the page declares pagination_key, the SDK's element type cannot " +
				"carry it, and the live response sends it. A second name-level view of the same " +
				"container disagreement as the top-level-shape-mismatch row for " +
				"data.GetStockInstruments, recorded separately because the rule is per name and a " +
				"name the SDK cannot reach is one caller-visible field rather than an abstraction.",
			"Defers to the top-level-shape-mismatch row for data.GetStockInstruments, which " +
				"names the same defect and the same unblock requirement: a decision on which " +
				"path the endpoint serves, then a credential for that environment.", true

	case "trade.GetOrderDetail|missing-required-name|client_order_id":
		return "The page marks client_order_id required and the SDK tags the field, so the " +
				"documented comparison is green. The live response is {\"orders\": []} and carries " +
				"neither that name nor combo_type: the sandbox held no order for the probed " +
				"client_order_id, so it answered an empty group. This is recorded as a live finding " +
				"because the evidence is genuinely absent rather than mismatched, and stating it that " +
				"way is the honest description -- it is a statement about one probe against an empty " +
				"account, not a demonstrated defect.",
			"A probe against a sandbox account holding at least one order, which requires a " +
				"credential for an account with order history. Until then this row records what one " +
				"empty-account response looked like, and the gate will report it changing if a future " +
				"run answers with a populated group.", true

	case "trade.GetOrderDetail|missing-required-name|combo_type":
		return "As for client_order_id, and for the same reason: the page marks combo_type " +
				"required and the live {\"orders\": []} carries it on no member. Recorded per name " +
				"because the rule is per name, so the two omissions stay separately visible.",
			"Defers to the missing-required-name row for client_order_id on the same " +
				"endpoint: a probe against a sandbox account that holds at least one order.", true

	case "trade.GetPositions|missing-required-name|option_strategy":
		return "The page marks option_strategy required on the positions list and the SDK tags " +
				"it, so the documented comparison is green. The live position carries cost_price, " +
				"currency, instrument_type, last_price, position_id, quantity, symbol and " +
				"unrealized_profit_loss, and no option_strategy: the probed account held an equity " +
				"position, and option_strategy is an option-only field. So this is a documentation " +
				"observation rather than a demonstrated SDK defect -- a required list naming a field " +
				"that is absent for every non-option position is a promise the server does not keep " +
				"for equities, and the finding is recorded as a live-only observation for exactly " +
				"that reason.",
			"A probe against an account holding an option position, and a statement from Webull " +
				"on whether option_strategy is required only for options. An account with options " +
				"requires a sandbox credential the HK default account does not have.", true
	}
	return "", "", false
}
