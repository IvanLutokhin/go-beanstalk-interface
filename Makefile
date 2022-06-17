# Go related variables
GOBASE := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# Targets
.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help screen
	@echo "Available targets:"
	@grep -h -E '(^[a-zA-Z_-]+:.*?##.*$$)|(^##)' $(firstword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-24s\033[0m %s\n", $$1, $$2}'

.PHONY: generate
generate: ## Generate code
	@ROOT_DIR=$(GOBASE) go generate ./...

.PHONY: test
test: ## Test
	@go test -v -race -cover -coverprofile=coverage.out ./...

.PHONY: coverage-html
coverage-html: coverage.out ## Make html coverage report
	@go tool cover -html=coverage.out -o coverage.html

build: goreleaser ## Build binaries
	@goreleaser build --snapshot --rm-dist

.PHONY: goreleaser
goreleaser:
	@which goreleaser 2>&1 >/dev/null || go install github.com/goreleaser/goreleaser@latest
