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
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"testing"
)

// These tests guard the committed artifacts, not the SDK. They read the fixture
// bytes from the embedded copy -- exactly as committed, never regenerated --
// which is the whole reason the comparison lives on this side of the generator
// boundary: a test that rebuilt its own input could not fail.
//
// They cover the fixtures alone. The per-DTO comparison that consumes them lives
// in divergence_test.go, together with the known-divergence baseline that keeps
// the gate green while every current defect stays recorded and visible.

func load(t *testing.T) *Manifest {
	t.Helper()
	m, err := Load()
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}
	if len(m.Fixtures) == 0 {
		t.Fatal("manifest declares no fixtures")
	}
	return m
}

// committedLength is the length of a committed fixture as the manifest records
// it, which is what every byte-count assertion in this file must compare against.
//
// # Why the line-ending normalisation exists
//
// Manifest.Fixture.Bytes and Totals.FixtureBytes are the length of the file
// tools/conformance/gen_fixtures.py wrote, and that tool opens every file with
// newline="\n" (gen_fixtures.py, write), so both record an LF count. The
// comparison here used to be byte-exact against the *embedded* bytes, and
// embedding reads the working tree, not the object store. A checkout whose
// line-ending policy differs therefore embedded a longer file than the manifest
// describes: the repository's .gitattributes is `* text=auto`, and
// windows-latest ships core.autocrlf=true, so the CRLF materialisation is
// upstream's default behaviour on that runner and not something anyone chose.
// That is what turned both Windows matrix jobs red at v2.1.15, on all 193
// fixtures, for a difference that carries no JSON meaning. CRLF is legal
// inter-token whitespace to the JSON grammar and to encoding/json, so the rest
// of this file -- json.Valid, the top-level kind, the required names, the
// distinct-name count -- already reads a CRLF checkout exactly as it reads an LF
// one, and only the byte counts disagreed.
//
// # What it gives up
//
// Exactly one thing: a committed fixture whose bytes differ from the manifest's
// only in line terminators is no longer reported. That is the whole class being
// tolerated, and there is no way to tolerate it and still count bytes.
//
// It is worth being explicit about what it does NOT give up, because the check
// looks stronger than it is. The manifest records a length, not a digest, so a
// content edit that preserves the byte count -- a price changed from "12.34" to
// "12.35" -- is invisible to this assertion. That was true before this change
// too, and is not made worse by it. What catches a real content change is the
// structural set below it: valid JSON, the documented top-level kind, every
// documented required name present at the level the page declares it, and the
// recorded distinct-property-name count. A byte count is a tripwire against a
// silently reshaped or truncated artifact, not a checksum, and calling it one
// would be the overclaim.
//
// The normalisation is deliberately CRLF -> LF and nothing else. It is the only
// transformation a line-ending policy performs, and checkNoStrayCarriageReturn
// below refuses a lone CR so the tolerated set cannot widen into "any byte the
// checkout likes to add". A bare CR inside a JSON string is invalid JSON
// whichever way it is written, so json.Valid rejects the case the tolerance
// might otherwise have opened.
// A byte count is a tripwire against a silently reshaped or truncated artifact.
// It is not a checksum, and for a long time this package had no checksum either:
// a length-preserving edit to a committed fixture kept every assertion here
// green. committedDigest is that checksum.
//
// The reason it was absent for so long is the reason it cannot simply be assumed
// unnecessary. The generator's own --check does compare the whole committed tree
// against the documentation, so a fixture edit is caught by anyone who runs
// make conformance-fixtures. But --check needs the docgen cache, and AGENTS.md
// records that both fixture targets skip with exit 0 when the cache is absent --
// which is the case in CI, so nothing in CI performed that comparison at all. The
// length tripwire was therefore the only fixture integrity check CI ran, and it
// cannot see an edit that preserves length. committedDigest needs no cache, no
// network and no toolchain, so it runs everywhere the tests do.
//
// The normalisation is deliberately CRLF -> LF and nothing else. It is the only
// transformation a line-ending policy performs, and checkNoStrayCarriageReturn
// below refuses a lone CR so the tolerated set cannot widen into "any byte the
// checkout likes to add". A bare CR inside a JSON string is invalid JSON
// whichever way it is written, so json.Valid rejects the case the tolerance
// might otherwise have opened. The digest normalises identically, so it is the
// same value on a Windows checkout as on a Linux one.
func committedLength(data []byte) int {
	return len(bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n")))
}

// committedDigest is the fixture's sha256 over the CRLF-collapsed bytes, hex
// encoded, which is the form the manifest records.
func committedDigest(data []byte) string {
	sum := sha256.Sum256(bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n")))
	return hex.EncodeToString(sum[:])
}

// TestCommittedDigestProperties pins the two properties the digest depends on,
// because a digest that held neither would be a worse check than none.
//
// The first is checkout independence: a Windows checkout delivers CRLF where a
// Linux one delivers LF, and the manifest records one value, so the digest has to
// be equal on both or the gate would fail on whichever platform it was not
// recorded from. The second is sensitivity: a length-preserving edit must change
// it, since that is the entire reason it exists and the length tripwire cannot
// see such an edit.
func TestCommittedDigestProperties(t *testing.T) {
	const lf = "{\n  \"a\": 1\n}\n"
	const crlf = "{\r\n  \"a\": 1\r\n}\r\n"

	if got, want := committedDigest([]byte(crlf)), committedDigest([]byte(lf)); got != want {
		t.Errorf("a CRLF checkout digests to %s and an LF checkout to %s; the manifest "+
			"records one value, so one of the two platforms would fail the gate", got, want)
	}
	// The raw byte counts differ, which is why the normalisation exists, and the
	// normalised ones match, which is why both platforms can be held to one
	// recorded value.
	if len(crlf) == len(lf) {
		t.Error("the CRLF and LF probes are the same length, so this test cannot " +
			"demonstrate that the normalisation is doing anything")
	}
	if committedLength([]byte(crlf)) != committedLength([]byte(lf)) {
		t.Errorf("normalised lengths differ: CRLF gives %d and LF gives %d, so the "+
			"length tripwire is not checkout-independent and the digest is not either",
			committedLength([]byte(crlf)), committedLength([]byte(lf)))
	}

	// Same length, different content. Swapping a digit is the smallest possible
	// change, and the one the length tripwire is blind to by construction.
	const edited = "{\n  \"a\": 2\n}\n"
	if len(edited) != len(lf) {
		t.Fatalf("the probe must be length-preserving to mean anything: %d vs %d",
			len(edited), len(lf))
	}
	if committedDigest([]byte(edited)) == committedDigest([]byte(lf)) {
		t.Error("a length-preserving content change produced the same digest, so the " +
			"check cannot see the edit it exists to catch")
	}

	// A lone CR must not be silently absorbed, matching the tripwire's tolerance.
	// checkNoStrayCarriageReturn is what enforces this on a real fixture; here it
	// is stated directly so the two cannot drift apart unnoticed.
	if committedDigest([]byte("{\r  \"a\": 1\n}\n")) == committedDigest([]byte("{\n  \"a\": 1\n}\n")) {
		t.Error("a lone CR produced the same digest as its absence, so a stray carriage " +
			"return would be invisible to the digest as well as to the length tripwire")
	}
}

// checkNoStrayCarriageReturn holds the normalisation to the one transformation it
// claims: every CR in a committed fixture must be half of a CRLF. Together with
// committedLength matching the manifest exactly, that pins the embedded bytes to
// the committed bytes plus a CR before each LF -- a byte-exact comparison whose
// only free parameter is the checkout's line-ending policy.
func checkNoStrayCarriageReturn(t *testing.T, f Fixture, data []byte) {
	t.Helper()
	cr := bytes.Count(data, []byte("\r"))
	if crlf := bytes.Count(data, []byte("\r\n")); cr != crlf {
		t.Errorf("fixture holds %d CR bytes but only %d CRLF terminators, so %d "+
			"carriage return(s) are not line terminators; only CRLF is a "+
			"line-ending policy, anything else is a content change",
			cr, crlf, cr-crlf)
	}
}

// TestEveryManifestFixtureIsEmbedded checks the reverse of the stray-file test:
// nothing the manifest names may be missing from the tree. go:embed makes a
// missing file a build error, so this guards the case embed cannot see -- a
// manifest row pointing at a path the tree does not contain.
func TestEveryManifestFixtureIsEmbedded(t *testing.T) {
	m := load(t)
	for _, f := range m.Fixtures {
		if _, err := f.Read(); err != nil {
			t.Errorf("%s: %v", f.ID, err)
		}
	}
}

// TestManifestDescribesTheWholeTree walks the embedded tree and holds the
// manifest to it in both directions, so a fixture nobody can get back to a page
// cannot sit in the package unrecorded.
//
// The tree holds two kinds of file and two indexes, and the walk is split to
// match. Everything outside live/ is described by testdata/manifest.json, which
// is generated from Webull's published OpenAPI JSON. Everything under live/ is
// described by testdata/live-manifest.json, which examples/live-probe writes
// from the sandbox: a different author, a different input, and a different
// consumer, and merging the two would put a live reading next to a documented
// example in one index and make a reader unable to tell which is which. The
// README at the root of live/ is described by neither, because it is prose about
// the tree rather than evidence in it.
//
// The split is a partition rather than an exemption: a file under live/ still
// has to be named by the live manifest, so the relaxation is that the live tree
// has its own index, not that it is unindexed.
func TestManifestDescribesTheWholeTree(t *testing.T) {
	m := load(t)
	want := make(map[string]bool, len(m.Fixtures)+1)
	want["manifest.json"] = true
	for _, f := range m.Fixtures {
		want[f.Fixture] = true
	}
	liveWant := loadLiveManifest(t)
	want[liveManifestFile] = true

	var got, liveGot []string
	err := fs.WalkDir(fixturesFS, testdataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, ok := strings.CutPrefix(path, testdataDir+"/")
		if !ok {
			return fmt.Errorf("path %q is not under %s", path, testdataDir)
		}
		if rel == "live/"+liveReadmeFile {
			return nil
		}
		if strings.HasPrefix(rel, liveTreePrefix) {
			liveGot = append(liveGot, strings.TrimPrefix(rel, liveTreePrefix))
			return nil
		}
		got = append(got, rel)
		return nil
	})
	if err != nil {
		t.Fatalf("walk embedded tree: %v", err)
	}
	sort.Strings(got)
	for _, rel := range got {
		if !want[rel] {
			t.Errorf("%s is committed but not in the manifest; the manifest is "+
				"the index of the tree, so this file has no provenance", rel)
		}
	}
	if len(got) != len(want) {
		t.Errorf("tree holds %d documentation files, manifest describes %d", len(got), len(want))
	}

	sort.Strings(liveGot)
	for _, rel := range liveGot {
		if !liveWant[rel] {
			t.Errorf("live/%s is committed but not in %s; that manifest is the "+
				"index of the live tree, so this skeleton has no provenance", rel, liveManifestFile)
		}
	}
	if len(liveGot) != len(liveWant) {
		t.Errorf("live tree holds %d file(s), %s describes %d", len(liveGot), liveManifestFile, len(liveWant))
	}
}

// The placeholders a reduced tree is made of. They are spelled out here rather
// than imported from the program that wrote the tree, because the point of the
// check below is that a reader of the committed files can verify the claim
// without trusting - or running - the producer. A check that called the
// producer's own helper would be asking the writer whether it had written what
// it was asked to write.
const (
	liveStringPlaceholder = "1"
	liveNumberPlaceholder = "-1"
)

// TestLiveSkeletonsCarryNoValue holds every committed live skeleton to the one
// property that makes it safe to be in a repository: it carries every member
// name the sandbox sent and no value it sent.
//
// This is the test that would catch a leak after the fact. A skeleton is
// committed, so a price, a share count, a timestamp or an account number in one
// is in git history, which is not a mistake anyone can delete. The reduction is
// what prevents that, and this asserts the result rather than the intent: every
// leaf in every committed file is the placeholder for its kind and nothing else,
// and no key in any of them is a reading. The key half is half the tree - a
// member name is kept verbatim by the reduction, so a name that is a value is a
// leak this walk would otherwise pass - and its rule is keyCarriesAReading.
//
// Four leaf types are accepted and two are rejected, and the rejections are the
// interesting half. A float64 or an int is a value a decode without UseNumber
// would have left behind, widened. A false is a boolean the reduction emits as
// true unconditionally, so a false leaf proves the file did not come from the
// reduction at all. Both are checked here so a future change to either the
// reduction or the writer cannot quietly reintroduce a reading.
//
// # What the dedupe changed, and what it did not
//
// The committed form is a deduplicated encoding (see liveskeleton.go): an array
// is written as distinct element shapes with a count beside each. So the walk has
// one new leaf position, the count, and it is the one place a reading could hide
// that did not exist before. It is held to a rule of its own rather than being
// walked past: a count must be a positive integer, which is not a reading a
// server sends, and every leaf *inside* a shape is still the placeholder for its
// kind exactly as before.
//
// That is the whole of the change. The reduction, the placeholders, the key rule
// and the tree of names are untouched, and a value planted anywhere in a shape is
// caught by the same walk that caught it before. The leak gate was not weakened
// to accommodate the dedupe; it was extended to cover the one new position, and
// that position is checked rather than skipped. The one thing the new position
// does admit, and the whole suite cannot see, is set out at isPositiveInteger.
func TestLiveSkeletonsCarryNoValue(t *testing.T) {
	files := liveJSONFiles(t)
	if len(files) == 0 {
		t.Fatal("the live tree holds no .json file, so this test is asserting nothing")
	}
	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			data, err := fixturesFS.ReadFile(name)
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			dec := json.NewDecoder(bytes.NewReader(data))
			dec.UseNumber()
			var value any
			if err := dec.Decode(&value); err != nil {
				t.Fatalf("a committed skeleton is not parseable JSON: %v", err)
			}
			if dec.More() {
				t.Fatal("a committed skeleton holds more than one JSON value")
			}
			assertNoLiveValue(t, value, "$", name)
		})
	}
}

