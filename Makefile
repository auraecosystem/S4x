# ==============================================================================
# Aura Moby Advanced CI/CD & Artifact Lifecycle Management Automation
# ==============================================================================

# Project Engine Configuration Metadata
BINARY_NAME=aura-moby
REGISTRY=ghcr.io/auraecosystem
IMAGE_TAG=latest
BUILD_DIR=build
NATIVE_DIR=$(BUILD_DIR)/native

# Primary Environment Directives
export CGO_ENABLED=1

.PHONY: all build-matrix docker-pack doc-sync clean help

# Default Target Execution Chain
all: clean build-matrix docker-pack doc-sync

## build-matrix: Cross-compile static C-ABI libraries and Go binaries for both AMD64 and ARM64
build-matrix:
	@echo "==> Initiating matrix compilation pipeline..."
	
	@echo "--> Building for Linux AMD64..."
	@mkdir -p $(NATIVE_DIR)/linux_amd64
	zig build -Dtarget=x86_64-linux -Doptimize=ReleaseFast --summary failures
	@mv -f zig-out/lib/* $(NATIVE_DIR)/linux_amd64/
	CGO_LDFLAGS="-L$(shell pwd)/$(NATIVE_DIR)/linux_amd64 -laurabridge -static" \
	GOOS=linux GOARCH=amd64 go build -ldflags="-extldflags=-static" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/dev

	@echo "--> Building for Linux ARM64 (Aarch64 Servers)..."
	@mkdir -p $(NATIVE_DIR)/linux_arm64
	zig build -Dtarget=aarch64-linux -Doptimize=ReleaseFast --summary failures
	@mv -f zig-out/lib/* $(NATIVE_DIR)/linux_arm64/
	CGO_LDFLAGS="-L$(shell pwd)/$(NATIVE_DIR)/linux_arm64 -laurabridge -static" \
	CC="zig cc -target aarch64-linux" GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/dev
	@echo "==> Matrix build complete. Artifacts stored in /$(BUILD_DIR)"

## docker-pack: Package compiled cross-platform binary distributions into local OCI/Docker layers
docker-pack: build-matrix
	@echo "==> Packaging local multi-arch production OCI container images..."
	docker buildx build --platform linux/amd64,linux/arm64 \
		-t $(REGISTRY)/$(BINARY_NAME):$(IMAGE_TAG) \
		--file ./Dockerfile . --load
	@echo "==> OCI multi-architecture container images securely built."

## doc-sync: Synchronize markdown source files to the local GitHub Pages layout
doc-sync:
	@echo "==> Refreshing static documentation assets..."
	@mkdir -p docs/doc
	cp README.md docs/index.md
	cp CLAUDE.md docs/claude_guidelines.md
	cp -R doc/* docs/doc/
	@echo "==> Docs synchronized. Commit and push the /docs folder to update auraecosystem.github.io/moby/."

## clean: Thoroughly purge localized binary assets, caching pipelines, and transient targets
clean:
	@echo "==> Resetting local construction state..."
	@rm -rf $(BUILD_DIR)
	@rm -rf docs
	@rm -rf .zig-cache
	@rm -rf zig-out

## help: Print all accessible structural targets with brief design parameters
help:
	@echo "Available structural lifecycle management commands:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'
