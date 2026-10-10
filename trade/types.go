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

package trade

import "github.com/shing1211/webullapi4go/pkg/domain/money"

// Market identifies a tradable market for order routing and instrument
// lookup.
type Market string

// Markets supported by the trading API.
const (
	// MarketUS identifies the United States market.
	MarketUS Market = "US"
	// MarketHK identifies the Hong Kong market.
	MarketHK Market = "HK"
	// MarketCN identifies the mainland China market.
	MarketCN Market = "CN"
)

// InstrumentType identifies the kind of financial instrument referenced by an
// order or held in a position. It is the trading enum and differs from the
// market-data instrument classification in
// [github.com/shing1211/webullapi4go/pkg/types].
type InstrumentType string

// Instrument types accepted by the trading API.
const (
	// InstrumentTypeEquity identifies stocks and other equity instruments.
	InstrumentTypeEquity InstrumentType = "EQUITY"
	// InstrumentTypeOption identifies listed option contracts.
	InstrumentTypeOption InstrumentType = "OPTION"
	// InstrumentTypeFutures identifies futures contracts.
	//
	// The trading API spells this "FUTURES". The market-data package
	// [github.com/shing1211/webullapi4go/pkg/types] spells the equivalent
	// classification "FUTURE" (singular). The two are intentionally separate
	// domains and their wire values are not interchangeable; neither is changed
	// without evidence from the API.
	InstrumentTypeFutures InstrumentType = "FUTURES"
	// InstrumentTypeEvent identifies event contract instruments (prediction markets).
	//
	// The trading API uses this value for event contract orders. Event contracts
	// are binary-outcome instruments with yes/no sides.
	InstrumentTypeEvent InstrumentType = "EVENT"
)

// AccountType identifies how an account is funded.
type AccountType string

// Account funding types.
const (
	// AccountTypeCash identifies a cash account, which settles trades from
	// its own cash balance.
	AccountTypeCash AccountType = "CASH"
	// AccountTypeMargin identifies a margin account, which can borrow
	// against its positions.
	AccountTypeMargin AccountType = "MARGIN"
)

// AccountClass is the regulatory classification of an account.
type AccountClass string

// Account classes returned by the account-list endpoint.
const (
	// AccountClassIndividualCash is an individual cash account.
	AccountClassIndividualCash AccountClass = "INDIVIDUAL_CASH"
	// AccountClassIndividualMargin is an individual margin account.
	AccountClassIndividualMargin AccountClass = "INDIVIDUAL_MRGN"
	// AccountClassFuturesMargin is a futures margin account.
	AccountClassFuturesMargin AccountClass = "FUTURES_MRGN"
	// AccountClassInstitutionalCash is an institutional cash account.
	AccountClassInstitutionalCash AccountClass = "INSTITUTIONAL_CASH"
	// AccountClassInstitutionalMargin is an institutional margin account.
	AccountClassInstitutionalMargin AccountClass = "INSTITUTIONAL_MRGN"
	// AccountClassInstitutionalFuturesMargin is an institutional futures
	// margin account.
	AccountClassInstitutionalFuturesMargin AccountClass = "INSTITUTIONAL_FUTURES_MRGN"
)

// Account is a single brokerage account.
type Account struct {
	// AccountID is the unique account identifier used by every account-scoped
	// endpoint.
	AccountID string `json:"account_id"`
	// AccountNumber is the broker-assigned account number.
	AccountNumber string `json:"account_number"`
	// AccountType is the funding type of the account.
	AccountType AccountType `json:"account_type"`
	// AccountClass is the regulatory classification of the account.
	AccountClass AccountClass `json:"account_class"`
}

// AssetsBalance is the multi-currency asset summary for one account.
type AssetsBalance struct {
	// TotalAssetCurrency is the currency the balance totals are expressed in.
	TotalAssetCurrency string `json:"total_asset_currency"`
	// TotalCashBalance is the total cash across currencies.
	TotalCashBalance money.Money `json:"total_cash_balance"`
	// TotalMarketValue is the total market value of held positions.
	TotalMarketValue money.Money `json:"total_market_value"`
	// TotalUnrealizedProfitLoss is the total open profit or loss.
	TotalUnrealizedProfitLoss money.Money `json:"total_unrealized_profit_loss"`
	// InitMargin is the total initial margin requirement.
	InitMargin money.Money `json:"init_margin"`
	// AccountCurrencyAssets breaks the balance down by currency.
	AccountCurrencyAssets []AssetsCurrencyAssets `json:"account_currency_assets"`
}

// AssetsCurrencyAssets is the asset breakdown for a single currency.
type AssetsCurrencyAssets struct {
	// Currency is the currency this breakdown applies to.
	Currency string `json:"currency"`
	// CashBalance is the total cash balance in Currency.
	CashBalance money.Money `json:"cash_balance"`
	// SettledCash is the settled portion of the cash balance.
	SettledCash money.Money `json:"settled_cash"`
	// UnsettledCash is the portion of the cash balance still settling.
	UnsettledCash money.Money `json:"unsettled_cash"`
	// MarketValue is the market value of positions held in Currency.
	MarketValue money.Money `json:"market_value"`
	// HeldAmount is the funds committed to open orders.
	HeldAmount money.Money `json:"held_amount"`
	// FrozenAmount is the funds frozen by the broker.
	FrozenAmount money.Money `json:"frozen_amount"`
	// BuyingPower is the funds available to open new positions.
	BuyingPower money.Money `json:"buying_power"`
	// UnrealizedProfitLoss is the open profit or loss in Currency.
	UnrealizedProfitLoss money.Money `json:"unrealized_profit_loss"`
	// AvailableWithdrawal is the amount available to withdraw.
	AvailableWithdrawal money.Money `json:"available_withdrawal"`
	// InterestsUnpaid is the accrued interest not yet paid.
	InterestsUnpaid money.Money `json:"interests_unpaid"`
	// InitMargin is the initial margin requirement in Currency.
	InitMargin money.Money `json:"init_margin"`
}

