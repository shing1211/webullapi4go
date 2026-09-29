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

// stringPlaceholder is what a JSON string reduces to. It is synthetic, so a live
// value cannot be mistaken for it, and it is not empty, so a reduced string is
// visibly a value the server sent rather than an absent one.
const stringPlaceholder = "1"

// boolPlaceholder is what a JSON boolean reduces to. It is emitted
// unconditionally: the reduction does not echo whether the server sent true or
// false, so a `false` anywhere in a reduced tree proves the tree did not come
// from Skeletonify. The type carries the kind and a bool has no room for
// anything else, so mirroring the input would leak nothing and would weaken that
// check.
const boolPlaceholder = true

// numberPlaceholder is what a JSON number reduces to. It is one fixed literal
// for every number, so a committed fixture is a file of this number repeated
// rather than of readings, and any other number in a reduced tree proves the
// tree did not come from Skeletonify.
//
// The literal is -1 because it is not a figure Webull sends for a price, a
// size, a timestamp, a ratio or a count, so a reduced tree cannot be read as a
// snapshot. 0 would be the obvious alternative and the wrong one: a zero
// balance is a real reading, so a leaked 0 would look like data.
const numberPlaceholder = json.Number("-1")

// Skeletonify decodes raw as JSON and reduces it to a value-free type skeleton.
//
// The reduction keeps every name the server sent and discards every value. An
// object reduces to a map from each member name to that member's reduced form,
// an array reduces to a slice of the reduced elements, and every scalar leaf
// reduces to a typed placeholder carrying its kind and nothing else:
//
//   - a JSON string  -> the Go string "1"
//   - a JSON number  -> the json.Number "-1"
//   - a JSON boolean -> the Go bool true
//   - a JSON null    -> a nil interface
//
// The placeholders are typed, not named, because the consumer of a skeleton is
// conformance.CompareBody, which classifies a decoded value by its Go type in
// jsonKind (conformance/shapes.go:361). A leaf spelled as the string "number"
// would be classified "string" there, so every numeric leaf the server sent would
// be demanded of a Go string, a live null would yield no leaf row at all, and a
// number where a string is documented would be reported as agreement. The
// harness would then be unable to detect the wire change it exists to find. Each
// placeholder above is a value jsonKind classifies as exactly the kind it stands
// for, so a reduced tree is fixture-shaped input to that comparison unchanged.
//
// A reduced tree therefore contains no value the server sent: it carries every
// name and no reading. A consumer may serialise it, write it to a committed
// fixture, and put it in the repository, and nothing it writes is a live
// account number, a price, or a timestamp. That is a property of the
// representation rather than an obligation on the consumer, and it is why every
// leaf is synthetic: the literal a number decoded to is evidence of what the
// server sent, and evidence is what the reduction consumes, not what it emits.
//
// The decoder is configured with UseNumber, so a JSON integer is neither widened
// to a float64 nor reported as something a decimal could equally be. The SDK's
// own permissiveness must not be allowed to mask what the server sent:
// money.Money and data.QuoteTime each decode a JSON string or a number, and a
// skeleton that called both a string would hide a wire change the SDK would
// absorb silently. UseNumber is what makes the discard above possible rather
// than merely tidy: without it a 40-digit integer is rounded to a float64 that
// no longer equals the integer, and 1e400 is a decode error, which drops the
// body from the evidence set instead of reducing it.
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
// for a slice, "string" for the string placeholder, "number" for a json.Number,
// "boolean" for a bool, and "null" for a nil. Those six are the whole answer set
// and each is derived from the value's own type, so a kind is never an arbitrary
// string echoed back from the input.
//
// Any other value is reported as "", so a one-line summary never claims a kind
// it cannot support. That guard is closed: an int, a float64, a struct, or any
// other type Skeletonify does not produce is a value that skipped the reduction,
// and saying so beats naming a kind. A slice of unreduced live values is the one
// case the guard cannot catch at the top level, because its container type really
// is "array"; the live values are inside it, and the caller that walks the tree
// is what rejects them.
func SkeletonKind(v any) string {
	switch v.(type) {
	case map[string]any:
		return "object"
	case []any:
		return "array"
	case string:
		return "string"
	case json.Number:
		return "number"
	case bool:
		return "boolean"
	case nil:
		return "null"
	}
	// A type Skeletonify does not produce has no kind.
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
		return stringPlaceholder
	case json.Number:
		// The literal is evidence, not output. Its type already carries the
		// kind, which is all the comparison reads, so returning it would buy no
		// accuracy and would put a price, a share count, a timestamp or an
		// account number into every fixture this tree is written to. A live
		// value in a committed file cannot be taken back out of the history.
		return numberPlaceholder
	case bool:
		return boolPlaceholder
	}
	// A JSON null decodes to a nil interface, and nil is the only value left
	// once the cases above are matched.
	return nil
}
