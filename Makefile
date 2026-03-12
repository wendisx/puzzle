BINARY_NAME = puzzle
BUILD_DIR   = ./bin
DIST_DIR    = ./dist

VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0")
BUILD     = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S')
MODULE_PATH = github.com/wendisx/puzzle/internal/command

GO          = go
GO_FLAGS    = -v
LDFLAGS     = -s -w \
              -X $(MODULE_PATH).VERSION=$(VERSION) \
              -X $(MODULE_PATH).BUILD=$(BUILD) \
              -X $(MODULE_PATH).DATE=$(DATE)

.PHONY: all build clean test fmt help

all: fmt test build

help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

fmt:
	$(GO) fmt ./...

test:
	$(GO) test ./... -cover

build: clean
	@mkdir -p $(BUILD_DIR)
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	$(GO) build $(GO_FLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) .

clean:
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)

cross-compile:
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .
