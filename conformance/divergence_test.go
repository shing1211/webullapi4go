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
	"testing"

	"github.com/shing1211/webullapi4go/data"
	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

// These tests compare the SDK's own decoded structs against the documented wire
// shape in the committed fixtures. They read the fixture bytes from the embedded
// copy, so nothing a run does can regenerate the input it is checking.
//
// There is a baseline, and the reason for it is in baseline.go: the current tree
// carries about a hundred documented divergences, so a gate that reported them
// would be red on day one and a gate that skipped them would have thrown the
// findings away. Instead the divergences are recorded, the gate fails only on a
// *change* in the set, and every entry carries a reason.

// loadManifest reads the committed manifest for a comparison test.
func loadManifest(t *testing.T) *Manifest {
	t.Helper()
	m, err := Load()
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}
	return m
}

// TestSymbolTableCoversTheManifest holds the symbol table to the manifest in
// both directions.
//
// The forward direction is the one that matters. Every manifest symbol in the
// root module must have a type; a symbol without one is skipped by the
// comparison, and a skipped symbol is indistinguishable from a passing one --
// the exact failure mode this instrument exists to remove. So the unmapped list
// is asserted empty, by name, rather than merely computed and reported.
//
// The reverse direction catches a stale row. A table entry for a symbol the
// manifest no longer names is either a typo that hides a real unmapped symbol or
// a leftover to delete, and both are invisible without this direction.
func TestSymbolTableCoversTheManifest(t *testing.T) {
	m := loadManifest(t)
	u := UnmappedReport(m)

	if len(u.Symbols) != 0 {
		for _, sym := range u.Symbols {
			t.Errorf("manifest names %s but SDKTypes has no entry, so no comparison "+
				"runs for it", sym)
		}
	}

	// Subject normalises a bare alternative against its row's package, so the
	// set of symbols the comparison will use is built the same way here.
	named := map[string]bool{}
	for _, f := range m.Fixtures {
		subject, alts, ok := Subject(f.SDKSymbol)
		if !ok {
			continue
		}
		named[subject] = true
		for _, a := range alts {
			named[a] = true
		}
	}
	for _, sym := range TableSymbols() {
		if !named[sym] {
			t.Errorf("SDKTypes has an entry for %s, which the manifest does not name; "+
				"a stale row hides a real unmapped symbol", sym)
		}
	}

	t.Logf("manifest fixtures: %d", len(m.Fixtures))
	t.Logf("compared against a Go type: %d", comparedCount(t, m))
	t.Logf("unmapped (a table hole, and the gate fails on one): %d %v", len(u.Symbols), u.Symbols)
	t.Logf("out of scope (%s is a separate Go module, so NOT covered here): %d %s",
		strings.TrimSuffix(outOfScopePackage, "."), len(u.OutOfScope), strings.Join(u.OutOfScope, ", "))
	t.Logf("symbol table entries: %d", len(SDKTypes))

	// Every manifest row must land in exactly one of the two buckets, or a row
	// is being passed over and nothing would say so.
	outcomes, notCompared, err := CompareAll(m)
	if err != nil {
		t.Fatalf("CompareAll: %v", err)
	}
	if len(outcomes)+len(notCompared) != len(m.Fixtures) {
		t.Errorf("CompareAll accounted for %d of %d fixtures: %d compared, %d not "+
			"compared", len(outcomes)+len(notCompared), len(m.Fixtures),
			len(outcomes), len(notCompared))
	}
	byReason := map[string]int{}
	for _, n := range notCompared {
		byReason[n.Reason]++
		if n.Reason == ReasonUntableable {
			t.Errorf("fixture %s names %s, which SDKTypes has no entry for; a "+
				"symbol with no type is a hole in the table, not coverage",
				n.Fixture, n.Symbol)
		}
	}
	for _, reason := range sortedKeys(byReason) {
		t.Logf("not compared: %3d x %s", byReason[reason], reason)
	}

	// A row reported not comparable is a coverage statement, so it is counted
	// here rather than left to be inferred from an absence of findings.
	notComparable := 0
	for _, o := range outcomes {
		if o.NotComparable == nil {
			continue
		}
		notComparable++
		if len(o.Divergences) != 0 {
			t.Errorf("%s: reported not comparable but carries %d divergences; the "+
				"two states are exclusive", o.Fixture, len(o.Divergences))
		}
		t.Logf("not comparable: %s (%s)", o.Fixture, o.NotComparable.Reason)
	}
	t.Logf("compared and NOT comparable (the page is not this call's contract): %d", notComparable)
}

// TestSymbolTableResolvesItsDecodeTargets holds every entry to the table's own
// rule, and is where an unexported envelope is read from the SDK sources.
//
// Two directions, because they fail differently. Forward: every entry resolves to
// a type, which is what makes an envelope the source cannot represent a loud
// table defect instead of a row that quietly compares against nil. Reverse: every
// entry that names an unexported envelope is checked field by field against the
// declaration the resolver read, so a renamed field, a changed json tag or a
// moved declaration all turn the gate red.
//
// The reverse direction is the reason a mirror struct was not used instead. A
// hand-written mirror of data.corpActionResponse would have been a second copy of
// the same fact, and nothing in the build would have compared the two; this way
// there is only one copy.
func TestSymbolTableResolvesItsDecodeTargets(t *testing.T) {
	for _, sym := range TableSymbols() {
		entry := SDKTypes[sym]
		target, err := entry.DecodeTarget()
		if err != nil {
			t.Errorf("%s: %v", sym, err)
			continue
		}
		if target == nil && entry.Envelope == "" && entry.Note == "" {
			t.Errorf("%s: a nil decode target with no note saying the method "+
				"decodes no body is indistinguishable from a hole", sym)
		}
		if entry.Envelope == "" {
			continue
		}
		assertEnvelopeMatchesSource(t, sym, entry.Envelope, target)
	}
}

// assertEnvelopeMatchesSource re-reads an envelope's declaration and compares it,
// field by field, with the type the resolver built.
func assertEnvelopeMatchesSource(t *testing.T, sym, spelling string, built reflect.Type) {
	t.Helper()
	pkg, name, ok := strings.Cut(spelling, ".")
	if !ok {
		t.Fatalf("%s: envelope %q is not spelled package.TypeName", sym, spelling)
	}
	dir, err := packageDir(pkg)
	if err != nil {
		t.Fatalf("%s: %v", sym, err)
	}
	fields, declPos, err := parseEnvelopeFields(dir, pkg, name)
	if err != nil {
		t.Fatalf("%s: %v", sym, err)
	}
	t.Logf("%s: %s declared at %s, %d field(s)", sym, spelling, declPos, len(fields))
	if got := built.NumField(); got != len(fields) {
		t.Fatalf("%s: the rebuilt type has %d fields, the declaration has %d",
			sym, got, len(fields))
	}
	for i, f := range fields {
		got := built.Field(i)
		if got.Name != f.Name {
			t.Errorf("%s: field %d is %s, the declaration says %s", sym, i, got.Name, f.Name)
		}
		if string(got.Tag) != f.Tag {
			t.Errorf("%s: field %s carries tag %q, the declaration says %q",
				sym, f.Name, got.Tag, f.Tag)
		}
	}
}

