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

package errs_test

import (
	"errors"
	"testing"

	"github.com/shing1211/webullapi4go/client"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/pkg/transport/mqtt"
)

func TestPublicSemanticSentinelsAreSpecific(t *testing.T) {
	genericTransport := errs.New(errs.CodeTransport, "dial failed")
	genericAuth := errs.New(errs.CodeAuth, "signing failed")

	if errors.Is(genericTransport, errs.ErrConnectionLimitExceeded) {
		t.Fatal("generic transport error matches ErrConnectionLimitExceeded")
	}
	if errors.Is(genericTransport, client.ErrCircuitOpen) {
		t.Fatal("generic transport error matches ErrCircuitOpen")
	}
	if errors.Is(genericTransport, mqtt.ErrConnectionLimit) {
		t.Fatal("generic transport error matches mqtt.ErrConnectionLimit")
	}
	if errors.Is(genericAuth, errs.ErrSubscriptionExpired) {
		t.Fatal("generic auth error matches ErrSubscriptionExpired")
	}
	if errors.Is(genericAuth, client.ErrAccessTokenRequired) {
		t.Fatal("generic auth error matches ErrAccessTokenRequired")
	}

	cases := []struct {
		name   string
		err    error
		target error
	}{
		{name: "connection limit", err: errs.ErrConnectionLimitExceeded, target: errs.ErrConnectionLimitExceeded},
		{name: "subscription expired", err: errs.ErrSubscriptionExpired, target: errs.ErrSubscriptionExpired},
		{name: "circuit open", err: client.ErrCircuitOpen, target: client.ErrCircuitOpen},
		{name: "access token required", err: client.ErrAccessTokenRequired, target: client.ErrAccessTokenRequired},
		{name: "mqtt limit", err: mqtt.ErrConnectionLimit, target: mqtt.ErrConnectionLimit},
		{name: "mqtt refused", err: mqtt.ErrConnectionRefused, target: mqtt.ErrConnectionRefused},
	}
	for _, tc := range cases {
		if !errors.Is(tc.err, tc.target) {
			t.Errorf("errors.Is(%s, sentinel) = false, want true", tc.name)
		}
		wrapped := errs.Wrap(errs.CodeTransport, "wrapped "+tc.name, tc.err)
		if !errors.Is(wrapped, tc.target) {
			t.Errorf("wrapped %s does not match its sentinel", tc.name)
		}
	}
}

func TestPublicSemanticSentinelsDoNotCrossMatch(t *testing.T) {
	if errors.Is(client.ErrCircuitOpen, mqtt.ErrConnectionLimit) {
		t.Fatal("ErrCircuitOpen matches mqtt.ErrConnectionLimit")
	}
	if errors.Is(mqtt.ErrConnectionRefused, mqtt.ErrConnectionLimit) {
		t.Fatal("mqtt.ErrConnectionRefused matches mqtt.ErrConnectionLimit")
	}
	if errors.Is(errs.ErrConnectionLimitExceeded, errs.ErrSubscriptionExpired) {
		t.Fatal("pkg/errors semantic sentinels cross-match")
	}
}

func TestPublicSemanticWrapperPreservesTypedError(t *testing.T) {
	err := errs.Wrap(errs.CodeTransport, "outer", client.ErrCircuitOpen)
	var typed *errs.Error
	if !errors.As(err, &typed) {
		t.Fatal("errors.As did not find the typed outer error")
	}
	if typed.Code != errs.CodeTransport {
		t.Fatalf("typed error code = %q, want %q", typed.Code, errs.CodeTransport)
	}
	if !errors.Is(err, client.ErrCircuitOpen) {
		t.Fatal("wrapped semantic cause was not preserved")
	}
}
