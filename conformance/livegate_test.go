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
	"reflect"
	"sort"
	"strings"
	"testing"
)

// TestLiveDivergenceBaseline is the gate: the classified set must exactly equal
// the recorded set, and it must fail in both directions.
//
// A new finding is not recorded, which is the ordinary case and the reason the
// harness is worth running. A recorded finding that no longer reproduces also
// fails, and that is the direction this design exists for: a finding that was
// fixed and a check that changed meaning look identical from here, and a gate
// that reported only the first would have told a reader a defect was closed when
// the comparison had merely stopped looking.
//
// The coverage assertions are here for the same reason the documented baseline
// carries a coverage block: a gate that only compared findings would report green
// while checking less. If a probe reached fewer endpoints the counts would move
// and this would fail, rather than the gate quietly shrinking.
func TestLiveDivergenceBaseline(t *testing.T) {
	baseline, err := LoadLiveBaseline()
	if err != nil {
		t.Fatalf("load %s: %v", LiveBaselineName, err)
	}
	run := CompareAllLive()
	for _, e := range run.Errs {
		t.Errorf("the live comparison could not cover everything: %v", e)
	}

	assertLiveCoverage(t, baseline, run)

	// The gate matches on the pair key, not the record key. A classification is
	// one finding observed on one or both sides, so that is the identity a
	// recorded entry stands for; LiveDivergence.Key additionally carries the
	// direction, which is what a *single observation* is keyed by and what makes
	// the two sides of one finding two records rather than one.
	observed := make(map[string]LiveClassification, len(run.Classifications))
	for _, c := range run.Classifications {
		observed[c.Key] = c
	}
	recorded := make(map[string]LiveEntry, len(baseline.Entries))
	for _, e := range baseline.Entries {
		recorded[e.PairKey()] = e
	}

	for key, c := range observed {
		entry, ok := recorded[key]
		if !ok {
			t.Errorf("a live finding is not recorded in %s: [%s] %s / %s / %q -- %s\n"+
				"Add it with a reason and an unblock requirement, or fix the SDK; a finding "+
				"nobody has agreed to look at is a suppressed one",
				LiveBaselineName, c.Bucket, c.Symbol, c.Kind, c.Name, c.Detail)
			continue
		}
		if entry.Detail != c.Detail {
			t.Errorf("%s: the recorded detail no longer matches what the comparison "+
				"reports for %s / %s / %q.\nrecorded: %s\nobserved: %s",
				LiveBaselineName, c.Symbol, c.Kind, c.Name, entry.Detail, c.Detail)
		}
		if entry.Bucket != c.Bucket {
			t.Errorf("%s: %s / %s / %q is recorded in bucket %q and observed in %q",
				LiveBaselineName, c.Symbol, c.Kind, c.Name, entry.Bucket, c.Bucket)
		}
	}
	for key, e := range recorded {
		if _, ok := observed[key]; !ok {
			t.Errorf("%s records %s / %s / %q in bucket %q and the comparison no longer "+
				"reproduces it.\nEither the finding is fixed and the entry should be removed, "+
				"or a check changed meaning and this needs a person to look at the tree; the "+
				"gate cannot tell which",
				LiveBaselineName, e.Symbol, e.Kind, e.Name, e.Bucket)
		}
	}
	// A recorded entry carries the direction its finding was observed on, and the
	// bucket that direction produced. Both are checked above against the
	// classification, so a file that got either wrong fails rather than being
	// matched on a key it happens to share.
}

// assertLiveCoverage holds the baseline's coverage block to the run that produced
// it.
func assertLiveCoverage(t *testing.T, baseline *LiveBaseline, run LiveRun) {
	t.Helper()
	got, want := run.Coverage, baseline.Coverage
	if got.Probed != want.Probed {
		t.Errorf("coverage.probed = %d, the baseline records %d", got.Probed, want.Probed)
	}
	if got.Compared != want.Compared {
		t.Errorf("coverage.compared = %d, the baseline records %d", got.Compared, want.Compared)
	}
	if got.DecodedCleanly != want.DecodedCleanly {
		t.Errorf("coverage.decodedCleanly = %d, the baseline records %d",
			got.DecodedCleanly, want.DecodedCleanly)
	}
	if got.LiveArrayElements != want.LiveArrayElements {
		t.Errorf("coverage.liveArrayElements = %d, the baseline records %d; the committed "+
			"tree no longer expands to the number of array elements the capture observed",
			got.LiveArrayElements, want.LiveArrayElements)
	}
	if got.TreeBytes != want.TreeBytes {
		t.Errorf("coverage.treeBytes = %d, the baseline records %d", got.TreeBytes, want.TreeBytes)
	}
	if !reflect.DeepEqual(got.NotComparable, want.NotComparable) {
		t.Errorf("coverage.notComparable = %v, the baseline records %v",
			got.NotComparable, want.NotComparable)
	}
	if !reflect.DeepEqual(got.DecodeRejected, want.DecodeRejected) {
		t.Errorf("coverage.decodeRejected = %v, the baseline records %v",
			got.DecodeRejected, want.DecodeRejected)
	}
}

