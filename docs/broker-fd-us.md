# Broker FD API US

The `brokerfd` package provides the Webull Broker FD (US) HTTP API. It shares the
core `client.Client` for signing and transport. The `brokerfd/events` sub-package
provides the corresponding gRPC event stream.

!!! note "Prerequisites"
    - A [Webull account](https://developer.webull.hk/apis/docs/sdk#test-accounts) (sandbox or production)
    - Go 1.26+
    - An authenticated client — see [Authentication](authentication.md)

## Construct

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

```go
import "github.com/shing1211/webullapi4go/brokerfd/events"

ec, err := events.New(cl)
if err != nil {
    log.Fatal(err)
}
defer ec.Close()

ec.OnConnect(func() { log.Println("connected") })
ec.OnError(func(err error) { log.Printf("error: %v", err) })
ec.OnData(func(subscribeType uint32, contentType string, payload []byte) {
    log.Printf("event type=%d content=%s", subscribeType, contentType)
})

if err := ec.Run(ctx); err != nil {
    log.Fatal(err)
}
```

Event types: `SubscribeSuccess`, `Ping`, `AuthError`, `NumOfConnExceed`, `SubscribeExpired`.

## Related

- [`brokerfd` reference](https://pkg.go.dev/github.com/shing1211/webullapi4go/brokerfd)
- [`brokerfd/events` reference](https://pkg.go.dev/github.com/shing1211/webullapi4go/brokerfd/events)
- [Getting Started](getting-started.md)
- [Authentication](authentication.md)
