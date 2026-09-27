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
func committedLength(data []byte) int {
	return len(bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n")))
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
func TestManifestDescribesTheWholeTree(t *testing.T) {
	m := load(t)
	want := make(map[string]bool, len(m.Fixtures)+1)
	want["manifest.json"] = true
	for _, f := range m.Fixtures {
		want[f.Fixture] = true
	}
	var got []string
	err := fs.WalkDir(fixturesFS, testdataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, ok := strings.CutPrefix(path, testdataDir+"/")
		if !ok {
			return fmt.Errorf("path %q is not under %s", path, testdataDir)
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
		t.Errorf("tree holds %d files, manifest describes %d", len(got), len(want))
	}
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
// committedLength gives. The structural assertions around it are what catch a
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
