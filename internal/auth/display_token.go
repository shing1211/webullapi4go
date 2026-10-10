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

import "time"

// DisplayCreateTokenRequest is the request body for POST /auth/client-tokens/create.
type DisplayCreateTokenRequest struct {
	ClientUserID string `json:"client_user_id"`
}

// DisplayToken is the access token returned by the Display Solution token endpoint.
type DisplayToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresAt        int64  `json:"expires_at"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresAt int64  `json:"refresh_expires_at"`
}

// AccessToken returns the bearer token value.
func (t DisplayToken) AccessTokenValue() string {
	return t.AccessToken
}

// ExpiresAtTime returns the expiry as time.Time.
func (t DisplayToken) ExpiresAtTime() time.Time {
	if t.ExpiresAt == 0 {
		return time.Time{}
	}
	return time.UnixMilli(t.ExpiresAt)
}

// IsValid reports whether the token is usable at now.
func (t DisplayToken) IsValid(now time.Time) bool {
	if t.AccessToken == "" {
		return false
	}
	if t.ExpiresAt == 0 {
		return true
	}
	return now.Before(t.ExpiresAtTime())
}
