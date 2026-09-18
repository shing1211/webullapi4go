# Plan — webullapi4go v0.1 (Foundation + Market Data)

## Run identity
- Date: 2026-09-18
- Slug: webull-sdk-bootstrap
- Mode: BUILD
- Repo: `D:\github\webullapi4go`
- Module: `github.com/shing1211/webullapi4go`
- Go directive: `go 1.26` (local toolchain go1.26.8)
- License: Apache-2.0
- Remotes: GitHub `github.com/shing1211/webullapi4go` + Gitee (exists)
- Main branch: `main`

## Goal
A public, idiomatic Go SDK wrapping the Webull OpenAPI (HK), with GitHub Pages docs and
Discussions. v0.1 ships **Auth + core HTTP client + Market Data (HTTP + MQTT streaming)**,
verified against Webull's shared sandbox test accounts.

## Scope

### In scope (v0.1)
- Module scaffolding, package layout, CI, community files, docs site.
- HMAC-SHA1 request signing (MD5 body, URL-encode, Go HTML-unescape).
- Token lifecycle: create / check / poll / store (sandbox auto-`NORMAL`).
- Core HTTP transport: request builder, signing middleware, compact JSON, typed errors.
- Resilience primitives: retry, rate limiter, circuit breaker.
- Multi-region / multi-endpoint configuration.
- Market Data HTTP: instrument, snapshot, tick, quotes, bars (single/batch), footprint,
  NOII, screener, watchlist, futures, options, news.
- MQTT streaming: connect, subscribe/unsubscribe, protobuf parsers, reconnect + re-subscribe.
- Examples, README, docs sync, v0.1.0 release.

### Out of scope (v0.2+)
- Trading HTTP (accounts, assets, orders, combo orders, options orders).
- gRPC trade events.
- Display Solution API.
- Broker API.
- MCP server surface.

## Approach
Three alternatives considered:

1. **Big-bang** — implement all 105 endpoints, release once. Long feedback loop, high risk. Rejected.
2. **Vertical slice → release v0.1 → expand** — architecture designed up-front, but ship a
   small testable slice early. *Chosen.*
3. **Horizontal infra-first** — all layers before any endpoint. Nothing usable until late;
   its architectural benefit is captured by (2).

**Pick (2)** because it produces a public, sandbox-verified artifact quickly, de-risks the
unknown protobuf work early (Phase 3 spike), and keeps later phases as roadmap items that
are detailed just-in-time via `next-phase.md`.

## Assumptions
- Webull publishes MQTT `Quote`/`Snapshot`/`Tick` `.proto` in docs; gRPC event protos are
  unpublished and must (if needed) be extracted from the Apache-2.0 Python SDK.
- Sandbox test credentials are valid and tokens are auto-`NORMAL` (no 2FA).
- Sandbox market data is limited to `AAPL`.
- Extracted `.proto` files may be vendored with a `NOTICE`/attribution file.

## v0.1 Task breakdown

