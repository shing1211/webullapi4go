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

package webullapi4go

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/shing1211/webullapi4go/conformance"
)

// The status documents quote their own divergence counts, and those numbers are the
// first thing a reader trusts and the first thing that goes stale. Three separate
// premises in the document turned out to be false during v2.1.29, and the class table
// had to be corrected by hand once, in v2.1.27, because a hand-counted figure said 29
// symbols where the baseline and the gate both said 28.
//
// A number a person recomputes by hand is a number that will be wrong. This file makes
// the class table a derived value instead: it is parsed out of both status documents and
// compared against the committed baselines, so a count can only change by changing the
// baseline or the table together, and a table left behind by a baseline edit fails the
// build rather than misleading a reader.
//
// It checks the table and not the prose deliberately. The prose carries the same figures
// in many sentence forms, and a regexp over prose would either miss them or match the
// wrong occurrence; the table is the one place the numbers are structured. The prose
// quotes the table, and item 12 in the document records the correction made when it did
// not.

// statusDocs are the two documents that carry the table, maintained in step by hand.
var statusDocs = []string{"IMPLEMENTATION_STATUS.md", "docs/implementation-status.md"}

// tableRow is one parsed row of the class table.
type tableRow struct {
	// Label is the first cell, verbatim.
	Label string
	// Cells are the numeric cells after the label, in order.
	Cells []int
}

var (
	tableHeaderRE = regexp.MustCompile(`^\|\s*Class\s*\|\s*Rows\s*\|`)
	// A table cell is a number, possibly bolded, possibly an em dash for "none".
	cellRE = regexp.MustCompile(`^\*?\*?(\d+|—|-)\*?\*?$`)
	// Backticked kind names inside a label, e.g. the container-kind row names two.
	kindRE = regexp.MustCompile("`([a-z-]+)`")
)

// parseClassTable reads the class table out of a status document.
func parseClassTable(t *testing.T, path string) []tableRow {
	t.Helper()
	// readRepoFile is the package's existing reader for repository documents, and
	// carries the note about why a variable path is safe here; the paths come from
	// the fixed statusDocs slice and are not request-derived.
	text := readRepoFile(t, path)
	lines := strings.Split(text, "\n")
	start := -1
	for i, l := range lines {
		if tableHeaderRE.MatchString(strings.TrimSpace(l)) {
			start = i
			break
		}
	}
	if start < 0 {
		t.Fatalf("%s has no '| Class | Rows |' table", path)
	}
	var rows []tableRow
	for _, l := range lines[start+2:] {
		s := strings.TrimSpace(l)
		if !strings.HasPrefix(s, "|") {
			break
		}
		cells := strings.Split(strings.Trim(s, "|"), "|")
		label := strings.TrimSpace(cells[0])
		if label == "" {
			break
		}
		row := tableRow{Label: label}
		for _, c := range cells[1:] {
			c = strings.TrimSpace(c)
			if m := cellRE.FindStringSubmatch(c); m != nil {
				switch m[1] {
				case "—", "-":
					row.Cells = append(row.Cells, 0)
				default:
					n, _ := strconv.Atoi(m[1])
					row.Cells = append(row.Cells, n)
				}
				continue
			}
			// A trailing description column holds prose, not a number.
			break
		}
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		t.Fatalf("%s has a table header but no rows", path)
	}
	return rows
}

// kindsNamedIn returns the divergence kinds a row's label names, so the composite
// container-kind row is checked as the sum of the two kinds it names rather than as a
// kind of its own.
func kindsNamedIn(label string) []string {
	var out []string
	for _, m := range kindRE.FindAllStringSubmatch(label, -1) {
		out = append(out, m[1])
	}
	return out
}

// baselineEntry is the subset of a recorded divergence this test reads.
type baselineEntry struct {
	Symbol string
	Kind   string
}

func loadRootBaseline(t *testing.T) []baselineEntry {
	t.Helper()
	b, err := conformance.LoadBaseline()
	if err != nil {
		t.Fatalf("conformance.LoadBaseline: %v", err)
	}
	out := make([]baselineEntry, 0, len(b.Entries))
	for _, e := range b.Entries {
		out = append(out, baselineEntry{Symbol: e.Symbol, Kind: string(e.Kind)})
	}
	return out
}

