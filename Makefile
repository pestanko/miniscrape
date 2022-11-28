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

.PHONY: help lint test install-tools build run-serve

all: help

## help: show help
help: Makefile
	@echo
	@echo " Choose a command run in:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: test
## test: run the tests in the Github pipeline
test: unit-test

.PHONY: unit-test
## unit-test: run the all unit tests
unit-test:
	@printf '\n$(YELLOW)Unit tests:$(NC) $(APPNAME) \n'
	mkdir -p reports/
	go test -v -coverprofile=reports/coverage.out $$(go list ./... | grep -v integration | grep -v mocks) 2>&1 | tee reports/unit.out
	cat reports/unit.out | go tool test2json > reports/unit.json
	go-junit-report -in reports/unit.out -set-exit-code > reports/unit.xml

## test-coverage: show test coverage
test-coverage: test
	go tool cover -func=reports/coverage.out
	go tool cover -html=reports/coverage.out

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
	rm bin/${BINARY_NAME}*

## run-serve: Run the server
run-serve:
	go run main.go serve