| ID | Objective | Role | Depends on | Acceptance | Size |
|----|-----------|------|-----------|-----------|------|
| T0.1 | Scaffold: `go.mod` (go 1.26, module path), package layout, `doc.go` | architect | — | `go build ./...` + `go vet ./...` pass | S |
| T0.2 | CI: GitHub Actions build/vet/test + `golangci-lint` | devops | T0.1 | workflow green on PR | S |
| T0.3 | Community: CONTRIBUTING, SECURITY, CODE_OF_CONDUCT; Discussions enabled | docs | — | files present, links valid | S |
| T0.4 | Docs site: MkDocs Material + Pages workflow | devops | T0.1 | Pages deploys | M |
| T1.1 | HMAC-SHA1 signer (MD5 body, URL-encode, Go HTML-unescape) | security | T0.1 | golden test = `kvlS6opdZDhEBo5jq40nHYXaLvM=` | M |
| T1.2 | Token lifecycle (create/check, sandbox auto-NORMAL, poll, store) | backend | T0.1, T2.1 | sandbox token status `NORMAL` | M |
| T1.3 | Credentials: env vars + functional options | backend | T0.1 | unit tests pass | S |
| T2.1 | Core HTTP transport: request builder, signing middleware, compact JSON, typed errors | backend | T0.1, T1.1 | integration vs sandbox `account/list` | L |
| T2.2 | Resilience: retry, rate limiter, circuit breaker (mirror futuapi4go `pkg/*`) | backend | T2.1 | unit tests pass | M |
| T2.3 | Multi-region / endpoint config | architect | T0.1 | region-switch test passes | S |
| T3.1 | **Proto spike:** MQTT `Quote/Snapshot/Tick` → generate Go | data | T0.1 | proto compiles; decode sample | M |
| T3.2 | **gRPC proto availability** ADR (informs v0.2/v0.3) | architect | — | ADR in `docs/adr/` | S |
| T4.1 | Market Data HTTP: instrument + company profile/analyst | backend | T2.1, T2.3 | sandbox test `AAPL` | M |
| T4.2 | Market Data HTTP: snapshot, tick, quotes, bars (single/batch) | backend | T2.1, T2.3 | sandbox test `AAPL` | L |
| T4.3 | Market Data HTTP: footprint, NOII, screener | backend | T2.1, T2.3 | sandbox tests | M |
| T4.4 | Market Data HTTP: watchlist CRUD, futures, options, news | backend | T2.1, T2.3 | sandbox tests | L |
| T5.1 | MQTT client (paho): connect/auth/session management | backend | T0.1, T3.1 | connects to sandbox MQTT | L |
| T5.2 | MQTT subscribe/unsubscribe via HTTP + parsers | backend | T5.1, T4.2 | live `AAPL` quote/snapshot/tick | L |
| T5.3 | Reconnect + re-subscribe + limits (5 conn, session_id rules) | backend | T5.2 | reconnect test passes | M |
| T6.1 | Examples + README + docs sync | docs | T1–T5 | documented commands run | M |
| T6.2 | Release v0.1.0 to GitHub + Gitee `main` | release | T6.1 | both remotes show commit/tag | S |

## Full phase roadmap
| Phase | Theme | Scope | Depends on | Risk |
|-------|-------|-------|-----------|------|
| v0.1 | Foundation + Market Data | Auth, core HTTP, resilience, regions; Market Data HTTP; MQTT | — | proto extraction; MQTT re-subscribe |
| v0.2 | Trading HTTP | Accounts, assets, orders, US/HK/CN rules, options, combo orders | v0.1 | order-rule complexity, BCAN, A-share limits |
| v0.3 | Events (gRPC) | Trade events + order status | v0.1 + T3.2 ADR | gRPC protos unpublished |
| v0.4 | Display Solution | Client token, display quotes, screener/news/corp-actions/instruments | v0.1 | separate C2S auth model |
| v0.5 | Broker API | Virtual accounts, activities, order, funding, journals, master data, events | v0.1, v0.2 | institutional access, 7 event streams |
| v1.0 | Stabilization | API freeze + semver, benchmarks, fuzz, security review, region parity | all | scope creep |
| v1.1 (opt) | MCP Server | Expose SDK as MCP tools | v1.0 | alignment with hosted MCP |

v0.1–v0.5 together cover all **105 endpoints**.

## Risks & mitigations
1. **gRPC event protos unpublished** → T3.2 ADR early; may extract from Python SDK for v0.2+.
2. **Signature JSON-escaping bugs** → golden signature test in T1.1.
3. **MQTT connection limits / no auto re-subscribe** → explicit reconnect logic in T5.3.
4. **Order-rule complexity** → deferred entirely to v0.2 (v0.1 scope contained).
5. **Rate limits** → rate limiter in T2.2 keyed per-endpoint.

## Order of work
T0.1 → (T0.2, T0.3, T0.4, T1.1, T2.3, T3.1, T3.2 in parallel) → T2.1 → (T1.2, T2.2) →
T4.1 → T4.2 → (T4.3, T4.4) → T5.1 → T5.2 → T5.3 → T6.1 → T6.2.

## Verification
- `go build ./...`
- `go vet ./...`
- `go test ./...`
- `golangci-lint run`
- Sandbox integration tests gated by env (`WEBULL_SANDBOX=1`).
