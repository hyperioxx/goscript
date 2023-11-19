.PHONY: help build test clean

BINARY := goscript

GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin

# Ensure GOPATH is correctly set.
GOPATH=$(go env GOPATH)

help:
	@echo "  make build - Compile and build the project."
	@echo "  make test  - Run tests."
	@echo "  make clean - Clean build files and caches."

build:
	@echo "  >  Building binary..."
	@go build -o $(GOBIN)/$(BINARY) ./cmd/slp

test:
	@echo "  >  Running tests..."
	go test -v ./...

clean:
	@echo "  >  Cleaning build cache"
	@go clean

install: build
	./install/install.sh