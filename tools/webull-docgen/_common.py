# Copyright 2026 shing1211
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Shared manifest and rendering helpers for the Webull doc generator.

Fetches Webull's machine-readable ``.md`` documentation (which embeds the
OpenAPI definition JSON), caches it locally, and renders the SDK-mapped
reference pages, the verbatim master guides/reference, and the SDK-API
reconciliation. No SDK-specific knowledge lives here beyond the area manifest.
"""
import datetime
import json
import os
import re
import sys
import time
import urllib.request
from collections import OrderedDict

HK = "https://developer.webull.hk/apis/docs/reference/"
US = "https://developer.webull.com/apis/docs/reference/"
HK_DOCS = "https://developer.webull.hk/apis/"
US_DOCS = "https://developer.webull.com/apis/"

_HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.abspath(os.path.join(_HERE, "..", ".."))
DOCS = os.path.join(ROOT, "docs")
REFERENCE_OUT = os.path.join(DOCS, "webull-api")
MASTER_OUT = os.path.join(DOCS, "webull-api")
RECON_OUT = os.path.join(DOCS, "reconciliation.md")
CACHE = os.path.join(_HERE, ".cache")

UA = "Mozilla/5.0 (compatible; webullapi4go-docgen/1.0)"

# area key -> (page title, blurb, [ (label, official .md url, sdk func, note) ])
AREAS = OrderedDict()

AREAS["authentication"] = (
    "Authentication",
    "Webull uses a dual layer: an HMAC-SHA1 request signature plus an access "
    "token. Server-to-server endpoints sign every request; Display Solution "
    "(Client-to-Server) uses an OAuth-style client token.",
    [
        ("Create Token", HK + "create-token.md", "client.CreateToken", "Returns a token in `PENDING`; completes 2FA in the Webull App. Sandbox tokens are `NORMAL` immediately."),
        ("Check Token", HK + "check-token.md", "client.CheckToken", "Polled by `client.EnsureToken` until the token is `NORMAL`."),
        ("Create Client Token (Display)", HK + "market-display-solution-data-api/create-client-token.md", "display.Service.EnsureToken", "Internal to `display.Service`; `client_user_id` is fixed to `openapi_client`."),
        ("Refresh Client Token (Display)", HK + "market-display-solution-data-api/refresh-client-token.md", "display.Service.RefreshClientToken", "The SDK re-creates the client token instead of refreshing it."),
    ],
)

AREAS["market-data-stock"] = (
    "Market Data — Stock",
    "HTTP on-demand stock/ETF market data (Non-Display Solution).",
    [
        ("Stock Snapshot", HK + "snapshot.md", "data.GetSnapshot", ""),
        ("Stock Tick", HK + "tick.md", "data.GetTick", ""),
        ("Stock Quotes (Depth)", HK + "quotes.md", "data.GetQuotes", ""),
        ("Stock Historical Bars (Batch)", HK + "historical-bars.md", "data.GetBatchBars", "`data.GetBars` POSTs a one-symbol batch against this same endpoint."),
        ("Stock Footprint", HK + "footprint.md", "data.GetFootprint", "Requires a paid entitlement."),
        ("NOII Bars", HK + "get-noii-bars.md", "data.GetNOIIBars", ""),
        ("NOII Snapshot", HK + "get-noii-snapshot.md", "data.GetNOIISnapshot", ""),
    ],
)

AREAS["market-data-option"] = (
    "Market Data — Option",
    "Option tick, snapshot and historical bars (Non-Display Solution), plus the "
    "option contract list used to build an option chain.",
    [
        ("Option Tick", HK + "option-tick.md", "data.GetOptionTick", ""),
        ("Option Snapshot", HK + "option-snapshot.md", "data.GetOptionSnapshot", "At most 20 symbols per query."),
        ("Option Historical Bars", HK + "option-historical-bars.md", "data.GetOptionBars", "At most 20 symbols per query."),
        ("Option Contract List (Chain)", HK + "option-contract-list.md", "data.GetOptionContracts", "Trading API surface."),
    ],
)

AREAS["market-data-crypto"] = (
    "Market Data — Crypto",
    "Crypto market data and instruments (US only).",
    [
        ("Crypto Snapshot", US + "crypto-snapshot.md", "data.GetCryptoSnapshot", ""),
        ("Crypto Bars", US + "crypto-bars.md", "data.GetCryptoBars", ""),
        ("Crypto Instruments", US + "crypto-instrument-list.md", "data.GetCryptoInstruments", "Trading API surface."),
    ],
)

AREAS["market-data-futures"] = (
    "Market Data — Futures",
    "Futures market data and instruments.",
    [
        ("Futures Tick", HK + "futures-tick.md", "data.GetFuturesTick", ""),
        ("Futures Snapshot", HK + "futures-snapshot.md", "data.GetFuturesSnapshot", ""),
        ("Futures Footprint", HK + "futures-footprint.md", "data.GetFuturesFootprint", ""),
        ("Futures Quotes (Depth)", HK + "futures-depth-of-book.md", "data.GetFuturesDepth", ""),
        ("Futures Historical Bars", HK + "futures-historical-bars.md", "data.GetFuturesBars", ""),
        ("Futures Instrument List", HK + "futures-instrument-list.md", "data.GetFuturesInstruments", ""),
        ("Futures Product Codes", HK + "futures-products.md", "data.GetFuturesProductCodes", ""),
        ("Futures Product Classes", HK + "futures-products-class.md", "data.GetFuturesProductClasses", ""),
    ],
)

AREAS["market-data-news"] = (
    "Market Data — News",
    "News summaries. The Non-Display endpoint streams Server-Sent Events.",
    [
        ("News Summary", HK + "news-summary.md", "data.GetNewsSummary", "SSE stream via `client.DoStream`; HK sandbox upstream returns 504."),
    ],
)

AREAS["market-data-screener"] = (
    "Market Data — Screener",
    "Ranked lists and sector data. Some rankers require the Display Solution "
    "entitlement (marked below).",
    [
        ("Top Gainers/Losers", HK + "get-gainers-losers.md", "data.GetTopGainersLosers", ""),
        ("Top Actives", HK + "get-top-active.md", "data.GetMostActive", ""),
        ("Market Sectors", HK + "get-market-sectors.md", "data.GetMarketSectors", "Display Solution; HK sandbox host blocked (403)."),
        ("Market Sector Detail", HK + "get-market-sectors-detail.md", "data.GetMarketSectorDetail", "Display Solution."),
        ("High Dividend Rank", HK + "get-high-dividend.md", "data.GetHighDividendRank", "Display Solution."),
        ("52-Week High/Low", HK + "get-week-52-high-low.md", "data.GetWeek52HighLow", "Display Solution."),
    ],
)

AREAS["market-data-watchlist"] = (
    "Market Data — Watchlist",
    "Watchlist CRUD and instrument membership.",
    [
        ("Get Watchlists", HK + "get-watchlist.md", "data.GetWatchlists", ""),
        ("Create Watchlist", HK + "create-watchlist.md", "data.CreateWatchlist", ""),
        ("Update Watchlist", HK + "update-watchlist.md", "data.UpdateWatchlist", ""),
        ("Delete Watchlist", HK + "delete-watchlist.md", "data.DeleteWatchlist", ""),
        ("Get Watchlist Instruments", HK + "get-watchlist-instruments.md", "data.GetWatchlistInstruments", ""),
        ("Add Watchlist Instruments", HK + "add-watchlist-instruments.md", "data.AddWatchlistInstruments", ""),
        ("Remove Watchlist Instruments", HK + "remove-watchlist-instruments.md", "data.RemoveWatchlistInstruments", ""),
        ("Update Watchlist Instruments", HK + "update-watchlist-instruments.md", "data.UpdateWatchlistInstruments", ""),
    ],
)

AREAS["fundamentals"] = (
    "Fundamentals and Fund Data",
    "Company fundamentals, analyst data, financial statements and fund data "
    "under the Non-Display Solution.",
    [
        ("Company Profile", HK + "get-company-profile.md", "data.GetCompanyProfile", "US only in practice."),
        ("Analyst Target Price", HK + "get-analyst-target-price.md", "data.GetAnalystTargetPrice", ""),
        ("Analyst Rating", HK + "get-analyst-rating.md", "data.GetAnalystRating", ""),
        ("Capital Flow", HK + "capital-flow.md", "data.GetCapitalFlow", ""),
        ("Industry Comparison", HK + "industry-comparison.md", "data.GetIndustryComparison", ""),
        ("Earnings Calendar", HK + "earnings-calendar.md", "data.GetEarningsCalendar", ""),
        ("Dividend Calendar", HK + "dividend-calendar.md", "data.GetDividendCalendar", ""),
        ("Filings", HK + "filings.md", "data.GetFilings", "US/SEC only."),
        ("Income Statement", HK + "financial-income.md", "data.GetIncomeStatement", "Returns `map[string]any` per period."),
        ("Balance Sheet", HK + "financial-balancesheet.md", "data.GetBalanceSheet", "Returns `map[string]any` per period."),
        ("Cash Flow", HK + "financial-cashflow.md", "data.GetCashFlow", "Returns `map[string]any` per period."),
        ("Financial Indicators", HK + "financial-indicators.md", "data.GetFinancialIndicators", ""),
        ("Financial Alert", HK + "financial-alert.md", "data.GetFinancialAlert", ""),
        ("Forecast EPS", HK + "forecast-eps.md", "data.GetForecastEPS", "Most recent 5 quarters."),
        ("Fund Brief", HK + "fund-brief.md", "data.GetFundInfo", ""),
        ("Fund Performance", HK + "fund-performance.md", "data.GetFundPerformance", ""),
        ("Fund Net Value", HK + "fund-net-value.md", "data.GetFundNav", ""),
        ("Fund Holdings", HK + "fund-holdings.md", "data.GetFundHoldings", ""),
        ("Fund Dividends", HK + "fund-dividends.md", "data.GetFundDividends", ""),
        ("Fund Rating", HK + "fund-rating.md", "data.GetFundRating", ""),
        ("Fund Splits", HK + "fund-splits.md", "data.GetFundSplits", ""),
        ("Fund Files", HK + "fund-files.md", "data.GetFundFiles", ""),
        ("Fund Allocation", HK + "fund-allocation.md", "data.GetFundAllocation", ""),
    ],
)

AREAS["event-contracts"] = (
    "Event Contracts",
    "Event-contract instruments and market data (US documentation).",
    [
        ("Event Contract Categories", US + "event-categories-list.md", "data.GetEventContractCategories", "Documented on the US site."),
        ("Event Contract Series", US + "event-series-list.md", "data.GetEventContractSeries", ""),
        ("Event Contract Events", US + "event-events-list.md", "data.GetEventContractEvents", ""),
        ("Event Contract Instruments", US + "event-market-list.md", "data.GetEventContractMarkets", ""),
        ("Event Snapshot", US + "event-snapshot.md", "data.GetEventSnapshot", ""),
        ("Event Depth", US + "event-depth.md", "data.GetEventDepth", ""),
        ("Event Bars", US + "event-bars.md", "data.GetEventBars", ""),
        ("Event Tick", US + "event-tick.md", "data.GetEventTick", ""),
        ("Event Contract Tags", US + "broker-market-data-api/all-tags-using-get.md", "data.GetEventContractTags", ""),
        ("Event Contract Events List", US + "broker-market-data-api/event-list-using-get.md", "data.GetEventContractEventsList", ""),
        ("Event Contract Milestones", US + "broker-market-data-api/milestones-using-get.md", "data.GetEventContractMilestones", ""),
        ("Event Contract Series List", US + "broker-market-data-api/series-list-using-get.md", "data.GetEventContractSeriesList", ""),
        ("Event Contract Sports Filters", US + "broker-market-data-api/sports-filter-using-get.md", "data.GetEventContractSportsFilters", ""),
        ("Event Game Stats", US + "broker-market-data-api/event-game-stats-using-get.md", "data.GetEventGameStats", ""),
        ("Event Live Data", US + "broker-market-data-api/event-live-data-using-get.md", "data.GetEventLiveData", ""),
        ("Event Market Bars", US + "broker-market-data-api/event-market-bars-using-get.md", "data.GetEventMarketBars", ""),
        ("Event Market Bars By Event", US + "broker-market-data-api/event-market-bars-by-event-using-get.md", "data.GetEventMarketBarsByEvent", ""),
        ("Event Market Depth", US + "broker-market-data-api/event-market-depth-using-get.md", "data.GetEventMarketDepth", ""),
        ("Event Market Snapshot", US + "broker-market-data-api/event-market-snapshot-using-get.md", "data.GetEventMarketSnapshot", ""),
    ],
)

AREAS["trading"] = (
    "Trading API",
    "Accounts, assets, the order lifecycle, order queries and instruments. "
    "Requests require an access token and default to API version `v3`.",
    [
        ("Get Instruments", HK + "instrument-list.md", "data.GetStockInstruments", ""),
        ("Account List", HK + "account-list.md", "trade.ListAccounts", ""),
        ("Account Balance", HK + "query-account-balance.md", "trade.GetBalance", ""),
        ("Account Positions", HK + "query-account-position.md", "trade.GetPositions", ""),
        ("Order Preview", HK + "common-order-preview.md", "trade.PreviewOrder", "Validates and enforces order guardrails."),
        ("Order Place", HK + "common-order-place.md", "trade.PlaceOrder", "Creates live orders."),
        ("Batch Place Orders", US + "order-batch-place.md", "trade.BatchPlaceOrder", "Equity only, up to 50 orders."),
        ("Order Replace", HK + "common-order-replace.md", "trade.ReplaceOrder", ""),
        ("Order Cancel", HK + "common-order-cancel.md", "trade.CancelOrder", ""),
        ("Open Orders", HK + "order-open.md", "trade.GetOpenOrders", ""),
        ("Order History", HK + "order-history.md", "trade.GetOrderHistory", ""),
        ("Order Detail", HK + "order-detail.md", "trade.GetOrderDetail", ""),
        ("Cash Activities", US + "trade-cash-activity-by-type.md", "trade.GetCashActivities", "Documented on the US site."),
    ],
)

AREAS["broker-hk"] = (
    "Broker API — HK",
    "Institutional Broker API for Hong Kong. Uses the `DoBroker` transport. The "
    "HK sandbox returns `401 ROUTE_NOT_PERMITTED` (missing app scope).",
    [
        ("Create Virtual Account", HK + "broker-api/broker-account-create.md", "broker.CreateVirtualAccount", ""),
        ("Update Virtual Account", HK + "broker-api/broker-account-update.md", "broker.UpdateVirtualAccount", ""),
        ("Get Virtual Account Detail", HK + "broker-api/broker-account-detail.md", "broker.GetVirtualAccount", ""),
        ("List Virtual Accounts", HK + "broker-api/broker-account-list.md", "broker.ListVirtualAccounts", ""),
        ("Get Stock Instrument", HK + "broker-api/broker-instrument-list.md", "broker.GetStockInstruments", ""),
        ("Get Stock Locate Detail", HK + "broker-api/broker-stock-locate-detail.md", "broker.GetStockLocate", ""),
        ("Get Corporate Actions Detail", HK + "broker-api/broker-corporate-actions-detail.md", "broker.GetCorporateActionsDetail", ""),
        ("Account Activities By Type", HK + "broker-api/broker-activity-by-type.md", "broker.GetCashActivities", ""),
        ("Account Balance", HK + "broker-api/broker-assets-balance.md", "broker.GetBalance", ""),
        ("Account Positions", HK + "broker-api/broker-assets-positions.md", "broker.GetPositions", ""),
        ("Order Preview", HK + "broker-api/broker-order-preview.md", "broker.PreviewOrder", ""),
        ("Order Place", HK + "broker-api/broker-order-place.md", "broker.PlaceOrder", ""),
        ("Order Replace", HK + "broker-api/broker-order-replace.md", "broker.ReplaceOrder", ""),
        ("Order Cancel", HK + "broker-api/broker-order-cancel.md", "broker.CancelOrder", ""),
        ("Order Detail", HK + "broker-api/broker-order-detail.md", "broker.GetOrderDetail", ""),
        ("Order History", HK + "broker-api/broker-order-history.md", "broker.GetOrderHistory", ""),
        ("Open Orders", HK + "broker-api/broker-order-open.md", "broker.GetOpenOrders", ""),
        ("Get FX Rate", HK + "broker-api/broker-funding-query-rate.md", "broker.GetFXRate", ""),
        ("Create FX Request", HK + "broker-api/broker-funding-create-fx.md", "broker.CreateFXExchange", ""),
        ("Get FX Detail", HK + "broker-api/broker-funding-query-fx.md", "broker.GetFXExchangeDetail", ""),
        ("Create Instant Exchange", HK + "broker-api/broker-funding-create-instant-fx.md", "broker.CreateInstantExchange", ""),
        ("Get Instant Exchange Detail", HK + "broker-api/broker-funding-query-instant-fx.md", "broker.GetInstantExchangeDetail", ""),
        ("Create Instant Funding", HK + "broker-api/broker-funding-instant-create.md", "broker.CreateInstantFunding", ""),
        ("Get Instant Funding Detail", HK + "broker-api/broker-funding-instant-query.md", "broker.GetInstantFundingDetail", ""),
        ("Create Cash Journal", HK + "broker-api/broker-journal-cash-create.md", "broker.CreateCashJournal", ""),
        ("Get Cash Journal Detail", HK + "broker-api/broker-journal-cash-query.md", "broker.GetCashJournalDetail", ""),
        ("Create Position Journal", HK + "broker-api/broker-journal-position-create.md", "broker.CreatePositionJournal", ""),
        ("Get Position Journal Detail", HK + "broker-api/broker-journal-position-query.md", "broker.GetPositionJournalDetail", ""),
        ("Trade Calendar", HK + "broker-api/broker-trade-calendar.md", "broker.GetTradeCalendar", ""),
        ("Account Events", HK + "custom/broker-account-events.md", "—", "gRPC subscription."),
        ("Instrument Events", HK + "custom/broker-instrument-events.md", "—", "gRPC subscription."),
        ("Corporate Actions Events", HK + "custom/broker-ca-events.md", "—", "gRPC subscription."),
        ("Trade Events", HK + "custom/broker-trade-events.md", "—", "gRPC subscription."),
        ("Funding Events", HK + "custom/broker-funding-events.md", "—", "gRPC subscription."),
        ("Journal Events", HK + "custom/broker-journal-events.md", "—", "gRPC subscription."),
        ("Master Data Events", HK + "custom/broker-master-data-events.md", "—", "gRPC subscription."),
    ],
)

AREAS["broker-fd-us"] = (
    "Broker API — FD (US)",
    "US Broker FD surface. Documented on the US site only.",
    [
        ("List Accounts", US + "broker-fd-api/list-accounts.md", "brokerfd.ListFDAccounts", ""),
        ("Account Detail", US + "broker-fd-api/get-account-detail.md", "brokerfd.GetFDAccountDetail", ""),
        ("Create Account", US + "broker-fd-api/create-account-apply.md", "brokerfd.CreateFDAccount", ""),
        ("Update Account", US + "broker-fd-api/update-account-apply.md", "brokerfd.UpdateFDAccount", ""),
        ("Close Account", US + "broker-fd-api/close-account.md", "brokerfd.CloseFDAccount", ""),
        ("Account Application Detail", US + "broker-fd-api/get-account-application-detail.md", "—", ""),
        ("List Forms", US + "broker-fd-api/get-form-list.md", "brokerfd.ListAccountForms", ""),
        ("List Form Versions", US + "broker-fd-api/get-form-version-list.md", "—", ""),
        ("Form Content", US + "broker-fd-api/get-form-content.md", "—", ""),
        ("Upload Document", US + "broker-fd-api/document-upload.md", "brokerfd.UploadDocument", ""),
        ("Download Document", US + "broker-fd-api/document-download.md", "brokerfd.DownloadDocument", ""),
        ("Assets Summary", US + "broker-fd-api/summary.md", "brokerfd.GetAccountsSummary / GetFDAssetsSummary", ""),
        ("Assets Detail", US + "broker-fd-api/account-balance.md", "brokerfd.GetFDAssetsDetail", ""),
        ("Positions", US + "broker-fd-api/account-position.md", "brokerfd.GetFDPositions / GetPositions", ""),
        ("Cash Activities By Type", US + "broker-fd-api/broker-cash-activity-by-type.md", "brokerfd.GetFDActivities", ""),
        ("Create Bank Relationship", US + "broker-fd-api/create-bank-relationship.md", "brokerfd.AddFDBankAccount", ""),
        ("Delete Bank Relationship", US + "broker-fd-api/delete-bank-relationship.md", "brokerfd.RemoveFDBankAccount", ""),
        ("List Bank Accounts", US + "broker-fd-api/list-linked-bank-accounts.md", "brokerfd.ListFDBankAccounts", ""),
        ("Create ACH Relationship", US + "broker-fd-api/create-ach-relationship.md", "brokerfd.AddFDAchAccount", ""),
        ("Delete ACH Relationship", US + "broker-fd-api/delete-ach-relationship.md", "brokerfd.RemoveFDAchAccount", ""),
        ("List ACH Relationships", US + "broker-fd-api/list-ach-relationships.md", "brokerfd.ListFDAchAccounts", ""),
        ("Create Transfer", US + "broker-fd-api/create-transfer.md", "brokerfd.InitiateFDTransfer", ""),
        ("Transfer List", US + "broker-fd-api/transfer-list.md", "brokerfd.ListFDTransfers", ""),
        ("Transfer Detail", US + "broker-fd-api/transfer-detail.md", "brokerfd.GetFDTransferDetail", ""),
        ("Cancel Transfer", US + "broker-fd-api/cancel-transfer.md", "—", ""),
        ("Create Instant Funding", US + "broker-fd-api/broker-funding-instant-create.md", "brokerfd.CreateFDInstantFunding", ""),
        ("Instant Funding Detail", US + "broker-fd-api/broker-funding-instant-query.md", "brokerfd.GetFDInstantFundingDetail", ""),
        ("Create Fee", US + "broker-fd-api/broker-funding-fee-create.md", "—", ""),
        ("Get Fee", US + "broker-fd-api/broker-funding-fee-query.md", "brokerfd.GetFDTransferFees", ""),
        ("Create Credit", US + "broker-fd-api/broker-funding-credit-create.md", "—", ""),
        ("Get Credit", US + "broker-fd-api/broker-funding-credit-query.md", "brokerfd.GetFDCreditInfo", ""),
        ("List Stock Instruments", US + "broker-fd-api/list-stock-instruments.md", "brokerfd.GetFDStockInstruments", ""),
        ("Event Contract Categories", US + "broker-fd-api/broker-event-categories-list.md", "—", ""),
        ("Event Contract Series", US + "broker-fd-api/broker-event-series-list.md", "—", ""),
        ("Event Contract Events", US + "broker-fd-api/broker-event-events-list.md", "—", ""),
        ("Event Contract Instruments", US + "broker-fd-api/broker-event-market-list.md", "brokerfd.GetFDECInstruments", ""),
        ("Corporate Actions Detail", US + "broker-fd-api/broker-corporate-actions-detail.md", "brokerfd.GetFDCorporateActions", ""),
        ("Order Preview", US + "broker-fd-api/common-order-preview.md", "brokerfd.PreviewFDOrder", ""),
        ("Order Place", US + "broker-fd-api/common-order-place.md", "brokerfd.PlaceFDOrder", ""),
        ("Order Replace", US + "broker-fd-api/common-order-replace.md", "brokerfd.ReplaceFDOrder", ""),
        ("Order Cancel", US + "broker-fd-api/common-order-cancel.md", "brokerfd.CancelFDOrder", ""),
        ("Open Orders", US + "broker-fd-api/order-open.md", "brokerfd.GetFDOpenOrders", ""),
        ("Order Detail", US + "broker-fd-api/order-detail.md", "brokerfd.GetFDOrderDetail", ""),
        ("Order History", US + "broker-fd-api/order-history.md", "brokerfd.GetFDOrderHistory", ""),
        ("Create Cash Journal", US + "broker-fd-api/broker-journal-cash-create.md", "—", ""),
        ("Cash Journal Detail", US + "broker-fd-api/broker-journal-cash-query.md", "brokerfd.GetFDCashJournalDetail", ""),
        ("List Enums", US + "broker-fd-api/list-enums.md", "brokerfd.GetFDEnums", ""),
        ("Trade Calendar", US + "broker-fd-api/list-trade-calendar.md", "brokerfd.GetFDTradeCalendar", ""),
        ("List Agreements", US + "broker-fd-api/broker-list-agreements-by-type.md", "brokerfd.ListAgreements", ""),
        ("Agreement Details", US + "broker-fd-api/broker-get-agreement-details.md", "brokerfd.GetAgreementDetail", ""),
    ],
)

AREAS["display-solution"] = (
    "Display Solution",
    "Hosted Display Solution: a separate entitlement and host with "
    "Client-to-Server (Bearer) authentication. The SDK routes these through "
    "`display.Service`.",
    [
        ("Stock Top Gainers/Losers", HK + "market-display-solution-data-api/top-gainers-using-get-new.md", "data.GetDisplayGainersLosers", "SDK rank_type values: `MIN_3`, `MIN_5`, `DAY_1`, `DAY_5`, `MONTH_1`, `MONTH_3`, `WEEK_52` (not `M3`, `D1`, etc.)."),
        ("Top Active", HK + "market-display-solution-data-api/top-active-using-get-new.md", "data.GetDisplayTopActive", ""),
        ("Snapshot", HK + "market-display-solution-data-api/snapshot-using-get.md", "data.GetDisplaySnapshot", "SDK sends GET with query params (`symbols`, `category`, `extend_hour_required`, `overnight_required`), not POST with `category_symbols` body."),
        ("Historical Bars (Batch)", HK + "market-display-solution-data-api/query-batch-bars-using-post.md", "data.GetDisplayBars", "SDK field name is `timespan` (not `interval`)."),
        ("Historical Bars (Single)", HK + "market-display-solution-data-api/bars-using-get.md", "data.GetDisplayBarsSingle", "SDK query param is `timespan` (not `interval`). Omits `last_time`, `real_time_required`, `trading_sessions`."),
        ("Tick", HK + "market-display-solution-data-api/tick-using-get.md", "data.GetDisplayTick", ""),
        ("Quotes Depth", HK + "market-display-solution-data-api/quotes-using-get.md", "data.GetDisplayDepth", ""),
        ("News Summary", HK + "market-display-solution-data-api/watchlist-summary-using-post.md", "data.GetDSNewsSummary", "SDK sends bare `[]string` body (not `{category_symbols, lang}`)."),
        ("Market News", HK + "market-display-solution-data-api/list-news-by-market-using-get.md", "data.GetDSMarketNews", "SDK sends only `category` query param (not `market`, `language`, `last_news_id`, `page_size`)."),
        ("Symbol News", HK + "market-display-solution-data-api/list-news-by-ticker-using-get.md", "data.GetDSSymbolNews", "SDK sends only `symbol` query param (not `category`, `language`, `last_news_id`, `page_size`)."),
        ("Latest News", HK + "market-display-solution-data-api/list-latest-news-using-get.md", "data.GetDSLatestNews", "SDK sends no query params (not `language`, `last_news_id`, `page_size`)."),
        ("Corporate Actions", HK + "market-display-solution-data-api/corp-action-using-get.md", "data.GetCorporateActions", ""),
        ("Corporate Actions By Market", HK + "market-display-solution-data-api/corp-market-using-get.md", "data.GetCorporateActionsByMarket", ""),
        ("Get Instruments", HK + "market-display-solution-data-api/list-using-get.md", "data.GetStockProfilesV3", ""),
        ("Batch Logos", HK + "market-display-solution-data-api/batch-logo-using-post.md", "data.GetLogos", "SDK sends `symbols` as query param with nil body (not `{category_symbols}` body). Response field is `logo` (not `logo_url`)."),
        ("Company Profile", HK + "market-display-solution-data-api/list-company-profile-using-get.md", "data.GetDSCompanyProfile", "SDK sends only `symbol` query param (not `category`)."),
        ("Analyst Target Price", HK + "market-display-solution-data-api/list-analyst-target-price-using-get.md", "data.GetDSAnalystTargetPrice", "SDK sends only `symbol` query param (not `category`)."),
        ("Analyst Rating", HK + "market-display-solution-data-api/list-analyst-rating-using-get.md", "data.GetDSAnalystRating", "SDK sends only `symbol` query param (not `category`)."),
        ("Streaming Subscribe", US + "broker-market-data-api/subscribe-using-post.md", "data.DSSubscribe", "US-site reference. SDK sends `{symbols: []string}` body (not `{session_id, category_symbols, sub_types, depth, overnight_required}`)."),
        ("Streaming Unsubscribe", US + "broker-market-data-api/unsubscribe-using-post.md", "data.DSUnsubscribe", "US-site reference. SDK sends `{symbols: []string}` body (not `{session_id, category_symbols, sub_types, unsubscribe_all}`)."),
    ],
)

AREAS["streaming"] = (
    "Streaming (MQTT)",
    "Real-time market data over MQTT. Subscribe/unsubscribe are HTTP calls that "
    "register the session; pushes arrive over MQTT. At most 5 concurrent "
    "connections per App Key.",
    [
        ("Subscribe", HK + "subscribe.md", "stream.Subscribe", "HTTP; registers the MQTT session."),
        ("Unsubscribe", HK + "unsubscribe.md", "stream.Unsubscribe", "HTTP; `UnsubscribeAll` clears everything."),
    ],
)

AREAS["events"] = (
    "Events (gRPC)",
    "Server-streaming gRPC subscriptions, signed with HMAC-SHA256 over the "
    "serialized request.",
    [
        ("Subscribe Trade Events", HK + "custom/subscribe-trade-events.md", "events.New / events.Run", "Order status changes."),
        ("Subscribe Position Events", US + "custom/subscribe-position-events.md", "events + SubscribePosition", "Position changes."),
        ("Subscribe Events (FD)", US + "fd-events/subscribe-events.md", "brokerfd/events", ""),
    ],
)

AREAS["connect-api"] = (
    "Connect API (OAuth)",
    "OAuth 2.0 authorization-code flow for third-party applications (US only).",
    [
        ("Authorization Code", US + "connect-api/get-authorization-code.md", "connect.AuthorizationURL", "Browser redirect URL builder."),
        ("Create Token", US + "connect-api/create-and-refresh-token.md", "connect.CreateToken", "Takes `TokenRequest` with fields: `GrantType` (`\"authorization_code\"` or `\"refresh_token\"`), `Code`, `RedirectURI`, `RefreshToken`."),
    ],
)

# section -> [ (label, official guide .md url) ]
GUIDE_SECTIONS = [
    ("Getting Started, About and SDKs", [
        ("Getting Started", HK_DOCS + "docs/getting-started.md"),
        ("About Webull", HK_DOCS + "docs/about.md"),
        ("About Webull OpenAPI", HK_DOCS + "docs/about-open-api.md"),
        ("SDKs and Tools", HK_DOCS + "docs/sdk.md"),
        ("Additional Resources", HK_DOCS + "docs/resources.md"),
    ]),
    ("Authentication", [
        ("Authentication Overview", HK_DOCS + "docs/authentication/overview.md"),
        ("Signature", HK_DOCS + "docs/authentication/signature.md"),
        ("Token", HK_DOCS + "docs/authentication/token.md"),
        ("Client Token (Display Solution)", US_DOCS + "docs/authentication/client-token.md"),
        ("Trading API Application", HK_DOCS + "docs/authentication/TradingAPIApplication.md"),
        ("Broker API Application", HK_DOCS + "docs/authentication/BrokerAPIapplication.md"),
    ]),
    ("Market Data API", [
        ("Market Data API Overview", HK_DOCS + "docs/market-data-api/overview.md"),
        ("Market Data API Getting Started", HK_DOCS + "docs/market-data-api/getting-started.md"),
        ("Data API", HK_DOCS + "docs/market-data-api/data-api.md"),
        ("Data Streaming API", HK_DOCS + "docs/market-data-api/data-streaming-api.md"),
        ("Subscribe Advanced Quotes", HK_DOCS + "docs/market-data-api/subscribe-quotes.md"),
        ("Hosted Display Solution", US_DOCS + "docs/market-data-api/hosted-display-solution.md"),
        ("Market Data API FAQ", HK_DOCS + "docs/market-data-api/faq.md"),
    ]),
    ("Trading API", [
        ("Trading API Overview", HK_DOCS + "docs/trade-api/overview.md"),
        ("Trading API Getting Started", HK_DOCS + "docs/trade-api/getting-started.md"),
        ("Accounts", HK_DOCS + "docs/trade-api/account.md"),
        ("Stock Trading", HK_DOCS + "docs/trade-api/stock.md"),
        ("Options Trading", HK_DOCS + "docs/trade-api/options.md"),
        ("Futures Trading", HK_DOCS + "docs/trade-api/futures.md"),
        ("Crypto Trading", US_DOCS + "docs/trade-api/crypto.md"),
        ("Event Contract Trading", US_DOCS + "docs/trade-api/event-contract.md"),
        ("Trading API FAQ", HK_DOCS + "docs/trade-api/faq.md"),
    ]),
    ("Broker API", [
        ("About Broker API", HK_DOCS + "docs/broker-api/about-broker-api.md"),
        ("Broker API Getting Started", HK_DOCS + "docs/broker-api/getting-started.md"),
        ("Event Contract Guidance", US_DOCS + "docs/broker-api/event-contract-guidance.md"),
    ]),
    ("Connect API", [
        ("About Connect API", US_DOCS + "docs/connect-api/about-connect-api.md"),
        ("OAuth Authentication", US_DOCS + "docs/connect-api/authentication.md"),
        ("Connect API Getting Started", HK_DOCS + "docs/connect-api/getting-started.md"),
    ]),
    ("Errors, FAQ and Changelog", [
        ("Error Codes", HK_DOCS + "docs/error-codes.md"),
        ("General FAQ", HK_DOCS + "docs/faq.md"),
        ("Documentation Changelog", HK_DOCS + "docs/changelog.md"),
    ]),
    ("Appendix — AI-Friendly Resources", [
        ("llms.txt", HK_DOCS + "docs/AI-friendly-Resources/llm.md"),
        ("Webull MCP Server", HK_DOCS + "docs/AI-friendly-Resources/mcp.md"),
        ("Webull Agent Skills", HK_DOCS + "docs/AI-friendly-Resources/skills.md"),
    ]),
]


# --------------------------------------------------------------------------
# Fetch / cache
# --------------------------------------------------------------------------
def fetch(url):
    os.makedirs(CACHE, exist_ok=True)
    safe = re.sub(r"[^A-Za-z0-9]+", "_", url)[-120:]
    path = os.path.join(CACHE, safe + ".md")
    if os.path.exists(path) and os.path.getsize(path) > 0:
        with open(path, "r", encoding="utf-8") as fh:
            return fh.read()
    req = urllib.request.Request(url, headers={"User-Agent": UA})
    with urllib.request.urlopen(req, timeout=45) as resp:
        data = resp.read().decode("utf-8", "replace")
    with open(path, "w", encoding="utf-8") as fh:
        fh.write(data)
    time.sleep(0.2)
    return data


def extract_json(md):
    m = re.search(r"```json\s*\n(.*?)\n```", md, re.S)
    if not m:
        return None
    try:
        return json.loads(m.group(1))
    except Exception:
        return None


# --------------------------------------------------------------------------
# Rendering helpers (SDK-mapped pages)
# --------------------------------------------------------------------------
def clean(text):
    if not text:
        return ""
    text = text.replace("<br/>", " ").replace("<br>", " ").replace("&bull;", "-")
    text = text.replace("&nbsp;", " ").replace("&minus;", "-").replace("&amp;", "&")
    text = re.sub(r"\[([^\]]+)\]\((?:/|\.\./|\./)[^)]*\)", r"\1", text)
    text = text.replace("|", "\\|")
    text = re.sub(r"\s+", " ", text).strip()
    return text


def type_of(schema):
    t = schema.get("type", "")
    if t == "array":
        it = schema.get("items", {})
        return "array<" + (it.get("type", "object")) + ">"
    return t or "object"


def short_enum(schema):
    e = schema.get("enum")
    if e and schema.get("type") == "string":
        return " — one of: " + ", ".join("`%s`" % v for v in e)
    return ""


def render_props(lines, props, required):
    req = set(required or [])
    nested = []
    for name, sch in props.items():
        desc = clean(sch.get("description", ""))
        lines.append("| `%s` | %s | %s | %s%s |" % (
            name, type_of(sch), "yes" if name in req else "", desc, short_enum(sch)))
        if sch.get("type") == "array" and sch.get("items", {}).get("type") == "object":
            sub = sch["items"].get("properties")
            if sub:
                nested.append((name, sub, sch["items"].get("required", [])))
        elif sch.get("type") == "object" and sch.get("properties"):
            nested.append((name, sch["properties"], sch.get("required", [])))
    for nm, sub, reqq in nested:
        lines.append("")
        lines.append("*Nested — `%s`:*" % nm)
        lines.append("")
        lines.append("| Field | Type | Required | Description |")
        lines.append("|---|---|---|---|")
        render_props(lines, sub, reqq)


def render_schema_block(title, schema):
    out = []
    if not schema:
        return out
    out.append("**%s**" % title)
    out.append("")
    if schema.get("type") == "array":
        items = schema.get("items", {})
        props = items.get("properties")
        if props:
            out.append("Array of objects:")
            out.append("")
            out.append("| Field | Type | Required | Description |")
            out.append("|---|---|---|---|")
            render_props(out, props, items.get("required", []))
        elif items:
            out.append("Array of `%s`." % type_of(items))
    elif schema.get("properties"):
        out.append("| Field | Type | Required | Description |")
        out.append("|---|---|---|---|")
        render_props(out, schema["properties"], schema.get("required", []))
    elif schema.get("type") == "string":
        out.append("Plain string.")
    out.append("")
    return out


SKIP_HEADERS = {
    "x-app-key", "x-app-secret", "x-timestamp", "x-signature-version",
    "x-signature-algorithm", "x-signature-nonce", "x-access-token",
    "x-version", "x-signature", "access_token", "reqid", "Accept",
    "Content-Type",
}


def sanitize_mdx(md):
    lines = []
    for line in md.splitlines():
        s = line.strip()
        if s.startswith("import ") or s.startswith("export "):
            continue
        if "@theme/" in s:
            continue
        if re.fullmatch(r"</?(Tabs|TabItem|Tab|figure|img|p|div|span)[^>]*>", s):
            continue
        if s in ("<Tabs>", "</Tabs>"):
            continue
        line = re.sub(r"^\s*<TabItem[^>]*>\s*$", "", line)
        line = re.sub(r"^\s*</TabItem>\s*$", "", line)
        line = re.sub(r'^#{1,6}\s+(.*)$', r'#### \1', line)
        lines.append(line)
    text = "\n".join(lines)
    text = re.sub(r"\[([^\]]+)\]\((?:/|\.\./|\./)[^)]*\)", r"\1", text)
    text = re.sub(r"\n{3,}", "\n\n", text)
    return text.strip()


def render_endpoint(title, url, sdk, note):
    md = fetch(url)
    spec = extract_json(md)
    out = ["## %s" % title, ""]
    if not spec:
        out.append("**Reference:** [%s](%s)" % (url.rsplit("/", 1)[-1], url))
        out.append("")
        if sdk:
            out.append("**SDK:** `%s`" % sdk)
            out.append("")
        if note:
            out.append("**Note:** %s" % note)
            out.append("")
        body = sanitize_mdx(re.sub(r"^#\s.*$", "", md, count=1, flags=re.M).strip())
        if body:
            out.append(body)
            out.append("")
        return "\n".join(out)

    method = spec.get("method", "").upper()
    path = spec.get("path", "")
    desc = clean(spec.get("description", ""))
    out.append("`%s %s`" % (method, path))
    out.append("")
    if desc:
        out.append("> %s" % desc)
        out.append("")
    out.append("| | |")
    out.append("|---|---|")
    out.append("| **SDK** | `%s` |" % sdk)
    out.append("| **Reference** | [%s](%s) |" % (url.rsplit("/", 1)[-1], url))
    if note:
        out.append("| **Note** | %s |" % note)
    out.append("")

    qp = [p for p in spec.get("parameters", [])
          if p.get("in") in ("query", "path") and p.get("name") not in SKIP_HEADERS]
    if qp:
        out.append("**Request — parameters**")
        out.append("")
        out.append("| Name | In | Type | Required | Description |")
        out.append("|---|---|---|---|---|")
        for p in qp:
            sch = p.get("schema", {})
            out.append("| `%s` | %s | %s | %s | %s%s |" % (
                p.get("name"), p.get("in"), type_of(sch),
                "yes" if p.get("required") else "",
                clean(p.get("description", "")), short_enum(sch)))
        out.append("")

    rb = spec.get("requestBody", {})
    if rb:
        schema = rb.get("content", {}).get("application/json", {}).get("schema")
        if schema:
            out += render_schema_block("Request body", schema)

    rschema = spec.get("responses", {}).get("200", {}).get("content", {}) \
                   .get("application/json", {}).get("schema")
    if rschema:
        out += render_schema_block("Response 200", rschema)

    out.append("**Errors** — `401` unauthorized, `417` business error, `500` "
               "server error. See [Errors](../errors.md).")
    out.append("")
    return "\n".join(out)


# --------------------------------------------------------------------------
# Verbatim helpers
# --------------------------------------------------------------------------
def sanitize_body(md, offset=2):
    """Verbatim-ish page body: MDX directives stripped, headings demoted by
    ``offset``, unresolvable internal links flattened, non-http images dropped.
    Fenced code blocks are preserved byte-for-byte."""
    lines = md.splitlines()
    out = []
    in_fence = False
    for line in lines:
        if line.lstrip().startswith("```"):
            in_fence = not in_fence
            out.append(line)
            continue
        if in_fence:
            out.append(line)
            continue
        s = line.strip()
        if s.startswith("import ") or s.startswith("export "):
            continue
        if "@theme/" in s:
            continue
        if re.fullmatch(r"</?(Tabs|TabItem|Tab)>", s) or re.fullmatch(r"<TabItem[^>]*>", s):
            continue
        if re.match(r"^:::\w*", s):
            continue
        m = re.match(r"^(#{1,6})(\s+)(.*)$", line)
        if m:
            line = "#" * (len(m.group(1)) + offset) + m.group(2) + m.group(3)
            out.append(line)
            continue
        line = re.sub(r"!\[[^\]]*\]\((?!https?://)[^)]*\)", "", line)
        line = re.sub(r"\[([^\]]+)\]\((?!https?://)[^)]*\)", r"\1", line)
        out.append(line)
    text = "\n".join(out)
    text = re.sub(r"\n{3,}", "\n\n", text).strip()
    return text


def render_verbatim_page(url):
    return sanitize_body(fetch(url))


# --------------------------------------------------------------------------
# Go source introspection (for reconciliation)
# --------------------------------------------------------------------------
GO_DIRS = ["data", "trade", "broker", "brokerfd", "display", "stream", "events", "client"]
PKG_DIR = {d: d for d in GO_DIRS}
_CONST_LINE_RE = re.compile(r'\b([A-Za-z_][A-Za-z0-9_]*)\s*=\s*"(/[^"]*)"')


def parse_consts():
    consts = {}
    for d in GO_DIRS:
        pkg = consts.setdefault(d, {})
        for root, _, files in os.walk(os.path.join(ROOT, d)):
            for f in files:
                if not f.endswith(".go") or f.endswith("_test.go"):
                    continue
                with open(os.path.join(root, f), encoding="utf-8") as fh:
                    for m in _CONST_LINE_RE.finditer(fh.read()):
                        pkg[m.group(1)] = m.group(2)
    return consts


def find_method_body(funcname, pkg):
    d = PKG_DIR.get(pkg)
    if not d:
        return None
    for root, _, files in os.walk(os.path.join(ROOT, d)):
        for f in files:
            if not f.endswith(".go") or f.endswith("_test.go"):
                continue
            with open(os.path.join(root, f), encoding="utf-8") as fh:
                txt = fh.read()
            m = re.search(r'func \([a-z] \*\w+\) ' + re.escape(funcname) + r'\(', txt)
            if not m:
                continue
            i = txt.find("{", m.end() - 1)
            depth = 0
            while i < len(txt):
                ch = txt[i]
                if ch == "{":
                    depth += 1
                elif ch == "}":
                    depth -= 1
                    if depth == 0:
                        return txt[m.start():i + 1]
                i += 1
    return None


_SKIP_DELEGATE = {"get", "post", "put", "delete", "do", "DoStream", "DisplayService",
                  "Core", "Close", "New", "EnsureToken"}


def resolve_method_path(sdk_func, consts, depth=0):
    if sdk_func in (None, "", "—", "not implemented", "not exposed"):
        return None, None
    first = (sdk_func or "").split(" ")[0].strip()
    parts = first.split(".")
    pkg, funcname = parts[0], parts[-1]
    body = find_method_body(funcname, pkg)
    if not body:
        return None, None
    for name, path in consts.get(pkg, {}).items():
        if re.search(r'[\(,]\s*' + re.escape(name) + r'\b', body):
            return path, name
    m = re.search(r'c\.\w+\(ctx,\s*"(/[^"]+)"', body)
    if m:
        return m.group(1), None
    if depth < 3:
        for m in re.finditer(r'c\.(\w+)\(', body):
            sub = m.group(1)
            if sub in _SKIP_DELEGATE:
                continue
            path, name = resolve_method_path(pkg + "." + sub, consts, depth + 1)
            if path:
                return path, name
    return None, None


def normalise_path(p):
    return re.sub(r"/\{[^}]+\}", "/{var}", p).rstrip("/").lower()


def norm_url(u):
    return (u.replace("https://developer.webull.hk//apis/docs/", "https://developer.webull.hk/apis/docs/")
             .replace("https://developer.webull.com//apis/docs/", "https://developer.webull.com/apis/docs/"))


def p(*parts):
    print(*parts, file=sys.stderr)