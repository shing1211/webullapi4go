# Plan — v0.5 Phase 2: Fund Data, Crypto Data, Screener v2

## Goal
Ship best-effort stubs for Fund Data, Crypto Data, and Screener v2 endpoints
in `data/`, matching existing package conventions. All three use `map[string]string`
for unconfirmed fields; schemas will be refined when live probe credentials are
available.

## Scope

### T1 — Fund Data (`data/fund_data.go`)
**Endpoints to implement (guessed from Python SDK / API reference patterns):**
- `GET /market-data/fund/{symbol}/nav` — ETF/fund NAV history
- `GET /market-data/fund/{symbol}/info` — Fund basic info
- `GET /market-data/fund/{symbol}/dividends` — Fund dividend history
- `GET /market-data/fund/list` — Fund list by market

**Approach:** Use `map[string]string Extra` on each struct for unconfirmed fields.
ETF/fund fields differ from equities — expected fields: `nav`, `nav_date`,
`aum`, `expense_ratio`, `dividend_yield`, `inception_date`, etc.

### T2 — Crypto Data (`data/crypto_data.go`)
**Endpoints to implement (guessed from Webull API reference patterns):**
- `GET /market-data/crypto/{symbol}/bars` — Candlestick data
- `GET /market-data/crypto/{symbol}/tick` — Tick data
- `GET /market-data/crypto/{symbol}/depth` — Market depth / order book
- `GET /market-data/crypto/{symbol}/snapshot` — Crypto snapshot

**Approach:** Mirror existing `bars.go`, `tick.go`, `snapshot.go` patterns for
crypto tickers. Expected fields: `symbol`, `price`, `bid`, `ask`, `volume`,
`open`, `high`, `low`, `close`, `turnover`, etc.

### T3 — Screener v2 (`data/screener_v2.go`)
**Endpoint:**
- `POST /wlas/screener/ng/query` — Screener v2 (next-gen), POST with filter body

**Approach:** Follow `screener.go` conventions. Screener v2 uses a POST body with
filter criteria. Struct fields for filter options with `map[string]string Extra`
for unconfirmed fields. Response: array of `ScreenerStock` (re-use existing type).

## Approach
Best-effort schema only — no live probe. All paths guessed from API reference
patterns and Python SDK. Refine with live data in a follow-up spike.

## Risks
- Paths may be wrong (e.g., `/wlas/` prefix vs `/market-data/`)
- Response schema fields unknown — `map[string]string` fallback on every struct
- Crypto endpoints likely require different market/category from equities

## Order of Work
T1, T2, T3 run in parallel (independent packages). Verification after all three.

## Verification
```
go build ./...
go vet ./...
gofmt -l .
go test ./...
golangci-lint run ./...
mkdocs build --strict
```
