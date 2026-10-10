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
	stderrors "errors"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestErrorWrappingAndMatching(t *testing.T) {
	cause := stderrors.New("cause")
	err := Wrap(CodeTransport, "request failed", cause)
	if !stderrors.Is(err, cause) {
		t.Fatal("wrapped error does not match its cause")
	}
	if !Is(err, CodeTransport) {
		t.Fatal("wrapped error does not match its code")
	}
	if stderrors.Is(err, ErrConnectionLimitExceeded) {
		t.Fatal("generic transport error matches ErrConnectionLimitExceeded")
	}
	wrappedSemantic := Wrap(CodeTransport, "connection limit", ErrConnectionLimitExceeded)
	if !stderrors.Is(wrappedSemantic, ErrConnectionLimitExceeded) {
		t.Fatal("wrapped semantic cause does not match its sentinel")
	}
	var typed *Error
	if !stderrors.As(err, &typed) || typed != err {
		t.Fatalf("errors.As() = %#v, want original error", typed)
	}
}

func TestSentinelErrorIdentity(t *testing.T) {
	if !stderrors.Is(ErrValidation, ErrValidation) {
		t.Fatal("ErrValidation does not match itself")
	}
	if !stderrors.Is(Wrap(CodeValidation, "failed", stderrors.New("cause")), ErrValidation) {
		t.Fatal("CodeValidation error does not match ErrValidation")
	}
}

func TestSemanticSentinelSpecificity(t *testing.T) {
	genericTransport := New(CodeTransport, "dial failed")
	genericAuth := New(CodeAuth, "signing failed")

	for _, target := range []error{ErrConnectionLimitExceeded, ErrSubscriptionExpired} {
		if stderrors.Is(genericTransport, target) {
			t.Errorf("generic transport error matches %v", target)
		}
		if stderrors.Is(genericAuth, target) {
			t.Errorf("generic auth error matches %v", target)
		}
	}
	if stderrors.Is(ErrConnectionLimitExceeded, ErrSubscriptionExpired) {
		t.Fatal("connection-limit sentinel matches subscription-expired sentinel")
	}
	if stderrors.Is(ErrSubscriptionExpired, ErrConnectionLimitExceeded) {
		t.Fatal("subscription-expired sentinel matches connection-limit sentinel")
	}
	if !stderrors.Is(ErrConnectionLimitExceeded, ErrConnectionLimitExceeded) {
		t.Fatal("connection-limit sentinel does not match itself")
	}
	if !stderrors.Is(ErrSubscriptionExpired, ErrSubscriptionExpired) {
		t.Fatal("subscription-expired sentinel does not match itself")
	}
	if !stderrors.Is(Wrap(CodeTransport, "outer", ErrConnectionLimitExceeded), ErrConnectionLimitExceeded) {
		t.Fatal("wrapped connection-limit cause does not match its sentinel")
	}
	wrappedTarget := Wrap(CodeTransport, "target wrapper", ErrConnectionLimitExceeded)
	if stderrors.Is(genericTransport, wrappedTarget) {
		t.Fatal("generic transport error matches a wrapper carrying a semantic cause")
	}
	if !stderrors.Is(Wrap(CodeAuth, "outer", ErrSubscriptionExpired), ErrSubscriptionExpired) {
		t.Fatal("wrapped subscription-expired cause does not match its sentinel")
	}
	if !Is(ErrConnectionLimitExceeded, CodeTransport) {
		t.Fatal("connection-limit sentinel lost its transport category")
	}
	if !Is(ErrSubscriptionExpired, CodeAuth) {
		t.Fatal("subscription-expired sentinel lost its auth category")
	}
}

func TestFromHTTPStatus(t *testing.T) {
	cases := []struct {
		status int
		body   string
		code   Code
		msg    string
	}{
		{status: http.StatusUnauthorized, body: `{"message":"expired"}`, code: CodeUnauthorized, msg: "expired"},
		{status: http.StatusForbidden, body: `not-json`, code: CodeForbidden, msg: "not-json"},
		{status: http.StatusExpectationFailed, body: `{"error_description":"token"}`, code: CodeInvalidToken, msg: "token"},
		{status: http.StatusTooManyRequests, body: `{"msg":"slow down"}`, code: CodeRateLimited, msg: "slow down"},
		{status: http.StatusBadGateway, body: `{"error_msg":"upstream"}`, code: CodeServer, msg: "upstream"},
		{status: http.StatusBadRequest, body: `{"errorMessage":"bad"}`, code: CodeAPI, msg: "bad"},
	}
	for _, tc := range cases {
		err := FromHTTPStatus(tc.status, []byte(tc.body))
		if err.Code != tc.code || err.Status != tc.status {
			t.Errorf("FromHTTPStatus(%d) = %+v, want code %s status %d", tc.status, err, tc.code, tc.status)
		}
		wantMessage := "http " + strconv.Itoa(tc.status) + ": " + tc.msg
		if err.Message != wantMessage {
			t.Errorf("FromHTTPStatus(%d).Message = %q, want %q", tc.status, err.Message, wantMessage)
		}
	}
}

