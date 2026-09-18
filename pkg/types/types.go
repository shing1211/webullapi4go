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

// Package types defines shared Webull domain types that are safe for external
// use.
package types

// Market identifies a tradable market.
type Market string

const (
	// MarketHK is the Hong Kong market.
	MarketHK Market = "HK"
	// MarketUS is the United States market.
	MarketUS Market = "US"
	// MarketJP is the Japan market.
	MarketJP Market = "JP"
	// MarketSG is the Singapore market.
	MarketSG Market = "SG"
	// MarketCN is the mainland China market.
	MarketCN Market = "CN"
)

// String returns the market code.
func (m Market) String() string { return string(m) }

// InstrumentType identifies the kind of instrument.
type InstrumentType string

const (
	// InstrumentStock is a common stock.
	InstrumentStock InstrumentType = "STOCK"
	// InstrumentETF is an exchange-traded fund.
	InstrumentETF InstrumentType = "ETF"
	// InstrumentOption is a listed option contract.
	InstrumentOption InstrumentType = "OPTION"
	// InstrumentFuture is a futures contract.
	InstrumentFuture InstrumentType = "FUTURE"
	// InstrumentIndex is an index.
	InstrumentIndex InstrumentType = "INDEX"
	// InstrumentCrypto is a cryptocurrency pair.
	InstrumentCrypto InstrumentType = "CRYPTO"
)

// String returns the instrument type code.
func (t InstrumentType) String() string { return string(t) }
