# streaming

Connects to the streaming broker over MQTT-over-WebSocket, subscribes to `AAPL`
`QUOTE`, `SNAPSHOT`, and `TICK` pushes, prints each message, and shuts down on
Ctrl+C (or `SIGTERM`).

## Run

```sh
go run ./examples/streaming
```

## Notes

- Credentials — see [../README.md](../README.md); the example obtains an access
  token with `EnsureToken`.
- WebSocket is used because plain MQTT on port `1883` is blocked on some
  networks; the example targets `wss://data-api.sandbox.webull.hk:8883/mqtt` in
  the sandbox.
- Sandbox market data is limited to `AAPL`.
