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
	"reflect"
	"testing"
)

func TestSkeletonify(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want any
	}{
		{"object", `{"a":"1"}`, map[string]any{"a": "string"}},
		// A number must stay a number: Money and QuoteTime accept both forms, so
		// conflating them is the one mistake this function exists to prevent.
		{"number is not string", `{"close":385.6}`, map[string]any{"close": "number"}},
		// A present-but-null field is a name the SDK may not carry, so it must
		// survive as "null" rather than vanish.
		{"null is recorded", `{"outstanding":null}`, map[string]any{"outstanding": "null"}},
		{"array of objects", `[{"a":1}]`, []any{map[string]any{"a": "number"}}},
		{"empty array is not a wrong kind", `[]`, []any{}},
		{"bool", `{"ok":true}`, map[string]any{"ok": "boolean"}},
		{"nested", `{"d":{"e":[]}}`, map[string]any{"d": map[string]any{"e": []any{}}}},
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
		{"null", `null`, "null"},
		{"integer stays a number", `42`, "number"},
		{"decimal", `385.6`, "number"},
		{"quoted number is a string", `"385.6"`, "string"},
		{"bool", `true`, "boolean"},
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

func TestSkeletonKind(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want string
	}{
		{"object", map[string]any{"a": "string"}, "object"},
		{"array", []any{"number"}, "array"},
		{"kind name is itself", "number", "number"},
		{"empty array is still an array", []any{}, "array"},
		{"nil has no kind", nil, ""},
		// A bool is not a kind: Skeletonify reduces one to "boolean" first, so
		// reaching this function with a bool means a caller skipped the
		// reduction, and reporting "" says so instead of guessing.
		{"unreduced bool has no kind", true, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SkeletonKind(tc.in); got != tc.want {
				t.Errorf("SkeletonKind(%#v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
