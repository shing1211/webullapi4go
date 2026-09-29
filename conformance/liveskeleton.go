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
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// # Why the committed skeletons are not the reduced trees
//
// The tree examples/live-probe writes is the shape of what the sandbox answered,
// member for member. It is also 1.97 MB, and 92.1% of those bytes are identical
// sibling elements: the sandbox answers a screener or an instrument list with
// hundreds of rows drawn from a handful of shapes, and the reduction keeps every
// row because a reduction that dropped rows would be a truncation nobody could
// detect.
//
// So the tree stores the evidence rather than the repetition. Every array in the
// committed form is a list of records, one per distinct element shape, each
// carrying how many elements took that shape:
//
//	"data": {"$array": [{"$count": 997, "$shape": {...}}, {"$count": 3, "$shape": {...}}]}
//
// Three properties make this lossless for what the harness reads, and each is
// stated here because the size win is worthless without them:
//
//   - Every distinct shape is kept. There is no cap and no truncation. Nine arrays
//     in the captured tree are heterogeneous, and in high-dividend-ranks the odd
//     element first appears at index 144 of 200; any bound low enough to matter
//     drops a shape, and because the counts still sum to the original length,
//     nothing in the encoding would report the loss.
//   - The multiplicity is written down. A count of 997 beside one shape carries
//     what 997 copies of it carried, so an array's length survives the reduction.
//   - The order is canonical rather than observed. Records are sorted by their
//     expanded encoding, so a server that returns the same names in a different
//     order produces a diff of nothing. Order was never evidence: a list of
//     identical shapes in a different order is not a different list.
//
// # What is lost, stated plainly
//
// The original interleaving of shapes inside an array is not recoverable from the
// committed form, and the same is true of the expanded tree ExpandSkeleton
// returns: a record's count copies are written together, so an array that was
// [A, B, A, A, B] expands to [A, A, A, B, B]. Nothing in the bytes records where
// the B's sat.
//
// That is stated because a reader must not assume the opposite. What *is*
// preserved exactly is the multiset: every distinct shape, and how many elements
// took it, which is what the comparison reads and what the element-count
// conservation check pins. Order within one array was never a property the
// harness used, and a server that reordered its rows must produce a diff of
// nothing, which is the property above. But "the tree does not preserve order" is
// a different sentence from "the tree does not care about order", and only the
// second is true.
//
// The committed bytes are therefore not a CompareBody input. ExpandSkeleton turns
// them back into the reduced tree, and that tree is what the comparison reads, so
// the value-free guarantee and the leak gate are unchanged by the encoding: every
// leaf is still the placeholder for its kind, and the gate in fixtures_test.go
// walks the committed bytes and checks them.

// The three members the dedupe encoding owns. A reduced tree is an object whose
// members are Webull's names, so a marker that could collide with one would make
// a committed skeleton ambiguous with a live response. The reserved prefix is
// chosen for that reason alone: no Webull member name in the committed tree
// begins with it, and EncodeSkeleton refuses any that does rather than depending
// on that staying true.
const (
	// SkeletonArrayKey wraps an array in the committed form.
	SkeletonArrayKey = "$array"
	// SkeletonCountKey is how many elements took the shape beside it.
	SkeletonCountKey = "$count"
	// SkeletonShapeKey is the element shape the count applies to.
	SkeletonShapeKey = "$shape"
	// SkeletonReservedPrefix is the prefix that makes the three keys above
	// unambiguous, and the rule a member name is held to.
	SkeletonReservedPrefix = "$"
)

// EncodeSkeleton renders a reduced tree as the committed form.
//
// The output is a function of the shape evidence alone: object members are
// written in sorted order by encoding/json, and array records are sorted by
// their expanded encoding, so the same response always produces the same bytes
// and a diff after a re-capture means the server sent something different.
//
// It returns an error for a member name carrying the reserved prefix, because
// such a name cannot be written: the encoder would produce a document in which
// the response and the encoding of the response are the same object, and the
// expander would silently reinterpret the response's own member.
func EncodeSkeleton(reduced any) ([]byte, error) {
	encoded, err := encodeValue(reduced, "$")
	if err != nil {
		return nil, err
	}
	out, err := json.MarshalIndent(encoded, "", " ")
	if err != nil {
		return nil, fmt.Errorf("conformance: encode skeleton: %w", err)
	}
	return append(out, '\n'), nil
}

