GOLANGCI_LINT_VERSION := v2.14.0
GOBIN := $(shell go env GOPATH)/bin
GOLANGCI_LINT := $(GOBIN)/golangci-lint

.PHONY: build vet fmt fmt-check lint test check

build:
	go build ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

fmt-check:
	@out="$$(gofmt -l .)"; if [ -n "$$out" ]; then echo "gofmt required for:"; echo "$$out"; exit 1; fi

# go install is a cheap no-op when the pinned version is already built
lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	$(GOLANGCI_LINT) run ./...

test:
	go test ./... -race -cover

check: build vet fmt-check lint test
