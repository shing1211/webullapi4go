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

package connect_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/connect"
)

func newTestConnect(t *testing.T, baseURL string) *connect.Client {
	t.Helper()
	core, err := client.New(
		client.WithBaseURL(baseURL),
		client.WithAppKey("test-key"),
		client.WithAppSecret("test-secret"),
	)
	if err != nil {
		t.Fatalf("client.New error = %v", err)
	}
	return connect.New(core)
}

func TestAuthorizationURL(t *testing.T) {
	t.Parallel()
	c := newTestConnect(t, "https://api.example.com")
	got := c.AuthorizationURL(connect.AuthorizationParams{
		ClientID:    "cid",
		Scope:       "trade",
		State:       "xyz",
		RedirectURI: "https://app/cb",
	})
	u, err := url.Parse(got)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if u.Path != "/oauth2/auth-codes/get" {
		t.Errorf("path = %q", u.Path)
	}
	q := u.Query()
	if q.Get("response_type") != "code" || q.Get("client_id") != "cid" ||
		q.Get("scope") != "trade" || q.Get("state") != "xyz" || q.Get("redirect_uri") != "https://app/cb" {
		t.Errorf("query = %v", q)
	}
}

func TestCreateToken(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/oauth2/tokens/create"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		var req map[string]string
		_ = json.Unmarshal(body, &req)
		if req["grant_type"] != "authorization_code" || req["code"] != "abc" {
			t.Errorf("body = %s", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"at","token_type":"Bearer","expires_in":"3600",` +
			`"rt_expires_in":"7200","refresh_token":"rt","created_at":"t","identity_id":"id"}`))
	}))
	defer srv.Close()

	c := newTestConnect(t, srv.URL)
	tok, err := c.CreateToken(context.Background(), connect.TokenRequest{
		GrantType:   connect.GrantTypeAuthorizationCode,
		Code:        "abc",
		RedirectURI: "https://app/cb",
	})
	if err != nil {
		t.Fatalf("CreateToken() error = %v", err)
	}
	if tok.AccessToken != "at" || tok.RefreshToken != "rt" || tok.TokenType != "Bearer" {
		t.Fatalf("token = %+v", tok)
	}
}
