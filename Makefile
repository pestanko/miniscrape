.DEFAULT_GOAL := lint
NC=\033[0m
YELLOW=\033[1;33m
GOVERSION=$(shell go version | awk  '{print $$3}' )
GOPATH=$(shell go env GOPATH)

BINARY_NAME=miniscrape

## lint: check all sources for errors
lint:
	@printf "\n$(YELLOW)linting$(NC)\n"
	revive -formatter friendly ./...
	gosec ./...

.PHONY: help lint test install-tools build

all: help

## help: show help
help: Makefile
	@echo
	@echo " Choose a command run in:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## test: run all unit tests
test:
	go test -coverprofile=coverage.out $$(go list ./... | grep -v -E '^github.com/CloudTalk-io/ctgo/services/|integration')

## test-coverage: show test coverage
test-coverage: test
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

## install-tools: install all golang tools
install-tools:
	go install github.com/mgechev/revive@latest
	go install golang.org/x/tools/gopls@latest
	curl -sfL \
		https://raw.githubusercontent.com/securego/gosec/master/install.sh | \
		sh -s -- -b $(GOPATH)/bin

## build: build all the binaries
build:
	mkdir -p bin/
	go build -o bin/${BINARY_NAME} main.go

## clean: clean all the binary files
clean:
	go clean
	rm bin/*
