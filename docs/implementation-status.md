# Implementation Status

Last updated: 2026-09-29

- Latest repository tag: **`v2.1.35`** (2026-09-28)
- Current hardening: **last declared in repository `v2.1.4`**; introduced in `v2.1.1`.
  This is the release at which the hardening was last *re-declared*, not the newest
  tag carrying SDK code — those differ, because `v2.1.4` changed no `.go` file at
  all. No release since has re-declared it, so the field is current rather than
  stale. For the record: `v2.1.1` and `v2.1.3` carry the hardening code, `v2.1.5`
  is the last tag to change any SDK `.go` file, and nothing from `v2.1.6` onward has
  changed SDK code — only the conformance harness, fixtures and documents.
- Module path: **`github.com/shing1211/webullapi4go`**, kept on the v1 import path
  **by decision**; no `/v2` migration is planned
- Installable version: **`v1.1.1`** — the module proxy serves only the `v1.x`
  line, so `v2.x` tags are not installable with `go get`. Consume newer work by
  pinning a commit
- Offline evidence recorded 2026-09-26: module-aware race/vet checks passed;
  `make cover` measured 73.6% aggregate root coverage and 80.8% in nested
  `broker/` (measurements, not guarantees)
- Live verification: partial historical evidence only; the current hardening
  was not newly live-verified. The repository tag is not a published
  Go-semver v2 module.
- Live evidence recorded 2026-09-29: one authorised reachability walk of all 193
  documented endpoints against the HK sandbox, a value-free capture of the 55
  response shapes that answered HTTP 200, and a comparison of each against the SDK
  type and the documentation. **Of 55 probed, 2 disagree with their documentation
  and 12 disagree with the SDK.** It measures response contracts, not the
  repository-tagged hardening, and it is one host on one day — see item 28.

## How to read this page

| State | Meaning |
|---|---|
| Implemented | The surface exists in the current tree |
| Offline-tested | Credential-free unit tests or local fakes pass |
| Live-verified | Successfully exercised against an available Webull environment |
| Blocked | Credentials, scope, entitlement, host, or sandbox data prevent verification |
| Repository-tagged | Included in the authorized repository Git tag `v2.1.1`; this is not a published Go-semver v2 module |
| Released | Included in a published module release; the unchanged root path means `v2.x` Git tags are not published Go-semver v2 modules |

Implementation and offline tests do not imply live verification.

## Architecture status

| Layer | Canonical location | Status |
|---|---|---|
| Core client | `client/` | Preserved; current pipeline hardening is offline-tested |
| Public services | `data/`, `stream/`, `trade/`, `events/`, `connect/`, `display/`, `brokerfd/`, `brokerfd/events/` | Preserved at repository root |
| Broker HK | `broker/` | Preserved as a separate Go module |
| Shared foundations | `pkg/errors`, `pkg/observability`, `pkg/resilience`, `pkg/transport`, `pkg/types` | Implemented and offline-tested |
| Domain | `pkg/domain/money`, `pkg/domain/order` | `money.Money` and OMS reconciliation implemented and offline-tested |
| Convenience | `webull/` | Thin core-client aliases, not an aggregate facade |

The proposed full relocation of service packages under `pkg/`, root-service deprecated shims, and raw `decimal.Decimal` DTO migration are superseded and closed. The stream path has no `sync.Pool`.

## Service matrix

| Area | Implemented | Offline-tested | Live-verified | Blocked or unverified |
|---|---:|---:|---|---|
| Core `client/` | Yes | Yes | Core auth and selected HK HTTP paths | Current pipeline hardening is not newly live-verified |
| Authentication | Yes | Yes | Token lifecycle in HK sandbox | Production activation requires the Webull App 2FA flow |
| Market Data HTTP | Yes | Yes | Selected AAPL/fundamentals/watchlist calls | US-only and entitlement-gated calls remain partial |
| Market Data MQTT | Yes | Yes | Basic MQTT-over-WebSocket paths | Current channel/health/shutdown hardening is offline-tested only |
| Trading HTTP | Yes | Yes | Accounts, assets, previews, guarded order paths | Multi-leg/futures/event live behavior is unavailable or rejected in HK |
| Trading Events | Yes | Yes | Basic stream path was previously exercised | Current telemetry is offline-tested only |
| Connect OAuth | Yes | Yes | — | US live access unavailable |
| Display Solution | Yes | Yes | — | HK host returns `403` |
| Broker API HK | Yes | Yes | — | HK returns `401 ROUTE_NOT_PERMITTED` because scope is missing |
| Broker FD HTTP | Yes | Yes | — | US-only; HK returns `404`, no US credentials |
| Broker FD Events | Yes | Yes | — | US-only; no live verification available |
| Shared foundations | Yes | Yes | Not applicable | — |

## Endpoint reconciliation

The generated [SDK ↔ API Reconciliation](reconciliation.md) is authoritative. Its 2026-09-26 snapshot reports:

| Measure | Count |
|---|---:|
| Implemented endpoints | 209 |
| Documented-only gaps | 0 |
| ✅ Exact OpenAPI JSON path match | 184 |
| 🟡 Matches the docs summary only | 4 |
| ⚠️ Path differs from both sources | 1 |
| ❓ Unresolved SDK path | 0 |
| 📄 No OpenAPI schema on page (gRPC page, not a REST endpoint) | 3 |
| ➖ No SDK symbol (manifest entry unmapped) | 17 |
| ℹ️ Intentionally not implemented | 0 |

The states partition the 209 implemented endpoints: 184 + 4 + 1 = 189 rows carry a verified path, and the remaining 20 are the 3 rows labelled `📄 No OpenAPI schema on page` plus 17 manifest entries deliberately mapped to no SDK symbol. 189 + 3 + 17 = 209, so no endpoint is missing. **The 3 is a label count, not a page count**: 7 gRPC reference pages embed no OpenAPI schema, and the 4 that the manifest also maps to no SDK symbol are labelled `➖ No SDK symbol` instead, because `_reconcile_data()` in `tools/webull-docgen/docgen.py` evaluates the unmapped branch before the no-openapi branch and the status table must partition the 209 rows. A further 3 rows have a JSON block that yields no `path` and are likewise labelled unmapped, so 10 rows in total have no usable official path. Those two categories are why the previously quoted "25 unresolved" was wrong: 20 of those 25 entries were generator artifacts, and the remaining 5 were investigated individually — 4 were correct SDK code the generator could not statically follow, and 1 is the real path mismatch now reported as ⚠️.

**Resolution, 2026-09-28: all four summary-only rows are the `/openapi/*` namespace,
and in each the SDK path equals the official `llms.txt` summary and differs only from
the OpenAPI JSON.** `llms.txt` is one of the two official machine-readable sources,
not a third-party index, so the SDK sides with an official source and the two Webull
sources disagree with each other:

| SDK | OpenAPI JSON | llms.txt summary (what the SDK sends) |
|---|---|---|
| `client.CreateToken` | `POST /auth/tokens/create` | `POST /openapi/auth/token/create` |
| `client.CheckToken` | `POST /auth/tokens/check` | `POST /openapi/auth/token/check` |
| `data.GetStockInstruments` | `GET /trading/instruments/stocks/profiles/list` | `GET /openapi/instrument/stock/list` |
| `data.GetDisplaySnapshot` | `POST /market-data/stocks/snapshots/list` | `POST /openapi/market-data/stock/snapshot` |

**The step that proposed to align these to the OpenAPI JSON is withdrawn, because
taken literally it moves `client.CreateToken` and every authenticated call in this
SDK depends on that path.** The OpenAPI JSON spelling has never been exercised by
this code; the summary spelling is what the SDK sends today. Two official sources
disagreeing makes the SDK path *supported*, not *proved* - no live call in this
session tested either spelling, and the sandbox-gated suites that would exercise
`GetStockInstruments` were not run. Unresolved SDK paths remain `0`, and the single
row differing from *both* sources is the Broker FD assets summary, which is the
`/broker-fd/*` item below.

There are no documented-only endpoint gaps, but the snapshot is not a zero-discrepancy report: 4 summary-only and 1 differing remain. Generated pages under `webull-api/` are not hand-edited; the label correction that produced this partition lives in `tools/webull-docgen/`.

A `✅ match` row is a path comparison, not a correctness verdict. The reconciler compares path strings only: it reports `broker.UpdateVirtualAccount` as a clean `✅ match` although the request behind that path is defective, and it cannot observe HTTP verbs, request bodies, transport-host routing, or response schemas at all, so it is blind to the `brokerfd` host-routing and `data.GetDisplaySnapshot` defects recorded below as well, and to the `brokerfd.GetFDPositions` response-schema defect, which is its sharpest case: a correct path, a clean `✅ match`, and a response DTO that cannot receive three of the eight properties the endpoint requires, so three values are silently zeroed. That is a limit of what a path comparison can show rather than a defect in the generator, which is doing its stated job of comparing documented paths against SDK paths. A green row in `reconciliation.md` is not evidence that an endpoint is correct. Nothing in that table speaks to the response contract at all; the response contracts are measured separately, and largely diverge — see item 21.

## v2.1.1 repository-tagged hardening

