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
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// These tests hold the committed form of the live tree, which is a deduplicated
// encoding of the reduced trees examples/live-probe produces. The invariant they
// protect is the one the whole harness rests on: the committed bytes and the
// expanded bytes carry the same member names and the same JSON kinds, so a
// finding about the SDK is a finding about what the sandbox sent.
//
// A test in this file is worthless if it would still pass with the comparison
// inverted, so each one below has a stated way to break it and the breakage is
// checked to actually produce a failure.

// decodeSkeletons parses committed or expanded skeleton bytes the way the
// comparison reads them: with UseNumber, so a number keeps the literal it was
// written with rather than a float64 round trip.
func decodeSkeletons(t *testing.T, raw []byte) any {
	t.Helper()
	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		t.Fatalf("parse skeleton: %v", err)
	}
	return v
}

// compact encodes v the way two shapes are compared for equality: one line, no
// spaces, object members in sorted order. encoding/json sorts map keys, so this
// is a function of the value alone.
func compact(t *testing.T, v any) string {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(raw)
}

// skeletonArray returns the dedupe records of the array at the named member path
// of a decoded committed skeleton, failing when the path holds no array.
func skeletonArray(t *testing.T, tree any, path ...string) []any {
	t.Helper()
	cursor := tree
	for _, step := range path {
		obj, ok := cursor.(map[string]any)
		if !ok {
			t.Fatalf("%s: %v is not an object", strings.Join(path, "."), cursor)
		}
		cursor, ok = obj[step]
		if !ok {
			t.Fatalf("no member %q at %s", step, strings.Join(path, "."))
		}
	}
	wrapper, ok := cursor.(map[string]any)
	if !ok {
		t.Fatalf("%s is not a dedupe wrapper", strings.Join(path, "."))
	}
	records, ok := wrapper[SkeletonArrayKey].([]any)
	if !ok {
		t.Fatalf("%s holds no %q list", strings.Join(path, "."), SkeletonArrayKey)
	}
	return records
}

// recordShape returns the expanded shape of one dedupe record and its count.
func recordShape(t *testing.T, record any) (any, int) {
	t.Helper()
	rec, ok := record.(map[string]any)
	if !ok {
		t.Fatalf("record %v is not an object", record)
	}
	count, ok := rec[SkeletonCountKey].(json.Number)
	if !ok {
		t.Fatalf("record %v holds no %q", record, SkeletonCountKey)
	}
	n, err := count.Int64()
	if err != nil {
		t.Fatalf("count %s is not an integer: %v", count, err)
	}
	if n < 1 {
		t.Fatalf("count %d is not positive", n)
	}
	shape, ok := rec[SkeletonShapeKey]
	if !ok {
		t.Fatalf("record %v holds no %q", record, SkeletonShapeKey)
	}
	return shape, int(n)
}

// TestEncodeSkeletonDeduplicatesRepeatedSiblingShapes is the size claim itself.
//
// 92.1% of the live tree's bytes are identical sibling elements, so the encoding
// keeps one record per distinct shape and a count beside it. The way to break
// this test is to emit every element, and the way to break the size win is to
// collapse two genuinely different shapes into one record; both are checked
// here, because a test that only counted bytes would pass on the second.
func TestEncodeSkeletonDeduplicatesRepeatedSiblingShapes(t *testing.T) {
	leaf := map[string]any{"symbol": "1", "close": json.Number("-1")}
	odd := map[string]any{"symbol": "1", "close": json.Number("-1"), "volume": json.Number("-1")}
	repeated := make([]any, 0, 197)
	for i := 0; i < 197; i++ {
		repeated = append(repeated, leaf)
	}

	encoded, err := EncodeSkeleton(map[string]any{"data": repeated})
	if err != nil {
		t.Fatalf("EncodeSkeleton: %v", err)
	}
	records := skeletonArray(t, decodeSkeletons(t, encoded), "data")
	if got, want := len(records), 1; got != want {
		t.Fatalf("an array of %d identical elements produced %d record(s), want %d",
			len(repeated), got, want)
	}
	shape, count := recordShape(t, records[0])
	if count != 197 {
		t.Errorf("count = %d, want 197: the multiplicity is the evidence the size "+
			"reduction costs, and losing it is silent truncation of a different kind", count)
	}
	if !reflect.DeepEqual(shape, leaf) {
		t.Errorf("shape = %v, want %v", shape, leaf)
	}

	// The same array with one odd element must keep two records, so the dedupe
	// cannot be a truncation wearing a count.
	mixed := append([]any{}, repeated...)
	mixed = append(mixed, odd)
	encoded, err = EncodeSkeleton(map[string]any{"data": mixed})
	if err != nil {
		t.Fatalf("EncodeSkeleton: %v", err)
	}
	if got := len(skeletonArray(t, decodeSkeletons(t, encoded), "data")); got != 2 {
		t.Errorf("an array of %d identical elements plus one different produced %d record(s), "+
			"want 2: two distinct shapes are two records", len(mixed), got)
	}
}

