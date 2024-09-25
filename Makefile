.PHONY: default
default: test

include common.mk

.PHONY: generate
generate:
	buf generate

.PHONY: test
test: go-test-all

.PHONY: lint
lint: go-lint-all git-clean-check

.PHONY: build-server
build-server:
	go build -o ./bin/server ./server/cmd/
