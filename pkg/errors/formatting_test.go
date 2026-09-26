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

package errs

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestSentinelIdentityError(t *testing.T) {
	t.Parallel()

	id := &sentinelIdentity{name: "sentinel"}
	if got := id.Error(); got != "sentinel" {
		t.Fatalf("sentinelIdentity.Error() = %q, want %q", got, "sentinel")
	}
}

func TestErrorFormatting(t *testing.T) {
	t.Parallel()

	cause := errors.New("dial tcp: timeout")

	cases := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "message and cause",
			err:  &Error{Code: CodeTransport, Message: "request failed", Err: cause},
			want: "webull: transport: request failed: dial tcp: timeout",
		},
		{
			name: "message only",
			err:  &Error{Code: CodeAuth, Message: "token rejected"},
			want: "webull: auth: token rejected",
		},
		{
			name: "cause only",
			err:  &Error{Code: CodeServer, Err: cause},
			want: "webull: SERVER_ERROR: dial tcp: timeout",
		},
		{
			name: "neither message nor cause",
			err:  &Error{Code: CodeInvalidToken},
			want: "webull: INVALID_TOKEN",
		},
		{
			name: "identity marker is not rendered as a cause",
			err:  &Error{Code: CodeAuth, Message: "no access", Err: &sentinelIdentity{name: "sentinel"}},
			want: "webull: auth: no access",
		},
		{
			name: "identity marker alone",
			err:  &Error{Code: CodeAuth, Err: &sentinelIdentity{name: "sentinel"}},
			want: "webull: auth",
		},
		{
			name: "status is not part of the message",
			err:  &Error{Code: CodeRateLimited, Message: "slow down", Status: http.StatusTooManyRequests},
			want: "webull: RATE_LIMITED: slow down",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.err.Error(); got != tc.want {
				t.Fatalf("Error() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestErrorNilReceiver(t *testing.T) {
	t.Parallel()

	var e *Error
	if got := e.Error(); got != "<nil>" {
		t.Fatalf("nil Error.Error() = %q, want %q", got, "<nil>")
	}
	if got := e.Unwrap(); got != nil {
		t.Fatalf("nil Error.Unwrap() = %v, want nil", got)
	}
	if got := sentinelID(e); got != nil {
		t.Fatalf("sentinelID(nil) = %v, want nil", got)
	}
}

func TestErrorUnwrap(t *testing.T) {
	t.Parallel()

	cause := errors.New("real cause")
	if got := (&Error{Code: CodeTransport, Err: cause}).Unwrap(); !errors.Is(got, cause) {
		t.Fatalf("Unwrap() = %v, want the real cause %v", got, cause)
	}

	// The identity marker carries no diagnostic value, so it must not surface
	// as an unwrapped cause.
	if got := (&Error{Code: CodeAuth, Err: &sentinelIdentity{name: "sentinel"}}).Unwrap(); got != nil {
		t.Fatalf("Unwrap() = %v, want nil for an identity marker", got)
	}

	if got := (&Error{Code: CodeAuth}).Unwrap(); got != nil {
		t.Fatalf("Unwrap() with no cause = %v, want nil", got)
	}
}

func TestSentinelIDOnlyForIdentityMarker(t *testing.T) {
	t.Parallel()

	if got := sentinelID(&Error{Code: CodeAuth, Err: errors.New("plain")}); got != nil {
		t.Fatalf("sentinelID(plain cause) = %v, want nil", got)
	}
	if got := sentinelID(&Error{Code: CodeAuth}); got != nil {
		t.Fatalf("sentinelID(no cause) = %v, want nil", got)
	}
	if got := sentinelID(&Error{Code: CodeAuth, Err: &sentinelIdentity{name: "sentinel"}}); got == nil {
		t.Fatal("sentinelID(identity marker) = nil, want the marker")
	}
}

// TestNewSentinelCarriesIdentityMarker pins that a semantic sentinel is an
// *Error whose cause is the identity marker, which is what makes it match only
// itself rather than a category peer.
func TestNewSentinelCarriesIdentityMarker(t *testing.T) {
	t.Parallel()

	s := NewSentinel(CodeAuth, "no access")
	if s.Code != CodeAuth {
		t.Fatalf("Code = %q, want %q", s.Code, CodeAuth)
	}
	if sentinelID(s) == nil {
		t.Fatal("NewSentinel() did not attach an identity marker")
	}
	if got := s.Error(); got != "webull: auth: no access" {
		t.Fatalf("Error() = %q, want %q", got, "webull: auth: no access")
	}
	if got := s.Unwrap(); got != nil {
		t.Fatalf("Unwrap() = %v, want nil for a sentinel", got)
	}
}

// TestIsDoesNotMatchUnrelatedCode pins that category matching is by code and
// never by message text.
func TestIsDoesNotMatchUnrelatedCode(t *testing.T) {
	t.Parallel()

	decoy := &Error{Code: CodeAuth, Message: "rate limit exceeded, RATE_LIMITED"}
	if Is(decoy, CodeRateLimited) {
		t.Fatal("Is() matched a code that only appears in the message text")
	}
	if !Is(decoy, CodeAuth) {
		t.Fatal("Is() did not match the actual code")
	}

	if Is(nil, CodeAuth) {
		t.Fatal("Is(nil, code) = true, want false")
	}
}

// TestIsUnwrapsThroughLayers asserts a typed error is still found when it is
// wrapped by an unrelated layer.
func TestIsUnwrapsThroughLayers(t *testing.T) {
	t.Parallel()

	inner := New(CodeRateLimited, "slow down")
	wrapped := errors.Join(errors.New("context lost"), inner)

	if !Is(wrapped, CodeRateLimited) {
		t.Fatal("Is() did not find a typed error inside errors.Join")
	}
}

func TestFromHTTPStatusMessageSources(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		status      int
		body        string
		wantMessage string
	}{
		{
			name:        "message field",
			status:      http.StatusBadRequest,
			body:        `{"message":"symbol is invalid"}`,
			wantMessage: "symbol is invalid",
		},
		{
			name:        "msg field",
			status:      http.StatusBadRequest,
			body:        `{"msg":"bad msg"}`,
			wantMessage: "bad msg",
		},
		{
			name:        "error_msg field",
			status:      http.StatusBadRequest,
			body:        `{"error_msg":"legacy field"}`,
			wantMessage: "legacy field",
		},
		{
			name:        "errorMessage field",
			status:      http.StatusBadRequest,
			body:        `{"errorMessage":"camel field"}`,
			wantMessage: "camel field",
		},
		{
			name:        "oauth error_description field",
			status:      http.StatusUnauthorized,
			body:        `{"error":"invalid_grant","error_description":"expired code"}`,
			wantMessage: "expired code",
		},
		{
			name:        "empty body falls back to status text",
			status:      http.StatusNotFound,
			body:        "",
			wantMessage: http.StatusText(http.StatusNotFound),
		},
		{
			name:        "whitespace body falls back to status text",
			status:      http.StatusNotFound,
			body:        "   \n ",
			wantMessage: http.StatusText(http.StatusNotFound),
		},
		{
			name:        "unrecognized body is kept verbatim",
			status:      http.StatusBadRequest,
			body:        `{}`,
			wantMessage: "{}",
		},
		{
			name:        "unknown status with no status text",
			status:      599,
			body:        "",
			wantMessage: "unexpected HTTP status",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := FromHTTPStatus(tc.status, []byte(tc.body))
			if got.Status != tc.status {
				t.Fatalf("Status = %d, want %d", got.Status, tc.status)
			}
			if !strings.Contains(got.Message, tc.wantMessage) {
				t.Fatalf("Message = %q, want it to contain %q", got.Message, tc.wantMessage)
			}
			if !strings.Contains(got.Message, "http") {
				t.Fatalf("Message = %q, want it to record the status", got.Message)
			}
		})
	}
}

