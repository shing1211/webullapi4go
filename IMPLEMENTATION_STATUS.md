# Implementation Status

Last updated: 2026-09-28

- Latest repository tag: **`v2.1.22`** (2026-09-28)
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
  `broker/` (measurements, not behavior guarantees)
- Live verification: partial historical evidence only; the current hardening
  was not newly live-verified. The repository tag is not a published
  Go-semver v2 module.

## Status definitions

| State | Meaning |
|---|---|
| Implemented | The public SDK surface exists in the current tree |
| Offline-tested | Credential-free unit tests or local HTTP/gRPC/MQTT fakes pass |
| Live-verified | Successfully exercised against an available Webull environment |
| Blocked | Live verification is prevented by credentials, scope, entitlement, host, or sandbox data |
| Repository-tagged | Included in the authorized repository Git tag `v2.1.1`; this is not a published Go-semver v2 module |
| Released | Included in a published module release; the unchanged root path means `v2.x` Git tags are not published Go-semver v2 modules |

Implementation and offline tests do not imply live verification.

## Current architecture

| Layer | Package | Purpose |
|---|---|---|
| Core | `client/` | Config, signing, tokens, HTTP transport, request pipeline, resilience |
| Services | `data/`, `stream/`, `trade/`, `events/`, `connect/`, `display/`, `brokerfd/`, `brokerfd/events/` | Public service clients remain at the repository root |
| Broker HK | `broker/` | Separate Go module |
| Shared foundations | `pkg/errors`, `pkg/observability`, `pkg/resilience`, `pkg/transport`, `pkg/types` | Reusable public primitives |
| Domain | `pkg/domain/money`, `pkg/domain/order` | `money.Money` and the OMS state machine |
| Convenience | `webull/` | Thin aliases for the core client; not a replacement service facade |

The full-service relocation to `pkg/`, root-service deprecation shims, and raw `decimal.Decimal` DTO migration are superseded and closed. The streaming path does not use `sync.Pool`.

## Service status

| Area | Implemented | Offline-tested | Live-verified | Blocked or unverified |
|---|---:|---:|---|---|
| Core `client/` | Yes | Yes | Core auth and selected HK HTTP paths | Current request-pipeline hardening is not newly live-verified |
| Authentication | Yes | Yes | Token create/check/ensure in HK sandbox | Production token activation still requires the Webull App 2FA flow |
| Market Data HTTP `data/` | Yes | Yes | Selected AAPL, fundamentals, watchlist, and related HK calls | US-only data and entitlement-gated calls remain partial |
| Market Data MQTT `stream/` | Yes | Yes | Basic MQTT-over-WebSocket receive/reconnect paths | Current channel, shutdown, and health hardening is offline-tested only |
| Trading HTTP `trade/` | Yes | Yes | Accounts, balances, positions, previews, and guarded order paths | Multi-leg/futures/event-contract live behavior is unavailable or rejected in HK |
| Trading Events `events/` | Yes | Yes | Basic stream path was previously exercised | Current telemetry/reconciliation behavior is offline-tested only |
| Connect OAuth `connect/` | Yes | Yes | — | US live access unavailable |
| Display Solution `display/` | Yes | Yes | — | HK host returns `403` before endpoint verification |
| Broker API HK `broker/` | Yes | Yes, separate module | — | HK sandbox returns `401 ROUTE_NOT_PERMITTED` because scope is missing |
| Broker FD HTTP `brokerfd/` | Yes | Yes | — | US-only; no US sandbox credentials, HK returns `404` |
| Broker FD Events `brokerfd/events/` | Yes | Yes | — | US-only; no live verification available |
| Shared foundations | Yes | Yes | Not applicable | — |

## Endpoint coverage

The generated [SDK ↔ API reconciliation](docs/reconciliation.md) is authoritative. Its 2026-09-26 snapshot reports:

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

The states partition the 209 implemented endpoints: 184 + 4 + 1 = 189 rows carry a verified path, and the remaining 20 are the 3 rows labelled `📄 No OpenAPI schema on page` plus 17 manifest entries deliberately mapped to no SDK symbol. 189 + 3 + 17 = 209, so no endpoint is missing. **The 3 is a label count, not a page count**: 7 gRPC reference pages embed no OpenAPI schema, and the 4 that the manifest also maps to no SDK symbol are labelled `➖ No SDK symbol` instead, because `_reconcile_data()` in `tools/webull-docgen/docgen.py` evaluates the unmapped branch before the no-openapi branch and the status table must partition the 209 rows. A further 3 rows have a JSON block that yields no `path` and are likewise labelled unmapped, so 10 rows in total have no usable official path. The two non-defect categories are why the previously quoted "25 unresolved" was wrong: 20 of those 25 entries were generator artifacts, and the remaining 5 were investigated individually — 4 were correct SDK code the generator could not statically follow, and 1 is the real path mismatch now reported as ⚠️.

The SDK therefore has no documented-only endpoint gaps, but it is inaccurate to call the snapshot a zero-discrepancy report: 4 summary-only and 1 differing remain. Do not hand-edit generated pages under `docs/webull-api/`; regenerate derived documentation with the doc generator when its inputs or the SDK change. The label correction that produced this partition lives in `tools/webull-docgen/`, not in the generated output.

**A `✅ match` row is a path comparison, not a correctness verdict.** The reconciler compares path strings only. It therefore reports `broker.UpdateVirtualAccount` as a clean `✅ match` even though the request behind that path is defective (item 18), and it cannot observe HTTP verbs, request bodies, transport-host routing, or response schemas at all, which makes it blind to items 16, 19, and 20 as well. Item 20 is the sharpest case: a correct path, a clean `✅ match`, and a response DTO that cannot receive three of the eight properties the endpoint requires, so three values are silently zeroed. The generator is doing its stated job of comparing documented paths against SDK paths, and this is a limit of what a path comparison can show rather than a defect in it. The consequence for readers is that a green row in `docs/reconciliation.md` is not evidence that an endpoint is correct. Nothing in that table speaks to the response contract at all; the response contracts are measured separately, and largely diverge — see item 21.

## v2.1.1 repository-tagged work

