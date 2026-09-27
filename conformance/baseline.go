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
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

// # Why a baseline rather than a skip list
//
// The comparison in shapes.go is correct about the current tree and the current
// tree is wrong in about a hundred places: brokerfd tags average_cost where the
// page requires cost_price, ten brokerfd rows decode an object where the page
// documents an array, and data.Quote.QuoteTime is an int64 where the page
// requires a string. A gate that reports those fails on day one, and a gate
// nobody runs is worth nothing. A gate that skips them has thrown the findings
// away, which is the specific outcome this instrument exists to prevent.
//
// So the divergences are recorded instead of filtered. The gate passes when the
// observed set *exactly equals* the committed set, and fails in both
// directions:
//
//   - a new divergence fails immediately, because it is not in the baseline;
//   - a baseline entry that no longer reproduces also fails, because a defect
//     that was fixed and a check that changed meaning look identical from here,
//     and both need a human to look at the tree.
//
// Every entry carries a Reason and a RecordedIn. That is not decoration: an
// entry with an empty reason is a suppressed finding wearing a comment, and
// LoadBaseline refuses to accept one. Requiring a reason does not make a reason
// true, and that is a separate failure this file has been bitten by: a false
// reason, repeated across many entries, satisfies the guard perfectly. The guard
// is kept because a reasonless entry is unambiguous, and the reasons are held to
// being per-entry true by review rather than by the schema.
//
// Editing the baseline is the whole of the workflow for adopting a change. It is
// a hand edit on purpose: nothing in a test run writes it, so a finding can only
// be adopted by a person who has read it.

// baselineFile is the committed baseline, embedded so the gate reads exactly the
// bytes that were committed and not whatever is on disk beside the source.
const baselineFile = "known-divergences.json"

//go:embed known-divergences.json
var baselineFS embed.FS

// KnownDivergence is one recorded divergence: the observation, plus why it is
// tolerated and where a reader can find the fuller account.
type KnownDivergence struct {
	// Divergence is the observation as the harness reported it. Every field
	// participates in the match, Detail included, so a check that starts
	// describing the same defect differently stops reproducing.
	Divergence
	// Reason says why the divergence is recorded rather than fixed. It is
	// required and must be non-empty; see the type comment.
	Reason string `json:"reason"`
	// RecordedIn points at the status document that carries the defect's
	// file:line, impact and unblock requirement. It is required so no entry is
	// merely "known" with nowhere to read about it.
	RecordedIn string `json:"recordedIn"`
}

// Key is the identity the gate matches an observation against. It is the
// Divergence key, so Reason and RecordedIn cannot mask a change in what was
// observed: those two fields are compared for presence, not for equality with
// anything the harness produces.
func (k KnownDivergence) Key() string { return k.Divergence.Key() }

// Coverage is the shape of the comparison the baseline was taken against. It is
// part of the baseline rather than a derived statistic so that coverage cannot
// shrink without the gate noticing: a manifest row added without a table entry
// changes these numbers, and a gate that only compared divergences would report
// green while checking less.
type Coverage struct {
	// Fixtures is the manifest's fixture count.
	Fixtures int `json:"fixtures"`
	// Compared is how many manifest rows a symbol was resolved for.
	Compared int `json:"compared"`
	// Unmapped must be empty. A non-empty value is a manifest symbol the table
	// has no type for, which would otherwise be a silent loss of coverage.
	Unmapped []string `json:"unmapped"`
	// OutOfScope names the broker/ module's symbols, unreachable from the root
	// module. Recorded so the gap is a number in a committed file rather than a
	// caveat in prose.
	OutOfScope []string `json:"outOfScope"`
	// NoSymbolRows counts the docgen em-dash rows, documented endpoints
	// deliberately mapped to no SDK symbol.
	NoSymbolRows int `json:"noSymbolRows"`
	// NotComparable counts rows that resolve to a Go type but whose page is not
	// the contract for the SDK call they are compared against, because the method
	// sends a different path. These rows are recorded here as a number rather than
	// left to be inferred from an absence of findings, because "not comparable" is
	// the outcome most easily misread as a pass and the one a reader most needs to
	// count.
	NotComparable int `json:"notComparable"`
	// SymbolsInTable is the size of SDKTypes, so a table edit is visible.
	SymbolsInTable int `json:"symbolsInTable"`
}

