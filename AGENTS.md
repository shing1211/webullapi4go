# AGENTS.md

Repository guide for AI coding agents working in `github.com/shing1211/webullapi4go`,
an idiomatic Go SDK for the Webull OpenAPI. Read this before making changes.

## Layout

```
client/                         Core SDK: config, options, signing, tokens, transport entry point
data/                           Market Data HTTP endpoints (typed requests and responses)
stream/                         Market Data streaming over MQTT (reconnect + resubscribe)
trade/                          Trading HTTP API: accounts, balances, positions, orders, rules
events/                         Trading events over gRPC (order, position, option streams)
broker/                         Broker API HK (separate Go module; `broker/go.mod`: `replace github.com/shing1211/webullapi4go => ../`)
brokerfd/                       Broker FD US HTTP endpoints
brokerfd/events/                Broker FD events over gRPC + protobuf payloads
display/                        Display Solution client-token authentication
connect/                        Connect API OAuth 2.0 authorization-code flow
gen/webull/...                  Generated protobuf types (committed; do not hand-edit)
pkg/errors/                     Canonical public typed errors, codes, and sentinels
pkg/observability/              Public OpenTelemetry helpers and instruments
pkg/resilience/                 Public retry, rate-limit, circuit-breaker, and clock primitives
pkg/transport/                  Public HTTP and MQTT transport foundations
pkg/types/                      Shared public market and instrument types
pkg/domain/money/               Public money.Money decimal type
pkg/domain/order/               Public order state machine and reconciliation model
internal/...                    Authentication and region implementation details
proto/, buf.gen.yaml, buf.yaml  Protobuf sources and codegen configuration
examples/                       Runnable main programs and nested example modules
tools/webull-docgen/            Doc generator for docs/webull-api/** and docs/reconciliation.md
docs/                           Documentation site sources (MkDocs Material)
docs/adr/                       Architecture Decision Records
docs/runs/                      Internal run artifacts; excluded from the built site
```

Root files: `doc.go` (module package doc), `mkdocs.yml`, `requirements-docs.txt`,
`Makefile`, `README.md`, `CHANGELOG.md`, `CONTRIBUTING.md`, `SECURITY.md`,
`CODE_OF_CONDUCT.md`, `THIRD_PARTY_NOTICES.md`, `IMPLEMENTATION_STATUS.md`,
`AGENTS.md`, `LICENSE`.

## Official Webull documentation

Canonical links for API paths, schemas, and guides. The official `path` in the
OpenAPI JSON (embedded in the `.md` variant of each reference page) is
authoritative; regenerate derived docs with
`python tools/webull-docgen/docgen.py all`.

**Machine-readable indexes** (fetched by `tools/webull-docgen`)

- HK: <https://developer.webull.hk/apis/llms.txt>
- US: <https://developer.webull.com/apis/llms.txt>
- AI-friendly resources: <https://developer.webull.hk/apis/docs/AI-friendly-Resources/llm.md>

**Docs roots**

- HK: <https://developer.webull.hk/apis/docs>
- US: <https://developer.webull.com/apis/docs>

**Per-endpoint reference bases** (the `// Reference:` comments in Go code)

- HK: <https://developer.webull.hk/apis/docs/reference/>
- US: <https://developer.webull.com/apis/docs/reference/>

**Guides cited in code and docs**

- Getting started: <https://developer.webull.hk/apis/docs/getting-started>
- Display Solution: <https://developer.webull.hk/apis/docs/market-data-api/Hosted-Display-Solution>
- Connect API: <https://developer.webull.com/apis/docs/connect-api/about-connect-api>
- Test accounts: <https://developer.webull.hk/apis/docs/sdk#test-accounts> (HK),
  <https://developer.webull.com/apis/docs/sdk#test-accounts> (US)

**Related**

- Official Python SDK (see `THIRD_PARTY_NOTICES.md`):
  <https://github.com/webull-inc/webull-openapi-python-sdk>
- Project docs site: <https://shing1211.github.io/webullapi4go/>
- Go package docs: <https://pkg.go.dev/github.com/shing1211/webullapi4go>
  (per-package table in `docs/api.md`)

