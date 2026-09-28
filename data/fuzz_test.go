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

package data_test

import (
	"encoding/json"
	"testing"

	"github.com/shing1211/webullapi4go/data"
)

// The three fuzz targets existed with one seed each, all the literal
// {"symbol":"AAPL"}, and each discarded the decode error. A single trivial seed
// with the result thrown away can only catch a panic, so `make fuzz` was a one-second
// exercise of almost nothing. The seeds below are the shapes worth running, and
// because Go replays a fuzz target's seeds as ordinary tests, every one of them is
// also a standing regression case in `go test ./data/`.
//
// Three groups, chosen from what this SDK actually gets wrong rather than from what
// is easy to type:
//
//   - The documented shape, and the shape the sandbox really sends. money.Money
//     decodes a JSON string or a number because Webull sends both, and Quote.QuoteTime
//     is declared int64 while the reference page documents a string, which is the
//     open leaf-type defect in item 21. A seed per side of each disagreement is the
//     cheapest way to keep both paths exercised.
//   - The container inversions the pages are ambiguous about. Several documented
//     endpoints put an array on one path and a single item on another, so every type
//     here is seeded with the wrong shape for one of its fields.
//   - The degenerate inputs a hostile or merely careless server produces: null,
//     empty, wrong-typed, truncated, and non-ASCII. None should panic, and a decode
//     error is a legitimate outcome rather than a failure, which is why the fuzz
//     body discards it deliberately.

func FuzzDecodeQuote(f *testing.F) {
	// The documented shape: a book with levels on both sides.
	f.Add(`{"symbol":"AAPL","instrument_id":"913243251","quote_time":1755486723000,` +
		`"asks":[{"price":"175.00","volume":"100"}],"bids":[{"price":"174.90","volume":"250"}]}`)

	// quote_time as a string, which is what the reference page documents and what
	// the SDK's int64 field cannot hold. This is the leaf-type divergence recorded
	// as item 21, and it currently decodes to the zero value with no error.
	f.Add(`{"symbol":"AAPL","quote_time":"1755486723000"}`)

	// The book fields as a single object rather than an array of levels, which is
	// the inversion several documented pages carry.
	f.Add(`{"symbol":"AAPL","asks":{"price":"175.00"},"bids":{"price":"174.90"}}`)

	// quote_time as a non-numeric string, an array, and an object: each must not
	// panic and each must leave the field at a usable zero.
	f.Add(`{"symbol":"AAPL","quote_time":"not-a-timestamp"}`)
	f.Add(`{"symbol":"AAPL","quote_time":[1755486723000]}`)
	f.Add(`{"symbol":"AAPL","quote_time":{"ms":1755486723000}}`)

	// Degenerate bodies.
	f.Add(`{}`)
	f.Add(`[]`)
	f.Add(`null`)
	f.Add(`{"symbol":null,"asks":null,"bids":null}`)
	f.Add(`{"symbol":"","asks":[],"bids":[]}`)

	// Truncated and trailing-garbage bodies: a server that closes the connection
	// mid-write must not crash a decoder.
	f.Add(`{"symbol":"AAPL"`)
	f.Add(`{"symbol":"AAPL"} trailing`)
	f.Add(`{"symbol":"AAPL","asks":[}`)

	// Non-ASCII and escaped content in a field the decoder copies through.
	f.Add(`{"symbol":"腾讯","instrument_id":"00700","asks":[{"price":"1.00","volume":"1"}]}`)
	f.Add(`{"symbol":"A\\u0041PL","quote_time":0}`)

	f.Fuzz(func(t *testing.T, input string) {
		var v data.Quote
		_ = json.Unmarshal([]byte(input), &v)
	})
}