// TestStatusClassTableAgreesWithTheBaseline is the gate.
//
// It recomputes every class from the committed baseline and compares. Both documents
// are checked, so the mirror cannot drift from the root, and each is checked against
// the baseline, so neither can drift from reality.
func TestStatusClassTableAgreesWithTheBaseline(t *testing.T) {
	base := loadRootBaseline(t)

	// Recompute per class, and per package, straight from the baseline.
	rowsByKind := map[string]int{}
	symsByKind := map[string]map[string]bool{}
	pkgByKind := map[string]map[string]int{}
	allSyms := map[string]bool{}
	pkgTotal := map[string]int{}
	for _, e := range base {
		rowsByKind[e.Kind]++
		if symsByKind[e.Kind] == nil {
			symsByKind[e.Kind] = map[string]bool{}
		}
		symsByKind[e.Kind][e.Symbol] = true
		if pkgByKind[e.Kind] == nil {
			pkgByKind[e.Kind] = map[string]int{}
		}
		pkg := strings.SplitN(e.Symbol, ".", 2)[0]
		pkgByKind[e.Kind][pkg]++
		allSyms[e.Symbol] = true
		pkgTotal[pkg]++
	}

	for _, doc := range statusDocs {
		t.Run(strings.ReplaceAll(doc, "/", "_"), func(t *testing.T) {
			rows := parseClassTable(t, doc)
			var totalRow *tableRow
			seen := map[string]bool{}
			for i := range rows {
				row := &rows[i]
				if strings.Contains(row.Label, "Total") {
					totalRow = row
					continue
				}
				kinds := kindsNamedIn(row.Label)
				if len(kinds) == 0 {
					t.Errorf("row %q names no divergence kind, so this test cannot check it; "+
						"a new class has to be added here or it is unchecked", row.Label)
					continue
				}
				wantRows, wantSyms := 0, map[string]bool{}
				wantPkg := map[string]int{}
				for _, k := range kinds {
					wantRows += rowsByKind[k]
					for s := range symsByKind[k] {
						wantSyms[s] = true
					}
					for p, n := range pkgByKind[k] {
						wantPkg[p] += n
					}
				}
				seen[kinds[0]] = true
				if len(row.Cells) < 5 {
					t.Errorf("row %q has %d numeric cell(s), want 5 (rows, symbols, and "+
						"three package columns)", row.Label, len(row.Cells))
					continue
				}
				got := row.Cells[:5]
				if got[0] != wantRows {
					t.Errorf("%s: %s rows = %d, baseline has %d", doc, row.Label, got[0], wantRows)
				}
				if got[1] != len(wantSyms) {
					t.Errorf("%s: %s symbols = %d, baseline has %d", doc, row.Label, got[1], len(wantSyms))
				}
				for i, p := range []string{"brokerfd", "data", "trade"} {
					if got[2+i] != wantPkg[p] {
						t.Errorf("%s: %s %s = %d, baseline has %d", doc, row.Label, p, got[2+i], wantPkg[p])
					}
				}
			}
			if totalRow == nil {
				t.Fatalf("%s has no Total row", doc)
			}
			if len(totalRow.Cells) < 5 {
				t.Fatalf("%s: the Total row has %d numeric cell(s), want 5",
					doc, len(totalRow.Cells))
			}
			got := totalRow.Cells[:5]
			if got[0] != len(base) {
				t.Errorf("%s: Total rows = %d, baseline has %d", doc, got[0], len(base))
			}
			if got[1] != len(allSyms) {
				t.Errorf("%s: Total symbols = %d, baseline has %d", doc, got[1], len(allSyms))
			}
			for i, p := range []string{"brokerfd", "data", "trade"} {
				if got[2+i] != pkgTotal[p] {
					t.Errorf("%s: Total %s = %d, baseline has %d", doc, p, got[2+i], pkgTotal[p])
				}
			}
		})
	}
}

// TestBothStatusDocumentsCarryTheSameTable stops the mirror drifting, which the gate
// above cannot see on its own because both would be equally wrong.
func TestBothStatusDocumentsCarryTheSameTable(t *testing.T) {
	a := parseClassTable(t, statusDocs[0])
	b := parseClassTable(t, statusDocs[1])
	if len(a) != len(b) {
		t.Fatalf("the two documents have %d and %d class-table rows", len(a), len(b))
	}
	for i := range a {
		if a[i].Label != b[i].Label {
			t.Errorf("row %d: %s has label %q, %s has %q", i, statusDocs[0], a[i].Label, statusDocs[1], b[i].Label)
			continue
		}
		if fmt.Sprint(a[i].Cells) != fmt.Sprint(b[i].Cells) {
			t.Errorf("row %q: %s has cells %v, %s has %v",
				a[i].Label, statusDocs[0], a[i].Cells, statusDocs[1], b[i].Cells)
		}
	}
}

// TestStatusTableComparisonBites is how the gate is known to work. A parser that found
// no table, or found one and compared nothing, would pass the gate above.
//
// The fixtures are written in the table's own shape, including the composite
// container-kind row that names two kinds in its label, so the two awkward shapes are
// the ones exercised.
func TestStatusTableComparisonBites(t *testing.T) {
	rows := []tableRow{
		{Label: "`missing-required-name`", Cells: []int{2, 2, 1, 1, 0}},
		{Label: "Container kind (`top-level-shape-mismatch` 1, `element-type-mismatch` 1)", Cells: []int{2, 2, 1, 1, 0}},
		{Label: "**Total**", Cells: []int{4, 2, 2, 2, 0}},
	}
	if got := cellsOf(rows, "`missing-required-name`"); got == nil {
		t.Fatal("the label was not found, so the comparison would check nothing")
	}
	// The composite row must resolve to the sum of the two kinds it names, not to
	// either one alone: 1 + 1 = 2 rows.
	if got := cellsOf(rows, "Container kind (`top-level-shape-mismatch` 1, `element-type-mismatch` 1)"); len(got.Cells) < 1 || got.Cells[0] != 2 {
		t.Errorf("the composite row resolved to rows %v, want 2 as the sum of two kinds", got.Cells)
	}
	if got := cellsOf(rows, "**Total**"); got == nil || got.Cells[0] != 4 {
		t.Errorf("the Total row resolved to %v, want 4 rows", got)
	}
	if got := cellsOf(rows, "`no-such-class`"); got != nil {
		t.Error("an absent label resolved, so the lookup cannot fail")
	}
}

func cellsOf(rows []tableRow, label string) *tableRow {
	for i := range rows {
		if rows[i].Label == label {
			return &rows[i]
		}
	}
	return nil
}
