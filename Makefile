.PHONY: all build run test clean test-compile test-run test-run-verbose

.DEFAULT_GOAL := all

BINARY_DIR := bin
BINARY := $(BINARY_DIR)/cfgctl.exe
PKG := ./...
CMD_PATH := cmd/cfgctl/main.go

GOOS ?= windows
GOARCH ?= amd64
GOFLAGS ?=

TESTBIN_DIR := testbin
PKGS := $(shell go list $(PKG))
STRIP_PKG_PREFIX_SED := 's|^github\.com/[^/]*/||; s|/|_|g'

all: build

build: | $(BINARY_DIR)
	@echo "Building $(BINARY) for $(GOOS)/$(GOARCH)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GOFLAGS) -o $(BINARY) $(CMD_PATH)
	@echo "Built $(BINARY)"

$(BINARY_DIR):
	@mkdir -p $(BINARY_DIR)

run: build
	@echo "Running $(BINARY)..."
	./$(BINARY)

test-compile:
	@echo "Compiling test binaries into $(TESTBIN_DIR) per package..."
	@mkdir -p $(TESTBIN_DIR)
	@for pkg in $(PKGS); do \
		pkgpath=$$(echo $$pkg | sed $(STRIP_PKG_PREFIX_SED)); \
		out="$(TESTBIN_DIR)/$${pkgpath}.test.exe"; \
		echo "Building $$pkg -> $${out}"; \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go test -c $(GOFLAGS) -o $$out $$pkg || true; \
		if [ -f "$$out" ]; then chmod +x "$$out" || true; else echo "No test binary for $$pkg (skipping)"; fi; \
	done
	@echo "Compiled test binaries in $(TESTBIN_DIR)"

test-run: test-compile
	@echo "Running compiled test binaries via cmd.exe..."
	@sh -c '\
	set -eu; \
	for winbin in $(TESTBIN_DIR)/*.test.exe; do \
		[ -e "$$winbin" ] || continue; \
		winpath=$$(cygpath -w "$$winbin" 2>/dev/null || echo "$$winbin"); \
		printf "Running %s via cmd.exe\n" "$$winpath"; \
		"$$winpath" -test.v || exit $$?; \
	done; \
	echo "Test binaries ran successfully."; \
'

test-run-verbose: test-compile
	@echo "Running compiled test binaries (verbose with coverage)..."
	@sh -c '\
	set -eu; \
	for testbin in $(TESTBIN_DIR)/*.test.exe; do \
		[ -e "$$testbin" ] || continue; \
		if [ ! -x "$$testbin" ]; then echo "Skipping non-executable $$testbin (likely cross-compiled)"; continue; fi; \
		echo "Running $$testbin"; \
		coverfile="$${testbin%.test.exe}.cover.out"; \
		"$$testbin" -test.v -test.coverprofile="$$coverfile" || exit $$?; \
		echo "Wrote $$coverfile"; \
	done; \
	echo "All test binaries ran successfully.";\
'

clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY)
	@rm -rf $(TESTBIN_DIR)
	@echo "Cleaned."
