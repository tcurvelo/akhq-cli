.PHONY: build test lint

# keep in sync with .github/workflows/ci.yml
GOLANGCI_LINT_VERSION := v2.15.0

build:
	go build -o akhq-cli .

test:
	go test ./...

lint:
	go vet ./...
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run ./...
