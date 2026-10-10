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

package display_test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/shing1211/webullapi4go/display"
)

const (
	testAppKey    = "display-test-key"
	testAppSecret = "display-test-secret"
)

func newTestService(srv *httptest.Server) *display.Service {
	return display.NewService(testAppKey, testAppSecret,
		display.WithBaseURL(srv.URL),
		display.WithHTTPClient(srv.Client()),
	)
}

func setupTestServer(t *testing.T) (*httptest.Server, *display.Service, *captureVars) {
	t.Helper()
	cv := &captureVars{}

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/client-tokens/create", func(w http.ResponseWriter, r *http.Request) {
		handleTokenRequest(cv, w, r)
	})
	mux.HandleFunc("/market-data/instruments/stocks/corporate-actions/list", func(w http.ResponseWriter, r *http.Request) {
		handleAPIRequest(cv, w, r)
	})
	mux.HandleFunc("/some/path", func(w http.ResponseWriter, r *http.Request) {
		handleAPIRequest(cv, w, r)
	})
	mux.HandleFunc("/any/path", func(w http.ResponseWriter, r *http.Request) {
		handleAPIRequest(cv, w, r)
	})
	mux.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {
		handlePathEndpoint(cv, w, r)
	})

	srv := httptest.NewServer(mux)
	return srv, newTestService(srv), cv
}

type captureVars struct {
	authHeader  string
	contentType string
	queryRaw    string
	body        map[string]string
	callCount   int
	lastStatus  int
}

func handleTokenRequest(cv *captureVars, w http.ResponseWriter, r *http.Request) {
	cv.callCount++
	cv.authHeader = r.Header.Get("Authorization")
	cv.contentType = r.Header.Get("Content-Type")

	if r.Header.Get("x-app-key") != testAppKey {
		http.Error(w, "wrong app key", http.StatusUnauthorized)
		return
	}
	timestamp := r.Header.Get("x-timestamp")
	nonce := r.Header.Get("x-signature-nonce")
	if r.Header.Get("x-signature-version") != "1.0" || r.Header.Get("x-signature-algorithm") != "HMAC-SHA1" {
		http.Error(w, "bad signature params", http.StatusBadRequest)
		return
	}

	body, _ := io.ReadAll(r.Body)
	gotSig := r.Header.Get("x-signature")
	wantSig := signDisplayToken("POST", "/auth/client-tokens/create", nil, r.Host, timestamp, nonce, body, testAppSecret)
	if !hmac.Equal([]byte(gotSig), []byte(wantSig)) {
		http.Error(w, "signature mismatch", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"access_token":       "test-access-token",
		"expires_at":         time.Now().Add(1 * time.Hour).UnixMilli(),
		"refresh_token":      "test-refresh-token",
		"refresh_expires_at": time.Now().Add(24 * time.Hour).UnixMilli(),
	})
}