| Work | State | Verification |
|---|---|---|
| Category matching plus identity-specific semantic sentinels | Implemented | Offline-tested across public client/MQTT/event sentinels |
| HTTP 417 compatibility mapping and MQTT/event terminal mappings | Implemented | Offline-tested; 417 API-message caveat preserved |
| Shared `Do`, `DoBroker`, and `DoStream` attempt pipeline | Implemented | Offline-tested; not newly live-verified |
| One-based hook attempts, stable correlation IDs, status/latency telemetry, W3C propagation, and `SafeErrorText` | Implemented | Offline-tested with real OTel readers and local servers |
| Provider-order-independent resilience metrics and caller-owned custom breakers | Implemented | Offline-tested |
| OMS status snapshots, stale/terminal reconciliation, and correct failure-scene mapping | Implemented | Offline-tested |
| Account-scoped order registry and successful replace/cancel transitions | Implemented | Offline-tested |
| Stable auto-generated client order IDs without mutating caller requests | Implemented | Offline-tested |
| Compare-and-swap stream state and deterministic data-only health recovery | Implemented | Offline-tested with deterministic watchdog ticks/barriers |
| Terminal MQTT/channel shutdown, cancellation unblocking, and synchronous registration-order dispatch | Implemented | Offline-tested for quote/snapshot/tick channels |
| MQTT Connect cancellation, Connect/Close races, and late-callback suppression | Implemented | Offline-tested |
| Serialized resubscription and active-subscription replay | Implemented | Offline-tested |
| Trading and Broker FD gRPC event telemetry, including cancellation/failure outcomes | Implemented | Offline-tested with local gRPC servers |
| Trading `Close` cancellation of all active `Run` calls | Implemented | Offline-tested |
| Cancellation, strict data/event leak checks, and public resilience/transport/type/`webull` coverage | Implemented | Offline race/vet evidence recorded 2026-09-25 |
| Account monitor and updated order/stream examples | Implemented | Included in offline module build/test |
| Documentation reconciliation | Implemented | Final strict MkDocs/link/diff gates recorded in the F1 report |

These items are included in repository tag `v2.1.1` (2026-09-25). They are
not a published Go-semver v2 module, and the current hardening was not newly
live-verified.

## Version history

The `v2.x` rows below are repository-tag records, not published Go-semver v2
modules; `v2.1.17` is the current authorized repository tag and earlier rows are
historical. The root module path remains `github.com/shing1211/webullapi4go` and
stays on the v1 import path by decision.

| Version | Date | Scope | Status |
|---|---|---|---|
| v0.1.0 | 2026-09-18 | Auth, core HTTP, Market Data HTTP, and MQTT streaming | Released |
| v0.2.1 | 2026-09-18 | Accounts, balances, and positions | Released |
| v0.2.2 | 2026-09-18 | Stock-order lifecycle and queries | Released |
| v0.2.3 | 2026-09-18 | Market-specific rules and HK BCAN | Released |
| v0.2.4 | 2026-09-18 | Single-leg options orders | Released |
| v0.2.5 | 2026-09-18 | US combo orders | Released |
| v0.2.6 | 2026-09-18 | News SSE through `DoStream` | Released |
| v0.3.0 | 2026-09-18 | Trading events over gRPC | Released |
| v0.4.0 | 2026-09-18 | Market Data fundamentals | Released |
| v0.5.0 | 2026-09-19 | Display, corporate actions, fund, crypto, and screener surfaces | Released; historical provisional paths later reconciled |
| v0.6.0 | 2026-09-20 | Multi-leg options, futures validation, and discovery helpers | Released; historical provisional behavior later reconciled |
| v0.7.0 | 2026-09-21 | Event contracts, Broker HK, Broker FD, and Broker FD events | Released |
| v0.8.0 | 2026-09-21 | Reserved | No code release |
| v0.9.0 | 2026-09-22 | GoDoc coverage, HK derivatives helpers, and examples | Released |
| v0.9.1 | 2026-09-21 | Watchlist boolean response, `DoBroker`, and examples | Released |
| v0.9.2 | 2026-09-21 | Broker HK path correction and live scope finding | Released |
| v1.0.0 | 2026-09-22 | Futures fixes and API stability audit | Released |
| v1.0.1 | 2026-09-22 | Flexible futures unit and order-ID example fix | Released |
| v1.0.2 | 2026-09-22 | Removed undocumented endpoints and corrected the option-contract path | Released |
| v1.0.3 | 2026-09-22 | Generated API reference, reconciliation, and path probe | Released |
| v1.1.0 | 2026-09-22 | Full documented endpoint coverage and zero documented-only gaps | Released |
| v1.1.1 | 2026-09-23 | DevOps, security checks, multi-OS CI, leak tests, and fuzzing | Released |
| v2.0.0 | 2026-09-23 | Public error/transport/resilience/domain foundations and thin `webull` aliases; root services retained | Historical tag |
| v2.0.1 | 2026-09-23 | `money.Money` conversion in `data/` and `trade/` | Historical tag |
| v2.0.2 | 2026-09-23 | `money.Money` conversion in `brokerfd/` | Historical tag |
| v2.0.3 | 2026-09-23 | MQTT/paho goroutine cleanup compatibility in leak tests | Historical tag |
| v2.0.4 | 2026-09-24 | Clock correction, idempotency helpers, transport tuning, resilience preset | Historical tag |
| v2.0.5 | 2026-09-24 | Request interceptors/hooks and initial OMS integration | Historical tag |
| v2.0.6 | 2026-09-24 | Stream state/channels, structured logging, and OTel tracing | Historical tag |
| v2.0.7 | 2026-09-24 | HTTP latency, breaker, reconnect, and drop metrics | Historical tag |
| v2.0.8 | 2026-09-24 | MQTT and resubscription context cleanup | Historical tag |
| v2.0.9 | 2026-09-24 | Structured public error codes and wrapped sentinels | Historical tag |
| v2.1.0 | 2026-09-24 | `go vet` mutex-copy fixes | Historical tag; not a published v2 module |
| v2.1.1 | 2026-09-25 | Error specificity, request/OMS, stream/MQTT, event telemetry, cancellation/leak coverage, and documentation hardening | Historical tag; not a published v2 module |
| v2.1.2 | 2026-09-26 | OpenTelemetry v1.46.0 across all modules, GitHub Actions bumps, pinned `govulncheck`, Dependabot grouping, and documentation policy updates | Historical tag; not a published v2 module |
| v2.1.3 | 2026-09-26 | Honest nightly live-verification signal, removal of the dead `internal/` shims and duplicate resilience tree, and direct tests for the order reconciliation, MQTT, money, error, and region surfaces | Historical tag; not a published v2 module |
| v2.1.4 | 2026-09-26 | Nested module coverage gate for `broker`, per-module `govulncheck`, a strict MkDocs gate, and run-index and roadmap corrections | Historical tag; not a published v2 module |
| v2.1.5 | 2026-09-26 | Additive Broker FD event metadata delivery, the first stream dispatch benchmarks with recorded head-of-line measurements, and classification of every non-exact reconciliation state | Historical tag; not a published v2 module; tagged with a red CI run |
| v2.1.6 | 2026-09-26 | Correction of the nested coverage CI job: artifact-safe matrix labels, tolerant artifact upload, and a guarded coverage report step | Repository tag; not a published v2 module; first tag with fully green CI |

## Previously exercised HK surface