// TestLiveBaselineRecordsEveryBucketAndItsCounts is the counting claim.
//
// The buckets are the deliverable -- "of N probed, Y disagree with the SDK",
// split by which side disagrees -- and a number that is computed from the
// classifications rather than transcribed cannot drift from them.
//
// The way to break this test is to compute a count from somewhere other than the
// classifications, which is how a report and a file come to disagree.
func TestLiveBaselineRecordsEveryBucketAndItsCounts(t *testing.T) {
	baseline, err := LoadLiveBaseline()
	if err != nil {
		t.Fatalf("load %s: %v", LiveBaselineName, err)
	}
	run := CompareAllLive()
	observed := BucketCounts(run.Classifications)

	recorded := map[LiveBucket]int{}
	for _, e := range baseline.Entries {
		recorded[e.Bucket]++
	}
	for _, bucket := range []LiveBucket{LiveBoth, LiveBothDetailChanged, LiveOnly, LiveDocsOnly} {
		if recorded[bucket] != observed[bucket] {
			t.Errorf("bucket %q holds %d recorded entr(ies) and %d observed finding(s)",
				bucket, recorded[bucket], observed[bucket])
		}
	}
	// A bucket nothing fell into is still a bucket the reader can count, so the
	// notes must name it rather than leaving it out of the file. The note keys are
	// camelCase because that is what the notes above spell them as, and a mismatch
	// here means a note exists that no reader will look for.
	noteKey := map[LiveBucket]string{
		LiveBoth:              "both",
		LiveBothDetailChanged: "bothDetailChanged",
		LiveOnly:              "liveOnly",
		LiveDocsOnly:          "docsOnly",
	}
	for _, bucket := range []LiveBucket{LiveBoth, LiveBothDetailChanged, LiveOnly, LiveDocsOnly} {
		key := noteKey[bucket]
		if _, ok := baseline.Notes[key]; !ok {
			t.Errorf("the baseline has no note for bucket %q (looked for the key %q), so a "+
				"reader cannot tell an empty bucket from one nobody counted", bucket, key)
		}
	}
}

// TestLiveBaselineNamesTheFourDecodeRejections is the payload assertion.
//
// The four bodies the SDK's own type could not unmarshal are the strongest signal
// this harness can produce, and the finding that each is a container-kind
// disagreement with the server is only worth anything if the four stay visible. A
// finding that is not in the recorded set is a finding nobody has agreed to look
// at, so this asserts each one has an entry rather than trusting that it happens
// to be there.
//
// The way to break this test is a filter that drops them: the comparability rule
// applied to the live direction drops data.GetStockInstruments, which is one of
// the four, and a filter that looked like this one would hide it without saying so.
func TestLiveBaselineNamesTheFourDecodeRejections(t *testing.T) {
	manifest := readLiveManifest(t)
	rejected := manifest.decodeRejected()
	if len(rejected) == 0 {
		t.Fatal("the live manifest records no decode rejections, so this test is asserting nothing")
	}

	baseline, err := LoadLiveBaseline()
	if err != nil {
		t.Fatalf("load %s: %v", LiveBaselineName, err)
	}
	recorded := map[string]bool{}
	for _, e := range baseline.Entries {
		recorded[e.Symbol] = true
	}
	for _, symbol := range rejected {
		if !recorded[symbol] {
			t.Errorf("%s: the live manifest records that %s's body did not decode into the "+
				"SDK's own type, and no live finding is recorded for it. That is the strongest "+
				"signal this harness produces, and the documented comparison cannot see it",
				LiveBaselineName, symbol)
		}
	}
	if !reflect.DeepEqual(rejected, baseline.Coverage.DecodeRejected) {
		t.Errorf("the manifest records decode rejections %v and the baseline records %v",
			rejected, baseline.Coverage.DecodeRejected)
	}
}

