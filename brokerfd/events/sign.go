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
	"sort"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/shing1211/webullapi4go/internal/auth"
)

const metadataSignature = "x-signature"

func subscribeMetadata(appKey, appSecret string, now time.Time, body []byte) (metadata.MD, error) {
	nonce, err := auth.NewNonce()
	if err != nil {
		return nil, err
	}
	ts := now.UTC().Format(auth.TimestampFormat)

	params := map[string]string{
		auth.HeaderAppKey:             appKey,
		auth.HeaderSignatureAlgorithm: auth.SignatureAlgorithmSHA256,
		auth.HeaderSignatureVersion:   auth.SignatureVersion,
		auth.HeaderSignatureNonce:     nonce,
		auth.HeaderTimestamp:          ts,
	}

	sig, err := eventSignature(params, appSecret, body)
	if err != nil {
		return nil, err
	}

	return metadata.Pairs(
		auth.HeaderAppKey, params[auth.HeaderAppKey],
		auth.HeaderSignatureAlgorithm, params[auth.HeaderSignatureAlgorithm],
		auth.HeaderSignatureVersion, params[auth.HeaderSignatureVersion],
		auth.HeaderSignatureNonce, params[auth.HeaderSignatureNonce],
		auth.HeaderTimestamp, params[auth.HeaderTimestamp],
		metadataSignature, sig,
	), nil
}

func eventSignature(params map[string]string, appSecret string, body []byte) (string, error) {
	encoded := percentEncode(eventCanonicalString(params, body))
	mac := hmac.New(sha256.New, []byte(appSecret+"&"))
	mac.Write([]byte(encoded))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

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
	return strings.Join(entries, "=") + "&" + strings.ToLower(hexEncode(sum[:]))
}

func hexEncode(b []byte) string {
	const hexChars = "0123456789abcdef"
	result := make([]byte, len(b)*2)
	for i, c := range b {
		result[i*2] = hexChars[c>>4]
		result[i*2+1] = hexChars[c&0xf]
	}
	return string(result)
}

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
		b.WriteByte(upperhex[c&0xf])
	}
	return b.String()
}

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
