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

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
)

// # What this file is for
//
// Everything else in this package compares the SDK against Webull's published
// documentation. That can only ever prove the SDK matches the documentation, and
// the two can agree with each other and both be wrong: a page and the SDK type
// that were made to match decode a body neither describes. That is not a
// hypothetical here. brokerfd tagged average_cost where the page requires
// cost_price and reconciled as a clean match, and a live call returns a zero
// money.Money with no error at all.
//
// The fix is a body that came from somewhere else. examples/live-probe calls the
// sandbox, reduces each response to a value-free skeleton of the member names
// and kinds it sent, and commits it under testdata/live/. This file runs the same
// five checks in shapes.go against those skeletons and records what the server
// said the SDK cannot hold.
//
// # The two runs, and why there are two
//
// For one endpoint there are two bodies and one type, so CompareBody is called
// twice:
//
//   - the live skeleton, which is what the sandbox sent. The direction is
//     LiveSDKDirection, and a finding here is a claim about Webull's server.
//   - the documented fixture, which is the minimal instance of the published
//     schema. The direction is LiveDocsDirection, and a finding here is a claim
//     about the SDK against the page -- the same claim known-divergences.json
//     already records.
//
// A page and a fixture disagree with a type in one way and the server disagrees
// with it in another, and the difference is the whole point. A name the page
// declares and the SDK cannot reach is a real defect whether or not the server
// sent the field; a name the server sent and the SDK cannot reach is a defect
// the page does not describe and the documented run is structurally blind to.
// Collapsing the two would file 25 documented rows plus N live rows as one
// backlog of 25+N, half of which are claims about Webull's server rather than
// about this SDK.
//
// So the directions are separate fields on the finding and separate parts of the
// key, and this file has its own kinds, its own baseline and its own gate. It
// does not reuse known-divergences.json's loader, its kinds or its exact-set
// rule, and a live finding recorded there would be a category error: that file is
// the SDK's backlog, and a row in it says the SDK is wrong.
//
// # What the live run can and cannot see
//
// This is a limit of CompareBody, not a choice made here, and it is stated
// because the two runs are otherwise easy to over-read.
//
// Four of the five checks read the body they are given. The required-name check
// reads the page's required list and looks each name up *in the body*; the
// declared-name check reads the page's declared names and compares them with the
// type's tags and the body's; the leaf check reads the body's value kinds; the
// decode check unmarshals the body. Against a skeleton those four are live
// evidence.
//
// checkShape does not read the body. It compares f.Checks.TopLevel -- a field of
// the Fixture, recorded from the page by the generator -- with the SDK type's
// shape. So a TopLevelMismatch or ElementTypeMismatch is the *same* observation
// in both runs, whatever body each was handed, and it is a claim about the page
// dressed as a claim about the wire. The neutral phrasing below is what makes
// that visible: both runs produce one row, they pair, and the row lands in
// LiveBoth. It is reported rather than suppressed, because a caller reading the
// baseline should see that the top-level disagreement is a documentation fact and
// not mistake the bucket for live corroboration of it.
//
// The consequence for the strongest signal in the tree is worth stating plainly.
// When the server's top-level kind differs from the SDK's, the top-level check
// cannot see it, and the *decode* check is what reports it -- which is exactly
// why live-manifest.json records decodedCleanly, and why all four captured
// bodies the SDK type could not decode are top-level-kind mismatches. The live
// evidence for a container-kind disagreement arrives as a decode rejection, not
// as a shape mismatch, and the two are recorded as different kinds because they
// are different checks. A reader looking for a TopLevelMismatch and finding none
// for those four rows has not been given a false negative: the shape check never
// could have reported them.
//
// # Where the checks come from
//
// CompareBody, unmodified. A second implementation of the five checks would be a
// second opinion that could be wrong in a different direction, and the whole
// value of this file is that it is the documented comparison pointed at a
// different body. Nothing in shapes.go knows a skeleton exists: the live run is
// CompareBody with a different fourth argument.
//
// That has one consequence worth stating, because it is a limit rather than a
// choice. The required-name check reads the page's required list out of the
// Fixture and asks whether the type carries a tag for it. Against a live body
// that list is the page's, not the server's, so the check answers "could the SDK
// hold what the page promises, and did the server send it" and not "what did the
// server send that the SDK cannot hold". The second question is answered by the
// declared-name check and by the shape and leaf checks, which do read the body.
// The required list is the page's by construction and stays that way.

// LiveDirection says which of the two bodies a finding was observed against, and
// it is part of a finding's identity rather than decoration.
//
// The two directions are different claims about different parties, so a consumer
// keying on kind and name alone would collapse them: a name the page declares
// and the server both send, which the SDK drops, is two findings, and reporting
// one of them loses half of what happened.
type LiveDirection string

const (
	// LiveSDKDirection is a finding about the SDK type against the live
	// skeleton: a claim about Webull's server and about this SDK, observed on
	// the wire. It is the new signal this harness exists to produce.
	LiveSDKDirection LiveDirection = "live-sdk"

	// LiveDocsDirection is a finding about the SDK type against the documented
	// fixture: the same class of claim known-divergences.json records, restated
	// here so the two can be diffed against each other.
	LiveDocsDirection LiveDirection = "live-docs"
)