// TestLiveFindingsAreAttributable is the honesty check on the recorded set.
//
// A live finding is a claim about Webull's server, and the whole reason this set
// is a separate file is that the claims have different authors and different
// unblock requirements. A recorded entry whose reason and unblock are the same
// sentence repeated is the documented baseline's recorded failure -- a false
// reason satisfies a required-field check perfectly -- so the reasons are held to
// being per-entry distinct here.
//
// The way to break this test is a template: one reason and one unblock formatted
// per finding, which reads as coverage and is not.
func TestLiveFindingsAreAttributable(t *testing.T) {
	baseline, err := LoadLiveBaseline()
	if err != nil {
		t.Fatalf("load %s: %v", LiveBaselineName, err)
	}
	reasons := map[string]string{}
	unblocks := map[string]string{}
	for _, e := range baseline.Entries {
		where := e.Symbol + " / " + string(e.Kind) + " / " + e.Name
		if prev, dup := reasons[e.Reason]; dup {
			t.Errorf("%s: the reason for %s is byte-identical to the one for %s; a reason "+
				"repeated across entries satisfies a required-field check perfectly while "+
				"saying nothing about this one", LiveBaselineName, where, prev)
		}
		reasons[e.Reason] = where

		if prev, dup := unblocks[e.Unblock]; dup {
			// A cross-reference is a legitimate way to close several findings, so
			// this is a warning about the shape of the entry rather than a failure:
			// the repeated unblock must say which finding it defers to, or it is a
			// placeholder that points at nothing.
			if !strings.Contains(e.Unblock, prev) {
				t.Errorf("%s: the unblock requirement for %s is byte-identical to the "+
					"one for %s and does not say which finding it defers to; an unblock "+
					"requirement that is the same sentence everywhere is a placeholder "+
					"rather than a requirement",
					LiveBaselineName, where, prev)
			}
		}
		unblocks[e.Unblock] = where

		if len(e.Unblock) < 40 {
			t.Errorf("%s: %s: the unblock requirement is %q, which is too short to name a "+
				"host, an entitlement or a decision; \"a credential\" is not a requirement",
				LiveBaselineName, where, e.Unblock)
		}
	}
}

// TestLiveBaselineDetailsNameWhatTheyRead holds every recorded detail to the claim
// it makes, rather than to the absence of a word.
//
// The first version of this test grepped for "documented" and stopped there, and
// it passed over a detail asserting a page fact about a live body. The word was
// gone and the claim was not: a decode-failure detail is a wrapper plus a cause
// list, and decodeCauses builds the cause for a container inversion out of
// f.Checks -- the page -- while a uniform substitution rewrote only the word
// around it. The detail read "the compared body's top level is object" about a
// body that was an array. The committed file escaped by coincidence:
// data.GetStockInstruments' live body really is an object and
// data.GetCapitalFlow's decode failure is docs-only, so no recorded row was wrong.
// One kind of body away it would have been, and this package's whole purpose is
// not to assert a page fact about the wire.
//
// So the check is on the claim. Every container kind a detail states is held
// against the fact it can be checked against: a "recorded" statement against the
// page's own Checks, and a "the compared body's" statement against the kind of
// the body the finding was observed on. A kind stated about the body that is not
// the body's kind fails, which is what a renamed word cannot hide.
//
// The two halves are different failures and both are kept. The word check catches
// a phrase the substitution list has not caught up with, which is a stale sentence
// a reviewer would read; the claim check catches a substituted word wrapping an
// unsubstituted value, which is a false statement. Neither implies the other, so
// a test that only did one would leave the other's failure invisible.
//
// The way to break this test is to leave a kind out of quotesTheRecordedCheck, or
// to re-introduce an early return that skips the body substitutions for a kind
// that quotes the page.
func TestLiveBaselineDetailsNameWhatTheyRead(t *testing.T) {
	baseline, err := LoadLiveBaseline()
	if err != nil {
		t.Fatalf("load %s: %v", LiveBaselineName, err)
	}
	byID := fixturesByID(t)
	live, err := LoadLive("")
	if err != nil {
		t.Fatalf("LoadLive: %v", err)
	}

	for _, e := range baseline.Entries {
		where := e.Symbol + " / " + string(e.Kind) + " / " + e.Name
		if strings.Contains(e.Detail, "documented") {
			t.Errorf("%s: %s has a detail naming the documentation: %q\n"+
				"Both runs are compared against a body and neither body is the "+
				"documentation; a detail that names it attributes a server observation "+
				"to a Webull page",
				LiveBaselineName, where, e.Detail)
		}
		f, ok := byID[e.Fixture]
		if !ok {
			t.Errorf("%s: %s names fixture %q, which the documentation manifest does "+
				"not, so its detail cannot be held to anything", LiveBaselineName, where, e.Fixture)
			continue
		}
		var body []byte
		if e.Direction == LiveSDKDirection {
			sk, ok := live[e.Symbol]
			if !ok {
				t.Errorf("%s: %s is a live finding and the live tree holds no skeleton "+
					"for it, so its detail cannot be held to the body it describes",
					LiveBaselineName, where)
				continue
			}
			body = sk.Body
		} else {
			raw, err := f.Read()
			if err != nil {
				t.Errorf("%s: %s: read the documented fixture: %v", LiveBaselineName, where, err)
				continue
			}
			body = raw
		}
		if !assertDetailKindsAreTrue(t, f, body, e.Detail) {
			t.Errorf("%s: %s states a kind the evidence contradicts: %q", LiveBaselineName, where, e.Detail)
		}
	}
}

// fixturesByID indexes the documentation manifest by fixture id, which is what a
// recorded live entry names.
func fixturesByID(t *testing.T) map[string]Fixture {
	t.Helper()
	m := loadManifest(t)
	byID := make(map[string]Fixture, len(m.Fixtures))
	for _, f := range m.Fixtures {
		byID[f.ID] = f
	}
	return byID
}

