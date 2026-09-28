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

package conformance

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/shing1211/webullapi4go/brokerfd"
	"github.com/shing1211/webullapi4go/connect"
	"github.com/shing1211/webullapi4go/data"
	"github.com/shing1211/webullapi4go/display"
	"github.com/shing1211/webullapi4go/internal/auth"
	"github.com/shing1211/webullapi4go/pkg/domain/order"
	"github.com/shing1211/webullapi4go/trade"
)

// This file is the harness's weak point, and it is written out longhand rather
// than derived so that a reader can audit it. Four things follow from that.
//
//  1. The mapping is not derivable by reflection. A manifest row names an SDK
//     method -- brokerfd.GetFDPositions is a method, not a type -- and nothing
//     in the Go type system connects a method to the type it decodes into. The
//     method body says `var out FDPosition`, and reading method bodies is not
//     something a reflection-based engine can do. So the table is the
//     association, spelled out, and it is reviewable in a way a
//     string-to-string search of source would not be.
//
//  2. Every manifest symbol must appear. A symbol with no entry is reported as
//     unmapped, by name, with a count -- never skipped. A silently skipped
//     symbol is indistinguishable from a passing one, which is the exact
//     failure mode the whole exercise exists to remove. The table is therefore
//     asserted against the manifest in both directions by a test, and
//     Unmapped reports the gap on demand.
//
//  3. The type named is the one the method's `var out` declares, not the one it
//     returns. Most methods return *T after decoding into a T, so the JSON
//     target is *T; naming T makes Describe transparent to pointers and makes
//     the decode target reflect.New(T) exactly what the SDK hands its transport.
//     Where a method wraps the decoded value, the wrapped type is named and the
//     wrap is recorded in the entry's Note. Naming the *returned* type instead is
//     the mistake this rule exists to prevent, and it is easy to make twice over:
//     a method that returns `[]T` because it decodes an envelope and hands back
//     the envelope's slice, and a method that returns a public projection built
//     from an unexported envelope. Both were in the table; TestEnvelopeTypesMatch
//     TheSource and the gate both notice, because a wrong name reports a shape
//     inversion against a page the SDK never claimed to implement.
//
//  4. A nil type is a fact, not a gap. Five methods return only an error, or raw
//     bytes, and decode no body at all. They are in the table with Type nil and
//     a Note saying so, which is different from being absent: the shape checks
//     are then reported as not applicable with a stated reason instead of
//     quietly passing.
//
//  5. Six envelopes are unexported, so an entry can name one only as a spelling.
//     Envelope holds "package.TypeName" and is read from the SDK sources by
//     envelopes.go, because a type another package does not export cannot be
//     named in a reflect.TypeOf call at all. An entry names one or the other,
//     never both; DecodeTarget refuses a pair.
//
// # Symbols this table cannot reach
//
// The 29 broker.* symbols are absent, and that is a module boundary rather than
// a gap in anyone's attention: broker/ is a separate Go module, so a type in it
// is not importable from the root module's conformance package without a
// go.mod change. Every other manifest symbol is in the root module and is in
// the table below. Separately, the ten rows the docgen manifest maps to an em
// dash name no SDK symbol at all, so they are neither mapped nor unmapped --
// see NoSDKSymbol.

// SDKType is one row of the table: what an SDK symbol decodes a documented
// response into.
type SDKType struct {
	// Type is the type the method's `var out` declares, or nil when the method
	// decodes no body.
	Type reflect.Type
	// Envelope is the "package.TypeName" spelling of a struct the SDK package
	// declares unexported, for the case Type cannot hold. It is empty whenever
	// Type is set, and it is read from the SDK sources; see envelopes.go.
	Envelope string
	// Note records anything a reader could not derive from Type alone: a wrap,
	// a deliberate second entry point, or why Type is nil.
	Note string
}

// DecodeTarget is the reflect.Type the symbol's method hands its transport, which
// is the type the five checks compare a documented body against.
//
// An error here is a defect in this table rather than in the SDK, and it is
// deliberately an error and not a nil: a nil type would be reported by the
// comparison as "the method decodes no response body", which is a different and
// untrue statement.
func (s SDKType) DecodeTarget() (reflect.Type, error) {
	switch {
	case s.Type != nil && s.Envelope != "":
		return nil, fmt.Errorf("entry names both %s and the unexported envelope %s; "+
			"a decode target is one type or the other", s.Type, s.Envelope)
	case s.Envelope != "":
		return envelopeType(s.Envelope)
	default:
		return s.Type, nil
	}
}

