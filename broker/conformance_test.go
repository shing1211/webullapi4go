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

package broker

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/shing1211/webullapi4go/conformance"
)

// This module carries the Broker API HK surface, which is a separate Go module and
// was therefore outside the reach of the wire-conformance harness in the root
// module: its 29 documented endpoints appeared in the manifest but no symbol in
// the root module could name a broker type, so every one of them was an
// out-of-scope placeholder. That is a larger blind spot than the one the root
// harness spent v2.1.22 and v2.1.23 working on, and it was open by construction
// rather than by accident.
//
// What closes it is that the harness is reachable from here. This module already
// requires the root module through a replace directive, so conformance.Load,
// conformance.Compare and the envelope resolver are all importable, and every check
// they perform is the same code the root module runs against its own 154 rows.
// Nothing in this file reimplements a check, and a change to a check therefore
// reaches these 29 rows as well as the root module's 285 recorded divergences.
//
// What is duplicated, and has to be, is the symbol table. conformance.SDKTypes is
// a map literal of reflect.TypeOf calls in the root module, which is exactly why
// it cannot name a broker type: the same unexported-type problem the envelope
// machinery exists to solve. brokerDecodeTypes below is the same fact written from
// the other side of the module boundary, and TestBrokerTableCoversItsManifest
// asserts it in both directions so a hole fails the build rather than quietly
// reducing coverage.

//go:embed conformance-divergences.json
var brokerBaselineJSON []byte

// brokerBaseline is the recorded divergence set for this module.
//
// The gate is the same shape as the root module's and for the same reason: a
// divergence not listed is a new defect and fails, and a listed entry that no
// longer reproduces also fails, because a fixed defect and a check that changed
// meaning are indistinguishable from here. Every entry carries a reason, and an
// empty one is refused rather than accepted, because a reasonless entry is a
// suppressed finding wearing a comment.
type brokerBaseline struct {
	// Policy states the gate rule, so the rule travels with the data it governs.
	Policy string `json:"policy"`
	// Coverage is what was compared, so coverage cannot shrink unnoticed.
	Coverage brokerCoverage `json:"coverage"`
	// Entries are the recorded divergences, sorted by key.
	Entries []brokerEntry `json:"entries"`
}

type brokerCoverage struct {
	// Rows is how many manifest rows name a broker symbol.
	Rows int `json:"rows"`
	// Compared is how many of those resolved to a type and were compared.
	Compared int `json:"compared"`
	// Envelopes is how many compared an unexported source-read envelope.
	Envelopes int `json:"envelopes"`
	// NoBody is how many decode no response body at all.
	NoBody int `json:"noBody"`
	// NotComparable is how many claim nothing because the page is not the
	// contract for the call the manifest names it against.
	NotComparable int `json:"notComparable"`
}

type brokerEntry struct {
	conformance.Divergence
	// Reason says why the divergence is recorded rather than fixed, and is
	// required to be non-empty.
	Reason string `json:"reason"`
	// RecordedIn points at the status document carrying the defect's fuller
	// account. It is required so no entry is merely "known" with nowhere to read.
	RecordedIn string `json:"recordedIn"`
}

// brokerSymbolPrefix is the manifest subject prefix this module owns.
const brokerSymbolPrefix = "broker."

func init() {
	// The two envelopes below are unexported, so neither can appear in a map
	// literal in the root module. They are registered here, once, before any test
	// runs, and the root module's parser then reads the declarations from these
	// sources and rebuilds them exactly as it does its own six. Registering a name
	// the root module already resolved is refused, so this cannot redefine a fact.
	for spelling, typ := range map[string]reflect.Type{
		"broker.Position":       reflect.TypeOf(Position{}),
		"broker.VirtualAccount": reflect.TypeOf(VirtualAccount{}),
	} {
		if err := conformance.RegisterEnvelopeNamedType(spelling, typ); err != nil {
			// A registration failure here would otherwise surface much later as an
			// unresolvable field type, which reads like a missing field rather than
			// like a wiring mistake.
			panic("broker: register envelope type " + spelling + ": " + err.Error())
		}
	}
}

