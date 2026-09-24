# order

Previews a small AAPL limit buy and, only when `WEBULL_ORDER_PLACE=1` and a
caller-supplied `WEBULL_ORDER_IDEMPOTENCY_KEY` are set, places it far below the
market and immediately cancels it. Previewing mutates nothing; placing is off by
default and refuses to start without an idempotency key.

## Run

```sh
# Preview only (mutates nothing)
go run ./examples/order

# Choose and persist a new key before the first mutating request.
# Reuse this exact key for retries of the same logical order.
WEBULL_ORDER_IDEMPOTENCY_KEY=phase7-aapl-001 \
WEBULL_ORDER_PLACE=1 \
go run ./examples/order
```

```powershell
$env:WEBULL_ORDER_IDEMPOTENCY_KEY = "phase7-aapl-001"
$env:WEBULL_ORDER_PLACE = "1"
go run ./examples/order
```

## Notes

- Use a sandbox account only. `WEBULL_ACCOUNT_ID` is optional (the first
  account is used when unset). Credentials — see [../README.md](../README.md).
- The idempotency key becomes Webull's `client_order_id`; it must be at most 32
  characters using only letters, digits, `-`, and `_`. Persist it before the
  first request, reuse it after an unknown outcome, and use a new key for each
  intentionally new order.
- Preview-only runs use `trade.ClientOrderIDFrom` to derive a deterministic
  preview identifier. Mutating runs use the caller's key verbatim and fail
  before any network call when it is missing or invalid.
- The example sets `WithMaxOrderQuantity("10")` and
  `WithMaxOrderNotional("2500.00")` guardrails, enforced by `PreviewOrder` and
  `PlaceOrder` before any network call.
- After placement, cancellation uses an independent 30-second timeout so Ctrl+C
  cannot prevent cleanup. If cancellation fails, the error includes the key to
  use for manual reconciliation.
- It never sends a market order.