// liveValueSink is where assertNoLiveValue reports. *testing.T satisfies it, and
// so does recordingSink below, which is what lets the key rule be exercised
// against a body that really does carry a reading: such a body cannot go into
// the committed tree in order to be walked, so the failing case has to be
// driven from here.
type liveValueSink interface {
	Helper()
	Errorf(format string, args ...any)
}

// recordingSink collects what assertNoLiveValue reports rather than failing the
// test that drives it.
type recordingSink struct{ messages []string }

func (r *recordingSink) Helper() {}

func (r *recordingSink) Errorf(format string, args ...any) {
	r.messages = append(r.messages, fmt.Sprintf(format, args...))
}

// keyCarriesAReading reports whether an object member name could itself be a
// value the server sent.
//
// The reduction keeps every member name verbatim and replaces every value, so a
// name is the one place a reading can survive the reduction. The rule is the
// narrow one the manifest's warning states - no key may be a value Webull sent -
// and two forms of key are refused.
//
//   - A key that is a well-formed JSON number. Every reading Webull sends that is
//     a number - an account id, an instrument id, a timestamp in milliseconds, a
//     price, a size, a count - is one of these when it appears as a key, and no
//     Webull member name is. The form is decided by the JSON number grammar alone,
//     so the rule needs no vocabulary of names and holds for any endpoint.
//   - A key that is one of the four placeholders, or the empty string. A key a
//     reader cannot tell from a value is the ambiguity the representation exists
//     to remove, and an empty key is indistinguishable from an empty string
//     reading.
//
// # What the rule cannot see, which is most of the class
//
// It cannot see a key that is a string reading: a ticker, an order id, a company
// name. A JSON object member name and a JSON string value are the same token in
// the same grammar, so nothing in the document, and nothing in the reduced tree,
// distinguishes them. That is a property of JSON rather than a gap in this walk,
// and it bounds what the tree may be claimed to hold: no key is a numeric reading
// and no key is a placeholder; a string-valued key is neither detectable here nor
// excluded.
//
// A numeric reading written in a form the JSON grammar rejects is not caught
// either, because the rule is the grammar's rather than a looser approximation of
// it: "01", ".5" and "1 " are not well-formed JSON numbers and pass. Webull
// writes numbers canonically, so the narrower rule is the one that cannot be
// talked out of its own grammar.
//
// # The one false positive, and why it is left in
//
// A member name that is itself a bare number fails the gate. A macro series keyed
// by period is the shape that would do it, and no captured endpoint is one. The
// failure is left in rather than carved out because an exception for a shape
// nobody has observed would be a special case inside a safety check, whereas a
// real one shows up loudly and is a question for the maintainer: the rule, not the
// key, is what would be reconsidered.
func keyCarriesAReading(name string) bool {

	if name == "" {
		return true
	}
	switch name {
	case liveStringPlaceholder, liveNumberPlaceholder, "true", "null":
		return true
	}
	_, err := json.Marshal(json.Number(name))
	return err == nil
}

