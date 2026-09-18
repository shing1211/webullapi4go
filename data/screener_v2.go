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
	"net/http"
)

// pathScreenerV2 is the Screener v2 (next-gen) endpoint.
//
// Path and schema unconfirmed; live probe required to verify.
const pathScreenerV2 = "/wlas/screener/ng/query"

// ScreenerV2Operator is the comparison operator for a filter condition.
type ScreenerV2Operator string

// Screener v2 filter operators.
const (
	ScreenerV2OpEQ      ScreenerV2Operator = "eq"
	ScreenerV2OpNE      ScreenerV2Operator = "ne"
	ScreenerV2OpGT      ScreenerV2Operator = "gt"
	ScreenerV2OpGTE     ScreenerV2Operator = "gte"
	ScreenerV2OpLT      ScreenerV2Operator = "lt"
	ScreenerV2OpLTE     ScreenerV2Operator = "lte"
	ScreenerV2OpIn      ScreenerV2Operator = "in"
	ScreenerV2OpBetween ScreenerV2Operator = "between"
)

// ScreenerV2Logic is the logic connector between filter conditions.
type ScreenerV2Logic string

// Screener v2 logic connectors.
const (
	ScreenerV2LogicAnd ScreenerV2Logic = "AND"
	ScreenerV2LogicOr  ScreenerV2Logic = "OR"
)

// ScreenerV2Condition is a single filter condition in a Screener v2 query.
type ScreenerV2Condition struct {
	Field    string             `json:"field"`
	Operator ScreenerV2Operator `json:"operator"`
	Value    any                `json:"value"`
}

// ScreenerV2Filter is a filter group that can contain nested conditions and sub-filters.
type ScreenerV2Filter struct {
	Logic      ScreenerV2Logic       `json:"logic,omitempty"`
	Conditions []ScreenerV2Condition `json:"conditions,omitempty"`
	Filters    []ScreenerV2Filter    `json:"filters,omitempty"`
}

// ScreenerV2Sort specifies the sort field and direction for a Screener v2 query.
type ScreenerV2Sort struct {
	Field     string        `json:"field"`
	Direction SortDirection `json:"direction"`
}

// ScreenerV2Query is the request body sent to the Screener v2 POST endpoint.
type ScreenerV2Query struct {
	Filter        *ScreenerV2Filter `json:"filter,omitempty"`
	Sort          []ScreenerV2Sort  `json:"sort,omitempty"`
	PageSize      int               `json:"page_size,omitempty"`
	PaginationKey string            `json:"pagination_key,omitempty"`
}

// screenerV2Resp is the response envelope for the Screener v2 endpoint.
type screenerV2Resp struct {
	Data          []ScreenerStock `json:"data"`
	PaginationKey string          `json:"pagination_key"`
}

// ScreenerV2Result is the result of [Client.GetScreenerV2].
type ScreenerV2Result struct {
	Stocks        []ScreenerStock
	PaginationKey string
	Extra         map[string]string `json:"-"`
}

// GetScreenerV2 retrieves stocks using the Screener v2 (next-gen) endpoint.
// The endpoint accepts filter conditions, sort, and pagination in the request
// body and returns a paginated {data, pagination_key} envelope.
//
// Path and schema are unconfirmed; live probe required to verify.
func (c *Client) GetScreenerV2(ctx context.Context, q ScreenerV2Query) (*ScreenerV2Result, error) {
	var resp screenerV2Resp
	if err := c.do(ctx, http.MethodPost, pathScreenerV2, nil, q, &resp); err != nil {
		return nil, err
	}
	return &ScreenerV2Result{
		Stocks:        resp.Data,
		PaginationKey: resp.PaginationKey,
		Extra:         make(map[string]string),
	}, nil
}
