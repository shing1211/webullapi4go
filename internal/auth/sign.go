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
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"hash"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// Signature algorithm and version identifiers published by Webull.
const (
	// SignatureAlgorithm is the value of the x-signature-algorithm header for
	// the default [HMAC_SHA1] algorithm.
	SignatureAlgorithm = "HMAC-SHA1"
	// SignatureAlgorithmSHA256 is the value of the x-signature-algorithm header
	// for the [HMAC_SHA256] algorithm used by the gRPC events API.
	SignatureAlgorithmSHA256 = "HMAC-SHA256"
	// SignatureVersion is the value of the x-signature-version header. It is
	// shared by every supported algorithm.
	SignatureVersion = "1.0"
	// TimestampFormat is the ISO-8601 UTC layout of the x-timestamp header.
	TimestampFormat = "2006-01-02T15:04:05Z"
)

// Algorithm selects a Webull request-signing algorithm. The zero value is
// [HMAC_SHA1], so a [SignParams] with no Algorithm set preserves the original
// HTTP signing behavior.
type Algorithm int

const (
	// HMAC_SHA1 is Webull's original request-signing algorithm: the body is
	// hashed with MD5 and the canonical string is signed with HMAC-SHA1.
	HMAC_SHA1 Algorithm = iota
	// HMAC_SHA256 hashes the body with SHA-256 and signs the canonical string
	// with HMAC-SHA256. It is required by the gRPC events API.
	HMAC_SHA256
)

// String returns the x-signature-algorithm header value for a. Any value other
// than [HMAC_SHA256] reports the default [HMAC_SHA1], so an out-of-range
// Algorithm never produces an invalid header.
func (a Algorithm) String() string {
	if a == HMAC_SHA256 {
		return SignatureAlgorithmSHA256
	}
	return SignatureAlgorithm
}

// Version returns the x-signature-version header value for a. Every supported
// algorithm currently publishes version "1.0".
func (a Algorithm) Version() string {
	return SignatureVersion
}

// hashFunc returns the constructor of the hash used by the HMAC, defaulting to
// SHA-1 for an out-of-range Algorithm.
func (a Algorithm) hashFunc() func() hash.Hash {
	if a == HMAC_SHA256 {
		return sha256.New
	}
	return sha1.New
}

// DigestCase selects the hex case for the body digest. The zero value is
// [DigestUpper], which preserves the original behavior (uppercase hex for all
// algorithms). The events API requires lowercase hex.
type DigestCase int

const (
	DigestUpper DigestCase = iota
	DigestLower
)

// bodyDigest returns the hexadecimal digest of body for a. HMAC-SHA1 uses MD5
// and HMAC-SHA256 uses SHA-256, matching Webull's composers. The case of the
// hex digits follows c (default [DigestUpper]).
func (a Algorithm) bodyDigest(body []byte, c DigestCase) string {
	var sum []byte
	if a == HMAC_SHA256 {
		h := sha256.Sum256(body)
		sum = h[:]
	} else {
		h := md5.Sum(body)
		sum = h[:]
	}
	hexStr := hex.EncodeToString(sum)
	if c == DigestLower {
		return strings.ToLower(hexStr)
	}
	return strings.ToUpper(hexStr)
}

// Names of the headers that participate in the signature. Any other header,
// including x-signature and x-version, is ignored.
const (
	HeaderAppKey             = "x-app-key"
	HeaderSignatureAlgorithm = "x-signature-algorithm"
	HeaderSignatureVersion   = "x-signature-version"
	HeaderSignatureNonce     = "x-signature-nonce"
	HeaderTimestamp          = "x-timestamp"
	HeaderHost               = "host"
)

// signingHeaderNames lists, in canonical signing order, the only headers that
// contribute to the canonical string.
var signingHeaderNames = []string{
	HeaderAppKey,
	HeaderSignatureAlgorithm,
	HeaderSignatureVersion,
	HeaderSignatureNonce,
	HeaderTimestamp,
	HeaderHost,
}

