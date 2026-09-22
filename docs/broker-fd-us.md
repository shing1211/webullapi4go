# Broker FD API US

The `brokerfd` package provides the Webull Broker FD (US) HTTP API. It shares the
core `client.Client` for signing and transport. The `brokerfd/events` sub-package
provides the corresponding gRPC event stream.

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

conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
svc := eventsevents.NewEventServiceClient(conn)
ec := events.NewClient(svc)

req := events.NewSubscribeRequest(subscribeType, accounts)
client, err := ec.Subscribe(ctx, req)
for {
    resp, err := client.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        return err
    }
    // process resp.GetEventType(), resp.GetPayload() ...
}
```

Event types: `SubscribeSuccess`, `Ping`, `AuthError`, `NumOfConnExceed`, `SubscribeExpired`.

## Related

- [`brokerfd` reference](https://pkg.go.dev/github.com/shing1211/webullapi4go/brokerfd)
- [`brokerfd/events` reference](https://pkg.go.dev/github.com/shing1211/webullapi4go/brokerfd/events)
- [Getting Started](getting-started.md)
- [Authentication](authentication.md)
