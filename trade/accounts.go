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

package trade

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/shing1211/webullapi4go/pkg/errors"
)

// Trading account and asset endpoint paths.
const (
	// pathAccountsList lists the accounts available to the caller.
	pathAccountsList = "/trading/accounts/list"
	// pathBalanceGet returns the asset balance of one account.
	pathBalanceGet = "/trading/assets/balances/get"
	// pathPositionsList lists the open positions of one account.
	pathPositionsList = "/trading/assets/positions/list"
	// pathCashActivitiesList returns the cash activities of one account.
	//
	// Reference: https://developer.webull.com/apis/docs/reference/trade-cash-activity-by-type.md
	pathCashActivitiesList = "/trading/activities/cash-activities/list"
)

// ListAccounts returns the accounts available to the authenticated user.
//
// Reference: https://developer.webull.hk/apis/docs/reference/account-list.md
func (c *Client) ListAccounts(ctx context.Context) ([]Account, error) {
	var out []Account
	if err := c.get(ctx, pathAccountsList, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetBalance returns the multi-currency asset balance of accountID, including
// the per-currency breakdown.
//
// accountID is required; an empty or whitespace-only value is rejected with an
// [errs.CodeInvalidConfig] error before any network call.
//
// Reference: https://developer.webull.hk/apis/docs/reference/query-account-balance.md
func (c *Client) GetBalance(ctx context.Context, accountID string) (*AssetsBalance, error) {
	query, err := accountQuery(accountID)
	if err != nil {
		return nil, err
	}
	var out AssetsBalance
	if err := c.get(ctx, pathBalanceGet, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPositions returns the open positions of accountID. The result is empty
// when the account holds no positions.
//
// accountID is required; an empty or whitespace-only value is rejected with an
// [errs.CodeInvalidConfig] error before any network call.
//
// Reference: https://developer.webull.hk/apis/docs/reference/query-account-position.md
func (c *Client) GetPositions(ctx context.Context, accountID string) ([]Position, error) {
	query, err := accountQuery(accountID)
	if err != nil {
		return nil, err
	}
	var out []Position
	if err := c.get(ctx, pathPositionsList, query, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// accountQuery builds the account_id query shared by the account-scoped asset
// endpoints. It returns an [errs.CodeInvalidConfig] error when accountID is
// empty or whitespace-only, so callers never send an account-less request.
func accountQuery(accountID string) (url.Values, error) {
	if strings.TrimSpace(accountID) == "" {
		return nil, errs.New(errs.CodeInvalidConfig, "account id is required")
	}
	query := url.Values{}
	query.Set("account_id", accountID)
	return query, nil
}

// CashActivityType identifies the type of cash activity.
type CashActivityType string

// Cash activity types.
const (
	CashActivityTypeTrade    CashActivityType = "TRADE"
	CashActivityTypeDividend CashActivityType = "DIVIDEND"
	CashActivityTypeInterest CashActivityType = "INTEREST"
	CashActivityTypeTransfer CashActivityType = "TRANSFER"
)

// CashActivity is a single cash flow record for an account.
type CashActivity struct {
	ID              string           `json:"id"`
	AccountID       string           `json:"account_id"`
	AccountNumber   string           `json:"account_number"`
	ActivityType    CashActivityType `json:"activity_type"`
	ActivitySubType string           `json:"activity_sub_type"`
	Currency        string           `json:"currency"`
	Market          string           `json:"market"`
	Symbol          string           `json:"symbol"`
	TradeDate       string           `json:"trade_date"`
	NetAmount       string           `json:"net_amount"`
	BizTime         string           `json:"biz_time"`
}

// CashActivityPage is a single page of cash activities.
type CashActivityPage struct {
	Activities    []CashActivity `json:"data"`
	PaginationKey string         `json:"pagination_key"`
}

// CashActivityQuery parameterizes cash activity queries. AccountID is required.
type CashActivityQuery struct {
	AccountID     string
	ActivityType  CashActivityType
	PaginationKey string
}

func (q CashActivityQuery) values() (url.Values, error) {
	query, err := accountQuery(q.AccountID)
	if err != nil {
		return nil, err
	}
	if q.ActivityType != "" {
		query.Set("activity_type", string(q.ActivityType))
	}
	if q.PaginationKey != "" {
		query.Set("pagination_key", q.PaginationKey)
	}
	return query, nil
}

// GetCashActivities returns the first page of cash activities for the account.
//
// q.AccountID is required; an empty or whitespace-only value is rejected with
// an [errs.CodeInvalidConfig] error before any network call.
//
// Reference: https://developer.webull.com/apis/docs/reference/trade-cash-activity-by-type.md
func (c *Client) GetCashActivities(ctx context.Context, q CashActivityQuery) ([]CashActivity, error) {
	page, err := c.GetCashActivitiesPage(ctx, q)
	if err != nil {
		return nil, err
	}
	return page.Activities, nil
}

// GetCashActivitiesPage returns one page of cash activities. Use q.PaginationKey
// to resume from a previous page. The returned page preserves the cursor so
// callers can page explicitly.
//
// q.AccountID is required; an empty or whitespace-only value is rejected with
// an [errs.CodeInvalidConfig] error before any network call.
//
// Reference: https://developer.webull.com/apis/docs/reference/trade-cash-activity-by-type.md
func (c *Client) GetCashActivitiesPage(ctx context.Context, q CashActivityQuery) (*CashActivityPage, error) {
	query, err := q.values()
	if err != nil {
		return nil, err
	}
	var out CashActivityPage
	if err := c.get(ctx, pathCashActivitiesList, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAllCashActivities follows the pagination cursor to exhaustion and returns
// all cash activities matching q. The walk is bounded by [MaxOrderQueryPages];
// a server that keeps returning a non-empty pagination_key fails the call with
// an [errs.CodeAPI] error rather than looping forever.
//
// q.AccountID is required; an empty or whitespace-only value is rejected with
// an [errs.CodeInvalidConfig] error before any network call.
//
// Reference: https://developer.webull.com/apis/docs/reference/trade-cash-activity-by-type.md
func (c *Client) GetAllCashActivities(ctx context.Context, q CashActivityQuery) ([]CashActivity, error) {
	var (
		all []CashActivity
		key string
	)
	for page := 1; ; page++ {
		if page > MaxOrderQueryPages {
			return nil, errs.New(errs.CodeAPI,
				fmt.Sprintf("cash activity query exceeded the maximum of %d pages", MaxOrderQueryPages))
		}
		q.PaginationKey = key
		res, err := c.GetCashActivitiesPage(ctx, q)
		if err != nil {
			return nil, err
		}
		all = append(all, res.Activities...)
		if res.PaginationKey == "" {
			return all, nil
		}
		key = res.PaginationKey
	}
}