// TestEnvelopeTagsAreUsableAsStructTags holds the rebuilt envelope tags to the one
// property assertEnvelopeMatchesSource cannot check.
//
// That helper compares the rebuilt tag against the tag the parser read, so it is
// faithful to its input whether or not that input is usable. Both sides were wrong
// together: go/ast reports a tag as its source literal, so a raw string arrived
// wrapped in backticks, reflect.StructTag failed every lookup on it, and wireNameOf
// fell back to the Go field name. That agrees for camelCase and silently loses a
// snake_case name, which is why the round-trip stayed green while a documented
// snake_case name compared as absent from an SDK that carries it.
//
// So this asserts usability rather than fidelity: every tag the parser reports must
// yield a readable json key, and at least one of them must be snake_case, since that
// is the case the fallback could not see.
func TestEnvelopeTagsAreUsableAsStructTags(t *testing.T) {
	var tags, snake int
	for _, sym := range TableSymbols() {
		entry := SDKTypes[sym]
		if entry.Envelope == "" {
			continue
		}
		pkg, name, ok := strings.Cut(entry.Envelope, ".")
		if !ok {
			t.Fatalf("%s: envelope %q is not spelled package.TypeName", sym, entry.Envelope)
		}
		dir, err := packageDir(pkg)
		if err != nil {
			t.Fatalf("%s: %v", sym, err)
		}
		fields, _, err := parseEnvelopeFields(dir, pkg, name)
		if err != nil {
			t.Fatalf("%s: %v", sym, err)
		}
		for _, f := range fields {
			if f.Tag == "" {
				continue
			}
			tags++
			jsonName, ok := reflect.StructTag(f.Tag).Lookup("json")
			if !ok {
				t.Errorf("%s: %s.%s field %s at %s carries tag %q, from which no json "+
					"key can be read, so a name lookup would fall back to the Go field "+
					"name and report a carried name as absent",
					sym, pkg, name, f.Name, f.Pos, f.Tag)
				continue
			}
			if strings.Contains(jsonName, "_") {
				snake++
			}
		}
	}
	// Both counters guard against a pass that checked nothing. An empty tag set
	// means the envelopes stopped declaring one, and a snake_case count of zero
	// means the tags stopped being snake_case, either of which would leave this
	// test green while the thing it exists to protect was gone.
	if tags == 0 {
		t.Error("no envelope field reported a tag, so nothing was checked")
	}
	if snake == 0 {
		t.Error("no envelope tag resolved a snake_case json name, so the case a " +
			"Go-field-name fallback cannot represent is no longer covered")
	}
	t.Logf("checked %d envelope tag(s), %d resolving a snake_case json name", tags, snake)
}

// TestKnownDivergenceBaseline is the gate.
//
// It passes only when the observed divergence set is exactly the committed set.
// A new divergence fails; a baseline entry that no longer reproduces also fails,
// because a fixed defect and a check that changed meaning are indistinguishable
// from here and both need a person to read the tree.
func TestKnownDivergenceBaseline(t *testing.T) {
	m := loadManifest(t)
	b, err := LoadBaseline()
	if err != nil {
		t.Fatalf("%v", err)
	}

	observed, notCompared, err := CompareAll(m)
	if err != nil {
		t.Fatalf("CompareAll: %v", err)
	}

	got := map[string]Divergence{}
	for _, o := range observed {
		for _, d := range o.Divergences {
			if _, dup := got[d.Key()]; dup {
				t.Errorf("divergence reported twice: %s", d.Key())
			}
			got[d.Key()] = d
		}
	}
	want := b.ByKey()

	var newKeys, staleKeys []string
	for key, d := range got {
		if _, ok := want[key]; !ok {
			newKeys = append(newKeys, fmt.Sprintf("%s\n    %s %s %s\n      %s",
				key, d.Symbol, d.Kind, d.Name, d.Detail))
		}
	}
	for key, e := range want {
		if _, ok := got[key]; !ok {
			staleKeys = append(staleKeys, fmt.Sprintf("%s\n    %s %s %s\n      %s\n      reason on file: %s",
				key, e.Symbol, e.Kind, e.Name, e.Detail, e.Reason))
		}
	}
	sort.Strings(newKeys)
	sort.Strings(staleKeys)

	for _, s := range newKeys {
		t.Errorf("NEW DIVERGENCE, not in the baseline: %s", s)
	}
	for _, s := range staleKeys {
		t.Errorf("BASELINE ENTRY NO LONGER REPRODUCES: %s\n      Either the defect was "+
			"fixed, or a check changed what it means. Both need a person to look. Run "+
			"TestObservedDivergenceReport to see what the harness observes now.", s)
	}
	if len(newKeys) == 0 && len(staleKeys) == 0 {
		t.Logf("observed %d divergences, matching all %d baseline entries exactly",
			len(got), len(want))
	}

	assertCoverage(t, b, m, observed, notCompared)
}

// assertCoverage pins how much the comparison looked at, so the gate cannot go
// green by checking less than it did when the baseline was taken. A manifest row
// added without a table entry, or a table entry removed, moves these numbers and
// fails here even when no divergence changed.
func assertCoverage(t *testing.T, b *Baseline, m *Manifest, observed []Outcome, notCompared []NotCompared) {
	t.Helper()
	u := UnmappedReport(m)
	noSymbol := 0
	for _, n := range notCompared {
		if n.Reason == ReasonNoSymbol {
			noSymbol++
		}
	}
	notComparable := 0
	for _, o := range observed {
		if o.NotComparable != nil {
			notComparable++
		}
	}
	if b.Coverage.Fixtures != len(m.Fixtures) {
		t.Errorf("baseline was taken over %d fixtures, the manifest has %d",
			b.Coverage.Fixtures, len(m.Fixtures))
	}
	if b.Coverage.Compared != len(observed) {
		t.Errorf("baseline compared %d manifest rows, this run compared %d",
			b.Coverage.Compared, len(observed))
	}
	if b.Coverage.NoSymbolRows != noSymbol {
		t.Errorf("baseline recorded %d em-dash rows, this run found %d",
			b.Coverage.NoSymbolRows, noSymbol)
	}
	if b.Coverage.NotComparable != notComparable {
		t.Errorf("baseline recorded %d rows as not comparable, this run found %d; "+
			"the set of endpoints whose page is not the SDK call's contract moved",
			b.Coverage.NotComparable, notComparable)
	}
	if b.Coverage.SymbolsInTable != len(SDKTypes) {
		t.Errorf("baseline was taken with a %d-entry symbol table, this build has %d; "+
			"the baseline is stale", b.Coverage.SymbolsInTable, len(SDKTypes))
	}
	if strings.Join(b.Coverage.OutOfScope, ",") != strings.Join(u.OutOfScope, ",") {
		t.Errorf("baseline recorded out-of-scope symbols %v, this run found %v",
			b.Coverage.OutOfScope, u.OutOfScope)
	}
	if len(b.Coverage.Unmapped) != 0 {
		t.Errorf("baseline claims %d unmapped symbols; the gate requires none, "+
			"because an unmapped symbol is a comparison that never ran",
			len(b.Coverage.Unmapped))
	}
}

