# options-multi-leg

Probes every multi-leg option strategy value (VERTICAL, STRADDLE, STRANGLE,
IRON_CONDOR, IRON_BUTTERFLY, BUTTERFLY, CALENDAR, DIAGONAL, RATIO, COLLAR)
through the read-only `PreviewOrder`. No live orders are placed. Own Go module.
Off by default.

## Run

```sh
cd examples/options-multi-leg
export WEBULL_ENVIRONMENT="sandbox"
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
WEBULL_OPTIONS_TEST=1 go run .
```

## Notes

- The gate (`WEBULL_OPTIONS_TEST=1`) must be set or the program exits without
  doing anything. `WEBULL_ACCOUNT_ID` is optional. Credentials — see
  [../README.md](../README.md).
- The HK sandbox accepts only `SINGLE` (`417` for every other strategy); a US
  sandbox is needed to exercise the combo strategies.
