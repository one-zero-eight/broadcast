GO ?= go

.PHONY: all fmt lint vet test vuln ci

all: fmt lint vet test vuln

fmt:
	$(GO) fmt ./...

lint:
	golangci-lint run

vet:
	$(GO) vet ./...

test:
	$(GO) test -race -cover ./...

vuln:
	$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	$$GOBIN/govulncheck ./...

ci: all