// containerKindPhrases are every way a detail can state a JSON kind about
// something, and what each one has to be true of.
//
// A phrase about the recorded check is the page's, so it is checked against the
// page. A phrase about the compared body is that body's, so it is checked against
// the body -- and that is the assertion that catches a container kind attributed
// to a body the check never read.
var containerKindPhrases = []struct {
	// phrase is the text as it appears in a detail, subject included.
	phrase string
	// aspect is which of the container kinds the phrase is about.
	aspect string
	// recorded reports whether the phrase is about the page's recorded check.
	recorded bool
}{
	{recordedTopLevel, "top level", true},
	{recordedElementType, "element type", true},
	{neutralDetailPhrase + " top level", "top level", false},
	{neutralDetailPhrase + " element type", "element type", false},
}

// The kinds a container phrase can state. checkShape and decodeCauses both write
// one of jsonKind's outputs, and a kind outside this set means a phrase the
// checker cannot read -- which is a failure, not a skip.
var containerKinds = []string{"array", "boolean", "null", "number", "object", "string"}

// assertDetailKindsAreTrue holds every container kind a detail states to the
// fact it can be checked against, and reports whether the detail was fully
// accounted for.
//
// It is one function rather than two assertions so the committed baseline and the
// driven case below are held to exactly the same rule; a second copy of the rule
// would be a second thing to keep true.
//
// Every occurrence of "top level" or "element type" has to be part of a phrase
// this knows and a kind it can read. That is what makes a new phrase in
// shapes.go fail here rather than pass unnoticed: a detail the checker cannot
// account for is reported, so adding a fifth way to state a kind does not quietly
// widen what the baseline is allowed to say.
//
// The leaf kinds a detail may quote are deliberately not checked. Those come from
// documentedNames, which reads the fixture instance the check was handed, so in
// the live run they are read from the live body and are the body's to state. The
// container kinds are the ones the page supplies, which is the whole of the axis
// this guards.
func assertDetailKindsAreTrue(t liveValueSink, f Fixture, body []byte, detail string) bool {
	t.Helper()
	ok := true
	accounted := detail
	for _, p := range containerKindPhrases {
		for {
			at := strings.Index(accounted, p.phrase)
			if at < 0 {
				break
			}
			rest := accounted[at+len(p.phrase):]
			kind, width := leadingKind(rest)
			if kind == "" {
				t.Errorf("the detail states %q and then no JSON kind this checker can "+
					"read: %q", p.phrase, detail)
				ok = false
				break
			}
			want := f.Checks.TopLevel
			if p.aspect == "element type" {
				want = f.Checks.ElementType
			}
			if p.recorded {
				// The page's own record. Nothing to compare against but the page,
				// and the value came from there, so this asserts the detail quotes
				// the field it names rather than a body it never read.
				if kind != want {
					t.Errorf("the detail says the recorded %s is %q and the manifest "+
						"records %q", p.aspect, kind, want)
					ok = false
				}
			} else if got := bodyKind(body); kind != got {
				t.Errorf("the detail says the compared body's %s is %q and the body it "+
					"was observed on is %q: the sentence attributes to the wire a fact "+
					"the check read from the page", p.aspect, kind, got)
				ok = false
			}
			accounted = rest[width:]
		}
	}
	for _, aspect := range []string{"top level", "element type"} {
		if strings.Contains(accounted, aspect) {
			t.Errorf("the detail mentions the %s in a form this checker does not "+
				"account for, so the claim cannot be held to anything: %q", aspect, detail)
			ok = false
		}
	}
	return ok
}