// SDKTypes is the symbol table, keyed by the manifest's symbol spelling.
//
// Keys are single symbols. Two docgen rows list two symbols separated by
// " / "; Subject picks the first, and NormaliseSymbols expands the rest.
var SDKTypes = map[string]SDKType{
	// -- client: token lifecycle ------------------------------------------
	// The wire DTO is internal/auth.Token, not the public client.Token. The
	// public type is a projection with no json tags at all, so naming it would
	// report three phantom missing names on every token endpoint; client.Token
	// is built from auth.Token in newToken and never touches the wire.
	"client.CreateToken": {Type: reflect.TypeOf(auth.Token{}), Note: "decodes auth.Token; the public client.Token is a tagless projection of it"},
	"client.CheckToken":  {Type: reflect.TypeOf(auth.Token{}), Note: "decodes auth.Token; the public client.Token is a tagless projection of it"},

	// -- connect: OAuth 2.0 authorization-code exchange --------------------
	"connect.CreateToken": {Type: reflect.TypeOf(connect.Token{})},

	// -- display: Hosted Display Solution client tokens --------------------
	"display.Service.EnsureToken": {
		Type: reflect.TypeOf(auth.DisplayToken{}),
		Note: "decodes auth.DisplayToken and returns only its access-token value; display.ClientToken is the identical public shape used by the refresh path",
	},
	"display.Service.RefreshClientToken": {Type: reflect.TypeOf(display.ClientToken{})},

	// -- brokerfd: Broker FD US -------------------------------------------
	"brokerfd.AddFDAchAccount":           {Type: reflect.TypeOf(brokerfd.ACHAccount{})},
	"brokerfd.AddFDBankAccount":          {Type: reflect.TypeOf(brokerfd.BankAccount{})},
	"brokerfd.CancelFDOrder":             {Type: nil, Note: "returns only error; the documented body is discarded"},
	"brokerfd.CloseFDAccount":            {Type: nil, Note: "returns only error; the documented body is discarded"},
	"brokerfd.CreateFDAccount":           {Type: reflect.TypeOf(brokerfd.FDAccount{})},
	"brokerfd.CreateFDInstantFunding":    {Type: reflect.TypeOf(brokerfd.InstantFunding{})},
	"brokerfd.DownloadDocument":          {Type: nil, Note: "returns the raw []byte body, not a typed DTO"},
	"brokerfd.GetAgreementDetail":        {Type: reflect.TypeOf(brokerfd.Agreement{})},
	"brokerfd.GetFDAccountDetail":        {Type: reflect.TypeOf(brokerfd.FDAccount{})},
	"brokerfd.GetFDActivities":           {Type: reflect.TypeOf([]brokerfd.FDActivity{})},
	"brokerfd.GetFDAssetsDetail":         {Type: reflect.TypeOf(brokerfd.FDAssetsDetail{})},
	"brokerfd.GetFDCashJournalDetail":    {Type: reflect.TypeOf(brokerfd.FDCashJournal{})},
	"brokerfd.GetFDCorporateActions":     {Type: reflect.TypeOf([]brokerfd.FDCorporateAction{})},
	"brokerfd.GetFDCreditInfo":           {Type: reflect.TypeOf(brokerfd.CreditInfo{})},
	"brokerfd.GetFDECInstruments":        {Type: reflect.TypeOf([]brokerfd.FDECInstrument{})},
	"brokerfd.GetFDEnums":                {Type: reflect.TypeOf([]brokerfd.FDEnum{})},
	"brokerfd.GetFDInstantFundingDetail": {Type: reflect.TypeOf(brokerfd.InstantFunding{})},
	"brokerfd.GetFDOpenOrders":           {Type: reflect.TypeOf([]brokerfd.FDOrder{})},
	"brokerfd.GetFDOrderDetail":          {Type: reflect.TypeOf(brokerfd.FDOrder{})},
	"brokerfd.GetFDOrderHistory":         {Type: reflect.TypeOf([]brokerfd.FDOrder{})},
	"brokerfd.GetFDStockInstruments":     {Type: reflect.TypeOf([]brokerfd.FDStockInstrument{})},
	"brokerfd.GetFDTradeCalendar":        {Type: reflect.TypeOf([]brokerfd.FDTradeCalendarEntry{})},
	"brokerfd.GetFDTransferDetail":       {Type: reflect.TypeOf(brokerfd.Transfer{})},
	"brokerfd.GetFDTransferFees":         {Type: reflect.TypeOf([]brokerfd.TransferFee{})},
	"brokerfd.InitiateFDTransfer":        {Type: reflect.TypeOf(brokerfd.Transfer{})},
	"brokerfd.ListAccountForms":          {Type: reflect.TypeOf([]brokerfd.AccountForm{})},
	"brokerfd.ListAgreements":            {Type: reflect.TypeOf([]brokerfd.Agreement{})},
	"brokerfd.ListFDAccounts":            {Type: reflect.TypeOf([]brokerfd.FDAccount{})},
	"brokerfd.ListFDAchAccounts":         {Type: reflect.TypeOf([]brokerfd.ACHAccount{})},
	"brokerfd.ListFDBankAccounts":        {Type: reflect.TypeOf([]brokerfd.BankAccount{})},
	"brokerfd.ListFDTransfers":           {Type: reflect.TypeOf([]brokerfd.Transfer{})},
	"brokerfd.PlaceFDOrder":              {Type: reflect.TypeOf(brokerfd.FDOrder{})},
	"brokerfd.PreviewFDOrder":            {Type: reflect.TypeOf(brokerfd.FDOrderPreview{})},
	"brokerfd.RemoveFDAchAccount":        {Type: nil, Note: "returns only error; the documented body is discarded"},
	"brokerfd.RemoveFDBankAccount":       {Type: nil, Note: "returns only error; the documented body is discarded"},
	"brokerfd.ReplaceFDOrder":            {Type: reflect.TypeOf(brokerfd.FDOrder{})},
	"brokerfd.UpdateFDAccount":           {Type: reflect.TypeOf(brokerfd.FDAccount{})},
	"brokerfd.UploadDocument":            {Type: reflect.TypeOf(brokerfd.Document{})},
	"brokerfd.GetAccountsSummary": {
		Type: reflect.TypeOf([]brokerfd.AccountSummary{}),
		Note: "first of two symbols the docgen row lists; the alternative GetFDAssetsSummary decodes a flat *FDAssetsSummary instead",
	},
	"brokerfd.GetFDPositions": {
		Type: reflect.TypeOf([]brokerfd.FDPosition{}),
		Note: "first of two symbols the docgen row lists; the alternative GetPositions decodes []PositionSummary and sends an undocumented path",
	},
	"brokerfd.GetFDAssetsSummary": {
		Type: reflect.TypeOf(brokerfd.FDAssetsSummary{}),
		Note: "alternative on the two-symbol summaries row; not compared, because the row's subject decodes the same page as an array",
	},
	"brokerfd.GetPositions": {
		Type: reflect.TypeOf([]brokerfd.PositionSummary{}),
		Note: "alternative on the two-symbol positions row; not compared, because it sends an undocumented path",
	},

	// -- trade: Trading API ------------------------------------------------
	"trade.BatchPlaceOrder": {Type: reflect.TypeOf(trade.BatchPlaceOrderResponse{})},
	"trade.CancelOrder":     {Type: reflect.TypeOf(trade.CancelOrderResult{})},
	"trade.GetBalance":      {Type: reflect.TypeOf(trade.AssetsBalance{})},
	"trade.GetCashActivities": {
		Type: reflect.TypeOf(trade.CashActivityPage{}),
		Note: "decodes trade.CashActivityPage in GetCashActivitiesPage (trade/accounts.go:177) and " +
			"returns its Activities slice; the page documents the {data, pagination_key} envelope, " +
			"so naming the returned []trade.CashActivity reported an inversion the SDK does not have",
	},
	"trade.GetOpenOrders": {
		Type: reflect.TypeOf(trade.OrderPage{}),
		Note: "decodes trade.OrderPage in GetOpenOrdersPage (trade/query.go:218) and returns its " +
			"Orders slice; see trade.GetCashActivities for the same shape",
	},
	"trade.GetOrderDetail": {Type: reflect.TypeOf(trade.OrderGroup{})},
	"trade.GetOrderHistory": {
		Type: reflect.TypeOf(trade.OrderPage{}),
		Note: "decodes trade.OrderPage in GetOrderHistoryPage (trade/query.go:275) and returns its " +
			"Orders slice; see trade.GetCashActivities for the same shape",
	},
	"trade.GetPositions": {Type: reflect.TypeOf([]trade.Position{})},
	"trade.ListAccounts": {Type: reflect.TypeOf([]trade.Account{})},
	"trade.PreviewOrder": {Type: reflect.TypeOf(trade.PreviewResult{})},
	"trade.ReplaceOrder": {Type: reflect.TypeOf(trade.ReplaceOrderResult{})},
	"trade.PlaceOrder": {
		Type: reflect.TypeOf(order.PlaceOrderResult{}),
		Note: "decodes order.PlaceOrderResult and wraps it in *order.Order together with machine state that never came off the wire",
	},

	// -- data: market data, fundamentals, instruments, watchlists ----------
	//
	// The five watchlist mutations below decode data.BoolOrSuccess, not
	// data.SuccessResponse. BoolOrSuccess brings its own UnmarshalJSON that takes
	// either a bare JSON boolean or a {"success": bool} object, and SuccessResponse
	// is the tagless public projection each method then builds. Naming the
	// projection compared a type the SDK never decodes against the page's body.
	"data.AddWatchlistInstruments": {Type: reflect.TypeOf(data.BoolOrSuccess{})},
	"data.CreateWatchlist":         {Type: reflect.TypeOf(data.CreateWatchlistResult{})},
	"data.DeleteWatchlist":         {Type: reflect.TypeOf(data.BoolOrSuccess{})},
	"data.GetAnalystRating":        {Type: reflect.TypeOf(data.AnalystRating{})},
	"data.GetAnalystTargetPrice":   {Type: reflect.TypeOf(data.AnalystTargetPrice{})},
	"data.GetBalanceSheet":         {Type: reflect.TypeOf([]data.FinancialsItem{})},
	"data.GetBatchBars":            {Type: reflect.TypeOf(data.BatchBars{})},
	"data.GetCapitalFlow":          {Type: reflect.TypeOf([]data.CapitalFlowEntry{})},
	"data.GetCashFlow":             {Type: reflect.TypeOf([]data.FinancialsItem{})},
	"data.GetCompanyProfile":       {Type: reflect.TypeOf(data.CompanyProfile{})},
	"data.GetCorporateActions": {
		Envelope: "data.corpActionResponse",
		Note: "decodes data.corpActionResponse (data/corporate_actions.go:44) and returns its " +
			"Data slice; naming the returned []data.CorporateAction reported an object-vs-array " +
			"inversion the SDK does not have",
	},
	"data.GetCorporateActionsByMarket": {
		Envelope: "data.corpActionResponse",
		Note: "decodes the same data.corpActionResponse (data/corporate_actions.go:131) as " +
			"data.GetCorporateActions, from the list-by-market path",
	},
	"data.GetCryptoBars": {Type: reflect.TypeOf([]data.CryptoSymbolBars{})},
	"data.GetCryptoInstruments": {
		Envelope: "data.cryptoInstrumentsResponse",
		Note: "decodes data.cryptoInstrumentsResponse (data/crypto.go:154) and projects it into " +
			"the tagless public data.CryptoInstrumentsResult",
	},
	"data.GetCryptoSnapshot":             {Type: reflect.TypeOf([]data.CryptoSnapshot{})},
	"data.GetDSAnalystRating":            {Type: reflect.TypeOf(data.AnalystRating{})},
	"data.GetDSAnalystTargetPrice":       {Type: reflect.TypeOf(data.AnalystTargetPrice{})},
	"data.GetDSCompanyProfile":           {Type: reflect.TypeOf(data.CompanyProfile{})},
	"data.GetDSLatestNews":               {Type: reflect.TypeOf([]data.DSNewsSummaryItem{})},
	"data.GetDSMarketNews":               {Type: reflect.TypeOf([]data.DSNewsSummaryItem{})},
	"data.GetDSNewsSummary":              {Type: reflect.TypeOf([]data.DSNewsSummaryItem{})},
	"data.GetDSSymbolNews":               {Type: reflect.TypeOf([]data.DSNewsSummaryItem{})},
	"data.GetDisplayBars":                {Type: reflect.TypeOf(data.BatchBars{})},
	"data.GetDisplayBarsSingle":          {Type: reflect.TypeOf(data.StockBars{})},
	"data.GetDisplayDepth":               {Type: reflect.TypeOf(data.Quote{})},
	"data.GetDisplayGainersLosers":       {Type: reflect.TypeOf([]data.ScreenerStock{})},
	"data.GetDisplaySnapshot":            {Type: reflect.TypeOf([]data.Snapshot{})},
	"data.GetDisplayTick":                {Type: reflect.TypeOf(data.StockTicks{})},
	"data.GetDisplayTopActive":           {Type: reflect.TypeOf([]data.ScreenerStock{})},
	"data.GetDividendCalendar":           {Type: reflect.TypeOf([]data.DividendCalendarEntry{})},
	"data.GetEarningsCalendar":           {Type: reflect.TypeOf([]data.EarningsCalendarEntry{})},
	"data.GetEventBars":                  {Type: reflect.TypeOf([]data.EventBarsResult{})},
	"data.GetEventContractCategories":    {Type: reflect.TypeOf([]data.EventContractCategory{})},
	"data.GetEventContractEvents":        {Type: reflect.TypeOf([]data.EventContractEvent{})},
	"data.GetEventContractEventsList":    {Type: reflect.TypeOf(data.EventListPage{})},
	"data.GetEventContractMarkets":       {Type: reflect.TypeOf([]data.EventContractMarket{})},
	"data.GetEventContractMilestones":    {Type: reflect.TypeOf(data.MilestonePage{})},
	"data.GetEventContractSeries":        {Type: reflect.TypeOf([]data.EventContractSeries{})},
	"data.GetEventContractSeriesList":    {Type: reflect.TypeOf(data.EventSeriesListPage{})},
	"data.GetEventContractSportsFilters": {Type: reflect.TypeOf([]data.EventSportsFilter{})},
	"data.GetEventContractTags":          {Type: reflect.TypeOf([]data.EventContractTag{})},
	"data.GetEventDepth":                 {Type: reflect.TypeOf([]data.EventDepth{})},
	"data.GetEventGameStats":             {Type: reflect.TypeOf(data.EventGameStats{})},
	"data.GetEventLiveData":              {Type: reflect.TypeOf(data.EventLiveData{})},
	"data.GetEventMarketBars":            {Type: reflect.TypeOf([]data.EventMarketBar{})},
	"data.GetEventMarketBarsByEvent":     {Type: reflect.TypeOf([]data.EventMarketBar{})},
	"data.GetEventMarketDepth":           {Type: reflect.TypeOf(data.EventMarketDepth{})},
	"data.GetEventMarketSnapshot":        {Type: reflect.TypeOf(data.EventMarketSnapshot{})},
	"data.GetEventSnapshot":              {Type: reflect.TypeOf([]data.EventSnapshot{})},
	"data.GetEventTick":                  {Type: reflect.TypeOf([]data.EventTickResult{})},
	"data.GetFilings":                    {Type: reflect.TypeOf(data.FilingsResponse{})},
	"data.GetFinancialAlert":             {Type: reflect.TypeOf(data.FinancialAlert{})},
	"data.GetFinancialIndicators":        {Type: reflect.TypeOf(data.FinancialIndicator{})},
	"data.GetFootprint":                  {Type: reflect.TypeOf([]data.StockFootprint{})},
	"data.GetForecastEPS":                {Type: reflect.TypeOf([]data.ForecastEPSEntry{})},
	"data.GetFundAllocation":             {Type: reflect.TypeOf([]data.FundAllocation{})},
	"data.GetFundDividends":              {Type: reflect.TypeOf([]data.FundDividend{})},
	"data.GetFundFiles":                  {Type: reflect.TypeOf([]data.FundFile{})},
	"data.GetFundHoldings":               {Type: reflect.TypeOf([]data.FundHolding{})},
	"data.GetFundInfo":                   {Type: reflect.TypeOf(data.FundInfo{})},
	"data.GetFundNav":                    {Type: reflect.TypeOf([]data.FundNav{})},
	"data.GetFundPerformance":            {Type: reflect.TypeOf(data.FundPerformance{})},
	"data.GetFundRating":                 {Type: reflect.TypeOf([]data.FundRating{})},
	"data.GetFundSplits":                 {Type: reflect.TypeOf([]data.FundSplit{})},
	"data.GetFuturesBars":                {Type: reflect.TypeOf(data.BatchBars{})},
	"data.GetFuturesDepth":               {Type: reflect.TypeOf(data.Quote{})},
	"data.GetFuturesFootprint":           {Type: reflect.TypeOf([]data.StockFootprint{})},
	"data.GetFuturesInstruments":         {Type: reflect.TypeOf([]data.FuturesInstrument{})},
	"data.GetFuturesProductClasses":      {Type: reflect.TypeOf([]data.FuturesProductClass{})},
	"data.GetFuturesProductCodes":        {Type: reflect.TypeOf([]data.FuturesProduct{})},
	"data.GetFuturesSnapshot":            {Type: reflect.TypeOf([]data.Snapshot{})},
	"data.GetFuturesTick":                {Type: reflect.TypeOf(data.StockTicks{})},
	"data.GetHighDividendRank":           {Type: reflect.TypeOf([]data.ScreenerStock{})},
	"data.GetIncomeStatement":            {Type: reflect.TypeOf([]data.FinancialsItem{})},
	"data.GetIndustryComparison":         {Type: reflect.TypeOf(data.IndustryComparison{})},
	"data.GetLogos":                      {Type: reflect.TypeOf([]data.Logo{})},
	"data.GetMarketSectorDetail":         {Type: reflect.TypeOf([]data.ScreenerStock{})},
	"data.GetMarketSectors":              {Type: reflect.TypeOf([]data.MarketSector{})},
	"data.GetMostActive":                 {Type: reflect.TypeOf([]data.ScreenerStock{})},
	"data.GetNOIIBars":                   {Type: reflect.TypeOf([]data.NOIIBar{})},
	"data.GetNOIISnapshot":               {Type: reflect.TypeOf(data.NOIISnapshot{})},
	"data.GetOptionBars": {
		Envelope: "data.optionBarsResponse",
		Note: "decodes data.optionBarsResponse (data/options.go:231) and returns its Result slice; " +
			"the page requires the top-level `result` name, which only the envelope carries",
	},
	"data.GetOptionContracts": {
		Envelope: "data.optionContractsResponse",
		Note: "decodes data.optionContractsResponse (data/options.go:367) and projects it into " +
			"the tagless public data.OptionContractsResult",
	},
	"data.GetOptionSnapshot":   {Type: reflect.TypeOf([]data.OptionSnapshot{})},
	"data.GetOptionTick":       {Type: reflect.TypeOf(data.OptionTickResult{})},
	"data.GetQuotes":           {Type: reflect.TypeOf(data.Quote{})},
	"data.GetSnapshot":         {Type: reflect.TypeOf([]data.Snapshot{})},
	"data.GetStockInstruments": {Type: reflect.TypeOf([]data.StockInstrument{})},
	"data.GetStockProfilesV3": {
		Envelope: "data.instrumentProfilesV3Resp",
		Note: "decodes data.instrumentProfilesV3Resp (data/instrument_v3.go:42) and projects it into " +
			"the tagless public data.StockProfilesV3Result. The remaining array-vs-object " +
			"disagreement on this row is the SDK's own, not the harness's: the method's " +
			"comment says the endpoint returns a {data, pagination_key} envelope while the page " +
			"documents a bare array",
	},
	"data.GetTick":                    {Type: reflect.TypeOf(data.StockTicks{})},
	"data.GetTopGainersLosers":        {Type: reflect.TypeOf([]data.ScreenerStock{})},
	"data.GetWatchlistInstruments":    {Type: reflect.TypeOf(data.WatchlistInstruments{})},
	"data.GetWatchlists":              {Type: reflect.TypeOf([]data.Watchlist{})},
	"data.GetWeek52HighLow":           {Type: reflect.TypeOf([]data.ScreenerStock{})},
	"data.RemoveWatchlistInstruments": {Type: reflect.TypeOf(data.BoolOrSuccess{})},
	"data.UpdateWatchlist":            {Type: reflect.TypeOf(data.BoolOrSuccess{})},
	"data.UpdateWatchlistInstruments": {Type: reflect.TypeOf(data.BoolOrSuccess{})},
}

