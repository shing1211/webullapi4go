# Todos — v02-trading-http

Single source of truth. Statuses: `todo`, `doing`, `blocked`, `review`, `done`.

| ID | Task | Role | Status | Depends On | Acceptance |
|----|------|------|--------|------------|------------|
| T9.1 | `trade` scaffold + types + guardrails + `/trading` v3 + accounts/assets | backend | done | v0.1 | sandbox #2 read-only: accounts+balance+11 positions OK |
| T9.2 | Tests: httptest + gated dedicated-sandbox accounts/assets | tester | done | T9.1 | delivered within T9.1 (unit + sandbox tests) |
| T9.3 | Trading docs page + accounts example | docs | done | T9.1 | docs/trading.md + examples/account; mkdocs --strict OK |
| T9.4 | Release v0.2.1 | release | T9.2, T9.3 | both remotes show tag |
| T10.1 | Order enums/types + validation | backend | todo | 0.2.1 | table-driven tests |
| T10.2 | Preview + Place (+ guardrails) | backend | T10.1 | live place | 
| T10.3 | Replace + Cancel | backend | T10.2 | live replace/cancel |
| T10.4 | Open/History/Detail queries + pagination | backend | T10.2 | query tests |
| T10.5 | Tests incl. live place→query→cancel | tester | T10.2–4 | cleanup guaranteed |
| T10.6 | Order docs + example | docs | T10.1–4 | documented commands run |
| T10.7 | Release v0.2.2 | release | T10.5, T10.6 | tag pushed |
| T11.1 | US/HK/CN rules + BCAN + tests | backend | T10.1 | table-driven tests |
| T11.2 | Docs + release v0.2.3 | docs/release | T11.1 | tag pushed |
| T12.1 | Options `SINGLE` orders | backend | T10.2 | unit + gated live |
| T12.2 | Docs + release v0.2.4 | docs/release | T12.1 | tag pushed |
| T13.1 | Combo orders (TP/SL, OTO, OCO, OTOCO) | backend | T10.1 | unit + gated live |
| T13.2 | Docs + release v0.2.5 | docs/release | T13.1 | tag pushed |
| T14.1 | Refactor news SSE through client pipeline | backend | v0.1 | tests + gated live |
| T14.2 | Release v0.2.6 | release | T14.1 | tag pushed |
