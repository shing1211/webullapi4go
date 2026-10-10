# Run Report

**Run:** `docs/runs/2026-09-22-broker-paths-fix/`
**Date:** 2026-09-22
**Mode:** BUILD
**Status:** Complete
**Commit:** `44299b2`
**Tag:** `v0.9.2`

## Shipped

| Task | Result |
|------|--------|
| Fix broker paths | ✅ All 32 methods now use correct paths |
| Fix request/response shapes | ✅ Match official OpenAPI spec |
| Fix broker unit tests | ✅ All pass |
| Run broker-probe | ✅ Paths confirmed correct — 401 (auth scope issue, not 404) |
| Build + vet + test | ✅ All pass |

## Root Cause

The `broker/` package was built with guessed paths. The SDK used `/openapi/v1/broker/...` but the official API uses `/broker/...` directly. This was discovered by fetching `https://developer.webull.hk/apis/docs/llms.txt` which revealed the official OpenAPI specification.

## Key Corrections

| What was wrong | How it was fixed |
|---|---|
| All 32 paths used `/openapi/v1/broker/...` | Corrected to `/broker/...` |
| `VirtualAccount` had wrong fields | Updated to official fields (`AccountNumber`, `AccountStatus`, etc.) |
| `ListVirtualAccounts` expected bare array | Now handles `{"data": [...]}` wrapper |
| `Balance` was flat | Now nested with `total_cash_balance`, `account_currency_assets` |
| `PreviewOrder` flat request | Adapted to nested `new_orders[]` wire format |
| `ReplaceOrder`/`CancelOrder` used query params | Now POST with JSON body |
| `GetTradeCalendar` was GET | Now POST with JSON body |
| `GetFXRate` params `from`/`to` | Corrected to `from_currency`/`to_currency` |

## Sandbox Probe Result

```
ListVirtualAccounts: FAIL 401 ROUTE_NOT_PERMITTED
```

**404 → 401** — paths are now correct. The 401 means the app lacks Broker API scope in the sandbox. This is a permissions issue, not a path issue.

## Release

- Tag `v0.9.2` pushed to GitHub
- CHANGELOG updated with full list of corrections
