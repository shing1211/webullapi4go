# Options Example

This example demonstrates a **preview-only** multi-leg US options order.

It builds a vertical call spread (buy the lower-strike call, sell the higher-strike call on the same expiration) and calls `PreviewOrder` to estimate its cost. No order is ever placed.

```text
AAPL vertical call spread: long 220 call / short 230 call
```

## Credentials

```sh
export WEBULL_APP_KEY=...
export WEBULL_APP_SECRET=...
export WEBULL_ACCOUNT_ID=...   # optional; uses first account if unset
export WEBULL_ENVIRONMENT=sandbox
```

## Run

```sh
go run ./examples/options
```

## Notes

- US sandbox market data is limited to `AAPL`.
- Option contracts for `AAPL` may not exist in the sandbox (`417 Invalid Symbol`).
- This is **preview only** — it never submits an order to the account.
