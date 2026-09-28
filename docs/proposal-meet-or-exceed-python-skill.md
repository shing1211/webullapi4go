# Proposal: Positioning This Go SDK Against the Official Python Variant

**Status**: Proposal / planning document
**Repository tag at time of writing**: `v2.1.19`
**Evidence base**: repository-local documentation only — `docs/implementation-status.md`, `docs/reconciliation.md`, `README.md`, and the `docs/*.md` reference set. No external repository was cloned and no external content was executed.

## Scope

**In scope**: a gap analysis comparing this SDK's observable surface against the officially published Python variant, followed by an ordered remediation plan. Every claim below cites a file in this repository so a reader can verify it independently.

**Out of scope**: this document does not propose changing any `.go` file. It records what is wrong, what is already strong, and in what order the corrections should be made. Each remediation step is a proposal, not a completed change.

**Not a claim of superiority**: this document compares verifiable capabilities. Where this SDK is ahead, that is stated. Where it is behind or defective, that is stated too.

---

## 1. Current state

### 1.1 Module identity

```
Module        github.com/shing1211/webullapi4go
Import path   v1 (retained by decision; no /v2 migration planned)
License       Apache-2.0
Go            1.26+, no cgo
```

The tagged tree is not an installable published Go-semver v2 module. Consumers needing current work pin a commit; consumers wanting the released line use the `v1.x` proxy version. See `README.md` § Install.

### 1.2 Regions

The SDK supports nine regions: US, HK, JP, SG, MY, UK, MX, BR, ZA — the same set the Python variant advertises. Region selection is via `client.WithRegion(...)` and the `WEBULL_REGION` environment variable (`README.md` § Configuration).

### 1.3 Endpoint coverage

Per the generated reconciliation snapshot (`docs/reconciliation.md`, authoritative):

| Measure | Count |
|---|---:|
| Implemented endpoints | 209 |
| Documented-only gaps | 0 |
| Exact OpenAPI path match | 184 |
| Summary-only match | 4 |
| Path differs from both sources | 1 |
| Unresolved SDK path | 0 |
| No OpenAPI schema on page (gRPC) | 3 |
| No SDK symbol (manifest entry unmapped) | 17 |
| Intentionally not implemented | 0 |

The states partition the 209 rows; no endpoint is missing. A `✅ match` row is a **path-string comparison only** — it is not a correctness verdict. The reconciler does not observe HTTP verbs, request bodies, transport-host routing, or response schemas. Section 3 records defects that sit behind clean `✅` rows for exactly this reason.

### 1.4 Verification state

Implemented and offline-tested is not live-verified. `docs/implementation-status.md` records, per area, which surfaces were exercised against a live environment and which were blocked by credentials, scope, entitlement, or host.

Blocked at time of writing:

- Broker API HK — scope-protected, `401 ROUTE_NOT_PERMITTED`
- Broker FD US — no US sandbox credentials
- Display Solution — host returns `403` without entitlement
- Connect OAuth — US live access unavailable
- Options multi-leg — non-`SINGLE` strategies rejected with `417`

---

## 2. Where this SDK is already strong

Verified against repository documentation; the Python variant's published README does not document equivalents.

- **Structured error contracts.** `docs/errors.md` defines provider-refusal, replay-safe, and timeout-shaped categories with explicit handling rules.
- **Module separation.** `broker/` is a separate Go module with its own `go.mod` replacing the root. This permits independent release cadence and version isolation. The Python variant is a single installed package.
- **Test-selection gate.** `docs/sandbox.md` describes a `Sandbox` test selector that separates *test selection* from *environment configuration*, keeping credentials out of CI configuration.
- **Endpoint reconciliation.** `docs/reconciliation.md` is a generated, per-endpoint comparison against the official API surface. No equivalent document is published for the Python variant.
- **Verification transparency.** `docs/implementation-status.md` records what was live-verified, what was not, and why — including defects that fail silently.
- **Import-path stability.** The v1 path is retained by decision, so existing consumers are not broken by a major-version bump.

These should be preserved by any remediation work.

---

## 3. Known defects

All items below are recorded in `docs/implementation-status.md` and were re-confirmed against the current tree. Each is a genuine divergence, not a documentation nit. None is fixed by this proposal.

### 3.1 `brokerfd` sends every request to the core host

`brokerfd/client.go` routes its shared request helper through the core client's `Do` method, so all Broker FD traffic is sent to the trading/market-data host instead of the documented Broker FD host. The HK `broker` package routes correctly through `DoBroker` by comparison.

- **Impact**: every method in the `brokerfd` package.
- **Minimal fix**: route the shared helper through the Broker host.
- **Unblock**: US sandbox credentials. The HK host does not serve the FD surface, so HK alone cannot distinguish "wrong host" from "wrong region".

### 3.2 `brokerfd` uses undocumented `/broker-fd/*` paths

Fourteen non-test path literals remain under a `/broker-fd/` prefix while cached `broker-fd-api` reference pages use the `/broker/...` namespace. Of the fourteen, four align mechanically to documented pages, one is a probable duplicate symbol pair, six are plausibly ambiguous, and three have no documented counterpart at all.

- **Impact**: the affected `brokerfd` endpoints. Failures are loud (wrong path) rather than silent.
- **Minimal fix**: align the four unambiguous literals, reconcile the duplicate symbol pair, and investigate the other nine. The three with no counterpart must not be guessed.
- **Unblock**: US sandbox credentials, plus a written answer from the provider on whether those three endpoints exist at all. A `404` cannot distinguish "undocumented path" from "endpoint not offered".

