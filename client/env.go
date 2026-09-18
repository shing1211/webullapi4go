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
	"os"
	"strings"
)

// Environment variables read by [WithEnv]. They mirror the names the SDK's
// integration tests use, so the same shell environment configures both.
const (
	// envAppKey is the Webull OpenAPI app key.
	envAppKey = "WEBULL_APP_KEY"
	// envAppSecret is the Webull OpenAPI app secret.
	envAppSecret = "WEBULL_APP_SECRET"
	// envRegion is the deployment region ("hk", "us", "jp", ...).
	envRegion = "WEBULL_REGION"
	// envEnvironment is the deployment environment ("prod", "production",
	// "uat", or "sandbox").
	envEnvironment = "WEBULL_ENVIRONMENT"
	// envBaseURL overrides only the REST base URL.
	envBaseURL = "WEBULL_BASE_URL"
	// envMQTTURL overrides only the MQTT broker address.
	envMQTTURL = "WEBULL_MQTT_URL"
)

// WithEnv returns an [Option] that fills the configuration from the process
// environment. It is a convenience for deployment: the same variables the
// sandbox integration tests read configure a production client.
//
// The variables read are:
//
//   - WEBULL_APP_KEY, WEBULL_APP_SECRET — credentials.
//   - WEBULL_REGION — deployment region, matched case-insensitively against
//     the known regions (for example "hk" or "us").
//   - WEBULL_ENVIRONMENT — "prod" or "production" selects [Production];
//     "uat" or "sandbox" selects [Sandbox].
//   - WEBULL_BASE_URL — optional override of the REST base URL only.
//   - WEBULL_MQTT_URL — optional override of the MQTT broker address only.
//
// Variables that are unset or empty are ignored. An invalid WEBULL_REGION or
// WEBULL_ENVIRONMENT value is never silently ignored: [New] fails with a typed
// *[errs.Error] whose code is [errs.CodeInvalidConfig] (the same code returned
// when required credentials are missing).
//
// Precedence: WithEnv acts as a lower-precedence default source. It fills a
// field only when that field has not already been set by another option, so
// explicit functional options always win over the environment regardless of
// the order in which the options are passed to [New]:
//
//	// The explicit credentials win; region and environment come from the
//	// environment because no explicit option sets them.
//	client.New(client.WithCredentials(key, secret), client.WithEnv())
//
// Applying WithEnv() before or after an explicit option is therefore
// equivalent. Similarly, an endpoint option ([WithEndpoints] or
// [WithBaseURL]) suppresses the WEBULL_BASE_URL and WEBULL_MQTT_URL overrides.
//
// The environment is read when the option is applied inside [New], not when
// WithEnv is called.
func WithEnv() Option {
	return func(c *Config) { applyEnv(c) }
}

// applyEnv fills the environment-overridable fields of c that were not set by
// an explicit option. It is called by [WithEnv] and keeps the precedence rule
// in one place.
func applyEnv(c *Config) {
	if v, ok := envValue(envAppKey); ok && !c.appKeySet {
		c.AppKey = v
	}
	if v, ok := envValue(envAppSecret); ok && !c.appSecretSet {
		c.AppSecret = v
	}
	if v, ok := envValue(envRegion); ok && !c.regionSet {
		c.Region = normalizeRegion(v)
	}
	if v, ok := envValue(envEnvironment); ok && !c.environmentSet {
		c.Environment = normalizeEnvironment(v)
	}
	if c.endpointsOverride {
		return
	}
	baseURL, hasBaseURL := envValue(envBaseURL)
	mqttURL, hasMQTTURL := envValue(envMQTTURL)
	if !hasBaseURL && !hasMQTTURL {
		return
	}
	// Recompute the endpoints from the effective region and environment so
	// that, for example, WEBULL_REGION=us with WEBULL_MQTT_URL set uses the US
	// MQTT host as the base when the variable is applied.
	endpoints := EndpointsFor(c.Region, c.Environment)
	if hasBaseURL {
		endpoints.HTTP = baseURL
	}
	if hasMQTTURL {
		endpoints.MQTT = mqttURL
	}
	c.Endpoints = endpoints
	c.endpointsOverride = true
}

// envValue returns the trimmed value of the named environment variable and
// reports whether it was present and non-empty.
func envValue(name string) (string, bool) {
	v, ok := os.LookupEnv(name)
	if !ok {
		return "", false
	}
	v = strings.TrimSpace(v)
	if v == "" {
		return "", false
	}
	return v, true
}

// normalizeRegion lowercases a WEBULL_REGION value. An unknown value is
// preserved so that [Config.Validate] reports it as an error rather than the
// client silently falling back to a default region.
func normalizeRegion(v string) Region {
	return Region(strings.ToLower(v))
}

// normalizeEnvironment maps a WEBULL_ENVIRONMENT value to an [Environment].
// The aliases "prod" and "production" select [Production]; "uat" and
// "sandbox" select [Sandbox]. Any other value is preserved so that
// [Config.Validate] reports it as an error.
func normalizeEnvironment(v string) Environment {
	switch strings.ToLower(v) {
	case "prod", "production":
		return Production
	case "uat", "sandbox":
		return Sandbox
	default:
		return Environment(strings.ToLower(v))
	}
}