// encodeValue is EncodeSkeleton over one node. The three container cases return
// the reduced value itself where the shape needs no wrapper, so a tree with no
// array is byte-for-byte the tree examples/live-probe wrote.
//
// A dedupe wrapper is passed through rather than refused, which is what makes
// EncodeSkeleton idempotent: re-encoding a committed file reproduces it, so a
// migration can be re-run and a canonicality check can be written as
// encode(decode(committed)) == committed. Without that, the reserved-name rule
// would make the committed tree un-encodable and the only way to check it would
// be to trust it.
func encodeValue(v any, path string) (any, error) {
	switch node := v.(type) {
	case map[string]any:
		if records, wrapped := node[SkeletonArrayKey]; wrapped {
			return encodeWrapper(records, path)
		}
		out := make(map[string]any, len(node))
		for name, member := range node {
			if strings.HasPrefix(name, SkeletonReservedPrefix) {
				return nil, fmt.Errorf("conformance: skeleton member %q at %s carries the "+
					"reserved prefix %q, so a response holding it cannot be written without "+
					"being indistinguishable from its own encoding", name, path, SkeletonReservedPrefix)
			}
			encoded, err := encodeValue(member, path+"."+name)
			if err != nil {
				return nil, err
			}
			out[name] = encoded
		}
		return out, nil
	case []any:
		return encodeArray(node, path)
	default:
		return v, nil
	}
}

// encodeWrapper re-encodes a dedupe wrapper, so a committed file passes through
// EncodeSkeleton to itself. The records are re-deduped through encodeArray rather
// than copied, so a wrapper whose records disagreed with each other collapses
// here rather than being preserved.
func encodeWrapper(records any, path string) (any, error) {
	list, ok := records.([]any)
	if !ok {
		return nil, fmt.Errorf("conformance: skeleton at %s holds a %s that is not a list of records",
			path, SkeletonArrayKey)
	}
	var elements []any
	for i, r := range list {
		record, ok := r.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("conformance: skeleton record at %s.%s[%d] is a %T, not an object",
				path, SkeletonArrayKey, i, r)
		}
		count, err := recordCount(record[SkeletonCountKey], path)
		if err != nil {
			return nil, err
		}
		shape, hasShape := record[SkeletonShapeKey]
		if !hasShape {
			return nil, fmt.Errorf("conformance: skeleton record at %s.%s[%d] holds no %s",
				path, SkeletonArrayKey, i, SkeletonShapeKey)
		}
		for n := 0; n < count; n++ {
			elements = append(elements, shape)
		}
	}
	return encodeArray(elements, path)
}

// encodeArray is the dedupe: distinct shapes with a count each, sorted by the
// shape's own encoding so the order is canonical.
func encodeArray(elements []any, path string) (any, error) {
	order := make([]string, 0, len(elements))
	byShape := make(map[string]any, len(elements))
	counts := make(map[string]int, len(elements))
	for _, element := range elements {
		encoded, err := encodeValue(element, path+"[]")
		if err != nil {
			return nil, err
		}
		key, err := shapeKey(encoded)
		if err != nil {
			return nil, err
		}
		if _, seen := byShape[key]; !seen {
			order = append(order, key)
			byShape[key] = encoded
		}
		counts[key]++
	}
	// A stable sort on a computed key is what makes the file a function of the
	// evidence rather than of the response's ordering.
	sort.Strings(order)

	records := make([]any, 0, len(order))
	for _, key := range order {
		records = append(records, map[string]any{
			SkeletonCountKey: counts[key],
			SkeletonShapeKey: byShape[key],
		})
	}
	return map[string]any{SkeletonArrayKey: records}, nil
}

// shapeKey is the encoding two element shapes are compared by. encoding/json
// writes an object's members in sorted order, so the key is a function of the
// shape alone and does not depend on the order the members happen to appear in.
func shapeKey(v any) (string, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("conformance: encode skeleton shape: %w", err)
	}
	return string(raw), nil
}

// ExpandSkeleton reconstructs the reduced tree a committed skeleton encodes, as
// the bytes the comparison reads.
//
// Every record becomes count copies of its shape, in canonical order, so the
// result is the shape evidence of the response with its repetition restored. It
// is a valid CompareBody input: the leaves are the same typed placeholders
// Skeletonify produced, and jsonKind classifies each of them as the kind it
// stands for, which is the property the whole live comparison rests on.
//
// A committed file it cannot read is an error rather than a partial tree. A
// silently skipped record is a finding that vanishes, which is the failure mode
// this package exists to end.
func ExpandSkeleton(committed []byte) ([]byte, error) {
	dec := json.NewDecoder(bytes.NewReader(committed))
	dec.UseNumber()
	var tree any
	if err := dec.Decode(&tree); err != nil {
		return nil, fmt.Errorf("conformance: parse skeleton: %w", err)
	}
	expanded, err := expandValue(tree, "$")
	if err != nil {
		return nil, err
	}
	out, err := json.Marshal(expanded)
	if err != nil {
		return nil, fmt.Errorf("conformance: expand skeleton: %w", err)
	}
	return out, nil
}

