PROJECT_NAME := "apimsprdev"
PKG := "gitlab.com/Lapli/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all dep build clean test coverage coverhtml lint

all: build

lint: @golint -set_exit_status ${PKG_LIST}

test: @go test -short ${PKG_LIST}

race: dep  @go test -race -short ${PKG_LIST}

msan: dep @go test -msan -short ${PKG_LIST}

coverage: ./tools/coverage.sh;

coverhtml: ./tools/coverage.sh html;

dep: @go get -v -d ./...

build: dep @go build -i -v $(PKG)

clean: @rm -f $(PROJECT_NAME)

help: ## Display this help screen
  @grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'