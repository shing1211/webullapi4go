# Todos — v0.4 (Market Data Completeness)

## T17.0: Recon Spike
- [ ] Read API reference for all candidate endpoints
- [ ] Probe sandbox for each candidate endpoint
- [ ] Classify each as (a) usable, (b) 403/417, (c) unknown
- [ ] Update plan.md findings table

## T17.1: Signer Reconciliation
- [ ] Add DigestCase to SignParams in internal/auth
- [ ] Update events to use internal/auth with lowercase digest
- [ ] Delete events/sign.go local duplicate
- [ ] Verify tests pass

## T17.2: Fundamentals
- [ ] Capital Flow
- [ ] Industry Comparison
- [ ] Earnings Calendar
- [ ] Dividend Calendar
- [ ] SEC Filings

## T17.3: Financial Statements
- [ ] Income Statements
- [ ] Balance Sheets
- [ ] Cash Flows
- [ ] Indicators
- [ ] Financial Alerts
- [ ] Forecast EPS

## T17.4: Tests + Docs
- [ ] Offline fixture tests
- [ ] docs/fundamentals.md
- [ ] examples/data-fundamentals/main.go
- [ ] README/CHANGELOG update
- [ ] mkdocs build --strict

## T17.5: Release v0.4.0
- [ ] Commit and tag on both remotes
