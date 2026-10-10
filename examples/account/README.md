# account

Lists the authenticated user's trading accounts and, for one account, prints its
balance and open positions. Read-only: it never places, replaces, or cancels
orders.

## Run

```sh
go run ./examples/account
```

## Notes

- Credentials come from the environment via
  [`client.WithEnv()`](https://pkg.go.dev/github.com/shing1211/webullapi4go/client#WithEnv)
  — see [../README.md](../README.md).
- Uses `WEBULL_ACCOUNT_ID` when set; otherwise the first account returned.
- Trading requests need an access token, which the example obtains with
  `EnsureToken`.
