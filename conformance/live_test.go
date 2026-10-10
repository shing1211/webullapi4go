// Copyright 2026 shing1211
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "as IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conformance

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// These tests hold the live comparison: the same five checks shapes.go runs, run
// a second time against the shape the sandbox actually sent. They matter for one
// reason the documented comparison cannot serve: a page and a server can agree
// with each other and both disagree with the SDK, and only one of those two
// disagreements is what a caller experiences.
//
// Every test here is written so that inverting the comparison makes it fail. A
// test that would pass with the two directions swapped is not evidence.

// A liveSkeleton turns a body literal into the bytes the harness reads, so a test
// states the wire rather than a path into the committed tree.
func liveSkeleton(t *testing.T, body string) []byte {
	t.Helper()
	if !json.Valid([]byte(body)) {
		t.Fatalf("the test body is not valid JSON: %s", body)
	}
	return []byte(body)
}

// liveFixture builds the minimal Fixture a comparison needs, from the documented
// checks a test states outright. The checks are the page's side of the
// comparison; the body is the live side, and CompareLive must not confuse them.
//
// The documented half of the comparison reads the fixture's committed instance
// through Fixture.Read, so a synthetic fixture has to point at a real one or that
// half cannot run. The tests below that need a documented half therefore use
// mustLiveFixture, which is a stronger test than this anyway; this is for the ones
// that need to control the checks, and it names the identity they share.
func liveFixture(id, topLevel string, required, declared []string) Fixture {
	return Fixture{
		ID: id,
		Checks: Checks{
			TopLevel:              topLevel,
			DeclaresRequired:      len(required) > 0,
			RequiredNames:         required,
			DeclaredTopLevelNames: declared,
		},
		Documented:   Documented{Method: "GET", Path: "/example/get"},
		SDKPath:      "/example/get",
		SDKPathMatch: PathMatchSame,
	}
}

// findings indexes a slice of live divergences by their key, so a test can assert
// on presence and absence by name rather than on a count.
func findings(rows []LiveDivergence) map[string][]LiveDivergence {
	out := map[string][]LiveDivergence{}
	for _, r := range rows {
		out[r.Key()] = append(out[r.Key()], r)
	}
	return out
}

// hasDirection reports whether rows holds a finding with the given pair key and
// direction.
//
// It matches on the pair key and then on the direction separately, because the
// two are separate fields: the pair key is what identifies a finding across the
// two runs, and the direction is what says which run observed it. A helper that
// matched on the record key alone could not answer the question at all, since
// that key already contains the direction.
func hasDirection(rows []LiveDivergence, pairKey string, d LiveDirection) bool {
	for _, r := range rows {
		if r.PairKey() == pairKey && r.Direction == d {
			return true
		}
	}
	return false
}

// TestCompareLiveReportsTheSDKDirection is the direction that matters.
//
// A name the SDK cannot reach is a silent zero for a caller, and a skeleton is
// the only evidence that says the server sent it. The comparison must report it
// as a live-sdk finding, and it must not report it against the documentation
// either: the fixture is a statement about a page, and a live disagreement the
// page does not describe is precisely what the two runs are there to separate.
//
// It is driven from the committed tree rather than from a synthetic fixture,
// because the documented direction needs a committed instance to compare against
// and a made-up one has none. That also makes it a test of the real thing: this
// is the exact shape of the finding the baseline records for
// data.GetFuturesTick, where the page requires instrument_id, the SDK tags it,
// and the live body sends instrumentId. The SDK has since grown a decoder that
// accepts either spelling, and the row is still reported, which is the point
// this test holds: the check reads tag inventory, and a custom UnmarshalJSON
// does not change a tag.
//
// The way to break this test is to run the checks once, against the fixture, and
// report the result in both directions.
func TestCompareLiveReportsTheSDKDirection(t *testing.T) {
	const symbol = "data.GetFuturesTick"
	m := loadManifest(t)
	f := mustLiveFixture(t, m, symbol)
	sk := mustLiveSkeleton(t, symbol)

	rows := CompareLive(f, symbol, mustType(t, symbol), sk)

	const name = "instrument_id"
	want := LiveDivergence{
		Symbol:    symbol,
		Fixture:   f.ID,
		Direction: LiveSDKDirection,
		Kind:      MissingRequiredName,
		Name:      name,
	}
	byKey := findings(rows)
	got, ok := byKey[want.Key()]
	if !ok {
		t.Fatalf("no finding for the name the SDK cannot reach; got %v", rows)
	}
	if len(got) != 1 {
		t.Errorf("the key %q carries %d findings, want 1: the same disagreement reported "+
			"twice is a duplicate that reads as two", want.Key(), len(got))
	}
	if got[0].Symbol != want.Symbol || got[0].Fixture != want.Fixture ||
		got[0].Kind != want.Kind || got[0].Name != want.Name {
		t.Errorf("finding = %+v, want %+v", got[0], want)
	}
	if hasDirection(rows, want.PairKey(), LiveDocsDirection) {
		t.Errorf("the live wire is also reported against the documentation, so a "+
			"server-only finding reads as a page-versus-SDK one: %v", rows)
	}

	// The detail has to name the observation rather than restate the kind, and it
	// has to say what the server sent instead, because "the name is missing" on
	// its own does not tell a reader there is a near-miss spelling of it. It stops
	// at the tag on purpose: what a caller then reads is not a thing this check
	// inspects, and the SDK now decodes the member its tag cannot match.
	if got[0].Detail == "" {
		t.Fatal("the finding carries no detail, so nothing says what disagrees")
	}
	if strings.Contains(got[0].Detail, "documented") {
		t.Errorf("a live finding describes the documentation: %q", got[0].Detail)
	}
	if !strings.Contains(got[0].Detail, "instrumentId") {
		t.Errorf("the detail does not name the member the server actually sent, so a "+
			"reader is told a name is absent and not that a differently-spelled one is "+
			"there instead: %q", got[0].Detail)
	}
}