// symbolSeparator is how the docgen manifest joins the two symbols some rows
// list.
const symbolSeparator = " / "

// Subject returns the symbol a manifest row is compared under, and the
// alternatives it names.
//
// The first listed symbol is the subject, on the docgen manifest's authority: it
// names the symbol the endpoint's path belongs to, so an alternative left behind
// by a rename cannot displace it. A bare alternative -- "GetPositions" beside
// "brokerfd.GetFDPositions" -- is normalised to its primary's package, which is
// how the docgen manifest itself reads that row.
//
// That authority is thinner than it looks, and the reason is worth recording.
// The order is a documentation convention, not a fact derivable from the SDK, so
// it picks correctly where the first-listed symbol is the one that sends the
// documented path and picks wrongly where it is not -- on the Broker FD
// summaries row it selects brokerfd.GetAccountsSummary, whose alternatives send
// an undocumented /broker-fd/* path, so the page is compared against a call the
// page does not describe. A rule like that cannot be repaired by choosing
// differently, because nothing in the manifest says which symbol is which.
//
// What bounds the damage is the comparability rule in shapes.go: a row whose
// subject does not send the page's path is reported not comparable and claims
// nothing. That is why the summaries row is safe to leave as it is, and why
// changing this function would trade a recorded coverage limit for a second
// convention of the same authority and no more.
func Subject(manifestSymbols string) (subject string, alternatives []string, named bool) {
	parts := strings.Split(manifestSymbols, symbolSeparator)
	if len(parts) == 0 {
		return "", nil, false
	}
	primary := strings.TrimSpace(parts[0])
	if primary == "" {
		return "", nil, false
	}
	pkg, _, qualified := strings.Cut(primary, ".")
	for _, raw := range parts[1:] {
		alt := strings.TrimSpace(raw)
		if alt == "" {
			continue
		}
		if !strings.Contains(alt, ".") && qualified {
			alt = pkg + "." + alt
		}
		alternatives = append(alternatives, alt)
	}
	return primary, alternatives, true
}

