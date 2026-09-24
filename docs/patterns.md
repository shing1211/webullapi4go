# Patterns

Shared patterns and conventions used across the SDK.

!!! note "Prerequisites"
    - A [Webull account](https://developer.webull.hk/apis/docs/sdk#test-accounts) (sandbox or production)
    - Go 1.26+
    - Basic familiarity with Go's functional options pattern

## Client construction

Every API package follows the same construction pattern:

```go
cl, err := client.New(client.WithEnv())
if err != nil {
    log.Fatal(err)
}
defer func() { _ = cl.Close() }()

ctx := context.Background()
if _, err := cl.EnsureToken(ctx); err != nil {
    log.Fatal(err)
}
```

`WithEnv()` reads credentials from the environment. Alternative constructors:

| Option | Purpose |
|--------|---------|
| `client.WithCredentials(key, secret)` | Explicit credentials |
| `client.WithSandbox()` | Point at the HK sandbox host |
| `client.WithRegion(client.HK)` / `client.WithRegion(client.US)` | Select region |
| `client.WithBaseURL(url)` | Override the base URL |
| `client.WithTimeout(d)` | Set HTTP timeout |
| `client.WithAutoToken(true)` | Auto-obtain token on first sandbox request |

Options are applied in order. `WithSandbox()` is equivalent to
`WithEnvironment(client.Sandbox)`.

`client`, `data`, `trade`, `stream`, and `events` remain the canonical service
packages. `pkg/` contains shared foundations such as `errors`, `observability`,
`resilience`, `transport`, `domain/money`, and `domain/order`; it is not a
relocated service hierarchy. The optional `webull` package only aliases the
core client and does not replace the root service constructors.

## Option pattern

All packages use the functional options pattern:

```go
package, err := pkg.New(cl,
    pkg.WithOption1(value),
    pkg.WithOption2(value),
)
```

Each `WithX` constructor returns a `func(*Config)` that is applied during
construction. Options are validated immediately; invalid options fail the
constructor.

## Pagination

List endpoints use cursor-based pagination with three tiers:

| Method | Behavior |
|--------|----------|
| `GetXxx(ctx, q)` | Returns the first page only |
| `GetXxxPage(ctx, q)` | Returns one page with `PaginationKey` for the next page |
| `GetAllXxx(ctx, q)` | Follows the cursor to exhaustion, returning all results |

```go
// First page
open, err := trading.GetOpenOrders(ctx, accountID)

// Manual pagination
page, err := trading.GetOpenOrdersPage(ctx, accountID, "")
for _, group := range page.Orders {
    // process group
}
if page.PaginationKey != "" {
    next, err := trading.GetOpenOrdersPage(ctx, accountID, page.PaginationKey)
    // ...
}

// Auto-pagination (all pages)
all, err := trading.GetAllOpenOrders(ctx, accountID)
```

`GetAllXxx` is bounded by `MaxOrderQueryPages` (100) to prevent infinite loops.

## Decimal money

Prices, quantities, balances, and other financial values are decimal strings on
the wire. Public DTOs preserve that precision with `money.Money` for required
and response values and `*money.Money` for optional and request values. Raw
`decimal.Decimal` is not the DTO API.

```go
import "github.com/shing1211/webullapi4go/pkg/domain/money"

quantity := money.MustNew("100")
limitPrice := money.MustNew("150.25")

req := trade.OrderRequest{
    Quantity:   &quantity,
    LimitPrice: &limitPrice,
}

log.Printf("price=%s qty=%s", req.LimitPrice, req.Quantity)
```

`MustNew` is appropriate for trusted constants and panics on invalid input.
For external input, handle the error:

```go
amount, err := money.NewFromString(input)
if err != nil {
    return err
}
```

`Money` implements `String`, comparison, arithmetic, and JSON methods. It
marshals back to a JSON string, so the wire format is unchanged.

## Error handling

Use the public `pkg/errors` codes and sentinels with `errors.Is`/`errors.As`;
never branch on `err.Error()` text.

```go
placed, err := trading.PlaceOrder(ctx, req)
switch {
case err == nil:
    log.Printf("placed %s", placed.OrderID)
case errors.Is(err, client.ErrAccessTokenRequired):
    // Complete token activation explicitly before retrying.
case errors.Is(err, client.ErrCircuitOpen):
    // Back off until the breaker permits another attempt.
case errors.Is(err, errs.ErrOrderGuardrail):
    // Reduce size/notional or change the application guardrail.
case errs.Is(err, errs.CodeRateLimited):
    // Back off; do not automatically retry a non-idempotent placement.
default:
    var sdkErr *errs.Error
    if errors.As(err, &sdkErr) {
        log.Printf("order failed: code=%s status=%d: %v",
            sdkErr.Code, sdkErr.Status, err)
    }
}
```

Import the standard library as `errors` and the SDK package with an alias:

```go
import (
    "errors"

    "github.com/shing1211/webullapi4go/client"
    errs "github.com/shing1211/webullapi4go/pkg/errors"
)
```

See [Errors](errors.md) for the code table and OMS error behavior.

## Display Service

Display Solution endpoints use a separate signing mechanism. Access them
through `data.Client.DisplayService()`:

```go
market := data.New(cl)
// The display service uses client-token auth, not the regular access token
profile, err := market.GetDSCompanyProfile(ctx, "AAPL")
```

The display service manages its own client token automatically. No separate
authentication setup is required.

## Streaming

MQTT streaming and gRPC events use callback or channel delivery:

```go
s, err := stream.New(cl,
    stream.WithWebSocket(true),
    stream.WithAutoReconnect(true),
    stream.WithHealthWatchdog(30*time.Second),
)
if err != nil {
    return err
}
defer func() { _ = s.Close() }()

s.OnQuote(func(q *marketdatav1.Quote) {
    // Keep callback work short; it runs on the MQTT pump.
})

quotes, cancel := s.SubscribeQuoteChan(stream.ChannelConfig{
    Policy:     stream.DropOldest,
    BufferSize: 256,
})
defer cancel() // idempotent; closes the channel

s.OnStateChange(func(previous, next stream.State) {
    log.Printf("stream %s -> %s", previous, next)
})
```

Only quote, snapshot, and tick data refresh the health watchdog. Notice and
echo traffic do not. `DropBlock` is the default; `DropOldest` preserves the
newest unread value, while `DropSample` randomly discards under pressure.
`Client.Close` closes all active channel subscriptions and unblocks a dispatch
waiting in `DropBlock`.

For gRPC Trading Events:

```go
ev, err := events.New(cl)
if err != nil {
    return err
}
defer func() { _ = ev.Close() }()

ev.OnOrder(func(e *events.OrderEvent) {
    // Decode and apply quickly; Run blocks on the receive loop.
})
if err := ev.Run(ctx); err != nil {
    return err
}
```

See [Streaming](streaming.md) and [Observability](observability.md).

## HK BCAN party IDs

Hong Kong equity orders require BCAN party IDs. Use the helper:

```go
NoPartyIDs: []trade.PartyID{{
    PartyID:       os.Getenv("WEBULL_TRADE_PARTY_ID"),
    PartyIDSource: trade.HKDerivativesPartyIDSource, // "D"
    PartyRole:     trade.HKDerivativesPartyRole,     // "3"
}},
```

Or use the builder:

```go
NoPartyIDs: trade.BuildHKDerivativesPartyIDs(partyID),
```

## Related

- [Getting Started](getting-started.md) — install and first call.
- [Authentication](authentication.md) — signing and tokens.
- [Errors](errors.md) — typed errors and classification.
- [Observability](observability.md) — logging, tracing, metrics, and correlation.
- [Sandbox](sandbox.md) — environments and test credentials.