// LiveDivergence is one observed disagreement on one side of the comparison.
//
// It is deliberately not conformance.Divergence. That type's Detail is part of a
// baseline entry's identity, and a wording that changes with the evidence would
// delete a finding and insert another. Here Detail is a human-readable account
// and the key is the tuple below, so a check can be reworded without the
// baseline moving.
type LiveDivergence struct {
	// Symbol is the SDK method whose type was compared, e.g. data.GetFuturesBars.
	Symbol string `json:"symbol"`
	// Fixture is the documentation fixture ID, which is how a live finding and a
	// documented one about the same endpoint join up.
	Fixture string `json:"fixture"`
	// Direction says which body the finding was observed against.
	Direction LiveDirection `json:"direction"`
	// Kind is the specific disagreement, from the harness's own vocabulary in
	// shapes.go. A live finding is not a new kind of disagreement; it is the same
	// kinds observed against a different body, which is why they are reused.
	Kind DivergenceKind `json:"kind"`
	// Name is the property the finding is about, or "" for one about the shape
	// rather than any single name.
	Name string `json:"name,omitempty"`
	// Detail states the disagreement in terms a reviewer can check against the
	// two sources without rerunning anything.
	Detail string `json:"detail"`
}

// Key is the identity a single observation is recorded by, and it includes the
// direction.
//
// Detail is excluded on purpose, for the reason LiveDivergence states: a
// reworded detail is the same finding, and including it would make every
// rephrasing of a check a deleted entry and an inserted one. Direction is
// included because a live-sdk finding and a live-docs finding about the same
// name are different claims about different parties, and a baseline that held
// one key for both would silently drop one of them.
func (d LiveDivergence) Key() string {
	return strings.Join([]string{d.Symbol, d.Fixture, string(d.Direction), string(d.Kind), d.Name}, "\x1f")
}

// PairKey is the identity the two runs are matched on, and it excludes the
// direction. It is what makes the comparison possible at all.
//
// If Direction were part of this key too, no key could ever appear on both sides
// and the two "both" buckets would be unreachable rather than merely empty. The
// direction is then not lost: it is what tells the two sides apart once they have
// been matched, and it is what a bucket is derived from.
func (d LiveDivergence) PairKey() string {
	return strings.Join([]string{d.Symbol, d.Fixture, string(d.Kind), d.Name}, "\x1f")
}

// pairKeyOf is the key of a finding without its direction, which is the identity
// the classifier matches the two runs on.
func pairKeyOf(symbol, fixture string, kind DivergenceKind, name string) string {
	return LiveDivergence{Symbol: symbol, Fixture: fixture, Kind: kind, Name: name}.PairKey()
}

// CompareLive runs the documented comparison against both bodies and returns
// every finding, in both directions.
//
// The second run is CompareBody against the committed fixture, which for a row
// whose path matches is what divergence_test.go already runs. It is run again
// here rather than read from the documented baseline because this file makes no
// claim about rows the live tree does not cover, and because re-running it is
// what makes the two result sets comparable: the same checks, the same inputs,
// the same code, differing only in the body.
//
// # Comparability governs one direction and not the other
//
// shapes.go's rule -- a row whose SDK method sends a path the page does not
// document is reported not comparable -- is a rule about the *page*. It exists
// because a schema comparison answers a question that presupposes the two sides
// describe the same HTTP call, and a page describing a different call is not that
// call's contract.
//
// The live run does not presuppose that. A skeleton was captured by calling the
// SDK method itself, so the live body is that method's own response by
// construction: the path question is not in doubt, whatever the page says about
// it. So the live direction runs regardless, and only the documented direction is
// gated.
//
// This is not a corner case. data.GetStockInstruments sends
// /openapi/instrument/stock/list where its page documents
// /trading/instruments/stocks/profiles/list, so the rule marks it not comparable
// and the documented run is rightly silent. But it is one of only four captured
// responses the SDK's own type could not decode at all, and gating the live run
// on the same rule would drop the strongest finding in the whole tree on a
// technicality about which page documents which path. The row is reported, and
// the coverage block names the three rows whose documented side was skipped, so a
// reader can see that the documented half of those rows claims nothing rather
// than having found nothing.
func CompareLive(f Fixture, symbol string, t reflect.Type, skeleton []byte) []LiveDivergence {
	var out []LiveDivergence
	live := CompareBody(f, symbol, t, skeleton)
	for _, d := range live.Divergences {
		out = append(out, liveDivergence(d, LiveSDKDirection, f, symbol, skeleton))
	}
	if compareability(f) != nil {
		return out
	}
	documented, err := f.Read()
	if err != nil {
		return out
	}
	doc := CompareBody(f, symbol, t, documented)
	for _, d := range doc.Divergences {
		out = append(out, liveDivergence(d, LiveDocsDirection, f, symbol, nil))
	}
	return out
}

// liveDivergence maps one of CompareBody's findings onto a LiveDivergence.
//
// The detail is rewritten rather than copied, and the rewrite is deliberately
// source-neutral: it replaces the phrase naming the body with one naming neither.
// CompareBody writes "the documented instance" and "documented top level" because
// it was written for a fixture, and leaving those in place would attribute every
// live observation to a Webull page. But replacing them with "the live response"
// in one direction and nothing in the other would be worse: the two runs would
// then differ in wording for findings that are the *same* finding, and
// ClassifyLive would report every shared divergence as both-detail-changed --
// manufacturing a signal, and drowning the one that matters.
//
// So both directions get the same neutral phrasing, and the two differ only when
// the two bodies genuinely disagree about the fact. That is the whole mechanism
// by which both-detail-changed becomes meaningful rather than an artefact.
func liveDivergence(d Divergence, direction LiveDirection, f Fixture, symbol string, skeleton []byte) LiveDivergence {
	return LiveDivergence{
		Symbol:    symbol,
		Fixture:   f.ID,
		Direction: direction,
		Kind:      d.Kind,
		Name:      d.Name,
		Detail:    nameNearMiss(neutralDetail(d), d, direction, skeleton),
	}
}

