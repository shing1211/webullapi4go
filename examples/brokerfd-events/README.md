# brokerfd-events

Subscribes to Broker FD order, option, and position events over gRPC and prints
each decoded payload until interrupted with Ctrl+C (or SIGTERM).

## Run

```sh
go run ./examples/brokerfd-events
```

## Notes

- Credentials — see [../README.md](../README.md). Exits gracefully if
  `WEBULL_APP_KEY` is unset.
- Optionally set `WEBULL_TRADE_ACCOUNT_ID` to filter events to one account.
- Logs `connected`, `ping`, and `data: type=N ...` lines; the gRPC connection is
  closed on shutdown.
- US-only: the HK sandbox has no FD event stream.