// TestCompareLiveReportsTheDocumentationDirection is the other half, driven from
// the same tree.
//
// A name the page declares and the SDK cannot reach is a real SDK defect, and
// reporting it as a live finding would blame the server for a page the SDK
// disagrees with. data.GetStockInstruments is the case: the page declares data
// and pagination_key, []data.StockInstrument carries no field tags for either,
// and the live body does send both -- so this is a disagreement with the page's
// shape, and the row that says so belongs to the documentation.
//
// The way to break this test is to attribute every finding to the live wire.
func TestCompareLiveReportsTheDocumentationDirection(t *testing.T) {
	const symbol = "data.GetStockInstruments"
	m := loadManifest(t)
	f := mustLiveFixture(t, m, symbol)
	sk := mustLiveSkeleton(t, symbol)

	rows := CompareLive(f, symbol, mustType(t, symbol), sk)

	// The row is not comparable on the documented side, so the documented half is
	// expected to be silent. This asserts what the comparison does produce, which
	// is that the SDK cannot reach a name the page declares.
	if len(rows) == 0 {
		t.Fatal("no live finding at all, so nothing is being attributed to either side")
	}
	for _, r := range rows {
		if r.Direction != LiveSDKDirection {
			t.Errorf("a finding is attributed to the documentation on a row the "+
				"comparability rule skips: %+v", r)
		}
		if r.Name == "" {
			continue
		}
		if !strings.Contains(r.Detail, "no json tag") {
			t.Errorf("the name-level finding for %q does not say the SDK has no tag "+
				"for it, which is the claim it makes: %q", r.Name, r.Detail)
		}
	}
	// And the coverage claim: the row is listed as not comparable, so a reader can
	// see its documented side was skipped rather than having found nothing.
	if compareability(f) == nil {
		t.Errorf("%s is comparable, so this test no longer covers a not-comparable row "+
			"and the reason it is chosen no longer holds", symbol)
	}
}

// TestCompareLiveReportsTheContainerDisagreementAsADecodeRejection is the signal
// this harness was built for, and the class the documented comparison is
// structurally blind to.
//
// When the page and the SDK agree, the documented run is green by construction:
// the fixture and the type were made to match. The only way to see that the
// server sent something else is to run the same checks against the server's own
// body. All four captured bodies the SDK could not decode are exactly this.
//
// # The finding is a decode failure, not a shape mismatch, and that is not a detail
//
// checkShape compares f.Checks.TopLevel -- the manifest's record of the page --
// with the SDK type's shape, and never reads the body. So a top-level shape
// mismatch is the *same* observation in both runs whatever body each was handed,
// and it cannot report a server that disagrees with the page. The check that does
// read the body is the decode check, and it is the one that reports these four.
//
// This test pins that, because the obvious implementation -- reporting a
// top-level mismatch when the skeleton's kind differs from the type's -- would be
// a new comparator rather than the existing one, and would be blind in exactly
// the way this finding is not.
//
// The way to break this test is to look for a shape mismatch against the live
// body instead of the decode rejection the comparison actually produces.
func TestCompareLiveReportsTheContainerDisagreementAsADecodeRejection(t *testing.T) {
	type envelope struct {
		Data []string `json:"data"`
	}
	// The page documents an object and the SDK decodes an object: the documented
	// comparison is green and has nothing to say.
	f := liveFixture("test/GET-example", "object", nil, []string{"data"})
	// The server sent an array.
	sk := liveSkeleton(t, `["1"]`)

	rows := CompareLive(f, "test", reflect.TypeOf(envelope{}), sk)

	var failures []LiveDivergence
	for _, r := range rows {
		if r.Kind == DecodeFailure && r.Direction == LiveSDKDirection {
			failures = append(failures, r)
		}
	}
	if len(failures) != 1 {
		t.Fatalf("got %d live decode finding(s), want 1; rows: %v", len(failures), rows)
	}
	d := failures[0].Detail
	if !strings.Contains(d, "envelope") {
		t.Errorf("the detail does not name the type the response failed to decode into: %q", d)
	}
	// The documented half must be silent: the page and the SDK agree, so there is
	// nothing to report against the documentation and reporting there would blame
	// the page for the SDK's shape.
	if hasDirection(rows, failures[0].PairKey(), LiveDocsDirection) {
		t.Errorf("the same rejection is reported against the documentation, where the "+
			"page and the SDK agree: %v", rows)
	}

	// And the shape check must not have produced a live finding, because it cannot
	// see one. A shape mismatch here would mean the comparator had started reading
	// the body, which is the change this test exists to catch.
	for _, r := range rows {
		if r.Kind == TopLevelMismatch {
			t.Errorf("a top-level-shape-mismatch was reported from a live body: %q. "+
				"checkShape reads the manifest rather than the body, so this can only "+
				"mean a second comparator has crept in", r.Detail)
		}
	}
}

// TestCompareLiveReportsTheTopLevelShapeMismatchAsAPageClaim is the other half,
// and it is the honest reading of what the shape check produces.
//
// The SDK decodes an array, the page documents an object, and checkShape reports
// it. So the finding is made twice -- once per run -- and the two details differ,
// because the decode check's account of *why* the body was rejected differs
// between an array and an object. That is the both-detail-changed bucket: not
// agreement, and not two independent findings.
//
// The way to break this test is to treat a same-key-different-detail finding as
// agreement, which is precisely how a live change gets reported as a clean run.
func TestCompareLiveReportsTheTopLevelShapeMismatchAsAPageClaim(t *testing.T) {
	type row struct {
		Symbol string `json:"symbol"`
	}
	// The page documents an object, the SDK decodes an array, and the live body is
	// an object -- so the shape check reports the page's disagreement identically in
	// both runs, while the decode check rejects only the documented fixture.
	//
	// It is data.GetCapitalFlow's real situation, driven against its real fixture so
	// the documented half has a committed instance to read.
	m := loadManifest(t)
	f := mustLiveFixture(t, m, "data.GetCapitalFlow")

	rows := CompareLive(f, "data.GetCapitalFlow", reflect.TypeOf([]row{}),
		liveSkeleton(t, `[{"date":"1"}]`))

	var shapeRows []LiveDivergence
	for _, r := range rows {
		if r.Kind == TopLevelMismatch {
			shapeRows = append(shapeRows, r)
		}
	}
	if len(shapeRows) != 2 {
		t.Fatalf("got %d top-level finding(s), want one per run: %v", len(shapeRows), shapeRows)
	}
	for _, r := range shapeRows {
		if !strings.Contains(r.Detail, recordedTopLevel) {
			t.Errorf("the detail does not name the manifest field the check read: %q; a "+
				"detail naming the body would assert something the check never examined",
				r.Detail)
		}
	}
}

