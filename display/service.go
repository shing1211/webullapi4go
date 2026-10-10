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

// Package display provides a client for the Webull Display Solution API,
// a separate product from the core Market Data API with its own host,
// authentication, and endpoint paths.
//
// Display Solution endpoints (Corporate Actions, Screener v2, Fund Data)
// require a different auth flow: a client-credentials token obtained via
// HMAC-SHA1 signed request, then used as Authorization: Bearer header.
//
// Reference: https://developer.webull.hk/apis/docs/market-data-api/Hosted-Display-Solution
package display

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/shing1211/webullapi4go/internal/auth"
	errs "github.com/shing1211/webullapi4go/pkg/errors"
	"github.com/shing1211/webullapi4go/pkg/transport"
)

const (
	HKProductionHTTP = "https://co-branding-openapi.webull.hk"
	HKSandboxHTTP    = "https://hk-co-branding-openapi.uat.webullbroker.com"

	displayTokenCreatePath  = "/auth/client-tokens/create"
	displayTokenRefreshPath = "/auth/client-tokens/refresh"
)

// ClientToken is an access/refresh token pair issued by the Display Solution.
type ClientToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresAt        int64  `json:"expires_at"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresAt int64  `json:"refresh_expires_at"`
}

type Service struct {
	appKey    string
	appSecret string
	baseURL   string
	http      *http.Client
	transport *transport.Client

	mu    sync.RWMutex
	token *auth.DisplayToken
}

func NewService(appKey, appSecret string, opts ...Option) *Service {
	cfg := config{baseURL: HKProductionHTTP}
	for _, o := range opts {
		o.apply(&cfg)
	}
	hc := cfg.httpClnt
	if hc == nil {
		hc = &http.Client{Timeout: 30 * time.Second}
	}
	t, _ := transport.New(cfg.baseURL, hc, "webullapi4go")
	return &Service{
		appKey:    appKey,
		appSecret: appSecret,
		baseURL:   cfg.baseURL,
		http:      hc,
		transport: t,
	}
}

type Option interface{ apply(*config) }

type config struct {
	baseURL  string
	httpClnt *http.Client
}

type optionFunc func(*config)

func (f optionFunc) apply(c *config) { f(c) }

func WithBaseURL(baseURL string) Option {
	return optionFunc(func(c *config) { c.baseURL = baseURL })
}

func WithSandbox(v bool) Option {
	return optionFunc(func(c *config) {
		if v {
			c.baseURL = HKSandboxHTTP
		} else {
			c.baseURL = HKProductionHTTP
		}
	})
}

func WithHTTPClient(hc *http.Client) Option {
	return optionFunc(func(c *config) { c.httpClnt = hc })
}

func (s *Service) EnsureToken(ctx context.Context) (string, error) {
	s.mu.RLock()
	cached := s.token
	s.mu.RUnlock()
	if cached != nil && cached.IsValid(time.Now()) {
		return cached.AccessTokenValue(), nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.token != nil && s.token.IsValid(time.Now()) {
		return s.token.AccessTokenValue(), nil
	}
	return s.fetchToken(ctx)
}

func (s *Service) fetchToken(ctx context.Context) (string, error) {
	reqBody := auth.DisplayCreateTokenRequest{ClientUserID: "openapi_client"}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", errs.Wrap(errs.CodeTransport, "encoding display token request", err)
	}

	req, err := s.buildSignedRequest(ctx, http.MethodPost, displayTokenCreatePath, nil, bodyBytes)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.http.Do(req)
	if err != nil {
		return "", errs.Wrap(errs.CodeTransport, displayTokenCreatePath, err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errs.Wrap(errs.CodeTransport, "reading display token response", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", errs.FromHTTPStatus(resp.StatusCode, data)
	}

	var tok auth.DisplayToken
	if err := json.Unmarshal(data, &tok); err != nil {
		return "", errs.Wrap(errs.CodeAPI, "decoding display token response", err)
	}
	s.token = &tok
	return tok.AccessTokenValue(), nil
}

// RefreshClientToken exchanges a refresh token for a new access token using the
// Client-to-Server credentials.
func (s *Service) RefreshClientToken(ctx context.Context, refreshToken string) (*ClientToken, error) {
	bodyBytes, err := json.Marshal(map[string]string{"refresh_token": refreshToken})
	if err != nil {
		return nil, errs.Wrap(errs.CodeTransport, "encoding display refresh request", err)
	}
	req, err := s.buildSignedRequest(ctx, http.MethodPost, displayTokenRefreshPath, nil, bodyBytes)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, errs.Wrap(errs.CodeTransport, displayTokenRefreshPath, err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.Wrap(errs.CodeTransport, "reading display refresh response", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errs.FromHTTPStatus(resp.StatusCode, data)
	}
	var tok ClientToken
	if err := json.Unmarshal(data, &tok); err != nil {
		return nil, errs.Wrap(errs.CodeAPI, "decoding display refresh response", err)
	}
	return &tok, nil
}

func (s *Service) Do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	token, err := s.EnsureToken(ctx)
	if err != nil {
		return err
	}

	bodyBytes, err := s.marshalBody(body)
	if err != nil {
		return err
	}

	reqPath := path
	if len(query) > 0 {
		reqPath += "?" + query.Encode()
	}

	req, err := s.transport.NewRequest(ctx, method, reqPath, nil, bodyBytes)
	if err != nil {
		return errs.Wrap(errs.CodeInvalidConfig, "building request", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.http.Do(req)
	if err != nil {
		return errs.Wrap(errs.CodeTransport, method+" "+path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errs.Wrap(errs.CodeTransport, "reading response body", err)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		s.mu.Lock()
		s.token = nil
		s.mu.Unlock()
		return errs.New(errs.CodeAuth, "display token expired")
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return errs.FromHTTPStatus(resp.StatusCode, data)
	}
	if out == nil || len(data) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return errs.Wrap(errs.CodeAPI, "decoding response body", err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, path string, query url.Values, out any) error {
	return s.Do(ctx, http.MethodGet, path, query, nil, out)
}

func (s *Service) Close() error { return nil }

func (s *Service) buildSignedRequest(ctx context.Context, method, reqPath string, query url.Values, bodyBytes []byte) (*http.Request, error) {
	u, err := url.Parse(s.baseURL)
	if err != nil {
		return nil, err
	}
	host := u.Hostname()

	nonce, err := newNonce()
	if err != nil {
		return nil, errs.Wrap(errs.CodeAuth, "generating nonce", err)
	}
	timestamp := time.Now().UTC().Format(auth.TimestampFormat)

	sig, err := s.signDisplay(method, reqPath, query, host, timestamp, nonce, bodyBytes)
	if err != nil {
		return nil, errs.Wrap(errs.CodeAuth, "signing display token request", err)
	}

	fullURL := s.baseURL + reqPath
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}
	var reader io.Reader
	if len(bodyBytes) > 0 {
		reader = bytes.NewReader(bodyBytes)
	}
	req, err := http.NewRequestWithContext(ctx, method, fullURL, reader)
	if err != nil {
		return nil, err
	}
	req.Host = host
	req.Header.Set("x-app-key", s.appKey)
	req.Header.Set("x-app-secret", s.appSecret)
	req.Header.Set("x-timestamp", timestamp)
	req.Header.Set("x-signature-nonce", nonce)
	req.Header.Set("x-signature-version", "1.0")
	req.Header.Set("x-signature-algorithm", "HMAC-SHA1")
	req.Header.Set("x-signature", sig)
	return req, nil
}

func (s *Service) signDisplay(method, reqPath string, query url.Values, host, timestamp, nonce string, body []byte) (string, error) {
	entries := []string{
		"x-app-key=" + s.appKey,
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
	canonical := percentEncode(strings.Join(entries, "&"))

	mac := hmac.New(sha1.New, []byte(s.appSecret+"&"))
	mac.Write([]byte(canonical))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func newNonce() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
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
		b.WriteByte(upperhex[c&0x0f])
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

func (s *Service) marshalBody(body any) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	if raw, ok := body.([]byte); ok {
		return raw, nil
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(body); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}
