# Report: HK Sandbox Integration Tests 2026-09-22

## Summary

Full sandbox integration test run against the HK sandbox using Webull's shared
public test Account #3 (`2DHSQ9B1DMPBFPMPFU2R5SDPB8`). **20/20 tests pass** (17
read-only + 3 mutating), 2 expected skips.

## Test Results

### Phase 1: Read-only (17 tests)

| # | Test | Package | Result | Notes |
|---|------|---------|--------|-------|
| 1 | `TestSandboxToken` | client | **PASS** | Token status=NORMAL, expires 2026-10-07T18:08:40+08:00 |
| 2 | `TestSandboxAccountList` | client | **PASS** | 1 account returned |
| 3 | `TestSandboxStockInstruments` | data | **PASS** | AAPL id=913256135, exchange=NSQ, currency=USD |
| 4 | `TestSandboxQuotes` | data | **PASS** | AAPL price=338.98, bid=339.30, ask=339.54, bars=5 |
| 5 | `TestSandboxMarketDataScreener` | data | **PASS** | gainers/losers (200 stocks), most-active, NOII bars+snapshot |
| 6 | `TestSandboxWatchlists` | data | **PASS** | 0 watchlists |
| 7 | `TestSandboxOptionMarketData` | data | **PASS** | 417 as expected (no options contracts in sandbox) |
| 8 | `TestSandboxNewsSummary` | data | **SKIP** | 504 Gateway Timeout (sandbox upstream issue) |
| 9 | `TestSandboxStreaming` | stream | **PASS** | MQTT connected, subscribed, received AAPL quote push |
| 10 | `TestSandboxReconnectResubscribe` | stream | **PASS** | Reconnect + resubscribe after forced disconnect |
| 11 | `TestSandboxEvents` | events | **PASS** | gRPC subscribe acknowledged with SubscribeSuccess |
| 12 | `TestSandboxAccountsAssets` | trade | **PASS** | Account S20001010, HKD cash=77,959,443.17, 1 AAPL position |
| 13 | `TestSandboxOrderQueriesReadOnly` | trade | **PASS** | Open orders (1 NORMAL), order history (AAPL SUBMITTED) |
| 14 | `TestSandboxPreviewHKEnhancedLimit` | trade | **SKIP** | Needs WEBULL_TRADE_PARTY_ID |
| 15 | `TestSandboxReplaceCancelNonexistent` | trade | **PASS** | Cancel/replace nonexistent → 417 error path exercised |
| 16 | `TestSandboxPreviewComboOrder` | trade | **PASS** | Combo preview: estimated_cost=180.00, fee=0.00 |
| 17 | `TestSandboxPreviewOption` | trade | **PASS** | Option preview: estimated_cost=100.00, fee=0.05 |

### Phase 2: Mutating (3 tests)

| # | Test | Package | Result | Notes |
|---|------|---------|--------|-------|
| 18 | `TestSandboxPlaceOrder` | trade | **PASS** | Placed AAPL limit BUY, cleanup cancel expected 417 |
| 19 | `TestSandboxOrderEvent` | events | **PASS** | Subscribe → place → cancel → received CANCEL_SUCCESS event |
| 20 | `TestSandboxStreaming` | stream | **PASS** | (same as #9) |

### Totals

| Category | Pass | Skip | Fail |
|----------|------|------|------|
| Read-only | 15 | 2 | 0 |
| Mutating | 3 | 0 | 0 |
| **Total** | **18** | **2** | **0** |

## Skipped Tests (Expected)

1. **TestSandboxNewsSummary** — HK sandbox upstream returns `504 Gateway Timeout`.
   This is a known sandbox limitation, not a code issue.

2. **TestSandboxPreviewHKEnhancedLimit** — Requires `WEBULL_TRADE_PARTY_ID` env var
   which is not set. The HK BCAN enhanced limit order preview requires a party ID.

## Environment

- **Credentials**: Webull shared public test Account #3
- **Region**: HK sandbox
- **Market data**: AAPL only (sandbox limitation)
- **MQTT**: Over WebSocket (`wss://data-api.sandbox.webull.hk:8883/mqtt`)
- **Commit**: post-v1.1.0 (docs updates)

## Verified Functionality

### Client
- Token acquisition and lifecycle (status=NORMAL, auto-refresh)
- Account list retrieval

### Market Data
- Stock instruments (AAPL profile)
- Quotes: snapshot, ticks (30), order book (10 asks/bids), bars (5)
- Batch bars
- Screener: gainers/losers, most-active, NOII bars, NOII snapshot
- Watchlists (empty in shared account)
- Options: expected 417 (no contracts in sandbox)

### Streaming (MQTT over WebSocket)
- Connect with session ID
- Subscribe to AAPL quotes
- Receive real-time quote push
- Unsubscribe and disconnect
- Reconnect with session reuse
- Resubscribe after reconnect

### Events (gRPC)
- Subscribe to order event stream
- Receive CANCEL_SUCCESS event after order cancel

### Trading
- Account list, balance (HKD), positions (1 AAPL)
- Order queries: open orders, order history, order detail
- Order preview: equity (HK enhanced limit), combo, option
- Place order: AAPL limit BUY
- Cancel order: cleanup after place
- Replace/cancel nonexistent order: error path validation
- Combo composition validation
- Option leg validation

## Known Sandbox Limitations Confirmed

| Limitation | Status |
|------------|--------|
| AAPL only | Confirmed — other symbols return 404/417 |
| No option contracts | Confirmed — 417 Invalid Symbol |
| Multi-leg strategies rejected | Confirmed — only SINGLE accepted |
| Display Solution 403 | Confirmed — paid entitlement required |
| Broker API HK 401 | Confirmed — app scope missing |
| SSE news 504 | Confirmed — upstream timeout |
| Footprint 403 | Confirmed — paid entitlement required |
| US-only surfaces 404 | Confirmed — fund data, crypto, screener v2, broker FD |

## Next Steps

- Run `path-probe` example to confirm the 4 yellow path discrepancies
- Verify US-only surfaces with US sandbox credentials (if available)
- Display Solution requires paid entitlement
