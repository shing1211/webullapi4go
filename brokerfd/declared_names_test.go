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

package brokerfd

import (
	"encoding/json"
	"testing"
)

// The Broker FD half of the v2.1.34 batch: 17 declared names across three types. Every
// page here publishes its properties but no required list, so a name may be optional or
// conditionally sent: the fields exist because the page publishes them, not because a
// server must send them. None is live-verified, because the Broker FD host returns 404
// in the HK sandbox and no US credential was available here.

// TestFDEnumReadsTheDeclaredNames covers the three names, two of which duplicate facts
// this type already carries under its own spelling.
func TestFDEnumReadsTheDeclaredNames(t *testing.T) {
	var got FDEnum
	if err := json.Unmarshal([]byte(`{
		"code": "LIMITED", "name": "Limited", "parent_code": "121"
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Code != "LIMITED" {
		t.Errorf("Code = %q, want LIMITED: this type spells the same fact Value, so a "+
			"caller reading the page's name got an empty string before", got.Code)
	}
	if got.Name != "Limited" {
		t.Errorf("Name = %q, want Limited: this type spells the same fact Label", got.Name)
	}
	if got.ParentCode != "121" {
		t.Errorf("ParentCode = %q, want 121: without it a caller cannot tell which "+
			"category a code belongs to, which is what the hierarchy is for", got.ParentCode)
	}
	// A body sending only the page's names must not populate the SDK's own pair, which
	// is what carrying both sets means.
	if got.Value != "" || got.Label != "" {
		t.Errorf("Value/Label = %q/%q, want them still empty", got.Value, got.Label)
	}
}

// TestFDOrderReadsTheDocumentedCorrelationKey covers the name the place and replace
// pages declare. It is the caller's own key, which is what makes an order traceable to
// the request that placed it.
func TestFDOrderReadsTheDocumentedCorrelationKey(t *testing.T) {
	var got FDOrder
	if err := json.Unmarshal([]byte(`{
		"order_id": "0352U72LQI6DT0KF41GK000000",
		"symbol": "AAPL", "status": "SUBMITTED",
		"client_order_id": "THI82O5JB7MQ2K76LL5FSDS2CB"
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.ClientRequestID != "THI82O5JB7MQ2K76LL5FSDS2CB" {
		t.Errorf("ClientRequestID = %q, want the documented value: it is the caller's own "+
			"correlation key, so a response echoing it is how an order is matched back "+
			"to the request that placed it", got.ClientRequestID)
	}
}

// TestFDCorporateActionReadsTheDeclaredNames covers the largest single batch, including
// the two structured names where the page nests an object and an array.
func TestFDCorporateActionReadsTheDeclaredNames(t *testing.T) {
	var got FDCorporateAction
	if err := json.Unmarshal([]byte(`{
		"action_id": "A1", "symbol": "AAPL", "action_type": "DIVIDEND",
		"category": "US_STOCK",
		"country_code": "US", "issuer_country_code": "US",
		"listing_country_of_code": "US",
		"event_id": "1234567890", "event_type": "DIVIDEND", "event_version": "1",
		"payment_date": "2025-01-15", "final_pay_date": "2025-02-15",
		"instrument_id": "943a9802f6c14983b3b4755c69c0",
		"from": {"symbol": "AAPL", "name": "Apple Inc", "exchange": "NAQ"},
		"to": [{"default_option_flag": "true", "description": "Cash",
			"option_number": "1", "payouts": {"amount": "1.00"}}]
	}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for name, pair := range map[string][2]string{
		"Category":             {got.Category, "US_STOCK"},
		"CountryCode":          {got.CountryCode, "US"},
		"IssuerCountryCode":    {got.IssuerCountryCode, "US"},
		"ListingCountryOfCode": {got.ListingCountryOfCode, "US"},
		"EventID":              {got.EventID, "1234567890"},
		"EventType":            {got.EventType, "DIVIDEND"},
		"EventVersion":         {got.EventVersion, "1"},
		"PaymentDate":          {got.PaymentDate, "2025-01-15"},
		"FinalPayDate":         {got.FinalPayDate, "2025-02-15"},
		"InstrumentID":         {got.InstrumentID, "943a9802f6c14983b3b4755c69c0"},
	} {
		if pair[0] != pair[1] {
			t.Errorf("%s = %q, want %q", name, pair[0], pair[1])
		}
	}
	// from and to are the two sides of the event, and without them a caller cannot
	// tell a split from a rename: a split moves holders from one instrument to
	// several, and the page's to is an array for exactly that reason.
	if got.From.Symbol != "AAPL" || got.From.Name != "Apple Inc" || got.From.Exchange != "NAQ" {
		t.Errorf("From = %+v, want the published position instrument", got.From)
	}
	if len(got.To) != 1 {
		t.Fatalf("len(To) = %d, want 1: the page documents an array, so a split with "+
			"several targets must all be readable", len(got.To))
	}
	if got.To[0].OptionNumber != "1" || got.To[0].Description != "Cash" {
		t.Errorf("To[0] = %+v", got.To[0])
	}
	// Payouts is a free-form object the page declares without naming any property
	// inside, so it is carried untyped. That is a weaker result than a typed struct and
	// the reason is recorded on the field rather than hidden.
	if got.To[0].Payouts["amount"] != "1.00" {
		t.Errorf("To[0].Payouts = %v, want the published value", got.To[0].Payouts)
	}
}
