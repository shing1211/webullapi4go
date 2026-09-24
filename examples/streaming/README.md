# streaming

Connects to the streaming broker over MQTT-over-WebSocket, subscribes to `AAPL`
`QUOTE`, `SNAPSHOT`, and `TICK` pushes, and consumes each push from a typed
channel. The run is bounded to two minutes and also stops on Ctrl+C (or
`SIGTERM`).

## Run

```sh
go run ./examples/streaming
```

## Notes

- Credentials — see [../README.md](../README.md); the example obtains an access
  token with `EnsureToken`.
- Each channel has a 32-message buffer and uses `DropOldest`, so a slow console
  cannot apply unbounded backpressure to MQTT dispatch.
- Channel cancellation functions and `Close` release the subscriptions during
  shutdown.
- WebSocket is used because plain MQTT on port `1883` is blocked on some
  networks; the example targets `wss://data-api.sandbox.webull.hk:8883/mqtt` in
  the sandbox.
- Sandbox market data is limited to `AAPL`.