| Work | State |
|---|---|
| Category matching and identity-specific semantic sentinels | Implemented, offline-tested |
| HTTP 417 compatibility mapping and MQTT/event terminal mappings | Implemented, offline-tested |
| Shared `Do` / `DoBroker` / `DoStream` request pipeline | Implemented, offline-tested |
| One-based attempt telemetry, stable correlation IDs, W3C propagation, and `SafeErrorText` | Implemented with real OTel/local-server assertions |
| Provider-order-independent resilience metrics and caller-owned custom breakers | Implemented, offline-tested |
| OMS snapshot reconciliation and failure-scene mapping | Implemented, offline-tested |
| Account-scoped tracking and stable auto-generated client IDs | Implemented, offline-tested |
| Compare-and-swap stream state and deterministic data-only health recovery | Implemented with deterministic watchdog tests |
| Idempotent MQTT/channel shutdown, blocked-dispatch cancellation, and registration-order head-of-line behavior | Implemented, offline-tested for all channel types |
| MQTT Connect cancellation, Connect/Close races, and late-callback suppression | Implemented, offline-tested |
| Serialized resubscription and active-subscription replay | Implemented, offline-tested |
| Trading and Broker FD event telemetry, including failure/cancellation outcomes | Implemented, local-gRPC tested |
| Trading `Close` cancellation of all active `Run` calls | Implemented, offline-tested |
| Cancellation, strict data/event leak checks, and public resilience/transport/type/`webull` coverage | Implemented, offline evidence recorded |
| Account monitor and updated order/stream examples | Implemented, included in offline module build/test |
| Documentation reconciliation | Implemented; final F1 report records strict docs/link/diff gates |

These items are included in repository tag `v2.1.1` (2026-09-25). They are
not a published Go-semver v2 module, and the current hardening was not newly
live-verified.

## Milestones and historical tags

The `v2.x` rows below are repository-tag records, not published Go-semver v2
modules; `v2.1.17` is the current authorized repository tag and earlier rows are
historical. The root module path remains `github.com/shing1211/webullapi4go` and
stays on the v1 import path by decision.

| Version | Scope |
|---|---|
| v1.1.0 | Full documented endpoint coverage; 209 implemented, 0 documented-only gaps |
| v1.1.1 | DevOps, security checks, multi-OS CI, leak tests, and fuzzing |
| v2.0.0 | Public error/transport/resilience/domain foundations and thin `webull` aliases; root services retained |
| v2.0.1–v2.0.2 | `money.Money` DTO conversion in `data/`, `trade/`, and `brokerfd/` |
| v2.0.3–v2.0.4 | MQTT cleanup, clock correction, idempotency helpers, transport tuning, and resilience preset |
| v2.0.5–v2.0.7 | Interceptors/hooks, initial OMS, stream state/channels, slog, OTel tracing, and metrics |
| v2.0.8–v2.0.9 | Context hygiene and structured public errors |
| v2.1.0 | `go vet` mutex-copy fixes |
| v2.1.1 | Current hardening, including error contracts, request/OMS reconciliation, stream/MQTT lifecycle, event telemetry, cancellation/leak coverage, and documentation reconciliation; repository tag, not a published v2 module |
| v2.1.2 | OpenTelemetry v1.46.0 across all modules, GitHub Actions bumps, pinned `govulncheck`, Dependabot grouping, and documentation policy updates; repository tag, not a published v2 module |
| v2.1.3 | Honest nightly live-verification signal, removal of the dead `internal/` shims and duplicate resilience tree, and direct tests for the order reconciliation, MQTT, money, error, and region surfaces; repository tag, not a published v2 module |
| v2.1.4 | Nested module coverage gate for `broker`, per-module `govulncheck`, a strict MkDocs gate, and run-index and roadmap corrections; repository tag, not a published v2 module |
| v2.1.5 | Additive Broker FD event metadata delivery, the first stream dispatch benchmarks with recorded head-of-line measurements, and classification of every non-exact reconciliation state; repository tag, not a published v2 module; tagged with a red CI run |
| v2.1.6 | Correction of the nested coverage CI job: artifact-safe matrix labels, tolerant artifact upload, and a guarded coverage report step; repository tag, not a published v2 module; first tag with fully green CI |

Earlier v0.x and v1.0 milestones remain recorded in the root `IMPLEMENTATION_STATUS.md`.

## Verification performed

- Module-aware offline race/vet checks were recorded as passing on 2026-09-25.
  `make test` and `make test-race` traverse the root, `broker/`, and nested
  example modules; root-only `go test ./...` does not.
- `make cover` measured 73.6% aggregate root coverage and 80.8% in `broker/` on
  2026-09-26. These dated measurements are not correctness guarantees and are not
  combined into one repository-wide percentage.
- Unit tests are offline and credential-free. Live tests use the `Sandbox`
  selector and remain environment-gated.
- Data and both event packages use strict goroutine-leak checks; event tests
  cancel and wait for completed runs.
- CI's 60% coverage gate is root-only. Nested coverage and
  `mkdocs build --strict` remain Makefile/release gates rather than CI gates.
- No new live validation was performed for the v2.1.1 repository-tagged hardening.

## Live-verification boundaries

Previously exercised HK paths include core tokens, selected AAPL data and fundamentals, watchlists, accounts, balances, positions, order previews, MQTT streaming, and Trading Events. That does not cover all 209 endpoints or the current hardening.

The 2026-09-29 run is the first to walk the documented surface systematically rather
than by hand, and it changes two of those statements rather than one: it reached
**55** of the 193 documented endpoints, so "selected" is now a measured 55 rather
than a list of examples; and 6 of the 18 `display-solution` endpoints answered HTTP
200 from the core sandbox host, which is new, is not a statement about the
`display/` package, and does not bear on the Display Solution host named below.

Neither this list nor the live-blocked defects below is the complete set of known
divergences. Both are request-side or single-DTO findings. The response contracts
across the whole API surface have now been measured, and they are largely
divergent: see item 21 and its subsection. A further limit sits underneath that
measurement, and it is a limit of the evidence rather than of the SDK: 90 of the
154 compared rows come from pages that publish no `required` list. A weaker name
check now runs over the property names those pages do declare, and it is not a
formality - it produced 104 further rows, 9 of them on rows that had recorded no
divergence at all - but 6 of the 90 remain unnameable and the evidence on the rest
is a description rather than a promise. That is item 22. A **third** set now exists
and is also not in this list: 13 findings measured against a live server rather
than against a page, recorded separately because a live finding is a claim about
Webull's server rather than about this SDK. That is item 28.

Blocked or unverified areas:

- Display Solution: host-level `403` on the Display Solution host. Measured on
  2026-09-29: 6 of the 18 `display-solution` endpoints answered HTTP 200 from the
  **core** sandbox host `api.sandbox.webull.hk`, so the `403` is a property of that
  host and not of every Display Solution endpoint; whether the core host's answers
  are the same contract is a question for Webull.
- US-only crypto, fund, screener, and Broker FD: no US sandbox credentials.
- Broker API HK: missing scope (`401 ROUTE_NOT_PERMITTED`).
- Footprint: missing entitlement (`403`).
- Options/multi-leg: limited sandbox contracts and non-`SINGLE` strategies rejected with `417`. The SDK's `INVALID_TOKEN` category is a compatibility mapping for every HTTP 417 and does not prove a token failure.
- Futures and event-contract trading: validation is offline-tested; live product behavior is unverified.
- SSE news: upstream `504`.
- Generated paths: four summary-only matches, one path differing from both, and
  `0` unresolved SDK paths. The previously quoted 25 unresolved paths were 20
  generator artifacts plus 5 individually investigated entries.

### Live-blocked SDK defects

Found by static analysis on 2026-09-26, with the `brokerfd.GetFDPositions`
response-schema defect added on 2026-09-27, also statically. None is
live-verified — no pass touched the network — and each may only be changed once
the credential or entitlement it names is available, because the fix cannot be
confirmed without it.

- **RESOLVED: `brokerfd` now sends every request to the Broker host.** The shared `do` helper calls `c.core.DoBroker` (`brokerfd/client.go:53`), which is what the sibling `broker/` package has always done at `broker/client.go:63`.
  The defect was a fact about the code rather than about the server: the transport,
  the `client.Endpoints.BrokerHTTP` field and the sibling module's use of it all
  already existed, and the two packages simply disagreed. The recorded unblock - a US
  credential to tell "wrong host" from "wrong region" - was needed to observe the
  symptom, not to establish the cause, which is why this is closed without one.
  `client.EndpointsFor` populates `BrokerHTTP` for every region and environment, so
  the change is correct by default; `TestEndpointsForAlwaysSuppliesBrokerHost` asserts
  that table rather than assuming it, because an empty field would make the package
  unusable rather than merely wrong.

  The one behaviour change is that a caller building `client.Endpoints` by hand and
  omitting `BrokerHTTP` now gets a typed configuration error naming the missing field,
  where the old code sent the request to the core host and returned its 404.

  **The fix exposed a live-traffic hazard.** `client.WithBaseURL` overrides one field
  and sets the internal override flag, so `client.New` leaves every other endpoint at
  its default - the *production* Broker host. 54 call sites in `brokerfd`'s tests used
  it, so correcting the routing turned the whole suite into calls against a production
  API. All 54 now use `client.WithEndpoints` with both hosts, as `broker/`'s tests
  already did, and `client/client_test.go` asserts the `WithBaseURL` behaviour so it
  cannot change silently.

  Verification: `brokerfd/host_routing_test.go` runs two test servers, one per host,
  and asserts which is reached. Reverting the one line fails three of its tests.

