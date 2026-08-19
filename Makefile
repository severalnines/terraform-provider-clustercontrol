# =============================================================================
# Makefile - terraform-provider-clustercontrol
#
# Drop this file at the repo root (replaces the existing darwin-only Makefile).
#
# Tested against: Ubuntu 24.04 LTS, go 1.26.2 (see go.mod).
# =============================================================================

SHELL := /bin/bash

# --- Provider identity (used for the local dev-override install path) ------
HOSTNAME  := severalnines.com
NAMESPACE := severalnines
NAME      := clustercontrol
BINARY    := terraform-provider-$(NAME)
VERSION   ?= 0.2.23

# --- Toolchain / platform detection (was hardcoded to darwin_amd64) --------
GO      ?= go
GOOS    := $(shell $(GO) env GOOS)
GOARCH  := $(shell $(GO) env GOARCH)
OS_ARCH := $(GOOS)_$(GOARCH)

TARGET_DIR := $(HOME)/.terraform.d/plugins/$(HOSTNAME)/$(NAMESPACE)/$(NAME)/$(VERSION)/$(OS_ARCH)
TARGET     := ./bin/$(BINARY)

GOLANGCI_LINT ?= golangci-lint
TFPLUGINDOCS  := $(GO) run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

.DEFAULT_GOAL := install

.PHONY: default all build install clean \
        fmt vet lint tidy \
        test testacc \
        docs \
        release \
        sdk-update \
        help

default: install

## build: Compile the provider binary for the current OS/ARCH into ./bin
build:
	CGO_ENABLED=0 $(GO) build -o $(TARGET)

all: $(TARGET)

## install: Build and install into the local Terraform plugin dev-override dir
##          i.e. ~/.terraform.d/plugins/$(HOSTNAME)/$(NAMESPACE)/$(NAME)/$(VERSION)/$(OS_ARCH)
install: build
	mkdir -p $(TARGET_DIR)
	mv $(TARGET) $(TARGET_DIR)

## clean: Remove the installed dev-override binary and local build artifacts
clean:
	rm -rf $(TARGET_DIR)
	rm -rf ./bin

## fmt: gofmt the whole module
fmt:
	$(GO) fmt ./...

## vet: go vet the whole module
vet:
	$(GO) vet ./...

## lint: golangci-lint (skips gracefully if not installed)
lint:
	@if command -v $(GOLANGCI_LINT) >/dev/null 2>&1; then \
		$(GOLANGCI_LINT) run ./...; \
	else \
		echo "golangci-lint not found; install via: https://golangci-lint.run/welcome/install/"; \
	fi

## tidy: go mod tidy
tidy:
	$(GO) mod tidy

## test: Run unit tests (no ClusterControl instance required)
test:
	$(GO) test ./... -v -count=1

## testacc: Run acceptance tests against a real ClusterControl instance
##          Requires TF_ACC=1 plus cc_api_url / cc_api_user / cc_api_user_password
##          equivalents exported as env vars for the provider under test.
testacc:
	TF_ACC=1 $(GO) test ./... -v -count=1 -timeout 120m

## docs: Regenerate ./docs from schema descriptions + ./examples
docs:
	$(TFPLUGINDOCS)

## release: Local snapshot release via goreleaser (no publish/sign)
release:
	goreleaser release --clean --snapshot --skip=publish --skip=sign

## sdk-update: Bump the clustercontrol-client-sdk dependency and tidy
sdk-update:
	$(GO) get -u github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi
	$(GO) mod tidy

## help: Show this help
help:
	@echo "terraform-provider-clustercontrol targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'
	@echo ""
	@echo "Detected OS_ARCH: $(OS_ARCH)"