// Baseline is the committed set of known divergences plus the coverage it was
// taken against.
type Baseline struct {
	// BaselineCommit is the tree this set was observed on. It is provenance, not
	// a mechanism: nothing reads it, because a commit hash is not a substitute
	// for the set it describes.
	BaselineCommit string `json:"baselineCommit"`
	// Policy states the gate rule in the file itself, so the rule travels with
	// the data it governs.
	Policy string `json:"policy"`
	// Notes carry the free-text explanations the file needs, chiefly the ones a
	// reader cannot get from a count: what the name-coverage number is a floor
	// against, and why there is no waived entry.
	Notes map[string]string `json:"notes"`
	// Coverage is what was compared.
	Coverage Coverage `json:"coverage"`
	// Entries are the recorded divergences, sorted by key.
	Entries []KnownDivergence `json:"entries"`
}

// ByKey indexes the entries for lookup. The gate matches an observation against
// it, so a missing key is a new divergence and an unused key is a baseline entry
// that no longer reproduces.
func (b *Baseline) ByKey() map[string]KnownDivergence {
	out := make(map[string]KnownDivergence, len(b.Entries))
	for _, e := range b.Entries {
		out[e.Key()] = e
	}
	return out
}

// LoadBaseline reads and validates the committed baseline.
func LoadBaseline() (*Baseline, error) {
	raw, err := baselineFS.ReadFile(baselineFile)
	if err != nil {
		return nil, fmt.Errorf("conformance: read baseline: %w", err)
	}
	var b Baseline
	if err := json.Unmarshal(raw, &b); err != nil {
		return nil, fmt.Errorf("conformance: parse baseline: %w", err)
	}
	if err := b.Validate(); err != nil {
		return nil, fmt.Errorf("conformance: baseline: %w", err)
	}
	return &b, nil
}

// Validate refuses a baseline that could act as a suppression.
//
// The two required fields are the whole point. An entry whose Reason is empty
// records a defect nobody has agreed to look at, and one whose RecordedIn is
// empty records a defect nobody can find written down; both are how a finding
// disappears while the file still looks maintained.
func (b *Baseline) Validate() error {
	if len(b.Entries) == 0 {
		return fmt.Errorf("no entries: a baseline with nothing in it is a gate that checks nothing")
	}
	if strings.TrimSpace(b.Policy) == "" {
		return fmt.Errorf("no policy statement")
	}
	seen := make(map[string]KnownDivergence, len(b.Entries))
	prev := ""
	for _, e := range b.Entries {
		if strings.TrimSpace(e.Reason) == "" {
			return fmt.Errorf("%s / %s / %s %q: entry has no reason, which is a "+
				"suppressed finding rather than a recorded one",
				e.Symbol, e.Fixture, e.Kind, e.Name)
		}
		if strings.TrimSpace(e.RecordedIn) == "" {
			return fmt.Errorf("%s / %s / %s %q: entry names no document recording the defect",
				e.Symbol, e.Fixture, e.Kind, e.Name)
		}
		if strings.TrimSpace(e.Detail) == "" {
			return fmt.Errorf("%s / %s / %s %q: entry has no detail, so nothing says what disagrees",
				e.Symbol, e.Fixture, e.Kind, e.Name)
		}
		key := e.Key()
		if prev != "" && key < prev {
			return fmt.Errorf("entries are not sorted: %q comes after %q", key, prev)
		}
		prev = key
		if dup, ok := seen[key]; ok {
			return fmt.Errorf("%s is recorded twice: %q and %q", key, dup.Reason, e.Reason)
		}
		seen[key] = e
	}
	return nil
}
