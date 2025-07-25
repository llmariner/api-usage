.PHONY: default
default: test

include common.mk

.PHONY: generate
generate: buf-generate-all typescript-compile

.PHONY: test
test: go-test-all

.PHONY: lint
lint: go-lint-all helm-lint git-clean-check

.PHONY: build-cleaner
build-cleaner:
	go build -o ./bin/cleaner ./cleaner/cmd/

.PHONY: build-server
build-server:
	go build -o ./bin/server ./server/cmd/

.PHONY: build-docker-cleaner
build-docker-cleaner:
	docker build --build-arg TARGETARCH=amd64 -t llmariner/api-usage-cleaner:latest -f build/cleaner/Dockerfile .

.PHONY: build-docker-server
build-docker-server:
	docker build --build-arg TARGETARCH=amd64 -t llmariner/api-usage-server:latest -f build/server/Dockerfile .

.PHONY: build-all
build-all: build-server build-cleaner

.PHONY: build-docker-all
build-docker-all: build-docker-server build-docker-cleaner

.PHONY: check-helm-tool
check-helm-tool:
	@command -v helm-tool >/dev/null 2>&1 || $(MAKE) install-helm-tool

.PHONY: install-helm-tool
install-helm-tool:
	go install github.com/cert-manager/helm-tool@latest

.PHONY: generate-chart-schema
generate-chart-schema: check-helm-tool
	@cd ./deployments/server && helm-tool schema > values.schema.json
	@cd ./deployments/cleaner && helm-tool schema > values.schema.json

.PHONY: helm-lint
helm-lint: generate-chart-schema
	cd ./deployments/server && helm-tool lint
	helm lint ./deployments/server
	cd ./deployments/cleaner && helm-tool lint
	helm lint ./deployments/cleaner
