BINARY_NAME=anydev

.PHONY=generate

generate:
	@go generate ./...

.PHONY=run

run:
	@go run ./cmd

.PHONY=install

install:
	@go mod download
