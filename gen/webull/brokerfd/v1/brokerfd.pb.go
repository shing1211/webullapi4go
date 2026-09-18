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

package brokerfdv1

type BrokerFDEventType int

const (
	BrokerFDEventTypeAccountPush  BrokerFDEventType = 0
	BrokerFDEventTypeOrderPush    BrokerFDEventType = 1
	BrokerFDEventTypePositionPush BrokerFDEventType = 2
	BrokerFDEventTypeTradePush    BrokerFDEventType = 3
	BrokerFDEventTypeAssetDetail  BrokerFDEventType = 4
	BrokerFDEventTypeRiskPush     BrokerFDEventType = 5
	BrokerFDEventTypeOrderFill    BrokerFDEventType = 6
	BrokerFDEventTypePositionSync BrokerFDEventType = 7
	BrokerFDEventTypeAssetSync    BrokerFDEventType = 8
	BrokerFDEventTypeAssetsPush   BrokerFDEventType = 9
	BrokerFDEventTypeOrdersPush   BrokerFDEventType = 10
	BrokerFDEventTypeCashSecLiab  BrokerFDEventType = 11
)

func (BrokerFDEventType) EnumDescriptor() ([]byte, []int) { return nil, nil }
func (x BrokerFDEventType) String() string {
	switch x {
	case BrokerFDEventTypeAccountPush:
		return "AccountPush"
	case BrokerFDEventTypeOrderPush:
		return "OrderPush"
	case BrokerFDEventTypePositionPush:
		return "PositionPush"
	case BrokerFDEventTypeTradePush:
		return "TradePush"
	case BrokerFDEventTypeAssetDetail:
		return "AssetDetail"
	case BrokerFDEventTypeRiskPush:
		return "RiskPush"
	case BrokerFDEventTypeOrderFill:
		return "OrderFill"
	case BrokerFDEventTypePositionSync:
		return "PositionSync"
	case BrokerFDEventTypeAssetSync:
		return "AssetSync"
	case BrokerFDEventTypeAssetsPush:
		return "AssetsPush"
	case BrokerFDEventTypeOrdersPush:
		return "OrdersPush"
	case BrokerFDEventTypeCashSecLiab:
		return "CashSecLiab"
	default:
		return "Unknown"
	}
}

type BrokerFDEvent struct {
	EventType BrokerFDEventType
	Data      []byte
}

type AccountPush struct {
	AccountID    string
	Currency     string
	NetLiquidity string
	CashBalance  string
	MarketValue  string
}

type OrderPush struct {
	OrderID   string
	Symbol    string
	Side      string
	Status    string
	Quantity  string
	FilledQty string
}

type PositionPush struct {
	Symbol       string
	Quantity     string
	MarketValue  string
	CostBasis    string
	UnrealizedPL string
}

type TradePush struct {
	TradeID   string
	OrderID   string
	Symbol    string
	Side      string
	Price     string
	Quantity  string
	TradeTime string
}

type AssetDetail struct {
	AccountID   string
	TotalAssets string
	CashBalance string
	MarketValue string
	NetLiq      string
	Currency    string
}

type RiskPush struct {
	AccountID     string
	BuyingPower   string
	MarginEquity  string
	RegTMarReq    string
	DayTradesLeft int
}

type OrderFillPush struct {
	FillId     string
	OrderId    string
	Symbol     string
	Side       string
	Price      string
	Quantity   string
	Commission string
}

type PositionSync struct {
	Positions []*PositionPush
}

type AssetSync struct {
	Assets []*AssetDetail
}

type AssetsPush struct {
	Assets []*AssetDetail
}

type OrdersPush struct {
	Orders []*OrderPush
}

type CashSecLiability struct {
	AccountID     string
	CashBalance   string
	SecuritiesLia string
	Currency      string
}