Earlier sandbox runs exercised the core token flow, selected AAPL market-data and fundamentals calls, watchlists, accounts, balances, positions, order preview paths, MQTT streaming, and Trading Events. This is not a claim that all 209 endpoints or all current hardening behavior passed a live test.

## Test status

- The module-aware offline race/vet checks recorded for this run passed on
  2026-09-25. `make test` and `make test-race` traverse the root, `broker/`, and
  nested example modules; root-only `go test ./...` does not.
- `make cover` measured 73.6% aggregate coverage for the root module and 80.8%
  for `broker/` on 2026-09-26. These dated measurements are not a guarantee of
  behavior or correctness and must not be combined into one aggregate
  percentage.
- Unit tests are offline and credential-free. Live tests use the actual
  `Sandbox` selector and remain environment-gated and skipped by default.
- Data and both event packages use strict goroutine-leak checks; DNS-dependent
  failure tests use local deterministic dialers.
- CI runs root race tests across three operating systems and build/vet/race for
  nested modules. Its 60% coverage gate is root-only; nested coverage and a
  strict MkDocs build are not CI gates.
- No live verification was added for the v2.1.1 repository-tagged hardening.

## Examples

| Example | Directory | Purpose |
|---|---|---|
| auth | `examples/auth/` | Token create/check/reuse |
| marketdata | `examples/marketdata/` | AAPL snapshot and bars |
| streaming | `examples/streaming/` | MQTT-over-WebSocket subscription |
| watchlist | `examples/watchlist/` | Read-only watchlists |
| account | `examples/account/` | Accounts, balances, and positions |
| account-monitor | `examples/account-monitor/` | Bounded, read-only periodic balance monitor |
| order | `examples/order/` | Guarded AAPL preview/place/cancel flow with persisted idempotency key |
| events | `examples/events/` | Trading order events over gRPC |
| data-fundamentals | `examples/data-fundamentals/` | Fundamentals endpoints |
| watchlist-cmd | `examples/watchlist-cmd/` | Explicitly gated watchlist mutations |
| broker-probe | `examples/broker-probe/` | Broker HK read-only probe |
| futures-probe | `examples/futures-probe/` | Futures discovery and market-data probe |
| options-multi-leg | `examples/options-multi-leg/` | Multi-leg strategy probe |
| options | `examples/options/` | HK options discovery probe |
| brokerfd | `examples/brokerfd/` | Broker FD US probe |
| brokerfd-events | `examples/brokerfd-events/` | Broker FD gRPC probe |
| probe | `examples/probe/` | General sandbox probe |
| path-probe | `examples/path-probe/` | Generated path discrepancy probe |

## Blocked and unverified surfaces

This list and the live-blocked defects under it are not the complete set of known
divergences. Both are request-side or single-DTO findings. The response contracts
across the whole API surface have now been measured, and they are largely
divergent: see item 21 and its subsection. A further limit sits underneath that
measurement, and it is a limit of the evidence rather than of the SDK: 90 of the
154 compared rows come from pages that publish no `required` list. A weaker name
check now runs over the property names those pages do declare, and it is not a
formality — it produced 104 further rows, 9 of them on rows that had recorded no
divergence at all — but 6 of the 90 remain unnameable and the evidence on the rest
is a description rather than a promise. That is item 22.

1. **Display Solution:** the HK host returns `403`, including token creation.
2. **US surfaces:** no US sandbox credentials are available for crypto, funds, screener, Broker FD, or other US-only behavior.
3. **Broker API HK:** returns `401 ROUTE_NOT_PERMITTED` because the app lacks the required scope.
4. **HK symbols:** sandbox data is limited to `AAPL`.
5. **Footprint:** requires an entitlement and returns `403` in the sandbox.
6. **Option contracts:** may be absent for the sandbox account/symbol and return `417 Invalid Symbol`. HTTP 417 is mapped to `INVALID_TOKEN` for compatibility even when the API message reports business validation.
7. **Multi-leg options:** HK rejects non-`SINGLE` strategies with `417`; US behavior is unverified.
8. **Futures/event-contract trading:** request validation is offline-tested; live product behavior is unverified.
9. **Order-book depth:** may be empty outside regular trading hours.
10. **Plain MQTT:** port `1883` may be blocked; prefer the configured MQTT-over-WebSocket endpoint.
11. **SSE news:** the HK upstream currently returns `504`.
12. **Generated path status:** four paths match only the docs summary and one differs from both; unresolved SDK paths are now `0`. See the generated reconciliation rather than claiming zero discrepancies.
13. **Broker FD event API:** data remains raw, `OnData` omits response request ID/timestamp, only one active `Run` is supported, and no public option injects a non-zero raw subscribe bitmask.
14. **Stream backpressure:** callbacks and channel dispatch are synchronous; a slow handler or full `DropBlock` subscriber causes head-of-line delay until cancellation or terminal close.
15. **CI gaps:** nested coverage and strict documentation builds are Makefile/release gates, not CI gates; measured percentages are not behavior guarantees.

### Live-blocked SDK defects

The first four were found by static analysis on 2026-09-26 and item 20 on 2026-09-27, also statically. None is live-verified — no pass touched the network — and each may only be changed once the credential or entitlement it names is available, because the fix cannot be confirmed without it.

16. **`brokerfd` sends every request to the core host.** `brokerfd/client.go:43` calls `c.core.Do(ctx, method, path, body, out)`, sending all Broker FD traffic to the Trading/Market Data host. The Broker FD reference pages document `https://broker-api.sandbox.webull.com`, which is the `BrokerHTTP` host in `internal/region/region.go:191`; the HK `broker` package routes correctly through `DoBroker` at `broker/client.go:63`.
    - Impact: every method in the `brokerfd` package. This is the highest-impact item of the four.
    - Minimal fix: route `brokerfd`'s shared `do` helper through the Broker host.
    - Unblock: US sandbox credentials. The HK host does not serve the FD surface (`404`), so HK cannot distinguish "wrong host" from "wrong region".
