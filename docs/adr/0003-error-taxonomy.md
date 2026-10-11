---
Status: Accepted
Date: 2026-10-11
Deciders: maintainer
---

# 0003. Error taxonomy

## Context

`webullapi4go` wraps the Webull OpenAPI. API calls can fail in several structurally distinct ways:
HTTP-level errors (401, 403, 429, 5xx), network failures, protocol-level errors (MQTT disconnect, gRPC status errors), and business-logic errors encoded in the response body.

Callers need to handle failures without matching on error strings. A stable, typed
error taxonomy is required for production-grade error handling.

## Decision

All SDK errors implement the shared [`Error`](https://pkg.go.dev/github.com/shing1211/webullapi4go/pkg/errors#Error) interface defined in `pkg/errors/errors.go`.

An `Error` carries:

- A stable [`Code`](https://pkg.go.dev/github.com/shing1211/webullapi4go/pkg/errors#Code) — a `string` machine-readable classification.
- A human-readable `Message`.
- An optional wrapped `Cause` (`error`).

`Error` participates in `errors.Is` and `errors.As` through `Unwrap`, so callers
can use `errors.Is(err, errs.CodeUnauthorized)` without string matching.

### Error codes

| Code | Meaning | HTTP analogue |
|------|---------|---------------|
| `invalid_config` | Client configuration is invalid | — |
| `unsupported` | Operation not supported by the SDK | — |
| `auth` | Authentication, signing, or token failure | 401 |
| `transport` | Network or HTTP transport failure | — |
| `api` | Non-success response from Webull API | 4xx / 5xx |
| `unauthorized` | HTTP 401 — not authenticated | 401 |
| `forbidden` | HTTP 403 — lacks permission or subscription | 403 |
| `not_found` | HTTP 404 — resource not found | 404 |
| `rate_limited` | HTTP 429 — rate limit exceeded | 429 |
| `internal` | HTTP 5xx — Webull server error | 5xx |
| `unavailable` | Service temporarily unavailable | 503 |

### Wrapping

SDK packages may wrap SDK errors with additional context (e.g. endpoint path,
request ID). The wrapped error is retrievable via `errors.Unwrap()`.

Third-party errors (network timeouts, TLS handshake failures) are wrapped with
`CodeTransport`. They are not re-typed to `Error`; callers checking
`errors.Is(err, context.DeadlineExceeded)` still works.

## Rationale

A single `Error` type with a `Code` field gives callers a stable, branchable
surface without requiring `errors.As` type assertions. The `Code` constants are
string literals (not integers) so they serialize predictably to logs and
observability systems.

Wrapping third-party errors with `CodeTransport` rather than re-typing them
preserves the original error chain and does not lose `errors.Is` behavior for
callers who check context or stdlib errors.

## Consequences

- Callers can classify all SDK errors by inspecting `Code` without string matching.
- `errors.Is` and `errors.As` work correctly on all SDK errors.
- Error messages are human-readable but not machine-contractual; callers must not
  parse `Message` strings.
- The code taxonomy is frozen; adding a new `Code` is an ADR.
