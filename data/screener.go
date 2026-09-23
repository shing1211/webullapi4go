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
	"context"
	"net/url"

	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

// Screener endpoints.
const (
	// pathGainersLosers is the top gainers/losers endpoint.
	//
	// Reference: https://developer.webull.hk/apis/docs/reference/get-gainers-losers.md
	pathGainersLosers = "/market-data/screeners/gainers-losers/list"
	// pathTopActives is the most-active endpoint. The documented path is
	// /market-data/screeners/top-actives/list with operationId "getTopActive";
	// the reference URL get-most-active.md does not exist.
	//
	// Reference: https://developer.webull.hk/apis/docs/reference/get-top-active.md
	pathTopActives = "/market-data/screeners/top-actives/list"
	// pathMarketSectors is the market-sectors endpoint.
	//
	// Reference: https://developer.webull.hk/apis/docs/reference/get-market-sectors.md
	pathMarketSectors = "/market-data/screeners/market-sectors/list"
	// pathMarketSectorDetail is the market-sector detail endpoint.
	//
	// Reference: https://developer.webull.hk/apis/docs/reference/get-market-sector-detail.md
	pathMarketSectorDetail = "/market-data/screeners/market-sectors/get"
	// pathHighDividend is the high-dividend-ranks endpoint.
	//
	// Reference: https://developer.webull.hk/apis/docs/reference/get-high-dividend-ranks.md
	pathHighDividend = "/market-data/screeners/high-dividend-ranks/list"
	// pathWeek52HighLow is the week52-high-low endpoint.
	//
	// Reference: https://developer.webull.hk/apis/docs/reference/get-week52-high-low.md
	pathWeek52HighLow = "/market-data/screeners/week52-high-low/list"
)

// GainersLosersRankType is the time window over which gainers and losers are
// ranked.
type GainersLosersRankType string

// Gainers/losers ranking windows.
const (
	// GainersLosersRankPreMarket ranks by the pre-market move.
	GainersLosersRankPreMarket GainersLosersRankType = "PRE_MARKET"
	// GainersLosersRankAfterMarket ranks by the after-market move.
	GainersLosersRankAfterMarket GainersLosersRankType = "AFTER_MARKET"
	// GainersLosersRankMin3 ranks by the three-minute move.
	GainersLosersRankMin3 GainersLosersRankType = "MIN_3"
	// GainersLosersRankMin5 ranks by the five-minute move.
	GainersLosersRankMin5 GainersLosersRankType = "MIN_5"
	// GainersLosersRankDay1 ranks by the one-day move.
	GainersLosersRankDay1 GainersLosersRankType = "DAY_1"
	// GainersLosersRankDay5 ranks by the five-day move.
	GainersLosersRankDay5 GainersLosersRankType = "DAY_5"
	// GainersLosersRankMonth1 ranks by the one-month move.
	GainersLosersRankMonth1 GainersLosersRankType = "MONTH_1"
	// GainersLosersRankMonth3 ranks by the three-month move.
	GainersLosersRankMonth3 GainersLosersRankType = "MONTH_3"
	// GainersLosersRankWeek52 ranks by the 52-week move.
	GainersLosersRankWeek52 GainersLosersRankType = "WEEK_52"
)

// MostActiveRankType is the activity metric used to rank the most active
// stocks.
type MostActiveRankType string

// Most-active ranking metrics.
const (
	// MostActiveRankVolume ranks by cumulative volume.
	MostActiveRankVolume MostActiveRankType = "VOLUME"
	// MostActiveRankRelativeVolume10D ranks by relative volume against the
	// ten-day average.
	MostActiveRankRelativeVolume10D MostActiveRankType = "RELATIVE_VOLUME_10D"
	// MostActiveRankTurnover ranks by cumulative turnover.
	MostActiveRankTurnover MostActiveRankType = "TURNOVER"
	// MostActiveRankTurnoverRate ranks by turnover rate.
	MostActiveRankTurnoverRate MostActiveRankType = "TURNOVER_RATE"
	// MostActiveRankAmplitude ranks by price amplitude.
	MostActiveRankAmplitude MostActiveRankType = "AMPLITUDE"
)

// ScreenerSortBy is the secondary sort field shared by the screener endpoints.
type ScreenerSortBy string

