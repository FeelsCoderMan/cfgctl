.PHONY: all build run clean test-compile test-run clean

.DEFAULT_GOAL := all

BINARY := bin/cfgctl.exe

PKG := ./...
CMD_PATH := cmd/cfgctl/main.go

all: build run

build:
	go build -o $(BINARY) $(CMD_PATH)

run: build
	./$(BINARY)

test:
	ifeq ($(origin TEST_TMP), undefined)
		@echo "Running: go test ./..."
		go test $(PKG)
	else
		@echo "Using TEST_TMP: $(TEST_TMP)"
		@TEMP=$(TEST_TMP) TMP=$(TEST_TMP) go test $(PKG)
	endif

test-compile:
	@echo "Compiling test binaries into ./testbin per package..."
	@mkdir -p testbin
	@for pkg in $(shell go list ./...); do \
		pkgpath=$$(echo $$pkg | sed 's|/|_|g'); \
		out=./testbin/$${pkgpath}.test.exe; \
		go test -c -o $$out $$pkg || exit $$?; \
	done
	@echo "Compiled test binaries in /.testbin"

test-run: test-compile
	@echo "Running compiled test binaries..."
	@for testbin in ./testbin/*.test.exe; do \
		echo "Running $$testbin"; \
		"$$testbin" -test.v || exit $$?; \
	done
	@echo "Test binaries ran successfully."

clean:
	rm -f $(BINARY)
	rm -rf testbin
