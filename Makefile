# Go related variables
GOBASE := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
GOAPI := $(GOBASE)/api

# Targets
.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help screen
	@echo "Available targets:"
	@grep -h -E '(^[a-zA-Z_-]+:.*?##.*$$)|(^##)' $(firstword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-24s\033[0m %s\n", $$1, $$2}'

.PHONY: generate-api-spec
generate-api-spec: swagger ## Generate API specification
	@swagger generate spec -o $(GOAPI)/system/v1/swagger.json --scan-models

.PHONY: swagger
swagger:
	@which swagger 2>&1 >/dev/null || go install github.com/go-swagger/go-swagger/cmd/swagger@latest

build: goreleaser ## Build binaries
	@goreleaser build --snapshot --rm-dist

.PHONY: goreleaser
goreleaser:
	@which goreleaser 2>&1 >/dev/null || go install github.com/goreleaser/goreleaser@latest
