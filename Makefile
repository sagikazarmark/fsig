# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

# Project variables
PACKAGE = $(shell echo $${PWD\#\#*src/})
BINARY_NAME = $(shell basename $$PWD)

# Build variables
BUILD_DIR = build
VERSION ?= $(shell git rev-parse --abbrev-ref HEAD)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)
LDFLAGS = -ldflags "-w -X main.Version=${VERSION} -X main.CommitHash=${COMMIT_HASH} -X main.BuildDate=${BUILD_DATE}"

# Dependency versions
DEP_VERSION = 0.5.0
GOLANGCI_VERSION = 1.10.2
GORELEASER_VERSION = 0.84.0

bin/dep: ## Install dep
	@mkdir -p ./bin/
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | INSTALL_DIRECTORY=./bin DEP_RELEASE_TAG=v${DEP_VERSION} sh

.PHONY: vendor
vendor: bin/dep ## Install dependencies
	bin/dep ensure -vendor-only

.PHONY: clean
clean: ## Clean the working area and the project
	rm -rf bin/ ${BUILD_DIR}/ vendor/

.PHONY: build
build: ## Build a binary
	CGO_ENABLED=0 go build -tags '${TAGS}' ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ${PACKAGE}

.PHONY: check
check: test lint ## Run tests and linters

.PHONY: test
test: ## Run all tests
	go test -tags 'unit integration acceptance' ${ARGS} ./...

.PHONY: test-%
test-%: ## Run a specific test suite
	go test -tags '$*' ${ARGS} ./...

bin/golangci-lint: ## Install golangci linter
	@mkdir -p ./bin/
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b ./bin/ v${GOLANGCI_VERSION}

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	bin/golangci-lint run

bin/goreleaser: ## Install goreleaser
	@mkdir -p ./bin/
	curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | bash -s -- v${GORELEASER_VERSION}

.PHONY: release
release: bin/goreleaser ## Release current tag
	bin/goreleaser

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Variable outputting/exporting rules
var-%: ; @echo $($*)
varexport-%: ; @echo $*=$($*)
