# watchlist-cmd

Watchlist CRUD example: create a watchlist, add instruments, update it, remove
instruments, and delete it. Mutating — gated by `WEBULL_WATCHLIST_TEST=1`. Own
Go module.

## Run

```sh
cd examples/watchlist-cmd
export WEBULL_ENVIRONMENT="sandbox"
export WEBULL_APP_KEY="your-sandbox-app-key"
export WEBULL_APP_SECRET="your-sandbox-app-secret"
WEBULL_WATCHLIST_TEST=1 go run .
```

## Notes

- The gate must be set or the program exits without doing anything, because the
  example mutates the user's watchlists. Credentials — see
  [../README.md](../README.md).
- Use a sandbox account; the example creates and then deletes its own watchlist.
