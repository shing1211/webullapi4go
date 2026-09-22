# Todos: HK Sandbox Probe 2026-09-22

## Probe Execution
- [x] Set up HK sandbox credentials in .env
- [x] Probe `GetFuturesProductCodes` (HK) — PASS, 89 products returned
- [x] Probe `GetFuturesInstruments` (HK) — FAIL: `json: cannot unmarshal number into Go struct field FuturesInstrument.unit of type string`
- [x] Probe multi-leg options strategies via `PreviewOrder`
- [x] Probe HK options discovery endpoints
- [x] Probe Broker API HK paths
- [x] Confirm US-only endpoints return 404

## Bug Fixes
- [ ] `FuturesInstrument.Unit` flexible type
- [ ] `examples/options-multi-leg/main.go:193` `client_order_id` length fix

## Documentation
- [x] Write plan.md
- [ ] Write todos.md (this file)
- [ ] Write report.md
- [ ] Update `docs/runs/index.md`
- [ ] Update CHANGELOG.md

## Verification
- [ ] Run `go build ./... && go vet ./...`
- [ ] Run `gofmt -l .`
- [ ] Run `go test ./...`
- [ ] Run `go test -race -count=1 ./...`
- [ ] Run `golangci-lint run ./...`
- [ ] Run `mkdocs build --strict`

## Commit
- [ ] Commit all changes with Conventional Commit
- [ ] Push to both remotes
