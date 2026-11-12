.PHONY: build test lint

build:
	go build -o akhq-cli .

test:
	go test ./...

lint:
	go vet ./...
	golangci-lint run ./...
