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
	"strings"
	"testing"
)

// The verb check is new, and a new check that has never been observed failing is
// not evidence of anything. broker.UpdateVirtualAccount reconciled as a clean path
// match for the whole life of the harness while issuing PUT where the page
// documents POST, so the three tests below are what make a green run mean
// something: the extractor is pinned to specific methods, the check is required to
// cover every symbol, and the comparison is shown to fire on a disagreement.

// verbSummary renders a verb finding for a failure message. summary takes an
// Outcome and this check returns a slice, so it gets its own joiner rather than a
// wrapper that would obscure which function produced the finding.
func verbSummary(ds []Divergence) string {
	if len(ds) == 0 {
		return "no findings"
	}
	parts := make([]string, 0, len(ds))
	for _, d := range ds {
		parts = append(parts, string(d.Kind)+" "+d.Name+": "+d.Detail)
	}
	return strings.Join(parts, "; ")
}

// TestVerbExtractorReadsKnownMethods pins the extraction itself.
//
// Each case names a method whose verb is settled by reading its source, and the
// interesting ones are the awkward shapes rather than the easy ones: a helper on
// the client's own receiver, a helper on a sub-service reached through a call, a
// net/http constant, and a method that also reads query parameters, which is the
// case a naive extractor mistakes for a GET.
func TestVerbExtractorReadsKnownMethods(t *testing.T) {
	tests := []struct {
		symbol string
		want   string
		why    string
	}{
		{"data.GetSnapshot", "GET",
			"c.get on the client's own receiver"},
		{"data.GetLogos", "POST",
			"a sub-service helper reached as c.DisplayService().Post(...)"},
		{"trade.GetOpenOrders", "GET",
			"trade's own get helper"},
		{"broker.UpdateVirtualAccount", "POST",
			"the method item 18 corrected; it issued PUT while the page documents POST"},
		{"broker.CancelOrder", "POST",
			"a write issued as POST, which is the house style across the broker package"},
	}
	for _, tc := range tests {
		t.Run(tc.symbol, func(t *testing.T) {
			got := SDKVerbOf(tc.symbol)
			if !got.Known() {
				t.Fatalf("SDKVerbOf(%s) is unknown: %s", tc.symbol, got.Reason)
			}
			if got.Verb != tc.want {
				t.Errorf("SDKVerbOf(%s) = %s, want %s (%s)", tc.symbol, got.Verb, tc.want, tc.why)
			}
		})
	}
}

// TestVerbExtractorCoversEverySymbol requires the extractor to read a verb for
// every symbol the manifest maps.
//
// This is the guard on the design decision in CompareVerb: a symbol with no
// readable verb produces no divergence, so its row would go green on a check that
// did not run. That is safe only while unreadable is zero, and only if something
// enforces it. The check reports rather than fails, because the baseline must stay
// a record of facts about the SDK; this test is where a gap in the extractor
// belongs, since that is a defect in the tool rather than in the code under test.
//
// A method that sends several verbs is counted separately and only logged. It is
// not unreadable and not a disagreement: client.CreateToken creates a token, polls
// its status and revokes it, so it has no single counterpart verb to compare
// against one page. It is listed so the skip is visible rather than assumed.
func TestVerbExtractorCoversEverySymbol(t *testing.T) {
	var unreadable, multi []string
	checked := 0
	for _, f := range loadManifest(t).Fixtures {
		// The em dash is the docgen placeholder for an endpoint deliberately mapped
		// to no SDK symbol. It is not an unreadable verb, and counting it as one
		// would be a finding about the manifest rather than about the extractor.
		if f.SDKSymbol == "" || f.SDKSymbol == NoSDKSymbol {
			continue
		}
		sym, _, named := Subject(f.SDKSymbol)
		if !named {
			unreadable = append(unreadable, f.SDKSymbol+
				": Subject() did not recognise the manifest symbol")
			continue
		}
		checked++
		v := SDKVerbOf(sym)
		switch {
		case !v.Known():
			unreadable = append(unreadable, sym+": "+v.Reason)
		case !v.Single():
			multi = append(multi, sym+": "+v.Reason)
		}
	}
	if checked == 0 {
		t.Fatal("no mapped symbols, so this checked nothing")
	}
	if len(unreadable) > 0 {
		t.Errorf("%d of %d mapped symbols have no readable verb, so their rows are green "+
			"on the verb check without having been examined:", len(unreadable), checked)
		for _, u := range unreadable {
			t.Errorf("  %s", u)
		}
	}
	t.Logf("%d symbols: all readable, %d send several verbs and are not comparable on one",
		checked, len(multi))
	for _, m := range multi {
		t.Logf("  not comparable on verb: %s", m)
	}
}

