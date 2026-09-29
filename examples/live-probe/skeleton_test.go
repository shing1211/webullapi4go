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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestSkeletonify(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want any
	}{
		{"object", `{"a":"1"}`, map[string]any{"a": "1"}},
		// A number must stay a number: Money and QuoteTime accept both forms, so
		// conflating them is the one mistake this function exists to prevent. The
		// kind is carried by the json.Number itself, because the consumer
		// classifies a value by its Go type and a leaf spelled "number" would be
		// read as a string. The literal is the shared placeholder, not 385.6: the
		// kind is identical either way and the server's figure is not carried.
		{"number is not string", `{"close":385.6}`, map[string]any{"close": numberPlaceholder}},
		// A present-but-null field is a name the SDK may not carry, so it must
		// survive as a null rather than vanish. A nil is what the consumer
		// classifies as "null"; the string "null" would be classified "string".
		{"null is recorded", `{"outstanding":null}`, map[string]any{"outstanding": nil}},
		{"array of objects", `[{"a":1}]`, []any{map[string]any{"a": numberPlaceholder}}},
		{"empty array is not a wrong kind", `[]`, []any{}},
		{"bool", `{"ok":true}`, map[string]any{"ok": true}},
		{"nested", `{"d":{"e":[]}}`, map[string]any{"d": map[string]any{"e": []any{}}}},
		// This one case proves two things. Without UseNumber, 1e400 is valid
		// JSON that a float64 decoder cannot represent, so the body fails to
		// decode and the probe records an unparseable live body instead of a
		// skeleton: a body silently dropped from the evidence set, which is the
		// one failure this harness cannot detect about itself. With UseNumber and
		// a synthetic leaf, the same body decodes and reduces to a number, which
		// is why the expected value is the placeholder and not "1e400". The
		// literal the decoder preserved is discarded by the reduction; the case
		// keeps the input precisely because it is the one that proves the
		// preservation was necessary first.
		{"number a float64 cannot hold", `{"big":1e400}`, map[string]any{"big": numberPlaceholder}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Skeletonify([]byte(tc.in))
			if err != nil {
				t.Fatalf("err = %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestSkeletonifyRejectsInvalidJSON(t *testing.T) {
	if _, err := Skeletonify([]byte(`{"a":`)); err == nil {
		t.Error("invalid JSON returned no error")
	}
}

// A body of one JSON value is a whole body. A second value would be reduced
// away without a word, so the reduction has to refuse it; surrounding
// whitespace is not a second value and must still be accepted.
func TestSkeletonifyRejectsTrailingData(t *testing.T) {
	if _, err := Skeletonify([]byte(`{"a":1} {"b":2}`)); err == nil {
		t.Error("two JSON values returned no error")
	}
	if _, err := Skeletonify([]byte("{\"a\":1}\n")); err != nil {
		t.Errorf("surrounding whitespace: err = %v", err)
	}
}

// The brief fixes the reduction of a body's members but not of a body that is
// itself a scalar, which is the shape several Webull endpoints answer with.
func TestSkeletonifyTopLevelScalar(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want any
	}{
		{"null", `null`, nil},
		// Three inputs, one output. A scalar body is a number as surely as a
		// member is, so it reduces to the same placeholder and carries no
		// reading either.
		{"integer stays a number", `42`, numberPlaceholder},
		{"decimal", `385.6`, numberPlaceholder},
		{"quoted number is a string", `"385.6"`, "1"},
		{"bool", `true`, true},
		// The boolean placeholder does not echo its input. Emitting one
		// unconditional value is what makes a `false` leaf in a reduced tree
		// proof that the tree was not reduced, so the case is pinned here rather
		// than assumed.
		{"false is the same placeholder as true", `false`, true},
		// A scalar body of a number no float64 could hold reduces like any
		// other, which is the top-level half of the UseNumber argument. The
		// overflow is the input that needs UseNumber; the placeholder is the
		// output that proves the reduction threw the literal away.
		{"unrepresentable number", `1e400`, numberPlaceholder},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Skeletonify([]byte(tc.in))
			if err != nil {
				t.Fatalf("err = %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}

// An array reduces element-wise, and each element keeps its own kind. The
// brief's array case has a single object element, so a heterogeneous array is
// unpinned by it, and a heterogeneous array is where a reduction that collapsed
// kinds would show.
func TestSkeletonifyMixedKindArray(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want any
	}{
		{
			"every element kind is kept",
			`[1,"1",true,null,[],{}]`,
			[]any{numberPlaceholder, "1", true, nil, []any{}, map[string]any{}},
		},
		{
			"a null element is an element",
			`[null,null]`,
			[]any{nil, nil},
		},
		{
			"nested arrays reduce element-wise",
			`[[1,["2"]],[]]`,
			[]any{[]any{numberPlaceholder, []any{"1"}}, []any{}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Skeletonify([]byte(tc.in))
			if err != nil {
				t.Fatalf("err = %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}

// The property the whole harness rests on: a reduced tree can carry a name the
// server sent and cannot carry a value it sent, so a consumer may write one to
// a committed fixture without a leak. The walk is the assertion, and the corpus
// is chosen for the inputs that break a reduction rather than for tidiness.
func TestSkeletonifyNoLiveValueSurvives(t *testing.T) {
	corpus := []struct {
		name string
		in   string
	}{
		// A 40-digit integer is exact as a json.Number and already rounded as a
		// float64, so it is where a widened decode would show first. It is also
		// the most account-number-shaped literal available, which is why it is
		// in the corpus: it is the number a reader of a committed fixture would
		// have to notice.
		{"forty digit integer", `{"id":1234567890123456789012345678901234567890,"acct":"9110101000000000001"}`},
		// 1e400 overflows float64 to +Inf and is a decode error without
		// UseNumber, so it is the case that makes the option load-bearing. With
		// it, both numbers decode and neither literal may survive.
		{"overflowing number", `{"big":1e400,"small":1e-400}`},
		// Escapes are where a reduction that decoded and re-encoded a string
		// could change its bytes, and where a live value hides most cheaply.
		{"escaped quotes and backslashes", `{"q":"he said \"acct-123\" \\ done\n\t\r\u0000"}`},
		{"large unicode string", `{"s":"` + strings.Repeat("香港交易所-Ünïcödé-\U0001f9ee-", 64) + `","t":"é"}`},
		// A body shaped like the ones this harness exists to adjudicate: every
		// kind, at depth, with a false boolean and a null that must not vanish.
		{
			"webull shaped body",
			`{"symbol":"AAPL","close":385.6,"closeTime":1756000000000,"ratio":null,` +
				`"active":false,"lots":[{"symbol":"AAPL","filledQty":100,"price":"1.0"},{},[]],` +
				`"empty":[],"emptyObj":{}}`,
		},
		{"top level scalar string", `"acct-123"`},
		{"top level scalar number", `1e400`},
	}
	for _, tc := range corpus {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Skeletonify([]byte(tc.in))
			if err != nil {
				t.Fatalf("err = %v", err)
			}
			assertNoLiveValue(t, got, "$")
		})
	}

	// The other half of the property, in the same test because it is the same
	// claim: a member name is live data the harness needs, and stripping it
	// would empty the evidence set just as silently as a surviving value would
	// pollute it.
	t.Run("member names are not values", func(t *testing.T) {
		got, err := Skeletonify([]byte(`{"symbol":"AAPL","close":385.6,"ratio":null,"lots":[]}`))
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		body, ok := got.(map[string]any)
		if !ok {
			t.Fatalf("got %T, want a map", got)
		}
		for _, name := range []string{"symbol", "close", "ratio", "lots"} {
			if _, present := body[name]; !present {
				t.Errorf("member %q did not survive the reduction", name)
			}
		}
		if len(body) != 4 {
			t.Errorf("got %d members, want 4: a name was invented or lost", len(body))
		}
	})
}

// liveValueSink is where assertNoLiveValue reports. *testing.T satisfies it, and
// so does recordingSink below, which is what lets the key rule be exercised
// against a reduced tree that does carry a reading as a key: such a tree cannot
// be committed in order to be walked, so the failing case has to be driven here.
type liveValueSink interface {
	Helper()
	Errorf(format string, args ...any)
}

// recordingSink collects what assertNoLiveValue reports rather than failing the
// test driving it.
type recordingSink struct{ messages []string }

func (r *recordingSink) Helper() {}

func (r *recordingSink) Errorf(format string, args ...any) {
	r.messages = append(r.messages, fmt.Sprintf(format, args...))
}

// keyCarriesAReading reports whether an object member name could itself be a
// value the server sent.
//
// The reduction keeps every member name verbatim and replaces every value, so a
// name is the one place a reading can survive it. Two forms of key are refused,
// and the rule is deliberately the JSON number grammar rather than a list of
// readings: a key that is a well-formed JSON number is an account id, an
// instrument id, a millisecond timestamp, a price, a size or a count when it is
// one of those things, and no Webull member name is; a key that is one of the
// four placeholders, or the empty string, is a key a reader cannot tell from a
// value.
//
// The rule cannot see a key that is a string reading, and that is stated rather
// than implied: a JSON member name and a JSON string value are the same token in
// the same grammar, so nothing in the document or the reduced tree separates
// them. What the tree therefore carries is no numeric reading and no
// placeholder, and a string-valued key is not excluded - see the identical rule
// and its limits in conformance/fixtures_test.go, which holds the committed
// bytes to it as well.
func keyCarriesAReading(name string) bool {
	if name == "" {
		return true
	}
	switch name {
	case stringPlaceholder, string(numberPlaceholder), "true", "null":
		return true
	}
	_, err := json.Marshal(json.Number(name))
	return err == nil
}

// assertNoLiveValue walks a reduced tree and fails on any key or leaf that could
// carry something the server sent.
//
// A string leaf must be the string placeholder, so a symbol or an account id
// surviving is a failure and the path names the member it survived in. A key is
// checked before the walk descends, because a member name is kept verbatim and a
// name that is a reading is a leak the leaf checks cannot see. A float64 or an
// int is a failure too, and of a different kind: neither is a type this reduction
// emits, so one appearing means the value was widened on the way through, which
// is what a decode without UseNumber does. A bool leaf cannot carry a value,
// because the reduction emits one unconditional placeholder, so `false` is a
// failure: it proves the tree did not come from Skeletonify.
//
// A number leaf must be the number placeholder, and this is the assertion that
// used to be missing. The kind is carried by the type, so a server literal buys
// no accuracy, and a reduced tree is about to be written to a committed fixture
// for every endpoint: a surviving number is a committed account number, which
// is a value in git history rather than a mistake that can be deleted. So every
// number leaf is compared against the one literal the reduction emits, and the
// message names both the path and the literal that got through.
func assertNoLiveValue(t liveValueSink, v any, path string) {
	t.Helper()
	switch leaf := v.(type) {
	case map[string]any:
		for name, member := range leaf {
			if keyCarriesAReading(name) {
				t.Errorf("%s holds the member name %q: no Webull member name is a JSON "+
					"number or one of the four placeholders, so a key in either form is a "+
					"reading that survived the reduction", path, name)
			}
			assertNoLiveValue(t, member, path+"."+name)
		}
	case []any:
		for i, element := range leaf {
			assertNoLiveValue(t, element, fmt.Sprintf("%s[%d]", path, i))
		}
	case string:
		if leaf != stringPlaceholder {
			t.Errorf("%s = %q: a live string survived the reduction", path, leaf)
		}
	case json.Number:
		if leaf != numberPlaceholder {
			t.Errorf("%s = %s: a live number survived the reduction", path, leaf)
		}
	case bool:
		if !leaf {
			t.Errorf("%s = false: the reduction emits %v for every boolean, so this leaf was not reduced", path, boolPlaceholder)
		}
	case nil:
	case float64, int:
		t.Errorf("%s is a %T: a widened value, which no reduction emits", path, leaf)
	default:
		t.Errorf("%s is a %T: no reduction emits that type", path, leaf)
	}
}

// The key half of the invariant, on its own. The corpus in TestSkeletonify proves
// the rule holds for the reduced trees real bodies produce; this table is what
// proves the rule fires, because it is the only place a key that is a reading can
// be presented at all.
//
// The last row pins the stated limit. A ticker key is a string reading, and a
// string reading is not distinguishable from a member name, so it passes and the
// case fails if the rule ever quietly narrows its own claim.
func TestSkeletonifyRejectsAReadingCarriedAsAKey(t *testing.T) {
	cases := []struct {
		name        string
		in          string
		wantReports int
		because     string
	}{
		{name: "member names", in: `{"symbol":"AAPL","orders":[{"symbol":"AAPL"}]}`,
			wantReports: 0, because: "a name is evidence and must survive"},
		{name: "a camelCase member name", in: `{"instrumentId":"9132750001"}`,
			wantReports: 0, because: "the rule is not a snake_case rule"},
		{name: "an account id as a key", in: `{"9110101000000000001":{"symbol":"AAPL"}}`,
			wantReports: 1, because: "an account id as a key is committed verbatim"},
		{name: "an instrument id as a key", in: `{"9132750001":{"symbol":"AAPL"}}`,
			wantReports: 1, because: "a ten digit id is the JSON number grammar"},
		{name: "a timestamp as a key", in: `{"1756000000000":{"symbol":"AAPL"}}`,
			wantReports: 1, because: "a millisecond epoch is a number like any other"},
		{name: "a price as a key", in: `{"385.6":{"symbol":"AAPL"}}`,
			wantReports: 1, because: "the fraction is still that grammar"},
		{name: "the string placeholder as a key", in: `{"1":{"symbol":"AAPL"}}`,
			wantReports: 1, because: "a key indistinguishable from a leaf is the ambiguity " +
				"the representation exists to remove"},
		{name: "the number placeholder as a key", in: `{"-1":{"symbol":"AAPL"}}`,
			wantReports: 1, because: "caught by the number rule too, and both agreeing is the point"},
		{name: "an empty key", in: `{"":{"symbol":"AAPL"}}`,
			wantReports: 1, because: "an empty key is indistinguishable from an empty string reading"},
		{name: "a reading as a key at depth", in: `{"data":{"9110101000000000001":{"symbol":"AAPL"}}}`,
			wantReports: 1, because: "the walk descends, so a nested map cannot hide a key"},
		{name: "a ticker as a key", in: `{"AAPL":{"symbol":"AAPL"},"MSFT":{"symbol":"AAPL"}}`,
			wantReports: 0, because: "the stated limit, pinned: a JSON member name and a " +
				"JSON string value are the same token, so a string reading as a key is " +
				"not detectable and keyCarriesAReading says so"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// The reduction itself preserves every key, which is what makes the
			// planted bodies above a fair test of the walk rather than of a fixture.
			got, err := Skeletonify([]byte(tc.in))
			if err != nil {
				t.Fatalf("err = %v", err)
			}
			if _, isObject := got.(map[string]any); !isObject {
				t.Fatalf("the corpus reduced to %T, want an object", got)
			}
			rec := &recordingSink{}
			assertNoLiveValue(rec, got, "$")
			if n := len(rec.messages); n != tc.wantReports {
				t.Errorf("assertNoLiveValue reported %d time(s), want %d: %v; %s",
					n, tc.wantReports, rec.messages, tc.because)
			}
		})
	}

	// The rule is not vacuous, and TestSkeletonify is the other half of that
	// evidence: it walks every reduced tree a real body produces.
	if keyCarriesAReading("symbol") || keyCarriesAReading("orders") {
		t.Error("the key rule reports ordinary member names, so it would fail every " +
			"reduction rather than describe a leak")
	}
}

// A reduced tree is written to a committed fixture, so it has to marshal, and
// what comes out has to be the placeholders rather than the inputs. The walk
// above proves the tree in memory; this proves the bytes that reach the
// repository, which is where a leak would actually land.
func TestSkeletonifyMarshalsToPlaceholders(t *testing.T) {
	const in = `{"symbol":"AAPL","close":385.6,"closeTime":1756000000000,"acct":` +
		`1234567890123456789012345678901234567890,"ratio":null,"active":false,` +
		`"big":1e400,"lots":[{"filledQty":100}],"empty":[]}`

	got, err := Skeletonify([]byte(in))
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	// json.Number is validated on the way out, so a placeholder that is not
	// valid number syntax would fail every fixture write the harness performs.
	encoded, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("a reduced tree does not marshal: err = %v", err)
	}

	// A round trip through the wire shape, because a committed fixture is read
	// back by a decoder and must still be one jsonKind can classify.
	var reread any
	dec := json.NewDecoder(bytes.NewReader(encoded))
	dec.UseNumber()
	if err := dec.Decode(&reread); err != nil {
		t.Fatalf("a marshalled tree does not decode: err = %v", err)
	}
	assertNoLiveValue(t, reread, "$")

	// Every literal in the output is one of the four the reduction emits, so a
	// server figure appearing in a committed file is visible in a diff rather
	// than only to this walk.
	for _, literal := range []string{"-1", `"1"`, "true", "null"} {
		if !strings.Contains(string(encoded), literal) {
			t.Errorf("marshalled tree %s does not contain the %s placeholder", encoded, literal)
		}
	}
	for _, leaked := range []string{"AAPL", "385.6", "1756000000000", "1e400", "100",
		"1234567890123456789012345678901234567890"} {
		if strings.Contains(string(encoded), leaked) {
			t.Errorf("marshalled tree %s contains the server literal %q", encoded, leaked)
		}
	}
}

func TestSkeletonKind(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want string
	}{
		{"object", map[string]any{"a": "1"}, "object"},
		{"array", []any{numberPlaceholder}, "array"},
		{"string placeholder", "1", "string"},
		{"number placeholder", numberPlaceholder, "number"},
		{"boolean", true, "boolean"},
		{"empty array is still an array", []any{}, "array"},
		{"nil is null", nil, "null"},
		// The consumer classifies a number by its Go type and never reads the
		// literal, which is the whole reason the reduction can throw the literal
		// away without losing anything: 1e400 and the placeholder are the same
		// answer. Pinned so a future change to jsonKind is caught here rather
		// than in a fixture 193 endpoints deep.
		{"any number is a number", json.Number("1e400"), "number"},
		// The guard is closed. An int is what the brief's stale interface line
		// promised and no reduction emits; a float64 is what a decode without
		// UseNumber would have left; a struct is a type the consumer never
		// produces. All three skipped the reduction, and "" says so instead of
		// naming a kind.
		{"int has no kind", 1, ""},
		{"float64 has no kind", 1.5, ""},
		{"struct has no kind", struct{}{}, ""},
		// The one case the guard cannot catch here: the container type really is
		// "array", and the live values are its members, which the walk in
		// TestSkeletonifyNoLiveValueSurvives is what rejects.
		{"unreduced array is still an array", []any{"acct-123", 1.5}, "array"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SkeletonKind(tc.in); got != tc.want {
				t.Errorf("SkeletonKind(%#v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