// expandValue is ExpandSkeleton over one node, and the mirror of encodeValue: a
// wrapper becomes an array of count copies, and everything else is itself.
//
// The presence of the reserved member is what makes an object a wrapper, not its
// shape. A wrapper whose records are malformed is a corrupted committed file and
// is refused; checking only that the object *looks* like a well-formed wrapper
// would let a corrupt one fall through to the object case and be passed through
// unexpanded, which is a finding that vanishes rather than an error.
func expandValue(v any, path string) (any, error) {
	switch node := v.(type) {
	case map[string]any:
		if records, wrapped := node[SkeletonArrayKey]; wrapped {
			return expandArray(records, path)
		}
		for name := range node {
			if strings.HasPrefix(name, SkeletonReservedPrefix) {
				return nil, fmt.Errorf("conformance: skeleton member %q at %s carries the "+
					"reserved prefix %q outside a dedupe wrapper, so the file cannot be read "+
					"without guessing which it is", name, path, SkeletonReservedPrefix)
			}
		}
		out := make(map[string]any, len(node))
		for name, member := range node {
			expanded, err := expandValue(member, path+"."+name)
			if err != nil {
				return nil, err
			}
			out[name] = expanded
		}
		return out, nil
	case []any:
		out := make([]any, 0, len(node))
		for i, element := range node {
			expanded, err := expandValue(element, fmt.Sprintf("%s[%d]", path, i))
			if err != nil {
				return nil, err
			}
			out = append(out, expanded)
		}
		return out, nil
	default:
		return v, nil
	}
}

// expandArray is the mirror of encodeArray: each record becomes count copies.
func expandArray(records any, path string) (any, error) {
	list, ok := records.([]any)
	if !ok {
		return nil, fmt.Errorf("conformance: skeleton at %s holds a %s that is not a list of records",
			path, SkeletonArrayKey)
	}
	var out []any
	for i, r := range list {
		at := fmt.Sprintf("%s.%s[%d]", path, SkeletonArrayKey, i)
		record, ok := r.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("conformance: skeleton record at %s is a %T, not an object", at, r)
		}
		rawCount, hasCount := record[SkeletonCountKey]
		if !hasCount {
			return nil, fmt.Errorf("conformance: skeleton record at %s holds no %s", at, SkeletonCountKey)
		}
		shape, hasShape := record[SkeletonShapeKey]
		if !hasShape {
			return nil, fmt.Errorf("conformance: skeleton record at %s holds no %s", at, SkeletonShapeKey)
		}
		count, err := recordCount(rawCount, at)
		if err != nil {
			return nil, err
		}
		expanded, err := expandValue(shape, at+"."+SkeletonShapeKey)
		if err != nil {
			return nil, err
		}
		for n := 0; n < count; n++ {
			out = append(out, expanded)
		}
	}
	// An empty array stays an empty slice rather than becoming nil, so a body
	// whose top level is an array of nothing stays distinguishable from one
	// whose kind is wrong.
	if out == nil {
		out = []any{}
	}
	return out, nil
}

// recordCount reads a record's count, requiring a positive integer.
//
// A count of zero would be an element the capture did not observe, and a
// non-integer would be a shape the encoder cannot have written. Both are
// refused rather than coerced, because either would make the expanded array
// shorter than the response was without saying so.
func recordCount(v any, at string) (int, error) {
	number, ok := v.(json.Number)
	if !ok {
		return 0, fmt.Errorf("conformance: skeleton record at %s holds a %T as its %s, not an "+
			"integer element count", at, v, SkeletonCountKey)
	}
	n, err := number.Int64()
	if err != nil {
		return 0, fmt.Errorf("conformance: skeleton record at %s holds %s as its %s, which is not "+
			"an integer element count", at, number, SkeletonCountKey)
	}
	if n < 1 {
		return 0, fmt.Errorf("conformance: skeleton record at %s holds a %s of %d, which is not a "+
			"positive element count", at, SkeletonCountKey, n)
	}
	return int(n), nil
}

