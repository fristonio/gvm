GO := go
pkgs  = $(shell $(GO) list ./... | grep -v vendor)

build:
	@./build/build.sh

check_format:
	@echo "[*] Checking for formatting errors using gofmt"
	@./build/check_gofmt.sh

test: check_format

format:
	@echo "[*] Formatting code"
	@$(GO) fmt $(pkgs)

.PHONY: build format test check_format
