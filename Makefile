SHELL := /bin/sh

PROTO_SRC := docs/api/gateway.proto
PROTO_OUT := gateway/gen

.PHONY: proto gateway contracts-test lint fmt

proto:
	@mkdir -p $(PROTO_OUT)
	protoc --go_out=$(PROTO_OUT) --go-grpc_out=$(PROTO_OUT) $(PROTO_SRC)

gateway:
	cd gateway && go build ./...

contracts-test:
	forge test

lint:
	cd gateway && go vet ./...

fmt:
	cd contracts && forge fmt || true
	cd gateway && gofmt -w ./
	cd web && npm run fmt || true
