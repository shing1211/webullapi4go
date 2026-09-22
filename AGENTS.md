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
broker/                         Broker API HK (own Go module; root module via `replace`)
brokerfd/                       Broker FD US HTTP endpoints
brokerfd/events/                Broker FD events over gRPC + protobuf payloads
display/                        Display Solution client-token authentication
connect/                        Connect API OAuth 2.0 authorization-code flow
gen/webull/...                  Generated protobuf types (committed; do not hand-edit)
pkg/types/                      Shared public domain types
internal/...                    auth, errs, region, transport, resilience, mqtt
proto/, buf.gen.yaml, buf.yaml  Protobuf sources and codegen configuration
examples/                       Runnable main programs
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

```sh
go build ./...
go vet ./...
gofmt -l .                 # must print nothing
go test ./...
go test -race -count=1 ./...
golangci-lint run ./...
```

- `go build ./...` and `go vet ./...` must pass without credentials.
- `gofmt -l .` must be empty. Generated code under `gen/` is excluded from the CI
  gofmt check and from golangci-lint (see `.golangci.yml`).
- Unit tests are offline and credential-free. They must not require network
  access.
- On Windows, if `go build ./...` fails with a file-lock error on
  `a.out.exe`, set `GOTMPDIR` to a writable, non-scanned directory and retry;
  this is a host-level antivirus/indexing issue, not a code problem.

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
  signatures in `client`, `data`, `stream`, `gen/...`, and `pkg/types` must use
  public types only. When a public mirror of an internal type is needed, define
  it in the public package (see `client.Region` and `client.Endpoints`) and do
  not expose `internal/*`.
- **Errors.** Use `internal/errs` for typed errors: `errs.New(code, msg)` and
  `errs.Wrap(code, msg, cause)`. Never match on error strings in library code.
  Wrap underlying causes so `errors.Is`/`errors.As` traversal keeps working.
- **Options.** Follow the `Option func(*Config)` pattern; new options get a
  `WithX` constructor with a GoDoc comment. Options are applied in order on top
  of the defaults.
- **Numeric precision.** Market-data prices and sizes are strings on the wire;
  keep them as strings in DTOs.
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
  **Exception:** Shared public test accounts published by Webull on their
  official docs (<https://developer.webull.hk/apis/docs/sdk#test-accounts>)
  may appear in documentation and test fixtures.
- The only sandbox host that may appear in committed material is
  `api.sandbox.webull.hk` (and the corresponding `data-api.sandbox.webull.hk`
  MQTT hosts). App keys, app secrets, and access tokens are per-account secrets.
- Sandbox integration tests are gated by `WEBULL_SANDBOX=1` plus
  `WEBULL_APP_KEY` and `WEBULL_APP_SECRET`; MQTT-over-WebSocket tests also need
  `WEBULL_MQTT_WEBSOCKET=1`.
- `.env` is gitignored. Never paste private credentials into issues, PRs, docs,
  or tests.

## Known constraints

- **v1.1.0 brings full coverage of the Webull OpenAPI**: every documented
  endpoint is implemented and all paths follow the official OpenAPI definition
  (`docs/reconciliation.md`: 209 implemented, 0 gaps, 0 path discrepancies).
  Earlier version history lives in `CHANGELOG.md`.
- US-only surfaces are blocked in this environment: the HK sandbox returns `404`
  (fund data, crypto data, screener v2, broker FD, instrument v3/logos) or `417`
  (crypto category), and no US sandbox credentials are available. Those items stay
  unverified until `WEBULL_APP_KEY` and `WEBULL_APP_SECRET` for the US sandbox are
  supplied.
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
  `docs/runs/**`. Accepted ADRs (0001, 0002) are immutable; supersede them with
  a new ADR instead.
- Generated docs (`docs/webull-api/**` and `docs/reconciliation.md`) start with
  a "Generated file — do not edit" banner. Change `tools/webull-docgen/` (or its
  manifest `_common.py`) and regenerate instead of editing them by hand.
- Do not edit `.github/workflows/ci.yml`, `.golangci.yml`, or `LICENSE` without
  an explicit request.

## Commit style

Conventional Commits with a `Signed-off-by` trailer (`git commit -s`). Keep each
commit focused. Do not commit unless explicitly asked.
