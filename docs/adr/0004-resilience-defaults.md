---
Status: Accepted
Date: 2026-10-11
Deciders: maintainer
---

# 0004. Resilience defaults

## Context

SDK users make network calls that can fail transiently. The SDK provides resilience primitives (`pkg/resilience/retry`, `pkg/resilience/breaker`, `pkg/resilience/ratelimit`) but callers must configure them. Without documented defaults, every caller invents their own values, leading to inconsistency and misconfiguration.

## Decision

Each resilience primitive ships with a `DefaultConfig()` constructor and named constants for all tunable parameters. Callers who do not configure a primitive get the documented defaults.

### Retry defaults

| Parameter | Default | Rationale |
|-----------|---------|-----------|
| `MaxAttempts` | 2 | One retry; total 2 attempts |
| `BaseDelay` | 100ms | Minimises unnecessary latency on success-path callers |
| `MaxDelay` | 2s | Caps worst-case wait |
| `Jitter` | true (full jitter) | Reduces thundering herd |
| `IsRetryable` | 5xx + context-deadline + rate-limit | Idempotent-safe only |

Retry is disabled by default in clients (`New` does not wire a retry policy). Callers
who want retry behaviour must explicitly pass a configured `Retry` option.

### Circuit breaker defaults

| Parameter | Default | Rationale |
|-----------|---------|-----------|
| `Threshold` | 5 consecutive failures | Avoids false-positives on single bad days |
| `Cooldown` | 30s | Webull API stability is undocumented; 30s is conservative |
| `HalfOpenMax` | 1 probe | Single probe limits exposure if backend is still unhealthy |

The circuit breaker is disabled by default in clients. Callers who want circuit-breaking must explicitly wire it.

### Rate limiter defaults

No global default rate limit is enforced by SDK clients. Per-client rate limiters are available in `pkg/resilience/ratelimit` for callers who want per-endpoint or per-account budgets.

### Observability

All primitives emit OpenTelemetry metrics (when a `MeterProvider` is supplied):
retry attempts, circuit state transitions, and limiter wait times. See `pkg/observability`.

## Rationale

The defaults are conservative (2 attempts, 30s breaker cooldown) because:
1. Webull API stability is undocumented; aggressive retry could amplify load during incidents.
2. The HK sandbox is used for development; aggressive defaults could trigger rate limits on shared sandbox credentials.

Defaults are opt-in (not automatically wired) so production callers make an explicit
decision about retry/breaker policy rather than silently inheriting SDK defaults.

## Consequences

- All resilience parameters are named constants — callers can reference them without magic numbers.
- Defaults are documented in code and in this ADR; no hidden behaviour.
- Adding or changing defaults is an ADR-level decision.