func handleAPIRequest(cv *captureVars, w http.ResponseWriter, r *http.Request) {
	cv.callCount++
	cv.authHeader = r.Header.Get("Authorization")
	cv.queryRaw = r.URL.RawQuery

	if cv.authHeader != "Bearer test-access-token" {
		cv.lastStatus = http.StatusUnauthorized
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	cv.lastStatus = http.StatusOK
	json.NewEncoder(w).Encode(map[string]any{
		"data": []map[string]any{
			{
				"instrument_id": 913256135,
				"symbol":        "AAPL",
				"exchange_code": "NSQ",
				"event_type":    "DIVIDEND",
				"event_action":  "DISTRIBUTE",
				"event_id":      12345,
				"source":        "WEBULL",
				"ratio_old":     "1.0",
				"ratio_new":     "1.0",
				"event_date":    "2026-09-19",
				"update_time":   "2026-09-19T10:00:00Z",
				"create_time":   "2026-09-19T10:00:00Z",
			},
		},
		"pagination_key": "next-page-key",
	})
}

func handlePathEndpoint(cv *captureVars, w http.ResponseWriter, r *http.Request) {
	cv.callCount++
	cv.authHeader = r.Header.Get("Authorization")
	cv.contentType = r.Header.Get("Content-Type")

	if r.Header.Get("Content-Type") == "application/json" && r.Body != nil {
		var got map[string]string
		json.NewDecoder(r.Body).Decode(&got)
		cv.body = got
	}

	w.Header().Set("Content-Type", "application/json")
	cv.lastStatus = http.StatusOK
	w.Write([]byte("{}"))
}

func signDisplayToken(method, reqPath string, query url.Values, host, timestamp, nonce string, body []byte, secret string) string {
	entries := []string{
		"x-app-key=" + testAppKey,
		"x-signature-algorithm=HMAC-SHA1",
		"x-signature-version=1.0",
		"x-signature-nonce=" + nonce,
		"x-timestamp=" + timestamp,
		"host=" + host,
	}
	if len(body) > 0 {
		h := md5.Sum(body)
		entries = append(entries, hex.EncodeToString(h[:]))
	}
	entries = append(entries, reqPath)
	canonical := percentEncodeTest(strings.Join(entries, "&"))
	mac := hmac.New(sha1.New, []byte(secret+"&"))
	mac.Write([]byte(canonical))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func percentEncodeTest(s string) string {
	const upperhex = "0123456789ABCDEF"
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isUnreservedTest(c) {
			b.WriteByte(c)
			continue
		}
		b.WriteByte('%')
		b.WriteByte(upperhex[c>>4])
		b.WriteByte(upperhex[c&0x0f])
	}
	return b.String()
}

func isUnreservedTest(c byte) bool {
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

func TestNewService_happy(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()
	svc := newTestService(srv)
	if svc == nil {
		t.Fatal("NewService returned nil")
	}
}

func TestNewService_defaultBaseURL(t *testing.T) {
	t.Parallel()
	svc := display.NewService(testAppKey, testAppSecret)
	if svc == nil {
		t.Fatal("NewService returned nil")
	}
}

func TestWithSandbox_switchesBaseURL(t *testing.T) {
	t.Parallel()
	svc := display.NewService(testAppKey, testAppSecret, display.WithSandbox(true))
	if svc == nil {
		t.Fatal("NewService with WithSandbox returned nil")
	}
}

func TestWithBaseURL_setsBaseURL(t *testing.T) {
	t.Parallel()
	svc := display.NewService(testAppKey, testAppSecret, display.WithBaseURL("https://custom.example.com"))
	if svc == nil {
		t.Fatal("NewService with WithBaseURL returned nil")
	}
}

func TestWithHTTPClient_setsHTTPClient(t *testing.T) {
	t.Parallel()
	hc := &http.Client{Timeout: 10 * time.Second}
	svc := display.NewService(testAppKey, testAppSecret, display.WithHTTPClient(hc))
	if svc == nil {
		t.Fatal("NewService with WithHTTPClient returned nil")
	}
}

func TestEnsureToken_fetchesAndCaches(t *testing.T) {
	t.Parallel()
	srv, svc, cv := setupTestServer(t)
	defer srv.Close()

	token1, err := svc.EnsureToken(context.Background())
	if err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}
	if token1 != "test-access-token" {
		t.Errorf("token = %q, want %q", token1, "test-access-token")
	}

	token2, err := svc.EnsureToken(context.Background())
	if err != nil {
		t.Fatalf("EnsureToken() second call error = %v", err)
	}
	if token2 != token1 {
		t.Errorf("cached token = %q, want same as first %q", token2, token1)
	}
	if cv.callCount != 1 {
		t.Errorf("callCount = %d, want 1 (second call should be cached)", cv.callCount)
	}
}

func TestEnsureToken_expiredTokenForcesRefresh(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/client-tokens/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "new-token",
			"expires_at":    time.Now().Add(1 * time.Hour).UnixMilli(),
			"refresh_token": "new-refresh",
		})
	})
	mux.HandleFunc("/any/path", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"data": []any{}})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	svc := newTestService(srv)

	_, err := svc.EnsureToken(context.Background())
	if err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}
}

func TestEnsureToken_noExpiryAlwaysValid(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/client-tokens/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "no-expiry-token",
			"expires_at":    0,
			"refresh_token": "refresh",
		})
	})
	mux.HandleFunc("/any/path", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": []any{}})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	svc := newTestService(srv)

	token1, err := svc.EnsureToken(context.Background())
	if err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}
	if token1 != "no-expiry-token" {
		t.Errorf("token = %q, want no-expiry-token", token1)
	}
}

func TestDo_getRequest_includesBearerToken(t *testing.T) {
	t.Parallel()
	srv, svc, cv := setupTestServer(t)
	defer srv.Close()

	err := svc.Get(context.Background(), "/some/path", nil, nil)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if cv.authHeader != "Bearer test-access-token" {
		t.Errorf("Authorization = %q, want %q", cv.authHeader, "Bearer test-access-token")
	}
}

func TestDo_postRequest_includesBearerToken(t *testing.T) {
	t.Parallel()
	srv, svc, cv := setupTestServer(t)
	defer srv.Close()

	err := svc.Do(context.Background(), http.MethodPost, "/some/path", nil, map[string]string{"key": "value"}, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if cv.authHeader != "Bearer test-access-token" {
		t.Errorf("Authorization = %q, want %q", cv.authHeader, "Bearer test-access-token")
	}
	if cv.contentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", cv.contentType, "application/json")
	}
}

