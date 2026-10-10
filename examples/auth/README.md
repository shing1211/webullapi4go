# auth

Creates or reuses an access token and prints its status and expiry. Run this
first to confirm your credentials work.

## Run

```sh
go run ./examples/auth
```

## Notes

- Credentials come from the environment — see [../README.md](../README.md).
- In the sandbox the token is `NORMAL` immediately; in production the call waits
  for the Webull App 2FA verification window (up to five minutes by default).
