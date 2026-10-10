# brokerfd

Read-only Broker FD (US) endpoint probe: accounts summary, positions, and open
(unfilled) fractional orders for one account. It never places, replaces, or
cancels orders.

## Run

```sh
go run ./examples/brokerfd
```

## Notes

- Credentials — see [../README.md](../README.md). The program exits early with a
  graceful message if `WEBULL_APP_KEY` is not set.
- Set `WEBULL_TRADE_ACCOUNT_ID` for the open-orders call; otherwise that step
  is skipped.
- US-only: the HK sandbox returns `404` for Broker FD endpoints.
