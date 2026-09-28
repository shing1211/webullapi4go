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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// ErrQuoteTimeFormat is returned when a quote_time value is neither a number nor a
// quoted decimal, which is not a shape Webull publishes. It is a sentinel so a
// caller can recognise a malformed timestamp without matching on message text.
var ErrQuoteTimeFormat = errors.New("quote_time is not a Unix millisecond timestamp as a " +
	"number or a quoted decimal")

// QuoteTime is a Unix epoch millisecond timestamp that decodes from either of the
// two JSON representations Webull publishes for the same field.
//
// Four documented pages carry a `quote_time` property, and Webull's own published
// examples disagree about its type: the stock depth and Display Solution depth pages
// show a quoted string ("1640688000000") while the futures and event-contract depth
// pages show a bare number (1761131409276). An int64 field therefore decodes two of
// them and fails the other two with an error that abandons the whole response, and
// the SDK cannot know which a live server will send.
//
// This type is the resolution that needs no probe: it reads both, so it is correct
// whichever the server does. A second field carrying the same name is not an
// alternative - two fields claiming one json tag is the defect described in
// conformance/tags.go, where encoding/json drops both.
//
// The zero value is the Unix epoch, matching the int64 this replaces. A field the
// server omits or sends as null decodes to it without an error, which is what an
// absent optional timestamp should do.
type QuoteTime int64

// epoch is the zero value in Unix milliseconds, used to turn the zero QuoteTime
// back into a time without a magic number at each call site.
const quoteTimeEpochMillis = 0

// UnmarshalJSON accepts a quoted decimal, a bare number, null, and an empty string.
//
// It is written the way money.Money.UnmarshalJSON is: ask for a string, and if that
// fails ask for a number. That is simpler than walking the bytes and handles escapes
// correctly, and the float fallback covers a server sending 1.761131409276e12 for the
// same instant.
func (q *QuoteTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		// null is not a number and not a string; an absent or null timestamp is the
		// zero value rather than an error, which is what an optional field should do.
		if string(bytes.TrimSpace(data)) == "null" {
			*q = 0
			return nil
		}
		var n float64
		if err2 := json.Unmarshal(data, &n); err2 != nil {
			// Neither a string nor a number: a bool, an object, an array. Wrapping
			// the sentinel keeps one recognisable error for every malformed shape,
			// rather than leaking the inner type error for some of them.
			return fmt.Errorf("%w: %v", ErrQuoteTimeFormat, err)
		}
		*q = QuoteTime(int64(n))
		return nil
	}
	if s == "" {
		*q = 0
		return nil
	}
	ms, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		// A quoted float, e.g. "1761131409276.0", is still the same instant.
		f, ferr := strconv.ParseFloat(s, 64)
		if ferr != nil {
			return ErrQuoteTimeFormat
		}
		ms = int64(f)
	}
	*q = QuoteTime(ms)
	return nil
}

// MarshalJSON writes the number, not the string.
//
// No endpoint in this SDK sends a Quote back to Webull, so nothing here is
// constrained by the wire format; the number is chosen because it is what the field
// held before this type existed, so a caller that round-trips a decoded value
// through encoding/json sees no change.
func (q QuoteTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(q), 10)), nil
}

// Int64 returns the timestamp in Unix milliseconds, for a caller that needs the
// underlying value.
func (q QuoteTime) Int64() int64 { return int64(q) }

// Time returns the timestamp as a UTC time. The zero QuoteTime is the Unix epoch.
func (q QuoteTime) Time() time.Time {
	return time.UnixMilli(int64(q) - quoteTimeEpochMillis).UTC()
}

// String renders the timestamp in Unix milliseconds, matching how it arrives.
func (q QuoteTime) String() string { return strconv.FormatInt(int64(q), 10) }