// TestVerbAgreement is the gate: every comparable row sends the verb its page
// documents.
//
// It reports rather than fails, because a verb disagreement is a fact about the
// SDK that belongs in the baseline rather than a broken build. A row that
// disagrees without a baseline entry is caught by the baseline's own exact-set
// test, which is the one place that can say "you found something and did not
// record it".
func TestVerbAgreement(t *testing.T) {
	var rows []string
	agree, pathSkipped := 0, 0
	for _, f := range loadManifest(t).Fixtures {
		sym, _, named := Subject(f.SDKSymbol)
		if !named || f.Documented.Method == "" {
			continue
		}
		if f.SDKPathMatch != PathMatchSame {
			// Not this page's call, so there is no counterpart verb to compare.
			pathSkipped++
			continue
		}
		ds := CompareVerb(f, sym)
		if len(ds) == 0 {
			agree++
			continue
		}
		rows = append(rows, sym+"\n      "+ds[0].Detail)
	}
	t.Logf("verb agreement: %d agree, %d disagree, %d not comparable on path", agree, len(rows), pathSkipped)
	if len(rows) > 0 {
		t.Errorf("%d row(s) send a verb their page does not document:", len(rows))
		for _, r := range rows {
			t.Errorf("  %s", r)
		}
	}
}

// TestVerbComparisonBites is how the check is known to fire.
//
// The disagreement is manufactured by changing the documented verb on a real
// fixture, so the SDK source, the extractor and the real symbol are all exercised;
// only the page's side is falsified. An extractor that returned "" for everything
// would pass TestVerbAgreement, which is exactly why this test exists.
func TestVerbComparisonBites(t *testing.T) {
	const id = "market-data-stock/GET-market-data-stocks-snapshots-list"
	f := mustFixture(t, loadManifest(t), id)
	if f.Documented.Method != "GET" {
		t.Fatalf("fixture %s documents %s, want GET; the fixture changed and this "+
			"test's premise no longer holds", id, f.Documented.Method)
	}
	// Unmodified: no finding.
	if ds := CompareVerb(f, "data.GetSnapshot"); len(ds) != 0 {
		t.Fatalf("unmodified fixture: %d finding(s), want 0: %s", len(ds), verbSummary(ds))
	}

	// The page claims DELETE. The SDK still sends GET, so this must be reported,
	// and the report must name both sides.
	falsified := f
	falsified.Documented.Method = "DELETE"
	ds := CompareVerb(falsified, "data.GetSnapshot")
	if len(ds) != 1 {
		t.Fatalf("a falsified verb produced %d finding(s), want exactly 1: %s", len(ds), verbSummary(ds))
	}
	d := ds[0]
	if d.Kind != VerbMismatch {
		t.Errorf("kind = %q, want %q", d.Kind, VerbMismatch)
	}
	if checkForKind(d.Kind) != CheckVerb {
		t.Errorf("checkForKind(%q) = %q, want %q", d.Kind, checkForKind(d.Kind), CheckVerb)
	}
	if d.Name != "DELETE" {
		t.Errorf("name = %q, want DELETE", d.Name)
	}
	if !strings.Contains(d.Detail, "DELETE") || !strings.Contains(d.Detail, "GET") {
		t.Errorf("detail must carry both sides so a report needs no lookup, got %q", d.Detail)
	}

	// Lower-case on the page must not hide a disagreement: the manifest writes the
	// verb upper-cased, but the comparison must not depend on that.
	lower := f
	lower.Documented.Method = "delete"
	if ds := CompareVerb(lower, "data.GetSnapshot"); len(ds) != 1 {
		t.Errorf("a lower-case documented verb produced %d finding(s), want 1", len(ds))
	}
}

// TestVerbCheckSkipsRowsItCannotCompare states the two exclusions, because a
// check that silently skips is a check whose skips have to be readable.
func TestVerbCheckSkipsRowsItCannotCompare(t *testing.T) {
	m := loadManifest(t)

	// A row whose path differs is not this page's call, so no verb claim is made
	// about it; the path divergence already reports the row as not comparable.
	notSame, ok := m.ByID("market-data-stock/GET-market-data-stocks-snapshots-list")
	if !ok {
		t.Fatal("fixture not found")
	}
	notSame.SDKPathMatch = PathMatchNone
	if ds := CompareVerb(notSame, "data.GetSnapshot"); len(ds) != 0 {
		t.Errorf("a path-mismatched row produced %d verb finding(s), want 0: %s",
			len(ds), verbSummary(ds))
	}

	// A page that documents no verb cannot be disagreed with.
	noMethod := notSame
	noMethod.SDKPathMatch = PathMatchSame
	noMethod.Documented.Method = ""
	if ds := CompareVerb(noMethod, "data.GetSnapshot"); len(ds) != 0 {
		t.Errorf("a page documenting no verb produced %d finding(s), want 0", len(ds))
	}
}

// TestVerbExtractorRejectsAnUnreadableSymbol keeps the unknown case honest: a
// symbol that does not resolve must say so rather than default to a verb.
func TestVerbExtractorRejectsAnUnreadableSymbol(t *testing.T) {
	for _, sym := range []string{
		"data.NoSuchMethod",
		"nosuchpackage.SomeMethod",
		"notASymbol",
	} {
		v := SDKVerbOf(sym)
		if v.Known() {
			t.Errorf("SDKVerbOf(%q) = %s, want unknown", sym, v.Verb)
		}
		if v.Reason == "" {
			t.Errorf("SDKVerbOf(%q) gave no reason for being unknown", sym)
		}
	}
}
