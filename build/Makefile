# ==============================================================================
# Aura Moby Core Daemon Engine Automation Makefile
# ==============================================================================

# Project Build Metadata Config
BINARY_NAME=aura-moby
BUILD_DIR=build
NATIVE_DIR=$(BUILD_DIR)/native

# Environment Compilation Toggles
export CGO_ENABLED=1
export CGO_LDFLAGS=-L$(shell pwd)/$(NATIVE_DIR) -laurabridge

.PHONY: all build native-lib go-daemon clean test help

# Default Target Execution Sequence
all: clean native-lib go-daemon

## native-lib: Invoke Zig build toolchain to assemble the 21.7% Assembly footprint
native-lib:
	@echo "==> Compiling native C-ABI assembly artifacts via Zig toolchain..."
	@mkdir -p $(NATIVE_DIR)
	zig build -Doptimize=ReleaseFast --summary failures

## go-daemon: Statically link compiled native objects into the core Go daemon binary
go-daemon: native-lib
	@echo "==> Linking native libraries and compiling primary Go application framework..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-extldflags=-static" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/dev
	@echo "==> Success! Monolithic binary generated at: $(BUILD_DIR)/$(BINARY_NAME)"

## test: Execute comprehensive tests across the entire Go system layer package spaces
test: native-lib
	@echo "==> Initiating unit test blocks on all engine sub-modules..."
	go test -v -race ./engine/... ./ai/...

## clean: Wipe localized compiler directories, object binaries, and testing artifacts
clean:
	@echo "==> Purging build directories and native cache repositories..."
	@rm -rf $(BUILD_DIR)
	@rm -rf .zig-cache
	@rm -f zig-out

## help: Print scannable targets and functional usage documentation details
help:
	@echo "Available target triggers for this Aura Moby deployment runtime:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'