- **`brokerfd` uses undocumented `/broker-fd/*` paths.** 14 non-test literals
  remain (`brokerfd/accounts.go:46-48`, `brokerfd/assets.go:25`,
  `brokerfd/brokerfd.go:58,68`, `brokerfd/documents.go:23-24`,
  `brokerfd/funding.go:26,30`, `brokerfd/instruments.go:27,29,31`,
  `brokerfd/journals.go:25`) while every other cached `broker-fd-api` page uses
  the `/broker/...` namespace. `brokerfd/assets.go:25` sends
  `/broker-fd/assets/summary` against a documented
  `GET /broker/assets/summaries/get` — the single ⚠️ row in the snapshot.

  **All 14 are now classified against the cache, and none is alignable.** The recorded
  judgement that "4 can be aligned" was optimistic; against the cached pages the answer
  is 0. The cache holds 50 Broker FD pages and every one is under `/broker/...`, so none
  of the 14 `/broker-fd/*` literals appears in the documentation at all.

  - **Contested, 6** - two SDK symbols plausibly own the one documented page, so it
    cannot be assigned to either without a decision the docs do not make:
    `GetAccountsSummary` / `GetFDAssetsSummary`, `GetFDPositions` / `GetPositions`
    (where `GetFDPositions` already sends the exact documented path), and
    `GetFDCorporateActions` / `GetFDCorporateActionDetail`,
    `GetFDCashJournalDetail` / `ListFDCashJournals`, and `CreateFDAccount` /
    `SubmitAccountForm`.
  - **No documented endpoint, 6** - and therefore unresolvable by any credential,
    because no probe can return a path no page describes: `ListDocuments` and
    `GetDocumentDetail` (only upload and download are documented),
    `GetFDAchAccountDetail` and `GetFDBankAccountDetail` (only create, delete, list),
    `GetFDStockLocate` (no page in any namespace), and `GetFDECInstrumentDetail` (only
    category, series, event and market lists).
  - **Ambiguous, 1** - `GetAccountFormStatus`, whose two candidates are
    `GET /broker/accounts/applications/get` and `GET /broker/forms/versions/list`,
    neither named for a status.
  - **The near miss, recorded so it is not re-attempted from a name match.**
    `GetAccountFormDetail` looked like the one alignable case: the cached Form Content
    page declares `GET /broker/forms/get` and is the only single-form fetch. It was
    aligned, the fixture regenerated, and the change reverted. That endpoint
    "retrieves the JSON schema for the specified form code and version": it takes
    `form_code` and `version` as query parameters and its 200 body is a **JSON Schema
    fragment** whose properties are the schema keywords `required`, `type`, `format`,
    `min_items`, `max_items`, `max_length`, `enum_values`, `description` and `example`.
    The method fetches a form *instance* by `form_id` and decodes `brokerfd.AccountForm`.
    The harness agreed, reporting `form_code` and `version` as required response names
    that no `AccountForm` field carries - a true statement about the page and a false
    one about the method. Matching on the word "form" paired different *kinds* of
    endpoint, and the same trap applies to every candidate.

  **What this changes about the block.** Six of the 14 are not "blocked on a credential"
  but "unanswerable from the documentation", a different category that should not be
  waited on. The other 8 need the live probe already named, and the contested 6
  additionally need a decision about which symbol owns a shared page. The reasoning is
  recorded in the code next to the three form path constants, and the docgen manifest's
  Form Content row names the near miss.
  `brokerfd.GetPositions` (`brokerfd/brokerfd.go:67-68`) maps to no documented
  page, and the assets-summary manifest entry at
  `tools/webull-docgen/_common.py:318` names two symbols for that page:
  `brokerfd.GetAccountsSummary` (`brokerfd/brokerfd.go:57-58`) has no page of
  its own and `brokerfd.GetFDAssetsSummary`'s DTO (`brokerfd/assets.go:32-41`)
  does not match the documented `balance`/`positions` envelope. A separate
  defect on the same Broker FD surface is the `brokerfd.GetFDPositions`
  response-schema mismatch below; it is recorded separately rather than merged
  here because it is the one defect on this surface that fails silently, where
  every literal in this bullet fails loudly. Impact: the affected `brokerfd`
  endpoints. The size is larger than the literal count suggests, because only a
  minority of the 14 have an unambiguous documented counterpart. 4 align
  mechanically — `brokerfd/assets.go:25` → `/broker/assets/summaries/get`,
  `brokerfd/brokerfd.go:68` → `/broker/assets/positions/list`,
  `brokerfd/instruments.go:29` →
  `/broker/instruments/stocks/corporate-actions/get`, and
  `brokerfd/journals.go:25` → `/broker/journals/cash-journals/get` — while 1
  (`brokerfd/brokerfd.go:58`) is a probable duplicate, 6 are plausibly ambiguous
  (`brokerfd/accounts.go:46-48`, `brokerfd/funding.go:26,30`,
  `brokerfd/instruments.go:31`, where a documented page exists in the same area
  but not for the same operation), and 3 have no documented counterpart at all
  (`brokerfd/instruments.go:27` stock-locate, `brokerfd/documents.go:23`
  documents, `brokerfd/documents.go:24` documents/detail; the cache holds only
  `/broker/documents/download` and `/broker/documents/upload` in the documents
  namespace, and no stock-locate page). Those buckets are a plausibility
  assessment summing to the 14 literals, not settled mappings. Minimal fix:
  align the 4 unambiguous literals, reconcile the duplicate symbol pair, and
  investigate the other 9; do not guess the 3 that have no documented
  counterpart. Unblock: US sandbox credentials; probe
  `/broker/assets/summaries/get` and `/broker-fd/assets/summary` side by side.
  Credentials alone cannot settle those 3, because a US `404` does not
  distinguish an undocumented path from an endpoint Webull does not offer, so
  they also need a written answer from Webull about whether the endpoints exist
  at all.
- **RESOLVED: `broker.UpdateVirtualAccount` now sends the documented verb, body and fields.** It issues POST with a JSON body (`broker/accounts.go:157`) where it issued PUT with `account_id` in the query, and `UpdateVirtualAccountRequest` now carries the documented fields (`broker/accounts.go:67`).
  **The recorded defect was larger than it read.** The status recorded the wrong verb
  and body shape. The official page says more: the verb is POST, `account_id` is a
  **body** field, the request takes no query parameters other than the auth headers,
  and two body fields are required - `account_id` and `client_request_id`. The old
  request type had exactly one field, `AccountName`, which the page does not declare
  at all. So it was wrong in every respect at once.

  **Source-level break, taken deliberately.** `AccountName` is gone rather than kept
  and ignored, because a field that was never in the contract and is silently dropped
  is the worst of the three options: the caller compiles, sets it, and learns nothing.
  `account_id` stays a method argument so the required pair cannot be built apart. The
  unexported `put` helper had this as its only call site and was removed with it.

  Unblock, for the live confirmation only: a production or US-scoped Broker credential.
  The fix is derived from the official page rather than from an observation, and **has
  not been live-verified**. It was invisible for so long because the reconciler
  compares *paths* and the path was correct; see the verb check below.

- **`data.GetDisplaySnapshot`: the recorded premise was wrong, and correcting it weakens the case for the breaking change.**
  `data/display_quotes.go:33` sets
  `pathDSSnapshot = "/openapi/market-data/stock/snapshot"` and `:63` sends it
  with GET, while the OpenAPI JSON says

  **Correction, 2026-09-28: the recorded premise was wrong.** This said the method
  differs from *both* official sources. It does not. The `llms.txt` summary for this
  endpoint is `/openapi/market-data/stock/snapshot`, which is exactly the path the SDK
  sends, so the method matches one official source and differs from the other. The error
  was inherited from the same miscount as the four summary-only rows above.

  **What that changes.** The path was the tidiest-looking part of a large breaking
  change and is now the least likely to be part of the fix, because aligning it would
  move away from the spelling an official source documents. The unblock is unchanged:
  the documented request is a required JSON body whose only required property is a
  `category_symbols` array, while the SDK models one category with one flat symbol list
  and encodes it as query parameters. That is still the real work, and it still needs
  the paid entitlement to confirm which contract the host honours.
  `POST /market-data/stocks/snapshots/list`. It is the only `/openapi/…`
  holdout in a const block whose four siblings (`pathDSBars`,
  `pathDSBarsSingle`, `pathDSTick`, `pathDSDepth`, lines 35-41) were aligned in
  commit `2b29c88`, the same commit that deleted the
  `TODO(ds): Confirm exact paths against US sandbox` marker covering it. The
  reference page is self-contradictory — `operationId` `snapshotUsingGET` against
  `method` `post` — so this is genuinely ambiguous. Impact:
  `data.GetDisplaySnapshot`. The size is larger than a path-and-verb change, and
  correcting the request is a **breaking public API change**. The documented
  request is a JSON body with `requestBody` `required: true`, whose only required
  property is `category_symbols`: an **array** whose items each require a
  `category` string enum and a `symbols` **array of strings**, at most 100
  symbols per query, with `extend_hour_required` and `overnight_required`
  documented as **strings** defaulting to `"false"`. The SDK models one category
  with one flat symbol list (`data/snapshot.go:35-47`) and encodes it as query
  parameters (`data/display_quotes.go:52-60`): `:54` comma-joins `Symbols` into a
  single string, `:57` sends one flat `category` instead of a list of
  `{category, symbols}` objects, and `:63` sends a GET instead of the documented
  POST. The array-versus-scalar and GET-versus-POST mismatches are the real work;
  the path is the smallest part of it. The two flags are already
  type-compatible, since `strconv.FormatBool` at `data/display_quotes.go:59-60`
  emits `"true"`/`"false"`, matching the documented string type. The breaking
  element is `data.SnapshotQuery` itself: it is the public parameter type of both
  `GetDisplaySnapshot` (`data/display_quotes.go:51`) and `GetSnapshot`
  (`data/snapshot.go:142`), so moving it to the documented array breaks every
  consumer of either method. Minimal fix: not a patch — it needs a new public
  request shape, and therefore a maintainer decision on how to version a
  breaking change on a module that stays on the v1 import path (an added field
  with a deprecation window, a new method alongside the retained old one, or a
  `/v2` migration declined by decision so far). The path can still be aligned
  once probed; interim, restore an explicit unverified marker. Unblock: a paid
  Display Solution entitlement; the host returns `403` without one. Probing alone
  does not settle which contract the host honours, because the page contradicts
  itself and a second cached page for the same path, `reference/snapshot.md`,
  documents it as `method: get` with flat `symbols` and `category` query
  parameters; which contract applies to a Display client is a question for
  Webull. Every `file:line` in this list is anchored to the cited file as it
  stood on 2026-09-27 — the `data/display_quotes.go` numbers in this bullet, the
  `tools/webull-docgen/_common.py` numbers in the bullet above, and any
  `tools/webull-docgen/docgen.py` numbers — because a comment insertion shifts
  every line below it. A single 4-line comment added above the `broker-fd-us`
  manifest entries in one release invalidated both `_common.py` citations in that
  bullet, and a second pass that corrected one of them left it one line short, so
  re-check these against the file after any edit to it rather than assuming they
  still hold.
