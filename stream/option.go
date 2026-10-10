// Copyright 2026 shing1211
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stream

import (
	"crypto/tls"
	"time"

	"go.opentelemetry.io/otel/metric"

	errs "github.com/shing1211/webullapi4go/pkg/errors"
	mqtt "github.com/shing1211/webullapi4go/pkg/transport/mqtt"
)

// Default streaming parameters. Most mirror the low-level MQTT defaults.
const (
	// DefaultKeepAlive is the MQTT keep-alive interval.
	DefaultKeepAlive = mqtt.DefaultKeepAlive
	// DefaultConnectTimeout bounds a single MQTT connection attempt.
	DefaultConnectTimeout = mqtt.DefaultConnectTimeout
	// DefaultWriteTimeout bounds writing an MQTT control packet.
	DefaultWriteTimeout = mqtt.DefaultWriteTimeout
	// DefaultMessageChannelDepth is the inbound message buffer size.
	DefaultMessageChannelDepth = mqtt.DefaultMessageChannelDepth
	// DefaultMaxReconnectInterval caps the reconnect backoff.
	DefaultMaxReconnectInterval = mqtt.DefaultMaxReconnectInterval
	// DefaultResubscribeTimeout bounds the whole re-subscription sequence
	// issued after a reconnect.
	DefaultResubscribeTimeout = 30 * time.Second
)

// Option configures a [Client] during [New]. Options are applied in order.
type Option func(*config)

// config is the resolved streaming configuration.
type config struct {
	sessionID              string
	mqttURL                string
	websocket              bool
	keepAlive              time.Duration
	connectTimeout         time.Duration
	writeTimeout           time.Duration
	messageChannelDepth    uint
	cleanSession           bool
	autoReconnect          bool
	autoResubscribe        bool
	resubscribeTimeout     time.Duration
	healthWatchdogInterval time.Duration
	tlsConfig              *tls.Config
	meter                  metric.Meter
}

// defaultConfig returns the streaming defaults.
func defaultConfig() config {
	return config{
		keepAlive:           DefaultKeepAlive,
		connectTimeout:      DefaultConnectTimeout,
		writeTimeout:        DefaultWriteTimeout,
		messageChannelDepth: DefaultMessageChannelDepth,
		cleanSession:        true,
		autoResubscribe:     true,
		resubscribeTimeout:  DefaultResubscribeTimeout,
	}
}

// validate reports whether cfg is usable.
func (cfg config) validate() error {
	if cfg.sessionID == "" {
		return errs.New(errs.CodeInvalidConfig, "stream: session id is required")
	}
	if cfg.keepAlive <= 0 {
		return errs.New(errs.CodeInvalidConfig, "stream: keep-alive must be positive")
	}
	if cfg.connectTimeout <= 0 {
		return errs.New(errs.CodeInvalidConfig, "stream: connect timeout must be positive")
	}
	if cfg.writeTimeout <= 0 {
		return errs.New(errs.CodeInvalidConfig, "stream: write timeout must be positive")
	}
	if cfg.messageChannelDepth == 0 {
		return errs.New(errs.CodeInvalidConfig, "stream: message channel depth must be positive")
	}
	if cfg.autoResubscribe && cfg.resubscribeTimeout <= 0 {
		return errs.New(errs.CodeInvalidConfig, "stream: re-subscribe timeout must be positive")
	}
	return nil
}

// WithSessionID sets the session id used both as the MQTT client id and in the
// subscribe/unsubscribe calls. When unset, [New] generates a unique id.
//
// Do not reuse a session id across connections under one App Key: Webull
// disconnects the previous connection when a new one connects with the same
// session id. Reusing an id therefore silently drops another client's stream.
// The generated default is unique per [New] call; only override it when a
// stable id is required and connection exclusivity is guaranteed.
func WithSessionID(id string) Option {
	return func(cfg *config) { cfg.sessionID = id }
}

// WithClientID sets the MQTT client id independently of the session id. Webull
// requires the MQTT client id to equal the session id used in the
// subscribe/unsubscribe calls, so this is an alias for [WithSessionID]; it is
// provided for callers that express the connection in MQTT terms.
func WithClientID(id string) Option {
	return func(cfg *config) { cfg.sessionID = id }
}