// positionsMissingCurrency is every name trade.GetPositions' page requires except
// currency, so the only finding either run can make about a name is about the one
// this type deliberately cannot reach.
//
// A smaller type would not isolate anything: it would report every name it lacks
// in both runs and the one under test would be lost among them.
type positionsMissingCurrency struct {
	PositionID     string `json:"position_id"`
	Symbol         string `json:"symbol"`
	Quantity       string `json:"quantity"`
	CostPrice      string `json:"cost_price"`
	LastPrice      string `json:"last_price"`
	UnrealizedPL   string `json:"unrealized_profit_loss"`
	InstrumentType string `json:"instrument_type"`
	OptionStrategy string `json:"option_strategy"`
}

// TestCompareLiveDistinguishesNamesTheServerOmittedFromNamesItSent is the
// distinction the whole two-run design exists to draw, stated on both sides at
// once.
//
// The same key -- one name the SDK cannot reach -- means two different things
// depending on whether the server sent it. A caller who never receives the field
// has a zero; a caller who receives a field the SDK drops also has a zero. The
// evidence differs and the finding must say which.
//
// It runs against a real committed fixture and a purpose-built type, so the
// documented half has a committed instance to read and both directions actually
// execute.
//
// The way to break this test is to key on the name alone, which is exactly what
// would make the two collapse into one row.
func TestCompareLiveDistinguishesNamesTheServerOmittedFromNamesItSent(t *testing.T) {
	const omitted = "currency"

	m := loadManifest(t)
	f := mustLiveFixture(t, m, "trade.GetPositions")
	if !slicesContains(f.Checks.RequiredNames, omitted) {
		t.Fatalf("the fixture does not require %q, so this test no longer covers the "+
			"case it was written for: %v", omitted, f.Checks.RequiredNames)
	}
	target := reflect.TypeOf(positionsMissingCurrency{})

	// The server sent the name and the SDK cannot reach it: the SDK is the problem
	// and the server is not.
	sentBody := `{"symbol":"1","currency":"1","position_id":"1","quantity":"1",` +
		`"cost_price":"1","last_price":"1","unrealized_profit_loss":"1",` +
		`"instrument_type":"1","option_strategy":"1"}`
	sent := CompareLive(f, "trade.GetPositions", target, liveSkeleton(t, sentBody))
	if !hasDirection(sent, pairKeyOf("trade.GetPositions", f.ID, MissingRequiredName, omitted), LiveSDKDirection) {
		t.Errorf("a name the server sent and the SDK dropped is not reported against "+
			"the live wire: %v", sent)
	}

	// The server did not send the name. The required-name check asks whether the
	// name is present in the body it was given, and it is not, so the live run
	// reports it -- and that row is live-only, because the documented fixture does
	// carry the name. It is evidence about the server, which is the whole point:
	// a page that marks a name required and a response that omits it is a
	// documentation-versus-server disagreement, and the documented run cannot see
	// it because its own fixture was built to satisfy the page.
	omittedBody := `{"symbol":"1","position_id":"1","quantity":"1",` +
		`"cost_price":"1","last_price":"1","unrealized_profit_loss":"1",` +
		`"instrument_type":"1","option_strategy":"1"}`
	notSent := CompareLive(f, "trade.GetPositions", target, liveSkeleton(t, omittedBody))
	pair := pairKeyOf("trade.GetPositions", f.ID, MissingRequiredName, omitted)

	var sdkOmission, docsMissingTag []LiveDivergence
	for _, r := range notSent {
		if r.PairKey() != pair {
			continue
		}
		switch r.Direction {
		case LiveSDKDirection:
			sdkOmission = append(sdkOmission, r)
		case LiveDocsDirection:
			docsMissingTag = append(docsMissingTag, r)
		}
	}

	// The two rows are the same key with different details, and the difference is
	// the finding: one says the server did not send the name, the other says the
	// SDK has no tag for it. Two accounts of one name, and neither is the other's.
	if len(sdkOmission) != 1 {
		t.Fatalf("got %d live omission finding(s), want 1: %v", len(sdkOmission), notSent)
	}
	if len(docsMissingTag) != 1 {
		t.Fatalf("got %d documented finding(s), want 1: %v", len(docsMissingTag), notSent)
	}
	if !strings.Contains(sdkOmission[0].Detail, "omits") {
		t.Errorf("the live row does not say the body omits the name: %q", sdkOmission[0].Detail)
	}
	if !strings.Contains(docsMissingTag[0].Detail, "no json tag") {
		t.Errorf("the documented row does not say the SDK has no tag for the name: %q",
			docsMissingTag[0].Detail)
	}

	// And the classifier must call that a live change, not agreement.
	cls := ClassifyLive(notSent)
	for _, c := range cls {
		if c.Key != pair {
			continue
		}
		if c.Bucket != LiveBothDetailChanged {
			t.Errorf("the shared name is classified %q; the server's omission and the "+
				"SDK's missing tag are different facts about one name, so reporting it "+
				"as agreement would report a live difference as a clean run", c.Bucket)
		}
		return
	}
	t.Errorf("the shared name was classified under no key at all: %+v", cls)
}

// slicesContains reports whether a slice holds a value, so a test can assert the
// fixture it depends on still declares the name it is about.
func slicesContains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// mustLiveFixture is the committed documentation fixture a live symbol belongs
// to, resolved through the same manifest the comparison uses.
func mustLiveFixture(t *testing.T, m *Manifest, symbol string) Fixture {
	t.Helper()
	entries, err := liveEntries(fixturesFS, testdataDir+"/"+liveManifestFile)
	if err != nil {
		t.Fatalf("read the live manifest: %v", err)
	}
	for _, e := range entries {
		if e.Symbol != symbol {
			continue
		}
		for _, f := range m.Fixtures {
			if f.Fixture == e.Fixture {
				return f
			}
		}
		t.Fatalf("%s: the live manifest names fixture %q, which the documentation "+
			"manifest does not", symbol, e.Fixture)
	}
	t.Fatalf("%s is not in the live manifest", symbol)
	return Fixture{}
}

// mustLiveSkeleton is the committed, expanded body a live symbol was captured
// from, so a test states the real wire rather than a body it made up.
func mustLiveSkeleton(t *testing.T, symbol string) []byte {
	t.Helper()
	skeletons, err := LoadLive("")
	if err != nil {
		t.Fatalf("LoadLive: %v", err)
	}
	s, ok := skeletons[symbol]
	if !ok {
		t.Fatalf("no committed skeleton for %s", symbol)
	}
	return s.Body
}

