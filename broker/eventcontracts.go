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
	"net/url"
)

const (
	pathECCategories  = "/broker/event-contracts/categories/list"
	pathECSeries      = "/broker/event-contracts/series/get"
	pathECEvents      = "/broker/event-contracts/events/get"
	pathECInstruments = "/broker/event-contracts/instruments/get"
)

type EventContractCategory struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type EventContractSeries struct {
	SeriesID   string `json:"series_id"`
	CategoryID string `json:"category_id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
}

type EventContractEvent struct {
	EventID   string `json:"event_id"`
	SeriesID  string `json:"series_id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	EventDate string `json:"event_date"`
}

type EventContractInstrument struct {
	Symbol         string `json:"symbol"`
	EventID        string `json:"event_id"`
	SeriesID       string `json:"series_id"`
	Status         string `json:"status"`
	StrikePrice    string `json:"strike_price"`
	ExpirationDate string `json:"expiration_date"`
}

func (c *Client) GetEventContractCategories(ctx context.Context) ([]EventContractCategory, error) {
	var out []EventContractCategory
	if err := c.get(ctx, pathECCategories, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetEventContractSeries(ctx context.Context, categoryID string) ([]EventContractSeries, error) {
	q := url.Values{}
	q.Set("category_id", categoryID)
	var out []EventContractSeries
	if err := c.get(ctx, pathECSeries, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetEventContractEvents(ctx context.Context, seriesID string) ([]EventContractEvent, error) {
	q := url.Values{}
	q.Set("series_id", seriesID)
	var out []EventContractEvent
	if err := c.get(ctx, pathECEvents, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetEventContractInstruments(ctx context.Context, eventID string) ([]EventContractInstrument, error) {
	q := url.Values{}
	q.Set("event_id", eventID)
	var out []EventContractInstrument
	if err := c.get(ctx, pathECInstruments, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}
