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
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

// DivergenceKind names one specific disagreement between what a documentation
// page declares and what an SDK type can decode. The kinds are the harness's
// vocabulary; nothing outside this package may invent another, because a
// divergence is identified in the committed baseline by
// (symbol, fixture, kind, name, detail) and an unrecognised kind would make the
// baseline unresolvable rather than wrong.
type DivergenceKind string

const (
	// MissingRequiredName means a name the page marks required is not carried
	// through to a value the caller can read. It covers both halves of that one
	// obligation, and the two halves get distinct Detail text: the name is
	// absent from the SDK's json tags, so a conforming response decodes into a
	// zero value with no error; or the name is absent from the documented
	// instance, so nothing can be checked at all. The first half is the defect
	// class a self round-trip test cannot see and that costs a caller real
	// money; the second half is what a stale or hand-edited fixture looks like,
	// and it never fires on a committed tree.
	MissingRequiredName DivergenceKind = "missing-required-name"

	// MissingDeclaredName means a name the page declares but does not mark
	// required is not carried through to a value the caller can read.
	//
	// It is a distinct kind rather than an extension of MissingRequiredName
	// because the evidence behind it is weaker, and collapsing the two would
	// make the strength unrecoverable from a baseline entry. A required name is
	// a promise: the page says the field is always sent. A declared name is a
	// description of one response, so a name missing from it may be optional,
	// conditionally present, or simply not sent in the example the page chose.
	// The defect it reports is the same silent zero the required half reports,
	// and it is recorded rather than assumed away, because 89 of the 90 pages
	// carrying no required list do publish property names and nothing else
	// examines them.
	//
	// On a row that also records TopLevelMismatch some of these names restate
	// that finding rather than adding to it: when the page documents an
	// envelope and the SDK decodes the bare payload, the envelope's own keys
	// read as missing because the wrapper is absent. They are kept, because a
	// name-level test that told the two apart could not be made mechanical. The
	// manifest records no type per declared name and no property name at all in
	// the committed instances, so nothing distinguishes a wrapper key from a
	// payload field. Suppressing the whole row instead would have discarded
	// real findings beside the restatements, so the overlap is documented here
	// and counted where it is reported.
	MissingDeclaredName DivergenceKind = "missing-declared-name"

	// DeclaredInventoryEmpty means a page carries no required list and declares
	// no property name either, so neither the required nor the declared name
	// check has anything to look at. It is a gap in the evidence base rather
	// than in the SDK, and it is recorded for that reason: filed as a pass, the
	// row would read as examined when no name on the page could be compared.
	DeclaredInventoryEmpty DivergenceKind = "declared-inventory-empty"

	// TopLevelMismatch means the documented response is an array where the SDK
	// decodes an object, or the reverse. One bit of shape decides whether the
	// SDK can hold the documented body at all.
	TopLevelMismatch DivergenceKind = "top-level-shape-mismatch"

	// ElementTypeMismatch means both sides agree the response is an array and
	// disagree on what an element is.
	ElementTypeMismatch DivergenceKind = "element-type-mismatch"

	// LeafTypeMismatch means a name both sides carry has a documented scalar
	// the SDK field cannot hold, so the documented value is dropped or the
	// whole decode fails.
	LeafTypeMismatch DivergenceKind = "leaf-type-mismatch"

	// DecodeFailure means unmarshalling the documented instance into the SDK
	// type returned an error. It is reported last and is the weakest of the
	// checks: a decode that succeeds proves almost nothing, and one that fails
	// is usually the same defect the sharper checks already named.
	DecodeFailure DivergenceKind = "decode-failure"

	// VerbMismatch means the SDK method sends an HTTP verb the page does not
	// document for that path, or that its verb could not be read from the source
	// at all. It is the only kind about the request rather than the response, and
	// it exists because the path check cannot see it: a method that sends the
	// documented path with the wrong verb reconciles as a clean match. See
	// verbs.go.
	VerbMismatch DivergenceKind = "verb-mismatch"
)

// checkOf is the harness's priority order. The two shape kinds share rank 1
// because they are one check: the element type only exists once the top level
// is known to be an array. The declared-name kind sits directly below the
// required-name kind because it reports the same silent-zero defect on weaker
// evidence, so a reader meets the promise-backed finding first.
func checkOf(k DivergenceKind) int {
	switch k {
	case MissingRequiredName:
		return 0
	case MissingDeclaredName, DeclaredInventoryEmpty:
		return 1
	case TopLevelMismatch, ElementTypeMismatch:
		return 2
	case LeafTypeMismatch:
		return 3
	case DecodeFailure:
		return 4
	default:
		return 5
	}
}

// Check is one of the five comparisons, named for reporting. It exists so an
// Outcome can say a check did not apply and why, rather than leaving an absent
// check indistinguishable from a passing one.
type Check string

const (
	// CheckRequiredNames is "does every documented required name reach a json tag".
	CheckRequiredNames Check = "required-name-coverage"
	// CheckDeclaredNames is the same question asked of the names a page declares
	// without marking required, and only runs where the required list is absent.
	// It is weaker evidence, which is why it is a separate check rather than a
	// second arm of CheckRequiredNames: an Outcome has to be able to say that a
	// row was examined against a description rather than against a promise.
	CheckDeclaredNames Check = "declared-name-coverage"
	// CheckShape is "does the documented top level and element type match what
	// the SDK decodes into".
	CheckShape Check = "top-level-shape"
	// CheckLeafTypes is "does every required name the SDK does tag have a
	// compatible type".
	CheckLeafTypes Check = "leaf-type-agreement"
	// CheckDecodes is "does the documented instance unmarshal without error".
	CheckDecodes Check = "decodes-without-error"
	// CheckNameDepth is "was every required name reached at the level the page
	// declares it at". It is a check that does not run, and it is reported as
	// such, because the depth-scoped form of the name check would report names
	// the SDK does decode. See tagsOf.
	CheckNameDepth Check = "required-name-depth"
	// CheckVerb is "does the SDK method send the HTTP verb the page declares". It
	// is the only check here that compares the *request* rather than the response,
	// and it is the one whose absence let broker.UpdateVirtualAccount reconcile
	// as a clean match while issuing PUT where the page documents POST. See
	// verbs.go.
	CheckVerb Check = "verb-agreement"
)