- **RESOLVED: `brokerfd.GetFDPositions` now receives all eight required response properties,
  so three values are silently zeroed.** The request is correct:
  `GetFDPositions` (`brokerfd/assets.go:92-100`) requests
  `pathFDAssetsPositions` (`brokerfd/assets.go:27`), exactly the documented
  `GET /broker/assets/positions/list`, so the reconciliation row renders as a
  clean `✅ match` and emits no signal — the same path-only limitation recorded
  above, which compares path strings and cannot observe a response schema at all.
  The defect is entirely on the response side. The documented `200` for that path
  is a `type: array` whose `items` require eight properties — `cost_price`,
  `currency`, `instrument_type`, `last_price`, `position_id`, `quantity`,
  `symbol`, `unrealized_profit_loss` — while `FDPosition`
  (`brokerfd/assets.go:98-121`) declares tags for `position_id`, `account_id`,
  `symbol`, `quantity`, `average_cost`, `market_value`, `unrealized_pl`,
  `realized_pl`, `instrument_type`, `currency`. **Three of the eight required
  names have no matching tag**: `cost_price` (the struct has `average_cost`),
  `last_price` (the struct has `market_value`), and `unrealized_profit_loss`
  (the struct has `unrealized_pl`).
  - Impact: the failure is **silent**, and that is what makes it the most
    dangerous of the five. All three fields are `money.Money`, so a missing JSON
    key leaves the value at its zero value and `encoding/json` reports no error:
    the caller gets a successful call, a non-nil slice, and three wrong zeroes
    where position cost basis, last price, and open P&L belong. Nothing in the
    return path distinguishes that from a genuinely zero position, so the defect
    is undetectable without a second source of truth. This is the opposite of the
    path defects above, which fail loudly with a `404` or a wrong-host error a
    caller cannot miss, and the reconciler cannot detect it either.
  - Size, restated: a retag is not automatically the fix, and the direction of
    the retag is the open question. The two shapes diverge in both directions,
    which is why a bare retag cannot be assumed correct: the struct also declares
    `account_id` and `realized_pl`, which the documented properties do not
    include at all, while the documented properties include `event_outcome`,
    which the struct omits. If the server sends the published names, the fix is
    to retag the three fields to `cost_price`, `last_price`, and
    `unrealized_profit_loss`; if the server actually sends the SDK's names, the
    documentation is describing a different payload, the struct is right, and
    retagging would break a call that currently works. The probe below has to
    establish which before anything is applied.
- Adjacent finding, **RESOLVED in v2.1.33**: `GetFDAssetsDetail`
    (`brokerfd/assets.go:97-105`), bound to `pathFDAssetsDetail =
    "/broker/assets/balances/get"` (`brokerfd/assets.go:26`), requests that
    documented path, whose documented `200` is a single `type: object`
    (`AssetsBalanceResult`, requiring `account_currency_assets`,
    `total_asset_currency`, and `total_cash_balance`), but the method used
    to decode into `[]FDAssetDetail`, a slice. A JSON object cannot be
    decoded into a Go slice, so that one failed loudly at decode time with
    an `encoding/json` type error rather than returning zeroes, which is why
    it was recorded here so the two were not mistaken for the same severity.
    The method now returns `*FDAssetsDetail` (`brokerfd/assets.go:82-91`),
    the object the page documents, which closed the shape and decode rows
    along with the three required names. **This was a breaking return-type
    change**; a caller reading the per-currency entries now reads
    `out.AccountCurrencyAssets[i]`. The container-type mismatch was certain
    from the source and the cached page; the corrected shape is not
    live-verified.  - Minimal fix: retag the three fields to the documented names, or add the
    documented names alongside the existing ones, **only after** the probe below
    establishes which set the server sends. Do not retag speculatively.
  - Unblock: a US sandbox credential — the same Broker FD grant that unblocks
    the first two `brokerfd` bullets — and, because whether the server sends the
    published names or the SDK's cannot be resolved from the documentation, one
    probe question added to the existing Webull enquiry: *for
    `GET /broker/assets/positions/list`, is the response the documented
    `AssetsPositionResult` carrying `cost_price`, `last_price`, and
    `unrealized_profit_loss`, or a payload carrying `average_cost`,
    `market_value`, and `unrealized_pl`?* That is a single addition to the
    enquiry the second bullet already requires for its three undocumented paths,
    so it needs no new round trip.
  - **Resolution, and a defect found while resolving it.** The type now carries
    `cost_price`, `last_price` and `unrealized_profit_loss` **alongside** the
    `average_cost`, `market_value` and `unrealized_pl` it always had, rather than one
    set replacing the other. That route was chosen because it is the one that needs no
    probe: whichever set the live server sends, the caller now reads a populated
    value, so the change cannot break a call that works today. A retag would have
needed the probe this bullet names, and would have been breaking. A response
populates whichever of the two names it carries, so a caller reading either name
reads a value when the server sends it; which of the two a live server sends is
precisely what stays unverified. `broker.Position` carried the identical
three names for the identical reason and is fixed the same way. These 6 baseline
    rows are deleted, and `brokerfd.GetFDPositions` records no divergence of any kind
    for the first time.
  - **A second, worse defect was found in the same struct, and the harness could not
    see it.** `RealizedPL` was tagged `json:"unrealized_pl"`, duplicating
    `UnrealizedPL`. When two fields of one struct claim the same name,
    `encoding/json` drops **both** rather than picking one — so open P&L had been
    silently zero for every caller, and `realized_pl` was never decodable at all.
    Nothing reads as *missing* in that state, which is why looking for missing names
    never found it. `TestFDPositionHasNoDuplicateWireName` now pins the class, and a
    scan of all 365 structs in the root module and examples found this to be the only
    occurrence. A tag-lookup harness has a blind spot at the *field* level that a
    name-lookup harness does not.
  - Verification status: the static mismatch is **certain** — the struct tags
    and the cached page's `required` list were compared directly and no
    interpretation is involved. Which side the live server honours is
    **unverified**: no endpoint was called, no US credential was available, and
    nothing in this bullet is live-verified.

#### The measured response-contract class

Four of the bullets above are request defects; the `brokerfd.GetFDPositions` bullet
is the first response-side one. The distinction matters more than the count,
because it decides how a defect reaches a caller. A request defect fails loudly —
a `404`, a wrong-host error, a decode error no caller can miss. A missing response
tag does not fail at all: `encoding/json` ignores a key it has no tag for, leaves
the field at its zero value, and reports no error. That is why that bullet is the
sharpest of the five, and why the class below is the one to read first.

`conformance/` measures that class. Each of the 193 documented endpoints carries a
committed fixture that is a minimal conforming instance of the response schema
Webull publishes for it, generated from the documentation by
`tools/conformance/gen_fixtures.py` and never from an SDK type, so a fixture is a
snapshot of the documentation and the SDK type is what has to agree with it.
`conformance/known-divergences.json` records every divergence the comparison
finds, one entry per row, each with a reason and a pointer to where it is written
down; the gate in `go test ./conformance/` fails when the observed set stops being
exactly that set, in either direction. `make conformance-report` prints the whole
observed set grouped by check.

**What the instrument does not prove.** It shows whether a documented response
shape is representable in the SDK's types, and nothing more. A row with no recorded
divergence is not a verified-correct row. The name check can only bite where the
page publishes a `required` list, and 90 of the 154 compared rows come from pages
that publish none, so their fixtures carry no name from Webull at all - 101 of the
193 fixtures are in that state. **That limit is now partly closed**: a second,
weaker name check runs over the property names those 90 pages do declare, which is
item 22's third direction and is applied as of this release. It is not closed
completely - 6 of the 90 remain unnameable, and why is in item 22.
`data.GetDisplaySnapshot` is the counter-example in the other direction: it
reconciles as a clean path match and is still defective.

Of the 154 compared rows, 63 carry at least one recorded divergence and 91 record
none; 5 of those 91 are rows the harness reports as not comparable, because the
method sends a path other than the one the page documents, so 86 comparable rows
record no divergence. All 25 divergences fall on 12 symbols, in three packages:
`brokerfd` 11, `data` 14, `trade` 0. The `trade` column is now **empty**: the four
rows it held were all `missing-required-name` on `trade.BatchPlaceOrder`, whose page
documents `total`, `success`, `failed` and `batch_orders` where the SDK had a single
`results` field appearing on no page. They are all carried now. The remaining `trade`
rows, if any, are the non-name classes, and the table below is the authority on that. The
4 `trade` rows were all on
`trade.BatchPlaceOrder`; 12 of the 13 trading endpoints the harness compared record
no divergence, and 4 of the 5 that publish a `required` list are among them, so the
trading API's response types are the most conformant part of the surface measured.
That is a statement about what the harness did not find, on the 5 trading pages that
give it something to check — not a correctness verdict.


- **The harness could not see a wrong HTTP verb, and now can.** All five
  pre-existing checks compare a documented *response* against what the SDK decodes;
  none looks at the *request* the SDK sends. So a method that sends the right path with
  the wrong verb reconciles as a clean path match - which is what
  `broker.UpdateVirtualAccount` did, reported as a `match`. `conformance.CompareVerb`
  now reads the verb the method actually sends **from the Go source** and compares it
  with `documented.method`, which the generator already recorded for all 193 fixtures.
  Reading it from source rather than a hand-maintained table follows the precedent
  `conformance/envelopes.go` set for unexported types: a table is a second place to
  forget to update. The extractor follows a method, a sub-service helper, an
  `http.MethodX` passed as an argument, a verb written as a string literal, and
  same-package delegation - which it must, because `trade.GetOpenOrders` delegates to
  `GetOpenOrdersPage` and `display.Service.EnsureToken` to `fetchToken`.

  On the committed surface: **183 mapped symbols, all readable; 178 rows agree, 0
  disagree, 15 not comparable on path.** Zero is the honest outcome only because the
  `UpdateVirtualAccount` fix landed first; reverting that one line produces exactly one
  finding naming it with both sides. **Two honest refusals:** `client.CreateToken` and
  `client.CheckToken` each send several verbs, so they have no single counterpart verb
  and are reported as not comparable rather than forced into an answer; and a verb the
  extractor cannot read yields no divergence at all, because that is a defect in the
  extractor, so it is required to be zero by a test rather than recorded in the
  baseline, which would put a tool defect in the one file whose purpose is to record
  SDK defects against the documentation.

