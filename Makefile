# Ensure Make is run with bash shell as some syntax below is bash-specific
SHELL:=/usr/bin/env bash

# allow overwriting the default `go` value with the custom path to the go executable
GOEXE ?= go

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell $(GOEXE) env GOBIN))
	GOBIN=$(shell $(GOEXE) env GOPATH)/bin
else
	GOBIN=$(shell $(GOEXE) env GOBIN)
endif

GOROOT=$(shell $(GOEXE) env GOROOT)

.PHONY: all test clean lint gosec

all: wif-go

SRCS = $(shell find pkg -iname "*.go"|grep -v pkg/generated) $(GENSRC)
TOOLS_DIR := utils/tools
TOOLS_BIN_DIR := $(abspath $(TOOLS_DIR)/bin)

KO_PREFIX ?= ghcr.io/loicsikidi/wif-go

GO_MODULE=$(shell head -1 go.mod | cut -f2 -d ' ')
GOTEST=go test
GOCOVER=go tool cover

GENSRC = pkg/generated/protobuf/%.go
PROTOBUF_DEPS = $(shell find . -iname "*.proto")

PLAYGROUND_DIR := web
PLAYGROUND_PUBLIC_DIR := $(abspath $(PLAYGROUND_DIR)/public)
PLAYGROUND_JS_DIR := $(abspath $(PLAYGROUND_DIR)/src/js)
WASM_DIR := cmd/wasm
SERVER_DIR := cmd/server
SERVER_STATIC_DIR := $(abspath $(SERVER_DIR)/kodata)
SERVER_EMBED_STATIC_DIR := $(abspath $(SERVER_DIR)/static)

# Set version variables for LDFLAGS
GIT_VERSION ?= $(shell git describe --tags --always --dirty)
GIT_HASH ?= $(shell git rev-parse HEAD)

SOURCE_DATE_EPOCH ?= $(shell git log -1 --no-show-signature --pretty=%ct)
ifdef SOURCE_DATE_EPOCH
    BUILD_DATE ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u "$(DATE_FMT)")
else
    BUILD_DATE ?= $(shell date "$(DATE_FMT)")
endif

MODULE_VERSION=main
LDFLAGS=-s -w -X $(MODULE_VERSION).Version=$(GIT_VERSION)

# Binaries
PROTOC-GEN-GO := $(TOOLS_BIN_DIR)/protoc-gen-go

$(GENSRC): $(PROTOC-GEN-GO) $(PROTOBUF_DEPS)
	mkdir -p pkg/generated/protobuf
	protoc --plugin=protoc-gen-go=$(TOOLS_BIN_DIR)/protoc-gen-go \
	       --go_opt=module=$(GO_MODULE) --go_out=. \
		   -I . $(PROTOBUF_DEPS)

lint: ## Runs golangci-lint
	$(GOBIN)/golangci-lint run -v ./...

gosec: ## Runs gosec
	$(GOBIN)/gosec ./...

gen-proto: $(GENSRC) ## Generate protobuf files

gen-hack: ## [HACK] Copy kodata to static dir
	cp -pr $(SERVER_STATIC_DIR) $(SERVER_EMBED_STATIC_DIR)

wif-go: $(SRCS) ## Build Playground for local tests
	$(GOEXE) build -trimpath -ldflags "$(LDFLAGS)" ./pkg/compiler

test: ## Runs go test
	$(GOTEST) -v -coverprofile=coverage.txt -covermode=atomic ./...

coverage: ## Runs go test with coverage
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out

clean: ## Clean the workspace
	rm -rf dist
	rm -rf wif-go

clean-gen: clean
	rm -rf $(shell find pkg/generated -iname "*.go")

## --------------------------------------
## Tooling Binaries
## --------------------------------------

$(PROTOC-GEN-GO): $(TOOLS_DIR)/go.mod
	cd $(TOOLS_DIR); go build -trimpath -tags=tools -o $(TOOLS_BIN_DIR)/protoc-gen-go google.golang.org/protobuf/cmd/protoc-gen-go

## --------------------------------------
## Modules
## --------------------------------------

.PHONY: modules
modules: ## Runs go mod to ensure modules are up to date.
	$(GOEXE) mod tidy

## --------------------------------------
## Generate playground sources
## --------------------------------------
.PHONY: build-playground
build-wasm: ## Build wasm components
	cd $(WASM_DIR); \
	GOOS=js GOARCH=wasm $(GOEXE) build -trimpath -ldflags "$(LDFLAGS)" -o $(PLAYGROUND_PUBLIC_DIR)/wif-go.wasm \
	&& cp "$(GOROOT)/misc/wasm/wasm_exec.js" $(PLAYGROUND_JS_DIR)
build-vue: ## Build the playground (frontend)
	cd $(PLAYGROUND_DIR); npm i && npm run build -- --outDir=$(SERVER_STATIC_DIR)
build-playground: build-wasm build-vue ## Build the playground static sources
clean-playground: ## Clean playground static sources
	rm -rf $(SERVER_STATIC_DIR)

## --------------------------------------
## Devtools
## --------------------------------------
.PHONY: dev
dev: build-wasm ## Run all steps required for local development (cf. playground).
	cd $(PLAYGROUND_DIR); npm i && npm run dev

## --------------------------------------
## Images with ko
## --------------------------------------
export KO_DOCKER_REPO=$(KO_PREFIX)

KOCACHE_PATH=/tmp/ko
OCI_LABELS=--image-label org.opencontainers.image.created=$(BUILD_DATE) \
           --image-label org.opencontainers.image.description="Container hosting a playground in order to test Workload Identity Federation mapping"
define create_kocache_path
  mkdir -p $(KOCACHE_PATH)
endef

.PHONY: ko
ko: build-playground
	$(create_kocache_path)
	LDFLAGS="$(LDFLAGS)" GIT_HASH=$(GIT_HASH) GIT_VERSION=$(GIT_VERSION) \
	KOCACHE=$(KOCACHE_PATH) ko build --base-import-paths \
		--platform=all --tags $(GIT_VERSION) --tags $(GIT_HASH) \
		$(ARTIFACT_HUB_LABELS) \
		"github.com/loicsikidi/wif-go/cmd/server"

.PHONY: ko-local
ko-local: build-playground
	$(create_kocache_path)
	KO_DOCKER_REPO=ko.local LDFLAGS="$(LDFLAGS)" GIT_HASH=$(GIT_HASH) GIT_VERSION=$(GIT_VERSION) \
	KOCACHE=$(KOCACHE_PATH) ko build --base-import-paths \
		--platform=linux/amd64 --tags $(GIT_VERSION) --tags $(GIT_HASH) \
		$(ARTIFACT_HUB_LABELS) \
		"github.com/loicsikidi/wif-go/cmd/server"

##################
# help
##################

help: ## Display help
	@awk -F ':|##' \
		'/^[^\t].+?:.*?##/ {\
			printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
		}' $(MAKEFILE_LIST) | sort
