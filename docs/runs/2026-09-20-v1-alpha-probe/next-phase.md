# Next Phase — v0.6: Credentialed Verification & Broker API

## What was completed

This run shipped the two largest missing Trading features: **multi-leg options
orders** (10 new `OptionStrategy` values, multi-leg validation, notional
handling) and **futures order validation** (dedicated per-market order-type
matrix, QTY-only/whole-contract/DAY-GTC rules), plus speculative **option
chain/expiration discovery** functions and full tests and docs. A quality
review closed five findings, chiefly adding provisional markers to unverified
rules and fixing numeric strike comparison.

## Gaps and deferred items

- **US sandbox is the single blocker** for everything below. HK sandbox returns
  404 for every unconfirmed endpoint; the option chain endpoint is not published
  at all.
- Multi-leg (`TODO(t8)`), futures (`TODO(t9)`), and option-chain (`TODO(t10)`)
  wire values are provisional and untested against the live API.
- No top-level `Market` validation for OPTION orders (review finding 5).
- `WithMaxOrderNotional` is bypassable by attaching a second leg (documented,
  accepted for now).
- Broker FD HTTP and gRPC remain unimplemented/guessed.
- HK/CN derivatives are not covered.

## Candidate next-phase items

| # | Title | Why now | Effort | Dependencies | Risks |
|---|-------|---------|--------|--------------|-------|
| 1 | **US sandbox account + credential verification** | Unblocks 8 tasks | S | Webull account signup | Approval wait |
| 2 | Verify T8/T9/T10 against US sandbox | Convert provisional to confirmed | M | #1 | Real rules may differ; schema drift |
| 3 | Fund + Crypto + Screener v2 confirmation | Close v0.5 stubs | M | #1 | Sandbox data limits |
| 4 | Broker FD HTTP path discovery | Unimplemented area | M | #1 | Paths may not exist |
| 5 | Broker FD gRPC streaming client | Unimplemented area | L | #1, proto capture | Proto schema unconfirmed |
| 6 | OPTION top-level market validation | Close review finding 5 | S | — | Low |
| 7 | Release v0.6 with confirmed schemas | Ship the above | S | #2–#5 | — |

## Recommended next phase

**Phase v0.6: Credentialed Verification.**

Prerequisite: a US sandbox account. Then:

1. **T1 — Provision US credentials**: sign up, export `WEBULL_APP_KEY` /
   `WEBULL_APP_SECRET`, confirm token acquisition against the US sandbox.
2. **T2 — Verify T8/T9/T10**: probe multi-leg option placement, futures order
   validation, and the option chain endpoint; replace `TODO(t8|t9|t10)` with
   confirmed values.
3. **T3 — Verify Fund/Crypto/Screener v2**: run `examples/probe` against US;
   finalize those stubs.
4. **T4 — Broker FD discovery**: find working HTTP paths; capture gRPC proto.
5. **T5 — Release v0.6.0**.

Draft task breakdown:

| ID | Objective | Role | Acceptance |
|----|-----------|------|-----------|
| N1 | Provision US sandbox credentials | ops | Token acquired against US sandbox |
| N2 | Confirm multi-leg option rules | backend | `TODO(t8)` removed; sandbox integration test passes |
| N3 | Confirm futures rules | backend | `TODO(t9)` removed; sandbox integration test passes |
| N4 | Confirm/remove option chain | backend | Either confirmed path or functions removed |
| N5 | Verify Fund/Crypto/Screener v2 | backend | Stubs finalized; probes recorded |
| N6 | Broker FD HTTP paths | backend | `brokerfd` functions hit live paths |
| N7 | Broker FD gRPC client | backend | Streaming client with tests |
| N8 | Release v0.6.0 | release | Tag pushed to GitHub + Gitee |

## Open questions for the human

1. When will the US sandbox account be available?
2. Keep the **speculative** option-chain functions until then, or remove them?
3. Is Broker FD (the full Broker API) in scope for v1.0, or post-v1.0?
4. Should HK/CN derivatives be part of v1.0 or deferred further?
