# Plan — v0.5 Live Probe

## Goal
Confirm (or correct) endpoint paths and response schemas for all v0.5 best-effort stubs by making authenticated real API calls against the Webull HK sandbox using shared public test accounts.

## Shared HK Test Accounts
| Account ID | App Key | App Secret |
|---|---|---|
| `V4H6R3L4VRI33UQ4TGR2NM1VI9` | `4b2b7acd2bf0d30d8aea173fceefa238` | `840b4353a6a31ce3ab91e2f99a510272` |

## Approach
Single Python script using standard library only (no dependencies). Why Python: no compilation → no AV false positives, user can audit signing logic.

## Files to create
- `examples/probe/probe.py` — self-contained probe script
- `examples/probe/README.md` — usage instructions
- `examples/probe/results/` — captured JSON responses

## Probe targets

### Priority 1: Corporate Actions (confirmed path from API ref)
- `GET /market-data/instruments/stocks/corporate-actions/list?symbols=AAPL`

### Priority 2: Fund Data (guessed paths — may be wrong)
- `GET /market-data/fund/AAPL/nav`
- `GET /market-data/fund/list`
- `GET /market-data/fund/AAPL/info`
- `GET /market-data/fund/AAPL/dividends`

### Priority 3: Crypto Data (guessed paths — may be wrong)
- `GET /market-data/crypto/BTCUSD/bars`
- `GET /market-data/crypto/BTCUSD/snapshot`
- `GET /market-data/crypto/ETHUSD/tick`

### Priority 4: Screener v2 (guessed path — may be wrong)
- `POST /wlas/screener/ng/query` with filter body

### Priority 5: Broker FD (guessed paths — returned 404 before)
- `GET /broker-fd/accounts`
- `GET /broker-fd/positions`

## Probe workflow
1. Get access token (HMAC-SHA1 signed POST to `/openapi/auth/token/create`)
2. For each endpoint: sign request, send HTTP call, capture response
3. Output: status code + raw JSON → `results/{category}_{endpoint}.json`
4. If 404: try alternate path guesses before moving on

## Rate limit handling
- Token endpoint: 10 req / 30s — sleep 4s between token requests
- General: sleep 1s between requests to avoid hitting limits
- Use `category=US` with `AAPL` symbol for equity endpoints

## Verification
- `python examples/probe/probe.py` runs without errors
- JSON files captured in `examples/probe/results/`
- Compare field names against Go struct field names
- Update `findings.md` with confirmed/corrected paths and schemas

## Risks
- Sandbox market data limited to AAPL
- Crypto may not be available in sandbox
- Fund data endpoints may not exist in HK sandbox
- Broker FD paths may differ between HK and US sandboxes
