# events

Connects to the Webull gRPC trade-event stream, subscribes to order events, and
prints each decoded `OrderEvent` until Ctrl+C. Read-only.

## Run

```sh
go run ./examples/events
```

## Notes

- Each `Subscribe` is signed with HMAC-SHA256 — no access token is required.
- `WEBULL_TRADE_ACCOUNT_ID` scopes the subscription when set; the account must
  belong to the App Key. The stream reconnects and re-subscribes automatically
  after a transient drop.
- In the sandbox, placement events may not be pushed for a resting order; a
  cancellation produces the observed `CANCEL_SUCCESS` event.
- Credentials otherwise — see [../README.md](../README.md).
