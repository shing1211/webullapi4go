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

// Package errs is a deprecated alias for [github.com/shing1211/webullapi4go/pkg/errors].
// New code should import that package directly. This package will be removed in v3.
package errs

import (
	pkgres "github.com/shing1211/webullapi4go/pkg/errors"
)

// Code is a stable, machine-readable classification of an error.
//
// Deprecated: use [pkgres.Code].
type Code = pkgres.Code

const (
	// CodeInvalidConfig is deprecated.
	//
	// Deprecated: use [pkgres.CodeInvalidConfig].
	CodeInvalidConfig = pkgres.CodeInvalidConfig
	// CodeUnsupported is deprecated.
	//
	// Deprecated: use [pkgres.CodeUnsupported].
	CodeUnsupported = pkgres.CodeUnsupported
	// CodeAuth is deprecated.
	//
	// Deprecated: use [pkgres.CodeAuth].
	CodeAuth = pkgres.CodeAuth
	// CodeTransport is deprecated.
	//
	// Deprecated: use [pkgres.CodeTransport].
	CodeTransport = pkgres.CodeTransport
	// CodeAPI is deprecated.
	//
	// Deprecated: use [pkgres.CodeAPI].
	CodeAPI = pkgres.CodeAPI
	// CodeUnauthorized is deprecated.
	//
	// Deprecated: use [pkgres.CodeUnauthorized].
	CodeUnauthorized = pkgres.CodeUnauthorized
	// CodeForbidden is deprecated.
	//
	// Deprecated: use [pkgres.CodeForbidden].
	CodeForbidden = pkgres.CodeForbidden
	// CodeInvalidToken is deprecated.
	//
	// Deprecated: use [pkgres.CodeInvalidToken].
	CodeInvalidToken = pkgres.CodeInvalidToken
	// CodeRateLimited is deprecated.
	//
	// Deprecated: use [pkgres.CodeRateLimited].
	CodeRateLimited = pkgres.CodeRateLimited
	// CodeServer is deprecated.
	//
	// Deprecated: use [pkgres.CodeServer].
	CodeServer = pkgres.CodeServer
)

// Error is the SDK's typed error.
//
// Deprecated: use [pkgres.Error].
type Error = pkgres.Error

// ErrUnsupported is the sentinel for unsupported operations.
//
// Deprecated: use [pkgres.ErrUnsupported].
var ErrUnsupported = pkgres.ErrUnsupported

// ErrUnauthorized matches CodeUnauthorized.
//
// Deprecated: use [pkgres.ErrUnauthorized].
var ErrUnauthorized = pkgres.ErrUnauthorized

// ErrForbidden matches CodeForbidden.
//
// Deprecated: use [pkgres.ErrForbidden].
var ErrForbidden = pkgres.ErrForbidden

// ErrInvalidToken matches CodeInvalidToken.
//
// Deprecated: use [pkgres.ErrInvalidToken].
var ErrInvalidToken = pkgres.ErrInvalidToken

// ErrRateLimited matches CodeRateLimited.
//
// Deprecated: use [pkgres.ErrRateLimited].
var ErrRateLimited = pkgres.ErrRateLimited

// ErrServer matches CodeServer.
//
// Deprecated: use [pkgres.ErrServer].
var ErrServer = pkgres.ErrServer

// New returns an *Error with the given code and message.
//
// Deprecated: use [pkgres.New].
func New(code Code, message string) *Error { return pkgres.New(code, message) }

// Wrap returns an *Error with the given code and message that wraps cause.
//
// Deprecated: use [pkgres.Wrap].
func Wrap(code Code, message string, cause error) *Error { return pkgres.Wrap(code, message, cause) }

// Is reports whether err is an *Error carrying the given code.
//
// Deprecated: use [pkgres.Is].
func Is(err error, code Code) bool { return pkgres.Is(err, code) }

// FromHTTPStatus maps a non-2xx HTTP response status to a typed Error.
//
// Deprecated: use [pkgres.FromHTTPStatus].
func FromHTTPStatus(status int, body []byte) *Error { return pkgres.FromHTTPStatus(status, body) }
