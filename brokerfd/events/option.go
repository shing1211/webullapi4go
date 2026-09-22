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

// DefaultGRPCPort is the default port for gRPC connections (443).
const DefaultGRPCPort = 443

// DefaultDialTimeout is the default timeout for establishing gRPC connections (10 seconds).
const DefaultDialTimeout = 10 * time.Second

// DefaultReconnectBaseDelay is the default base delay for exponential backoff reconnection (1 second).
const DefaultReconnectBaseDelay = time.Second

// DefaultReconnectMaxDelay is the default maximum delay for exponential backoff reconnection (30 seconds).
const DefaultReconnectMaxDelay = 30 * time.Second

// Option configures a BrokerFD events client.
type Option func(*config)

// config holds the configuration for a BrokerFD events client.
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

// defaultConfig returns the default configuration for a BrokerFD events client.
func defaultConfig() config {
	return config{
		port:               DefaultGRPCPort,
		tls:                true,
		dialTimeout:        DefaultDialTimeout,
		autoReconnect:      true,
		reconnectBaseDelay: DefaultReconnectBaseDelay,
		reconnectMaxDelay:  DefaultReconnectMaxDelay,
	}
}

// validate checks that the configuration values are valid.
func (cfg config) validate() error {
	if cfg.port <= 0 || cfg.port > 65535 {
		return errs.New(errs.CodeInvalidConfig, "brokerfd/events: port must be between 1 and 65535")
	}
	if cfg.dialTimeout <= 0 {
		return errs.New(errs.CodeInvalidConfig, "brokerfd/events: dial timeout must be positive")
	}
	if cfg.reconnectBaseDelay <= 0 {
		return errs.New(errs.CodeInvalidConfig, "brokerfd/events: reconnect base delay must be positive")
	}
	if cfg.reconnectMaxDelay < cfg.reconnectBaseDelay {
		return errs.New(errs.CodeInvalidConfig, "brokerfd/events: reconnect max delay must not be less than the base delay")
	}
	if cfg.maxReconnectAttempts < 0 {
		return errs.New(errs.CodeInvalidConfig, "brokerfd/events: max reconnect attempts must not be negative")
	}
	return nil
}

// WithGRPCEndpoint sets the gRPC server host address.
func WithGRPCEndpoint(host string) Option {
	return func(cfg *config) { cfg.endpoint = host }
}

// WithGRPCPort sets the gRPC server port.
func WithGRPCPort(port int) Option {
	return func(cfg *config) {
		if port > 0 {
			cfg.port = port
		}
	}
}

// WithTLS enables or disables TLS for the gRPC connection.
func WithTLS(enabled bool) Option {
	return func(cfg *config) { cfg.tls = enabled }
}

// WithDialTimeout sets the timeout for establishing gRPC connections.
func WithDialTimeout(d time.Duration) Option {
	return func(cfg *config) {
		if d > 0 {
			cfg.dialTimeout = d
		}
	}
}

// WithGRPCDialOption appends a gRPC dial option to the connection.
func WithGRPCDialOption(opt grpc.DialOption) Option {
	return func(cfg *config) {
		if opt != nil {
			cfg.dialOptions = append(cfg.dialOptions, opt)
		}
	}
}

// WithAccounts sets the account IDs to filter events.
func WithAccounts(accounts []string) Option {
	return func(cfg *config) {
		cfg.accounts = append([]string(nil), accounts...)
	}
}

// WithAutoReconnect enables or disables automatic reconnection on connection loss.
func WithAutoReconnect(enabled bool) Option {
	return func(cfg *config) { cfg.autoReconnect = enabled }
}

// WithReconnectBaseDelay sets the base delay for exponential backoff reconnection.
func WithReconnectBaseDelay(d time.Duration) Option {
	return func(cfg *config) {
		if d > 0 {
			cfg.reconnectBaseDelay = d
		}
	}
}

// WithReconnectMaxDelay sets the maximum delay for exponential backoff reconnection.
func WithReconnectMaxDelay(d time.Duration) Option {
	return func(cfg *config) {
		if d > 0 {
			cfg.reconnectMaxDelay = d
		}
	}
}

// WithMaxReconnectAttempts sets the maximum number of reconnection attempts.
// A value of 0 means unlimited attempts.
func WithMaxReconnectAttempts(n int) Option {
	return func(cfg *config) {
		if n > 0 {
			cfg.maxReconnectAttempts = n
			return
		}
		cfg.maxReconnectAttempts = 0
	}
}