// TestCompareLiveIsEmptyWhenBothSidesAgree is the negative case, and it is the
// one that makes the positive cases mean anything.
//
// The way to break this test is to report a finding the harness cannot justify.
func TestCompareLiveIsEmptyWhenBothSidesAgree(t *testing.T) {
	type payload struct {
		A string `json:"a"`
		B string `json:"b"`
	}
	f := liveFixture("test/GET-example", "object", []string{"a"}, []string{"a", "b"})
	sk := liveSkeleton(t, `{"a":"1","b":"1"}`)

	if rows := CompareLive(f, "test", reflect.TypeOf(payload{}), sk); len(rows) != 0 {
		t.Errorf("a response the SDK can hold produced %d finding(s): %v", len(rows), rows)
	}
}

// TestCompareLiveAcceptsATypedSkeleton is the reason the probe reduces to typed
// placeholders rather than to kind strings, checked rather than assumed.
//
// jsonKind classifies a decoded value by its Go type, so a leaf spelled as the
// string "number" would be classified a string and a number where a string is
// documented would be reported as agreement. The harness would then be unable to
// detect the wire change it exists to find. A number the SDK holds as a string
// has to be reported, and it is here.
//
// The way to break this test is to compare kinds as strings, which passes for
// every string-shaped disagreement and fails for this one.
func TestCompareLiveAcceptsATypedSkeleton(t *testing.T) {
	type quoted struct {
		Price string `json:"price"`
	}
	f := liveFixture("test/GET-example", "object", []string{"price"}, nil)

	// A number on the wire, a string in the SDK: json.Number is what the probe
	// writes, and jsonKind must classify it "number" rather than "string".
	rows := CompareLive(f, "test", reflect.TypeOf(quoted{}), liveSkeleton(t, `{"price":-1}`))
	want := pairKeyOf("test", f.ID, LeafTypeMismatch, "price")
	if !hasDirection(rows, want, LiveSDKDirection) {
		t.Errorf("a numeric leaf against a string field produced no live leaf finding; "+
			"rows: %v", rows)
	}

	// And the reverse, so the classification is symmetric rather than a rule that
	// happens to fire on numbers.
	type numeric struct {
		Price float64 `json:"price"`
	}
	rows = CompareLive(f, "test", reflect.TypeOf(numeric{}), liveSkeleton(t, `{"price":"1"}`))
	if !hasDirection(rows, want, LiveSDKDirection) {
		t.Errorf("a string leaf against a numeric field produced no live leaf finding; "+
			"rows: %v", rows)
	}
}

// TestCompareLiveReportsADocumentAndLiveDifferenceAsOneRow is the collision
// case, and it is the one a naive implementation gets wrong.
//
// A required name the server omits is a live finding; a required name the page
// requires and the SDK cannot reach is a documented one. When both are true of
// the same name, the name appears once in each run -- and the two rows carry the
// same kind, so a consumer keying on (kind, name) alone would collapse them into
// one and lose which side disagreed.
//
// The way to break this test is to key without the direction.
func TestCompareLiveReportsADocumentAndLiveDifferenceAsOneRow(t *testing.T) {
	const shared = "currency"

	m := loadManifest(t)
	f := mustLiveFixture(t, m, "trade.GetPositions")
	if !slicesContains(f.Checks.RequiredNames, shared) {
		t.Fatalf("the fixture does not require %q, so this test no longer covers the "+
			"case it was written for: %v", shared, f.Checks.RequiredNames)
	}
	// The server did not send the name either, so the live run reports the
	// omission and the documented run reports the missing tag.
	sk := liveSkeleton(t, `{"symbol":"1","position_id":"1","quantity":"1",`+
		`"cost_price":"1","last_price":"1","unrealized_profit_loss":"1",`+
		`"instrument_type":"1","option_strategy":"1"}`)

	rows := CompareLive(f, "trade.GetPositions", reflect.TypeOf(positionsMissingCurrency{}), sk)

	var sdkRows, docsRows []LiveDivergence
	for _, r := range rows {
		if r.Name != shared {
			continue
		}
		switch r.Direction {
		case LiveSDKDirection:
			sdkRows = append(sdkRows, r)
		case LiveDocsDirection:
			docsRows = append(docsRows, r)
		}
	}
	if len(sdkRows) == 0 {
		t.Fatalf("the live omission of a required name is not reported: %v", rows)
	}
	if len(docsRows) == 0 {
		t.Errorf("the SDK's inability to reach a required name is not reported against "+
			"the documentation: %v", rows)
		return
	}
	if sdkRows[0].Key() == docsRows[0].Key() {
		t.Error("the two rows share a key, so one overwrites the other in a map and " +
			"one side of the finding is lost")
	}
	if sdkRows[0].PairKey() != docsRows[0].PairKey() {
		t.Errorf("the two rows do not share a pair key, so the classifier treats them "+
			"as unrelated findings rather than as the two sides of one: %q vs %q",
			sdkRows[0].PairKey(), docsRows[0].PairKey())
	}
	// And the classifier must say so: the same name, two different accounts, one
	// finding that is not agreement.
	cls := ClassifyLive(rows)
	found := false
	for _, c := range cls {
		if c.Name != shared {
			continue
		}
		found = true
		if c.Bucket != LiveBothDetailChanged {
			t.Errorf("the shared name is classified %q; the two runs describe it "+
				"differently -- one says the body omits it, the other says the SDK has "+
				"no tag -- so reporting it as agreement would be reporting a live "+
				"difference as a clean run", c.Bucket)
		}
	}
	if !found {
		t.Errorf("the shared name was classified under no key at all: %+v", cls)
	}
}