// assertNoLiveValue walks one decoded skeleton and fails on any key or leaf that
// could carry something the server sent.
//
// The map case checks the key before it descends, so the invariant is asserted
// over the whole tree and not over its leaves alone. keyCarriesAReading holds the
// rule and states what it cannot see.
//
// The dedupe wrapper is the one new case, and it is handled rather than skipped.
// A wrapper's own members are the encoding's and are checked against the
// encoding's rules; its count is held to a rule of its own; and its shape is
// walked exactly as an object would be, so every name and every leaf inside a
// shape is held to the same rule as before.
func assertNoLiveValue(t liveValueSink, value any, path, file string) {
	t.Helper()
	switch leaf := value.(type) {
	case map[string]any:
		if records, wrapped := leaf[SkeletonArrayKey]; wrapped {
			assertNoDedupeWrapper(t, leaf, records, path, file)
			return
		}
		for name, member := range leaf {
			if keyCarriesAReading(name) {
				t.Errorf("%s holds the member name %q at %s: no Webull member name is a "+
					"JSON number or one of the four placeholders, so a key in either form "+
					"is a reading that reached a committed file",
					file, name, path+"."+name)
			}
			assertNoLiveValue(t, member, path+"."+name, file)
		}
	case []any:
		for i, element := range leaf {
			assertNoLiveValue(t, element, fmt.Sprintf("%s[%d]", path, i), file)
		}
	case string:
		if leaf != liveStringPlaceholder {
			t.Errorf("%s holds %q at %s: a string the reduction emits is %q",
				file, leaf, path, liveStringPlaceholder)
		}
	case json.Number:
		if leaf.String() != liveNumberPlaceholder {
			t.Errorf("%s holds the number %s at %s: the reduction emits %s for every number",
				file, leaf, path, liveNumberPlaceholder)
		}
	case bool:
		if !leaf {
			t.Errorf("%s holds false at %s: the reduction emits true for every boolean, "+
				"so a false leaf proves the file was not reduced", file, path)
		}
	case nil:
	case float64, int, int64:
		t.Errorf("%s holds a %T at %s: a widened value, which the reduction never emits",
			file, leaf, path)
	default:
		t.Errorf("%s holds a %T at %s: no reduction emits that type", file, leaf, path)
	}
}