// brokerEnvelopes maps a symbol to the unexported type its method's `var out`
// declares, for the rows where the public projection is not the compared type.
//
// The conformance table's own rule is that the type named is the one the method's
// `var out` declares and not the one it returns, which is why these two rows
// compare an envelope rather than the slice their methods hand back: the page
// documents an object and the SDK decodes one, and comparing against []Position
// would report a shape inversion that the method does not have.
var brokerEnvelopes = map[string]string{
	"broker.GetPositions":        "broker.positionsResponse",
	"broker.ListVirtualAccounts": "broker.virtualAccountsResponse",
}

// brokerDecodeTypes is the symbol table written from this side of the boundary.
// broker.CancelOrder is absent because it decodes no response body: it passes nil
// as the out parameter and returns only an error, so there is no type to compare
// and the row is counted as a no-body row instead.
var brokerDecodeTypes = map[string]reflect.Type{
	"broker.CreateCashJournal":         reflect.TypeOf(CashJournal{}),
	"broker.CreateFXExchange":          reflect.TypeOf(FXExchange{}),
	"broker.CreateInstantExchange":     reflect.TypeOf(InstantExchange{}),
	"broker.CreateInstantFunding":      reflect.TypeOf(InstantFunding{}),
	"broker.CreatePositionJournal":     reflect.TypeOf(PositionJournal{}),
	"broker.CreateVirtualAccount":      reflect.TypeOf(VirtualAccount{}),
	"broker.GetBalance":                reflect.TypeOf(Balance{}),
	"broker.GetCashActivities":         reflect.TypeOf([]CashActivity{}),
	"broker.GetCashJournalDetail":      reflect.TypeOf(CashJournal{}),
	"broker.GetCorporateActionsDetail": reflect.TypeOf([]CorporateActionDetail{}),
	"broker.GetFXExchangeDetail":       reflect.TypeOf(FXExchange{}),
	"broker.GetFXRate":                 reflect.TypeOf(FXRate{}),
	"broker.GetInstantExchangeDetail":  reflect.TypeOf(InstantExchange{}),
	"broker.GetInstantFundingDetail":   reflect.TypeOf(InstantFunding{}),
	"broker.GetOpenOrders":             reflect.TypeOf([]BrokerOrder{}),
	"broker.GetOrderDetail":            reflect.TypeOf(BrokerOrder{}),
	"broker.GetOrderHistory":           reflect.TypeOf([]BrokerOrder{}),
	"broker.GetPositionJournalDetail":  reflect.TypeOf(PositionJournal{}),
	"broker.GetStockInstruments":       reflect.TypeOf([]StockInstrument{}),
	"broker.GetStockLocate":            reflect.TypeOf([]StockLocate{}),
	"broker.GetTradeCalendar":          reflect.TypeOf([]TradeCalendar{}),
	"broker.GetVirtualAccount":         reflect.TypeOf(VirtualAccount{}),
	"broker.PlaceOrder":                reflect.TypeOf(BrokerOrder{}),
	"broker.PreviewOrder":              reflect.TypeOf(OrderPreview{}),
	"broker.ReplaceOrder":              reflect.TypeOf(BrokerOrder{}),
	"broker.UpdateVirtualAccount":      reflect.TypeOf(VirtualAccount{}),
}

// brokerNoBodySymbols are the rows whose method decodes no response body, listed
// explicitly rather than left as a gap in the table. A symbol missing from both
// maps and from this one is a table hole, and TestBrokerTableCoversItsManifest
// fails on it; the alternative is a row that silently stops being compared.
var brokerNoBodySymbols = []string{"broker.CancelOrder"}

func loadBrokerManifest(t *testing.T) *conformance.Manifest {
	t.Helper()
	m, err := conformance.Load()
	if err != nil {
		t.Fatalf("conformance.Load: %v", err)
	}
	return m
}

func brokerRows(m *conformance.Manifest) []conformance.Fixture {
	var rows []conformance.Fixture
	for _, f := range m.Fixtures {
		if strings.HasPrefix(f.SDKSymbol, brokerSymbolPrefix) {
			rows = append(rows, f)
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].SDKSymbol < rows[j].SDKSymbol })
	return rows
}

