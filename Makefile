# Makefile for gotech-labs api
.DEFAULT_GOAL := help

# -----------------------------------------------------------------
#    ENV VARIABLE
# -----------------------------------------------------------------
SOURCEDIR  := .
SOURCES    := $(shell find . -type f -name '*.go' | grep -v vendor)
NOVENDOR   := $(shell go list $(SOURCEDIR)/... | grep -v vendor)
TOOLBINDIR := ./tools/bin

# Tools version
GOLANGCI_LINT_VERSION  := 1.45.0

# -----------------------------------------------------------------
#    Main targets
# -----------------------------------------------------------------

.PHONY: clean
clean: ## Remove temporary files
	@rm -rf cover.*
	@go clean --cache --testcache

.PHONY: fmt
fmt: ## Format all packages
	@go fmt $(NOVENDOR)

.PHONY: lint
lint: ## Code check
	@$(TOOLBINDIR)/golangci-lint run -v ./...

.PHONY: test
test: ## Run all the tests
	@go test -race -cover $(SOURCEDIR)/...

.PHONY: cover
cover: ## Run unit test and out coverage file for local environment
	@go test -race -timeout 10m -coverprofile=cover.out -covermode=atomic $(SOURCEDIR)/...
	@go tool cover -html=cover.out -o cover.html

.PHONY: mod-download
mod-download: ## Download go module packages
	@go mod download

.PHONY: mod-tidy
mod-tidy: ## Remove unnecessary go module packages
	@go mod tidy
	@go mod verify

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# -----------------------------------------------------------------
#    Setup targets
# -----------------------------------------------------------------

.PHONY: setup
setup: ## Setup dev tools
	@rm -f $(TOOLBINDIR)/golangci-lint
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLBINDIR) v$(GOLANGCI_LINT_VERSION)
