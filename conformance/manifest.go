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
	"sort"
)

// testdataDir is the committed fixture tree, relative to this package. It is a
// const rather than a path walked at run time so a test cannot silently pass
// against a different tree.
const testdataDir = "testdata"

// PathMatch reports how the path an SDK method sends compares with the path its
// documentation page declares. It is a three-state type rather than a bool
// because "the SDK sends a different path" and "this endpoint is deliberately
// mapped to no SDK symbol" are different facts, and a bool renders them
// identically.
type PathMatch string

const (
	// PathMatchSame means the SDK method sends the documented path.
	PathMatchSame PathMatch = "same"
	// PathMatchDiffers means the SDK method sends some other path. Each such
	// row is a path defect candidate; see IMPLEMENTATION_STATUS.md.
	PathMatchDiffers PathMatch = "differs"
	// PathMatchNone means the manifest entry names no SDK symbol, so there is
	// no SDK path to compare. The docgen AREAS manifest spells this as an em
	// dash, which survives into Fixture.SDKSymbol verbatim.
	PathMatchNone PathMatch = "noSdkPath"
)

// Source locates a fixture in the documentation it was derived from. Every
// field exists so a reviewer can re-derive the fixture by hand: open the page,
// take the JSON block at this index, and read the schema at this pointer.
type Source struct {
	// URL is the page's canonical .md address on developer.webull.hk or
	// developer.webull.com.
	URL string `json:"url"`
	// CacheFile is the filename the page has inside the docgen cache, which is
	// the only way to read a page without network access.
	CacheFile string `json:"cacheFile"`
	// JSONBlockIndex is the zero-based index of the ```json block holding the
	// OpenAPI document. Zero for every page in the tree; a guide or gRPC page
	// with several blocks would select a later one, which is why it is recorded
	// rather than assumed.
	JSONBlockIndex int `json:"jsonBlockIndex"`
	// SchemaPointer is a JSON Pointer to the 200 application/json schema within
	// that block.
	SchemaPointer string `json:"schemaPointer"`
}

// Documented is the request the page itself documents.
type Documented struct {
	// Method is the HTTP verb the page declares, upper-cased.
	Method string `json:"method"`
	// Path is the request path the page declares, exactly as written there.
	Path string `json:"path"`
}

// Checks records which comparisons a fixture actually supports. The fields exist
// because an empty fixture and a fixture full of names are both legal output of
// the generator, and only the manifest distinguishes them.
type Checks struct {
	// TopLevel is "object" or "array", matching the documented 200 schema.
	TopLevel string `json:"topLevel"`
	// ElementType is the type of an array's element, or "" when TopLevel is
	// "object".
	ElementType string `json:"elementType"`
	// DeclaresRequired reports whether the page declares a required list at the
	// top level or on the array element. False does not mean the names were
	// verified; it means the page named none, so the minimal instance has
	// nothing from Webull to check. See the package documentation.
	DeclaresRequired bool `json:"declaresRequired"`
	// RequiredNamesSource is "topLevel", "items", or "" when DeclaresRequired is
	// false. It says which of the two levels RequiredNames came from.
	RequiredNamesSource string `json:"requiredNamesSource"`
	// RequiredNames are the names the page marks required, sorted. Every one of
	// them is present in the fixture, which is what makes the fixture minimal
	// rather than arbitrary.
	RequiredNames []string `json:"requiredNames"`
	// DeclaredTopLevelNames are all the names the top-level object declares,
	// including optional ones the minimal instance omits.
	DeclaredTopLevelNames []string `json:"declaredTopLevelNames"`
	// DeclaredElementNames are all the names an array's element object declares,
	// including optional ones.
	DeclaredElementNames []string `json:"declaredElementNames"`
	// DeclaredPropertyNameCount counts every name the schema declares at any
	// depth, deduplicated.
	DeclaredPropertyNameCount int `json:"declaredPropertyNameCount"`
	// PropertyNameCountInFixture counts the distinct names actually present in
	// the committed instance, at any depth.
	PropertyNameCountInFixture int `json:"propertyNameCountInFixture"`
	// SynthesizedLeafCount counts emitted leaves that carried no inline example
	// and so took a stand-in value derived from the declared type. Any non-zero
	// value means part of the fixture is the generator's invention, not
	// Webull's.
	SynthesizedLeafCount int `json:"synthesizedLeafCount"`
	// RequiredWithoutExample names required properties the page declares but
	// gives no schema for. Empty across the committed tree.
	RequiredWithoutExample []string `json:"requiredWithoutExample"`
}

