# Next Phase — v0.3

## What completed this run
v0.2.x delivered the full Trading HTTP surface for stocks, options, and US combo orders,
plus market rules, HK BCAN, order guardrails, `client.DoStream`, and a news refactor —
shipped as six patch releases across GitHub and Gitee.

## Gaps, tech debt, deferred items
- **gRPC trade events** (order status push) — the headline remaining v0.x gap. Event
  `.proto` is unpublished (ADR-0002 recommends deriving from the Apache-2.0 Python SDK).
- **Display Solution API** and **Broker API** — not started.
- Remaining fundamentals endpoints exposed by Webull's MCP tool list (financial statements,
  capital flow, dividends, earnings, SEC filings, fund data) — some may not be in the HK
  OpenAPI reference; needs a coverage audit.
- Sandbox test flakiness on shared accounts (watchlist 401; token churn) — consolidate to a
  single cached token per test package.
- No CI job runs the env-gated live suite; it is manual.

## Candidate next-phase items
| # | Title | Objective | Why now | Effort | Deps | Risk |
|---|-------|-----------|---------|--------|------|------|
| 1 | gRPC event proto spike | Extract/derive trade-event `.proto` per ADR-0002 | Unblocks real-time trading | M | — | high (schemas) |
| 2 | gRPC trade events client | Subscribe to order-status events (filled/cancelled/…) | Core v0.3 feature | L | #1 | high |
| 3 | Display Solution API | Client token + display quotes/news/screener/instruments | Broadens market-data reach | L | v0.1 | med |
| 4 | Broker API | Virtual accounts, funding, journals, events | Institutional use | L | v0.2 | med |
| 5 | Fundamentals coverage audit | Map Webull MCP/fundamental endpoints vs HK OpenAPI; implement gaps | Fills doc-completeness gaps | M | v0.1 | low |
| 6 | Test reliability + live CI gate | Shared token per package; optional scheduled live job | Trust + regression safety | S | v0.2 | low |
| 7 | MCP server (Go) | Expose the SDK as MCP tools | AI-native distribution | L | v1.0 | med |

## Recommended next phase: **v0.3.0 — gRPC trade events**
Start with the ADR-0002 spike (#1), then the events client (#2); fold in test reliability
(#6) so live coverage is trustworthy.

### Draft task breakdown
| ID | Objective | Role | Acceptance |
|----|-----------|------|-----------|
| T15.1 | Spike: obtain/derive trade-event `.proto`; generate Go; document provenance | data | proto compiles; decode a sample; NOTICE attribution |
| T15.2 | `events` package: connect to `events-api.<domain>` with signing; subscribe by account | backend | connects (sandbox) |
| T15.3 | Typed order-event handlers + reconnect/resubscribe | backend | live order event received (mutate-gated) |
| T15.4 | Tests + docs + example | tester/docs | documented commands run |
| T15.5 | Release v0.3.0 | release | tag pushed to both remotes |

## Open questions
1. Confirm the ADR-0002 approach (derive `.proto` from the Python SDK) before T15.1.
2. Is the sandbox `events-api.sandbox.webull.hk` reachable from the build network?
3. Should Display Solution (#3) or Broker API (#4) be pulled ahead of gRPC for your goals?
4. Add a scheduled (nightly) live-sandbox CI job, or keep live tests manual?
