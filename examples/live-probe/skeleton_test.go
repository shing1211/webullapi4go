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
		// read as a string.
		{"number is not string", `{"close":385.6}`, map[string]any{"close": json.Number("385.6")}},
		// A present-but-null field is a name the SDK may not carry, so it must
		// survive as a null rather than vanish. A nil is what the consumer
		// classifies as "null"; the string "null" would be classified "string".
		{"null is recorded", `{"outstanding":null}`, map[string]any{"outstanding": nil}},
		{"array of objects", `[{"a":1}]`, []any{map[string]any{"a": json.Number("1")}}},
		{"empty array is not a wrong kind", `[]`, []any{}},
		{"bool", `{"ok":true}`, map[string]any{"ok": true}},
		{"nested", `{"d":{"e":[]}}`, map[string]any{"d": map[string]any{"e": []any{}}}},
		// UseNumber is load-bearing for a reason 1 and 1.0 cannot show: 1e400 is
		// valid JSON that a float64 decoder cannot represent, so without
		// UseNumber this body fails to decode and the probe records an
		// unparseable live body instead of a skeleton. That is a body silently
		// dropped from the evidence set, which is the one failure this harness
		// cannot detect about itself.
		{"number a float64 cannot hold", `{"big":1e400}`, map[string]any{"big": json.Number("1e400")}},
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
		{"integer stays a number", `42`, json.Number("42")},
		{"decimal", `385.6`, json.Number("385.6")},
		{"quoted number is a string", `"385.6"`, "1"},
		{"bool", `true`, true},
		// The boolean placeholder does not echo its input. Emitting one
		// unconditional value is what makes a `false` leaf in a reduced tree
		// proof that the tree was not reduced, so the case is pinned here rather
		// than assumed.
		{"false is the same placeholder as true", `false`, true},
		// A scalar body of a number no float64 could hold reduces like any
		// other, which is the top-level half of the UseNumber argument.
		{"unrepresentable number", `1e400`, json.Number("1e400")},
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
			[]any{json.Number("1"), "1", true, nil, []any{}, map[string]any{}},
		},
		{
			"a null element is an element",
			`[null,null]`,
			[]any{nil, nil},
		},
		{
			"nested arrays reduce element-wise",
			`[[1,["2"]],[]]`,
			[]any{[]any{json.Number("1"), []any{"1"}}, []any{}},
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
// server sent and cannot carry a value it sent. The walk is the assertion, and
// the corpus is chosen for the inputs that break a reduction rather than for
// tidiness.
func TestSkeletonifyNoLiveValueSurvives(t *testing.T) {
	corpus := []struct {
		name string
		in   string
	}{
		// A 40-digit integer is exact as a json.Number and already rounded as a
		// float64, so it is where a widened decode would show first.
		{"forty digit integer", `{"id":1234567890123456789012345678901234567890,"acct":"9110101000000000001"}`},
		// 1e400 overflows float64 to +Inf and is a decode error without
		// UseNumber, so it is the case that makes the option load-bearing.
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

// assertNoLiveValue walks a reduced tree and fails on any leaf that could carry
// something the server sent.
//
// A string leaf must be the string placeholder, so a symbol or an account id
// surviving is a failure and the path names the member it survived in. A
// float64 or an int is a failure too, and of a different kind: neither is a type
// this reduction emits, so one appearing means the value was widened on the way
// through, which is what a decode without UseNumber does. A bool leaf cannot
// carry a value, because the reduction emits one unconditional placeholder, so
// `false` is a failure: it proves the tree did not come from Skeletonify. A
// json.Number is accepted and is the one leaf that keeps the server's own text;
// the caller must extract kinds from a tree rather than serialise it, which is
// the recorded consequence of the representation.
func assertNoLiveValue(t *testing.T, v any, path string) {
	t.Helper()
	switch leaf := v.(type) {
	case map[string]any:
		for name, member := range leaf {
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
		// The literal text is the kind's evidence, and this is the only leaf
		// type that keeps anything the server wrote.
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

func TestSkeletonKind(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want string
	}{
		{"object", map[string]any{"a": "1"}, "object"},
		{"array", []any{json.Number("1")}, "array"},
		{"string placeholder", "1", "string"},
		{"number", json.Number("385.6"), "number"},
		{"boolean", true, "boolean"},
		{"empty array is still an array", []any{}, "array"},
		{"nil is null", nil, "null"},
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