- **A duplicated `json` tag made a value silently unreadable, and only a
  whole-module scan finds the class.** `brokerfd.FDPosition` tagged both
  `UnrealizedPL` and `RealizedPL` `json:"unrealized_pl"`, and `encoding/json` drops
  **both** fields when two fields of one struct claim the same name - with no error
  reported - so open P&L had been silently zero for every caller. The duplicate-tag
  defect and the missing-name defect look identical from outside, a value the caller
  cannot read, and only the level at which you look tells them apart: this harness
  looks for names that are *missing*, and nothing was missing, because `tagsOf` keeps
  one of the two colliding names and reported the name as covered while the decoder
  filled neither field. `conformance.ScanDuplicateTags` now walks every non-test struct
  in the root module, `broker/` and the four nested example modules - **370 structs
  across 6 module directories** - and matches on the name before the comma, so
  `json:"a"` and `json:"a,omitempty"` are correctly a collision. A nested module is
  excluded from the parent walk by its own `go.mod` and scanned as its own entry, so
  no struct is counted twice. Reintroducing the original duplicate tag is confirmed to
  be caught, naming both structs the one careless edit would have hit. The two
  per-type tests stay, because a per-type test pins the specific defect with its
  reasoning while the scan only says the class is empty.
- **`missing-required-name` — 2 rows across 1 symbol, the silent class.** A
  name the page marks `required` has no matching json tag anywhere in the type the
  method decodes into, so the documented value decodes to the zero value and no
  error is reported. This is the `brokerfd.GetFDPositions` mechanism at scale, and
  it is the largest class. The check considers the union of json tags at every depth
  of the compared type, because `encoding/json` flattens a response body into one name space, so the class
size is a
floor rather than a total. When this class stood at 119 rows, 19 of them
carried 41 of
their 90 required names only through a nested object, which the harness
printed per row
under the `required-name-depth` skip; the 2 rows left do not, so that figure
is history
rather than a property of the class today. Depth scoping is deliberately not
  enforced, because a documented top-level name the SDK reaches one level down
  still decodes and reporting it would trade a floor for a class of false
  positives — so names sitting at the wrong depth are not measured at all, by
  design. The check can fire only on a page that publishes a `required` list: 64 of
  the 154 compared rows do, and it fired on 29 of them. The worst-affected symbols
  are `brokerfd.GetFDCreditInfo` and `data.GetEventSnapshot` at 10 rows each,
  `brokerfd.AddFDAchAccount` and `brokerfd.AddFDBankAccount` at 8,
  `brokerfd.ListFDAchAccounts` and `brokerfd.ListFDBankAccounts` at 7, and
  `data.GetEventDepth` at 4 of its 6. Minimal fix direction: align the tags DTO by
  DTO against the documented `required` list, in the direction a live probe
  establishes — the `brokerfd.GetFDPositions` bullet records at length why a retag
  is not automatically the fix, and that same question governs every row in this
  class. Unblock: the US sandbox grant the first two `brokerfd` bullets already
  name, for the 89 `brokerfd` rows; HK sandbox credentials for the 29 `data` rows,
  less the US-only `data` surfaces among them, which need the same US grant; and the
  same HK credentials for the 4 `trade` rows.
- **`missing-declared-name` - 0 rows. The class is closed.** A property the
page
declares but does not mark `required` has no matching json tag in the type
the method
decodes into, so a response carrying it would decode the value to its zero
value and
report no error. The mechanism is identical to the class above and the
severity is the
same; what differs is the evidence, and the two must never be summed as one
figure. A
`required` name is a promise the page makes; a declared name is a
description of one
response, so a name absent from a page that does not require it may be
optional,
conditionally sent, or simply absent from the example the page chose. The
check runs
only where a page publishes no `required` list, which is exactly where the
class above
cannot run at all, and it is a separate check and a separate kind so the two
strengths
never collapse into one number.

**It held 104 rows across 30 symbols, and it closed in three releases.** 76
of them
were added as fields in v2.1.33 and v2.1.34 -- genuine gaps in the same
sense as the
required half, since a caller reading a name the page publishes got an empty
string, a
zero or an empty slice -- but on the weaker evidence, and each field says so
in its own
GoDoc rather than being presented as a documented requirement. The remaining
28 were
not fields at all: every one was `data` or `pagination_key` on the 14
methods that
returned a bare slice where the page documents a `{data, pagination_key}`
envelope, so
the name check had stopped finding missing fields and started reporting the
same
missing wrapper 28 times. v2.1.35 gave those methods the envelope, and the
class went
to zero. **The count reaching zero is a statement about this class only, not
about the
surface: it means no page's declared property is now unreachable, which is a
narrower
claim than the surface being correct.**


- **`declared-inventory-empty` — 1 row, 1 symbol, and it is a hole rather than a
  defect.** `brokerfd.ListAccountForms` on `broker-fd-us/GET-broker-forms-list`
  declares no property name at all, so neither name check had anything to look at. It
  is a gap in Webull's published documentation, not in the SDK, and it is recorded as
  its own kind precisely so it does not read as a row that examined clean when it was
  not examined — the specific failure this instrument exists to remove. The row
  already carried a container-kind and a `decode-failure` row, so this adds no new
  divergent row; it makes an existing one say why it cannot be examined further.
  Unblock: a written answer from Webull, or one captured body.
- - **Container kind - 10 rows across 10 symbols** (9
    `top-level-shape-mismatch` and 1
`element-type-mismatch`). The documented top level, or an array's element
kind, and
the type the method decodes into disagree. Unlike the name classes, this
fails loudly,
at decode time, with an `encoding/json` type error. That is the same
loud-versus-silent distinction the `brokerfd.GetFDPositions` bullet draws
for its
adjacent finding, and it is why the two severities must not be merged into
one number.
**All 10 are indeterminate, and are more likely documentation errors than
SDK
defects.** 5 sit on pages whose fixture name ends `-list`:
`data.GetDSLatestNews`,
`data.GetDSMarketNews`, `data.GetDSSymbolNews` and `data.GetLogos` document
a single
item's own fields at the top level of a path ending `/list`, which is
implausible for
a list endpoint and is the strongest documentation-error candidate in the
set, and
`data.GetDSNewsSummary` documents a `ChatStreamResponse` the SDK does not
model. The
other 5 sit on pages whose name does not end `-list`:
`brokerfd.GetAgreementDetail`,
`brokerfd.GetFDCorporateActions` and `data.GetCapitalFlow` document an array
where the
SDK returns one object, `brokerfd.GetFDTransferFees` documents one object
where the SDK
returns an array, and `brokerfd.ListAccountForms` is the free-form row. The
12 that
used to be an unwrapped `{data, pagination_key}` envelope are gone: v2.1.35
unwrapped
it. The harness has no basis to prefer the page over the type, so it records
all 10 as
open. **Do not read all 10 as SDK defects**: one live body per pattern would
settle
them, and the 4 single-item pages are the ones to put to Webull first.
Minimal fix
direction: establish the direction first, then correct whichever side is
wrong - a
retag is not the fix in this class. Unblock: the same credentials, plus a
written
answer from Webull on the 4 pages that document a single item where the path
says
list.

- - **`decode-failure` - 11 rows across 11 symbols, the weakest check.** A
    decode that
succeeds is consistent with a type that ignores every documented name, and a
decode
that fails is usually the same defect one of the other checks already named,
which is
why this check runs last. All 11 restate a row for the same symbol: 10
duplicate a
container-kind row and 1 duplicates a `declared-inventory-empty` row. They
are recorded
so the count stays stable, not because they are 11 additional defects, and
no separate
work is attached to them. This is the arithmetic reason the class totals
must not be
read as 25 independent problems.
- **`leaf-type-mismatch` — no rows.** The class had two, both
`data.Quote.QuoteTime`, and both are fixed; see item 25. It is the only class in
this table that is now empty, which is worth stating plainly rather than leaving
a row of zeroes to be read as a class that was measured and found clean.
  `data.GetQuotes` and `data.GetDisplayDepth` each publish `quote_time` as
  `type: string` and require it, and both SDK fields are `int64`. `data.GetQuotes`
  cites at its own GoDoc the same reference page the harness read, so the SDK
  contradicts a contract it names. These two rows were previously waived on the
  reasoning that the sandbox sends a JSON number; the waiver was removed as
  unsound, because a published type does not stop applying because one environment
  disagrees, and for the Display host there was no evidence at all — that host
  returns `403` at host level here. Minimal fix direction: a decode that accepts
  both a string and a number, which is a decision that wants live traffic rather
  than a mechanical retag. Unblock: HK sandbox credentials for the
  `data.GetQuotes` row, a paid Display Solution entitlement for the other.

Per-symbol detail is deliberately kept out of this page.
`conformance/known-divergences.json` carries every row with its own reason and its
own `recordedIn` pointer, `make conformance-report` prints the whole observed set
grouped by check, and `conformance/doc.go` states what the instrument covers and
    what it does not. No row of this class has a home here, and the 2 that had one were
    `brokerfd.GetFDAssetsDetail` pair recorded above; the other three that once had
    one, the `brokerfd.GetFDPositions` missing-name rows, no longer exist, because
    that type now carries the documented names. Each of the remaining 280 points at
    item 21, which records the class rather than the row, and says so in its own
    `recordedIn`.

