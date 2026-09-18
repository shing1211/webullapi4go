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
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// Signature algorithm and version identifiers published by Webull.
const (
	// SignatureAlgorithm is the value of the x-signature-algorithm header.
	SignatureAlgorithm = "HMAC-SHA1"
	// SignatureVersion is the value of the x-signature-version header.
	SignatureVersion = "1.0"
	// TimestampFormat is the ISO-8601 UTC layout of the x-timestamp header.
	TimestampFormat = "2006-01-02T15:04:05Z"
)

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
	// body. It must be nil or empty for a bodyless request.
	Body []byte
	// AppSecret is the Webull app secret used as the HMAC key.
	AppSecret string
}

// Sign computes the base64-encoded HMAC-SHA1 signature for p.
//
// The canonical string is built from the signing headers, the query
// parameters and an MD5 digest of the body; see [StringToSign]. The HMAC key
// is the app secret followed by "&".
func Sign(p SignParams) (string, error) {
	encoded, err := StringToSign(p)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha1.New, []byte(p.AppSecret+"&"))
	mac.Write([]byte(encoded))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

// StringToSign returns the percent-encoded canonical string that is fed to
// HMAC-SHA1 by [Sign]. It is exposed for debugging and tests.
//
// It implements Webull's algorithm: the participating query parameters and
// signing headers are merged, sorted by name and joined with "&" to form the
// header string; a non-empty body contributes the uppercase hex MD5 of the
// body. The two are joined with the path and then every byte except the
// RFC 3986 unreserved set (A-Z, a-z, 0-9, "-", "_", ".", "~") is
// percent-encoded with uppercase hex digits, matching Python's
// urllib.parse.quote(s, safe="").
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
			str3 += "&" + bodyDigest(p.Body)
		}
	} else if len(p.Body) > 0 {
		str3 = p.Path + "&" + strings.Join(entries, "&") + "&" + bodyDigest(p.Body)
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

// bodyDigest returns the uppercase hexadecimal MD5 digest of body.
func bodyDigest(body []byte) string {
	sum := md5.Sum(body)
	return strings.ToUpper(hex.EncodeToString(sum[:]))
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

// NewSigningHeaders builds the standard header set that participates in a
// signature for appKey and host, using now (converted to UTC) for the
// timestamp and a fresh random nonce. Pass an empty host to omit the host
// header.
func NewSigningHeaders(appKey, host string, now time.Time) (http.Header, error) {
	nonce, err := NewNonce()
	if err != nil {
		return nil, err
	}
	h := make(http.Header, len(signingHeaderNames))
	h.Set(HeaderAppKey, appKey)
	h.Set(HeaderSignatureAlgorithm, SignatureAlgorithm)
	h.Set(HeaderSignatureVersion, SignatureVersion)
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