// brokerDecodeTarget resolves a symbol to the type its method's `var out`
// declares, or nil when the method decodes no body.
func brokerDecodeTarget(t *testing.T, sym string) reflect.Type {
	t.Helper()
	if spelling, ok := brokerEnvelopes[sym]; ok {
		typ, err := conformance.EnvelopeType(spelling)
		if err != nil {
			t.Fatalf("%s: resolve envelope %s: %v", sym, spelling, err)
		}
		return typ
	}
	if typ, ok := brokerDecodeTypes[sym]; ok {
		return typ
	}
	for _, noBody := range brokerNoBodySymbols {
		if noBody == sym {
			return nil
		}
	}
	t.Fatalf("%s is in neither brokerDecodeTypes, brokerEnvelopes nor "+
		"brokerNoBodySymbols, so this file has a hole where the manifest has a row", sym)
	return nil
}

// TestBrokerResponseContracts is the gate for the 29 documented Broker API HK
// endpoints: the same five checks the root module runs, over this module's types.
func TestBrokerResponseContracts(t *testing.T) {
	m := loadBrokerManifest(t)
	rows := brokerRows(m)
	if len(rows) == 0 {
		t.Fatal("the manifest names no broker symbol, so this test would check nothing")
	}
	base, err := loadBrokerBaseline()
	if err != nil {
		t.Fatal(err)
	}

	got := map[string]conformance.Divergence{}
	var cov brokerCoverage
	for _, f := range rows {
		cov.Rows++
		sym := f.SDKSymbol
		target := brokerDecodeTarget(t, sym)
		if target == nil {
			cov.NoBody++
			continue
		}
		if _, ok := brokerEnvelopes[sym]; ok {
			cov.Envelopes++
		}
		outcome := conformance.Compare(f, sym, target)
		if outcome.NotComparable != nil {
			cov.NotComparable++
			continue
		}
		cov.Compared++
		for _, d := range outcome.Divergences {
			if _, dup := got[d.Key()]; dup {
				t.Errorf("divergence reported twice: %s", d.Key())
			}
			got[d.Key()] = d
		}
		// The verb check runs here as well as in the root module, because this is
		// the module where it earned its place: UpdateVirtualAccount issued PUT
		// where the page documents POST, and Compare above reports its path as a
		// clean match because it compares paths and not verbs. Its finding joins the
		// same exact-set gate rather than being asserted separately, so a future
		// verb disagreement cannot be recorded in one baseline and forgotten in the
		// other.
		for _, d := range conformance.CompareVerb(f, sym) {
			if _, dup := got[d.Key()]; dup {
				t.Errorf("divergence reported twice: %s", d.Key())
			}
			got[d.Key()] = d
		}
	}
	t.Logf("compared %d of %d rows: %d envelope, %d no-body, %d not-comparable; "+
		"%d divergence(s) observed against %d recorded",
		cov.Compared, cov.Rows, cov.Envelopes, cov.NoBody, cov.NotComparable,
		len(got), len(base.Entries))

	// Coverage first, so a table that lost a row cannot quietly compare less.
	if cov != base.Coverage {
		t.Errorf("coverage is %+v, baseline says %+v; a manifest row, a table entry or "+
			"the baseline was changed, and comparing less must not read as green", cov, base.Coverage)
	}

	want := map[string]brokerEntry{}
	for _, e := range base.Entries {
		want[e.Key()] = e
	}
	var newKeys, staleKeys []string
	for key, d := range got {
		if _, ok := want[key]; !ok {
			newKeys = append(newKeys, fmt.Sprintf("%s\n    %s %s %s\n      %s",
				key, d.Symbol, d.Kind, d.Name, d.Detail))
		}
	}
	for key, e := range want {
		if _, ok := got[key]; !ok {
			staleKeys = append(staleKeys, fmt.Sprintf("%s\n    %s %s %s\n      %s\n"+
				"      reason on file: %s", key, e.Symbol, e.Kind, e.Name, e.Detail, e.Reason))
		}
	}
	sort.Strings(newKeys)
	sort.Strings(staleKeys)
	for _, s := range newKeys {
		t.Errorf("NEW DIVERGENCE, not in the baseline: %s", s)
	}
	for _, s := range staleKeys {
		t.Errorf("BASELINE ENTRY NO LONGER REPRODUCES: %s\n      Either the defect was "+
			"fixed, or a check changed what it means. Both need a person to look.", s)
	}
	if len(newKeys) == 0 && len(staleKeys) == 0 {
		t.Logf("observed %d divergences, matching all %d baseline entries exactly",
			len(got), len(want))
	}
}