// nameNearMiss adds, to a name-level finding in the live direction, the body
// member the server sent that the check could not match.
//
// This is the one place a live finding is materially more informative than a
// documented one, and it is worth the machinery. A required name the live body
// omits is, on its own, an unhelpfully quiet fact: the row says a name is
// missing and stops. But the captured futures tick response is
//
//	{"delay_minutes": -1, "instrumentId": "1", "result": [], "symbol": "1"}
//
// while the SDK tags the field instrument_id. encoding/json matches a member
// exactly and then case-insensitively, and "instrumentId" folds to neither
// "instrument_id" nor "instrumentId" -- the underscore is not a case -- so the
// field is silently left zero. A documented comparison cannot find this at all,
// because the documented instance carries instrument_id and agrees with the SDK.
//
// So when the live body holds a member that differs from the missing name only
// by underscores and case, the detail says so. It is a *report* of what the body
// holds, not a second check: the row exists or does not on CompareBody's verdict
// alone, and this only makes the existing row readable. It is deliberately narrow
// -- a member that differs by any character other than underscore or case is not
// reported, because that is a different finding and naming it here would be a
// check wearing a report's clothes.
func nameNearMiss(detail string, d Divergence, direction LiveDirection, skeleton []byte) string {
	if direction != LiveSDKDirection || d.Name == "" {
		return detail
	}
	if near, ok := nearMissName(skeleton, d.Name); ok {
		return detail + "; the live response sends " + near +
			" instead, which encoding/json matches to this field by neither an exact " +
			"nor a case-insensitive comparison, so the field is left at its zero value"
	}
	return detail
}

// nearMissName returns the member name a body holds that differs from want only
// by underscores and letter case.
func nearMissName(body []byte, want string) (string, bool) {
	var tree any
	if err := decodeNumbered(body, &tree); err != nil {
		return "", false
	}
	var found string
	ok := false
	var walk func(any)
	walk = func(v any) {
		obj, isObject := v.(map[string]any)
		if isObject {
			for name := range obj {
				if !ok && foldsTo(name, want) {
					found, ok = name, true
				}
			}
			for _, member := range obj {
				walk(member)
			}
			return
		}
		if list, isList := v.([]any); isList {
			for _, element := range list {
				walk(element)
			}
		}
	}
	walk(tree)
	return found, ok
}

// foldsTo reports whether two member names differ only by underscores and case.
//
// Underscores are removed and the rest is upper-cased, so "instrumentId" and
// "instrument_id" both become "INSTRUMENTID" and match, while "instrument" and
// "instrument_id" do not. That asymmetry is the point: the underscore is what
// encoding/json will not fold, and a member differing by anything else is a
// different name rather than a spelling of this one.
func foldsTo(got, want string) bool {
	if got == want {
		return false
	}
	return fold(got) == fold(want)
}

// fold is the case-and-underscore-insensitive form of a wire name.
func fold(name string) string {
	var b strings.Builder
	b.Grow(len(name))
	for _, r := range name {
		if r == '_' {
			continue
		}
		if r >= 'a' && r <= 'z' {
			r -= 'a' - 'A'
		}
		b.WriteRune(r)
	}
	return b.String()
}

// neutralDetail restates a CompareBody finding without naming a source.
//
// The substitution is per-kind rather than global, and it has to be: the checks
// do not all read the same inputs. Four of the five read the body they were
// handed, so naming it is accurate for them. checkShape does not -- it reads
// f.Checks.TopLevel, a field the generator recorded from the page -- so a detail
// saying "the compared body's top level is object" would assert something about
// the live body the check never looked at.
//
// That is not a hypothetical. data.GetCapitalFlow's live response is a JSON array
// while its manifest row records the top level as an object, so the live run's
// shape finding would claim the live body is an object while the committed
// skeleton beside it is an array. A uniform substitution produces exactly that
// sentence. So the shape check's wording names what it actually read -- the
// recorded check, which is the page's -- and that is true in both directions.
//
// # Why the decode check is in the same set, and it is not the same fix
//
// The obvious reading of "the shape check does not read the body" is that the
// substitution belongs to the two shape kinds, and that is where it started. The
// decode check quotes the same page fact through a second path, so it needs the
// same substitution, and the reason is worth stating because the two kinds fail
// differently if it is missed.
//
// A shape finding is the *whole* detail, so renaming the word renames the claim:
// the sentence has no other content and cannot be half right. A decode-failure
// detail is a wrapper plus a cause list, and decodeCauses builds the cause for a
// top-level inversion out of f.Checks.TopLevel -- the page's -- while the leaf
// causes are built out of the fixture instance, which in the live run is the live
// body. So the detail mixes a page fact and a body fact, and a substitution that
// rewrites only the word leaves the value attributed to the wrong source: the
// detail reads "the compared body's top level is object" about a body that is an
// array.
//
// The committed file escapes that by coincidence rather than by design.
// data.GetStockInstruments' live body really is an object, and
// data.GetCapitalFlow's decode failure is docs-only, so no recorded row is wrong
// today. That is one kind of body away from being wrong, and this file's whole
// purpose is not to assert a page fact about the wire, so the decode kind is
// rewritten on the same rule as the shape kinds: a container kind in a detail is
// the page's until the wording says it is the body's.
//
// The per-kind wording is longer than a uniform one and it is the right trade: the
// strongest findings in this set are the container disagreements, and the shape
// check is the one that cannot see them, so its sentences are the ones most likely
// to be read as saying more than they do.
//
// Every phrase shapes.go writes for a fixture is replaced, and a test holds each
// detail's stated kind against the fact it can actually be checked against -- so a
// check that gains a new phrase fails visibly rather than being silently
// mis-bucketed, and a rewritten word with an unrewritten value fails too. A stale
// phrase is a misleading sentence a reviewer can read, which is the direction this
// failure has to fail in; a mis-bucketed finding is one nobody can see.
const (
	// neutralDetailPhrase names the body, for the checks that read one.
	neutralDetailPhrase = "the compared body's"
	// recordedTopLevel and recordedElementType name the manifest field the shape
	// check reads, which is the page's and not either body's.
	recordedTopLevel    = "the recorded top level"
	recordedElementType = "the recorded element type"
)

