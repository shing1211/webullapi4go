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

package data

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const (
	pathCryptoBars     = "/market-data/crypto/%s/bars"
	pathCryptoTick     = "/market-data/crypto/%s/tick"
	pathCryptoDepth    = "/market-data/crypto/%s/depth"
	pathCryptoSnapshot = "/market-data/crypto/%s/snapshot"
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
	path := fmt.Sprintf(pathCryptoBars, url.PathEscape(q.Symbol))
	query := url.Values{}
	if q.Interval != "" {
		query.Set("timespan", string(q.Interval))
	}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}
	var out []CryptoBar
	if err := c.get(ctx, path, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetCryptoTick(ctx context.Context, q CryptoTickQuery) ([]CryptoTick, error) {
	path := fmt.Sprintf(pathCryptoTick, url.PathEscape(q.Symbol))
	query := url.Values{}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}
	var out []CryptoTick
	if err := c.get(ctx, path, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetCryptoDepth(ctx context.Context, q CryptoDepthQuery) (*CryptoDepth, error) {
	path := fmt.Sprintf(pathCryptoDepth, url.PathEscape(q.Symbol))
	query := url.Values{}
	if q.Count > 0 {
		query.Set("count", strconv.Itoa(q.Count))
	}
	var out CryptoDepth
	if err := c.get(ctx, path, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetCryptoSnapshot(ctx context.Context, q CryptoSnapshotQuery) (*CryptoSnapshot, error) {
	path := fmt.Sprintf(pathCryptoSnapshot, url.PathEscape(q.Symbol))
	var out CryptoSnapshot
	if err := c.get(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
