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
	"reflect"
	"testing"
)

// This file exists because the response type for the documented
// `GET /broker/assets/positions/list` carries two sets of names for the same three
// values, and the tests beside it could not tell whether either set was read.
//
// The existing test marshals []FDPosition and decodes it back, which proves the
// type is self-consistent and nothing else. It cannot observe a name the server
// never sends, and it did not assert UnrealizedPL or RealizedPL at all — which is
// why it passed for the whole time those two fields shared one json tag, which is
// a state in which encoding/json drops BOTH. A self-round-trip test is green on a
// type that reads nothing at all, so these decode literal bodies instead.

const documentedPositionsBody = `[
  {
    "position_id": "P1",
    "symbol": "AAPL",
    "quantity": "100",
    "instrument_type": "STOCK",
    "currency": "USD",
    "cost_price": "150.00",
    "last_price": "175.00",
    "unrealized_profit_loss": "2500.00"
  }
]`

const sdkNamedPositionsBody = `[
  {
    "position_id": "P1",
    "symbol": "AAPL",
    "quantity": "100",
    "instrument_type": "STOCK",
    "currency": "USD",
    "average_cost": "150.00",
    "market_value": "17500",
    "unrealized_pl": "2500.00",
    "realized_pl": "10.00"
  }
]`

const bothSetsPositionsBody = `[
  {
    "position_id": "P1",
    "symbol": "AAPL",
    "quantity": "100",
    "instrument_type": "STOCK",
    "currency": "USD",
    "average_cost": "150.00",
    "market_value": "17500",
    "unrealized_pl": "2500.00",
    "realized_pl": "10.00",
    "cost_price": "150.00",
    "last_price": "175.00",
    "unrealized_profit_loss": "2500.00"
  }
]`

// TestFDPositionReadsTheDocumentedNames is the case that motivated the three
// added fields: a server that sends the published names must not leave a caller
// holding three zeroes.
func TestFDPositionReadsTheDocumentedNames(t *testing.T) {
	var got []FDPosition
	if err := json.Unmarshal([]byte(documentedPositionsBody), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	p := got[0]
	for _, c := range []struct {
		field, want string
		got         string
	}{
		{"CostPrice", "150", p.CostPrice.String()},
		{"LastPrice", "175", p.LastPrice.String()},
		{"UnrealizedProfitLoss", "2500", p.UnrealizedProfitLoss.String()},
	} {
		if c.got == "" || c.got == "0" {
			t.Errorf("%s = %q, want %q: the documented name reached no field", c.field, c.got, c.want)
			continue
		}
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.field, c.got, c.want)
		}
	}
	// The SDK's own names are absent from this body, so their fields must stay at
	// their zero value rather than borrowing the documented ones. That is the
	// decoder property carrying both sets rests on: a name the body did not send
	// populates nothing, and a name it did send populates exactly its own field.
	// It is not a claim about which set a live server sends, which is unverified.
	if p.AverageCost.String() != "0" {
		t.Errorf("AverageCost = %q, want 0: a name the body did not send was populated", p.AverageCost.String())
	}
	if p.MarketValue.String() != "0" {
		t.Errorf("MarketValue = %q, want 0: a name the body did not send was populated", p.MarketValue.String())
	}
}

// TestFDPositionStillReadsTheSDKNames is the other half. Which set a live server
// sends is unverified, so a body carrying the SDK's own names must still decode —
// that is what makes the change additive rather than a retag.
func TestFDPositionStillReadsTheSDKNames(t *testing.T) {
	var got []FDPosition
	if err := json.Unmarshal([]byte(sdkNamedPositionsBody), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	p := got[0]
	for _, c := range []struct {
		field, want string
		got         string
	}{
		{"AverageCost", "150", p.AverageCost.String()},
		{"MarketValue", "17500", p.MarketValue.String()},
		{"UnrealizedPL", "2500", p.UnrealizedPL.String()},
		{"RealizedPL", "10", p.RealizedPL.String()},
	} {
		if c.got == "" || c.got == "0" {
			t.Errorf("%s = %q, want %q", c.field, c.got, c.want)
		}
	}
	if p.CostPrice.String() != "0" {
		t.Errorf("CostPrice = %q, want 0: a name the body did not send was populated", p.CostPrice.String())
	}
}

// TestFDPositionReadsBothSetsAtOnce covers a server that sends both, which is the
// case a retag would have made impossible.
func TestFDPositionReadsBothSetsAtOnce(t *testing.T) {
	var got []FDPosition
	if err := json.Unmarshal([]byte(bothSetsPositionsBody), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	p := got[0]
	if p.CostPrice.String() != "150" || p.AverageCost.String() != "150" {
		t.Errorf("cost basis: CostPrice=%q AverageCost=%q, want both 150",
			p.CostPrice.String(), p.AverageCost.String())
	}
	if p.UnrealizedPL.String() != "2500" || p.UnrealizedProfitLoss.String() != "2500" {
		t.Errorf("open P&L: UnrealizedPL=%q UnrealizedProfitLoss=%q, want both 2500",
			p.UnrealizedPL.String(), p.UnrealizedProfitLoss.String())
	}
	if p.RealizedPL.String() != "10" {
		t.Errorf("RealizedPL = %q, want 10", p.RealizedPL.String())
	}
}

// TestFDPositionHasNoDuplicateWireName pins the defect this change corrected.
//
// When two fields of one struct encode under the same name, encoding/json drops
// BOTH rather than picking one, so a field can be present, tagged, and reported as
// covered by a tag lookup while never being filled. FDPosition shipped that way:
// UnrealizedPL and RealizedPL both carried `json:"unrealized_pl"`, so open P&L was
// silently zero for every caller and realized_pl was unreachable. No amount of
// looking for *missing* names finds it, because nothing is missing.
//
// This walks the struct rather than hard-coding the pair, so a different collision
// is caught too.
func TestFDPositionHasNoDuplicateWireName(t *testing.T) {
	seen := map[string]string{}
	rt := reflect.TypeOf(FDPosition{})
	for i := range rt.NumField() {
		f := rt.Field(i)
		raw, ok := f.Tag.Lookup("json")
		if !ok {
			continue
		}
		for i := 0; i < len(raw); i++ {
			if raw[i] == ',' {
				raw = raw[:i]
				break
			}
		}
		if raw == "" || raw == "-" {
			continue
		}
		if prev, dup := seen[raw]; dup {
			t.Errorf("FDPosition fields %s and %s both encode under %q; encoding/json "+
				"drops both when two fields of one struct claim the same name, so "+
				"neither is ever filled and neither reads as missing",
				prev, f.Name, raw)
			continue
		}
		seen[raw] = f.Name
	}
	if len(seen) == 0 {
		t.Fatal("no json tags were read, so this checked nothing")
	}
}
