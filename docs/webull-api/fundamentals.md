# Fundamentals and Fund Data

Company fundamentals, analyst data, financial statements and fund data under the Non-Display Solution.

[<- Webull API Reference](../webull-api.md)

## Company Profile

`GET /market-data/fundamentals/company-profiles/get`

> Retrieves company profile information.

| | |
|---|---|
| **SDK** | `data.GetCompanyProfile` |
| **Reference** | [get-company-profile.md](https://developer.webull.hk/apis/docs/reference/get-company-profile.md) |
| **Note** | US only in practice. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Possible values: US_STOCK. — one of: `US_STOCK` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Security symbol |
| `category` | string |  | Security type |
| `company_name` | string |  | Company name |
| `establish_date` | string |  | Date of incorporation |
| `exhibition_code` | string |  | The current market where the target is located |
| `profile` | string |  | Company profile |
| `employees` | string |  | Number of employees |
| `address` | string |  | Headquarters address |
| `ceo` | string |  | Company CEO |
| `industries` | array<string> |  | Company industries |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Analyst Target Price

`GET /market-data/fundamentals/analysis/target-prices/get`

> Retrieves analyst target price data.

| | |
|---|---|
| **SDK** | `data.GetAnalystTargetPrice` |
| **Reference** | [get-analyst-target-price.md](https://developer.webull.hk/apis/docs/reference/get-analyst-target-price.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Possible values: US_STOCK. — one of: `US_STOCK` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Security symbol |
| `category` | string |  | Security type |
| `mean` | string |  | Average target price |
| `low` | string |  | Lowest target price |
| `high` | string |  | Highest target price |
| `median` | string |  | Median price |
| `currency` | string |  | Currency |
| `effective_start_date` | string |  | Effective start date |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Analyst Rating

`GET /market-data/fundamentals/analysis/ratings/get`

> Retrieves analyst rating data.

| | |
|---|---|
| **SDK** | `data.GetAnalystRating` |
| **Reference** | [get-analyst-rating.md](https://developer.webull.hk/apis/docs/reference/get-analyst-rating.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Possible values: US_STOCK. — one of: `US_STOCK` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Security symbol |
| `category` | string |  | Security type |
| `number` | string |  | Total number of analysts |
| `under_perform` | string |  | Under perform count |
| `buy` | string |  | Buy count |
| `sell` | string |  | Sell count |
| `strong_buy` | string |  | Strong buy count |
| `hold` | string |  | Hold (neutral) count |
| `effective_start_date` | string |  | Effective start date |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Capital Flow

`GET /market-data/fundamentals/capital-flows/get`

> Retrieves stock capital flow data.

| | |
|---|---|
| **SDK** | `data.GetCapitalFlow` |
| **Reference** | [capital-flow.md](https://developer.webull.hk/apis/docs/reference/capital-flow.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. Only single symbol supported. |
| `category` | query | string | yes | Security type. Supports US_STOCK, HK_STOCK, CN_STOCK. — one of: `US_STOCK`, `HK_STOCK`, `CN_STOCK` |
| `count` | query | string |  | Number of records to return, range 1~5, default 5. Fetches the most recent N trading days' capital flow in reverse chronological order, response sorted in ascending order. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `date` | string |  | Date YYYYMMDD |
| `large_in` | string |  | Large inflow amount, in the base currency of the market (e.g., CNY, USD) |
| `large_out` | string |  | Large outflow amount |
| `medium_in` | string |  | Medium inflow amount |
| `medium_out` | string |  | Medium outflow amount |
| `small_in` | string |  | Small inflow amount |
| `small_out` | string |  | Small outflow amount |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Industry Comparison

`GET /market-data/fundamentals/industry-comparisons/get`

> Retrieves stock industry comparison data.

| | |
|---|---|
| **SDK** | `data.GetIndustryComparison` |
| **Reference** | [industry-comparison.md](https://developer.webull.hk/apis/docs/reference/industry-comparison.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. Only single symbol supported. |
| `category` | query | string | yes | Security type. Supports US_STOCK, HK_STOCK, CN_STOCK. — one of: `US_STOCK`, `HK_STOCK`, `CN_STOCK` |
| `sort_by` | query | string |  | Sort by financial metric, default EPS_TTM. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `fiscal_year` | integer |  | Fiscal year |
| `fiscal_period` | integer |  | Fiscal period (quarter) |
| `industry_name` | string |  | Industry name |
| `type` | string |  | Financial metric type |
| `data` | array<object> |  | Comparison data list |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Security symbol |
| `name` | string |  | Security name |
| `rank` | integer |  | Rank |
| `value` | string |  | Metric value |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Earnings Calendar

`GET /market-data/fundamentals/earnings-calendars/list`

> Retrieves stock earnings calendar data.

| | |
|---|---|
| **SDK** | `data.GetEarningsCalendar` |
| **Reference** | [earnings-calendar.md](https://developer.webull.hk/apis/docs/reference/earnings-calendar.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. Only single symbol supported. |
| `category` | query | string | yes | Security type. Supports US_STOCK, HK_STOCK, CN_STOCK. — one of: `US_STOCK`, `HK_STOCK`, `CN_STOCK` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `fiscal_year` | integer |  | Fiscal year |
| `fiscal_period` | integer |  | Fiscal period (quarter) |
| `currency` | string |  | Currency |
| `expected_publish_date` | string |  | Expected publish date, format YYYY-MM-DD |
| `eps_actual` | string |  | Actual EPS |
| `eps_est` | string |  | Estimated EPS |
| `rev_actual` | string |  | Actual revenue |
| `rev_est` | string |  | Estimated revenue |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Dividend Calendar

`GET /market-data/fundamentals/dividend-calendars/list`

> Retrieves stock dividend calendar data.

| | |
|---|---|
| **SDK** | `data.GetDividendCalendar` |
| **Reference** | [dividend-calendar.md](https://developer.webull.hk/apis/docs/reference/dividend-calendar.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. Only single symbol supported. |
| `category` | query | string | yes | Security type. Supports US_STOCK, HK_STOCK, CN_STOCK. — one of: `US_STOCK`, `HK_STOCK`, `CN_STOCK` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Symbol |
| `market` | string |  | Market, e.g. US, HK, CN |
| `currency` | string |  | Currency |
| `amount` | string |  | Dividend amount per share |
| `div_type` | string |  | Dividend type — one of: `CASH_DIVIDEND`, `STOCK_DIVIDEND`, `SCRIP_DIVIDEND`, `RETURN_OF_CAPITAL`, `INTEREST_PRINCIPAL_PAYMENT`, `DIVIDEND_REINVESTMENT`, `TRUST_INCOME_DISTRIBUTION` |
| `declare_date` | string |  | Declare date, format YYYY-MM-DD |
| `ex_div_date` | string |  | Ex-dividend date, format YYYY-MM-DD |
| `record_date` | string |  | Record date, format YYYY-MM-DD |
| `pay_date` | string |  | Pay date, format YYYY-MM-DD |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Filings

`GET /market-data/fundamentals/filings/list`

> Retrieves stock filings data.

| | |
|---|---|
| **SDK** | `data.GetFilings` |
| **Reference** | [filings.md](https://developer.webull.hk/apis/docs/reference/filings.md) |
| **Note** | US/SEC only. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. Only single symbol supported. |
| `category` | query | string | yes | Security type. Only US_STOCK is supported currently. — one of: `US_STOCK` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `symbol` | string |  | Security symbol |
| `category` | string |  | Security type |
| `filings` | array<object> |  | Filings list |

*Nested — `filings`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `title` | string |  | Filing title |
| `url` | string |  | Filing URL |
| `publish_date` | string |  | Publish date, format YYYY-MM-DD |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Income Statement

`GET /market-data/fundamentals/income-statements/get`

> Retrieves financial income statement data.

| | |
|---|---|
| **SDK** | `data.GetIncomeStatement` |
| **Reference** | [financial-income.md](https://developer.webull.hk/apis/docs/reference/financial-income.md) |
| **Note** | Returns `map[string]any` per period. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Currently only US_STOCK is supported. — one of: `US_STOCK` |
| `type` | query | string |  | Financial type: ANNUAL or QUARTERLY. |
| `count` | query | string |  | The number of each query, default value is 5, maximum value is 20. |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `fiscal_year` | integer |  | Fiscal year |
| `fiscal_period` | integer |  | Fiscal period (0=FY, 1=Q1, 2=Q2, 3=Q3, 4=Q4) |
| `end_date` | string |  | Report end date |
| `currency` | string |  | Currency |
| `publish_date` | string |  | Publish date |
| `total_revenue` | string |  | Total revenue |
| `revenue` | string |  | Revenue |
| `cost_of_revenue` | string |  | Total cost of revenue |
| `gross_profit` | string |  | Gross profit |
| `opex` | string |  | Operating expenses |
| `sga_exp` | string |  | Selling, general and administrative expenses |
| `rnd_exp` | string |  | Research and development expenses |
| `op_income` | string |  | Operating income |
| `other_net_income` | string |  | Other net income |
| `ebt` | string |  | Net income before tax |
| `income_tax` | string |  | Income tax |
| `eat` | string |  | Net income after tax |
| `ni_pre_extra` | string |  | Net income before extraordinary items |
| `extra_items` | string |  | Total extraordinary items |
| `net_income` | string |  | Net income |
| `ni_common_excl_extra` | string |  | Income available to common shareholders excluding extraordinary items |
| `ni_common_incl_extra` | string |  | Income available to common shareholders including extraordinary items |
| `diluted_ni` | string |  | Diluted net income |
| `diluted_avg_shares` | string |  | Diluted weighted average shares |
| `diluted_eps_excl_extra` | string |  | Diluted EPS excluding extraordinary items |
| `diluted_eps_incl_extra` | string |  | Diluted EPS including extraordinary items |
| `dps` | string |  | Dividends per share |
| `diluted_norm_eps` | string |  | Diluted normalized EPS |
| `op_profit` | string |  | Operating profit |
| `eat_alt` | string |  | Earnings after tax |
| `ebt_alt` | string |  | Earnings before tax |
| `unusual_expense_income` | string |  | Non-recurring expenses (income) |
| `inter_inc_expse_net_non_oper` | string |  | Net interest expense (income), non-operating |
| `gain_loss_on_sale_of_assets` | string |  | Gain (Loss) from Asset Sale |
| `minority_interest` | string |  | Minority shareholders' equity |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Balance Sheet

`GET /market-data/fundamentals/balance-sheets/get`

> Retrieves financial balance sheet data.

| | |
|---|---|
| **SDK** | `data.GetBalanceSheet` |
| **Reference** | [financial-balancesheet.md](https://developer.webull.hk/apis/docs/reference/financial-balancesheet.md) |
| **Note** | Returns `map[string]any` per period. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Currently only US_STOCK is supported. — one of: `US_STOCK` |
| `type` | query | string |  | Financial type: ANNUAL or QUARTERLY. |
| `count` | query | string |  | The number of each query, default value is 5, maximum value is 20. |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `fiscal_year` | integer |  | Fiscal year |
| `fiscal_period` | integer |  | Fiscal period (0=FY, 1=Q1, 2=Q2, 3=Q3, 4=Q4) |
| `end_date` | string |  | Report end date |
| `currency` | string |  | Currency |
| `publish_date` | string |  | Publish date |
| `total_assets` | string |  | Total assets |
| `total_cur_assets` | string |  | Total current assets |
| `cash_st_invest` | string |  | Cash and short-term investments |
| `cash` | string |  | Cash |
| `cash_equiv` | string |  | Cash equivalents |
| `st_invest` | string |  | Short-term investments |
| `total_recv_net` | string |  | Total receivables (net) |
| `ar_trade_net` | string |  | Trade receivables (net) |
| `total_inv` | string |  | Total inventory |
| `other_cur_assets` | string |  | Other current assets |
| `total_non_cur_assets` | string |  | Total non-current assets |
| `ppe_net` | string |  | Property, plant and equipment (net) |
| `ppe_gross` | string |  | Property, plant and equipment (gross) |
| `acc_depre` | string |  | Accumulated depreciation |
| `lt_invest` | string |  | Long-term investments |
| `other_lt_assets` | string |  | Other long-term assets |
| `total_liab` | string |  | Total liabilities |
| `total_cur_liab` | string |  | Total current liabilities |
| `ap` | string |  | Accounts payable |
| `notes_st_debt` | string |  | Short-term debt |
| `cur_lt_debt_lease` | string |  | Current portion of long-term debt |
| `other_cur_liab` | string |  | Other current liabilities |
| `total_non_cur_liab` | string |  | Total non-current liabilities |
| `total_lt_debt` | string |  | Total long-term debt |
| `lt_debt` | string |  | Long-term debt |
| `total_debt` | string |  | Total debt |
| `other_liab` | string |  | Other liabilities |
| `total_equity` | string |  | Total equity |
| `total_sh_equity` | string |  | Total shareholders' equity |
| `common_stock` | string |  | Common stock |
| `apic` | string |  | Additional paid-in capital |
| `retained_earnings` | string |  | Retained earnings |
| `other_equity` | string |  | Other equity |
| `total_liab_sh_equity` | string |  | Total liabilities and shareholders' equity |
| `common_shares_out` | string |  | Total common shares outstanding |
| `prepaid_expenses` | string |  | Advance payment for expenses |
| `accrued_expenses` | string |  | Accrued expenses |
| `goodwill_net` | string |  | Net goodwill value |
| `intangibles_net` | string |  | Net value of intangible assets |
| `note_rece_long_term` | string |  | Long-term receivable bills |
| `capital_lease_obligations` | string |  | Long-term debt in capital lease transactions |
| `minority_interest` | string |  | Minority shareholders' equity |
| `non_redeemable_preferred_stock` | string |  | Non-redeemable preferred stocks in total |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Cash Flow

`GET /market-data/fundamentals/cash-flows/get`

> Retrieves financial cash flow statement data.

| | |
|---|---|
| **SDK** | `data.GetCashFlow` |
| **Reference** | [financial-cashflow.md](https://developer.webull.hk/apis/docs/reference/financial-cashflow.md) |
| **Note** | Returns `map[string]any` per period. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Currently only US_STOCK is supported. — one of: `US_STOCK` |
| `type` | query | string |  | Financial type: ANNUAL or QUARTERLY. |
| `count` | query | string |  | The number of each query, default value is 5, maximum value is 20. |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `fiscal_year` | integer |  | Fiscal year |
| `fiscal_period` | integer |  | Fiscal period (0=FY, 1=Q1, 2=Q2, 3=Q3, 4=Q4) |
| `end_date` | string |  | Report end date |
| `currency` | string |  | Currency |
| `publish_date` | string |  | Publish date |
| `cfo` | string |  | Cash from operating activities |
| `net_income` | string |  | Net income |
| `dna` | string |  | Depreciation and amortization |
| `deferred_tax` | string |  | Deferred taxes |
| `non_cash_items` | string |  | Non-cash items |
| `wc_change` | string |  | Changes in working capital |
| `cfi` | string |  | Cash from investing activities |
| `capex` | string |  | Capital expenditures |
| `other_cfi_items` | string |  | Other investing cash flow items |
| `cff` | string |  | Cash from financing activities |
| `cff_items` | string |  | Financing cash flow items |
| `net_stock_iss_ret` | string |  | Net issuance or repurchase of stock |
| `net_debt_iss_ret` | string |  | Net issuance or repayment of debt |
| `fx_effects` | string |  | Foreign exchange effects on cash |
| `net_change_cash` | string |  | Net change in cash |
| `interest_paid` | string |  | Cash interest paid |
| `taxes_paid` | string |  | Cash taxes paid |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Financial Indicators

`GET /market-data/fundamentals/indicators/get`

> Retrieves financial indicators data.

| | |
|---|---|
| **SDK** | `data.GetFinancialIndicators` |
| **Reference** | [financial-indicators.md](https://developer.webull.hk/apis/docs/reference/financial-indicators.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Currently only US_STOCK is supported. — one of: `US_STOCK` |
| `type` | query | string |  | Financial type: ANNUAL or QUARTERLY. |
| `count` | query | string |  | The number of each query, default value is 5, maximum value is 20. |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `currency` | string |  | Currency |
| `values` | object |  | Financial report factors. Keys: roa (Return on Total Assets), roe (Return on Net Assets), diluted_eps_incl_extra (Earnings per share), net_margin (Net profit margin), debt_to_assets (Debt ratio), naps (Per-share net asset value), ocf_ps (Per-share cash flow), cap_surplus_ps (Per-share reserve fund) |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Financial Alert

`GET /market-data/fundamentals/financial-alerts/get`

> Retrieves financial alert data.

| | |
|---|---|
| **SDK** | `data.GetFinancialAlert` |
| **Reference** | [financial-alert.md](https://developer.webull.hk/apis/docs/reference/financial-alert.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Currently only US_STOCK is supported. — one of: `US_STOCK` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `start_date` | string |  | Earliest possible date for financial report release |
| `end_date` | string |  | Latest possible date for financial report release |
| `fiscal_year` | integer |  | Fiscal year |
| `fiscal_period` | integer |  | Fiscal quarter (1=Q1, 2=Q2, 3=Q3, 4=Q4, 5=Pre-release) |
| `currency` | string |  | Currency |
| `eps_est` | string |  | Current period projected earnings per share |
| `eps_ly` | string |  | Earnings per share for the same period last year |
| `rev_est` | string |  | Current period projected revenue |
| `rev_ly` | string |  | Revenue for the corresponding period of the previous year |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Forecast EPS

`GET /market-data/fundamentals/forecast-eps/get`

> Retrieves stock forecast EPS data.

| | |
|---|---|
| **SDK** | `data.GetForecastEPS` |
| **Reference** | [forecast-eps.md](https://developer.webull.hk/apis/docs/reference/forecast-eps.md) |
| **Note** | Most recent 5 quarters. |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. Only single symbol supported. |
| `category` | query | string | yes | Security type. Supports US_STOCK, HK_STOCK, CN_STOCK. — one of: `US_STOCK`, `HK_STOCK`, `CN_STOCK` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `fiscal_year` | integer |  | Fiscal year |
| `fiscal_period` | integer |  | Fiscal period (quarter) |
| `actual` | string |  | Actual EPS |
| `est` | string |  | Estimated EPS |
| `reported` | boolean |  | Whether the earnings have been reported |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Fund Brief

`GET /market-data/fundamentals/fund-brief/get`

> Retrieves fund brief information.

| | |
|---|---|
| **SDK** | `data.GetFundInfo` |
| **Reference** | [fund-brief.md](https://developer.webull.hk/apis/docs/reference/fund-brief.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Currently only US_STOCK is supported. — one of: `US_STOCK` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string |  | Fund Name |
| `launch_date` | string |  | Inception Date |
| `benchmark` | string |  | Benchmark |
| `investment_objective` | string |  | Investment Objective |
| `aum` | string |  | Assets |
| `issuer` | string |  | Advisor Company |
| `custodian` | string |  | Custodian |
| `managers` | array<object> |  | Details of the Fund Manager |

*Nested — `managers`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string |  | Name |
| `start_date` | string |  | Start date of employment |
| `end_date` | string |  | End date of employment |
| `is_incumbent` | integer |  | Is incumbent (1=yes) |
| `tenure_return` | string |  | Term-based returns |
| `title` | string |  | Fund Manager Position |
| `tenure_years` | string |  | Period of Service |
| `tenure_days` | integer |  | Term days |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Fund Performance

`GET /market-data/fundamentals/fund-performances/get`

> Retrieves fund performance data.

| | |
|---|---|
| **SDK** | `data.GetFundPerformance` |
| **Reference** | [fund-performance.md](https://developer.webull.hk/apis/docs/reference/fund-performance.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Currently only US_STOCK is supported. — one of: `US_STOCK` |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `currency` | string |  | Currency |
| `end_date` | string |  | End Date |
| `return_1m` | string |  | One Month return |
| `return_3m` | string |  | Three Month return |
| `return_6m` | string |  | Six Month return |
| `return_1y` | string |  | One Year return |
| `return_3y` | string |  | Three Year return |
| `return_5y` | string |  | Five Year return |
| `return_10y` | string |  | Ten Year return |
| `return_si` | string |  | Since its establishment |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Fund Net Value

`GET /market-data/fundamentals/fund-net-values/get`

> Retrieves fund net value data.

| | |
|---|---|
| **SDK** | `data.GetFundNav` |
| **Reference** | [fund-net-value.md](https://developer.webull.hk/apis/docs/reference/fund-net-value.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Currently only US_STOCK is supported. — one of: `US_STOCK` |
| `last_date` | query | string |  | Last Query Date. |
| `count` | query | string |  | The number of each query, default value is 5, maximum value is 20. |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `date` | string |  | Net Date |
| `currency` | string |  | Currency |
| `net_value` | string |  | Net Value |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Fund Holdings

`GET /market-data/fundamentals/fund-holdings/get`

> Retrieves fund holdings data.

| | |
|---|---|
| **SDK** | `data.GetFundHoldings` |
| **Reference** | [fund-holdings.md](https://developer.webull.hk/apis/docs/reference/fund-holdings.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Currently only US_STOCK is supported. — one of: `US_STOCK` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `target_symbol` | string |  | Target Symbol |
| `stock_name` | string |  | Security Name |
| `share_held_pct` | string |  | Share Held Pct |
| `share_held_chg_pct` | string |  | Share Held Chg Pct |
| `maturity_date` | string |  | Maturity date of securities held |
| `update_time` | string |  | Update Time |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Fund Dividends

`GET /market-data/fundamentals/fund-dividends/get`

> • Function description: Get fund dividend history. • Frequency limit: Rate limit 60 requests every 60 seconds

| | |
|---|---|
| **SDK** | `data.GetFundDividends` |
| **Reference** | [fund-dividends.md](https://developer.webull.hk/apis/docs/reference/fund-dividends.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Currently only US_STOCK is supported. — one of: `US_STOCK` |
| `pagination_key` | query | string |  | Pagination key from previous response for next page |

**Response 200**

| Field | Type | Required | Description |
|---|---|---|---|
| `data` | array<object> |  | List of fund dividend records |
| `pagination_key` | string |  | Pagination key for next page. If absent, indicates this is the last page. |

*Nested — `data`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `share_date` | string |  | Ex-Date |
| `publish_date` | string |  | Declaration Date |
| `pay_date` | string |  | Pay Date |
| `record_date` | string |  | Record Date |
| `dps` | string |  | Dividend Per Share |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Fund Rating

`GET /market-data/fundamentals/fund-ratings/get`

> Retrieves fund rating data.

| | |
|---|---|
| **SDK** | `data.GetFundRating` |
| **Reference** | [fund-rating.md](https://developer.webull.hk/apis/docs/reference/fund-rating.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Currently only US_STOCK is supported. — one of: `US_STOCK` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `rating_date` | string |  | Rating Date |
| `rating_agency` | string |  | Name of the rating agency |
| `rating_cycle` | string |  | Rating cycle. 0:Since establishment, 3:3 years, 5:5 years, 10:10 years |
| `rating_results` | integer |  | Overall rating, 1:Low, 2:Below Average, 3:Average, 4:Above Average, 5:High |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Fund Splits

`GET /market-data/fundamentals/fund-splits/get`

> Retrieves fund splits data.

| | |
|---|---|
| **SDK** | `data.GetFundSplits` |
| **Reference** | [fund-splits.md](https://developer.webull.hk/apis/docs/reference/fund-splits.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Currently only US_STOCK is supported. — one of: `US_STOCK` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `split_date` | string |  | Split date |
| `split_type` | string |  | Split type, e.g. MERGE, SPLIT |
| `split_ratio` | string |  | Split ratio |
| `from` | number |  | Split From |
| `to` | number |  | Split To |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Fund Files

`GET /market-data/fundamentals/fund-files/get`

> Retrieves fund files data.

| | |
|---|---|
| **SDK** | `data.GetFundFiles` |
| **Reference** | [fund-files.md](https://developer.webull.hk/apis/docs/reference/fund-files.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Currently only US_STOCK is supported. — one of: `US_STOCK` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `publish_date` | string |  | Publish Date |
| `url` | string |  | Url |
| `type` | integer |  | Fund document type, e.g, 1:Prospectus, 4:Annual Report, 5:Semi-Annual Report, 14:Quarterly Report, 17:Prospectus Summary, 37:Announcement, 51:Rulebook/Statutes, 52:Factsheet, 57:Rulebook Summary, 58:Custodian Agreement, 74:Investor Information Document, 76:Key Fact Statement, 77:Product Highlight Sheet, 112:Share Sale Announcement |
| `file_name` | string |  | File Name |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

## Fund Allocation

`GET /market-data/fundamentals/fund-allocations/get`

> Retrieves fund allocation data.

| | |
|---|---|
| **SDK** | `data.GetFundAllocation` |
| **Reference** | [fund-allocation.md](https://developer.webull.hk/apis/docs/reference/fund-allocation.md) |

**Request — parameters**

| Name | In | Type | Required | Description |
|---|---|---|---|---|
| `symbol` | query | string | yes | Security symbol. |
| `category` | query | string | yes | Security type. Currently only US_STOCK is supported. — one of: `US_STOCK` |

**Response 200**

Array of objects:

| Field | Type | Required | Description |
|---|---|---|---|
| `date` | string |  | End Date |
| `aum` | string |  | Total Assets |
| `cash` | object |  | Fund Asset |
| `bond` | object |  | Fund Asset |
| `stock` | object |  | Fund Asset |
| `preferred` | object |  | Fund Asset |
| `convertible` | object |  | Fund Asset |
| `other` | object |  | Fund Asset |

*Nested — `cash`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `value` | string |  | Assets |
| `ratio` | string |  | Ratio |

*Nested — `bond`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `value` | string |  | Assets |
| `ratio` | string |  | Ratio |

*Nested — `stock`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `value` | string |  | Assets |
| `ratio` | string |  | Ratio |

*Nested — `preferred`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `value` | string |  | Assets |
| `ratio` | string |  | Ratio |

*Nested — `convertible`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `value` | string |  | Assets |
| `ratio` | string |  | Ratio |

*Nested — `other`:*

| Field | Type | Required | Description |
|---|---|---|---|
| `value` | string |  | Assets |
| `ratio` | string |  | Ratio |

**Errors** — `401` unauthorized, `417` business error, `500` server error. See [Errors](../errors.md).