17. **`brokerfd` uses undocumented `/broker-fd/*` paths.** 14 non-test literals matching `/broker-fd/` remain in the package (`brokerfd/accounts.go:29-31`, `brokerfd/assets.go:25`, `brokerfd/brokerfd.go:58,68`, `brokerfd/documents.go:23-24`, `brokerfd/funding.go:26,30`, `brokerfd/instruments.go:27,29,31`, `brokerfd/journals.go:25`), while every other cached `broker-fd-api` page uses the `/broker/...` namespace. Concretely `brokerfd/assets.go:25` sends `/broker-fd/assets/summary` against a documented `GET /broker/assets/summaries/get`; that is the single ⚠️ path-differs-from-both row in the current snapshot. `brokerfd.GetPositions` (`brokerfd/brokerfd.go:67-68`) maps to no documented page, and the assets-summary manifest entry at `tools/webull-docgen/_common.py:314` names two SDK symbols for that page — `brokerfd.GetAccountsSummary` (`brokerfd/brokerfd.go:57-58`, path `/broker-fd/accounts`) has no page of its own, and `brokerfd.GetFDAssetsSummary`'s response DTO (`brokerfd/assets.go:32-41`) does not match the documented `balance`/`positions` envelope. A separate defect on the same Broker FD surface is item 20, a `brokerfd` FD positions response-schema mismatch rather than a request-path defect; it is recorded separately rather than merged here because it is the one defect on this surface that fails silently, where every literal in this item fails loudly.
    - Impact: the affected `brokerfd` endpoints, including the assets summary.
    - Size, restated: only a minority of the 14 literals have an unambiguous documented counterpart, so this is not a uniform one-line path fix. 4 can be aligned mechanically — `brokerfd/assets.go:25` → `/broker/assets/summaries/get`, `brokerfd/brokerfd.go:68` → `/broker/assets/positions/list`, `brokerfd/instruments.go:29` → `/broker/instruments/stocks/corporate-actions/get`, and `brokerfd/journals.go:25` → `/broker/journals/cash-journals/get`, each confirmed against the cached `broker-fd-api` page carrying that path. The remaining 10 split three ways: 1 is a probable duplicate (`brokerfd/brokerfd.go:58` overlaps the already-aligned `/broker/accounts/list` at `brokerfd/accounts.go:23`, consistent with the single-symbol List Accounts manifest entry at `tools/webull-docgen/_common.py:299`, which shows only `brokerfd.ListFDAccounts` and so cannot surface `GetAccountsSummary`'s different path); 6 are plausibly ambiguous, because a documented page exists in the same area but not for the same operation (`brokerfd/accounts.go:29-31` form detail/submit/status, `brokerfd/funding.go:26,30`, and `brokerfd/instruments.go:31`, against pages such as `/broker/forms/get` and `/broker/accounts/applications/get`); and 3 have no documented counterpart at all — `brokerfd/instruments.go:27` stock-locate, `brokerfd/documents.go:23` documents, and `brokerfd/documents.go:24` documents/detail. For those 3 the cache holds only `/broker/documents/download` and `/broker/documents/upload` in the documents namespace and no stock-locate page. The 4 + 1 + 6 + 3 buckets are a plausibility assessment of the 14 literals, not a claim that any mapping is settled.
    - Minimal fix: align the 4 unambiguous literals, reconcile the duplicate symbol pair in the docgen manifest, and investigate the other 9 before changing them. Do not guess the 3 that have no documented counterpart.
    - Unblock: US sandbox credentials; probe `/broker/assets/summaries/get` and `/broker-fd/assets/summary` side by side. Credentials alone cannot settle the 3 with no documented counterpart, because a `404` from a US host does not distinguish an undocumented path from an endpoint Webull does not offer; those 3 need a written answer from Webull about whether the endpoints exist at all.
18. **`broker.UpdateVirtualAccount` sends the wrong verb and body shape.** `broker/accounts.go:66-68` builds `pathVirtualAccountsUpdate + "?account_id=" + accountID` and issues `c.put`. The documented endpoint is POST and requires `account_id` AND `client_request_id` in the JSON body, with no `account_name` field at all; `UpdateVirtualAccountRequest.AccountName` (`broker/accounts.go:54`) is undocumented. `GetVirtualAccount` does document `account_id` as a query parameter, so the POST appears to have copied the GET's convention. The path string itself is correct, which is why reconciliation reports it as a match — the generator compares paths, not verbs or bodies. The existing test `broker/accounts_test.go:151-160` currently certifies the wrong contract.
    - Impact: `broker.UpdateVirtualAccount`.
    - Minimal fix: issue POST with the documented body fields and update the test alongside it.
    - Unblock: a production or US-scoped Broker credential; the HK sandbox returns `401 ROUTE_NOT_PERMITTED` for the whole Broker API.
19. **`data.GetDisplaySnapshot` differs from both official sources.** `data/display_quotes.go:33` sets `pathDSSnapshot = "/openapi/market-data/stock/snapshot"` and `:63` sends it with GET. Both official sources say `POST /market-data/stocks/snapshots/list`. It is the only `/openapi/…` holdout in a const block whose four siblings (`pathDSBars`, `pathDSBarsSingle`, `pathDSTick`, `pathDSDepth`, lines 35-41) were aligned to `/market-data/stocks/…` in commit `2b29c88`, and the `TODO(ds): Confirm exact paths against US sandbox` marker that covered it was deleted in that same commit. The reference page is self-contradictory — its `operationId` is `snapshotUsingGET` while its `method` is `post` — so this is genuinely ambiguous rather than settled.
    - Impact: `data.GetDisplaySnapshot`; the four sibling Display endpoints are aligned.
    - Size, restated: larger than a path-and-verb change, and correcting the request is a **breaking public API change**. The documented request is a JSON body whose `requestBody` is `required: true`, and whose only required property is `category_symbols` — an **array** whose items each require a `category` string enum and a `symbols` **array of strings**, at most 100 symbols per query. `extend_hour_required` and `overnight_required` are documented as **strings** defaulting to `"false"`. The SDK models one category with one flat symbol list (`data/snapshot.go:35-47`) and encodes it as query parameters (`data/display_quotes.go:52-60`): `:54` comma-joins `Symbols` into a single string rather than sending a JSON array, `:57` sends one flat `category` rather than a list of `{category, symbols}` objects, and `:63` sends a GET rather than the documented POST. The array-versus-scalar and GET-versus-POST mismatches are therefore the real work, and the path alone is the smallest part of it. The two flags are the one part already type-compatible: `strconv.FormatBool` at `data/display_quotes.go:59-60` emits `"true"`/`"false"`, which coincides with the documented string type.
    - Breaking element: `data.SnapshotQuery` is the public parameter type of both `GetDisplaySnapshot` (`data/display_quotes.go:51`) and `GetSnapshot` (`data/snapshot.go:142`), so moving it to the documented `category_symbols` array breaks every consumer of either method, not only the Display one. That makes the fix semver-major for a module deliberately held on the v1 import path.
    - Minimal fix: not a patch, and not a maintainer-mechanical change either. Aligning the request needs a new public request shape — a changed `data.SnapshotQuery` or a separate Display-specific type — which requires a maintainer decision on how to version a breaking change on a module that stays on the v1 path: an added field with a deprecation window, a new method alongside the retained old one, or a `/v2` migration that has so far been declined by decision. The path can still be aligned once probed. Interim recommendation: restore an explicit unverified marker instead of leaving the divergence silent.
    - Unblock: a paid Display Solution entitlement; the host returns `403` without one. Probing alone does not resolve which contract the host honours, because the reference page is self-contradictory — its `operationId` is `snapshotUsingGET` while its `method` is `post` — and a second cached page for the same path, `reference/snapshot.md`, documents that path as `method: get` with flat `symbols` and `category` query parameters. Which contract applies to a Display Solution client is a question for Webull.
    - Citation scope: every `file:line` in this list is anchored to the cited file as it stood on 2026-09-27 — the `data/display_quotes.go` numbers in item 19, the `tools/webull-docgen/_common.py` numbers in item 17, and any `tools/webull-docgen/docgen.py` numbers — because a comment insertion shifts every line below it. A single 4-line comment added above the `broker-fd-us` manifest entries in one release invalidated both `_common.py` citations in item 17, and a second pass that corrected one of them left it one line short, so re-check these against the file after any edit to it rather than assuming they still hold.
20. **`brokerfd.GetFDPositions` cannot receive three required response properties, so three values are silently zeroed.** The request is correct: `GetFDPositions` (`brokerfd/assets.go:92-100`) requests `pathFDAssetsPositions` (`brokerfd/assets.go:27`), which is exactly the documented `GET /broker/assets/positions/list`, so the reconciliation row renders as a clean `✅ match` and emits no signal — the same path-only limitation recorded in the caveat above, which compares path strings and cannot observe a response schema at all. The defect is entirely on the response side. The documented `200` for that path is a `type: array` whose `items` require eight properties — `cost_price`, `currency`, `instrument_type`, `last_price`, `position_id`, `quantity`, `symbol`, `unrealized_profit_loss` — while `FDPosition` (`brokerfd/assets.go:79-90`) declares tags for `position_id`, `account_id`, `symbol`, `quantity`, `average_cost`, `market_value`, `unrealized_pl`, `realized_pl`, `instrument_type`, `currency`. **Three of the eight required names have no matching tag**: `cost_price` (the struct has `average_cost`), `last_price` (the struct has `market_value`), and `unrealized_profit_loss` (the struct has `unrealized_pl`).
    - Impact: the failure is **silent**, and that is what makes this the most dangerous of the five. All three fields are `money.Money`, so a missing JSON key leaves the value at its zero value and `encoding/json` reports no error: the caller gets a successful call, a non-nil slice, and three wrong zeroes where position cost basis, last price, and open P&L belong. Nothing in the return path distinguishes that from a genuinely zero position, so the defect is undetectable without a second source of truth. This is the opposite of the path defects in items 16 and 17, which fail loudly with a `404` or a wrong-host error a caller cannot miss, and the reconciler cannot detect it either.
    - Size, restated: a retag is not automatically the fix, and the direction of the retag is the open question. The two shapes diverge in both directions, which is why a bare retag cannot be assumed correct: the struct also declares `account_id` and `realized_pl`, which the documented properties do not include at all, while the documented properties include `event_outcome`, which the struct omits. If the server sends the published names, the fix is to retag the three fields to `cost_price`, `last_price`, and `unrealized_profit_loss`; if the server actually sends the SDK's names, the documentation is describing a different payload, the struct is right, and retagging would break a call that currently works. The probe below has to establish which before anything is applied.
    - Adjacent finding, loud rather than silent: `GetFDAssetsDetail` (`brokerfd/assets.go:65-75`), bound to `pathFDAssetsDetail = "/broker/assets/balances/get"` (`brokerfd/assets.go:26`), requests that documented path, whose documented `200` is a single `type: object` (`AssetsBalanceResult`, requiring `account_currency_assets`, `total_asset_currency`, and `total_cash_balance`), but the method decodes into `[]FDAssetDetail` (`brokerfd/assets.go:70`), a slice. A JSON object cannot be decoded into a Go slice, so this one fails loudly at decode time with an `encoding/json` type error instead of returning zeroes; it is recorded here so the two are not mistaken for the same severity, and not as an addition to their count. The container-type mismatch is certain from the source and the cached page; the behaviour is not live-verified.
    - Minimal fix: retag the three fields to the documented names, or add the documented names alongside the existing ones, **only after** the probe below establishes which set the server sends. Do not retag speculatively.
    - Unblock: a US sandbox credential — the same Broker FD grant that unblocks items 16 and 17 — and, because whether the server sends the published names or the SDK's cannot be resolved from the documentation, one probe question added to the existing Webull enquiry: *for `GET /broker/assets/positions/list`, is the response the documented `AssetsPositionResult` carrying `cost_price`, `last_price`, and `unrealized_profit_loss`, or a payload carrying `average_cost`, `market_value`, and `unrealized_pl`?* That is a single addition to the enquiry item 17 already requires for its three undocumented paths, so it needs no new round trip.
    - Verification status: the static mismatch is **certain** — the struct tags and the cached page's `required` list were compared directly and no interpretation is involved. Which side the live server honours is **unverified**: no endpoint was called, no US credential was available, and nothing in this entry is live-verified.

#### The measured response-contract class

Items 16 to 19 are request defects; 20 is the first response-side one. The
distinction matters more than the count, because it decides how a defect reaches
a caller. A request defect fails loudly — a `404`, a wrong-host error, a decode
error no caller can miss. A missing response tag does not fail at all:
`encoding/json` ignores a key it has no tag for, leaves the field at its zero
value, and reports no error. That is why 20 is the sharpest of the five, and why
the class below is the one to read first.

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
shape is representable in the SDK's types, and nothing more. A row with no
recorded divergence is not a verified-correct row. The name check can only bite
where the page publishes a `required` list, and 90 of the 154 compared rows come
from pages that publish none, so their fixtures carry no name from Webull at all
— 101 of the 193 fixtures are in that state. **That limit is now partly closed**:
a second, weaker name check runs over the property names those 90 pages do
declare, which is item 22's third direction and is applied as of this release.
It is not closed completely — 6 of the 90 remain unnameable, and why is in item 22.
Item 19 is the counter-example in the other direction: `data.GetDisplaySnapshot`
reconciles as a clean path match and is still defective.

Of the 154 compared rows, 63 carry at least one recorded divergence and 91 record
none; 5 of those 91 are rows the harness reports as not comparable, because the
method sends a path other than the one the page documents, so 86 comparable rows
record no divergence. All 285 divergences fall on 63 symbols, in three packages:
`brokerfd` 145, `data` 136, `trade` 4. The 4 `trade` rows are all on
`trade.BatchPlaceOrder`; 12 of the 13 trading endpoints the harness compared record
no divergence, and 4 of the 5 that publish a `required` list are among them, so the
trading API's response types are the most conformant part of the surface measured.
That is a statement about what the harness did not find, on the 5 trading pages that
give it something to check — not a correctness verdict.

**`missing-required-name` — 122 rows across 29 symbols, the silent class.** A name
the page marks `required` has no matching json tag anywhere in the type the method
decodes into, so the documented value decodes to the zero value and no error is
reported. This is item 20's mechanism at scale, and it is the largest class. The
check considers the union of json tags at every depth of the compared type, because
`encoding/json` flattens a response body into one name space, so **122 is a floor
rather than a total**: 19 rows carry 41 of their 90 required names only through a
nested object, which the harness prints per row under the `required-name-depth`
skip. Depth scoping is deliberately not enforced, because a documented top-level
name the SDK reaches one level down still decodes and reporting it would trade a
floor for a class of false positives — so names sitting at the wrong depth are not
measured at all, by design. The check can fire only on a page that publishes a
`required` list: 64 of the 154 compared rows do, and it fired on 29 of them. The
worst-affected symbols are `brokerfd.GetFDCreditInfo` and `data.GetEventSnapshot`
at 10 rows each, `brokerfd.AddFDAchAccount` and `brokerfd.AddFDBankAccount` at 8,
`brokerfd.ListFDAchAccounts` and `brokerfd.ListFDBankAccounts` at 7, and
`data.GetEventDepth` at 4 of its 6. Minimal fix direction: align the tags DTO by
DTO against the documented `required` list, in the direction a live probe
establishes — item 20 records at length why a retag is not automatically the fix,
and that same question governs every row in this class. Unblock: the US sandbox
grant the first two `brokerfd` items already name, for the 89 `brokerfd` rows; HK
sandbox credentials for the 29 `data` rows, less the US-only `data` surfaces among
them, which need the same US grant; and the same HK credentials for the 4 `trade`
rows.

**`missing-declared-name` — 104 rows across 30 symbols, the same silent class on
weaker evidence.** A property the page declares but does not mark `required` has no
matching json tag in the type the method decodes into, so a response carrying it
would decode the value to its zero value and report no error. The mechanism is
identical to the class above and the severity is the same; **what differs is the
evidence, and the count must not be added to the 122 as if it were.** A `required`
name is a promise the page makes; a declared name is a description of one response,
so a name absent from a page that does not require it may be optional,
conditionally sent, or simply absent from the example the page chose. Each of the
104 reasons says so, per entry. The check runs only where a page publishes no
`required` list, which is exactly where the class above cannot run at all, and it is
a separate check and a separate kind so the two strengths never collapse into one
number. The worst-affected rows are `brokerfd.GetFDCorporateActions` at 12 of the 14
names its page declares, `data.GetFinancialAlert` at 9 of 14 — carrying `eps_est`,
`rev_est`, `fiscal_year` and four more — and `data.GetFundInfo`,
`data.GetHighDividendRank` and `data.GetWeek52HighLow` at 6 each. The two screener
rows are the structurally interesting ones: `data.ScreenerStock` is shared with
`data.GetMostActive` and `data.GetTopGainersLosers`, which report no divergence at
all, so those pages declare names a shared DTO cannot carry for every one of its
callers. Minimal fix direction: the same as the class above, and item 20's warning
applies identically — no retag on the strength of the documentation alone, because
a name the page never promises is also a name that may simply be absent. Unblock:
the same credentials as the class above.
- **68 of the 104 sit on rows that also record a container-kind mismatch, and 26 of
those are the `data`/`pagination_key` pair** a bare-array SDK cannot carry. They are
recorded rather than suppressed, and the dead end is worth recording because the
obvious suppression is wrong. `brokerfd.ListFDAccounts` declares 2 names against a
page of 6, so its object is an envelope and its names are the wrapper;
`data.GetDSLatestNews` declares 6 against a page of 6, so its object is the payload
and its 5 missing names are a second, independent defect. `data.GetMarketSectorDetail`
is both at once, carrying 4 content names beside 2 wrapper keys. Suppressing per
row would discard the content names, and no per-name test exists — the manifest
records no type per declared name, and `propertyNameCountInFixture` is 0 on these
rows because the committed minimal instances are empty. This is the same treatment
the 29 `decode-failure` rows already get.
- **The floor applies here too, and is not printed.** 104 is a floor for the reason
122 is: the union of tags at every depth counts a name as covered. Unlike the
required half, the depth accounting is deliberately not reported per row, because
`tagsOf` is called with the slice type on an array row, so the element's own fields
land one level down and the same test would fire on nearly every declared name of
every array row as a false "reached from below".

**`declared-inventory-empty` — 1 row, 1 symbol, and it is a hole rather than a
defect.** `brokerfd.ListAccountForms` on `broker-fd-us/GET-broker-forms-list`
declares no property name at all, so neither name check had anything to look at. It
is a gap in Webull's published documentation, not in the SDK, and it is recorded as
its own kind precisely so it does not read as a row that examined clean when it was
not examined — the specific failure this instrument exists to remove. The row already
carried a container-kind and a `decode-failure` row, so this adds no new divergent
row; it makes an existing one say why it cannot be examined further. Unblock: a
written answer from Webull, or one captured body.


`element-type-mismatch`). The documented top level, or an array's element kind, and
the type the method decodes into disagree: 23 of the 27 are a documented bare object
against an SDK array, 3 are the inverse, and 1 is an element kind. Unlike the class
above, this fails loudly, at decode time, with an `encoding/json` type error. That
is the same loud-versus-silent distinction item 20 draws for its adjacent finding,
and it is why the two severities must not be merged into one number.

