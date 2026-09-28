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
	"testing"
)

// The Display Solution bars pages require times and special_times, and neither field
// existed on StockBars or BatchBars, so the display host's trading calendar decoded to
// nil with no error reported. These tests decode the documented bodies, because the
// round-trip tests that already existed for these types could not see it: marshalling
// the SDK's own type and reading it back is green whether the shape is right or wrong.

const (
	documentedStockBars = `{
		"instrument_id": "913256135",
		"symbol": "AAPL",
		"result": [{"timestamp": 1768872168870, "open": "1.00", "high": "1.10",
		            "low": "0.90", "close": "1.05", "volume": "1000", "amount": "1050.0"}],
		"times": [
			{"start": "09:30:00", "end": "16:00:00", "trading_session": "RTH"},
			{"start": "04:00:00", "end": "09:30:00", "trading_session": "PRE"}
		],
		"special_times": [
			{"start": 0, "end": 0, "trading_session": "RTH"}
		]
	}`

	documentedBatchBars = `{
		"result": [{"symbol": "AAPL", "instrument_id": "913256135", "result": []}],
		"times": [{"start": "09:30:00", "end": "16:00:00", "trading_session": "RTH"}],
		"special_times": [{"start": 0, "end": 0, "trading_session": "RTH"}]
	}`
)

// TestStockBarsDecodeTheDocumentedHours covers the point of the two new types. The page
// documents the two hours arrays with the same three required names but different types
// for start and end, so they are two types and not one with a widened field; a caller
// reading either array must still get the value the page sends.
func TestStockBarsDecodeTheDocumentedHours(t *testing.T) {
	var got StockBars
	if err := json.Unmarshal([]byte(documentedStockBars), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got.Times) != 2 {
		t.Fatalf("len(Times) = %d, want 2: the exchange calendar a response was "+
			"computed against decoded to nil with no error reported", len(got.Times))
	}
	if got.Times[0].Start != "09:30:00" || got.Times[0].End != "16:00:00" {
		t.Errorf("Times[0] = %+v, want the documented wall-clock strings", got.Times[0])
	}
	if got.Times[0].TradingSession != "RTH" {
		t.Errorf("Times[0].TradingSession = %q, want RTH", got.Times[0].TradingSession)
	}
	if got.Times[1].TradingSession != "PRE" {
		t.Errorf("Times[1].TradingSession = %q, want PRE: the page carries more than "+
			"one session, so a single session field could not hold them all",
			got.Times[1].TradingSession)
	}
	// The special-hours array documents start and end as integers. If SpecialExchangeTimes
	// reused ExchangeTimes these would have failed to decode, which is why the two are
	// separate types.
	if len(got.SpecialTimes) != 1 {
		t.Fatalf("len(SpecialTimes) = %d, want 1", len(got.SpecialTimes))
	}
	if got.SpecialTimes[0].Start != 0 || got.SpecialTimes[0].End != 0 {
		t.Errorf("SpecialTimes[0] = %+v, want the documented integer form", got.SpecialTimes[0])
	}
}

func TestBatchBarsDecodeTheDocumentedHours(t *testing.T) {
	var got BatchBars
	if err := json.Unmarshal([]byte(documentedBatchBars), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got.Times) != 1 || got.Times[0].Start != "09:30:00" {
		t.Errorf("Times = %+v, want one entry starting 09:30:00", got.Times)
	}
	if len(got.SpecialTimes) != 1 || got.SpecialTimes[0].TradingSession != "RTH" {
		t.Errorf("SpecialTimes = %+v, want one RTH entry", got.SpecialTimes)
	}
}

// TestExchangeTimesRejectTheOtherArrayShape pins the distinction between the two
// element types. This is a negative test: it states the price the two schemas disagree
// on, so a future attempt to merge them into one type fails here rather than silently
// dropping whichever form the merged field could not hold.
func TestExchangeTimesRejectTheOtherArrayShape(t *testing.T) {
	// A string start cannot decode into the integer field the page documents for
	// special hours, and the reverse likewise. That is the intended behaviour, and it
	// is why SpecialExchangeTimes is not an alias of ExchangeTimes.
	var special SpecialExchangeTimes
	if err := json.Unmarshal([]byte(`{"start":"09:30:00","end":"16:00:00","trading_session":"RTH"}`), &special); err == nil {
		t.Error("SpecialExchangeTimes accepted a quoted start: the page documents " +
			"special-hours start and end as integers, so accepting a string here would " +
			"hide a real disagreement between the two arrays")
	}
}