// checks is the report order, and the order the gate's summary groups by.
var checks = []Check{CheckRequiredNames, CheckDeclaredNames, CheckShape, CheckLeafTypes, CheckDecodes}

// checkForKind names the check a divergence belongs to.
func checkForKind(k DivergenceKind) Check {
	switch k {
	case MissingRequiredName:
		return CheckRequiredNames
	case MissingDeclaredName, DeclaredInventoryEmpty:
		return CheckDeclaredNames
	case TopLevelMismatch, ElementTypeMismatch:
		return CheckShape
	case LeafTypeMismatch:
		return CheckLeafTypes
	case DecodeFailure:
		return CheckDecodes
	case VerbMismatch:
		return CheckVerb
	default:
		return Check("unknown:" + string(k))
	}
}

// Divergence is one observed disagreement. Detail is part of the identity the
// gate compares, not decoration: if a check starts reporting the same
// disagreement in different words, the baseline stops reproducing and a human
// has to decide whether the check or the SDK changed. A Detail that merely
// restated the kind would make that impossible to see.
type Divergence struct {
	// Symbol is the SDK method whose type was compared, e.g. brokerfd.GetFDPositions.
	Symbol string `json:"symbol"`
	// Fixture is the fixture ID, e.g.
	// broker-fd-us/GET-broker-assets-positions-list.
	Fixture string `json:"fixture"`
	// Kind is the specific disagreement.
	Kind DivergenceKind `json:"kind"`
	// Name is the documented property the divergence is about, or "" for a
	// disagreement that is about the shape rather than any one name.
	Name string `json:"name,omitempty"`
	// Detail states the disagreement in terms a reviewer can check against the
	// two sources without rerunning anything.
	Detail string `json:"detail"`
}

// Key is the identity the gate matches on. Detail is included deliberately;
// see the type comment.
func (d Divergence) Key() string {
	return strings.Join([]string{d.Symbol, d.Fixture, string(d.Kind), d.Name, d.Detail}, "\x1f")
}

// SkippedCheck records a comparison that could not run, and why. A check that
// does not apply must be visible: the failure mode this instrument exists to
// remove is a finding that quietly stops being looked for.
type SkippedCheck struct {
	Check  Check  `json:"check"`
	Reason string `json:"reason"`
}

// ReasonPathNotContract is why a row is reported not comparable rather than
// divergent.
//
// It is a constant, not a per-row string, so the rows group in a report; the two
// paths that make it apply are on NotComparable.
const ReasonPathNotContract = "the SDK method sends a path other than the one the " +
	"page documents, so the page is not this call's contract and a schema difference " +
	"between them is not a defect claim about this SDK"

// NotComparable says a documented page cannot be treated as the contract for the
// SDK call it is compared against, and carries the two paths so a reader can see
// the mismatch without rerunning anything.
//
// It exists because a schema comparison answers a question that presupposes the
// two sides describe the same HTTP call. Where they do not, the answer is not
// "no divergence" and not "divergence": it is that the question was not asked. A
// row that is not comparable is a coverage statement, and folding it into either
// of the other two is how a wrong finding gets reported as a real defect.
type NotComparable struct {
	// Reason is ReasonPathNotContract. A field rather than a method so a report
	// can print it next to the paths without a type switch.
	Reason string `json:"reason"`
	// DocumentedPath is the path the page declares.
	DocumentedPath string `json:"documentedPath"`
	// SDKPath is the path the SDK method sends.
	SDKPath string `json:"sdkPath"`
}

// Outcome is the full result of comparing one fixture against one Go type,
// including everything the comparison could not decide.
type Outcome struct {
	// Symbol is the manifest subject this comparison was made for.
	Symbol string `json:"symbol"`
	// Fixture is the fixture ID compared.
	Fixture string `json:"fixture"`
	// GoType is the compared type as source text, or "" when the symbol decodes
	// no body.
	GoType string `json:"goType,omitempty"`
	// Shape is what the SDK type accepts at the top level.
	Shape WireShape `json:"-"`
	// NotComparable is set, and the row carries no divergences, when the page and
	// the SDK call do not describe the same endpoint. Distinct from a conforming
	// row, which reports five checks that found nothing.
	NotComparable *NotComparable `json:"notComparable,omitempty"`
	// Skipped lists the checks that did not apply.
	Skipped []SkippedCheck `json:"skipped,omitempty"`
	// Divergences are sorted by check priority, then kind, then name.
	Divergences []Divergence `json:"divergences,omitempty"`
}

// WireShape is what an SDK type accepts at the top level of a JSON body.
type WireShape struct {
	// Go is the type itself, as the symbol table declares it.
	Go reflect.Type `json:"-"`
	// TopLevel is the JSON kind: "object" or "array".
	TopLevel string `json:"topLevel"`
	// Element is the array element type, nil unless TopLevel is "array".
	Element reflect.Type `json:"-"`
	// ElementKind is the JSON kind of Element: "object", "string", "number",
	// "boolean", or "unknown".
	ElementKind string `json:"elementKind,omitempty"`
	// Carrier is the struct type whose json tags can carry documented names.
	// Nil when the SDK type is a map or a scalar, which is the fact that
	// decides whether the required-name check can apply at all.
	Carrier reflect.Type `json:"-"`
}