// TestCompareLiveSkipsOnlyTheDocumentedDirection is the honesty case, and it is
// asymmetric on purpose.
//
// shapes.go's comparability rule -- a row whose SDK method sends a path the page
// does not document is reported not comparable -- is a rule about the *page*. The
// documented half of this comparison is a claim about the page, so the rule
// applies to it. The live half is not: a skeleton was captured by calling the SDK
// method itself, so the live body is that method's own response whatever any page
// says, and dropping it would discard live evidence over a question about
// documentation that the live evidence does not bear on.
//
// This is not a corner case. data.GetStockInstruments is one of only four captured
// bodies the SDK's own type could not decode, and it is one of three rows the
// comparability rule marks not comparable; gating both directions on the rule
// would drop the strongest finding in the tree on a technicality.
//
// The way to break this test is to apply the rule to both directions, which is
// what the first version of CompareLive did.
func TestCompareLiveSkipsOnlyTheDocumentedDirection(t *testing.T) {
	type payload struct {
		A string `json:"a"`
	}
	f := liveFixture("test/GET-example", "object", nil, []string{"a", "b"})
	f.Documented.Path = "/example/get"
	f.SDKPath = "/broker-fd/assets/get"
	f.SDKPathMatch = PathMatchDiffers

	rows := CompareLive(f, "test", reflect.TypeOf(payload{}), liveSkeleton(t, `{"a":"1","b":"1"}`))

	if len(rows) == 0 {
		t.Fatal("no live finding at all, so the live direction was gated by a rule about " +
			"the page: the server sent a name the SDK cannot reach and that is live evidence")
	}
	for _, r := range rows {
		if r.Direction != LiveSDKDirection {
			t.Errorf("a finding is attributed to the documentation on a row the "+
				"comparability rule skips, so a page that is not this call's contract is "+
				"being compared anyway: %+v", r)
		}
	}
}

// TestLiveDivergenceKeyExcludesDetail is the key discipline, stated as a test.
//
// Detail is prose; it changes when someone improves a sentence. Identity must
// survive that, or every reword of a check is a deleted finding and an inserted
// one. But Direction must be in the key, or the two runs collide -- which is the
// case the previous test pins.
//
// The way to break this test is to include Detail in the key, which is what the
// documented baseline deliberately does and what this set must not copy.
func TestLiveDivergenceKeyExcludesDetail(t *testing.T) {
	a := LiveDivergence{Symbol: "s", Fixture: "f", Direction: LiveSDKDirection,
		Kind: MissingDeclaredName, Name: "b", Detail: "one phrasing"}
	b := a
	b.Detail = "a different phrasing"
	if a.Key() != b.Key() {
		t.Errorf("rephrasing a detail changed the key: %q vs %q; a rewording would "+
			"delete a finding and insert another", a.Key(), b.Key())
	}

	c := a
	c.Direction = LiveDocsDirection
	if a.Key() == c.Key() {
		t.Error("the two directions share a key, so a live finding and a documented " +
			"one would overwrite each other")
	}

	d := a
	d.Name = "c"
	if a.Key() == d.Key() {
		t.Error("two different names share a key")
	}
}

// TestClassifyLiveSplitsTheThreeBuckets is the deliverable: "of N probed, Y
// disagree with the SDK", split by which side disagrees.
//
// The rule is a set difference over keyed rows, and the case that makes it
// non-trivial is a key present in both runs carrying *different* detail. That is
// not agreement and it is not two findings: it is one disagreement the two
// sources describe differently, and losing it is how a live change would be
// reported as agreement.
//
// The way to break this test is to compare whole slices, or to key on the
// full detail string, and each of those puts a same-key-different-detail row in
// the wrong bucket silently.
func TestClassifyLiveSplitsTheThreeBuckets(t *testing.T) {
	fixture := "test/GET-example"
	rows := []LiveDivergence{
		{Symbol: "s", Fixture: fixture, Direction: LiveSDKDirection, Kind: MissingDeclaredName,
			Name: "shared", Detail: "same wording"},
		{Symbol: "s", Fixture: fixture, Direction: LiveDocsDirection, Kind: MissingDeclaredName,
			Name: "shared", Detail: "same wording"},

		{Symbol: "s", Fixture: fixture, Direction: LiveSDKDirection, Kind: LeafTypeMismatch,
			Name: "onlylive", Detail: "the server sent a number where the SDK holds a string"},
		{Symbol: "s", Fixture: fixture, Direction: LiveDocsDirection, Kind: LeafTypeMismatch,
			Name: "onlydocs", Detail: "the page documents a number where the SDK holds a string"},

		// The same disagreement, described differently by the two sources. It is
		// in neither "only" bucket and it is not agreement.
		{Symbol: "s", Fixture: fixture, Direction: LiveSDKDirection, Kind: TopLevelMismatch,
			Name: "moved", Detail: "live sent an array"},
		{Symbol: "s", Fixture: fixture, Direction: LiveDocsDirection, Kind: TopLevelMismatch,
			Name: "moved", Detail: "the page documents an object"},
	}

	got := ClassifyLive(rows)
	if len(got) != 4 {
		t.Fatalf("got %d classification(s), want 4: rows: %+v", len(got), got)
	}
	byKey := map[string]LiveBucket{}
	for _, c := range got {
		byKey[c.Key] = c.Bucket
	}
	// The keys are pair keys, which is what makes the four buckets reachable: a key
	// carrying the direction could never appear on both sides.
	want := map[string]LiveBucket{
		pairKeyOf("s", fixture, MissingDeclaredName, "shared"): LiveBoth,
		pairKeyOf("s", fixture, LeafTypeMismatch, "onlylive"):  LiveOnly,
		pairKeyOf("s", fixture, LeafTypeMismatch, "onlydocs"):  LiveDocsOnly,
		pairKeyOf("s", fixture, TopLevelMismatch, "moved"):     LiveBothDetailChanged,
	}
	if len(byKey) != len(want) {
		t.Fatalf("classified %d key(s), want %d: %+v", len(byKey), len(want), got)
	}
	for key, bucket := range want {
		if byKey[key] != bucket {
			t.Errorf("key %q classified %q, want %q", key, byKey[key], bucket)
		}
	}
}

// TestClassifyLiveIsNotFooledByOrdering is the reason the classification is a
// set difference. CompareBody returns a slice, the two runs will not produce the
// same order, and an implementation that zips or compares element-wise reports a
// different answer for the same evidence.
//
// The way to break this test is to compare the two inputs positionally.
func TestClassifyLiveIsNotFooledByOrdering(t *testing.T) {
	const fixture = "test/GET-example"
	sdk := []LiveDivergence{
		{Symbol: "s", Fixture: fixture, Direction: LiveSDKDirection, Kind: LeafTypeMismatch, Name: "a", Detail: "d"},
		{Symbol: "s", Fixture: fixture, Direction: LiveSDKDirection, Kind: TopLevelMismatch, Name: "b", Detail: "d"},
		{Symbol: "s", Fixture: fixture, Direction: LiveSDKDirection, Kind: MissingRequiredName, Name: "c", Detail: "d"},
	}
	docs := []LiveDivergence{
		{Symbol: "s", Fixture: fixture, Direction: LiveDocsDirection, Kind: TopLevelMismatch, Name: "b", Detail: "d"},
		{Symbol: "s", Fixture: fixture, Direction: LiveDocsDirection, Kind: LeafTypeMismatch, Name: "a", Detail: "d"},
	}

	forward := ClassifyLive(append(append([]LiveDivergence{}, sdk...), docs...))
	reversed := ClassifyLive(append(append([]LiveDivergence{}, docs...), sdk...))
	if !reflect.DeepEqual(forward, reversed) {
		t.Errorf("classification depends on input order.\nforward %+v\nreversed %+v",
			forward, reversed)
	}
}

