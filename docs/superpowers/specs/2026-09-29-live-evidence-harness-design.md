# Live evidence harness — design

**Date:** 2026-09-29
**Status:** approved in conversation; awaiting written-spec review
**Supersedes:** nothing. First proposal for empirical verification of the SDK.

## 1. Problem

The conformance harness checks 209 documented endpoints against fixtures generated
from Webull's own OpenAPI JSON. It can therefore only ever confirm that the SDK agrees
with the documentation. It cannot confirm the documentation agrees with a real
server, and no test in the repository does.

That gap now carries live risk. Three consecutive releases changed response contracts
based on documentation alone:

| Release | Change | Rows closed |
|---|---|---|
| v2.1.33 | 62 required names across 20 types; `GetFDAssetsDetail` envelope | 157 → 81 |
| v2.1.34 | 76 declared names across 17 types | 81 → 25 |
| v2.1.35 | 14 response envelopes; `types.Page[T]` | 81 → 25 |

Not one of those 204 decisions has been checked against a running server. v2.1.35 is
the sharpest case: the old code provably fails against the documented shape
(`cannot unmarshal object into Go value of type []T`), but the new code has never been
shown to succeed against a real response. If Webull actually returns a bare array for
`GetMarketSectors`, that release shipped a regression and nothing in the tree detects it.

The remaining 25 baseline rows are almost entirely *container-kind* disagreements
where the page and the SDK disagree about the top-level shape. The harness has no basis
to prefer one side, so it records them rather than resolving them.

## 2. What the investigation found

Recorded because each one changed the design.

**HK credentials can settle 1 of the 10 remaining shape rows.** Nine are Display
Solution (403 at the host, needs a paid subscription) or Broker FD (404 in the HK
sandbox). The probe's value is therefore *not* draining the baseline. It is verifying
the surface the last three releases touched.

**Reachability cannot be determined statically.** `/market-data/screeners/*` returns
404 in the HK sandbox while `/market-data/stocks/*` works. A path-prefix heuristic over
all 193 endpoints classified only ~20 confidently and left 131 unclassified. Only
calling settles it, which makes the census the first deliverable rather than a
preliminary.

**127 of 193 endpoints require a query parameter and 50 require a request body.** The
census cannot send empty requests; minimum valid requests must be synthesised from each
page's parameter schema. Of the required query parameters, 82 carry an `example` or an
`enum`; the remaining 144 include 34 that depend on an `account_id` the SDK cannot
know without asking the server, which is why phase 0 exists. See §5.5.

**`Money.UnmarshalJSON` and `QuoteTime` accept both a JSON string and a number.** A
redaction that replaced every numeric value with a string placeholder would have
destroyed the single fact most worth verifying — whether Webull sends `close` as
`385.6` or `"385.6"`. This is why the artifact is a type skeleton rather than a
redacted body.

**A type skeleton is already a valid fixture input.** `jsonKind` (`shapes.go:361`)
classifies decoded values and `leafSamples` (`shapes.go:1068`) already reasons with
kind-level literals, so the existing checks consume a skeleton unchanged. No new
comparison code is required.

**`CompareBody(f Fixture, symbol string, t reflect.Type, body []byte)` takes raw
bytes**, so the live path reuses the same five checks with no parallel implementation.

**Unexported response types are not a blocker for capture.** `conformance/envelopes.go`
exists because six envelopes are unexported, but a probe that *calls* a method never
names its return type. Unexported types matter only for the reflect comparison, which
the symbol table already bridges.

## 3. Goals

1. Establish which of the 209 documented endpoints answer on a given credential, and
   record that as a durable artifact rather than prose.
2. Compare every reachable live response against both the SDK type and the documented
   fixture, and report the two directions separately.
3. Keep the result honest when it is negative — including "nothing here is verifiable
   with this credential".

## 4. Non-goals

- **Not** a fix for any of the 25 remaining baseline rows. Credentials cannot settle
  them.
- **Not** a general-purpose API fuzzing or contract-testing tool.
- **Not** a replacement for the documentation-derived harness, which remains the
  primary check and covers 159 endpoints HK cannot reach.
