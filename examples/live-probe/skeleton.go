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
	"errors"
)

// Skeletonify decodes raw as JSON and reduces it to a value-free type skeleton.
//
// The reduction keeps every name the server sent and discards every value. An
// object reduces to a map from each member name to that member's reduced form,
// an array reduces to a slice of the reduced elements, and every scalar leaf
// reduces to its kind name: "string", "number", "boolean", or "null". A member
// that is present and null therefore survives as "null" instead of vanishing,
// because a name the SDK does not carry is precisely what the harness is
// looking for.
//
// The decoder is configured with UseNumber, so a JSON integer is not widened to
// a float64 and then reported as something a decimal could equally be. The
// SDK's own permissiveness must not be allowed to mask what the server sent:
// money.Money and data.QuoteTime each decode a JSON string or a number, and a
// skeleton that called both "string" would hide a wire change the SDK would
// absorb silently.
//
// An array reduces to a slice and never to nil, so an empty body stays
// distinguishable from a body whose top-level kind is wrong. A repeated member
// name collapses to one entry, as encoding/json already collapses it.
//
// An error is returned only when raw is not a single well-formed JSON value, in
// which case a valid value never fails to reduce. The decoder's error is
// returned unwrapped, so errors.As still reaches the *json.SyntaxError and
// *json.UnmarshalTypeError it carries.
func Skeletonify(raw []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()

	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	// One JSON value is a whole body. A second value means the reduction below
	// would describe part of what was sent while appearing to describe all of
	// it, which is the kind of quiet success this probe exists to rule out.
	if dec.More() {
		return nil, errors.New("live-probe: trailing data after the JSON body")
	}
	return reduceValue(v), nil
}

// SkeletonKind names the kind of a reduced value: "object" for a map, "array"
// for a slice, and the kind name itself for the string a scalar leaf reduces
// to. A nil or unrecognised value is reported as "", so a one-line summary
// never claims a kind it cannot support. Skeletonify reduces a bool to the
// string "boolean" and never returns an int, so neither is a kind here.
func SkeletonKind(v any) string {
	switch t := v.(type) {
	case map[string]any:
		return "object"
	case []any:
		return "array"
	case string:
		return t
	}
	// A nil value, and any type Skeletonify does not produce, has no kind.
	return ""
}

// reduceValue reduces one decoded JSON value to its value-free form. It is
// total over the types a json.Decoder configured with UseNumber produces:
// map[string]any, []any, string, json.Number, bool, and nil.
func reduceValue(v any) any {
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for name, member := range t {
			out[name] = reduceValue(member)
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, element := range t {
			out[i] = reduceValue(element)
		}
		return out
	case string:
		return "string"
	case json.Number:
		return "number"
	case bool:
		return "boolean"
	}
	// A JSON null decodes to a nil interface, and nil is the only value left
	// once the cases above are matched.
	return "null"
}