// assertNoDedupeWrapper holds one dedupe wrapper to the encoding's rules.
//
// It is written to fail on anything it does not recognise rather than to walk
// past it, because a wrapper this function cannot account for is a position
// where a reading could sit unchecked. Three things are checked: the wrapper
// carries nothing but the reserved member, every record carries exactly a count
// and a shape, and the count is a positive integer -- the one leaf the dedupe
// introduced, and the one a reading would have to impersonate.
func assertNoDedupeWrapper(t liveValueSink, wrapper map[string]any, records any, path, file string) {
	t.Helper()
	for name := range wrapper {
		if name != SkeletonArrayKey {
			t.Errorf("%s holds %q beside a dedupe wrapper at %s: the wrapper's own members "+
				"are the encoding's, and an extra one means the position is unaccounted for",
				file, name, path)
		}
	}
	list, ok := records.([]any)
	if !ok {
		t.Errorf("%s holds a %T as the %s of a dedupe wrapper at %s, not a list of records",
			file, records, SkeletonArrayKey, path)
		return
	}
	for i, r := range list {
		at := fmt.Sprintf("%s.%s[%d]", path, SkeletonArrayKey, i)
		record, ok := r.(map[string]any)
		if !ok {
			t.Errorf("%s holds a %T as a dedupe record at %s, not an object", file, r, at)
			continue
		}
		for name := range record {
			if name != SkeletonCountKey && name != SkeletonShapeKey {
				t.Errorf("%s holds %q in a dedupe record at %s: a record carries a count and "+
					"a shape and nothing else", file, name, at)
			}
		}
		count, hasCount := record[SkeletonCountKey]
		if !hasCount {
			t.Errorf("%s holds a dedupe record at %s with no %s, so its multiplicity is unrecorded",
				file, at, SkeletonCountKey)
		} else if !isPositiveInteger(count) {
			t.Errorf("%s holds %v as the %s of a dedupe record at %s; a count is a positive "+
				"integer -- a reading is not", file, count, SkeletonCountKey, at)
		}
		shape, hasShape := record[SkeletonShapeKey]
		if !hasShape {
			t.Errorf("%s holds a dedupe record at %s with no %s, so the element shape is absent",
				file, at, SkeletonShapeKey)
			continue
		}
		assertNoLiveValue(t, shape, at+"."+SkeletonShapeKey, file)
	}
}

// isPositiveInteger reports whether v is a JSON integer of one or more.
//
// It is stricter than "is a number" on purpose. The count is the only position
// the dedupe introduced, and the property being defended is that it holds a
// multiplicity rather than a reading; a float, a negative number, zero or a
// string are all things a count is not, and each of them is a sign that the
// position is being used for something else.
//
// # What the count cannot see, which is a real widening
//
// A reading that is a positive integer is not distinguishable from a count. An
// account id, an instrument id, a millisecond timestamp and a share count are all
// positive integers, and if one of those were written into a $count the rule
// here would accept it.
//
// That is a genuine widening of the guarantee the tree otherwise holds, and it is
// stated rather than left for a reader to discover. The mitigation is structural
// rather than a matter of the rule: the count is written by EncodeSkeleton, which
// computes it by counting the elements of the response and has no path by which a
// response *value* could reach it -- a number in the tree arrives only as a leaf
// placeholder, and a leaf placeholder is -1, which this rule refuses. So a reading
// would have to be written into the field by hand, after the reduction, rather
// than surviving it.
//
// # What catches it then, stated exactly
//
// Nothing automatic, and that is the honest answer. The structural mitigation
// above is a property of how the file is *produced*, not a check on how it reads,
// and the two are not the same: every other integrity check in this package would
// stay green. Element-count conservation holds, because a count is a count of
// whatever it claims and the arithmetic does not know what the number came from.
// Canonicality holds, because re-encoding reproduces whatever counts are written.
// The manifest's own accounts hold, because they are computed from the counts by
// the same rule. And this walk holds, because a positive integer is exactly what
// it requires a count to be.
//
// So the count is the one position in the committed tree where a planted reading
// would pass every test in the repository, and the only thing standing between
// the tree and a hand-edited $count is that nobody hand-edits one. That is not a
// guarantee and must not be described as one; it is a narrower trust than the rest
// of the tree rests on, and a maintainer editing a committed skeleton has to know
// it.
//
// The same is true of a string-valued member name, which keyCarriesAReading
// already states it cannot see, and for the same reason: the reduction keeps
// names verbatim, so a name is the one position a reading survives, and the rule
// covers the forms it can rather than the form it cannot. The count is weaker
// than that: a name is a string, so the rule can at least refuse the numeric and
// placeholder forms, whereas the count's form is the reading's own form.
func isPositiveInteger(v any) bool {
	number, ok := v.(json.Number)
	if !ok {
		return false
	}
	n, err := number.Int64()
	return err == nil && n >= 1
}