// TestDetailKindClaimCheckerBites is the checker held to its own claims, because a
// rule that silently accepts everything would make the two tests that use it
// worthless.
//
// The cases that matter are the last four. A detail that states the page's kind
// correctly passes; one that states it wrongly fails; one that states the body's
// kind correctly passes; one that states it wrongly fails -- and that last pair is
// the one this checker exists for, because a kind attributed to the wire that came
// from the page is exactly the defect the rewrite fixed. The fifth case is the
// guard on the guard: a phrasing the checker does not know is reported rather than
// passed, so adding a fifth way for shapes.go to state a kind fails here instead of
// quietly widening what the baseline may say.
//
// The way to break this test is to make the checker lenient -- skip the body
// comparison, or stop reporting an unaccounted mention -- at which point the last
// four cases pass for the wrong reason.
func TestDetailKindClaimCheckerBites(t *testing.T) {
	f := liveFixture("test/GET-example", "object", nil, nil)
	f.Checks.ElementType = "object"
	array := []byte(`[{"a":"1"}]`)
	object := []byte(`{"a":"1"}`)

	cases := []struct {
		name        string
		body        []byte
		detail      string
		wantReports int
		because     string
	}{
		{
			name:   "a detail that states no kind",
			body:   array,
			detail: "the compared body's instance omits a name the page requires",
			because: "most findings say nothing about a container kind, and the checker " +
				"must not report those or every name-level row would fail",
		},
		{
			name:   "a leaf kind, which the body supplies",
			body:   object,
			detail: "the compared body's string, SDK field P.A is a number",
			because: "a leaf kind is read from the instance the check was handed, so it is " +
				"the body's to state and there is nothing to hold it to",
		},
		{
			name:   "the recorded top level, stated as the manifest records it",
			body:   array,
			detail: "the recorded top level is object, []row decodes it as array",
			because: "the page records an object and the detail says so; this is what the " +
				"rewrite is for and it must pass",
		},
		{
			name:        "the recorded top level, stated wrongly",
			body:        array,
			detail:      "the recorded top level is array, []row decodes it as array",
			wantReports: 1,
			because: "the page records an object, so a detail quoting an array is false " +
				"about the page whether or not the body agrees",
		},
		{
			name:   "the body's top level, stated as the body is",
			body:   array,
			detail: "the compared body's top level is array, []row decodes it as array",
			because: "true of the body, and the checker has to be able to say so rather " +
				"than refusing every body-phrased claim",
		},
		{
			name:        "the body's top level, stated as the page records it",
			body:        array,
			detail:      "the compared body's top level is object, []row decodes it as array",
			wantReports: 1,
			because: "this is the exact sentence the old substitution produced, and the " +
				"body is an array, so the claim is false",
		},
		{
			name:        "a phrasing the checker does not know",
			body:        array,
			detail:      "the wire's top level is array, []row decodes it as array",
			wantReports: 1,
			because: "a new phrasing has to be reported rather than passed, or a check " +
				"that grows one widens the baseline silently",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := &recordingSink{}
			if got := assertDetailKindsAreTrue(rec, f, tc.body, tc.detail); got != (tc.wantReports == 0) {
				t.Errorf("assertDetailKindsAreTrue returned %t with %d report(s), want %t "+
					"with %d: %v; %s", got, len(rec.messages), tc.wantReports == 0,
					tc.wantReports, rec.messages, tc.because)
			}
			if got := len(rec.messages); got != tc.wantReports {
				t.Errorf("assertDetailKindsAreTrue reported %d time(s), want %d: %v; %s",
					got, tc.wantReports, rec.messages, tc.because)
			}
		})
	}
}

// leadingKind reads the JSON kind a container phrase is followed by, and how many
// bytes of rest it took. shapes.go writes " is <kind>" immediately after the
// phrase, so the word is bounded by the next space or comma.
func leadingKind(rest string) (kind string, width int) {
	trimmed := strings.TrimLeft(rest, " ")
	consumed := len(rest) - len(trimmed) + len("is ")
	if !strings.HasPrefix(trimmed, "is ") {
		return "", 0
	}
	trimmed = trimmed[len("is "):]
	for _, candidate := range containerKinds {
		if !strings.HasPrefix(trimmed, candidate) {
			continue
		}
		after := trimmed[len(candidate):]
		if after == "" || after[0] == ' ' || after[0] == ',' {
			return candidate, consumed + len(candidate)
		}
	}
	return "", 0
}

// bodyKind is the JSON kind of a compared body's top level, as the check that read
// it would have classified it.
func bodyKind(body []byte) string {
	var tree any
	if err := decodeNumbered(body, &tree); err != nil {
		return "unreadable"
	}
	return jsonKind(tree)
}