// TestBrokerTableCoversItsManifest holds the table to the manifest in both
// directions, so a hole in it fails the build rather than quietly reducing what
// this module is measured against.
//
// The forward direction is the point: every manifest row naming a broker symbol
// must resolve. The reverse direction matters because a stale entry would
// otherwise be dead weight that still looks like coverage.
func TestBrokerTableCoversItsManifest(t *testing.T) {
	m := loadBrokerManifest(t)
	rows := brokerRows(m)
	subjects := map[string]bool{}
	for _, f := range rows {
		subjects[f.SDKSymbol] = true
		if _, ok := brokerDecodeTypes[f.SDKSymbol]; ok {
			continue
		}
		if _, ok := brokerEnvelopes[f.SDKSymbol]; ok {
			continue
		}
		found := false
		for _, noBody := range brokerNoBodySymbols {
			if noBody == f.SDKSymbol {
				found = true
			}
		}
		if !found {
			t.Errorf("manifest row %s names %s, which no table here can resolve",
				f.ID, f.SDKSymbol)
		}
	}
	for sym := range brokerDecodeTypes {
		if !subjects[sym] {
			t.Errorf("brokerDecodeTypes names %s, which no manifest row refers to, so "+
				"it is dead weight that still looks like coverage", sym)
		}
	}
	for sym := range brokerEnvelopes {
		if !subjects[sym] {
			t.Errorf("brokerEnvelopes names %s, which no manifest row refers to", sym)
		}
	}
	for _, sym := range brokerNoBodySymbols {
		if !subjects[sym] {
			t.Errorf("brokerNoBodySymbols names %s, which no manifest row refers to", sym)
		}
	}
	t.Logf("%d manifest row(s) for %d symbols, all resolved", len(rows), len(subjects))
}

func loadBrokerBaseline() (*brokerBaseline, error) {
	var b brokerBaseline
	if err := json.Unmarshal(brokerBaselineJSON, &b); err != nil {
		return nil, fmt.Errorf("decode conformance-divergences.json: %w", err)
	}
	if strings.TrimSpace(b.Policy) == "" {
		return nil, fmt.Errorf("conformance-divergences.json states no policy, so the " +
			"gate rule does not travel with the data it governs")
	}
	seen := map[string]string{}
	for i, e := range b.Entries {
		if strings.TrimSpace(e.Reason) == "" {
			return nil, fmt.Errorf("entry %d (%s %s %s) has an empty reason, which is a "+
				"suppressed finding wearing a comment", i, e.Symbol, e.Kind, e.Name)
		}
		if strings.TrimSpace(e.RecordedIn) == "" {
			return nil, fmt.Errorf("entry %d (%s %s %s) has an empty recordedIn, so it is "+
				"merely known with nowhere to read about it", i, e.Symbol, e.Kind, e.Name)
		}
		if prev, dup := seen[e.Key()]; dup {
			return nil, fmt.Errorf("entry %d duplicates an earlier entry: %s", i, prev)
		}
		seen[e.Key()] = fmt.Sprintf("%s %s %s", e.Symbol, e.Kind, e.Name)
	}
	if !sort.SliceIsSorted(b.Entries, func(i, j int) bool {
		return b.Entries[i].Key() < b.Entries[j].Key()
	}) {
		return nil, fmt.Errorf("conformance-divergences.json is not sorted by key, so two " +
			"people reading it would not see the same order")
	}
	return &b, nil
}