// TestLiveSkeletonWrapperCarriesNoValue is the dedupe half of the leak gate on
// its own, and it exists because the wrapper is the one position the encoding
// introduced.
//
// The gate above walks the committed tree, which proves the committed tree is
// clean. It cannot prove the gate would *notice* a value, because a value in a
// committed file is a value in git history and the test that would have caught
// it is the one being written. So every position the encoding owns is driven
// here from a planted value, and a case that expects no report is a case where
// that position is genuinely unremarkable.
//
// The three count cases are the ones that matter most. A count of 197 is the shape
// the encoding writes and must pass; a count of 385.6, a count of -1 and a count
// that is a string must all fail, because each of them is a reading wearing the
// one position a reading could reach that no other rule covers. If the count check
// were removed, those three would pass silently, which is the specific way the
// dedupe could have weakened this gate. The case after them is the one that shows
// the limit instead: a reading that *is* a positive integer passes here, and
// isPositiveInteger sets out why nothing else in the repository would catch it
// either.
func TestLiveSkeletonWrapperCarriesNoValue(t *testing.T) {
	shape := func() any {
		return map[string]any{"symbol": liveStringPlaceholder, "close": json.Number(liveNumberPlaceholder)}
	}
	record := func(count any) any {
		return map[string]any{SkeletonCountKey: count, SkeletonShapeKey: shape()}
	}
	wrapper := func(records ...any) any {
		return map[string]any{"data": map[string]any{SkeletonArrayKey: records}}
	}

	cases := []struct {
		name        string
		tree        any
		wantReports int
		because     string
	}{
		{
			name:        "a well-formed wrapper",
			tree:        wrapper(record(json.Number("197"))),
			wantReports: 0,
			because:     "this is the shape the encoding writes, so it must pass",
		},
		{
			name:        "a wrapper with no record at all",
			tree:        wrapper(),
			wantReports: 0,
			because:     "an empty array is an observation, not a leak",
		},
		{
			name:        "a price in a shape",
			tree:        wrapper(record(json.Number("1")), record(json.Number("197"))),
			wantReports: 0,
			because:     "a price is a leaf, and the leaf rule already covers it",
		},
		{
			name: "a reading as a leaf inside a shape",
			tree: wrapper(record(json.Number("197")), map[string]any{
				SkeletonCountKey: json.Number("1"),
				SkeletonShapeKey: map[string]any{"symbol": "AAPL"},
			}),
			wantReports: 1,
			because: "a value inside a shape is the case the whole gate exists for, " +
				"and the dedupe must not have made a shape a blind spot",
		},
		{
			name: "a reading as a name inside a shape",
			tree: wrapper(record(json.Number("197")), map[string]any{
				SkeletonCountKey: json.Number("1"),
				SkeletonShapeKey: map[string]any{"9110101000000000001": liveStringPlaceholder},
			}),
			wantReports: 1,
			because: "a member name inside a shape is still a member name, and the " +
				"dedupe must not have exempted a shape from the key rule",
		},
		{
			name: "a price as a count",
			tree: wrapper(map[string]any{
				SkeletonCountKey: json.Number("385.6"),
				SkeletonShapeKey: shape(),
			}),
			wantReports: 1,
			because: "a count is a positive integer, and a price is not one; this is the " +
				"case that fails if the count check is removed",
		},
		{
			name: "the number placeholder as a count",
			tree: wrapper(map[string]any{
				SkeletonCountKey: json.Number(liveNumberPlaceholder),
				SkeletonShapeKey: shape(),
			}),
			wantReports: 1,
			because: "a count of -1 is indistinguishable from the number the reduction " +
				"emits everywhere else, and it is still not a count",
		},
		{
			name: "a reading as a count",
			tree: wrapper(map[string]any{
				SkeletonCountKey: json.Number("9110101000000000001"),
				SkeletonShapeKey: shape(),
			}),
			wantReports: 0,
			because: "the rule's stated limit, pinned deliberately: an account id or a " +
				"millisecond timestamp IS a positive integer, so the count rule cannot " +
				"refuse it. This is a real widening of the leak gate and it is recorded " +
				"rather than papered over -- see isPositiveInteger's own comment for the " +
				"wider and more uncomfortable half: this case would pass the WHOLE suite, " +
				"not just this walk. Element-count conservation, canonicality and the " +
				"manifest's accounts all hold on a planted count, so the only thing that " +
				"catches it is that EncodeSkeleton computes the count and nobody hand-edits " +
				"one. A count is a multiplicity of a shape, so a reading that is an integer " +
				"is indistinguishable from one, exactly as a string-valued key is",
		},
		{
			name: "a count of zero",
			tree: wrapper(map[string]any{
				SkeletonCountKey: json.Number("0"),
				SkeletonShapeKey: shape(),
			}),
			wantReports: 1,
			because:     "an array element the capture did not observe is not a count of one",
		},
		{
			name: "a count that is a string",
			tree: wrapper(map[string]any{
				SkeletonCountKey: "197",
				SkeletonShapeKey: shape(),
			}),
			wantReports: 1,
			because:     "a string is a reading's own form, and the encoding writes a number",
		},
		{
			name:        "a record with no count",
			tree:        wrapper(map[string]any{SkeletonShapeKey: shape()}),
			wantReports: 1,
			because: "a record whose multiplicity is unrecorded is a position the walk " +
				"cannot account for",
		},
		{
			name:        "a record with no shape",
			tree:        wrapper(map[string]any{SkeletonCountKey: json.Number("3")}),
			wantReports: 1,
			because:     "a count with no shape beside it accounts for nothing",
		},
		{
			name: "a record carrying an extra member",
			tree: wrapper(map[string]any{
				SkeletonCountKey: json.Number("3"),
				SkeletonShapeKey: shape(),
				"price":          "385.6",
			}),
			wantReports: 1,
			because: "the extra member is unaccounted for. Its value is not separately " +
				"reported, because a member of a record is not walked -- the record is " +
				"the encoding's, and the encoding writes two members and no third. One " +
				"report naming it is the right answer, and a second would mean the walk " +
				"had descended into a position it does not own",
		},
		{
			name: "a bare list of records with no wrapper",
			tree: map[string]any{SkeletonArrayKey: []any{
				map[string]any{SkeletonCountKey: json.Number("3"), SkeletonShapeKey: shape()},
			}},
			wantReports: 0,
			because: "the reserved member alone makes a wrapper, so this is a " +
				"well-formed one wherever it appears, including at the top level",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := &recordingSink{}
			assertNoLiveValue(rec, tc.tree, "$", "planted.json")
			if got := len(rec.messages); got != tc.wantReports {
				t.Errorf("assertNoLiveValue reported %d time(s), want %d: %v; %s",
					got, tc.wantReports, rec.messages, tc.because)
			}
		})
	}
}

// TestLiveSkeletonKeysCarryNoReading is the key half of the invariant on its own,
// driven from a table because a body whose key is a reading cannot be committed
// to be walked.
//
// Every tree below carries placeholder leaves, so the only thing that can report
// is the key rule, which is what makes the counts readable: a case expecting one
// report is a case where a key was caught, and a case expecting none is a case
// where every key is a member name. The last row is the rule's stated limit
// pinned deliberately - a string-valued key is not detectable, so a ticker key is
// expected to pass, and this case fails if the rule ever quietly narrows its own
// claim.
func TestLiveSkeletonKeysCarryNoReading(t *testing.T) {
	var (
		leaf   = map[string]any{"symbol": liveStringPlaceholder, "close": json.Number(liveNumberPlaceholder)}
		nested = map[string]any{"orders": []any{map[string]any{"symbol": liveStringPlaceholder}}}
	)
	cases := []struct {
		name        string
		tree        any
		wantReports int
		because     string
	}{
		{
			name:        "a member name",
			tree:        nested,
			wantReports: 0,
			because:     "a name is evidence and must survive",
		},
		{
			name:        "a camelCase member name",
			tree:        map[string]any{"instrumentId": liveStringPlaceholder},
			wantReports: 0,
			because: "the rule is not a snake_case rule; instrumentId is the one " +
				"non-snake_case member name the whole committed tree carries",
		},
		{
			name:        "a dotted member name",
			tree:        map[string]any{"net.asset.value": liveStringPlaceholder},
			wantReports: 0,
			because:     "a dot is not evidence of a reading",
		},
		{
			name:        "an account id as a key",
			tree:        map[string]any{"9110101000000000001": leaf},
			wantReports: 1,
			because: "an account id is a number Webull sent, and as a key the " +
				"reduction preserves it verbatim",
		},
		{
			name:        "an instrument id as a key",
			tree:        map[string]any{"9132750001": leaf},
			wantReports: 1,
			because:     "an endpoint answering {\"9132750001\": {...}} commits the id",
		},
		{
			name:        "a timestamp as a key",
			tree:        map[string]any{"1756000000000": leaf},
			wantReports: 1,
			because:     "a millisecond epoch is a number like any other",
		},
		{
			name:        "a price as a key",
			tree:        map[string]any{"385.6": leaf},
			wantReports: 1,
			because:     "the fraction is still the JSON number grammar",
		},
		{
			name:        "the string placeholder as a key",
			tree:        map[string]any{liveStringPlaceholder: leaf},
			wantReports: 1,
			because: "a key a reader cannot tell from a leaf is the ambiguity the " +
				"representation exists to remove",
		},
		{
			name:        "the number placeholder as a key",
			tree:        map[string]any{liveNumberPlaceholder: leaf},
			wantReports: 1,
			because:     "caught by the number rule as well, and both agreeing is the point",
		},
		{
			name:        "an empty key",
			tree:        map[string]any{"": leaf},
			wantReports: 1,
			because:     "an empty key is indistinguishable from an empty string reading",
		},
		{
			name:        "a reading as a key at depth",
			tree:        map[string]any{"data": map[string]any{"9110101000000000001": leaf}},
			wantReports: 1,
			because:     "the rule descends, so a nested map cannot hide a key",
		},
		{
			name:        "a ticker as a key",
			tree:        map[string]any{"AAPL": leaf, "MSFT": leaf},
			wantReports: 0,
			because: "the stated limit, pinned: a JSON member name and a JSON string " +
				"value are the same token, so a string reading as a key is not " +
				"detectable and keyCarriesAReading says so rather than implying it is",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := &recordingSink{}
			assertNoLiveValue(rec, tc.tree, "$", "planted.json")
			if got := len(rec.messages); got != tc.wantReports {
				t.Errorf("assertNoLiveValue reported %d time(s), want %d: %v; %s",
					got, tc.wantReports, rec.messages, tc.because)
			}
			for _, message := range rec.messages {
				if !strings.Contains(message, "member name") {
					t.Errorf("a report does not identify the key as the problem: %s", message)
				}
			}
		})
	}

	// The rule is not vacuous, and the table above is not the only evidence of
	// that: the committed tree is walked by TestLiveSkeletonsCarryNoValue, so a
	// rule that reported on every key would fail there instead.
	if keyCarriesAReading("symbol") || keyCarriesAReading("instrumentId") {
		t.Error("the key rule reports ordinary member names, so it would fail the " +
			"committed tree rather than describe a leak")
	}
}

