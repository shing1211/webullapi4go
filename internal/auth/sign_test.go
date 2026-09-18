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

package auth_test

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/internal/auth"
)

// testHeaders returns a deterministic signing header set used by the
// non-golden cases.
func testHeaders() http.Header {
	h := make(http.Header)
	h.Set("x-app-key", "test-app-key")
	h.Set(auth.HeaderSignatureAlgorithm, auth.SignatureAlgorithm)
	h.Set(auth.HeaderSignatureVersion, auth.SignatureVersion)
	h.Set(auth.HeaderSignatureNonce, "0123456789abcdef0123456789abcdef")
	h.Set(auth.HeaderTimestamp, "2022-01-04T03:55:31Z")
	h.Set(auth.HeaderHost, "api.webull.com")
	return h
}

func TestSignGoldenVector(t *testing.T) {
	t.Parallel()

	params := auth.SignParams{
		Method: "POST",
		Path:   "/trade/place_order",
		Query: url.Values{
			"a1": {"webull"},
			"a2": {"123"},
			"a3": {"xxx"},
			"q1": {"yyy"},
		},
		Headers: func() http.Header {
			h := make(http.Header)
			h.Set("x-app-key", "776da210ab4a452795d74e726ebd74b6")
			h.Set(auth.HeaderTimestamp, "2022-01-04T03:55:31Z")
			h.Set(auth.HeaderSignatureVersion, "1.0")
			h.Set(auth.HeaderSignatureAlgorithm, "HMAC-SHA1")
			h.Set(auth.HeaderSignatureNonce, "48ef5afed43d4d91ae514aaeafbc29ba")
			h.Set(auth.HeaderHost, "api.webull.com")
			return h
		}(),
		Body:      []byte(`{"k1":123,"k2":"this is the api request body","k3":true,"k4":{"foo":[1,2]}}`),
		AppSecret: "0f50a2e853334a9aae1a783bee120c1f",
	}

	const want = "kvlS6opdZDhEBo5jq40nHYXaLvM="
	got, err := auth.Sign(params)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	if got != want {
		t.Fatalf("Sign() = %q, want %q", got, want)
	}

	// The URL-encoded canonical string underpinning the golden signature.
	const wantStringToSign = "%2Ftrade%2Fplace_order%26a1%3Dwebull%26a2%3D123%26a3%3Dxxx%26host%3Dapi.webull.com%26q1%3Dyyy%26x-app-key%3D776da210ab4a452795d74e726ebd74b6%26x-signature-algorithm%3DHMAC-SHA1%26x-signature-nonce%3D48ef5afed43d4d91ae514aaeafbc29ba%26x-signature-version%3D1.0%26x-timestamp%3D2022-01-04T03%3A55%3A31Z%26E296C96787E1A309691CEF3692F5EEDD"
	if got := mustStringToSign(t, params); got != wantStringToSign {
		t.Fatalf("StringToSign() =\n%q\nwant\n%q", got, wantStringToSign)
	}
}

func TestStringToSignEmptyBody(t *testing.T) {
	t.Parallel()

	params := auth.SignParams{
		Path:      "/market-data/stock/quotes",
		Query:     url.Values{"symbols": {"AAPL", "MSFT"}},
		Headers:   testHeaders(),
		AppSecret: "secret",
	}

	const wantSign = "38xSjgdHuwqRpXxO+zzXjF+rJ5s="
	got, err := auth.Sign(params)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	if got != wantSign {
		t.Fatalf("Sign() = %q, want %q", got, wantSign)
	}

	const wantStringToSign = "%2Fmarket-data%2Fstock%2Fquotes%26host%3Dapi.webull.com%26symbols%3DAAPL%26MSFT%26x-app-key%3Dtest-app-key%26x-signature-algorithm%3DHMAC-SHA1%26x-signature-nonce%3D0123456789abcdef0123456789abcdef%26x-signature-version%3D1.0%26x-timestamp%3D2022-01-04T03%3A55%3A31Z"
	if got := mustStringToSign(t, params); got != wantStringToSign {
		t.Fatalf("StringToSign() =\n%q\nwant\n%q", got, wantStringToSign)
	}
}

