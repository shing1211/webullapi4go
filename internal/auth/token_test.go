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

package auth

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTokenStatusValid(t *testing.T) {
	t.Parallel()

	valid := []TokenStatus{StatusPending, StatusNormal, StatusInvalid, StatusExpired}
	for _, s := range valid {
		if !s.Valid() {
			t.Errorf("TokenStatus(%q).Valid() = false, want true", s)
		}
		if got := s.String(); got != string(s) {
			t.Errorf("TokenStatus(%q).String() = %q, want %q", s, got, string(s))
		}
	}
	for _, s := range []TokenStatus{"", "normal", "UNKNOWN"} {
		if s.Valid() {
			t.Errorf("TokenStatus(%q).Valid() = true, want false", s)
		}
	}
}

func TestTokenUnmarshalJSON(t *testing.T) {
	t.Parallel()

	const expires = int64(1755486723000)
	cases := []struct {
		name string
		raw  string
		want Token
	}{
		{
			name: "documented expires_at",
			raw:  `{"token":"abc123","expires_at":1755486723000,"status":"NORMAL"}`,
			want: Token{Value: "abc123", ExpiresAt: expires, Status: StatusNormal},
		},
		{
			name: "legacy expires fallback",
			raw:  `{"token":"abc123","expires":1755486723000,"status":"PENDING"}`,
			want: Token{Value: "abc123", ExpiresAt: expires, Status: StatusPending},
		},
		{
			name: "expires_at wins when both present",
			raw:  `{"token":"abc123","expires":1,"expires_at":1755486723000,"status":"EXPIRED"}`,
			want: Token{Value: "abc123", ExpiresAt: expires, Status: StatusExpired},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var got Token
			if err := json.Unmarshal([]byte(tc.raw), &got); err != nil {
				t.Fatalf("Unmarshal(%s) error = %v", tc.raw, err)
			}
			if got != tc.want {
				t.Fatalf("Unmarshal(%s) = %+v, want %+v", tc.raw, got, tc.want)
			}
		})
	}
}

func TestTokenIsValid(t *testing.T) {
	t.Parallel()

	now := time.UnixMilli(1755486723000)
	future := now.Add(time.Hour).UnixMilli()
	past := now.Add(-time.Hour).UnixMilli()

	cases := []struct {
		name  string
		token Token
		now   time.Time
		want  bool
	}{
		{"normal future", Token{Value: "t", ExpiresAt: future, Status: StatusNormal}, now, true},
		{"normal past", Token{Value: "t", ExpiresAt: past, Status: StatusNormal}, now, false},
		{"normal zero expiry", Token{Value: "t", Status: StatusNormal}, now, true},
		{"pending future", Token{Value: "t", ExpiresAt: future, Status: StatusPending}, now, false},
		{"expired future", Token{Value: "t", ExpiresAt: future, Status: StatusExpired}, now, false},
		{"normal empty value", Token{ExpiresAt: future, Status: StatusNormal}, now, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.token.IsValid(tc.now); got != tc.want {
				t.Fatalf("IsValid() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestTokenExpiry(t *testing.T) {
	t.Parallel()

	if got := (Token{}).Expiry(); !got.IsZero() {
		t.Fatalf("zero ExpiresAt yields %v, want zero time", got)
	}
	ms := int64(1755486723000)
	if got := (Token{ExpiresAt: ms}).Expiry(); got.UnixMilli() != ms {
		t.Fatalf("Expiry() = %v, want unix millis %d", got, ms)
	}
}

func TestCheckTokenRequestMarshal(t *testing.T) {
	t.Parallel()

	data, err := json.Marshal(CheckTokenRequest{Token: "ccb071"})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if got, want := string(data), `{"token":"ccb071"}`; got != want {
		t.Fatalf("Marshal() = %s, want %s", got, want)
	}
}

func TestCreateTokenRequestMarshal(t *testing.T) {
	t.Parallel()

	data, err := json.Marshal(CreateTokenRequest{})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if got, want := string(data), `{}`; got != want {
		t.Fatalf("Marshal() = %s, want %s", got, want)
	}
}