// OptionStrategy identifies the structure of an options position.
type OptionStrategy string

// Option strategy values.
const (
	// OptionStrategySingle is a single-leg options position.
	OptionStrategySingle OptionStrategy = "SINGLE"
	// OptionStrategyVertical is a vertical spread: long and short options of
	// the same type and expiration at different strikes.
	OptionStrategyVertical OptionStrategy = "VERTICAL"
	// OptionStrategyStraddle is a straddle: a call and a put at the same
	// strike and expiration on the same side.
	OptionStrategyStraddle OptionStrategy = "STRADDLE"
	// OptionStrategyStrangle is a strangle: a call and a put on the same side
	// at different strikes and the same expiration.
	OptionStrategyStrangle OptionStrategy = "STRANGLE"
	// OptionStrategyIronCondor is an iron condor: a short strangle bracketed by
	// a wider long strangle.
	OptionStrategyIronCondor OptionStrategy = "IRON_CONDOR"
	// OptionStrategyIronButterfly is an iron butterfly: a short straddle
	// bracketed by a long strangle.
	OptionStrategyIronButterfly OptionStrategy = "IRON_BUTTERFLY"
	// OptionStrategyButterfly is a butterfly spread: a long option at a low
	// strike, two short options at a middle strike, and a long option at a high
	// strike, all of the same type and expiration.
	OptionStrategyButterfly OptionStrategy = "BUTTERFLY"
	// OptionStrategyCollar is a collar: a long put financed by a short call on
	// the same underlying.
	OptionStrategyCollar OptionStrategy = "COLLAR"
	// OptionStrategyCalendar is a calendar spread: options of the same type
	// and strike at different expirations.
	OptionStrategyCalendar OptionStrategy = "CALENDAR"
	// OptionStrategyDiagonal is a diagonal spread: options of the same type at
	// different strikes and different expirations.
	OptionStrategyDiagonal OptionStrategy = "DIAGONAL"
	// OptionStrategyRatio is a ratio spread: an unequal number of long and
	// short options of the same type.
	OptionStrategyRatio OptionStrategy = "RATIO"
)

// EventOutcome identifies the side of an event contract position.
type EventOutcome string

// Event outcomes.
const (
	// EventOutcomeYes represents the "yes" side of an event contract.
	EventOutcomeYes EventOutcome = "yes"
	// EventOutcomeNo represents the "no" side of an event contract.
	EventOutcomeNo EventOutcome = "no"
)

// OptionType identifies whether an options leg is a call or a put.
type OptionType string

// Option types.
const (
	// OptionTypeCall is a call option: the right to buy the underlying.
	OptionTypeCall OptionType = "CALL"
	// OptionTypePut is a put option: the right to sell the underlying.
	OptionTypePut OptionType = "PUT"
)

// Position is a single holding in an account.
type Position struct {
	// PositionID is the unique position identifier.
	PositionID string `json:"position_id"`
	// Currency is the currency the position is denominated in.
	Currency string `json:"currency"`
	// Quantity is the held quantity, as a decimal string.
	Quantity money.Money `json:"quantity"`
	// Symbol is the trading symbol of the held instrument.
	Symbol string `json:"symbol"`
	// OptionStrategy is the options strategy, empty for non-option
	// positions.
	OptionStrategy OptionStrategy `json:"option_strategy"`
	// InstrumentType is the kind of instrument held.
	InstrumentType InstrumentType `json:"instrument_type"`
	// LastPrice is the latest market price, as a decimal string.
	LastPrice money.Money `json:"last_price"`
	// CostPrice is the average cost basis, as a decimal string.
	CostPrice money.Money `json:"cost_price"`
	// UnrealizedProfitLoss is the open profit or loss, as a decimal string.
	UnrealizedProfitLoss money.Money `json:"unrealized_profit_loss"`
	// Legs lists the legs of a multi-leg option position. It is empty for
	// single-leg and non-option positions.
	Legs []PositionLeg `json:"legs"`
}

// PositionLeg is one leg of a multi-leg option position.
type PositionLeg struct {
	// Symbol is the trading symbol of the leg.
	Symbol string `json:"symbol"`
	// Quantity is the leg quantity, as a decimal string.
	Quantity money.Money `json:"quantity"`
	// OptionType is whether the leg is a call or a put.
	OptionType OptionType `json:"option_type"`
	// OptionExpireDate is the leg expiration date in yyyy-MM-dd form.
	OptionExpireDate string `json:"option_expire_date"`
	// OptionExercisePrice is the strike price, as a decimal string.
	OptionExercisePrice money.Money `json:"option_exercise_price"`
	// OptionContractMultiplier is the number of shares one contract
	// represents, as a decimal string.
	OptionContractMultiplier money.Money `json:"option_contract_multiplier"`
	// OptionContractDeliverable is the number of shares delivered on
	// exercise of one contract, as a decimal string.
	OptionContractDeliverable money.Money `json:"option_contract_deliverable"`
	// ExpirationType is the option expiration style, for example "AM".
	ExpirationType string `json:"expiration_type"`
}
