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

// Package connect provides a client for the Webull Connect API, an OAuth 2.0
// flow for third-party applications. The authorization-code step is performed
// by redirecting the user's browser to the URL returned by
// [Client.AuthorizationURL]; the code is then exchanged for tokens with
// [Client.CreateToken].
//
// Reference: https://developer.webull.com/apis/docs/connect-api/about-connect-api
package connect

import (
	"context"
	"net/url"

	"github.com/shing1211/webullapi4go/client"
)

const (
	pathAuthorizationCode = "/oauth2/auth-codes/get"
	pathTokenCreate       = "/oauth2/tokens/create"

	// ResponseTypeCode is the OAuth 2.0 authorization-code response type.
	ResponseTypeCode = "code"
	// GrantTypeAuthorizationCode exchanges an authorization code for tokens.
	GrantTypeAuthorizationCode = "authorization_code"
	// GrantTypeRefreshToken exchanges a refresh token for new tokens.
	GrantTypeRefreshToken = "refresh_token"
)

// Client is the Webull Connect (OAuth) API client.
type Client struct {
	core *client.Client
}

// New binds a Connect API client to an existing core client. The core client is
// owned by the caller and is not closed by [Client.Close].
func New(c *client.Client) *Client {
	return &Client{core: c}
}

// Core returns the underlying core client.
func (c *Client) Core() *client.Client { return c.core }

// Close releases resources held by the Connect client. The underlying core
// client is not closed.
func (c *Client) Close() error { return nil }

// AuthorizationParams parameterizes [Client.AuthorizationURL].
type AuthorizationParams struct {
	// ResponseType is the OAuth response type. Empty means [ResponseTypeCode].
	ResponseType string
	// ClientID is the OAuth client identifier. Required.
	ClientID string
	// Scope is the requested scope. Required.
	Scope string
	// State is an opaque anti-CSRF value echoed back by the server. Required.
	State string
	// RedirectURI is the registered callback URI. Required.
	RedirectURI string
}

// AuthorizationURL builds the browser URL that starts the OAuth 2.0
// authorization-code flow. The user is redirected here to grant access; the
// server then redirects back to RedirectURI with the authorization code.
func (c *Client) AuthorizationURL(p AuthorizationParams) string {
	rt := p.ResponseType
	if rt == "" {
		rt = ResponseTypeCode
	}
	q := url.Values{}
	q.Set("response_type", rt)
	q.Set("client_id", p.ClientID)
	q.Set("scope", p.Scope)
	q.Set("state", p.State)
	q.Set("redirect_uri", p.RedirectURI)
	return c.core.Endpoints().HTTP + pathAuthorizationCode + "?" + q.Encode()
}

// TokenRequest parameterizes [Client.CreateToken].
type TokenRequest struct {
	// GrantType is [GrantTypeAuthorizationCode] or [GrantTypeRefreshToken].
	// Required.
	GrantType string
	// Code is the authorization code, required for the authorization-code grant.
	Code string
	// RedirectURI is the callback URI, required for the authorization-code grant.
	RedirectURI string
	// RefreshToken is the refresh token, required for the refresh-token grant.
	RefreshToken string
}

// Token is an OAuth 2.0 token pair returned by [Client.CreateToken].
type Token struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        string `json:"expires_in"`
	RefreshExpiresIn string `json:"rt_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	CreatedAt        string `json:"created_at"`
	IdentityID       string `json:"identity_id"`
}

// CreateToken exchanges an authorization code or refresh token for an access
// token via the signed Connect API.
func (c *Client) CreateToken(ctx context.Context, req TokenRequest) (*Token, error) {
	body := struct {
		GrantType    string `json:"grant_type"`
		Code         string `json:"code,omitempty"`
		RedirectURI  string `json:"redirect_uri,omitempty"`
		RefreshToken string `json:"refresh_token,omitempty"`
	}{
		GrantType:    req.GrantType,
		Code:         req.Code,
		RedirectURI:  req.RedirectURI,
		RefreshToken: req.RefreshToken,
	}
	var out Token
	if err := c.core.Do(ctx, "POST", pathTokenCreate, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