func TestDo_401_clearsToken(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	callCount := 0
	mux.HandleFunc("/auth/client-tokens/create", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "token-" + itoa(callCount),
			"expires_at":    time.Now().Add(1 * time.Hour).UnixMilli(),
			"refresh_token": "refresh",
		})
	})
	mux.HandleFunc("/any/path", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	svc := newTestService(srv)

	_, err := svc.EnsureToken(context.Background())
	if err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}

	err = svc.Get(context.Background(), "/any/path", nil, nil)
	if err == nil {
		t.Fatal("Get() expected error after 401, got nil")
	}
}

func TestDo_non2xx_returnsError(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/client-tokens/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "valid-token",
			"expires_at":    time.Now().Add(1 * time.Hour).UnixMilli(),
			"refresh_token": "refresh",
		})
	})
	mux.HandleFunc("/any/path", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	svc := newTestService(srv)

	_, err := svc.EnsureToken(context.Background())
	if err != nil {
		t.Fatalf("EnsureToken() error = %v", err)
	}
	err = svc.Get(context.Background(), "/any/path", nil, nil)
	if err == nil {
		t.Fatal("Get() expected error, got nil")
	}
}

func TestClose_returnsNil(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()
	svc := newTestService(srv)
	if err := svc.Close(); err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}

func TestDo_getRequest_queryParams(t *testing.T) {
	t.Parallel()
	srv, svc, cv := setupTestServer(t)
	defer srv.Close()

	q := url.Values{"symbol": []string{"AAPL"}, "category": []string{"US"}}
	err := svc.Get(context.Background(), "/some/path", q, nil)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if cv.queryRaw == "" {
		t.Error("query params not sent")
	}
	if cv.queryRaw != "category=US&symbol=AAPL" {
		t.Errorf("query = %q, want %q", cv.queryRaw, "category=US&symbol=AAPL")
	}
}

func TestDo_responseBody_unmarshaled(t *testing.T) {
	t.Parallel()
	srv, svc, _ := setupTestServer(t)
	defer srv.Close()

	type result struct {
		Data []struct {
			Symbol string `json:"symbol"`
		} `json:"data"`
	}
	var out result
	err := svc.Get(context.Background(), "/market-data/instruments/stocks/corporate-actions/list", nil, &out)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if len(out.Data) != 1 || out.Data[0].Symbol != "AAPL" {
		t.Errorf("out = %+v, want data with AAPL", out)
	}
}

func TestDo_nilBody(t *testing.T) {
	t.Parallel()
	srv, svc, _ := setupTestServer(t)
	defer srv.Close()

	err := svc.Do(context.Background(), http.MethodGet, "/some/path", nil, nil, nil)
	if err != nil {
		t.Fatalf("Do() nil body error = %v", err)
	}
}

func TestDo_emptyResponse(t *testing.T) {
	t.Parallel()
	srv, svc, _ := setupTestServer(t)
	defer srv.Close()

	err := svc.Do(context.Background(), http.MethodGet, "/some/path", nil, nil, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
}

func TestMarshalBody_escapesHTML(t *testing.T) {
	t.Parallel()
	srv, svc, cv := setupTestServer(t)
	defer srv.Close()

	err := svc.Do(context.Background(), http.MethodPost, "/path", nil, map[string]string{"name": "AAPL <B>"}, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if cv.body == nil {
		t.Fatal("body not captured")
	}
	if cv.body["name"] != "AAPL \u003cB\u003e" {
		t.Errorf("html escaped = %q, want %q", cv.body["name"], "AAPL \u003cB\u003e")
	}
}

func TestMarshalBody_trimsTrailingNewline(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if bytes.HasSuffix(body, []byte("\n")) {
			t.Errorf("body ends with \\n: %q", string(body))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}))
	defer srv.Close()
	svc := newTestService(srv)
	err := svc.Do(context.Background(), http.MethodPost, "/path", nil, map[string]string{"key": "value"}, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[pos:])
}

func TestRefreshClientToken(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/auth/client-tokens/refresh"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), "rt1") {
			t.Errorf("body = %s", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"at2","expires_at":123,"refresh_token":"rt2","refresh_expires_at":456}`))
	}))
	defer srv.Close()

	svc := newTestService(srv)
	tok, err := svc.RefreshClientToken(context.Background(), "rt1")
	if err != nil {
		t.Fatalf("RefreshClientToken() error = %v", err)
	}
	if tok.AccessToken != "at2" || tok.RefreshToken != "rt2" || tok.ExpiresAt != 123 {
		t.Fatalf("token = %+v", tok)
	}
}
