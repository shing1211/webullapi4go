# Implementation Status

Last updated: 2026-09-26

- Latest repository tag: **`v2.1.8`** (2026-09-26)
- Current hardening: **tagged in repository `v2.1.4`**; introduced in `v2.1.1`
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

There are no documented-only endpoint gaps, but the snapshot is not a zero-discrepancy report: 4 summary-only and 1 differing remain. Generated pages under `webull-api/` are not hand-edited; the label correction that produced this partition lives in `tools/webull-docgen/`.

A `✅ match` row is a path comparison, not a correctness verdict. The reconciler compares path strings only: it reports `broker.UpdateVirtualAccount` as a clean `✅ match` although the request behind that path is defective, and it cannot observe HTTP verbs, request bodies, or transport-host routing at all, so it is blind to the `brokerfd` host-routing and `data.GetDisplaySnapshot` defects recorded below as well. That is a limit of what a path comparison can show rather than a defect in the generator, which is doing its stated job of comparing documented paths against SDK paths. A green row in `reconciliation.md` is not evidence that an endpoint is correct.

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
modules; `v2.1.6` is the current authorized repository tag and earlier rows are
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

Blocked or unverified areas:

- Display Solution: host-level `403`.
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

Found by static analysis on 2026-09-26. None is live-verified — nothing in that
pass touched the network — and each may only be changed once the credential or
entitlement it names is available, because the fix cannot be confirmed without
it.

- **`brokerfd` sends every request to the core host.** `brokerfd/client.go:43`
  calls `c.core.Do(...)`, so all Broker FD traffic goes to the Trading/Market
  Data host instead of the documented `https://broker-api.sandbox.webull.com`
  (`BrokerHTTP` in `internal/region/region.go:191`); the HK `broker` package
  routes correctly through `DoBroker` at `broker/client.go:63`. Impact: every
  method in the package, and the highest-impact item of the four. Minimal fix:
  route the shared `do` helper through the Broker host. Unblock: US sandbox
  credentials — the HK host does not serve the FD surface (`404`), so HK cannot
  distinguish "wrong host" from "wrong region".
