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

package data_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/data"
	"github.com/shing1211/webullapi4go/display"
)

func TestDSSubscribe(t *testing.T) {
	t.Parallel()
	var capturedBody []byte
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/client-tokens/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":       "test-access-token",
			"expires_at":         0,
			"refresh_token":      "refresh",
			"refresh_expires_at": 0,
		})
	})
	mux.HandleFunc("/market-data/streaming/subscribe", func(w http.ResponseWriter, r *http.Request) {
		capturedBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	dsvc := display.NewService(testAppKey, testAppSecret,
		display.WithBaseURL(srv.URL),
		display.WithHTTPClient(srv.Client()),
	)
	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	defer cl.Close()
	dc := data.New(cl)
	dc.SetDisplaySvcForTesting(dsvc)

	err = dc.DSSubscribe(context.Background(), data.DSSubscribeRequest{
		Symbols: []string{"AAPL", "GOOG"},
	})
	if err != nil {
		t.Fatalf("DSSubscribe() error = %v", err)
	}

	var got data.DSSubscribeRequest
	if err := json.Unmarshal(capturedBody, &got); err != nil {
		t.Fatalf("decoding captured body: %v", err)
	}
	if len(got.Symbols) != 2 || got.Symbols[0] != "AAPL" || got.Symbols[1] != "GOOG" {
		t.Errorf("body symbols = %v, want [AAPL GOOG]", got.Symbols)
	}
}

func TestDSUnsubscribe(t *testing.T) {
	t.Parallel()
	var capturedBody []byte
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/client-tokens/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":       "test-access-token",
			"expires_at":         0,
			"refresh_token":      "refresh",
			"refresh_expires_at": 0,
		})
	})
	mux.HandleFunc("/market-data/streaming/unsubscribe", func(w http.ResponseWriter, r *http.Request) {
		capturedBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	dsvc := display.NewService(testAppKey, testAppSecret,
		display.WithBaseURL(srv.URL),
		display.WithHTTPClient(srv.Client()),
	)
	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	defer cl.Close()
	dc := data.New(cl)
	dc.SetDisplaySvcForTesting(dsvc)

	err = dc.DSUnsubscribe(context.Background(), data.DSUnsubscribeRequest{
		Symbols: []string{"AAPL"},
	})
	if err != nil {
		t.Fatalf("DSUnsubscribe() error = %v", err)
	}

	var got data.DSUnsubscribeRequest
	if err := json.Unmarshal(capturedBody, &got); err != nil {
		t.Fatalf("decoding captured body: %v", err)
	}
	if len(got.Symbols) != 1 || got.Symbols[0] != "AAPL" {
		t.Errorf("body symbols = %v, want [AAPL]", got.Symbols)
	}
}

func TestDSSubscribe_emptySymbols(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/client-tokens/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":       "test-access-token",
			"expires_at":         0,
			"refresh_token":      "refresh",
			"refresh_expires_at": 0,
		})
	})
	mux.HandleFunc("/market-data/streaming/subscribe", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	dsvc := display.NewService(testAppKey, testAppSecret,
		display.WithBaseURL(srv.URL),
		display.WithHTTPClient(srv.Client()),
	)
	cl, err := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
	)
	if err != nil {
		t.Fatalf("client.New() error = %v", err)
	}
	defer cl.Close()
	dc := data.New(cl)
	dc.SetDisplaySvcForTesting(dsvc)

	err = dc.DSSubscribe(context.Background(), data.DSSubscribeRequest{
		Symbols: []string{},
	})
	if err != nil {
		t.Fatalf("DSSubscribe() with empty symbols: %v", err)
	}
}
