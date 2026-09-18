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
	"time"
)

// TokenStatus is the lifecycle state of a Webull access token as reported by
// the create and check token endpoints.
type TokenStatus string

// Token lifecycle states. A token is created PENDING, becomes NORMAL once 2FA
// has been completed in the Webull App, and is otherwise INVALID or EXPIRED.
// Tokens created in the sandbox environment are NORMAL immediately.
const (
	// StatusPending means the token was created and is awaiting verification.
	StatusPending TokenStatus = "PENDING"
	// StatusNormal means the token is valid and usable for API calls.
	StatusNormal TokenStatus = "NORMAL"
	// StatusInvalid means the token is invalid: it was never used for 15
	// consecutive days, or it does not exist.
	StatusInvalid TokenStatus = "INVALID"
	// StatusExpired means 2FA was not completed within the five minute window.
	StatusExpired TokenStatus = "EXPIRED"
)

// String returns the wire representation of the status.
func (s TokenStatus) String() string { return string(s) }

// Valid reports whether s is one of the four statuses defined by Webull's
// TokenRespVo schema.
func (s TokenStatus) Valid() bool {
	switch s {
	case StatusPending, StatusNormal, StatusInvalid, StatusExpired:
		return true
	default:
		return false
	}
}

// Token is the access-token resource returned by both the create and check
// token endpoints. It doubles as the decoded response DTO (Webull's
// "TokenRespVo") and as the domain representation stored and reused by callers.
type Token struct {
	// Value is the 32-character hexadecimal access token string. It is sent
	// on subsequent requests as the x-access-token header.
	Value string `json:"token"`
	// ExpiresAt is the token expiry as a Unix timestamp in milliseconds.
	ExpiresAt int64 `json:"expires_at"`
	// Status is the current lifecycle state.
	Status TokenStatus `json:"status"`
}

// TokenResponse is the decoded body of the create and check token endpoints.
// Both endpoints share the same TokenRespVo shape, so it is an alias of
// [Token].
type TokenResponse = Token

// Expiry returns the token expiry as a [time.Time]. A zero ExpiresAt is an
// unspecified expiry and yields the zero time.
func (t Token) Expiry() time.Time {
	if t.ExpiresAt == 0 {
		return time.Time{}
	}
	return time.UnixMilli(t.ExpiresAt)
}

// IsValid reports whether the token is usable at now: its value is non-empty,
// its status is [StatusNormal], and it has not passed its expiry. A token with
// an unspecified expiry (ExpiresAt == 0) is treated as not yet expired.
func (t Token) IsValid(now time.Time) bool {
	if t.Value == "" || t.Status != StatusNormal {
		return false
	}
	if t.ExpiresAt == 0 {
		return true
	}
	return now.Before(t.Expiry())
}

// UnmarshalJSON decodes a token response. The documented field is "expires_at";
// some Webull deployments also emit the equivalent "expires" field, which is
// accepted as a fallback so both wire shapes decode correctly.
func (t *Token) UnmarshalJSON(data []byte) error {
	var wire struct {
		Token     string      `json:"token"`
		ExpiresAt int64       `json:"expires_at"`
		Expires   int64       `json:"expires"`
		Status    TokenStatus `json:"status"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	t.Value = wire.Token
	t.Status = wire.Status
	switch {
	case wire.ExpiresAt != 0:
		t.ExpiresAt = wire.ExpiresAt
	default:
		t.ExpiresAt = wire.Expires
	}
	return nil
}

// CreateTokenRequest is the request body of POST /openapi/auth/token/create.
// The endpoint takes no parameters today; the type is defined so the SDK has an
// explicit DTO and can grow fields without changing call sites.
type CreateTokenRequest struct{}

// CheckTokenRequest is the request body of POST /openapi/auth/token/check.
type CheckTokenRequest struct {
	// Token is the access token whose status is being checked.
	Token string `json:"token"`
}