- **`brokerfd` uses undocumented `/broker-fd/*` paths.** 14 non-test literals
  remain (`brokerfd/accounts.go:29-31`, `assets.go:25`, `brokerfd.go:58,68`,
  `documents.go:23-24`, `funding.go:26,30`, `instruments.go:27,29,31`,
  `journals.go:25`) while every other cached `broker-fd-api` page uses the
  `/broker/...` namespace. `brokerfd/assets.go:25` sends
  `/broker-fd/assets/summary` against a documented
  `GET /broker/assets/summaries/get` — the single ⚠️ row in the snapshot.
  `brokerfd.GetPositions` (`brokerfd/brokerfd.go:67-68`) maps to no documented
  page, and `tools/webull-docgen/_common.py:298` names two symbols for that one
  page: `brokerfd.GetAccountsSummary` (`brokerfd/brokerfd.go:57-58`) has no page
  of its own and `brokerfd.GetFDAssetsSummary`'s DTO
  (`brokerfd/assets.go:32-41`) does not match the documented
  `balance`/`positions` envelope. Impact: the affected `brokerfd` endpoints.
  The size is larger than the literal count suggests, because only a minority of
  the 14 have an unambiguous documented counterpart. 4 align mechanically —
  `assets.go:25` → `/broker/assets/summaries/get`, `brokerfd.go:68` →
  `/broker/assets/positions/list`, `instruments.go:29` →
  `/broker/instruments/stocks/corporate-actions/get`, and `journals.go:25` →
  `/broker/journals/cash-journals/get` — while 1 (`brokerfd.go:58`) is a probable
  duplicate, 6 are plausibly ambiguous (`accounts.go:29-31`, `funding.go:26,30`,
  `instruments.go:31`, where a documented page exists in the same area but not
  for the same operation), and 3 have no documented counterpart at all
  (`instruments.go:27` stock-locate, `documents.go:23` documents,
  `documents.go:24` documents/detail; the cache holds only
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
- **`broker.UpdateVirtualAccount` sends the wrong verb and body shape.**
  `broker/accounts.go:66-68` builds
  `pathVirtualAccountsUpdate + "?account_id=" + accountID` and issues
  `c.put`, while the documented endpoint is POST and requires `account_id` AND
  `client_request_id` in the JSON body with no `account_name` field;
  `UpdateVirtualAccountRequest.AccountName` (`broker/accounts.go:54`) is
  undocumented. `GetVirtualAccount` does document `account_id` as a query
  parameter, so the POST appears to have copied the GET's convention. The path
  is correct, which is why reconciliation reports a match — the generator
  compares paths, not verbs or bodies — and
  `broker/accounts_test.go:151-160` currently certifies the wrong contract.
  Impact: `broker.UpdateVirtualAccount`. Minimal fix: issue POST with the
  documented body fields and update the test in the same change. Unblock: a
  production or US-scoped Broker credential; the HK sandbox returns
  `401 ROUTE_NOT_PERMITTED` for the whole Broker API.
- **`data.GetDisplaySnapshot` differs from both official sources.**
  `data/display_quotes.go:28` sets
  `pathDSSnapshot = "/openapi/market-data/stock/snapshot"` and `:52` sends it
  with GET, while both official sources say
  `POST /market-data/stocks/snapshots/list`. It is the only `/openapi/…`
  holdout in a const block whose four siblings (`pathDSBars`,
  `pathDSBarsSingle`, `pathDSTick`, `pathDSDepth`, lines 30-36) were aligned in
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
  parameters (`data/display_quotes.go:41-49`): `:43` comma-joins `Symbols` into a
  single string, `:46` sends one flat `category` instead of a list of
  `{category, symbols}` objects, and `:52` sends a GET instead of the documented
  POST. The array-versus-scalar and GET-versus-POST mismatches are the real work;
  the path is the smallest part of it. The two flags are already
  type-compatible, since `strconv.FormatBool` at `data/display_quotes.go:48-49`
  emits `"true"`/`"false"`, matching the documented string type. The breaking
  element is `data.SnapshotQuery` itself: it is the public parameter type of both
  `GetDisplaySnapshot` (`data/display_quotes.go:40`) and `GetSnapshot`
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
  Webull.

## Remaining risks

- Current v2.1.1 repository-tagged behavior has not been newly live-verified;
  US-only, Display, Broker HK, entitlement-gated, SSE, and unsupported sandbox
  paths remain blocked or unverified.
- Broker FD events are raw-only, `OnData` omits request ID/timestamp, one active
  `Run` is supported, and no public option injects a non-zero raw subscribe
  bitmask.
- Stream callbacks/channels are synchronous; a slow handler or full
  `DropBlock` subscriber causes head-of-line delay until cancellation or close.
- Four summary-only matches and one differing path remain; unresolved SDK paths
  are `0`.
- Four live-blocked SDK defects are recorded above, none live-verified; the two
  `brokerfd` items need US sandbox credentials, and three of the second item's
  undocumented paths additionally need a written answer from Webull,
  `broker.UpdateVirtualAccount` needs a production or US-scoped Broker credential,
  and `data.GetDisplaySnapshot` needs a paid Display Solution entitlement plus a
  maintainer decision on versioning a breaking public API change.
- Nested coverage and strict docs are not CI gates; percentages are measurements,
  not behavior guarantees.

## Next steps

1. Live-verify the v2.1.1 repository-tagged request, OMS, stream, and event telemetry with suitable non-production access.
2. Add US sandbox verification for US-only surfaces; this also unblocks most of the second `brokerfd` defect, and the three paths with no documented counterpart need a written answer from Webull.
3. Resolve the four summary-only and one differing path states through the doc generator and official sources. The label fix that cleared the false unresolved flags is in `tools/webull-docgen/`; do not hand-edit the generated report.
4. Correct `broker.UpdateVirtualAccount` only against the credential it names, and update its test in the same change. Treat `data.GetDisplaySnapshot` as a breaking-API decision rather than a patch: obtain the maintainer decision on versioning a public `data.SnapshotQuery` change first, then probe before touching the path.
5. Decide whether Broker FD needs public subscribe-bitmask, richer raw metadata, and all-runs lifecycle APIs before release.
6. Run module-aware race, vet, formatting, lint, and `mkdocs build --strict` before any separately approved release tag.
