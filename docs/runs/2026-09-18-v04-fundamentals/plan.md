# Plan — webullapi4go v0.4 (Market Data Completeness)

## Run identity
- Date: 2026-09-18
- Slug: v04-fundamentals
- Mode: BUILD
- Repo: D:\github\webullapi4go
- Base: 99d505a (v0.3.0)
- Module: github.com/shing1211/webullapi4go
- Remotes: GitHub origin, Gitee gitee

## Goal
Complete the data/ package by implementing all remaining documented Market Data HTTP
endpoints (fundamentals, financial statements, fund data, crypto, extended screener,
corporate actions), fix the v0.3 signer tech debt, and update documentation.
Release as v0.4.0.

## Architecture decisions (locked)
1. All new endpoints use the existing data.Client with standard HMAC-SHA1 signing.
2. New endpoint files go in data/fundamentals.go (augmented), data/crypto.go (new),
   data/screener.go (augmented).
3. No new packages; no new dependencies beyond existing google.golang.org/grpc.

## Tasks
| ID  | Objective                                                              | Role      | Depends | Acceptance                                      | Size |
|-----|-----------------------------------------------------------------------|-----------|---------|-------------------------------------------------|------|
| T17.0 | Recon spike: probe sandbox availability for all candidate endpoints, classify a/b/c, update findings here | research | — | findings table written | M |
| T17.1 | Signer reconciliation: add DigestCase to SignParams, make events use internal/auth, delete events/sign.go | security | — | go test internal/auth... -v; go test events... -v | S |
| T17.2 | Fundamentals: capital flow, industry comparison, earnings/dividend calendar, SEC filings → data/fundamentals.go | data | T17.0 | offline fixtures; sandbox where T17.0 said usable | M |
| T17.3 | Financial statements: income, balance sheet, cashflow, indicators, alert, forecast EPS → data/fundamentals.go | data | T17.2 | offline fixtures; sandbox where usable | M |
| T17.4 | Tests + docs/fundamentals.md + example + README/CHANGELOG | tester/docs | T17.3 | go build/vet/lint/test; mkdocs --strict | M |
| T17.5 | Release v0.4.0 | release | T17.4 | tag on both remotes | S |

## Order of work
T17.0 (recon gate) ──┬── T17.1 (signer, independent)
                      └── T17.2 → T17.3 → T17.4 → T17.5

fundamentals.go is single-owner serial (T17.2→17.3). T17.1 is independent and can
run in parallel with the recon spike.

## Verification
- go build ./...
- go vet ./...
- go test ./... -count=1
- golangci-lint run ./...
- mkdocs build --strict

## What's NOT in v0.4
- Display Solution API — not a separate surface; covered by existing data/
- Broker API — deferred (separate surface)
- Go MCP server — deferred (post-v1.0)
- News REST list — no REST news endpoint; only SSE stream in data/news.go
- Crypto tick/depth/footprint — not in Python SDK
