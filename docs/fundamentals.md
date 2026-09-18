# Fundamentals

Fundamental data provides company financial and ownership information: earnings
and dividend calendars, capital flows, industry comparisons, SEC filings, and
full financial statements.

All fundamentals endpoints live in the `data` package alongside the rest of the
Market Data HTTP API and share the same authentication and transport.

## Coverage

| Group | Methods |
|-------|---------|
| Company and analyst | `GetCompanyProfile`, `GetAnalystTargetPrice`, `GetAnalystRating` |
| Capital and industry | `GetCapitalFlow`, `GetIndustryComparison` |
| Calendars | `GetEarningsCalendar`, `GetDividendCalendar` |
| Filings | `GetFilings` |
| Financial statements | `GetIncomeStatement`, `GetBalanceSheet`, `GetCashFlow` |
| Metrics and forecasts | `GetFinancialIndicators`, `GetFinancialAlert`, `GetForecastEPS` |

## Capital Flow

Capital flow data shows net buying and selling activity broken down by trade
size (large, medium, small) over a trailing window.

```go
flows, err := market.GetCapitalFlow(ctx, "AAPL", data.StockCategoryUS, 5)
if err != nil {
	return err
}
for _, f := range flows {
	fmt.Println(f.Date, f.LargeIn, f.LargeOut)
}
```

## Industry Comparison

Relative performance metrics within an industry. The `compType` parameter selects
the comparison type (e.g. `EPS_TTM`, `ROE`, `P_E_RATIO`).

```go
comp, err := market.GetIndustryComparison(ctx, "AAPL", data.StockCategoryUS, "EPS_TTM")
if err != nil {
	return err
}
fmt.Println(comp.IndustryName)
for _, d := range comp.Data {
	fmt.Println(d.Rank, d.Symbol, d.Value)
}
```

## Earnings Calendar

Upcoming and historical earnings-release dates with EPS and revenue actuals vs.
estimates.

```go
cal, err := market.GetEarningsCalendar(ctx, "AAPL", data.StockCategoryUS)
if err != nil {
	return err
}
for _, e := range cal {
	fmt.Printf("%s Q%d: eps_actual=%s eps_est=%s\n",
		e.FiscalYear, e.FiscalPeriod, e.EPSActual, e.EPSEst)
}
```

## Dividend Calendar

Dividend and split events with declare, ex-div, record, and pay dates.

```go
divs, err := market.GetDividendCalendar(ctx, "AAPL", data.StockCategoryUS)
if err != nil {
	return err
}
for _, d := range divs {
	fmt.Printf("%s %s: amount=%s pay_date=%s\n",
		d.Symbol, d.DivType, d.Amount, d.PayDate)
}
```

## SEC Filings

SEC regulatory filings (8-K, 10-K, 10-Q, etc.) with titles, URLs, and publish
dates.

```go
filings, err := market.GetFilings(ctx, "AAPL", data.StockCategoryUS)
if err != nil {
	return err
}
for _, f := range filings.Filings {
	fmt.Println(f.Title, f.URL, f.PublishDate)
}
```

## Financial Statements

Multi-period financial statement data. Each method returns a slice of period
entries keyed by fiscal year and period.

### Income Statement

```go
income, err := market.GetIncomeStatement(ctx, "AAPL", data.StockCategoryUS)
if err != nil {
	return err
}
for _, p := range income {
	fmt.Printf("%s %s: revenue=%s net_income=%s\n",
		p["fiscal_year"], p["fiscal_period"],
		p["revenue"], p["net_income"])
}
```

### Balance Sheet

```go
balance, err := market.GetBalanceSheet(ctx, "AAPL", data.StockCategoryUS)
if err != nil {
	return err
}
for _, p := range balance {
	fmt.Printf("%s %s: assets=%s liabilities=%s\n",
		p["fiscal_year"], p["fiscal_period"],
		p["total_assets"], p["total_liabilities"])
}
```

### Cash Flow

```go
cf, err := market.GetCashFlow(ctx, "AAPL", data.StockCategoryUS)
if err != nil {
	return err
}
for _, p := range cf {
	fmt.Printf("%s %s: operating=%s free=%s\n",
		p["fiscal_year"], p["fiscal_period"],
		p["operating_cash_flow"], p["free_cash_flow"])
}
```

## Financial Indicators

Key ratios and metrics (ROA, ROE, EPS, net margin, debt ratio).

```go
ind, err := market.GetFinancialIndicators(ctx, "AAPL", data.StockCategoryUS)
if err != nil {
	return err
}
fmt.Printf("ROE=%s ROA=%s EPS=%s net_margin=%s\n",
	ind.ROE, ind.ROA, ind.EPS, ind.NetMargin)
```

## Financial Alert

Upcoming earnings-release alert with expected report date and estimated vs.
last-year EPS.

```go
alert, err := market.GetFinancialAlert(ctx, "AAPL", data.StockCategoryUS)
if err != nil {
	return err
}
fmt.Printf("expected %s: est_eps=%s last_year=%s\n",
	alert.ExpectedReportDate, alert.EstimatedEPS, alert.LastYearEPS)
```

## Forecast EPS

Analyst consensus EPS forecasts for the next five quarters.

```go
forecast, err := market.GetForecastEPS(ctx, "AAPL", data.StockCategoryUS)
if err != nil {
	return err
}
for _, f := range forecast {
	fmt.Printf("%d Q%d: actual=%s est=%s reported=%t\n",
		f.FiscalYear, f.FiscalPeriod,
		f.Actual, f.Est, f.Reported)
}
```

## Related

- [Market Data](market-data.md) — snapshot, bars, quotes, streaming.
- [Getting Started](getting-started.md) — install and credentials.
- [Authentication](authentication.md) — signing and tokens.
