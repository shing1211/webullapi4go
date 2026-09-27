# Webull SDK toolchain.
#
# Usage:
#   make build      compile all packages in every Go module
#   make vet        run go vet in every Go module
#   make test       run unit tests (sandbox tests need WEBULL_SANDBOX=1 and valid creds)
#   make test-race  run tests with the race detector in every Go module
#   make cover      run tests with coverage profiling in every Go module
#   make lint       run golangci-lint (includes gosec)
#   make fuzz       fuzz the data package deserializers
#   make vuln       run govulncheck in every Go module (pinned tool version)
#   make docs       build the MkDocs site (strict)
#   make citations  validate the file:line citations in the status documents
#   make conformance-fixtures         gate the committed fixtures against the docs cache
#   make conformance-fixtures-update  regenerate them, for adopting an intentional change
#   make conformance-gate             gate each SDK DTO against its documented wire shape
#   make conformance-report           print the observed divergence set, read-only
#   make generate   regenerate protobuf code
#   make proto-tools install protobuf codegen tools

MODULES := . broker examples/watchlist-cmd examples/broker-probe examples/futures-probe examples/options-multi-leg

.PHONY: build vet test test-race cover lint fuzz vuln docs citations conformance-fixtures conformance-fixtures-update conformance-gate conformance-report generate proto-tools

build:
	@set -e; for module in $(MODULES); do go -C "$$module" build -mod=readonly ./...; done

vet:
	@set -e; for module in $(MODULES); do go -C "$$module" vet -mod=readonly ./...; done

test:
	@set -e; for module in $(MODULES); do go -C "$$module" test -mod=readonly ./...; done

test-race:
	@set -e; for module in $(MODULES); do go -C "$$module" test -mod=readonly -race -count=1 ./...; done

cover:
	@set -e; for module in $(MODULES); do go -C "$$module" test -mod=readonly ./... -cover; done

lint:
	golangci-lint run

fuzz:
	go test -fuzz=FuzzDecodeQuote -fuzztime=1s ./data/...
	go test -fuzz=FuzzDecodeSnapshot -fuzztime=1s ./data/...
	go test -fuzz=FuzzDecodeTick -fuzztime=1s ./data/...

# Pinned so scan results are reproducible across runs. The vulnerability
# database is still fetched live from vuln.go.dev, so pinning the tool does not
# freeze or stale the reported findings. Bump GOVULNCHECK_VERSION deliberately.
GOVULNCHECK_VERSION := v1.8.0

vuln:
	@set -e; for module in $(MODULES); do go -C "$$module" run -mod=mod golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...; done

docs:
	mkdocs build --strict

# Reports only; a non-zero exit fails the build. Cited paths resolve against the
# repository root derived from the tool's own location, not the current
# directory, so `make -C <repo> citations` works from anywhere.
citations:
	python tools/citations/check.py

# Wire-conformance fixtures. The tree under conformance/testdata is a committed
# snapshot of Webull's published response schemas, so the interesting failure is
# the tree drifting away from the documentation without anybody noticing.
# conformance-fixtures therefore regenerates in memory and fails on any
# difference, which makes drift a signal; it never rewrites the tree, so a stale
# fixture cannot quietly become a fresh one. conformance-fixtures-update is the
# only way to adopt new output, and it is deliberately a separate verb.
#
# Both skip with exit 0 when the docgen cache is absent. The cache is gitignored,
# so a fresh checkout has none and cannot judge drift either way; failing there
# would train the gate to be ignored. Populate it with tools/webull-docgen, or
# point WEBULL_DOCGEN_CACHE at a populated directory.
conformance-fixtures:
	python tools/conformance/gen_fixtures.py --check --self-test

conformance-fixtures-update:
	python tools/conformance/gen_fixtures.py --write --self-test

# The per-DTO wire-conformance gate. It compares each SDK response type against
# the documented shape in conformance/testdata and fails unless the observed
# divergence set is exactly the set recorded in conformance/known-divergences.json,
# so a new defect fails at once and a fixed one fails too, because both need a
# person to read the tree. conformance-report is the same comparison run in
# read-only mode: it prints the observed set, grouped by check, which is what a
# baseline edit has to be made against.
conformance-gate:
	go test ./conformance/ -run 'TestSymbolTableCoversTheManifest|TestSymbolTableResolvesItsDecodeTargets|TestKnownDivergenceBaseline|TestComparisonBites' -count=1

conformance-report:
	go test ./conformance/ -run TestObservedDivergenceReport -count=1 -v

generate:
	buf generate

proto-tools:
	go install github.com/bufbuild/buf/cmd/buf@v1.73.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.2
