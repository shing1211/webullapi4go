# futures-probe

Probes HK futures discovery and market data: instrument list, product codes,
product classes, snapshot, and bars. Own Go module. Off by default.

## Run

```sh
cd examples/futures-probe
export WEBULL_ENVIRONMENT="sandbox"
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
WEBULL_FUTURES_TEST=1 go run .
```

## Notes

- The gate (`WEBULL_FUTURES_TEST=1`) must be set or the program exits without
  doing anything. Credentials — see [../README.md](../README.md).
- Read-only probe against live sandbox endpoints.