// TestEncodeSkeletonKeepsEveryDistinctShape is the losslessness half, and it is
// pinned on the case the review named rather than on a convenient one: in the
// captured futures product-codes response the five distinct shapes first appear
// at indices 0, 42, 221, 953 and 969 of 2527. Any bound low enough to matter
// drops a shape, and nothing in the encoding would report the loss.
//
// The way to break this test is any cap on the number of records, or a truncation
// of the input before dedupe.
func TestEncodeSkeletonKeepsEveryDistinctShape(t *testing.T) {
	firstSeen := []int{0, 42, 221, 953, 969}
	const total = 2527
	shapes := make([]any, total)
	for i := range shapes {
		switch shapeIndexFor(i, firstSeen) {
		case 0:
			shapes[i] = map[string]any{"symbol": "1"}
		case 1:
			shapes[i] = map[string]any{"symbol": "1", "exchange": "1"}
		case 2:
			shapes[i] = map[string]any{"symbol": "1", "exchange": "1", "tick": json.Number("-1")}
		case 3:
			shapes[i] = map[string]any{"symbol": "1", "exchange": "1", "tick": json.Number("-1"), "lot": json.Number("-1")}
		default:
			shapes[i] = map[string]any{"symbol": "1", "exchange": "1", "tick": json.Number("-1"),
				"lot": json.Number("-1"), "currency": "1"}
		}
	}

	encoded, err := EncodeSkeleton(map[string]any{"codes": shapes})
	if err != nil {
		t.Fatalf("EncodeSkeleton: %v", err)
	}
	records := skeletonArray(t, decodeSkeletons(t, encoded), "codes")
	if got, want := len(records), len(firstSeen); got != want {
		t.Fatalf("kept %d distinct shape(s), want %d: a shape is dropped and the "+
			"count of the rest still sums to %d, so nothing would report it",
			got, want, total)
	}

	// Every shape that was in the input is in the encoding, matched on its
	// expansion rather than on its record, so a record holding the wrong shape
	// fails even when the record count is right.
	want := map[string]bool{}
	for _, s := range shapes {
		want[compact(t, s)] = true
	}
	got := map[string]bool{}
	for _, r := range records {
		shape, _ := recordShape(t, r)
		got[compact(t, shape)] = true
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("the encoded shapes are not the input shapes.\n got %v\nwant %v", got, want)
	}
}

// shapeIndexFor is the shape an element at index i takes, given the first index
// each shape appears at. It reproduces the captured distribution rather than
// interleaving the shapes, because an interleaved array is the easy case and the
// clustered one is what the sandbox actually sent.
func shapeIndexFor(i int, firstSeen []int) int {
	index := 0
	for n, at := range firstSeen {
		if i >= at {
			index = n
		}
	}
	return index
}

