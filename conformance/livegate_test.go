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
// this harness can produce, and the review that established they are all
// top-level-kind mismatches is only worth anything if the four stay visible. A
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

// TestLiveBaselineDetailsNameWhatTheyRead pins the neutral phrasing.
//
// A live finding that says "the documented instance" is attributing a server
// observation to a Webull page, which would make the whole set misread. The
// rewrite is mechanical, so a check that grows a new phrase would leave a stale
// "documented" behind, and that has to be visible rather than merely wrong in a
// way nobody reads.
//
// The second half is the part that is easy to get wrong in the other direction. A
// detail must name *what the check read*, and the checks do not all read the same
// thing: four read the body and the shape check reads the manifest. So a shape
// finding's detail says "the recorded top level" and a body finding's says "the
// compared body's". Demanding the body phrasing of a shape finding would be
// demanding a false statement -- the shape check never looked at the body, and
// data.GetCapitalFlow's live body is an array while its manifest row records an
// object.
//
// The way to break this test is to drop the substitution list, which leaves
// CompareBody's fixture-specific wording in the file.
func TestLiveBaselineDetailsNameWhatTheyRead(t *testing.T) {
	baseline, err := LoadLiveBaseline()
	if err != nil {
		t.Fatalf("load %s: %v", LiveBaselineName, err)
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
		switch e.Kind {
		case TopLevelMismatch, ElementTypeMismatch:
			// These two read the manifest, so their detail must say so.
			if !strings.Contains(e.Detail, "recorded") {
				t.Errorf("%s: %s is a shape finding whose detail does not say it read "+
					"the manifest's recorded check: %q; the shape check never reads a "+
					"body, so a detail implying otherwise asserts something unexamined",
					LiveBaselineName, where, e.Detail)
			}
		default:
			// Every other check reads the body it was handed.
			if !strings.Contains(e.Detail, neutralDetailPhrase) {
				t.Errorf("%s: %s has a detail that never says it read a body: %q",
					LiveBaselineName, where, e.Detail)
			}
		}
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
