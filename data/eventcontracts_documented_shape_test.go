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

	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

// Three event-contract endpoints decoded one nesting level too deep, or the wrong
// shape, until v2.1.32. The pages document:
//
//	GET /market-data/event-contracts/bars/list   200: array<{instrument_id, symbol, result: array<bar>}>
//	GET /market-data/event-contracts/tick/list   200: array<{instrument_id, symbol, result: array<tick>}>
//	GET /market-data/event-contracts/depths/list 200: array<{instrument_id, symbol, quote_time,
//	                                                              yes_bids, yes_asks, no_bids, no_asks}>
//
// The SDK returned the elements of the inner result array, so the grouping key was
// discarded and bars or ticks from several instruments arrived in one flat list. These
// tests decode the documented bodies, because the existing round-trip tests could not
// see this: marshalling the SDK's own shape and reading it back is green whether the
// shape is right or wrong.

const (
	documentedBars = `[{"instrument_id":"504279491","symbol":"KXCPI-26JAN-T0.3","result":[
		{"open":"3.40","high":"3.60","low":"3.35","close":"3.50","volume":"5000","time":"1768872168870"},
		{"open":"3.50","high":"3.70","low":"3.45","close":"3.60","volume":"4000","time":"1768872228870"}]}]`

	documentedTicks = `[{"instrument_id":"504279491","symbol":"KXCPI-26JAN-T0.3","result":[
		{"yes_price":"0.52","no_price":"0.48","side":"B","volume":"10","trade_id":"T1","time":"1768872168870"},
		{"yes_price":"0.51","no_price":"0.49","side":"S","volume":"5","trade_id":"T2","time":"1768872169870"}]}]`

	documentedDepth = `[{"instrument_id":"504279491","symbol":"KXCPI-26JAN-T0.3","quote_time":1768872168870,
		"yes_bids":[{"price":"0.53","size":"100"}],
		"yes_asks":[{"price":"0.55","size":"150"}],
		"no_bids":[{"price":"0.47","size":"300"}],
		"no_asks":[{"price":"0.45","size":"400"}]}]`

	documentedSnapshot = `[{"instrument_id":"504279491","symbol":"KXCPI-26JAN-T0.3",
		"name":"CPI January 2026","price":"0.52","volume":"1200","last_trade_time":1768872168870,
		"open_interest":"5000","yes_bid":"0.53","yes_bid_size":"100","yes_ask":"0.55","yes_ask_size":"150",
		"no_bid":"0.47","no_bid_size":"300","no_ask":"0.45","no_ask_size":"400"}]`
)

// TestEventBarsDecodeTheDocumentedGrouping is the depth fix. Two bars for one
// instrument arrive under the instrument they belong to, which is what the previous
// return type could not represent.
func TestEventBarsDecodeTheDocumentedGrouping(t *testing.T) {
	var got []EventBarsResult
	if err := json.Unmarshal([]byte(documentedBars), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1 element per instrument", len(got))
	}
	if got[0].InstrumentID != "504279491" {
		t.Errorf("InstrumentID = %q, want 504279491: the grouping key is the whole point "+
			"of the documented shape", got[0].InstrumentID)
	}
	if got[0].Symbol != "KXCPI-26JAN-T0.3" {
		t.Errorf("Symbol = %q", got[0].Symbol)
	}
	if len(got[0].Result) != 2 {
		t.Fatalf("len(Result) = %d, want 2", len(got[0].Result))
	}
	first := got[0].Result[0]
	if first.Open.Cmp(money.Must(money.NewFromString("3.40"))) != 0 {
		t.Errorf("first bar Open = %s, want 3.40", first.Open)
	}
	// The page requires time on every bar and does not mention timestamp, so a
	// conforming response reaches Time and leaves the SDK's own field empty.
	if first.Time != "1768872168870" {
		t.Errorf("Time = %q, want 1768872168870", first.Time)
	}
	if first.Timestamp != "" {
		t.Errorf("Timestamp = %q, want empty: the documented bar has no timestamp field, "+
			"so a populated value could only have come from the wrong wire name", first.Timestamp)
	}
}

