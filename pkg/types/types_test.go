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

package types_test

import (
	"testing"

	"github.com/shing1211/webullapi4go/pkg/types"
)

func TestMarketString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		market types.Market
		want   string
	}{
		{name: "Hong Kong", market: types.MarketHK, want: "HK"},
		{name: "United States", market: types.MarketUS, want: "US"},
		{name: "Japan", market: types.MarketJP, want: "JP"},
		{name: "Singapore", market: types.MarketSG, want: "SG"},
		{name: "China", market: types.MarketCN, want: "CN"},
		{name: "zero", market: types.Market(""), want: ""},
		{name: "unknown", market: types.Market("EU"), want: "EU"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.market.String(); got != tt.want {
				t.Fatalf("Market(%q).String() = %q, want %q", tt.market, got, tt.want)
			}
		})
	}
}

func TestInstrumentTypeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		instrumentType types.InstrumentType
		want           string
	}{
		{name: "stock", instrumentType: types.InstrumentStock, want: "STOCK"},
		{name: "ETF", instrumentType: types.InstrumentETF, want: "ETF"},
		{name: "option", instrumentType: types.InstrumentOption, want: "OPTION"},
		{name: "future", instrumentType: types.InstrumentFuture, want: "FUTURE"},
		{name: "index", instrumentType: types.InstrumentIndex, want: "INDEX"},
		{name: "crypto", instrumentType: types.InstrumentCrypto, want: "CRYPTO"},
		{name: "zero", instrumentType: types.InstrumentType(""), want: ""},
		{name: "unknown", instrumentType: types.InstrumentType("BOND"), want: "BOND"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.instrumentType.String(); got != tt.want {
				t.Fatalf("InstrumentType(%q).String() = %q, want %q", tt.instrumentType, got, tt.want)
			}
		})
	}
}
