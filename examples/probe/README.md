# probe

Sandbox endpoint probe (Go and Python equivalents): hits live sandbox endpoints
to verify paths and response schemas. Not a user-facing example — a development
tool used while building the SDK.

## Run

```sh
export WEBULL_ENVIRONMENT="sandbox"
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
go run ./examples/probe

# or the Python equivalent
python examples/probe/probe.py
```

## Notes

- Results are written to `examples/probe/results/` as JSON files
  (`{category}_{endpoint}.json`).
- Shared HK sandbox test accounts:
  <https://developer.webull.hk/apis/docs/sdk#test-accounts>
- Rate limits: token endpoint 10 req / 30 s (the script enforces a 4 s
  cooldown); 1 s between general requests.
- To add targets, extend the `tests` list in each `probe_*` function; a `200`
  response stops alternate-path attempts for that base path.
