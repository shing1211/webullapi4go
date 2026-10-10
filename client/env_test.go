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

package client_test

import (
	"strings"
	"testing"

	"github.com/shing1211/webullapi4go/client"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
)

// envVarNames are the variables WithEnv reads. Tests clear all of them so the
// result never depends on the ambient shell.
var envVarNames = []string{
	"WEBULL_APP_KEY",
	"WEBULL_APP_SECRET",
	"WEBULL_REGION",
	"WEBULL_ENVIRONMENT",
	"WEBULL_BASE_URL",
	"WEBULL_MQTT_URL",
}

// clearWebullEnv blanks every variable WithEnv reads. Setting a variable to the
// empty string is how t.Setenv removes a value: WithEnv treats empty as absent.
func clearWebullEnv(t *testing.T) {
	t.Helper()
	for _, name := range envVarNames {
		t.Setenv(name, "")
	}
}

// setEnv sets several environment variables at once.
func setEnv(t *testing.T, pairs map[string]string) {
	t.Helper()
	for name, value := range pairs {
		t.Setenv(name, value)
	}
}

func TestEnvMissingCredentialsFailsCleanly(t *testing.T) {
	clearWebullEnv(t)

	cl, err := client.New(client.WithEnv())
	if err == nil {
		_ = cl.Close()
		t.Fatal("client.New(WithEnv()) succeeded without credentials, want error")
	}
	if !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("client.New(WithEnv()) error = %v, want code %s", err, errs.CodeInvalidConfig)
	}
	if !strings.Contains(err.Error(), "app key is required") {
		t.Fatalf("client.New(WithEnv()) error = %v, want it to mention the missing app key", err)
	}
}

func TestEnvLoadsCredentials(t *testing.T) {
	clearWebullEnv(t)
	setEnv(t, map[string]string{
		"WEBULL_APP_KEY":    "env-key",
		"WEBULL_APP_SECRET": "env-secret",
	})

	cl, err := client.New(client.WithEnv())
	if err != nil {
		t.Fatalf("client.New(WithEnv()) error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	cfg := cl.Config()
	if cfg.AppKey != "env-key" {
		t.Errorf("AppKey = %q, want %q", cfg.AppKey, "env-key")
	}
	if cfg.AppSecret != "env-secret" {
		t.Errorf("AppSecret = %q, want %q", cfg.AppSecret, "env-secret")
	}
	// Region and environment are absent, so the SDK defaults apply.
	if cfg.Region != client.HK {
		t.Errorf("Region = %q, want default %q", cfg.Region, client.HK)
	}
	if cfg.Environment != client.Production {
		t.Errorf("Environment = %q, want default %q", cfg.Environment, client.Production)
	}
}

func TestEnvRegionAndEnvironment(t *testing.T) {
	tests := []struct {
		name    string
		region  string
		env     string
		wantReg client.Region
		wantEnv client.Environment
	}{
		{name: "lowercase hk uat", region: "hk", env: "uat", wantReg: client.HK, wantEnv: client.Sandbox},
		{name: "uppercase us prod", region: "US", env: "PROD", wantReg: client.US, wantEnv: client.Production},
		{name: "production alias", region: "jp", env: "production", wantReg: client.JP, wantEnv: client.Production},
		{name: "sandbox alias", region: "sg", env: "sandbox", wantReg: client.SG, wantEnv: client.Sandbox},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearWebullEnv(t)
			setEnv(t, map[string]string{
				"WEBULL_REGION":      tt.region,
				"WEBULL_ENVIRONMENT": tt.env,
			})

			cl, err := client.New(client.WithEnv(), client.WithCredentials("k", "s"))
			if err != nil {
				t.Fatalf("client.New() error = %v", err)
			}
			t.Cleanup(func() { _ = cl.Close() })

			if got := cl.Region(); got != tt.wantReg {
				t.Errorf("Region = %q, want %q", got, tt.wantReg)
			}
			if got := cl.Environment(); got != tt.wantEnv {
				t.Errorf("Environment = %q, want %q", got, tt.wantEnv)
			}
		})
	}
}

func TestEnvInvalidRegionFails(t *testing.T) {
	clearWebullEnv(t)
	setEnv(t, map[string]string{"WEBULL_REGION": "mars"})

	cl, err := client.New(client.WithEnv(), client.WithCredentials("k", "s"))
	if err == nil {
		_ = cl.Close()
		t.Fatal("client.New() succeeded with an invalid region, want error")
	}
	if !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("error = %v, want code %s", err, errs.CodeInvalidConfig)
	}
	if !strings.Contains(err.Error(), "unknown region") || !strings.Contains(err.Error(), "mars") {
		t.Fatalf("error = %v, want it to name the unknown region", err)
	}
}

func TestEnvInvalidEnvironmentFails(t *testing.T) {
	clearWebullEnv(t)
	setEnv(t, map[string]string{"WEBULL_ENVIRONMENT": "staging"})

	cl, err := client.New(client.WithEnv(), client.WithCredentials("k", "s"))
	if err == nil {
		_ = cl.Close()
		t.Fatal("client.New() succeeded with an invalid environment, want error")
	}
	if !errs.Is(err, errs.CodeInvalidConfig) {
		t.Fatalf("error = %v, want code %s", err, errs.CodeInvalidConfig)
	}
	if !strings.Contains(err.Error(), "unknown environment") || !strings.Contains(err.Error(), "staging") {
		t.Fatalf("error = %v, want it to name the unknown environment", err)
	}
}