## Build, test, and lint

The Makefile is module-aware: its targets run the root module and every nested
module listed in `MODULES`. Root `go build ./...` and `go test ./...` do not
traverse nested modules.

```sh
make build
make vet
make test
make test-race
make cover
make lint
gofmt -l .                 # must print nothing
make docs
```

- `make build`, `make vet`, and `make test` must pass without credentials.
- `make test-race` runs the race detector in every module.
- `make cover` produces per-module measurements, not a guarantee. Record the
  command/date and do not present an aggregate percentage as proof of behavior.
- CI runs root race tests across three operating systems and build/vet/race for
  nested modules. Its 60% coverage gate is root-only; it does not enforce nested
  coverage or a strict documentation build, so the Makefile gates remain the
  local/release verification source.
- `gofmt -l .` must be empty. Generated code under `gen/` is excluded from the CI
  gofmt check and from golangci-lint (see `.golangci.yml`).
- Unit tests are offline and credential-free. They must not require network
  access.
- On Windows, if a Go build fails with a file-lock error on `a.out.exe`, set
  `GOTMPDIR` to a writable, non-scanned directory and retry; this is a host-level
  antivirus/indexing issue, not a code problem.

Docs:

```sh
python -m venv .venv
.venv/Scripts/pip install -r requirements-docs.txt   # Windows
pip install -r requirements-docs.txt                 # macOS / Linux
mkdocs build --strict
```

The strict build must succeed. Confirm `docs/runs/` is absent from the generated
`site/`.

Protobuf regeneration (not required for builds; generated code is committed):

```sh
make proto-tools
make generate
```

## Conventions

- **License headers.** Every hand-written `.go` file starts with the Apache-2.0
  header (Copyright 2026 shing1211). Copy the header from a neighboring file.
  Generated `*.pb.go` files carry their own generator header.
- **GoDoc.** Every exported identifier has a GoDoc comment that starts with the
  identifier name and explains behavior, not just restates the signature. Flag
  edge cases and defaults.
- **No internal types in public signatures.** Exported types and function
  signatures in the public service packages, `gen/...`, and `pkg/...` must use
  public types only. When a public mirror of an internal type is needed, define
  it in the public package (see `client.Region` and `client.Endpoints`) and do
  not expose `internal/*`.
- **Errors.** `pkg/errors` (imported as `errs`) is the canonical typed-error
  package. Use `errs.New(code, msg)` and `errs.Wrap(code, msg, cause)` for
  category errors; use a package-level `errs.NewSentinel` value only when one
  category contains an identity-specific meaning. `errs.Is(err, code)` is
  category matching, while semantic sentinels match only themselves or wrappers
  preserving their identity. Never match on error strings in library code. Wrap
  underlying causes so `errors.Is`/`errors.As` traversal keeps working. The
  former `internal/errs` compatibility shim was removed as dead code; an
  `internal/` path is unreachable from outside the module, so it could never
  have served that purpose.
- **Options.** Follow the package's functional-option convention; new options
  get a `WithX` constructor with a GoDoc comment. Options are applied in order
  on top of that package's defaults. Constructors that return an error validate
  the resolved configuration; other packages may validate at operation time.
- **Numeric precision.** Webull prices, sizes, balances, and other decimal
  financial values are JSON strings on the wire. Public DTOs expose
  `money.Money` for required/response values and `*money.Money` for
  optional/request values; raw `decimal.Decimal` is not a DTO type.
- **Comments.** Do not add comments that merely restate the code. Comment
  intent, invariants, and non-obvious decisions. The project deliberately omits
  inline comments where the code is self-explanatory.
- **Formatting.** `gofmt`/`goimports` with `local-prefixes:
  github.com/shing1211/webullapi4go`, so the local module imports form the last
  import group. Run gofmt before finishing.
- **Testing.** Add `_test.go` coverage for behavior changes. Prefer table-driven
  tests. Integration tests that touch the network live behind build-independent
  env gates and must be skipped by default.

