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

package brokerfd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shing1211/webullapi4go/client"
)

// Two servers, deliberately: one stands for the Broker host and one for the core
// host, and the tests assert which of them is reached. A single test server would
// be green whether the package routed to the Broker host or the core host, which
// is the state this file exists to rule out - brokerfd sent every request to the
// core host while its sibling broker/ module used the Broker host, and nothing in
// either package could tell the difference.
type hostRecorder struct {
	brokerHits int
	coreHits   int
}

func newTwoHostClient(t *testing.T) (*Client, *hostRecorder, func()) {
	t.Helper()
	rec := &hostRecorder{}

	brokerSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.brokerHits++
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"data": []any{}})
	}))
	coreSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.coreHits++
		w.WriteHeader(http.StatusNotFound)
	}))

	core, err := client.New(
		client.WithCredentials(testAppKey, testAppSecret),
		client.WithAutoToken(false),
		client.WithEndpoints(client.Endpoints{
			HTTP:       coreSrv.URL,
			BrokerHTTP: brokerSrv.URL,
		}),
	)
	if err != nil {
		t.Fatalf("client.New: %v", err)
	}
	return New(core), rec, func() {
		brokerSrv.Close()
		coreSrv.Close()
		_ = core.Close()
	}
}

func assertBrokerHostOnly(t *testing.T, rec *hostRecorder, wantBroker int) {
	t.Helper()
	if rec.brokerHits != wantBroker {
		t.Errorf("Broker host hits = %d, want %d: the request did not go to the Broker endpoint",
			rec.brokerHits, wantBroker)
	}
	if rec.coreHits != 0 {
		t.Errorf("core host hits = %d, want 0: a Broker FD request reached the core host", rec.coreHits)
	}
}

// TestGetReachesTheBrokerHost is the assertion itself, through the get helper.
// Against the pre-fix code this fails: the core server answers 404 and
// brokerHits stays 0.
func TestGetReachesTheBrokerHost(t *testing.T) {
	c, rec, done := newTwoHostClient(t)
	defer done()

	// The core host answers 404, so a nil error is itself evidence the request
	// did not go there.
	if _, err := c.GetFDAssetsSummary(context.Background(), "A1"); err != nil {
		t.Fatalf("GetFDAssetsSummary: %v", err)
	}
	assertBrokerHostOnly(t, rec, 1)
}

// TestPostReachesTheBrokerHost covers the post helper, so the assertion does not
// rest on whichever single path get happened to take. The package has no put
// helper, so get and post are the whole surface.
func TestPostReachesTheBrokerHost(t *testing.T) {
	c, rec, done := newTwoHostClient(t)
	defer done()

	req := AddFDAchAccountRequest{
		AccountID:     "A1",
		BankName:      "Test Bank",
		AccountNumber: "000123456789",
	}
	if _, err := c.AddFDAchAccount(context.Background(), req); err != nil {
		t.Fatalf("AddFDAchAccount: %v", err)
	}
	assertBrokerHostOnly(t, rec, 1)
}

// TestEveryHelperReachesTheBrokerHost exercises both helpers in one call so a
// regression in either is caught by a single test.
func TestEveryHelperReachesTheBrokerHost(t *testing.T) {
	c, rec, done := newTwoHostClient(t)
	defer done()

	ctx := context.Background()
	// Two gets and a post, all decoding into a struct, because the stand-in
	// Broker host answers {"data":[]} for every path and a method that decodes a
	// slice would fail on the shape rather than on the routing under test.
	if _, err := c.GetFDAccountDetail(ctx, "A1"); err != nil {
		t.Fatalf("GetFDAccountDetail: %v", err)
	}
	if _, err := c.GetFDAssetsSummary(ctx, "A1"); err != nil {
		t.Fatalf("GetFDAssetsSummary: %v", err)
	}
	if _, err := c.AddFDAchAccount(ctx, AddFDAchAccountRequest{AccountID: "A1"}); err != nil {
		t.Fatalf("AddFDAchAccount: %v", err)
	}
	assertBrokerHostOnly(t, rec, 3)
}

// TestBrokerHostUnconfiguredIsAConfigError records the one behaviour change the
// routing fix introduces, because it is reachable: DoBroker returns a typed
// configuration error where the old code quietly sent the request to the core
// host and returned that server's 404. A caller who configures Endpoints by hand
// and omits BrokerHTTP now gets a diagnosis naming the missing field instead of a
// status code that reads like an absent endpoint.
//
// client.EndpointsFor populates BrokerHTTP for every region and environment, so
// this only arises from a hand-built Endpoints value.
func TestBrokerHostUnconfiguredIsAConfigError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("the core host was called, which is the pre-fix behaviour: %s", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	core, err := client.New(
		client.WithCredentials(testAppKey, testAppSecret),
		client.WithAutoToken(false),
		client.WithEndpoints(client.Endpoints{HTTP: srv.URL}),
	)
	if err != nil {
		t.Fatalf("client.New: %v", err)
	}
	defer func() { _ = core.Close() }()

	_, err = New(core).ListFDAccounts(context.Background(), "")
	if err == nil {
		t.Fatal("want an error when BrokerHTTP is unconfigured, got nil")
	}
	if !strings.Contains(err.Error(), "broker HTTP endpoint") {
		t.Errorf("error = %v, want one naming the missing broker HTTP endpoint so a caller can diagnose it", err)
	}
}

// TestEndpointsForAlwaysSuppliesBrokerHost is the guard on the risk the fix
// introduces. If any region or environment ever left BrokerHTTP empty, routing
// brokerfd through DoBroker would make the whole package unusable by default
// rather than merely wrong, so the table is asserted rather than assumed.
func TestEndpointsForAlwaysSuppliesBrokerHost(t *testing.T) {
	for _, r := range []client.Region{"hk", "us", "jp"} {
		for _, env := range []client.Environment{"production", "sandbox"} {
			eps := client.EndpointsFor(r, env)
			if eps.BrokerHTTP == "" {
				t.Errorf("region %q environment %q: BrokerHTTP is empty, so brokerfd would "+
					"fail to route at all rather than route wrongly", r, env)
			}
			if eps.BrokerHTTP == eps.HTTP {
				t.Errorf("region %q environment %q: BrokerHTTP equals HTTP (%q), so routing "+
					"through the Broker transport would be indistinguishable from the core host",
					r, env, eps.HTTP)
			}
		}
	}
}