**19 of the 27 are indeterminate, and are more likely documentation errors than SDK
defects.** 19 sit on pages whose fixture name ends `-list`, and they divide three
ways: 12 document a `{data, pagination_key}` pagination envelope that the SDK does
not unwrap, 4 document a single item's own fields at the top level of a path ending
`/list` (`data.GetDSLatestNews`, `data.GetDSMarketNews`, `data.GetDSSymbolNews`,
`data.GetLogos`) — implausible for a list endpoint, and the strongest
documentation-error candidate in the set — and 3 are the inverse or element-kind
rows. The harness has no basis to prefer the page over the type, so it records all
27 as open. **Do not read all 27 as SDK defects**: one live body per pattern would
settle them, and the 4 single-item pages are the ones to put to Webull first. The
other 8 sit on pages whose name does not end `-list`; `brokerfd.GetFDAssetsDetail`
is among them and is already recorded under item 20. Minimal fix direction:
establish the direction first, then either unwrap the documented envelope in the SDK
or correct the page — a retag is not the fix in this class. Unblock: the same
credentials, plus a written answer from Webull on the 4 pages that document a single
item where the path says list.

**`decode-failure` — 29 rows across 29 symbols, the weakest check.** A decode that
succeeds is consistent with a type that ignores every documented name, and a decode
that fails is usually the same defect one of the other checks already named, which
is why this check runs last. All 29 rows restate a row for the same symbol: 27
duplicate a container-kind row and 2 duplicate a `leaf-type-mismatch` row. They are
recorded so the count stays stable, not because they are 29 additional defects, and
no separate work is attached to them. This is the arithmetic reason the class
totals must not be read as 285 independent problems.