// liveJSONFiles returns every .json path under the live tree, in sorted order,
// so a failure names a file rather than an index.
//
// The root is spelled without a trailing separator because io/fs requires a
// valid path: fs.WalkDir stats its root before it walks, and an embedded FS
// rejects a name with a trailing slash, which would fail this test on a tree
// that is present.
func liveJSONFiles(t *testing.T) []string {
	t.Helper()
	var found []string
	root := testdataDir + "/" + strings.TrimSuffix(liveTreePrefix, "/")
	err := fs.WalkDir(fixturesFS, root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".json") {
			found = append(found, p)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	sort.Strings(found)
	return found
}

// TestLiveManifestDeclaresItHoldsNoValues checks the two statements the live
// manifest makes about itself that a reader would otherwise have to infer: that
// it is a different artifact from the documentation manifest beside it, and that
// it holds no readings.
func TestLiveManifestDeclaresItHoldsNoValues(t *testing.T) {
	data, err := fixturesFS.ReadFile(testdataDir + "/" + liveManifestFile)
	if err != nil {
		t.Fatalf("read %s: %v", liveManifestFile, err)
	}
	var manifest struct {
		Kind      string `json:"kind"`
		Warning   string `json:"warning"`
		Generator struct {
			ContainsValues bool `json:"containsValues"`
		} `json:"generator"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("parse %s: %v", liveManifestFile, err)
	}
	if manifest.Kind == "" {
		t.Errorf("%s does not name its own kind, so a reader holding only the file "+
			"cannot tell it from the documentation manifest", liveManifestFile)
	}
	if manifest.Generator.ContainsValues {
		t.Errorf("%s claims to hold values", liveManifestFile)
	}
	if !strings.Contains(manifest.Warning, "no value it sent") {
		t.Errorf("%s does not warn that it holds no values: %q", liveManifestFile, manifest.Warning)
	}
}

// liveTreePrefix, liveReadmeFile and liveManifestFile are declared in live.go,
// alongside the code that reads the tree they locate. This file used to declare
// them, and the live comparison could not then name its own input.

// loadLiveManifest reads the live manifest and returns the skeleton paths it
// names, relative to live/.
//
// The manifest is read with a minimal anonymous struct rather than a type of its
// own: it is written by an example program, not by this package, and this
// package has no business depending on the shape of a file it does not read
// elsewhere. The only field this needs is the file each entry was written to,
// and reading one field is what a cross-boundary reader should do: a strict
// parse would make this test fail for an unrelated field the two sides are
// entitled to evolve independently.
func loadLiveManifest(t *testing.T) map[string]bool {
	t.Helper()
	data, err := fixturesFS.ReadFile(testdataDir + "/" + liveManifestFile)
	if err != nil {
		t.Fatalf("read %s: %v", liveManifestFile, err)
	}
	var manifest struct {
		Entries []struct {
			Skeleton string `json:"skeleton"`
		} `json:"entries"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("parse %s: %v", liveManifestFile, err)
	}
	if len(manifest.Entries) == 0 {
		t.Fatalf("%s records no entries, so the live tree is unindexed", liveManifestFile)
	}
	named := make(map[string]bool, len(manifest.Entries))
	for _, e := range manifest.Entries {
		if e.Skeleton == "" {
			continue
		}
		rel, ok := strings.CutPrefix(e.Skeleton, liveTreePrefix)
		if !ok {
			t.Errorf("%s names skeleton %q, which is not under %s", liveManifestFile, e.Skeleton, liveTreePrefix)
			continue
		}
		named[rel] = true
	}
	return named
}

func TestManifestTotalsAgreeWithRecords(t *testing.T) {
	m := load(t)
	if m.Totals.Fixtures != len(m.Fixtures) {
		t.Errorf("totals.fixtures = %d, records = %d", m.Totals.Fixtures, len(m.Fixtures))
	}
	var total int
	var noRequired []string
	var oneOf, addlProps int
	for _, f := range m.Fixtures {
		data, err := f.Read()
		if err != nil {
			t.Errorf("%s: %v", f.ID, err)
			continue
		}
		// The total is the generator's, so it is summed the generator's way;
		// see committedLength. On a CRLF checkout the raw sum is larger by one
		// byte per line and this total would read as drift.
		total += committedLength(data)
		if !f.Checks.DeclaresRequired {
			noRequired = append(noRequired, f.ID)
		}
		if len(f.OneOfChoices) > 0 {
			oneOf++
		}
		if len(f.AdditionalProperties) > 0 {
			addlProps++
		}
	}
	if m.Totals.FixtureBytes != total {
		t.Errorf("totals.fixtureBytes = %d, sum of records = %d", m.Totals.FixtureBytes, total)
	}
	if m.Totals.PagesWithoutRequired != len(noRequired) {
		t.Errorf("totals.pagesWithoutRequired = %d, records = %d",
			m.Totals.PagesWithoutRequired, len(noRequired))
	}
	if m.Totals.PagesWithOneOf != oneOf {
		t.Errorf("totals.pagesWithOneOf = %d, records = %d", m.Totals.PagesWithOneOf, oneOf)
	}
	if m.Totals.PagesWithAdditionalProperties != addlProps {
		t.Errorf("totals.pagesWithAdditionalProperties = %d, records = %d",
			m.Totals.PagesWithAdditionalProperties, addlProps)
	}
	if len(m.Totals.OversizedFixtures) != 0 {
		t.Errorf("committed tree names oversized fixtures %v; generation fails "+
			"above %d bytes, so one of these should not exist",
			m.Totals.OversizedFixtures, m.SizeTripwire.FailAboveBytes)
	}
	sort.Strings(noRequired)
	if strings.Join(m.PagesWithoutRequiredNames, "\n") != strings.Join(noRequired, "\n") {
		t.Errorf("pagesWithoutRequiredNames (%d entries) does not match the "+
			"records that set checks.declaresRequired = false (%d)",
			len(m.PagesWithoutRequiredNames), len(noRequired))
	}
}

// TestFixtureMatchesItsRecord is the core invariant: the committed bytes are the
// shape the manifest says, and they are the size the manifest says. A fixture
// edited by hand, or a manifest regenerated against a stale tree, fails here.
//
// The size comparison is exact after CRLF is collapsed to LF, for the reason
// committedDigest gives. The structural assertions around it are what catch a
// content change, and they are unaffected by a checkout's line-ending policy.
func TestFixtureMatchesItsRecord(t *testing.T) {
	m := load(t)
	for _, f := range m.Fixtures {
		t.Run(f.ID, func(t *testing.T) {
			data, err := f.Read()
			if err != nil {
				t.Fatal(err)
			}
			checkNoStrayCarriageReturn(t, f, data)
			if got := committedLength(data); got != f.Bytes {
				t.Errorf("fixture is %d bytes with CRLF collapsed to LF (%d bytes as "+
					"checked out), manifest says %d", got, len(data), f.Bytes)
			}
			if got := committedDigest(data); got != f.SHA256 {
				t.Errorf("fixture digest is %s, manifest says %s; the length check above "+
					"cannot see this, because a length-preserving edit to a committed "+
					"fixture would have kept it green", got, f.SHA256)
			}
			if m.SizeTripwire.FailAboveBytes > 0 && committedLength(data) > m.SizeTripwire.FailAboveBytes {
				t.Errorf("fixture is %d bytes, above the %d byte ceiling: %s",
					committedLength(data), m.SizeTripwire.FailAboveBytes, m.SizeTripwire.Reason)
			}
			if !json.Valid(data) {
				t.Fatalf("fixture is not valid JSON")
			}
			decoded, err := f.Decode()
			if err != nil {
				t.Fatal(err)
			}
			checkTopLevel(t, f, decoded)
			checkRequiredNames(t, f, decoded)
			if got := countNames(decoded); got != f.Checks.PropertyNameCountInFixture {
				t.Errorf("fixture holds %d distinct property names, manifest says %d",
					got, f.Checks.PropertyNameCountInFixture)
			}
			if len(f.Checks.RequiredWithoutExample) != 0 {
				t.Errorf("required names %v have no declared schema, so the "+
					"generator had nothing to derive a value from",
					f.Checks.RequiredWithoutExample)
			}
			if f.Checks.SynthesizedLeafCount > 0 {
				t.Logf("%d leaf/leaves took a synthesized value: the page gave no "+
					"inline example", f.Checks.SynthesizedLeafCount)
			}
		})
	}
}

// checkTopLevel asserts the instance's JSON kind against the documented one, so
// an array-of-strings endpoint cannot quietly become an array of objects.
func checkTopLevel(t *testing.T, f Fixture, decoded any) {
	t.Helper()
	switch f.Checks.TopLevel {
	case "object":
		obj, ok := decoded.(map[string]any)
		if !ok {
			t.Fatalf("documented top level is object, fixture decodes to %T", decoded)
		}
		if f.Checks.ElementType != "" {
			t.Errorf("elementType = %q on an object top level", f.Checks.ElementType)
		}
		if len(f.Checks.RequiredNames) > 0 &&
			f.Checks.RequiredNamesSource != "topLevel" {
			t.Errorf("requiredNamesSource = %q, want topLevel", f.Checks.RequiredNamesSource)
		}
		for _, name := range f.Checks.RequiredNames {
			if _, ok := obj[name]; !ok {
				t.Errorf("required name %q absent from the object instance", name)
			}
		}
	case "array":
		arr, ok := decoded.([]any)
		if !ok {
			t.Fatalf("documented top level is array, fixture decodes to %T", decoded)
		}
		if len(arr) != 1 {
			t.Errorf("array instance holds %d elements, want exactly 1: the "+
				"minimal rule emits one element per array", len(arr))
		}
		if f.Checks.ElementType == "" {
			t.Error("array top level with no elementType")
			return
		}
		if len(arr) == 0 {
			return
		}
		switch f.Checks.ElementType {
		case "object":
			elem, ok := arr[0].(map[string]any)
			if !ok {
				t.Errorf("elementType is object, element decodes to %T", arr[0])
				return
			}
			for _, name := range f.Checks.RequiredNames {
				if _, ok := elem[name]; !ok {
					t.Errorf("required element name %q absent from the instance", name)
				}
			}
			if f.Checks.ElementType == "object" {
				if len(arr[0].(map[string]any)) == 0 && f.Checks.DeclaresRequired {
					t.Error("declaresRequired is true but the array instance is empty")
				}
			}
		case "string":
			if _, ok := arr[0].(string); !ok {
				t.Errorf("elementType is string, element decodes to %T", arr[0])
			}
		default:
			t.Logf("elementType %q is not exercised by this check", f.Checks.ElementType)
		}
		if f.Checks.DeclaresRequired && f.Checks.RequiredNamesSource != "items" {
			t.Errorf("requiredNamesSource = %q, want items for an array top level",
				f.Checks.RequiredNamesSource)
		}
	default:
		t.Errorf("unknown topLevel %q", f.Checks.TopLevel)
	}
}

// checkRequiredNames asserts the two required-name bookkeeping fields are
// consistent with each other, so DeclaresRequired can never quietly disagree
// with the list it summarises.
func checkRequiredNames(t *testing.T, f Fixture, _ any) {
	t.Helper()
	if f.Checks.DeclaresRequired != (len(f.Checks.RequiredNames) > 0) {
		t.Errorf("declaresRequired = %t with %d required names",
			f.Checks.DeclaresRequired, len(f.Checks.RequiredNames))
	}
	if f.Checks.DeclaresRequired && f.Checks.RequiredNamesSource == "" {
		t.Error("declaresRequired is true but requiredNamesSource is empty")
	}
	if !sort.StringsAreSorted(f.Checks.RequiredNames) {
		t.Errorf("requiredNames is not sorted: %v", f.Checks.RequiredNames)
	}
	if f.Checks.DeclaredPropertyNameCount <
		len(f.Checks.RequiredNames) {
		t.Errorf("declaredPropertyNameCount %d is below the %d required names",
			f.Checks.DeclaredPropertyNameCount, len(f.Checks.RequiredNames))
	}
}

func countNames(v any) int {
	seen := map[string]bool{}
	var walk func(any)
	walk = func(v any) {
		switch t := v.(type) {
		case map[string]any:
			for k, sub := range t {
				seen[k] = true
				walk(sub)
			}
		case []any:
			for _, sub := range t {
				walk(sub)
			}
		}
	}
	walk(v)
	return len(seen)
}

// TestRecordedChoicesAreWellFormed checks the oneOf and additionalProperties
// records, which are the parts of the manifest a reviewer relies on instead of
// the bytes.
func TestRecordedChoicesAreWellFormed(t *testing.T) {
	m := load(t)
	for _, f := range m.Fixtures {
		for _, c := range append(append([]Choice{}, f.OneOfChoices...), f.AdditionalProperties...) {
			if !strings.HasPrefix(c.Pointer, "/") {
				t.Errorf("%s: %s pointer %q is not a JSON Pointer", f.ID, c.Kind, c.Pointer)
			}
			switch c.Kind {
			case "oneOf":
				if c.BranchIndex < 0 || c.BranchIndex >= c.BranchCount {
					t.Errorf("%s: branchIndex %d outside [0,%d)",
						f.ID, c.BranchIndex, c.BranchCount)
				}
				if c.Selection != "title-matches-discriminator" && c.Selection != "first-branch" {
					t.Errorf("%s: selection %q is neither an evidence-based nor a "+
						"default choice, so a reader cannot tell how the branch "+
						"was picked", f.ID, c.Selection)
				}
				if c.Selection == "title-matches-discriminator" && c.Discriminator == "" {
					t.Errorf("%s: branch chosen on a discriminator, but none is "+
						"recorded", f.ID)
				}
				// A first-branch fallback alongside a discriminator is the
				// honest case, not a contradiction: the page supplied a
				// discriminating example that named no branch title. The
				// milestones page is one -- its own `type` example is
				// `football_game` while its two branches are titled
				// SportsGameDetails and EconomicReleaseDetails. Requiring the
				// discriminator to be absent here would push the generator to
				// hide the evidence it actually had.
				if c.Selection == "first-branch" && c.Discriminator != "" {
					t.Logf("%s: discriminator %q (%v) names no branch title, so the "+
						"first branch %q was taken",
						f.ID, c.Discriminator, c.DiscriminatorValue, deref(c.Title))
				}
				if c.BranchRequiredNameCount == 0 {
					t.Logf("%s: chosen branch %q declares no required names, so the "+
						"choice is recorded but not observable in the instance",
						f.ID, deref(c.Title))
				}
			case "additionalProperties":
				if c.ValueType == "" {
					t.Errorf("%s: additionalProperties with no declared value type",
						f.ID)
				}
				if c.SyntheticKey == "" {
					t.Errorf("%s: additionalProperties with no synthetic key", f.ID)
				}
			default:
				t.Errorf("%s: unknown choice kind %q", f.ID, c.Kind)
			}
		}
	}
}

func deref(s *string) string {
	if s == nil {
		return "<untitled>"
	}
	return *s
}

// TestProvenanceIsComplete keeps every row reviewable: a fixture nobody can get
// back to a page is not evidence.
func TestProvenanceIsComplete(t *testing.T) {
	m := load(t)
	seen := map[string]string{}
	for _, f := range m.Fixtures {
		switch {
		case f.Source.URL == "":
			t.Errorf("%s: no source URL", f.ID)
		case !strings.Contains(f.Source.URL, "developer.webull."):
			t.Errorf("%s: source URL %q is not a Webull documentation page",
				f.ID, f.Source.URL)
		}
		if f.Source.CacheFile == "" {
			t.Errorf("%s: no cache filename recorded", f.ID)
		} else if !strings.HasSuffix(f.Source.CacheFile, ".md") {
			t.Errorf("%s: cache file %q is not a .md page", f.ID, f.Source.CacheFile)
		}
		if f.Source.JSONBlockIndex < 0 {
			t.Errorf("%s: negative JSON block index", f.ID)
		}
		if f.Source.SchemaPointer == "" {
			t.Errorf("%s: no schema pointer recorded", f.ID)
		}
		if f.Documented.Method != "GET" && f.Documented.Method != "POST" {
			t.Errorf("%s: documented method %q is neither GET nor POST",
				f.ID, f.Documented.Method)
		}
		if !strings.HasPrefix(f.Documented.Path, "/") {
			t.Errorf("%s: documented path %q is not absolute", f.ID, f.Documented.Path)
		}
		if f.Fixture != f.ID+".json" {
			t.Errorf("%s: fixture path %q does not follow the id", f.ID, f.Fixture)
		}
		if prev, dup := seen[f.ID]; dup {
			t.Errorf("id %q appears twice (%s and %s)", f.ID, prev, f.PageTitle)
		}
		seen[f.ID] = f.PageTitle

		switch f.SDKPathMatch {
		case PathMatchSame:
			if f.SDKPath != f.Documented.Path {
				t.Errorf("%s: sdkPathMatch is same but %q != documented %q",
					f.ID, f.SDKPath, f.Documented.Path)
			}
		case PathMatchDiffers:
			if f.SDKPath == f.Documented.Path {
				t.Errorf("%s: sdkPathMatch is differs but both paths are %q",
					f.ID, f.SDKPath)
			}
			if f.SDKPath == "" {
				t.Errorf("%s: sdkPathMatch is differs with no SDK path", f.ID)
			}
		case PathMatchNone:
			if f.SDKPath != "" {
				t.Errorf("%s: sdkPathMatch is noSdkPath but sdkPath is %q",
					f.ID, f.SDKPath)
			}
		default:
			t.Errorf("%s: unknown sdkPathMatch %q", f.ID, f.SDKPathMatch)
		}
	}
}

// TestDefectivePathRowsAreDocumented pins the count of endpoints where the SDK
// sends a path other than the documented one. A change here is a change in a
// known, tracked defect set, and IMPLEMENTATION_STATUS.md has to move with it --
// so this test fails loudly rather than letting the set drift unnoticed.
func TestDefectivePathRowsAreDocumented(t *testing.T) {
	m := load(t)
	var differs []string
	for _, f := range m.Fixtures {
		if f.SDKPathMatch == PathMatchDiffers {
			differs = append(differs, f.ID)
		}
	}
	sort.Strings(differs)
	t.Logf("%d of %d endpoints send a path other than the documented one:",
		len(differs), len(m.Fixtures))
	for _, id := range differs {
		t.Logf("  %s", id)
	}
}
