SHELL := /bin/bash
GO ?= go
GOBIN := $(shell $(GO) env GOBIN)
ifeq ($(GOBIN),)
  GOBIN := $(shell $(GO) env GOPATH)/bin
endif

GOLANGCI_LINT := $(GOBIN)/golangci-lint
GOVULNCHECK := $(GOBIN)/govulncheck

.PHONY: all tools fmt fmt-check lint vet test vuln ci tidy

all: tools fmt lint vet test vuln

tools:
	@echo "Installing golangci-lint v2 (latest)..."
	@rm -f $(GOLANGCI_LINT)
	@$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	@echo "Installing govulncheck (latest)..."
	@$(GO) install golang.org/x/vuln/cmd/govulncheck@latest

fmt:
	$(GO) fmt ./...

fmt-check:
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "These files need gofmt:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

lint: tools
	$(GOLANGCI_LINT) run

vet:
	$(GO) vet ./...

test:
	$(GO) test -race -cover ./...

vuln: tools
	$(GOVULNCHECK) ./...

tidy:
	$(GO) mod tidy

ci: fmt-check lint vet test vuln