func TestEnvEmptyValuesAreIgnored(t *testing.T) {
	clearWebullEnv(t)
	setEnv(t, map[string]string{
		"WEBULL_REGION":      "   ",
		"WEBULL_ENVIRONMENT": "",
	})

	// Blank region/environment must not be treated as invalid; the only
	// failure is the genuinely missing credentials.
	cl, err := client.New(client.WithEnv())
	if err == nil {
		_ = cl.Close()
		t.Fatal("client.New() succeeded without credentials, want error")
	}
	if !strings.Contains(err.Error(), "app key is required") {
		t.Fatalf("error = %v, want the missing app key error", err)
	}
}

func TestEnvBaseAndMQTTURL(t *testing.T) {
	clearWebullEnv(t)
	setEnv(t, map[string]string{
		"WEBULL_APP_KEY":    "env-key",
		"WEBULL_APP_SECRET": "env-secret",
		"WEBULL_BASE_URL":   "https://base.example",
		"WEBULL_MQTT_URL":   "mqtt.example:1883",
	})

	cl, err := client.New(client.WithEnv())
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	endpoints := cl.Endpoints()
	if endpoints.HTTP != "https://base.example" {
		t.Errorf("Endpoints.HTTP = %q, want %q", endpoints.HTTP, "https://base.example")
	}
	if endpoints.MQTT != "mqtt.example:1883" {
		t.Errorf("Endpoints.MQTT = %q, want %q", endpoints.MQTT, "mqtt.example:1883")
	}
}

// TestEnvPrecedence verifies the documented rule: an explicit option always
// wins over the environment, regardless of the order in which the options are
// applied.
func TestEnvPrecedence(t *testing.T) {
	env := map[string]string{
		"WEBULL_APP_KEY":     "env-key",
		"WEBULL_APP_SECRET":  "env-secret",
		"WEBULL_REGION":      "us",
		"WEBULL_ENVIRONMENT": "sandbox",
		"WEBULL_BASE_URL":    "https://env.example",
	}
	explicit := []client.Option{
		client.WithCredentials("explicit-key", "explicit-secret"),
		client.WithRegion(client.HK),
		client.WithEnvironment(client.Production),
		client.WithBaseURL("https://explicit.example"),
	}

	orders := map[string][]client.Option{
		"explicit after env":  append([]client.Option{client.WithEnv()}, explicit...),
		"explicit before env": append(append([]client.Option{}, explicit...), client.WithEnv()),
	}
	for name, opts := range orders {
		t.Run(name, func(t *testing.T) {
			clearWebullEnv(t)
			setEnv(t, env)

			cl, err := client.New(opts...)
			if err != nil {
				t.Fatalf("client.New() error = %v", err)
			}
			t.Cleanup(func() { _ = cl.Close() })

			cfg := cl.Config()
			if cfg.AppKey != "explicit-key" {
				t.Errorf("AppKey = %q, want explicit value", cfg.AppKey)
			}
			if cfg.AppSecret != "explicit-secret" {
				t.Errorf("AppSecret = %q, want explicit value", cfg.AppSecret)
			}
			if cfg.Region != client.HK {
				t.Errorf("Region = %q, want explicit %q", cfg.Region, client.HK)
			}
			if cfg.Environment != client.Production {
				t.Errorf("Environment = %q, want explicit %q", cfg.Environment, client.Production)
			}
			if cfg.Endpoints.HTTP != "https://explicit.example" {
				t.Errorf("Endpoints.HTTP = %q, want explicit base URL", cfg.Endpoints.HTTP)
			}
		})
	}
}

// TestEnvPartialPrecedence checks that environment fields not overridden by an
// explicit option still take effect when only some fields are set explicitly.
func TestEnvPartialPrecedence(t *testing.T) {
	clearWebullEnv(t)
	setEnv(t, map[string]string{
		"WEBULL_APP_KEY":    "env-key",
		"WEBULL_APP_SECRET": "env-secret",
		"WEBULL_REGION":     "us",
		"WEBULL_BASE_URL":   "https://env.example",
	})

	cl, err := client.New(client.WithAppKey("explicit-key"), client.WithEnv())
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	t.Cleanup(func() { _ = cl.Close() })

	cfg := cl.Config()
	if cfg.AppKey != "explicit-key" {
		t.Errorf("AppKey = %q, want explicit value", cfg.AppKey)
	}
	if cfg.AppSecret != "env-secret" {
		t.Errorf("AppSecret = %q, want env value %q", cfg.AppSecret, "env-secret")
	}
	if cfg.Region != client.US {
		t.Errorf("Region = %q, want env %q", cfg.Region, client.US)
	}
	if cfg.Endpoints.HTTP != "https://env.example" {
		t.Errorf("Endpoints.HTTP = %q, want env base URL", cfg.Endpoints.HTTP)
	}
}
