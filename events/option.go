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

package events

import (
	"time"

	"google.golang.org/grpc"

	"github.com/shing1211/webullapi4go/internal/errs"
)

// Default connection parameters.
const (
	// DefaultGRPCPort is the TLS port of the Webull event service.
	DefaultGRPCPort = 443
	// DefaultDialTimeout bounds how long [Client.Run] waits for the gRPC
	// connection to become ready before giving up.
	DefaultDialTimeout = 10 * time.Second
	// DefaultReconnectBaseDelay is the first delay before [Client.Run]
	// reconnects after a stream ends or fails.
	DefaultReconnectBaseDelay = time.Second
	// DefaultReconnectMaxDelay caps the exponential reconnect delay.
	DefaultReconnectMaxDelay = 30 * time.Second
)

// Option configures a [Client] during [New]. Options are applied in order on
// top of the defaults.
type Option func(*config)

// config is the resolved event-client configuration.
type config struct {
	endpoint       string
	port           int
	tls            bool
	dialTimeout    time.Duration
	dialOptions    []grpc.DialOption
	subscribeTypes SubscribeType
	accounts       []string

	autoReconnect        bool
	reconnectBaseDelay   time.Duration
	reconnectMaxDelay    time.Duration
	maxReconnectAttempts int
}

// defaultConfig returns the event-client defaults: TLS on port 443, a
// ten-second dial timeout, every event category subscribed, and automatic
// reconnect with exponential backoff.
func defaultConfig() config {
	return config{
		port:               DefaultGRPCPort,
		tls:                true,
		dialTimeout:        DefaultDialTimeout,
		subscribeTypes:     SubscribeAll,
		autoReconnect:      true,
		reconnectBaseDelay: DefaultReconnectBaseDelay,
		reconnectMaxDelay:  DefaultReconnectMaxDelay,
	}
}

// validate reports whether cfg is usable.
func (cfg config) validate() error {
	if cfg.port <= 0 || cfg.port > 65535 {
		return errs.New(errs.CodeInvalidConfig, "events: port must be between 1 and 65535")
	}
	if cfg.dialTimeout <= 0 {
		return errs.New(errs.CodeInvalidConfig, "events: dial timeout must be positive")
	}
	if cfg.subscribeTypes == 0 {
		return errs.New(errs.CodeInvalidConfig, "events: at least one subscribe type is required")
	}
	if cfg.reconnectBaseDelay <= 0 {
		return errs.New(errs.CodeInvalidConfig, "events: reconnect base delay must be positive")
	}
	if cfg.reconnectMaxDelay < cfg.reconnectBaseDelay {
		return errs.New(errs.CodeInvalidConfig, "events: reconnect max delay must not be less than the base delay")
	}
	if cfg.maxReconnectAttempts < 0 {
		return errs.New(errs.CodeInvalidConfig, "events: max reconnect attempts must not be negative")
	}
	return nil
}

// WithGRPCEndpoint overrides the event-service host resolved from the client
// endpoints. It accepts a bare host (for example
// "events-api.sandbox.webull.hk") or a full gRPC target that carries a resolver
// scheme (for example "passthrough:///bufnet"); the port is not included.
func WithGRPCEndpoint(host string) Option {
	return func(cfg *config) { cfg.endpoint = host }
}

// WithGRPCPort overrides the event-service port. Non-positive values are
// ignored, keeping the default ([DefaultGRPCPort]).
func WithGRPCPort(port int) Option {
	return func(cfg *config) {
		if port > 0 {
			cfg.port = port
		}
	}
}

// WithTLS selects whether the connection uses TLS. It defaults to true; setting
// it to false uses plaintext credentials and is intended for tests and
// unencrypted local servers.
func WithTLS(enabled bool) Option {
	return func(cfg *config) { cfg.tls = enabled }
}

// WithDialTimeout sets how long [Client.Run] waits for the connection to become
// ready. Non-positive values are ignored, keeping the default
// ([DefaultDialTimeout]).
func WithDialTimeout(d time.Duration) Option {
	return func(cfg *config) {
		if d > 0 {
			cfg.dialTimeout = d
		}
	}
}

// WithGRPCDialOption appends a dial option, for example a custom context dialer
// or keepalive parameters. Options are passed to the gRPC client after the
// transport credentials derived from [WithTLS]. A nil option is ignored.
func WithGRPCDialOption(opt grpc.DialOption) Option {
	return func(cfg *config) {
		if opt != nil {
			cfg.dialOptions = append(cfg.dialOptions, opt)
		}
	}
}

// WithSubscribeTypes restricts the subscription to the given categories. The
// values are OR-ed together, so both [SubscribeOrder] | [SubscribePosition] and
// the variadic form select the same set. Calling it with no values leaves the
// default ([SubscribeAll]) unchanged.
func WithSubscribeTypes(types ...SubscribeType) Option {
	return func(cfg *config) {
		if len(types) == 0 {
			return
		}
		var combined SubscribeType
		for _, t := range types {
			combined |= t
		}
		cfg.subscribeTypes = combined
	}
}

// WithAccounts sets the trading account ids to subscribe for. An empty list
// subscribes without an account filter, which the server accepts for
// application-level events. The slice is copied, so later caller mutations do
// not affect the client.
func WithAccounts(accounts []string) Option {
	return func(cfg *config) {
		cfg.accounts = append([]string(nil), accounts...)
	}
}

// WithAutoReconnect selects whether [Client.Run] reconnects and re-subscribes
// after the stream ends or fails. It defaults to true. When enabled, Run
// retries transient failures and clean server ends with exponential backoff and
// jitter; terminal failures (authentication, permission, account, or
// configuration errors) end Run instead of looping. When disabled, Run returns
// after the first stream ends.
func WithAutoReconnect(enabled bool) Option {
	return func(cfg *config) { cfg.autoReconnect = enabled }
}

// WithReconnectBaseDelay sets the first delay before [Client.Run] reconnects.
// The delay doubles after each failed attempt up to [WithReconnectMaxDelay],
// with jitter applied. Non-positive values are ignored, keeping the default
// ([DefaultReconnectBaseDelay]).
func WithReconnectBaseDelay(d time.Duration) Option {
	return func(cfg *config) {
		if d > 0 {
			cfg.reconnectBaseDelay = d
		}
	}
}

// WithReconnectMaxDelay caps the reconnect delay. Non-positive values are
// ignored, keeping the default ([DefaultReconnectMaxDelay]). A value below the
// base delay is rejected by [New].
func WithReconnectMaxDelay(d time.Duration) Option {
	return func(cfg *config) {
		if d > 0 {
			cfg.reconnectMaxDelay = d
		}
	}
}

// WithMaxReconnectAttempts bounds the consecutive reconnect attempts [Client.Run]
// makes before it gives up and returns the last error. A non-positive value
// means unlimited retries, which is the default; negative values are treated as
// unlimited.
func WithMaxReconnectAttempts(n int) Option {
	return func(cfg *config) {
		if n > 0 {
			cfg.maxReconnectAttempts = n
			return
		}
		cfg.maxReconnectAttempts = 0
	}
}
