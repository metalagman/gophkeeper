lint:
	@echo "Running linter checks"
	golangci-lint run

test:
	@echo "Running UNIT tests"
	@go clean -testcache
	go test -cover -race -short ./... | { grep -v 'no test files'; true; }

cover-html:
	@echo "Running test coverage"
	@go clean -testcache
	go test -cover -coverprofile=coverage.out -race -short ./... | grep -v 'no test files'
	go tool cover -html=coverage.out

cover:
	@echo "Running test coverage"
	@go clean -testcache
	go test -cover -coverprofile=coverage.out -race -short ./... | grep -v 'no test files'
	go tool cover -func=coverage.out

generate:
	@echo "Generating mocks"
	go generate ./...

.PHONY: build
build: build-client build-server

EXECUTABLE=executable-name
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN=$(EXECUTABLE)_darwin_amd64
VERSION=$(shell git describe --tags --always --long --dirty)

GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILD_DATE := $(shell date +%FT%T%z)
#VERSION := $(shell git describe --tags --abbrev=0 --always)

build-client:
	@echo "Building the client app to the bin dir"
	@echo version: $(VERSION)
	CGO_ENABLED=1 go build -o ./bin/gkcli \
		-ldflags="-X 'gophkeeper/pkg/version.Revision=$(GIT_COMMIT)'\
		 -X 'gophkeeper/pkg/version.Version=$(VERSION)'\
		 -X 'gophkeeper/pkg/version.Branch=$(GIT_BRANCH)'\
		 -X 'gophkeeper/pkg/version.BuildUser=$(USER)'\
		  -X 'gophkeeper/pkg/version.BuildDate=$(BUILD_DATE)'" \
		./cmd/client/*.go

build-server:
	@echo "Building the server app to the bin dir"
	CGO_ENABLED=1 go build -o ./bin/gk \
		-ldflags="-X 'gophkeeper/pkg/version.Revision=$(GIT_COMMIT)'\
		 -X 'gophkeeper/pkg/version.Version=$(VERSION)'\
		 -X 'gophkeeper/pkg/version.Branch=$(GIT_BRANCH)'\
		 -X 'gophkeeper/pkg/version.BuildUser=$(USER)'\
		  -X 'gophkeeper/pkg/version.BuildDate=$(BUILD_DATE)'" \
		./cmd/server/*.go
