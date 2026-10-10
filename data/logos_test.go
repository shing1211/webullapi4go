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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shing1211/webullapi4go/client"
	"github.com/shing1211/webullapi4go/data"
	"github.com/shing1211/webullapi4go/display"
)

func newLogosTestServer(t *testing.T) (*httptest.Server, *display.Service) {
	t.Helper()
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
	mux.HandleFunc("/market-data/fundamentals/logos/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]string{
			{"symbol": "AAPL", "logo": "https://logo.example.com/AAPL.png"},
			{"symbol": "TSLA", "logo": "https://logo.example.com/TSLA.png"},
		})
	})
	srv := httptest.NewServer(mux)
	dsvc := display.NewService(testAppKey, testAppSecret,
		display.WithBaseURL(srv.URL),
		display.WithHTTPClient(srv.Client()),
	)
	return srv, dsvc
}

func TestGetLogos_pathAndQuery(t *testing.T) {
	t.Parallel()
	srv, dsvc := newLogosTestServer(t)
	defer srv.Close()

	cl, _ := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithBaseURL(srv.URL),
	)
	dc := data.New(cl)
	dc.SetDisplaySvcForTesting(dsvc)

	_, err := dc.GetLogos(context.Background(), data.LogoQuery{Symbols: []string{"AAPL", "TSLA"}})
	if err != nil {
		t.Fatalf("GetLogos() error = %v", err)
	}
}

func TestGetLogos_responseFields(t *testing.T) {
	t.Parallel()
	srv, dsvc := newLogosTestServer(t)
	defer srv.Close()

	cl, _ := client.New(
		client.WithAppKey(testAppKey),
		client.WithAppSecret(testAppSecret),
		client.WithBaseURL(srv.URL),
	)
	dc := data.New(cl)
	dc.SetDisplaySvcForTesting(dsvc)

	logos, err := dc.GetLogos(context.Background(), data.LogoQuery{Symbols: []string{"AAPL", "TSLA"}})
	if err != nil {
		t.Fatalf("GetLogos() error = %v", err)
	}
	if len(logos) != 2 {
		t.Fatalf("len(logos) = %d, want 2", len(logos))
	}
	if logos[0].Symbol != "AAPL" || logos[0].Logo != "https://logo.example.com/AAPL.png" {
		t.Errorf("logo[0] = %+v", logos[0])
	}
	if logos[1].Symbol != "TSLA" {
		t.Errorf("logo[1] = %+v", logos[1])
	}
}