// ErrAppSecretRequired is returned by [Sign] and [StringToSign] when no app
// secret is supplied.
var ErrAppSecretRequired = errors.New("auth: app secret is required")

// SignParams carries the inputs to Webull's request-signing algorithm.
type SignParams struct {
	// Method is the HTTP method (e.g. "GET" or "POST"). It does not appear in
	// the canonical string, but is accepted so a caller can describe a full
	// request.
	Method string
	// Path is the request path, for example "/trade/place_order". Leave empty
	// for requests without a path (the gRPC/no-path case).
	Path string
	// Query holds the URL query parameters. It may be nil. A name may repeat;
	// the values of a repeated name are sorted ascending and joined with "&".
	Query url.Values
	// Headers holds the request headers. Only the headers named in
	// [signingHeaderNames] participate; all others are ignored.
	Headers http.Header
	// Body is the exact byte sequence that will be transmitted as the request
	// body. It must be nil or empty for a bodyless request. For the gRPC case
	// it is the serialized protobuf message.
	Body []byte
	// AppSecret is the Webull app secret used as the HMAC key.
	AppSecret string
	// Algorithm selects the signing algorithm. The zero value is [HMAC_SHA1],
	// which preserves the original HTTP behavior; gRPC events use
	// [HMAC_SHA256].
	Algorithm Algorithm
	// DigestCase selects the hex case of the body digest. The zero value is
	// [DigestUpper], which produces uppercase hex to match the published REST
	// API. The gRPC events API requires lowercase hex; pass [DigestLower].
	DigestCase DigestCase
}

// Sign computes the base64-encoded signature for p using the algorithm in
// p.Algorithm (default [HMAC_SHA1]).
//
// The canonical string is built from the signing headers, the query
// parameters and a digest of the body; see [StringToSign]. The HMAC key is the
// app secret followed by "&".
func Sign(p SignParams) (string, error) {
	encoded, err := StringToSign(p)
	if err != nil {
		return "", err
	}
	mac := hmac.New(p.Algorithm.hashFunc(), []byte(p.AppSecret+"&"))
	mac.Write([]byte(encoded))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

// StringToSign returns the percent-encoded canonical string that is fed to the
// HMAC by [Sign]. It is exposed for debugging and tests.
//
// It implements Webull's algorithm: the participating query parameters and
// signing headers are merged, sorted by name and joined with "&" to form the
// header string; a non-empty body contributes the uppercase hex digest of the
// body (MD5 for [HMAC_SHA1], SHA-256 for [HMAC_SHA256]). The two are joined
// with the path and then every byte except the RFC 3986 unreserved set (A-Z,
// a-z, 0-9, "-", "_", ".", "~") is percent-encoded with uppercase hex digits,
// matching Python's urllib.parse.quote(s, safe=""). A pathless request (the
// gRPC form) omits the path and joins the entries with "=" instead.
func StringToSign(p SignParams) (string, error) {
	if p.AppSecret == "" {
		return "", ErrAppSecretRequired
	}
	entries := canonicalEntries(p)
	var str3 string
	if p.Path == "" {
		// No-path (gRPC) form: pairs are joined with "=" and the body hash,
		// when present, is appended with "&".
		str3 = strings.Join(entries, "=")
		if len(p.Body) > 0 {
			str3 += "&" + p.Algorithm.bodyDigest(p.Body, p.DigestCase)
		}
	} else if len(p.Body) > 0 {
		str3 = p.Path + "&" + strings.Join(entries, "&") + "&" + p.Algorithm.bodyDigest(p.Body, p.DigestCase)
	} else {
		str3 = p.Path + "&" + strings.Join(entries, "&")
	}
	return percentEncode(str3), nil
}

// canonicalEntries merges the query parameters and participating headers into
// "name=value" entries sorted by name. A repeated name has its values sorted
// ascending and joined with "&" inside a single entry.
func canonicalEntries(p SignParams) []string {
	values := make(map[string][]string, len(p.Query)+len(signingHeaderNames))
	for name, vs := range p.Query {
		values[name] = append(values[name], vs...)
	}
	for _, name := range signingHeaderNames {
		values[name] = append(values[name], p.Headers.Values(name)...)
	}
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	sort.Strings(names)
	entries := make([]string, 0, len(names))
	for _, name := range names {
		vs := values[name]
		sort.Strings(vs)
		entries = append(entries, name+"="+strings.Join(vs, "&"))
	}
	return entries
}

// percentEncode percent-encodes s using the RFC 3986 unreserved set, leaving
// those bytes untouched and encoding every other byte (including non-ASCII
// UTF-8 bytes) as "%XX" with uppercase hex digits. This reproduces Python's
// urllib.parse.quote(s, safe=""). In particular it encodes "&", "=", "/", ":",
// "?" and " ", so it differs from net/url.QueryEscape (which emits "+" for
// space) and net/url.PathEscape (which leaves "/" and ":" untouched).
func percentEncode(s string) string {
	const upperhex = "0123456789ABCDEF"
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isUnreserved(c) {
			b.WriteByte(c)
			continue
		}
		b.WriteByte('%')
		b.WriteByte(upperhex[c>>4])
		b.WriteByte(upperhex[c&0x0f])
	}
	return b.String()
}

