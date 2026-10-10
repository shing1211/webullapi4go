# marketdata

Fetches a snapshot and recent daily bars for `AAPL` on the US market, then
prints the fields. Read-only.

## Run

```sh
go run ./examples/marketdata
```

## Notes

- Credentials — see [../README.md](../README.md).
- The sandbox serves a limited symbol set (currently `AAPL`). Bars are requested
  through the batch endpoint because the historical single-symbol endpoint has
  been retired.
