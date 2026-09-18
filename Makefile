# Webull market-data protobuf codegen.
#
# Requires:
#   - buf:               go install github.com/bufbuild/buf/cmd/buf@v1.73.0
#   - protoc-gen-go:     go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.1
#   - protoc-gen-go-grpc: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.2
#
# The generated code is committed under gen/ so builds do not require codegen.

.PHONY: generate
generate:
	buf generate

.PHONY: proto-tools
proto-tools:
	go install github.com/bufbuild/buf/cmd/buf@v1.73.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.2