// Describe derives the wire shape of t. A pointer is transparent because every
// SDK method does `var out T` and hands the transport a *T, so the type that
// matters is the one the method names.
func Describe(t reflect.Type) WireShape {
	if t == nil {
		return WireShape{TopLevel: "none"}
	}
	s := WireShape{Go: t}
	switch t.Kind() {
	case reflect.Pointer:
		return Describe(t.Elem())
	case reflect.Slice, reflect.Array:
		s.TopLevel = "array"
		s.Element = t.Elem()
		s.ElementKind = jsonKindOf(t.Elem())
		s.Carrier = carrierOf(t.Elem())
	case reflect.Map:
		s.TopLevel = "object"
	case reflect.Struct:
		s.TopLevel = "object"
		s.Carrier = t
	case reflect.Interface:
		s.TopLevel = "unknown"
	default:
		s.TopLevel = "scalar"
	}
	return s
}

// carrierOf is the struct whose tags carry names for a decoded body: the
// element struct for an array, the struct itself for an object. A map carries
// every name by construction and declares none, so it is not a carrier.
func carrierOf(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() == reflect.Struct {
		return t
	}
	return nil
}

// jsonKindOf classifies a Go type as the JSON kind a decoder would accept for
// it. It is the mirror of jsonKind, which classifies a decoded JSON value, and
// the two are compared against each other in the leaf check.
func jsonKindOf(t reflect.Type) string {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Map, reflect.Struct:
		return "object"
	case reflect.Interface:
		return "unknown"
	default:
		return "unknown"
	}
}

// jsonKind classifies a value decoded from the committed fixture.
func jsonKind(v any) string {
	switch v.(type) {
	case string:
		return "string"
	case bool:
		return "boolean"
	case float64, json.Number:
		return "number"
	case map[string]any:
		return "object"
	case []any:
		return "array"
	case nil:
		return "null"
	default:
		return "unknown"
	}
}

// tagSite is one json tag discovered in a type tree.
type tagSite struct {
	Field reflect.StructField
	Type  reflect.Type
	// Depth is the number of struct levels below the type the comparison was
	// given: 0 for the type's own fields, 1 for those of its element or of an
	// untagged embedded struct, and so on.
	//
	// Depth is read in two places, and both are worth being precise about. It
	// picks which field a duplicated wire name resolves to, because the
	// shallowest occurrence is the one a documented top-level or element name is
	// referring to, and that choice is load-bearing: the leaf check type-checks
	// the field Depth selects. And it is reported, so a reader can see how much
	// of the name coverage came from a nested field.
	//
	// What Depth does NOT do is scope the check. Nothing consults it to decide
	// whether a name is carried at the depth the page declares it, so a
	// documented top-level name satisfied only by a tag three levels down counts
	// as covered. That makes the missing-required-name count a floor rather than
	// a total; see the nameDepth check and the baseline note.
	Depth int
}

// maxTagDepth bounds the type-tree walk. Nothing in the SDK nests a response
// body deeper than a handful of levels; the bound exists so a self-referential
// type cannot make the walk unbounded.
const maxTagDepth = 12

// fieldEncoding is what encoding/json does with one struct field.
type fieldEncoding int

const (
	// fieldIgnored is not encoded: unexported and not embedded, or tagged "-".
	fieldIgnored fieldEncoding = iota
	// fieldInlined contributes no wire name of its own. encoding/json flattens an
	// anonymous field that has no name in its tag, so the walk descends into it
	// and collects its own fields instead.
	fieldInlined
	// fieldNamed encodes under the returned name.
	fieldNamed
)

// wireNameOf classifies f the way encoding/json does.
//
// The rules are not the same as "does the field have a json tag", and getting
// them wrong is how a harness invents a missing name or misses a real one. A
// field with no tag encodes under its Go name. A field tagged "-" is ignored. An
// anonymous field with no name in its tag is inlined, so its Go type name is not
// a wire name -- but its exported fields are, and they are what the SDK will
// decode into, which is why the walk descends rather than skipping. An unexported
// field is ignored, except an embedded struct, whose exported fields encoding/json
// promotes and will therefore set; skipping those wholesale would report every
// promoted name as missing.
func wireNameOf(f reflect.StructField) (fieldEncoding, string) {
	raw, tagged := f.Tag.Lookup("json")
	if tagged {
		if i := strings.IndexByte(raw, ','); i >= 0 {
			raw = raw[:i]
		}
		switch {
		case raw == "-":
			return fieldIgnored, ""
		case raw != "":
			// An explicitly named field is a normal field, and an unexported one is
			// still never set from the wire.
			if !f.IsExported() {
				return fieldIgnored, ""
			}
			return fieldNamed, raw
		}
	}
	if f.Anonymous {
		// No name of its own, so the struct is inlined: the type's Go name is not a
		// wire name, and its exported fields are promoted, so the walk descends.
		// An embedded field that is unexported is inlined too, which is why the
		// unexported test here is not a skip -- though an unexported embedded field
		// that IS named in its tag is not promoted and cannot be set, so it is
		// reported as ignored above. That reading is conservative in the safe
		// direction: a name it carries is reported as missing rather than covered.
		return fieldInlined, ""
	}
	if !f.IsExported() {
		return fieldIgnored, ""
	}
	return fieldNamed, f.Name
}

// tagsOf collects every json tag reachable from t, keyed by the wire name.
//
// The key set is the union over all depths, which is what makes the
// missing-required-name count a floor: encoding/json flattens a response body
// into one name space, so a name the SDK carries anywhere can answer to a
// documented name, and this function cannot tell whether it was carried where
// the page says. Scoping the walk to the depth a page declares would close the
// gap, and it is not done, for a reason that is a judgement rather than a
// difficulty: a documented top-level name the SDK reaches one level down is not a
// defect a caller can observe, it decodes. Reporting it would trade a floor for a
// class of false positives, and this instrument is built to prefer the one to the
// other. The names that are covered only from below are reported instead, by
// CheckNameDepth, so the floor is a number rather than a caveat.
func tagsOf(t reflect.Type) map[string]tagSite {
	out := map[string]tagSite{}
	seen := map[reflect.Type]bool{}
	var walk func(reflect.Type, int)
	walk = func(t reflect.Type, depth int) {
		if t == nil || depth > maxTagDepth || seen[t] {
			return
		}
		seen[t] = true
		for t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		switch t.Kind() {
		case reflect.Struct:
			for i := range t.NumField() {
				f := t.Field(i)
				switch encoding, name := wireNameOf(f); encoding {
				case fieldIgnored:
					// Nothing is decoded from this field and nothing inside it, so
					// the walk stops here rather than collecting names the wire
					// will never reach.
					continue
				case fieldNamed:
					if prev, ok := out[name]; !ok || prev.Depth > depth {
						out[name] = tagSite{Field: f, Type: f.Type, Depth: depth}
					}
					// A named field is a nested object, so its own names are one
					// level below it.
					walk(f.Type, depth+1)
				case fieldInlined:
					// An inlined field is not a level at all: encoding/json flattens
					// it into the enclosing object, so its names are carried at the
					// enclosing depth. Descending with depth+1 here would report
					// every promoted name as reaching the type from below, which is
					// the opposite of what inlining means.
					walk(f.Type, depth)
				}
			}
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Pointer:
			walk(t.Elem(), depth+1)
		}
	}
	walk(t, 0)
	return out
}

