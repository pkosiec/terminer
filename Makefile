.DEFAULT_GOAL = all

export GO_LINTER_VERSION = v1.41.1

all: test lint build
.PHONY: all

install-dev-tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(shell go env GOPATH)/bin" ${GO_LINTER_VERSION}
.PHONY: install-dev-tools

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
.PHONY: test

lint:
	 golangci-lint run "./..."
.PHONY: lint

fix-lint-issues:
	golangci-lint run --fix "./..."
.PHONY: lint-fix

build:
	CGO_ENABLED=0 go build -o ./bin/terminer main.go
.PHONY: build