// Screener sort fields.
const (
	// ScreenerSortChangeRatio sorts by price change percentage.
	ScreenerSortChangeRatio ScreenerSortBy = "CHANGE_RATIO"
	// ScreenerSortRelativeVolume10D sorts by relative ten-day volume.
	ScreenerSortRelativeVolume10D ScreenerSortBy = "RELATIVE_VOLUME_10D"
	// ScreenerSortMarketValue sorts by market capitalization.
	ScreenerSortMarketValue ScreenerSortBy = "MARKET_VALUE"
	// ScreenerSortClose sorts by the latest close.
	ScreenerSortClose ScreenerSortBy = "CLOSE"
	// ScreenerSortPrice sorts by the latest price.
	ScreenerSortPrice ScreenerSortBy = "PRICE"
	// ScreenerSortPETTM sorts by trailing-twelve-month price/earnings.
	ScreenerSortPETTM ScreenerSortBy = "PE_TTM"
	// ScreenerSortHigh sorts by the intraday high.
	ScreenerSortHigh ScreenerSortBy = "HIGH"
	// ScreenerSortLow sorts by the intraday low.
	ScreenerSortLow ScreenerSortBy = "LOW"
	// ScreenerSortAmplitude sorts by price amplitude.
	ScreenerSortAmplitude ScreenerSortBy = "AMPLITUDE"
	// ScreenerSortTurnover sorts by turnover.
	ScreenerSortTurnover ScreenerSortBy = "TURNOVER"
	// ScreenerSortVolume sorts by volume.
	ScreenerSortVolume ScreenerSortBy = "VOLUME"
)

// SortDirection is the direction of a screener sort.
type SortDirection string

// Sort directions.
const (
	// SortDirectionAsc sorts ascending.
	SortDirectionAsc SortDirection = "ASC"
	// SortDirectionDesc sorts descending.
	SortDirectionDesc SortDirection = "DESC"
)

// GainersLosersQuery parameterizes [Client.GetTopGainersLosers]. RankType,
// Category and SortBy are required. Direction defaults to descending on the
// server; pass [SortDirectionAsc] to retrieve losers rather than gainers.
type GainersLosersQuery struct {
	// RankType is the time window used to rank the price change. Required.
	RankType GainersLosersRankType
	// Category is the security market. Required, and only [StockCategoryUS]
	// is supported.
	Category StockCategory
	// SortBy is the secondary sort field. Required.
	SortBy ScreenerSortBy
	// Direction is the sort direction. Empty uses the server default.
	Direction SortDirection
}

// MostActiveQuery parameterizes [Client.GetMostActive]. Category is required;
// the remaining fields default on the server.
type MostActiveQuery struct {
	// Category is the security market. Required, and only [StockCategoryUS]
	// is supported.
	Category StockCategory
	// RankType is the activity metric used to rank the result. Empty uses the
	// server default of [MostActiveRankVolume].
	RankType MostActiveRankType
	// SortBy is the secondary sort field. Empty uses the server default of
	// [ScreenerSortVolume].
	SortBy ScreenerSortBy
	// Direction is the sort direction. Empty uses the server default.
	Direction SortDirection
}

// MarketSector represents a market sector overview.
type MarketSector struct {
	SectorName  string          `json:"sector_name"`
	ChangeRatio string          `json:"change_ratio"`
	Volume      string          `json:"volume"`
	MarketValue string          `json:"market_value"`
	Stocks      []ScreenerStock `json:"stocks"`
}

// MarketSectorDetailQuery parameterizes [Client.GetMarketSectorDetail].
type MarketSectorDetailQuery struct {
	// SectorName is the sector to query. Required.
	SectorName string
	// Category is the security market. Required.
	Category StockCategory
	// SortBy is the secondary sort field. Empty uses the server default.
	SortBy ScreenerSortBy
	// Direction is the sort direction. Empty uses the server default.
	Direction SortDirection
}

// HighDividendQuery parameterizes [Client.GetHighDividendRank].
type HighDividendQuery struct {
	// Category is the security market. Required.
	Category StockCategory
	// SortBy is the secondary sort field. Empty uses the server default.
	SortBy ScreenerSortBy
	// Direction is the sort direction. Empty uses the server default.
	Direction SortDirection
}

// Week52HighLowQuery parameterizes [Client.GetWeek52HighLow].
type Week52HighLowQuery struct {
	// Category is the security market. Required.
	Category StockCategory
	// SortBy is the secondary sort field. Empty uses the server default.
	SortBy ScreenerSortBy
	// Direction is the sort direction. Empty uses the server default.
	Direction SortDirection
}

// ScreenerStock is a single stock row returned by the screener endpoints. Not
// every field is populated by every endpoint; for example RelativeVolume10D is
// only reported by the most-active endpoint.
type ScreenerStock struct {
	// InstrumentID is the unique identifier of the tradable instrument.
	InstrumentID string `json:"instrument_id"`
	// Symbol is the trading symbol, for example "AAPL".
	Symbol string `json:"symbol"`
	// Name is the full name of the instrument.
	Name string `json:"name"`
	// ExchangeCode is the standardized exchange code, for example "NSQ".
	ExchangeCode string `json:"exchange_code"`
	// CurrencyCode is the denomination currency (ISO 4217), for example "USD".
	CurrencyCode string `json:"currency_code"`
	// PreClose is the previous trading day's closing price.
	PreClose money.Money `json:"pre_close"`
	// Open is the opening price for the current trading day.
	Open money.Money `json:"open"`
	// High is the intraday high for the current trading day.
	High money.Money `json:"high"`
	// Low is the intraday low for the current trading day.
	Low money.Money `json:"low"`
	// Close is the latest traded price for the current trading day.
	Close money.Money `json:"close"`
	// Price is the most recent quoted price within the selected time interval.
	Price money.Money `json:"price"`
	// Change is the absolute price change within the selected time interval.
	Change money.Money `json:"change"`
	// ChangeRatio is the price change percentage within the selected time
	// interval, as a decimal ratio.
	ChangeRatio string `json:"change_ratio"`
	// Volume is the cumulative traded volume for the current day.
	Volume string `json:"volume"`
	// Turnover is the cumulative turnover amount in the denomination currency.
	Turnover money.Money `json:"turnover"`
	// TurnoverRate is the turnover rate as a decimal ratio.
	TurnoverRate string `json:"turnover_rate"`
	// MarketValue is the total market capitalization in the denomination
	// currency.
	MarketValue money.Money `json:"market_value"`
	// Amplitude is (high-low)/pre_close as a decimal ratio.
	Amplitude string `json:"amplitude"`
	// RelativeVolume10D is the current-day volume divided by the ten-day
	// average volume.
	RelativeVolume10D string `json:"relative_volume_10d"`
}