- **Item 21 — the response contracts across the API surface are now measured, and
  they are largely divergent.** The five bullets above are one host, some path
  literals, a verb and body, a request shape, and one response DTO. This entry is
  the class that last one belongs to, taken across the surface: 25 recorded
    divergences on 12 symbols, in `brokerfd` (11), `data` (14), and `trade` (0).

  | Class | Rows | Symbols | `brokerfd` | `data` | `trade` | How it fails |
  |---|---:|---:|---:|---:|---:|---|
| `missing-required-name` | 2 | 1 | 2 | 0 | 0 | **Silently** - the documented value decodes to its zero value, no error reported |
| `missing-declared-name` | 0 | 0 | 0 | 0 | 0 | **Silently**, on weaker evidence - a name the page describes but does not require |
| `declared-inventory-empty` | 1 | 1 | 1 | 0 | 0 | Nothing was examined; the page declares no property at all |
| Container kind (`top-level-shape-mismatch` 26, `element-type-mismatch` 1) | 11 | 11 | 4 | 7 | 0 | Loudly, at decode time |
| `decode-failure` | 11 | 11 | 4 | 7 | 0 | Restates a container-kind or leaf-type row for the same symbol |
| `leaf-type-mismatch` | 0 | 0 | 0 | 0 | 0 | **Empty.** Both rows were `data.Quote.QuoteTime`, fixed by a type that reads both shapes Webull publishes; see item 25 |
| **Total** | **25** | **12** | **11** | **14** | **0** |

  Three caveats travel with those numbers. **The name checks report a floor,
not a
total**: each considers the union of json tags at every depth, so a name
reachable
only through a nested object counts as covered. The harness prints the
required half
per row under the `required-name-depth` skip and deliberately does not print
the
declared half, because the same test on an array row would fire on nearly
every
declared name. **The two silent name classes must not be summed as one
figure.** The
2 required-name rows rest on a promise the page makes, and the 28
declared-name rows
on a description of one response, and only the first is a defect of the same
certainty as the class above. The declared half is now empty, so the 25 rows
are
the 2 required-name rows and the 23 container-kind and decode rows. **All 10
of the
container-kind rows are indeterminate, and are more likely documentation
errors than SDK
defects**: 5 sit on pages whose name ends `-list`, of which 4 document a
single item
where the path says list, and the other 5 sit elsewhere. The envelope group
that used to
make up 12 of the 27 is gone: v2.1.35 unwrapped it. The harness has no basis
to prefer
the page over the type, so it records them as open; do not read all 10 as
SDK defects.
  - Impact: every SDK response type in `brokerfd`, and much of `data`. The strong
    name check could bite on only 64 of the 154 compared rows, the ones whose page
    publishes a `required` list, and it fired on 29 of those 64. The weak name
    check reaches the 90 pages that publish none, and fired on 30 of them, 9 of
    which recorded no divergence of any kind before it ran. The classes that matter
    most are the two silent ones, because nothing in the return path
    distinguishes a zeroed field from a genuine zero.
  - **No SDK fix is applied.** No `*.go` file under `data/`, `trade/`, `brokerfd/`
    or `broker/` was changed. The only code changed is the `conformance/` harness
    itself, and it changed behaviour only for the envelopes whose tags it could not
    previously read - see the envelope note below.
  - **One harness defect was fixed in the same release, and it changed the set.**
    `conformance/envelopes.go` stored a reconstructed envelope's struct tag as
    `go/ast` reports it, which is the source literal with its backticks, so
    `reflect.StructTag` could not read a single tag and `wireNameOf` fell back to
    the Go field name. That fallback agrees for camelCase and silently loses a
    snake_case name, so 4 rows were false positives - `data.GetCorporateActions`,
    `data.GetCorporateActionsByMarket`, `data.GetCryptoInstruments` and
    `data.GetOptionContracts`, each on `pagination_key`. The fault was latent
    because the test that round-trips an envelope compared the rebuilt tag against
    the tag the parser read, so both sides were wrong together. It is fixed,
    `TestEnvelopeTagsAreUsableAsStructTags` fails on the previous code, and 4 rows
    this class would have recorded are absent.
  - **Nothing here is live-verified.** No endpoint was called and no credential was
    used. The static comparison is certain; which side a live server honours is
    unverified for every row, exactly as for the `brokerfd.GetFDPositions` bullet.
  - - Not 25 defects: 11 rows are the `decode-failure` class restating another
    row for the same symbol, 2 are required-name rows resting on a container-kind row
    for the same symbol, 19 are the indeterminate container-kind rows, and 1 is an
    evidence-base hole.
  - Minimal fix direction: settle the direction per class before changing a tag, as
    the subsection above sets out. The order is the silent name class first, then
    the container-kind class once one live body establishes which shape the server
    sends, then the two leaf-type rows. The 104 declared-name rows come after the
    119 rather than beside them, because a name the page does not require is also a
    name that may legitimately be absent.
  - Unblock: US sandbox credentials for the `brokerfd` rows and for the US-only
    `data` surfaces; HK sandbox credentials for the rest of `data` and for `trade`;
    a paid Display Solution entitlement for `data.GetDisplayDepth`; and a written
    answer from Webull on the 4 pages that document a single item where the path
    says list.
  - **The `broker/` module is now measured too, by the same code.** Its 29
    documented endpoints were out of scope here until this release, because the
    module is separate and no symbol in the root module can name a broker type.
    That is no longer the case: `broker/conformance_test.go` resolves all 29 rows
    from inside the module and runs the identical five checks, and
    `broker/conformance-divergences.json` records what it found — **59 further
    divergences over 28 compared rows, 1 row decoding no body and 0 not
    comparable**. The counts above remain the root module's own; this entry's 81
    does not include the 59, and the two sets are reported separately rather than
    merged, because merging them would lose the fact that one comes from a module
    the root cannot see. Two findings are worth naming because they are not new
    defects but the *same* defect on a second surface: `broker.GetPositions` is
    missing `cost_price`, `last_price` and `unrealized_profit_loss`, the identical
    three names `brokerfd.GetFDPositions` is missing, so a caller reading a Broker
    HK position gets the same three silently zeroed values it already got from
    Broker FD. `broker.UpdateVirtualAccount` is additionally missing
    `client_request_id`, which is the field item 18's verb-and-body defect is
    about. All 59 remain open and none is live-verified; Broker API HK returns
    `401 ROUTE_NOT_PERMITTED` in the HK sandbox, so the unblock is a production or
    US-scoped Broker credential.
  - Scope limit: nothing in the 193-fixture manifest is now unexamined. The 29
    `broker/` rows are covered as above. What remains out of reach is different:
    the root module still cannot name a broker type, so a *change* to a broker type
    is only visible to the broker module's own job, not to the root gate. Item 22
    is the limit underneath this entry's counts rather than beside them: it records
    what a row with no recorded divergence does and does not establish.

- **Item 22 — for most of the measured surface, a clean conformance row meant
  unexamined rather than correct, and the weaker half of that gap is now closed.**
  Both name checks and the leaf-type check walk documented property names, so all
  three can bite only on a page that publishes a `required` list. **90 of the 154
  compared rows come from pages that publish none, and 101 of the 193 fixtures are
  in that state.** On those 90 rows the harness had no name from Webull to check a
  json tag against, so a clean result there carried no information about names. This
  is promoted out of the caveat in the subsection above to a tracked entry, because
  a caveat can be missed and this was the most consequential thing the whole
  measurement uncovered.
  - Impact: the limit was in the evidence base — in Webull's published
    documentation, not in the SDK and not in the instrument. Those pages declare
    no `required` array, so there was no documented property list to check the
    types against and no fixture change could put a name where no list existed to
    take one from. It bore on every clean row among the 90, and it is why rows
    recording no divergence across all 154 never read as correct endpoints. It
    remains a coverage limit and never a defect count.
  - **The third direction is now applied rather than recorded.** Those 90 pages do
    publish property names — 548 in total, against the 281 required names across
    the 64 — and the manifest already carried them as
    `checks.declaredTopLevelNames` and `checks.declaredElementNames`. A second
    name check now runs over that inventory, in its own check and its own
    divergence kind, on every row where the stronger check cannot run. It found
    **104 rows over 30 symbols, and 9 of the 30 rows recorded no divergence of
    any kind before it ran** — so the 67 fully clean rows were not clean.
  - **The 90 therefore split 32/58 rather than 23/67**, and 58 rows still record no
    divergence of any kind. All 58 of those are rows the weak check examined and
    agreed with, on names rather than on a promise.
  - **What a clean row does and does not establish, stated so it cannot be
    misread.** It shows that no divergence was recorded against the checks that
    could run on that row. It is not a correctness verdict, and it is **not a
    claim that the SDK is correct on those rows** — unexamined and correct are
    different states and only one of them is evidenced. `data.GetDisplaySnapshot`
    is the standing counter-example in the other direction: a clean path match
    that is still defective. A green row is a statement about the harness, not
    about the endpoint. The limit is now narrower, not gone: a name the page never
    requires is weaker evidence than one it does, so a clean row among the 58
    attests that every name the page described is carried, which is a real
    statement but not the same as being right about the endpoint.
  - **6 of the 90 are still unnameable, and the reasons are different from one
    another.** `data.GetBalanceSheet`, `data.GetCashFlow` and `data.GetIncomeStatement`
    decode into a free-form `map[string]any`, which carries every documented name
    and declares none of them; their pages declare **105 names between them**, so
    a missing tag there would be the harness's own error rather than a finding, and
    the check reports itself as not applicable instead. Without that guard, 68% of
    the whole 104 would have been false positives. `brokerfd.DownloadDocument` and
    `brokerfd.CancelFDOrder` decode no body at all. `data.GetStockInstruments` is one
    of the 5 not-comparable rows, so its page is not that call's contract.
    `brokerfd.ListAccountForms` is the 1 `declared-inventory-empty` row: its page
    declares no property name whatever. That leaves **84 of the 90 actually
    checkable**, and the 6 are named here rather than left as an absence.
  - **The `trade` split, so the 90 is not read as the whole surface.** 13 trading
    endpoints are compared; 5 publish a `required` list and 8 do not, 12 of the 13
    record no divergence, and all 4 `trade` divergence rows sit on the single one
    that does publish a list (`trade.BatchPlaceOrder`). Of the 12 clean, 4 are
    name-checkable and 8 are not. The trading API is the part of the surface where
    the instrument both ran and had documented names to check, and it found almost
    nothing — which is a statement about what the harness did not find on 5 pages,
    not a correctness verdict.
  - Minimal fix: there is no SDK change here, and the remaining honest direction is
    about evidence rather than code — obtain a documented `required` list for those
    endpoints, which is a Webull question whose natural form is to ask that the
    pages mark `required` what they always send. That is the only route that would
    upgrade 104 rows of weaker evidence into rows of the same standing as the 119.
    Failing it, the live probe in the next steps still moves rows, because a
    captured response body carries the server's own property names whether or not
    the page marked them required. A fixture change cannot help, and a further
    weaker check cannot either: the declared inventory is now fully consumed and 6
    rows have nothing left to consume.
  - Unblock: a written answer from Webull marking `required` the properties those
    pages always send. Failing that, live credentials for the probe, which is what
    the next steps describe.
  - **Nothing here is live-verified.** No endpoint was called and no credential was
    used.

