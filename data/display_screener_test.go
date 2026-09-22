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

	"github.com/shing1211/webullapi4go/data"
	"github.com/shing1211/webullapi4go/display"
)

const dsScreenerBody = `[{"instrument_id":"913256135","symbol":"AAPL",` +
	`"name":"Apple Inc.","exchange_code":"NSQ","currency_code":"USD",` +
	`"pre_close":"380.2","open":"382.0","high":"388.6","low":"380.0",` +
	`"close":"385.6","price":"385.6","change":"5.4","change_ratio":"0.0142",` +
	`"volume":"12345678","turnover":"4756789012","turnover_rate":"0.0013",` +
	`"market_value":"3650000000000","amplitude":"0.0226",` +
	`"relative_volume_10d":"11.91"}]`

type dsScreenerVars struct {
	path      string
	rankType  string
	category  string
	sortBy    string
	direction string
}

func newDisplayScreenerTestServer(t *testing.T) (*httptest.Server, *display.Service, *dsScreenerVars) {
	t.Helper()
	sv := &dsScreenerVars{}
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/client-tokens/create", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":       "test-access-token",
			"expires_at":         0,
			"refresh_token":      "refresh",
			"refresh_expires_at": 0,
		})
	})

	mux.HandleFunc("/market-data/screeners/gainers-losers/list", func(w http.ResponseWriter, r *http.Request) {
		sv.path = r.URL.Path
		q := r.URL.Query()
		sv.rankType = q.Get("rank_type")
		sv.category = q.Get("category")
		sv.sortBy = q.Get("sort_by")
		sv.direction = q.Get("direction")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(dsScreenerBody))
	})

	mux.HandleFunc("/market-data/screeners/top-actives/list", func(w http.ResponseWriter, r *http.Request) {
		sv.path = r.URL.Path
		q := r.URL.Query()
		sv.category = q.Get("category")
		sv.rankType = q.Get("rank_type")
		sv.sortBy = q.Get("sort_by")
		sv.direction = q.Get("direction")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(dsScreenerBody))
	})

	srv := httptest.NewServer(mux)
	dsvc := display.NewService(testAppKey, testAppSecret,
		display.WithBaseURL(srv.URL),
		display.WithHTTPClient(srv.Client()),
	)
	return srv, dsvc, sv
}

func TestGetDisplayGainersLosers(t *testing.T) {
	t.Parallel()
	srv, dsvc, sv := newDisplayScreenerTestServer(t)
	defer srv.Close()
	dc := newDataClientWithDisplay(t, dsvc)

	got, err := dc.GetDisplayGainersLosers(context.Background(), data.GainersLosersQuery{
		RankType:  data.GainersLosersRankDay1,
		Category:  data.StockCategoryUS,
		SortBy:    data.ScreenerSortChangeRatio,
		Direction: data.SortDirectionAsc,
	})
	if err != nil {
		t.Fatalf("GetDisplayGainersLosers() error = %v", err)
	}
	if sv.path != "/market-data/screeners/gainers-losers/list" {
		t.Errorf("path = %q, want %q", sv.path, "/market-data/screeners/gainers-losers/list")
	}
	if sv.rankType != "DAY_1" {
		t.Errorf("rank_type = %q, want %q", sv.rankType, "DAY_1")
	}
	if sv.category != "US_STOCK" {
		t.Errorf("category = %q, want %q", sv.category, "US_STOCK")
	}
	if sv.sortBy != "CHANGE_RATIO" {
		t.Errorf("sort_by = %q, want %q", sv.sortBy, "CHANGE_RATIO")
	}
	if sv.direction != "ASC" {
		t.Errorf("direction = %q, want %q", sv.direction, "ASC")
	}
	if len(got) != 1 {
		t.Fatalf("got %d stocks, want 1", len(got))
	}
	if got[0].Symbol != "AAPL" || got[0].InstrumentID != "913256135" {
		t.Errorf("identity = %q/%q, want AAPL/913256135", got[0].Symbol, got[0].InstrumentID)
	}
	if got[0].ChangeRatio != "0.0142" || got[0].RelativeVolume10D != "11.91" {
		t.Errorf("decoded ratios = %+v", got[0])
	}
}

func TestGetDisplayTopActive(t *testing.T) {
	t.Parallel()
	srv, dsvc, sv := newDisplayScreenerTestServer(t)
	defer srv.Close()
	dc := newDataClientWithDisplay(t, dsvc)

	got, err := dc.GetDisplayTopActive(context.Background(), data.MostActiveQuery{
		Category:  data.StockCategoryUS,
		RankType:  data.MostActiveRankVolume,
		SortBy:    data.ScreenerSortVolume,
		Direction: data.SortDirectionDesc,
	})
	if err != nil {
		t.Fatalf("GetDisplayTopActive() error = %v", err)
	}
	if sv.path != "/market-data/screeners/top-actives/list" {
		t.Errorf("path = %q, want %q", sv.path, "/market-data/screeners/top-actives/list")
	}
	if sv.category != "US_STOCK" {
		t.Errorf("category = %q, want %q", sv.category, "US_STOCK")
	}
	if sv.rankType != "VOLUME" {
		t.Errorf("rank_type = %q, want %q", sv.rankType, "VOLUME")
	}
	if sv.sortBy != "VOLUME" {
		t.Errorf("sort_by = %q, want %q", sv.sortBy, "VOLUME")
	}
	if sv.direction != "DESC" {
		t.Errorf("direction = %q, want %q", sv.direction, "DESC")
	}
	if len(got) != 1 {
		t.Fatalf("got %d stocks, want 1", len(got))
	}
	if got[0].Symbol != "AAPL" {
		t.Errorf("Symbol = %q, want AAPL", got[0].Symbol)
	}
}

func TestGetDisplayGainersLosers_omitsOptionalParams(t *testing.T) {
	t.Parallel()
	srv, dsvc, sv := newDisplayScreenerTestServer(t)
	defer srv.Close()
	dc := newDataClientWithDisplay(t, dsvc)

	_, err := dc.GetDisplayGainersLosers(context.Background(), data.GainersLosersQuery{
		Category: data.StockCategoryUS,
	})
	if err != nil {
		t.Fatalf("GetDisplayGainersLosers() error = %v", err)
	}
	if sv.rankType != "" {
		t.Errorf("rank_type = %q, want empty", sv.rankType)
	}
	if sv.sortBy != "" {
		t.Errorf("sort_by = %q, want empty", sv.sortBy)
	}
	if sv.direction != "" {
		t.Errorf("direction = %q, want empty", sv.direction)
	}
}