func TestStringToSignDuplicateQueryParams(t *testing.T) {
	t.Parallel()

	params := auth.SignParams{
		Path: "/x",
		Query: url.Values{
			"b": {"y"},
			"a": {"2", "1", "0"},
		},
		Headers:   testHeaders(),
		AppSecret: "secret",
	}

	// a's values are sorted ascending and joined with "&" inside one entry.
	const wantStringToSign = "%2Fx%26a%3D0%261%262%26b%3Dy%26host%3Dapi.webull.com%26x-app-key%3Dtest-app-key%26x-signature-algorithm%3DHMAC-SHA1%26x-signature-nonce%3D0123456789abcdef0123456789abcdef%26x-signature-version%3D1.0%26x-timestamp%3D2022-01-04T03%3A55%3A31Z"
	if got := mustStringToSign(t, params); got != wantStringToSign {
		t.Fatalf("StringToSign() =\n%q\nwant\n%q", got, wantStringToSign)
	}

	// Order of the duplicate values must not change the signature.
	shuffled := params
	shuffled.Query = url.Values{
		"b": {"y"},
		"a": {"1", "0", "2"},
	}
	if got, want := mustSign(t, shuffled), mustSign(t, params); got != want {
		t.Fatalf("Sign() with reordered duplicates = %q, want %q", got, want)
	}
}

func TestStringToSignNoPath(t *testing.T) {
	t.Parallel()

	params := auth.SignParams{
		Query:     url.Values{"b": {"2"}, "a": {"1"}},
		Headers:   testHeaders(),
		Body:      []byte(`{"k":"v"}`),
		AppSecret: "secret",
	}

	const wantSign = "SPAQ9WAdGrZtzlNDxit0oCfRH50="
	got, err := auth.Sign(params)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	if got != wantSign {
		t.Fatalf("Sign() = %q, want %q", got, wantSign)
	}

	// The no-path form joins pairs with "=" and appends "&<md5>" for the body.
	const wantStringToSign = "a%3D1%3Db%3D2%3Dhost%3Dapi.webull.com%3Dx-app-key%3Dtest-app-key%3Dx-signature-algorithm%3DHMAC-SHA1%3Dx-signature-nonce%3D0123456789abcdef0123456789abcdef%3Dx-signature-version%3D1.0%3Dx-timestamp%3D2022-01-04T03%3A55%3A31Z%2644244CE1A15EE6D4DC270001564CB759"
	if got := mustStringToSign(t, params); got != wantStringToSign {
		t.Fatalf("StringToSign() =\n%q\nwant\n%q", got, wantStringToSign)
	}
}

func TestSignBodyWithHTMLEscapeCharacters(t *testing.T) {
	t.Parallel()

	body, err := auth.MarshalBody(map[string]string{"x": "<a>&b"})
	if err != nil {
		t.Fatalf("MarshalBody() error = %v", err)
	}
	if string(body) != `{"x":"<a>&b"}` {
		t.Fatalf("MarshalBody() = %q, want %q", body, `{"x":"<a>&b"}`)
	}

	params := auth.SignParams{
		Path:      "/y",
		Query:     url.Values{"q": {"1"}},
		Headers:   testHeaders(),
		Body:      body,
		AppSecret: "secret",
	}

	const wantSign = "rklZ8HzkUl+GMiSY/x5AuyqqzgM="
	if got, err := auth.Sign(params); err != nil || got != wantSign {
		t.Fatalf("Sign() = %q, %v; want %q, nil", got, err, wantSign)
	}

	const wantStringToSign = "%2Fy%26host%3Dapi.webull.com%26q%3D1%26x-app-key%3Dtest-app-key%26x-signature-algorithm%3DHMAC-SHA1%26x-signature-nonce%3D0123456789abcdef0123456789abcdef%26x-signature-version%3D1.0%26x-timestamp%3D2022-01-04T03%3A55%3A31Z%26D9DC7DD92EA2092DC2CC754B64A50985"
	if got := mustStringToSign(t, params); got != wantStringToSign {
		t.Fatalf("StringToSign() =\n%q\nwant\n%q", got, wantStringToSign)
	}

	// json.Marshal escapes <, > and &, which yields a different body digest and
	// therefore a different signature.
	escaped, err := json.Marshal(map[string]string{"x": "<a>&b"})
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	escapedParams := params
	escapedParams.Body = escaped
	if escapedSig, _ := auth.Sign(escapedParams); escapedSig == wantSign {
		t.Fatalf("json.Marshal body unexpectedly produced the same signature as the unescaped body")
	}
}

