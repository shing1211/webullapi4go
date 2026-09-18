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

// Package data provides market data HTTP endpoints for the Webull OpenAPI.
// Crypto endpoints use category=CRYPTO and symbol=<ticker> as query parameters
// on /market-data/bars, /market-data/tick, etc. HK sandbox does not support
// the CRYPTO category (returns 417). US sandbox credentials are required to
// verify response schemas.
package data

import (
	"context"
	"net/url"
	"strconv"
)

const (
	pathCryptoBars     = "/market-data/bars"
	pathCryptoTick     = "/market-data/tick"
	pathCryptoDepth    = "/market-data/depth"
	pathCryptoSnapshot = "/market-data/snapshot"
)

type CryptoBarsQuery struct {
	Symbol   string
	Interval BarTimespan
	Count    int
}

type CryptoBar struct {
	Symbol    string            `json:"symbol"`
	Exchange  string            `json:"exchange"`
	Currency  string            `json:"currency"`
	Open      string            `json:"open"`
	High      string            `json:"high"`
	Low       string            `json:"low"`
	Close     string            `json:"close"`
	Volume    string            `json:"volume"`
	Turnover  string            `json:"turnover"`
	Timestamp string            `json:"timestamp"`
	Timespan  BarTimespan       `json:"timespan"`
	Extra     map[string]string `json:"-"`
}

type CryptoTickQuery struct {
	Symbol string
	Count  int
}

type CryptoTick struct {
	Symbol    string            `json:"symbol"`
	Exchange  string            `json:"exchange"`
	Currency  string            `json:"currency"`
	Price     string            `json:"price"`
	Volume    string            `json:"volume"`
	Turnover  string            `json:"turnover"`
	Timestamp string            `json:"timestamp"`
	Direction string            `json:"direction"`
	TradeID   string            `json:"trade_id"`
	Extra     map[string]string `json:"-"`
}

type CryptoDepthQuery struct {
	Symbol string
	Count  int
}

type CryptoDepthLevel struct {
	BidPrice string `json:"bid_price"`
	BidSize  string `json:"bid_size"`
	AskPrice string `json:"ask_price"`
	AskSize  string `json:"ask_size"`
}

type CryptoDepth struct {
	Symbol    string             `json:"symbol"`
	Exchange  string             `json:"exchange"`
	Currency  string             `json:"currency"`
	Timestamp string             `json:"timestamp"`
	Levels    []CryptoDepthLevel `json:"levels"`
	Extra     map[string]string  `json:"-"`
}

type CryptoSnapshotQuery struct {
	Symbol string
}

type CryptoSnapshot struct {
	Symbol    string            `json:"symbol"`
	Exchange  string            `json:"exchange"`
	Currency  string            `json:"currency"`
	LastPrice string            `json:"last_price"`
	Open      string            `json:"open"`
	High      string            `json:"high"`
	Low       string            `json:"low"`
	Close     string            `json:"close"`
	Volume    string            `json:"volume"`
	Turnover  string            `json:"turnover"`
	Bid       string            `json:"bid"`
	Ask       string            `json:"ask"`
	Timestamp string            `json:"timestamp"`
	Extra     map[string]string `json:"-"`
}

func (c *Client) GetCryptoBars(ctx context.Context, q CryptoBarsQuery) ([]CryptoBar, error) {
	query := url.Values{}
	query.Set("category", "CRYPTO")
	query.Set("symbol", q.Symbol)
	if q.Interval != "" {
		query.Set("timespan", string(q.Interval))
	}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}
	var out []CryptoBar
	if err := c.get(ctx, pathCryptoBars, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetCryptoTick(ctx context.Context, q CryptoTickQuery) ([]CryptoTick, error) {
	query := url.Values{}
	query.Set("category", "CRYPTO")
	query.Set("symbol", q.Symbol)
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}
	var out []CryptoTick
	if err := c.get(ctx, pathCryptoTick, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetCryptoDepth(ctx context.Context, q CryptoDepthQuery) (*CryptoDepth, error) {
	query := url.Values{}
	query.Set("category", "CRYPTO")
	query.Set("symbol", q.Symbol)
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}
	var out CryptoDepth
	if err := c.get(ctx, pathCryptoDepth, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetCryptoSnapshot(ctx context.Context, q CryptoSnapshotQuery) (*CryptoSnapshot, error) {
	query := url.Values{}
	query.Set("category", "CRYPTO")
	query.Set("symbol", q.Symbol)
	var out CryptoSnapshot
	if err := c.get(ctx, pathCryptoSnapshot, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