func TestAPIMessageTruncatesLongBodies(t *testing.T) {
	t.Parallel()

	long := strings.Repeat("x", 600)
	got := FromHTTPStatus(http.StatusBadGateway, []byte(long))
	if !strings.HasSuffix(got.Message, "...") {
		t.Fatalf("Message = %q, want a truncation marker", got.Message)
	}
	if strings.Contains(got.Message, strings.Repeat("x", 513)) {
		t.Fatal("Message was not truncated to the documented limit")
	}
}

// TestFromHTTPStatus417PreservesAPIReason guards the documented contract that a
// 417 keeps its status and API message rather than being flattened into a token
// failure by the helper.
func TestFromHTTPStatus417PreservesAPIReason(t *testing.T) {
	t.Parallel()

	got := FromHTTPStatus(http.StatusExpectationFailed, []byte(`{"message":"Invalid Symbol"}`))
	if got.Status != http.StatusExpectationFailed {
		t.Fatalf("Status = %d, want %d", got.Status, http.StatusExpectationFailed)
	}
	if got.Code != CodeInvalidToken {
		t.Fatalf("Code = %q, want the historical %q for compatibility", got.Code, CodeInvalidToken)
	}
	if !strings.Contains(got.Message, "Invalid Symbol") {
		t.Fatalf("Message = %q, want the API reason preserved", got.Message)
	}
	if !strings.Contains(got.Error(), "417") {
		t.Fatalf("Error() = %q, want it to mention the status", got.Error())
	}
}
