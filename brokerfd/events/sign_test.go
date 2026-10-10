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
	"testing"
	"time"
)

func TestSubscribeMetadataHasRequiredKeys(t *testing.T) {
	md, err := subscribeMetadata("app-key", "app-secret", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), []byte(`{}`))
	if err != nil {
		t.Fatalf("subscribeMetadata returned error: %v", err)
	}

	expectedKeys := []string{
		"x-app-key",
		"x-signature-algorithm",
		"x-signature-version",
		"x-signature-nonce",
		"x-timestamp",
		"x-signature",
	}
	for _, key := range expectedKeys {
		if vals := md.Get(key); len(vals) == 0 {
			t.Errorf("metadata missing key %q", key)
		}
	}
}

func TestEventSignatureNonEmpty(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	md, err := subscribeMetadata("app-key", "app-secret", now, []byte(`{}`))
	if err != nil {
		t.Fatalf("subscribeMetadata failed: %v", err)
	}
	sig := md.Get("x-signature")
	if len(sig) != 1 || sig[0] == "" {
		t.Errorf("signature is empty or missing: %v", sig)
	}
}

func TestEventSignatureDifferentForDifferentSecrets(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	body := []byte(`{}`)

	md1, _ := subscribeMetadata("app-key", "secret1", now, body)
	md2, _ := subscribeMetadata("app-key", "secret2", now, body)

	if len(md1.Get("x-signature")) == 1 && len(md2.Get("x-signature")) == 1 && md1.Get("x-signature")[0] == md2.Get("x-signature")[0] {
		t.Error("different secrets produced same signature")
	}
}

func TestEventSignatureDifferentForDifferentBodies(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	md1, _ := subscribeMetadata("app-key", "app-secret", now, []byte(`{}`))
	md2, _ := subscribeMetadata("app-key", "app-secret", now, []byte(`{"a":1}`))

	if len(md1.Get("x-signature")) == 1 && len(md2.Get("x-signature")) == 1 && md1.Get("x-signature")[0] == md2.Get("x-signature")[0] {
		t.Error("different bodies produced same signature")
	}
}

func TestPercentEncodeUnreserved(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"abcDEF123", "abcDEF123"},
		{"-_.~", "-_.~"},
		{"%", "%25"},
		{"=", "%3D"},
		{"&", "%26"},
		{"key=value", "key%3Dvalue"},
	}

	for _, tc := range tests {
		got := percentEncode(tc.input)
		if got != tc.expected {
			t.Errorf("percentEncode(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestIsUnreserved(t *testing.T) {
	valid := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.~"
	for i := 0; i < 256; i++ {
		c := byte(i)
		isValid := isUnreserved(c)
		isContained := len(valid) > 0 && containsByte(valid, c)
		if isValid != isContained {
			t.Errorf("isUnreserved(%d/%c) = %v, want %v", i, c, isValid, isContained)
		}
	}
}

func containsByte(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}