// SkeletonSurvey is what one committed skeleton holds, for the live manifest's
// size accounts.
type SkeletonSurvey struct {
	// ArrayElements is every element of every array the skeleton expands to. It
	// is read from the counts, so it is the number the capture observed rather
	// than the number of records the file happens to hold.
	ArrayElements int
	// DistinctElementShapes is the encodings of the distinct element shapes, over
	// the whole file, so a consumer can count the tree-wide set.
	DistinctElementShapes map[string]bool
	// CarriesElementShape is false when the skeleton holds an array with no
	// element anywhere, which is a real observation: the server sent an envelope
	// and no row. The manifest records which files those are.
	CarriesElementShape bool

	// sawArray and sawPopulatedArray record what the probe's survey records: the
	// first says the file held an array at all, the second that some array had an
	// element. Together they decide CarriesElementShape.
	sawArray          bool
	sawPopulatedArray bool
}

// SurveySkeleton reads one committed skeleton and reports what it holds.
//
// The accounts are computed from the *expanded* tree, and that is not an
// implementation detail: the numbers the live manifest records were computed by
// examples/live-probe walking the reduced trees it had just written, so an
// account computed over the deduplicated form would count records rather than
// elements and would disagree with its own index on every nested array. The
// multiplicity has to multiply -- an element repeated 197 times, each carrying a
// three-element array, is 591 inner elements and not three -- and expanding first
// is what makes that come out right without a second traversal to get it wrong.
//
// Expansion is cheap relative to the tree it reads: the whole committed tree is
// about 30 KB and expands to roughly 2 MB in memory, once per call.
func SurveySkeleton(committed []byte) SkeletonSurvey {
	s := SkeletonSurvey{DistinctElementShapes: map[string]bool{}, CarriesElementShape: true}
	expanded, err := ExpandSkeleton(committed)
	if err != nil {
		// An unreadable file is reported as carrying an element, which is the
		// conservative direction: the manifest's list names the skeletons that
		// hold none, and adding one to that list would claim a capture observed
		// nothing when the harness simply could not read the file.
		return s
	}
	var tree any
	if err := decodeNumbered(expanded, &tree); err != nil {
		return s
	}
	walkSurvey(tree, &s)
	// The probe's rule is that a skeleton with no array at all carries an element
	// shape, not that it fails to: the field is fed a list of skeletons whose
	// every array is empty, and adding one that holds a member-name-and-kind shape
	// would tell a reader it holds nothing when it holds something. So the
	// inversion here is deliberate and matches surveySkeleton in the probe:
	// only a skeleton that holds an array and never populated one is listed.
	if s.sawArray && !s.sawPopulatedArray {
		s.CarriesElementShape = false
	}
	return s
}

// decodeNumbered parses a body keeping every number as the literal that was
// written. UseNumber is not tidiness here: a shape key is the encoding of an
// element, so a -1 widened to a float64 and re-encoded is "-1" while a -1 that
// round-tripped through float64 can come back as "-1" or "-1e+00" depending on
// the value, and a key that changes shape with the magnitude makes the distinct
// shape count a function of the values rather than of the shapes.
func decodeNumbered(raw []byte, out any) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	return dec.Decode(out)
}

// walkSurvey accumulates one skeleton's accounts, over an expanded reduced tree.
//
// The probe that writes the tree calls SurveySkeleton rather than keeping its own
// copy of this walk, so the two are one definition rather than two that can
// disagree. That is the direction the dependency already runs: examples/live-probe
// imports this package for the symbol table and for the encoder, and nothing here
// imports it. A reimplementation across that boundary would be a second opinion
// about the same counts, and the disagreement it could produce is a manifest that
// describes a tree other than the one in the repository -- which is exactly the
// class of defect this harness exists to end. The rule is restated here in one
// sentence: an array with an element anywhere in the file carries an element shape,
// and only a file that holds an array and never populated one is listed as
// carrying none.
func walkSurvey(v any, s *SkeletonSurvey) {
	switch node := v.(type) {
	case map[string]any:
		for _, member := range node {
			walkSurvey(member, s)
		}
	case []any:
		s.sawArray = true
		if len(node) > 0 {
			s.sawPopulatedArray = true
		}
		for _, element := range node {
			s.ArrayElements++
			if encoded, err := json.Marshal(element); err == nil {
				s.DistinctElementShapes[string(encoded)] = true
			}
			walkSurvey(element, s)
		}
	}
}
