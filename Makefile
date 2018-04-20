GO := go
pkgs  = $(shell $(GO) list ./... | grep -v vendor)

format:
	@echo ">>> Formatting code"
	@$(GO) fmt $(pkgs)

build:
	@./build/build.sh

.PHONY: build
