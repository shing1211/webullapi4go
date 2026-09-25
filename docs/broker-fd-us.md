# Broker FD API US

The `brokerfd` package provides the Webull Broker FD (US) HTTP API. It shares the
core `client.Client` for signing and transport. The `brokerfd/events` sub-package
provides the corresponding gRPC event stream.

!!! note "Prerequisites"
    - A [Webull account](https://developer.webull.hk/apis/docs/sdk#test-accounts) (sandbox or production)
    - Go 1.26+
    - An authenticated client — see [Authentication](authentication.md)

## Construct

`brokerfd.New` wraps the shared core client for the Broker FD HTTP API. The
`brokerfd/events` client is constructed independently from the same core
client; it does not require an HTTP `brokerfd.Client` instance.

```go
cl, err := client.New(client.WithEnv())
if err != nil {
    return err
}
defer func() { _ = cl.Close() }()

if _, err := cl.EnsureToken(ctx); err != nil {
    return err
}

fd := brokerfd.New(cl)
_ = fd
```

## Coverage

| Group | Methods |
|-------|---------|
| Summary | `GetAccountsSummary` (aggregate account data), `GetPositions` (aggregate position data) |
| Agreements | `ListAgreements`, `GetAgreementDetail` |
| Accounts | `ListFDAccounts`, `GetFDAccountDetail`, `CreateFDAccount`, `UpdateFDAccount`, `CloseFDAccount`, `ListAccountForms`, `GetAccountFormDetail`, `SubmitAccountForm`, `GetAccountFormStatus` |
| Documents | `UploadDocument`, `DownloadDocument`, `ListDocuments`, `GetDocumentDetail` |
| Assets | `GetFDAssetsSummary`, `GetFDAssetsDetail`, `GetFDPositions` |
| Activity | `GetFDActivities` |
| Funding | `ListFDBankAccounts`, `GetFDBankAccountDetail`, `AddFDBankAccount`, `RemoveFDBankAccount`, `ListFDAchAccounts`, `GetFDAchAccountDetail`, `AddFDAchAccount`, `RemoveFDAchAccount`, `ListFDTransfers`, `GetFDTransferDetail`, `InitiateFDTransfer`, `CreateFDInstantFunding`, `GetFDInstantFundingDetail`, `GetFDTransferFees`, `GetFDCreditInfo` |
| Instruments | `GetFDStockInstruments`, `GetFDStockLocate`, `GetFDCorporateActions`, `GetFDCorporateActionDetail`, `GetFDECInstruments`, `GetFDECInstrumentDetail` |
| Orders | `PreviewFDOrder`, `PlaceFDOrder`, `ReplaceFDOrder`, `CancelFDOrder`, `GetFDOrderDetail`, `GetFDOrderHistory`, `GetFDOpenOrders` |
| Journals | `ListFDCashJournals`, `GetFDCashJournalDetail` |
| Master data | `GetFDEnums`, `GetFDTradeCalendar` |

## Events (gRPC)

The Broker FD event client signs each `Subscribe` request with HMAC-SHA256 and
exposes lifecycle callbacks plus raw data delivery. It does not decode Broker
FD data into typed event structs.

```go
import "github.com/shing1211/webullapi4go/brokerfd/events"

ec, err := events.New(cl,
    events.WithAccounts([]string{"<account-id>"}), // optional
)
if err != nil {
    log.Fatal(err)
}
defer func() { _ = ec.Close() }()

ec.OnConnect(func() { log.Println("connected") })
ec.OnPing(func() { log.Println("ping") })
ec.OnError(func(err error) { log.Printf("error: %v", err) })
ec.OnData(func(subscribeType uint32, contentType string, payload []byte) {
    log.Printf("event type=%d content=%s bytes=%d", subscribeType, contentType, len(payload))
})

if err := ec.Run(ctx); err != nil {
    log.Fatal(err)
}
```

`SubscribeSuccess` and `Ping` invoke `OnConnect` and `OnPing`.
`AuthError`, `NumOfConnExceed`, and `SubscribeExpired` end `Run` with the same
typed categories and semantic identities used by Trading Events; see
[Errors](errors.md#mqtt-and-event-terminal-mappings). Other responses are
delivered to `OnData` with three values: the raw category value, MIME content
type, and payload bytes. The callback does not currently receive the response
request ID or timestamp. Although `DataEvent` and `SubscribeResponse.ToDataEvent`
model those fields, the client does not expose the received response object to
`OnData`; decode the available raw payload and treat those fields as a current
API limitation.

Handlers run synchronously in registration order on the stream receive loop.
Keep them short so one handler does not delay later Broker FD events.

`WithAccounts` is optional and copies the supplied IDs. An empty or omitted
list leaves the request's account filter empty; it does not itself select a
non-zero event-category bitmask.

The package exposes `WithGRPCEndpoint`, `WithGRPCPort`, `WithTLS`,
`WithDialTimeout`, `WithGRPCDialOption`, `WithAccounts`,
`WithAutoReconnect`, `WithReconnectBaseDelay`, `WithReconnectMaxDelay`, and
`WithMaxReconnectAttempts`. `WithMaxReconnectAttempts(n)` bounds retries after
the initial attempt; zero means unlimited.

!!! warning "Raw subscribe bitmask is not configurable through `Client` yet"
    Broker FD uses a raw `SubscribeType` bitmask, but the client has no
    `WithSubscribeTypes` option and `Run` does not accept a caller-built
    request. `NewSubscribeRequest` constructs the generated wire request for
    inspection or lower-level use; it does not inject that request into an
    existing `Client`. Consequently, the current high-level client cannot set
    a non-zero category bitmask. This remains an API gap pending a public
    option or request-injection API.

Unlike Trading Events, Broker FD tracks one active `Run`; concurrent `Run` calls
are not supported. `Close` is idempotent, cancels the active run, and closes
the gRPC connection. Context cancellation and `Close` make `Run` return
`context.Canceled` and are not sent to `OnError`. The client does not currently
provide the same all-runs registry or explicit post-close guard as Trading
Events.

Each attempt emits telemetry in scope `webullapi4go/brokerfd/events`, including
`/grpc.event.EventService/Subscribe` client spans and the
`event_stream_attempts` / `event_stream_attempt_duration` metrics. Failed
attempts record sanitized error text; cancellation records gRPC `Canceled` and
outcome `error` even though it is not an `OnError` callback. See
[Observability](observability.md#event-telemetry).

## Related

- [`brokerfd` reference](https://pkg.go.dev/github.com/shing1211/webullapi4go/brokerfd)
- [`brokerfd/events` reference](https://pkg.go.dev/github.com/shing1211/webullapi4go/brokerfd/events)
- [Getting Started](getting-started.md)
- [Authentication](authentication.md)
