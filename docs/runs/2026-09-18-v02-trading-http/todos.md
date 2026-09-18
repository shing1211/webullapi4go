# Todos — v02-trading-http

Single source of truth. Statuses: `todo`, `doing`, `blocked`, `review`, `done`.

| ID | Task | Role | Status | Depends On | Acceptance |
|----|------|------|--------|------------|------------|
| T9.1 | `trade` scaffold + types + guardrails + `/trading` v3 + accounts/assets | backend | done | v0.1 | sandbox #2 read-only: accounts+balance+11 positions OK |
| T9.2 | Tests: httptest + gated dedicated-sandbox accounts/assets | tester | done | T9.1 | delivered within T9.1 |
| T9.3 | Trading docs page + accounts example | docs | done | T9.1 | docs/trading.md + examples/account; mkdocs --strict OK |
| T9.4 | Release v0.2.1 | release | done | T9.2, T9.3 | `5958dcf`; tag v0.2.1 on both remotes |
| T10.1 | Order enums/types + validation + Preview + Place (+ guardrails) | backend | done | 0.2.1 | live preview OK (est. cost 180.00) |
| T10.2 | Preview + Place (+ guardrails) | backend | done | T10.1 | absorbed into T10.1 |
| T10.3 | Replace + Cancel | backend | done | T10.2 | live place→cancel; `/trading/orders/{replace,cancel}` |
| T10.4 | Open/History/Detail queries + pagination | backend | done | T10.2 | live: open-orders/list, historical-orders/list, /orders/get |
| T10.5 | Tests incl. live place→query→cancel | tester | done | T10.2–4 | delivered within T10.1/3/4 |
| T10.6 | Order docs + example | docs | done | T10.1–4 | docs/trading.md + examples/order; mkdocs --strict OK |
| T10.7 | Release v0.2.2 | release | done | T10.5, T10.6 | `025588b`; tag v0.2.2 on both remotes |
| T11.1 | US/HK/CN rules + BCAN + tests | backend | done | T10.1 | 74 rule subtests; live HK preview OK |
| T11.2 | Docs + release v0.2.3 | docs/release | done | T11.1 | `d9ba14a`; tag v0.2.3 on both remotes |
| T12.1 | Options `SINGLE` orders | backend | done | T10.2 | rules + tests; live option preview OK |
| T12.2 | Docs + release v0.2.4 | docs/release | done | T12.1 | `b08b83c`; tag v0.2.4 on both remotes |
| T13.1 | Combo orders (TP/SL, OTO, OCO, OTOCO) | backend | done | T10.1 | composition rules + tests; live combo preview OK |
| T13.2 | Docs + release v0.2.5 | docs/release | doing | T13.1 | tag pushed |
| T14.1 | Refactor news SSE through client pipeline | backend | todo | v0.1 | tests + gated live |
| T14.2 | Release v0.2.6 | release | todo | T14.1 | tag pushed |