// TestClassifyLiveSortsItsOutput is what makes a baseline file reviewable as a
// diff: the same evidence must produce the same bytes.
//
// The way to break this test is to emit rows in the order they were classified.
func TestClassifyLiveSortsItsOutput(t *testing.T) {
	const fixture = "test/GET-example"
	rows := []LiveDivergence{
		{Symbol: "s", Fixture: fixture, Direction: LiveSDKDirection, Kind: TopLevelMismatch, Name: "z"},
		{Symbol: "s", Fixture: fixture, Direction: LiveSDKDirection, Kind: LeafTypeMismatch, Name: "a"},
		{Symbol: "s", Fixture: fixture, Direction: LiveSDKDirection, Kind: MissingRequiredName, Name: "m"},
	}
	got := ClassifyLive(rows)
	for i := 1; i < len(got); i++ {
		if got[i-1].Key > got[i].Key {
			t.Fatalf("output is not sorted by key: %q before %q", got[i-1].Key, got[i].Key)
		}
	}
}

// TestLiveBaselineRefusesAnEntryWithoutAReason is the discipline the third set
// inherits from the documented one.
//
// A required reason does not make a reason true, and this file has been bitten by
// a false one repeated across many entries. The guard is kept anyway because a
// reasonless entry is unambiguous, and the loader has to be the thing that
// refuses it rather than a reviewer's memory.
//
// The way to break this test is to accept a blank reason, and the file then
// becomes a suppression list wearing a comment.
func TestLiveBaselineRefusesAnEntryWithoutAReason(t *testing.T) {
	cases := []struct {
		name string
		doc  string
		want string
	}{
		{
			name: "no reason",
			doc:  `{"policy":"p","entries":[{"symbol":"s","fixture":"f","direction":"live-sdk","kind":"top-level-shape-mismatch","detail":"d","unblock":"u"}]}`,
			want: "reason",
		},
		{
			name: "a whitespace reason",
			doc:  `{"policy":"p","entries":[{"symbol":"s","fixture":"f","direction":"live-sdk","kind":"top-level-shape-mismatch","detail":"d","reason":"  ","unblock":"u"}]}`,
			want: "reason",
		},
		{
			name: "no unblock",
			doc:  `{"policy":"p","entries":[{"symbol":"s","fixture":"f","direction":"live-sdk","kind":"top-level-shape-mismatch","detail":"d","reason":"r"}]}`,
			want: "unblock",
		},
		{
			name: "no detail",
			doc:  `{"policy":"p","entries":[{"symbol":"s","fixture":"f","direction":"live-sdk","kind":"top-level-shape-mismatch","reason":"r","unblock":"u"}]}`,
			want: "detail",
		},
		{
			name: "an unknown direction",
			doc:  `{"policy":"p","entries":[{"symbol":"s","fixture":"f","direction":"sideways","kind":"top-level-shape-mismatch","detail":"d","reason":"r","unblock":"u"}]}`,
			want: "direction",
		},
		{
			name: "an unknown kind",
			doc:  `{"policy":"p","entries":[{"symbol":"s","fixture":"f","direction":"live-sdk","kind":"invented","detail":"d","reason":"r","unblock":"u"}]}`,
			want: "kind",
		},
		{
			name: "no policy",
			doc:  `{"entries":[{"symbol":"s","fixture":"f","direction":"live-sdk","kind":"top-level-shape-mismatch","detail":"d","reason":"r","unblock":"u"}]}`,
			want: "policy",
		},
		{
			name: "no entries",
			doc:  `{"policy":"p","entries":[]}`,
			want: "no entries",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseLiveBaseline([]byte(tc.doc))
			if err == nil {
				t.Fatal("the baseline was accepted, so the guard does not guard")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error does not mention %q: %v", tc.want, err)
			}
		})
	}
}

// TestLiveBaselineRefusesADuplicateKey stops a finding being recorded twice,
// which would read as two independent observations of one fact.
//
// The way to break this test is to accept the second entry.
func TestLiveBaselineRefusesADuplicateKey(t *testing.T) {
	const entry = `{"symbol":"s","fixture":"f","direction":"live-sdk","kind":"top-level-shape-mismatch","detail":"d","reason":"r","unblock":"u"}`
	doc := `{"policy":"p","entries":[` + entry + `,` + entry + `]}`
	if _, err := ParseLiveBaseline([]byte(doc)); err == nil {
		t.Fatal("a duplicated key was accepted")
	} else if !strings.Contains(err.Error(), "twice") {
		t.Errorf("error does not name the duplication: %v", err)
	}
}