**`leaf-type-mismatch` — 2 rows across 2 symbols, both `quote_time`.**
`data.GetQuotes` and `data.GetDisplayDepth` each publish `quote_time` as
`type: string` and require it, and both SDK fields are `int64`. `data.GetQuotes`
cites at its own GoDoc the same reference page the harness read, so the SDK
contradicts a contract it names. These two rows were previously waived on the
reasoning that the sandbox sends a JSON number; the waiver was removed as unsound,
because a published type does not stop applying because one environment disagrees,
and for the Display host there was no evidence at all — that host returns `403` at
host level here. Minimal fix direction: a decode that accepts both a string and a
number, which is a decision that wants live traffic rather than a mechanical retag.
Unblock: HK sandbox credentials for the `data.GetQuotes` row, a paid Display
Solution entitlement for the other.

Per-symbol detail is deliberately kept out of this document.
`conformance/known-divergences.json` carries every row with its own reason and its
own `recordedIn` pointer, `make conformance-report` prints the whole observed set
grouped by check, and `conformance/doc.go` states what the instrument covers and
what it does not. 5 of the 285 rows already had a home here before this subsection
existed — the three `brokerfd.GetFDPositions` missing-name rows and the two
`brokerfd.GetFDAssetsDetail` rows, all under item 20 — and each of the remaining
280 points at item 21, which records the class rather than the row, and says so in
its own `recordedIn`.