### 3.3 `broker.UpdateVirtualAccount` sends the wrong verb and body shape

`broker/accounts.go` builds a query-string path and issues `PUT`. The documented endpoint is `POST` and requires `account_id` and `client_request_id` in a JSON body, with no `account_name` field. The request struct carries an undocumented `AccountName` field. `GetVirtualAccount` does document `account_id` as a query parameter, which is the likely origin of the copy.

The path is correct, which is why the reconciler reports a clean match — it compares paths, not verbs. The corresponding test currently certifies the incorrect contract.

- **Impact**: `broker.UpdateVirtualAccount`.
- **Minimal fix**: issue POST with the documented body fields, remove `AccountName`, and update the test in the same change.
- **Unblock**: a production or US-scoped Broker credential.

### 3.4 `data.GetDisplaySnapshot` differs from both official sources

The path is `/openapi/market-data/stock/snapshot` and the verb is GET. Both official sources indicate `POST /market-data/stocks/snapshots/list`. This is the last remaining `/openapi/…` holdout in a const block whose siblings were already aligned.

The reference page is self-contradictory (`operationId` indicates GET against `method: post`), so the verb is genuinely ambiguous. The request-shape mismatch is not: the documented request is a JSON body carrying an **array** of `{category, symbols}` objects, while the SDK models one category with one flat symbol list encoded as query parameters. Correcting this changes the public parameter type, which is shared by two public methods.

- **Impact**: `data.GetDisplaySnapshot`; the fix is a breaking public API change.
- **Minimal fix**: not a patch. Requires a maintainer decision on how to version a breaking change on a module that remains on the v1 import path — an added field with a deprecation window, a new method alongside the retained one, or a `/v2` migration.
- **Unblock**: a Display Solution entitlement. Probing alone cannot settle it, because the provider documentation contradicts itself and a second cached page documents the opposite contract.

### 3.5 `brokerfd.GetFDPositions` response DTO cannot receive required properties

The response struct cannot receive three of the eight properties the endpoint requires; the three values are silently zeroed. This is the sharpest illustration that a `✅` reconciliation row is not evidence of correctness: correct path, clean match, and a response DTO that drops required data.

Recorded separately from 3.2 because it is the one defect on that surface that fails silently, where every path literal in 3.2 fails loudly.

- **Minimal fix**: align the DTO to the documented response envelope.
- **Unblock**: US sandbox credentials.

---

## 4. Ordered remediation plan

Proposed, not executed.

**Phase 1 — confirm before changing anything**

1. Re-verify each citation in section 3 against the current tree. Line numbers cited in `docs/implementation-status.md` shift when comments are inserted above them; re-check rather than assume.
2. Confirm the `docs/reconciliation.md` snapshot is current with respect to the repository tag.
3. Establish which blocked surfaces can be unblocked with credentials that are actually obtainable, and which need a written answer from the provider.

**Phase 2 — correct the defects, smallest first**

4. §3.3 `broker.UpdateVirtualAccount` — smallest change, no public-type impact.
5. §3.1 `brokerfd` host routing — one-line routing change, unlocks the whole package.
6. §3.2 and §3.5 `brokerfd` paths and response DTO — the four unambiguous path literals, the duplicate symbol pair, then the ambiguous nine. Requires credentials.
7. §3.4 `data.GetDisplaySnapshot` — deferred until the versioning decision is made, because it is a breaking public API change.

**Phase 3 — surfaces blocked on access**

8. Once §3.2, §3.5 and the blocked areas can be exercised, re-run the response-contract measurement described in `docs/implementation-status.md` and update the implementation status document with the results.

**Not proposed for change**: the structural advantages in section 2.

---

## 5. Feature-level comparison against the Python variant

The Python variant's published feature list includes multi-leg option strategies, combo orders (OTO/OCO/OTOCO), option leg in/out, and algorithmic orders (TWAP/VWAP/POV), all scoped to US.

This SDK's feature matrix (`README.md`) does not list those. Before treating that as a gap, note that it may be an access boundary rather than a missing implementation: `docs/implementation-status.md` records multi-leg and event-contract live behavior as unavailable or rejected in HK, and non-`SINGLE` strategies as rejected with `417`.

**Recommendation**: establish whether those paths are reachable with obtainable credentials before implementing SDK support for them. Adding client-side support for a surface that cannot be exercised produces untested code.

---

## 6. Verification notes

Every claim above cites a file in this repository. A reader can confirm each one by reading the cited file or by running `git grep` against the cited symbol.

| Claim | Source |
|---|---|
| Module identity, install semantics, region configuration, feature matrix | `README.md` |
| Verification state per area; blocked surfaces; defects in section 3 | `docs/implementation-status.md` |
| Endpoint counts and their partition; what a match row does and does not mean | `docs/reconciliation.md` |
| Error contract categories | `docs/errors.md` |
| Sandbox test selection | `docs/sandbox.md` |
| Observability setup, metric names, span names, correlation propagation | `docs/observability.md` |
| Streaming lifecycle: reconnect, resubscribe, health, watchdog | `docs/streaming.md` |
| Trading surface, cancellation, guardrails | `docs/trading.md` |

Where this document describes the Python variant, it is characterised from its published README as a single installed package exposing a CLI, configured through a `.env` file, covering the same nine regions. No external repository was cloned and no external content was executed.