- **Item 28 — a live run now measures the response contracts against a server and
  not only against a page, and it found one demonstrated defect the documented
  comparison is structurally blind to.** One authorised walk of all 193 documented
  endpoints against `api.sandbox.webull.hk` on 2026-09-29 recorded one outcome per
  endpoint (`examples/live-probe/census.go:635`); the 55 that answered HTTP 200
  (`conformance/testdata/live-manifest.json:210`) were reduced to value-free
  skeletons and committed under `conformance/testdata/live/`, and each was then
  compared against the SDK type **and** against the documentation
  (`conformance/live.go:243`). In one line: **of 55 probed, 2 disagree with their
  documentation and 12 disagree with the SDK.** That is 14 directional observations
  over **13 distinct findings**, because the one finding both runs reported is
  counted in each figure. The 13 are recorded one per row in
  `conformance/live-divergences.json:31`, each with a `reason` saying why it is
  recorded rather than fixed and an `unblock` saying what would close it, and the
  gate fails in both directions: a new finding is not recorded, and a recorded one
  that stops reproducing also fails, because a fixed finding and a check that
  changed meaning look identical from the gate.
  - **The set is separate from `known-divergences.json` on purpose.** A live finding
    is a claim about Webull's server. Recording one in the SDK's own backlog would
    make the documented rows and the live rows read as one list of defects, when
    the live rows are not claims about this SDK. The documented baseline is
    byte-unchanged by this run, and it stays the authority on the documented
    comparison.
  - **The census, and the bound it puts on this item.** 193 endpoints walked;
    the host answered 158 of them and 55 of those answered HTTP 200. The other 35
    are accounted for rather than missing: 34 are mutating endpoints the walk
    deliberately did not call
    (`conformance/testdata/live-manifest.json:40`) and 1 is blocked because no JSON
    request body is documented for it, so the probe must not invent one. Of the 103
    non-200 answers, 73 were `404`
    (`conformance/testdata/live-manifest.json:212`), 13 `417`, 9 `500`, 7 `403` and
    1 `504`; a non-200 body is Webull's error shape, so reducing one would compare
    an error against a documented success response and manufacture a divergence. By
    area, the surfaces that stayed unreachable are `broker-fd-us` (32 endpoints, all
    `404`), `broker-hk` (18, all `404`), `event-contracts` (11 `404`, 8 `500`),
    `display-solution` (6 `200`, 9 `404`, 2 `417`, 1 `504`), `market-data-crypto`
    (2 `403`, 1 `500`), `market-data-stock` (3 `200`, 3 `403`, 1 `417`),
    `market-data-futures` (5 `200`, 2 `403`, 1 `417`), `market-data-option`
    (1 `200`, 3 `417`) and `market-data-watchlist` (2 `200`, 6 `417`).
  - **The census cannot be reproduced from a clean checkout, and every coverage
    sentence here carries that.** The committed manifest is the endpoint inventory
    and is response-only; the request schemas the walk needs — query parameters,
    required fields, bodies — exist only in the gitignored docgen cache, which the
    probe reads at `examples/live-probe/census.go:1128` and refuses to proceed
    without. A fresh checkout can read the census and the captures and can
    re-derive every number in this item from them; it cannot re-derive them from a
    second run.
  - **The 13 are not 13 defects, and each row's `reason` says which one it is.** A
    reader who quotes the count as a defect count is misreading it.
    - **4 are decode rejections**, one per captured body the SDK's own type could
      not unmarshal, and all four are container-kind disagreements with the SDK.
      Three of them the sandbox answered a bare array where the SDK decodes an
      object: `data.GetDisplayGainersLosers` and `data.GetDisplayTopActive` decode
      `types.Page[data.ScreenerStock]` and `data.GetFuturesBars` decodes
      `data.BatchBars`, and for those three the page and the SDK agree with each
      other, which is why they carry no `top-level-shape-mismatch` and the shape
      check cannot see them. The fourth is the other way round — the sandbox
      answered the documented object and `data.GetStockInstruments` decodes
      `[]data.StockInstrument`, so there the page disagrees with the type and the
      row is recorded as a `top-level-shape-mismatch` as well. A decode that fails is the
      strongest signal this harness can produce, and it is the one class a
      documented comparison cannot see at all: a page and a fixture made to agree
      with each other decode cleanly whether or not either matches the server. It
      is not 4 demonstrated SDK defects either — whether the SDK should read the
      bare payload or the endpoint should wrap it is a question for Webull.
    - **3 are one container disagreement seen from three angles**, all on
      `data.GetStockInstruments`, where the live body is an object carrying `data`
      and `pagination_key` and the SDK decodes a slice. The shape check names the
      inversion, the decode check reports what a caller experiences, and the name
      check reports the two names an element type has nowhere to put. One
      disagreement, three rows, one fix.
    - **3 are absence of evidence rather than disagreement**, and are documentation
      observations. Two are on `trade.GetOrderDetail`, where the page requires
      `client_order_id` and `combo_type` and the sandbox held no order for the
      probed `client_order_id`, answering `{"orders": []}` — a statement about one
      probe against an empty account. One is on `trade.GetPositions`, where the
      page marks `option_strategy` required on the positions list and the probed
      account held an equity position, so the field is option-only and absent for
      every non-option holding: a required list naming a field the server does not
      send for equities is a promise the page does not keep, and that is a claim
      about the page.
    - **2 are the one documentation divergence this run settled**, below.
    - **1 is the single demonstrated defect**, below.
  - **The demonstrated defect: `data.GetFuturesTick` decodes `instrument_id` to `""`
    on every response it gets.** The page requires `instrument_id` and
    `data.StockTicks` tags the field `instrument_id` (`data/tick.go:66`), so the
    documented comparison is green and stays green — its fixture carries the name
    the page documents
    (`conformance/testdata/market-data-futures/GET-market-data-futures-ticks-list.json:1`).
    The sandbox sent `instrumentId`
    (`conformance/testdata/live/market-data-futures/GET-market-data-futures-ticks-list.json:3`).
    `encoding/json` matches a member name exactly and then case-insensitively, and
    an underscore is not a case, so the tag matches neither what the server sent nor
    what the page documents: `GetFuturesTick` (`data/futures_market.go:47`) returns
    a `StockTicks` whose `InstrumentID` is the empty string, with no error and
    nothing to indicate a value was missed. A doc-versus-SDK harness reports this
    endpoint as clean forever, because the fixture and the type were made to match
    each other.
    - **Why this row and not the other twelve.** It is the only one whose mechanism
      is provable from the committed tree without first deciding which side is
      right. The wire name cannot match the tag; that is a fact about
      `encoding/json`, not an inference about either party's intent. Every other row
      records an observation and a direction that is still undecided.
    - **Why it is recorded and not fixed.** Changing a public DTO's wire name is a
      breaking change for anyone already written against `StockTicks.InstrumentID`,
      and one host's answer cannot establish whether `instrumentId` is specific to
      this endpoint or the spelling this environment uses across the surface.
      Unblock: one probe of the same endpoint against a second host. Nothing else is
      needed to fix it, and nothing else will decide it. **No tag is changed by this
      item**, and none should be until that probe exists.
  - **What the run settled in the documented set, and what it could not.** Item 21's
    25 recorded documentation divergences cover 12 symbols, and the census reached
    **2 of the 25 rows**.
    - **`data.GetCapitalFlow` is settled, in the direction that the page is the
      outlier.** The page documents an object, the SDK decodes
      `[]data.CapitalFlowEntry`, and the sandbox answered an array — so the server
      and the SDK agree and the documentation is the odd one out. This is the one
      row the run **decides** rather than leaves open: item 21 recorded all 10
      container-kind rows as "more likely documentation errors than SDK defects"
      with no basis to choose, and this is one body that chooses. Both rows on that
      symbol now carry the direction, and the documented half closes when the page
      is corrected, which needs no credential.
    - **The other 23 rows could not be reached.** 11 sit on the `brokerfd` US-only
      surface, whose 32 probed endpoints all answered `404`. The other 12 sit on six
      `display-solution` symbols, in an area whose 18 probed endpoints answered
      6 `200`, 9 `404`, 2 `417` and 1 `504`: 10 of the 12 rows are on the `404`s and
      the remaining 2 on the single `504`. That last attribution is a property of the
      run rather than something a reader can re-derive — the committed tree carries
      the per-area status counts but not a per-endpoint status map for the endpoints
      that were never captured, because the census run artefact is gitignored. Until
      this run, "those endpoints are not reachable from the HK sandbox" was an
      assertion in this document; it is now a measurement with a per-row blocker.
  - **The reachable count is a bound, not a total, and what it leaves out is named
    here rather than discovered later.** The live tree covers 55 of the 193
    documented endpoints; the 138 it does not cover are the 34 mutating, 1 blocked
    and 103 non-200 already accounted for. Three further exclusions belong on the
    record.
    - **The 29 `broker/` rows have no live coverage whatever.** They are part of the
      193 walked, and every one of them falls in the `404` or not-called part of the
      census, so that module's separately recorded divergences are exactly as
      live-unverified as they were before this run.
    - **5 of the 7 `data` methods v2.1.35 gave their documented
      `{data, pagination_key}` envelope were reachable, and 2 of those 5 do not
      decode the live body** — `data.GetDisplayGainersLosers`
      (`data/display_screener.go:42`) and `data.GetDisplayTopActive`
      (`data/display_screener.go:77`), which now return
      `types.Page[data.ScreenerStock]` where the sandbox answered a bare array. The
      other 3 decode cleanly, and the 7 `brokerfd` methods are in the unreachable
      part. So a breaking change two days old was exercised against a server for
      the first time here, and on 2 of its 14 methods the server answers a different
      shape than the page.
    - **The 34 mutating endpoints were not called**, deliberately: a census that
      reported its own refusals as the server's would report a sandbox restriction
      as a Webull behaviour. They are the one part of the documented surface this
      run says nothing about beyond their existence, and they are the reason a
      re-run of the census needs an explicit decision rather than a credential.
  - **What this evidence is, stated so it cannot be quoted as more than it is.**
    - **One run, one environment, one day.** Every figure here comes from one
      authorised walk of `api.sandbox.webull.hk` on 2026-09-29, and there is no second
      run to compare it against. A finding here is what that host answered once.
      Nothing here says what production answers, and a sandbox is a different
      deployment from the one a caller uses — which is the same reason the
      live-blocked defects above are not closed by it.
    - **The value-free tree is not provably value-free.** Every leaf is a fixed
      placeholder the reduction chose — `"1"`, `-1`, `true`, `null` — and a body that
      happened to contain those literals reduces to the same bytes as any other
      reading of that kind, so a positive integer a server sent is indistinguishable
      from a count the reduction wrote. The claim is that every leaf is the constant
      the reduction chose
      (`conformance/testdata/live-manifest.json:3`), and it is held by the reducer
      and by no one hand-editing the file, not by a proof. Do not describe the tree
      as proven to hold nothing.
    - **A live finding is an observation, not a verdict.** That is what the separate
      file, the per-row `reason` and the separate `unblock` requirements all encode.
  - Minimal fix direction: none of the 13 has a fix applied. The order is
    `data.GetFuturesTick` first, because it is the only one whose fix is known and
    whose remaining obstacle is a maintainer's decision rather than more evidence;
    then the four decode rejections, which need Webull to say which contract the
    host honours before any type changes; then the `data.GetStockInstruments`
    container question, which is a path decision first and a shape decision second.
    **No retag in this class should be applied on the strength of one sandbox body
    alone**, which is item 20's argument again, and the `instrumentId` row is the
    case where the tempting move is exactly that.
  - Unblock: a second host for the futures tick wire name, which is the only single
    question standing between the demonstrated defect and a fix; US sandbox or
    production credentials for the futures, screener and Broker FD rows; a paid
    Display Solution entitlement for the two `GetDisplay*` rows and the six
    `display-solution` pages; a written answer from Webull on the 4 pages that
    document a single item where the path says list; and an account holding at
    least one order and one option position for the two absence-of-evidence rows. A
    second run of the census needs no new credential at all — it needs the docgen
    cache, which is gitignored and has to be rebuilt.

