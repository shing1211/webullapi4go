# Webull SDK toolchain.
#
# Usage:
#   make build      compile all packages
#   make test       run unit tests (sandbox tests need WEBULL_SANDBOX=1 and valid creds)
#   make test-race  run tests with the race detector
#   make cover      run tests with coverage profiling
#   make lint       run golangci-lint (includes gosec)
#   make fuzz       fuzz the data package deserializers
#   make vuln       run govulncheck
#   make docs       build the MkDocs site (strict)
#   make generate   regenerate protobuf code
#   make proto-tools install protobuf codegen tools

.PHONY: build test test-race cover lint fuzz vuln docs generate proto-tools

build:
	go build ./...

test:
	go test ./...

test-race:
	go test -race -count=1 ./...

cover:
	go test ./... -cover

lint:
	golangci-lint run

fuzz:
	go test -fuzz=FuzzDecodeQuote -fuzztime=1s ./data/...

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

docs:
	mkdocs build --strict

generate:
	buf generate

proto-tools:
	go install github.com/bufbuild/buf/cmd/buf@v1.73.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.2