// TestEncodeSkeletonIsCanonicalUnderReordering is the determinism claim: the
// committed bytes are a function of the shape evidence alone, so a server that
// returns the same names in a different order produces a diff of nothing.
//
// The way to break this test is to preserve input order in the encoding, which
// would make every capture run a diff against the last one for reasons that
// carry no evidence.
func TestEncodeSkeletonIsCanonicalUnderReordering(t *testing.T) {
	a := map[string]any{"symbol": "1", "close": json.Number("-1")}
	b := map[string]any{"symbol": "1", "open": json.Number("-1")}
	c := map[string]any{"symbol": "1", "high": json.Number("-1"), "low": json.Number("-1")}

	forward, err := EncodeSkeleton(map[string]any{"data": []any{a, b, c, a, c}})
	if err != nil {
		t.Fatalf("EncodeSkeleton: %v", err)
	}
	shuffled, err := EncodeSkeleton(map[string]any{"data": []any{c, a, a, c, b}})
	if err != nil {
		t.Fatalf("EncodeSkeleton: %v", err)
	}
	if string(forward) != string(shuffled) {
		t.Errorf("reordering the elements changed the committed bytes.\nforward %s\nshuffled %s",
			forward, shuffled)
	}

	// The records are in a defined order rather than an incidental one, so the
	// assertion above is a statement about a rule and not about one run.
	records := skeletonArray(t, decodeSkeletons(t, forward), "data")
	var order []string
	for _, r := range records {
		shape, _ := recordShape(t, r)
		order = append(order, compact(t, shape))
	}
	if !sortStringsAscending(order) {
		t.Errorf("records are not in ascending order of their expanded shape: %v", order)
	}
}

// sortStringsAscending reports whether s is sorted, and is a test helper rather
// than sort.StringsAreSorted so the failure message can name the order.
func sortStringsAscending(s []string) bool {
	for i := 1; i < len(s); i++ {
		if s[i-1] > s[i] {
			return false
		}
	}
	return true
}

// TestEncodeSkeletonRefusesTheReservedNames is the ambiguity guard.
//
// A dedupe wrapper is an object with a reserved member, so a Webull response
// carrying a member of the same name would be indistinguishable from the
// encoding. No member name in the committed tree begins with the reserved
// prefix, but "no name today" is not a property the encoding may depend on, so
// the encoder refuses rather than hopes.
//
// The way to break this test is to encode such a member, after which a live
// skeleton and a committed one stop meaning the same thing.
func TestEncodeSkeletonRefusesTheReservedNames(t *testing.T) {
	for _, name := range []string{SkeletonArrayKey, SkeletonCountKey, SkeletonShapeKey} {
		t.Run(name, func(t *testing.T) {
			tree := map[string]any{"data": []any{map[string]any{name: "1"}}}
			encoded, err := EncodeSkeleton(tree)
			if err == nil {
				t.Fatalf("a member named %q was encoded, producing %s; the wrapper and "+
					"the response would then be the same object", name, encoded)
			}
			if !strings.Contains(err.Error(), name) {
				t.Errorf("the error does not name the offending member: %v", err)
			}
		})
	}
}

// TestEncodeSkeletonIsIdempotent is the property that makes the committed tree
// checkable at all: encoding what a committed file decodes to must reproduce the
// committed file byte for byte. Without it a regenerated tree is a diff nobody
// can read, and the tree drifts one capture at a time.
//
// The way to break this test is any part of the encoding that is not a function
// of the shape evidence, which is the same class as the reordering test but
// catches the cases that one does not reach.
func TestEncodeSkeletonIsIdempotent(t *testing.T) {
	for _, name := range liveJSONFiles(t) {
		t.Run(name, func(t *testing.T) {
			committed, err := fixturesFS.ReadFile(name)
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			reencoded, err := EncodeSkeleton(decodeSkeletons(t, committed))
			if err != nil {
				t.Fatalf("EncodeSkeleton: %v", err)
			}
			if string(reencoded) != string(committed) {
				t.Errorf("the committed file is not canonical: re-encoding it differs.\n"+
					"committed  %d bytes\nre-encoded %d bytes", len(committed), len(reencoded))
			}
		})
	}
}

