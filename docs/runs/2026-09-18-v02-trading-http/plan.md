# Plan — webullapi4go v0.2 (Trading HTTP), patch-staged

## Run identity
- Date: 2026-09-18
- Slug: v02-trading-http
- Mode: BUILD
- Repo: `D:\github\webullapi4go`
- Base: v0.1.0 (`ee30d6f`)
- Module: `github.com/shing1211/webullapi4go`
- Remotes: GitHub `origin`, Gitee `gitee`

## Goal
Add Trading HTTP support (accounts, assets, and order lifecycle) as a public `trade`
package, verified against a sandbox account, released in staged patch versions
`0.2.1 … 0.2.6`.

## Decisions (locked)
1. **Dedicated sandbox account for trading tests**: documented shared Account #2
   (`OGG4RRLC6EDE98HI920KRBVSKB`) so Account #1 stays reserved for market-data smoke
   tests. Values are env-only and must never be committed.
2. **Patch-staged releases** `0.2.1`, `0.2.2`, … each pushed to GitHub + Gitee.
3. **Order guardrails**: opt-in `MaxOrderNotional` / `MaxOrderQuantity`.
4. **Package name**: `trade`.
5. **`x-version`**: `/trading/**` defaults to `v3` (documented current); `/market-data/**`
   stays `v2`; overridable.

## Architecture
- Public `trade` package: `trade.New(cl *client.Client, opts ...Option) *Client`,
  grouped files (`client.go`, `accounts.go`, `orders.go`, `query.go`, `options.go`,
  `combo.go`, `types.go`), all requests via `client.Do`.
- Reuses v0.1 `client` (signing, token, retry/ratelimit/breaker, multi-region).
- Numeric values stay strings; typed enums; GoDoc; Apache-2.0; no `internal/*` in public
  signatures. Mirror `data` package conventions.
- Client change (additive): per-path default API version so `/trading/**` → `v3`.
- Guardrails enforced at order build time (preview/place), not in generic `Do`.

## Reference endpoints (confirm exact path per task)
| Area | Method / Path | Response |
|---|---|---|
| Accounts | `GET /trading/accounts/list` | `[]Account{account_id, account_number, account_type, account_class}` |
| Balance | `GET /trading/assets/balances/get?account_id=` | `AssetsBalanceResult` |
| Positions | `GET /trading/assets/positions/list?account_id=` | `[]AssetsPositionResult` |
| Place | `POST /trading/orders/place` | `{client_order_id, order_id}` |
| Preview | `POST /trading/orders/preview` (confirm) | estimate/costs |
| Replace | `POST /trading/orders/replace` (confirm) | result |
| Cancel | `POST /trading/orders/cancel` (confirm) | result |
| Open orders | `GET /trading/orders/open-orders/list` | `{data, pagination_key}` |
| History | `GET /trading/orders/history/...` (confirm) | paginated |
| Detail | `GET /trading/orders/detail/...` (confirm) | order |

Statuses: `PENDING, SUBMITTED, CANCELLED, FILLED, FAILED, PARTIAL_FILLED`.
Combo types: `NORMAL, MASTER, STOP_PROFIT, STOP_LOSS, OTO, OCO, OTOCO` (OTOx EQUITY-only).

## Patch 0.2.1 — Foundation + Accounts/Assets
| ID | Objective | Role | Depends | Acceptance | Size |
|----|-----------|------|---------|-----------|------|
| T9.1 | `trade` scaffold + shared types/enums + guardrail options + `/trading` v3 default; implement `ListAccounts`, `GetBalance`, `GetPositions` | backend | v0.1 | build/lint clean; httptest + sandbox read-only OK | M |
| T9.2 | Tests: httptest + env-gated dedicated-sandbox read-only integration | tester | T9.1 | unit passes offline; sandbox test gated | S |
| T9.3 | Trading docs page + accounts/positions example | docs | T9.1 | documented commands run | S |
| T9.4 | Release v0.2.1 | release | T9.2, T9.3 | both remotes show tag | S |

