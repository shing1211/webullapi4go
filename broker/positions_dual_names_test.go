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

package broker_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/shing1211/webullapi4go/broker"
)

// The documented response for `GET /broker/assets/positions/list` requires
// `cost_price`, `last_price` and `unrealized_profit_loss`, while this type has
// always carried `average_cost`, `market_value` and `unrealized_pl`. The two sets
// are now held side by side because which a live server sends is unverified: a
// retag would break a call that currently works for anyone whose server sends the
// SDK's own names. These tests decode literal bodies rather than round-tripping the
// SDK's own type, which is the form of test that cannot tell whether either set is
// read — see brokerfd/positions_dual_names_test.go for the longer account of why.

const documentedPositionsBody = `{
  "data": [
    {
      "symbol": "AAPL",
      "quantity": "100",
      "instrument_type": "STOCK",
      "currency": "USD",
      "cost_price": "150.00",
      "last_price": "175.00",
      "unrealized_profit_loss": "2500.00"
    }
  ]
}`

const sdkNamedPositionsBody = `{
  "data": [
    {
      "symbol": "AAPL",
      "quantity": "100",
      "instrument_type": "STOCK",
      "currency": "USD",
      "average_cost": "150.00",
      "market_value": "17500",
      "unrealized_pl": "2500.00"
    }
  ]
}`

func decodePositions(t *testing.T, body string) broker.Position {
	t.Helper()
	var envelope struct {
		Data []broker.Position `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(envelope.Data) != 1 {
		t.Fatalf("decoded %d positions, want 1", len(envelope.Data))
	}
	return envelope.Data[0]
}

func TestPositionReadsTheDocumentedNames(t *testing.T) {
	p := decodePositions(t, documentedPositionsBody)
	// These are string fields, not money.Money, so the value decodes verbatim and
	// keeps the scale the body wrote. That is the type's existing behaviour and
	// the test records it rather than normalising it.
	for _, c := range []struct{ field, got, want string }{
		{"CostPrice", p.CostPrice, "150.00"},
		{"LastPrice", p.LastPrice, "175.00"},
		{"UnrealizedProfitLoss", p.UnrealizedProfitLoss, "2500.00"},
	} {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q: the documented name reached no field", c.field, c.got, c.want)
		}
	}
	if p.AverageCost != "" {
		t.Errorf("AverageCost = %q, want empty: a name the body did not send was populated", p.AverageCost)
	}
	if p.MarketValue != "" {
		t.Errorf("MarketValue = %q, want empty: a name the body did not send was populated", p.MarketValue)
	}
}

func TestPositionStillReadsTheSDKNames(t *testing.T) {
	p := decodePositions(t, sdkNamedPositionsBody)
	for _, c := range []struct{ field, got, want string }{
		{"AverageCost", p.AverageCost, "150.00"},
		{"MarketValue", p.MarketValue, "17500"},
		{"UnrealizedPL", p.UnrealizedPL, "2500.00"},
	} {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.field, c.got, c.want)
		}
	}
	if p.CostPrice != "" {
		t.Errorf("CostPrice = %q, want empty: a name the body did not send was populated", p.CostPrice)
	}
}

// TestPositionHasNoDuplicateWireName pins the class rather than this struct's past.
// brokerfd.FDPosition shipped with two fields claiming `unrealized_pl`, which
// encoding/json resolves by dropping both, so its open-P&L field was silently zero
// and no lookup for a *missing* name would ever have found it. broker.Position does
// not have that defect; this says so in a way that fails if it ever acquires one.
func TestPositionHasNoDuplicateWireName(t *testing.T) {
	seen := map[string]string{}
	rt := reflect.TypeOf(broker.Position{})
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
			t.Errorf("Position fields %s and %s both encode under %q; encoding/json "+
				"drops both when two fields of one struct claim the same name, so "+
				"neither is ever filled and neither reads as missing", prev, f.Name, raw)
			continue
		}
		seen[raw] = f.Name
	}
	if len(seen) == 0 {
		t.Fatal("no json tags were read, so this checked nothing")
	}
}
