# Report — webull-sdk-bootstrap (v0.1.0)

## Outcome
v0.1.0 of **webullapi4go** is implemented and verified: authentication, core HTTP client,
market-data HTTP, and MQTT streaming, with examples, docs site, CI, and community files.

## Shipped (vs plan)
All v0.1 tasks (T0.1–T6.1b) completed. Highlights:

- **Module**: `github.com/shing1211/webullapi4go`, Go 1.26, Apache-2.0, package layout
  (`client`, `data`, `stream`, `pkg/types`, `internal/*`, `gen/...`).
- **Auth**: HMAC-SHA1 signer (documented golden vector `kvlS6opdZDhEBo5jq40nHYXaLvM=`),
  token create/check/poll/cache (`EnsureToken`), env credentials (`WithEnv`).
- **Core client**: `client.New`, `Do`, typed errors, multi-region endpoints (12 regions),
  configurable API version (`v2`/`v3`), auto-token, retry/rate-limit/breaker.
- **Market Data HTTP**: instruments/profile/analyst, futures static, snapshot, tick, quotes,
  bars (single+batch), footprint, NOII, screener, watchlist CRUD, options, news.
- **Streaming**: MQTT (paho) connect over TCP/WebSocket, HTTP subscribe/unsubscribe,
  protobuf parsing, reconnect + re-subscribe registry, connection-limit handling.
- **Quality**: CI (`ci.yml`, `golangci-lint`), 60+ tests, docs site (MkDocs Material,
  `--strict`), 4 runnable examples, README/CHANGELOG/AGENTS/CONTRIBUTING/SECURITY/CoC, ADRs.

## Verification evidence
- `go build ./...`, `go vet ./...`, `gofmt -l .`, `golangci-lint run ./...` all clean.
- `go test ./... -count=1` passes.
- Live sandbox (shared test accounts) verified: `account/list` 200, token `NORMAL`,
  `AAPL` snapshot/tick/bars/instrument/screener/NOII, MQTT-over-WebSocket stream +
  forced reconnect with re-subscribe (2 subscribe calls).

## Deviations from the original plan
1. **Documented v3 paths**: market-data endpoints use current `/market-data/...` and
   `/trading/instruments/...` paths (the `/openapi/...` aliases in llms.txt are legacy).
2. **`x-version` default remains `v2`**, configurable via `WithAPIVersion`; live probes
   confirmed both v2 and v3 work.
3. **News** is an SSE stream and bypasses `client.Do` (own signing path) — token/version/
   resilience are not applied. Recorded as technical debt (see next-phase).
4. **MQTT plain TCP 1883** is egress-blocked from the build network; WebSocket 8883 is
   verified and recommended in docs.

## Sandbox limitations observed
- Market data limited to `AAPL`; footprint `403` (entitlement); AAPL option symbol `417`;
  depth may be empty outside regular hours; token endpoint 10 req/30s; 5 MQTT connections.

## Risks / follow-ups
- Residual flakiness of shared sandbox accounts (401s under rapid token creation) —
  test-only; SDK-side per-client token state was fixed in T2.5.
- `news` SSE duplication of the signing/transport path.
- gRPC event protos remain unpublished (ADR-0002 proposes extracting from the Apache-2.0
  Python SDK).
- Coverage: trading, gRPC events, Display Solution, Broker API are deferred to v0.2+.

## Post-fix note
During the run, T2.5 fixed a real bug: per-client token state stored in
package-level `sync.Map`s keyed by `*Client` caused stale state under pointer-address
reuse (intermittent 401). State now lives on the `Client` struct with regression tests.