21. **The response contracts across the API surface are now measured, and they are largely divergent.** The five items above are one host, some path literals, a verb and body, a request shape, and one response DTO. This item is the class that last one belongs to, taken across the surface: 285 recorded divergences on 63 symbols, in `brokerfd` (145), `data` (136), and `trade` (4).

    | Class | Rows | Symbols | `brokerfd` | `data` | `trade` | How it fails |
    |---|---:|---:|---:|---:|---:|---|
    | `missing-required-name` | 122 | 29 | 89 | 29 | 4 | **Silently** — the documented value decodes to its zero value, no error reported |
    | `missing-declared-name` | 104 | 30 | 31 | 73 | 0 | **Silently**, on weaker evidence — a name the page describes but does not require |
    | `declared-inventory-empty` | 1 | 1 | 1 | 0 | 0 | Nothing was examined; the page declares no property at all |
    | Container kind (`top-level-shape-mismatch` 26, `element-type-mismatch` 1) | 27 | 27 | 12 | 15 | 0 | Loudly, at decode time |
    | `decode-failure` | 29 | 29 | 12 | 17 | 0 | Restates a container-kind or leaf-type row for the same symbol |
    | `leaf-type-mismatch` | 2 | 2 | 0 | 2 | 0 | Loudly — a JSON type error on the documented value |
    | **Total** | **285** | **63** | **145** | **136** | **4** | |

    Three caveats travel with those numbers. **122 and 104 are both floors, not
    totals**: the name checks consider the union of json tags at every depth, so 19
    rows carry 41 of their 90 required names only through a nested object and count
    as covered; the harness prints each required half under the `required-name-depth`
    skip, and the declared half is deliberately not printed, because the same test on
    an array row would fire on nearly every declared name. **The two silent name
    classes must not be summed as one figure.** 122 rests on a promise and 104 on a
    description, and 68 of the 104 sit on rows that also record a container-kind
    mismatch, 26 of them restating that mismatch. **19 of the 27 container-kind rows
    are indeterminate, and are more likely documentation errors than SDK defects**:
    19 sit on pages whose name ends `-list`, of which 12 document a
    `{data, pagination_key}` envelope the SDK does not unwrap and 4 document a single
    item where the path says list. The harness has no basis to prefer the page over
    the type, so it records them as open; do not read all 27 as
    SDK defects.
    - Impact: every SDK response type in `brokerfd`, and much of `data`. The strong name check could bite on only 64 of the 154 compared rows, the ones whose page publishes a `required` list, and it fired on 29 of those 64. The weak name check reaches the 90 pages that publish none, and fired on 30 of them, 9 of which recorded no divergence of any kind before it ran. The classes that matter most are the two silent ones, because nothing in the return path distinguishes a zeroed field from a genuine zero.
    - **No SDK fix is applied.** No `*.go` file under `data/`, `trade/`, `brokerfd/` or `broker/` was changed. The only code changed is the `conformance/` harness itself, and it changed behaviour only for the envelopes whose tags it could not previously read — see the envelope note below.
    - **One harness defect was fixed in the same release, and it changed the set.** `conformance/envelopes.go` stored a reconstructed envelope's struct tag as `go/ast` reports it, which is the source literal with its backticks, so `reflect.StructTag` could not read a single tag and `wireNameOf` fell back to the Go field name. That fallback agrees for camelCase and silently loses a snake_case name, so 4 rows were false positives — `data.GetCorporateActions`, `data.GetCorporateActionsByMarket`, `data.GetCryptoInstruments` and `data.GetOptionContracts`, each on `pagination_key`. The fault was latent because the test that round-trips an envelope compared the rebuilt tag against the tag the parser read, so both sides were wrong together. It is fixed, `TestEnvelopeTagsAreUsableAsStructTags` fails on the previous code, and 4 rows this class would have recorded are absent.
    - **Nothing here is live-verified.** No endpoint was called and no credential was used. The static comparison is certain; which side a live server honours is unverified for every row, exactly as for item 20.
    - Not 285 defects: 29 rows are the `decode-failure` class restating another row for the same symbol, 68 are declared-name rows resting on a container-kind row for the same symbol, 19 are the indeterminate container-kind rows, and 1 is an evidence-base hole.
    - Minimal fix direction: settle the direction per class before changing a tag, as the subsection above sets out. The order is the silent name class first, then the container-kind class once one live body establishes which shape the server sends, then the two leaf-type rows. The 104 declared-name rows come after the 122 rather than beside them, because a name the page does not require is also a name that may legitimately be absent.
    - Unblock: US sandbox credentials for the `brokerfd` rows and for the US-only `data` surfaces; HK sandbox credentials for the rest of `data` and for `trade`; a paid Display Solution entitlement for `data.GetDisplayDepth`; and a written answer from Webull on the 4 pages that document a single item where the path says list.
    - Scope limit: the 29 `broker/` endpoints are not covered at all, being a separate Go module, so this item is not a statement about them. Item 22 is the limit underneath this item's counts rather than beside them: it records what a row with no recorded divergence does and does not establish.

