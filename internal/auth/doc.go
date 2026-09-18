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

// Package auth implements Webull request signing and token lifecycle
// management.
//
// [Sign] and [StringToSign] implement Webull's HMAC-SHA1 request signature:
// the participating query parameters and signing headers form a sorted
// canonical string, a non-empty body contributes its MD5 digest, and the whole
// string is percent-encoded before being HMAC'd with the app secret.
// [NewSigningHeaders] builds the required header set. The token
// create/check/poll/store flow is implemented in token.go.
package auth

// Credentials holds a Webull OpenAPI key pair.
type Credentials struct {
	// AppKey is the Webull OpenAPI app key.
	AppKey string
	// AppSecret is the Webull OpenAPI app secret.
	AppSecret string
}