// TestExpandSkeletonIsTheInverseOfEncode is the losslessness claim stated over
// the whole tree rather than over one array.
//
// The reduction to distinct shapes is lossless for member names and JSON kinds,
// and lossless for multiplicity, because the counts are written down. This walks
// every committed file, expands it, and compares the result with the tree the
// file encodes: same names, same kinds, same element count at every position.
//
// The way to break this test is to drop a count, to drop a record, or to expand
// a record to the wrong shape, and each of those is silent in the bytes.
func TestExpandSkeletonIsTheInverseOfEncode(t *testing.T) {
	for _, name := range liveJSONFiles(t) {
		t.Run(name, func(t *testing.T) {
			committed, err := fixturesFS.ReadFile(name)
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			expanded, err := ExpandSkeleton(committed)
			if err != nil {
				t.Fatalf("ExpandSkeleton: %v", err)
			}
			reencoded, err := EncodeSkeleton(decodeSkeletons(t, expanded))
			if err != nil {
				t.Fatalf("EncodeSkeleton(expanded): %v", err)
			}
			if string(reencoded) != string(committed) {
				t.Errorf("expanding the committed file and re-encoding it did not "+
					"reproduce the file, so the encoding loses or gains evidence.\n"+
					"committed %d bytes\nround trip %d bytes", len(committed), len(reencoded))
			}
		})
	}
}

// TestExpandSkeletonReproducesTheElementCount is the multiplicity claim, pinned
// on the tree's own numbers rather than on a synthetic array. The live manifest
// records how many array elements the captured responses held; a committed tree
// that expands to fewer is a truncated tree that still parses, and the counts
// beside the shapes are the only thing that would say so.
//
// The way to break this test is to write a count that does not match the elements
// it stands for, or to drop a record whose count was the only thing keeping the
// total.
func TestExpandSkeletonReproducesTheElementCount(t *testing.T) {
	manifest := readLiveManifest(t)
	entries := manifest.Entries
	if len(entries) == 0 {
		t.Fatal("the live manifest records no entries, so this test is asserting nothing")
	}
	if manifest.Size.ArrayElements == 0 {
		t.Fatal("the live manifest records no array elements, so this test is asserting nothing")
	}

	var elements int
	distinct := map[string]bool{}
	for _, e := range entries {
		raw, err := fixturesFS.ReadFile(testdataDir + "/" + e.Skeleton)
		if err != nil {
			t.Errorf("%s: %v", e.Skeleton, err)
			continue
		}
		expanded, err := ExpandSkeleton(raw)
		if err != nil {
			t.Errorf("%s: ExpandSkeleton: %v", e.Skeleton, err)
			continue
		}
		elements += countArrayElements(decodeSkeletons(t, expanded))
		// Unioned, not summed: the manifest counts distinct shapes over the whole
		// tree, and two files holding the same shape contribute one.
		for shape := range SurveySkeleton(raw).DistinctElementShapes {
			distinct[shape] = true
		}
	}
	if elements != manifest.Size.ArrayElements {
		t.Errorf("the committed tree expands to %d array element(s), the capture "+
			"recorded %d: the encoding is not lossless", elements, manifest.Size.ArrayElements)
	}
	if len(distinct) != manifest.Size.DistinctElementShapes {
		t.Errorf("the committed tree holds %d distinct element shape(s), the capture "+
			"recorded %d", len(distinct), manifest.Size.DistinctElementShapes)
	}
}

// countArrayElements counts every element of every array in an expanded tree.
func countArrayElements(v any) int {
	switch node := v.(type) {
	case map[string]any:
		total := 0
		for _, member := range node {
			total += countArrayElements(member)
		}
		return total
	case []any:
		total := len(node)
		for _, element := range node {
			total += countArrayElements(element)
		}
		return total
	default:
		return 0
	}
}