// Choice is a recorded schema decision: which oneOf branch was taken, or which
// synthetic key stands in for a free-form object's declared value type.
type Choice struct {
	// Kind is "oneOf" or "additionalProperties".
	Kind string `json:"kind"`
	// Pointer locates the decision inside the documented schema.
	Pointer string `json:"pointer"`
	// Emitted reports whether the minimal instance actually reached this node.
	// False means the node sits under an optional property, which the minimal
	// rule omits, so the decision is recorded but not exercised.
	Emitted bool `json:"emitted"`

	// BranchIndex is the zero-based oneOf branch taken.
	BranchIndex int `json:"branchIndex,omitempty"`
	// Title is that branch's schema title, when it has one.
	Title *string `json:"title,omitempty"`
	// Selection is "title-matches-discriminator" when a sibling type example
	// named the branch, or "first-branch" when no branch was named and the
	// first was taken.
	Selection string `json:"selection,omitempty"`
	// BranchCount is how many branches the oneOf declared.
	BranchCount int `json:"branchCount,omitempty"`
	// BranchRequiredNameCount is how many required names the chosen branch
	// declares. Zero means the choice cannot be observed in the instance,
	// because the minimal rule emits only required names.
	BranchRequiredNameCount int `json:"branchRequiredNameCount,omitempty"`
	// Discriminator is the sibling property that selects the branch, or "".
	Discriminator string `json:"discriminator,omitempty"`
	// DiscriminatorValue is that property's documented example, or "".
	DiscriminatorValue any `json:"discriminatorValue,omitempty"`

	// ValueType is the declared type of a typed additionalProperties.
	ValueType string `json:"valueType,omitempty"`
	// SyntheticKey is the key used to give a free-form object a shape.
	SyntheticKey string `json:"syntheticKey,omitempty"`
}

// Fixture is one documented endpoint: where its fixture came from, which SDK
// symbol documents it, and what the committed instance supports.
type Fixture struct {
	// ID is the fixture's identity and the first half of its filename.
	ID string `json:"id"`
	// Fixture is the instance's path under testdata.
	Fixture string `json:"fixture"`
	// Area is the docgen AREAS manifest key this endpoint belongs to.
	Area string `json:"area"`
	// PageTitle is the endpoint's name in the docgen manifest.
	PageTitle string `json:"pageTitle"`
	// SDKSymbol is the manifest's SDK symbol, which may list alternatives
	// separated by " / " and may be the docgen em-dash placeholder when the
	// endpoint is deliberately mapped to no symbol.
	SDKSymbol string `json:"sdkSymbol"`
	// Source locates the documented schema.
	Source Source `json:"source"`
	// Documented is the request the page declares.
	Documented Documented `json:"documented"`
	// SDKPath is the path the SDK method actually sends, or "" when none.
	SDKPath string `json:"sdkPath"`
	// SDKPathConst is the Go path constant SDKPath came from, or "".
	SDKPathConst string `json:"sdkPathConst"`
	// SDKPathMatch compares SDKPath with Documented.Path.
	SDKPathMatch PathMatch `json:"sdkPathMatch"`
	// Checks records which comparisons the instance supports.
	Checks Checks `json:"checks"`
	// OneOfChoices is every oneOf branch decision on the page.
	OneOfChoices []Choice `json:"oneOfChoices"`
	// AdditionalProperties is every typed additionalProperties on the page.
	AdditionalProperties []Choice `json:"additionalProperties"`
	// Bytes is the committed instance's length in bytes.
	Bytes int `json:"bytes"`
}

// Path returns the fixture's location inside the embedded tree, using the
// forward-slash form the manifest writes.
func (f Fixture) Path() string {
	return testdataDir + "/" + f.Fixture
}

// Read returns the committed fixture bytes.
//
// The bytes come from the embedded copy, so they are what was committed at build
// time and cannot be regenerated by anything a caller does. That is the property
// that lets a reviewer trust them.
func (f Fixture) Read() ([]byte, error) {
	data, err := fixturesFS.ReadFile(f.Path())
	if err != nil {
		return nil, fmt.Errorf("conformance: read fixture %s: %w", f.ID, err)
	}
	return data, nil
}

