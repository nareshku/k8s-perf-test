# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=k8s-perf-test
BINARY_UNIX=$(BINARY_NAME)_unix
BIN_DIR=bin

# Build information
VERSION?=1.0.0
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: all build clean test coverage deps tidy fmt vet lint build-linux help

all: test build

build: 
	mkdir -p $(BIN_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) -v

clean: 
	$(GOCLEAN)
	rm -rf $(BIN_DIR)

test: 
	$(GOTEST) -v ./...

coverage: 
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

deps:
	$(GOGET) -v -t -d ./...

tidy:
	$(GOMOD) tidy

fmt:
	$(GOCMD) fmt ./...

vet:
	$(GOCMD) vet ./...

lint:
	golangci-lint run

build-linux:
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_UNIX) -v

# Install golangci-lint
install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin

# Run the application
run: build
	./$(BIN_DIR)/$(BINARY_NAME) -config config.yaml

# Show help
help:
	@echo "Available targets:"
	@echo "  all        : Run tests and build binary"
	@echo "  build      : Build binary for current platform in ./bin"
	@echo "  clean      : Clean build files"
	@echo "  test       : Run tests"
	@echo "  coverage   : Generate test coverage report"
	@echo "  deps       : Download dependencies"
	@echo "  tidy       : Tidy go.mod"
	@echo "  fmt        : Format code"
	@echo "  vet        : Run go vet"
	@echo "  lint       : Run golangci-lint"
	@echo "  build-linux: Build binary for Linux in ./bin"
	@echo "  run        : Build and run the application"
	@echo "  help       : Show this help message"

# Default target
.DEFAULT_GOAL := help 