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

## Numeric strings

Prices, quantities, and other numeric fields are **strings** on the wire and
remain strings in Go types. This preserves decimal precision that `float64`
would lose:

```go
log.Printf("price=%s qty=%s", order.LimitPrice, order.Quantity)
```

When constructing requests, pass numeric values as strings:

```go
Quantity:  "100",
LimitPrice: "150.25",
```

## Error handling

All SDK errors are typed and support `errors.Is`/`errors.As`:

```go
result, err := trading.PlaceOrder(ctx, req)
if err != nil {
    if errors.Is(err, client.ErrAccessTokenRequired) {
        // Token expired; refresh and retry
        if _, err := cl.EnsureToken(ctx); err != nil {
            log.Fatal(err)
        }
        // retry...
    } else if errors.Is(err, client.ErrCircuitOpen) {
        // Circuit breaker open; back off
        time.Sleep(backoff)
    } else {
        log.Printf("order failed: %v", err)
    }
}
```

See [Errors](errors.md) for the full code table and handling guidance.

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

MQTT streaming and gRPC events follow a callback + run pattern:

```go
// MQTT (Market Data)
s, err := stream.New(cl, stream.WithWebSocket(true))
s.Connect(ctx)
s.Subscribe(ctx, stream.SubscribeRequest{...})
s.OnQuote(func(q *marketdatav1.Quote) { /* handle */ })

// gRPC (Trading Events)
ev, err := events.New(cl)
ev.OnOrder(func(o *events.OrderEvent) { /* handle */ })
ev.Run(ctx) // blocks
```

Register handlers before calling `Connect`/`Run`. Handlers run on the stream
pump and must not block.

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
- [Sandbox](sandbox.md) — environments and test credentials.
