# order

Previews a small AAPL limit buy and, only when `WEBULL_ORDER_PLACE=1` is set,
places it far below the market and immediately cancels it. Previewing mutates
nothing; placing is off by default.

## Run

```sh
# Preview only (mutates nothing)
go run ./examples/order

# Place and immediately cancel a non-marketable limit order (mutating)
WEBULL_ORDER_PLACE=1 go run ./examples/order
```

```powershell
$env:WEBULL_ORDER_PLACE = "1"
go run ./examples/order
```

## Notes

- Use a sandbox account only. `WEBULL_ACCOUNT_ID` is optional (the first
  account is used when unset). Credentials — see [../README.md](../README.md).
- The example sets `WithMaxOrderQuantity("10")` and
  `WithMaxOrderNotional("2500.00")` guardrails, enforced by `PreviewOrder` and
  `PlaceOrder` before any network call.
- It never sends a market order.
