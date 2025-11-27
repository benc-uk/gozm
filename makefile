ROOT_DIR := $(shell git rev-parse --show-toplevel)
DEV_DIR := $(ROOT_DIR)/.dev

.DEFAULT_GOAL := help

.PHONY: help build test run watch lint tidy install

help: # 💬 Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: # 🔨 Build the Go binary
	go build -o bin/gozm main.go

test: # 🧪 Run tests
	go test -v ./...

run: # 🚀 Run the application
	go run cmd/cli/main.go ./test/hello.z3

watch: # 👀 Watch for file changes and rebuild
	go tool -modfile=.dev/tools.mod air -c $(DEV_DIR)/air.toml

lint: # ✨ Run golangci-lint
	go tool -modfile=.dev/tools.mod golangci-lint run --config $(DEV_DIR)/golangci.yaml

tidy: # 🧹 Tidy Go modules
	go mod tidy
	go mod tidy -modfile=$(DEV_DIR)/tools.mod

install: # 📦 Install dependencies
	go mod download
	go mod download -modfile=$(DEV_DIR)/tools.mod

hello:
	./tools/inform6.exe -v3 ./test/hello.inf ./test/hello.z3
	./tools/unz ./test/hello.z3 > ./test/hello.dump 