# Plan — webullapi4go Phase 0: API Recon Spike

## Run Identity
- Date: 2026-09-18
- Slug: v05-recon
- Mode: BUILD
- Repo: D:\github\webullapi4go
- Base: adb4243 (v0.4.0 + docs reconciliation)
- Module: github.com/shing1211/webullapi4go
- Remotes: GitHub origin, Gitee gitee

## Goal
Fill every unknown gap in the v0.5 + v0.6 surface via live API probe using the partner
account. Outputs: findings.md with confirmed paths, schemas, and auth details for all
undocumented endpoints. No implementation code written in this phase.

## Architecture Decisions (locked)
1. New `brokerfd/` top-level package for Broker API (FD) HTTP + gRPC events.
2. New `data/fund.go`, `data/crypto.go`, `data/corporate_actions.go`, `data/screener_v2.go`.
3. All new REST endpoints use HMAC-SHA1 signing (same as existing `data/`).
4. All new gRPC events use HMAC-SHA256 signing with lowercase hex (same as `events/`).
5. Corporate Actions / Fund / Crypto use `map[string]string` rows until schema confirmed.
6. Partner account credentials available for live probing.

## What's Known vs Unknown

### Known (implement directly)
- Corporate Actions: 2 endpoints, paths confirmed from API reference
- Instrument v3: new path `/market-data/instruments/stocks/list` confirmed
- Existing Display Solution endpoints: same paths as existing data/ — no new work
- Broker FD gRPC events: 12 string-based event types, proto differs from existing events.proto

### Unknown (probe live)
- Fund data endpoints and schemas
- Crypto tick/depth/orderbook endpoints and schemas
- Screener v2 path and schema (`wlas/screener/ng/query` in Python SDK)
- Broker FD HTTP schemas for all 12 categories
- Broker FD gRPC proto files (event contract structures)
- Display Solution host confirmation (`co-branding-openapi...`)

## Phase 0 Tasks

| ID | Objective | Role | Depends | Size |
|----|-----------|-------|---------|------|
| T0.1 | Display Solution host + auth probe | backend | — | S |
| T0.2 | Corporate Actions schemas | backend | T0.1 | S |
| T0.3 | Fund data endpoints | backend | T0.1 | M |
| T0.4 | Crypto data endpoints | backend | T0.1 | M |
| T0.5 | Screener v2 endpoint | backend | T0.1 | S |
| T0.6 | Broker FD HTTP schemas (all 12 categories) | backend | T0.1 | L |
| T0.7 | Broker FD gRPC proto + event types | architect | T0.1 | M |

T0.1 is the gate — it must run first to establish which host/creds work. T0.2–T0.7
can run in parallel once T0.1 completes.

## Verification
- `docs/runs/2026-09-18-v05-recon/findings.md` written with all confirmed paths,
  schemas, and auth details.
- No implementation code written.
- Run folder contains: plan.md, todos.md, findings.md.
