# Errors

Every SDK function returns a normal Go `error`. SDK errors are typed, carry a
stable `errs.Code`, and preserve wrapped causes for `errors.Is` and
`errors.As`.

!!! note "Prerequisites"
    - Familiarity with Go's `errors` package
    - An SDK client returning errors — see [Getting Started](getting-started.md)

## Typed error shape

`pkg/errors.Error` is public:

```go
type Error struct {
    Code    errs.Code
    Message string
    Status  int
    Err     error
}
```

`Status` is the HTTP status when the error came from an HTTP response and zero
otherwise. `Unwrap` preserves the underlying network, context, gRPC, or
operation error.

```go
import (
    "context"
    "errors"
    "fmt"

    errs "github.com/shing1211/webullapi4go/pkg/errors"
)

func classify(err error) error {
    if err == nil {
        return nil
    }

    var sdkErr *errs.Error
    if errors.As(err, &sdkErr) {
        fmt.Printf("code=%s status=%d message=%s\n",
            sdkErr.Code, sdkErr.Status, sdkErr.Message)
    }

    switch {
    case errors.Is(err, context.Canceled):
        return err
    case errs.Is(err, errs.CodeRateLimited):
        return err
    case errors.Is(err, errs.ErrForbidden):
        return err
    default:
        return err
    }
}
```

Import the standard `errors` package and the SDK package with an alias such as
`errs`. Do not compare `err.Error()` text.

## Error format

A typed error formats as one of:

```text
webull: <code>: <message>
webull: <code>: <message>: <cause>
```

Examples:

```text
webull: FORBIDDEN: http 403: Insufficient permission
webull: transport: GET /market-data/stocks/snapshots/list: dial tcp: connection refused
```

The code is the stable classification. The message is for diagnostics and may
change.

## Codes

### HTTP-derived codes

| Code | Source |
|---|---|
| `UNAUTHORIZED` | HTTP 401 |
| `FORBIDDEN` | HTTP 403 |
| `INVALID_TOKEN` | HTTP 417 |
| `RATE_LIMITED` | HTTP 429 |
| `SERVER_ERROR` | HTTP 5xx |
| `api` | Any other non-2xx response or response decode failure |

### SDK codes

| Code | Meaning |
|---|---|
| `invalid_config` | Invalid client or request configuration |
| `auth` | Signing, token, subscription, or authentication failure |
| `transport` | Network, HTTP transport, or reconnectable stream failure |
| `unsupported` | The requested operation or protocol feature is unsupported |
| `validation` | An input value failed validation outside client construction |
| `not_initialized` | A required client or tracked resource was not initialized |
| `invalid_transition` | An operation is illegal for the current order state |
| `order_guardrail` | A risk guardrail rejected an order; normally wrapped by an outer `invalid_config` error for compatibility |

HTTP messages use the API's `message`, `msg`, `error_msg`, `errorMessage`, or
`error_description` field when present. A short raw body is retained when no
recognized field exists.

## Sentinels and matching

Export a sentinel when callers need semantic matching; do not parse a code from
text.

```go
if errors.Is(err, client.ErrAccessTokenRequired) {
    // Complete the production token/2FA flow explicitly.
}
if errors.Is(err, client.ErrCircuitOpen) {
    // Back off until the breaker permits another attempt.
}
if errs.Is(err, errs.CodeRateLimited) {
    // Stable code-based classification.
}
if errors.Is(err, errs.ErrOrderGuardrail) {
    // A quantity or notional guardrail stopped the order.
}
```

Public sentinels in `pkg/errors` include:

- `ErrUnsupported`
- `ErrUnauthorized`, `ErrForbidden`, `ErrInvalidToken`, `ErrRateLimited`,
  `ErrServer`
- `ErrValidation`, `ErrNotInitialized`, `ErrInvalidTransition`
- `ErrOrderGuardrail`
- `ErrSubscriptionExpired`, `ErrConnectionLimitExceeded`

`errors.Is` on two `*errs.Error` values matches when their `Code` values are
equal, so both of these work:

```go
errors.Is(err, errs.ErrRateLimited)
errs.Is(err, errs.CodeRateLimited)
```

The underlying cause remains in the chain. Context cancellation and deadlines
match `context.Canceled` and `context.DeadlineExceeded`; network failures keep
their `net` error chain.

## Order errors

OMS methods add a resource-specific layer while retaining the public typed
error model:

```go
import (
    "errors"

    "github.com/shing1211/webullapi4go/pkg/domain/order"
    errs "github.com/shing1211/webullapi4go/pkg/errors"
    "github.com/shing1211/webullapi4go/trade"
)

state, err := trading.ApplyOrderEvent(accountID, clientOrderID, order.EventAcknowledge)
switch {
case err == nil:
    // state was updated
case errs.Is(err, errs.CodeNotInitialized):
    // The order is not tracked by this client instance.
case errs.Is(err, errs.CodeInvalidTransition):
    // The event is not legal for the current state.
case errors.Is(err, order.ErrInvalidTransition):
    // The underlying domain transition error is also preserved.
}
```

A guardrail error intentionally keeps the historical outer
`errs.CodeInvalidConfig` classification and also wraps
`errs.ErrOrderGuardrail`. Existing checks for either classification continue
to work.

## Handling guidance

| Category | Transient? | Recommended handling |
|---|---|---|
| `UNAUTHORIZED`, `INVALID_TOKEN` | Authentication failure | Create or refresh a token, then retry once where safe |
| `FORBIDDEN` | No | Check account scope and data entitlement |
| `RATE_LIMITED` | Yes | Back off; the SDK retries eligible idempotent requests |
| `SERVER_ERROR`, `transport` | Usually | Retry idempotent operations with backoff; inspect connectivity and breaker state |
| `invalid_config`, `validation` | No | Fix the request or client configuration |
| `unsupported` | No | Select another operation or SDK version |
| `not_initialized` | No | Construct/run/register the required resource first |
| `invalid_transition` | No | Reconcile the current resource state before retrying |
| `order_guardrail` | No | Reduce size/notional or change the application guardrail explicitly |

Only GET and HEAD are retried by default. Enabling non-idempotent retries can
duplicate an order placement unless the caller has a safe idempotency policy.
A stable `client_order_id` is still required by the application.

### Context-aware rate-limit loop

```go
func getBalanceWithBackoff(
    ctx context.Context,
    trading *trade.Client,
    accountID string,
) (*trade.AssetsBalance, error) {
    backoff := 200 * time.Millisecond
    for attempt := 1; attempt <= 3; attempt++ {
        balance, err := trading.GetBalance(ctx, accountID)
        if err == nil {
            return balance, nil
        }
        if !errs.Is(err, errs.CodeRateLimited) {
            return nil, err
        }
        if attempt == 3 {
            return nil, err
        }

        timer := time.NewTimer(backoff)
        select {
        case <-ctx.Done():
            timer.Stop()
            return nil, ctx.Err()
        case <-timer.C:
        }
        backoff *= 2
    }
    return nil, errs.New(errs.CodeRateLimited, "balance retries exhausted")
}
```

## Related

- [Observability](observability.md) — errors in logs, spans, and metrics
- [Trading](trading.md) — OMS reconciliation and guardrails
- [Streaming](streaming.md) — asynchronous and reconnect errors
- [Authentication](authentication.md) — token lifecycle
- [Troubleshooting](troubleshooting.md) — sandbox-specific failures