// NoSDKSymbol is the docgen AREAS manifest's em dash: an endpoint deliberately
// mapped to no SDK symbol. It is neither mapped nor unmapped, because there is
// no symbol to map, and calling it unmapped would overstate what is missing.
const NoSDKSymbol = "—"

// outOfScopePackage is the only manifest package the root module cannot import.
// broker/ is a separate Go module, so its types are unreachable from here
// without a go.mod change; the constant exists so the reason is stated in one
// place instead of in prose at each use.
const outOfScopePackage = "broker."

// Unmapped is the coverage report of the symbol table: which manifest symbols
// the harness could not compare, and why.
//
// Only the first field is a defect in the table. OutOfScope is the broker/ module
// boundary, which no edit to this file can cross, and NoSymbol is the docgen em
// dash, which names no symbol to compare in the first place. A non-empty
// Symbols is a hole in the table, and TestSymbolTableCoversTheManifest fails on
// one by name so a hole cannot read as coverage.
type Unmapped struct {
	// Symbols are manifest symbols the table has no entry for.
	Symbols []string
	// OutOfScope are the broker/ module's, unreachable from the root module.
	OutOfScope []string
	// NoSymbol is the em-dash placeholder, the one spelling that names no SDK
	// symbol. It is a single value rather than a list of rows because the count
	// of em-dash rows is a coverage number reported from CompareAll.
	NoSymbol string
}

