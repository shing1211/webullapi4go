# Webull Broker FD gRPC Event Streaming

Subscribes to Webull Broker FD events over a gRPC stream and prints received
data events until interrupted with Ctrl+C.

## Usage

```sh
go run ./examples/brokerfd-events
```

## Credentials

Requires `WEBULL_APP_KEY` and `WEBULL_APP_SECRET`. If `WEBULL_APP_KEY` is not
set, the program logs a message and exits gracefully without error.

Optionally set `WEBULL_TRADE_ACCOUNT_ID` to filter events to a specific
account.

## Output

The program logs:

- `"connected"` when the subscription is established
- `"ping"` upon receiving a keep-alive ping from the server
- `"event stream error: ..."` if an error occurs on the stream
- `"data: type=N content=... payload=..."` for each received data event

## Graceful Shutdown

Send SIGINT (Ctrl+C) or SIGTERM to exit cleanly. The underlying gRPC connection
is closed before the program exits.
