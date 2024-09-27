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

.PHONY: build-docker-server
build-docker-server:
	docker build --build-arg TARGETARCH=amd64 -t llmariner/api-usage-server:latest -f build/server/Dockerfile .