// quotesTheRecordedCheck reports whether a kind's detail states a fact the harness
// read from the page rather than from the body it was handed.
//
// The two shape kinds are the page's own observation. DecodeFailure is included
// because decodeCauses builds its top-level and element-type causes from
// f.Checks, the same field checkShape reads; its leaf causes are built from the
// instance the check was given, and those are handled by the body substitutions
// below rather than by this one.
func quotesTheRecordedCheck(kind DivergenceKind) bool {
	switch kind {
	case TopLevelMismatch, ElementTypeMismatch, DecodeFailure:
		return true
	}
	return false
}

func neutralDetail(d Divergence) string {
	out := d.Detail
	if quotesTheRecordedCheck(d.Kind) {
		out = strings.ReplaceAll(out, "the documented top level", recordedTopLevel)
		out = strings.ReplaceAll(out, "the documented element type", recordedElementType)
		out = strings.ReplaceAll(out, "documented top level", recordedTopLevel)
		out = strings.ReplaceAll(out, "documented element type", recordedElementType)
	}
	// The long phrases are replaced before the bare "documented" that stands alone
	// in the leaf check's "documented %s, SDK field ...", so the word does not fire
	// inside a phrase that was already rewritten and no detail ends up with two
	// articles. A decode-failure detail reaches here too, so the substitution is
	// no longer restricted to the kinds that are wholly body-derived.
	out = strings.ReplaceAll(out, "the documented instance", neutralDetailPhrase+" instance")
	out = strings.ReplaceAll(out, "the documented value", neutralDetailPhrase+" value")
	out = strings.ReplaceAll(out, "the documented", neutralDetailPhrase)
	out = strings.ReplaceAll(out, "documented", neutralDetailPhrase)
	return out
}

// LiveBucket is where a finding falls once the two runs are diffed.
type LiveBucket string

const (
	// LiveBoth means the same finding was observed on both sides with the same
	// detail. It is a pre-existing documentation-versus-SDK disagreement that the
	// live server reproduces, and it is the one bucket that is not new evidence.
	LiveBoth LiveBucket = "both"

	// LiveBothDetailChanged means the same key was observed on both sides with
	// *different* detail. This is the case a naive set difference loses, and
	// losing it is how a live change gets reported as agreement.
	//
	// The two runs compare the same type against two different bodies, so a
	// finding that appears in both is the same disagreement -- but a detail that
	// differs says the two bodies disagree about it. A top-level mismatch whose
	// live detail names an array where the documented one names an object is a
	// server that answers differently from the page: not agreement, and not two
	// independent findings either.
	LiveBothDetailChanged LiveBucket = "both-detail-changed"

	// LiveOnly means the finding was observed against the live wire and not
	// against the documentation. This is the new signal: the server did
	// something the documentation does not describe, which is precisely the class
	// the documented comparison is structurally blind to.
	LiveOnly LiveBucket = "live-only"

	// LiveDocsOnly means the finding was observed against the documentation and
	// not against the live wire. It says the documentation and the SDK disagree
	// where the server and the SDK do not, which is evidence about the page
	// rather than about the SDK.
	LiveDocsOnly LiveBucket = "docs-only"
)

// LiveClassification is one keyed finding and the bucket it falls in.
type LiveClassification struct {
	// Key is the finding's identity, as LiveDivergence.Key builds it.
	Key string `json:"key"`
	// Bucket is where it falls.
	Bucket LiveBucket `json:"bucket"`
	// Symbol, Fixture, Kind and Name are the finding's fields, carried so a
	// baseline entry is readable without splitting the key.
	Symbol  string `json:"symbol"`
	Fixture string `json:"fixture"`
	Kind    string `json:"kind"`
	Name    string `json:"name,omitempty"`
	// Detail is the detail of the observed finding. On LiveBothDetailChanged
	// there are two, and both are recorded.
	Detail string `json:"detail,omitempty"`
	// OtherDetail is the detail from the other side, set only on
	// LiveBothDetailChanged, where the two disagree.
	OtherDetail string `json:"otherDetail,omitempty"`
	// Colliding holds every detail this pair key was observed with beyond the one
	// Detail carries, on either side.
	//
	// It is set only where one side contributed two rows under a single pair key,
	// which the checks this harness runs cannot do: each CompareBody emits at most
	// one finding per (kind, name), and the pair key is exactly that tuple plus the
	// endpoint. So a collision means the caller passed rows this harness did not
	// produce, and the two are recorded rather than one of them quietly kept. An
	// empty list is the normal case and is omitted from the file.
	Colliding []string `json:"colliding,omitempty"`
}