// tagName returns the wire name a field encodes under, reproducing
// encoding/json's fallback to the Go field name. It does not apply the anonymous
// or unexported cases, which wireNameOf decides; it is the single-field helper the
// tests use to locate a field by the name it is known on the wire.
func tagName(f reflect.StructField) string {
	raw, ok := f.Tag.Lookup("json")
	if !ok {
		return f.Name
	}
	name := raw
	if i := strings.IndexByte(raw, ','); i >= 0 {
		name = raw[:i]
	}
	if name == "" {
		return f.Name
	}
	return name
}

// lookupTag finds the site carrying a documented wire name, with the same
// case-insensitive fallback encoding/json applies when no tag matches exactly.
//
// The fallback is not a nicety. A field with no json tag encodes under its Go
// name, and a page that spells the same name snake_case -- `success` against a
// field named Success -- still lands in it, because the decoder folds case on a
// miss. Without this lookup the harness reports such a name as missing, which is
// the same class of error as naming a projection the SDK never decodes: a finding
// produced by the harness's own model of encoding/json rather than by the SDK.
//
// When two tags fold to the same name the shallowest wins, as in tagsOf, and ties
// are broken on the wire name so the answer does not depend on map order. Two tags
// differing only in case cannot both be filled by one key, so this is a choice,
// and a deterministic one.
func lookupTag(tags map[string]tagSite, name string) (tagSite, bool) {
	if site, ok := tags[name]; ok {
		return site, true
	}
	var best tagSite
	bestWire := ""
	found := false
	for _, wire := range sortedTagNames(tags) {
		if !strings.EqualFold(wire, name) {
			continue
		}
		site := tags[wire]
		if !found || site.Depth < best.Depth || (site.Depth == best.Depth && wire < bestWire) {
			best, bestWire, found = site, wire, true
		}
	}
	return best, found
}

