# Streaming

The `stream` package delivers Market Data over MQTT. Subscribe and unsubscribe
are HTTP calls made through the core client; the MQTT connection carries pushes.
A valid access token is required. See [Authentication](authentication.md).

!!! note "Prerequisites"
    - A [Webull account](https://developer.webull.hk/apis/docs/sdk#test-accounts) (sandbox or production)
    - Go 1.26+
    - An authenticated core client

!!! tip "Error handling"
    Stream callbacks and connection methods return or receive typed errors.
    See [Errors](errors.md).

## Construct and connect

```go
cl, err := client.New(client.WithEnv())
if err != nil {
    return err
}
defer func() { _ = cl.Close() }()

if _, err := cl.EnsureToken(ctx); err != nil {
    return err
}

s, err := stream.New(cl,
    stream.WithWebSocket(true),
    stream.WithAutoReconnect(true),
    stream.WithHealthWatchdog(30*time.Second),
)
if err != nil {
    return err
}
defer func() { _ = s.Close() }()

if err := s.Connect(ctx); err != nil {
    return err
}
```

The broker address comes from the core client's resolved endpoints.
`WithWebSocket(true)` selects MQTT over WebSocket; otherwise plain TCP is used.
`WithMQTTURL` overrides both. In constrained networks, prefer the configured
WebSocket endpoint because plain MQTT port `1883` may be blocked.

### Connect, cancellation, and close

`Connect(ctx)` blocks until the broker accepts the connection, the attempt
fails, the context ends, or the client is closed. The low-level MQTT adapter
checks an already-cancelled context before starting a broker connection, so
that case performs no network attempt. If cancellation occurs after the broker
attempt starts, the adapter disconnects, returns `context.Canceled` (preserved
through the stream client's transport wrapper), and does not leave a token
waiter goroutine behind.

`Close` is terminal and idempotent. It first makes the stream state
`StateClosed`, cancels the watchdog and replay contexts, closes every channel
subscription, and then closes the MQTT transport. A `Connect` racing `Close`
returns without starting another connection; a connection attempt already in
progress is released. Late MQTT connect, reconnect, disconnect, error, and
message callbacks are suppressed after close. Calling `Connect`, `Subscribe`,
or `Unsubscribe` after stream close returns `errs.CodeInvalidConfig`.

## Subscribe and unsubscribe

```go
if err := s.Subscribe(ctx, stream.SubscribeRequest{
    Symbols:  []string{"AAPL"},
    Category: stream.CategoryUSStock,
    SubTypes: []stream.SubType{
        stream.SubTypeQuote,
        stream.SubTypeSnapshot,
        stream.SubTypeTick,
    },
    Grab: true,
}); err != nil {
    return err
}
```

`SubscribeRequest` fields:

| Field | Purpose |
|---|---|
| `SessionID` | Overrides the client session ID for this call; rarely needed |
| `Symbols` | Symbols to subscribe to; at most 100 per request |
| `Category` | Security category such as `CategoryUSStock`, `CategoryUSETF`, `CategoryHKStock`, or `CategoryCNStock` |
| `SubTypes` | `SubTypeQuote`, `SubTypeSnapshot`, and/or `SubTypeTick` |
| `Grab` | Request an immediate snapshot push |
| `Depth` | Level-2 depth string; optional |
| `OvernightRequired` | Include the US overnight session |

`Unsubscribe` mirrors the request and accepts `UnsubscribeAll: true` to clear
the active subscription set for the session.

Subscription mutations and reconnect replay are serialized. A reconnect cannot
interleave with a concurrent `Subscribe` or `Unsubscribe` request.

## Typed handlers

Handlers may be registered before or after `Connect`; each registration is
retained and registration is concurrency-safe.

| Handler | Payload | Fires on |
|---|---|---|
| `OnQuote` | `*marketdatav1.Quote` | Quote pushes |
| `OnSnapshot` | `*marketdatav1.Snapshot` | Snapshot pushes |
| `OnTick` | `*marketdatav1.Tick` | Tick pushes |
| `OnNotice` | `[]byte` | Server notice JSON |
| `OnError` | `error` | Decode, topic, connection, or resubscribe errors |
| `OnConnect` | none | Initial connection and successful reconnects |
| `OnDisconnect` | `error` | Loss of an established connection |
| `OnReconnecting` | none | Start of a reconnect attempt |
| `OnStateChange` | `(previous, next State)` | Every actual lifecycle state change |

Callbacks run synchronously on the stream receive/dispatch path, in registration
order. Keep them short. A callback that blocks delays MQTT dispatch and health
processing. For each data topic, channel subscribers are then dispatched
synchronously in their registration order.

## Connection state and health

`Client.State()` returns one of six states:

| State | Meaning |
|---|---|
| `StateDisconnected` | No live connection, or disconnected without an active reconnect |
| `StateConnecting` | A connection attempt is in progress |
| `StateConnected` | MQTT is live and the data-age health signal is fresh |
| `StateReconnecting` | A lost connection is being re-established |
| `StateDegraded` | MQTT is live, but no data message arrived within the watchdog interval |
| `StateClosed` | `Close` was called; this state is terminal |

```go
s.OnStateChange(func(previous, next stream.State) {
    log.Printf("stream state %s -> %s", previous, next)
})
```

State writes use compare-and-swap transitions. Health and data-arrival updates
also require the expected current state, so a delayed watchdog tick or late
message cannot overwrite a newer reconnect/close state. `StateClosed` is never
left by a transition. Handlers are called synchronously, in registration
order, only after the new state has been committed.

`WithHealthWatchdog(interval)` starts one watchdog after the first successful
`Connect`. Only quote, snapshot, and tick messages refresh the data timestamp.
Notices, echo heartbeats, and other non-data messages do not.

When the age of the last data message exceeds the interval, a stream whose
current state is still `StateConnected` becomes `StateDegraded`. Only a new
quote, snapshot, or tick can move a stream whose current state is
`StateDegraded` back to `StateConnected`; repeated stale watchdog ticks cannot.
This compare-and-swap behavior makes recovery deterministic even when a tick,
message callback, reconnect, and close race. `StateDegraded` is a health
signal, not an automatic disconnect, and `IsConnected` may still be true.

Webull MQTT payloads do not include sequence numbers. The SDK cannot prove
message-sequence continuity, so this watchdog measures data age only. Use
`Grab: true` when subscribing or re-snapshotting after a gap when the endpoint
supports it.

`StateClosed` cannot transition back to a live state. `Connect`, `Subscribe`,
and `Unsubscribe` return an `invalid_config` error after close; late MQTT
callbacks and messages are ignored.

## Reconnection and resubscribe

With `WithAutoReconnect(true)`, a connection loss starts the background
reconnect loop. After a successful reconnect, the client:

1. marks the connection live;
2. snapshots the active subscription registry;
3. re-issues each active HTTP subscription exactly once, bounded by
   `WithResubscribeTimeout`;
4. invokes `OnConnect` handlers after replay finishes.

Re-subscription is idempotent and failures are delivered to `OnError`; one
failed request does not prevent the remaining requests from being attempted.
`Reconnecting` reports an in-progress reconnect and `IsConnected` reports false
while reconnecting.

Webull limits an App Key to five concurrent MQTT connections. Code 105 is
mapped to a transport error explaining the limit. After a disconnect, server
state may be retained for about a minute, so avoid immediate session-ID reuse.
Reusing a session ID can disconnect the previous connection.

## Channel delivery

Channel subscriptions provide a bounded alternative to callbacks:

```go
quotes, cancel := s.SubscribeQuoteChan(stream.ChannelConfig{
    Policy:     stream.DropOldest,
    BufferSize: 256,
})
defer cancel()

for quote := range quotes {
    if quote == nil {
        continue
    }
    log.Printf("%s", quote.GetBasic().GetSymbol())
}
```

The three constructors are:

- `SubscribeQuoteChan`
- `SubscribeSnapshotChan`
- `SubscribeTickChan`

Each returns a receive-only channel and an idempotent cancel function. The
channel closes when its cancel function runs or when `Client.Close` shuts down
the stream. Calling cancel after close is safe. `BufferSize <= 0` uses 100.

| Policy | Full-buffer behavior | Consequence |
|---|---|---|
| `DropBlock` | Wait for buffer space or cancellation | Backpressure; can pause the MQTT pump |
| `DropOldest` | Remove the oldest unread value, then enqueue the new value | Preserves the newest value; older values are discarded |
| `DropSample` | Randomly discard an arriving value when full | Reduces load with probabilistic loss |

The default is `DropBlock`. Choose a drop policy only when the consumer can
tolerate data loss. Dispatch is deliberately synchronous: a full `DropBlock`
subscriber causes head-of-line blocking and prevents later subscribers of the
same topic from receiving that message until the blocked send completes. Use
`DropOldest` or `DropSample` when one slow consumer must not hold the MQTT
pump. A blocked dispatch is unblocked by that channel's idempotent cancel
function or by `Client.Close`, so a full channel does not permanently wedge
shutdown.

Dropped messages are counted internally and, when a meter is configured,
recorded as OTel `channel_drops` with a `topic` attribute. There is no public
per-channel drop-count accessor; export the metric for operational visibility.

## Options

| Option | Purpose |
|---|---|
| `WithSessionID`, `WithClientID` | Set the MQTT client/session ID; a unique ID is generated by default |
| `WithMQTTURL` | Override the broker address |
| `WithWebSocket` | Select MQTT over WebSocket |
| `WithAutoReconnect` | Enable the reconnect loop; disabled by default |
| `WithAutoResubscribe` | Replay active HTTP subscriptions after reconnect; enabled by default |
| `WithResubscribeTimeout` | Bound the complete replay sequence |
| `WithKeepAlive` | MQTT keep-alive interval |
| `WithConnectTimeout` | Bound one connection attempt |
| `WithWriteTimeout` | Bound a control-packet write |
| `WithMessageChannelDepth` | Low-level MQTT inbound buffer depth |
| `WithCleanSession` | Request a clean MQTT session; default true |
| `WithTLSConfig` | Override TLS configuration |
| `WithHealthWatchdog` | Enable the data-age health watchdog; zero disables it |
| `WithMeter` | Override the stream OTel meter; otherwise the core client's meter is inherited |

## Topics

| Constant | Value | Encoding |
|---|---|---|
| `TopicQuote` | `quote` | Protobuf `Quote` |
| `TopicSnapshot` | `snapshot` | Protobuf `Snapshot` |
| `TopicTick` | `tick` | Protobuf `Tick` |
| `TopicNotice` | `notice` | JSON bytes |
| `TopicEcho` | `echo` | Heartbeat; ignored |

## Observability

When the core client has an OTel tracer, each MQTT message dispatch can create
a consumer span named `mqtt.dispatch` with messaging system, destination, and
payload-size attributes. Reconnects and channel drops are optional metrics. The
core logger can record connect, disconnect, and reconnect transitions.

See [Observability](observability.md) for setup and metric names.

## Full example

The runnable
[`examples/streaming`](https://github.com/shing1211/webullapi4go/tree/main/examples/streaming)
program connects over WebSocket, subscribes to AAPL quote/snapshot/tick pushes,
prints messages, and handles Ctrl+C. Live behavior is environment-dependent;
the current health/channel hardening is offline-tested, not newly live-verified.

## Related

- [Observability](observability.md) — MQTT spans and metrics
- [Market Data](market-data.md) — HTTP queries
- [Errors](errors.md) — typed stream errors
- [Troubleshooting](troubleshooting.md) — connection and entitlement issues