func TestFromHTTPStatusUsesStatusText(t *testing.T) {
	err := FromHTTPStatus(http.StatusTeapot, nil)
	if err.Code != CodeAPI || err.Status != http.StatusTeapot {
		t.Fatalf("FromHTTPStatus() = %+v, want API/teapot", err)
	}
	if err.Message != "http 418: I'm a teapot" {
		t.Fatalf("FromHTTPStatus().Message = %q, want status text", err.Message)
	}
}

func TestAPIMessageFieldVariants(t *testing.T) {
	for _, body := range [][]byte{
		[]byte(`{"message":"message"}`),
		[]byte(`{"msg":"msg"}`),
		[]byte(`{"error_msg":"error_msg"}`),
		[]byte(`{"errorMessage":"errorMessage"}`),
		[]byte(`{"error_description":"error_description"}`),
	} {
		if got := apiMessage(body); got == "" {
			t.Errorf("apiMessage(%s) is empty", body)
		}
	}
}

func TestAPIMessagePrecedence(t *testing.T) {
	body := []byte(`{"error_description":"last","message":"first","msg":"second"}`)
	if got := apiMessage(body); got != "first" {
		t.Fatalf("apiMessage() = %q, want first field", got)
	}
	body = []byte(`{"message":42,"msg":"fallback"}`)
	if got := apiMessage(body); got != "fallback" {
		t.Fatalf("apiMessage() = %q, want fallback field", got)
	}
	body = []byte(`{"message":"","msg":"fallback"}`)
	if got := apiMessage(body); got != "fallback" {
		t.Fatalf("apiMessage() = %q, want non-empty fallback field", got)
	}
}

func TestAPIMessageRawFallbackAndTruncation(t *testing.T) {
	if got := apiMessage([]byte(" \n malformed body \t")); got != "malformed body" {
		t.Fatalf("apiMessage() = %q, want trimmed raw body", got)
	}
	if got := apiMessage([]byte(`{"unrecognized":"value"}`)); got != `{"unrecognized":"value"}` {
		t.Fatalf("apiMessage() = %q, want raw JSON body", got)
	}
	raw := strings.Repeat("x", 513)
	want := strings.Repeat("x", 512) + "..."
	if got := apiMessage([]byte(raw)); got != want {
		t.Fatalf("apiMessage() length/content = %d/%q, want %d/%q", len(got), got, len(want), want)
	}
}

func TestFromHTTPStatusEmptyAndMalformedBodies(t *testing.T) {
	empty := FromHTTPStatus(http.StatusTeapot, nil)
	if empty.Message != "http 418: I'm a teapot" {
		t.Fatalf("empty body message = %q, want status text", empty.Message)
	}
	malformed := FromHTTPStatus(http.StatusBadRequest, []byte(`{"message":`))
	if malformed.Message != `http 400: {"message":` {
		t.Fatalf("malformed body message = %q, want raw body", malformed.Message)
	}
}

func TestHTTPSentinelMatching(t *testing.T) {
	cases := []struct {
		status int
		code   Code
		target error
	}{
		{status: http.StatusUnauthorized, code: CodeUnauthorized, target: ErrUnauthorized},
		{status: http.StatusForbidden, code: CodeForbidden, target: ErrForbidden},
		{status: http.StatusExpectationFailed, code: CodeInvalidToken, target: ErrInvalidToken},
		{status: http.StatusTooManyRequests, code: CodeRateLimited, target: ErrRateLimited},
		{status: http.StatusBadGateway, code: CodeServer, target: ErrServer},
	}
	for _, tc := range cases {
		err := FromHTTPStatus(tc.status, nil)
		if !stderrors.Is(err, tc.target) {
			t.Errorf("errors.Is(FromHTTPStatus(%d), %v) = false, want true", tc.status, tc.target)
		}
		if !Is(err, tc.code) {
			t.Errorf("Is(FromHTTPStatus(%d), %s) = false, want true", tc.status, tc.code)
		}
	}
}
