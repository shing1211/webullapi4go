# Release Plan — v0.1.0

## Pre-flight (all verified)
- `go build ./...`, `go vet ./...`, `gofmt -l .`, `golangci-lint run ./...` clean.
- `go test ./... -count=1` passes.
- `mkdocs build --strict` succeeds; `docs/runs/` excluded from the site.
- Working tree contains only intended additions/modifications on `main`.
- Remotes: `origin` → `https://github.com/shing1211/webullapi4go.git` (GitHub).
  **Gitee remote not yet configured — URL to be confirmed.**

## Intended changes
All new source, tests, examples, docs, CI, and community files, plus:
- Modified: `.gitignore`, `README.md`.

## Steps
1. Stage intended files (no secrets; verify no credential strings).
2. Commit with Conventional Commits + `Signed-off-by` (`git commit -s`):
   ```
   feat: v0.1.0 — auth, market data HTTP, and MQTT streaming

   - HMAC-SHA1 request signing and token lifecycle (EnsureToken)
   - core HTTP client: typed errors, retry, rate limiting, circuit breaker, multi-region
   - market data HTTP: instruments, fundamentals, snapshot/tick/quotes/bars, footprint,
     NOII, screener, watchlist, options, news
   - MQTT streaming: connect, subscribe/unsubscribe, protobuf parsing, reconnect + resubscribe
   - examples, MkDocs docs site, CI, and community files

   Run: docs/runs/2026-09-18-webull-sdk-bootstrap
   ```
3. Tag annotated `v0.1.0`.
4. Push to GitHub `main` and tags.
5. Push to Gitee `main` and tags (after remote is configured).
6. Verify both remotes show the commit/tag.

## Safety
- No force-push, no rebase of shared branches. Stop and report on any rejection.
- Do not commit credentials or the shared sandbox keys.

## Rollback
- If a push partially succeeds, report the exact remote state; do not force-push.