func TestEventTicksDecodeTheDocumentedGrouping(t *testing.T) {
	var got []EventTickResult
	if err := json.Unmarshal([]byte(documentedTicks), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1 element per instrument", len(got))
	}
	if got[0].InstrumentID != "504279491" || got[0].Symbol != "KXCPI-26JAN-T0.3" {
		t.Errorf("grouping key = %q/%q", got[0].InstrumentID, got[0].Symbol)
	}
	if len(got[0].Result) != 2 {
		t.Fatalf("len(Result) = %d, want 2", len(got[0].Result))
	}
	if got[0].Result[0].TradeID != "T1" || got[0].Result[1].Side != "S" {
		t.Errorf("ticks = %+v", got[0].Result)
	}
	if got[0].Result[0].Time != "1768872168870" {
		t.Errorf("Time = %q, want 1768872168870", got[0].Result[0].Time)
	}
}

// TestEventDepthReadsBothSides is the part that was simply missing: the SDK carried no
// no-side book, so a caller could not read the side an event contract is named for.
func TestEventDepthReadsBothSides(t *testing.T) {
	var got []EventDepth
	if err := json.Unmarshal([]byte(documentedDepth), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1 element per instrument", len(got))
	}
	book := got[0]
	if book.InstrumentID != "504279491" {
		t.Errorf("InstrumentID = %q, want 504279491", book.InstrumentID)
	}
	if book.QuoteTime != 1768872168870 {
		t.Errorf("QuoteTime = %d, want 1768872168870: the page requires quote_time, which "+
			"the SDK named timestamp", book.QuoteTime)
	}
	if book.Timestamp != "" {
		t.Errorf("Timestamp = %q, want empty: the documented field is quote_time", book.Timestamp)
	}
	for name, side := range map[string][]DepthLevel{
		"YesBids": book.YesBids, "YesAsks": book.YesAsks,
		"NoBids": book.NoBids, "NoAsks": book.NoAsks,
	} {
		if len(side) != 1 {
			t.Errorf("len(%s) = %d, want 1: the no-side book is what this change makes "+
				"readable", name, len(side))
		}
	}
	if len(book.NoBids) == 1 && book.NoBids[0].Size != "300" {
		t.Errorf("NoBids[0].Size = %q, want 300", book.NoBids[0].Size)
	}
}

// TestEventSnapshotReadsEveryDocumentedName covers the 10 required names the type
// lacked, including the four size companions that make a snapshot's depth readable.
func TestEventSnapshotReadsEveryDocumentedName(t *testing.T) {
	var got []EventSnapshot
	if err := json.Unmarshal([]byte(documentedSnapshot), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	s := got[0]
	if s.InstrumentID != "504279491" {
		t.Errorf("InstrumentID = %q", s.InstrumentID)
	}
	if s.Name != "CPI January 2026" {
		t.Errorf("Name = %q", s.Name)
	}
	if s.Price.Cmp(money.Must(money.NewFromString("0.52"))) != 0 {
		t.Errorf("Price = %s, want 0.52: the page requires price, which the SDK named "+
			"last_price", s.Price)
	}
	if s.LastTradeTime != 1768872168870 {
		t.Errorf("LastTradeTime = %d, want 1768872168870", s.LastTradeTime)
	}
	for name, got := range map[string]string{
		"YesBidSize": s.YesBidSize, "YesAskSize": s.YesAskSize,
		"NoBidSize": s.NoBidSize, "NoAskSize": s.NoAskSize,
	} {
		if got == "" {
			t.Errorf("%s is empty: a snapshot's sizes are what make its depth readable", name)
		}
	}
	for name, ok := range map[string]bool{
		"NoBid": s.NoBid.IsZero(), "NoAsk": s.NoAsk.IsZero(),
	} {
		if ok {
			t.Errorf("%s is zero: the page requires the no side", name)
		}
	}
	// LastPrice is the SDK's own spelling; the documented body sends price, so the two
	// must not both read.
	if !s.LastPrice.IsZero() {
		t.Errorf("LastPrice = %s, want zero: the documented body sends price, not "+
			"last_price", s.LastPrice)
	}
}
