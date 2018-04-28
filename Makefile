GO := go
pkgs  = $(shell $(GO) list ./... | grep -v vendor)

check_format:
	@echo "[*] Checking for formatting errors using gofmt"
	@./build/check_gofmt.sh

test: check_format

format:
	@echo "[*] Formatting code"
	@$(GO) fmt $(pkgs)

build:
	@./build/build.sh

.PHONY: build format test check_format
