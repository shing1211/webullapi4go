# Next Phase — v0.4

## What completed this run
v0.3.0 added gRPC trade-event streaming end to end: vendored + attributed `events.proto`, a
HMAC-SHA256-capable signer, the `events` package (subscribe, typed order/position/option handlers,
reconnect), docs/example, a nightly read-only live CI job, and ADR-0002 promoted to Accepted.

## Gaps, tech debt, deferred items
- **Signer reconciliation**: `internal/auth` HMAC-SHA256 (uppercase body digest) is unused;
  `events/sign.go` implements the lowercase events variant. Unify under one algorithm-aware signer.
- **Display Solution API** — not started (client token, display quotes/news/screener/instruments).
- **Broker API** — not started (virtual accounts, funding, journals, events).
- **Go MCP server** — not started.
- Nightly CI secrets need configuring; consider a scheduled alert on failure.
- Position/option event schemas not verified against live payloads.
- Fundamentals coverage audit (financial statements, capital flow, dividends, earnings, SEC
  filings, fund data) still open.

## Candidate next-phase items
| # | Title | Objective | Why now | Effort | Deps | Risk |
|---|-------|-----------|---------|--------|------|------|
| 1 | Signer reconciliation | Unify HMAC algorithms/case under `internal/auth`; drop duplication | Removes debt before more surfaces | S | v0.3 | low |
| 2 | Display Solution API | Client token + display market data | Completes market-data story | L | v0.1 | med |
| 3 | Broker API | Virtual accounts, funding, journals, events | Institutional reach | L | v0.2/v0.3 | med |
| 4 | Go MCP server | Expose SDK as MCP tools | AI-native distribution | L | v1.0 | med |
| 5 | Broker gRPC events | Account/instrument/CA/trade/funding/journal/master-data streams | Reuses events infra | L | v0.3 | med |
| 6 | Fundamentals audit | Map/implement remaining documented endpoints | Completeness | M | v0.1 | low |
| 7 | Live CI hardening | Configure secrets; alert on nightly failure | Trust | S | v0.3 | low |

## Recommended next phase: **v0.4.0 — Display Solution API**
It is the largest remaining documented HTTP surface and shares the existing client/signing
machinery; fold in #1 (signer reconciliation) and #7 (CI secrets) as small prep tasks.

### Draft task breakdown
| ID | Objective | Role | Acceptance |
|----|-----------|------|-----------|
| T16.1 | Signer reconciliation under `internal/auth` | security | single signer supports SHA1+SHA256 (upper/lower body); events uses it |
| T16.2 | Display client-token create/refresh | backend | sandbox token pair obtained |
| T16.3 | Display quotes/snapshot/bars/tick | backend | sandbox data returned |
| T16.4 | Display screener/news/corporate actions/instruments/logos/profile/analyst | backend | endpoint tests |
| T16.5 | Tests + docs + example | tester/docs | documented commands run |
| T16.6 | Release v0.4.0 | release | tag on both remotes |

## Open questions
1. Display Solution uses a separate host (`co-branding-openapi…`) and a client-token (C2S) model —
   confirm sandbox access/credentials for it.
2. Prioritize Broker API gRPC events (#5) before Display Solution if institutional use matters more.
3. Should the nightly job post failure notifications (e.g. GitHub issue/Slack)? 
4. Reconcile the case-sensitivity by adding a lowercase mode, or keep `events/sign.go`?
