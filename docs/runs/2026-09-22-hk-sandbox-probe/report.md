# Report: HK Sandbox Probe 2026-09-22

## Summary
HK sandbox probe completed 2026-09-22. Key findings: futures product codes path confirmed,
futures instruments `Unit` field needs flexible type (live data sends numeric), multi-leg
option strategies fail in HK sandbox (strategy validation rejects anything but SINGLE),
and Broker API HK returns `401 ROUTE_NOT_PERMITTED` due to missing app scope.

## Confirmed Working (HK Sandbox)
- `GetFuturesProductCodes` — PASS — path `/trading/instruments/futures/product-codes/list`
  returns 89 products for HK_FUTURES
- `GetFuturesProductClasses` — path confirmed
- Watchlist CRUD — PASS when not rate-limited (429) or auth-limited (401)

## Bugs Found

### 1. FuturesInstrument.Unit type mismatch
**File**: `data/futures.go:111`
**Symptom**: `json: cannot unmarshal number into Go struct field FuturesInstrument.unit of type string`
**Cause**: HK sandbox sends numeric `unit` (e.g., `1`) instead of string (e.g., `"1-index points"`)
**Fix**: Add `StringOrNumber` flexible type with `UnmarshalJSON`

### 2. client_order_id length overflow in options-multi-leg example
**File**: `examples/options-multi-leg/main.go:193`
**Symptom**: `client_order_id must be at most 32 characters, got 34–40`
**Cause**: `time.Now().UnixNano()` produces 19 digits; combined with "probe-STRATEGY-" prefix
  yields 34–40 characters
**Fix**: Use `time.Now().Unix()` (10 digits) instead of `UnixNano()`

## Confirmed HK Sandbox Limitations

### Multi-leg option strategies
HK sandbox rejects all multi-leg strategies except SINGLE:
- VERTICAL → 417 `invalid option_strategy`
- STRADDLE → 417 `invalid option_strategy`
- STRANGLE → 417 `invalid option_strategy`
- IRON_CONDOR → 417 `invalid option_strategy`
- IRON_BUTTERFLY → 417 `invalid option_strategy`
- BUTTERFLY → 417 `invalid option_strategy`
- CALENDAR → 417 `invalid option_strategy`
- DIAGONAL → 417 `invalid option_strategy`
- RATIO → 417 `option_strategy only support SINGLE`
- COLLAR → 417 `invalid option_strategy` (HK); confirmed working in US sandbox from prior run

### Broker API HK
Returns `401 ROUTE_NOT_PERMITTED` — app lacks required scope. Not a path issue.

### US-only surfaces
Return 404 in HK sandbox:
- Fund data endpoints
- Crypto data endpoints
- Screener v2
- Broker FD endpoints
- Instrument v3/logos
- Crypto category → 417

### Display Solution
Returns `403 Forbidden` in HK sandbox — paid entitlement required (not a path issue).

### SSE News
Upstream returns `504 Gateway Timeout` in HK sandbox.

## Confirmed Futures Paths
- `GetFuturesProductCodes` — PASS — path confirmed
- `GetFuturesProductClasses` — PASS — path confirmed
- `GetFuturesInstruments` — FAIL (Unit type bug only)

## Open TODOs Updated
- `TODO(t8)` in `trade/types.go:158`: HK sandbox confirms multi-leg strategies are
  rejected except SINGLE; strategy wire values remain provisional for multi-leg
- `TODO(t8)` in `trade/options.go:40,239`: multi-leg order-type matrix and structural
  rules unconfirmed for HK; US sandbox needed for full confirmation
- `TODO(futures)` in `data/futures_market.go`: market data paths (tick, snapshot, bars,
  depth, footprint) remain unconfirmed; product codes path confirmed

## Next Steps
- US sandbox credentials needed to confirm multi-leg strategy wire values and
  futures market data paths
- v0.10 release blocked on US sandbox access