- **Not** a mechanism for endpoints it cannot reach. If phase 4 reports zero SDK
  disagreements, the correct outcome is a smaller artifact, not a larger one.

## 5. Design

Three new artifacts. No existing file changes.

### 5.1 Type skeleton — the committed evidence

For each reachable endpoint, the live response reduced to names, JSON kinds and
nesting. No values at all, not placeholders.

```
live      {"data":[{"symbol":"0700.HK","close":385.6,"outstanding":null}]}
skeleton  {"data":[{"symbol":"string","close":"number","outstanding":"null"}]}
```

`string` and `number` are preserved as distinct kinds, which is the fact `Money` and
`QuoteTime` currently mask. Written to `conformance/testdata/live/<fixture-id>.json`,
a sibling of the committed documentation fixtures.

### 5.2 Live divergences — a third set, kept separate

`conformance/live-divergences.json`, deliberately **not** merged into
`known-divergences.json`.

A live finding is a claim about *Webull*: "its server disagrees with its own OpenAPI
JSON". Every existing row is a claim about *the SDK*. Merging them would let a reader
take 25 documented rows plus N live rows as a single backlog of 25+N SDK defects,
which is a different and much weaker claim than either part. The two sets stay
separate, exactly as `missing-required-name` and `missing-declared-name` are kept
separate today because they rest on different evidence.

Per endpoint the harness runs the existing `Compare` twice:

- **skeleton vs SDK type** → what the SDK cannot read
- **skeleton vs documented fixture** → what the server disagrees with the docs about

### 5.3 Live manifest — provenance

`conformance/live-manifest.json` records, per endpoint: the probe timestamp, the base
host, the HTTP verb and path actually sent, the HTTP status received, and whether the
raw body decoded into the SDK type cleanly. This is what makes a later failure
attributable — a rotated credential surfaces as a 401 at the top of the report rather
than as a conformance finding.

### 5.4 Credentials

Read from the environment, matching the existing `WEBULL_SANDBOX` / `WEBULL_APP_KEY` /
`WEBULL_APP_SECRET` gates.

Never committed. `AGENTS.md` forbids inlining app keys, secrets and account IDs in
hand-written files, and the reason applies with extra force here: Webull can rotate
the published test credentials at any time, and a probe whose credential silently went
stale would report missing fields as conformance findings. For a tool whose entire job
is deciding whether the SDK matches reality, a stale shared credential is a
false-positive generator. Environment variables make that failure loud and immediate,
at the auth step, where it cannot be mistaken for a regression.

### 5.5 Parameter resolution — the one place phase 1 could stall

"Send a minimum valid request" is only meaningful if the value is derivable. Of the
required query parameters across the 193 endpoints, 82 carry an `example` or an
`enum` and can be taken from the page. The remaining 144 fall into four classes, and
they need different treatment:

| Class | Parameters | Resolution |
|---|---|---|
| Caller-supplied correlation key | `client_request_id` (10), `client_order_id` (3) | Synthesise a unique value; the server accepts anything |
| Declared scalar | `count` (4), `depth` (3), `real_time_required` (4), `overnight_required` (2) | A value from the declared type, respecting `minimum`/`maximum` where published |
| Caller-identity, derivable at runtime | `account_id` (34) | **Discovered, not synthesised** — see below |
| Domain constant | `symbol` (42), `symbols` (13), `series_symbol` (3) | The sandbox carries AAPL only, so `AAPL` where the value is a ticker |

`account_id` is the blocking detail: no document records it, because it is issued by
the server per credential. 34 endpoints depend on it. Phase 1 therefore begins with a
**pre-flight account discovery** — call the accounts endpoint once, take the first
`account_id` from the response, and thread it through. A census that started by
guessing would report 34 endpoints as unreachable for a reason that is not the
environment's fault, and that is precisely the kind of false negative this harness
exists to avoid.

If discovery itself fails, those 34 endpoints are recorded as
`blocked: account discovery failed`, which is a different and more useful statement
than "404".

## 6. Phases