// GetTopGainersLosers retrieves the top gaining or losing US stocks for a
// ranking window. Pass [SortDirectionAsc] to rank losers.
//
// Reference: https://developer.webull.hk/apis/docs/reference/get-gainers-losers.md
func (c *Client) GetTopGainersLosers(ctx context.Context, q GainersLosersQuery) ([]ScreenerStock, error) {
	query := url.Values{}
	if q.RankType != "" {
		query.Set("rank_type", string(q.RankType))
	}
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	if q.SortBy != "" {
		query.Set("sort_by", string(q.SortBy))
	}
	if q.Direction != "" {
		query.Set("direction", string(q.Direction))
	}

	var out []ScreenerStock
	if err := c.get(ctx, pathGainersLosers, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMostActive retrieves the most actively traded US stocks, ranked by the
// metric selected with [MostActiveQuery.RankType].
//
// The Webull reference URL get-most-active.md does not exist; the endpoint is
// documented as "List Top Actives" (operationId getTopActive).
//
// Reference: https://developer.webull.hk/apis/docs/reference/get-top-active.md
func (c *Client) GetMostActive(ctx context.Context, q MostActiveQuery) ([]ScreenerStock, error) {
	query := url.Values{}
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	if q.RankType != "" {
		query.Set("rank_type", string(q.RankType))
	}
	if q.SortBy != "" {
		query.Set("sort_by", string(q.SortBy))
	}
	if q.Direction != "" {
		query.Set("direction", string(q.Direction))
	}

	var out []ScreenerStock
	if err := c.get(ctx, pathTopActives, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMarketSectors retrieves the list of market sectors with their aggregate
// statistics and constituent stocks.
//
// Reference: https://developer.webull.hk/apis/docs/reference/get-market-sectors.md
func (c *Client) GetMarketSectors(ctx context.Context) ([]MarketSector, error) {
	var out []MarketSector
	if err := c.get(ctx, pathMarketSectors, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetMarketSectorDetail retrieves the constituent stocks of a single market
// sector, with optional sorting.
//
// Reference: https://developer.webull.hk/apis/docs/reference/get-market-sector-detail.md
func (c *Client) GetMarketSectorDetail(ctx context.Context, q MarketSectorDetailQuery) ([]ScreenerStock, error) {
	query := url.Values{}
	if q.SectorName != "" {
		query.Set("sector", q.SectorName)
	}
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	if q.SortBy != "" {
		query.Set("sort_by", string(q.SortBy))
	}
	if q.Direction != "" {
		query.Set("direction", string(q.Direction))
	}

	var out []ScreenerStock
	if err := c.get(ctx, pathMarketSectorDetail, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetHighDividendRank retrieves US stocks ranked by dividend yield.
//
// Reference: https://developer.webull.hk/apis/docs/reference/get-high-dividend-ranks.md
func (c *Client) GetHighDividendRank(ctx context.Context, q HighDividendQuery) ([]ScreenerStock, error) {
	query := url.Values{}
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	if q.SortBy != "" {
		query.Set("sort_by", string(q.SortBy))
	}
	if q.Direction != "" {
		query.Set("direction", string(q.Direction))
	}

	var out []ScreenerStock
	if err := c.get(ctx, pathHighDividend, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetWeek52HighLow retrieves US stocks based on their 52-week high/low
// performance.
//
// Reference: https://developer.webull.hk/apis/docs/reference/get-week52-high-low.md
func (c *Client) GetWeek52HighLow(ctx context.Context, q Week52HighLowQuery) ([]ScreenerStock, error) {
	query := url.Values{}
	if q.Category != "" {
		query.Set("category", string(q.Category))
	}
	if q.SortBy != "" {
		query.Set("sort_by", string(q.SortBy))
	}
	if q.Direction != "" {
		query.Set("direction", string(q.Direction))
	}

	var out []ScreenerStock
	if err := c.get(ctx, pathWeek52HighLow, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}