// sortedTagNames returns the tag names in a stable order, so a caller that picks
// among several matches reports the same one on every run.
func sortedTagNames(tags map[string]tagSite) []string {
	out := make([]string, 0, len(tags))
	for name := range tags {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// typeName renders t the way a Go reader would name it, without the package
// path, so a message matches what a reviewer sees in godoc.
//
// A type rebuilt from a source-declared envelope is the one case where the Go type
// system has no name to offer: reflect.StructOf produces an anonymous struct, and
// printing that would put the whole field list and its escaped tags into a
// divergence detail. The declared spelling is used instead; see envelopes.go.
func typeName(t reflect.Type) string {
	if t == nil {
		return "<nil>"
	}
	if name, ok := declaredName(t); ok {
		return name
	}
	return t.String()
}

// Compare runs all five checks of one fixture against one Go type, using the
// committed bytes.
//
// It applies the comparability rule first. A fixture whose SDK method sends a
// path other than the one its page documents is reported not comparable and no
// check runs, because the page is not that call's contract; see NotComparable.
// The rule lives here rather than in CompareAll so that it cannot be bypassed by
// calling the comparison directly.
func Compare(f Fixture, symbol string, t reflect.Type) Outcome {
	if nc := compareability(f); nc != nil {
		return Outcome{
			Symbol: symbol, Fixture: f.ID, Shape: Describe(t),
			GoType:        typeName(t),
			NotComparable: nc,
			Skipped:       allChecksSkipped(nc.Reason),
		}
	}
	body, err := f.Read()
	if err != nil {
		return unreadable(f, symbol, t,
			"the committed fixture "+f.Path()+" is not readable from the embedded tree, "+
				"so no check could run",
			err.Error())
	}
	return CompareBody(f, symbol, t, body)
}

// compareability reports why a fixture's page cannot be compared against the SDK
// call it names, or nil when it can.
//
// PathMatchNone is not a case here: it is the em dash, an endpoint mapped to no
// SDK symbol at all, and CompareAll routes those to NotCompared before any type
// is resolved. A row that reaches this function has a type, so the only question
// left is whether the two sides describe the same HTTP call.
func compareability(f Fixture) *NotComparable {
	if f.SDKPathMatch == PathMatchSame {
		return nil
	}
	return &NotComparable{
		Reason:         ReasonPathNotContract,
		DocumentedPath: f.Documented.Path,
		SDKPath:        f.SDKPath,
	}
}

// allChecksSkipped is the five-check skip list for a reason that applies to all
// of them, so a not-comparable row never reads as a partial pass.
func allChecksSkipped(reason string) []SkippedCheck {
	out := make([]SkippedCheck, 0, len(checks))
	for _, c := range checks {
		out = append(out, SkippedCheck{Check: c, Reason: reason})
	}
	return out
}

// CompareBody is Compare against an explicitly supplied body, with no
// comparability rule.
//
// It exists so a test can prove the checks fire: a comparison that has never
// been observed failing is not evidence of anything. The embedded tree stays
// immutable either way, because the substitution is an argument rather than a
// write, and no code on the committed-fixture path calls this with anything
// other than f.Read's bytes. Because a caller supplies the body, this is also the
// only entry point that does not consult Fixture.SDKPathMatch, and a test that
// wants a check to fire on a not-same fixture is expected to use it.
//
// The checks are ordered by how much they prove, and the order is also the
// reporting order: required-name coverage and top-level shape are the two that
// find real defects, leaf types find a handful, and a successful decode proves
// so little that it is reported last and described as such wherever it appears.
func CompareBody(f Fixture, symbol string, t reflect.Type, body []byte) Outcome {
	o := Outcome{Symbol: symbol, Fixture: f.ID, Shape: Describe(t)}
	if t == nil {
		// A method that returns only an error decodes no body, so no shape can
		// be compared. Saying so beats emitting four vacuous passes.
		o.Skipped = append(o.Skipped,
			SkippedCheck{CheckRequiredNames, "the SDK method decodes no response body"},
			SkippedCheck{CheckDeclaredNames, "the SDK method decodes no response body"},
			SkippedCheck{CheckShape, "the SDK method decodes no response body"},
			SkippedCheck{CheckLeafTypes, "the SDK method decodes no response body"},
			SkippedCheck{CheckDecodes, "the SDK method decodes no response body"})
		return o
	}
	o.GoType = typeName(t)

	var instance any
	if err := json.Unmarshal(body, &instance); err != nil {
		return unreadable(f, symbol, t,
			"the committed instance "+f.Path()+" is not parseable JSON, so no check could run",
			err.Error())
	}

	documented := documentedNames(f, instance)
	tags := tagsOf(t)

	// One pass per check, in the harness's priority order, so a report reads in
	// the same sequence the checks are documented in.
	//
	// Each pass is also handed the divergences the earlier ones recorded. The
	// decode check cannot localise a rejection by itself -- it only learns that
	// one happened -- so it is told which sharper finding already named the
	// disagreement, and restates that finding's own facts rather than quoting the
	// decoder's wording. See decodeCauses for why the wording is not quotable.
	for _, c := range []func(prior []Divergence) ([]Divergence, []SkippedCheck){
		func([]Divergence) ([]Divergence, []SkippedCheck) {
			return checkNames(symbol, f, o.Shape, tags, documented)
		},
		func([]Divergence) ([]Divergence, []SkippedCheck) {
			return checkDeclaredNames(symbol, f, o.Shape, tags)
		},
		func([]Divergence) ([]Divergence, []SkippedCheck) {
			return checkShape(symbol, f, o.Shape)
		},
		func([]Divergence) ([]Divergence, []SkippedCheck) {
			return checkLeaves(symbol, f, tags, documented)
		},
		func(prior []Divergence) ([]Divergence, []SkippedCheck) {
			return checkDecode(symbol, f, t, body,
				decodeCauses(f, o.Shape, tags, documented, prior))
		},
	} {
		divergences, skipped := c(o.Divergences)
		o.Divergences = append(o.Divergences, divergences...)
		o.Skipped = append(o.Skipped, skipped...)
	}
	sortDivergences(&o)
	return o
}

// unreadable is the one path through Compare that cannot compare anything: the
// body is missing or unparseable. It is reported as a decode failure because
// that is what a caller would see, and the other three checks are recorded as
// not applicable with the cause named, so the outcome never reads as a pass.
//
// why and cause are kept apart on purpose. why becomes a Divergence.Detail, and
// Detail is part of a baseline entry's identity, so it is written from the
// fixture's own identity and says nothing the decoder or the filesystem phrases.
// cause is the underlying message, which is diagnostic and deliberately not
// stable across operating systems or Go releases: an embed miss reads
// "no such file or directory" on one and "The system cannot find the file
// specified" on the other. It reaches SkippedCheck.Reason, which nothing
// compares, and never the identity.
func unreadable(f Fixture, symbol string, t reflect.Type, why, cause string) Outcome {
	o := Outcome{Symbol: symbol, Fixture: f.ID, Shape: Describe(t), GoType: typeName(t)}
	o.Divergences = append(o.Divergences, Divergence{
		Symbol: symbol, Fixture: f.ID, Kind: DecodeFailure, Detail: why,
	})
	for _, c := range []Check{CheckRequiredNames, CheckDeclaredNames, CheckShape, CheckLeafTypes} {
		o.Skipped = append(o.Skipped, SkippedCheck{c, cause})
	}
	sortDivergences(&o)
	return o
}

// sortDivergences puts an outcome's findings in check-priority order so two
// runs, and two people reading the same report, see the same order.
func sortDivergences(o *Outcome) {
	sort.SliceStable(o.Divergences, func(i, j int) bool {
		a, b := o.Divergences[i], o.Divergences[j]
		if ra, rb := checkOf(a.Kind), checkOf(b.Kind); ra != rb {
			return ra < rb
		}
		if a.Kind != b.Kind {
			return a.Kind < b.Kind
		}
		return a.Name < b.Name
	})
}

// documentedName is a required name paired with the value the committed fixture
// gives it, which is where the documented leaf type comes from. The manifest
// records the names; the fixture records their JSON kinds. Neither alone is
// enough for the leaf check.
type documentedName struct {
	Name  string
	Value any
}

// documentedNames reads every required name and its documented value out of the
// committed instance. Required names are recorded by the generator at exactly
// two levels -- the top-level object, or the single array element -- so those
// are the two places a value can be found.
func documentedNames(f Fixture, instance any) map[string]documentedName {
	out := map[string]documentedName{}
	if len(f.Checks.RequiredNames) == 0 {
		return out
	}
	var obj map[string]any
	switch v := instance.(type) {
	case map[string]any:
		obj = v
	case []any:
		if len(v) > 0 {
			obj, _ = v[0].(map[string]any)
		}
	}
	if obj == nil {
		return out
	}
	for _, name := range f.Checks.RequiredNames {
		if value, ok := obj[name]; ok {
			out[name] = documentedName{Name: name, Value: value}
		}
	}
	return out
}

// checkNames is check 1: every name the page marks required must reach a json
// tag, and the documented instance must actually carry it.
//
// The two halves are separate findings because they have separate causes and
// separate fixes. A missing tag is an SDK defect that decodes silently. A name
// the body omits is a fixture defect, which fixtures_test.go also catches, but
// asserting it here as well is what makes the body half of this check
// falsifiable: without it the check would be insensitive to the very bytes it
// claims to be reading, and a check insensitive to its input cannot be trusted
// when it reports nothing.
//
// A third thing is reported, not as a finding: a name that reached a tag only
// from below the level the page declares it at. That is not a defect, because the
// name still decodes, but it is why the count of missing names is a floor, and a
// floor that is not printed is indistinguishable from an exact total.
func checkNames(symbol string, f Fixture, shape WireShape, tags map[string]tagSite, documented map[string]documentedName) ([]Divergence, []SkippedCheck) {
	switch {
	case len(f.Checks.RequiredNames) == 0:
		return nil, []SkippedCheck{{CheckRequiredNames,
			"the page declares no required list at the top level or on the array " +
				"element, so " + string(CheckDeclaredNames) + " carries the weaker check instead"}}
	case shape.Carrier == nil:
		// A map decodes every documented name and declares none of them, so a
		// missing tag is not a defect it could have. Calling that a pass would
		// overstate what was checked, so it is reported as not applicable.
		return nil, []SkippedCheck{{CheckRequiredNames, fmt.Sprintf(
			"%s carries no field tags, so no required name can be missing from it",
			typeName(shape.Go))}}
	}

	var out []Divergence
	var coveredFromBelow []string
	for _, name := range f.Checks.RequiredNames {
		if _, ok := documented[name]; !ok {
			out = append(out, Divergence{
				Symbol: symbol, Fixture: f.ID, Kind: MissingRequiredName, Name: name,
				Detail: "the documented instance omits a name the page requires, " +
					"so there is nothing to check the SDK against",
			})
			continue
		}
		site, ok := lookupTag(tags, name)
		if !ok {
			out = append(out, Divergence{
				Symbol: symbol, Fixture: f.ID, Kind: MissingRequiredName, Name: name,
				Detail: fmt.Sprintf("no json tag in %s, so the documented value decodes to the zero value",
					typeName(shape.Carrier)),
			})
			continue
		}
		if site.Depth > 0 {
			coveredFromBelow = append(coveredFromBelow,
				fmt.Sprintf("%s at depth %d in %s", name, site.Depth, typeName(shape.Carrier)))
		}
	}
	if len(coveredFromBelow) > 0 {
		// The depth-scoped form of this check is not run, because the manifest
		// records only the two levels it read names from and not the depth of each
		// one. Saying so keeps the missing-name count honest about being a floor.
		return out, []SkippedCheck{{CheckNameDepth, fmt.Sprintf(
			"%d of %d required names reached a tag only from below the level the "+
				"page declares them at, so the missing-name count is a floor: %s",
			len(coveredFromBelow), len(f.Checks.RequiredNames),
			strings.Join(coveredFromBelow, ", "))}}
	}
	return out, nil
}

// checkDeclaredNames asks the required-name question of the names a page declares
// without promising, and runs only where that stronger check cannot.
//
// It is a separate check rather than a second arm of checkNames for two reasons.
// One is evidential: a required name is a promise and a declared name is a
// description, and an Outcome that conflated them could not say which it had
// examined a row against. The other is coverage: 90 of the 154 compared rows come
// from pages publishing no required list, so on those the strong check has nothing
// to look at and this is the only name check the harness applies at all.
//
// It is tag-only. checkNames also reports a required name missing from the
// committed instance, but the generator builds that instance from required names
// alone, so every declared name would be absent from it and the half would report
// 100% false positives. The depth accounting checkNames performs is not repeated
// either: tagsOf is called with the slice type on an array row, so the element's
// own fields land one level down and the Depth test would fire on nearly every
// declared name of every array row as a false "reached from below". The count this
// produces is a floor for the same reason the required count is, and that is
// stated in MissingDeclaredName rather than worked around.
func checkDeclaredNames(symbol string, f Fixture, shape WireShape, tags map[string]tagSite) ([]Divergence, []SkippedCheck) {
	if len(f.Checks.RequiredNames) > 0 {
		return nil, []SkippedCheck{{CheckDeclaredNames,
			"the page publishes a required list, so the stronger required-name check applies"}}
	}
	switch {
	case shape.Go == nil:
		return nil, []SkippedCheck{{CheckDeclaredNames, "the SDK method decodes no response body"}}
	case shape.Carrier == nil:
		// A map decodes every documented name and declares none of them, so a
		// missing tag is not a defect it could have, and reporting one would be
		// the harness's own error rather than a finding. Three data.* rows decode
		// into a free-form map and carry 105 declared names between them, so this
		// guard is worth more rows than the one carrying it.
		return nil, []SkippedCheck{{CheckDeclaredNames, fmt.Sprintf(
			"%s carries no field tags, so no declared name can be missing from it",
			typeName(shape.Go))}}
	}

	names := declaredNames(f)
	if len(names) == 0 {
		if f.Checks.DeclaredPropertyNameCount == 0 {
			// The page offers nothing to check with, so this row is a hole in the
			// evidence base rather than a pass. Filing it as one would make a page
			// nobody can examine read as a page that examined clean.
			return []Divergence{{
				Symbol: symbol, Fixture: f.ID, Kind: DeclaredInventoryEmpty,
				Detail: "the page publishes no required list and declares no property " +
					"name, so no name on it could be compared against the SDK at all",
			}}, nil
		}
		// The page does declare names, just not at the level this row is read at:
		// an array of scalars has no element object to name. That is a shape fact
		// for the shape check to carry, not a gap in the name inventory.
		return nil, []SkippedCheck{{CheckDeclaredNames, fmt.Sprintf(
			"the page declares %d property name(s) but none at the level this row reads",
			f.Checks.DeclaredPropertyNameCount)}}
	}

	var out []Divergence
	for _, name := range names {
		if _, ok := lookupTag(tags, name); ok {
			continue
		}
		out = append(out, Divergence{
			Symbol: symbol, Fixture: f.ID, Kind: MissingDeclaredName, Name: name,
			Detail: fmt.Sprintf("no json tag in %s, so the documented value decodes to the zero value",
				typeName(shape.Carrier)),
		})
	}
	return out, nil
}

// declaredNames returns the property names the page declares at the level the
// generator reads for this row: the single array element when the page's top
// level is an array, and the top-level object otherwise. It mirrors the choice
// the generator makes between the two requiredNamesSource values, so the weaker
// check looks where the stronger one would have looked.
func declaredNames(f Fixture) []string {
	if f.Checks.TopLevel == "array" {
		return f.Checks.DeclaredElementNames
	}
	return f.Checks.DeclaredTopLevelNames
}

// checkShape is check 3: the documented response's JSON kind, and an array's
// element kind, against what the SDK type accepts. Ten brokerfd rows and
// GetAgreementDetail turn on this one bit.
func checkShape(symbol string, f Fixture, shape WireShape) ([]Divergence, []SkippedCheck) {
	var out []Divergence
	if shape.TopLevel != "none" && f.Checks.TopLevel != "" && shape.TopLevel != f.Checks.TopLevel {
		out = append(out, Divergence{
			Symbol: symbol, Fixture: f.ID, Kind: TopLevelMismatch,
			Detail: fmt.Sprintf("documented top level is %s, %s decodes it as %s",
				f.Checks.TopLevel, typeName(shape.Go), shape.TopLevel),
		})
		// The element type is meaningless once the top levels disagree, and
		// reporting it too would double-count one defect.
		return out, []SkippedCheck{{CheckShape, "the top level already disagrees"}}
	}
	if shape.ElementKind == "unknown" && f.Checks.ElementType != "" {
		return out, []SkippedCheck{{CheckShape, fmt.Sprintf(
			"%s element is an interface, so its element kind cannot be compared",
			typeName(shape.Go))}}
	}
	if f.Checks.ElementType != "" && shape.ElementKind != "" &&
		shape.ElementKind != f.Checks.ElementType {
		out = append(out, Divergence{
			Symbol: symbol, Fixture: f.ID, Kind: ElementTypeMismatch,
			Detail: fmt.Sprintf("documented element type is %s, %s element is %s (%s)",
				f.Checks.ElementType, typeName(shape.Go), shape.ElementKind, typeName(shape.Element)),
		})
	}
	return out, nil
}

// checkLeaves is check 4: for every required name the SDK does tag, does the Go
// field hold the JSON kind the documented value has.
//
// Only names the SDK tags are considered, because a name it does not tag is
// already reported by check 1 and has no field to type-check.
func checkLeaves(symbol string, f Fixture, tags map[string]tagSite, documented map[string]documentedName) ([]Divergence, []SkippedCheck) {
	if len(f.Checks.RequiredNames) == 0 {
		return nil, []SkippedCheck{{CheckLeafTypes,
			"no documented required name to type-check"}}
	}

	var out []Divergence
	var notJudgeable []string
	for _, name := range f.Checks.RequiredNames {
		doc, ok := documented[name]
		if !ok {
			// fixtures_test.go already asserts every required name is present
			// in the instance, so this cannot happen on a passing tree.
			notJudgeable = append(notJudgeable, name)
			continue
		}
		want := jsonKind(doc.Value)
		if want == "object" || want == "array" || want == "null" {
			notJudgeable = append(notJudgeable, name)
			continue
		}
		site, ok := lookupTag(tags, name)
		if !ok {
			continue // already reported by check 1
		}
		ok, reason := leafVerdict(site.Type, want)
		if ok {
			continue
		}
		if reason != "" {
			// A field that decides for itself is not a mismatch; it is a field
			// this check cannot judge, and saying so beats guessing.
			notJudgeable = append(notJudgeable, name)
			continue
		}
		out = append(out, Divergence{
			Symbol: symbol, Fixture: f.ID, Kind: LeafTypeMismatch, Name: name,
			Detail: fmt.Sprintf("documented %s, SDK field %s.%s is %s",
				want, typeName(site.Field.Type), site.Field.Name, jsonKindOf(site.Type)),
		})
	}
	if len(notJudgeable) > 0 {
		return out, []SkippedCheck{{CheckLeafTypes, fmt.Sprintf(
			"not a documented scalar, or the field carries its own decoder, for: %s",
			strings.Join(notJudgeable, ", "))}}
	}
	if len(documented) == 0 {
		return nil, []SkippedCheck{{CheckLeafTypes, "no documented required name to type-check"}}
	}
	return out, nil
}

// moneyType is the type the SDK uses for a decimal financial value, and
// jsonUnmarshalerType is any type that brings its own decoder. Both are
// resolved by name against this package's imports so the leaf check can treat
// them as facts about the type rather than as guesses about its name.
var (
	moneyType           = reflect.TypeOf(money.Money{})
	jsonUnmarshalerType = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
)

// leafSamples maps a documented JSON kind to a literal of that kind, used to ask a
// type that brings its own decoder whether it can actually be built from the shape
// the page documents.
//
// The values are deliberately small and unremarkable: the question is which JSON
// kind the type accepts, not whether it accepts this particular magnitude.
var leafSamples = map[string]string{
	"string":  `"1"`,
	"number":  `1`,
	"boolean": `true`,
	"object":  `{}`,
	"array":   `[]`,
}

// leafVerdict decides whether a Go field can hold a documented JSON kind.
//
// The three outcomes are deliberate. ok means the field can hold it. A non-empty
// reason means the field is not judgeable from outside -- it decodes itself, and its
// decoder rejected the sample -- which is a different fact from disagreement. An
// empty reason with ok false is a mismatch.
//
// A type that brings its own decoder is asked rather than assumed. The earlier
// version named money.Money as the one exempt type, which is a fact about the SDK
// that goes stale the moment a second such type appears, and it had to be extended
// the moment data.QuoteTime did. Sending the documented kind through the type's own
// UnmarshalJSON answers the question directly, and it is deliberately one-sided: a
// type that accepts the sample becomes positively confirmed, and a type that
// rejects it stays not judgeable rather than being reported as a mismatch. That way
// this can only reduce the set of unexamined fields, never invent a finding about a
// decoder the harness has not understood.
func leafVerdict(field reflect.Type, want string) (ok bool, reason string) {
	t := field
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t == moneyType {
		// money.Money is the SDK's decimal carrier, and its UnmarshalJSON
		// (pkg/domain/money/money.go:177) tries a JSON string and then a JSON
		// number, because Webull sends the same field both ways across endpoints.
		// So it agrees with a documented string and a documented number, and with
		// nothing else: it cannot be built from a boolean or an array, and saying
		// otherwise would be an exemption rather than a statement about the type.
		return want == "string" || want == "number", ""
	}
	if reflect.PointerTo(t).Implements(jsonUnmarshalerType) || t.Implements(jsonUnmarshalerType) {
		if sample, has := leafSamples[want]; has {
			ptr := reflect.New(t)
			if err := json.Unmarshal([]byte(sample), ptr.Interface()); err == nil {
				return true, ""
			}
			return false, "the SDK field type supplies its own json.Unmarshaler, and that " +
				"decoder rejected a " + want + " of the documented kind"
		}
		return false, "the SDK field type supplies its own json.Unmarshaler"
	}
	return jsonKindOf(t) == want, ""
}

// checkDecode is check 5: the documented instance must unmarshal into the SDK
// type. It is the weakest check and is described as such in DivergenceKind: a
// successful decode is consistent with a type that ignores every name, so a
// pass here carries little information. A failure, on the other hand, is
// normally the same defect a sharper check already named, and both are kept
// because the failure is what a caller actually sees.
//
// causes is that sharper diagnosis, built by decodeCauses and not by this
// function: a decode that fails says only that a failure happened, so the detail
// this records is a restatement of what the harness already knows rather than of
// what the decoder printed.
func checkDecode(symbol string, f Fixture, t reflect.Type, body []byte, causes []string) ([]Divergence, []SkippedCheck) {
	ptr := reflect.New(t)
	if err := json.Unmarshal(body, ptr.Interface()); err != nil {
		return []Divergence{{
			Symbol: symbol, Fixture: f.ID, Kind: DecodeFailure,
			Detail: fmt.Sprintf("the documented instance does not unmarshal into %s: %s",
				typeName(t), strings.Join(causes, "; ")),
		}}, nil
	}
	return nil, nil
}

// decodeCauses names, from the harness's own inputs, why the documented instance
// does not unmarshal into the compared type, so a decode-failure Detail is a
// function of the SDK and the page rather than of the decoder's wording.
//
// The distinction is load-bearing, because Detail is part of a baseline entry's
// identity (see Divergence). encoding/json's message is not a stable string. It
// was quoted verbatim into 29 decode-failure details at v2.1.15, and the next Go
// release reworded
//
//	json: cannot unmarshal string into Go value of type brokerfd.AccountForm
//
// to
//
//	json: cannot unmarshal string into .0 of type brokerfd.AccountForm
//
// which is what turned the two Go-stable matrix jobs red while the same tree
// passed on the pinned toolchain. Each of the 29 entries then read as a new
// divergence *and* as a baseline entry that had stopped reproducing. A string
// that moves with the toolchain cannot be part of an identity meant to outlive
// one, and re-recording the details against whatever wording the current host
// produces would only move the breakage to the next release.
//
// What the stdlib message carries -- what was sent, and what was expected to take
// it -- is available here without the decoder: the fixture's recorded top-level
// kind and element type, and the compared type's shape by reflection. Each cause
// below is therefore the same fact the corresponding sharper finding already
// states, phrased the same way, so a reader meets one account of the defect in
// two entries rather than two accounts of it.
//
// The clauses come out in the order the sharper checks ran, which is deterministic
// without sorting: the name and leaf checks walk f.Checks.RequiredNames, which is
// sorted and asserted to be sorted, and at most one shape finding exists per
// outcome. They are built from each finding's Kind and Name, never from its
// Detail, so rewording a sharper check does not silently reword this one.
func decodeCauses(f Fixture, shape WireShape, tags map[string]tagSite, documented map[string]documentedName, prior []Divergence) []string {
	var out []string
	for _, d := range prior {
		switch d.Kind {
		case TopLevelMismatch:
			out = append(out, fmt.Sprintf("the documented top level is %s and %s decodes it as %s",
				f.Checks.TopLevel, typeName(shape.Go), shape.TopLevel))
		case ElementTypeMismatch:
			out = append(out, fmt.Sprintf("the documented element type is %s and the %s element is %s",
				f.Checks.ElementType, typeName(shape.Go), shape.ElementKind))
		case LeafTypeMismatch:
			doc, ok := documented[d.Name]
			if !ok {
				continue
			}
			site, ok := lookupTag(tags, d.Name)
			if !ok {
				continue
			}
			out = append(out, fmt.Sprintf("%s: documented %s, SDK field %s.%s is %s",
				d.Name, jsonKind(doc.Value), typeName(site.Field.Type), site.Field.Name,
				jsonKindOf(site.Type)))
		}
	}
	if len(out) == 0 {
		// Reachable in principle and unreachable on the committed tree: a rejection
		// none of the three sharper checks names, such as a type disagreement on a
		// name the page does not require. It is stated rather than left blank,
		// because a Detail that says nothing is the case this function exists to
		// prevent, and a blanket clause would read as a claim the harness checked
		// for a cause and did not find.
		out = append(out, "no name, shape or leaf check disagrees, so the rejection is at "+
			"a name or a type that none of the three checks compares")
	}
	return out
}
