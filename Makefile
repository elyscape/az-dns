OS = $(shell uname -s)

GO = go
GOPATH ?= $(shell $(GO) env GOPATH)

TARGET_DIR = .targets
TARGET_FILES = deps deps-vendor

SOURCE_FILES = ./...
PACKAGE_NAME = github.com/elyscape/az-dns

BIN_DIR = /usr/local/bin
DEP = $(BIN_DIR)/dep
GOMETALINTER = $(BIN_DIR)/gometalinter
GORELEASER = $(BIN_DIR)/goreleaser

COVER_FILE = $(TARGET_DIR)/coverage.out
COVER_FLAGS = -coverpkg=./... -coverprofile=$(COVER_FILE)
TEST_FLAGS = -v -race
TEST_PATTERN = ''

TIMESTAMP ?= $(shell date -u +%FT%TZ)

ifeq ($(OS), Darwin)
ifeq ($(BIN_DIR), $(shell brew --prefix)/bin)
USE_HOMEBREW = 1
endif
endif

$(DEP):
ifdef USE_HOMEBREW
	brew install dep
else
	curl -fsSL https://raw.githubusercontent.com/golang/dep/master/install.sh | INSTALL_DIRECTORY=$(BIN_DIR) sh
endif

$(GOMETALINTER):
	curl -fsSL https://install.goreleaser.com/github.com/alecthomas/gometalinter.sh | sh -s -- -b $(BIN_DIR)

$(GORELEASER):
	curl -fsSL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh -s -- -b $(BIN_DIR)

setup: $(DEP) $(GOMETALINTER) $(GORELEASER)
.PHONY: setup

$(TARGET_DIR):
	@mkdir -p $(TARGET_DIR)

$(addprefix $(TARGET_DIR)/,$(TARGET_FILES)): $(TARGET_DIR)

$(TARGET_DIR)/deps: $(DEP) Gopkg.toml Gopkg.lock
	$(DEP) ensure
	@touch $@

$(TARGET_DIR)/deps-vendor: $(DEP) Gopkg.lock
	$(DEP) ensure -vendor-only
	@touch $@

deps: $(TARGET_DIR)/deps
.PHONY: deps

deps-vendor: $(TARGET_DIR)/deps-vendor
.PHONY: deps-vendor

test: deps-vendor
test $(COVER_FILE):
	$(GO) test $(TEST_FLAGS) $(COVER_FLAGS) $(SOURCE_FILES) -run $(TEST_PATTERN)
.PHONY: test

open-cover: $(COVER_FILE)
	$(GO) tool cover -html=$(COVER_FILE) -o $(TARGET_DIR)/cover.html
.PHONY: open-cover

lint: $(GOMETALINTER)
	$(GOMETALINTER) -t --vendor $(SOURCE_FILES)
.PHONY: lint

az-dns: $(TARGET_DIR)/deps-vendor
	$(GO) build

build: az-dns
.PHONY: build

snapshot: $(GORELEASER)
	$(GORELEASER) --snapshot --rm-dist
.PHONY: snapshot

release: test $(GORELEASER)
	$(GORELEASER) --release-notes=<(./generate-changelog.rb)
.PHONY: release

clean:
	-rm -rvf az-dns dist $(TARGET_DIR)
.PHONY: clean

all: test lint build
.PHONY: all

ci: test lint snapshot
.PHONY: ci

ci-release: test lint release
.PHONY: ci-release

.DEFAULT_GOAL = all