// ClassifyLive diffs the two runs and puts every keyed finding in a bucket.
//
// The comparison is a set difference over *pair keys*, never over slices and
// never over the detail. Four properties follow and each is load-bearing:
//
//   - Ordering is irrelevant. CompareBody returns a slice, the two runs will not
//     produce the same order, and an element-wise comparison would report a
//     different answer for the same evidence. A side that contributed more than one
//     detail for one key is therefore reduced by sorting those details rather than
//     by keeping the last one seen.
//   - The pair key excludes the direction, so a finding both runs observed is one
//     finding with two details. Including the direction would make the two "both"
//     buckets unreachable and every shared divergence read as two unrelated ones.
//   - A key in both runs with the same detail is agreement, and a key in both runs
//     with a different detail is not. Keying on the detail would turn the second
//     case into one live-only and one docs-only finding, which is the specific
//     wrong answer this function exists to avoid; keying on the key alone and
//     ignoring the detail would report it as agreement, which is worse.
//   - The output is sorted by key, so a baseline file is a readable diff rather
//     than a reordering.
//
// # One classification per pair key, and what happens when a key arrives twice
//
// The fourth property -- one classification per pair key, so no finding is
// counted twice and none is dropped -- is a claim about *one row per side*, and
// it was documented as a claim about no row at all, which it did not deliver. A
// map keyed on the pair key keeps the last row it sees, so two rows sharing a
// pair key on one side became one classification with the first detail gone and
// nothing said. This version says so rather than hiding it: every distinct detail
// a pair key was observed with is reported, the classification carries the
// lexicographically first of them so the choice does not depend on row order, and
// the rest are listed in LiveClassification.Colliding.
//
// Detect-and-report was chosen over last-wins for the same reason every other
// check in this file reports rather than tolerates: a finding that disappears
// without saying so is the failure mode this package exists to end, and a caller
// that reads the contract would otherwise rely on it.
//
// # The collision is unreachable through CompareAllLive
//
// Each CompareBody emits at most one finding per (kind, name), so the two rows
// that would collide cannot both come from one run of the checks. That makes the
// property of the *inputs* rather than of this function, and it is why
// CompareAllLive's output cannot contain one -- a test asserts it rather than
// this comment being taken on trust. It is still exported with a contract, and an
// exported contract that is false for a reachable input is worse than a private
// one, because a caller will rely on it.
func ClassifyLive(rows []LiveDivergence) []LiveClassification {
	live := map[string]liveSide{}
	docs := map[string]liveSide{}
	identity := map[string]LiveDivergence{}
	for _, r := range rows {
		key := r.PairKey()
		if _, seen := identity[key]; !seen {
			identity[key] = r
		}
		switch r.Direction {
		case LiveSDKDirection:
			s := live[key]
			s.add(r.Detail)
			live[key] = s
		case LiveDocsDirection:
			s := docs[key]
			s.add(r.Detail)
			docs[key] = s
		}
	}

	keys := make([]string, 0, len(live)+len(docs))
	seen := map[string]bool{}
	for key := range live {
		if !seen[key] {
			keys = append(keys, key)
			seen[key] = true
		}
	}
	for key := range docs {
		if !seen[key] {
			keys = append(keys, key)
			seen[key] = true
		}
	}
	sort.Strings(keys)

	out := make([]LiveClassification, 0, len(keys))
	for _, key := range keys {
		l, inLive := live[key]
		d, inDocs := docs[key]
		lDetail, dDetail := l.kept(), d.kept()
		c := LiveClassification{Key: key}
		if id, ok := identity[key]; ok {
			c.Symbol, c.Fixture, c.Kind, c.Name = id.Symbol, id.Fixture, string(id.Kind), id.Name
		}
		switch {
		case inLive && inDocs && lDetail == dDetail:
			c.Bucket, c.Detail = LiveBoth, lDetail
		case inLive && inDocs:
			// The interesting case: the same disagreement, described differently
			// by the two bodies. Reporting it as LiveBoth would call it agreement;
			// reporting it as two findings would double-count one defect.
			c.Bucket, c.Detail, c.OtherDetail = LiveBothDetailChanged, lDetail, dDetail
		case inLive:
			c.Bucket, c.Detail = LiveOnly, lDetail
		default:
			c.Bucket, c.Detail = LiveDocsOnly, dDetail
		}
		// Whatever the bucket, the details that lost the ordering decision are
		// reported rather than dropped, so the classification accounts for every row
		// it was built from.
		if extra := append(l.unkept(), d.unkept()...); len(extra) > 0 {
			sort.Strings(extra)
			c.Colliding = extra
		}
		out = append(out, c)
	}
	return out
}

// liveSide is every distinct detail one direction contributed for one pair key.
//
// It is a list rather than a single string because a collision has to survive.
// CompareBody emits at most one finding per (kind, name), so a second entry here
// cannot come from the checks this harness runs; when one appears anyway, the
// detail it carries is a finding that would otherwise be lost without a word.
type liveSide struct{ details []string }

// add records a detail, ignoring an exact repeat.
//
// The repeat is ignored rather than doubled so that a caller passing the same row
// twice gets one classification saying one thing, not one saying the same thing
// twice, and so a genuine second *account* of one key is what surfaces.
func (s *liveSide) add(detail string) {
	for _, have := range s.details {
		if have == detail {
			return
		}
	}
	s.details = append(s.details, detail)
}

// kept is the detail the classification carries. It is the lexicographically first
// rather than the first seen, so the answer does not depend on the order rows
// arrived in -- the same property that makes the whole function order-blind.
func (s liveSide) kept() string {
	if len(s.details) == 0 {
		return ""
	}
	sorted := make([]string, len(s.details))
	copy(sorted, s.details)
	sort.Strings(sorted)
	return sorted[0]
}

// unkept is every detail kept did not choose, in the same order-free way.
func (s liveSide) unkept() []string {
	if len(s.details) < 2 {
		return nil
	}
	sorted := make([]string, len(s.details))
	copy(sorted, s.details)
	sort.Strings(sorted)
	return sorted[1:]
}

// LiveBaseline is the committed record of what the live comparison observes, and
// it is separate from known-divergences.json for the reason the file's own
// comment states: the two hold claims about different parties.
//
// A live finding is a claim about Webull's server. Recorded in the SDK's
// baseline it would read as an SDK defect, and the documented backlog plus the
// live backlog would read as one list of 25+N defects when the live rows are not
// claims about this SDK at all.
type LiveBaseline struct {
	// Policy states the gate rule in the file itself, so the rule travels with
	// the data it governs.
	Policy string `json:"policy"`
	// Notes carry the free-text a reader cannot get from a count, chiefly what
	// each bucket means and what the coverage numbers are.
	Notes map[string]string `json:"notes"`
	// Coverage is what was compared, so coverage cannot shrink without the gate
	// noticing.
	Coverage LiveCoverage `json:"coverage"`
	// Entries are the recorded findings, sorted by key.
	Entries []LiveEntry `json:"entries"`
}

// LiveCoverage is the shape of the comparison the baseline was taken against.
//
// It is a committed number rather than a derived one for the reason the documented
// baseline's coverage block is: a probe that reaches fewer endpoints must change
// the file, and a gate that only compared findings would report green while
// checking less.
type LiveCoverage struct {
	// Probed is how many endpoints the live manifest records.
	Probed int `json:"probed"`
	// Compared is how many a symbol and a type were resolved for.
	Compared int `json:"compared"`
	// NotComparable names the rows the page is not the SDK call's contract for.
	// They claim nothing in either direction.
	NotComparable []string `json:"notComparable"`
	// DecodedCleanly is how many live bodies the SDK's own type decoded. A body
	// the type rejects is the strongest signal this harness can produce, and it
	// is counted here so a reader can see how many rows carried one.
	DecodedCleanly int `json:"decodedCleanly"`
	// DecodeRejected names the rows the SDK type could not decode.
	DecodeRejected []string `json:"decodeRejected"`
	// LiveArrayElements is the number of array elements the captured responses
	// held, which is the number the committed tree expands to.
	LiveArrayElements int `json:"liveArrayElements"`
	// TreeBytes is the committed live tree's size, so a re-encoding that changed
	// it would show up as a coverage change rather than passing unnoticed.
	TreeBytes int `json:"treeBytes"`
}

// LiveEntry is one recorded live finding: the observation, plus why it is
// recorded rather than fixed and what would settle it.
type LiveEntry struct {
	LiveDivergence
	// Reason says why the finding is recorded rather than fixed. It is required
	// and must be non-empty; see ParseLiveBaseline.
	Reason string `json:"reason"`
	// Unblock says what would have to happen before the finding could be closed.
	// It is required, and it is required to be specific: "a credential" is not an
	// unblock requirement, a named host and entitlement is.
	Unblock string `json:"unblock"`
	// Bucket records which side the finding was observed on, so a reader can see
	// a server claim apart from a page claim without splitting the key.
	Bucket LiveBucket `json:"bucket"`
}

// liveTreePrefix, liveReadmeFile and liveManifestFile locate the live evidence
// tree and its index inside the committed testdata directory.
//
// The three are spelled out rather than derived, and the prefix ends in a
// separator so a file called "liver" is not in the live tree. They are facts
// about the committed layout, and a fact a consumer re-derives from the tree it
// is reading is a fact the tree can change underneath it.
const (
	liveTreePrefix   = "live/"
	liveReadmeFile   = "README.md"
	liveManifestFile = "live-manifest.json"
)

// LiveBaselineName is the committed file, deliberately a different name from
// known-divergences.json so a reader holding both files cannot mistake one for
// the other, and embedded so the gate reads the bytes that were committed rather
// than whatever is on disk beside the source.
const LiveBaselineName = "live-divergences.json"

//go:embed live-divergences.json
var liveBaselineFS embed.FS

// ParseLiveBaseline reads and validates a live baseline.
//
// The required fields are the point. An entry with an empty Reason records a
// finding nobody has agreed to look at, and one with an empty Unblock records one
// nobody can ever close, and both are how a finding disappears while the file
// still looks maintained. Requiring a reason does not make a reason true, and
// this file has been bitten by a false one: the guard is kept because a reasonless
// entry is unambiguous, and the reasons are held to being per-entry true by review
// rather than by the schema.
func ParseLiveBaseline(raw []byte) (*LiveBaseline, error) {
	var b LiveBaseline
	if err := json.Unmarshal(raw, &b); err != nil {
		return nil, fmt.Errorf("conformance: parse %s: %w", LiveBaselineName, err)
	}
	if len(b.Entries) == 0 {
		return nil, fmt.Errorf("conformance: %s has no entries: a baseline with nothing in "+
			"it is a gate that checks nothing", LiveBaselineName)
	}
	if strings.TrimSpace(b.Policy) == "" {
		return nil, fmt.Errorf("conformance: %s states no policy", LiveBaselineName)
	}
	seen := map[string]bool{}
	prev := ""
	for _, e := range b.Entries {
		where := e.Symbol + " / " + e.Fixture + " / " + string(e.Direction) + " / " + string(e.Kind) + " / " + e.Name
		if strings.TrimSpace(e.Reason) == "" {
			return nil, fmt.Errorf("conformance: %s: %s: entry has no reason, which is a "+
				"suppressed finding rather than a recorded one", LiveBaselineName, where)
		}
		if strings.TrimSpace(e.Unblock) == "" {
			return nil, fmt.Errorf("conformance: %s: %s: entry names no unblock requirement, "+
				"so the finding can never be closed", LiveBaselineName, where)
		}
		if strings.TrimSpace(e.Detail) == "" {
			return nil, fmt.Errorf("conformance: %s: %s: entry has no detail, so nothing "+
				"says what disagrees", LiveBaselineName, where)
		}
		switch e.Direction {
		case LiveSDKDirection, LiveDocsDirection:
		default:
			return nil, fmt.Errorf("conformance: %s: %s: entry has an unknown direction %q",
				LiveBaselineName, where, e.Direction)
		}
		if !knownLiveKind(e.Kind) {
			return nil, fmt.Errorf("conformance: %s: %s: entry has a kind the comparison "+
				"cannot produce, so the entry is unresolvable", LiveBaselineName, where)
		}
		if prev != "" && e.Key() < prev {
			return nil, fmt.Errorf("conformance: %s: entries are not sorted: %q comes after %q",
				LiveBaselineName, e.Key(), prev)
		}
		prev = e.Key()
		if seen[e.Key()] {
			return nil, fmt.Errorf("conformance: %s: %s is recorded twice, so it reads as "+
				"two independent observations of one fact", LiveBaselineName, e.Key())
		}
		seen[e.Key()] = true
	}
	return &b, nil
}

