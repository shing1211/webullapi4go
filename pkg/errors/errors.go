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

// Package errs defines the typed error model shared by every webullapi4go
// package.
//
// An [Error] carries a stable [Code] for programmatic classification, a
// human-readable message, and an optional wrapped cause. It participates in
// [errors.Is] and [errors.As] through [Error.Unwrap] and [Error.Is], so callers
// can branch on a category without matching strings.
package errs

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Code is a stable, machine-readable classification of an error.
type Code string

const (
	// CodeInvalidConfig indicates that the client configuration is invalid.
	CodeInvalidConfig Code = "invalid_config"
	// CodeUnsupported indicates an operation the SDK does not support.
	CodeUnsupported Code = "unsupported"
	// CodeAuth indicates an authentication, signing, or token failure.
	CodeAuth Code = "auth"
	// CodeTransport indicates a network or HTTP transport failure.
	CodeTransport Code = "transport"
	// CodeAPI indicates a non-success response returned by the Webull API.
	CodeAPI Code = "api"
	// CodeUnauthorized indicates an HTTP 401 response: the request was not
	// authenticated.
	CodeUnauthorized Code = "UNAUTHORIZED"
	// CodeForbidden indicates an HTTP 403 response: the caller lacks the
	// permission or data subscription required by the endpoint.
	CodeForbidden Code = "FORBIDDEN"
	// CodeInvalidToken indicates an HTTP 417 response: the access token is
	// missing, expired, or otherwise invalid.
	CodeInvalidToken Code = "INVALID_TOKEN"
	// CodeRateLimited indicates an HTTP 429 response: the request exceeded the
	// endpoint's rate limit.
	CodeRateLimited Code = "RATE_LIMITED"
	// CodeServer indicates an HTTP 5xx response: the Webull service failed.
	CodeServer Code = "SERVER_ERROR"
	// CodeInvalidTransition indicates a request that is invalid given the
	// current state of the target resource (for example cancelling an already-filled
	// order).
	CodeInvalidTransition Code = "invalid_transition"
	// CodeNotInitialized indicates a client or resource that is required to be
	// initialized has not been (for example, a stream or events client whose
	// Run method has not been called).
	CodeNotInitialized Code = "not_initialized"
	// CodeValidation indicates an input value failed validation (distinct from
	// a client configuration error).
	CodeValidation Code = "validation"
	// CodeOrderGuardrail indicates an order was rejected because it exceeded a
	// risk guardrail (max quantity, max notional, or similar).
	CodeOrderGuardrail Code = "order_guardrail"
)

// Error is the SDK's typed error. It is safe to return by value as *Error and
// to compare by [Code] via [errors.Is].
type Error struct {
	// Code is the stable error category.
	Code Code
	// Message is a human-readable description of what went wrong.
	Message string
	// Status is the HTTP status code of the response that caused the error, or
	// zero when the error did not originate from an HTTP response.
	Status int
	// Err is the optional underlying cause.
	Err error
}

type sentinelIdentity struct {
	name string
}

func (s *sentinelIdentity) Error() string { return s.name }

// Error implements the error interface.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	cause := e.Err
	if _, marker := cause.(*sentinelIdentity); marker {
		cause = nil
	}
	switch {
	case cause != nil && e.Message != "":
		return fmt.Sprintf("webull: %s: %s: %v", e.Code, e.Message, cause)
	case e.Message != "":
		return fmt.Sprintf("webull: %s: %s", e.Code, e.Message)
	case cause != nil:
		return fmt.Sprintf("webull: %s: %v", e.Code, cause)
	default:
		return fmt.Sprintf("webull: %s", e.Code)
	}
}

// Unwrap returns the wrapped cause, if any, so the error can be inspected with
// [errors.Is] and [errors.As].
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	if _, marker := e.Err.(*sentinelIdentity); marker {
		return nil
	}
	return e.Err
}

func sentinelID(e *Error) *sentinelIdentity {
	if e == nil {
		return nil
	}
	id, _ := e.Err.(*sentinelIdentity)
	return id
}

func semanticID(err error) *sentinelIdentity {
	e, ok := err.(*Error)
	if !ok || e == nil {
		return nil
	}
	if id := sentinelID(e); id != nil {
		return id
	}
	next, ok := e.Err.(*Error)
	if !ok || next == e {
		return nil
	}
	return semanticID(next)
}

// NewSentinel returns a semantic sentinel with the given category and message.
// Unlike [New], it does not match an unrelated error carrying the same category.
func NewSentinel(code Code, message string) *Error {
	return &Error{Code: code, Message: message, Err: &sentinelIdentity{name: "sentinel"}}
}

