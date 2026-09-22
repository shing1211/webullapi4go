# options

Preview-only multi-leg US options order: builds an AAPL vertical call spread
(long 220 call / short 230 call) and calls `PreviewOrder` to estimate its cost.
No order is ever placed.

## Run

```sh
go run ./examples/options
```

## Notes

- Credentials — see [../README.md](../README.md); `WEBULL_ACCOUNT_ID` is optional
  (the first account is used when unset).
- The HK sandbox accepts only `SINGLE` strategies (`417` for multi-leg) and may
  not have option contracts for `AAPL` (`417 Invalid Symbol`) — a US sandbox is
  needed to exercise this end to end.
- This is preview only; it never submits an order.