## Remaining risks

- Current v2.1.1 repository-tagged behavior has not been newly live-verified;
  US-only, Display, Broker HK, entitlement-gated, SSE, and unsupported sandbox
  paths remain blocked or unverified.
- Broker FD events are raw-only, `OnData` omits request ID/timestamp, one active
  `Run` is supported, and no public option injects a non-zero raw subscribe
  bitmask.
- Stream callbacks/channels are synchronous; a slow handler or full
  `DropBlock` subscriber causes head-of-line delay until cancellation or close.
- The four summary-only matches are resolved as documentation drift: the SDK matches the official `llms.txt` summary and differs only from the OpenAPI JSON. One path still differs from both, the Broker FD assets summary; unresolved SDK paths
  are `0`.
- Five live-blocked SDK defects are recorded above, none live-verified; the
  `brokerfd` items need US sandbox credentials, and three of the second item's
  undocumented paths additionally need a written answer from Webull, which the
  silent response-schema defect extends with one further question;
  `broker.UpdateVirtualAccount` needs a production or US-scoped Broker credential,
  and `data.GetDisplaySnapshot` needs a paid Display Solution entitlement plus a
  maintainer decision on versioning a breaking public API change.
- **Those five are not the complete list.** Item 21 records a sixth entry covering
  25 measured response-contract divergences on 12 symbols, of which 23 were
  previously unwritten. Most of that class fails silently, which is worse than
  failing loudly. It now spans two silent name classes on different evidence, 119
  resting on a `required` promise and 104 on a name the page merely describes, and
  **the two must not be summed as one figure.**
- **A limit sits underneath that class, and item 22 is now its own entry.** 90 of
  the 154 compared rows come from pages that publish no `required` list at all, so
  the strong name check and the leaf-type check could not run there and a clean row
  was closer to unexamined than to correct. It is a coverage limit in the evidence
  base, not 90 defects, and not a claim the SDK is right on those rows. **A second,
  weaker name check now runs over the property names those pages do declare, and it
  is not a formality: it produced 104 rows over 30 symbols, 9 of them on rows that
  had recorded no divergence of any kind before it ran.** The 90 now split 32/58
  rather than 23/67, and 58 rows are clean against names rather than merely
  unexamined. 6 rows remain unnameable for three different reasons, named in item
  22, which leaves 84 of the 90 checkable. The `trade` surface splits — 5 of its 13
  compared endpoints publish a `required` list and 8 do not — so the 90 does not
  describe the whole surface. Nothing in the 193-fixture manifest is now
  unexamined: the 29 `broker/` endpoints, out of reach for want of a nameable
  type until this release, are covered by `broker/conformance_test.go` under the
  same checks and carry **59 further recorded divergences of their own**.
- Nested coverage and strict docs are not CI gates; percentages are measurements,
  not behavior guarantees.
- **A third divergence set now exists and is measured against a server, not against
  a page.** Item 28 records 13 findings from one HK sandbox walk on 2026-09-29, and
  the document above is right that the other entries are not the complete list: the
  13 are not in the live-blocked set either, because a live finding is a claim about
  Webull's server. **Only one of the 13 is a demonstrated SDK defect** —
  `data.GetFuturesTick` sends `instrumentId` where `data.StockTicks` tags
  `instrument_id`, so the field decodes to `""` silently and a documented
  comparison reports the endpoint as clean forever. The other 12 are observations:
  4 decode rejections whose direction is undecided, 3 rows of one container
  disagreement, 3 absence-of-evidence rows against an empty account and an equity
  position, and the 2 rows of the one documentation divergence the run settled. The
  run covers 55 of 193 documented endpoints, so it bounds rather than closes the
  surface, and it is one host on one day.

## Next steps

1. Live-verify the v2.1.1 repository-tagged request, OMS, stream, and event telemetry with suitable non-production access.
2. Add US sandbox verification for US-only surfaces; this also unblocks most of the second `brokerfd` defect, the `brokerfd.GetFDPositions` response-schema defect, and the 89 `brokerfd` rows of the response-contract class in item 21. The three paths with no documented counterpart need a written answer from Webull, and that enquiry should also ask whether `GET /broker/assets/positions/list` returns the documented `AssetsPositionResult` names or the SDK's own, which the documentation cannot settle.
3. Work the response-contract class in item 21, and settle the direction before changing a tag. The cheapest step is the 4 pages that document a single item where the path says list, which one live body or one written answer would resolve; the `brokerfd` silent-name rows come next, and no retag should be applied on the strength of the documentation alone. The 104 `missing-declared-name` rows come **after** the 119 rather than beside them, because a name the page never requires is also a name that may legitimately be absent, so they need the same live evidence before a retag.
4. Shrink the 90 name-unexamined rows of item 22, which is now a residual rather than an untouched surface: **84 of the 90 are already examined and 6 are named as unnameable**, and the declared-inventory route those 6 would need is fully consumed. What is left is evidence, not instrumentation. The 84 checkable rows fall into only 3 documented shape families — 60 a bare object, 29 an array of objects, 1 an array of strings — so capturing one live response body per family is a bounded task, and a captured body carries the server's own property names whether or not the page marked them `required`. Add to the existing Webull enquiry a request that the affected pages mark `required` the properties they always send, which is the only route that would upgrade the 104 rows of weaker evidence to the standing of the 119.
5. **Withdrawn, and do not perform it as written.** It proposed resolving the four summary-only matches to the OpenAPI JSON paths, which would move `client.CreateToken` off the path every authenticated call depends on. All four are the `/openapi/*` namespace and the SDK matches the official `llms.txt` summary in each; the OpenAPI JSON records a path reorganisation. The only genuinely open path is the Broker FD assets summary, which is the `/broker-fd/*` item. The label fix that cleared the false unresolved flags is in `tools/webull-docgen/`; do not hand-edit the generated report.
7. Decide whether Broker FD needs public subscribe-bitmask, richer raw metadata, and all-runs lifecycle APIs before release.
8. Run module-aware race, vet, formatting, lint, and `mkdocs build --strict` before any separately approved release tag.
9. Probe `data.GetFuturesTick` against a second host, and only then decide the `instrumentId` retag. It is the one demonstrated defect the 2026-09-29 run produced (item 28), and the second host is the only thing that decides whether the camelCase spelling is endpoint-specific or this environment's convention. Do not retag on one host's answer: `data.StockTicks` is a public type and the tag is a breaking change for anyone reading it.
