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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shing1211/webullapi4go/data"
	"github.com/shing1211/webullapi4go/pkg/domain/money"
)

func eventServer(t *testing.T, wantPath, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Path; got != wantPath {
			t.Errorf("path = %q, want %q", got, wantPath)
		}
		_, _ = w.Write([]byte(body))
	}))
}

func TestGetEventContractTags(t *testing.T) {
	t.Parallel()
	srv := eventServer(t, "/market-data/instruments/event-contracts/categories/tags/list",
		`[{"tags":["nba"],"category_id":1,"category_name":"Sports","category_code":"SPORTS"}]`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetEventContractTags(context.Background())
	if err != nil {
		t.Fatalf("GetEventContractTags() error = %v", err)
	}
	if len(got) != 1 || got[0].CategoryCode != "SPORTS" {
		t.Fatalf("tags = %+v", got)
	}
}

func TestGetEventContractEventsList(t *testing.T) {
	t.Parallel()
	srv := eventServer(t, "/market-data/instruments/event-contracts/events/list",
		`{"data":[{"event_symbol":"E1"}],"pagination_key":"p2"}`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetEventContractEventsList(context.Background(), data.EventListQuery{SeriesSymbol: "S1"})
	if err != nil {
		t.Fatalf("GetEventContractEventsList() error = %v", err)
	}
	if len(got.Data) != 1 || got.PaginationKey != "p2" {
		t.Fatalf("events = %+v", got)
	}
}

func TestGetEventContractMilestones(t *testing.T) {
	t.Parallel()
	srv := eventServer(t, "/market-data/instruments/event-contracts/milestones/list",
		`{"data":[{"milestone_id":"m1"}],"pagination_key":""}`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetEventContractMilestones(context.Background(), data.MilestoneQuery{Competition: "NBA"})
	if err != nil {
		t.Fatalf("GetEventContractMilestones() error = %v", err)
	}
	if len(got.Data) != 1 {
		t.Fatalf("milestones = %+v", got)
	}
}

func TestGetEventContractSeriesList(t *testing.T) {
	t.Parallel()
	srv := eventServer(t, "/market-data/instruments/event-contracts/series/list",
		`{"data":[{"series_symbol":"S1"}],"pagination_key":""}`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetEventContractSeriesList(context.Background(), data.EventSeriesListQuery{Category: "SPORTS"})
	if err != nil {
		t.Fatalf("GetEventContractSeriesList() error = %v", err)
	}
	if len(got.Data) != 1 {
		t.Fatalf("series = %+v", got)
	}
}

func TestGetEventContractSportsFilters(t *testing.T) {
	t.Parallel()
	srv := eventServer(t, "/market-data/instruments/event-contracts/sports-filters/list",
		`[{"tag":"NBA","competitions":["x"],"scopes":["y"]}]`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetEventContractSportsFilters(context.Background(), "")
	if err != nil {
		t.Fatalf("GetEventContractSportsFilters() error = %v", err)
	}
	if len(got) != 1 || got[0].Tag != "NBA" {
		t.Fatalf("filters = %+v", got)
	}
}

func TestGetEventGameStats(t *testing.T) {
	t.Parallel()
	srv := eventServer(t, "/market-data/event-contracts/game-stats/get",
		`{"milestone_id":"m1","periods":[{"q":1}]}`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetEventGameStats(context.Background(), "m1", "US_EVENT")
	if err != nil {
		t.Fatalf("GetEventGameStats() error = %v", err)
	}
	if got.MilestoneID != "m1" || len(got.Periods) != 1 {
		t.Fatalf("stats = %+v", got)
	}
}

func TestGetEventLiveData(t *testing.T) {
	t.Parallel()
	srv := eventServer(t, "/market-data/event-contracts/live-data/get",
		`{"type":"x","milestone_id":"m1","status":"LIVE","winner":"","last_play":{},"last_updated_ts":1,"details":{}}`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetEventLiveData(context.Background(), "m1", "US_EVENT")
	if err != nil {
		t.Fatalf("GetEventLiveData() error = %v", err)
	}
	if got.Status != "LIVE" {
		t.Fatalf("live = %+v", got)
	}
}

func TestGetEventMarketBars(t *testing.T) {
	t.Parallel()
	srv := eventServer(t, "/market-data/event-contracts/markets/bars/list",
		`[{"end_period_time":"t","volume":"1","open":"1","high":"2","low":"0","close":"1.5"}]`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetEventMarketBars(context.Background(), data.EventMarketBarsQuery{
		Symbols: []string{"S1"}, Timespan: "M1",
	})
	if err != nil {
		t.Fatalf("GetEventMarketBars() error = %v", err)
	}
	if len(got) != 1 || got[0].Close.Cmp(money.Must(money.NewFromString("1.5"))) != 0 {
		t.Fatalf("bars = %+v", got)
	}
}

func TestGetEventMarketBarsByEvent(t *testing.T) {
	t.Parallel()
	srv := eventServer(t, "/market-data/event-contracts/markets/bars/list-by-event",
		`[{"end_period_time":"t","volume":"1","open":"1","high":"2","low":"0","close":"1.5"}]`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetEventMarketBarsByEvent(context.Background(), data.EventMarketBarsByEventQuery{
		EventSymbol: "E1", Timespan: "M1",
	})
	if err != nil {
		t.Fatalf("GetEventMarketBarsByEvent() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("bars = %+v", got)
	}
}

func TestGetEventMarketDepth(t *testing.T) {
	t.Parallel()
	srv := eventServer(t, "/market-data/event-contracts/markets/depths/list",
		`{"symbol":"S1","instrument_id":"1","yes_asks":[],"yes_bids":[],"no_asks":[],"no_bids":[]}`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetEventMarketDepth(context.Background(), "S1", "US_EVENT", 0)
	if err != nil {
		t.Fatalf("GetEventMarketDepth() error = %v", err)
	}
	if got.Symbol != "S1" {
		t.Fatalf("depth = %+v", got)
	}
}

func TestGetEventMarketSnapshot(t *testing.T) {
	t.Parallel()
	srv := eventServer(t, "/market-data/event-contracts/markets/snapshots/list",
		`{"symbol":"S1","instrument_id":"1","event_symbol":"E1","yes_sub_title":"Y","no_sub_title":"N",`+
			`"status":"OPEN","yes_bid":"1","yes_ask":"2","no_bid":"9","no_ask":"8","price":"1.5",`+
			`"volume":"10","open_interest":"5","last_trade_time":"t"}`)
	defer srv.Close()
	c := newTestClient(t, srv.URL)
	got, err := c.GetEventMarketSnapshot(context.Background(), "S1", "US_EVENT")
	if err != nil {
		t.Fatalf("GetEventMarketSnapshot() error = %v", err)
	}
	if got.Price.Cmp(money.Must(money.NewFromString("1.5"))) != 0 || got.YesBid.Cmp(money.Must(money.NewFromString("1"))) != 0 {
		t.Fatalf("snapshot = %+v", got)
	}
}