// Decode reads the committed fixture and unmarshals it into a generic Go value:
// map[string]any for an object top level, []any for an array. Numbers become
// float64, so a caller comparing against a typed DTO should compare via
// json.Number or re-marshal, not by float equality.
func (f Fixture) Decode() (any, error) {
	data, err := f.Read()
	if err != nil {
		return nil, err
	}
	var out any
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("conformance: decode fixture %s: %w", f.ID, err)
	}
	return out, nil
}

// Tripwire is the fixture size policy the generator enforces.
type Tripwire struct {
	// LogAboveBytes is the size above which a fixture is reported.
	LogAboveBytes int `json:"logAboveBytes"`
	// FailAboveBytes is the size above which generation fails.
	FailAboveBytes int `json:"failAboveBytes"`
	// Reason states why the ceiling exists.
	Reason string `json:"reason"`
}

// Totals summarises the committed tree.
type Totals struct {
	// Fixtures counts the instances.
	Fixtures int `json:"fixtures"`
	// FixtureBytes is their combined length.
	FixtureBytes int `json:"fixtureBytes"`
	// PagesWithoutRequired counts fixtures from pages that declared no
	// required list at the top level or on the array element.
	PagesWithoutRequired int `json:"pagesWithoutRequired"`
	// PagesWithOneOf counts pages carrying a oneOf.
	PagesWithOneOf int `json:"pagesWithOneOf"`
	// PagesWithAdditionalProperties counts pages carrying a typed
	// additionalProperties.
	PagesWithAdditionalProperties int `json:"pagesWithAdditionalProperties"`
	// OversizedFixtures names any instance over FailAboveBytes. Empty in a
	// committed tree, since generation fails before one can be written.
	OversizedFixtures []string `json:"oversizedFixtures"`
}

// Generator records how the tree was produced.
type Generator struct {
	Tool                string `json:"tool"`
	Command             string `json:"command"`
	Source              string `json:"source"`
	DerivedFrom         string `json:"derivedFrom"`
	NeverReads          string `json:"neverReads"`
	MinimalInstanceRule string `json:"minimalInstanceRule"`
}

// SkippedPage is a documented endpoint that produced no fixture, with the
// reason. These are the gRPC subscription references and guide pages that carry
// no 200 application/json schema.
type SkippedPage struct {
	Area   string `json:"area"`
	Label  string `json:"label"`
	Page   string `json:"page"`
	Reason string `json:"reason"`
}

// Manifest is the committed provenance record for the fixture tree. It is
// regenerated only by tools/conformance/gen_fixtures.py and verified by
// make conformance-fixtures; a test run never rewrites it.
type Manifest struct {
	Generator                 Generator         `json:"generator"`
	SizeTripwire              Tripwire          `json:"sizeTripwire"`
	Totals                    Totals            `json:"totals"`
	PagesWithoutRequiredNames []string          `json:"pagesWithoutRequiredNames"`
	Notes                     map[string]string `json:"notes"`
	SkippedPages              []SkippedPage     `json:"skippedPages"`
	Fixtures                  []Fixture         `json:"fixtures"`
}

// ByID returns the fixture with the given ID.
func (m *Manifest) ByID(id string) (Fixture, bool) {
	for _, f := range m.Fixtures {
		if f.ID == id {
			return f, true
		}
	}
	return Fixture{}, false
}

// Load reads the manifest from the embedded tree.
func Load() (*Manifest, error) {
	data, err := fixturesFS.ReadFile(testdataDir + "/manifest.json")
	if err != nil {
		return nil, fmt.Errorf("conformance: read manifest: %w", err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("conformance: parse manifest: %w", err)
	}
	// The manifest is emitted with sorted-by-construction fixture order; sorting
	// again here keeps Load independent of that and makes report output stable.
	sort.Slice(m.Fixtures, func(i, j int) bool { return m.Fixtures[i].ID < m.Fixtures[j].ID })
	return &m, nil
}

// Names returns every distinct object key in the committed fixture, at any
// depth. These are the names a decoder will see on the wire, which is the set a
// DTO's json tags have to cover.
func (f Fixture) Names() ([]string, error) {
	decoded, err := f.Decode()
	if err != nil {
		return nil, err
	}
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
	walk(decoded)
	names := make([]string, 0, len(seen))
	for k := range seen {
		names = append(names, k)
	}
	sort.Strings(names)
	return names, nil
}
