# Broker API HK

The `broker` package provides the Webull Broker API for Hong Kong. It is a
separate Go module that shares the core `client.Client` for signing and transport.

!!! note "Prerequisites"
    - A [Webull account](https://developer.webull.hk/apis/docs/sdk#test-accounts) (sandbox or production)
    - Go 1.26+
    - An authenticated client — see [Authentication](authentication.md)

## Install

`broker` is a separate Go module. Publication of a Go-semver-compatible v2
module is deferred, so use a checkout of this repository and the local module
replacement rather than assuming that an unqualified `go get` installs the
current tree.
The `broker/go.mod` file contains
`replace github.com/shing1211/webullapi4go => ../`.

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

bk := broker.New(cl)
```

## Coverage

| Group | Methods |
|-------|---------|
| Virtual accounts | `CreateVirtualAccount`, `UpdateVirtualAccount`, `GetVirtualAccount`, `ListVirtualAccounts` |
| Instruments | `GetStockInstruments`, `GetStockLocate`, `GetCorporateActionsDetail` |
| Assets | `GetBalance`, `GetPositions` |
| Orders | `PreviewOrder`, `PlaceOrder`, `ReplaceOrder`, `CancelOrder`, `GetOrderDetail`, `GetOrderHistory`, `GetOpenOrders` |
| Cash activities | `GetCashActivities` |
| Funding FX | `GetFXRate`, `CreateFXExchange`, `GetFXExchangeDetail`, `CreateInstantExchange`, `GetInstantExchangeDetail` |
| Instant funding | `CreateInstantFunding`, `GetInstantFundingDetail` |
| Journals | `CreateCashJournal`, `GetCashJournalDetail`, `CreatePositionJournal`, `GetPositionJournalDetail` |
| Master data | `GetTradeCalendar` |
| Event contracts | `GetEventContractCategories`, `GetEventContractSeries`, `GetEventContractEvents`, `GetEventContractInstruments` |

## Hosts

| Environment | Host |
|-------------|------|
| Production (HK) | `https://broker-api.webull.hk` |
| Sandbox (HK) | `https://broker-api.sandbox.webull.hk` |

## Related

- [`broker` reference](https://pkg.go.dev/github.com/shing1211/webullapi4go/broker)
- [Getting Started](getting-started.md)
- [Authentication](authentication.md)
