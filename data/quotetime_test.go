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

package data

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

// Webull's own published pages disagree about quote_time: the stock depth and
// Display Solution depth examples show a quoted string, the futures and
// event-contract depth examples show a bare number. These tests pin that the SDK
// reads both, because the two failure modes are not symmetric. An int64 field
// decoded the number form and failed the string form with an error that abandoned
// the whole response, so a server honouring half its own documentation returned
// nothing at all.
func TestQuoteTimeDecodesEveryDocumentedShape(t *testing.T) {
	const want = QuoteTime(1640688000000)
	tests := []struct {
		name string
		body string
		want QuoteTime
	}{
		// The stock depth page's own example.
		{"documented string", `{"quote_time": "1640688000000"}`, want},
		// The futures and event-contract pages' own examples.
		{"documented number", `{"quote_time": 1761131409276}`, QuoteTime(1761131409276)},
		{"number with exponent", `{"quote_time": 1.761131409276e12}`, QuoteTime(1761131409276)},
		{"quoted number with a fraction", `{"quote_time": "1761131409276.0"}`, QuoteTime(1761131409276)},
		// An optional field the server omits or nulls is the zero value, not an
		// error: a caller asking for depth outside trading hours gets an empty book
		// rather than a failure.
		{"null", `{"quote_time": null}`, 0},
		{"omitted", `{}`, 0},
		{"empty string", `{"quote_time": ""}`, 0},
		{"zero", `{"quote_time": 0}`, 0},
		{"negative", `{"quote_time": -1}`, QuoteTime(-1)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got struct {
				QuoteTime QuoteTime `json:"quote_time"`
			}
			if err := json.Unmarshal([]byte(tc.body), &got); err != nil {
				t.Fatalf("unmarshal %s: %v", tc.body, err)
			}
			if got.QuoteTime != tc.want {
				t.Errorf("QuoteTime = %d, want %d", int64(got.QuoteTime), int64(tc.want))
			}
		})
	}
}

// TestQuoteTimeRejectsWhatIsNotATimestamp is the other half: accepting both
// documented shapes must not mean accepting anything.
func TestQuoteTimeRejectsWhatIsNotATimestamp(t *testing.T) {
	for _, body := range []string{
		`{"quote_time": "not a number"}`,
		`{"quote_time": "12x34"}`,
		`{"quote_time": true}`,
		`{"quote_time": {}}`,
		`{"quote_time": []}`,
	} {
		var got struct {
			QuoteTime QuoteTime `json:"quote_time"`
		}
		err := json.Unmarshal([]byte(body), &got)
		if err == nil {
			t.Errorf("unmarshal %s succeeded with %d, want an error", body, int64(got.QuoteTime))
			continue
		}
		if !errors.Is(err, ErrQuoteTimeFormat) {
			t.Errorf("unmarshal %s gave %v, want one wrapping ErrQuoteTimeFormat", body, err)
		}
	}
}

// TestQuoteTimeAccessors pins the conversions, so the type is usable without a
// caller reaching for strconv.
func TestQuoteTimeAccessors(t *testing.T) {
	q := QuoteTime(1640688000000)
	if got := q.Int64(); got != 1640688000000 {
		t.Errorf("Int64() = %d, want 1640688000000", got)
	}
	if got := q.String(); got != "1640688000000" {
		t.Errorf("String() = %q, want 1640688000000", got)
	}
	if got, want := q.Time().UTC(), time.UnixMilli(1640688000000).UTC(); !got.Equal(want) {
		t.Errorf("Time() = %s, want %s", got, want)
	}
	// The zero value is the Unix epoch, so Time() is well defined rather than
	// undefined.
	if got, want := QuoteTime(0).Time().UTC(), time.UnixMilli(0).UTC(); !got.Equal(want) {
		t.Errorf("zero Time() = %s, want %s", got, want)
	}
}

// TestQuoteTimeMarshalsAsANumber records the send-side choice. Nothing in this SDK
// posts a Quote to Webull, so the wire form is unconstrained; a number is what the
// field was before this type existed, so a caller that re-encodes a decoded value
// sees no change.
func TestQuoteTimeMarshalsAsANumber(t *testing.T) {
	b, err := json.Marshal(struct {
		QuoteTime QuoteTime `json:"quote_time"`
	}{QuoteTime: 1640688000000})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if got, want := string(b), `{"quote_time":1640688000000}`; got != want {
		t.Errorf("marshal = %s, want %s", got, want)
	}
}

// TestQuoteRoundTripsThroughBothForms is the property the type exists to provide,
// stated as a round trip rather than as two separate decode assertions.
func TestQuoteRoundTripsThroughBothForms(t *testing.T) {
	const stringForm = `{"symbol":"AAPL","quote_time":"1640688000000","asks":[],"bids":[]}`
	const numberForm = `{"symbol":"AAPL","quote_time":1640688000000,"asks":[],"bids":[]}`

	for _, body := range []string{stringForm, numberForm} {
		var q Quote
		if err := json.Unmarshal([]byte(body), &q); err != nil {
			t.Fatalf("unmarshal %s: %v", body, err)
		}
		if q.QuoteTime != 1640688000000 {
			t.Errorf("%s gave QuoteTime %d, want 1640688000000", body, int64(q.QuoteTime))
		}
	}
}