22. **For most of the measured surface, a clean conformance row meant unexamined rather than correct — and the weaker half of that gap is now closed.** Both name checks and the leaf-type check walk documented property names, so all three can bite only on a page that publishes a `required` list. **90 of the 154 compared rows come from pages that publish none, and 101 of the 193 fixtures are in that state.** On those 90 rows the harness had no name from Webull to check a json tag against, so a clean result there carried no information about names. This is promoted out of the caveat in the subsection above to a tracked item, because a caveat can be missed and this was the most consequential thing the whole measurement uncovered.
    - Impact: the limit was in the evidence base — in Webull's published documentation, not in the SDK and not in the instrument. Those pages declare no `required` array, so there was no documented property list to check the types against and no fixture change could put a name where no list existed to take one from. It bore on every clean row among the 90, and it is why rows recording no divergence across all 154 never read as correct endpoints. It remains a coverage limit and never a defect count.
    - **The third direction is now applied rather than recorded.** Those 90 pages do publish property names — 548 in total, against the 281 required names across the 64 — and the manifest already carried them as `checks.declaredTopLevelNames` and `checks.declaredElementNames`. A second name check now runs over that inventory, in its own check and its own divergence kind, on every row where the stronger check cannot run. It found **104 rows over 30 symbols, and 9 of the 30 rows recorded no divergence of any kind before it ran** — so the 67 fully clean rows were not clean. `data.GetFDCorporateActions`-class findings were there all along behind a documentation page that never marked anything required.
    - **The 90 therefore split 32/58 rather than 23/67**, and 58 rows still record no divergence of any kind. All 58 of those are rows the weak check examined and agreed with, on names rather than on a promise.
    - **What a clean row does and does not establish, stated so it cannot be misread.** It shows that no divergence was recorded against the checks that could run on that row. It is not a correctness verdict, and it is **not a claim that the SDK is correct on those rows** — unexamined and correct are different states and only one of them is evidenced. `data.GetDisplaySnapshot` is the standing counter-example in the other direction: a clean path match that is still defective. A green row is a statement about the harness, not about the endpoint. The limit is now narrower, not gone: a name the page never requires is weaker evidence than one it does, so a clean row among the 58 attests that every name the page described is carried, which is a real statement but not the same as being right about the endpoint.
    - **6 of the 90 are still unnameable, and the reasons are different from one another.** `data.GetBalanceSheet`, `data.GetCashFlow` and `data.GetIncomeStatement` decode into a free-form `map[string]any`, which carries every documented name and declares none of them; their pages declare **105 names between them**, so a missing tag there would be the harness's own error rather than a finding, and the check reports itself as not applicable instead. Without that guard, 68% of the whole 104 would have been false positives. `brokerfd.DownloadDocument` and `brokerfd.CancelFDOrder` decode no body at all. `data.GetStockInstruments` is one of the 5 not-comparable rows, so its page is not that call's contract. `brokerfd.ListAccountForms` is the 1 `declared-inventory-empty` row: its page declares no property name whatever. That leaves **84 of the 90 actually checkable**, and the 6 are named here rather than left as an absence.
    - **The `trade` split, so the 90 is not read as the whole surface.** 13 trading endpoints are compared; 5 publish a `required` list and 8 do not, 12 of the 13 record no divergence, and all 4 `trade` divergence rows sit on the single one that does publish a list (`trade.BatchPlaceOrder`). Of the 12 clean, 4 are name-checkable and 8 are not. The trading API is the part of the surface where the instrument both ran and had documented names to check, and it found almost nothing — which is a statement about what the harness did not find on 5 pages, not a correctness verdict.
    - Minimal fix: there is no SDK change here, and the remaining honest direction is about evidence rather than code — obtain a documented `required` list for those endpoints, which is a Webull question whose natural form is to ask that the pages mark `required` what they always send. That is the only route that would upgrade 104 rows of weaker evidence into rows of the same standing as the 122. Failing it, the live probe in the next steps still moves rows, because a captured response body carries the server's own property names whether or not the page marked them required. A fixture change cannot help, and a further weaker check cannot either: the declared inventory is now fully consumed and 6 rows have nothing left to consume.
    - Unblock: a written answer from Webull marking `required` the properties those pages always send. Failing that, live credentials for the probe, which is what the next steps describe.
    - **Nothing here is live-verified.** No endpoint was called and no credential was used.

## Next steps

1. Live-verify the v2.1.1 repository-tagged request, OMS, stream, and event-telemetry changes when suitable credentials and non-production test access are available.
2. Supply US sandbox credentials for the US-only surfaces, which is also what unblocks items 16, most of 17, 20, and the 89 `brokerfd` rows of the response-contract class in 21. The three literals in item 17 with no documented counterpart additionally need a written answer from Webull, and item 20 needs one further question added to that same enquiry, namely whether `GET /broker/assets/positions/list` returns the documented `AssetsPositionResult` names or the SDK's own, which the documentation cannot settle.
3. Work the response-contract class in item 21, and settle the direction before changing a tag. The cheapest step is the 4 pages that document a single item where the path says list, which one live body or one written answer would resolve; the `brokerfd` silent-name rows come next, and no retag should be applied on the strength of the documentation alone. The 104 `missing-declared-name` rows come **after** the 122 rather than beside them, because a name the page never requires is also a name that may legitimately be absent, so they need the same live evidence before a retag.
4. Shrink the 90 name-unexamined rows of item 22, which is now a residual rather than an untouched surface: **84 of the 90 are already examined and 6 are named as unnameable**, and the declared-inventory route those 6 would need is fully consumed. What is left is evidence, not instrumentation. The 84 checkable rows fall into only 3 documented shape families — 60 a bare object, 29 an array of objects, 1 an array of strings — so capturing one live response body per family is a bounded task, and a captured body carries the server's own property names whether or not the page marked them `required`. Add to the item 20 enquiry a request that the affected pages mark `required` the properties they always send, which is the only route that would upgrade the 104 rows of weaker evidence to the standing of the 122.
5. Resolve the four summary-only matches and the one differing path through the doc generator and official OpenAPI sources. The label fix that cleared the false unresolved flags is in `tools/webull-docgen/`; do not hand-edit the generated report.
6. Correct `broker.UpdateVirtualAccount` (item 18) only against the credentials it names, and correct its test in the same change. Treat `data.GetDisplaySnapshot` (item 19) as a breaking-API decision rather than a patch: obtain the maintainer decision on versioning first, then probe before touching the path.
7. Decide whether Broker FD needs a public raw-subscribe option, richer `OnData` metadata, and all-runs lifecycle parity before any future release tag.
8. Benchmark a separately approved asynchronous stream-dispatch design only if synchronous head-of-line latency is unacceptable.
9. Run the full race, vet, formatting, lint, and strict documentation gates before any future release tag.
