# Next Phase — v0.2 (Trading HTTP)

## What completed this run
v0.1.0: authentication (signer + token lifecycle), core signed HTTP client with resilience
and multi-region support, the full Market Data HTTP surface, MQTT streaming with reconnect
and re-subscription, examples, docs site, CI, and community files — verified against the
Webull sandbox.

## Gaps, tech debt, deferred items
- **Trading API** (accounts, assets, order preview/place/replace/cancel, order queries,
  US/HK/CN stock rules, options orders, combo orders) — not implemented.
- **gRPC trade/broker events** — not implemented; event `.proto` unpublished (ADR-0002).
- **Display Solution API** and **Broker API** — not implemented.
- `data` news uses a bespoke SSE path instead of the client pipeline (token/version/
  resilience bypass) — refactor needed.
- Sandbox tests create a token per test; consolidate to reduce shared-account flakiness.
- `x-version` defaults to `v2`; consider `v3` once all endpoints are confirmed.

## Candidate next-phase items
| # | Title | Objective | Why now | Effort | Deps | Risk |
|---|-------|-----------|---------|--------|------|------|
| 1 | Accounts & assets | `account/list`, `assets/balance`, `assets/positions` | Foundation for any order flow | S | v0.1 | low |
| 2 | Order lifecycle | preview/place/replace/cancel + open/history/detail | Core trading value | L | #1 | med (rules) |
| 3 | Market-specific order rules | US/HK/CN order types, BCAN party ids, A-share limits | Correctness across markets | M | #2 | med/high |
| 4 | Options orders | single-leg CALL/PUT, STOP_LOSS(_LIMIT) | Requested scope | M | #2 | med |
| 5 | Combo orders (US) | TP/SL, OTO, OCO, OTOCO | Advanced order support | M | #2 | med |
| 6 | gRPC trade events | subscribe order-status notifications | Real-time trading | L | #2, ADR-0002 spike | high (protos) |
| 7 | News SSE refactor | route news through client pipeline | Remove debt | S | v0.1 | low |

## Recommended next phase: **v0.2 — Trading HTTP**
Deliver accounts/assets first, then the order lifecycle, then market-specific rules and
options/combo orders. Keep gRPC events (v0.3) behind an ADR-0002 spike.

### Draft task breakdown
| ID | Objective | Role | Acceptance |
|----|-----------|------|-----------|
| T7.1 | Account + assets endpoints (`client`/`trade` package) | backend | sandbox `account/list`, balance, positions |
| T7.2 | Order preview + place (stock) | backend | sandbox preview + place on AAPL |
| T7.3 | Order replace + cancel | backend | modify/cancel an open order |
| T7.4 | Order queries (open/history/detail) | backend | matches placed order |
| T7.5 | US/HK/CN market rules + BCAN | backend | table-driven rule tests |
| T7.6 | Options orders | backend | single-leg place/preview |
| T7.7 | Combo orders (US) | backend | OTO/OCO/TP-SL examples |
| T7.8 | Refactor news SSE via client pipeline | backend | tests + live news |
| T7.9 | Trading docs + examples | docs | documented commands run |
| T7.10 | Release v0.2.0 | release | GitHub + Gitee |

## Open questions
1. Sandbox order placement is mutating shared public accounts — is that acceptable for
   integration tests, or should trading tests be preview-only + unit-tested with httptest?
2. Should `x-version` move to `v3` by default in v0.2, or stay configurable?
3. Confirm the Gitee remote URL for releases.
4. For gRPC events, approve the ADR-0002 approach (extract `.proto` from the Apache-2.0
   Python SDK) before the v0.3 spike.