**Phase 0 — Account discovery.** One call to the accounts endpoint; record the
`account_id` it returns. Cheap, and it unblocks 34 endpoints in phase 1.

**Phase 1 — Reachability census.** Synthesise a minimum valid request per endpoint
using the resolution order above; call it; record the outcome. Capture nothing.
Produces the reachability matrix, which ends the path-literal guessing in `AGENTS.md`.

**Phase 2 — Skeleton capture.** For every endpoint that answered, walk the raw body
into a skeleton and commit it with its live manifest entry.

**Phase 3 — Live comparison.** `conformance/live.go` loads the skeletons and runs the
existing checks in both directions, writing `live-divergences.json`.

**Phase 4 — Report.** One summary: *of N endpoints probed, X disagree with their own
documentation and Y disagree with the SDK.*

### 6.1 The gate on phase 4

Build phases 0, 1 and 4 first. If phase 4 reports **Y = 0**, phases 2 and 3 were a
mechanism for a problem that turned out not to exist, and the honest outcome is the
census plus a small drift check — not the full machinery. The order is deliberate: it costs one
probe run to learn whether the remaining work is worth doing.

## 7. Bounds

Recorded so they are design constraints rather than mid-project surprises.

- **The 14 envelope methods from v2.1.35 are unreachable.** `GetMarketSectors`,
  `GetDisplayGainersLosers`, `GetDisplayTopActive`, `GetMarketSectorDetail`,
  `GetFundDividends`, `GetEventContractMarkets`, `GetEventContractSeries` and the
  seven `brokerfd` methods sit on screener or Broker FD paths that HK returns 404
  for. The highest-risk recent change is the one this harness cannot verify.
- **HK-only caps the probe at roughly 50 of 209 endpoints.** Four single-item
  news/logos pages, all 40 `brokerfd` endpoints and all 29 `broker` endpoints are
  *permanently* unverifiable with this credential — not deferred, because no
  obtainable credential fixes them.
- **Sandbox market data is AAPL-only.** Some endpoints may answer 200 with sparse or
  empty bodies. A 200 with an empty array is recorded as an answer, not a failure.
- **The census may return mostly 401/404.** The finding is then "the HK sandbox cannot
  verify this SDK", which is a real result and is written up as one.
- **Skeletons are value-free but the manifest is disclosure-adjacent** — it records
  which host and when. Both are reviewed before commit.

## 8. Risks

| Risk | Mitigation |
|---|---|
| Skeletons lose the string/number distinction | Skeleton kind is read from the raw decoded value, never from the SDK type, so SDK permissiveness cannot mask it |
| A redacted projection would be circular | Explicitly rejected: the skeleton is built from the raw body, so both comparisons stay honest |
| Stale credentials read as conformance findings | Environment-only credentials; the manifest records HTTP status so auth failure is visible at the top of the report |
| Reachability assumptions baked in | The census *is* the measurement; no endpoint is assumed reachable |
| Mechanism built for a non-problem | Phase 4's gate, with Y = 0 as an explicit stopping condition |
| Live findings merged into the SDK backlog | Separate file and separate kind, with the reasoning recorded in both the spec and the report |

## 9. Verification

How this work is itself checked:

- The census is deterministic for a fixed host and credential, so a second run must
  produce the same reachability matrix; a difference means the environment moved.
- A skeleton must decode cleanly into its SDK type. A skeleton that does not is a bug
  in the capture step, not a finding, and must be reported as such.
- `make conformance-gate` and the full local gate stay green: the live path adds files
  and a new test binary, and must not alter the existing baseline.
- No live body, in any form, reaches a commit. A review check greps the diff for
  high-entropy tokens before commit.

## 10. Open questions

None outstanding. Resolved during brainstorming:

- *Mechanism or report?* Mechanism (capture-and-compare), gated on the phase 4 result.
- *Redact values or commit raw?* Type skeleton — stronger than either, since it
  contains no value.
- *Third divergence kind or separate report?* Separate file and separate kind.
- *Where do credentials come from?* Environment, never committed.
- *Can the harness settle the remaining 25 rows?* No, and that is a stated bound.