// TestLiveDecodeDetailDoesNotAssertAPageFactAboutTheWire is the case the word
// check could not see, driven end to end through the real pipeline.
//
// The page records an object, the SDK decodes an array, and the server sent an
// array -- so the two facts differ, and the decode detail has to say which one it
// is quoting. Every kind of body away from the two the committed tree happens to
// hold, and the substitution that keeps the word out of the detail is exactly what
// would put the page's value inside the body's claim.
//
// The body is an array whose element does not decode, so the failure is a genuine
// rejection and the cause list is the top-level one rather than the fallback
// clause. The detail is then held to the same rule the committed baseline is held
// to, against the same page record and the same body, so this is not a word check
// wearing a new hat: it asserts that the kind the detail states about the body is
// the kind the body has.
//
// The way to break this test is to drop DecodeFailure from
// quotesTheRecordedCheck, which leaves the detail saying the body's top level is
// an object while the body is an array. No test in the repository fails on that:
// the detail is prose, the bucket is unchanged, the counts are unchanged and the
// key is unchanged.
func TestLiveDecodeDetailDoesNotAssertAPageFactAboutTheWire(t *testing.T) {
	type row struct {
		Close float64 `json:"close"`
	}
	// The page records an object; the SDK decodes a slice, so the shape check
	// reports an inversion and that report becomes the decode check's cause.
	f := liveFixture("test/GET-example", "object", nil, nil)
	if recorded := f.Checks.TopLevel; recorded != "object" {
		t.Fatalf("the fixture records top level %q, so this test no longer covers the "+
			"case it was written for", recorded)
	}
	// The server sent an array, and an element of it does not decode, so the
	// rejection is real and its stated cause is the page's top level.
	sk := liveSkeleton(t, `[{"close":"not-a-number"}]`)
	if got, want := bodyKind(sk), "array"; got != want {
		t.Fatalf("the body is a %s, so this test no longer covers the case it was "+
			"written for: the detail would be right", got)
	}

	rows := CompareLive(f, "test", reflect.TypeOf([]row{}), sk)

	var failures []LiveDivergence
	for _, r := range rows {
		if r.Kind == DecodeFailure && r.Direction == LiveSDKDirection {
			failures = append(failures, r)
		}
	}
	if len(failures) != 1 {
		t.Fatalf("got %d live decode finding(s), want 1; rows: %v", len(failures), rows)
	}
	detail := failures[0].Detail
	if !strings.Contains(detail, "does not unmarshal") {
		t.Fatalf("the detail does not report a rejection: %q", detail)
	}
	if !assertDetailKindsAreTrue(t, f, sk, detail) {
		t.Error("the detail states a kind the evidence contradicts: " + detail)
	}
	// Spelled out as well, so the failure message is readable without the
	// checker's vocabulary: the page's object must appear as the page's and never
	// as the body's.
	if strings.Contains(detail, neutralDetailPhrase+" top level is object") {
		t.Errorf("the detail says the live body is an object; it is an array, and the "+
			"object is the page's: %q", detail)
	}
	if !strings.Contains(detail, recordedTopLevel+" is object") {
		t.Errorf("the detail does not say the object is the page's recorded top level: %q", detail)
	}
}

// TestLiveBaselineIsCanonical is the reviewability claim: the committed file is
// what the harness would write, so a regeneration is a diff of evidence rather
// than of formatting.
//
// The way to break this test is a renderer that drops, reorders or rephrases a
// field, which would make every regeneration a diff nobody can read.
func TestLiveBaselineIsCanonical(t *testing.T) {
	raw, err := liveBaselineFS.ReadFile(LiveBaselineName)
	if err != nil {
		t.Fatalf("read %s: %v", LiveBaselineName, err)
	}
	baseline, err := ParseLiveBaseline(raw)
	if err != nil {
		t.Fatalf("parse %s: %v", LiveBaselineName, err)
	}
	rendered, err := baseline.Marshal()
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(rendered) != string(raw) {
		t.Errorf("%s is not canonical: re-rendering it differs.\ncommitted:\n%s\nrendered:\n%s",
			LiveBaselineName, raw, rendered)
	}
}

// TestLiveBaselineNotesDescribeTheRun holds the whole notes block to the run it
// describes, sentence for sentence.
//
// The notes are the part of the file a reader quotes, and a note is prose, so
// nothing else here would notice one that had stopped being true: the gate
// compares findings and coverage, the bucket test compares counts it computes
// itself, and the notes are only checked for having the right keys. A live-only
// note still reading "It holds 11 finding(s)" after the bucket emptied, or still
// claiming a decomposition the tree no longer has, would pass every test.
//
// The reason this matters rather than being tidiness is the live-only note
// specifically. It is the headline a reader takes away, and the number in it is
// the one number in the file that is not a count of defects. Rendering it from
// countLiveOnly rather than transcribing it means the qualification travels with
// the count, and comparing the whole map means the qualification cannot survive
// the evidence it qualifies.
//
// The way to break this test is to edit the committed notes without editing
// liveBaselineNotes, which is the whole point: one of the two is then wrong and
// this says which.
func TestLiveBaselineNotesDescribeTheRun(t *testing.T) {
	baseline, err := LoadLiveBaseline()
	if err != nil {
		t.Fatalf("load %s: %v", LiveBaselineName, err)
	}
	run := CompareAllLive()
	if len(run.Errs) > 0 {
		t.Fatalf("the run could not compare everything: %v", run.Errs)
	}
	want := liveBaselineNotes(run)
	for _, key := range sortedKeys(want) {
		if got, ok := baseline.Notes[key]; !ok {
			t.Errorf("%s has no note for %q", LiveBaselineName, key)
		} else if got != want[key] {
			t.Errorf("%s: the note for %q does not describe this run.\ncommitted: %s\n"+
				"derived:   %s", LiveBaselineName, key, got, want[key])
		}
	}
	for key := range baseline.Notes {
		if _, ok := want[key]; !ok {
			t.Errorf("%s records a note for %q, which nothing derives any more: a note "+
				"no run produces is a sentence about a run that did not happen",
				LiveBaselineName, key)
		}
	}
}

