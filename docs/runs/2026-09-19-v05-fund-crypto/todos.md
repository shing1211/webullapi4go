# Todos — v0.5 Phase 2: Fund Data, Crypto Data, Screener v2

Single source of truth. Statuses: `todo`, `doing`, `blocked`, `review`, `done`.

| ID | Task | Role | Status | Depends On | Acceptance |
|----|------|------|--------|-----------|------------|
| T1 | `data/fund_data.go` — Fund data endpoints, `map[string]string` rows | backend | todo | — | Compiles, passes vet/lint, endpoint paths documented |
| T2 | `data/crypto_data.go` — Crypto data endpoints, `map[string]string` rows | backend | todo | — | Compiles, passes vet/lint, endpoint paths documented |
| T3 | `data/screener_v2.go` — Screener v2 POST endpoint, best-effort schema | backend | todo | — | Compiles, passes vet/lint, path documented |
| T4 | Verification: go build, vet, fmt, test, lint | tester | todo | T1,T2,T3 | All commands pass |
| T5 | Update `docs/api.md` with new endpoints | docs | todo | T4 | api.md reflects new endpoints |
| T6 | `mkdocs build --strict` | tester | todo | T5 | Strict build passes |
| T7 | Commit and push to GitHub + Gitee | release | todo | T6 | Both remotes show new commit |
| T8 | Write `next-phase.md` | planner | todo | T7 | Next-phase candidates documented |