// isUnreserved reports whether c is in the RFC 3986 unreserved set.
func isUnreserved(c byte) bool {
	switch {
	case 'a' <= c && c <= 'z', 'A' <= c && c <= 'Z', '0' <= c && c <= '9':
		return true
	}
	switch c {
	case '-', '_', '.', '~':
		return true
	}
	return false
}

// NewSigningHeaders builds the standard [HMAC_SHA1] header set that
// participates in a signature for appKey and host, using now (converted to
// UTC) for the timestamp and a fresh random nonce. Pass an empty host to omit
// the host header.
//
// It is equivalent to [NewSigningHeadersWithAlgorithm] with [HMAC_SHA1].
func NewSigningHeaders(appKey, host string, now time.Time) (http.Header, error) {
	return NewSigningHeadersWithAlgorithm(appKey, host, now, HMAC_SHA1)
}

// NewSigningHeadersWithAlgorithm builds the standard header set that
// participates in a signature for appKey and host, using now (converted to
// UTC) for the timestamp and a fresh random nonce, and selecting alg for the
// x-signature-algorithm header. Pass an empty host to omit the host header.
// The returned header set must be the same one passed to [Sign] so that the
// canonical string and the transmitted headers agree.
func NewSigningHeadersWithAlgorithm(appKey, host string, now time.Time, alg Algorithm) (http.Header, error) {
	nonce, err := NewNonce()
	if err != nil {
		return nil, err
	}
	h := make(http.Header, len(signingHeaderNames))
	h.Set(HeaderAppKey, appKey)
	h.Set(HeaderSignatureAlgorithm, alg.String())
	h.Set(HeaderSignatureVersion, alg.Version())
	h.Set(HeaderSignatureNonce, nonce)
	h.Set(HeaderTimestamp, now.UTC().Format(TimestampFormat))
	if host != "" {
		h.Set(HeaderHost, host)
	}
	return h, nil
}

// NewNonce returns a fresh 128-bit nonce as 32 lowercase hexadecimal
// characters, matching the format of Webull's x-signature-nonce header.
func NewNonce() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

// MarshalBody encodes v as compact JSON without HTML escaping and without the
// trailing newline that [json.Encoder] adds. Webull signs the exact bytes that
// are transmitted, so callers must sign the value returned here rather than a
// plain json.Marshal result: json.Marshal escapes "<", ">" and "&" as
// "\u003c", "\u003e" and "\u0026", which would change the MD5 body digest.
func MarshalBody(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}
