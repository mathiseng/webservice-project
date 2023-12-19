SHELL := /usr/bin/env bash -euo pipefail

.DEFAULT_GOAL := default

MKFILE_DIR = $(abspath $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
LOCAL_DIR := $(abspath $(MKFILE_DIR)/.local)

BIN_DIR	  := $(LOCAL_DIR)/bin
TMP_DIR	  := $(LOCAL_DIR)/tmp
SRC_DIR	  := $(MKFILE_DIR)


export PATH := $(BIN_DIR):$(PATH)

export GOMODCACHE = $(LOCAL_DIR)/cache/go
export GOTMPDIR = $(TMP_DIR)/go
export GOBIN = $(BIN_DIR)


default: clean install build run


.PHONY: install
install:
	mkdir -p $(GOTMPDIR)
	cd $(SRC_DIR) \
		&& go get -t ./...

.PHONY: run
run:
	cd $(SRC_DIR) \
	&& go run .

.PHONY: build $(BIN_DIR)/artifact.bin
build: $(BIN_DIR)/artifact.bin
$(BIN_DIR)/artifact.bin:
	cd $(SRC_DIR) \
	&& go build \
		-o $(@) \
		-ldflags "-X webservice/configuration.version=0.0.1" \
		$(SRC_DIR)/*.go

.PHONY: build-linux
build-linux: export GOOS := linux
build-linux: export GOARCH := amd64
build-linux: export CGO_ENABLED := 0
build-linux: $(BIN_DIR)/artifact.bin
	sha256sum $(BIN_DIR)/artifact.bin


.PHONY: test
.SILENT: test
test:
	cd $(SRC_DIR) \
	&& go clean -testcache \
	&& go test \
		-race \
		-v \
		$(SRC_DIR)/...


.PHONY: exec
exec:
	chmod +x $(BIN_DIR)/artifact.bin
	exec $(BIN_DIR)/artifact.bin


.PHONY: clean
clean:
	go clean -modcache
	rm -rf \
		$(LOCAL_DIR)