// TestLiveBucketCountsAreTotal is the conservation check.
//
// Every finding the comparison produced lands in exactly one bucket, and the
// buckets sum to the number of classifications. A finding that fell out of the
// classification -- dropped by a key collision, or overwritten in a map -- would
// make the sum short, and the whole report is a count.
//
// The way to break this test is a bucket assignment that does not cover every
// key, which is precisely the silent loss this arithmetic catches.
func TestLiveBucketCountsAreTotal(t *testing.T) {
	run := CompareAllLive()
	counts := BucketCounts(run.Classifications)
	total := 0
	for _, bucket := range []LiveBucket{LiveBoth, LiveBothDetailChanged, LiveOnly, LiveDocsOnly} {
		total += counts[bucket]
	}
	if total != len(run.Classifications) {
		t.Errorf("the buckets hold %d finding(s) and the comparison classified %d: %d "+
			"finding(s) fell out of the classification", total, len(run.Classifications),
			len(run.Classifications)-total)
	}
	// Every bucket is present at zero rather than absent, so a reader can tell an
	// empty bucket from one nobody counted.
	if len(counts) != 4 {
		t.Errorf("BucketCounts reports %d bucket(s), want 4: %v", len(counts), counts)
	}
}

// TestLiveComparisonCollidesOnNoPairKey is the arithmetic's blind spot, named.
//
// TestLiveBucketCountsAreTotal sums buckets over classifications, so a key that
// collapsed on the way into ClassifyLive has already become one classification by
// the time it counts: the sum stays right and one detail is gone. This looks at
// the rows instead, and it is what makes "none is dropped" a statement about the
// real run rather than a claim about a function that cannot be reached this way.
//
// Each CompareBody emits at most one finding per (kind, name), so a collision
// cannot arise from CompareAllLive. That is the property being checked rather than
// the reason for trusting the check.
//
// The way to break this test is any change that lets two rows of one direction
// share a pair key -- a check that reports a name twice, a run that visits a
// fixture under two symbols.
func TestLiveComparisonCollidesOnNoPairKey(t *testing.T) {
	run := CompareAllLive()
	seen := map[string]LiveDivergence{}
	for _, r := range run.Rows {
		key := string(r.Direction) + "\x1f" + r.PairKey()
		if prev, dup := seen[key]; dup {
			t.Errorf("two %s findings share the pair key %q, so one of them is a detail "+
				"the classification can only keep one of; the bucket counts are computed "+
				"from classifications and would not notice.\nfirst:  %s\nsecond: %s",
				r.Direction, r.PairKey(), prev.Detail, r.Detail)
		}
		seen[key] = r
	}
	for _, c := range run.Classifications {
		if len(c.Colliding) > 0 {
			t.Errorf("%s / %s / %q is classified with %d further detail(s) it did not "+
				"keep: %v. CompareBody emits at most one finding per (kind, name), so "+
				"this cannot come from the checks and means a row was produced twice",
				c.Symbol, c.Kind, c.Name, len(c.Colliding), c.Colliding)
		}
	}
}

// TestClassifyLiveReportsACollidingPairKey is the exported contract, driven.
//
// ClassifyLive documents that no finding is dropped, and the version that
// documented it that way could drop one: two rows on one side sharing a pair key
// became one classification with the last row's detail and the first gone, with
// no error and no word. This is the input that does it, and the assertion is
// that every detail the classifier was handed comes back out -- carried in Detail
// or listed in Colliding.
//
// The assertion is deliberately about the *set* of details rather than about the
// field, so it states the claim and not the implementation: a classifier that
// kept the first detail and reported the second as a collision passes, one that
// kept the first and dropped the second does not, and so does one that keeps the
// second. Which detail is carried is a presentation choice; losing one is not.
//
// It is unreachable through CompareAllLive and is here anyway, because the
// function is exported and a caller will read its contract rather than discover
// which inputs its producer happens to allow.
//
// The way to break this test is the map assignment the last row of the pair key
// sees, which is what the previous version did.
func TestClassifyLiveReportsACollidingPairKey(t *testing.T) {
	const fixture = "test/GET-example"
	rows := []LiveDivergence{
		{Symbol: "s", Fixture: fixture, Direction: LiveSDKDirection, Kind: DecodeFailure,
			Detail: "the compared body's instance does not unmarshal: the recorded top level is object"},
		{Symbol: "s", Fixture: fixture, Direction: LiveSDKDirection, Kind: DecodeFailure,
			Detail: "the compared body's instance does not unmarshal: no name, shape or leaf check disagrees"},
		// An exact repeat is one observation, not two, and is held separately below.
		{Symbol: "s", Fixture: fixture, Direction: LiveSDKDirection, Kind: LeafTypeMismatch,
			Name: "price", Detail: "the compared body's number, SDK field P.Price is string"},
		{Symbol: "s", Fixture: fixture, Direction: LiveSDKDirection, Kind: LeafTypeMismatch,
			Name: "price", Detail: "the compared body's number, SDK field P.Price is string"},
	}

	got := ClassifyLive(rows)
	if len(got) != 2 {
		t.Fatalf("got %d classification(s), want 2: a pair key observed three times is "+
			"still one finding: %+v", len(got), got)
	}
	byKey := map[string]LiveClassification{}
	for _, c := range got {
		byKey[c.Key] = c
	}

	decode, ok := byKey[pairKeyOf("s", fixture, DecodeFailure, "")]
	if !ok {
		t.Fatalf("the colliding key was not classified at all: %+v", got)
	}
	// Every detail handed in is accounted for, which is the whole claim.
	carried := map[string]bool{decode.Detail: true}
	for _, extra := range decode.Colliding {
		carried[extra] = true
	}
	for _, r := range rows {
		if r.Kind != DecodeFailure {
			continue
		}
		if !carried[r.Detail] {
			t.Errorf("a detail handed to the classifier is in neither Detail (%q) nor "+
				"Colliding (%v): the finding was dropped without a word", r.Detail, decode.Colliding)
		}
	}
	if len(decode.Colliding) != 1 {
		t.Errorf("the colliding key reports %d further detail(s), want 1: %v",
			len(decode.Colliding), decode.Colliding)
	}
	// The choice of which to carry is order-free, because the whole function is.
	forward := ClassifyLive(rows)
	reversed := ClassifyLive([]LiveDivergence{rows[3], rows[2], rows[1], rows[0]})
	if !reflect.DeepEqual(forward, reversed) {
		t.Errorf("which detail a collision carries depends on row order.\nforward %+v\n"+
			"reversed %+v", forward, reversed)
	}

	// An exact repeat is not a collision: it is the same observation twice, and
	// reporting it as two would turn one finding into an apparent disagreement.
	leaf, ok := byKey[pairKeyOf("s", fixture, LeafTypeMismatch, "price")]
	if !ok {
		t.Fatalf("the repeated key was not classified: %+v", got)
	}
	if len(leaf.Colliding) != 0 {
		t.Errorf("one detail supplied twice is reported as a collision: %v", leaf.Colliding)
	}
}

