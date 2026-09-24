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
	if !stderrors.Is(err, ErrConnectionLimitExceeded) {
		t.Fatal("transport error does not match ErrConnectionLimitExceeded")
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
		if err.Message == "" {
			t.Errorf("FromHTTPStatus(%d) has empty message", tc.status)
		}
	}
}

func TestFromHTTPStatusUsesStatusText(t *testing.T) {
	err := FromHTTPStatus(http.StatusTeapot, nil)
	if err.Code != CodeAPI || err.Status != http.StatusTeapot {
		t.Fatalf("FromHTTPStatus() = %+v, want API/teapot", err)
	}
	if err.Message == "" {
		t.Fatal("FromHTTPStatus() has empty fallback message")
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
