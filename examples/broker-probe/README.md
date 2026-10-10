# broker-probe

Broker HK read-only endpoint probe: virtual account, instrument, asset, order,
cash activity, FX, and journal endpoints against the sandbox. Own Go module.

## Run

```sh
cd examples/broker-probe
export WEBULL_ENVIRONMENT="sandbox"
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
go run .
```

## Notes

- Credentials — see [../README.md](../README.md).
- The HK sandbox returns `401 ROUTE_NOT_PERMITTED` for `/broker/...` — the app
  lacks the required scope, not a path issue.
- Read-only: it never mutates broker state.
