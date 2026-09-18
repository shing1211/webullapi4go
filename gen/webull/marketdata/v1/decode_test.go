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

package marketdatav1

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestDecodeSnapshotRoundTrip(t *testing.T) {
	want := &Snapshot{
		Basic: &Basic{
			Symbol:       "AAPL",
			InstrumentId: "913256135",
			Timestamp:    "1700000000000",
		},
		TradeTime:      "1700000001000",
		Price:          "190.12",
		Open:           "189.50",
		High:           "191.00",
		Low:            "188.75",
		PreClose:       "189.00",
		Volume:         "1234567",
		Change:         "1.12",
		ChangeRatio:    "0.0059",
		ExtTradeTime:   "1700010001000",
		ExtPrice:       "190.55",
		ExtHigh:        "190.80",
		ExtLow:         "189.90",
		ExtVolume:      "23456",
		ExtChange:      "0.43",
		ExtChangeRatio: "0.0023",
		OvnTradeTime:   "1699990001000",
		OvnPrice:       "189.80",
		OvnHigh:        "190.00",
		OvnLow:         "189.20",
		OvnVolume:      "3456",
		OvnChange:      "0.30",
		OvnChangeRatio: "0.0016",
	}

	data, err := proto.Marshal(want)
	if err != nil {
		t.Fatalf("proto.Marshal: %v", err)
	}

	got, err := DecodeSnapshot(data)
	if err != nil {
		t.Fatalf("DecodeSnapshot: %v", err)
	}

	if !proto.Equal(want, got) {
		t.Fatalf("round-trip mismatch:\n want %v\n  got %v", want, got)
	}
	if got.GetBasic().GetSymbol() != "AAPL" {
		t.Errorf("Symbol = %q, want %q", got.GetBasic().GetSymbol(), "AAPL")
	}
	if got.GetPrice() != "190.12" {
		t.Errorf("Price = %q, want %q", got.GetPrice(), "190.12")
	}
	if got.GetOvnChangeRatio() != "0.0016" {
		t.Errorf("OvnChangeRatio = %q, want %q", got.GetOvnChangeRatio(), "0.0016")
	}
}

func TestDecodeQuoteRoundTrip(t *testing.T) {
	want := &Quote{
		Basic: &Basic{Symbol: "AAPL", InstrumentId: "913256135"},
		Asks: []*AskBid{
			{
				Price:  "190.20",
				Size:   "300",
				Order:  []*Order{{Mpid: "NSDQ", Size: "200"}, {Mpid: "ARCA", Size: "100"}},
				Broker: []*Broker{{Bid: "1", Name: "Broker A"}},
			},
		},
		Bids: []*AskBid{{Price: "190.10", Size: "500"}},
	}

	data, err := proto.Marshal(want)
	if err != nil {
		t.Fatalf("proto.Marshal: %v", err)
	}

	got, err := DecodeQuote(data)
	if err != nil {
		t.Fatalf("DecodeQuote: %v", err)
	}
	if got.GetBasic().GetSymbol() != "AAPL" {
		t.Errorf("Symbol = %q, want %q", got.GetBasic().GetSymbol(), "AAPL")
	}
	if len(got.GetAsks()) != 1 || got.GetAsks()[0].GetPrice() != "190.20" {
		t.Errorf("Asks = %v, want one level at 190.20", got.GetAsks())
	}
	if len(got.GetAsks()[0].GetOrder()) != 2 {
		t.Errorf("Ask orders = %d, want 2", len(got.GetAsks()[0].GetOrder()))
	}
	if len(got.GetBids()) != 1 || got.GetBids()[0].GetSize() != "500" {
		t.Errorf("Bids = %v, want one level of size 500", got.GetBids())
	}
}

func TestDecodeTickRoundTrip(t *testing.T) {
	want := &Tick{
		Basic:  &Basic{Symbol: "AAPL", InstrumentId: "913256135"},
		Time:   "1700000001000",
		Price:  "190.12",
		Volume: "100",
		Side:   "B",
	}

	data, err := proto.Marshal(want)
	if err != nil {
		t.Fatalf("proto.Marshal: %v", err)
	}

	got, err := DecodeTick(data)
	if err != nil {
		t.Fatalf("DecodeTick: %v", err)
	}
	if !proto.Equal(want, got) {
		t.Fatalf("round-trip mismatch:\n want %v\n  got %v", want, got)
	}
}

func TestDecodeSnapshotInvalid(t *testing.T) {
	if _, err := DecodeSnapshot([]byte{0xff, 0xff, 0xff, 0xff}); err == nil {
		t.Fatal("DecodeSnapshot(invalid) = nil error, want error")
	}
}
