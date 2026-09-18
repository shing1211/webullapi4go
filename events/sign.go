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

package events

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"sort"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/shing1211/webullapi4go/internal/auth"
)

// subscribeMetadata builds the signed gRPC metadata for body.
//
// The event service signs a different canonical string from the REST API. Only
// the five x-signature parameters participate (there is no host, path, or
// query), they are joined with "=" in sorted name order, and the body is
// appended as "&" followed by the lowercase hexadecimal SHA-256 digest of the
// serialized request. The canonical string is then percent-encoded with the
// RFC 3986 unreserved set and signed with HMAC-SHA256 using "app secret &" as
// the key. The computed signature is carried under x-signature.
//
// The events canonical string hashes the body to lowercase hex, which is why
// this package does not call [auth.Sign]: that signer uppercases the SHA-256
// digest for the REST API and would be rejected by the event service.
func subscribeMetadata(appKey, appSecret string, now time.Time, body []byte) (metadata.MD, error) {
	nonce, err := auth.NewNonce()
	if err != nil {
		return nil, err
	}
	params := map[string]string{
		auth.HeaderAppKey:             appKey,
		auth.HeaderSignatureAlgorithm: auth.SignatureAlgorithmSHA256,
		auth.HeaderSignatureVersion:   auth.SignatureVersion,
		auth.HeaderSignatureNonce:     nonce,
		auth.HeaderTimestamp:          now.UTC().Format(auth.TimestampFormat),
	}
	signature := eventSignature(params, appSecret, body)
	return metadata.Pairs(
		auth.HeaderAppKey, params[auth.HeaderAppKey],
		auth.HeaderSignatureAlgorithm, params[auth.HeaderSignatureAlgorithm],
		auth.HeaderSignatureVersion, params[auth.HeaderSignatureVersion],
		auth.HeaderSignatureNonce, params[auth.HeaderSignatureNonce],
		auth.HeaderTimestamp, params[auth.HeaderTimestamp],
		metadataSignature, signature,
	), nil
}

// eventSignature returns the base64 HMAC-SHA256 signature of the event
// canonical string derived from params and body.
func eventSignature(params map[string]string, appSecret string, body []byte) string {
	encoded := percentEncode(eventCanonicalString(params, body))
	mac := hmac.New(sha256.New, []byte(appSecret+"&"))
	mac.Write([]byte(encoded))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// eventCanonicalString builds the event canonical string: the parameters sorted
// by name and joined as "name=value", then "&" and the lowercase hex SHA-256
// digest of body.
func eventCanonicalString(params map[string]string, body []byte) string {
	names := make([]string, 0, len(params))
	for name := range params {
		names = append(names, name)
	}
	sort.Strings(names)
	entries := make([]string, 0, len(names))
	for _, name := range names {
		entries = append(entries, name+"="+params[name])
	}
	sum := sha256.Sum256(body)
	return strings.Join(entries, "=") + "&" + hex.EncodeToString(sum[:])
}

// percentEncode percent-encodes s using the RFC 3986 unreserved set, matching
// the encoding Webull applies to its canonical strings. It mirrors the encoder
// in internal/auth, which is not exported.
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
