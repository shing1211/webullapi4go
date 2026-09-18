# Todos — webull-sdk-bootstrap

Single source of truth. Statuses: `todo`, `doing`, `blocked`, `review`, `done`.

| ID | Task | Role | Status | Depends On | Acceptance |
|----|------|------|--------|------------|------------|
| T0.1 | Scaffold module + package layout + doc.go | architect | done | — | build+vet+gofmt clean; fixed to `.webull.hk` region hosts |
| T0.2 | GitHub Actions CI (build/vet/test/lint) | devops | done | T0.1 | `.github/workflows/ci.yml` + `.golangci.yml`; lint 0 issues |
| T0.3 | Community files + Discussions | docs | done | — | 6 files added; Discussions enabled (`has_discussions=true`) |
| T0.4 | MkDocs Material docs site + Pages | devops | done | T0.1 | `mkdocs build --strict` OK; runs/ excluded |
| T1.1 | HMAC-SHA1 signer + golden test | security | done | T0.1 | `kvlS6opdZDhEBo5jq40nHYXaLvM=` matches; 9 tests pass |
| T1.2 | Token lifecycle (create/check/poll/store) | backend | done | T0.1, T2.1 | sandbox token `NORMAL`; `client/token.go` + RoundTripper injection seam |
| T1.3 | Credentials (env + functional options) | backend | done | T2.1, T2.2 | `WithEnv()` + 9 tests; explicit options > env |
| T2.1 | Core HTTP transport + typed errors + public `client` facade | backend | done | T0.1, T1.1 | live sandbox `account/list` OK; httptest sign/verify |
| T2.2 | Retry / rate limiter / circuit breaker | backend | done | T2.1 | `internal/resilience/*`; retry default GET; lint 0 |
| T2.4 | API version (v2/v3) + auto token attach + live-verify T4.1 endpoints | backend | done | T1.2, T2.2 | v3 paths OK; `WithAPIVersion`/`WithAutoToken`; live-verified |
| T2.5 | Fix per-client token state (remove `sync.Map` keyed by `*Client`) | backend | doing | T1.2, T2.4 | full sandbox suite no longer 401s |
| T2.3 | Multi-region / endpoint config | architect | done | T0.1 | region-switch test passes (`internal/region`) |
| T3.1 | Proto spike: MQTT Quote/Snapshot/Tick → Go | data | done | T0.1 | buf 1.73/protoc 29.2; generated code + round-trip tests pass |
| T3.2 | gRPC proto availability ADR | architect | done | — | ADRs 0001/0002; recommend vendor from Apache-2.0 Python SDK |
| T4.1 | Market Data HTTP: instrument/profile/analyst | backend | done | T2.1, T2.3 | 7 endpoints; v3 paths verified vs docs; 8 unit tests; AAPL live |
| T4.2 | Market Data HTTP: snapshot/tick/quotes/bars | backend | done | T2.1, T2.3 | v3 paths; AAPL live (snapshot/tick/bars); quotes depth empty off-RTH |
| T4.3 | Market Data HTTP: footprint/NOII/screener | backend | done | T2.1, T2.3 | live: gainers/losers, most-active, NOII; footprint 403 (entitlement) |
| T4.4 | Market Data HTTP: watchlist/options/news | backend | done | T2.1, T2.3 | watchlist live; options 417 (sandbox symbol); news SSE (bypasses client.Do — debt) |
| T5.1 | MQTT: connect/auth/session + HTTP subscribe + protobuf parsers (public `stream`) | backend | done | T0.1, T3.1 | paho v1.5.1; WebSocket sandbox live PASS (TCP 1883 blocked here) |
| T5.2 | MQTT reconnect + re-subscribe + connection limits | backend | done | T5.1 | registry + resubscribe; live reconnect PASS |
| T5.3 | MQTT live E2E hardening | backend | done | T5.2 | absorbed into T5.2 (live reconnect verified) |
| T6.1 | Examples + README + docs sync | docs | done | T1–T5 | 4 examples; README/CHANGELOG/AGENTS/docs site; mkdocs --strict OK |
| T6.1b | Final GoDoc/comment polish + consistency review | reviewer | doing | T6.1 | no stale "in progress"/placeholder comments |
| T6.2 | Release v0.1.0 (GitHub + Gitee) | release | todo | T6.1 | both remotes show commit/tag |