func TestSigningHeadersDoNotParticipate(t *testing.T) {
	t.Parallel()

	base := auth.SignParams{
		Path:      "/x",
		Headers:   testHeaders(),
		AppSecret: "secret",
	}
	want := mustSign(t, base)

	// x-signature and x-version are explicitly excluded from signing.
	withExtras := base
	withExtras.Headers = testHeaders()
	withExtras.Headers.Set("x-signature", "ignored")
	withExtras.Headers.Set("x-version", "1.0")
	withExtras.Headers.Set("x-request-id", "also-ignored")
	if got := mustSign(t, withExtras); got != want {
		t.Fatalf("Sign() changed when non-participating headers were added: %q != %q", got, want)
	}
}

func TestSignRequiresAppSecret(t *testing.T) {
	t.Parallel()

	if _, err := auth.Sign(auth.SignParams{Path: "/x"}); err == nil {
		t.Fatalf("Sign() with empty secret returned nil error")
	}
	if _, err := auth.StringToSign(auth.SignParams{Path: "/x"}); err == nil {
		t.Fatalf("StringToSign() with empty secret returned nil error")
	}
}

func TestNewSigningHeaders(t *testing.T) {
	t.Parallel()

	now := time.Date(2022, 1, 4, 11, 55, 31, 0, time.FixedZone("CST", 8*60*60))
	h, err := auth.NewSigningHeaders("my-key", "api.webull.com", now)
	if err != nil {
		t.Fatalf("NewSigningHeaders() error = %v", err)
	}

	if got := h.Get(auth.HeaderAppKey); got != "my-key" {
		t.Fatalf("x-app-key = %q, want %q", got, "my-key")
	}
	if got := h.Get(auth.HeaderSignatureAlgorithm); got != auth.SignatureAlgorithm {
		t.Fatalf("x-signature-algorithm = %q, want %q", got, auth.SignatureAlgorithm)
	}
	if got := h.Get(auth.HeaderSignatureVersion); got != auth.SignatureVersion {
		t.Fatalf("x-signature-version = %q, want %q", got, auth.SignatureVersion)
	}
	if got := h.Get(auth.HeaderHost); got != "api.webull.com" {
		t.Fatalf("host = %q, want %q", got, "api.webull.com")
	}
	if got := h.Get(auth.HeaderTimestamp); got != "2022-01-04T03:55:31Z" {
		t.Fatalf("x-timestamp = %q, want %q", got, "2022-01-04T03:55:31Z")
	}
	assertHexNonce(t, h.Get(auth.HeaderSignatureNonce))
}

func TestNewNonceIsRandomHex(t *testing.T) {
	t.Parallel()

	first, err := auth.NewNonce()
	if err != nil {
		t.Fatalf("NewNonce() error = %v", err)
	}
	assertHexNonce(t, first)

	second, err := auth.NewNonce()
	if err != nil {
		t.Fatalf("NewNonce() error = %v", err)
	}
	if first == second {
		t.Fatalf("NewNonce() returned the same value twice: %q", first)
	}
}

func mustSign(t *testing.T, p auth.SignParams) string {
	t.Helper()
	sig, err := auth.Sign(p)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	return sig
}

func mustStringToSign(t *testing.T, p auth.SignParams) string {
	t.Helper()
	s, err := auth.StringToSign(p)
	if err != nil {
		t.Fatalf("StringToSign() error = %v", err)
	}
	return s
}

func assertHexNonce(t *testing.T, nonce string) {
	t.Helper()
	if len(nonce) != 32 {
		t.Fatalf("nonce = %q, want 32 hex characters", nonce)
	}
	if _, err := hex.DecodeString(nonce); err != nil {
		t.Fatalf("nonce = %q is not valid hex: %v", nonce, err)
	}
	if nonce != strings.ToLower(nonce) {
		t.Fatalf("nonce = %q, want lowercase hex", nonce)
	}
}
