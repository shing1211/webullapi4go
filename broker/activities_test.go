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

package broker_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shing1211/webullapi4go/broker"
)

func TestGetCashActivities(t *testing.T) {
	t.Parallel()

	const body = `[{"activity_id":"ACT1","account_id":"ACC1","type":"DEPOSIT",` +
		`"amount":"1000.00","currency":"HKD","status":"COMPLETED",` +
		`"create_time":"2026-01-15T10:30:00Z","description":"Initial deposit"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/broker/activities/list"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("account_id"), "ACC1"; got != want {
			t.Errorf("account_id = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := broker.New(newTestClient(t, srv.URL))
	got, err := c.GetCashActivities(context.Background(), "ACC1")
	if err != nil {
		t.Fatalf("GetCashActivities() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d activities, want 1", len(got))
	}
	if got[0].ActivityID != "ACT1" || got[0].Type != "DEPOSIT" {
		t.Errorf("activity = %+v", got[0])
	}
	if got[0].Amount != "1000.00" || got[0].Currency != "HKD" {
		t.Errorf("amount/currency = %+v", got[0])
	}
}
