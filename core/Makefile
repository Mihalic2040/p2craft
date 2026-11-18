.DEFAULT_GOAL := all

.PHONY: all build install clean
.PHONY: linux-amd64 linux-arm64 windows-amd64

GO ?= go

BUILD_DIR := build
NATIVE_DIR := ../p2craft/src/main/resources/natives

# Cross-compilers
LINUX_AMD64_CC ?= gcc
LINUX_ARM64_CC ?= aarch64-linux-gnu-gcc
WINDOWS_AMD64_CC ?= x86_64-w64-mingw32-gcc


# Build everything and install into the mod.
all: install


# Build all targets.
build: linux-amd64 linux-arm64 windows-amd64


# Linux x86-64
linux-amd64:
	@mkdir -p $(BUILD_DIR)/linux-x86-64
	CGO_ENABLED=1 \
	GOOS=linux GOARCH=amd64 \
	CC=$(LINUX_AMD64_CC) \
	$(GO) build -buildmode=c-shared \
	-o $(BUILD_DIR)/linux-x86-64/libp2craft.so ./lib


# Linux ARM64
linux-arm64:
	@mkdir -p $(BUILD_DIR)/linux-aarch64
	CGO_ENABLED=1 \
	GOOS=linux GOARCH=arm64 \
	CC=$(LINUX_ARM64_CC) \
	$(GO) build -buildmode=c-shared \
	-o $(BUILD_DIR)/linux-aarch64/libp2craft.so ./lib


# Windows x86-64
windows-amd64:
	@mkdir -p $(BUILD_DIR)/windows-x86-64
	CGO_ENABLED=1 \
	GOOS=windows GOARCH=amd64 \
	CC=$(WINDOWS_AMD64_CC) \
	$(GO) build -buildmode=c-shared \
	-o $(BUILD_DIR)/windows-x86-64/p2craft.dll ./lib


# Install all native libraries into P2Craft resources.
install: build
	install -Dm755 $(BUILD_DIR)/linux-x86-64/libp2craft.so \
		$(NATIVE_DIR)/linux-x86-64/libp2craft.so

	install -Dm755 $(BUILD_DIR)/linux-aarch64/libp2craft.so \
		$(NATIVE_DIR)/linux-aarch64/libp2craft.so

	install -Dm755 $(BUILD_DIR)/windows-x86-64/p2craft.dll \
		$(NATIVE_DIR)/windows-x86-64/p2craft.dll

	@echo "All P2Craft native libraries installed."


clean:
	rm -rf $(BUILD_DIR)

	rm -f $(NATIVE_DIR)/linux-x86-64/libp2craft.so
	rm -f $(NATIVE_DIR)/linux-aarch64/libp2craft.so
	rm -f $(NATIVE_DIR)/windows-x86-64/p2craft.dll