func FuzzDecodeSnapshot(f *testing.F) {
	// The documented shape, with every money field as the decimal string the
	// reference gives.
	f.Add(`{"instrument_id":"913243251","symbol":"AAPL","pre_close":"174.00",` +
		`"change_ratio":"0.0057","price":"175.00","open":"174.50","last_trade_time":1755486723000,` +
		`"ask":[{"price":"175.00","volume":"100"}],"bid":[{"price":"174.90","volume":"250"}],` +
		`"ask_p":[1755486723000,1755486723000],"bid_p":[1755486723000,1755486723000]}`)

	// The same money fields as JSON numbers, because money.Money accepts both and
	// Webull sends both across endpoints. If only one path is exercised, the other
	// is untested until a caller hits it.
	f.Add(`{"symbol":"AAPL","pre_close":174.00,"price":175.00,"open":174.50,"change_ratio":0.0057}`)

	// Money fields carrying a number, a bool, null, an array and an object, none of
	// which is a decimal.
	f.Add(`{"symbol":"AAPL","pre_close":true}`)
	f.Add(`{"symbol":"AAPL","price":null}`)
	f.Add(`{"symbol":"AAPL","open":[]}`)
	f.Add(`{"symbol":"AAPL","pre_close":{"value":"1.00"}}`)

	// last_trade_time as the documented string, and as the other non-numeric
	// shapes an int64 cannot take.
	f.Add(`{"symbol":"AAPL","last_trade_time":"1755486723000"}`)
	f.Add(`{"symbol":"AAPL","last_trade_time":1.5e300}`)
	f.Add(`{"symbol":"AAPL","last_trade_time":[1,2,3]}`)

	// ask/bid as a single object rather than an array of levels, and the timestamp
	// arrays as a single number, which is the inversion the documented shape carries
	// between endpoints on this surface.
	f.Add(`{"symbol":"AAPL","ask":{"price":"175.00"},"bid":{"price":"174.90"}}`)
	f.Add(`{"symbol":"AAPL","ask_p":1755486723000,"bid_p":1755486723000}`)

	// Degenerate, truncated and non-ASCII bodies.
	f.Add(`{}`)
	f.Add(`[]`)
	f.Add(`null`)
	f.Add(`{"symbol":"腾讯","price":"1.00"}`)
	f.Add(`{"symbol":"AAPL"`)
	f.Add(`{"symbol":"AAPL"} junk`)

	f.Fuzz(func(t *testing.T, input string) {
		var v data.Snapshot
		_ = json.Unmarshal([]byte(input), &v)
	})
}

func FuzzDecodeTick(f *testing.F) {
	// The documented shape: a result list of executed trades.
	f.Add(`{"symbol":"AAPL","instrument_id":"913243251","result":[` +
		`{"time":1755486723000,"price":"175.00","volume":"100","type":1},` +
		`{"time":1755486722000,"price":"174.95","volume":"50","type":1}]}`)

	// result as a single tick rather than a list, which is the container inversion
	// these documented pages carry.
	f.Add(`{"symbol":"AAPL","result":{"time":1755486723000,"price":"175.00","volume":"100"}}`)

	// Each trade field carrying a shape the decoder must survive.
	f.Add(`{"symbol":"AAPL","result":[{"time":"1755486723000","price":"175.00","volume":"100"}]}`)
	f.Add(`{"symbol":"AAPL","result":[{"time":null,"price":null,"volume":null,"type":null}]}`)
	f.Add(`{"symbol":"AAPL","result":[{"time":[1],"price":{"a":1},"volume":"abc","type":{}}]}`)

	// Degenerate, truncated and non-ASCII bodies.
	f.Add(`{}`)
	f.Add(`[]`)
	f.Add(`null`)
	f.Add(`{"symbol":"腾讯","result":[]}`)
	f.Add(`{"symbol":"AAPL","result":[{`)
	f.Add(`{"symbol":"AAPL","result":[]} trailing`)

	f.Fuzz(func(t *testing.T, input string) {
		var v data.StockTicks
		_ = json.Unmarshal([]byte(input), &v)
	})
}