// knownLiveKind reports whether the comparison can produce this kind. It is the
// same set shapes.go can emit, and a kind outside it would make the entry
// unresolvable rather than wrong -- which is the worse of the two for a gate.
func knownLiveKind(k DivergenceKind) bool {
	for _, candidate := range []DivergenceKind{
		MissingRequiredName, MissingDeclaredName, DeclaredInventoryEmpty,
		TopLevelMismatch, ElementTypeMismatch, LeafTypeMismatch, DecodeFailure,
	} {
		if k == candidate {
			return true
		}
	}
	return false
}

// Marshal renders the baseline the way it is committed: entries sorted by key,
// two-space indent, one trailing newline, so a regeneration is a readable diff
// rather than a reformatting.
func (b *LiveBaseline) Marshal() ([]byte, error) {
	sorted := &LiveBaseline{Policy: b.Policy, Notes: b.Notes, Coverage: b.Coverage}
	sorted.Entries = append([]LiveEntry(nil), b.Entries...)
	sort.Slice(sorted.Entries, func(i, j int) bool {
		return sorted.Entries[i].Key() < sorted.Entries[j].Key()
	})
	out, err := json.MarshalIndent(sorted, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("conformance: render %s: %w", LiveBaselineName, err)
	}
	return append(out, '\n'), nil
}

// LoadLiveBaseline reads the committed live baseline from the embedded tree, so
// the gate reads the bytes that were committed and not whatever is on disk beside
// the source.
func LoadLiveBaseline() (*LiveBaseline, error) {
	raw, err := liveBaselineFS.ReadFile(LiveBaselineName)
	if err != nil {
		return nil, fmt.Errorf("conformance: read %s: %w", LiveBaselineName, err)
	}
	return ParseLiveBaseline(raw)
}

// LiveRun is the whole live comparison: what it covered, what it found, and what
// it could not decide.
type LiveRun struct {
	// Rows is every finding from both runs, unclassified.
	Rows []LiveDivergence
	// Classifications is Rows diffed into buckets, sorted by key.
	Classifications []LiveClassification
	// Coverage is what the comparison reached, for the baseline's coverage block.
	Coverage LiveCoverage
	// Errs are the endpoints no comparison could be made for, each naming why.
	// They are a list rather than a count because a dropped endpoint is the
	// failure mode this instrument exists to end, and a count hides which.
	Errs []error
}

// CompareAllLive runs the live comparison for every endpoint the live tree
// covers, and returns the findings, the classification and the coverage.
//
// It joins the two trees on the fixture name each live manifest entry records,
// resolves a Go type through the same hand-written table the documented
// comparison uses, and runs CompareLive per row. A row it cannot compare is
// collected in Errs and is *not* silently skipped: an endpoint nobody examined
// and an endpoint that examined clean are the two things a reader must be able to
// tell apart, and the only way to tell them apart is for the difference to be
// visible.
func CompareAllLive() LiveRun {
	run := LiveRun{}

	manifest, err := Load()
	if err != nil {
		run.Errs = append(run.Errs, err)
		return run
	}
	byFixture := make(map[string]Fixture, len(manifest.Fixtures))
	for _, f := range manifest.Fixtures {
		byFixture[f.Fixture] = f
	}

	// The entries are read rather than LoadLive's map, because the loop needs the
	// committed bytes per row to survey them, and loading the tree a second time
	// only to iterate it would expand every skeleton twice for one result.
	entries, err := liveEntries(fixturesFS, testdataDir+"/"+liveManifestFile)
	if err != nil {
		run.Errs = append(run.Errs, err)
		return run
	}

	run.Coverage.Probed = len(entries)
	run.Coverage.Compared = 0
	distinctShapes := map[string]bool{}
	for _, e := range entries {
		if e.DecodedCleanly {
			run.Coverage.DecodedCleanly++
		} else {
			run.Coverage.DecodeRejected = append(run.Coverage.DecodeRejected, e.Symbol)
		}

		raw, err := fixturesFS.ReadFile(testdataDir + "/" + e.Skeleton)
		if err == nil {
			survey := SurveySkeleton(raw)
			run.Coverage.LiveArrayElements += survey.ArrayElements
			run.Coverage.TreeBytes += len(raw)
			for shape := range survey.DistinctElementShapes {
				distinctShapes[shape] = true
			}
		}
		if err != nil {
			run.Errs = append(run.Errs, err)
			continue
		}

		f, ok := byFixture[e.Fixture]
		if !ok {
			run.Errs = append(run.Errs, fmt.Errorf("conformance: %s: the live manifest names "+
				"fixture %q, which the documentation manifest does not", e.Symbol, e.Fixture))
			continue
		}
		if compareability(f) != nil {
			// The live direction still runs; only the documented half is skipped.
			// Recording the row makes a reader see that its documented side
			// claims nothing rather than having found nothing.
			run.Coverage.NotComparable = append(run.Coverage.NotComparable, e.Symbol)
		}
		subject, _, named := Subject(f.SDKSymbol)
		if !named || subject == NoSDKSymbol {
			run.Errs = append(run.Errs, fmt.Errorf("conformance: %s: fixture %s names no SDK symbol",
				e.Symbol, e.Fixture))
			continue
		}
		entry, found := SDKTypes[subject]
		if !found {
			run.Errs = append(run.Errs, fmt.Errorf("conformance: %s: no decode target in the "+
				"symbol table", subject))
			continue
		}
		target, err := entry.DecodeTarget()
		if err != nil {
			run.Errs = append(run.Errs, err)
			continue
		}
		body, err := ExpandSkeleton(raw)
		if err != nil {
			run.Errs = append(run.Errs, err)
			continue
		}
		run.Coverage.Compared++
		run.Rows = append(run.Rows, CompareLive(f, subject, target, body)...)
	}

	sort.Strings(run.Coverage.NotComparable)
	sort.Strings(run.Coverage.DecodeRejected)
	run.Classifications = ClassifyLive(run.Rows)
	return run
}

