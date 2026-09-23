package data

import (
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m,
		// paho.mqtt.golang's WebSocket transport uses net/http's HTTP/2
		// client. After Disconnect the read-loop goroutines can persist
		// briefly; they are cleaned up asynchronously by the Go runtime
		// and are not a leak in our code.
		goleak.IgnoreAnyFunction("net/http.(*http2clientConnReadLoop).run"),
	)
}
