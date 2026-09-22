# Errors

Every SDK function returns a normal Go `error`. Errors produced by the SDK are
typed, carry a stable machine-readable code, and wrap their underlying cause, so
they work with the standard library's `errors.Is` and `errors.As` traversal.

## Error shape

A Webull error formats as:

```
webull: <code>: <message>
webull: <code>: <message>: <cause>
```

For example:

```
webull: FORBIDDEN: http 403: Insufficient permission
webull: transport: GET /market-data/stocks/snapshots/list: dial tcp: ...: connectex: connection refused
```

The `<code>` token is stable and intended for classification. HTTP-derived
errors classify as follows:

| Code | Source |
|------|--------|
| `UNAUTHORIZED` | HTTP 401 |
| `FORBIDDEN` | HTTP 403 (missing permission or data entitlement) |
| `INVALID_TOKEN` | HTTP 417 (missing, expired, or invalid access token) |
| `RATE_LIMITED` | HTTP 429 |
| `SERVER_ERROR` | HTTP 5xx |
| `api` | Any other non-2xx response, or a response that failed to decode |

Errors that do not come from an HTTP response use:

| Code | Meaning |
|------|---------|
| `invalid_config` | Invalid client or request configuration |
| `auth` | Signing, token, or authentication failure |
| `transport` | Network or transport failure |
| `unsupported` | Scaffolded API surface that is not implemented yet |

HTTP error messages include the API's own message when the response body carries
one (`message`, `msg`, `error_msg`, `errorMessage`, or `error_description`),
otherwise the HTTP status text.

## Matching errors

Two sentinels are exported for programmatic matching:

- `client.ErrAccessTokenRequired` — returned by `Client.Do` in production when
  automatic token handling is enabled but no usable token is cached.
- `client.ErrCircuitOpen` — wrapped into the error returned by `Client.Do` when a
  configured circuit breaker rejects a call.

```go
if errors.Is(err, client.ErrAccessTokenRequired) {
	// Production: run cl.EnsureToken(ctx) once to complete 2FA, then retry.
}
if errors.Is(err, client.ErrCircuitOpen) {
	// A circuit breaker is open; back off before retrying.
}
```

Because every SDK error implements `Unwrap`, `errors.Is` also traverses to the
underlying cause. Context cancellation surfaces as an error that matches
`context.Canceled` or `context.DeadlineExceeded`, and transport failures keep
their `net` error chain.

```go
if errors.Is(err, context.Canceled) {
	// The caller cancelled the request.
}
```

The concrete typed error is an implementation detail and is not exported.
Classification beyond the two sentinels above is by the stable code token in the
message, as shown in the tables.

## Handling guidance

- Treat `UNAUTHORIZED` and `INVALID_TOKEN` as authentication failures: obtain a
  fresh token with `Client.EnsureToken` and retry once.
- Treat `FORBIDDEN` as a configuration or entitlement problem; retrying will not
  help. See [Troubleshooting](troubleshooting.md).
- Treat `RATE_LIMITED` as transient. The SDK retries transient errors for
  idempotent requests by default; otherwise back off and retry.
- Treat `SERVER_ERROR` and `transport` as transient. Configure `WithRetry`,
  `WithRateLimiter`, or `WithBreaker` if the defaults are not enough.
- Treat `invalid_config` as a programming error and fail fast.

### Transient vs permanent errors

| Error code | Transient? | Retry? |
|------------|-----------|--------|
| `UNAUTHORIZED` | No | Fix credentials, then retry once |
| `FORBIDDEN` | No | Never — check entitlements and permissions |
| `INVALID_TOKEN` | No | Call `EnsureToken`, then retry once |
| `RATE_LIMITED` | Yes | Back off exponentially; SDK retries idempotent requests |
| `SERVER_ERROR` | Yes | Retry with backoff; configure `WithRetry` |
| `transport` | Yes | Retry with backoff; check network connectivity |
| `invalid_config` | No | Never — fix the code |

### Rate limit pattern

```go
// The SDK retries automatically for idempotent requests. For non-idempotent
// requests (PlaceOrder, etc.), implement exponential backoff:
backoff := time.Second
for i := 0; i < maxRetries; i++ {
    _, err := trading.PlaceOrder(ctx, req)
    if err == nil {
        break
    }
    if !strings.Contains(err.Error(), "RATE_LIMITED") {
        return err
    }
    time.Sleep(backoff)
    backoff *= 2
}
```

See [Authentication](authentication.md) for the token lifecycle and
[Troubleshooting](troubleshooting.md) for sandbox-specific error causes.