## Patch 0.2.2 — Stock orders
| ID | Objective | Role | Depends | Acceptance | Size |
|----|-----------|------|---------|-----------|------|
| T10.1 | Order enums/types + validation | backend | 0.2.1 | table-driven type tests | M |
| T10.2 | `PreviewOrder`, `PlaceOrder` (+ guardrails) | backend | T10.1 | httptest + live preview/place | L |
| T10.3 | `ReplaceOrder`, `CancelOrder` | backend | T10.2 | live replace/cancel+cleanup | M |
| T10.4 | `GetOpenOrders` (paginated), `GetOrderHistory`, `GetOrderDetail` | backend | T10.2 | query tests | M |
| T10.5 | Tests incl. env-gated live place→query→cancel | tester | T10.2–4 | cleanup guaranteed | M |
| T10.6 | Order docs + example | docs | T10.1–4 | documented commands run | S |
| T10.7 | Release v0.2.2 | release | T10.5, T10.6 | both remotes show tag | S |

## Patch 0.2.3 — Market-specific rules
| ID | Objective | Role | Depends | Acceptance | Size |
|----|-----------|------|---------|-----------|------|
| T11.1 | US/HK/CN order-type matrices, HK BCAN `no_party_ids`, A-share notes + tests | backend | T10.1 | table-driven tests | M |
| T11.2 | Docs + release v0.2.3 | docs/release | T11.1 | docs run; tag pushed | S |

## Patch 0.2.4 — Options orders
| ID | Objective | Role | Depends | Acceptance | Size |
|----|-----------|------|---------|-----------|------|
| T12.1 | `SINGLE` option orders (LIMIT/STOP_LOSS/STOP_LOSS_LIMIT), sell-side DAY rule | backend | T10.2 | unit + gated live | M |
| T12.2 | Docs + release v0.2.4 | docs/release | T12.1 | tag pushed | S |

## Patch 0.2.5 — Combo orders (US only)
| ID | Objective | Role | Depends | Acceptance | Size |
|----|-----------|------|---------|-----------|------|
| T13.1 | TP/SL, OTO, OCO, OTOCO with group/leg-count validation | backend | T10.1 | unit + gated live | L |
| T13.2 | Docs + release v0.2.5 | docs/release | T13.1 | tag pushed | S |

## Patch 0.2.6 — v0.1 debt
| ID | Objective | Role | Depends | Acceptance | Size |
|----|-----------|------|---------|-----------|------|
| T14.1 | Refactor `data` news SSE through the client pipeline (token/version/resilience) | backend | v0.1 | tests + gated live | M |
| T14.2 | Release v0.2.6 | release | T14.1 | tag pushed | S |

## Test & credential strategy
- Env-gated, never committed:
  `WEBULL_TRADE_SANDBOX=1`, `WEBULL_TRADE_APP_KEY`, `WEBULL_TRADE_APP_SECRET`,
  `WEBULL_TRADE_ACCOUNT_ID`.
- Mutation tests additionally gated by `WEBULL_TRADE_MUTATE=1`; use tiny AAPL limit orders
  away from market; always cancel/cleanup in `t.Cleanup`.
- Default `go test ./...` stays offline (httptest).
- Documented shared sandbox accounts are public but must still never be written to files.

## Risks & mitigations
1. **Path/schema uncertainty** (preview/replace/cancel/history/detail) → each task fetches
   its reference page first and posts both v2/v3 where relevant.
2. **Shared account interference** → dedicated Account #2 for trading.
3. **Rate limits** (US 15/s, HK/CN 1/s; token 10/30s) → space live tests; reuse token.
4. **HK BCAN / A-share disabled** → validated and documented, not silently sent.
5. **Partial release failures** → no force-push; stop and report.

## Verification
- `go build ./...`, `go vet ./...`, `gofmt -l .`, `golangci-lint run ./...` clean.
- `go test ./... -count=1` passes (offline).
- Env-gated sandbox tests pass when configured.
- `mkdocs build --strict` still succeeds.

## Order of work
T9.1 → T9.2 → T9.3 → T9.4 (v0.2.1) → T10.1 → T10.2 → (T10.3, T10.4) → T10.5 → T10.6 →
T10.7 (v0.2.2) → T11 → T12 → T13 → T14.
