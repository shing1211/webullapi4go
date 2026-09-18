# Webull API Live Probe

Authenticates against the Webull HK sandbox and probes endpoint paths + schemas
for v0.5 best-effort stubs (Fund Data, Crypto Data, Screener v2, Broker FD).

## Usage

```sh
python examples/probe/probe.py
```

Results are written to `examples/probe/results/` as JSON files.

## Credentials

Uses the shared HK sandbox test account hardcoded in `probe.py`.
Override with environment variables:

```sh
WEBULL_APP_KEY=xxx WEBULL_APP_SECRET=yyy python examples/probe/probe.py
```

Credentials: https://developer.webull.hk/apis/docs/sdk#test-accounts

## Output

Each probe result is saved as `{category}_{endpoint}.json` with structure:

```json
{
  "path": "/market-data/fund/AAPL/nav",
  "query": {"page_size": "5"},
  "result": {
    "status": 200,
    "data": { ... }
  }
}
```

## Rate limits

- Token endpoint: max 10 req / 30s — script enforces 4s cooldown
- General: 1s between requests

## Adding new probe targets

Add entries to the `tests` list in each `probe_*` function. If a path returns 200,
the function stops trying alternates for that base path.