// TestObservedDivergenceReport prints the full observed set and the coverage
// numbers, grouped by check. It never fails, and it is how the baseline is read
// back: `go test ./conformance/ -run TestObservedDivergenceReport -v` shows
// exactly what the harness sees today, so a baseline edit can be made against a
// fresh observation rather than against a memory of one.
func TestObservedDivergenceReport(t *testing.T) {
	m := loadManifest(t)
	observed, notCompared, err := CompareAll(m)
	if err != nil {
		t.Fatalf("CompareAll: %v", err)
	}
	u := UnmappedReport(m)

	byReason := map[string]int{}
	for _, n := range notCompared {
		byReason[n.Reason]++
	}
	t.Logf("manifest fixtures: %d   rows compared: %d   rows not compared: %d",
		len(m.Fixtures), len(observed), len(notCompared))
	t.Logf("unmapped symbols (would be a table defect): %d %v", len(u.Symbols), u.Symbols)
	t.Logf("out-of-scope symbols (broker/ is a separate module): %d", len(u.OutOfScope))
	t.Logf("symbol table entries: %d", len(SDKTypes))
	for _, reason := range sortedKeys(byReason) {
		t.Logf("not compared: %3d x %s", byReason[reason], reason)
	}

	// Rows reported not comparable, with both paths. This is the third state,
	// and it is the one most easily read as a pass, so it is printed in full.
	var incomparable []Outcome
	for _, o := range observed {
		if o.NotComparable != nil {
			incomparable = append(incomparable, o)
		}
	}
	t.Logf("compared and NOT comparable: %d of %d compared rows", len(incomparable), len(observed))
	for _, o := range incomparable {
		t.Logf("  %-58s %s", o.Fixture, o.Symbol)
		t.Logf("      %s", o.NotComparable.Reason)
		t.Logf("      page documents %s; the method sends %s",
			o.NotComparable.DocumentedPath, o.NotComparable.SDKPath)
	}

	grouped := map[Check]map[string][]Divergence{}
	for _, o := range observed {
		for _, d := range o.Divergences {
			if grouped[checkForKind(d.Kind)] == nil {
				grouped[checkForKind(d.Kind)] = map[string][]Divergence{}
			}
			g := grouped[checkForKind(d.Kind)]
			g[d.Symbol] = append(g[d.Symbol], d)
		}
	}
	total := 0
	for _, check := range checks {
		g := grouped[check]
		if len(g) == 0 {
			t.Logf("%s: none", check)
			continue
		}
		n := 0
		for _, ds := range g {
			n += len(ds)
		}
		total += n
		t.Logf("%s: %d divergences across %d symbols", check, n, len(g))
		for _, sym := range sortedKeys(g) {
			ds := g[sym]
			t.Logf("  %s (%d)", sym, len(ds))
			for _, d := range ds {
				t.Logf("    %s %s", d.Kind, d.Detail)
			}
		}
	}
	t.Logf("TOTAL observed divergences: %d", total)

	// The checks that could not run are reported too. A check that silently
	// stops applying is how a whole endpoint class goes unchecked without
	// anybody deciding to stop checking it.
	skipped := map[string]int{}
	for _, o := range observed {
		for _, s := range o.Skipped {
			reason := s.Reason
			if len(reason) > 72 {
				reason = reason[:72] + "..."
			}
			skipped[fmt.Sprintf("%s: %s", s.Check, reason)]++
		}
	}
	if len(skipped) > 0 {
		t.Logf("checks that did not apply, by reason:")
		for _, reason := range sortedKeys(skipped) {
			t.Logf("  %3d x %s", skipped[reason], reason)
		}
	}
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func comparedCount(t *testing.T, m *Manifest) int {
	t.Helper()
	observed, _, err := CompareAll(m)
	if err != nil {
		t.Fatalf("CompareAll: %v", err)
	}
	return len(observed)
}

// TestComparisonBites proves the comparison can fail, by changing one side of
// each check and showing the check's verdict moves.
//
// Without this, a green gate is unfalsifiable: a check that never fires passes
// exactly as well as a correct one. Every case below changes either a *copy* of
// the SDK type, built with reflect.StructOf, or a *copy* of the documented body
// supplied through CompareBody, or a *copy* of a manifest row. No production file
// and no committed fixture is edited, so TestFixtureMatchesItsRecord would fail
// if one were.
func TestComparisonBites(t *testing.T) {
	m := loadManifest(t)

	t.Run("RequiredNameCoverage", func(t *testing.T) {
		const sym = "brokerfd.GetFDPositions"
		f := mustFixture(t, m, "broker-fd-us/GET-broker-assets-positions-list")
		entry := mustType(t, sym)

		// The three names this row used to be missing are now carried. The page
		// requires cost_price, last_price and unrealized_profit_loss, and
		// FDPosition used to tag only average_cost, market_value and unrealized_pl,
		// so the documented values decoded to zero money.Money with no error
		// reported. FDPosition now carries both sets, so neither the documented
		// names nor the SDK's own are missing, and the row records no divergence
		// for the first time. Asserting the fix rather than the defect is what
		// keeps a retag from silently reintroducing it, and the three mutation
		// checks below still prove the check itself bites.
		before := Compare(f, sym, entry)
		for _, name := range []string{"cost_price", "last_price", "unrealized_profit_loss"} {
			if hasKind(before, MissingRequiredName, name) {
				t.Errorf("the documented name %s is reported missing again, so the fix "+
					"did not hold: %s", name, summary(before))
			}
		}
		if len(before.Divergences) != 0 {
			t.Errorf("expected this row to record no divergence now that both name sets are "+
				"carried, but it has %d: %s", len(before.Divergences), summary(before))
		}
		t.Logf("unmodified fixture: %d divergences", len(before.Divergences))
		for _, d := range before.Divergences {
			t.Logf("    %s %s %s", d.Kind, d.Name, d.Detail)
		}

		// Tag side: drop a json tag the SDK does have, from a copy of the type.
		// The check must then name it, which is what proves the check reads tags
		// rather than agreeing with whatever the page says.
		untagged := withField(t, entry, "symbol", func(f reflect.StructField) reflect.StructField {
			f.Tag = `json:"-"`
			return f
		})
		afterTag := Compare(f, sym, untagged)
		if !hasKind(afterTag, MissingRequiredName, "symbol") {
			t.Fatalf("removing the symbol tag did not make the name check fail: %s",
				summary(afterTag))
		}
		if countKind(afterTag, MissingRequiredName) != countKind(before, MissingRequiredName)+1 {
			t.Errorf("expected exactly one more missing name: before %d, after %d",
				countKind(before, MissingRequiredName), countKind(afterTag, MissingRequiredName))
		}
		t.Logf("copy of the type with the symbol tag removed: %s", summary(afterTag))

		// Body side: delete cost_price from a copy of the documented instance.
		// The name is still required by the page, and it is now carried by
		// nothing, so the check has to say so.
		edited := deleteName(t, f, "cost_price")
		afterBody := CompareBody(f, sym, entry, edited)
		if !hasKind(afterBody, MissingRequiredName, "cost_price") {
			t.Fatalf("deleting cost_price from the body did not make the name check "+
				"fail: %s", summary(afterBody))
		}
		if d := detailFor(afterBody, MissingRequiredName, "cost_price"); !strings.Contains(d, "documented instance omits") {
			t.Errorf("expected the body-omission half of the check, got %q", d)
		}
		t.Logf(`copy of the fixture with "cost_price" deleted: %s`, summary(afterBody))

		// Case side: a field with no json tag encodes under its Go name, and
		// encoding/json folds case when no tag matches exactly, so a page that
		// spells the name in a different case still decodes into it. Stripping the
		// tag must therefore NOT turn a covered name into a reported missing one.
		folded := withField(t, entry, "symbol", func(sf reflect.StructField) reflect.StructField {
			sf.Tag = ``
			return sf
		})
		afterCase := Compare(f, sym, folded)
		if hasKind(afterCase, MissingRequiredName, "symbol") {
			t.Errorf("a copy whose symbol tag was stripped reports `symbol` missing, "+
				"but encoding/json falls back to the Go field name and folds case, "+
				"so the decoder still fills it: %s", summary(afterCase))
		}
		if got, want := countKind(afterCase, MissingRequiredName), countKind(before, MissingRequiredName); got != want {
			t.Errorf("removing a tag that only changes the case of a Go name changed "+
				"the missing-name count: before %d, after %d", want, got)
		}
		t.Logf("copy of the type with the symbol tag stripped (Go name only): %s",
			summary(afterCase))
	})

	t.Run("DeclaredNameCoverage", func(t *testing.T) {
		// GetFDCorporateActions is one of the 34 rows the declared-name check
		// fires on: the page declares 14 properties and FDCorporateAction carries
		// 7, so 12 are missing and the two declared names it does carry are
		// ex_date and record_date. This row is on a page with no required list, so
		// the weaker check
		// is the only name check that applies to it.
		const sym = "brokerfd.GetFDCorporateActions"
		f := mustFixture(t, m, "broker-fd-us/GET-broker-instruments-stocks-corporate-actions-get")
		entry := mustType(t, sym)
		before := Compare(f, sym, entry)
		for _, name := range []string{"event_id", "payment_date", "final_pay_date", "event_type"} {
			if !hasKind(before, MissingDeclaredName, name) {
				t.Fatalf("the recorded defect on %s does not reproduce: %s",
					name, summary(before))
			}
		}
		// The row must not also report the required half: it is the stronger check
		// that is not applicable here, and a page that publishes a required list
		// is what the weaker check defers to.
		if countKind(before, MissingRequiredName) != 0 {
			t.Errorf("a page with no required list reported %d required-name rows: %s",
				countKind(before, MissingRequiredName), summary(before))
		}
		t.Logf("unmodified fixture: %d divergences", len(before.Divergences))

		// Tag side: drop a json tag the SDK does have, from a copy of the type.
		// The check must then name it, which is what proves it reads tags rather
		// than agreeing with whatever the page declares.
		untagged := withField(t, entry, "ex_date", func(f reflect.StructField) reflect.StructField {
			f.Tag = `json:"-"`
			return f
		})
		afterTag := Compare(f, sym, untagged)
		if !hasKind(afterTag, MissingDeclaredName, "ex_date") {
			t.Fatalf("removing the ex_date tag did not make the declared-name check "+
				"fail, so the check is not reading tags: %s", summary(afterTag))
		}
		if countKind(afterTag, MissingDeclaredName) != countKind(before, MissingDeclaredName)+1 {
			t.Errorf("expected exactly one more declared name: before %d, after %d",
				countKind(before, MissingDeclaredName), countKind(afterTag, MissingDeclaredName))
		}
		t.Logf("copy of the type with the ex_date tag removed: %s", summary(afterTag))

		// Case side: encoding/json folds case on a tag miss, so a page spelling a
		// carried name differently still decodes into it. Stripping the tag must
		// therefore NOT report a covered name as missing.
		// Name side: strip the tag entirely, so the field encodes under its Go name
		// ExDate, and the check must still report ex_date missing. This is not the
		// case-fold rescue the required-name subtest relies on, and the difference
		// is worth pinning: strings.EqualFold folds case but not punctuation, so
		// ex_date cannot match ExDate, and neither can encoding/json. Verified
		// against the decoder: an untagged struct is filled from `symbol` and
		// `success` but left empty by `ex_date`. Every declared name on this page
		// is snake_case, so on this row the fallback can never apply, and a check
		// that assumed it would have passed a type the decoder drops the value for.
		folded := withField(t, entry, "ex_date", func(sf reflect.StructField) reflect.StructField {
			sf.Tag = ``
			return sf
		})
		afterCase := Compare(f, sym, folded)
		if !hasKind(afterCase, MissingDeclaredName, "ex_date") {
			t.Errorf("a copy whose ex_date tag was stripped does not report it missing, "+
				"but the Go name ExDate cannot be filled from the wire key ex_date, so "+
				"the value would be dropped: %s", summary(afterCase))
		}
		if got, want := countKind(afterCase, MissingDeclaredName), countKind(before, MissingDeclaredName)+1; got != want {
			t.Errorf("expected exactly one more declared name: before %d, after %d",
				want, got)
		}
		t.Logf("copy of the type with the ex_date tag stripped (Go name only): %s",
			summary(afterCase))

		// Guard side: a free-form map decodes every documented name and declares
		// none of them, so a missing tag on one is the harness's own error. The
		// three data.* rows below declare 105 names between them, so failing to
		// guard this would have produced 105 false positives.
		const mapSym = "data.GetBalanceSheet"
		mapFixture := mustFixture(t, m, "fundamentals/GET-market-data-fundamentals-balance-sheets-get")
		onMap := Compare(mapFixture, mapSym, mustType(t, mapSym))
		if countKind(onMap, MissingDeclaredName) != 0 {
			t.Errorf("a free-form map reported %d declared names missing, which is "+
				"the harness's error rather than a finding: %s",
				countKind(onMap, MissingDeclaredName), summary(onMap))
		}
		if !skippedCheck(onMap, CheckDeclaredNames) {
			t.Errorf("the map row does not report the declared-name check as skipped: %s",
				summary(onMap))
		}
		t.Logf("free-form map row, %d declared names, reported as skipped: %s",
			mustFixture(t, m, "fundamentals/GET-market-data-fundamentals-balance-sheets-get").
				Checks.DeclaredPropertyNameCount, summary(onMap))

		// Inventory side: a page declaring no property at all is a hole in the
		// evidence base, and has to be reported as one rather than as a pass.
		const emptySym = "brokerfd.ListAccountForms"
		emptyFixture := mustFixture(t, m, "broker-fd-us/GET-broker-forms-list")
		onEmpty := Compare(emptyFixture, emptySym, mustType(t, emptySym))
		if !hasKind(onEmpty, DeclaredInventoryEmpty) {
			t.Fatalf("a page declaring no property name is not reported as an empty "+
				"inventory: %s", summary(onEmpty))
		}
		if countKind(onEmpty, MissingDeclaredName) != 0 {
			t.Errorf("an empty inventory also reported %d missing names",
				countKind(onEmpty, MissingDeclaredName))
		}
		t.Logf("page with an empty declared inventory: %s", summary(onEmpty))

		// Deference side: a page that does publish a required list is the
		// stronger check's business, and this one must stand aside entirely.
		required := Compare(mustFixture(t, m, "broker-fd-us/GET-broker-assets-positions-list"),
			"brokerfd.GetFDPositions", mustType(t, "brokerfd.GetFDPositions"))
		if countKind(required, MissingDeclaredName) != 0 {
			t.Errorf("a page publishing a required list also reported %d declared-name rows",
				countKind(required, MissingDeclaredName))
		}
		if !skippedCheck(required, CheckDeclaredNames) {
			t.Errorf("a page publishing a required list does not report the " +
				"declared-name check as skipped")
		}
	})

	t.Run("TopLevelShape", func(t *testing.T) {
		// GetFDTransferFees decodes a slice where the page documents an object.
		//
		// This subtest used to bite on brokerfd.GetFDAssetsDetail, which inverted
		// the same way, but v2.1.33 gave that method the documented envelope and the
		// defect it was constructed to detect no longer exists. Asserting it would
		// mean asserting a fixed defect, so the bite moved here: the class still has
		// live rows, so a real one is a better fixture than a manufactured copy.
		// The same direction applies on both ends, so the helpers below are unchanged.
		const sym = "brokerfd.GetFDTransferFees"
		f := mustFixture(t, m, "broker-fd-us/GET-broker-fees-get")
		entry := mustType(t, sym)
		before := Compare(f, sym, entry)
		if !hasKind(before, TopLevelMismatch) {
			t.Fatalf("the recorded shape inversion does not reproduce: %s", summary(before))
		}
		t.Logf("unmodified fixture: %d divergences", len(before.Divergences))
		for _, d := range before.Divergences {
			t.Logf("    %s %s %s", d.Kind, d.Name, d.Detail)
		}

		// A documented object must fail to decode into the SDK's slice. That is
		// what a caller sees, and it is the decode check firing on the same bit
		// the shape check named.
		if err := decodeCommitted(f, entry); err == nil {
			t.Fatal("the documented object decoded into []brokerfd.TransferFee, so " +
				"the recorded shape inversion is not real")
		} else {
			t.Logf("documented object into the SDK slice: %v", err)
		}

		// SDK side: rebuild the slice as the struct the page documents, from the
		// slice element's own fields. The top-level bit must flip to agreeing and
		// the object must then decode. Names the element lacks stay missing --
		// that is a separate check reporting a separate defect, and pretending
		// otherwise would hide it.
		asStruct := structOfElem(t, entry)
		after := Compare(f, sym, asStruct)
		if hasKind(after, TopLevelMismatch) {
			t.Errorf("a struct-typed copy still reports a top-level mismatch: %s", summary(after))
		}
		if err := decodeCommitted(f, asStruct); err != nil {
			t.Errorf("a struct-typed copy still fails to decode: %v", err)
		}
		t.Logf("copy of the type as a struct: %s", summary(after))

		// Documented side: a copy of the fixture whose recorded contract is an
		// array, with an array body. The check must pass on both ends.
		asArrayFixture, asArray := reDescribe(t, f, "array", "object", []any{map[string]any{}})
		if hasKind(CompareBody(asArrayFixture, sym, entry, asArray), TopLevelMismatch) {
			t.Error("an array-described copy still reports a top-level mismatch")
		}
		t.Logf("copy of the fixture described as an array: no top-level mismatch")
	})

	t.Run("LeafType", func(t *testing.T) {
		// The class this subtest proves is currently empty: its only two rows were
		// data.Quote.QuoteTime, and QuoteTime is now a type that reads both the
		// quoted-string and bare-number forms Webull's four depth pages publish
		// between them. So the check is no longer exercised by a live defect and
		// the bite has to be manufactured, or this subtest would go on asserting a
		// defect that no longer exists.
		//
		// The manufactured type below is a plain int64, which is exactly the field
		// the SDK carried before, against the same page that documents a string.
		// That is the question the check was written to answer.
		const sym = "data.GetQuotes"
		f := mustFixture(t, m, "market-data-stock/GET-market-data-stocks-depths-list")
		entry := mustType(t, sym)

		// The fix holds: the real type reads the page's quoted string, so the row
		// reports no leaf disagreement.
		before := Compare(f, sym, entry)
		if hasKind(before, LeafTypeMismatch, "quote_time") {
			t.Errorf("the real QuoteTime type still reports a leaf mismatch: %s", summary(before))
		}
		if err := decodeCommitted(f, entry); err != nil {
			t.Errorf("the real QuoteTime type still fails to decode the committed "+
				"fixture, whose quote_time is a quoted string: %v", err)
		}
		t.Logf("unmodified fixture against the real type: %s", summary(before))

		// The bite: re-type quote_time to int64, which is the pre-fix field, and the
		// check must name the disagreement again.
		retyped := withField(t, entry, "quote_time", func(sf reflect.StructField) reflect.StructField {
			sf.Type = reflect.TypeOf(int64(0))
			return sf
		})
		afterType := Compare(f, sym, retyped)
		if !hasKind(afterType, LeafTypeMismatch, "quote_time") {
			t.Fatalf("an int64 quote_time against a documented string no longer "+
				"reports a leaf mismatch, so this check is not proven to bite: %s",
				summary(afterType))
		}
		if err := decodeCommitted(f, retyped); err == nil {
			t.Error("an int64 quote_time still decodes the page's quoted string, " +
				"so the check is not detecting the real disagreement")
		}
		t.Logf("copy of the type with quote_time as an int64: %s", summary(afterType))

		// Documented side: the same body presented as a number, which is what the
		// futures and event-contract pages show. The leaf check must go quiet, which
		// proves it reads the body rather than the page's prose.
		asNumber := setLeaf(t, f, "quote_time", json.Number("1640688000000"))
		afterBody := CompareBody(f, sym, entry, asNumber)
		if hasKind(afterBody, LeafTypeMismatch, "quote_time") {
			t.Errorf("a numeric quote_time still reports a leaf mismatch: %s", summary(afterBody))
		}
		t.Logf(`copy of the fixture with quote_time as a number: %s`, summary(afterBody))

		// The two types that accept more than one documented shape must be exempt
		// for exactly the shapes they accept, or the exemption is a blanket one and
		// the comment claiming otherwise is false.
		for _, tc := range []struct {
			name string
			typ  reflect.Type
		}{
			{"money.Money", reflect.TypeOf(money.Money{})},
			{"data.QuoteTime", reflect.TypeOf(data.QuoteTime(0))},
		} {
			for _, kind := range []string{"string", "number"} {
				if ok, reason := leafVerdict(tc.typ, kind); !ok {
					t.Errorf("%s is rejected for a documented %s: %s", tc.name, kind, reason)
				}
			}
			if ok, _ := leafVerdict(tc.typ, "boolean"); ok {
				t.Errorf("%s is accepted for a documented boolean, which its "+
					"UnmarshalJSON cannot build", tc.name)
			}
		}
		t.Log("money.Money and data.QuoteTime each agree with a documented string " +
			"and number, and not with a boolean")
	})

	t.Run("NotComparableWhenThePathDiffers", func(t *testing.T) {
		// The rule under test: a page whose path the SDK method does not send is
		// not that call's contract, so a schema difference is not a finding about
		// the SDK. Both halves matter. The row must claim nothing, and it must say
		// so distinctly rather than read as a pass.
		const sym = "data.GetStockInstruments"
		f := mustFixture(t, m, "trading/GET-trading-instruments-stocks-profiles-list")
		if f.SDKPathMatch == PathMatchSame {
			t.Fatalf("%s now reports sdkPathMatch same; the case under test needs a "+
				"row that differs", f.ID)
		}
		entry := mustType(t, sym)

		out := Compare(f, sym, entry)
		if out.NotComparable == nil {
			t.Fatalf("a differing-path row was compared as if the page were the "+
				"contract: %s", summary(out))
		}
		if len(out.Divergences) != 0 {
			t.Errorf("a not-comparable row carries %d divergences: %s",
				len(out.Divergences), summary(out))
		}
		if out.NotComparable.Reason != ReasonPathNotContract {
			t.Errorf("reason = %q, want %q", out.NotComparable.Reason, ReasonPathNotContract)
		}
		if out.NotComparable.DocumentedPath != f.Documented.Path ||
			out.NotComparable.SDKPath != f.SDKPath {
			t.Errorf("the row does not carry both paths: page %q, SDK %q",
				out.NotComparable.DocumentedPath, out.NotComparable.SDKPath)
		}
		if len(out.Skipped) != len(checks) {
			t.Errorf("%d checks skipped, want all %d: a not-comparable row must not "+
				"read as a partial pass", len(out.Skipped), len(checks))
		}
		t.Logf("%s: %s", f.ID, out.NotComparable.Reason)
		t.Logf("  page documents %s; the method sends %s",
			out.NotComparable.DocumentedPath, out.NotComparable.SDKPath)
		t.Logf("  %d divergences claimed, %d checks recorded as not applicable",
			len(out.Divergences), len(out.Skipped))

		// Converse: the same type against a copy of the row whose sdkPathMatch says
		// same. The comparison must now run, and the shape check must speak, which
		// is the pair of facts that keep the rule from being a blanket suppression.
		asSame := f
		asSame.SDKPathMatch = PathMatchSame
		asSame.SDKPath = asSame.Documented.Path
		compared := Compare(asSame, sym, entry)
		if compared.NotComparable != nil {
			t.Fatalf("a same-path row is still reported not comparable: %s",
				compared.NotComparable.Reason)
		}
		if len(compared.Divergences) == 0 {
			t.Error("a same-path row reports no divergences, so the rule is not what " +
				"decides the outcome; it cannot be proven to bite")
		}
		t.Logf("copy of the row with sdkPathMatch same: %s", summary(compared))
	})

	t.Run("EmbeddedFieldsAreInlined", func(t *testing.T) {
		// No production DTO in brokerfd, data or trade uses an anonymous field, so
		// the shapes under test are declared here. They are declared as real types
		// rather than built with reflect.StructOf because one of the two cases needs
		// an *unexported* embedded field, which the reflector will not build.
		//
		// Both cases are ways of getting the tag set wrong, and both are silent:
		// recording an embedded type's Go name as a wire name that does not exist,
		// and skipping an unexported embedded struct whose exported fields
		// encoding/json promotes and does set. A wrong tag set produces a missing
		// name the SDK decodes, or hides one it does not.
		tags := tagsOf(reflect.TypeOf(embedProbe{}))
		if _, ok := tags["embedBase"]; ok {
			t.Error("an unexported untagged embedded field is recorded under its Go " +
				"type name; encoding/json inlines it, so no such key exists on the wire")
		}
		if _, ok := tags["promoted"]; !ok {
			t.Error("the exported field of an unexported embedded struct is missing " +
				"from the tag set, so a name the SDK decodes into it would be reported " +
				"as missing")
		}
		if _, ok := tags["own"]; !ok {
			t.Error("the enclosing struct's own tag is missing")
		}
		if got := tags["promoted"].Depth; got != 0 {
			t.Errorf("the promoted field sits at depth %d, want 0: an inlined field is "+
				"flattened into the enclosing object, so its names are carried at the "+
				"enclosing depth", got)
		}
		t.Logf("unexported untagged embed: names %v", sortedTagNames(tags))

		// The same for an exported embedded struct, which was the other half of the
		// bug: its Go type name is not a wire name either.
		exportedTags := tagsOf(reflect.TypeOf(exportedEmbedProbe{}))
		if _, ok := exportedTags["ExportedEmbedBase"]; ok {
			t.Error("an exported untagged embedded field is recorded under its Go " +
				"type name; encoding/json inlines it")
		}
		if _, ok := exportedTags["promoted"]; !ok {
			t.Error("the exported field of an exported embedded struct is missing")
		}
		t.Logf("exported untagged embed: names %v", sortedTagNames(exportedTags))

		// An embedded field named in its tag is a nested object, so the name is
		// real and must be kept. Dropping it would report a name the SDK decodes.
		tagged := tagsOf(reflect.TypeOf(taggedEmbedProbe{}))
		if _, ok := tagged["inner"]; !ok {
			t.Error("an embedded field named in its tag is a nested object, so its " +
				"name is real, but it is missing from the tag set")
		}
		t.Logf("tagged embed: names %v", sortedTagNames(tagged))

		// A field tagged "-" is not encoded at all.
		if _, ok := tagsOf(reflect.TypeOf(skippedFieldProbe{}))["Hidden"]; ok {
			t.Error(`a field tagged json:"-" is in the tag set`)
		}
	})

	t.Run("BaselineRefusesAnUnexplainedEntry", func(t *testing.T) {
		// The mechanism, not just the data. A baseline that accepted a
		// reasonless entry would make a suppression indistinguishable from a
		// record, which is the failure mode the whole design exists to remove.
		b := &Baseline{
			Policy: "test",
			Entries: []KnownDivergence{{
				Divergence: Divergence{Symbol: "s", Fixture: "f", Kind: MissingRequiredName,
					Name: "n", Detail: "d"},
				RecordedIn: "IMPLEMENTATION_STATUS.md",
			}},
		}
		err := b.Validate()
		if err == nil {
			t.Fatal("a baseline entry with no reason validated; an unexplained entry " +
				"is a suppressed finding and must be refused")
		}
		t.Logf("reasonless entry refused: %v", err)

		b.Entries[0].Reason = "   "
		if err := b.Validate(); err == nil {
			t.Fatal("a whitespace-only reason validated")
		}
		b.Entries[0].Reason = "because"
		b.Entries[0].RecordedIn = ""
		if err := b.Validate(); err == nil {
			t.Fatal("an entry naming no recording document validated")
		}
		t.Logf("the committed baseline loads and validates with every reason present")
		if committed, err := LoadBaseline(); err == nil {
			t.Logf("  %d entries, all with a reason and a recording document", len(committed.Entries))
		}
	})
}

// TestDecodeFailureDetailsSurviveTheDecoder proves the decode-failure Detail is
// written from the check's own inputs, so a Go release that rewords
// encoding/json's message cannot turn every recorded decode failure into a new
// divergence and an entry that stopped reproducing at the same time.
//
// That is not hypothetical: v2.1.15 shipped 29 details quoted from the decoder
// verbatim, and the next Go release rephrasing "cannot unmarshal string into Go
// value of type brokerfd.AccountForm" to "cannot unmarshal string into .0 of
// type brokerfd.AccountForm" turned the two Go-stable matrix jobs red, on a tree
// that passed on the pinned toolchain.
//
// Two directions, because a detail can be wrong in two ways. First, no committed
// detail may contain a fragment of a decoder message, and every one must name the
// type it failed on: a detail that quotes the decoder is one release away from
// breaking the gate, and a detail that names nothing is indistinguishable from an
// emptied one. Second, and empirically: for every row whose diagnosis is a shape,
// the same comparison is re-run against a body that also fails to decode while
// every check's verdict is held fixed -- a top-level JSON number, which the
// diagnosis never reads -- and the Detail must come out byte-identical even though
// the decoder's message is demonstrably different. The two leaf-only rows have no
// such body: for them the decoder's message already names the leaf the harness
// names, so changing the wording means inventing a second real disagreement. They
// are held to the first direction plus an exact restatement of the leaf verdicts.
func TestDecodeFailureDetailsSurviveTheDecoder(t *testing.T) {
	m := loadManifest(t)
	observed, _, err := CompareAll(m)
	if err != nil {
		t.Fatalf("CompareAll: %v", err)
	}
	committed, err := LoadBaseline()
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Fragments of an encoding/json message, across the wordings seen so far.
	fragments := []string{
		"json:", "cannot unmarshal", "Go value of type", "Go struct field",
		"of type", "invalid character", "unexpected end of JSON",
		"into .0 of type", "into Go value",
	}

	var recorded, bySecondBody, byInspection int
	for _, e := range committed.Entries {
		if e.Kind != DecodeFailure {
			continue
		}
		recorded++
		for _, frag := range fragments {
			if strings.Contains(e.Detail, frag) {
				t.Errorf("%s: the recorded detail quotes the decoder (%q): %q",
					e.Symbol, frag, e.Detail)
			}
		}
		want := "the documented instance does not unmarshal into " +
			typeName(mustType(t, e.Symbol)) + ": "
		if !strings.HasPrefix(e.Detail, want) {
			t.Errorf("%s: the recorded detail does not name the type it failed on, "+
				"so it is either quoting the decoder or saying nothing; want a %q "+
				"prefix, got %q", e.Symbol, want, e.Detail)
		}
	}
	t.Logf("recorded decode-failure entries: %d, none quoting the decoder", recorded)

	for _, o := range observed {
		detail, ok := detailOf(o, DecodeFailure)
		if !ok {
			continue
		}
		target, err := SDKTypes[o.Symbol].DecodeTarget()
		if err != nil {
			t.Fatalf("%s: %v", o.Symbol, err)
		}
		f, ok := m.ByID(o.Fixture)
		if !ok {
			t.Fatalf("no fixture %q", o.Fixture)
		}
		body, err := f.Read()
		if err != nil {
			t.Fatal(err)
		}
		first := decodeErrorText(body, target)

		if !hasKind(o, TopLevelMismatch) && !hasKind(o, ElementTypeMismatch) {
			// Leaf-only: the detail must be exactly the leaf verdicts, and the
			// decoder's message must contribute nothing to it.
			for _, d := range o.Divergences {
				if d.Kind != LeafTypeMismatch {
					continue
				}
				if clause := d.Name + ": " + d.Detail; !strings.Contains(detail, clause) {
					t.Errorf("%s: the detail does not carry the leaf verdict %q: %q",
						o.Symbol, clause, detail)
				}
			}
			for _, frag := range fragments {
				if strings.Contains(detail, frag) {
					t.Errorf("%s: the detail quotes the decoder (%q): %q", o.Symbol, frag, detail)
				}
			}
			byInspection++
			t.Logf("%s: leaf-only, decoder said %q, detail is %q", o.Symbol, first, detail)
			continue
		}

		// A top-level JSON number is rejected by every SDK response type, and the
		// diagnosis reads the fixture's recorded top-level kind and the type's
		// shape, neither of which the body carries.
		second := decodeErrorText([]byte("42"), target)
		if first == second {
			t.Errorf("%s: both bodies draw the same message from the decoder, so "+
				"this row proves nothing either way: %q", o.Symbol, first)
			continue
		}
		alt := CompareBody(f, o.Symbol, target, []byte("42"))
		got, ok := detailOf(alt, DecodeFailure)
		if !ok {
			t.Errorf("%s: a top-level number decoded, so the row cannot be "+
				"re-proved here", o.Symbol)
			continue
		}
		if got != detail {
			t.Errorf("%s: the detail moved when only the decoder's wording could "+
				"have changed.\n  committed body: %q\n  number body:    %q",
				o.Symbol, detail, got)
			continue
		}
		bySecondBody++
		t.Logf("%s: decoder said %q then %q, detail stayed %q", o.Symbol, first, second, detail)
	}
	t.Logf("rows re-proved with a second body: %d, rows proved by inspection: %d",
		bySecondBody, byInspection)
	if bySecondBody+byInspection == 0 {
		t.Log("no decode failure is observed on this tree, so there is nothing to re-prove")
	}
}

// detailOf returns the Detail of the outcome's first divergence of that kind.
func detailOf(o Outcome, kind DivergenceKind) (string, bool) {
	for _, d := range o.Divergences {
		if d.Kind == kind {
			return d.Detail, true
		}
	}
	return "", false
}

// decodeErrorText is what a caller would see from unmarshalling body into t.
func decodeErrorText(body []byte, t reflect.Type) string {
	err := json.Unmarshal(body, reflect.New(t).Interface())
	if err == nil {
		return "<no error>"
	}
	return err.Error()
}

func mustFixture(t *testing.T, m *Manifest, id string) Fixture {
	t.Helper()
	f, ok := m.ByID(id)
	if !ok {
		t.Fatalf("no fixture %q", id)
	}
	return f
}

// mustType is the symbol table's entry for sym, resolved to the type its method
// decodes. It goes through DecodeTarget rather than reading Type, so a test
// cannot accidentally compare against a projection for a symbol whose decode
// target is an unexported envelope.
func mustType(t *testing.T, sym string) reflect.Type {
	t.Helper()
	entry, found := SDKTypes[sym]
	if !found {
		t.Fatalf("%s is not in the symbol table", sym)
	}
	target, err := entry.DecodeTarget()
	if err != nil {
		t.Fatalf("%s: %v", sym, err)
	}
	return target
}

// The four shapes TestComparisonBites/EmbeddedFieldsAreInlined declares, because
// reflect.StructOf cannot build an unexported field and the case under test is
// precisely an unexported one. Two base types are needed because "embedded" and
// "unexported" are independent: an embed of an unexported type name and an embed of
// an exported one both inline, and the two are the two halves of the bug.
type embedBase struct {
	Promoted string `json:"promoted"`
}

type ExportedEmbedBase struct {
	Promoted string `json:"promoted"`
}

type embedProbe struct {
	embedBase
	Own string `json:"own"`
}

type exportedEmbedProbe struct {
	ExportedEmbedBase
	Own string `json:"own"`
}

type taggedEmbedProbe struct {
	ExportedEmbedBase `json:"inner"`
	Own               string `json:"own"`
}

type skippedFieldProbe struct {
	Hidden string `json:"-"`
}

// mustJSON reads the committed bytes of a fixture.
func mustJSON(t *testing.T, f Fixture) []byte {
	t.Helper()
	raw, err := f.Read()
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

// deleteName removes a name from the documented instance and returns the
// rendered body. The fixture's recorded Checks are deliberately left alone: the
// page still requires the name, and the case under test is what the comparison
// says about a name the instance no longer carries.
func deleteName(t *testing.T, f Fixture, name string) []byte {
	t.Helper()
	var instance any
	if err := json.Unmarshal(mustJSON(t, f), &instance); err != nil {
		t.Fatal(err)
	}
	obj := topObject(t, instance)
	if _, ok := obj[name]; !ok {
		t.Fatalf("%s carries no %q at the compared level", f.ID, name)
	}
	delete(obj, name)
	return mustMarshal(t, instance)
}

// setLeaf replaces one documented value at the compared level, which is how a
// documented leaf's JSON kind is changed without touching the SDK.
func setLeaf(t *testing.T, f Fixture, name string, value any) []byte {
	t.Helper()
	var instance any
	if err := json.Unmarshal(mustJSON(t, f), &instance); err != nil {
		t.Fatal(err)
	}
	obj := topObject(t, instance)
	if _, ok := obj[name]; !ok {
		t.Fatalf("%s carries no %q at the compared level", f.ID, name)
	}
	obj[name] = value
	return mustMarshal(t, instance)
}

// reDescribe returns a copy of f whose recorded contract says topLevel/element,
// together with a body matching it. The point is to change the documented side
// of the shape check without editing an SDK file.
func reDescribe(t *testing.T, f Fixture, topLevel, elementType string, instance any) (Fixture, []byte) {
	t.Helper()
	edited := f
	edited.Checks.TopLevel = topLevel
	edited.Checks.ElementType = elementType
	return edited, mustMarshal(t, instance)
}

func mustMarshal(t *testing.T, v any) []byte {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func topObject(t *testing.T, instance any) map[string]any {
	t.Helper()
	switch v := instance.(type) {
	case map[string]any:
		return v
	case []any:
		if len(v) == 0 {
			t.Fatal("array body is empty")
		}
		obj, ok := v[0].(map[string]any)
		if !ok {
			t.Fatalf("array element is %T, not an object", v[0])
		}
		return obj
	default:
		t.Fatalf("body top level is %T", instance)
		return nil
	}
}

// decodeCommitted unmarshals the committed bytes into t, returning the error a
// caller would see.
func decodeCommitted(f Fixture, t reflect.Type) error {
	raw, err := f.Read()
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, reflect.New(t).Interface())
}

// withField rebuilds typ with one field replaced by edit, keeping every other
// field. reflect.StructOf is the only way to state "what if this field were a
// string, or carried no tag" without editing a production file, which this task
// must not do. The rebuilt type is a distinct type from the one it came from, so
// the SDK's own type is untouched and the compiled binary is unaffected.
func withField(t *testing.T, typ reflect.Type, field string, edit func(reflect.StructField) reflect.StructField) reflect.Type {
	t.Helper()
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	switch typ.Kind() {
	case reflect.Slice, reflect.Array:
		return reflect.SliceOf(withField(t, typ.Elem(), field, edit))
	case reflect.Struct:
		fields := make([]reflect.StructField, typ.NumField())
		replaced := false
		for i := range fields {
			fields[i] = typ.Field(i)
			if tagName(fields[i]) != field {
				continue
			}
			fields[i] = edit(fields[i])
			replaced = true
		}
		if !replaced {
			t.Fatalf("%s has no field tagged %q", typ, field)
		}
		return reflect.StructOf(fields)
	default:
		t.Fatalf("cannot rebuild %s", typ)
		return nil
	}
}

// structOfElem collapses a slice type to a struct built from its element's
// fields, which is how a top-level inversion is un-made without an SDK edit.
func structOfElem(t *testing.T, typ reflect.Type) reflect.Type {
	t.Helper()
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Slice && typ.Kind() != reflect.Array {
		t.Fatalf("%s is not a slice; the case under test is a slice decoding an object", typ)
	}
	elem := typ.Elem()
	for elem.Kind() == reflect.Pointer {
		elem = elem.Elem()
	}
	if elem.Kind() != reflect.Struct {
		t.Fatalf("%s element is %s, which has no fields to rebuild", typ, elem)
	}
	fields := make([]reflect.StructField, elem.NumField())
	for i := range fields {
		fields[i] = elem.Field(i)
	}
	return reflect.StructOf(fields)
}

// hasKind reports whether an outcome carries a divergence of that kind, and of
// that name when one is given.
func hasKind(o Outcome, kind DivergenceKind, name ...string) bool {
	for _, d := range o.Divergences {
		if d.Kind != kind {
			continue
		}
		if len(name) == 0 || d.Name == name[0] {
			return true
		}
	}
	return false
}

// skippedCheck reports whether an outcome records that a check could not run, so
// a test can assert that a check stood aside rather than silently passing. The
// failure mode this guards is the one the whole package exists to remove: a check
// that quietly stops being looked for.
func skippedCheck(o Outcome, c Check) bool {
	for _, s := range o.Skipped {
		if s.Check == c {
			return true
		}
	}
	return false
}

func countKind(o Outcome, kind DivergenceKind) int {
	n := 0
	for _, d := range o.Divergences {
		if d.Kind == kind {
			n++
		}
	}
	return n
}

// detailFor returns the Detail of the named divergence, or "" when absent.
func detailFor(o Outcome, kind DivergenceKind, name string) string {
	for _, d := range o.Divergences {
		if d.Kind == kind && d.Name == name {
			return d.Detail
		}
	}
	return ""
}

// summary renders an outcome on one line, for a log line or a failure message.
func summary(o Outcome) string {
	if len(o.Divergences) == 0 {
		return "no divergences"
	}
	parts := make([]string, 0, len(o.Divergences))
	for _, d := range o.Divergences {
		label := string(d.Kind)
		if d.Name != "" {
			label += "(" + d.Name + ")"
		}
		parts = append(parts, label)
	}
	return strings.Join(parts, ", ")
}
