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

BIN_DIR=./bin/
CLIENT_EXECUTABLE=gkcli
CLIENT_WINDOWS=$(BIN_DIR)$(CLIENT_EXECUTABLE)_windows_amd64.exe
CLIENT_LINUX=$(BIN_DIR)$(CLIENT_EXECUTABLE)_linux_amd64
CLIENT_DARWIN=$(BIN_DIR)$(CLIENT_EXECUTABLE)_darwin_amd64
#VERSION=$(shell git describe --tags --always --long --dirty)

GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILD_DATE := $(shell date +%FT%T%z)
VERSION := $(shell git describe --tags --abbrev=0 --always)

$(CLIENT_WINDOWS):
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o $(CLIENT_WINDOWS) \
		-ldflags="-X 'gophkeeper/pkg/version.Revision=$(GIT_COMMIT)'\
		 -X 'gophkeeper/pkg/version.Version=$(VERSION)'\
		 -X 'gophkeeper/pkg/version.Branch=$(GIT_BRANCH)'\
		 -X 'gophkeeper/pkg/version.BuildUser=$(USER)'\
		  -X 'gophkeeper/pkg/version.BuildDate=$(BUILD_DATE)'" \
		./cmd/client/*.go

$(CLIENT_LINUX):
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(CLIENT_LINUX) \
		-ldflags="-X 'gophkeeper/pkg/version.Revision=$(GIT_COMMIT)'\
		 -X 'gophkeeper/pkg/version.Version=$(VERSION)'\
		 -X 'gophkeeper/pkg/version.Branch=$(GIT_BRANCH)'\
		 -X 'gophkeeper/pkg/version.BuildUser=$(USER)'\
		  -X 'gophkeeper/pkg/version.BuildDate=$(BUILD_DATE)'" \
		./cmd/client/*.go

$(CLIENT_DARWIN):
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o $(CLIENT_DARWIN) \
		-ldflags="-X 'gophkeeper/pkg/version.Revision=$(GIT_COMMIT)'\
		 -X 'gophkeeper/pkg/version.Version=$(VERSION)'\
		 -X 'gophkeeper/pkg/version.Branch=$(GIT_BRANCH)'\
		 -X 'gophkeeper/pkg/version.BuildUser=$(USER)'\
		  -X 'gophkeeper/pkg/version.BuildDate=$(BUILD_DATE)'" \
		./cmd/client/*.go

build-client: $(CLIENT_LINUX) $(CLIENT_WINDOWS) $(CLIENT_DARWIN)
	@echo Version: $(VERSION)

build-server:
	@echo "Building the server app to the bin dir"
	CGO_ENABLED=1 go build -o ./bin/gk \
		-ldflags="-X 'gophkeeper/pkg/version.Revision=$(GIT_COMMIT)'\
		 -X 'gophkeeper/pkg/version.Version=$(VERSION)'\
		 -X 'gophkeeper/pkg/version.Branch=$(GIT_BRANCH)'\
		 -X 'gophkeeper/pkg/version.BuildUser=$(USER)'\
		  -X 'gophkeeper/pkg/version.BuildDate=$(BUILD_DATE)'" \
		./cmd/server/*.go

clean: ## Remove previous build
	rm -f $(CLIENT_LINUX) $(CLIENT_WINDOWS) $(CLIENT_DARWIN)