## Sandbox and secrets

- Credentials are never committed. Read them from the environment:
  `WEBULL_APP_KEY`, `WEBULL_APP_SECRET`, `WEBULL_REGION`,
  `WEBULL_ENVIRONMENT`, `WEBULL_BASE_URL`, `WEBULL_MQTT_URL`.
  **Exception:** the shared public test accounts published by Webull may be
  *linked* from documentation
  (<https://developer.webull.hk/apis/docs/sdk#test-accounts> for HK,
  <https://developer.webull.com/apis/docs/sdk#test-accounts> for US), but their
  account IDs, app keys, and app secrets must not be inlined into hand-written
  files. Webull can rotate or retire them, so a committed copy goes stale
  silently. The generated pages under `docs/webull-api/**` are exempt because
  they are a verbatim snapshot of Webull's own published material; do not edit
  them by hand to remove it.
- The only sandbox host that may appear in committed material is
  `api.sandbox.webull.hk` (and the corresponding `data-api.sandbox.webull.hk`
  MQTT hosts). App keys, app secrets, and access tokens are per-account secrets.
- Core, data, stream, and events sandbox tests are gated by `WEBULL_SANDBOX=1`
  plus `WEBULL_APP_KEY` and `WEBULL_APP_SECRET`. `WEBULL_BASE_URL` may select a
  different regional sandbox; MQTT-over-WebSocket tests also need
  `WEBULL_MQTT_WEBSOCKET=1`.
- Trading-package sandbox tests use a separate gate: `WEBULL_TRADE_SANDBOX=1`,
  `WEBULL_TRADE_APP_KEY`, `WEBULL_TRADE_APP_SECRET`, and
  `WEBULL_TRADE_ACCOUNT_ID`. The mutating trade test additionally requires
  `WEBULL_TRADE_MUTATE=1`; the HK BCAN preview additionally reads
  `WEBULL_TRADE_PARTY_ID`.
- The mutating order-event test uses the generic `WEBULL_SANDBOX` credentials
  plus `WEBULL_TRADE_ACCOUNT_ID` and `WEBULL_TRADE_MUTATE=1`.
- `.env` is gitignored. Never paste private credentials into issues, PRs, docs,
  or tests.

## Known constraints

- **Release status:** current request, OMS, streaming, event-telemetry, and
  documentation hardening is tagged in repository `v2.1.4` (2026-09-26), an
  authorized repository Git patch release; the hardening itself was introduced
  in `v2.1.1` (2026-09-25). The root module path remains
  `github.com/shing1211/webullapi4go`; the module stays on the v1 import path
  **by decision** and no `/v2` migration is planned. Go-semver-compatible v2
  module publication is therefore not pending: it is declined. Because the
  import path carries no major-version suffix, the module proxy serves only the
  `v1.x` line, `v2.x` Git tags are not installable with `go get`, and `v1.1.1`
  is the newest installable version. Work after that tag is consumed by pinning
  a commit, for example `go get github.com/shing1211/webullapi4go@78c164c`.
  Treat any future claim that a `v2.x` tag can be installed as inaccurate.
- A `/v2` migration is a breaking change that alters every consumer import path
  and the meaning of the historical tags. Do not migrate the module path
  implicitly; it requires explicit maintainer approval and a documented
  migration. The current hardening was not newly live-verified.
- **v1.1.0 brought full documented-endpoint coverage**: every documented
  endpoint is implemented. The generated 2026-09-26 reconciliation snapshot
  reports 209 implemented endpoints, 0 documented-only gaps, and the partition
  184 exact OpenAPI JSON path matches, 4 summary-only matches, 1 path differing
  from both sources, 0 unresolved SDK paths, 3 rows carrying the
  `no OpenAPI schema on page` label, and 17 manifest entries deliberately
  mapped to no SDK symbol. The 189 rows with a verified path plus those 3 and 17
  account for all 209, so no endpoint is missing; the previously quoted "25
  unresolved" was 20 generator artifacts and 5 individually investigated
  entries. Do not describe that snapshot as a zero-discrepancy report: 4
  summary-only and 1 differing remain. Earlier version history lives in
  `CHANGELOG.md`.
- **That label count is not a page count.** 7 gRPC reference pages embed no
  OpenAPI schema at all, but only 3 rows carry the `no OpenAPI schema on page`
  label. `_reconcile_data()` in `tools/webull-docgen/docgen.py` tests the
  `unmapped` branch before the `no-openapi` branch, so a page that is both
  schema-less and unmapped — 4 of the 7 are — is counted only as unmapped;
  that precedence is deliberate, because the status table partitions the 209
  rows and must not double-count. A further 3 rows have a JSON block that yields
  no `path`, and all 3 are labelled unmapped, so 10 rows in total have no
  usable official path. Quote 3 as a label count and 7 as the page count.
- Four live-blocked SDK defects found by static analysis on 2026-09-26 are
  recorded with `file:line`, impact, minimal fix, and unblock requirement in
  `IMPLEMENTATION_STATUS.md` and `docs/implementation-status.md`:
  `brokerfd/client.go:43` routes the whole package to the core host instead of
  the Broker host, `brokerfd` still uses 14 undocumented `/broker-fd/*` path
  literals (`brokerfd/assets.go:25` is the one the generated report flags),
  `broker.UpdateVirtualAccount` sends the wrong verb and body
  (`broker/accounts.go:66-68`), and `data.GetDisplaySnapshot` differs from both
  official sources (`data/display_quotes.go:28`, `:52`). None is live-verified
  and none may be changed without the credentials or entitlement each entry
  names.
- US-only surfaces are blocked in this environment: the HK sandbox returns `404`
  (fund data, crypto data, screener v2, broker FD, instrument v3/logos) or `417`
  (crypto category), and no US sandbox credentials are available. Those items stay
  unverified until `WEBULL_APP_KEY` and `WEBULL_APP_SECRET` for the US sandbox are
  supplied.
- HTTP `417` maps to `errs.CodeInvalidToken` for compatibility, but Webull also
  uses it for business validation such as invalid symbols, unsupported
  categories, and rejected strategies. Preserve status/message for diagnostics;
  never infer every 417 is a token failure or match the message string.
- Display Solution requires a paid Webull subscription; the HK sandbox host
  (`hk-co-branding-openapi.uat.webullbroker.com`) returns `403 Forbidden` at the
  host level, blocking all Display Solution endpoints even with valid credentials.
- Broker API HK (`/broker/...`) returns `401 ROUTE_NOT_PERMITTED` in the HK
  sandbox — the app lacks the required scope, not a path issue. Broker HK remains
  unverified pending production or US sandbox access.
- SSE news upstream returns `504 Gateway Timeout` in the HK sandbox.
- Sandbox market data is limited to `AAPL`.
- Footprint requires a paid entitlement; the sandbox returns `403 Insufficient
  permission`.
- Option contracts for `AAPL` may not exist in the sandbox (`417 Invalid
  Symbol`).
- Order-book depth can be empty outside regular trading hours.
- Plain MQTT on port `1883` can be blocked; use MQTT over WebSocket on
  `wss://...:8883/mqtt`.
- The token endpoint allows 10 requests per 30 seconds, and MQTT allows at most
  5 concurrent connections per App Key.
- `docs/runs/**` is excluded from the published docs site. Do not edit
  `docs/runs/**` unless the task explicitly authorizes a release/status artifact
  update. Accepted ADRs (0001, 0002) are immutable; supersede them with a new ADR
  instead.
- Generated docs (`docs/webull-api/**` and `docs/reconciliation.md`) start with
  a "Generated file — do not edit" banner. Change `tools/webull-docgen/` (or its
  manifest `_common.py`) and regenerate instead of editing them by hand.
- Do not edit `.github/workflows/ci.yml`, `.golangci.yml`, or `LICENSE` without
  an explicit request.

## Commit style

Conventional Commits with a `Signed-off-by` trailer (`git commit -s`). Keep each
commit focused. Do not commit unless explicitly asked.