// WithMQTTURL overrides the broker address resolved from the client endpoints.
// It accepts a paho URL (tcp://host:port, wss://host:port/path) or a bare
// host:port, which is treated as tcp://host:port.
func WithMQTTURL(rawURL string) Option {
	return func(cfg *config) { cfg.mqttURL = rawURL }
}

// WithWebSocket selects the MQTT-over-WebSocket endpoint instead of the plain
// TCP endpoint. It is ignored when [WithMQTTURL] is also set.
func WithWebSocket(useWebSocket bool) Option {
	return func(cfg *config) { cfg.websocket = useWebSocket }
}

// WithKeepAlive sets the MQTT keep-alive interval. Non-positive values are
// ignored by [New], which keeps the previous value.
func WithKeepAlive(d time.Duration) Option {
	return func(cfg *config) {
		if d > 0 {
			cfg.keepAlive = d
		}
	}
}

// WithConnectTimeout sets how long a single MQTT connection attempt may take.
// Non-positive values are ignored by [New], which keeps the previous value.
func WithConnectTimeout(d time.Duration) Option {
	return func(cfg *config) {
		if d > 0 {
			cfg.connectTimeout = d
		}
	}
}

// WithWriteTimeout sets the timeout for writing an MQTT control packet.
// Non-positive values are ignored by [New], which keeps the previous value.
func WithWriteTimeout(d time.Duration) Option {
	return func(cfg *config) {
		if d > 0 {
			cfg.writeTimeout = d
		}
	}
}

// WithMessageChannelDepth sets the number of inbound messages buffered by the
// MQTT client before traffic is throttled. A zero value is ignored by [New].
func WithMessageChannelDepth(depth uint) Option {
	return func(cfg *config) {
		if depth > 0 {
			cfg.messageChannelDepth = depth
		}
	}
}

// WithCleanSession requests a clean MQTT session. Streaming connections are
// stateless, so the default is true.
func WithCleanSession(clean bool) Option {
	return func(cfg *config) { cfg.cleanSession = clean }
}

// WithAutoReconnect enables the client's background reconnect loop. When a lost
// connection is re-established, the client automatically re-issues the active
// subscriptions (see [WithAutoResubscribe]) because Webull does not restore
// them. [Client.OnConnect] handlers fire again after a successful reconnect.
func WithAutoReconnect(auto bool) Option {
	return func(cfg *config) { cfg.autoReconnect = auto }
}

// WithAutoResubscribe controls whether the client re-issues its active HTTP
// subscriptions after a reconnect. It is enabled by default and is only
// meaningful together with [WithAutoReconnect]. Re-subscription is idempotent
// and never duplicates a subscription.
func WithAutoResubscribe(auto bool) Option {
	return func(cfg *config) { cfg.autoResubscribe = auto }
}

// WithResubscribeTimeout bounds the whole re-subscription sequence issued after
// a reconnect. Non-positive values are ignored by [New], which keeps the
// default ([DefaultResubscribeTimeout]).
func WithResubscribeTimeout(d time.Duration) Option {
	return func(cfg *config) {
		if d > 0 {
			cfg.resubscribeTimeout = d
		}
	}
}

// WithTLSConfig overrides the TLS configuration used for WebSocket (wss) or
// TLS connections. A nil value lets the MQTT client use a default
// configuration.
func WithTLSConfig(tc *tls.Config) Option {
	return func(cfg *config) { cfg.tlsConfig = tc }
}

// WithHealthWatchdog sets the interval at which the client checks message
// throughput. If no data message (quote, snapshot, or tick) arrives within this
// window, the connection transitions to [StateDegraded]; it recovers to
// [StateConnected] automatically when a message arrives. A zero interval
// (the default) disables the watchdog.
func WithHealthWatchdog(interval time.Duration) Option {
	return func(cfg *config) { cfg.healthWatchdogInterval = interval }
}

// WithMeter overrides the core client's OpenTelemetry meter for stream metrics
// (reconnects and channel drops). A nil value inherits the core client's meter,
// including its nil default; a non-nil value explicitly overrides it.
func WithMeter(m metric.Meter) Option {
	return func(cfg *config) { cfg.meter = m }
}
