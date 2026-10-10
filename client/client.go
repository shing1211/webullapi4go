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

package client

import (
	"net/http"
	"sync"
	"time"

	errs "github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/pkg/observability"
	"github.com/shing1211/webullapi4go/pkg/transport"
)

// Client is the Webull OpenAPI HTTP client. It owns the resolved [Config] and
// the shared transport, and is the SDK's public entry point: API-specific
// clients (for example the market-data client in package data) are constructed
// from it.
//
// A Client is safe for concurrent use. All network access flows through
// [Client.Do], [Client.DoBroker], or [Client.DoStream].
type Client struct {
	cfg       *Config
	transport *transport.Client
	brokerTr  *transport.Client

	// tok is this client's token cache, polling configuration, and injection
	// installation state. It lives in the struct rather than in a package-level
	// map keyed by *Client so that a freed client's state can never be observed
	// by a later client allocated at the same address. The zero value is usable:
	// unset polling parameters fall back to their defaults.
	tok tokenState

	// autoTokenMu serialises automatic token acquisition for this client so that
	// concurrent first requests do not create several tokens.
	autoTokenMu sync.Mutex

	// clockOffset and clockOffsetMu implement clock-drift correction: when
	// enabled, the offset between the local clock and the server's clock is
	// learned from the Date response header and applied to subsequent signing
	// timestamps.
	clockOffset   time.Duration
	clockOffsetMu sync.Mutex
}

// New returns a Client configured by opts. Options are applied in order on top
// of [DefaultConfig]; nil options are ignored. New validates the resulting
// configuration and returns an error when the credentials are missing or the
// region, environment, or endpoint is invalid.
func New(opts ...Option) (*Client, error) {
	cfg := DefaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	cfg.applyResiliencePreset()
	// Endpoints default to production Hong Kong; recompute them from the
	// resolved region and environment unless the caller overrode them.
	if !cfg.endpointsOverride {
		cfg.Endpoints = EndpointsFor(cfg.Region, cfg.Environment)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: cfg.Timeout, Transport: defaultHTTPTransport()}
	}
	if cfg.httpTransport != nil {
		cfg.HTTPClient.Transport = cfg.httpTransport
	}
	t, err := transport.New(cfg.Endpoints.HTTP, cfg.HTTPClient, cfg.UserAgent)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInvalidConfig, "invalid HTTP endpoint", err)
	}
	c := &Client{cfg: &cfg, transport: t}
	if cfg.Endpoints.BrokerHTTP != "" {
		c.brokerTr, _ = transport.New(cfg.Endpoints.BrokerHTTP, cfg.HTTPClient, cfg.UserAgent)
	}
	return c, nil
}

func defaultHTTPTransport() http.RoundTripper {
	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return http.DefaultTransport
	}
	tuned := transport.Clone()
	tuned.MaxIdleConns = 100
	tuned.MaxIdleConnsPerHost = 10
	tuned.IdleConnTimeout = 90 * time.Second
	tuned.TLSHandshakeTimeout = 10 * time.Second
	tuned.ExpectContinueTimeout = time.Second
	return tuned
}

// Config returns a copy of the client configuration.
func (c *Client) Config() Config { return *c.cfg }

// ObservabilityConfig returns the telemetry configuration shared by this client
// and service clients created from it. The returned value must be treated as
// read-only after the client is constructed.
func (c *Client) ObservabilityConfig() *observability.Config {
	if c == nil {
		return nil
	}
	return c.cfg.otel
}

// Region returns the configured region.
func (c *Client) Region() Region { return c.cfg.Region }

// Environment returns the configured environment.
func (c *Client) Environment() Environment { return c.cfg.Environment }

// Endpoints returns the resolved service endpoints.
func (c *Client) Endpoints() Endpoints { return c.cfg.Endpoints }

// HTTPClient returns the HTTP client used for REST calls.
func (c *Client) HTTPClient() *http.Client { return c.cfg.HTTPClient }

// AppKey returns the configured Webull app key.
func (c *Client) AppKey() string { return c.cfg.AppKey }

// AppSecret returns the configured Webull app secret.
func (c *Client) AppSecret() string { return c.cfg.AppSecret }

// Close releases resources owned by the client, in particular the idle
// connections held by the underlying HTTP transport. It does not close
// connections that are still in use, and may be called more than once.
func (c *Client) Close() error {
	c.transport.CloseIdleConnections()
	return nil
}
