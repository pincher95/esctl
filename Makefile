# Makefile for esctl

BINARY    := esctl
PKG       := github.com/pincher95/esctl
BUILDINFO := $(PKG)/internal/buildinfo

# Where `make install` moves the binary. Defaults to `go env GOBIN`, falling back
# to `$(go env GOPATH)/bin` (e.g. ~/go/bin). Override with `make install INSTALL_DIR=...`.
GOBIN       := $(shell go env GOBIN)
INSTALL_DIR ?= $(if $(GOBIN),$(GOBIN),$(shell go env GOPATH)/bin)

# Version metadata injected into internal/buildinfo via -ldflags -X.
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT   ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE     ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
BUILT_BY ?= make

LDFLAGS := -s -w \
	-X $(BUILDINFO).Version=$(VERSION) \
	-X $(BUILDINFO).Commit=$(COMMIT) \
	-X $(BUILDINFO).Date=$(DATE) \
	-X $(BUILDINFO).BuiltBy=$(BUILT_BY)

.DEFAULT_GOAL := build

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

.PHONY: fmt
fmt: ## Format code (gofmt -s -w)
	gofmt -s -w .

.PHONY: fmt-check
fmt-check: ## Fail if any file is not gofmt-ed
	@out=$$(gofmt -l -s .); \
	if [ -n "$$out" ]; then echo "not gofmt-ed:"; echo "$$out"; exit 1; fi

.PHONY: lint
lint: fmt-check ## Run go vet (gating) + golangci-lint (advisory)
	go vet ./...
	@# golangci-lint currently reports pre-existing findings, so it is advisory
	@# (non-fatal) here. Remove the trailing "|| true" to make it a hard gate.
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run || true; \
	else \
		echo "golangci-lint not found; ran go vet only"; \
	fi

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: test-race
test-race: ## Run tests with the race detector
	go test -race ./...

.PHONY: cover
cover: ## Run tests and print a coverage summary
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

.PHONY: build
build: ## Build the esctl binary into the repo root
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

.PHONY: install
install: build ## Build and move the binary into INSTALL_DIR ($(INSTALL_DIR))
	@mkdir -p "$(INSTALL_DIR)"
	mv $(BINARY) "$(INSTALL_DIR)/$(BINARY)"
	@echo "installed $(BINARY) -> $(INSTALL_DIR)/$(BINARY)"

.PHONY: tidy
tidy: ## Tidy go.mod / go.sum
	go mod tidy

.PHONY: clean
clean: ## Remove build artifacts
	rm -f $(BINARY) coverage.out

.PHONY: all
all: fmt lint test build ## Format, lint, test, and build