// UnmappedReport classifies every distinct symbol the manifest names.
//
// The split goes through Subject so a bare alternative is normalised against its
// row's package, exactly as the comparison does. Splitting the raw string here
// instead would report "GetPositions" as an unmapped symbol when the table has
// brokerfd.GetPositions, which is a false alarm in the one report whose job is
// to have no false alarms.
func UnmappedReport(m *Manifest) Unmapped {
	var u Unmapped
	seen := map[string]bool{}
	for _, f := range m.Fixtures {
		if f.SDKSymbol == "" || f.SDKSymbol == NoSDKSymbol {
			continue
		}
		subject, alts, named := Subject(f.SDKSymbol)
		if !named {
			continue
		}
		for _, sym := range append([]string{subject}, alts...) {
			if seen[sym] {
				continue
			}
			seen[sym] = true
			if strings.HasPrefix(sym, outOfScopePackage) {
				u.OutOfScope = append(u.OutOfScope, sym)
				continue
			}
			if _, found := SDKTypes[sym]; !found {
				u.Symbols = append(u.Symbols, sym)
			}
		}
	}
	sort.Strings(u.Symbols)
	sort.Strings(u.OutOfScope)
	u.NoSymbol = NoSDKSymbol
	return u
}

// TableSymbols returns every symbol in the table, sorted.
func TableSymbols() []string {
	out := make([]string, 0, len(SDKTypes))
	for sym := range SDKTypes {
		out = append(out, sym)
	}
	sort.Strings(out)
	return out
}

