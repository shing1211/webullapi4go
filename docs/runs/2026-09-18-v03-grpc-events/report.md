# Report — v03-grpc-events (v0.3.0)

## Outcome
v0.3.0 is released on both GitHub and Gitee: real-time gRPC trade events (order/position/option)
with signed subscriptions, typed JSON payloads, and reconnect/re-subscribe.

- Release commit `c46884f`; run-record `99d505a`; annotated tag **`v0.3.0`** (both remotes at `99d505a`).

## Shipped
- **Vendored proto**: `proto/webull/trade/events/v1/events.proto`, byte-identical to upstream
  (git blob `7ebc4728…`), Apache-2.0 with `THIRD_PARTY_NOTICES.md`; Go types + gRPC stub generated.
- **Signer**: `internal/auth` algorithm abstraction (HMAC-SHA1 default, HMAC-SHA256) with a
  SHA256 golden vector; REST behaviour unchanged.
- **`events` package**: TLS gRPC connect to `events-api.<domain>:443`, HMAC-SHA256 signed metadata
  (no `host`), `SubscribeRequest` with subscribe bitmask/accounts/timestamp, dispatch of
  `SubscribeSuccess`/`Ping`/`AuthError`/`NumOfConnExceed`/`SubscribeExpired`, typed
  `OrderEvent`/`PositionEvent`/`OptionEvent` handlers, reconnect with exponential backoff + jitter.
- **Docs/example**: `docs/events.md`, `examples/events`, README/CHANGELOG/docs-site updates.
- **CI**: `.github/workflows/nightly-live.yml` — scheduled read-only/preview live-sandbox job
  (mutation excluded; skips cleanly when secrets are absent).
- **ADR-0002** promoted to **Accepted** (Option A confirmed: real `.proto` sources exist).

## Verification / live evidence
- `go build`/`vet`/`gofmt`/`golangci-lint` clean; `go test ./... -count=1` and `-race` green.
- Sandbox: `SubscribeSuccess` (matched App Key #2 ↔ Account #2) and a live `OrderEvent`
  decoded (`status=CANCELLED scene=CANCEL_SUCCESS symbol=AAPL`).

## Deviations / notes (important)
1. **Events digest is lowercase SHA-256**, not uppercase. The task brief said uppercase (mirroring
   the REST signer); the sandbox rejected that. `events/sign.go` implements the lowercase events
   canonical string locally (pinned by a Python-derived golden test).
2. Consequently, the `internal/auth` HMAC-SHA256 path (uppercase) added in T15.2 is **currently
   unused**. Reconcile: either add a lowercase mode to `internal/auth` and delete `events/sign.go`,
   or keep the split and document it. Tracked as tech debt.
3. **App key ↔ account pairing matters**: Account #2 (`OGG4…`) is owned by App Key #2, not #1;
   using the wrong pair returns `PermissionDenied unauthorized account`.
4. Sandbox did **not** push placement events for a resting limit order — only `CANCEL_SUCCESS` was
   observed; other `scene_type` values are covered by fixtures, not live.
5. Accounts are effectively required by the sandbox subscribe (`InvalidArgument` when empty),
   despite being optional in the schema.
6. The nightly workflow needs repository secrets configured (list in the workflow header) or it skips.

## Risks / follow-ups
- Position/option payload schemas were derived from docs/analogy, not a captured live sample.
- Nightly live CI depends on shared sandbox stability (documented flakiness on watchlist 401).
- Display Solution and Broker API remain unimplemented.