// TestLiveComparisonIsDeterministic runs the whole comparison twice and requires
// the same answer, because the gate compares sets and a set that moves between
// runs would make the gate flaky rather than strict.
//
// The way to break this test is anything ordered -- a map range whose result
// reaches the output, or a comparison whose result depends on map iteration.
func TestLiveComparisonIsDeterministic(t *testing.T) {
	first := CompareAllLive()
	second := CompareAllLive()
	if !reflect.DeepEqual(first.Classifications, second.Classifications) {
		t.Errorf("two runs of the same committed tree classified differently.\nfirst:  %+v\nsecond: %+v",
			first.Classifications, second.Classifications)
	}
	if !reflect.DeepEqual(first.Coverage, second.Coverage) {
		t.Errorf("two runs reported different coverage.\nfirst:  %+v\nsecond: %+v",
			first.Coverage, second.Coverage)
	}
	if len(first.Rows) != len(second.Rows) {
		t.Errorf("two runs produced %d and %d raw finding(s)", len(first.Rows), len(second.Rows))
	}
}

// TestLiveComparisonCoversEveryManifestRow is the tree-index claim.
//
// Every row the live manifest records must be a row the comparison reached, and
// every row it reached must be one the manifest records. A skeleton the
// comparison skipped is a finding that vanishes, which is the failure mode this
// package exists to end, and a manifest row with no comparison is a claim about
// the sandbox that nothing checked.
//
// The way to break this test is a filter over the rows, or a join that drops one
// side.
func TestLiveComparisonCoversEveryManifestRow(t *testing.T) {
	manifest := readLiveManifest(t)
	run := CompareAllLive()
	if len(run.Errs) != 0 {
		t.Fatalf("the comparison reported %d error(s): %v", len(run.Errs), run.Errs)
	}
	if run.Coverage.Probed != len(manifest.Entries) {
		t.Errorf("the comparison probed %d endpoint(s), the manifest records %d",
			run.Coverage.Probed, len(manifest.Entries))
	}
	if run.Coverage.Compared != len(manifest.Entries) {
		t.Errorf("the comparison resolved a type for %d of %d endpoint(s); a row without "+
			"one is a hole in coverage and not a pass",
			run.Coverage.Compared, len(manifest.Entries))
	}
}

// TestLiveRunCarriesNoUnclassifiedRows is the plumbing check: every row the
// comparison produced was classified, so nothing was produced and dropped.
func TestLiveRunCarriesNoUnclassifiedRows(t *testing.T) {
	run := CompareAllLive()
	classified := map[string]bool{}
	for _, c := range run.Classifications {
		classified[c.Key] = true
	}
	var unclassified []string
	for _, r := range run.Rows {
		if !classified[r.PairKey()] {
			unclassified = append(unclassified, r.PairKey())
		}
	}
	sort.Strings(unclassified)
	if len(unclassified) > 0 {
		t.Errorf("%d raw finding(s) were classified under no key: %v", len(unclassified), unclassified)
	}
}