// BucketCounts is how many findings fell in each bucket, for the coverage block
// and for a report. Every bucket is present even at zero, so a reader can see
// that a bucket was counted and found empty rather than not looked at.
func BucketCounts(classifications []LiveClassification) map[LiveBucket]int {
	counts := map[LiveBucket]int{
		LiveBoth: 0, LiveBothDetailChanged: 0, LiveOnly: 0, LiveDocsOnly: 0,
	}
	for _, c := range classifications {
		counts[c.Bucket]++
	}
	return counts
}

// LoadLive reads the committed live tree and returns the expanded body of every
// skeleton it holds, keyed by the SDK symbol that produced it.
//
// The bodies are expanded, so a caller feeds CompareLive the same reduced tree
// the probe wrote and the comparison reads typed placeholders. That is what makes
// a skeleton a valid comparison input at all: jsonKind classifies a decoded value
// by its Go type, so a leaf spelled as a kind string would be demanded of a Go
// string and a numeric disagreement would be reported as agreement.
//
// An empty dir reads the embedded tree, which is what every test and the gate
// want, and a non-empty dir reads from disk for a caller regenerating the
// findings. An error is returned rather than an empty map, because a comparison
// that reached nothing and reported nothing would be indistinguishable from a
// clean run.
func LoadLive(dir string) (map[string]Skeleton, error) {
	if dir != "" {
		onDisk := osDirFS(dir)
		entries, err := liveEntries(onDisk, testdataDir+"/"+liveManifestFile)
		if err != nil {
			return nil, err
		}
		return loadLiveSkeletons(onDisk, entries, testdataDir+"/"+liveTreePrefix)
	}
	entries, err := liveEntries(fixturesFS, testdataDir+"/"+liveManifestFile)
	if err != nil {
		return nil, err
	}
	return loadLiveSkeletons(fixturesFS, entries, testdataDir+"/"+liveTreePrefix)
}

// Skeleton is one captured response: which symbol produced it, which documented
// endpoint it belongs to, and its body as the comparison reads it.
type Skeleton struct {
	// Symbol is the SDK method the probe called to obtain the response.
	Symbol string
	// Fixture is the documentation fixture ID for the same endpoint, which is how
	// the two trees join row for row.
	Fixture string
	// Body is the expanded skeleton, ready for CompareLive.
	Body []byte
}

// liveEntry is the part of a live manifest entry this package reads. It is a
// minimal struct rather than a type mirroring the file's, because the file is
// written by an example program and this package has no business depending on
// the shape of a field it does not read.
type liveEntry struct {
	Symbol         string `json:"symbol"`
	Fixture        string `json:"fixture"`
	Skeleton       string `json:"skeleton"`
	DecodedCleanly bool   `json:"decodedCleanly"`
}

// liveEntries reads the live manifest's entries from an fs.
func liveEntries(fsys fs.FS, path string) ([]liveEntry, error) {
	raw, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("conformance: read %s: %w", path, err)
	}
	var doc struct {
		Entries []liveEntry `json:"entries"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("conformance: parse %s: %w", path, err)
	}
	if len(doc.Entries) == 0 {
		return nil, fmt.Errorf("conformance: %s records no entries, so the live tree is "+
			"unindexed and nothing could be compared against it", path)
	}
	return doc.Entries, nil
}

// loadLiveSkeletons reads and expands every skeleton an entry names.
func loadLiveSkeletons(fsys fs.FS, entries []liveEntry, root string) (map[string]Skeleton, error) {
	out := make(map[string]Skeleton, len(entries))
	for _, e := range entries {
		rel, ok := strings.CutPrefix(e.Skeleton, liveTreePrefix)
		if !ok {
			return nil, fmt.Errorf("conformance: %s names skeleton %q, which is not under %s",
				liveManifestFile, e.Skeleton, liveTreePrefix)
		}
		raw, err := fs.ReadFile(fsys, path.Join(root, rel))
		if err != nil {
			return nil, fmt.Errorf("conformance: read skeleton %s: %w", e.Skeleton, err)
		}
		body, err := ExpandSkeleton(raw)
		if err != nil {
			return nil, fmt.Errorf("conformance: %s: %w", e.Skeleton, err)
		}
		out[e.Symbol] = Skeleton{Symbol: e.Symbol, Fixture: e.Fixture, Body: body}
	}
	return out, nil
}

// osDirFS adapts a directory on disk to the reader interface loadLiveSkeletons
// uses, so the embedded and the on-disk path share one implementation rather than
// two. The path is built from the manifest's own names and confined to the tree
// root, which is why LoadLive is not a package-level var over a bare path.
func osDirFS(dir string) fs.FS { return os.DirFS(dir) }