// NotCompared records a manifest row the harness did not compare, and why.
//
// Every reason is a fact about reachability, recorded rather than assumed. A row
// the harness silently passed over would be indistinguishable from a row it
// checked and found correct, which is the one confusion this whole instrument
// exists to end.
type NotCompared struct {
	// Fixture is the manifest row.
	Fixture string `json:"fixture"`
	// Symbol is the row's subject, or "" when it names none.
	Symbol string `json:"symbol,omitempty"`
	// Reason states why no comparison ran.
	Reason string `json:"reason"`
}

// Reasons CompareAll reports. They are constants so a caller can tell a coverage
// boundary from a table defect without matching on prose.
const (
	// ReasonOutOfScope is the broker/ module boundary.
	ReasonOutOfScope = "broker/ is a separate Go module, so its types are not " +
		"importable from the root module without a go.mod change"
	// ReasonNoSymbol is the docgen em dash.
	ReasonNoSymbol = "the docgen manifest maps this endpoint to no SDK symbol"
	// ReasonUntableable is a hole in SDKTypes and must never appear.
	ReasonUntableable = "the manifest names a symbol SDKTypes has no entry for"
)

// CompareAll runs the comparison for every manifest row.
//
// Every row lands in exactly one of the two returned slices, so the union is the
// whole manifest: an outcome for a row compared, a NotCompared for a row not
// compared, each with a reason. A row whose symbol is neither in the table nor
// in an out-of-scope package yields a NotCompared carrying ReasonUntableable, and
// the gate turns that into a failure -- a table hole must never read as coverage.
//
// The error return is for the table itself, not for the SDK. A row whose decode
// target cannot be resolved -- an unexported envelope the source parser refuses,
// or an entry that names both a type and an envelope -- is a defect in symbols.go
// and is reported as one, because returning a nil type instead would make the
// comparison announce that the method decodes no body, which is a different and
// untrue statement about a method that plainly decodes one.
func CompareAll(m *Manifest) (outcomes []Outcome, notCompared []NotCompared, err error) {
	for _, f := range m.Fixtures {
		if f.SDKSymbol == "" || f.SDKSymbol == NoSDKSymbol {
			notCompared = append(notCompared, NotCompared{
				Fixture: f.ID, Reason: ReasonNoSymbol})
			continue
		}
		subject, _, named := Subject(f.SDKSymbol)
		if !named {
			notCompared = append(notCompared, NotCompared{
				Fixture: f.ID, Reason: ReasonNoSymbol})
			continue
		}
		if strings.HasPrefix(subject, outOfScopePackage) {
			notCompared = append(notCompared, NotCompared{
				Fixture: f.ID, Symbol: subject, Reason: ReasonOutOfScope})
			continue
		}
		entry, found := SDKTypes[subject]
		if !found {
			notCompared = append(notCompared, NotCompared{
				Fixture: f.ID, Symbol: subject, Reason: ReasonUntableable})
			continue
		}
		target, resolveErr := entry.DecodeTarget()
		if resolveErr != nil {
			return nil, nil, fmt.Errorf("conformance: %s (%s): %w", subject, f.ID, resolveErr)
		}
		outcomes = append(outcomes, Compare(f, subject, target))
	}
	return outcomes, notCompared, nil
}