// TestExpandSkeletonRejectsAMalformedCommittedFile is the refusal half.
//
// The expander is the only path from the committed bytes to the comparison, so a
// committed file it cannot read must be an error rather than a partial tree. A
// silently skipped record would be a finding that vanishes, which is the failure
// mode the whole package exists to end.
//
// The way to break this test is to return a partial expansion on any error.
func TestExpandSkeletonRejectsAMalformedCommittedFile(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "not JSON at all",
			raw:  `<html>gateway error</html>`,
			want: "parse",
		},
		{
			name: "a wrapper that is not a list",
			raw:  `{"data": {"` + SkeletonArrayKey + `": {"symbol": "1"}}}`,
			want: "not a list",
		},
		{
			name: "a record that is not an object",
			raw:  `{"data": {"` + SkeletonArrayKey + `": ["nope"]}}`,
			want: "not an object",
		},
		{
			name: "a wrapper with no count",
			raw:  `{"data": {"` + SkeletonArrayKey + `": [{"` + SkeletonShapeKey + `": {"symbol": "1"}}]}}`,
			want: SkeletonCountKey,
		},
		{
			name: "a count that is not a number",
			raw:  `{"data": {"` + SkeletonArrayKey + `": [{"` + SkeletonCountKey + `": "many", "` + SkeletonShapeKey + `": {"symbol": "1"}}]}}`,
			want: "integer",
		},
		{
			name: "a count of zero",
			raw:  `{"data": {"` + SkeletonArrayKey + `": [{"` + SkeletonCountKey + `": 0, "` + SkeletonShapeKey + `": {"symbol": "1"}}]}}`,
			want: "positive",
		},
		{
			name: "a record with no shape",
			raw:  `{"data": {"` + SkeletonArrayKey + `": [{"` + SkeletonCountKey + `": 3}]}}`,
			want: SkeletonShapeKey,
		},
		{
			name: "a shape carrying a reserved name",
			raw:  `{"data": {"` + SkeletonArrayKey + `": [{"` + SkeletonCountKey + `": 3, "` + SkeletonShapeKey + `": {"` + SkeletonCountKey + `": "1"}}]}}`,
			want: SkeletonCountKey,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expanded, err := ExpandSkeleton([]byte(tc.raw))
			if err == nil {
				t.Fatalf("a malformed committed file expanded to %s", expanded)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error does not mention %q: %v", tc.want, err)
			}
		})
	}
}

// TestSurveySkeletonCountsWhatTheManifestRecords is the accounting claim, over
// the two numbers that must not move.
//
// live-manifest.json records the array elements and the distinct element shapes
// of the tree. The dedupe changes the byte count, which is expected, but it must
// change neither of these: a tree whose element count no longer agrees with its
// own index is a tree whose index describes something else, and nothing else in
// the package would notice.
//
// The way to break this test is to change the encoding without checking the
// numbers survive, which is exactly the drift this catches.
func TestSurveySkeletonCountsWhatTheManifestRecords(t *testing.T) {
	manifest := readLiveManifest(t)
	if manifest.Size.ArrayElements == 0 {
		t.Fatal("the live manifest records no array elements, so this test is asserting nothing")
	}

	var elements int
	distinct := map[string]bool{}
	var carries []string
	for _, e := range manifest.Entries {
		raw, err := fixturesFS.ReadFile(testdataDir + "/" + e.Skeleton)
		if err != nil {
			t.Errorf("%s: %v", e.Skeleton, err)
			continue
		}
		survey := SurveySkeleton(raw)
		elements += survey.ArrayElements
		for shape := range survey.DistinctElementShapes {
			distinct[shape] = true
		}
		if !survey.CarriesElementShape {
			carries = append(carries, e.Skeleton)
		}
	}
	if elements != manifest.Size.ArrayElements {
		t.Errorf("the committed tree holds %d array element(s), the manifest records %d",
			elements, manifest.Size.ArrayElements)
	}
	if len(distinct) != manifest.Size.DistinctElementShapes {
		t.Errorf("the committed tree holds %d distinct element shape(s), the manifest records %d",
			len(distinct), manifest.Size.DistinctElementShapes)
	}
	// The four envelopes that carried no row must still be identified as such:
	// deduping an empty array must not make it look like an array with an element.
	sortStrings(carries)
	if strings.Join(carries, "\n") != strings.Join(manifest.Totals.WithoutElementShape, "\n") {
		t.Errorf("the skeletons holding no element shape are\n%v\nthe manifest records\n%v",
			carries, manifest.Totals.WithoutElementShape)
	}
}

// sortStrings sorts in place, so a test can compare two lists without the
// comparison depending on the order it walked the manifest in.
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