// Is reports whether e matches target. Ordinary errors match another *Error
// with the same [Code]. A semantic sentinel matches only itself, a copy of
// itself, or an error that wraps that sentinel; matching otherwise delegates
// to the wrapped cause.
func (e *Error) Is(target error) bool {
	if e == nil {
		return false
	}
	if e == target {
		return true
	}
	if id := semanticID(target); id != nil {
		if semanticID(e) == id {
			return true
		}
		return errors.Is(e.Err, target)
	}
	var t *Error
	if errors.As(target, &t) && t != nil && e.Code == t.Code {
		return true
	}
	return errors.Is(e.Err, target)
}

// New returns an *Error with the given code and message.
func New(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

// Wrap returns an *Error with the given code and message that wraps cause.
func Wrap(code Code, message string, cause error) *Error {
	return &Error{Code: code, Message: message, Err: cause}
}

// Is reports whether err is an *Error carrying the given code.
func Is(err error, code Code) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == code
}

// ErrUnsupported is the sentinel for operations the SDK does not support. It
// also satisfies errors.Is(err, errors.ErrUnsupported).
var ErrUnsupported = Wrap(CodeUnsupported, "operation not supported", errors.ErrUnsupported)

// Sentinel errors for the HTTP status classifications returned by
// [FromHTTPStatus]. [errors.Is] matches them by [Code], so callers can write
// errors.Is(err, errs.ErrUnauthorized) without depending on message text.
var (
	// ErrUnauthorized matches [CodeUnauthorized].
	ErrUnauthorized = New(CodeUnauthorized, "unauthorized")
	// ErrForbidden matches [CodeForbidden].
	ErrForbidden = New(CodeForbidden, "forbidden")
	// ErrInvalidToken matches [CodeInvalidToken].
	ErrInvalidToken = New(CodeInvalidToken, "invalid token")
	// ErrRateLimited matches [CodeRateLimited].
	ErrRateLimited = New(CodeRateLimited, "rate limited")
	// ErrServer matches [CodeServer].
	ErrServer = New(CodeServer, "server error")
	// ErrInvalidTransition matches [CodeInvalidTransition].
	ErrInvalidTransition = New(CodeInvalidTransition, "invalid state transition")
	// ErrNotInitialized matches [CodeNotInitialized].
	ErrNotInitialized = New(CodeNotInitialized, "client not initialized")
	// ErrValidation matches [CodeValidation].
	ErrValidation = New(CodeValidation, "validation failed")
	// ErrOrderGuardrail matches [CodeOrderGuardrail].
	ErrOrderGuardrail = New(CodeOrderGuardrail, "order guardrail exceeded")
	// ErrSubscriptionExpired matches [CodeAuth] for subscription/tocket expiry in
	// streaming event clients.
	ErrSubscriptionExpired = NewSentinel(CodeAuth, "subscription expired")
	// ErrConnectionLimitExceeded matches [CodeTransport] for gRPC-stream connection
	// limit rejection (at most 5 concurrent connections per App Key).
	ErrConnectionLimitExceeded = NewSentinel(CodeTransport, "connection limit exceeded")
)

// FromHTTPStatus maps a non-2xx HTTP response status to a typed [Error]. When
// the response body carries an API error message (a JSON object with a
// "message", "msg", "error_msg", "errorMessage", or "error_description"
// field), that message is included; otherwise the raw body is included when it
// is short enough to be readable, and the standard status text is used as a
// last resort. The returned error's [Error.Status] is always status.
func FromHTTPStatus(status int, body []byte) *Error {
	e := &Error{Code: httpStatusCode(status), Status: status}
	msg := apiMessage(body)
	if msg == "" {
		msg = http.StatusText(status)
	}
	if msg == "" {
		msg = "unexpected HTTP status"
	}
	e.Message = fmt.Sprintf("http %d: %s", status, msg)
	return e
}

// httpStatusCode classifies status as one of the SDK's HTTP error codes.
func httpStatusCode(status int) Code {
	switch {
	case status == http.StatusUnauthorized:
		return CodeUnauthorized
	case status == http.StatusForbidden:
		return CodeForbidden
	case status == http.StatusExpectationFailed:
		return CodeInvalidToken
	case status == http.StatusTooManyRequests:
		return CodeRateLimited
	case status >= 500:
		return CodeServer
	default:
		return CodeAPI
	}
}

// apiMessage extracts a human-readable message from a Webull error body. It
// recognizes the JSON field names used across Webull services and falls back to
// the raw, trimmed body when no recognized field is present.
func apiMessage(body []byte) string {
	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" {
		return ""
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(body, &fields) == nil {
		for _, name := range []string{"message", "msg", "error_msg", "errorMessage", "error_description"} {
			raw, ok := fields[name]
			if !ok {
				continue
			}
			var s string
			if json.Unmarshal(raw, &s) == nil && s != "" {
				return s
			}
		}
	}
	const maxRaw = 512
	if len(trimmed) > maxRaw {
		return trimmed[:maxRaw] + "..."
	}
	return trimmed
}