// TestParseLiveBaselineIsTheInverseOfRendering is the reviewability property: the
// committed file is what the harness would write, so a diff is a change in
// evidence rather than in formatting.
//
// The way to break this test is a renderer that drops, reorders or rephrases a
// field, which would make every regeneration a diff.
func TestParseLiveBaselineIsTheInverseOfRendering(t *testing.T) {
	in := &LiveBaseline{
		Policy: "the gate passes when the observed set exactly equals this set",
		Notes:  map[string]string{"bucketCounts": "live-only is the new signal"},
		Entries: []LiveEntry{{
			LiveDivergence: LiveDivergence{
				Symbol: "data.GetStockInstruments", Fixture: "trading/GET-x.json",
				Direction: LiveSDKDirection, Kind: TopLevelMismatch,
				Detail: "the live response top level is object, []data.StockInstrument decodes it as array",
			},
			Reason:  "the SDK decodes an array and the server sent an object",
			Unblock: "a US sandbox credential, or a production account, would settle it",
		}},
	}
	first, err := in.Marshal()
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	parsed, err := ParseLiveBaseline(first)
	if err != nil {
		t.Fatalf("ParseLiveBaseline: %v", err)
	}
	second, err := parsed.Marshal()
	if err != nil {
		t.Fatalf("Marshal(parsed): %v", err)
	}
	if string(first) != string(second) {
		t.Errorf("render -> parse -> render is not stable.\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

// TestLoadLiveReadsTheCommittedTree is the entry point the whole live
// comparison runs through, checked against the tree as committed.
//
// The way to break this test is to read a subset, or to fall back to generating
// skeletons, which would make the comparison check the harness's own output.
func TestLoadLiveReadsTheCommittedTree(t *testing.T) {
	got, err := LoadLive("")
	if err != nil {
		t.Fatalf("LoadLive: %v", err)
	}
	manifest := readLiveManifest(t)
	if len(got) != len(manifest.Entries) {
		t.Fatalf("loaded %d skeleton(s), the manifest records %d", len(got), len(manifest.Entries))
	}
	for _, e := range manifest.Entries {
		s, ok := got[e.Symbol]
		if !ok {
			t.Errorf("no skeleton loaded for %s", e.Symbol)
			continue
		}
		if s.Fixture != e.Fixture {
			t.Errorf("%s: fixture = %q, manifest says %q", e.Symbol, s.Fixture, e.Fixture)
		}
		if len(s.Body) == 0 {
			t.Errorf("%s: loaded skeleton is empty", e.Symbol)
		}
	}
}

// TestReadLiveSkeletonRefusesARowThatCapturedNothing is the latent crash, named.
//
// A row whose capture produced no body names no file, and the empty string is a
// name this package's own types permit. CompareAllLive built the path from it
// unconditionally, so such a row read the tree root: an error reading "is a
// directory", once per row, naming no endpoint. The committed manifest has no such
// row, so nothing reached it - which is why the assertion is on the sentence and
// not on the error existing, since the unguarded code also returned an error.
//
// The second half is the guard not having become a refusal of everything: the
// comparison is useless if it reads nothing, so a real row is read here too.
func TestReadLiveSkeletonRefusesARowThatCapturedNothing(t *testing.T) {
	_, err := readLiveSkeleton(liveEntry{Symbol: "data.GetNothing", Fixture: "market-data-stock/GET-x.json"})
	if err == nil {
		t.Fatal("readLiveSkeleton read a row that names no skeleton")
	}
	for _, want := range []string{"data.GetNothing", "no skeleton", "nothing to compare"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the refusal does not say %q: %v", want, err)
		}
	}
	for _, never := range []string{"is a directory", testdataDir} {
		if strings.Contains(err.Error(), never) {
			t.Errorf("the refusal carries %q, which is what the unguarded read produced: %v",
				never, err)
		}
	}

	manifest := readLiveManifest(t)
	if len(manifest.Entries) == 0 {
		t.Fatal("the live manifest records no entries")
	}
	raw, err := readLiveSkeleton(manifest.Entries[0])
	if err != nil {
		t.Fatalf("readLiveSkeleton on a captured row: %v", err)
	}
	if len(raw) == 0 {
		t.Error("readLiveSkeleton returned an empty body for a captured row")
	}
}

// TestLoadLiveRefusesAMissingTree is the refusal half: a harness that reported
// zero skeletons rather than an error would compare nothing and pass.
//
// The way to break this test is to return an empty map and a nil error.
func TestLoadLiveRefusesAMissingTree(t *testing.T) {
	if _, err := LoadLive(t.TempDir()); err == nil {
		t.Fatal("a directory holding no live tree was accepted, so an empty " +
			"comparison would read as a clean run")
	}
}

// TestLoadLiveExpandsToBodiesCompareBodyCanRead is the whole chain in one
// assertion, and it is the assumption the design rests on.
//
// A skeleton is only usable because Task 1's reducer emits typed placeholders:
// jsonKind classifies a decoded value by its Go type, so a number written as -1
// is classified "number" and compared as such. If the committed form were kind
// strings, every numeric leaf would be demanded of a Go string and the harness
// would report agreement for a wire change it exists to find.
//
// The way to break this test is to expand into a tree whose leaves are strings
// spelling kinds, which parses cleanly and would pass a weaker check.
func TestLoadLiveExpandsToBodiesCompareBodyCanRead(t *testing.T) {
	skeletons, err := LoadLive("")
	if err != nil {
		t.Fatalf("LoadLive: %v", err)
	}
	f, ok := skeletons["data.GetStockInstruments"]
	if !ok {
		t.Skip("data.GetStockInstruments is not in the committed live tree")
	}
	if !json.Valid(f.Body) {
		t.Fatalf("the expanded body is not valid JSON: %s", f.Body)
	}
	dec := json.NewDecoder(strings.NewReader(string(f.Body)))
	dec.UseNumber()
	var tree any
	if err := dec.Decode(&tree); err != nil {
		t.Fatalf("decode expanded body: %v", err)
	}
	// Every leaf must be classified by jsonKind as the kind it stands for, and a
	// kind spelled as a string is the failure this guards.
	assertSkeletonLeavesAreTyped(t, tree, "$", f.Symbol)
}

// assertSkeletonLeavesAreTyped walks an expanded body and fails on any leaf
// whose JSON kind is not the one the reduction's placeholder carries.
func assertSkeletonLeavesAreTyped(t *testing.T, v any, path, symbol string) {
	t.Helper()
	switch node := v.(type) {
	case map[string]any:
		for name, member := range node {
			assertSkeletonLeavesAreTyped(t, member, path+"."+name, symbol)
		}
	case []any:
		for i, element := range node {
			assertSkeletonLeavesAreTyped(t, element, path+"["+itoa(i)+"]", symbol)
		}
	default:
		kind := jsonKind(v)
		switch kind {
		case "string", "number", "boolean", "null", "object", "array":
		default:
			t.Errorf("%s: leaf at %s is a %s, which jsonKind cannot classify", symbol, path, kind)
		}
		if _, isString := v.(string); isString && v == "number" {
			t.Errorf("%s: the leaf at %s is the kind as a string, so jsonKind would "+
				"classify it a string and the comparison would read it as agreement", symbol, path)
		}
	}
}

// The report is the sentence the status documents quote: "of N probed, X disagree
// with their documentation and Y disagree with the SDK". It is the deliverable
// phase 4 is gated on, because the gate is on Y: a run in which nothing disagreed
// with the SDK closes a question, and a run in which something did opens a defect
// item per disagreeing endpoint.
//
// It lives in a test file because it is a report, not a library surface. Nothing
// calls it and no non-test code depends on it, so exporting it would put a
// sentence format into the package's public contract for no reason.

// liveReport holds the two counts, derived rather than written down.
type liveReport struct {
	// Probed is the number of captured responses the run compared.
	Probed int
	// Docs is how many findings were observed against the documentation, so the
	// SDK and the page disagree.
	Docs int
	// SDK is how many were observed against the live wire, so the SDK and the
	// server disagree.
	SDK int
	// Both is how many of the findings each direction contributed, which is what
	// keeps Docs and SDK from being added into a total.
	Both int
}

// summariseLive derives the report from a run.
//
// Every count comes from the classifications or from the raw rows, never from a
// literal, so a report and the set it describes cannot drift apart. The two
// directions are read from different places on purpose: the classification bucket
// says which side *carried* a finding and the raw row's Direction says which
// *body* it was observed against, and a classifier that quietly changed what a
// bucket means would move one and not the other.
func summariseLive(run LiveRun) liveReport {
	rep := liveReport{Probed: run.Coverage.Probed}
	for _, c := range run.Classifications {
		switch c.Bucket {
		case LiveDocsOnly:
			rep.Docs++
		case LiveOnly:
			rep.SDK++
		case LiveBoth, LiveBothDetailChanged:
			// One finding, observed on both sides. It is counted in both
			// directions because it is a disagreement with each, and separately
			// because Docs+SDK would otherwise over-count the set by one.
			rep.Docs++
			rep.SDK++
			rep.Both++
		default:
			// BucketCounts already requires every bucket to be accounted for; a
			// fifth bucket here would be a new meaning rather than a new count.
			panic("conformance: live report does not know bucket " + string(c.Bucket))
		}
	}
	return rep
}

// Line renders the one sentence the status documents quote.
func (r liveReport) Line() string {
	return "of " + itoa(r.Probed) + " probed, " + itoa(r.Docs) +
		" disagree with their documentation and " + itoa(r.SDK) +
		" disagree with the SDK"
}

// reportLineRE reads the two counts back out of a rendered line, so a renderer
// holding a literal rather than the derived count fails here instead of in a
// document.
var reportLineRE = regexp.MustCompile(`^of (\d+) probed, (\d+) disagree with their documentation and (\d+) disagree with the SDK$`)

// TestLiveReportCountsBothDirections is the report gate.
//
// It asserts three things. The two directions are distinguished, so a reader
// cannot collapse "disagrees with the page" into "disagrees with the server". The
// numbers the line prints are the counts the run found, checked by reading them
// back out of the rendered text, so a constant in the renderer fails here rather
// than in a status document. And the run has a non-zero count in each direction,
// which is the claim phase 4 rests on: Y is not zero, so a defect item is opened
// per disagreeing endpoint rather than the question being declared closed.
//
// The way to break this test is a renderer that formats literals, or a summary
// that reads one direction twice, and both are the arithmetic errors a report
// makes when a person writes the sentence by hand.
func TestLiveReportCountsBothDirections(t *testing.T) {
	run := CompareAllLive()
	for _, err := range run.Errs {
		t.Errorf("the run could not compare an endpoint, so the report would be "+
			"a report on a subset that does not say so: %v", err)
	}

	rep := summariseLive(run)

	// The two counts are computed from independent sources: the buckets say which
	// side carried each finding, the raw rows say which body it was observed
	// against. Agreement between them is the check that the classifier still means
	// what the report claims it means.
	byDirection := map[LiveDirection]int{}
	for _, r := range run.Rows {
		byDirection[r.Direction]++
	}
	if byDirection[LiveDocsDirection] != rep.Docs {
		t.Errorf("%d row(s) carry the documentation direction but the report counts "+
			"%d; the two are read from different places and must agree",
			byDirection[LiveDocsDirection], rep.Docs)
	}
	if byDirection[LiveSDKDirection] != rep.SDK {
		t.Errorf("%d row(s) carry the live direction but the report counts %d; the "+
			"two are read from different places and must agree",
			byDirection[LiveSDKDirection], rep.SDK)
	}

	// Everything found is in exactly one bucket, and the overlap between the two
	// counts is exactly the findings both directions contributed. Without this a
	// finding could be in neither count and the report would still print.
	if sum := rep.Docs + rep.SDK - rep.Both; sum != len(run.Classifications) {
		t.Errorf("the two directions account for %d finding(s) between them, the "+
			"run classified %d; a finding in neither count is missing from the "+
			"report", sum, len(run.Classifications))
	}
	if rep.Both == rep.Docs || rep.Both == rep.SDK {
		t.Errorf("every finding is in both directions (%d docs, %d SDK, %d both), so "+
			"the report cannot be distinguishing them", rep.Docs, rep.SDK, rep.Both)
	}

	// The rendered numbers must be the derived ones. Reading them back is what
	// makes this a check rather than a restatement.
	m := reportLineRE.FindStringSubmatch(rep.Line())
	if m == nil {
		t.Fatalf("the report line does not have the expected shape: %q", rep.Line())
	}
	for _, c := range []struct {
		what string
		got  string
		want int
	}{
		{"probed", m[1], rep.Probed},
		{"documentation", m[2], rep.Docs},
		{"SDK", m[3], rep.SDK},
	} {
		if c.got != itoa(c.want) {
			t.Errorf("the report prints %s = %s; the run found %d", c.what, c.got, c.want)
		}
	}

	// The claim the phase-4 gate rests on, in both directions. Y is not zero, so
	// the run did not close the question and each disagreeing endpoint is opened
	// as a defect item; X is not zero either, so the documented half is not a
	// clean row either.
	if rep.SDK == 0 {
		t.Error("no finding disagreed with the SDK, so the report would claim the " +
			"gate closed the question; the committed tree disagrees with that")
	}
	if rep.Docs == 0 {
		t.Error("no finding disagreed with the documentation, so the report would " +
			"claim the documented comparison found nothing to say about the SDK")
	}
	if rep.Probed != run.Coverage.Compared {
		t.Errorf("the report prints %d probed and the run compared %d; the line is "+
			"what a reader would take as the denominator", rep.Probed, run.Coverage.Compared)
	}
	t.Log(rep.Line())
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
