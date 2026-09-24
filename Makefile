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
#   make vuln       run govulncheck in every Go module
#   make docs       build the MkDocs site (strict)
#   make generate   regenerate protobuf code
#   make proto-tools install protobuf codegen tools

MODULES := . broker examples/watchlist-cmd examples/broker-probe examples/futures-probe examples/options-multi-leg

.PHONY: build vet test test-race cover lint fuzz vuln docs generate proto-tools

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

vuln:
	@set -e; for module in $(MODULES); do go -C "$$module" run -mod=mod golang.org/x/vuln/cmd/govulncheck@latest ./...; done

docs:
	mkdocs build --strict

generate:
	buf generate

proto-tools:
	go install github.com/bufbuild/buf/cmd/buf@v1.73.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.2
