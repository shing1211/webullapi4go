# Plan — webullapi4go v0.3 (gRPC trade events)

## Run identity
- Date: 2026-09-18
- Slug: v03-grpc-events
- Mode: BUILD
- Repo: `D:\github\webullapi4go`
- Base: v0.2.6 (`df6040d` + docs `12d87f8`)
- Module: `github.com/shing1211/webullapi4go`
- Remotes: GitHub `origin`, Gitee `gitee`

## Goal
Subscribe to Webull's gRPC order/position/option event stream and expose it as a public
`events` package, verified against the sandbox, released as **v0.3.0**.

## Recon findings (verified this run)
- Python SDK ships a real `.proto`: `webull/trade/events/events.proto` (Apache-2.0, `LICENSE` + `NOTICE`, `Copyright 2022 Webull`).
- Service: `grpc.trade.event.EventService/Subscribe` (server-streaming).
- `SubscribeRequest{subscribeType, timestamp(millis), contentType, payload, accounts[]}`.
- `SubscribeResponse{eventType, subscribeType, contentType, payload, requestId, timestamp}`.
- `EventType`: SubscribeSuccess=0, Ping=1, AuthError=2, NumOfConnExceed=3, SubscribeExpired=4.
- Data `subscribeType`: order=1024, position=1028, option=1032; bitmask request `7` = all.
- Payload is **JSON** when `contentType=application/json`; order status change types 1/2/4.
- Endpoint `events-api.sandbox.webull.hk:443` is reachable (TLS gRPC).
- Events auth uses **HMAC-SHA256** (algorithm `HMAC-SHA256`, version `1.0`), body hash `UPPER(SHA256(proto_bytes))`, signature `base64(HMAC-SHA256(app_secret+"&", quote(string_to_sign)))`.

## Decisions (locked)
1. ADR-0002 **Option A** confirmed viable (real proto sources) — vendor with attribution.
2. Add dependency `google.golang.org/grpc`.
3. Include a **nightly read-only/preview** live-sandbox CI job (mutation off).

## Approach (alternatives)
1. **Vendor the published `events.proto` + gRPC** (chosen) — highest fidelity, small surface, Apache-2.0 clean.
2. Hand-author protos from wire bytes — error-prone, rejected.
3. Defer to hosted Cloud MCP — gives up a planned capability, rejected.

## Architecture
- `proto/webull/trade/events/v1/events.proto` (vendored, unmodified) + `NOTICE`/`THIRD_PARTY_NOTICES`.
- Generated Go in `gen/webull/trade/events/v1/` (non-public path is fine; keep under `gen/`).
- `internal/auth`: add an algorithm-parameterised signer (HMAC-SHA1 default, HMAC-SHA256 for gRPC) with SHA256/MD5 body hashing; keep existing API working.
- Public `events` package: `events.New(cl *client.Client, opts...)`, `Connect(ctx)`, `Subscribe(ctx, accounts, types)`, handlers (`OnOrder`, `OnPosition`, `OnOption`, `OnError`, `OnConnect`), `Close`.
- Reconnect/resubscribe mirroring the Python SDK's retry policy.

## Tasks
| ID | Objective | Role | Depends | Acceptance | Size |
|----|-----------|------|---------|-----------|------|
| T15.1 | Vendor `events.proto`, `NOTICE`/attribution, generate Go; promote ADR-0002 to Accepted | data | — | proto compiles; decode sample; attribution present | M |
| T15.2 | Algorithm-parameterised signer (`HMAC-SHA256` + SHA256 body) with golden tests | security | — | SHA1 unchanged; SHA256 vectors pass | M |
| T15.3 | `events` package: TLS gRPC dial, signed metadata, `Subscribe`, stream dispatch, Ping/error handling | backend | T15.1, T15.2 | sandbox `SubscribeSuccess` | L |
| T15.4 | Typed order/position/option JSON payloads + reconnect/resubscribe | backend | T15.3 | live order event (mutate-gated) | L |
| T15.5 | Tests + `docs/events.md` + example | tester/docs | T15.4 | documented commands run | M |
| T15.6 | Nightly read-only live-sandbox CI job | devops | T15.3 | workflow valid; mutation off | S |
| T15.7 | Release v0.3.0 | release | T15.5, T15.6 | tag on both remotes | S |

## Test & credential strategy
- Unit tests offline (httptest / bufconn gRPC).
- Live gated: `WEBULL_SANDBOX=1`, `WEBULL_APP_KEY`/`WEBULL_APP_SECRET`; event subscription uses dedicated trading Account #2 env (`WEBULL_TRADE_*`) for mutation-gated event triggering.
- CI nightly job: read-only + preview only; no mutation; uses repository secrets.

## Risks & mitigations
1. **Unknown JSON payload schemas** → derive from a captured sandbox event in T15.4; keep decoders tolerant of unknown fields.
2. **Signer regression** → T15.2 golden tests must prove SHA1 path unchanged.
3. **Sandbox event silence** → trigger an event via a tiny non-marketable order under the mutate gate.
4. **gRPC/reconnect complexity** → mirror the Python retry policy; unit-test with bufconn.

## Verification
- `go build ./...`, `go vet ./...`, `gofmt -l .`, `golangci-lint run ./...` clean.
- `go test ./... -count=1` offline passes; gated live tests pass when configured.
- `mkdocs build --strict` succeeds.

## Order of work
T15.1 ‖ T15.2 → T15.3 → T15.4 → T15.5 ‖ T15.6 → T15.7.
