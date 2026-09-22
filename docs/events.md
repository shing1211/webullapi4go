# Trading events

The `events` package delivers Webull's trade events (order status changes,
event-contract position settlements, and option status changes) over a
persistent, server-streaming gRPC connection. Requires a trade account ID
(see [Trading](trading.md)). Unlike the REST API, the event
service authenticates each `Subscribe` call with an HMAC-SHA256 signature over
the serialized request, so no access token is involved. Build a client from a
configured [`*client.Client`](api.md#client), which supplies the credentials,
region, and resolved gRPC endpoint:

```go
cl, err := client.New(client.WithEnv())
if err != nil {
	return err
}
defer func() { _ = cl.Close() }()

ev, err := events.New(cl)
if err != nil {
	return err
}
defer func() { _ = ev.Close() }()
```

!!! note "Prerequisites"
    - A [Webull account](https://developer.webull.hk/apis/docs/sdk#test-accounts) (sandbox or production)
    - Go 1.26+
    - App key and app secret (HMAC-SHA256 signing, no access token required)
    - A trading account ID — see [Trading](trading.md)

!!! tip "Error handling"
    All SDK functions return `error`. See [Errors](errors.md) for the typed error model, transient vs permanent classification, and retry patterns.

## Service and endpoint

The service is the server-streaming RPC
`grpc.trade.event.EventService/Subscribe`. The endpoint is
`events-api.<domain>:443` over TLS, resolved from the client's region and
environment:

| Environment | Endpoint |
|-------------|----------|
| Production (HK) | `events-api.webull.hk:443` |
| Sandbox (HK) | `events-api.sandbox.webull.hk:443` |

`WithGRPCEndpoint` overrides the host and `WithGRPCPort` the port;
`WithTLS(false)` selects plaintext for tests and local servers only.

## Signing

The event service signs a different canonical string from the REST API. Only the
five `x-signature` parameters participate — there is **no `host`**, path, or
query — and the serialized request body is hashed into the string:

1. Join the parameters, sorted by name, as `name=value` with `=` separators.
2. Append `&` and the **lowercase** hexadecimal SHA-256 digest of the
   serialized `SubscribeRequest` proto bytes.
3. Percent-encode the result with the RFC 3986 unreserved set.
4. Sign it with HMAC-SHA256 keyed by `app_secret + "&"`, then base64-encode.

The five parameters and the signature are sent as gRPC metadata:

| Metadata key | Value |
|--------------|-------|
| `x-app-key` | App key |
| `x-signature-algorithm` | `HMAC-SHA256` |
| `x-signature-version` | `1.0` |
| `x-signature-nonce` | Random nonce |
| `x-timestamp` | UTC timestamp, `2006-01-02T15:04:05Z` |
| `x-signature` | Base64 HMAC-SHA256 signature |

!!! warning "Lowercase digest"
    The event canonical string hashes the body to **lowercase** hex. The REST
    signer (`internal/auth`) uppercases the SHA-256 digest and would be rejected
    by the event service, so `events` computes its signature directly rather than
    calling that signer.

## Subscribe types

`WithSubscribeTypes` selects the event categories. The value is a bitmask, so
the individual constants OR together; `SubscribeAll` is the default.

| Constant | Value | Events |
|----------|-------|--------|
| `SubscribeOrder` | `1` | Order status changes |
| `SubscribePosition` | `2` | Event-contract position settlements |
| `SubscribeOption` | `4` | Option status changes |
| `SubscribeAll` | `7` | Every category |

`WithAccounts` restricts the stream to specific trading account ids. The
accounts must belong to the App Key: an account that does not match the
credentials is not streamed.

```go
ev, err := events.New(cl,
	events.WithSubscribeTypes(events.SubscribeOrder),
	events.WithAccounts([]string{"<account-id>"}),
)
```

## Dispatch model

`Run` connects, subscribes, and blocks pumping events until the context is
cancelled or the stream ends. Each response is routed by its event type:

| Server event | Delivered to |
|--------------|--------------|
| `SubscribeSuccess` | `OnConnect` |
| `Ping` | `OnPing` (heartbeat) |
| `AuthError`, `NumOfConnExceed`, `SubscribeExpired` | `OnError` with a typed error; ends `Run` |
| Data event (`1024` order, `1028` position, `1032` option) | `OnEvent` raw, then the matching typed handler |

The raw event kinds dispatched on the wire:

| Constant | Value | Delivered to |
|----------|-------|-------------|
| `EventOrder` | `1024` | `OnOrder` |
| `EventPosition` | `1028` | `OnPosition` |
| `EventOption` | `1032` | `OnOption` |

Data events are delivered to `OnEvent` with the kind and MIME content type and,
when the content type is `application/json`, decoded and delivered to `OnOrder`,
`OnPosition`, or `OnOption`. A payload that fails to decode is reported to
`OnError` and the stream continues. Registering a handler is optional; `OnEvent`
fires before the typed handler so both may be used.

Handlers run on the stream pump, so they must not block: a slow handler delays
later events.

## Order event payload

Webull does not publish a formal schema for the JSON payload. The `OrderEvent`
type mirrors the documented and observed keys; unknown keys are ignored and
`OrderEvent.Raw` always holds the exact bytes the server sent. Every numeric
value is kept as a string to preserve precision.

Captured sandbox example:

```json
{"secAccountId":66600004338,"requestId":"<event-request-id>","account_id":"<acct>","request_id":"<order-id>","order_id":"<order-id>","client_order_id":"sdk-evt-<ts>","instrument_id":"913256135","order_status":"CANCELLED","symbol":"AAPL","short_name":"Apple Inc","qty":"1.0000000000","filled_price":"0E-10","filled_qty":"0E-10","side":"BUY","scene_type":"CANCEL_SUCCESS","category":"US_STOCK","order_type":"LIMIT","actual_commission":"0E-10","receivable_commission":"0E-10"}
```

| JSON key | Go field | Type | Description |
|----------|----------|------|-------------|
| `secAccountId` | `SecAccountID` | `json.Number` | Securities account identifier |
| `requestId` | `EventRequestID` | string | Event-level request identifier |
| `account_id` | `AccountID` | string | Account the order belongs to |
| `request_id` | `RequestID` | string | Order request identifier |
| `order_id` | `OrderID` | string | System-generated order identifier |
| `client_order_id` | `ClientOrderID` | string | Caller-supplied order identifier |
| `instrument_id` | `InstrumentID` | string | Webull instrument identifier |
| `order_status` | `OrderStatus` | string | Lifecycle state, for example `SUBMITTED`, `FILLED`, `CANCELLED` |
| `symbol` | `Symbol` | string | Trading symbol |
| `short_name` | `ShortName` | string | Human-readable instrument name |
| `qty` | `Quantity` | string | Total order quantity |
| `filled_qty` | `FilledQuantity` | string | Quantity executed so far |
| `filled_price` | `FilledPrice` | string | Average execution price |
| `filled_time` | `FilledTime` | string | Time of the last execution, when supplied |
| `side` | `Side` | string | `BUY` or `SELL` |
| `scene_type` | `SceneType` | string | Scenario, for example `CANCEL_SUCCESS`, `FILLED`, `FINAL_FILLED` |
| `category` | `Category` | string | Instrument category, for example `US_STOCK` |
| `order_type` | `OrderType` | string | Execution instruction, for example `LIMIT` |
| `biz_type` | `BizType` | string | Business type, `TRADE` for order events |
| `actual_commission` | `ActualCommission` | string | Commission collected |
| `receivable_commission` | `ReceivableCommission` | string | Commission still receivable |

`PositionEvent` and `OptionEvent` model the same style of JSON for the position
and option streams. `OptionEvent` carries the order fields above with
`category` `US_OPTION`; `PositionEvent` carries `account_id`, `symbol`, and the
settlement fields (`position_id`, `event_name`, `yes_condition`, `settle_result`,
`settle_side`, `quantity`, `cost`, `settle_amount`, and `biz_type`).

## Reconnect and backoff

`Run` reconnects and re-subscribes by default. A clean server end and transient
transport failures are retried with exponential backoff plus jitter; terminal
failures — authentication, permission, account, and configuration errors, plus
the terminal server events — end `Run`. Consecutive attempts are bounded by
`WithMaxReconnectAttempts` (unlimited by default).

| Option | Default | Purpose |
|--------|---------|---------|
| `WithAutoReconnect(bool)` | `true` | Reconnect and re-subscribe after a stream ends or fails |
| `WithReconnectBaseDelay(d)` | `1s` | First delay before reconnecting |
| `WithReconnectMaxDelay(d)` | `30s` | Cap on the exponential delay |
| `WithMaxReconnectAttempts(n)` | `0` (unlimited) | Consecutive attempts before giving up |
| `WithDialTimeout(d)` | `10s` | Bound on waiting for the connection to become ready |
| `WithGRPCDialOption(opt)` | — | Extra gRPC dial options, such as a custom dialer |
| `WithAccounts([]string)` | none | Account filter (copied) |
| `WithSubscribeTypes(...)` | `SubscribeAll` | Event categories |

Defaults: `DefaultGRPCPort` = 443, `DefaultDialTimeout` = 10s,
`DefaultReconnectBaseDelay` = 1s, `DefaultReconnectMaxDelay` = 30s.

`Close` cancels an in-progress `Run` and closes the gRPC connection; it is
idempotent. The underlying `client.Client` is owned by the caller and is not
closed by `Close`.

## Full example

The runnable example in
[`examples/events`](https://github.com/shing1211/webullapi4go/tree/main/examples/events)
subscribes to order events, prints each decoded event, and handles Ctrl+C:

```go
ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
defer cancel()

opts := []events.Option{events.WithSubscribeTypes(events.SubscribeOrder)}
if accountID != "" {
    opts = append(opts, events.WithAccounts([]string{accountID}))
}

ev, err := events.New(cl, opts...)
if err != nil {
    log.Fatal(err)
}
defer func() { _ = ev.Close() }()

ev.OnConnect(func() { log.Println("subscribed") })
ev.OnError(func(err error) { log.Printf("event stream error: %v", err) })
ev.OnOrder(func(o *events.OrderEvent) {
    log.Printf("order %s status=%s side=%s type=%s qty=%s filled_qty=%s filled_price=%s symbol=%s",
        o.OrderID, o.OrderStatus, o.Side, o.OrderType,
        o.Quantity, o.FilledQuantity, o.FilledPrice, o.Symbol)
})

if err := ev.Run(ctx); err != nil {
    log.Fatal(err)
}
```

## Sandbox caveats

- The sandbox may not push a placement event for a resting order; only
  `CANCEL_SUCCESS` has been observed, so trigger a cancel to see an event.
- An account must belong to the App Key. A mismatched account is not streamed.
- The sandbox accepts the same `events-api.sandbox.webull.hk:443` endpoint; TLS
  is always on.
- Webull allows at most **5 concurrent event connections** per App Key;
  exceeding the limit surfaces `NumOfConnExceed` on `OnError`.
- The live tests are gated by `WEBULL_SANDBOX=1`, `WEBULL_APP_KEY`,
  `WEBULL_APP_SECRET`, and (for the mutating order-event test)
  `WEBULL_TRADE_ACCOUNT_ID` and `WEBULL_TRADE_MUTATE=1`. Never commit
  credentials.

## Related

- [Trading](trading.md) — placing and cancelling the orders that emit events.
- [Streaming](streaming.md) — Market Data pushes over MQTT.
- [Sandbox](sandbox.md) — environments, test credentials, limitations.
- [Errors](errors.md) — typed errors and classification.
