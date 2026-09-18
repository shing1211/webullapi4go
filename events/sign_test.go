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

import "testing"

// testSignParams returns the fixed sign parameters shared by the golden tests.
func testSignParams() map[string]string {
	return map[string]string{
		"x-app-key":             "test-app-key",
		"x-signature-algorithm": "HMAC-SHA256",
		"x-signature-version":   "1.0",
		"x-signature-nonce":     "0123456789abcdef0123456789abcdef",
		"x-timestamp":           "2022-01-04T03:55:31Z",
	}
}

// TestEventCanonicalStringGolden pins the events canonical string, in
// particular that the SHA-256 body digest is lowercase hex (unlike the REST
// signer). The expected values were produced by Webull's Python
// signature_composer using the same parameters and body.
func TestEventCanonicalStringGolden(t *testing.T) {
	const (
		wantCanonical = "x-app-key=test-app-key=x-signature-algorithm=HMAC-SHA256=x-signature-nonce=0123456789abcdef0123456789abcdef=x-signature-version=1.0=x-timestamp=2022-01-04T03:55:31Z&d9d865cc54ec60678f1b119084ad79ae7f9357d1c4519c6457de3314b7fbba8a"
		wantEncoded   = "x-app-key%3Dtest-app-key%3Dx-signature-algorithm%3DHMAC-SHA256%3Dx-signature-nonce%3D0123456789abcdef0123456789abcdef%3Dx-signature-version%3D1.0%3Dx-timestamp%3D2022-01-04T03%3A55%3A31Z%26d9d865cc54ec60678f1b119084ad79ae7f9357d1c4519c6457de3314b7fbba8a"
		wantSignature = "K6aegSfwYt5gRNoOB9rJlXVgw49y+aUQn+D7u18eGIo="
	)
	body := []byte("test-body")
	params := testSignParams()

	if got := eventCanonicalString(params, body); got != wantCanonical {
		t.Errorf("eventCanonicalString() = %q, want %q", got, wantCanonical)
	}
	if got := percentEncode(wantCanonical); got != wantEncoded {
		t.Errorf("percentEncode() = %q, want %q", got, wantEncoded)
	}
	if got := eventSignature(params, "test-app-secret", body); got != wantSignature {
		t.Errorf("eventSignature() = %q, want %q", got, wantSignature)
	}
}
