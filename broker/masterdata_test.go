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

package broker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTradeCalendar(t *testing.T) {
	t.Parallel()

	const body = `[{"date":"2026-01-02","market":"HK","status":"TRADING",` +
		`"open_time":"09:30","close_time":"16:00"}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got, want := r.URL.Path, "/openapi/v1/broker/master-data/trade-calendar"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("market"), "HK"; got != want {
			t.Errorf("market = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("start_date"), "2026-01-01"; got != want {
			t.Errorf("start_date = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("end_date"), "2026-01-31"; got != want {
			t.Errorf("end_date = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := newTestBrokerClient(t, srv.URL)
	got, err := c.GetTradeCalendar(context.Background(), "HK", "2026-01-01", "2026-01-31")
	if err != nil {
		t.Fatalf("GetTradeCalendar() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d entries, want 1", len(got))
	}
	if got[0].Date != "2026-01-02" || got[0].Market != "HK" {
		t.Errorf("calendar = %+v", got[0])
	}
}
