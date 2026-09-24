# account-monitor

Periodically prints a trading account's total balance and per-currency buying
power and initial margin. The monitor is read-only and runs until Ctrl+C (or
`SIGTERM`).

## Run

```sh
go run ./examples/account-monitor
```

The default interval is 30 seconds. Set a different interval of at least one
second when needed:

```sh
WEBULL_MONITOR_INTERVAL=10s go run ./examples/account-monitor
```

## Notes

- Credentials — see [../README.md](../README.md); the example obtains an access
  token with `EnsureToken`.
- `WEBULL_ACCOUNT_ID` is optional; when unset, the first account returned by the
  API is monitored.
- A single worker performs the read-only `GetBalance` requests. Ticks received
  while a request is active are skipped, so requests never overlap.
- Each request has a 15-second timeout. A successful poll resets the error
  count; three consecutive failures stop the monitor with an error.
- Ctrl+C or `SIGTERM` cancels the active request and waits for the worker to
  exit.
