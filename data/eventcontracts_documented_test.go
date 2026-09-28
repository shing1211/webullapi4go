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

// The event-contract category and event types gained the names their pages document,
// alongside the SDK's own spellings. Both sets are carried because a response that
// honours the documentation previously left every required value unreadable, and a
// required name is one a conforming server always sends - so the addition cannot break
// a call that works today and needs no probe.
//
// These tests decode literal documented bodies, so a tag typo is caught here rather
// than by the conformance gate noticing a name went missing again.

func TestEventContractCategoryReadsTheDocumentedNames(t *testing.T) {
	const body = `[{"category_id": 1, "category_code": "ECONOMICS", "category_name": "Economics"}]`
	var got []EventContractCategory
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	c := got[0]
	if c.CategoryID != 1 {
		t.Errorf("CategoryID = %d, want 1: the documented name reached no field", c.CategoryID)
	}
	if c.CategoryCode != "ECONOMICS" {
		t.Errorf("CategoryCode = %q, want ECONOMICS", c.CategoryCode)
	}
	if c.CategoryName != "Economics" {
		t.Errorf("CategoryName = %q, want Economics", c.CategoryName)
	}
	// The SDK's own spellings are absent from the body, so they must stay zero rather
	// than borrowing the documented values.
	if c.Category != "" || c.Name != "" {
		t.Errorf("Category = %q and Name = %q, want both empty: a name the body did not "+
			"send was populated", c.Category, c.Name)
	}
}

func TestEventContractEventReadsTheDocumentedNames(t *testing.T) {
	const body = `[{
		"symbol": "KXCPI-26JAN-T0.3",
		"series_id": "KXCPI",
		"name": "CPI January 2026",
		"status": "ACTIVE",
		"short_name": "CPI Jan 26",
		"strike_date": "2026-01-01",
		"strike_period": "MONTH",
		"mutually_exclusive": true
	}]`
	var got []EventContractEvent
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	e := got[0]
	for _, c := range []struct {
		field, got, want string
	}{
		{"Symbol", e.Symbol, "KXCPI-26JAN-T0.3"},
		{"SeriesID", e.SeriesID, "KXCPI"},
		{"Name", e.Name, "CPI January 2026"},
		{"Status", e.Status, "ACTIVE"},
		{"ShortName", e.ShortName, "CPI Jan 26"},
		{"StrikeDate", e.StrikeDate, "2026-01-01"},
		{"StrikePeriod", e.StrikePeriod, "MONTH"},
	} {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.field, c.got, c.want)
		}
	}
	if !e.MutuallyExclusive {
		t.Error("MutuallyExclusive = false, want true: the documented boolean reached no field")
	}
	// The SDK's own spellings are absent from the body.
	if e.EventSymbol != "" || e.SeriesSymbol != "" {
		t.Errorf("EventSymbol = %q and SeriesSymbol = %q, want both empty: a name the body "+
			"did not send was populated", e.EventSymbol, e.SeriesSymbol)
	}
}

// TestEventContractEventStillReadsTheSDKNames is the other half, and it is what makes
// the change additive rather than a rename. Which naming a live server sends is
// unverified, so a body carrying the SDK's own names must still decode.
func TestEventContractEventStillReadsTheSDKNames(t *testing.T) {
	const body = `[{
		"event_symbol": "KXCPI-26JAN-T0.3",
		"series_symbol": "KXCPI",
		"name": "CPI January 2026",
		"status": "ACTIVE"
	}]`
	var got []EventContractEvent
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].EventSymbol != "KXCPI-26JAN-T0.3" || got[0].SeriesSymbol != "KXCPI" {
		t.Errorf("EventSymbol = %q, SeriesSymbol = %q, want the SDK's own names to still read",
			got[0].EventSymbol, got[0].SeriesSymbol)
	}
	if got[0].Symbol != "" || got[0].SeriesID != "" {
		t.Errorf("Symbol = %q, SeriesID = %q, want both empty: the body sent the SDK's own "+
			"spellings, so the documented ones must not borrow them",
			got[0].Symbol, got[0].SeriesID)
	}
}
