# path-probe

Compares the paths the SDK calls against the official Webull OpenAPI definition
for the endpoints where Webull's `llms.txt` summary and its own OpenAPI JSON
disagree (Known Issue 12 in `IMPLEMENTATION_STATUS.md`). Env-gated and skipped
by default.

## Run

```sh
export WEBULL_SANDBOX="1"
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
go run ./examples/path-probe
```

## Notes

- Requires sandbox credentials; without them the probe skips itself and exits.
- It probes both path variants and reports which one the sandbox accepts.
- Read-only: it only calls token and instrument endpoints.
