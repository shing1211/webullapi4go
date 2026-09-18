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
	"net/url"
	"strings"

	"github.com/shing1211/webullapi4go/internal/errs"
)

// Trading account and asset endpoint paths.
const (
	// pathAccountsList lists the accounts available to the caller.
	pathAccountsList = "/trading/accounts/list"
	// pathBalanceGet returns the asset balance of one account.
	pathBalanceGet = "/trading/assets/balances/get"
	// pathPositionsList lists the open positions of one account.
	pathPositionsList = "/trading/assets/positions/list"
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
